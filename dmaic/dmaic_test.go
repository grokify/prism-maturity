package dmaic

import (
	"testing"
	"time"
)

func TestNewDMAICDocument(t *testing.T) {
	doc := New("DMAIC-2025-001", "Security Metrics", "Alice Smith")

	if doc.Metadata == nil {
		t.Fatal("Expected metadata to be initialized")
	}
	if doc.Metadata.ID != "DMAIC-2025-001" {
		t.Errorf("Expected ID 'DMAIC-2025-001', got '%s'", doc.Metadata.ID)
	}
	if doc.Metadata.Name != "Security Metrics" {
		t.Errorf("Expected name 'Security Metrics', got '%s'", doc.Metadata.Name)
	}
	if doc.Metadata.Owner != "Alice Smith" {
		t.Errorf("Expected owner 'Alice Smith', got '%s'", doc.Metadata.Owner)
	}
	if doc.Metadata.Status != DocumentStatusDraft {
		t.Errorf("Expected status 'Draft', got '%s'", doc.Metadata.Status)
	}
	if len(doc.Categories) != 0 {
		t.Errorf("Expected empty categories, got %d", len(doc.Categories))
	}
	if len(doc.Initiatives) != 0 {
		t.Errorf("Expected empty initiatives, got %d", len(doc.Initiatives))
	}
}

func TestGenerateID(t *testing.T) {
	id := GenerateID()
	if len(id) < 14 { // "DMAIC-YYYY-DDD"
		t.Errorf("Generated ID too short: %s", id)
	}
	if id[:6] != "DMAIC-" {
		t.Errorf("Generated ID should start with 'DMAIC-': %s", id)
	}
}

func TestAllMetrics(t *testing.T) {
	doc := &DMAICDocument{
		Categories: []Category{
			{
				Name: "Application Security",
				Metrics: []Metric{
					{Name: "Critical SAST Findings"},
					{Name: "DAST Coverage"},
				},
			},
			{
				Name: "Infrastructure Security",
				Metrics: []Metric{
					{Name: "Cloud Misconfigurations"},
				},
			},
		},
	}

	metrics := doc.AllMetrics()
	if len(metrics) != 3 {
		t.Errorf("Expected 3 metrics, got %d", len(metrics))
	}
}

func TestMetricsByPhase(t *testing.T) {
	doc := &DMAICDocument{
		Categories: []Category{
			{
				Name: "Security",
				Metrics: []Metric{
					{Name: "Metric 1", Phase: PhaseDefine},
					{Name: "Metric 2", Phase: PhaseMeasure},
					{Name: "Metric 3", Phase: PhaseControl},
					{Name: "Metric 4", Phase: PhaseControl},
				},
			},
		},
	}

	byPhase := doc.MetricsByPhase()

	if len(byPhase[PhaseDefine]) != 1 {
		t.Errorf("Expected 1 Define metric, got %d", len(byPhase[PhaseDefine]))
	}
	if len(byPhase[PhaseMeasure]) != 1 {
		t.Errorf("Expected 1 Measure metric, got %d", len(byPhase[PhaseMeasure]))
	}
	if len(byPhase[PhaseControl]) != 2 {
		t.Errorf("Expected 2 Control metrics, got %d", len(byPhase[PhaseControl]))
	}
}

func TestMetricsByStatus(t *testing.T) {
	doc := &DMAICDocument{
		Categories: []Category{
			{
				Name: "Security",
				Metrics: []Metric{
					{Name: "Metric 1", Status: StatusGreen},
					{Name: "Metric 2", Status: StatusYellow},
					{Name: "Metric 3", Status: StatusRed},
					{Name: "Metric 4", Status: StatusGreen},
				},
			},
		},
	}

	byStatus := doc.MetricsByStatus()

	if len(byStatus[StatusGreen]) != 2 {
		t.Errorf("Expected 2 Green metrics, got %d", len(byStatus[StatusGreen]))
	}
	if len(byStatus[StatusYellow]) != 1 {
		t.Errorf("Expected 1 Yellow metric, got %d", len(byStatus[StatusYellow]))
	}
	if len(byStatus[StatusRed]) != 1 {
		t.Errorf("Expected 1 Red metric, got %d", len(byStatus[StatusRed]))
	}
}

