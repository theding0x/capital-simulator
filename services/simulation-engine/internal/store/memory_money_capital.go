package store

import (
	"context"
	"sync"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// memoryMoneyCapital is the in-memory backing store for Vol. II Ch. 18 types.
type memoryMoneyCapital struct {
	mu                   sync.RWMutex
	apportionments       map[repro.MoneySupplyApportionmentID]repro.MoneySupplyApportionment
	departmentReserves   map[repro.DepartmentMoneyReserveID]repro.DepartmentMoneyReserve
	circulatingMasses    map[repro.CirculatingMoneyMassID]repro.CirculatingMoneyMass
	wageRotationFunds    map[repro.WageRotationFundID]repro.WageRotationFund
	interDeptSettlements map[repro.InterDepartmentSettlementID]repro.InterDepartmentSettlement
}

func newMemoryMoneyCapital() *memoryMoneyCapital {
	return &memoryMoneyCapital{
		apportionments:       make(map[repro.MoneySupplyApportionmentID]repro.MoneySupplyApportionment),
		departmentReserves:   make(map[repro.DepartmentMoneyReserveID]repro.DepartmentMoneyReserve),
		circulatingMasses:    make(map[repro.CirculatingMoneyMassID]repro.CirculatingMoneyMass),
		wageRotationFunds:    make(map[repro.WageRotationFundID]repro.WageRotationFund),
		interDeptSettlements: make(map[repro.InterDepartmentSettlementID]repro.InterDepartmentSettlement),
	}
}

// --- MoneySupplyApportionment -----------------------------------------------

func (m *Memory) CreateMoneySupplyApportionment(_ context.Context, a repro.MoneySupplyApportionment) (repro.MoneySupplyApportionment, error) {
	if err := a.Validate(); err != nil {
		return repro.MoneySupplyApportionment{}, err
	}
	if a.ID == "" {
		a.ID = repro.NewMoneySupplyApportionmentID()
	}
	m.moneyCapital.mu.Lock()
	defer m.moneyCapital.mu.Unlock()
	if _, exists := m.moneyCapital.apportionments[a.ID]; exists {
		return repro.MoneySupplyApportionment{}, ErrAlreadyExists
	}
	m.moneyCapital.apportionments[a.ID] = a
	return a, nil
}

func (m *Memory) GetMoneySupplyApportionment(_ context.Context, id repro.MoneySupplyApportionmentID) (repro.MoneySupplyApportionment, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	a, ok := m.moneyCapital.apportionments[id]
	if !ok {
		return repro.MoneySupplyApportionment{}, ErrNotFound
	}
	return a, nil
}

func (m *Memory) ListMoneySupplyApportionments(_ context.Context) ([]repro.MoneySupplyApportionment, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	out := make([]repro.MoneySupplyApportionment, 0, len(m.moneyCapital.apportionments))
	for _, a := range m.moneyCapital.apportionments {
		out = append(out, a)
	}
	return out, nil
}

// CheckApportionmentBalance returns true when the four parts of the given
// apportionment sum exactly to TotalSocialMoneyPence.
func (m *Memory) CheckApportionmentBalance(_ context.Context, id repro.MoneySupplyApportionmentID) (bool, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	a, ok := m.moneyCapital.apportionments[id]
	if !ok {
		return false, ErrNotFound
	}
	parts := a.DeptIReservePence +
		a.DeptIIReservePence +
		a.WageRotationPence +
		a.IdleHoardPence
	return parts == a.TotalSocialMoneyPence, nil
}

// --- DepartmentMoneyReserve -------------------------------------------------

func (m *Memory) CreateDepartmentMoneyReserve(_ context.Context, r repro.DepartmentMoneyReserve) (repro.DepartmentMoneyReserve, error) {
	if err := r.Department.Validate(); err != nil {
		return repro.DepartmentMoneyReserve{}, err
	}
	if err := r.ReservePurpose.Validate(); err != nil {
		return repro.DepartmentMoneyReserve{}, err
	}
	if r.ID == "" {
		r.ID = repro.NewDepartmentMoneyReserveID()
	}
	m.moneyCapital.mu.Lock()
	defer m.moneyCapital.mu.Unlock()
	if _, exists := m.moneyCapital.departmentReserves[r.ID]; exists {
		return repro.DepartmentMoneyReserve{}, ErrAlreadyExists
	}
	m.moneyCapital.departmentReserves[r.ID] = r
	return r, nil
}

func (m *Memory) GetDepartmentMoneyReserve(_ context.Context, id repro.DepartmentMoneyReserveID) (repro.DepartmentMoneyReserve, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	r, ok := m.moneyCapital.departmentReserves[id]
	if !ok {
		return repro.DepartmentMoneyReserve{}, ErrNotFound
	}
	return r, nil
}

func (m *Memory) ListDepartmentMoneyReserves(_ context.Context) ([]repro.DepartmentMoneyReserve, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	out := make([]repro.DepartmentMoneyReserve, 0, len(m.moneyCapital.departmentReserves))
	for _, r := range m.moneyCapital.departmentReserves {
		out = append(out, r)
	}
	return out, nil
}

// --- CirculatingMoneyMass ---------------------------------------------------

func (m *Memory) CreateCirculatingMoneyMass(_ context.Context, mass repro.CirculatingMoneyMass) (repro.CirculatingMoneyMass, error) {
	if mass.MoneyStockPence < 0 {
		return repro.CirculatingMoneyMass{}, repro.ErrNegativeMoneyStock
	}
	if mass.VelocityPerYearBasisPoints <= 0 {
		return repro.CirculatingMoneyMass{}, repro.ErrZeroVelocity
	}
	if mass.ID == "" {
		mass.ID = repro.NewCirculatingMoneyMassID()
	}
	m.moneyCapital.mu.Lock()
	defer m.moneyCapital.mu.Unlock()
	if _, exists := m.moneyCapital.circulatingMasses[mass.ID]; exists {
		return repro.CirculatingMoneyMass{}, ErrAlreadyExists
	}
	m.moneyCapital.circulatingMasses[mass.ID] = mass
	return mass, nil
}

func (m *Memory) GetCirculatingMoneyMass(_ context.Context, id repro.CirculatingMoneyMassID) (repro.CirculatingMoneyMass, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	mass, ok := m.moneyCapital.circulatingMasses[id]
	if !ok {
		return repro.CirculatingMoneyMass{}, ErrNotFound
	}
	return mass, nil
}

func (m *Memory) ListCirculatingMoneyMasses(_ context.Context) ([]repro.CirculatingMoneyMass, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	out := make([]repro.CirculatingMoneyMass, 0, len(m.moneyCapital.circulatingMasses))
	for _, mass := range m.moneyCapital.circulatingMasses {
		out = append(out, mass)
	}
	return out, nil
}

// --- WageRotationFund -------------------------------------------------------

func (m *Memory) CreateWageRotationFund(_ context.Context, w repro.WageRotationFund) (repro.WageRotationFund, error) {
	if w.ID == "" {
		w.ID = repro.NewWageRotationFundID()
	}
	m.moneyCapital.mu.Lock()
	defer m.moneyCapital.mu.Unlock()
	if _, exists := m.moneyCapital.wageRotationFunds[w.ID]; exists {
		return repro.WageRotationFund{}, ErrAlreadyExists
	}
	m.moneyCapital.wageRotationFunds[w.ID] = w
	return w, nil
}

func (m *Memory) GetWageRotationFund(_ context.Context, id repro.WageRotationFundID) (repro.WageRotationFund, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	w, ok := m.moneyCapital.wageRotationFunds[id]
	if !ok {
		return repro.WageRotationFund{}, ErrNotFound
	}
	return w, nil
}

func (m *Memory) ListWageRotationFunds(_ context.Context) ([]repro.WageRotationFund, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	out := make([]repro.WageRotationFund, 0, len(m.moneyCapital.wageRotationFunds))
	for _, w := range m.moneyCapital.wageRotationFunds {
		out = append(out, w)
	}
	return out, nil
}

// --- InterDepartmentSettlement ----------------------------------------------

func (m *Memory) CreateInterDepartmentSettlement(_ context.Context, s repro.InterDepartmentSettlement) (repro.InterDepartmentSettlement, error) {
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
	m.moneyCapital.mu.Lock()
	defer m.moneyCapital.mu.Unlock()
	if _, exists := m.moneyCapital.interDeptSettlements[s.ID]; exists {
		return repro.InterDepartmentSettlement{}, ErrAlreadyExists
	}
	m.moneyCapital.interDeptSettlements[s.ID] = s
	return s, nil
}

func (m *Memory) GetInterDepartmentSettlement(_ context.Context, id repro.InterDepartmentSettlementID) (repro.InterDepartmentSettlement, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	s, ok := m.moneyCapital.interDeptSettlements[id]
	if !ok {
		return repro.InterDepartmentSettlement{}, ErrNotFound
	}
	return s, nil
}

func (m *Memory) ListInterDepartmentSettlements(_ context.Context) ([]repro.InterDepartmentSettlement, error) {
	m.moneyCapital.mu.RLock()
	defer m.moneyCapital.mu.RUnlock()
	out := make([]repro.InterDepartmentSettlement, 0, len(m.moneyCapital.interDeptSettlements))
	for _, s := range m.moneyCapital.interDeptSettlements {
		out = append(out, s)
	}
	return out, nil
}
