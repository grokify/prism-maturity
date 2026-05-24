package export

import (
	"strings"
	"testing"

	"github.com/grokify/prism-maturity"
)

func TestConvertToOKREmpty(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metadata: &prism.Metadata{Name: "Test"},
	}

	okr := ConvertToOKR(doc)

	if okr == nil {
		t.Fatal("ConvertToOKR returned nil")
	}
	if okr.Metadata.Name != "Test OKRs" {
		t.Errorf("Metadata.Name = %q, want %q", okr.Metadata.Name, "Test OKRs")
	}
	if okr.Metadata.Source != "prism" {
		t.Errorf("Metadata.Source = %q, want %q", okr.Metadata.Source, "prism")
	}
	if len(okr.Objectives) != 0 {
		t.Errorf("len(Objectives) = %d, want 0", len(okr.Objectives))
	}
}

func TestConvertToOKRWithGoals(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metadata: &prism.Metadata{Name: "Test", Author: "Test Author"},
		Goals: []prism.Goal{
			{
				ID:           "goal-1",
				Name:         "Reliability",
				Description:  "Improve reliability",
				CurrentLevel: 2,
				TargetLevel:  4,
				Owner:        "SRE Team",
			},
		},
	}

	okr := ConvertToOKR(doc)

	if len(okr.Objectives) != 1 {
		t.Fatalf("len(Objectives) = %d, want 1", len(okr.Objectives))
	}

	obj := okr.Objectives[0]
	if obj.ID != "obj-goal-1" {
		t.Errorf("Objectives[0].ID = %q, want %q", obj.ID, "obj-goal-1")
	}
	if obj.Title != "Reliability" {
		t.Errorf("Objectives[0].Title = %q, want %q", obj.Title, "Reliability")
	}
	if obj.Owner != "SRE Team" {
		t.Errorf("Objectives[0].Owner = %q, want %q", obj.Owner, "SRE Team")
	}
	if !strings.Contains(obj.Rationale, "M2 to M4") {
		t.Errorf("Objectives[0].Rationale = %q, want to contain M2 to M4", obj.Rationale)
	}
}

func TestConvertToOKRWithPeriod(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metadata: &prism.Metadata{Name: "Test"},
		Phases: []prism.Phase{
			{ID: "p1", Quarter: "Q1", Year: 2024},
			{ID: "p2", Quarter: "Q4", Year: 2024},
		},
	}

	okr := ConvertToOKR(doc)

	if okr.Metadata.Period != "Q1 2024 - Q4 2024" {
		t.Errorf("Metadata.Period = %q, want %q", okr.Metadata.Period, "Q1 2024 - Q4 2024")
	}
}

func TestGoalToObjectiveWithSLOs(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metrics: []prism.Metric{
			{
				ID:       "m1",
				Name:     "Availability",
				Current:  99.9,
				Baseline: 99.0,
				Target:   99.95,
				SLO: &prism.SLO{
					Target:   ">=99.95%",
					Operator: "gte",
					Value:    99.95,
					Window:   "30d",
				},
			},
		},
		Goals: []prism.Goal{
			{
				ID:           "goal-1",
				Name:         "Reliability",
				CurrentLevel: 2,
				TargetLevel:  3,
				MaturityModel: &prism.GoalMaturityModel{
					Levels: []prism.GoalMaturityLevel{
						{
							Level: 3,
							RequiredSLOs: []prism.SLORequirement{
								{MetricID: "m1"},
							},
						},
					},
				},
			},
		},
	}

	obj := GoalToObjective(doc.Goals[0], doc)

	if len(obj.KeyResults) != 1 {
		t.Fatalf("len(KeyResults) = %d, want 1", len(obj.KeyResults))
	}

	kr := obj.KeyResults[0]
	if kr.Metric != "m1" {
		t.Errorf("KeyResults[0].Metric = %q, want %q", kr.Metric, "m1")
	}
	if !strings.Contains(kr.Title, "Availability") {
		t.Errorf("KeyResults[0].Title = %q, want to contain Availability", kr.Title)
	}
	if kr.Target != ">=99.95%" {
		t.Errorf("KeyResults[0].Target = %q, want %q", kr.Target, ">=99.95%")
	}
}

func TestCalculateCriterionProgress(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		value    float64
		baseline float64
		current  float64
		want     float64
	}{
		{"gte 50%", "gte", 100, 0, 50, 0.5},
		{"gte 100%", "gte", 100, 0, 100, 1.0},
		{"gte exceeded", "gte", 100, 0, 150, 1.0},
		{"gte negative", "gte", 100, 0, -50, 0.0},
		{"lte 50%", "lte", 0, 100, 50, 0.5},
		{"lte 100%", "lte", 0, 100, 0, 1.0},
		{"lte exceeded", "lte", 0, 100, -50, 1.0},
		{"invalid baseline gte", "gte", 50, 100, 75, 0.0}, // target <= baseline
		{"invalid baseline lte", "lte", 100, 50, 75, 0.0}, // target >= baseline
		{"unknown operator", "eq", 100, 0, 50, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := prism.MetricCriterion{
				MetricID: "test",
				Operator: tt.operator,
				Value:    tt.value,
			}
			metric := &prism.Metric{
				Baseline: tt.baseline,
				Current:  tt.current,
			}

			got := CalculateCriterionProgress(mc, metric)
			if got != tt.want {
				t.Errorf("CalculateCriterionProgress() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestCalculateCriterionProgressNilMetric(t *testing.T) {
	mc := prism.MetricCriterion{Operator: "gte", Value: 100}
	got := CalculateCriterionProgress(mc, nil)
	if got != 0 {
		t.Errorf("CalculateCriterionProgress(nil) = %f, want 0", got)
	}
}

func TestGeneratePhaseTargets(t *testing.T) {
	phases := []prism.Phase{
		{
			ID:     "p1",
			Status: "completed",
			GoalTargets: []prism.PhaseGoalTarget{
				{GoalID: "goal-1", EnterLevel: 2, ExitLevel: 3},
			},
		},
		{
			ID:     "p2",
			Status: "in_progress",
			GoalTargets: []prism.PhaseGoalTarget{
				{GoalID: "goal-1", EnterLevel: 3, ExitLevel: 4},
			},
		},
		{
			ID:     "p3",
			Status: "planned",
			GoalTargets: []prism.PhaseGoalTarget{
				{GoalID: "goal-2", EnterLevel: 1, ExitLevel: 2}, // different goal
			},
		},
	}

	targets := GeneratePhaseTargets(phases, "goal-1", 4)

	// Should only include phases where exit level >= target level (4)
	if len(targets) != 1 {
		t.Fatalf("len(targets) = %d, want 1", len(targets))
	}

	if targets[0].PhaseID != "p2" {
		t.Errorf("targets[0].PhaseID = %q, want %q", targets[0].PhaseID, "p2")
	}
	if targets[0].Status != "in_progress" {
		t.Errorf("targets[0].Status = %q, want %q", targets[0].Status, "in_progress")
	}
}

func TestGeneratePhaseTargetsEmpty(t *testing.T) {
	phases := []prism.Phase{}
	targets := GeneratePhaseTargets(phases, "goal-1", 3)

	if len(targets) != 0 {
		t.Errorf("len(targets) = %d, want 0", len(targets))
	}
}
