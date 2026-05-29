package store

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// Memory is an in-memory Store for unit tests and local development.
type Memory struct {
	mu          sync.RWMutex
	now         func() time.Time
	costPrices  map[profit.CostPriceID]profit.CostPrice
	profitRates map[profit.ProfitRateID]profit.ProfitRateAnalysis
	variations  map[profit.VariationAnalysisID]profit.VariationAnalysis
	turnovers   map[profit.TurnoverAnalysisID]profit.TurnoverAnalysis
}

// NewMemory returns an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{
		now:         time.Now,
		costPrices:  make(map[profit.CostPriceID]profit.CostPrice),
		profitRates: make(map[profit.ProfitRateID]profit.ProfitRateAnalysis),
		variations:  make(map[profit.VariationAnalysisID]profit.VariationAnalysis),
		turnovers:   make(map[profit.TurnoverAnalysisID]profit.TurnoverAnalysis),
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
