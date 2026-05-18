// Package export provides converters for exporting PRISM data to other formats.
package export

import (
	"fmt"
	"time"

	"github.com/grokify/prism-execution/goals/okr"
	"github.com/grokify/prism-execution/roadmap"
	"github.com/grokify/prism-intelligence"
	"github.com/grokify/prism-intelligence/output"
)

// RoadmapExport contains both the roadmap and associated OKRs.
// This combined export captures the full PRISM model:
// - Roadmap: phases, deliverables, rollout tracking
// - OKRs: maturity-based objectives with SLO key results
type RoadmapExport struct {
	Roadmap *roadmap.Roadmap `json:"roadmap"`
	OKRs    *okr.OKRDocument `json:"okrs,omitempty"`
}

// ConvertToRoadmap converts a PRISM document to a structured-plan roadmap.
func ConvertToRoadmap(doc *prism.PRISMDocument) *roadmap.Roadmap {
	rm := &roadmap.Roadmap{
		Phases: make([]roadmap.Phase, 0, len(doc.Phases)),
	}

	for _, phase := range doc.GetPhasesSorted() {
		rm.Phases = append(rm.Phases, convertPhase(phase, doc))
	}

	return rm
}

// ConvertToRoadmapWithOKRs converts a PRISM document to both roadmap and OKRs.
// This provides the complete export with:
// - Phases → roadmap.Phase
// - Initiatives → roadmap.Deliverable with RolloutStatus
// - Goals/Maturity → OKR Objectives
// - SLOs → OKR Key Results
func ConvertToRoadmapWithOKRs(doc *prism.PRISMDocument) *RoadmapExport {
	return &RoadmapExport{
		Roadmap: ConvertToRoadmap(doc),
		OKRs:    ConvertToStructuredOKR(doc),
	}
}

// convertPhase converts a PRISM phase to a structured-plan phase.
func convertPhase(phase prism.Phase, doc *prism.PRISMDocument) roadmap.Phase {
	p := roadmap.Phase{
		ID:     phase.ID,
		Name:   phase.Name,
		Type:   convertPhaseType(phase),
		Status: convertPhaseStatus(phase.Status),
		Goals:  make([]string, 0),
		Tags:   []string{"prism"},
	}

	// Set dates
	if phase.StartDate != "" {
		if t, err := time.Parse("2006-01-02", phase.StartDate); err == nil {
			p.StartDate = &t
		}
	}
	if phase.EndDate != "" {
		if t, err := time.Parse("2006-01-02", phase.EndDate); err == nil {
			p.EndDate = &t
		}
	}

	// Convert goal targets to goals list with maturity info
	for _, gt := range phase.GoalTargets {
		goal := doc.GetGoalByID(gt.GoalID)
		goalName := gt.GoalID
		if goal != nil {
			goalName = goal.Name
		}
		p.Goals = append(p.Goals, fmt.Sprintf("%s (M%d→M%d)", goalName, gt.EnterLevel, gt.ExitLevel))
	}

	// Add success criteria from goal targets
	for _, gt := range phase.GoalTargets {
		p.SuccessCriteria = append(p.SuccessCriteria,
			fmt.Sprintf("Achieve M%d for %s", gt.ExitLevel, gt.GoalID))
	}

	// Convert initiatives to deliverables
	phaseInitiatives := doc.GetInitiativesForPhase(phase.ID)
	for _, init := range phaseInitiatives {
		p.Deliverables = append(p.Deliverables, convertInitiative(init, doc))
	}

	// Calculate progress
	if len(phaseInitiatives) > 0 {
		var totalProgress float64
		for _, init := range phaseInitiatives {
			totalProgress += init.DevCompletionPercent
		}
		progress := int(totalProgress / float64(len(phaseInitiatives)))
		p.Progress = &progress
	}

	return p
}

// convertPhaseType determines the phase type from PRISM phase data.
func convertPhaseType(phase prism.Phase) roadmap.PhaseType {
	if phase.Quarter != "" {
		return roadmap.PhaseTypeQuarter
	}
	return roadmap.PhaseTypeGeneric
}

