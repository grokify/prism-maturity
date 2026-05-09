/**
 * Framework View Component
 *
 * Renders maturity model organized by compliance framework:
 * Framework → Control → Criteria mappings
 */

import type { Spec, Criterion, DomainAssessment, FrameworkMapping } from '../schema/maturity';

export interface FrameworkViewOptions {
  /** Filter to specific frameworks (empty = all) */
  frameworks?: string[];
  /** Custom CSS class for the container */
  className?: string;
}

interface CriterionRef {
  domain: string;
  level: number;
  criterion: Criterion;
  mapping: FrameworkMapping;
}

/**
 * Render the framework view as HTML string
 */
export function renderFrameworkView(spec: Spec, options: FrameworkViewOptions = {}): string {
  const frameworks = collectFrameworks(spec, options.frameworks);

  if (frameworks.length === 0) {
    return '<div class="prism-framework-view"><p>No framework mappings found.</p></div>';
  }

  let html = `<div class="prism-framework-view ${options.className || ''}">`;
  html += '<h2>Maturity Model by Framework</h2>';
  html += '<p>This view shows how maturity criteria map to compliance framework controls.</p>';

  for (const framework of frameworks) {
    html += renderFramework(spec, framework);
  }

  html += '</div>';
  return html;
}

function collectFrameworks(spec: Spec, filter?: string[]): string[] {
  const frameworkSet = new Set<string>();

  for (const domain of Object.values(spec.domains)) {
    for (const level of domain.levels) {
      for (const criterion of level.criteria || []) {
        for (const fm of criterion.frameworkMappings || []) {
          if (!filter || filter.length === 0 || filter.includes(fm.framework)) {
            frameworkSet.add(fm.framework);
          }
        }
      }
    }
  }

  return Array.from(frameworkSet).sort();
}

function renderFramework(spec: Spec, framework: string): string {
  const refs = collectCriteriaForFramework(spec, framework);

  let html = `
    <section class="prism-framework" data-framework="${escapeHtml(framework)}">
      <h3>${formatFrameworkName(framework)}</h3>
  `;

  if (refs.length === 0) {
    html += '<p>No criteria mapped to this framework.</p>';
    html += '</section>';
    return html;
  }

  // Sort by control reference
  refs.sort((a, b) => a.mapping.reference.localeCompare(b.mapping.reference));

  // Render table
  html += `
    <table class="prism-table">
      <thead>
        <tr>
          <th>Control</th>
          <th>Name</th>
          <th>Domain</th>
          <th>Level</th>
          <th>Criterion</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
  `;

  let metCount = 0;
  for (const ref of refs) {
    const assessment = spec.assessments?.[ref.domain];
    const isMet = checkCriterionMet(ref.criterion, assessment);
    if (isMet) metCount++;

    const statusIcon = isMet ? '✅' : '⏳';
    const controlName = ref.mapping.name || '-';
    const baseline = ref.mapping.baseline ? ` [${ref.mapping.baseline}]` : '';

    html += `
      <tr class="${isMet ? 'prism-criterion-met' : 'prism-criterion-pending'}">
        <td><code>${escapeHtml(ref.mapping.reference)}${baseline}</code></td>
        <td>${escapeHtml(controlName)}</td>
        <td>${escapeHtml(ref.domain)}</td>
        <td>M${ref.level}</td>
        <td>${escapeHtml(ref.criterion.name)}</td>
        <td>${statusIcon}</td>
      </tr>
    `;
  }

  html += '</tbody></table>';

  // Coverage summary
  const percentage = refs.length > 0 ? Math.round((metCount / refs.length) * 100) : 0;
  html += `
    <div class="prism-framework-coverage">
      <strong>Coverage:</strong> ${metCount}/${refs.length} controls satisfied (${percentage}%)
    </div>
  `;

  html += '</section>';
  return html;
}

function collectCriteriaForFramework(spec: Spec, framework: string): CriterionRef[] {
  const refs: CriterionRef[] = [];

  for (const [domainName, domain] of Object.entries(spec.domains)) {
    for (const level of domain.levels) {
      for (const criterion of level.criteria || []) {
        for (const fm of criterion.frameworkMappings || []) {
          if (fm.framework === framework) {
            refs.push({
              domain: domainName,
              level: level.level,
              criterion,
              mapping: fm,
            });
          }
        }
      }
    }
  }

  return refs;
}

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

function formatFrameworkName(fw: string): string {
  const names: Record<string, string> = {
    'NIST_CSF': 'NIST CSF 1.1',
    'NIST_CSF_2': 'NIST CSF 2.0',
    'NIST_800_53': 'NIST SP 800-53',
    'NIST_RMF': 'NIST RMF',
    'NIST_AI_RMF': 'NIST AI RMF',
    'NIST_800_171': 'NIST SP 800-171',
    'FEDRAMP': 'FedRAMP',
    'FEDRAMP_HIGH': 'FedRAMP High',
    'FEDRAMP_MOD': 'FedRAMP Moderate',
    'FEDRAMP_LOW': 'FedRAMP Low',
    'MITRE_ATTACK': 'MITRE ATT&CK',
    'CIS_CONTROLS': 'CIS Controls',
    'SOC_2': 'SOC 2',
    'ISO_27001': 'ISO 27001',
    'DORA': 'DORA',
    'SRE': 'SRE',
  };
  return names[fw] || fw;
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}
