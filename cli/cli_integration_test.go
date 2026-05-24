package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/oscompat/testutil"
)

// getTeamTopologyFile returns the path to the team-topology.json example file
func getTeamTopologyFile() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "..", "examples", "prism-documents", "team-topology.json")
}

// getOperationsLayersFile returns the path to the operations-layers.json example file
func getOperationsLayersFile() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "..", "examples", "prism-documents", "operations-layers.json")
}

// TestLayerListCommand tests the layer list command
func TestLayerListCommand(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Reset flags
	layerJSONOutput = false

	err := runLayerList(layerListCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runLayerList failed: %v", err)
	}
}

// TestLayerListCommandJSON tests the layer list command with JSON output
func TestLayerListCommandJSON(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	layerJSONOutput = true
	defer func() { layerJSONOutput = false }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runLayerList(layerListCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runLayerList with JSON failed: %v", runErr)
	}

	// Check that output contains JSON array
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Errorf("Expected JSON array output, got: %s", output)
	}
}

// TestLayerShowCommand tests the layer show command
func TestLayerShowCommand(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	layerJSONOutput = false

	err := runLayerShow(layerShowCmd, []string{exampleFile, "runtime"})
	if err != nil {
		t.Errorf("runLayerShow failed: %v", err)
	}
}

// TestLayerShowCommandNotFound tests the layer show command with a non-existent layer
func TestLayerShowCommandNotFound(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	layerJSONOutput = false

	err := runLayerShow(layerShowCmd, []string{exampleFile, "nonexistent-layer"})
	if err == nil {
		t.Error("Expected error for nonexistent layer, got nil")
	}
}

// TestLayerShowWithSignals tests the layer show command with a layer that has golden signals
func TestLayerShowWithSignals(t *testing.T) {
	exampleFile := getOperationsLayersFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	layerJSONOutput = false

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runLayerShow(layerShowCmd, []string{exampleFile, "runtime"})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runLayerShow failed: %v", runErr)
	}

	// Check that output contains golden signals
	if !strings.Contains(output, "Golden Signals") {
		t.Errorf("Expected Golden Signals in output, got: %s", output)
	}
}

// TestTeamListCommand tests the team list command
func TestTeamListCommand(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	teamJSONOutput = false

	err := runTeamList(teamListCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runTeamList failed: %v", err)
	}
}

// TestTeamListCommandJSON tests the team list command with JSON output
func TestTeamListCommandJSON(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	teamJSONOutput = true
	defer func() { teamJSONOutput = false }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runTeamList(teamListCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runTeamList with JSON failed: %v", runErr)
	}

	// Check that output contains JSON array
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Errorf("Expected JSON array output, got: %s", output)
	}
}

// TestTeamShowCommand tests the team show command
func TestTeamShowCommand(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	teamJSONOutput = false

	err := runTeamShow(teamShowCmd, []string{exampleFile, "platform-team"})
	if err != nil {
		t.Errorf("runTeamShow failed: %v", err)
	}
}

// TestTeamShowCommandNotFound tests the team show command with a non-existent team
func TestTeamShowCommandNotFound(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	teamJSONOutput = false

	err := runTeamShow(teamShowCmd, []string{exampleFile, "nonexistent-team"})
	if err == nil {
		t.Error("Expected error for nonexistent team, got nil")
	}
}

// TestServiceListCommand tests the service list command
func TestServiceListCommand(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	serviceJSONOutput = false

	err := runServiceList(serviceListCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runServiceList failed: %v", err)
	}
}

// TestServiceListCommandJSON tests the service list command with JSON output
func TestServiceListCommandJSON(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	serviceJSONOutput = true
	defer func() { serviceJSONOutput = false }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runServiceList(serviceListCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runServiceList with JSON failed: %v", runErr)
	}

	// Check that output contains JSON array
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Errorf("Expected JSON array output, got: %s", output)
	}
}

// TestServiceShowCommand tests the service show command
func TestServiceShowCommand(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	serviceJSONOutput = false

	err := runServiceShow(serviceShowCmd, []string{exampleFile, "payments-api"})
	if err != nil {
		t.Errorf("runServiceShow failed: %v", err)
	}
}

