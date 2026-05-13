package store

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
)

// Memory is the in-memory implementation of MachineStore and FactoryStore.
type Memory struct {
	mu        sync.RWMutex
	machines  map[machinery.MachineID]machinery.Machine
	factories map[machinery.FactoryID]machinery.Factory
	ticks     map[machinery.FactoryID][]engine.Tick
	now       func() time.Time
}

func NewMemory() *Memory {
	return &Memory{
		machines:  make(map[machinery.MachineID]machinery.Machine),
		factories: make(map[machinery.FactoryID]machinery.Factory),
		ticks:     make(map[machinery.FactoryID][]engine.Tick),
		now:       time.Now,
	}
}

func (m *Memory) CreateMachine(_ context.Context, mc machinery.Machine) (machinery.Machine, error) {
	if err := mc.Validate(); err != nil {
		return machinery.Machine{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if mc.ID.IsZero() {
		mc.ID = machinery.NewMachineID()
	}
	if _, ok := m.machines[mc.ID]; ok {
		return machinery.Machine{}, ErrAlreadyExists
	}
	mc.CreatedAt = m.now().UTC()
	m.machines[mc.ID] = mc
	return mc, nil
}

func (m *Memory) GetMachine(_ context.Context, id machinery.MachineID) (machinery.Machine, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mc, ok := m.machines[id]
	if !ok {
		return machinery.Machine{}, ErrNotFound
	}
	return mc, nil
}

func (m *Memory) ListMachines(_ context.Context) ([]machinery.Machine, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]machinery.Machine, 0, len(m.machines))
	for _, mc := range m.machines {
		out = append(out, mc)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (m *Memory) UpdateMachine(_ context.Context, id machinery.MachineID, u MachineUpdate) (machinery.Machine, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.machines[id]
	if !ok {
		return machinery.Machine{}, ErrNotFound
	}
	if u.AccumulatedWear != nil {
		cur.AccumulatedWear = *u.AccumulatedWear
	}
	if u.AccumulatedDepreciation != nil {
		cur.AccumulatedDepreciation = *u.AccumulatedDepreciation
	}
	m.machines[id] = cur
	return cur, nil
}

func (m *Memory) CreateFactory(ctx context.Context, f machinery.Factory) (machinery.Factory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if f.ID.IsZero() {
		f.ID = machinery.NewFactoryID()
	}
	if _, ok := m.factories[f.ID]; ok {
		return machinery.Factory{}, ErrAlreadyExists
	}
	// Resolve machine references against the store. Factory may be created
	// either with embedded Machine records or with bare IDs; we always
	// persist the embedded form.
	resolved := make([]machinery.Machine, 0, len(f.Machines))
	for _, mc := range f.Machines {
		if !mc.ID.IsZero() {
			if existing, ok := m.machines[mc.ID]; ok {
				resolved = append(resolved, existing)
				continue
			}
		}
		// Inline machine: persist it transparently.
		if err := mc.Validate(); err != nil {
			return machinery.Factory{}, err
		}
		if mc.ID.IsZero() {
			mc.ID = machinery.NewMachineID()
		}
		mc.CreatedAt = m.now().UTC()
		m.machines[mc.ID] = mc
		resolved = append(resolved, mc)
	}
	f.Machines = resolved
	if err := f.Validate(); err != nil {
		return machinery.Factory{}, err
	}
	f.CreatedAt = m.now().UTC()
	m.factories[f.ID] = f
	return f, nil
}

func (m *Memory) GetFactory(_ context.Context, id machinery.FactoryID) (machinery.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.factories[id]
	if !ok {
		return machinery.Factory{}, ErrNotFound
	}
	// Refresh embedded machines from the canonical store map so accumulated
	// wear reads come back current.
	refreshed := make([]machinery.Machine, 0, len(f.Machines))
	for _, mc := range f.Machines {
		if cur, ok := m.machines[mc.ID]; ok {
			refreshed = append(refreshed, cur)
		} else {
			refreshed = append(refreshed, mc)
		}
	}
	f.Machines = refreshed
	return f, nil
}

func (m *Memory) ListFactories(_ context.Context) ([]machinery.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]machinery.Factory, 0, len(m.factories))
	for _, f := range m.factories {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (m *Memory) AdvanceTick(_ context.Context, id machinery.FactoryID) (machinery.Factory, engine.Tick, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f, ok := m.factories[id]
	if !ok {
		return machinery.Factory{}, engine.Tick{}, ErrNotFound
	}
	// Refresh embedded machines.
	live := make([]machinery.Machine, 0, len(f.Machines))
	for _, mc := range f.Machines {
		if cur, ok := m.machines[mc.ID]; ok {
			live = append(live, cur)
		} else {
			live = append(live, mc)
		}
	}
	f.Machines = live
	result := f.RunTick()
	// Accumulate wear per machine.
	for i, mc := range f.Machines {
		dwt := machinery.DailyWearAndTear(mc)
		mc.AccumulatedWear.Value += dwt
		f.Machines[i] = mc
		m.machines[mc.ID] = mc
	}
	f.TickCount++
	m.factories[id] = f
	tick := engine.Tick{
		FactoryID:        string(id),
		Sequence:         f.TickCount,
		ValueTransferred: int64(result.ValueTransferred),
		UnitsProduced:    result.UnitsProduced,
		HandLabourSaved:  int64(result.HandLabourSaved),
		OccurredAt:       m.now().UTC(),
	}
	m.ticks[id] = append(m.ticks[id], tick)
	return f, tick, nil
}

func (m *Memory) ListTicks(_ context.Context, id machinery.FactoryID, limit int) ([]engine.Tick, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	all := m.ticks[id]
	if limit <= 0 || limit > len(all) {
		limit = len(all)
	}
	// Return the latest `limit` ticks in ascending-sequence order.
	start := len(all) - limit
	out := make([]engine.Tick, limit)
	copy(out, all[start:])
	return out, nil
}
