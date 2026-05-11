# PRD: PRISM Schema Restructure (v0.6.0)

**Status:** Planned
**Target Version:** v0.6.0
**Author:** PRISM Team
**Created:** 2026-05-10

## Overview

Restructure PRISM schemas into a clear three-part taxonomy: Model, State, and Plan. This creates a unified naming convention and clarifies the purpose of each document type.

## Problem Statement

Current schema naming is inconsistent:

- `maturity.schema.json` - Maturity model definitions
- `prism.schema.json` - Mix of metrics, goals, phases, initiatives

The relationship between document types is unclear, and the naming doesn't reflect the Model → State → Plan workflow.

## Proposed Solution

### Schema Taxonomy

| Schema | Purpose | Focus | Question Answered |
|--------|---------|-------|-------------------|
| `prism-maturity-model.schema.json` | Definitions | What does good look like? | "M4 availability means ≥99.9%" |
| `prism-maturity-state.schema.json` | Measurement | Where are we now? | "Current availability is 99.7%" |
| `prism-maturity-plan.schema.json` | Execution | How do we get there? | "Reach M4 by Q4 via these initiatives" |

### Document Relationships

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

### Schema Contents

#### prism-maturity-model.schema.json

Defines what maturity levels mean. No current state.

```json
{
  "$schema": "prism-maturity-model.schema.json",
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
        { "level": 1, "criteria": [...] },
        { "level": 2, "criteria": [...] }
      ]
    }
  }
}
```

#### prism-maturity-state.schema.json

Tracks current state with temporal windows. References a model.

```json
{
  "$schema": "prism-maturity-state.schema.json",
  "metadata": {
    "name": "Operations State Q2 2026",
    "maturityModelRef": "./model.json"
  },
  "sloWindows": ["7d", "30d", "90d", "quarterly"],
  "sliState": {
    "sli-availability": {
      "windows": {
        "30d": { "value": 99.7, "timestamp": "2026-05-10T00:00:00Z" }
      },
      "history": [...]
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

#### prism-maturity-plan.schema.json

Defines improvement roadmap. References state (and transitively, model).

```json
{
  "$schema": "prism-maturity-plan.schema.json",
  "metadata": {
    "name": "Operations Improvement Plan 2026",
    "maturityStateRef": "./state-q2-2026.json"
  },
  "goals": [
    {
      "id": "goal-reliability",
      "name": "Achieve High Reliability",
      "targetLevel": 4,
      "targetDate": "2026-12-31"
    }
  ],
  "phases": [
    {
      "id": "phase-q2-2026",
      "name": "Q2 2026",
      "goalTargets": [...]
    }
  ],
  "initiatives": [
    {
      "id": "init-observability",
      "name": "Observability Platform",
      "phaseId": "phase-q2-2026",
      "goalIds": ["goal-reliability"]
    }
  ],
  "teams": [...],
  "services": [...]
}
```

## Examples Directory Restructure

### Current Structure (by type)

```
examples/
├── prism-documents/
│   ├── operations-metrics.json
│   ├── operations-layers.json
│   └── goal-roadmap.json
├── maturity-models/
│   ├── operations/model.json
│   └── security/model.json
└── maturity-state/
    ├── operations/state-q2-2026.json
    └── security/state-q2-2026.json
```

### Proposed Structure (by domain)

```
examples/
├── operations/
│   ├── model.json           # prism-maturity-model
│   ├── state-q2-2026.json   # prism-maturity-state
│   └── plan-2026.json       # prism-maturity-plan
├── security/
│   ├── model.json
│   ├── state-q2-2026.json
│   └── plan-2026.json
├── organization/
│   ├── model.json
│   └── state-q1-2026.json
└── README.md
```

## Migration Path

### Schema File Changes

| Current | New |
|---------|-----|
| `maturity.schema.json` | `prism-maturity-model.schema.json` |
| `prism.schema.json` | `prism-maturity-plan.schema.json` |
| (new) | `prism-maturity-state.schema.json` |

### Code Changes

1. **Schema generation** - Update `schema/generate*.go` files
2. **Schema embedding** - Update `schema/embed.go` with new file names
3. **CLI commands** - Update `$schema` references in scaffolds
4. **Validation** - Update schema validation logic
5. **Examples** - Restructure and update `$schema` references

### Breaking Changes

- `$schema` URLs change for all document types
- Example file paths change
- CLI output may reference new schema names

## Implementation Phases

### Phase 1: Schema Files

1. Create `prism-maturity-model.schema.json` (rename from `maturity.schema.json`)
2. Create `prism-maturity-state.schema.json` (formalize existing state format)
3. Create `prism-maturity-plan.schema.json` (rename from `prism.schema.json`)
4. Update `schema/embed.go` with new embeddings

### Phase 2: Go Types

1. Review and update Go types to match new schema organization
2. Add `MaturityStateRef` and `MaturityModelRef` fields
3. Update validation functions

### Phase 3: Examples

1. Restructure `examples/` to domain-centric layout
2. Update all `$schema` references
3. Create missing state and plan files for each domain

### Phase 4: CLI

1. Update `prism init` to generate new schema references
2. Update `prism validate` to recognize all three schema types
3. Consider new commands: `prism model`, `prism state`, `prism plan`

### Phase 5: Documentation

1. Update README with new schema taxonomy
2. Update MkDocs site with new structure
3. Add migration guide for v0.5.0 → v0.6.0

## Success Criteria

1. All three schemas are formally defined and generated from Go types
2. Examples demonstrate the Model → State → Plan workflow
3. CLI commands work with all three document types
4. Documentation clearly explains the taxonomy
5. Migration guide helps users upgrade

## Timeline

| Phase | Description | Estimate |
|-------|-------------|----------|
| 1 | Schema files | - |
| 2 | Go types | - |
| 3 | Examples | - |
| 4 | CLI | - |
| 5 | Documentation | - |

## Open Questions

1. Should Plan reference State directly, or should it reference Model and State separately?
2. Should we support inline state in Plan documents for simpler use cases?
3. How should CLI commands be organized? (`prism model`, `prism state`, `prism plan` vs. current structure)
4. Should we add a `prism migrate` command to help users upgrade?

## References

- [v0.5.0 Release Notes](../../releases/v0.5.0.md) - Model/State separation
- [REFACTOR_MATURITY_STATE.md](../REFACTOR_MATURITY_STATE.md) - Initial separation design
- [Core PRD](../core/PRD.md) - Overall PRISM vision
