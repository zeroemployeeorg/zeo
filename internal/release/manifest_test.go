package release

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// All versions, URLs and payloads below are synthetic regression fixtures.
// They are not an available release, install artifact or acceptance receipt.
func fixture() Manifest {
	m := Manifest{
		Schema: "zeo.release/v1", Version: "0.1.0",
		Platform:    Platform{OS: "linux", Architecture: "amd64", LinuxID: "test-linux", LinuxVersion: "1"},
		Environment: Environment{Profile: "test-profile", Version: "1"},
		UpgradeFrom: []string{}, AgentTools: []VersionPin{},
		Provenance: Provenance{Recipe: "test-only", SourceRepository: "https://example.invalid/distribution.git", SourceRevision: strings.Repeat("a", 40), InputsSHA256: strings.Repeat("0", 64)},
	}
	for _, capability := range requiredCapabilities {
		m.Components = append(m.Components, Component{
			Name: capability, Version: "1", SourceRepository: "https://example.invalid/" + capability + ".git", SourceRevision: strings.Repeat("b", 40),
			Provides: []string{capability}, Requires: []string{}, InstallAfter: []string{},
			Contracts:    []VersionPin{{Name: capability + "-api", Version: "1"}},
			StateSchemas: []StateSchema{}, RuntimeDependencies: []VersionPin{},
			Artifact: Artifact{Filename: capability + ".bin", URL: "https://example.invalid/" + capability + ".bin", SHA256: strings.Repeat("0", 64), SizeBytes: 1},
		})
	}
	return m
}

func target(m Manifest) Target { return Target{Platform: m.Platform, Environment: m.Environment} }

func TestFullCapabilitySetAndDeterministicOrdering(t *testing.T) {
	m := fixture()
	// Runtime/messaging integration cycles do not imply an installation cycle.
	m.Components[0].Requires = []string{"messaging"}
	m.Components[2].Requires = []string{"runtime"}
	plan, err := m.ValidateFor(target(m))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"agenda", "application", "controlling", "coordination", "identity", "messaging", "runtime"}
	if !reflect.DeepEqual(plan.InstallOrder, want) {
		t.Fatalf("order %v", plan.InstallOrder)
	}
	if plan.ReadyForInstall || plan.ArtifactsVerified {
		t.Fatal("metadata validation promoted itself to installation or artifact proof")
	}
	for i := 0; i < 30; i++ {
		again, err := m.ValidateFor(target(m))
		if err != nil || !reflect.DeepEqual(again, plan) {
			t.Fatal("plan is not deterministic")
		}
	}
}

func TestInstallOrderingIsSeparateAndHonored(t *testing.T) {
	m := fixture()
	m.Components[0].InstallAfter = []string{"messaging"}
	plan, err := m.ValidateFor(target(m))
	if err != nil {
		t.Fatal(err)
	}
	if plan.InstallOrder[len(plan.InstallOrder)-1] != "runtime" {
		t.Fatal(plan.InstallOrder)
	}
	m.Components[2].InstallAfter = []string{"runtime"}
	if _, err := m.ValidateFor(target(m)); err == nil {
		t.Fatal("accepted an installation cycle")
	}
}

func TestPlanBindsProfilePlatformAndMetadata(t *testing.T) {
	m := fixture()
	plan, err := m.ValidateFor(target(m))
	if err != nil {
		t.Fatal(err)
	}
	if plan.Platform != m.Platform || plan.Environment != m.Environment || len(plan.MetadataSHA256) != 64 {
		t.Fatal("plan omitted target or metadata binding")
	}
	m.Components[0].Contracts[0].Version = "2"
	changed, err := m.ValidateFor(target(m))
	if err != nil {
		t.Fatal(err)
	}
	if changed.MetadataSHA256 == plan.MetadataSHA256 {
		t.Fatal("changed contract retained old metadata binding")
	}
}

