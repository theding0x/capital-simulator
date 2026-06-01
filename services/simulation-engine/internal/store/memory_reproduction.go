package store

import (
	"context"
	"sort"
	"sync"
	"time"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// memoryReproduction is the in-memory backing store for Vol. II Ch. 17 types.
type memoryReproduction struct {
	mu                 sync.RWMutex
	circulations       map[repro.SurplusCirculationID]repro.SurplusCirculation
	createdAt          map[repro.SurplusCirculationID]time.Time
	realisationSources map[repro.SurplusCirculationID][]repro.RealisationSourceEntry
	aggregates         []repro.SocialCapitalAggregate
	realisationPuzzles []repro.RealisationPuzzle
}

func newMemoryReproduction() *memoryReproduction {
	return &memoryReproduction{
		circulations:       make(map[repro.SurplusCirculationID]repro.SurplusCirculation),
		createdAt:          make(map[repro.SurplusCirculationID]time.Time),
		realisationSources: make(map[repro.SurplusCirculationID][]repro.RealisationSourceEntry),
	}
}

// CreateSurplusCirculation persists a new SurplusCirculation. The
// RealisationSources slice is ignored on creation; sources are added via
// RecordRealisationSource.
func (m *Memory) CreateSurplusCirculation(_ context.Context, sc repro.SurplusCirculation) (repro.SurplusCirculation, error) {
	if sc.TotalSurplusPence < 0 {
		return repro.SurplusCirculation{}, repro.ErrNegativeSurplus
	}
	if sc.ID.IsZero() {
		sc.ID = repro.NewSurplusCirculationID()
	}
	sc.RealisedSurplusPence = 0
	sc.UnrealisedSurplusPence = sc.TotalSurplusPence
	sc.RealisationSources = []repro.RealisationSourceEntry{}

	m.reproduction.mu.Lock()
	defer m.reproduction.mu.Unlock()
	if _, exists := m.reproduction.circulations[sc.ID]; exists {
		return repro.SurplusCirculation{}, ErrAlreadyExists
	}
	m.reproduction.circulations[sc.ID] = sc
	m.reproduction.createdAt[sc.ID] = m.now().UTC()
	m.reproduction.realisationSources[sc.ID] = []repro.RealisationSourceEntry{}
	return sc, nil
}

// GetSurplusCirculation returns a SurplusCirculation with all source entries.
func (m *Memory) GetSurplusCirculation(_ context.Context, id repro.SurplusCirculationID) (repro.SurplusCirculation, error) {
	m.reproduction.mu.RLock()
	defer m.reproduction.mu.RUnlock()
	sc, ok := m.reproduction.circulations[id]
	if !ok {
		return repro.SurplusCirculation{}, ErrNotFound
	}
	sources := m.reproduction.realisationSources[id]
	sc.RealisationSources = append([]repro.RealisationSourceEntry(nil), sources...)
	return sc, nil
}

// ListSurplusCirculations returns all SurplusCirculation records ordered by
// creation time ascending.
func (m *Memory) ListSurplusCirculations(_ context.Context) ([]repro.SurplusCirculation, error) {
	m.reproduction.mu.RLock()
	defer m.reproduction.mu.RUnlock()
	type item struct {
		sc repro.SurplusCirculation
		t  time.Time
	}
	items := make([]item, 0, len(m.reproduction.circulations))
	for _, sc := range m.reproduction.circulations {
		sources := m.reproduction.realisationSources[sc.ID]
		sc.RealisationSources = append([]repro.RealisationSourceEntry(nil), sources...)
		items = append(items, item{sc: sc, t: m.reproduction.createdAt[sc.ID]})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].t.Before(items[j].t)
	})
	out := make([]repro.SurplusCirculation, len(items))
	for i, it := range items {
		out[i] = it.sc
	}
	return out, nil
}

// RecordRealisationSource appends a source entry to a SurplusCirculation,
// updating realised/unrealised pence atomically.
func (m *Memory) RecordRealisationSource(_ context.Context, id repro.SurplusCirculationID, entry repro.RealisationSourceEntry) (repro.RealisationSourceEntry, error) {
	if entry.Pence < 0 {
		return repro.RealisationSourceEntry{}, repro.ErrNegativePence
	}
	if err := entry.Source.Validate(); err != nil {
		return repro.RealisationSourceEntry{}, err
	}

	m.reproduction.mu.Lock()
	defer m.reproduction.mu.Unlock()

	sc, ok := m.reproduction.circulations[id]
	if !ok {
		return repro.RealisationSourceEntry{}, ErrNotFound
	}

	entry.SurplusCirculationID = id
	if entry.Source == repro.SourceCredit {
		entry.DeferredRealisationPence = entry.Pence
	}

	newRealised := sc.RealisedSurplusPence + entry.Pence
	if newRealised > sc.TotalSurplusPence {
		return repro.RealisationSourceEntry{}, repro.ErrRealisationExceedsTotal
	}
	sc.RealisedSurplusPence = newRealised
	sc.UnrealisedSurplusPence = sc.TotalSurplusPence - newRealised
	m.reproduction.circulations[id] = sc
	m.reproduction.realisationSources[id] = append(m.reproduction.realisationSources[id], entry)
	return entry, nil
}

// SnapshotAggregate returns the SocialCapitalAggregate for the given period.
// If an explicit aggregate row exists (inserted by the seed migration), it is
// returned; otherwise TotalAnnualSurplusPence is derived by summing all
// SurplusCirculation rows whose period matches.
func (m *Memory) SnapshotAggregate(_ context.Context, period string) (repro.SocialCapitalAggregate, error) {
	m.reproduction.mu.RLock()
	defer m.reproduction.mu.RUnlock()

	for _, agg := range m.reproduction.aggregates {
		if agg.Period == period {
			return agg, nil
		}
	}

	agg := repro.SocialCapitalAggregate{Period: period}
	for _, sc := range m.reproduction.circulations {
		if sc.Period == period {
			agg.TotalAnnualSurplusPence += sc.TotalSurplusPence
		}
	}
	return agg, nil
}

// UpsertAggregateShares sets the Department I/II share split on the aggregate
// for the given period, updating the existing row in place or appending a new
// one. Other aggregate fields are left untouched.
func (m *Memory) UpsertAggregateShares(_ context.Context, period string, deptIShareBP, deptIIShareBP int64) (repro.SocialCapitalAggregate, error) {
	m.reproduction.mu.Lock()
	defer m.reproduction.mu.Unlock()

	for i := range m.reproduction.aggregates {
		if m.reproduction.aggregates[i].Period == period {
			m.reproduction.aggregates[i].DepartmentIShareBP = deptIShareBP
			m.reproduction.aggregates[i].DepartmentIIShareBP = deptIIShareBP
			return m.reproduction.aggregates[i], nil
		}
	}

	agg := repro.SocialCapitalAggregate{
		Period:              period,
		DepartmentIShareBP:  deptIShareBP,
		DepartmentIIShareBP: deptIIShareBP,
	}
	m.reproduction.aggregates = append(m.reproduction.aggregates, agg)
	return agg, nil
}

// ListRealisationPuzzles returns all canonical RealisationPuzzle rows.
func (m *Memory) ListRealisationPuzzles(_ context.Context) ([]repro.RealisationPuzzle, error) {
	m.reproduction.mu.RLock()
	defer m.reproduction.mu.RUnlock()
	out := make([]repro.RealisationPuzzle, len(m.reproduction.realisationPuzzles))
	copy(out, m.reproduction.realisationPuzzles)
	return out, nil
}
