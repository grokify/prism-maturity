# SLIs & SLOs

PRISM supports Service Level Indicators (SLIs) and Service Level Objectives (SLOs) following SRE best practices.

## SLI vs SLO Architecture

PRISM distinguishes between SLIs (what is measured) and SLOs (target thresholds):

| Concept | Definition | Example | Framework Mappings |
|---------|------------|---------|-------------------|
| **SLI** | The metric being measured | "MTTR" (Mean Time to Recovery) | Defined here |
| **SLO** | Target threshold for a specific maturity level | "MTTR ≤ 7 days for M4" | Inherited from SLI |

**Key principle:** Framework mappings are defined once at the SLI level, not repeated for each SLO. The metric "Security MTTR" maps to NIST/SOC2 controls regardless of the maturity level threshold.

## SLI (Service Level Indicator)

An SLI defines what is being measured, including any compliance framework mappings.

### SLI Fields (Maturity Model)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | SLI name |
| `description` | string | No | SLI description |
| `metricName` | string | Yes | The metric being measured |
| `unit` | string | No | Unit of measurement (%, hours, days) |
| `type` | string | No | `quantitative` (default) or `qualitative` |
| `layer` | string | No | Value stream layer (code, infra, runtime, etc.) |
| `category` | string | No | Metric category (reliability, efficiency, security) |
| `frameworkMappings` | array | No | Compliance framework references |

### SLI Example (Maturity Model)

```json
{
  "slis": {
    "security-mttr": {
      "id": "security-mttr",
      "name": "Security MTTR",
      "description": "Mean time to remediate critical security findings",
      "metricName": "security_mttr_days",
      "unit": "days",
      "type": "quantitative",
      "layer": "runtime",
      "category": "response",
      "frameworkMappings": [
        {"framework": "NIST_800_53", "reference": "IR-6", "name": "Incident Reporting"},
        {"framework": "SOC_2", "reference": "CC7.4", "name": "Security Incident Response"}
      ]
    }
  }
}
```

### SLI Fields (Embedded in Metric)

For standalone metrics (not maturity model), SLIs can be embedded:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | No | SLI name |
| `description` | string | No | SLI description |
| `formula` | string | No | Calculation formula |

### Example (Embedded SLI)

```json
{
  "sli": {
    "name": "Availability",
    "description": "Percentage of successful requests",
    "formula": "successful_requests / total_requests * 100"
  }
}
```

## SLO (Service Level Objective)

An SLO defines the target for an SLI.

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | No | SLO type: `quantitative` (default) or `qualitative` |
| `target` | string | Yes | Human-readable target (e.g., ">=99.9%" or "Tracked") |
| `operator` | string | No | Comparison operator for machine evaluation |
| `value` | number | No | Numeric target value (quantitative only) |
| `window` | string | No | Measurement window (quantitative only) |
| `status` | string | No | Current qualitative status (qualitative only) |
| `thresholds` | object | No | Additional thresholds (quantitative only) |

### SLO Types

PRISM supports two types of SLOs:

| Type | Description | Example Use Case |
|------|-------------|------------------|
| `quantitative` | Numeric comparison (default) | Availability ≥99.99%, Latency ≤200ms |
| `qualitative` | Binary state tracking | "Metric is tracked", "Process is defined" |

### SLO Operators

PRISM supports machine-evaluable SLOs with these operators:

| Operator | Symbol | Type | Description | Example |
|----------|--------|------|-------------|---------|
| `gte` | `>=` | Quantitative | Greater than or equal | Availability ≥99.99% |
| `lte` | `<=` | Quantitative | Less than or equal | Latency ≤200ms |
| `gt` | `>` | Quantitative | Greater than | Score >80 |
| `lt` | `<` | Quantitative | Less than | Error rate <0.1% |
| `eq` | `=` | Quantitative | Equal to | Target exactly 100 |
| `exists` | Tracked | Qualitative | Metric exists/is tracked | Metric is being tracked |

### Example with Quantitative SLO

