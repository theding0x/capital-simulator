package store

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// Memory is an in-memory Store for unit tests and local development.
type Memory struct {
	mu           sync.RWMutex
	now          func() time.Time
	costPrices   map[profit.CostPriceID]profit.CostPrice
	profitRates  map[profit.ProfitRateID]profit.ProfitRateAnalysis
	variations   map[profit.VariationAnalysisID]profit.VariationAnalysis
	turnovers    map[profit.TurnoverAnalysisID]profit.TurnoverAnalysis
	economies    map[profit.EconomyAnalysisID]profit.EconomyAnalysis
	priceFlux    map[profit.PriceFluctuationAnalysisID]profit.PriceFluctuationAnalysis
	compositions map[profit.CompositionEffectID]profit.CompositionEffectAnalysis
	magnitudes   map[profit.MagnitudeChangeID]profit.MagnitudeChangeAnalysis
	spheres      map[avgprofit.ProductionSphereID]avgprofit.ProductionSphere
}

// NewMemory returns an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{
		now:          time.Now,
		costPrices:   make(map[profit.CostPriceID]profit.CostPrice),
		profitRates:  make(map[profit.ProfitRateID]profit.ProfitRateAnalysis),
		variations:   make(map[profit.VariationAnalysisID]profit.VariationAnalysis),
		turnovers:    make(map[profit.TurnoverAnalysisID]profit.TurnoverAnalysis),
		economies:    make(map[profit.EconomyAnalysisID]profit.EconomyAnalysis),
		priceFlux:    make(map[profit.PriceFluctuationAnalysisID]profit.PriceFluctuationAnalysis),
		compositions: make(map[profit.CompositionEffectID]profit.CompositionEffectAnalysis),
		magnitudes:   make(map[profit.MagnitudeChangeID]profit.MagnitudeChangeAnalysis),
		spheres:      make(map[avgprofit.ProductionSphereID]avgprofit.ProductionSphere),
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
