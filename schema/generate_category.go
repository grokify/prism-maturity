//go:build ignore

// This file generates the JSON Schema for PRISM Goal Category types from Go struct definitions.
// Run from the schema directory:
//
//	cd schema && go run generate_category.go
//
// The generated schema is used by schema/embed.go for runtime access.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-maturity/category"
	"github.com/invopop/jsonschema"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	r := &jsonschema.Reflector{
		DoNotReference: false,
		ExpandedStruct: false,
	}

	// Generate RevenueTarget schema
	if err := generateSchema(r, &category.RevenueTarget{},
		"https://github.com/grokify/prism-maturity/schema/revenue-target.schema.json",
		"Revenue Target",
		"Revenue-specific goal targets with ARR tracking",
		"revenue-target.schema.json"); err != nil {
		return err
	}

	// Generate AdoptionTarget schema
	if err := generateSchema(r, &category.AdoptionTarget{},
		"https://github.com/grokify/prism-maturity/schema/adoption-target.schema.json",
		"Adoption Target",
		"Adoption-specific goal targets with user engagement metrics",
		"adoption-target.schema.json"); err != nil {
		return err
	}

	return nil
}

func generateSchema(r *jsonschema.Reflector, typ any, id, title, description, filename string) error {
	schema := r.Reflect(typ)
	schema.ID = jsonschema.ID(id)
	schema.Title = title
	schema.Description = description

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s schema: %w", title, err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}

	fmt.Printf("Generated %s\n", filename)
	return nil
}
