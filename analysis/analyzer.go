// Package analysis provides analysis and reporting capabilities for PRISM documents.
package analysis

import (
	"fmt"

	"github.com/grokify/prism-maturity"
	"github.com/grokify/prism-maturity/output"
)

// Result holds the structured analysis of a PRISM document.
type Result struct {
	Summary         Summary          `json:"summary"`
	Goals           []GoalAnalysis   `json:"goals"`
	Phases          []PhaseAnalysis  `json:"phases"`
	Gaps            []Gap            `json:"gaps"`
	Recommendations []Recommendation `json:"recommendations,omitempty"`
}

// Summary provides high-level statistics.
type Summary struct {
	TotalGoals     int     `json:"totalGoals"`
	TotalPhases    int     `json:"totalPhases"`
	TotalSLOs      int     `json:"totalSLOs"`
	SLOsMet        int     `json:"slosMet"`
	SLOCompliance  float64 `json:"sloCompliance"`
	AvgMaturityGap float64 `json:"avgMaturityGap"`
}

// GoalAnalysis analyzes a single goal.
type GoalAnalysis struct {
	GoalID       string           `json:"goalId"`
	GoalName     string           `json:"goalName"`
	CurrentLevel int              `json:"currentLevel"`
	TargetLevel  int              `json:"targetLevel"`
	Gap          int              `json:"gap"`
	Status       string           `json:"status"`
	SLOsRequired []SLORequirement `json:"slosRequired"`
	SLOsMet      int              `json:"slosMet"`
	SLOsTotal    int              `json:"slosTotal"`
}

// SLORequirement captures an SLO needed for a maturity level.
type SLORequirement struct {
	MetricID   string  `json:"metricId"`
	MetricName string  `json:"metricName"`
	Level      int     `json:"level"`
	Target     string  `json:"target"`
	Current    float64 `json:"current"`
	IsMet      bool    `json:"isMet"`
}

// PhaseAnalysis analyzes a single phase.
type PhaseAnalysis struct {
	PhaseID     string       `json:"phaseId"`
	PhaseName   string       `json:"phaseName"`
	Period      string       `json:"period"`
	Status      string       `json:"status"`
	GoalTargets []GoalTarget `json:"goalTargets"`
	Initiatives int          `json:"initiatives"`
	Completion  float64      `json:"completion"`
}

// GoalTarget shows maturity progression in a phase.
type GoalTarget struct {
	GoalID     string `json:"goalId"`
	GoalName   string `json:"goalName"`
	EnterLevel int    `json:"enterLevel"`
	ExitLevel  int    `json:"exitLevel"`
	SLOsNeeded int    `json:"slosNeeded"`
}

// Recommendation suggests an initiative to close a gap.
type Recommendation struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	GoalIDs     []string `json:"goalIds"`
	PhaseID     string   `json:"phaseId"`
	Priority    string   `json:"priority"`
	SLOsEnabled []string `json:"slosEnabled"`
}

// Analyze performs a comprehensive analysis of a PRISM document.
func Analyze(doc *prism.PRISMDocument) *Result {
	result := &Result{}

	// Analyze goals
	var totalGap float64
	for _, goal := range doc.Goals {
		currentLevel := goal.CurrentLevel
		if currentLevel == 0 {
			currentLevel = goal.CurrentMaturityLevel(doc)
		}

		gap := goal.TargetLevel - currentLevel
		totalGap += float64(gap)

		ga := GoalAnalysis{
			GoalID:       goal.ID,
			GoalName:     goal.Name,
			CurrentLevel: currentLevel,
			TargetLevel:  goal.TargetLevel,
			Gap:          gap,
			Status:       output.GoalStatus(currentLevel, goal.TargetLevel),
		}

		// Collect SLO requirements for levels between current and target
		for level := currentLevel + 1; level <= goal.TargetLevel; level++ {
			if goal.MaturityModel != nil {
				for _, ml := range goal.MaturityModel.Levels {
					if ml.Level == level {
						for _, slo := range ml.RequiredSLOs {
							metric := doc.GetMetricByID(slo.MetricID)
							req := SLORequirement{
								MetricID: slo.MetricID,
								Level:    level,
							}
							if metric != nil {
								req.MetricName = metric.Name
								req.Current = metric.Current
								req.IsMet = metric.MeetsSLO()
								if metric.SLO != nil {
									req.Target = metric.SLO.Target
								}
							}
							ga.SLOsRequired = append(ga.SLOsRequired, req)
							ga.SLOsTotal++
							if req.IsMet {
								ga.SLOsMet++
							}
						}
					}
				}
			}
		}

		result.Goals = append(result.Goals, ga)
	}

	// Analyze phases
	phases := doc.GetPhasesSorted()
	for _, phase := range phases {
		period := ""
		if phase.Quarter != "" {
			period = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
		}

		pa := PhaseAnalysis{
			PhaseID:   phase.ID,
			PhaseName: phase.Name,
			Period:    period,
			Status:    phase.Status,
		}

		if pa.Status == "" {
			pa.Status = "planned"
		}

		// Analyze goal targets in this phase
		for _, gt := range phase.GoalTargets {
			goal := doc.GetGoalByID(gt.GoalID)
			name := gt.GoalID
			if goal != nil {
				name = goal.Name
			}

			// Count SLOs needed for this phase's maturity jump
			slosNeeded := CountSLOsForLevelRange(goal, gt.EnterLevel, gt.ExitLevel)

			pa.GoalTargets = append(pa.GoalTargets, GoalTarget{
				GoalID:     gt.GoalID,
				GoalName:   name,
				EnterLevel: gt.EnterLevel,
				ExitLevel:  gt.ExitLevel,
				SLOsNeeded: slosNeeded,
			})
		}

		// Count initiatives
		for _, init := range doc.Initiatives {
			if init.PhaseID == phase.ID {
				pa.Initiatives++
			}
		}

		// Calculate completion
		view := doc.GeneratePhaseRoadmapView(phase.ID)
		if view != nil {
			pa.Completion = view.OverallCompletion
		}

		result.Phases = append(result.Phases, pa)
	}

	// Identify gaps
	result.Gaps = IdentifyGaps(result)

	// Summary
	totalSLOs := 0
	slosMet := 0
	for _, m := range doc.Metrics {
		if m.SLO != nil {
			totalSLOs++
			if m.MeetsSLO() {
				slosMet++
			}
		}
	}

	sloCompliance := 0.0
	if totalSLOs > 0 {
		sloCompliance = float64(slosMet) / float64(totalSLOs) * 100
	}

	avgGap := 0.0
	if len(doc.Goals) > 0 {
		avgGap = totalGap / float64(len(doc.Goals))
	}

	result.Summary = Summary{
		TotalGoals:     len(doc.Goals),
		TotalPhases:    len(phases),
		TotalSLOs:      totalSLOs,
		SLOsMet:        slosMet,
		SLOCompliance:  sloCompliance,
		AvgMaturityGap: avgGap,
	}

	return result
}

// CountSLOsForLevelRange counts SLOs required between enter and exit levels.
func CountSLOsForLevelRange(goal *prism.Goal, enterLevel, exitLevel int) int {
	if goal == nil || goal.MaturityModel == nil {
		return 0
	}

	count := 0
	for _, ml := range goal.MaturityModel.Levels {
		if ml.Level > enterLevel && ml.Level <= exitLevel {
			count += len(ml.RequiredSLOs)
		}
	}
	return count
}
