package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-intelligence"
	"github.com/spf13/cobra"
)

var layerCmd = &cobra.Command{
	Use:   "layer",
	Short: "Work with PRISM layers",
	Long:  `Commands for listing and showing layer definitions in a PRISM document.`,
}

var layerListCmd = &cobra.Command{
	Use:   "list <prism-file>",
	Short: "List all layers in a PRISM document",
	Long:  `Display all layer definitions from a PRISM document.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runLayerList,
}

var layerShowCmd = &cobra.Command{
	Use:   "show <prism-file> <layer-id>",
	Short: "Show details of a specific layer",
	Long:  `Display detailed information about a specific layer by ID.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runLayerShow,
}

var layerJSONOutput bool

func init() {
	layerCmd.AddCommand(layerListCmd)
	layerCmd.AddCommand(layerShowCmd)

	layerListCmd.Flags().BoolVar(&layerJSONOutput, "json", false, "Output in JSON format")
	layerShowCmd.Flags().BoolVar(&layerJSONOutput, "json", false, "Output in JSON format")
}

func runLayerList(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	if len(doc.Layers) == 0 {
		fmt.Println("No layers defined in document.")
		fmt.Println("\nDefault layers available:")
		for _, l := range prism.DefaultLayers() {
			fmt.Printf("  - %s: %s\n", l.ID, l.Name)
		}
		return nil
	}

	if layerJSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(doc.Layers)
	}

	fmt.Println("Layers:")
	fmt.Println("=======")
	for _, l := range doc.Layers {
		fmt.Printf("\n%s (%s)\n", l.Name, l.ID)
		if l.Description != "" {
			fmt.Printf("  Description: %s\n", l.Description)
		}
		if l.Weight > 0 {
			fmt.Printf("  Weight: %.2f\n", l.Weight)
		}
		if l.Signals.Latency != "" || l.Signals.Traffic != "" ||
			l.Signals.Errors != "" || l.Signals.Saturation != "" {
			fmt.Println("  Golden Signals:")
			if l.Signals.Latency != "" {
				fmt.Printf("    - Latency: %s\n", l.Signals.Latency)
			}
			if l.Signals.Traffic != "" {
				fmt.Printf("    - Traffic: %s\n", l.Signals.Traffic)
			}
			if l.Signals.Errors != "" {
				fmt.Printf("    - Errors: %s\n", l.Signals.Errors)
			}
			if l.Signals.Saturation != "" {
				fmt.Printf("    - Saturation: %s\n", l.Signals.Saturation)
			}
		}
	}

	return nil
}

func runLayerShow(cmd *cobra.Command, args []string) error {
	doc, err := loadPRISMDocument(args[0])
	if err != nil {
		return err
	}

	layerID := args[1]
	layer := doc.GetLayerByID(layerID)
	if layer == nil {
		return fmt.Errorf("layer not found: %s", layerID)
	}

	if layerJSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(layer)
	}

	fmt.Printf("Layer: %s\n", layer.Name)
	fmt.Printf("ID: %s\n", layer.ID)
	if layer.Description != "" {
		fmt.Printf("Description: %s\n", layer.Description)
	}
	if layer.Weight > 0 {
		fmt.Printf("Weight: %.2f\n", layer.Weight)
	}

	// Show golden signals
	if layer.Signals.Latency != "" || layer.Signals.Traffic != "" ||
		layer.Signals.Errors != "" || layer.Signals.Saturation != "" {
		fmt.Println("\nGolden Signals:")
		if layer.Signals.Latency != "" {
			fmt.Printf("  Latency: %s\n", layer.Signals.Latency)
		}
		if layer.Signals.Traffic != "" {
			fmt.Printf("  Traffic: %s\n", layer.Signals.Traffic)
		}
		if layer.Signals.Errors != "" {
			fmt.Printf("  Errors: %s\n", layer.Signals.Errors)
		}
		if layer.Signals.Saturation != "" {
			fmt.Printf("  Saturation: %s\n", layer.Signals.Saturation)
		}
	}

	// Show metrics in this layer
	metrics := doc.GetMetricsByLayer(layerID)
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

	return nil
}

func loadPRISMDocument(filename string) (*prism.PRISMDocument, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var doc prism.PRISMDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &doc, nil
}
