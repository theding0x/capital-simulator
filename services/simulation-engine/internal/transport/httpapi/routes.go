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

	// Ch. 25 — The General Law of Capitalist Accumulation
	s.HandleFunc("POST /v1/accumulation/organic-composition", h.ComputeOrganicComposition)
	s.HandleFunc("POST /v1/accumulation/labour-demand", h.ComputeLabourDemandEndpoint)
	s.HandleFunc("POST /v1/accumulation/reserve-army", h.ComputeReserveArmyEndpoint)
	s.HandleFunc("POST /v1/accumulation/scenarios", h.CreateGeneralLawScenario)
	s.HandleFunc("GET /v1/accumulation/scenarios/{id}", h.GetGeneralLawScenario)

	// Ch. 26 — The Secret of Primitive Accumulation
	s.HandleFunc("POST /v1/historical-stages", h.CreateHistoricalStage)
	s.HandleFunc("GET /v1/historical-stages", h.ListHistoricalStages)
	s.HandleFunc("POST /v1/historical-stages/{id}/seed-scenario", h.SeedScenarioFromStage)

	// Ch. 27 — Expropriation of the Agricultural Population
	s.HandleFunc("POST /v1/enclosure-events", h.CreateEnclosureEvent)
	s.HandleFunc("GET /v1/enclosure-events", h.ListEnclosureEvents)

	// Ch. 28 — Bloody Legislation Against the Expropriated
	s.HandleFunc("POST /v1/historical-stages/{id}/wage-statutes", h.CreateWageStatute)
	s.HandleFunc("POST /v1/historical-stages/{id}/vagrancy-laws", h.CreateVagrancyLaw)
	s.HandleFunc("GET /v1/historical-stages/{id}/labour-discipline", h.GetLabourDiscipline)
	s.HandleFunc("POST /v1/statutory-wages/compare", h.CompareStatutoryWage)

	// Ch. 29 — Genesis of the Capitalist Farmer
	s.HandleFunc("POST /v1/historical-stages/{id}/farm-tenures", h.CreateFarmTenure)
	s.HandleFunc("GET /v1/historical-stages/{id}/farm-tenures", h.ListFarmTenures)
	s.HandleFunc("POST /v1/farm-tenures/real-rent", h.ComputeRealRent)

	// Ch. 30 — Reaction of the Agricultural Revolution on Industry
	s.HandleFunc("POST /v1/historical-stages/{id}/domestic-industries", h.CreateDomesticIndustry)
	s.HandleFunc("GET /v1/historical-stages/{id}/home-market", h.GetHomeMarket)
	s.HandleFunc("POST /v1/market-formation", h.ComputeMarketFormation)

	// Ch. 31 — Genesis of the Industrial Capitalist
	s.HandleFunc("POST /v1/historical-stages/{id}/capital-origins", h.CreateCapitalOrigin)
	s.HandleFunc("POST /v1/historical-stages/{id}/colonial-transfers", h.CreateColonialTransfer)
	s.HandleFunc("POST /v1/historical-stages/{id}/national-debts", h.CreateNationalDebt)
	s.HandleFunc("GET /v1/historical-stages/{id}/genesis", h.GetIndustrialCapitalGenesis)

	// Ch. 32 — Historical Tendency of Capitalist Accumulation
	s.HandleFunc("POST /v1/accumulation/centralisation", h.RunCentralisation)
	s.HandleFunc("GET /v1/accumulation/negation-of-negation", h.GetNegationOfNegation)
	s.HandleFunc("GET /v1/accumulation/trajectories", h.ListAccumulationTrajectories)
	s.HandleFunc("GET /v1/accumulation/trajectories/{id}", h.GetAccumulationTrajectory)

	// Ch. 33 — The Modern Theory of Colonisation
	s.HandleFunc("POST /v1/colonial-markets", h.CreateColonialMarket)
	s.HandleFunc("GET /v1/colonial-markets", h.ListColonialMarkets)
	s.HandleFunc("GET /v1/colonial-markets/{id}", h.GetColonialMarket)
	s.HandleFunc("POST /v1/colonial-markets/{id}/regulate", h.RegulateColonialMarket)
	s.HandleFunc("POST /v1/colonial-markets/{id}/independence", h.ComputeIndependence)
}
