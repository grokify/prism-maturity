package maturity

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateXLSX(t *testing.T) {
	// Read the security maturity model
	specFile := filepath.Join("..", "examples", "security", "model.json")
	spec, err := ReadSpecFile(specFile)
	if err != nil {
		t.Fatalf("Failed to read spec file: %v", err)
	}

	// Verify spec was parsed correctly
	if spec == nil {
		t.Fatal("Spec is nil")
	}

	if len(spec.Domains) == 0 {
		t.Fatal("No domains in spec")
	}

	securityDomain, ok := spec.Domains["security"]
	if !ok {
		t.Fatal("Security domain not found")
	}

	if len(securityDomain.Levels) != 5 {
		t.Errorf("Expected 5 levels, got %d", len(securityDomain.Levels))
	}

	// Test level criteria
	level3, ok := securityDomain.GetLevel(3)
	if !ok {
		t.Fatal("Level 3 not found")
	}

	if len(level3.Criteria) == 0 {
		t.Error("Level 3 has no criteria")
	}

	if len(level3.Enablers) == 0 {
		t.Error("Level 3 has no enablers")
	}

	// Test criterion checking
	for _, c := range level3.Criteria {
		if c.ID == "" {
			t.Error("Criterion has empty ID")
		}
		// MetricName can be resolved from SLI reference
		if c.GetMetricName(spec) == "" && c.SLIID == "" {
			t.Errorf("Criterion %s has no MetricName and no SLI reference", c.ID)
		}
	}

	// Generate XLSX
	gen := NewXLSXGenerator(spec)
	if err := gen.Generate(); err != nil {
		t.Fatalf("Failed to generate XLSX: %v", err)
	}

	// Save to temp file
	tmpFile := filepath.Join(os.TempDir(), "maturity-test.xlsx")
	if err := gen.SaveAs(tmpFile); err != nil {
		t.Fatalf("Failed to save XLSX: %v", err)
	}

	// Verify file was created
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Output file is empty")
	}

	t.Logf("Generated XLSX file: %s (%d bytes)", tmpFile, info.Size())

	// Cleanup
	os.Remove(tmpFile)
}

func TestCriterionCheckMet(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		target   float64
		current  float64
		expected bool
	}{
		{"gte met", OpGTE, 80, 85, true},
		{"gte not met", OpGTE, 80, 75, false},
		{"gte equal", OpGTE, 80, 80, true},
		{"lte met", OpLTE, 7, 5, true},
		{"lte not met", OpLTE, 7, 10, false},
		{"eq met", OpEQ, 0, 0, true},
		{"eq not met", OpEQ, 0, 5, false},
		{"gt met", OpGT, 50, 51, true},
		{"gt not met", OpGT, 50, 50, false},
		{"lt met", OpLT, 100, 99, true},
		{"lt not met", OpLT, 100, 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Criterion{
				Operator: tt.operator,
				Target:   tt.target,
			}
			result := c.CheckMet(tt.current)
			if result != tt.expected {
				t.Errorf("CheckMet(%v) = %v, expected %v", tt.current, result, tt.expected)
			}
		})
	}
}

func TestOperatorSymbol(t *testing.T) {
	tests := []struct {
		op       string
		expected string
	}{
		{OpGTE, ">="},
		{OpLTE, "<="},
		{OpGT, ">"},
		{OpLT, "<"},
		{OpEQ, "="},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.op, func(t *testing.T) {
			result := OperatorSymbol(tt.op)
			if result != tt.expected {
				t.Errorf("OperatorSymbol(%q) = %q, expected %q", tt.op, result, tt.expected)
			}
		})
	}
}

