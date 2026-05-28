package store

import (
	"context"
	"sort"
	"sync"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// memoryExtendedReproduction is the in-memory backing store for Vol. II Ch. 21 types.
type memoryExtendedReproduction struct {
	mu              sync.RWMutex
	schemes         map[repro.ExtendedReproductionSchemeID]extendedReproSchemeRow
	reinvestments   map[repro.ExtendedReproductionSchemeID][]repro.Reinvestment
	multiPeriods    map[repro.MultiPeriodSchemeID]repro.MultiPeriodScheme
	compShifts      map[repro.MultiPeriodSchemeID][]repro.CompositionShift
	growthLeads     map[repro.MultiPeriodSchemeID][]repro.DepartmentIGrowthLead
}

// extendedReproSchemeRow stores the scheme root plus the two optional departments.
type extendedReproSchemeRow struct {
	scheme repro.ExtendedReproductionScheme
	deptI  *repro.DepartmentalCapital
	deptII *repro.DepartmentalCapital
}

func newMemoryExtendedReproduction() *memoryExtendedReproduction {
	return &memoryExtendedReproduction{
		schemes:       make(map[repro.ExtendedReproductionSchemeID]extendedReproSchemeRow),
		reinvestments: make(map[repro.ExtendedReproductionSchemeID][]repro.Reinvestment),
		multiPeriods:  make(map[repro.MultiPeriodSchemeID]repro.MultiPeriodScheme),
		compShifts:    make(map[repro.MultiPeriodSchemeID][]repro.CompositionShift),
		growthLeads:   make(map[repro.MultiPeriodSchemeID][]repro.DepartmentIGrowthLead),
	}
}

// assembleExtendedScheme builds an ExtendedReproductionScheme from a row.
// Must be called with the read lock held.
func assembleExtendedScheme(row extendedReproSchemeRow) repro.ExtendedReproductionScheme {
	s := row.scheme
	s.DepartmentI = row.deptI
	s.DepartmentII = row.deptII
	return s
}

// --- ExtendedReproductionScheme ---------------------------------------------

func (m *Memory) CreateExtendedReproductionScheme(_ context.Context, s repro.ExtendedReproductionScheme) (repro.ExtendedReproductionScheme, error) {
	if s.ID == "" {
		s.ID = repro.NewExtendedReproductionSchemeID()
	}
	s.IsBalanced = false
	s.TickCount = 0
	m.extendedRepro.mu.Lock()
	defer m.extendedRepro.mu.Unlock()
	if _, exists := m.extendedRepro.schemes[s.ID]; exists {
		return repro.ExtendedReproductionScheme{}, ErrAlreadyExists
	}
	m.extendedRepro.schemes[s.ID] = extendedReproSchemeRow{scheme: s}
	return s, nil
}

func (m *Memory) GetExtendedReproductionScheme(_ context.Context, id repro.ExtendedReproductionSchemeID) (repro.ExtendedReproductionScheme, error) {
	m.extendedRepro.mu.RLock()
	defer m.extendedRepro.mu.RUnlock()
	row, ok := m.extendedRepro.schemes[id]
	if !ok {
		return repro.ExtendedReproductionScheme{}, ErrNotFound
	}
	return assembleExtendedScheme(row), nil
}

func (m *Memory) ListExtendedReproductionSchemes(_ context.Context, period string) ([]repro.ExtendedReproductionScheme, error) {
	m.extendedRepro.mu.RLock()
	defer m.extendedRepro.mu.RUnlock()
	out := make([]repro.ExtendedReproductionScheme, 0, len(m.extendedRepro.schemes))
	for _, row := range m.extendedRepro.schemes {
		s := assembleExtendedScheme(row)
		if period != "" && s.Period != period {
			continue
		}
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return string(out[i].ID) < string(out[j].ID) })
	return out, nil
}

// --- DepartmentalCapital ----------------------------------------------------

func (m *Memory) AddExtendedDepartmentToScheme(_ context.Context, schemeID repro.ExtendedReproductionSchemeID, dept repro.DepartmentalCapital) (repro.DepartmentalCapital, error) {
	if err := dept.Validate(); err != nil {
		return repro.DepartmentalCapital{}, err
	}
	if dept.ID == "" {
		dept.ID = repro.NewDepartmentalCapitalID()
	}
	m.extendedRepro.mu.Lock()
	defer m.extendedRepro.mu.Unlock()
	row, ok := m.extendedRepro.schemes[schemeID]
	if !ok {
		return repro.DepartmentalCapital{}, ErrNotFound
	}
	switch dept.Department {
	case repro.DepartmentI:
		if row.deptI != nil {
			return repro.DepartmentalCapital{}, ErrAlreadyExists
		}
		row.deptI = &dept
	case repro.DepartmentII:
		if row.deptII != nil {
			return repro.DepartmentalCapital{}, ErrAlreadyExists
		}
		row.deptII = &dept
	default:
		return repro.DepartmentalCapital{}, repro.ErrInvalidDepartment
	}
	m.extendedRepro.schemes[schemeID] = row
	return dept, nil
}

