package prism

import (
	"strconv"

	core "github.com/grokify/prism-core"
)

// Domain constants represent the three primary domains in PRISM.
// Note: prism-core has 10 domains; prism-intelligence uses these 3 primary ones.
const (
	DomainSecurity   = core.DomainSecurity
	DomainOperations = core.DomainOperations
	DomainQuality    = core.DomainQuality
)

// AllDomains returns all valid domain values for prism-intelligence.
// Returns the 3 primary domains used in this module.
func AllDomains() []string {
	return []string{DomainSecurity, DomainOperations, DomainQuality}
}

// ValidDomain checks if a domain value is valid.
func ValidDomain(domain string) bool {
	return core.ValidDomain(domain)
}

// Layer constants imported from prism-core.
const (
	LayerRequirements = core.LayerRequirements
	LayerCode         = core.LayerCode
	LayerInfra        = core.LayerInfra
	LayerRuntime      = core.LayerRuntime
	LayerAdoption     = core.LayerAdoption
	LayerSupport      = core.LayerSupport
)

// AllLayers returns all valid layer values in value stream order.
func AllLayers() []string {
	return core.AllLayers()
}

// ValidLayer checks if a layer value is valid.
func ValidLayer(layer string) bool {
	return core.ValidLayer(layer)
}

// QualityVertical constants based on ISO 25010 quality characteristics.
// These are prism-intelligence-specific.
const (
	QualityVerticalFunctional      = "functional"
	QualityVerticalReliability     = "reliability"
	QualityVerticalPerformance     = "performance"
	QualityVerticalSecurity        = "security"
	QualityVerticalUsability       = "usability"
	QualityVerticalMaintainability = "maintainability"
)

// AllQualityVerticals returns all valid ISO 25010 quality vertical values.
func AllQualityVerticals() []string {
	return []string{
		QualityVerticalFunctional,
		QualityVerticalReliability,
		QualityVerticalPerformance,
		QualityVerticalSecurity,
		QualityVerticalUsability,
		QualityVerticalMaintainability,
	}
}

// Lifecycle stage constants imported from prism-core.
const (
	StageDesign   = core.StageDesign
	StageBuild    = core.StageBuild
	StageTest     = core.StageTest
	StageRuntime  = core.StageRuntime
	StageResponse = core.StageResponse
)

// AllStages returns all valid stage values.
func AllStages() []string {
	return core.AllStages()
}

// ValidStage checks if a stage value is valid.
func ValidStage(stage string) bool {
	return core.ValidStage(stage)
}

// Category constants represent metric categories.
// These are prism-intelligence-specific.
const (
	CategoryPrevention  = "prevention"
	CategoryDetection   = "detection"
	CategoryResponse    = "response"
	CategoryReliability = "reliability"
	CategoryEfficiency  = "efficiency"
	CategoryQuality     = "quality"
)

// AllCategories returns all valid category values.
func AllCategories() []string {
	return []string{
		CategoryPrevention,
		CategoryDetection,
		CategoryResponse,
		CategoryReliability,
		CategoryEfficiency,
		CategoryQuality,
	}
}

// Maturity level constants imported from prism-core.
const (
	MaturityLevel1 = core.MaturityLevel1 // Reactive
	MaturityLevel2 = core.MaturityLevel2 // Basic
	MaturityLevel3 = core.MaturityLevel3 // Defined
	MaturityLevel4 = core.MaturityLevel4 // Managed
	MaturityLevel5 = core.MaturityLevel5 // Optimizing
)

// MaturityLevelName returns the name for a maturity level.
func MaturityLevelName(level int) string {
	return core.MaturityLevelName(level)
}

// MaturityLevelDescription returns a description for a maturity level.
func MaturityLevelDescription(level int) string {
	return core.MaturityLevelDescription(level)
}

// ValidMaturityLevel checks if a maturity level is valid (1-5).
func ValidMaturityLevel(level int) bool {
	return core.ValidMaturityLevel(level)
}

// Customer awareness state constants.
// These are prism-intelligence-specific.
const (
	AwarenessUnaware          = "unaware"
	AwarenessAwareNotActing   = "aware_not_remediating"
	AwarenessAwareRemediating = "aware_remediating"
	AwarenessAwareRemediated  = "aware_remediated"
)

