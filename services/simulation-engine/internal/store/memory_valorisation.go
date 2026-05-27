package store

import (
	"context"
	"sync"

	val "github.com/theding0x/capital-simulator/services/simulation-engine/internal/valorisation"
)

// memoryValorisation is the in-memory backing store for Vol. II Ch. 16 types.
type memoryValorisation struct {
	mu              sync.RWMutex
	advances        map[val.VariableCapitalAdvanceID]val.VariableCapitalAdvance
	perCircuitRates map[val.VariableCapitalAdvanceID]val.PerCircuitSurplusRate
	annualRates     []val.AnnualSurplusRate
	contrasts       []val.AnnualSurplusRateContrast
}

func newMemoryValorisation() *memoryValorisation {
	return &memoryValorisation{
		advances:        make(map[val.VariableCapitalAdvanceID]val.VariableCapitalAdvance),
		perCircuitRates: make(map[val.VariableCapitalAdvanceID]val.PerCircuitSurplusRate),
	}
}

// CreateAdvance persists a VariableCapitalAdvance after validation.
func (m *Memory) CreateAdvance(_ context.Context, a val.VariableCapitalAdvance) (val.VariableCapitalAdvance, error) {
	if err := a.Validate(); err != nil {
		return val.VariableCapitalAdvance{}, err
	}
	if a.ID == "" {
		a.ID = val.NewVariableCapitalAdvanceID()
	}
	m.valorisation.mu.Lock()
	defer m.valorisation.mu.Unlock()
	if _, exists := m.valorisation.advances[a.ID]; exists {
		return val.VariableCapitalAdvance{}, ErrAlreadyExists
	}
	m.valorisation.advances[a.ID] = a
	return a, nil
}

// GetAdvance returns a VariableCapitalAdvance by ID.
func (m *Memory) GetAdvance(_ context.Context, id val.VariableCapitalAdvanceID) (val.VariableCapitalAdvance, error) {
	m.valorisation.mu.RLock()
	defer m.valorisation.mu.RUnlock()
	a, ok := m.valorisation.advances[id]
	if !ok {
		return val.VariableCapitalAdvance{}, ErrNotFound
	}
	return a, nil
}

// ListAdvances returns all advances for an industrial capital.
func (m *Memory) ListAdvances(_ context.Context, industrialCapitalID string) ([]val.VariableCapitalAdvance, error) {
	m.valorisation.mu.RLock()
	defer m.valorisation.mu.RUnlock()
	var out []val.VariableCapitalAdvance
	for _, a := range m.valorisation.advances {
		if industrialCapitalID == "" || a.IndustrialCapitalID == industrialCapitalID {
			out = append(out, a)
		}
	}
	if out == nil {
		out = []val.VariableCapitalAdvance{}
	}
	return out, nil
}

// RecordPerCircuitRate persists (or replaces) the per-circuit surplus rate.
func (m *Memory) RecordPerCircuitRate(_ context.Context, rate val.PerCircuitSurplusRate) (val.PerCircuitSurplusRate, error) {
	m.valorisation.mu.Lock()
	defer m.valorisation.mu.Unlock()
	if _, ok := m.valorisation.advances[rate.VariableCapitalAdvanceID]; !ok {
		return val.PerCircuitSurplusRate{}, ErrNotFound
	}
	if rate.ID == "" {
		rate.ID = val.NewPerCircuitSurplusRateID()
	}
	m.valorisation.perCircuitRates[rate.VariableCapitalAdvanceID] = rate
	return rate, nil
}

// GetPerCircuitRate returns the per-circuit surplus rate for an advance.
func (m *Memory) GetPerCircuitRate(_ context.Context, advanceID val.VariableCapitalAdvanceID) (val.PerCircuitSurplusRate, error) {
	m.valorisation.mu.RLock()
	defer m.valorisation.mu.RUnlock()
	r, ok := m.valorisation.perCircuitRates[advanceID]
	if !ok {
		return val.PerCircuitSurplusRate{}, ErrNotFound
	}
	return r, nil
}

// RecordAnnualRate persists a computed AnnualSurplusRate snapshot.
func (m *Memory) RecordAnnualRate(_ context.Context, rate val.AnnualSurplusRate) (val.AnnualSurplusRate, error) {
	m.valorisation.mu.Lock()
	defer m.valorisation.mu.Unlock()
	if rate.ID == "" {
		rate.ID = val.NewAnnualSurplusRateID()
	}
	m.valorisation.annualRates = append(m.valorisation.annualRates, rate)
	return rate, nil
}

// ListAnnualRates returns annual-rate snapshots, optionally filtered.
func (m *Memory) ListAnnualRates(_ context.Context, industrialCapitalID, period string) ([]val.AnnualSurplusRate, error) {
	m.valorisation.mu.RLock()
	defer m.valorisation.mu.RUnlock()

	// Build a set of advance IDs for the requested IC (if filtering by IC).
	var filterByAdvance map[val.VariableCapitalAdvanceID]bool
	if industrialCapitalID != "" {
		filterByAdvance = make(map[val.VariableCapitalAdvanceID]bool)
		for _, a := range m.valorisation.advances {
			if a.IndustrialCapitalID == industrialCapitalID {
				filterByAdvance[a.ID] = true
			}
		}
	}

	var out []val.AnnualSurplusRate
	for _, r := range m.valorisation.annualRates {
		if filterByAdvance != nil && !filterByAdvance[r.VariableCapitalAdvanceID] {
			continue
		}
		if period != "" && r.Period != period {
			continue
		}
		out = append(out, r)
	}
	if out == nil {
		out = []val.AnnualSurplusRate{}
	}
	return out, nil
}

// RecordContrast persists a canonical AnnualSurplusRateContrast.
func (m *Memory) RecordContrast(_ context.Context, c val.AnnualSurplusRateContrast) (val.AnnualSurplusRateContrast, error) {
	m.valorisation.mu.Lock()
	defer m.valorisation.mu.Unlock()
	if c.ID == "" {
		c.ID = val.NewAnnualSurplusRateContrastID()
	}
	m.valorisation.contrasts = append(m.valorisation.contrasts, c)
	return c, nil
}
