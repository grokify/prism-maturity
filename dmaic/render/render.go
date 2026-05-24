// Package render provides interfaces and utilities for rendering DMAIC documents
// to various output formats including Marp slides.
package render

import "github.com/grokify/prism-maturity/dmaic"

// Renderer defines the interface for output format renderers.
type Renderer interface {
	// Format returns the output format name (e.g., "marp").
	Format() string
	// FileExtension returns the file extension for this format (e.g., ".md").
	FileExtension() string
	// Render converts a DMAIC document to the output format.
	Render(doc *dmaic.DMAICDocument, opts *Options) ([]byte, error)
}

// Options contains rendering options common to all renderers.
type Options struct {
	// Theme name (renderer-specific, e.g., "default", "corporate", "minimal")
	Theme string

	// IncludeDataPoints includes historical data point slides
	IncludeDataPoints bool

	// IncludeRootCauses includes root cause analysis slides
	IncludeRootCauses bool

	// IncludeInitiatives includes improvement initiatives slides
	IncludeInitiatives bool

	// ShowControlCharts shows SPC control chart visualizations
	ShowControlCharts bool

	// ShowCapabilityMetrics shows Six Sigma capability metrics
	ShowCapabilityMetrics bool

	// GroupByPhase groups metrics by DMAIC phase
	GroupByPhase bool

	// GroupByStatus groups metrics by status (Green/Yellow/Red)
	GroupByStatus bool

	// MaxDataPoints limits data points shown (0 = all)
	MaxDataPoints int

	// Custom CSS (for Marp/HTML renderers)
	CustomCSS string

	// Additional metadata (renderer-specific)
	Metadata map[string]string
}

// DefaultOptions returns sensible default rendering options.
func DefaultOptions() *Options {
	return &Options{
		Theme:                 "default",
		IncludeDataPoints:     false,
		IncludeRootCauses:     true,
		IncludeInitiatives:    true,
		ShowControlCharts:     false,
		ShowCapabilityMetrics: false,
		GroupByPhase:          false,
		GroupByStatus:         false,
		MaxDataPoints:         10,
		Metadata:              make(map[string]string),
	}
}

// ExecutiveOptions returns options for executive-focused slides (fewer details).
func ExecutiveOptions() *Options {
	return &Options{
		Theme:                 "corporate",
		IncludeDataPoints:     false,
		IncludeRootCauses:     false,
		IncludeInitiatives:    true,
		ShowControlCharts:     false,
		ShowCapabilityMetrics: false,
		GroupByPhase:          false,
		GroupByStatus:         true,
		MaxDataPoints:         0,
		Metadata:              make(map[string]string),
	}
}

// SixSigmaOptions returns options for Six Sigma detailed analysis.
func SixSigmaOptions() *Options {
	return &Options{
		Theme:                 "default",
		IncludeDataPoints:     true,
		IncludeRootCauses:     true,
		IncludeInitiatives:    true,
		ShowControlCharts:     true,
		ShowCapabilityMetrics: true,
		GroupByPhase:          true,
		GroupByStatus:         false,
		MaxDataPoints:         20,
		Metadata:              make(map[string]string),
	}
}
