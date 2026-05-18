// Package dashboard provides HTML dashboard generation for PRISM maturity models.
package dashboard

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/grokify/prism-intelligence"
	"github.com/grokify/prism-intelligence/maturity"
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
	spec     *maturity.Spec
	stateDoc *prism.PRISMDocument // Optional: PRISM Maturity State document
	widgets  []Widget
	data     []DataSource
	row      int
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

// WithStateDocument adds a PRISM Maturity State document for state tracking.
// When set, state is read from this document instead of the legacy assessments.
func (g *Generator) WithStateDocument(doc *prism.PRISMDocument) *Generator {
	g.stateDoc = doc
	return g
}

// domainState represents maturity state for a domain (from either source).
type domainState struct {
	CurrentLevel int
	TargetLevel  int
	AssessedAt   string
	AssessedBy   string
}

// getDomainState returns maturity state for a domain.
// Reads from PRISM Maturity State document.
func (g *Generator) getDomainState(domainKey string) *domainState {
	// Try state document
	if g.stateDoc != nil && g.stateDoc.MaturityState != nil {
		if state, ok := g.stateDoc.MaturityState[domainKey]; ok && state != nil {
			ds := &domainState{
				CurrentLevel: 1,
				TargetLevel:  5,
			}
			if state.Current != nil {
				ds.CurrentLevel = state.Current.Level
				ds.AssessedAt = state.Current.AchievedAt
				ds.AssessedBy = state.Current.AssessedBy
			}
			if state.Target != nil {
				ds.TargetLevel = state.Target.Level
			}
			return ds
		}
	}

	// Default values (no state document provided)
	return &domainState{
		CurrentLevel: 1,
		TargetLevel:  5,
	}
}

// getSLIValue returns the current value for an SLI.
// Reads from PRISM Maturity State document.
// Returns the value and whether it was found.
func (g *Generator) getSLIValue(sliID string, window string) (float64, bool) {
	// Try state document
	if g.stateDoc != nil && g.stateDoc.SLIState != nil {
		if state, ok := g.stateDoc.SLIState[sliID]; ok && state != nil {
			// Try specific window first
			if window != "" && state.Windows != nil {
				if ws, ok := state.Windows[window]; ok && ws != nil {
					return ws.Value, true
				}
			}
			// Try default window (30d)
			if state.Windows != nil {
				if ws, ok := state.Windows["30d"]; ok && ws != nil {
					return ws.Value, true
				}
				// Return any available window
				for _, ws := range state.Windows {
					if ws != nil {
						return ws.Value, true
					}
				}
			}
		}
	}

	return 0, false
}

// getSLIQualitativeState returns the qualitative state for an SLI.
// Reads from PRISM Maturity State document.
func (g *Generator) getSLIQualitativeState(sliID string) string {
	// Try state document
	if g.stateDoc != nil && g.stateDoc.SLIState != nil {
		if state, ok := g.stateDoc.SLIState[sliID]; ok && state != nil {
			return state.QualitativeState
		}
	}

	return ""
}

// getCriterionValue returns the current value for a criterion.
// Reads from PRISM Maturity State document by SLI ID.
func (g *Generator) getCriterionValue(_ string, criterion maturity.Criterion) (float64, bool) {
	// Try by SLI ID from state document
	if criterion.SLIID != "" {
		if val, ok := g.getSLIValue(criterion.SLIID, ""); ok {
			return val, true
		}
	}

	return 0, false
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
		state := g.getDomainState(domainKey)

		// Create data source for this domain
		dataID := fmt.Sprintf("domain-%s-data", domainKey)
		currentLevel := state.CurrentLevel
		targetLevel := state.TargetLevel
		progressPercent := 0.0

		if targetLevel > 1 {
			progressPercent = float64(currentLevel-1) / float64(targetLevel-1) * 100
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

		// Build level progress data
		levelData := g.buildLevelProgressData(domainKey, domain)
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

		// Calculate height based on number of SLIs (1 row per SLI + 2 for legend/padding)
		height := len(levelData) + 2
		if height < 4 {
			height = 4
		}
		if height > 12 {
			height = 12
		}

		g.widgets = append(g.widgets, Widget{
			ID:           fmt.Sprintf("level-progress-%s", domainKey),
			Type:         "chart",
			Title:        fmt.Sprintf("%s - SLI Progress to Target", domain.Name),
			Position:     Position{X: 0, Y: g.row, W: 12, H: height},
			DataSourceID: dataID,
			Config:       config,
		})
		g.row += height
	}
}

