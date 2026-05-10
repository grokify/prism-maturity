// Package maturity provides types and functions for maturity model management.
package maturity

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Spec defines a complete maturity specification for an organization.
// It contains maturity models for multiple domains.
type Spec struct {
	Schema        string                       `json:"$schema,omitempty"`
	Metadata      *SpecMetadata                `json:"metadata,omitempty"`
	SLIs          map[string]*SLI              `json:"slis,omitempty"`          // Service Level Indicators with framework mappings
	KPIThresholds map[string][]KPIThreshold    `json:"kpiThresholds,omitempty"` // Deprecated: use SLIs instead
	Domains       map[string]*DomainModel      `json:"domains"`
	Assessments   map[string]*DomainAssessment `json:"assessments,omitempty"`
}

// SLI defines a Service Level Indicator (the metric being measured).
// Framework mappings are defined here since they apply to the metric itself,
// not to specific targets at different maturity levels.
type SLI struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// Metric definition
	MetricName string `json:"metricName"`     // Human-readable metric description
	Unit       string `json:"unit,omitempty"` // %, days, count, seconds, etc.
	Type       string `json:"type,omitempty"` // "quantitative" (default), "qualitative"

	// Classification
	Layer    string `json:"layer,omitempty"`    // requirements, code, infra, runtime, adoption, support
	Category string `json:"category,omitempty"` // prevention, detection, response
	SLIType  string `json:"sliType,omitempty"`  // Observability type: availability, latency, error_rate, throughput, saturation, utilization, quality, freshness

	// Framework mappings - defined once on the SLI, inherited by all SLOs
	FrameworkMappings []FrameworkMapping `json:"frameworkMappings,omitempty"`
}

// KPIThreshold defines the progression of a KPI across maturity levels.
type KPIThreshold struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Unit        string          `json:"unit,omitempty"`
	Operator    string          `json:"operator,omitempty"` // gte (default), lte for "lower is better"
	Thresholds  LevelThresholds `json:"thresholds"`
	Current     any             `json:"current,omitempty"` // Current value for assessment
}

// LevelThresholds holds threshold values for each maturity level.
type LevelThresholds struct {
	M1 any `json:"m1,omitempty"`
	M2 any `json:"m2,omitempty"`
	M3 any `json:"m3,omitempty"`
	M4 any `json:"m4,omitempty"`
	M5 any `json:"m5,omitempty"`
}

