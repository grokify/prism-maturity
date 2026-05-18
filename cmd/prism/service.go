package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-intelligence"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Work with PRISM services",
	Long:  `Commands for listing and showing service definitions in a PRISM document.`,
}

var serviceListCmd = &cobra.Command{
	Use:   "list <prism-file>",
	Short: "List all services in a PRISM document",
	Long:  `Display all service definitions from a PRISM document.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runServiceList,
}

var serviceShowCmd = &cobra.Command{
	Use:   "show <prism-file> <service-id>",
	Short: "Show details of a specific service",
	Long:  `Display detailed information about a specific service by ID.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runServiceShow,
}

var serviceJSONOutput bool

func init() {
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceShowCmd)

	serviceListCmd.Flags().BoolVar(&serviceJSONOutput, "json", false, "Output in JSON format")
	serviceShowCmd.Flags().BoolVar(&serviceJSONOutput, "json", false, "Output in JSON format")
}

func runServiceList(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	if len(doc.Services) == 0 {
		fmt.Println("No services defined in document.")
		return nil
	}

	if serviceJSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(doc.Services)
	}

	fmt.Println("Services:")
	fmt.Println("=========")

	// Group by layer
	servicesByLayer := make(map[string][]prism.Service)
	noLayer := []prism.Service{}
	for _, s := range doc.Services {
		if s.LayerID != "" {
			servicesByLayer[s.LayerID] = append(servicesByLayer[s.LayerID], s)
		} else {
			noLayer = append(noLayer, s)
		}
	}

	for _, layer := range prism.AllLayers() {
		services := servicesByLayer[layer]
		if len(services) == 0 {
			continue
		}
		fmt.Printf("\n%s layer:\n", layer)
		for _, s := range services {
			printServiceSummary(doc, s)
		}
	}

	if len(noLayer) > 0 {
		fmt.Println("\nNo layer assigned:")
		for _, s := range noLayer {
			printServiceSummary(doc, s)
		}
	}

	return nil
}

func printServiceSummary(doc *prism.PRISMDocument, s prism.Service) {
	fmt.Printf("  %s (%s)\n", s.Name, s.ID)
	if s.OwnerTeamID != "" {
		team := doc.GetTeamByID(s.OwnerTeamID)
		if team != nil {
			fmt.Printf("    Owner: %s\n", team.Name)
		} else {
			fmt.Printf("    Owner: %s (team not found)\n", s.OwnerTeamID)
		}
	}
	if s.Tier != "" {
		fmt.Printf("    Tier: %s\n", s.Tier)
	}
}

func runServiceShow(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	serviceID := args[1]
	service := doc.GetServiceByID(serviceID)
	if service == nil {
		return fmt.Errorf("service not found: %s", serviceID)
	}

	if serviceJSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(service)
	}

	fmt.Printf("Service: %s\n", service.Name)
	fmt.Printf("ID: %s\n", service.ID)
	if service.Description != "" {
		fmt.Printf("Description: %s\n", service.Description)
	}
	if service.LayerID != "" {
		fmt.Printf("Layer: %s\n", service.LayerID)
	}
	if service.Tier != "" {
		fmt.Printf("Tier: %s\n", service.Tier)
	}
	if service.Repository != "" {
		fmt.Printf("Repository: %s\n", service.Repository)
	}

	// Owner team
	if service.OwnerTeamID != "" {
		team := doc.GetTeamByID(service.OwnerTeamID)
		if team != nil {
			fmt.Printf("\nOwner Team: %s (%s)\n", team.Name, team.ID)
		} else {
			fmt.Printf("\nOwner Team ID: %s (not found)\n", service.OwnerTeamID)
		}
	}

	// Metrics
	metrics := doc.GetMetricsByService(serviceID)
	if len(metrics) > 0 {
		fmt.Printf("\nMetrics (%d):\n", len(metrics))
		for _, m := range metrics {
			status := ""
			if m.Status != "" {
				status = fmt.Sprintf(" [%s]", m.Status)
			}
			fmt.Printf("  - %s%s\n", m.Name, status)
		}
	}

	// Also show metrics from service.MetricIDs if different
	if len(service.MetricIDs) > 0 {
		fmt.Printf("\nLinked Metrics (%d):\n", len(service.MetricIDs))
		for _, metricID := range service.MetricIDs {
			metric := doc.GetMetricByID(metricID)
			if metric != nil {
				status := ""
				if metric.Status != "" {
					status = fmt.Sprintf(" [%s]", metric.Status)
				}
				fmt.Printf("  - %s%s\n", metric.Name, status)
			} else {
				fmt.Printf("  - %s (not found)\n", metricID)
			}
		}
	}

	return nil
}
