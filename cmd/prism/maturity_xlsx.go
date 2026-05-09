package main

import (
	"fmt"

	"github.com/grokify/prism/fileutil"
	"github.com/grokify/prism/maturity"
	"github.com/spf13/cobra"
)

var maturityXLSXCmd = &cobra.Command{
	Use:   "xlsx <maturity-spec-file>",
	Short: "Generate Excel report from a maturity specification",
	Long: `Generate an Excel (XLSX) report from a maturity model specification.

The report includes multiple sheets:
  - Requirements:         Enablers with domain, level, status
  - SLOs:                 Criteria with framework columns
  - Framework Mappings:   Detailed framework control mappings
  - Progress:             Assessment status by domain
  - Level Definitions:    M1-M5 level descriptions

Examples:
  prism maturity xlsx spec.json                       # Generate spec.xlsx
  prism maturity xlsx spec.json -o report.xlsx        # Generate report.xlsx`,
	Args: cobra.ExactArgs(1),
	RunE: runMaturityXLSX,
}

var xlsxOutput string

func init() {
	maturityCmd.AddCommand(maturityXLSXCmd)

	maturityXLSXCmd.Flags().StringVarP(&xlsxOutput, "output", "o", "", "Output file (default: <input>.xlsx)")
}

func runMaturityXLSX(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Read and parse maturity spec
	spec, err := maturity.ReadSpecFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read maturity spec: %w", err)
	}

	// Determine output filename
	outputFile := xlsxOutput
	if outputFile == "" {
		// Default to input filename with .xlsx extension
		outputFile = fileutil.ReplaceExtension(filename, ".xlsx")
	}

	// Generate XLSX
	gen := maturity.NewXLSXGenerator(spec)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("failed to generate XLSX: %w", err)
	}

	if err := gen.SaveAs(outputFile); err != nil {
		return fmt.Errorf("failed to save XLSX: %w", err)
	}

	fmt.Printf("Excel report written to %s\n", outputFile)
	return nil
}
