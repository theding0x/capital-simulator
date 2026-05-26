package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	wp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

func (m *MySQL) CreateWorkingPeriod(ctx context.Context, w wp.WorkingPeriod) (wp.WorkingPeriod, error) {
	if err := w.Validate(); err != nil {
		return wp.WorkingPeriod{}, err
	}
	// Check natural constraint floor.
	row := m.db.QueryRowContext(ctx,
		`SELECT min_working_days FROM natural_working_period_constraints WHERE commodity_id = ?`,
		w.CommodityID)
	var minDays int64
	if err := row.Scan(&minDays); err == nil {
		if w.WorkingDayCount < minDays {
			return wp.WorkingPeriod{}, wp.ErrPrematureDelivery
		}
	}
	if w.ID == "" {
		w.ID = wp.NewWorkingPeriodID()
	}
	if w.CreatedAt.IsZero() {
		w.CreatedAt = m.now().UTC()
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO working_periods (id, industrial_capital_id, commodity_id, working_day_count, mode, working_day_length_minutes, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		string(w.ID), w.IndustrialCapitalID, w.CommodityID,
		w.WorkingDayCount, string(w.Mode), w.WorkingDayLengthMinutes, w.CreatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return wp.WorkingPeriod{}, ErrAlreadyExists
		}
		return wp.WorkingPeriod{}, err
	}
	return w, nil
}

func (m *MySQL) GetWorkingPeriod(ctx context.Context, id wp.WorkingPeriodID) (wp.WorkingPeriod, error) {
	row := m.db.QueryRowContext(ctx,
		`SELECT id, industrial_capital_id, commodity_id, working_day_count, mode, working_day_length_minutes, created_at
		 FROM working_periods WHERE id = ?`, string(id))
	w, err := scanWorkingPeriod(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wp.WorkingPeriod{}, ErrNotFound
		}
		return wp.WorkingPeriod{}, err
	}
	intrs, err := m.listInterruptions(ctx, id)
	if err != nil {
		return wp.WorkingPeriod{}, err
	}
	w.Interruptions = intrs
	shorts, err := m.listShortenings(ctx, id)
	if err != nil {
		return wp.WorkingPeriod{}, err
	}
	w.Shortenings = shorts
	cf, err := m.getCreditFinancing(ctx, id)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return wp.WorkingPeriod{}, err
	}
	if err == nil {
		w.CreditFinancing = &cf
	}
	return w, nil
}

