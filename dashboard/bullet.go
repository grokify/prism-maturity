package dashboard

import (
	"encoding/json"
	"fmt"
)

// MaturityBullet represents a D3 bullet chart for maturity visualization.
// Ranges are configured for M1-M3 (red), M4 (yellow), M5 (green) zones.
type MaturityBullet struct {
	Title       string           `json:"title,omitempty"`
	Subtitle    string           `json:"subtitle,omitempty"`
	Ranges      []float64        `json:"ranges,omitempty"`      // Zone boundaries [M3, M4, M5]
	Measures    []float64        `json:"measures,omitempty"`    // Current maturity level(s)
	Markers     []float64        `json:"markers,omitempty"`     // Target level(s)
	ActualValue string           `json:"actualValue,omitempty"` // Actual SLI value with unit (e.g., "65%", "120ms")
	Thresholds  []LevelThreshold `json:"thresholds,omitempty"`  // Thresholds for each maturity level
}

// LevelThreshold defines the threshold for a specific maturity level.
type LevelThreshold struct {
	Level    int     `json:"level"`              // Maturity level (1-5)
	Operator string  `json:"operator,omitempty"` // >=, <=, ==, etc.
	Value    float64 `json:"value"`              // Threshold value
	ValueStr string  `json:"valueStr,omitempty"` // Formatted value with unit
}

// MaturityBulletData holds a collection of maturity bullet charts.
type MaturityBulletData struct {
	Bullets []MaturityBullet `json:"bullets"`
}

// ToJSON returns the bullet data as JSON for JavaScript consumption.
func (d *MaturityBulletData) ToJSON() ([]byte, error) {
	return json.Marshal(d.Bullets)
}

// NewMaturityBullet creates a bullet chart for a maturity metric.
// currentLevel is the actual achieved level (0-5).
// targetLevel is the target level (0-5), use 0 for no target marker.
func NewMaturityBullet(title, subtitle string, currentLevel, targetLevel float64) MaturityBullet {
	bullet := MaturityBullet{
		Title:    title,
		Subtitle: subtitle,
		Ranges:   []float64{3, 4, 5}, // M1-3 | M4 | M5 zones
		Measures: []float64{currentLevel},
	}
	if targetLevel > 0 {
		bullet.Markers = []float64{targetLevel}
	}
	return bullet
}

// NewMaturityBulletWithDetails creates a bullet chart with actual SLI value and thresholds.
func NewMaturityBulletWithDetails(title string, currentLevel, targetLevel, actualValue float64, unit string, thresholds []LevelThreshold) MaturityBullet {
	// Format actual value with unit
	actualValueStr := ""
	if actualValue != 0 || unit != "" {
		if actualValue == float64(int(actualValue)) {
			actualValueStr = fmt.Sprintf("%d%s", int(actualValue), unit)
		} else {
			actualValueStr = fmt.Sprintf("%.1f%s", actualValue, unit)
		}
	}

	// Build subtitle with level and actual value
	subtitle := MaturityLevel(currentLevel)
	if actualValueStr != "" {
		subtitle = fmt.Sprintf("%s (%s)", actualValueStr, MaturityLevel(currentLevel))
	}

	bullet := MaturityBullet{
		Title:       title,
		Subtitle:    subtitle,
		Ranges:      []float64{3, 4, 5},
		Measures:    []float64{currentLevel},
		ActualValue: actualValueStr,
		Thresholds:  thresholds,
	}
	if targetLevel > 0 && targetLevel != currentLevel {
		bullet.Markers = []float64{targetLevel}
	}
	return bullet
}

// NewMaturityBulletWithProjection creates a bullet chart with current and projected levels.
func NewMaturityBulletWithProjection(title, subtitle string, currentLevel, projectedLevel, targetLevel float64) MaturityBullet {
	bullet := MaturityBullet{
		Title:    title,
		Subtitle: subtitle,
		Ranges:   []float64{3, 4, 5},
		Measures: []float64{currentLevel, projectedLevel},
	}
	if targetLevel > 0 {
		bullet.Markers = []float64{targetLevel}
	}
	return bullet
}

// GetMaturityBulletCSS returns CSS for maturity-colored bullet charts.
// D3 sorts ranges descending: s0=highest(M5), s1=middle(M4), s2=lowest(M1-3).
// Zone colors: M5 = green (right), M4 = yellow (middle), M1-M3 = red (left).
func GetMaturityBulletCSS() string {
	return `.bullet { font: 10px sans-serif; }
.bullet .marker { stroke: #000; stroke-width: 2px; }
.bullet .tick line { stroke: #666; stroke-width: .5px; }
.bullet .range.s0 { fill: #dcfce7; }
.bullet .range.s1 { fill: #fef3c7; }
.bullet .range.s2 { fill: #fee2e2; }
.bullet .measure.s0 { fill: #3b82f6; }
.bullet .measure.s1 { fill: #60a5fa; }
.bullet .title { font-size: 14px; font-weight: bold; }
.bullet .subtitle { fill: #999; }`
}

// GetMaturityBulletCSSStyled returns CSS with the style tags included.
func GetMaturityBulletCSSStyled() string {
	return "<style>\n" + GetMaturityBulletCSS() + "\n</style>"
}

// MaturityLevel returns the maturity level label (M1-M5) for a numeric value.
func MaturityLevel(value float64) string {
	switch {
	case value >= 5:
		return "M5"
	case value >= 4:
		return "M4"
	case value >= 3:
		return "M3"
	case value >= 2:
		return "M2"
	case value >= 1:
		return "M1"
	default:
		return "M0"
	}
}

// MaturityStatus returns the status (green/yellow/red) for a maturity level.
func MaturityStatus(level float64) string {
	switch {
	case level >= 5:
		return "green"
	case level >= 4:
		return "yellow"
	default:
		return "red"
	}
}

// MaturityStatusEmoji returns an emoji indicator for a maturity level.
func MaturityStatusEmoji(level float64) string {
	switch {
	case level >= 5:
		return "🟢"
	case level >= 4:
		return "🟡"
	default:
		return "🔴"
	}
}
