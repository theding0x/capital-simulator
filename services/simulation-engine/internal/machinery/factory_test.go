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
	// 145,000 is Marx's (§8A); the 600_000/1,000 magnitudes are illustrative.
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
	// 145,000 is Marx's (§8A); the 600_000/1,000 magnitudes are illustrative.
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

// Capital Vol. I, Ch. 15, §2 — Mr. Baynes' self-acting mule. One horse-power
// drives 450 mule spindles tended by 2½ operatives; each spindle, "working ten
// hours, produces 13 ounces" of yarn, which a hand-spinner needs "sixty hours"
// (3,600 min) to spin. His text is ambiguous (ten hours per spindle, yet a weekly
// 365 5/8 lb for the 2½ operatives); the seed takes the ten-hour reading, so one
// tick is one day and one unit is one spindle's day (450/day), hence
// hand_labour_per_unit 3,600. The operatives' ten-hour day is 2½ × 600 = 1,500 min =
// 5 attendants × 300 min, half necessary (the project's 100% convention; Marx's 100%
// is a hypothetical in the plough passage). Machine value 480_000, life 2,000 days,
// the steam prime mover and the mill's name are illustrative: Marx gives none of them.
func baynesMuleMill(intensity machinery.IntensityFactor) machinery.Factory {
	return machinery.Factory{
		ID:         machinery.NewFactoryID(),
		Name:       "Baynes Mule Mill",
		PrimeMover: machinery.PrimeMover{Kind: machinery.PrimeMoverSteam, Horsepower: 1},
		Machines: []machinery.Machine{{
			ID:                machinery.NewMachineID(),
			MachineValue:      machinery.MachineValue(480_000),
			LifespanDays:      2000,
			ProductivePower:   450,
			HandLabourPerUnit: machinery.LabourMinutes(3600),
		}},
		IntensityFactor:             intensity,
		Attendants:                  5,
		LivingLabourPerAttendant:    machinery.LabourMinutes(300),
		NecessaryLabourPerAttendant: machinery.LabourMinutes(150),
	}
}

// Only living labour creates new value: W = c + v + s, with c from wear.
// c = 480,000 / 2,000 = 240; living labour 5 × 300 = 1,500 (= 2½ × 600);
// v = 5 × 150 = 750 and s = 750.
func TestFactory_TickPartitionsValueIntoCVS(t *testing.T) {
	t.Parallel()
	r := baynesMuleMill(0).RunTick()
	if r.ValueTransferred != 240 || r.LivingLabour != 1500 || r.VariableCapital != 750 || r.SurplusValue != 750 {
		t.Fatalf("c, living, v, s = %d, %d, %d, %d; want 240, 1500, 750, 750",
			r.ValueTransferred, r.LivingLabour, r.VariableCapital, r.SurplusValue)
	}
	if r.LivingLabour != r.VariableCapital+r.SurplusValue {
		t.Fatalf("LivingLabour %d != v + s %d", r.LivingLabour, r.VariableCapital+r.SurplusValue)
	}
}

// §2: the attendants' labour belongs in the cost of the machinery. Per day the
// mill costs 240 (wear) + 1,500 (living labour) = 1,740 against 450 × 3,600 =
// 1,620,000 hand-minutes displaced: a factor of about 931. That figure is a modelling
// consequence of reading one tick as a day (25 operative-hours against 27,000
// hand-hours), not Marx's: his two figures, 27,000 h and 150 h, give 180×. Without
// attendants it would be 1,620,000 / 240 = 6,750.
func TestFactoryProductivityFactor_CountsAttendantLabour(t *testing.T) {
	t.Parallel()
	f := baynesMuleMill(0)
	want := float64(450*3600) / float64(240+1500)
	if got := machinery.FactoryProductivityFactor(f); got != want {
		t.Fatalf("factor = %v, want %v", got, want)
	}
	f.Attendants = 0
	if got := machinery.FactoryProductivityFactor(f); got <= want {
		t.Fatalf("without attendants factor = %v, want > %v", got, want)
	}
	if got := machinery.FactoryProductivityFactor(f); got != 6750 {
		t.Fatalf("without attendants factor = %v, want 6750", got)
	}
}

