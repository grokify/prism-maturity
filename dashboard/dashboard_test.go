package dashboard

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/prism/maturity"
)

func TestGenerateDashboard(t *testing.T) {
	// Load the operations maturity model
	specFile := filepath.Join("..", "maturity-models", "operations.json")
	spec, err := maturity.ReadSpecFile(specFile)
	if err != nil {
		t.Fatalf("Failed to read spec file: %v", err)
	}

	gen := NewGenerator(spec)
	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard: %v", err)
	}

	if dashboard.ID == "" {
		t.Error("Dashboard ID is empty")
	}

	if dashboard.Title == "" {
		t.Error("Dashboard title is empty")
	}

	if len(dashboard.Widgets) == 0 {
		t.Error("Dashboard has no widgets")
	}

	if len(dashboard.DataSources) == 0 {
		t.Error("Dashboard has no data sources")
	}

	// Verify we have different widget types
	widgetTypes := make(map[string]int)
	for _, w := range dashboard.Widgets {
		widgetTypes[w.Type]++
	}

	if widgetTypes["metric"] == 0 {
		t.Error("Dashboard has no metric widgets")
	}

	if widgetTypes["chart"] == 0 {
		t.Error("Dashboard has no chart widgets")
	}

	// Export to JSON for inspection
	jsonBytes, err := dashboard.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal dashboard: %v", err)
	}

	// Verify it's valid JSON
	var check map[string]any
	if err := json.Unmarshal(jsonBytes, &check); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	// Save to temp file for inspection
	tmpFile := filepath.Join(os.TempDir(), "prism-dashboard-test.json")
	if err := os.WriteFile(tmpFile, jsonBytes, 0600); err != nil { //nolint:gosec
		t.Logf("Could not write temp file: %v", err)
	} else {
		t.Logf("Generated dashboard: %s (%d bytes)", tmpFile, len(jsonBytes))
	}
}

func TestGenerateSecurityDashboard(t *testing.T) {
	// Load the security maturity model
	specFile := filepath.Join("..", "maturity-models", "security.json")
	spec, err := maturity.ReadSpecFile(specFile)
	if err != nil {
		t.Fatalf("Failed to read spec file: %v", err)
	}

	gen := NewGenerator(spec)
	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard: %v", err)
	}

	// Check SLI tables are generated
	hasSLITable := false
	for _, w := range dashboard.Widgets {
		if w.Type == "table" {
			hasSLITable = true
			break
		}
	}

	if !hasSLITable {
		t.Error("Dashboard has no SLI tables")
	}

	// Export to JSON
	jsonBytes, err := dashboard.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal dashboard: %v", err)
	}

	tmpFile := filepath.Join(os.TempDir(), "prism-security-dashboard.json")
	if err := os.WriteFile(tmpFile, jsonBytes, 0600); err != nil { //nolint:gosec
		t.Logf("Could not write temp file: %v", err)
	} else {
		t.Logf("Generated security dashboard: %s (%d bytes)", tmpFile, len(jsonBytes))
	}
}

func TestEmptySpec(t *testing.T) {
	gen := NewGenerator(&maturity.Spec{
		Domains: map[string]*maturity.DomainModel{},
	})

	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate empty dashboard: %v", err)
	}

	if dashboard == nil {
		t.Error("Dashboard is nil")
	}
}

func TestNilSpec(t *testing.T) {
	gen := NewGenerator(nil)
	_, err := gen.Generate()
	if err == nil {
		t.Error("Expected error for nil spec")
	}
}