// AllAwarenessStates returns all valid awareness state values.
func AllAwarenessStates() []string {
	return []string{
		AwarenessUnaware,
		AwarenessAwareNotActing,
		AwarenessAwareRemediating,
		AwarenessAwareRemediated,
	}
}

// Framework constants imported from prism-core.
const (
	// NIST Frameworks
	FrameworkNISTCSF    = core.FrameworkNISTCSF
	FrameworkNISTCSF2   = core.FrameworkNISTCSF2
	FrameworkNIST80053  = core.FrameworkNIST80053
	FrameworkNISTRMF    = core.FrameworkNISTRMF
	FrameworkNISTAIRMF  = core.FrameworkNISTAIRMF
	FrameworkNIST800171 = core.FrameworkNIST800171

	// FedRAMP (uses NIST 800-53 controls)
	FrameworkFEDRAMP     = "FEDRAMP" // FedRAMP (general) - prism-intelligence-specific
	FrameworkFEDRAMPHigh = core.FrameworkFedRAMPHigh
	FrameworkFEDRAMPMod  = core.FrameworkFedRAMPMod
	FrameworkFEDRAMPLow  = core.FrameworkFedRAMPLow

	// Other Security Frameworks
	FrameworkMITREATTACK = core.FrameworkMITREATTACK
	FrameworkCISControls = core.FrameworkCISControls
	FrameworkSOC2        = core.FrameworkSOC2
	FrameworkISO27001    = core.FrameworkISO27001

	// Engineering Frameworks
	FrameworkDORA = core.FrameworkDORA
	FrameworkSRE  = core.FrameworkSRE
)

// Framework baseline/impact levels for NIST 800-53 and FedRAMP.
const (
	BaselineHigh     = "high"
	BaselineModerate = "moderate"
	BaselineLow      = "low"
)

// AllFrameworks returns all valid framework values.
func AllFrameworks() []string {
	return []string{
		FrameworkNISTCSF,
		FrameworkNISTCSF2,
		FrameworkNIST80053,
		FrameworkNISTRMF,
		FrameworkNISTAIRMF,
		FrameworkNIST800171,
		FrameworkFEDRAMP,
		FrameworkFEDRAMPHigh,
		FrameworkFEDRAMPMod,
		FrameworkFEDRAMPLow,
		FrameworkMITREATTACK,
		FrameworkCISControls,
		FrameworkSOC2,
		FrameworkISO27001,
		FrameworkDORA,
		FrameworkSRE,
	}
}

// NISTFrameworks returns NIST-specific frameworks.
func NISTFrameworks() []string {
	return core.NISTFrameworks()
}

// ComplianceFrameworks returns compliance-focused frameworks.
func ComplianceFrameworks() []string {
	return []string{
		FrameworkNISTCSF,
		FrameworkNISTCSF2,
		FrameworkNIST80053,
		FrameworkFEDRAMP,
		FrameworkFEDRAMPHigh,
		FrameworkFEDRAMPMod,
		FrameworkFEDRAMPLow,
		FrameworkSOC2,
		FrameworkISO27001,
		FrameworkCISControls,
	}
}

// ValidFramework checks if a framework value is valid.
func ValidFramework(framework string) bool {
	for _, f := range AllFrameworks() {
		if f == framework {
			return true
		}
	}
	return false
}

// FrameworkDisplayName returns a human-readable name for a framework.
func FrameworkDisplayName(framework string) string {
	return core.FrameworkDisplayName(framework)
}

// Metric type constants.
// These are prism-intelligence-specific.
const (
	MetricTypeCoverage     = "coverage"
	MetricTypeRate         = "rate"
	MetricTypeLatency      = "latency"
	MetricTypeRatio        = "ratio"
	MetricTypeCount        = "count"
	MetricTypeDistribution = "distribution"
	MetricTypeScore        = "score"
)

