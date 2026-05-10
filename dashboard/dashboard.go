// Package dashboard provides HTML dashboard generation for PRISM maturity models.
package dashboard

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/grokify/prism/maturity"
)

// Dashboard represents a Dashforge-compatible dashboard.
type Dashboard struct {
	Schema      string       `json:"$schema,omitempty"`
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Layout      Layout       `json:"layout"`
	DataSources []DataSource `json:"dataSources"`
	Widgets     []Widget     `json:"widgets"`
	Theme       *Theme       `json:"theme,omitempty"`
}

// Layout defines the dashboard grid layout.
type Layout struct {
	Type      string `json:"type"`
	Columns   int    `json:"columns"`
	RowHeight int    `json:"rowHeight"`
	Gap       int    `json:"gap"`
	Padding   int    `json:"padding"`
}

// DataSource defines a data source for widgets.
type DataSource struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data,omitempty"`
	URL    string          `json:"url,omitempty"`
	Format string          `json:"format,omitempty"`
}

// Widget represents a dashboard widget.
type Widget struct {
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Title        string          `json:"title"`
	Position     Position        `json:"position"`
	DataSourceID string          `json:"dataSourceId,omitempty"`
	Config       json.RawMessage `json:"config"`
}

// Position defines widget grid position.
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// Theme defines dashboard visual styling.
type Theme struct {
	Mode            string `json:"mode"`
	PrimaryColor    string `json:"primaryColor"`
	BackgroundColor string `json:"backgroundColor,omitempty"`
}

// MetricConfig configures a metric widget.
type MetricConfig struct {
	ValueField    string           `json:"valueField"`
	Format        string           `json:"format"`
	FormatOptions *FormatOptions   `json:"formatOptions,omitempty"`
	Thresholds    []ThresholdValue `json:"thresholds,omitempty"`
	Subtitle      string           `json:"subtitle,omitempty"`
}

// FormatOptions defines number formatting.
type FormatOptions struct {
	Decimals int    `json:"decimals,omitempty"`
	Prefix   string `json:"prefix,omitempty"`
	Suffix   string `json:"suffix,omitempty"`
}

// ThresholdValue defines a color threshold.
type ThresholdValue struct {
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

// TableConfig configures a table widget.
type TableConfig struct {
	Columns  []TableColumn `json:"columns"`
	Sortable bool          `json:"sortable,omitempty"`
	Striped  bool          `json:"striped,omitempty"`
}

// TableColumn defines a table column.
type TableColumn struct {
	Field  string `json:"field"`
	Header string `json:"header"`
	Width  string `json:"width,omitempty"`
	Align  string `json:"align,omitempty"`
	Format string `json:"format,omitempty"`
}

// ChartConfig configures a chart widget (ChartIR subset).
type ChartConfig struct {
	Marks   []Mark   `json:"marks"`
	Axes    []Axis   `json:"axes,omitempty"`
	Grid    *Grid    `json:"grid,omitempty"`
	Legend  *Legend  `json:"legend,omitempty"`
	Tooltip *Tooltip `json:"tooltip,omitempty"`
}

// Mark represents a chart series.
type Mark struct {
	ID       string `json:"id"`
	Geometry string `json:"geometry"`
	Encode   Encode `json:"encode"`
	Style    *Style `json:"style,omitempty"`
	Stack    string `json:"stack,omitempty"`
	Name     string `json:"name,omitempty"`
}

// Encode maps data to visual channels.
type Encode struct {
	X     string `json:"x,omitempty"`
	Y     string `json:"y,omitempty"`
	Value string `json:"value,omitempty"`
	Name  string `json:"name,omitempty"`
}

// Axis defines a chart axis.
type Axis struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Position string   `json:"position"`
	Name     string   `json:"name,omitempty"`
	Min      *float64 `json:"min,omitempty"`
	Max      *float64 `json:"max,omitempty"`
}

// Grid defines chart container positioning.
type Grid struct {
	Left         string `json:"left,omitempty"`
	Right        string `json:"right,omitempty"`
	Top          string `json:"top,omitempty"`
	Bottom       string `json:"bottom,omitempty"`
	ContainLabel bool   `json:"containLabel,omitempty"`
}

// Legend configures chart legend.
type Legend struct {
	Show     bool   `json:"show"`
	Position string `json:"position,omitempty"`
}

// Tooltip configures chart tooltips.
type Tooltip struct {
	Show    bool   `json:"show"`
	Trigger string `json:"trigger,omitempty"`
}

