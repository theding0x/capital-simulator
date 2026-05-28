package store

import (
	"context"
	"database/sql"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// Vol. II Ch. 18 — The Role of Money-Capital in Reproduction.

// --- MoneySupplyApportionment -----------------------------------------------

func (m *MySQL) CreateMoneySupplyApportionment(ctx context.Context, a repro.MoneySupplyApportionment) (repro.MoneySupplyApportionment, error) {
	if err := a.Validate(); err != nil {
		return repro.MoneySupplyApportionment{}, err
	}
	if a.ID == "" {
		a.ID = repro.NewMoneySupplyApportionmentID()
	}
	const q = `INSERT INTO money_supply_apportionments
		(id, total_social_money_pence, dept_i_reserve_pence,
		 dept_ii_reserve_pence, wage_rotation_pence, idle_hoard_pence,
		 period, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID),
		a.TotalSocialMoneyPence,
		a.DeptIReservePence,
		a.DeptIIReservePence,
		a.WageRotationPence,
		a.IdleHoardPence,
		a.Period,
		m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return repro.MoneySupplyApportionment{}, ErrAlreadyExists
		}
		return repro.MoneySupplyApportionment{}, err
	}
	return a, nil
}

func (m *MySQL) GetMoneySupplyApportionment(ctx context.Context, id repro.MoneySupplyApportionmentID) (repro.MoneySupplyApportionment, error) {
	const q = `SELECT id, total_social_money_pence, dept_i_reserve_pence,
		dept_ii_reserve_pence, wage_rotation_pence, idle_hoard_pence, period
		FROM money_supply_apportionments WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanMoneySupplyApportionment(row.Scan)
}

func (m *MySQL) ListMoneySupplyApportionments(ctx context.Context) ([]repro.MoneySupplyApportionment, error) {
	const q = `SELECT id, total_social_money_pence, dept_i_reserve_pence,
		dept_ii_reserve_pence, wage_rotation_pence, idle_hoard_pence, period
		FROM money_supply_apportionments ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]repro.MoneySupplyApportionment, 0)
	for rows.Next() {
		a, err := scanMoneySupplyApportionment(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func (m *MySQL) CheckApportionmentBalance(ctx context.Context, id repro.MoneySupplyApportionmentID) (bool, error) {
	a, err := m.GetMoneySupplyApportionment(ctx, id)
	if err != nil {
		return false, err
	}
	parts := a.DeptIReservePence +
		a.DeptIIReservePence +
		a.WageRotationPence +
		a.IdleHoardPence
	return parts == a.TotalSocialMoneyPence, nil
}

// --- DepartmentMoneyReserve -------------------------------------------------

func (m *MySQL) CreateDepartmentMoneyReserve(ctx context.Context, r repro.DepartmentMoneyReserve) (repro.DepartmentMoneyReserve, error) {
	if err := r.Department.Validate(); err != nil {
		return repro.DepartmentMoneyReserve{}, err
	}
	if err := r.ReservePurpose.Validate(); err != nil {
		return repro.DepartmentMoneyReserve{}, err
	}
	if r.ID == "" {
		r.ID = repro.NewDepartmentMoneyReserveID()
	}
	const q = `INSERT INTO department_money_reserves
		(id, department, reserve_pence, reserve_purpose, period, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(r.ID),
		string(r.Department),
		r.ReservePence,
		string(r.ReservePurpose),
		r.Period,
		m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return repro.DepartmentMoneyReserve{}, ErrAlreadyExists
		}
		return repro.DepartmentMoneyReserve{}, err
	}
	return r, nil
}

