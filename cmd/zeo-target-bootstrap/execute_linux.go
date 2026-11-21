package main

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

func sealedSnapshot(raw []byte) (*os.File, error) {
	if memfdSyscall == 0 {
		return nil, errRefused
	}
	name := append([]byte("zeo-registered-command"), 0)
	// CLOEXEC | ALLOW_SEALING; ExtraFiles explicitly maps only child fd 3.
	fd, _, errno := syscall.Syscall(uintptr(memfdSyscall), uintptr(unsafe.Pointer(&name[0])), 3, 0)
	if errno != 0 {
		return nil, errno
	}
	f := os.NewFile(fd, "sealed-command")
	ok := false
	defer func() {
		if !ok {
			f.Close()
		}
	}()
	if _, e := f.Write(raw); e != nil {
		return nil, e
	}
	// F_ADD_SEALS: seal, shrink, grow, write. The child receives these exact
	// captured bytes even if an administrator replaces the registered path.
	if _, _, errno = syscall.Syscall(syscall.SYS_FCNTL, fd, 1033, 15); errno != 0 {
		return nil, errno
	}
	if _, e := f.Seek(0, 0); e != nil {
		return nil, e
	}
	ok = true
	return f, nil
}

type cappedOutput struct {
	mu             sync.Mutex
	left           int
	exceeded       bool
	cancel         context.CancelFunc
	stdout, stderr io.Writer
}
type outputSide struct {
	parent *cappedOutput
	stderr bool
}

func (w outputSide) Write(raw []byte) (int, error) {
	b := w.parent
	b.mu.Lock()
	defer b.mu.Unlock()
	out := b.stdout
	if w.stderr {
		out = b.stderr
	}
	count := len(raw)
	if count > b.left {
		count = b.left
		b.exceeded = true
		b.cancel()
	}
	n, e := out.Write(raw[:count])
	b.left -= n
	if e != nil {
		return n, e
	}
	if n != len(raw) {
		return n, errRefused
	}
	return n, nil
}
func (b *cappedOutput) overflow() bool { b.mu.Lock(); defer b.mu.Unlock(); return b.exceeded }

func runCaptured(parent context.Context, raw []byte, timeout time.Duration, limit int, stdout, stderr io.Writer) error {
	snapshot, e := sealedSnapshot(raw)
	if e != nil {
		return e
	}
	defer snapshot.Close()
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	output := &cappedOutput{left: limit, cancel: cancel, stdout: stdout, stderr: stderr}
	c := exec.Command("/bin/sh", "/proc/self/fd/3")
	c.ExtraFiles = []*os.File{snapshot}
	c.Dir = "/"
	c.Env = []string{"PATH=/usr/bin:/bin", "LC_ALL=C"}
	c.Stdin = nil
	c.Stdout = outputSide{output, false}
	c.Stderr = outputSide{output, true}
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.WaitDelay = 2 * time.Second
	if e = c.Start(); e != nil {
		return e
	}
	done := make(chan error, 1)
	go func() { done <- c.Wait() }()
	var result error
	select {
	case result = <-done:
		// A successful leader must not leave background descendants alive.
		if syscall.Kill(-c.Process.Pid, 0) == nil {
			_ = syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
			result = errRefused
		}
	case <-ctx.Done():
		_ = syscall.Kill(-c.Process.Pid, syscall.SIGTERM)
		waited := false
		select {
		case <-done:
			waited = true
		case <-time.After(250 * time.Millisecond):
		}
		_ = syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
		// WaitDelay bounds inherited-pipe closure. Never start a replacement.
		if !waited {
			select {
			case <-done:
			case <-time.After(3 * time.Second):
			}
		}
		result = errRefused
	}
	if output.overflow() || ctx.Err() != nil {
		return errRefused
	}
	return result
}

