package export

import (
	"strings"
	"testing"

	"github.com/grokify/prism-intelligence"
)

func TestConvertToV2MOMEmpty(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metadata: &prism.Metadata{Name: "Test"},
	}

	v2mom := ConvertToV2MOM(doc)

	if v2mom == nil {
		t.Fatal("ConvertToV2MOM returned nil")
	}
	if v2mom.Metadata.Name != "Test V2MOM" {
		t.Errorf("Metadata.Name = %q, want %q", v2mom.Metadata.Name, "Test V2MOM")
	}
	if v2mom.Metadata.Source != "prism" {
		t.Errorf("Metadata.Source = %q, want %q", v2mom.Metadata.Source, "prism")
	}
	if len(v2mom.Methods) != 0 {
		t.Errorf("len(Methods) = %d, want 0", len(v2mom.Methods))
	}
}

func TestConvertToV2MOMWithGoals(t *testing.T) {
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
			{
				ID:           "goal-2",
				Name:         "Security",
				CurrentLevel: 3,
				TargetLevel:  5,
			},
		},
	}

	v2mom := ConvertToV2MOM(doc)

	// Check vision includes goal count
	if !strings.Contains(v2mom.Vision, "2 strategic goals") {
		t.Errorf("Vision = %q, want to contain '2 strategic goals'", v2mom.Vision)
	}

	// Check default values
	if len(v2mom.Values) != 3 {
		t.Errorf("len(Values) = %d, want 3", len(v2mom.Values))
	}

	// Check methods
	if len(v2mom.Methods) != 2 {
		t.Fatalf("len(Methods) = %d, want 2", len(v2mom.Methods))
	}

	m1 := v2mom.Methods[0]
	if m1.ID != "method-goal-1" {
		t.Errorf("Methods[0].ID = %q, want %q", m1.ID, "method-goal-1")
	}
	if m1.Name != "Reliability" {
		t.Errorf("Methods[0].Name = %q, want %q", m1.Name, "Reliability")
	}
	if m1.Owner != "SRE Team" {
		t.Errorf("Methods[0].Owner = %q, want %q", m1.Owner, "SRE Team")
	}
	if m1.Priority != "P1" {
		t.Errorf("Methods[0].Priority = %q, want %q", m1.Priority, "P1")
	}

	m2 := v2mom.Methods[1]
	if m2.Priority != "P2" {
		t.Errorf("Methods[1].Priority = %q, want %q", m2.Priority, "P2")
	}
}

func TestConvertToV2MOMWithPhases(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metadata: &prism.Metadata{Name: "Test"},
		Phases: []prism.Phase{
			{ID: "p1", Quarter: "Q2", Year: 2024},
		},
	}

	v2mom := ConvertToV2MOM(doc)

	if v2mom.Metadata.FiscalYear != 2024 {
		t.Errorf("Metadata.FiscalYear = %d, want 2024", v2mom.Metadata.FiscalYear)
	}
	if v2mom.Metadata.Quarter != "Q2" {
		t.Errorf("Metadata.Quarter = %q, want %q", v2mom.Metadata.Quarter, "Q2")
	}
}

func TestGoalToMethodWithSLOs(t *testing.T) {
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

	method := GoalToMethod(doc.Goals[0], doc, 1)

	if len(method.Measures) != 1 {
		t.Fatalf("len(Measures) = %d, want 1", len(method.Measures))
	}

	measure := method.Measures[0]
	if !strings.Contains(measure.Name, "Availability") {
		t.Errorf("Measures[0].Name = %q, want to contain Availability", measure.Name)
	}
	if measure.Target != ">=99.95%" {
		t.Errorf("Measures[0].Target = %q, want %q", measure.Target, ">=99.95%")
	}
	if measure.Timeline != "30d" {
		t.Errorf("Measures[0].Timeline = %q, want %q", measure.Timeline, "30d")
	}
}

func TestGoalToMethodWithPhases(t *testing.T) {
	doc := &prism.PRISMDocument{
		Goals: []prism.Goal{
			{
				ID:           "goal-1",
				Name:         "Reliability",
				CurrentLevel: 2,
				TargetLevel:  3,
			},
		},
		Phases: []prism.Phase{
			{
				ID:        "p1",
				StartDate: "2024-01-01",
				EndDate:   "2024-03-31",
				GoalTargets: []prism.PhaseGoalTarget{
					{GoalID: "goal-1", EnterLevel: 2, ExitLevel: 3},
				},
			},
		},
	}

	method := GoalToMethod(doc.Goals[0], doc, 1)

	if method.StartDate != "2024-01-01" {
		t.Errorf("StartDate = %q, want %q", method.StartDate, "2024-01-01")
	}
	if method.EndDate != "2024-03-31" {
		t.Errorf("EndDate = %q, want %q", method.EndDate, "2024-03-31")
	}
}

