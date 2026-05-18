// Package dmaic provides types and utilities for DMAIC (Define, Measure, Analyze,
// Improve, Control) metrics framework documents.
//
// DMAIC is a data-driven quality strategy used for improving, optimizing, and
// stabilizing business processes and designs. This framework is commonly used
// in Six Sigma and continuous improvement methodologies.
//
// Key characteristics of DMAIC metrics:
//   - Metrics progress through five phases: Define, Measure, Analyze, Improve, Control
//   - Statistical Process Control (SPC) with UCL/LCL/CenterLine for Control phase
//   - Six Sigma metrics: Cp, Cpk, Sigma level, DPMO
//   - Root cause tracking for Analyze phase
//   - Initiative linkage connects improvements to metrics
//   - Historical data points support trend analysis
package dmaic

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Phase constants for DMAIC methodology.
const (
	PhaseDefine  = "Define"
	PhaseMeasure = "Measure"
	PhaseAnalyze = "Analyze"
	PhaseImprove = "Improve"
	PhaseControl = "Control"
)

// Trend direction constants for threshold interpretation.
const (
	TrendHigherBetter = "higher_better" // e.g., coverage percentages, adoption rates
	TrendLowerBetter  = "lower_better"  // e.g., vulnerabilities, incidents, response times
	TrendTargetValue  = "target_value"  // e.g., exact SLA targets
)

// Status indicator constants.
const (
	StatusGreen  = "Green"  // Within control limits / meeting targets
	StatusYellow = "Yellow" // Warning threshold breached
	StatusRed    = "Red"    // Critical threshold breached
)

// Document status constants for lifecycle.
const (
	DocumentStatusDraft     = "Draft"
	DocumentStatusActive    = "Active"
	DocumentStatusCompleted = "Completed"
	DocumentStatusArchived  = "Archived"
)

// Initiative status constants.
const (
	InitiativeStatusPlanned    = "Planned"
	InitiativeStatusInProgress = "In Progress"
	InitiativeStatusCompleted  = "Completed"
	InitiativeStatusCancelled  = "Cancelled"
)

// DMAICDocument represents a complete DMAIC metrics document.
type DMAICDocument struct {
	Schema      string       `json:"$schema,omitempty"`
	Metadata    *Metadata    `json:"metadata,omitempty"`
	Categories  []Category   `json:"categories"`            // Metric categories (ordered)
	Initiatives []Initiative `json:"initiatives,omitempty"` // Improvement projects
}

// Metadata contains document metadata.
type Metadata struct {
	ID            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	Owner         string    `json:"owner,omitempty"`
	Team          string    `json:"team,omitempty"`
	Period        string    `json:"period,omitempty"` // e.g., "2025-Q1", "FY2025"
	Version       string    `json:"version,omitempty"`
	Status        string    `json:"status,omitempty"`
	ReviewCadence string    `json:"reviewCadence,omitempty"` // e.g., "Weekly", "Monthly", "Quarterly"
	CreatedAt     time.Time `json:"createdAt,omitzero"`
	UpdatedAt     time.Time `json:"updatedAt,omitzero"`
}

// Category represents a group of related metrics.
type Category struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Order       int      `json:"order,omitempty"` // Display order
	Owner       string   `json:"owner,omitempty"`
	Metrics     []Metric `json:"metrics"`
}

