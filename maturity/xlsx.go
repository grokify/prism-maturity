package maturity

import (
	"fmt"
	"sort"
	"strings"

	"github.com/plexusone/omniframe"
	"github.com/xuri/excelize/v2"
)

// XLSXGenerator generates Excel reports from maturity specs.
type XLSXGenerator struct {
	spec *Spec
	file *excelize.File
}

// NewXLSXGenerator creates a new XLSX generator.
func NewXLSXGenerator(spec *Spec) *XLSXGenerator {
	return &XLSXGenerator{
		spec: spec,
		file: excelize.NewFile(),
	}
}

// Generate creates the XLSX file with all sheets.
func (g *XLSXGenerator) Generate() error {
	// Create sheets
	if err := g.createRequirementsSheet(); err != nil {
		return fmt.Errorf("failed to create requirements sheet: %w", err)
	}

	if err := g.createSLOsSheet(); err != nil {
		return fmt.Errorf("failed to create SLOs sheet: %w", err)
	}

	if err := g.createFrameworkMappingsSheet(); err != nil {
		return fmt.Errorf("failed to create framework mappings sheet: %w", err)
	}

	if err := g.createProgressSheet(); err != nil {
		return fmt.Errorf("failed to create progress sheet: %w", err)
	}

	if err := g.createLevelDefinitionsSheet(); err != nil {
		return fmt.Errorf("failed to create level definitions sheet: %w", err)
	}

	// Delete the default "Sheet1"
	_ = g.file.DeleteSheet("Sheet1")

	return nil
}

// SaveAs saves the XLSX file to the specified path.
func (g *XLSXGenerator) SaveAs(filename string) error {
	return g.file.SaveAs(filename)
}

// createRequirementsSheet creates the Requirements (Enablers) sheet.
func (g *XLSXGenerator) createRequirementsSheet() error {
	sheetName := "Requirements"
	index, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}
	g.file.SetActiveSheet(index)

	// Headers
	headers := []string{
		"ID", "Domain", "Level", "Name", "Description", "Type",
		"Layer", "Team", "Effort", "Status", "Enables", "Depends On",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	g.setCellStyle(sheetName, "A1", "L1", headerStyle)

	// Data rows
	row := 2
	domainNames := g.sortedDomainNames()

	for _, domainName := range domainNames {
		domain := g.spec.Domains[domainName]
		for _, level := range domain.Levels {
			for _, e := range level.Enablers {
				g.setCellValue(sheetName, fmt.Sprintf("A%d", row), e.ID)
				g.setCellValue(sheetName, fmt.Sprintf("B%d", row), domainName)
				g.setCellValue(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("M%d", level.Level))
				g.setCellValue(sheetName, fmt.Sprintf("D%d", row), e.Name)
				g.setCellValue(sheetName, fmt.Sprintf("E%d", row), e.Description)
				g.setCellValue(sheetName, fmt.Sprintf("F%d", row), e.Type)
				g.setCellValue(sheetName, fmt.Sprintf("G%d", row), e.Layer)
				g.setCellValue(sheetName, fmt.Sprintf("H%d", row), e.Team)
				g.setCellValue(sheetName, fmt.Sprintf("I%d", row), e.Effort)

				// Get status from assessment if available
				status := e.Status
				if assessment, ok := g.spec.Assessments[domainName]; ok {
					if s, ok := assessment.EnablerStatus[e.ID]; ok {
						status = s
					}
				}
				g.setCellValue(sheetName, fmt.Sprintf("J%d", row), status)

				g.setCellValue(sheetName, fmt.Sprintf("K%d", row), strings.Join(e.CriteriaIDs, ", "))
				g.setCellValue(sheetName, fmt.Sprintf("L%d", row), strings.Join(e.DependsOn, ", "))

				// Color code status
				if status != "" {
					statusStyle := g.statusStyle(status)
					if statusStyle != 0 {
						g.setCellStyle(sheetName, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), statusStyle)
					}
				}

				row++
			}
		}
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 25, "B": 12, "C": 8, "D": 35, "E": 50, "F": 15,
		"G": 12, "H": 20, "I": 12, "J": 15, "K": 30, "L": 25,
	}
	for col, width := range colWidths {
		g.setColWidth(sheetName, col, col, width)
	}

	// Auto filter
	g.setAutoFilter(sheetName, "A1:L1")

	return nil
}

