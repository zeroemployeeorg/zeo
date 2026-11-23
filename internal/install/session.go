// Package install contains distribution service packaging. Rendering a unit
// does not install an artifact, allocate an identity or qualify an environment.
package install

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	identity = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)
	pin      = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

// SessionService binds an already allocated Runtime identity to a pinned
// Pomar executable. Limits are mandatory profile inputs, not product minima.
// Each organization/seat pair requires its own UID/GID; the caller must check
// this globally before activation. Names alone cannot enforce that invariant.
type SessionService struct {
	Organization string
	Seat         string
	User         string
	Group        string
	PomarSHA256  string
	MemoryBytes  uint64
	Tasks        uint64
	CPUPercent   uint64
}

// RenderSessionUnit renders a SYSTEM unit for installation INSIDE the seat's
// Linux environment. No host-agent or App unit owns its lifecycle. The keeper's
// own restart still loses PTYs; systemd supervision is not session recovery.
func RenderSessionUnit(c SessionService) (string, error) {
	for _, value := range []string{c.Organization, c.Seat, c.User, c.Group} {
		if !identity.MatchString(value) || value == "root" {
			return "", fmt.Errorf("organization, seat and non-root identities require lowercase safe names")
		}
	}
	if !pin.MatchString(c.PomarSHA256) {
		return "", fmt.Errorf("Pomar artifact requires a lowercase SHA256 pin")
	}
	if c.MemoryBytes == 0 || c.Tasks == 0 || c.CPUPercent == 0 {
		return "", fmt.Errorf("memory, process and CPU limits must be explicitly positive")
	}
	base := "zeo-seats/" + c.Organization + "/" + c.Seat
	home := "/var/lib/" + base + "/home"
	workspace := "/var/lib/" + base + "/workspace"
	sessions := "/var/lib/" + base + "/sessions"
	binary := "/opt/zeo/artifacts/" + c.PomarSHA256 + "/pomar"
	// All substituted strings have a restricted alphabet; no shell, systemd
	// specifiers, environment expansion or configuration line injection occurs.
	return strings.Join([]string{
		"# Candidate guest-local session keeper; installed artifact/custody must be qualified.",
		"# Restart of THIS unit loses live sessions; App/host-agent restart must not restart it.",
		"[Unit]",
		"Description=ZEO guest session keeper",
		"After=local-fs.target",
		"",
		"[Service]",
		"Type=exec",
		"User=" + c.User,
		"Group=" + c.Group,
		"WorkingDirectory=" + workspace,
		"Environment=HOME=" + home,
		"Environment=POMAR_SESSION_ROOT=" + sessions,
		"StateDirectory=" + base + "/home " + base + "/workspace " + base + "/sessions",
		"StateDirectoryMode=0700",
		"ExecStart=" + binary + " session serve -root " + sessions,
		"Restart=on-failure",
		"RestartSec=3s",
		"TimeoutStopSec=30s",
		"KillMode=control-group",
		"UMask=0077",
		"NoNewPrivileges=true",
		"PrivateTmp=true",
		"ProtectSystem=strict",
		"ProtectHome=read-only",
		"ReadWritePaths=" + home + " " + workspace + " " + sessions,
		fmt.Sprintf("MemoryMax=%d", c.MemoryBytes),
		fmt.Sprintf("TasksMax=%d", c.Tasks),
		fmt.Sprintf("CPUQuota=%d%%", c.CPUPercent),
		"StandardOutput=journal",
		"StandardError=journal",
		"",
		"[Install]",
		"WantedBy=multi-user.target",
		"",
	}, "\n"), nil
}