func TestLevelProgress(t *testing.T) {
	level := Level{
		Level: 3,
		Criteria: []Criterion{
			{ID: "c1", Operator: OpGTE, Target: 80, Required: true},
			{ID: "c2", Operator: OpEQ, Target: 0, Required: true},
			{ID: "c3", Operator: OpGTE, Target: 100, Required: true},
		},
		Enablers: []Enabler{
			{ID: "e1"},
			{ID: "e2"},
		},
	}

	values := map[string]float64{
		"c1": 85, // Met
		"c2": 0,  // Met
		"c3": 80, // Not met
	}

	enablerStatus := map[string]string{
		"e1": StatusCompleted,
		"e2": StatusInProgress,
	}

	progress := level.CalculateLevelProgress(values, enablerStatus)

	if progress.CriteriaMet != 2 {
		t.Errorf("Expected 2 criteria met, got %d", progress.CriteriaMet)
	}

	if progress.CriteriaTotal != 3 {
		t.Errorf("Expected 3 total criteria, got %d", progress.CriteriaTotal)
	}

	if progress.EnablersDone != 1 {
		t.Errorf("Expected 1 enabler done, got %d", progress.EnablersDone)
	}

	if progress.ProgressPercent < 66 || progress.ProgressPercent > 67 {
		t.Errorf("Expected ~66.67%% progress, got %.2f%%", progress.ProgressPercent)
	}
}

func TestXLSXGenerator_NilAssessments(t *testing.T) {
	// Test that XLSX generation works when assessments are nil
	spec := &Spec{
		Domains: map[string]*DomainModel{
			"Security": {
				Name: "Security",
				Levels: []Level{
					{
						Level:       1,
						Name:        "Reactive",
						Description: "Basic security",
						Criteria: []Criterion{
							{
								ID:         "SEC-001",
								Name:       "Asset Inventory",
								MetricName: "asset_coverage",
								Operator:   "gte",
								Target:     80,
							},
						},
					},
				},
			},
		},
		Assessments: nil, // Explicitly nil
	}

	gen := NewXLSXGenerator(spec)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed with nil assessments: %v", err)
	}
}

func TestXLSXGenerator_EmptyAssessments(t *testing.T) {
	// Test that XLSX generation works when assessments map is empty
	spec := &Spec{
		Domains: map[string]*DomainModel{
			"Operations": {
				Name: "Operations",
				Levels: []Level{
					{
						Level:       1,
						Name:        "Reactive",
						Description: "Basic operations",
					},
				},
			},
		},
		Assessments: map[string]*DomainAssessment{}, // Empty map
	}

	gen := NewXLSXGenerator(spec)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed with empty assessments: %v", err)
	}
}

func TestXLSXGenerator_WithAssessments(t *testing.T) {
	// Test that XLSX generation works with full assessments
	spec := &Spec{
		Domains: map[string]*DomainModel{
			"Security": {
				Name: "Security",
				Levels: []Level{
					{
						Level:       1,
						Name:        "Reactive",
						Description: "Basic security",
						Criteria: []Criterion{
							{
								ID:         "SEC-001",
								Name:       "Asset Inventory",
								MetricName: "asset_coverage",
								Operator:   "gte",
								Target:     80,
							},
						},
						Enablers: []Enabler{
							{
								ID:   "SEC-E-001",
								Name: "Deploy asset scanner",
							},
						},
					},
				},
			},
		},
		Assessments: map[string]*DomainAssessment{
			"Security": {
				Domain:       "Security",
				CurrentLevel: 1,
				TargetLevel:  3,
				CriteriaValues: map[string]float64{
					"SEC-001": 85.0,
				},
				EnablerStatus: map[string]string{
					"SEC-E-001": "completed",
				},
			},
		},
	}

	gen := NewXLSXGenerator(spec)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed with assessments: %v", err)
	}
}

