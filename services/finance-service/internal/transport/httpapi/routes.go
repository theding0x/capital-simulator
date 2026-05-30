package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

// Register wires finance-service routes onto s. Each Vol. III chapter PR
// appends its route block here.
func Register(s *httpx.Server, h *Handler) {
	// Vol. III Ch. 1 — Cost-Price and Profit
	s.HandleFunc("POST /v1/profit/cost-price", h.CreateCostPrice)
	s.HandleFunc("GET /v1/profit/cost-price", h.ListCostPrices)
	s.HandleFunc("GET /v1/profit/cost-price/{id}", h.GetCostPrice)
	s.HandleFunc("POST /v1/profit/profit-form", h.ComputeProfitForm)

	// Vol. III Ch. 2 — The Rate of Profit
	s.HandleFunc("POST /v1/profit/rate", h.CreateProfitRate)
	s.HandleFunc("GET /v1/profit/rate", h.ListProfitRates)
	s.HandleFunc("GET /v1/profit/rate/{id}", h.GetProfitRate)

	// Vol. III Ch. 3 — Relation of the Rate of Profit to the Rate of Surplus-Value
	s.HandleFunc("POST /v1/profit/variation", h.CreateVariation)
	s.HandleFunc("GET /v1/profit/variation", h.ListVariations)
	s.HandleFunc("GET /v1/profit/variation/{id}", h.GetVariation)
	s.HandleFunc("POST /v1/profit/compare", h.CompareProfitRates)

	// Vol. III Ch. 4 — The Effect of the Turnover on the Rate of Profit
	s.HandleFunc("POST /v1/profit/turnover-analysis", h.CreateTurnoverAnalysis)
	s.HandleFunc("GET /v1/profit/turnover-analysis", h.ListTurnoverAnalyses)
	s.HandleFunc("GET /v1/profit/turnover-analysis/{id}", h.GetTurnoverAnalysis)

	// Vol. III Ch. 5 — Economy in the Employment of Constant Capital
	s.HandleFunc("POST /v1/profit/economy", h.CreateEconomyAnalysis)
	s.HandleFunc("GET /v1/profit/economy", h.ListEconomyAnalyses)
	s.HandleFunc("GET /v1/profit/economy/{id}", h.GetEconomyAnalysis)

	// Vol. III Ch. 6 — The Effect of Price Fluctuation on the Rate of Profit
	s.HandleFunc("POST /v1/profit/price-fluctuation", h.CreatePriceFluctuationAnalysis)
	s.HandleFunc("GET /v1/profit/price-fluctuation", h.ListPriceFluctuationAnalyses)
	s.HandleFunc("GET /v1/profit/price-fluctuation/{id}", h.GetPriceFluctuationAnalysis)

	// Vol. III Ch. 7 — Supplementary Remarks (on the rate of profit)
	s.HandleFunc("POST /v1/profit/composition-effect", h.CreateCompositionEffect)
	s.HandleFunc("GET /v1/profit/composition-effect", h.ListCompositionEffects)
	s.HandleFunc("GET /v1/profit/composition-effect/{id}", h.GetCompositionEffect)
	s.HandleFunc("POST /v1/profit/magnitude-change", h.CreateMagnitudeChange)
	s.HandleFunc("GET /v1/profit/magnitude-change", h.ListMagnitudeChanges)
	s.HandleFunc("GET /v1/profit/magnitude-change/{id}", h.GetMagnitudeChange)
	s.HandleFunc("GET /v1/profit/summary", h.GetPartISummary)
}
