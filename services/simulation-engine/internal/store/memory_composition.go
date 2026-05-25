package store

import (
	"context"
	"sync"
	"time"

	comp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
)

type memoryComposition struct {
	mu                sync.RWMutex
	components        map[comp.CapitalComponentID]comp.CapitalComponent
	fixedItems        map[comp.FixedCapitalItemID]comp.FixedCapitalItem
	subcomponents     map[comp.FixedCapitalItemID][]comp.FixedCapitalSubcomponent
	wearAndTear       map[comp.FixedCapitalItemID][]comp.WearAndTear
	sinkingFunds      map[comp.FixedCapitalItemID]comp.SinkingFundForItem
	repairs           map[comp.FixedCapitalItemID][]comp.Repair
	circulatingCycles map[comp.CapitalComponentID][]comp.CirculatingCycle
	reinvestments     map[comp.FixedCapitalItemID][]comp.SinkingFundReinvestment
}

func newMemoryComposition() *memoryComposition {
	return &memoryComposition{
		components:        make(map[comp.CapitalComponentID]comp.CapitalComponent),
		fixedItems:        make(map[comp.FixedCapitalItemID]comp.FixedCapitalItem),
		subcomponents:     make(map[comp.FixedCapitalItemID][]comp.FixedCapitalSubcomponent),
		wearAndTear:       make(map[comp.FixedCapitalItemID][]comp.WearAndTear),
		sinkingFunds:      make(map[comp.FixedCapitalItemID]comp.SinkingFundForItem),
		repairs:           make(map[comp.FixedCapitalItemID][]comp.Repair),
		circulatingCycles: make(map[comp.CapitalComponentID][]comp.CirculatingCycle),
		reinvestments:     make(map[comp.FixedCapitalItemID][]comp.SinkingFundReinvestment),
	}
}

func (m *Memory) CreateComponent(_ context.Context, c comp.CapitalComponent) (comp.CapitalComponent, error) {
	if err := c.Validate(); err != nil {
		return comp.CapitalComponent{}, err
	}
	if c.ID == "" {
		c.ID = comp.NewCapitalComponentID()
	}
	if c.EnteredAt.IsZero() {
		c.EnteredAt = time.Now().UTC()
	}
	m.composition.mu.Lock()
	defer m.composition.mu.Unlock()
	if _, ok := m.composition.components[c.ID]; ok {
		return comp.CapitalComponent{}, ErrAlreadyExists
	}
	m.composition.components[c.ID] = c
	return c, nil
}

func (m *Memory) GetComponent(_ context.Context, id comp.CapitalComponentID) (comp.CapitalComponent, error) {
	m.composition.mu.RLock()
	defer m.composition.mu.RUnlock()
	c, ok := m.composition.components[id]
	if !ok {
		return comp.CapitalComponent{}, ErrNotFound
	}
	return c, nil
}

func (m *Memory) ListComponents(_ context.Context, industrialCapitalID string, kind comp.CapitalKind, role comp.RoleInProcess) ([]comp.CapitalComponent, error) {
	m.composition.mu.RLock()
	defer m.composition.mu.RUnlock()
	var out []comp.CapitalComponent
	for _, c := range m.composition.components {
		if industrialCapitalID != "" && c.IndustrialCapitalID != industrialCapitalID {
			continue
		}
		if kind != "" && c.Kind != kind {
			continue
		}
		if role != "" && c.Role != role {
			continue
		}
		out = append(out, c)
	}
	if out == nil {
		out = []comp.CapitalComponent{}
	}
	return out, nil
}

func (m *Memory) RegisterFixedItem(_ context.Context, item comp.FixedCapitalItem) (comp.FixedCapitalItem, error) {
	if err := item.Validate(); err != nil {
		return comp.FixedCapitalItem{}, err
	}
	if item.ID == "" {
		item.ID = comp.NewFixedCapitalItemID()
	}
	if item.PurchasedAt.IsZero() {
		item.PurchasedAt = time.Now().UTC()
	}
	if item.WearModel == "" {
		item.WearModel = comp.WearLinear
	}
	m.composition.mu.Lock()
	defer m.composition.mu.Unlock()
	if _, ok := m.composition.fixedItems[item.ID]; ok {
		return comp.FixedCapitalItem{}, ErrAlreadyExists
	}
	m.composition.fixedItems[item.ID] = item
	m.composition.sinkingFunds[item.ID] = comp.SinkingFundForItem{
		ID:                 comp.NewSinkingFundForItemID(),
		FixedCapitalItemID: item.ID,
		AccumulatedPence:   0,
		TargetPence:        item.PencePurchased,
	}
	return item, nil
}