func (m *MySQL) ListWorkingPeriods(ctx context.Context, industrialCapitalID, mode, commodityID string) ([]wp.WorkingPeriod, error) {
	q := `SELECT id, industrial_capital_id, commodity_id, working_day_count, mode, working_day_length_minutes, created_at
	      FROM working_periods WHERE 1=1`
	var args []any
	if industrialCapitalID != "" {
		q += " AND industrial_capital_id = ?"
		args = append(args, industrialCapitalID)
	}
	if mode != "" {
		q += " AND mode = ?"
		args = append(args, mode)
	}
	if commodityID != "" {
		q += " AND commodity_id = ?"
		args = append(args, commodityID)
	}
	q += " ORDER BY created_at ASC"
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []wp.WorkingPeriod{}
	for rows.Next() {
		w, err := scanWorkingPeriod(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func (m *MySQL) RecordInterruption(ctx context.Context, id wp.WorkingPeriodID, intr wp.WorkingPeriodInterruption) (wp.WorkingPeriodInterruption, error) {
	row := m.db.QueryRowContext(ctx, `SELECT mode FROM working_periods WHERE id = ?`, string(id))
	var modeStr string
	if err := row.Scan(&modeStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wp.WorkingPeriodInterruption{}, ErrNotFound
		}
		return wp.WorkingPeriodInterruption{}, err
	}
	if err := intr.Validate(wp.WorkingPeriodMode(modeStr)); err != nil {
		return wp.WorkingPeriodInterruption{}, err
	}
	if intr.ID == "" {
		intr.ID = wp.NewWorkingPeriodInterruptionID()
	}
	intr.WorkingPeriodID = id
	if intr.InterruptedAt.IsZero() {
		intr.InterruptedAt = m.now().UTC()
	}
	var resumedAt any
	if intr.ResumedAt != nil {
		resumedAt = intr.ResumedAt.UTC()
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO working_period_interruptions
		 (id, working_period_id, interrupted_at, resumed_at, deterioration_pence, waste_of_means_of_production_pence)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		string(intr.ID), string(id), intr.InterruptedAt, resumedAt,
		intr.DeteriorationPence, intr.WasteOfMeansOfProductionPence)
	if err != nil {
		return wp.WorkingPeriodInterruption{}, err
	}
	return intr, nil
}

func (m *MySQL) RecordShortening(ctx context.Context, id wp.WorkingPeriodID, s wp.WorkingPeriodShortening) (wp.WorkingPeriodShortening, error) {
	var exists int
	if err := m.db.QueryRowContext(ctx, `SELECT 1 FROM working_periods WHERE id = ?`, string(id)).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wp.WorkingPeriodShortening{}, ErrNotFound
		}
		return wp.WorkingPeriodShortening{}, err
	}
	if err := s.Validate(); err != nil {
		return wp.WorkingPeriodShortening{}, err
	}
	if s.ID == "" {
		s.ID = wp.NewWorkingPeriodShorteningID()
	}
	s.WorkingPeriodID = id
	if s.OccurredAt.IsZero() {
		s.OccurredAt = m.now().UTC()
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO working_period_shortenings
		 (id, working_period_id, kind, old_working_day_count, new_working_day_count, additional_fixed_capital_pence, additional_labour_count, occurred_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		string(s.ID), string(id), string(s.Kind),
		s.OldWorkingDayCount, s.NewWorkingDayCount,
		s.AdditionalFixedCapitalPence, s.AdditionalLabourCount, s.OccurredAt)
	if err != nil {
		return wp.WorkingPeriodShortening{}, err
	}
	return s, nil
}

func (m *MySQL) RecordCreditFinancing(ctx context.Context, id wp.WorkingPeriodID, cf wp.WorkingPeriodCreditFinancing) (wp.WorkingPeriodCreditFinancing, error) {
	var exists int
	if err := m.db.QueryRowContext(ctx, `SELECT 1 FROM working_periods WHERE id = ?`, string(id)).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wp.WorkingPeriodCreditFinancing{}, ErrNotFound
		}
		return wp.WorkingPeriodCreditFinancing{}, err
	}
	if err := cf.Validate(); err != nil {
		return wp.WorkingPeriodCreditFinancing{}, err
	}
	if cf.ID == "" {
		cf.ID = wp.NewWorkingPeriodCreditFinancingID()
	}
	cf.WorkingPeriodID = id
	if cf.CreatedAt.IsZero() {
		cf.CreatedAt = m.now().UTC()
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO working_period_credit_financings
		 (id, working_period_id, own_capital_pence, borrowed_capital_pence, mortgage_provider_id, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE
		   own_capital_pence = VALUES(own_capital_pence),
		   borrowed_capital_pence = VALUES(borrowed_capital_pence),
		   mortgage_provider_id = VALUES(mortgage_provider_id)`,
		string(cf.ID), string(id),
		cf.OwnCapitalPence, cf.BorrowedCapitalPence, cf.MortgageProviderID, cf.CreatedAt)
	if err != nil {
		return wp.WorkingPeriodCreditFinancing{}, err
	}
	return cf, nil
}

func (m *MySQL) CreateNaturalConstraint(ctx context.Context, nc wp.NaturalWorkingPeriodConstraint) (wp.NaturalWorkingPeriodConstraint, error) {
	if err := nc.Validate(); err != nil {
		return wp.NaturalWorkingPeriodConstraint{}, err
	}
	if nc.ID == "" {
		nc.ID = wp.NewNaturalWorkingPeriodConstraintID()
	}
	if nc.CreatedAt.IsZero() {
		nc.CreatedAt = m.now().UTC()
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO natural_working_period_constraints (id, commodity_id, min_working_days, reason, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		string(nc.ID), nc.CommodityID, nc.MinWorkingDays, nc.Reason, nc.CreatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return wp.NaturalWorkingPeriodConstraint{}, ErrAlreadyExists
		}
		return wp.NaturalWorkingPeriodConstraint{}, err
	}
	return nc, nil
}

func (m *MySQL) ListNaturalConstraints(ctx context.Context) ([]wp.NaturalWorkingPeriodConstraint, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT id, commodity_id, min_working_days, reason, created_at
		 FROM natural_working_period_constraints ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []wp.NaturalWorkingPeriodConstraint{}
	for rows.Next() {
		var nc wp.NaturalWorkingPeriodConstraint
		if err := rows.Scan(&nc.ID, &nc.CommodityID, &nc.MinWorkingDays, &nc.Reason, &nc.CreatedAt); err != nil {
			return nil, err
		}
		nc.CreatedAt = nc.CreatedAt.UTC()
		out = append(out, nc)
	}
	return out, rows.Err()
}

// --- helpers ----------------------------------------------------------------

func scanWorkingPeriod(scan func(...any) error) (wp.WorkingPeriod, error) {
	var w wp.WorkingPeriod
	var modeStr string
	var createdAt time.Time
	if err := scan(&w.ID, &w.IndustrialCapitalID, &w.CommodityID,
		&w.WorkingDayCount, &modeStr, &w.WorkingDayLengthMinutes, &createdAt); err != nil {
		return wp.WorkingPeriod{}, err
	}
	w.Mode = wp.WorkingPeriodMode(modeStr)
	w.CreatedAt = createdAt.UTC()
	w.Interruptions = []wp.WorkingPeriodInterruption{}
	w.Shortenings = []wp.WorkingPeriodShortening{}
	return w, nil
}

func (m *MySQL) listInterruptions(ctx context.Context, id wp.WorkingPeriodID) ([]wp.WorkingPeriodInterruption, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT id, working_period_id, interrupted_at, resumed_at, deterioration_pence, waste_of_means_of_production_pence
		 FROM working_period_interruptions WHERE working_period_id = ? ORDER BY interrupted_at ASC`,
		string(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []wp.WorkingPeriodInterruption{}
	for rows.Next() {
		var intr wp.WorkingPeriodInterruption
		var resumedAt sql.NullTime
		if err := rows.Scan(&intr.ID, &intr.WorkingPeriodID,
			&intr.InterruptedAt, &resumedAt,
			&intr.DeteriorationPence, &intr.WasteOfMeansOfProductionPence); err != nil {
			return nil, err
		}
		if resumedAt.Valid {
			t := resumedAt.Time.UTC()
			intr.ResumedAt = &t
		}
		intr.InterruptedAt = intr.InterruptedAt.UTC()
		out = append(out, intr)
	}
	return out, rows.Err()
}

func (m *MySQL) listShortenings(ctx context.Context, id wp.WorkingPeriodID) ([]wp.WorkingPeriodShortening, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT id, working_period_id, kind, old_working_day_count, new_working_day_count,
		        additional_fixed_capital_pence, additional_labour_count, occurred_at
		 FROM working_period_shortenings WHERE working_period_id = ? ORDER BY occurred_at ASC`,
		string(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []wp.WorkingPeriodShortening{}
	for rows.Next() {
		var s wp.WorkingPeriodShortening
		var kindStr string
		if err := rows.Scan(&s.ID, &s.WorkingPeriodID, &kindStr,
			&s.OldWorkingDayCount, &s.NewWorkingDayCount,
			&s.AdditionalFixedCapitalPence, &s.AdditionalLabourCount, &s.OccurredAt); err != nil {
			return nil, err
		}
		s.Kind = wp.WorkingPeriodShorteningKind(kindStr)
		s.OccurredAt = s.OccurredAt.UTC()
		out = append(out, s)
	}
	return out, rows.Err()
}

func (m *MySQL) getCreditFinancing(ctx context.Context, id wp.WorkingPeriodID) (wp.WorkingPeriodCreditFinancing, error) {
	row := m.db.QueryRowContext(ctx,
		`SELECT id, working_period_id, own_capital_pence, borrowed_capital_pence, mortgage_provider_id, created_at
		 FROM working_period_credit_financings WHERE working_period_id = ?`, string(id))
	var cf wp.WorkingPeriodCreditFinancing
	if err := row.Scan(&cf.ID, &cf.WorkingPeriodID,
		&cf.OwnCapitalPence, &cf.BorrowedCapitalPence, &cf.MortgageProviderID, &cf.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wp.WorkingPeriodCreditFinancing{}, ErrNotFound
		}
		return wp.WorkingPeriodCreditFinancing{}, err
	}
	cf.CreatedAt = cf.CreatedAt.UTC()
	return cf, nil
}
