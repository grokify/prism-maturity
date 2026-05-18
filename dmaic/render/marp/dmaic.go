// Package marp provides a Marp markdown renderer for DMAIC documents.
// Marp is a presentation ecosystem that converts Markdown to slides.
// See https://marp.app/ for more information.
package marp

import (
	"bytes"
	"fmt"
	"text/template"

	sdmarp "github.com/grokify/structureddocs/marp"

	"github.com/grokify/prism-intelligence/dmaic"
	"github.com/grokify/prism-intelligence/dmaic/render"
)

// Renderer implements the render.Renderer interface for DMAIC Marp output.
type Renderer struct{}

// New creates a new DMAIC Marp renderer.
func New() *Renderer {
	return &Renderer{}
}

// Format returns the output format name.
func (r *Renderer) Format() string {
	return "marp"
}

// FileExtension returns the file extension for Marp output.
func (r *Renderer) FileExtension() string {
	return ".md"
}

// Render converts a DMAIC document to Marp markdown slides.
func (r *Renderer) Render(doc *dmaic.DMAICDocument, opts *render.Options) ([]byte, error) {
	if opts == nil {
		opts = render.DefaultOptions()
	}

	data := &templateData{
		Doc:            doc,
		Options:        opts,
		Theme:          sdmarp.GetTheme(opts.Theme),
		OverallHealth:  doc.CalculateOverallHealth(),
		HasInitiatives: opts.IncludeInitiatives && len(doc.Initiatives) > 0,
	}

	var buf bytes.Buffer

	// Render front matter
	if err := frontMatterTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering front matter: %w", err)
	}

	// Render title slide
	if err := titleSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering title slide: %w", err)
	}

	// Render executive summary slide
	if err := executiveSummaryTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering executive summary: %w", err)
	}

	// Render category overview slides
	for i, cat := range doc.Categories {
		catData := &categoryData{
			templateData: data,
			Category:     cat,
			Index:        i + 1,
		}
		if err := categoryOverviewTmpl.Execute(&buf, catData); err != nil {
			return nil, fmt.Errorf("rendering category slide %d: %w", i+1, err)
		}
	}

	// Render metric detail slides for Yellow/Red status metrics
	for _, cat := range doc.Categories {
		for _, m := range cat.Metrics {
			status := m.CalculateStatus()
			if status == dmaic.StatusYellow || status == dmaic.StatusRed {
				metricData := &metricDetailData{
					templateData: data,
					Metric:       m,
					CategoryName: cat.Name,
				}
				if err := metricDetailTmpl.Execute(&buf, metricData); err != nil {
					return nil, fmt.Errorf("rendering metric detail: %w", err)
				}
			}
		}
	}

	// Render phase summary if grouping by phase
	if opts.GroupByPhase {
		if err := phaseSummaryTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering phase summary: %w", err)
		}
	}

	// Render root cause analysis if enabled
	if opts.IncludeRootCauses {
		hasRootCauses := false
		for _, cat := range doc.Categories {
			for _, m := range cat.Metrics {
				if len(m.RootCauses) > 0 {
					hasRootCauses = true
					break
				}
			}
		}
		if hasRootCauses {
			if err := rootCausesTmpl.Execute(&buf, data); err != nil {
				return nil, fmt.Errorf("rendering root causes: %w", err)
			}
		}
	}

	// Render initiatives summary
	if data.HasInitiatives {
		if err := initiativesTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering initiatives: %w", err)
		}
	}

	// Render summary/questions slide
	if err := summarySlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering summary: %w", err)
	}

	return buf.Bytes(), nil
}

// templateData holds data for DMAIC template rendering.
type templateData struct {
	Doc            *dmaic.DMAICDocument
	Options        *render.Options
	Theme          sdmarp.ThemeConfig
	OverallHealth  float64
	HasInitiatives bool
}

// categoryData holds data for category slide rendering.
type categoryData struct {
	*templateData
	Category dmaic.Category
	Index    int
}

// metricDetailData holds data for metric detail slide rendering.
type metricDetailData struct {
	*templateData
	Metric       dmaic.Metric
	CategoryName string
}

// funcMap merges structureddocs CommonFuncMap with DMAIC-specific functions.
var funcMap = mergeFuncMaps(sdmarp.CommonFuncMap, template.FuncMap{
	"statusIcon": func(status string) string {
		switch status {
		case dmaic.StatusGreen:
			return "[GREEN]"
		case dmaic.StatusYellow:
			return "[YELLOW]"
		case dmaic.StatusRed:
			return "[RED]"
		default:
			return ""
		}
	},
	"statusColor": func(status string) string {
		switch status {
		case dmaic.StatusGreen:
			return "#38a169"
		case dmaic.StatusYellow:
			return "#d69e2e"
		case dmaic.StatusRed:
			return "#e53e3e"
		default:
			return "#718096"
		}
	},
	"phaseIcon": func(phase string) string {
		switch phase {
		case dmaic.PhaseDefine:
			return "[D]"
		case dmaic.PhaseMeasure:
			return "[M]"
		case dmaic.PhaseAnalyze:
			return "[A]"
		case dmaic.PhaseImprove:
			return "[I]"
		case dmaic.PhaseControl:
			return "[C]"
		default:
			return ""
		}
	},
	"healthPercent": func(h float64) string {
		return fmt.Sprintf("%.0f%%", h*100)
	},
	"formatValue": func(v float64, unit string) string {
		if unit == "%" {
			return fmt.Sprintf("%.1f%%", v)
		}
		if v == float64(int(v)) {
			return fmt.Sprintf("%.0f", v)
		}
		return fmt.Sprintf("%.2f", v)
	},
	"calculateStatus": func(m dmaic.Metric) string {
		return m.CalculateStatus()
	},
	"categoryHealth": func(c dmaic.Category) float64 {
		return c.CalculateCategoryHealth()
	},
})

