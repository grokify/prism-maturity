package analysis

import (
	"testing"

	"github.com/grokify/prism-intelligence"
)

func TestAnalyzeEmptyDocument(t *testing.T) {
	doc := &prism.PRISMDocument{}
	result := Analyze(doc)

	if result == nil {
		t.Fatal("Analyze returned nil")
	}

	if result.Summary.TotalGoals != 0 {
		t.Errorf("TotalGoals = %d, want 0", result.Summary.TotalGoals)
	}
	if result.Summary.TotalPhases != 0 {
		t.Errorf("TotalPhases = %d, want 0", result.Summary.TotalPhases)
	}
	if result.Summary.TotalSLOs != 0 {
		t.Errorf("TotalSLOs = %d, want 0", result.Summary.TotalSLOs)
	}
}

func TestAnalyzeWithGoals(t *testing.T) {
	doc := &prism.PRISMDocument{
		Goals: []prism.Goal{
			{
				ID:           "goal-1",
				Name:         "Reliability",
				CurrentLevel: 2,
				TargetLevel:  4,
			},
			{
				ID:           "goal-2",
				Name:         "Security",
				CurrentLevel: 3,
				TargetLevel:  5,
			},
		},
	}

	result := Analyze(doc)

	if result.Summary.TotalGoals != 2 {
		t.Errorf("TotalGoals = %d, want 2", result.Summary.TotalGoals)
	}

	if len(result.Goals) != 2 {
		t.Errorf("len(Goals) = %d, want 2", len(result.Goals))
	}

	// Check first goal
	g1 := result.Goals[0]
	if g1.GoalID != "goal-1" {
		t.Errorf("Goals[0].GoalID = %q, want %q", g1.GoalID, "goal-1")
	}
	if g1.Gap != 2 {
		t.Errorf("Goals[0].Gap = %d, want 2", g1.Gap)
	}
	if g1.Status != "Behind" {
		t.Errorf("Goals[0].Status = %q, want %q", g1.Status, "Behind")
	}

	// Check average gap
	expectedAvgGap := 2.0 // (2 + 2) / 2
	if result.Summary.AvgMaturityGap != expectedAvgGap {
		t.Errorf("AvgMaturityGap = %f, want %f", result.Summary.AvgMaturityGap, expectedAvgGap)
	}
}

func TestAnalyzeWithSLOs(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metrics: []prism.Metric{
			{
				ID:      "metric-1",
				Name:    "Availability",
				Current: 99.9,
				Target:  99.5,
				SLO: &prism.SLO{
					Target:   ">=99.5%",
					Operator: "gte",
					Value:    99.5,
				},
			},
			{
				ID:      "metric-2",
				Name:    "Latency",
				Current: 250,
				Target:  200,
				SLO: &prism.SLO{
					Target:   "<=200ms",
					Operator: "lte",
					Value:    200,
				},
			},
		},
	}

	result := Analyze(doc)

	if result.Summary.TotalSLOs != 2 {
		t.Errorf("TotalSLOs = %d, want 2", result.Summary.TotalSLOs)
	}
	if result.Summary.SLOsMet != 1 {
		t.Errorf("SLOsMet = %d, want 1", result.Summary.SLOsMet)
	}
	if result.Summary.SLOCompliance != 50 {
		t.Errorf("SLOCompliance = %f, want 50", result.Summary.SLOCompliance)
	}
}

func TestAnalyzeWithPhases(t *testing.T) {
	doc := &prism.PRISMDocument{
		Goals: []prism.Goal{
			{
				ID:           "goal-1",
				Name:         "Reliability",
				CurrentLevel: 2,
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

	result := Analyze(doc)

	if result.Summary.TotalPhases != 2 {
		t.Errorf("TotalPhases = %d, want 2", result.Summary.TotalPhases)
	}

	if len(result.Phases) != 2 {
		t.Errorf("len(Phases) = %d, want 2", len(result.Phases))
	}

	// Check first phase
	p1 := result.Phases[0]
	if p1.PhaseID != "phase-q1" {
		t.Errorf("Phases[0].PhaseID = %q, want %q", p1.PhaseID, "phase-q1")
	}
	if p1.Period != "Q1 2024" {
		t.Errorf("Phases[0].Period = %q, want %q", p1.Period, "Q1 2024")
	}
	if p1.Status != "completed" {
		t.Errorf("Phases[0].Status = %q, want %q", p1.Status, "completed")
	}
	if len(p1.GoalTargets) != 1 {
		t.Errorf("len(Phases[0].GoalTargets) = %d, want 1", len(p1.GoalTargets))
	}
}

func TestCountSLOsForLevelRange(t *testing.T) {
	goal := &prism.Goal{
		ID:   "test-goal",
		Name: "Test Goal",
		MaturityModel: &prism.GoalMaturityModel{
			Levels: []prism.GoalMaturityLevel{
				{
					Level: 2,
					RequiredSLOs: []prism.SLORequirement{
						{MetricID: "m1"},
						{MetricID: "m2"},
					},
				},
				{
					Level: 3,
					RequiredSLOs: []prism.SLORequirement{
						{MetricID: "m3"},
					},
				},
				{
					Level: 4,
					RequiredSLOs: []prism.SLORequirement{
						{MetricID: "m4"},
						{MetricID: "m5"},
						{MetricID: "m6"},
					},
				},
			},
		},
	}

	tests := []struct {
		name       string
		enterLevel int
		exitLevel  int
		want       int
	}{
		{"level 1 to 2", 1, 2, 2},
		{"level 2 to 3", 2, 3, 1},
		{"level 2 to 4", 2, 4, 4}, // levels 3 and 4
		{"level 1 to 4", 1, 4, 6}, // levels 2, 3, and 4
		{"same level", 3, 3, 0},
		{"level 4 to 5", 4, 5, 0}, // no level 5 defined
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountSLOsForLevelRange(goal, tt.enterLevel, tt.exitLevel)
			if got != tt.want {
				t.Errorf("CountSLOsForLevelRange(%d, %d) = %d, want %d", tt.enterLevel, tt.exitLevel, got, tt.want)
			}
		})
	}
}

func TestCountSLOsForLevelRangeNilGoal(t *testing.T) {
	got := CountSLOsForLevelRange(nil, 1, 3)
	if got != 0 {
		t.Errorf("CountSLOsForLevelRange(nil) = %d, want 0", got)
	}
}

func TestCountSLOsForLevelRangeNoMaturityModel(t *testing.T) {
	goal := &prism.Goal{
		ID:   "test-goal",
		Name: "Test Goal",
	}
	got := CountSLOsForLevelRange(goal, 1, 3)
	if got != 0 {
		t.Errorf("CountSLOsForLevelRange(goal without maturity model) = %d, want 0", got)
	}
}
