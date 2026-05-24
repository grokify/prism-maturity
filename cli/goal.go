package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-maturity"
	"github.com/spf13/cobra"
)

var goalCmd = &cobra.Command{
	Use:   "goal",
	Short: "Work with PRISM goals",
	Long:  `Commands for listing, showing, and tracking progress of goals in a PRISM document.`,
}

var goalListCmd = &cobra.Command{
	Use:   "list <prism-file>",
	Short: "List all goals in a PRISM document",
	Long: `Display all goals from a PRISM document with their current and target maturity levels.

Example:
  prism goal list prism.json`,
	Args: cobra.ExactArgs(1),
	RunE: runGoalList,
}

var goalShowCmd = &cobra.Command{
	Use:   "show <prism-file> <goal-id>",
	Short: "Show details of a specific goal",
	Long: `Display detailed information about a specific goal including maturity model and SLO requirements.

Example:
  prism goal show prism.json goal-reliability`,
	Args: cobra.ExactArgs(2),
	RunE: runGoalShow,
}

var goalProgressCmd = &cobra.Command{
	Use:   "progress <prism-file> [goal-id]",
	Short: "Show goal progress across phases",
	Long: `Display progress for a goal across all phases, or for all goals if no goal ID specified.

Example:
  prism goal progress prism.json
  prism goal progress prism.json goal-reliability`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runGoalProgress,
}

var goalOutputFormat string

func init() {
	goalCmd.AddCommand(goalListCmd)
	goalCmd.AddCommand(goalShowCmd)
	goalCmd.AddCommand(goalProgressCmd)

	goalListCmd.Flags().StringVarP(&goalOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
	goalShowCmd.Flags().StringVarP(&goalOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
	goalProgressCmd.Flags().StringVarP(&goalOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
}

func runGoalList(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	if len(doc.Goals) == 0 {
		fmt.Println("No goals defined in document.")
		return nil
	}

	formatter := NewFormatter(goalOutputFormat)

	// JSON outputs the full goal objects
	if goalOutputFormat == "json" {
		return formatter.WriteJSON(doc.Goals)
	}

	// Build table data
	rows := make([][]string, 0, len(doc.Goals))
	for _, goal := range doc.Goals {
		currentLevel := goal.CurrentLevel
		if currentLevel == 0 {
			currentLevel = goal.CurrentMaturityLevel(doc)
		}

		status := getGoalStatus(currentLevel, goal.TargetLevel)
		rows = append(rows, []string{
			truncateString(goal.Name, 25),
			truncateString(goal.Owner, 12),
			fmt.Sprintf("M%d", currentLevel),
			fmt.Sprintf("M%d", goal.TargetLevel),
			status,
		})
	}

	return formatter.WriteTable(&TableData{
		Title:   "Goals",
		Headers: []string{"NAME", "OWNER", "CURRENT", "TARGET", "STATUS"},
		Rows:    rows,
		Summary: fmt.Sprintf("Total: %d goals", len(doc.Goals)),
	})
}

func runGoalShow(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	goalID := args[1]
	goal := doc.GetGoalByID(goalID)
	if goal == nil {
		return fmt.Errorf("goal not found: %s", goalID)
	}

	if goalOutputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(goal)
	}

	currentLevel := goal.CurrentLevel
	if currentLevel == 0 {
		currentLevel = goal.CurrentMaturityLevel(doc)
	}

	fmt.Printf("Goal: %s\n", goal.Name)
	fmt.Printf("ID: %s\n", goal.ID)
	if goal.Description != "" {
		fmt.Printf("Description: %s\n", goal.Description)
	}
	if goal.Owner != "" {
		fmt.Printf("Owner: %s\n", goal.Owner)
	}
	if goal.Status != "" {
		fmt.Printf("Status: %s\n", goal.Status)
	}
	fmt.Printf("\nMaturity: M%d (current) -> M%d (target)\n", currentLevel, goal.TargetLevel)

	// Show maturity model levels
	if goal.MaturityModel != nil && len(goal.MaturityModel.Levels) > 0 {
		fmt.Println("\nMaturity Levels:")
		for _, level := range goal.MaturityModel.Levels {
			marker := " "
			if level.Level == currentLevel {
				marker = ">"
			}
			fmt.Printf("%s M%d: %s\n", marker, level.Level, level.Name)

			// Show required SLOs
			if len(level.RequiredSLOs) > 0 {
				for _, slo := range level.RequiredSLOs {
					metric := doc.GetMetricByID(slo.MetricID)
					status := "?"
					if metric != nil {
						if metric.MeetsSLO() {
							status = "Met"
						} else {
							status = "Not Met"
						}
					}
					fmt.Printf("     - SLO: %s [%s]\n", slo.MetricID, status)
				}
			}

			// Show metric criteria
			if len(level.MetricCriteria) > 0 {
				for _, mc := range level.MetricCriteria {
					metric := doc.GetMetricByID(mc.MetricID)
					status := "?"
					if metric != nil {
						if mc.IsMet(metric.Current) {
							status = "Met"
						} else {
							status = "Not Met"
						}
					}
					fmt.Printf("     - Criterion: %s %s %.1f [%s]\n",
						mc.MetricID, operatorSymbol(mc.Operator), mc.Value, status)
				}
			}
		}
	}

	// Show SLO compliance summary
	met, total := goal.SLOsMetForLevel(goal.TargetLevel, doc)
	if total > 0 {
		fmt.Printf("\nSLO Compliance for Target (M%d): %d/%d (%.0f%%)\n",
			goal.TargetLevel, met, total, float64(met)/float64(total)*100)
	}

	// Show associated initiatives
	initiatives := doc.GetInitiativesForGoal(goalID)
	if len(initiatives) > 0 {
		fmt.Printf("\nInitiatives (%d):\n", len(initiatives))
		for _, init := range initiatives {
			status := init.Status
			if status == "" {
				status = "not_started"
			}
			fmt.Printf("  - %s [%s]\n", init.Name, status)
		}
	}

	return nil
}

func runGoalProgress(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	// If specific goal ID provided, show just that goal
	if len(args) == 2 {
		goalID := args[1] //nolint:gosec // Safe: len(args)==2 checked above
		view := doc.GenerateGoalRoadmapView(goalID)
		if view == nil {
			return fmt.Errorf("goal not found: %s", goalID)
		}

		if goalOutputFormat == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(view)
		}

		return printGoalProgress(view)
	}

	// Show all goals
	if len(doc.Goals) == 0 {
		fmt.Println("No goals defined in document.")
		return nil
	}

	var views []prism.GoalRoadmapView
	for _, goal := range doc.Goals {
		view := doc.GenerateGoalRoadmapView(goal.ID)
		if view != nil {
			views = append(views, *view)
		}
	}

	if goalOutputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(views)
	}

	fmt.Println("Goal Progress Across Phases")
	fmt.Println("===========================")

	for i, view := range views {
		if i > 0 {
			fmt.Println()
		}
		if err := printGoalProgress(&view); err != nil {
			return err
		}
	}

	return nil
}