func TestGenerateHTML(t *testing.T) {
	specFile := filepath.Join("..", "maturity-models", "operations.json")
	spec, err := maturity.ReadSpecFile(specFile)
	if err != nil {
		t.Fatalf("Failed to read spec file: %v", err)
	}

	gen := NewGenerator(spec)
	dashboard, err := gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate dashboard: %v", err)
	}

	html, err := dashboard.ToHTML(DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("Failed to generate HTML: %v", err)
	}

	if len(html) == 0 {
		t.Error("HTML output is empty")
	}

	// Check for required elements
	if !contains(html, "<!DOCTYPE html>") {
		t.Error("HTML missing DOCTYPE")
	}
	if !contains(html, "echarts") {
		t.Error("HTML missing ECharts")
	}
	if !contains(html, "dashboard") {
		t.Error("HTML missing dashboard data")
	}

	// Save for inspection
	tmpFile := filepath.Join(os.TempDir(), "prism-dashboard.html")
	if err := os.WriteFile(tmpFile, []byte(html), 0600); err != nil { //nolint:gosec
		t.Logf("Could not write temp file: %v", err)
	} else {
		t.Logf("Generated HTML: %s (%d bytes)", tmpFile, len(html))
	}
}

func TestMaturityBullet(t *testing.T) {
	// Test NewMaturityBullet
	bullet := NewMaturityBullet("Availability", "99.5% uptime", 3.5, 5)

	if bullet.Title != "Availability" {
		t.Errorf("Expected title 'Availability', got '%s'", bullet.Title)
	}
	if bullet.Subtitle != "99.5% uptime" {
		t.Errorf("Expected subtitle '99.5%% uptime', got '%s'", bullet.Subtitle)
	}
	if len(bullet.Ranges) != 3 {
		t.Errorf("Expected 3 ranges, got %d", len(bullet.Ranges))
	}
	if len(bullet.Measures) != 1 || bullet.Measures[0] != 3.5 {
		t.Errorf("Expected measures [3.5], got %v", bullet.Measures)
	}
	if len(bullet.Markers) != 1 || bullet.Markers[0] != 5 {
		t.Errorf("Expected markers [5], got %v", bullet.Markers)
	}

	// Test MaturityLevel
	cases := []struct {
		value    float64
		expected string
	}{
		{5.0, "M5"},
		{4.5, "M4"},
		{4.0, "M4"},
		{3.5, "M3"},
		{2.0, "M2"},
		{1.0, "M1"},
		{0.5, "M0"},
	}

	for _, tc := range cases {
		got := MaturityLevel(tc.value)
		if got != tc.expected {
			t.Errorf("MaturityLevel(%v) = %s, want %s", tc.value, got, tc.expected)
		}
	}

	// Test MaturityStatus
	statusCases := []struct {
		value    float64
		expected string
	}{
		{5.0, "green"},
		{4.5, "yellow"},
		{3.0, "red"},
	}

	for _, tc := range statusCases {
		got := MaturityStatus(tc.value)
		if got != tc.expected {
			t.Errorf("MaturityStatus(%v) = %s, want %s", tc.value, got, tc.expected)
		}
	}
}

func TestMaturityBulletCSS(t *testing.T) {
	css := GetMaturityBulletCSS()

	// Check for required CSS classes
	requiredClasses := []string{
		".bullet",
		".range.s0",
		".range.s1",
		".range.s2",
		".measure.s0",
		"#fee2e2", // red
		"#fef3c7", // yellow
		"#dcfce7", // green
	}

	for _, class := range requiredClasses {
		if !contains(css, class) {
			t.Errorf("CSS missing '%s'", class)
		}
	}
}

func TestGenerateBullets(t *testing.T) {
	specFile := filepath.Join("..", "maturity-models", "operations.json")
	spec, err := maturity.ReadSpecFile(specFile)
	if err != nil {
		t.Fatalf("Failed to read spec file: %v", err)
	}

	gen := NewGenerator(spec)
	bulletData := gen.GenerateMaturityBullets()

	if len(bulletData.Bullets) == 0 {
		t.Error("Expected bullets to be generated")
	}

	// Check JSON serialization
	jsonBytes, err := bulletData.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal bullet data: %v", err)
	}

	if len(jsonBytes) == 0 {
		t.Error("JSON output is empty")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(len(s) >= len(substr) && (s == substr ||
			len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
