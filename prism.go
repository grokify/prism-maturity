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
