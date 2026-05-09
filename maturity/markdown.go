package maturity

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// MarkdownOptions configures maturity model markdown generation.
type MarkdownOptions struct {
	Title           string
	Author          string
	Date            string
	IncludeYAMLMeta bool     // Include YAML front matter for Pandoc/MkDocs
	IncludeTOC      bool     // Include table of contents
	ViewType        string   // "domain", "framework", "both"
	IncludeDetails  bool     // Include criterion details
	Frameworks      []string // Filter to specific frameworks (empty = all)
}

// DefaultMarkdownOptions returns sensible defaults.
func DefaultMarkdownOptions() *MarkdownOptions {
	return &MarkdownOptions{
		Title:           "Maturity Model",
		Date:            time.Now().Format("2006-01-02"),
		IncludeYAMLMeta: true,
		IncludeTOC:      true,
		ViewType:        "both",
		IncludeDetails:  true,
	}
}

// GenerateMarkdown generates a Markdown document from a maturity spec.
func (s *Spec) GenerateMarkdown(opts *MarkdownOptions) string {
	if opts == nil {
		opts = DefaultMarkdownOptions()
	}

	var sb strings.Builder

	// YAML front matter
	if opts.IncludeYAMLMeta {
		sb.WriteString("---\n")
		title := opts.Title
		if title == "" && s.Metadata != nil {
			title = s.Metadata.Name
		}
		sb.WriteString(fmt.Sprintf("title: \"%s\"\n", escapeYAML(title)))
		if opts.Author != "" {
			sb.WriteString(fmt.Sprintf("author: \"%s\"\n", escapeYAML(opts.Author)))
		}
		sb.WriteString(fmt.Sprintf("date: \"%s\"\n", opts.Date))
		sb.WriteString("---\n\n")
	}

	// Title
	title := opts.Title
	if title == "" && s.Metadata != nil {
		title = s.Metadata.Name
	}
	sb.WriteString(fmt.Sprintf("# %s\n\n", title))

	// Description
	if s.Metadata != nil && s.Metadata.Description != "" {
		sb.WriteString(fmt.Sprintf("> %s\n\n", s.Metadata.Description))
	}

	// Table of contents
	if opts.IncludeTOC {
		s.writeTOC(&sb, opts)
	}

	// SLI Catalog (if SLIs are defined)
	if len(s.SLIs) > 0 {
		s.writeSLICatalog(&sb)
	}

	// Generate views
	switch opts.ViewType {
	case "domain":
		s.writeDomainView(&sb, opts)
	case "framework":
		s.writeFrameworkView(&sb, opts)
	default: // "both"
		s.writeDomainView(&sb, opts)
		sb.WriteString("\n---\n\n")
		s.writeFrameworkView(&sb, opts)
	}

	return sb.String()
}

func (s *Spec) writeTOC(sb *strings.Builder, opts *MarkdownOptions) {
	sb.WriteString("## Table of Contents\n\n")

	// SLI Catalog link
	if len(s.SLIs) > 0 {
		sb.WriteString("- [SLI Catalog](#sli-catalog)\n")
	}

	domainNames := s.sortedDomainNames()

	if opts.ViewType == "domain" || opts.ViewType == "both" {
		sb.WriteString("### By Domain\n\n")
		for _, name := range domainNames {
			domain := s.Domains[name]
			anchor := strings.ToLower(strings.ReplaceAll(domain.Name, " ", "-"))
			sb.WriteString(fmt.Sprintf("- [%s](#%s)\n", domain.Name, anchor))
			for _, level := range domain.Levels {
				levelAnchor := fmt.Sprintf("%s-level-%d", anchor, level.Level)
				sb.WriteString(fmt.Sprintf("  - [Level %d: %s](#%s)\n", level.Level, level.Name, levelAnchor))
			}
		}
		sb.WriteString("\n")
	}

	if opts.ViewType == "framework" || opts.ViewType == "both" {
		frameworks := s.collectFrameworks(opts.Frameworks)
		if len(frameworks) > 0 {
			sb.WriteString("### By Framework\n\n")
			for _, fw := range frameworks {
				anchor := strings.ToLower(strings.ReplaceAll(fw, "_", "-"))
				sb.WriteString(fmt.Sprintf("- [%s](#framework-%s)\n", formatFrameworkName(fw), anchor))
			}
			sb.WriteString("\n")
		}
	}
}

