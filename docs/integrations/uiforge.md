# UIForge Integration

!!! note "Coming Soon"
    UIForge integration is planned for a future release.

## Overview

UIForge integration enables PRISM dashboards as:

- **Standalone pages** - Full-page PRISM dashboards
- **Embedded widgets** - PRISM components in UIForge sites
- **Trend visualization** - Historical score tracking

## Planned Features

### Standalone Dashboard

A full-page dashboard showing:

- Overall PRISM score with interpretation
- Domain scores (Security, Operations)
- Stage breakdown heatmap
- Metric status table
- Score trends over time

### Dashboard Widgets

Embeddable components for UIForge sites:

| Widget | Description |
|--------|-------------|
| Score Gauge | PRISM score with color indicator |
| Domain Cards | Security/Operations score cards |
| Heatmap | Domain × Stage matrix |
| Trend Chart | Score history line chart |
| Metric Table | Filterable metric list |

### Dashboard Layout

The maturity dashboard is organized into hierarchical views:

#### Executive Summary

Top-level metrics and domain overview:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  OPERATIONS MATURITY DASHBOARD                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐  │
│  │ OVERALL MATURITY    │  │ SLO COMPLIANCE      │  │ METHODOLOGY         │  │
│  │                     │  │                     │  │                     │  │
│  │   ████ 3.2 / 5      │  │   ████░ 78%         │  │ Golden Signals: 4/4 │  │
│  │   ▰▰▰▱▱             │  │   Target: 95%       │  │ RED: 3/3            │  │
│  │   🟡 Defined        │  │   🟡 Below Target   │  │ USE: 3/3            │  │
│  └─────────────────────┘  └─────────────────────┘  └─────────────────────┘  │
│                                                                             │
│  DOMAIN MATURITY                                                            │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │ Reliability     [██████████░░░░░░░░░░] M3  🟢 Target: M3               ││
│  │ Incident Mgmt   [████████░░░░░░░░░░░░] M2  🟡 Target: M3               ││
│  │ Deployment      [████████████████░░░░] M4  🟢 Target: M4               ││
│  │ Monitoring      [██████░░░░░░░░░░░░░░] M2  🔴 Target: M4               ││
│  └─────────────────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Domain Drill-Down with Bullet Charts

Detailed view for each domain showing SLI progress toward maturity levels:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  RELIABILITY DOMAIN                                          Current: M3    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Availability SLI          ┌──┬──┬──┬──┬──┐                                 │
│  Current: 99.5%            │M1│M2│M3│M4│M5│  🟢 Meeting M3 (99%)            │
│                            │██│██│██│░░│░░│                                 │
│                            └──┴──┴──┴──┴──┘                                 │
│                                                                             │
│  Error Rate SLI            ┌──┬──┬──┬──┬──┐                                 │
│  Current: 0.3%             │M1│M2│M3│M4│M5│  🟡 Between M2-M3               │
│                            │██│██│▓▓│░░│░░│                                 │
│                            └──┴──┴──┴──┴──┘                                 │
│                                                                             │
│  Latency P99 SLI           ┌──┬──┬──┬──┬──┐                                 │
│  Current: 180ms            │M1│M2│M3│M4│M5│  🟢 Meeting M3 (200ms)          │
│                            │██│██│██│░░│░░│                                 │
│                            └──┴──┴──┴──┴──┘                                 │
│                                                                             │
│  Throughput SLI            ┌──┬──┬──┬──┬──┐                                 │
│  Current: 850 rps          │M1│M2│M3│M4│M5│  🟡 Between M2-M3               │
│                            │██│██│▓▓│░░│░░│                                 │
│                            └──┴──┴──┴──┴──┘                                 │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### SLI Category View (Grouped by SLI Type)

