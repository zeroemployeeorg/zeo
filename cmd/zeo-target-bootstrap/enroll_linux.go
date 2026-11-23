package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type FilePin struct {
	SHA  string `json:"sha256"`
	Mode uint32 `json:"mode"`
}
type Allocation struct {
	Schema    string `json:"schema"`
	PacketSHA string `json:"packet_sha256"`
	BinarySHA string `json:"binary_sha256"`
	UID       uint32 `json:"uid"`
	GID       uint32 `json:"gid"`
	LoginUser string `json:"login_user"`
}
type Receipt struct {
	Allocation
	Files map[string]FilePin `json:"files"`
}
type Content struct {
	Raw  []byte
	Mode os.FileMode
}
type Manager struct {
	Paths    Paths
	Accounts Accounts
	Hostname func() (string, error)
	Binary   []byte
	Guard    func() error
}

func (m Manager) check() error {
	if m.Guard != nil && m.Guard() != nil {
		return errRefused
	}
	return nil
}

func encode(value any) []byte { raw, _ := json.Marshal(value); return append(raw, '\n') }
func (p Paths) policy() []byte {
	return []byte(loginUser + " ALL=(root) NOPASSWD: " + p.tool() + " run\n")
}
func (p Paths) keyLine(packet Packet) []byte {
	// Defense at the serialization boundary as well as packet parsing. Only the
	// template may introduce quotes or the final newline.
	if !validBrokerIP(packet.SourceIP) || strings.ContainsAny(packet.PublicKey, "\r\n\"") || strings.ContainsAny(p.tool(), "\r\n\"") {
		return nil
	}
	return []byte("restrict,from=\"" + packet.SourceIP + "\",command=\"" + p.tool() + " dispatch\" " + packet.PublicKey + "\n")
}
func (m Manager) boundHost(p Packet) error {
	host, e := m.Hostname()
	if e != nil || strings.Split(host, ".")[0] != p.Target {
		return errRefused
	}
	return nil
}
func validAccount(a Account, p Paths) bool {
	return a.UID != 0 && a.GID != 0 && a.Home == p.Home && a.Shell == "/bin/sh" && len(a.Groups) == 1 && a.Groups[0] == strconv.FormatUint(uint64(a.GID), 10)
}

func (m Manager) contents(packet Packet, raw, script []byte, a Allocation) map[string]Content {
	p := m.Paths
	files := map[string]Content{
		p.tool():                                 {m.Binary, 0755},
		filepath.Join(p.Base, "command.sh"):      {script, 0755},
		filepath.Join(p.Base, "packet.json"):     {raw, 0644},
		filepath.Join(p.Base, "allocation.json"): {encode(a), 0644},
		filepath.Join(p.Base, "command.lock"):    {nil, 0644},
		p.key():                                  {p.keyLine(packet), 0644},
	}
	if packet.Privilege == "root" {
		files[p.Sudoers] = Content{p.policy(), 0440}
	}
	return files
}
func pins(files map[string]Content) map[string]FilePin {
	out := map[string]FilePin{}
	for path, c := range files {
		out[path] = FilePin{hash(c.Raw), uint32(c.Mode)}
	}
	return out
}

func (m Manager) verify(r Receipt, packet Packet, raw, script []byte, revoked bool) error {
	p := m.Paths
	if r.Schema != "zeo.ssh-target-enrollment/v1" || r.LoginUser != loginUser || r.PacketSHA != hash(raw) || r.BinarySHA != hash(m.Binary) {
		return errRefused
	}
	a, found, e := m.Accounts.Lookup()
	if e != nil || !found || !validAccount(a, p) || a.UID != r.UID || a.GID != r.GID {
		return errRefused
	}
	expected := pins(m.contents(packet, raw, script, r.Allocation))
	if !reflect.DeepEqual(r.Files, expected) {
		return errRefused
	}
	for _, dir := range []string{p.Base, p.Home, filepath.Join(p.Home, ".ssh"), p.Journal} {
		if p.trustedDir(dir) != nil {
			return errRefused
		}
	}
	for name, pin := range expected {
		if (name == p.key() || name == p.Sudoers) && revoked {
			if exists(name) {
				return errRefused
			}
			name += ".revoked"
		}
		b, e := readStable(name, 64<<20, p.Owner, os.FileMode(pin.Mode))
		if e != nil || hash(b) != pin.SHA {
			return errRefused
		}
	}
	return nil
}
func (m Manager) receipt() (Receipt, error) {
	var r Receipt
	raw, e := readStable(filepath.Join(m.Paths.Journal, "COMPLETE.json"), maxPacket, m.Paths.Owner, 0600)
	if e != nil || decodeClosed(raw, &r) != nil {
		return r, errRefused
	}
	return r, nil
}