// createSLOsSheet creates the SLOs (Criteria) sheet with framework columns.
func (g *XLSXGenerator) createSLOsSheet() error {
	sheetName := "SLOs"
	_, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Collect all unique frameworks across all criteria (sorted alphabetically)
	frameworks := g.collectAllFrameworks()

	// Headers - base columns plus framework columns
	headers := []string{
		"ID", "Domain", "Level", "Name", "Metric", "Type", "Operator",
		"Target", "Unit", "Current", "Met", "Layer", "Category", "Required",
	}
	// Add framework columns
	for _, fw := range frameworks {
		headers = append(headers, fw)
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"548235"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	endHeaderCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
	g.setCellStyle(sheetName, "A1", endHeaderCell, headerStyle)

	// Data rows
	row := 2
	domainNames := g.sortedDomainNames()

	for _, domainName := range domainNames {
		domain := g.spec.Domains[domainName]

		// Get assessment for domain (may be nil)
		var assessment *DomainAssessment
		if g.spec.Assessments != nil {
			assessment = g.spec.Assessments[domainName]
		}

		for _, level := range domain.Levels {
			for _, c := range level.Criteria {
				g.setCellValue(sheetName, fmt.Sprintf("A%d", row), c.ID)
				g.setCellValue(sheetName, fmt.Sprintf("B%d", row), domainName)
				g.setCellValue(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("M%d", level.Level))
				g.setCellValue(sheetName, fmt.Sprintf("D%d", row), c.Name)
				g.setCellValue(sheetName, fmt.Sprintf("E%d", row), c.GetMetricName(g.spec))

				// Type column - Quantitative or Qualitative (resolve from SLI)
				isQual := c.IsQualitativeWithSpec(g.spec)
				criterionType := "Quantitative"
				if isQual {
					criterionType = "Qualitative"
				}
				g.setCellValue(sheetName, fmt.Sprintf("F%d", row), criterionType)

				g.setCellValue(sheetName, fmt.Sprintf("G%d", row), OperatorSymbol(c.Operator))

				// Target - different display for qualitative
				if isQual {
					g.setCellValue(sheetName, fmt.Sprintf("H%d", row), "Tracked")
				} else {
					g.setCellValue(sheetName, fmt.Sprintf("H%d", row), c.Target)
				}

				g.setCellValue(sheetName, fmt.Sprintf("I%d", row), c.GetUnit(g.spec))

				// Get current value/status from assessment
				var isMet bool
				if isQual {
					// For qualitative, check if status is set
					status := c.Status
					if assessment != nil && assessment.CriteriaStatus != nil {
						if s, ok := assessment.CriteriaStatus[c.ID]; ok {
							status = s
						}
					}
					isMet = IsQualitativeStatusMet(status)
					g.setCellValue(sheetName, fmt.Sprintf("J%d", row), formatQualitativeStatus(status))
				} else {
					// For quantitative, use numeric value
					var current float64
					if assessment != nil && assessment.CriteriaValues != nil {
						if v, ok := assessment.CriteriaValues[c.ID]; ok {
							current = v
							isMet = c.CheckMet(current)
						}
					}
					g.setCellValue(sheetName, fmt.Sprintf("J%d", row), current)
				}

				metStatus := "No"
				if isMet {
					metStatus = "Yes"
				}
				g.setCellValue(sheetName, fmt.Sprintf("K%d", row), metStatus)

				g.setCellValue(sheetName, fmt.Sprintf("L%d", row), c.GetLayer(g.spec))
				g.setCellValue(sheetName, fmt.Sprintf("M%d", row), c.GetCategory(g.spec))

				required := "Yes"
				if !c.Required && c.Weight > 0 {
					required = "No"
				}
				g.setCellValue(sheetName, fmt.Sprintf("N%d", row), required)

				// Color code met status
				metStyle := g.metStyle(isMet)
				if metStyle != 0 {
					g.setCellStyle(sheetName, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), metStyle)
				}

				// Framework columns - show control reference if mapped (resolve from SLI)
				frameworkRefs := make(map[string]string)
				for _, fm := range c.GetFrameworkMappings(g.spec) {
					frameworkRefs[fm.Framework] = fm.Reference
				}
				for i, fw := range frameworks {
					col, _ := excelize.CoordinatesToCellName(15+i, row) // Column O onwards
					if ref, ok := frameworkRefs[fw]; ok {
						g.setCellValue(sheetName, col, ref)
					} else {
						g.setCellValue(sheetName, col, "-")
					}
				}

				row++
			}
		}
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 25, "B": 12, "C": 8, "D": 30, "E": 35, "F": 12,
		"G": 10, "H": 12, "I": 10, "J": 12, "K": 8, "L": 12, "M": 12, "N": 10,
	}
	for col, width := range colWidths {
		g.setColWidth(sheetName, col, col, width)
	}
	// Set framework column widths
	for i := range frameworks {
		col, _ := excelize.ColumnNumberToName(15 + i)
		g.setColWidth(sheetName, col, col, 15)
	}

	// Auto filter
	endFilterCell, _ := excelize.CoordinatesToCellName(len(headers), 1)
	g.setAutoFilter(sheetName, "A1:"+endFilterCell)

	return nil
}

