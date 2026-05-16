package dashboard

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
)

// HTMLOptions configures HTML output.
type HTMLOptions struct {
	EmbedData bool   // Embed data in HTML instead of external JSON
	Title     string // Page title override
	Theme     string // "light" or "dark"
}

// DefaultHTMLOptions returns default HTML options.
func DefaultHTMLOptions() HTMLOptions {
	return HTMLOptions{
		EmbedData: true,
		Theme:     "light",
	}
}

// WriteHTML writes the dashboard as a standalone HTML file.
func (d *Dashboard) WriteHTML(w io.Writer, opts HTMLOptions) error {
	tmpl, err := template.New("dashboard").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	dashboardJSON, err := json.Marshal(d)
	if err != nil {
		return err
	}

	title := d.Title
	if opts.Title != "" {
		title = opts.Title
	}

	data := map[string]any{
		"Title":         title,
		"DashboardJSON": template.JS(dashboardJSON), //nolint:gosec // Safe: marshaling our own data
		"Theme":         opts.Theme,
		"EmbedData":     opts.EmbedData,
	}

	return tmpl.Execute(w, data)
}

// ToHTML returns the dashboard as an HTML string.
func (d *Dashboard) ToHTML(opts HTMLOptions) (string, error) {
	var buf bytes.Buffer
	if err := d.WriteHTML(&buf, opts); err != nil {
		return "", err
	}
	return buf.String(), nil
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Title}}</title>
  <script src="https://cdn.jsdelivr.net/npm/echarts@5.5.0/dist/echarts.min.js"></script>
  <script src="https://d3js.org/d3.v3.min.js"></script>
  <style>
    :root {
      --bg-primary: {{if eq .Theme "dark"}}#1a1a2e{{else}}#f8fafc{{end}};
      --bg-card: {{if eq .Theme "dark"}}#16213e{{else}}#ffffff{{end}};
      --text-primary: {{if eq .Theme "dark"}}#e2e8f0{{else}}#1e293b{{end}};
      --text-secondary: {{if eq .Theme "dark"}}#94a3b8{{else}}#64748b{{end}};
      --border-color: {{if eq .Theme "dark"}}#334155{{else}}#e2e8f0{{end}};
      --green: #22c55e;
      --yellow: #f59e0b;
      --red: #ef4444;
      --blue: #3b82f6;
    }

    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }

    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
      background-color: var(--bg-primary);
      color: var(--text-primary);
      line-height: 1.5;
    }

    .dashboard {
      display: grid;
      grid-template-columns: repeat(12, 1fr);
      gap: 16px;
      padding: 24px;
      max-width: 1600px;
      margin: 0 auto;
    }

    .widget {
      background: var(--bg-card);
      border-radius: 8px;
      box-shadow: 0 1px 3px rgba(0,0,0,0.1);
      padding: 16px;
      border: 1px solid var(--border-color);
    }

    .widget-title {
      font-size: 14px;
      font-weight: 600;
      color: var(--text-secondary);
      margin-bottom: 12px;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    .metric-value {
      font-size: 48px;
      font-weight: 700;
      line-height: 1;
    }

    .metric-subtitle {
      font-size: 14px;
      color: var(--text-secondary);
      margin-top: 8px;
    }

    .chart-container {
      width: 100%;
      height: 100%;
      min-height: 280px;
    }

    .table-container {
      overflow-x: auto;
    }

    table {
      width: 100%;
      border-collapse: collapse;
      font-size: 14px;
    }

    th {
      text-align: left;
      padding: 12px 8px;
      border-bottom: 2px solid var(--border-color);
      font-weight: 600;
      color: var(--text-secondary);
    }

    td {
      padding: 10px 8px;
      border-bottom: 1px solid var(--border-color);
    }

    tr.striped:nth-child(even) {
      background: {{if eq .Theme "dark"}}rgba(255,255,255,0.02){{else}}rgba(0,0,0,0.02){{end}};
    }

    .status-badge {
      display: inline-block;
      padding: 2px 8px;
      border-radius: 4px;
      font-size: 12px;
      font-weight: 500;
    }

    .status-green { background: rgba(34, 197, 94, 0.1); color: var(--green); }
    .status-yellow { background: rgba(245, 158, 11, 0.1); color: var(--yellow); }
    .status-red { background: rgba(239, 68, 68, 0.1); color: var(--red); }

    .level-indicator {
      display: flex;
      align-items: center;
      gap: 4px;
      margin-top: 8px;
    }

    .level-dot {
      width: 12px;
      height: 12px;
      border-radius: 50%;
      background: var(--border-color);
    }

    .level-dot.active { background: var(--blue); }
    .level-dot.completed { background: var(--green); }

    .text-content {
      line-height: 1.6;
    }

    .text-content h1 {
      font-size: 24px;
      margin-bottom: 8px;
    }

    .text-content p {
      color: var(--text-secondary);
    }

    @media (max-width: 768px) {
      .dashboard {
        grid-template-columns: 1fr;
        padding: 16px;
      }
      .widget {
        grid-column: span 1 !important;
      }
    }

    /* Maturity Bullet Chart Styles */
    /* D3 sorts ranges descending: s0=highest(M5), s1=middle(M4), s2=lowest(M1-3) */
    .bullet { font: 12px sans-serif; }
    .bullet .marker { stroke: #000; stroke-width: 2px; }
    .bullet .tick line { stroke: #666; stroke-width: .5px; }
    .bullet .range.s0 { fill: #dcfce7; }
    .bullet .range.s1 { fill: #fef3c7; }
    .bullet .range.s2 { fill: #fee2e2; }
    .bullet .measure.s0 { fill: #3b82f6; }
    .bullet .measure.s1 { fill: #60a5fa; }
    .bullet-container { padding: 8px 0; }
    .bullet-row {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
      gap: 12px;
    }
    .bullet-label {
      flex: 0 0 auto;
      text-align: left;
    }
    .bullet-label-title {
      font-size: 14px;
      font-weight: 600;
      color: var(--text-primary);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .bullet-label-subtitle {
      font-size: 12px;
      color: var(--text-secondary);
    }
    .bullet-chart {
      flex: 1 1 auto;
      min-width: 200px;
    }
    .bullet-legend {
      display: flex;
      gap: 16px;
      font-size: 11px;
      margin-bottom: 8px;
      color: var(--text-secondary);
    }
    .bullet-legend-item {
      display: flex;
      align-items: center;
      gap: 4px;
    }
    .bullet-legend-swatch {
      width: 12px;
      height: 12px;
      border-radius: 2px;
    }
  </style>
</head>
<body>
  <div id="dashboard" class="dashboard"></div>

  <script>
    const dashboard = {{.DashboardJSON}};

    function render() {
      const container = document.getElementById('dashboard');

      // Create data lookup
      const dataMap = {};
      for (const ds of dashboard.dataSources || []) {
        if (ds.data) {
          dataMap[ds.id] = ds.data;
        }
      }

      // Render widgets
      for (const widget of dashboard.widgets || []) {
        const el = document.createElement('div');
        el.className = 'widget';
        el.style.gridColumn = ` + "`" + `span ${widget.position.w}` + "`" + `;
        el.style.gridRow = ` + "`" + `span ${widget.position.h}` + "`" + `;

        const data = dataMap[widget.dataSourceId] || [];

        switch (widget.type) {
          case 'text':
            renderText(el, widget);
            break;
          case 'metric':
            renderMetric(el, widget, data);
            break;
          case 'chart':
            renderChart(el, widget, data);
            break;
          case 'table':
            renderTable(el, widget, data);
            break;
          case 'bullet':
            renderBullet(el, widget, data);
            break;
        }

        container.appendChild(el);
      }
    }

    function renderText(el, widget) {
      const config = widget.config || {};
      el.innerHTML = ` + "`" + `<div class="text-content">${parseMarkdown(config.content || '')}</div>` + "`" + `;
    }

    function parseMarkdown(text) {
      return text
        .replace(/^# (.+)$/gm, '<h1>$1</h1>')
        .replace(/^## (.+)$/gm, '<h2>$1</h2>')
        .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
        .replace(/\n\n/g, '</p><p>')
        .replace(/^(.+)$/gm, (m, p1) => p1.startsWith('<') ? p1 : ` + "`" + `<p>${p1}</p>` + "`" + `);
    }

    function renderMetric(el, widget, data) {
      const config = widget.config || {};
      const row = Array.isArray(data) && data[0] ? data[0] : data;
      const value = row[config.valueField] || 0;

      const opts = config.formatOptions || {};
      const formatted = (opts.prefix || '') + value + (opts.suffix || '');

      // Find threshold color
      let color = 'var(--text-primary)';
      if (config.thresholds) {
        for (const t of config.thresholds.reverse()) {
          if (value >= t.value) {
            color = t.color;
            break;
          }
        }
      }

      el.innerHTML = ` + "`" + `
        <div class="widget-title">${widget.title}</div>
        <div class="metric-value" style="color: ${color}">${formatted}</div>
        ${config.subtitle ? ` + "`" + `<div class="metric-subtitle">${config.subtitle}</div>` + "`" + ` : ''}
        <div class="level-indicator">
          ${[1,2,3,4,5].map(l => ` + "`" + `
            <div class="level-dot ${l <= value ? 'completed' : ''}" title="M${l}"></div>
          ` + "`" + `).join('')}
        </div>
      ` + "`" + `;
    }

    function renderChart(el, widget, data) {
      // Calculate chart height based on grid position (rowHeight=80, gap=16, minus title/padding)
      const rowHeight = dashboard.layout?.rowHeight || 80;
      const gap = dashboard.layout?.gap || 16;
      const chartHeight = (widget.position.h * rowHeight) + ((widget.position.h - 1) * gap) - 60;

      el.innerHTML = ` + "`" + `
        <div class="widget-title">${widget.title}</div>
        <div class="chart-container" id="chart-${widget.id}" style="height: ${chartHeight}px;"></div>
      ` + "`" + `;

      setTimeout(() => {
        const chartEl = document.getElementById('chart-' + widget.id);
        if (!chartEl) return;

        const chart = echarts.init(chartEl);
        const config = widget.config || {};
        const option = compileChartIR(config, data);
        chart.setOption(option);

        window.addEventListener('resize', () => chart.resize());
      }, 0);
    }

    function compileChartIR(ir, data) {
      const option = {
        dataset: { source: data },
        tooltip: {
          trigger: ir.tooltip?.trigger || 'axis',
          show: ir.tooltip?.show !== false
        },
        legend: ir.legend?.show ? {
          top: ir.legend.position === 'bottom' ? 'bottom' : 'top'
        } : undefined,
        grid: ir.grid || { left: '3%', right: '4%', bottom: '3%', containLabel: true },
        xAxis: buildAxis(ir.axes?.find(a => a.position === 'bottom' || a.position === 'top'), 'x'),
        yAxis: buildAxis(ir.axes?.find(a => a.position === 'left' || a.position === 'right'), 'y'),
        series: (ir.marks || []).map(mark => ({
          type: mark.geometry === 'area' ? 'line' : mark.geometry,
          name: mark.name || mark.id,
          encode: mark.encode,
          smooth: mark.smooth,
          stack: mark.stack,
          areaStyle: mark.geometry === 'area' ? {} : undefined,
          itemStyle: mark.style?.color ? { color: mark.style.color } : undefined,
          barWidth: mark.style?.barWidth,
          emphasis: { focus: 'series' }
        }))
      };
      return option;
    }

    function buildAxis(axis, dim) {
      if (!axis) {
        return dim === 'x' ? { type: 'category' } : { type: 'value' };
      }
      return {
        type: axis.type || (dim === 'x' ? 'category' : 'value'),
        name: axis.name,
        min: axis.min,
        max: axis.max,
        axisLabel: { interval: 0 }
      };
    }

    function renderTable(el, widget, data) {
      const config = widget.config || {};
      const columns = config.columns || [];

      el.innerHTML = ` + "`" + `
        <div class="widget-title">${widget.title}</div>
        <div class="table-container">
          <table>
            <thead>
              <tr>
                ${columns.map(c => ` + "`" + `<th style="width:${c.width || 'auto'};text-align:${c.align || 'left'}">${c.header}</th>` + "`" + `).join('')}
              </tr>
            </thead>
            <tbody>
              ${(Array.isArray(data) ? data : []).map((row, i) => ` + "`" + `
                <tr class="${config.striped ? 'striped' : ''}">
                  ${columns.map(c => ` + "`" + `<td style="text-align:${c.align || 'left'}">${formatCell(row[c.field], c)}</td>` + "`" + `).join('')}
                </tr>
              ` + "`" + `).join('')}
            </tbody>
          </table>
        </div>
      ` + "`" + `;
    }

    function formatCell(value, column) {
      if (value === undefined || value === null) return '-';
      if (column.format === 'number') {
        return Number(value).toLocaleString();
      }
      return value;
    }

    function renderBullet(el, widget, data) {
      el.innerHTML = ` + "`" + `
        <div class="widget-title">${widget.title}</div>
        <div class="bullet-legend">
          <div class="bullet-legend-item">
            <div class="bullet-legend-swatch" style="background: #fee2e2"></div>
            <span>M1-M3</span>
          </div>
          <div class="bullet-legend-item">
            <div class="bullet-legend-swatch" style="background: #fef3c7"></div>
            <span>M4</span>
          </div>
          <div class="bullet-legend-item">
            <div class="bullet-legend-swatch" style="background: #dcfce7"></div>
            <span>M5</span>
          </div>
          <div class="bullet-legend-item">
            <div class="bullet-legend-swatch" style="background: #3b82f6"></div>
            <span>Current</span>
          </div>
          <div class="bullet-legend-item">
            <span style="border-left: 2px solid #000; height: 12px; margin-right: 4px;"></span>
            <span>Target</span>
          </div>
        </div>
        <div class="bullet-container" id="bullet-${widget.id}"></div>
      ` + "`" + `;

      setTimeout(() => {
        const container = document.getElementById('bullet-' + widget.id);
        if (!container || !Array.isArray(data)) return;

        // Create bullet rows with HTML labels and SVG charts
        data.forEach((d, i) => {
          const row = document.createElement('div');
          row.className = 'bullet-row';
          row.innerHTML = ` + "`" + `
            <div class="bullet-label">
              <div class="bullet-label-title" title="${d.title}">${d.title}</div>
              <div class="bullet-label-subtitle">${d.subtitle || ''}</div>
            </div>
            <div class="bullet-chart" id="bullet-chart-${widget.id}-${i}"></div>
          ` + "`" + `;
          container.appendChild(row);
        });

        // Calculate max label width and apply uniform width
        setTimeout(() => {
          const labels = container.querySelectorAll('.bullet-label');
          let maxWidth = 0;
          labels.forEach(label => {
            maxWidth = Math.max(maxWidth, label.offsetWidth);
          });
          // Add padding and apply uniform width
          const uniformWidth = Math.ceil(maxWidth) + 16;
          labels.forEach(label => {
            label.style.flex = '0 0 ' + uniformWidth + 'px';
            label.style.width = uniformWidth + 'px';
          });

          // Render SVG bullet charts after layout is complete
          setTimeout(() => {
            data.forEach((d, i) => {
              const chartContainer = document.getElementById(` + "`" + `bullet-chart-${widget.id}-${i}` + "`" + `);
              if (!chartContainer) return;

              const margin = {top: 5, right: 40, bottom: 35, left: 0};
              const width = chartContainer.offsetWidth - margin.left - margin.right;
              const height = 32;

              const chart = d3MaturityBullet().width(width).height(height);

              d3.select(chartContainer).append('svg')
                  .attr('class', 'bullet')
                  .attr('width', width + margin.left + margin.right)
                  .attr('height', height + margin.top + margin.bottom)
                .append('g')
                  .attr('transform', 'translate(' + margin.left + ',' + margin.top + ')')
                  .datum(d)
                  .call(chart);
            });
          }, 0);
        }, 0);
      }, 0);
    }

    // D3 Maturity Bullet Chart - specialized for M1-M5 levels
    function d3MaturityBullet() {
      let width = 380;
      let height = 30;

      function bullet(g) {
        g.each(function(d) {
          const rangez = (d.ranges || [3, 4, 5]).slice().sort(d3.descending);
          const markerz = (d.markers || []).slice();
          const measurez = (d.measures || []).slice();
          const g = d3.select(this);

          // Fixed scale 0-5 for maturity levels
          const x = d3.scale.linear()
              .domain([0, 5])
              .range([0, width]);

          // Background ranges (colored zones)
          const rangeData = [
            { start: 0, end: 3, class: 's2' },   // M1-M3: red
            { start: 3, end: 4, class: 's1' },   // M4: yellow
            { start: 4, end: 5, class: 's0' }    // M5: green
          ];

          g.selectAll('rect.range')
              .data(rangeData)
            .enter().append('rect')
              .attr('class', d => 'range ' + d.class)
              .attr('x', d => x(d.start))
              .attr('width', d => x(d.end) - x(d.start))
              .attr('height', height);

          // Measure bar (current level)
          if (measurez.length > 0) {
            const currentLevel = measurez[0];
            g.append('rect')
                .attr('class', 'measure s0')
                .attr('x', 0)
                .attr('width', x(currentLevel))
                .attr('height', height / 3)
                .attr('y', height / 3);

            // Current level label on bar
            g.append('text')
                .attr('class', 'measure-label')
                .attr('x', x(currentLevel) + 4)
                .attr('y', height / 2 + 5)
                .style('font-size', '13px')
                .style('font-weight', 'bold')
                .style('fill', '#3b82f6')
                .text(currentLevel.toFixed(1));
          }

          // Target marker
          markerz.forEach(m => {
            g.append('line')
                .attr('class', 'marker')
                .attr('x1', x(m))
                .attr('x2', x(m))
                .attr('y1', height / 6)
                .attr('y2', height * 5 / 6);
          });

          // M1-M5 tick marks and labels with thresholds
          const ticks = [1, 2, 3, 4, 5];
          const thresholds = d.thresholds || [];

          // Create threshold lookup by level
          const thresholdByLevel = {};
          thresholds.forEach(t => {
            thresholdByLevel[t.level] = t.valueStr || '';
          });

          const tickG = g.selectAll('g.tick')
              .data(ticks)
            .enter().append('g')
              .attr('class', 'tick')
              .attr('transform', d => 'translate(' + x(d) + ',0)');

          tickG.append('line')
              .attr('y1', height)
              .attr('y2', height + 6)
              .style('stroke', '#666')
              .style('stroke-width', '1px');

          // Level label (M1, M2, etc.)
          tickG.append('text')
              .attr('text-anchor', 'middle')
              .attr('y', height + 16)
              .style('font-size', '11px')
              .style('font-weight', 'bold')
              .style('fill', '#333')
              .text(d => 'M' + d);

          // Threshold value below level label
          // Show "-" for levels without thresholds (only if SLI has any thresholds)
          const hasThresholds = thresholds.length > 0;
          tickG.append('text')
              .attr('text-anchor', 'middle')
              .attr('y', height + 28)
              .style('font-size', '10px')
              .style('fill', '#666')
              .text(d => thresholdByLevel[d] || (hasThresholds ? '-' : ''));
        });
      }

      bullet.width = function(x) {
        if (!arguments.length) return width;
        width = x;
        return bullet;
      };

      bullet.height = function(x) {
        if (!arguments.length) return height;
        height = x;
        return bullet;
      };

      return bullet;
    }

    // D3 Bullet Chart implementation (generic)
    function d3Bullet() {
      let orient = 'left';
      let reverse = false;
      let duration = 0;
      let ranges = d => d.ranges;
      let markers = d => d.markers;
      let measures = d => d.measures;
      let width = 380;
      let height = 30;
      let tickFormat = null;

      function bullet(g) {
        g.each(function(d, i) {
          const rangez = ranges.call(this, d, i).slice().sort(d3.descending);
          const markerz = markers.call(this, d, i).slice().sort(d3.descending);
          const measurez = measures.call(this, d, i).slice().sort(d3.descending);
          const g = d3.select(this);

          const x1 = d3.scale.linear()
              .domain([0, Math.max(rangez[0], markerz[0] || 0, measurez[0] || 0)])
              .range(reverse ? [width, 0] : [0, width]);

          const x0 = this.__chart__ || d3.scale.linear()
              .domain([0, Infinity])
              .range(x1.range());

          this.__chart__ = x1;

          const w0 = bulletWidth(x0);
          const w1 = bulletWidth(x1);

          // Ranges
          let range = g.selectAll('rect.range')
              .data(rangez);

          range.enter().append('rect')
              .attr('class', (d, i) => 'range s' + i)
              .attr('width', w0)
              .attr('height', height)
              .attr('x', reverse ? x0 : 0)
            .transition()
              .duration(duration)
              .attr('width', w1)
              .attr('x', reverse ? x1 : 0);

          range.transition()
              .duration(duration)
              .attr('x', reverse ? x1 : 0)
              .attr('width', w1)
              .attr('height', height);

          // Measures
          let measure = g.selectAll('rect.measure')
              .data(measurez);

          measure.enter().append('rect')
              .attr('class', (d, i) => 'measure s' + i)
              .attr('width', w0)
              .attr('height', height / 3)
              .attr('x', reverse ? x0 : 0)
              .attr('y', height / 3)
            .transition()
              .duration(duration)
              .attr('width', w1)
              .attr('x', reverse ? x1 : 0);

          measure.transition()
              .duration(duration)
              .attr('width', w1)
              .attr('height', height / 3)
              .attr('x', reverse ? x1 : 0)
              .attr('y', height / 3);

          // Markers
          let marker = g.selectAll('line.marker')
              .data(markerz);

          marker.enter().append('line')
              .attr('class', 'marker')
              .attr('x1', x0)
              .attr('x2', x0)
              .attr('y1', height / 6)
              .attr('y2', height * 5 / 6)
            .transition()
              .duration(duration)
              .attr('x1', x1)
              .attr('x2', x1);

          marker.transition()
              .duration(duration)
              .attr('x1', x1)
              .attr('x2', x1)
              .attr('y1', height / 6)
              .attr('y2', height * 5 / 6);

          // Ticks
          const format = tickFormat || x1.tickFormat(8);
          let tick = g.selectAll('g.tick')
              .data(x1.ticks(8), d => this.textContent || format(d));

          const tickEnter = tick.enter().append('g')
              .attr('class', 'tick')
              .attr('transform', bulletTranslate(x0))
              .style('opacity', 1e-6);

          tickEnter.append('line')
              .attr('y1', height)
              .attr('y2', height * 7 / 6);

          tickEnter.append('text')
              .attr('text-anchor', 'middle')
              .attr('dy', '1em')
              .attr('y', height * 7 / 6)
              .text(format);

          tickEnter.transition()
              .duration(duration)
              .attr('transform', bulletTranslate(x1))
              .style('opacity', 1);

          const tickUpdate = tick.transition()
              .duration(duration)
              .attr('transform', bulletTranslate(x1))
              .style('opacity', 1);

          tickUpdate.select('line')
              .attr('y1', height)
              .attr('y2', height * 7 / 6);

          tickUpdate.select('text')
              .attr('y', height * 7 / 6);

          tick.exit().transition()
              .duration(duration)
              .attr('transform', bulletTranslate(x1))
              .style('opacity', 1e-6)
              .remove();
        });
        d3.timer.flush();
      }

      function bulletWidth(x) {
        const x0 = x(0);
        return function(d) {
          return Math.abs(x(d) - x0);
        };
      }

      function bulletTranslate(x) {
        return function(d) {
          return 'translate(' + x(d) + ',0)';
        };
      }

      bullet.orient = function(x) {
        if (!arguments.length) return orient;
        orient = x;
        reverse = orient === 'right' || orient === 'bottom';
        return bullet;
      };

      bullet.ranges = function(x) {
        if (!arguments.length) return ranges;
        ranges = x;
        return bullet;
      };

      bullet.markers = function(x) {
        if (!arguments.length) return markers;
        markers = x;
        return bullet;
      };

      bullet.measures = function(x) {
        if (!arguments.length) return measures;
        measures = x;
        return bullet;
      };

      bullet.width = function(x) {
        if (!arguments.length) return width;
        width = x;
        return bullet;
      };

      bullet.height = function(x) {
        if (!arguments.length) return height;
        height = x;
        return bullet;
      };

      bullet.tickFormat = function(x) {
        if (!arguments.length) return tickFormat;
        tickFormat = x;
        return bullet;
      };

      bullet.duration = function(x) {
        if (!arguments.length) return duration;
        duration = x;
        return bullet;
      };

      return bullet;
    }

    // Initialize
    render();
  </script>
</body>
</html>`
