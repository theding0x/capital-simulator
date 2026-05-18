package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestMemory_CreateAndGetMachine(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	mc, err := st.CreateMachine(context.Background(), machinery.Machine{
		Name:            "needle-machine",
		MachineValue:    machinery.MachineValue(600_000),
		LifespanDays:    1000,
		ProductivePower: 145_000,
	})
	if err != nil {
		t.Fatalf("CreateMachine: %v", err)
	}
	if mc.ID.IsZero() {
		t.Fatal("CreateMachine should assign an ID")
	}
	got, err := st.GetMachine(context.Background(), mc.ID)
	if err != nil {
		t.Fatalf("GetMachine: %v", err)
	}
	if got.ID != mc.ID {
		t.Fatalf("GetMachine id=%s, want %s", got.ID, mc.ID)
	}
}

func TestMemory_AdvanceTick_AccumulatesWear(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st := store.NewMemory()
	mc, err := st.CreateMachine(ctx, machinery.Machine{
		Name:            "needle-machine",
		MachineValue:    machinery.MachineValue(600_000),
		LifespanDays:    1000,
		ProductivePower: 145_000,
	})
	if err != nil {
		t.Fatalf("CreateMachine: %v", err)
	}
	f, err := st.CreateFactory(ctx, machinery.Factory{
		Name:       "Needle Works",
		PrimeMover: machinery.PrimeMover{Kind: machinery.PrimeMoverSteam, Horsepower: 30},
		Machines:   []machinery.Machine{mc},
	})
	if err != nil {
		t.Fatalf("CreateFactory: %v", err)
	}
	advanced, tick, err := st.AdvanceTick(ctx, f.ID)
	if err != nil {
		t.Fatalf("AdvanceTick: %v", err)
	}
	if advanced.TickCount != 1 {
		t.Fatalf("TickCount = %d, want 1", advanced.TickCount)
	}
	if tick.ValueTransferred != 600 {
		t.Fatalf("ValueTransferred = %d, want 600", tick.ValueTransferred)
	}
	if tick.UnitsProduced != 145_000 {
		t.Fatalf("UnitsProduced = %d, want 145000", tick.UnitsProduced)
	}
	if tick.Sequence != 1 {
		t.Fatalf("Sequence = %d, want 1", tick.Sequence)
	}
	hist, err := st.ListTicks(ctx, f.ID, 10)
	if err != nil {
		t.Fatalf("ListTicks: %v", err)
	}
	if len(hist) != 1 {
		t.Fatalf("ListTicks len = %d, want 1", len(hist))
	}
	if hist[0].Sequence != 1 {
		t.Fatalf("ListTicks[0].Sequence = %d, want 1", hist[0].Sequence)
	}
	updated, err := st.GetMachine(ctx, mc.ID)
	if err != nil {
		t.Fatalf("GetMachine: %v", err)
	}
	if updated.AccumulatedWear.Value != machinery.LabourMinutes(600) {
		t.Fatalf("AccumulatedWear = %d, want 600", updated.AccumulatedWear.Value)
	}
}

