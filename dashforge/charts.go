package dashforge

import (
	"fmt"

	"github.com/grokify/echartify/chartir"
	"github.com/grokify/prism-intelligence"
)

// Helper to create float64 pointer
func floatPtr(f float64) *float64 {
	return &f
}

// createMaturityRadarChart creates a radar chart showing maturity by goal.
func createMaturityRadarChart(scorecard []prism.GoalMaturityStatus) *chartir.ChartIR {
	// Build dataset with goal maturity data
	columns := []chartir.Column{
		{Name: "goal", Type: chartir.ColumnTypeString},
		{Name: "current", Type: chartir.ColumnTypeNumber},
		{Name: "target", Type: chartir.ColumnTypeNumber},
	}

	var rows [][]string
	for _, g := range scorecard {
		rows = append(rows, []string{
			g.GoalName,
			fmt.Sprintf("%d", g.CurrentLevel),
			fmt.Sprintf("%d", g.TargetLevel),
		})
	}

	return &chartir.ChartIR{
		Title: "Maturity by Goal",
		Datasets: []chartir.Dataset{
			{
				ID:      "maturity",
				Columns: columns,
				Rows:    rows,
			},
		},
		Marks: []chartir.Mark{
			{
				ID:        "current",
				DatasetID: "maturity",
				Geometry:  chartir.GeometryRadar,
				Name:      "Current",
				Encode: chartir.Encode{
					Name:  "goal",
					Value: "current",
				},
				Style: &chartir.Style{
					Color:   "#1890ff",
					Opacity: floatPtr(0.6),
				},
			},
			{
				ID:        "target",
				DatasetID: "maturity",
				Geometry:  chartir.GeometryRadar,
				Name:      "Target",
				Encode: chartir.Encode{
					Name:  "goal",
					Value: "target",
				},
				Style: &chartir.Style{
					Color:   "#52c41a",
					Opacity: floatPtr(0.3),
				},
			},
		},
		Legend: &chartir.Legend{
			Show:     true,
			Position: "bottom",
		},
		Tooltip: &chartir.Tooltip{
			Show:    true,
			Trigger: "item",
		},
	}
}

// createComplianceBarChart creates a bar chart showing SLO compliance by category.
func createComplianceBarChart(compliance prism.SLOComplianceSummary) *chartir.ChartIR {
	columns := []chartir.Column{
		{Name: "category", Type: chartir.ColumnTypeString},
		{Name: "met", Type: chartir.ColumnTypeNumber},
		{Name: "missed", Type: chartir.ColumnTypeNumber},
	}

	var rows [][]string
	for _, cat := range compliance.Categories {
		rows = append(rows, []string{
			titleCase(cat.Category),
			fmt.Sprintf("%d", cat.Met),
			fmt.Sprintf("%d", cat.Missed),
		})
	}

	return &chartir.ChartIR{
		Title: "SLO Compliance by Category",
		Datasets: []chartir.Dataset{
			{
				ID:      "compliance",
				Columns: columns,
				Rows:    rows,
			},
		},
		Marks: []chartir.Mark{
			{
				ID:        "met",
				DatasetID: "compliance",
				Geometry:  chartir.GeometryBar,
				Name:      "Met",
				Encode: chartir.Encode{
					X: "category",
					Y: "met",
				},
				Stack: "total",
				Style: &chartir.Style{
					Color: "#52c41a",
				},
			},
			{
				ID:        "missed",
				DatasetID: "compliance",
				Geometry:  chartir.GeometryBar,
				Name:      "Missed",
				Encode: chartir.Encode{
					X: "category",
					Y: "missed",
				},
				Stack: "total",
				Style: &chartir.Style{
					Color: "#f5222d",
				},
			},
		},
		Axes: []chartir.Axis{
			{
				ID:       "x",
				Type:     chartir.AxisTypeCategory,
				Position: "bottom",
			},
			{
				ID:       "y",
				Type:     chartir.AxisTypeValue,
				Position: "left",
				Name:     "SLOs",
			},
		},
		Legend: &chartir.Legend{
			Show:     true,
			Position: "top",
		},
		Tooltip: &chartir.Tooltip{
			Show:    true,
			Trigger: "axis",
		},
	}
}

