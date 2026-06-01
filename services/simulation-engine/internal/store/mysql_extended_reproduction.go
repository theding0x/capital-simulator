package store

import (
	"context"
	"database/sql"
	"sort"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// Vol. II Ch. 21 — Accumulation and Reproduction on an Extended Scale.

// --- ExtendedReproductionScheme ---------------------------------------------

func (m *MySQL) CreateExtendedReproductionScheme(ctx context.Context, s repro.ExtendedReproductionScheme) (repro.ExtendedReproductionScheme, error) {
	if s.ID == "" {
		s.ID = repro.NewExtendedReproductionSchemeID()
	}
	s.IsBalanced = false
	s.TickCount = 0
	const q = `INSERT INTO extended_reproduction_schemes
		(id, period, accumulation_rate_i_bps, accumulation_rate_ii_bps, is_balanced, tick_count, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID), s.Period,
		s.AccumulationRateIBps, s.AccumulationRateIIBps,
		s.IsBalanced, s.TickCount, m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return repro.ExtendedReproductionScheme{}, ErrAlreadyExists
		}
		return repro.ExtendedReproductionScheme{}, err
	}
	return s, nil
}

func (m *MySQL) GetExtendedReproductionScheme(ctx context.Context, id repro.ExtendedReproductionSchemeID) (repro.ExtendedReproductionScheme, error) {
	const q = `SELECT id, period, accumulation_rate_i_bps, accumulation_rate_ii_bps, tick_count
		FROM extended_reproduction_schemes WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var s repro.ExtendedReproductionScheme
	var sid string
	if err := row.Scan(&sid, &s.Period, &s.AccumulationRateIBps, &s.AccumulationRateIIBps, &s.TickCount); err != nil {
		if err == sql.ErrNoRows {
			return repro.ExtendedReproductionScheme{}, ErrNotFound
		}
		return repro.ExtendedReproductionScheme{}, err
	}
	s.ID = repro.ExtendedReproductionSchemeID(sid)

	// Load departments.
	deptI, err := m.getExtendedDepartmentalCapital(ctx, id, repro.DepartmentI)
	if err == nil {
		s.DepartmentI = &deptI
	} else if err != ErrNotFound {
		return repro.ExtendedReproductionScheme{}, err
	}
	deptII, err := m.getExtendedDepartmentalCapital(ctx, id, repro.DepartmentII)
	if err == nil {
		s.DepartmentII = &deptII
	} else if err != ErrNotFound {
		return repro.ExtendedReproductionScheme{}, err
	}

	// Recompute IsBalanced if both departments present.
	if s.DepartmentI != nil && s.DepartmentII != nil {
		// Use a zero reinvestment for the balance check at rest (no pending tick).
		rI := repro.Reinvestment{}
		rII := repro.Reinvestment{}
		s.IsBalanced = repro.ComputeExtendedBalance(*s.DepartmentI, *s.DepartmentII, rI, rII)
	}
	return s, nil
}

func (m *MySQL) ListExtendedReproductionSchemes(ctx context.Context, period string) ([]repro.ExtendedReproductionScheme, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if period != "" {
		rows, err = m.db.QueryContext(ctx,
			`SELECT id, period, accumulation_rate_i_bps, accumulation_rate_ii_bps, tick_count
			FROM extended_reproduction_schemes WHERE period = ? ORDER BY id`, period)
	} else {
		rows, err = m.db.QueryContext(ctx,
			`SELECT id, period, accumulation_rate_i_bps, accumulation_rate_ii_bps, tick_count
			FROM extended_reproduction_schemes ORDER BY id`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []repro.ExtendedReproductionScheme
	for rows.Next() {
		var s repro.ExtendedReproductionScheme
		var sid string
		if err2 := rows.Scan(&sid, &s.Period, &s.AccumulationRateIBps, &s.AccumulationRateIIBps, &s.TickCount); err2 != nil {
			return nil, err2
		}
		s.ID = repro.ExtendedReproductionSchemeID(sid)
		id := s.ID
		deptI, e := m.getExtendedDepartmentalCapital(ctx, id, repro.DepartmentI)
		if e == nil {
			s.DepartmentI = &deptI
		} else if e != ErrNotFound {
			return nil, e
		}
		deptII, e := m.getExtendedDepartmentalCapital(ctx, id, repro.DepartmentII)
		if e == nil {
			s.DepartmentII = &deptII
		} else if e != ErrNotFound {
			return nil, e
		}
		if s.DepartmentI != nil && s.DepartmentII != nil {
			rI := repro.Reinvestment{}
			rII := repro.Reinvestment{}
			s.IsBalanced = repro.ComputeExtendedBalance(*s.DepartmentI, *s.DepartmentII, rI, rII)
		}
		out = append(out, s)
	}
	if out == nil {
		out = []repro.ExtendedReproductionScheme{}
	}
	return out, rows.Err()
}

// getExtendedDepartmentalCapital loads a single DepartmentalCapital from the
// extended_departmental_capitals table for a given scheme + department.
func (m *MySQL) getExtendedDepartmentalCapital(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID, dept repro.CapitalDepartment) (repro.DepartmentalCapital, error) {
	const q = `SELECT id, scheme_id, department, constant_pence, variable_pence,
		surplus_pence, total_pence, turnover_number_bp
		FROM extended_departmental_capitals WHERE scheme_id = ? AND department = ?`
	row := m.db.QueryRowContext(ctx, q, string(schemeID), string(dept))
	var d repro.DepartmentalCapital
	var did, sid, department string
	if err := row.Scan(&did, &sid, &department,
		&d.ConstantPence, &d.VariablePence, &d.SurplusPence, &d.TotalPence,
		&d.TurnoverNumberBP); err != nil {
		if err == sql.ErrNoRows {
			return repro.DepartmentalCapital{}, ErrNotFound
		}
		return repro.DepartmentalCapital{}, err
	}
	d.ID = repro.DepartmentalCapitalID(did)
	d.Department = repro.CapitalDepartment(department)
	return d, nil
}

// --- DepartmentalCapital ----------------------------------------------------

func (m *MySQL) AddExtendedDepartmentToScheme(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID, dept repro.DepartmentalCapital) (repro.DepartmentalCapital, error) {
	if err := dept.Validate(); err != nil {
		return repro.DepartmentalCapital{}, err
	}
	if dept.ID == "" {
		dept.ID = repro.NewDepartmentalCapitalID()
	}

	// Check that the scheme exists.
	var dummy string
	if err := m.db.QueryRowContext(ctx,
		`SELECT id FROM extended_reproduction_schemes WHERE id = ?`, string(schemeID),
	).Scan(&dummy); err != nil {
		if err == sql.ErrNoRows {
			return repro.DepartmentalCapital{}, ErrNotFound
		}
		return repro.DepartmentalCapital{}, err
	}

	// Check for an existing department row.
	if _, err := m.getExtendedDepartmentalCapital(ctx, schemeID, dept.Department); err == nil {
		return repro.DepartmentalCapital{}, ErrAlreadyExists
	} else if err != ErrNotFound {
		return repro.DepartmentalCapital{}, err
	}

	const q = `INSERT INTO extended_departmental_capitals
		(id, scheme_id, department, constant_pence, variable_pence, surplus_pence, total_pence, turnover_number_bp, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := m.db.ExecContext(ctx, q,
		string(dept.ID), string(schemeID), string(dept.Department),
		dept.ConstantPence, dept.VariablePence, dept.SurplusPence, dept.TotalPence,
		dept.TurnoverNumberBP, m.now(),
	); err != nil {
		if isDuplicateKey(err) {
			return repro.DepartmentalCapital{}, ErrAlreadyExists
		}
		return repro.DepartmentalCapital{}, err
	}
	return dept, nil
}

// --- Reinvestment -----------------------------------------------------------

func (m *MySQL) CreateReinvestment(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID, r repro.Reinvestment) (repro.Reinvestment, error) {
	if r.ID == "" {
		r.ID = repro.NewReinvestmentID()
	}
	r.SchemeID = schemeID

	// Verify scheme exists.
	var dummy string
	if err := m.db.QueryRowContext(ctx,
		`SELECT id FROM extended_reproduction_schemes WHERE id = ?`, string(schemeID),
	).Scan(&dummy); err != nil {
		if err == sql.ErrNoRows {
			return repro.Reinvestment{}, ErrNotFound
		}
		return repro.Reinvestment{}, err
	}

	const q = `INSERT INTO reinvestments
		(id, scheme_id, department, delta_constant_pence, delta_variable_pence, consumed_surplus_pence, cycle_number, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := m.db.ExecContext(ctx, q,
		string(r.ID), string(schemeID), string(r.Department),
		r.DeltaConstantPence, r.DeltaVariablePence, r.ConsumedSurplusPence, r.CycleNumber,
		m.now(),
	); err != nil {
		if isDuplicateKey(err) {
			return repro.Reinvestment{}, ErrAlreadyExists
		}
		return repro.Reinvestment{}, err
	}
	return r, nil
}

func (m *MySQL) ListReinvestments(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID) ([]repro.Reinvestment, error) {
	const q = `SELECT id, scheme_id, department, delta_constant_pence, delta_variable_pence,
		consumed_surplus_pence, cycle_number
		FROM reinvestments WHERE scheme_id = ? ORDER BY cycle_number ASC, created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(schemeID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]repro.Reinvestment, 0)
	for rows.Next() {
		var r repro.Reinvestment
		var rid, sid, dept string
		if err := rows.Scan(&rid, &sid, &dept,
			&r.DeltaConstantPence, &r.DeltaVariablePence, &r.ConsumedSurplusPence, &r.CycleNumber); err != nil {
			return nil, err
		}
		r.ID = repro.ReinvestmentID(rid)
		r.SchemeID = repro.ExtendedReproductionSchemeID(sid)
		r.Department = repro.CapitalDepartment(dept)
		out = append(out, r)
	}
	return out, nil
}

// --- TickExtendedScheme -----------------------------------------------------

func (m *MySQL) TickExtendedScheme(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID) (repro.ExtendedReproductionScheme, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}
	defer func() { _ = tx.Rollback() }()

	// Lock and read the scheme row.
	var sid, period string
	var rateI, rateII, tickCount int64
	if err := tx.QueryRowContext(ctx,
		`SELECT id, period, accumulation_rate_i_bps, accumulation_rate_ii_bps, tick_count
		 FROM extended_reproduction_schemes WHERE id = ? FOR UPDATE`,
		string(schemeID),
	).Scan(&sid, &period, &rateI, &rateII, &tickCount); err != nil {
		if err == sql.ErrNoRows {
			return repro.ExtendedReproductionScheme{}, ErrNotFound
		}
		return repro.ExtendedReproductionScheme{}, err
	}

	// Lock and read both department rows.
	deptI, err := m.getExtendedDeptForUpdate(ctx, tx, schemeID, repro.DepartmentI)
	if err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}
	deptII, err := m.getExtendedDeptForUpdate(ctx, tx, schemeID, repro.DepartmentII)
	if err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}

	cycleNumber := tickCount + 1

	// Compute reinvestments.
	rI, err := repro.ComputeReinvestment(deptI, rateI)
	if err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}
	rI.SchemeID = schemeID
	rI.CycleNumber = cycleNumber

	rII, err := repro.ComputeReinvestment(deptII, rateII)
	if err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}
	rII.SchemeID = schemeID
	rII.CycleNumber = cycleNumber

	// Insert reinvestment rows.
	const insertReinvest = `INSERT INTO reinvestments
		(id, scheme_id, department, delta_constant_pence, delta_variable_pence, consumed_surplus_pence, cycle_number, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := tx.ExecContext(ctx, insertReinvest,
		string(rI.ID), string(schemeID), string(rI.Department),
		rI.DeltaConstantPence, rI.DeltaVariablePence, rI.ConsumedSurplusPence, rI.CycleNumber, m.now(),
	); err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}
	if _, err := tx.ExecContext(ctx, insertReinvest,
		string(rII.ID), string(schemeID), string(rII.Department),
		rII.DeltaConstantPence, rII.DeltaVariablePence, rII.ConsumedSurplusPence, rII.CycleNumber, m.now(),
	); err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}

	// Advance departmental capitals.
	newDeptI := repro.ComputeNextPeriodCapital(deptI, rI)
	newDeptII := repro.ComputeNextPeriodCapital(deptII, rII)

	const updateDept = `UPDATE extended_departmental_capitals
		SET constant_pence = ?, variable_pence = ?, surplus_pence = ?, total_pence = ?
		WHERE scheme_id = ? AND department = ?`
	if _, err := tx.ExecContext(ctx, updateDept,
		newDeptI.ConstantPence, newDeptI.VariablePence, newDeptI.SurplusPence, newDeptI.TotalPence,
		string(schemeID), string(repro.DepartmentI),
	); err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}
	if _, err := tx.ExecContext(ctx, updateDept,
		newDeptII.ConstantPence, newDeptII.VariablePence, newDeptII.SurplusPence, newDeptII.TotalPence,
		string(schemeID), string(repro.DepartmentII),
	); err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}

	// Update scheme.
	isBalanced := repro.ComputeExtendedBalance(newDeptI, newDeptII, rI, rII)
	if _, err := tx.ExecContext(ctx,
		`UPDATE extended_reproduction_schemes SET is_balanced = ?, tick_count = ? WHERE id = ?`,
		isBalanced, cycleNumber, string(schemeID),
	); err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}

	if err := tx.Commit(); err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}

	return repro.ExtendedReproductionScheme{
		ID:                    repro.ExtendedReproductionSchemeID(sid),
		Period:                period,
		DepartmentI:           &newDeptI,
		DepartmentII:          &newDeptII,
		AccumulationRateIBps:  rateI,
		AccumulationRateIIBps: rateII,
		IsBalanced:            isBalanced,
		TickCount:             cycleNumber,
	}, nil
}

