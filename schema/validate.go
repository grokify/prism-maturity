// Package schema provides JSON Schema validation for PRISM types.
package schema

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// SchemaType identifies which PRISM schema to use for validation.
type SchemaType string

const (
	SchemaMaturityModel SchemaType = "maturity-model"
	SchemaMaturityState SchemaType = "maturity-state"
	SchemaMaturityPlan  SchemaType = "maturity-plan"
)

// Validator validates JSON documents against PRISM schemas.
type Validator struct {
	schemas map[SchemaType]*jsonschema.Schema
}

// NewValidator creates a new validator with compiled schemas.
func NewValidator() (*Validator, error) {
	v := &Validator{
		schemas: make(map[SchemaType]*jsonschema.Schema),
	}

	schemaData := map[SchemaType][]byte{
		SchemaMaturityModel: maturityModelSchemaJSON,
		SchemaMaturityState: maturityStateSchemaJSON,
		SchemaMaturityPlan:  maturityPlanSchemaJSON,
	}

	for schemaType, data := range schemaData {
		schema, err := compileSchema(data)
		if err != nil {
			return nil, fmt.Errorf("compiling %s schema: %w", schemaType, err)
		}
		v.schemas[schemaType] = schema
	}

	return v, nil
}

// compileSchema compiles a JSON schema from bytes.
func compileSchema(data []byte) (*jsonschema.Schema, error) {
	// Unmarshal the schema JSON first
	var schemaDoc any
	if err := json.Unmarshal(data, &schemaDoc); err != nil {
		return nil, fmt.Errorf("unmarshaling schema: %w", err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", schemaDoc); err != nil {
		return nil, fmt.Errorf("adding schema resource: %w", err)
	}
	return c.Compile("schema.json")
}

// ValidationError represents a validation error with location and message.
type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// ValidationResult contains the result of schema validation.
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// Validate validates a JSON document against the specified schema.
func (v *Validator) Validate(schemaType SchemaType, data []byte) (*ValidationResult, error) {
	schema, ok := v.schemas[schemaType]
	if !ok {
		return nil, fmt.Errorf("unknown schema type: %s", schemaType)
	}

	var doc any
	if err := json.Unmarshal(data, &doc); err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{{
				Path:    "",
				Message: fmt.Sprintf("invalid JSON: %v", err),
			}},
		}, nil
	}

	err := schema.Validate(doc)
	if err == nil {
		return &ValidationResult{Valid: true}, nil
	}

	result := &ValidationResult{Valid: false}

	// Extract validation errors
	if ve, ok := err.(*jsonschema.ValidationError); ok {
		result.Errors = extractErrors(ve)
	} else {
		result.Errors = []ValidationError{{
			Path:    "",
			Message: err.Error(),
		}}
	}

	return result, nil
}

// extractErrors recursively extracts validation errors.
func extractErrors(ve *jsonschema.ValidationError) []ValidationError {
	var errors []ValidationError

	// Build path from InstanceLocation slice
	path := ""
	if len(ve.InstanceLocation) > 0 {
		path = "/" + strings.Join(ve.InstanceLocation, "/")
	}

	// Add this error if it has an ErrorKind
	if ve.ErrorKind != nil {
		errors = append(errors, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("%v", ve.ErrorKind),
		})
	}

	// Process nested causes
	for _, cause := range ve.Causes {
		errors = append(errors, extractErrors(cause)...)
	}

	return errors
}

// ValidateJSON is a convenience function that validates JSON against a schema type.
func ValidateJSON(schemaType SchemaType, data []byte) (*ValidationResult, error) {
	v, err := NewValidator()
	if err != nil {
		return nil, err
	}
	return v.Validate(schemaType, data)
}

// ValidateMaturityPlan validates a JSON document as a Maturity Plan.
func ValidateMaturityPlan(data []byte) (*ValidationResult, error) {
	return ValidateJSON(SchemaMaturityPlan, data)
}

// ValidateMaturityState validates a JSON document as a Maturity State.
func ValidateMaturityState(data []byte) (*ValidationResult, error) {
	return ValidateJSON(SchemaMaturityState, data)
}

// ValidateMaturityModel validates a JSON document as a Maturity Model.
func ValidateMaturityModel(data []byte) (*ValidationResult, error) {
	return ValidateJSON(SchemaMaturityModel, data)
}

// Error returns a formatted error string for the validation result.
func (r *ValidationResult) Error() string {
	if r.Valid {
		return ""
	}
	var msgs []string
	for _, e := range r.Errors {
		if e.Path != "" {
			msgs = append(msgs, fmt.Sprintf("%s: %s", e.Path, e.Message))
		} else {
			msgs = append(msgs, e.Message)
		}
	}
	return strings.Join(msgs, "; ")
}
