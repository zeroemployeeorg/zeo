// Package release checks distribution assembly inputs. It does not grant
// deployment authority or establish that a component or organization works.
package release

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	MaxManifestBytes       = 1 << 20
	MaxArtifactBytes int64 = 1 << 40
)

var (
	namePattern      = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)
	fieldNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	pinNamePattern   = regexp.MustCompile(`^[A-Za-z0-9@][A-Za-z0-9@._+:/-]{0,127}$`)
	versionPattern   = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+-]{0,127}$`)
	digestPattern    = regexp.MustCompile(`^[0-9a-f]{64}$`)
	revisionPattern  = regexp.MustCompile(`^([0-9a-f]{40}|[0-9a-f]{64})$`)
)

// RequiredCapabilities describes product responsibilities, not daemon counts.
// One component may own multiple capabilities. Requires may contain cycles;
// InstallAfter is a separate, acyclic installation-order constraint.
var requiredCapabilities = []string{
	"runtime", "coordination", "messaging", "agenda", "controlling", "identity", "application",
}

type Platform struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	LinuxID      string `json:"linux_id"`
	LinuxVersion string `json:"linux_version"`
}

type Environment struct {
	Profile string `json:"profile"`
	Version string `json:"version"`
}

type VersionPin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type StateSchema struct {
	Name           string   `json:"name"`
	Version        string   `json:"version"`
	CompatibleFrom []string `json:"compatible_from"`
	BackupRequired bool     `json:"backup_required"`
}

type Artifact struct {
	Filename  string `json:"filename"`
	URL       string `json:"url"`
	SHA256    string `json:"sha256"`
	SizeBytes int64  `json:"size_bytes"`
}

type Component struct {
	Name                string        `json:"name"`
	Version             string        `json:"version"`
	SourceRepository    string        `json:"source_repository"`
	SourceRevision      string        `json:"source_revision"`
	Provides            []string      `json:"provides"`
	Requires            []string      `json:"requires"`
	InstallAfter        []string      `json:"install_after"`
	Artifact            Artifact      `json:"artifact"`
	Contracts           []VersionPin  `json:"contracts"`
	StateSchemas        []StateSchema `json:"state_schemas"`
	RuntimeDependencies []VersionPin  `json:"runtime_dependencies"`
}

type Provenance struct {
	Recipe           string `json:"recipe"`
	SourceRepository string `json:"source_repository"`
	SourceRevision   string `json:"source_revision"`
	InputsSHA256     string `json:"inputs_sha256"`
}

// Manifest pins a release for one platform and environment profile. Version
// identifiers are opaque: upgrade compatibility is always listed explicitly.
// Organization settings, credentials and experiment treatments are not inputs.
type Manifest struct {
	Schema      string       `json:"schema"`
	Version     string       `json:"version"`
	Platform    Platform     `json:"platform"`
	Environment Environment  `json:"environment"`
	Components  []Component  `json:"components"`
	AgentTools  []VersionPin `json:"agent_tools"`
	UpgradeFrom []string     `json:"upgrade_from"`
	Provenance  Provenance   `json:"provenance"`
}

type Target struct {
	Platform    Platform
	Environment Environment
}

type Plan struct {
	Schema            string      `json:"schema"`
	ReleaseVersion    string      `json:"release_version"`
	Platform          Platform    `json:"platform"`
	Environment       Environment `json:"environment"`
	MetadataSHA256    string      `json:"metadata_sha256"`
	InstallOrder      []string    `json:"install_order"`
	ArtifactsVerified bool        `json:"artifacts_verified"`
	ReadyForInstall   bool        `json:"ready_for_install"`
}

// Decode rejects ambiguous JSON, unknown fields and oversized documents. Errors
// deliberately omit untrusted values, which could themselves contain secrets.
func Decode(r io.Reader) (Manifest, error) {
	data, err := io.ReadAll(io.LimitReader(r, MaxManifestBytes+1))
	if err != nil || len(data) > MaxManifestBytes {
		return Manifest{}, errors.New("manifest unreadable or exceeds size limit")
	}
	keys := json.NewDecoder(bytes.NewReader(data))
	if err := uniqueJSON(keys, 0); err != nil {
		return Manifest{}, err
	}
	if _, err := keys.Token(); err != io.EOF {
		return Manifest{}, errors.New("manifest must contain exactly one JSON value")
	}
	var m Manifest
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&m); err != nil {
		return Manifest{}, errors.New("manifest has unknown fields or invalid field types")
	}
	return m, nil
}

func uniqueJSON(d *json.Decoder, depth int) error {
	if depth > 32 {
		return errors.New("manifest nesting exceeds limit")
	}
	token, err := d.Token()
	if err != nil {
		return errors.New("manifest JSON is invalid")
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for d.More() {
			t, err := d.Token()
			key, ok := t.(string)
			if err != nil || !ok || !fieldNamePattern.MatchString(key) || seen[key] {
				return errors.New("manifest JSON object keys must be canonical and unique")
			}
			seen[key] = true
			if err := uniqueJSON(d, depth+1); err != nil {
				return err
			}
		}
	case '[':
		for d.More() {
			if err := uniqueJSON(d, depth+1); err != nil {
				return err
			}
		}
	default:
		return errors.New("manifest JSON is invalid")
	}
	_, err = d.Token()
	if err != nil {
		return errors.New("manifest JSON is invalid")
	}
	return nil
}

func validPin(p VersionPin) bool {
	return pinNamePattern.MatchString(p.Name) && versionPattern.MatchString(p.Version)
}

func validVersions(values []string) bool {
	if values == nil {
		return false
	}
	seen := map[string]bool{}
	for _, v := range values {
		if !versionPattern.MatchString(v) || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}

func validURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Opaque != "" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	if host == "localhost" || strings.HasSuffix(host, ".localhost") || strings.HasSuffix(host, ".local") || strings.HasSuffix(host, ".internal") {
		return false
	}
	if ip := net.ParseIP(host); ip != nil && (!ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast()) {
		return false
	}
	return true
}

// ValidateFor checks metadata compatibility, not artifact provenance, public
// availability, executable architecture or product acceptance. Those require
// independent release verification. All Linux/profile target fields must match.
func (m Manifest) ValidateFor(target Target) (Plan, error) {
	bad := func(reason string) (Plan, error) { return Plan{}, errors.New(reason) }
	if m.Schema != "zeo.release/v1" || !versionPattern.MatchString(m.Version) {
		return bad("release schema or version is invalid")
	}
	if m.Platform.OS != "linux" || (m.Platform.Architecture != "arm64" && m.Platform.Architecture != "amd64") || !namePattern.MatchString(m.Platform.LinuxID) || !versionPattern.MatchString(m.Platform.LinuxVersion) {
		return bad("release platform is invalid")
	}
	if !namePattern.MatchString(m.Environment.Profile) || !versionPattern.MatchString(m.Environment.Version) || m.Platform != target.Platform || m.Environment != target.Environment {
		return bad("release does not match the requested Linux target and environment profile")
	}
	if !namePattern.MatchString(m.Provenance.Recipe) || !validURL(m.Provenance.SourceRepository) || !revisionPattern.MatchString(m.Provenance.SourceRevision) || !digestPattern.MatchString(m.Provenance.InputsSHA256) {
		return bad("build provenance is incomplete")
	}
	if !validVersions(m.UpgradeFrom) || m.AgentTools == nil || len(m.Components) == 0 {
		return bad("release upgrade paths, agent tools and components must be explicit")
	}
	agents := map[string]bool{}
	for _, pin := range m.AgentTools {
		if !validPin(pin) || agents[pin.Name] {
			return bad("agent tool pins are invalid or duplicated")
		}
		agents[pin.Name] = true
	}
	components := map[string]Component{}
	capabilities, filenames, contracts, dependencies := map[string]bool{}, map[string]bool{}, map[string]string{}, map[string]string{}
	for _, c := range m.Components {
		if !namePattern.MatchString(c.Name) || !versionPattern.MatchString(c.Version) || !validURL(c.SourceRepository) || !revisionPattern.MatchString(c.SourceRevision) {
			return bad("component identity or source pin is invalid")
		}
		if _, exists := components[c.Name]; exists {
			return bad("component names must be unique")
		}
		components[c.Name] = c
		if c.Provides == nil || c.Requires == nil || c.InstallAfter == nil || len(c.Contracts) == 0 || c.StateSchemas == nil || c.RuntimeDependencies == nil {
			return bad("component capabilities, ordering, contracts, state schemas and runtime dependencies must be explicit")
		}
		for _, cap := range c.Provides {
			if !namePattern.MatchString(cap) || capabilities[cap] {
				return bad("capabilities must have one explicit owner")
			}
			capabilities[cap] = true
		}
		a := c.Artifact
		if a.Filename == "" || a.Filename == "." || a.Filename == ".." || filepath.Base(a.Filename) != a.Filename || strings.ContainsAny(a.Filename, `\/:`) || !versionPattern.MatchString(a.Filename) || filenames[a.Filename] || !digestPattern.MatchString(a.SHA256) || a.SizeBytes <= 0 || a.SizeBytes > MaxArtifactBytes || !validURL(a.URL) {
			return bad("artifact identity, public URL, size or digest is invalid")
		}
		filenames[a.Filename] = true
		for _, pins := range []struct {
			values []VersionPin
			common map[string]string
		}{{c.Contracts, contracts}, {c.RuntimeDependencies, dependencies}} {
			seen := map[string]bool{}
			for _, pin := range pins.values {
				if !validPin(pin) || seen[pin.Name] || (pins.common[pin.Name] != "" && pins.common[pin.Name] != pin.Version) {
					return bad("contract or runtime dependency versions conflict")
				}
				seen[pin.Name], pins.common[pin.Name] = true, pin.Version
			}
		}
		seenSchemas := map[string]bool{}
		for _, schema := range c.StateSchemas {
			if !pinNamePattern.MatchString(schema.Name) || !versionPattern.MatchString(schema.Version) || !validVersions(schema.CompatibleFrom) || seenSchemas[schema.Name] {
				return bad("state schema compatibility is invalid or ambiguous")
			}
			seenSchemas[schema.Name] = true
		}
	}
	for _, cap := range requiredCapabilities {
		if !capabilities[cap] {
			return bad("release is missing a required product capability")
		}
	}
	for _, c := range m.Components {
		for _, refs := range [][]string{c.Requires, c.InstallAfter} {
			seen := map[string]bool{}
			for _, ref := range refs {
				if _, exists := components[ref]; !exists || ref == c.Name || seen[ref] {
					return bad("component dependency is unknown, duplicated or self-referential")
				}
				seen[ref] = true
			}
		}
	}
	order, done := []string{}, map[string]bool{}
	for len(order) < len(components) {
		available := []string{}
		for name, c := range components {
			if done[name] {
				continue
			}
			ready := true
			for _, before := range c.InstallAfter {
				if !done[before] {
					ready = false
				}
			}
			if ready {
				available = append(available, name)
			}
		}
		if len(available) == 0 {
			return bad("installation ordering contains a cycle")
		}
		sort.Strings(available)
		for _, name := range available {
			order = append(order, name)
			done[name] = true
		}
	}
	// This binds the plan to this typed metadata representation, not to a
	// signed/raw manifest. It is not a release trust anchor or attestation.
	encoded, err := json.Marshal(m)
	if err != nil {
		return bad("release metadata could not be encoded")
	}
	digest := sha256.Sum256(encoded)
	// Validation can never promote metadata to an install/health verdict.
	return Plan{Schema: "zeo.release-plan/v1", ReleaseVersion: m.Version,
		Platform: m.Platform, Environment: m.Environment,
		MetadataSHA256: hex.EncodeToString(digest[:]),
		InstallOrder:   order, ReadyForInstall: false}, nil
}

// VerifyArtifacts checks regular local artifact bytes against the manifest.
// It never extracts an archive, fetches a URL or executes a component.
func (m Manifest) VerifyArtifacts(directory string) error {
	for _, c := range m.Components {
		if !namePattern.MatchString(c.Name) || !versionPattern.MatchString(c.Artifact.Filename) || filepath.Base(c.Artifact.Filename) != c.Artifact.Filename || strings.ContainsAny(c.Artifact.Filename, `\/:`) || !digestPattern.MatchString(c.Artifact.SHA256) || c.Artifact.SizeBytes <= 0 || c.Artifact.SizeBytes > MaxArtifactBytes {
			return errors.New("artifact metadata must be valid before reading files")
		}
		path := filepath.Join(directory, c.Artifact.Filename)
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() {
			return errors.New("artifact must be an available regular file")
		}
		file, err := os.Open(path)
		if err != nil {
			return errors.New("artifact is unreadable")
		}
		opened, err := file.Stat()
		if err != nil || !os.SameFile(info, opened) {
			file.Close()
			return errors.New("artifact changed before reading")
		}
		hash := sha256.New()
		n, readErr := io.Copy(hash, io.LimitReader(file, c.Artifact.SizeBytes+1))
		final, statErr := file.Stat()
		closeErr := file.Close()
		named, nameErr := os.Lstat(path)
		if readErr != nil || statErr != nil || closeErr != nil || nameErr != nil || !named.Mode().IsRegular() || !os.SameFile(opened, named) || !opened.ModTime().Equal(final.ModTime()) || opened.Size() != final.Size() || n != c.Artifact.SizeBytes || hex.EncodeToString(hash.Sum(nil)) != c.Artifact.SHA256 {
			return fmt.Errorf("artifact verification failed for component %s", c.Name)
		}
	}
	return nil
}
