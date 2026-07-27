package mcp

import (
	"context"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	prism "github.com/grokify/prism-maturity"
)

// Maturity Tool Input/Output types

// MaturitySummaryInput is the input for maturity_summary.
type MaturitySummaryInput struct {
	File string `json:"file" jsonschema:"description=Path to PRISM document JSON file"`
}

// MaturitySummaryOutput is the output for maturity_summary.
type MaturitySummaryOutput struct {
	File           string           `json:"file"`
	TotalLayers    int              `json:"total_layers"`
	TotalServices  int              `json:"total_services"`
	TotalGoals     int              `json:"total_goals"`
	TotalMetrics   int              `json:"total_metrics"`
	LayerSummaries []LayerSummary   `json:"layer_summaries,omitempty"`
	ServiceSummary []ServiceSummary `json:"service_summaries,omitempty"`
	GoalsByStatus  map[string]int   `json:"goals_by_status,omitempty"`
	SLOCompliance  float64          `json:"slo_compliance_percent"`
	Error          string           `json:"error,omitempty"`
}

// LayerSummary provides summary for a layer.
type LayerSummary struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Weight      float64 `json:"weight,omitempty"`
}

// ServiceSummary provides summary for a service.
type ServiceSummary struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	LayerID string `json:"layer_id,omitempty"`
	Tier    string `json:"tier,omitempty"`
}

// ListMetricsInput is the input for list_metrics.
type ListMetricsInput struct {
	File      string `json:"file" jsonschema:"description=Path to PRISM document JSON file"`
	SLOStatus string `json:"slo_status,omitempty" jsonschema:"description=Filter by SLO status (met, not_met)"`
	Limit     int    `json:"limit,omitempty" jsonschema:"description=Maximum results (default 50)"`
}

// MetricSummary provides a brief summary of a metric.
type MetricSummary struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Current float64 `json:"current"`
	Target  float64 `json:"target"`
	SLOMet  bool    `json:"slo_met"`
}

// ListMetricsOutput is the output for list_metrics.
type ListMetricsOutput struct {
	File    string          `json:"file"`
	Metrics []MetricSummary `json:"metrics"`
	Total   int             `json:"total"`
	Error   string          `json:"error,omitempty"`
}

// GetMetricInput is the input for get_metric.
type GetMetricInput struct {
	File string `json:"file" jsonschema:"description=Path to PRISM document JSON file"`
	ID   string `json:"id" jsonschema:"description=Metric ID"`
}

// GetMetricOutput is the output for get_metric.
type GetMetricOutput struct {
	File   string        `json:"file"`
	Metric *prism.Metric `json:"metric,omitempty"`
	SLOMet bool          `json:"slo_met"`
	Gap    float64       `json:"gap"`
	Error  string        `json:"error,omitempty"`
}

// ListLayersInput is the input for list_layers.
type ListLayersInput struct {
	File string `json:"file" jsonschema:"description=Path to PRISM document JSON file"`
}

// ListLayersOutput is the output for list_layers.
type ListLayersOutput struct {
	File   string         `json:"file"`
	Layers []LayerSummary `json:"layers"`
	Total  int            `json:"total"`
	Error  string         `json:"error,omitempty"`
}

// ListServicesInput is the input for list_services.
type ListServicesInput struct {
	File    string `json:"file" jsonschema:"description=Path to PRISM document JSON file"`
	LayerID string `json:"layer_id,omitempty" jsonschema:"description=Filter by layer ID"`
	Tier    string `json:"tier,omitempty" jsonschema:"description=Filter by service tier"`
	Limit   int    `json:"limit,omitempty" jsonschema:"description=Maximum results (default 50)"`
}

// ListServicesOutput is the output for list_services.
type ListServicesOutput struct {
	File     string           `json:"file"`
	Services []ServiceSummary `json:"services"`
	Total    int              `json:"total"`
	Error    string           `json:"error,omitempty"`
}

// RegisterMaturityTools registers maturity-related MCP tools.
func (s *Server) RegisterMaturityTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "maturity_summary",
		Description: "Get overall maturity summary including layer counts, goal status, and SLO compliance.",
	}, s.maturitySummary)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_metrics",
		Description: "List metrics with optional filtering by SLO status.",
	}, s.listMetrics)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_metric",
		Description: "Get details for a specific metric including SLO status and gap.",
	}, s.getMetric)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_layers",
		Description: "List all layers defined in the document.",
	}, s.listLayers)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_services",
		Description: "List services with optional filtering by layer or tier.",
	}, s.listServices)
}

// Tool implementations

