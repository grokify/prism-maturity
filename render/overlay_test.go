package render

import (
	"testing"
)

func TestFormatMaturityBadge(t *testing.T) {
	tests := []struct {
		name     string
		level    float64
		expected string
	}{
		{"integer level 1", 1.0, "M1"},
		{"integer level 3", 3.0, "M3"},
		{"integer level 5", 5.0, "M5"},
		{"decimal level", 2.5, "M2.5"},
		{"decimal level high", 4.7, "M4.7"},
		{"decimal level low", 1.2, "M1.2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatMaturityBadge(tt.level)
			if result != tt.expected {
				t.Errorf("formatMaturityBadge(%v) = %q, want %q", tt.level, result, tt.expected)
			}
		})
	}
}

func TestMaturityColors(t *testing.T) {
	// Verify all levels 1-5 have colors defined
	for level := 1; level <= 5; level++ {
		if maturityColors[level] == "" {
			t.Errorf("maturityColors[%d] is empty", level)
		}
		if maturityTextColors[level] == "" {
			t.Errorf("maturityTextColors[%d] is empty", level)
		}
	}

	// Verify specific colors
	expectedColors := map[int]string{
		1: "#ef4444", // red
		2: "#f59e0b", // amber
		3: "#eab308", // yellow
		4: "#22c55e", // green
		5: "#3b82f6", // blue
	}

	for level, expected := range expectedColors {
		if maturityColors[level] != expected {
			t.Errorf("maturityColors[%d] = %q, want %q", level, maturityColors[level], expected)
		}
	}
}

func TestBuildMaturityOverlayNil(t *testing.T) {
	result := BuildMaturityOverlay(nil)
	if result != nil {
		t.Errorf("BuildMaturityOverlay(nil) = %v, want nil", result)
	}
}

func TestBuildLayerOverlayNil(t *testing.T) {
	result := BuildLayerOverlay(nil)
	if result != nil {
		t.Errorf("BuildLayerOverlay(nil) = %v, want nil", result)
	}
}
