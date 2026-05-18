package analysis

import (
	"testing"

	"github.com/grokify/prism-intelligence"
)

func TestCalculateRoadmapProgressEmpty(t *testing.T) {
	doc := &prism.PRISMDocument{}
	progress := CalculateRoadmapProgress(doc)

	if progress == nil {
		t.Fatal("CalculateRoadmapProgress returned nil")
	}

	if progress.OverallCompletion != 0 {
		t.Errorf("OverallCompletion = %f, want 0", progress.OverallCompletion)
	}
	if len(progress.PhaseProgress) != 0 {
		t.Errorf("len(PhaseProgress) = %d, want 0", len(progress.PhaseProgress))
	}
	if len(progress.GoalProgress) != 0 {
		t.Errorf("len(GoalProgress) = %d, want 0", len(progress.GoalProgress))
	}
}

func TestCalculateRoadmapProgressWithGoals(t *testing.T) {
	doc := &prism.PRISMDocument{
		Goals: []prism.Goal{
			{
				ID:           "goal-1",
				Name:         "Reliability",
				CurrentLevel: 3,
				TargetLevel:  5,
			},
			{
				ID:           "goal-2",
				Name:         "Security",
				CurrentLevel: 5,
				TargetLevel:  5,
			},
		},
	}

	progress := CalculateRoadmapProgress(doc)

	if len(progress.GoalProgress) != 2 {
		t.Errorf("len(GoalProgress) = %d, want 2", len(progress.GoalProgress))
	}

	// Check goal statuses
	g1 := progress.GetGoalByID("goal-1")
	if g1 == nil {
		t.Fatal("GetGoalByID(goal-1) returned nil")
	}
	if g1.Status != "Behind" {
		t.Errorf("goal-1 Status = %q, want %q", g1.Status, "Behind")
	}

	g2 := progress.GetGoalByID("goal-2")
	if g2 == nil {
		t.Fatal("GetGoalByID(goal-2) returned nil")
	}
	if g2.Status != "Achieved" {
		t.Errorf("goal-2 Status = %q, want %q", g2.Status, "Achieved")
	}
}

func TestCalculateRoadmapProgressWithPhases(t *testing.T) {
	doc := &prism.PRISMDocument{
		Goals: []prism.Goal{
			{
				ID:           "goal-1",
				Name:         "Reliability",
				CurrentLevel: 3,
				TargetLevel:  4,
			},
		},
		Phases: []prism.Phase{
			{
				ID:      "phase-q1",
				Name:    "Q1 2024",
				Quarter: "Q1",
				Year:    2024,
				Status:  "completed",
				GoalTargets: []prism.PhaseGoalTarget{
					{GoalID: "goal-1", EnterLevel: 2, ExitLevel: 3},
				},
			},
			{
				ID:      "phase-q2",
				Name:    "Q2 2024",
				Quarter: "Q2",
				Year:    2024,
				Status:  "in_progress",
				GoalTargets: []prism.PhaseGoalTarget{
					{GoalID: "goal-1", EnterLevel: 3, ExitLevel: 4},
				},
			},
		},
	}

	progress := CalculateRoadmapProgress(doc)

	if len(progress.PhaseProgress) != 2 {
		t.Errorf("len(PhaseProgress) = %d, want 2", len(progress.PhaseProgress))
	}

	// Check phase lookup
	p1 := progress.GetPhaseByID("phase-q1")
	if p1 == nil {
		t.Fatal("GetPhaseByID(phase-q1) returned nil")
	}
	if p1.Period != "Q1 2024" {
		t.Errorf("phase-q1 Period = %q, want %q", p1.Period, "Q1 2024")
	}
	if p1.Status != "completed" {
		t.Errorf("phase-q1 Status = %q, want %q", p1.Status, "completed")
	}
}

func TestRoadmapProgressGetByIDNotFound(t *testing.T) {
	progress := &RoadmapProgress{
		PhaseProgress: []PhaseProgressSummary{
			{PhaseID: "phase-1"},
		},
		GoalProgress: []GoalProgressSummary{
			{GoalID: "goal-1"},
		},
	}

	if progress.GetPhaseByID("nonexistent") != nil {
		t.Error("GetPhaseByID(nonexistent) should return nil")
	}
	if progress.GetGoalByID("nonexistent") != nil {
		t.Error("GetGoalByID(nonexistent) should return nil")
	}
}

func TestRoadmapProgressCompletedPhases(t *testing.T) {
	progress := &RoadmapProgress{
		PhaseProgress: []PhaseProgressSummary{
			{PhaseID: "p1", Status: "completed", Completion: 100},
			{PhaseID: "p2", Status: "in_progress", Completion: 50},
			{PhaseID: "p3", Status: "planned", Completion: 100}, // completion=100 counts as completed
			{PhaseID: "p4", Status: "planned", Completion: 0},
		},
	}

	completed := progress.CompletedPhases()

	if len(completed) != 2 {
		t.Errorf("len(CompletedPhases) = %d, want 2", len(completed))
	}

	// Check which phases are in the list
	ids := make(map[string]bool)
	for _, p := range completed {
		ids[p.PhaseID] = true
	}
	if !ids["p1"] {
		t.Error("CompletedPhases should include p1 (status=completed)")
	}
	if !ids["p3"] {
		t.Error("CompletedPhases should include p3 (completion=100)")
	}
}

func TestRoadmapProgressAchievedGoals(t *testing.T) {
	progress := &RoadmapProgress{
		GoalProgress: []GoalProgressSummary{
			{GoalID: "g1", Status: "Achieved", CurrentLevel: 5, TargetLevel: 5},
			{GoalID: "g2", Status: "On Track", CurrentLevel: 4, TargetLevel: 5},
			{GoalID: "g3", Status: "Behind", CurrentLevel: 5, TargetLevel: 4}, // exceeded target
			{GoalID: "g4", Status: "Behind", CurrentLevel: 2, TargetLevel: 5},
		},
	}

	achieved := progress.AchievedGoals()

	if len(achieved) != 2 {
		t.Errorf("len(AchievedGoals) = %d, want 2", len(achieved))
	}

	// Check which goals are in the list
	ids := make(map[string]bool)
	for _, g := range achieved {
		ids[g.GoalID] = true
	}
	if !ids["g1"] {
		t.Error("AchievedGoals should include g1 (status=Achieved)")
	}
	if !ids["g3"] {
		t.Error("AchievedGoals should include g3 (current >= target)")
	}
}

func TestRoadmapProgressEmptyFilters(t *testing.T) {
	progress := &RoadmapProgress{
		PhaseProgress: []PhaseProgressSummary{
			{PhaseID: "p1", Status: "in_progress", Completion: 50},
		},
		GoalProgress: []GoalProgressSummary{
			{GoalID: "g1", Status: "Behind", CurrentLevel: 2, TargetLevel: 5},
		},
	}

	completed := progress.CompletedPhases()
	if len(completed) != 0 {
		t.Errorf("CompletedPhases should be empty, got %d", len(completed))
	}

	achieved := progress.AchievedGoals()
	if len(achieved) != 0 {
		t.Errorf("AchievedGoals should be empty, got %d", len(achieved))
	}
}