func TestCalculateOverallHealth(t *testing.T) {
	tests := []struct {
		name     string
		doc      *DMAICDocument
		expected float64
	}{
		{
			name: "all green",
			doc: &DMAICDocument{
				Categories: []Category{{
					Metrics: []Metric{
						{Status: StatusGreen},
						{Status: StatusGreen},
					},
				}},
			},
			expected: 1.0,
		},
		{
			name: "all yellow",
			doc: &DMAICDocument{
				Categories: []Category{{
					Metrics: []Metric{
						{Status: StatusYellow},
						{Status: StatusYellow},
					},
				}},
			},
			expected: 0.5,
		},
		{
			name: "all red",
			doc: &DMAICDocument{
				Categories: []Category{{
					Metrics: []Metric{
						{Status: StatusRed},
						{Status: StatusRed},
					},
				}},
			},
			expected: 0.0,
		},
		{
			name: "mixed",
			doc: &DMAICDocument{
				Categories: []Category{{
					Metrics: []Metric{
						{Status: StatusGreen},  // 1.0
						{Status: StatusYellow}, // 0.5
						{Status: StatusRed},    // 0.0
					},
				}},
			},
			expected: 0.5, // (1.0 + 0.5 + 0.0) / 3
		},
		{
			name:     "empty",
			doc:      &DMAICDocument{},
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			health := tt.doc.CalculateOverallHealth()
			if health < tt.expected-0.01 || health > tt.expected+0.01 {
				t.Errorf("Expected health ~%.2f, got %.4f", tt.expected, health)
			}
		})
	}
}

func TestCalculateCategoryHealth(t *testing.T) {
	cat := Category{
		Metrics: []Metric{
			{Status: StatusGreen},  // 1.0
			{Status: StatusYellow}, // 0.5
		},
	}

	// Expected: (1.0 + 0.5) / 2 = 0.75
	health := cat.CalculateCategoryHealth()
	expected := 0.75

	if health < expected-0.01 || health > expected+0.01 {
		t.Errorf("Expected health ~%.2f, got %.4f", expected, health)
	}
}

