package auth

import "testing"

func TestIsComputePath(t *testing.T) {
	t.Parallel()
	cases := []struct {
		path string
		want bool
	}{
		// Basic exact match.
		{"/v1/observatory/levers", true},
		// Extra segment — no such endpoint; exact-length rejects it.
		{"/v1/observatory/levers/extra", false},
		// Persisting CRUD roots — always blocked.
		{"/v1/commodities", false},
		{"/v1/owners", false},
		// GET-only path, not a write allowlist entry.
		{"/v1/observatory/snapshot", false},
		// Templated entries — positive: ID segment matches wildcard.
		{"/v1/cooperations/abc123/collective-working-day", true},
		{"/v1/cooperations/abc123/average-social-labour", true},
		{"/v1/manufactures/xyz/scale", true},
		{"/v1/manufactures/xyz/proportional-group-size", true},
		{"/v1/colonial-markets/abc123/independence", true},
		// Templated entries — negative: persisting siblings and roots must stay blocked.
		{"/v1/cooperations", false},
		{"/v1/cooperations/abc123", false},
		{"/v1/manufactures", false},
		{"/v1/manufactures/xyz", false},
		{"/v1/colonial-markets", false},
		{"/v1/colonial-markets/abc123", false},
		// Wildcard must not match empty segment.
		{"/v1/cooperations//collective-working-day", false},
	}
	for _, c := range cases {
		if got := IsComputePath(c.path); got != c.want {
			t.Errorf("IsComputePath(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}

// TestComputeAllowlistAudited locks in the full Task 8 audit result.
// mustAllow: confirmed compute (no store mutation) paths from the audit.
// mustBlock: persisting CRUD roots that must stay owner-only.
func TestComputeAllowlistAudited(t *testing.T) {
	t.Parallel()
	mustAllow := []string{
		// simulation-engine — stateless calculators
		"/v1/production/extra-surplus-value",
		"/v1/production/working-day",
		"/v1/production/working-day/shorten",
		"/v1/production/relative-surplus",
		"/v1/production/relative-surplus-from-productivity",
		"/v1/surplus-value/absolute",
		"/v1/surplus-value/relative",
		"/v1/surplus-value/rates",
		"/v1/surplus/mass",
		"/v1/reproductions/simple",
		"/v1/reproductions/repayment-period",
		"/v1/reproductions/extended",
		"/v1/reproductions/split-surplus",
		"/v1/accumulation/organic-composition",
		"/v1/accumulation/labour-demand",
		"/v1/accumulation/reserve-army",
		"/v1/accumulation/centralisation",
		"/v1/statutory-wages/compare",
		"/v1/farm-tenures/real-rent",
		"/v1/market-formation",
		// agent-service — stateless calculators
		"/v1/labour-scenarios",
		"/v1/cooperations/minimum-capital",
		"/v1/circuit-probes",
		"/v1/exchange-simulations",
		"/v1/piece-price",
		"/v1/time-wages/hourly-price",
		"/v1/time-wages/nominal-wage",
		// agent-service — templated compute (read+compute, no mutation)
		"/v1/cooperations/abc123/collective-working-day",
		"/v1/cooperations/abc123/average-social-labour",
		"/v1/manufactures/abc123/proportional-group-size",
		"/v1/manufactures/abc123/scale",
		"/v1/colonial-markets/abc123/independence",
		// finance-service — stateless calculators
		"/v1/profit/profit-form",
		"/v1/profit/compare",
		"/v1/avgprofit/social-aggregate",
		"/v1/avgprofit/compensation-ground",
		"/v1/merchant/turnover-effect",
		"/v1/credit/interest-rate-analysis",
		// simulation-engine — observatory
		"/v1/observatory/levers",
	}
	mustBlock := []string{
		// persisting CRUD — must stay owner-only
		"/v1/commodities",
		"/v1/owners",
		"/v1/agents",
		"/v1/profit/rate",
		"/v1/profit/variation",
		"/v1/avgprofit/spheres",
		"/v1/avgprofit/general-rate",
		"/v1/credit/bills-of-exchange",
		"/v1/rent/parcels",
		"/v1/tendency/trajectory",
		"/v1/merchant/commercial-capital",
		"/v1/engine/start",
		"/v1/engine/stop",
		// templated roots and bare ID paths must stay blocked
		"/v1/cooperations",
		"/v1/cooperations/abc123",
		"/v1/manufactures",
		"/v1/manufactures/abc123",
		"/v1/colonial-markets",
		"/v1/colonial-markets/abc123",
	}
	for _, p := range mustAllow {
		if !IsComputePath(p) {
			t.Errorf("expected compute path on allowlist: %s", p)
		}
	}
	for _, p := range mustBlock {
		if IsComputePath(p) {
			t.Errorf("persisting path must NOT be on compute allowlist: %s", p)
		}
	}
}
