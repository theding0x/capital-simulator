package store

import (
	"context"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
)

// RecordEngineTick inserts one scheduler-pass audit row into engine_ticks. A
// missing ID is assigned and a zero OccurredAt is stamped with the store clock.
func (m *MySQL) RecordEngineTick(ctx context.Context, t engine.ScheduledTick) (engine.ScheduledTick, error) {
	if t.ID == "" {
		t.ID = engine.NewScheduledTickID()
	}
	if t.OccurredAt.IsZero() {
		t.OccurredAt = m.now().UTC()
	}
	const q = `INSERT INTO engine_ticks
		(id, sequence, occurred_at, duration_ms, tickers_run, entities_advanced, error_count)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	if _, err := m.db.ExecContext(ctx, q,
		t.ID, t.Sequence, t.OccurredAt, t.DurationMS, t.TickersRun, t.EntitiesAdvanced, t.ErrorCount,
	); err != nil {
		return engine.ScheduledTick{}, err
	}
	return t, nil
}

// ListEngineTicks returns up to limit most-recent scheduler passes, newest
// first (by sequence). A non-positive limit returns every row. Never nil.
func (m *MySQL) ListEngineTicks(ctx context.Context, limit int) ([]engine.ScheduledTick, error) {
	q := `SELECT id, sequence, occurred_at, duration_ms, tickers_run, entities_advanced, error_count
		FROM engine_ticks ORDER BY sequence DESC`
	var args []any
	if limit > 0 {
		q += " LIMIT ?"
		args = append(args, limit)
	}
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]engine.ScheduledTick, 0)
	for rows.Next() {
		var t engine.ScheduledTick
		if err := rows.Scan(
			&t.ID, &t.Sequence, &t.OccurredAt, &t.DurationMS,
			&t.TickersRun, &t.EntitiesAdvanced, &t.ErrorCount,
		); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}
