# FEAT_MATURITYROADMAP: Goal-Driven Maturity Roadmap

## Problem Statement

Organizations need to:

1. **Track maturity progression** against strategic Goals across a multi-phase roadmap
2. **Define measurable criteria** (SLOs) that must be met to achieve each maturity level
3. **Monitor initiative execution** and customer adoption as they work toward Goals
4. **Report progress** at phase boundaries with clear enter/exit metrics

Currently, PRISM supports metrics, SLOs, and a global maturity model (per domain/stage), but lacks:

- Goal-level planning and maturity tracking
- Phase/quarter-based roadmap organization
- Initiative-to-Goal linkage
- SLO requirements per maturity level per Goal

## Solution Overview

Extend PRISM with a **Goal-driven Maturity Roadmap** layer that:

1. Defines **Goals** with their own 5-level maturity models
2. Organizes work into **Phases** (quarters) with swimlanes
3. Links **Initiatives** to Goals and Phases
4. Specifies **SLO requirements** for each maturity level
5. Tracks **progress metrics** at phase boundaries

## Core Concepts

### Concept Hierarchy

```
Roadmap
├── Goals (strategic objectives)
│   └── GoalMaturityModel (5 levels per goal)
│       └── GoalMaturityLevel
│           └── RequiredSLOs (must be met for this level)
├── Phases (Q1, Q2, Q3, Q4...)
│   └── Swimlanes (by domain or stage)
│       └── Initiatives (work items)
└── PhaseMetrics (enter/exit criteria)
```

### Relationship Model

```
Goal 1:N Initiative    (A goal has many initiatives)
Phase 1:N Initiative   (A phase contains many initiatives)
Goal 1:1 GoalMaturityModel (Each goal has its own maturity model)
GoalMaturityLevel 1:N SLO (Each level requires specific SLOs to be met)
```

## Requirements

### R1: Goals

A **Goal** represents a strategic objective that the organization wants to achieve.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Goal name |
| `description` | string | No | Business context and rationale |
| `owner` | string | No | Goal owner/sponsor |
| `priority` | int | No | Priority ranking (1 = highest) |
| `status` | string | No | active, on_hold, completed, cancelled |
| `startDate` | string | No | Goal start date (ISO 8601) |
| `targetDate` | string | No | Target completion date |
| `maturityModel` | GoalMaturityModel | Yes | 5-level maturity definition |
| `currentLevel` | int | No | Current maturity level (1-5) |
| `targetLevel` | int | No | Target maturity level for current period |

**Example:**

```json
{
  "id": "goal-reliability",
  "name": "Achieve High Reliability",
  "description": "Achieve industry-leading reliability practices across all lifecycle stages",
  "owner": "VP Engineering",
  "priority": 1,
  "status": "active",
  "targetDate": "2026-12-31",
  "currentLevel": 2,
  "targetLevel": 4,
  "maturityModel": { ... }
}
```

### R2: Goal Maturity Model

Each Goal has its own **5-level maturity model** defining what it means to progress from Level 1 (Reactive) to Level 5 (Optimizing) for that specific goal.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `levels` | []GoalMaturityLevel | Yes | Array of 5 level definitions |

### R3: Goal Maturity Level with SLO Requirements

Each level specifies:

1. What the level means for this goal
2. Which **SLOs must be met** to achieve this level
3. Optional metrics criteria

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `level` | int | Yes | Level number (1-5) |
| `name` | string | Yes | Level name |
| `description` | string | No | What this level means for the goal |
| `requiredSLOs` | []SLORequirement | Yes | SLOs that must be met |
| `metricCriteria` | []MetricCriterion | No | Additional metric requirements |

**SLORequirement:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `metricId` | string | Yes | Reference to a Metric with SLO |
| `description` | string | No | Why this SLO matters for this level |

**MetricCriterion:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `metricId` | string | Yes | Reference to a Metric |
| `operator` | string | Yes | gte, lte, gt, lt, eq |
| `value` | float64 | Yes | Required value |

**Example:**

