package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/grokify/prism/maturity"
	"github.com/spf13/cobra"
)

var maturityCmd = &cobra.Command{
	Use:   "maturity",
	Short: "Maturity model commands",
	Long:  `Commands for working with maturity model specifications.`,
}

var maturityReportCmd = &cobra.Command{
	Use:   "report <maturity-spec-file>",
	Short: "Generate markdown report from a maturity specification",
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
  prism maturity report spec.json                       # Markdown to stdout
  prism maturity report spec.json -o report.md          # Markdown to file
  prism maturity report spec.json --view domain         # Domain view only
  prism maturity report spec.json --view framework      # Framework view only
  prism maturity report spec.json --frameworks NIST_CSF_2,NIST_800_53`,
	Args: cobra.ExactArgs(1),
	RunE: runMaturityReport,
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
	maturityCmd.AddCommand(maturityReportCmd)

	maturityReportCmd.Flags().StringVarP(&maturityOutput, "output", "o", "", "Output file (default: stdout)")
	maturityReportCmd.Flags().StringVarP(&maturityFormat, "format", "f", "markdown", "Output format: markdown, json")
	maturityReportCmd.Flags().StringVarP(&maturityView, "view", "v", "both", "View type: both, domain, framework")
	maturityReportCmd.Flags().StringVar(&maturityTitle, "title", "", "Report title (default: from metadata or 'Maturity Model')")
	maturityReportCmd.Flags().StringVar(&maturityAuthor, "author", "", "Report author")
	maturityReportCmd.Flags().BoolVar(&maturityNoMeta, "no-meta", false, "Omit YAML front matter")
	maturityReportCmd.Flags().BoolVar(&maturityNoTOC, "no-toc", false, "Omit table of contents")
	maturityReportCmd.Flags().BoolVar(&maturityNoDetail, "no-detail", false, "Omit criterion details (framework mappings)")
	maturityReportCmd.Flags().StringVar(&maturityFrameworks, "frameworks", "", "Filter to specific frameworks (comma-separated)")

	rootCmd.AddCommand(maturityCmd)
}

func runMaturityReport(cmd *cobra.Command, args []string) error {
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
