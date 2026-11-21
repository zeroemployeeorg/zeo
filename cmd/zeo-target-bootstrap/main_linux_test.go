package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Unprivileged state fixtures: real temporary files and native kernel/process
// primitives, modeled NSS/useradd/visudo. Never modify a real user or SSH host.
type fakeAccounts struct {
	account                 Account
	present, group, expired bool
	creates                 int
	createErr               error
	policyDir               string
}

func (f *fakeAccounts) Lookup() (Account, bool, error) { return f.account, f.present, nil }
func (f *fakeAccounts) GroupExists() (bool, error)     { return f.group, nil }
func (f *fakeAccounts) Create(home string) error {
	f.creates++
	f.present = true
	f.account = Account{12345, 12345, home, "/bin/sh", []string{"12345"}}
	return f.createErr
}
func (f *fakeAccounts) Expire() error { f.expired = true; return nil }
func (f *fakeAccounts) CheckPolicy(path string) error {
	// Simulate a split /var and /etc: a stage on /var cannot be hard-linked
	// to /etc. Validate the DESTINATION-side stage before identity creation.
	if filepath.Dir(path) != f.policyDir {
		return syscall.EXDEV
	}
	_, e := os.ReadFile(path)
	return e
}

func testPacket() (Packet, []byte, []byte) {
	public := append([]byte("\x00\x00\x00\x0bssh-ed25519\x00\x00\x00\x20"), make([]byte, 32)...)
	script := []byte("#!/bin/sh\nprintf 'fixed metadata\\n'\n")
	p := Packet{Schema: "zeo.ssh-target-bootstrap/v1", Target: "test-target", SourceIP: "192.0.2.11", SourceNetwork: "public", LoginUser: loginUser, PublicKey: "ssh-ed25519 " + base64.StdEncoding.EncodeToString(public), Command: "fixed-metadata", Privilege: "read", Timeout: 5, OutputMax: 4096, ScriptBase64: base64.StdEncoding.EncodeToString(script), ScriptSHA: hash(script)}
	return p, encode(p), script
}
func fixture(t *testing.T) (Manager, *fakeAccounts) {
	t.Helper()
	root := t.TempDir()
	p := Paths{filepath.Join(root, "usr/target"), filepath.Join(root, "var/home"), filepath.Join(root, "var/journal"), filepath.Join(root, "etc/policy"), root, uint32(os.Getuid())}
	for _, path := range []string{filepath.Dir(p.Base), filepath.Dir(p.Home), filepath.Dir(p.Journal), filepath.Dir(p.Sudoers)} {
		if e := os.MkdirAll(path, 0755); e != nil {
			t.Fatal(e)
		}
	}
	a := &fakeAccounts{policyDir: filepath.Dir(p.Sudoers)}
	return Manager{Paths: p, Accounts: a, Hostname: func() (string, error) { return "test-target", nil }, Binary: []byte("fixture binary, never executed")}, a
}

func TestClosedPacket(t *testing.T) {
	p, raw, script := testPacket()
	got, b, e := parsePacket(raw, hash(raw))
	if e != nil || got != p || !bytes.Equal(b, script) {
		t.Fatal("valid public packet refused")
	}
	for _, change := range []func(*Packet){
		func(p *Packet) { p.Privilege = "ROOT" }, func(p *Packet) { p.LoginUser = "root" }, func(p *Packet) { p.SourceIP = "192.0.2.0/24" },
		func(p *Packet) { p.Command = "fixed-metadata; id" }, func(p *Packet) { p.PublicKey += " extra" }, func(p *Packet) { p.ScriptSHA = strings.Repeat("0", 64) },
		func(p *Packet) { p.OutputMax = 65537 }, func(p *Packet) { p.Timeout = 301 }, func(p *Packet) { p.SourceNetwork = "unbound" },
	} {
		bad := p
		change(&bad)
		b := encode(bad)
		if _, _, e := parsePacket(b, hash(b)); e == nil {
			t.Fatal("unsafe packet accepted")
		}
	}
}
func TestDuplicateUnknownAndWrongPins(t *testing.T) {
	_, raw, _ := testPacket()
	if _, _, e := parsePacket(raw, strings.Repeat("0", 64)); e == nil {
		t.Fatal("wrong digest accepted")
	}
	for _, extra := range []string{`,"schema":"zeo.ssh-target-bootstrap/v1"}`, `,"unrecognized":true}`} {
		b := append(append([]byte(nil), bytes.TrimSpace(raw)[:len(bytes.TrimSpace(raw))-1]...), []byte(extra)...)
		if _, _, e := parsePacket(b, hash(b)); e == nil {
			t.Fatal("duplicate/unknown field accepted")
		}
	}
}