```json
{
  "levels": [
    {
      "level": 1,
      "name": "Reactive",
      "description": "Ad-hoc practices, no systematic monitoring",
      "requiredSLOs": []
    },
    {
      "level": 2,
      "name": "Basic",
      "description": "Basic monitoring in place, manual alerting",
      "requiredSLOs": [
        { "metricId": "ops-availability", "description": "Basic availability tracking" }
      ],
      "metricCriteria": [
        { "metricId": "ops-availability", "operator": "gte", "value": 99.0 }
      ]
    },
    {
      "level": 3,
      "name": "Defined",
      "description": "SLOs tracked, incident response documented",
      "requiredSLOs": [
        { "metricId": "ops-availability", "description": "Production SLO" },
        { "metricId": "ops-mttr", "description": "Recovery time tracked" }
      ],
      "metricCriteria": [
        { "metricId": "ops-availability", "operator": "gte", "value": 99.5 },
        { "metricId": "ops-mttr", "operator": "lte", "value": 4 }
      ]
    },
    {
      "level": 4,
      "name": "Managed",
      "description": "Proactive monitoring, error budgets active",
      "requiredSLOs": [
        { "metricId": "ops-availability", "description": "99.9% availability" },
        { "metricId": "ops-mttr", "description": "1-hour MTTR" },
        { "metricId": "ops-p99-latency", "description": "Latency SLO" }
      ],
      "metricCriteria": [
        { "metricId": "ops-availability", "operator": "gte", "value": 99.9 },
        { "metricId": "ops-mttr", "operator": "lte", "value": 1 },
        { "metricId": "ops-p99-latency", "operator": "lte", "value": 200 }
      ]
    },
    {
      "level": 5,
      "name": "Optimizing",
      "description": "Self-healing systems, chaos engineering",
      "requiredSLOs": [
        { "metricId": "ops-availability", "description": "99.99% availability" },
        { "metricId": "ops-mttr", "description": "15-minute MTTR" },
        { "metricId": "ops-p99-latency", "description": "Sub-100ms latency" },
        { "metricId": "ops-deploy-frequency", "description": "On-demand deployments" }
      ],
      "metricCriteria": [
        { "metricId": "ops-availability", "operator": "gte", "value": 99.99 },
        { "metricId": "ops-mttr", "operator": "lte", "value": 0.25 },
        { "metricId": "ops-p99-latency", "operator": "lte", "value": 100 },
        { "metricId": "ops-deploy-frequency", "operator": "gte", "value": 10 }
      ]
    }
  ]
}
```

### R4: Phases (Roadmap Quarters)

A **Phase** represents a time-bounded period (typically a quarter) for planning and execution.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Phase name (e.g., "Q1 2026") |
| `quarter` | string | No | Quarter identifier (Q1, Q2, Q3, Q4) |
| `year` | int | No | Year |
| `startDate` | string | Yes | Phase start date |
| `endDate` | string | Yes | Phase end date |
| `status` | string | No | planning, in_progress, completed |
| `goalTargets` | []PhaseGoalTarget | No | Goal maturity targets for this phase |
| `swimlanes` | []Swimlane | No | Initiatives organized by domain/stage |

**PhaseGoalTarget:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `goalId` | string | Yes | Reference to Goal |
| `enterLevel` | int | Yes | Expected maturity level at phase start |
| `exitLevel` | int | Yes | Target maturity level at phase end |

### R5: Swimlanes

Swimlanes organize initiatives within a phase by domain or stage.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Swimlane name |
| `domain` | string | No | security, operations |
| `stage` | string | No | design, build, test, runtime, response |
| `initiativeIds` | []string | Yes | Initiatives in this swimlane |

### R6: Enhanced Initiatives

Extend the existing Initiative type with Goal and Phase linkage.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `goalIds` | []string | No | Goals this initiative contributes to |
| `phaseId` | string | No | Phase this initiative is scheduled in |
| `devCompletionPercent` | float64 | No | Development completion (0-100) |
| `deploymentStatus` | DeploymentStatus | No | Customer adoption tracking |

**DeploymentStatus:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `status` | string | Yes | not_started, in_progress, completed |
| `totalCustomers` | int | No | Total customers to deploy to |
| `deployedCustomers` | int | No | Customers deployed |
| `adoptionPercent` | float64 | No | Calculated adoption percentage |

### R7: Phase Metrics (Enter/Exit Criteria)

Track progress at phase boundaries.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `phaseId` | string | Yes | Reference to Phase |
| `goalProgress` | []GoalProgress | No | Progress per goal |
| `initiativeMetrics` | InitiativeMetrics | No | Initiative completion stats |
| `sloCompliance` | []SLOCompliance | No | SLO status at phase boundary |

**GoalProgress:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `goalId` | string | Yes | Reference to Goal |
| `enterLevel` | int | Yes | Maturity level at phase entry |
| `currentLevel` | int | Yes | Current maturity level |
| `targetLevel` | int | Yes | Target level for phase exit |
| `initiativesTotal` | int | Yes | Total initiatives for this goal |
| `initiativesCompleted` | int | Yes | Completed initiatives |
| `completionPercent` | float64 | Yes | initiativesCompleted / initiativesTotal |
| `slosRequired` | int | Yes | SLOs required for target level |
| `slosMet` | int | Yes | SLOs currently met |

**InitiativeMetrics:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `total` | int | Yes | Total initiatives in phase |
| `completed` | int | Yes | Dev-complete initiatives |
| `deployed` | int | Yes | Fully deployed initiatives |
| `avgAdoptionPercent` | float64 | No | Average adoption across completed |

## Calculations

### Goal Completion Percentage

```
completionPercent = initiativesCompleted / initiativesTotal * 100
```

Where:
- `initiativesCompleted` = count of initiatives with `status == "completed"` and `goalIds` contains this goal
- `initiativesTotal` = total count of initiatives where `goalIds` contains this goal

### Maturity Level Achievement

A Goal achieves maturity level N when:

1. **All required SLOs for level N are met** (`metric.MeetsSLO() == true` for all SLOs in `requiredSLOs`)
2. **All metric criteria for level N are satisfied** (metric values meet operator/value requirements)

