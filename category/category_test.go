package category

import "testing"

func TestGoalCategoryDescription(t *testing.T) {
	tests := []struct {
		category GoalCategory
		wantLen  int // minimum description length
	}{
		{GoalCategoryOperations, 20},
		{GoalCategoryRevenue, 20},
		{GoalCategoryAdoption, 20},
		{GoalCategoryGrowth, 20},
		{GoalCategoryQuality, 20},
		{GoalCategorySecurity, 20},
		{GoalCategoryCompliance, 20},
	}

	for _, tt := range tests {
		desc := tt.category.Description()
		if len(desc) < tt.wantLen {
			t.Errorf("GoalCategory(%s).Description() = %q, want len >= %d", tt.category, desc, tt.wantLen)
		}
	}
}

func TestValidGoalCategories(t *testing.T) {
	categories := ValidGoalCategories()
	if len(categories) != 7 {
		t.Errorf("ValidGoalCategories() returned %d categories, want 7", len(categories))
	}
}

func TestIsValidGoalCategory(t *testing.T) {
	if !IsValidGoalCategory(GoalCategoryRevenue) {
		t.Error("IsValidGoalCategory(revenue) = false, want true")
	}
	if IsValidGoalCategory("invalid") {
		t.Error("IsValidGoalCategory(invalid) = true, want false")
	}
}

