package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-intelligence"
	"github.com/grokify/prism-intelligence/scaffold"
	"github.com/spf13/cobra"
)

var (
	initDomain string
	initOutput string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new PRISM document",
	Long: `Create a new PRISM document scaffold with default structure.

Examples:
  prism init                          # Create default prism.json with operations metrics
  prism init -d operations -o ops.json # Create ops-focused document
  prism init -d security -o sec.json  # Create security-focused document`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVarP(&initDomain, "domain", "d", "", "Focus domain (operations or security)")
	initCmd.Flags().StringVarP(&initOutput, "output", "o", "prism.json", "Output file path")
}

func runInit(cmd *cobra.Command, args []string) error {
	var doc *prism.PRISMDocument

	// Create document based on domain selection
	switch initDomain {
	case "":
		// Default: both security and operations
		doc = scaffold.NewDocument(prism.DomainSecurity, prism.DomainOperations)
		doc.Metrics = append(doc.Metrics, scaffold.OperationsMetrics()...)
	case prism.DomainSecurity:
		doc = scaffold.NewDocument(prism.DomainSecurity)
		doc.Metrics = append(doc.Metrics, scaffold.SecurityMetrics()...)
	case prism.DomainOperations:
		doc = scaffold.NewDocument(prism.DomainOperations)
		doc.Metrics = append(doc.Metrics, scaffold.OperationsMetrics()...)
	default:
		return fmt.Errorf("invalid domain %q, must be 'security' or 'operations'", initDomain)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// Write to file (0644 is appropriate for shareable config files)
	if err := os.WriteFile(initOutput, data, 0644); err != nil { //nolint:gosec // G306: PRISM docs are shareable configs
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Created %s\n", initOutput)
	return nil
}

// createOperationsMetrics moved to scaffold package as scaffold.OperationsMetrics()