// collectAllFrameworks returns all unique frameworks across all criteria, sorted alphabetically.
// Resolves framework mappings from both inline criterion mappings and referenced SLIs.
func (g *XLSXGenerator) collectAllFrameworks() []string {
	frameworkSet := make(map[string]bool)

	for _, domain := range g.spec.Domains {
		for _, level := range domain.Levels {
			for _, c := range level.Criteria {
				// Use GetFrameworkMappings to resolve from SLI if needed
				for _, fm := range c.GetFrameworkMappings(g.spec) {
					frameworkSet[fm.Framework] = true
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

// createFrameworkMappingsSheet creates a detailed Framework Mappings sheet (Option 4).
func (g *XLSXGenerator) createFrameworkMappingsSheet() error {
	sheetName := "Framework Mappings"
	_, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	headers := []string{
		"SLO ID", "SLO Name", "Domain", "Level", "Framework", "Reference",
		"Control Name", "Baseline", "Version", "Status",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"7030A0"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	g.setCellStyle(sheetName, "A1", "J1", headerStyle)

	// Data rows - one row per SLO-framework mapping
	row := 2
	domainNames := g.sortedDomainNames()

	for _, domainName := range domainNames {
		domain := g.spec.Domains[domainName]

		// Get assessment for domain (may be nil)
		var assessment *DomainAssessment
		if g.spec.Assessments != nil {
			assessment = g.spec.Assessments[domainName]
		}

		for _, level := range domain.Levels {
			for _, c := range level.Criteria {
				// Get framework mappings (resolve from SLI if needed)
				frameworkMappings := c.GetFrameworkMappings(g.spec)
				if len(frameworkMappings) == 0 {
					continue
				}

				// Determine status
				var isMet bool
				if c.IsQualitativeWithSpec(g.spec) {
					status := c.Status
					if assessment != nil && assessment.CriteriaStatus != nil {
						if s, ok := assessment.CriteriaStatus[c.ID]; ok {
							status = s
						}
					}
					isMet = IsQualitativeStatusMet(status)
				} else if assessment != nil && assessment.CriteriaValues != nil {
					if v, ok := assessment.CriteriaValues[c.ID]; ok {
						isMet = c.CheckMet(v)
					}
				}

				metStatus := "Pending"
				if isMet {
					metStatus = "Met"
				}

				// Create a row for each framework mapping
				for _, fm := range frameworkMappings {
					g.setCellValue(sheetName, fmt.Sprintf("A%d", row), c.ID)
					g.setCellValue(sheetName, fmt.Sprintf("B%d", row), c.Name)
					g.setCellValue(sheetName, fmt.Sprintf("C%d", row), domainName)
					g.setCellValue(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("M%d", level.Level))
					g.setCellValue(sheetName, fmt.Sprintf("E%d", row), fm.Framework)
					g.setCellValue(sheetName, fmt.Sprintf("F%d", row), fm.Reference)
					g.setCellValue(sheetName, fmt.Sprintf("G%d", row), fm.Name)
					g.setCellValue(sheetName, fmt.Sprintf("H%d", row), fm.Baseline)
					g.setCellValue(sheetName, fmt.Sprintf("I%d", row), fm.Version)
					g.setCellValue(sheetName, fmt.Sprintf("J%d", row), metStatus)

					// Color code status
					statusStyle := g.frameworkStatusStyle(isMet)
					if statusStyle != 0 {
						g.setCellStyle(sheetName, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), statusStyle)
					}

					row++
				}
			}
		}
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 25, "B": 30, "C": 12, "D": 8, "E": 15, "F": 15,
		"G": 35, "H": 12, "I": 10, "J": 10,
	}
	for col, width := range colWidths {
		g.setColWidth(sheetName, col, col, width)
	}

	// Auto filter
	g.setAutoFilter(sheetName, "A1:J1")

	return nil
}

func (g *XLSXGenerator) frameworkStatusStyle(isMet bool) int {
	var color string
	if isMet {
		color = "C6EFCE" // Green
	} else {
		color = "FFEB9C" // Yellow (pending)
	}

	style, _ := g.file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
	})
	return style
}

// createProgressSheet creates the Progress Summary sheet.
func (g *XLSXGenerator) createProgressSheet() error {
	sheetName := "Progress"
	_, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Headers
	headers := []string{
		"Domain", "Current Level", "Target Level",
		"M2 Progress", "M3 Progress", "M4 Progress", "M5 Progress",
		"Next Actions",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"7030A0"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	g.setCellStyle(sheetName, "A1", "H1", headerStyle)

	// Data rows
	row := 2
	domainNames := g.sortedDomainNames()

	for _, domainName := range domainNames {
		domain := g.spec.Domains[domainName]

		// Get assessment for domain (may be nil)
		var assessment *DomainAssessment
		if g.spec.Assessments != nil {
			assessment = g.spec.Assessments[domainName]
		}

		g.setCellValue(sheetName, fmt.Sprintf("A%d", row), domain.Name)

		// Show current/target levels if assessment exists
		if assessment != nil {
			g.setCellValue(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("M%d", assessment.CurrentLevel))
			g.setCellValue(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("M%d", assessment.TargetLevel))
		} else {
			g.setCellValue(sheetName, fmt.Sprintf("B%d", row), "-")
			g.setCellValue(sheetName, fmt.Sprintf("C%d", row), "-")
		}

		// Calculate progress for each level
		for level := 2; level <= 5; level++ {
			col, _ := excelize.CoordinatesToCellName(level+2, row)

			levelDef, found := domain.GetLevel(level)
			if !found || len(levelDef.Criteria) == 0 {
				g.setCellValue(sheetName, col, "N/A")
				continue
			}

			var criteriaValues map[string]float64
			var enablerStatus map[string]string
			if assessment != nil {
				criteriaValues = assessment.CriteriaValues
				enablerStatus = assessment.EnablerStatus
			}

			progress := levelDef.CalculateLevelProgress(criteriaValues, enablerStatus)
			g.setCellValue(sheetName, col, fmt.Sprintf("%.0f%%", progress.ProgressPercent))

			// Color code based on progress
			progressStyle := g.progressStyle(progress.ProgressPercent)
			if progressStyle != 0 {
				g.setCellStyle(sheetName, col, col, progressStyle)
			}
		}

		// Next actions: list incomplete enablers for next level
		nextLevel := 1
		if assessment != nil {
			nextLevel = assessment.CurrentLevel + 1
		}
		if nextLevel <= 5 {
			enablers := domain.EnablersForLevel(nextLevel)
			var nextActions []string
			for _, e := range enablers {
				status := StatusNotStarted
				if assessment != nil && assessment.EnablerStatus != nil {
					if s, ok := assessment.EnablerStatus[e.ID]; ok {
						status = s
					}
				}
				if status != StatusCompleted {
					nextActions = append(nextActions, e.Name)
				}
			}
			if len(nextActions) > 3 {
				nextActions = nextActions[:3]
				nextActions = append(nextActions, "...")
			}
			g.setCellValue(sheetName, fmt.Sprintf("H%d", row), strings.Join(nextActions, "; "))
		}

		row++
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 15, "B": 15, "C": 15, "D": 12, "E": 12, "F": 12, "G": 12, "H": 60,
	}
	for col, width := range colWidths {
		g.setColWidth(sheetName, col, col, width)
	}

	return nil
}

// createLevelDefinitionsSheet creates the Level Definitions sheet.
func (g *XLSXGenerator) createLevelDefinitionsSheet() error {
	sheetName := "Level Definitions"
	_, err := g.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Build headers dynamically from domains
	domainNames := g.sortedDomainNames()
	headers := []string{"Level", "Name"}
	for _, d := range domainNames {
		headers = append(headers, g.spec.Domains[d].Name)
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		g.setCellValue(sheetName, cell, h)
	}

	// Style header row
	headerStyle, _ := g.file.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"305496"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	endCol, _ := excelize.CoordinatesToCellName(len(headers), 1)
	g.setCellStyle(sheetName, "A1", endCol, headerStyle)

	// Data rows for levels M1-M5
	levelNames := DefaultLevelNames()
	for level := 1; level <= 5; level++ {
		row := level + 1
		g.setCellValue(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("M%d", level))
		g.setCellValue(sheetName, fmt.Sprintf("B%d", row), levelNames[level])

		// Description for each domain
		for col, domainName := range domainNames {
			domain := g.spec.Domains[domainName]
			levelDef, found := domain.GetLevel(level)
			desc := ""
			if found {
				desc = levelDef.Description
			}
			cell, _ := excelize.CoordinatesToCellName(col+3, row)
			g.setCellValue(sheetName, cell, desc)
		}
	}

	// Set column widths
	g.setColWidth(sheetName, "A", "A", 8)
	g.setColWidth(sheetName, "B", "B", 12)
	for i := range domainNames {
		col, _ := excelize.ColumnNumberToName(i + 3)
		g.setColWidth(sheetName, col, col, 50)
	}

	// Enable text wrap for description columns
	wrapStyle, _ := g.file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "top"},
	})
	endCol, _ = excelize.CoordinatesToCellName(len(headers), 6)
	g.setCellStyle(sheetName, "C2", endCol, wrapStyle)

	return nil
}

// Helper methods

// setCellValue wraps SetCellValue and ignores errors for simplicity.
// Cell value setting errors in excelize are rare and non-fatal.
func (g *XLSXGenerator) setCellValue(sheet, cell string, value interface{}) {
	_ = g.file.SetCellValue(sheet, cell, value)
}

// setCellStyle wraps SetCellStyle and ignores errors.
func (g *XLSXGenerator) setCellStyle(sheet, startCell, endCell string, styleID int) {
	_ = g.file.SetCellStyle(sheet, startCell, endCell, styleID)
}

// setColWidth wraps SetColWidth and ignores errors.
func (g *XLSXGenerator) setColWidth(sheet, startCol, endCol string, width float64) {
	_ = g.file.SetColWidth(sheet, startCol, endCol, width)
}

// setAutoFilter wraps AutoFilter and ignores errors.
func (g *XLSXGenerator) setAutoFilter(sheet, rangeRef string) {
	_ = g.file.AutoFilter(sheet, rangeRef, nil)
}

func (g *XLSXGenerator) sortedDomainNames() []string {
	var names []string
	for name := range g.spec.Domains {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (g *XLSXGenerator) statusStyle(status string) int {
	var color string
	switch status {
	case StatusCompleted:
		color = "C6EFCE" // Green
	case StatusInProgress:
		color = "FFEB9C" // Yellow
	case StatusBlocked:
		color = "FFC7CE" // Red
	case StatusNotStarted:
		color = "DDDDDD" // Gray
	default:
		return 0
	}

	style, _ := g.file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
	})
	return style
}

func (g *XLSXGenerator) metStyle(isMet bool) int {
	var color string
	if isMet {
		color = "C6EFCE" // Green
	} else {
		color = "FFC7CE" // Red
	}

	style, _ := g.file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
	})
	return style
}

