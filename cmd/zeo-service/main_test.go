package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderCommandHasNoOutputBeforeValidation(t *testing.T) {
	var out, stderr bytes.Buffer
	code := run([]string{"--organization", "example", "--seat", "pilot", "--user", "root"}, &out, &stderr)
	if code == 0 || out.Len() != 0 || stderr.Len() == 0 {
		t.Fatal("invalid configuration must fail without partial unit output")
	}
}

func TestRenderCommandUsesExplicitProfileInputs(t *testing.T) {
	var out, stderr bytes.Buffer
	code := run([]string{"--organization", "example", "--seat", "pilot", "--user", "pilot_user", "--group", "pilot_group", "--pomar-sha256", strings.Repeat("b", 64), "--memory-bytes", "2147483648", "--tasks", "256", "--cpu-percent", "200"}, &out, &stderr)
	if code != 0 || stderr.Len() != 0 || !strings.Contains(out.String(), "CPUQuota=200%\n") {
		t.Fatalf("unit command failed: code=%d stderr=%s", code, stderr.String())
	}
}