// convertPhaseStatus converts PRISM status to roadmap status.
func convertPhaseStatus(status string) roadmap.PhaseStatus {
	switch status {
	case "completed":
		return roadmap.PhaseStatusCompleted
	case "in_progress":
		return roadmap.PhaseStatusInProgress
	case "delayed":
		return roadmap.PhaseStatusDelayed
	case "cancelled":
		return roadmap.PhaseStatusCancelled
	default:
		return roadmap.PhaseStatusPlanned
	}
}

// convertInitiative converts a PRISM initiative to a roadmap deliverable.
func convertInitiative(init prism.Initiative, _ *prism.PRISMDocument) roadmap.Deliverable {
	d := roadmap.Deliverable{
		ID:          init.ID,
		Title:       init.Name,
		Description: init.Description,
		Type:        roadmap.DeliverableFeature,
		Status:      convertDeliverableStatus(init.Status),
		Tags:        []string{"prism", "initiative"},
	}

	// Add goal tags
	for _, goalID := range init.GoalIDs {
		d.Tags = append(d.Tags, fmt.Sprintf("goal:%s", goalID))
	}

	// Convert deployment status to rollout
	if init.DeploymentStatus != nil {
		ds := init.DeploymentStatus
		d.Rollout = &roadmap.RolloutStatus{
			TotalCustomers:    ds.TotalCustomers,
			DeployedCustomers: ds.DeployedCustomers,
			AdoptedCustomers:  calculateAdoptedCustomers(ds),
			Status:            convertRolloutStage(ds.Status, init.DevCompletionPercent),
			Notes:             fmt.Sprintf("Dev completion: %.0f%%", init.DevCompletionPercent),
		}

		// Add start/target dates if available
		if init.StartDate != "" {
			d.Rollout.StartDate = init.StartDate
		}
		if init.EndDate != "" {
			d.Rollout.TargetDate = init.EndDate
		}
	} else if init.DevCompletionPercent > 0 {
		// Create rollout status from dev completion if no deployment status
		d.Rollout = &roadmap.RolloutStatus{
			Status: convertRolloutStage("", init.DevCompletionPercent),
			Notes:  fmt.Sprintf("Dev completion: %.0f%%", init.DevCompletionPercent),
		}
	}

	return d
}

// convertDeliverableStatus converts PRISM initiative status to deliverable status.
func convertDeliverableStatus(status string) roadmap.DeliverableStatus {
	switch status {
	case prism.InitiativeStatusCompleted:
		return roadmap.DeliverableCompleted
	case prism.InitiativeStatusInProgress:
		return roadmap.DeliverableInProgress
	case prism.InitiativeStatusNotStarted, prism.InitiativeStatusPlanned:
		return roadmap.DeliverableNotStarted
	default:
		return roadmap.DeliverableNotStarted
	}
}

// convertRolloutStage determines rollout stage from status and completion.
func convertRolloutStage(status string, devCompletion float64) roadmap.RolloutStage {
	switch status {
	case "deployed", "completed":
		return roadmap.RolloutStageDeployed
	case "rolling_out", "in_progress":
		return roadmap.RolloutStageRollingOut
	case "paused":
		return roadmap.RolloutStagePaused
	case "rolled_back":
		return roadmap.RolloutStageRolledBack
	default:
		if devCompletion >= 100 {
			return roadmap.RolloutStageDeployed
		} else if devCompletion > 0 {
			return roadmap.RolloutStageRollingOut
		}
		return roadmap.RolloutStageNotStarted
	}
}

// calculateAdoptedCustomers estimates adopted customers from adoption percent.
func calculateAdoptedCustomers(ds *prism.DeploymentStatus) int {
	if ds.AdoptionPercent > 0 && ds.TotalCustomers > 0 {
		return int(float64(ds.TotalCustomers) * ds.AdoptionPercent / 100)
	}
	return 0
}

