package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grokify/prism-intelligence"
	"github.com/grokify/prism-intelligence/dashforge"
	"github.com/spf13/cobra"
)

var dashforgeCmd = &cobra.Command{
	Use:   "dashforge <file>",
	Short: "Generate Dashforge dashboards from a PRISM document",
	Long: `Generate a set of Dashforge dashboard JSON files from a PRISM document.

This creates a multi-page dashboard set for visualization:
  - executive.json   - Executive summary with key metrics
  - maturity.json    - Maturity scorecard with gauges
  - slo-matrix.json  - SLO requirements matrix
  - roadmap.json     - Roadmap progress timeline
  - gaps.json        - Gap analysis and priorities
  - goals/*.json     - Individual goal deep-dives

Examples:
  prism dashforge prism.json -o dashboards/
  prism dashforge prism.json --base-id appsec -o dashboards/
  prism dashforge prism.json --no-goals -o dashboards/`,
	Args: cobra.ExactArgs(1),
	RunE: runDashforge,
}

var (
	dashforgeOutput   string
	dashforgeBaseID   string
	dashforgeDataPath string
	dashforgeNoGoals  bool
)

func init() {
	dashforgeCmd.Flags().StringVarP(&dashforgeOutput, "output", "o", "./dashboards", "Output directory")
	dashforgeCmd.Flags().StringVar(&dashforgeBaseID, "base-id", "prism", "Base ID for dashboard IDs")
	dashforgeCmd.Flags().StringVar(&dashforgeDataPath, "data-path", "./data/prism.json", "Path to PRISM data file (for URL data source)")
	dashforgeCmd.Flags().BoolVar(&dashforgeNoGoals, "no-goals", false, "Skip generating individual goal dashboards")
	rootCmd.AddCommand(dashforgeCmd)
}

func runDashforge(cmd *cobra.Command, args []string) error {
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

	// Configure conversion
	opts := dashforge.DefaultConvertOptions()
	opts.BaseID = dashforgeBaseID
	opts.DataSourcePath = dashforgeDataPath
	opts.GenerateGoalDashboards = !dashforgeNoGoals

	// Convert to dashboards
	set, err := dashforge.Convert(&doc, opts)
	if err != nil {
		return fmt.Errorf("failed to convert: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(dashforgeOutput, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write dashboards
	dashboards := map[string]interface{}{
		"executive.json":  set.Executive,
		"maturity.json":   set.Maturity,
		"slo-matrix.json": set.SLOMatrix,
		"roadmap.json":    set.Roadmap,
		"gaps.json":       set.Gaps,
	}

	for name, dash := range dashboards {
		if err := writeDashboard(filepath.Join(dashforgeOutput, name), dash); err != nil {
			return err
		}
		fmt.Printf("  ✓ %s\n", name)
	}

	// Write goal dashboards
	if !dashforgeNoGoals && len(set.Goals) > 0 {
		goalsDir := filepath.Join(dashforgeOutput, "goals")
		if err := os.MkdirAll(goalsDir, 0755); err != nil {
			return fmt.Errorf("failed to create goals directory: %w", err)
		}

		for id, dash := range set.Goals {
			name := fmt.Sprintf("%s.json", id)
			if err := writeDashboard(filepath.Join(goalsDir, name), dash); err != nil {
				return err
			}
			fmt.Printf("  ✓ goals/%s\n", name)
		}
	}

	// Write index file with dashboard manifest
	index := map[string]interface{}{
		"dashboards": []map[string]string{
			{"id": set.Executive.ID, "title": set.Executive.Title, "file": "executive.json"},
			{"id": set.Maturity.ID, "title": set.Maturity.Title, "file": "maturity.json"},
			{"id": set.SLOMatrix.ID, "title": set.SLOMatrix.Title, "file": "slo-matrix.json"},
			{"id": set.Roadmap.ID, "title": set.Roadmap.Title, "file": "roadmap.json"},
			{"id": set.Gaps.ID, "title": set.Gaps.Title, "file": "gaps.json"},
		},
	}
	if err := writeDashboard(filepath.Join(dashforgeOutput, "index.json"), index); err != nil {
		return err
	}
	fmt.Printf("  ✓ index.json\n")

	fmt.Printf("\nDashboards written to %s/\n", dashforgeOutput)
	return nil
}

func writeDashboard(path string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", path, err)
	}
	if err := os.WriteFile(path, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}
