package store

import (
	"context"
	"sync"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// memorySimpleReproduction is the in-memory backing store for Vol. II Ch. 20 types.
type memorySimpleReproduction struct {
	mu        sync.RWMutex
	schemes   map[repro.SimpleReproductionSchemeID]simpleReproSchemeRow
	exchanges map[repro.SimpleReproductionSchemeID][]repro.InterDepartmentExchange
	ticks     map[repro.SimpleReproductionSchemeID][]repro.ReproductionTick
	loops     map[repro.SimpleReproductionSchemeID][]repro.MoneyClosedLoop
}

// simpleReproSchemeRow stores the scheme root plus the two optional departments.
type simpleReproSchemeRow struct {
	scheme repro.SimpleReproductionScheme
	deptI  *repro.DepartmentalCapital
	deptII *repro.DepartmentalCapital
}

func newMemorySimpleReproduction() *memorySimpleReproduction {
	return &memorySimpleReproduction{
		schemes:   make(map[repro.SimpleReproductionSchemeID]simpleReproSchemeRow),
		exchanges: make(map[repro.SimpleReproductionSchemeID][]repro.InterDepartmentExchange),
		ticks:     make(map[repro.SimpleReproductionSchemeID][]repro.ReproductionTick),
		loops:     make(map[repro.SimpleReproductionSchemeID][]repro.MoneyClosedLoop),
	}
}

// assembleScheme builds a SimpleReproductionScheme from a row, recomputing
// IsBalanced. Must be called with the read lock held.
func assembleScheme(row simpleReproSchemeRow) repro.SimpleReproductionScheme {
	s := row.scheme
	s.DepartmentI = row.deptI
	s.DepartmentII = row.deptII
	s.IsBalanced = s.CheckBalance()
	return s
}

// --- SimpleReproductionScheme -----------------------------------------------

func (m *Memory) CreateSimpleReproductionScheme(_ context.Context, s repro.SimpleReproductionScheme) (repro.SimpleReproductionScheme, error) {
	if s.ID == "" {
		s.ID = repro.NewSimpleReproductionSchemeID()
	}
	s.IsBalanced = false
	s.TickCount = 0
	m.simpleRepro.mu.Lock()
	defer m.simpleRepro.mu.Unlock()
	if _, exists := m.simpleRepro.schemes[s.ID]; exists {
		return repro.SimpleReproductionScheme{}, ErrAlreadyExists
	}
	m.simpleRepro.schemes[s.ID] = simpleReproSchemeRow{scheme: s}
	return s, nil
}

func (m *Memory) GetSimpleReproductionScheme(_ context.Context, id repro.SimpleReproductionSchemeID) (repro.SimpleReproductionScheme, error) {
	m.simpleRepro.mu.RLock()
	defer m.simpleRepro.mu.RUnlock()
	row, ok := m.simpleRepro.schemes[id]
	if !ok {
		return repro.SimpleReproductionScheme{}, ErrNotFound
	}
	return assembleScheme(row), nil
}

func (m *Memory) ListSimpleReproductionSchemes(_ context.Context, period string) ([]repro.SimpleReproductionScheme, error) {
	m.simpleRepro.mu.RLock()
	defer m.simpleRepro.mu.RUnlock()
	out := make([]repro.SimpleReproductionScheme, 0, len(m.simpleRepro.schemes))
	for _, row := range m.simpleRepro.schemes {
		s := assembleScheme(row)
		if period != "" && s.Period != period {
			continue
		}
		out = append(out, s)
	}
	return out, nil
}

// --- DepartmentalCapital ----------------------------------------------------

func (m *Memory) AddDepartmentToScheme(_ context.Context, id repro.SimpleReproductionSchemeID, d repro.DepartmentalCapital) (repro.DepartmentalCapital, error) {
	if err := d.Validate(); err != nil {
		return repro.DepartmentalCapital{}, err
	}
	if d.ID == "" {
		d.ID = repro.NewDepartmentalCapitalID()
	}
	d.SchemeID = id
	m.simpleRepro.mu.Lock()
	defer m.simpleRepro.mu.Unlock()
	row, ok := m.simpleRepro.schemes[id]
	if !ok {
		return repro.DepartmentalCapital{}, ErrNotFound
	}
	switch d.Department {
	case repro.DepartmentI:
		if row.deptI != nil {
			return repro.DepartmentalCapital{}, repro.ErrDepartmentAlreadySet
		}
		row.deptI = &d
	case repro.DepartmentII:
		if row.deptII != nil {
			return repro.DepartmentalCapital{}, repro.ErrDepartmentAlreadySet
		}
		row.deptII = &d
	default:
		return repro.DepartmentalCapital{}, repro.ErrInvalidDepartment
	}
	// Recompute balance.
	s := row.scheme
	s.DepartmentI = row.deptI
	s.DepartmentII = row.deptII
	s.IsBalanced = s.CheckBalance()
	row.scheme = s
	m.simpleRepro.schemes[id] = row
	return d, nil
}

// --- InterDepartmentExchange ------------------------------------------------

func (m *Memory) RecordInterDepartmentExchange(_ context.Context, id repro.SimpleReproductionSchemeID, e repro.InterDepartmentExchange) (repro.InterDepartmentExchange, error) {
	if err := e.Validate(); err != nil {
		return repro.InterDepartmentExchange{}, err
	}
	if e.ID == "" {
		e.ID = repro.NewInterDepartmentExchangeID()
	}
	e.SchemeID = id
	m.simpleRepro.mu.Lock()
	defer m.simpleRepro.mu.Unlock()
	if _, ok := m.simpleRepro.schemes[id]; !ok {
		return repro.InterDepartmentExchange{}, ErrNotFound
	}
	m.simpleRepro.exchanges[id] = append(m.simpleRepro.exchanges[id], e)
	return e, nil
}

func (m *Memory) ListInterDepartmentExchanges(_ context.Context, schemeID repro.SimpleReproductionSchemeID) ([]repro.InterDepartmentExchange, error) {
	m.simpleRepro.mu.RLock()
	defer m.simpleRepro.mu.RUnlock()
	if _, ok := m.simpleRepro.schemes[schemeID]; !ok {
		return nil, ErrNotFound
	}
	src := m.simpleRepro.exchanges[schemeID]
	out := make([]repro.InterDepartmentExchange, len(src))
	copy(out, src)
	return out, nil
}

// --- ReproductionTick -------------------------------------------------------

func (m *Memory) AdvanceReproductionTick(_ context.Context, id repro.SimpleReproductionSchemeID) (repro.ReproductionTick, error) {
	m.simpleRepro.mu.Lock()
	defer m.simpleRepro.mu.Unlock()
	row, ok := m.simpleRepro.schemes[id]
	if !ok {
		return repro.ReproductionTick{}, ErrNotFound
	}
	s := row.scheme
	s.TickCount++
	row.scheme = s
	m.simpleRepro.schemes[id] = row

	assembled := assembleScheme(row)
	tick := repro.ReproductionTick{
		ID:         repro.NewReproductionTickID(),
		SchemeID:   id,
		TickNumber: s.TickCount,
		Period:     s.Period,
		IsBalanced: assembled.IsBalanced,
	}
	m.simpleRepro.ticks[id] = append(m.simpleRepro.ticks[id], tick)
	return tick, nil
}

// --- BalanceCheck -----------------------------------------------------------

func (m *Memory) CheckSchemeBalance(_ context.Context, id repro.SimpleReproductionSchemeID) (repro.BalanceCheckResult, error) {
	m.simpleRepro.mu.RLock()
	defer m.simpleRepro.mu.RUnlock()
	row, ok := m.simpleRepro.schemes[id]
	if !ok {
		return repro.BalanceCheckResult{}, ErrNotFound
	}
	s := assembleScheme(row)
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
	return res, nil
}

// --- MoneyClosedLoop --------------------------------------------------------

func (m *Memory) RecordMoneyClosedLoop(_ context.Context, id repro.SimpleReproductionSchemeID, loop repro.MoneyClosedLoop) (repro.MoneyClosedLoop, error) {
	if loop.ID == "" {
		loop.ID = repro.NewMoneyClosedLoopID()
	}
	loop.SchemeID = id
	m.simpleRepro.mu.Lock()
	defer m.simpleRepro.mu.Unlock()
	if _, ok := m.simpleRepro.schemes[id]; !ok {
		return repro.MoneyClosedLoop{}, ErrNotFound
	}
	m.simpleRepro.loops[id] = append(m.simpleRepro.loops[id], loop)
	return loop, nil
}
