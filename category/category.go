// Package category provides goal category types for different business domains.
//
// Goals can be categorized by their focus area:
//   - Operations: SLI/SLO targets, reliability, performance
//   - Revenue: ARR/MRR targets, growth metrics
//   - Adoption: MAU/DAU, retention, NPS targets
//   - Growth: Market expansion, new segments
//   - Quality: Reliability, performance, defect rates
//   - Security: Compliance, risk reduction
package category

import "fmt"

// GoalCategory categorizes goals by business domain.
type GoalCategory string

const (
	GoalCategoryOperations GoalCategory = "operations" // SLI/SLO, reliability, performance
	GoalCategoryRevenue    GoalCategory = "revenue"    // ARR, MRR, financial growth
	GoalCategoryAdoption   GoalCategory = "adoption"   // MAU, DAU, retention, NPS
	GoalCategoryGrowth     GoalCategory = "growth"     // Market expansion, new segments
	GoalCategoryQuality    GoalCategory = "quality"    // Reliability, performance, defects
	GoalCategorySecurity   GoalCategory = "security"   // Compliance, risk reduction
	GoalCategoryCompliance GoalCategory = "compliance" // Regulatory, audit requirements
)

// String returns the string representation.
func (c GoalCategory) String() string {
	return string(c)
}

// Description returns a human-readable description of the category.
func (c GoalCategory) Description() string {
	switch c {
	case GoalCategoryOperations:
		return "Operational excellence - SLI/SLO targets, reliability, performance"
	case GoalCategoryRevenue:
		return "Revenue growth - ARR/MRR targets, financial metrics"
	case GoalCategoryAdoption:
		return "User adoption - MAU/DAU, retention, NPS, engagement"
	case GoalCategoryGrowth:
		return "Business growth - market expansion, new segments"
	case GoalCategoryQuality:
		return "Quality improvement - reliability, performance, defect rates"
	case GoalCategorySecurity:
		return "Security posture - risk reduction, vulnerability management"
	case GoalCategoryCompliance:
		return "Compliance - regulatory requirements, audit readiness"
	default:
		return "Unknown category"
	}
}

// ValidGoalCategories returns all valid goal categories.
func ValidGoalCategories() []GoalCategory {
	return []GoalCategory{
		GoalCategoryOperations,
		GoalCategoryRevenue,
		GoalCategoryAdoption,
		GoalCategoryGrowth,
		GoalCategoryQuality,
		GoalCategorySecurity,
		GoalCategoryCompliance,
	}
}

// IsValidGoalCategory returns true if the category is valid.
func IsValidGoalCategory(c GoalCategory) bool {
	for _, valid := range ValidGoalCategories() {
		if c == valid {
			return true
		}
	}
	return false
}

// ParseGoalCategory parses a string to GoalCategory.
func ParseGoalCategory(s string) (GoalCategory, error) {
	switch s {
	case "operations", "Operations", "OPERATIONS", "ops":
		return GoalCategoryOperations, nil
	case "revenue", "Revenue", "REVENUE", "rev":
		return GoalCategoryRevenue, nil
	case "adoption", "Adoption", "ADOPTION":
		return GoalCategoryAdoption, nil
	case "growth", "Growth", "GROWTH":
		return GoalCategoryGrowth, nil
	case "quality", "Quality", "QUALITY":
		return GoalCategoryQuality, nil
	case "security", "Security", "SECURITY", "sec":
		return GoalCategorySecurity, nil
	case "compliance", "Compliance", "COMPLIANCE":
		return GoalCategoryCompliance, nil
	default:
		return "", fmt.Errorf("invalid goal category: %s", s)
	}
}

// RevenueTarget defines targets for revenue-category goals.
type RevenueTarget struct {
	TargetARR     int64   `json:"target_arr"`       // in cents
	CurrentARR    int64   `json:"current_arr"`      // in cents
	GrowthPercent float64 `json:"growth_percent"`   // target growth percentage
	Currency      string  `json:"currency"`         // USD, EUR, etc.
	Period        string  `json:"period,omitempty"` // e.g., "Q3 2026", "FY2026"
	Notes         string  `json:"notes,omitempty"`
}

// TargetARRInDollars returns the target ARR in dollars.
func (r *RevenueTarget) TargetARRInDollars() float64 {
	return float64(r.TargetARR) / 100
}

// CurrentARRInDollars returns the current ARR in dollars.
func (r *RevenueTarget) CurrentARRInDollars() float64 {
	return float64(r.CurrentARR) / 100
}

// GapARR returns the gap between target and current ARR in cents.
func (r *RevenueTarget) GapARR() int64 {
	return r.TargetARR - r.CurrentARR
}

// ProgressPercent returns progress toward target as a percentage.
func (r *RevenueTarget) ProgressPercent() float64 {
	if r.TargetARR == 0 {
		return 0
	}
	return float64(r.CurrentARR) / float64(r.TargetARR) * 100
}

