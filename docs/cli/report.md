# prism maturity plan report

Generate roadmap reports from a PRISM plan document.

## Synopsis

```bash
prism maturity plan report <plan-file> [options]
```

> **Note:** As of v0.6.0, `prism report` has been moved to `prism maturity plan report`.

## Description

Generate a roadmap report in Markdown or JSON format. The report can be generated in different views:

- **by-phase**: Phase → Goal → Initiative (timeline view)
- **by-goal**: Goal → Phase → Initiative (strategic view)
- **both**: Both views in a single document (default)

## Options

| Option | Description |
|--------|-------------|
| `-o`, `--output <file>` | Output file (default: stdout) |
| `-f`, `--format <format>` | Output format: `markdown`, `json` (default: markdown) |
| `-v`, `--view <view>` | View type: `both`, `by-phase`, `by-goal` (default: both) |
| `--title <title>` | Report title (default: from metadata) |
| `--author <author>` | Report author (default: from metadata) |
| `--no-meta` | Omit YAML front matter (Markdown only) |
| `--no-detail` | Omit initiative details |

## Examples

Generate markdown report to stdout:

```bash
prism maturity plan report plan.json
```

Generate markdown report to file:

```bash
prism maturity plan report plan.json -o report.md
```

Generate JSON report:

```bash
prism maturity plan report plan.json --format json
```

Generate phase-centric view only:

```bash
prism maturity plan report plan.json --view by-phase
```

Generate goal-centric view only:

```bash
prism maturity plan report plan.json --view by-goal
```

Generate report with custom title:

```bash
prism maturity plan report plan.json --title "Q1 2026 Roadmap Report"
```

## Output Format

### Markdown Output

The markdown output includes:

- YAML front matter (title, date, author)
- Summary section with goal and phase counts
- View sections (by-phase and/or by-goal)
- Initiative details with status and deployment info

### JSON Output

The JSON output includes a structured roadmap report object:

```json
{
  "title": "Roadmap Report",
  "generatedAt": "2026-04-18T10:00:00Z",
  "summary": {
    "totalGoals": 2,
    "totalPhases": 4,
    "totalInitiatives": 8
  },
  "goals": [...],
  "phases": [...]
}
```

## See Also

- [prism slo-report](slo-report.md) - SLO compliance reports
- [prism maturity plan dashboard](dashboard.md) - Executive dashboards
- [prism roadmap](roadmap.md) - Roadmap overview commands
