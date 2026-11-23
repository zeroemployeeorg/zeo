package install

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func stageFixture(t *testing.T) (*os.Root, string, string) {
	t.Helper()
	path := t.TempDir()
	if err := os.Chmod(path, 0o700); err != nil {
		t.Fatal(err)
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { root.Close() })
	payload := "synthetic test bytes; never execute"
	sum := sha256.Sum256([]byte(payload))
	return root, payload, hex.EncodeToString(sum[:])
}

func TestStagePublishesVerifiedBytesAndRerunPreservesObject(t *testing.T) {
	root, payload, digest := stageFixture(t)
	path, err := StagePomarArtifact(root, strings.NewReader(payload), digest, int64(len(payload)))
	if err != nil {
		t.Fatal(err)
	}
	before, err := root.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := StagePomarArtifact(root, strings.NewReader(payload), digest, int64(len(payload))); err != nil {
		t.Fatal(err)
	}
	after, err := root.Stat(path)
	if err != nil || !os.SameFile(before, after) || after.Mode().Perm() != 0o555 {
		t.Fatal("rerun replaced object or failed executable custody")
	}
	data, err := root.ReadFile(path)
	if err != nil || string(data) != payload {
		t.Fatal("published bytes differ")
	}
}

func TestStageRejectsBadInputsWithoutPublication(t *testing.T) {
	for _, kind := range []string{"digest", "short", "long", "read-error"} {
		t.Run(kind, func(t *testing.T) {
			root, payload, digest := stageFixture(t)
			size := int64(len(payload))
			if kind == "digest" {
				digest = strings.Repeat("0", 64)
			}
			if kind == "short" {
				size++
			}
			if kind == "long" {
				size--
			}
			var err error
			if kind == "read-error" {
				_, err = StagePomarArtifact(root, brokenArtifactReader{}, digest, size)
			} else {
				_, err = StagePomarArtifact(root, strings.NewReader(payload), digest, size)
			}
			if err == nil {
				t.Fatal("invalid source accepted")
			}
			entries, err := os.ReadDir(root.Name())
			if err != nil || len(entries) != 0 {
				t.Fatal("failed copy published or left staging files")
			}
		})
	}
}

type brokenArtifactReader struct{}

func (brokenArtifactReader) Read([]byte) (int, error) { return 0, errors.New("source interrupted") }

func TestStagePreservesConflictingExistingArtifact(t *testing.T) {
	root, payload, digest := stageFixture(t)
	if err := root.Mkdir(digest, 0o755); err != nil {
		t.Fatal(err)
	}
	path := digest + "/pomar"
	if err := root.WriteFile(path, []byte("existing different object"), 0o555); err != nil {
		t.Fatal(err)
	}
	if _, err := StagePomarArtifact(root, strings.NewReader(payload), digest, int64(len(payload))); err == nil {
		t.Fatal("existing conflict accepted")
	}
	data, err := root.ReadFile(path)
	if err != nil || string(data) != "existing different object" {
		t.Fatal("conflicting object was overwritten")
	}
}

func TestStageRefusesSymlinkEscapeAndWritableRoot(t *testing.T) {
	root, payload, digest := stageFixture(t)
	outside := t.TempDir()
	if err := root.Symlink(outside, digest); err != nil {
		t.Fatal(err)
	}
	if _, err := StagePomarArtifact(root, strings.NewReader(payload), digest, int64(len(payload))); err == nil {
		t.Fatal("symlink target accepted")
	}
	if _, err := os.Stat(filepath.Join(outside, "pomar")); !os.IsNotExist(err) {
		t.Fatal("artifact escaped root")
	}
	if err := root.Remove(digest); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(root.Name(), 0o777); err != nil {
		t.Fatal(err)
	}
	if _, err := StagePomarArtifact(root, strings.NewReader(payload), digest, int64(len(payload))); err == nil {
		t.Fatal("world-writable staging root accepted")
	}
}

func TestConcurrentMatchingStagesConverge(t *testing.T) {
	root, payload, digest := stageFixture(t)
	results := make(chan error, 2)
	for range 2 {
		go func() {
			_, err := StagePomarArtifact(root, strings.NewReader(payload), digest, int64(len(payload)))
			results <- err
		}()
	}
	for range 2 {
		if err := <-results; err != nil {
			t.Fatal(err)
		}
	}
	data, err := root.ReadFile(digest + "/pomar")
	if err != nil || string(data) != payload {
		t.Fatal("concurrent publication did not preserve pinned bytes")
	}
}

func TestStageRefusesSetIDExistingArtifact(t *testing.T) {
	root, payload, digest := stageFixture(t)
	path, err := StagePomarArtifact(root, strings.NewReader(payload), digest, int64(len(payload)))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(filepath.Join(root.Name(), path), os.ModeSetuid|0o555); err != nil {
		t.Fatal(err)
	}
	if _, err := StagePomarArtifact(root, strings.NewReader(payload), digest, int64(len(payload))); err == nil {
		t.Fatal("existing privileged executable accepted")
	}
}