func (m Manager) apply(packet Packet, raw, script []byte) (result string, err error) {
	p := m.Paths
	if m.check() != nil || m.boundHost(packet) != nil || len(p.keyLine(packet)) == 0 {
		return "", errRefused
	}
	if exists(p.Journal) {
		if p.trustedDir(p.Journal) != nil || exists(filepath.Join(p.Journal, "REVOKED.json")) {
			return "", errRefused
		}
		r, e := m.receipt()
		if e != nil || m.verify(r, packet, raw, script, false) != nil {
			return "", errRefused
		}
		return "ALREADY_INSTALLED", nil
	}
	for _, dir := range []string{filepath.Dir(p.Base), filepath.Dir(p.Home), filepath.Dir(p.Journal), filepath.Dir(p.Sudoers)} {
		if p.trustedDir(dir) != nil {
			return "", errRefused
		}
	}
	for _, path := range []string{p.Base, p.Home, p.Sudoers} {
		if exists(path) {
			return "", errRefused
		}
	}
	// A sudoers hard-link must remain on /etc's own filesystem. /var may be
	// a separate mount. The dot-named stage is ignored by sudo's includedir
	// and is checked explicitly before any identity is created.
	staged := filepath.Join(filepath.Dir(p.Sudoers), ".zeo-ssh-target-"+hash(raw)+".pending")
	if packet.Privilege == "root" && exists(staged) {
		return "", errRefused
	}
	_, found, e := m.Accounts.Lookup()
	if e != nil || found {
		return "", errRefused
	}
	found, e = m.Accounts.GroupExists()
	if e != nil || found {
		return "", errRefused
	}
	if len(m.Binary) == 0 || len(m.Binary) > 64<<20 {
		return "", errRefused
	}
	if e = mkdirNew(p.Journal, 0700); e != nil {
		return "", e
	}
	if e = writeRecord(filepath.Join(p.Journal, "INTENT.json"), map[string]string{"schema": "zeo.ssh-target-bootstrap-intent/v1", "packet_sha256": hash(raw), "binary_sha256": hash(m.Binary), "state": "PREPARE"}); e != nil {
		return "", e
	}
	defer func() {
		if err != nil {
			_ = writeRecord(filepath.Join(p.Journal, "HELD.json"), map[string]string{"state": "PARTIAL_RECONCILE_NO_RETRY", "packet_sha256": hash(raw)})
		}
	}()
	if m.check() != nil {
		return "", errRefused
	}
	if packet.Privilege == "root" {
		if e = createFile(staged, p.policy(), 0440); e != nil {
			return "", e
		}
		if e = m.Accounts.CheckPolicy(staged); e != nil {
			return "", e
		}
	}
	if m.check() != nil {
		return "", errRefused
	}
	if e = m.Accounts.Create(p.Home); e != nil {
		return "", e
	}
	if m.check() != nil {
		return "", errRefused
	}
	a, found, e := m.Accounts.Lookup()
	if e != nil || !found || !validAccount(a, p) {
		return "", errRefused
	}
	for _, dir := range []string{p.Base, p.Home, filepath.Join(p.Home, ".ssh")} {
		if m.check() != nil {
			return "", errRefused
		}
		if e = mkdirNew(dir, 0755); e != nil {
			return "", e
		}
	}
	allocation := Allocation{"zeo.ssh-target-enrollment/v1", hash(raw), hash(m.Binary), a.UID, a.GID, loginUser}
	files := m.contents(packet, raw, script, allocation)
	for _, name := range []string{p.tool(), filepath.Join(p.Base, "command.sh"), filepath.Join(p.Base, "packet.json"), filepath.Join(p.Base, "allocation.json"), filepath.Join(p.Base, "command.lock")} {
		if m.check() != nil {
			return "", errRefused
		}
		c := files[name]
		if e = createFile(name, c.Raw, c.Mode); e != nil {
			return "", e
		}
	}
	if packet.Privilege == "root" {
		if m.check() != nil {
			return "", errRefused
		}
		if e = os.Link(staged, p.Sudoers); e != nil {
			return "", e
		}
		if e = syncDir(filepath.Dir(p.Sudoers)); e != nil {
			return "", e
		}
		if e = os.Remove(staged); e != nil {
			return "", e
		}
		if e = syncDir(filepath.Dir(staged)); e != nil {
			return "", e
		}
	}
	// Public authorization is published LAST. No actor receives a private key.
	if m.check() != nil {
		return "", errRefused
	}
	key := files[p.key()]
	if e = createFile(p.key(), key.Raw, key.Mode); e != nil {
		return "", e
	}
	r := Receipt{allocation, pins(files)}
	if e = m.verify(r, packet, raw, script, false); e != nil {
		return "", e
	}
	if m.check() != nil {
		return "", errRefused
	}
	if e = writeRecord(filepath.Join(p.Journal, "COMPLETE.json"), r); e != nil {
		return "", e
	}
	return "ENROLLED_NOT_BROKER_OR_LOGIN_QUALIFIED", nil
}

