package uiforge

import (
	"encoding/json"

	"github.com/plexusone/uiforge/dashboardir"
)

// createMetricWidget creates a single-value metric widget.
func createMetricWidget(id, title, dataSourceID, valueField string, pos dashboardir.Position, thresholds []dashboardir.MetricThreshold) dashboardir.Widget {
	config := dashboardir.MetricConfig{
		ValueField: valueField,
		Format:     "number",
		FormatOptions: &dashboardir.FormatOptions{
			Decimals: 1,
		},
		Thresholds: thresholds,
	}
	configBytes, _ := json.Marshal(config)

	return dashboardir.Widget{
		ID:           id,
		Title:        title,
		Type:         dashboardir.WidgetTypeMetric,
		Position:     pos,
		DataSourceID: dataSourceID,
		Config:       configBytes,
	}
}

// createScorecardTable creates the maturity scorecard table config.
func createScorecardTable() dashboardir.TableConfig {
	return dashboardir.TableConfig{
		Columns: []dashboardir.TableColumn{
			{Field: "goalName", Header: "Goal", Width: "30%"},
			{Field: "currentLevel", Header: "Current", Align: "center"},
			{Field: "targetLevel", Header: "Target", Align: "center"},
			{Field: "gap", Header: "Gap", Align: "center"},
			{Field: "slosMetPercent", Header: "SLOs Met", Align: "center", Format: "percent"},
			{Field: "status", Header: "Status", Align: "center"},
		},
		Sortable: true,
		Striped:  true,
	}
}

// createMaturityDetailTable creates detailed maturity table config.
func createMaturityDetailTable() dashboardir.TableConfig {
	return dashboardir.TableConfig{
		Columns: []dashboardir.TableColumn{
			{Field: "goalName", Header: "Goal", Width: "25%"},
			{Field: "currentLevel", Header: "Current", Align: "center"},
			{Field: "targetLevel", Header: "Target", Align: "center"},
			{Field: "gap", Header: "Gap", Align: "center"},
			{Field: "slosMet", Header: "SLOs Met", Align: "center"},
			{Field: "slosTotal", Header: "SLOs Total", Align: "center"},
			{Field: "slosMetPercent", Header: "% Met", Align: "center", Format: "percent"},
			{Field: "status", Header: "Status", Align: "center"},
		},
		Sortable: true,
		Striped:  true,
		Pagination: &dashboardir.TablePagination{
			Enabled:  true,
			PageSize: 10,
		},
	}
}

// createSLODetailTable creates the SLO requirements table config.
func createSLODetailTable() dashboardir.TableConfig {
	return dashboardir.TableConfig{
		Columns: []dashboardir.TableColumn{
			{Field: "category", Header: "Category", Width: "15%"},
			{Field: "metricName", Header: "Metric", Width: "25%"},
			{Field: "level", Header: "Level", Align: "center"},
			{Field: "levelName", Header: "Level Name", Width: "12%"},
			{Field: "requirement", Header: "Requirement", Align: "center"},
			{Field: "goalName", Header: "Goal", Width: "20%"},
		},
		Sortable: true,
		Striped:  true,
		Pagination: &dashboardir.TablePagination{
			Enabled:  true,
			PageSize: 15,
		},
	}
}

// createPhaseDetailTable creates the phase progress table config.
func createPhaseDetailTable() dashboardir.TableConfig {
	return dashboardir.TableConfig{
		Columns: []dashboardir.TableColumn{
			{Field: "phaseName", Header: "Phase", Width: "20%"},
			{Field: "status", Header: "Status", Align: "center"},
			{Field: "completionPct", Header: "Progress", Align: "center", Format: "percent"},
			{Field: "goalsTargeted", Header: "Goals", Align: "center"},
			{Field: "goalsAchieved", Header: "Achieved", Align: "center"},
			{Field: "initCompleted", Header: "Initiatives Done", Align: "center"},
			{Field: "initTotal", Header: "Total Initiatives", Align: "center"},
		},
		Sortable: true,
		Striped:  true,
	}
}

// createGapDetailTable creates the gap analysis table config.
func createGapDetailTable() dashboardir.TableConfig {
	return dashboardir.TableConfig{
		Columns: []dashboardir.TableColumn{
			{Field: "priority", Header: "Priority", Width: "10%", Align: "center"},
			{Field: "metricName", Header: "Metric", Width: "25%"},
			{Field: "category", Header: "Category", Width: "12%"},
			{Field: "currentVal", Header: "Current", Align: "right", Format: "number"},
			{Field: "targetVal", Header: "Target", Align: "right", Format: "number"},
			{Field: "gapPercent", Header: "Gap %", Align: "right", Format: "percent"},
			{Field: "goalName", Header: "Goal", Width: "20%"},
		},
		Sortable: true,
		Striped:  true,
		Pagination: &dashboardir.TablePagination{
			Enabled:  true,
			PageSize: 10,
		},
	}
}

// createGoalSLOTable creates the goal-specific SLO table config.
func createGoalSLOTable() dashboardir.TableConfig {
	return dashboardir.TableConfig{
		Columns: []dashboardir.TableColumn{
			{Field: "level", Header: "Level", Align: "center", Width: "10%"},
			{Field: "levelName", Header: "Level Name", Width: "15%"},
			{Field: "metricName", Header: "Metric", Width: "35%"},
			{Field: "requirement", Header: "Requirement", Align: "center", Width: "20%"},
		},
		Sortable: true,
		Striped:  true,
	}
}
