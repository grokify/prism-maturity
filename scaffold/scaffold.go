// Package scaffold provides templates for creating new PRISM documents.
package scaffold

import (
	"github.com/grokify/prism-intelligence"
)

// NewDocument creates a new PRISM document with the specified domains.
// If no domains are specified, it defaults to both security and operations.
func NewDocument(domains ...string) *prism.PRISMDocument {
	if len(domains) == 0 {
		domains = []string{prism.DomainSecurity, prism.DomainOperations}
	}

	var domainDefs []prism.DomainDef
	for _, d := range domains {
		weight := 1.0 / float64(len(domains))
		desc := ""
		switch d {
		case prism.DomainSecurity:
			desc = "Security metrics and controls"
		case prism.DomainOperations:
			desc = "Operational metrics and SLOs"
		case prism.DomainQuality:
			desc = "Quality metrics and standards"
		}
		domainDefs = append(domainDefs, prism.DomainDef{
			Name:        d,
			Description: desc,
			Weight:      weight,
		})
	}

	doc := &prism.PRISMDocument{
		Schema: "https://github.com/grokify/prism-intelligence/schema/prism.schema.json",
		Metadata: &prism.Metadata{
			Name:        "My PRISM Document",
			Description: "PRISM metrics for SaaS health monitoring",
			Version:     "1.0.0",
		},
		Domains:  domainDefs,
		Maturity: prism.NewMaturityModelForDomains(domains),
		Metrics:  make([]prism.Metric, 0),
	}

	return doc
}

// OperationsMetrics returns a set of example operations metrics.
func OperationsMetrics() []prism.Metric {
	return []prism.Metric{
		{
			ID:             "ops-availability-01",
			Name:           "Service Availability",
			Description:    "Percentage of time the service is available",
			Domain:         prism.DomainOperations,
			Stage:          prism.StageRuntime,
			Category:       prism.CategoryReliability,
			MetricType:     prism.MetricTypeRate,
			TrendDirection: prism.TrendHigherBetter,
			Unit:           "%",
			Baseline:       99.0,
			Current:        99.9,
			Target:         99.95,
			Thresholds:     &prism.Thresholds{Green: 99.9, Yellow: 99.5, Red: 99.0},
			SLI:            &prism.SLI{Name: "Availability", Formula: "successful_requests / total_requests"},
			SLO:            &prism.SLO{Target: ">=99.95%", Window: prism.Window30Days},
			FrameworkMappings: []prism.FrameworkMapping{
				{Framework: prism.FrameworkSRE, Reference: "availability-slo"},
				{Framework: prism.FrameworkDORA, Reference: "availability"},
			},
		},
		{
			ID:             "ops-latency-01",
			Name:           "P99 Latency",
			Description:    "99th percentile response latency",
			Domain:         prism.DomainOperations,
			Stage:          prism.StageRuntime,
			Category:       prism.CategoryEfficiency,
			MetricType:     prism.MetricTypeLatency,
			TrendDirection: prism.TrendLowerBetter,
			Unit:           "ms",
			Baseline:       500,
			Current:        200,
			Target:         100,
			Thresholds:     &prism.Thresholds{Green: 150, Yellow: 300, Red: 500},
			SLI:            &prism.SLI{Name: "Latency", Formula: "percentile(response_time, 99)"},
			SLO:            &prism.SLO{Target: "<=100ms", Window: prism.Window7Days},
		},
	}
}

// SecurityMetrics returns a set of example security metrics.
func SecurityMetrics() []prism.Metric {
	return []prism.Metric{
		{
			ID:             "sec-vuln-01",
			Name:           "Critical Vulnerabilities",
			Description:    "Number of unresolved critical vulnerabilities",
			Domain:         prism.DomainSecurity,
			Stage:          prism.StageRuntime,
			Category:       prism.CategoryDetection,
			MetricType:     prism.MetricTypeCount,
			TrendDirection: prism.TrendLowerBetter,
			Unit:           "count",
			Baseline:       10,
			Current:        2,
			Target:         0,
			Thresholds:     &prism.Thresholds{Green: 0, Yellow: 3, Red: 5},
		},
		{
			ID:             "sec-patch-01",
			Name:           "Patch Compliance",
			Description:    "Percentage of systems with current security patches",
			Domain:         prism.DomainSecurity,
			Stage:          prism.StageRuntime,
			Category:       prism.CategoryPrevention,
			MetricType:     prism.MetricTypeRate,
			TrendDirection: prism.TrendHigherBetter,
			Unit:           "%",
			Baseline:       85,
			Current:        95,
			Target:         99,
			Thresholds:     &prism.Thresholds{Green: 98, Yellow: 95, Red: 90},
		},
	}
}