// SpecMetadata contains metadata about the maturity specification.
type SpecMetadata struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Version      string `json:"version,omitempty"`
	Organization string `json:"organization,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
}

// DomainModel defines maturity levels for a specific domain.
type DomainModel struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Owner       string  `json:"owner,omitempty"`
	Levels      []Level `json:"levels"`
}

// Level defines a maturity level (M1-M5) for a domain.
type Level struct {
	Level       int         `json:"level"` // 1-5
	Name        string      `json:"name"`  // Reactive, Basic, Defined, Managed, Optimizing
	Description string      `json:"description"`
	Criteria    []Criterion `json:"criteria,omitempty"` // SLOs that define the level
	Enablers    []Enabler   `json:"enablers,omitempty"` // Tasks to achieve the level
}

// Criterion is a measurable SLO (Service Level Objective) that defines level achievement.
// It references an SLI and specifies a target threshold for a specific maturity level.
// Framework mappings are inherited from the referenced SLI.
type Criterion struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// SLI Reference - links to the metric definition with framework mappings
	SLIID string `json:"sliId,omitempty"` // Reference to SLI by ID

	// Inline SLI fields (for simple cases without separate SLI definition)
	Type       string `json:"type,omitempty"`       // "quantitative" (default), "qualitative"
	MetricName string `json:"metricName,omitempty"` // Human-readable metric description
	Unit       string `json:"unit,omitempty"`       // %, days, count, seconds, etc.

	// SLO Definition (target for this level)
	Operator string  `json:"operator"` // gte, lte, gt, lt, eq, exists
	Target   float64 `json:"target"`   // Target value (for quantitative)

	// Qualitative fields
	Status string `json:"status,omitempty"` // For qualitative: tracked, implemented, defined, etc.

	// Framework mappings - DEPRECATED: use SLI.FrameworkMappings instead
	// Kept for backward compatibility; if set, takes precedence over SLI mappings
	FrameworkMappings []FrameworkMapping `json:"frameworkMappings,omitempty"`

	// Assessment (populated during evaluation)
	Current float64 `json:"current,omitempty"` // Current value (for quantitative)
	IsMet   bool    `json:"isMet,omitempty"`   // Calculated: meets target?

	// Weighting
	Weight   float64 `json:"weight,omitempty"`   // Relative importance (default 1.0)
	Required bool    `json:"required,omitempty"` // Must pass for level (default true if omitted)
}

// FrameworkMapping maps a criterion to an external framework control.
type FrameworkMapping struct {
	Framework   string `json:"framework"`             // Framework identifier (e.g., "NIST_CSF_2", "NIST_800_53")
	Reference   string `json:"reference"`             // Control or function ID (e.g., "PR.DS-1", "AC-2")
	Name        string `json:"name,omitempty"`        // Human-readable control name
	Description string `json:"description,omitempty"` // Control description
	Baseline    string `json:"baseline,omitempty"`    // Required baseline level (e.g., "high", "moderate", "low")
	Version     string `json:"version,omitempty"`     // Framework version (e.g., "2.0", "Rev 5")
}

// IsQualitative returns true if this is a qualitative criterion.
func (c *Criterion) IsQualitative() bool {
	return c.Type == CriterionTypeQualitative || c.Operator == OperatorExists
}

// IsQualitativeWithSpec returns true if this is a qualitative criterion,
// checking both inline type and resolved SLI type.
func (c *Criterion) IsQualitativeWithSpec(spec *Spec) bool {
	if c.Operator == OperatorExists {
		return true
	}
	return c.GetType(spec) == CriterionTypeQualitative
}

// Criterion type constants.
const (
	CriterionTypeQuantitative = "quantitative" // Default: numeric comparison
	CriterionTypeQualitative  = "qualitative"  // Binary state tracking
)

// Enabler is implementation work to achieve criteria.
type Enabler struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// Classification
	Type  string `json:"type,omitempty"`  // implementation, process, training, tooling
	Layer string `json:"layer,omitempty"` // requirements, code, infra, runtime, adoption, support

	// Effort
	Effort string `json:"effort,omitempty"` // T-shirt size or duration
	Team   string `json:"team,omitempty"`   // Responsible team

	// Tracking
	Status      string   `json:"status,omitempty"`      // not_started, in_progress, completed
	CriteriaIDs []string `json:"criteriaIds,omitempty"` // Which SLOs this enables

	// Dependencies
	DependsOn []string `json:"dependsOn,omitempty"` // Other enabler IDs
}

// DomainAssessment captures current state against a domain's maturity model.
type DomainAssessment struct {
	Domain         string             `json:"domain"`
	AssessedAt     string             `json:"assessedAt,omitempty"`
	AssessedBy     string             `json:"assessedBy,omitempty"`
	CurrentLevel   int                `json:"currentLevel"`             // Achieved level (1-5)
	TargetLevel    int                `json:"targetLevel"`              // Goal level
	CriteriaValues map[string]float64 `json:"criteriaValues,omitempty"` // Current values by criterion ID (quantitative)
	CriteriaStatus map[string]string  `json:"criteriaStatus,omitempty"` // Current status by criterion ID (qualitative)
	EnablerStatus  map[string]string  `json:"enablerStatus,omitempty"`  // Status by enabler ID
}

// Enabler status constants.
const (
	StatusNotStarted = "not_started"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusBlocked    = "blocked"
)

// Enabler type constants.
const (
	TypeImplementation = "implementation"
	TypeProcess        = "process"
	TypeTraining       = "training"
	TypeTooling        = "tooling"
)

// Operator constants.
const (
	OpGTE          = "gte"    // Greater than or equal
	OpLTE          = "lte"    // Less than or equal
	OpGT           = "gt"     // Greater than
	OpLT           = "lt"     // Less than
	OpEQ           = "eq"     // Equal
	OperatorExists = "exists" // Qualitative: metric exists/is tracked
)

// Qualitative status constants for criteria.
const (
	QualStatusTracked     = "tracked"     // Metric is being tracked
	QualStatusImplemented = "implemented" // Control/feature is implemented
	QualStatusDefined     = "defined"     // Process/policy is defined
	QualStatusDocumented  = "documented"  // Documentation exists
	QualStatusCompliant   = "compliant"   // Meets compliance requirement
	QualStatusEnabled     = "enabled"     // Feature/capability is enabled
	QualStatusNotTracked  = "not_tracked" // Not yet being tracked
	QualStatusPartial     = "partial"     // Partially implemented
	QualStatusPlanned     = "planned"     // Planned but not started
)

// IsQualitativeStatusMet returns whether a qualitative status indicates compliance.
func IsQualitativeStatusMet(status string) bool {
	switch status {
	case QualStatusTracked, QualStatusImplemented,
		QualStatusDefined, QualStatusDocumented,
		QualStatusCompliant, QualStatusEnabled:
		return true
	default:
		return false
	}
}

// Level name constants.
const (
	LevelNameReactive   = "Reactive"
	LevelNameBasic      = "Basic"
	LevelNameDefined    = "Defined"
	LevelNameManaged    = "Managed"
	LevelNameOptimizing = "Optimizing"
)

// DefaultLevelNames returns the standard M1-M5 level names.
func DefaultLevelNames() map[int]string {
	return map[int]string{
		1: LevelNameReactive,
		2: LevelNameBasic,
		3: LevelNameDefined,
		4: LevelNameManaged,
		5: LevelNameOptimizing,
	}
}

// ReadSpecFile reads a maturity spec from a JSON file.
func ReadSpecFile(filename string) (*Spec, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &spec, nil
}

// WriteSpecFile writes a maturity spec to a JSON file.
func (s *Spec) WriteSpecFile(filename string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal spec: %w", err)
	}

	if err := os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}

	return nil
}

// GetDomain returns the domain model by name.
func (s *Spec) GetDomain(name string) (*DomainModel, bool) {
	domain, ok := s.Domains[strings.ToLower(name)]
	return domain, ok
}

// GetLevel returns the level definition for a domain.
func (d *DomainModel) GetLevel(level int) (*Level, bool) {
	for i := range d.Levels {
		if d.Levels[i].Level == level {
			return &d.Levels[i], true
		}
	}
	return nil, false
}

// OperatorSymbol returns the symbol for an operator.
func OperatorSymbol(op string) string {
	switch op {
	case OpGTE:
		return ">="
	case OpLTE:
		return "<="
	case OpGT:
		return ">"
	case OpLT:
		return "<"
	case OpEQ:
		return "="
	case OperatorExists:
		return "Tracked"
	default:
		return op
	}
}

// CheckMet checks if a criterion is met given a current value.
// For qualitative criteria, use CheckQualitativeMet instead.
func (c *Criterion) CheckMet(current float64) bool {
	// For qualitative criteria, check status instead
	if c.IsQualitative() {
		return c.CheckQualitativeMet()
	}

	switch c.Operator {
	case OpGTE:
		return current >= c.Target
	case OpLTE:
		return current <= c.Target
	case OpGT:
		return current > c.Target
	case OpLT:
		return current < c.Target
	case OpEQ:
		return current == c.Target
	default:
		return false
	}
}

// CheckQualitativeMet checks if a qualitative criterion is met based on its status.
func (c *Criterion) CheckQualitativeMet() bool {
	return IsQualitativeStatusMet(c.Status)
}

// GetSLI returns the SLI for this criterion from the spec.
// Returns nil if no SLI is referenced or not found.
func (c *Criterion) GetSLI(spec *Spec) *SLI {
	if c.SLIID == "" || spec == nil || spec.SLIs == nil {
		return nil
	}
	return spec.SLIs[c.SLIID]
}

// GetFrameworkMappings returns framework mappings for this criterion.
// If the criterion has inline mappings, those are returned.
// Otherwise, mappings are resolved from the referenced SLI.
func (c *Criterion) GetFrameworkMappings(spec *Spec) []FrameworkMapping {
	// Inline mappings take precedence (backward compatibility)
	if len(c.FrameworkMappings) > 0 {
		return c.FrameworkMappings
	}

	// Resolve from SLI
	if sli := c.GetSLI(spec); sli != nil {
		return sli.FrameworkMappings
	}

	return nil
}

// GetMetricName returns the metric name for this criterion.
// Resolves from SLI if not set inline.
func (c *Criterion) GetMetricName(spec *Spec) string {
	if c.MetricName != "" {
		return c.MetricName
	}
	if sli := c.GetSLI(spec); sli != nil {
		return sli.MetricName
	}
	return ""
}

// GetUnit returns the unit for this criterion.
// Resolves from SLI if not set inline.
func (c *Criterion) GetUnit(spec *Spec) string {
	if c.Unit != "" {
		return c.Unit
	}
	if sli := c.GetSLI(spec); sli != nil {
		return sli.Unit
	}
	return ""
}

// GetType returns the type for this criterion.
// Resolves from SLI if not set inline.
func (c *Criterion) GetType(spec *Spec) string {
	if c.Type != "" {
		return c.Type
	}
	if sli := c.GetSLI(spec); sli != nil {
		return sli.Type
	}
	return CriterionTypeQuantitative // default
}

// GetLayer returns the layer for this criterion from its SLI.
func (c *Criterion) GetLayer(spec *Spec) string {
	if sli := c.GetSLI(spec); sli != nil {
		return sli.Layer
	}
	return ""
}

// GetCategory returns the category for this criterion from its SLI.
func (c *Criterion) GetCategory(spec *Spec) string {
	if sli := c.GetSLI(spec); sli != nil {
		return sli.Category
	}
	return ""
}

// GetSLI returns an SLI by ID.
func (s *Spec) GetSLI(id string) *SLI {
	if s.SLIs == nil {
		return nil
	}
	return s.SLIs[id]
}

// AllSLIs returns all SLIs in the spec.
func (s *Spec) AllSLIs() []*SLI {
	var slis []*SLI
	for _, sli := range s.SLIs {
		slis = append(slis, sli)
	}
	return slis
}

// TargetString returns a formatted target string like ">=95%" or "Tracked" for qualitative.
func (c *Criterion) TargetString() string {
	if c.IsQualitative() {
		return "Tracked"
	}
	symbol := OperatorSymbol(c.Operator)
	if c.Unit != "" {
		return fmt.Sprintf("%s%.0f%s", symbol, c.Target, c.Unit)
	}
	return fmt.Sprintf("%s%.0f", symbol, c.Target)
}

// CurrentString returns a formatted current value string or status for qualitative.
func (c *Criterion) CurrentString() string {
	if c.IsQualitative() {
		if c.Status == "" {
			return "Not Tracked"
		}
		return formatQualitativeStatus(c.Status)
	}
	if c.Unit != "" {
		return fmt.Sprintf("%.1f%s", c.Current, c.Unit)
	}
	return fmt.Sprintf("%.1f", c.Current)
}

// formatQualitativeStatus returns a human-readable status string.
func formatQualitativeStatus(status string) string {
	switch status {
	case QualStatusTracked:
		return "Tracked"
	case QualStatusImplemented:
		return "Implemented"
	case QualStatusDefined:
		return "Defined"
	case QualStatusDocumented:
		return "Documented"
	case QualStatusCompliant:
		return "Compliant"
	case QualStatusEnabled:
		return "Enabled"
	case QualStatusNotTracked:
		return "Not Tracked"
	case QualStatusPartial:
		return "Partial"
	case QualStatusPlanned:
		return "Planned"
	default:
		return status
	}
}

// LevelProgress tracks progress toward a maturity level.
type LevelProgress struct {
	Level           int     `json:"level"`
	CriteriaMet     int     `json:"criteriaMet"`
	CriteriaTotal   int     `json:"criteriaTotal"`
	ProgressPercent float64 `json:"progressPercent"`
	EnablersDone    int     `json:"enablersDone"`
	EnablersTotal   int     `json:"enablersTotal"`
}

// CalculateLevelProgress calculates progress for a level given current values.
func (l *Level) CalculateLevelProgress(values map[string]float64, enablerStatus map[string]string) LevelProgress {
	progress := LevelProgress{
		Level:         l.Level,
		CriteriaTotal: len(l.Criteria),
		EnablersTotal: len(l.Enablers),
	}

	for _, c := range l.Criteria {
		if current, ok := values[c.ID]; ok {
			if c.CheckMet(current) {
				progress.CriteriaMet++
			}
		}
	}

	for _, e := range l.Enablers {
		if status, ok := enablerStatus[e.ID]; ok && status == StatusCompleted {
			progress.EnablersDone++
		}
	}

	if progress.CriteriaTotal > 0 {
		progress.ProgressPercent = float64(progress.CriteriaMet) / float64(progress.CriteriaTotal) * 100
	} else {
		progress.ProgressPercent = 100 // No criteria means level is achieved
	}

	return progress
}

// IsLevelAchieved checks if all required criteria for a level are met.
func (l *Level) IsLevelAchieved(values map[string]float64) bool {
	for _, c := range l.Criteria {
		// Default to required if not specified
		required := c.Required || c.Weight == 0
		if !required {
			continue
		}

		current, ok := values[c.ID]
		if !ok {
			return false // Missing value for required criterion
		}

		if !c.CheckMet(current) {
			return false
		}
	}
	return true
}

// AllCriteria returns all criteria across all levels for a domain.
func (d *DomainModel) AllCriteria() []Criterion {
	var all []Criterion
	for _, level := range d.Levels {
		for _, c := range level.Criteria {
			all = append(all, c)
		}
	}
	return all
}

// AllEnablers returns all enablers across all levels for a domain.
func (d *DomainModel) AllEnablers() []Enabler {
	var all []Enabler
	for _, level := range d.Levels {
		for _, e := range level.Enablers {
			all = append(all, e)
		}
	}
	return all
}

// CriteriaForLevel returns criteria for a specific level.
func (d *DomainModel) CriteriaForLevel(level int) []Criterion {
	for _, l := range d.Levels {
		if l.Level == level {
			return l.Criteria
		}
	}
	return nil
}

// EnablersForLevel returns enablers for a specific level.
func (d *DomainModel) EnablersForLevel(level int) []Enabler {
	for _, l := range d.Levels {
		if l.Level == level {
			return l.Enablers
		}
	}
	return nil
}
