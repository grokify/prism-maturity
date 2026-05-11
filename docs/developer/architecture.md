# Architecture

This document describes the internal architecture and package structure of PRISM.

## Package Overview

### Core Package (`github.com/grokify/prism`)

The root package contains all core types and primary document operations:

| File | Purpose |
|------|---------|
| `prism.go` | `PRISMDocument` type and document operations |
| `constants.go` | Domain, stage, category, and other constants |
| `score.go` | PRISM score calculation |
| `validation.go` | Document validation |
| `goal.go` | Goal types and maturity level calculations |
| `phase.go` | Phase management and timeline operations |
| `layer.go` | Value stream layer definitions |
| `team.go` | Team topology types |
| `service.go` | Service definitions and ownership |
| `awareness.go` | Customer awareness scoring |
| `maturity.go` | Core maturity model types |

### Analysis Package (`analysis/`)

Provides document analysis and gap detection:

```go
// Key types
type Result struct {
    Summary         Summary
    Goals           []GoalAnalysis
    Phases          []PhaseAnalysis
    Gaps            []Gap
    Recommendations []Recommendation
}

// Primary function
func Analyze(doc *prism.PRISMDocument) *Result
```

**Files:**

- `analyzer.go` - Main analysis logic
- `gaps.go` - Gap identification (maturity gaps, SLO gaps, coverage gaps)
- `roadmap.go` - Roadmap analysis and recommendations

### Export Package (`export/`)

Converts PRISM documents to structured-plan formats:

```go
// OKR export
func ConvertToOKR(doc *prism.PRISMDocument) *OKRDocument

// V2MOM export
func ConvertToV2MOM(doc *prism.PRISMDocument) *V2MOMDocument

// Roadmap export
func ConvertToRoadmap(doc *prism.PRISMDocument) *RoadmapDocument
```

**Integration with structured-plan:**

