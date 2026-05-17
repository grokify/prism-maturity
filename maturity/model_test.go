package maturity

import (
	"sort"
	"testing"
)

func TestValidateTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		wantErr bool
	}{
		// Valid tags
		{"simple lowercase", "ai", false},
		{"with hyphen", "shift-left", false},
		{"multiple hyphens", "supply-chain-security", false},
		{"with numbers", "cloud9", false},
		{"hyphen and numbers", "v2-api", false},
		{"max length", "abcdefghijklmnopqrstuvwxyz123456", false}, // 32 chars

		// Invalid tags
		{"empty", "", true},
		{"too long", "abcdefghijklmnopqrstuvwxyz1234567", true}, // 33 chars
		{"uppercase", "AI", true},
		{"mixed case", "ShiftLeft", true},
		{"starts with number", "123abc", true},
		{"starts with hyphen", "-invalid", true},
		{"ends with hyphen", "invalid-", true},
		{"consecutive hyphens", "shift--left", true},
		{"contains underscore", "shift_left", true},
		{"contains space", "shift left", true},
		{"contains special char", "shift@left", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTag(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTag(%q) error = %v, wantErr %v", tt.tag, err, tt.wantErr)
			}
		})
	}
}

func TestNormalizeTags(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		want []string
	}{
		{
			name: "empty slice",
			tags: []string{},
			want: nil,
		},
		{
			name: "nil slice",
			tags: nil,
			want: nil,
		},
		{
			name: "single tag",
			tags: []string{"ai"},
			want: []string{"ai"},
		},
		{
			name: "sorted output",
			tags: []string{"supply-chain", "ai", "shift-left"},
			want: []string{"ai", "shift-left", "supply-chain"},
		},
		{
			name: "deduplication",
			tags: []string{"ai", "shift-left", "ai", "supply-chain", "shift-left"},
			want: []string{"ai", "shift-left", "supply-chain"},
		},
		{
			name: "trims whitespace",
			tags: []string{" ai ", "  shift-left  "},
			want: []string{"ai", "shift-left"},
		},
		{
			name: "lowercases input",
			tags: []string{"AI", "Shift-Left", "SUPPLY-CHAIN"},
			want: []string{"ai", "shift-left", "supply-chain"},
		},
		{
			name: "filters invalid tags",
			tags: []string{"ai", "123invalid", "shift-left", ""},
			want: []string{"ai", "shift-left"},
		},
		{
			name: "all invalid returns nil",
			tags: []string{"", "  ", "123", "-invalid"},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeTags(tt.tags)
			if !slicesEqual(got, tt.want) {
				t.Errorf("NormalizeTags(%v) = %v, want %v", tt.tags, got, tt.want)
			}
		})
	}
}

