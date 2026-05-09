package maturity

import (
	"strings"
	"testing"
)

func TestGenerateMarkdown(t *testing.T) {
	spec := &Spec{
		Metadata: &SpecMetadata{
			Name:        "Test Maturity Model",
			Description: "A test maturity model for unit testing",
		},
		Domains: map[string]*DomainModel{
			"security": {
				Name:        "Security",
				Description: "Security domain",
				Levels: []Level{
					{
						Level:       1,
						Name:        "Reactive",
						Description: "Ad-hoc security",
					},
					{
						Level:       2,
						Name:        "Basic",
						Description: "Basic security controls",
						Criteria: []Criterion{
							{
								ID:         "sec-sast",
								Name:       "SAST Coverage",
								MetricName: "Percentage of repos with SAST",
								Operator:   OpGTE,
								Target:     50,
								Unit:       "%",
								FrameworkMappings: []FrameworkMapping{
									{
										Framework: "NIST_800_53",
										Reference: "SA-11",
										Name:      "Developer Testing",
									},
								},
							},
							{
								ID:         "sec-enc",
								Name:       "Encryption",
								Type:       CriterionTypeQualitative,
								Operator:   OperatorExists,
								Status:     QualStatusImplemented,
								MetricName: "Encryption at rest enabled",
								FrameworkMappings: []FrameworkMapping{
									{
										Framework: "NIST_CSF_2",
										Reference: "PR.DS-01",
										Name:      "Data-at-rest is protected",
									},
								},
							},
						},
						Enablers: []Enabler{
							{
								ID:     "e1",
								Name:   "Deploy SAST",
								Type:   TypeTooling,
								Status: StatusCompleted,
							},
						},
					},
				},
			},
		},
		Assessments: map[string]*DomainAssessment{
			"security": {
				Domain:       "security",
				CurrentLevel: 2,
				TargetLevel:  4,
				CriteriaValues: map[string]float64{
					"sec-sast": 60,
				},
				CriteriaStatus: map[string]string{
					"sec-enc": QualStatusImplemented,
				},
			},
		},
	}

	t.Run("default options", func(t *testing.T) {
		md := spec.GenerateMarkdown(nil)

		// Check YAML front matter
		if !strings.Contains(md, "---\n") {
			t.Error("Expected YAML front matter")
		}
		// Default options use "Maturity Model" as title
		if !strings.Contains(md, "title: \"Maturity Model\"") {
			t.Error("Expected default title in front matter")
		}

		// Check TOC
		if !strings.Contains(md, "## Table of Contents") {
			t.Error("Expected table of contents")
		}

		// Check domain view
		if !strings.Contains(md, "# Maturity Model by Domain") {
			t.Error("Expected domain view header")
		}
		if !strings.Contains(md, "## Security") {
			t.Error("Expected Security domain")
		}

		// Check framework view
		if !strings.Contains(md, "# Maturity Model by Framework") {
			t.Error("Expected framework view header")
		}
		if !strings.Contains(md, "NIST SP 800-53") {
			t.Error("Expected NIST 800-53 framework")
		}
		if !strings.Contains(md, "NIST CSF 2.0") {
			t.Error("Expected NIST CSF 2.0 framework")
		}
	})

	t.Run("domain view only", func(t *testing.T) {
		opts := &MarkdownOptions{
			ViewType:        "domain",
			IncludeYAMLMeta: false,
			IncludeTOC:      false,
		}
		md := spec.GenerateMarkdown(opts)

		if strings.Contains(md, "---\n") {
			t.Error("Should not have YAML front matter")
		}
		if !strings.Contains(md, "# Maturity Model by Domain") {
			t.Error("Expected domain view")
		}
		if strings.Contains(md, "# Maturity Model by Framework") {
			t.Error("Should not have framework view")
		}
	})

	t.Run("framework view only", func(t *testing.T) {
		opts := &MarkdownOptions{
			ViewType:        "framework",
			IncludeYAMLMeta: false,
		}
		md := spec.GenerateMarkdown(opts)

		if strings.Contains(md, "# Maturity Model by Domain") {
			t.Error("Should not have domain view")
		}
		if !strings.Contains(md, "# Maturity Model by Framework") {
			t.Error("Expected framework view")
		}
	})

	t.Run("filter frameworks", func(t *testing.T) {
		opts := &MarkdownOptions{
			ViewType:        "framework",
			IncludeYAMLMeta: false,
			Frameworks:      []string{"NIST_800_53"},
		}
		md := spec.GenerateMarkdown(opts)

		if !strings.Contains(md, "NIST SP 800-53") {
			t.Error("Expected NIST 800-53 framework")
		}
		// NIST CSF 2.0 should be filtered out
		if strings.Contains(md, "## NIST CSF 2.0") {
			t.Error("Should not have NIST CSF 2.0 framework")
		}
	})

	t.Run("criteria status indicators", func(t *testing.T) {
		md := spec.GenerateMarkdown(nil)

		// SAST coverage should be met (60 >= 50)
		if !strings.Contains(md, "| sec-sast | SAST Coverage | Quantitative | >=50% | ✅ |") {
			t.Error("Expected SAST criterion with met status")
		}

		// Encryption should be met (implemented)
		if !strings.Contains(md, "| sec-enc | Encryption | Qualitative | Tracked | ✅ |") {
			t.Error("Expected Encryption criterion with met status")
		}
	})
}

