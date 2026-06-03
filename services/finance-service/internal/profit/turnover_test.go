package profit

import (
	"errors"
	"testing"
)

// Marx's §main worked capitals. Capital A: 80c + 20v = 100C, n = 2, s' = 100%
// → annual p' = 40%. Capital B: 160c + 40v = 200C, n = 1, s' = 100% → annual
// p' = 20% (equal composition to A, half the rate). Capital I: fixed £10,000,
// circulating c £500, v £500, n = 10 → C = £11,000, p' = 5,000/11,000 = 45 5/11%.
func TestComputeTurnoverAnalysis_Marx(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		c          TotalCapital
		v          VariableCapitalPerTurnover
		sRate      RateOfSurplusValue
		n          TurnoverCount
		wantAnnual int64
	}{
		{"Capital A: 100C, n=2, s'=100%", 100, 20, 10000, 2, 4000},
		{"Capital B: 200C, n=1, s'=100%", 200, 40, 10000, 1, 2000},
		{"Capital I: 11000C, n=10, s'=100%", 11000, 500, 10000, 10, 4545},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := ComputeTurnoverAnalysis(tc.c, tc.v, tc.sRate, tc.n)
			if a.AnnualProfitRate.BasisPoints != tc.wantAnnual {
				t.Errorf("annual p' = %d, want %d", a.AnnualProfitRate.BasisPoints, tc.wantAnnual)
			}
		})
	}
}

// The single-turnover rate is p' without the turnover factor, s'·(v/C); the
// annual rate multiplies it by n. Capital A turns over twice, so its annual
// rate (40%) is twice its single-turnover rate (20%) — and that single-turnover
// rate is exactly Capital B's annual rate, B turning over but once.
func TestComputeTurnoverAnalysis_SingleVersusAnnual(t *testing.T) {
	t.Parallel()
	a := ComputeTurnoverAnalysis(100, 20, 10000, 2)
	if a.SingleTurnoverProfitRate != 2000 {
		t.Errorf("single-turnover p' = %d, want 2000", a.SingleTurnoverProfitRate)
	}
	if a.AnnualProfitRate.BasisPoints != a.SingleTurnoverProfitRate*int64(a.N) {
		t.Errorf("annual p' %d != single %d × n %d",
			a.AnnualProfitRate.BasisPoints, a.SingleTurnoverProfitRate, a.N)
	}
	b := ComputeTurnoverAnalysis(200, 40, 10000, 1)
	if a.SingleTurnoverProfitRate != b.AnnualProfitRate.BasisPoints {
		t.Errorf("A's single-turnover rate %d should equal B's annual rate %d",
			a.SingleTurnoverProfitRate, b.AnnualProfitRate.BasisPoints)
	}
}

// §main invariant: for capitals of equal composition and equal s', the rates of
// profit "are related inversely as their periods of turnover" — p'_A : p'_B =
// n_A : n_B. A (n=2) and B (n=1) share composition v/C = 20% and s' = 100%.
func TestComputeTurnoverAnalysis_InverseAsTurnoverPeriods(t *testing.T) {
	t.Parallel()
	a := ComputeTurnoverAnalysis(100, 20, 10000, 2)
	b := ComputeTurnoverAnalysis(200, 40, 10000, 1)
	// Equal composition: both v/C = 20%.
	if a.SingleTurnoverProfitRate != b.SingleTurnoverProfitRate {
		t.Fatalf("compositions differ: A single %d, B single %d",
			a.SingleTurnoverProfitRate, b.SingleTurnoverProfitRate)
	}
	// p'_A : p'_B = n_A : n_B — cross-multiply to avoid division.
	lhs := a.AnnualProfitRate.BasisPoints * int64(b.N)
	rhs := b.AnnualProfitRate.BasisPoints * int64(a.N)
	if lhs != rhs {
		t.Errorf("p'_A:p'_B (%d:%d) not as n_A:n_B (%d:%d)",
			a.AnnualProfitRate.BasisPoints, b.AnnualProfitRate.BasisPoints, a.N, b.N)
	}
}

