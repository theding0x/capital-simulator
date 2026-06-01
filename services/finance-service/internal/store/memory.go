package store

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/tendency"
)

// Memory is an in-memory Store for unit tests and local development.
type Memory struct {
	mu                         sync.RWMutex
	now                        func() time.Time
	costPrices                 map[profit.CostPriceID]profit.CostPrice
	profitRates                map[profit.ProfitRateID]profit.ProfitRateAnalysis
	variations                 map[profit.VariationAnalysisID]profit.VariationAnalysis
	turnovers                  map[profit.TurnoverAnalysisID]profit.TurnoverAnalysis
	economies                  map[profit.EconomyAnalysisID]profit.EconomyAnalysis
	priceFlux                  map[profit.PriceFluctuationAnalysisID]profit.PriceFluctuationAnalysis
	compositions               map[profit.CompositionEffectID]profit.CompositionEffectAnalysis
	magnitudes                 map[profit.MagnitudeChangeID]profit.MagnitudeChangeAnalysis
	spheres                    map[avgprofit.ProductionSphereID]avgprofit.ProductionSphere
	generalRates               map[avgprofit.GeneralProfitRateID]avgprofit.GeneralProfitRate
	pricesOfProduction         map[avgprofit.PriceOfProductionID]avgprofit.PriceOfProduction
	marketValues               map[avgprofit.MarketValueID]avgprofit.MarketValue
	surplusProfits             map[avgprofit.SurplusProfitID]avgprofit.SurplusProfit
	equalisations              map[avgprofit.EqualisationID]avgprofit.Equalisation
	wageEffects                map[avgprofit.WageEffectAnalysisID]avgprofit.WageEffectAnalysis
	priceChanges               map[avgprofit.PriceOfProductionChangeID]avgprofit.PriceOfProductionChange
	trajectories               map[tendency.CompositionTrajectoryID]tendency.CompositionTrajectory
	rateMasses                 map[tendency.RateMassContradictionID]tendency.RateMassContradiction
	counterForces              map[tendency.CounteractingForceID]tendency.CounteractingForce
	counterScenarios           map[tendency.CounteractingScenarioID]tendency.CounteractingScenario
	crises                     map[tendency.CrisisID]tendency.Crisis
	contradictions             map[tendency.InternalContradictionID]tendency.InternalContradiction
	commercialCapitals         map[merchant.CommercialCapitalID]merchant.CommercialCapital
	commercialProfits          map[merchant.CommercialProfitID]merchant.CommercialProfit
	turnoversM                 map[merchant.MerchantTurnoverID]merchant.MerchantTurnover
	moneyDealingCapitals       map[merchant.MoneyDealingCapitalID]merchant.MoneyDealingCapital
	historicalMerchantCapitals map[merchant.HistoricalMerchantCapitalID]merchant.HistoricalMerchantCapital
	interestBearingCapitals    map[credit.InterestBearingCapitalID]credit.InterestBearingCapital
	ratesOfInterest            map[credit.RateOfInterestID]credit.RateOfInterest
	profitDivisions            map[credit.ProfitDivisionID]credit.ProfitDivision
	compoundInterestSchedules  map[credit.CompoundInterestID]credit.CompoundInterestSchedule
	billsOfExchange            map[credit.BillOfExchangeID]credit.BillOfExchange
	fictitiousCapitals         map[credit.FictitiousCapitalID]credit.FictitiousCapital
	moneyCapitalAccumulations  map[credit.MoneyCapitalAccumulationID]credit.MoneyCapitalAccumulation
	stockCompanies             map[credit.StockCompanyID]credit.StockCompany
	cooperativeFactories       map[credit.CooperativeFactoryID]credit.CooperativeFactory
	currencyObservations       map[credit.CurrencyObservationID]credit.CurrencyObservation
	bankCapitals               map[credit.BankCapitalID]credit.BankCapital
	fictitiousValuations       map[credit.FictitiousCapitalValuationID]credit.FictitiousCapitalValuation
	realCapitalAccumulations   map[credit.RealCapitalAccumulationID]credit.RealCapitalAccumulation
	floatingCapitals           map[credit.FloatingCapitalID]credit.FloatingCapital
	capitalReleases            map[credit.CapitalReleaseID]credit.CapitalRelease
	clearingHouseSettlements   map[credit.ClearingHouseSettlementID]credit.ClearingHouseSettlement
	noteIssueConstraints       map[credit.NoteIssueConstraintID]credit.NoteIssueConstraint
	goldReserves               map[credit.GoldReserveID]credit.GoldReserve
	ratesOfExchange            map[credit.RateOfExchangeID]credit.RateOfExchange
	usurersCapitals            map[credit.UsurersCapitalID]credit.UsurersCapital
	landParcels                map[rent.LandParcelID]rent.LandParcel
	landowners                 map[rent.LandownerID]rent.Landowner
	agriculturalCapitalists    map[rent.AgriculturalCapitalistID]rent.AgriculturalCapitalist
	groundRents                map[rent.GroundRentID]rent.GroundRent
	popSurplusProfits          map[rent.PriceOfProductionSurplusProfitID]rent.PriceOfProductionSurplusProfit
	monopolisedForces          map[rent.MonopolisedNaturalForceID]rent.MonopolisedNaturalForce
	differentialRents          map[rent.DifferentialRentID]rent.DifferentialRent
	capitalisedRentPrices      map[rent.CapitalisedRentPriceID]rent.CapitalisedRentPrice
	dr1Tables                  map[rent.DifferentialRentITableID]rent.DifferentialRentITable
	dr1Entries                 map[rent.DifferentialRentITableID][]rent.DifferentialRentIEntry
	locationRentFactors        map[rent.LocationRentFactorID]rent.LocationRentFactor
	successiveInvestments      map[rent.SuccessiveInvestmentID]rent.SuccessiveInvestment
	leaseTerms                 map[rent.LeaseTermID]rent.LeaseTerm
	dr2Outcomes                map[rent.DifferentialRentIIOutcomeID]rent.DifferentialRentIIOutcome
	intensiveExtensive         map[rent.IntensiveExtensiveComparisonID]rent.IntensiveExtensiveComparison
	soilExclusionEvents        map[rent.SoilExclusionEventID]rent.SoilExclusionEvent
	fallingPriceDR2Outcomes    map[rent.FallingPriceDR2OutcomeID]rent.FallingPriceDR2Outcome
	normalCapitalPerAcre       map[rent.NormalCapitalPerAcreID]rent.NormalCapitalPerAcre
	risingPriceDR2Outcomes     map[rent.RisingPriceDR2OutcomeID]rent.RisingPriceDR2Outcome
	rentSequences              map[rent.RentSequenceID]rent.RentSequence
	engelsTableEntries         map[rent.EngelsTableEntryID]rent.EngelsTableEntry
	rentOnWorstSoils           map[rent.RentOnWorstSoilID]rent.RentOnWorstSoil
	soilImprovementCapitals    map[rent.SoilImprovementCapitalID]rent.SoilImprovementCapital
	regulatingPriceShifts      map[rent.RegulatingPriceShiftID]rent.RegulatingPriceShift
	agriculturalCompositions   map[rent.AgriculturalCompositionID]rent.AgriculturalComposition
	valuePriceGaps             map[rent.ValuePriceGapID]rent.ValuePriceGap
	absoluteRents              map[rent.AbsoluteRentID]rent.AbsoluteRent
	absoluteRentLimits         map[rent.AbsoluteRentLimitID]rent.AbsoluteRentLimit
	buildingSiteRents          map[rent.BuildingSiteRentID]rent.BuildingSiteRent
	miningRents                map[rent.MiningRentID]rent.MiningRent
	monopolyPriceRents         map[rent.MonopolyPriceRentID]rent.MonopolyPriceRent
	landPriceScenarios         map[rent.LandPriceScenarioID]rent.LandPriceScenario
	rentHistoricalForms        map[rent.RentHistoricalFormID]rent.RentHistoricalForm
	historicalRentTransitions  map[rent.HistoricalRentTransitionID]rent.HistoricalRentTransition
	smallPeasantProductions    map[rent.SmallPeasantProductionID]rent.SmallPeasantProduction
	revenueStreams             map[revenue.RevenueStreamID]revenue.RevenueStream
	trinityFormulas            map[revenue.TrinityFormulaID]revenue.TrinityFormula
	revenueFetishForms         map[revenue.RevenueFetishFormID]revenue.RevenueFetishForm
}

// NewMemory returns an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{
		now:                        time.Now,
		costPrices:                 make(map[profit.CostPriceID]profit.CostPrice),
		profitRates:                make(map[profit.ProfitRateID]profit.ProfitRateAnalysis),
		variations:                 make(map[profit.VariationAnalysisID]profit.VariationAnalysis),
		turnovers:                  make(map[profit.TurnoverAnalysisID]profit.TurnoverAnalysis),
		economies:                  make(map[profit.EconomyAnalysisID]profit.EconomyAnalysis),
		priceFlux:                  make(map[profit.PriceFluctuationAnalysisID]profit.PriceFluctuationAnalysis),
		compositions:               make(map[profit.CompositionEffectID]profit.CompositionEffectAnalysis),
		magnitudes:                 make(map[profit.MagnitudeChangeID]profit.MagnitudeChangeAnalysis),
		spheres:                    make(map[avgprofit.ProductionSphereID]avgprofit.ProductionSphere),
		generalRates:               make(map[avgprofit.GeneralProfitRateID]avgprofit.GeneralProfitRate),
		pricesOfProduction:         make(map[avgprofit.PriceOfProductionID]avgprofit.PriceOfProduction),
		marketValues:               make(map[avgprofit.MarketValueID]avgprofit.MarketValue),
		surplusProfits:             make(map[avgprofit.SurplusProfitID]avgprofit.SurplusProfit),
		equalisations:              make(map[avgprofit.EqualisationID]avgprofit.Equalisation),
		wageEffects:                make(map[avgprofit.WageEffectAnalysisID]avgprofit.WageEffectAnalysis),
		priceChanges:               make(map[avgprofit.PriceOfProductionChangeID]avgprofit.PriceOfProductionChange),
		trajectories:               make(map[tendency.CompositionTrajectoryID]tendency.CompositionTrajectory),
		rateMasses:                 make(map[tendency.RateMassContradictionID]tendency.RateMassContradiction),
		counterForces:              make(map[tendency.CounteractingForceID]tendency.CounteractingForce),
		counterScenarios:           make(map[tendency.CounteractingScenarioID]tendency.CounteractingScenario),
		crises:                     make(map[tendency.CrisisID]tendency.Crisis),
		contradictions:             make(map[tendency.InternalContradictionID]tendency.InternalContradiction),
		commercialCapitals:         make(map[merchant.CommercialCapitalID]merchant.CommercialCapital),
		commercialProfits:          make(map[merchant.CommercialProfitID]merchant.CommercialProfit),
		turnoversM:                 make(map[merchant.MerchantTurnoverID]merchant.MerchantTurnover),
		moneyDealingCapitals:       make(map[merchant.MoneyDealingCapitalID]merchant.MoneyDealingCapital),
		historicalMerchantCapitals: make(map[merchant.HistoricalMerchantCapitalID]merchant.HistoricalMerchantCapital),
		interestBearingCapitals:    make(map[credit.InterestBearingCapitalID]credit.InterestBearingCapital),
		ratesOfInterest:            make(map[credit.RateOfInterestID]credit.RateOfInterest),
		profitDivisions:            make(map[credit.ProfitDivisionID]credit.ProfitDivision),
		compoundInterestSchedules:  make(map[credit.CompoundInterestID]credit.CompoundInterestSchedule),
		billsOfExchange:            make(map[credit.BillOfExchangeID]credit.BillOfExchange),
		fictitiousCapitals:         make(map[credit.FictitiousCapitalID]credit.FictitiousCapital),
		moneyCapitalAccumulations:  make(map[credit.MoneyCapitalAccumulationID]credit.MoneyCapitalAccumulation),
		stockCompanies:             make(map[credit.StockCompanyID]credit.StockCompany),
		cooperativeFactories:       make(map[credit.CooperativeFactoryID]credit.CooperativeFactory),
		currencyObservations:       make(map[credit.CurrencyObservationID]credit.CurrencyObservation),
		bankCapitals:               make(map[credit.BankCapitalID]credit.BankCapital),
		fictitiousValuations:       make(map[credit.FictitiousCapitalValuationID]credit.FictitiousCapitalValuation),
		realCapitalAccumulations:   make(map[credit.RealCapitalAccumulationID]credit.RealCapitalAccumulation),
		floatingCapitals:           make(map[credit.FloatingCapitalID]credit.FloatingCapital),
		capitalReleases:            make(map[credit.CapitalReleaseID]credit.CapitalRelease),
		clearingHouseSettlements:   make(map[credit.ClearingHouseSettlementID]credit.ClearingHouseSettlement),
		noteIssueConstraints:       make(map[credit.NoteIssueConstraintID]credit.NoteIssueConstraint),
		goldReserves:               make(map[credit.GoldReserveID]credit.GoldReserve),
		ratesOfExchange:            make(map[credit.RateOfExchangeID]credit.RateOfExchange),
		usurersCapitals:            make(map[credit.UsurersCapitalID]credit.UsurersCapital),
		landParcels:                make(map[rent.LandParcelID]rent.LandParcel),
		landowners:                 make(map[rent.LandownerID]rent.Landowner),
		agriculturalCapitalists:    make(map[rent.AgriculturalCapitalistID]rent.AgriculturalCapitalist),
		groundRents:                make(map[rent.GroundRentID]rent.GroundRent),
		popSurplusProfits:          make(map[rent.PriceOfProductionSurplusProfitID]rent.PriceOfProductionSurplusProfit),
		monopolisedForces:          make(map[rent.MonopolisedNaturalForceID]rent.MonopolisedNaturalForce),
		differentialRents:          make(map[rent.DifferentialRentID]rent.DifferentialRent),
		capitalisedRentPrices:      make(map[rent.CapitalisedRentPriceID]rent.CapitalisedRentPrice),
		dr1Tables:                  make(map[rent.DifferentialRentITableID]rent.DifferentialRentITable),
		dr1Entries:                 make(map[rent.DifferentialRentITableID][]rent.DifferentialRentIEntry),
		locationRentFactors:        make(map[rent.LocationRentFactorID]rent.LocationRentFactor),
		successiveInvestments:      make(map[rent.SuccessiveInvestmentID]rent.SuccessiveInvestment),
		leaseTerms:                 make(map[rent.LeaseTermID]rent.LeaseTerm),
		dr2Outcomes:                make(map[rent.DifferentialRentIIOutcomeID]rent.DifferentialRentIIOutcome),
		intensiveExtensive:         make(map[rent.IntensiveExtensiveComparisonID]rent.IntensiveExtensiveComparison),
		soilExclusionEvents:        make(map[rent.SoilExclusionEventID]rent.SoilExclusionEvent),
		fallingPriceDR2Outcomes:    make(map[rent.FallingPriceDR2OutcomeID]rent.FallingPriceDR2Outcome),
		normalCapitalPerAcre:       make(map[rent.NormalCapitalPerAcreID]rent.NormalCapitalPerAcre),
		risingPriceDR2Outcomes:     make(map[rent.RisingPriceDR2OutcomeID]rent.RisingPriceDR2Outcome),
		rentSequences:              make(map[rent.RentSequenceID]rent.RentSequence),
		engelsTableEntries:         make(map[rent.EngelsTableEntryID]rent.EngelsTableEntry),
		rentOnWorstSoils:           make(map[rent.RentOnWorstSoilID]rent.RentOnWorstSoil),
		soilImprovementCapitals:    make(map[rent.SoilImprovementCapitalID]rent.SoilImprovementCapital),
		regulatingPriceShifts:      make(map[rent.RegulatingPriceShiftID]rent.RegulatingPriceShift),
		agriculturalCompositions:   make(map[rent.AgriculturalCompositionID]rent.AgriculturalComposition),
		valuePriceGaps:             make(map[rent.ValuePriceGapID]rent.ValuePriceGap),
		absoluteRents:              make(map[rent.AbsoluteRentID]rent.AbsoluteRent),
		absoluteRentLimits:         make(map[rent.AbsoluteRentLimitID]rent.AbsoluteRentLimit),
		buildingSiteRents:          make(map[rent.BuildingSiteRentID]rent.BuildingSiteRent),
		miningRents:                make(map[rent.MiningRentID]rent.MiningRent),
		monopolyPriceRents:         make(map[rent.MonopolyPriceRentID]rent.MonopolyPriceRent),
		landPriceScenarios:         make(map[rent.LandPriceScenarioID]rent.LandPriceScenario),
		rentHistoricalForms:        make(map[rent.RentHistoricalFormID]rent.RentHistoricalForm),
		historicalRentTransitions:  make(map[rent.HistoricalRentTransitionID]rent.HistoricalRentTransition),
		smallPeasantProductions:    make(map[rent.SmallPeasantProductionID]rent.SmallPeasantProduction),
		revenueStreams:             make(map[revenue.RevenueStreamID]revenue.RevenueStream),
		trinityFormulas:            make(map[revenue.TrinityFormulaID]revenue.TrinityFormula),
		revenueFetishForms:         make(map[revenue.RevenueFetishFormID]revenue.RevenueFetishForm),
	}
}