// createMaturityGauge creates a gauge chart for a single goal's maturity.
func createMaturityGauge(goal prism.GoalMaturityStatus) *chartir.ChartIR {
	return &chartir.ChartIR{
		Title: goal.GoalName,
		Datasets: []chartir.Dataset{
			{
				ID: "gauge",
				Columns: []chartir.Column{
					{Name: "name", Type: chartir.ColumnTypeString},
					{Name: "value", Type: chartir.ColumnTypeNumber},
				},
				Rows: [][]string{
					{"Maturity", fmt.Sprintf("%d", goal.CurrentLevel)},
				},
			},
		},
		Marks: []chartir.Mark{
			{
				ID:        "gauge",
				DatasetID: "gauge",
				Geometry:  chartir.GeometryGauge,
				Encode: chartir.Encode{
					Name:  "name",
					Value: "value",
				},
				Style: &chartir.Style{
					GaugeMin:   floatPtr(0),
					GaugeMax:   floatPtr(5),
					StartAngle: floatPtr(225),
					EndAngle:   floatPtr(-45),
				},
			},
		},
		Tooltip: &chartir.Tooltip{
			Show: true,
		},
	}
}

// createGoalGauge creates a detailed gauge for a goal dashboard.
func createGoalGauge(goal *prism.Goal) *chartir.ChartIR {
	return &chartir.ChartIR{
		Title: "Current Maturity Level",
		Datasets: []chartir.Dataset{
			{
				ID: "gauge",
				Columns: []chartir.Column{
					{Name: "name", Type: chartir.ColumnTypeString},
					{Name: "value", Type: chartir.ColumnTypeNumber},
				},
				Rows: [][]string{
					{"Level", fmt.Sprintf("%d", goal.CurrentLevel)},
				},
			},
		},
		Marks: []chartir.Mark{
			{
				ID:        "gauge",
				DatasetID: "gauge",
				Geometry:  chartir.GeometryGauge,
				Encode: chartir.Encode{
					Name:  "name",
					Value: "value",
				},
				Style: &chartir.Style{
					GaugeMin:   floatPtr(0),
					GaugeMax:   floatPtr(5),
					StartAngle: floatPtr(225),
					EndAngle:   floatPtr(-45),
				},
			},
		},
	}
}

// createPhaseProgressChart creates a bar chart showing phase completion.
func createPhaseProgressChart(phases []prism.PhaseProgressSummary) *chartir.ChartIR {
	columns := []chartir.Column{
		{Name: "phase", Type: chartir.ColumnTypeString},
		{Name: "completion", Type: chartir.ColumnTypeNumber},
	}

	var rows [][]string
	for _, p := range phases {
		phaseName := p.PhaseName
		if p.Quarter != "" && p.Year > 0 {
			phaseName = fmt.Sprintf("%s %d", p.Quarter, p.Year)
		}
		rows = append(rows, []string{
			phaseName,
			fmt.Sprintf("%.1f", p.CompletionPct),
		})
	}

	return &chartir.ChartIR{
		Title: "Phase Completion",
		Datasets: []chartir.Dataset{
			{
				ID:      "phases",
				Columns: columns,
				Rows:    rows,
			},
		},
		Marks: []chartir.Mark{
			{
				ID:        "completion",
				DatasetID: "phases",
				Geometry:  chartir.GeometryBar,
				Encode: chartir.Encode{
					X: "phase",
					Y: "completion",
				},
				Style: &chartir.Style{
					Color: "#1890ff",
				},
			},
		},
		Axes: []chartir.Axis{
			{
				ID:       "x",
				Type:     chartir.AxisTypeCategory,
				Position: "bottom",
			},
			{
				ID:       "y",
				Type:     chartir.AxisTypeValue,
				Position: "left",
				Name:     "Completion %",
				Max:      floatPtr(100),
			},
		},
		Tooltip: &chartir.Tooltip{
			Show:    true,
			Trigger: "axis",
		},
	}
}