// --- Reinvestment -----------------------------------------------------------

func (m *Memory) CreateReinvestment(_ context.Context, schemeID repro.ExtendedReproductionSchemeID, r repro.Reinvestment) (repro.Reinvestment, error) {
	if r.ID == "" {
		r.ID = repro.NewReinvestmentID()
	}
	r.SchemeID = schemeID
	m.extendedRepro.mu.Lock()
	defer m.extendedRepro.mu.Unlock()
	if _, ok := m.extendedRepro.schemes[schemeID]; !ok {
		return repro.Reinvestment{}, ErrNotFound
	}
	m.extendedRepro.reinvestments[schemeID] = append(m.extendedRepro.reinvestments[schemeID], r)
	return r, nil
}

func (m *Memory) ListReinvestments(_ context.Context, schemeID repro.ExtendedReproductionSchemeID) ([]repro.Reinvestment, error) {
	m.extendedRepro.mu.RLock()
	defer m.extendedRepro.mu.RUnlock()
	if _, ok := m.extendedRepro.schemes[schemeID]; !ok {
		return nil, ErrNotFound
	}
	src := m.extendedRepro.reinvestments[schemeID]
	out := make([]repro.Reinvestment, len(src))
	copy(out, src)
	sort.Slice(out, func(i, j int) bool { return out[i].CycleNumber < out[j].CycleNumber })
	return out, nil
}

// --- TickExtendedScheme -----------------------------------------------------

func (m *Memory) TickExtendedScheme(_ context.Context, schemeID repro.ExtendedReproductionSchemeID) (repro.ExtendedReproductionScheme, error) {
	m.extendedRepro.mu.Lock()
	defer m.extendedRepro.mu.Unlock()
	row, ok := m.extendedRepro.schemes[schemeID]
	if !ok {
		return repro.ExtendedReproductionScheme{}, ErrNotFound
	}
	if row.deptI == nil || row.deptII == nil {
		return repro.ExtendedReproductionScheme{}, ErrNotFound
	}

	cycleNumber := row.scheme.TickCount + 1

	rI, err := repro.ComputeReinvestment(*row.deptI, row.scheme.AccumulationRateIBps)
	if err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}
	rI.SchemeID = schemeID
	rI.CycleNumber = cycleNumber

	rII, err := repro.ComputeReinvestment(*row.deptII, row.scheme.AccumulationRateIIBps)
	if err != nil {
		return repro.ExtendedReproductionScheme{}, err
	}
	rII.SchemeID = schemeID
	rII.CycleNumber = cycleNumber

	m.extendedRepro.reinvestments[schemeID] = append(m.extendedRepro.reinvestments[schemeID], rI, rII)

	newDeptI := repro.ComputeNextPeriodCapital(*row.deptI, rI)
	newDeptII := repro.ComputeNextPeriodCapital(*row.deptII, rII)

	row.deptI = &newDeptI
	row.deptII = &newDeptII

	balanced := repro.ComputeExtendedBalance(newDeptI, newDeptII, rI, rII)
	row.scheme.IsBalanced = balanced
	row.scheme.TickCount = cycleNumber

	m.extendedRepro.schemes[schemeID] = row
	return assembleExtendedScheme(row), nil
}

// --- GetExtendedMoneyRequirement --------------------------------------------

func (m *Memory) GetExtendedMoneyRequirement(_ context.Context, schemeID repro.ExtendedReproductionSchemeID, baseMoneyPence int64) (repro.ExtendedMoneyRequirement, error) {
	m.extendedRepro.mu.RLock()
	defer m.extendedRepro.mu.RUnlock()
	row, ok := m.extendedRepro.schemes[schemeID]
	if !ok {
		return repro.ExtendedMoneyRequirement{}, ErrNotFound
	}
	s := assembleExtendedScheme(row)
	return repro.ComputeExtendedMoneyRequirement(s, baseMoneyPence)
}

// --- MultiPeriodScheme ------------------------------------------------------

