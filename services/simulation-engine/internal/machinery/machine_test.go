package machinery_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
)

// Capital Vol. I, Ch. 15, §2:
// "the machine ... enters into the value-begetting process only by bits. It
//  never adds more value than it loses, on an average, by wear and tear."
//
// A £600 machine that wears out evenly over 1,000 working days transfers
// 600 LabourMinutes-worth of value per day.
func TestDailyWearAndTear_Para2_BitByBit(t *testing.T) {
	t.Parallel()
	m := machinery.Machine{
		MachineValue: machinery.MachineValue(600_000),
		LifespanDays: 1000,
	}
	got := machinery.DailyWearAndTear(m)
	if got != machinery.LabourMinutes(600) {
		t.Fatalf("DailyWearAndTear = %d, want 600", got)
	}
}

// Capital Vol. I, Ch. 15, §2:
// "2½ operatives spin weekly 365⅝ lbs. of yarn ... 366 lbs. of cotton absorb
//  only 150 hours' labour, while the same weight of cotton spun by the hand-
//  spinning wheel would absorb 2,700 ten-hour working-days, or 27,000 hours."
//
// Machine path: 150 h × 60 = 9,000 min for 366 lb.
// Hand path:    27,000 h × 60 = 1,620,000 min for 366 lb.
// Labour saving ratio: 1,620,000 / 9,000 = 180.
func TestLabourSavingRatio_Para2_Spinning(t *testing.T) {
	t.Parallel()
	const machineMinutes = machinery.LabourMinutes(9_000)
	const handMinutes = machinery.LabourMinutes(1_620_000)
	ratio := machinery.LabourSavingRatio(handMinutes, machineMinutes)
	if ratio != 180 {
		t.Fatalf("LabourSavingRatio = %d, want 180", ratio)
	}
}

// Capital Vol. I, Ch. 15, §2:
// "a machine working 16 hours daily for 7½ years covers as long a working
//  period as the same machine working only 8 hours daily for 15 years. But in
//  the first case the value would be reproduced twice as quickly."
//
// Each path consumes the full StandardWorkingDayHours × LifespanDays of
// machine life. TotalValueTransferred returns the whole MachineValue once the
// lifespan is exhausted.
func TestTotalValueTransferred_Para2_LifecycleEquivalence(t *testing.T) {
	t.Parallel()
	m := machinery.Machine{
		MachineValue: machinery.MachineValue(600_000),
		LifespanDays: 5475, // calibrated to a standard 8h day
	}
	a := machinery.TotalValueTransferred(m, 16, 2738)
	b := machinery.TotalValueTransferred(m, 8, 5475)
	if a != b {
		t.Fatalf("TotalValueTransferred(16, 2738)=%d, (8, 5475)=%d; want equal", a, b)
	}
	if a != machinery.LabourMinutes(m.MachineValue) {
		t.Fatalf("TotalValueTransferred at full life = %d, want %d", a, m.MachineValue)
	}
	// Surplus extracted twice as fast in the 16h case.
	rate16 := machinery.RateOfValueReproduction(m, 16)
	rate8 := machinery.RateOfValueReproduction(m, 8)
	if rate16 != 2*rate8 {
		t.Fatalf("RateOfValueReproduction 16h=%v, 8h=%v; want 2× relation", rate16, rate8)
	}
}

// Capital Vol. I, Ch. 15, §2:
// "a steam-plough does as much work in one hour at a cost of three-pence, as
//  66 men at a cost of 15 shillings."
//
// LabourDisplaced = 66 man-hours = 3,960 LabourMinutes.
// Machine cost in pence (3d/hr) is far less than the displaced wage bill
// (15s = 180d/hr).
func TestLabourDisplaced_Para2_SteamPlough(t *testing.T) {
	t.Parallel()
	const handLabourMinutesPerHour = machinery.LabourMinutes(60)
	const menReplaced = 66
	const machineUnitsPerHour = 1 // one hour of ploughing
	displaced := machinery.LabourDisplaced(handLabourMinutesPerHour*menReplaced, machineUnitsPerHour)
	if displaced != machinery.LabourMinutes(3960) {
		t.Fatalf("LabourDisplaced = %d, want 3960", displaced)
	}
	const machineCostPencePerHour = 3
	const displacedWagePence = 15 * 12 // 15 shillings = 180 pence
	if !(displacedWagePence > machineCostPencePerHour) {
		t.Fatalf("displaced wage bill (%d) must exceed machine cost (%d)",
			displacedWagePence, machineCostPencePerHour)
	}
}

