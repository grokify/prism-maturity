package maturity

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateOrganizationXLSX(t *testing.T) {
	specFile := "../examples/organization/model.json"
	outputFile := filepath.Join(os.TempDir(), "organization-maturity.xlsx")

	// Use the simple XLSX generator (omniframe-based)
	if err := GenerateSimpleXLSX(specFile, outputFile); err != nil {
		t.Fatalf("Failed to generate simple XLSX: %v", err)
	}

	info, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	t.Logf("Generated: %s (%d bytes)", outputFile, info.Size())

	// Also generate the full-featured version for comparison
	fullOutputFile := filepath.Join(os.TempDir(), "organization-maturity-full.xlsx")
	if err := GenerateXLSX(specFile, fullOutputFile); err != nil {
		t.Fatalf("Failed to generate full XLSX: %v", err)
	}

	fullInfo, err := os.Stat(fullOutputFile)
	if err != nil {
		t.Fatalf("Failed to stat full output file: %v", err)
	}

	t.Logf("Generated (full): %s (%d bytes)", fullOutputFile, fullInfo.Size())
}
