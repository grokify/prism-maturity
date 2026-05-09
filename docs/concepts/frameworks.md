# Framework Mappings

PRISM metrics and SLOs can be mapped to industry frameworks to provide context and enable compliance reporting. This allows you to view your maturity model through the lens of any compliance framework.

## Supported Frameworks

### NIST Frameworks

| Framework | Constant | Description |
|-----------|----------|-------------|
| NIST CSF 1.1 | `NIST_CSF` | NIST Cybersecurity Framework 1.1 |
| NIST CSF 2.0 | `NIST_CSF_2` | NIST Cybersecurity Framework 2.0 |
| NIST 800-53 | `NIST_800_53` | Security and Privacy Controls |
| NIST RMF | `NIST_RMF` | Risk Management Framework |
| NIST AI RMF | `NIST_AI_RMF` | AI Risk Management Framework |
| NIST 800-171 | `NIST_800_171` | CUI Protection |

### FedRAMP (based on NIST 800-53)

| Framework | Constant | Description |
|-----------|----------|-------------|
| FedRAMP | `FEDRAMP` | General FedRAMP |
| FedRAMP High | `FEDRAMP_HIGH` | High impact baseline |
| FedRAMP Moderate | `FEDRAMP_MOD` | Moderate impact baseline |
| FedRAMP Low | `FEDRAMP_LOW` | Low impact baseline |

### Other Security Frameworks

| Framework | Constant | Description |
|-----------|----------|-------------|
| MITRE ATT&CK | `MITRE_ATTACK` | Threat framework |
| CIS Controls | `CIS_CONTROLS` | Critical Security Controls |
| SOC 2 | `SOC_2` | Trust Services Criteria |
| ISO 27001 | `ISO_27001` | Information Security Management |

### Engineering Frameworks

| Framework | Constant | Description |
|-----------|----------|-------------|
| DORA | `DORA` | DevOps Research and Assessment |
| SRE | `SRE` | Site Reliability Engineering |

## Framework Mapping Structure

Metrics and SLOs can have multiple framework mappings:

```json
{
  "frameworkMappings": [
    {
      "framework": "NIST_CSF_2",
      "reference": "PR.DS-01",
      "name": "Data-at-rest is protected",
      "version": "2.0"
    },
    {
      "framework": "NIST_800_53",
      "reference": "SC-28",
      "name": "Protection of Information at Rest",
      "baseline": "moderate"
    },
    {
      "framework": "FEDRAMP_MOD",
      "reference": "SC-28",
      "name": "Protection of Information at Rest"
    }
  ]
}
```

### Framework Mapping Fields

| Field | Description | Example |
|-------|-------------|---------|
| `framework` | Framework identifier | `NIST_800_53` |
| `reference` | Control/function ID | `AC-2`, `PR.DS-01` |
| `name` | Human-readable control name | `Account Management` |
| `description` | Control description | Optional details |
| `baseline` | Required baseline level | `high`, `moderate`, `low` |
| `version` | Framework version | `2.0`, `Rev 5` |

## SLO Framework Mappings

SLOs can be mapped to framework controls to show how objectives satisfy compliance requirements:

```json
{
  "metrics": [{
    "id": "encryption-at-rest",
    "name": "Data Encryption at Rest",
    "slo": {
      "id": "slo-encryption",
      "name": "100% encryption coverage",
      "type": "quantitative",
      "target": ">=100%",
      "operator": "gte",
      "value": 100,
      "frameworkMappings": [
        {
          "framework": "NIST_800_53",
          "reference": "SC-28",
          "name": "Protection of Information at Rest",
          "baseline": "moderate"
        },
        {
          "framework": "FEDRAMP_HIGH",
          "reference": "SC-28",
          "name": "Protection of Information at Rest"
        }
      ]
    }
  }]
}
```

## Maturity Model Framework Mappings

In maturity models, framework mappings should be defined at the **SLI level** (the metric being measured), not repeated on each criterion (SLO) per maturity level. This avoids duplication and ensures consistent mappings.

### Recommended: SLI-Level Mappings

Define the SLI with framework mappings once:

