package store_test

import (
	"context"
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
