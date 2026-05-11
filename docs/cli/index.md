# CLI Reference

The PRISM CLI provides commands for creating, validating, and scoring PRISM documents.

## Installation

```bash
go install github.com/grokify/prism/cmd/prism@latest
```

## Commands

### Core Commands

| Command | Description |
|---------|-------------|
| [`prism init`](init.md) | Initialize a new PRISM document |
| [`prism validate`](validate.md) | Validate a PRISM document |
| [`prism score`](score.md) | Calculate the PRISM score |
| [`prism catalog`](catalog.md) | List available constants |

### Roadmap Commands (v0.2.0)

| Command | Description |
|---------|-------------|
| [`prism goal`](goal.md) | Manage and inspect goals |
| [`prism phase`](phase.md) | Manage and inspect phases |
| [`prism roadmap`](roadmap.md) | View roadmap overview |
| [`prism initiative`](initiative.md) | List and inspect initiatives |

### SLO Reporting (v0.2.0)

| Command | Description |
|---------|-------------|
| [`prism slo-report`](slo-report.md) | Generate SLO compliance reports |
| `prism dashforge` | Convert to dashforge format |

### Maturity Commands (v0.6.0)

#### Model Commands

| Command | Description |
|---------|-------------|
| [`prism maturity model report`](maturity-report.md) | Generate maturity model reports (Markdown) |
| [`prism maturity model xlsx`](maturity-xlsx.md) | Generate maturity model reports (Excel) |
| `prism maturity model dashboard` | Generate model dashboards with state integration |
| `prism maturity model validate` | Validate model documents |
| [`prism maturity model lint`](maturity-lint.md) | Check for dashboard display issues and threshold gaps |

#### State Commands

| Command | Description |
|---------|-------------|
| `prism maturity state validate` | Validate state documents (with optional model cross-validation) |
| `prism maturity state show` | Display state summary (text/json) |

#### Plan Commands

| Command | Description |
|---------|-------------|
| [`prism maturity plan dashboard`](dashboard.md) | Generate executive dashboards |
| [`prism maturity plan report`](report.md) | Generate roadmap reports (Markdown/JSON) |

### Organizational Commands (v0.3.0)

| Command | Description |
|---------|-------------|
| [`prism layer`](layer.md) | List and inspect value stream layers |
| [`prism team`](team.md) | List and inspect teams |
| [`prism service`](service.md) | List and inspect services |

### Analysis & Export Commands (v0.3.0)

| Command | Description |
|---------|-------------|
| [`prism analyze`](analyze.md) | Analyze document and generate recommendations |
| [`prism export`](export.md) | Export to OKR/V2MOM formats |

## Global Flags

| Flag | Description |
|------|-------------|
| `--help`, `-h` | Show help for any command |
| `--version` | Show version information |

## Usage

```bash
# Get help
prism --help
prism <command> --help

# Initialize a document
prism init -o prism.json

# Validate a document
prism validate prism.json

# Calculate score
prism score prism.json

# List constants
prism catalog
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Validation errors or general failure |
| 2 | File not found or I/O error |
