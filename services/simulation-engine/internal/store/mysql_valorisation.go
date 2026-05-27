package store

import (
	"context"
	"database/sql"

	val "github.com/theding0x/capital-simulator/services/simulation-engine/internal/valorisation"
)

// CreateAdvance inserts a new VariableCapitalAdvance.
func (m *MySQL) CreateAdvance(ctx context.Context, a val.VariableCapitalAdvance) (val.VariableCapitalAdvance, error) {
	if err := a.Validate(); err != nil {
		return val.VariableCapitalAdvance{}, err
	}
	if a.ID == "" {
		a.ID = val.NewVariableCapitalAdvanceID()
	}
	const q = `INSERT INTO variable_capital_advances
		(id, industrial_capital_id, pence, turnover_period_minutes,
		 turnovers_per_year_basis_points, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID), a.IndustrialCapitalID,
		int64(a.Pence), a.TurnoverPeriodMinutes,
		a.TurnoversPerYearBasisPoints, m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return val.VariableCapitalAdvance{}, ErrAlreadyExists
		}
		return val.VariableCapitalAdvance{}, err
	}
	return a, nil
}

// GetAdvance returns a VariableCapitalAdvance by ID.
func (m *MySQL) GetAdvance(ctx context.Context, id val.VariableCapitalAdvanceID) (val.VariableCapitalAdvance, error) {
	const q = `SELECT id, industrial_capital_id, pence,
		turnover_period_minutes, turnovers_per_year_basis_points
		FROM variable_capital_advances WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanAdvance(row.Scan)
}