func (s *Spec) writeSLICatalog(sb *strings.Builder) {
	sb.WriteString("## SLI Catalog\n\n")
	sb.WriteString("Service Level Indicators (SLIs) define the metrics being measured. ")
	sb.WriteString("Framework mappings are defined at the SLI level and inherited by all criteria.\n\n")

	// Group SLIs by category
	byCategory := make(map[string][]*SLI)
	var categories []string
	categorySet := make(map[string]bool)

	for _, sli := range s.SLIs {
		cat := sli.Category
		if cat == "" {
			cat = "Uncategorized"
		}
		if !categorySet[cat] {
			categorySet[cat] = true
			categories = append(categories, cat)
		}
		byCategory[cat] = append(byCategory[cat], sli)
	}

	// Sort categories
	sort.Strings(categories)

	// Write each category
	for _, cat := range categories {
		slis := byCategory[cat]

		// Sort SLIs by name within category
		sort.Slice(slis, func(i, j int) bool {
			return slis[i].Name < slis[j].Name
		})

		sb.WriteString(fmt.Sprintf("### %s\n\n", formatCategoryName(cat)))
		sb.WriteString("| SLI | Metric | Unit | Type | Frameworks |\n")
		sb.WriteString("|-----|--------|------|------|------------|\n")

		for _, sli := range slis {
			unit := sli.Unit
			if unit == "" {
				unit = "-"
			}
			sliType := "Quantitative"
			if sli.Type == CriterionTypeQualitative {
				sliType = "Qualitative"
			}
			frameworks := formatSLIFrameworks(sli.FrameworkMappings)
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				sli.Name, sli.MetricName, unit, sliType, frameworks))
		}
		sb.WriteString("\n")
	}
}

func formatCategoryName(cat string) string {
	switch cat {
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
	case "quality":
		return "Quality"
	default:
		// Title case the first letter
		if len(cat) > 0 {
			return strings.ToUpper(cat[:1]) + cat[1:]
		}
		return cat
	}
}

func formatSLIFrameworks(mappings []FrameworkMapping) string {
	if len(mappings) == 0 {
		return "-"
	}
	var parts []string
	for _, fm := range mappings {
		parts = append(parts, fmt.Sprintf("%s:%s", fm.Framework, fm.Reference))
	}
	return strings.Join(parts, ", ")
}

