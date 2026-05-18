package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grokify/prism-intelligence/export"
	"github.com/spf13/cobra"
)

// Export commands

var exportOKRCmd = &cobra.Command{
	Use:   "okr <prism-file>",
	Short: "Export PRISM as OKR document for structured-plan",
	Long: `Export PRISM goals, SLOs, and phases as an OKR document compatible
with the structured-plan repository.

Mapping:
  PRISM Goal        → OKR Objective
  PRISM SLO         → OKR Key Result
  PRISM Phase       → PhaseTargets in Key Results
  PRISM Initiative  → Referenced in Key Result descriptions

Examples:
  prism export okr prism.json -o roadmap.okr.json
  prism export okr prism.json -o ./out/ -f json,markdown`,
	Args: cobra.ExactArgs(1),
	RunE: runExportOKR,
}

var exportV2MOMCmd = &cobra.Command{
	Use:   "v2mom <prism-file>",
	Short: "Export PRISM as V2MOM document for structured-plan",
	Long: `Export PRISM goals, SLOs, and phases as a V2MOM document compatible
with the structured-plan repository.

Mapping:
  PRISM Goal        → V2MOM Method
  PRISM SLO         → V2MOM Measure
  PRISM Phase       → Method timeline
  PRISM Initiative  → Referenced in Method description

Examples:
  prism export v2mom prism.json -o roadmap.v2mom.json`,
	Args: cobra.ExactArgs(1),
	RunE: runExportV2MOM,
}

func init() {
	exportCmd.AddCommand(exportOKRCmd)
	exportCmd.AddCommand(exportV2MOMCmd)

	// Add output flag to OKR and V2MOM commands
	// These use exportOutputDir from export.go
	exportOKRCmd.Flags().StringVarP(&exportOutputDir, "output", "o", "", "Output file or directory")
	exportV2MOMCmd.Flags().StringVarP(&exportOutputDir, "output", "o", "", "Output file or directory")
}

func runExportOKR(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	okrDoc := export.ConvertToOKR(doc)

	// Determine output path
	outputPath := exportOutputDir
	if outputPath == "" {
		// If no output dir specified, output to stdout
		return outputOKRToStdout(okrDoc)
	}

	// Check if outputPath is a directory or file
	info, err := os.Stat(outputPath)
	if err == nil && info.IsDir() {
		outputPath = filepath.Join(outputPath, "roadmap.okr.json")
	}

	return outputOKRToFile(okrDoc, outputPath)
}

func outputOKRToStdout(okrDoc *export.OKRDocument) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(okrDoc)
}

func outputOKRToFile(okrDoc *export.OKRDocument, path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(okrDoc); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	fmt.Printf("Exported: %s\n", path)
	return nil
}

// V2MOM export

func runExportV2MOM(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	v2momDoc := export.ConvertToV2MOM(doc)

	// Determine output path
	outputPath := exportOutputDir
	if outputPath == "" {
		return outputV2MOMToStdout(v2momDoc)
	}

	// Check if outputPath is a directory or file
	info, err := os.Stat(outputPath)
	if err == nil && info.IsDir() {
		outputPath = filepath.Join(outputPath, "roadmap.v2mom.json")
	}

	return outputV2MOMToFile(v2momDoc, outputPath)
}

func outputV2MOMToStdout(v2momDoc *export.V2MOMDocument) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v2momDoc)
}

func outputV2MOMToFile(v2momDoc *export.V2MOMDocument, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v2momDoc); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	fmt.Printf("Exported: %s\n", path)
	return nil
}
