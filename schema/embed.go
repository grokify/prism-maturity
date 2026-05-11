// Package schema provides embedded JSON Schema for PRISM types.
package schema

import (
	_ "embed"
	"encoding/json"
)

// Maturity Model schema (what good looks like)
//
//go:embed prism-maturity-model.schema.json
var maturityModelSchemaJSON []byte

// Maturity State schema (where we are now)
//
//go:embed prism-maturity-state.schema.json
var maturityStateSchemaJSON []byte

// Maturity Plan schema (how we get there)
//
//go:embed prism-maturity-plan.schema.json
var maturityPlanSchemaJSON []byte

// MaturityModelSchemaJSON returns the raw JSON Schema bytes for Maturity Model.
func MaturityModelSchemaJSON() []byte {
	return maturityModelSchemaJSON
}

// MaturityStateSchemaJSON returns the raw JSON Schema bytes for Maturity State.
func MaturityStateSchemaJSON() []byte {
	return maturityStateSchemaJSON
}

// MaturityPlanSchemaJSON returns the raw JSON Schema bytes for Maturity Plan.
func MaturityPlanSchemaJSON() []byte {
	return maturityPlanSchemaJSON
}

// MaturityModelSchemaMap returns the Maturity Model schema as a map for programmatic access.
func MaturityModelSchemaMap() (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(maturityModelSchemaJSON, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// MaturityStateSchemaMap returns the Maturity State schema as a map for programmatic access.
func MaturityStateSchemaMap() (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(maturityStateSchemaJSON, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// MaturityPlanSchemaMap returns the Maturity Plan schema as a map for programmatic access.
func MaturityPlanSchemaMap() (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(maturityPlanSchemaJSON, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Schema IDs
const (
	MaturityModelSchemaID = "https://github.com/grokify/prism/schema/prism-maturity-model.schema.json"
	MaturityStateSchemaID = "https://github.com/grokify/prism/schema/prism-maturity-state.schema.json"
	MaturityPlanSchemaID  = "https://github.com/grokify/prism/schema/prism-maturity-plan.schema.json"
)