// getExtendedDeptForUpdate reads a departmental capital row with SELECT FOR UPDATE.
func (m *MySQL) getExtendedDeptForUpdate(ctx context.Context, tx *sql.Tx, schemeID repro.ExtendedReproductionSchemeID, dept repro.CapitalDepartment) (repro.DepartmentalCapital, error) {
	const q = `SELECT id, scheme_id, department, constant_pence, variable_pence,
		surplus_pence, total_pence, turnover_number_bp
		FROM extended_departmental_capitals WHERE scheme_id = ? AND department = ? FOR UPDATE`
	row := tx.QueryRowContext(ctx, q, string(schemeID), string(dept))
	var d repro.DepartmentalCapital
	var did, sid, department string
	if err := row.Scan(&did, &sid, &department,
		&d.ConstantPence, &d.VariablePence, &d.SurplusPence, &d.TotalPence,
		&d.TurnoverNumberBP); err != nil {
		if err == sql.ErrNoRows {
			return repro.DepartmentalCapital{}, ErrNotFound
		}
		return repro.DepartmentalCapital{}, err
	}
	d.ID = repro.DepartmentalCapitalID(did)
	d.Department = repro.CapitalDepartment(department)
	return d, nil
}

// --- GetExtendedMoneyRequirement --------------------------------------------

