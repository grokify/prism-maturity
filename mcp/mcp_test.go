//nolint:dupl,gosec // test functions have similar structure by design; G306 ok for test fixtures
package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	prism "github.com/grokify/prism-maturity"
	"github.com/grokify/prism-maturity/category"
)

func createTestDocument(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "mcp-test-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}

	doc := &prism.PRISMDocument{
		Layers: []prism.LayerDef{
			{ID: "code", Name: "Code Quality", Weight: 0.3},
			{ID: "infra", Name: "Infrastructure", Weight: 0.4},
			{ID: "ops", Name: "Operations", Weight: 0.3},
		},
		Services: []prism.Service{
			{ID: "api", Name: "API Service", LayerID: "code", Tier: "tier1"},
			{ID: "db", Name: "Database", LayerID: "infra", Tier: "tier1"},
			{ID: "cache", Name: "Cache Service", LayerID: "infra", Tier: "tier2"},
		},
		Goals: []prism.Goal{
			{ID: "goal-1", Name: "Improve Reliability", Status: prism.GoalStatusActive, CurrentLevel: 2, TargetLevel: 4},
			{ID: "goal-2", Name: "Scale Performance", Status: prism.GoalStatusActive, CurrentLevel: 3, TargetLevel: 5},
			{ID: "goal-3", Name: "Security Hardening", Status: prism.GoalStatusOnHold, CurrentLevel: 1, TargetLevel: 3},
		},
		Metrics: []prism.Metric{
			{ID: "uptime", Name: "Uptime", Current: 99.5, Target: 99.9, SLO: &prism.SLO{Value: 99.9, Operator: prism.SLOOperatorGTE}},
			{ID: "latency", Name: "P99 Latency", Current: 150, Target: 100, SLO: &prism.SLO{Value: 100, Operator: prism.SLOOperatorLTE}},
			{ID: "errors", Name: "Error Rate", Current: 0.5, Target: 1.0, SLO: &prism.SLO{Value: 1.0, Operator: prism.SLOOperatorLTE}},
		},
	}

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("marshaling document: %v", err)
	}

	filePath := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		t.Fatalf("writing file: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return filePath
}

func TestServer_listCategories(t *testing.T) {
	srv := NewServer()

	_, output, err := srv.listCategories(context.Background(), nil, ListCategoriesInput{})
	if err != nil {
		t.Fatalf("listCategories() error = %v", err)
	}

	if output.Total != len(category.ValidGoalCategories()) {
		t.Errorf("Total = %d, want %d", output.Total, len(category.ValidGoalCategories()))
	}

	// Verify operations category exists
	found := false
	for _, cat := range output.Categories {
		if cat.Name == "operations" {
			found = true
			break
		}
	}
	if !found {
		t.Error("operations category not found")
	}
}

func TestServer_validateCategory(t *testing.T) {
	srv := NewServer()

	t.Run("valid category", func(t *testing.T) {
		_, output, err := srv.validateCategory(context.Background(), nil, ValidateCategoryInput{
			Category: "revenue",
		})
		if err != nil {
			t.Fatalf("validateCategory() error = %v", err)
		}

		if !output.Valid {
			t.Error("Valid = false, want true")
		}
		if output.Category != "revenue" {
			t.Errorf("Category = %s, want revenue", output.Category)
		}
	})

	t.Run("invalid category", func(t *testing.T) {
		_, output, err := srv.validateCategory(context.Background(), nil, ValidateCategoryInput{
			Category: "invalid",
		})
		if err != nil {
			t.Fatalf("validateCategory() error = %v", err)
		}

		if output.Valid {
			t.Error("Valid = true, want false")
		}
		if output.Error == "" {
			t.Error("expected error message")
		}
	})
}

func TestServer_createRevenueTarget(t *testing.T) {
	srv := NewServer()

	_, output, err := srv.createRevenueTarget(context.Background(), nil, CreateRevenueTargetInput{
		TargetARR:  1000000, // $10,000
		CurrentARR: 750000,  // $7,500
		Currency:   "USD",
		Period:     "Q3 2026",
	})
	if err != nil {
		t.Fatalf("createRevenueTarget() error = %v", err)
	}

	if output.Target == nil {
		t.Fatal("Target is nil")
	}
	if output.ProgressPercent != 75 {
		t.Errorf("ProgressPercent = %f, want 75", output.ProgressPercent)
	}
	if output.GapARRCents != 250000 {
		t.Errorf("GapARRCents = %d, want 250000", output.GapARRCents)
	}
	if output.IsMet {
		t.Error("IsMet = true, want false")
	}
}

