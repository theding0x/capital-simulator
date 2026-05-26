package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	comp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
	tv "github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

func (s *MySQL) CreateAggregateTurnover(ctx context.Context, at tv.AggregateTurnover) (tv.AggregateTurnover, error) {
	if err := at.Validate(); err != nil {
		return tv.AggregateTurnover{}, err
	}
	if at.ID == "" {
		at.ID = tv.NewAggregateTurnoverID()
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO aggregate_turnovers (id, industrial_capital_id, lens, advanced_total_pence, annual_turned_over_pence, aggregate_number_bp)
		 VALUES (?,?,?,?,?,?)`,
		string(at.ID), at.IndustrialCapitalID, string(at.Lens),
		int64(at.AdvancedTotalPence), int64(at.AnnualTurnedOverPence), at.AggregateNumberBasisPoints)
	if err != nil {
		if isDuplicateKey(err) {
			return tv.AggregateTurnover{}, ErrAlreadyExists
		}
		return tv.AggregateTurnover{}, err
	}
	at.Contributions = []tv.ComponentTurnoverContribution{}
	return at, nil
}

func (s *MySQL) GetAggregateTurnover(ctx context.Context, id tv.AggregateTurnoverID) (tv.AggregateTurnover, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, industrial_capital_id, lens, advanced_total_pence, annual_turned_over_pence, aggregate_number_bp
		 FROM aggregate_turnovers WHERE id=?`, string(id))
	var at tv.AggregateTurnover
	var lens string
	if err := row.Scan(&at.ID, &at.IndustrialCapitalID, &lens, &at.AdvancedTotalPence, &at.AnnualTurnedOverPence, &at.AggregateNumberBasisPoints); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tv.AggregateTurnover{}, ErrNotFound
		}
		return tv.AggregateTurnover{}, err
	}
	at.Lens = tv.CircuitLens(lens)

	contribs, err := s.listContributions(ctx, id)
	if err != nil {
		return tv.AggregateTurnover{}, err
	}
	at.Contributions = contribs

	lc, err := s.getLifetimeCycle(ctx, id)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return tv.AggregateTurnover{}, err
	}
	if err == nil {
		phases, err := s.listCrisisPhases(ctx, lc.ID)
		if err != nil {
			return tv.AggregateTurnover{}, err
		}
		lc.CrisisPhases = phases
		at.LifetimeCycle = &lc
		if len(phases) > 0 {
			latest := phases[len(phases)-1]
			at.LatestCrisisPhase = &latest
		}
	}
	return at, nil
}

func (s *MySQL) ListAggregateTurnovers(ctx context.Context, industrialCapitalID string) ([]tv.AggregateTurnover, error) {
	q := `SELECT id, industrial_capital_id, lens, advanced_total_pence, annual_turned_over_pence, aggregate_number_bp FROM aggregate_turnovers`
	args := []any{}
	if industrialCapitalID != "" {
		q += ` WHERE industrial_capital_id=?`
		args = append(args, industrialCapitalID)
	}
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []tv.AggregateTurnover
	for rows.Next() {
		var at tv.AggregateTurnover
		var lens string
		if err := rows.Scan(&at.ID, &at.IndustrialCapitalID, &lens, &at.AdvancedTotalPence, &at.AnnualTurnedOverPence, &at.AggregateNumberBasisPoints); err != nil {
			return nil, err
		}
		at.Lens = tv.CircuitLens(lens)
		contribs, err := s.listContributions(ctx, at.ID)
		if err != nil {
			return nil, err
		}
		at.Contributions = contribs
		out = append(out, at)
	}
	if out == nil {
		out = []tv.AggregateTurnover{}
	}
	return out, rows.Err()
}

