package agent

import (
	"errors"
	"math"
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

// Marx Ch. 17 §1 Case I, derived from the domain: a 12-hour day (720 min)
// with necessary labour 480 / surplus 240 (LPV 4s, s 2s) and productivity
// 4/3 yields necessary 360 / surplus 360 (LPV 3s, s 3s). The proportional
// change is asymmetric — LPV falls 25%, surplus-value rises 50% — which is
// Marx's correction of Ricardo.
func TestComputeOutcome_ProductivityRise_MarxCaseI(t *testing.T) {
	t.Parallel()
	day := WorkingDay{NecessaryLabourMinutes: 480, SurplusLabourMinutes: 240}
	before := ComputeOutcome(LabourScenario{WorkingDay: day, Intensity: NormalLabourIntensity, Productivity: NormalLabourProductivity})
	after := ComputeOutcome(LabourScenario{WorkingDay: day, Intensity: NormalLabourIntensity, Productivity: LabourProductivity{Factor: 4.0 / 3.0}})

	if after.NecessaryLabourMinutes != 360 || after.SurplusLabourMinutes != 360 {
		t.Fatalf("partition: got %d/%d, want 360/360", after.NecessaryLabourMinutes, after.SurplusLabourMinutes)
	}
	if after.LabourPowerValueMinutes != 360 {
		t.Fatalf("LPV: got %d, want 360", after.LabourPowerValueMinutes)
	}
	if after.DailyValueMinutes != before.DailyValueMinutes {
		t.Fatalf("daily value must stay constant: %d vs %d", after.DailyValueMinutes, before.DailyValueMinutes)
	}
	if !LawInverseRelation(before, after) {
		t.Fatal("derived outcomes must satisfy LawInverseRelation")
	}
	lpvFall := float64(before.LabourPowerValueMinutes-after.LabourPowerValueMinutes) / float64(before.LabourPowerValueMinutes)
	svRise := float64(after.SurplusLabourMinutes-before.SurplusLabourMinutes) / float64(before.SurplusLabourMinutes)
	if math.Abs(lpvFall-0.25) > 1e-9 || math.Abs(svRise-0.50) > 1e-9 {
		t.Fatalf("asymmetry: LPV fall %.4f (want 0.25), surplus rise %.4f (want 0.50)", lpvFall, svRise)
	}
}

// A productivity fall is the mirror case: necessary labour lengthens and
// surplus labour shrinks within the same day.
func TestComputeOutcome_ProductivityFall(t *testing.T) {
	t.Parallel()
	day := WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360}
	out := ComputeOutcome(LabourScenario{WorkingDay: day, Intensity: NormalLabourIntensity, Productivity: LabourProductivity{Factor: 0.75}})
	if out.NecessaryLabourMinutes != 480 || out.SurplusLabourMinutes != 240 {
		t.Fatalf("partition: got %d/%d, want 480/240", out.NecessaryLabourMinutes, out.SurplusLabourMinutes)
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

// Ch. 17 §2 Case II: a 12-hour day yielding 6s at normal intensity yields 8s
// at 4/3 intensity, and the two parts rise together from 3s to 4s.
func TestComputeOutcome_IntensityScalesBothPartsOfDailyValue(t *testing.T) {
	t.Parallel()
	out := ComputeOutcome(LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360},
		Intensity:    LabourIntensity{Factor: 4.0 / 3.0},
		Productivity: NormalLabourProductivity,
	})
	if out.DailyValueMinutes != 960 {
		t.Errorf("daily value: got %d, want 960 (8s)", out.DailyValueMinutes)
	}
	if out.LabourPowerValueMinutes != 480 || out.SurplusLabourMinutes != 480 {
		t.Errorf("parts: got %d/%d, want 480/480 (4s/4s)", out.LabourPowerValueMinutes, out.SurplusLabourMinutes)
	}
}

func TestComputeOutcome_PartitionHoldsAtEveryIntensity(t *testing.T) {
	t.Parallel()
	for _, f := range []float64{0.5, 1.0, 1.1, 4.0 / 3.0, 1.37, 2.0} {
		out := ComputeOutcome(LabourScenario{
			WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360},
			Intensity:    LabourIntensity{Factor: f},
			Productivity: NormalLabourProductivity,
		})
		if out.SurplusLabourMinutes+out.LabourPowerValueMinutes != out.DailyValueMinutes {
			t.Errorf("factor %v: %d + %d != %d", f, out.SurplusLabourMinutes, out.LabourPowerValueMinutes, out.DailyValueMinutes)
		}
	}
}

func TestDailyValueCreated_RoundsHalfUp(t *testing.T) {
	t.Parallel()
	day := WorkingDay{NecessaryLabourMinutes: 1, SurplusLabourMinutes: 2} // 3 min × 1.5 = 4.5
	if got := DailyValueCreated(day, LabourIntensity{Factor: 1.5}); got != 5 {
		t.Errorf("got %d, want 5", got)
	}
}

