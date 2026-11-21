package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestStageCommandCopiesOnlyPinnedLocalArtifact(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chmod(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(t.TempDir(), "pomar")
	data := []byte("synthetic fixture, not an executable")
	if err := os.WriteFile(source, data, 0o600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	digest := hex.EncodeToString(sum[:])
	var out, stderr bytes.Buffer
	code := run([]string{"--pomar", source, "--sha256", digest, "--size", strconv.Itoa(len(data)), "--artifact-root", dir}, &out, &stderr)
	if code != 0 || out.String() != digest+"/pomar\n" || stderr.Len() != 0 {
		t.Fatalf("staging command failed: %d %s", code, stderr.String())
	}
	staged, err := os.ReadFile(filepath.Join(dir, digest, "pomar"))
	if err != nil || !bytes.Equal(staged, data) {
		t.Fatal("staged bytes differ")
	}
}

func TestStageCommandRefusesSourceSymlink(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alias")
	if err := os.Symlink(filepath.Join(dir, "missing"), path); err != nil {
		t.Fatal(err)
	}
	var out, stderr bytes.Buffer
	if run([]string{"--pomar", path, "--artifact-root", dir}, &out, &stderr) == 0 || out.Len() != 0 {
		t.Fatal("source link accepted or success output emitted")
	}
}
