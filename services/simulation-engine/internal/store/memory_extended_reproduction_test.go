package store_test

import (
	"context"
	"errors"
	"testing"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// ch21 canonical fixture — Dept I: 4000c+1000v+1000s, Dept II: 1500c+750v+750s.

func newCh21ExtDeptI(schemeID repro.ExtendedReproductionSchemeID) repro.DepartmentalCapital {
	return repro.DepartmentalCapital{
		ID:            repro.NewDepartmentalCapitalID(),
		SchemeID:      repro.SimpleReproductionSchemeID(schemeID),
		Department:    repro.DepartmentI,
		ConstantPence: 4_000_000,
		VariablePence: 1_000_000,
		SurplusPence:  1_000_000,
		TotalPence:    6_000_000,
	}
}

func newCh21ExtDeptII(schemeID repro.ExtendedReproductionSchemeID) repro.DepartmentalCapital {
	return repro.DepartmentalCapital{
		ID:            repro.NewDepartmentalCapitalID(),
		SchemeID:      repro.SimpleReproductionSchemeID(schemeID),
		Department:    repro.DepartmentII,
		ConstantPence: 1_500_000,
		VariablePence: 750_000,
		SurplusPence:  750_000,
		TotalPence:    3_000_000,
	}
}

// buildExtendedSchemeWithDepts creates a fully-populated extended scheme using
// Marx's canonical Ch. 21 fixture with the given accumulation rates.
func buildExtendedSchemeWithDepts(
	t *testing.T,
	mem *store.Memory,
	accI, accII int64,
) repro.ExtendedReproductionScheme {
	t.Helper()
	ctx := context.Background()

	s, err := mem.CreateExtendedReproductionScheme(ctx, repro.ExtendedReproductionScheme{
		ID:                    repro.NewExtendedReproductionSchemeID(),
		Period:                "Year I",
		AccumulationRateIBps:  accI,
		AccumulationRateIIBps: accII,
	})
	if err != nil {
		t.Fatalf("CreateExtendedReproductionScheme: %v", err)
	}

	dI := newCh21ExtDeptI(s.ID)
	if _, err := mem.AddExtendedDepartmentToScheme(ctx, s.ID, dI); err != nil {
		t.Fatalf("AddExtendedDepartmentToScheme(I): %v", err)
	}

	dII := newCh21ExtDeptII(s.ID)
	if _, err := mem.AddExtendedDepartmentToScheme(ctx, s.ID, dII); err != nil {
		t.Fatalf("AddExtendedDepartmentToScheme(II): %v", err)
	}

	got, err := mem.GetExtendedReproductionScheme(ctx, s.ID)
	if err != nil {
		t.Fatalf("GetExtendedReproductionScheme: %v", err)
	}
	return got
}

// --- CreateExtendedReproductionScheme / GetExtendedReproductionScheme --------

func TestMemoryExtendedReproduction_CreateAndGet(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	created, err := mem.CreateExtendedReproductionScheme(ctx, repro.ExtendedReproductionScheme{
		ID:                    repro.NewExtendedReproductionSchemeID(),
		Period:                "Year I",
		AccumulationRateIBps:  5000,
		AccumulationRateIIBps: 3000,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.TickCount != 0 {
		t.Fatalf("TickCount: want 0 on creation, got %d", created.TickCount)
	}

	got, err := mem.GetExtendedReproductionScheme(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Period != "Year I" {
		t.Fatalf("Period: want %q, got %q", "Year I", got.Period)
	}
	if got.AccumulationRateIBps != 5000 {
		t.Fatalf("AccumulationRateIBps: want 5000, got %d", got.AccumulationRateIBps)
	}
	if got.AccumulationRateIIBps != 3000 {
		t.Fatalf("AccumulationRateIIBps: want 3000, got %d", got.AccumulationRateIIBps)
	}
}

func TestMemoryExtendedReproduction_GetNotFound(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	_, err := mem.GetExtendedReproductionScheme(ctx, "ghost-extended-id")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// --- AddExtendedDepartmentToScheme ------------------------------------------

func TestMemoryExtendedReproduction_AddDepartment_BothDepts(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	s := buildExtendedSchemeWithDepts(t, mem, 5000, 3000)

	if s.DepartmentI == nil {
		t.Fatal("expected DepartmentI to be non-nil")
	}
	if s.DepartmentII == nil {
		t.Fatal("expected DepartmentII to be non-nil")
	}
	if s.DepartmentI.ConstantPence != 4_000_000 {
		t.Fatalf("DeptI.ConstantPence: want 4_000_000, got %d", s.DepartmentI.ConstantPence)
	}
	if s.DepartmentII.ConstantPence != 1_500_000 {
		t.Fatalf("DeptII.ConstantPence: want 1_500_000, got %d", s.DepartmentII.ConstantPence)
	}
}

func TestMemoryExtendedReproduction_AddDepartment_AlreadyExists(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	s, err := mem.CreateExtendedReproductionScheme(ctx, repro.ExtendedReproductionScheme{
		ID:                    repro.NewExtendedReproductionSchemeID(),
		Period:                "Year I",
		AccumulationRateIBps:  5000,
		AccumulationRateIIBps: 3000,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	dI := newCh21ExtDeptI(s.ID)
	if _, err := mem.AddExtendedDepartmentToScheme(ctx, s.ID, dI); err != nil {
		t.Fatalf("first add Dept I: %v", err)
	}
	// Adding the same department again must return ErrAlreadyExists.
	_, err = mem.AddExtendedDepartmentToScheme(ctx, s.ID, dI)
	if !errors.Is(err, store.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists on duplicate, got %v", err)
	}
}

func TestMemoryExtendedReproduction_AddDepartment_SchemeNotFound(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	ghostID := repro.ExtendedReproductionSchemeID("ghost-scheme-extended")
	dI := newCh21ExtDeptI(ghostID)
	_, err := mem.AddExtendedDepartmentToScheme(ctx, ghostID, dI)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// --- CreateReinvestment / ListReinvestments ----------------------------------

func TestMemoryExtendedReproduction_CreateReinvestment_AndList(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	s, err := mem.CreateExtendedReproductionScheme(ctx, repro.ExtendedReproductionScheme{
		ID:                    repro.NewExtendedReproductionSchemeID(),
		Period:                "Year I",
		AccumulationRateIBps:  5000,
		AccumulationRateIIBps: 3000,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	r := repro.Reinvestment{
		SchemeID:             s.ID,
		Department:           repro.DepartmentI,
		DeltaConstantPence:   333_333,
		DeltaVariablePence:   166_667,
		ConsumedSurplusPence: 500_000,
		CycleNumber:          1,
	}
	created, err := mem.CreateReinvestment(ctx, s.ID, r)
	if err != nil {
		t.Fatalf("CreateReinvestment: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty Reinvestment ID")
	}
	if created.SchemeID != s.ID {
		t.Fatalf("SchemeID: want %q, got %q", s.ID, created.SchemeID)
	}

	list, err := mem.ListReinvestments(ctx, s.ID)
	if err != nil {
		t.Fatalf("ListReinvestments: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 reinvestment, got %d", len(list))
	}
	if list[0].DeltaConstantPence != 333_333 {
		t.Fatalf("DeltaConstantPence: want 333_333, got %d", list[0].DeltaConstantPence)
	}
}

func TestMemoryExtendedReproduction_ListReinvestments_Empty(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	s, err := mem.CreateExtendedReproductionScheme(ctx, repro.ExtendedReproductionScheme{
		ID:                    repro.NewExtendedReproductionSchemeID(),
		Period:                "Year I",
		AccumulationRateIBps:  5000,
		AccumulationRateIIBps: 3000,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	list, err := mem.ListReinvestments(ctx, s.ID)
	if err != nil {
		t.Fatalf("ListReinvestments on empty: %v", err)
	}
	if list == nil {
		t.Fatal("expected non-nil (possibly empty) slice from ListReinvestments")
	}
	if len(list) != 0 {
		t.Fatalf("expected 0 reinvestments, got %d", len(list))
	}
}

func TestMemoryExtendedReproduction_CreateReinvestment_SchemeNotFound(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	ghostID := repro.ExtendedReproductionSchemeID("ghost-scheme")
	r := repro.Reinvestment{
		SchemeID:             ghostID,
		Department:           repro.DepartmentI,
		DeltaConstantPence:   100_000,
		DeltaVariablePence:   25_000,
		ConsumedSurplusPence: 375_000,
		CycleNumber:          1,
	}
	_, err := mem.CreateReinvestment(ctx, ghostID, r)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// --- TickExtendedScheme -----------------------------------------------------

func TestMemoryExtendedReproduction_TickScheme(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	// Create scheme with both departments.
	s := buildExtendedSchemeWithDepts(t, mem, 5000, 3000)

	// Tick once.
	ticked, err := mem.TickExtendedScheme(ctx, s.ID)
	if err != nil {
		t.Fatalf("TickExtendedScheme: %v", err)
	}
	if ticked.TickCount != 1 {
		t.Fatalf("TickCount: want 1, got %d", ticked.TickCount)
	}
	if ticked.DepartmentI == nil {
		t.Fatal("expected DepartmentI to be non-nil after tick")
	}
	if ticked.DepartmentII == nil {
		t.Fatal("expected DepartmentII to be non-nil after tick")
	}
	// After accumulation Dept I must have grown (50% rate)
	if ticked.DepartmentI.TotalPence <= 6_000_000 {
		t.Fatalf("Dept I must have grown after 50%% accumulation: got %d",
			ticked.DepartmentI.TotalPence)
	}
	// Reinvestments must have been recorded (2 entries: one per dept).
	reinvs, err := mem.ListReinvestments(ctx, s.ID)
	if err != nil {
		t.Fatalf("ListReinvestments after tick: %v", err)
	}
	if len(reinvs) != 2 {
		t.Fatalf("expected 2 reinvestments (one per dept) after tick, got %d", len(reinvs))
	}
}

func TestMemoryExtendedReproduction_TickScheme_NoDepts(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	s, err := mem.CreateExtendedReproductionScheme(ctx, repro.ExtendedReproductionScheme{
		ID:                    repro.NewExtendedReproductionSchemeID(),
		Period:                "Year I",
		AccumulationRateIBps:  5000,
		AccumulationRateIIBps: 3000,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// Tick without adding departments — expect an error.
	_, err = mem.TickExtendedScheme(ctx, s.ID)
	if err == nil {
		t.Fatal("expected error when ticking scheme with no departments")
	}
	// Should be ErrNotFound (departments are missing).
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for missing departments, got %v", err)
	}
}

func TestMemoryExtendedReproduction_TickScheme_NotFound(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	_, err := mem.TickExtendedScheme(ctx, "ghost-scheme")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for ghost scheme, got %v", err)
	}
}

// --- GetExtendedMoneyRequirement --------------------------------------------

func TestMemoryExtendedReproduction_GetMoneyRequirement(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	s := buildExtendedSchemeWithDepts(t, mem, 5000, 3000)

	req, err := mem.GetExtendedMoneyRequirement(ctx, s.ID, 1_000_000)
	if err != nil {
		t.Fatalf("GetExtendedMoneyRequirement: %v", err)
	}
	if req.BaseMoneyPence != 1_000_000 {
		t.Fatalf("BaseMoneyPence: want 1_000_000, got %d", req.BaseMoneyPence)
	}
	if req.TotalMoneyPence != req.BaseMoneyPence+req.AdditionalMoneyPence {
		t.Fatalf("TotalMoney invariant: %d + %d != %d",
			req.BaseMoneyPence, req.AdditionalMoneyPence, req.TotalMoneyPence)
	}
	if req.AdditionalMoneyPence < 0 {
		t.Fatalf("AdditionalMoneyPence must be >= 0, got %d", req.AdditionalMoneyPence)
	}
}

func TestMemoryExtendedReproduction_GetMoneyRequirement_NotFound(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	_, err := mem.GetExtendedMoneyRequirement(ctx, "ghost-scheme", 1_000_000)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// --- CreateMultiPeriodScheme / GetMultiPeriodScheme -------------------------

func TestMemoryExtendedReproduction_CreateMultiPeriodScheme(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	// Build s1: Marx's Year I baseline (4000c+1000v+1000s, 1500c+750v+750s).
	s1 := buildExtendedSchemeWithDepts(t, mem, 5000, 3000)

	// Build s2: Year II with larger capitals (representing growth from accumulation).
	// Dept I grows from 6_000_000 to ~6_500_000 (≈8.3% growth).
	// Dept II grows from 3_000_000 to ~3_225_000 (≈7.5% growth).
	// This confirms Dept I growth lead (8.3% > 7.5%).
	s2, err := mem.CreateExtendedReproductionScheme(ctx, repro.ExtendedReproductionScheme{
		ID:                    repro.NewExtendedReproductionSchemeID(),
		Period:                "Year II",
		AccumulationRateIBps:  5000,
		AccumulationRateIIBps: 3000,
	})
	if err != nil {
		t.Fatalf("CreateExtendedReproductionScheme(s2): %v", err)
	}
	// Year II Dept I: 4333333c + 1166667v + 1166667s ≈ 6666667 total (c/v≈4)
	dIS2 := repro.DepartmentalCapital{
		Department:    repro.DepartmentI,
		ConstantPence: 4_333_333,
		VariablePence: 1_166_667,
		SurplusPence:  1_166_667,
		TotalPence:    6_666_667,
	}
	if _, err := mem.AddExtendedDepartmentToScheme(ctx, s2.ID, dIS2); err != nil {
		t.Fatalf("AddExtendedDepartmentToScheme(s2, I): %v", err)
	}
	// Year II Dept II: 1575000c + 787500v + 787500s = 3_150_000 total
	dIIS2 := repro.DepartmentalCapital{
		Department:    repro.DepartmentII,
		ConstantPence: 1_575_000,
		VariablePence: 787_500,
		SurplusPence:  787_500,
		TotalPence:    3_150_000,
	}
	if _, err := mem.AddExtendedDepartmentToScheme(ctx, s2.ID, dIIS2); err != nil {
		t.Fatalf("AddExtendedDepartmentToScheme(s2, II): %v", err)
	}

	mp, err := mem.CreateMultiPeriodScheme(ctx, repro.MultiPeriodScheme{
		ID:        repro.NewMultiPeriodSchemeID(),
		Label:     "Marx Year I → Year II",
		SchemeIDs: []repro.ExtendedReproductionSchemeID{s1.ID, s2.ID},
	})
	if err != nil {
		t.Fatalf("CreateMultiPeriodScheme: %v", err)
	}
	if mp.ID == "" {
		t.Fatal("expected non-empty MultiPeriodScheme ID")
	}

	// ListCompositionShifts should have exactly 1 shift (one pair).
	shifts, err := mem.ListCompositionShifts(ctx, mp.ID)
	if err != nil {
		t.Fatalf("ListCompositionShifts: %v", err)
	}
	if shifts == nil {
		t.Fatal("expected non-nil shifts slice")
	}
	if len(shifts) != 1 {
		t.Fatalf("expected 1 CompositionShift for two schemes, got %d", len(shifts))
	}
	if shifts[0].CycleNumber != 1 {
		t.Fatalf("CycleNumber: want 1, got %d", shifts[0].CycleNumber)
	}

	// ListGrowthLeads should have exactly 1 lead.
	leads, err := mem.ListGrowthLeads(ctx, mp.ID)
	if err != nil {
		t.Fatalf("ListGrowthLeads: %v", err)
	}
	if leads == nil {
		t.Fatal("expected non-nil growth leads slice")
	}
	if len(leads) != 1 {
		t.Fatalf("expected 1 DepartmentIGrowthLead for two schemes, got %d", len(leads))
	}
	// Dept I: 6_000_000 → 6_666_667 ≈ +11.1%; Dept II: 3_000_000 → 3_150_000 = +5%.
	// Dept I growth > Dept II growth → LeadConfirmed=true.
	if !leads[0].LeadConfirmed {
		t.Fatalf("expected LeadConfirmed=true: DeptI=%d bps, DeptII=%d bps",
			leads[0].DeptIGrowthBps, leads[0].DeptIIGrowthBps)
	}
}

func TestMemoryExtendedReproduction_MultiPeriod_GetNotFound(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	_, err := mem.GetMultiPeriodScheme(ctx, "ghost-multi-period-id")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryExtendedReproduction_ListCompositionShifts_Empty(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	// A MultiPeriodScheme with a single scheme produces no pairs → no shifts.
	s := buildExtendedSchemeWithDepts(t, mem, 5000, 3000)
	mp, err := mem.CreateMultiPeriodScheme(ctx, repro.MultiPeriodScheme{
		ID:        repro.NewMultiPeriodSchemeID(),
		Label:     "single period",
		SchemeIDs: []repro.ExtendedReproductionSchemeID{s.ID},
	})
	if err != nil {
		t.Fatalf("CreateMultiPeriodScheme: %v", err)
	}

	shifts, err := mem.ListCompositionShifts(ctx, mp.ID)
	if err != nil {
		t.Fatalf("ListCompositionShifts: %v", err)
	}
	if shifts == nil {
		t.Fatal("expected non-nil (empty) shifts slice, got nil")
	}
}

func TestMemoryExtendedReproduction_ListGrowthLeads_Empty(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	s := buildExtendedSchemeWithDepts(t, mem, 5000, 3000)
	mp, err := mem.CreateMultiPeriodScheme(ctx, repro.MultiPeriodScheme{
		ID:        repro.NewMultiPeriodSchemeID(),
		Label:     "single period",
		SchemeIDs: []repro.ExtendedReproductionSchemeID{s.ID},
	})
	if err != nil {
		t.Fatalf("CreateMultiPeriodScheme: %v", err)
	}

	leads, err := mem.ListGrowthLeads(ctx, mp.ID)
	if err != nil {
		t.Fatalf("ListGrowthLeads: %v", err)
	}
	if leads == nil {
		t.Fatal("expected non-nil (empty) growth leads slice, got nil")
	}
}

// --- Zero accumulation rate — simple reproduction limiting case -------------

func TestMemoryExtendedReproduction_ZeroAccumulationRate_SimpleLimitingCase(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	// With 0% accumulation, no capital growth should occur.
	s := buildExtendedSchemeWithDepts(t, mem, 0, 0)

	prevDeptI := *s.DepartmentI
	prevDeptII := *s.DepartmentII

	ticked, err := mem.TickExtendedScheme(ctx, s.ID)
	if err != nil {
		t.Fatalf("TickExtendedScheme with 0%% rate: %v", err)
	}
	if ticked.DepartmentI == nil || ticked.DepartmentII == nil {
		t.Fatal("expected both departments after tick")
	}
	// With zero accumulation, constant and variable pence must not change.
	if ticked.DepartmentI.ConstantPence != prevDeptI.ConstantPence {
		t.Fatalf("DeptI.ConstantPence changed at 0%% rate: before=%d, after=%d",
			prevDeptI.ConstantPence, ticked.DepartmentI.ConstantPence)
	}
	if ticked.DepartmentI.VariablePence != prevDeptI.VariablePence {
		t.Fatalf("DeptI.VariablePence changed at 0%% rate: before=%d, after=%d",
			prevDeptI.VariablePence, ticked.DepartmentI.VariablePence)
	}
	if ticked.DepartmentII.ConstantPence != prevDeptII.ConstantPence {
		t.Fatalf("DeptII.ConstantPence changed at 0%% rate: before=%d, after=%d",
			prevDeptII.ConstantPence, ticked.DepartmentII.ConstantPence)
	}
	// Reinvestments should show DeltaC=0, DeltaV=0 for both departments.
	reinvs, err := mem.ListReinvestments(ctx, s.ID)
	if err != nil {
		t.Fatalf("ListReinvestments: %v", err)
	}
	for _, rv := range reinvs {
		if rv.DeltaConstantPence != 0 {
			t.Fatalf("DeltaConstantPence must be 0 at zero rate, got %d (dept=%s)",
				rv.DeltaConstantPence, rv.Department)
		}
		if rv.DeltaVariablePence != 0 {
			t.Fatalf("DeltaVariablePence must be 0 at zero rate, got %d (dept=%s)",
				rv.DeltaVariablePence, rv.Department)
		}
	}
}

// --- ListCompositionShifts sorted order -------------------------------------

func TestMemoryExtendedReproduction_ListCompositionShifts_SortedByCycle(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	// Build three schemes: s1 → s2 → s3.
	s1 := buildExtendedSchemeWithDepts(t, mem, 5000, 3000)

	buildSuccessor := func(prev repro.ExtendedReproductionScheme) repro.ExtendedReproductionScheme {
		t.Helper()
		ticked, err := mem.TickExtendedScheme(ctx, prev.ID)
		if err != nil {
			t.Fatalf("tick: %v", err)
		}
		s, err := mem.CreateExtendedReproductionScheme(ctx, repro.ExtendedReproductionScheme{
			ID:                    repro.NewExtendedReproductionSchemeID(),
			Period:                "next",
			AccumulationRateIBps:  5000,
			AccumulationRateIIBps: 3000,
		})
		if err != nil {
			t.Fatalf("create successor: %v", err)
		}
		if ticked.DepartmentI != nil {
			d := *ticked.DepartmentI
			d.ID = repro.NewDepartmentalCapitalID()
			if _, err := mem.AddExtendedDepartmentToScheme(ctx, s.ID, d); err != nil {
				t.Fatalf("add deptI: %v", err)
			}
		}
		if ticked.DepartmentII != nil {
			d := *ticked.DepartmentII
			d.ID = repro.NewDepartmentalCapitalID()
			if _, err := mem.AddExtendedDepartmentToScheme(ctx, s.ID, d); err != nil {
				t.Fatalf("add deptII: %v", err)
			}
		}
		got, err := mem.GetExtendedReproductionScheme(ctx, s.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		return got
	}

	s2 := buildSuccessor(s1)
	s3 := buildSuccessor(s2)

	mp, err := mem.CreateMultiPeriodScheme(ctx, repro.MultiPeriodScheme{
		ID:    repro.NewMultiPeriodSchemeID(),
		Label: "three periods",
		SchemeIDs: []repro.ExtendedReproductionSchemeID{
			s1.ID, s2.ID, s3.ID,
		},
	})
	if err != nil {
		t.Fatalf("CreateMultiPeriodScheme: %v", err)
	}

	shifts, err := mem.ListCompositionShifts(ctx, mp.ID)
	if err != nil {
		t.Fatalf("ListCompositionShifts: %v", err)
	}
	if len(shifts) != 2 {
		t.Fatalf("expected 2 shifts for 3 schemes, got %d", len(shifts))
	}
	// Must be sorted by CycleNumber ascending.
	for i := 1; i < len(shifts); i++ {
		if shifts[i].CycleNumber <= shifts[i-1].CycleNumber {
			t.Fatalf("shifts not sorted: cycle[%d]=%d <= cycle[%d]=%d",
				i, shifts[i].CycleNumber, i-1, shifts[i-1].CycleNumber)
		}
	}
}