```json
{
  "slis": {
    "data-encryption": {
      "id": "data-encryption",
      "name": "Data Encryption Coverage",
      "description": "Percentage of data encrypted at rest",
      "metricName": "encryption_coverage_pct",
      "unit": "%",
      "type": "quantitative",
      "frameworkMappings": [
        {"framework": "NIST_CSF_2", "reference": "PR.DS-01", "name": "Data-at-rest is protected"},
        {"framework": "NIST_800_53", "reference": "SC-28", "baseline": "moderate"}
      ]
    }
  },
  "domains": {
    "security": {
      "levels": [
        {
          "level": 3,
          "criteria": [{
            "id": "sec-m3-encryption",
            "name": "Basic Encryption",
            "sliId": "data-encryption",
            "operator": "gte",
            "target": 80
          }]
        },
        {
          "level": 4,
          "criteria": [{
            "id": "sec-m4-encryption",
            "name": "Full Encryption",
            "sliId": "data-encryption",
            "operator": "gte",
            "target": 100
          }]
        }
      ]
    }
  }
}
```

Both M3 and M4 criteria inherit the NIST CSF and NIST 800-53 mappings from the SLI.

### Legacy: Inline Criterion Mappings

For backward compatibility, criteria can still have inline mappings:

```json
{
  "domains": {
    "security": {
      "levels": [{
        "level": 3,
        "criteria": [{
          "id": "sec-encryption",
          "name": "Data Encryption",
          "type": "qualitative",
          "operator": "exists",
          "status": "implemented",
          "frameworkMappings": [
            {"framework": "NIST_CSF_2", "reference": "PR.DS-01"},
            {"framework": "NIST_800_53", "reference": "SC-28", "baseline": "moderate"}
          ]
        }]
      }]
    }
  }
}
```

!!! note "Resolution Order"
    When both are present, inline `frameworkMappings` on the criterion take precedence over SLI mappings.

## DORA Metrics

DORA (DevOps Research and Assessment) defines four key metrics for software delivery performance.

### DORA Metrics Mapping

| DORA Metric | PRISM Domain | Stage | Category |
|-------------|--------------|-------|----------|
| Deployment Frequency | Operations | Build | Efficiency |
| Lead Time for Changes | Operations | Build | Efficiency |
| Mean Time to Recovery | Operations | Response | Reliability |
| Change Failure Rate | Operations | Build | Quality |

### Example Document

```json
{
  "metrics": [
    {
      "id": "dora-deploy-frequency",
      "name": "Deployment Frequency",
      "domain": "operations",
      "stage": "build",
      "category": "efficiency",
      "metricType": "rate",
      "unit": "deploys/day",
      "current": 5,
      "target": 10,
      "frameworkMappings": [
        {"framework": "DORA", "reference": "deployment-frequency"}
      ]
    },
    {
      "id": "dora-lead-time",
      "name": "Lead Time for Changes",
      "domain": "operations",
      "stage": "build",
      "category": "efficiency",
      "metricType": "latency",
      "unit": "hours",
      "current": 24,
      "target": 1,
      "frameworkMappings": [
        {"framework": "DORA", "reference": "lead-time"}
      ]
    },
    {
      "id": "dora-mttr",
      "name": "Mean Time to Recovery",
      "domain": "operations",
      "stage": "response",
      "category": "reliability",
      "metricType": "latency",
      "unit": "hours",
      "current": 4,
      "target": 1,
      "frameworkMappings": [
        {"framework": "DORA", "reference": "mttr"}
      ]
    },
    {
      "id": "dora-change-failure",
      "name": "Change Failure Rate",
      "domain": "operations",
      "stage": "build",
      "category": "quality",
      "metricType": "rate",
      "unit": "%",
      "current": 10,
      "target": 5,
      "frameworkMappings": [
        {"framework": "DORA", "reference": "change-failure-rate"}
      ]
    }
  ]
}
```

### DORA Performance Levels

| Metric | Elite | High | Medium | Low |
|--------|-------|------|--------|-----|
| Deploy Frequency | On-demand | 1/day-1/week | 1/week-1/month | 1/month-6/month |
| Lead Time | <1 hour | 1 day-1 week | 1 week-1 month | 1-6 months |
| MTTR | <1 hour | <1 day | <1 day | 1 week-1 month |
| Change Failure | 0-15% | 16-30% | 16-30% | 16-30% |

## SRE (Site Reliability Engineering)

SRE practices focus on reliability through SLIs, SLOs, and error budgets.

### SRE Concepts in PRISM

| SRE Concept | PRISM Field | Description |
|-------------|-------------|-------------|
| SLI | `sli` | Service Level Indicator |
| SLO | `slo` | Service Level Objective |
| Error Budget | `thresholds` | Acceptable variance |

