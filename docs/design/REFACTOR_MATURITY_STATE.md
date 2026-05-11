# Refactor: Separate Maturity Model from State

## Overview

Split the current maturity model into two document types:

1. **PRISM Maturity Model** - Pure reference (levels, criteria, thresholds, SLI definitions)
2. **PRISM Maturity State** - State tracking (past, present, future values and targets)

## Document Relationship

```
┌─────────────────────────┐         ┌─────────────────────────┐
│  PRISM Maturity Model   │◄────────│  PRISM Maturity State   │
│  (Reference/Definition) │         │  (Tracking/Assessment)  │
├─────────────────────────┤         ├─────────────────────────┤
│ • SLI definitions       │         │ • Current SLI values    │
│ • Qualitative states    │         │ • Historical values     │
│ • Domain levels (M1-M5) │         │ • Future targets        │
│ • Criteria & thresholds │         │ • Maturity level state  │
│ • Enabler definitions   │         │ • Enabler progress      │
│ • Framework mappings    │         │ • Temporal windows      │
└─────────────────────────┘         └─────────────────────────┘
```

## Current State Analysis

### Functionality Using `maturity.Assessments`

| File | Usage | Migration Path |
|------|-------|----------------|
| `dashboard/dashboard.go` | `currentLevel`, `targetLevel`, `criteriaValues` | Read from Maturity State |
| `maturity/xlsx.go` | Current values for spreadsheet | Read from Maturity State |
| `maturity/marp.go` | Values for presentations | Read from Maturity State |
| `maturity/xlsx_test.go` | Test data | Update test fixtures |

### Current `DomainAssessment` Fields

```go
type DomainAssessment struct {
    Domain         string             // → Move to Maturity State
    AssessedAt     string             // → Move to Maturity State
    AssessedBy     string             // → Move to Maturity State
    CurrentLevel   int                // → Move to Maturity State
    TargetLevel    int                // → Move to Maturity State
    CriteriaValues map[string]float64 // → Move to Maturity State (with windows)
    CriteriaStatus map[string]string  // → Move to Maturity State (qualitative)
    EnablerStatus  map[string]string  // → Move to Maturity State
}
```

## New Schema Design

### 1. PRISM Maturity Model (Reference Only)

The Maturity Model defines **what** is measured and **how** levels are achieved. It contains no current state.

```json
{
  "$schema": "prism-maturity-model.schema.json",
  "metadata": { "name": "Operations Maturity", "version": "1.0.0" },

  "slis": {
    "sli-availability": {
      "name": "Service Availability",
      "unit": "%",
      "sliType": "availability",
      "measurementType": "hybrid",
      "qualitativeStates": [
        { "id": "none", "label": "Not tracked", "order": 0 },
        { "id": "adhoc", "label": "Ad-hoc", "order": 1 },
        { "id": "tracked", "label": "Tracked", "order": 2 },
        { "id": "measured", "label": "Measured with SLO", "order": 3 },
        { "id": "alerting", "label": "SLO + Alerting", "order": 4 },
        { "id": "optimized", "label": "Optimized", "order": 5 }
      ]
    }
  },

  "domains": {
    "reliability": {
      "levels": [
        {
          "level": 1,
          "criteria": [{
            "sliId": "sli-availability",
            "type": "qualitative",
            "operator": ">=",
            "target": "tracked"
          }]
        },
        {
          "level": 2,
          "criteria": [{
            "sliId": "sli-availability",
            "type": "quantitative",
            "operator": ">=",
            "target": 95
          }]
        }
      ]
    }
  }
}
```

### 2. PRISM Maturity State (State + Planning)

The Maturity State tracks **where** you are, **where** you've been, and **where** you're going.

