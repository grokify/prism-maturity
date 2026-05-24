// Package dashforge provides conversion from PRISM documents to Dashforge dashboard IR.
package dashforge

import (
	"encoding/json"
	"fmt"

	"github.com/grokify/prism-maturity"
	"github.com/plexusone/dashforge/dashboardir"
)

// DashboardSet contains all generated dashboards for a PRISM document.
type DashboardSet struct {
	// Executive is the main executive summary dashboard.
	Executive *dashboardir.Dashboard `json:"executive"`

	// Maturity is the maturity scorecard dashboard.
	Maturity *dashboardir.Dashboard `json:"maturity"`

	// SLOMatrix is the SLO compliance matrix dashboard.
	SLOMatrix *dashboardir.Dashboard `json:"sloMatrix"`

	// Roadmap is the roadmap progress dashboard.
	Roadmap *dashboardir.Dashboard `json:"roadmap"`

	// Gaps is the gap analysis dashboard.
	Gaps *dashboardir.Dashboard `json:"gaps"`

	// Goals are individual goal deep-dive dashboards.
	Goals map[string]*dashboardir.Dashboard `json:"goals,omitempty"`
}

// ConvertOptions configures the dashboard conversion.
type ConvertOptions struct {
	// BaseID is the prefix for dashboard IDs.
	BaseID string

	// DataSourcePath is the path to the PRISM JSON file (for URL data source).
	DataSourcePath string

	// GenerateGoalDashboards creates individual goal dashboards.
	GenerateGoalDashboards bool

	// Theme customizes the visual appearance.
	Theme *dashboardir.Theme
}

// DefaultConvertOptions returns sensible defaults.
func DefaultConvertOptions() *ConvertOptions {
	return &ConvertOptions{
		BaseID:                 "prism",
		DataSourcePath:         "./data/prism.json",
		GenerateGoalDashboards: true,
		Theme: &dashboardir.Theme{
			Mode:            "light",
			PrimaryColor:    "#1890ff",
			BackgroundColor: "#f5f5f5",
		},
	}
}

// Convert generates Dashforge dashboards from a PRISM document.
func Convert(doc *prism.PRISMDocument, opts *ConvertOptions) (*DashboardSet, error) {
	if opts == nil {
		opts = DefaultConvertOptions()
	}

	// Generate the executive dashboard data
	execDash := doc.GenerateExecutiveDashboard()

	set := &DashboardSet{
		Goals: make(map[string]*dashboardir.Dashboard),
	}

	// Generate each dashboard
	set.Executive = generateExecutiveDashboard(doc, execDash, opts)
	set.Maturity = generateMaturityDashboard(doc, execDash, opts)
	set.SLOMatrix = generateSLOMatrixDashboard(doc, opts)
	set.Roadmap = generateRoadmapDashboard(doc, execDash, opts)
	set.Gaps = generateGapsDashboard(doc, execDash, opts)

	// Generate goal dashboards
	if opts.GenerateGoalDashboards {
		for _, goal := range doc.Goals {
			set.Goals[goal.ID] = generateGoalDashboard(doc, &goal, opts)
		}
	}

	return set, nil
}

