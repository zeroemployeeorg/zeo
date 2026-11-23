// zeo-job-checks composes checks from a reviewed catalog. It emits shell source;
// it never executes checks. Catalogs and retirement decisions are trusted inputs.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

type source struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type check struct {
	ID       string   `json:"id"`
	Source   string   `json:"source"`
	Contains string   `json:"contains"`
	Argv     []string `json:"argv"`
}

type retirement struct {
	ID     string `json:"id"`
	Reason string `json:"reason"`
}

type amendment struct {
	ID             string `json:"id"`
	Reason         string `json:"reason"`
	PreviousSHA256 string `json:"previous_sha256"`
	CurrentSHA256  string `json:"current_sha256"`
}

type catalog struct {
	Schema  string       `json:"schema"`
	Sources []source     `json:"sources"`
	Checks  []check      `json:"checks"`
	Retired []retirement `json:"retired"`
	Amended []amendment  `json:"amended"`
}

var idPattern = regexp.MustCompile(`^[a-z][a-z0-9.-]{1,127}$`)
var variablePattern = regexp.MustCompile(`^\$\{([A-Z][A-Z0-9_]*)\}$`)

func digest(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }

func definition(k check) string { b, _ := json.Marshal(k); return digest(b) }

// Validate keys before encoding/json can discard duplicates or match struct
// tags case-insensitively. Walk every nested catalog object, not only its root.
func canonicalCatalog(raw []byte) error {
	d := json.NewDecoder(bytes.NewReader(raw))
	if err := catalogValue(d, reflect.TypeOf(catalog{})); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		return errors.New("catalog has trailing data")
	}
	return nil
}
func catalogValue(d *json.Decoder, typ reflect.Type) error {
	token, err := d.Token()
	if err != nil {
		return err
	}
	switch typ.Kind() {
	case reflect.Struct:
		if token != json.Delim('{') {
			return errors.New("catalog object required")
		}
		fields := map[string]reflect.Type{}
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fields[strings.Split(field.Tag.Get("json"), ",")[0]] = field.Type
		}
		seen := map[string]bool{}
		for d.More() {
			keyToken, err := d.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			field, exists := fields[key]
			if !ok || !exists || seen[key] {
				return errors.New("duplicate or noncanonical catalog key")
			}
			seen[key] = true
			if err := catalogValue(d, field); err != nil {
				return err
			}
		}
		close, err := d.Token()
		if err != nil || close != json.Delim('}') {
			return errors.New("catalog object incomplete")
		}
	case reflect.Slice:
		if token == nil {
			return nil
		} // Optional arrays marshal as null.
		if token != json.Delim('[') {
			return errors.New("catalog array required")
		}
		for d.More() {
			if err := catalogValue(d, typ.Elem()); err != nil {
				return err
			}
		}
		close, err := d.Token()
		if err != nil || close != json.Delim(']') {
			return errors.New("catalog array incomplete")
		}
	default:
		if _, nested := token.(json.Delim); nested {
			return errors.New("catalog scalar required")
		}
	}
	return nil // Typed decoding and semantic validation follow this key check.
}

