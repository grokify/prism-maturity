package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-intelligence"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a PRISM document",
	Long: `Validate a PRISM document against the schema and check for errors.

Examples:
  prism validate prism.json
  prism validate operations-metrics.json`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var doc prism.PRISMDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate document
	errs := doc.Validate()
	if errs.HasErrors() {
		fmt.Println("Validation errors:")
		for _, e := range errs {
			fmt.Printf("  - %s\n", e.Error())
		}
		return fmt.Errorf("document has %d validation errors", len(errs))
	}

	// Print summary
	fmt.Printf("✓ %s is valid\n", filename)
	fmt.Printf("  Metrics: %d\n", len(doc.Metrics))

	// Count by domain
	domainCounts := make(map[string]int)
	for _, m := range doc.Metrics {
		domainCounts[m.Domain]++
	}
	for domain, count := range domainCounts {
		fmt.Printf("    %s: %d\n", domain, count)
	}

	// Count by stage
	stageCounts := make(map[string]int)
	for _, m := range doc.Metrics {
		stageCounts[m.Stage]++
	}
	fmt.Println("  By stage:")
	for _, stage := range prism.AllStages() {
		if count, ok := stageCounts[stage]; ok {
			fmt.Printf("    %s: %d\n", stage, count)
		}
	}

	if doc.Maturity != nil {
		fmt.Printf("  Maturity cells: %d\n", len(doc.Maturity.Cells))
	}

	return nil
}