func TestServer_createAdoptionTarget(t *testing.T) {
	srv := NewServer()

	_, output, err := srv.createAdoptionTarget(context.Background(), nil, CreateAdoptionTargetInput{
		TargetMAU:        100000,
		CurrentMAU:       80000,
		TargetRetention:  90,
		CurrentRetention: 85,
		TargetNPS:        50,
		CurrentNPS:       45,
	})
	if err != nil {
		t.Fatalf("createAdoptionTarget() error = %v", err)
	}

	if output.Target == nil {
		t.Fatal("Target is nil")
	}
	if output.MAUProgress != 80 {
		t.Errorf("MAUProgress = %f, want 80", output.MAUProgress)
	}
	if output.RetentionGap != 5 {
		t.Errorf("RetentionGap = %f, want 5", output.RetentionGap)
	}
	if output.NPSGap != 5 {
		t.Errorf("NPSGap = %d, want 5", output.NPSGap)
	}
}

func TestServer_listGoals(t *testing.T) {
	srv := NewServer()
	filePath := createTestDocument(t)

	t.Run("list all goals", func(t *testing.T) {
		_, output, err := srv.listGoals(context.Background(), nil, ListGoalsInput{
			File: filePath,
		})
		if err != nil {
			t.Fatalf("listGoals() error = %v", err)
		}

		if output.Error != "" {
			t.Errorf("unexpected error: %s", output.Error)
		}
		if output.Total != 3 {
			t.Errorf("Total = %d, want 3", output.Total)
		}
	})

	t.Run("filter by status", func(t *testing.T) {
		_, output, err := srv.listGoals(context.Background(), nil, ListGoalsInput{
			File:   filePath,
			Status: prism.GoalStatusActive,
		})
		if err != nil {
			t.Fatalf("listGoals() error = %v", err)
		}

		if output.Total != 2 {
			t.Errorf("Total = %d, want 2", output.Total)
		}
	})
}

func TestServer_getGoal(t *testing.T) {
	srv := NewServer()
	filePath := createTestDocument(t)

	t.Run("get existing goal", func(t *testing.T) {
		_, output, err := srv.getGoal(context.Background(), nil, GetGoalInput{
			File: filePath,
			ID:   "goal-1",
		})
		if err != nil {
			t.Fatalf("getGoal() error = %v", err)
		}

		if output.Error != "" {
			t.Errorf("unexpected error: %s", output.Error)
		}
		if output.Goal == nil {
			t.Fatal("Goal is nil")
		}
		if output.Goal.Name != "Improve Reliability" {
			t.Errorf("Name = %s, want Improve Reliability", output.Goal.Name)
		}
	})

	t.Run("get nonexistent goal", func(t *testing.T) {
		_, output, err := srv.getGoal(context.Background(), nil, GetGoalInput{
			File: filePath,
			ID:   "nonexistent",
		})
		if err != nil {
			t.Fatalf("getGoal() error = %v", err)
		}

		if output.Error == "" {
			t.Error("expected error for nonexistent goal")
		}
	})
}

func TestServer_goalStatuses(t *testing.T) {
	srv := NewServer()

	_, output, err := srv.goalStatuses(context.Background(), nil, GoalStatusesInput{})
	if err != nil {
		t.Fatalf("goalStatuses() error = %v", err)
	}

	if len(output.Statuses) != 4 {
		t.Errorf("len(Statuses) = %d, want 4", len(output.Statuses))
	}
}

func TestServer_maturitySummary(t *testing.T) {
	srv := NewServer()
	filePath := createTestDocument(t)

	_, output, err := srv.maturitySummary(context.Background(), nil, MaturitySummaryInput{
		File: filePath,
	})
	if err != nil {
		t.Fatalf("maturitySummary() error = %v", err)
	}

	if output.Error != "" {
		t.Errorf("unexpected error: %s", output.Error)
	}
	if output.TotalLayers != 3 {
		t.Errorf("TotalLayers = %d, want 3", output.TotalLayers)
	}
	if output.TotalServices != 3 {
		t.Errorf("TotalServices = %d, want 3", output.TotalServices)
	}
	if output.TotalGoals != 3 {
		t.Errorf("TotalGoals = %d, want 3", output.TotalGoals)
	}
	if output.TotalMetrics != 3 {
		t.Errorf("TotalMetrics = %d, want 3", output.TotalMetrics)
	}
}

