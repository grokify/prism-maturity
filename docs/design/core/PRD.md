# PRISM Product Requirements Document

## Vision

PRISM (Program Rating & Improvement System for Maturity) provides a unified framework for measuring and improving organizational maturity across multiple domains: Security, Operations, Quality, and Product.

## Problem Statement

Organizations struggle to:

1. **Measure maturity consistently** - Different teams use different metrics and definitions
2. **Map to compliance frameworks** - NIST, SOC 2, FedRAMP requirements are tracked separately
3. **Track improvement over time** - No unified view of progress across domains
4. **Communicate to stakeholders** - Technical metrics don't translate to executive dashboards

## Solution

PRISM provides:

1. **Unified maturity model** - 5-level (M1-M5) scale applicable to all domains
2. **SLI/SLO-backed levels** - Each level defined by measurable criteria, not vague requirements
3. **Framework mappings** - Built-in support for NIST CSF, NIST 800-53, FedRAMP, SOC 2, DORA, SRE
4. **Machine-readable specs** - JSON/YAML specs that generate reports, dashboards, presentations
5. **Export formats** - Markdown, XLSX, Marp slides for different audiences

## Key Design Decisions

### SLI vs SLO Separation

**Decision:** Define metrics (SLIs) separately from level-specific targets (Criteria/SLOs).

**Rationale:**

- Framework mappings defined once per metric, not repeated per level
- Cleaner exports with consistent metadata
- Enables SLI Catalog for quick metric reference

### Quantitative + Qualitative SLOs

**Decision:** Support both numeric thresholds and binary state tracking.

**Rationale:**

- Quantitative: "MTTR ≤ 7 days" - measurable outcomes
- Qualitative: "Encryption enabled" - binary compliance states
- Both are valid maturity indicators

### Domain-Centric Organization

**Decision:** Organize by domain (Security, Operations, Quality) rather than team or technology.

**Rationale:**

- Domains align with organizational concerns
- Cross-cutting view of capability
- Easier to map to compliance frameworks

## Domains

| Domain | Description | Example Metrics |
|--------|-------------|-----------------|
| Security | Application and infrastructure security | SAST coverage, MTTR, secret detection |
| Operations | Reliability, efficiency, deployment | Availability, change failure rate, lead time |
| Quality | Code quality, test coverage | Test coverage, defect density, code review |
| Product | Feature delivery, customer outcomes | Cycle time, feature adoption |

## Maturity Levels

| Level | Name | Description |
|-------|------|-------------|
| M1 | Reactive | Ad-hoc processes, firefighting mode |
| M2 | Basic | Basic controls in place, some documentation |
| M3 | Defined | Standardized processes, consistent execution |
| M4 | Managed | Data-driven, measured and controlled |
| M5 | Optimizing | Continuous improvement, proactive automation |

## Framework Support

PRISM maps criteria to industry standards:

- **NIST CSF 2.0** - Cybersecurity Framework
- **NIST SP 800-53** - Security and Privacy Controls
- **FedRAMP** - Federal Risk and Authorization Management
- **SOC 2** - Service Organization Controls
- **ISO 27001** - Information Security Management
- **DORA** - Digital Operational Resilience Act
- **CIS Controls** - Center for Internet Security
- **SRE** - Site Reliability Engineering practices

## Quality Model Alignment

PRISM quality metrics align with ISO/IEC 25010:

| ISO 25010 Characteristic | PRISM Category | Example Metrics |
|--------------------------|----------------|-----------------|
| Functional Suitability | quality | Test coverage, requirement coverage |
| Reliability | reliability | Availability, MTBF, error rate |
| Performance Efficiency | efficiency | Latency, throughput, resource usage |
| Security | security | Vulnerability count, SAST coverage |
| Maintainability | quality | Code complexity, documentation |
| Portability | operations | Deployment frequency, containerization |

## Export Formats

| Format | Use Case | Command |
|--------|----------|---------|
| Markdown | Documentation sites, GitHub | `prism maturity report spec.json -o report.md` |
| XLSX | Executive dashboards, compliance | `prism maturity xlsx spec.json -o report.xlsx` |
| JSON | API integration, automation | `prism maturity report spec.json -f json` |
| Marp | Presentations | `prism export marp spec.json -o slides.md` |

## Success Metrics

1. **Adoption** - Teams using PRISM for maturity tracking
2. **Coverage** - Percentage of services with defined maturity levels
3. **Accuracy** - Automated vs manual assessment correlation
4. **Time savings** - Reduction in compliance reporting effort

## Future Directions

1. **Automated data collection** - Pull metrics from observability platforms
2. **Trend analysis** - Historical tracking with improvement visualization
3. **Benchmark data** - Anonymous industry comparisons
4. **AI recommendations** - Suggested enablers based on gaps

## References

- [Maturity Model Concepts](../../concepts/maturity.md)
- [SLI/SLO Architecture](../../schema/slos.md)
- [Framework Mappings](../../concepts/frameworks.md)
- [Maturity Report CLI](../../cli/maturity-report.md)