// TestServiceShowCommandNotFound tests the service show command with a non-existent service
func TestServiceShowCommandNotFound(t *testing.T) {
	exampleFile := getTeamTopologyFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	serviceJSONOutput = false

	err := runServiceShow(serviceShowCmd, []string{exampleFile, "nonexistent-service"})
	if err == nil {
		t.Error("Expected error for nonexistent service, got nil")
	}
}

// TestAnalyzeCommand tests the analyze command
func TestAnalyzeCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Reset flags
	analyzeOutputFormat = "text"

	err := runAnalyze(analyzeCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runAnalyze failed: %v", err)
	}
}

// TestAnalyzeCommandJSON tests the analyze command with JSON output
func TestAnalyzeCommandJSON(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	analyzeOutputFormat = "json"
	defer func() { analyzeOutputFormat = "text" }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runAnalyze(analyzeCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runAnalyze with JSON failed: %v", runErr)
	}

	// Check that output contains JSON structure
	if !strings.Contains(output, "summary") {
		t.Errorf("Expected JSON with summary in output, got: %s", output)
	}
}

// TestAnalyzeCommandPrompt tests the analyze command with prompt output format
func TestAnalyzeCommandPrompt(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	analyzeOutputFormat = "prompt"
	defer func() { analyzeOutputFormat = "text" }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runAnalyze(analyzeCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runAnalyze with prompt failed: %v", runErr)
	}

	// Check that output contains prompt structure
	if !strings.Contains(output, "# PRISM Analysis Prompt") {
		t.Errorf("Expected prompt header in output, got: %s", output)
	}
}

// TestExportOKRCommand tests the export okr command
func TestExportOKRCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Reset flags - output to stdout
	exportOutputDir = ""

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runExportOKR(exportOKRCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runExportOKR failed: %v", runErr)
	}

	// Check that output contains OKR structure
	if !strings.Contains(output, "objectives") {
		t.Errorf("Expected OKR with objectives in output, got: %s", output)
	}
	if !strings.Contains(output, "keyResults") {
		t.Errorf("Expected OKR with keyResults in output, got: %s", output)
	}
}

// TestExportV2MOMCommand tests the export v2mom command
func TestExportV2MOMCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Reset flags - output to stdout
	exportOutputDir = ""

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runExportV2MOM(exportV2MOMCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runExportV2MOM failed: %v", runErr)
	}

	// Check that output contains V2MOM structure
	if !strings.Contains(output, "vision") {
		t.Errorf("Expected V2MOM with vision in output, got: %s", output)
	}
	if !strings.Contains(output, "methods") {
		t.Errorf("Expected V2MOM with methods in output, got: %s", output)
	}
}

// TestExportOKRToFile tests the export okr command writing to a file
func TestExportOKRToFile(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Create temp dir for output
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-okr.json")
	exportOutputDir = outputPath
	defer func() { exportOutputDir = "" }()

	err := runExportOKR(exportOKRCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runExportOKR to file failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to be created at %s", outputPath)
	}

	// Read and verify content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "objectives") {
		t.Error("Expected OKR with objectives in output file")
	}
}

// TestExportV2MOMToFile tests the export v2mom command writing to a file
func TestExportV2MOMToFile(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Create temp dir for output
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-v2mom.json")
	exportOutputDir = outputPath
	defer func() { exportOutputDir = "" }()

	err := runExportV2MOM(exportV2MOMCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runExportV2MOM to file failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to be created at %s", outputPath)
	}

	// Read and verify content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "methods") {
		t.Error("Expected V2MOM with methods in output file")
	}
}

// TestExportOKRToDirectory tests the export okr command writing to a directory
func TestExportOKRToDirectory(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Create temp dir for output
	tmpDir := t.TempDir()
	exportOutputDir = tmpDir
	defer func() { exportOutputDir = "" }()

	err := runExportOKR(exportOKRCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runExportOKR to directory failed: %v", err)
	}

	// Verify file was created with default name
	outputPath := filepath.Join(tmpDir, "roadmap.okr.json")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to be created at %s", outputPath)
	}
}

// TestInitiativeListCommand tests the initiative list command
func TestInitiativeListCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Reset flags
	initiativeOutputFormat = "text"
	initiativeByPhase = false
	initiativeByGoal = false

	err := runInitiativeList(initiativeListCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runInitiativeList failed: %v", err)
	}
}