func (s *Server) maturitySummary(ctx context.Context, req *mcp.CallToolRequest, input MaturitySummaryInput) (*mcp.CallToolResult, MaturitySummaryOutput, error) {
	doc, err := loadDocument(input.File)
	if err != nil {
		return nil, MaturitySummaryOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	output := MaturitySummaryOutput{
		File:          filepath.Base(input.File),
		TotalLayers:   len(doc.Layers),
		TotalServices: len(doc.Services),
		TotalGoals:    len(doc.Goals),
		TotalMetrics:  len(doc.Metrics),
		GoalsByStatus: make(map[string]int),
	}

	// Layer summaries
	for _, layer := range doc.Layers {
		output.LayerSummaries = append(output.LayerSummaries, LayerSummary{
			ID:          layer.ID,
			Name:        layer.Name,
			Description: layer.Description,
			Weight:      layer.Weight,
		})
	}

	// Service summaries
	for _, svc := range doc.Services {
		output.ServiceSummary = append(output.ServiceSummary, ServiceSummary{
			ID:      svc.ID,
			Name:    svc.Name,
			LayerID: svc.LayerID,
			Tier:    svc.Tier,
		})
	}

	// Count goals by status
	for _, goal := range doc.Goals {
		status := goal.Status
		if status == "" {
			status = "unknown"
		}
		output.GoalsByStatus[status]++
	}

	// Calculate SLO compliance
	var metricsWithSLO, slosMet int
	for _, metric := range doc.Metrics {
		if metric.Target > 0 {
			metricsWithSLO++
			if metric.MeetsSLO() {
				slosMet++
			}
		}
	}

	if metricsWithSLO > 0 {
		output.SLOCompliance = float64(slosMet) / float64(metricsWithSLO) * 100
	}

	return nil, output, nil
}

func (s *Server) listMetrics(ctx context.Context, req *mcp.CallToolRequest, input ListMetricsInput) (*mcp.CallToolResult, ListMetricsOutput, error) {
	doc, err := loadDocument(input.File)
	if err != nil {
		return nil, ListMetricsOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	if input.Limit == 0 {
		input.Limit = 50
	}

	var summaries []MetricSummary
	for _, metric := range doc.Metrics {
		sloMet := metric.MeetsSLO()

		// Apply SLO status filter
		if input.SLOStatus != "" {
			if input.SLOStatus == "met" && !sloMet {
				continue
			}
			if input.SLOStatus == "not_met" && sloMet {
				continue
			}
		}

		summaries = append(summaries, MetricSummary{
			ID:      metric.ID,
			Name:    metric.Name,
			Current: metric.Current,
			Target:  metric.Target,
			SLOMet:  sloMet,
		})

		if len(summaries) >= input.Limit {
			break
		}
	}

	return nil, ListMetricsOutput{
		File:    filepath.Base(input.File),
		Metrics: summaries,
		Total:   len(summaries),
	}, nil
}

func (s *Server) getMetric(ctx context.Context, req *mcp.CallToolRequest, input GetMetricInput) (*mcp.CallToolResult, GetMetricOutput, error) {
	doc, err := loadDocument(input.File)
	if err != nil {
		return nil, GetMetricOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	metric := doc.GetMetricByID(input.ID)
	if metric == nil {
		return nil, GetMetricOutput{
			File:  filepath.Base(input.File),
			Error: "metric not found: " + input.ID,
		}, nil
	}

	return nil, GetMetricOutput{
		File:   filepath.Base(input.File),
		Metric: metric,
		SLOMet: metric.MeetsSLO(),
		Gap:    metric.Target - metric.Current,
	}, nil
}

func (s *Server) listLayers(ctx context.Context, req *mcp.CallToolRequest, input ListLayersInput) (*mcp.CallToolResult, ListLayersOutput, error) {
	doc, err := loadDocument(input.File)
	if err != nil {
		return nil, ListLayersOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	var summaries []LayerSummary
	for _, layer := range doc.Layers {
		summaries = append(summaries, LayerSummary{
			ID:          layer.ID,
			Name:        layer.Name,
			Description: layer.Description,
			Weight:      layer.Weight,
		})
	}

	return nil, ListLayersOutput{
		File:   filepath.Base(input.File),
		Layers: summaries,
		Total:  len(summaries),
	}, nil
}

func (s *Server) listServices(ctx context.Context, req *mcp.CallToolRequest, input ListServicesInput) (*mcp.CallToolResult, ListServicesOutput, error) {
	doc, err := loadDocument(input.File)
	if err != nil {
		return nil, ListServicesOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	if input.Limit == 0 {
		input.Limit = 50
	}

	var summaries []ServiceSummary
	for _, svc := range doc.Services {
		// Apply layer filter
		if input.LayerID != "" && svc.LayerID != input.LayerID {
			continue
		}

		// Apply tier filter
		if input.Tier != "" && svc.Tier != input.Tier {
			continue
		}

		summaries = append(summaries, ServiceSummary{
			ID:      svc.ID,
			Name:    svc.Name,
			LayerID: svc.LayerID,
			Tier:    svc.Tier,
		})

		if len(summaries) >= input.Limit {
			break
		}
	}

	return nil, ListServicesOutput{
		File:     filepath.Base(input.File),
		Services: summaries,
		Total:    len(summaries),
	}, nil
}
