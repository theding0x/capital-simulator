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
