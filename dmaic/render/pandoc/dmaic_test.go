package pandoc

import (
	"strings"
	"testing"

	"github.com/grokify/prism-maturity/dmaic"
	"github.com/grokify/prism-maturity/dmaic/render"
)

func TestRenderer_Format(t *testing.T) {
	r := New()
	if r.Format() != "pandoc" {
		t.Errorf("Expected format 'pandoc', got '%s'", r.Format())
	}
}

func TestRenderer_FileExtension(t *testing.T) {
	r := New()
	if r.FileExtension() != ".md" {
		t.Errorf("Expected extension '.md', got '%s'", r.FileExtension())
	}
}

func TestRenderer_Render(t *testing.T) {
	doc := &dmaic.DMAICDocument{
		Metadata: &dmaic.Metadata{
			ID:     "DMAIC-TEST",
			Name:   "Test DMAIC Report",
			Owner:  "Test Owner",
			Team:   "Test Team",
			Period: "2025-Q1",
			Status: dmaic.DocumentStatusActive,
		},
		Categories: []dmaic.Category{
			{
				Name:        "Test Category",
				Description: "A test category",
				Owner:       "Category Owner",
				Metrics: []dmaic.Metric{
					{
						Name:           "Test Metric",
						Baseline:       100,
						Current:        85,
						Target:         90,
						Unit:           "%",
						Phase:          dmaic.PhaseControl,
						TrendDirection: dmaic.TrendHigherBetter,
						Status:         dmaic.StatusGreen,
					},
				},
			},
		},
	}

	r := New()
	output, err := r.Render(doc, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	content := string(output)

	// Check YAML frontmatter
	if !strings.Contains(content, "---") {
		t.Error("Expected YAML frontmatter delimiter")
	}
	if !strings.Contains(content, "title:") {
		t.Error("Expected title in frontmatter")
	}
	if !strings.Contains(content, "geometry:") {
		t.Error("Expected geometry in frontmatter")
	}
	if !strings.Contains(content, "margin=2cm") {
		t.Error("Expected 2cm margin in frontmatter")
	}
	if !strings.Contains(content, "fontfamily: helvet") {
		t.Error("Expected sans-serif font family")
	}
	if !strings.Contains(content, "\\renewcommand{\\familydefault}{\\sfdefault}") {
		t.Error("Expected sans-serif default font")
	}

	// Check content
	if !strings.Contains(content, "Test DMAIC Report") {
		t.Error("Expected document name in output")
	}
	if !strings.Contains(content, "Test Category") {
		t.Error("Expected category name in output")
	}
	if !strings.Contains(content, "Test Metric") {
		t.Error("Expected metric name in output")
	}
}

func TestRenderer_RenderWithOptions(t *testing.T) {
	doc := &dmaic.DMAICDocument{
		Metadata: &dmaic.Metadata{
			Name: "Options Test",
		},
		Categories: []dmaic.Category{
			{
				Name: "Category 1",
				Metrics: []dmaic.Metric{
					{
						Name:    "Metric 1",
						Phase:   dmaic.PhaseDefine,
						Current: 50,
						Target:  100,
					},
				},
			},
		},
		Initiatives: []dmaic.Initiative{
			{
				Name:           "Test Initiative",
				Status:         dmaic.InitiativeStatusInProgress,
				ExpectedImpact: "50% improvement",
			},
		},
	}

	r := New()

	// Test with initiatives enabled
	opts := &render.Options{
		IncludeInitiatives: true,
		GroupByPhase:       true,
	}
	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render with options failed: %v", err)
	}

	content := string(output)
	if !strings.Contains(content, "Test Initiative") {
		t.Error("Expected initiative in output")
	}
	if !strings.Contains(content, "DMAIC Phase Distribution") {
		t.Error("Expected phase summary in output")
	}
}