func (m *MySQL) GetExtendedMoneyRequirement(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID, baseMoneyPence int64) (repro.ExtendedMoneyRequirement, error) {
	s, err := m.GetExtendedReproductionScheme(ctx, schemeID)
	if err != nil {
		return repro.ExtendedMoneyRequirement{}, err
	}
	return repro.ComputeExtendedMoneyRequirement(s, baseMoneyPence)
}

// --- MultiPeriodScheme ------------------------------------------------------

func (m *MySQL) CreateMultiPeriodScheme(ctx context.Context, mp repro.MultiPeriodScheme) (repro.MultiPeriodScheme, error) {
	if mp.ID == "" {
		mp.ID = repro.NewMultiPeriodSchemeID()
	}

	// Collect all schemes first to validate existence.
	schemes := make([]repro.ExtendedReproductionScheme, len(mp.SchemeIDs))
	for i, sid := range mp.SchemeIDs {
		s, err := m.GetExtendedReproductionScheme(ctx, sid)
		if err != nil {
			return repro.MultiPeriodScheme{}, err
		}
		schemes[i] = s
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return repro.MultiPeriodScheme{}, err
	}
	defer func() { _ = tx.Rollback() }()

	// Insert multi_period_schemes header.
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO multi_period_schemes (id, label, created_at) VALUES (?, ?, ?)`,
		string(mp.ID), mp.Label, m.now(),
	); err != nil {
		if isDuplicateKey(err) {
			return repro.MultiPeriodScheme{}, ErrAlreadyExists
		}
		return repro.MultiPeriodScheme{}, err
	}

	// Insert multi_period_scheme_entries.
	for i, sid := range mp.SchemeIDs {
		eid := repro.NewExtendedReproductionSchemeID()
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO multi_period_scheme_entries (id, multi_period_id, scheme_id, position) VALUES (?, ?, ?, ?)`,
			string(eid), string(mp.ID), string(sid), int64(i),
		); err != nil {
			return repro.MultiPeriodScheme{}, err
		}
	}

	// Compute CompositionShift and DepartmentIGrowthLead for consecutive pairs.
	for i := 1; i < len(schemes); i++ {
		prev := schemes[i-1]
		curr := schemes[i]
		cycle := int64(i)

		var deptIComp, deptIIComp int64
		if curr.DepartmentI != nil {
			deptIComp = repro.ComputeCompositionBps(*curr.DepartmentI)
		}
		if curr.DepartmentII != nil {
			deptIIComp = repro.ComputeCompositionBps(*curr.DepartmentII)
		}
		socialComp := (deptIComp + deptIIComp) / 2

		csID := repro.NewCompositionShiftID()
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO composition_shifts
			(id, multi_period_id, cycle_number, dept_i_composition_bps, dept_ii_composition_bps, social_composition_bps, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			string(csID), string(mp.ID), cycle, deptIComp, deptIIComp, socialComp, m.now(),
		); err != nil {
			return repro.MultiPeriodScheme{}, err
		}

		var deptIGrowth, deptIIGrowth int64
		if prev.DepartmentI != nil && curr.DepartmentI != nil {
			deptIGrowth = repro.ComputeGrowthBps(*prev.DepartmentI, *curr.DepartmentI)
		}
		if prev.DepartmentII != nil && curr.DepartmentII != nil {
			deptIIGrowth = repro.ComputeGrowthBps(*prev.DepartmentII, *curr.DepartmentII)
		}
		leadConfirmed := deptIGrowth > deptIIGrowth

		glID := repro.NewDepartmentIGrowthLeadID()
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO dept_i_growth_leads
			(id, multi_period_id, cycle_number, dept_i_growth_bps, dept_ii_growth_bps, lead_confirmed, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			string(glID), string(mp.ID), cycle, deptIGrowth, deptIIGrowth, leadConfirmed, m.now(),
		); err != nil {
			return repro.MultiPeriodScheme{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return repro.MultiPeriodScheme{}, err
	}
	return mp, nil
}

func (m *MySQL) GetMultiPeriodScheme(ctx context.Context, id repro.MultiPeriodSchemeID) (repro.MultiPeriodScheme, error) {
	var label string
	if err := m.db.QueryRowContext(ctx,
		`SELECT label FROM multi_period_schemes WHERE id = ?`, string(id),
	).Scan(&label); err != nil {
		if err == sql.ErrNoRows {
			return repro.MultiPeriodScheme{}, ErrNotFound
		}
		return repro.MultiPeriodScheme{}, err
	}

	rows, err := m.db.QueryContext(ctx,
		`SELECT scheme_id FROM multi_period_scheme_entries WHERE multi_period_id = ? ORDER BY position ASC`,
		string(id),
	)
	if err != nil {
		return repro.MultiPeriodScheme{}, err
	}
	defer rows.Close()
	var schemeIDs []repro.ExtendedReproductionSchemeID
	for rows.Next() {
		var sid string
		if err := rows.Scan(&sid); err != nil {
			return repro.MultiPeriodScheme{}, err
		}
		schemeIDs = append(schemeIDs, repro.ExtendedReproductionSchemeID(sid))
	}
	return repro.MultiPeriodScheme{
		ID:        id,
		Label:     label,
		SchemeIDs: schemeIDs,
	}, nil
}

func (m *MySQL) ListCompositionShifts(ctx context.Context, mpID repro.MultiPeriodSchemeID) ([]repro.CompositionShift, error) {
	const q = `SELECT id, multi_period_id, cycle_number, dept_i_composition_bps,
		dept_ii_composition_bps, social_composition_bps
		FROM composition_shifts WHERE multi_period_id = ? ORDER BY cycle_number ASC`
	rows, err := m.db.QueryContext(ctx, q, string(mpID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]repro.CompositionShift, 0)
	for rows.Next() {
		var cs repro.CompositionShift
		var csID, mpid string
		if err := rows.Scan(&csID, &mpid, &cs.CycleNumber,
			&cs.DeptICompositionBps, &cs.DeptIICompositionBps, &cs.SocialCompositionBps); err != nil {
			return nil, err
		}
		cs.ID = repro.CompositionShiftID(csID)
		cs.MultiPeriodID = repro.MultiPeriodSchemeID(mpid)
		out = append(out, cs)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CycleNumber < out[j].CycleNumber })
	return out, nil
}

func (m *MySQL) ListGrowthLeads(ctx context.Context, mpID repro.MultiPeriodSchemeID) ([]repro.DepartmentIGrowthLead, error) {
	const q = `SELECT id, multi_period_id, cycle_number, dept_i_growth_bps, dept_ii_growth_bps, lead_confirmed
		FROM dept_i_growth_leads WHERE multi_period_id = ? ORDER BY cycle_number ASC`
	rows, err := m.db.QueryContext(ctx, q, string(mpID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]repro.DepartmentIGrowthLead, 0)
	for rows.Next() {
		var gl repro.DepartmentIGrowthLead
		var glID, mpid string
		if err := rows.Scan(&glID, &mpid, &gl.CycleNumber,
			&gl.DeptIGrowthBps, &gl.DeptIIGrowthBps, &gl.LeadConfirmed); err != nil {
			return nil, err
		}
		gl.ID = repro.DepartmentIGrowthLeadID(glID)
		gl.MultiPeriodID = repro.MultiPeriodSchemeID(mpid)
		out = append(out, gl)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CycleNumber < out[j].CycleNumber })
	return out, nil
}
