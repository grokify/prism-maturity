// Package prism provides the PRISM (Platform for Reliability, Improvement, and Strategic Maturity)
// framework for COO-level organizational health monitoring combining SLOs, maturity modeling, and OKRs.
package prism

import "time"

// PRISMDocument represents the top-level PRISM document.
type PRISMDocument struct {
	Schema      string         `json:"$schema,omitempty"`
	Metadata    *Metadata      `json:"metadata,omitempty"`
	Domains     []DomainDef    `json:"domains,omitempty"`
	Layers      []LayerDef     `json:"layers,omitempty"`
	Teams       []Team         `json:"teams,omitempty"`
	Services    []Service      `json:"services,omitempty"`
	Metrics     []Metric       `json:"metrics"`
	Maturity    *MaturityModel `json:"maturity,omitempty"`
	OKRs        []OKRMapping   `json:"okrs,omitempty"`
	Initiatives []Initiative   `json:"initiatives,omitempty"`

	// Goal-driven Maturity Roadmap (FEAT_MATURITYROADMAP)
	Goals   []Goal         `json:"goals,omitempty"`
	Phases  []Phase        `json:"phases,omitempty"`
	Roadmap *RoadmapConfig `json:"roadmap,omitempty"`

	// Temporal State Tracking (REFACTOR_MATURITY_STATE)
	// Reference to the maturity model this document tracks against
	MaturityModelRef string `json:"maturityModelRef,omitempty"`

	// Standard SLO windows used in this document
	SLOWindows []string `json:"sloWindows,omitempty"`

	// SLI state with temporal tracking (past, present, future)
	SLIState SLIStateMap `json:"sliState,omitempty"`

	// Maturity level state per domain
	MaturityState MaturityStateMap `json:"maturityState,omitempty"`

	// Enabler/initiative state
	EnablerState EnablerStateMap `json:"enablerState,omitempty"`
}

// Metadata contains document-level metadata.
type Metadata struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
	Author      string `json:"author,omitempty"`
	Created     string `json:"created,omitempty"`
	Updated     string `json:"updated,omitempty"`
}

// DomainDef defines a PRISM domain (security or operations).
type DomainDef struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Weight      float64 `json:"weight,omitempty"`
}

// Metric represents a PRISM metric with SLO, maturity, and framework mappings.
type Metric struct {
	// Core identity
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// PRISM classification
	Domain          string `json:"domain"`
	Stage           string `json:"stage"`
	Category        string `json:"category"`
	Layer           string `json:"layer,omitempty"`           // code, infra, runtime
	QualityVertical string `json:"qualityVertical,omitempty"` // ISO 25010: functional, reliability, performance, security, usability, maintainability

	// Measurement
	MetricType     string  `json:"metricType"`
	TrendDirection string  `json:"trendDirection,omitempty"`
	Unit           string  `json:"unit,omitempty"`
	Baseline       float64 `json:"baseline"`
	Current        float64 `json:"current"`
	Target         float64 `json:"target"`

	// SLI/SLO
	SLI *SLI `json:"sli,omitempty"`
	SLO *SLO `json:"slo,omitempty"`

	// Thresholds & Status
	Thresholds *Thresholds `json:"thresholds,omitempty"`
	Status     string      `json:"status,omitempty"`

	// Maturity mapping
	MaturityMapping *MaturityMapping `json:"maturityMapping,omitempty"`

	// DMAIC mapping
	DMAIC *DMAICMapping `json:"dmaic,omitempty"`

	// Customer awareness
	CustomerAwareness *CustomerAwarenessConfig `json:"customerAwareness,omitempty"`

	// Framework mappings
	FrameworkMappings []FrameworkMapping `json:"frameworkMappings,omitempty"`

	// Ownership
	Owner      string `json:"owner,omitempty"`
	DataSource string `json:"dataSource,omitempty"`
	ServiceID  string `json:"serviceId,omitempty"` // Associated service

	// History
	DataPoints []DataPoint `json:"dataPoints,omitempty"`
}

