# PRISM

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/prism/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/prism/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/prism/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/prism/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/prism/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/prism/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/prism
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/prism
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/prism
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/prism
 [docs-mkdoc-svg]: https://img.shields.io/badge/docs-guide-blue.svg
 [docs-mkdoc-url]: https://grokify.dev/prism
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/prism/blob/main/LICENSE

**Platform for Reliability, Improvement, and Strategic Maturity**

PRISM is an Operational Product Management platform for COO-level organizational health monitoring. It provides a unified framework for B2B SaaS health metrics that combines SLOs, maturity modeling, and OKRs into a single coherent system. PRISM enables organizations to understand current health and drive improvement projects across operations, security, quality, product, and AI domains.

## Concepts

PRISM organizes metrics using a multi-dimensional model that clarifies ownership, accountability, and measurement across your organization.

### The PRISM Model

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              VALUE STREAM LAYERS                            │
│                          (Where in the value stream?)                       │
│                                                                             │
│  Requirements → Code → Infra → Runtime → Adoption → Support                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
        ┌─────────────────────────────┼─────────────────────────────┐
        │                             │                             │
        ▼                             ▼                             ▼
  ┌───────────┐                ┌───────────┐                ┌───────────┐
  │ OPERATIONS│                │ SECURITY  │                │  QUALITY  │
  │  Domain   │                │  Domain   │                │  Domain   │
  └───────────┘                └───────────┘                └───────────┘
        │                             │                             │
        └─────────────────────────────┼─────────────────────────────┘
                                      │
┌─────────────────────────────────────────────────────────────────────────────┐
│                           LIFECYCLE STAGES                                  │
│                        (When in delivery cycle?)                            │
│                                                                             │
│            Design → Build → Test → Runtime → Response                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Domains

Domains represent functional areas with their own standards and overlay teams:

| Domain | Description | Overlay Team |
|--------|-------------|--------------|
| `operations` | Reliability, performance, efficiency | SRE/Platform |
| `security` | AppSec, CloudSec, compliance | Security |
| `quality` | Testing, code quality, defects | QE |

### Layers

Layers represent the full value stream from ideation to support:

| Layer | Description | Typical Owner |
|-------|-------------|---------------|
| `requirements` | Product ideation, specs, design | Product/Design |
| `code` | Application code, libraries, dependencies | Stream-aligned teams |
| `infra` | Cloud resources, networking, platform | Platform team |
| `runtime` | Running services, production workloads | Stream-aligned + SRE |
| `adoption` | Product analytics, user engagement | Product/Growth |
| `support` | Customer support, incident management | Support/CS |

### Teams (Team Topologies)

PRISM supports Team Topologies patterns for clear accountability:

| Type | Role | Example |
|------|------|---------|
| `stream_aligned` | Build and run services end-to-end | Payments Team |
| `platform` | Provide infrastructure as a product | Platform Engineering |
| `enabling` | Help teams adopt new practices | Developer Experience |
| `overlay` | Define standards across organization | Security Team, QE Team |

### Services

Services are deployable units owned by teams with associated metrics:

```json
{
  "id": "payments-api",
  "name": "Payments API",
  "ownerTeamId": "payments-team",
  "layerId": "runtime",
  "tier": "tier1",
  "metricIds": ["slo-payments-availability", "slo-payments-latency"]
}
```

### Connecting to SLOs and Maturity

The organizational model connects to SLOs and maturity roadmaps:

1. **Metrics** belong to a domain, layer, and optionally a service
2. **SLOs** are defined on metrics with machine-evaluable targets
3. **Goals** aggregate SLO requirements into maturity levels
4. **Teams** own services and are accountable for their SLOs
5. **Phases** organize goal progression over time (quarters)

```
Team owns → Service has → Metrics with → SLOs required by → Goals tracked in → Phases
```

## Installation

```bash
go install github.com/grokify/prism/cmd/prism@latest
```

Or add as a library dependency:

```bash
go get github.com/grokify/prism
```

## CLI Usage

### Initialize a new PRISM document

```bash
# Create default document with operations metrics
prism init

# Create operations-focused document
prism init -d operations -o ops.json
```


### Validate a document

```bash
prism validate prism.json
```

### Calculate PRISM score

```bash
# Basic score
prism score prism.json

# Detailed breakdown
prism score prism.json --detailed

# JSON output
prism score prism.json --json
```

### List available constants

```bash
prism catalog
```

### Goal commands (v0.2.0)

```bash
# List all goals
prism goal list prism.json

# Show goal details
prism goal show goal-reliability prism.json

# Show goal progress with SLO compliance
prism goal progress goal-reliability prism.json
```

### Phase commands (v0.2.0)

```bash
# List all phases
prism phase list prism.json

# Show phase details
prism phase show phase-q1-2026 prism.json

# Show enter/exit metrics for a phase
prism phase metrics phase-q1-2026 prism.json
```

