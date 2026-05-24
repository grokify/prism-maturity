//go:build ignore

// This file generates the JSON Schema for maturity model types from Go struct definitions.
// Run from the schema directory:
//
//	cd schema && go run generate_maturity.go
//
// Or generate all schemas:
//
//	cd schema && go run generate.go && go run generate_maturity.go && go run generate_state.go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-maturity/maturity"
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

	// Generate schema for maturity.Spec
	schema := r.Reflect(&maturity.Spec{})

	// Set schema metadata
	schema.ID = "https://github.com/grokify/prism-maturity/schema/prism-maturity-model.schema.json"
	schema.Title = "PRISM Maturity Model"
	schema.Description = "PRISM Maturity Model specification defining what good looks like at each maturity level (M1-M5)"

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Write to file
	filename := "prism-maturity-model.schema.json"
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write schema file: %w", err)
	}

	fmt.Printf("Generated %s\n", filename)
	return nil
}