func TestIsInControl(t *testing.T) {
	tests := []struct {
		name     string
		metric   Metric
		expected bool
	}{
		{
			name:     "no limits",
			metric:   Metric{Current: 50},
			expected: true,
		},
		{
			name: "in control",
			metric: Metric{
				Current:       50,
				ControlLimits: &ControlLimits{UCL: 60, LCL: 40, CenterLine: 50},
			},
			expected: true,
		},
		{
			name: "at UCL",
			metric: Metric{
				Current:       60,
				ControlLimits: &ControlLimits{UCL: 60, LCL: 40, CenterLine: 50},
			},
			expected: true,
		},
		{
			name: "at LCL",
			metric: Metric{
				Current:       40,
				ControlLimits: &ControlLimits{UCL: 60, LCL: 40, CenterLine: 50},
			},
			expected: true,
		},
		{
			name: "above UCL",
			metric: Metric{
				Current:       65,
				ControlLimits: &ControlLimits{UCL: 60, LCL: 40, CenterLine: 50},
			},
			expected: false,
		},
		{
			name: "below LCL",
			metric: Metric{
				Current:       35,
				ControlLimits: &ControlLimits{UCL: 60, LCL: 40, CenterLine: 50},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metric.IsInControl()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCalculateStatus(t *testing.T) {
	tests := []struct {
		name     string
		metric   Metric
		expected string
	}{
		{
			name:     "explicit status",
			metric:   Metric{Status: StatusYellow},
			expected: StatusYellow,
		},
		{
			name:     "no thresholds, target met",
			metric:   Metric{Current: 100, Target: 100, TrendDirection: TrendHigherBetter},
			expected: StatusGreen,
		},
		{
			name:     "higher better, meeting target",
			metric:   Metric{Current: 100, Target: 90, TrendDirection: TrendHigherBetter},
			expected: StatusGreen,
		},
		{
			name: "higher better, above warning",
			metric: Metric{
				Current:        85,
				Target:         90,
				TrendDirection: TrendHigherBetter,
				Thresholds:     &Thresholds{Warning: 80, Critical: 70},
			},
			expected: StatusGreen,
		},
		{
			name: "higher better, below warning",
			metric: Metric{
				Current:        75,
				Target:         90,
				TrendDirection: TrendHigherBetter,
				Thresholds:     &Thresholds{Warning: 80, Critical: 70},
			},
			expected: StatusYellow,
		},
		{
			name: "higher better, below critical",
			metric: Metric{
				Current:        65,
				Target:         90,
				TrendDirection: TrendHigherBetter,
				Thresholds:     &Thresholds{Warning: 80, Critical: 70},
			},
			expected: StatusRed,
		},
		{
			name:     "lower better, meeting target",
			metric:   Metric{Current: 5, Target: 10, TrendDirection: TrendLowerBetter},
			expected: StatusGreen,
		},
		{
			name: "lower better, above critical",
			metric: Metric{
				Current:        25,
				Target:         10,
				TrendDirection: TrendLowerBetter,
				Thresholds:     &Thresholds{Warning: 15, Critical: 20},
			},
			expected: StatusRed,
		},
		{
			name:     "target value, exact match",
			metric:   Metric{Current: 99.9, Target: 99.9, TrendDirection: TrendTargetValue},
			expected: StatusGreen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metric.CalculateStatus()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestValidateBasic(t *testing.T) {
	// Valid minimal document
	doc := &DMAICDocument{
		Categories: []Category{
			{
				Name: "Application Security",
				Metrics: []Metric{
					{Name: "Critical SAST Findings", Phase: PhaseControl},
				},
			},
		},
	}

	errs := doc.Validate(nil)
	errorCount := len(Errors(errs))
	if errorCount != 0 {
		t.Errorf("Expected no errors for valid document, got %d: %v", errorCount, errs)
	}
}

func TestValidateMissingCategories(t *testing.T) {
	doc := &DMAICDocument{
		Categories: []Category{},
	}

	errs := doc.Validate(nil)
	if len(Errors(errs)) == 0 {
		t.Error("Expected error for missing categories")
	}
}

func TestValidateMissingMetrics(t *testing.T) {
	doc := &DMAICDocument{
		Categories: []Category{
			{
				Name:    "Empty Category",
				Metrics: []Metric{},
			},
		},
	}

	errs := doc.Validate(nil)
	hasError := false
	for _, e := range errs {
		if e.Path == "categories[0].metrics" && e.IsError {
			hasError = true
			break
		}
	}
	if !hasError {
		t.Error("Expected error for missing metrics")
	}
}

func TestValidateMissingCategoryName(t *testing.T) {
	doc := &DMAICDocument{
		Categories: []Category{
			{
				Name: "", // Missing
				Metrics: []Metric{
					{Name: "Test Metric"},
				},
			},
		},
	}

	errs := doc.Validate(nil)
	hasError := false
	for _, e := range errs {
		if e.Path == "categories[0].name" && e.IsError {
			hasError = true
			break
		}
	}
	if !hasError {
		t.Error("Expected error for missing category name")
	}
}

func TestValidateInvalidPhase(t *testing.T) {
	doc := &DMAICDocument{
		Categories: []Category{
			{
				Name: "Test",
				Metrics: []Metric{
					{Name: "Test Metric", Phase: "InvalidPhase"},
				},
			},
		},
	}

	errs := doc.Validate(nil)
	hasError := false
	for _, e := range errs {
		if e.Path == "categories[0].metrics[0].phase" && e.IsError {
			hasError = true
			break
		}
	}
	if !hasError {
		t.Error("Expected error for invalid phase")
	}
}

func TestValidateControlLimits(t *testing.T) {
	doc := &DMAICDocument{
		Categories: []Category{
			{
				Name: "Test",
				Metrics: []Metric{
					{
						Name: "Test Metric",
						ControlLimits: &ControlLimits{
							UCL:        40, // Invalid: UCL < CenterLine
							LCL:        30,
							CenterLine: 50,
						},
					},
				},
			},
		},
	}

	opts := DefaultValidationOptions()
	opts.ValidateControlLimits = true
	errs := doc.Validate(opts)

	hasError := false
	for _, e := range errs {
		if e.Path == "categories[0].metrics[0].controlLimits" && e.IsError {
			hasError = true
			break
		}
	}
	if !hasError {
		t.Error("Expected error for invalid control limits")
	}
}

func TestValidateThresholdsHigherBetter(t *testing.T) {
	doc := &DMAICDocument{
		Categories: []Category{
			{
				Name: "Test",
				Metrics: []Metric{
					{
						Name:           "Test Metric",
						TrendDirection: TrendHigherBetter,
						Thresholds: &Thresholds{
							Warning:  70, // Warning < Critical (inverted for higher_better)
							Critical: 80,
						},
					},
				},
			},
		},
	}

	opts := DefaultValidationOptions()
	opts.ValidateThresholds = true
	errs := doc.Validate(opts)

	hasWarning := false
	for _, e := range errs {
		if e.Path == "categories[0].metrics[0].thresholds" && !e.IsError {
			hasWarning = true
			break
		}
	}
	if !hasWarning {
		t.Error("Expected warning for inverted thresholds")
	}
}

func TestValidateSixSigmaOptions(t *testing.T) {
	doc := &DMAICDocument{
		Categories: []Category{
			{
				Name: "Test",
				Metrics: []Metric{
					{
						Name:  "Control Phase Metric",
						Phase: PhaseControl,
						// Missing ControlLimits and ProcessCapability
					},
				},
			},
		},
	}

	opts := SixSigmaValidationOptions()
	errs := doc.Validate(opts)

	// Should have warnings for missing control limits and process capability
	if len(Warnings(errs)) == 0 {
		t.Error("Expected warnings for Six Sigma validation")
	}
}

func TestJSONRoundTrip(t *testing.T) {
	original := &DMAICDocument{
		Metadata: &Metadata{
			ID:     "DMAIC-TEST",
			Name:   "Test DMAIC",
			Owner:  "Test Owner",
			Status: DocumentStatusActive,
		},
		Categories: []Category{
			{
				ID:   "cat1",
				Name: "Security",
				Metrics: []Metric{
					{
						ID:             "m1",
						Name:           "Critical Findings",
						Phase:          PhaseControl,
						Baseline:       10,
						Current:        5,
						Target:         0,
						TrendDirection: TrendLowerBetter,
					},
				},
			},
		},
		Initiatives: []Initiative{
			{
				ID:     "init1",
				Name:   "Remediation Sprint",
				Status: InitiativeStatusInProgress,
			},
		},
	}

	// Convert to JSON
	data, err := original.JSON()
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Parse back
	parsed, err := Parse(data)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify
	if parsed.Metadata.ID != original.Metadata.ID {
		t.Errorf("ID mismatch: got %s, want %s", parsed.Metadata.ID, original.Metadata.ID)
	}
	if len(parsed.Categories) != len(original.Categories) {
		t.Errorf("Categories count mismatch: got %d, want %d", len(parsed.Categories), len(original.Categories))
	}
	if parsed.Categories[0].Metrics[0].Phase != original.Categories[0].Metrics[0].Phase {
		t.Errorf("Phase mismatch: got %s, want %s",
			parsed.Categories[0].Metrics[0].Phase,
			original.Categories[0].Metrics[0].Phase)
	}
	if len(parsed.Initiatives) != len(original.Initiatives) {
		t.Errorf("Initiatives count mismatch: got %d, want %d", len(parsed.Initiatives), len(original.Initiatives))
	}
}

func TestPhases(t *testing.T) {
	phases := Phases()
	if len(phases) != 5 {
		t.Errorf("Expected 5 phases, got %d", len(phases))
	}
	if phases[0] != PhaseDefine {
		t.Errorf("Expected first phase to be Define, got %s", phases[0])
	}
	if phases[4] != PhaseControl {
		t.Errorf("Expected last phase to be Control, got %s", phases[4])
	}
}

func TestValidPhase(t *testing.T) {
	tests := []struct {
		phase    string
		expected bool
	}{
		{PhaseDefine, true},
		{PhaseMeasure, true},
		{PhaseAnalyze, true},
		{PhaseImprove, true},
		{PhaseControl, true},
		{"Invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		result := ValidPhase(tt.phase)
		if result != tt.expected {
			t.Errorf("ValidPhase(%q) = %v, expected %v", tt.phase, result, tt.expected)
		}
	}
}

func TestValidTrendDirection(t *testing.T) {
	tests := []struct {
		trend    string
		expected bool
	}{
		{TrendHigherBetter, true},
		{TrendLowerBetter, true},
		{TrendTargetValue, true},
		{"", true}, // Empty is valid (defaults)
		{"invalid", false},
	}

	for _, tt := range tests {
		result := ValidTrendDirection(tt.trend)
		if result != tt.expected {
			t.Errorf("ValidTrendDirection(%q) = %v, expected %v", tt.trend, result, tt.expected)
		}
	}
}

func TestValidStatus(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{StatusGreen, true},
		{StatusYellow, true},
		{StatusRed, true},
		{"", true}, // Empty is valid
		{"invalid", false},
	}

	for _, tt := range tests {
		result := ValidStatus(tt.status)
		if result != tt.expected {
			t.Errorf("ValidStatus(%q) = %v, expected %v", tt.status, result, tt.expected)
		}
	}
}

func TestStatusDescription(t *testing.T) {
	tests := []struct {
		status   string
		expected string
	}{
		{StatusGreen, "Meeting target"},
		{StatusYellow, "Warning - attention needed"},
		{StatusRed, "Critical - immediate action required"},
		{"", "Unknown status"},
	}

	for _, tt := range tests {
		result := StatusDescription(tt.status)
		if result != tt.expected {
			t.Errorf("StatusDescription(%q) = %q, expected %q", tt.status, result, tt.expected)
		}
	}
}

func TestPhaseDescription(t *testing.T) {
	tests := []struct {
		phase    string
		expected string
	}{
		{PhaseDefine, "Define the problem and project goals"},
		{PhaseMeasure, "Measure current performance"},
		{PhaseAnalyze, "Analyze root causes"},
		{PhaseImprove, "Implement improvements"},
		{PhaseControl, "Sustain the gains"},
		{"", "Unknown phase"},
	}

	for _, tt := range tests {
		result := PhaseDescription(tt.phase)
		if result != tt.expected {
			t.Errorf("PhaseDescription(%q) = %q, expected %q", tt.phase, result, tt.expected)
		}
	}
}

func TestDataPointsHistory(t *testing.T) {
	metric := Metric{
		Name: "Test Metric",
		DataPoints: []DataPoint{
			{Timestamp: time.Now().Add(-24 * time.Hour), Value: 10, Notes: "Initial"},
			{Timestamp: time.Now(), Value: 15, Notes: "Updated"},
		},
	}

	if len(metric.DataPoints) != 2 {
		t.Errorf("Expected 2 data points, got %d", len(metric.DataPoints))
	}
	if metric.DataPoints[1].Value != 15 {
		t.Errorf("Expected latest value 15, got %.2f", metric.DataPoints[1].Value)
	}
}

func TestRootCauses(t *testing.T) {
	metric := Metric{
		Name: "Test Metric",
		RootCauses: []RootCause{
			{ID: "rc1", Description: "Process bottleneck", Category: "Process", Impact: "High", Validated: true},
			{ID: "rc2", Description: "Skill gap", Category: "People", Impact: "Medium", Validated: false},
		},
	}

	if len(metric.RootCauses) != 2 {
		t.Errorf("Expected 2 root causes, got %d", len(metric.RootCauses))
	}
	if !metric.RootCauses[0].Validated {
		t.Error("Expected first root cause to be validated")
	}
}

func TestErrorsAndWarnings(t *testing.T) {
	errs := []ValidationError{
		{Path: "test1", Message: "error 1", IsError: true},
		{Path: "test2", Message: "warning 1", IsError: false},
		{Path: "test3", Message: "error 2", IsError: true},
		{Path: "test4", Message: "warning 2", IsError: false},
	}

	errors := Errors(errs)
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errors))
	}

	warnings := Warnings(errs)
	if len(warnings) != 2 {
		t.Errorf("Expected 2 warnings, got %d", len(warnings))
	}

	if IsValid(errs) {
		t.Error("Expected IsValid to return false when errors exist")
	}

	if !IsValid(warnings) {
		t.Error("Expected IsValid to return true when only warnings exist")
	}
}
