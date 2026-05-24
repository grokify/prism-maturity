package scaffold

import (
	"testing"

	"github.com/grokify/prism-maturity"
)

func TestNewDocument(t *testing.T) {
	// Test with no domains (defaults to security + operations)
	doc := NewDocument()
	if doc == nil {
		t.Fatal("NewDocument() returned nil")
	}
	if len(doc.Domains) != 2 {
		t.Errorf("Expected 2 domains, got %d", len(doc.Domains))
	}
	if doc.Metadata == nil {
		t.Error("Expected metadata, got nil")
	}
	if doc.Maturity == nil {
		t.Error("Expected maturity model, got nil")
	}

	// Test with single domain
	doc = NewDocument(prism.DomainOperations)
	if len(doc.Domains) != 1 {
		t.Errorf("Expected 1 domain, got %d", len(doc.Domains))
	}
	if doc.Domains[0].Name != prism.DomainOperations {
		t.Errorf("Expected operations domain, got %s", doc.Domains[0].Name)
	}
	if doc.Domains[0].Weight != 1.0 {
		t.Errorf("Expected weight 1.0 for single domain, got %f", doc.Domains[0].Weight)
	}

	// Test with multiple domains
	doc = NewDocument(prism.DomainSecurity, prism.DomainOperations)
	if len(doc.Domains) != 2 {
		t.Errorf("Expected 2 domains, got %d", len(doc.Domains))
	}
	expectedWeight := 0.5
	for _, d := range doc.Domains {
		if d.Weight != expectedWeight {
			t.Errorf("Expected weight %f, got %f for domain %s", expectedWeight, d.Weight, d.Name)
		}
	}
}

func TestOperationsMetrics(t *testing.T) {
	metrics := OperationsMetrics()
	if len(metrics) == 0 {
		t.Error("OperationsMetrics() returned empty slice")
	}

	for _, m := range metrics {
		if m.ID == "" {
			t.Error("Metric ID is empty")
		}
		if m.Name == "" {
			t.Error("Metric Name is empty")
		}
		if m.Domain != prism.DomainOperations {
			t.Errorf("Expected operations domain, got %s", m.Domain)
		}
	}

	// Check that availability metric exists
	found := false
	for _, m := range metrics {
		if m.ID == "ops-availability-01" {
			found = true
			if m.SLO == nil {
				t.Error("Availability metric should have SLO")
			}
			if m.SLI == nil {
				t.Error("Availability metric should have SLI")
			}
		}
	}
	if !found {
		t.Error("Expected ops-availability-01 metric")
	}
}

func TestSecurityMetrics(t *testing.T) {
	metrics := SecurityMetrics()
	if len(metrics) == 0 {
		t.Error("SecurityMetrics() returned empty slice")
	}

	for _, m := range metrics {
		if m.ID == "" {
			t.Error("Metric ID is empty")
		}
		if m.Name == "" {
			t.Error("Metric Name is empty")
		}
		if m.Domain != prism.DomainSecurity {
			t.Errorf("Expected security domain, got %s", m.Domain)
		}
	}

	// Check that vulnerability metric exists
	found := false
	for _, m := range metrics {
		if m.ID == "sec-vuln-01" {
			found = true
			if m.TrendDirection != prism.TrendLowerBetter {
				t.Error("Vulnerability metric should have lower-is-better trend")
			}
		}
	}
	if !found {
		t.Error("Expected sec-vuln-01 metric")
	}
}
