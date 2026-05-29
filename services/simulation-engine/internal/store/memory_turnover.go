package store

import (
	"context"
	"errors"

	tv "github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

func (m *Memory) CreateTurnover(_ context.Context, t tv.Turnover) (tv.Turnover, error) {
	if err := tv.ValidateLens(t.Lens); err != nil {
		return tv.Turnover{}, err
	}
	if t.ID == "" {
		t.ID = tv.NewTurnoverID()
	}
	t.Cycles = []tv.TurnoverCycleID{}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.turnovers[t.ID]; ok {
		return tv.Turnover{}, ErrAlreadyExists
	}
	m.turnovers[t.ID] = t
	return t, nil
}

func (m *Memory) GetTurnover(_ context.Context, id tv.TurnoverID) (tv.Turnover, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.turnovers[id]
	if !ok {
		return tv.Turnover{}, ErrNotFound
	}
	cycles := m.turnoverCycles[id]
	ids := make([]tv.TurnoverCycleID, len(cycles))
	for i, c := range cycles {
		ids[i] = c.ID
	}
	t.Cycles = ids
	return t, nil
}

func (m *Memory) ListTurnovers(_ context.Context, industrialCapitalID string, lens tv.CircuitLens) ([]tv.Turnover, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []tv.Turnover
	for _, t := range m.turnovers {
		if industrialCapitalID != "" && t.IndustrialCapitalID != industrialCapitalID {
			continue
		}
		if lens != "" && t.Lens != lens {
			continue
		}
		out = append(out, t)
	}
	if out == nil {
		out = []tv.Turnover{}
	}
	return out, nil
}

func (m *Memory) RecordCycle(_ context.Context, id tv.TurnoverID, c tv.TurnoverCycle) (tv.TurnoverCycle, error) {
	if err := c.Validate(); err != nil {
		return tv.TurnoverCycle{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.turnovers[id]; !ok {
		return tv.TurnoverCycle{}, ErrNotFound
	}
	if c.ID == "" {
		c.ID = tv.NewTurnoverCycleID()
	}
	c.TurnoverID = id
	m.turnoverCycles[id] = append(m.turnoverCycles[id], c)
	return c, nil
}

func (m *Memory) RecomputeNumber(_ context.Context, id tv.TurnoverID) (tv.TurnoverNumber, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.turnovers[id]; !ok {
		return tv.TurnoverNumber{}, ErrNotFound
	}
	cycles := m.turnoverCycles[id]
	if len(cycles) == 0 {
		return tv.TurnoverNumber{}, tv.ErrNoCompletedCycles
	}
	avg := tv.AverageTurnoverTime(cycles)
	tn := tv.ComputeTurnoverNumber(id, avg)
	m.turnoverNumbers[id] = tn
	t := m.turnovers[id]
	t.TurnoverTimeMinutes = avg
	m.turnovers[id] = t
	return tn, nil
}

func (m *Memory) GetNumber(_ context.Context, id tv.TurnoverID) (tv.TurnoverNumber, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.turnovers[id]; !ok {
		return tv.TurnoverNumber{}, ErrNotFound
	}
	tn, ok := m.turnoverNumbers[id]
	if !ok {
		return tv.TurnoverNumber{}, errors.New("no turnover number computed yet")
	}
	return tn, nil
}
