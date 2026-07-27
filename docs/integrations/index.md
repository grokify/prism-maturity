# Integrations Overview

PRISM Maturity integrates with various tools and platforms for reporting, visualization, and export.

## Available Integrations

| Integration | Status | Description |
|-------------|--------|-------------|
| [PRISM Roadmap](prism-roadmap.md) | Available | Roadmap, OKR, and V2MOM export |
| [UIForge](uiforge.md) | Available | Dashboard generation and embedding |
| [MCP](mcp.md) | Available | AI assistant integration via Model Context Protocol |
| [Marp](marp.md) | Planned | Presentation generation |
| [Excel](excel.md) | Planned | XLSX export for stakeholders |

## Integration Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                 PRISM Intelligence Document                 │
│                       (prism.json)                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   PRISM Maturity Library                    │
│            (github.com/grokify/prism-maturity)              │
└─────────────────────────────────────────────────────────────┘
                              │
     ┌────────────┬───────────┼───────────┬────────────┐
     ▼            ▼           ▼           ▼            ▼
┌──────────┐┌──────────┐┌──────────┐┌──────────┐┌──────────┐
│  PRISM   ││ UIForge  ││   Marp   ││   Excel  ││   CLI    │
│Execution ││Dashboard ││  Slides  ││   XLSX   ││ Reports  │
└──────────┘└──────────┘└──────────┘└──────────┘└──────────┘
     │            │           │           │            │
     ▼            ▼           ▼           ▼            ▼
┌──────────┐┌──────────┐┌──────────┐┌──────────┐┌──────────┐
│ Roadmap  ││   HTML   ││   PDF    ││   XLSX   ││ Markdown │
│ OKR/V2MOM││Dashboard ││  Slides  ││  Report  ││  Reports │
└──────────┘└──────────┘└──────────┘└──────────┘└──────────┘
```

## Use Cases

### Executive Reporting

Generate presentations and spreadsheets for leadership:

1. Export PRISM score to Marp slides
2. Export detailed metrics to Excel
3. Embed trend charts in dashboards

### Team Dashboards

Create team-specific views:

1. Security team: Security domain dashboard
2. SRE team: Operations domain dashboard
3. Combined: Full PRISM dashboard

### Compliance Reporting

Generate framework-specific reports:

1. Filter metrics by framework mapping
2. Export to Excel with framework references
3. Generate audit-ready documentation

## Planned Features

### UIForge Integration

- Standalone PRISM dashboard pages
- Embedded dashboard widgets
- Trend visualization
- Domain/stage heatmaps

### Marp Integration

- Executive summary slides
- Score trend charts
- Domain breakdown visuals
- Metric status tables

### Excel Integration

- Full document export
- Filtered exports by domain/stage
- Score calculation worksheets
- Chart-ready data format

## Configuration

Future integrations will be configured in the PRISM document:

```json
{
  "metadata": {
    "name": "Acme PRISM",
    "version": "1.0.0"
  },
  "metrics": [...],
  "integrations": {
    "uiforge": {
      "enabled": true,
      "theme": "corporate",
      "refreshInterval": "1h"
    },
    "marp": {
      "enabled": true,
      "template": "executive"
    },
    "excel": {
      "enabled": true,
      "includeCharts": true
    }
  }
}
```

## Coming Soon

These integrations are in development. Check the [GitHub repository](https://github.com/grokify/prism-maturity) for updates.
