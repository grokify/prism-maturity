package cli

import (
	"io"
	"os"

	"github.com/grokify/prism-maturity/output"
)

// Re-export types from output package for CLI use
type (
	OutputFormat  = output.Format
	Formatter     = output.Formatter
	TableData     = output.TableData
	DetailData    = output.DetailData
	DetailField   = output.DetailField
	DetailSection = output.DetailSection
)

// Re-export constants
const (
	FormatText     = output.FormatText
	FormatJSON     = output.FormatJSON
	FormatMarkdown = output.FormatMarkdown
	FormatTOON     = output.FormatTOON
)

// ValidFormats returns all valid format strings.
func ValidFormats() []string {
	return output.ValidFormats()
}

// IsValidFormat checks if a format string is valid.
func IsValidFormat(f string) bool {
	return output.IsValidFormat(f)
}

// NewFormatter creates a formatter for the given format string.
func NewFormatter(format string) *Formatter {
	return output.NewFormatter(format)
}

// NewFormatterWithWriter creates a formatter with a custom writer.
func NewFormatterWithWriter(format string, w io.Writer) *Formatter {
	return output.NewFormatterWithWriter(format, w)
}

// Helper functions re-exported from output package
var (
	TruncateString         = output.TruncateString
	OperatorSymbol         = output.OperatorSymbol
	SafePercent            = output.SafePercent
	GoalStatus             = output.GoalStatus
	MaturityLevelName      = output.MaturityLevelName
	StatusSymbol           = output.StatusSymbol
	FormatInitiativeStatus = output.FormatInitiativeStatus
)

// getGoalStatus is the old name, kept for backwards compatibility in CLI
func getGoalStatus(current, target int) string {
	return output.GoalStatus(current, target)
}

// truncateString is the old name, kept for backwards compatibility in CLI
func truncateString(s string, maxLen int) string {
	return output.TruncateString(s, maxLen)
}

// operatorSymbol is the old name, kept for backwards compatibility in CLI
func operatorSymbol(op string) string {
	return output.OperatorSymbol(op)
}

// safePercent is the old name, kept for backwards compatibility in CLI
func safePercent(value, total int) float64 {
	return output.SafePercent(value, total)
}

// FormatToStdout creates a formatter that writes to stdout
func FormatToStdout(format string) *Formatter {
	return &Formatter{
		Format: OutputFormat(format),
		Writer: os.Stdout,
	}
}