// Marx Ch. 17 §3: lengthening the day from 12h to 14h with necessary labour
// fixed at 6h raises surplus-value "both absolutely and relatively". The
// relative rise is the rate s/v: 360/360 = 1.0 before, 480/360 = 4/3 after.
func TestComputeOutcome_S3_LengtheningRaisesRateOfSurplusValue(t *testing.T) {
	t.Parallel()
	twelve := ComputeOutcome(LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360},
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	})
	fourteen := ComputeOutcome(LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 480},
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	})
	if fourteen.SurplusLabourMinutes <= twelve.SurplusLabourMinutes {
		t.Errorf("absolute surplus must rise: 12h=%d, 14h=%d", twelve.SurplusLabourMinutes, fourteen.SurplusLabourMinutes)
	}
	if twelve.RateOfSurplusValue != 1.0 {
		t.Errorf("12h rate: got %v, want 1.0", twelve.RateOfSurplusValue)
	}
	if want := 480.0 / 360.0; math.Abs(fourteen.RateOfSurplusValue-want) > 1e-9 {
		t.Errorf("14h rate: got %v, want %v", fourteen.RateOfSurplusValue, want)
	}
	if fourteen.LabourPowerValueMinutes != twelve.LabourPowerValueMinutes {
		t.Errorf("labour-power value must not change with day length: %d vs %d",
			twelve.LabourPowerValueMinutes, fourteen.LabourPowerValueMinutes)
	}
}

// Derived from the §3 setup (not a Marx quotation): shortening the day from
// 12h to 10h with necessary labour fixed at 6h lowers surplus-value both
// absolutely (360 → 240) and relatively (rate 1.0 → 2/3).
func TestComputeOutcome_S3_ShorteningLowersSurplusAbsolutelyAndRelatively(t *testing.T) {
	t.Parallel()
	before := ComputeOutcome(LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360},
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	})
	after := ComputeOutcome(LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 240},
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	})
	if after.SurplusLabourMinutes != 240 || after.DailyValueMinutes != 600 {
		t.Errorf("10h day: surplus=%d daily=%d, want 240 and 600", after.SurplusLabourMinutes, after.DailyValueMinutes)
	}
	if after.RateOfSurplusValue >= before.RateOfSurplusValue {
		t.Errorf("rate must fall: before=%v after=%v", before.RateOfSurplusValue, after.RateOfSurplusValue)
	}
}

// Marx Ch. 17 §4A: lengthening the day by 4h exactly offsets the rise in the
// value of labour-power from 3s to 4s. Surplus-value rises absolutely
// (360 → 480) but the rate s/v is unchanged at 1.0 — the rise is absolute only.
func TestComputeOutcome_S4A_AbsoluteRiseWithUnchangedRate(t *testing.T) {
	t.Parallel()
	before := ComputeOutcome(LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360},
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	})
	after := ComputeOutcome(LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 480, SurplusLabourMinutes: 480},
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	})
	if after.SurplusLabourMinutes <= before.SurplusLabourMinutes {
		t.Errorf("absolute surplus must rise: %d → %d", before.SurplusLabourMinutes, after.SurplusLabourMinutes)
	}
	if before.RateOfSurplusValue != after.RateOfSurplusValue {
		t.Errorf("rate must be unchanged: %v → %v", before.RateOfSurplusValue, after.RateOfSurplusValue)
	}
}

// Marx Ch. 17 §4B: increased productiveness and intensity both shorten the
// necessary labour, so the working-day can be shortened while surplus-value
// holds up. Day 12h → 10h (720 → 600 min), productiveness 4/3 (necessary
// 360 → 270 clock minutes) and intensity 6/5. Daily value stays 720, the
// value of labour-power is 324, and surplus-value rises 360 → 396 and the rate
// 1.0 → 11/9 even though the day is shorter.
func TestComputeOutcome_S4B_ShorterDayWithIntensityAndProductivity(t *testing.T) {
	t.Parallel()
	before := ComputeOutcome(LabourScenario{
		WorkingDay:   WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360},
		Intensity:    NormalLabourIntensity,
		Productivity: NormalLabourProductivity,
	})
	shorter := WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 240}
	after := ComputeOutcome(LabourScenario{
		WorkingDay:   shorter,
		Intensity:    LabourIntensity{Factor: 1.2},
		Productivity: LabourProductivity{Factor: 4.0 / 3.0},
	})

	if shorter.TotalMinutes() >= 720 {
		t.Fatalf("setup: day must be shorter than 12h, got %d", shorter.TotalMinutes())
	}
	if after.NecessaryLabourMinutes != 270 {
		t.Errorf("necessary labour: got %d, want 270", after.NecessaryLabourMinutes)
	}
	if after.DailyValueMinutes != 720 || after.LabourPowerValueMinutes != 324 {
		t.Errorf("daily value=%d, labour-power value=%d, want 720 and 324",
			after.DailyValueMinutes, after.LabourPowerValueMinutes)
	}
	if after.SurplusLabourMinutes != 396 {
		t.Errorf("surplus-value: got %d, want 396", after.SurplusLabourMinutes)
	}
	if after.SurplusLabourMinutes <= before.SurplusLabourMinutes {
		t.Errorf("surplus-value must rise despite the shorter day: %d → %d",
			before.SurplusLabourMinutes, after.SurplusLabourMinutes)
	}
	if want := 396.0 / 324.0; math.Abs(after.RateOfSurplusValue-want) > 1e-9 {
		t.Errorf("rate: got %v, want %v", after.RateOfSurplusValue, want)
	}
	// The minimum day is fixed by the necessary labour: surplus labour never
	// vanishes under capital, so necessary labour stays strictly below the day.
	if int64(after.NecessaryLabourMinutes) >= shorter.TotalMinutes() {
		t.Errorf("necessary labour %d must stay below the day %d",
			after.NecessaryLabourMinutes, shorter.TotalMinutes())
	}
}
