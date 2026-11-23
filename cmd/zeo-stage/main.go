// zeo-stage stages a local pinned Pomar executable. It does not activate it.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/zeroemployeeorg/zeo/internal/install"
)

const usage = "zeo-stage --pomar FILE --sha256 HEX --size BYTES --artifact-root DIR\n"

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, stdout, stderr io.Writer) int {
	f := flag.NewFlagSet("zeo-stage", flag.ContinueOnError)
	f.SetOutput(io.Discard)
	source := f.String("pomar", "", "absolute local executable source")
	digest := f.String("sha256", "", "qualified release pin")
	size := f.Int64("size", 0, "exact release artifact size")
	rootPath := f.String("artifact-root", "", "existing trusted staging directory")
	if err := f.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, usage)
			return 0
		}
		fmt.Fprint(stderr, usage)
		return 2
	}
	if f.NArg() != 0 || !filepath.IsAbs(*source) || !filepath.IsAbs(*rootPath) {
		fmt.Fprint(stderr, usage)
		return 2
	}
	// Nonblocking and no-follow refuse FIFO/symlink substitution at open.
	file, err := os.OpenFile(*source, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		fmt.Fprintln(stderr, "source executable is unreadable")
		return 1
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		fmt.Fprintln(stderr, "source must be a regular file")
		return 1
	}
	root, err := os.OpenRoot(*rootPath)
	if err != nil {
		fmt.Fprintln(stderr, "artifact root is unreadable")
		return 1
	}
	defer root.Close()
	path, err := install.StagePomarArtifact(root, file, *digest, *size)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if _, err := fmt.Fprintln(stdout, path); err != nil {
		return 1
	}
	return 0
}
