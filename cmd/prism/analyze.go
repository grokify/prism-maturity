package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/grokify/prism-intelligence"
	"github.com/grokify/prism-intelligence/analysis"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze <prism-file>",
	Short: "Analyze PRISM document and generate initiative recommendations",
	Long: `Analyze a PRISM document to generate initiative recommendations for achieving
maturity targets. Outputs a structured analysis that can be used directly or
fed to an LLM for detailed planning.

Examples:
  # Generate analysis prompt for LLM
  prism analyze prism.json

  # Output as JSON for programmatic use
  prism analyze prism.json -f json

  # Generate detailed recommendations (requires LLM integration)
  prism analyze prism.json --recommend`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

var (
	analyzeOutputFormat string
	analyzeRecommend    bool
)

func init() {
	analyzeCmd.Flags().StringVarP(&analyzeOutputFormat, "format", "f", "text", "Output format (text|json|prompt)")
	analyzeCmd.Flags().BoolVar(&analyzeRecommend, "recommend", false, "Generate detailed recommendations (placeholder for LLM integration)")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	result := analysis.Analyze(doc)

	switch analyzeOutputFormat {
	case "json":
		return outputAnalysisJSON(result)
	case "prompt":
		return outputAnalysisPrompt(doc, result)
	default:
		return outputAnalysisText(result)
	}
}

func outputAnalysisJSON(result *analysis.Result) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

func outputAnalysisText(result *analysis.Result) error {
	fmt.Println("PRISM Analysis")
	fmt.Println("==============")
	fmt.Println()

	// Summary
	fmt.Println("Summary")
	fmt.Println("-------")
	fmt.Printf("Goals: %d | Phases: %d | SLOs: %d/%d met (%.0f%%)\n",
		result.Summary.TotalGoals,
		result.Summary.TotalPhases,
		result.Summary.SLOsMet,
		result.Summary.TotalSLOs,
		result.Summary.SLOCompliance)
	fmt.Printf("Average Maturity Gap: %.1f levels\n", result.Summary.AvgMaturityGap)
	fmt.Println()

	// Goals
	fmt.Println("Goal Analysis")
	fmt.Println("-------------")
	for _, ga := range result.Goals {
		fmt.Printf("\n%s (M%d → M%d) [%s]\n", ga.GoalName, ga.CurrentLevel, ga.TargetLevel, ga.Status)
		if ga.Gap > 0 {
			fmt.Printf("  Gap: %d levels | SLOs: %d/%d met\n", ga.Gap, ga.SLOsMet, ga.SLOsTotal)
		}
		if len(ga.SLOsRequired) > 0 {
			fmt.Println("  Required SLOs:")
			for _, slo := range ga.SLOsRequired {
				status := "Not Met"
				if slo.IsMet {
					status = "Met"
				}
				fmt.Printf("    - [M%d] %s: %s [%s]\n", slo.Level, slo.MetricName, slo.Target, status)
			}
		}
	}
	fmt.Println()

	// Phases
	fmt.Println("Phase Analysis")
	fmt.Println("--------------")
	for _, pa := range result.Phases {
		period := pa.Period
		if period == "" {
			period = pa.PhaseName
		}
		fmt.Printf("\n%s [%s] - %.0f%% complete\n", period, pa.Status, pa.Completion)
		fmt.Printf("  Initiatives: %d\n", pa.Initiatives)
		if len(pa.GoalTargets) > 0 {
			fmt.Println("  Goal Targets:")
			for _, gt := range pa.GoalTargets {
				fmt.Printf("    - %s: M%d → M%d (%d SLOs needed)\n", gt.GoalName, gt.EnterLevel, gt.ExitLevel, gt.SLOsNeeded)
			}
		}
	}
	fmt.Println()

	// Gaps
	if len(result.Gaps) > 0 {
		fmt.Println("Identified Gaps")
		fmt.Println("---------------")
		for _, gap := range result.Gaps {
			fmt.Printf("  [%s] %s: %s\n", strings.ToUpper(string(gap.Severity)), gap.Type, gap.Description)
		}
	}

	return nil
}

func outputAnalysisPrompt(_ *prism.PRISMDocument, result *analysis.Result) error {
	fmt.Println("# PRISM Analysis Prompt")
	fmt.Println()
	fmt.Println("You are an operational planning assistant. Analyze the following PRISM document")
	fmt.Println("and recommend initiatives to achieve the maturity targets.")
	fmt.Println()
	fmt.Println("## Current State")
	fmt.Println()

	// Goals summary
	fmt.Println("### Goals")
	fmt.Println()
	for _, ga := range result.Goals {
		fmt.Printf("- **%s**: Currently at M%d, targeting M%d (%d level gap)\n", ga.GoalName, ga.CurrentLevel, ga.TargetLevel, ga.Gap)
		if len(ga.SLOsRequired) > 0 {
			fmt.Println("  - Required SLOs:")
			for _, slo := range ga.SLOsRequired {
				status := "NOT MET"
				if slo.IsMet {
					status = "MET"
				}
				fmt.Printf("    - [M%d] %s: target %s, current %.2f (%s)\n", slo.Level, slo.MetricName, slo.Target, slo.Current, status)
			}
		}
	}
	fmt.Println()

	// Phases summary
	fmt.Println("### Phases")
	fmt.Println()
	for _, pa := range result.Phases {
		period := pa.Period
		if period == "" {
			period = pa.PhaseName
		}
		fmt.Printf("- **%s** [%s]\n", period, pa.Status)
		for _, gt := range pa.GoalTargets {
			fmt.Printf("  - %s: M%d → M%d\n", gt.GoalName, gt.EnterLevel, gt.ExitLevel)
		}
		fmt.Printf("  - Current initiatives: %d\n", pa.Initiatives)
	}
	fmt.Println()

	// Gaps
	if len(result.Gaps) > 0 {
		fmt.Println("### Identified Gaps")
		fmt.Println()
		for _, gap := range result.Gaps {
			fmt.Printf("- [%s] %s\n", strings.ToUpper(string(gap.Severity)), gap.Description)
		}
		fmt.Println()
	}

	// Request
	fmt.Println("## Request")
	fmt.Println()
	fmt.Println("Based on the above analysis, please recommend initiatives that will:")
	fmt.Println()
	fmt.Println("1. Enable achievement of SLOs required for each maturity level progression")
	fmt.Println("2. Be appropriately sequenced across phases (dependencies considered)")
	fmt.Println("3. Address the identified gaps")
	fmt.Println()
	fmt.Println("For each initiative, provide:")
	fmt.Println()
	fmt.Println("- **Title**: Clear, actionable name")
	fmt.Println("- **Description**: What will be delivered")
	fmt.Println("- **Phase**: Which phase to execute in")
	fmt.Println("- **Goals**: Which goals this supports")
	fmt.Println("- **SLOs Enabled**: Which SLOs this helps achieve")
	fmt.Println("- **Priority**: High/Medium/Low")
	fmt.Println("- **Dependencies**: Other initiatives this depends on")
	fmt.Println()
	fmt.Println("Output as JSON array of recommendations.")

	return nil
}