// generateExecutiveDashboard creates the main executive summary dashboard.
func generateExecutiveDashboard(doc *prism.PRISMDocument, exec *prism.ExecutiveDashboard, opts *ConvertOptions) *dashboardir.Dashboard {
	title := "Executive Security Dashboard"
	if doc.Metadata != nil && doc.Metadata.Name != "" {
		title = doc.Metadata.Name
	}

	dash := &dashboardir.Dashboard{
		ID:          fmt.Sprintf("%s-executive", opts.BaseID),
		Title:       title,
		Description: "High-level maturity progress and SLO compliance overview",
		Layout: dashboardir.Layout{
			Type:      dashboardir.LayoutTypeGrid,
			Columns:   12,
			RowHeight: 80,
			Gap:       16,
			Padding:   20,
		},
		Theme:       opts.Theme,
		DataSources: []dashboardir.DataSource{},
		Widgets:     []dashboardir.Widget{},
	}

	// Add inline data source with executive summary data
	summaryData, _ := json.Marshal(exec.Summary)
	dash.DataSources = append(dash.DataSources, dashboardir.DataSource{
		ID:   "summary",
		Name: "Executive Summary",
		Type: dashboardir.DataSourceTypeInline,
		Data: summaryData,
	})

	scorecardData, _ := json.Marshal(exec.MaturityScorecard)
	dash.DataSources = append(dash.DataSources, dashboardir.DataSource{
		ID:   "scorecard",
		Name: "Maturity Scorecard",
		Type: dashboardir.DataSourceTypeInline,
		Data: scorecardData,
	})

	complianceData, _ := json.Marshal(exec.SLOCompliance)
	dash.DataSources = append(dash.DataSources, dashboardir.DataSource{
		ID:   "compliance",
		Name: "SLO Compliance",
		Type: dashboardir.DataSourceTypeInline,
		Data: complianceData,
	})

	// Row 0: Key metric cards
	dash.Widgets = append(dash.Widgets,
		createMetricWidget("w-maturity", "Overall Maturity", "summary", "overallMaturity",
			dashboardir.Position{X: 0, Y: 0, W: 3, H: 2},
			[]dashboardir.MetricThreshold{
				{Value: 3, Color: "#52c41a"},
				{Value: 2, Color: "#faad14"},
				{Value: 0, Color: "#f5222d"},
			},
		),
		createMetricWidget("w-compliance", "SLO Compliance", "compliance", "overallCompliance",
			dashboardir.Position{X: 3, Y: 0, W: 3, H: 2},
			[]dashboardir.MetricThreshold{
				{Value: 90, Color: "#52c41a"},
				{Value: 70, Color: "#faad14"},
				{Value: 0, Color: "#f5222d"},
			},
		),
		createMetricWidget("w-goals-on-track", "Goals On Track", "summary", "goalsOnTrack",
			dashboardir.Position{X: 6, Y: 0, W: 3, H: 2}, nil),
		createMetricWidget("w-initiatives", "Initiatives", "summary", "initiativesTotal",
			dashboardir.Position{X: 9, Y: 0, W: 3, H: 2}, nil),
	)

	// Row 2: Radar chart for maturity by goal
	radarChart := createMaturityRadarChart(exec.MaturityScorecard)
	radarConfig, _ := json.Marshal(radarChart)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:           "w-maturity-radar",
		Title:        "Maturity by Goal",
		Type:         dashboardir.WidgetTypeChart,
		Position:     dashboardir.Position{X: 0, Y: 2, W: 6, H: 4},
		DataSourceID: "scorecard",
		Config:       radarConfig,
	})

	// Row 2: SLO compliance bar chart
	complianceChart := createComplianceBarChart(exec.SLOCompliance)
	complianceConfig, _ := json.Marshal(complianceChart)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:           "w-compliance-chart",
		Title:        "SLO Compliance by Category",
		Type:         dashboardir.WidgetTypeChart,
		Position:     dashboardir.Position{X: 6, Y: 2, W: 6, H: 4},
		DataSourceID: "compliance",
		Config:       complianceConfig,
	})

	// Row 6: Maturity scorecard table
	scorecardTable := createScorecardTable()
	scorecardConfig, _ := json.Marshal(scorecardTable)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:           "w-scorecard-table",
		Title:        "Goal Status",
		Type:         dashboardir.WidgetTypeTable,
		Position:     dashboardir.Position{X: 0, Y: 6, W: 12, H: 4},
		DataSourceID: "scorecard",
		Config:       scorecardConfig,
		DrillDown: &dashboardir.DrillDown{
			Type:   dashboardir.DrillDownTypeDashboard,
			Target: fmt.Sprintf("%s-maturity", opts.BaseID),
			Params: map[string]string{"goalId": "goalId"},
		},
	})

	return dash
}