// Style defines visual styling for marks.
type Style struct {
	Color        string   `json:"color,omitempty"`
	BorderRadius any      `json:"borderRadius,omitempty"`
	BarWidth     any      `json:"barWidth,omitempty"`
	Opacity      *float64 `json:"opacity,omitempty"`
}

// TextConfig configures a text widget.
type TextConfig struct {
	Content   string `json:"content"`
	Format    string `json:"format,omitempty"`
	Variables bool   `json:"variables,omitempty"`
}

// Generator creates dashboards from maturity specs.
type Generator struct {
	spec    *maturity.Spec
	widgets []Widget
	data    []DataSource
	row     int
}

// NewGenerator creates a dashboard generator for a maturity spec.
func NewGenerator(spec *maturity.Spec) *Generator {
	return &Generator{
		spec:    spec,
		widgets: []Widget{},
		data:    []DataSource{},
		row:     0,
	}
}

// Generate creates a complete dashboard.
func (g *Generator) Generate() (*Dashboard, error) {
	if g.spec == nil {
		return nil, fmt.Errorf("spec is nil")
	}

	title := "Maturity Dashboard"
	description := "PRISM Maturity Model Dashboard"
	if g.spec.Metadata != nil {
		if g.spec.Metadata.Name != "" {
			title = g.spec.Metadata.Name + " Dashboard"
		}
		if g.spec.Metadata.Description != "" {
			description = g.spec.Metadata.Description
		}
	}

	// Add header text
	g.addHeaderWidget(title, description)

	// Add domain summary cards
	g.addDomainSummaryRow()

	// Add maturity bullet charts for each domain
	g.addBulletWidgets()

	// Add level progress charts for each domain
	g.addLevelProgressCharts()

	// Add SLI tables by category
	g.addSLITables()

	return &Dashboard{
		Schema:      "https://github.com/plexusone/dashforge/schema/dashboard.schema.json",
		ID:          "prism-maturity-dashboard",
		Title:       title,
		Description: description,
		Layout: Layout{
			Type:      "grid",
			Columns:   12,
			RowHeight: 80,
			Gap:       16,
			Padding:   24,
		},
		DataSources: g.data,
		Widgets:     g.widgets,
		Theme: &Theme{
			Mode:         "light",
			PrimaryColor: "#3b82f6",
		},
	}, nil
}

func (g *Generator) addHeaderWidget(title, description string) {
	content := fmt.Sprintf("# %s\n\n%s", title, description)
	config, _ := json.Marshal(TextConfig{
		Content: content,
		Format:  "markdown",
	})

	g.widgets = append(g.widgets, Widget{
		ID:       "header",
		Type:     "text",
		Title:    "",
		Position: Position{X: 0, Y: g.row, W: 12, H: 1},
		Config:   config,
	})
	g.row++
}

func (g *Generator) addDomainSummaryRow() {
	domains := g.getSortedDomains()
	if len(domains) == 0 {
		return
	}

	width := 12 / len(domains)
	if width < 3 {
		width = 3
	}

	for i, domainKey := range domains {
		domain := g.spec.Domains[domainKey]
		assessment := g.spec.Assessments[domainKey]

		// Create data source for this domain
		dataID := fmt.Sprintf("domain-%s-data", domainKey)
		currentLevel := 1
		targetLevel := 5
		progressPercent := 0.0

		if assessment != nil {
			currentLevel = assessment.CurrentLevel
			targetLevel = assessment.TargetLevel
			if targetLevel > 1 {
				progressPercent = float64(currentLevel-1) / float64(targetLevel-1) * 100
			}
		}

		dataRow := map[string]any{
			"domain":       domain.Name,
			"current":      currentLevel,
			"target":       targetLevel,
			"progress":     progressPercent,
			"levelDisplay": fmt.Sprintf("M%d → M%d", currentLevel, targetLevel),
		}
		dataBytes, _ := json.Marshal([]any{dataRow})

		g.data = append(g.data, DataSource{
			ID:   dataID,
			Type: "inline",
			Data: dataBytes,
		})

		// Create metric widget
		config, _ := json.Marshal(MetricConfig{
			ValueField: "current",
			Format:     "number",
			FormatOptions: &FormatOptions{
				Prefix: "M",
			},
			Subtitle: fmt.Sprintf("Target: M%d", targetLevel),
			Thresholds: []ThresholdValue{
				{Value: 0, Color: "#ef4444"}, // Red
				{Value: 2, Color: "#f59e0b"}, // Yellow
				{Value: 3, Color: "#22c55e"}, // Green
				{Value: 4, Color: "#3b82f6"}, // Blue
			},
		})

		g.widgets = append(g.widgets, Widget{
			ID:           fmt.Sprintf("domain-%s-metric", domainKey),
			Type:         "metric",
			Title:        domain.Name,
			Position:     Position{X: i * width, Y: g.row, W: width, H: 2},
			DataSourceID: dataID,
			Config:       config,
		})
	}
	g.row += 2
}

