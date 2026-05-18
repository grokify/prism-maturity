package export

import (
	"fmt"
	"time"

	"github.com/grokify/prism-intelligence"
	"github.com/grokify/prism-intelligence/output"
)

// V2MOMDocument represents a structured-plan V2MOM document.
type V2MOMDocument struct {
	Schema    string          `json:"$schema,omitempty"`
	Metadata  V2MOMMetadata   `json:"metadata"`
	Vision    string          `json:"vision"`
	Values    []V2MOMValue    `json:"values,omitempty"`
	Methods   []V2MOMMethod   `json:"methods"`
	Obstacles []V2MOMObstacle `json:"obstacles,omitempty"`
	Measures  []V2MOMMeasure  `json:"measures,omitempty"`
}

// V2MOMMetadata holds V2MOM document metadata.
type V2MOMMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Author      string `json:"author,omitempty"`
	Team        string `json:"team,omitempty"`
	FiscalYear  int    `json:"fiscalYear,omitempty"`
	Quarter     string `json:"quarter,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	Source      string `json:"source,omitempty"`
}

// V2MOMValue represents a V2MOM value.
type V2MOMValue struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"`
}

// V2MOMMethod represents a V2MOM method (like an OKR objective).
type V2MOMMethod struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Priority    string          `json:"priority,omitempty"`
	Status      string          `json:"status,omitempty"`
	Owner       string          `json:"owner,omitempty"`
	StartDate   string          `json:"startDate,omitempty"`
	EndDate     string          `json:"endDate,omitempty"`
	Measures    []V2MOMMeasure  `json:"measures,omitempty"`
	Obstacles   []V2MOMObstacle `json:"obstacles,omitempty"`
	Projects    []string        `json:"projects,omitempty"`
}

// V2MOMMeasure represents a V2MOM measure (like an OKR key result).
type V2MOMMeasure struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Baseline    string  `json:"baseline,omitempty"`
	Target      string  `json:"target"`
	Current     string  `json:"current,omitempty"`
	Unit        string  `json:"unit,omitempty"`
	Progress    float64 `json:"progress,omitempty"`
	Timeline    string  `json:"timeline,omitempty"`
	Status      string  `json:"status,omitempty"`
}

// V2MOMObstacle represents a V2MOM obstacle.
type V2MOMObstacle struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Impact      string `json:"impact,omitempty"`
	Mitigation  string `json:"mitigation,omitempty"`
	Status      string `json:"status,omitempty"`
}

// ConvertToV2MOM converts a PRISM document to a V2MOM document.
func ConvertToV2MOM(doc *prism.PRISMDocument) *V2MOMDocument {
	now := time.Now().Format(time.RFC3339)

	v2momDoc := &V2MOMDocument{
		Schema: "https://github.com/grokify/structured-plan/schema/v2mom.schema.json",
		Metadata: V2MOMMetadata{
			ID:          fmt.Sprintf("v2mom-%s", doc.Metadata.Name),
			Name:        fmt.Sprintf("%s V2MOM", doc.Metadata.Name),
			Description: "Generated from PRISM document",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
			Source:      "prism",
		},
	}

	if doc.Metadata.Author != "" {
		v2momDoc.Metadata.Author = doc.Metadata.Author
	}

	// Determine fiscal year from phases
	phases := doc.GetPhasesSorted()
	if len(phases) > 0 {
		v2momDoc.Metadata.FiscalYear = phases[0].Year
		v2momDoc.Metadata.Quarter = phases[0].Quarter
	}

	// Generate vision from goals
	if len(doc.Goals) > 0 {
		v2momDoc.Vision = fmt.Sprintf("Achieve operational excellence through %d strategic goals, progressing organizational maturity across reliability, delivery, and quality domains.", len(doc.Goals))
	}

	// Default values
	v2momDoc.Values = []V2MOMValue{
		{ID: "value-1", Name: "Reliability First", Description: "Prioritize system reliability and customer trust", Priority: 1},
		{ID: "value-2", Name: "Continuous Improvement", Description: "Systematically progress maturity levels", Priority: 2},
		{ID: "value-3", Name: "Data-Driven Decisions", Description: "Use SLOs and metrics to guide priorities", Priority: 3},
	}

	// Convert goals to methods
	for i, goal := range doc.Goals {
		method := GoalToMethod(goal, doc, i+1)
		v2momDoc.Methods = append(v2momDoc.Methods, method)
	}

	return v2momDoc
}

