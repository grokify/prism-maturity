package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-maturity"
	"github.com/grokify/prism-maturity/maturity"
	"github.com/spf13/cobra"
)

// maturityStateCmd is the parent command for state subcommands.
var maturityStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Work with maturity state documents",
	Long: `Commands for working with prism-maturity-state documents that track
current state and measurements against a maturity model.

State documents contain:
  - SLI state with temporal tracking (past, present, future targets)
  - Maturity level state per domain
  - Enabler/initiative progress tracking`,
}

// State subcommand flags.
var (
	stateModelFile  string // --model flag
	stateOutputFile string // -o flag
	stateFormat     string // -f flag
)

var maturityStateValidateCmd = &cobra.Command{
	Use:   "validate <state-file>",
	Short: "Validate a maturity state document",
	Long: `Validate a maturity state document for structural correctness.

When a model file is provided (--model), cross-validation is performed:
  - SLI IDs in sliState exist in model SLIs
  - Domain keys in maturityState exist in model domains
  - Enabler IDs in enablerState exist in model enablers

Examples:
  prism maturity state validate state.json
  prism maturity state validate state.json --model model.json`,
	Args: cobra.ExactArgs(1),
	RunE: runMaturityStateValidate,
}

var maturityStateShowCmd = &cobra.Command{
	Use:   "show <state-file>",
	Short: "Display maturity state summary",
	Long: `Display a summary of the maturity state document.

Shows current maturity levels, SLI states, and enabler progress.
When a model file is provided (--model), displays names from the model.

Output formats:
  --format text    Human-readable text output (default)
  --format json    JSON output

Examples:
  prism maturity state show state.json
  prism maturity state show state.json --model model.json
  prism maturity state show state.json -f json`,
	Args: cobra.ExactArgs(1),
	RunE: runMaturityStateShow,
}

func init() {
	// Add state subcommand to maturity
	maturityCmd.AddCommand(maturityStateCmd)

	// Add subcommands to state
	maturityStateCmd.AddCommand(maturityStateValidateCmd)
	maturityStateCmd.AddCommand(maturityStateShowCmd)

	// Validate flags
	maturityStateValidateCmd.Flags().StringVar(&stateModelFile, "model", "", "Model document for cross-validation")

	// Show flags
	maturityStateShowCmd.Flags().StringVar(&stateModelFile, "model", "", "Model document for resolving names")
	maturityStateShowCmd.Flags().StringVarP(&stateOutputFile, "output", "o", "", "Output file (default: stdout)")
	maturityStateShowCmd.Flags().StringVarP(&stateFormat, "format", "f", "text", "Output format: text, json")
}

func runMaturityStateValidate(cmd *cobra.Command, args []string) error {
	stateFile := args[0]

	// Load state document
	stateData, err := os.ReadFile(stateFile)
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}

	var stateDoc prism.MaturityStateDocument
	if err := json.Unmarshal(stateData, &stateDoc); err != nil {
		return fmt.Errorf("failed to parse state: %w", err)
	}

	var errors []string

	// Basic validation
	if stateDoc.Metadata.Name == "" {
		errors = append(errors, "metadata.name is required")
	}

	// Validate SLI state entries
	for sliID, state := range stateDoc.SLIState {
		if state == nil {
			errors = append(errors, fmt.Sprintf("sliState[%q] is nil", sliID))
			continue
		}
		if state.SLIID == "" {
			errors = append(errors, fmt.Sprintf("sliState[%q].sliId is empty", sliID))
		}
	}

	// Validate maturity state entries
	for domainID, state := range stateDoc.MaturityState {
		if state == nil {
			errors = append(errors, fmt.Sprintf("maturityState[%q] is nil", domainID))
			continue
		}
		if state.DomainID == "" {
			errors = append(errors, fmt.Sprintf("maturityState[%q].domainId is empty", domainID))
		}
		if state.Current != nil {
			if state.Current.Level < 1 || state.Current.Level > 5 {
				errors = append(errors, fmt.Sprintf("maturityState[%q].current.level must be 1-5, got %d", domainID, state.Current.Level))
			}
		}
		if state.Target != nil {
			if state.Target.Level < 1 || state.Target.Level > 5 {
				errors = append(errors, fmt.Sprintf("maturityState[%q].target.level must be 1-5, got %d", domainID, state.Target.Level))
			}
		}
	}

	// Validate enabler state entries
	for enablerID, state := range stateDoc.EnablerState {
		if state == nil {
			errors = append(errors, fmt.Sprintf("enablerState[%q] is nil", enablerID))
			continue
		}
		if state.EnablerID == "" {
			errors = append(errors, fmt.Sprintf("enablerState[%q].enablerId is empty", enablerID))
		}
		validStatuses := map[string]bool{
			"not_started": true, "in_progress": true, "completed": true, "blocked": true,
		}
		if state.Status != "" && !validStatuses[state.Status] {
			errors = append(errors, fmt.Sprintf("enablerState[%q].status %q is invalid", enablerID, state.Status))
		}
		if state.Progress < 0 || state.Progress > 100 {
			errors = append(errors, fmt.Sprintf("enablerState[%q].progress must be 0-100, got %.1f", enablerID, state.Progress))
		}
	}

	// Cross-validation with model if provided
	if stateModelFile != "" {
		spec, err := maturity.ReadSpecFile(stateModelFile)
		if err != nil {
			return fmt.Errorf("failed to read model: %w", err)
		}

		// Validate SLI references
		for sliID := range stateDoc.SLIState {
			if spec.SLIs != nil {
				if _, ok := spec.SLIs[sliID]; !ok {
					errors = append(errors, fmt.Sprintf("sliState references unknown SLI %q", sliID))
				}
			}
		}

		// Validate domain references
		for domainID := range stateDoc.MaturityState {
			if _, ok := spec.Domains[domainID]; !ok {
				errors = append(errors, fmt.Sprintf("maturityState references unknown domain %q", domainID))
			}
		}

		// Validate enabler references
		enablerIDs := make(map[string]bool)
		for _, domain := range spec.Domains {
			for _, level := range domain.Levels {
				for _, enabler := range level.Enablers {
					enablerIDs[enabler.ID] = true
				}
			}
		}
		for enablerID := range stateDoc.EnablerState {
			if !enablerIDs[enablerID] {
				errors = append(errors, fmt.Sprintf("enablerState references unknown enabler %q", enablerID))
			}
		}
	}

	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
		return fmt.Errorf("state document has %d validation errors", len(errors))
	}

	// Print summary
	fmt.Printf("✓ %s is valid\n", stateFile)
	fmt.Printf("  SLI states: %d\n", len(stateDoc.SLIState))
	fmt.Printf("  Domain states: %d\n", len(stateDoc.MaturityState))
	fmt.Printf("  Enabler states: %d\n", len(stateDoc.EnablerState))

	return nil
}

