# prism maturity report

Generate markdown reports from a maturity model specification file.

## Synopsis

```bash
prism maturity report <maturity-spec-file> [flags]
```

## Description

The `prism maturity report` command generates comprehensive markdown documentation from a maturity model specification. Reports can be viewed from two perspectives:

- **Domain view**: Maturity levels organized by domain (security, operations, quality)
- **Framework view**: Criteria organized by compliance framework mappings

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | stdout | Output file path |
| `--format` | `-f` | `markdown` | Output format: `markdown`, `json` |
| `--view` | `-v` | `both` | View type: `both`, `domain`, `framework` |
| `--title` | | from metadata | Report title |
| `--author` | | | Report author |
| `--no-meta` | | `false` | Omit YAML front matter |
| `--no-toc` | | `false` | Omit table of contents |
| `--no-detail` | | `false` | Omit criterion details (framework mappings) |
| `--frameworks` | | all | Filter to specific frameworks (comma-separated) |

## SLI Catalog

When SLIs are defined in the spec, the report includes an **SLI Catalog** section near the top that provides a quick reference of all metrics grouped by category:

```markdown
## SLI Catalog

### Prevention

| SLI | Metric | Unit | Type | Frameworks |
|-----|--------|------|------|------------|
| Asset Coverage | asset_coverage_pct | % | Quantitative | NIST_CSF_2:ID.AM-1, NIST_800_53:CM-8 |

### Detection

| SLI | Metric | Unit | Type | Frameworks |
|-----|--------|------|------|------------|
| Vuln Scanning | vuln_scan_coverage | % | Quantitative | NIST_CSF_2:ID.RA-1, SOC_2:CC7.1 |

### Response

| SLI | Metric | Unit | Type | Frameworks |
|-----|--------|------|------|------------|
| Security MTTR | security_mttr_days | days | Quantitative | NIST_800_53:IR-6, SOC_2:CC7.4 |
```

This section helps:

- **Compliance teams** quickly see framework coverage
- **Executives** understand what's being measured
- **New team members** get oriented to the metrics landscape

## View Types

### Domain View (`--view domain`)

Shows maturity model organized by domain:

```
Domain
└── Level (M1-M5)
    ├── Description
    ├── Criteria (SLOs)
    │   ├── Quantitative: >=95%, <=7 days
    │   └── Qualitative: Tracked, Implemented
    └── Enablers (implementation tasks)
```

### Framework View (`--view framework`)

Shows criteria organized by compliance framework:

```
Framework (NIST CSF 2.0, NIST 800-53, FedRAMP, etc.)
└── Control
    ├── Domain
    ├── Maturity Level
    ├── Criterion
    └── Status (✅ met, ⏳ pending)
```

## Examples

### Generate to stdout

```bash
prism maturity report maturity-spec.json
```

### Generate to file

```bash
prism maturity report maturity-spec.json -o maturity-report.md
```

### Domain view only

```bash
prism maturity report maturity-spec.json --view domain -o domain-report.md
```

### Framework view only

```bash
prism maturity report maturity-spec.json --view framework -o compliance-report.md
```

### Filter to specific frameworks

```bash
prism maturity report maturity-spec.json --frameworks NIST_CSF_2,NIST_800_53
```

### Without YAML front matter

```bash
prism maturity report maturity-spec.json --no-meta
```

### JSON output

```bash
prism maturity report maturity-spec.json -f json -o maturity-spec.json
```

## Output Format

### YAML Front Matter

By default, reports include YAML front matter for Pandoc/MkDocs compatibility:

```yaml
---
title: "Security Maturity Model"
author: "Security Team"
date: "2024-01-15"
---
```

### Table of Contents

Reports include a navigable table of contents with anchor links:

```markdown
## Table of Contents

### By Domain

- [Security](#security)
  - [Level 1: Reactive](#security-level-1)
  - [Level 2: Basic](#security-level-2)
  ...

### By Framework

- [NIST CSF 2.0](#framework-nist-csf-2)
- [NIST SP 800-53](#framework-nist-800-53)
...
```

### Criteria Tables

Criteria are displayed in tables with status indicators:

```markdown
| ID | Name | Type | Target | Status | Frameworks |
|----|------|------|--------|--------|------------|
| sec-sast | SAST Coverage | Quantitative | >=100% | ✅ | NIST_800_53:SA-11 |
| sec-enc | Encryption at Rest | Qualitative | Tracked | ⏳ | NIST_CSF_2:PR.DS-01 |
```

### Framework Coverage

Framework view shows compliance coverage statistics:

```markdown
**Coverage:** 12/15 controls satisfied (80%)
```

### SLI Resolution

Framework mappings are resolved from SLIs when criteria use `sliId` references:

- Criteria with `sliId` inherit framework mappings from the referenced SLI
- Inline `frameworkMappings` on criteria take precedence if present
- See [SLIs & SLOs](../schema/slos.md) for the SLI architecture

## Supported Frameworks

| Framework | Constant |
|-----------|----------|
| NIST CSF 1.1 | `NIST_CSF` |
| NIST CSF 2.0 | `NIST_CSF_2` |
| NIST SP 800-53 | `NIST_800_53` |
| NIST RMF | `NIST_RMF` |
| NIST AI RMF | `NIST_AI_RMF` |
| FedRAMP High | `FEDRAMP_HIGH` |
| FedRAMP Moderate | `FEDRAMP_MOD` |
| FedRAMP Low | `FEDRAMP_LOW` |
| CIS Controls | `CIS_CONTROLS` |
| SOC 2 | `SOC_2` |
| ISO 27001 | `ISO_27001` |
| DORA | `DORA` |
| SRE | `SRE` |

## Use Cases

### Compliance Reporting

Generate framework-specific compliance reports:

```bash
prism maturity report spec.json --view framework --frameworks FEDRAMP_MOD -o fedramp-moderate.md
```

### Executive Summaries

Generate domain-focused summaries without implementation details:

```bash
prism maturity report spec.json --view domain --no-detail -o executive-summary.md
```

### Documentation Site

Generate MkDocs-compatible documentation:

```bash
prism maturity report spec.json -o docs/maturity/index.md
```

## See Also

- [Maturity Model Concepts](../concepts/maturity.md)
- [Framework Mappings](../concepts/frameworks.md)
- [Schema: Maturity Specification](../schema/maturity.md)