func (s *MySQL) RecordContribution(ctx context.Context, id tv.AggregateTurnoverID, c tv.ComponentTurnoverContribution) (tv.ComponentTurnoverContribution, error) {
	var exists int
	if err := s.db.QueryRowContext(ctx, `SELECT 1 FROM aggregate_turnovers WHERE id=?`, string(id)).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tv.ComponentTurnoverContribution{}, ErrNotFound
		}
		return tv.ComponentTurnoverContribution{}, err
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
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO component_turnover_contributions (id, aggregate_turnover_id, capital_component_id, kind, advanced_pence, turnover_number_bp, annual_return_pence, difference_kind)
		 VALUES (?,?,?,?,?,?,?,?)`,
		string(c.ID), string(id), string(c.CapitalComponentID), string(c.Kind),
		int64(c.AdvancedPence), c.TurnoverNumberBasisPoints, int64(c.AnnualReturnPence), string(c.DifferenceKind))
	if err != nil {
		if isDuplicateKey(err) {
			return tv.ComponentTurnoverContribution{}, ErrAlreadyExists
		}
		return tv.ComponentTurnoverContribution{}, err
	}
	return c, nil
}

func (s *MySQL) RecomputeAggregate(ctx context.Context, id tv.AggregateTurnoverID) (tv.AggregateTurnover, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return tv.AggregateTurnover{}, err
	}
	defer func() { _ = tx.Rollback() }()

	row := tx.QueryRowContext(ctx,
		`SELECT id, industrial_capital_id, lens, advanced_total_pence, annual_turned_over_pence, aggregate_number_bp
		 FROM aggregate_turnovers WHERE id=? FOR UPDATE`, string(id))
	var at tv.AggregateTurnover
	var lens string
	if err := row.Scan(&at.ID, &at.IndustrialCapitalID, &lens, &at.AdvancedTotalPence, &at.AnnualTurnedOverPence, &at.AggregateNumberBasisPoints); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tv.AggregateTurnover{}, ErrNotFound
		}
		return tv.AggregateTurnover{}, err
	}
	at.Lens = tv.CircuitLens(lens)

	contribs, err := s.listContributionsTx(ctx, tx, id)
	if err != nil {
		return tv.AggregateTurnover{}, err
	}

	at = tv.RecomputeAggregate(at, contribs)

	_, err = tx.ExecContext(ctx,
		`UPDATE aggregate_turnovers SET advanced_total_pence=?, annual_turned_over_pence=?, aggregate_number_bp=? WHERE id=?`,
		int64(at.AdvancedTotalPence), int64(at.AnnualTurnedOverPence), at.AggregateNumberBasisPoints, string(id))
	if err != nil {
		return tv.AggregateTurnover{}, err
	}
	if err := tx.Commit(); err != nil {
		return tv.AggregateTurnover{}, err
	}
	at.Contributions = contribs
	return at, nil
}

func (s *MySQL) RecordLifetimeCycle(ctx context.Context, id tv.AggregateTurnoverID, lc tv.LifetimeCycle) (tv.LifetimeCycle, error) {
	var exists int
	if err := s.db.QueryRowContext(ctx, `SELECT 1 FROM aggregate_turnovers WHERE id=?`, string(id)).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tv.LifetimeCycle{}, ErrNotFound
		}
		return tv.LifetimeCycle{}, err
	}
	if lc.ID == "" {
		lc.ID = tv.NewLifetimeCycleID()
	}
	lc.AggregateTurnoverID = id
	if lc.StartedAt.IsZero() {
		lc.StartedAt = time.Now().UTC()
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO lifetime_cycles (id, aggregate_turnover_id, duration_minutes, started_at, completed_cycles)
		 VALUES (?,?,?,?,?)
		 ON DUPLICATE KEY UPDATE duration_minutes=VALUES(duration_minutes), started_at=VALUES(started_at)`,
		string(lc.ID), string(id), lc.DurationMinutes, lc.StartedAt.UTC(), lc.CompletedCycles)
	if err != nil {
		return tv.LifetimeCycle{}, err
	}
	return lc, nil
}