func TestXLSXGenerator_FrameworkColumns(t *testing.T) {
	// Test that framework columns are collected correctly
	spec := &Spec{
		Domains: map[string]*DomainModel{
			"Security": {
				Name: "Security",
				Levels: []Level{
					{
						Level:       2,
						Name:        "Managed",
						Description: "Managed security",
						Criteria: []Criterion{
							{
								ID:         "SEC-001",
								Name:       "Asset Inventory",
								MetricName: "asset_coverage",
								Operator:   "gte",
								Target:     80,
								FrameworkMappings: []FrameworkMapping{
									{Framework: "NIST_CSF_2", Reference: "ID.AM-1"},
									{Framework: "NIST_800_53", Reference: "CM-8"},
								},
							},
							{
								ID:         "SEC-002",
								Name:       "Vulnerability Scanning",
								MetricName: "vuln_scan_coverage",
								Operator:   "gte",
								Target:     90,
								FrameworkMappings: []FrameworkMapping{
									{Framework: "NIST_CSF_2", Reference: "ID.RA-1"},
									{Framework: "SOC_2", Reference: "CC7.1"},
								},
							},
						},
					},
				},
			},
		},
	}

	gen := NewXLSXGenerator(spec)

	// Test collectAllFrameworks
	frameworks := gen.collectAllFrameworks()

	if len(frameworks) != 3 {
		t.Errorf("collectAllFrameworks() returned %d frameworks, want 3", len(frameworks))
	}

	// Should be sorted alphabetically
	expected := []string{"NIST_800_53", "NIST_CSF_2", "SOC_2"}
	for i, fw := range expected {
		if i >= len(frameworks) || frameworks[i] != fw {
			t.Errorf("frameworks[%d] = %q, want %q", i, frameworks[i], fw)
		}
	}
}

func TestXLSXGenerator_QualitativeCriteria(t *testing.T) {
	// Test handling of qualitative criteria
	spec := &Spec{
		Domains: map[string]*DomainModel{
			"Security": {
				Name: "Security",
				Levels: []Level{
					{
						Level:       1,
						Name:        "Reactive",
						Description: "Basic security",
						Criteria: []Criterion{
							{
								ID:         "SEC-Q-001",
								Name:       "Security Policy",
								MetricName: "security_policy",
								Type:       "qualitative",
								Operator:   "exists",
								Target:     0,
								Status:     "documented",
							},
						},
					},
				},
			},
		},
		Assessments: map[string]*DomainAssessment{
			"Security": {
				Domain:       "Security",
				CurrentLevel: 1,
				TargetLevel:  2,
				CriteriaStatus: map[string]string{
					"SEC-Q-001": "implemented",
				},
			},
		},
	}

	gen := NewXLSXGenerator(spec)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed with qualitative criteria: %v", err)
	}
}

