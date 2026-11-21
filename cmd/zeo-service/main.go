// zeo-service renders candidate Linux service packaging without changing the OS.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/zeroemployeeorg/zeo/internal/install"
)

const usage = "zeo-service --organization ID --seat ID --user USER --group GROUP --pomar-sha256 HEX --memory-bytes N --tasks N --cpu-percent N\n"

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, stdout, stderr io.Writer) int {
	var c install.SessionService
	f := flag.NewFlagSet("zeo-service", flag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.StringVar(&c.Organization, "organization", "", "organization identity")
	f.StringVar(&c.Seat, "seat", "", "seat identity")
	f.StringVar(&c.User, "user", "", "Runtime-allocated OS user")
	f.StringVar(&c.Group, "group", "", "Runtime-allocated OS group")
	f.StringVar(&c.PomarSHA256, "pomar-sha256", "", "qualified Pomar binary digest")
	f.Uint64Var(&c.MemoryBytes, "memory-bytes", 0, "profile memory cap")
	f.Uint64Var(&c.Tasks, "tasks", 0, "profile task cap")
	f.Uint64Var(&c.CPUPercent, "cpu-percent", 0, "profile CPU cap; 100 is one CPU")
	if err := f.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, usage)
			return 0
		}
		fmt.Fprint(stderr, usage)
		return 2
	}
	if f.NArg() != 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	unit, err := install.RenderSessionUnit(c)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if _, err := io.WriteString(stdout, unit); err != nil {
		return 1
	}
	return 0
}
