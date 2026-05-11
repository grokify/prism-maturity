package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism"
	"github.com/grokify/prism/dashboard"
	"github.com/grokify/prism/maturity"
	"github.com/spf13/cobra"
)

// maturityModelCmd is the parent command for model subcommands.
var maturityModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Work with maturity model documents",
	Long: `Commands for working with prism-maturity-model documents that define
what good looks like at each maturity level (M1-M5).

Model documents contain:
  - SLIs (Service Level Indicators) with framework mappings
  - Domain definitions with maturity levels
  - Criteria (SLOs) that define level achievement
  - Enablers (tasks) to achieve criteria`,
}

// Model subcommand flags.
var (
	modelStateFile  string // --state flag for state document
	modelOutputFile string // -o flag
	modelFormat     string // -f flag for dashboard
)

var maturityModelDashboardCmd = &cobra.Command{
	Use:   "dashboard <model-file>",
	Short: "Generate HTML dashboard from a maturity model",
	Long: `Generate an interactive HTML dashboard from a maturity model specification.

The dashboard displays:
  - Domain summary cards with current/target maturity levels
  - Bullet charts showing progress per SLI by methodology (RED/USE/Golden Signals)
  - Progress charts showing SLI advancement toward targets
  - SLI tables grouped by category

When a state document is provided (--state), current values are read from
the state document instead of relying on inline values in the model.

Output formats:
  --format json    JSON data (Dashforge format)
  --format html    Standalone HTML dashboard (default)

Examples:
  prism maturity model dashboard model.json
  prism maturity model dashboard model.json --state state.json
  prism maturity model dashboard model.json --state state.json -o dashboard.html
  prism maturity model dashboard model.json -f json -o dashboard.json`,
	Args: cobra.ExactArgs(1),
	RunE: runMaturityModelDashboard,
}

var maturityModelValidateCmd = &cobra.Command{
	Use:   "validate <model-file>",
	Short: "Validate a maturity model specification",
	Long: `Validate a maturity model specification for structural correctness.

Checks include:
  - Valid JSON structure
  - Required fields present (domains, levels, criteria)
  - SLI references resolve correctly
  - Level numbers are sequential (1-5)
  - Criterion operators are valid

Examples:
  prism maturity model validate model.json`,
	Args: cobra.ExactArgs(1),
	RunE: runMaturityModelValidate,
}

func init() {
	// Add model subcommand to maturity
	maturityCmd.AddCommand(maturityModelCmd)

	// Add subcommands to model
	maturityModelCmd.AddCommand(maturityModelDashboardCmd)
	maturityModelCmd.AddCommand(maturityModelValidateCmd)

	// Dashboard flags
	maturityModelDashboardCmd.Flags().StringVar(&modelStateFile, "state", "", "State document to read current values from")
	maturityModelDashboardCmd.Flags().StringVarP(&modelOutputFile, "output", "o", "", "Output file (default: stdout)")
	maturityModelDashboardCmd.Flags().StringVarP(&modelFormat, "format", "f", "html", "Output format: html, json")
}

func runMaturityModelDashboard(cmd *cobra.Command, args []string) error {
	modelFile := args[0]

	// Read and parse maturity spec
	spec, err := maturity.ReadSpecFile(modelFile)
	if err != nil {
		return fmt.Errorf("failed to read model: %w", err)
	}

	// Create dashboard generator
	gen := dashboard.NewGenerator(spec)

	// Optionally load state document
	if modelStateFile != "" {
		stateData, err := os.ReadFile(modelStateFile)
		if err != nil {
			return fmt.Errorf("failed to read state: %w", err)
		}
		var stateDoc prism.PRISMDocument
		if err := json.Unmarshal(stateData, &stateDoc); err != nil {
			return fmt.Errorf("failed to parse state: %w", err)
		}
		gen.WithStateDocument(&stateDoc)
	}

	// Generate dashboard
	dash, err := gen.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate dashboard: %w", err)
	}

	var output string

	switch modelFormat {
	case "json":
		jsonData, err := dash.ToJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		output = string(jsonData)

	case "html":
		html, err := dash.ToHTML(dashboard.DefaultHTMLOptions())
		if err != nil {
			return fmt.Errorf("failed to generate HTML: %w", err)
		}
		output = html

	default:
		return fmt.Errorf("unknown format: %s (must be: html, json)", modelFormat)
	}

	// Write output
	if modelOutputFile != "" {
		if err := os.WriteFile(modelOutputFile, []byte(output), 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Dashboard written to %s\n", modelOutputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}

func runMaturityModelValidate(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Read and parse maturity spec
	spec, err := maturity.ReadSpecFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read model: %w", err)
	}

	// Validate structure
	var errors []string

	// Check domains exist
	if len(spec.Domains) == 0 {
		errors = append(errors, "no domains defined")
	}

	// Validate each domain
	for domainKey, domain := range spec.Domains {
		if domain == nil {
			errors = append(errors, fmt.Sprintf("domain %q is nil", domainKey))
			continue
		}

		if domain.Name == "" {
			errors = append(errors, fmt.Sprintf("domain %q has no name", domainKey))
		}

		if len(domain.Levels) == 0 {
			errors = append(errors, fmt.Sprintf("domain %q has no levels defined", domainKey))
			continue
		}

		// Check level sequence
		levelNums := make(map[int]bool)
		for _, level := range domain.Levels {
			if level.Level < 1 || level.Level > 5 {
				errors = append(errors, fmt.Sprintf("domain %q has invalid level %d (must be 1-5)", domainKey, level.Level))
			}
			if levelNums[level.Level] {
				errors = append(errors, fmt.Sprintf("domain %q has duplicate level %d", domainKey, level.Level))
			}
			levelNums[level.Level] = true

			// Check criteria
			for _, criterion := range level.Criteria {
				if criterion.ID == "" {
					errors = append(errors, fmt.Sprintf("domain %q level %d has criterion without ID", domainKey, level.Level))
				}

				// Validate SLI reference
				if criterion.SLIID != "" && spec.SLIs != nil {
					if _, ok := spec.SLIs[criterion.SLIID]; !ok {
						errors = append(errors, fmt.Sprintf("criterion %q references unknown SLI %q", criterion.ID, criterion.SLIID))
					}
				}

				// Validate operator
				if criterion.Operator != "" {
					validOps := map[string]bool{
						"gte": true, "lte": true, "gt": true, "lt": true, "eq": true, "exists": true,
					}
					if !validOps[criterion.Operator] {
						errors = append(errors, fmt.Sprintf("criterion %q has invalid operator %q", criterion.ID, criterion.Operator))
					}
				}
			}
		}
	}

	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
		return fmt.Errorf("model has %d validation errors", len(errors))
	}

	// Print summary
	fmt.Printf("✓ %s is valid\n", filename)
	fmt.Printf("  Domains: %d\n", len(spec.Domains))
	fmt.Printf("  SLIs: %d\n", len(spec.SLIs))

	for domainKey, domain := range spec.Domains {
		totalCriteria := 0
		totalEnablers := 0
		for _, level := range domain.Levels {
			totalCriteria += len(level.Criteria)
			totalEnablers += len(level.Enablers)
		}
		fmt.Printf("  %s: %d levels, %d criteria, %d enablers\n", domainKey, len(domain.Levels), totalCriteria, totalEnablers)
	}

	return nil
}