### Example SRE Mapping

```json
{
  "id": "sre-availability",
  "name": "Service Availability",
  "domain": "operations",
  "stage": "runtime",
  "category": "reliability",
  "sli": {
    "name": "Availability",
    "formula": "successful_requests / total_requests"
  },
  "slo": {
    "target": ">=99.99%",
    "operator": "gte",
    "value": 99.99,
    "window": "30d"
  },
  "frameworkMappings": [
    {"framework": "SRE", "reference": "availability-slo"}
  ]
}
```

## Using Framework Mappings

### Query by Framework

```go
// Find all metrics mapped to DORA
for _, metric := range doc.Metrics {
    for _, mapping := range metric.FrameworkMappings {
        if mapping.Framework == prism.FrameworkDORA {
            fmt.Printf("%s → %s\n", metric.Name, mapping.Reference)
        }
    }
}
```

### Generate Compliance Reports

Framework mappings enable automated compliance reporting:

1. Extract metrics by framework
2. Calculate coverage per framework category
3. Identify gaps in framework coverage
4. Generate framework-specific reports

## Framework-Based Maturity View

With framework mappings on SLOs and criteria, you can view your maturity model through the lens of any compliance framework:

### Example: NIST CSF 2.0 View

```
NIST CSF 2.0 Control Coverage
=============================
GOVERN (GV)
  GV.OC-01: Organizational Context    [M3] ✓ Implemented
  GV.RM-01: Risk Management Strategy  [M2] ✓ Implemented

IDENTIFY (ID)
  ID.AM-01: Asset Inventory           [M4] ✓ Implemented
  ID.RA-01: Risk Assessment           [M3] ○ Partial

PROTECT (PR)
  PR.DS-01: Data Protection           [M4] ✓ Implemented
  PR.AA-01: Identity Management       [M3] ✓ Implemented

DETECT (DE)
  DE.CM-01: Network Monitoring        [M3] ✓ Implemented
  DE.AE-01: Adverse Event Analysis    [M2] ○ Partial

RESPOND (RS)
  RS.MA-01: Incident Management       [M3] ✓ Implemented
  RS.CO-01: Incident Communications   [M2] ○ Planned

RECOVER (RC)
  RC.RP-01: Recovery Planning         [M2] ✓ Implemented
```

### Example: FedRAMP Moderate View

```
FedRAMP Moderate Control Status
================================
Access Control (AC)
  AC-2  Account Management           ✓ Implemented  [SLO: slo-iam-coverage]
  AC-3  Access Enforcement           ✓ Implemented  [SLO: slo-rbac]
  AC-6  Least Privilege              ○ Partial      [SLO: slo-privilege]

System & Comms Protection (SC)
  SC-7  Boundary Protection          ✓ Implemented  [SLO: slo-network-seg]
  SC-8  Transmission Confidentiality ✓ Implemented  [SLO: slo-tls]
  SC-28 Protection at Rest           ✓ Implemented  [SLO: slo-encryption]

Summary: 85% of Moderate controls mapped
         92% of mapped controls satisfied
```

## NIST CSF 2.0 Functions

The NIST Cybersecurity Framework 2.0 organizes controls into six functions:

| Function | Code | Description |
|----------|------|-------------|
| Govern | GV | Organizational context and risk strategy |
| Identify | ID | Asset and risk identification |
| Protect | PR | Safeguards and access control |
| Detect | DE | Monitoring and anomaly detection |
| Respond | RS | Incident response |
| Recover | RC | Recovery planning and improvements |

## NIST 800-53 Control Families

Common NIST 800-53 control families:

| Family | Code | Description |
|--------|------|-------------|
| Access Control | AC | Identity and access management |
| Audit & Accountability | AU | Logging and monitoring |
| Configuration Management | CM | System configuration |
| Identification & Auth | IA | Authentication controls |
| Incident Response | IR | Incident handling |
| System & Comms Protection | SC | Encryption, network security |
| System & Info Integrity | SI | Malware protection, patching |

## Best Practices

1. **Map to multiple frameworks** - A single control often satisfies multiple frameworks
2. **Include baselines** - Specify FedRAMP High/Moderate/Low requirements
3. **Link SLOs to controls** - Show how SLOs provide evidence for compliance
4. **Track coverage** - Calculate percentage of framework controls addressed
5. **Use qualitative SLOs** - Many compliance controls are binary (implemented/not)

