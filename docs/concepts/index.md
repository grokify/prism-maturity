# Concepts Overview

PRISM provides a multi-dimensional model for organizing metrics, ownership, and maturity tracking across your organization.

## The PRISM Model

PRISM organizes metrics using four key dimensions:

| Dimension | Purpose | Values |
|-----------|---------|--------|
| **Domain** | Functional area with standards | operations, security, quality |
| **Layer** | Ownership boundary in stack | code, infra, runtime |
| **Stage** | Lifecycle phase | design, build, test, runtime, response |
| **Category** | Type of control | prevention, detection, response, reliability, efficiency, quality |

These dimensions work together to classify every metric and clarify accountability.

## Organizational Model

### [Domains](../schema/domains.md)

Domains represent functional areas with overlay teams that define standards:

| Domain | Description | Overlay Team |
|--------|-------------|--------------|
| Operations | Reliability, performance, efficiency | SRE/Platform |
| Security | AppSec, CloudSec, compliance | Security Team |
| Quality | Testing, code quality, defects | QE Team |

### [Layers](layers.md)

Layers represent the full value stream from ideation to support:

| Layer | Description | Typical Owner |
|-------|-------------|---------------|
| Requirements | Product ideation, specs, design | Product/Design |
| Code | Application code, libraries, dependencies | Stream-aligned teams |
| Infra | Cloud resources, networking, platform | Platform team |
| Runtime | Running services, production workloads | Stream-aligned + SRE |
| Adoption | Product analytics, user engagement | Product/Growth |
| Support | Customer support, incident management | Support/CS |

Each layer can define **golden signals** (latency, traffic, errors, saturation) pointing to specific metrics.

### [Teams](teams.md)

Teams follow the Team Topologies model for clear accountability:

| Type | Role | Owns |
|------|------|------|
| Stream-Aligned | Build and run services | Services, code/runtime metrics |
| Platform | Provide infrastructure | Infrastructure layer |
| Enabling | Help adopt practices | Cross-team capabilities |
| Overlay | Define standards | Domain standards |

### [Services](services.md)

Services are deployable units that connect teams, layers, and metrics:

- Owned by a team
- Deployed in a layer
- Associated with metrics and SLOs

## Metrics and SLOs

### Metric Classification

Every metric belongs to:

- A **domain** (what functional area)
- A **stage** (when in lifecycle)
- A **category** (what type of control)
- Optionally a **layer** (where in stack)
- Optionally a **service** (which system)
- Optionally a **quality vertical** (ISO 25010 characteristic)

### SLO Definition

Metrics can have machine-evaluable SLOs:

```json
{
  "slo": {
    "target": ">=99.99%",
    "operator": "gte",
    "value": 99.99,
    "window": "30d"
  }
}
```

## Maturity Roadmap

### [Maturity Model](maturity.md)

Maturity models define levels (M1-M5) with criteria and enablers:

- **SLIs** - Service Level Indicators define metrics with framework mappings
- **Criteria** - SLO-based requirements that reference SLIs
- **Enablers** - Tasks that help achieve criteria

### [Goals](goals.md)

Goals represent strategic objectives with SLO-backed maturity levels:

- Define 5-level maturity progression
- Specify which SLOs must be met at each level
- Track progress from Reactive to Optimizing

### [Phases](phases.md)

Phases organize work into time-bounded periods (quarters):

- Set goal maturity targets (enter/exit levels)
- Group initiatives into swimlanes
- Track completion and SLO compliance

### Connecting It All

The PRISM model creates a clear chain of accountability:

```
Team owns → Service has → Metrics with → SLOs required by → Goals tracked in → Phases
```

**Example flow:**

1. **Payments Team** (stream-aligned) owns **Payments API** (service)
2. **Payments API** has **Availability** metric (operations/runtime/reliability)
3. **Availability** has SLO: `>=99.99%` over 30 days
4. **Reliability Goal** requires this SLO at Level 4
5. **Q2 2026 Phase** targets Reliability Goal from Level 3 to Level 4
6. **Security Team** (overlay) defines security standards that apply to all services

## Scoring

### [PRISM Score](scoring.md)

A composite health score (0.0-1.0) combining:

- **Maturity scores** (40% weight) - organizational capability
- **Performance scores** (60% weight) - metric achievement
- **Awareness multiplier** - customer communication effectiveness

### [Customer Awareness](awareness.md)

Track customer awareness through four states:

| State | Weight | Description |
|-------|--------|-------------|
| Unaware | 0.0 | Customer not aware |
| Aware (not acting) | 0.25 | Aware but not remediating |
| Remediating | 0.5 | Actively working on fix |
| Remediated | 1.0 | Issue resolved |

### [Framework Mappings](frameworks.md)

Map metrics to industry standards:

- DORA metrics
- SRE practices
- NIST Cybersecurity Framework
- MITRE ATT&CK

## Getting Started

1. **Define your organizational model**
   - Identify domains relevant to your org
   - Map teams to Team Topologies types
   - List services and their ownership

2. **Define metrics**
   - Classify by domain, stage, layer
   - Set baselines and targets
   - Define SLOs with operators

3. **Create goals**
   - Define maturity levels
   - Specify SLO requirements per level
   - Link initiatives to goals

4. **Plan phases**
   - Set quarterly goal targets
   - Organize initiatives into swimlanes
   - Track progress and SLO compliance

5. **Calculate and report**
   - Generate PRISM scores
   - Create roadmap reports
   - Track maturity progression