func (g *Generator) addLevelProgressCharts() {
	domains := g.getSortedDomains()

	for _, domainKey := range domains {
		domain := g.spec.Domains[domainKey]
		assessment := g.spec.Assessments[domainKey]

		// Build level progress data
		levelData := g.buildLevelProgressData(domain, assessment)
		dataID := fmt.Sprintf("level-progress-%s", domainKey)
		dataBytes, _ := json.Marshal(levelData)

		g.data = append(g.data, DataSource{
			ID:   dataID,
			Type: "inline",
			Data: dataBytes,
		})

		// Create stacked horizontal bar chart for level progress
		config, _ := json.Marshal(ChartConfig{
			Marks: []Mark{
				{
					ID:       "completed",
					Geometry: "bar",
					Encode:   Encode{X: "completed", Y: "sli"},
					Stack:    "progress",
					Name:     "Completed",
					Style:    &Style{Color: "#22c55e"},
				},
				{
					ID:       "inProgress",
					Geometry: "bar",
					Encode:   Encode{X: "inProgress", Y: "sli"},
					Stack:    "progress",
					Name:     "In Progress",
					Style:    &Style{Color: "#f59e0b"},
				},
				{
					ID:       "remaining",
					Geometry: "bar",
					Encode:   Encode{X: "remaining", Y: "sli"},
					Stack:    "progress",
					Name:     "Remaining",
					Style:    &Style{Color: "#e5e7eb", Opacity: float64Ptr(0.5)},
				},
			},
			Axes: []Axis{
				{ID: "x", Type: "value", Position: "bottom", Min: float64Ptr(0), Max: float64Ptr(100)},
				{ID: "y", Type: "category", Position: "left"},
			},
			Grid: &Grid{
				Left:         "20%",
				Right:        "5%",
				Top:          "10%",
				Bottom:       "15%",
				ContainLabel: true,
			},
			Legend:  &Legend{Show: true, Position: "top"},
			Tooltip: &Tooltip{Show: true, Trigger: "axis"},
		})

		g.widgets = append(g.widgets, Widget{
			ID:           fmt.Sprintf("level-progress-%s", domainKey),
			Type:         "chart",
			Title:        fmt.Sprintf("%s - SLI Progress to Target", domain.Name),
			Position:     Position{X: 0, Y: g.row, W: 12, H: 4},
			DataSourceID: dataID,
			Config:       config,
		})
		g.row += 4
	}
}

func (g *Generator) buildLevelProgressData(domain *maturity.DomainModel, assessment *maturity.DomainAssessment) []map[string]any {
	var data []map[string]any

	if domain == nil {
		return data
	}

	// Collect all SLIs used in criteria
	sliProgress := make(map[string]map[string]any)

	for _, level := range domain.Levels {
		for _, criterion := range level.Criteria {
			sliID := criterion.SLIID
			if sliID == "" {
				sliID = criterion.ID
			}

			sliName := criterion.Name
			if sli, ok := g.spec.SLIs[criterion.SLIID]; ok && sli != nil {
				sliName = sli.Name
			}

			if _, exists := sliProgress[sliID]; !exists {
				sliProgress[sliID] = map[string]any{
					"sli":        sliName,
					"completed":  0.0,
					"inProgress": 0.0,
					"remaining":  100.0,
				}
			}

			// Check if criterion is met
			if assessment != nil {
				if val, ok := assessment.CriteriaValues[criterion.ID]; ok {
					isMet := criterion.CheckMet(val)
					if isMet {
						sliProgress[sliID]["completed"] = 100.0
						sliProgress[sliID]["inProgress"] = 0.0
						sliProgress[sliID]["remaining"] = 0.0
					} else {
						// Partial progress
						progress := calculateProgress(criterion, val)
						sliProgress[sliID]["completed"] = 0.0
						sliProgress[sliID]["inProgress"] = progress
						sliProgress[sliID]["remaining"] = 100.0 - progress
					}
				}
			}
		}
	}

	// Convert to slice
	for _, v := range sliProgress {
		data = append(data, v)
	}

	// Sort by SLI name
	sort.Slice(data, func(i, j int) bool {
		return data[i]["sli"].(string) < data[j]["sli"].(string)
	})

	return data
}

