package machinery_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
)

// Capital Vol. I, Ch. 15, §1c — the Factory as a system of machines.
// "An organised system of machines, to which motion is communicated by the
//  transmitting mechanism from a central automaton, is the most developed
//  form of production by machinery."
func TestFactory_AggregateProductivePower(t *testing.T) {
	t.Parallel()
	a := machinery.Machine{
		ID:              machinery.NewMachineID(),
		ProductivePower: 145_000,
		MachineValue:    machinery.MachineValue(600_000),
		LifespanDays:    1000,
	}
	b := machinery.Machine{
		ID:              machinery.NewMachineID(),
		ProductivePower: 145_000,
		MachineValue:    machinery.MachineValue(600_000),
		LifespanDays:    1000,
	}
	f := machinery.Factory{
		ID:         machinery.NewFactoryID(),
		Name:       "Needle Works",
		PrimeMover: machinery.PrimeMover{Kind: machinery.PrimeMoverSteam, Horsepower: 30},
		Machines:   []machinery.Machine{a, b},
	}
	if got := f.TotalProductivePower(); got != 290_000 {
		t.Fatalf("TotalProductivePower = %d, want 290000", got)
	}
	if got := f.DailyValueTransfer(); got != machinery.LabourMinutes(1200) {
		t.Fatalf("DailyValueTransfer = %d, want 1200", got)
	}
}

func TestFactory_TickAccumulatesWearAndOutput(t *testing.T) {
	t.Parallel()
	m := machinery.Machine{
		ID:              machinery.NewMachineID(),
		ProductivePower: 145_000,
		MachineValue:    machinery.MachineValue(600_000),
		LifespanDays:    1000,
	}
	f := machinery.Factory{
		ID:         machinery.NewFactoryID(),
		PrimeMover: machinery.PrimeMover{Kind: machinery.PrimeMoverSteam, Horsepower: 30},
		Machines:   []machinery.Machine{m},
	}
	result := f.RunTick()
	if result.ValueTransferred != machinery.LabourMinutes(600) {
		t.Fatalf("tick ValueTransferred = %d, want 600", result.ValueTransferred)
	}
	if result.UnitsProduced != 145_000 {
		t.Fatalf("tick UnitsProduced = %d, want 145000", result.UnitsProduced)
	}
}

func TestFactory_Validate(t *testing.T) {
	t.Parallel()
	bad := machinery.Factory{}
	if err := bad.Validate(); err == nil {
		t.Fatal("Validate must reject empty factory")
	}
	noMover := machinery.Factory{Name: "x", Machines: []machinery.Machine{{}}}
	if err := noMover.Validate(); err == nil {
		t.Fatal("Validate must reject factory with no prime mover kind")
	}
}

// Capital Vol. I, Ch. 15, §2 — the power-loom. One weaver tends the loom
// (hand_labour_per_unit 4400, 366 units/day, 480000 over 2000 days).
func powerLoomShed(intensity machinery.IntensityFactor) machinery.Factory {
	return machinery.Factory{
		ID:         machinery.NewFactoryID(),
		Name:       "Manchester Weaving Shed",
		PrimeMover: machinery.PrimeMover{Kind: machinery.PrimeMoverSteam, Horsepower: 60},
		Machines: []machinery.Machine{{
			ID:                machinery.NewMachineID(),
			MachineValue:      machinery.MachineValue(480_000),
			LifespanDays:      2000,
			ProductivePower:   366,
			HandLabourPerUnit: machinery.LabourMinutes(4400),
		}},
		IntensityFactor:             intensity,
		Attendants:                  1,
		LivingLabourPerAttendant:    machinery.LabourMinutes(480),
		NecessaryLabourPerAttendant: machinery.LabourMinutes(240),
	}
}

// Only living labour creates new value: W = c + v + s, with c from wear.
func TestFactory_TickPartitionsValueIntoCVS(t *testing.T) {
	t.Parallel()
	r := powerLoomShed(0).RunTick()
	if r.ValueTransferred != 240 || r.VariableCapital != 240 || r.SurplusValue != 240 {
		t.Fatalf("c, v, s = %d, %d, %d; want 240, 240, 240", r.ValueTransferred, r.VariableCapital, r.SurplusValue)
	}
	if r.LivingLabour != r.VariableCapital+r.SurplusValue {
		t.Fatalf("LivingLabour %d != v + s %d", r.LivingLabour, r.VariableCapital+r.SurplusValue)
	}
}

// §2: the attendants' labour belongs in the cost of the machinery.
func TestFactoryProductivityFactor_CountsAttendantLabour(t *testing.T) {
	t.Parallel()
	f := powerLoomShed(0)
	want := float64(366*4400) / float64(240+480)
	if got := machinery.FactoryProductivityFactor(f); got != want {
		t.Fatalf("factor = %v, want %v", got, want)
	}
	f.Attendants = 0
	if got := machinery.FactoryProductivityFactor(f); got <= want {
		t.Fatalf("without attendants factor = %v, want > %v", got, want)
	}
}

// §3C: intensity densifies living labour, not the hand labour displaced.
func TestFactory_IntensityScalesLivingLabourNotHandLabour(t *testing.T) {
	t.Parallel()
	base := powerLoomShed(1.0).RunTick()
	dense := powerLoomShed(machinery.NewIntensityFactor(1.5)).RunTick()
	if dense.HandLabourSaved != base.HandLabourSaved {
		t.Fatalf("HandLabourSaved changed with intensity: %d vs %d", dense.HandLabourSaved, base.HandLabourSaved)
	}
	if dense.LivingLabour != 720 {
		t.Fatalf("LivingLabour at 1.5x = %d, want 720", dense.LivingLabour)
	}
	if dense.VariableCapital != base.VariableCapital || dense.SurplusValue != 480 {
		t.Fatalf("v, s at 1.5x = %d, %d; want 240, 480", dense.VariableCapital, dense.SurplusValue)
	}
}

// §2 social limit: machinery costing more than the hand labour it replaces
// confers no productivity gain.
func TestFactoryProductivityFactor_FlooredAtOne(t *testing.T) {
	t.Parallel()
	f := powerLoomShed(0)
	f.Machines[0].HandLabourPerUnit = 1
	if got := machinery.FactoryProductivityFactor(f); got != 1.0 {
		t.Fatalf("factor = %v, want 1.0", got)
	}
}

func TestFactory_ValidateAttendantLabour(t *testing.T) {
	t.Parallel()
	f := powerLoomShed(0)
	f.NecessaryLabourPerAttendant = f.LivingLabourPerAttendant + 1
	if err := f.Validate(); err != machinery.ErrFactoryLabour {
		t.Fatalf("necessary > living: err = %v, want ErrFactoryLabour", err)
	}
	f = powerLoomShed(0)
	f.Attendants = -1
	if err := f.Validate(); err != machinery.ErrFactoryLabour {
		t.Fatalf("negative attendants: err = %v, want ErrFactoryLabour", err)
	}
}