// TestInitiativeListCommandByPhase tests the initiative list command grouped by phase
func TestInitiativeListCommandByPhase(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	initiativeOutputFormat = "text"
	initiativeByPhase = true
	initiativeByGoal = false
	defer func() { initiativeByPhase = false }()

	err := runInitiativeList(initiativeListCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runInitiativeList by phase failed: %v", err)
	}
}

// TestInitiativeListCommandByGoal tests the initiative list command grouped by goal
func TestInitiativeListCommandByGoal(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	initiativeOutputFormat = "text"
	initiativeByPhase = false
	initiativeByGoal = true
	defer func() { initiativeByGoal = false }()

	err := runInitiativeList(initiativeListCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runInitiativeList by goal failed: %v", err)
	}
}

// TestInitiativeListCommandJSON tests the initiative list command with JSON output
func TestInitiativeListCommandJSON(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	initiativeOutputFormat = "json"
	initiativeByPhase = false
	initiativeByGoal = false
	defer func() { initiativeOutputFormat = "text" }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runInitiativeList(initiativeListCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runInitiativeList with JSON failed: %v", runErr)
	}

	// Check that output contains JSON array
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Errorf("Expected JSON array output, got: %s", output)
	}
}

// TestInitiativeShowCommand tests the initiative show command
func TestInitiativeShowCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	initiativeOutputFormat = "text"

	err := runInitiativeShow(initiativeShowCmd, []string{exampleFile, "init-monitoring"})
	if err != nil {
		t.Errorf("runInitiativeShow failed: %v", err)
	}
}

// TestInitiativeShowCommandJSON tests the initiative show command with JSON output
func TestInitiativeShowCommandJSON(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	initiativeOutputFormat = "json"
	defer func() { initiativeOutputFormat = "text" }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runInitiativeShow(initiativeShowCmd, []string{exampleFile, "init-ci-cd"})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runInitiativeShow with JSON failed: %v", runErr)
	}

	// Check that output contains initiative fields
	if !strings.Contains(output, "init-ci-cd") {
		t.Errorf("Expected initiative ID in output, got: %s", output)
	}
}

// TestInitiativeShowCommandNotFound tests the initiative show command with non-existent initiative
func TestInitiativeShowCommandNotFound(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	initiativeOutputFormat = "text"

	err := runInitiativeShow(initiativeShowCmd, []string{exampleFile, "nonexistent-initiative"})
	if err == nil {
		t.Error("Expected error for nonexistent initiative, got nil")
	}
}

// TestInitCommand tests the init command
func TestInitCommand(t *testing.T) {
	// Create temp dir for output
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-prism.json")

	// Reset flags
	initDomain = ""
	initOutput = outputPath

	err := runInit(initCmd, []string{})
	if err != nil {
		t.Errorf("runInit failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to be created at %s", outputPath)
	}

	// Read and verify content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "schema") {
		t.Error("Expected schema field in output file")
	}
	if !strings.Contains(string(content), "metadata") {
		t.Error("Expected metadata field in output file")
	}
}

// TestInitCommandWithOperationsDomain tests the init command with operations domain
func TestInitCommandWithOperationsDomain(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "ops-prism.json")

	initDomain = "operations"
	initOutput = outputPath
	defer func() { initDomain = "" }()

	err := runInit(initCmd, []string{})
	if err != nil {
		t.Errorf("runInit with operations domain failed: %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "operations") {
		t.Error("Expected operations domain in output file")
	}
}

// TestInitCommandWithSecurityDomain tests the init command with security domain
func TestInitCommandWithSecurityDomain(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "sec-prism.json")

	initDomain = "security"
	initOutput = outputPath
	defer func() { initDomain = "" }()

	err := runInit(initCmd, []string{})
	if err != nil {
		t.Errorf("runInit with security domain failed: %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "security") {
		t.Error("Expected security domain in output file")
	}
}

// TestInitCommandWithInvalidDomain tests the init command with invalid domain
func TestInitCommandWithInvalidDomain(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "invalid-prism.json")

	initDomain = "invalid-domain"
	initOutput = outputPath
	defer func() { initDomain = "" }()

	err := runInit(initCmd, []string{})
	if err == nil {
		t.Error("Expected error for invalid domain, got nil")
	}
}