func (m *Memory) GetFixedItem(_ context.Context, id comp.FixedCapitalItemID) (comp.FixedCapitalItem, error) {
	m.composition.mu.RLock()
	defer m.composition.mu.RUnlock()
	item, ok := m.composition.fixedItems[id]
	if !ok {
		return comp.FixedCapitalItem{}, ErrNotFound
	}
	return item, nil
}

func (m *Memory) RecordSubcomponent(_ context.Context, sub comp.FixedCapitalSubcomponent) (comp.FixedCapitalSubcomponent, error) {
	if sub.ID == "" {
		sub.ID = comp.NewFixedCapitalSubcomponentID()
	}
	if sub.WearModel == "" {
		sub.WearModel = comp.WearLinear
	}
	m.composition.mu.Lock()
	defer m.composition.mu.Unlock()
	m.composition.subcomponents[sub.FixedCapitalItemID] = append(m.composition.subcomponents[sub.FixedCapitalItemID], sub)
	return sub, nil
}

func (m *Memory) RecordWear(_ context.Context, w comp.WearAndTear) (comp.WearAndTear, error) {
	if w.ID == "" {
		w.ID = comp.NewWearAndTearID()
	}
	if w.CalculatedAt.IsZero() {
		w.CalculatedAt = time.Now().UTC()
	}
	m.composition.mu.Lock()
	defer m.composition.mu.Unlock()
	m.composition.wearAndTear[w.FixedCapitalItemID] = append(m.composition.wearAndTear[w.FixedCapitalItemID], w)
	if sf, ok := m.composition.sinkingFunds[w.FixedCapitalItemID]; ok {
		sf.AccumulatedPence += w.Pence
		m.composition.sinkingFunds[w.FixedCapitalItemID] = sf
	}
	return w, nil
}

func (m *Memory) GetSinkingFund(_ context.Context, itemID comp.FixedCapitalItemID) (comp.SinkingFundForItem, error) {
	m.composition.mu.RLock()
	defer m.composition.mu.RUnlock()
	sf, ok := m.composition.sinkingFunds[itemID]
	if !ok {
		return comp.SinkingFundForItem{}, ErrNotFound
	}
	return sf, nil
}

func (m *Memory) RecordRepair(_ context.Context, r comp.Repair) (comp.Repair, error) {
	if err := r.Validate(); err != nil {
		return comp.Repair{}, err
	}
	if r.ID == "" {
		r.ID = comp.NewRepairID()
	}
	if r.OccurredAt.IsZero() {
		r.OccurredAt = time.Now().UTC()
	}
	m.composition.mu.Lock()
	defer m.composition.mu.Unlock()
	if r.Kind == comp.RepairLarge {
		if item, ok := m.composition.fixedItems[r.FixedCapitalItemID]; ok {
			item.ServiceLifeMinutes += r.ServiceLifeExtensionMinutes
			m.composition.fixedItems[r.FixedCapitalItemID] = item
		}
	}
	m.composition.repairs[r.FixedCapitalItemID] = append(m.composition.repairs[r.FixedCapitalItemID], r)
	return r, nil
}

func (m *Memory) RecordCirculatingCycle(_ context.Context, c comp.CirculatingCycle) (comp.CirculatingCycle, error) {
	if c.ID == "" {
		c.ID = comp.NewCirculatingCycleID()
	}
	m.composition.mu.Lock()
	defer m.composition.mu.Unlock()
	m.composition.circulatingCycles[c.CapitalComponentID] = append(m.composition.circulatingCycles[c.CapitalComponentID], c)
	return c, nil
}

func (m *Memory) RecordReinvestment(_ context.Context, r comp.SinkingFundReinvestment) (comp.SinkingFundReinvestment, error) {
	if r.ID == "" {
		r.ID = comp.NewSinkingFundReinvestmentID()
	}
	if r.OccurredAt.IsZero() {
		r.OccurredAt = time.Now().UTC()
	}
	m.composition.mu.Lock()
	defer m.composition.mu.Unlock()
	m.composition.reinvestments[r.FixedCapitalItemID] = append(m.composition.reinvestments[r.FixedCapitalItemID], r)
	return r, nil
}