// §3C: intensity densifies living labour, not the hand labour displaced.
// At 1.5× the 1,500 min of living labour become 2,250: v stays 750, s doubles
// to 1,500.
func TestFactory_IntensityScalesLivingLabourNotHandLabour(t *testing.T) {
	t.Parallel()
	base := baynesMuleMill(1.0).RunTick()
	dense := baynesMuleMill(machinery.NewIntensityFactor(1.5)).RunTick()
	if dense.HandLabourSaved != base.HandLabourSaved {
		t.Fatalf("HandLabourSaved changed with intensity: %d vs %d", dense.HandLabourSaved, base.HandLabourSaved)
	}
	if dense.LivingLabour != 2250 {
		t.Fatalf("LivingLabour at 1.5x = %d, want 2250", dense.LivingLabour)
	}
	if dense.VariableCapital != base.VariableCapital || dense.SurplusValue != 1500 {
		t.Fatalf("v, s at 1.5x = %d, %d; want 750, 1500", dense.VariableCapital, dense.SurplusValue)
	}
}

// §2 social limit: machinery costing more than the hand labour it replaces
// confers no productivity gain. Here 450 min saved against 1,740 min of cost.
func TestFactoryProductivityFactor_FlooredAtOne(t *testing.T) {
	t.Parallel()
	f := baynesMuleMill(0)
	f.Machines[0].HandLabourPerUnit = 1
	if got := machinery.FactoryProductivityFactor(f); got != 1.0 {
		t.Fatalf("factor = %v, want 1.0", got)
	}
}

func TestFactory_ValidateAttendantLabour(t *testing.T) {
	t.Parallel()
	f := baynesMuleMill(0)
	f.NecessaryLabourPerAttendant = f.LivingLabourPerAttendant + 1
	if err := f.Validate(); err != machinery.ErrFactoryLabour {
		t.Fatalf("necessary > living: err = %v, want ErrFactoryLabour", err)
	}
	f = baynesMuleMill(0)
	f.Attendants = -1
	if err := f.Validate(); err != machinery.ErrFactoryLabour {
		t.Fatalf("negative attendants: err = %v, want ErrFactoryLabour", err)
	}
}

// Capital Vol. I, Ch. 15, §2 — his two spinning figures give 180×.
// Marx: 366 lbs. of cotton absorb "only 150 hours' labour", where the same weight
// by hand-spinning would absorb "2,700 days' labour of ten hours each, or 27,000
// hours' labour".
//
// One tick (a day, on the ten-hour reading) of the Baynes mill spins 450 × 13 oz:
// 450 units, each displacing 3,600 hand-minutes, so 450 × 3,600 = 1,620,000 min =
// Marx's 27,000 hours.
func TestFactory_BaynesMule_Para2_HandLabourDisplaced(t *testing.T) {
	t.Parallel()
	f := baynesMuleMill(0)
	r := f.RunTick()
	if r.UnitsProduced != 450 {
		t.Fatalf("UnitsProduced = %d, want 450", r.UnitsProduced)
	}
	if r.HandLabourSaved != machinery.LabourMinutes(27_000*60) {
		t.Fatalf("HandLabourSaved = %d, want %d (27,000 h)", r.HandLabourSaved, 27_000*60)
	}
	if r.LivingLabour != 1500 {
		t.Fatalf("LivingLabour = %d, want 1500", r.LivingLabour)
	}
	// His two figures, 27,000 h and 150 h (a 60-hour week); not the per-day factor.
	if got := machinery.LabourSavingRatio(r.HandLabourSaved, machinery.LabourMinutes(150*60)); got != 180 {
		t.Fatalf("LabourSavingRatio = %d, want 180", got)
	}
	if err := f.Validate(); err != nil {
		t.Fatalf("Validate = %v, want nil", err)
	}
}
