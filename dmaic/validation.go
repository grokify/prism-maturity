package dmaic

import (
	"fmt"
	"strings"
)

// ValidationError represents a validation issue.
type ValidationError struct {
	Path    string // JSON path to the problematic field
	Message string
	IsError bool // true for errors, false for warnings
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationOptions configures validation behavior.
type ValidationOptions struct {
	RequireCategories        bool // Require at least one category
	RequireMetrics           bool // Require at least one metric per category
	RequirePhase             bool // Require phase to be set on metrics
	RequireTarget            bool // Require target values on metrics
	RequireThresholds        bool // Require thresholds on metrics
	RequireControlLimits     bool // Require control limits for Control phase metrics
	RequireProcessCapability bool // Require process capability metrics
	ValidateControlLimits    bool // Validate UCL > CenterLine > LCL
	ValidateThresholds       bool // Validate threshold logic based on trend direction
	ValidateOrdering         bool // Validate ordering uniqueness
	MinMetricsPerCategory    int  // Minimum metrics per category
	MaxMetricsPerCategory    int  // Maximum metrics per category (0 = no limit)
	MaxCategories            int  // Maximum categories per document (0 = no limit)
}

// DefaultValidationOptions returns permissive defaults.
func DefaultValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		RequireCategories:        true,
		RequireMetrics:           true,
		RequirePhase:             false,
		RequireTarget:            false,
		RequireThresholds:        false,
		RequireControlLimits:     false,
		RequireProcessCapability: false,
		ValidateControlLimits:    true,
		ValidateThresholds:       true,
		ValidateOrdering:         false,
		MinMetricsPerCategory:    1,
		MaxMetricsPerCategory:    0,
		MaxCategories:            0,
	}
}

// StrictValidationOptions returns production-ready validation settings.
func StrictValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		RequireCategories:        true,
		RequireMetrics:           true,
		RequirePhase:             true,
		RequireTarget:            true,
		RequireThresholds:        false,
		RequireControlLimits:     false,
		RequireProcessCapability: false,
		ValidateControlLimits:    true,
		ValidateThresholds:       true,
		ValidateOrdering:         true,
		MinMetricsPerCategory:    1,
		MaxMetricsPerCategory:    10,
		MaxCategories:            10,
	}
}

// SixSigmaValidationOptions returns Six Sigma compliance validation settings.
func SixSigmaValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		RequireCategories:        true,
		RequireMetrics:           true,
		RequirePhase:             true,
		RequireTarget:            true,
		RequireThresholds:        true,
		RequireControlLimits:     true, // Required for Control phase
		RequireProcessCapability: true,
		ValidateControlLimits:    true,
		ValidateThresholds:       true,
		ValidateOrdering:         true,
		MinMetricsPerCategory:    1,
		MaxMetricsPerCategory:    0,
		MaxCategories:            0,
	}
}

// Validate checks the DMAIC document for issues.
func (doc *DMAICDocument) Validate(opts *ValidationOptions) []ValidationError {
	if opts == nil {
		opts = DefaultValidationOptions()
	}

	var errs []ValidationError

	// Validate categories exist
	if opts.RequireCategories && len(doc.Categories) == 0 {
		errs = append(errs, ValidationError{
			Path:    "categories",
			Message: "at least one category is required",
			IsError: true,
		})
	}

	// Check max categories
	if opts.MaxCategories > 0 && len(doc.Categories) > opts.MaxCategories {
		errs = append(errs, ValidationError{
			Path:    "categories",
			Message: fmt.Sprintf("too many categories: %d (max: %d)", len(doc.Categories), opts.MaxCategories),
			IsError: false,
		})
	}

	// Validate ordering uniqueness for categories
	if opts.ValidateOrdering {
		errs = append(errs, validateCategoryOrdering(doc.Categories)...)
	}

	// Validate each category
	for i, cat := range doc.Categories {
		catPath := fmt.Sprintf("categories[%d]", i)
		errs = append(errs, validateCategory(cat, catPath, opts)...)
	}

	// Validate metadata if present
	if doc.Metadata != nil {
		errs = append(errs, validateMetadata(doc.Metadata)...)
	}

	// Validate initiatives
	for i, init := range doc.Initiatives {
		initPath := fmt.Sprintf("initiatives[%d]", i)
		errs = append(errs, validateInitiative(init, initPath)...)
	}

	return errs
}

