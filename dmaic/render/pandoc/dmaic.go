// Package pandoc provides a Pandoc markdown renderer for DMAIC documents.
// The output includes YAML frontmatter configured for PDF generation via LuaLaTeX
// with sans-serif fonts and proper margins.
//
// To convert the output to PDF:
//
//	pandoc output.md -o output.pdf --pdf-engine=lualatex
package pandoc

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/grokify/prism-intelligence/dmaic"
	"github.com/grokify/prism-intelligence/dmaic/render"
)

// Renderer implements the render.Renderer interface for Pandoc Markdown output.
type Renderer struct{}

// New creates a new Pandoc Markdown renderer.
func New() *Renderer {
	return &Renderer{}
}

// Format returns the output format name.
func (r *Renderer) Format() string {
	return "pandoc"
}

// FileExtension returns the file extension for Pandoc output.
func (r *Renderer) FileExtension() string {
	return ".md"
}

// Render converts a DMAIC document to Pandoc Markdown with YAML frontmatter.
func (r *Renderer) Render(doc *dmaic.DMAICDocument, opts *render.Options) ([]byte, error) {
	if opts == nil {
		opts = render.DefaultOptions()
	}

	data := &templateData{
		Doc:            doc,
		Options:        opts,
		OverallHealth:  doc.CalculateOverallHealth(),
		HasInitiatives: opts.IncludeInitiatives && len(doc.Initiatives) > 0,
		GeneratedAt:    time.Now().Format("2006-01-02"),
	}

	var buf bytes.Buffer

	// Render YAML frontmatter
	if err := frontMatterTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering front matter: %w", err)
	}

	// Render title section
	if err := titleSectionTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering title section: %w", err)
	}

	// Render executive summary
	if err := executiveSummaryTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering executive summary: %w", err)
	}

	// Render each category
	for i, cat := range doc.Categories {
		catData := &categoryData{
			templateData: data,
			Category:     cat,
			Index:        i + 1,
		}
		if err := categorySectionTmpl.Execute(&buf, catData); err != nil {
			return nil, fmt.Errorf("rendering category section %d: %w", i+1, err)
		}
	}

	// Render DMAIC phase summary if grouping by phase
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
			if hasRootCauses {
				break
			}
		}
		if hasRootCauses {
			if err := rootCausesTmpl.Execute(&buf, data); err != nil {
				return nil, fmt.Errorf("rendering root causes: %w", err)
			}
		}
	}

	// Render initiatives
	if data.HasInitiatives {
		if err := initiativesTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering initiatives: %w", err)
		}
	}

	// Render appendix with metric details if showing capability metrics
	if opts.ShowCapabilityMetrics {
		if err := appendixTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering appendix: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// templateData holds data for template rendering.
type templateData struct {
	Doc            *dmaic.DMAICDocument
	Options        *render.Options
	OverallHealth  float64
	HasInitiatives bool
	GeneratedAt    string
}

// categoryData holds data for category section rendering.
type categoryData struct {
	*templateData
	Category dmaic.Category
	Index    int
}

// funcMap provides template functions for rendering.
var funcMap = template.FuncMap{
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
	"statusEmoji": func(status string) string {
		switch status {
		case dmaic.StatusGreen:
			return "\\textcolor{green}{$\\bullet$}"
		case dmaic.StatusYellow:
			return "\\textcolor{yellow}{$\\bullet$}"
		case dmaic.StatusRed:
			return "\\textcolor{red}{$\\bullet$}"
		default:
			return "$\\circ$"
		}
	},
	"statusText": func(status string) string {
		switch status {
		case dmaic.StatusGreen:
			return "Green"
		case dmaic.StatusYellow:
			return "Yellow"
		case dmaic.StatusRed:
			return "Red"
		default:
			return "Unknown"
		}
	},
	"phaseAbbrev": func(phase string) string {
		switch phase {
		case dmaic.PhaseDefine:
			return "D"
		case dmaic.PhaseMeasure:
			return "M"
		case dmaic.PhaseAnalyze:
			return "A"
		case dmaic.PhaseImprove:
			return "I"
		case dmaic.PhaseControl:
			return "C"
		default:
			return "?"
		}
	},
	"calculateStatus": func(m dmaic.Metric) string {
		return m.CalculateStatus()
	},
	"categoryHealth": func(c dmaic.Category) float64 {
		return c.CalculateCategoryHealth()
	},
	"truncate": func(s string, max int) string {
		if len(s) <= max {
			return s
		}
		return s[:max-3] + "..."
	},
	"escapeLatex": func(s string) string {
		// Escape special LaTeX characters
		replacer := strings.NewReplacer(
			"&", "\\&",
			"%", "\\%",
			"$", "\\$",
			"#", "\\#",
			"_", "\\_",
			"{", "\\{",
			"}", "\\}",
			"~", "\\textasciitilde{}",
			"^", "\\textasciicircum{}",
		)
		return replacer.Replace(s)
	},
	"add": func(a, b int) int {
		return a + b
	},
}

// Templates

var frontMatterTmpl = template.Must(template.New("frontMatter").Funcs(funcMap).Parse(`---
title: "{{if and .Doc.Metadata .Doc.Metadata.Name}}{{escapeLatex .Doc.Metadata.Name}}{{else}}DMAIC Metrics Report{{end}}"
{{- if and .Doc.Metadata .Doc.Metadata.Owner}}
author: "{{escapeLatex .Doc.Metadata.Owner}}"
{{- end}}
date: "{{.GeneratedAt}}"
documentclass: article
papersize: a4
geometry:
  - margin=2cm
fontsize: 11pt
mainfont: "DejaVu Sans"
sansfont: "DejaVu Sans"
monofont: "DejaVu Sans Mono"
mathfont: "DejaVu Math TeX Gyre"
fontfamily: helvet
fontfamilyoptions: scaled
header-includes:
  - \usepackage[T1]{fontenc}
  - \usepackage{helvet}
  - \renewcommand{\familydefault}{\sfdefault}
  - \usepackage{xcolor}
  - \definecolor{green}{RGB}{56,161,105}
  - \definecolor{yellow}{RGB}{214,158,46}
  - \definecolor{red}{RGB}{229,62,62}
  - \usepackage{booktabs}
  - \usepackage{longtable}
  - \usepackage{array}
  - \usepackage{multirow}
  - \usepackage{float}
  - \floatplacement{table}{H}
colorlinks: true
linkcolor: blue
urlcolor: blue
toc: true
toc-depth: 3
numbersections: true
---

`))

var titleSectionTmpl = template.Must(template.New("titleSection").Funcs(funcMap).Parse(`
# Executive Overview

{{- if and .Doc.Metadata .Doc.Metadata.Description}}

{{.Doc.Metadata.Description}}
{{- end}}

| Field | Value |
|:------|:------|
{{- if and .Doc.Metadata .Doc.Metadata.Team}}
| **Team** | {{escapeLatex .Doc.Metadata.Team}} |
{{- end}}
{{- if and .Doc.Metadata .Doc.Metadata.Period}}
| **Period** | {{.Doc.Metadata.Period}} |
{{- end}}
{{- if and .Doc.Metadata .Doc.Metadata.Status}}
| **Status** | {{.Doc.Metadata.Status}} |
{{- end}}
{{- if and .Doc.Metadata .Doc.Metadata.ReviewCadence}}
| **Review Cadence** | {{.Doc.Metadata.ReviewCadence}} |
{{- end}}
| **Overall Health** | {{healthPercent .OverallHealth}} |
| **Total Categories** | {{len .Doc.Categories}} |
| **Total Metrics** | {{len .Doc.AllMetrics}} |

`))

var executiveSummaryTmpl = template.Must(template.New("executiveSummary").Funcs(funcMap).Parse(`
## Health Summary by Category

| Category | Health | Metrics | Green | Yellow | Red |
|:---------|-------:|--------:|------:|-------:|----:|
{{- range .Doc.Categories}}
{{- $health := categoryHealth .}}
{{- $green := 0}}{{$yellow := 0}}{{$red := 0}}
{{- range .Metrics}}
{{- $status := calculateStatus .}}
{{- if eq $status "Green"}}{{$green = add $green 1}}{{end}}
{{- if eq $status "Yellow"}}{{$yellow = add $yellow 1}}{{end}}
{{- if eq $status "Red"}}{{$red = add $red 1}}{{end}}
{{- end}}
| {{escapeLatex .Name}} | {{healthPercent $health}} | {{len .Metrics}} | {{$green}} | {{$yellow}} | {{$red}} |
{{- end}}

`))

var categorySectionTmpl = template.Must(template.New("categorySection").Funcs(funcMap).Parse(`
# {{.Category.Name}}

{{- if .Category.Description}}

{{.Category.Description}}
{{- end}}

{{- if .Category.Owner}}

**Owner:** {{escapeLatex .Category.Owner}}
{{- end}}

**Category Health:** {{healthPercent (categoryHealth .Category)}}

## Metrics

| Metric | Phase | Baseline | Current | Target | Unit | Status |
|:-------|:-----:|---------:|--------:|-------:|:-----|:------:|
{{- range .Category.Metrics}}
| {{escapeLatex (truncate .Name 30)}} | {{phaseAbbrev .Phase}} | {{formatValue .Baseline .Unit}} | {{formatValue .Current .Unit}} | {{formatValue .Target .Unit}} | {{if .Unit}}{{.Unit}}{{else}}-{{end}} | {{statusText (calculateStatus .)}} |
{{- end}}

{{- range .Category.Metrics}}
{{- $status := calculateStatus .}}
{{- if or (eq $status "Yellow") (eq $status "Red")}}

### {{escapeLatex .Name}} ({{$status}})

{{- if .Description}}

{{.Description}}
{{- end}}

| Measure | Value |
|:--------|------:|
| Baseline | {{formatValue .Baseline .Unit}} |
| Current | {{formatValue .Current .Unit}} |
| Target | {{formatValue .Target .Unit}} |
| Trend Direction | {{if .TrendDirection}}{{.TrendDirection}}{{else}}N/A{{end}} |
{{- if .Thresholds}}
| Warning Threshold | {{formatValue .Thresholds.Warning .Unit}} |
| Critical Threshold | {{formatValue .Thresholds.Critical .Unit}} |
{{- end}}

{{- if .ControlLimits}}

**Control Limits:**

| Limit | Value |
|:------|------:|
| Upper Control Limit (UCL) | {{formatValue .ControlLimits.UCL .Unit}} |
| Center Line | {{formatValue .ControlLimits.CenterLine .Unit}} |
| Lower Control Limit (LCL) | {{formatValue .ControlLimits.LCL .Unit}} |
{{- if .ControlLimits.Sigma}}
| Sigma | {{formatValue .ControlLimits.Sigma ""}} |
{{- end}}
{{- end}}

{{- if .RootCauses}}

**Root Causes:**

| Description | Category | Impact | Validated |
|:------------|:---------|:-------|:---------:|
{{- range .RootCauses}}
| {{escapeLatex (truncate .Description 40)}} | {{.Category}} | {{.Impact}} | {{if .Validated}}Yes{{else}}No{{end}} |
{{- end}}
{{- end}}
{{- end}}
{{- end}}

`))

var phaseSummaryTmpl = template.Must(template.New("phaseSummary").Funcs(funcMap).Parse(`
# DMAIC Phase Distribution

| Phase | Description | Count |
|:------|:------------|------:|
| **D** - Define | Define the problem and project goals | {{len (index .Doc.MetricsByPhase "Define")}} |
| **M** - Measure | Measure current performance | {{len (index .Doc.MetricsByPhase "Measure")}} |
| **A** - Analyze | Analyze root causes | {{len (index .Doc.MetricsByPhase "Analyze")}} |
| **I** - Improve | Implement improvements | {{len (index .Doc.MetricsByPhase "Improve")}} |
| **C** - Control | Sustain the gains | {{len (index .Doc.MetricsByPhase "Control")}} |

`))

var rootCausesTmpl = template.Must(template.New("rootCauses").Funcs(funcMap).Parse(`
# Root Cause Analysis

| Metric | Root Cause | Category | Impact | Validated |
|:-------|:-----------|:---------|:-------|:---------:|
{{- range .Doc.Categories}}
{{- range .Metrics}}
{{- $metricName := .Name}}
{{- range .RootCauses}}
| {{escapeLatex (truncate $metricName 20)}} | {{escapeLatex (truncate .Description 35)}} | {{.Category}} | {{.Impact}} | {{if .Validated}}Yes{{else}}No{{end}} |
{{- end}}
{{- end}}
{{- end}}

`))

var initiativesTmpl = template.Must(template.New("initiatives").Funcs(funcMap).Parse(`
# Improvement Initiatives

{{range .Doc.Initiatives}}
## {{escapeLatex .Name}}

{{- if .Description}}

{{.Description}}
{{- end}}

| Field | Value |
|:------|:------|
{{- if .Owner}}
| **Owner** | {{escapeLatex .Owner}} |
{{- end}}
| **Status** | {{if .Status}}{{.Status}}{{else}}N/A{{end}} |
{{- if .StartDate}}
| **Start Date** | {{.StartDate}} |
{{- end}}
{{- if .EndDate}}
| **End Date** | {{.EndDate}} |
{{- end}}
{{- if .ExpectedImpact}}
| **Expected Impact** | {{escapeLatex .ExpectedImpact}} |
{{- end}}
{{- if .ActualImpact}}
| **Actual Impact** | {{escapeLatex .ActualImpact}} |
{{- end}}

{{- if .MetricIDs}}

**Linked Metrics:** {{range $i, $id := .MetricIDs}}{{if $i}}, {{end}}{{$id}}{{end}}
{{- end}}

{{end}}
`))

var appendixTmpl = template.Must(template.New("appendix").Funcs(funcMap).Parse(`
# Appendix: Process Capability Metrics

| Metric | Cp | Cpk | Sigma Level | DPMO |
|:-------|---:|----:|------------:|-----:|
{{- range .Doc.Categories}}
{{- range .Metrics}}
{{- if .ProcessCapability}}
| {{escapeLatex (truncate .Name 25)}} | {{formatValue .ProcessCapability.Cp ""}} | {{formatValue .ProcessCapability.Cpk ""}} | {{formatValue .ProcessCapability.SigmaLevel ""}} | {{formatValue .ProcessCapability.DPMO ""}} |
{{- end}}
{{- end}}
{{- end}}

## Control Limits Summary

| Metric | UCL | Center | LCL | In Control |
|:-------|----:|-------:|----:|:----------:|
{{- range .Doc.Categories}}
{{- range .Metrics}}
{{- if .ControlLimits}}
| {{escapeLatex (truncate .Name 25)}} | {{formatValue .ControlLimits.UCL .Unit}} | {{formatValue .ControlLimits.CenterLine .Unit}} | {{formatValue .ControlLimits.LCL .Unit}} | {{if .IsInControl}}Yes{{else}}No{{end}} |
{{- end}}
{{- end}}
{{- end}}

`))
