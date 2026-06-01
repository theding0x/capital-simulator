package store

import (
	"context"
	"database/sql"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// Vol. II Ch. 20 — Simple Reproduction.

// --- SimpleReproductionScheme -----------------------------------------------

func (m *MySQL) CreateSimpleReproductionScheme(ctx context.Context, s repro.SimpleReproductionScheme) (repro.SimpleReproductionScheme, error) {
	if s.ID == "" {
		s.ID = repro.NewSimpleReproductionSchemeID()
	}
	s.IsBalanced = false
	s.TickCount = 0
	const q = `INSERT INTO simple_reproduction_schemes
		(id, period, is_balanced, tick_count, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID), s.Period, s.IsBalanced, s.TickCount, m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return repro.SimpleReproductionScheme{}, ErrAlreadyExists
		}
		return repro.SimpleReproductionScheme{}, err
	}
	return s, nil
}

func (m *MySQL) GetSimpleReproductionScheme(ctx context.Context, id repro.SimpleReproductionSchemeID) (repro.SimpleReproductionScheme, error) {
	const q = `SELECT id, period, tick_count FROM simple_reproduction_schemes WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var s repro.SimpleReproductionScheme
	var sid string
	if err := row.Scan(&sid, &s.Period, &s.TickCount); err != nil {
		if err == sql.ErrNoRows {
			return repro.SimpleReproductionScheme{}, ErrNotFound
		}
		return repro.SimpleReproductionScheme{}, err
	}
	s.ID = repro.SimpleReproductionSchemeID(sid)

	// Load departments.
	deptI, err := m.getDepartmentalCapital(ctx, id, repro.DepartmentI)
	if err == nil {
		s.DepartmentI = &deptI
	} else if err != ErrNotFound {
		return repro.SimpleReproductionScheme{}, err
	}
	deptII, err := m.getDepartmentalCapital(ctx, id, repro.DepartmentII)
	if err == nil {
		s.DepartmentII = &deptII
	} else if err != ErrNotFound {
		return repro.SimpleReproductionScheme{}, err
	}
	s.IsBalanced = s.CheckBalance()
	return s, nil
}

func (m *MySQL) ListSimpleReproductionSchemes(ctx context.Context, period string) ([]repro.SimpleReproductionScheme, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if period != "" {
		const q = `SELECT id FROM simple_reproduction_schemes WHERE period = ? ORDER BY created_at ASC`
		rows, err = m.db.QueryContext(ctx, q, period)
	} else {
		const q = `SELECT id FROM simple_reproduction_schemes ORDER BY created_at ASC`
		rows, err = m.db.QueryContext(ctx, q)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]repro.SimpleReproductionScheme, 0)
	for rows.Next() {
		var sid string
		if err := rows.Scan(&sid); err != nil {
			return nil, err
		}
		s, err := m.GetSimpleReproductionScheme(ctx, repro.SimpleReproductionSchemeID(sid))
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

// --- DepartmentalCapital ----------------------------------------------------

func (m *MySQL) AddDepartmentToScheme(ctx context.Context, id repro.SimpleReproductionSchemeID, d repro.DepartmentalCapital) (repro.DepartmentalCapital, error) {
	if err := d.Validate(); err != nil {
		return repro.DepartmentalCapital{}, err
	}
	if d.ID == "" {
		d.ID = repro.NewDepartmentalCapitalID()
	}
	d.SchemeID = id

	// Check that the scheme exists.
	var dummy string
	if err := m.db.QueryRowContext(ctx,
		`SELECT id FROM simple_reproduction_schemes WHERE id = ?`, string(id),
	).Scan(&dummy); err != nil {
		if err == sql.ErrNoRows {
			return repro.DepartmentalCapital{}, ErrNotFound
		}
		return repro.DepartmentalCapital{}, err
	}

	// Check for an existing department row.
	if _, err := m.getDepartmentalCapital(ctx, id, d.Department); err == nil {
		return repro.DepartmentalCapital{}, repro.ErrDepartmentAlreadySet
	} else if err != ErrNotFound {
		return repro.DepartmentalCapital{}, err
	}

	const q = `INSERT INTO departmental_capitals
		(id, scheme_id, department, constant_pence, variable_pence,
		 surplus_pence, total_pence, turnover_number_bp, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := m.db.ExecContext(ctx, q,
		string(d.ID), string(id), string(d.Department),
		d.ConstantPence, d.VariablePence, d.SurplusPence, d.TotalPence,
		d.TurnoverNumberBP, m.now(),
	); err != nil {
		if isDuplicateKey(err) {
			return repro.DepartmentalCapital{}, ErrAlreadyExists
		}
		return repro.DepartmentalCapital{}, err
	}

	// Recompute IsBalanced on the scheme.
	s, err := m.GetSimpleReproductionScheme(ctx, id)
	if err == nil {
		_, _ = m.db.ExecContext(ctx,
			`UPDATE simple_reproduction_schemes SET is_balanced = ? WHERE id = ?`,
			s.IsBalanced, string(id),
		)
	}
	return d, nil
}

// getDepartmentalCapital loads a single DepartmentalCapital for a scheme + department.
func (m *MySQL) getDepartmentalCapital(ctx context.Context, schemeID repro.SimpleReproductionSchemeID, dept repro.CapitalDepartment) (repro.DepartmentalCapital, error) {
	const q = `SELECT id, scheme_id, department, constant_pence, variable_pence,
		surplus_pence, total_pence, turnover_number_bp
		FROM departmental_capitals WHERE scheme_id = ? AND department = ?`
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
	d.SchemeID = repro.SimpleReproductionSchemeID(sid)
	d.Department = repro.CapitalDepartment(department)
	return d, nil
}

// --- InterDepartmentExchange ------------------------------------------------

func (m *MySQL) RecordInterDepartmentExchange(ctx context.Context, id repro.SimpleReproductionSchemeID, e repro.InterDepartmentExchange) (repro.InterDepartmentExchange, error) {
	if err := e.Validate(); err != nil {
		return repro.InterDepartmentExchange{}, err
	}
	if e.ID == "" {
		e.ID = repro.NewInterDepartmentExchangeID()
	}
	e.SchemeID = id
	const q = `INSERT INTO interdepartment_exchanges
		(id, scheme_id, from_department, to_department, pence, kind, description, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(e.ID), string(id),
		string(e.FromDepartment), string(e.ToDepartment),
		e.Pence, string(e.Kind), e.Description,
		m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return repro.InterDepartmentExchange{}, ErrAlreadyExists
		}
		return repro.InterDepartmentExchange{}, err
	}
	return e, nil
}

func (m *MySQL) ListInterDepartmentExchanges(ctx context.Context, schemeID repro.SimpleReproductionSchemeID) ([]repro.InterDepartmentExchange, error) {
	const q = `SELECT id, scheme_id, from_department, to_department, pence, kind, description
		FROM interdepartment_exchanges WHERE scheme_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(schemeID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]repro.InterDepartmentExchange, 0)
	for rows.Next() {
		var e repro.InterDepartmentExchange
		var eid, sid, from, to, kind string
		if err := rows.Scan(&eid, &sid, &from, &to, &e.Pence, &kind, &e.Description); err != nil {
			return nil, err
		}
		e.ID = repro.InterDepartmentExchangeID(eid)
		e.SchemeID = repro.SimpleReproductionSchemeID(sid)
		e.FromDepartment = repro.CapitalDepartment(from)
		e.ToDepartment = repro.CapitalDepartment(to)
		e.Kind = repro.ExchangeKind(kind)
		out = append(out, e)
	}
	return out, nil
}

// --- ReproductionTick -------------------------------------------------------

func (m *MySQL) AdvanceReproductionTick(ctx context.Context, id repro.SimpleReproductionSchemeID) (repro.ReproductionTick, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return repro.ReproductionTick{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var period string
	var tickCount int64
	if err := tx.QueryRowContext(ctx,
		`SELECT period, tick_count FROM simple_reproduction_schemes WHERE id = ? FOR UPDATE`,
		string(id),
	).Scan(&period, &tickCount); err != nil {
		if err == sql.ErrNoRows {
			return repro.ReproductionTick{}, ErrNotFound
		}
		return repro.ReproductionTick{}, err
	}
	tickCount++
	if _, err := tx.ExecContext(ctx,
		`UPDATE simple_reproduction_schemes SET tick_count = ? WHERE id = ?`,
		tickCount, string(id),
	); err != nil {
		return repro.ReproductionTick{}, err
	}

	s, balErr := m.GetSimpleReproductionScheme(ctx, id)
	isBalanced := balErr == nil && s.IsBalanced

	tick := repro.ReproductionTick{
		ID:         repro.NewReproductionTickID(),
		SchemeID:   id,
		TickNumber: tickCount,
		Period:     period,
		IsBalanced: isBalanced,
	}
	const q = `INSERT INTO reproduction_ticks
		(id, scheme_id, tick_number, period, is_balanced, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	if _, err := tx.ExecContext(ctx, q,
		string(tick.ID), string(id), tick.TickNumber, tick.Period, tick.IsBalanced, m.now(),
	); err != nil {
		return repro.ReproductionTick{}, err
	}

	if err := tx.Commit(); err != nil {
		return repro.ReproductionTick{}, err
	}
	return tick, nil
}

// --- BalanceCheck -----------------------------------------------------------

func (m *MySQL) CheckSchemeBalance(ctx context.Context, id repro.SimpleReproductionSchemeID) (repro.BalanceCheckResult, error) {
	s, err := m.GetSimpleReproductionScheme(ctx, id)
	if err != nil {
		return repro.BalanceCheckResult{}, err
	}
	res := repro.BalanceCheckResult{
		SchemeID:   string(id),
		IsBalanced: s.IsBalanced,
	}
	if s.DepartmentI != nil {
		res.DeptIVplusSPence = s.DepartmentI.VplusSPence()
	}
	if s.DepartmentII != nil {
		res.DeptIIConstantPence = s.DepartmentII.ConstantPence
	}
	res.DeficitPence = res.DeptIVplusSPence - res.DeptIIConstantPence
	repro.ApplyAnnualisedSettlement(&res, s.DepartmentI, s.DepartmentII)
	return res, nil
}

// --- MoneyClosedLoop --------------------------------------------------------

func (m *MySQL) RecordMoneyClosedLoop(ctx context.Context, id repro.SimpleReproductionSchemeID, loop repro.MoneyClosedLoop) (repro.MoneyClosedLoop, error) {
	if loop.ID == "" {
		loop.ID = repro.NewMoneyClosedLoopID()
	}
	loop.SchemeID = id
	const q = `INSERT INTO money_closed_loops
		(id, scheme_id, period, is_closed, net_flow_pence, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(loop.ID), string(id), loop.Period, loop.IsClosed, loop.NetFlowPence, m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return repro.MoneyClosedLoop{}, ErrAlreadyExists
		}
		return repro.MoneyClosedLoop{}, err
	}
	return loop, nil
}
