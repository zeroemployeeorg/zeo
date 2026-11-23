package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/netip"
	"reflect"
	"regexp"
	"strings"
)

const loginUser = "zeoconnect-broker-run"
const maxPacket = 100000
const maxScript = 65536

var hexPin = regexp.MustCompile(`^[0-9a-f]{64}$`)
var safeName = regexp.MustCompile(`^[a-z][a-z0-9-]{1,63}$`)
var errRefused = errors.New("refused or held; preserve receipts, no blind retry")

type Packet struct {
	Schema        string `json:"schema"`
	Target        string `json:"target_hostname"`
	SourceIP      string `json:"broker_source_ip"`
	SourceNetwork string `json:"broker_source_network"`
	LoginUser     string `json:"login_user"`
	PublicKey     string `json:"public_key"`
	Command       string `json:"command_id"`
	Privilege     string `json:"privilege"`
	Timeout       int    `json:"timeout_seconds"`
	OutputMax     int    `json:"output_max_bytes"`
	ScriptBase64  string `json:"script_base64"`
	ScriptSHA     string `json:"script_sha256"`
}

func hash(raw []byte) string { value := sha256.Sum256(raw); return hex.EncodeToString(value[:]) }

// Closed, duplicate-free JSON; never execute values while validating a packet.
func uniqueValue(d *json.Decoder) error {
	t, err := d.Token()
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); ok {
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				k, e := d.Token()
				if e != nil {
					return e
				}
				key, ok := k.(string)
				if !ok || seen[key] {
					return errRefused
				}
				seen[key] = true
				if e = uniqueValue(d); e != nil {
					return e
				}
			}
		case '[':
			for d.More() {
				if e := uniqueValue(d); e != nil {
					return e
				}
			}
		default:
			return errRefused
		}
		_, err = d.Token()
		return err
	}
	return nil
}

func decodeClosed(raw []byte, out any) error {
	d := json.NewDecoder(bytes.NewReader(raw))
	if err := uniqueValue(d); err != nil {
		return errRefused
	}
	if _, err := d.Token(); err != io.EOF {
		return errRefused
	}
	// encoding/json matches struct fields case-insensitively even with
	// DisallowUnknownFields. Check exact schema tags BEFORE typed decoding,
	// including embedded allocation fields and nested receipt file pins.
	var value any
	if json.Unmarshal(raw, &value) != nil || !canonicalShape(value, reflect.TypeOf(out).Elem()) {
		return errRefused
	}
	d = json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if err := d.Decode(out); err != nil {
		return errRefused
	}
	return nil
}

func schemaFields(t reflect.Type, fields map[string]reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		tag := strings.Split(f.Tag.Get("json"), ",")[0]
		if tag == "-" {
			continue
		}
		if f.Anonymous && tag == "" && f.Type.Kind() == reflect.Struct {
			schemaFields(f.Type, fields)
			continue
		}
		if tag == "" {
			tag = f.Name
		}
		fields[tag] = f.Type
	}
}
func canonicalShape(value any, t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Struct:
		object, ok := value.(map[string]any)
		if !ok {
			return false
		}
		fields := map[string]reflect.Type{}
		schemaFields(t, fields)
		if len(object) != len(fields) {
			return false
		}
		for name, v := range object {
			field, ok := fields[name]
			if !ok || !canonicalShape(v, field) {
				return false
			}
		}
		return true
	case reflect.Map:
		object, ok := value.(map[string]any)
		if !ok || t.Key().Kind() != reflect.String {
			return false
		}
		for _, v := range object {
			if !canonicalShape(v, t.Elem()) {
				return false
			}
		}
		return true
	default:
		return true // The subsequent typed decoder checks scalar types.
	}
}

func parsePacket(raw []byte, pin string) (Packet, []byte, error) {
	var p Packet
	if len(raw) > maxPacket || !hexPin.MatchString(pin) || hash(raw) != pin || decodeClosed(raw, &p) != nil {
		return p, nil, errRefused
	}
	if !validBrokerIP(p.SourceIP) {
		return p, nil, errRefused
	}
	if p.Schema != "zeo.ssh-target-bootstrap/v1" || p.LoginUser != loginUser || !safeName.MatchString(p.Target) || !safeName.MatchString(p.Command) || !hexPin.MatchString(p.ScriptSHA) {
		return p, nil, errRefused
	}
	if (p.SourceNetwork != "tailnet" && p.SourceNetwork != "public") || (p.Privilege != "read" && p.Privilege != "write" && p.Privilege != "root") || p.Timeout < 1 || p.Timeout > 300 || p.OutputMax < 1 || p.OutputMax > 65536 {
		return p, nil, errRefused
	}
	prefix := "ssh-ed25519 "
	if len(p.PublicKey) <= len(prefix) || p.PublicKey[:len(prefix)] != prefix {
		return p, nil, errRefused
	}
	encoded := p.PublicKey[len(prefix):]
	key, e := base64.StdEncoding.Strict().DecodeString(encoded)
	if e != nil || base64.StdEncoding.EncodeToString(key) != encoded || len(key) != 51 || !bytes.Equal(key[:19], []byte("\x00\x00\x00\x0bssh-ed25519\x00\x00\x00\x20")) {
		return p, nil, errRefused
	}
	script, e := base64.StdEncoding.Strict().DecodeString(p.ScriptBase64)
	if e != nil || base64.StdEncoding.EncodeToString(script) != p.ScriptBase64 || len(script) > maxScript || !bytes.HasPrefix(script, []byte("#!/bin/sh\n")) || bytes.IndexByte(script, 0) >= 0 || hash(script) != p.ScriptSHA {
		return p, nil, errRefused
	}
	return p, script, nil
}

// IPv6 zone identifiers are interface metadata, not a broker address. Go accepts
// arbitrary zone text, including authorized_keys delimiters and line breaks.
func validBrokerIP(value string) bool {
	ip, err := netip.ParseAddr(value)
	return err == nil && ip.Zone() == "" && ip.String() == value && !ip.IsUnspecified() && !ip.IsMulticast() && !ip.IsLoopback() && !ip.Is4In6()
}
