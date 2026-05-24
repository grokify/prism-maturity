package cli

import (
	"testing"
)

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a longer string", 10, "this is..."},
		{"", 10, ""},
		{"abc", 3, "abc"},
		{"abcd", 3, "abc"}, // maxLen <= 3: just truncate without ellipsis
	}

	for _, tc := range tests {
		result := truncateString(tc.input, tc.maxLen)
		if result != tc.expected {
			t.Errorf("truncateString(%q, %d) = %q, want %q",
				tc.input, tc.maxLen, result, tc.expected)
		}
	}
}

func TestGetGoalStatus(t *testing.T) {
	tests := []struct {
		current  int
		target   int
		expected string
	}{
		{5, 5, "Achieved"},
		{6, 5, "Achieved"},
		{4, 5, "On Track"},
		{3, 5, "Behind"},
		{1, 5, "Behind"},
		{0, 1, "On Track"}, // 0 is only 1 below target 1, so "On Track"
		{1, 1, "Achieved"},
		{0, 2, "Behind"}, // 0 is 2 below target 2
	}

	for _, tc := range tests {
		result := getGoalStatus(tc.current, tc.target)
		if result != tc.expected {
			t.Errorf("getGoalStatus(%d, %d) = %q, want %q",
				tc.current, tc.target, result, tc.expected)
		}
	}
}

func TestOperatorSymbol(t *testing.T) {
	tests := []struct {
		op       string
		expected string
	}{
		{"gte", ">="},
		{"lte", "<="},
		{"gt", ">"},
		{"lt", "<"},
		{"eq", "="},
		{"unknown", "unknown"},
		{"", ""},
	}

	for _, tc := range tests {
		result := operatorSymbol(tc.op)
		if result != tc.expected {
			t.Errorf("operatorSymbol(%q) = %q, want %q",
				tc.op, result, tc.expected)
		}
	}
}

func TestSafePercent(t *testing.T) {
	tests := []struct {
		num      int
		denom    int
		expected float64
	}{
		{50, 100, 50.0},
		{100, 100, 100.0},
		{0, 100, 0.0},
		{0, 0, 0.0},
		{1, 4, 25.0},
		{3, 4, 75.0},
	}

	for _, tc := range tests {
		result := safePercent(tc.num, tc.denom)
		if result != tc.expected {
			t.Errorf("safePercent(%d, %d) = %f, want %f",
				tc.num, tc.denom, result, tc.expected)
		}
	}
}