// SLI represents a Service Level Indicator.
type SLI struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Formula     string `json:"formula,omitempty"`
	SLIType     string `json:"sliType,omitempty"` // Observability type: availability, latency, error_rate, etc.
}

// IsGoldenSignal returns true if the SLI type is part of Google's Golden Signals.
func (s *SLI) IsGoldenSignal() bool {
	if s == nil || s.SLIType == "" {
		return false
	}
	for _, t := range GoldenSignalsSLITypes() {
		if t == s.SLIType {
			return true
		}
	}
	return false
}

// IsRED returns true if the SLI type is part of the RED methodology.
func (s *SLI) IsRED() bool {
	if s == nil || s.SLIType == "" {
		return false
	}
	for _, t := range REDSLITypes() {
		if t == s.SLIType {
			return true
		}
	}
	return false
}

// IsUSE returns true if the SLI type is part of the USE methodology.
func (s *SLI) IsUSE() bool {
	if s == nil || s.SLIType == "" {
		return false
	}
	for _, t := range USESLITypes() {
		if t == s.SLIType {
			return true
		}
	}
	return false
}

// Methodologies returns all methodologies that include this SLI's type.
func (s *SLI) Methodologies() []string {
	if s == nil || s.SLIType == "" {
		return nil
	}
	return MethodologiesForSLIType(s.SLIType)
}

// SLO represents a Service Level Objective.
// SLOs can be quantitative (numeric targets) or qualitative (tracked states).
type SLO struct {
	// Identity
	ID   string `json:"id,omitempty"`   // Unique identifier for the SLO
	Name string `json:"name,omitempty"` // Human-readable name

	// Type distinguishes quantitative vs qualitative SLOs.
	// Quantitative SLOs have numeric targets (>=99.9%).
	// Qualitative SLOs track binary states (tracked, implemented, defined).
	Type string `json:"type,omitempty"` // "quantitative" (default), "qualitative"

	// Quantitative fields
	Target     string      `json:"target"`             // Display string (e.g., ">=99.9%")
	Operator   string      `json:"operator,omitempty"` // Machine-readable: "gte", "lte", "eq", "gt", "lt", "exists"
	Value      float64     `json:"value,omitempty"`    // Numeric target value
	Window     string      `json:"window,omitempty"`   // "7d", "30d", "90d"
	Thresholds *Thresholds `json:"thresholds,omitempty"`

	// Qualitative fields
	Status string `json:"status,omitempty"` // For qualitative: "tracked", "implemented", "defined", "documented", "not_tracked"

	// Framework mappings - maps this SLO to compliance framework controls
	FrameworkMappings []FrameworkMapping `json:"frameworkMappings,omitempty"`
}

// IsQualitative returns true if this is a qualitative SLO.
func (s *SLO) IsQualitative() bool {
	return s.Type == SLOTypeQualitative || s.Operator == SLOOperatorExists
}

// IsMet returns whether the SLO is met.
// For qualitative SLOs, checks if status indicates compliance.
// For quantitative SLOs, compares current value against target.
func (s *SLO) IsMet(current float64) bool {
	if s.IsQualitative() {
		return s.IsQualitativeStatusMet()
	}
	return s.isQuantitativeMet(current)
}

// IsQualitativeStatusMet returns whether a qualitative SLO status indicates compliance.
func (s *SLO) IsQualitativeStatusMet() bool {
	switch s.Status {
	case QualitativeStatusTracked, QualitativeStatusImplemented,
		QualitativeStatusDefined, QualitativeStatusDocumented,
		QualitativeStatusCompliant, QualitativeStatusEnabled:
		return true
	default:
		return false
	}
}

