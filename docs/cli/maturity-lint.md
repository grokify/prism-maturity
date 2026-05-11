# prism maturity model lint

Lint a maturity model document for issues that affect dashboard display and completeness.

## Synopsis

```bash
prism maturity model lint <model-file> [flags]
```

## Description

Unlike `validate` which checks structural correctness, `lint` checks for practical issues that affect dashboard generation and user experience:

- **Criteria without sliId** - Won't appear in dashboard
- **Criteria missing operator or target** - Threshold can't display properly
- **SLIs without any criteria** - Unused SLI definitions
- **SLIs missing unit** - Affects threshold formatting
- **SLIs missing sliType** - Affects methodology grouping in dashboard
- **Incomplete threshold coverage** - Gaps between maturity levels

## Flags

| Flag | Description |
|------|-------------|
| `--strict` | Treat warnings as errors (exit code 2) |
| `-h, --help` | Help for lint |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | No issues found |
| 1 | Warnings found (non-blocking) |
| 2 | Errors found (blocking issues, or warnings with `--strict`) |

## Checks Performed

### Errors

| Check | Description |
|-------|-------------|
| Missing `sliId` | Criteria without an SLI reference won't appear in the dashboard |
| Missing `operator` | Threshold comparison operator required for display |

### Warnings

| Check | Description |
|-------|-------------|
| `target=0` (quantitative) | May be unintentional for non-qualitative metrics |
| Qualitative without `exists` | Qualitative criteria should use `operator: "exists"` |
| Orphan SLIs | SLIs defined but not referenced by any criteria |
| Threshold gaps | Missing levels between defined ones (e.g., M1 and M3 defined, but M2 missing) |

### Info

| Check | Description |
|-------|-------------|
| Missing `unit` | SLI lacks unit field, affects threshold formatting |
| Missing `sliType` | SLI lacks sliType field, affects methodology grouping |

## Examples

### Basic usage

```bash
prism maturity model lint model.json
```

Output:

```
WARNING product_security.M3.ps-secrets_m3
        Criterion has target=0 which may be unintentional for quantitative metrics
        → Set target to the threshold value, or use operator='exists' for qualitative criteria

─────────────────────────────────────────
Lint summary for model.json:
  Errors:   0
  Warnings: 1
  Info:     0
```

### Strict mode for CI/CD

```bash
prism maturity model lint model.json --strict
```

In strict mode, any warning will cause the command to exit with code 2, making it suitable for CI/CD pipelines.

### Lint before dashboard generation

```bash
# Check for issues first
prism maturity model lint model.json

# Generate dashboard if no errors
prism maturity model dashboard model.json --state state.json -f html -o dashboard.html
```

## Common Issues and Fixes

### Missing sliId

**Problem:** Criteria defines a threshold but doesn't reference an SLI.

```json
{
  "id": "avail-m3",
  "description": "99.9% availability",
  "operator": ">=",
  "target": 99.9
}
```

**Fix:** Add the `sliId` field:

```json
{
  "id": "avail-m3",
  "sliId": "sli-availability",
  "description": "99.9% availability",
  "operator": ">=",
  "target": 99.9
}
```

### Qualitative without exists operator

**Problem:** SLI is qualitative but criteria uses numeric operator.

**Fix:** Use `operator: "exists"` for qualitative criteria:

```json
{
  "id": "threat-model-m2",
  "sliId": "sli-threat-model",
  "operator": "exists",
  "target": 1
}
```

### Threshold gaps

**Problem:** SLI has thresholds for M1 and M3, but not M2.

**Fix:** Either add the missing M2 threshold, or document why it's intentionally skipped.

## See Also

- [prism maturity model validate](validate.md) - Structural validation
- [prism maturity model dashboard](dashboard.md) - Dashboard generation
- [v0.6.0 Release Notes](../releases/v0.6.0.md) - Feature introduction