```go
func (g *Goal) CurrentMaturityLevel() int {
    for level := 5; level >= 1; level-- {
        if g.MeetsLevelRequirements(level) {
            return level
        }
    }
    return 1
}

func (g *Goal) MeetsLevelRequirements(level int) bool {
    levelDef := g.MaturityModel.GetLevel(level)

    // Check all required SLOs
    for _, sloReq := range levelDef.RequiredSLOs {
        metric := g.GetMetric(sloReq.MetricID)
        if !metric.MeetsSLO() {
            return false
        }
    }

    // Check all metric criteria
    for _, criterion := range levelDef.MetricCriteria {
        metric := g.GetMetric(criterion.MetricID)
        if !criterion.IsMet(metric.Current) {
            return false
        }
    }

    return true
}
```

### Phase Goal Progress

```go
func (p *Phase) GoalProgress(goalId string) *GoalProgress {
    goal := p.GetGoal(goalId)
    initiatives := p.GetInitiativesForGoal(goalId)

    completed := 0
    for _, init := range initiatives {
        if init.Status == "completed" {
            completed++
        }
    }

    targetLevel := p.GetGoalTarget(goalId).ExitLevel
    levelDef := goal.MaturityModel.GetLevel(targetLevel)

    slosMet := 0
    for _, sloReq := range levelDef.RequiredSLOs {
        if goal.GetMetric(sloReq.MetricID).MeetsSLO() {
            slosMet++
        }
    }

    return &GoalProgress{
        GoalID:               goalId,
        EnterLevel:           p.GetGoalTarget(goalId).EnterLevel,
        CurrentLevel:         goal.CurrentMaturityLevel(),
        TargetLevel:          targetLevel,
        InitiativesTotal:     len(initiatives),
        InitiativesCompleted: completed,
        CompletionPercent:    float64(completed) / float64(len(initiatives)) * 100,
        SLOsRequired:         len(levelDef.RequiredSLOs),
        SLOsMet:              slosMet,
    }
}
```

## Schema Changes

### PRISMDocument Extensions

```go
type PRISMDocument struct {
    // Existing fields...
    Schema      string          `json:"$schema,omitempty"`
    Metadata    *Metadata       `json:"metadata,omitempty"`
    Metrics     []Metric        `json:"metrics"`
    Maturity    *MaturityModel  `json:"maturity,omitempty"`

    // New fields
    Goals       []Goal          `json:"goals,omitempty"`
    Phases      []Phase         `json:"phases,omitempty"`
    Roadmap     *RoadmapConfig  `json:"roadmap,omitempty"`
}
```

### New Types Summary

| Type | File | Purpose |
|------|------|---------|
| `Goal` | goal.go | Strategic objective with maturity model |
| `GoalMaturityModel` | goal.go | 5-level maturity definition per goal |
| `GoalMaturityLevel` | goal.go | Level definition with SLO requirements |
| `SLORequirement` | goal.go | SLO reference for maturity level |
| `MetricCriterion` | goal.go | Metric value requirement |
| `Phase` | phase.go | Time-bounded planning period |
| `PhaseGoalTarget` | phase.go | Goal maturity target per phase |
| `Swimlane` | phase.go | Initiative grouping within phase |
| `DeploymentStatus` | initiative.go | Customer adoption tracking |
| `PhaseMetrics` | phase_metrics.go | Enter/exit criteria |
| `GoalProgress` | phase_metrics.go | Goal progress within phase |
| `InitiativeMetrics` | phase_metrics.go | Initiative stats |
| `SLOCompliance` | phase_metrics.go | SLO status tracking |

## CLI Extensions

### New Commands

```bash
# Goal operations
prism goal list                     # List all goals
prism goal show <goal-id>           # Show goal details with maturity status
prism goal progress <goal-id>       # Show goal progress (initiatives, SLOs)

# Phase operations
prism phase list                    # List all phases
prism phase show <phase-id>         # Show phase details with swimlanes
prism phase metrics <phase-id>      # Show enter/exit metrics

# Roadmap operations
prism roadmap show                  # Show roadmap overview
prism roadmap progress              # Show progress across all phases/goals
```

### Enhanced Score Command

```bash
prism score prism.json --goals      # Include goal maturity in score
prism score prism.json --phase Q1   # Score for specific phase
```

## Example Document

See `examples/goal-roadmap.json` for a complete example.

## Success Criteria

1. **Goals** can be defined with their own maturity models
2. **Maturity levels** are backed by specific SLOs
3. **Initiatives** can be linked to Goals and Phases
4. **Phase metrics** show enter/exit maturity levels and SLO compliance
5. **Progress** is calculable: completion %, adoption %, SLOs met
6. **Validation** ensures SLO references are valid
7. **CLI** supports goal/phase/roadmap operations

## Out of Scope (Future)

- Visual roadmap rendering (defer to uiforge)
- Dependency tracking between initiatives
- Resource/capacity planning
- Automated SLO target suggestions
- Historical trend analysis across phases

## References

- Current PRISM types: `prism.go`, `maturity.go`
- SLO implementation: `prism.go` (SLO, MeetsSLO)
- Maturity model: `maturity.go` (MaturityModel, MaturityCell)
