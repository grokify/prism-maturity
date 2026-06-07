package render

import (
	"io"

	capstack "github.com/grokify/prism-capability"
	caprender "github.com/grokify/prism-capability/render"
	"github.com/grokify/prism-maturity/dashboard"
)

// LitOptionsWithMaturity creates LitOptions with maturity overlay data.
func LitOptionsWithMaturity(baseOpts caprender.LitOptions, agg *dashboard.MaturityAggregator) caprender.LitOptions {
	opts := baseOpts
	opts.Overlays = BuildMaturityOverlay(agg)
	return opts
}

// RenderLitHTMLWithMaturity renders a capability stack with maturity data to a Lit HTML page.
func RenderLitHTMLWithMaturity(w io.Writer, doc *capstack.CapabilityStack, opts caprender.LitOptions, agg *dashboard.MaturityAggregator) error {
	opts.Overlays = BuildMaturityOverlay(agg)
	return caprender.RenderLitHTML(w, doc, opts)
}

// RenderJSONWithMaturity renders a capability stack with maturity data as JSON.
func RenderJSONWithMaturity(w io.Writer, doc *capstack.CapabilityStack, opts caprender.LitOptions, agg *dashboard.MaturityAggregator) error {
	opts.Overlays = BuildMaturityOverlay(agg)
	return caprender.RenderJSON(w, doc, opts)
}

// ToLitGridDataWithMaturity converts a CapabilityStack to LitGridData with maturity overlay.
func ToLitGridDataWithMaturity(doc *capstack.CapabilityStack, opts caprender.LitOptions, agg *dashboard.MaturityAggregator) *caprender.LitGridData {
	opts.Overlays = BuildMaturityOverlay(agg)
	return caprender.ToLitGridData(doc, opts)
}
