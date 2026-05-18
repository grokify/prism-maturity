package export

import (
	"testing"

	"github.com/grokify/prism-execution/roadmap"
	"github.com/grokify/prism-intelligence"
)

func TestConvertToRoadmap(t *testing.T) {
	doc := createTestPRISMDocument()

	rm := ConvertToRoadmap(doc)
	if rm == nil {
		t.Fatal("ConvertToRoadmap returned nil")
	}

	if len(rm.Phases) != 2 {
		t.Errorf("Expected 2 phases, got %d", len(rm.Phases))
	}

	// Check first phase
	p1 := rm.Phases[0]
	if p1.ID != "phase-q1-2026" {
		t.Errorf("Expected phase ID 'phase-q1-2026', got %q", p1.ID)
	}
	if p1.Type != roadmap.PhaseTypeQuarter {
		t.Errorf("Expected quarter phase type, got %v", p1.Type)
	}
	if len(p1.Goals) != 1 {
		t.Errorf("Expected 1 goal, got %d", len(p1.Goals))
	}
	if len(p1.Deliverables) != 1 {
		t.Errorf("Expected 1 deliverable, got %d", len(p1.Deliverables))
	}
}

func TestConvertToRoadmapWithOKRs(t *testing.T) {
	doc := createTestPRISMDocument()

	export := ConvertToRoadmapWithOKRs(doc)
	if export == nil {
		t.Fatal("ConvertToRoadmapWithOKRs returned nil")
	}

	if export.Roadmap == nil {
		t.Error("Expected roadmap, got nil")
	}
	if export.OKRs == nil {
		t.Error("Expected OKRs, got nil")
	}

	// Should have objectives for each maturity level to achieve
	if len(export.OKRs.Objectives) == 0 {
		t.Error("Expected at least one objective")
	}
}

func TestConvertInitiativeWithRollout(t *testing.T) {
	doc := createTestPRISMDocument()

	rm := ConvertToRoadmap(doc)

	// Find the initiative with deployment status
	var found *roadmap.Deliverable
	for _, phase := range rm.Phases {
		for i := range phase.Deliverables {
			if phase.Deliverables[i].ID == "init-monitoring" {
				found = &phase.Deliverables[i]
				break
			}
		}
	}

	if found == nil {
		t.Fatal("Could not find init-monitoring deliverable")
	}

	if found.Rollout == nil {
		t.Fatal("Expected rollout status, got nil")
	}

	if found.Rollout.TotalCustomers != 50 {
		t.Errorf("Expected 50 total customers, got %d", found.Rollout.TotalCustomers)
	}

	if found.Rollout.DeployedCustomers != 45 {
		t.Errorf("Expected 45 deployed customers, got %d", found.Rollout.DeployedCustomers)
	}

	if found.Rollout.DeploymentPercent() != 90 {
		t.Errorf("Expected 90%% deployment, got %.0f%%", found.Rollout.DeploymentPercent())
	}
}

func TestConvertToStructuredOKR(t *testing.T) {
	doc := createTestPRISMDocument()

	okrDoc := ConvertToStructuredOKR(doc)
	if okrDoc == nil {
		t.Fatal("ConvertToStructuredOKR returned nil")
	}

	if okrDoc.Metadata == nil {
		t.Error("Expected metadata, got nil")
	}

	if len(okrDoc.Objectives) == 0 {
		t.Error("Expected objectives")
	}

	// Goal is M3→M5, so should have objectives for M4 and M5
	expectedObjectives := 2 // M4 and M5
	if len(okrDoc.Objectives) != expectedObjectives {
		t.Errorf("Expected %d objectives, got %d", expectedObjectives, len(okrDoc.Objectives))
	}

	// Check first objective (M4)
	obj := okrDoc.Objectives[0]
	if obj.Category != "Operational Maturity" {
		t.Errorf("Expected category 'Operational Maturity', got %q", obj.Category)
	}

	// Should have key results from SLOs
	if len(obj.KeyResults) == 0 {
		t.Error("Expected key results from SLOs")
	}
}

func TestConvertPhaseStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected roadmap.PhaseStatus
	}{
		{"completed", roadmap.PhaseStatusCompleted},
		{"in_progress", roadmap.PhaseStatusInProgress},
		{"delayed", roadmap.PhaseStatusDelayed},
		{"cancelled", roadmap.PhaseStatusCancelled},
		{"planned", roadmap.PhaseStatusPlanned},
		{"", roadmap.PhaseStatusPlanned},
		{"unknown", roadmap.PhaseStatusPlanned},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := convertPhaseStatus(tt.input)
			if result != tt.expected {
				t.Errorf("convertPhaseStatus(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertDeliverableStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected roadmap.DeliverableStatus
	}{
		{prism.InitiativeStatusCompleted, roadmap.DeliverableCompleted},
		{prism.InitiativeStatusInProgress, roadmap.DeliverableInProgress},
		{prism.InitiativeStatusNotStarted, roadmap.DeliverableNotStarted},
		{prism.InitiativeStatusPlanned, roadmap.DeliverableNotStarted},
		{"", roadmap.DeliverableNotStarted},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := convertDeliverableStatus(tt.input)
			if result != tt.expected {
				t.Errorf("convertDeliverableStatus(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertRolloutStage(t *testing.T) {
	tests := []struct {
		status        string
		devCompletion float64
		expected      roadmap.RolloutStage
	}{
		{"deployed", 100, roadmap.RolloutStageDeployed},
		{"completed", 100, roadmap.RolloutStageDeployed},
		{"rolling_out", 50, roadmap.RolloutStageRollingOut},
		{"in_progress", 50, roadmap.RolloutStageRollingOut},
		{"paused", 50, roadmap.RolloutStagePaused},
		{"rolled_back", 50, roadmap.RolloutStageRolledBack},
		{"", 100, roadmap.RolloutStageDeployed},
		{"", 50, roadmap.RolloutStageRollingOut},
		{"", 0, roadmap.RolloutStageNotStarted},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := convertRolloutStage(tt.status, tt.devCompletion)
			if result != tt.expected {
				t.Errorf("convertRolloutStage(%q, %.0f) = %v, want %v",
					tt.status, tt.devCompletion, result, tt.expected)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Test-Document", "test-document"},
		{"My_Project_Name", "my-project-name"},
		{"UPPERCASE", "uppercase"},
		{"with  spaces", "with-spaces"},
		{"special@#$chars", "specialchars"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := slugify(tt.input)
			if result != tt.expected {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaturityLevelName(t *testing.T) {
	tests := []struct {
		level    int
		expected string
	}{
		{1, "Reactive"},
		{2, "Basic"},
		{3, "Defined"},
		{4, "Managed"},
		{5, "Optimizing"},
		{0, "Level 0"},
		{6, "Level 6"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := maturityLevelName(tt.level)
			if result != tt.expected {
				t.Errorf("maturityLevelName(%d) = %q, want %q", tt.level, result, tt.expected)
			}
		})
	}
}

func TestRoadmapCalculateCriterionProgress(t *testing.T) {
	tests := []struct {
		name     string
		mc       prism.MetricCriterion
		metric   *prism.Metric
		expected float64
	}{
		{
			"gte 50% progress",
			prism.MetricCriterion{Operator: "gte", Value: 100},
			&prism.Metric{Baseline: 0, Current: 50},
			0.5,
		},
		{
			"gte achieved",
			prism.MetricCriterion{Operator: "gte", Value: 100},
			&prism.Metric{Baseline: 0, Current: 100},
			1.0,
		},
		{
			"lte 50% progress",
			prism.MetricCriterion{Operator: "lte", Value: 0},
			&prism.Metric{Baseline: 100, Current: 50},
			0.5,
		},
		{
			"nil metric",
			prism.MetricCriterion{Operator: "gte", Value: 100},
			nil,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateCriterionProgress(tt.mc, tt.metric)
			if result != tt.expected {
				t.Errorf("calculateCriterionProgress() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// createTestPRISMDocument creates a test PRISM document for testing.
func createTestPRISMDocument() *prism.PRISMDocument {
	return &prism.PRISMDocument{
		Metadata: &prism.Metadata{
			Name:        "Test Project",
			Description: "Test PRISM document",
			Author:      "Test Author",
		},
		Goals: []prism.Goal{
			{
				ID:           "goal-reliability",
				Name:         "Achieve High Reliability",
				Description:  "Improve system reliability",
				Owner:        "SRE Team",
				CurrentLevel: 3,
				TargetLevel:  5,
				MaturityModel: &prism.GoalMaturityModel{
					Levels: []prism.GoalMaturityLevel{
						{
							Level: 4,
							Name:  "Managed",
							RequiredSLOs: []prism.SLORequirement{
								{MetricID: "availability"},
							},
							MetricCriteria: []prism.MetricCriterion{
								{MetricID: "availability", Operator: "gte", Value: 99.9},
							},
						},
						{
							Level: 5,
							Name:  "Optimizing",
							RequiredSLOs: []prism.SLORequirement{
								{MetricID: "availability"},
							},
							MetricCriteria: []prism.MetricCriterion{
								{MetricID: "availability", Operator: "gte", Value: 99.99},
							},
						},
					},
				},
			},
		},
		Phases: []prism.Phase{
			{
				ID:        "phase-q1-2026",
				Name:      "Q1 2026",
				Quarter:   "Q1",
				Year:      2026,
				StartDate: "2026-01-01",
				EndDate:   "2026-03-31",
				Status:    "in_progress",
				GoalTargets: []prism.PhaseGoalTarget{
					{GoalID: "goal-reliability", EnterLevel: 3, ExitLevel: 4},
				},
			},
			{
				ID:        "phase-q2-2026",
				Name:      "Q2 2026",
				Quarter:   "Q2",
				Year:      2026,
				StartDate: "2026-04-01",
				EndDate:   "2026-06-30",
				Status:    "planned",
				GoalTargets: []prism.PhaseGoalTarget{
					{GoalID: "goal-reliability", EnterLevel: 4, ExitLevel: 5},
				},
			},
		},
		Initiatives: []prism.Initiative{
			{
				ID:                   "init-monitoring",
				Name:                 "Observability Platform",
				Description:          "Deploy comprehensive monitoring",
				Status:               prism.InitiativeStatusInProgress,
				PhaseID:              "phase-q1-2026",
				GoalIDs:              []string{"goal-reliability"},
				DevCompletionPercent: 90,
				DeploymentStatus: &prism.DeploymentStatus{
					Status:            "rolling_out",
					TotalCustomers:    50,
					DeployedCustomers: 45,
					AdoptionPercent:   80,
				},
			},
		},
		Metrics: []prism.Metric{
			{
				ID:       "availability",
				Name:     "Service Availability",
				Domain:   prism.DomainOperations,
				Stage:    prism.StageRuntime,
				Unit:     "%",
				Baseline: 99.0,
				Current:  99.5,
				Target:   99.99,
				SLO: &prism.SLO{
					Target: ">=99.9%",
					Window: prism.Window30Days,
				},
			},
		},
	}
}