func (s *SLO) isQuantitativeMet(current float64) bool {
	switch s.Operator {
	case SLOOperatorGTE:
		return current >= s.Value
	case SLOOperatorLTE:
		return current <= s.Value
	case SLOOperatorEQ:
		return current == s.Value
	case SLOOperatorGT:
		return current > s.Value
	case SLOOperatorLT:
		return current < s.Value
	default:
		return false
	}
}

// SLO type constants.
const (
	SLOTypeQuantitative = "quantitative" // Default: numeric comparison
	SLOTypeQualitative  = "qualitative"  // Binary state tracking
)

// SLO operator constants.
const (
	SLOOperatorGTE    = "gte"    // Greater than or equal
	SLOOperatorLTE    = "lte"    // Less than or equal
	SLOOperatorEQ     = "eq"     // Equal
	SLOOperatorGT     = "gt"     // Greater than
	SLOOperatorLT     = "lt"     // Less than
	SLOOperatorExists = "exists" // Qualitative: metric exists/is tracked
)

// Qualitative status constants.
const (
	QualitativeStatusTracked     = "tracked"     // Metric is being tracked
	QualitativeStatusImplemented = "implemented" // Control/feature is implemented
	QualitativeStatusDefined     = "defined"     // Process/policy is defined
	QualitativeStatusDocumented  = "documented"  // Documentation exists
	QualitativeStatusCompliant   = "compliant"   // Meets compliance requirement
	QualitativeStatusEnabled     = "enabled"     // Feature/capability is enabled
	QualitativeStatusNotTracked  = "not_tracked" // Not yet being tracked
	QualitativeStatusPartial     = "partial"     // Partially implemented
	QualitativeStatusPlanned     = "planned"     // Planned but not started
)

// AllQualitativeStatuses returns all valid qualitative status values.
func AllQualitativeStatuses() []string {
	return []string{
		QualitativeStatusTracked,
		QualitativeStatusImplemented,
		QualitativeStatusDefined,
		QualitativeStatusDocumented,
		QualitativeStatusCompliant,
		QualitativeStatusEnabled,
		QualitativeStatusNotTracked,
		QualitativeStatusPartial,
		QualitativeStatusPlanned,
	}
}

// Thresholds defines threshold values for status calculation.
type Thresholds struct {
	Green  float64 `json:"green"`
	Yellow float64 `json:"yellow"`
	Red    float64 `json:"red"`
}

// MaturityMapping maps metric values to maturity levels.
type MaturityMapping struct {
	Level1 string `json:"level1,omitempty"`
	Level2 string `json:"level2,omitempty"`
	Level3 string `json:"level3,omitempty"`
	Level4 string `json:"level4,omitempty"`
	Level5 string `json:"level5,omitempty"`
}

// DMAICMapping maps the metric to DMAIC phases.
type DMAICMapping struct {
	Define  string `json:"define,omitempty"`
	Measure string `json:"measure,omitempty"`
	Analyze string `json:"analyze,omitempty"`
	Improve string `json:"improve,omitempty"`
	Control string `json:"control,omitempty"`
}

// FrameworkMapping maps a metric or SLO to an external framework reference.
// Supports compliance frameworks like NIST CSF, NIST 800-53, FedRAMP, etc.
type FrameworkMapping struct {
	Framework   string `json:"framework"`             // Framework identifier (e.g., "NIST_CSF_2", "NIST_800_53")
	Reference   string `json:"reference"`             // Control or function ID (e.g., "PR.DS-1", "AC-2")
	Name        string `json:"name,omitempty"`        // Human-readable control name
	Description string `json:"description,omitempty"` // Control description
	Baseline    string `json:"baseline,omitempty"`    // Required baseline level (e.g., "high", "moderate", "low" for FedRAMP)
	Version     string `json:"version,omitempty"`     // Framework version (e.g., "2.0", "Rev 5")
}