// generateMaturityDashboard creates the maturity scorecard dashboard.
func generateMaturityDashboard(_ *prism.PRISMDocument, exec *prism.ExecutiveDashboard, opts *ConvertOptions) *dashboardir.Dashboard {
	dash := &dashboardir.Dashboard{
		ID:          fmt.Sprintf("%s-maturity", opts.BaseID),
		Title:       "Maturity Scorecard",
		Description: "Detailed view of maturity levels and progress by goal",
		Layout: dashboardir.Layout{
			Type:      dashboardir.LayoutTypeGrid,
			Columns:   12,
			RowHeight: 80,
			Gap:       16,
			Padding:   20,
		},
		Theme:       opts.Theme,
		DataSources: []dashboardir.DataSource{},
		Widgets:     []dashboardir.Widget{},
	}

	// Add scorecard data
	scorecardData, _ := json.Marshal(exec.MaturityScorecard)
	dash.DataSources = append(dash.DataSources, dashboardir.DataSource{
		ID:   "scorecard",
		Name: "Maturity Scorecard",
		Type: dashboardir.DataSourceTypeInline,
		Data: scorecardData,
	})

	// Create gauge widgets for each goal
	row := 0
	for i, goal := range exec.MaturityScorecard {
		col := (i % 3) * 4
		if i > 0 && i%3 == 0 {
			row += 3
		}

		gaugeChart := createMaturityGauge(goal)
		gaugeConfig, _ := json.Marshal(gaugeChart)

		dash.Widgets = append(dash.Widgets, dashboardir.Widget{
			ID:       fmt.Sprintf("w-gauge-%s", goal.GoalID),
			Title:    goal.GoalName,
			Type:     dashboardir.WidgetTypeChart,
			Position: dashboardir.Position{X: col, Y: row, W: 4, H: 3},
			Config:   gaugeConfig,
			DrillDown: &dashboardir.DrillDown{
				Type:   dashboardir.DrillDownTypeDashboard,
				Target: fmt.Sprintf("%s-goal-%s", opts.BaseID, goal.GoalID),
			},
		})
	}

	// Add detailed table at the bottom
	row += 3
	detailTable := createMaturityDetailTable()
	detailConfig, _ := json.Marshal(detailTable)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:           "w-maturity-detail",
		Title:        "Maturity Details",
		Type:         dashboardir.WidgetTypeTable,
		Position:     dashboardir.Position{X: 0, Y: row, W: 12, H: 4},
		DataSourceID: "scorecard",
		Config:       detailConfig,
	})

	return dash
}

// generateSLOMatrixDashboard creates the SLO compliance matrix dashboard.
func generateSLOMatrixDashboard(doc *prism.PRISMDocument, opts *ConvertOptions) *dashboardir.Dashboard {
	sloReport := doc.GenerateSLOReport()

	dash := &dashboardir.Dashboard{
		ID:          fmt.Sprintf("%s-slo-matrix", opts.BaseID),
		Title:       "SLO Requirements Matrix",
		Description: "SLO requirements by category and maturity level",
		Layout: dashboardir.Layout{
			Type:      dashboardir.LayoutTypeGrid,
			Columns:   12,
			RowHeight: 80,
			Gap:       16,
			Padding:   20,
		},
		Theme:       opts.Theme,
		DataSources: []dashboardir.DataSource{},
		Widgets:     []dashboardir.Widget{},
	}

	// Add SLO entries data
	entriesData, _ := json.Marshal(sloReport.Entries)
	dash.DataSources = append(dash.DataSources, dashboardir.DataSource{
		ID:   "slo-entries",
		Name: "SLO Entries",
		Type: dashboardir.DataSourceTypeInline,
		Data: entriesData,
	})

	categoriesData, _ := json.Marshal(sloReport.Categories)
	dash.DataSources = append(dash.DataSources, dashboardir.DataSource{
		ID:   "slo-categories",
		Name: "SLO Categories",
		Type: dashboardir.DataSourceTypeInline,
		Data: categoriesData,
	})

	// Heatmap chart showing SLO coverage
	heatmapChart := createSLOHeatmap(sloReport)
	heatmapConfig, _ := json.Marshal(heatmapChart)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:           "w-slo-heatmap",
		Title:        "SLO Coverage by Level",
		Type:         dashboardir.WidgetTypeChart,
		Position:     dashboardir.Position{X: 0, Y: 0, W: 12, H: 5},
		DataSourceID: "slo-entries",
		Config:       heatmapConfig,
	})

	// SLO detail table
	sloTable := createSLODetailTable()
	sloConfig, _ := json.Marshal(sloTable)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:           "w-slo-table",
		Title:        "SLO Requirements Detail",
		Type:         dashboardir.WidgetTypeTable,
		Position:     dashboardir.Position{X: 0, Y: 5, W: 12, H: 5},
		DataSourceID: "slo-entries",
		Config:       sloConfig,
	})

	return dash
}