// mergeFuncMaps merges multiple template.FuncMaps into one.
func mergeFuncMaps(fmaps ...template.FuncMap) template.FuncMap {
	result := make(template.FuncMap)
	for _, fm := range fmaps {
		for k, v := range fm {
			result[k] = v
		}
	}
	return result
}

// Templates
var frontMatterTmpl = template.Must(template.New("frontMatter").Parse(`---
marp: true
theme: {{.Theme.Name}}
paginate: true
{{- if and .Doc.Metadata .Doc.Metadata.Team}}
header: "{{.Doc.Metadata.Team}}{{if .Doc.Metadata.Period}} | {{.Doc.Metadata.Period}}{{end}}"
{{- end}}
{{- if and .Doc.Metadata .Doc.Metadata.Name}}
footer: "DMAIC | {{.Doc.Metadata.Name}}"
{{- end}}
style: |
  section {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
  section.title {
    text-align: center;
    background: linear-gradient(135deg, {{.Theme.PrimaryBgColor}} 0%, {{.Theme.AccentColor}} 100%);
    color: {{.Theme.PrimaryTextColor}};
  }
  section.title h1 {
    font-size: 2.5em;
    color: {{.Theme.PrimaryTextColor}};
  }
  section.category {
    background: linear-gradient(135deg, {{.Theme.PrimaryBgColor}} 0%, {{.Theme.AccentColor}} 100%);
    color: {{.Theme.PrimaryTextColor}};
  }
  table {
    font-size: 0.85em;
    width: 100%;
  }
  th {
    background: #f7fafc;
  }
  .status-green { color: {{.Theme.SuccessColor}}; }
  .status-yellow { color: {{.Theme.WarningColor}}; }
  .status-red { color: {{.Theme.DangerColor}}; }
  blockquote {
    font-size: 1.2em;
    border-left: 4px solid {{.Theme.AccentColor}};
    padding-left: 1em;
    font-style: italic;
  }
  .health-bar {
    background: #e2e8f0;
    border-radius: 4px;
    padding: 2px;
  }
  .health-fill {
    background: {{.Theme.AccentColor}};
    border-radius: 2px;
    height: 20px;
  }
---

`))

var titleSlideTmpl = template.Must(template.New("titleSlide").Funcs(funcMap).Parse(`<!-- _class: title -->

# {{if and .Doc.Metadata .Doc.Metadata.Name}}{{.Doc.Metadata.Name}}{{else}}DMAIC Metrics{{end}}

**DMAIC Continuous Improvement Framework**

{{- if and .Doc.Metadata .Doc.Metadata.Owner}}
**Owner:** {{.Doc.Metadata.Owner}}
{{- end}}
{{- if and .Doc.Metadata .Doc.Metadata.Team}}
**Team:** {{.Doc.Metadata.Team}}
{{- end}}
{{- if and .Doc.Metadata .Doc.Metadata.Period}}
**Period:** {{.Doc.Metadata.Period}}
{{- end}}
{{- if and .Doc.Metadata .Doc.Metadata.Status}}
**Status:** {{.Doc.Metadata.Status}}
{{- end}}

---

`))

var executiveSummaryTmpl = template.Must(template.New("executiveSummary").Funcs(funcMap).Parse(`## Executive Summary

{{- if and .Doc.Metadata .Doc.Metadata.Description}}

{{.Doc.Metadata.Description}}
{{- end}}

**Overall Health:** [{{progressBar .OverallHealth}}] {{healthPercent .OverallHealth}}

### Categories Overview

| Category | Health | Metrics | Status |
|----------|--------|---------|--------|
{{- range $i, $cat := .Doc.Categories}}
| {{$cat.Name}} | {{healthPercent (categoryHealth $cat)}} | {{len $cat.Metrics}} | {{- range $j, $m := $cat.Metrics}}{{if $j}}, {{end}}{{statusIcon (calculateStatus $m)}}{{end}} |
{{- end}}

---

`))