// TestCatalogCommand tests the catalog command
func TestCatalogCommand(t *testing.T) {
	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runCatalog(catalogCmd, []string{})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runCatalog failed: %v", runErr)
	}

	// Verify output contains expected sections
	if !strings.Contains(output, "Domains:") {
		t.Error("Expected Domains section in catalog output")
	}
	if !strings.Contains(output, "Lifecycle Stages:") {
		t.Error("Expected Lifecycle Stages section in catalog output")
	}
	if !strings.Contains(output, "Categories:") {
		t.Error("Expected Categories section in catalog output")
	}
	if !strings.Contains(output, "Maturity Levels:") {
		t.Error("Expected Maturity Levels section in catalog output")
	}
}

// TestPlanReportCommand tests the maturity plan report command
func TestPlanReportCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Reset flags
	planFormat = "markdown"
	planOutputFile = ""
	planView = "both"

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runMaturityPlanReport(maturityPlanReportCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runMaturityPlanReport failed: %v", runErr)
	}

	// Verify markdown output
	if !strings.Contains(output, "#") {
		t.Error("Expected markdown headers in output")
	}
}

// TestPlanReportCommandJSON tests the maturity plan report command with JSON output
func TestPlanReportCommandJSON(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	planFormat = "json"
	planOutputFile = ""
	planView = "both"
	defer func() { planFormat = "markdown" }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runMaturityPlanReport(maturityPlanReportCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runMaturityPlanReport with JSON failed: %v", runErr)
	}

	// Verify JSON output
	if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Error("Expected JSON object in output")
	}
}

// TestPlanReportCommandToFile tests the maturity plan report command writing to a file
func TestPlanReportCommandToFile(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "report.md")

	planFormat = "markdown"
	planOutputFile = outputPath
	planView = "both"
	defer func() { planOutputFile = "" }()

	err := runMaturityPlanReport(maturityPlanReportCmd, []string{exampleFile})
	if err != nil {
		t.Errorf("runMaturityPlanReport to file failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to be created at %s", outputPath)
	}
}

// TestPlanDashboardCommand tests the maturity plan dashboard command
func TestPlanDashboardCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Reset flags
	planFormat = "json"
	planOutputFile = ""

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runMaturityPlanDashboard(maturityPlanDashboardCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runMaturityPlanDashboard failed: %v", runErr)
	}

	// Verify JSON output
	if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Error("Expected JSON object in output")
	}
}

// TestPlanDashboardCommandMarkdown tests the maturity plan dashboard command with markdown output
func TestPlanDashboardCommandMarkdown(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	planFormat = "markdown"
	planOutputFile = ""
	defer func() { planFormat = "json" }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runMaturityPlanDashboard(maturityPlanDashboardCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runMaturityPlanDashboard markdown failed: %v", runErr)
	}

	// Verify markdown output
	if !strings.Contains(output, "#") {
		t.Error("Expected markdown headers in output")
	}
}

// TestSLOReportCommand tests the slo-report command
func TestSLOReportCommand(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	// Reset flags
	sloReportFormat = "json"

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runSLOReport(sloReportCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runSLOReport failed: %v", runErr)
	}

	// Verify JSON output
	if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Error("Expected JSON object in output")
	}
}

// TestSLOReportCommandMarkdown tests the slo-report command with markdown output
func TestSLOReportCommandMarkdown(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	sloReportFormat = "markdown"
	defer func() { sloReportFormat = "json" }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runSLOReport(sloReportCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runSLOReport markdown failed: %v", runErr)
	}

	// Verify markdown output
	if !strings.Contains(output, "#") {
		t.Error("Expected markdown headers in output")
	}
}

// TestSLOReportCommandMatrix tests the slo-report command with matrix output
func TestSLOReportCommandMatrix(t *testing.T) {
	exampleFile := getExampleFile()
	if _, err := os.Stat(exampleFile); os.IsNotExist(err) {
		t.Skip("example file not found, skipping integration test")
	}

	sloReportFormat = "matrix"
	defer func() { sloReportFormat = "json" }()

	var runErr error
	output, captureErr := testutil.CaptureStdout(func() {
		runErr = runSLOReport(sloReportCmd, []string{exampleFile})
	})
	if captureErr != nil {
		t.Fatalf("CaptureStdout failed: %v", captureErr)
	}

	if runErr != nil {
		t.Errorf("runSLOReport matrix failed: %v", runErr)
	}

	// Verify matrix output (should have table-like structure)
	if len(output) == 0 {
		t.Error("Expected non-empty matrix output")
	}
}