// AllMetricTypes returns all valid metric type values.
func AllMetricTypes() []string {
	return []string{
		MetricTypeCoverage,
		MetricTypeRate,
		MetricTypeLatency,
		MetricTypeRatio,
		MetricTypeCount,
		MetricTypeDistribution,
		MetricTypeScore,
	}
}

// Trend direction constants for threshold interpretation.
// Note: These are comparison semantics, distinct from prism-core's trend direction types.
const (
	TrendHigherBetter = "higher_better"
	TrendLowerBetter  = "lower_better"
	TrendTargetValue  = "target_value"
)

// AllTrendDirections returns all valid trend direction values.
func AllTrendDirections() []string {
	return []string{TrendHigherBetter, TrendLowerBetter, TrendTargetValue}
}

// Status constants for metric health.
const (
	StatusGreen  = "Green"
	StatusYellow = "Yellow"
	StatusRed    = "Red"
)

// AllStatuses returns all valid status values.
func AllStatuses() []string {
	return []string{StatusGreen, StatusYellow, StatusRed}
}

// SLO window constants.
const (
	Window7Days  = "7d"
	Window30Days = "30d"
	Window90Days = "90d"
)

// AllWindows returns all valid SLO window values.
func AllWindows() []string {
	return []string{Window7Days, Window30Days, Window90Days}
}

// SLI type constants classify SLIs by observability type.
// Note: prism-core has a similar set; these are prism-intelligence's SLI types.
const (
	SLITypeAvailability = core.SLITypeAvailability
	SLITypeLatency      = core.SLITypeLatency
	SLITypeErrorRate    = core.SLITypeErrorRate
	SLITypeThroughput   = core.SLITypeThroughput
	SLITypeSaturation   = "saturation"  // prism-intelligence-specific
	SLITypeUtilization  = "utilization" // prism-intelligence-specific
	SLITypeQuality      = "quality"     // prism-intelligence-specific
	SLITypeFreshness    = core.SLITypeFreshness
)

// AllSLITypes returns all valid SLI type values.
func AllSLITypes() []string {
	return []string{
		SLITypeAvailability,
		SLITypeLatency,
		SLITypeErrorRate,
		SLITypeThroughput,
		SLITypeSaturation,
		SLITypeUtilization,
		SLITypeQuality,
		SLITypeFreshness,
	}
}

// SLITypeDirection returns the default comparison direction for an SLI type.
func SLITypeDirection(sliType string) string {
	return core.SLITypeDirection(sliType)
}

// Methodology constants for standard observability methodologies.
const (
	MethodologyGoldenSignals = "GOLDEN_SIGNALS" // Google Golden Signals: latency, traffic, errors, saturation
	MethodologyRED           = "RED"            // Rate, Errors, Duration
	MethodologyUSE           = "USE"            // Utilization, Saturation, Errors
)

// AllMethodologies returns all valid methodology values.
func AllMethodologies() []string {
	return []string{
		MethodologyGoldenSignals,
		MethodologyRED,
		MethodologyUSE,
	}
}

// GoldenSignalsSLITypes returns the SLI types that map to Google's Golden Signals.
// Golden Signals: Latency, Traffic, Errors, Saturation.
func GoldenSignalsSLITypes() []string {
	return []string{
		SLITypeLatency,
		SLITypeThroughput,
		SLITypeErrorRate,
		SLITypeSaturation,
	}
}

// REDSLITypes returns the SLI types that map to the RED methodology.
// RED: Rate (throughput), Errors, Duration (latency).
func REDSLITypes() []string {
	return []string{
		SLITypeThroughput,
		SLITypeErrorRate,
		SLITypeLatency,
	}
}

// USESLITypes returns the SLI types that map to the USE methodology.
// USE: Utilization, Saturation, Errors.
func USESLITypes() []string {
	return []string{
		SLITypeUtilization,
		SLITypeSaturation,
		SLITypeErrorRate,
	}
}

// SLITypesForMethodology returns the SLI types associated with a methodology.
// Returns nil for unknown methodologies.
func SLITypesForMethodology(methodology string) []string {
	switch methodology {
	case MethodologyGoldenSignals:
		return GoldenSignalsSLITypes()
	case MethodologyRED:
		return REDSLITypes()
	case MethodologyUSE:
		return USESLITypes()
	default:
		return nil
	}
}