func TestRenderer_RenderWithRootCauses(t *testing.T) {
	doc := &dmaic.DMAICDocument{
		Categories: []dmaic.Category{
			{
				Name: "Security",
				Metrics: []dmaic.Metric{
					{
						Name:   "Vulnerabilities",
						Phase:  dmaic.PhaseAnalyze,
						Status: dmaic.StatusYellow,
						RootCauses: []dmaic.RootCause{
							{
								Description: "Insufficient testing",
								Category:    "Process",
								Impact:      "High",
								Validated:   true,
							},
						},
					},
				},
			},
		},
	}

	r := New()
	opts := &render.Options{
		IncludeRootCauses: true,
	}
	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render with root causes failed: %v", err)
	}

	content := string(output)
	if !strings.Contains(content, "Root Cause Analysis") {
		t.Error("Expected root cause analysis section")
	}
	if !strings.Contains(content, "Insufficient testing") {
		t.Error("Expected root cause description")
	}
}

func TestRenderer_RenderWithCapabilityMetrics(t *testing.T) {
	doc := &dmaic.DMAICDocument{
		Categories: []dmaic.Category{
			{
				Name: "Quality",
				Metrics: []dmaic.Metric{
					{
						Name:  "Defect Rate",
						Phase: dmaic.PhaseControl,
						ProcessCapability: &dmaic.ProcessCapability{
							Cp:         1.33,
							Cpk:        1.25,
							SigmaLevel: 4.0,
							DPMO:       6210,
						},
						ControlLimits: &dmaic.ControlLimits{
							UCL:        10,
							LCL:        2,
							CenterLine: 5,
						},
					},
				},
			},
		},
	}

	r := New()
	opts := &render.Options{
		ShowCapabilityMetrics: true,
	}
	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render with capability metrics failed: %v", err)
	}

	content := string(output)
	if !strings.Contains(content, "Process Capability Metrics") {
		t.Error("Expected process capability appendix")
	}
	if !strings.Contains(content, "Control Limits Summary") {
		t.Error("Expected control limits summary")
	}
}

func TestRenderer_EscapeLatex(t *testing.T) {
	doc := &dmaic.DMAICDocument{
		Metadata: &dmaic.Metadata{
			Name:  "Test & Report with $pecial #chars",
			Owner: "User_Name",
		},
		Categories: []dmaic.Category{
			{
				Name: "Category with % and &",
				Metrics: []dmaic.Metric{
					{Name: "Metric_1", Phase: dmaic.PhaseDefine},
				},
			},
		},
	}

	r := New()
	output, err := r.Render(doc, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	content := string(output)
	// Check that special characters are escaped
	if !strings.Contains(content, "\\&") {
		t.Error("Expected ampersand to be escaped")
	}
	if !strings.Contains(content, "\\%") {
		t.Error("Expected percent to be escaped")
	}
	if !strings.Contains(content, "\\_") {
		t.Error("Expected underscore to be escaped")
	}
}

func TestRenderer_StatusText(t *testing.T) {
	tests := []struct {
		status   string
		expected string
	}{
		{dmaic.StatusGreen, "Green"},
		{dmaic.StatusYellow, "Yellow"},
		{dmaic.StatusRed, "Red"},
		{"", "Unknown"},
	}

	for _, tt := range tests {
		result := funcMap["statusText"].(func(string) string)(tt.status)
		if result != tt.expected {
			t.Errorf("statusText(%q) = %q, expected %q", tt.status, result, tt.expected)
		}
	}
}

func TestRenderer_PhaseAbbrev(t *testing.T) {
	tests := []struct {
		phase    string
		expected string
	}{
		{dmaic.PhaseDefine, "D"},
		{dmaic.PhaseMeasure, "M"},
		{dmaic.PhaseAnalyze, "A"},
		{dmaic.PhaseImprove, "I"},
		{dmaic.PhaseControl, "C"},
		{"", "?"},
	}

	for _, tt := range tests {
		result := funcMap["phaseAbbrev"].(func(string) string)(tt.phase)
		if result != tt.expected {
			t.Errorf("phaseAbbrev(%q) = %q, expected %q", tt.phase, result, tt.expected)
		}
	}
}

func TestRenderer_Truncate(t *testing.T) {
	truncate := funcMap["truncate"].(func(string, int) string)

	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly ten", 11, "exactly ten"},
		{"this is a long string", 10, "this is..."},
	}

	for _, tt := range tests {
		result := truncate(tt.input, tt.max)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, expected %q", tt.input, tt.max, result, tt.expected)
		}
	}
}
