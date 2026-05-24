# Layer-Based Views

Layer views aggregate maturity metrics across the value stream, providing visibility into how different parts of your system are performing.

## Overview

Layer views group Service Level Indicators (SLIs) by their assigned layer in the value stream. This enables:

- **Clear accountability** - See which layers need attention
- **Trend analysis** - Track maturity progress per layer over time
- **Reporting** - Export layer-specific maturity reports

## Enabling Layer Views

Layer views are automatically enabled when:

1. **Capability stacks include layer definitions** - Layers are defined in the stack with unique IDs
2. **SLIs reference layers** - Each SLI in the maturity model specifies a `layer` field
3. **Capabilities map to PRISM SLIs** - Capabilities reference SLI IDs via `prism.sliIds`

### Example Layer Definition

```json
{
  "layers": [
    { "id": "code", "name": "Code", "order": 1 },
    { "id": "infra", "name": "Infrastructure", "order": 2 },
    { "id": "runtime", "name": "Runtime", "order": 3 }
  ]
}
```

### Example SLI with Layer

```json
{
  "id": "availability",
  "name": "Service Availability",
  "layer": "runtime",
  "category": "reliability"
}
```

## Aggregation Methods

When displaying maturity at the layer level, individual SLI levels are aggregated using one of two methods:

| Method | Description | CLI Flag |
|--------|-------------|----------|
| `min` | Minimum SLI level (conservative, default) | `--aggregation=min` |
| `avg` | Average SLI level (balanced view) | `--aggregation=avg` |

### Minimum Aggregation

The `min` method returns the lowest maturity level among all SLIs in a layer. This is the conservative approach:

- **Use case**: Ensuring all metrics meet a baseline before claiming a maturity level
- **Interpretation**: "The layer is only as mature as its weakest metric"

```
Layer: Runtime
  - Availability: M4
  - Latency: M3
  - Error Rate: M4

Aggregate (min): M3
```

### Average Aggregation

The `avg` method calculates the mean maturity level across all SLIs:

- **Use case**: Tracking overall progress when some metrics lag behind
- **Interpretation**: "The layer's average maturity across all metrics"

```
Layer: Runtime
  - Availability: M4
  - Latency: M3
  - Error Rate: M4

Aggregate (avg): M3.67
```

## PDF Export

Layer views include PDF export functionality for reporting and documentation.

### Single Layer Export

Click the **PDF** button on any layer header to export that layer's metrics:

- Layer name and generation date
- Table with all metrics, current levels, and M1-M5 thresholds

### Export All Layers

Click **Export All Layers** to generate a comprehensive report:

- **Title page** with summary table showing all layers
- **Per-layer pages** with detailed metric tables

### PDF Contents

Each PDF includes:

| Column | Description |
|--------|-------------|
| Metric | SLI name |
| Current | Current maturity level (M1-M5) |
| M1-M5 | Threshold values for each maturity level |

## CLI Usage

Generate a site with layer views:

```bash
prism site generate \
  --stack=./security/ \
  --aggregation=min \
  --output=./dist
```

Switch to average aggregation:

```bash
prism site generate \
  --stack=./security/ \
  --aggregation=avg \
  --output=./dist
```

## Best Practices

1. **Define layers consistently** - Use the same layer IDs across all capability stacks
2. **Assign all SLIs to layers** - Ensure every SLI has a `layer` field
3. **Choose aggregation based on use case**:
   - Use `min` for compliance and readiness assessments
   - Use `avg` for progress tracking and trend analysis
4. **Export regularly** - Generate PDFs for quarterly reviews and stakeholder reports

## See Also

- [prism site generate](../cli/site-generate.md) - CLI reference
- [Layers](layers.md) - Layer concept documentation
- [Maturity Model](maturity.md) - Maturity model documentation