func loadInstalled(p Paths) (Packet, Allocation, []byte, error) {
	var a Allocation
	var packet Packet
	if p.trustedDir(p.Base) != nil {
		return packet, a, nil, errRefused
	}
	allocation, e := readStable(filepath.Join(p.Base, "allocation.json"), maxPacket, p.Owner, 0644)
	if e != nil || decodeClosed(allocation, &a) != nil || a.Schema != "zeo.ssh-target-enrollment/v1" || a.LoginUser != loginUser || a.UID == 0 || a.GID == 0 {
		return packet, a, nil, errRefused
	}
	raw, e := readStable(filepath.Join(p.Base, "packet.json"), maxPacket, p.Owner, 0644)
	if e != nil {
		return packet, a, nil, e
	}
	packet, script, e := parsePacket(raw, a.PacketSHA)
	if e != nil {
		return packet, a, nil, e
	}
	binary, e := readStable(p.tool(), 64<<20, p.Owner, 0755)
	if e != nil || hash(binary) != a.BinarySHA {
		return packet, a, nil, errRefused
	}
	registered, e := readStable(filepath.Join(p.Base, "command.sh"), maxScript, p.Owner, 0755)
	if e != nil || hash(registered) != packet.ScriptSHA || hash(registered) != hash(script) {
		return packet, a, nil, errRefused
	}
	if (Manager{Hostname: os.Hostname}).boundHost(packet) != nil {
		return packet, a, nil, errRefused
	}
	return packet, a, registered, nil
}

// The active public authorization is also a local execution gate. It is readable
// by the restricted account; the private enrollment journal is never exposed.
// Revocation removes this name before account expiry or sudo-policy cleanup.
func activeAuthorization(p Paths, packet Packet, checkPolicy bool) error {
	if p.trustedDir(filepath.Join(p.Home, ".ssh")) != nil {
		return errRefused
	}
	expected := p.keyLine(packet)
	key, err := readStable(p.key(), maxPacket, p.Owner, 0644)
	if err != nil || len(expected) == 0 || hash(key) != hash(expected) {
		return errRefused
	}
	if packet.Privilege == "root" && checkPolicy {
		if p.trustedDir(filepath.Dir(p.Sudoers)) != nil {
			return errRefused
		}
		policy, err := readStable(p.Sudoers, maxPacket, p.Owner, 0440)
		if err != nil || hash(policy) != hash(p.policy()) {
			return errRefused
		}
	}
	return nil
}

func dispatch(parent context.Context, p Paths, rootMode bool, stdout, stderr io.Writer) error {
	packet, a, script, e := loadInstalled(p)
	if e != nil {
		return e
	}
	if activeAuthorization(p, packet, rootMode) != nil {
		return errRefused
	}
	if rootMode {
		if packet.Privilege != "root" || os.Getuid() != 0 || os.Geteuid() != 0 || os.Getgid() != 0 || os.Getegid() != 0 {
			return errRefused
		}
	} else {
		if uint32(os.Getuid()) != a.UID || uint32(os.Geteuid()) != a.UID || uint32(os.Getgid()) != a.GID || uint32(os.Getegid()) != a.GID || os.Getenv("SSH_ORIGINAL_COMMAND") != packet.Command {
			return errRefused
		}
		groups, e := os.Getgroups()
		if e != nil {
			return errRefused
		}
		for _, group := range groups {
			if uint32(group) != a.GID {
				return errRefused
			}
		}
		if packet.Privilege == "root" {
			return syscall.Exec("/usr/bin/sudo", []string{"sudo", "-n", "--", p.tool(), "run"}, []string{"PATH=/usr/bin:/bin", "LC_ALL=C"})
		}
		// read/write classes remain UNPRIVILEGED. The broker owns effect approval.
	}
	lock, e := lockCommand(p)
	if e != nil {
		return e
	}
	defer lock.Close()
	// Recheck after taking the execution lock. Already running commands are not
	// drained by revocation; no new invocation may start after its gate closes.
	if activeAuthorization(p, packet, rootMode) != nil {
		return errRefused
	}
	return runCaptured(parent, script, time.Duration(packet.Timeout)*time.Second, packet.OutputMax, stdout, stderr)
}