// FrameworkRequirement specifies a framework control requirement with satisfaction criteria.
type FrameworkRequirement struct {
	Framework      string   `json:"framework"`                // Framework identifier
	ControlID      string   `json:"controlId"`                // Control ID (e.g., "AC-2", "PR.DS-1")
	ControlName    string   `json:"controlName,omitempty"`    // Human-readable name
	Baseline       string   `json:"baseline,omitempty"`       // Required baseline (high/moderate/low)
	MetricIDs      []string `json:"metricIds,omitempty"`      // Metrics that satisfy this control
	SLOIDs         []string `json:"sloIds,omitempty"`         // SLOs that satisfy this control
	Status         string   `json:"status,omitempty"`         // implemented, partial, planned, not_applicable
	Evidence       string   `json:"evidence,omitempty"`       // Evidence or documentation reference
	Implementation string   `json:"implementation,omitempty"` // Implementation description
}

// DataPoint represents a historical measurement.
type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Note      string    `json:"note,omitempty"`
}

// OKRMapping represents alignment between metrics and OKRs.
type OKRMapping struct {
	ObjectiveID   string   `json:"objectiveId,omitempty"`
	ObjectiveName string   `json:"objectiveName"`
	KeyResultID   string   `json:"keyResultId,omitempty"`
	KeyResultName string   `json:"keyResultName,omitempty"`
	MetricIDs     []string `json:"metricIds,omitempty"`
}

// Initiative represents an improvement initiative.
type Initiative struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name"`
	Description    string   `json:"description,omitempty"`
	Status         string   `json:"status,omitempty"`
	Priority       int      `json:"priority,omitempty"`
	MetricIDs      []string `json:"metricIds,omitempty"`
	Owner          string   `json:"owner,omitempty"`
	Team           string   `json:"team,omitempty"`
	DependentTeams []string `json:"dependentTeams,omitempty"`
	StartDate      string   `json:"startDate,omitempty"`
	EndDate        string   `json:"endDate,omitempty"`

	// Goal and Phase linkage (FEAT_MATURITYROADMAP)
	GoalIDs              []string          `json:"goalIds,omitempty"`
	PhaseID              string            `json:"phaseId,omitempty"`
	DevCompletionPercent float64           `json:"devCompletionPercent,omitempty"`
	DeploymentStatus     *DeploymentStatus `json:"deploymentStatus,omitempty"`

	// Service linkage
	ServiceID string `json:"serviceId,omitempty"` // Associated service
}

// Initiative status constants.
const (
	InitiativeStatusPlanned    = "planned"
	InitiativeStatusNotStarted = "not_started"
	InitiativeStatusInProgress = "in_progress"
	InitiativeStatusCompleted  = "completed"
	InitiativeStatusCancelled  = "cancelled"
)

// AllInitiativeStatuses returns all valid initiative status values.
func AllInitiativeStatuses() []string {
	return []string{
		InitiativeStatusPlanned,
		InitiativeStatusNotStarted,
		InitiativeStatusInProgress,
		InitiativeStatusCompleted,
		InitiativeStatusCancelled,
	}
}

// DeploymentStatus tracks customer adoption for an initiative.
type DeploymentStatus struct {
	Status            string  `json:"status"`                      // not_started, in_progress, completed
	TotalCustomers    int     `json:"totalCustomers,omitempty"`    // Total customers to deploy to
	DeployedCustomers int     `json:"deployedCustomers,omitempty"` // Customers deployed
	AdoptionPercent   float64 `json:"adoptionPercent,omitempty"`   // Calculated adoption percentage
}

// CalculateAdoptionPercent calculates and updates the adoption percentage.
func (ds *DeploymentStatus) CalculateAdoptionPercent() float64 {
	if ds == nil || ds.TotalCustomers == 0 {
		return 0
	}
	ds.AdoptionPercent = float64(ds.DeployedCustomers) / float64(ds.TotalCustomers) * 100
	return ds.AdoptionPercent
}

// IsDevComplete returns whether the initiative is development complete.
func (i *Initiative) IsDevComplete() bool {
	return i.Status == InitiativeStatusCompleted || i.DevCompletionPercent >= 100
}

