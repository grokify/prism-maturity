package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-intelligence"
	"github.com/grokify/prism-intelligence/analysis"
	"github.com/spf13/cobra"
)

var roadmapCmd = &cobra.Command{
	Use:   "roadmap",
	Short: "Work with PRISM roadmaps",
	Long:  `Commands for viewing roadmap progress across phases and goals.`,
}

var roadmapShowCmd = &cobra.Command{
	Use:   "show <prism-file>",
	Short: "Show the roadmap overview",
	Long: `Display the roadmap showing all phases and their associated goals.

Example:
  prism roadmap show prism.json
  prism roadmap show prism.json --by-goal`,
	Args: cobra.ExactArgs(1),
	RunE: runRoadmapShow,
}

var roadmapProgressCmd = &cobra.Command{
	Use:   "progress <prism-file>",
	Short: "Show overall roadmap progress",
	Long: `Display progress metrics across all phases and goals.

Example:
  prism roadmap progress prism.json`,
	Args: cobra.ExactArgs(1),
	RunE: runRoadmapProgress,
}

var (
	roadmapOutputFormat string
	roadmapByGoal       bool
)

func init() {
	roadmapCmd.AddCommand(roadmapShowCmd)
	roadmapCmd.AddCommand(roadmapProgressCmd)

	roadmapShowCmd.Flags().StringVarP(&roadmapOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
	roadmapShowCmd.Flags().BoolVar(&roadmapByGoal, "by-goal", false, "Organize roadmap by goal instead of by phase")
	roadmapProgressCmd.Flags().StringVarP(&roadmapOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
}

func runRoadmapShow(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	if roadmapByGoal {
		return showRoadmapByGoal(doc)
	}
	return showRoadmapByPhase(doc)
}

func showRoadmapByPhase(doc *prism.PRISMDocument) error {
	phases := doc.GetPhasesSorted()

	if len(phases) == 0 {
		fmt.Println("No phases defined in document.")
		return nil
	}

	if roadmapOutputFormat == "json" {
		var views []prism.PhaseRoadmapView
		for _, phase := range phases {
			view := doc.GeneratePhaseRoadmapView(phase.ID)
			if view != nil {
				views = append(views, *view)
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(views)
	}

	fmt.Println("Roadmap by Phase")
	fmt.Println("================")
	fmt.Println()

	for _, phase := range phases {
		view := doc.GeneratePhaseRoadmapView(phase.ID)
		if view == nil {
			continue
		}

		// Phase header
		period := ""
		if view.Quarter != "" {
			period = fmt.Sprintf(" (%s %d)", view.Quarter, view.Year)
		}
		fmt.Printf("%s%s\n", view.PhaseName, period)

		if view.StartDate != "" {
			fmt.Printf("  Dates: %s to %s\n", view.StartDate, view.EndDate)
		}
		if view.Status != "" {
			fmt.Printf("  Status: %s\n", view.Status)
		}
		fmt.Printf("  Completion: %.0f%%\n", view.OverallCompletion)

		// Goals in this phase
		if len(view.GoalViews) > 0 {
			fmt.Println("  Goals:")
			for _, gv := range view.GoalViews {
				fmt.Printf("    - %s (M%d -> M%d): %.0f%% complete\n",
					gv.GoalName, gv.EnterLevel, gv.ExitLevel, gv.CompletionPercent)

				// Show initiatives
				for _, init := range gv.Initiatives {
					status := init.Status
					if status == "" {
						status = "pending"
					}
					fmt.Printf("      * %s [%s]\n", init.Name, status)
				}
			}
		}
		fmt.Println()
	}

	return nil
}

func showRoadmapByGoal(doc *prism.PRISMDocument) error {
	if len(doc.Goals) == 0 {
		fmt.Println("No goals defined in document.")
		return nil
	}

	if roadmapOutputFormat == "json" {
		var views []prism.GoalRoadmapView
		for _, goal := range doc.Goals {
			view := doc.GenerateGoalRoadmapView(goal.ID)
			if view != nil {
				views = append(views, *view)
			}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(views)
	}

	fmt.Println("Roadmap by Goal")
	fmt.Println("===============")
	fmt.Println()

	for _, goal := range doc.Goals {
		view := doc.GenerateGoalRoadmapView(goal.ID)
		if view == nil {
			continue
		}

		fmt.Printf("%s (M%d -> M%d)\n", view.GoalName, view.CurrentLevel, view.TargetLevel)
		if view.Description != "" {
			fmt.Printf("  %s\n", view.Description)
		}

		if len(view.PhaseProgress) == 0 {
			fmt.Println("  No phases defined for this goal.")
		} else {
			fmt.Println("  Phase Progress:")
			for _, pp := range view.PhaseProgress {
				phaseName := pp.PhaseName
				if pp.Quarter != "" {
					phaseName = fmt.Sprintf("%s %d", pp.Quarter, pp.Year)
				}
				fmt.Printf("    %s: M%d -> M%d (%.0f%%)\n",
					phaseName, pp.EnterLevel, pp.ExitLevel, pp.CompletionPercent)

				for _, init := range pp.Initiatives {
					status := init.Status
					if status == "" {
						status = "pending"
					}
					fmt.Printf("      * %s [%s]\n", init.Name, status)
				}
			}
		}
		fmt.Println()
	}

	return nil
}

func runRoadmapProgress(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	progress := analysis.CalculateRoadmapProgress(doc)
	formatter := NewFormatter(roadmapOutputFormat)

	if roadmapOutputFormat == "json" {
		return formatter.WriteJSON(progress)
	}

	// Build phase rows
	phaseRows := make([][]string, 0, len(progress.PhaseProgress))
	for _, ps := range progress.PhaseProgress {
		phaseRows = append(phaseRows, []string{
			truncateString(ps.PhaseName, 15),
			ps.Period,
			ps.Status,
			ps.GoalSummary,
			fmt.Sprintf("%.0f%%", ps.Completion),
		})
	}

	// Build goal rows
	goalRows := make([][]string, 0, len(progress.GoalProgress))
	for _, gs := range progress.GoalProgress {
		goalRows = append(goalRows, []string{
			truncateString(gs.GoalName, 25),
			fmt.Sprintf("M%d", gs.CurrentLevel),
			fmt.Sprintf("M%d", gs.TargetLevel),
			gs.Status,
		})
	}

	// For markdown/toon, output both tables
	if roadmapOutputFormat == "markdown" || roadmapOutputFormat == "toon" {
		fmt.Fprintf(formatter.Writer, "Overall Completion: %.0f%%\n\n", progress.OverallCompletion)
		_ = formatter.WriteTable(&TableData{
			Title:   "Phase Summary",
			Headers: []string{"PHASE", "PERIOD", "STATUS", "GOALS", "COMPLETION"},
			Rows:    phaseRows,
		})
		fmt.Fprintln(formatter.Writer)
		return formatter.WriteTable(&TableData{
			Title:   "Goal Summary",
			Headers: []string{"GOAL", "CURRENT", "TARGET", "STATUS"},
			Rows:    goalRows,
		})
	}

	// Text output
	fmt.Println("Roadmap Progress")
	fmt.Println("================")
	fmt.Println()
	fmt.Printf("Overall Completion: %.0f%%\n", progress.OverallCompletion)
	fmt.Println()

	fmt.Println("Phase Summary:")
	fmt.Printf("  %-15s %-12s %-10s %-12s %-10s\n", "PHASE", "PERIOD", "STATUS", "GOALS", "COMPLETION")
	fmt.Printf("  %-15s %-12s %-10s %-12s %-10s\n", "-----", "------", "------", "-----", "----------")

	for _, ps := range progress.PhaseProgress {
		nameShort := truncateString(ps.PhaseName, 15)
		fmt.Printf("  %-15s %-12s %-10s %-12s %.0f%%\n",
			nameShort, ps.Period, ps.Status, ps.GoalSummary, ps.Completion)
	}

	fmt.Println()
	fmt.Println("Goal Summary:")
	fmt.Printf("  %-25s %-10s %-10s %-10s\n", "GOAL", "CURRENT", "TARGET", "STATUS")
	fmt.Printf("  %-25s %-10s %-10s %-10s\n", "----", "-------", "------", "------")

	for _, gs := range progress.GoalProgress {
		nameShort := truncateString(gs.GoalName, 25)
		fmt.Printf("  %-25s M%-9d M%-9d %-10s\n",
			nameShort, gs.CurrentLevel, gs.TargetLevel, gs.Status)
	}

	return nil
}