func (g *XLSXGenerator) progressStyle(percent float64) int {
	var color string
	switch {
	case percent >= 100:
		color = "C6EFCE" // Green
	case percent >= 75:
		color = "C6EFCE" // Light green
	case percent >= 50:
		color = "FFEB9C" // Yellow
	case percent >= 25:
		color = "FFCC99" // Orange
	default:
		color = "FFC7CE" // Red
	}

	style, _ := g.file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
	})
	return style
}

// GenerateXLSX is a convenience function to generate XLSX from a spec file.
func GenerateXLSX(specFile, outputFile string) error {
	spec, err := ReadSpecFile(specFile)
	if err != nil {
		return err
	}

	gen := NewXLSXGenerator(spec)
	if err := gen.Generate(); err != nil {
		return err
	}

	return gen.SaveAs(outputFile)
}

// GenerateSimpleXLSX generates a simple XLSX report using omniframe.
// This provides a basic export without conditional styling.
func GenerateSimpleXLSX(specFile, outputFile string) error {
	spec, err := ReadSpecFile(specFile)
	if err != nil {
		return err
	}

	fs := omniframe.NewFrameSet(spec.Metadata.Name)

	// Create Requirements frame
	reqFrame, err := buildRequirementsFrame(spec)
	if err != nil {
		return fmt.Errorf("failed to build requirements frame: %w", err)
	}
	if err := fs.AddFrame(reqFrame); err != nil {
		return err
	}

	// Create SLOs frame
	sloFrame, err := buildSLOsFrame(spec)
	if err != nil {
		return fmt.Errorf("failed to build SLOs frame: %w", err)
	}
	if err := fs.AddFrame(sloFrame); err != nil {
		return err
	}

	// Create Progress frame
	progressFrame, err := buildProgressFrame(spec)
	if err != nil {
		return fmt.Errorf("failed to build progress frame: %w", err)
	}
	if err := fs.AddFrame(progressFrame); err != nil {
		return err
	}

	return fs.WriteXLSX(outputFile)
}