func (s *Spec) writeDomainView(sb *strings.Builder, opts *MarkdownOptions) {
	sb.WriteString("# Maturity Model by Domain\n\n")

	domainNames := s.sortedDomainNames()

	for _, name := range domainNames {
		domain := s.Domains[name]
		assessment := s.Assessments[name]

		// Domain header
		sb.WriteString(fmt.Sprintf("## %s\n\n", domain.Name))

		if domain.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", domain.Description))
		}

		// Current assessment
		if assessment != nil {
			sb.WriteString(fmt.Sprintf("**Current Level:** %d (%s)  \n", assessment.CurrentLevel, levelName(assessment.CurrentLevel)))
			sb.WriteString(fmt.Sprintf("**Target Level:** %d (%s)\n\n", assessment.TargetLevel, levelName(assessment.TargetLevel)))
		}

		// Levels
		for _, level := range domain.Levels {
			anchor := strings.ToLower(strings.ReplaceAll(domain.Name, " ", "-"))
			sb.WriteString(fmt.Sprintf("### Level %d: %s {#%s-level-%d}\n\n", level.Level, level.Name, anchor, level.Level))

			if level.Description != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", level.Description))
			}

			// Criteria table
			if len(level.Criteria) > 0 {
				sb.WriteString("#### Criteria (SLOs)\n\n")
				sb.WriteString("| ID | Name | Type | Target | Status | Frameworks |\n")
				sb.WriteString("|----|------|------|--------|--------|------------|\n")

				for _, c := range level.Criteria {
					criterionType := "Quantitative"
					target := c.TargetString()
					if c.IsQualitative() {
						criterionType = "Qualitative"
						target = "Tracked"
					}

					status := "⏳"
					if assessment != nil {
						if c.IsQualitative() {
							if assessment.CriteriaStatus != nil {
								if st, ok := assessment.CriteriaStatus[c.ID]; ok && IsQualitativeStatusMet(st) {
									status = "✅"
								}
							}
						} else if assessment.CriteriaValues != nil {
							if v, ok := assessment.CriteriaValues[c.ID]; ok && c.CheckMet(v) {
								status = "✅"
							}
						}
					}

					// Resolve framework mappings from SLI if needed
					fwMappings := c.GetFrameworkMappings(s)
					frameworks := formatCriterionFrameworks(fwMappings)
					sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
						c.ID, c.Name, criterionType, target, status, frameworks))
				}
				sb.WriteString("\n")

				// Detailed criterion info
				if opts.IncludeDetails {
					for _, c := range level.Criteria {
						fwMappings := c.GetFrameworkMappings(s)
						if len(fwMappings) > 0 || c.Description != "" {
							sb.WriteString(fmt.Sprintf("##### %s\n\n", c.Name))
							if c.Description != "" {
								sb.WriteString(fmt.Sprintf("%s\n\n", c.Description))
							}
							if len(fwMappings) > 0 {
								sb.WriteString("**Framework Mappings:**\n\n")
								for _, fm := range fwMappings {
									baseline := ""
									if fm.Baseline != "" {
										baseline = fmt.Sprintf(" (%s)", fm.Baseline)
									}
									name := fm.Name
									if name == "" {
										name = fm.Reference
									}
									sb.WriteString(fmt.Sprintf("- **%s**: %s%s\n", formatFrameworkName(fm.Framework), name, baseline))
								}
								sb.WriteString("\n")
							}
						}
					}
				}
			}

			// Enablers
			if len(level.Enablers) > 0 {
				sb.WriteString("#### Enablers\n\n")
				sb.WriteString("| ID | Name | Type | Status |\n")
				sb.WriteString("|----|------|------|--------|\n")

				for _, e := range level.Enablers {
					status := e.Status
					if assessment != nil && assessment.EnablerStatus != nil {
						if s, ok := assessment.EnablerStatus[e.ID]; ok {
							status = s
						}
					}
					sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
						e.ID, e.Name, e.Type, formatEnablerStatus(status)))
				}
				sb.WriteString("\n")
			}
		}
	}
}

func (s *Spec) writeFrameworkView(sb *strings.Builder, opts *MarkdownOptions) {
	frameworks := s.collectFrameworks(opts.Frameworks)
	if len(frameworks) == 0 {
		return
	}

	sb.WriteString("# Maturity Model by Framework\n\n")
	sb.WriteString("This view shows how maturity criteria map to compliance framework controls.\n\n")

	for _, fw := range frameworks {
		anchor := strings.ToLower(strings.ReplaceAll(fw, "_", "-"))
		sb.WriteString(fmt.Sprintf("## %s {#framework-%s}\n\n", formatFrameworkName(fw), anchor))

		// Collect all criteria mapped to this framework
		type criterionRef struct {
			Domain    string
			Level     int
			Criterion *Criterion
			Mapping   FrameworkMapping
		}

		var refs []criterionRef
		for domainName, domain := range s.Domains {
			for _, level := range domain.Levels {
				for i := range level.Criteria {
					c := &level.Criteria[i]
					// Resolve framework mappings from SLI if needed
					for _, fm := range c.GetFrameworkMappings(s) {
						if fm.Framework == fw {
							refs = append(refs, criterionRef{
								Domain:    domainName,
								Level:     level.Level,
								Criterion: c,
								Mapping:   fm,
							})
						}
					}
				}
			}
		}

		if len(refs) == 0 {
			sb.WriteString("No criteria mapped to this framework.\n\n")
			continue
		}

		// Sort by control ID
		sort.Slice(refs, func(i, j int) bool {
			return refs[i].Mapping.Reference < refs[j].Mapping.Reference
		})

		// Group by control ID
		sb.WriteString("| Control | Name | Domain | Level | Criterion | Status |\n")
		sb.WriteString("|---------|------|--------|-------|-----------|--------|\n")

		for _, ref := range refs {
			controlName := ref.Mapping.Name
			if controlName == "" {
				controlName = "-"
			}

			baseline := ""
			if ref.Mapping.Baseline != "" {
				baseline = fmt.Sprintf(" [%s]", ref.Mapping.Baseline)
			}

			status := "⏳"
			assessment := s.Assessments[ref.Domain]
			if assessment != nil {
				if ref.Criterion.IsQualitative() {
					if assessment.CriteriaStatus != nil {
						if st, ok := assessment.CriteriaStatus[ref.Criterion.ID]; ok && IsQualitativeStatusMet(st) {
							status = "✅"
						}
					}
				} else if assessment.CriteriaValues != nil {
					if v, ok := assessment.CriteriaValues[ref.Criterion.ID]; ok && ref.Criterion.CheckMet(v) {
						status = "✅"
					}
				}
			}

			sb.WriteString(fmt.Sprintf("| %s%s | %s | %s | M%d | %s | %s |\n",
				ref.Mapping.Reference, baseline, controlName, ref.Domain, ref.Level, ref.Criterion.Name, status))
		}
		sb.WriteString("\n")

		// Summary
		met := 0
		for _, ref := range refs {
			assessment := s.Assessments[ref.Domain]
			if assessment != nil {
				if ref.Criterion.IsQualitative() {
					if assessment.CriteriaStatus != nil {
						if st, ok := assessment.CriteriaStatus[ref.Criterion.ID]; ok && IsQualitativeStatusMet(st) {
							met++
						}
					}
				} else if assessment.CriteriaValues != nil {
					if v, ok := assessment.CriteriaValues[ref.Criterion.ID]; ok && ref.Criterion.CheckMet(v) {
						met++
					}
				}
			}
		}
		pct := float64(met) / float64(len(refs)) * 100
		sb.WriteString(fmt.Sprintf("**Coverage:** %d/%d controls satisfied (%.0f%%)\n\n", met, len(refs), pct))
	}
}

