package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type Paths struct {
	Base, Home, Journal, Sudoers, Anchor string
	Owner                                uint32
}

func nativePaths() Paths {
	return Paths{"/usr/local/lib/zeo-ssh-target", "/var/lib/zeoconnect-broker-run", "/var/lib/zeo-ssh-target-bootstrap", "/etc/sudoers.d/zeo-ssh-target", "/", 0}
}
func (p Paths) key() string  { return filepath.Join(p.Home, ".ssh/authorized_keys") }
func (p Paths) tool() string { return filepath.Join(p.Base, "zeo-target-bootstrap") }

func sameStat(a, b *syscall.Stat_t) bool {
	return a.Dev == b.Dev && a.Ino == b.Ino && a.Mode == b.Mode && a.Uid == b.Uid && a.Gid == b.Gid && a.Nlink == b.Nlink && a.Size == b.Size && a.Mtim == b.Mtim && a.Ctim == b.Ctim
}

func readStable(path string, limit int64, owner uint32, mode os.FileMode) ([]byte, error) {
	f, e := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if e != nil {
		return nil, errRefused
	}
	defer f.Close()
	a, e := f.Stat()
	if e != nil {
		return nil, errRefused
	}
	s, ok := a.Sys().(*syscall.Stat_t)
	if !ok || !a.Mode().IsRegular() || s.Uid != owner || s.Nlink != 1 || (owner == 0 && s.Gid != 0) || a.Mode().Perm()&0022 != 0 || a.Size() > limit || (mode != 0 && a.Mode().Perm() != mode) {
		return nil, errRefused
	}
	raw, e := io.ReadAll(io.LimitReader(f, limit+1))
	if e != nil || int64(len(raw)) != a.Size() {
		return nil, errRefused
	}
	b, e := f.Stat()
	if e != nil || !sameStat(s, b.Sys().(*syscall.Stat_t)) {
		return nil, errRefused
	}
	return raw, nil
}
func ioReadBounded(r io.Reader, limit int64) ([]byte, error) {
	b, e := io.ReadAll(io.LimitReader(r, limit+1))
	if e != nil || int64(len(b)) > limit {
		return nil, errRefused
	}
	return b, nil
}

func (p Paths) trustedDir(path string) error {
	path = filepath.Clean(path)
	anchor := filepath.Clean(p.Anchor)
	if path != anchor && !strings.HasPrefix(path, anchor+string(os.PathSeparator)) && anchor != "/" {
		return errRefused
	}
	for {
		s, e := os.Lstat(path)
		if e != nil {
			return errRefused
		}
		st, ok := s.Sys().(*syscall.Stat_t)
		if !ok || !s.IsDir() || st.Uid != p.Owner || s.Mode().Perm()&0022 != 0 {
			return errRefused
		}
		buf := make([]byte, 8192)
		n, e := syscall.Listxattr(path, buf)
		if e != nil {
			return errRefused
		}
		if strings.Contains(string(buf[:n]), "system.posix_acl_") {
			return errRefused
		}
		if path == anchor {
			return nil
		}
		next := filepath.Dir(path)
		if next == path {
			return errRefused
		}
		path = next
	}
}

func syncDir(path string) error {
	f, e := os.OpenFile(path, os.O_RDONLY|syscall.O_DIRECTORY|syscall.O_NOFOLLOW, 0)
	if e != nil {
		return e
	}
	defer f.Close()
	return f.Sync()
}
func exists(path string) bool { _, e := os.Lstat(path); return !os.IsNotExist(e) }
func createFile(path string, raw []byte, mode os.FileMode) error {
	f, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL|syscall.O_NOFOLLOW, mode)
	if e != nil {
		return e
	}
	_, writeErr := f.Write(raw)
	if writeErr == nil {
		writeErr = f.Chmod(mode)
	}
	if writeErr == nil {
		writeErr = f.Sync()
	}
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return closeErr
	}
	return syncDir(filepath.Dir(path))
}
func writeRecord(path string, value any) error {
	raw, e := json.Marshal(value)
	if e != nil {
		return e
	}
	return createFile(path, append(raw, '\n'), 0600)
}
func mkdirNew(path string, mode os.FileMode) error {
	if e := os.Mkdir(path, mode); e != nil {
		return e
	}
	if e := os.Chmod(path, mode); e != nil {
		return e
	}
	return syncDir(filepath.Dir(path))
}

func lockCommand(p Paths) (*os.File, error) {
	path := filepath.Join(p.Base, "command.lock")
	f, e := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if e != nil {
		return nil, e
	}
	info, e := f.Stat()
	if e != nil {
		f.Close()
		return nil, e
	}
	s := info.Sys().(*syscall.Stat_t)
	if !info.Mode().IsRegular() || s.Uid != p.Owner || (p.Owner == 0 && s.Gid != 0) || s.Nlink != 1 || info.Size() != 0 || info.Mode().Perm() != 0644 {
		f.Close()
		return nil, errRefused
	}
	if e = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
		f.Close()
		return nil, errRefused
	}
	return f, nil
}

func verifyLinkedAuthorization(p Paths, expected string) error {
	return verifyLinkedFile(p.key(), expected, 0644, p.Owner)
}
func verifyLinkedFile(path string, expected string, mode os.FileMode, owner uint32) error {
	f, e := os.OpenFile(path+".revoked", os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if e != nil {
		return errRefused
	}
	defer f.Close()
	a, e := f.Stat()
	if e != nil {
		return errRefused
	}
	st := a.Sys().(*syscall.Stat_t)
	if !a.Mode().IsRegular() || st.Uid != owner || (owner == 0 && st.Gid != 0) || st.Nlink != 2 || a.Mode().Perm() != mode || a.Size() > maxPacket {
		return errRefused
	}
	raw, e := ioReadBounded(f, maxPacket)
	if e != nil || int64(len(raw)) != a.Size() || hash(raw) != expected {
		return errRefused
	}
	b, e := f.Stat()
	if e != nil || !sameStat(st, b.Sys().(*syscall.Stat_t)) {
		return errRefused
	}
	for _, path := range []string{path, path + ".revoked"} {
		current, e := os.Lstat(path)
		if e != nil || !sameStat(st, current.Sys().(*syscall.Stat_t)) {
			return errRefused
		}
	}
	return nil
}