func TestFoldedPrivilegeRejectedInBothOrders(t *testing.T) {
	_, raw, _ := testPacket()
	body := bytes.TrimSpace(raw)
	for _, b := range [][]byte{
		append(append([]byte(nil), body[:len(body)-1]...), []byte(`,"PRIVILEGE":"root"}`)...),
		append([]byte(`{"PRIVILEGE":"root",`), body[1:]...),
		bytes.ReplaceAll(raw, []byte(`"privilege"`), []byte(`"PRIVILEGE"`)),
		bytes.ReplaceAll(raw, []byte(`"login_user"`), []byte(`"LOGIN_USER"`)),
	} {
		if _, _, e := parsePacket(b, hash(b)); e == nil {
			t.Fatal("case-folded schema key accepted")
		}
	}
}

func TestAllocationAndNestedReceiptRequireCanonicalKeys(t *testing.T) {
	a := Allocation{"zeo.ssh-target-enrollment/v1", strings.Repeat("a", 64), strings.Repeat("b", 64), 12345, 12345, loginUser}
	for _, b := range [][]byte{bytes.ReplaceAll(encode(a), []byte(`"uid"`), []byte(`"UID"`)), bytes.ReplaceAll(encode(a), []byte(`"login_user"`), []byte(`"Login_User"`))} {
		var got Allocation
		if decodeClosed(b, &got) == nil {
			t.Fatal("folded allocation key accepted")
		}
	}
	r := Receipt{a, map[string]FilePin{"/registered": {strings.Repeat("c", 64), 0644}}}
	raw := encode(r)
	var got Receipt
	if decodeClosed(raw, &got) != nil {
		t.Fatal("canonical embedded receipt refused")
	}
	if decodeClosed(bytes.ReplaceAll(raw, []byte(`"mode"`), []byte(`"MODE"`)), &got) == nil {
		t.Fatal("folded nested pin key accepted")
	}
}
func TestReadEnrollmentNoSudoAndIdempotent(t *testing.T) {
	m, a := fixture(t)
	p, raw, script := testPacket()
	state, e := m.apply(p, raw, script)
	if e != nil || state != "ENROLLED_NOT_BROKER_OR_LOGIN_QUALIFIED" {
		t.Fatalf("apply: %v", e)
	}
	if exists(m.Paths.Sudoers) {
		t.Fatal("read class received sudo")
	}
	key, e := os.ReadFile(m.Paths.key())
	if e != nil || !bytes.Contains(key, []byte("restrict,from=\"192.0.2.11\"")) {
		t.Fatal("forced key missing")
	}
	before, _ := os.ReadFile(filepath.Join(m.Paths.Journal, "COMPLETE.json"))
	state, e = m.apply(p, raw, script)
	after, _ := os.ReadFile(filepath.Join(m.Paths.Journal, "COMPLETE.json"))
	if e != nil || state != "ALREADY_INSTALLED" || a.creates != 1 || !bytes.Equal(before, after) {
		t.Fatal("rerun changed completed enrollment")
	}
}
func TestExistingIdentityAndPathRefuseBeforeIntent(t *testing.T) {
	for _, identity := range []bool{false, true} {
		m, a := fixture(t)
		p, raw, script := testPacket()
		if identity {
			a.present = true
		} else {
			if e := os.Mkdir(m.Paths.Base, 0755); e != nil {
				t.Fatal(e)
			}
		}
		if _, e := m.apply(p, raw, script); e == nil || exists(m.Paths.Journal) || a.creates != 0 {
			t.Fatal("existing identity/path adopted")
		}
	}
}
func TestPartialCreateRetainedNoRetry(t *testing.T) {
	m, a := fixture(t)
	p, raw, script := testPacket()
	a.createErr = errors.New("simulated interrupted identity allocation")
	if _, e := m.apply(p, raw, script); e == nil {
		t.Fatal("partial operation claimed success")
	}
	if !exists(filepath.Join(m.Paths.Journal, "INTENT.json")) || !exists(filepath.Join(m.Paths.Journal, "HELD.json")) || exists(m.Paths.key()) {
		t.Fatal("partial evidence or key ordering")
	}
	a.createErr = nil
	if _, e := m.apply(p, raw, script); e == nil || a.creates != 1 {
		t.Fatal("partial operation blindly retried")
	}
}