// IsFullyDeployed returns whether the initiative is fully deployed to all customers.
func (i *Initiative) IsFullyDeployed() bool {
	if i.DeploymentStatus == nil {
		return false
	}
	return i.DeploymentStatus.Status == InitiativeStatusCompleted ||
		(i.DeploymentStatus.TotalCustomers > 0 && i.DeploymentStatus.DeployedCustomers >= i.DeploymentStatus.TotalCustomers)
}

// CalculateStatus computes the status based on current value and thresholds.
// For higher_better trends: value >= green threshold = Green, etc.
// For lower_better trends: value <= green threshold = Green, etc.
func (m *Metric) CalculateStatus() string {
	if m.Thresholds == nil {
		return ""
	}

	current := m.Current
	th := m.Thresholds

	switch m.TrendDirection {
	case TrendHigherBetter:
		if current >= th.Green {
			return StatusGreen
		} else if current >= th.Yellow {
			return StatusYellow
		}
		return StatusRed
	case TrendLowerBetter:
		if current <= th.Green {
			return StatusGreen
		} else if current <= th.Yellow {
			return StatusYellow
		}
		return StatusRed
	case TrendTargetValue:
		// For target value, check distance from target
		diff := abs(current - m.Target)
		greenRange := abs(th.Green - m.Target)
		yellowRange := abs(th.Yellow - m.Target)
		if diff <= greenRange {
			return StatusGreen
		} else if diff <= yellowRange {
			return StatusYellow
		}
		return StatusRed
	default:
		// Default to higher_better behavior
		if current >= th.Green {
			return StatusGreen
		} else if current >= th.Yellow {
			return StatusYellow
		}
		return StatusRed
	}
}

// ProgressToTarget returns the progress as a ratio (0.0-1.0) toward the target.
func (m *Metric) ProgressToTarget() float64 {
	if m.Target == m.Baseline {
		if m.Current >= m.Target {
			return 1.0
		}
		return 0.0
	}

	progress := (m.Current - m.Baseline) / (m.Target - m.Baseline)
	if progress < 0 {
		return 0.0
	}
	if progress > 1 {
		return 1.0
	}
	return progress
}

// MeetsSLO returns whether the metric's current value meets its SLO.
// Returns true if no SLO is defined or if Operator/Value are not set.
// Uses the structured Operator and Value fields for evaluation.
func (m *Metric) MeetsSLO() bool {
	if m.SLO == nil {
		return true // No SLO defined
	}
	if m.SLO.Operator == "" {
		return true // No machine-readable operator defined
	}

	current := m.Current
	target := m.SLO.Value

	switch m.SLO.Operator {
	case SLOOperatorGTE:
		return current >= target
	case SLOOperatorLTE:
		return current <= target
	case SLOOperatorGT:
		return current > target
	case SLOOperatorLT:
		return current < target
	case SLOOperatorEQ:
		return current == target
	default:
		return true // Unknown operator, assume met
	}
}

// GetMetricsByDomain returns all metrics for the specified domain.
func (doc *PRISMDocument) GetMetricsByDomain(domain string) []Metric {
	var result []Metric
	for _, m := range doc.Metrics {
		if m.Domain == domain {
			result = append(result, m)
		}
	}
	return result
}

// GetMetricsByStage returns all metrics for the specified stage.
func (doc *PRISMDocument) GetMetricsByStage(stage string) []Metric {
	var result []Metric
	for _, m := range doc.Metrics {
		if m.Stage == stage {
			result = append(result, m)
		}
	}
	return result
}

// GetMetricsByCategory returns all metrics for the specified category.
func (doc *PRISMDocument) GetMetricsByCategory(category string) []Metric {
	var result []Metric
	for _, m := range doc.Metrics {
		if m.Category == category {
			result = append(result, m)
		}
	}
	return result
}

