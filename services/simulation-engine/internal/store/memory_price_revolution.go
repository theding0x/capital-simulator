package store

import (
	"context"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// memoryPriceRevolution is the in-memory backing store for Vol. II Ch. 15 types.
type memoryPriceRevolution struct {
	mu            sync.RWMutex
	events        map[circulation.ValueRevolutionEventID]circulation.ValueRevolutionEvent
	revaluedCap   map[circulation.ValueRevolutionEventID]circulation.RevaluedCapital
	realisedDelta map[circulation.ValueRevolutionEventID]circulation.RealisedValueDelta
	cases         map[circulation.ValueRevolutionEventID]circulation.PriceCaseRecord
	revaluations  map[circulation.IndustrialCapitalID][]circulation.InventoryRevaluation
	speculative   map[circulation.SpeculativeHoldID]circulation.SpeculativeHold
}

func newMemoryPriceRevolution() *memoryPriceRevolution {
	return &memoryPriceRevolution{
		events:        make(map[circulation.ValueRevolutionEventID]circulation.ValueRevolutionEvent),
		revaluedCap:   make(map[circulation.ValueRevolutionEventID]circulation.RevaluedCapital),
		realisedDelta: make(map[circulation.ValueRevolutionEventID]circulation.RealisedValueDelta),
		cases:         make(map[circulation.ValueRevolutionEventID]circulation.PriceCaseRecord),
		revaluations:  make(map[circulation.IndustrialCapitalID][]circulation.InventoryRevaluation),
		speculative:   make(map[circulation.SpeculativeHoldID]circulation.SpeculativeHold),
	}
}

// CreatePriceRevolution persists a ValueRevolutionEvent and its derived records.
func (m *Memory) CreatePriceRevolution(_ context.Context, rec circulation.PriceRevolutionRecord) (circulation.PriceRevolutionRecord, error) {
	m.priceRevolution.mu.Lock()
	defer m.priceRevolution.mu.Unlock()

	evt := rec.Event
	if evt.ID == "" {
		evt.ID = circulation.NewValueRevolutionEventID()
	}
	if evt.OccurredAt.IsZero() {
		evt.OccurredAt = time.Now().UTC()
	}
	if _, exists := m.priceRevolution.events[evt.ID]; exists {
		return circulation.PriceRevolutionRecord{}, ErrAlreadyExists
	}
	m.priceRevolution.events[evt.ID] = evt
	rec.Event = evt

	if rec.RevaluedCapital != nil {
		rc := *rec.RevaluedCapital
		if rc.ID == "" {
			rc.ID = circulation.NewRevaluedCapitalID()
		}
		m.priceRevolution.revaluedCap[evt.ID] = rc
		rec.RevaluedCapital = &rc
	}

	if rec.RealisedDelta != nil {
		rvd := *rec.RealisedDelta
		if rvd.ID == "" {
			rvd.ID = circulation.NewRealisedValueDeltaID()
		}
		m.priceRevolution.realisedDelta[evt.ID] = rvd
		rec.RealisedDelta = &rvd
	}

	if rec.Revaluations == nil {
		rec.Revaluations = []circulation.InventoryRevaluation{}
	}

	return rec, nil
}

// GetPriceRevolution assembles a PriceRevolutionRecord by event ID.
func (m *Memory) GetPriceRevolution(_ context.Context, id circulation.ValueRevolutionEventID) (circulation.PriceRevolutionRecord, error) {
	m.priceRevolution.mu.RLock()
	defer m.priceRevolution.mu.RUnlock()

	evt, ok := m.priceRevolution.events[id]
	if !ok {
		return circulation.PriceRevolutionRecord{}, ErrNotFound
	}
	rec := circulation.PriceRevolutionRecord{
		Event:        evt,
		Revaluations: m.revaluationsForCapital(evt.IndustrialCapitalID),
	}
	if rc, ok := m.priceRevolution.revaluedCap[id]; ok {
		rc2 := rc
		rec.RevaluedCapital = &rc2
	}
	if rvd, ok := m.priceRevolution.realisedDelta[id]; ok {
		rvd2 := rvd
		rec.RealisedDelta = &rvd2
	}
	if pc, ok := m.priceRevolution.cases[id]; ok {
		if pc.FallCase != nil {
			fc := *pc.FallCase
			rec.FallCase = &fc
		}
		if pc.RiseCase != nil {
			rcase := *pc.RiseCase
			rec.RiseCase = &rcase
		}
	}
	return rec, nil
}

// revaluationsForCapital returns inventory revaluations for a capital ID.
// Caller must hold the read lock.
func (m *Memory) revaluationsForCapital(id circulation.IndustrialCapitalID) []circulation.InventoryRevaluation {
	revs := m.priceRevolution.revaluations[id]
	if revs == nil {
		return []circulation.InventoryRevaluation{}
	}
	out := make([]circulation.InventoryRevaluation, len(revs))
	copy(out, revs)
	return out
}

// RecordPriceCase stores the operator-selected case for an event.
func (m *Memory) RecordPriceCase(_ context.Context, pc circulation.PriceCaseRecord) (circulation.PriceCaseRecord, error) {
	if err := pc.Validate(); err != nil {
		return circulation.PriceCaseRecord{}, err
	}
	m.priceRevolution.mu.Lock()
	defer m.priceRevolution.mu.Unlock()

	if _, ok := m.priceRevolution.events[pc.ValueRevolutionEventID]; !ok {
		return circulation.PriceCaseRecord{}, ErrNotFound
	}
	m.priceRevolution.cases[pc.ValueRevolutionEventID] = pc
	return pc, nil
}

// RecordInventoryRevaluation appends an InventoryRevaluation for a capital.
func (m *Memory) RecordInventoryRevaluation(_ context.Context, inv circulation.InventoryRevaluation) (circulation.InventoryRevaluation, error) {
	if err := inv.Validate(); err != nil {
		return circulation.InventoryRevaluation{}, err
	}
	m.priceRevolution.mu.Lock()
	defer m.priceRevolution.mu.Unlock()

	if inv.ID == "" {
		inv.ID = circulation.NewInventoryRevaluationID()
	}
	m.priceRevolution.revaluations[inv.IndustrialCapitalID] = append(
		m.priceRevolution.revaluations[inv.IndustrialCapitalID], inv)
	return inv, nil
}

// RecordSpeculativeHold stores a speculative hold by its ID.
func (m *Memory) RecordSpeculativeHold(_ context.Context, sh circulation.SpeculativeHold) (circulation.SpeculativeHold, error) {
	if err := sh.Validate(); err != nil {
		return circulation.SpeculativeHold{}, err
	}
	m.priceRevolution.mu.Lock()
	defer m.priceRevolution.mu.Unlock()

	if sh.ID == "" {
		sh.ID = circulation.NewSpeculativeHoldID()
	}
	if sh.HeldSince.IsZero() {
		sh.HeldSince = time.Now().UTC()
	}
	m.priceRevolution.speculative[sh.ID] = sh
	return sh, nil
}

// AggregateCompound sums all delta streams for a capital+period and returns
// a CompoundCapitalAdjustment. If no deltas exist, returns a zero adjustment.
func (m *Memory) AggregateCompound(_ context.Context, icID circulation.IndustrialCapitalID, period string) (circulation.CompoundCapitalAdjustment, error) {
	m.priceRevolution.mu.RLock()
	defer m.priceRevolution.mu.RUnlock()

	var revDelta, reaDelta, invDelta circulation.Pence

	for _, evt := range m.priceRevolution.events {
		if evt.IndustrialCapitalID != icID {
			continue
		}
		if rc, ok := m.priceRevolution.revaluedCap[evt.ID]; ok {
			revDelta += rc.DeltaPence
		}
		if rvd, ok := m.priceRevolution.realisedDelta[evt.ID]; ok {
			reaDelta += rvd.DeltaPence
		}
	}

	for _, inv := range m.priceRevolution.revaluations[icID] {
		invDelta += inv.DeltaPence
	}

	net := revDelta + reaDelta + invDelta
	return circulation.CompoundCapitalAdjustment{
		ID:                    circulation.NewCompoundCapitalAdjustmentID(),
		IndustrialCapitalID:   icID,
		Period:                period,
		RevaluationDeltaPence: revDelta,
		RealisationDeltaPence: reaDelta,
		InventoryDeltaPence:   invDelta,
		NetEffectPence:        net,
	}, nil
}
