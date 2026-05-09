/**
 * Styles exports
 *
 * Provides CSS as importable strings for embedding.
 */

// CSS file path for direct import
export const PRISM_CSS_PATH = new URL('./prism.css', import.meta.url).href;

// Inline CSS for embedding (populated at build time or can be imported)
export const prismStyles = `
/* See prism.css for full stylesheet */
/* This is a minimal inline version for basic rendering */

.prism-domain-view, .prism-framework-view {
  font-family: system-ui, sans-serif;
  line-height: 1.6;
  max-width: 1200px;
  margin: 0 auto;
}

.prism-table {
  width: 100%;
  border-collapse: collapse;
}

.prism-table th, .prism-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.prism-table th {
  background: #f3f4f6;
  font-weight: 600;
}

.prism-criterion-met { background: #f0fdf4; }
.prism-criterion-pending { background: #fefce8; }
`;