func (s *MySQL) RecordCrisisPhase(ctx context.Context, cycleID tv.LifetimeCycleID, phase tv.CrisisCyclePhaseRecord) (tv.CrisisCyclePhaseRecord, error) {
	var exists int
	if err := s.db.QueryRowContext(ctx, `SELECT 1 FROM lifetime_cycles WHERE id=?`, string(cycleID)).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tv.CrisisCyclePhaseRecord{}, ErrNotFound
		}
		return tv.CrisisCyclePhaseRecord{}, err
	}
	if phase.ID == "" {
		phase.ID = tv.NewCrisisCyclePhaseRecordID()
	}
	phase.LifetimeCycleID = cycleID
	if phase.ObservedAt.IsZero() {
		phase.ObservedAt = time.Now().UTC()
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO crisis_cycle_phase_records (id, lifetime_cycle_id, phase, observed_at, notes)
		 VALUES (?,?,?,?,?)`,
		string(phase.ID), string(cycleID), string(phase.Phase), phase.ObservedAt.UTC(), phase.Notes)
	if err != nil {
		if isDuplicateKey(err) {
			return tv.CrisisCyclePhaseRecord{}, ErrAlreadyExists
		}
		return tv.CrisisCyclePhaseRecord{}, err
	}
	return phase, nil
}

// --- helpers -----------------------------------------------------------------

func (s *MySQL) listContributions(ctx context.Context, id tv.AggregateTurnoverID) ([]tv.ComponentTurnoverContribution, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, aggregate_turnover_id, capital_component_id, kind, advanced_pence, turnover_number_bp, annual_return_pence, difference_kind
		 FROM component_turnover_contributions WHERE aggregate_turnover_id=? ORDER BY recorded_at`, string(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanContributions(rows)
}

func (s *MySQL) listContributionsTx(ctx context.Context, tx *sql.Tx, id tv.AggregateTurnoverID) ([]tv.ComponentTurnoverContribution, error) {
	rows, err := tx.QueryContext(ctx,
		`SELECT id, aggregate_turnover_id, capital_component_id, kind, advanced_pence, turnover_number_bp, annual_return_pence, difference_kind
		 FROM component_turnover_contributions WHERE aggregate_turnover_id=? ORDER BY recorded_at`, string(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanContributions(rows)
}

func scanContributions(rows *sql.Rows) ([]tv.ComponentTurnoverContribution, error) {
	var out []tv.ComponentTurnoverContribution
	for rows.Next() {
		var c tv.ComponentTurnoverContribution
		var kind, diffKind, compID string
		if err := rows.Scan(&c.ID, &c.AggregateTurnoverID, &compID, &kind, &c.AdvancedPence, &c.TurnoverNumberBasisPoints, &c.AnnualReturnPence, &diffKind); err != nil {
			return nil, err
		}
		c.CapitalComponentID = comp.CapitalComponentID(compID)
		c.Kind = comp.CapitalKind(kind)
		c.DifferenceKind = tv.TurnoverDifferenceKind(diffKind)
		out = append(out, c)
	}
	if out == nil {
		out = []tv.ComponentTurnoverContribution{}
	}
	return out, rows.Err()
}

func (s *MySQL) getLifetimeCycle(ctx context.Context, id tv.AggregateTurnoverID) (tv.LifetimeCycle, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, aggregate_turnover_id, duration_minutes, started_at, completed_cycles
		 FROM lifetime_cycles WHERE aggregate_turnover_id=?`, string(id))
	var lc tv.LifetimeCycle
	var startedAt []byte
	if err := row.Scan(&lc.ID, &lc.AggregateTurnoverID, &lc.DurationMinutes, &startedAt, &lc.CompletedCycles); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tv.LifetimeCycle{}, ErrNotFound
		}
		return tv.LifetimeCycle{}, err
	}
	t, err := time.Parse("2006-01-02 15:04:05.999999", string(startedAt))
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05Z", string(startedAt))
		if err != nil {
			return tv.LifetimeCycle{}, err
		}
	}
	lc.StartedAt = t.UTC()
	return lc, nil
}

func (s *MySQL) listCrisisPhases(ctx context.Context, cycleID tv.LifetimeCycleID) ([]tv.CrisisCyclePhaseRecord, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, lifetime_cycle_id, phase, observed_at, notes
		 FROM crisis_cycle_phase_records WHERE lifetime_cycle_id=? ORDER BY observed_at`, string(cycleID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []tv.CrisisCyclePhaseRecord
	for rows.Next() {
		var p tv.CrisisCyclePhaseRecord
		var phase string
		var observedAt []byte
		if err := rows.Scan(&p.ID, &p.LifetimeCycleID, &phase, &observedAt, &p.Notes); err != nil {
			return nil, err
		}
		p.Phase = tv.CrisisCyclePhase(phase)
		t, err := time.Parse("2006-01-02 15:04:05.999999", string(observedAt))
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05Z", string(observedAt))
			if err != nil {
				return nil, err
			}
		}
		p.ObservedAt = t.UTC()
		out = append(out, p)
	}
	if out == nil {
		out = []tv.CrisisCyclePhaseRecord{}
	}
	return out, rows.Err()
}
