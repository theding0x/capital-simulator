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
	// Ch. 10 — The Working-Day
	s.HandleFunc("POST /v1/working-days", h.CreateWorkingDay)
	s.HandleFunc("GET /v1/working-days/{id}", h.GetWorkingDay)
	s.HandleFunc("POST /v1/working-days/validate", h.ValidateWorkingDay)
	s.HandleFunc("POST /v1/relay-schedules", h.CreateRelaySchedule)
	s.HandleFunc("GET /v1/relay-schedules/{id}", h.GetRelaySchedule)
	// Ch. 13 — Co-operation
	s.HandleFunc("POST /v1/cooperations", h.CreateCooperation)
	s.HandleFunc("GET /v1/cooperations", h.ListCooperations)
	s.HandleFunc("GET /v1/cooperations/{id}", h.GetCooperation)
	s.HandleFunc("POST /v1/cooperations/{id}/collective-working-day", h.ComputeCollectiveWorkingDay)
	s.HandleFunc("POST /v1/cooperations/{id}/average-social-labour", h.ComputeAverageSocialLabour)
	s.HandleFunc("POST /v1/cooperations/minimum-capital", h.ComputeMinimumCapital)
	// Ch. 14 — Division of Labour and Manufacture
	s.HandleFunc("POST /v1/manufactures", h.CreateManufacture)
	s.HandleFunc("GET /v1/manufactures", h.ListManufactures)
	s.HandleFunc("GET /v1/manufactures/{id}", h.GetManufacture)
	s.HandleFunc("POST /v1/manufactures/{id}/proportional-group-size", h.ProportionalGroupSize)
	s.HandleFunc("POST /v1/manufactures/{id}/scale", h.ScaleManufacture)
	s.HandleFunc("GET /v1/manufactures/{id}/collective-labourer", h.GetCollectiveLabourer)
	s.HandleFunc("GET /v1/manufactures/{id}/minimum-capital", h.GetManufactureMinimumCapital)
	// Ch. 17 — Changes of Magnitude in the Price of Labour-Power and in Surplus-Value
	s.HandleFunc("POST /v1/labour-scenarios", h.ComputeLabourScenario)
	// Ch. 19 — The Transformation of the Value of Labour-Power into Wages
	s.HandleFunc("POST /v1/wage-forms", h.CreateWageForm)
	s.HandleFunc("GET /v1/wage-forms", h.ListWageForms)
	s.HandleFunc("GET /v1/wage-forms/{agentID}", h.GetWageForm)
	// Ch. 20 — Time-Wages
	s.HandleFunc("POST /v1/time-wages/hourly-price", h.ComputeHourlyPrice)
	s.HandleFunc("POST /v1/time-wages/sessions", h.CreateWorkingSession)
	s.HandleFunc("GET /v1/time-wages/sessions/{id}", h.GetWorkingSession)
	s.HandleFunc("GET /v1/agents/{id}/time-wages/sessions", h.ListWorkingSessions)
	// Ch. 21 — Piece-Wages
	s.HandleFunc("POST /v1/agents/{id}/piece-wages", h.CreatePieceWage)
	s.HandleFunc("GET /v1/agents/{id}/piece-wages", h.GetPieceWage)
	s.HandleFunc("POST /v1/piece-price", h.ComputePiecePrice)
	s.HandleFunc("POST /v1/sub-contracts", h.CreateSubContract)
	s.HandleFunc("GET /v1/sub-contracts/{id}", h.GetSubContract)
	// Ch. 22 — National Differences in Wages
	s.HandleFunc("POST /v1/intensities", h.RegisterIntensity)
	s.HandleFunc("GET /v1/intensities", h.ListIntensities)
	s.HandleFunc("POST /v1/wages", h.RegisterDayWage)
	s.HandleFunc("GET /v1/wages/{country}/standardised", h.GetStandardisedWage)
	s.HandleFunc("GET /v1/comparisons", h.GetWageComparison)
}
