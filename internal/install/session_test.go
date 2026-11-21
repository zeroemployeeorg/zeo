package install

import (
	"strings"
	"testing"
)

func validSession() SessionService {
	return SessionService{Organization: "example", Seat: "research", User: "seat_research", Group: "seat_research", PomarSHA256: strings.Repeat("a", 64), MemoryBytes: 1 << 30, Tasks: 128, CPUPercent: 100}
}

func TestSessionLifecycleAndCustody(t *testing.T) {
	unit, err := RenderSessionUnit(validSession())
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"Type=exec\n", "User=seat_research\n", "Group=seat_research\n",
		"StateDirectoryMode=0700\n", "UMask=0077\n", "NoNewPrivileges=true\n",
		"MemoryMax=1073741824\n", "TasksMax=128\n", "CPUQuota=100%\n",
		"ExecStart=/opt/zeo/artifacts/" + strings.Repeat("a", 64) + "/pomar session serve -root /var/lib/zeo-seats/example/research/sessions\n",
		"Environment=POMAR_SESSION_ROOT=/var/lib/zeo-seats/example/research/sessions\n",
		"Restart=on-failure\n", "KillMode=control-group\n", "WantedBy=multi-user.target\n",
	} {
		if !strings.Contains(unit, required) {
			t.Errorf("missing service boundary %q", required)
		}
	}
	// App and host-agent disruption must never propagate through a service
	// dependency. This is a packaging invariant, not a live survival proof.
	for _, forbidden := range []string{"PartOf=", "BindsTo=", "Requires=", "ExecStopPost=", "User=root", "/bin/sh", "tmux", "%i", "$HOME"} {
		if strings.Contains(unit, forbidden) {
			t.Errorf("unexpected lifecycle/authority coupling %q", forbidden)
		}
	}
}

func TestSessionRejectsUnsafeIdentityAndUnpinnedArtifact(t *testing.T) {
	cases := []struct {
		name string
		edit func(*SessionService)
	}{
		{"root-user", func(c *SessionService) { c.User = "root" }},
		{"root-group", func(c *SessionService) { c.Group = "root" }},
		{"numeric-uid", func(c *SessionService) { c.User = "0" }},
		{"directive-injection", func(c *SessionService) { c.Seat = "seat\nUser=root" }},
		{"path-traversal", func(c *SessionService) { c.Organization = "../other" }},
		{"specifier", func(c *SessionService) { c.User = "%i" }},
		{"expansion", func(c *SessionService) { c.Group = "$USER" }},
		{"empty-seat", func(c *SessionService) { c.Seat = "" }},
		{"long-identity", func(c *SessionService) { c.User = strings.Repeat("a", 33) }},
		{"unpinned", func(c *SessionService) { c.PomarSHA256 = "latest" }},
		{"digest-path", func(c *SessionService) { c.PomarSHA256 = strings.Repeat("a", 63) + "/" }},
		{"no-memory-limit", func(c *SessionService) { c.MemoryBytes = 0 }},
		{"no-process-limit", func(c *SessionService) { c.Tasks = 0 }},
		{"no-cpu-limit", func(c *SessionService) { c.CPUPercent = 0 }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := validSession()
			tc.edit(&c)
			unit, err := RenderSessionUnit(c)
			if err == nil || unit != "" {
				t.Fatal("unsafe configuration produced an installable-looking unit")
			}
		})
	}
}

func TestSessionStateSeparatedAcrossOrganizationsAndSeats(t *testing.T) {
	first, err := RenderSessionUnit(validSession())
	if err != nil {
		t.Fatal(err)
	}
	for _, change := range []func(*SessionService){
		func(c *SessionService) { c.Organization = "another" },
		func(c *SessionService) { c.Seat = "another" },
	} {
		c := validSession()
		change(&c)
		other, err := RenderSessionUnit(c)
		if err != nil {
			t.Fatal(err)
		}
		if other == first || strings.Contains(other, "/example/research/") {
			t.Fatal("separate identity reused the original state location")
		}
	}
}