func TestServer_listMetrics(t *testing.T) {
	srv := NewServer()
	filePath := createTestDocument(t)

	t.Run("list all metrics", func(t *testing.T) {
		_, output, err := srv.listMetrics(context.Background(), nil, ListMetricsInput{
			File: filePath,
		})
		if err != nil {
			t.Fatalf("listMetrics() error = %v", err)
		}

		if output.Total != 3 {
			t.Errorf("Total = %d, want 3", output.Total)
		}
	})

	t.Run("filter by SLO met", func(t *testing.T) {
		_, output, err := srv.listMetrics(context.Background(), nil, ListMetricsInput{
			File:      filePath,
			SLOStatus: "met",
		})
		if err != nil {
			t.Fatalf("listMetrics() error = %v", err)
		}

		// uptime: 99.5 < 99.9 (not met, gte)
		// latency: 150 > 100 (not met, lte)
		// errors: 0.5 <= 1.0 (met, lte)
		if output.Total != 1 {
			t.Errorf("Total = %d, want 1 (only errors SLO is met)", output.Total)
		}
	})
}

func TestServer_getMetric(t *testing.T) {
	srv := NewServer()
	filePath := createTestDocument(t)

	t.Run("get existing metric", func(t *testing.T) {
		_, output, err := srv.getMetric(context.Background(), nil, GetMetricInput{
			File: filePath,
			ID:   "uptime",
		})
		if err != nil {
			t.Fatalf("getMetric() error = %v", err)
		}

		if output.Error != "" {
			t.Errorf("unexpected error: %s", output.Error)
		}
		if output.Metric == nil {
			t.Fatal("Metric is nil")
		}
		if output.Metric.Current != 99.5 {
			t.Errorf("Current = %f, want 99.5", output.Metric.Current)
		}
	})

	t.Run("get nonexistent metric", func(t *testing.T) {
		_, output, err := srv.getMetric(context.Background(), nil, GetMetricInput{
			File: filePath,
			ID:   "nonexistent",
		})
		if err != nil {
			t.Fatalf("getMetric() error = %v", err)
		}

		if output.Error == "" {
			t.Error("expected error for nonexistent metric")
		}
	})
}

func TestServer_listLayers(t *testing.T) {
	srv := NewServer()
	filePath := createTestDocument(t)

	_, output, err := srv.listLayers(context.Background(), nil, ListLayersInput{
		File: filePath,
	})
	if err != nil {
		t.Fatalf("listLayers() error = %v", err)
	}

	if output.Total != 3 {
		t.Errorf("Total = %d, want 3", output.Total)
	}
}

func TestServer_listServices(t *testing.T) {
	srv := NewServer()
	filePath := createTestDocument(t)

	t.Run("list all services", func(t *testing.T) {
		_, output, err := srv.listServices(context.Background(), nil, ListServicesInput{
			File: filePath,
		})
		if err != nil {
			t.Fatalf("listServices() error = %v", err)
		}

		if output.Total != 3 {
			t.Errorf("Total = %d, want 3", output.Total)
		}
	})

	t.Run("filter by layer", func(t *testing.T) {
		_, output, err := srv.listServices(context.Background(), nil, ListServicesInput{
			File:    filePath,
			LayerID: "infra",
		})
		if err != nil {
			t.Fatalf("listServices() error = %v", err)
		}

		if output.Total != 2 {
			t.Errorf("Total = %d, want 2", output.Total)
		}
	})

	t.Run("filter by tier", func(t *testing.T) {
		_, output, err := srv.listServices(context.Background(), nil, ListServicesInput{
			File: filePath,
			Tier: "tier1",
		})
		if err != nil {
			t.Fatalf("listServices() error = %v", err)
		}

		if output.Total != 2 {
			t.Errorf("Total = %d, want 2", output.Total)
		}
	})
}
