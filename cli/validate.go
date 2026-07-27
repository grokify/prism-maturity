package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grokify/prism-maturity"
	"github.com/grokify/prism-maturity/schema"
	"github.com/spf13/cobra"
)

var (
	validateSchema     bool
	validateSchemaType string
)

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a PRISM document",
	Long: `Validate a PRISM document against the schema and check for errors.

Use --schema to also validate against the JSON Schema.
Use --type to specify which schema type (default: maturity-plan).

Schema types:
  maturity-plan   - Full PRISM plan document (default)
  maturity-state  - Current state snapshot
  maturity-model  - Maturity model definition

Examples:
  prism validate prism.json
  prism validate --schema prism.json
  prism validate --schema --type maturity-state state.json`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().BoolVar(&validateSchema, "schema", false, "Also validate against JSON Schema")
	validateCmd.Flags().StringVar(&validateSchemaType, "type", "maturity-plan", "Schema type: maturity-plan, maturity-state, maturity-model")
}

func runValidate(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// JSON Schema validation (if requested)
	if validateSchema {
		schemaType := schema.SchemaType(validateSchemaType)
		result, err := schema.ValidateJSON(schemaType, data)
		if err != nil {
			return fmt.Errorf("schema validation error: %w", err)
		}
		if !result.Valid {
			fmt.Println("JSON Schema validation errors:")
			for _, e := range result.Errors {
				if e.Path != "" {
					fmt.Printf("  - %s: %s\n", e.Path, e.Message)
				} else {
					fmt.Printf("  - %s\n", e.Message)
				}
			}
			return fmt.Errorf("document failed JSON Schema validation with %d errors", len(result.Errors))
		}
		fmt.Printf("✓ JSON Schema validation passed (%s)\n", validateSchemaType)
	}

	// Parse JSON
	var doc prism.PRISMDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate document (structural validation)
	errs := doc.Validate()
	if errs.HasErrors() {
		fmt.Println("Structural validation errors:")
		for _, e := range errs {
			fmt.Printf("  - %s\n", e.Error())
		}
		return fmt.Errorf("document has %d validation errors", len(errs))
	}

	// Print summary
	fmt.Printf("✓ %s is valid\n", filename)
	fmt.Printf("  Metrics: %d\n", len(doc.Metrics))

	// Count by domain
	domainCounts := make(map[string]int)
	for _, m := range doc.Metrics {
		domainCounts[m.Domain]++
	}
	for domain, count := range domainCounts {
		fmt.Printf("    %s: %d\n", domain, count)
	}

	// Count by stage
	stageCounts := make(map[string]int)
	for _, m := range doc.Metrics {
		stageCounts[m.Stage]++
	}
	fmt.Println("  By stage:")
	for _, stage := range prism.AllStages() {
		if count, ok := stageCounts[stage]; ok {
			fmt.Printf("    %s: %d\n", stage, count)
		}
	}

	if doc.Maturity != nil {
		fmt.Printf("  Maturity cells: %d\n", len(doc.Maturity.Cells))
	}

	return nil
}
