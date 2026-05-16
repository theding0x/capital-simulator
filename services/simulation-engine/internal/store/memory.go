package store

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// Memory is the in-memory implementation of MachineStore, FactoryStore,
// GeneralLawStore, HistoricalStageStore, and EnclosureEventStore.
type Memory struct {
	mu               sync.RWMutex
	machines         map[machinery.MachineID]machinery.Machine
	factories        map[machinery.FactoryID]machinery.Factory
	ticks            map[machinery.FactoryID][]engine.Tick
	generalLaw       map[simulation.GeneralLawScenarioID]simulation.GeneralLawScenario
	historicalStages map[simulation.HistoricalStageID]simulation.HistoricalStage
	stageNames       map[string]simulation.HistoricalStageID
	enclosureEvents  []simulation.EnclosureEvent
	now              func() time.Time
}

func NewMemory() *Memory {
	return &Memory{
		machines:         make(map[machinery.MachineID]machinery.Machine),
		factories:        make(map[machinery.FactoryID]machinery.Factory),
		ticks:            make(map[machinery.FactoryID][]engine.Tick),
		generalLaw:       make(map[simulation.GeneralLawScenarioID]simulation.GeneralLawScenario),
		historicalStages: make(map[simulation.HistoricalStageID]simulation.HistoricalStage),
		stageNames:       make(map[string]simulation.HistoricalStageID),
		now:              time.Now,
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

func (m *Memory) CreateGeneralLawScenario(_ context.Context, s simulation.GeneralLawScenario) (simulation.GeneralLawScenario, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s.ID.IsZero() {
		s.ID = simulation.NewGeneralLawScenarioID()
	}
	if _, ok := m.generalLaw[s.ID]; ok {
		return simulation.GeneralLawScenario{}, ErrAlreadyExists
	}
	s.CreatedAt = m.now().UTC()
	m.generalLaw[s.ID] = s
	return s, nil
}

func (m *Memory) GetGeneralLawScenario(_ context.Context, id simulation.GeneralLawScenarioID) (simulation.GeneralLawScenario, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.generalLaw[id]
	if !ok {
		return simulation.GeneralLawScenario{}, ErrNotFound
	}
	return s, nil
}

func (m *Memory) CreateHistoricalStage(_ context.Context, h simulation.HistoricalStage) (simulation.HistoricalStage, error) {
	if err := h.Validate(); err != nil {
		return simulation.HistoricalStage{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.stageNames[strings.ToLower(h.Name)]; ok {
		return simulation.HistoricalStage{}, ErrAlreadyExists
	}
	if h.ID.IsZero() {
		h.ID = simulation.NewHistoricalStageID()
	}
	if _, ok := m.historicalStages[h.ID]; ok {
		return simulation.HistoricalStage{}, ErrAlreadyExists
	}
	h.CreatedAt = m.now().UTC()
	stored := simulation.HistoricalStage{
		ID:                     h.ID,
		Name:                   h.Name,
		Description:            h.Description,
		PrimitiveAccumulations: append([]simulation.PrimitiveAccumulation(nil), h.PrimitiveAccumulations...),
		CreatedAt:              h.CreatedAt,
	}
	m.historicalStages[h.ID] = stored
	m.stageNames[strings.ToLower(h.Name)] = h.ID
	return stored, nil
}

func (m *Memory) GetHistoricalStage(_ context.Context, id simulation.HistoricalStageID) (simulation.HistoricalStage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h, ok := m.historicalStages[id]
	if !ok {
		return simulation.HistoricalStage{}, ErrNotFound
	}
	cp := h
	cp.PrimitiveAccumulations = append([]simulation.PrimitiveAccumulation(nil), h.PrimitiveAccumulations...)
	return cp, nil
}

func (m *Memory) ListHistoricalStages(_ context.Context) ([]simulation.HistoricalStage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.HistoricalStage, 0, len(m.historicalStages))
	for _, h := range m.historicalStages {
		cp := h
		cp.PrimitiveAccumulations = append([]simulation.PrimitiveAccumulation(nil), h.PrimitiveAccumulations...)
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
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

func (m *Memory) CreateEnclosureEvent(_ context.Context, e simulation.EnclosureEvent) (simulation.EnclosureEvent, error) {
	if err := e.Validate(); err != nil {
		return simulation.EnclosureEvent{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if e.ID.IsZero() {
		e.ID = simulation.NewEnclosureEventID()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = m.now()
	}
	m.enclosureEvents = append(m.enclosureEvents, e)
	return e, nil
}

func (m *Memory) ListEnclosureEvents(_ context.Context) ([]simulation.EnclosureEvent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.EnclosureEvent, len(m.enclosureEvents))
	copy(out, m.enclosureEvents)
	return out, nil
}