func decode(b []byte) (catalog, error) {
	var c catalog
	if err := canonicalCatalog(b); err != nil {
		return c, err
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if err := d.Decode(&c); err != nil {
		return c, err
	}
	var extra any
	if err := d.Decode(&extra); !errors.Is(err, io.EOF) {
		return c, errors.New("catalog has trailing data")
	}
	if c.Schema != "zeo.job-checks/v1" || len(c.Checks) == 0 {
		return c, errors.New("catalog schema or checks missing")
	}
	seen := map[string]bool{}
	for _, k := range c.Checks {
		if !idPattern.MatchString(k.ID) || seen[k.ID] || k.Source == "" || strings.TrimSpace(k.Contains) == "" || len(k.Argv) == 0 {
			return c, errors.New("invalid or duplicate check")
		}
		seen[k.ID] = true
		for _, a := range k.Argv {
			if a == "" || strings.ContainsRune(a, 0) {
				return c, errors.New("invalid command argument")
			}
		}
		if variablePattern.MatchString(k.Argv[0]) {
			return c, errors.New("executable must be literal")
		}
	}
	return c, nil
}

func compose(previous, current catalog, read func(string) ([]byte, error)) ([]byte, error) {
	active := map[string]bool{}
	for _, k := range current.Checks {
		active[k.ID] = true
	}
	prior := map[string]bool{}
	for _, k := range previous.Checks {
		prior[k.ID] = true
	}
	// Keeping an ID is insufficient if its only invocation or assertion changes.
	oldDefinitions := map[string]string{}
	for _, k := range previous.Checks {
		oldDefinitions[k.ID] = definition(k)
	}
	newDefinitions := map[string]string{}
	for _, k := range current.Checks {
		newDefinitions[k.ID] = definition(k)
	}
	amended := map[string]bool{}
	for _, a := range current.Amended {
		old, had := oldDefinitions[a.ID]
		next, has := newDefinitions[a.ID]
		if !had || !has || old == next || amended[a.ID] || strings.TrimSpace(a.Reason) == "" || a.PreviousSHA256 != old || a.CurrentSHA256 != next {
			return nil, fmt.Errorf("invalid reviewed amendment: %s", a.ID)
		}
		amended[a.ID] = true
	}
	for _, k := range current.Checks {
		if old, had := oldDefinitions[k.ID]; had && old != definition(k) && !amended[k.ID] {
			return nil, fmt.Errorf("approved check definition changed without reviewed amendment: %s", k.ID)
		}
	}
	retired := map[string]bool{}
	for _, k := range current.Retired {
		if !prior[k.ID] || active[k.ID] || retired[k.ID] || strings.TrimSpace(k.Reason) == "" {
			return nil, fmt.Errorf("invalid retirement: %s", k.ID)
		}
		retired[k.ID] = true
	}
	for _, k := range previous.Checks {
		if !active[k.ID] && !retired[k.ID] {
			return nil, fmt.Errorf("dropped approved check without retirement: %s", k.ID)
		}
	}
	files := map[string][]byte{}
	for _, s := range current.Sources {
		if _, exists := files[s.Path]; exists {
			return nil, errors.New("duplicate source pin")
		}
		b, err := read(s.Path)
		if err != nil || digest(b) != s.SHA256 {
			return nil, fmt.Errorf("source pin mismatch: %s", s.Path)
		}
		files[s.Path] = b
	}
	var out bytes.Buffer
	out.WriteString("# Generated by zeo-job-checks; reviewed catalog required.\nset -eu\n")
	invoked := map[string]bool{}
	for _, k := range current.Checks {
		b, pinned := files[k.Source]
		if !pinned || !bytes.Contains(b, []byte(k.Contains)) {
			return nil, fmt.Errorf("missing pinned check anchor: %s", k.ID)
		}
		fmt.Fprintf(&out, "# check %s\n", k.ID)
		key, _ := json.Marshal(k.Argv)
		if invoked[string(key)] {
			continue
		}
		invoked[string(key)] = true
		for i, a := range k.Argv {
			if i > 0 {
				out.WriteByte(' ')
			}
			if v := variablePattern.FindStringSubmatch(a); v != nil {
				fmt.Fprintf(&out, "\"${%s:?required by check catalog}\"", v[1])
			} else {
				out.WriteString("'" + strings.ReplaceAll(a, "'", "'\"'\"'") + "'")
			}
		}
		out.WriteByte('\n')
	}
	return out.Bytes(), nil
}

func boundedRead(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, (1<<20)+1))
	if len(b) > 1<<20 {
		return nil, errors.New("catalog exceeds 1MiB")
	}
	return b, err
}

