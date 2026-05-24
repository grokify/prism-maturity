// Package report provides report generation for PRISM documents.
package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/grokify/prism-maturity"
)

// MarkdownOptions configures markdown report generation.
type MarkdownOptions struct {
	Title           string
	Author          string
	Date            string
	IncludeYAMLMeta bool   // Include YAML front matter for Pandoc
	ViewType        string // "both", "by-goal", "by-phase"
	IncludeDetails  bool   // Include initiative details
}

// DefaultMarkdownOptions returns sensible defaults.
func DefaultMarkdownOptions() *MarkdownOptions {
	return &MarkdownOptions{
		Title:           "PRISM Roadmap Report",
		Date:            time.Now().Format("2006-01-02"),
		IncludeYAMLMeta: true,
		ViewType:        "both",
		IncludeDetails:  true,
	}
}

// GenerateMarkdown generates a Pandoc-compatible Markdown report.
func GenerateMarkdown(doc *prism.PRISMDocument, opts *MarkdownOptions) string {
	if opts == nil {
		opts = DefaultMarkdownOptions()
	}

	var sb strings.Builder

	// YAML front matter for Pandoc
	if opts.IncludeYAMLMeta {
		sb.WriteString("---\n")
		sb.WriteString(fmt.Sprintf("title: \"%s\"\n", escapeYAML(opts.Title)))
		if opts.Author != "" {
			sb.WriteString(fmt.Sprintf("author: \"%s\"\n", escapeYAML(opts.Author)))
		} else if doc.Metadata != nil && doc.Metadata.Author != "" {
			sb.WriteString(fmt.Sprintf("author: \"%s\"\n", escapeYAML(doc.Metadata.Author)))
		}
		sb.WriteString(fmt.Sprintf("date: \"%s\"\n", opts.Date))
		sb.WriteString("---\n\n")
	}

	// Document metadata section
	if doc.Metadata != nil {
		if doc.Metadata.Description != "" {
			sb.WriteString(fmt.Sprintf("> %s\n\n", doc.Metadata.Description))
		}
	}

	report := doc.GenerateRoadmapReport()

	// Generate views based on option
	switch opts.ViewType {
	case "by-goal":
		writeGoalView(&sb, report, opts)
	case "by-phase":
		writePhaseView(&sb, report, opts)
	default: // "both"
		writePhaseView(&sb, report, opts)
		sb.WriteString("\n---\n\n")
		writeGoalView(&sb, report, opts)
	}

	return sb.String()
}

func writePhaseView(sb *strings.Builder, report *prism.RoadmapReport, opts *MarkdownOptions) {
	sb.WriteString("# Roadmap by Phase\n\n")
	sb.WriteString("This view shows progress organized by time period (Phase → Goal → Initiative).\n\n")

	for _, phase := range report.ByPhase {
		// Phase header
		phaseName := phase.PhaseName
		if phase.Quarter != "" && phase.Year > 0 {
			phaseName = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
		}
		sb.WriteString(fmt.Sprintf("## %s\n\n", phaseName))

		// Phase metadata
		sb.WriteString(fmt.Sprintf("- **Period:** %s to %s\n", phase.StartDate, phase.EndDate))
		sb.WriteString(fmt.Sprintf("- **Status:** %s\n", formatStatus(phase.Status)))
		sb.WriteString(fmt.Sprintf("- **Overall Completion:** %.1f%%\n\n", phase.OverallCompletion))

		// Progress bar
		sb.WriteString(fmt.Sprintf("%s\n\n", progressBar(phase.OverallCompletion)))

		// Goals table
		if len(phase.GoalViews) > 0 {
			sb.WriteString("### Goals\n\n")
			sb.WriteString("| Goal | Maturity | Progress | Complete |\n")
			sb.WriteString("|------|----------|----------|----------|\n")

			for _, goal := range phase.GoalViews {
				maturity := fmt.Sprintf("L%d → L%d", goal.EnterLevel, goal.ExitLevel)
				if goal.CurrentLevel > 0 {
					maturity = fmt.Sprintf("L%d → L%d (now L%d)", goal.EnterLevel, goal.ExitLevel, goal.CurrentLevel)
				}
				progress := fmt.Sprintf("%d/%d", goal.InitiativesCompleted, goal.InitiativesTotal)
				complete := fmt.Sprintf("%.0f%%", goal.CompletionPercent)
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
					goal.GoalName, maturity, progress, complete))
			}
			sb.WriteString("\n")

			// Initiative details per goal
			if opts.IncludeDetails {
				for _, goal := range phase.GoalViews {
					if len(goal.Initiatives) > 0 {
						sb.WriteString(fmt.Sprintf("#### %s Initiatives\n\n", goal.GoalName))
						writeInitiativeTable(sb, goal.Initiatives)
					}
				}
			}
		}
	}
}

