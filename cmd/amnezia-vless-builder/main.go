// Package main is the entry point for the Amnezia VLESS profile builder CLI.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"amnezia-vless-builder/internal/keyconv"
)

func main() {
	flag.Usage = func() {
		bin := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s <source_file>\n", bin)
		fmt.Fprintln(os.Stderr, "Decode an encrypted provisioning file into a usable profile.")
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	outPath, err := keyconv.ProcessFile(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Processing failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Profile saved as %s\n", outPath)
}
