# PRISM Maturity

**Platform for Reliability, Improvement, and Strategic Maturity**

PRISM Maturity is an Operational Product Management platform for COO-level organizational health monitoring. It provides a unified framework for B2B SaaS health metrics that combines SLOs, maturity modeling, and OKRs into a single coherent system. PRISM enables organizations to understand current health and drive improvement projects across operations, security, quality, product, and AI domains.

## Key Features

- **Unified Metrics Framework** - Combine operations, security, and quality metrics in a single document
- **5-Level Maturity Model** - Track organizational capability from Reactive to Optimizing
- **Dynamic Priority (P0-P3)** - Calculate priority from importance and maturity gap
- **Composite Scoring** - Calculate weighted PRISM scores across domains and lifecycle stages
- **Customer Awareness Tracking** - Model customer awareness states for proactive communication
- **Framework Mappings** - Map metrics to DORA, SRE, NIST CSF, and MITRE ATT&CK frameworks
- **Machine-Evaluable SLOs** - Define SLOs with operators for programmatic checking
- **CLI Tool** - Initialize, validate, and score PRISM documents from the command line
- **JSON Schema** - Auto-generated schema for editor validation and IDE support

## Dynamic Priority

Priority is calculated dynamically from static importance weights and maturity gap:

```
Priority Score = Importance Weight × (Target Level - Current Level)
```

| Score | Priority | Description |
|-------|----------|-------------|
| ≥8 | P0 | Immediate action required |
| ≥4 | P1 | High priority improvement |
| ≥2 | P2 | Scheduled improvement |
| <2 | P3 | Low priority enhancement |

The SLI state tracks priority and rationale:

```json
{
  "sliStates": [{
    "sliId": "sli-sast-coverage",
    "currentLevel": 2,
    "targetLevel": 4,
    "priority": 0,
    "priorityRationale": "Critical security capability with 2-level gap"
  }]
}
```

## Roadmap Planning (v0.2.0)

- **Goal-Driven Maturity** - Define strategic goals with SLO-backed maturity levels
- **Phase Planning** - Organize work into quarters with enter/exit maturity targets
- **Swimlanes** - Group initiatives by domain within phases
- **Initiative Tracking** - Track deployment status and customer adoption
- **Progress Reports** - Generate roadmap reports and SLO compliance summaries

## Operations Management (v0.3.0)

- **Value Stream Layers** - Full lifecycle from requirements through support (requirements, code, infra, runtime, adoption, support)
- **Quality Domain** - Add quality alongside operations and security with ISO 25010 verticals
- **Team Topology** - Model stream-aligned, platform, enabling, and overlay teams
- **Service Ownership** - Define services with team ownership and layer assignment
- **Golden Signals** - Associate latency, traffic, errors, saturation metrics per layer
- **Adoption Metrics** - Track product analytics (Pendo, Amplitude) and user engagement
- **Support Metrics** - Monitor ticket volume, resolution time, customer satisfaction

## Analysis & Export (v0.3.0)

- **Document Analysis** - Analyze goals, phases, and SLOs to identify gaps and recommendations
- **Gap Identification** - Find maturity gaps, SLO compliance issues, and initiative gaps
- **Roadmap Export** - Export to structured-plan/roadmap format with deployment/adoption tracking
- **OKR Export** - Export to Objectives and Key Results format for structured-plan
- **V2MOM Export** - Export to Vision, Values, Methods, Obstacles, Measures format
- **LLM Prompts** - Generate prompts for AI-assisted initiative planning

## Enterprise Features

- **Dashforge Integration** - Embed PRISM dashboards in dashforge sites or standalone pages
- **Marp Presentations** - Generate executive presentations from PRISM documents
- **Excel Export** - Export metrics and scores to XLSX for stakeholder reporting

## Quick Example