func TestContractAndPackageIdentifiersKeepTheirNamespaces(t *testing.T) {
	m := fixture()
	m.Components[0].Contracts = []VersionPin{{Name: "zeo.runtime.api", Version: "v1"}}
	m.Components[1].Contracts = []VersionPin{{Name: "zeo.coordination/api", Version: "v1"}}
	m.Components[0].RuntimeDependencies = []VersionPin{{Name: "libstdc++6", Version: "14"}}
	m.Components[0].StateSchemas = []StateSchema{{Name: "zeo.runtime/events", Version: "1", CompatibleFrom: []string{}}}
	if _, err := m.ValidateFor(target(m)); err != nil {
		t.Fatal("valid owner contract or package namespace was lost", err)
	}
}

func TestRejectIncompleteOrConflictingRelease(t *testing.T) {
	cases := map[string]func(*Manifest){
		"missing capability":            func(m *Manifest) { m.Components = m.Components[1:] },
		"duplicate authority":           func(m *Manifest) { m.Components[1].Provides = []string{"runtime"} },
		"duplicate component":           func(m *Manifest) { m.Components[1].Name = m.Components[0].Name },
		"unknown required component":    func(m *Manifest) { m.Components[0].Requires = []string{"absent"} },
		"unknown install predecessor":   func(m *Manifest) { m.Components[0].InstallAfter = []string{"absent"} },
		"self dependency":               func(m *Manifest) { m.Components[0].Requires = []string{"runtime"} },
		"implicit runtime requirements": func(m *Manifest) { m.Components[0].RuntimeDependencies = nil },
		"missing source pin":            func(m *Manifest) { m.Components[0].SourceRevision = "main" },
		"unqualified source repository": func(m *Manifest) { m.Components[0].SourceRepository = "" },
		"conflicting contracts":         func(m *Manifest) { m.Components[1].Contracts = []VersionPin{{Name: "runtime-api", Version: "2"}} },
		"conflicting runtime dependencies": func(m *Manifest) {
			m.Components[0].RuntimeDependencies = []VersionPin{{Name: "node", Version: "1"}}
			m.Components[1].RuntimeDependencies = []VersionPin{{Name: "node", Version: "2"}}
		},
		"implicit schema compatibility": func(m *Manifest) { m.Components[0].StateSchemas = []StateSchema{{Name: "state", Version: "1"}} },
		"implicit upgrade paths":        func(m *Manifest) { m.UpgradeFrom = nil },
		"duplicate upgrade path":        func(m *Manifest) { m.UpgradeFrom = []string{"0.0.1", "0.0.1"} },
		"invalid architecture":          func(m *Manifest) { m.Platform.Architecture = "unknown" },
		"non Linux platform":            func(m *Manifest) { m.Platform.OS = "darwin" },
		"missing provenance":            func(m *Manifest) { m.Provenance.InputsSHA256 = "" },
		"short digest":                  func(m *Manifest) { m.Components[0].Artifact.SHA256 = "123" },
		"overflowing size":              func(m *Manifest) { m.Components[0].Artifact.SizeBytes = 1<<63 - 1 },
		"zero size":                     func(m *Manifest) { m.Components[0].Artifact.SizeBytes = 0 },
		"path traversal":                func(m *Manifest) { m.Components[0].Artifact.Filename = "../elsewhere" },
		"duplicate filename":            func(m *Manifest) { m.Components[1].Artifact.Filename = m.Components[0].Artifact.Filename },
		"HTTP artifact":                 func(m *Manifest) { m.Components[0].Artifact.URL = "http://example.invalid/file" },
		"embedded credential":           func(m *Manifest) { m.Components[0].Artifact.URL = "https://user:PRIVATE_MARKER@example.invalid/file" },
		"query credential":              func(m *Manifest) { m.Components[0].Artifact.URL = "https://example.invalid/file?token=PRIVATE_MARKER" },
		"loopback":                      func(m *Manifest) { m.Components[0].Artifact.URL = "https://127.0.0.1/file" },
		"private address":               func(m *Manifest) { m.Components[0].Artifact.URL = "https://10.1.2.3/file" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			m := fixture()
			mutate(&m)
			_, err := m.ValidateFor(target(m))
			if err == nil {
				t.Fatal("accepted invalid release")
			}
			if strings.Contains(err.Error(), "PRIVATE_MARKER") {
				t.Fatal("error disclosed an input value")
			}
		})
	}
}