func buildRequirementsFrame(spec *Spec) (*omniframe.Frame, error) {
	columns := []string{
		"ID", "Domain", "Level", "Name", "Description", "Type",
		"Layer", "Team", "Effort", "Status", "Enables", "Depends On",
	}

	var rows [][]any
	domainNames := sortedDomainNamesFromSpec(spec)

	for _, domainName := range domainNames {
		domain := spec.Domains[domainName]
		for _, level := range domain.Levels {
			for _, e := range level.Enablers {
				status := e.Status
				if assessment, ok := spec.Assessments[domainName]; ok {
					if s, ok := assessment.EnablerStatus[e.ID]; ok {
						status = s
					}
				}

				rows = append(rows, []any{
					e.ID,
					domainName,
					fmt.Sprintf("M%d", level.Level),
					e.Name,
					e.Description,
					e.Type,
					e.Layer,
					e.Team,
					e.Effort,
					status,
					strings.Join(e.CriteriaIDs, ", "),
					strings.Join(e.DependsOn, ", "),
				})
			}
		}
	}

	frame, err := omniframe.FromRows("Requirements", columns, rows)
	if err != nil {
		return nil, err
	}

	// Set column widths
	_ = frame.SetColumnWidth("ID", 25)
	_ = frame.SetColumnWidth("Name", 35)
	_ = frame.SetColumnWidth("Description", 50)

	return frame, nil
}