func (s *Spec) sortedDomainNames() []string {
	var names []string
	for name := range s.Domains {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (s *Spec) collectFrameworks(filter []string) []string {
	frameworkSet := make(map[string]bool)

	for _, domain := range s.Domains {
		for _, level := range domain.Levels {
			for _, c := range level.Criteria {
				// Resolve framework mappings from SLI if needed
				for _, fm := range c.GetFrameworkMappings(s) {
					if len(filter) == 0 || contains(filter, fm.Framework) {
						frameworkSet[fm.Framework] = true
					}
				}
			}
		}
	}

	var frameworks []string
	for fw := range frameworkSet {
		frameworks = append(frameworks, fw)
	}
	sort.Strings(frameworks)
	return frameworks
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func formatFrameworkName(fw string) string {
	switch fw {
	case "NIST_CSF":
		return "NIST CSF 1.1"
	case "NIST_CSF_2":
		return "NIST CSF 2.0"
	case "NIST_800_53":
		return "NIST SP 800-53"
	case "NIST_RMF":
		return "NIST RMF"
	case "NIST_AI_RMF":
		return "NIST AI RMF"
	case "NIST_800_171":
		return "NIST SP 800-171"
	case "FEDRAMP":
		return "FedRAMP"
	case "FEDRAMP_HIGH":
		return "FedRAMP High"
	case "FEDRAMP_MOD":
		return "FedRAMP Moderate"
	case "FEDRAMP_LOW":
		return "FedRAMP Low"
	case "MITRE_ATTACK":
		return "MITRE ATT&CK"
	case "CIS_CONTROLS":
		return "CIS Controls"
	case "SOC_2":
		return "SOC 2"
	case "ISO_27001":
		return "ISO 27001"
	case "DORA":
		return "DORA"
	case "SRE":
		return "SRE"
	default:
		return fw
	}
}

func formatCriterionFrameworks(mappings []FrameworkMapping) string {
	if len(mappings) == 0 {
		return "-"
	}
	var parts []string
	for _, fm := range mappings {
		parts = append(parts, fmt.Sprintf("%s:%s", fm.Framework, fm.Reference))
	}
	return strings.Join(parts, ", ")
}

func formatEnablerStatus(status string) string {
	switch status {
	case StatusCompleted:
		return "✅ Completed"
	case StatusInProgress:
		return "🔄 In Progress"
	case StatusBlocked:
		return "🚫 Blocked"
	case StatusNotStarted, "":
		return "⏳ Not Started"
	default:
		return status
	}
}

func levelName(level int) string {
	names := DefaultLevelNames()
	if name, ok := names[level]; ok {
		return name
	}
	return fmt.Sprintf("Level %d", level)
}

func escapeYAML(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
