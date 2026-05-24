// Command-line tool for generating maturity reports.
//
// Usage:
//
//	go run maturity/cmd/main.go maturity-models/security.json -o security-maturity.xlsx
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/grokify/prism-maturity/maturity"
)

func main() {
	output := flag.String("o", "maturity-report.xlsx", "Output XLSX file")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: maturity-xlsx <spec.json> [-o output.xlsx]")
		os.Exit(1)
	}

	specFile := flag.Arg(0)

	fmt.Printf("Reading spec from: %s\n", specFile)
	spec, err := maturity.ReadSpecFile(specFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading spec: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d domain(s)\n", len(spec.Domains))
	for name, domain := range spec.Domains {
		fmt.Printf("  - %s: %d levels\n", name, len(domain.Levels))

		// Count criteria and enablers
		var criteria, enablers int
		for _, level := range domain.Levels {
			criteria += len(level.Criteria)
			enablers += len(level.Enablers)
		}
		fmt.Printf("    Criteria: %d, Enablers: %d\n", criteria, enablers)
	}

	fmt.Printf("Generating XLSX: %s\n", *output)
	gen := maturity.NewXLSXGenerator(spec)
	if err := gen.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating XLSX: %v\n", err)
		os.Exit(1)
	}

	if err := gen.SaveAs(*output); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving XLSX: %v\n", err)
		os.Exit(1)
	}

	info, _ := os.Stat(*output)
	fmt.Printf("Generated %s (%d bytes)\n", *output, info.Size())
	fmt.Println("Sheets created:")
	fmt.Println("  - Requirements: All enablers with status")
	fmt.Println("  - SLOs: All criteria with current values")
	fmt.Println("  - Progress: Summary by domain and level")
	fmt.Println("  - Level Definitions: M1-M5 descriptions")
}
