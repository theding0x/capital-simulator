package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

func Register(s *httpx.Server, h *Handler) {
	// Vol. II Ch. 1 — The Circuit of Money-Capital
	s.HandleFunc("POST /v1/money-circuits", h.CreateMoneyCircuit)
	s.HandleFunc("GET /v1/money-circuits", h.ListMoneyCircuits)
	s.HandleFunc("GET /v1/money-circuits/{id}", h.GetMoneyCircuit)
	s.HandleFunc("POST /v1/money-circuits/{id}/purchase", h.RecordMoneyPurchase)
	s.HandleFunc("POST /v1/money-circuits/{id}/produce", h.RecordMoneyProduce)
	s.HandleFunc("POST /v1/money-circuits/{id}/commodity", h.RecordMoneyCommodity)
	s.HandleFunc("POST /v1/money-circuits/{id}/realise", h.RecordMoneyRealise)

	// Vol. II Ch. 3 — The Circuit of Commodity-Capital
	s.HandleFunc("POST /v1/commodity-circuits", h.CreateCommodityCircuit)
	s.HandleFunc("GET /v1/commodity-circuits", h.ListCommodityCircuits)
	s.HandleFunc("GET /v1/commodity-circuits/aggregate", h.GetCommodityCircuitAggregate)
	s.HandleFunc("GET /v1/commodity-circuits/{id}", h.GetCommodityCircuit)
	s.HandleFunc("POST /v1/commodity-circuits/{id}/partial-sales", h.RecordCommodityPartialSale)
	s.HandleFunc("POST /v1/commodity-circuits/{id}/mp-source", h.LinkCommodityMPSource)
	s.HandleFunc("POST /v1/commodity-circuits/{id}/close", h.CloseCommodityCircuit)

	// Vol. II Ch. 7 — The Turnover Time and the Number of Turnovers
	s.HandleFunc("POST /v1/turnovers", h.CreateTurnover)
	s.HandleFunc("GET /v1/turnovers", h.ListTurnovers)
	s.HandleFunc("GET /v1/turnovers/{id}", h.GetTurnover)
	s.HandleFunc("POST /v1/turnovers/{id}/cycles", h.RecordCycle)
	s.HandleFunc("POST /v1/turnovers/{id}/recompute-number", h.RecomputeNumber)
	s.HandleFunc("GET /v1/turnovers/{id}/number", h.GetNumber)

	// Vol. II Ch. 4 — The Three Formulas of the Circuit
	s.HandleFunc("POST /v1/industrial-capitals", h.CreateIndustrialCapital)
	s.HandleFunc("GET /v1/industrial-capitals", h.ListIndustrialCapitals)
	s.HandleFunc("GET /v1/industrial-capitals/supply-demand/aggregate", h.AggregateSupplyDemand)
	s.HandleFunc("GET /v1/industrial-capitals/{id}", h.GetIndustrialCapital)
	s.HandleFunc("POST /v1/industrial-capitals/{id}/parts", h.RecordCapitalPart)
	s.HandleFunc("POST /v1/industrial-capitals/{id}/snapshot", h.SnapshotStageDistribution)
	s.HandleFunc("POST /v1/industrial-capitals/{id}/blocks", h.OpenStageBlock)
	s.HandleFunc("POST /v1/industrial-capitals/{id}/blocks/{block_id}/close", h.CloseStageBlock)
	s.HandleFunc("POST /v1/industrial-capitals/{id}/value-revolution", h.RecordValueRevolution)
	s.HandleFunc("POST /v1/industrial-capitals/{id}/interlocks", h.RecordInterlock)
	s.HandleFunc("POST /v1/industrial-capitals/{id}/supply-demand", h.ComputeAndRecordSupplyDemand)
	s.HandleFunc("GET /v1/industrial-capitals/{id}/supply-demand", h.GetSupplyDemand)
	s.HandleFunc("POST /v1/industrial-capitals/{id}/sinking-fund", h.SetSinkingFund)
	s.HandleFunc("POST /v1/industrial-capitals/{id}/sinking-fund/tick", h.TickSinkingFund)

	// Vol. II Ch. 2 — The Circuit of Productive Capital
	s.HandleFunc("POST /v1/productive-circuits", h.CreateProductiveCircuit)
	s.HandleFunc("GET /v1/productive-circuits", h.ListProductiveCircuits)
	s.HandleFunc("GET /v1/productive-circuits/{id}", h.GetProductiveCircuit)
	s.HandleFunc("POST /v1/productive-circuits/{id}/revenue", h.RecordRevenue)
	s.HandleFunc("POST /v1/productive-circuits/{id}/accumulate", h.AccumulateLatentCapital)
	s.HandleFunc("POST /v1/productive-circuits/{id}/capitalise", h.CapitaliseCircuit)
	s.HandleFunc("POST /v1/productive-circuits/{id}/reserve/deposit", h.DepositReserve)
	s.HandleFunc("POST /v1/productive-circuits/{id}/reserve/draw", h.DrawReserve)

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

	// Vol. II Ch. 9 — The Aggregate Turnover of Advanced Capital. Cycles of Turnover.
	s.HandleFunc("POST /v1/aggregate-turnovers", h.CreateAggregateTurnover)
	s.HandleFunc("GET /v1/aggregate-turnovers", h.ListAggregateTurnovers)
	s.HandleFunc("GET /v1/aggregate-turnovers/{id}", h.GetAggregateTurnover)
	s.HandleFunc("POST /v1/aggregate-turnovers/{id}/recompute", h.RecomputeAggregateTurnover)
	s.HandleFunc("POST /v1/aggregate-turnovers/{id}/lifetime-cycle", h.RecordAggregateTurnoverLifetimeCycle)
	s.HandleFunc("POST /v1/aggregate-turnovers/{id}/crisis-phase", h.RecordCrisisPhase)

	// Vol. II Ch. 8 — Fixed Capital and Circulating Capital
	s.HandleFunc("POST /v1/capital-components", h.CreateComponent)
	s.HandleFunc("GET /v1/capital-components", h.ListComponents)
	s.HandleFunc("GET /v1/capital-components/{id}", h.GetComponent)
	s.HandleFunc("POST /v1/fixed-items", h.RegisterFixedItem)
	s.HandleFunc("GET /v1/fixed-items/{id}", h.GetFixedItem)
	s.HandleFunc("POST /v1/fixed-items/{id}/subcomponents", h.RecordSubcomponent)
	s.HandleFunc("POST /v1/fixed-items/{id}/wear", h.RecordWear)
	s.HandleFunc("GET /v1/fixed-items/{id}/sinking-fund", h.GetSinkingFund)
	s.HandleFunc("POST /v1/fixed-items/{id}/repairs", h.RecordRepair)
	s.HandleFunc("POST /v1/fixed-items/{id}/reinvestments", h.RecordReinvestment)
	s.HandleFunc("POST /v1/circulating-cycles", h.RecordCirculatingCycle)
}
