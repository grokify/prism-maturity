# PRISM Project Review

**Date**: 2026-04-07
**Purpose**: Evaluate PRISM for building a maturity model-based roadmap

## Summary

PRISM (Platform for Reliability, Improvement, and Strategic Maturity) is a Go library + CLI for COO-level organizational health monitoring. It combines SLOs, maturity modeling, and OKRs into a single JSON document format for tracking maturity across operations, security, quality, product, and AI domains.

The goal-driven maturity roadmap feature (FEAT_MATURITYROADMAP) has been implemented at the data model level but the CLI commands for it are missing.

## Current State

### What's Working

- **Core types**: `PRISMDocument`, `Metric`, `SLO`, `MaturityModel`, `Goal`, `Phase`, `PhaseMetrics` — all implemented
- **Validation**: Comprehensive cross-reference validation (metric IDs, goal IDs, phase IDs, initiative linkage)
- **Scoring**: Composite PRISM score with domain/stage cell breakdown, awareness multiplier
- **Goal maturity calculation**: `Goal.CurrentMaturityLevel()` checks SLOs + metric criteria from level 5 down
- **Phase metrics**: `CalculateGoalProgress()`, `CalculateInitiativeMetrics()`, `CalculateSLOCompliance()` all work
- **Tests**: All pass (55.9% coverage), all example files validate
- **CLI**: `init`, `validate`, `score`, `catalog` commands work

### What's Missing (from the PRD plan)

#### 1. CLI Commands for Roadmap Features

The PRD specifies `prism goal list/show/progress`, `prism phase list/show/metrics`, and `prism roadmap show/progress`. None of these exist. The `cmd/prism/` directory only has `init.go`, `validate.go`, `score.go`, `catalog.go`.

#### 2. Score Doesn't Incorporate Goals

The `--goals` and `--phase` flags mentioned in the PRD aren't implemented. The current score is 27.5% for the example because all maturity cells are 0 (the global `MaturityModel` has no cells in the goal-roadmap example). Goal-level maturity is tracked separately but doesn't feed into the composite score.

#### 3. The Score Is Misleading

The example shows "Weak" at 27.5% despite metrics performing well (SAST at 95%, availability at 99.95%). Root causes:

- Maturity weight is 40% but all maturity cells are 0.0 (no global maturity model populated)
- Empty domain/stage cells (e.g., security/runtime, operations/design) score 0.0 and drag down the average
- Goal-level maturity (both goals at level 3) isn't reflected in the composite score

#### 4. No Visual/Export Output

The PRD defers visual rendering to uiforge, but there's no markdown, HTML, or slide export for roadmap views.

## Strengths for Maturity-Based Roadmap

- The **Goal → MaturityModel → RequiredSLOs → MetricCriteria** chain is exactly what's needed. Each goal has its own 5-level model with concrete, measurable criteria.
- **Phase/quarter organization** with enter/exit maturity targets maps directly to roadmap planning.
- **Initiative tracking** with deployment status (customer adoption %) is useful for B2B rollout tracking.
- The **JSON schema** is clean and well-documented.

## Key Issues to Address

### 1. Scoring Gap

The composite score needs to incorporate goal maturity, not just the global maturity model. Right now you can have goals at level 3-4 and still get a "Weak" score because the global maturity cells are empty. Options:

- Auto-populate global maturity cells from goal maturity levels
- Add a goal-weighted score path (the `--goals` flag from the PRD)

### 2. Sparse Metric Coverage Penalizes Unfairly

If you only have security metrics (no operations/design, operations/test, etc.), those empty cells score 0.0 and tank the overall score. Options:

- Only score cells that have metrics (skip empty cells)
- Explicitly define which domains/stages are in scope

### 3. Missing CLI for Roadmap Workflows

Need to implement the `goal`, `phase`, and `roadmap` subcommands to make this usable beyond programmatic Go usage. Estimated ~200-300 lines each based on the existing `score.go` pattern.

### 4. Test Coverage

55.9% for the core library is decent, but the CLI has 0% coverage and validation/scoring edge cases could use more tests.

## Score Breakdown (examples/goal-roadmap.json)

```
Overall Score: 27.5% 🔴 Weak

Domain Scores:
  Security:   33.2%
  Operations: 21.9%

Cell Breakdown:
  Domain      | Stage    | Maturity | Performance | Cell Score
  ------------|----------|----------|-------------|----------
  security    | design   |    0.0%  |      68.8%  |     41.2%
  security    | build    |    0.0%  |      90.0%  |     54.0%
  security    | test     |    0.0%  |      66.7%  |     40.0%
  security    | runtime  |    0.0%  |       0.0%  |      0.0%
  security    | response |    0.0%  |      85.2%  |     51.1%
  operations  | design   |    0.0%  |       0.0%  |      0.0%
  operations  | build    |    0.0%  |      44.4%  |     26.7%
  operations  | test     |    0.0%  |       0.0%  |      0.0%
  operations  | runtime  |    0.0%  |      91.8%  |     55.1%
  operations  | response |    0.0%  |       0.0%  |      0.0%
```

All maturity scores are 0.0% because the global `MaturityModel` has no cells — goal-level maturity (both at level 3) is not reflected here.

## Recommendation

The data model and core calculations are ready to use as a Go library for building a maturity-based roadmap. The main work needed:

1. **Implement 3 missing CLI command groups** (`goal`, `phase`, `roadmap`)
2. **Fix scoring** to account for goal maturity and sparse metrics
3. **Add `--goals` flag** to `prism score` that uses goal-level maturity instead of global cells

If consuming this as a library (not CLI), it's usable now — `Goal.CurrentMaturityLevel()`, `CalculateGoalProgress()`, and `CalculatePhaseMetrics()` work correctly. The scoring issue is isolated to `CalculatePRISMScore()`.
