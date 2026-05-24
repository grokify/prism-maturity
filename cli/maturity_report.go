package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/grokify/prism-maturity/maturity"
	"github.com/spf13/cobra"
)

var maturityCmd = &cobra.Command{
	Use:   "maturity",
	Short: "Maturity model, state, and plan commands",
	Long: `Commands for working with PRISM maturity documents:

  model  - Maturity model documents that define "what good looks like"
  state  - State documents that track current measurements
  plan   - Plan documents that define goals, phases, and initiatives`,
}

var maturityModelReportCmd = &cobra.Command{
	Use:   "report <model-file>",
	Short: "Generate markdown report from a maturity model",
	Long: `Generate a markdown report from a maturity model specification.

The report can be generated in three views:
  - domain:    Maturity levels organized by domain
  - framework: Criteria organized by compliance framework
  - both:      Both views in a single document (default)

Output includes:
  - YAML front matter (for Pandoc/MkDocs)
  - Table of contents
  - Domain view: domains → levels → criteria (SLOs) → enablers
  - Framework view: frameworks → controls → criteria mappings

Examples:
  prism maturity model report model.json                       # Markdown to stdout
  prism maturity model report model.json -o report.md          # Markdown to file
  prism maturity model report model.json --view domain         # Domain view only
  prism maturity model report model.json --view framework      # Framework view only
  prism maturity model report model.json --frameworks NIST_CSF_2,NIST_800_53`,
	Args: cobra.ExactArgs(1),
	RunE: runMaturityModelReport,
}

var (
	maturityOutput     string
	maturityFormat     string
	maturityView       string
	maturityTitle      string
	maturityAuthor     string
	maturityNoMeta     bool
	maturityNoTOC      bool
	maturityNoDetail   bool
	maturityFrameworks string
)

func init() {
	maturityModelCmd.AddCommand(maturityModelReportCmd)

	maturityModelReportCmd.Flags().StringVarP(&maturityOutput, "output", "o", "", "Output file (default: stdout)")
	maturityModelReportCmd.Flags().StringVarP(&maturityFormat, "format", "f", "markdown", "Output format: markdown, json")
	maturityModelReportCmd.Flags().StringVarP(&maturityView, "view", "v", "both", "View type: both, domain, framework")
	maturityModelReportCmd.Flags().StringVar(&maturityTitle, "title", "", "Report title (default: from metadata or 'Maturity Model')")
	maturityModelReportCmd.Flags().StringVar(&maturityAuthor, "author", "", "Report author")
	maturityModelReportCmd.Flags().BoolVar(&maturityNoMeta, "no-meta", false, "Omit YAML front matter")
	maturityModelReportCmd.Flags().BoolVar(&maturityNoTOC, "no-toc", false, "Omit table of contents")
	maturityModelReportCmd.Flags().BoolVar(&maturityNoDetail, "no-detail", false, "Omit criterion details (framework mappings)")
	maturityModelReportCmd.Flags().StringVar(&maturityFrameworks, "frameworks", "", "Filter to specific frameworks (comma-separated)")
}

func runMaturityModelReport(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Read and parse maturity spec
	spec, err := maturity.ReadSpecFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read maturity spec: %w", err)
	}

	// Validate view type
	switch maturityView {
	case "both", "domain", "framework":
		// Valid
	default:
		return fmt.Errorf("invalid view type: %s (must be: both, domain, framework)", maturityView)
	}

	var output string

	switch maturityFormat {
	case "markdown", "md":
		opts := maturity.DefaultMarkdownOptions()
		opts.ViewType = maturityView
		opts.IncludeYAMLMeta = !maturityNoMeta
		opts.IncludeTOC = !maturityNoTOC
		opts.IncludeDetails = !maturityNoDetail

		if maturityTitle != "" {
			opts.Title = maturityTitle
		}

		if maturityAuthor != "" {
			opts.Author = maturityAuthor
		}

		if maturityFrameworks != "" {
			opts.Frameworks = parseFrameworksList(maturityFrameworks)
		}

		output = spec.GenerateMarkdown(opts)

	case "json":
		jsonData, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		output = string(jsonData)

	default:
		return fmt.Errorf("invalid format: %s (must be: markdown, json)", maturityFormat)
	}

	// Write output
	if maturityOutput != "" {
		if err := os.WriteFile(maturityOutput, []byte(output), 0600); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Report written to %s\n", maturityOutput)
	} else {
		fmt.Print(output)
	}

	return nil
}

// parseFrameworksList parses a comma-separated list of framework identifiers.
func parseFrameworksList(input string) []string {
	if input == "" {
		return nil
	}
	var frameworks []string
	for _, fw := range strings.Split(input, ",") {
		fw = strings.TrimSpace(fw)
		if fw != "" {
			frameworks = append(frameworks, fw)
		}
	}
	return frameworks
}