func runMaturityStateShow(cmd *cobra.Command, args []string) error {
	stateFile := args[0]

	// Load state document
	stateData, err := os.ReadFile(stateFile)
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}

	var stateDoc prism.MaturityStateDocument
	if err := json.Unmarshal(stateData, &stateDoc); err != nil {
		return fmt.Errorf("failed to parse state: %w", err)
	}

	// Optionally load model for name resolution
	var spec *maturity.Spec
	if stateModelFile != "" {
		spec, err = maturity.ReadSpecFile(stateModelFile)
		if err != nil {
			return fmt.Errorf("failed to read model: %w", err)
		}
	}

	var output string

	switch stateFormat {
	case "json":
		jsonData, err := json.MarshalIndent(stateDoc, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		output = string(jsonData)

	case "text":
		output = formatStateText(&stateDoc, spec)

	default:
		return fmt.Errorf("unknown format: %s (must be: text, json)", stateFormat)
	}

	// Write output
	if stateOutputFile != "" {
		if err := os.WriteFile(stateOutputFile, []byte(output), 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("State written to %s\n", stateOutputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}

func formatStateText(state *prism.MaturityStateDocument, spec *maturity.Spec) string {
	var result string

	// Header
	result += fmt.Sprintf("=== %s ===\n", state.Metadata.Name)
	if state.Metadata.Description != "" {
		result += fmt.Sprintf("%s\n", state.Metadata.Description)
	}
	if state.Metadata.AssessedAt != "" {
		result += fmt.Sprintf("Assessed: %s", state.Metadata.AssessedAt)
		if state.Metadata.AssessedBy != "" {
			result += fmt.Sprintf(" by %s", state.Metadata.AssessedBy)
		}
		result += "\n"
	}
	result += "\n"

	// Maturity State
	if len(state.MaturityState) > 0 {
		result += "MATURITY LEVELS\n"
		result += "---------------\n"
		for domainID, domainState := range state.MaturityState {
			domainName := domainID
			if spec != nil {
				if domain, ok := spec.Domains[domainID]; ok && domain != nil {
					domainName = domain.Name
				}
			}

			current := 0
			target := 0
			if domainState.Current != nil {
				current = domainState.Current.Level
			}
			if domainState.Target != nil {
				target = domainState.Target.Level
			}

			result += fmt.Sprintf("  %s: M%d", domainName, current)
			if target > 0 {
				result += fmt.Sprintf(" → M%d", target)
			}
			result += "\n"
		}
		result += "\n"
	}

	// SLI State
	if len(state.SLIState) > 0 {
		result += "SLI STATE\n"
		result += "---------\n"
		for sliID, sliState := range state.SLIState {
			sliName := sliID
			if spec != nil && spec.SLIs != nil {
				if sli, ok := spec.SLIs[sliID]; ok && sli != nil {
					sliName = sli.Name
				}
			}

			result += fmt.Sprintf("  %s: ", sliName)
			if sliState.QualitativeState != "" {
				result += sliState.QualitativeState
			}
			if len(sliState.Windows) > 0 {
				result += " ["
				first := true
				for window, ws := range sliState.Windows {
					if !first {
						result += ", "
					}
					result += fmt.Sprintf("%s: %.1f", window, ws.Value)
					first = false
				}
				result += "]"
			}
			result += "\n"
		}
		result += "\n"
	}

	// Enabler State
	if len(state.EnablerState) > 0 {
		// Count by status
		counts := map[string]int{
			"completed":   0,
			"in_progress": 0,
			"not_started": 0,
			"blocked":     0,
		}
		for _, es := range state.EnablerState {
			counts[es.Status]++
		}

		result += "ENABLER PROGRESS\n"
		result += "----------------\n"
		result += fmt.Sprintf("  Completed: %d\n", counts["completed"])
		result += fmt.Sprintf("  In Progress: %d\n", counts["in_progress"])
		result += fmt.Sprintf("  Not Started: %d\n", counts["not_started"])
		if counts["blocked"] > 0 {
			result += fmt.Sprintf("  Blocked: %d\n", counts["blocked"])
		}
		result += "\n"
	}

	return result
}
