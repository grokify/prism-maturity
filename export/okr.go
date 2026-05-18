// Package export provides converters for exporting PRISM data to other formats.
package export

import (
	"fmt"
	"time"

	"github.com/grokify/prism-intelligence"
	"github.com/grokify/prism-intelligence/output"
)

// OKRDocument represents a structured-plan OKR document.
type OKRDocument struct {
	Schema     string        `json:"$schema,omitempty"`
	Metadata   OKRMetadata   `json:"metadata"`
	Theme      string        `json:"theme,omitempty"`
	Objectives []Objective   `json:"objectives"`
	Alignment  *OKRAlignment `json:"alignment,omitempty"`
}

// OKRMetadata holds document metadata.
type OKRMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Team        string `json:"team,omitempty"`
	Period      string `json:"period,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	Source      string `json:"source,omitempty"`
}

// OKRAlignment links to parent OKRs.
type OKRAlignment struct {
	ParentOKRID string   `json:"parentOkrId,omitempty"`
	CompanyOKRs []string `json:"companyOkrs,omitempty"`
}

// Objective represents an OKR objective.
type Objective struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Rationale   string      `json:"rationale,omitempty"`
	Category    string      `json:"category,omitempty"`
	Owner       string      `json:"owner,omitempty"`
	Timeframe   string      `json:"timeframe,omitempty"`
	Status      string      `json:"status,omitempty"`
	KeyResults  []KeyResult `json:"keyResults"`
	Progress    float64     `json:"progress,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
}

// KeyResult represents an OKR key result.
type KeyResult struct {
	ID                string        `json:"id"`
	Title             string        `json:"title"`
	Description       string        `json:"description,omitempty"`
	Owner             string        `json:"owner,omitempty"`
	Metric            string        `json:"metric,omitempty"`
	Baseline          string        `json:"baseline,omitempty"`
	Target            string        `json:"target"`
	Current           string        `json:"current,omitempty"`
	Unit              string        `json:"unit,omitempty"`
	MeasurementMethod string        `json:"measurementMethod,omitempty"`
	Score             float64       `json:"score,omitempty"`
	Confidence        string        `json:"confidence,omitempty"`
	Status            string        `json:"status,omitempty"`
	DueDate           string        `json:"dueDate,omitempty"`
	PhaseTargets      []PhaseTarget `json:"phaseTargets,omitempty"`
	Tags              []string      `json:"tags,omitempty"`
}

