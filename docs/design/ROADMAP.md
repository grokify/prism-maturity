# PRISM Roadmap

## Vision

PRISM is an **Operational Product Management** platform for COO-level organizational health monitoring across multi-department, multi-function organizations.

**Primary Goals:**

1. **Understand Current Health** - Unified view of organizational maturity across all operational domains
2. **Drive Improvement Projects** - Prioritized roadmap to advance maturity and achieve SLOs

This is analogous to **Growth Product Management** where projects are prioritized to achieve growth metrics (MAU, DAU/MAU, activation, retention). PRISM enables **Operational Product Management** where projects are prioritized to achieve operational health metrics (SLOs, maturity levels, DORA metrics, security posture).

## Architecture

PRISM uses a three-part document model:

```
┌─────────────────────────────────────────────────────────────────┐
│                    prism-maturity-model                         │
│              "What does M4 availability mean?"                  │
│         SLIs, domains, levels, criteria, enablers               │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ references
┌─────────────────────────────┴───────────────────────────────────┐
│                    prism-maturity-state                         │
│              "We're at M3, availability is 99.7%"               │
│       Current values, temporal windows, history, gaps           │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ references
┌─────────────────────────────┴───────────────────────────────────┐
│                    prism-maturity-plan                          │
│              "Reach M4 by Q4 via these initiatives"             │
│         Goals, phases, initiatives, teams, roadmap              │
└─────────────────────────────────────────────────────────────────┘
```

## Target Users

| Role | Use Case |
|------|----------|
| **COO/CTO** | Executive dashboard showing organizational health |
| **VP Engineering** | Domain-level maturity and SLO compliance |
| **Engineering Manager** | Team metrics and improvement roadmap |
| **SRE/Platform Lead** | SLO definition, monitoring, alerting |
| **Security Lead** | Security maturity, vulnerability tracking |
| **QE Lead** | Quality metrics, test coverage, defect density |

## Completed Releases

### v0.7.0 (May 2026)

- Threshold Matrix XLSX sheet with M1-M5 columns
- SLI tags with kebab-case validation
- NIST CSF category ordering (govern → identify → protect → detect → respond → recover)
- Dashboard bullet charts grouped by category
- Dashboard layout improvements (flexbox labels, dynamic chart heights)

### v0.6.0 (May 2026)

- **Three-part schema taxonomy**: Model, State, Plan documents
- **CLI restructure**: `prism maturity model`, `prism maturity state`, `prism maturity plan`
- MaturityStateDocument type for state tracking
- Model lint command for dashboard display issues
- Dashboard qualitative state display
- Security methodology groupings (Prevention, Detection, Response)

### v0.5.0 (May 2026)

- Separated maturity model (definition) from maturity state (tracking)
- Temporal SLO window support (7d, 30d, 90d, quarterly, annual)
- Hybrid qualitative/quantitative criteria
- Dashboard dual-read for backward compatibility

### v0.4.0 (May 2026)

- SLI type with framework mappings
- SLI Catalog in Markdown reports
- Expanded framework support (NIST CSF 2.0, FedRAMP, SOC 2, ISO 27001)
- TypeScript/Zod schema for web viewer

### v0.3.0 (April 2026)

- Extended value stream layers (requirements → code → infra → runtime → adoption → support)
- Quality domain with ISO 25010 alignment
- Team Topologies support
- Analysis package for gap detection
- Export to OKR/V2MOM formats
- CLI commands: layer, team, service, analyze, export, initiative

### v0.2.0 (April 2026)

- Goal-driven maturity roadmap
- Phase-based planning with quarters and swimlanes
- Initiative tracking with deployment status
- CLI commands: goal, phase, roadmap, slo-report, dashboard
- Dashforge integration

### v0.1.0 (March 2026)

- Core data model (PRISMDocument, Metric, SLI, SLO)
- 5-level maturity model
- PRISM score calculation
- JSON Schema generation
- CLI: init, validate, score, catalog

## Planned: v0.8.0

### Dashboard Enhancements

- [ ] **No-state indicator** - Visual indicator when no state file is provided (currently shows 0/0/100)
- [ ] **Threshold Matrix in dashboard** - Add pivot-style threshold view to HTML dashboard
- [ ] **Tag filtering** - Filter SLIs by tag in the HTML dashboard
- [ ] **Category collapse/expand** - Collapsible category sections

### State Validation

- [ ] **Cross-validation CLI** - `prism maturity state validate` checks SLI IDs exist in model
- [ ] **State completeness report** - Identify SLIs in model without state values
- [ ] **State freshness warning** - Alert when state timestamps are stale

### Reporting

- [ ] **Markdown roadmap export** - Generate roadmap by phase or goal
- [ ] **Executive dashboard improvements** - Maturity trends, SLO compliance trends
- [ ] **SLO compliance enhancements** - Historical tracking, breach alerting, error budgets

## Planned: v0.9.0+

### Visualization & Automation

- [ ] ECharts visualizations (radar, gauge, timeline, heatmap)
- [ ] Automated metric collection (Datadog, Prometheus)
- [ ] SLO breach alerting
- [ ] Scheduled reporting

### Data Import/Export

- [ ] CSV import/export for metrics
- [ ] Datadog SLO import
- [ ] OpsLevel maturity import
- [ ] Grafana SLO import

### Dependencies

- [ ] Publish echartify, dashforge, omniframe packages
- [ ] Remove local replace directives from go.mod

## Related Projects

| Project | Repository | Purpose |
|---------|------------|---------|
| **PRISM** | `github.com/grokify/prism` | Health model, maturity, SLOs |
| **Structured Plan** | `github.com/grokify/structured-plan` | OKRs, V2MOM, roadmap |
| **Structured Changelog** | `github.com/grokify/structured-changelog` | Release management |

## Contributing

See [GitHub Issues](https://github.com/grokify/prism/issues) for detailed task tracking and discussion.