func TestSLI_GetNormalizedTags(t *testing.T) {
	tests := []struct {
		name string
		sli  SLI
		want []string
	}{
		{
			name: "with tags",
			sli:  SLI{Tags: []string{"supply-chain", "ai"}},
			want: []string{"ai", "supply-chain"},
		},
		{
			name: "no tags",
			sli:  SLI{Tags: nil},
			want: nil,
		},
		{
			name: "empty tags",
			sli:  SLI{Tags: []string{}},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sli.GetNormalizedTags()
			if !slicesEqual(got, tt.want) {
				t.Errorf("SLI.GetNormalizedTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategorySortWeight(t *testing.T) {
	weights := CategorySortWeight()

	// Test NIST CSF canonical order
	nistOrder := []string{
		CategoryGovern,   // 1
		CategoryIdentify, // 2
		CategoryProtect,  // 3
		CategoryDetect,   // 4
		CategoryRespond,  // 5
		CategoryRecover,  // 6
	}

	for i := 0; i < len(nistOrder)-1; i++ {
		current := nistOrder[i]
		next := nistOrder[i+1]
		if weights[current] >= weights[next] {
			t.Errorf("CategorySortWeight: %s (%d) should be less than %s (%d)",
				current, weights[current], next, weights[next])
		}
	}

	// Test variations map to same priority as canonical
	variations := map[string]string{
		"governance": CategoryGovern,
		"prevention": CategoryProtect,
		"detection":  CategoryDetect,
		"response":   CategoryRespond,
		"recovery":   CategoryRecover,
	}

	for variant, canonical := range variations {
		if weights[variant] != weights[canonical] {
			t.Errorf("CategorySortWeight: %s (%d) should equal %s (%d)",
				variant, weights[variant], canonical, weights[canonical])
		}
	}

	// Test operations categories sort after NIST CSF
	opsCats := []string{"reliability", "efficiency", "quality", "availability"}
	for _, opsCat := range opsCats {
		for _, nistCat := range nistOrder {
			if weights[opsCat] <= weights[nistCat] {
				t.Errorf("CategorySortWeight: %s (%d) should be greater than %s (%d)",
					opsCat, weights[opsCat], nistCat, weights[nistCat])
			}
		}
	}
}

func TestDefaultCategoryOrder(t *testing.T) {
	order := DefaultCategoryOrder()

	expected := []string{
		CategoryGovern,
		CategoryIdentify,
		CategoryProtect,
		CategoryDetect,
		CategoryRespond,
		CategoryRecover,
	}

	if !slicesEqual(order, expected) {
		t.Errorf("DefaultCategoryOrder() = %v, want %v", order, expected)
	}
}

func TestSpec_GetCategoryOrder(t *testing.T) {
	tests := []struct {
		name string
		spec Spec
		want []string
	}{
		{
			name: "with custom categories",
			spec: Spec{
				Categories: []Category{
					{ID: "detect"},
					{ID: "respond"},
					{ID: "protect"},
				},
			},
			want: []string{"detect", "respond", "protect"},
		},
		{
			name: "empty categories uses default",
			spec: Spec{
				Categories: []Category{},
			},
			want: DefaultCategoryOrder(),
		},
		{
			name: "nil categories uses default",
			spec: Spec{
				Categories: nil,
			},
			want: DefaultCategoryOrder(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.spec.GetCategoryOrder()
			if !slicesEqual(got, tt.want) {
				t.Errorf("Spec.GetCategoryOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpec_GetCategoryByID(t *testing.T) {
	spec := Spec{
		Categories: []Category{
			{ID: "detect", Name: "Detection", Description: "Detect threats"},
			{ID: "respond", Name: "Response", Description: "Respond to incidents"},
		},
	}

	tests := []struct {
		name     string
		id       string
		wantNil  bool
		wantName string
	}{
		{"found", "detect", false, "Detection"},
		{"found second", "respond", false, "Response"},
		{"not found", "govern", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := spec.GetCategoryByID(tt.id)
			if tt.wantNil {
				if got != nil {
					t.Errorf("Spec.GetCategoryByID(%q) = %v, want nil", tt.id, got)
				}
			} else {
				if got == nil {
					t.Errorf("Spec.GetCategoryByID(%q) = nil, want non-nil", tt.id)
				} else if got.Name != tt.wantName {
					t.Errorf("Spec.GetCategoryByID(%q).Name = %q, want %q", tt.id, got.Name, tt.wantName)
				}
			}
		})
	}
}

func TestSpec_GetSLIOrderForCategory(t *testing.T) {
	spec := Spec{
		Categories: []Category{
			{ID: "detect", SLIOrder: []string{"sli-a", "sli-b", "sli-c"}},
			{ID: "respond", SLIOrder: nil},
		},
	}

	tests := []struct {
		name string
		id   string
		want []string
	}{
		{"with order", "detect", []string{"sli-a", "sli-b", "sli-c"}},
		{"nil order", "respond", nil},
		{"not found", "govern", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := spec.GetSLIOrderForCategory(tt.id)
			if !slicesEqual(got, tt.want) {
				t.Errorf("Spec.GetSLIOrderForCategory(%q) = %v, want %v", tt.id, got, tt.want)
			}
		})
	}
}

func TestCategorySortWeight_Sorting(t *testing.T) {
	// Test that sorting by weight produces NIST CSF order
	categories := []string{"recover", "govern", "detect", "protect", "respond", "identify"}

	weights := CategorySortWeight()
	sort.Slice(categories, func(i, j int) bool {
		return weights[categories[i]] < weights[categories[j]]
	})

	expected := []string{"govern", "identify", "protect", "detect", "respond", "recover"}
	if !slicesEqual(categories, expected) {
		t.Errorf("Sorted categories = %v, want %v", categories, expected)
	}
}

func TestRecommendedTags(t *testing.T) {
	tags := RecommendedTags()

	// Verify all recommended tags are valid
	for _, tag := range tags {
		if err := ValidateTag(tag); err != nil {
			t.Errorf("RecommendedTag %q is invalid: %v", tag, err)
		}
	}

	// Verify expected tags are present
	expected := []string{TagAI, TagShiftLeft, TagSupplyChain}
	for _, exp := range expected {
		found := false
		for _, tag := range tags {
			if tag == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("RecommendedTags() missing expected tag %q", exp)
		}
	}
}

// slicesEqual compares two string slices for equality
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
