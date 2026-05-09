# Maturity Model

PRISM uses a 5-level maturity model to assess organizational capability across domains and lifecycle stages.

## Maturity Levels

| Level | Name | Score | Description |
|-------|------|-------|-------------|
| 1 | Reactive | 0.2 | Ad-hoc processes, firefighting mode |
| 2 | Basic | 0.4 | Basic controls, some documentation |
| 3 | Defined | 0.6 | Standardized processes, consistent execution |
| 4 | Managed | 0.8 | Data-driven, measured and controlled |
| 5 | Optimizing | 1.0 | Continuous improvement, automated optimization |

## Level Descriptions

### Level 1: Reactive

- Processes are ad-hoc and chaotic
- Success depends on individual heroics
- No formal documentation
- Firefighting mode is common
- Results are unpredictable

**Reliability Example**: Incidents handled without runbooks
**Efficiency Example**: No deployment automation, manual releases

### Level 2: Basic

- Basic processes are documented
- Some repeatability exists
- Policies are defined but inconsistently applied
- Manual processes predominate
- Limited metrics collection

**Reliability Example**: Basic monitoring with manual alerting
**Efficiency Example**: Basic CI/CD exists but inconsistent

### Level 3: Defined

- Standardized processes across the organization
- Consistent execution of practices
- Documentation is maintained
- Roles and responsibilities are clear
- Metrics are collected systematically

**Reliability Example**: Standardized incident response procedures
**Efficiency Example**: Consistent CI/CD pipelines across all services

### Level 4: Managed

- Data-driven decision making
- Processes are measured and controlled
- Quantitative quality goals
- Variation is understood and addressed
- Predictable outcomes

**Reliability Example**: SLOs with error budgets and automated alerting
**Efficiency Example**: DORA metrics tracked with quantitative targets

### Level 5: Optimizing

- Continuous process improvement
- Innovation and optimization
- Automated remediation where possible
- Proactive risk management
- Industry-leading practices

**Reliability Example**: Self-healing systems, automated capacity management
**Efficiency Example**: On-demand deployments, sub-hour lead time

## Maturity Matrix

PRISM assesses maturity for each domain/stage combination:

|  | Design | Build | Test | Runtime | Response |
|--|--------|-------|------|---------|----------|
| **Reliability** | L3 | L4 | L3 | L4 | L4 |
| **Efficiency** | L3 | L4 | L3 | L4 | L3 |

## Maturity Model Structure

PRISM maturity models consist of:

- **Domains** - Areas of capability (security, operations, quality)
- **Levels** - Maturity stages (M1-M5)
- **Criteria** - SLOs that must be met for each level
- **Enablers** - Tasks/projects that help achieve criteria
- **SLIs** - Shared metric definitions with framework mappings

### MaturityModel

```json
{
  "maturity": {
    "levels": [
      {"level": 1, "name": "Reactive", "description": "..."},
      {"level": 2, "name": "Basic", "description": "..."},
      {"level": 3, "name": "Defined", "description": "..."},
      {"level": 4, "name": "Managed", "description": "..."},
      {"level": 5, "name": "Optimizing", "description": "..."}
    ],
    "cells": [...]
  }
}
```

### MaturityCell

Each cell represents a domain/stage intersection:

```json
{
  "domain": "operations",
  "stage": "build",
  "currentLevel": 4,
  "targetLevel": 5,
  "primaryKPI": "ops-deploy-frequency",
  "kpiTarget": ">=10/day"
}
```

## Maturity Score Calculation

The maturity score for a cell is:

```
MaturityScore = CurrentLevel / 5
```

| Level | Score |
|-------|-------|
| 1 | 0.2 |
| 2 | 0.4 |
| 3 | 0.6 |
| 4 | 0.8 |
| 5 | 1.0 |

## Using Maturity in Go

```go
// Create a maturity model
model := prism.NewMaturityModel()

// Get a specific cell
cell := model.GetCell("operations", "build")
cell.CurrentLevel = 4

// Calculate maturity score
score := cell.CalculateMaturityScore()
fmt.Printf("Maturity: %.1f%%\n", score*100) // 80%

// Create domain-filtered model
opsOnly := prism.NewMaturityModelForDomains([]string{"operations"})
```

## Assessment Guidelines

### Level 1 → 2 (Basic)

- Document existing processes
- Establish basic policies
- Implement foundational tools

### Level 2 → 3 (Defined)

- Standardize across teams
- Create formal procedures
- Establish consistent metrics

### Level 3 → 4 (Managed)

- Define quantitative goals
- Implement continuous monitoring
- Establish feedback loops

### Level 4 → 5 (Optimizing)

- Automate optimization
- Implement predictive capabilities
- Drive continuous improvement

## SLIs and Criteria

PRISM separates metrics (SLIs) from level-specific targets (Criteria/SLOs):

| Concept | Definition | Example |
|---------|------------|---------|
| **SLI** | The metric being measured | "Security MTTR" |
| **Criterion** | Target threshold for a level | "MTTR ≤ 7 days for M4" |

### Why Separate SLIs?

1. **No duplication** - Framework mappings defined once per metric
2. **Consistent metadata** - Unit, type, layer inherited by all criteria
3. **Cleaner exports** - XLSX/Markdown shows framework columns correctly

### Example

```json
{
  "slis": {
    "security-mttr": {
      "id": "security-mttr",
      "name": "Security MTTR",
      "metricName": "security_mttr_days",
      "unit": "days",
      "frameworkMappings": [
        {"framework": "NIST_800_53", "reference": "IR-6"},
        {"framework": "SOC_2", "reference": "CC7.4"}
      ]
    }
  },
  "domains": {
    "security": {
      "levels": [
        {
          "level": 4,
          "criteria": [{
            "id": "sec-m4-mttr",
            "name": "Fast MTTR",
            "sliId": "security-mttr",
            "operator": "lte",
            "target": 7
          }]
        },
        {
          "level": 5,
          "criteria": [{
            "id": "sec-m5-mttr",
            "name": "Rapid MTTR",
            "sliId": "security-mttr",
            "operator": "lte",
            "target": 1
          }]
        }
      ]
    }
  }
}
```

See [SLIs & SLOs](../schema/slos.md) for detailed schema documentation.

## Maturity Weight in PRISM Score

By default, maturity contributes 40% to the PRISM score:

```go
config := &prism.ScoreConfig{
    MaturityWeight:    0.4, // 40%
    PerformanceWeight: 0.6, // 60%
}
```

The cell score formula is:

```
CellScore = (0.4 × MaturityScore) + (0.6 × PerformanceScore)
```
