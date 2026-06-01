package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// CreateAnnualValueComposition inserts c, assigning an ID and timestamp when absent (Vol. III Ch. 49).
func (m *MySQL) CreateAnnualValueComposition(ctx context.Context, c revenue.AnnualValueComposition) (revenue.AnnualValueComposition, error) {
	if c.ID == "" {
		c.ID = revenue.NewAnnualValueCompositionID()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `annual_value_compositions`" +
		" (`id`, `constant_capital_labour_minutes`, `variable_capital_labour_minutes`, `surplus_value_labour_minutes`, `total_value_labour_minutes`, `new_value_labour_minutes`, `created_at`)" +
		" VALUES (?,?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), c.ConstantCapitalLabourMinutes, c.VariableCapitalLabourMinutes,
		c.SurplusValueLabourMinutes, c.TotalValueLabourMinutes, c.NewValueLabourMinutes, c.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return revenue.AnnualValueComposition{}, ErrAlreadyExists
		}
		return revenue.AnnualValueComposition{}, err
	}
	return c, nil
}

// GetAnnualValueComposition returns the composition with the given ID, or ErrNotFound (Vol. III Ch. 49).
func (m *MySQL) GetAnnualValueComposition(ctx context.Context, id revenue.AnnualValueCompositionID) (revenue.AnnualValueComposition, error) {
	const q = "SELECT `id`, `constant_capital_labour_minutes`, `variable_capital_labour_minutes`, `surplus_value_labour_minutes`, `total_value_labour_minutes`, `new_value_labour_minutes`, `created_at`" +
		" FROM `annual_value_compositions` WHERE `id` = ?"
	return scanAnnualValueComposition(m.db.QueryRowContext(ctx, q, string(id)))
}

func scanAnnualValueComposition(s rowScanner) (revenue.AnnualValueComposition, error) {
	var (
		id        string
		constant  int64
		variable  int64
		surplus   int64
		total     int64
		newValue  int64
		createdAt time.Time
	)
	if err := s.Scan(&id, &constant, &variable, &surplus, &total, &newValue, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return revenue.AnnualValueComposition{}, ErrNotFound
		}
		return revenue.AnnualValueComposition{}, err
	}
	return revenue.AnnualValueComposition{
		ID:                           revenue.AnnualValueCompositionID(id),
		ConstantCapitalLabourMinutes: constant,
		VariableCapitalLabourMinutes: variable,
		SurplusValueLabourMinutes:    surplus,
		TotalValueLabourMinutes:      total,
		NewValueLabourMinutes:        newValue,
		CreatedAt:                    createdAt,
	}, nil
}

// CreateRevenueLimit inserts l, assigning an ID and timestamp when absent (Vol. III Ch. 49).
func (m *MySQL) CreateRevenueLimit(ctx context.Context, l revenue.RevenueLimit) (revenue.RevenueLimit, error) {
	if l.ID == "" {
		l.ID = revenue.NewRevenueLimitID()
	}
	if l.CreatedAt.IsZero() {
		l.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `revenue_limits`" +
		" (`id`, `composition_id`, `max_distributable_revenue_lm`, `constant_capital_replacement`, `created_at`)" +
		" VALUES (?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q,
		string(l.ID), l.CompositionID, l.MaxDistributableRevenueLM, l.ConstantCapitalReplacement, l.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return revenue.RevenueLimit{}, ErrAlreadyExists
		}
		return revenue.RevenueLimit{}, err
	}
	return l, nil
}

// GetRevenueLimitByComposition returns the revenue limit for the given composition ID, or ErrNotFound (Vol. III Ch. 49).
func (m *MySQL) GetRevenueLimitByComposition(ctx context.Context, compositionID string) (revenue.RevenueLimit, error) {
	const q = "SELECT `id`, `composition_id`, `max_distributable_revenue_lm`, `constant_capital_replacement`, `created_at`" +
		" FROM `revenue_limits` WHERE `composition_id` = ? ORDER BY `created_at` ASC, `id` ASC LIMIT 1"
	return scanRevenueLimit(m.db.QueryRowContext(ctx, q, compositionID))
}

func scanRevenueLimit(s rowScanner) (revenue.RevenueLimit, error) {
	var (
		id            string
		compositionID string
		maxDist       int64
		replacement   int64
		createdAt     time.Time
	)
	if err := s.Scan(&id, &compositionID, &maxDist, &replacement, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return revenue.RevenueLimit{}, ErrNotFound
		}
		return revenue.RevenueLimit{}, err
	}
	return revenue.RevenueLimit{
		ID:                         revenue.RevenueLimitID(id),
		CompositionID:              compositionID,
		MaxDistributableRevenueLM:  maxDist,
		ConstantCapitalReplacement: replacement,
		CreatedAt:                  createdAt,
	}, nil
}

// CreateSurplusValuePartition inserts p, assigning an ID and timestamp when absent (Vol. III Ch. 49).
func (m *MySQL) CreateSurplusValuePartition(ctx context.Context, p revenue.SurplusValuePartition) (revenue.SurplusValuePartition, error) {
	if p.ID == "" {
		p.ID = revenue.NewSurplusValuePartitionID()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `surplus_value_partitions`" +
		" (`id`, `composition_id`, `total_surplus_value_lm`, `profit_enterprise_lm`, `interest_lm`, `ground_rent_lm`, `residual_surplus_lm`, `created_at`)" +
		" VALUES (?,?,?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q,
		string(p.ID), p.CompositionID, p.TotalSurplusValueLM, p.ProfitEnterpriseLM,
		p.InterestLM, p.GroundRentLM, p.ResidualSurplusLM, p.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return revenue.SurplusValuePartition{}, ErrAlreadyExists
		}
		return revenue.SurplusValuePartition{}, err
	}
	return p, nil
}

// ListSurplusValuePartitions returns all stored partitions, newest first (Vol. III Ch. 49).
func (m *MySQL) ListSurplusValuePartitions(ctx context.Context) ([]revenue.SurplusValuePartition, error) {
	const q = "SELECT `id`, `composition_id`, `total_surplus_value_lm`, `profit_enterprise_lm`, `interest_lm`, `ground_rent_lm`, `residual_surplus_lm`, `created_at`" +
		" FROM `surplus_value_partitions` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.SurplusValuePartition, 0)
	for rows.Next() {
		p, err := scanSurplusValuePartition(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func scanSurplusValuePartition(s rowScanner) (revenue.SurplusValuePartition, error) {
	var (
		id            string
		compositionID string
		total         int64
		profitEnt     int64
		interest      int64
		groundRent    int64
		residual      int64
		createdAt     time.Time
	)
	if err := s.Scan(&id, &compositionID, &total, &profitEnt, &interest, &groundRent, &residual, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return revenue.SurplusValuePartition{}, ErrNotFound
		}
		return revenue.SurplusValuePartition{}, err
	}
	return revenue.SurplusValuePartition{
		ID:                  revenue.SurplusValuePartitionID(id),
		CompositionID:       compositionID,
		TotalSurplusValueLM: total,
		ProfitEnterpriseLM:  profitEnt,
		InterestLM:          interest,
		GroundRentLM:        groundRent,
		ResidualSurplusLM:   residual,
		CreatedAt:           createdAt,
	}, nil
}
