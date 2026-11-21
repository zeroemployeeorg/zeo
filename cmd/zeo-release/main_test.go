package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zeroemployeeorg/zeo/internal/release"
)

func validInputs(t *testing.T) (string, string, []string) {
	t.Helper()
	dir := t.TempDir()
	payload := []byte("synthetic release payload")
	hash := sha256.Sum256(payload)
	m := release.Manifest{
		Schema: "zeo.release/v1", Version: "0.0.0-test",
		Platform:    release.Platform{OS: "linux", Architecture: "amd64", LinuxID: "test-linux", LinuxVersion: "1"},
		Environment: release.Environment{Profile: "test-profile", Version: "1"},
		UpgradeFrom: []string{}, AgentTools: []release.VersionPin{},
		Provenance: release.Provenance{Recipe: "test-only", SourceRepository: "https://example.invalid/distribution.git", SourceRevision: strings.Repeat("a", 40), InputsSHA256: strings.Repeat("0", 64)},
		Components: []release.Component{{
			Name: "platform", Version: "1", SourceRepository: "https://example.invalid/platform.git", SourceRevision: strings.Repeat("b", 40),
			Provides: []string{"runtime", "coordination", "messaging", "agenda", "controlling", "identity", "application"},
			Requires: []string{}, InstallAfter: []string{}, Contracts: []release.VersionPin{{Name: "test-api", Version: "1"}},
			StateSchemas: []release.StateSchema{}, RuntimeDependencies: []release.VersionPin{},
			Artifact: release.Artifact{Filename: "test.bin", URL: "https://example.invalid/test.bin", SHA256: hex.EncodeToString(hash[:]), SizeBytes: int64(len(payload))},
		}},
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(dir, "release.json")
	if err := os.WriteFile(manifest, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "test.bin"), payload, 0600); err != nil {
		t.Fatal(err)
	}
	args := []string{"--manifest", manifest, "--arch", "amd64", "--linux-id", "test-linux", "--linux-version", "1", "--profile", "test-profile", "--profile-version", "1", "--artifacts", dir}
	return manifest, dir, args
}

func TestCLIProducesVerifiedArtifactPlanWithoutInstallationClaim(t *testing.T) {
	_, _, args := validInputs(t)
	var stdout, stderr bytes.Buffer
	if exit := run(args, &stdout, &stderr); exit != 0 {
		t.Fatal("valid preflight failed", stderr.String())
	}
	var plan release.Plan
	if err := json.Unmarshal(stdout.Bytes(), &plan); err != nil {
		t.Fatal(err)
	}
	if !plan.ArtifactsVerified || plan.ReadyForInstall || len(plan.InstallOrder) != 1 || plan.InstallOrder[0] != "platform" {
		t.Fatal("incorrect evidence boundary in plan")
	}
}

func TestCLIRefusesToReplaceManifestOrArtifactInputs(t *testing.T) {
	for _, choice := range []string{"manifest", "artifact", "manifest-alias"} {
		t.Run(choice, func(t *testing.T) {
			manifest, dir, args := validInputs(t)
			output := manifest
			if choice == "artifact" {
				output = filepath.Join(dir, "test.bin")
			}
			if choice == "manifest-alias" {
				output = filepath.Join(dir, "alias.json")
				if err := os.Link(manifest, output); err != nil {
					t.Fatal(err)
				}
			}
			prior, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			if exit := run(append(args, "--out", output), &stdout, &stderr); exit != 1 || stdout.Len() != 0 {
				t.Fatal("accepted output over release input")
			}
			after, err := os.ReadFile(output)
			if err != nil || !bytes.Equal(prior, after) {
				t.Fatal("release input was changed")
			}
		})
	}
}

func TestRejectedManifestLeavesOutputUntouched(t *testing.T) {
	dir := t.TempDir()
	manifest, output := filepath.Join(dir, "release.json"), filepath.Join(dir, "plan.json")
	if err := os.WriteFile(manifest, []byte(`{"PRIVATE_MARKER":"PRIVATE_MARKER"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(output, []byte("prior plan"), 0600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	exit := run([]string{"--manifest", manifest, "--arch", "amd64", "--linux-id", "test-linux", "--linux-version", "1", "--profile", "test-profile", "--profile-version", "1", "--out", output}, &stdout, &stderr)
	if exit != 1 || stdout.Len() != 0 {
		t.Fatal("failure did not refuse cleanly")
	}
	if strings.Contains(stderr.String(), "PRIVATE_MARKER") {
		t.Fatal("error disclosed an input value")
	}
	data, err := os.ReadFile(output)
	if err != nil || string(data) != "prior plan" {
		t.Fatal("failure replaced existing output")
	}
}

func TestHelpAndUsageDoNotOpenManifest(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if exit := run([]string{"--help"}, &stdout, &stderr); exit != 0 || !strings.Contains(stdout.String(), "--manifest") {
		t.Fatal("help unavailable")
	}
	stdout.Reset()
	stderr.Reset()
	if exit := run([]string{"--unknown=PRIVATE_MARKER"}, &stdout, &stderr); exit != 2 || stdout.Len() != 0 {
		t.Fatal("usage failure was not bounded")
	}
	if strings.Contains(stderr.String(), "PRIVATE_MARKER") {
		t.Fatal("usage echoed an untrusted value")
	}
}

func TestAtomicPlanReplacementIsRepeatable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plan.json")
	if err := os.WriteFile(path, []byte("previous"), 0600); err != nil {
		t.Fatal(err)
	}
	data := []byte("{\"ready_for_install\":false}\n")
	for i := 0; i < 2; i++ {
		if err := writePlan(path, data); err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(got, data) {
			t.Fatal("plan replacement was not exact and repeatable")
		}
	}
}