func validateCategory(cat Category, path string, opts *ValidationOptions) []ValidationError {
	var errs []ValidationError

	// Name is required
	if strings.TrimSpace(cat.Name) == "" {
		errs = append(errs, ValidationError{
			Path:    path + ".name",
			Message: "category name is required",
			IsError: true,
		})
	}

	// Validate metrics exist
	if opts.RequireMetrics && len(cat.Metrics) == 0 {
		errs = append(errs, ValidationError{
			Path:    path + ".metrics",
			Message: "at least one metric is required",
			IsError: true,
		})
	}

	// Check minimum metrics
	if opts.MinMetricsPerCategory > 0 && len(cat.Metrics) < opts.MinMetricsPerCategory {
		errs = append(errs, ValidationError{
			Path:    path + ".metrics",
			Message: fmt.Sprintf("too few metrics: %d (min: %d)", len(cat.Metrics), opts.MinMetricsPerCategory),
			IsError: false,
		})
	}

	// Check maximum metrics
	if opts.MaxMetricsPerCategory > 0 && len(cat.Metrics) > opts.MaxMetricsPerCategory {
		errs = append(errs, ValidationError{
			Path:    path + ".metrics",
			Message: fmt.Sprintf("too many metrics: %d (max: %d)", len(cat.Metrics), opts.MaxMetricsPerCategory),
			IsError: false,
		})
	}

	// Validate ordering uniqueness for metrics
	if opts.ValidateOrdering {
		errs = append(errs, validateMetricOrdering(cat.Metrics, path)...)
	}

	// Validate each metric
	for j, m := range cat.Metrics {
		mPath := fmt.Sprintf("%s.metrics[%d]", path, j)
		errs = append(errs, validateMetric(m, mPath, opts)...)
	}

	return errs
}

func validateMetric(m Metric, path string, opts *ValidationOptions) []ValidationError {
	var errs []ValidationError

	// Name is required
	if strings.TrimSpace(m.Name) == "" {
		errs = append(errs, ValidationError{
			Path:    path + ".name",
			Message: "metric name is required",
			IsError: true,
		})
	}

	// Phase validation
	if opts.RequirePhase && strings.TrimSpace(m.Phase) == "" {
		errs = append(errs, ValidationError{
			Path:    path + ".phase",
			Message: "metric phase is required",
			IsError: true,
		})
	}

	if m.Phase != "" && !ValidPhase(m.Phase) {
		errs = append(errs, ValidationError{
			Path:    path + ".phase",
			Message: fmt.Sprintf("invalid phase: %s (expected: Define, Measure, Analyze, Improve, Control)", m.Phase),
			IsError: true,
		})
	}

	// Trend direction validation
	if m.TrendDirection != "" && !ValidTrendDirection(m.TrendDirection) {
		errs = append(errs, ValidationError{
			Path:    path + ".trendDirection",
			Message: fmt.Sprintf("invalid trend direction: %s (expected: higher_better, lower_better, target_value)", m.TrendDirection),
			IsError: true,
		})
	}

	// Status validation
	if m.Status != "" && !ValidStatus(m.Status) {
		errs = append(errs, ValidationError{
			Path:    path + ".status",
			Message: fmt.Sprintf("invalid status: %s (expected: Green, Yellow, Red)", m.Status),
			IsError: true,
		})
	}

	// Target validation
	if opts.RequireTarget && m.Target == 0 {
		errs = append(errs, ValidationError{
			Path:    path + ".target",
			Message: "target is required",
			IsError: false, // Warning - target might legitimately be 0
		})
	}

	// Thresholds validation
	if opts.RequireThresholds && m.Thresholds == nil {
		errs = append(errs, ValidationError{
			Path:    path + ".thresholds",
			Message: "thresholds are required",
			IsError: false,
		})
	}

	if m.Thresholds != nil && opts.ValidateThresholds {
		errs = append(errs, validateThresholds(m.Thresholds, m.TrendDirection, path+".thresholds")...)
	}

	// Control limits validation
	if opts.RequireControlLimits && m.Phase == PhaseControl && m.ControlLimits == nil {
		errs = append(errs, ValidationError{
			Path:    path + ".controlLimits",
			Message: "control limits are required for Control phase metrics",
			IsError: false,
		})
	}

	if m.ControlLimits != nil && opts.ValidateControlLimits {
		errs = append(errs, validateControlLimits(m.ControlLimits, path+".controlLimits")...)
	}

	// Process capability validation
	if opts.RequireProcessCapability && m.ProcessCapability == nil {
		errs = append(errs, ValidationError{
			Path:    path + ".processCapability",
			Message: "process capability metrics are required",
			IsError: false,
		})
	}

	return errs
}

func validateControlLimits(cl *ControlLimits, path string) []ValidationError {
	var errs []ValidationError

	// UCL > CenterLine > LCL
	if cl.UCL < cl.CenterLine {
		errs = append(errs, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("UCL (%.2f) must be greater than or equal to CenterLine (%.2f)", cl.UCL, cl.CenterLine),
			IsError: true,
		})
	}

	if cl.CenterLine < cl.LCL {
		errs = append(errs, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("CenterLine (%.2f) must be greater than or equal to LCL (%.2f)", cl.CenterLine, cl.LCL),
			IsError: true,
		})
	}

	// Sigma should be positive if set
	if cl.Sigma < 0 {
		errs = append(errs, ValidationError{
			Path:    path + ".sigma",
			Message: fmt.Sprintf("sigma must be non-negative, got %.2f", cl.Sigma),
			IsError: true,
		})
	}

	return errs
}