// PhaseTarget represents a per-phase target for a key result.
type PhaseTarget struct {
	PhaseID string `json:"phaseId"`
	Target  string `json:"target"`
	Status  string `json:"status,omitempty"`
	Actual  string `json:"actual,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

// ConvertToOKR converts a PRISM document to an OKR document.
func ConvertToOKR(doc *prism.PRISMDocument) *OKRDocument {
	now := time.Now().Format(time.RFC3339)

	okrDoc := &OKRDocument{
		Schema: "https://github.com/grokify/structured-plan/schema/okr.schema.json",
		Metadata: OKRMetadata{
			ID:          fmt.Sprintf("okr-%s", doc.Metadata.Name),
			Name:        fmt.Sprintf("%s OKRs", doc.Metadata.Name),
			Description: "Generated from PRISM document",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
			Source:      "prism",
		},
	}

	if doc.Metadata.Author != "" {
		okrDoc.Metadata.Owner = doc.Metadata.Author
	}

	// Determine period from phases
	phases := doc.GetPhasesSorted()
	if len(phases) > 0 {
		first := phases[0]
		last := phases[len(phases)-1]
		if first.Quarter != "" && last.Quarter != "" {
			okrDoc.Metadata.Period = fmt.Sprintf("%s %d - %s %d", first.Quarter, first.Year, last.Quarter, last.Year)
		}
	}

	// Convert goals to objectives
	for _, goal := range doc.Goals {
		obj := GoalToObjective(goal, doc)
		okrDoc.Objectives = append(okrDoc.Objectives, obj)
	}

	return okrDoc
}

// GoalToObjective converts a PRISM goal to an OKR objective.
func GoalToObjective(goal prism.Goal, doc *prism.PRISMDocument) Objective {
	currentLevel := goal.CurrentLevel
	if currentLevel == 0 {
		currentLevel = goal.CurrentMaturityLevel(doc)
	}

	obj := Objective{
		ID:          fmt.Sprintf("obj-%s", goal.ID),
		Title:       goal.Name,
		Description: goal.Description,
		Category:    "Operational Excellence",
		Owner:       goal.Owner,
		Status:      goal.Status,
		Tags:        []string{"prism", "maturity"},
	}

	if obj.Status == "" {
		obj.Status = "active"
	}

	// Add rationale based on maturity gap
	obj.Rationale = fmt.Sprintf("Progress from M%d to M%d maturity level", currentLevel, goal.TargetLevel)

	// Get phases that include this goal
	phases := doc.GetPhasesSorted()
	for _, phase := range phases {
		for _, gt := range phase.GoalTargets {
			if gt.GoalID == goal.ID && phase.EndDate != "" {
				obj.Timeframe = phase.EndDate[:10] // Use last phase end date
			}
		}
	}

	// Convert SLOs to key results
	if goal.MaturityModel != nil {
		krIndex := 1
		for level := currentLevel + 1; level <= goal.TargetLevel; level++ {
			for _, ml := range goal.MaturityModel.Levels {
				if ml.Level == level {
					// Add key results for required SLOs
					for _, slo := range ml.RequiredSLOs {
						kr := SLOToKeyResult(slo, level, doc, phases, goal.ID, krIndex)
						obj.KeyResults = append(obj.KeyResults, kr)
						krIndex++
					}

					// Add key results for metric criteria
					for _, mc := range ml.MetricCriteria {
						kr := MetricCriteriaToKeyResult(mc, level, doc, phases, goal.ID, krIndex)
						obj.KeyResults = append(obj.KeyResults, kr)
						krIndex++
					}
				}
			}
		}
	}

	// Calculate progress
	if len(obj.KeyResults) > 0 {
		var totalScore float64
		for _, kr := range obj.KeyResults {
			totalScore += kr.Score
		}
		obj.Progress = totalScore / float64(len(obj.KeyResults))
	}

	return obj
}

// SLOToKeyResult converts an SLO requirement to a key result.
func SLOToKeyResult(slo prism.SLORequirement, level int, doc *prism.PRISMDocument, phases []prism.Phase, goalID string, index int) KeyResult {
	metric := doc.GetMetricByID(slo.MetricID)

	kr := KeyResult{
		ID:     fmt.Sprintf("kr-%s-%d", goalID, index),
		Title:  fmt.Sprintf("[M%d] %s", level, slo.MetricID),
		Metric: slo.MetricID,
		Tags:   []string{fmt.Sprintf("M%d", level)},
	}

	if metric != nil {
		kr.Title = fmt.Sprintf("[M%d] %s", level, metric.Name)
		kr.Description = metric.Description
		kr.Baseline = fmt.Sprintf("%.2f", metric.Baseline)
		kr.Current = fmt.Sprintf("%.2f", metric.Current)
		kr.Unit = metric.Unit

		if metric.SLO != nil {
			kr.Target = metric.SLO.Target
			kr.MeasurementMethod = fmt.Sprintf("SLO window: %s", metric.SLO.Window)
		}

		// Calculate score
		if metric.MeetsSLO() {
			kr.Score = 1.0
			kr.Status = "achieved"
			kr.Confidence = "high"
		} else {
			// Calculate partial progress
			if metric.Target > metric.Baseline {
				progress := (metric.Current - metric.Baseline) / (metric.Target - metric.Baseline)
				if progress < 0 {
					progress = 0
				} else if progress > 1 {
					progress = 1
				}
				kr.Score = progress
			}
			kr.Status = "in_progress"
			kr.Confidence = "medium"
		}
	}

	// Add phase targets
	kr.PhaseTargets = GeneratePhaseTargets(phases, goalID, level)

	return kr
}

// MetricCriteriaToKeyResult converts a metric criterion to a key result.
func MetricCriteriaToKeyResult(mc prism.MetricCriterion, level int, doc *prism.PRISMDocument, phases []prism.Phase, goalID string, index int) KeyResult {
	metric := doc.GetMetricByID(mc.MetricID)

	kr := KeyResult{
		ID:     fmt.Sprintf("kr-%s-%d", goalID, index),
		Title:  fmt.Sprintf("[M%d] %s %s %.2f", level, mc.MetricID, output.OperatorSymbol(mc.Operator), mc.Value),
		Metric: mc.MetricID,
		Target: fmt.Sprintf("%s%.2f", output.OperatorSymbol(mc.Operator), mc.Value),
		Tags:   []string{fmt.Sprintf("M%d", level)},
	}

	if metric != nil {
		kr.Title = fmt.Sprintf("[M%d] %s %s %.2f%s", level, metric.Name, output.OperatorSymbol(mc.Operator), mc.Value, metric.Unit)
		kr.Description = metric.Description
		kr.Baseline = fmt.Sprintf("%.2f", metric.Baseline)
		kr.Current = fmt.Sprintf("%.2f", metric.Current)
		kr.Unit = metric.Unit

		// Check if criterion is met
		if mc.IsMet(metric.Current) {
			kr.Score = 1.0
			kr.Status = "achieved"
			kr.Confidence = "high"
		} else {
			kr.Status = "in_progress"
			kr.Confidence = "medium"
			// Calculate progress toward criterion
			kr.Score = CalculateCriterionProgress(mc, metric)
		}
	}

	// Add phase targets
	kr.PhaseTargets = GeneratePhaseTargets(phases, goalID, level)

	return kr
}

// CalculateCriterionProgress calculates progress toward a metric criterion.
func CalculateCriterionProgress(mc prism.MetricCriterion, metric *prism.Metric) float64 {
	if metric == nil {
		return 0
	}

	current := metric.Current
	baseline := metric.Baseline
	target := mc.Value

	switch mc.Operator {
	case "gte", "gt":
		if target <= baseline {
			return 0
		}
		progress := (current - baseline) / (target - baseline)
		if progress < 0 {
			return 0
		}
		if progress > 1 {
			return 1
		}
		return progress
	case "lte", "lt":
		if target >= baseline {
			return 0
		}
		progress := (baseline - current) / (baseline - target)
		if progress < 0 {
			return 0
		}
		if progress > 1 {
			return 1
		}
		return progress
	default:
		return 0
	}
}

// GeneratePhaseTargets generates phase targets for a key result.
func GeneratePhaseTargets(phases []prism.Phase, goalID string, targetLevel int) []PhaseTarget {
	var targets []PhaseTarget

	for _, phase := range phases {
		for _, gt := range phase.GoalTargets {
			if gt.GoalID == goalID {
				// Check if this phase progresses toward the target level
				if gt.ExitLevel >= targetLevel {
					status := "not_started"
					if phase.Status == "in_progress" {
						status = "in_progress"
					} else if phase.Status == "completed" {
						status = "achieved"
					}

					targets = append(targets, PhaseTarget{
						PhaseID: phase.ID,
						Target:  fmt.Sprintf("M%d", gt.ExitLevel),
						Status:  status,
						Notes:   fmt.Sprintf("Progress from M%d to M%d", gt.EnterLevel, gt.ExitLevel),
					})
				}
			}
		}
	}

	return targets
}
