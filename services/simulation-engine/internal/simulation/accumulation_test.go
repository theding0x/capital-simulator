package simulation_test

import (
	"math"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// SplitSurplus — § spinner example §1: £10,000 capital, 4/5 constant + 1/5 variable,
// surplus-value 100%; accumulate all: reinvest £1,600 in cotton/machinery, £400 in wages.
func TestSplitSurplus_SpinnerFullAccumulation(t *testing.T) {
	t.Parallel()
	got := simulation.SplitSurplus(2000, 1.0, 0.8)
	if got.Constant != 1600 {
		t.Errorf("constant = %d, want 1600", got.Constant)
	}
	if got.Variable != 400 {
		t.Errorf("variable = %d, want 400", got.Variable)
	}
}

// Under full accumulation Revenue must be zero — entire surplus converted to capital.
func TestSplitSurplus_FullAccumulationRevenue(t *testing.T) {
	t.Parallel()
	additional := simulation.SplitSurplus(2000, 1.0, 0.8)
	revenue := simulation.Pence(2000) - additional.Constant - additional.Variable
	if revenue != 0 {
		t.Errorf("revenue = %d, want 0 under full accumulation", revenue)
	}
}

// § partial accumulation §3: capitalist consumes half; AccumulationRate=0.5.
func TestSplitSurplus_PartialAccumulation(t *testing.T) {
	t.Parallel()
	got := simulation.SplitSurplus(2000, 0.5, 0.8)
	if got.Constant != 800 {
		t.Errorf("constant = %d, want 800", got.Constant)
	}
	if got.Variable != 200 {
		t.Errorf("variable = %d, want 200", got.Variable)
	}
	revenue := simulation.Pence(2000) - got.Constant - got.Variable
	if revenue != 1000 {
		t.Errorf("revenue = %d, want 1000", revenue)
	}
}

// Invariant: Constant == int64(float64(surplus) * rate * ratio).
func TestSplitSurplus_ConstantInvariant(t *testing.T) {
	t.Parallel()
	surplus := simulation.Pence(2000)
	rate := 1.0
	ratio := 0.8
	got := simulation.SplitSurplus(surplus, rate, ratio)
	want := simulation.Pence(int64(float64(surplus) * rate * ratio))
	if got.Constant != want {
		t.Errorf("constant = %d, want %d", got.Constant, want)
	}
}

// Invariant: Variable is the exact complement of Constant within the accumulated amount.
// Using math.Round on the accumulated total avoids float64 truncation (e.g. 2000*0.2=399.999...).
func TestSplitSurplus_VariableInvariant(t *testing.T) {
	t.Parallel()
	surplus := simulation.Pence(2000)
	rate := 0.5
	ratio := 0.8
	got := simulation.SplitSurplus(surplus, rate, ratio)
	accumulated := int64(math.Round(float64(surplus) * rate))
	wantConstant := int64(math.Round(float64(accumulated) * ratio))
	want := simulation.Pence(accumulated - wantConstant)
	if got.Variable != want {
		t.Errorf("variable = %d, want %d", got.Variable, want)
	}
}

// Zero accumulation rate: all surplus is revenue, no new capital formed.
func TestSplitSurplus_ZeroAccumulationRate(t *testing.T) {
	t.Parallel()
	got := simulation.SplitSurplus(2000, 0.0, 0.8)
	if got.Constant != 0 {
		t.Errorf("constant = %d, want 0", got.Constant)
	}
	if got.Variable != 0 {
		t.Errorf("variable = %d, want 0", got.Variable)
	}
}

// § compound §1: second cycle on new capital {1600c+400v} generates £400 surplus;
// SplitSurplus(400, 1.0, 0.8) == AdditionalCapital{320, 80}.
func TestSplitSurplus_CompoundSecondCycle(t *testing.T) {
	t.Parallel()
	got := simulation.SplitSurplus(400, 1.0, 0.8)
	if got.Constant != 320 {
		t.Errorf("constant = %d, want 320", got.Constant)
	}
	if got.Variable != 80 {
		t.Errorf("variable = %d, want 80", got.Variable)
	}
}

// § Abraham begat Isaac §1: £10,000 → £2,000 surplus → £400 new surplus → £80 further surplus.
// First three iterations of RunExtendedReproduction(CapitalStock{8000,2000}, 1.0, 1.0, 0.8, 3).
func TestRunExtendedReproduction_AbrahamBegatsIsaac(t *testing.T) {
	t.Parallel()
	initial := simulation.CapitalStock{ConstantCapital: 8000, VariableCapital: 2000}
	cycles := simulation.RunExtendedReproduction(initial, 1.0, 1.0, 0.8, 3)

	if len(cycles) != 3 {
		t.Fatalf("len = %d, want 3", len(cycles))
	}

	wantSurplus := []simulation.Pence{2000, 400, 80}
	for i, c := range cycles {
		if c.Period != int64(i+1) {
			t.Errorf("period[%d] = %d, want %d", i, c.Period, i+1)
		}
		if c.SurplusProduced != wantSurplus[i] {
			t.Errorf("surplus[%d] = %d, want %d", i, c.SurplusProduced, wantSurplus[i])
		}
	}
}

// § reinvestment §1: period 1 additional capital must match SplitSurplus(2000, 1.0, 0.8).
func TestRunExtendedReproduction_PeriodOneAdditional(t *testing.T) {
	t.Parallel()
	initial := simulation.CapitalStock{ConstantCapital: 8000, VariableCapital: 2000}
	cycles := simulation.RunExtendedReproduction(initial, 1.0, 1.0, 0.8, 1)

	c := cycles[0]
	if c.NewConstant != 1600 {
		t.Errorf("new_constant = %d, want 1600", c.NewConstant)
	}
	if c.NewVariable != 400 {
		t.Errorf("new_variable = %d, want 400", c.NewVariable)
	}
	if c.Revenue != 0 {
		t.Errorf("revenue = %d, want 0 (full accumulation)", c.Revenue)
	}
}

// Zero periods returns an empty slice, never nil.
func TestRunExtendedReproduction_ZeroPeriods(t *testing.T) {
	t.Parallel()
	initial := simulation.CapitalStock{ConstantCapital: 8000, VariableCapital: 2000}
	cycles := simulation.RunExtendedReproduction(initial, 1.0, 1.0, 0.8, 0)
	if cycles == nil {
		t.Error("want empty slice, got nil")
	}
	if len(cycles) != 0 {
		t.Errorf("len = %d, want 0", len(cycles))
	}
}
