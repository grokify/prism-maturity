package prism

// SLITypeCoverage represents coverage statistics for SLI types.
type SLITypeCoverage struct {
	TotalMetrics   int            `json:"totalMetrics"`
	MetricsWithSLI int            `json:"metricsWithSli"`
	ByType         map[string]int `json:"byType"`        // Count per SLI type
	MissingTypes   []string       `json:"missingTypes"`  // SLI types with no metrics
	CoverageRatio  float64        `json:"coverageRatio"` // Ratio of metrics with SLI type
}

// MethodologyCoverage represents coverage for a specific observability methodology.
type MethodologyCoverage struct {
	Methodology   string              `json:"methodology"`
	RequiredTypes []string            `json:"requiredTypes"`           // SLI types required by this methodology
	CoveredTypes  []string            `json:"coveredTypes"`            // SLI types that have metrics
	MissingTypes  []string            `json:"missingTypes"`            // SLI types missing metrics
	ByType        map[string]int      `json:"byType"`                  // Count per SLI type
	MetricsByType map[string][]string `json:"metricsByType,omitempty"` // Metric IDs per SLI type
	CoverageRatio float64             `json:"coverageRatio"`           // Ratio of required types covered
	IsComplete    bool                `json:"isComplete"`              // True if all required types are covered
}

// AnalyzeSLICoverage analyzes SLI type coverage across all metrics.
func (doc *PRISMDocument) AnalyzeSLICoverage() *SLITypeCoverage {
	coverage := &SLITypeCoverage{
		TotalMetrics: len(doc.Metrics),
		ByType:       make(map[string]int),
	}

	// Initialize all SLI types with zero count
	for _, t := range AllSLITypes() {
		coverage.ByType[t] = 0
	}

	// Count metrics by SLI type
	for _, m := range doc.Metrics {
		if m.SLI != nil && m.SLI.SLIType != "" {
			coverage.MetricsWithSLI++
			coverage.ByType[m.SLI.SLIType]++
		}
	}

	// Calculate coverage ratio
	if coverage.TotalMetrics > 0 {
		coverage.CoverageRatio = float64(coverage.MetricsWithSLI) / float64(coverage.TotalMetrics)
	}

	// Identify missing types (types with zero metrics)
	for _, t := range AllSLITypes() {
		if coverage.ByType[t] == 0 {
			coverage.MissingTypes = append(coverage.MissingTypes, t)
		}
	}

	return coverage
}

// AnalyzeSLICoverageByLayer analyzes SLI type coverage grouped by layer.
func (doc *PRISMDocument) AnalyzeSLICoverageByLayer() map[string]*SLITypeCoverage {
	result := make(map[string]*SLITypeCoverage)

	// Initialize coverage for each layer
	for _, layer := range AllLayers() {
		result[layer] = &SLITypeCoverage{
			ByType: make(map[string]int),
		}
		// Initialize all SLI types with zero count
		for _, t := range AllSLITypes() {
			result[layer].ByType[t] = 0
		}
	}

	// Count metrics by layer and SLI type
	for _, m := range doc.Metrics {
		layer := m.Layer
		if layer == "" {
			continue
		}
		if _, exists := result[layer]; !exists {
			continue
		}

		result[layer].TotalMetrics++
		if m.SLI != nil && m.SLI.SLIType != "" {
			result[layer].MetricsWithSLI++
			result[layer].ByType[m.SLI.SLIType]++
		}
	}

	// Calculate ratios and missing types for each layer
	for _, coverage := range result {
		if coverage.TotalMetrics > 0 {
			coverage.CoverageRatio = float64(coverage.MetricsWithSLI) / float64(coverage.TotalMetrics)
		}
		for _, t := range AllSLITypes() {
			if coverage.ByType[t] == 0 {
				coverage.MissingTypes = append(coverage.MissingTypes, t)
			}
		}
	}

	return result
}

// AnalyzeMethodologyCoverage analyzes coverage for a specific observability methodology.
func (doc *PRISMDocument) AnalyzeMethodologyCoverage(methodology string) *MethodologyCoverage {
	requiredTypes := SLITypesForMethodology(methodology)
	if requiredTypes == nil {
		return nil
	}

	coverage := &MethodologyCoverage{
		Methodology:   methodology,
		RequiredTypes: requiredTypes,
		ByType:        make(map[string]int),
		MetricsByType: make(map[string][]string),
	}

	// Initialize required types with zero count
	for _, t := range requiredTypes {
		coverage.ByType[t] = 0
		coverage.MetricsByType[t] = []string{}
	}

	// Count metrics by SLI type (only for required types)
	for _, m := range doc.Metrics {
		if m.SLI == nil || m.SLI.SLIType == "" {
			continue
		}
		sliType := m.SLI.SLIType
		if _, isRequired := coverage.ByType[sliType]; isRequired {
			coverage.ByType[sliType]++
			metricID := m.ID
			if metricID == "" {
				metricID = m.Name
			}
			coverage.MetricsByType[sliType] = append(coverage.MetricsByType[sliType], metricID)
		}
	}

	// Determine covered and missing types
	for _, t := range requiredTypes {
		if coverage.ByType[t] > 0 {
			coverage.CoveredTypes = append(coverage.CoveredTypes, t)
		} else {
			coverage.MissingTypes = append(coverage.MissingTypes, t)
		}
	}

	// Calculate coverage ratio
	if len(requiredTypes) > 0 {
		coverage.CoverageRatio = float64(len(coverage.CoveredTypes)) / float64(len(requiredTypes))
	}
	coverage.IsComplete = len(coverage.MissingTypes) == 0

	return coverage
}

// AnalyzeAllMethodologyCoverage analyzes coverage for all observability methodologies.
func (doc *PRISMDocument) AnalyzeAllMethodologyCoverage() map[string]*MethodologyCoverage {
	result := make(map[string]*MethodologyCoverage)
	for _, methodology := range AllMethodologies() {
		result[methodology] = doc.AnalyzeMethodologyCoverage(methodology)
	}
	return result
}

// GetMetricsBySLIType returns all metrics with the specified SLI type.
func (doc *PRISMDocument) GetMetricsBySLIType(sliType string) []Metric {
	var result []Metric
	for _, m := range doc.Metrics {
		if m.SLI != nil && m.SLI.SLIType == sliType {
			result = append(result, m)
		}
	}
	return result
}

// GetMetricsByMethodology returns all metrics that belong to a specific methodology.
func (doc *PRISMDocument) GetMetricsByMethodology(methodology string) []Metric {
	requiredTypes := SLITypesForMethodology(methodology)
	if requiredTypes == nil {
		return nil
	}

	typeSet := make(map[string]bool)
	for _, t := range requiredTypes {
		typeSet[t] = true
	}

	var result []Metric
	for _, m := range doc.Metrics {
		if m.SLI != nil && typeSet[m.SLI.SLIType] {
			result = append(result, m)
		}
	}
	return result
}
