/**
 * Domain View Component
 *
 * Renders maturity model organized by domain:
 * Domain → Level → Criteria (SLOs) → Enablers
 */

import type { Spec, DomainModel, Level, Criterion, DomainAssessment } from '../schema/maturity';

export interface DomainViewOptions {
  /** Include criteria details (framework mappings) */
  includeDetails?: boolean;
  /** Show enablers section */
  showEnablers?: boolean;
  /** Custom CSS class for the container */
  className?: string;
}

const defaultOptions: DomainViewOptions = {
  includeDetails: true,
  showEnablers: true,
};

/**
 * Render the domain view as HTML string
 */
export function renderDomainView(spec: Spec, options: DomainViewOptions = {}): string {
  const opts = { ...defaultOptions, ...options };
  const domainNames = Object.keys(spec.domains).sort();

  let html = `<div class="prism-domain-view ${opts.className || ''}">`;
  html += '<h2>Maturity Model by Domain</h2>';

  for (const domainName of domainNames) {
    const domain = spec.domains[domainName];
    const assessment = spec.assessments?.[domainName];
    html += renderDomain(domain, assessment, opts);
  }

  html += '</div>';
  return html;
}

function renderDomain(
  domain: DomainModel,
  assessment: DomainAssessment | undefined,
  opts: DomainViewOptions
): string {
  let html = `<section class="prism-domain" data-domain="${escapeHtml(domain.name)}">`;
  html += `<h3>${escapeHtml(domain.name)}</h3>`;

  if (domain.description) {
    html += `<p class="prism-domain-description">${escapeHtml(domain.description)}</p>`;
  }

  if (assessment) {
    html += `
      <div class="prism-assessment-summary">
        <span class="prism-current-level">Current: M${assessment.currentLevel}</span>
        <span class="prism-target-level">Target: M${assessment.targetLevel}</span>
      </div>
    `;
  }

  for (const level of domain.levels) {
    html += renderLevel(level, assessment, opts);
  }

  html += '</section>';
  return html;
}

function renderLevel(
  level: Level,
  assessment: DomainAssessment | undefined,
  opts: DomainViewOptions
): string {
  const isCurrentLevel = assessment?.currentLevel === level.level;
  const isAchieved = assessment ? assessment.currentLevel >= level.level : false;

  let html = `
    <div class="prism-level ${isCurrentLevel ? 'prism-level-current' : ''} ${isAchieved ? 'prism-level-achieved' : ''}"
         data-level="${level.level}">
      <h4>Level ${level.level}: ${escapeHtml(level.name)}</h4>
  `;

  if (level.description) {
    html += `<p class="prism-level-description">${escapeHtml(level.description)}</p>`;
  }

  // Criteria table
  if (level.criteria && level.criteria.length > 0) {
    html += renderCriteriaTable(level.criteria, assessment, opts);
  }

  // Enablers
  if (opts.showEnablers && level.enablers && level.enablers.length > 0) {
    html += renderEnablersTable(level.enablers, assessment);
  }

  html += '</div>';
  return html;
}

function renderCriteriaTable(
  criteria: Criterion[],
  assessment: DomainAssessment | undefined,
  opts: DomainViewOptions
): string {
  let html = `
    <div class="prism-criteria">
      <h5>Criteria (SLOs)</h5>
      <table class="prism-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Type</th>
            <th>Target</th>
            <th>Status</th>
            <th>Frameworks</th>
          </tr>
        </thead>
        <tbody>
  `;

  for (const criterion of criteria) {
    const isMet = checkCriterionMet(criterion, assessment);
    const statusIcon = isMet ? '✅' : '⏳';
    const criterionType = criterion.type === 'qualitative' ? 'Qualitative' : 'Quantitative';
    const target = criterion.type === 'qualitative' ? 'Tracked' : formatTarget(criterion);
    const frameworks = formatFrameworks(criterion.frameworkMappings);

    html += `
      <tr class="${isMet ? 'prism-criterion-met' : 'prism-criterion-pending'}">
        <td><code>${escapeHtml(criterion.id)}</code></td>
        <td>${escapeHtml(criterion.name)}</td>
        <td>${criterionType}</td>
        <td>${target}</td>
        <td>${statusIcon}</td>
        <td>${frameworks}</td>
      </tr>
    `;
  }

  html += '</tbody></table>';

  // Criterion details
  if (opts.includeDetails) {
    for (const criterion of criteria) {
      if (criterion.frameworkMappings && criterion.frameworkMappings.length > 0) {
        html += renderCriterionDetails(criterion);
      }
    }
  }

  html += '</div>';
  return html;
}

