package main

import (
	"fmt"

	"github.com/grokify/prism-intelligence"
	"github.com/spf13/cobra"
)

var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "List available constants and enumerations",
	Long: `Display all available PRISM constants including domains, stages,
categories, frameworks, metric types, and maturity levels.

This command helps when creating or editing PRISM documents.`,
	RunE: runCatalog,
}

func runCatalog(cmd *cobra.Command, args []string) error {
	fmt.Println("PRISM Constants Catalog")
	fmt.Println("=======================")
	fmt.Println()

	fmt.Println("Domains:")
	for _, d := range prism.AllDomains() {
		fmt.Printf("  - %s\n", d)
	}
	fmt.Println()

	fmt.Println("Lifecycle Stages:")
	for _, s := range prism.AllStages() {
		fmt.Printf("  - %s\n", s)
	}
	fmt.Println()

	fmt.Println("Categories:")
	for _, c := range prism.AllCategories() {
		fmt.Printf("  - %s\n", c)
	}
	fmt.Println()

	fmt.Println("Metric Types:")
	for _, t := range prism.AllMetricTypes() {
		fmt.Printf("  - %s\n", t)
	}
	fmt.Println()

	fmt.Println("Trend Directions:")
	for _, t := range prism.AllTrendDirections() {
		fmt.Printf("  - %s\n", t)
	}
	fmt.Println()

	fmt.Println("Status Values:")
	for _, s := range prism.AllStatuses() {
		fmt.Printf("  - %s\n", s)
	}
	fmt.Println()

	fmt.Println("SLO Windows:")
	for _, w := range prism.AllWindows() {
		fmt.Printf("  - %s\n", w)
	}
	fmt.Println()

	fmt.Println("Maturity Levels:")
	for level := prism.MaturityLevel1; level <= prism.MaturityLevel5; level++ {
		fmt.Printf("  - Level %d: %s\n", level, prism.MaturityLevelName(level))
	}
	fmt.Println()

	fmt.Println("Customer Awareness States:")
	for _, s := range prism.AllAwarenessStates() {
		fmt.Printf("  - %s\n", s)
	}
	fmt.Println()

	fmt.Println("Framework Mappings:")
	for _, f := range prism.AllFrameworks() {
		fmt.Printf("  - %s\n", f)
	}
	fmt.Println()

	fmt.Println("Layers:")
	for _, l := range prism.AllLayers() {
		fmt.Printf("  - %s\n", l)
	}
	fmt.Println()

	fmt.Println("Quality Verticals (ISO 25010):")
	for _, v := range prism.AllQualityVerticals() {
		fmt.Printf("  - %s\n", v)
	}
	fmt.Println()

	fmt.Println("Team Types:")
	for _, t := range prism.AllTeamTypes() {
		fmt.Printf("  - %s\n", t)
	}

	return nil
}