```json
{
  "slo": {
    "type": "quantitative",
    "target": ">=99.99%",
    "operator": "gte",
    "value": 99.99,
    "window": "30d"
  }
}
```

### Example with Qualitative SLO

Qualitative SLOs track binary states rather than numeric values:

```json
{
  "slo": {
    "type": "qualitative",
    "target": "Tracked",
    "operator": "exists",
    "status": "tracked"
  }
}
```

### Qualitative Status Values

| Status | Indicates Met | Description |
|--------|---------------|-------------|
| `tracked` | Yes | Metric is being tracked |
| `implemented` | Yes | Control/feature is implemented |
| `defined` | Yes | Process/policy is defined |
| `documented` | Yes | Documentation exists |
| `compliant` | Yes | Meets compliance requirement |
| `enabled` | Yes | Feature/capability is enabled |
| `not_tracked` | No | Not yet being tracked |
| `partial` | No | Partially implemented |
| `planned` | No | Planned but not started |

### Qualitative SLO Use Cases

Qualitative SLOs are useful for:

- **Process maturity** - "Incident runbooks are documented"
- **Compliance tracking** - "SOC 2 controls are implemented"
- **Capability existence** - "Alerting is enabled for critical services"
- **Observability** - "Distributed tracing is tracked"

### Measurement Windows

Common window values:

| Window | Description |
|--------|-------------|
| `7d` | 7-day rolling window |
| `30d` | 30-day rolling window |
| `90d` | 90-day rolling window |
| `monthly` | Calendar month |
| `quarterly` | Calendar quarter |

## Programmatic SLO Checking

PRISM provides a `MeetsSLO()` method for programmatic checking:

```go
metric := prism.Metric{
    Current: 99.95,
    SLO: &prism.SLO{
        Target:   ">=99.99%",
        Operator: prism.SLOOperatorGTE,
        Value:    99.99,
    },
}

if metric.MeetsSLO() {
    fmt.Println("SLO met!")
} else {
    fmt.Println("SLO not met")
}
```

### Operator Behavior

| Operator | Current | Value | MeetsSLO() |
|----------|---------|-------|------------|
| `gte` | 99.99 | 99.99 | true |
| `gte` | 99.95 | 99.99 | false |
| `lte` | 200 | 250 | true |
| `lte` | 300 | 250 | false |
| `eq` | 100 | 100 | true |

## Complete Metric Example

```json
{
  "id": "ops-availability",
  "name": "Service Availability",
  "description": "Percentage of time the service is available",
  "domain": "operations",
  "stage": "runtime",
  "category": "reliability",
  "metricType": "rate",
  "trendDirection": "higher_better",
  "unit": "%",
  "baseline": 99.0,
  "current": 99.95,
  "target": 99.99,
  "sli": {
    "name": "Availability",
    "description": "Successful requests / total requests",
    "formula": "1 - (error_count / total_requests)"
  },
  "slo": {
    "target": ">=99.99%",
    "operator": "gte",
    "value": 99.99,
    "window": "30d"
  },
  "thresholds": {
    "green": 99.95,
    "yellow": 99.9,
    "red": 99.0
  },
  "frameworkMappings": [
    {"framework": "SRE", "reference": "availability-slo"},
    {"framework": "DORA", "reference": "availability"}
  ]
}
```

## Qualitative Metric Example

For metrics that track existence or state rather than numeric values:

```json
{
  "id": "sec-incident-runbooks",
  "name": "Incident Runbooks",
  "description": "Documented runbooks for incident response",
  "domain": "security",
  "stage": "response",
  "category": "response",
  "metricType": "boolean",
  "slo": {
    "type": "qualitative",
    "target": "Documented",
    "operator": "exists",
    "status": "documented"
  }
}
```

## XLSX Export

When exporting to XLSX, qualitative metrics display differently:

| Metric | Type | Target | Current | Met |
|--------|------|--------|---------|-----|
| Service Availability | Quantitative | >=99.99% | 99.95 | No |
| Incident Runbooks | Qualitative | Tracked | Documented | Yes |
| Alerting Coverage | Qualitative | Tracked | Partial | No |