// CreateCostPrice stores cp, assigning an ID and timestamp when absent.
func (m *Memory) CreateCostPrice(_ context.Context, cp profit.CostPrice) (profit.CostPrice, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cp.ID.IsZero() {
		cp.ID = profit.NewCostPriceID()
	}
	if _, exists := m.costPrices[cp.ID]; exists {
		return profit.CostPrice{}, ErrAlreadyExists
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = m.now().UTC()
	}
	m.costPrices[cp.ID] = cp
	return cp, nil
}

// GetCostPrice returns the cost-price with id, or ErrNotFound.
func (m *Memory) GetCostPrice(_ context.Context, id profit.CostPriceID) (profit.CostPrice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cp, ok := m.costPrices[id]
	if !ok {
		return profit.CostPrice{}, ErrNotFound
	}
	return cp, nil
}

// ListCostPrices returns all stored cost-prices, newest first.
func (m *Memory) ListCostPrices(_ context.Context) ([]profit.CostPrice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]profit.CostPrice, 0, len(m.costPrices))
	for _, cp := range m.costPrices {
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateProfitRate stores a, assigning an ID and timestamp when absent.
func (m *Memory) CreateProfitRate(_ context.Context, a profit.ProfitRateAnalysis) (profit.ProfitRateAnalysis, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID.IsZero() {
		a.ID = profit.NewProfitRateID()
	}
	if _, exists := m.profitRates[a.ID]; exists {
		return profit.ProfitRateAnalysis{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	m.profitRates[a.ID] = a
	return a, nil
}

// GetProfitRate returns the analysis with id, or ErrNotFound.
func (m *Memory) GetProfitRate(_ context.Context, id profit.ProfitRateID) (profit.ProfitRateAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.profitRates[id]
	if !ok {
		return profit.ProfitRateAnalysis{}, ErrNotFound
	}
	return a, nil
}

// ListProfitRates returns all stored analyses, newest first.
func (m *Memory) ListProfitRates(_ context.Context) ([]profit.ProfitRateAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]profit.ProfitRateAnalysis, 0, len(m.profitRates))
	for _, a := range m.profitRates {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateVariation stores a, assigning an ID and timestamp when absent.
func (m *Memory) CreateVariation(_ context.Context, a profit.VariationAnalysis) (profit.VariationAnalysis, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID.IsZero() {
		a.ID = profit.NewVariationAnalysisID()
	}
	if _, exists := m.variations[a.ID]; exists {
		return profit.VariationAnalysis{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	m.variations[a.ID] = a
	return a, nil
}

// GetVariation returns the variation analysis with id, or ErrNotFound.
func (m *Memory) GetVariation(_ context.Context, id profit.VariationAnalysisID) (profit.VariationAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.variations[id]
	if !ok {
		return profit.VariationAnalysis{}, ErrNotFound
	}
	return a, nil
}

// ListVariations returns all stored variation analyses, newest first.
func (m *Memory) ListVariations(_ context.Context) ([]profit.VariationAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]profit.VariationAnalysis, 0, len(m.variations))
	for _, a := range m.variations {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateTurnoverAnalysis stores a, assigning an ID and timestamp when absent.
func (m *Memory) CreateTurnoverAnalysis(_ context.Context, a profit.TurnoverAnalysis) (profit.TurnoverAnalysis, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID.IsZero() {
		a.ID = profit.NewTurnoverAnalysisID()
	}
	if _, exists := m.turnovers[a.ID]; exists {
		return profit.TurnoverAnalysis{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	m.turnovers[a.ID] = a
	return a, nil
}

// GetTurnoverAnalysis returns the analysis with id, or ErrNotFound.
func (m *Memory) GetTurnoverAnalysis(_ context.Context, id profit.TurnoverAnalysisID) (profit.TurnoverAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.turnovers[id]
	if !ok {
		return profit.TurnoverAnalysis{}, ErrNotFound
	}
	return a, nil
}

// ListTurnoverAnalyses returns all stored turnover analyses, newest first.
func (m *Memory) ListTurnoverAnalyses(_ context.Context) ([]profit.TurnoverAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]profit.TurnoverAnalysis, 0, len(m.turnovers))
	for _, a := range m.turnovers {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreatePriceFluctuationAnalysis stores a, assigning an ID and timestamp when absent.
func (m *Memory) CreatePriceFluctuationAnalysis(_ context.Context, a profit.PriceFluctuationAnalysis) (profit.PriceFluctuationAnalysis, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID.IsZero() {
		a.ID = profit.NewPriceFluctuationAnalysisID()
	}
	if _, exists := m.priceFlux[a.ID]; exists {
		return profit.PriceFluctuationAnalysis{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	m.priceFlux[a.ID] = a
	return a, nil
}

// GetPriceFluctuationAnalysis returns the analysis with id, or ErrNotFound.
func (m *Memory) GetPriceFluctuationAnalysis(_ context.Context, id profit.PriceFluctuationAnalysisID) (profit.PriceFluctuationAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.priceFlux[id]
	if !ok {
		return profit.PriceFluctuationAnalysis{}, ErrNotFound
	}
	return a, nil
}

// ListPriceFluctuationAnalyses returns all stored analyses, newest first.
func (m *Memory) ListPriceFluctuationAnalyses(_ context.Context) ([]profit.PriceFluctuationAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]profit.PriceFluctuationAnalysis, 0, len(m.priceFlux))
	for _, a := range m.priceFlux {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateEconomyAnalysis stores a, assigning an ID and timestamp when absent.
func (m *Memory) CreateEconomyAnalysis(_ context.Context, a profit.EconomyAnalysis) (profit.EconomyAnalysis, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID.IsZero() {
		a.ID = profit.NewEconomyAnalysisID()
	}
	if _, exists := m.economies[a.ID]; exists {
		return profit.EconomyAnalysis{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	m.economies[a.ID] = a
	return a, nil
}

// GetEconomyAnalysis returns the analysis with id, or ErrNotFound.
func (m *Memory) GetEconomyAnalysis(_ context.Context, id profit.EconomyAnalysisID) (profit.EconomyAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.economies[id]
	if !ok {
		return profit.EconomyAnalysis{}, ErrNotFound
	}
	return a, nil
}

// ListEconomyAnalyses returns all stored economy analyses, newest first.
func (m *Memory) ListEconomyAnalyses(_ context.Context) ([]profit.EconomyAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]profit.EconomyAnalysis, 0, len(m.economies))
	for _, a := range m.economies {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCompositionEffect stores a, assigning an ID and timestamp when absent.
func (m *Memory) CreateCompositionEffect(_ context.Context, a profit.CompositionEffectAnalysis) (profit.CompositionEffectAnalysis, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID.IsZero() {
		a.ID = profit.NewCompositionEffectID()
	}
	if _, exists := m.compositions[a.ID]; exists {
		return profit.CompositionEffectAnalysis{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	m.compositions[a.ID] = a
	return a, nil
}

// GetCompositionEffect returns the comparison with id, or ErrNotFound.
func (m *Memory) GetCompositionEffect(_ context.Context, id profit.CompositionEffectID) (profit.CompositionEffectAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.compositions[id]
	if !ok {
		return profit.CompositionEffectAnalysis{}, ErrNotFound
	}
	return a, nil
}

// ListCompositionEffects returns all stored comparisons, newest first.
func (m *Memory) ListCompositionEffects(_ context.Context) ([]profit.CompositionEffectAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]profit.CompositionEffectAnalysis, 0, len(m.compositions))
	for _, a := range m.compositions {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateMagnitudeChange stores a, assigning an ID and timestamp when absent.
func (m *Memory) CreateMagnitudeChange(_ context.Context, a profit.MagnitudeChangeAnalysis) (profit.MagnitudeChangeAnalysis, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID.IsZero() {
		a.ID = profit.NewMagnitudeChangeID()
	}
	if _, exists := m.magnitudes[a.ID]; exists {
		return profit.MagnitudeChangeAnalysis{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	m.magnitudes[a.ID] = a
	return a, nil
}

// GetMagnitudeChange returns the change with id, or ErrNotFound.
func (m *Memory) GetMagnitudeChange(_ context.Context, id profit.MagnitudeChangeID) (profit.MagnitudeChangeAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.magnitudes[id]
	if !ok {
		return profit.MagnitudeChangeAnalysis{}, ErrNotFound
	}
	return a, nil
}

// ListMagnitudeChanges returns all stored changes, newest first.
func (m *Memory) ListMagnitudeChanges(_ context.Context) ([]profit.MagnitudeChangeAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]profit.MagnitudeChangeAnalysis, 0, len(m.magnitudes))
	for _, a := range m.magnitudes {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateProductionSphere stores s, assigning an ID and timestamp when absent.
func (m *Memory) CreateProductionSphere(_ context.Context, s avgprofit.ProductionSphere) (avgprofit.ProductionSphere, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID.IsZero() {
		s.ID = avgprofit.NewProductionSphereID()
	}
	if _, exists := m.spheres[s.ID]; exists {
		return avgprofit.ProductionSphere{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	m.spheres[s.ID] = s
	return s, nil
}

// GetProductionSphere returns the sphere with id, or ErrNotFound.
func (m *Memory) GetProductionSphere(_ context.Context, id avgprofit.ProductionSphereID) (avgprofit.ProductionSphere, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.spheres[id]
	if !ok {
		return avgprofit.ProductionSphere{}, ErrNotFound
	}
	return s, nil
}

// ListProductionSpheres returns all stored spheres, newest first.
func (m *Memory) ListProductionSpheres(_ context.Context) ([]avgprofit.ProductionSphere, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]avgprofit.ProductionSphere, 0, len(m.spheres))
	for _, s := range m.spheres {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateGeneralProfitRate stores g, assigning an ID and timestamp when absent.
func (m *Memory) CreateGeneralProfitRate(_ context.Context, g avgprofit.GeneralProfitRate) (avgprofit.GeneralProfitRate, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if g.ID.IsZero() {
		g.ID = avgprofit.NewGeneralProfitRateID()
	}
	if _, exists := m.generalRates[g.ID]; exists {
		return avgprofit.GeneralProfitRate{}, ErrAlreadyExists
	}
	if g.CreatedAt.IsZero() {
		g.CreatedAt = m.now().UTC()
	}
	m.generalRates[g.ID] = g
	return g, nil
}

// GetGeneralProfitRate returns the general rate with id, or ErrNotFound.
func (m *Memory) GetGeneralProfitRate(_ context.Context, id avgprofit.GeneralProfitRateID) (avgprofit.GeneralProfitRate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	g, ok := m.generalRates[id]
	if !ok {
		return avgprofit.GeneralProfitRate{}, ErrNotFound
	}
	return g, nil
}

// CreatePriceOfProduction stores p, assigning an ID and timestamp when absent.
func (m *Memory) CreatePriceOfProduction(_ context.Context, p avgprofit.PriceOfProduction) (avgprofit.PriceOfProduction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p.ID.IsZero() {
		p.ID = avgprofit.NewPriceOfProductionID()
	}
	if _, exists := m.pricesOfProduction[p.ID]; exists {
		return avgprofit.PriceOfProduction{}, ErrAlreadyExists
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = m.now().UTC()
	}
	m.pricesOfProduction[p.ID] = p
	return p, nil
}

// GetPriceOfProduction returns the record with id, or ErrNotFound.
func (m *Memory) GetPriceOfProduction(_ context.Context, id avgprofit.PriceOfProductionID) (avgprofit.PriceOfProduction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.pricesOfProduction[id]
	if !ok {
		return avgprofit.PriceOfProduction{}, ErrNotFound
	}
	return p, nil
}

// ListPricesOfProduction returns all stored records, newest first.
func (m *Memory) ListPricesOfProduction(_ context.Context) ([]avgprofit.PriceOfProduction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]avgprofit.PriceOfProduction, 0, len(m.pricesOfProduction))
	for _, p := range m.pricesOfProduction {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateMarketValue stores v, assigning an ID and timestamp when absent.
func (m *Memory) CreateMarketValue(_ context.Context, v avgprofit.MarketValue) (avgprofit.MarketValue, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if v.ID.IsZero() {
		v.ID = avgprofit.NewMarketValueID()
	}
	if _, exists := m.marketValues[v.ID]; exists {
		return avgprofit.MarketValue{}, ErrAlreadyExists
	}
	if v.CreatedAt.IsZero() {
		v.CreatedAt = m.now().UTC()
	}
	m.marketValues[v.ID] = v
	return v, nil
}

// GetMarketValue returns the market-value record with id, or ErrNotFound.
func (m *Memory) GetMarketValue(_ context.Context, id avgprofit.MarketValueID) (avgprofit.MarketValue, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.marketValues[id]
	if !ok {
		return avgprofit.MarketValue{}, ErrNotFound
	}
	return v, nil
}

// CreateSurplusProfit stores s, assigning an ID and timestamp when absent.
func (m *Memory) CreateSurplusProfit(_ context.Context, s avgprofit.SurplusProfit) (avgprofit.SurplusProfit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID.IsZero() {
		s.ID = avgprofit.NewSurplusProfitID()
	}
	if _, exists := m.surplusProfits[s.ID]; exists {
		return avgprofit.SurplusProfit{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	m.surplusProfits[s.ID] = s
	return s, nil
}

// GetSurplusProfit returns the surplus-profit record with id, or ErrNotFound.
func (m *Memory) GetSurplusProfit(_ context.Context, id avgprofit.SurplusProfitID) (avgprofit.SurplusProfit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.surplusProfits[id]
	if !ok {
		return avgprofit.SurplusProfit{}, ErrNotFound
	}
	return s, nil
}

// CreateEqualisation stores e, assigning an ID and timestamp when absent.
func (m *Memory) CreateEqualisation(_ context.Context, e avgprofit.Equalisation) (avgprofit.Equalisation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if e.ID.IsZero() {
		e.ID = avgprofit.NewEqualisationID()
	}
	if _, exists := m.equalisations[e.ID]; exists {
		return avgprofit.Equalisation{}, ErrAlreadyExists
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = m.now().UTC()
	}
	m.equalisations[e.ID] = e
	return e, nil
}

// GetEqualisation returns the equalisation record with id, or ErrNotFound.
func (m *Memory) GetEqualisation(_ context.Context, id avgprofit.EqualisationID) (avgprofit.Equalisation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	e, ok := m.equalisations[id]
	if !ok {
		return avgprofit.Equalisation{}, ErrNotFound
	}
	return e, nil
}

// CreateWageEffectAnalysis stores a, assigning an ID and timestamp when absent.
func (m *Memory) CreateWageEffectAnalysis(_ context.Context, a avgprofit.WageEffectAnalysis) (avgprofit.WageEffectAnalysis, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID.IsZero() {
		a.ID = avgprofit.NewWageEffectAnalysisID()
	}
	if _, exists := m.wageEffects[a.ID]; exists {
		return avgprofit.WageEffectAnalysis{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	if a.Outcomes == nil {
		a.Outcomes = []avgprofit.SphereWageOutcome{}
	}
	m.wageEffects[a.ID] = a
	return a, nil
}

// GetWageEffectAnalysis returns the analysis with id, or ErrNotFound.
func (m *Memory) GetWageEffectAnalysis(_ context.Context, id avgprofit.WageEffectAnalysisID) (avgprofit.WageEffectAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.wageEffects[id]
	if !ok {
		return avgprofit.WageEffectAnalysis{}, ErrNotFound
	}
	return a, nil
}

// ListWageEffectAnalyses returns all stored analyses, newest first.
func (m *Memory) ListWageEffectAnalyses(_ context.Context) ([]avgprofit.WageEffectAnalysis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]avgprofit.WageEffectAnalysis, 0, len(m.wageEffects))
	for _, a := range m.wageEffects {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreatePriceOfProductionChange stores c, assigning an ID and timestamp when absent.
func (m *Memory) CreatePriceOfProductionChange(_ context.Context, c avgprofit.PriceOfProductionChange) (avgprofit.PriceOfProductionChange, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c.ID.IsZero() {
		c.ID = avgprofit.NewPriceOfProductionChangeID()
	}
	if _, exists := m.priceChanges[c.ID]; exists {
		return avgprofit.PriceOfProductionChange{}, ErrAlreadyExists
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}
	m.priceChanges[c.ID] = c
	return c, nil
}

// GetPriceOfProductionChange returns the change with id, or ErrNotFound.
func (m *Memory) GetPriceOfProductionChange(_ context.Context, id avgprofit.PriceOfProductionChangeID) (avgprofit.PriceOfProductionChange, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.priceChanges[id]
	if !ok {
		return avgprofit.PriceOfProductionChange{}, ErrNotFound
	}
	return c, nil
}

// CreateCompositionTrajectory stores t, assigning an ID and timestamp when
// absent. It ensures Periods is non-nil and re-derives ProfitRates before store.
func (m *Memory) CreateCompositionTrajectory(_ context.Context, t tendency.CompositionTrajectory) (tendency.CompositionTrajectory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t.ID.IsZero() {
		t.ID = tendency.NewCompositionTrajectoryID()
	}
	if _, exists := m.trajectories[t.ID]; exists {
		return tendency.CompositionTrajectory{}, ErrAlreadyExists
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = m.now().UTC()
	}
	if t.Periods == nil {
		t.Periods = []tendency.TrajectoryPeriod{}
	}
	t.ProfitRates = t.DeriveProfitRates()
	m.trajectories[t.ID] = t
	return t, nil
}

// GetCompositionTrajectory returns the trajectory with id, or ErrNotFound.
func (m *Memory) GetCompositionTrajectory(_ context.Context, id tendency.CompositionTrajectoryID) (tendency.CompositionTrajectory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	t, ok := m.trajectories[id]
	if !ok {
		return tendency.CompositionTrajectory{}, ErrNotFound
	}
	t.ProfitRates = t.DeriveProfitRates()
	return t, nil
}

// ListCompositionTrajectories returns all stored trajectories, newest first.
func (m *Memory) ListCompositionTrajectories(_ context.Context) ([]tendency.CompositionTrajectory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]tendency.CompositionTrajectory, 0, len(m.trajectories))
	for _, t := range m.trajectories {
		t.ProfitRates = t.DeriveProfitRates()
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateRateMassContradiction stores r, assigning an ID and timestamp when absent.
func (m *Memory) CreateRateMassContradiction(_ context.Context, r tendency.RateMassContradiction) (tendency.RateMassContradiction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r.ID.IsZero() {
		r.ID = tendency.NewRateMassContradictionID()
	}
	if _, exists := m.rateMasses[r.ID]; exists {
		return tendency.RateMassContradiction{}, ErrAlreadyExists
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = m.now().UTC()
	}
	m.rateMasses[r.ID] = r
	return r, nil
}

// GetRateMassContradiction returns the record with id, or ErrNotFound.
func (m *Memory) GetRateMassContradiction(_ context.Context, id tendency.RateMassContradictionID) (tendency.RateMassContradiction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.rateMasses[id]
	if !ok {
		return tendency.RateMassContradiction{}, ErrNotFound
	}
	return r, nil
}

// ListRateMassContradictions returns all stored records, newest first.
func (m *Memory) ListRateMassContradictions(_ context.Context) ([]tendency.RateMassContradiction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]tendency.RateMassContradiction, 0, len(m.rateMasses))
	for _, r := range m.rateMasses {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCounteractingForce stores f, assigning an ID and timestamp when absent.
func (m *Memory) CreateCounteractingForce(_ context.Context, f tendency.CounteractingForce) (tendency.CounteractingForce, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if f.ID.IsZero() {
		f.ID = tendency.NewCounteractingForceID()
	}
	if _, exists := m.counterForces[f.ID]; exists {
		return tendency.CounteractingForce{}, ErrAlreadyExists
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = m.now().UTC()
	}
	m.counterForces[f.ID] = f
	return f, nil
}

// GetCounteractingForce returns the force with id, or ErrNotFound.
func (m *Memory) GetCounteractingForce(_ context.Context, id tendency.CounteractingForceID) (tendency.CounteractingForce, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	f, ok := m.counterForces[id]
	if !ok {
		return tendency.CounteractingForce{}, ErrNotFound
	}
	return f, nil
}

// CreateCounteractingScenario stores s, assigning an ID and timestamp when absent.
func (m *Memory) CreateCounteractingScenario(_ context.Context, s tendency.CounteractingScenario) (tendency.CounteractingScenario, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID.IsZero() {
		s.ID = tendency.NewCounteractingScenarioID()
	}
	if _, exists := m.counterScenarios[s.ID]; exists {
		return tendency.CounteractingScenario{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	if s.Forces == nil {
		s.Forces = []tendency.CounteractingForce{}
	}
	if s.ModifiedRates == nil {
		s.ModifiedRates = []int64{}
	}
	m.counterScenarios[s.ID] = s
	return s, nil
}

// GetCounteractingScenario returns the scenario with id, or ErrNotFound.
func (m *Memory) GetCounteractingScenario(_ context.Context, id tendency.CounteractingScenarioID) (tendency.CounteractingScenario, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.counterScenarios[id]
	if !ok {
		return tendency.CounteractingScenario{}, ErrNotFound
	}
	return s, nil
}

// ListCounteractingScenarios returns all stored scenarios, newest first.
func (m *Memory) ListCounteractingScenarios(_ context.Context) ([]tendency.CounteractingScenario, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]tendency.CounteractingScenario, 0, len(m.counterScenarios))
	for _, s := range m.counterScenarios {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCrisis stores c, assigning an ID and timestamp when absent.
func (m *Memory) CreateCrisis(_ context.Context, c tendency.Crisis) (tendency.Crisis, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c.ID.IsZero() {
		c.ID = tendency.NewCrisisID()
	}
	if _, exists := m.crises[c.ID]; exists {
		return tendency.Crisis{}, ErrAlreadyExists
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}
	m.crises[c.ID] = c
	return c, nil
}

// GetCrisis returns the crisis with id, or ErrNotFound.
func (m *Memory) GetCrisis(_ context.Context, id tendency.CrisisID) (tendency.Crisis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.crises[id]
	if !ok {
		return tendency.Crisis{}, ErrNotFound
	}
	return c, nil
}

// ListCrises returns all stored crises, newest first.
func (m *Memory) ListCrises(_ context.Context) ([]tendency.Crisis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]tendency.Crisis, 0, len(m.crises))
	for _, c := range m.crises {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateInternalContradiction stores c, assigning an ID and timestamp when absent.
func (m *Memory) CreateInternalContradiction(_ context.Context, c tendency.InternalContradiction) (tendency.InternalContradiction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c.ID.IsZero() {
		c.ID = tendency.NewInternalContradictionID()
	}
	if _, exists := m.contradictions[c.ID]; exists {
		return tendency.InternalContradiction{}, ErrAlreadyExists
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}
	m.contradictions[c.ID] = c
	return c, nil
}

// GetInternalContradiction returns the contradiction with id, or ErrNotFound.
func (m *Memory) GetInternalContradiction(_ context.Context, id tendency.InternalContradictionID) (tendency.InternalContradiction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.contradictions[id]
	if !ok {
		return tendency.InternalContradiction{}, ErrNotFound
	}
	return c, nil
}

// ListInternalContradictions returns all stored contradictions, newest first.
func (m *Memory) ListInternalContradictions(_ context.Context) ([]tendency.InternalContradiction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]tendency.InternalContradiction, 0, len(m.contradictions))
	for _, c := range m.contradictions {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCommercialCapital stores cc, assigning an ID and timestamp when absent.
func (m *Memory) CreateCommercialCapital(_ context.Context, cc merchant.CommercialCapital) (merchant.CommercialCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cc.ID.IsZero() {
		cc.ID = merchant.NewCommercialCapitalID()
	}
	if _, exists := m.commercialCapitals[cc.ID]; exists {
		return merchant.CommercialCapital{}, ErrAlreadyExists
	}
	if cc.CreatedAt.IsZero() {
		cc.CreatedAt = m.now().UTC()
	}
	m.commercialCapitals[cc.ID] = cc
	return cc, nil
}

// GetCommercialCapital returns the commercial-capital record with id, or ErrNotFound.
func (m *Memory) GetCommercialCapital(_ context.Context, id merchant.CommercialCapitalID) (merchant.CommercialCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cc, ok := m.commercialCapitals[id]
	if !ok {
		return merchant.CommercialCapital{}, ErrNotFound
	}
	return cc, nil
}

// ListCommercialCapitals returns all stored commercial-capital records, newest first.
func (m *Memory) ListCommercialCapitals(_ context.Context) ([]merchant.CommercialCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]merchant.CommercialCapital, 0, len(m.commercialCapitals))
	for _, cc := range m.commercialCapitals {
		out = append(out, cc)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCommercialProfit stores cp, assigning an ID and timestamp when absent.
func (m *Memory) CreateCommercialProfit(_ context.Context, cp merchant.CommercialProfit) (merchant.CommercialProfit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cp.ID.IsZero() {
		cp.ID = merchant.NewCommercialProfitID()
	}
	if _, exists := m.commercialProfits[cp.ID]; exists {
		return merchant.CommercialProfit{}, ErrAlreadyExists
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = m.now().UTC()
	}
	m.commercialProfits[cp.ID] = cp
	return cp, nil
}

// GetCommercialProfit returns the commercial-profit record with id, or ErrNotFound.
func (m *Memory) GetCommercialProfit(_ context.Context, id merchant.CommercialProfitID) (merchant.CommercialProfit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cp, ok := m.commercialProfits[id]
	if !ok {
		return merchant.CommercialProfit{}, ErrNotFound
	}
	return cp, nil
}

// ListCommercialProfits returns all stored commercial-profit records, newest first.
func (m *Memory) ListCommercialProfits(_ context.Context) ([]merchant.CommercialProfit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]merchant.CommercialProfit, 0, len(m.commercialProfits))
	for _, cp := range m.commercialProfits {
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateMerchantTurnover stores mt, assigning an ID and timestamp when absent.
func (m *Memory) CreateMerchantTurnover(_ context.Context, mt merchant.MerchantTurnover) (merchant.MerchantTurnover, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if mt.ID.IsZero() {
		mt.ID = merchant.NewMerchantTurnoverID()
	}
	if _, exists := m.turnoversM[mt.ID]; exists {
		return merchant.MerchantTurnover{}, ErrAlreadyExists
	}
	if mt.CreatedAt.IsZero() {
		mt.CreatedAt = m.now().UTC()
	}
	m.turnoversM[mt.ID] = mt
	return mt, nil
}

// GetMerchantTurnover returns the merchant-turnover record with id, or ErrNotFound.
func (m *Memory) GetMerchantTurnover(_ context.Context, id merchant.MerchantTurnoverID) (merchant.MerchantTurnover, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mt, ok := m.turnoversM[id]
	if !ok {
		return merchant.MerchantTurnover{}, ErrNotFound
	}
	return mt, nil
}

// ListMerchantTurnovers returns all stored merchant-turnover records, newest first.
func (m *Memory) ListMerchantTurnovers(_ context.Context) ([]merchant.MerchantTurnover, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]merchant.MerchantTurnover, 0, len(m.turnoversM))
	for _, mt := range m.turnoversM {
		out = append(out, mt)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateMoneyDealingCapital stores m, assigning an ID and timestamp when absent.
func (m *Memory) CreateMoneyDealingCapital(_ context.Context, md merchant.MoneyDealingCapital) (merchant.MoneyDealingCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if md.ID.IsZero() {
		md.ID = merchant.NewMoneyDealingCapitalID()
	}
	if _, exists := m.moneyDealingCapitals[md.ID]; exists {
		return merchant.MoneyDealingCapital{}, ErrAlreadyExists
	}
	if md.CreatedAt.IsZero() {
		md.CreatedAt = m.now().UTC()
	}
	m.moneyDealingCapitals[md.ID] = md
	return md, nil
}

// GetMoneyDealingCapital returns the money-dealing-capital record with id, or ErrNotFound.
func (m *Memory) GetMoneyDealingCapital(_ context.Context, id merchant.MoneyDealingCapitalID) (merchant.MoneyDealingCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	md, ok := m.moneyDealingCapitals[id]
	if !ok {
		return merchant.MoneyDealingCapital{}, ErrNotFound
	}
	return md, nil
}

// ListMoneyDealingCapitals returns all stored money-dealing-capital records, newest first.
func (m *Memory) ListMoneyDealingCapitals(_ context.Context) ([]merchant.MoneyDealingCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]merchant.MoneyDealingCapital, 0, len(m.moneyDealingCapitals))
	for _, md := range m.moneyDealingCapitals {
		out = append(out, md)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateHistoricalMerchantCapital stores hm, assigning an ID and timestamp when absent.
func (m *Memory) CreateHistoricalMerchantCapital(_ context.Context, hm merchant.HistoricalMerchantCapital) (merchant.HistoricalMerchantCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hm.ID.IsZero() {
		hm.ID = merchant.NewHistoricalMerchantCapitalID()
	}
	if _, exists := m.historicalMerchantCapitals[hm.ID]; exists {
		return merchant.HistoricalMerchantCapital{}, ErrAlreadyExists
	}
	if hm.CreatedAt.IsZero() {
		hm.CreatedAt = m.now().UTC()
	}
	m.historicalMerchantCapitals[hm.ID] = hm
	return hm, nil
}

// GetHistoricalMerchantCapital returns the historical-merchant-capital record with id, or ErrNotFound.
func (m *Memory) GetHistoricalMerchantCapital(_ context.Context, id merchant.HistoricalMerchantCapitalID) (merchant.HistoricalMerchantCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hm, ok := m.historicalMerchantCapitals[id]
	if !ok {
		return merchant.HistoricalMerchantCapital{}, ErrNotFound
	}
	return hm, nil
}

// ListHistoricalMerchantCapitals returns all stored historical-merchant-capital records, newest first.
func (m *Memory) ListHistoricalMerchantCapitals(_ context.Context) ([]merchant.HistoricalMerchantCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]merchant.HistoricalMerchantCapital, 0, len(m.historicalMerchantCapitals))
	for _, hm := range m.historicalMerchantCapitals {
		out = append(out, hm)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateInterestBearingCapital stores ibc, assigning an ID and timestamp when absent.
func (m *Memory) CreateInterestBearingCapital(_ context.Context, ibc credit.InterestBearingCapital) (credit.InterestBearingCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ibc.ID.IsZero() {
		ibc.ID = credit.NewInterestBearingCapitalID()
	}
	if _, exists := m.interestBearingCapitals[ibc.ID]; exists {
		return credit.InterestBearingCapital{}, ErrAlreadyExists
	}
	if ibc.CreatedAt.IsZero() {
		ibc.CreatedAt = m.now().UTC()
	}
	m.interestBearingCapitals[ibc.ID] = ibc
	return ibc, nil
}

// GetInterestBearingCapital returns the interest-bearing-capital record with id, or ErrNotFound.
func (m *Memory) GetInterestBearingCapital(_ context.Context, id credit.InterestBearingCapitalID) (credit.InterestBearingCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ibc, ok := m.interestBearingCapitals[id]
	if !ok {
		return credit.InterestBearingCapital{}, ErrNotFound
	}
	return ibc, nil
}

// ListInterestBearingCapitals returns all stored interest-bearing-capital records, newest first.
func (m *Memory) ListInterestBearingCapitals(_ context.Context) ([]credit.InterestBearingCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.InterestBearingCapital, 0, len(m.interestBearingCapitals))
	for _, ibc := range m.interestBearingCapitals {
		out = append(out, ibc)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateRateOfInterest stores r, assigning an ID and timestamp when absent.
func (m *Memory) CreateRateOfInterest(_ context.Context, r credit.RateOfInterest) (credit.RateOfInterest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r.ID.IsZero() {
		r.ID = credit.NewRateOfInterestID()
	}
	if _, exists := m.ratesOfInterest[r.ID]; exists {
		return credit.RateOfInterest{}, ErrAlreadyExists
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = m.now().UTC()
	}
	m.ratesOfInterest[r.ID] = r
	return r, nil
}

// GetRateOfInterest returns the rate-of-interest record with id, or ErrNotFound.
func (m *Memory) GetRateOfInterest(_ context.Context, id credit.RateOfInterestID) (credit.RateOfInterest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.ratesOfInterest[id]
	if !ok {
		return credit.RateOfInterest{}, ErrNotFound
	}
	return r, nil
}

// ListRatesOfInterest returns all stored rate-of-interest records, newest first.
func (m *Memory) ListRatesOfInterest(_ context.Context) ([]credit.RateOfInterest, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.RateOfInterest, 0, len(m.ratesOfInterest))
	for _, r := range m.ratesOfInterest {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateProfitDivision stores pd, assigning an ID and timestamp when absent.
func (m *Memory) CreateProfitDivision(_ context.Context, pd credit.ProfitDivision) (credit.ProfitDivision, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if pd.ID.IsZero() {
		pd.ID = credit.NewProfitDivisionID()
	}
	if _, exists := m.profitDivisions[pd.ID]; exists {
		return credit.ProfitDivision{}, ErrAlreadyExists
	}
	if pd.CreatedAt.IsZero() {
		pd.CreatedAt = m.now().UTC()
	}
	m.profitDivisions[pd.ID] = pd
	return pd, nil
}

// GetProfitDivision returns the profit-division record with id, or ErrNotFound.
func (m *Memory) GetProfitDivision(_ context.Context, id credit.ProfitDivisionID) (credit.ProfitDivision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pd, ok := m.profitDivisions[id]
	if !ok {
		return credit.ProfitDivision{}, ErrNotFound
	}
	return pd, nil
}

// ListProfitDivisions returns all stored profit-division records, newest first.
func (m *Memory) ListProfitDivisions(_ context.Context) ([]credit.ProfitDivision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.ProfitDivision, 0, len(m.profitDivisions))
	for _, pd := range m.profitDivisions {
		out = append(out, pd)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCompoundInterestSchedule stores s, assigning an ID and timestamp when absent.
func (m *Memory) CreateCompoundInterestSchedule(_ context.Context, s credit.CompoundInterestSchedule) (credit.CompoundInterestSchedule, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID.IsZero() {
		s.ID = credit.NewCompoundInterestID()
	}
	if _, exists := m.compoundInterestSchedules[s.ID]; exists {
		return credit.CompoundInterestSchedule{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	m.compoundInterestSchedules[s.ID] = s
	return s, nil
}

// GetCompoundInterestSchedule returns the schedule with id, or ErrNotFound.
func (m *Memory) GetCompoundInterestSchedule(_ context.Context, id credit.CompoundInterestID) (credit.CompoundInterestSchedule, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.compoundInterestSchedules[id]
	if !ok {
		return credit.CompoundInterestSchedule{}, ErrNotFound
	}
	return s, nil
}

// ListCompoundInterestSchedules returns all stored schedules, newest first.
func (m *Memory) ListCompoundInterestSchedules(_ context.Context) ([]credit.CompoundInterestSchedule, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.CompoundInterestSchedule, 0, len(m.compoundInterestSchedules))
	for _, s := range m.compoundInterestSchedules {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateBillOfExchange stores b, assigning an ID and timestamp when absent.
func (m *Memory) CreateBillOfExchange(_ context.Context, b credit.BillOfExchange) (credit.BillOfExchange, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if b.ID.IsZero() {
		b.ID = credit.NewBillOfExchangeID()
	}
	if _, exists := m.billsOfExchange[b.ID]; exists {
		return credit.BillOfExchange{}, ErrAlreadyExists
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = m.now().UTC()
	}
	m.billsOfExchange[b.ID] = b
	return b, nil
}

// GetBillOfExchange returns the bill with id, or ErrNotFound.
func (m *Memory) GetBillOfExchange(_ context.Context, id credit.BillOfExchangeID) (credit.BillOfExchange, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, ok := m.billsOfExchange[id]
	if !ok {
		return credit.BillOfExchange{}, ErrNotFound
	}
	return b, nil
}

// ListBillsOfExchange returns all stored bills, newest first.
func (m *Memory) ListBillsOfExchange(_ context.Context) ([]credit.BillOfExchange, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.BillOfExchange, 0, len(m.billsOfExchange))
	for _, b := range m.billsOfExchange {
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateFictitiousCapital stores fc, assigning an ID and timestamp when absent.
func (m *Memory) CreateFictitiousCapital(_ context.Context, fc credit.FictitiousCapital) (credit.FictitiousCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if fc.ID.IsZero() {
		fc.ID = credit.NewFictitiousCapitalID()
	}
	if _, exists := m.fictitiousCapitals[fc.ID]; exists {
		return credit.FictitiousCapital{}, ErrAlreadyExists
	}
	if fc.CreatedAt.IsZero() {
		fc.CreatedAt = m.now().UTC()
	}
	m.fictitiousCapitals[fc.ID] = fc
	return fc, nil
}

// GetFictitiousCapital returns the record with id, or ErrNotFound.
func (m *Memory) GetFictitiousCapital(_ context.Context, id credit.FictitiousCapitalID) (credit.FictitiousCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fc, ok := m.fictitiousCapitals[id]
	if !ok {
		return credit.FictitiousCapital{}, ErrNotFound
	}
	return fc, nil
}

// ListFictitiousCapitals returns all stored records, newest first.
func (m *Memory) ListFictitiousCapitals(_ context.Context) ([]credit.FictitiousCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.FictitiousCapital, 0, len(m.fictitiousCapitals))
	for _, fc := range m.fictitiousCapitals {
		out = append(out, fc)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateMoneyCapitalAccumulation stores a, assigning an ID and timestamp when absent (Vol. III Ch. 26).
func (m *Memory) CreateMoneyCapitalAccumulation(_ context.Context, a credit.MoneyCapitalAccumulation) (credit.MoneyCapitalAccumulation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID.IsZero() {
		a.ID = credit.NewMoneyCapitalAccumulationID()
	}
	if _, exists := m.moneyCapitalAccumulations[a.ID]; exists {
		return credit.MoneyCapitalAccumulation{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	m.moneyCapitalAccumulations[a.ID] = a
	return a, nil
}

// GetMoneyCapitalAccumulation returns the record with id, or ErrNotFound.
func (m *Memory) GetMoneyCapitalAccumulation(_ context.Context, id credit.MoneyCapitalAccumulationID) (credit.MoneyCapitalAccumulation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.moneyCapitalAccumulations[id]
	if !ok {
		return credit.MoneyCapitalAccumulation{}, ErrNotFound
	}
	return a, nil
}

// ListMoneyCapitalAccumulations returns all stored records, newest first.
func (m *Memory) ListMoneyCapitalAccumulations(_ context.Context) ([]credit.MoneyCapitalAccumulation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.MoneyCapitalAccumulation, 0, len(m.moneyCapitalAccumulations))
	for _, a := range m.moneyCapitalAccumulations {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateStockCompany stores sc, assigning an ID and timestamp when absent (Vol. III Ch. 27).
func (m *Memory) CreateStockCompany(_ context.Context, sc credit.StockCompany) (credit.StockCompany, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sc.ID.IsZero() {
		sc.ID = credit.NewStockCompanyID()
	}
	if _, exists := m.stockCompanies[sc.ID]; exists {
		return credit.StockCompany{}, ErrAlreadyExists
	}
	if sc.CreatedAt.IsZero() {
		sc.CreatedAt = m.now().UTC()
	}
	m.stockCompanies[sc.ID] = sc
	return sc, nil
}

// GetStockCompany returns the stock company with id, or ErrNotFound.
func (m *Memory) GetStockCompany(_ context.Context, id credit.StockCompanyID) (credit.StockCompany, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sc, ok := m.stockCompanies[id]
	if !ok {
		return credit.StockCompany{}, ErrNotFound
	}
	return sc, nil
}

// ListStockCompanies returns all stored stock companies, newest first.
func (m *Memory) ListStockCompanies(_ context.Context) ([]credit.StockCompany, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.StockCompany, 0, len(m.stockCompanies))
	for _, sc := range m.stockCompanies {
		out = append(out, sc)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCooperativeFactory stores cf, assigning an ID and timestamp when absent (Vol. III Ch. 27).
func (m *Memory) CreateCooperativeFactory(_ context.Context, cf credit.CooperativeFactory) (credit.CooperativeFactory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cf.ID.IsZero() {
		cf.ID = credit.NewCooperativeFactoryID()
	}
	if _, exists := m.cooperativeFactories[cf.ID]; exists {
		return credit.CooperativeFactory{}, ErrAlreadyExists
	}
	if cf.CreatedAt.IsZero() {
		cf.CreatedAt = m.now().UTC()
	}
	m.cooperativeFactories[cf.ID] = cf
	return cf, nil
}

// GetCooperativeFactory returns the cooperative factory with id, or ErrNotFound.
func (m *Memory) GetCooperativeFactory(_ context.Context, id credit.CooperativeFactoryID) (credit.CooperativeFactory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cf, ok := m.cooperativeFactories[id]
	if !ok {
		return credit.CooperativeFactory{}, ErrNotFound
	}
	return cf, nil
}

// ListCooperativeFactories returns all stored cooperative factories, newest first.
func (m *Memory) ListCooperativeFactories(_ context.Context) ([]credit.CooperativeFactory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.CooperativeFactory, 0, len(m.cooperativeFactories))
	for _, cf := range m.cooperativeFactories {
		out = append(out, cf)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCurrencyObservation stores o, assigning an ID and timestamp when absent (Vol. III Ch. 28).
func (m *Memory) CreateCurrencyObservation(_ context.Context, o credit.CurrencyObservation) (credit.CurrencyObservation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if o.ID.IsZero() {
		o.ID = credit.NewCurrencyObservationID()
	}
	if _, exists := m.currencyObservations[o.ID]; exists {
		return credit.CurrencyObservation{}, ErrAlreadyExists
	}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = m.now().UTC()
	}
	m.currencyObservations[o.ID] = o
	return o, nil
}

// GetCurrencyObservation returns the record with id, or ErrNotFound.
func (m *Memory) GetCurrencyObservation(_ context.Context, id credit.CurrencyObservationID) (credit.CurrencyObservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	o, ok := m.currencyObservations[id]
	if !ok {
		return credit.CurrencyObservation{}, ErrNotFound
	}
	return o, nil
}

// ListCurrencyObservations returns all stored records, newest first.
func (m *Memory) ListCurrencyObservations(_ context.Context) ([]credit.CurrencyObservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.CurrencyObservation, 0, len(m.currencyObservations))
	for _, o := range m.currencyObservations {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// cloneBankComponents returns a deep copy of the component slice (nil → empty),
// so the in-memory store never aliases the caller's slice — matching the
// JSON round-trip behaviour of the MySQL store.
func cloneBankComponents(in []credit.BankCapitalComponent) []credit.BankCapitalComponent {
	out := make([]credit.BankCapitalComponent, len(in))
	copy(out, in)
	return out
}

// CreateBankCapital stores bc, assigning an ID and timestamp when absent (Vol. III Ch. 29).
func (m *Memory) CreateBankCapital(_ context.Context, bc credit.BankCapital) (credit.BankCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if bc.ID.IsZero() {
		bc.ID = credit.NewBankCapitalID()
	}
	if _, exists := m.bankCapitals[bc.ID]; exists {
		return credit.BankCapital{}, ErrAlreadyExists
	}
	if bc.CreatedAt.IsZero() {
		bc.CreatedAt = m.now().UTC()
	}
	bc.Components = cloneBankComponents(bc.Components)
	m.bankCapitals[bc.ID] = bc
	return bc, nil
}

// GetBankCapital returns the bank-capital record with id, or ErrNotFound.
func (m *Memory) GetBankCapital(_ context.Context, id credit.BankCapitalID) (credit.BankCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	bc, ok := m.bankCapitals[id]
	if !ok {
		return credit.BankCapital{}, ErrNotFound
	}
	bc.Components = cloneBankComponents(bc.Components)
	return bc, nil
}

// ListBankCapitals returns all stored bank-capital records, newest first.
func (m *Memory) ListBankCapitals(_ context.Context) ([]credit.BankCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.BankCapital, 0, len(m.bankCapitals))
	for _, bc := range m.bankCapitals {
		bc.Components = cloneBankComponents(bc.Components)
		out = append(out, bc)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateFictitiousCapitalValuation stores v, assigning an ID and timestamp when absent (Vol. III Ch. 29).
func (m *Memory) CreateFictitiousCapitalValuation(_ context.Context, v credit.FictitiousCapitalValuation) (credit.FictitiousCapitalValuation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if v.ID.IsZero() {
		v.ID = credit.NewFictitiousCapitalValuationID()
	}
	if _, exists := m.fictitiousValuations[v.ID]; exists {
		return credit.FictitiousCapitalValuation{}, ErrAlreadyExists
	}
	if v.CreatedAt.IsZero() {
		v.CreatedAt = m.now().UTC()
	}
	m.fictitiousValuations[v.ID] = v
	return v, nil
}

// GetFictitiousCapitalValuation returns the valuation with id, or ErrNotFound.
func (m *Memory) GetFictitiousCapitalValuation(_ context.Context, id credit.FictitiousCapitalValuationID) (credit.FictitiousCapitalValuation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.fictitiousValuations[id]
	if !ok {
		return credit.FictitiousCapitalValuation{}, ErrNotFound
	}
	return v, nil
}

// ListFictitiousCapitalValuations returns all stored valuations, newest first.
func (m *Memory) ListFictitiousCapitalValuations(_ context.Context) ([]credit.FictitiousCapitalValuation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.FictitiousCapitalValuation, 0, len(m.fictitiousValuations))
	for _, v := range m.fictitiousValuations {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateRealCapitalAccumulation stores a, assigning an ID and timestamp when absent (Vol. III Ch. 30).
func (m *Memory) CreateRealCapitalAccumulation(_ context.Context, a credit.RealCapitalAccumulation) (credit.RealCapitalAccumulation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID.IsZero() {
		a.ID = credit.NewRealCapitalAccumulationID()
	}
	if _, exists := m.realCapitalAccumulations[a.ID]; exists {
		return credit.RealCapitalAccumulation{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	m.realCapitalAccumulations[a.ID] = a
	return a, nil
}

// GetRealCapitalAccumulation returns the record with id, or ErrNotFound.
func (m *Memory) GetRealCapitalAccumulation(_ context.Context, id credit.RealCapitalAccumulationID) (credit.RealCapitalAccumulation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	a, ok := m.realCapitalAccumulations[id]
	if !ok {
		return credit.RealCapitalAccumulation{}, ErrNotFound
	}
	return a, nil
}

// ListRealCapitalAccumulations returns all stored records, newest first.
func (m *Memory) ListRealCapitalAccumulations(_ context.Context) ([]credit.RealCapitalAccumulation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.RealCapitalAccumulation, 0, len(m.realCapitalAccumulations))
	for _, a := range m.realCapitalAccumulations {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateFloatingCapital stores f, assigning an ID and timestamp when absent (Vol. III Ch. 31).
func (m *Memory) CreateFloatingCapital(_ context.Context, f credit.FloatingCapital) (credit.FloatingCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if f.ID.IsZero() {
		f.ID = credit.NewFloatingCapitalID()
	}
	if _, exists := m.floatingCapitals[f.ID]; exists {
		return credit.FloatingCapital{}, ErrAlreadyExists
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = m.now().UTC()
	}
	m.floatingCapitals[f.ID] = f
	return f, nil
}

// GetFloatingCapital returns the record with id, or ErrNotFound.
func (m *Memory) GetFloatingCapital(_ context.Context, id credit.FloatingCapitalID) (credit.FloatingCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	f, ok := m.floatingCapitals[id]
	if !ok {
		return credit.FloatingCapital{}, ErrNotFound
	}
	return f, nil
}

// ListFloatingCapitals returns all stored records, newest first.
func (m *Memory) ListFloatingCapitals(_ context.Context) ([]credit.FloatingCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.FloatingCapital, 0, len(m.floatingCapitals))
	for _, f := range m.floatingCapitals {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCapitalRelease stores c, assigning an ID and timestamp when absent (Vol. III Ch. 32).
func (m *Memory) CreateCapitalRelease(_ context.Context, c credit.CapitalRelease) (credit.CapitalRelease, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c.ID.IsZero() {
		c.ID = credit.NewCapitalReleaseID()
	}
	if _, exists := m.capitalReleases[c.ID]; exists {
		return credit.CapitalRelease{}, ErrAlreadyExists
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}
	m.capitalReleases[c.ID] = c
	return c, nil
}

// GetCapitalRelease returns the record with id, or ErrNotFound.
func (m *Memory) GetCapitalRelease(_ context.Context, id credit.CapitalReleaseID) (credit.CapitalRelease, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.capitalReleases[id]
	if !ok {
		return credit.CapitalRelease{}, ErrNotFound
	}
	return c, nil
}

// ListCapitalReleases returns all stored records, newest first.
func (m *Memory) ListCapitalReleases(_ context.Context) ([]credit.CapitalRelease, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.CapitalRelease, 0, len(m.capitalReleases))
	for _, c := range m.capitalReleases {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateClearingHouseSettlement stores s, assigning an ID and timestamp when absent (Vol. III Ch. 33).
func (m *Memory) CreateClearingHouseSettlement(_ context.Context, s credit.ClearingHouseSettlement) (credit.ClearingHouseSettlement, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID.IsZero() {
		s.ID = credit.NewClearingHouseSettlementID()
	}
	if _, exists := m.clearingHouseSettlements[s.ID]; exists {
		return credit.ClearingHouseSettlement{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	m.clearingHouseSettlements[s.ID] = s
	return s, nil
}

// GetClearingHouseSettlement returns the record with id, or ErrNotFound.
func (m *Memory) GetClearingHouseSettlement(_ context.Context, id credit.ClearingHouseSettlementID) (credit.ClearingHouseSettlement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.clearingHouseSettlements[id]
	if !ok {
		return credit.ClearingHouseSettlement{}, ErrNotFound
	}
	return s, nil
}

// ListClearingHouseSettlements returns all stored records, newest first.
func (m *Memory) ListClearingHouseSettlements(_ context.Context) ([]credit.ClearingHouseSettlement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.ClearingHouseSettlement, 0, len(m.clearingHouseSettlements))
	for _, s := range m.clearingHouseSettlements {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateNoteIssueConstraint stores n, assigning an ID and timestamp when absent (Vol. III Ch. 34).
func (m *Memory) CreateNoteIssueConstraint(_ context.Context, n credit.NoteIssueConstraint) (credit.NoteIssueConstraint, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if n.ID.IsZero() {
		n.ID = credit.NewNoteIssueConstraintID()
	}
	if _, exists := m.noteIssueConstraints[n.ID]; exists {
		return credit.NoteIssueConstraint{}, ErrAlreadyExists
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = m.now().UTC()
	}
	m.noteIssueConstraints[n.ID] = n
	return n, nil
}

// GetNoteIssueConstraint returns the record with id, or ErrNotFound.
func (m *Memory) GetNoteIssueConstraint(_ context.Context, id credit.NoteIssueConstraintID) (credit.NoteIssueConstraint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	n, ok := m.noteIssueConstraints[id]
	if !ok {
		return credit.NoteIssueConstraint{}, ErrNotFound
	}
	return n, nil
}

// ListNoteIssueConstraints returns all stored records, newest first.
func (m *Memory) ListNoteIssueConstraints(_ context.Context) ([]credit.NoteIssueConstraint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.NoteIssueConstraint, 0, len(m.noteIssueConstraints))
	for _, n := range m.noteIssueConstraints {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateGoldReserve stores g, assigning an ID and timestamp when absent (Vol. III Ch. 35).
func (m *Memory) CreateGoldReserve(_ context.Context, g credit.GoldReserve) (credit.GoldReserve, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if g.ID.IsZero() {
		g.ID = credit.NewGoldReserveID()
	}
	if _, exists := m.goldReserves[g.ID]; exists {
		return credit.GoldReserve{}, ErrAlreadyExists
	}
	if g.CreatedAt.IsZero() {
		g.CreatedAt = m.now().UTC()
	}
	m.goldReserves[g.ID] = g
	return g, nil
}

// GetGoldReserve returns the record with id, or ErrNotFound.
func (m *Memory) GetGoldReserve(_ context.Context, id credit.GoldReserveID) (credit.GoldReserve, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	g, ok := m.goldReserves[id]
	if !ok {
		return credit.GoldReserve{}, ErrNotFound
	}
	return g, nil
}

// ListGoldReserves returns all stored records, newest first.
func (m *Memory) ListGoldReserves(_ context.Context) ([]credit.GoldReserve, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.GoldReserve, 0, len(m.goldReserves))
	for _, g := range m.goldReserves {
		out = append(out, g)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateRateOfExchange stores r, assigning an ID and timestamp when absent (Vol. III Ch. 35).
func (m *Memory) CreateRateOfExchange(_ context.Context, r credit.RateOfExchange) (credit.RateOfExchange, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r.ID.IsZero() {
		r.ID = credit.NewRateOfExchangeID()
	}
	if _, exists := m.ratesOfExchange[r.ID]; exists {
		return credit.RateOfExchange{}, ErrAlreadyExists
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = m.now().UTC()
	}
	m.ratesOfExchange[r.ID] = r
	return r, nil
}

// GetRateOfExchange returns the record with id, or ErrNotFound.
func (m *Memory) GetRateOfExchange(_ context.Context, id credit.RateOfExchangeID) (credit.RateOfExchange, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.ratesOfExchange[id]
	if !ok {
		return credit.RateOfExchange{}, ErrNotFound
	}
	return r, nil
}

// ListRatesOfExchange returns all stored records, newest first.
func (m *Memory) ListRatesOfExchange(_ context.Context) ([]credit.RateOfExchange, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.RateOfExchange, 0, len(m.ratesOfExchange))
	for _, r := range m.ratesOfExchange {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// cloneBorrowers returns a deep copy of a borrowers slice.
func cloneBorrowers(src []string) []string {
	out := make([]string, len(src))
	copy(out, src)
	return out
}

// CreateUsurersCapital stores u, assigning an ID and timestamp when absent (Vol. III Ch. 36).
func (m *Memory) CreateUsurersCapital(_ context.Context, u credit.UsurersCapital) (credit.UsurersCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if u.ID.IsZero() {
		u.ID = credit.NewUsurersCapitalID()
	}
	if _, exists := m.usurersCapitals[u.ID]; exists {
		return credit.UsurersCapital{}, ErrAlreadyExists
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = m.now().UTC()
	}
	u.Borrowers = cloneBorrowers(u.Borrowers)
	m.usurersCapitals[u.ID] = u
	return u, nil
}

// GetUsurersCapital returns the record with id, or ErrNotFound.
func (m *Memory) GetUsurersCapital(_ context.Context, id credit.UsurersCapitalID) (credit.UsurersCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	u, ok := m.usurersCapitals[id]
	if !ok {
		return credit.UsurersCapital{}, ErrNotFound
	}
	u.Borrowers = cloneBorrowers(u.Borrowers)
	return u, nil
}

// ListUsurersCapitals returns all stored records, newest first. Never returns nil.
func (m *Memory) ListUsurersCapitals(_ context.Context) ([]credit.UsurersCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]credit.UsurersCapital, 0, len(m.usurersCapitals))
	for _, u := range m.usurersCapitals {
		u.Borrowers = cloneBorrowers(u.Borrowers)
		out = append(out, u)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateLandParcel stores p, assigning an ID and timestamp when absent (Vol. III Ch. 37).
func (m *Memory) CreateLandParcel(_ context.Context, p rent.LandParcel) (rent.LandParcel, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p.ID == "" {
		p.ID = rent.NewLandParcelID()
	}
	if _, exists := m.landParcels[p.ID]; exists {
		return rent.LandParcel{}, ErrAlreadyExists
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = m.now().UTC()
	}
	m.landParcels[p.ID] = p
	return p, nil
}

// GetLandParcel returns the land-parcel with id, or ErrNotFound.
func (m *Memory) GetLandParcel(_ context.Context, id rent.LandParcelID) (rent.LandParcel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.landParcels[id]
	if !ok {
		return rent.LandParcel{}, ErrNotFound
	}
	return p, nil
}

// ListLandParcels returns all stored land-parcels, newest first. Never returns nil.
func (m *Memory) ListLandParcels(_ context.Context) ([]rent.LandParcel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.LandParcel, 0, len(m.landParcels))
	for _, p := range m.landParcels {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateLandowner stores lo, assigning an ID and timestamp when absent (Vol. III Ch. 37).
func (m *Memory) CreateLandowner(_ context.Context, lo rent.Landowner) (rent.Landowner, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if lo.ID == "" {
		lo.ID = rent.NewLandownerID()
	}
	if _, exists := m.landowners[lo.ID]; exists {
		return rent.Landowner{}, ErrAlreadyExists
	}
	if lo.CreatedAt.IsZero() {
		lo.CreatedAt = m.now().UTC()
	}
	m.landowners[lo.ID] = lo
	return lo, nil
}

// GetLandowner returns the landowner with id, or ErrNotFound.
func (m *Memory) GetLandowner(_ context.Context, id rent.LandownerID) (rent.Landowner, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lo, ok := m.landowners[id]
	if !ok {
		return rent.Landowner{}, ErrNotFound
	}
	return lo, nil
}

// CreateAgriculturalCapitalist stores ac, assigning an ID and timestamp when absent (Vol. III Ch. 37).
func (m *Memory) CreateAgriculturalCapitalist(_ context.Context, ac rent.AgriculturalCapitalist) (rent.AgriculturalCapitalist, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ac.ID == "" {
		ac.ID = rent.NewAgriculturalCapitalistID()
	}
	if _, exists := m.agriculturalCapitalists[ac.ID]; exists {
		return rent.AgriculturalCapitalist{}, ErrAlreadyExists
	}
	if ac.CreatedAt.IsZero() {
		ac.CreatedAt = m.now().UTC()
	}
	m.agriculturalCapitalists[ac.ID] = ac
	return ac, nil
}

// GetAgriculturalCapitalist returns the agricultural capitalist with id, or ErrNotFound.
func (m *Memory) GetAgriculturalCapitalist(_ context.Context, id rent.AgriculturalCapitalistID) (rent.AgriculturalCapitalist, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ac, ok := m.agriculturalCapitalists[id]
	if !ok {
		return rent.AgriculturalCapitalist{}, ErrNotFound
	}
	return ac, nil
}

// CreateGroundRent stores g, assigning an ID and timestamp when absent (Vol. III Ch. 37).
func (m *Memory) CreateGroundRent(_ context.Context, g rent.GroundRent) (rent.GroundRent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if g.ID == "" {
		g.ID = rent.NewGroundRentID()
	}
	if _, exists := m.groundRents[g.ID]; exists {
		return rent.GroundRent{}, ErrAlreadyExists
	}
	if g.CreatedAt.IsZero() {
		g.CreatedAt = m.now().UTC()
	}
	m.groundRents[g.ID] = g
	return g, nil
}

// ListGroundRents returns all stored ground-rent records, newest first. Never returns nil.
func (m *Memory) ListGroundRents(_ context.Context) ([]rent.GroundRent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.GroundRent, 0, len(m.groundRents))
	for _, g := range m.groundRents {
		out = append(out, g)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreatePoPSurplusProfit stores s, assigning an ID and timestamp when absent (Vol. III Ch. 38).
func (m *Memory) CreatePoPSurplusProfit(_ context.Context, s rent.PriceOfProductionSurplusProfit) (rent.PriceOfProductionSurplusProfit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID == "" {
		s.ID = rent.NewPoPSurplusProfitID()
	}
	if _, exists := m.popSurplusProfits[s.ID]; exists {
		return rent.PriceOfProductionSurplusProfit{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	m.popSurplusProfits[s.ID] = s
	return s, nil
}

// ListPoPSurplusProfits returns all stored PoP surplus-profit records, newest first. Never returns nil.
func (m *Memory) ListPoPSurplusProfits(_ context.Context) ([]rent.PriceOfProductionSurplusProfit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.PriceOfProductionSurplusProfit, 0, len(m.popSurplusProfits))
	for _, s := range m.popSurplusProfits {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateMonopolisedNaturalForce stores f, assigning an ID and timestamp when absent (Vol. III Ch. 38).
func (m *Memory) CreateMonopolisedNaturalForce(_ context.Context, f rent.MonopolisedNaturalForce) (rent.MonopolisedNaturalForce, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if f.ID == "" {
		f.ID = rent.NewMonopolisedNaturalForceID()
	}
	if _, exists := m.monopolisedForces[f.ID]; exists {
		return rent.MonopolisedNaturalForce{}, ErrAlreadyExists
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = m.now().UTC()
	}
	m.monopolisedForces[f.ID] = f
	return f, nil
}

// ListMonopolisedNaturalForces returns all stored monopolised-natural-force records, newest first. Never returns nil.
func (m *Memory) ListMonopolisedNaturalForces(_ context.Context) ([]rent.MonopolisedNaturalForce, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.MonopolisedNaturalForce, 0, len(m.monopolisedForces))
	for _, f := range m.monopolisedForces {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateDifferentialRent stores dr, assigning an ID and timestamp when absent (Vol. III Ch. 38).
func (m *Memory) CreateDifferentialRent(_ context.Context, dr rent.DifferentialRent) (rent.DifferentialRent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if dr.ID == "" {
		dr.ID = rent.NewDifferentialRentID()
	}
	if _, exists := m.differentialRents[dr.ID]; exists {
		return rent.DifferentialRent{}, ErrAlreadyExists
	}
	if dr.CreatedAt.IsZero() {
		dr.CreatedAt = m.now().UTC()
	}
	m.differentialRents[dr.ID] = dr
	return dr, nil
}

// ListDifferentialRents returns all stored differential-rent records, newest first. Never returns nil.
func (m *Memory) ListDifferentialRents(_ context.Context) ([]rent.DifferentialRent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.DifferentialRent, 0, len(m.differentialRents))
	for _, dr := range m.differentialRents {
		out = append(out, dr)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateCapitalisedRentPrice stores crp, assigning an ID and timestamp when absent (Vol. III Ch. 38).
func (m *Memory) CreateCapitalisedRentPrice(_ context.Context, crp rent.CapitalisedRentPrice) (rent.CapitalisedRentPrice, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if crp.ID == "" {
		crp.ID = rent.NewCapitalisedRentPriceID()
	}
	if _, exists := m.capitalisedRentPrices[crp.ID]; exists {
		return rent.CapitalisedRentPrice{}, ErrAlreadyExists
	}
	if crp.CreatedAt.IsZero() {
		crp.CreatedAt = m.now().UTC()
	}
	m.capitalisedRentPrices[crp.ID] = crp
	return crp, nil
}

// ListCapitalisedRentPrices returns all stored capitalised-rent-price records, newest first. Never returns nil.
func (m *Memory) ListCapitalisedRentPrices(_ context.Context) ([]rent.CapitalisedRentPrice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.CapitalisedRentPrice, 0, len(m.capitalisedRentPrices))
	for _, crp := range m.capitalisedRentPrices {
		out = append(out, crp)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateDifferentialRentITable stores tbl and its entries (Vol. III Ch. 39).
// Assigns IDs and timestamps when absent. Returns ErrAlreadyExists on table ID collision.
func (m *Memory) CreateDifferentialRentITable(_ context.Context, tbl rent.DifferentialRentITable, entries []rent.DifferentialRentIEntry) (rent.DifferentialRentITable, []rent.DifferentialRentIEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if tbl.ID == "" {
		tbl.ID = rent.NewDifferentialRentITableID()
	}
	if _, exists := m.dr1Tables[tbl.ID]; exists {
		return rent.DifferentialRentITable{}, nil, ErrAlreadyExists
	}
	if tbl.CreatedAt.IsZero() {
		tbl.CreatedAt = m.now().UTC()
	}
	filled := make([]rent.DifferentialRentIEntry, len(entries))
	for i, e := range entries {
		if e.ID == "" {
			e.ID = rent.NewDifferentialRentIEntryID()
		}
		if e.CreatedAt.IsZero() {
			e.CreatedAt = m.now().UTC()
		}
		e.TableID = string(tbl.ID)
		filled[i] = e
	}
	m.dr1Tables[tbl.ID] = tbl
	m.dr1Entries[tbl.ID] = filled
	return tbl, filled, nil
}

// GetDifferentialRentITable returns the table and its entries, or ErrNotFound (Vol. III Ch. 39).
func (m *Memory) GetDifferentialRentITable(_ context.Context, id rent.DifferentialRentITableID) (rent.DifferentialRentITable, []rent.DifferentialRentIEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tbl, ok := m.dr1Tables[id]
	if !ok {
		return rent.DifferentialRentITable{}, nil, ErrNotFound
	}
	entries := m.dr1Entries[id]
	if entries == nil {
		entries = []rent.DifferentialRentIEntry{}
	}
	return tbl, entries, nil
}

// ListDifferentialRentITables returns all table headers, newest first. Never returns nil (Vol. III Ch. 39).
func (m *Memory) ListDifferentialRentITables(_ context.Context) ([]rent.DifferentialRentITable, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.DifferentialRentITable, 0, len(m.dr1Tables))
	for _, tbl := range m.dr1Tables {
		out = append(out, tbl)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// ListDifferentialRentIEntries returns entries for tableID, or ErrNotFound if the table is unknown (Vol. III Ch. 39).
func (m *Memory) ListDifferentialRentIEntries(_ context.Context, tableID rent.DifferentialRentITableID) ([]rent.DifferentialRentIEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.dr1Tables[tableID]; !ok {
		return nil, ErrNotFound
	}
	entries := m.dr1Entries[tableID]
	if entries == nil {
		return []rent.DifferentialRentIEntry{}, nil
	}
	return entries, nil
}

// CreateLocationRentFactor stores f, assigning an ID and timestamp when absent (Vol. III Ch. 39).
func (m *Memory) CreateLocationRentFactor(_ context.Context, f rent.LocationRentFactor) (rent.LocationRentFactor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if f.ID == "" {
		f.ID = rent.NewLocationRentFactorID()
	}
	if _, exists := m.locationRentFactors[f.ID]; exists {
		return rent.LocationRentFactor{}, ErrAlreadyExists
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = m.now().UTC()
	}
	m.locationRentFactors[f.ID] = f
	return f, nil
}

// ListLocationRentFactors returns all stored location-rent-factor records, newest first. Never returns nil (Vol. III Ch. 39).
func (m *Memory) ListLocationRentFactors(_ context.Context) ([]rent.LocationRentFactor, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.LocationRentFactor, 0, len(m.locationRentFactors))
	for _, f := range m.locationRentFactors {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateSuccessiveInvestment stores inv, assigning an ID and timestamp when absent (Vol. III Ch. 40).
func (m *Memory) CreateSuccessiveInvestment(_ context.Context, inv rent.SuccessiveInvestment) (rent.SuccessiveInvestment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if inv.ID == "" {
		inv.ID = rent.NewSuccessiveInvestmentID()
	}
	if _, exists := m.successiveInvestments[inv.ID]; exists {
		return rent.SuccessiveInvestment{}, ErrAlreadyExists
	}
	if inv.CreatedAt.IsZero() {
		inv.CreatedAt = m.now().UTC()
	}
	m.successiveInvestments[inv.ID] = inv
	return inv, nil
}

// ListSuccessiveInvestments returns all stored successive-investment records, newest first. Never returns nil (Vol. III Ch. 40).
func (m *Memory) ListSuccessiveInvestments(_ context.Context) ([]rent.SuccessiveInvestment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.SuccessiveInvestment, 0, len(m.successiveInvestments))
	for _, inv := range m.successiveInvestments {
		out = append(out, inv)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// ListSuccessiveInvestmentsByParcel returns successive-investment records matching parcelID, newest first.
// Returns an empty slice (not ErrNotFound) when no rows match (Vol. III Ch. 40).
func (m *Memory) ListSuccessiveInvestmentsByParcel(_ context.Context, parcelID string) ([]rent.SuccessiveInvestment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.SuccessiveInvestment, 0)
	for _, inv := range m.successiveInvestments {
		if inv.ParcelID == parcelID {
			out = append(out, inv)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateLeaseTerm stores lt, assigning an ID and timestamp when absent (Vol. III Ch. 40).
func (m *Memory) CreateLeaseTerm(_ context.Context, lt rent.LeaseTerm) (rent.LeaseTerm, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if lt.ID == "" {
		lt.ID = rent.NewLeaseTermID()
	}
	if _, exists := m.leaseTerms[lt.ID]; exists {
		return rent.LeaseTerm{}, ErrAlreadyExists
	}
	if lt.CreatedAt.IsZero() {
		lt.CreatedAt = m.now().UTC()
	}
	m.leaseTerms[lt.ID] = lt
	return lt, nil
}

// ListLeaseTerms returns all stored lease-term records, newest first. Never returns nil (Vol. III Ch. 40).
func (m *Memory) ListLeaseTerms(_ context.Context) ([]rent.LeaseTerm, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.LeaseTerm, 0, len(m.leaseTerms))
	for _, lt := range m.leaseTerms {
		out = append(out, lt)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateDifferentialRentIIOutcome stores o, assigning an ID and timestamp when absent (Vol. III Ch. 41).
func (m *Memory) CreateDifferentialRentIIOutcome(_ context.Context, o rent.DifferentialRentIIOutcome) (rent.DifferentialRentIIOutcome, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if o.ID == "" {
		o.ID = rent.NewDifferentialRentIIOutcomeID()
	}
	if _, exists := m.dr2Outcomes[o.ID]; exists {
		return rent.DifferentialRentIIOutcome{}, ErrAlreadyExists
	}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = m.now().UTC()
	}
	m.dr2Outcomes[o.ID] = o
	return o, nil
}

// ListDifferentialRentIIOutcomes returns all stored outcomes, newest first. Never returns nil (Vol. III Ch. 41).
func (m *Memory) ListDifferentialRentIIOutcomes(_ context.Context) ([]rent.DifferentialRentIIOutcome, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.DifferentialRentIIOutcome, 0, len(m.dr2Outcomes))
	for _, o := range m.dr2Outcomes {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateIntensiveExtensiveComparison stores c, assigning an ID and timestamp when absent (Vol. III Ch. 41).
func (m *Memory) CreateIntensiveExtensiveComparison(_ context.Context, c rent.IntensiveExtensiveComparison) (rent.IntensiveExtensiveComparison, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c.ID == "" {
		c.ID = rent.NewIntensiveExtensiveComparisonID()
	}
	if _, exists := m.intensiveExtensive[c.ID]; exists {
		return rent.IntensiveExtensiveComparison{}, ErrAlreadyExists
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}
	m.intensiveExtensive[c.ID] = c
	return c, nil
}

// ListIntensiveExtensiveComparisons returns all stored comparisons, newest first. Never returns nil (Vol. III Ch. 41).
func (m *Memory) ListIntensiveExtensiveComparisons(_ context.Context) ([]rent.IntensiveExtensiveComparison, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.IntensiveExtensiveComparison, 0, len(m.intensiveExtensive))
	for _, c := range m.intensiveExtensive {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateSoilExclusionEvent stores e, assigning an ID and timestamp when absent (Vol. III Ch. 42).
func (m *Memory) CreateSoilExclusionEvent(_ context.Context, e rent.SoilExclusionEvent) (rent.SoilExclusionEvent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if e.ID == "" {
		e.ID = rent.NewSoilExclusionEventID()
	}
	if _, exists := m.soilExclusionEvents[e.ID]; exists {
		return rent.SoilExclusionEvent{}, ErrAlreadyExists
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = m.now().UTC()
	}
	m.soilExclusionEvents[e.ID] = e
	return e, nil
}

// GetSoilExclusionEvent returns the soil-exclusion event with id, or ErrNotFound (Vol. III Ch. 42).
func (m *Memory) GetSoilExclusionEvent(_ context.Context, id rent.SoilExclusionEventID) (rent.SoilExclusionEvent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	e, ok := m.soilExclusionEvents[id]
	if !ok {
		return rent.SoilExclusionEvent{}, ErrNotFound
	}
	return e, nil
}

// ListSoilExclusionEvents returns all stored soil-exclusion events, newest first. Never returns nil (Vol. III Ch. 42).
func (m *Memory) ListSoilExclusionEvents(_ context.Context) ([]rent.SoilExclusionEvent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.SoilExclusionEvent, 0, len(m.soilExclusionEvents))
	for _, e := range m.soilExclusionEvents {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateFallingPriceDR2Outcome stores o, assigning an ID and timestamp when absent (Vol. III Ch. 42).
func (m *Memory) CreateFallingPriceDR2Outcome(_ context.Context, o rent.FallingPriceDR2Outcome) (rent.FallingPriceDR2Outcome, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if o.ID == "" {
		o.ID = rent.NewFallingPriceDR2OutcomeID()
	}
	if _, exists := m.fallingPriceDR2Outcomes[o.ID]; exists {
		return rent.FallingPriceDR2Outcome{}, ErrAlreadyExists
	}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = m.now().UTC()
	}
	m.fallingPriceDR2Outcomes[o.ID] = o
	return o, nil
}

// ListFallingPriceDR2Outcomes returns all stored falling-price DR II outcomes, newest first. Never returns nil (Vol. III Ch. 42).
func (m *Memory) ListFallingPriceDR2Outcomes(_ context.Context) ([]rent.FallingPriceDR2Outcome, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.FallingPriceDR2Outcome, 0, len(m.fallingPriceDR2Outcomes))
	for _, o := range m.fallingPriceDR2Outcomes {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateNormalCapitalPerAcre stores n, assigning an ID and timestamp when absent (Vol. III Ch. 42).
func (m *Memory) CreateNormalCapitalPerAcre(_ context.Context, n rent.NormalCapitalPerAcre) (rent.NormalCapitalPerAcre, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if n.ID == "" {
		n.ID = rent.NewNormalCapitalPerAcreID()
	}
	if _, exists := m.normalCapitalPerAcre[n.ID]; exists {
		return rent.NormalCapitalPerAcre{}, ErrAlreadyExists
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = m.now().UTC()
	}
	m.normalCapitalPerAcre[n.ID] = n
	return n, nil
}

// ListNormalCapitalPerAcre returns all stored normal-capital-per-acre records, newest first. Never returns nil (Vol. III Ch. 42).
func (m *Memory) ListNormalCapitalPerAcre(_ context.Context) ([]rent.NormalCapitalPerAcre, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.NormalCapitalPerAcre, 0, len(m.normalCapitalPerAcre))
	for _, n := range m.normalCapitalPerAcre {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateRisingPriceDR2Outcome stores o, assigning an ID and timestamp when absent (Vol. III Ch. 43).
func (m *Memory) CreateRisingPriceDR2Outcome(_ context.Context, o rent.RisingPriceDR2Outcome) (rent.RisingPriceDR2Outcome, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if o.ID == "" {
		o.ID = rent.NewRisingPriceDR2OutcomeID()
	}
	if _, exists := m.risingPriceDR2Outcomes[o.ID]; exists {
		return rent.RisingPriceDR2Outcome{}, ErrAlreadyExists
	}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = m.now().UTC()
	}
	m.risingPriceDR2Outcomes[o.ID] = o
	return o, nil
}

// ListRisingPriceDR2Outcomes returns all stored rising-price DR II outcomes, newest first. Never returns nil (Vol. III Ch. 43).
func (m *Memory) ListRisingPriceDR2Outcomes(_ context.Context) ([]rent.RisingPriceDR2Outcome, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.RisingPriceDR2Outcome, 0, len(m.risingPriceDR2Outcomes))
	for _, o := range m.risingPriceDR2Outcomes {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateRentSequence stores rs, assigning an ID and timestamp when absent (Vol. III Ch. 43).
func (m *Memory) CreateRentSequence(_ context.Context, rs rent.RentSequence) (rent.RentSequence, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if rs.ID == "" {
		rs.ID = rent.NewRentSequenceID()
	}
	if _, exists := m.rentSequences[rs.ID]; exists {
		return rent.RentSequence{}, ErrAlreadyExists
	}
	if rs.CreatedAt.IsZero() {
		rs.CreatedAt = m.now().UTC()
	}
	m.rentSequences[rs.ID] = rs
	return rs, nil
}

// ListRentSequences returns all stored rent-sequence records, newest first. Never returns nil (Vol. III Ch. 43).
func (m *Memory) ListRentSequences(_ context.Context) ([]rent.RentSequence, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.RentSequence, 0, len(m.rentSequences))
	for _, rs := range m.rentSequences {
		out = append(out, rs)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateEngelsTableEntry stores e, assigning an ID and timestamp when absent (Vol. III Ch. 43).
func (m *Memory) CreateEngelsTableEntry(_ context.Context, e rent.EngelsTableEntry) (rent.EngelsTableEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if e.ID == "" {
		e.ID = rent.NewEngelsTableEntryID()
	}
	if _, exists := m.engelsTableEntries[e.ID]; exists {
		return rent.EngelsTableEntry{}, ErrAlreadyExists
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = m.now().UTC()
	}
	m.engelsTableEntries[e.ID] = e
	return e, nil
}

// ListEngelsTableEntries returns all stored Engels-table entries, newest first. Never returns nil (Vol. III Ch. 43).
func (m *Memory) ListEngelsTableEntries(_ context.Context) ([]rent.EngelsTableEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.EngelsTableEntry, 0, len(m.engelsTableEntries))
	for _, e := range m.engelsTableEntries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// ── Vol. III Ch. 44 — Differential Rent on the Worst Cultivated Soil ──

// CreateRentOnWorstSoil stores r, assigning an ID and timestamp when absent (Vol. III Ch. 44).
func (m *Memory) CreateRentOnWorstSoil(_ context.Context, r rent.RentOnWorstSoil) (rent.RentOnWorstSoil, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r.ID == "" {
		r.ID = rent.NewRentOnWorstSoilID()
	}
	if _, exists := m.rentOnWorstSoils[r.ID]; exists {
		return rent.RentOnWorstSoil{}, ErrAlreadyExists
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = m.now().UTC()
	}
	m.rentOnWorstSoils[r.ID] = r
	return r, nil
}

// ListRentOnWorstSoils returns all stored rent-on-worst-soil records, newest first. Never returns nil.
func (m *Memory) ListRentOnWorstSoils(_ context.Context) ([]rent.RentOnWorstSoil, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.RentOnWorstSoil, 0, len(m.rentOnWorstSoils))
	for _, r := range m.rentOnWorstSoils {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateSoilImprovementCapital stores s, assigning an ID and timestamp when absent (Vol. III Ch. 44).
func (m *Memory) CreateSoilImprovementCapital(_ context.Context, s rent.SoilImprovementCapital) (rent.SoilImprovementCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID == "" {
		s.ID = rent.NewSoilImprovementCapitalID()
	}
	if _, exists := m.soilImprovementCapitals[s.ID]; exists {
		return rent.SoilImprovementCapital{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	m.soilImprovementCapitals[s.ID] = s
	return s, nil
}

// ListSoilImprovementCapitals returns all stored soil-improvement-capital records, newest first. Never returns nil.
func (m *Memory) ListSoilImprovementCapitals(_ context.Context) ([]rent.SoilImprovementCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.SoilImprovementCapital, 0, len(m.soilImprovementCapitals))
	for _, s := range m.soilImprovementCapitals {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateRegulatingPriceShift stores s, assigning an ID and timestamp when absent (Vol. III Ch. 44).
func (m *Memory) CreateRegulatingPriceShift(_ context.Context, s rent.RegulatingPriceShift) (rent.RegulatingPriceShift, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID == "" {
		s.ID = rent.NewRegulatingPriceShiftID()
	}
	if _, exists := m.regulatingPriceShifts[s.ID]; exists {
		return rent.RegulatingPriceShift{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	m.regulatingPriceShifts[s.ID] = s
	return s, nil
}

// ListRegulatingPriceShifts returns all stored regulating-price-shift records, newest first. Never returns nil.
func (m *Memory) ListRegulatingPriceShifts(_ context.Context) ([]rent.RegulatingPriceShift, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.RegulatingPriceShift, 0, len(m.regulatingPriceShifts))
	for _, s := range m.regulatingPriceShifts {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// ── Vol. III Ch. 45 — Absolute Ground-Rent ──

// CreateAgriculturalComposition stores ac, assigning an ID and timestamp when absent (Vol. III Ch. 45).
func (m *Memory) CreateAgriculturalComposition(_ context.Context, ac rent.AgriculturalComposition) (rent.AgriculturalComposition, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ac.ID == "" {
		ac.ID = rent.NewAgriculturalCompositionID()
	}
	if _, exists := m.agriculturalCompositions[ac.ID]; exists {
		return rent.AgriculturalComposition{}, ErrAlreadyExists
	}
	if ac.CreatedAt.IsZero() {
		ac.CreatedAt = m.now().UTC()
	}
	m.agriculturalCompositions[ac.ID] = ac
	return ac, nil
}

// ListAgriculturalCompositions returns all stored records, newest first. Never returns nil.
func (m *Memory) ListAgriculturalCompositions(_ context.Context) ([]rent.AgriculturalComposition, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.AgriculturalComposition, 0, len(m.agriculturalCompositions))
	for _, ac := range m.agriculturalCompositions {
		out = append(out, ac)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateValuePriceGap stores g, assigning an ID and timestamp when absent (Vol. III Ch. 45).
func (m *Memory) CreateValuePriceGap(_ context.Context, g rent.ValuePriceGap) (rent.ValuePriceGap, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if g.ID == "" {
		g.ID = rent.NewValuePriceGapID()
	}
	if _, exists := m.valuePriceGaps[g.ID]; exists {
		return rent.ValuePriceGap{}, ErrAlreadyExists
	}
	if g.CreatedAt.IsZero() {
		g.CreatedAt = m.now().UTC()
	}
	m.valuePriceGaps[g.ID] = g
	return g, nil
}

// ListValuePriceGaps returns all stored records, newest first. Never returns nil.
func (m *Memory) ListValuePriceGaps(_ context.Context) ([]rent.ValuePriceGap, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.ValuePriceGap, 0, len(m.valuePriceGaps))
	for _, g := range m.valuePriceGaps {
		out = append(out, g)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateAbsoluteRent stores ar, assigning an ID and timestamp when absent (Vol. III Ch. 45).
func (m *Memory) CreateAbsoluteRent(_ context.Context, ar rent.AbsoluteRent) (rent.AbsoluteRent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ar.ID == "" {
		ar.ID = rent.NewAbsoluteRentID()
	}
	if _, exists := m.absoluteRents[ar.ID]; exists {
		return rent.AbsoluteRent{}, ErrAlreadyExists
	}
	if ar.CreatedAt.IsZero() {
		ar.CreatedAt = m.now().UTC()
	}
	m.absoluteRents[ar.ID] = ar
	return ar, nil
}

// GetAbsoluteRent returns the record with id, or ErrNotFound (Vol. III Ch. 45).
func (m *Memory) GetAbsoluteRent(_ context.Context, id rent.AbsoluteRentID) (rent.AbsoluteRent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ar, ok := m.absoluteRents[id]
	if !ok {
		return rent.AbsoluteRent{}, ErrNotFound
	}
	return ar, nil
}

// ListAbsoluteRents returns all stored records, newest first. Never returns nil.
func (m *Memory) ListAbsoluteRents(_ context.Context) ([]rent.AbsoluteRent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.AbsoluteRent, 0, len(m.absoluteRents))
	for _, ar := range m.absoluteRents {
		out = append(out, ar)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateAbsoluteRentLimit stores l, assigning an ID and timestamp when absent (Vol. III Ch. 45).
func (m *Memory) CreateAbsoluteRentLimit(_ context.Context, l rent.AbsoluteRentLimit) (rent.AbsoluteRentLimit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if l.ID == "" {
		l.ID = rent.NewAbsoluteRentLimitID()
	}
	if _, exists := m.absoluteRentLimits[l.ID]; exists {
		return rent.AbsoluteRentLimit{}, ErrAlreadyExists
	}
	if l.CreatedAt.IsZero() {
		l.CreatedAt = m.now().UTC()
	}
	m.absoluteRentLimits[l.ID] = l
	return l, nil
}

// ListAbsoluteRentLimits returns all stored records, newest first. Never returns nil.
func (m *Memory) ListAbsoluteRentLimits(_ context.Context) ([]rent.AbsoluteRentLimit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.AbsoluteRentLimit, 0, len(m.absoluteRentLimits))
	for _, l := range m.absoluteRentLimits {
		out = append(out, l)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateBuildingSiteRent stores r, assigning an ID and timestamp when absent (Vol. III Ch. 46).
func (m *Memory) CreateBuildingSiteRent(_ context.Context, r rent.BuildingSiteRent) (rent.BuildingSiteRent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r.ID == "" {
		r.ID = rent.NewBuildingSiteRentID()
	}
	if _, exists := m.buildingSiteRents[r.ID]; exists {
		return rent.BuildingSiteRent{}, ErrAlreadyExists
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = m.now().UTC()
	}
	m.buildingSiteRents[r.ID] = r
	return r, nil
}

// ListBuildingSiteRents returns all stored records, newest first. Never returns nil.
func (m *Memory) ListBuildingSiteRents(_ context.Context) ([]rent.BuildingSiteRent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.BuildingSiteRent, 0, len(m.buildingSiteRents))
	for _, r := range m.buildingSiteRents {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateMiningRent stores r, assigning an ID and timestamp when absent (Vol. III Ch. 46).
func (m *Memory) CreateMiningRent(_ context.Context, r rent.MiningRent) (rent.MiningRent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r.ID == "" {
		r.ID = rent.NewMiningRentID()
	}
	if _, exists := m.miningRents[r.ID]; exists {
		return rent.MiningRent{}, ErrAlreadyExists
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = m.now().UTC()
	}
	m.miningRents[r.ID] = r
	return r, nil
}

// ListMiningRents returns all stored records, newest first. Never returns nil.
func (m *Memory) ListMiningRents(_ context.Context) ([]rent.MiningRent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.MiningRent, 0, len(m.miningRents))
	for _, r := range m.miningRents {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateMonopolyPriceRent stores r, assigning an ID and timestamp when absent (Vol. III Ch. 46).
func (m *Memory) CreateMonopolyPriceRent(_ context.Context, r rent.MonopolyPriceRent) (rent.MonopolyPriceRent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r.ID == "" {
		r.ID = rent.NewMonopolyPriceRentID()
	}
	if _, exists := m.monopolyPriceRents[r.ID]; exists {
		return rent.MonopolyPriceRent{}, ErrAlreadyExists
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = m.now().UTC()
	}
	m.monopolyPriceRents[r.ID] = r
	return r, nil
}

// ListMonopolyPriceRents returns all stored records, newest first. Never returns nil.
func (m *Memory) ListMonopolyPriceRents(_ context.Context) ([]rent.MonopolyPriceRent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.MonopolyPriceRent, 0, len(m.monopolyPriceRents))
	for _, r := range m.monopolyPriceRents {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateLandPriceScenario stores s, assigning an ID and timestamp when absent (Vol. III Ch. 46).
func (m *Memory) CreateLandPriceScenario(_ context.Context, s rent.LandPriceScenario) (rent.LandPriceScenario, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID == "" {
		s.ID = rent.NewLandPriceScenarioID()
	}
	if _, exists := m.landPriceScenarios[s.ID]; exists {
		return rent.LandPriceScenario{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	m.landPriceScenarios[s.ID] = s
	return s, nil
}

// ListLandPriceScenarios returns all stored records, newest first. Never returns nil.
func (m *Memory) ListLandPriceScenarios(_ context.Context) ([]rent.LandPriceScenario, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.LandPriceScenario, 0, len(m.landPriceScenarios))
	for _, s := range m.landPriceScenarios {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateRentHistoricalForm stores f, assigning an ID and timestamp when absent (Vol. III Ch. 47).
func (m *Memory) CreateRentHistoricalForm(_ context.Context, f rent.RentHistoricalForm) (rent.RentHistoricalForm, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if f.ID == "" {
		f.ID = rent.NewRentHistoricalFormID()
	}
	if _, exists := m.rentHistoricalForms[f.ID]; exists {
		return rent.RentHistoricalForm{}, ErrAlreadyExists
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}
	m.rentHistoricalForms[f.ID] = f
	return f, nil
}

// ListRentHistoricalForms returns all stored records, newest first. Never returns nil.
func (m *Memory) ListRentHistoricalForms(_ context.Context) ([]rent.RentHistoricalForm, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.RentHistoricalForm, 0, len(m.rentHistoricalForms))
	for _, f := range m.rentHistoricalForms {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return string(out[i].ID) < string(out[j].ID)
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateHistoricalRentTransition stores t, assigning an ID and timestamp when absent (Vol. III Ch. 47).
func (m *Memory) CreateHistoricalRentTransition(_ context.Context, t rent.HistoricalRentTransition) (rent.HistoricalRentTransition, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t.ID == "" {
		t.ID = rent.NewHistoricalRentTransitionID()
	}
	if _, exists := m.historicalRentTransitions[t.ID]; exists {
		return rent.HistoricalRentTransition{}, ErrAlreadyExists
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
	m.historicalRentTransitions[t.ID] = t
	return t, nil
}

// ListHistoricalRentTransitions returns all stored records, newest first. Never returns nil.
func (m *Memory) ListHistoricalRentTransitions(_ context.Context) ([]rent.HistoricalRentTransition, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.HistoricalRentTransition, 0, len(m.historicalRentTransitions))
	for _, t := range m.historicalRentTransitions {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return string(out[i].ID) < string(out[j].ID)
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateSmallPeasantProduction stores p, assigning an ID and timestamp when absent (Vol. III Ch. 47).
func (m *Memory) CreateSmallPeasantProduction(_ context.Context, p rent.SmallPeasantProduction) (rent.SmallPeasantProduction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p.ID == "" {
		p.ID = rent.NewSmallPeasantProductionID()
	}
	if _, exists := m.smallPeasantProductions[p.ID]; exists {
		return rent.SmallPeasantProduction{}, ErrAlreadyExists
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	m.smallPeasantProductions[p.ID] = p
	return p, nil
}

// ListSmallPeasantProductions returns all stored records, newest first. Never returns nil.
func (m *Memory) ListSmallPeasantProductions(_ context.Context) ([]rent.SmallPeasantProduction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]rent.SmallPeasantProduction, 0, len(m.smallPeasantProductions))
	for _, p := range m.smallPeasantProductions {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return string(out[i].ID) < string(out[j].ID)
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}