// IsMet returns true if the target is met.
func (r *RevenueTarget) IsMet() bool {
	return r.CurrentARR >= r.TargetARR
}

// AdoptionTarget defines targets for adoption-category goals.
type AdoptionTarget struct {
	// User metrics
	TargetMAU  int64 `json:"target_mau,omitempty"`
	CurrentMAU int64 `json:"current_mau,omitempty"`
	TargetDAU  int64 `json:"target_dau,omitempty"`
	CurrentDAU int64 `json:"current_dau,omitempty"`

	// Engagement metrics
	TargetRetention  float64 `json:"target_retention,omitempty"`  // percentage
	CurrentRetention float64 `json:"current_retention,omitempty"` // percentage

	// Satisfaction metrics
	TargetNPS  int `json:"target_nps,omitempty"`  // -100 to 100
	CurrentNPS int `json:"current_nps,omitempty"` // -100 to 100

	// Time period
	Period string `json:"period,omitempty"` // e.g., "Q3 2026"
	Notes  string `json:"notes,omitempty"`
}

// MAUProgressPercent returns MAU progress as a percentage.
func (a *AdoptionTarget) MAUProgressPercent() float64 {
	if a.TargetMAU == 0 {
		return 0
	}
	return float64(a.CurrentMAU) / float64(a.TargetMAU) * 100
}

// DAUProgressPercent returns DAU progress as a percentage.
func (a *AdoptionTarget) DAUProgressPercent() float64 {
	if a.TargetDAU == 0 {
		return 0
	}
	return float64(a.CurrentDAU) / float64(a.TargetDAU) * 100
}

// RetentionGap returns the gap between target and current retention.
func (a *AdoptionTarget) RetentionGap() float64 {
	return a.TargetRetention - a.CurrentRetention
}

// NPSGap returns the gap between target and current NPS.
func (a *AdoptionTarget) NPSGap() int {
	return a.TargetNPS - a.CurrentNPS
}

// MAUIsMet returns true if MAU target is met.
func (a *AdoptionTarget) MAUIsMet() bool {
	return a.CurrentMAU >= a.TargetMAU
}

// DAUIsMet returns true if DAU target is met.
func (a *AdoptionTarget) DAUIsMet() bool {
	return a.CurrentDAU >= a.TargetDAU
}

// RetentionIsMet returns true if retention target is met.
func (a *AdoptionTarget) RetentionIsMet() bool {
	return a.CurrentRetention >= a.TargetRetention
}

// NPSIsMet returns true if NPS target is met.
func (a *AdoptionTarget) NPSIsMet() bool {
	return a.CurrentNPS >= a.TargetNPS
}

// DAUMAURatio returns the DAU/MAU stickiness ratio.
func (a *AdoptionTarget) DAUMAURatio() float64 {
	if a.CurrentMAU == 0 {
		return 0
	}
	return float64(a.CurrentDAU) / float64(a.CurrentMAU)
}

// CategoryGoal extends a goal with category-specific information.
type CategoryGoal struct {
	// Base goal ID reference
	GoalID string `json:"goal_id"`

	// Category
	Category GoalCategory `json:"category"`

	// Category-specific targets (only one should be populated based on category)
	RevenueTarget  *RevenueTarget  `json:"revenue_target,omitempty"`
	AdoptionTarget *AdoptionTarget `json:"adoption_target,omitempty"`

	// Notes
	Notes string `json:"notes,omitempty"`
}

// Validate returns an error if the category goal is invalid.
func (cg *CategoryGoal) Validate() error {
	if cg.GoalID == "" {
		return fmt.Errorf("goal_id is required")
	}
	if cg.Category == "" {
		return fmt.Errorf("category is required")
	}
	if !IsValidGoalCategory(cg.Category) {
		return fmt.Errorf("invalid category: %s", cg.Category)
	}

	// Validate that category-specific target matches category
	if cg.RevenueTarget != nil && cg.Category != GoalCategoryRevenue {
		return fmt.Errorf("revenue_target set but category is %s", cg.Category)
	}
	if cg.AdoptionTarget != nil && cg.Category != GoalCategoryAdoption {
		return fmt.Errorf("adoption_target set but category is %s", cg.Category)
	}

	return nil
}

// NewRevenueTarget creates a new revenue target.
func NewRevenueTarget(targetARRCents, currentARRCents int64, currency string) *RevenueTarget {
	target := &RevenueTarget{
		TargetARR:  targetARRCents,
		CurrentARR: currentARRCents,
		Currency:   currency,
	}
	if currentARRCents > 0 {
		target.GrowthPercent = float64(targetARRCents-currentARRCents) / float64(currentARRCents) * 100
	}
	return target
}

// NewAdoptionTarget creates a new adoption target.
func NewAdoptionTarget() *AdoptionTarget {
	return &AdoptionTarget{}
}