// Metric represents a single measurable metric in the DMAIC framework.
type Metric struct {
	// Core fields
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order,omitempty"` // Display order within category

	// Values
	Baseline float64 `json:"baseline"`       // Starting value
	Current  float64 `json:"current"`        // Current value
	Target   float64 `json:"target"`         // Target value
	Unit     string  `json:"unit,omitempty"` // Unit of measurement

	// DMAIC phase
	Phase string `json:"phase"` // Define, Measure, Analyze, Improve, Control

	// Trend direction for threshold interpretation
	TrendDirection string `json:"trendDirection,omitempty"` // higher_better, lower_better, target_value

	// Statistical Process Control
	ControlLimits *ControlLimits `json:"controlLimits,omitempty"` // UCL, LCL, CenterLine

	// Thresholds
	Thresholds *Thresholds `json:"thresholds,omitempty"` // Warning and critical thresholds

	// Tracking
	Frequency  string `json:"frequency,omitempty"`  // Measurement frequency (e.g., "Daily", "Weekly")
	DataSource string `json:"dataSource,omitempty"` // Where data comes from
	Owner      string `json:"owner,omitempty"`
	Status     string `json:"status,omitempty"` // Green, Yellow, Red

	// Six Sigma
	ProcessCapability *ProcessCapability `json:"processCapability,omitempty"`

	// History
	DataPoints []DataPoint `json:"dataPoints,omitempty"` // Historical measurements

	// Analysis
	RootCauses []RootCause `json:"rootCauses,omitempty"` // Identified root causes

	// Links
	InitiativeIDs []string          `json:"initiativeIds,omitempty"` // Linked improvement initiatives
	ExternalLinks map[string]string `json:"externalLinks,omitempty"` // External system links (Jira, Datadog, etc.)
}

// ControlLimits represents Statistical Process Control limits.
type ControlLimits struct {
	UCL        float64 `json:"ucl"`             // Upper Control Limit
	LCL        float64 `json:"lcl"`             // Lower Control Limit
	CenterLine float64 `json:"centerLine"`      // Center Line (mean)
	Sigma      float64 `json:"sigma,omitempty"` // Standard deviation
}

// Thresholds represents warning and critical threshold values.
type Thresholds struct {
	Warning  float64 `json:"warning,omitempty"`  // Warning threshold
	Critical float64 `json:"critical,omitempty"` // Critical threshold
}

// ProcessCapability represents Six Sigma process capability metrics.
type ProcessCapability struct {
	Cp         float64 `json:"cp,omitempty"`         // Process Capability
	Cpk        float64 `json:"cpk,omitempty"`        // Process Capability Index
	SigmaLevel float64 `json:"sigmaLevel,omitempty"` // Sigma level (e.g., 6.0)
	DPMO       float64 `json:"dpmo,omitempty"`       // Defects Per Million Opportunities
}

// DataPoint represents a single historical measurement.
type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Notes     string    `json:"notes,omitempty"`
}

// RootCause represents an identified root cause for a metric issue.
type RootCause struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
	Category    string `json:"category,omitempty"`  // e.g., "Process", "People", "Technology"
	Impact      string `json:"impact,omitempty"`    // e.g., "High", "Medium", "Low"
	Validated   bool   `json:"validated,omitempty"` // Whether root cause has been validated
}

// Initiative represents an improvement project linked to metrics.
type Initiative struct {
	ID             string            `json:"id,omitempty"`
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	Owner          string            `json:"owner,omitempty"`
	Status         string            `json:"status,omitempty"` // Planned, In Progress, Completed, Cancelled
	StartDate      string            `json:"startDate,omitempty"`
	EndDate        string            `json:"endDate,omitempty"`
	MetricIDs      []string          `json:"metricIds,omitempty"`      // Linked metrics
	ExpectedImpact string            `json:"expectedImpact,omitempty"` // Description of expected improvement
	ActualImpact   string            `json:"actualImpact,omitempty"`   // Description of actual improvement
	ExternalLinks  map[string]string `json:"externalLinks,omitempty"`  // External system links
}

// DefaultFilename is the standard DMAIC filename.
const DefaultFilename = "dmaic.json"