// ListAdvances returns advances, optionally filtered by industrial_capital_id.
func (m *MySQL) ListAdvances(ctx context.Context, industrialCapitalID string) ([]val.VariableCapitalAdvance, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if industrialCapitalID == "" {
		rows, err = m.db.QueryContext(ctx,
			`SELECT id, industrial_capital_id, pence,
			 turnover_period_minutes, turnovers_per_year_basis_points
			 FROM variable_capital_advances ORDER BY created_at ASC`)
	} else {
		rows, err = m.db.QueryContext(ctx,
			`SELECT id, industrial_capital_id, pence,
			 turnover_period_minutes, turnovers_per_year_basis_points
			 FROM variable_capital_advances
			 WHERE industrial_capital_id = ?
			 ORDER BY created_at ASC`,
			industrialCapitalID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []val.VariableCapitalAdvance
	for rows.Next() {
		a, err := scanAdvance(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	if out == nil {
		out = []val.VariableCapitalAdvance{}
	}
	return out, nil
}

func scanAdvance(scan func(...any) error) (val.VariableCapitalAdvance, error) {
	var (
		id                          string
		industrialCapitalID         string
		pence                       int64
		turnoverPeriodMinutes       int64
		turnoversPerYearBasisPoints int64
	)
	if err := scan(&id, &industrialCapitalID, &pence,
		&turnoverPeriodMinutes, &turnoversPerYearBasisPoints); err != nil {
		if err == sql.ErrNoRows {
			return val.VariableCapitalAdvance{}, ErrNotFound
		}
		return val.VariableCapitalAdvance{}, err
	}
	return val.VariableCapitalAdvance{
		ID:                          val.VariableCapitalAdvanceID(id),
		IndustrialCapitalID:         industrialCapitalID,
		Pence:                       val.Pence(pence),
		TurnoverPeriodMinutes:       turnoverPeriodMinutes,
		TurnoversPerYearBasisPoints: turnoversPerYearBasisPoints,
	}, nil
}

// RecordPerCircuitRate inserts or replaces the per-circuit surplus rate.
func (m *MySQL) RecordPerCircuitRate(ctx context.Context, rate val.PerCircuitSurplusRate) (val.PerCircuitSurplusRate, error) {
	if rate.ID == "" {
		rate.ID = val.NewPerCircuitSurplusRateID()
	}
	const q = `INSERT INTO per_circuit_surplus_rates
		(id, variable_capital_advance_id, surplus_pence, variable_pence, basis_points)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		  id = VALUES(id),
		  surplus_pence = VALUES(surplus_pence),
		  variable_pence = VALUES(variable_pence),
		  basis_points = VALUES(basis_points)`
	_, err := m.db.ExecContext(ctx, q,
		string(rate.ID), string(rate.VariableCapitalAdvanceID),
		int64(rate.SurplusPence), int64(rate.VariablePence), rate.BasisPoints,
	)
	if err != nil {
		return val.PerCircuitSurplusRate{}, err
	}
	return rate, nil
}

// GetPerCircuitRate returns the per-circuit surplus rate for an advance.
func (m *MySQL) GetPerCircuitRate(ctx context.Context, advanceID val.VariableCapitalAdvanceID) (val.PerCircuitSurplusRate, error) {
	const q = `SELECT id, variable_capital_advance_id, surplus_pence, variable_pence, basis_points
		FROM per_circuit_surplus_rates WHERE variable_capital_advance_id = ?`
	row := m.db.QueryRowContext(ctx, q, string(advanceID))
	var (
		id     string
		advID  string
		sPence int64
		vPence int64
		bp     int64
	)
	if err := row.Scan(&id, &advID, &sPence, &vPence, &bp); err != nil {
		if err == sql.ErrNoRows {
			return val.PerCircuitSurplusRate{}, ErrNotFound
		}
		return val.PerCircuitSurplusRate{}, err
	}
	return val.PerCircuitSurplusRate{
		ID:                       val.PerCircuitSurplusRateID(id),
		VariableCapitalAdvanceID: val.VariableCapitalAdvanceID(advID),
		SurplusPence:             val.Pence(sPence),
		VariablePence:            val.Pence(vPence),
		BasisPoints:              bp,
	}, nil
}

// RecordAnnualRate inserts an AnnualSurplusRate snapshot.
func (m *MySQL) RecordAnnualRate(ctx context.Context, rate val.AnnualSurplusRate) (val.AnnualSurplusRate, error) {
	if rate.ID == "" {
		rate.ID = val.NewAnnualSurplusRateID()
	}
	const q = `INSERT INTO annual_surplus_rates
		(id, variable_capital_advance_id, period,
		 n_basis_points, per_circuit_basis_points, annual_basis_points)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(rate.ID), string(rate.VariableCapitalAdvanceID), rate.Period,
		rate.NBasisPoints, rate.PerCircuitBasisPoints, rate.AnnualBasisPoints,
	)
	if err != nil {
		if isDuplicateKey(err) {
			return val.AnnualSurplusRate{}, ErrAlreadyExists
		}
		return val.AnnualSurplusRate{}, err
	}
	return rate, nil
}

// ListAnnualRates returns annual-rate snapshots, optionally filtered.
func (m *MySQL) ListAnnualRates(ctx context.Context, industrialCapitalID, period string) ([]val.AnnualSurplusRate, error) {
	joinClause := ""
	var args []any
	whereClause := ""

	if industrialCapitalID != "" {
		joinClause = ` JOIN variable_capital_advances vca ON asr.variable_capital_advance_id = vca.id`
		whereClause = "WHERE vca.industrial_capital_id = ?"
		args = append(args, industrialCapitalID)
		if period != "" {
			whereClause += " AND asr.period = ?"
			args = append(args, period)
		}
	} else if period != "" {
		whereClause = "WHERE asr.period = ?"
		args = append(args, period)
	}

	query := `SELECT asr.id, asr.variable_capital_advance_id, asr.period,
		asr.n_basis_points, asr.per_circuit_basis_points, asr.annual_basis_points
		FROM annual_surplus_rates asr` +
		joinClause + " " + whereClause

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []val.AnnualSurplusRate
	for rows.Next() {
		var (
			id    string
			advID string
			per   string
			nBP   int64
			pcBP  int64
			annBP int64
		)
		if err := rows.Scan(&id, &advID, &per, &nBP, &pcBP, &annBP); err != nil {
			return nil, err
		}
		out = append(out, val.AnnualSurplusRate{
			ID:                       val.AnnualSurplusRateID(id),
			VariableCapitalAdvanceID: val.VariableCapitalAdvanceID(advID),
			Period:                   per,
			NBasisPoints:             nBP,
			PerCircuitBasisPoints:    pcBP,
			AnnualBasisPoints:        annBP,
		})
	}
	if out == nil {
		out = []val.AnnualSurplusRate{}
	}
	return out, nil
}

// RecordContrast inserts an AnnualSurplusRateContrast.
func (m *MySQL) RecordContrast(ctx context.Context, c val.AnnualSurplusRateContrast) (val.AnnualSurplusRateContrast, error) {
	if c.ID == "" {
		c.ID = val.NewAnnualSurplusRateContrastID()
	}
	const q = `INSERT INTO annual_surplus_rate_contrasts
		(id, capital_a_id, capital_b_id,
		 capital_a_annual_basis_points, capital_b_annual_basis_points, ratio_basis_points)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), string(c.CapitalAID), string(c.CapitalBID),
		c.CapitalAAnnualBasisPoints, c.CapitalBAnnualBasisPoints, c.RatioBasisPoints,
	)
	if err != nil {
		if isDuplicateKey(err) {
			return val.AnnualSurplusRateContrast{}, ErrAlreadyExists
		}
		return val.AnnualSurplusRateContrast{}, err
	}
	return c, nil
}
