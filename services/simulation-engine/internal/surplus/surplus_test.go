package surplus_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
)

func TestAbsoluteWorkdayLimit(t *testing.T) {
	t.Parallel()
	if surplus.AbsoluteWorkdayLimit != 1440 {
		t.Fatalf("expected 1440, got %d", surplus.AbsoluteWorkdayLimit)
	}
}

func TestSurplusValueRate_Rate(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	if got := rate.Rate(); got != 1.0 {
		t.Fatalf("expected 1.0, got %f", got)
	}
	rate200 := surplus.SurplusValueRate{SurplusLabour: 12, NecessaryLabour: 6}
	if got := rate200.Rate(); got != 2.0 {
		t.Fatalf("expected 2.0, got %f", got)
	}
}

// §1 fixture: "if the rate of surplus-value be = 100%, this variable capital of
// 3s. produces a mass of surplus-value of 3s."
func TestMassByRate_100pct(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	got := surplus.MassByRate(rate, surplus.VariableCapital(3))
	if got != surplus.SurplusValueMass(3) {
		t.Fatalf("expected 3, got %d", got)
	}
}

// §1 fixture: 100 labourers at rate 100%, V=300s → S=300s
func TestMassByWorkers_100workers(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	got := surplus.MassByWorkers(surplus.LabourPowerValue(3), rate, surplus.WorkerCount(100))
	if got != surplus.SurplusValueMass(300) {
		t.Fatalf("expected 300, got %d", got)
	}
}

// §1 compensation law: rate doubles (100%→200%), V halves (300s→150s), workers halve (100→50).
// S stays constant at 300 (not 150 — the spec fixture has a typo; the formula is correct).
func TestMassByRate_CompensationLaw(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 12, NecessaryLabour: 6}
	got := surplus.MassByRate(rate, surplus.VariableCapital(150))
	// S = (12/6) × 150 = 2 × 150 = 300
	if got != surplus.SurplusValueMass(300) {
		t.Fatalf("expected 300, got %d", got)
	}
}

// §1 fixture: V=1500s, rate=100% → S=1500s
func TestMassByRate_Large(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	got := surplus.MassByRate(rate, surplus.VariableCapital(1500))
	if got != surplus.SurplusValueMass(1500) {
		t.Fatalf("expected 1500, got %d", got)
	}
}

// §1 fixture: "A capital of 300s. that employs 100 labourers a day with a rate of
// surplus-value of 200% … produces only a mass of surplus-value of 600s."
func TestMassByRate_200pct(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 12, NecessaryLabour: 6}
	got := surplus.MassByRate(rate, surplus.VariableCapital(300))
	if got != surplus.SurplusValueMass(600) {
		t.Fatalf("expected 600, got %d", got)
	}
}

// §1 minimum capital: "he would have to employ two labourers in order to live …"
func TestMinimumCapital(t *testing.T) {
	t.Parallel()
	got := surplus.MinimumCapital(surplus.LabourPowerValue(3), surplus.WorkerCount(2))
	if got != surplus.VariableCapital(6) {
		t.Fatalf("expected 6, got %d", got)
	}
}

// Invariant: MassByRate and MassByWorkers agree when V = v × n [§1]
func TestInvariant_FormulaAgreement(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	v := surplus.LabourPowerValue(3)
	n := surplus.WorkerCount(100)
	V := surplus.VariableCapital(int64(v) * int64(n)) // 300
	byRate := surplus.MassByRate(rate, V)
	byWorkers := surplus.MassByWorkers(v, rate, n)
	if byRate != byWorkers {
		t.Fatalf("formula disagreement: MassByRate=%d MassByWorkers=%d", byRate, byWorkers)
	}
}

// Invariant: IndividualSurplus × n == MassByWorkers [§1]
func TestInvariant_IndividualScaled(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	v := surplus.LabourPowerValue(3)
	n := surplus.WorkerCount(100)
	individual := surplus.IndividualSurplus(rate, v)
	mass := surplus.MassByWorkers(v, rate, n)
	if surplus.SurplusValueMass(int64(individual)*int64(n)) != mass {
		t.Fatalf("scale invariant failed: %d × %d ≠ %d", individual, n, mass)
	}
}

// Invariant: working day (necessary + surplus) < AbsoluteWorkdayLimit [§1]
func TestInvariant_WorkingDayBelowAbsoluteLimit(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	total := surplus.LabourMinutes(rate.NecessaryLabour + rate.SurplusLabour)
	if total >= surplus.AbsoluteWorkdayLimit {
		t.Fatalf("working day %d must be < %d", total, surplus.AbsoluteWorkdayLimit)
	}
}