// Invariant [§formula]: S' = s'·n. Capital A: 100% over two turnovers → 200%.
func TestComputeTurnoverAnalysis_AnnualSurplusValueRate(t *testing.T) {
	t.Parallel()
	a := ComputeTurnoverAnalysis(100, 20, 10000, 2)
	if a.AnnualSurplusValueRate != 20000 {
		t.Errorf("S' = %d, want 20000 (s'·n)", a.AnnualSurplusValueRate)
	}
	// The annual wage bill vn — what the capitalist's books show — is v·n.
	if a.AnnualWages != 40 {
		t.Errorf("annual wages = %d, want 40 (v·n)", a.AnnualWages)
	}
}

// Invariant [§formula]: annual p' = S'·(v/C) ≤ S', since v ≤ C. The annual rate
// of profit can never exceed the annual rate of surplus-value.
func TestComputeTurnoverAnalysis_ProfitRateNeverExceedsSurplusRate(t *testing.T) {
	t.Parallel()
	for _, a := range []TurnoverAnalysis{
		ComputeTurnoverAnalysis(100, 20, 10000, 2),
		ComputeTurnoverAnalysis(11000, 500, 10000, 10),
		ComputeTurnoverAnalysis(100, 100, 15000, 3), // limiting case v = C
	} {
		if a.AnnualProfitRate.BasisPoints > int64(a.AnnualSurplusValueRate) {
			t.Errorf("annual p' %d exceeds S' %d (v=%d, C=%d)",
				a.AnnualProfitRate.BasisPoints, a.AnnualSurplusValueRate, a.V, a.C)
		}
	}
	// At v = C the two rates coincide: p' = s'·n·1 = S'.
	limiting := ComputeTurnoverAnalysis(100, 100, 15000, 3)
	if limiting.AnnualProfitRate.BasisPoints != int64(limiting.AnnualSurplusValueRate) {
		t.Errorf("v=C: annual p' %d should equal S' %d",
			limiting.AnnualProfitRate.BasisPoints, limiting.AnnualSurplusValueRate)
	}
}

// §main: "should capital I have only 5 instead of 10 turnovers ... the rate of
// profit has fallen one-half, because the period of turnover has doubled."
// Halving n halves annual p' when composition and s' are held constant — exact
// for Capital A (no integer-truncation drift).
func TestComputeTurnoverAnalysis_HalvingTurnoverHalvesRate(t *testing.T) {
	t.Parallel()
	full := ComputeTurnoverAnalysis(100, 20, 10000, 2)
	half := ComputeTurnoverAnalysis(100, 20, 10000, 1)
	if full.AnnualProfitRate.BasisPoints != 2*half.AnnualProfitRate.BasisPoints {
		t.Errorf("n=2 rate %d should be twice n=1 rate %d",
			full.AnnualProfitRate.BasisPoints, half.AnnualProfitRate.BasisPoints)
	}
}

// Boundaries: no capital advanced and no turnover both yield a zero annual rate.
func TestComputeTurnoverAnalysis_ZeroBoundaries(t *testing.T) {
	t.Parallel()
	if got := ComputeTurnoverAnalysis(0, 0, 10000, 2); got.AnnualProfitRate.BasisPoints != 0 {
		t.Errorf("C=0: annual p' = %d, want 0", got.AnnualProfitRate.BasisPoints)
	}
	noTurnover := ComputeTurnoverAnalysis(100, 20, 10000, 0)
	if noTurnover.AnnualProfitRate.BasisPoints != 0 || noTurnover.AnnualSurplusValueRate != 0 {
		t.Errorf("n=0: annual p' = %d, S' = %d, want 0/0",
			noTurnover.AnnualProfitRate.BasisPoints, noTurnover.AnnualSurplusValueRate)
	}
	if noTurnover.SingleTurnoverProfitRate != 2000 {
		t.Errorf("n=0: single-turnover p' = %d, want 2000 (independent of n)",
			noTurnover.SingleTurnoverProfitRate)
	}
}