func TestTargetMustMatchEveryDimension(t *testing.T) {
	m := fixture()
	for _, field := range []string{"architecture", "distribution", "linux-version", "profile", "profile-version"} {
		t.Run(field, func(t *testing.T) {
			other := target(m)
			switch field {
			case "architecture":
				other.Platform.Architecture = "arm64"
			case "distribution":
				other.Platform.LinuxID = "different"
			case "linux-version":
				other.Platform.LinuxVersion = "2"
			case "profile":
				other.Environment.Profile = "different"
			case "profile-version":
				other.Environment.Version = "2"
			}
			if _, err := m.ValidateFor(other); err == nil {
				t.Fatal("accepted a mismatched target")
			}
		})
	}
}

func TestDecodeRejectsAmbiguityUnknownFieldsAndOversize(t *testing.T) {
	raw, err := json.Marshal(fixture())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Decode(bytes.NewReader(raw)); err != nil {
		t.Fatal(err)
	}
	cases := [][]byte{
		[]byte(`{"schema":"a","schema":"b"}`),
		[]byte(`{"schema":"zeo.release/v1","Schema":"other"}`),
		[]byte(`{"schema":"zeo.release/v1","ſchema":"other"}`),
		[]byte(`{"private_marker":"PRIVATE_MARKER"}`),
		[]byte(`{"components":[{"artifact":{"size_bytes":1,"size_bytes":2}}]}`),
		[]byte(`{"PRIVATE_MARKER":"PRIVATE_MARKER"}`),
		append(append([]byte{}, raw...), []byte(` {}`)...),
		[]byte(strings.Repeat("[", 34) + "0" + strings.Repeat("]", 34)),
		[]byte(strings.Repeat(" ", MaxManifestBytes+1)),
		[]byte(`{"components":[{"artifact":{"size_bytes":true}}]}`),
	}
	for i, data := range cases {
		_, err := Decode(bytes.NewReader(data))
		if err == nil {
			t.Fatalf("accepted invalid document %d", i)
		}
		if strings.Contains(err.Error(), "PRIVATE_MARKER") {
			t.Fatal("decoder error disclosed untrusted fields")
		}
	}
}

func artifactFixture(t *testing.T) (Manifest, string) {
	t.Helper()
	m, dir := fixture(), t.TempDir()
	for i := range m.Components {
		data := []byte("synthetic artifact for " + m.Components[i].Name)
		hash := sha256.Sum256(data)
		m.Components[i].Artifact.SHA256 = hex.EncodeToString(hash[:])
		m.Components[i].Artifact.SizeBytes = int64(len(data))
		if err := os.WriteFile(filepath.Join(dir, m.Components[i].Artifact.Filename), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	return m, dir
}

func TestArtifactBytesAndSizeMustMatch(t *testing.T) {
	m, dir := artifactFixture(t)
	if err := m.VerifyArtifacts(dir); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, m.Components[0].Artifact.Filename)
	if err := os.WriteFile(file, bytes.Repeat([]byte("x"), int(m.Components[0].Artifact.SizeBytes)), 0600); err != nil {
		t.Fatal(err)
	}
	if err := m.VerifyArtifacts(dir); err == nil {
		t.Fatal("accepted changed artifact with unchanged size")
	}
}

func TestArtifactSizeMustMatchIndependentlyOfDigest(t *testing.T) {
	m, dir := artifactFixture(t)
	m.Components[0].Artifact.SizeBytes++
	if err := m.VerifyArtifacts(dir); err == nil {
		t.Fatal("accepted incorrect declared size")
	}
}

func TestArtifactSymlinkAndUnvalidatedPathAreRefused(t *testing.T) {
	m, dir := artifactFixture(t)
	file := filepath.Join(dir, m.Components[0].Artifact.Filename)
	if err := os.Rename(file, file+".original"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(file+".original", file); err != nil {
		t.Fatal(err)
	}
	if err := m.VerifyArtifacts(dir); err == nil {
		t.Fatal("accepted artifact symlink")
	}
	m.Components[0].Artifact.Filename = "../unrelated"
	if err := m.VerifyArtifacts(dir); err == nil {
		t.Fatal("accepted an unchecked file path")
	}
}