### Roadmap commands (v0.2.0)

```bash
# Show roadmap overview
prism roadmap show prism.json

# Show progress across all phases and goals
prism roadmap progress prism.json
```

### Initiative commands (v0.3.0)

```bash
# List all initiatives by status
prism initiative list prism.json

# List initiatives by phase
prism initiative list prism.json --by-phase

# List initiatives by goal
prism initiative list prism.json --by-goal

# Show initiative details
prism initiative show prism.json init-monitoring
```

### Report commands (v0.2.0)

```bash
# Generate markdown report (both views)
prism report prism.json -o report.md

# Generate phase-centric view only
prism report prism.json --view by-phase

# Generate JSON report
prism report prism.json --format json

# Generate SLO compliance report
prism slo-report prism.json
```

### Dashboard commands (v0.2.0)

```bash
# Generate executive dashboard
prism dashboard prism.json -o dashboard.json

# Convert to dashforge format
prism dashforge prism.json -o dashforge.json
```

### Layer commands (v0.3.0)

```bash
# List all layers
prism layer list prism.json

# Show layer details with metrics
prism layer show prism.json runtime
```

### Team commands (v0.3.0)

```bash
# List all teams grouped by type
prism team list prism.json

# Show team details with services
prism team show prism.json payments-team
```

### Service commands (v0.3.0)

```bash
# List all services grouped by layer
prism service list prism.json

# Show service details with metrics
prism service show prism.json payments-api
```

### Analyze command (v0.3.0)

```bash
# Analyze document and show gaps
prism analyze prism.json

# Output as JSON for automation
prism analyze prism.json -f json

# Generate LLM prompt for initiative recommendations
prism analyze prism.json -f prompt
```

### Export commands (v0.3.0)

```bash
# Export as OKR document
prism export okr prism.json -o roadmap.okr.json

# Export as V2MOM document
prism export v2mom prism.json -o roadmap.v2mom.json
```

### Maturity commands (v0.4.0)

```bash
# Generate Markdown report with SLI Catalog
prism maturity report maturity-spec.json -o report.md

# Generate domain view only
prism maturity report maturity-spec.json --view domain

# Generate framework view only
prism maturity report maturity-spec.json --view framework

# Generate Excel workbook
prism maturity xlsx maturity-spec.json -o report.xlsx
```

## Schema Overview

### Domains

PRISM organizes metrics into three primary domains:

| Domain | Description |
|--------|-------------|
| `operations` | Reliability, performance, and efficiency metrics |
| `security` | Application and infrastructure security metrics |
| `quality` | Code quality, testing, and defect management metrics |

### Layers

Metrics are classified by value stream layer:

| Layer | Description |
|-------|-------------|
| `requirements` | Product ideation, specifications, and design |
| `code` | Application code, libraries, and dependencies |
| `infra` | Cloud resources, networking, and platform services |
| `runtime` | Running services, containers, and workloads |
| `adoption` | Product analytics, user engagement, and self-service |
| `support` | Customer support, incident management, and escalations |

### Lifecycle Stages

Metrics are mapped to software delivery lifecycle stages:

| Stage | Description |
|-------|-------------|
| `design` | Architecture, requirements, planning |
| `build` | CI/CD, code quality, dependency management |
| `test` | Testing coverage, quality assurance |
| `runtime` | Production monitoring, availability, performance |
| `response` | Incident response, remediation, recovery |

### Categories

| Category | Description |
|----------|-------------|
| `prevention` | Proactive controls that prevent issues |
| `detection` | Monitoring and alerting capabilities |
| `response` | Incident handling and remediation |
| `reliability` | Availability and durability |
| `efficiency` | Performance and resource utilization |
| `quality` | Code and process quality |

### Metric Types

| Type | Description | Example |
|------|-------------|---------|
| `coverage` | Percentage of coverage | Test coverage |
| `rate` | Frequency or percentage | Error rate |
| `latency` | Time duration | P99 latency, MTTR |
| `ratio` | Proportion | Success ratio |
| `count` | Absolute count | Deployment count |
| `distribution` | Statistical distribution | Latency percentiles |
| `score` | Composite score | Health score |

## Example Metric

```json
{
  "id": "ops-availability",
  "name": "Service Availability",
  "description": "Percentage of time the service is available",
  "domain": "operations",
  "stage": "runtime",
  "category": "reliability",
  "layer": "runtime",
  "serviceId": "payments-api",
  "metricType": "rate",
  "trendDirection": "higher_better",
  "unit": "%",
  "baseline": 99.0,
  "current": 99.95,
  "target": 99.99,
  "thresholds": {
    "green": 99.95,
    "yellow": 99.9,
    "red": 99.0
  },
  "slo": {
    "target": ">=99.99%",
    "operator": "gte",
    "value": 99.99,
    "window": "30d"
  },
  "frameworkMappings": [
    {"framework": "SRE", "reference": "availability-slo"},
    {"framework": "DORA", "reference": "availability"}
  ]
}
```