PRISM exports to [structured-plan](https://github.com/grokify/structured-plan) formats:

| PRISM Concept | Structured-Plan Concept |
|---------------|-------------------------|
| Goal | OKR Objective |
| Goal.TargetLevel | Objective Target |
| SLO (per level) | Key Result |
| Phase.GoalTargets | PhaseTargets in Key Results |
| Initiative | Deliverable (in roadmap phase) |

### Maturity Package (`maturity/`)

Handles maturity model specifications and evaluation using a three-part schema:

| Schema | Type | Purpose |
|--------|------|---------|
| `prism-maturity-model` | `Spec` | Defines what good looks like (M1-M5 levels) |
| `prism-maturity-state` | `MaturityStateDocument` | Tracks current state and measurements |
| `prism-maturity-plan` | `PRISMDocument` | Plans how to achieve target levels |

```go
// Spec defines a maturity model specification
type Spec struct {
    Schema   string                  `json:"$schema,omitempty"`
    Metadata *SpecMetadata           `json:"metadata,omitempty"`
    SLIs     map[string]*SLI         `json:"slis,omitempty"`
    Domains  map[string]*DomainModel `json:"domains"`
}

// DomainModel defines maturity levels for a domain
type DomainModel struct {
    Name        string
    Description string
    Levels      []Level
}

// Level defines a maturity level (M1-M5)
type Level struct {
    Level       int
    Name        string
    Criteria    []Criterion  // SLOs that define the level
    Enablers    []Enabler    // Tasks to achieve the level
}
```

**Key operations:**

- `ReadSpecFile()` - Load maturity specs from JSON
- `Level.IsLevelAchieved()` - Check if criteria are met
- `Level.CalculateLevelProgress()` - Calculate progress toward a level

**State tracking:**

State is tracked separately from the model using `MaturityStateDocument`:

```go
type MaturityStateDocument struct {
    Schema        string                `json:"$schema,omitempty"`
    Metadata      MaturityStateMetadata `json:"metadata"`
    SLIState      SLIStateMap           `json:"sliState,omitempty"`
    MaturityState MaturityStateMap      `json:"maturityState,omitempty"`
    EnablerState  EnablerStateMap       `json:"enablerState,omitempty"`
}
```

### Output Package (`output/`)

Provides multi-format output capabilities:

```go
// Supported formats
const (
    FormatText     Format = "text"
    FormatJSON     Format = "json"
    FormatMarkdown Format = "markdown"
    FormatTOON     Format = "toon"  // Token-Optimized Object Notation
)

// Formatter handles output
type Formatter struct {
    Format Format
    Writer io.Writer
}

// Output methods
func (f *Formatter) WriteTable(data *TableData) error
func (f *Formatter) WriteDetail(data *DetailData) error
func (f *Formatter) WriteJSON(data interface{}) error
```

**TOON Format:**

TOON (Token-Optimized Object Notation) is a compact format designed for LLM token efficiency:

```
title;h1,h2,h3;r1c1,r1c2,r1c3;r2c1,r2c2,r2c3
```

### Scaffold Package (`scaffold/`)

Provides templates for creating new PRISM documents:

```go
// Create a new document with specified domains
func NewDocument(domains ...string) *prism.PRISMDocument

// Get example metrics
func OperationsMetrics() []prism.Metric
func SecurityMetrics() []prism.Metric
```

### Dashforge Package (`dashforge/`)

Converts PRISM data to dashboard format:

```go
// Convert to Dashforge format
func Convert(doc *prism.PRISMDocument) *dashforge.Dashboard

// Generate specific chart types
func GaugeChart(score float64, title string) *echartify.Chart
func RadarChart(scores map[string]float64) *echartify.Chart
```

### Report Package (`report/`)

Generates various report formats:

- `exec_dashboard.go` - Executive dashboard data
- `markdown.go` - Markdown report generation

## Data Flow

```
┌─────────────────────────────────────────────────────────────┐
│                     prism.json (Input)                      │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                   PRISMDocument (Core)                       │
│  - Validate()                                                │
│  - CalculatePRISMScore()                                    │
│  - GetMetricByID(), GetGoalByID(), etc.                     │
└─────────────────────────┬───────────────────────────────────┘
                          │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  analysis/  │  │   export/   │  │  dashforge/ │
│             │  │             │  │             │
│ - Analyze() │  │ - ToOKR()   │  │ - Convert() │
│ - Gaps      │  │ - ToV2MOM() │  │ - Charts    │
└─────────────┘  └─────────────┘  └─────────────┘
         │               │               │
         ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────┐
│                    output/ (Formatter)                       │
│  - text, json, markdown, toon                                │
└─────────────────────────────────────────────────────────────┘
```

## Design Principles

### 1. Go-First Schema

JSON Schema is generated from Go types, not hand-written:

```bash
cd schema && go run generate.go
```

This ensures the schema always matches the Go types.

### 2. Sparse Data Handling

PRISM handles incomplete data gracefully:

- Metrics without SLOs are still valid
- Goals can have partial maturity models
- Missing values don't break calculations

### 3. Multi-Domain Architecture

Metrics are organized by orthogonal dimensions:

- **Domain**: operations, security, quality
- **Layer**: requirements, code, infra, runtime, adoption, support
- **Stage**: design, build, test, runtime, response
- **Category**: prevention, detection, response, reliability, efficiency

### 4. Framework Agnostic

PRISM defines its own model but maps to external frameworks:

- DORA (DevOps Research and Assessment)
- SRE (Site Reliability Engineering)
- NIST CSF (Cybersecurity Framework)
- ISO 25010 (Quality Model)

## Testing

### Unit Tests

Each package has corresponding `*_test.go` files:

```bash
go test -v ./...
```

### Integration Tests

CLI integration tests are in `cmd/prism/cli_integration_test.go`:

```bash
go test -v ./cmd/prism/ -run Integration
```

### Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Dependencies

Key external dependencies:

| Dependency | Purpose |
|------------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/xuri/excelize/v2` | Excel export |
| `github.com/grokify/structured-plan` | OKR/V2MOM/Roadmap formats |
| `github.com/plexusone/dashforge` | Dashboard generation |
| `github.com/grokify/echartify` | Chart definitions |
