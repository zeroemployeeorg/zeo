// zeo-release performs offline release preflight. It never installs software,
// starts services, contacts component APIs or reads organizational settings.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zeroemployeeorg/zeo/internal/release"
)

const usage = "zeo-release --manifest FILE --arch arm64|amd64 --linux-id ID --linux-version VERSION --profile ID --profile-version VERSION [--artifacts DIR] [--out FILE]\n"

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("zeo-release", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	manifest := flags.String("manifest", "", "release manifest")
	arch := flags.String("arch", "", "target architecture")
	linuxID := flags.String("linux-id", "", "Linux distribution")
	linuxVersion := flags.String("linux-version", "", "Linux version")
	profile := flags.String("profile", "", "Pomar profile identifier")
	profileVersion := flags.String("profile-version", "", "Pomar profile version")
	artifacts := flags.String("artifacts", "", "optional local artifact directory")
	output := flags.String("out", "", "atomic plan output; stdout by default")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, usage)
			return 0
		}
		fmt.Fprint(stderr, usage)
		return 2
	}
	if flags.NArg() != 0 || *manifest == "" || *arch == "" || *linuxID == "" || *linuxVersion == "" || *profile == "" || *profileVersion == "" {
		fmt.Fprint(stderr, usage)
		return 2
	}
	file, err := os.Open(*manifest)
	if err != nil {
		fmt.Fprintln(stderr, "release manifest is unreadable")
		return 1
	}
	m, decodeErr := release.Decode(file)
	closeErr := file.Close()
	if decodeErr != nil {
		fmt.Fprintln(stderr, decodeErr)
		return 1
	}
	if closeErr != nil {
		fmt.Fprintln(stderr, "release manifest read failed")
		return 1
	}
	plan, err := m.ValidateFor(release.Target{
		Platform:    release.Platform{OS: "linux", Architecture: *arch, LinuxID: *linuxID, LinuxVersion: *linuxVersion},
		Environment: release.Environment{Profile: *profile, Version: *profileVersion},
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if *artifacts != "" {
		if err := m.VerifyArtifacts(*artifacts); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		plan.ArtifactsVerified = true
	}
	if *output != "" {
		inputs := []string{*manifest}
		if *artifacts != "" {
			for _, component := range m.Components {
				inputs = append(inputs, filepath.Join(*artifacts, component.Artifact.Filename))
			}
		}
		for _, input := range inputs {
			if samePath(*output, input) {
				fmt.Fprintln(stderr, "plan output would replace a release input")
				return 1
			}
		}
	}
	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		fmt.Fprintln(stderr, "release plan encoding failed")
		return 1
	}
	data = append(data, '\n')
	if *output == "" {
		if _, err := stdout.Write(data); err != nil {
			return 1
		}
	} else if err := writePlan(*output, data); err != nil {
		fmt.Fprintln(stderr, "release plan could not be written")
		return 1
	}
	return 0
}

func samePath(a, b string) bool {
	absoluteA, errA := filepath.Abs(a)
	absoluteB, errB := filepath.Abs(b)
	if errA != nil || errB != nil || absoluteA == absoluteB {
		return true
	}
	infoA, errA := os.Stat(a)
	infoB, errB := os.Stat(b)
	return errA == nil && errB == nil && os.SameFile(infoA, infoB)
}

// All validation and optional artifact checks finish before touching output.
func writePlan(path string, data []byte) error {
	file, err := os.CreateTemp(filepath.Dir(path), ".zeo-release-plan-*")
	if err != nil {
		return err
	}
	name := file.Name()
	defer os.Remove(name)
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}
