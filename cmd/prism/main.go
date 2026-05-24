// Package main provides the splan CLI tool for working with PRISM documents.
package main

import (
	"fmt"
	"os"

	"github.com/grokify/prism-maturity/cli"
)

func main() {
	// Override the command name for standalone use
	cli.RootCmd.Use = "splan"
	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