func TestMemory_GetMissingReturnsNotFound(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	if _, err := st.GetMachine(context.Background(), machinery.MachineID("nope")); err != store.ErrNotFound {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if _, err := st.GetFactory(context.Background(), machinery.FactoryID("nope")); err != store.ErrNotFound {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if _, err := st.GetGeneralLawScenario(context.Background(), simulation.GeneralLawScenarioID("nope")); err != store.ErrNotFound {
		t.Fatalf("GetGeneralLawScenario missing: err = %v, want ErrNotFound", err)
	}
}

func TestMemory_GeneralLaw_CreateAndGet(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	ctx := context.Background()

	s := simulation.GeneralLawScenario{
		Name:               "§1 unchanged OC",
		ConstantCapital:    simulation.Pence(8000),
		VariableCapital:    simulation.Pence(2000),
		SurplusRate:        1.0,
		AccumulationRate:   1.0,
		ProductivityGrowth: 0.0,
		WagePence:          200,
		WorkerSupply:       50,
		Periods:            3,
	}
	created, err := st.CreateGeneralLawScenario(ctx, s)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if string(created.ID) == "" {
		t.Fatal("id should be non-empty")
	}

	got, err := st.GetGeneralLawScenario(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("id = %q, want %q", got.ID, created.ID)
	}
	if got.Name != s.Name {
		t.Errorf("name = %q, want %q", got.Name, s.Name)
	}
	if got.ConstantCapital != s.ConstantCapital {
		t.Errorf("constant_capital = %d, want %d", got.ConstantCapital, s.ConstantCapital)
	}
}

func TestMemory_ColonialLabourMarket_CRUD(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st := store.NewMemory()

	market := simulation.ColonialLabourMarket{
		Colony:          "Swan River, 1829",
		FreeLabourers:   300,
		AnnualWagePence: 4800,
		LandAvailable:   true,
	}
	created, err := st.CreateColonialLabourMarket(ctx, market)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("ID should be assigned")
	}
	if created.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be assigned")
	}

	got, err := st.GetColonialLabourMarket(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Colony != market.Colony {
		t.Errorf("Colony = %q, want %q", got.Colony, market.Colony)
	}
	if got.FreeLabourers != market.FreeLabourers {
		t.Errorf("FreeLabourers = %d, want %d", got.FreeLabourers, market.FreeLabourers)
	}

	all, err := st.ListColonialLabourMarkets(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("list len = %d, want 1", len(all))
	}
}

func TestMemory_ColonialLabourMarket_DuplicateColonyRejected(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st := store.NewMemory()

	first := simulation.ColonialLabourMarket{
		Colony:          "Saint Domingo",
		FreeLabourers:   50,
		AnnualWagePence: 3600,
	}
	if _, err := st.CreateColonialLabourMarket(ctx, first); err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := simulation.ColonialLabourMarket{
		Colony:          "saint domingo", // case-insensitive uniqueness
		FreeLabourers:   100,
		AnnualWagePence: 3600,
	}
	_, err := st.CreateColonialLabourMarket(ctx, dup)
	if !errors.Is(err, store.ErrAlreadyExists) {
		t.Fatalf("second create error = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ColonialLabourMarket_GetMissingReturnsNotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st := store.NewMemory()
	_, err := st.GetColonialLabourMarket(ctx, "nonexistent")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_ColonialLabourMarket_UpdateAppliesPartialFields(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st := store.NewMemory()

	created, err := st.CreateColonialLabourMarket(ctx, simulation.ColonialLabourMarket{
		Colony:          "South Australia, 1836",
		FreeLabourers:   5000,
		AnnualWagePence: 4800,
		LandAvailable:   true,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	applied := true
	years := int64(5)
	extractable := true
	updated, err := st.UpdateColonialLabourMarket(ctx, created.ID, store.ColonialLabourMarketUpdate{
		WakefieldSchemeApplied:   &applied,
		IndependenceYears:        &years,
		SurplusLabourExtractable: &extractable,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if !updated.WakefieldSchemeApplied {
		t.Errorf("WakefieldSchemeApplied = false, want true")
	}
	if updated.IndependenceYears != 5 {
		t.Errorf("IndependenceYears = %d, want 5", updated.IndependenceYears)
	}
	if !updated.SurplusLabourExtractable {
		t.Errorf("SurplusLabourExtractable = false, want true")
	}
	if updated.FreeLabourers != created.FreeLabourers {
		t.Errorf("FreeLabourers changed: %d -> %d", created.FreeLabourers, updated.FreeLabourers)
	}
}

func TestMemory_RegulateColonialLabourMarket_AtomicReadComputeWrite(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st := store.NewMemory()

	// South Australia, 1836: £20/year wage; savings 2400/yr.
	// Ransom 240 * 50 = 12000 -> 5 years.
	created, err := st.CreateColonialLabourMarket(ctx, simulation.ColonialLabourMarket{
		Colony:          "South Australia, 1836",
		FreeLabourers:   5000,
		AnnualWagePence: 4800,
		LandAvailable:   true,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	regulated, err := st.RegulateColonialLabourMarket(ctx, created.ID, simulation.SystematicColonisation{
		SufficientPrice: simulation.SufficientPrice{
			PricePerAcrePence: 240,
			DesiredAcres:      50,
		},
	})
	if err != nil {
		t.Fatalf("regulate: %v", err)
	}
	if !regulated.WakefieldSchemeApplied {
		t.Errorf("WakefieldSchemeApplied = false, want true")
	}
	if regulated.IndependenceYears != 5 {
		t.Errorf("IndependenceYears = %d, want 5", regulated.IndependenceYears)
	}
	if !regulated.SurplusLabourExtractable {
		t.Errorf("SurplusLabourExtractable = false, want true")
	}

	// Re-read the persisted market: the regulated state must be the
	// state the store now returns from Get (the write committed).
	persisted, err := st.GetColonialLabourMarket(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !persisted.WakefieldSchemeApplied || persisted.IndependenceYears != 5 || !persisted.SurplusLabourExtractable {
		t.Errorf("persisted state not regulated: %+v", persisted)
	}

	// Zero sufficient price is no scheme: the market is returned
	// unchanged. The fact that the previous /regulate already set
	// WakefieldSchemeApplied=true makes this case observable.
	noop, err := st.RegulateColonialLabourMarket(ctx, created.ID, simulation.SystematicColonisation{})
	if err != nil {
		t.Fatalf("regulate (empty scheme): %v", err)
	}
	if noop.WakefieldSchemeApplied != persisted.WakefieldSchemeApplied ||
		noop.IndependenceYears != persisted.IndependenceYears ||
		noop.SurplusLabourExtractable != persisted.SurplusLabourExtractable {
		t.Errorf("empty-scheme regulate changed state: before=%+v after=%+v", persisted, noop)
	}
}

func TestMemory_RegulateColonialLabourMarket_NotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st := store.NewMemory()
	_, err := st.RegulateColonialLabourMarket(ctx, "missing", simulation.SystematicColonisation{
		SufficientPrice: simulation.SufficientPrice{PricePerAcrePence: 240, DesiredAcres: 50},
	})
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
