//go:build ignore

// This file generates the JSON Schema for PRISM Maturity Plan types from Go struct definitions.
// Run from the schema directory:
//
//	cd schema && go run generate.go
//
// The generated schema is used by schema/embed.go for runtime access.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-intelligence"
	"github.com/invopop/jsonschema"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Create reflector with custom settings
	r := &jsonschema.Reflector{
		DoNotReference: false,
		ExpandedStruct: false,
	}

	// Generate schema for PRISMDocument (MaturityPlanDocument)
	schema := r.Reflect(&prism.PRISMDocument{})

	// Set schema metadata
	schema.ID = "https://github.com/grokify/prism-intelligence/schema/prism-maturity-plan.schema.json"
	schema.Title = "PRISM Maturity Plan"
	schema.Description = "PRISM Maturity Plan document defining goals, phases, initiatives, and roadmap for achieving maturity targets"

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Write to file
	filename := "prism-maturity-plan.schema.json"
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write schema file: %w", err)
	}

	fmt.Printf("Generated %s\n", filename)
	return nil
}