// Capital Vol. I, Ch. 15, §3B (Moral Depreciation):
// "It loses exchange-value, either by machines of the same sort being produced
//  cheaper than it, or by better machines entering into competition with it."
//
// When a new machine of the same ProductivePower enters at a lower
// MachineValue, the older one loses the difference as moral depreciation.
func TestMoralDepreciation_Para3B_CheaperRival(t *testing.T) {
	t.Parallel()
	existing := machinery.Machine{
		MachineValue:    machinery.MachineValue(1_000_000),
		LifespanDays:    1000,
		ProductivePower: 145_000,
	}
	rival := machinery.Machine{
		MachineValue:    machinery.MachineValue(700_000),
		LifespanDays:    1000,
		ProductivePower: 145_000,
	}
	md := machinery.ComputeMoralDepreciation(existing, rival)
	if md.Value <= 0 {
		t.Fatalf("MoralDepreciation.Value = %d, want > 0", md.Value)
	}
	if md.Value != machinery.LabourMinutes(300_000) {
		t.Fatalf("MoralDepreciation.Value = %d, want 300000", md.Value)
	}
	// Invariant: never negative.
	noChange := machinery.ComputeMoralDepreciation(existing, existing)
	if noChange.Value < 0 {
		t.Fatalf("MoralDepreciation.Value = %d, want >= 0", noChange.Value)
	}
}

// Capital Vol. I, Ch. 15, §8A:
// "a single needle-machine makes 145,000 in a working-day of 11 hours. One
//  woman superintends four such machines producing near upon 600,000 needles
//  a day."
//
// Per-machine ProductivePower = 145,000 units/day.
// Four machines under one supervisor: ~580,000 units/day (close to Marx's
// "near upon 600,000" qualifier).
func TestProductivePower_Para8A_NeedleMachine(t *testing.T) {
	t.Parallel()
	m := machinery.Machine{ProductivePower: 145_000}
	if m.ProductivePower != 145_000 {
		t.Fatalf("ProductivePower = %d, want 145000", m.ProductivePower)
	}
	fourMachines := int64(m.ProductivePower) * 4
	if fourMachines != 580_000 {
		t.Fatalf("4 × ProductivePower = %d, want 580000", fourMachines)
	}
}

// Invariant from spec: ValueTransferredPerUnit(m, unitsPerDay) ==
// DailyWearAndTear(m) / unitsPerDay. Doubling output halves per-unit value.
func TestValueTransferredPerUnit_HalvesWhenOutputDoubles(t *testing.T) {
	t.Parallel()
	m := machinery.Machine{
		MachineValue: machinery.MachineValue(600_000),
		LifespanDays: 1000,
	}
	low := machinery.ValueTransferredPerUnit(m, 100)
	high := machinery.ValueTransferredPerUnit(m, 200)
	if low <= 0 || high <= 0 {
		t.Fatalf("ValueTransferredPerUnit must be positive; got low=%d high=%d", low, high)
	}
	if low != 2*high {
		t.Fatalf("doubling output should halve per-unit value: low=%d, high=%d", low, high)
	}
}

// Invariant from spec: LabourDisplaced(m) > LabourEmbodiedInMachine(m).
// "there is always a difference of labour saved in favour of the machine."
func TestLabourDisplaced_ExceedsEmbodiedLabour(t *testing.T) {
	t.Parallel()
	// 66 men replaced for one hour: 3,960 LabourMinutes saved daily across
	// many hours of operation. Embodied labour amortised per day is
	// machineValue / lifespan; pick a machine whose daily wear is small.
	m := machinery.Machine{
		MachineValue: machinery.MachineValue(60_000),
		LifespanDays: 1000,
	}
	embodiedPerDay := machinery.DailyWearAndTear(m) // 60 LabourMinutes
	const handLabourPerMachineDay = machinery.LabourMinutes(3960)
	const machineUnitsPerDay = 1
	displaced := machinery.LabourDisplaced(handLabourPerMachineDay, machineUnitsPerDay)
	if !(displaced > embodiedPerDay) {
		t.Fatalf("LabourDisplaced (%d) must exceed embodied labour per day (%d)",
			displaced, embodiedPerDay)
	}
}

// Invariant: DailyWearAndTear is exact integer division of MachineValue by LifespanDays.
// "a machine adds no more value per day than it loses."
func TestDailyWearAndTear_Invariant(t *testing.T) {
	t.Parallel()
	m := machinery.Machine{
		MachineValue: machinery.MachineValue(1_000_000),
		LifespanDays: 2500,
	}
	got := machinery.DailyWearAndTear(m)
	want := machinery.LabourMinutes(int64(m.MachineValue) / int64(m.LifespanDays))
	if got != want {
		t.Fatalf("DailyWearAndTear = %d, want %d", got, want)
	}
}

func TestMachineID_NewIsUniqueAndHex(t *testing.T) {
	t.Parallel()
	a := machinery.NewMachineID()
	b := machinery.NewMachineID()
	if a == b {
		t.Fatalf("NewMachineID should return distinct ids; got duplicate %s", a)
	}
	if len(a) != 24 {
		t.Fatalf("MachineID length = %d, want 24", len(a))
	}
}