func validateThresholds(t *Thresholds, trendDirection string, path string) []ValidationError {
	var errs []ValidationError

	// Validate threshold logic based on trend direction
	switch trendDirection {
	case TrendHigherBetter:
		// Warning should be >= Critical (higher is better, so warning is less severe)
		if t.Warning > 0 && t.Critical > 0 && t.Warning < t.Critical {
			errs = append(errs, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("for higher_better trend, warning (%.2f) should be >= critical (%.2f)", t.Warning, t.Critical),
				IsError: false, // Warning - might be intentional
			})
		}
	case TrendLowerBetter:
		// Warning should be <= Critical (lower is better, so warning is less severe)
		if t.Warning > 0 && t.Critical > 0 && t.Warning > t.Critical {
			errs = append(errs, ValidationError{
				Path:    path,
				Message: fmt.Sprintf("for lower_better trend, warning (%.2f) should be <= critical (%.2f)", t.Warning, t.Critical),
				IsError: false, // Warning - might be intentional
			})
		}
	}

	return errs
}

func validateCategoryOrdering(categories []Category) []ValidationError {
	var errs []ValidationError

	// Check for duplicate order values (excluding 0 which means unset)
	seen := make(map[int]int) // order -> first index seen
	for i, cat := range categories {
		if cat.Order > 0 {
			if firstIdx, exists := seen[cat.Order]; exists {
				errs = append(errs, ValidationError{
					Path:    fmt.Sprintf("categories[%d].order", i),
					Message: fmt.Sprintf("duplicate order value %d (also at categories[%d])", cat.Order, firstIdx),
					IsError: false,
				})
			} else {
				seen[cat.Order] = i
			}
		}
	}

	return errs
}

func validateMetricOrdering(metrics []Metric, categoryPath string) []ValidationError {
	var errs []ValidationError

	// Check for duplicate order values (excluding 0 which means unset)
	seen := make(map[int]int) // order -> first index seen
	for i, m := range metrics {
		if m.Order > 0 {
			if firstIdx, exists := seen[m.Order]; exists {
				errs = append(errs, ValidationError{
					Path:    fmt.Sprintf("%s.metrics[%d].order", categoryPath, i),
					Message: fmt.Sprintf("duplicate order value %d (also at index %d)", m.Order, firstIdx),
					IsError: false,
				})
			} else {
				seen[m.Order] = i
			}
		}
	}

	return errs
}

func validateMetadata(meta *Metadata) []ValidationError {
	var errs []ValidationError

	// Name should be present
	if strings.TrimSpace(meta.Name) == "" {
		errs = append(errs, ValidationError{
			Path:    "metadata.name",
			Message: "DMAIC document name is recommended",
			IsError: false,
		})
	}

	// Status should be valid if present
	if meta.Status != "" &&
		meta.Status != DocumentStatusDraft &&
		meta.Status != DocumentStatusActive &&
		meta.Status != DocumentStatusCompleted &&
		meta.Status != DocumentStatusArchived {
		errs = append(errs, ValidationError{
			Path:    "metadata.status",
			Message: fmt.Sprintf("invalid status: %s (expected: Draft, Active, Completed, Archived)", meta.Status),
			IsError: true,
		})
	}

	return errs
}

func validateInitiative(init Initiative, path string) []ValidationError {
	var errs []ValidationError

	// Name is required
	if strings.TrimSpace(init.Name) == "" {
		errs = append(errs, ValidationError{
			Path:    path + ".name",
			Message: "initiative name is required",
			IsError: true,
		})
	}

	// Status should be valid if present
	if init.Status != "" &&
		init.Status != InitiativeStatusPlanned &&
		init.Status != InitiativeStatusInProgress &&
		init.Status != InitiativeStatusCompleted &&
		init.Status != InitiativeStatusCancelled {
		errs = append(errs, ValidationError{
			Path:    path + ".status",
			Message: fmt.Sprintf("invalid initiative status: %s (expected: Planned, In Progress, Completed, Cancelled)", init.Status),
			IsError: true,
		})
	}

	return errs
}

// Errors returns only error-level validation results.
func Errors(errs []ValidationError) []ValidationError {
	var result []ValidationError
	for _, e := range errs {
		if e.IsError {
			result = append(result, e)
		}
	}
	return result
}

// Warnings returns only warning-level validation results.
func Warnings(errs []ValidationError) []ValidationError {
	var result []ValidationError
	for _, e := range errs {
		if !e.IsError {
			result = append(result, e)
		}
	}
	return result
}

// IsValid returns true if there are no error-level validation issues.
func IsValid(errs []ValidationError) bool {
	return len(Errors(errs)) == 0
}