func buildSLOsFrame(spec *Spec) (*omniframe.Frame, error) {
	columns := []string{
		"ID", "Domain", "Level", "Name", "Metric", "Type", "Operator",
		"Target", "Unit", "Current", "Met", "Layer", "Category", "Required",
	}

	var rows [][]any
	domainNames := sortedDomainNamesFromSpec(spec)

	for _, domainName := range domainNames {
		domain := spec.Domains[domainName]
		assessment := spec.Assessments[domainName]

		for _, level := range domain.Levels {
			for _, c := range level.Criteria {
				// Determine type (resolve from SLI if needed)
				isQual := c.IsQualitativeWithSpec(spec)
				criterionType := "Quantitative"
				if isQual {
					criterionType = "Qualitative"
				}

				// Determine target display
				var targetDisplay any
				if isQual {
					targetDisplay = "Tracked"
				} else {
					targetDisplay = c.Target
				}

				// Get current value/status
				var currentDisplay any
				var isMet bool
				if isQual {
					status := c.Status
					if assessment != nil && assessment.CriteriaStatus != nil {
						if s, ok := assessment.CriteriaStatus[c.ID]; ok {
							status = s
						}
					}
					isMet = IsQualitativeStatusMet(status)
					currentDisplay = formatQualitativeStatus(status)
				} else {
					var current float64
					if assessment != nil && assessment.CriteriaValues != nil {
						if v, ok := assessment.CriteriaValues[c.ID]; ok {
							current = v
							isMet = c.CheckMet(current)
						}
					}
					currentDisplay = current
				}

				metStatus := "No"
				if isMet {
					metStatus = "Yes"
				}

				required := "Yes"
				if !c.Required && c.Weight > 0 {
					required = "No"
				}

				rows = append(rows, []any{
					c.ID,
					domainName,
					fmt.Sprintf("M%d", level.Level),
					c.Name,
					c.GetMetricName(spec),
					criterionType,
					OperatorSymbol(c.Operator),
					targetDisplay,
					c.GetUnit(spec),
					currentDisplay,
					metStatus,
					c.GetLayer(spec),
					c.GetCategory(spec),
					required,
				})
			}
		}
	}

	frame, err := omniframe.FromRows("SLOs", columns, rows)
	if err != nil {
		return nil, err
	}

	_ = frame.SetColumnWidth("ID", 25)
	_ = frame.SetColumnWidth("Name", 30)
	_ = frame.SetColumnWidth("Metric", 35)
	_ = frame.SetColumnWidth("Type", 12)

	return frame, nil
}

