package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-intelligence"
	"github.com/spf13/cobra"
)

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Work with PRISM teams",
	Long:  `Commands for listing and showing team definitions in a PRISM document.`,
}

var teamListCmd = &cobra.Command{
	Use:   "list <prism-file>",
	Short: "List all teams in a PRISM document",
	Long:  `Display all team definitions from a PRISM document.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runTeamList,
}

var teamShowCmd = &cobra.Command{
	Use:   "show <prism-file> <team-id>",
	Short: "Show details of a specific team",
	Long:  `Display detailed information about a specific team by ID.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runTeamShow,
}

var teamJSONOutput bool

func init() {
	teamCmd.AddCommand(teamListCmd)
	teamCmd.AddCommand(teamShowCmd)

	teamListCmd.Flags().BoolVar(&teamJSONOutput, "json", false, "Output in JSON format")
	teamShowCmd.Flags().BoolVar(&teamJSONOutput, "json", false, "Output in JSON format")
}

func runTeamList(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	if len(doc.Teams) == 0 {
		fmt.Println("No teams defined in document.")
		return nil
	}

	if teamJSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(doc.Teams)
	}

	fmt.Println("Teams:")
	fmt.Println("======")

	// Group by type
	teamsByType := make(map[string][]prism.Team)
	for _, t := range doc.Teams {
		teamsByType[t.Type] = append(teamsByType[t.Type], t)
	}

	for _, teamType := range prism.AllTeamTypes() {
		teams := teamsByType[teamType]
		if len(teams) == 0 {
			continue
		}
		fmt.Printf("\n%s:\n", formatTeamType(teamType))
		for _, t := range teams {
			fmt.Printf("  %s (%s)\n", t.Name, t.ID)
			if t.Domain != "" {
				fmt.Printf("    Domain: %s\n", t.Domain)
			}
			if len(t.ServiceIDs) > 0 {
				fmt.Printf("    Services: %d\n", len(t.ServiceIDs))
			}
		}
	}

	return nil
}

func runTeamShow(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	teamID := args[1]
	team := doc.GetTeamByID(teamID)
	if team == nil {
		return fmt.Errorf("team not found: %s", teamID)
	}

	if teamJSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(team)
	}

	fmt.Printf("Team: %s\n", team.Name)
	fmt.Printf("ID: %s\n", team.ID)
	fmt.Printf("Type: %s\n", formatTeamType(team.Type))
	if team.Description != "" {
		fmt.Printf("Description: %s\n", team.Description)
	}
	if team.Domain != "" {
		fmt.Printf("Domain: %s\n", team.Domain)
	}
	if team.Owner != "" {
		fmt.Printf("Owner: %s\n", team.Owner)
	}

	// Contact info
	if team.Slack != "" || team.Email != "" {
		fmt.Println("\nContact:")
		if team.Slack != "" {
			fmt.Printf("  Slack: %s\n", team.Slack)
		}
		if team.Email != "" {
			fmt.Printf("  Email: %s\n", team.Email)
		}
	}

	// Layer accountability
	if len(team.LayerAccountability) > 0 {
		fmt.Println("\nLayer Accountability:")
		for _, layer := range team.LayerAccountability {
			fmt.Printf("  - %s\n", layer)
		}
	}

	// Services owned
	if len(team.ServiceIDs) > 0 {
		fmt.Printf("\nServices (%d):\n", len(team.ServiceIDs))
		for _, serviceID := range team.ServiceIDs {
			service := doc.GetServiceByID(serviceID)
			if service != nil {
				fmt.Printf("  - %s (%s)\n", service.Name, service.ID)
			} else {
				fmt.Printf("  - %s (not found)\n", serviceID)
			}
		}
	}

	return nil
}

func formatTeamType(teamType string) string {
	switch teamType {
	case prism.TeamTypeStreamAligned:
		return "Stream-Aligned"
	case prism.TeamTypePlatform:
		return "Platform"
	case prism.TeamTypeEnabling:
		return "Enabling"
	case prism.TeamTypeOverlay:
		return "Overlay"
	default:
		return teamType
	}
}
