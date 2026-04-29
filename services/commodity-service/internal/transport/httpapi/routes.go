package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

// Register attaches the commodity-service routes to s. Patterns use Go 1.22
// ServeMux method+path syntax.
func Register(s *httpx.Server, h *Handler) {
	s.HandleFunc("POST /v1/commodities", h.Create)
	s.HandleFunc("GET /v1/commodities", h.List)
	s.HandleFunc("GET /v1/commodities/{id}", h.Get)
	s.HandleFunc("PATCH /v1/commodities/{id}", h.Update)
	s.HandleFunc("DELETE /v1/commodities/{id}", h.Delete)
	s.HandleFunc("POST /v1/commodities/{id}/value", h.Value)
	s.HandleFunc("GET /v1/commodities/{id}/value-form", h.ValueForm)
	s.HandleFunc("GET /v1/commodities/{id}/social-relations", h.SocialRelations)
	s.HandleFunc("POST /v1/exchange-ratio", h.ExchangeRatio)
}
