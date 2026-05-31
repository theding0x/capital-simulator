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
	profitDivisions             map[credit.ProfitDivisionID]credit.ProfitDivision
	compoundInterestSchedules  map[credit.CompoundInterestID]credit.CompoundInterestSchedule
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
