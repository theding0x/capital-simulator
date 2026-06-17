package auth

import "testing"

func TestIsComputePath(t *testing.T) {
	t.Parallel()
	cases := []struct {
		path string
		want bool
	}{
		{"/v1/observatory/levers", true},
		{"/v1/observatory/levers/extra", true},
		{"/v1/commodities", false},
		{"/v1/owners", false},
		{"/v1/observatory/snapshot", false}, // GET-only path, not a write allowlist entry
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
		// finance-service — stateless calculators
		"/v1/profit/profit-form",
		"/v1/profit/compare",
		"/v1/avgprofit/social-aggregate",
		"/v1/avgprofit/compensation-ground",
		"/v1/merchant/turnover-effect",
		"/v1/credit/interest-rate-analysis",
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
