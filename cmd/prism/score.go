package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-intelligence"
	"github.com/spf13/cobra"
)

var (
	scoreOutputFormat string
	scoreDetailed     bool
	scoreGoals        bool
	scoreLegacy       bool
)

var scoreCmd = &cobra.Command{
	Use:   "score <file>",
	Short: "Calculate PRISM score for a document",
	Long: `Calculate the composite PRISM score for a document.

The score combines maturity levels, metric performance, and optionally
customer awareness data into a single health score (0.0-1.0).

By default, empty cells (no metrics and no maturity data) are skipped
to avoid penalizing sparse data. Use --legacy to include all cells.

Examples:
  prism score prism.json
  prism score prism.json -f json
  prism score prism.json --detailed
  prism score prism.json --goals     # Use goal maturity instead of global cells
  prism score prism.json --legacy    # Include empty cells (original behavior)`,
	Args: cobra.ExactArgs(1),
	RunE: runScore,
}

func init() {
	scoreCmd.Flags().StringVarP(&scoreOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
	scoreCmd.Flags().BoolVar(&scoreDetailed, "detailed", false, "Show detailed breakdown")
	scoreCmd.Flags().BoolVar(&scoreGoals, "goals", false, "Use goal-level maturity instead of global maturity model")
	scoreCmd.Flags().BoolVar(&scoreLegacy, "legacy", false, "Include empty cells in scoring (original behavior)")
}

func runScore(cmd *cobra.Command, args []string) error {
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

	// Validate first
	errs := doc.Validate()
	if errs.HasErrors() {
		return fmt.Errorf("document validation failed: %v", errs)
	}

	// Calculate score
	var config *prism.ScoreConfig
	if scoreLegacy {
		config = prism.LegacyScoreConfig()
	} else {
		config = prism.DefaultScoreConfig()
	}
	if scoreGoals {
		config.UseGoalMaturity = true
	}
	score := doc.CalculatePRISMScore(config, nil)

	if scoreOutputFormat == "json" {
		return outputScoreJSON(score)
	}

	return outputScoreText(score, scoreDetailed)
}

func outputScoreJSON(score *prism.PRISMScore) error {
	data, err := json.MarshalIndent(score, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal score: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func outputScoreText(score *prism.PRISMScore, detailed bool) error {
	health := score.GetHealthStatus()

	fmt.Println("PRISM Score Report")
	fmt.Println("==================")
	fmt.Println()

	// Overall score with color indicator
	colorIndicator := ""
	switch health.Color {
	case prism.StatusGreen:
		colorIndicator = "🟢"
	case prism.StatusYellow:
		colorIndicator = "🟡"
	case prism.StatusRed:
		colorIndicator = "🔴"
	}

	fmt.Printf("Overall Score: %.1f%% %s %s\n", score.Overall*100, colorIndicator, score.Interpretation)
	fmt.Println()

	fmt.Println("Component Scores:")
	fmt.Printf("  Base Score:      %.1f%%\n", score.BaseScore*100)
	fmt.Printf("  Awareness Score: %.1f%%\n", score.AwarenessScore*100)
	fmt.Println()

	fmt.Println("Domain Scores:")
	fmt.Printf("  Security:   %.1f%%\n", score.SecurityScore*100)
	fmt.Printf("  Operations: %.1f%%\n", score.OperationsScore*100)
	fmt.Println()

	if detailed {
		fmt.Println("Cell Breakdown:")
		fmt.Println("  Domain      | Stage    | Maturity | Performance | Cell Score | Weight")
		fmt.Println("  ------------|----------|----------|-------------|------------|-------")
		for _, cs := range score.CellScores {
			fmt.Printf("  %-11s | %-8s | %6.1f%% | %9.1f%% | %8.1f%% | %.2f\n",
				cs.Domain,
				cs.Stage,
				cs.MaturityScore*100,
				cs.PerformanceScore*100,
				cs.CellScore*100,
				cs.Weight,
			)
		}
		fmt.Println()

		breakdown := score.GetScoreBreakdown()
		fmt.Println("Stage Summary:")
		for stage, sb := range breakdown.StageBreakdown {
			fmt.Printf("  %s: %.1f%% (weight: %.2f)\n", stage, sb.Score*100, sb.Weight)
		}
	}

	fmt.Println()
	fmt.Printf("Health Status: %s\n", health.Description)

	return nil
}