func run(args []string, stdout, stderr io.Writer) int {
	f := flag.NewFlagSet("zeo-job-checks", flag.ContinueOnError)
	f.SetOutput(stderr)
	baseline := f.String("baseline", "", "last approved catalog")
	pin := f.String("baseline-sha256", "", "reviewed baseline digest")
	current := f.String("catalog", "", "owned current catalog")
	dir := f.String("inputs", "", "root of pinned current inputs")
	if err := f.Parse(args); err != nil {
		return 2
	}
	if f.NArg() != 0 || *baseline == "" || *pin == "" || *current == "" || *dir == "" {
		fmt.Fprintln(stderr, "baseline, baseline-sha256, catalog and inputs are required")
		return 2
	}
	p, err := boundedRead(*baseline)
	if err != nil || digest(p) != *pin {
		fmt.Fprintln(stderr, "approved baseline digest mismatch")
		return 1
	}
	old, err := decode(p)
	if err != nil {
		fmt.Fprintln(stderr, "invalid approved baseline:", err)
		return 1
	}
	b, err := boundedRead(*current)
	if err != nil {
		fmt.Fprintln(stderr, "current catalog unreadable")
		return 1
	}
	next, err := decode(b)
	if err != nil {
		fmt.Fprintln(stderr, "invalid current catalog:", err)
		return 1
	}
	absolute, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintln(stderr, "input root unresolved")
		return 1
	}
	canonical, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		fmt.Fprintln(stderr, "input root unresolved")
		return 1
	}
	root, err := os.OpenRoot(canonical)
	if err != nil {
		fmt.Fprintln(stderr, "input root unreadable")
		return 1
	}
	defer root.Close()
	out, err := compose(old, next, func(path string) ([]byte, error) {
		f, err := root.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil || !info.Mode().IsRegular() || info.Size() > 16<<20 {
			return nil, errors.New("source must be a regular file under16MiB")
		}
		return io.ReadAll(f)
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	// The same canonical root supplies both checked source bytes and execution cwd.
	// Review still owns invocation semantics and concurrent-mutation assumptions.
	bound := "cd -- '" + strings.ReplaceAll(canonical, "'", "'\"'\"'") + "' || exit 1\n"
	if _, err := stdout.Write(append([]byte(bound), out...)); err != nil {
		return 1
	}
	return 0
}

// Publish only a complete composition, on the destination filesystem. Never
// replace an old fragment: the caller chooses a fresh output and gates on exit.
func publishFragment(path string, raw []byte, write func(*os.File, []byte) (int, error)) error {
	if !filepath.IsAbs(path) {
		return errors.New("output path must be absolute")
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".zeo-job-checks-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	n, err := write(f, raw)
	if err != nil {
		return err
	}
	if n != len(raw) {
		return io.ErrShortWrite
	}
	if err = f.Sync(); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	// link is no-replace; the temporary file is complete before publication.
	if err = os.Link(f.Name(), path); err != nil {
		return err
	}
	return nil
}

func emit(args []string, stdout, stderr io.Writer) int {
	var output string
	filtered := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--help" || arg == "-help" || arg == "-h" {
			fmt.Fprintln(stdout, "zeo-job-checks --baseline FILE --baseline-sha256 HEX --catalog FILE --inputs DIR --out ABSOLUTE_NEW_FILE")
			return 0
		}
		if arg == "--out" || arg == "-out" {
			if output != "" || i+1 == len(args) {
				fmt.Fprintln(stderr, "one output path is required")
				return 2
			}
			i++
			output = args[i]
		} else if strings.HasPrefix(arg, "--out=") || strings.HasPrefix(arg, "-out=") {
			if output != "" {
				fmt.Fprintln(stderr, "one output path is required")
				return 2
			}
			output = strings.SplitN(arg, "=", 2)[1]
		} else {
			filtered = append(filtered, arg)
		}
	}
	if output == "" || !filepath.IsAbs(output) {
		fmt.Fprintln(stderr, "--out ABSOLUTE_NEW_FILE required; stdout is not an executable fragment")
		return 2
	}
	var fragment bytes.Buffer
	if code := run(filtered, &fragment, stderr); code != 0 {
		return code
	}
	if err := publishFragment(output, fragment.Bytes(), func(f *os.File, b []byte) (int, error) { return f.Write(b) }); err != nil {
		fmt.Fprintln(stderr, "fragment publication failed; no retry or execution without reconciliation")
		return 1
	}
	return 0
}
func main() { os.Exit(emit(os.Args[1:], os.Stdout, os.Stderr)) }
