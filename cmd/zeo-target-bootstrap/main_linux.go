package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	defer cancel()
	code := runCLI(ctx, os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(code)
}

func runCLI(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	p := nativePaths()
	if len(args) == 1 && (args[0] == "dispatch" || args[0] == "run") {
		exe, e := os.Executable()
		if e != nil || filepath.Clean(exe) != p.tool() {
			return 1
		}
		if dispatch(ctx, p, args[0] == "run", stdout, stderr) != nil {
			return 1
		}
		return 0
	}
	flags := flag.NewFlagSet("zeo-target-bootstrap", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	packetPath := flags.String("packet", "", "reviewed PUBLIC packet")
	packetPin := flags.String("packet-sha256", "", "packet digest, not release trust")
	binaryPin := flags.String("binary-sha256", "", "qualified executable digest")
	apply := flags.Bool("apply", false, "create initial target enrollment")
	revoke := flags.Bool("revoke", false, "revoke only the owned enrollment")
	if flags.Parse(args) != nil || flags.NArg() != 0 || *packetPath == "" || (*apply && *revoke) {
		io.WriteString(stderr, "zeo-target-bootstrap --packet FILE --packet-sha256 HEX --binary-sha256 HEX [--apply|--revoke]\n")
		return 2
	}
	raw, e := readStable(*packetPath, maxPacket, uint32(os.Getuid()), 0)
	if e != nil {
		return fail(stderr)
	}
	packet, script, e := parsePacket(raw, *packetPin)
	if e != nil {
		return fail(stderr)
	}
	if !*apply && !*revoke {
		_ = json.NewEncoder(stdout).Encode(map[string]string{"state": "PLAN_ONLY", "target": packet.Target, "login_user": packet.LoginUser, "privilege": packet.Privilege, "packet_sha256": *packetPin, "admission": "NOT_INFERRED"})
		return 0
	}
	if os.Getuid() != 0 || os.Geteuid() != 0 || os.Getgid() != 0 || os.Getegid() != 0 || !hexPin.MatchString(*binaryPin) {
		return fail(stderr)
	}
	syscall.Umask(0077)
	exe, e := os.Executable()
	if e != nil {
		return fail(stderr)
	}
	binary, e := readStable(exe, 64<<20, 0, 0)
	if e != nil || hash(binary) != *binaryPin {
		return fail(stderr)
	}
	if *apply {
		// OS tools are preconditions, never downloaded or substituted here.
		tools := []string{"/usr/sbin/useradd", "/usr/sbin/usermod", "/usr/bin/getent", "/bin/sh"}
		if packet.Privilege == "root" {
			tools = append(tools, "/usr/bin/sudo", "/usr/sbin/visudo")
		}
		for _, name := range tools {
			resolved, e := filepath.EvalSymlinks(name)
			if e != nil || p.trustedDir(filepath.Dir(resolved)) != nil {
				return fail(stderr)
			}
			if _, e = readStable(resolved, 64<<20, 0, 0); e != nil {
				return fail(stderr)
			}
		}
		f, e := sealedSnapshot(script)
		if e != nil {
			return fail(stderr)
		}
		f.Close()
	}
	m := Manager{Paths: p, Accounts: systemAccounts{}, Hostname: os.Hostname, Binary: binary, Guard: ctx.Err}
	var state string
	if *apply {
		state, e = m.apply(packet, raw, script)
	} else {
		state, e = m.revoke(packet, raw, script)
	}
	if e != nil {
		return fail(stderr)
	}
	_ = json.NewEncoder(stdout).Encode(map[string]string{"state": state, "packet_sha256": *packetPin, "binary_sha256": *binaryPin, "login_user": packet.LoginUser, "privilege": packet.Privilege, "command_id": packet.Command})
	return 0
}
func fail(stderr io.Writer) int {
	io.WriteString(stderr, "{\"state\":\"REFUSED_OR_HELD\",\"action\":\"PRESERVE_RECEIPTS_NO_BLIND_RETRY\"}\n")
	return 1
}
