package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

func Register(s *httpx.Server, h *Handler) {
	s.HandleFunc("POST /v1/surplus/mass", h.ComputeMass)
	s.HandleFunc("GET /v1/surplus/limits", h.GetLimits)
}
