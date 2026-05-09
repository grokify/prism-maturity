/**
 * PRISM Maturity Viewer
 *
 * TypeScript library for rendering PRISM maturity specifications.
 * Provides HTML rendering components for domain and framework views.
 *
 * @packageDocumentation
 */

// Schema types and validators
export * from './schema';

// View components
export * from './components/domain-view';
export * from './components/framework-view';

// Re-export common types for convenience
export type {
  Spec,
  SLI,
  DomainModel,
  Level,
  Criterion,
  Enabler,
  DomainAssessment,
  FrameworkMapping,
  SpecMetadata,
} from './schema/maturity';
