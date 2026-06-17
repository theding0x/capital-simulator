package auth

import "strings"

// computePaths are write-method endpoints that perform NO database mutation —
// pure calculators and the per-session in-memory Atlas observatory. Non-owners
// (guests) may call these so they can play with the simulation.
//
// FAIL-CLOSED: anything not listed here is owner-only for write methods. The
// only failure mode is "a forgotten compute panel degrades to owner-only for
// guests", never "a forgotten write leaks to guests". Add a prefix here only
// after confirming the handler does not call a store mutation (see Task 8).
//
// Task 8 audit: 109 POST routes classified; 37 confirmed compute (listed here),
// 72 persisting (omitted). All finance/credit/rent/tendency/merchant Create*
// handlers call store.Create* and are owner-only. Engine start/stop mutate
// in-memory ticker state and are owner-only. All agent/sim Create* handlers
// that persist are owner-only.
var computePaths = []string{
	// Atlas observatory — per-session in-memory, zero DB writes (simulation-engine)
	"/v1/observatory/levers",

	// Vol. I Ch. 11 — Rate and Mass of Surplus-Value (simulation-engine, handler.go)
	"/v1/surplus/mass",

	// Vol. I Ch. 12 — Relative Surplus-Value working-day calculators (simulation-engine)
	"/v1/production/working-day",
	"/v1/production/working-day/shorten",
	"/v1/production/extra-surplus-value",

	// Vol. I Ch. 12/13/14 bridge — productivity → relative surplus (simulation-engine)
	"/v1/production/relative-surplus",
	"/v1/production/relative-surplus-from-productivity",

	// Vol. I Ch. 16 — Absolute and Relative Surplus-Value (simulation-engine)
	"/v1/surplus-value/absolute",
	"/v1/surplus-value/relative",

	// Vol. I Ch. 18 — Three Formulae for Rate of Surplus-Value (simulation-engine)
	"/v1/surplus-value/rates",

	// Vol. I Ch. 17 — Labour-scenario calculator (agent-service)
	"/v1/labour-scenarios",

	// Vol. I Ch. 13 — Co-operation compute endpoints (agent-service, read store then compute)
	"/v1/cooperations/minimum-capital",
	"/v1/cooperations/{id}/collective-working-day",
	"/v1/cooperations/{id}/average-social-labour",

	// Vol. I Ch. 14 — Manufacture compute endpoints (agent-service, read store then compute)
	"/v1/manufactures/{id}/proportional-group-size",
	"/v1/manufactures/{id}/scale",

	// Vol. I Ch. 4 — Circuit probe and exchange simulation (agent-service, pure compute)
	"/v1/circuit-probes",
	"/v1/exchange-simulations",

	// Vol. I Ch. 20/21 — Piece-price and time-wage calculators (agent-service)
	"/v1/piece-price",
	"/v1/time-wages/hourly-price",
	"/v1/time-wages/nominal-wage",

	// Vol. I Ch. 23 — Simple Reproduction compute (simulation-engine, stateless)
	"/v1/reproductions/simple",
	"/v1/reproductions/repayment-period",

	// Vol. I Ch. 24 — Extended Reproduction compute (simulation-engine, stateless)
	"/v1/reproductions/extended",
	"/v1/reproductions/split-surplus",

	// Vol. I Ch. 25 — General Law compute (simulation-engine, stateless)
	"/v1/accumulation/organic-composition",
	"/v1/accumulation/labour-demand",
	"/v1/accumulation/reserve-army",

	// Vol. I Ch. 32 — Historical tendency centralisation (simulation-engine, stateless)
	"/v1/accumulation/centralisation",

	// Vol. I Ch. 28/29/30 — Historical-stage compute (simulation-engine, reads store then computes)
	"/v1/statutory-wages/compare",
	"/v1/farm-tenures/real-rent",
	"/v1/market-formation",

	// Vol. I Ch. 33 — Colonial market independence compute (simulation-engine, reads store)
	"/v1/colonial-markets/{id}/independence",

	// Vol. III Ch. 1 — Profit-form calculator (finance-service, pure compute)
	"/v1/profit/profit-form",

	// Vol. III Ch. 3 — Compare profit rates (finance-service, pure compute)
	"/v1/profit/compare",

	// Vol. III Ch. 9 — Social capital aggregate (finance-service, pure compute, no store)
	"/v1/avgprofit/social-aggregate",

	// Vol. III Ch. 12 — Compensation ground (finance-service, pure compute)
	"/v1/avgprofit/compensation-ground",

	// Vol. III Ch. 18 — Merchant turnover effect (finance-service, pure compute)
	"/v1/merchant/turnover-effect",

	// Vol. III Ch. 22 — Interest rate analysis (finance-service, reads store then computes)
	"/v1/credit/interest-rate-analysis",
}

// IsComputePath reports whether path is a non-persisting compute endpoint that
// non-owners may invoke with a write method.
//
// Matching is exact-length segment-by-segment: the entry and the request path
// must have the same number of "/" segments. A segment of the form {anything}
// is a wildcard that matches any single non-empty path segment. This lets
// templated entries like /v1/cooperations/{id}/collective-working-day match
// real runtime paths while keeping persisting siblings (e.g. /v1/cooperations)
// blocked.
func IsComputePath(path string) bool {
	pathSegs := splitSegs(path)
	for _, entry := range computePaths {
		if matchSegs(splitSegs(entry), pathSegs) {
			return true
		}
	}
	return false
}

func splitSegs(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

func matchSegs(entry, path []string) bool {
	if len(entry) != len(path) {
		return false
	}
	for i, e := range entry {
		if strings.HasPrefix(e, "{") && strings.HasSuffix(e, "}") {
			if path[i] == "" {
				return false
			}
			continue
		}
		if e != path[i] {
			return false
		}
	}
	return true
}