```json
{
  "$schema": "https://github.com/grokify/prism-maturity/schema/prism-maturity-plan.schema.json",
  "metadata": {
    "name": "Acme Corp PRISM",
    "version": "1.0.0"
  },
  "teams": [
    {
      "id": "payments-team",
      "name": "Payments Team",
      "type": "stream_aligned",
      "serviceIds": ["payments-api"]
    }
  ],
  "services": [
    {
      "id": "payments-api",
      "name": "Payments API",
      "ownerTeamId": "payments-team",
      "layerId": "runtime",
      "tier": "tier1"
    }
  ],
  "metrics": [
    {
      "id": "ops-availability",
      "name": "Service Availability",
      "domain": "operations",
      "stage": "runtime",
      "category": "reliability",
      "layer": "runtime",
      "serviceId": "payments-api",
      "metricType": "rate",
      "current": 99.95,
      "target": 99.99,
      "slo": {
        "target": ">=99.99%",
        "operator": "gte",
        "value": 99.99,
        "window": "30d"
      }
    }
  ]
}
```


## PRISM Score

The PRISM score provides a single composite metric (0.0-1.0) representing organizational health:

| Score | Level | Description |
|-------|-------|-------------|
| ≥0.90 | Elite | Industry-leading practices |
| ≥0.75 | Strong | Well-managed, proactive |
| ≥0.50 | Medium | Adequate, room for improvement |
| ≥0.25 | Weak | Significant gaps |
| <0.25 | Critical | Immediate attention required |

## Integration with PRISM Execution

PRISM Maturity integrates with [prism-execution](https://github.com/grokify/prism-execution) to provide a complete operational planning workflow. PRISM Maturity serves as the source of truth for requirements (maturity models, SLOs), while prism-execution handles execution tracking (OKRs, roadmaps).

### Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                  PRISM Maturity (Source of Truth)          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │    Goals     │  │   Maturity   │  │     SLOs     │          │
│  │              │  │    Models    │  │              │          │
│  │ "Reliability"│  │  M1→M2→M3→M4 │  │ avail>=99.9% │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│                      LLM Analysis Layer                        │
│                                                                │
│  Input: PRISM requirements + context (team capacity, etc.)     │
│  Output: Suggested initiatives, phase mapping, dependencies    │
│                                                                │
│  "To reach M3 reliability by Q2, you need:                     │
│   - Q1: Observability platform (enables SLO measurement)       │
│   - Q1: Incident runbooks (reduces MTTR)                       │
│   - Q2: SLO dashboards (tracks compliance)"                    │
└───────────────────────────┬────────────────────────────────────┘
                            │
                            ▼
┌────────────────────────────────────────────────────────────────┐
│                   PRISM Execution (Execution)                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │     OKR      │  │   V2MOM      │  │   Roadmap    │          │
│  │              │  │              │  │              │          │
│  │ Objectives   │  │ Methods      │  │ Phases       │          │
│  │ Key Results  │  │ Measures     │  │ Deliverables │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└────────────────────────────────────────────────────────────────┘
```

### Data Flow

| PRISM Concept | Roadmap Export | OKR Export |
|---------------|----------------|------------|
| Phase | Phase | — |
| Phase.GoalTargets | Phase.Goals, SuccessCriteria | PhaseTargets |
| Goal | — | Objective (one per maturity level) |
| SLO/MetricCriterion | SuccessCriteria | Key Result |
| Initiative | Deliverable | — |
| DeploymentStatus | RolloutStatus | — |

**RolloutStatus** tracks B2B SaaS deployment and adoption:

- **deployedCustomers** - Customers with feature available (rolled out)
- **adoptedCustomers** - Customers actively using the feature
- **status** - `not_started`, `rolling_out`, `deployed`, `adopted`, `paused`, `rolled_back`

### Workflow

1. **Define requirements in PRISM Maturity** - Goals, maturity models, SLOs
2. **Analyze with LLM** - Generate initiative recommendations to achieve targets
3. **Export to prism-execution** - Roadmap/OKR/V2MOM format
4. **Track execution** - Monitor progress against phase targets

```bash
# Analyze PRISM document and suggest initiatives
prism analyze prism.json

# Export as roadmap with deployment tracking
prism export roadmap prism.json --with-okrs -o roadmap.json

# Or export as OKR-only document
prism export okr prism.json -o roadmap.okr.json
```

## Getting Started

```bash
# Install the CLI
go install github.com/grokify/prism-maturity/cmd/prism@latest

# Initialize a new document
prism init -o prism.json

# Validate the document
prism validate prism.json

# Calculate the PRISM score
prism score prism.json --detailed
```

## License

MIT
