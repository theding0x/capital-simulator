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

	// Vol. III Ch. 2 — The Rate of Profit
	s.HandleFunc("POST /v1/profit/rate", h.CreateProfitRate)
	s.HandleFunc("GET /v1/profit/rate", h.ListProfitRates)
	s.HandleFunc("GET /v1/profit/rate/{id}", h.GetProfitRate)

	// Vol. III Ch. 3 — Relation of the Rate of Profit to the Rate of Surplus-Value
	s.HandleFunc("POST /v1/profit/variation", h.CreateVariation)
	s.HandleFunc("GET /v1/profit/variation", h.ListVariations)
	s.HandleFunc("GET /v1/profit/variation/{id}", h.GetVariation)
	s.HandleFunc("POST /v1/profit/compare", h.CompareProfitRates)

	// Vol. III Ch. 4 — The Effect of the Turnover on the Rate of Profit
	s.HandleFunc("POST /v1/profit/turnover-analysis", h.CreateTurnoverAnalysis)
	s.HandleFunc("GET /v1/profit/turnover-analysis", h.ListTurnoverAnalyses)
	s.HandleFunc("GET /v1/profit/turnover-analysis/{id}", h.GetTurnoverAnalysis)

	// Vol. III Ch. 5 — Economy in the Employment of Constant Capital
	s.HandleFunc("POST /v1/profit/economy", h.CreateEconomyAnalysis)
	s.HandleFunc("GET /v1/profit/economy", h.ListEconomyAnalyses)
	s.HandleFunc("GET /v1/profit/economy/{id}", h.GetEconomyAnalysis)

	// Vol. III Ch. 6 — The Effect of Price Fluctuation on the Rate of Profit
	s.HandleFunc("POST /v1/profit/price-fluctuation", h.CreatePriceFluctuationAnalysis)
	s.HandleFunc("GET /v1/profit/price-fluctuation", h.ListPriceFluctuationAnalyses)
	s.HandleFunc("GET /v1/profit/price-fluctuation/{id}", h.GetPriceFluctuationAnalysis)

	// Vol. III Ch. 7 — Supplementary Remarks (on the rate of profit)
	s.HandleFunc("POST /v1/profit/composition-effect", h.CreateCompositionEffect)
	s.HandleFunc("GET /v1/profit/composition-effect", h.ListCompositionEffects)
	s.HandleFunc("GET /v1/profit/composition-effect/{id}", h.GetCompositionEffect)
	s.HandleFunc("POST /v1/profit/magnitude-change", h.CreateMagnitudeChange)
	s.HandleFunc("GET /v1/profit/magnitude-change", h.ListMagnitudeChanges)
	s.HandleFunc("GET /v1/profit/magnitude-change/{id}", h.GetMagnitudeChange)
	s.HandleFunc("GET /v1/profit/summary", h.GetPartISummary)

	// Vol. III Ch. 8 — Different Compositions of Capitals in Different Branches
	s.HandleFunc("POST /v1/avgprofit/spheres", h.CreateProductionSphere)
	s.HandleFunc("GET /v1/avgprofit/spheres", h.ListProductionSpheres)
	s.HandleFunc("GET /v1/avgprofit/spheres/{id}", h.GetProductionSphere)

	// Vol. III Ch. 9 — Formation of a General Rate of Profit
	s.HandleFunc("POST /v1/avgprofit/general-rate", h.CreateGeneralProfitRate)
	s.HandleFunc("GET /v1/avgprofit/general-rate/{id}", h.GetGeneralProfitRate)
	s.HandleFunc("POST /v1/avgprofit/price-of-production", h.CreatePriceOfProduction)
	s.HandleFunc("GET /v1/avgprofit/price-of-production", h.ListPricesOfProduction)
	s.HandleFunc("GET /v1/avgprofit/price-of-production/{id}", h.GetPriceOfProduction)
	s.HandleFunc("POST /v1/avgprofit/social-aggregate", h.ComputeSocialAggregate)

	// Vol. III Ch. 10 — Equalisation of the General Rate of Profit Through Competition
	s.HandleFunc("POST /v1/avgprofit/market-value", h.CreateMarketValue)
	s.HandleFunc("GET /v1/avgprofit/market-value/{id}", h.GetMarketValue)
	s.HandleFunc("POST /v1/avgprofit/surplus-profit", h.CreateSurplusProfit)
	s.HandleFunc("GET /v1/avgprofit/surplus-profit/{id}", h.GetSurplusProfit)
	s.HandleFunc("POST /v1/avgprofit/capital-flow", h.CreateCapitalFlow)
	s.HandleFunc("GET /v1/avgprofit/equalisation/{id}", h.GetEqualisation)

	// Vol. III Ch. 11 — Effects of General Wage Fluctuations on Prices of Production
	s.HandleFunc("POST /v1/avgprofit/wage-effect", h.CreateWageEffectAnalysis)
	s.HandleFunc("GET /v1/avgprofit/wage-effect", h.ListWageEffectAnalyses)
	s.HandleFunc("GET /v1/avgprofit/wage-effect/{id}", h.GetWageEffectAnalysis)

	// Vol. III Ch. 12 — Supplementary Remarks (on prices of production)
	s.HandleFunc("POST /v1/avgprofit/price-change", h.CreatePriceOfProductionChange)
	s.HandleFunc("GET /v1/avgprofit/price-change/{id}", h.GetPriceOfProductionChange)
	s.HandleFunc("POST /v1/avgprofit/compensation-ground", h.ComputeCompensationGround)
	s.HandleFunc("GET /v1/avgprofit/summary", h.GetPartIISummary)

	// Vol. III Ch. 13 — The Law As Such (Tendential Fall in the Rate of Profit)
	s.HandleFunc("POST /v1/tendency/trajectory", h.CreateCompositionTrajectory)
	s.HandleFunc("GET /v1/tendency/trajectory", h.ListCompositionTrajectories)
	s.HandleFunc("GET /v1/tendency/trajectory/{id}", h.GetCompositionTrajectory)
	s.HandleFunc("POST /v1/tendency/rate-mass", h.CreateRateMassContradiction)
	s.HandleFunc("GET /v1/tendency/rate-mass/{id}", h.GetRateMassContradiction)

	// Vol. III Ch. 14 — Counteracting Influences
	s.HandleFunc("POST /v1/tendency/counteracting-force", h.CreateCounteractingForce)
	s.HandleFunc("GET /v1/tendency/counteracting-force/{id}", h.GetCounteractingForce)
	s.HandleFunc("POST /v1/tendency/scenario", h.CreateCounteractingScenario)
	s.HandleFunc("GET /v1/tendency/scenario", h.ListCounteractingScenarios)
	s.HandleFunc("GET /v1/tendency/scenario/{id}", h.GetCounteractingScenario)

	// Vol. III Ch. 15 — Internal Contradictions of the Law (COMPLETES Part III)
	s.HandleFunc("POST /v1/tendency/crisis", h.CreateCrisis)
	s.HandleFunc("GET /v1/tendency/crisis", h.ListCrises)
	s.HandleFunc("GET /v1/tendency/crisis/{id}", h.GetCrisis)
	s.HandleFunc("POST /v1/tendency/contradiction", h.CreateInternalContradiction)
	s.HandleFunc("GET /v1/tendency/contradiction/{id}", h.GetInternalContradiction)
	s.HandleFunc("GET /v1/tendency/summary", h.GetPartIIISummary)

	// Vol. III Ch. 16 — Commercial Capital (opens Part IV)
	s.HandleFunc("POST /v1/merchant/commercial-capital", h.CreateCommercialCapital)
	s.HandleFunc("GET /v1/merchant/commercial-capital", h.ListCommercialCapitals)
	s.HandleFunc("GET /v1/merchant/commercial-capital/{id}", h.GetCommercialCapital)

	// Vol. III Ch. 17 — Commercial Profit
	s.HandleFunc("POST /v1/merchant/commercial-profit", h.CreateCommercialProfit)
	s.HandleFunc("GET /v1/merchant/commercial-profit", h.ListCommercialProfits)
	s.HandleFunc("GET /v1/merchant/commercial-profit/{id}", h.GetCommercialProfit)

	// Vol. III Ch. 18 — Merchant Turnover
	s.HandleFunc("POST /v1/merchant/turnover", h.CreateMerchantTurnover)
	s.HandleFunc("GET /v1/merchant/turnover", h.ListMerchantTurnovers)
	s.HandleFunc("GET /v1/merchant/turnover/{id}", h.GetMerchantTurnover)
	s.HandleFunc("POST /v1/merchant/turnover-effect", h.ComputeTurnoverEffect)

	// Vol. III Ch. 19 — Money-Dealing Capital
	s.HandleFunc("POST /v1/merchant/money-dealing", h.CreateMoneyDealingCapital)
	s.HandleFunc("GET /v1/merchant/money-dealing", h.ListMoneyDealingCapitals)
	s.HandleFunc("GET /v1/merchant/money-dealing/{id}", h.GetMoneyDealingCapital)

	// Vol. III Ch. 20 — Historical Facts about Merchant's Capital (COMPLETES Part IV)
	s.HandleFunc("POST /v1/merchant/historical-forms", h.CreateHistoricalMerchantCapital)
	s.HandleFunc("GET /v1/merchant/historical-forms", h.ListHistoricalMerchantCapitals)
	s.HandleFunc("GET /v1/merchant/historical-forms/{id}", h.GetHistoricalMerchantCapital)

	// Vol. III Ch. 21 — Interest-Bearing Capital (opens Part V)
	s.HandleFunc("POST /v1/credit/interest-bearing-capital", h.CreateInterestBearingCapital)
	s.HandleFunc("GET /v1/credit/interest-bearing-capital", h.ListInterestBearingCapitals)
	s.HandleFunc("GET /v1/credit/interest-bearing-capital/{id}", h.GetInterestBearingCapital)

	// Vol. III Ch. 22 — Division of Profit; Rate of Interest
	s.HandleFunc("POST /v1/credit/rate-of-interest", h.CreateRateOfInterest)
	s.HandleFunc("GET /v1/credit/rate-of-interest", h.ListRatesOfInterest)
	s.HandleFunc("GET /v1/credit/rate-of-interest/{id}", h.GetRateOfInterest)
	s.HandleFunc("POST /v1/credit/interest-rate-analysis", h.ComputeInterestRateAnalysis)

	// Vol. III Ch. 23 — Interest and Profit of Enterprise
	s.HandleFunc("POST /v1/credit/profit-division", h.CreateProfitDivision)
	s.HandleFunc("GET /v1/credit/profit-division", h.ListProfitDivisions)
	s.HandleFunc("GET /v1/credit/profit-division/{id}", h.GetProfitDivision)

	// Vol. III Ch. 24 — Externalisation / Fetish Capital
	s.HandleFunc("POST /v1/credit/compound-interest", h.CreateCompoundInterestSchedule)
	s.HandleFunc("GET /v1/credit/compound-interest", h.ListCompoundInterestSchedules)
	s.HandleFunc("GET /v1/credit/compound-interest/{id}", h.GetCompoundInterestSchedule)

	// Vol. III Ch. 25 — Credit and Fictitious Capital
	s.HandleFunc("POST /v1/credit/bills-of-exchange", h.CreateBillOfExchange)
	s.HandleFunc("GET /v1/credit/bills-of-exchange", h.ListBillsOfExchange)
	s.HandleFunc("GET /v1/credit/bills-of-exchange/{id}", h.GetBillOfExchange)
	s.HandleFunc("POST /v1/credit/fictitious-capital", h.CreateFictitiousCapital)
	s.HandleFunc("GET /v1/credit/fictitious-capital", h.ListFictitiousCapitals)
	s.HandleFunc("GET /v1/credit/fictitious-capital/{id}", h.GetFictitiousCapital)

	// Vol. III Ch. 26 — Accumulation of Money-Capital
	s.HandleFunc("POST /v1/credit/money-capital-accumulation", h.CreateMoneyCapitalAccumulation)
	s.HandleFunc("GET /v1/credit/money-capital-accumulation", h.ListMoneyCapitalAccumulations)
	s.HandleFunc("GET /v1/credit/money-capital-accumulation/{id}", h.GetMoneyCapitalAccumulation)

	// Vol. III Ch. 27 — The Role of Credit in Capitalist Production
	s.HandleFunc("POST /v1/credit/stock-companies", h.CreateStockCompany)
	s.HandleFunc("GET /v1/credit/stock-companies", h.ListStockCompanies)
	s.HandleFunc("GET /v1/credit/stock-companies/{id}", h.GetStockCompany)
	s.HandleFunc("POST /v1/credit/cooperative-factories", h.CreateCooperativeFactory)
	s.HandleFunc("GET /v1/credit/cooperative-factories", h.ListCooperativeFactories)

	// Vol. III Ch. 28 — Medium of Circulation and Capital (Tooke and Fullarton)
	s.HandleFunc("POST /v1/credit/currency-observations", h.CreateCurrencyObservation)
	s.HandleFunc("GET /v1/credit/currency-observations", h.ListCurrencyObservations)
	s.HandleFunc("GET /v1/credit/currency-observations/{id}", h.GetCurrencyObservation)

	// Vol. III Ch. 29 — Component Parts of Bank Capital
	s.HandleFunc("POST /v1/credit/bank-capital", h.CreateBankCapital)
	s.HandleFunc("GET /v1/credit/bank-capital", h.ListBankCapitals)
	s.HandleFunc("GET /v1/credit/bank-capital/{id}", h.GetBankCapital)
	s.HandleFunc("POST /v1/credit/fictitious-capital-valuation", h.CreateFictitiousCapitalValuation)
	s.HandleFunc("GET /v1/credit/fictitious-capital-valuation", h.ListFictitiousCapitalValuations)
	s.HandleFunc("GET /v1/credit/fictitious-capital-valuation/{id}", h.GetFictitiousCapitalValuation)

	// Vol. III Ch. 30 — Money-Capital and Real Capital, I
	s.HandleFunc("POST /v1/credit/real-capital-accumulation", h.CreateRealCapitalAccumulation)
	s.HandleFunc("GET /v1/credit/real-capital-accumulation", h.ListRealCapitalAccumulations)
	s.HandleFunc("GET /v1/credit/real-capital-accumulation/{id}", h.GetRealCapitalAccumulation)

	// Vol. III Ch. 31 — Money-Capital and Real Capital, II
	s.HandleFunc("POST /v1/credit/floating-capital", h.CreateFloatingCapital)
	s.HandleFunc("GET /v1/credit/floating-capital", h.ListFloatingCapitals)
	s.HandleFunc("GET /v1/credit/floating-capital/{id}", h.GetFloatingCapital)

	// Vol. III Ch. 32 — Money-Capital and Real Capital, III (Conclusion)
	s.HandleFunc("POST /v1/credit/capital-release", h.CreateCapitalRelease)
	s.HandleFunc("GET /v1/credit/capital-release", h.ListCapitalReleases)
	s.HandleFunc("GET /v1/credit/capital-release/{id}", h.GetCapitalRelease)

	// Vol. III Ch. 33 — The Medium of Circulation in the Credit System
	s.HandleFunc("POST /v1/credit/clearing-house-settlements", h.CreateClearingHouseSettlement)
	s.HandleFunc("GET /v1/credit/clearing-house-settlements", h.ListClearingHouseSettlements)
	s.HandleFunc("GET /v1/credit/clearing-house-settlements/{id}", h.GetClearingHouseSettlement)

	// Vol. III Ch. 34 — The Currency Principle and the Bank Acts of 1844 and 1845
	s.HandleFunc("POST /v1/credit/note-issue-constraints", h.CreateNoteIssueConstraint)
	s.HandleFunc("GET /v1/credit/note-issue-constraints", h.ListNoteIssueConstraints)
	s.HandleFunc("GET /v1/credit/note-issue-constraints/{id}", h.GetNoteIssueConstraint)

	// Vol. III Ch. 35 — Precious Metal and Rate of Exchange
	s.HandleFunc("POST /v1/credit/gold-reserves", h.CreateGoldReserve)
	s.HandleFunc("GET /v1/credit/gold-reserves", h.ListGoldReserves)
	s.HandleFunc("GET /v1/credit/gold-reserves/{id}", h.GetGoldReserve)
	s.HandleFunc("POST /v1/credit/rates-of-exchange", h.CreateRateOfExchange)
	s.HandleFunc("GET /v1/credit/rates-of-exchange", h.ListRatesOfExchange)
	s.HandleFunc("GET /v1/credit/rates-of-exchange/{id}", h.GetRateOfExchange)

	// Vol. III Ch. 36 — Pre-Capitalist Relationships (FINAL chapter)
	s.HandleFunc("POST /v1/credit/usurers-capital", h.CreateUsurersCapital)
	s.HandleFunc("GET /v1/credit/usurers-capital", h.ListUsurersCapitals)
	s.HandleFunc("GET /v1/credit/usurers-capital/{id}", h.GetUsurersCapital)

	// Vol. III Ch. 37 — Introduction to Ground-Rent (opens Part VI)
	s.HandleFunc("POST /v1/rent/parcels", h.CreateLandParcel)
	s.HandleFunc("GET /v1/rent/parcels", h.ListLandParcels)
	s.HandleFunc("GET /v1/rent/parcels/{id}", h.GetLandParcel)
	s.HandleFunc("POST /v1/rent/rents", h.CreateGroundRent)
	s.HandleFunc("GET /v1/rent/rents", h.ListGroundRents)

	// Vol. III Ch. 38 — Differential Rent: General Remarks
	s.HandleFunc("POST /v1/rent/surplus-profits", h.CreatePoPSurplusProfit)
	s.HandleFunc("GET /v1/rent/surplus-profits", h.ListPoPSurplusProfits)
	s.HandleFunc("POST /v1/rent/monopolised-forces", h.CreateMonopolisedNaturalForce)
	s.HandleFunc("GET /v1/rent/monopolised-forces", h.ListMonopolisedNaturalForces)
	s.HandleFunc("POST /v1/rent/differential-rents", h.CreateDifferentialRent)
	s.HandleFunc("GET /v1/rent/differential-rents", h.ListDifferentialRents)
	s.HandleFunc("POST /v1/rent/capitalised-prices", h.CreateCapitalisedRentPrice)
	s.HandleFunc("GET /v1/rent/capitalised-prices", h.ListCapitalisedRentPrices)

	// Vol. III Ch. 39 — Differential Rent I (First Form)
	s.HandleFunc("POST /v1/rent/dr1-tables", h.CreateDifferentialRentITable)
	s.HandleFunc("GET /v1/rent/dr1-tables", h.ListDifferentialRentITables)
	s.HandleFunc("GET /v1/rent/dr1-tables/{id}", h.GetDifferentialRentITable)
	s.HandleFunc("GET /v1/rent/dr1-tables/{id}/entries", h.ListDifferentialRentIEntries)
	s.HandleFunc("POST /v1/rent/location-rent-factors", h.CreateLocationRentFactor)
	s.HandleFunc("GET /v1/rent/location-rent-factors", h.ListLocationRentFactors)

	// Vol. III Ch. 40 — Differential Rent II (Second Form)
	s.HandleFunc("POST /v1/rent/successive-investments", h.CreateSuccessiveInvestment)
	s.HandleFunc("GET /v1/rent/successive-investments", h.ListSuccessiveInvestments)
	s.HandleFunc("GET /v1/rent/dr2-tables/{parcel_id}", h.GetDR2Table)
	s.HandleFunc("POST /v1/rent/lease-terms", h.CreateLeaseTerm)
	s.HandleFunc("GET /v1/rent/lease-terms", h.ListLeaseTerms)

	// Vol. III Ch. 41 — DR II, First Case: Constant Price of Production
	s.HandleFunc("POST /v1/rent/dr2-outcomes", h.CreateDR2Outcome)
	s.HandleFunc("GET /v1/rent/dr2-outcomes", h.ListDR2Outcomes)
	s.HandleFunc("POST /v1/rent/intensive-extensive", h.CreateIntensiveExtensiveComparison)
	s.HandleFunc("GET /v1/rent/intensive-extensive", h.ListIntensiveExtensiveComparisons)

	// Vol. III Ch. 42 — DR II, Second Case: Falling Price of Production
	s.HandleFunc("POST /v1/rent/exclusion-events", h.CreateSoilExclusionEvent)
	s.HandleFunc("GET /v1/rent/exclusion-events", h.ListSoilExclusionEvents)
	s.HandleFunc("POST /v1/rent/dr2-falling-outcomes", h.CreateFallingPriceDR2Outcome)
	s.HandleFunc("GET /v1/rent/dr2-falling-outcomes", h.ListFallingPriceDR2Outcomes)
	s.HandleFunc("POST /v1/rent/normal-capital", h.CreateNormalCapitalPerAcre)
	s.HandleFunc("GET /v1/rent/normal-capital", h.ListNormalCapitalPerAcre)

	// Vol. III Ch. 43 — DR II, Third Case: Rising Price of Production
	s.HandleFunc("POST /v1/rent/rising-price-outcomes", h.CreateRisingPriceDR2Outcome)
	s.HandleFunc("GET /v1/rent/rising-price-outcomes", h.ListRisingPriceDR2Outcomes)
	s.HandleFunc("POST /v1/rent/rent-sequences", h.CreateRentSequence)
	s.HandleFunc("GET /v1/rent/rent-sequences", h.ListRentSequences)
	s.HandleFunc("POST /v1/rent/engels-tables", h.CreateEngelsTableEntry)
	s.HandleFunc("GET /v1/rent/engels-tables", h.ListEngelsTableEntries)

	// Vol. III Ch. 44 — Differential Rent on the Worst Cultivated Soil
	s.HandleFunc("POST /v1/rent/worst-soil-rents", h.CreateRentOnWorstSoil)
	s.HandleFunc("GET /v1/rent/worst-soil-rents", h.ListRentOnWorstSoils)
	s.HandleFunc("POST /v1/rent/soil-improvements", h.CreateSoilImprovementCapital)
	s.HandleFunc("GET /v1/rent/soil-improvements", h.ListSoilImprovementCapitals)
	s.HandleFunc("POST /v1/rent/regulating-price-shifts", h.CreateRegulatingPriceShift)
	s.HandleFunc("GET /v1/rent/regulating-price-shifts", h.ListRegulatingPriceShifts)

	// Vol. III Ch. 45 — Absolute Ground-Rent
	s.HandleFunc("POST /v1/rent/compositions", h.CreateAgriculturalComposition)
	s.HandleFunc("GET /v1/rent/compositions", h.ListAgriculturalCompositions)
	s.HandleFunc("POST /v1/rent/value-gaps", h.CreateValuePriceGap)
	s.HandleFunc("GET /v1/rent/value-gaps", h.ListValuePriceGaps)
	s.HandleFunc("POST /v1/rent/absolute-rents", h.CreateAbsoluteRent)
	s.HandleFunc("GET /v1/rent/absolute-rents", h.ListAbsoluteRents)
	s.HandleFunc("GET /v1/rent/absolute-rents/{id}", h.GetAbsoluteRent)
	s.HandleFunc("POST /v1/rent/rent-limits", h.CreateAbsoluteRentLimit)
	s.HandleFunc("GET /v1/rent/rent-limits", h.ListAbsoluteRentLimits)

	// Vol. III Ch. 46 — Building-Site Rent, Mining Rent, Monopoly-Price Rent, Price of Land
	s.HandleFunc("POST /v1/rent/building-rents", h.CreateBuildingSiteRent)
	s.HandleFunc("GET /v1/rent/building-rents", h.ListBuildingSiteRents)
	s.HandleFunc("POST /v1/rent/mining-rents", h.CreateMiningRent)
	s.HandleFunc("GET /v1/rent/mining-rents", h.ListMiningRents)
	s.HandleFunc("POST /v1/rent/monopoly-rents", h.CreateMonopolyPriceRent)
	s.HandleFunc("GET /v1/rent/monopoly-rents", h.ListMonopolyPriceRents)
	s.HandleFunc("POST /v1/rent/land-price-scenarios", h.CreateLandPriceScenario)
	s.HandleFunc("GET /v1/rent/land-price-scenarios", h.ListLandPriceScenarios)
	s.HandleFunc("POST /v1/rent/historical-forms", h.CreateRentHistoricalForm)
	s.HandleFunc("GET /v1/rent/historical-forms", h.ListRentHistoricalForms)
	s.HandleFunc("POST /v1/rent/historical-transitions", h.CreateHistoricalRentTransition)
	s.HandleFunc("GET /v1/rent/historical-transitions", h.ListHistoricalRentTransitions)
	s.HandleFunc("POST /v1/rent/small-peasant-productions", h.CreateSmallPeasantProduction)
	s.HandleFunc("GET /v1/rent/small-peasant-productions", h.ListSmallPeasantProductions)

	// Vol. III Ch. 48 — The Trinity Formula (Part VII opens)
	s.HandleFunc("POST /v1/revenue/trinity-formulas", h.CreateTrinityFormula)
	s.HandleFunc("GET /v1/revenue/trinity-formulas/{id}", h.GetTrinityFormula)
	s.HandleFunc("GET /v1/revenue/fetish-forms", h.ListRevenueFetishForms)

	// Vol. III Ch. 49 — Concerning the Analysis of the Process of Production
	s.HandleFunc("POST /v1/revenue/annual-compositions", h.CreateAnnualValueComposition)
	s.HandleFunc("GET /v1/revenue/annual-compositions/{id}", h.GetAnnualValueComposition)
	s.HandleFunc("POST /v1/revenue/surplus-partitions", h.CreateSurplusValuePartition)

	// Vol. III Ch. 50 — Illusions Created by Competition
	s.HandleFunc("GET /v1/revenue/competition-illusions", h.ListCompetitionIllusions)
	s.HandleFunc("POST /v1/revenue/cost-price-fetishes", h.CreateCostPriceFetish)
	s.HandleFunc("POST /v1/revenue/annual-accounts", h.CreateAnnualRevenueAccount)
	s.HandleFunc("GET /v1/revenue/value-contrasts", h.ListValueAppearanceContrasts)
}
