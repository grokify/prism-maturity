# PRISM Examples

This directory contains example documents demonstrating PRISM's three-part document taxonomy.

## Document Types

PRISM uses three document types that form a Model → State → Plan workflow:

| Schema | Purpose | Question Answered |
|--------|---------|-------------------|
| `prism-maturity-model` | Definitions | What does good look like? |
| `prism-maturity-state` | Measurement | Where are we now? |
| `prism-maturity-plan` | Execution | How do we get there? |

## Directory Structure

Examples are organized by domain:

```
examples/
├── operations/
│   ├── model.json           # Maturity Model - level definitions
│   ├── state-q2-2026.json   # Maturity State - current tracking
│   ├── plan.json            # Maturity Plan - goals and roadmap
│   └── plan-layers.json     # Alternative plan with layer organization
├── security/
│   ├── model.json
│   └── state-q2-2026.json
├── organization/
│   ├── model.json
│   └── state-q1-2026.json
├── quality/
│   └── plan.json            # Quality domain metrics and goals
├── prism-documents/         # Additional plan examples
│   ├── operations-metrics.json
│   ├── operations-layers.json
│   ├── goal-roadmap.json
│   ├── project-scores.json
│   ├── quality-metrics.json
│   └── team-topology.json
└── README.md
```

## Document Relationships

```
┌─────────────────────────────────────────────────────────────────┐
│                    prism-maturity-model                         │
│              "What does M4 availability mean?"                  │
│         SLIs, domains, levels, criteria, enablers               │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ references (maturityModelRef)
┌─────────────────────────────┴───────────────────────────────────┐
│                    prism-maturity-state                         │
│              "We're at M3, availability is 99.7%"               │
│       Current values, temporal windows, history, gaps           │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ references (maturityStateRef)
┌─────────────────────────────┴───────────────────────────────────┐
│                    prism-maturity-plan                          │
│              "Reach M4 by Q4 via these initiatives"             │
│         Goals, phases, initiatives, teams, roadmap              │
└─────────────────────────────────────────────────────────────────┘
```

## 1. PRISM Maturity Model

**Schema:** `prism-maturity-model.schema.json`

Maturity Models define the **criteria for each maturity level (M1-M5)**. They are reference documents that describe what "good" looks like at each level. They contain no current state.

**Example:** `operations/model.json`

```json
{
  "$schema": "https://github.com/grokify/prism/schema/prism-maturity-model.schema.json",
  "metadata": { "name": "Operations Maturity Model" },
  "slis": {
    "sli-availability": {
      "name": "Service Availability",
      "unit": "%",
      "sliType": "availability",
      "measurementType": "quantitative"
    }
  },
  "domains": {
    "reliability": {
      "levels": [
        { "level": 1, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 95}] },
        { "level": 2, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99}] },
        { "level": 3, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99.5}] },
        { "level": 4, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99.9}] },
        { "level": 5, "criteria": [{"sliId": "sli-availability", "operator": ">=", "target": 99.99}] }
      ]
    }
  }
}
```

## 2. PRISM Maturity State

**Schema:** `prism-maturity-state.schema.json`

Maturity State documents track the **current state** of your system with temporal windows and historical values.

**Example:** `operations/state-q2-2026.json`

```json
{
  "$schema": "https://github.com/grokify/prism/schema/prism-maturity-state.schema.json",
  "metadata": {
    "name": "Operations State Q2 2026",
    "maturityModelRef": "./model.json"
  },
  "sloWindows": ["7d", "30d", "90d", "quarterly"],
  "sliState": {
    "sli-availability": {
      "qualitativeState": "alerting",
      "windows": {
        "7d":  { "value": 99.7, "timestamp": "2026-05-09T00:00:00Z" },
        "30d": { "value": 99.5, "timestamp": "2026-05-09T00:00:00Z" }
      },
      "targets": {
        "Q2_2026": { "value": 99.5, "maturityLevel": 3 },
        "Q4_2026": { "value": 99.9, "maturityLevel": 4 }
      }
    }
  },
  "maturityState": {
    "reliability": {
      "currentLevel": 3,
      "targetLevel": 4
    }
  }
}
```

## 3. PRISM Maturity Plan

**Schema:** `prism-maturity-plan.schema.json`

Maturity Plan documents define **goals, phases, and initiatives** for achieving maturity targets.

**Example:** `operations/plan.json`

```json
{
  "$schema": "https://github.com/grokify/prism/schema/prism-maturity-plan.schema.json",
  "metadata": {
    "name": "Operations Improvement Plan 2026"
  },
  "goals": [
    {
      "id": "goal-reliability",
      "name": "Achieve High Reliability",
      "domain": "operations",
      "targetLevel": 4,
      "currentLevel": 3
    }
  ],
  "phases": [
    {
      "id": "phase-q2-2026",
      "name": "Q2 2026",
      "startDate": "2026-04-01",
      "endDate": "2026-06-30"
    }
  ],
  "initiatives": [
    {
      "id": "init-observability",
      "name": "Observability Platform",
      "phaseId": "phase-q2-2026"
    }
  ]
}
```

## Temporal Windows

PRISM Maturity State supports multiple SLO windows for tracking:

| Window | Description |
|--------|-------------|
| `7d` | Rolling 7-day window |
| `30d` | Rolling 30-day window |
| `90d` | Rolling 90-day window |
| `quarterly` | Calendar quarter |
| `annual` | Calendar year |

## Measurement Types

SLIs can be classified by measurement type:

| Type | Description | Example |
|------|-------------|---------|
| `quantitative` | Numeric values only | Availability: 99.9% |
| `qualitative` | State-based only | Monitoring: "tracked" |
| `hybrid` | Both numeric and state | Start with "tracked", progress to 99.9% |

## Qualitative States

Standard qualitative state progression:

| State | Order | Description |
|-------|-------|-------------|
| `none` | 0 | Not started |
| `adhoc` | 1 | Ad-hoc/informal |
| `tracked` | 2 | Being tracked |
| `measured` | 3 | Measured with targets |
| `alerting` | 4 | Alerting enabled |
| `optimized` | 5 | Continuously optimized |

## Example Domains

| Domain | Files | Description |
|--------|-------|-------------|
| Operations | `operations/` | Reliability, deployment, monitoring |
| Security | `security/` | Prevention, detection, response |
| Organization | `organization/` | Multi-domain organizational maturity |
| Quality | `quality/` | Code quality, testing metrics |

## CLI Commands

```bash
# Validate a maturity model
prism maturity report operations/model.json

# Generate XLSX report
prism maturity xlsx operations/model.json -o report.xlsx

# Validate a plan document
prism validate prism-documents/operations-metrics.json
```
