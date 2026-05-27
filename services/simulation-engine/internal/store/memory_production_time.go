package store

import (
	"context"
	"sync"
	"time"

	wp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

// memoryProductionTime is the in-memory backing store for Ch. 13 types.
type memoryProductionTime struct {
	mu          sync.RWMutex
	gaps        map[wp.WorkingPeriodID]wp.ProductionLabourGap
	activations map[wp.WorkingPeriodID][]wp.NaturalProcessActivation
	subjects    map[wp.WorkingPeriodID]wp.NaturalSubject
	industries  map[wp.NaturalProcessIndustryID]wp.NaturalProcessIndustry
}

func newMemoryProductionTime() *memoryProductionTime {
	return &memoryProductionTime{
		gaps:        make(map[wp.WorkingPeriodID]wp.ProductionLabourGap),
		activations: make(map[wp.WorkingPeriodID][]wp.NaturalProcessActivation),
		subjects:    make(map[wp.WorkingPeriodID]wp.NaturalSubject),
		industries:  make(map[wp.NaturalProcessIndustryID]wp.NaturalProcessIndustry),
	}
}

// RecordProductionLabourGap persists a gap record (one per working period).
func (m *Memory) RecordProductionLabourGap(_ context.Context, gap wp.ProductionLabourGap) (wp.ProductionLabourGap, error) {
	if err := gap.Validate(); err != nil {
		return wp.ProductionLabourGap{}, err
	}
	m.productionTime.mu.Lock()
	defer m.productionTime.mu.Unlock()
	if gap.ID == "" {
		gap.ID = wp.NewProductionLabourGapID()
	}
	if gap.CreatedAt.IsZero() {
		gap.CreatedAt = time.Now().UTC()
	}
	if _, ok := m.productionTime.gaps[gap.WorkingPeriodID]; ok {
		return wp.ProductionLabourGap{}, ErrAlreadyExists
	}
	m.productionTime.gaps[gap.WorkingPeriodID] = gap
	return gap, nil
}

// GetProductionLabourGap returns the gap for a given working period.
func (m *Memory) GetProductionLabourGap(_ context.Context, wpID wp.WorkingPeriodID) (wp.ProductionLabourGap, error) {
	m.productionTime.mu.RLock()
	defer m.productionTime.mu.RUnlock()
	g, ok := m.productionTime.gaps[wpID]
	if !ok {
		return wp.ProductionLabourGap{}, ErrNotFound
	}
	return g, nil
}

// RecordNaturalProcessActivation appends an activation for a working period.
func (m *Memory) RecordNaturalProcessActivation(_ context.Context, act wp.NaturalProcessActivation) (wp.NaturalProcessActivation, error) {
	if err := act.Validate(); err != nil {
		return wp.NaturalProcessActivation{}, err
	}
	m.productionTime.mu.Lock()
	defer m.productionTime.mu.Unlock()
	if act.ID == "" {
		act.ID = wp.NewNaturalProcessActivationID()
	}
	if act.StartedAt.IsZero() {
		act.StartedAt = time.Now().UTC()
	}
	m.productionTime.activations[act.WorkingPeriodID] = append(
		m.productionTime.activations[act.WorkingPeriodID], act)
	return act, nil
}

// ListNaturalProcessActivations returns all activations for a working period.
func (m *Memory) ListNaturalProcessActivations(_ context.Context, wpID wp.WorkingPeriodID) ([]wp.NaturalProcessActivation, error) {
	m.productionTime.mu.RLock()
	defer m.productionTime.mu.RUnlock()
	acts := m.productionTime.activations[wpID]
	if acts == nil {
		return []wp.NaturalProcessActivation{}, nil
	}
	out := make([]wp.NaturalProcessActivation, len(acts))
	copy(out, acts)
	return out, nil
}

// RecordNaturalSubject persists a subject description (one per working period).
func (m *Memory) RecordNaturalSubject(_ context.Context, sub wp.NaturalSubject) (wp.NaturalSubject, error) {
	if err := sub.Validate(); err != nil {
		return wp.NaturalSubject{}, err
	}
	m.productionTime.mu.Lock()
	defer m.productionTime.mu.Unlock()
	if sub.ID == "" {
		sub.ID = wp.NewNaturalSubjectID()
	}
	if sub.CreatedAt.IsZero() {
		sub.CreatedAt = time.Now().UTC()
	}
	if _, ok := m.productionTime.subjects[sub.WorkingPeriodID]; ok {
		return wp.NaturalSubject{}, ErrAlreadyExists
	}
	m.productionTime.subjects[sub.WorkingPeriodID] = sub
	return sub, nil
}

// GetNaturalSubject returns the subject for a given working period.
func (m *Memory) GetNaturalSubject(_ context.Context, wpID wp.WorkingPeriodID) (wp.NaturalSubject, error) {
	m.productionTime.mu.RLock()
	defer m.productionTime.mu.RUnlock()
	s, ok := m.productionTime.subjects[wpID]
	if !ok {
		return wp.NaturalSubject{}, ErrNotFound
	}
	return s, nil
}

// CreateIndustryBenchmark persists a new NaturalProcessIndustry benchmark.
func (m *Memory) CreateIndustryBenchmark(_ context.Context, ind wp.NaturalProcessIndustry) (wp.NaturalProcessIndustry, error) {
	if err := ind.Validate(); err != nil {
		return wp.NaturalProcessIndustry{}, err
	}
	m.productionTime.mu.Lock()
	defer m.productionTime.mu.Unlock()
	if ind.ID == "" {
		ind.ID = wp.NewNaturalProcessIndustryID()
	}
	if ind.CreatedAt.IsZero() {
		ind.CreatedAt = time.Now().UTC()
	}
	if _, ok := m.productionTime.industries[ind.ID]; ok {
		return wp.NaturalProcessIndustry{}, ErrAlreadyExists
	}
	m.productionTime.industries[ind.ID] = ind
	return ind, nil
}

// ListIndustryBenchmarks returns all industry benchmarks.
func (m *Memory) ListIndustryBenchmarks(_ context.Context) ([]wp.NaturalProcessIndustry, error) {
	m.productionTime.mu.RLock()
	defer m.productionTime.mu.RUnlock()
	out := make([]wp.NaturalProcessIndustry, 0, len(m.productionTime.industries))
	for _, ind := range m.productionTime.industries {
		out = append(out, ind)
	}
	return out, nil
}