## Best Practices

1. **Set Realistic Targets** - SLOs should be achievable but challenging
2. **Include Error Budgets** - Use thresholds to define acceptable ranges
3. **Document Formulas** - Include SLI formulas for clarity
4. **Use Machine-Evaluable Operators** - Enable automated SLO checking
5. **Define Measurement Windows** - Clarify the evaluation period
6. **Map to Frameworks** - Reference industry standards (SRE, DORA)
7. **Use Qualitative for Binary States** - Track existence/compliance with qualitative SLOs
8. **Combine Both Types** - Mix quantitative metrics with qualitative process controls
9. **Define SLIs Separately** - In maturity models, define SLIs with framework mappings once, then reference them from criteria

## Maturity Model: SLIs and Criteria

In PRISM maturity models, criteria (SLOs) can reference shared SLIs:

### SLI Definition

Define SLIs at the top level with framework mappings:

```json
{
  "slis": {
    "security-mttr": {
      "id": "security-mttr",
      "name": "Security MTTR",
      "description": "Mean time to remediate critical security findings",
      "metricName": "security_mttr_days",
      "unit": "days",
      "type": "quantitative",
      "frameworkMappings": [
        {"framework": "NIST_800_53", "reference": "IR-6"},
        {"framework": "SOC_2", "reference": "CC7.4"}
      ]
    }
  }
}
```

### Criterion (SLO) Referencing SLI

Criteria at each maturity level reference the SLI via `sliId`:

```json
{
  "domains": {
    "security": {
      "levels": [
        {
          "level": 4,
          "name": "Managed",
          "criteria": [
            {
              "id": "sec-m4-mttr",
              "name": "Fast Security MTTR",
              "sliId": "security-mttr",
              "operator": "lte",
              "target": 7
            }
          ]
        },
        {
          "level": 5,
          "name": "Optimizing",
          "criteria": [
            {
              "id": "sec-m5-mttr",
              "name": "Rapid Security MTTR",
              "sliId": "security-mttr",
              "operator": "lte",
              "target": 1
            }
          ]
        }
      ]
    }
  }
}
```

### How It Works

1. **SLI defined once** - The metric "security-mttr" with framework mappings is defined in the `slis` section
2. **Criteria reference SLI** - M4 and M5 both use `sliId: "security-mttr"` with different targets
3. **Framework mappings inherited** - Both criteria inherit the NIST/SOC2 mappings from the SLI
4. **Exports show mappings** - XLSX and Markdown exports display framework columns for each criterion

### Criterion Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier |
| `name` | string | Yes | Criterion name |
| `description` | string | No | Description |
| `sliId` | string | No | Reference to SLI (preferred) |
| `operator` | string | Yes | Comparison operator (gte, lte, gt, lt, eq, exists) |
| `target` | number | Yes | Target threshold value |
| `type` | string | No | `quantitative` (default) or `qualitative` |
| `metricName` | string | No | Inline metric name (if no sliId) |
| `unit` | string | No | Inline unit (if no sliId) |
| `frameworkMappings` | array | No | Inline mappings (if no sliId) |

### Backward Compatibility

Criteria can still define fields inline without referencing an SLI:

```json
{
  "criteria": [
    {
      "id": "sec-m4-mttr",
      "name": "Security MTTR",
      "metricName": "security_mttr_days",
      "unit": "days",
      "operator": "lte",
      "target": 7,
      "frameworkMappings": [
        {"framework": "NIST_800_53", "reference": "IR-6"}
      ]
    }
  ]
}
```

Inline `frameworkMappings` on a criterion take precedence over SLI mappings when both are present.

### Resolution Order

When resolving framework mappings for a criterion:

1. If the criterion has `frameworkMappings`, use them
2. Otherwise, look up the SLI via `sliId` and use its `frameworkMappings`
3. If neither exists, no framework mappings are shown