// GetMetricByID returns a metric by its ID.
func (doc *PRISMDocument) GetMetricByID(id string) *Metric {
	for i := range doc.Metrics {
		if doc.Metrics[i].ID == id {
			return &doc.Metrics[i]
		}
	}
	return nil
}

// GetGoalByID returns a goal by its ID.
func (doc *PRISMDocument) GetGoalByID(id string) *Goal {
	for i := range doc.Goals {
		if doc.Goals[i].ID == id {
			return &doc.Goals[i]
		}
	}
	return nil
}

// GetPhaseByID returns a phase by its ID.
func (doc *PRISMDocument) GetPhaseByID(id string) *Phase {
	for i := range doc.Phases {
		if doc.Phases[i].ID == id {
			return &doc.Phases[i]
		}
	}
	return nil
}

// GetInitiativeByID returns an initiative by its ID.
func (doc *PRISMDocument) GetInitiativeByID(id string) *Initiative {
	for i := range doc.Initiatives {
		if doc.Initiatives[i].ID == id {
			return &doc.Initiatives[i]
		}
	}
	return nil
}

// GetInitiativesForGoal returns all initiatives linked to the specified goal.
func (doc *PRISMDocument) GetInitiativesForGoal(goalID string) []Initiative {
	var result []Initiative
	for _, init := range doc.Initiatives {
		for _, gid := range init.GoalIDs {
			if gid == goalID {
				result = append(result, init)
				break
			}
		}
	}
	return result
}

// GetInitiativesForPhase returns all initiatives in the specified phase.
func (doc *PRISMDocument) GetInitiativesForPhase(phaseID string) []Initiative {
	var result []Initiative
	for _, init := range doc.Initiatives {
		if init.PhaseID == phaseID {
			result = append(result, init)
		}
	}
	return result
}

// GetLayerByID returns a layer definition by its ID.
func (doc *PRISMDocument) GetLayerByID(id string) *LayerDef {
	for i := range doc.Layers {
		if doc.Layers[i].ID == id {
			return &doc.Layers[i]
		}
	}
	return nil
}

// GetMetricsByLayer returns all metrics for the specified layer.
func (doc *PRISMDocument) GetMetricsByLayer(layer string) []Metric {
	var result []Metric
	for _, m := range doc.Metrics {
		if m.Layer == layer {
			result = append(result, m)
		}
	}
	return result
}

// GetTeamByID returns a team by its ID.
func (doc *PRISMDocument) GetTeamByID(id string) *Team {
	for i := range doc.Teams {
		if doc.Teams[i].ID == id {
			return &doc.Teams[i]
		}
	}
	return nil
}

// GetServiceByID returns a service by its ID.
func (doc *PRISMDocument) GetServiceByID(id string) *Service {
	for i := range doc.Services {
		if doc.Services[i].ID == id {
			return &doc.Services[i]
		}
	}
	return nil
}

// GetMetricsByService returns all metrics for the specified service.
func (doc *PRISMDocument) GetMetricsByService(serviceID string) []Metric {
	var result []Metric
	for _, m := range doc.Metrics {
		if m.ServiceID == serviceID {
			result = append(result, m)
		}
	}
	return result
}

// abs returns the absolute value of x.
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// =============================================================================
// Temporal State Tracking (REFACTOR_MATURITY_STATE)
// =============================================================================

// SLOWindow represents a standard SLO measurement window.
const (
	SLOWindow7Days     = "7d"
	SLOWindow30Days    = "30d"
	SLOWindow90Days    = "90d"
	SLOWindowQuarterly = "quarterly"
	SLOWindowAnnual    = "annual"
)

// AllSLOWindows returns all standard SLO windows.
func AllSLOWindows() []string {
	return []string{
		SLOWindow7Days,
		SLOWindow30Days,
		SLOWindow90Days,
		SLOWindowQuarterly,
		SLOWindowAnnual,
	}
}