func (m Manager) revoke(packet Packet, raw, script []byte) (string, error) {
	p := m.Paths
	if m.check() != nil || m.boundHost(packet) != nil || p.trustedDir(p.Journal) != nil {
		return "", errRefused
	}
	r, e := m.receipt()
	if e != nil {
		return "", e
	}
	if exists(filepath.Join(p.Journal, "REVOKED.json")) {
		if m.verify(r, packet, raw, script, true) != nil {
			return "", errRefused
		}
		return "ALREADY_REVOKED_NEW_CONNECTIONS_ONLY", nil
	}
	if exists(filepath.Join(p.Journal, "REVOCATION-INTENT.json")) || exists(p.key()+".revoked") || (packet.Privilege == "root" && exists(p.Sudoers+".revoked")) || m.verify(r, packet, raw, script, false) != nil {
		return "", errRefused
	}
	if e = writeRecord(filepath.Join(p.Journal, "REVOCATION-INTENT.json"), map[string]string{"state": "REVOKING_EXACT_OWNED_KEY", "packet_sha256": hash(raw)}); e != nil {
		return "", e
	}
	// Hard-link first, verify the retained bytes, then remove only the owned
	// authorization name. A partial revoke is HELD, never blindly repeated.
	if e = os.Link(p.key(), p.key()+".revoked"); e != nil {
		return "", e
	}
	if e = verifyLinkedAuthorization(p, r.Files[p.key()].SHA); e != nil {
		return "", errRefused
	}
	if e = os.Remove(p.key()); e != nil {
		return "", e
	}
	if e = syncDir(filepath.Dir(p.key())); e != nil {
		return "", e
	}
	// Remove the exact owned elevation rule too. The absent authorization already
	// closes both dispatcher modes, including an already authenticated connection.
	// Preserve the original policy without overwriting any prior retained object.
	if packet.Privilege == "root" {
		if e = os.Link(p.Sudoers, p.Sudoers+".revoked"); e != nil {
			return "", e
		}
		if e = verifyLinkedFile(p.Sudoers, r.Files[p.Sudoers].SHA, 0440, p.Owner); e != nil {
			return "", errRefused
		}
		if e = os.Remove(p.Sudoers); e != nil {
			return "", e
		}
		if e = syncDir(filepath.Dir(p.Sudoers)); e != nil {
			return "", e
		}
	}
	if e = m.Accounts.Expire(); e != nil {
		return "", e
	}
	if e = m.verify(r, packet, raw, script, true); e != nil {
		return "", e
	}
	if e = writeRecord(filepath.Join(p.Journal, "REVOKED.json"), map[string]string{"state": "NEW_CONNECTIONS_REVOKED_ACCOUNT_EXPIRE_REQUESTED_NOT_DRAINED", "packet_sha256": hash(raw)}); e != nil {
		return "", e
	}
	return "NEW_CONNECTIONS_REVOKED_ACCOUNT_EXPIRE_REQUESTED_NOT_DRAINED", nil
}
