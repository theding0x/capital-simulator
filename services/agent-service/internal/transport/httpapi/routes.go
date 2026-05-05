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
	s.HandleFunc("POST /v1/circuit-probes", h.ComputeCircuit)
	s.HandleFunc("POST /v1/exchange-simulations", h.ComputeExchange)
	// Ch. 6 — The Buying and Selling of Labour-Power
	s.HandleFunc("POST /v1/workers", h.CreateWorker)
	s.HandleFunc("GET /v1/workers", h.ListWorkers)
	s.HandleFunc("GET /v1/workers/{id}", h.GetWorker)
	s.HandleFunc("POST /v1/capitalists", h.CreateCapitalist)
	s.HandleFunc("GET /v1/capitalists", h.ListCapitalists)
	s.HandleFunc("GET /v1/capitalists/{id}", h.GetCapitalist)
	s.HandleFunc("POST /v1/labour-power/offerings", h.CreateOffering)
	s.HandleFunc("GET /v1/labour-power/offerings", h.ListOfferings)
	s.HandleFunc("POST /v1/labour-power/purchases", h.CreatePurchase)
	s.HandleFunc("GET /v1/labour-power/purchases", h.ListPurchases)
	s.HandleFunc("GET /v1/labour-power/purchases/{id}", h.GetPurchase)
	// Ch. 7 — The Labour-Process and the Production of Surplus-Value
	s.HandleFunc("POST /v1/labour-processes", h.RunLabourProcess)
	s.HandleFunc("GET /v1/labour-processes/{id}", h.GetLabourProcessRecord)
}
