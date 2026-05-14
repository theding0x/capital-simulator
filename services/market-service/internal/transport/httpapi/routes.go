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
}
