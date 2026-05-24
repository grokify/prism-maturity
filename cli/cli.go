// Package cli provides the exported Cobra command tree for the PRISM maturity CLI.
package cli

import (
	"github.com/spf13/cobra"
)

var version = "0.1.0"

// RootCmd is the root command for PRISM maturity operations.
// It can be imported and added as a subcommand to other CLI tools.
var RootCmd = &cobra.Command{
	Use:   "maturity",
	Short: "PRISM Maturity - Maturity modeling and analysis",
	Long: `PRISM Maturity provides maturity modeling, SLO tracking, and organizational
health analysis capabilities.

Use the subcommands to create, validate, score, and analyze PRISM documents.`,
	Version: version,
}

func init() {
	// Core document commands
	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(validateCmd)
	RootCmd.AddCommand(scoreCmd)
	RootCmd.AddCommand(analyzeCmd)

	// Entity commands
	RootCmd.AddCommand(catalogCmd)
	RootCmd.AddCommand(layerCmd)
	RootCmd.AddCommand(teamCmd)
	RootCmd.AddCommand(serviceCmd)
	RootCmd.AddCommand(goalCmd)
	RootCmd.AddCommand(phaseCmd)
	RootCmd.AddCommand(roadmapCmd)
	RootCmd.AddCommand(initiativeCmd)

	// Export and reporting commands
	RootCmd.AddCommand(exportCmd)
	RootCmd.AddCommand(dashforgeCmd)
	RootCmd.AddCommand(maturityCmd)
	RootCmd.AddCommand(sloReportCmd)

	// Capability stack with maturity overlay
	RootCmd.AddCommand(capstackCmd)
}
