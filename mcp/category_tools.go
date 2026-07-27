package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/prism-maturity/category"
)

// Category Tool Input/Output types

// ListCategoriesInput is the input for list_categories.
type ListCategoriesInput struct{}

// CategoryInfo describes a goal category.
type CategoryInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ListCategoriesOutput is the output for list_categories.
type ListCategoriesOutput struct {
	Categories []CategoryInfo `json:"categories"`
	Total      int            `json:"total"`
}

// ValidateCategoryInput is the input for validate_category.
type ValidateCategoryInput struct {
	Category string `json:"category" jsonschema:"description=Category name to validate"`
}

// ValidateCategoryOutput is the output for validate_category.
type ValidateCategoryOutput struct {
	Valid       bool   `json:"valid"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Error       string `json:"error,omitempty"`
}

// CreateRevenueTargetInput is the input for create_revenue_target.
type CreateRevenueTargetInput struct {
	TargetARR  int64  `json:"target_arr" jsonschema:"description=Target ARR in cents"`
	CurrentARR int64  `json:"current_arr" jsonschema:"description=Current ARR in cents"`
	Currency   string `json:"currency,omitempty" jsonschema:"description=Currency code (default USD)"`
	Period     string `json:"period,omitempty" jsonschema:"description=Target period (e.g., Q3 2026)"`
}

// RevenueTargetOutput is the output for create_revenue_target.
type RevenueTargetOutput struct {
	Target          *category.RevenueTarget `json:"target"`
	ProgressPercent float64                 `json:"progress_percent"`
	GapARRCents     int64                   `json:"gap_arr_cents"`
	GapARRDollars   float64                 `json:"gap_arr_dollars"`
	IsMet           bool                    `json:"is_met"`
}

// CreateAdoptionTargetInput is the input for create_adoption_target.
type CreateAdoptionTargetInput struct {
	TargetMAU        int64   `json:"target_mau,omitempty" jsonschema:"description=Target Monthly Active Users"`
	CurrentMAU       int64   `json:"current_mau,omitempty" jsonschema:"description=Current Monthly Active Users"`
	TargetDAU        int64   `json:"target_dau,omitempty" jsonschema:"description=Target Daily Active Users"`
	CurrentDAU       int64   `json:"current_dau,omitempty" jsonschema:"description=Current Daily Active Users"`
	TargetRetention  float64 `json:"target_retention,omitempty" jsonschema:"description=Target retention percentage"`
	CurrentRetention float64 `json:"current_retention,omitempty" jsonschema:"description=Current retention percentage"`
	TargetNPS        int     `json:"target_nps,omitempty" jsonschema:"description=Target NPS score (-100 to 100)"`
	CurrentNPS       int     `json:"current_nps,omitempty" jsonschema:"description=Current NPS score (-100 to 100)"`
	Period           string  `json:"period,omitempty" jsonschema:"description=Target period (e.g., Q3 2026)"`
}

// AdoptionTargetOutput is the output for create_adoption_target.
type AdoptionTargetOutput struct {
	Target         *category.AdoptionTarget `json:"target"`
	MAUProgress    float64                  `json:"mau_progress"`
	DAUProgress    float64                  `json:"dau_progress"`
	RetentionGap   float64                  `json:"retention_gap"`
	NPSGap         int                      `json:"nps_gap"`
	DAUMAURatio    float64                  `json:"dau_mau_ratio"`
	MAUIsMet       bool                     `json:"mau_is_met"`
	DAUIsMet       bool                     `json:"dau_is_met"`
	RetentionIsMet bool                     `json:"retention_is_met"`
	NPSIsMet       bool                     `json:"nps_is_met"`
}

// RegisterCategoryTools registers category-related MCP tools.
func (s *Server) RegisterCategoryTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_categories",
		Description: "List all valid goal categories with descriptions.",
	}, s.listCategories)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "validate_category",
		Description: "Validate a goal category string.",
	}, s.validateCategory)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_revenue_target",
		Description: "Create a revenue target with ARR values. Returns progress and gap calculations.",
	}, s.createRevenueTarget)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_adoption_target",
		Description: "Create an adoption target with MAU/DAU/retention/NPS values. Returns progress calculations.",
	}, s.createAdoptionTarget)
}

// Tool implementations

func (s *Server) listCategories(ctx context.Context, req *mcp.CallToolRequest, input ListCategoriesInput) (*mcp.CallToolResult, ListCategoriesOutput, error) {
	categories := category.ValidGoalCategories()

	var infos []CategoryInfo
	for _, c := range categories {
		infos = append(infos, CategoryInfo{
			Name:        string(c),
			Description: c.Description(),
		})
	}

	return nil, ListCategoriesOutput{
		Categories: infos,
		Total:      len(infos),
	}, nil
}

func (s *Server) validateCategory(ctx context.Context, req *mcp.CallToolRequest, input ValidateCategoryInput) (*mcp.CallToolResult, ValidateCategoryOutput, error) {
	cat, err := category.ParseGoalCategory(input.Category)
	if err != nil {
		return nil, ValidateCategoryOutput{
			Valid: false,
			Error: fmt.Sprintf("invalid category: %s", input.Category),
		}, nil
	}

	return nil, ValidateCategoryOutput{
		Valid:       true,
		Category:    string(cat),
		Description: cat.Description(),
	}, nil
}

func (s *Server) createRevenueTarget(ctx context.Context, req *mcp.CallToolRequest, input CreateRevenueTargetInput) (*mcp.CallToolResult, RevenueTargetOutput, error) {
	currency := input.Currency
	if currency == "" {
		currency = "USD"
	}

	target := category.NewRevenueTarget(input.TargetARR, input.CurrentARR, currency)
	target.Period = input.Period

	return nil, RevenueTargetOutput{
		Target:          target,
		ProgressPercent: target.ProgressPercent(),
		GapARRCents:     target.GapARR(),
		GapARRDollars:   float64(target.GapARR()) / 100,
		IsMet:           target.IsMet(),
	}, nil
}

func (s *Server) createAdoptionTarget(ctx context.Context, req *mcp.CallToolRequest, input CreateAdoptionTargetInput) (*mcp.CallToolResult, AdoptionTargetOutput, error) {
	target := category.NewAdoptionTarget()
	target.TargetMAU = input.TargetMAU
	target.CurrentMAU = input.CurrentMAU
	target.TargetDAU = input.TargetDAU
	target.CurrentDAU = input.CurrentDAU
	target.TargetRetention = input.TargetRetention
	target.CurrentRetention = input.CurrentRetention
	target.TargetNPS = input.TargetNPS
	target.CurrentNPS = input.CurrentNPS
	target.Period = input.Period

	return nil, AdoptionTargetOutput{
		Target:         target,
		MAUProgress:    target.MAUProgressPercent(),
		DAUProgress:    target.DAUProgressPercent(),
		RetentionGap:   target.RetentionGap(),
		NPSGap:         target.NPSGap(),
		DAUMAURatio:    target.DAUMAURatio(),
		MAUIsMet:       target.MAUIsMet(),
		DAUIsMet:       target.DAUIsMet(),
		RetentionIsMet: target.RetentionIsMet(),
		NPSIsMet:       target.NPSIsMet(),
	}, nil
}