Metrics grouped by observability type:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  SLI METRICS BY TYPE                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  LATENCY                                   THROUGHPUT                       │
│  ├─ API Response P99      180ms    🟢 M3   ├─ Request Rate     850 rps  🟡  │
│  ├─ DB Query P95           45ms    🟢 M3   ├─ Event Throughput 12k/min  🟢  │
│  └─ Cache Hit Latency       2ms    🟢 M4   └─ Batch Processing   5k/hr  🟡  │
│                                                                             │
│  ERROR RATE                                SATURATION                       │
│  ├─ API Error Rate        0.3%     🟢 M3   ├─ CPU Utilization     65%   🟢  │
│  ├─ Failed Deployments    2.1%     🟡 M2   ├─ Memory Pressure     72%   🟡  │
│  └─ Data Validation        0.1%    🟢 M4   └─ Queue Depth          45   🟢  │
│                                                                             │
│  AVAILABILITY                                                               │
│  ├─ Service Uptime       99.5%     🟢 M3                                    │
│  ├─ API Availability     99.8%     🟢 M4                                    │
│  └─ Database Uptime      99.9%     🟢 M4                                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Methodology Coverage

Coverage analysis for standard observability methodologies:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  METHODOLOGY COVERAGE                                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  GOLDEN SIGNALS (Google SRE)                              Coverage: 100%    │
│  ┌────────────────┬────────────────┬────────────────┬────────────────┐      │
│  │   LATENCY      │   TRAFFIC      │   ERRORS       │  SATURATION    │      │
│  │   ████ 🟢      │   ███░ 🟡      │   ████ 🟢      │   ███░ 🟡      │      │
│  │   3 metrics    │   2 metrics    │   3 metrics    │   3 metrics    │      │
│  └────────────────┴────────────────┴────────────────┴────────────────┘      │
│                                                                             │
│  RED METHOD (Microservices)                               Coverage: 100%    │
│  ┌──────────────────────┬──────────────────────┬──────────────────────┐     │
│  │       RATE           │       ERRORS         │       DURATION       │     │
│  │       ███░ 🟡        │       ████ 🟢        │       ████ 🟢        │     │
│  │       2 metrics      │       3 metrics      │       3 metrics      │     │
│  └──────────────────────┴──────────────────────┴──────────────────────┘     │
│                                                                             │
│  USE METHOD (Resources)                                   Coverage: 100%    │
│  ┌──────────────────────┬──────────────────────┬──────────────────────┐     │
│  │     UTILIZATION      │      SATURATION      │       ERRORS         │     │
│  │       ███░ 🟡        │       ███░ 🟡        │       ████ 🟢        │     │
│  │       3 metrics      │       3 metrics      │       3 metrics      │     │
│  └──────────────────────┴──────────────────────┴──────────────────────┘     │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Heatmap Visualization

A domain × stage heatmap showing cell scores:

```
           │ Design │ Build │ Test │ Runtime │ Response │
───────────┼────────┼───────┼──────┼─────────┼──────────┤
Security   │  🟡 65 │ 🟢 92 │ 🟡 70│  🟢 88  │   🟡 72  │
───────────┼────────┼───────┼──────┼─────────┼──────────┤
Operations │  🟡 68 │ 🟢 85 │ 🟡 75│  🟢 95  │   🟢 82  │
```

## Configuration

### Document Configuration

```json
{
  "integrations": {
    "uiforge": {
      "enabled": true,
      "theme": "default",
      "refreshInterval": "1h",
      "widgets": ["score", "heatmap", "trends"]
    }
  }
}
```

### UIForge Site Configuration

```yaml
# uiforge.yml
pages:
  - name: PRISM Dashboard
    type: prism
    source: prism.json
    layout: full
    widgets:
      - type: prism-score
        position: top
      - type: prism-heatmap
        position: center
      - type: prism-trends
        position: bottom
```

## Planned CLI Commands

```bash
# Generate standalone dashboard
prism dashboard prism.json -o dashboard.html

# Generate UIForge widget data
prism dashboard prism.json --format uiforge -o prism-widget.json

# Start live dashboard server
prism serve prism.json --port 8080
```

## Integration with UIForge

### Embedding in UIForge Site

```markdown
<!-- In a UIForge page -->
# Platform Health

{{< prism-score source="prism.json" >}}

## Domain Breakdown

{{< prism-heatmap source="prism.json" >}}
```

### API Access

```go
import (
    "github.com/grokify/prism"
    "github.com/plexusone/uiforge"
)

// Load PRISM document
doc, _ := prism.LoadDocument("prism.json")

// Generate UIForge widget
widget := uiforge.NewPRISMWidget(doc)
widget.Render(w)
```

## Roadmap

1. **Phase 1**: Static dashboard generation
2. **Phase 2**: UIForge widget integration
3. **Phase 3**: Live data refresh
4. **Phase 4**: Historical trend storage