func (g *Generator) buildLevelProgressData(domainKey string, domain *maturity.DomainModel) []map[string]any {
	var data []map[string]any

	if domain == nil {
		return data
	}

	// First pass: collect unique SLIs
	seenSLIs := make(map[string]bool)
	var sliIDs []string

	for _, level := range domain.Levels {
		for _, criterion := range level.Criteria {
			sliID := criterion.SLIID
			if sliID == "" {
				continue
			}
			if !seenSLIs[sliID] {
				seenSLIs[sliID] = true
				sliIDs = append(sliIDs, sliID)
			}
		}
	}

	// Second pass: for each SLI, find its value and calculate progress
	for _, sliID := range sliIDs {
		sliName := sliID
		category := "other"
		if sli, ok := g.spec.SLIs[sliID]; ok && sli != nil {
			sliName = sli.Name
			if sli.Category != "" {
				category = sli.Category
			}
		}

		// Search through all criteria to find a value (same as collectSLIsByType)
		var actualValue float64
		var foundValue bool
		var bestCriterion maturity.Criterion

		for _, lvl := range domain.Levels {
			for _, c := range lvl.Criteria {
				if c.SLIID == sliID {
					if val, ok := g.getCriterionValue(domainKey, c); ok {
						actualValue = val
						foundValue = true
						bestCriterion = c
						break
					}
				}
			}
			if foundValue {
				break
			}
		}

		row := map[string]any{
			"sli":        sliName,
			"category":   category,
			"completed":  0.0,
			"inProgress": 0.0,
			"remaining":  100.0,
		}

		if foundValue {
			isMet := bestCriterion.CheckMet(actualValue)
			if isMet {
				row["completed"] = 100.0
				row["inProgress"] = 0.0
				row["remaining"] = 0.0
			} else {
				progress := calculateProgress(bestCriterion, actualValue)
				row["completed"] = 0.0
				row["inProgress"] = progress
				row["remaining"] = 100.0 - progress
			}
		}

		data = append(data, row)
	}

	// Sort by category (NIST CSF order) then alphabetically by SLI name
	catWeights := maturity.CategorySortWeight()
	sort.Slice(data, func(i, j int) bool {
		catI := data[i]["category"].(string)
		catJ := data[j]["category"].(string)
		weightI := catWeights[catI]
		weightJ := catWeights[catJ]
		if weightI == 0 {
			weightI = 100 // Unknown categories sort last
		}
		if weightJ == 0 {
			weightJ = 100
		}
		if weightI != weightJ {
			return weightI < weightJ
		}
		// Within same category, sort alphabetically by SLI name
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

	// Sort categories by NIST CSF order
	var categories []string
	for cat := range byCategory {
		categories = append(categories, cat)
	}
	catWeights := maturity.CategorySortWeight()
	sort.Slice(categories, func(i, j int) bool {
		weightI := catWeights[categories[i]]
		weightJ := catWeights[categories[j]]
		if weightI == 0 {
			weightI = 100 // Unknown categories sort last
		}
		if weightJ == 0 {
			weightJ = 100
		}
		if weightI != weightJ {
			return weightI < weightJ
		}
		return categories[i] < categories[j] // Alphabetical fallback
	})

	// Build SLI order map from spec categories
	sliOrderMap := make(map[string]int)
	for _, cat := range g.spec.Categories {
		for idx, sliID := range cat.SLIOrder {
			sliOrderMap[sliID] = idx
		}
	}

	for _, category := range categories {
		slis := byCategory[category]

		// Sort SLIs by custom order if defined, otherwise by name
		sort.Slice(slis, func(i, j int) bool {
			orderI, hasI := sliOrderMap[slis[i].ID]
			orderJ, hasJ := sliOrderMap[slis[j].ID]
			if hasI && hasJ {
				return orderI < orderJ
			}
			if hasI {
				return true
			}
			if hasJ {
				return false
			}
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
	// NIST CSF 2.0 categories
	case "govern":
		return "Govern"
	case "identify":
		return "Identify"
	case "protect":
		return "Protect"
	case "detect":
		return "Detect"
	case "respond":
		return "Respond"
	case "recover":
		return "Recover"
	// Legacy/alternative category names
	case "governance":
		return "Governance"
	case "prevention":
		return "Prevention"
	case "detection":
		return "Detection"
	case "response":
		return "Response"
	case "recovery":
		return "Recovery"
	case "reliability":
		return "Reliability"
	case "efficiency":
		return "Efficiency"
	case "quality":
		return "Quality"
	case "availability":
		return "Availability"
	case "other":
		return "Other"
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
		state := g.getDomainState(domainKey)

		currentLevel := float64(state.CurrentLevel)
		targetLevel := float64(state.TargetLevel)

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
				sliLevel := g.calculateSLIMaturityLevel(domainKey, domain, criterion.SLIID)

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
// Uses dual-read to get values from either PRISM Maturity State or legacy assessments.
func (g *Generator) calculateSLIMaturityLevel(domainKey string, domain *maturity.DomainModel, sliID string) float64 {
	highestLevel := float64(1)

	for _, level := range domain.Levels {
		for _, criterion := range level.Criteria {
			if criterion.SLIID != sliID {
				continue
			}

			// Use dual-read helper to get criterion value
			if val, ok := g.getCriterionValue(domainKey, criterion); ok {
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

// addBulletWidgets adds maturity bullet chart widgets grouped by category.
// Categories are sorted by NIST CSF order (govern → identify → protect → detect → respond → recover).
// Works with any maturity model - gracefully handles missing SLIs, SLITypes, or assessments.
func (g *Generator) addBulletWidgets() {
	domains := g.getSortedDomains()
	if len(domains) == 0 {
		return
	}

	for _, domainKey := range domains {
		domain := g.spec.Domains[domainKey]
		state := g.getDomainState(domainKey)

		currentLevel := float64(state.CurrentLevel)
		targetLevel := float64(state.TargetLevel)

		// Add domain overview bullet
		g.addDomainOverviewBullet(domainKey, domain, currentLevel, targetLevel)

		// Collect SLIs for this domain with their types
		slisByType := g.collectSLIsByType(domainKey, domain, targetLevel)

		// Add bullets grouped by category (NIST CSF order)
		g.addFlatBulletList(domainKey, domain, slisByType)
	}
}

// sliInfo holds SLI data for grouping.
type sliInfo struct {
	ID               string
	Name             string
	SLIType          string
	Level            float64
	Target           float64
	ActualValue      float64          // The actual SLI value (e.g., 65 for 65%)
	Unit             string           // Unit of measurement (e.g., "%", "ms", "rps")
	QualitativeState string           // For qualitative/hybrid SLIs (e.g., "tracked", "measured")
	Thresholds       []LevelThreshold // Thresholds for each maturity level
}

func (g *Generator) addDomainOverviewBullet(domainKey string, domain *maturity.DomainModel, currentLevel, targetLevel float64) {
	bullet := NewMaturityBullet(
		domain.Name+" Overall",
		fmt.Sprintf("M%d → M%d", int(currentLevel), int(targetLevel)),
		currentLevel,
		targetLevel,
	)

	dataID := fmt.Sprintf("bullet-%s-overview", domainKey)
	dataBytes, _ := json.Marshal([]MaturityBullet{bullet})

	g.data = append(g.data, DataSource{
		ID:   dataID,
		Type: "inline",
		Data: dataBytes,
	})

	config, _ := json.Marshal(map[string]any{
		"bulletType": "overview",
	})

	g.widgets = append(g.widgets, Widget{
		ID:           dataID,
		Type:         "bullet",
		Title:        fmt.Sprintf("%s - Overall Maturity", domain.Name),
		Position:     Position{X: 0, Y: g.row, W: 12, H: 2},
		DataSourceID: dataID,
		Config:       config,
	})
	g.row += 2
}

func (g *Generator) collectSLIsByType(domainKey string, domain *maturity.DomainModel, targetLevel float64) map[string][]sliInfo {
	slisByType := make(map[string][]sliInfo)
	seenSLIs := make(map[string]bool)

	for _, level := range domain.Levels {
		for _, criterion := range level.Criteria {
			sliID := criterion.SLIID
			if sliID == "" {
				continue
			}

			// Get or create SLI info
			if !seenSLIs[sliID] {
				seenSLIs[sliID] = true

				sliName := criterion.Name
				sliType := ""
				unit := ""
				if sli, ok := g.spec.SLIs[sliID]; ok && sli != nil {
					sliName = sli.Name
					sliType = sli.SLIType
					unit = sli.Unit
				}

				// Use dual-read helpers for level and value
				sliLevel := g.calculateSLIMaturityLevel(domainKey, domain, sliID)
				actualValue := float64(0)

				// Get actual value from first matching criterion using dual-read
				for _, lvl := range domain.Levels {
					for _, c := range lvl.Criteria {
						if c.SLIID == sliID {
							if val, ok := g.getCriterionValue(domainKey, c); ok {
								actualValue = val
								break
							}
						}
					}
					if actualValue != 0 {
						break
					}
				}

				info := sliInfo{
					ID:               sliID,
					Name:             sliName,
					SLIType:          sliType,
					Level:            sliLevel,
					Target:           targetLevel,
					ActualValue:      actualValue,
					Unit:             unit,
					QualitativeState: g.getSLIQualitativeState(sliID),
					Thresholds:       []LevelThreshold{},
				}

				slisByType[sliType] = append(slisByType[sliType], info)
			}

			// Collect threshold for this level
			for i := range slisByType[g.getSLIType(sliID)] {
				if slisByType[g.getSLIType(sliID)][i].ID == sliID {
					threshold := LevelThreshold{
						Level:    level.Level,
						Operator: criterion.Operator,
						Value:    criterion.Target,
						ValueStr: formatThresholdValue(criterion.Target, criterion.Operator, g.getSLIUnit(sliID)),
					}
					slisByType[g.getSLIType(sliID)][i].Thresholds = append(
						slisByType[g.getSLIType(sliID)][i].Thresholds,
						threshold,
					)
					break
				}
			}
		}
	}

	return slisByType
}

func (g *Generator) getSLIType(sliID string) string {
	if sli, ok := g.spec.SLIs[sliID]; ok && sli != nil {
		return sli.SLIType
	}
	return ""
}

func (g *Generator) getSLIUnit(sliID string) string {
	if sli, ok := g.spec.SLIs[sliID]; ok && sli != nil {
		return sli.Unit
	}
	return ""
}

func formatThresholdValue(value float64, operator, unit string) string {
	// Handle qualitative criteria
	if operator == "exists" {
		return "Tracked"
	}

	// Format the value with appropriate precision
	var valStr string
	if value == float64(int(value)) {
		valStr = fmt.Sprintf("%d", int(value))
	} else {
		valStr = fmt.Sprintf("%.1f", value)
	}

	// Add unit
	if unit != "" {
		valStr += unit
	}

	// Add operator prefix
	switch operator {
	case ">=", "gte":
		return "≥" + valStr
	case "<=", "lte":
		return "≤" + valStr
	case ">", "gt":
		return ">" + valStr
	case "<", "lt":
		return "<" + valStr
	case "==", "eq":
		return "=" + valStr
	default:
		return valStr
	}
}

// addFlatBulletList adds bullets grouped by category with section headers.
// Categories are sorted by NIST CSF order (govern → identify → protect → detect → respond → recover).
func (g *Generator) addFlatBulletList(domainKey string, domain *maturity.DomainModel, slisByType map[string][]sliInfo) {
	// Collect all SLIs and group by category
	slisByCategory := make(map[string][]sliInfo)
	for _, infos := range slisByType {
		for _, info := range infos {
			// Get category from SLI definition
			category := "other"
			if sli, ok := g.spec.SLIs[info.ID]; ok && sli != nil && sli.Category != "" {
				category = sli.Category
			}
			slisByCategory[category] = append(slisByCategory[category], info)
		}
	}

	if len(slisByCategory) == 0 {
		return
	}

	// Sort categories by NIST CSF order
	var categories []string
	for cat := range slisByCategory {
		categories = append(categories, cat)
	}
	catWeights := maturity.CategorySortWeight()
	sort.Slice(categories, func(i, j int) bool {
		weightI := catWeights[categories[i]]
		weightJ := catWeights[categories[j]]
		if weightI == 0 {
			weightI = 100 // Unknown categories sort last
		}
		if weightJ == 0 {
			weightJ = 100
		}
		if weightI != weightJ {
			return weightI < weightJ
		}
		return categories[i] < categories[j] // Alphabetical fallback
	})

	// Build SLI order map from spec categories
	sliOrderMap := make(map[string]int)
	for _, cat := range g.spec.Categories {
		for idx, sliID := range cat.SLIOrder {
			sliOrderMap[sliID] = idx
		}
	}

	// Add bullet section for each category
	for _, category := range categories {
		infos := slisByCategory[category]

		// Sort SLIs by custom order if defined, otherwise by name
		sort.Slice(infos, func(i, j int) bool {
			orderI, hasI := sliOrderMap[infos[i].ID]
			orderJ, hasJ := sliOrderMap[infos[j].ID]
			if hasI && hasJ {
				return orderI < orderJ
			}
			if hasI {
				return true
			}
			if hasJ {
				return false
			}
			return infos[i].Name < infos[j].Name
		})

		// Build bullets for this category
		var bullets []MaturityBullet
		for _, info := range infos {
			bullet := NewMaturityBulletWithDetails(
				info.Name,
				info.Level,
				info.Target,
				info.ActualValue,
				info.Unit,
				info.QualitativeState,
				info.Thresholds,
			)
			bullets = append(bullets, bullet)
		}

		if len(bullets) == 0 {
			continue
		}

		dataID := fmt.Sprintf("bullet-%s-%s", domainKey, category)
		dataBytes, _ := json.Marshal(bullets)

		g.data = append(g.data, DataSource{
			ID:   dataID,
			Type: "inline",
			Data: dataBytes,
		})

		config, _ := json.Marshal(map[string]any{
			"bulletType": "maturity",
			"category":   category,
		})

		// Calculate height based on number of bullets
		height := len(bullets) + 1
		if height < 2 {
			height = 2
		}
		if height > 8 {
			height = 8
		}

		g.widgets = append(g.widgets, Widget{
			ID:           fmt.Sprintf("bullet-%s-%s", domainKey, category),
			Type:         "bullet",
			Title:        fmt.Sprintf("%s - %s", domain.Name, formatCategory(category)),
			Position:     Position{X: 0, Y: g.row, W: 12, H: height},
			DataSourceID: dataID,
			Config:       config,
		})
		g.row += height
	}
}
