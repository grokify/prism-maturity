# Schema Overview

PRISM uses a JSON schema to define the structure of documents. The schema is auto-generated from Go types, ensuring consistency between the library and documentation.

## Document Structure

```json
{
  "$schema": "https://github.com/grokify/prism/schema/prism.schema.json",
  "metadata": { ... },
  "metrics": [ ... ],
  "maturity": { ... }
}
```

## Top-Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `$schema` | string | No | JSON Schema URL for validation |
| `metadata` | object | No | Document metadata |
| `metrics` | array | Yes | List of metrics |
| `maturity` | object | No | Maturity model configuration |
| `okrs` | array | No | OKR mappings |
| `initiatives` | array | No | Related initiatives |

## Classification Hierarchy

PRISM organizes metrics in a three-level hierarchy:

```
Domain (security, operations)
  └── Stage (design, build, test, runtime, response)
       └── Category (prevention, detection, response, reliability, efficiency, quality)
```

## Schema Files

| File | Description |
|------|-------------|
| `schema/prism.schema.json` | Full JSON Schema |
| `schema/embed.go` | Go embed directives |
| `schema/generate.go` | Schema generator |

## Generating the Schema

The schema is auto-generated from Go types:

```bash
cd schema
go run generate.go
```

## Using the Schema

### In JSON Files

```json
{
  "$schema": "https://github.com/grokify/prism/schema/prism.schema.json",
  "metrics": [...]
}
```

### In Go Code

```go
import "github.com/grokify/prism/schema"

schemaJSON := schema.PRISMSchemaJSON()
```

## Related Pages

- [Domains](domains.md) - Security and operations domains
- [Lifecycle Stages](stages.md) - Software delivery stages
- [Categories](categories.md) - Metric categories
- [Metric Types](metrics.md) - Types of measurements
- [SLIs & SLOs](slos.md) - Service level indicators and objectives (including maturity model SLIs)
- [Thresholds](thresholds.md) - Status thresholds