## PRISM Score Calculation

The PRISM score combines maturity levels, metric performance, and customer awareness into a composite health score (0.0-1.0).

### Formula

```
CellScore = (MaturityWeight × MaturityScore) + (PerformanceWeight × PerformanceScore)
BaseScore = Σ(CellScore × Weight) / Σ(Weight)
Overall = BaseScore × AwarenessScore
```

### Default Weights

**Component weights:**

- Maturity: 40%
- Performance: 60%

**Stage weights:**

- Design: 15%
- Build: 20%
- Test: 15%
- Runtime: 30%
- Response: 20%

**Domain weights:**

- Security: 50%
- Operations: 50%

### Score Interpretation

| Score | Level | Description |
|-------|-------|-------------|
| ≥0.90 | Elite | Industry-leading practices |
| ≥0.75 | Strong | Well-managed, proactive |
| ≥0.50 | Medium | Adequate, room for improvement |
| ≥0.25 | Weak | Significant gaps |
| <0.25 | Critical | Immediate attention required |

## Maturity Levels

PRISM uses a 5-level maturity model:

| Level | Name | Description |
|-------|------|-------------|
| 1 | Reactive | Ad-hoc processes, firefighting mode |
| 2 | Basic | Basic controls, some documentation |
| 3 | Defined | Standardized processes, consistent execution |
| 4 | Managed | Data-driven, measured and controlled |
| 5 | Optimizing | Continuous improvement, automated optimization |

## Framework Mappings

PRISM metrics can be mapped to external frameworks:

| Framework | Constant | Description |
|-----------|----------|-------------|
| NIST CSF 1.1 | `NIST_CSF` | NIST Cybersecurity Framework 1.1 |
| NIST CSF 2.0 | `NIST_CSF_2` | NIST Cybersecurity Framework 2.0 |
| NIST 800-53 | `NIST_800_53` | Security and Privacy Controls |
| NIST RMF | `NIST_RMF` | Risk Management Framework |
| NIST AI RMF | `NIST_AI_RMF` | AI Risk Management Framework |
| FedRAMP High | `FEDRAMP_HIGH` | FedRAMP High baseline |
| FedRAMP Moderate | `FEDRAMP_MOD` | FedRAMP Moderate baseline |
| FedRAMP Low | `FEDRAMP_LOW` | FedRAMP Low baseline |
| MITRE ATT&CK | `MITRE_ATTACK` | Threat framework |
| CIS Controls | `CIS_CONTROLS` | Critical Security Controls |
| SOC 2 | `SOC_2` | Trust Services Criteria |
| ISO 27001 | `ISO_27001` | Information Security Management |
| DORA | `DORA` | DevOps Research and Assessment |
| SRE | `SRE` | Site Reliability Engineering |

## JSON Schema

The JSON Schema is auto-generated from Go types:

```bash
cd schema && go run generate.go
```

Schema location: `schema/prism.schema.json`

Use in your editor for validation:

```json
{
  "$schema": "https://github.com/grokify/prism/schema/prism.schema.json",
  "metrics": [...]
}
```

## Goal-Driven Maturity Roadmap

PRISM supports goal-driven maturity tracking with multi-phase roadmaps.

### Goals

Goals represent strategic objectives with their own 5-level maturity models:

```json
{
  "id": "goal-reliability",
  "name": "Achieve High Reliability",
  "owner": "VP Engineering",
  "currentLevel": 3,
  "targetLevel": 4,
  "maturityModel": {
    "levels": [
      {
        "level": 3,
        "name": "Defined",
        "requiredSLOs": [
          { "metricId": "ops-availability" },
          { "metricId": "ops-mttr" }
        ],
        "metricCriteria": [
          { "metricId": "ops-availability", "operator": "gte", "value": 99.5 }
        ]
      }
    ]
  }
}
```

Each maturity level specifies which SLOs must be met to achieve that level.

### Phases

Phases organize work into time-bounded periods (quarters) with goal targets:

```json
{
  "id": "phase-q1-2026",
  "name": "Q1 2026",
  "quarter": "Q1",
  "year": 2026,
  "startDate": "2026-01-01",
  "endDate": "2026-03-31",
  "goalTargets": [
    { "goalId": "goal-reliability", "enterLevel": 2, "exitLevel": 3 }
  ],
  "swimlanes": [
    {
      "name": "Platform Initiatives",
      "domain": "operations",
      "initiativeIds": ["init-monitoring", "init-ci-cd"]
    }
  ]
}
```

### Initiatives

Initiatives link to goals and phases with deployment tracking:

```json
{
  "id": "init-monitoring",
  "name": "Observability Platform",
  "goalIds": ["goal-reliability"],
  "phaseId": "phase-q1-2026",
  "status": "completed",
  "deploymentStatus": {
    "status": "completed",
    "totalCustomers": 50,
    "deployedCustomers": 50,
    "adoptionPercent": 100
  }
}
```

## Integration with Structured-Plan

PRISM integrates with [structured-plan](https://github.com/grokify/structured-plan) to provide a complete operational planning workflow. PRISM serves as the source of truth for requirements (maturity models, SLOs), while structured-plan handles execution tracking (OKRs, roadmaps).

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         PRISM (Source of Truth)                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │
│  │    Goals     │  │   Maturity   │  │     SLOs     │           │
│  │              │  │    Models    │  │              │           │
│  │ "Reliability"│  │  M1→M2→M3→M4 │  │ avail>=99.9% │           │
│  └──────────────┘  └──────────────┘  └──────────────┘           │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      LLM Analysis Layer                         │
│                                                                 │
│  Input: PRISM requirements + context (team capacity, etc.)      │
│  Output: Suggested initiatives, phase mapping, dependencies     │
│                                                                 │
│  "To reach M3 reliability by Q2, you need:                      │
│   - Q1: Observability platform (enables SLO measurement)        │
│   - Q1: Incident runbooks (reduces MTTR)                        │
│   - Q2: SLO dashboards (tracks compliance)"                     │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                 Structured-Plan (Execution)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │
│  │     OKR      │  │   V2MOM      │  │   Roadmap    │           │
│  │              │  │              │  │              │           │
│  │ Objectives   │  │ Methods      │  │ Phases       │           │
│  │ Key Results  │  │ Measures     │  │ Deliverables │           │
│  └──────────────┘  └──────────────┘  └──────────────┘           │
└─────────────────────────────────────────────────────────────────┘
```

### Data Flow

| PRISM Concept | Structured-Plan Concept |
|---------------|-------------------------|
| Goal | OKR Objective |
| Goal.TargetLevel | Objective Target |
| SLO (per maturity level) | Key Result |
| Phase.GoalTargets | PhaseTargets in Key Results |
| Initiative | Deliverable (in roadmap phase) |

### Workflow

1. **Define requirements in PRISM** - Goals, maturity models, SLOs
2. **Analyze with LLM** - Generate initiative recommendations to achieve targets
3. **Export to structured-plan** - OKR/V2MOM/Roadmap format
4. **Track execution** - Monitor progress against phase targets

```bash
# Analyze PRISM document and suggest initiatives
prism analyze prism.json

# Export as OKR document for structured-plan
prism export okr prism.json -o roadmap.okr.json

# Export as V2MOM document
prism export v2mom prism.json -o roadmap.v2mom.json
```

## Examples

See the `examples/` directory for PRISM document types:

### PRISM Documents (`examples/prism-documents/`)

Standard PRISM documents for metrics, goals, and roadmaps:

- `operations-metrics.json` - Operations-focused metrics (DORA metrics, SLOs, reliability)
- `operations-layers.json` - Layer-based metric organization with golden signals
- `team-topology.json` - Full team topology with services and ownership
- `quality-metrics.json` - Quality domain with ISO 25010 verticals
- `goal-roadmap.json` - Goal-driven maturity roadmap with phases and initiatives

### Maturity Models (`examples/maturity-models/`)

PRISM Maturity Models define criteria for each maturity level (M1-M5):

- `operations/model.json` - Reliability, deployment, monitoring maturity
- `security/model.json` - Prevention, detection, response maturity
- `organization/model.json` - Multi-domain organizational maturity

### Maturity State (`examples/maturity-state/`)

PRISM Maturity State documents track current state with temporal windows:

- `operations/state-q2-2026.json` - Q2 2026 operations state
- `organization/state-q1-2026.json` - Q1 2026 multi-domain state
- `security/state-q2-2026.json` - Q2 2026 security state


## Library Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/grokify/prism"
)

func main() {
    // Load document
    data, _ := os.ReadFile("prism.json")
    var doc prism.PRISMDocument
    json.Unmarshal(data, &doc)

    // Validate
    if errs := doc.Validate(); errs.HasErrors() {
        fmt.Println("Validation errors:", errs)
        return
    }

    // Calculate score
    score := doc.CalculatePRISMScore(nil, nil)
    fmt.Printf("PRISM Score: %.1f%% (%s)\n", score.Overall*100, score.Interpretation)

    // Check individual metrics
    for _, m := range doc.Metrics {
        status := m.CalculateStatus()
        meetsSLO := m.MeetsSLO()
        fmt.Printf("  %s: %s (SLO met: %v)\n", m.Name, status, meetsSLO)
    }
}
```

## License

MIT