// MethodologiesForSLIType returns the methodologies that include the given SLI type.
func MethodologiesForSLIType(sliType string) []string {
	var result []string
	methodologyTypes := map[string][]string{
		MethodologyGoldenSignals: GoldenSignalsSLITypes(),
		MethodologyRED:           REDSLITypes(),
		MethodologyUSE:           USESLITypes(),
	}

	for methodology, types := range methodologyTypes {
		for _, t := range types {
			if t == sliType {
				result = append(result, methodology)
				break
			}
		}
	}
	return result
}

// Importance constants define static weights for categories, layers, and capabilities.
// These represent the inherent importance of "-ilities" (security, availability, etc.)
// and are used in conjunction with current state to calculate dynamic priority.
const (
	ImportanceCritical = "critical"
	ImportanceHigh     = "high"
	ImportanceMedium   = "medium"
	ImportanceLow      = "low"
)

// AllImportanceLevels returns all valid importance levels in descending order.
func AllImportanceLevels() []string {
	return []string{
		ImportanceCritical,
		ImportanceHigh,
		ImportanceMedium,
		ImportanceLow,
	}
}

// ImportanceWeight returns a numeric weight for the importance level.
// Higher weights indicate higher importance.
func ImportanceWeight(importance string) int {
	switch importance {
	case ImportanceCritical:
		return 4
	case ImportanceHigh:
		return 3
	case ImportanceMedium:
		return 2
	case ImportanceLow:
		return 1
	default:
		return 2 // Default to medium
	}
}

// Priority constants define dynamic priority levels based on current state.
// These are calculated by combining importance with maturity gap.
const (
	PriorityP0 = "P0" // Immediate action required
	PriorityP1 = "P1" // High priority
	PriorityP2 = "P2" // Medium priority
	PriorityP3 = "P3" // Low priority
)

// AllPriorityLevels returns all valid priority levels in descending order.
func AllPriorityLevels() []string {
	return []string{
		PriorityP0,
		PriorityP1,
		PriorityP2,
		PriorityP3,
	}
}

// DynamicPriorityWeight returns a numeric weight for the dynamic priority level (P0-P3).
// Higher weights indicate higher priority.
func DynamicPriorityWeight(priority string) int {
	switch priority {
	case PriorityP0:
		return 4
	case PriorityP1:
		return 3
	case PriorityP2:
		return 2
	case PriorityP3:
		return 1
	default:
		return 2 // Default to P2
	}
}

// CalculatePriority determines dynamic priority based on importance and maturity gap.
// importance: the static importance level (critical, high, medium, low)
// currentLevel: current maturity level (1-5)
// targetLevel: target maturity level (1-5)
// Returns P0-P3 based on the combination.
func CalculatePriority(importance string, currentLevel, targetLevel int) string {
	if currentLevel >= targetLevel {
		return PriorityP3 // Already at or above target
	}

	gap := targetLevel - currentLevel
	weight := ImportanceWeight(importance)

	// Priority score: importance weight * gap
	score := weight * gap

	switch {
	case score >= 8: // Critical with 2+ gap, or High with 3+ gap
		return PriorityP0
	case score >= 4: // High with 2 gap, or Medium with 2+ gap
		return PriorityP1
	case score >= 2: // Any importance with small gap
		return PriorityP2
	default:
		return PriorityP3
	}
}

// PriorityRationale returns a human-readable explanation for the calculated priority.
func PriorityRationale(importance string, currentLevel, targetLevel int) string {
	if currentLevel >= targetLevel {
		return "At or above target maturity level"
	}

	gap := targetLevel - currentLevel
	priority := CalculatePriority(importance, currentLevel, targetLevel)

	rationale := ""
	switch priority {
	case PriorityP0:
		rationale = "Immediate action required"
	case PriorityP1:
		rationale = "High priority improvement"
	case PriorityP2:
		rationale = "Scheduled improvement"
	case PriorityP3:
		rationale = "Low priority enhancement"
	}

	return rationale + ": " + importance + " importance with " + strconv.Itoa(gap) + "-level gap"
}