function renderCriterionDetails(criterion: Criterion): string {
  if (!criterion.frameworkMappings || criterion.frameworkMappings.length === 0) {
    return '';
  }

  let html = `
    <details class="prism-criterion-details">
      <summary>${escapeHtml(criterion.name)}</summary>
  `;

  if (criterion.description) {
    html += `<p>${escapeHtml(criterion.description)}</p>`;
  }

  html += '<strong>Framework Mappings:</strong><ul>';
  for (const fm of criterion.frameworkMappings) {
    const baseline = fm.baseline ? ` (${fm.baseline})` : '';
    const name = fm.name || fm.reference;
    html += `<li><strong>${formatFrameworkName(fm.framework)}</strong>: ${escapeHtml(name)}${baseline}</li>`;
  }
  html += '</ul></details>';

  return html;
}

function renderEnablersTable(
  enablers: { id: string; name: string; type?: string; status?: string }[],
  assessment: DomainAssessment | undefined
): string {
  let html = `
    <div class="prism-enablers">
      <h5>Enablers</h5>
      <table class="prism-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Type</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
  `;

  for (const enabler of enablers) {
    const status = assessment?.enablerStatus?.[enabler.id] || enabler.status || 'not_started';
    const statusDisplay = formatEnablerStatus(status);

    html += `
      <tr>
        <td><code>${escapeHtml(enabler.id)}</code></td>
        <td>${escapeHtml(enabler.name)}</td>
        <td>${escapeHtml(enabler.type || '-')}</td>
        <td>${statusDisplay}</td>
      </tr>
    `;
  }

  html += '</tbody></table></div>';
  return html;
}

// Helper functions

function checkCriterionMet(criterion: Criterion, assessment: DomainAssessment | undefined): boolean {
  if (!assessment) return false;

  if (criterion.type === 'qualitative' || criterion.operator === 'exists') {
    const status = assessment.criteriaStatus?.[criterion.id];
    return isQualitativeStatusMet(status);
  }

  const value = assessment.criteriaValues?.[criterion.id];
  if (value === undefined) return false;

  const target = criterion.target ?? 0;
  switch (criterion.operator) {
    case 'gte': return value >= target;
    case 'lte': return value <= target;
    case 'gt': return value > target;
    case 'lt': return value < target;
    case 'eq': return value === target;
    default: return false;
  }
}

function isQualitativeStatusMet(status: string | undefined): boolean {
  const metStatuses = ['tracked', 'implemented', 'defined', 'documented', 'compliant', 'enabled'];
  return status !== undefined && metStatuses.includes(status.toLowerCase());
}

function formatTarget(criterion: Criterion): string {
  const opSymbols: Record<string, string> = {
    gte: '>=', lte: '<=', gt: '>', lt: '<', eq: '=', exists: 'Tracked'
  };
  const symbol = opSymbols[criterion.operator] || criterion.operator;
  const unit = criterion.unit || '';
  return `${symbol}${criterion.target ?? ''}${unit}`;
}

function formatFrameworks(mappings?: { framework: string; reference: string }[]): string {
  if (!mappings || mappings.length === 0) return '-';
  return mappings.map(m => `${m.framework}:${m.reference}`).join(', ');
}

function formatFrameworkName(fw: string): string {
  const names: Record<string, string> = {
    'NIST_CSF': 'NIST CSF 1.1',
    'NIST_CSF_2': 'NIST CSF 2.0',
    'NIST_800_53': 'NIST SP 800-53',
    'FEDRAMP_HIGH': 'FedRAMP High',
    'FEDRAMP_MOD': 'FedRAMP Moderate',
    'FEDRAMP_LOW': 'FedRAMP Low',
    'DORA': 'DORA',
    'SRE': 'SRE',
  };
  return names[fw] || fw;
}

function formatEnablerStatus(status: string): string {
  const icons: Record<string, string> = {
    'completed': '✅ Completed',
    'in_progress': '🔄 In Progress',
    'blocked': '🚫 Blocked',
    'not_started': '⏳ Not Started',
  };
  return icons[status] || status;
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}
