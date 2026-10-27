package surplus_test

import (
	"errors"
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

// Invariant: necessary + surplus < AbsoluteWorkdayLimit [§1, second law].
func TestRateValidate_WorkdayLimit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		rate    surplus.SurplusValueRate
		wantErr bool
	}{
		{"six and six", surplus.SurplusValueRate{SurplusLabour: 360, NecessaryLabour: 360}, false},
		{"one minute short of 24h", surplus.SurplusValueRate{SurplusLabour: 1079, NecessaryLabour: 360}, false},
		{"exactly 24h", surplus.SurplusValueRate{SurplusLabour: 1080, NecessaryLabour: 360}, true},
		{"beyond 24h", surplus.SurplusValueRate{SurplusLabour: 2000, NecessaryLabour: 360}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.rate.Validate()
			if tc.wantErr && !errors.Is(err, surplus.ErrExceedsWorkdayLimit) {
				t.Fatalf("expected ErrExceedsWorkdayLimit, got %v", err)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// §1 second law: 100 labourers at v=3s (6h necessary) "can never reach the sum
// of 1,200s. or 24 × 100 working-hours". The longest admissible day yields a
// total product (V+S) strictly below 1,200s; the 24-hour day that would reach it
// is rejected.
func TestInvariant_CompensationNeverReaches1200s(t *testing.T) {
	t.Parallel()
	v := surplus.LabourPowerValue(3)
	n := surplus.WorkerCount(100)
	V := int64(surplus.MinimumCapital(v, n)) // 300

	longest := surplus.SurplusValueRate{SurplusLabour: 1079, NecessaryLabour: 360}
	if err := longest.Validate(); err != nil {
		t.Fatalf("longest admissible day rejected: %v", err)
	}
	if total := V + int64(surplus.MassByWorkers(v, longest, n)); total >= 1200 {
		t.Fatalf("V+S = %d, must stay below 1200", total)
	}

	full := surplus.SurplusValueRate{SurplusLabour: 1080, NecessaryLabour: 360}
	if err := full.Validate(); !errors.Is(err, surplus.ErrExceedsWorkdayLimit) {
		t.Fatalf("24-hour day must be rejected, got %v", err)
	}
	if total := V + int64(surplus.MassByWorkers(v, full, n)); total != 1200 {
		t.Fatalf("fixture sanity: 24h day should give V+S = 1200, got %d", total)
	}
}

func TestRate_ZeroNecessaryLabour(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 0}
	if got := rate.Rate(); got != 0 {
		t.Fatalf("expected 0 for zero necessary_labour, got %f", got)
	}
}

func TestMassByRate_ZeroNecessaryLabour(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 0}
	if got := surplus.MassByRate(rate, surplus.VariableCapital(300)); got != 0 {
		t.Fatalf("expected 0 for zero necessary_labour, got %d", got)
	}
}

func TestMassByWorkers_ZeroNecessaryLabour(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 0}
	if got := surplus.MassByWorkers(surplus.LabourPowerValue(3), rate, surplus.WorkerCount(100)); got != 0 {
		t.Fatalf("expected 0 for zero necessary_labour, got %d", got)
	}
}

func TestIndividualSurplus_ZeroNecessaryLabour(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 0}
	if got := surplus.IndividualSurplus(rate, surplus.LabourPowerValue(3)); got != 0 {
		t.Fatalf("expected 0 for zero necessary_labour, got %d", got)
	}
}