func (m *MySQL) GetDepartmentMoneyReserve(ctx context.Context, id repro.DepartmentMoneyReserveID) (repro.DepartmentMoneyReserve, error) {
	const q = `SELECT id, department, reserve_pence, reserve_purpose, period
		FROM department_money_reserves WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanDepartmentMoneyReserve(row.Scan)
}

func (m *MySQL) ListDepartmentMoneyReserves(ctx context.Context) ([]repro.DepartmentMoneyReserve, error) {
	const q = `SELECT id, department, reserve_pence, reserve_purpose, period
		FROM department_money_reserves ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]repro.DepartmentMoneyReserve, 0)
	for rows.Next() {
		r, err := scanDepartmentMoneyReserve(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

// --- CirculatingMoneyMass ---------------------------------------------------

func (m *MySQL) CreateCirculatingMoneyMass(ctx context.Context, mass repro.CirculatingMoneyMass) (repro.CirculatingMoneyMass, error) {
	if mass.MoneyStockPence < 0 {
		return repro.CirculatingMoneyMass{}, repro.ErrNegativeMoneyStock
	}
	if mass.VelocityPerYearBasisPoints <= 0 {
		return repro.CirculatingMoneyMass{}, repro.ErrZeroVelocity
	}
	if mass.ID == "" {
		mass.ID = repro.NewCirculatingMoneyMassID()
	}
	const q = `INSERT INTO circulating_money_masses
		(id, money_stock_pence, velocity_per_year_basis_points,
		 effective_circulating_value_pence, period, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(mass.ID),
		mass.MoneyStockPence,
		mass.VelocityPerYearBasisPoints,
		mass.EffectiveCirculatingValuePence,
		mass.Period,
		m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return repro.CirculatingMoneyMass{}, ErrAlreadyExists
		}
		return repro.CirculatingMoneyMass{}, err
	}
	return mass, nil
}

func (m *MySQL) GetCirculatingMoneyMass(ctx context.Context, id repro.CirculatingMoneyMassID) (repro.CirculatingMoneyMass, error) {
	const q = `SELECT id, money_stock_pence, velocity_per_year_basis_points,
		effective_circulating_value_pence, period
		FROM circulating_money_masses WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanCirculatingMoneyMass(row.Scan)
}

func (m *MySQL) ListCirculatingMoneyMasses(ctx context.Context) ([]repro.CirculatingMoneyMass, error) {
	const q = `SELECT id, money_stock_pence, velocity_per_year_basis_points,
		effective_circulating_value_pence, period
		FROM circulating_money_masses ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]repro.CirculatingMoneyMass, 0)
	for rows.Next() {
		mass, err := scanCirculatingMoneyMass(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, mass)
	}
	return out, nil
}

// --- WageRotationFund -------------------------------------------------------

func (m *MySQL) CreateWageRotationFund(ctx context.Context, w repro.WageRotationFund) (repro.WageRotationFund, error) {
	if w.ID == "" {
		w.ID = repro.NewWageRotationFundID()
	}
	const q = `INSERT INTO wage_rotation_funds
		(id, fund_pence, wage_cycle_frequency, department, period, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(w.ID),
		w.FundPence,
		w.WageCycleFrequency,
		string(w.Department),
		w.Period,
		m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return repro.WageRotationFund{}, ErrAlreadyExists
		}
		return repro.WageRotationFund{}, err
	}
	return w, nil
}

func (m *MySQL) GetWageRotationFund(ctx context.Context, id repro.WageRotationFundID) (repro.WageRotationFund, error) {
	const q = `SELECT id, fund_pence, wage_cycle_frequency, department, period
		FROM wage_rotation_funds WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanWageRotationFund(row.Scan)
}

func (m *MySQL) ListWageRotationFunds(ctx context.Context) ([]repro.WageRotationFund, error) {
	const q = `SELECT id, fund_pence, wage_cycle_frequency, department, period
		FROM wage_rotation_funds ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]repro.WageRotationFund, 0)
	for rows.Next() {
		w, err := scanWageRotationFund(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, nil
}

// --- InterDepartmentSettlement ----------------------------------------------

func (m *MySQL) CreateInterDepartmentSettlement(ctx context.Context, s repro.InterDepartmentSettlement) (repro.InterDepartmentSettlement, error) {
	if err := s.FromDepartment.Validate(); err != nil {
		return repro.InterDepartmentSettlement{}, err
	}
	if err := s.ToDepartment.Validate(); err != nil {
		return repro.InterDepartmentSettlement{}, err
	}
	if err := s.SettlementPurpose.Validate(); err != nil {
		return repro.InterDepartmentSettlement{}, err
	}
	if s.ID == "" {
		s.ID = repro.NewInterDepartmentSettlementID()
	}
	const q = `INSERT INTO inter_department_settlements
		(id, from_department, to_department, settled_pence, settlement_purpose, period, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID),
		string(s.FromDepartment),
		string(s.ToDepartment),
		s.SettledPence,
		string(s.SettlementPurpose),
		s.Period,
		m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return repro.InterDepartmentSettlement{}, ErrAlreadyExists
		}
		return repro.InterDepartmentSettlement{}, err
	}
	return s, nil
}

func (m *MySQL) GetInterDepartmentSettlement(ctx context.Context, id repro.InterDepartmentSettlementID) (repro.InterDepartmentSettlement, error) {
	const q = `SELECT id, from_department, to_department, settled_pence, settlement_purpose, period
		FROM inter_department_settlements WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanInterDepartmentSettlement(row.Scan)
}

func (m *MySQL) ListInterDepartmentSettlements(ctx context.Context) ([]repro.InterDepartmentSettlement, error) {
	const q = `SELECT id, from_department, to_department, settled_pence, settlement_purpose, period
		FROM inter_department_settlements ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]repro.InterDepartmentSettlement, 0)
	for rows.Next() {
		s, err := scanInterDepartmentSettlement(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

// --- scan helpers -----------------------------------------------------------

func scanMoneySupplyApportionment(scan func(...any) error) (repro.MoneySupplyApportionment, error) {
	var a repro.MoneySupplyApportionment
	var id string
	if err := scan(&id,
		&a.TotalSocialMoneyPence,
		&a.DeptIReservePence,
		&a.DeptIIReservePence,
		&a.WageRotationPence,
		&a.IdleHoardPence,
		&a.Period,
	); err != nil {
		if err == sql.ErrNoRows {
			return repro.MoneySupplyApportionment{}, ErrNotFound
		}
		return repro.MoneySupplyApportionment{}, err
	}
	a.ID = repro.MoneySupplyApportionmentID(id)
	return a, nil
}

func scanDepartmentMoneyReserve(scan func(...any) error) (repro.DepartmentMoneyReserve, error) {
	var r repro.DepartmentMoneyReserve
	var id, department, reservePurpose string
	if err := scan(&id, &department, &r.ReservePence, &reservePurpose, &r.Period); err != nil {
		if err == sql.ErrNoRows {
			return repro.DepartmentMoneyReserve{}, ErrNotFound
		}
		return repro.DepartmentMoneyReserve{}, err
	}
	r.ID = repro.DepartmentMoneyReserveID(id)
	r.Department = repro.CapitalDepartment(department)
	r.ReservePurpose = repro.ReservePurpose(reservePurpose)
	return r, nil
}

func scanCirculatingMoneyMass(scan func(...any) error) (repro.CirculatingMoneyMass, error) {
	var mass repro.CirculatingMoneyMass
	var id string
	if err := scan(&id,
		&mass.MoneyStockPence,
		&mass.VelocityPerYearBasisPoints,
		&mass.EffectiveCirculatingValuePence,
		&mass.Period,
	); err != nil {
		if err == sql.ErrNoRows {
			return repro.CirculatingMoneyMass{}, ErrNotFound
		}
		return repro.CirculatingMoneyMass{}, err
	}
	mass.ID = repro.CirculatingMoneyMassID(id)
	return mass, nil
}

func scanWageRotationFund(scan func(...any) error) (repro.WageRotationFund, error) {
	var w repro.WageRotationFund
	var id, department string
	if err := scan(&id, &w.FundPence, &w.WageCycleFrequency, &department, &w.Period); err != nil {
		if err == sql.ErrNoRows {
			return repro.WageRotationFund{}, ErrNotFound
		}
		return repro.WageRotationFund{}, err
	}
	w.ID = repro.WageRotationFundID(id)
	w.Department = repro.CapitalDepartment(department)
	return w, nil
}

func scanInterDepartmentSettlement(scan func(...any) error) (repro.InterDepartmentSettlement, error) {
	var s repro.InterDepartmentSettlement
	var id, fromDept, toDept, settlementPurpose string
	if err := scan(&id, &fromDept, &toDept, &s.SettledPence, &settlementPurpose, &s.Period); err != nil {
		if err == sql.ErrNoRows {
			return repro.InterDepartmentSettlement{}, ErrNotFound
		}
		return repro.InterDepartmentSettlement{}, err
	}
	s.ID = repro.InterDepartmentSettlementID(id)
	s.FromDepartment = repro.CapitalDepartment(fromDept)
	s.ToDepartment = repro.CapitalDepartment(toDept)
	s.SettlementPurpose = repro.ReservePurpose(settlementPurpose)
	return s, nil
}
