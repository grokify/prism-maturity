package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	capstack "github.com/grokify/prism-capability"
	caprender "github.com/grokify/prism-capability/render"
	prism "github.com/grokify/prism-maturity"
	"github.com/grokify/prism-maturity/dashboard"
	"github.com/grokify/prism-maturity/maturity"
	"github.com/grokify/prism-maturity/render"
	"github.com/spf13/cobra"
)

var (
	capstackOutput       string
	capstackFormat       string
	capstackTitle        string
	capstackStyle        string
	capstackState        string
	capstackModel        string
	capstackAggregation  string
	capstackNoDeps       bool
	capstackNoFoundation bool
	capstackNoLegend     bool
	capstackColorBy      string
	capstackDarkTheme    bool
	capstackStandalone   bool
)

var capstackCmd = &cobra.Command{
	Use:   "capstack",
	Short: "Capability stack operations with maturity overlay",
	Long:  `Commands for rendering capability stacks with maturity level overlays.`,
}

var capstackRenderCmd = &cobra.Command{
	Use:   "render <capstack-file>",
	Short: "Render capability stack with maturity overlay",
	Long: `Renders a capability stack diagram with maturity level badges.

This command combines capability stack visualization with maturity data,
showing the maturity level (M1-M5) for each capability as a badge.

Required inputs:
  - Capability stack file (JSON)
  - Maturity model file (--model)
  - State file (--state) with SLI measurements

Supported formats:
  d2    D2 diagram language (https://d2lang.com)
  html  Static HTML (embeddable or standalone)

Examples:
  # Render with maturity overlay
  prism maturity capstack render stack.json --model model.json --state state.json -o stack.d2

  # Render to SVG via d2
  prism maturity capstack render stack.json --model model.json --state state.json | d2 - stack.svg

  # HTML with dark theme
  prism maturity capstack render stack.json --model model.json --state state.json -f html --dark -o stack.html

  # Grid style for executives
  prism maturity capstack render stack.json --model model.json --state state.json --style=grid -o exec.d2`,
	Args: cobra.ExactArgs(1),
	RunE: runCapstackRender,
}

func init() {
	capstackRenderCmd.Flags().StringVarP(&capstackOutput, "output", "o", "", "Output file (default: stdout)")
	capstackRenderCmd.Flags().StringVarP(&capstackFormat, "format", "f", "d2", "Output format (d2, html)")
	capstackRenderCmd.Flags().StringVarP(&capstackTitle, "title", "t", "", "Diagram title (default: from metadata)")
	capstackRenderCmd.Flags().StringVarP(&capstackStyle, "style", "s", "default", "Render style: default or grid (D2 only)")
	capstackRenderCmd.Flags().StringVar(&capstackModel, "model", "", "Maturity model file (required)")
	capstackRenderCmd.Flags().StringVar(&capstackState, "state", "", "State file with SLI measurements (required)")
	capstackRenderCmd.Flags().StringVar(&capstackAggregation, "aggregation", "min", "Aggregation method: min or avg")
	capstackRenderCmd.Flags().BoolVar(&capstackNoDeps, "no-deps", false, "Hide dependency arrows (D2 only)")
	capstackRenderCmd.Flags().BoolVar(&capstackNoFoundation, "no-foundational", false, "Hide foundational capabilities")
	capstackRenderCmd.Flags().BoolVar(&capstackNoLegend, "no-legend", false, "Hide status legend")
	capstackRenderCmd.Flags().StringVar(&capstackColorBy, "color-by", "status", "Color scheme: status or category (D2 only)")
	capstackRenderCmd.Flags().BoolVar(&capstackDarkTheme, "dark", false, "Use dark theme (HTML only)")
	capstackRenderCmd.Flags().BoolVar(&capstackStandalone, "standalone", false, "Generate complete HTML document (HTML only)")

	if err := capstackRenderCmd.MarkFlagRequired("model"); err != nil {
		panic(err)
	}
	if err := capstackRenderCmd.MarkFlagRequired("state"); err != nil {
		panic(err)
	}

	capstackCmd.AddCommand(capstackRenderCmd)
}