func TestParseGoalCategory(t *testing.T) {
	tests := []struct {
		input   string
		want    GoalCategory
		wantErr bool
	}{
		{"operations", GoalCategoryOperations, false},
		{"Operations", GoalCategoryOperations, false},
		{"ops", GoalCategoryOperations, false},
		{"revenue", GoalCategoryRevenue, false},
		{"rev", GoalCategoryRevenue, false},
		{"adoption", GoalCategoryAdoption, false},
		{"growth", GoalCategoryGrowth, false},
		{"quality", GoalCategoryQuality, false},
		{"security", GoalCategorySecurity, false},
		{"sec", GoalCategorySecurity, false},
		{"compliance", GoalCategoryCompliance, false},
		{"invalid", "", true},
	}

	for _, tt := range tests {
		got, err := ParseGoalCategory(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseGoalCategory(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseGoalCategory(%q) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestRevenueTargetMethods(t *testing.T) {
	target := &RevenueTarget{
		TargetARR:  1000000, // $10,000
		CurrentARR: 750000,  // $7,500
		Currency:   "USD",
	}

	if target.TargetARRInDollars() != 10000.0 {
		t.Errorf("TargetARRInDollars() = %f, want 10000.0", target.TargetARRInDollars())
	}

	if target.CurrentARRInDollars() != 7500.0 {
		t.Errorf("CurrentARRInDollars() = %f, want 7500.0", target.CurrentARRInDollars())
	}

	if target.GapARR() != 250000 {
		t.Errorf("GapARR() = %d, want 250000", target.GapARR())
	}

	if target.ProgressPercent() != 75.0 {
		t.Errorf("ProgressPercent() = %f, want 75.0", target.ProgressPercent())
	}

	if target.IsMet() {
		t.Error("IsMet() = true, want false")
	}

	// Test met
	target.CurrentARR = 1000000
	if !target.IsMet() {
		t.Error("IsMet() = false, want true (when equal)")
	}
}

func TestAdoptionTargetMethods(t *testing.T) {
	target := &AdoptionTarget{
		TargetMAU:        100000,
		CurrentMAU:       80000,
		TargetDAU:        20000,
		CurrentDAU:       16000,
		TargetRetention:  90.0,
		CurrentRetention: 85.0,
		TargetNPS:        50,
		CurrentNPS:       40,
	}

	if target.MAUProgressPercent() != 80.0 {
		t.Errorf("MAUProgressPercent() = %f, want 80.0", target.MAUProgressPercent())
	}

	if target.DAUProgressPercent() != 80.0 {
		t.Errorf("DAUProgressPercent() = %f, want 80.0", target.DAUProgressPercent())
	}

	if target.RetentionGap() != 5.0 {
		t.Errorf("RetentionGap() = %f, want 5.0", target.RetentionGap())
	}

	if target.NPSGap() != 10 {
		t.Errorf("NPSGap() = %d, want 10", target.NPSGap())
	}

	if target.MAUIsMet() {
		t.Error("MAUIsMet() = true, want false")
	}

	if target.DAUIsMet() {
		t.Error("DAUIsMet() = true, want false")
	}

	if target.RetentionIsMet() {
		t.Error("RetentionIsMet() = true, want false")
	}

	if target.NPSIsMet() {
		t.Error("NPSIsMet() = true, want false")
	}

	// DAU/MAU ratio (stickiness)
	ratio := target.DAUMAURatio()
	if ratio != 0.2 {
		t.Errorf("DAUMAURatio() = %f, want 0.2", ratio)
	}
}

func TestCategoryGoalValidate(t *testing.T) {
	tests := []struct {
		name    string
		goal    CategoryGoal
		wantErr bool
	}{
		{
			name: "valid operations",
			goal: CategoryGoal{
				GoalID:   "goal-1",
				Category: GoalCategoryOperations,
			},
			wantErr: false,
		},
		{
			name: "valid revenue with target",
			goal: CategoryGoal{
				GoalID:   "goal-2",
				Category: GoalCategoryRevenue,
				RevenueTarget: &RevenueTarget{
					TargetARR: 1000000,
					Currency:  "USD",
				},
			},
			wantErr: false,
		},
		{
			name: "missing goal_id",
			goal: CategoryGoal{
				Category: GoalCategoryRevenue,
			},
			wantErr: true,
		},
		{
			name: "missing category",
			goal: CategoryGoal{
				GoalID: "goal-1",
			},
			wantErr: true,
		},
		{
			name: "invalid category",
			goal: CategoryGoal{
				GoalID:   "goal-1",
				Category: "invalid",
			},
			wantErr: true,
		},
		{
			name: "mismatched target - revenue on adoption",
			goal: CategoryGoal{
				GoalID:   "goal-1",
				Category: GoalCategoryAdoption,
				RevenueTarget: &RevenueTarget{
					TargetARR: 1000000,
				},
			},
			wantErr: true,
		},
		{
			name: "mismatched target - adoption on revenue",
			goal: CategoryGoal{
				GoalID:   "goal-1",
				Category: GoalCategoryRevenue,
				AdoptionTarget: &AdoptionTarget{
					TargetMAU: 100000,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.goal.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRevenueTarget(t *testing.T) {
	target := NewRevenueTarget(1000000, 800000, "USD")

	if target.TargetARR != 1000000 {
		t.Errorf("TargetARR = %d, want 1000000", target.TargetARR)
	}
	if target.CurrentARR != 800000 {
		t.Errorf("CurrentARR = %d, want 800000", target.CurrentARR)
	}
	if target.Currency != "USD" {
		t.Errorf("Currency = %s, want USD", target.Currency)
	}
	if target.GrowthPercent != 25.0 {
		t.Errorf("GrowthPercent = %f, want 25.0", target.GrowthPercent)
	}
}

func TestAdoptionTargetZeroValues(t *testing.T) {
	target := &AdoptionTarget{}

	if target.MAUProgressPercent() != 0 {
		t.Errorf("MAUProgressPercent() = %f, want 0 (zero target)", target.MAUProgressPercent())
	}

	if target.DAUProgressPercent() != 0 {
		t.Errorf("DAUProgressPercent() = %f, want 0 (zero target)", target.DAUProgressPercent())
	}

	if target.DAUMAURatio() != 0 {
		t.Errorf("DAUMAURatio() = %f, want 0 (zero MAU)", target.DAUMAURatio())
	}
}