// ConvertToStructuredOKR converts a PRISM document to structured-plan OKR format.
// Each maturity level to achieve becomes an Objective.
// SLOs and metric criteria become Key Results.
func ConvertToStructuredOKR(doc *prism.PRISMDocument) *okr.OKRDocument {
	now := time.Now()

	okrDoc := &okr.OKRDocument{
		Schema: "https://github.com/grokify/structured-plan/schema/okr.schema.json",
		Metadata: &okr.Metadata{
			ID:        generateOKRID(doc),
			Name:      fmt.Sprintf("%s Maturity OKRs", doc.Metadata.Name),
			Status:    okr.StatusActive,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Objectives: make([]okr.Objective, 0),
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
			okrDoc.Metadata.Period = fmt.Sprintf("%s %d - %s %d",
				first.Quarter, first.Year, last.Quarter, last.Year)
		}
	}

	// Convert each goal to objectives (one per maturity level to achieve)
	for _, goal := range doc.Goals {
		objectives := goalToObjectives(goal, doc, phases)
		okrDoc.Objectives = append(okrDoc.Objectives, objectives...)
	}

	return okrDoc
}

// generateOKRID creates an OKR document ID from the PRISM document.
func generateOKRID(doc *prism.PRISMDocument) string {
	if doc.Metadata != nil && doc.Metadata.Name != "" {
		return fmt.Sprintf("okr-prism-%s", slugify(doc.Metadata.Name))
	}
	return fmt.Sprintf("okr-prism-%d", time.Now().Unix())
}

// slugify creates a URL-safe slug from a string.
// Only processes ASCII alphanumeric characters and converts to lowercase.
func slugify(s string) string {
	result := make([]byte, 0, len(s))
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			result = append(result, byte(c)) //nolint:gosec // ASCII chars guaranteed to fit in byte
		} else if c >= 'A' && c <= 'Z' {
			result = append(result, byte(c-'A'+'a')) //nolint:gosec // ASCII lowercase guaranteed to fit
		} else if c == ' ' || c == '-' || c == '_' {
			if len(result) > 0 && result[len(result)-1] != '-' {
				result = append(result, '-')
			}
		}
	}
	return string(result)
}

// goalToObjectives converts a PRISM goal to one or more OKR objectives.
// Each maturity level to achieve becomes a separate objective.
func goalToObjectives(goal prism.Goal, doc *prism.PRISMDocument, phases []prism.Phase) []okr.Objective {
	currentLevel := goal.CurrentLevel
	if currentLevel == 0 {
		currentLevel = goal.CurrentMaturityLevel(doc)
	}

	objectives := make([]okr.Objective, 0)

	// Create an objective for each maturity level to achieve
	for level := currentLevel + 1; level <= goal.TargetLevel; level++ {
		obj := okr.Objective{
			ID:          fmt.Sprintf("obj-%s-m%d", goal.ID, level),
			Title:       fmt.Sprintf("%s: Achieve M%d (%s)", goal.Name, level, maturityLevelName(level)),
			Description: goal.Description,
			Category:    "Operational Maturity",
			Owner:       goal.Owner,
			Status:      okr.StatusActive,
			Tags:        []string{"prism", "maturity", fmt.Sprintf("M%d", level), goal.ID},
			KeyResults:  make([]okr.KeyResult, 0),
		}

		// Find the phase that targets this level
		for _, phase := range phases {
			for _, gt := range phase.GoalTargets {
				if gt.GoalID == goal.ID && gt.ExitLevel >= level {
					if phase.EndDate != "" {
						obj.Timeframe = phase.EndDate
					}
					break
				}
			}
		}

		// Add rationale
		obj.Rationale = fmt.Sprintf("Progress from M%d (%s) to M%d (%s)",
			level-1, maturityLevelName(level-1), level, maturityLevelName(level))

		// Convert SLOs to key results for this level
		if goal.MaturityModel != nil {
			for _, ml := range goal.MaturityModel.Levels {
				if ml.Level == level {
					// Add required SLOs
					for i, slo := range ml.RequiredSLOs {
						kr := sloToKeyResult(slo, level, doc, phases, goal.ID, i+1)
						obj.KeyResults = append(obj.KeyResults, kr)
					}

					// Add metric criteria
					for i, mc := range ml.MetricCriteria {
						kr := metricCriteriaToKeyResult(mc, level, doc, phases, goal.ID, len(ml.RequiredSLOs)+i+1)
						obj.KeyResults = append(obj.KeyResults, kr)
					}
				}
			}
		}

		// Calculate progress
		obj.Progress = calculateObjectiveProgress(obj.KeyResults)

		objectives = append(objectives, obj)
	}

	return objectives
}