// SLIStateMap holds state for all SLIs in a PRISM document.
type SLIStateMap map[string]*SLIState

// SLIState tracks temporal state for a single SLI.
type SLIState struct {
	SLIID            string                  `json:"sliId"`                      // Reference to maturity model SLI
	QualitativeState string                  `json:"qualitativeState,omitempty"` // Current qualitative state (e.g., "tracked", "measured")
	Windows          map[string]*WindowState `json:"windows,omitempty"`          // State by window (7d, 30d, 90d, quarterly, annual)
	History          []HistoricalValue       `json:"history,omitempty"`          // Historical values over time
	Targets          map[string]*TargetValue `json:"targets,omitempty"`          // Future targets by period (e.g., "Q2_2026")
}

// WindowState represents the current value for a specific SLO window.
type WindowState struct {
	Value     float64   `json:"value"`               // Current measured value
	Target    float64   `json:"target,omitempty"`    // Target value for this window
	Met       bool      `json:"met,omitempty"`       // Whether target is met
	Timestamp time.Time `json:"timestamp,omitempty"` // When this was measured
}

// HistoricalValue represents a point-in-time measurement.
type HistoricalValue struct {
	Window    string    `json:"window"`         // Which SLO window (7d, 30d, etc.)
	Value     float64   `json:"value"`          // Measured value
	Timestamp time.Time `json:"timestamp"`      // When measured
	Note      string    `json:"note,omitempty"` // Optional note
}

// TargetValue represents a future target.
type TargetValue struct {
	Value         float64 `json:"value"`                   // Target value
	MaturityLevel int     `json:"maturityLevel,omitempty"` // Target maturity level (1-5)
	TargetDate    string  `json:"targetDate,omitempty"`    // Target date (ISO 8601)
}

// MaturityStateMap holds maturity state for all domains.
type MaturityStateMap map[string]*DomainMaturityState

// DomainMaturityState tracks maturity level progression for a domain.
type DomainMaturityState struct {
	DomainID string               `json:"domainId"`          // Reference to maturity model domain
	Current  *MaturityLevelState  `json:"current"`           // Current achieved level
	Target   *MaturityLevelTarget `json:"target,omitempty"`  // Target level
	History  []MaturityLevelState `json:"history,omitempty"` // Level progression history
}

// MaturityLevelState represents a maturity level at a point in time.
type MaturityLevelState struct {
	Level      int    `json:"level"`                // Maturity level (1-5)
	AchievedAt string `json:"achievedAt,omitempty"` // When this level was achieved (ISO 8601)
	AssessedBy string `json:"assessedBy,omitempty"` // Who performed the assessment
	Note       string `json:"note,omitempty"`       // Assessment notes
}

// MaturityLevelTarget represents a target maturity level.
type MaturityLevelTarget struct {
	Level      int    `json:"level"`                // Target maturity level (1-5)
	TargetDate string `json:"targetDate,omitempty"` // Target date (ISO 8601)
	Rationale  string `json:"rationale,omitempty"`  // Why this target was chosen
}

// EnablerStateMap holds state for all enablers.
type EnablerStateMap map[string]*EnablerState

// EnablerState tracks progress on an enabler (project/capability).
type EnablerState struct {
	EnablerID   string  `json:"enablerId"`             // Reference to maturity model enabler
	Status      string  `json:"status"`                // not_started, in_progress, completed, blocked
	Progress    float64 `json:"progress,omitempty"`    // Completion percentage (0-100)
	StartedAt   string  `json:"startedAt,omitempty"`   // When work started
	CompletedAt string  `json:"completedAt,omitempty"` // When completed
	Owner       string  `json:"owner,omitempty"`       // Who owns this enabler
	Note        string  `json:"note,omitempty"`        // Status notes
}

// EnablerStatus constants.
const (
	EnablerStatusNotStarted = "not_started"
	EnablerStatusInProgress = "in_progress"
	EnablerStatusCompleted  = "completed"
	EnablerStatusBlocked    = "blocked"
)

