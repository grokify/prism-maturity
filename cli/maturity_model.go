package cli

import (
	"encoding/json"
	"fmt"
	"os"

	capstack "github.com/grokify/prism-capability"
	"github.com/grokify/prism-maturity"
	"github.com/grokify/prism-maturity/dashboard"
	"github.com/grokify/prism-maturity/maturity"
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
	modelStateFile   string // --state flag for state document
	modelOutputFile  string // -o flag
	modelFormat      string // -f flag for dashboard
	modelCapStack    string // --capstack flag for capability stack
	modelAggregation string // --aggregation flag for aggregation method
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

When a capability stack is provided (--capstack), additional views are generated:
  - Layer-based maturity overview with aggregate maturity per layer
  - Layer summary cards showing capability counts and maturity
  - Capability bullets grouped by layer

Aggregation methods (--aggregation):
  min    Use the minimum value across SLIs/capabilities (conservative)
  avg    Use the average value across SLIs/capabilities

Output formats:
  --format json    JSON data (Dashforge format)
  --format html    Standalone HTML dashboard (default)

Examples:
  prism maturity model dashboard model.json
  prism maturity model dashboard model.json --state state.json
  prism maturity model dashboard model.json --state state.json -o dashboard.html
  prism maturity model dashboard model.json -f json -o dashboard.json
  prism maturity model dashboard model.json --capstack capstack.json --aggregation min`,
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

var maturityModelLintCmd = &cobra.Command{
	Use:   "lint <model-file>",
	Short: "Check maturity model for common issues",
	Long: `Lint a maturity model for issues that affect dashboard display and completeness.

Unlike 'validate' which checks structural correctness, 'lint' checks for:
  - Criteria without sliId (orphan criteria won't show in dashboard)
  - Criteria missing operator or target (threshold won't display)
  - SLIs without any criteria (unused SLIs)
  - SLIs missing unit (affects threshold formatting)
  - SLIs missing sliType (affects methodology grouping)
  - Incomplete threshold coverage across maturity levels

Exit codes:
  0 - No issues found
  1 - Warnings found (non-blocking)
  2 - Errors found (blocking issues)

Examples:
  prism maturity model lint model.json
  prism maturity model lint model.json --strict`,
	Args: cobra.ExactArgs(1),
	RunE: runMaturityModelLint,
}

var lintStrict bool // --strict flag

func init() {
	// Add model subcommand to maturity
	maturityCmd.AddCommand(maturityModelCmd)

	// Add subcommands to model
	maturityModelCmd.AddCommand(maturityModelDashboardCmd)
	maturityModelCmd.AddCommand(maturityModelValidateCmd)
	maturityModelCmd.AddCommand(maturityModelLintCmd)

	// Dashboard flags
	maturityModelDashboardCmd.Flags().StringVar(&modelStateFile, "state", "", "State document to read current values from")
	maturityModelDashboardCmd.Flags().StringVarP(&modelOutputFile, "output", "o", "", "Output file (default: stdout)")
	maturityModelDashboardCmd.Flags().StringVarP(&modelFormat, "format", "f", "html", "Output format: html, json")
	maturityModelDashboardCmd.Flags().StringVar(&modelCapStack, "capstack", "", "Capability stack document for layer-based views")
	maturityModelDashboardCmd.Flags().StringVar(&modelAggregation, "aggregation", "min", "Aggregation method: min, avg")

	// Lint flags
	maturityModelLintCmd.Flags().BoolVar(&lintStrict, "strict", false, "Treat warnings as errors")
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

	// Optionally load capability stack for layer-based views
	if modelCapStack != "" {
		cs, err := capstack.LoadFromFile(modelCapStack)
		if err != nil {
			return fmt.Errorf("failed to read capstack: %w", err)
		}
		gen.WithCapabilityStack(cs)

		// Set aggregation method
		switch modelAggregation {
		case "min":
			gen.WithAggregationMethod(dashboard.AggregationMin)
		case "avg":
			gen.WithAggregationMethod(dashboard.AggregationAvg)
		default:
			return fmt.Errorf("invalid aggregation method: %s (must be: min, avg)", modelAggregation)
		}
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

// LintIssue represents a linting issue found in the model.
type LintIssue struct {
	Severity string // "error", "warning", "info"
	Location string // e.g., "domain.level.criterion"
	Message  string
	Hint     string // Suggested fix
}

func runMaturityModelLint(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Read and parse maturity spec
	spec, err := maturity.ReadSpecFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read model: %w", err)
	}

	var issues []LintIssue

	// Track which SLIs are referenced by criteria
	referencedSLIs := make(map[string]bool)

	// Check each domain
	for domainKey, domain := range spec.Domains {
		if domain == nil {
			continue
		}

		// Track criteria per SLI per level for threshold coverage analysis
		sliLevelCoverage := make(map[string]map[int]bool) // sliId -> level -> hasThreshold

		for _, level := range domain.Levels {
			for _, criterion := range level.Criteria {
				loc := fmt.Sprintf("%s.M%d.%s", domainKey, level.Level, criterion.ID)

				// Check for missing sliId
				if criterion.SLIID == "" {
					issues = append(issues, LintIssue{
						Severity: "warning",
						Location: loc,
						Message:  "Criterion has no sliId - won't appear in dashboard bullet charts",
						Hint:     "Add sliId to link this criterion to an SLI definition",
					})
					continue
				}

				referencedSLIs[criterion.SLIID] = true

				// Initialize level coverage map for this SLI
				if sliLevelCoverage[criterion.SLIID] == nil {
					sliLevelCoverage[criterion.SLIID] = make(map[int]bool)
				}
				sliLevelCoverage[criterion.SLIID][level.Level] = true

				// Check for missing operator
				if criterion.Operator == "" {
					issues = append(issues, LintIssue{
						Severity: "error",
						Location: loc,
						Message:  "Criterion missing operator - threshold cannot be displayed",
						Hint:     "Add operator (gte, lte, eq, exists) to define the threshold type",
					})
				}

				// Check for missing target on non-qualitative criteria
				if criterion.Operator != "exists" && criterion.Target == 0 {
					// Get SLI to check if it's qualitative
					isQualitative := false
					if sli, ok := spec.SLIs[criterion.SLIID]; ok && sli != nil {
						isQualitative = sli.IsQualitativeOnly()
					}
					if !isQualitative {
						issues = append(issues, LintIssue{
							Severity: "warning",
							Location: loc,
							Message:  "Criterion has target=0 which may be unintentional for quantitative metrics",
							Hint:     "Set target to the threshold value, or use operator='exists' for qualitative criteria",
						})
					}
				}

				// Check for qualitative criteria not using 'exists' operator
				if criterion.Type == "qualitative" && criterion.Operator != "exists" {
					issues = append(issues, LintIssue{
						Severity: "warning",
						Location: loc,
						Message:  "Qualitative criterion should use operator='exists'",
						Hint:     "Change operator to 'exists' for qualitative criteria",
					})
				}
			}
		}

		// Check for incomplete threshold coverage (SLIs that skip levels)
		for sliID, levels := range sliLevelCoverage {
			// Find min and max levels with criteria
			minLevel, maxLevel := 6, 0
			for lvl := range levels {
				if lvl < minLevel {
					minLevel = lvl
				}
				if lvl > maxLevel {
					maxLevel = lvl
				}
			}

			// Check for gaps
			for lvl := minLevel; lvl <= maxLevel; lvl++ {
				if !levels[lvl] {
					issues = append(issues, LintIssue{
						Severity: "info",
						Location: fmt.Sprintf("%s.%s", domainKey, sliID),
						Message:  fmt.Sprintf("SLI has no criterion at M%d (gap between M%d and M%d)", lvl, minLevel, maxLevel),
						Hint:     fmt.Sprintf("Add a criterion for M%d or this level will show '-' in dashboard", lvl),
					})
				}
			}
		}
	}

	// Check SLIs
	for sliID, sli := range spec.SLIs {
		if sli == nil {
			continue
		}

		// Check for unreferenced SLIs
		if !referencedSLIs[sliID] {
			issues = append(issues, LintIssue{
				Severity: "warning",
				Location: fmt.Sprintf("slis.%s", sliID),
				Message:  "SLI is defined but not referenced by any criteria",
				Hint:     "Add criteria that reference this SLI, or remove the unused SLI definition",
			})
		}

		// Check for missing unit
		if sli.Unit == "" && !sli.IsQualitativeOnly() {
			issues = append(issues, LintIssue{
				Severity: "warning",
				Location: fmt.Sprintf("slis.%s", sliID),
				Message:  "SLI missing unit - thresholds will display without units (e.g., '50' instead of '50%')",
				Hint:     "Add unit field (e.g., '%', 'ms', 'hours', 'count')",
			})
		}

		// Check for missing sliType (info level - not required but helpful)
		if sli.SLIType == "" {
			issues = append(issues, LintIssue{
				Severity: "info",
				Location: fmt.Sprintf("slis.%s", sliID),
				Message:  "SLI missing sliType - cannot be grouped by methodology (RED/USE/Golden Signals)",
				Hint:     "Add sliType (availability, latency, error_rate, throughput, saturation, utilization, quality)",
			})
		}
	}

	// Output results
	errorCount := 0
	warningCount := 0
	infoCount := 0

	for _, issue := range issues {
		switch issue.Severity {
		case "error":
			errorCount++
			fmt.Printf("ERROR   %s\n", issue.Location)
		case "warning":
			warningCount++
			fmt.Printf("WARNING %s\n", issue.Location)
		case "info":
			infoCount++
			fmt.Printf("INFO    %s\n", issue.Location)
		}
		fmt.Printf("        %s\n", issue.Message)
		fmt.Printf("        → %s\n\n", issue.Hint)
	}

	// Summary
	fmt.Printf("─────────────────────────────────────────\n")
	fmt.Printf("Lint summary for %s:\n", filename)
	fmt.Printf("  Errors:   %d\n", errorCount)
	fmt.Printf("  Warnings: %d\n", warningCount)
	fmt.Printf("  Info:     %d\n", infoCount)

	// Exit code based on findings
	if errorCount > 0 {
		return fmt.Errorf("found %d errors", errorCount)
	}
	if lintStrict && warningCount > 0 {
		return fmt.Errorf("found %d warnings (strict mode)", warningCount)
	}
	if len(issues) == 0 {
		fmt.Println("\n✓ No issues found")
	}

	return nil
}