func TestComputeAnnualProfitRate_CarriesFactors(t *testing.T) {
	t.Parallel()
	r := ComputeAnnualProfitRate(10000, 2, 20, 100)
	if r.SurplusValueRate != 10000 || r.Turnovers != 2 || r.VariableCapital != 20 || r.TotalCapital != 100 {
		t.Errorf("factors not carried: %+v", r)
	}
	if r.BasisPoints != 4000 {
		t.Errorf("bp = %d, want 4000", r.BasisPoints)
	}
}

func TestTurnoverAnalysis_Validate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		a    TurnoverAnalysis
		want error
	}{
		{"valid", TurnoverAnalysis{C: 100, V: 20, SRate: 10000, N: 2}, nil},
		{"negative s'", TurnoverAnalysis{C: 100, V: 20, SRate: -1, N: 2}, ErrNegativeMagnitude},
		{"negative n", TurnoverAnalysis{C: 100, V: 20, SRate: 10000, N: -1}, ErrNegativeMagnitude},
		{"v exceeds C", TurnoverAnalysis{C: 100, V: 120, SRate: 10000, N: 2}, ErrVariableExceedsTotal},
		{"v equals C", TurnoverAnalysis{C: 100, V: 100, SRate: 10000, N: 2}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.a.Validate(); !errors.Is(err, tc.want) {
				t.Errorf("Validate() = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestNewTurnoverAnalysisID_Unique(t *testing.T) {
	t.Parallel()
	a, b := NewTurnoverAnalysisID(), NewTurnoverAnalysisID()
	if a == b {
		t.Error("two IDs collided")
	}
	if a.IsZero() || len(string(a)) != 24 {
		t.Errorf("id = %q, want 24 hex chars", a)
	}
}

// §Manchester — the 1871 Manchester cotton-spinnery Marx works through: total
// capital C = £12,500, variable capital advanced v = £318, s' = 153 11/13%
// (≈ 15385 bp), turning over n ≈ 8½ times a year. The integer TurnoverCount
// cannot hold a half-turnover, so the model floors Marx's 8½ to 8 — yielding
// p' = 3131 bp and S' = 123080 bp, the integer-model approximation of Marx's
// idealised 8½ → ≈ 3327 bp (33.27%) / ≈ 130709 bp (1307.09%). The signature
// worked example is thereby exercised, with its turnover limitation made
// explicit rather than silently dropped.
func TestComputeTurnoverAnalysis_Manchester1871(t *testing.T) {
	t.Parallel()
	const (
		c     TotalCapital               = 12500
		v     VariableCapitalPerTurnover = 318
		sRate RateOfSurplusValue         = 15385 // 153 11/13% in basis points
		n     TurnoverCount              = 8      // Marx's 8½ floored to the integer model
	)
	a := ComputeTurnoverAnalysis(c, v, sRate, n)

	if got := a.AnnualProfitRate.BasisPoints; got != 3131 {
		t.Errorf("annual p' = %d bp, want 3131 (s'·n·v/C = 15385·8·318/12500)", got)
	}
	if got := int64(a.AnnualSurplusValueRate); got != 123080 {
		t.Errorf("annual S' = %d bp, want 123080 (s'·n = 15385·8)", got)
	}
	if got := a.SingleTurnoverProfitRate; got != 391 {
		t.Errorf("single-turnover p' = %d bp, want 391 (s'·v/C)", got)
	}
	if got := int64(a.AnnualWages); got != 2544 {
		t.Errorf("annual wage bill vn = %d, want 2544 (v·n = 318·8)", got)
	}
	// The annual rate of profit never exceeds the annual rate of surplus-value,
	// since v ≤ C (the chapter's bounding invariant).
	if a.AnnualProfitRate.BasisPoints > int64(a.AnnualSurplusValueRate) {
		t.Errorf("annual p' %d must not exceed annual S' %d", a.AnnualProfitRate.BasisPoints, a.AnnualSurplusValueRate)
	}
}