// QualitativeStateDefinition defines a qualitative state for an SLI.
type QualitativeStateDefinition struct {
	ID          string `json:"id"`                    // State identifier (e.g., "tracked")
	Label       string `json:"label"`                 // Human-readable label
	Description string `json:"description,omitempty"` // What this state means
	Order       int    `json:"order"`                 // Progression order (0 = lowest)
}

// StandardQualitativeStates returns the default progression of qualitative states.
func StandardQualitativeStates() []QualitativeStateDefinition {
	return []QualitativeStateDefinition{
		{ID: "none", Label: "Not Tracked", Description: "Metric is not being tracked", Order: 0},
		{ID: "adhoc", Label: "Ad-hoc", Description: "Tracked inconsistently or manually", Order: 1},
		{ID: "tracked", Label: "Tracked", Description: "Tracked regularly but no SLO", Order: 2},
		{ID: "measured", Label: "Measured", Description: "Measured with defined SLO", Order: 3},
		{ID: "alerting", Label: "SLO + Alerting", Description: "SLO with automated alerting", Order: 4},
		{ID: "optimized", Label: "Optimized", Description: "Continuously optimized with automation", Order: 5},
	}
}

// CompareQualitativeStates compares two qualitative states.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
func CompareQualitativeStates(a, b string) int {
	states := StandardQualitativeStates()
	orderA, orderB := -1, -1
	for _, s := range states {
		if s.ID == a {
			orderA = s.Order
		}
		if s.ID == b {
			orderB = s.Order
		}
	}
	if orderA < orderB {
		return -1
	}
	if orderA > orderB {
		return 1
	}
	return 0
}

// MeetsQualitativeTarget returns true if current state meets or exceeds target.
func MeetsQualitativeTarget(current, target string) bool {
	return CompareQualitativeStates(current, target) >= 0
}

// MaturityStateDocument is a PRISM Maturity State document that tracks
// current state against a maturity model. This is the top-level document
// type for prism-maturity-state.schema.json.
type MaturityStateDocument struct {
	// Schema is the JSON Schema reference.
	Schema string `json:"$schema,omitempty"`

	// Metadata contains document identification and references.
	Metadata MaturityStateMetadata `json:"metadata"`

	// SLOWindows defines which temporal windows are tracked in this state document.
	SLOWindows []string `json:"sloWindows,omitempty"`

	// SLIState tracks current values for each SLI.
	SLIState SLIStateMap `json:"sliState,omitempty"`

	// MaturityState tracks maturity level progression for each domain.
	MaturityState MaturityStateMap `json:"maturityState,omitempty"`

	// EnablerState tracks progress on enablers/initiatives.
	EnablerState EnablerStateMap `json:"enablerState,omitempty"`
}

// MaturityStateMetadata contains identification and reference information.
type MaturityStateMetadata struct {
	// Name is the document name.
	Name string `json:"name"`

	// Description is an optional description.
	Description string `json:"description,omitempty"`

	// Version is the document version.
	Version string `json:"version,omitempty"`

	// MaturityModelRef is a reference to the maturity model this state tracks against.
	MaturityModelRef string `json:"maturityModelRef,omitempty"`

	// AssessedAt is when this state was assessed (ISO 8601).
	AssessedAt string `json:"assessedAt,omitempty"`

	// AssessedBy is who performed the assessment.
	AssessedBy string `json:"assessedBy,omitempty"`

	// Organization is the organization name.
	Organization string `json:"organization,omitempty"`
}

// MaturityPlanDocument is a PRISM Maturity Plan document that defines
// goals, phases, and initiatives for achieving maturity targets.
// This is the top-level document type for prism-maturity-plan.schema.json.
// Note: This is equivalent to the existing PRISMDocument type, which will
// be aliased to MaturityPlanDocument in a future version.
type MaturityPlanDocument = PRISMDocument