// createGapBarChart creates a horizontal bar chart for gap analysis.
func createGapBarChart(gaps []prism.GapAnalysisEntry) *chartir.ChartIR {
	columns := []chartir.Column{
		{Name: "metric", Type: chartir.ColumnTypeString},
		{Name: "gap", Type: chartir.ColumnTypeNumber},
		{Name: "priority", Type: chartir.ColumnTypeString},
	}

	// Limit to top 10 gaps
	maxGaps := 10
	if len(gaps) < maxGaps {
		maxGaps = len(gaps)
	}

	var rows [][]string
	for i := 0; i < maxGaps; i++ {
		g := gaps[i]
		rows = append(rows, []string{
			g.MetricName,
			fmt.Sprintf("%.1f", g.GapPercent),
			g.Priority,
		})
	}

	return &chartir.ChartIR{
		Title: "Top Gaps by Impact",
		Datasets: []chartir.Dataset{
			{
				ID:      "gaps",
				Columns: columns,
				Rows:    rows,
			},
		},
		Marks: []chartir.Mark{
			{
				ID:        "gap",
				DatasetID: "gaps",
				Geometry:  chartir.GeometryBar,
				Encode: chartir.Encode{
					Y: "metric",
					X: "gap",
				},
				Style: &chartir.Style{
					Color: "#f5222d",
				},
			},
		},
		Axes: []chartir.Axis{
			{
				ID:       "y",
				Type:     chartir.AxisTypeCategory,
				Position: "left",
			},
			{
				ID:       "x",
				Type:     chartir.AxisTypeValue,
				Position: "bottom",
				Name:     "Gap %",
			},
		},
		Tooltip: &chartir.Tooltip{
			Show:    true,
			Trigger: "axis",
		},
		Grid: &chartir.Grid{
			Left:         "25%",
			ContainLabel: true,
		},
	}
}

// createSLOHeatmap creates a heatmap showing SLO coverage by level.
func createSLOHeatmap(report *prism.SLOReport) *chartir.ChartIR {
	// Count SLOs per category per level
	type key struct {
		category string
		level    int
	}
	counts := make(map[key]int)
	var categories []string
	catSet := make(map[string]bool)

	for _, entry := range report.Entries {
		k := key{entry.Category, entry.Level}
		counts[k]++
		if !catSet[entry.Category] {
			categories = append(categories, entry.Category)
			catSet[entry.Category] = true
		}
	}

	// Build heatmap data
	columns := []chartir.Column{
		{Name: "category", Type: chartir.ColumnTypeString},
		{Name: "level", Type: chartir.ColumnTypeNumber},
		{Name: "count", Type: chartir.ColumnTypeNumber},
	}

	var rows [][]string
	for _, cat := range categories {
		for level := 1; level <= 5; level++ {
			k := key{cat, level}
			count := counts[k]
			rows = append(rows, []string{
				titleCase(cat),
				fmt.Sprintf("%d", level),
				fmt.Sprintf("%d", count),
			})
		}
	}

	return &chartir.ChartIR{
		Title: "SLO Requirements by Level",
		Datasets: []chartir.Dataset{
			{
				ID:      "heatmap",
				Columns: columns,
				Rows:    rows,
			},
		},
		Marks: []chartir.Mark{
			{
				ID:        "heatmap",
				DatasetID: "heatmap",
				Geometry:  chartir.GeometryHeatmap,
				Encode: chartir.Encode{
					X:    "level",
					Y:    "category",
					Heat: "count",
				},
			},
		},
		Axes: []chartir.Axis{
			{
				ID:       "x",
				Type:     chartir.AxisTypeCategory,
				Position: "bottom",
				Name:     "Maturity Level",
			},
			{
				ID:       "y",
				Type:     chartir.AxisTypeCategory,
				Position: "left",
			},
		},
		Tooltip: &chartir.Tooltip{
			Show: true,
		},
	}
}

// titleCase converts a string to title case.
func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	// Simple title case: uppercase first letter if lowercase
	if s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}
	return s
}
