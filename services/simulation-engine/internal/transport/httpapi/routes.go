package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

func Register(s *httpx.Server, h *Handler) {
	// Ch. 11 — Rate and Mass of Surplus-Value
	s.HandleFunc("POST /v1/surplus/mass", h.ComputeMass)
	s.HandleFunc("GET /v1/surplus/limits", h.GetLimits)

	// Ch. 12 — Relative Surplus-Value
	s.HandleFunc("POST /v1/production/working-day", h.RecordWorkingDay)
	s.HandleFunc("POST /v1/production/working-day/shorten", h.ShortenWorkingDay)
	s.HandleFunc("GET /v1/production/rate-of-surplus-value", h.GetProductionRate)
	s.HandleFunc("POST /v1/production/extra-surplus-value", h.ComputeExtraSurplusValue)

	// Part IV bridge — Ch. 13/14/15 productivity → Ch. 12 relative surplus-value
	s.HandleFunc("POST /v1/production/relative-surplus", h.RelativeSurplus)
	s.HandleFunc("POST /v1/production/relative-surplus-from-productivity", h.RelativeSurplusFromProductivity)

	// Ch. 15 — Machinery and Modern Industry
	s.HandleFunc("POST /v1/machines", h.CreateMachine)
	s.HandleFunc("GET /v1/machines", h.ListMachines)
	s.HandleFunc("GET /v1/machines/{id}", h.GetMachine)
	s.HandleFunc("GET /v1/machines/{id}/wear", h.GetMachineWear)
	s.HandleFunc("POST /v1/factories", h.CreateFactory)
	s.HandleFunc("GET /v1/factories", h.ListFactories)
	s.HandleFunc("GET /v1/factories/{id}", h.GetFactory)
	s.HandleFunc("POST /v1/factories/{id}/tick", h.TickFactory)
	s.HandleFunc("GET /v1/factories/{id}/ticks", h.ListFactoryTicks)

	// Ch. 16 — Absolute and Relative Surplus-Value
	s.HandleFunc("POST /v1/surplus-value/absolute", h.ComputeAbsoluteSurplusValue)
	s.HandleFunc("POST /v1/surplus-value/relative", h.ComputeRelativeSurplusValue)
	s.HandleFunc("GET /v1/surplus-value/rate", h.GetSurplusValueRate)

	// Ch. 18 — Various Formula for the Rate of Surplus-Value
	s.HandleFunc("POST /v1/surplus-value/rates", h.ComputeRatesOfSurplusValue)

	// Ch. 23 — Simple Reproduction
	s.HandleFunc("POST /v1/reproductions/simple", h.RunSimpleReproduction)
	s.HandleFunc("POST /v1/reproductions/repayment-period", h.ComputeRepaymentPeriod)

	// Ch. 24 — The Transformation of Surplus-Value into Capital
	s.HandleFunc("POST /v1/reproductions/extended", h.RunExtendedReproduction)
	s.HandleFunc("POST /v1/reproductions/split-surplus", h.SplitSurplusValue)
}