func TestGoalToMethodWithInitiatives(t *testing.T) {
	doc := &prism.PRISMDocument{
		Goals: []prism.Goal{
			{
				ID:           "goal-1",
				Name:         "Reliability",
				CurrentLevel: 2,
				TargetLevel:  3,
			},
		},
		Initiatives: []prism.Initiative{
			{ID: "init-1", GoalIDs: []string{"goal-1"}},
			{ID: "init-2", GoalIDs: []string{"goal-1", "goal-2"}},
			{ID: "init-3", GoalIDs: []string{"goal-2"}}, // different goal
		},
	}

	method := GoalToMethod(doc.Goals[0], doc, 1)

	if len(method.Projects) != 2 {
		t.Errorf("len(Projects) = %d, want 2", len(method.Projects))
	}
}

func TestSLOToMeasure(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metrics: []prism.Metric{
			{
				ID:       "m1",
				Name:     "Availability",
				Current:  99.95,
				Baseline: 99.0,
				Target:   99.95,
				SLO: &prism.SLO{
					Target: ">=99.95%",
					Window: "30d",
				},
			},
		},
	}

	slo := prism.SLORequirement{MetricID: "m1"}
	measure := SLOToMeasure(slo, 3, doc, 1)

	if measure.ID != "measure-1" {
		t.Errorf("ID = %q, want %q", measure.ID, "measure-1")
	}
	if !strings.Contains(measure.Name, "[M3]") {
		t.Errorf("Name = %q, want to contain [M3]", measure.Name)
	}
	if measure.Status != "achieved" {
		t.Errorf("Status = %q, want %q", measure.Status, "achieved")
	}
	if measure.Progress != 1.0 {
		t.Errorf("Progress = %f, want 1.0", measure.Progress)
	}
}

func TestSLOToMeasureNotMet(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metrics: []prism.Metric{
			{
				ID:       "m1",
				Name:     "Availability",
				Current:  99.5,
				Baseline: 99.0,
				Target:   99.95,
				SLO: &prism.SLO{
					Target:   ">=99.95%",
					Operator: "gte",
					Value:    99.95,
				},
			},
		},
	}

	slo := prism.SLORequirement{MetricID: "m1"}
	measure := SLOToMeasure(slo, 3, doc, 1)

	if measure.Status != "in_progress" {
		t.Errorf("Status = %q, want %q", measure.Status, "in_progress")
	}
	// Progress should be partial: (99.5 - 99.0) / (99.95 - 99.0) ≈ 0.526
	if measure.Progress < 0.5 || measure.Progress > 0.6 {
		t.Errorf("Progress = %f, want ~0.53", measure.Progress)
	}
}

func TestSLOToMeasureMetricNotFound(t *testing.T) {
	doc := &prism.PRISMDocument{}

	slo := prism.SLORequirement{MetricID: "nonexistent"}
	measure := SLOToMeasure(slo, 3, doc, 1)

	// Should still create measure with metric ID
	if !strings.Contains(measure.Name, "nonexistent") {
		t.Errorf("Name = %q, want to contain metric ID", measure.Name)
	}
}

func TestCriterionToMeasure(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metrics: []prism.Metric{
			{
				ID:       "m1",
				Name:     "Latency",
				Current:  150,
				Baseline: 500,
				Target:   200,
				Unit:     "ms",
			},
		},
	}

	mc := prism.MetricCriterion{
		MetricID: "m1",
		Operator: "lte",
		Value:    200,
	}
	measure := CriterionToMeasure(mc, 3, doc, 1)

	if !strings.Contains(measure.Name, "Latency") {
		t.Errorf("Name = %q, want to contain Latency", measure.Name)
	}
	if !strings.Contains(measure.Target, "<=") {
		t.Errorf("Target = %q, want to contain <=", measure.Target)
	}
	if measure.Status != "achieved" {
		t.Errorf("Status = %q, want %q", measure.Status, "achieved")
	}
}

func TestDefaultValues(t *testing.T) {
	doc := &prism.PRISMDocument{
		Metadata: &prism.Metadata{Name: "Test"},
		Goals: []prism.Goal{
			{ID: "g1", Name: "Goal 1"},
		},
	}

	v2mom := ConvertToV2MOM(doc)

	// Check default values are set
	if len(v2mom.Values) != 3 {
		t.Fatalf("len(Values) = %d, want 3", len(v2mom.Values))
	}

	expectedValues := []string{"Reliability First", "Continuous Improvement", "Data-Driven Decisions"}
	for i, expected := range expectedValues {
		if v2mom.Values[i].Name != expected {
			t.Errorf("Values[%d].Name = %q, want %q", i, v2mom.Values[i].Name, expected)
		}
		if v2mom.Values[i].Priority != i+1 {
			t.Errorf("Values[%d].Priority = %d, want %d", i, v2mom.Values[i].Priority, i+1)
		}
	}
}
