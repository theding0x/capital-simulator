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

	// Vol. III Ch. 8 — Different Compositions of Capitals in Different Branches
	s.HandleFunc("POST /v1/avgprofit/spheres", h.CreateProductionSphere)
	s.HandleFunc("GET /v1/avgprofit/spheres", h.ListProductionSpheres)
	s.HandleFunc("GET /v1/avgprofit/spheres/{id}", h.GetProductionSphere)

	// Vol. III Ch. 9 — Formation of a General Rate of Profit
	s.HandleFunc("POST /v1/avgprofit/general-rate", h.CreateGeneralProfitRate)
	s.HandleFunc("GET /v1/avgprofit/general-rate/{id}", h.GetGeneralProfitRate)
	s.HandleFunc("POST /v1/avgprofit/price-of-production", h.CreatePriceOfProduction)
	s.HandleFunc("GET /v1/avgprofit/price-of-production", h.ListPricesOfProduction)
	s.HandleFunc("GET /v1/avgprofit/price-of-production/{id}", h.GetPriceOfProduction)
	s.HandleFunc("POST /v1/avgprofit/social-aggregate", h.ComputeSocialAggregate)

	// Vol. III Ch. 10 — Equalisation of the General Rate of Profit Through Competition
	s.HandleFunc("POST /v1/avgprofit/market-value", h.CreateMarketValue)
	s.HandleFunc("GET /v1/avgprofit/market-value/{id}", h.GetMarketValue)
	s.HandleFunc("POST /v1/avgprofit/surplus-profit", h.CreateSurplusProfit)
	s.HandleFunc("GET /v1/avgprofit/surplus-profit/{id}", h.GetSurplusProfit)
	s.HandleFunc("POST /v1/avgprofit/capital-flow", h.CreateCapitalFlow)
	s.HandleFunc("GET /v1/avgprofit/equalisation/{id}", h.GetEqualisation)

	// Vol. III Ch. 11 — Effects of General Wage Fluctuations on Prices of Production
	s.HandleFunc("POST /v1/avgprofit/wage-effect", h.CreateWageEffectAnalysis)
	s.HandleFunc("GET /v1/avgprofit/wage-effect", h.ListWageEffectAnalyses)
	s.HandleFunc("GET /v1/avgprofit/wage-effect/{id}", h.GetWageEffectAnalysis)

	// Vol. III Ch. 12 — Supplementary Remarks (on prices of production)
	s.HandleFunc("POST /v1/avgprofit/price-change", h.CreatePriceOfProductionChange)
	s.HandleFunc("GET /v1/avgprofit/price-change/{id}", h.GetPriceOfProductionChange)
	s.HandleFunc("POST /v1/avgprofit/compensation-ground", h.ComputeCompensationGround)
	s.HandleFunc("GET /v1/avgprofit/summary", h.GetPartIISummary)
}
