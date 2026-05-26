package store

import (
	"context"
	"sync"
	"time"

	wp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

type memoryWorkingPeriod struct {
	mu              sync.RWMutex
	periods         map[wp.WorkingPeriodID]wp.WorkingPeriod
	interruptions   map[wp.WorkingPeriodID][]wp.WorkingPeriodInterruption
	shortenings     map[wp.WorkingPeriodID][]wp.WorkingPeriodShortening
	creditFinancing map[wp.WorkingPeriodID]wp.WorkingPeriodCreditFinancing
	constraints     map[string]wp.NaturalWorkingPeriodConstraint // keyed by commodity_id
}

func newMemoryWorkingPeriod() *memoryWorkingPeriod {
	return &memoryWorkingPeriod{
		periods:         make(map[wp.WorkingPeriodID]wp.WorkingPeriod),
		interruptions:   make(map[wp.WorkingPeriodID][]wp.WorkingPeriodInterruption),
		shortenings:     make(map[wp.WorkingPeriodID][]wp.WorkingPeriodShortening),
		creditFinancing: make(map[wp.WorkingPeriodID]wp.WorkingPeriodCreditFinancing),
		constraints:     make(map[string]wp.NaturalWorkingPeriodConstraint),
	}
}

func (m *Memory) CreateWorkingPeriod(_ context.Context, w wp.WorkingPeriod) (wp.WorkingPeriod, error) {
	if err := w.Validate(); err != nil {
		return wp.WorkingPeriod{}, err
	}
	m.workingPeriod.mu.Lock()
	defer m.workingPeriod.mu.Unlock()
	// Enforce natural constraint floor.
	if nc, ok := m.workingPeriod.constraints[w.CommodityID]; ok {
		if w.WorkingDayCount < nc.MinWorkingDays {
			return wp.WorkingPeriod{}, wp.ErrPrematureDelivery
		}
	}
	if w.ID == "" {
		w.ID = wp.NewWorkingPeriodID()
	}
	if w.CreatedAt.IsZero() {
		w.CreatedAt = time.Now().UTC()
	}
	if _, ok := m.workingPeriod.periods[w.ID]; ok {
		return wp.WorkingPeriod{}, ErrAlreadyExists
	}
	m.workingPeriod.periods[w.ID] = w
	return w, nil
}

func (m *Memory) GetWorkingPeriod(_ context.Context, id wp.WorkingPeriodID) (wp.WorkingPeriod, error) {
	m.workingPeriod.mu.RLock()
	defer m.workingPeriod.mu.RUnlock()
	w, ok := m.workingPeriod.periods[id]
	if !ok {
		return wp.WorkingPeriod{}, ErrNotFound
	}
	w.Interruptions = m.workingPeriod.interruptions[id]
	if w.Interruptions == nil {
		w.Interruptions = []wp.WorkingPeriodInterruption{}
	}
	w.Shortenings = m.workingPeriod.shortenings[id]
	if w.Shortenings == nil {
		w.Shortenings = []wp.WorkingPeriodShortening{}
	}
	if cf, ok := m.workingPeriod.creditFinancing[id]; ok {
		w.CreditFinancing = &cf
	}
	return w, nil
}

func (m *Memory) ListWorkingPeriods(_ context.Context, industrialCapitalID, mode, commodityID string) ([]wp.WorkingPeriod, error) {
	m.workingPeriod.mu.RLock()
	defer m.workingPeriod.mu.RUnlock()
	out := []wp.WorkingPeriod{}
	for _, w := range m.workingPeriod.periods {
		if industrialCapitalID != "" && w.IndustrialCapitalID != industrialCapitalID {
			continue
		}
		if mode != "" && string(w.Mode) != mode {
			continue
		}
		if commodityID != "" && w.CommodityID != commodityID {
			continue
		}
		out = append(out, w)
	}
	return out, nil
}

func (m *Memory) RecordInterruption(_ context.Context, id wp.WorkingPeriodID, intr wp.WorkingPeriodInterruption) (wp.WorkingPeriodInterruption, error) {
	m.workingPeriod.mu.Lock()
	defer m.workingPeriod.mu.Unlock()
	w, ok := m.workingPeriod.periods[id]
	if !ok {
		return wp.WorkingPeriodInterruption{}, ErrNotFound
	}
	if err := intr.Validate(w.Mode); err != nil {
		return wp.WorkingPeriodInterruption{}, err
	}
	if intr.ID == "" {
		intr.ID = wp.NewWorkingPeriodInterruptionID()
	}
	intr.WorkingPeriodID = id
	if intr.InterruptedAt.IsZero() {
		intr.InterruptedAt = time.Now().UTC()
	}
	m.workingPeriod.interruptions[id] = append(m.workingPeriod.interruptions[id], intr)
	return intr, nil
}

func (m *Memory) RecordShortening(_ context.Context, id wp.WorkingPeriodID, s wp.WorkingPeriodShortening) (wp.WorkingPeriodShortening, error) {
	m.workingPeriod.mu.Lock()
	defer m.workingPeriod.mu.Unlock()
	if _, ok := m.workingPeriod.periods[id]; !ok {
		return wp.WorkingPeriodShortening{}, ErrNotFound
	}
	if err := s.Validate(); err != nil {
		return wp.WorkingPeriodShortening{}, err
	}
	if s.ID == "" {
		s.ID = wp.NewWorkingPeriodShorteningID()
	}
	s.WorkingPeriodID = id
	if s.OccurredAt.IsZero() {
		s.OccurredAt = time.Now().UTC()
	}
	m.workingPeriod.shortenings[id] = append(m.workingPeriod.shortenings[id], s)
	return s, nil
}

func (m *Memory) RecordCreditFinancing(_ context.Context, id wp.WorkingPeriodID, cf wp.WorkingPeriodCreditFinancing) (wp.WorkingPeriodCreditFinancing, error) {
	m.workingPeriod.mu.Lock()
	defer m.workingPeriod.mu.Unlock()
	if _, ok := m.workingPeriod.periods[id]; !ok {
		return wp.WorkingPeriodCreditFinancing{}, ErrNotFound
	}
	if err := cf.Validate(); err != nil {
		return wp.WorkingPeriodCreditFinancing{}, err
	}
	if cf.ID == "" {
		cf.ID = wp.NewWorkingPeriodCreditFinancingID()
	}
	cf.WorkingPeriodID = id
	if cf.CreatedAt.IsZero() {
		cf.CreatedAt = time.Now().UTC()
	}
	m.workingPeriod.creditFinancing[id] = cf
	return cf, nil
}

func (m *Memory) CreateNaturalConstraint(_ context.Context, nc wp.NaturalWorkingPeriodConstraint) (wp.NaturalWorkingPeriodConstraint, error) {
	if err := nc.Validate(); err != nil {
		return wp.NaturalWorkingPeriodConstraint{}, err
	}
	m.workingPeriod.mu.Lock()
	defer m.workingPeriod.mu.Unlock()
	if _, ok := m.workingPeriod.constraints[nc.CommodityID]; ok {
		return wp.NaturalWorkingPeriodConstraint{}, ErrAlreadyExists
	}
	if nc.ID == "" {
		nc.ID = wp.NewNaturalWorkingPeriodConstraintID()
	}
	if nc.CreatedAt.IsZero() {
		nc.CreatedAt = time.Now().UTC()
	}
	m.workingPeriod.constraints[nc.CommodityID] = nc
	return nc, nil
}

func (m *Memory) ListNaturalConstraints(_ context.Context) ([]wp.NaturalWorkingPeriodConstraint, error) {
	m.workingPeriod.mu.RLock()
	defer m.workingPeriod.mu.RUnlock()
	out := []wp.NaturalWorkingPeriodConstraint{}
	for _, nc := range m.workingPeriod.constraints {
		out = append(out, nc)
	}
	return out, nil
}
