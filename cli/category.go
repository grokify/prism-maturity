package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/grokify/prism-maturity/category"
)

var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "Manage goal categories",
	Long:  `Commands for working with goal categories including revenue, adoption, operations, etc.`,
}

var categoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all goal categories",
	Long:  `Lists all valid goal categories with their descriptions.`,
	RunE:  runCategoryList,
}

var categoryTargetCmd = &cobra.Command{
	Use:   "target",
	Short: "Manage category-specific targets",
	Long:  `Create and validate category-specific targets like revenue or adoption targets.`,
}

var revenueTargetCmd = &cobra.Command{
	Use:   "revenue",
	Short: "Create a revenue target",
	Long: `Create a revenue target with ARR values.

Example:
  maturity category target revenue --target-arr 1000000 --current-arr 750000 --currency USD`,
	RunE: runRevenueTarget,
}

var adoptionTargetCmd = &cobra.Command{
	Use:   "adoption",
	Short: "Create an adoption target",
	Long: `Create an adoption target with MAU/DAU/retention/NPS values.

Example:
  maturity category target adoption --target-mau 100000 --current-mau 75000 --target-retention 85 --current-retention 78`,
	RunE: runAdoptionTarget,
}

var categoryValidateCmd = &cobra.Command{
	Use:   "validate <category>",
	Short: "Validate a goal category",
	Long:  `Validates that a category string is a valid goal category.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCategoryValidate,
}

var (
	// Revenue target flags
	targetARR  int64
	currentARR int64
	currency   string
	period     string

	// Adoption target flags
	targetMAU        int64
	currentMAU       int64
	targetDAU        int64
	currentDAU       int64
	targetRetention  float64
	currentRetention float64
	targetNPS        int
	currentNPS       int

	// Output flags
	categoryOutputJSON bool
)

func init() {
	categoryCmd.AddCommand(categoryListCmd)
	categoryCmd.AddCommand(categoryTargetCmd)
	categoryCmd.AddCommand(categoryValidateCmd)

	categoryTargetCmd.AddCommand(revenueTargetCmd)
	categoryTargetCmd.AddCommand(adoptionTargetCmd)

	// Revenue target flags
	revenueTargetCmd.Flags().Int64Var(&targetARR, "target-arr", 0, "Target ARR in cents")
	revenueTargetCmd.Flags().Int64Var(&currentARR, "current-arr", 0, "Current ARR in cents")
	revenueTargetCmd.Flags().StringVar(&currency, "currency", "USD", "Currency code (USD, EUR, etc.)")
	revenueTargetCmd.Flags().StringVar(&period, "period", "", "Target period (e.g., Q3 2026)")

	// Adoption target flags
	adoptionTargetCmd.Flags().Int64Var(&targetMAU, "target-mau", 0, "Target Monthly Active Users")
	adoptionTargetCmd.Flags().Int64Var(&currentMAU, "current-mau", 0, "Current Monthly Active Users")
	adoptionTargetCmd.Flags().Int64Var(&targetDAU, "target-dau", 0, "Target Daily Active Users")
	adoptionTargetCmd.Flags().Int64Var(&currentDAU, "current-dau", 0, "Current Daily Active Users")
	adoptionTargetCmd.Flags().Float64Var(&targetRetention, "target-retention", 0, "Target retention percentage")
	adoptionTargetCmd.Flags().Float64Var(&currentRetention, "current-retention", 0, "Current retention percentage")
	adoptionTargetCmd.Flags().IntVar(&targetNPS, "target-nps", 0, "Target NPS score (-100 to 100)")
	adoptionTargetCmd.Flags().IntVar(&currentNPS, "current-nps", 0, "Current NPS score (-100 to 100)")
	adoptionTargetCmd.Flags().StringVar(&period, "period", "", "Target period (e.g., Q3 2026)")

	// Output flags
	categoryListCmd.Flags().BoolVar(&categoryOutputJSON, "json", false, "Output as JSON")
	revenueTargetCmd.Flags().BoolVar(&categoryOutputJSON, "json", false, "Output as JSON")
	adoptionTargetCmd.Flags().BoolVar(&categoryOutputJSON, "json", false, "Output as JSON")
}

func runCategoryList(cmd *cobra.Command, args []string) error {
	categories := category.ValidGoalCategories()

	if categoryOutputJSON {
		type categoryInfo struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		var info []categoryInfo
		for _, c := range categories {
			info = append(info, categoryInfo{
				Name:        string(c),
				Description: c.Description(),
			})
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CATEGORY\tDESCRIPTION")
	fmt.Fprintln(w, "--------\t-----------")
	for _, c := range categories {
		fmt.Fprintf(w, "%s\t%s\n", c, c.Description())
	}
	return w.Flush()
}

func runCategoryValidate(cmd *cobra.Command, args []string) error {
	categoryStr := args[0]

	cat, err := category.ParseGoalCategory(categoryStr)
	if err != nil {
		fmt.Printf("Invalid category: %s\n", categoryStr)
		fmt.Println("\nValid categories:")
		for _, c := range category.ValidGoalCategories() {
			fmt.Printf("  - %s\n", c)
		}
		os.Exit(1)
	}

	fmt.Printf("Valid category: %s\n", cat)
	fmt.Printf("Description: %s\n", cat.Description())
	return nil
}

func runRevenueTarget(cmd *cobra.Command, args []string) error {
	target := category.NewRevenueTarget(targetARR, currentARR, currency)
	target.Period = period

	if categoryOutputJSON {
		result := map[string]interface{}{
			"target":              target,
			"progress_percent":    target.ProgressPercent(),
			"gap_arr_cents":       target.GapARR(),
			"gap_arr_dollars":     float64(target.GapARR()) / 100,
			"is_met":              target.IsMet(),
			"target_arr_dollars":  target.TargetARRInDollars(),
			"current_arr_dollars": target.CurrentARRInDollars(),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Println("Revenue Target Summary")
	fmt.Println("======================")
	fmt.Printf("Currency:     %s\n", target.Currency)
	if target.Period != "" {
		fmt.Printf("Period:       %s\n", target.Period)
	}
	fmt.Printf("Target ARR:   $%.2f\n", target.TargetARRInDollars())
	fmt.Printf("Current ARR:  $%.2f\n", target.CurrentARRInDollars())
	fmt.Printf("Gap:          $%.2f\n", float64(target.GapARR())/100)
	fmt.Printf("Progress:     %.1f%%\n", target.ProgressPercent())
	fmt.Printf("Growth:       %.1f%%\n", target.GrowthPercent)

	if target.IsMet() {
		fmt.Println("Status:       TARGET MET")
	} else {
		fmt.Println("Status:       In Progress")
	}

	return nil
}

func runAdoptionTarget(cmd *cobra.Command, args []string) error {
	target := category.NewAdoptionTarget()
	target.TargetMAU = targetMAU
	target.CurrentMAU = currentMAU
	target.TargetDAU = targetDAU
	target.CurrentDAU = currentDAU
	target.TargetRetention = targetRetention
	target.CurrentRetention = currentRetention
	target.TargetNPS = targetNPS
	target.CurrentNPS = currentNPS
	target.Period = period

	if categoryOutputJSON {
		result := map[string]interface{}{
			"target":           target,
			"mau_progress":     target.MAUProgressPercent(),
			"dau_progress":     target.DAUProgressPercent(),
			"retention_gap":    target.RetentionGap(),
			"nps_gap":          target.NPSGap(),
			"dau_mau_ratio":    target.DAUMAURatio(),
			"mau_is_met":       target.MAUIsMet(),
			"dau_is_met":       target.DAUIsMet(),
			"retention_is_met": target.RetentionIsMet(),
			"nps_is_met":       target.NPSIsMet(),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Println("Adoption Target Summary")
	fmt.Println("=======================")
	if target.Period != "" {
		fmt.Printf("Period:           %s\n", target.Period)
	}

	if target.TargetMAU > 0 {
		fmt.Printf("\nMAU:\n")
		fmt.Printf("  Target:         %d\n", target.TargetMAU)
		fmt.Printf("  Current:        %d\n", target.CurrentMAU)
		fmt.Printf("  Progress:       %.1f%%\n", target.MAUProgressPercent())
		if target.MAUIsMet() {
			fmt.Println("  Status:         TARGET MET")
		}
	}

	if target.TargetDAU > 0 {
		fmt.Printf("\nDAU:\n")
		fmt.Printf("  Target:         %d\n", target.TargetDAU)
		fmt.Printf("  Current:        %d\n", target.CurrentDAU)
		fmt.Printf("  Progress:       %.1f%%\n", target.DAUProgressPercent())
		if target.DAUIsMet() {
			fmt.Println("  Status:         TARGET MET")
		}
	}

	if target.CurrentMAU > 0 && target.CurrentDAU > 0 {
		fmt.Printf("\nStickiness (DAU/MAU): %.1f%%\n", target.DAUMAURatio()*100)
	}

	if target.TargetRetention > 0 {
		fmt.Printf("\nRetention:\n")
		fmt.Printf("  Target:         %.1f%%\n", target.TargetRetention)
		fmt.Printf("  Current:        %.1f%%\n", target.CurrentRetention)
		fmt.Printf("  Gap:            %.1f%%\n", target.RetentionGap())
		if target.RetentionIsMet() {
			fmt.Println("  Status:         TARGET MET")
		}
	}

	if target.TargetNPS != 0 || target.CurrentNPS != 0 {
		fmt.Printf("\nNPS:\n")
		fmt.Printf("  Target:         %d\n", target.TargetNPS)
		fmt.Printf("  Current:        %d\n", target.CurrentNPS)
		fmt.Printf("  Gap:            %d\n", target.NPSGap())
		if target.NPSIsMet() {
			fmt.Println("  Status:         TARGET MET")
		}
	}

	return nil
}
