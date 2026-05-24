package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-maturity"
	"github.com/grokify/prism-maturity/report"
	"github.com/spf13/cobra"
)

// maturityPlanCmd is the parent command for plan subcommands.
var maturityPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Work with maturity plan documents",
	Long: `Commands for working with prism-maturity-plan documents that define
how to achieve target maturity levels.

Plan documents contain:
  - Goals with maturity targets and timelines
  - Phases (quarters) with goal targets
  - Initiatives linked to goals and phases
  - Metrics with SLO definitions`,
}

// Plan subcommand flags.
var (
	planOutputFile string // -o flag
	planFormat     string // -f flag
	planTitle      string // --title flag
	planAuthor     string // --author flag
	planMaxGaps    int    // --max-gaps flag
	planView       string // --view flag
	planNoMeta     bool   // --no-meta flag
	planNoDetail   bool   // --no-detail flag
)

var maturityPlanDashboardCmd = &cobra.Command{
	Use:   "dashboard <plan-file>",
	Short: "Generate executive dashboard from a maturity plan",
	Long: `Generate an executive-level dashboard showing maturity progress,
SLO compliance, and roadmap status.

Output formats:
  --format json      JSON data (default)
  --format markdown  Pandoc-compatible markdown
  --format marp      Marp presentation slides
  --format html      Standalone HTML dashboard

Examples:
  prism maturity plan dashboard plan.json                    # JSON to stdout
  prism maturity plan dashboard plan.json -f markdown        # Pandoc markdown
  prism maturity plan dashboard plan.json -f html -o dash.html # HTML file
  prism maturity plan dashboard plan.json -f marp -o slides.md # Marp slides`,
	Args: cobra.ExactArgs(1),
	RunE: runMaturityPlanDashboard,
}

var maturityPlanReportCmd = &cobra.Command{
	Use:   "report <plan-file>",
	Short: "Generate roadmap report from a maturity plan",
	Long: `Generate a roadmap report in Markdown or JSON format.

The report can be generated in two views:
  - by-phase: Phase → Goal → Initiative (timeline view)
  - by-goal:  Goal → Phase → Initiative (strategic view)
  - both:     Both views in a single document (default)

Examples:
  prism maturity plan report plan.json                   # Markdown to stdout
  prism maturity plan report plan.json -o report.md      # Markdown to file
  prism maturity plan report plan.json --format json     # JSON output
  prism maturity plan report plan.json --view by-phase   # Phase-centric only
  prism maturity plan report plan.json --view by-goal    # Goal-centric only`,
	Args: cobra.ExactArgs(1),
	RunE: runMaturityPlanReport,
}

func init() {
	// Add plan subcommand to maturity
	maturityCmd.AddCommand(maturityPlanCmd)

	// Add subcommands to plan
	maturityPlanCmd.AddCommand(maturityPlanDashboardCmd)
	maturityPlanCmd.AddCommand(maturityPlanReportCmd)

	// Dashboard flags
	maturityPlanDashboardCmd.Flags().StringVarP(&planFormat, "format", "f", "json", "Output format: json, markdown, marp, html")
	maturityPlanDashboardCmd.Flags().StringVarP(&planOutputFile, "output", "o", "", "Output file (default: stdout)")
	maturityPlanDashboardCmd.Flags().StringVar(&planTitle, "title", "", "Dashboard title (default: from metadata)")
	maturityPlanDashboardCmd.Flags().StringVar(&planAuthor, "author", "", "Dashboard author")
	maturityPlanDashboardCmd.Flags().IntVar(&planMaxGaps, "max-gaps", 10, "Maximum gaps to show (0 = all)")

	// Report flags
	maturityPlanReportCmd.Flags().StringVarP(&planOutputFile, "output", "o", "", "Output file (default: stdout)")
	maturityPlanReportCmd.Flags().StringVarP(&planFormat, "format", "f", "markdown", "Output format: markdown, json")
	maturityPlanReportCmd.Flags().StringVarP(&planView, "view", "v", "both", "View type: both, by-phase, by-goal")
	maturityPlanReportCmd.Flags().StringVar(&planTitle, "title", "", "Report title (default: from metadata or 'PRISM Roadmap Report')")
	maturityPlanReportCmd.Flags().StringVar(&planAuthor, "author", "", "Report author (default: from metadata)")
	maturityPlanReportCmd.Flags().BoolVar(&planNoMeta, "no-meta", false, "Omit YAML front matter (Markdown only)")
	maturityPlanReportCmd.Flags().BoolVar(&planNoDetail, "no-detail", false, "Omit initiative details")
}

func runMaturityPlanDashboard(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Read and parse document
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var doc prism.PRISMDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Generate dashboard data
	dashboard := doc.GenerateExecutiveDashboard()

	// Apply title override
	if planTitle != "" {
		dashboard.Title = planTitle
	}

	var output string

	switch planFormat {
	case "json":
		jsonData, err := json.MarshalIndent(dashboard, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		output = string(jsonData)

	case "markdown", "md":
		opts := report.DefaultDashboardOptions()
		if planTitle != "" {
			opts.Title = planTitle
		} else if dashboard.Title != "" {
			opts.Title = dashboard.Title
		}
		opts.Author = planAuthor
		opts.MaxGaps = planMaxGaps
		output = report.GenerateDashboardMarkdown(dashboard, opts)

	case "marp":
		opts := report.DefaultDashboardOptions()
		if planTitle != "" {
			opts.Title = planTitle
		} else if dashboard.Title != "" {
			opts.Title = dashboard.Title
		}
		opts.Author = planAuthor
		opts.MaxGaps = planMaxGaps
		output = report.GenerateDashboardMarp(dashboard, opts)

	case "html":
		opts := report.DefaultDashboardOptions()
		if planTitle != "" {
			opts.Title = planTitle
		} else if dashboard.Title != "" {
			opts.Title = dashboard.Title
		}
		opts.Author = planAuthor
		opts.MaxGaps = planMaxGaps
		output = report.GenerateDashboardHTML(dashboard, opts)

	default:
		return fmt.Errorf("unknown format: %s (must be: json, markdown, marp, html)", planFormat)
	}

	// Write output
	if planOutputFile != "" {
		if err := os.WriteFile(planOutputFile, []byte(output), 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Dashboard written to %s\n", planOutputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}

func runMaturityPlanReport(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Read and parse document
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var doc prism.PRISMDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate view type
	switch planView {
	case "both", "by-phase", "by-goal":
		// Valid
	default:
		return fmt.Errorf("invalid view type: %s (must be: both, by-phase, by-goal)", planView)
	}

	var output string

	switch planFormat {
	case "markdown", "md":
		opts := report.DefaultMarkdownOptions()
		opts.ViewType = planView
		opts.IncludeYAMLMeta = !planNoMeta
		opts.IncludeDetails = !planNoDetail

		if planTitle != "" {
			opts.Title = planTitle
		} else if doc.Metadata != nil && doc.Metadata.Name != "" {
			opts.Title = doc.Metadata.Name + " Roadmap"
		}

		if planAuthor != "" {
			opts.Author = planAuthor
		}

		output = report.GenerateMarkdown(&doc, opts)

	case "json":
		roadmapReport := doc.GenerateRoadmapReport()
		jsonData, err := json.MarshalIndent(roadmapReport, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		output = string(jsonData)

	default:
		return fmt.Errorf("invalid format: %s (must be: markdown, json)", planFormat)
	}

	// Write output
	if planOutputFile != "" {
		if err := os.WriteFile(planOutputFile, []byte(output), 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Report written to %s\n", planOutputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}