func buildProgressFrame(spec *Spec) (*omniframe.Frame, error) {
	columns := []string{
		"Domain", "Current Level", "Target Level",
		"M2 Progress", "M3 Progress", "M4 Progress", "M5 Progress",
	}

	var rows [][]any
	domainNames := sortedDomainNamesFromSpec(spec)

	for _, domainName := range domainNames {
		domain := spec.Domains[domainName]
		assessment := spec.Assessments[domainName]

		row := []any{
			domain.Name,
			fmt.Sprintf("M%d", assessment.CurrentLevel),
			fmt.Sprintf("M%d", assessment.TargetLevel),
		}

		for level := 2; level <= 5; level++ {
			levelDef, found := domain.GetLevel(level)
			if !found || len(levelDef.Criteria) == 0 {
				row = append(row, "N/A")
				continue
			}

			progress := levelDef.CalculateLevelProgress(assessment.CriteriaValues, assessment.EnablerStatus)
			row = append(row, fmt.Sprintf("%.0f%%", progress.ProgressPercent))
		}

		rows = append(rows, row)
	}

	frame, err := omniframe.FromRows("Progress", columns, rows)
	if err != nil {
		return nil, err
	}

	_ = frame.SetColumnWidth("Domain", 15)

	return frame, nil
}

func sortedDomainNamesFromSpec(spec *Spec) []string {
	var names []string
	for name := range spec.Domains {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
