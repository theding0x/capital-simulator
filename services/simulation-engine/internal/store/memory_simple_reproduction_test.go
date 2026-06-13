package store_test

import (
	"context"
	"errors"
	"testing"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// newCh20Store returns a fresh in-memory store for Ch. 20 tests.
func newCh20Store(t *testing.T) *store.Memory {
	t.Helper()
	return store.NewMemory()
}

// buildMarxScheme creates and returns a fully-populated balanced scheme using
// Marx's canonical Ch. 20 fixture: Dept I 4000c+1000v+1000s | Dept II 2000c+500v+500s.
func buildMarxScheme(t *testing.T, mem *store.Memory) repro.SimpleReproductionScheme {
	t.Helper()
	ctx := context.Background()

	s, err := mem.CreateSimpleReproductionScheme(ctx, repro.SimpleReproductionScheme{
		ID:     repro.NewSimpleReproductionSchemeID(),
		Period: "1865",
	})
	if err != nil {
		t.Fatalf("CreateSimpleReproductionScheme: %v", err)
	}

	dI := repro.DepartmentalCapital{
		Department:    repro.DepartmentI,
		ConstantPence: 4_000_000,
		VariablePence: 1_000_000,
		SurplusPence:  1_000_000,
		TotalPence:    6_000_000,
	}
	if _, err := mem.AddDepartmentToScheme(ctx, s.ID, dI); err != nil {
		t.Fatalf("AddDepartmentToScheme(I): %v", err)
	}

	dII := repro.DepartmentalCapital{
		Department:    repro.DepartmentII,
		ConstantPence: 2_000_000,
		VariablePence: 500_000,
		SurplusPence:  500_000,
		TotalPence:    3_000_000,
	}
	if _, err := mem.AddDepartmentToScheme(ctx, s.ID, dII); err != nil {
		t.Fatalf("AddDepartmentToScheme(II): %v", err)
	}

	got, err := mem.GetSimpleReproductionScheme(ctx, s.ID)
	if err != nil {
		t.Fatalf("GetSimpleReproductionScheme: %v", err)
	}
	return got
}

// --- CreateSimpleReproductionScheme -----------------------------------------

func TestMemoryCreateGetSimpleReproductionScheme(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	ctx := context.Background()

	s, err := mem.CreateSimpleReproductionScheme(ctx, repro.SimpleReproductionScheme{
		ID:     repro.NewSimpleReproductionSchemeID(),
		Period: "1865",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if s.ID == "" {
		t.Fatal("expected non-empty ID")
	}

	got, err := mem.GetSimpleReproductionScheme(ctx, s.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Period != "1865" {
		t.Fatalf("Period: want 1865, got %q", got.Period)
	}
	if got.IsBalanced {
		t.Fatal("expected IsBalanced=false before departments are added")
	}

	// Missing ID must return ErrNotFound.
	_, err = mem.GetSimpleReproductionScheme(ctx, "nonexistent")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryCreateSchemeErrAlreadyExists(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	ctx := context.Background()

	id := repro.NewSimpleReproductionSchemeID()
	if _, err := mem.CreateSimpleReproductionScheme(ctx, repro.SimpleReproductionScheme{
		ID: id, Period: "1865",
	}); err != nil {
		t.Fatalf("first create: %v", err)
	}
	_, err := mem.CreateSimpleReproductionScheme(ctx, repro.SimpleReproductionScheme{
		ID: id, Period: "1865",
	})
	if !errors.Is(err, store.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

// --- AddDepartmentToScheme --------------------------------------------------

func TestMemoryAddDepartmentKeystone(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	s := buildMarxScheme(t, mem)

	if !s.IsBalanced {
		t.Fatal("expected IsBalanced=true for canonical Marx fixture")
	}
	if s.DepartmentI == nil || s.DepartmentII == nil {
		t.Fatal("expected both departments present")
	}
	if s.DepartmentI.VplusSPence() != s.DepartmentII.ConstantPence {
		t.Fatalf("keystone: I(v+s)=%d != II(c)=%d",
			s.DepartmentI.VplusSPence(), s.DepartmentII.ConstantPence)
	}
}

func TestMemoryAddDepartmentAlreadySet(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	ctx := context.Background()

	s, _ := mem.CreateSimpleReproductionScheme(ctx, repro.SimpleReproductionScheme{
		ID: repro.NewSimpleReproductionSchemeID(), Period: "1865",
	})
	dI := repro.DepartmentalCapital{
		Department:    repro.DepartmentI,
		ConstantPence: 4_000_000, VariablePence: 1_000_000,
		SurplusPence: 1_000_000, TotalPence: 6_000_000,
	}
	if _, err := mem.AddDepartmentToScheme(ctx, s.ID, dI); err != nil {
		t.Fatalf("first add: %v", err)
	}
	_, err := mem.AddDepartmentToScheme(ctx, s.ID, dI)
	if !errors.Is(err, repro.ErrDepartmentAlreadySet) {
		t.Fatalf("expected ErrDepartmentAlreadySet, got %v", err)
	}
}

func TestMemoryAddDepartmentSchemeNotFound(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	ctx := context.Background()

	dI := repro.DepartmentalCapital{
		Department:    repro.DepartmentI,
		ConstantPence: 4_000_000, VariablePence: 1_000_000,
		SurplusPence: 1_000_000, TotalPence: 6_000_000,
	}
	_, err := mem.AddDepartmentToScheme(ctx, "ghost-scheme", dI)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// --- ListSimpleReproductionSchemes ------------------------------------------

func TestMemoryListSimpleReproductionSchemes(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	ctx := context.Background()

	list, err := mem.ListSimpleReproductionSchemes(ctx, "")
	if err != nil {
		t.Fatalf("list empty: %v", err)
	}
	if list == nil {
		t.Fatal("expected non-nil slice")
	}

	_ = buildMarxScheme(t, mem)
	list, err = mem.ListSimpleReproductionSchemes(ctx, "")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 scheme, got %d", len(list))
	}

	filtered, _ := mem.ListSimpleReproductionSchemes(ctx, "1871")
	if len(filtered) != 0 {
		t.Fatalf("expected 0 schemes for period 1871, got %d", len(filtered))
	}
}

// --- RecordInterDepartmentExchange / ListInterDepartmentExchanges -----------

func TestMemoryRecordListExchanges(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	ctx := context.Background()

	s := buildMarxScheme(t, mem)
	ex := repro.InterDepartmentExchange{
		FromDepartment: repro.DepartmentI,
		ToDepartment:   repro.DepartmentII,
		Pence:          2_000_000,
		Kind:           repro.ExchangeGoodsForMoney,
		Description:    "I(v+s) → II",
	}
	created, err := mem.RecordInterDepartmentExchange(ctx, s.ID, ex)
	if err != nil {
		t.Fatalf("record exchange: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty exchange ID")
	}

	list, err := mem.ListInterDepartmentExchanges(ctx, s.ID)
	if err != nil {
		t.Fatalf("list exchanges: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 exchange, got %d", len(list))
	}
	if list[0].Pence != 2_000_000 {
		t.Fatalf("Pence: want 2000000, got %d", list[0].Pence)
	}
}

// --- AdvanceReproductionTick ------------------------------------------------

func TestMemoryAdvanceReproductionTick(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	ctx := context.Background()

	s := buildMarxScheme(t, mem)

	tick, err := mem.AdvanceReproductionTick(ctx, s.ID, 0, 0)
	if err != nil {
		t.Fatalf("advance tick: %v", err)
	}
	if tick.TickNumber != 1 {
		t.Fatalf("tick number: want 1, got %d", tick.TickNumber)
	}
	if !tick.IsBalanced {
		t.Fatal("expected balanced tick for canonical Marx fixture")
	}

	tick2, err := mem.AdvanceReproductionTick(ctx, s.ID, 0, 0)
	if err != nil {
		t.Fatalf("advance tick 2: %v", err)
	}
	if tick2.TickNumber != 2 {
		t.Fatalf("tick number 2: want 2, got %d", tick2.TickNumber)
	}

	_, err = mem.AdvanceReproductionTick(ctx, "ghost-id", 0, 0)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for ghost scheme, got %v", err)
	}
}

// --- CheckSchemeBalance -----------------------------------------------------

func TestMemoryCheckSchemeBalance(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	ctx := context.Background()

	s := buildMarxScheme(t, mem)
	res, err := mem.CheckSchemeBalance(ctx, s.ID)
	if err != nil {
		t.Fatalf("check balance: %v", err)
	}
	if !res.IsBalanced {
		t.Fatal("expected IsBalanced=true")
	}
	if res.DeficitPence != 0 {
		t.Fatalf("DeficitPence: want 0, got %d", res.DeficitPence)
	}
	if res.DeptIVplusSPence != 2_000_000 {
		t.Fatalf("DeptIVplusSPence: want 2000000, got %d", res.DeptIVplusSPence)
	}
	if res.DeptIIConstantPence != 2_000_000 {
		t.Fatalf("DeptIIConstantPence: want 2000000, got %d", res.DeptIIConstantPence)
	}

	_, err = mem.CheckSchemeBalance(ctx, "ghost-id")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// --- AdvanceReproductionTick: stock-constancy (FIX 1) -----------------------

// TestMemoryAdvanceReproductionTickStockConstancy verifies the defining law of
// simple reproduction: advancing a tick must NOT change departmental c, v, s
// magnitudes. The scheme must be identical period-to-period (no accumulation).
func TestMemoryAdvanceReproductionTickStockConstancy(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	ctx := context.Background()

	s := buildMarxScheme(t, mem)

	// Capture the pre-tick departmental magnitudes.
	wantIC := s.DepartmentI.ConstantPence
	wantIV := s.DepartmentI.VariablePence
	wantIS := s.DepartmentI.SurplusPence
	wantIIT := s.DepartmentI.TotalPence
	wantIIC := s.DepartmentII.ConstantPence
	wantIIV := s.DepartmentII.VariablePence
	wantIIS := s.DepartmentII.SurplusPence
	wantIIIT := s.DepartmentII.TotalPence

	// Advance two ticks.
	if _, err := mem.AdvanceReproductionTick(ctx, s.ID, 0, 0); err != nil {
		t.Fatalf("advance tick 1: %v", err)
	}
	if _, err := mem.AdvanceReproductionTick(ctx, s.ID, 0, 0); err != nil {
		t.Fatalf("advance tick 2: %v", err)
	}

	// Re-read the scheme and assert all magnitudes are unchanged.
	got, err := mem.GetSimpleReproductionScheme(ctx, s.ID)
	if err != nil {
		t.Fatalf("GetSimpleReproductionScheme: %v", err)
	}
	if got.DepartmentI == nil || got.DepartmentII == nil {
		t.Fatal("departments must still be present after ticking")
	}

	if got.DepartmentI.ConstantPence != wantIC {
		t.Errorf("DeptI.ConstantPence: want %d, got %d (tick must not mutate c)", wantIC, got.DepartmentI.ConstantPence)
	}
	if got.DepartmentI.VariablePence != wantIV {
		t.Errorf("DeptI.VariablePence: want %d, got %d (tick must not mutate v)", wantIV, got.DepartmentI.VariablePence)
	}
	if got.DepartmentI.SurplusPence != wantIS {
		t.Errorf("DeptI.SurplusPence: want %d, got %d (tick must not mutate s)", wantIS, got.DepartmentI.SurplusPence)
	}
	if got.DepartmentI.TotalPence != wantIIT {
		t.Errorf("DeptI.TotalPence: want %d, got %d", wantIIT, got.DepartmentI.TotalPence)
	}
	if got.DepartmentII.ConstantPence != wantIIC {
		t.Errorf("DeptII.ConstantPence: want %d, got %d (tick must not mutate c)", wantIIC, got.DepartmentII.ConstantPence)
	}
	if got.DepartmentII.VariablePence != wantIIV {
		t.Errorf("DeptII.VariablePence: want %d, got %d (tick must not mutate v)", wantIIV, got.DepartmentII.VariablePence)
	}
	if got.DepartmentII.SurplusPence != wantIIS {
		t.Errorf("DeptII.SurplusPence: want %d, got %d (tick must not mutate s)", wantIIS, got.DepartmentII.SurplusPence)
	}
	if got.DepartmentII.TotalPence != wantIIIT {
		t.Errorf("DeptII.TotalPence: want %d, got %d", wantIIIT, got.DepartmentII.TotalPence)
	}
}

// --- RecordMoneyClosedLoop --------------------------------------------------

func TestMemoryRecordMoneyClosedLoop(t *testing.T) {
	t.Parallel()
	mem := newCh20Store(t)
	ctx := context.Background()

	s := buildMarxScheme(t, mem)
	loop := repro.MoneyClosedLoop{
		Period:       "1865",
		IsClosed:     true,
		NetFlowPence: 0,
	}
	created, err := mem.RecordMoneyClosedLoop(ctx, s.ID, loop)
	if err != nil {
		t.Fatalf("record loop: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if !created.IsClosed {
		t.Fatal("expected IsClosed=true")
	}

	_, err = mem.RecordMoneyClosedLoop(ctx, "ghost-id", loop)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for ghost scheme, got %v", err)
	}
}