func TestCancellationAfterIdentityNeverPublishesKey(t *testing.T) {
	m, a := fixture(t)
	p, raw, script := testPacket()
	m.Guard = func() error {
		if a.creates > 0 {
			return context.Canceled
		}
		return nil
	}
	if _, e := m.apply(p, raw, script); e == nil {
		t.Fatal("canceled bootstrap completed")
	}
	if exists(m.Paths.key()) || !exists(filepath.Join(m.Paths.Journal, "HELD.json")) {
		t.Fatal("canceled allocation published authorization or lost evidence")
	}
}
func TestSudoStagingSupportsSplitFilesystems(t *testing.T) {
	m, a := fixture(t)
	p, _, script := testPacket()
	p.Privilege = "root"
	raw := encode(p)
	if _, e := m.apply(p, raw, script); e != nil {
		t.Fatalf("destination-side policy stage failed: %v", e)
	}
	if a.creates != 1 || !exists(m.Paths.Sudoers) {
		t.Fatal("root policy not published")
	}
	policy, _ := os.ReadFile(m.Paths.Sudoers)
	if !bytes.Contains(policy, []byte(m.Paths.tool()+" run\n")) {
		t.Fatal("sudo command not exact")
	}
	m2, a2 := fixture(t)
	a2.policyDir = filepath.Join(m2.Paths.Journal, "wrong-filesystem")
	if _, e := m2.apply(p, raw, script); !errors.Is(e, syscall.EXDEV) || a2.creates != 0 || exists(m2.Paths.key()) {
		t.Fatal("policy refusal created identity/key")
	}
}
func TestChangedPostimageAndForgedReceiptCannotRevoke(t *testing.T) {
	for _, forge := range []bool{false, true} {
		m, a := fixture(t)
		p, raw, script := testPacket()
		if _, e := m.apply(p, raw, script); e != nil {
			t.Fatal(e)
		}
		if forge {
			r, e := m.receipt()
			if e != nil {
				t.Fatal(e)
			}
			pin := r.Files[m.Paths.key()]
			pin.SHA = strings.Repeat("0", 64)
			r.Files[m.Paths.key()] = pin
			if e = os.WriteFile(filepath.Join(m.Paths.Journal, "COMPLETE.json"), encode(r), 0600); e != nil {
				t.Fatal(e)
			}
		} else {
			if e := os.WriteFile(m.Paths.key(), []byte("other owner's bytes"), 0644); e != nil {
				t.Fatal(e)
			}
		}
		if _, e := m.revoke(p, raw, script); e == nil || a.expired || exists(filepath.Join(m.Paths.Journal, "REVOCATION-INTENT.json")) || !exists(m.Paths.key()) {
			t.Fatal("changed postimage revoked")
		}
	}
}
func TestOwnedRevokePreservesGoldAndNamesDrainLimit(t *testing.T) {
	m, a := fixture(t)
	p, raw, script := testPacket()
	if _, e := m.apply(p, raw, script); e != nil {
		t.Fatal(e)
	}
	before, _ := os.ReadFile(m.Paths.key())
	state, e := m.revoke(p, raw, script)
	if e != nil || !strings.Contains(state, "NOT_DRAINED") || !a.expired || exists(m.Paths.key()) {
		t.Fatalf("revoke: %v", e)
	}
	after, _ := os.ReadFile(m.Paths.key() + ".revoked")
	if !bytes.Equal(before, after) {
		t.Fatal("revoked authorization not preserved")
	}
	state, e = m.revoke(p, raw, script)
	if e != nil || state != "ALREADY_REVOKED_NEW_CONNECTIONS_ONLY" {
		t.Fatal("revoke not idempotent")
	}
	if _, e = m.apply(p, raw, script); e == nil {
		t.Fatal("revoked credential reactivated")
	}
}
func TestNativeSealedSnapshotCannotBeChanged(t *testing.T) {
	raw := []byte("verified source bytes")
	f, e := sealedSnapshot(raw)
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	raw[0] = 'X'
	got, e := io.ReadAll(f)
	if e != nil || string(got) != "verified source bytes" {
		t.Fatal("snapshot followed mutable source")
	}
	if _, e = f.WriteAt([]byte("tamper"), 0); e == nil {
		t.Fatal("sealed bytes writable")
	}
}
func TestNativeOutputBound(t *testing.T) {
	var out, stderr bytes.Buffer
	e := runCaptured(context.Background(), []byte("#!/bin/sh\nwhile :; do printf 'xxxxxxxxxxxxxxxx'; done\n"), time.Second, 128, &out, &stderr)
	if e == nil || out.Len()+stderr.Len() > 128 {
		t.Fatal("output limit did not stop registered command")
	}
}
func TestNativeDeadlineCleansOwnedProcessGroup(t *testing.T) {
	var out, stderr bytes.Buffer
	start := time.Now()
	e := runCaptured(context.Background(), []byte("#!/bin/sh\ntrap '' TERM\nsleep 10 &\nprintf '%s\\n' \"$!\"\nwait\n"), 150*time.Millisecond, 4096, &out, &stderr)
	if e == nil || time.Since(start) > 4*time.Second {
		t.Fatal("owned process-group deadline failed")
	}
	pid, e := strconv.Atoi(strings.TrimSpace(out.String()))
	if e != nil {
		t.Fatal("child identity not observed")
	}
	// A reparented zombie is no longer executing; never kill an arbitrary PID.
	until := time.Now().Add(time.Second)
	for {
		b, e := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
		if os.IsNotExist(e) {
			break
		}
		if e != nil {
			t.Fatal("owned child drain not observed")
		}
		parts := strings.Fields(string(b))
		if len(parts) >= 3 && (parts[2] == "Z" || parts[2] == "X") {
			break
		}
		if time.Now().After(until) {
			t.Fatal("owned child still running after deadline")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
func TestNativeCommandLockIsExclusive(t *testing.T) {
	m, _ := fixture(t)
	if e := os.Mkdir(m.Paths.Base, 0755); e != nil {
		t.Fatal(e)
	}
	if e := createFile(filepath.Join(m.Paths.Base, "command.lock"), nil, 0644); e != nil {
		t.Fatal(e)
	}
	first, e := lockCommand(m.Paths)
	if e != nil {
		t.Fatal(e)
	}
	if second, e := lockCommand(m.Paths); e == nil {
		second.Close()
		t.Fatal("concurrent session admitted")
	}
	first.Close()
	third, e := lockCommand(m.Paths)
	if e != nil {
		t.Fatal(e)
	}
	third.Close()
}
func TestPacketJSONRoundTripIsNotTrust(t *testing.T) {
	p, raw, _ := testPacket()
	var other Packet
	if json.Unmarshal(raw, &other) != nil || other != p {
		t.Fatal("fixture mismatch")
	}
	if _, _, e := parsePacket(append(raw, ' '), hash(raw)); e == nil {
		t.Fatal("raw packet pin ignored")
	}
}

func TestBrokerIPv6ZoneCannotInjectAuthorizedKeys(t *testing.T) {
	m, _ := fixture(t)
	original, _, _ := testPacket()
	for _, ip := range []string{"fe80::1%eth0", "2001:db8::1%tailnet", "fe80::1%\"\nssh-ed25519 AAAA unrestricted\n#"} {
		p := original
		p.SourceIP = ip
		raw := encode(p)
		if _, _, err := parsePacket(raw, hash(raw)); err == nil {
			t.Fatal("IPv6 zone accepted", ip)
		}
		if line := m.Paths.keyLine(p); len(line) != 0 {
			t.Fatal("unsafe address serialized", string(line))
		}
	}
	for _, ip := range []string{"192.0.2.11", "2001:db8::1"} {
		p := original
		p.SourceIP = ip
		raw := encode(p)
		if _, _, err := parsePacket(raw, hash(raw)); err != nil {
			t.Fatal("canonical address refused", ip)
		}
		line := string(m.Paths.keyLine(p))
		want := "restrict,from=\"" + ip + "\",command=\"" + m.Paths.tool() + " dispatch\" " + p.PublicKey + "\n"
		if line != want || strings.Count(line, "\n") != 1 || strings.Count(line, "\"") != 4 {
			t.Fatal("key-line template changed", line)
		}
	}
}

func TestRevokeClosesLocalExecutionAndRetainsExactSudoPolicy(t *testing.T) {
	for _, privilege := range []string{"read", "root"} {
		t.Run(privilege, func(t *testing.T) {
			m, _ := fixture(t)
			p, _, script := testPacket()
			p.Privilege = privilege
			raw := encode(p)
			if _, err := m.apply(p, raw, script); err != nil {
				t.Fatal(err)
			}
			if err := activeAuthorization(m.Paths, p, true); err != nil {
				t.Fatal("active enrollment refused", err)
			}
			var policy []byte
			if privilege == "root" {
				var err error
				policy, err = os.ReadFile(m.Paths.Sudoers)
				if err != nil {
					t.Fatal(err)
				}
			}
			if _, err := m.revoke(p, raw, script); err != nil {
				t.Fatal(err)
			}
			if err := activeAuthorization(m.Paths, p, true); err == nil {
				t.Fatal("revoked enrollment can start local execution")
			}
			if privilege == "root" {
				if exists(m.Paths.Sudoers) {
					t.Fatal("NOPASSWD rule still active")
				}
				retained, err := os.ReadFile(m.Paths.Sudoers + ".revoked")
				if err != nil || !bytes.Equal(retained, policy) {
					t.Fatal("sudo policy not preserved exactly", err)
				}
			}
			if _, err := m.revoke(p, raw, script); err != nil {
				t.Fatal("completed revoke not idempotent", err)
			}
		})
	}
}

func TestPartialRevocationClosesExecutionBeforeAccountExpiry(t *testing.T) {
	m, _ := fixture(t)
	p, _, script := testPacket()
	p.Privilege = "root"
	raw := encode(p)
	if _, err := m.apply(p, raw, script); err != nil {
		t.Fatal(err)
	}
	// Model interruption immediately after the authorization removal. Even with
	// the original NOPASSWD policy and account still present, neither run mode
	// may pass the common authorization gate.
	if err := os.Remove(m.Paths.key()); err != nil {
		t.Fatal(err)
	}
	if !exists(m.Paths.Sudoers) {
		t.Fatal("fixture lacks surviving elevation rule")
	}
	if err := activeAuthorization(m.Paths, p, true); err == nil {
		t.Fatal("partial revocation allowed new execution")
	}
}
