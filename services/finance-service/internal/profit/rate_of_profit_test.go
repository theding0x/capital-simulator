package profit

import "testing"

// Capital III, Ch. 2: a capital of 400c + 100v + 100s. The same surplus-value
// of £100 yields a rate of profit p' = s/C = 100/500 = 20% and a rate of
// surplus-value s' = s/v = 100/100 = 100%. The profit-form thus conceals four
// fifths of the rate (mystification 80%).
func TestComputeProfitRateAnalysis_Marx(t *testing.T) {
	t.Parallel()
	a := ComputeProfitRateAnalysis(400, 100, 100)

	if a.ProfitRate.TotalCapital != 500 {
		t.Errorf("total capital C = %d, want 500", a.ProfitRate.TotalCapital)
	}
	if a.ProfitRate.BasisPoints != 2000 {
		t.Errorf("p' = %d bp, want 2000 (20%%)", a.ProfitRate.BasisPoints)
	}
	if a.SurplusValueRate != 10000 {
		t.Errorf("s' = %d bp, want 10000 (100%%)", a.SurplusValueRate)
	}
	if a.Mystification != 8000 {
		t.Errorf("mystification = %d bp, want 8000 (80%%)", a.Mystification)
	}
	if a.ProfitRate.SurplusValue != 100 {
		t.Errorf("surplus-value = %d, want 100", a.ProfitRate.SurplusValue)
	}
}

// p' : s' = v : C, so the rate of profit is always smaller than the rate of
// surplus-value whenever any constant capital is advanced (c > 0).
func TestComputeProfitRateAnalysis_ProfitRateBelowSurplusRate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		c ConstantCapital
		v VariableCapital
		s SurplusValue
	}{
		{400, 100, 100}, // Marx's example
		{100, 100, 100}, // equal composition
		{1, 100, 100},   // tiny constant capital
		{900, 100, 50},  // high organic composition
	}
	for _, tc := range cases {
		a := ComputeProfitRateAnalysis(tc.c, tc.v, tc.s)
		if a.ProfitRate.BasisPoints >= int64(a.SurplusValueRate) {
			t.Errorf("c=%d v=%d s=%d: p'=%d should be < s'=%d", tc.c, tc.v, tc.s, a.ProfitRate.BasisPoints, a.SurplusValueRate)
		}
	}
}

// The only case of equality is v = C, i.e. no constant capital. Then the rate
// of profit equals the rate of surplus-value and nothing is mystified.
func TestComputeProfitRateAnalysis_EqualWhenNoConstantCapital(t *testing.T) {
	t.Parallel()
	a := ComputeProfitRateAnalysis(0, 100, 100)

	if a.ProfitRate.BasisPoints != int64(a.SurplusValueRate) {
		t.Errorf("with c=0, p'=%d should equal s'=%d", a.ProfitRate.BasisPoints, a.SurplusValueRate)
	}
	if a.ProfitRate.BasisPoints != 10000 {
		t.Errorf("p' = %d bp, want 10000 (100%%)", a.ProfitRate.BasisPoints)
	}
	if a.Mystification != 0 {
		t.Errorf("mystification = %d bp, want 0 (nothing concealed without constant capital)", a.Mystification)
	}
}

// MystificationDegree is (s' - p') / s' in basis points across the board.
func TestComputeProfitRateAnalysis_MystificationFormula(t *testing.T) {
	t.Parallel()
	cases := []struct {
		c ConstantCapital
		v VariableCapital
		s SurplusValue
	}{
		{400, 100, 100},
		{100, 100, 100},
		{0, 100, 100},
		{900, 100, 50},
	}
	for _, tc := range cases {
		a := ComputeProfitRateAnalysis(tc.c, tc.v, tc.s)
		want := MystificationDegree((int64(a.SurplusValueRate) - a.ProfitRate.BasisPoints) * basisPoints / int64(a.SurplusValueRate))
		if a.Mystification != want {
			t.Errorf("c=%d v=%d s=%d: mystification = %d, want %d", tc.c, tc.v, tc.s, a.Mystification, want)
		}
	}
}

// A zero total capital (nothing advanced) yields a zero rate of profit rather
// than a division by zero.
func TestComputeRateOfProfit_ZeroTotalCapital(t *testing.T) {
	t.Parallel()
	p := ComputeRateOfProfit(100, 0, 0)
	if p.TotalCapital != 0 || p.BasisPoints != 0 {
		t.Errorf("got C=%d p'=%d, want 0/0 for an empty advance", p.TotalCapital, p.BasisPoints)
	}
}

// Surplus-value originates in variable capital; with v = 0 there is no rate of
// surplus-value to measure, and the analysis reports zero rather than dividing
// by zero.
func TestComputeRateOfSurplusValue_ZeroVariable(t *testing.T) {
	t.Parallel()
	if got := ComputeRateOfSurplusValue(100, 0); got != 0 {
		t.Errorf("s' with v=0 = %d, want 0", got)
	}
	a := ComputeProfitRateAnalysis(400, 0, 0)
	if a.SurplusValueRate != 0 || a.Mystification != 0 {
		t.Errorf("with v=0: s'=%d myst=%d, want 0/0", a.SurplusValueRate, a.Mystification)
	}
}

// ProfitRateID is 96-bit hex (24 chars) and non-empty.
func TestNewProfitRateID(t *testing.T) {
	t.Parallel()
	id := NewProfitRateID()
	if id.IsZero() {
		t.Fatal("NewProfitRateID returned the zero value")
	}
	if len(string(id)) != 24 {
		t.Errorf("id length = %d, want 24 hex chars", len(string(id)))
	}
}