// generateRoadmapDashboard creates the roadmap progress dashboard.
func generateRoadmapDashboard(_ *prism.PRISMDocument, exec *prism.ExecutiveDashboard, opts *ConvertOptions) *dashboardir.Dashboard {
	dash := &dashboardir.Dashboard{
		ID:          fmt.Sprintf("%s-roadmap", opts.BaseID),
		Title:       "Roadmap Progress",
		Description: "Phase-by-phase progress toward maturity goals",
		Layout: dashboardir.Layout{
			Type:      dashboardir.LayoutTypeGrid,
			Columns:   12,
			RowHeight: 80,
			Gap:       16,
			Padding:   20,
		},
		Theme:       opts.Theme,
		DataSources: []dashboardir.DataSource{},
		Widgets:     []dashboardir.Widget{},
	}

	// Add phase progress data
	phaseData, _ := json.Marshal(exec.PhaseProgress)
	dash.DataSources = append(dash.DataSources, dashboardir.DataSource{
		ID:   "phases",
		Name: "Phase Progress",
		Type: dashboardir.DataSourceTypeInline,
		Data: phaseData,
	})

	// Phase progress bar chart
	phaseChart := createPhaseProgressChart(exec.PhaseProgress)
	phaseConfig, _ := json.Marshal(phaseChart)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:           "w-phase-progress",
		Title:        "Phase Completion",
		Type:         dashboardir.WidgetTypeChart,
		Position:     dashboardir.Position{X: 0, Y: 0, W: 8, H: 4},
		DataSourceID: "phases",
		Config:       phaseConfig,
	})

	// Current phase metrics
	for i, phase := range exec.PhaseProgress {
		if phase.IsCurrent {
			dash.Widgets = append(dash.Widgets,
				createMetricWidget(
					fmt.Sprintf("w-current-phase-%d", i),
					"Current Phase",
					"phases",
					"phaseName",
					dashboardir.Position{X: 8, Y: 0, W: 4, H: 2},
					nil,
				),
			)
			break
		}
	}

	// Phase detail table
	phaseTable := createPhaseDetailTable()
	phaseTableConfig, _ := json.Marshal(phaseTable)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:           "w-phase-table",
		Title:        "Phase Details",
		Type:         dashboardir.WidgetTypeTable,
		Position:     dashboardir.Position{X: 0, Y: 4, W: 12, H: 4},
		DataSourceID: "phases",
		Config:       phaseTableConfig,
	})

	return dash
}

// generateGapsDashboard creates the gap analysis dashboard.
func generateGapsDashboard(_ *prism.PRISMDocument, exec *prism.ExecutiveDashboard, opts *ConvertOptions) *dashboardir.Dashboard {
	dash := &dashboardir.Dashboard{
		ID:          fmt.Sprintf("%s-gaps", opts.BaseID),
		Title:       "Gap Analysis",
		Description: "Priority gaps between current state and targets",
		Layout: dashboardir.Layout{
			Type:      dashboardir.LayoutTypeGrid,
			Columns:   12,
			RowHeight: 80,
			Gap:       16,
			Padding:   20,
		},
		Theme:       opts.Theme,
		DataSources: []dashboardir.DataSource{},
		Widgets:     []dashboardir.Widget{},
	}

	// Add gaps data
	gapsData, _ := json.Marshal(exec.Gaps)
	dash.DataSources = append(dash.DataSources, dashboardir.DataSource{
		ID:   "gaps",
		Name: "Gap Analysis",
		Type: dashboardir.DataSourceTypeInline,
		Data: gapsData,
	})

	// Gap bar chart (horizontal, sorted by impact)
	gapChart := createGapBarChart(exec.Gaps)
	gapConfig, _ := json.Marshal(gapChart)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:           "w-gap-chart",
		Title:        "Gaps by Impact",
		Type:         dashboardir.WidgetTypeChart,
		Position:     dashboardir.Position{X: 0, Y: 0, W: 12, H: 5},
		DataSourceID: "gaps",
		Config:       gapConfig,
	})

	// Gap detail table
	gapTable := createGapDetailTable()
	gapTableConfig, _ := json.Marshal(gapTable)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:           "w-gap-table",
		Title:        "Gap Details",
		Type:         dashboardir.WidgetTypeTable,
		Position:     dashboardir.Position{X: 0, Y: 5, W: 12, H: 5},
		DataSourceID: "gaps",
		Config:       gapTableConfig,
	})

	return dash
}