```json
{
  "$schema": "prism-maturity-state.schema.json",
  "metadata": {
    "name": "Operations Status Q2 2026",
    "maturityModelRef": "examples/maturity-models/operations/model.json"
  },

  "sloWindows": ["7d", "30d", "90d", "quarterly", "annual"],

  "sliState": {
    "sli-availability": {
      "qualitativeState": "measured",
      "windows": {
        "7d":   { "value": 99.97, "timestamp": "2026-05-10T00:00:00Z" },
        "30d":  { "value": 99.92, "timestamp": "2026-05-10T00:00:00Z" },
        "90d":  { "value": 99.88, "timestamp": "2026-05-10T00:00:00Z" }
      },
      "history": [
        { "window": "30d", "value": 99.88, "timestamp": "2026-04-01T00:00:00Z" },
        { "window": "30d", "value": 99.90, "timestamp": "2026-04-15T00:00:00Z" },
        { "window": "30d", "value": 99.92, "timestamp": "2026-05-01T00:00:00Z" }
      ],
      "targets": {
        "Q2_2026": { "value": 99.9 },
        "Q4_2026": { "value": 99.95 }
      }
    }
  },

  "maturityState": {
    "reliability": {
      "current": { "level": 3, "achievedAt": "2026-03-15" },
      "target": { "level": 4, "targetDate": "2026-06-30" },
      "history": [
        { "level": 2, "achievedAt": "2025-09-01" },
        { "level": 3, "achievedAt": "2026-03-15" }
      ]
    }
  },

  "enablerState": {
    "enabler-monitoring": { "status": "completed", "completedAt": "2026-02-01" },
    "enabler-alerting": { "status": "in_progress", "progress": 60 }
  },

  "goals": [...],
  "phases": [...],
  "initiatives": [...]
}
```

## File Organization

```
examples/
├── maturity-models/           # PRISM Maturity Model files
│   ├── operations/
│   │   └── model.json
│   └── organization/
│       └── model.json
│
└── maturity-state/            # PRISM Maturity State files
    ├── operations/
    │   ├── state-q1-2026.json
    │   └── state-q2-2026.json
    └── organization/
        └── state-current.json
```

## Migration Plan

### Phase 1: Add New Types (No Breaking Changes)

1. Add `SLIState` type to `prism.go`
2. Add `MaturityState` type to `prism.go`
3. Add `QualitativeState` to maturity model
4. Update schemas (additive only)

### Phase 2: Update Generators (Dual Support)

1. Dashboard reads from:
   - New: `PRISMDocument.SLIState` / `PRISMDocument.MaturityState`
   - Fallback: `maturity.Spec.Assessments` (deprecated)

2. XLSX generator supports both sources

### Phase 3: Update Examples

1. Remove `assessments` from maturity model examples
2. Create PRISM Maturity State examples

### Phase 4: Deprecation (User Approval Required)

1. Mark `maturity.Spec.Assessments` as deprecated
2. Add migration CLI command
3. Remove after confirmation

## Potential Functionality Loss

| Functionality | Risk | Mitigation |
|---------------|------|------------|
| Dashboard maturity display | LOW | Dual-read support |
| XLSX generation | LOW | Dual-read support |
| Marp presentation | LOW | Dual-read support |
| Existing JSON files | NONE | Old format still works |

## New Capabilities

1. **Temporal tracking**: 7d, 30d, 90d, quarterly, annual windows
2. **History**: Track progression over time
3. **Qualitative states**: Track "tracked", "measured", etc.
4. **Hybrid criteria**: Mix qualitative + quantitative
5. **Target planning**: Future state targets by date
6. **Model reference**: Link Maturity State to its Model

## Naming Summary

| Term | Description |
|------|-------------|
| **PRISM Maturity Model** | Reference document defining SLIs, levels, criteria, thresholds |
| **PRISM Maturity State** | State document tracking current values, history, and targets |
| **SLI** | Service Level Indicator - what is measured |
| **SLO** | Service Level Objective - target for an SLI |
| **Qualitative State** | Non-numeric state (e.g., "tracked", "measured", "alerting") |
| **Temporal Window** | Time period for SLO measurement (7d, 30d, 90d, quarterly, annual) |
