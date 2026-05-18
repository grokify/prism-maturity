package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/oscompat/testutil"
	"github.com/grokify/prism-intelligence/analysis"
)

// getExampleFile returns the path to the goal-roadmap.json example file
func getExampleFile() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "..", "..", "examples", "goal-roadmap.json")
}

func TestGoalListCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Reset flag to default
	goalOutputFormat = "text"

	err := runGoalList(goalListCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runGoalList failed: %v", err)
	}
}

func TestGoalListCommandJSON(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Set JSON output
	goalOutputFormat = "json"
	defer func() { goalOutputFormat = "text" }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runGoalList(goalListCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runGoalList with JSON failed: %v", runErr)
	}

	// Check that output contains JSON array
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Errorf("Expected JSON array output, got: %s", output)
	}
}

func TestGoalShowCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	goalOutputFormat = "text"

	err := runGoalShow(goalShowCmd, []string{exampleFile, "goal-reliability"})
	if err != nil {
		t.Errorf("runGoalShow failed: %v", err)
	}
}

func TestGoalShowCommandNotFound(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	goalOutputFormat = "text"

	err := runGoalShow(goalShowCmd, []string{exampleFile, "nonexistent-goal"})
	if err == nil {
		t.Error("Expected error for nonexistent goal, got nil")
	}
}

func TestGoalProgressCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	goalOutputFormat = "text"

	err := runGoalProgress(goalProgressCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runGoalProgress failed: %v", err)
	}
}

func TestGoalProgressCommandWithGoalID(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	goalOutputFormat = "text"

	err := runGoalProgress(goalProgressCmd, []string{exampleFile, "goal-reliability"})
	if err != nil {
		t.Errorf("runGoalProgress with goal ID failed: %v", err)
	}
}

func TestPhaseListCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	phaseOutputFormat = "text"

	err := runPhaseList(phaseListCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runPhaseList failed: %v", err)
	}
}

func TestPhaseShowCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	phaseOutputFormat = "text"

	err := runPhaseShow(phaseShowCmd, []string{exampleFile, "phase-q1-2026"})
	if err != nil {
		t.Errorf("runPhaseShow failed: %v", err)
	}
}

func TestPhaseShowCommandNotFound(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	phaseOutputFormat = "text"

	err := runPhaseShow(phaseShowCmd, []string{exampleFile, "nonexistent-phase"})
	if err == nil {
		t.Error("Expected error for nonexistent phase, got nil")
	}
}

func TestPhaseMetricsCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	phaseOutputFormat = "text"

	err := runPhaseMetrics(phaseMetricsCmd, []string{exampleFile, "phase-q1-2026"})
	if err != nil {
		t.Errorf("runPhaseMetrics failed: %v", err)
	}
}

func TestRoadmapShowCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	roadmapOutputFormat = "text"
	roadmapByGoal = false

	err := runRoadmapShow(roadmapShowCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runRoadmapShow failed: %v", err)
	}
}

func TestRoadmapShowByGoalCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	roadmapOutputFormat = "text"
	roadmapByGoal = true
	defer func() { roadmapByGoal = false }()

	err := runRoadmapShow(roadmapShowCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runRoadmapShow by goal failed: %v", err)
	}
}

func TestRoadmapProgressCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	roadmapOutputFormat = "text"

	err := runRoadmapProgress(roadmapProgressCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runRoadmapProgress failed: %v", err)
	}
}

func TestScoreCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Reset flags
	scoreOutputFormat = "text"
	scoreDetailed = false
	scoreGoals = false
	scoreLegacy = false

	err := runScore(scoreCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runScore failed: %v", err)
	}
}

func TestScoreCommandWithGoals(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Set flags
	scoreOutputFormat = "text"
	scoreDetailed = false
	scoreGoals = true
	scoreLegacy = false
	defer func() { scoreGoals = false }()

	err := runScore(scoreCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runScore with goals failed: %v", err)
	}
}

func TestScoreCommandWithLegacy(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Set flags
	scoreOutputFormat = "text"
	scoreDetailed = false
	scoreGoals = false
	scoreLegacy = true
	defer func() { scoreLegacy = false }()

	err := runScore(scoreCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runScore with legacy failed: %v", err)
	}
}

func TestValidateCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	err := runValidate(validateCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runValidate failed: %v", err)
	}
}

func TestLoadPRISMDocument(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	doc, err := loadPRISMDocument(exampleFile)
	if err != nil {
		t.Errorf("loadPRISMDocument failed: %v", err)
	}

	if doc == nil {
		t.Error("loadPRISMDocument returned nil document")
	}

	if len(doc.Goals) == 0 {
		t.Error("Expected goals in document, got none")
	}

	if len(doc.Phases) == 0 {
		t.Error("Expected phases in document, got none")
	}
}

func TestLoadPRISMDocumentNotFound(t *testing.T) {
	_, err := loadPRISMDocument("/nonexistent/file.json")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestCalculateRoadmapProgress(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	doc, err := loadPRISMDocument(exampleFile)
	if err != nil {
		t.Fatalf("loadPRISMDocument failed: %v", err)
	}

	progress := analysis.CalculateRoadmapProgress(doc)
	if progress == nil {
		t.Fatal("analysis.CalculateRoadmapProgress returned nil")
	}

	if len(progress.PhaseProgress) == 0 {
		t.Error("Expected phase progress, got none")
	}

	if len(progress.GoalProgress) == 0 {
		t.Error("Expected goal progress, got none")
	}
}
