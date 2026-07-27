# Goal Categories

The `category` package provides goal categorization with category-specific target types for revenue and adoption metrics.

## Overview

Goal categories help organize maturity goals by functional area:

- **Operations** - Reliability, efficiency, DORA metrics
- **Revenue** - ARR growth, monetization
- **Adoption** - MAU, retention, NPS
- **Growth** - Market expansion, user acquisition
- **Quality** - Testing, defect management
- **Security** - Security posture, compliance
- **Compliance** - Regulatory requirements

## GoalCategory Type

```go
type GoalCategory string

const (
    GoalCategoryOperations GoalCategory = "operations"
    GoalCategoryRevenue    GoalCategory = "revenue"
    GoalCategoryAdoption   GoalCategory = "adoption"
    GoalCategoryGrowth     GoalCategory = "growth"
    GoalCategoryQuality    GoalCategory = "quality"
    GoalCategorySecurity   GoalCategory = "security"
    GoalCategoryCompliance GoalCategory = "compliance"
)
```

## Category Methods

```go
// Check if category is valid
category := GoalCategoryRevenue
if category.IsValid() {
    // Use category
}

// Get all valid categories
allCategories := ValidGoalCategories()

// String representation
str := category.String() // "revenue"
```

## RevenueTarget Type

Revenue-specific goal targets with ARR tracking:

```go
type RevenueTarget struct {
    TargetARR     int64   `json:"target_arr"`      // Target ARR in cents
    CurrentARR    int64   `json:"current_arr"`     // Current ARR in cents
    GrowthPercent float64 `json:"growth_percent"`  // Target growth %
    Currency      string  `json:"currency"`        // ISO currency code
}
```

### RevenueTarget Methods

```go
target := &RevenueTarget{
    TargetARR:     10000000000, // $100M
    CurrentARR:    80000000000, // $80M
    GrowthPercent: 25.0,
    Currency:      "USD",
}

// Calculate progress toward target
progress := target.Progress() // 0.8 (80%)

// Get gap to target
gap := target.Gap() // $20M in cents

// Get ARR values in dollars
targetDollars := target.TargetARRInDollars()   // 100000000.0
currentDollars := target.CurrentARRInDollars() // 80000000.0
```

### JSON Example

```json
{
  "category": "revenue",
  "revenue_target": {
    "target_arr": 10000000000,
    "current_arr": 8000000000,
    "growth_percent": 25.0,
    "currency": "USD"
  }
}
```

## AdoptionTarget Type

Adoption-specific goal targets with user engagement metrics:

```go
type AdoptionTarget struct {
    TargetMAU          int64   `json:"target_mau"`
    CurrentMAU         int64   `json:"current_mau"`
    TargetRetention    float64 `json:"target_retention"`    // 0.0-1.0
    CurrentRetention   float64 `json:"current_retention"`   // 0.0-1.0
    TargetNPS          int     `json:"target_nps"`          // -100 to 100
    CurrentNPS         int     `json:"current_nps"`         // -100 to 100
}
```

### AdoptionTarget Methods

```go
target := &AdoptionTarget{
    TargetMAU:        1000000,
    CurrentMAU:       750000,
    TargetRetention:  0.85,
    CurrentRetention: 0.78,
    TargetNPS:        50,
    CurrentNPS:       42,
}

// Calculate MAU progress
mauProgress := target.MAUProgress() // 0.75 (75%)

// Calculate retention progress
retentionProgress := target.RetentionProgress() // 0.918 (91.8%)

// Calculate NPS progress (normalized to 0-1)
npsProgress := target.NPSProgress() // 0.71 (71%)

// Get overall progress (average of all metrics)
overallProgress := target.Progress() // Average of all metrics
```

### JSON Example

```json
{
  "category": "adoption",
  "adoption_target": {
    "target_mau": 1000000,
    "current_mau": 750000,
    "target_retention": 0.85,
    "current_retention": 0.78,
    "target_nps": 50,
    "current_nps": 42
  }
}
```

## Usage with Goals

Attach category-specific targets to goals:

```go
import "github.com/grokify/prism-maturity/category"

// Revenue goal with target
revenueGoal := Goal{
    ID:       "goal-revenue-q2",
    Name:     "Q2 Revenue Growth",
    Category: category.GoalCategoryRevenue,
}

revenueTarget := &category.RevenueTarget{
    TargetARR:     10000000000,
    CurrentARR:    8000000000,
    GrowthPercent: 25.0,
    Currency:      "USD",
}

// Adoption goal with target
adoptionGoal := Goal{
    ID:       "goal-adoption-q2",
    Name:     "Q2 User Adoption",
    Category: category.GoalCategoryAdoption,
}

adoptionTarget := &category.AdoptionTarget{
    TargetMAU:        1000000,
    CurrentMAU:       750000,
    TargetRetention:  0.85,
    CurrentRetention: 0.78,
}
```

## Category Descriptions

| Category | Focus Area | Example Goals |
|----------|------------|---------------|
| `operations` | Reliability, efficiency | MTTR < 1hr, Deploy frequency 10/day |
| `revenue` | Monetization | ARR growth 25%, Churn < 5% |
| `adoption` | User engagement | MAU 1M, NPS > 50, Retention > 85% |
| `growth` | Market expansion | New markets, User acquisition |
| `quality` | Product quality | Test coverage 80%, Critical bugs 0 |
| `security` | Security posture | SAST coverage 100%, No critical CVEs |
| `compliance` | Regulatory | SOC2 certified, GDPR compliant |

## Integration with Maturity Model

Categories align with maturity domains:

| Category | Primary Domain |
|----------|----------------|
| `operations` | operations |
| `revenue` | product |
| `adoption` | product |
| `growth` | product |
| `quality` | quality |
| `security` | security |
| `compliance` | security |

## See Also

- [Goals](../concepts/goals.md) - Goal concepts
- [Domains](domains.md) - Domain schema
- [Metrics](metrics.md) - Metric types
