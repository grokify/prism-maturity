# Developer Guide

This section provides documentation for developers who want to contribute to PRISM or extend its functionality.

## Overview

PRISM is written in Go and follows standard Go project conventions. The codebase is organized into:

- **Root package (`prism`)** - Core types and document handling
- **Internal packages** - Specialized functionality
- **CLI (`cmd/prism`)** - Command-line interface

## Quick Start

### Prerequisites

- Go 1.22 or later
- Git

### Clone and Build

```bash
git clone https://github.com/grokify/prism-intelligence.git
cd prism
go build ./cmd/prism
```

### Run Tests

```bash
go test -v ./...
```

### Lint Code

```bash
golangci-lint run
```

## Project Structure

```
prism/
├── cmd/prism/           # CLI commands
├── analysis/            # Document analysis and gap detection
├── uiforge/           # Dashboard generation
├── export/              # Format converters (OKR, V2MOM, Roadmap)
├── maturity/            # Maturity model types and evaluation
├── output/              # Output formatting utilities
├── report/              # Report generation
├── scaffold/            # Document templates
├── schema/              # JSON Schema generation
├── docs/                # Documentation (MkDocs)
├── examples/            # Example PRISM documents
└── maturity-models/     # Example maturity model specs
```

## Key Concepts for Contributors

### Document Structure

The `PRISMDocument` struct is the central type:

```go
type PRISMDocument struct {
    Schema      string        `json:"$schema,omitempty"`
    Metadata    *Metadata     `json:"metadata,omitempty"`
    Domains     []DomainDef   `json:"domains,omitempty"`
    Maturity    *MaturitySet  `json:"maturity,omitempty"`
    Metrics     []Metric      `json:"metrics"`
    Goals       []Goal        `json:"goals,omitempty"`
    Phases      []Phase       `json:"phases,omitempty"`
    Initiatives []Initiative  `json:"initiatives,omitempty"`
    Layers      []Layer       `json:"layers,omitempty"`
    Teams       []Team        `json:"teams,omitempty"`
    Services    []Service     `json:"services,omitempty"`
    Awareness   *Awareness    `json:"awareness,omitempty"`
}
```

### Adding a New CLI Command

1. Create a new file in `cmd/prism/` (e.g., `mycommand.go`)
2. Define the command using Cobra:

```go
var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Brief description",
    Long:  "Longer description",
    RunE:  runMyCommand,
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

3. Add documentation in `docs/cli/mycommand.md`
4. Update `mkdocs.yml` navigation

### Adding Export Formats

Export converters live in `export/`. To add a new format:

1. Create a new file (e.g., `export/myformat.go`)
2. Implement conversion from `*prism.PRISMDocument`
3. Add a CLI subcommand under `export`

## Documentation

- [Architecture](architecture.md) - Package details and design decisions
- [Contributing](contributing.md) - Contribution guidelines and code style

## Getting Help

- [GitHub Issues](https://github.com/grokify/prism-intelligence/issues) - Bug reports and feature requests
- [GitHub Discussions](https://github.com/grokify/prism-intelligence/discussions) - Questions and discussions
