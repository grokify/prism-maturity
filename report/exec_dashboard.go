package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/grokify/prism-maturity"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// titleCase converts a string to title case.
func titleCase(s string) string {
	return cases.Title(language.English).String(s)
}

// DashboardOptions configures dashboard report generation.
type DashboardOptions struct {
	Title           string
	Author          string
	Date            string
	IncludeYAMLMeta bool   // Include YAML front matter for Pandoc
	IncludeGaps     bool   // Include gap analysis section
	MaxGaps         int    // Maximum gaps to show (0 = all)
	Theme           string // HTML theme: "light", "dark"
}

// DefaultDashboardOptions returns sensible defaults.
func DefaultDashboardOptions() *DashboardOptions {
	return &DashboardOptions{
		Title:           "Executive Security Dashboard",
		Date:            time.Now().Format("2006-01-02"),
		IncludeYAMLMeta: true,
		IncludeGaps:     true,
		MaxGaps:         10,
		Theme:           "light",
	}
}

// GenerateDashboardMarkdown generates a Pandoc-compatible Markdown dashboard.
func GenerateDashboardMarkdown(dashboard *prism.ExecutiveDashboard, opts *DashboardOptions) string {
	if opts == nil {
		opts = DefaultDashboardOptions()
	}

	var sb strings.Builder

	// YAML front matter for Pandoc
	if opts.IncludeYAMLMeta {
		sb.WriteString("---\n")
		sb.WriteString(fmt.Sprintf("title: \"%s\"\n", escapeYAML(dashboard.Title)))
		if opts.Author != "" {
			sb.WriteString(fmt.Sprintf("author: \"%s\"\n", escapeYAML(opts.Author)))
		}
		sb.WriteString(fmt.Sprintf("date: \"%s\"\n", opts.Date))
		sb.WriteString("---\n\n")
	}

	// Executive Summary
	sb.WriteString("# Executive Summary\n\n")
	writeSummarySection(&sb, dashboard)

	// Maturity Scorecard
	sb.WriteString("\n# Maturity Scorecard\n\n")
	writeMaturityScorecard(&sb, dashboard)

	// SLO Compliance
	sb.WriteString("\n# SLO Compliance\n\n")
	writeSLOCompliance(&sb, dashboard)

	// Phase Progress
	sb.WriteString("\n# Roadmap Progress\n\n")
	writePhaseProgress(&sb, dashboard)

	// Gap Analysis
	if opts.IncludeGaps && len(dashboard.Gaps) > 0 {
		sb.WriteString("\n# Gap Analysis\n\n")
		writeGapAnalysis(&sb, dashboard, opts.MaxGaps)
	}

	return sb.String()
}

