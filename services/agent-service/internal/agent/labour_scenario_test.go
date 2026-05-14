package agent

import (
	"errors"
	"testing"
)

// Marx Ch. 17 §1: "A working day of 12 hours ... creates the same amount
// of value ... six shillings." At normal intensity the daily value
// equals the working-day minutes regardless of how it is partitioned.
func TestDailyValueCreated_NormalIntensityEqualsWorkingDay(t *testing.T) {
	t.Parallel()
	day := WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360} // 12 h
	got := DailyValueCreated(day, NormalLabourIntensity)
	if got != 720 {
		t.Fatalf("daily value at normal intensity: got %d, want 720", got)
	}
}

// Marx Ch. 17 §2: "a working-day of more intense labour is embodied in
// more products ... value created may be 7, 8, or more shillings." A
// factor of 4/3 (Marx's example of intensification by a third) produces
// 960 minutes of value out of a 720-minute day.
func TestDailyValueCreated_IntensifiedLabourScalesValue(t *testing.T) {
	t.Parallel()
	day := WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360}
	got := DailyValueCreated(day, LabourIntensity{Factor: 4.0 / 3.0})
	want := LabourMinutes(960)
	if got != want {
		t.Fatalf("daily value at 4/3 intensity: got %d, want %d", got, want)
	}
	if got <= DailyValueCreated(day, NormalLabourIntensity) {
		t.Fatal("intensified labour must produce more value than normal intensity")
	}
}

// Marx Ch. 17 §1: "If the value of the labour-power fall to 3 shillings,
// or the necessary labour to 6 hours, the surplus-value will rise to 3
// shillings." With a 12-hour day (720 min) and necessary = 360, surplus
// = 360. Rate of surplus-value = 100%.
func TestComputeOutcome_S1_PartitionAndRate(t *testing.T) {
	t.Parallel()
	s := LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360},
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	}
	out := ComputeOutcome(s)
	if out.DailyValueMinutes != 720 {
		t.Errorf("daily value: got %d, want 720", out.DailyValueMinutes)
	}
	if out.NecessaryLabourMinutes != 360 {
		t.Errorf("necessary: got %d, want 360", out.NecessaryLabourMinutes)
	}
	if out.SurplusLabourMinutes != 360 {
		t.Errorf("surplus: got %d, want 360", out.SurplusLabourMinutes)
	}
	if out.LabourPowerValueMinutes != 360 {
		t.Errorf("labour-power value: got %d, want 360", out.LabourPowerValueMinutes)
	}
	if out.RateOfSurplusValue != 1.0 {
		t.Errorf("rate of surplus-value: got %v, want 1.0", out.RateOfSurplusValue)
	}
}

// Marx Ch. 17 §3: "necessary labour time be 6 hours ... working-day be
// lengthened by 2 hours ... surplus-value increases both absolutely and
// relatively." Day = 14 h (840 min), necessary = 360 → surplus = 480.
func TestComputeOutcome_S3_WorkingDayLengthened(t *testing.T) {
	t.Parallel()
	s := LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 480}, // 14h total
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	}
	out := ComputeOutcome(s)
	if out.SurplusLabourMinutes != 480 {
		t.Errorf("surplus after lengthening: got %d, want 480", out.SurplusLabourMinutes)
	}
	if out.DailyValueMinutes != 840 {
		t.Errorf("daily value at 14h: got %d, want 840", out.DailyValueMinutes)
	}
}

// Marx Ch. 17 §4A: "working-day at 12 hours ... value of labour-power
// rises from 3 shillings to 4 ... if the day be lengthened by 4 hours ...
// absolute magnitude of surplus-value rises from 3 shillings to 4." Day
// = 16h (960 min), necessary = 480 → surplus = 480.
func TestComputeOutcome_S4A_LengthenedDayCompensatesRisingLPValue(t *testing.T) {
	t.Parallel()
	s := LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 480, SurplusLabourMinutes: 480},
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	}
	out := ComputeOutcome(s)
	if out.SurplusLabourMinutes != 480 {
		t.Errorf("surplus §4A: got %d, want 480", out.SurplusLabourMinutes)
	}
	// 33⅓% rise in absolute surplus (3 sh → 4 sh) — the chapter's claim.
	const oldSurplus = 360.0
	risePct := (float64(out.SurplusLabourMinutes) - oldSurplus) / oldSurplus
	if risePct < 0.33 || risePct > 0.34 {
		t.Errorf("surplus rise: got %.4f, want ~0.3333", risePct)
	}
}

// Marx Ch. 17 §1 Law 2: when productivity rises, value of labour-power
// falls and surplus rises — inverse relation.
func TestLawInverseRelation_ProductivityRise(t *testing.T) {
	t.Parallel()
	before := ScenarioOutcome{LabourPowerValueMinutes: 480, SurplusLabourMinutes: 240}
	after := ScenarioOutcome{LabourPowerValueMinutes: 360, SurplusLabourMinutes: 360}
	if !LawInverseRelation(before, after) {
		t.Fatal("LP value falls and surplus rises must satisfy LawInverseRelation")
	}
}

// Both moving in the same direction violates the §1 law 2.
func TestLawInverseRelation_SameDirectionFails(t *testing.T) {
	t.Parallel()
	before := ScenarioOutcome{LabourPowerValueMinutes: 360, SurplusLabourMinutes: 360}
	after := ScenarioOutcome{LabourPowerValueMinutes: 400, SurplusLabourMinutes: 400}
	if LawInverseRelation(before, after) {
		t.Fatal("both rising violates inverse relation")
	}
}

// §1 Law 1: at normal intensity, daily value equals working-day minutes
// regardless of productivity.
func TestLawConstantDailyValue_NormalIntensity(t *testing.T) {
	t.Parallel()
	day := WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360}
	if !LawConstantDailyValue(day, NormalLabourIntensity, 1.5) {
		t.Fatal("daily value must be constant under productivity change at normal intensity")
	}
	if !LawConstantDailyValue(day, NormalLabourIntensity, 0.5) {
		t.Fatal("daily value must be constant for productivity falls too")
	}
}

func TestLawConstantDailyValue_RequiresNormalIntensity(t *testing.T) {
	t.Parallel()
	day := WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360}
	if LawConstantDailyValue(day, LabourIntensity{Factor: 1.5}, 1.0) {
		t.Fatal("law only holds at normal intensity")
	}
}

func TestLabourScenario_Validate(t *testing.T) {
	t.Parallel()
	ok := LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360},
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	}
	if err := ok.Validate(); err != nil {
		t.Fatalf("valid scenario rejected: %v", err)
	}

	badIntensity := ok
	badIntensity.Intensity = LabourIntensity{Factor: 0}
	if err := badIntensity.Validate(); !errors.Is(err, ErrIntensityFactorNonPositive) {
		t.Errorf("zero intensity: got %v, want ErrIntensityFactorNonPositive", err)
	}

	badProductivity := ok
	badProductivity.Productivity = LabourProductivity{Factor: -1}
	if err := badProductivity.Validate(); !errors.Is(err, ErrProductivityFactorNonPositive) {
		t.Errorf("negative productivity: got %v, want ErrProductivityFactorNonPositive", err)
	}

	overLong := ok
	overLong.WorkingDay = WorkingDay{NecessaryLabourMinutes: 1500, SurplusLabourMinutes: 0} // > physical max
	if err := overLong.Validate(); !errors.Is(err, ErrWorkingDayExceedsPhysicalMax) {
		t.Errorf("over-long day: got %v, want ErrWorkingDayExceedsPhysicalMax", err)
	}
}
