package prism

// Domain constants represent the three primary domains in PRISM.
const (
	DomainSecurity   = "security"
	DomainOperations = "operations"
	DomainQuality    = "quality"
)

// AllDomains returns all valid domain values.
func AllDomains() []string {
	return []string{DomainSecurity, DomainOperations, DomainQuality}
}

// Layer constants represent value stream phases from ideation to support.
const (
	LayerRequirements = "requirements"
	LayerCode         = "code"
	LayerInfra        = "infra"
	LayerRuntime      = "runtime"
	LayerAdoption     = "adoption"
	LayerSupport      = "support"
)

// AllLayers returns all valid layer values in value stream order.
func AllLayers() []string {
	return []string{
		LayerRequirements,
		LayerCode,
		LayerInfra,
		LayerRuntime,
		LayerAdoption,
		LayerSupport,
	}
}

// QualityVertical constants based on ISO 25010 quality characteristics.
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

// Lifecycle stage constants represent stages in the software delivery lifecycle.
const (
	StageDesign   = "design"
	StageBuild    = "build"
	StageTest     = "test"
	StageRuntime  = "runtime"
	StageResponse = "response"
)

// AllStages returns all valid stage values.
func AllStages() []string {
	return []string{StageDesign, StageBuild, StageTest, StageRuntime, StageResponse}
}

// Category constants represent metric categories.
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

// Maturity level constants represent the 5-level maturity model.
const (
	MaturityLevel1 = 1 // Reactive
	MaturityLevel2 = 2 // Basic
	MaturityLevel3 = 3 // Defined
	MaturityLevel4 = 4 // Managed
	MaturityLevel5 = 5 // Optimizing
)

// MaturityLevelName returns the name for a maturity level.
func MaturityLevelName(level int) string {
	switch level {
	case MaturityLevel1:
		return "Reactive"
	case MaturityLevel2:
		return "Basic"
	case MaturityLevel3:
		return "Defined"
	case MaturityLevel4:
		return "Managed"
	case MaturityLevel5:
		return "Optimizing"
	default:
		return ""
	}
}

// Customer awareness state constants.
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

// Framework constants for external framework mappings.
const (
	// NIST Frameworks
	FrameworkNISTCSF    = "NIST_CSF"     // NIST Cybersecurity Framework (1.1)
	FrameworkNISTCSF2   = "NIST_CSF_2"   // NIST Cybersecurity Framework 2.0
	FrameworkNIST80053  = "NIST_800_53"  // NIST SP 800-53 (Security and Privacy Controls)
	FrameworkNISTRMF    = "NIST_RMF"     // NIST Risk Management Framework
	FrameworkNISTAIRMF  = "NIST_AI_RMF"  // NIST AI Risk Management Framework
	FrameworkNIST800171 = "NIST_800_171" // NIST SP 800-171 (CUI Protection)

	// FedRAMP (uses NIST 800-53 controls)
	FrameworkFEDRAMP     = "FEDRAMP"      // FedRAMP (general)
	FrameworkFEDRAMPHigh = "FEDRAMP_HIGH" // FedRAMP High baseline
	FrameworkFEDRAMPMod  = "FEDRAMP_MOD"  // FedRAMP Moderate baseline
	FrameworkFEDRAMPLow  = "FEDRAMP_LOW"  // FedRAMP Low baseline

	// Other Security Frameworks
	FrameworkMITREATTACK = "MITRE_ATTACK" // MITRE ATT&CK
	FrameworkCISControls = "CIS_CONTROLS" // CIS Critical Security Controls
	FrameworkSOC2        = "SOC_2"        // SOC 2 Trust Services Criteria
	FrameworkISO27001    = "ISO_27001"    // ISO/IEC 27001

	// Engineering Frameworks
	FrameworkDORA = "DORA" // DORA DevOps Metrics
	FrameworkSRE  = "SRE"  // Google SRE Practices
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
	return []string{
		FrameworkNISTCSF,
		FrameworkNISTCSF2,
		FrameworkNIST80053,
		FrameworkNISTRMF,
		FrameworkNISTAIRMF,
		FrameworkNIST800171,
	}
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

// Metric type constants.
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

// Trend direction constants.
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
const (
	SLITypeAvailability = "availability" // System/service availability
	SLITypeLatency      = "latency"      // Response time, duration
	SLITypeErrorRate    = "error_rate"   // Error percentage, failure rate
	SLITypeThroughput   = "throughput"   // Request rate, traffic volume
	SLITypeSaturation   = "saturation"   // Resource utilization/exhaustion
	SLITypeUtilization  = "utilization"  // Resource usage level
	SLITypeQuality      = "quality"      // Data quality, correctness
	SLITypeFreshness    = "freshness"    // Data age, staleness
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