// New creates a new DMAIC document with required fields initialized.
func New(id, name, owner string) *DMAICDocument {
	now := time.Now()
	return &DMAICDocument{
		Metadata: &Metadata{
			ID:        id,
			Name:      name,
			Owner:     owner,
			Status:    DocumentStatusDraft,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Categories:  []Category{},
		Initiatives: []Initiative{},
	}
}

// GenerateID generates a DMAIC ID based on the current date.
// Format: DMAIC-YYYY-DDD where DDD is the day of year.
func GenerateID() string {
	now := time.Now()
	return fmt.Sprintf("DMAIC-%d-%03d", now.Year(), now.YearDay())
}

// ReadFile reads a DMAIC document from a JSON file.
func ReadFile(filepath string) (*DMAICDocument, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return Parse(data)
}

// Parse parses DMAIC JSON data.
func Parse(data []byte) (*DMAICDocument, error) {
	var doc DMAICDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return &doc, nil
}

// JSON returns the DMAIC document as formatted JSON.
func (doc *DMAICDocument) JSON() ([]byte, error) {
	return json.MarshalIndent(doc, "", "  ")
}

// WriteFile writes the DMAIC document to a JSON file.
func (doc *DMAICDocument) WriteFile(filepath string) error {
	data, err := doc.JSON()
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// AllMetrics returns all metrics from all categories, flattened.
func (doc *DMAICDocument) AllMetrics() []Metric {
	var all []Metric
	for _, cat := range doc.Categories {
		all = append(all, cat.Metrics...)
	}
	return all
}

// MetricsByPhase returns metrics grouped by DMAIC phase.
func (doc *DMAICDocument) MetricsByPhase() map[string][]Metric {
	result := make(map[string][]Metric)
	for _, cat := range doc.Categories {
		for _, m := range cat.Metrics {
			phase := m.Phase
			if phase == "" {
				phase = PhaseDefine
			}
			result[phase] = append(result[phase], m)
		}
	}
	return result
}

// MetricsByStatus returns metrics grouped by status.
func (doc *DMAICDocument) MetricsByStatus() map[string][]Metric {
	result := make(map[string][]Metric)
	for _, cat := range doc.Categories {
		for _, m := range cat.Metrics {
			status := m.CalculateStatus()
			result[status] = append(result[status], m)
		}
	}
	return result
}

// CalculateOverallHealth calculates the overall health of the DMAIC document.
// Returns a value between 0.0 and 1.0 where:
//   - 1.0 = all metrics are Green
//   - 0.5 = all metrics are Yellow
//   - 0.0 = all metrics are Red
func (doc *DMAICDocument) CalculateOverallHealth() float64 {
	metrics := doc.AllMetrics()
	if len(metrics) == 0 {
		return 1.0 // No metrics = healthy
	}

	var total float64
	for _, m := range metrics {
		switch m.CalculateStatus() {
		case StatusGreen:
			total += 1.0
		case StatusYellow:
			total += 0.5
		case StatusRed:
			total += 0.0
		}
	}
	return total / float64(len(metrics))
}

// CalculateCategoryHealth calculates the health of a category.
// Returns a value between 0.0 and 1.0.
func (c *Category) CalculateCategoryHealth() float64 {
	if len(c.Metrics) == 0 {
		return 1.0 // No metrics = healthy
	}

	var total float64
	for _, m := range c.Metrics {
		switch m.CalculateStatus() {
		case StatusGreen:
			total += 1.0
		case StatusYellow:
			total += 0.5
		case StatusRed:
			total += 0.0
		}
	}
	return total / float64(len(c.Metrics))
}

// IsInControl checks if a metric is within its statistical control limits.
// Returns true if the metric has no control limits defined or is within limits.
func (m *Metric) IsInControl() bool {
	if m.ControlLimits == nil {
		return true // No limits = in control
	}
	return m.Current >= m.ControlLimits.LCL && m.Current <= m.ControlLimits.UCL
}

// CalculateStatus calculates the status based on thresholds and trend direction.
// Returns StatusGreen, StatusYellow, or StatusRed.
func (m *Metric) CalculateStatus() string {
	// If status is explicitly set, use it
	if m.Status != "" {
		return m.Status
	}

	// If no thresholds defined, check if target is met
	if m.Thresholds == nil {
		return m.calculateStatusFromTarget()
	}

	// Calculate status based on trend direction and thresholds
	switch m.TrendDirection {
	case TrendHigherBetter:
		return m.calculateStatusHigherBetter()
	case TrendLowerBetter:
		return m.calculateStatusLowerBetter()
	case TrendTargetValue:
		return m.calculateStatusTargetValue()
	default:
		// Default to higher is better
		return m.calculateStatusHigherBetter()
	}
}

// calculateStatusFromTarget determines status based on progress toward target.
func (m *Metric) calculateStatusFromTarget() string {
	if m.Target == 0 {
		return StatusGreen
	}

	var progress float64
	if m.TrendDirection == TrendLowerBetter {
		// For lower is better, invert the progress calculation
		if m.Baseline == 0 {
			if m.Current <= m.Target {
				return StatusGreen
			}
			return StatusRed
		}
		progress = (m.Baseline - m.Current) / (m.Baseline - m.Target)
	} else {
		// For higher is better or default
		if m.Target == m.Baseline {
			if m.Current >= m.Target {
				return StatusGreen
			}
			return StatusRed
		}
		progress = (m.Current - m.Baseline) / (m.Target - m.Baseline)
	}

	switch {
	case progress >= 1.0:
		return StatusGreen
	case progress >= 0.5:
		return StatusYellow
	default:
		return StatusRed
	}
}

// calculateStatusHigherBetter calculates status when higher values are better.
func (m *Metric) calculateStatusHigherBetter() string {
	if m.Thresholds == nil {
		return StatusGreen
	}

	switch {
	case m.Current >= m.Target:
		return StatusGreen
	case m.Thresholds.Warning > 0 && m.Current >= m.Thresholds.Warning:
		return StatusGreen
	case m.Thresholds.Critical > 0 && m.Current < m.Thresholds.Critical:
		return StatusRed
	case m.Thresholds.Warning > 0 && m.Current < m.Thresholds.Warning:
		return StatusYellow
	default:
		return StatusGreen
	}
}

// calculateStatusLowerBetter calculates status when lower values are better.
func (m *Metric) calculateStatusLowerBetter() string {
	if m.Thresholds == nil {
		return StatusGreen
	}

	switch {
	case m.Current <= m.Target:
		return StatusGreen
	case m.Thresholds.Warning > 0 && m.Current <= m.Thresholds.Warning:
		return StatusGreen
	case m.Thresholds.Critical > 0 && m.Current > m.Thresholds.Critical:
		return StatusRed
	case m.Thresholds.Warning > 0 && m.Current > m.Thresholds.Warning:
		return StatusYellow
	default:
		return StatusGreen
	}
}

// calculateStatusTargetValue calculates status when a specific value is the target.
func (m *Metric) calculateStatusTargetValue() string {
	if m.Thresholds == nil {
		if m.Current == m.Target {
			return StatusGreen
		}
		return StatusYellow
	}

	deviation := abs(m.Current - m.Target)

	switch {
	case deviation == 0:
		return StatusGreen
	case m.Thresholds.Critical > 0 && deviation > m.Thresholds.Critical:
		return StatusRed
	case m.Thresholds.Warning > 0 && deviation > m.Thresholds.Warning:
		return StatusYellow
	default:
		return StatusGreen
	}
}

// abs returns the absolute value of a float64.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// StatusDescription returns a description for a status.
func StatusDescription(status string) string {
	switch status {
	case StatusGreen:
		return "Meeting target"
	case StatusYellow:
		return "Warning - attention needed"
	case StatusRed:
		return "Critical - immediate action required"
	default:
		return "Unknown status"
	}
}

// PhaseDescription returns a description for a DMAIC phase.
func PhaseDescription(phase string) string {
	switch phase {
	case PhaseDefine:
		return "Define the problem and project goals"
	case PhaseMeasure:
		return "Measure current performance"
	case PhaseAnalyze:
		return "Analyze root causes"
	case PhaseImprove:
		return "Implement improvements"
	case PhaseControl:
		return "Sustain the gains"
	default:
		return "Unknown phase"
	}
}

// Phases returns all DMAIC phases in order.
func Phases() []string {
	return []string{PhaseDefine, PhaseMeasure, PhaseAnalyze, PhaseImprove, PhaseControl}
}

// ValidPhase checks if a phase value is valid.
func ValidPhase(phase string) bool {
	switch phase {
	case PhaseDefine, PhaseMeasure, PhaseAnalyze, PhaseImprove, PhaseControl:
		return true
	default:
		return false
	}
}

// ValidTrendDirection checks if a trend direction value is valid.
func ValidTrendDirection(trend string) bool {
	switch trend {
	case TrendHigherBetter, TrendLowerBetter, TrendTargetValue, "":
		return true
	default:
		return false
	}
}

// ValidStatus checks if a status value is valid.
func ValidStatus(status string) bool {
	switch status {
	case StatusGreen, StatusYellow, StatusRed, "":
		return true
	default:
		return false
	}
}
