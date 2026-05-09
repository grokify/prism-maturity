/**
 * Maturity Spec - Zod Schema
 *
 * Hand-written Zod schema matching Go types in maturity/types.go.
 * The auto-generated version doesn't handle $ref properly.
 *
 * Pipeline: Go types → JSON Schema → (manual) Zod Schema → TypeScript types
 *
 * To update: Sync with Go types in maturity/types.go
 */

import { z } from 'zod';

// Framework mapping for compliance frameworks
export const FrameworkMappingSchema = z.object({
  framework: z.string(),
  reference: z.string(),
  name: z.string().optional(),
  description: z.string().optional(),
  baseline: z.string().optional(),
  version: z.string().optional(),
});

export type FrameworkMapping = z.infer<typeof FrameworkMappingSchema>;

// SLI (Service Level Indicator) - the metric being measured
// Framework mappings are defined here since they apply to the metric itself
export const SLISchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().optional(),
  metricName: z.string(),
  unit: z.string().optional(),
  type: z.enum(['quantitative', 'qualitative']).optional(),
  layer: z.string().optional(),
  category: z.string().optional(),
  frameworkMappings: z.array(FrameworkMappingSchema).optional(),
});

export type SLI = z.infer<typeof SLISchema>;

// Criterion (SLO) within a maturity level
// References an SLI and specifies a target threshold for that level
export const CriterionSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().optional(),
  // SLI reference - links to the metric definition with framework mappings
  sliId: z.string().optional(),
  // Inline SLI fields (for simple cases without separate SLI definition)
  type: z.enum(['quantitative', 'qualitative']).optional(),
  metricName: z.string().optional(),
  unit: z.string().optional(),
  // SLO target for this level
  operator: z.enum(['gte', 'lte', 'gt', 'lt', 'eq', 'exists']),
  target: z.number(),
  status: z.string().optional(),
  // Deprecated: use SLI.frameworkMappings instead
  frameworkMappings: z.array(FrameworkMappingSchema).optional(),
  current: z.number().optional(),
  isMet: z.boolean().optional(),
  weight: z.number().optional(),
  required: z.boolean().optional(),
});

export type Criterion = z.infer<typeof CriterionSchema>;

// Enabler (implementation task) within a maturity level
export const EnablerSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().optional(),
  type: z.string().optional(),
  layer: z.string().optional(),
  effort: z.string().optional(),
  team: z.string().optional(),
  status: z.string().optional(),
  criteriaIds: z.array(z.string()).optional(),
  dependsOn: z.array(z.string()).optional(),
});

export type Enabler = z.infer<typeof EnablerSchema>;

// Maturity level (M1-M5)
export const LevelSchema = z.object({
  level: z.number().int(),
  name: z.string(),
  description: z.string(),
  criteria: z.array(CriterionSchema).optional(),
  enablers: z.array(EnablerSchema).optional(),
});

export type Level = z.infer<typeof LevelSchema>;

// Domain model (security, operations, quality)
export const DomainModelSchema = z.object({
  name: z.string(),
  description: z.string().optional(),
  owner: z.string().optional(),
  levels: z.array(LevelSchema),
});

export type DomainModel = z.infer<typeof DomainModelSchema>;

// Domain assessment (current state)
export const DomainAssessmentSchema = z.object({
  domain: z.string(),
  assessedAt: z.string().optional(),
  assessedBy: z.string().optional(),
  currentLevel: z.number().int(),
  targetLevel: z.number().int(),
  criteriaValues: z.record(z.string(), z.number()).optional(),
  criteriaStatus: z.record(z.string(), z.string()).optional(),
  enablerStatus: z.record(z.string(), z.string()).optional(),
});

export type DomainAssessment = z.infer<typeof DomainAssessmentSchema>;

// Level thresholds for KPIs
export const LevelThresholdsSchema = z.object({
  m1: z.any().optional(),
  m2: z.any().optional(),
  m3: z.any().optional(),
  m4: z.any().optional(),
  m5: z.any().optional(),
});

export type LevelThresholds = z.infer<typeof LevelThresholdsSchema>;

// KPI threshold definition
export const KPIThresholdSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().optional(),
  unit: z.string().optional(),
  operator: z.string().optional(),
  thresholds: LevelThresholdsSchema,
  current: z.any().optional(),
});

export type KPIThreshold = z.infer<typeof KPIThresholdSchema>;

// Spec metadata
export const SpecMetadataSchema = z.object({
  name: z.string().optional(),
  description: z.string().optional(),
  version: z.string().optional(),
  organization: z.string().optional(),
  createdAt: z.string().optional(),
  updatedAt: z.string().optional(),
});

export type SpecMetadata = z.infer<typeof SpecMetadataSchema>;

// Full maturity specification
export const SpecSchema = z.object({
  $schema: z.string().optional(),
  metadata: SpecMetadataSchema.optional(),
  slis: z.record(z.string(), SLISchema).optional(), // Service Level Indicators with framework mappings
  kpiThresholds: z.record(z.string(), z.array(KPIThresholdSchema)).optional(), // Deprecated: use SLIs
  domains: z.record(z.string(), DomainModelSchema),
  assessments: z.record(z.string(), DomainAssessmentSchema).optional(),
});

export type Spec = z.infer<typeof SpecSchema>;

// Validation helper
export function validateSpec(data: unknown): Spec {
  return SpecSchema.parse(data);
}

// Safe validation helper (returns result instead of throwing)
export function safeValidateSpec(data: unknown): z.SafeParseReturnType<unknown, Spec> {
  return SpecSchema.safeParse(data);
}