func TestFormatFrameworkName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"NIST_CSF", "NIST CSF 1.1"},
		{"NIST_CSF_2", "NIST CSF 2.0"},
		{"NIST_800_53", "NIST SP 800-53"},
		{"FEDRAMP_HIGH", "FedRAMP High"},
		{"FEDRAMP_MOD", "FedRAMP Moderate"},
		{"DORA", "DORA"},
		{"SRE", "SRE"},
		{"UNKNOWN", "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatFrameworkName(tt.input)
			if result != tt.expected {
				t.Errorf("formatFrameworkName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsQualitativeStatusMet(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{QualStatusTracked, true},
		{QualStatusImplemented, true},
		{QualStatusDefined, true},
		{QualStatusDocumented, true},
		{QualStatusCompliant, true},
		{QualStatusEnabled, true},
		{QualStatusNotTracked, false},
		{QualStatusPartial, false},
		{QualStatusPlanned, false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := IsQualitativeStatusMet(tt.status)
			if result != tt.expected {
				t.Errorf("IsQualitativeStatusMet(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestSLICatalog(t *testing.T) {
	spec := &Spec{
		SLIs: map[string]*SLI{
			"sli-asset-coverage": {
				ID:         "sli-asset-coverage",
				Name:       "Asset Coverage",
				MetricName: "asset_coverage_pct",
				Unit:       "%",
				Type:       CriterionTypeQuantitative,
				Category:   "prevention",
				FrameworkMappings: []FrameworkMapping{
					{Framework: "NIST_CSF_2", Reference: "ID.AM-1"},
					{Framework: "NIST_800_53", Reference: "CM-8"},
				},
			},
			"sli-mttr": {
				ID:         "sli-mttr",
				Name:       "Security MTTR",
				MetricName: "security_mttr_days",
				Unit:       "days",
				Type:       CriterionTypeQuantitative,
				Category:   "response",
				FrameworkMappings: []FrameworkMapping{
					{Framework: "NIST_800_53", Reference: "IR-6"},
				},
			},
			"sli-encryption": {
				ID:         "sli-encryption",
				Name:       "Encryption Status",
				MetricName: "encryption_enabled",
				Type:       CriterionTypeQualitative,
				Category:   "prevention",
				FrameworkMappings: []FrameworkMapping{
					{Framework: "NIST_CSF_2", Reference: "PR.DS-01"},
				},
			},
		},
		Domains: map[string]*DomainModel{
			"security": {
				Name:   "Security",
				Levels: []Level{{Level: 1, Name: "Reactive", Description: "Ad-hoc"}},
			},
		},
	}

	md := spec.GenerateMarkdown(nil)

	t.Run("has SLI Catalog section", func(t *testing.T) {
		if !strings.Contains(md, "## SLI Catalog") {
			t.Error("Expected SLI Catalog section")
		}
	})

	t.Run("groups by category", func(t *testing.T) {
		if !strings.Contains(md, "### Prevention") {
			t.Error("Expected Prevention category")
		}
		if !strings.Contains(md, "### Response") {
			t.Error("Expected Response category")
		}
	})

	t.Run("includes SLI details", func(t *testing.T) {
		if !strings.Contains(md, "Asset Coverage") {
			t.Error("Expected Asset Coverage SLI")
		}
		if !strings.Contains(md, "asset_coverage_pct") {
			t.Error("Expected metric name")
		}
		if !strings.Contains(md, "NIST_CSF_2:ID.AM-1") {
			t.Error("Expected framework mapping")
		}
	})

	t.Run("shows qualitative type", func(t *testing.T) {
		if !strings.Contains(md, "| Encryption Status | encryption_enabled | - | Qualitative |") {
			t.Error("Expected Encryption Status with Qualitative type")
		}
	})

	t.Run("TOC includes SLI Catalog link", func(t *testing.T) {
		if !strings.Contains(md, "- [SLI Catalog](#sli-catalog)") {
			t.Error("Expected SLI Catalog link in TOC")
		}
	})
}

func TestSLICatalog_NoSLIs(t *testing.T) {
	spec := &Spec{
		Domains: map[string]*DomainModel{
			"security": {
				Name:   "Security",
				Levels: []Level{{Level: 1, Name: "Reactive", Description: "Ad-hoc"}},
			},
		},
	}

	md := spec.GenerateMarkdown(nil)

	if strings.Contains(md, "## SLI Catalog") {
		t.Error("Should not have SLI Catalog when no SLIs defined")
	}
}