// GenerateDashboardMarp generates a Marp presentation dashboard.
func GenerateDashboardMarp(dashboard *prism.ExecutiveDashboard, opts *DashboardOptions) string {
	if opts == nil {
		opts = DefaultDashboardOptions()
	}

	var sb strings.Builder

	// Marp front matter
	sb.WriteString("---\n")
	sb.WriteString("marp: true\n")
	sb.WriteString("theme: default\n")
	sb.WriteString("paginate: true\n")
	sb.WriteString("backgroundColor: #fff\n")
	sb.WriteString("---\n\n")

	// Title slide
	sb.WriteString(fmt.Sprintf("# %s\n\n", dashboard.Title))
	if dashboard.Subtitle != "" {
		sb.WriteString(fmt.Sprintf("## %s\n\n", dashboard.Subtitle))
	}
	sb.WriteString(fmt.Sprintf("**%s**\n\n", opts.Date))
	sb.WriteString("---\n\n")

	// Executive Summary slide
	sb.WriteString("# Executive Summary\n\n")
	summary := dashboard.Summary

	sb.WriteString("<div style=\"display: flex; justify-content: space-around;\">\n\n")

	// Key metrics in boxes
	sb.WriteString(fmt.Sprintf("**Overall Maturity**\n%.1f / %.1f\n\n",
		summary.OverallMaturity, summary.TargetMaturity))
	sb.WriteString(fmt.Sprintf("**SLO Compliance**\n%.0f%%\n\n",
		dashboard.SLOCompliance.OverallCompliance))
	sb.WriteString(fmt.Sprintf("**Goals On Track**\n%d / %d\n\n",
		summary.GoalsOnTrack, summary.GoalsTotal))

	sb.WriteString("</div>\n\n")
	sb.WriteString("---\n\n")

	// Maturity Scorecard slide
	sb.WriteString("# Maturity Scorecard\n\n")
	sb.WriteString("| Goal | Current | Target | Gap | Status |\n")
	sb.WriteString("|------|:-------:|:------:|:---:|:------:|\n")

	for _, g := range dashboard.MaturityScorecard {
		status := statusEmoji(g.Status)
		sb.WriteString(fmt.Sprintf("| %s | L%d | L%d | %d | %s |\n",
			g.GoalName, g.CurrentLevel, g.TargetLevel, g.Gap, status))
	}
	sb.WriteString("\n---\n\n")

	// SLO Compliance slide
	sb.WriteString("# SLO Compliance by Category\n\n")
	sb.WriteString("| Category | Met | Missed | Compliance |\n")
	sb.WriteString("|----------|:---:|:------:|:----------:|\n")

	for _, cat := range dashboard.SLOCompliance.Categories {
		complianceBar := complianceIndicator(cat.Compliance)
		sb.WriteString(fmt.Sprintf("| %s | %d | %d | %s %.0f%% |\n",
			titleCase(cat.Category), cat.Met, cat.Missed, complianceBar, cat.Compliance))
	}
	sb.WriteString("\n---\n\n")

	// Phase Progress slide
	sb.WriteString("# Roadmap Progress\n\n")
	for _, phase := range dashboard.PhaseProgress {
		phaseName := phase.PhaseName
		if phase.Quarter != "" && phase.Year > 0 {
			phaseName = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
		}

		statusIcon := "⏳"
		if phase.Status == "completed" {
			statusIcon = "✅"
		} else if phase.IsCurrent {
			statusIcon = "🔄"
		}

		sb.WriteString(fmt.Sprintf("### %s %s\n\n", statusIcon, phaseName))
		sb.WriteString(fmt.Sprintf("%s %.0f%% complete\n\n", progressBarMarp(phase.CompletionPct), phase.CompletionPct))
	}
	sb.WriteString("\n---\n\n")

	// Gap Analysis slide (top 5)
	if opts.IncludeGaps && len(dashboard.Gaps) > 0 {
		sb.WriteString("# Priority Gaps\n\n")
		sb.WriteString("| Priority | Metric | Current | Target | Gap |\n")
		sb.WriteString("|:--------:|--------|:-------:|:------:|:---:|\n")

		maxGaps := 5
		if len(dashboard.Gaps) < maxGaps {
			maxGaps = len(dashboard.Gaps)
		}

		for i := 0; i < maxGaps; i++ {
			gap := dashboard.Gaps[i]
			priorityIcon := priorityEmoji(gap.Priority)
			sb.WriteString(fmt.Sprintf("| %s | %s | %.1f | %.1f | %.1f%% |\n",
				priorityIcon, gap.MetricName, gap.CurrentVal, gap.TargetVal, gap.GapPercent))
		}
		sb.WriteString("\n---\n\n")
	}

	// Thank you slide
	sb.WriteString("# Questions?\n\n")
	sb.WriteString(fmt.Sprintf("Dashboard generated: %s\n", opts.Date))

	return sb.String()
}