func TestXLSXGenerator_FrameworkMappingsSheet(t *testing.T) {
	// Test that framework mappings sheet is created correctly
	spec := &Spec{
		Domains: map[string]*DomainModel{
			"Security": {
				Name: "Security",
				Levels: []Level{
					{
						Level:       2,
						Name:        "Managed",
						Description: "Managed security",
						Criteria: []Criterion{
							{
								ID:         "SEC-001",
								Name:       "Asset Inventory",
								MetricName: "asset_coverage",
								Operator:   "gte",
								Target:     80,
								FrameworkMappings: []FrameworkMapping{
									{
										Framework: "NIST_CSF_2",
										Reference: "ID.AM-1",
										Name:      "Asset Management",
										Baseline:  "Low",
									},
									{
										Framework: "FEDRAMP_HIGH",
										Reference: "CM-8",
										Name:      "System Component Inventory",
										Baseline:  "High",
										Version:   "Rev 5",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	gen := NewXLSXGenerator(spec)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Verify sheets exist
	sheets := gen.file.GetSheetList()
	expectedSheets := []string{"Requirements", "SLOs", "Framework Mappings", "Progress", "Level Definitions"}

	for _, expected := range expectedSheets {
		found := false
		for _, sheet := range sheets {
			if sheet == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected sheet %q not found in %v", expected, sheets)
		}
	}
}

func TestXLSXGenerator_MultipleDomains(t *testing.T) {
	// Test with multiple domains where some have assessments and some don't
	spec := &Spec{
		Domains: map[string]*DomainModel{
			"Security": {
				Name: "Security",
				Levels: []Level{
					{Level: 1, Name: "Reactive", Description: "Basic"},
				},
			},
			"Operations": {
				Name: "Operations",
				Levels: []Level{
					{Level: 1, Name: "Reactive", Description: "Basic"},
				},
			},
			"Quality": {
				Name: "Quality",
				Levels: []Level{
					{Level: 1, Name: "Reactive", Description: "Basic"},
				},
			},
		},
		Assessments: map[string]*DomainAssessment{
			"Security": {
				Domain:       "Security",
				CurrentLevel: 2,
				TargetLevel:  4,
			},
			// Operations and Quality have no assessments
		},
	}

	gen := NewXLSXGenerator(spec)
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed with multiple domains: %v", err)
	}
}

func TestXLSXGenerator_SLIResolution(t *testing.T) {
	// Test that framework mappings are resolved from SLI
	spec := &Spec{
		SLIs: map[string]*SLI{
			"sli-asset-coverage": {
				ID:         "sli-asset-coverage",
				Name:       "Asset Coverage",
				MetricName: "Percentage of assets tracked",
				Unit:       "%",
				Type:       CriterionTypeQuantitative,
				Layer:      "requirements",
				Category:   "prevention",
				FrameworkMappings: []FrameworkMapping{
					{Framework: "NIST_CSF_2", Reference: "ID.AM-1", Name: "Asset Management"},
					{Framework: "NIST_800_53", Reference: "CM-8", Name: "System Component Inventory"},
				},
			},
			"sli-vuln-scan": {
				ID:         "sli-vuln-scan",
				Name:       "Vulnerability Scanning",
				MetricName: "Vulnerability scan coverage",
				Unit:       "%",
				Type:       CriterionTypeQuantitative,
				FrameworkMappings: []FrameworkMapping{
					{Framework: "NIST_CSF_2", Reference: "ID.RA-1"},
					{Framework: "SOC_2", Reference: "CC7.1"},
				},
			},
		},
		Domains: map[string]*DomainModel{
			"Security": {
				Name: "Security",
				Levels: []Level{
					{
						Level:       2,
						Name:        "Basic",
						Description: "Basic security",
						Criteria: []Criterion{
							{
								ID:       "SEC-001",
								Name:     "Asset Inventory M2",
								SLIID:    "sli-asset-coverage", // Reference to SLI
								Operator: "gte",
								Target:   80,
							},
							{
								ID:       "SEC-002",
								Name:     "Vuln Scan M2",
								SLIID:    "sli-vuln-scan", // Reference to SLI
								Operator: "gte",
								Target:   70,
							},
						},
					},
					{
						Level:       3,
						Name:        "Defined",
						Description: "Defined security",
						Criteria: []Criterion{
							{
								ID:       "SEC-003",
								Name:     "Asset Inventory M3",
								SLIID:    "sli-asset-coverage", // Same SLI, different target
								Operator: "gte",
								Target:   90,
							},
							{
								ID:       "SEC-004",
								Name:     "Vuln Scan M3",
								SLIID:    "sli-vuln-scan", // Same SLI, different target
								Operator: "gte",
								Target:   85,
							},
						},
					},
				},
			},
		},
	}

	gen := NewXLSXGenerator(spec)

	// Test collectAllFrameworks resolves from SLI
	frameworks := gen.collectAllFrameworks()
	if len(frameworks) != 3 {
		t.Errorf("collectAllFrameworks() returned %d frameworks, want 3", len(frameworks))
	}

	expected := []string{"NIST_800_53", "NIST_CSF_2", "SOC_2"}
	for i, fw := range expected {
		if i >= len(frameworks) || frameworks[i] != fw {
			t.Errorf("frameworks[%d] = %q, want %q", i, frameworks[i], fw)
		}
	}

	// Test that generation succeeds
	err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() failed with SLI references: %v", err)
	}

	// Verify sheets exist
	sheets := gen.file.GetSheetList()
	expectedSheets := []string{"Requirements", "SLOs", "Framework Mappings", "Progress", "Level Definitions"}

	for _, expected := range expectedSheets {
		found := false
		for _, sheet := range sheets {
			if sheet == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected sheet %q not found in %v", expected, sheets)
		}
	}
}

func TestCriterion_GetFrameworkMappings(t *testing.T) {
	// Test GetFrameworkMappings resolution logic
	sli := &SLI{
		ID:   "test-sli",
		Name: "Test SLI",
		FrameworkMappings: []FrameworkMapping{
			{Framework: "NIST_CSF_2", Reference: "ID.AM-1"},
		},
	}

	spec := &Spec{
		SLIs: map[string]*SLI{
			"test-sli": sli,
		},
	}

	t.Run("inline mappings take precedence", func(t *testing.T) {
		c := Criterion{
			ID:    "c1",
			SLIID: "test-sli",
			FrameworkMappings: []FrameworkMapping{
				{Framework: "SOC_2", Reference: "CC1.1"},
			},
		}

		mappings := c.GetFrameworkMappings(spec)
		if len(mappings) != 1 {
			t.Fatalf("expected 1 mapping, got %d", len(mappings))
		}
		if mappings[0].Framework != "SOC_2" {
			t.Errorf("expected SOC_2, got %s", mappings[0].Framework)
		}
	})

	t.Run("resolve from SLI when no inline mappings", func(t *testing.T) {
		c := Criterion{
			ID:    "c2",
			SLIID: "test-sli",
			// No inline FrameworkMappings
		}

		mappings := c.GetFrameworkMappings(spec)
		if len(mappings) != 1 {
			t.Fatalf("expected 1 mapping from SLI, got %d", len(mappings))
		}
		if mappings[0].Framework != "NIST_CSF_2" {
			t.Errorf("expected NIST_CSF_2, got %s", mappings[0].Framework)
		}
	})

	t.Run("nil when no SLI and no inline mappings", func(t *testing.T) {
		c := Criterion{
			ID: "c3",
			// No SLIID, no inline FrameworkMappings
		}

		mappings := c.GetFrameworkMappings(spec)
		if mappings != nil {
			t.Errorf("expected nil, got %v", mappings)
		}
	})
}

func TestCriterion_GetMetricName(t *testing.T) {
	sli := &SLI{
		ID:         "test-sli",
		MetricName: "SLI Metric Name",
		Unit:       "%",
		Type:       CriterionTypeQualitative,
		Layer:      "code",
		Category:   "detection",
	}

	spec := &Spec{
		SLIs: map[string]*SLI{
			"test-sli": sli,
		},
	}

	t.Run("inline metric name takes precedence", func(t *testing.T) {
		c := Criterion{
			ID:         "c1",
			SLIID:      "test-sli",
			MetricName: "Inline Metric",
		}

		if got := c.GetMetricName(spec); got != "Inline Metric" {
			t.Errorf("expected 'Inline Metric', got %q", got)
		}
	})

	t.Run("resolve from SLI", func(t *testing.T) {
		c := Criterion{
			ID:    "c2",
			SLIID: "test-sli",
		}

		if got := c.GetMetricName(spec); got != "SLI Metric Name" {
			t.Errorf("expected 'SLI Metric Name', got %q", got)
		}
	})

	t.Run("resolve unit from SLI", func(t *testing.T) {
		c := Criterion{
			ID:    "c2",
			SLIID: "test-sli",
		}

		if got := c.GetUnit(spec); got != "%" {
			t.Errorf("expected '%%', got %q", got)
		}
	})

	t.Run("resolve type from SLI", func(t *testing.T) {
		c := Criterion{
			ID:    "c2",
			SLIID: "test-sli",
		}

		if got := c.GetType(spec); got != CriterionTypeQualitative {
			t.Errorf("expected 'qualitative', got %q", got)
		}
	})

	t.Run("resolve layer from SLI", func(t *testing.T) {
		c := Criterion{
			ID:    "c2",
			SLIID: "test-sli",
		}

		if got := c.GetLayer(spec); got != "code" {
			t.Errorf("expected 'code', got %q", got)
		}
	})

	t.Run("resolve category from SLI", func(t *testing.T) {
		c := Criterion{
			ID:    "c2",
			SLIID: "test-sli",
		}

		if got := c.GetCategory(spec); got != "detection" {
			t.Errorf("expected 'detection', got %q", got)
		}
	})
}
