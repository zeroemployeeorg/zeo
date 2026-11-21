package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Reproduce the two checks lost during the historical v8 manual repin. These
// names are stable obligations, independent of the source commit or method name.
func fixture(t *testing.T) (string, []string, catalog) {
	t.Helper()
	dir := t.TempDir()
	body := []byte("def test_old_stamp_refused(): pass\nassert state.source == HEAD\n")
	if err := os.WriteFile(filepath.Join(dir, "owned.txt"), body, 0600); err != nil {
		t.Fatal(err)
	}
	c := catalog{Schema: "zeo.job-checks/v1", Sources: []source{{Path: "owned.txt", SHA256: digest(body)}}, Checks: []check{
		{ID: "app.old-stamp-refusal", Source: "owned.txt", Contains: "def test_old_stamp_refused", Argv: []string{"python3", "-m", "unittest", "scanner"}},
		{ID: "app.consumer-source-binding", Source: "owned.txt", Contains: "assert state.source == HEAD", Argv: []string{"node", "harness.mjs", "${OUT}"}},
	}}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"baseline.json", "current.json"} {
		if err := os.WriteFile(filepath.Join(dir, name), b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	args := []string{"--baseline", filepath.Join(dir, "baseline.json"), "--baseline-sha256", digest(b), "--catalog", filepath.Join(dir, "current.json"), "--inputs", dir}
	return dir, args, c
}

func writeCurrent(t *testing.T, dir string, c catalog) {
	t.Helper()
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "current.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestLostScannerAndConsumerChecksFailBeforeAnyOutput(t *testing.T) {
	for _, removed := range []int{0, 1} {
		t.Run([]string{"old-stamp-refusal", "consumer-binding"}[removed], func(t *testing.T) {
			dir, args, c := fixture(t)
			c.Checks = append(c.Checks[:removed:removed], c.Checks[removed+1:]...)
			writeCurrent(t, dir, c)
			var out, errs bytes.Buffer
			if run(args, &out, &errs) != 1 || out.Len() != 0 || !strings.Contains(errs.String(), "dropped approved check") {
				t.Fatal("lost check accepted or partial job emitted", out.String(), errs.String())
			}
		})
	}
}

func TestRenamingAnObligationNeedsExplicitRetirement(t *testing.T) {
	dir, args, c := fixture(t)
	c.Checks[0].ID = "app.renamed-refusal"
	writeCurrent(t, dir, c)
	var out, errs bytes.Buffer
	if run(args, &out, &errs) != 1 || out.Len() != 0 {
		t.Fatal("renamed ID silently retired approved obligation")
	}
}

func TestRetirementRequiresNamedMissingIDAndReason(t *testing.T) {
	for _, reason := range []string{"", "   ", "Retired by owner decision example-17; superseded by mandatory source binding"} {
		t.Run(reason, func(t *testing.T) {
			dir, args, c := fixture(t)
			id := c.Checks[0].ID
			c.Checks = c.Checks[1:]
			c.Retired = []retirement{{ID: id, Reason: reason}}
			writeCurrent(t, dir, c)
			var out, errs bytes.Buffer
			code := run(args, &out, &errs)
			if strings.TrimSpace(reason) == "" {
				if code != 1 || out.Len() != 0 {
					t.Fatal("empty retirement reason accepted")
				}
			} else if code != 0 || !strings.Contains(out.String(), "consumer-source-binding") {
				t.Fatal("explicit retirement failed", errs.String())
			}
		})
	}
	for _, id := range []string{"app.unknown", "app.consumer-source-binding"} {
		dir, args, c := fixture(t)
		c.Retired = []retirement{{ID: id, Reason: "example"}}
		writeCurrent(t, dir, c)
		var out, errs bytes.Buffer
		if run(args, &out, &errs) != 1 || out.Len() != 0 {
			t.Fatal("unknown or still-active retirement accepted")
		}
	}
}

func TestRepinnedOwnerFileStillMustContainEachApprovedAnchor(t *testing.T) {
	dir, args, c := fixture(t)
	b := []byte("def test_old_stamp_refused(): pass\n")
	if err := os.WriteFile(filepath.Join(dir, "owned.txt"), b, 0600); err != nil {
		t.Fatal(err)
	}
	c.Sources[0].SHA256 = digest(b)
	writeCurrent(t, dir, c)
	var out, errs bytes.Buffer
	if run(args, &out, &errs) != 1 || out.Len() != 0 || !strings.Contains(errs.String(), "missing pinned check anchor") {
		t.Fatal("new owner file silently lost consumer assertion", errs.String())
	}
}

func TestOrdinaryOwnerHashRepinKeepsUnchangedDefinitionsWithoutAmendment(t *testing.T) {
	dir, args, c := fixture(t)
	p := filepath.Join(dir, "owned.txt")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	b = append(b, []byte("# owner revision: same obligations\n")...)
	if err = os.WriteFile(p, b, 0600); err != nil {
		t.Fatal(err)
	}
	c.Sources[0].SHA256 = digest(b)
	writeCurrent(t, dir, c)
	var out, errs bytes.Buffer
	if run(args, &out, &errs) != 0 {
		t.Fatal("ordinary source hash repin incorrectly required definition amendment", errs.String())
	}
}

func TestChangedOwnerBytesAndChangedBaselineAreRefused(t *testing.T) {
	for _, name := range []string{"owned.txt", "baseline.json"} {
		dir, args, _ := fixture(t)
		if err := os.WriteFile(filepath.Join(dir, name), []byte("replacement"), 0600); err != nil {
			t.Fatal(err)
		}
		var out, errs bytes.Buffer
		if run(args, &out, &errs) != 1 || out.Len() != 0 {
			t.Fatal("changed pinned input accepted", name)
		}
	}
}

func TestGenerateAllChecksOncePerInvocationWithoutExecutingThem(t *testing.T) {
	dir, args, c := fixture(t)
	marker := filepath.Join(dir, "must-not-exist")
	old := c.Checks[0]
	c.Checks[0].Argv = []string{"touch", marker, "a'b", "$(must-stay-literal)"}
	c.Amended = []amendment{{ID: old.ID, Reason: "test-only nonexecuting invocation fixture", PreviousSHA256: definition(old), CurrentSHA256: definition(c.Checks[0])}}
	extra := c.Checks[0]
	extra.ID = "app.another-refusal"
	c.Checks = append(c.Checks, extra)
	writeCurrent(t, dir, c)
	var out, errs bytes.Buffer
	if run(args, &out, &errs) != 0 {
		t.Fatal(errs.String())
	}
	s := out.String()
	if strings.Count(s, "'touch'") != 1 || !strings.Contains(s, "# check app.another-refusal") || !strings.Contains(s, `"${OUT:?required by check catalog}"`) || !strings.Contains(s, `'a'"'"'b'`) || !strings.Contains(s, "'$(must-stay-literal)'") {
		t.Fatal("incomplete or unsafe invocation fragment", s)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("generator executed a check")
	}
}

func TestSameIDCannotReplaceConsumerInvocationWithTrue(t *testing.T) {
	dir, args, c := fixture(t)
	c.Checks[1].Argv = []string{"true"}
	writeCurrent(t, dir, c)
	var out, errs bytes.Buffer
	if run(args, &out, &errs) != 1 || out.Len() != 0 || !strings.Contains(errs.String(), "definition changed") {
		t.Fatal("same ID hid dropped harness invocation", errs.String())
	}
}

func TestDefinitionAmendmentsBindBothDefinitionsAndRequireReason(t *testing.T) {
	for _, fault := range []string{"valid", "reason", "old-digest", "new-digest", "unknown-id", "source", "anchor"} {
		t.Run(fault, func(t *testing.T) {
			dir, args, c := fixture(t)
			old := c.Checks[1]
			c.Checks[1].Argv = []string{"true"}
			a := amendment{ID: old.ID, Reason: "owner decision fixture-18", PreviousSHA256: definition(old), CurrentSHA256: definition(c.Checks[1])}
			switch fault {
			case "reason":
				a.Reason = " "
			case "old-digest":
				a.PreviousSHA256 = strings.Repeat("0", 64)
			case "new-digest":
				a.CurrentSHA256 = strings.Repeat("0", 64)
			case "unknown-id":
				a.ID = "app.unknown"
			case "source":
				c.Checks[1].Source = "other.txt"
			case "anchor":
				c.Checks[1].Contains = "source == HEAD"
			}
			c.Amended = []amendment{a}
			writeCurrent(t, dir, c)
			var out, errs bytes.Buffer
			code := run(args, &out, &errs)
			if fault == "valid" {
				if code != 0 {
					t.Fatal(errs.String())
				}
			} else if code != 1 || out.Len() != 0 {
				t.Fatal("invalid definition amendment accepted", fault)
			}
		})
	}
}

func TestGeneratedCommandsUseValidatedRootDespiteDifferentCallerDirectory(t *testing.T) {
	dir, args, c := fixture(t)
	// All test execution happens only when the test suite runs on its admitted host.
	// The command proves relative reads resolve to A, even when launched from B.
	c.Checks = c.Checks[:1]
	c.Checks[0].Argv = []string{"/bin/sh", "-c", "cat owned.txt"}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "baseline.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	args[3] = digest(b)
	writeCurrent(t, dir, c)
	var out, errs bytes.Buffer
	if run(args, &out, &errs) != 0 {
		t.Fatal(errs.String())
	}
	other := t.TempDir()
	if err = os.WriteFile(filepath.Join(other, "owned.txt"), []byte("WRONG ROOT"), 0600); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("/bin/sh", "-c", out.String())
	command.Dir = other
	got, err := command.CombinedOutput()
	if err != nil {
		t.Fatal(err, string(got))
	}
	want, err := os.ReadFile(filepath.Join(dir, "owned.txt"))
	if err != nil || !bytes.Equal(got, want) {
		t.Fatal("generated invocation ran outside validated root", string(got))
	}
}

func TestDuplicateIDsAndUnknownCatalogFieldsRefuse(t *testing.T) {
	dir, args, c := fixture(t)
	c.Checks = append(c.Checks, c.Checks[0])
	writeCurrent(t, dir, c)
	var out, errs bytes.Buffer
	if run(args, &out, &errs) != 1 || out.Len() != 0 {
		t.Fatal("duplicate obligation accepted")
	}
	_, _, valid := fixture(t)
	b, _ := json.Marshal(valid)
	b = append([]byte(`{"ignored":true,`), b[1:]...)
	if _, err := decode(b); err == nil {
		t.Fatal("unknown catalog field accepted")
	}
}

func TestCatalogRejectsDuplicateAndCaseFoldedKeysAtEveryLevel(t *testing.T) {
	_, _, c := fixture(t)
	raw, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range [][]byte{
		append([]byte(`{"schema":"zeo.job-checks/v1",`), raw[1:]...),
		bytes.Replace(raw, []byte(`"schema"`), []byte(`"SCHEMA"`), 1),
		bytes.Replace(raw, []byte(`"checks":`), []byte(`"CHECKS":[],"checks":`), 1),
		bytes.Replace(raw, []byte(`"id":"app.old-stamp-refusal"`), []byte(`"id":"app.old-stamp-refusal","id":"app.changed"`), 1),
		bytes.Replace(raw, []byte(`"argv":`), []byte(`"ARGV":["true"],"argv":`), 1),
		bytes.Replace(raw, []byte(`"path":`), []byte(`"PATH":`), 1),
		bytes.Replace(raw, []byte(`"retired":null`), []byte(`"retired":[{"id":"app.old-stamp-refusal","reason":"a","REASON":"b"}]`), 1),
		bytes.Replace(raw, []byte(`"amended":null`), []byte(`"amended":[{"id":"app.old-stamp-refusal","reason":"a","previous_sha256":"x","PREVIOUS_SHA256":"y","current_sha256":"z"}]`), 1),
	} {
		if bytes.Equal(bad, raw) {
			t.Fatal("fixture failed to mutate catalog")
		}
		if _, err := decode(bad); err == nil {
			t.Fatal("ambiguous catalog accepted", string(bad))
		}
	}
	if _, err := decode(raw); err != nil {
		t.Fatal("canonical catalog refused", err)
	}
}

func TestRefusedCompositionPublishesNoRunnableFragment(t *testing.T) {
	dir, args, c := fixture(t)
	c.Checks = c.Checks[1:]
	writeCurrent(t, dir, c)
	target := filepath.Join(dir, "checks.sh")
	args = append(args, "--out", target)
	var out, errs bytes.Buffer
	if emit(args, &out, &errs) == 0 {
		t.Fatal("refused composition returned success")
	}
	if out.Len() != 0 {
		t.Fatal("executable prefix emitted to stdout")
	}
	if _, err := os.Lstat(target); !os.IsNotExist(err) {
		t.Fatal("refused composition published a runnable fragment", err)
	}
	command := exec.Command("/bin/sh", target)
	if err := command.Run(); err == nil {
		t.Fatal("refused fragment could be executed successfully")
	}
}

func TestShortFragmentWriteNeverPublishesPrefix(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "checks.sh")
	err := publishFragment(target, []byte("true\necho complete\n"), func(f *os.File, b []byte) (int, error) {
		n, err := f.Write(b[:5])
		if err != nil {
			return n, err
		}
		return n, io.ErrShortWrite
	})
	if err == nil {
		t.Fatal("short write accepted")
	}
	if _, err := os.Lstat(target); !os.IsNotExist(err) {
		t.Fatal("short write published executable prefix", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 0 {
		t.Fatal("temporary prefix not removed", err)
	}
}

func TestCLIAtomicallyPublishesCompleteFragmentWithoutReplacement(t *testing.T) {
	dir, args, _ := fixture(t)
	target := filepath.Join(dir, "checks.sh")
	args = append(args, "--out", target)
	var out, errs bytes.Buffer
	if emit(args, &out, &errs) != 0 || out.Len() != 0 {
		t.Fatal("complete fragment not published", errs.String())
	}
	before, err := os.ReadFile(target)
	if err != nil || !bytes.Contains(before, []byte("consumer-source-binding")) {
		t.Fatal("incomplete output", err)
	}
	if emit(args, &out, &errs) == 0 {
		t.Fatal("existing fragment replaced")
	}
	after, err := os.ReadFile(target)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("existing fragment changed", err)
	}
}
