package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-maturity"
	"github.com/spf13/cobra"
)

var sloReportCmd = &cobra.Command{
	Use:   "slo-report <file>",
	Short: "Generate SLO report showing requirements by maturity level",
	Long: `Generate a report of all SLO requirements organized by category and maturity level.
This shows how SLOs become more stringent as maturity increases.

Output formats:
  --format json      JSON output (default)
  --format markdown  Pandoc-compatible markdown
  --format marp      Marp presentation markdown
  --format matrix    Matrix view showing all levels`,
	Args: cobra.ExactArgs(1),
	RunE: runSLOReport,
}

var sloReportFormat string

func init() {
	sloReportCmd.Flags().StringVarP(&sloReportFormat, "format", "f", "json", "Output format: json, markdown, marp, matrix")
}

func runSLOReport(cmd *cobra.Command, args []string) error {
	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var doc prism.PRISMDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	report := doc.GenerateSLOReport()

	switch sloReportFormat {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	case "markdown":
		fmt.Print(report.ToMarkdown())
	case "marp":
		fmt.Print(report.ToMarp())
	case "matrix":
		fmt.Print(report.ToMatrixMarkdown())
	default:
		return fmt.Errorf("unknown format: %s", sloReportFormat)
	}

	return nil
}
