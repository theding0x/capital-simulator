package store

import (
	"context"
	"database/sql"
	"time"

	wp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

// RecordProductionLabourGap inserts a ProductionLabourGap row.
// One gap per working period; duplicate working_period_id returns ErrAlreadyExists.
func (m *MySQL) RecordProductionLabourGap(ctx context.Context, gap wp.ProductionLabourGap) (wp.ProductionLabourGap, error) {
	if err := gap.Validate(); err != nil {
		return wp.ProductionLabourGap{}, err
	}
	if gap.ID == "" {
		gap.ID = wp.NewProductionLabourGapID()
	}
	if gap.CreatedAt.IsZero() {
		gap.CreatedAt = m.now().UTC()
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO production_labour_gaps
		   (id, working_period_id, production_time_nanos, labour_time_nanos, gap_nanos, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		string(gap.ID), string(gap.WorkingPeriodID),
		gap.ProductionTimeNanos, gap.LabourTimeNanos, gap.GapNanos,
		gap.CreatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return wp.ProductionLabourGap{}, ErrAlreadyExists
		}
		return wp.ProductionLabourGap{}, err
	}
	return gap, nil
}

// GetProductionLabourGap returns the gap registered for a working period.
func (m *MySQL) GetProductionLabourGap(ctx context.Context, wpID wp.WorkingPeriodID) (wp.ProductionLabourGap, error) {
	row := m.db.QueryRowContext(ctx,
		`SELECT id, working_period_id, production_time_nanos, labour_time_nanos, gap_nanos, created_at
		 FROM production_labour_gaps
		 WHERE working_period_id = ?`,
		string(wpID))
	return scanProductionLabourGap(row)
}

func scanProductionLabourGap(s rowScanner) (wp.ProductionLabourGap, error) {
	var g wp.ProductionLabourGap
	var id, wpID string
	var createdAt time.Time
	if err := s.Scan(&id, &wpID, &g.ProductionTimeNanos, &g.LabourTimeNanos, &g.GapNanos, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return wp.ProductionLabourGap{}, ErrNotFound
		}
		return wp.ProductionLabourGap{}, err
	}
	g.ID = wp.ProductionLabourGapID(id)
	g.WorkingPeriodID = wp.WorkingPeriodID(wpID)
	g.CreatedAt = createdAt.UTC()
	return g, nil
}

// RecordNaturalProcessActivation inserts a NaturalProcessActivation row.
func (m *MySQL) RecordNaturalProcessActivation(ctx context.Context, act wp.NaturalProcessActivation) (wp.NaturalProcessActivation, error) {
	if err := act.Validate(); err != nil {
		return wp.NaturalProcessActivation{}, err
	}
	if act.ID == "" {
		act.ID = wp.NewNaturalProcessActivationID()
	}
	if act.StartedAt.IsZero() {
		act.StartedAt = m.now().UTC()
	}
	var endedAt interface{}
	if act.EndedAt != nil {
		endedAt = act.EndedAt.UTC()
	}
	maintenance := 0
	if act.RequiresLabourMaintenance {
		maintenance = 1
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO natural_process_activations
		   (id, working_period_id, kind, started_at, ended_at, requires_labour_maintenance)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		string(act.ID), string(act.WorkingPeriodID), string(act.Kind),
		act.StartedAt, endedAt, maintenance)
	if err != nil {
		if isDuplicateKey(err) {
			return wp.NaturalProcessActivation{}, ErrAlreadyExists
		}
		return wp.NaturalProcessActivation{}, err
	}
	return act, nil
}

