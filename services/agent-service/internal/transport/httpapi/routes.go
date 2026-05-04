package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

func Register(s *httpx.Server, h *Handler) {
	s.HandleFunc("POST /v1/agents", h.Create)
	s.HandleFunc("GET /v1/agents", h.List)
	s.HandleFunc("GET /v1/agents/{id}", h.Get)
	s.HandleFunc("PATCH /v1/agents/{id}", h.Update)
	s.HandleFunc("DELETE /v1/agents/{id}", h.Delete)
	s.HandleFunc("POST /v1/agents/{id}/circuits", h.CreateCircuit)
	s.HandleFunc("GET /v1/agents/{id}/circuits", h.ListCircuits)
	s.HandleFunc("POST /v1/agents/{id}/reinvest", h.Reinvest)
	s.HandleFunc("POST /v1/agents/{id}/hoard", h.Hoard)
	s.HandleFunc("POST /v1/circuits", h.ComputeCircuit)
	s.HandleFunc("POST /v1/exchange-simulations", h.ComputeExchange)
}
