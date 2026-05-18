package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grokify/prism-intelligence"
	"github.com/grokify/prism-intelligence/analysis"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export PRISM data to files",
	Long:  `Export PRISM data in multiple formats to an output directory.`,
}

var exportGoalsCmd = &cobra.Command{
	Use:   "goals <prism-file>",
	Short: "Export goals to files",
	Long: `Export goals data in one or more formats to the output directory.

Examples:
  prism export goals prism.json -o ./out -f json
  prism export goals prism.json -o ./out -f json -f markdown
  prism export goals prism.json -o ./out -f json,markdown,toon`,
	Args: cobra.ExactArgs(1),
	RunE: runExportGoals,
}

var exportPhasesCmd = &cobra.Command{
	Use:   "phases <prism-file>",
	Short: "Export phases to files",
	Long: `Export phases data in one or more formats to the output directory.

Examples:
  prism export phases prism.json -o ./out -f json
  prism export phases prism.json -o ./out -f json,markdown`,
	Args: cobra.ExactArgs(1),
	RunE: runExportPhases,
}

var exportRoadmapCmd = &cobra.Command{
	Use:   "roadmap <prism-file>",
	Short: "Export roadmap to files",
	Long: `Export roadmap data in one or more formats to the output directory.

Examples:
  prism export roadmap prism.json -o ./out -f json
  prism export roadmap prism.json -o ./out -f json,markdown-pandoc,markdown-marp`,
	Args: cobra.ExactArgs(1),
	RunE: runExportRoadmap,
}

var (
	exportOutputDir string
	exportFormats   []string
)

func init() {
	exportCmd.AddCommand(exportGoalsCmd)
	exportCmd.AddCommand(exportPhasesCmd)
	exportCmd.AddCommand(exportRoadmapCmd)

	// Add flags to all export subcommands
	for _, cmd := range []*cobra.Command{exportGoalsCmd, exportPhasesCmd, exportRoadmapCmd} {
		cmd.Flags().StringVarP(&exportOutputDir, "output-dir", "o", "", "Output directory (required)")
		cmd.Flags().StringSliceVarP(&exportFormats, "format", "f", []string{"json"}, "Output format(s): json, markdown, markdown-pandoc, markdown-marp, toon, yaml, csv")
		_ = cmd.MarkFlagRequired("output-dir")
	}
}

// parseFormats parses the format flags, handling comma-separated values
func parseFormats(formats []string) []string {
	result := []string{}
	for _, f := range formats {
		// Handle comma-separated values
		parts := strings.Split(f, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				result = append(result, p)
			}
		}
	}
	return result
}

// getFileSuffix returns the file suffix (including format identifier) for a format
func getFileSuffix(format string) string {
	switch format {
	case "json":
		return ".json"
	case "markdown":
		return ".md"
	case "markdown-pandoc":
		return ".pandoc.md"
	case "markdown-marp":
		return ".marp.md"
	case "toon":
		return ".toon"
	case "yaml":
		return ".yaml"
	case "csv":
		return ".csv"
	default:
		return ".txt"
	}
}