func printGoalProgress(view *prism.GoalRoadmapView) error {
	fmt.Printf("\n%s (M%d -> M%d)\n", view.GoalName, view.CurrentLevel, view.TargetLevel)
	if view.Description != "" {
		fmt.Printf("  %s\n", view.Description)
	}

	if len(view.PhaseProgress) == 0 {
		fmt.Println("  No phases defined for this goal.")
		return nil
	}

	fmt.Println()
	fmt.Printf("  %-15s %-8s %-12s %-20s\n", "PHASE", "LEVELS", "INITIATIVES", "COMPLETION")
	fmt.Printf("  %-15s %-8s %-12s %-20s\n", "-----", "------", "-----------", "----------")

	for _, pp := range view.PhaseProgress {
		phaseName := pp.PhaseName
		if pp.Quarter != "" {
			phaseName = fmt.Sprintf("%s %d", pp.Quarter, pp.Year)
		}

		levels := fmt.Sprintf("M%d->M%d", pp.EnterLevel, pp.ExitLevel)
		initiatives := fmt.Sprintf("%d/%d", pp.InitiativesCompleted, pp.InitiativesTotal)
		completion := fmt.Sprintf("%.0f%%", pp.CompletionPercent)

		fmt.Printf("  %-15s %-8s %-12s %-20s\n", phaseName, levels, initiatives, completion)

		// Show individual initiatives
		for _, init := range pp.Initiatives {
			status := init.Status
			if status == "" {
				status = "pending"
			}
			fmt.Printf("    - %s [%s]\n", init.Name, status)
		}
	}

	return nil
}

// Helper functions are now in format.go (re-exported from output package)