// sloToKeyResult converts a PRISM SLO requirement to an OKR key result.
func sloToKeyResult(slo prism.SLORequirement, level int, doc *prism.PRISMDocument, phases []prism.Phase, goalID string, index int) okr.KeyResult {
	metric := doc.GetMetricByID(slo.MetricID)

	kr := okr.KeyResult{
		ID:     fmt.Sprintf("kr-%s-m%d-%d", goalID, level, index),
		Title:  fmt.Sprintf("[M%d] %s SLO", level, slo.MetricID),
		Metric: slo.MetricID,
		Tags:   []string{fmt.Sprintf("M%d", level), "slo"},
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
			kr.Confidence = okr.ConfidenceHigh
		} else {
			kr.Score = calculateMetricProgress(metric)
			kr.Status = "in_progress"
			kr.Confidence = okr.ConfidenceMedium
		}
	}

	// Add phase targets
	kr.PhaseTargets = generatePhaseTargets(phases, goalID, level)

	return kr
}

// metricCriteriaToKeyResult converts a PRISM metric criterion to an OKR key result.
func metricCriteriaToKeyResult(mc prism.MetricCriterion, level int, doc *prism.PRISMDocument, phases []prism.Phase, goalID string, index int) okr.KeyResult {
	metric := doc.GetMetricByID(mc.MetricID)

	kr := okr.KeyResult{
		ID:     fmt.Sprintf("kr-%s-m%d-%d", goalID, level, index),
		Title:  fmt.Sprintf("[M%d] %s %s %.2f", level, mc.MetricID, output.OperatorSymbol(mc.Operator), mc.Value),
		Metric: mc.MetricID,
		Target: fmt.Sprintf("%s%.2f", output.OperatorSymbol(mc.Operator), mc.Value),
		Tags:   []string{fmt.Sprintf("M%d", level), "criterion"},
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
			kr.Confidence = okr.ConfidenceHigh
		} else {
			kr.Score = calculateCriterionProgress(mc, metric)
			kr.Status = "in_progress"
			kr.Confidence = okr.ConfidenceMedium
		}
	}

	// Add phase targets
	kr.PhaseTargets = generatePhaseTargets(phases, goalID, level)

	return kr
}

// calculateObjectiveProgress calculates overall objective progress from key results.
func calculateObjectiveProgress(keyResults []okr.KeyResult) float64 {
	if len(keyResults) == 0 {
		return 0
	}
	var total float64
	for _, kr := range keyResults {
		total += kr.Score
	}
	return total / float64(len(keyResults))
}

// calculateMetricProgress calculates progress toward an SLO target.
func calculateMetricProgress(metric *prism.Metric) float64 {
	if metric == nil || metric.Target == metric.Baseline {
		return 0
	}
	progress := (metric.Current - metric.Baseline) / (metric.Target - metric.Baseline)
	if progress < 0 {
		return 0
	}
	if progress > 1 {
		return 1
	}
	return progress
}

// calculateCriterionProgress calculates progress toward a metric criterion.
func calculateCriterionProgress(mc prism.MetricCriterion, metric *prism.Metric) float64 {
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

// generatePhaseTargets creates phase targets for a key result.
func generatePhaseTargets(phases []prism.Phase, goalID string, targetLevel int) []okr.PhaseTarget {
	targets := make([]okr.PhaseTarget, 0)

	for _, phase := range phases {
		for _, gt := range phase.GoalTargets {
			if gt.GoalID == goalID && gt.ExitLevel >= targetLevel {
				status := "not_started"
				switch phase.Status {
				case "in_progress":
					status = "in_progress"
				case "completed":
					status = "achieved"
				}

				targets = append(targets, okr.PhaseTarget{
					PhaseID: phase.ID,
					Target:  fmt.Sprintf("M%d", gt.ExitLevel),
					Status:  status,
					Notes:   fmt.Sprintf("Progress from M%d to M%d", gt.EnterLevel, gt.ExitLevel),
				})
			}
		}
	}

	return targets
}

// maturityLevelName returns the name for a maturity level.
func maturityLevelName(level int) string {
	names := map[int]string{
		1: "Reactive",
		2: "Basic",
		3: "Defined",
		4: "Managed",
		5: "Optimizing",
	}
	if name, ok := names[level]; ok {
		return name
	}
	return fmt.Sprintf("Level %d", level)
}
