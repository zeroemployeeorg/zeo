package main

import (
	"bytes"
	"context"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"
)

type Account struct {
	UID, GID    uint32
	Home, Shell string
	Groups      []string
}
type Accounts interface {
	Lookup() (Account, bool, error)
	GroupExists() (bool, error)
	Create(string) error
	Expire() error
	CheckPolicy(string) error
}
type systemAccounts struct{}

func (systemAccounts) Lookup() (Account, bool, error) {
	u, e := user.Lookup(loginUser)
	if _, ok := e.(user.UnknownUserError); ok {
		return Account{}, false, nil
	}
	if e != nil {
		return Account{}, false, e
	}
	uid, e := strconv.ParseUint(u.Uid, 10, 32)
	if e != nil {
		return Account{}, false, e
	}
	gid, e := strconv.ParseUint(u.Gid, 10, 32)
	if e != nil {
		return Account{}, false, e
	}
	groups, e := u.GroupIds()
	if e != nil {
		return Account{}, false, e
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, "/usr/bin/getent", "passwd", loginUser)
	c.Env = []string{"PATH=/usr/bin:/bin", "LC_ALL=C"}
	out := &toolBuffer{}
	c.Stdout = out
	if e = c.Run(); e != nil {
		return Account{}, false, errRefused
	}
	fields := strings.Split(strings.TrimSuffix(out.String(), "\n"), ":")
	if len(fields) != 7 || fields[0] != loginUser || fields[2] != u.Uid || fields[3] != u.Gid || strings.Contains(fields[6], "\n") {
		return Account{}, false, errRefused
	}
	return Account{uint32(uid), uint32(gid), u.HomeDir, fields[6], groups}, true, nil
}

type toolBuffer struct{ data bytes.Buffer }

func (b *toolBuffer) Write(p []byte) (int, error) {
	if b.data.Len()+len(p) > 8192 {
		return 0, errRefused
	}
	return b.data.Write(p)
}
func (b *toolBuffer) String() string { return b.data.String() }
func (systemAccounts) GroupExists() (bool, error) {
	_, e := user.LookupGroup(loginUser)
	if _, ok := e.(user.UnknownGroupError); ok {
		return false, nil
	}
	return e == nil, e
}
func tool(argv ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, argv[0], argv[1:]...)
	c.Env = []string{"PATH=/usr/sbin:/usr/bin:/sbin:/bin", "LC_ALL=C"}
	if e := c.Run(); e != nil {
		return errRefused
	}
	return nil
}
func (systemAccounts) Create(home string) error {
	return tool("/usr/sbin/useradd", "--system", "--user-group", "--no-create-home", "--home-dir", home, "--shell", "/bin/sh", "--password", "!", loginUser)
}
func (systemAccounts) Expire() error {
	return tool("/usr/sbin/usermod", "-L", "--expiredate", "1", loginUser)
}
func (systemAccounts) CheckPolicy(path string) error { return tool("/usr/sbin/visudo", "-cf", path) }
