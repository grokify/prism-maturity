package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/grokify/prism-intelligence"
	"github.com/spf13/cobra"
)

var initiativeCmd = &cobra.Command{
	Use:   "initiative",
	Short: "Work with PRISM initiatives",
	Long:  `Commands for listing and showing initiatives in a PRISM document.`,
}

var initiativeListCmd = &cobra.Command{
	Use:   "list <prism-file>",
	Short: "List all initiatives in a PRISM document",
	Long: `Display all initiatives from a PRISM document grouped by status or phase.

Example:
  prism initiative list prism.json
  prism initiative list prism.json --by-phase
  prism initiative list prism.json --by-goal`,
	Args: cobra.ExactArgs(1),
	RunE: runInitiativeList,
}

var initiativeShowCmd = &cobra.Command{
	Use:   "show <prism-file> <initiative-id>",
	Short: "Show details of a specific initiative",
	Long: `Display detailed information about a specific initiative including
linked goals, phase, deployment status, and metrics.

Example:
  prism initiative show prism.json init-monitoring`,
	Args: cobra.ExactArgs(2),
	RunE: runInitiativeShow,
}

var (
	initiativeOutputFormat string
	initiativeByPhase      bool
	initiativeByGoal       bool
)

func init() {
	initiativeCmd.AddCommand(initiativeListCmd)
	initiativeCmd.AddCommand(initiativeShowCmd)

	initiativeListCmd.Flags().StringVarP(&initiativeOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
	initiativeListCmd.Flags().BoolVar(&initiativeByPhase, "by-phase", false, "Group initiatives by phase")
	initiativeListCmd.Flags().BoolVar(&initiativeByGoal, "by-goal", false, "Group initiatives by goal")

	initiativeShowCmd.Flags().StringVarP(&initiativeOutputFormat, "format", "f", "text", "Output format (json|text|markdown|toon)")
}

func runInitiativeList(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	if len(doc.Initiatives) == 0 {
		fmt.Println("No initiatives defined in document.")
		return nil
	}

	formatter := NewFormatter(initiativeOutputFormat)

	if initiativeOutputFormat == "json" {
		return formatter.WriteJSON(doc.Initiatives)
	}

	if initiativeByPhase {
		return printInitiativesByPhase(doc)
	}

	if initiativeByGoal {
		return printInitiativesByGoal(doc)
	}

	// Default: group by status
	return printInitiativesByStatus(doc)
}

func printInitiativesByStatus(doc *prism.PRISMDocument) error {
	// Group by status
	byStatus := make(map[string][]prism.Initiative)
	for _, init := range doc.Initiatives {
		status := init.Status
		if status == "" {
			status = "not_started"
		}
		byStatus[status] = append(byStatus[status], init)
	}

	fmt.Println("Initiatives by Status")
	fmt.Println("=====================")

	// Print in order: in_progress, planned, not_started, completed, cancelled
	statusOrder := []string{
		prism.InitiativeStatusInProgress,
		prism.InitiativeStatusPlanned,
		prism.InitiativeStatusNotStarted,
		prism.InitiativeStatusCompleted,
		prism.InitiativeStatusCancelled,
	}

	for _, status := range statusOrder {
		initiatives := byStatus[status]
		if len(initiatives) == 0 {
			continue
		}

		fmt.Printf("\n%s (%d):\n", FormatInitiativeStatus(status), len(initiatives))
		for _, init := range initiatives {
			printInitiativeSummary(doc, init, "  ")
		}
	}

	fmt.Printf("\nTotal: %d initiatives\n", len(doc.Initiatives))
	return nil
}

func printInitiativesByPhase(doc *prism.PRISMDocument) error {
	// Group by phase
	byPhase := make(map[string][]prism.Initiative)
	noPhase := []prism.Initiative{}

	for _, init := range doc.Initiatives {
		if init.PhaseID != "" {
			byPhase[init.PhaseID] = append(byPhase[init.PhaseID], init)
		} else {
			noPhase = append(noPhase, init)
		}
	}

	fmt.Println("Initiatives by Phase")
	fmt.Println("====================")

	// Print phases in order
	phases := doc.GetPhasesSorted()
	for _, phase := range phases {
		initiatives := byPhase[phase.ID]
		if len(initiatives) == 0 {
			continue
		}

		phaseName := phase.Name
		if phase.Quarter != "" {
			phaseName = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
		}
		fmt.Printf("\n%s (%d initiatives):\n", phaseName, len(initiatives))
		for _, init := range initiatives {
			printInitiativeSummary(doc, init, "  ")
		}
	}

	if len(noPhase) > 0 {
		fmt.Printf("\nNo Phase Assigned (%d):\n", len(noPhase))
		for _, init := range noPhase {
			printInitiativeSummary(doc, init, "  ")
		}
	}

	fmt.Printf("\nTotal: %d initiatives\n", len(doc.Initiatives))
	return nil
}

func printInitiativesByGoal(doc *prism.PRISMDocument) error {
	// Group by goal
	byGoal := make(map[string][]prism.Initiative)
	noGoal := []prism.Initiative{}

	for _, init := range doc.Initiatives {
		if len(init.GoalIDs) > 0 {
			for _, goalID := range init.GoalIDs {
				byGoal[goalID] = append(byGoal[goalID], init)
			}
		} else {
			noGoal = append(noGoal, init)
		}
	}

	fmt.Println("Initiatives by Goal")
	fmt.Println("===================")

	for _, goal := range doc.Goals {
		initiatives := byGoal[goal.ID]
		if len(initiatives) == 0 {
			continue
		}

		fmt.Printf("\n%s (%d initiatives):\n", goal.Name, len(initiatives))
		for _, init := range initiatives {
			printInitiativeSummary(doc, init, "  ")
		}
	}

	if len(noGoal) > 0 {
		fmt.Printf("\nNo Goal Linked (%d):\n", len(noGoal))
		for _, init := range noGoal {
			printInitiativeSummary(doc, init, "  ")
		}
	}

	fmt.Printf("\nTotal: %d initiatives\n", len(doc.Initiatives))
	return nil
}

//nolint:unparam // doc reserved for future use (team lookup, etc.)
func printInitiativeSummary(doc *prism.PRISMDocument, init prism.Initiative, indent string) {
	status := init.Status
	if status == "" {
		status = "not_started"
	}

	// Build progress indicator
	progress := ""
	if init.DevCompletionPercent > 0 {
		progress = fmt.Sprintf(" (%.0f%% dev", init.DevCompletionPercent)
		if init.DeploymentStatus != nil && init.DeploymentStatus.AdoptionPercent > 0 {
			progress += fmt.Sprintf(", %.0f%% deployed", init.DeploymentStatus.AdoptionPercent)
		}
		progress += ")"
	}

	fmt.Printf("%s%s [%s]%s\n", indent, init.Name, status, progress)

	// Show owner if present
	if init.Owner != "" {
		fmt.Printf("%s  Owner: %s\n", indent, init.Owner)
	}
}

func runInitiativeShow(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	initID := args[1]
	init := doc.GetInitiativeByID(initID)
	if init == nil {
		return fmt.Errorf("initiative not found: %s", initID)
	}

	if initiativeOutputFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(init)
	}

	// Text output
	fmt.Printf("Initiative: %s\n", init.Name)
	fmt.Printf("ID: %s\n", init.ID)
	if init.Description != "" {
		fmt.Printf("Description: %s\n", init.Description)
	}

	status := init.Status
	if status == "" {
		status = "not_started"
	}
	fmt.Printf("Status: %s\n", FormatInitiativeStatus(status))

	if init.Priority > 0 {
		fmt.Printf("Priority: %d\n", init.Priority)
	}
	if init.Owner != "" {
		fmt.Printf("Owner: %s\n", init.Owner)
	}
	if init.Team != "" {
		fmt.Printf("Team: %s\n", init.Team)
	}

	// Dates
	if init.StartDate != "" || init.EndDate != "" {
		dates := ""
		if init.StartDate != "" {
			dates = init.StartDate
		}
		if init.EndDate != "" {
			if dates != "" {
				dates += " to " + init.EndDate
			} else {
				dates = "to " + init.EndDate
			}
		}
		fmt.Printf("Dates: %s\n", dates)
	}

	// Development progress
	fmt.Println("\nProgress:")
	fmt.Printf("  Dev Completion: %.0f%%\n", init.DevCompletionPercent)

	// Deployment status
	if init.DeploymentStatus != nil {
		ds := init.DeploymentStatus
		fmt.Println("\nDeployment:")
		fmt.Printf("  Status: %s\n", FormatInitiativeStatus(ds.Status))
		if ds.TotalCustomers > 0 {
			fmt.Printf("  Customers: %d/%d (%.0f%%)\n",
				ds.DeployedCustomers, ds.TotalCustomers, ds.AdoptionPercent)
		}
	}

	// Linked goals
	if len(init.GoalIDs) > 0 {
		fmt.Printf("\nLinked Goals (%d):\n", len(init.GoalIDs))
		for _, goalID := range init.GoalIDs {
			goal := doc.GetGoalByID(goalID)
			if goal != nil {
				currentLevel := goal.CurrentLevel
				if currentLevel == 0 {
					currentLevel = goal.CurrentMaturityLevel(doc)
				}
				fmt.Printf("  - %s (M%d → M%d)\n", goal.Name, currentLevel, goal.TargetLevel)
			} else {
				fmt.Printf("  - %s (not found)\n", goalID)
			}
		}
	}

	// Linked phase
	if init.PhaseID != "" {
		phase := doc.GetPhaseByID(init.PhaseID)
		if phase != nil {
			phaseName := phase.Name
			if phase.Quarter != "" {
				phaseName = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
			}
			fmt.Printf("\nPhase: %s\n", phaseName)
		} else {
			fmt.Printf("\nPhase: %s (not found)\n", init.PhaseID)
		}
	}

	// Linked service
	if init.ServiceID != "" {
		service := doc.GetServiceByID(init.ServiceID)
		if service != nil {
			fmt.Printf("\nService: %s\n", service.Name)
		} else {
			fmt.Printf("\nService: %s (not found)\n", init.ServiceID)
		}
	}

	// Linked metrics
	if len(init.MetricIDs) > 0 {
		fmt.Printf("\nLinked Metrics (%d):\n", len(init.MetricIDs))
		for _, metricID := range init.MetricIDs {
			metric := doc.GetMetricByID(metricID)
			if metric != nil {
				sloStatus := ""
				if metric.SLO != nil {
					if metric.MeetsSLO() {
						sloStatus = " [SLO Met]"
					} else {
						sloStatus = " [SLO Not Met]"
					}
				}
				fmt.Printf("  - %s%s\n", metric.Name, sloStatus)
			} else {
				fmt.Printf("  - %s (not found)\n", metricID)
			}
		}
	}

	// Dependent teams
	if len(init.DependentTeams) > 0 {
		fmt.Printf("\nDependent Teams: %s\n", strings.Join(init.DependentTeams, ", "))
	}

	return nil
}

// formatStatus is re-exported from output package as FormatInitiativeStatus