func runCapstackRender(cmd *cobra.Command, args []string) error {
	capstackPath := args[0]

	// Load capability stack
	cs, err := capstack.LoadFromFile(capstackPath)
	if err != nil {
		return fmt.Errorf("failed to load capability stack %s: %w", capstackPath, err)
	}

	// Validate capability stack
	if errs := cs.Validate(); errs.HasErrors() {
		return fmt.Errorf("capability stack validation failed: %s", errs.Error())
	}

	// Load maturity model
	spec, err := maturity.ReadSpecFile(capstackModel)
	if err != nil {
		return fmt.Errorf("failed to load maturity model %s: %w", capstackModel, err)
	}

	// Load state document
	stateData, err := os.ReadFile(capstackState)
	if err != nil {
		return fmt.Errorf("failed to read state file %s: %w", capstackState, err)
	}
	var stateDoc prism.PRISMDocument
	if err := json.Unmarshal(stateData, &stateDoc); err != nil {
		return fmt.Errorf("failed to parse state file %s: %w", capstackState, err)
	}

	// Determine aggregation method
	var aggMethod dashboard.AggregationMethod
	switch strings.ToLower(capstackAggregation) {
	case "avg", "average":
		aggMethod = dashboard.AggregationAvg
	default:
		aggMethod = dashboard.AggregationMin
	}

	// Create maturity aggregator
	agg := dashboard.NewMaturityAggregator(spec, cs, &stateDoc, aggMethod)

	// Build maturity overlay
	overlays := render.BuildMaturityOverlay(agg)

	// Configure D2 options
	var d2Opts caprender.D2Options
	switch strings.ToLower(capstackStyle) {
	case "grid", "exec", "executive":
		d2Opts = caprender.GridD2Options()
	default:
		d2Opts = caprender.DefaultD2Options()
	}

	d2Opts.Title = capstackTitle
	d2Opts.ShowDependencies = !capstackNoDeps && d2Opts.Style != caprender.D2StyleGrid
	d2Opts.ShowFoundational = !capstackNoFoundation
	d2Opts.ShowLegend = !capstackNoLegend
	d2Opts.ColorByStatus = strings.ToLower(capstackColorBy) == "status"
	d2Opts.Overlays = overlays

	// Determine output
	var out *os.File
	if capstackOutput == "" {
		out = os.Stdout
	} else {
		// Auto-detect format from extension if not specified
		if capstackFormat == "d2" {
			ext := strings.ToLower(filepath.Ext(capstackOutput))
			switch ext {
			case ".d2":
				capstackFormat = "d2"
			case ".html", ".htm":
				capstackFormat = "html"
			}
		}

		var err error
		out, err = os.Create(capstackOutput)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer out.Close()
	}

	// Render
	switch strings.ToLower(capstackFormat) {
	case "d2":
		if err := caprender.RenderD2(out, cs, d2Opts); err != nil {
			return fmt.Errorf("render failed: %w", err)
		}
	case "html", "htm":
		htmlOpts := caprender.HTMLOptions{
			Title:            capstackTitle,
			ShowLegend:       !capstackNoLegend,
			ShowFoundational: !capstackNoFoundation,
			Standalone:       capstackStandalone,
			DarkTheme:        capstackDarkTheme,
			Overlays:         overlays,
		}
		if err := caprender.RenderHTML(out, cs, htmlOpts); err != nil {
			return fmt.Errorf("render failed: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s (supported: d2, html)", capstackFormat)
	}

	if capstackOutput != "" {
		fmt.Fprintf(os.Stderr, "Rendered capability stack with maturity overlay to %s\n", capstackOutput)
	}

	return nil
}
