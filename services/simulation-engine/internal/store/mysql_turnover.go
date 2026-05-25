package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	tv "github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

func (s *MySQL) CreateTurnover(ctx context.Context, t tv.Turnover) (tv.Turnover, error) {
	if err := tv.ValidateLens(t.Lens); err != nil {
		return tv.Turnover{}, err
	}
	if t.ID == "" {
		t.ID = tv.NewTurnoverID()
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO turnovers (id, industrial_capital_id, lens, turnover_time_minutes) VALUES (?,?,?,?)`,
		string(t.ID), t.IndustrialCapitalID, string(t.Lens), t.TurnoverTimeMinutes)
	if err != nil {
		if isDuplicateKey(err) {
			return tv.Turnover{}, ErrAlreadyExists
		}
		return tv.Turnover{}, err
	}
	t.Cycles = []tv.TurnoverCycleID{}
	return t, nil
}

func (s *MySQL) GetTurnover(ctx context.Context, id tv.TurnoverID) (tv.Turnover, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, industrial_capital_id, lens, turnover_time_minutes FROM turnovers WHERE id=?`, string(id))
	var t tv.Turnover
	var lens string
	if err := row.Scan(&t.ID, &t.IndustrialCapitalID, &lens, &t.TurnoverTimeMinutes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tv.Turnover{}, ErrNotFound
		}
		return tv.Turnover{}, err
	}
	t.Lens = tv.CircuitLens(lens)
	rows, err := s.db.QueryContext(ctx,
		`SELECT id FROM turnover_cycles WHERE turnover_id=? ORDER BY started_at`, string(id))
	if err != nil {
		return tv.Turnover{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var cid tv.TurnoverCycleID
		if err := rows.Scan(&cid); err != nil {
			return tv.Turnover{}, err
		}
		t.Cycles = append(t.Cycles, cid)
	}
	if t.Cycles == nil {
		t.Cycles = []tv.TurnoverCycleID{}
	}
	return t, nil
}

func (s *MySQL) ListTurnovers(ctx context.Context, industrialCapitalID string, lens tv.CircuitLens) ([]tv.Turnover, error) {
	q := `SELECT id, industrial_capital_id, lens, turnover_time_minutes FROM turnovers WHERE 1=1`
	var args []any
	if industrialCapitalID != "" {
		q += " AND industrial_capital_id=?"
		args = append(args, industrialCapitalID)
	}
	if lens != "" {
		q += " AND lens=?"
		args = append(args, string(lens))
	}
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []tv.Turnover
	for rows.Next() {
		var t tv.Turnover
		var l string
		if err := rows.Scan(&t.ID, &t.IndustrialCapitalID, &l, &t.TurnoverTimeMinutes); err != nil {
			return nil, err
		}
		t.Lens = tv.CircuitLens(l)
		t.Cycles = []tv.TurnoverCycleID{}
		out = append(out, t)
	}
	if out == nil {
		out = []tv.Turnover{}
	}
	return out, nil
}

func (s *MySQL) RecordCycle(ctx context.Context, id tv.TurnoverID, c tv.TurnoverCycle) (tv.TurnoverCycle, error) {
	if err := c.Validate(); err != nil {
		return tv.TurnoverCycle{}, err
	}
	if _, err := s.GetTurnover(ctx, id); err != nil {
		return tv.TurnoverCycle{}, err
	}
	if c.ID == "" {
		c.ID = tv.NewTurnoverCycleID()
	}
	c.TurnoverID = id
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO turnover_cycles
		 (id, turnover_id, started_at, ended_at, advance_pence, returned_pence, production_minutes, circulation_minutes)
		 VALUES (?,?,?,?,?,?,?,?)`,
		string(c.ID), string(id),
		c.StartedAt.UTC().Format(time.RFC3339Nano), c.EndedAt.UTC().Format(time.RFC3339Nano),
		int64(c.AdvancePence), int64(c.ReturnedPence), c.ProductionMinutes, c.CirculationMinutes)
	if err != nil {
		return tv.TurnoverCycle{}, err
	}
	return c, nil
}

func (s *MySQL) RecomputeNumber(ctx context.Context, id tv.TurnoverID) (tv.TurnoverNumber, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT AVG(production_minutes + circulation_minutes) FROM turnover_cycles WHERE turnover_id=?`, string(id))
	var avgF sql.NullFloat64
	if err := row.Scan(&avgF); err != nil {
		return tv.TurnoverNumber{}, err
	}
	if !avgF.Valid || avgF.Float64 == 0 {
		return tv.TurnoverNumber{}, tv.ErrNoCompletedCycles
	}
	avg := int64(avgF.Float64)
	tn := tv.ComputeTurnoverNumber(id, avg)
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO turnover_numbers (id, turnover_id, year_reference_minutes, turnover_time_minutes, basis_points)
		 VALUES (?,?,?,?,?)
		 ON DUPLICATE KEY UPDATE
		   turnover_time_minutes=VALUES(turnover_time_minutes),
		   basis_points=VALUES(basis_points)`,
		string(tn.ID), string(id), tn.YearReferenceMinutes, tn.TurnoverTimeMinutes, tn.BasisPoints)
	if err != nil {
		return tv.TurnoverNumber{}, err
	}
	_, _ = s.db.ExecContext(ctx, `UPDATE turnovers SET turnover_time_minutes=? WHERE id=?`, avg, string(id))
	return tn, nil
}

func (s *MySQL) GetNumber(ctx context.Context, id tv.TurnoverID) (tv.TurnoverNumber, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, turnover_id, year_reference_minutes, turnover_time_minutes, basis_points
		 FROM turnover_numbers WHERE turnover_id=?`, string(id))
	var tn tv.TurnoverNumber
	if err := row.Scan(&tn.ID, &tn.TurnoverID, &tn.YearReferenceMinutes, &tn.TurnoverTimeMinutes, &tn.BasisPoints); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tv.TurnoverNumber{}, ErrNotFound
		}
		return tv.TurnoverNumber{}, err
	}
	return tn, nil
}
