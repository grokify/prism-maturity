package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-maturity"
	"github.com/spf13/cobra"
)

var phaseCmd = &cobra.Command{
	Use:   "phase",
	Short: "Work with PRISM phases",
	Long:  `Commands for listing, showing, and tracking metrics of phases in a PRISM document.`,
}

var phaseListCmd = &cobra.Command{
	Use:   "list <prism-file>",
	Short: "List all phases in a PRISM document",
	Long: `Display all phases from a PRISM document sorted by year and quarter.

Example:
  prism phase list prism.json`,
	Args: cobra.ExactArgs(1),
	RunE: runPhaseList,
}

var phaseShowCmd = &cobra.Command{
	Use:   "show <prism-file> <phase-id>",
	Short: "Show details of a specific phase",
	Long: `Display detailed information about a specific phase including goals and initiatives.

Example:
  prism phase show prism.json phase-q1-2026`,
	Args: cobra.ExactArgs(2),
	RunE: runPhaseShow,
}

var phaseMetricsCmd = &cobra.Command{
	Use:   "metrics <prism-file> <phase-id>",
	Short: "Show metrics for a specific phase",
	Long: `Display comprehensive metrics for a phase including goal progress, initiative completion, and SLO compliance.

Example:
  prism phase metrics prism.json phase-q1-2026`,
	Args: cobra.ExactArgs(2),
	RunE: runPhaseMetrics,
}

var phaseOutputFormat string

