package store

import (
	"context"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
)

// RecordEngineTick appends one scheduler-pass audit row. A missing ID is
// assigned and a zero OccurredAt is stamped with the store clock. The log is
// append-only, so this never returns ErrAlreadyExists.
func (m *Memory) RecordEngineTick(_ context.Context, t engine.ScheduledTick) (engine.ScheduledTick, error) {
	if t.ID == "" {
		t.ID = engine.NewScheduledTickID()
	}
	if t.OccurredAt.IsZero() {
		t.OccurredAt = m.now().UTC()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.engineTicks = append(m.engineTicks, t)
	return t, nil
}

// ListEngineTicks returns up to limit most-recent scheduler passes, newest
// first. A non-positive limit returns every row. Never returns nil.
func (m *Memory) ListEngineTicks(_ context.Context, limit int) ([]engine.ScheduledTick, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]engine.ScheduledTick, 0, len(m.engineTicks))
	for i := len(m.engineTicks) - 1; i >= 0; i-- {
		out = append(out, m.engineTicks[i])
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}