// ListNaturalProcessActivations returns all activations for a working period ordered by started_at.
func (m *MySQL) ListNaturalProcessActivations(ctx context.Context, wpID wp.WorkingPeriodID) ([]wp.NaturalProcessActivation, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT id, working_period_id, kind, started_at, ended_at, requires_labour_maintenance
		 FROM natural_process_activations
		 WHERE working_period_id = ?
		 ORDER BY started_at ASC`,
		string(wpID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []wp.NaturalProcessActivation{}
	for rows.Next() {
		var a wp.NaturalProcessActivation
		var id, wid, kind string
		var startedAt time.Time
		var endedAt sql.NullTime
		var maintenance int
		if err := rows.Scan(&id, &wid, &kind, &startedAt, &endedAt, &maintenance); err != nil {
			return nil, err
		}
		a.ID = wp.NaturalProcessActivationID(id)
		a.WorkingPeriodID = wp.WorkingPeriodID(wid)
		a.Kind = wp.NaturalProcessKind(kind)
		a.StartedAt = startedAt.UTC()
		if endedAt.Valid {
			t := endedAt.Time.UTC()
			a.EndedAt = &t
		}
		a.RequiresLabourMaintenance = maintenance != 0
		out = append(out, a)
	}
	return out, rows.Err()
}

// RecordNaturalSubject inserts a NaturalSubject row.
// One subject per working period; duplicate working_period_id returns ErrAlreadyExists.
func (m *MySQL) RecordNaturalSubject(ctx context.Context, sub wp.NaturalSubject) (wp.NaturalSubject, error) {
	if err := sub.Validate(); err != nil {
		return wp.NaturalSubject{}, err
	}
	if sub.ID == "" {
		sub.ID = wp.NewNaturalSubjectID()
	}
	if sub.CreatedAt.IsZero() {
		sub.CreatedAt = m.now().UTC()
	}
	isLiving := 0
	if sub.IsLiving {
		isLiving = 1
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO natural_subjects
		   (id, working_period_id, commodity_id, description, is_living, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		string(sub.ID), string(sub.WorkingPeriodID), sub.CommodityID,
		sub.Description, isLiving, sub.CreatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return wp.NaturalSubject{}, ErrAlreadyExists
		}
		return wp.NaturalSubject{}, err
	}
	return sub, nil
}

// GetNaturalSubject returns the subject registered for a working period.
func (m *MySQL) GetNaturalSubject(ctx context.Context, wpID wp.WorkingPeriodID) (wp.NaturalSubject, error) {
	row := m.db.QueryRowContext(ctx,
		`SELECT id, working_period_id, commodity_id, description, is_living, created_at
		 FROM natural_subjects
		 WHERE working_period_id = ?`,
		string(wpID))
	var s wp.NaturalSubject
	var id, wid string
	var isLiving int
	var createdAt time.Time
	if err := row.Scan(&id, &wid, &s.CommodityID, &s.Description, &isLiving, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return wp.NaturalSubject{}, ErrNotFound
		}
		return wp.NaturalSubject{}, err
	}
	s.ID = wp.NaturalSubjectID(id)
	s.WorkingPeriodID = wp.WorkingPeriodID(wid)
	s.IsLiving = isLiving != 0
	s.CreatedAt = createdAt.UTC()
	return s, nil
}

// CreateIndustryBenchmark inserts a NaturalProcessIndustry row.
func (m *MySQL) CreateIndustryBenchmark(ctx context.Context, ind wp.NaturalProcessIndustry) (wp.NaturalProcessIndustry, error) {
	if err := ind.Validate(); err != nil {
		return wp.NaturalProcessIndustry{}, err
	}
	if ind.ID == "" {
		ind.ID = wp.NewNaturalProcessIndustryID()
	}
	if ind.CreatedAt.IsZero() {
		ind.CreatedAt = m.now().UTC()
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO natural_process_industries
		   (id, industry_name, typical_production_days_count, typical_labour_days_count, ratio_basis_points, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		string(ind.ID), ind.IndustryName,
		ind.TypicalProductionDaysCount, ind.TypicalLabourDaysCount,
		ind.RatioBasisPoints, ind.CreatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return wp.NaturalProcessIndustry{}, ErrAlreadyExists
		}
		return wp.NaturalProcessIndustry{}, err
	}
	return ind, nil
}

// ListIndustryBenchmarks returns all industry benchmarks ordered by ratio_basis_points ASC
// (most latent-capital-intensive industries first).
func (m *MySQL) ListIndustryBenchmarks(ctx context.Context) ([]wp.NaturalProcessIndustry, error) {
	rows, err := m.db.QueryContext(ctx,
		`SELECT id, industry_name, typical_production_days_count, typical_labour_days_count, ratio_basis_points, created_at
		 FROM natural_process_industries
		 ORDER BY ratio_basis_points ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []wp.NaturalProcessIndustry{}
	for rows.Next() {
		var ind wp.NaturalProcessIndustry
		var id string
		var createdAt time.Time
		if err := rows.Scan(&id, &ind.IndustryName,
			&ind.TypicalProductionDaysCount, &ind.TypicalLabourDaysCount,
			&ind.RatioBasisPoints, &createdAt); err != nil {
			return nil, err
		}
		ind.ID = wp.NaturalProcessIndustryID(id)
		ind.CreatedAt = createdAt.UTC()
		out = append(out, ind)
	}
	return out, rows.Err()
}