func (m *Memory) CreateMultiPeriodScheme(_ context.Context, mp repro.MultiPeriodScheme) (repro.MultiPeriodScheme, error) {
	if mp.ID == "" {
		mp.ID = repro.NewMultiPeriodSchemeID()
	}
	m.extendedRepro.mu.Lock()
	defer m.extendedRepro.mu.Unlock()

	// Verify all linked schemes exist and collect them.
	schemes := make([]repro.ExtendedReproductionScheme, len(mp.SchemeIDs))
	for i, sid := range mp.SchemeIDs {
		row, ok := m.extendedRepro.schemes[sid]
		if !ok {
			return repro.MultiPeriodScheme{}, ErrNotFound
		}
		schemes[i] = assembleExtendedScheme(row)
	}

	m.extendedRepro.multiPeriods[mp.ID] = mp

	// Compute CompositionShift and DepartmentIGrowthLead for each consecutive pair.
	for i := 1; i < len(schemes); i++ {
		prev := schemes[i-1]
		curr := schemes[i]
		cycle := int64(i)

		var deptICompPrev, deptICompCurr, deptIICompPrev, deptIICompCurr int64
		if prev.DepartmentI != nil {
			deptICompPrev = repro.ComputeCompositionBps(*prev.DepartmentI)
		}
		if curr.DepartmentI != nil {
			deptICompCurr = repro.ComputeCompositionBps(*curr.DepartmentI)
		}
		if prev.DepartmentII != nil {
			deptIICompPrev = repro.ComputeCompositionBps(*prev.DepartmentII)
		}
		if curr.DepartmentII != nil {
			deptIICompCurr = repro.ComputeCompositionBps(*curr.DepartmentII)
		}

		// Social composition: average of dept I and dept II compositions in the current period.
		socialComp := (deptICompCurr + deptIICompCurr) / 2
		_ = deptICompPrev
		_ = deptIICompPrev

		cs := repro.CompositionShift{
			ID:                   repro.NewCompositionShiftID(),
			MultiPeriodID:        mp.ID,
			CycleNumber:          cycle,
			DeptICompositionBps:  deptICompCurr,
			DeptIICompositionBps: deptIICompCurr,
			SocialCompositionBps: socialComp,
		}
		m.extendedRepro.compShifts[mp.ID] = append(m.extendedRepro.compShifts[mp.ID], cs)

		// Growth leads.
		var deptIGrowth, deptIIGrowth int64
		if prev.DepartmentI != nil && curr.DepartmentI != nil {
			deptIGrowth = repro.ComputeGrowthBps(*prev.DepartmentI, *curr.DepartmentI)
		}
		if prev.DepartmentII != nil && curr.DepartmentII != nil {
			deptIIGrowth = repro.ComputeGrowthBps(*prev.DepartmentII, *curr.DepartmentII)
		}
		gl := repro.DepartmentIGrowthLead{
			ID:              repro.NewDepartmentIGrowthLeadID(),
			MultiPeriodID:   mp.ID,
			CycleNumber:     cycle,
			DeptIGrowthBps:  deptIGrowth,
			DeptIIGrowthBps: deptIIGrowth,
			LeadConfirmed:   deptIGrowth > deptIIGrowth,
		}
		m.extendedRepro.growthLeads[mp.ID] = append(m.extendedRepro.growthLeads[mp.ID], gl)
	}

	return mp, nil
}

func (m *Memory) GetMultiPeriodScheme(_ context.Context, id repro.MultiPeriodSchemeID) (repro.MultiPeriodScheme, error) {
	m.extendedRepro.mu.RLock()
	defer m.extendedRepro.mu.RUnlock()
	mp, ok := m.extendedRepro.multiPeriods[id]
	if !ok {
		return repro.MultiPeriodScheme{}, ErrNotFound
	}
	return mp, nil
}

func (m *Memory) ListCompositionShifts(_ context.Context, mpID repro.MultiPeriodSchemeID) ([]repro.CompositionShift, error) {
	m.extendedRepro.mu.RLock()
	defer m.extendedRepro.mu.RUnlock()
	if _, ok := m.extendedRepro.multiPeriods[mpID]; !ok {
		return nil, ErrNotFound
	}
	src := m.extendedRepro.compShifts[mpID]
	out := make([]repro.CompositionShift, len(src))
	copy(out, src)
	sort.Slice(out, func(i, j int) bool { return out[i].CycleNumber < out[j].CycleNumber })
	return out, nil
}

func (m *Memory) ListGrowthLeads(_ context.Context, mpID repro.MultiPeriodSchemeID) ([]repro.DepartmentIGrowthLead, error) {
	m.extendedRepro.mu.RLock()
	defer m.extendedRepro.mu.RUnlock()
	if _, ok := m.extendedRepro.multiPeriods[mpID]; !ok {
		return nil, ErrNotFound
	}
	src := m.extendedRepro.growthLeads[mpID]
	out := make([]repro.DepartmentIGrowthLead, len(src))
	copy(out, src)
	sort.Slice(out, func(i, j int) bool { return out[i].CycleNumber < out[j].CycleNumber })
	return out, nil
}
