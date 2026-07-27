package schema

import (
	"encoding/json"
	"testing"

	prism "github.com/grokify/prism-maturity"
)

func TestNewValidator(t *testing.T) {
	v, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}
	if v == nil {
		t.Fatal("NewValidator() returned nil")
	}
	if len(v.schemas) != 3 {
		t.Errorf("expected 3 schemas, got %d", len(v.schemas))
	}
}

func TestValidateMaturityPlan_ValidDocument(t *testing.T) {
	doc := prism.PRISMDocument{
		Layers:   []prism.LayerDef{},
		Services: []prism.Service{},
		Goals: []prism.Goal{
			{ID: "goal-1", Name: "Test Goal", Status: prism.GoalStatusActive},
		},
		Metrics: []prism.Metric{},
	}

	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshaling document: %v", err)
	}

	result, err := ValidateMaturityPlan(data)
	if err != nil {
		t.Fatalf("ValidateMaturityPlan() error = %v", err)
	}

	if !result.Valid {
		t.Errorf("expected valid document, got errors: %s", result.Error())
	}
}

func TestValidateMaturityPlan_InvalidJSON(t *testing.T) {
	data := []byte(`{invalid json`)

	result, err := ValidateMaturityPlan(data)
	if err != nil {
		t.Fatalf("ValidateMaturityPlan() error = %v", err)
	}

	if result.Valid {
		t.Error("expected invalid result for invalid JSON")
	}
	if len(result.Errors) == 0 {
		t.Error("expected errors for invalid JSON")
	}
}

func TestValidateMaturityPlan_EmptyDocument(t *testing.T) {
	data := []byte(`{}`)

	result, err := ValidateMaturityPlan(data)
	if err != nil {
		t.Fatalf("ValidateMaturityPlan() error = %v", err)
	}

	// Empty document may be valid depending on schema requirements
	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestValidate_UnknownSchemaType(t *testing.T) {
	v, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}

	_, err = v.Validate("unknown", []byte(`{}`))
	if err == nil {
		t.Error("expected error for unknown schema type")
	}
}

func TestValidationResult_Error(t *testing.T) {
	t.Run("valid result", func(t *testing.T) {
		r := &ValidationResult{Valid: true}
		if r.Error() != "" {
			t.Errorf("expected empty error for valid result, got %q", r.Error())
		}
	})

	t.Run("single error", func(t *testing.T) {
		r := &ValidationResult{
			Valid:  false,
			Errors: []ValidationError{{Path: "/foo", Message: "required"}},
		}
		if r.Error() == "" {
			t.Error("expected non-empty error for invalid result")
		}
	})

	t.Run("multiple errors", func(t *testing.T) {
		r := &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{Path: "/foo", Message: "required"},
				{Path: "/bar", Message: "invalid type"},
			},
		}
		errStr := r.Error()
		if errStr == "" {
			t.Error("expected non-empty error for invalid result")
		}
	})
}

func TestValidateMaturityPlan_ComplexDocument(t *testing.T) {
	doc := prism.PRISMDocument{
		Layers: []prism.LayerDef{
			{ID: "code", Name: "Code Quality", Weight: 0.3},
			{ID: "infra", Name: "Infrastructure", Weight: 0.4},
		},
		Services: []prism.Service{
			{ID: "api", Name: "API Service", LayerID: "code", Tier: "tier1"},
		},
		Goals: []prism.Goal{
			{
				ID:           "goal-1",
				Name:         "Improve Reliability",
				Status:       prism.GoalStatusActive,
				CurrentLevel: 2,
				TargetLevel:  4,
			},
		},
		Metrics: []prism.Metric{
			{
				ID:      "uptime",
				Name:    "Uptime",
				Current: 99.5,
				Target:  99.9,
				SLO:     &prism.SLO{Value: 99.9, Operator: prism.SLOOperatorGTE},
			},
		},
	}

	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("marshaling document: %v", err)
	}

	result, err := ValidateMaturityPlan(data)
	if err != nil {
		t.Fatalf("ValidateMaturityPlan() error = %v", err)
	}

	if !result.Valid {
		t.Errorf("expected valid complex document, got errors: %s", result.Error())
	}
}

func TestValidator_Reuse(t *testing.T) {
	v, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}

	// Validate multiple documents with the same validator
	for i := 0; i < 3; i++ {
		doc := prism.PRISMDocument{
			Layers:   []prism.LayerDef{},
			Services: []prism.Service{},
			Goals: []prism.Goal{
				{ID: "goal-1", Name: "Test Goal"},
			},
			Metrics: []prism.Metric{},
		}
		data, _ := json.Marshal(doc)

		result, err := v.Validate(SchemaMaturityPlan, data)
		if err != nil {
			t.Fatalf("iteration %d: Validate() error = %v", i, err)
		}
		if !result.Valid {
			t.Errorf("iteration %d: expected valid, got errors: %s", i, result.Error())
		}
	}
}
