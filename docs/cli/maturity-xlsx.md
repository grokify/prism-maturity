# prism maturity model xlsx

Generate an Excel (XLSX) report from a maturity model specification.

## Synopsis

```bash
prism maturity model xlsx <model-file> [flags]
```

> **Note:** As of v0.6.0, `prism maturity xlsx` has been moved to `prism maturity model xlsx`.

## Description

The `prism maturity model xlsx` command generates a comprehensive Excel workbook from a PRISM maturity specification. The workbook contains multiple sheets for different views of the maturity data.

## Output Sheets

| Sheet | Description |
|-------|-------------|
| **Requirements** | Enablers with domain, level, type, team, effort, and status |
| **SLOs** | Criteria (SLOs) with framework columns showing control references |
| **Framework Mappings** | Detailed per-mapping rows with framework, reference, control name, baseline |
| **Progress** | Assessment status by domain with level progress percentages |
| **Level Definitions** | M1-M5 level descriptions for each domain |

## Framework Columns

The SLOs sheet dynamically adds columns for each compliance framework found in the maturity spec. For example, if criteria map to NIST CSF 2.0 and SOC 2, the sheet will include columns:

| ID | Domain | Level | Name | ... | NIST_CSF_2 | SOC_2 |
|----|--------|-------|------|-----|------------|-------|
| SEC-001 | Security | M2 | Asset Inventory | ... | ID.AM-1 | CC6.1 |
| SEC-002 | Security | M2 | Vuln Scanning | ... | ID.RA-1 | - |

The `-` indicates the criterion has no mapping to that framework.

### SLI Resolution

Framework mappings can be defined in two places:

1. **SLI level** (recommended) - Defined once in the `slis` section, referenced via `sliId`
2. **Criterion level** (legacy) - Defined inline on each criterion

The XLSX generator automatically resolves mappings from the SLI when a criterion uses `sliId`. Inline mappings on the criterion take precedence if both are present.

See [SLIs & SLOs](../schema/slos.md) for details on the SLI architecture.

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output file path (default: input filename with .xlsx extension) |
| `--help` | `-h` | Help for xlsx command |

## Examples

```bash
# Generate XLSX with default filename (model.xlsx)
prism maturity model xlsx model.json

# Generate XLSX with custom output filename
prism maturity model xlsx model.json -o security-maturity-report.xlsx

# Generate from YAML input
prism maturity model xlsx maturity.yaml -o report.xlsx
```

## Output File

If no output file is specified, the command generates a file with the same name as the input but with a `.xlsx` extension:

| Input | Output |
|-------|--------|
| `spec.json` | `spec.xlsx` |
| `maturity.yaml` | `maturity.xlsx` |
| `data/model.json` | `data/model.xlsx` |

## Color Coding

The generated workbook uses color coding for status fields:

**Requirements Sheet (Enabler Status):**

| Status | Color |
|--------|-------|
| Completed | Green |
| In Progress | Yellow |
| Blocked | Red |
| Not Started | Gray |

**SLOs Sheet (Met Status):**

| Status | Color |
|--------|-------|
| Yes (Met) | Green |
| No (Pending) | Yellow |

**Progress Sheet:**

| Progress | Color |
|----------|-------|
| 80-100% | Green |
| 50-79% | Yellow |
| 0-49% | Red |

## See Also

- [`prism maturity model report`](maturity-report.md) - Generate Markdown reports
- [`prism slo-report`](slo-report.md) - Generate SLO compliance reports
