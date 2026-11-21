package install

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
)

// StagePomarArtifact copies a release-pinned executable into a trusted artifact
// root as <sha256>/pomar. Publication is atomic and never replaces an existing
// object. The caller owns root and must keep other writers out of it. This is
// staging only: target, provenance and owner qualification precede activation.
func StagePomarArtifact(root *os.Root, source io.Reader, digest string, size int64) (string, error) {
	if !pin.MatchString(digest) || size <= 0 || size > 1<<30 {
		return "", fmt.Errorf("valid digest and positive size up to 1 GiB required")
	}
	info, err := root.Stat(".")
	if err != nil || !trusted(info, true) {
		return "", fmt.Errorf("artifact root must be owned by the caller and not writable by group or others")
	}
	var nonce [16]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", err
	}
	temp := ".stage-pomar-" + hex.EncodeToString(nonce[:])
	file, err := root.OpenFile(temp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", err
	}
	defer root.Remove(temp)
	defer file.Close()
	hash := sha256.New()
	n, err := io.Copy(io.MultiWriter(file, hash), io.LimitReader(source, size+1))
	if err != nil || n != size || hex.EncodeToString(hash.Sum(nil)) != digest {
		return "", fmt.Errorf("artifact copy failed size or SHA256 verification")
	}
	if err := file.Chmod(0o555); err != nil {
		return "", err
	}
	if err := file.Sync(); err != nil {
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	if err := root.Mkdir(digest, 0o755); err != nil && !errors.Is(err, os.ErrExist) {
		return "", err
	}
	// Root confines all paths; refuse alias/custody conflicts even within it.
	info, err = root.Lstat(digest)
	if err != nil || !trusted(info, true) {
		return "", fmt.Errorf("artifact directory has conflicting custody")
	}
	dir, err := root.OpenRoot(digest)
	if err != nil {
		return "", err
	}
	defer dir.Close()
	// Link is an atomic, no-replace publication. Concurrent matching stages
	// converge; a conflicting existing object is checked, never overwritten.
	if err := root.Link(temp, digest+"/pomar"); err != nil && !errors.Is(err, os.ErrExist) {
		return "", err
	}
	if err := verifyStagedPomar(dir, digest, size); err != nil {
		return "", err
	}
	if err := syncDirectory(dir); err != nil {
		return "", err
	}
	if err := syncDirectory(root); err != nil {
		return "", err
	}
	return digest + "/pomar", nil
}

func trusted(info os.FileInfo, directory bool) bool {
	if info == nil || info.Mode()&os.ModeSymlink != 0 || info.Mode().Perm()&0o022 != 0 || info.Mode()&(os.ModeSetuid|os.ModeSetgid|os.ModeSticky) != 0 {
		return false
	}
	if directory && !info.IsDir() || !directory && !info.Mode().IsRegular() {
		return false
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	return ok && st.Uid == uint32(os.Geteuid())
}

func verifyStagedPomar(dir *os.Root, digest string, size int64) error {
	before, err := dir.Lstat("pomar")
	if err != nil || !trusted(before, false) || before.Mode().Perm() != 0o555 {
		return fmt.Errorf("existing artifact is not a caller-owned read/execute-only regular executable")
	}
	file, err := dir.OpenFile("pomar", os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !os.SameFile(before, opened) || opened.Size() != size {
		return fmt.Errorf("artifact identity or size changed")
	}
	hash := sha256.New()
	n, err := io.Copy(hash, io.LimitReader(file, size+1))
	after, statErr := dir.Lstat("pomar")
	if err != nil || statErr != nil || n != size || !os.SameFile(opened, after) || !trusted(after, false) || after.Mode().Perm() != 0o555 || hex.EncodeToString(hash.Sum(nil)) != digest {
		return fmt.Errorf("existing artifact conflicts with release pin")
	}
	return nil
}

func syncDirectory(root *os.Root) error {
	file, err := root.Open(".")
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Sync()
}