var categoryOverviewTmpl = template.Must(template.New("categoryOverview").Funcs(funcMap).Parse(`<!-- _class: category -->

## {{.Category.Name}}

{{- if .Category.Description}}
{{.Category.Description}}
{{- end}}

{{- if .Category.Owner}}
**Owner:** {{.Category.Owner}}
{{- end}}

**Category Health:** {{healthPercent (categoryHealth .Category)}}

---

## {{.Category.Name}} - Metrics

| Metric | Phase | Current | Target | Status |
|--------|-------|---------|--------|--------|
{{- range .Category.Metrics}}
| {{truncate .Name 30}} | {{phaseIcon .Phase}} | {{formatValue .Current .Unit}} | {{formatValue .Target .Unit}} | {{statusIcon (calculateStatus .)}} |
{{- end}}

---

`))

var metricDetailTmpl = template.Must(template.New("metricDetail").Funcs(funcMap).Parse(`## {{statusIcon (calculateStatus .Metric)}} {{.Metric.Name}}

**Category:** {{.CategoryName}}
**Phase:** {{.Metric.Phase}}
**Status:** {{calculateStatus .Metric}}

| Measure | Value |
|---------|-------|
| Baseline | {{formatValue .Metric.Baseline .Metric.Unit}} |
| Current | {{formatValue .Metric.Current .Metric.Unit}} |
| Target | {{formatValue .Metric.Target .Metric.Unit}} |

{{- if .Metric.ControlLimits}}

### Control Limits

| Limit | Value |
|-------|-------|
| UCL | {{formatValue .Metric.ControlLimits.UCL .Metric.Unit}} |
| Center | {{formatValue .Metric.ControlLimits.CenterLine .Metric.Unit}} |
| LCL | {{formatValue .Metric.ControlLimits.LCL .Metric.Unit}} |
{{- end}}

{{- if .Metric.RootCauses}}

### Root Causes

{{- range .Metric.RootCauses}}
- **{{.Description}}** ({{.Category}}, {{.Impact}} impact){{if .Validated}} [Validated]{{end}}
{{- end}}
{{- end}}

---

`))

var phaseSummaryTmpl = template.Must(template.New("phaseSummary").Funcs(funcMap).Parse(`## DMAIC Phase Summary

| Phase | Description | Metrics |
|-------|-------------|---------|
| [D] Define | Define the problem and project goals | {{len (index .Doc.MetricsByPhase "Define")}} |
| [M] Measure | Measure current performance | {{len (index .Doc.MetricsByPhase "Measure")}} |
| [A] Analyze | Analyze root causes | {{len (index .Doc.MetricsByPhase "Analyze")}} |
| [I] Improve | Implement improvements | {{len (index .Doc.MetricsByPhase "Improve")}} |
| [C] Control | Sustain the gains | {{len (index .Doc.MetricsByPhase "Control")}} |

---

`))

var rootCausesTmpl = template.Must(template.New("rootCauses").Funcs(funcMap).Parse(`## Root Cause Analysis

| Metric | Root Cause | Category | Impact | Validated |
|--------|------------|----------|--------|-----------|
{{- range .Doc.Categories}}
{{- range .Metrics}}
{{- $metricName := .Name}}
{{- range .RootCauses}}
| {{truncate $metricName 20}} | {{truncate .Description 30}} | {{.Category}} | {{.Impact}} | {{if .Validated}}Yes{{else}}No{{end}} |
{{- end}}
{{- end}}
{{- end}}

---

`))

var initiativesTmpl = template.Must(template.New("initiatives").Funcs(funcMap).Parse(`## Improvement Initiatives

| Initiative | Owner | Status | Timeline |
|------------|-------|--------|----------|
{{- range .Doc.Initiatives}}
| {{truncate .Name 30}} | {{if .Owner}}{{.Owner}}{{else}}-{{end}} | {{if .Status}}{{.Status}}{{else}}-{{end}} | {{if .StartDate}}{{.StartDate}}{{if .EndDate}} - {{.EndDate}}{{end}}{{else}}-{{end}} |
{{- end}}

{{- range .Doc.Initiatives}}
{{- if .ExpectedImpact}}

### {{.Name}}

**Expected Impact:** {{.ExpectedImpact}}
{{- if .ActualImpact}}
**Actual Impact:** {{.ActualImpact}}
{{- end}}
{{- end}}
{{- end}}

---

`))

var summarySlideTmpl = template.Must(template.New("summarySlide").Funcs(funcMap).Parse(`<!-- _class: title -->

## Summary

**Overall Health:** {{healthPercent .OverallHealth}}

{{- range .Doc.Categories}}
- **{{.Name}}:** {{healthPercent (categoryHealth .)}} ({{len .Metrics}} metrics)
{{- end}}

---

## Questions?

{{- if and .Doc.Metadata .Doc.Metadata.Name}}
**{{.Doc.Metadata.Name}}**
{{- end}}
{{- if and .Doc.Metadata .Doc.Metadata.Period}}
{{.Doc.Metadata.Period}}
{{- end}}

{{- if and .Doc.Metadata .Doc.Metadata.Owner}}
Contact: {{.Doc.Metadata.Owner}}
{{- end}}

{{- if and .Doc.Metadata .Doc.Metadata.ReviewCadence}}
**Review Cadence:** {{.Doc.Metadata.ReviewCadence}}
{{- end}}

`))
