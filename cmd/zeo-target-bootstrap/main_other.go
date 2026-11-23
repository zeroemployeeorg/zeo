//go:build !linux

package main

import "os"

func main() {
	os.Stderr.WriteString("Linux target required\n")
	os.Exit(1)
}