func writeGoalView(sb *strings.Builder, report *prism.RoadmapReport, opts *MarkdownOptions) {
	sb.WriteString("# Roadmap by Goal\n\n")
	sb.WriteString("This view shows progress organized by strategic goal (Goal → Phase → Initiative).\n\n")

	for _, goal := range report.ByGoal {
		// Goal header
		sb.WriteString(fmt.Sprintf("## %s\n\n", goal.GoalName))

		if goal.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", goal.Description))
		}

		// Goal metadata
		sb.WriteString(fmt.Sprintf("- **Current Level:** %d\n", goal.CurrentLevel))
		sb.WriteString(fmt.Sprintf("- **Target Level:** %d\n\n", goal.TargetLevel))

		// Phase progress table
		if len(goal.PhaseProgress) > 0 {
			sb.WriteString("### Phase Progress\n\n")
			sb.WriteString("| Phase | Maturity Target | Progress | Complete |\n")
			sb.WriteString("|-------|-----------------|----------|----------|\n")

			for _, phase := range goal.PhaseProgress {
				phaseName := phase.PhaseName
				if phase.Quarter != "" && phase.Year > 0 {
					phaseName = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
				}
				maturity := fmt.Sprintf("L%d → L%d", phase.EnterLevel, phase.ExitLevel)
				progress := fmt.Sprintf("%d/%d", phase.InitiativesCompleted, phase.InitiativesTotal)
				complete := fmt.Sprintf("%.0f%%", phase.CompletionPercent)
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
					phaseName, maturity, progress, complete))
			}
			sb.WriteString("\n")

			// Initiative details per phase
			if opts.IncludeDetails {
				for _, phase := range goal.PhaseProgress {
					if len(phase.Initiatives) > 0 {
						phaseName := phase.PhaseName
						if phase.Quarter != "" && phase.Year > 0 {
							phaseName = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
						}
						sb.WriteString(fmt.Sprintf("#### %s Initiatives\n\n", phaseName))
						writeInitiativeTable(sb, phase.Initiatives)
					}
				}
			}
		}
	}
}

func writeInitiativeTable(sb *strings.Builder, initiatives []prism.InitiativeSummary) {
	sb.WriteString("| Initiative | Team | Status | Dev % |\n")
	sb.WriteString("|------------|------|--------|-------|\n")

	for _, init := range initiatives {
		team := init.Team
		if team == "" {
			team = "-"
		}
		status := formatStatus(init.Status)
		devPct := fmt.Sprintf("%.0f%%", init.DevCompletionPercent)
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
			init.Name, team, status, devPct))
	}
	sb.WriteString("\n")
}

func formatStatus(status string) string {
	switch status {
	case "completed":
		return "✅ Completed"
	case "in_progress":
		return "🔄 In Progress"
	case "planned", "planning":
		return "📋 Planned"
	case "not_started":
		return "⏳ Not Started"
	case "cancelled":
		return "❌ Cancelled"
	default:
		if status == "" {
			return "-"
		}
		return status
	}
}

func progressBar(percent float64) string {
	const width = 20
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("`[%s]` %.1f%%", bar, percent)
}

func escapeYAML(s string) string {
	// Escape quotes for YAML
	return strings.ReplaceAll(s, "\"", "\\\"")
}