// GenerateDashboardHTML generates a standalone HTML dashboard.
func GenerateDashboardHTML(dashboard *prism.ExecutiveDashboard, opts *DashboardOptions) string {
	if opts == nil {
		opts = DefaultDashboardOptions()
	}

	var sb strings.Builder

	// HTML header with embedded CSS
	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>`)
	sb.WriteString(dashboard.Title)
	sb.WriteString(`</title>
    <style>
        :root {
            --bg-primary: #ffffff;
            --bg-secondary: #f8f9fa;
            --text-primary: #212529;
            --text-secondary: #6c757d;
            --border-color: #dee2e6;
            --success: #28a745;
            --warning: #ffc107;
            --danger: #dc3545;
            --info: #17a2b8;
        }

        * { box-sizing: border-box; }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background-color: var(--bg-secondary);
            color: var(--text-primary);
            margin: 0;
            padding: 20px;
            line-height: 1.6;
        }

        .dashboard {
            max-width: 1400px;
            margin: 0 auto;
        }

        .header {
            text-align: center;
            margin-bottom: 30px;
        }

        .header h1 {
            margin: 0;
            font-size: 2rem;
        }

        .header .subtitle {
            color: var(--text-secondary);
            margin-top: 5px;
        }

        .summary-cards {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }

        .card {
            background: var(--bg-primary);
            border-radius: 8px;
            padding: 20px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }

        .card.full-width {
            grid-column: 1 / -1;
        }

        .card h2 {
            margin: 0 0 15px 0;
            font-size: 1.1rem;
            color: var(--text-secondary);
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .card h3 {
            margin: 20px 0 10px 0;
            font-size: 1rem;
        }

        .metric-value {
            font-size: 2.5rem;
            font-weight: bold;
        }

        .metric-label {
            color: var(--text-secondary);
            font-size: 0.9rem;
        }

        .status-green { color: var(--success); }
        .status-yellow { color: var(--warning); }
        .status-red { color: var(--danger); }

        .progress-bar {
            background: var(--bg-secondary);
            border-radius: 4px;
            height: 24px;
            overflow: hidden;
            margin: 10px 0;
        }

        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, var(--info), var(--success));
            display: flex;
            align-items: center;
            justify-content: flex-end;
            padding-right: 8px;
            color: white;
            font-size: 0.8rem;
            font-weight: bold;
            min-width: 40px;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 10px;
        }

        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid var(--border-color);
        }

        th {
            background: var(--bg-secondary);
            font-weight: 600;
            text-transform: uppercase;
            font-size: 0.8rem;
            letter-spacing: 0.5px;
        }

        tr:hover {
            background: var(--bg-secondary);
        }

        .badge {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.75rem;
            font-weight: bold;
            text-transform: uppercase;
        }

        .badge-success { background: #d4edda; color: #155724; }
        .badge-warning { background: #fff3cd; color: #856404; }
        .badge-danger { background: #f8d7da; color: #721c24; }
        .badge-info { background: #d1ecf1; color: #0c5460; }

        .priority-critical { color: var(--danger); font-weight: bold; }
        .priority-high { color: #e67700; }
        .priority-medium { color: var(--warning); }
        .priority-low { color: var(--text-secondary); }

        .grid-2 {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
            gap: 20px;
        }

        .heatmap {
            display: grid;
            grid-template-columns: repeat(5, 1fr);
            gap: 4px;
            margin-top: 10px;
        }

        .heatmap-cell {
            padding: 8px;
            text-align: center;
            border-radius: 4px;
            font-size: 0.8rem;
        }

        .heatmap-header {
            background: var(--bg-secondary);
            font-weight: bold;
        }

        .footer {
            text-align: center;
            color: var(--text-secondary);
            margin-top: 30px;
            font-size: 0.85rem;
        }
    </style>
</head>
<body>
    <div class="dashboard">
`)

	// Header
	sb.WriteString(fmt.Sprintf(`        <div class="header">
            <h1>%s</h1>
            <div class="subtitle">%s</div>
        </div>
`, dashboard.Title, dashboard.Subtitle))

	// Summary cards
	summary := dashboard.Summary
	sb.WriteString(`        <div class="summary-cards">
`)

	// Overall Maturity card
	maturityClass := "status-green"
	if summary.MaturityGap > 1 {
		maturityClass = "status-red"
	} else if summary.MaturityGap > 0 {
		maturityClass = "status-yellow"
	}
	sb.WriteString(fmt.Sprintf(`            <div class="card">
                <h2>Overall Maturity</h2>
                <div class="metric-value %s">L%.1f</div>
                <div class="metric-label">Target: L%.1f (Gap: %.1f)</div>
            </div>
`, maturityClass, summary.OverallMaturity, summary.TargetMaturity, summary.MaturityGap))

	// SLO Compliance card
	complianceClass := "status-green"
	if dashboard.SLOCompliance.OverallCompliance < 70 {
		complianceClass = "status-red"
	} else if dashboard.SLOCompliance.OverallCompliance < 90 {
		complianceClass = "status-yellow"
	}
	sb.WriteString(fmt.Sprintf(`            <div class="card">
                <h2>SLO Compliance</h2>
                <div class="metric-value %s">%.0f%%</div>
                <div class="metric-label">%d of %d SLOs met</div>
            </div>
`, complianceClass, dashboard.SLOCompliance.OverallCompliance,
		dashboard.SLOCompliance.OverallMet, dashboard.SLOCompliance.OverallTotal))

	// Goals card
	goalsClass := "status-green"
	if summary.GoalsOnTrack < summary.GoalsTotal/2 {
		goalsClass = "status-red"
	} else if summary.GoalsOnTrack < summary.GoalsTotal {
		goalsClass = "status-yellow"
	}
	sb.WriteString(fmt.Sprintf(`            <div class="card">
                <h2>Goals On Track</h2>
                <div class="metric-value %s">%d/%d</div>
                <div class="metric-label">%d at risk or behind</div>
            </div>
`, goalsClass, summary.GoalsOnTrack, summary.GoalsTotal, summary.GoalsAtRisk))

	// Initiatives card
	sb.WriteString(fmt.Sprintf(`            <div class="card">
                <h2>Initiatives</h2>
                <div class="metric-value">%d</div>
                <div class="metric-label">%d completed, %d in progress</div>
            </div>
`, summary.InitiativesTotal, summary.InitiativesCompleted, summary.InitiativesInProgress))

	sb.WriteString(`        </div>

        <div class="grid-2">
`)

	// Maturity Scorecard
	sb.WriteString(`            <div class="card">
                <h2>Maturity Scorecard</h2>
                <table>
                    <thead>
                        <tr>
                            <th>Goal</th>
                            <th>Current</th>
                            <th>Target</th>
                            <th>Gap</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody>
`)

	for _, g := range dashboard.MaturityScorecard {
		badgeClass := "badge-success"
		statusText := "On Track"
		if g.Status == "at_risk" {
			badgeClass = "badge-warning"
			statusText = "At Risk"
		} else if g.Status == "behind" {
			badgeClass = "badge-danger"
			statusText = "Behind"
		}
		sb.WriteString(fmt.Sprintf(`                        <tr>
                            <td>%s</td>
                            <td>L%d</td>
                            <td>L%d</td>
                            <td>%d</td>
                            <td><span class="badge %s">%s</span></td>
                        </tr>
`, g.GoalName, g.CurrentLevel, g.TargetLevel, g.Gap, badgeClass, statusText))
	}

	sb.WriteString(`                    </tbody>
                </table>
            </div>
`)

	// SLO Compliance by Category
	sb.WriteString(`            <div class="card">
                <h2>SLO Compliance by Category</h2>
                <table>
                    <thead>
                        <tr>
                            <th>Category</th>
                            <th>Met</th>
                            <th>Missed</th>
                            <th>Compliance</th>
                        </tr>
                    </thead>
                    <tbody>
`)

	for _, cat := range dashboard.SLOCompliance.Categories {
		compClass := "status-green"
		if cat.Compliance < 70 {
			compClass = "status-red"
		} else if cat.Compliance < 90 {
			compClass = "status-yellow"
		}
		sb.WriteString(fmt.Sprintf(`                        <tr>
                            <td>%s</td>
                            <td>%d</td>
                            <td>%d</td>
                            <td class="%s">%.0f%%</td>
                        </tr>
`, titleCase(cat.Category), cat.Met, cat.Missed, compClass, cat.Compliance))
	}

	sb.WriteString(`                    </tbody>
                </table>
            </div>
        </div>
`)

	// Phase Progress
	sb.WriteString(`
        <div class="card full-width" style="margin-top: 20px;">
            <h2>Roadmap Progress</h2>
`)

	for _, phase := range dashboard.PhaseProgress {
		phaseName := phase.PhaseName
		if phase.Quarter != "" && phase.Year > 0 {
			phaseName = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
		}

		statusBadge := "badge-info"
		statusText := "Planned"
		if phase.Status == "completed" {
			statusBadge = "badge-success"
			statusText = "Completed"
		} else if phase.IsCurrent {
			statusBadge = "badge-warning"
			statusText = "In Progress"
		}

		sb.WriteString(fmt.Sprintf(`            <h3>%s <span class="badge %s">%s</span></h3>
            <div class="progress-bar">
                <div class="progress-fill" style="width: %.0f%%;">%.0f%%</div>
            </div>
            <div class="metric-label">%d/%d initiatives complete | %d/%d goals achieved</div>
`, phaseName, statusBadge, statusText, phase.CompletionPct, phase.CompletionPct,
			phase.InitCompleted, phase.InitTotal, phase.GoalsAchieved, phase.GoalsTargeted))
	}

	sb.WriteString(`        </div>
`)

	// Gap Analysis
	if opts.IncludeGaps && len(dashboard.Gaps) > 0 {
		sb.WriteString(`
        <div class="card full-width" style="margin-top: 20px;">
            <h2>Priority Gaps</h2>
            <table>
                <thead>
                    <tr>
                        <th>Priority</th>
                        <th>Metric</th>
                        <th>Category</th>
                        <th>Current</th>
                        <th>Target</th>
                        <th>Gap %</th>
                        <th>Goal</th>
                    </tr>
                </thead>
                <tbody>
`)

		maxGaps := opts.MaxGaps
		if maxGaps == 0 || maxGaps > len(dashboard.Gaps) {
			maxGaps = len(dashboard.Gaps)
		}

		for i := 0; i < maxGaps; i++ {
			gap := dashboard.Gaps[i]
			priorityClass := fmt.Sprintf("priority-%s", gap.Priority)
			sb.WriteString(fmt.Sprintf(`                    <tr>
                        <td class="%s">%s</td>
                        <td>%s</td>
                        <td>%s</td>
                        <td>%.1f</td>
                        <td>%.1f</td>
                        <td>%.1f%%</td>
                        <td>%s</td>
                    </tr>
`, priorityClass, strings.ToUpper(gap.Priority), gap.MetricName, titleCase(gap.Category),
				gap.CurrentVal, gap.TargetVal, gap.GapPercent, gap.GoalName))
		}

		sb.WriteString(`                </tbody>
            </table>
        </div>
`)
	}

	// Footer
	sb.WriteString(fmt.Sprintf(`
        <div class="footer">
            Generated: %s | PRISM Executive Dashboard
        </div>
    </div>
</body>
</html>
`, opts.Date))

	return sb.String()
}

// Helper functions

func writeSummarySection(sb *strings.Builder, dashboard *prism.ExecutiveDashboard) {
	summary := dashboard.Summary

	sb.WriteString("## Key Metrics\n\n")
	sb.WriteString("| Metric | Value | Status |\n")
	sb.WriteString("|--------|-------|--------|\n")

	// Overall Maturity
	maturityStatus := "🟢 On Track"
	if summary.MaturityGap > 1 {
		maturityStatus = "🔴 Behind"
	} else if summary.MaturityGap > 0 {
		maturityStatus = "🟡 At Risk"
	}
	sb.WriteString(fmt.Sprintf("| Overall Maturity | L%.1f / L%.1f | %s |\n",
		summary.OverallMaturity, summary.TargetMaturity, maturityStatus))

	// SLO Compliance
	sloStatus := "🟢 Healthy"
	if dashboard.SLOCompliance.OverallCompliance < 70 {
		sloStatus = "🔴 Critical"
	} else if dashboard.SLOCompliance.OverallCompliance < 90 {
		sloStatus = "🟡 Needs Attention"
	}
	sb.WriteString(fmt.Sprintf("| SLO Compliance | %.0f%% | %s |\n",
		dashboard.SLOCompliance.OverallCompliance, sloStatus))

	// Goals
	sb.WriteString(fmt.Sprintf("| Goals On Track | %d / %d | %s |\n",
		summary.GoalsOnTrack, summary.GoalsTotal,
		fmt.Sprintf("%d at risk", summary.GoalsAtRisk)))

	// Initiatives
	sb.WriteString(fmt.Sprintf("| Initiatives | %d complete, %d in progress | %d total |\n",
		summary.InitiativesCompleted, summary.InitiativesInProgress, summary.InitiativesTotal))

	sb.WriteString("\n")
}

func writeMaturityScorecard(sb *strings.Builder, dashboard *prism.ExecutiveDashboard) {
	sb.WriteString("| Goal | Current | Target | Gap | SLOs Met | Status |\n")
	sb.WriteString("|------|:-------:|:------:|:---:|:--------:|:------:|\n")

	for _, g := range dashboard.MaturityScorecard {
		status := statusEmoji(g.Status)
		slos := fmt.Sprintf("%d/%d (%.0f%%)", g.SLOsMet, g.SLOsTotal, g.SLOsMetPercent)
		sb.WriteString(fmt.Sprintf("| %s | L%d | L%d | %d | %s | %s |\n",
			g.GoalName, g.CurrentLevel, g.TargetLevel, g.Gap, slos, status))
	}
	sb.WriteString("\n")
}

func writeSLOCompliance(sb *strings.Builder, dashboard *prism.ExecutiveDashboard) {
	sb.WriteString("| Category | Met | At Risk | Missed | Compliance |\n")
	sb.WriteString("|----------|:---:|:-------:|:------:|:----------:|\n")

	for _, cat := range dashboard.SLOCompliance.Categories {
		indicator := complianceIndicator(cat.Compliance)
		sb.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %s %.0f%% |\n",
			titleCase(cat.Category), cat.Met, cat.AtRisk, cat.Missed, indicator, cat.Compliance))
	}

	sb.WriteString(fmt.Sprintf("\n**Overall Compliance:** %.0f%% (%d/%d SLOs met)\n\n",
		dashboard.SLOCompliance.OverallCompliance,
		dashboard.SLOCompliance.OverallMet,
		dashboard.SLOCompliance.OverallTotal))
}

func writePhaseProgress(sb *strings.Builder, dashboard *prism.ExecutiveDashboard) {
	sb.WriteString("| Phase | Status | Progress | Goals | Initiatives |\n")
	sb.WriteString("|-------|:------:|:--------:|:-----:|:-----------:|\n")

	for _, phase := range dashboard.PhaseProgress {
		phaseName := phase.PhaseName
		if phase.Quarter != "" && phase.Year > 0 {
			phaseName = fmt.Sprintf("%s %d", phase.Quarter, phase.Year)
		}

		statusIcon := "⏳ Planned"
		if phase.Status == "completed" {
			statusIcon = "✅ Complete"
		} else if phase.IsCurrent {
			statusIcon = "🔄 Current"
		}

		progress := progressBarText(phase.CompletionPct)
		goals := fmt.Sprintf("%d/%d", phase.GoalsAchieved, phase.GoalsTargeted)
		inits := fmt.Sprintf("%d/%d", phase.InitCompleted, phase.InitTotal)

		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			phaseName, statusIcon, progress, goals, inits))
	}
	sb.WriteString("\n")
}

func writeGapAnalysis(sb *strings.Builder, dashboard *prism.ExecutiveDashboard, maxGaps int) {
	sb.WriteString("*Gaps sorted by priority (largest impact first)*\n\n")
	sb.WriteString("| Priority | Metric | Current | Target | Gap % | Goal |\n")
	sb.WriteString("|:--------:|--------|:-------:|:------:|:-----:|------|\n")

	count := maxGaps
	if count == 0 || count > len(dashboard.Gaps) {
		count = len(dashboard.Gaps)
	}

	for i := 0; i < count; i++ {
		gap := dashboard.Gaps[i]
		priority := priorityEmoji(gap.Priority)
		sb.WriteString(fmt.Sprintf("| %s | %s | %.1f | %.1f | %.1f%% | %s |\n",
			priority, gap.MetricName, gap.CurrentVal, gap.TargetVal, gap.GapPercent, gap.GoalName))
	}
	sb.WriteString("\n")
}

func statusEmoji(status string) string {
	switch status {
	case "on_track":
		return "🟢"
	case "at_risk":
		return "🟡"
	case "behind":
		return "🔴"
	default:
		return "⚪"
	}
}

func priorityEmoji(priority string) string {
	switch priority {
	case "critical":
		return "🔴 CRITICAL"
	case "high":
		return "🟠 HIGH"
	case "medium":
		return "🟡 MEDIUM"
	case "low":
		return "🟢 LOW"
	default:
		return priority
	}
}

func complianceIndicator(pct float64) string {
	if pct >= 90 {
		return "🟢"
	} else if pct >= 70 {
		return "🟡"
	}
	return "🔴"
}

func progressBarText(pct float64) string {
	filled := int(pct / 10)
	if filled > 10 {
		filled = 10
	}
	return fmt.Sprintf("`[%s%s]` %.0f%%",
		strings.Repeat("█", filled),
		strings.Repeat("░", 10-filled),
		pct)
}

func progressBarMarp(pct float64) string {
	filled := int(pct / 5)
	if filled > 20 {
		filled = 20
	}
	return fmt.Sprintf("`%s%s`",
		strings.Repeat("█", filled),
		strings.Repeat("░", 20-filled))
}