func init() {
	phaseCmd.AddCommand(phaseListCmd)
	phaseCmd.AddCommand(phaseShowCmd)
	phaseCmd.AddCommand(phaseMetricsCmd)

	phaseListCmd.Flags().StringVarP(&phaseOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
	phaseShowCmd.Flags().StringVarP(&phaseOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
	phaseMetricsCmd.Flags().StringVarP(&phaseOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
}

func runPhaseList(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	if len(doc.Phases) == 0 {
		fmt.Println("No phases defined in document.")
		return nil
	}

	phases := doc.GetPhasesSorted()
	formatter := NewFormatter(phaseOutputFormat)

	if phaseOutputFormat == "json" {
		return formatter.WriteJSON(phases)
	}

	// Build table data
	rows := make([][]string, 0, len(phases))
	for _, phase := range phases {
		quarter := ""
		if phase.Quarter != "" {
			quarter = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
		}

		status := phase.Status
		if status == "" {
			status = "planned"
		}

		startDate := phase.StartDate
		if len(startDate) > 10 {
			startDate = startDate[:10]
		}
		endDate := phase.EndDate
		if len(endDate) > 10 {
			endDate = endDate[:10]
		}

		rows = append(rows, []string{
			truncateString(phase.Name, 20),
			quarter,
			startDate,
			endDate,
			status,
			fmt.Sprintf("%d", len(phase.GoalTargets)),
		})
	}

	return formatter.WriteTable(&TableData{
		Title:   "Phases",
		Headers: []string{"NAME", "QUARTER", "START", "END", "STATUS", "GOALS"},
		Rows:    rows,
		Summary: fmt.Sprintf("Total: %d phases", len(phases)),
	})
}

func runPhaseShow(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	phaseID := args[1]
	view := doc.GeneratePhaseRoadmapView(phaseID)
	if view == nil {
		return fmt.Errorf("phase not found: %s", phaseID)
	}

	if phaseOutputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(view)
	}

	fmt.Printf("Phase: %s\n", view.PhaseName)
	fmt.Printf("ID: %s\n", view.PhaseID)
	if view.Quarter != "" {
		fmt.Printf("Period: %s %d\n", view.Quarter, view.Year)
	}
	if view.StartDate != "" {
		fmt.Printf("Dates: %s to %s\n", view.StartDate, view.EndDate)
	}
	if view.Status != "" {
		fmt.Printf("Status: %s\n", view.Status)
	}
	fmt.Printf("Overall Completion: %.0f%%\n", view.OverallCompletion)

	if len(view.GoalViews) > 0 {
		fmt.Println("\nGoals:")
		fmt.Printf("  %-25s %-10s %-10s %-12s %-10s\n", "GOAL", "ENTER", "EXIT", "INITIATIVES", "COMPLETE")
		fmt.Printf("  %-25s %-10s %-10s %-12s %-10s\n", "----", "-----", "----", "-----------", "--------")

		for _, gv := range view.GoalViews {
			nameShort := truncateString(gv.GoalName, 25)
			initiatives := fmt.Sprintf("%d/%d", gv.InitiativesCompleted, gv.InitiativesTotal)

			fmt.Printf("  %-25s M%-9d M%-9d %-12s %.0f%%\n",
				nameShort, gv.EnterLevel, gv.ExitLevel, initiatives, gv.CompletionPercent)

			// Show initiatives under each goal
			for _, init := range gv.Initiatives {
				status := init.Status
				if status == "" {
					status = "pending"
				}
				fmt.Printf("    - %s [%s]\n", init.Name, status)
			}
		}
	}

	// Show swimlanes if defined
	phase := doc.GetPhaseByID(phaseID)
	if phase != nil && len(phase.Swimlanes) > 0 {
		fmt.Println("\nSwimlanes:")
		for _, sl := range phase.Swimlanes {
			fmt.Printf("  %s", sl.Name)
			if sl.Domain != "" {
				fmt.Printf(" (%s)", sl.Domain)
			}
			fmt.Println()
			for _, initID := range sl.InitiativeIDs {
				init := doc.GetInitiativeByID(initID)
				if init != nil {
					status := init.Status
					if status == "" {
						status = "pending"
					}
					fmt.Printf("    - %s [%s]\n", init.Name, status)
				} else {
					fmt.Printf("    - %s [not found]\n", initID)
				}
			}
		}
	}

	return nil
}

func runPhaseMetrics(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	phaseID := args[1]
	phase := doc.GetPhaseByID(phaseID)
	if phase == nil {
		return fmt.Errorf("phase not found: %s", phaseID)
	}

	metrics := prism.CalculatePhaseMetrics(phase, doc)
	if metrics == nil {
		return fmt.Errorf("failed to calculate metrics for phase: %s", phaseID)
	}

	if phaseOutputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(metrics)
	}

	fmt.Printf("Phase Metrics: %s\n", phase.Name)
	fmt.Println("========================================")

	// Initiative Metrics
	if metrics.InitiativeMetrics != nil {
		im := metrics.InitiativeMetrics
		fmt.Println("\nInitiatives:")
		fmt.Printf("  Total:     %d\n", im.Total)
		fmt.Printf("  Completed: %d (%.0f%%)\n", im.Completed, safePercent(im.Completed, im.Total))
		fmt.Printf("  Deployed:  %d (%.0f%%)\n", im.Deployed, safePercent(im.Deployed, im.Total))
		if im.AvgAdoptionPercent > 0 {
			fmt.Printf("  Avg Adoption: %.0f%%\n", im.AvgAdoptionPercent)
		}
	}

	// Goal Progress
	if len(metrics.GoalProgress) > 0 {
		fmt.Println("\nGoal Progress:")
		fmt.Printf("  %-20s %-8s %-8s %-10s %-8s\n", "GOAL", "CURRENT", "TARGET", "SLOs", "INIT")
		fmt.Printf("  %-20s %-8s %-8s %-10s %-8s\n", "----", "-------", "------", "----", "----")

		for _, gp := range metrics.GoalProgress {
			goal := doc.GetGoalByID(gp.GoalID)
			name := gp.GoalID
			if goal != nil {
				name = truncateString(goal.Name, 20)
			}

			slos := fmt.Sprintf("%d/%d", gp.SLOsMet, gp.SLOsRequired)
			inits := fmt.Sprintf("%d/%d", gp.InitiativesCompleted, gp.InitiativesTotal)

			fmt.Printf("  %-20s M%-7d M%-7d %-10s %-8s\n",
				name, gp.CurrentLevel, gp.TargetLevel, slos, inits)
		}
	}

	// SLO Compliance
	if len(metrics.SLOCompliance) > 0 {
		fmt.Println("\nSLO Compliance:")
		fmt.Printf("  %-25s %-15s %-10s %-8s\n", "METRIC", "TARGET", "CURRENT", "STATUS")
		fmt.Printf("  %-25s %-15s %-10s %-8s\n", "------", "------", "-------", "------")

		met := 0
		for _, slo := range metrics.SLOCompliance {
			nameShort := truncateString(slo.MetricName, 25)
			if nameShort == "" {
				nameShort = slo.MetricID
			}

			status := "Not Met"
			if slo.IsMet {
				status = "Met"
				met++
			}

			target := slo.SLOTarget
			if len(target) > 15 {
				target = target[:12] + "..."
			}

			fmt.Printf("  %-25s %-15s %-10.2f %-8s\n",
				nameShort, target, slo.Current, status)
		}

		total := len(metrics.SLOCompliance)
		fmt.Printf("\n  SLO Compliance: %d/%d (%.0f%%)\n", met, total, safePercent(met, total))
	}

	return nil
}

// safePercent is now in format.go (re-exported from output package)