// generateGoalDashboard creates a dashboard for a specific goal.
func generateGoalDashboard(doc *prism.PRISMDocument, goal *prism.Goal, opts *ConvertOptions) *dashboardir.Dashboard {
	dash := &dashboardir.Dashboard{
		ID:          fmt.Sprintf("%s-goal-%s", opts.BaseID, goal.ID),
		Title:       goal.Name,
		Description: goal.Description,
		Layout: dashboardir.Layout{
			Type:      dashboardir.LayoutTypeGrid,
			Columns:   12,
			RowHeight: 80,
			Gap:       16,
			Padding:   20,
		},
		Theme:       opts.Theme,
		DataSources: []dashboardir.DataSource{},
		Widgets:     []dashboardir.Widget{},
	}

	// Goal metadata
	goalData, _ := json.Marshal(map[string]interface{}{
		"id":           goal.ID,
		"name":         goal.Name,
		"description":  goal.Description,
		"currentLevel": goal.CurrentLevel,
		"targetLevel":  goal.TargetLevel,
	})
	dash.DataSources = append(dash.DataSources, dashboardir.DataSource{
		ID:   "goal",
		Name: "Goal Data",
		Type: dashboardir.DataSourceTypeInline,
		Data: goalData,
	})

	// Maturity gauge
	gaugeChart := createGoalGauge(goal)
	gaugeConfig, _ := json.Marshal(gaugeChart)
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:       "w-goal-gauge",
		Title:    "Current Maturity",
		Type:     dashboardir.WidgetTypeChart,
		Position: dashboardir.Position{X: 0, Y: 0, W: 4, H: 3},
		Config:   gaugeConfig,
	})

	// Goal description text
	textConfig, _ := json.Marshal(dashboardir.TextConfig{
		Content: fmt.Sprintf("## %s\n\n%s\n\n**Current Level:** L%d\n**Target Level:** L%d",
			goal.Name, goal.Description, goal.CurrentLevel, goal.TargetLevel),
		Format: dashboardir.TextFormatMarkdown,
	})
	dash.Widgets = append(dash.Widgets, dashboardir.Widget{
		ID:       "w-goal-description",
		Type:     dashboardir.WidgetTypeText,
		Position: dashboardir.Position{X: 4, Y: 0, W: 8, H: 3},
		Config:   textConfig,
	})

	// SLO requirements for this goal
	if goal.MaturityModel != nil {
		var sloData []map[string]interface{}
		for _, level := range goal.MaturityModel.Levels {
			for _, criterion := range level.MetricCriteria {
				metric := doc.GetMetricByID(criterion.MetricID)
				metricName := criterion.MetricID
				if metric != nil {
					metricName = metric.Name
				}
				sloData = append(sloData, map[string]interface{}{
					"level":       level.Level,
					"levelName":   prism.MaturityLevelName(level.Level),
					"metricId":    criterion.MetricID,
					"metricName":  metricName,
					"operator":    criterion.Operator,
					"value":       criterion.Value,
					"requirement": formatRequirement(criterion.Operator, criterion.Value),
				})
			}
		}
		sloBytes, _ := json.Marshal(sloData)
		dash.DataSources = append(dash.DataSources, dashboardir.DataSource{
			ID:   "goal-slos",
			Name: "Goal SLOs",
			Type: dashboardir.DataSourceTypeInline,
			Data: sloBytes,
		})

		sloTable := createGoalSLOTable()
		sloConfig, _ := json.Marshal(sloTable)
		dash.Widgets = append(dash.Widgets, dashboardir.Widget{
			ID:           "w-goal-slos",
			Title:        "SLO Requirements",
			Type:         dashboardir.WidgetTypeTable,
			Position:     dashboardir.Position{X: 0, Y: 3, W: 12, H: 4},
			DataSourceID: "goal-slos",
			Config:       sloConfig,
		})
	}

	return dash
}

func formatRequirement(operator string, value float64) string {
	opSymbol := operator
	switch operator {
	case prism.SLOOperatorGTE:
		opSymbol = ">="
	case prism.SLOOperatorLTE:
		opSymbol = "<="
	case prism.SLOOperatorGT:
		opSymbol = ">"
	case prism.SLOOperatorLT:
		opSymbol = "<"
	case prism.SLOOperatorEQ:
		opSymbol = "="
	}
	return fmt.Sprintf("%s%.2f", opSymbol, value)
}