// ensureOutputDir creates the output directory if it doesn't exist
func ensureOutputDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func runExportGoals(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	if err := ensureOutputDir(exportOutputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	formats := parseFormats(exportFormats)

	for _, format := range formats {
		filename := filepath.Join(exportOutputDir, "goals"+getFileSuffix(format))

		if err := exportGoalsToFile(doc, filename, format); err != nil {
			return fmt.Errorf("failed to export %s: %w", format, err)
		}

		fmt.Printf("Exported: %s\n", filename)
	}

	return nil
}

func exportGoalsToFile(doc *prism.PRISMDocument, filename, format string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	switch format {
	case "json":
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return enc.Encode(doc.Goals)

	case "markdown", "markdown-pandoc":
		return exportGoalsMarkdown(doc, f, false)

	case "markdown-marp":
		return exportGoalsMarkdown(doc, f, true)

	case "toon":
		return exportGoalsTOON(doc, f)

	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

//nolint:unparam // error return for future extensibility
func exportGoalsMarkdown(doc *prism.PRISMDocument, f *os.File, marp bool) error {
	if marp {
		fmt.Fprintln(f, "---")
		fmt.Fprintln(f, "marp: true")
		fmt.Fprintln(f, "theme: default")
		fmt.Fprintln(f, "---")
		fmt.Fprintln(f)
		fmt.Fprintln(f, "# Goals Overview")
		fmt.Fprintln(f)
		fmt.Fprintln(f, "---")
		fmt.Fprintln(f)
	} else {
		fmt.Fprintln(f, "# Goals Overview")
		fmt.Fprintln(f)
	}

	// Summary table
	fmt.Fprintln(f, "## Summary")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "| Goal | Owner | Current | Target | Status |")
	fmt.Fprintln(f, "| --- | --- | --- | --- | --- |")

	for _, goal := range doc.Goals {
		currentLevel := goal.CurrentLevel
		if currentLevel == 0 {
			currentLevel = goal.CurrentMaturityLevel(doc)
		}
		status := getGoalStatus(currentLevel, goal.TargetLevel)
		fmt.Fprintf(f, "| %s | %s | M%d | M%d | %s |\n",
			goal.Name, goal.Owner, currentLevel, goal.TargetLevel, status)
	}

	// Detail sections
	for _, goal := range doc.Goals {
		if marp {
			fmt.Fprintln(f)
			fmt.Fprintln(f, "---")
		}
		fmt.Fprintln(f)
		fmt.Fprintf(f, "## %s\n\n", goal.Name)

		if goal.Description != "" {
			fmt.Fprintf(f, "%s\n\n", goal.Description)
		}

		currentLevel := goal.CurrentLevel
		if currentLevel == 0 {
			currentLevel = goal.CurrentMaturityLevel(doc)
		}

		fmt.Fprintf(f, "- **Owner:** %s\n", goal.Owner)
		fmt.Fprintf(f, "- **Current Level:** M%d\n", currentLevel)
		fmt.Fprintf(f, "- **Target Level:** M%d\n", goal.TargetLevel)

		if goal.MaturityModel != nil && len(goal.MaturityModel.Levels) > 0 {
			fmt.Fprintln(f)
			fmt.Fprintln(f, "### Maturity Levels")
			fmt.Fprintln(f)
			for _, level := range goal.MaturityModel.Levels {
				marker := ""
				if level.Level == currentLevel {
					marker = " **(current)**"
				}
				fmt.Fprintf(f, "- **M%d:** %s%s\n", level.Level, level.Name, marker)
			}
		}
	}

	return nil
}

//nolint:unparam // error return for future extensibility
func exportGoalsTOON(doc *prism.PRISMDocument, f *os.File) error {
	// Header
	fmt.Fprintln(f, "GOALS;NAME,OWNER,CURRENT,TARGET,STATUS")

	for _, goal := range doc.Goals {
		currentLevel := goal.CurrentLevel
		if currentLevel == 0 {
			currentLevel = goal.CurrentMaturityLevel(doc)
		}
		status := getGoalStatus(currentLevel, goal.TargetLevel)
		fmt.Fprintf(f, ";%s,%s,M%d,M%d,%s\n",
			goal.Name, goal.Owner, currentLevel, goal.TargetLevel, status)
	}

	return nil
}

func runExportPhases(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	if err := ensureOutputDir(exportOutputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	formats := parseFormats(exportFormats)
	phases := doc.GetPhasesSorted()

	for _, format := range formats {
		filename := filepath.Join(exportOutputDir, "phases"+getFileSuffix(format))

		if err := exportPhasesToFile(phases, doc, filename, format); err != nil {
			return fmt.Errorf("failed to export %s: %w", format, err)
		}

		fmt.Printf("Exported: %s\n", filename)
	}

	return nil
}

func exportPhasesToFile(phases []prism.Phase, doc *prism.PRISMDocument, filename, format string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	switch format {
	case "json":
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return enc.Encode(phases)

	case "markdown", "markdown-pandoc":
		return exportPhasesMarkdown(phases, doc, f, false)

	case "markdown-marp":
		return exportPhasesMarkdown(phases, doc, f, true)

	case "toon":
		return exportPhasesTOON(phases, f)

	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

//nolint:unparam // error return for future extensibility
func exportPhasesMarkdown(phases []prism.Phase, doc *prism.PRISMDocument, f *os.File, marp bool) error {
	if marp {
		fmt.Fprintln(f, "---")
		fmt.Fprintln(f, "marp: true")
		fmt.Fprintln(f, "theme: default")
		fmt.Fprintln(f, "---")
		fmt.Fprintln(f)
		fmt.Fprintln(f, "# Phases Overview")
		fmt.Fprintln(f)
		fmt.Fprintln(f, "---")
		fmt.Fprintln(f)
	} else {
		fmt.Fprintln(f, "# Phases Overview")
		fmt.Fprintln(f)
	}

	// Summary table
	fmt.Fprintln(f, "## Summary")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "| Phase | Period | Status | Goals |")
	fmt.Fprintln(f, "| --- | --- | --- | --- |")

	for _, phase := range phases {
		period := ""
		if phase.Quarter != "" {
			period = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
		}
		status := phase.Status
		if status == "" {
			status = "planned"
		}
		fmt.Fprintf(f, "| %s | %s | %s | %d |\n",
			phase.Name, period, status, len(phase.GoalTargets))
	}

	// Detail sections
	for _, phase := range phases {
		if marp {
			fmt.Fprintln(f)
			fmt.Fprintln(f, "---")
		}
		fmt.Fprintln(f)
		fmt.Fprintf(f, "## %s\n\n", phase.Name)

		if phase.Quarter != "" {
			fmt.Fprintf(f, "**Period:** %s %d\n\n", phase.Quarter, phase.Year)
		}
		if phase.StartDate != "" {
			fmt.Fprintf(f, "**Dates:** %s to %s\n\n", phase.StartDate, phase.EndDate)
		}

		if len(phase.GoalTargets) > 0 {
			fmt.Fprintln(f, "### Goal Targets")
			fmt.Fprintln(f)
			fmt.Fprintln(f, "| Goal | Enter | Exit |")
			fmt.Fprintln(f, "| --- | --- | --- |")
			for _, gt := range phase.GoalTargets {
				goal := doc.GetGoalByID(gt.GoalID)
				name := gt.GoalID
				if goal != nil {
					name = goal.Name
				}
				fmt.Fprintf(f, "| %s | M%d | M%d |\n", name, gt.EnterLevel, gt.ExitLevel)
			}
		}
	}

	return nil
}

//nolint:unparam // error return for future extensibility
func exportPhasesTOON(phases []prism.Phase, f *os.File) error {
	fmt.Fprintln(f, "PHASES;NAME,PERIOD,STATUS,GOALS")

	for _, phase := range phases {
		period := ""
		if phase.Quarter != "" {
			period = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
		}
		status := phase.Status
		if status == "" {
			status = "planned"
		}
		fmt.Fprintf(f, ";%s,%s,%s,%d\n",
			phase.Name, period, status, len(phase.GoalTargets))
	}

	return nil
}

func runExportRoadmap(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	if err := ensureOutputDir(exportOutputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	formats := parseFormats(exportFormats)
	progress := analysis.CalculateRoadmapProgress(doc)

	for _, format := range formats {
		filename := filepath.Join(exportOutputDir, "roadmap"+getFileSuffix(format))

		if err := exportRoadmapToFile(progress, doc, filename, format); err != nil {
			return fmt.Errorf("failed to export %s: %w", format, err)
		}

		fmt.Printf("Exported: %s\n", filename)
	}

	return nil
}

func exportRoadmapToFile(progress *analysis.RoadmapProgress, doc *prism.PRISMDocument, filename, format string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	switch format {
	case "json":
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return enc.Encode(progress)

	case "markdown", "markdown-pandoc":
		return exportRoadmapMarkdown(progress, doc, f, false)

	case "markdown-marp":
		return exportRoadmapMarkdown(progress, doc, f, true)

	case "toon":
		return exportRoadmapTOON(progress, f)

	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

//nolint:unparam // error return for future extensibility
func exportRoadmapMarkdown(progress *analysis.RoadmapProgress, _ *prism.PRISMDocument, f *os.File, marp bool) error {
	if marp {
		fmt.Fprintln(f, "---")
		fmt.Fprintln(f, "marp: true")
		fmt.Fprintln(f, "theme: default")
		fmt.Fprintln(f, "---")
		fmt.Fprintln(f)
	}

	fmt.Fprintln(f, "# Roadmap Progress")
	fmt.Fprintln(f)
	fmt.Fprintf(f, "**Overall Completion:** %.0f%%\n\n", progress.OverallCompletion)

	if marp {
		fmt.Fprintln(f, "---")
		fmt.Fprintln(f)
	}

	// Phase summary
	fmt.Fprintln(f, "## Phase Summary")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "| Phase | Period | Status | Goals | Completion |")
	fmt.Fprintln(f, "| --- | --- | --- | --- | --- |")

	for _, ps := range progress.PhaseProgress {
		fmt.Fprintf(f, "| %s | %s | %s | %s | %.0f%% |\n",
			ps.PhaseName, ps.Period, ps.Status, ps.GoalSummary, ps.Completion)
	}

	if marp {
		fmt.Fprintln(f)
		fmt.Fprintln(f, "---")
	}
	fmt.Fprintln(f)

	// Goal summary
	fmt.Fprintln(f, "## Goal Summary")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "| Goal | Current | Target | Status |")
	fmt.Fprintln(f, "| --- | --- | --- | --- |")

	for _, gs := range progress.GoalProgress {
		fmt.Fprintf(f, "| %s | M%d | M%d | %s |\n",
			gs.GoalName, gs.CurrentLevel, gs.TargetLevel, gs.Status)
	}

	return nil
}

//nolint:unparam // error return for future extensibility
func exportRoadmapTOON(progress *analysis.RoadmapProgress, f *os.File) error {
	fmt.Fprintf(f, "ROADMAP;completion=%.0f%%\n", progress.OverallCompletion)

	fmt.Fprintln(f, "PHASES;NAME,PERIOD,STATUS,GOALS,COMPLETION")
	for _, ps := range progress.PhaseProgress {
		fmt.Fprintf(f, ";%s,%s,%s,%s,%.0f%%\n",
			ps.PhaseName, ps.Period, ps.Status, ps.GoalSummary, ps.Completion)
	}

	fmt.Fprintln(f, "GOALS;NAME,CURRENT,TARGET,STATUS")
	for _, gs := range progress.GoalProgress {
		fmt.Fprintf(f, ";%s,M%d,M%d,%s\n",
			gs.GoalName, gs.CurrentLevel, gs.TargetLevel, gs.Status)
	}

	return nil
}
