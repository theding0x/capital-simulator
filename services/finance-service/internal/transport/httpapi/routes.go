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
}
