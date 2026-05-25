package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

// Register attaches the market-service routes to s. Patterns use Go 1.22
// ServeMux method+path syntax.
func Register(s *httpx.Server, h *Handler) {
	// Owners
	s.HandleFunc("POST /v1/owners", h.CreateOwner)
	s.HandleFunc("GET /v1/owners", h.ListOwners)
	s.HandleFunc("GET /v1/owners/{id}", h.GetOwner)

	// Offers
	s.HandleFunc("POST /v1/offers", h.CreateOffer)
	s.HandleFunc("GET /v1/offers", h.ListOffers)
	s.HandleFunc("DELETE /v1/offers/{id}", h.DeleteOffer)

	// Exchanges
	s.HandleFunc("POST /v1/exchanges", h.CreateExchange)
	s.HandleFunc("GET /v1/exchanges", h.ListExchanges)
	s.HandleFunc("GET /v1/exchanges/{id}", h.GetExchange)

	// Universal equivalent
	s.HandleFunc("POST /v1/universal-equivalent", h.SetUniversalEquivalent)
	s.HandleFunc("GET /v1/universal-equivalent", h.GetUniversalEquivalent)

	// Money commodity
	s.HandleFunc("POST /v1/money-commodity", h.SetMoneyCommodity)
	s.HandleFunc("GET /v1/money-commodity", h.GetMoneyCommodity)

	// Prices
	s.HandleFunc("POST /v1/prices", h.ComputePrice)
	s.HandleFunc("GET /v1/prices", h.ListPrices)
	s.HandleFunc("GET /v1/prices/{commodityID}", h.GetPrice)

	// Vol. II Ch. 6 — Costs of circulation
	s.HandleFunc("POST /v1/circulation-costs", h.CreateCirculationCost)
	s.HandleFunc("GET /v1/circulation-costs/aggregate", h.AggregateCirculationCosts)
	s.HandleFunc("GET /v1/circulation-costs/system-faux-frais", h.SystemFauxFrais)
	s.HandleFunc("GET /v1/circulation-costs/{id}", h.GetCirculationCost)
	s.HandleFunc("GET /v1/circulation-costs", h.ListCirculationCosts)
	s.HandleFunc("POST /v1/circulation-agents", h.CreateCirculationAgent)
	s.HandleFunc("POST /v1/money-as-faux-frais", h.CreateMoneyAsFauxFrais)
	s.HandleFunc("POST /v1/commodity-supplies", h.CreateCommoditySupply)
	s.HandleFunc("POST /v1/commodity-supplies/{id}/storage-cost", h.AddStorageCost)
	s.HandleFunc("POST /v1/transport-tariffs", h.SetTransportTariff)
	s.HandleFunc("POST /v1/transport-legs", h.CreateTransportLeg)

	// Vol. II Ch. 5 — Turnover time and circulation phases
	s.HandleFunc("POST /v1/turnover-time", h.CreateTurnoverTime)
	s.HandleFunc("GET /v1/turnover-time", h.ListTurnoverTimes)
	s.HandleFunc("GET /v1/turnover-time/{id}", h.GetTurnoverTime)
	s.HandleFunc("POST /v1/turnover-time/{id}/labour-time", h.AddLabourTime)
	s.HandleFunc("POST /v1/turnover-time/{id}/labour-interruption", h.AddLabourInterruption)
	s.HandleFunc("POST /v1/turnover-time/{id}/latent-mp", h.RecordLatentMP)
	s.HandleFunc("POST /v1/turnover-time/{id}/natural-process", h.RecordNaturalProcess)
	s.HandleFunc("POST /v1/turnover-time/{id}/selling-phase", h.RecordSellingPhase)
	s.HandleFunc("POST /v1/turnover-time/{id}/buying-phase", h.RecordBuyingPhase)
	s.HandleFunc("GET /v1/turnover-time/{id}/active-fraction", h.GetActiveFraction)
	s.HandleFunc("POST /v1/perishability", h.SetPerishability)
	s.HandleFunc("POST /v1/market-separation", h.SetMarketSeparation)
}
