package store

import (
	"context"
	"sync"
	"time"

	tv "github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

type memoryAggregateTurnover struct {
	mu             sync.RWMutex
	aggregates     map[tv.AggregateTurnoverID]tv.AggregateTurnover
	contributions  map[tv.AggregateTurnoverID][]tv.ComponentTurnoverContribution
	lifetimeCycles map[tv.AggregateTurnoverID]tv.LifetimeCycle
	crisisPhases   map[tv.LifetimeCycleID][]tv.CrisisCyclePhaseRecord
}

func newMemoryAggregateTurnover() *memoryAggregateTurnover {
	return &memoryAggregateTurnover{
		aggregates:     make(map[tv.AggregateTurnoverID]tv.AggregateTurnover),
		contributions:  make(map[tv.AggregateTurnoverID][]tv.ComponentTurnoverContribution),
		lifetimeCycles: make(map[tv.AggregateTurnoverID]tv.LifetimeCycle),
		crisisPhases:   make(map[tv.LifetimeCycleID][]tv.CrisisCyclePhaseRecord),
	}
}

func (m *Memory) CreateAggregateTurnover(_ context.Context, at tv.AggregateTurnover) (tv.AggregateTurnover, error) {
	if err := at.Validate(); err != nil {
		return tv.AggregateTurnover{}, err
	}
	if at.ID == "" {
		at.ID = tv.NewAggregateTurnoverID()
	}
	m.aggregateTurnovers.mu.Lock()
	defer m.aggregateTurnovers.mu.Unlock()
	if _, ok := m.aggregateTurnovers.aggregates[at.ID]; ok {
		return tv.AggregateTurnover{}, ErrAlreadyExists
	}
	at.Contributions = []tv.ComponentTurnoverContribution{}
	m.aggregateTurnovers.aggregates[at.ID] = at
	m.aggregateTurnovers.contributions[at.ID] = []tv.ComponentTurnoverContribution{}
	return at, nil
}

func (m *Memory) GetAggregateTurnover(_ context.Context, id tv.AggregateTurnoverID) (tv.AggregateTurnover, error) {
	m.aggregateTurnovers.mu.RLock()
	defer m.aggregateTurnovers.mu.RUnlock()
	at, ok := m.aggregateTurnovers.aggregates[id]
	if !ok {
		return tv.AggregateTurnover{}, ErrNotFound
	}
	// Attach contributions.
	contribs := make([]tv.ComponentTurnoverContribution, len(m.aggregateTurnovers.contributions[id]))
	copy(contribs, m.aggregateTurnovers.contributions[id])
	at.Contributions = contribs
	// Attach lifetime cycle and crisis phases if present.
	if lc, ok := m.aggregateTurnovers.lifetimeCycles[id]; ok {
		phases := make([]tv.CrisisCyclePhaseRecord, len(m.aggregateTurnovers.crisisPhases[lc.ID]))
		copy(phases, m.aggregateTurnovers.crisisPhases[lc.ID])
		lc.CrisisPhases = phases
		lcCopy := lc
		at.LifetimeCycle = &lcCopy
		if len(phases) > 0 {
			latest := phases[len(phases)-1]
			at.LatestCrisisPhase = &latest
		}
	}
	return at, nil
}

func (m *Memory) ListAggregateTurnovers(_ context.Context, industrialCapitalID string) ([]tv.AggregateTurnover, error) {
	m.aggregateTurnovers.mu.RLock()
	defer m.aggregateTurnovers.mu.RUnlock()
	var out []tv.AggregateTurnover
	for _, at := range m.aggregateTurnovers.aggregates {
		if industrialCapitalID != "" && at.IndustrialCapitalID != industrialCapitalID {
			continue
		}
		contribs := make([]tv.ComponentTurnoverContribution, len(m.aggregateTurnovers.contributions[at.ID]))
		copy(contribs, m.aggregateTurnovers.contributions[at.ID])
		at.Contributions = contribs
		out = append(out, at)
	}
	if out == nil {
		out = []tv.AggregateTurnover{}
	}
	return out, nil
}

func (m *Memory) RecordContribution(_ context.Context, id tv.AggregateTurnoverID, c tv.ComponentTurnoverContribution) (tv.ComponentTurnoverContribution, error) {
	m.aggregateTurnovers.mu.Lock()
	defer m.aggregateTurnovers.mu.Unlock()
	if _, ok := m.aggregateTurnovers.aggregates[id]; !ok {
		return tv.ComponentTurnoverContribution{}, ErrNotFound
	}
	if c.ID == "" {
		c.ID = tv.NewComponentTurnoverContributionID()
	}
	c.AggregateTurnoverID = id
	if c.AnnualReturnPence == 0 {
		c.AnnualReturnPence = c.ComputeAnnualReturn()
	}
	if c.DifferenceKind == "" {
		c.DifferenceKind = tv.DifferenceReal
	}
	m.aggregateTurnovers.contributions[id] = append(m.aggregateTurnovers.contributions[id], c)
	return c, nil
}

func (m *Memory) RecomputeAggregate(_ context.Context, id tv.AggregateTurnoverID) (tv.AggregateTurnover, error) {
	m.aggregateTurnovers.mu.Lock()
	defer m.aggregateTurnovers.mu.Unlock()
	at, ok := m.aggregateTurnovers.aggregates[id]
	if !ok {
		return tv.AggregateTurnover{}, ErrNotFound
	}
	contribs := m.aggregateTurnovers.contributions[id]
	at = tv.RecomputeAggregate(at, contribs)
	m.aggregateTurnovers.aggregates[id] = at
	result := at
	result.Contributions = make([]tv.ComponentTurnoverContribution, len(contribs))
	copy(result.Contributions, contribs)
	return result, nil
}

func (m *Memory) RecordLifetimeCycle(_ context.Context, id tv.AggregateTurnoverID, lc tv.LifetimeCycle) (tv.LifetimeCycle, error) {
	m.aggregateTurnovers.mu.Lock()
	defer m.aggregateTurnovers.mu.Unlock()
	if _, ok := m.aggregateTurnovers.aggregates[id]; !ok {
		return tv.LifetimeCycle{}, ErrNotFound
	}
	if lc.ID == "" {
		lc.ID = tv.NewLifetimeCycleID()
	}
	lc.AggregateTurnoverID = id
	if lc.StartedAt.IsZero() {
		lc.StartedAt = time.Now().UTC()
	}
	m.aggregateTurnovers.lifetimeCycles[id] = lc
	if _, ok := m.aggregateTurnovers.crisisPhases[lc.ID]; !ok {
		m.aggregateTurnovers.crisisPhases[lc.ID] = []tv.CrisisCyclePhaseRecord{}
	}
	return lc, nil
}

func (m *Memory) RecordCrisisPhase(_ context.Context, cycleID tv.LifetimeCycleID, phase tv.CrisisCyclePhaseRecord) (tv.CrisisCyclePhaseRecord, error) {
	m.aggregateTurnovers.mu.Lock()
	defer m.aggregateTurnovers.mu.Unlock()
	if _, ok := m.aggregateTurnovers.crisisPhases[cycleID]; !ok {
		return tv.CrisisCyclePhaseRecord{}, ErrNotFound
	}
	if phase.ID == "" {
		phase.ID = tv.NewCrisisCyclePhaseRecordID()
	}
	phase.LifetimeCycleID = cycleID
	if phase.ObservedAt.IsZero() {
		phase.ObservedAt = time.Now().UTC()
	}
	m.aggregateTurnovers.crisisPhases[cycleID] = append(m.aggregateTurnovers.crisisPhases[cycleID], phase)
	return phase, nil
}
