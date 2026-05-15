package simulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

func TestRunSimpleReproduction_CoreExample(t *testing.T) {
	t.Parallel()
	// "capital of £1,000 begets yearly a surplus-value of £200"
	// c=800, v=200, s'=100%, 5 periods
	cycles := simulation.RunSimpleReproduction(
		simulation.CapitalStock{ConstantCapital: 800, VariableCapital: 200},
		1.0, 5,
	)
	if len(cycles) != 5 {
		t.Fatalf("want 5 cycles, got %d", len(cycles))
	}
	for _, c := range cycles {
		if c.Fund.Total != 200 {
			t.Errorf("period %d: Fund.Total = %d, want 200", c.Period, c.Fund.Total)
		}
		if c.Fund.Revenue != 200 {
			t.Errorf("period %d: Fund.Revenue = %d, want 200", c.Period, c.Fund.Revenue)
		}
		if c.Fund.Accumulated != 0 {
			t.Errorf("period %d: Fund.Accumulated = %d, want 0", c.Period, c.Fund.Accumulated)
		}
	}
}

func TestRepaymentPeriod_CoreExample(t *testing.T) {
	t.Parallel()
	// "at the end of 5 years the surplus-value consumed will amount to 5×£200
	// or the £1,000 originally advanced"
	got := simulation.RepaymentPeriod(1000, 200)
	if got != 5 {
		t.Fatalf("RepaymentPeriod(1000, 200): want 5, got %d", got)
	}
}

func TestRepaymentPeriod_HalfConsumed(t *testing.T) {
	t.Parallel()
	// "if only a part, say one half, were consumed ... the same result would
	// follow at the end of 10 years"
	got := simulation.RepaymentPeriod(1000, 100)
	if got != 10 {
		t.Fatalf("RepaymentPeriod(1000, 100): want 10, got %d", got)
	}
}

func TestRepaymentPeriod_GeneralRule(t *testing.T) {
	t.Parallel()
	// "value of capital advanced / surplus-value annually consumed = number of
	// reproduction periods"
	cases := []struct {
		capital simulation.Pence
		revenue simulation.Pence
		want    int64
	}{
		{600, 200, 3},
		{800, 400, 2},
		{1200, 300, 4},
	}
	for _, tc := range cases {
		got := simulation.RepaymentPeriod(tc.capital, tc.revenue)
		if got != tc.want {
			t.Errorf("RepaymentPeriod(%d, %d): want %d, got %d", tc.capital, tc.revenue, tc.want, got)
		}
	}
}

func TestRepaymentPeriod_RoundsUp(t *testing.T) {
	t.Parallel()
	// 1000 / 300 = 3.33… → rounds up to 4
	got := simulation.RepaymentPeriod(1000, 300)
	if got != 4 {
		t.Fatalf("RepaymentPeriod(1000, 300): want 4, got %d", got)
	}
}

func TestRunSimpleReproduction_CapitalStockInvariant(t *testing.T) {
	t.Parallel()
	// Under simple reproduction the capital stock does not grow
	initial := simulation.CapitalStock{ConstantCapital: 800, VariableCapital: 200}
	cycles := simulation.RunSimpleReproduction(initial, 1.0, 5)
	for _, c := range cycles {
		if c.Capital.ConstantCapital != initial.ConstantCapital ||
			c.Capital.VariableCapital != initial.VariableCapital {
			t.Errorf("period %d: capital changed: got {%d,%d}, want {%d,%d}",
				c.Period,
				c.Capital.ConstantCapital, c.Capital.VariableCapital,
				initial.ConstantCapital, initial.VariableCapital)
		}
	}
}

func TestRunSimpleReproduction_SurplusConservation(t *testing.T) {
	t.Parallel()
	// After the repayment period, total consumed surplus-value >= original capital
	initialCapital := simulation.Pence(1000)
	annualRevenue := simulation.Pence(200)
	n := simulation.RepaymentPeriod(initialCapital, annualRevenue)
	cycles := simulation.RunSimpleReproduction(
		simulation.CapitalStock{ConstantCapital: 800, VariableCapital: 200},
		1.0, n,
	)
	var total simulation.Pence
	for _, c := range cycles {
		total += c.Fund.Revenue
	}
	if total < initialCapital {
		t.Errorf("consumed surplus-value %d < original capital %d", total, initialCapital)
	}
}

func TestRunSimpleReproduction_ZeroPeriods(t *testing.T) {
	t.Parallel()
	cycles := simulation.RunSimpleReproduction(
		simulation.CapitalStock{ConstantCapital: 800, VariableCapital: 200},
		1.0, 0,
	)
	if len(cycles) != 0 {
		t.Fatalf("want empty slice, got %d cycles", len(cycles))
	}
}

func TestRunSimpleReproduction_SimpleAccumulatedZero(t *testing.T) {
	t.Parallel()
	// simple reproduction invariant: Accumulated == 0, Revenue == Total
	cycles := simulation.RunSimpleReproduction(
		simulation.CapitalStock{ConstantCapital: 500, VariableCapital: 500},
		0.5, 3,
	)
	for _, c := range cycles {
		if c.Fund.Accumulated != 0 {
			t.Errorf("period %d: Accumulated = %d, want 0", c.Period, c.Fund.Accumulated)
		}
		if c.Fund.Revenue != c.Fund.Total {
			t.Errorf("period %d: Revenue %d != Total %d", c.Period, c.Fund.Revenue, c.Fund.Total)
		}
	}
}