// GoalToMethod converts a PRISM goal to a V2MOM method.
func GoalToMethod(goal prism.Goal, doc *prism.PRISMDocument, priority int) V2MOMMethod {
	currentLevel := goal.CurrentLevel
	if currentLevel == 0 {
		currentLevel = goal.CurrentMaturityLevel(doc)
	}

	method := V2MOMMethod{
		ID:          fmt.Sprintf("method-%s", goal.ID),
		Name:        goal.Name,
		Description: goal.Description,
		Priority:    fmt.Sprintf("P%d", priority),
		Status:      goal.Status,
		Owner:       goal.Owner,
	}

	if method.Status == "" {
		method.Status = "in_progress"
	}

	// Get date range from phases
	phases := doc.GetPhasesSorted()
	for _, phase := range phases {
		for _, gt := range phase.GoalTargets {
			if gt.GoalID == goal.ID {
				if method.StartDate == "" && phase.StartDate != "" {
					method.StartDate = phase.StartDate
				}
				if phase.EndDate != "" {
					method.EndDate = phase.EndDate
				}
			}
		}
	}

	// Convert SLOs to measures
	if goal.MaturityModel != nil {
		measureIndex := 1
		for level := currentLevel + 1; level <= goal.TargetLevel; level++ {
			for _, ml := range goal.MaturityModel.Levels {
				if ml.Level == level {
					for _, slo := range ml.RequiredSLOs {
						measure := SLOToMeasure(slo, level, doc, measureIndex)
						method.Measures = append(method.Measures, measure)
						measureIndex++
					}
					for _, mc := range ml.MetricCriteria {
						measure := CriterionToMeasure(mc, level, doc, measureIndex)
						method.Measures = append(method.Measures, measure)
						measureIndex++
					}
				}
			}
		}
	}

	// Add initiatives as projects
	for _, init := range doc.Initiatives {
		for _, gid := range init.GoalIDs {
			if gid == goal.ID {
				method.Projects = append(method.Projects, init.ID)
			}
		}
	}

	return method
}

// SLOToMeasure converts an SLO requirement to a V2MOM measure.
func SLOToMeasure(slo prism.SLORequirement, level int, doc *prism.PRISMDocument, index int) V2MOMMeasure {
	metric := doc.GetMetricByID(slo.MetricID)

	measure := V2MOMMeasure{
		ID:   fmt.Sprintf("measure-%d", index),
		Name: fmt.Sprintf("[M%d] %s", level, slo.MetricID),
	}

	if metric != nil {
		measure.Name = fmt.Sprintf("[M%d] %s", level, metric.Name)
		measure.Description = metric.Description
		measure.Baseline = fmt.Sprintf("%.2f", metric.Baseline)
		measure.Current = fmt.Sprintf("%.2f", metric.Current)
		measure.Unit = metric.Unit

		if metric.SLO != nil {
			measure.Target = metric.SLO.Target
			measure.Timeline = metric.SLO.Window
		}

		if metric.MeetsSLO() {
			measure.Status = "achieved"
			measure.Progress = 1.0
		} else {
			measure.Status = "in_progress"
			if metric.Target > metric.Baseline {
				measure.Progress = (metric.Current - metric.Baseline) / (metric.Target - metric.Baseline)
				if measure.Progress < 0 {
					measure.Progress = 0
				}
			}
		}
	}

	return measure
}

// CriterionToMeasure converts a metric criterion to a V2MOM measure.
func CriterionToMeasure(mc prism.MetricCriterion, level int, doc *prism.PRISMDocument, index int) V2MOMMeasure {
	metric := doc.GetMetricByID(mc.MetricID)

	measure := V2MOMMeasure{
		ID:     fmt.Sprintf("measure-%d", index),
		Name:   fmt.Sprintf("[M%d] %s %s %.2f", level, mc.MetricID, output.OperatorSymbol(mc.Operator), mc.Value),
		Target: fmt.Sprintf("%s%.2f", output.OperatorSymbol(mc.Operator), mc.Value),
	}

	if metric != nil {
		measure.Name = fmt.Sprintf("[M%d] %s %s %.2f%s", level, metric.Name, output.OperatorSymbol(mc.Operator), mc.Value, metric.Unit)
		measure.Description = metric.Description
		measure.Baseline = fmt.Sprintf("%.2f", metric.Baseline)
		measure.Current = fmt.Sprintf("%.2f", metric.Current)
		measure.Unit = metric.Unit

		if mc.IsMet(metric.Current) {
			measure.Status = "achieved"
			measure.Progress = 1.0
		} else {
			measure.Status = "in_progress"
			measure.Progress = CalculateCriterionProgress(mc, metric)
		}
	}

	return measure
}
