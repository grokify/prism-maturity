package analysis

import (
	"fmt"

	"github.com/grokify/prism-intelligence"
	"github.com/grokify/prism-intelligence/output"
)

// RoadmapProgress holds overall roadmap progress metrics.
type RoadmapProgress struct {
	OverallCompletion float64                `json:"overallCompletion"`
	PhaseProgress     []PhaseProgressSummary `json:"phaseProgress"`
	GoalProgress      []GoalProgressSummary  `json:"goalProgress"`
}

// PhaseProgressSummary holds summary progress for a phase.
type PhaseProgressSummary struct {
	PhaseID     string  `json:"phaseId"`
	PhaseName   string  `json:"phaseName"`
	Period      string  `json:"period"`
	Status      string  `json:"status"`
	GoalSummary string  `json:"goalSummary"`
	Completion  float64 `json:"completion"`
}

// GoalProgressSummary holds summary progress for a goal.
type GoalProgressSummary struct {
	GoalID       string `json:"goalId"`
	GoalName     string `json:"goalName"`
	CurrentLevel int    `json:"currentLevel"`
	TargetLevel  int    `json:"targetLevel"`
	Status       string `json:"status"`
}

// CalculateRoadmapProgress calculates overall roadmap progress.
func CalculateRoadmapProgress(doc *prism.PRISMDocument) *RoadmapProgress {
	progress := &RoadmapProgress{}

	// Calculate phase progress
	phases := doc.GetPhasesSorted()
	var totalCompletion float64
	var phaseCount int

	for _, phase := range phases {
		view := doc.GeneratePhaseRoadmapView(phase.ID)
		if view == nil {
			continue
		}

		period := ""
		if view.Quarter != "" {
			period = fmt.Sprintf("%s %d", view.Quarter, view.Year)
		}

		status := view.Status
		if status == "" {
			status = "planned"
		}

		goalSummary := fmt.Sprintf("%d goals", len(view.GoalViews))

		progress.PhaseProgress = append(progress.PhaseProgress, PhaseProgressSummary{
			PhaseID:     view.PhaseID,
			PhaseName:   view.PhaseName,
			Period:      period,
			Status:      status,
			GoalSummary: goalSummary,
			Completion:  view.OverallCompletion,
		})

		totalCompletion += view.OverallCompletion
		phaseCount++
	}

	if phaseCount > 0 {
		progress.OverallCompletion = totalCompletion / float64(phaseCount)
	}

	// Calculate goal progress
	for _, goal := range doc.Goals {
		currentLevel := goal.CurrentLevel
		if currentLevel == 0 {
			currentLevel = goal.CurrentMaturityLevel(doc)
		}

		status := output.GoalStatus(currentLevel, goal.TargetLevel)

		progress.GoalProgress = append(progress.GoalProgress, GoalProgressSummary{
			GoalID:       goal.ID,
			GoalName:     goal.Name,
			CurrentLevel: currentLevel,
			TargetLevel:  goal.TargetLevel,
			Status:       status,
		})
	}

	return progress
}

// GetPhaseByID returns the phase progress for a specific phase ID.
func (rp *RoadmapProgress) GetPhaseByID(phaseID string) *PhaseProgressSummary {
	for i := range rp.PhaseProgress {
		if rp.PhaseProgress[i].PhaseID == phaseID {
			return &rp.PhaseProgress[i]
		}
	}
	return nil
}

// GetGoalByID returns the goal progress for a specific goal ID.
func (rp *RoadmapProgress) GetGoalByID(goalID string) *GoalProgressSummary {
	for i := range rp.GoalProgress {
		if rp.GoalProgress[i].GoalID == goalID {
			return &rp.GoalProgress[i]
		}
	}
	return nil
}

// CompletedPhases returns phases that are completed.
func (rp *RoadmapProgress) CompletedPhases() []PhaseProgressSummary {
	var completed []PhaseProgressSummary
	for _, p := range rp.PhaseProgress {
		if p.Status == "completed" || p.Completion >= 100 {
			completed = append(completed, p)
		}
	}
	return completed
}

// AchievedGoals returns goals that have achieved their target.
func (rp *RoadmapProgress) AchievedGoals() []GoalProgressSummary {
	var achieved []GoalProgressSummary
	for _, g := range rp.GoalProgress {
		if g.Status == "Achieved" || g.CurrentLevel >= g.TargetLevel {
			achieved = append(achieved, g)
		}
	}
	return achieved
}