func calculateProgress(criterion maturity.Criterion, current float64) float64 {
	target := criterion.Target
	if target == 0 {
		return 0
	}

	switch criterion.Operator {
	case maturity.OpGTE, maturity.OpGT:
		if current >= target {
			return 100
		}
		return (current / target) * 100
	case maturity.OpLTE, maturity.OpLT:
		if current <= target {
			return 100
		}
		// Inverse progress for "lower is better"
		if current == 0 {
			return 100
		}
		return (target / current) * 100
	case maturity.OpEQ:
		if current == target {
			return 100
		}
		return 0
	}
	return 0
}

func (g *Generator) addSLITables() {
	if len(g.spec.SLIs) == 0 {
		return
	}

	// Group SLIs by category
	byCategory := make(map[string][]*maturity.SLI)
	for _, sli := range g.spec.SLIs {
		cat := sli.Category
		if cat == "" {
			cat = "other"
		}
		byCategory[cat] = append(byCategory[cat], sli)
	}

	// Sort categories
	var categories []string
	for cat := range byCategory {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, category := range categories {
		slis := byCategory[category]

		// Sort SLIs by name
		sort.Slice(slis, func(i, j int) bool {
			return slis[i].Name < slis[j].Name
		})

		// Build table data
		var tableData []map[string]any
		for _, sli := range slis {
			row := map[string]any{
				"name":    sli.Name,
				"type":    sli.SLIType,
				"layer":   sli.Layer,
				"unit":    sli.Unit,
				"metric":  sli.MetricName,
				"sliType": formatSLIType(sli.SLIType),
			}
			tableData = append(tableData, row)
		}

		dataID := fmt.Sprintf("sli-table-%s", category)
		dataBytes, _ := json.Marshal(tableData)

		g.data = append(g.data, DataSource{
			ID:   dataID,
			Type: "inline",
			Data: dataBytes,
		})

		config, _ := json.Marshal(TableConfig{
			Columns: []TableColumn{
				{Field: "name", Header: "SLI Name", Width: "25%"},
				{Field: "sliType", Header: "Type", Width: "15%"},
				{Field: "layer", Header: "Layer", Width: "15%"},
				{Field: "metric", Header: "Metric", Width: "30%"},
				{Field: "unit", Header: "Unit", Width: "15%"},
			},
			Sortable: true,
			Striped:  true,
		})

		g.widgets = append(g.widgets, Widget{
			ID:           fmt.Sprintf("sli-table-%s", category),
			Type:         "table",
			Title:        fmt.Sprintf("SLIs - %s", formatCategory(category)),
			Position:     Position{X: 0, Y: g.row, W: 12, H: 3},
			DataSourceID: dataID,
			Config:       config,
		})
		g.row += 3
	}
}

func (g *Generator) getSortedDomains() []string {
	var domains []string
	for k := range g.spec.Domains {
		domains = append(domains, k)
	}
	sort.Strings(domains)
	return domains
}

func formatSLIType(sliType string) string {
	switch sliType {
	case "availability":
		return "Availability"
	case "latency":
		return "Latency"
	case "error_rate":
		return "Error Rate"
	case "throughput":
		return "Throughput"
	case "saturation":
		return "Saturation"
	case "utilization":
		return "Utilization"
	case "quality":
		return "Quality"
	case "freshness":
		return "Freshness"
	default:
		return sliType
	}
}

func formatCategory(category string) string {
	switch category {
	case "prevention":
		return "Prevention"
	case "detection":
		return "Detection"
	case "response":
		return "Response"
	case "reliability":
		return "Reliability"
	case "efficiency":
		return "Efficiency"
	default:
		return category
	}
}

func float64Ptr(v float64) *float64 {
	return &v
}

// ToJSON returns the dashboard as JSON bytes.
func (d *Dashboard) ToJSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

// GenerateMaturityBullets creates bullet chart data for all domains.
func (g *Generator) GenerateMaturityBullets() *MaturityBulletData {
	if g.spec == nil {
		return &MaturityBulletData{}
	}

	bullets := []MaturityBullet{}

	for domainKey, domain := range g.spec.Domains {
		assessment := g.spec.Assessments[domainKey]

		currentLevel := float64(1)
		targetLevel := float64(5)

		if assessment != nil {
			currentLevel = float64(assessment.CurrentLevel)
			targetLevel = float64(assessment.TargetLevel)
		}

		bullets = append(bullets, NewMaturityBullet(
			domain.Name,
			fmt.Sprintf("M%d → M%d", int(currentLevel), int(targetLevel)),
			currentLevel,
			targetLevel,
		))

		// Add bullets for each SLI in the domain
		for _, level := range domain.Levels {
			for _, criterion := range level.Criteria {
				sliName := criterion.Name
				if sli, ok := g.spec.SLIs[criterion.SLIID]; ok && sli != nil {
					sliName = sli.Name
				}

				// Calculate current maturity level for this SLI
				sliLevel := float64(1)
				if assessment != nil {
					sliLevel = g.calculateSLIMaturityLevel(domain, assessment, criterion.SLIID)
				}

				bullets = append(bullets, NewMaturityBullet(
					sliName,
					MaturityLevel(sliLevel),
					sliLevel,
					targetLevel,
				))
			}
		}
	}

	return &MaturityBulletData{Bullets: bullets}
}

// calculateSLIMaturityLevel determines the highest level achieved for an SLI.
func (g *Generator) calculateSLIMaturityLevel(domain *maturity.DomainModel, assessment *maturity.DomainAssessment, sliID string) float64 {
	highestLevel := float64(1)

	for _, level := range domain.Levels {
		for _, criterion := range level.Criteria {
			if criterion.SLIID != sliID {
				continue
			}

			if val, ok := assessment.CriteriaValues[criterion.ID]; ok {
				if criterion.CheckMet(val) {
					levelNum := float64(level.Level)
					if levelNum > highestLevel {
						highestLevel = levelNum
					}
				} else {
					// Partial progress within level
					progress := calculateProgress(criterion, val)
					partialLevel := float64(level.Level-1) + (progress / 100)
					if partialLevel > highestLevel {
						highestLevel = partialLevel
					}
				}
			}
		}
	}

	return highestLevel
}

// addBulletWidgets adds maturity bullet chart widgets.
func (g *Generator) addBulletWidgets() {
	domains := g.getSortedDomains()

	for _, domainKey := range domains {
		domain := g.spec.Domains[domainKey]
		assessment := g.spec.Assessments[domainKey]

		// Collect SLI bullets for this domain
		var bullets []MaturityBullet

		currentLevel := float64(1)
		targetLevel := float64(5)
		if assessment != nil {
			currentLevel = float64(assessment.CurrentLevel)
			targetLevel = float64(assessment.TargetLevel)
		}

		// Domain-level bullet
		bullets = append(bullets, NewMaturityBullet(
			domain.Name,
			fmt.Sprintf("Overall: M%d → M%d", int(currentLevel), int(targetLevel)),
			currentLevel,
			targetLevel,
		))

		// SLI-level bullets
		seenSLIs := make(map[string]bool)
		for _, level := range domain.Levels {
			for _, criterion := range level.Criteria {
				sliID := criterion.SLIID
				if sliID == "" || seenSLIs[sliID] {
					continue
				}
				seenSLIs[sliID] = true

				sliName := criterion.Name
				if sli, ok := g.spec.SLIs[sliID]; ok && sli != nil {
					sliName = sli.Name
				}

				sliLevel := float64(1)
				if assessment != nil {
					sliLevel = g.calculateSLIMaturityLevel(domain, assessment, sliID)
				}

				bullets = append(bullets, NewMaturityBullet(
					sliName,
					MaturityLevel(sliLevel),
					sliLevel,
					targetLevel,
				))
			}
		}

		// Create data source
		dataID := fmt.Sprintf("bullet-%s", domainKey)
		dataBytes, _ := json.Marshal(bullets)

		g.data = append(g.data, DataSource{
			ID:   dataID,
			Type: "inline",
			Data: dataBytes,
		})

		// Create bullet widget
		config, _ := json.Marshal(map[string]any{
			"bulletType": "maturity",
		})

		g.widgets = append(g.widgets, Widget{
			ID:           fmt.Sprintf("bullet-%s", domainKey),
			Type:         "bullet",
			Title:        fmt.Sprintf("%s - Maturity Levels", domain.Name),
			Position:     Position{X: 0, Y: g.row, W: 12, H: 4},
			DataSourceID: dataID,
			Config:       config,
		})
		g.row += 4
	}
}
