package tendency

import (
	"encoding/hex"
	"errors"
	"testing"
)

// lawSteps is Marx's §law fixture for Ch. 13: s′ held at 100%, v constant at
// 100, and c rising 50 → 100 → 200 → 300 → 400. The rate of profit p′ = s/(c+v)
// falls monotonically across the five periods.
var lawSteps = [][2]int64{
	{50, 100},  // c=50, v=100 → C=150
	{100, 100}, // c=100 → C=200
	{200, 100}, // c=200 → C=300
	{300, 100}, // c=300 → C=400
	{400, 100}, // c=400 → C=500
}

// wantLawRates are the round-half-up basis-point rates for lawSteps. Period 1 is
// 1,000,000/150 = 6666.67 → 6667 (NOT 6666 — the house convention rounds half up).
var wantLawRates = []int64{6667, 5000, 3333, 2500, 2000}

// TestBuildCompositionTrajectory_LawRates asserts the §law fixture produces the
// documented falling sequence of profit rates, that the series is monotonically
// decreasing, and that the derived ProfitRates slice mirrors the per-period rate.
func TestBuildCompositionTrajectory_LawRates(t *testing.T) {
	t.Parallel()
	tr := BuildCompositionTrajectory("Marx §law s'=100%", 10000, lawSteps)

	if len(tr.Periods) != len(wantLawRates) {
		t.Fatalf("Periods len = %d, want %d", len(tr.Periods), len(wantLawRates))
	}
	if len(tr.ProfitRates) != len(wantLawRates) {
		t.Fatalf("ProfitRates len = %d, want %d", len(tr.ProfitRates), len(wantLawRates))
	}

	for i, want := range wantLawRates {
		if got := tr.Periods[i].ProfitRate; got != want {
			t.Errorf("period %d ProfitRate = %d, want %d", i+1, got, want)
		}
		if got := tr.ProfitRates[i]; got != want {
			t.Errorf("ProfitRates[%d] = %d, want %d (must mirror period rate)", i, got, want)
		}
	}

	// Monotonically decreasing across the whole trajectory.
	for i := 1; i < len(tr.ProfitRates); i++ {
		if tr.ProfitRates[i] >= tr.ProfitRates[i-1] {
			t.Errorf("series not strictly decreasing at %d: %d >= %d",
				i, tr.ProfitRates[i], tr.ProfitRates[i-1])
		}
	}
}

// TestBuildCompositionTrajectory_Period1IsRoundHalfUp pins the boundary case the
// spec calls out explicitly: Period 1 must be 6667, never 6666.
func TestBuildCompositionTrajectory_Period1IsRoundHalfUp(t *testing.T) {
	t.Parallel()
	tr := BuildCompositionTrajectory("p1", 10000, [][2]int64{{50, 100}})
	if got := tr.Periods[0].ProfitRate; got != 6667 {
		t.Fatalf("Period 1 ProfitRate = %d, want 6667 (round-half-up of 6666.67)", got)
	}
}

// TestBuildCompositionTrajectory_PeriodInvariant verifies, per period, the
// invariant ProfitRate == roundHalfUp(s*10000/(c+v)), computed independently of
// the production code via (s*10000 + (c+v)/2)/(c+v).
func TestBuildCompositionTrajectory_PeriodInvariant(t *testing.T) {
	t.Parallel()
	tr := BuildCompositionTrajectory("inv", 10000, lawSteps)
	for i, p := range tr.Periods {
		denom := p.ConstantCapital + p.VariableCapital
		want := (p.SurplusValue*10000 + denom/2) / denom
		if p.ProfitRate != want {
			t.Errorf("period %d: ProfitRate = %d, want roundHalfUp = %d (s=%d, c+v=%d)",
				i+1, p.ProfitRate, want, p.SurplusValue, denom)
		}
		// s = s′·v/10000 with s′ = 100% → s == v.
		if p.SurplusValue != p.VariableCapital {
			t.Errorf("period %d: SurplusValue = %d, want %d (s′=100%% ⇒ s=v)",
				i+1, p.SurplusValue, p.VariableCapital)
		}
	}
}

// TestComputeFallingProfitRateLaw_Rates checks the single-period law against the
// §law fixture rates.
func TestComputeFallingProfitRateLaw_Rates(t *testing.T) {
	t.Parallel()
	for i, step := range lawSteps {
		law := ComputeFallingProfitRateLaw(step[0], step[1], 10000)
		if law.ProfitRate != wantLawRates[i] {
			t.Errorf("c=%d v=%d: ProfitRate = %d, want %d",
				step[0], step[1], law.ProfitRate, wantLawRates[i])
		}
		if law.SurplusValue != step[1] { // s′=100% ⇒ s=v
			t.Errorf("c=%d v=%d: SurplusValue = %d, want %d", step[0], step[1], law.SurplusValue, step[1])
		}
	}
}

// TestComputeFallingProfitRateLaw_RateBelowSurplusRate asserts the law's core
// inequality: for every c > 0 the rate of profit is strictly below the rate of
// surplus-value (because the denominator c+v exceeds v).
func TestComputeFallingProfitRateLaw_RateBelowSurplusRate(t *testing.T) {
	t.Parallel()
	const sRate = 10000 // 100%
	for _, c := range []int64{1, 50, 100, 200, 300, 400, 1000} {
		law := ComputeFallingProfitRateLaw(c, 100, sRate)
		if law.ProfitRate >= law.SurplusValueRate {
			t.Errorf("c=%d: ProfitRate %d >= SurplusValueRate %d, want strictly less",
				c, law.ProfitRate, law.SurplusValueRate)
		}
	}
}

// TestComputeFallingProfitRateLaw_ZeroConstantCapital is the c=0 boundary: with
// no constant capital the rate equals the rate of surplus-value (p′ = s/v = s′).
func TestComputeFallingProfitRateLaw_ZeroConstantCapital(t *testing.T) {
	t.Parallel()
	law := ComputeFallingProfitRateLaw(0, 100, 10000)
	if law.ProfitRate != 10000 {
		t.Errorf("c=0: ProfitRate = %d, want 10000 (p′=s′ when c=0)", law.ProfitRate)
	}
}

// TestProfitRate_DivideByZeroGuard ensures c+v=0 yields a 0 rate and does not
// panic.
func TestProfitRate_DivideByZeroGuard(t *testing.T) {
	t.Parallel()
	law := ComputeFallingProfitRateLaw(0, 0, 10000)
	if law.ProfitRate != 0 {
		t.Errorf("c+v=0: ProfitRate = %d, want 0", law.ProfitRate)
	}
	// Through the trajectory builder too.
	tr := BuildCompositionTrajectory("zero", 10000, [][2]int64{{0, 0}})
	if tr.Periods[0].ProfitRate != 0 {
		t.Errorf("c+v=0 trajectory: ProfitRate = %d, want 0", tr.Periods[0].ProfitRate)
	}
}

// TestComputeRateMassContradiction_MassPreserved is the §rate-mass fixture where
// the rate halves (20% → 10%) but the capital doubles, so the mass of profit is
// exactly preserved and MassChange == 0.
func TestComputeRateMassContradiction_MassPreserved(t *testing.T) {
	t.Parallel()
	rm := ComputeRateMassContradiction(1_000_000, 2000, 2_000_000, 1000)
	if rm.OldMass != 200000 {
		t.Errorf("OldMass = %d, want 200000", rm.OldMass)
	}
	if rm.NewMass != 200000 {
		t.Errorf("NewMass = %d, want 200000", rm.NewMass)
	}
	if rm.MassChange != 0 {
		t.Errorf("MassChange = %d, want 0 (mass preserved while rate fell)", rm.MassChange)
	}
}

// TestComputeRateMassContradiction_MassRoseWhileRateFell is the §rate-mass case
// where the rate halves but the capital triples: the mass of profit grows even
// as the rate falls — the law's internal contradiction.
func TestComputeRateMassContradiction_MassRoseWhileRateFell(t *testing.T) {
	t.Parallel()
	rm := ComputeRateMassContradiction(1_000_000, 2000, 3_000_000, 1000)
	if rm.OldMass != 200000 {
		t.Errorf("OldMass = %d, want 200000", rm.OldMass)
	}
	if rm.NewMass != 300000 {
		t.Errorf("NewMass = %d, want 300000", rm.NewMass)
	}
	if rm.MassChange <= 0 {
		t.Errorf("MassChange = %d, want > 0 (mass rose while rate fell)", rm.MassChange)
	}
	if rm.NewRate >= rm.OldRate {
		t.Errorf("NewRate %d should be below OldRate %d", rm.NewRate, rm.OldRate)
	}
}

// TestComputeRateMassContradiction_MassFormula verifies the invariant
// NewMass == NewC*NewRate/10000 across several states.
func TestComputeRateMassContradiction_MassFormula(t *testing.T) {
	t.Parallel()
	cases := []struct{ oldC, oldRate, newC, newRate int64 }{
		{1_000_000, 2000, 2_000_000, 1000},
		{1_000_000, 2000, 3_000_000, 1000},
		{500_000, 1500, 750_000, 1200},
		{0, 0, 0, 0},
	}
	for _, c := range cases {
		rm := ComputeRateMassContradiction(c.oldC, c.oldRate, c.newC, c.newRate)
		wantOld := c.oldC * c.oldRate / basisPoints
		wantNew := c.newC * c.newRate / basisPoints
		if rm.OldMass != wantOld {
			t.Errorf("%+v: OldMass = %d, want %d", c, rm.OldMass, wantOld)
		}
		if rm.NewMass != wantNew {
			t.Errorf("%+v: NewMass = %d, want %d", c, rm.NewMass, wantNew)
		}
		if rm.MassChange != rm.NewMass-rm.OldMass {
			t.Errorf("%+v: MassChange = %d, want %d", c, rm.MassChange, rm.NewMass-rm.OldMass)
		}
	}
}

// TestComputeAbsoluteOveraccumulation covers the limit-case predicate: a capital
// increment is absolute over-accumulation exactly when the additional
// surplus-value it sets in motion is zero or negative.
func TestComputeAbsoluteOveraccumulation(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		deltaC, deltaS int64
		wantAbsolute   bool
	}{
		{"no additional surplus-value", 1000, 0, true},
		{"some additional surplus-value", 1000, 500, false},
		{"negative additional surplus-value", 1000, -1, true},
		{"zero increment, zero surplus", 0, 0, true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			oa := ComputeAbsoluteOveraccumulation(tc.deltaC, tc.deltaS)
			if oa.IsAbsolute != tc.wantAbsolute {
				t.Errorf("IsAbsolute = %v, want %v", oa.IsAbsolute, tc.wantAbsolute)
			}
			// Invariant: IsAbsolute == (DeltaSurplusValue <= 0).
			if oa.IsAbsolute != (oa.DeltaSurplusValue <= 0) {
				t.Errorf("IsAbsolute %v != (DeltaSurplusValue<=0) for Δs=%d", oa.IsAbsolute, oa.DeltaSurplusValue)
			}
		})
	}
}

// TestCompositionTrajectory_Validate covers the happy path and both invariant
// boundaries: empty periods and negative magnitudes.
func TestCompositionTrajectory_Validate(t *testing.T) {
	t.Parallel()

	valid := BuildCompositionTrajectory("ok", 10000, lawSteps)
	if err := valid.Validate(); err != nil {
		t.Errorf("valid trajectory: err = %v, want nil", err)
	}

	empty := CompositionTrajectory{Periods: nil, SurplusValueRate: 10000}
	if err := empty.Validate(); !errors.Is(err, ErrEmptyTrajectory) {
		t.Errorf("empty periods: err = %v, want ErrEmptyTrajectory", err)
	}

	negRate := CompositionTrajectory{
		SurplusValueRate: -1,
		Periods:          []TrajectoryPeriod{{ConstantCapital: 50, VariableCapital: 100}},
	}
	if err := negRate.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("negative s′: err = %v, want ErrNegativeMagnitude", err)
	}

	negPeriod := CompositionTrajectory{
		SurplusValueRate: 10000,
		Periods:          []TrajectoryPeriod{{ConstantCapital: -1, VariableCapital: 100}},
	}
	if err := negPeriod.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("negative c: err = %v, want ErrNegativeMagnitude", err)
	}
}

// TestRateMassContradiction_Validate covers the happy path and the negative
// magnitude rejection on each field.
func TestRateMassContradiction_Validate(t *testing.T) {
	t.Parallel()

	valid := ComputeRateMassContradiction(1_000_000, 2000, 2_000_000, 1000)
	if err := valid.Validate(); err != nil {
		t.Errorf("valid contradiction: err = %v, want nil", err)
	}

	for _, bad := range []RateMassContradiction{
		{OldC: -1, OldRate: 2000, NewC: 1, NewRate: 1000},
		{OldC: 1, OldRate: -1, NewC: 1, NewRate: 1000},
		{OldC: 1, OldRate: 2000, NewC: -1, NewRate: 1000},
		{OldC: 1, OldRate: 2000, NewC: 1, NewRate: -1},
	} {
		if err := bad.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
			t.Errorf("%+v: err = %v, want ErrNegativeMagnitude", bad, err)
		}
	}
}

// TestDeriveProfitRates verifies the derived slice equals the per-period rates
// and is the single source of truth recomputed from Periods.
func TestDeriveProfitRates(t *testing.T) {
	t.Parallel()
	tr := BuildCompositionTrajectory("derive", 10000, lawSteps)
	derived := tr.DeriveProfitRates()
	if len(derived) != len(tr.Periods) {
		t.Fatalf("derived len = %d, want %d", len(derived), len(tr.Periods))
	}
	for i, p := range tr.Periods {
		if derived[i] != p.ProfitRate {
			t.Errorf("derived[%d] = %d, want %d", i, derived[i], p.ProfitRate)
		}
	}
}

// TestNewCompositionTrajectoryID_HexUnique asserts IDs are non-empty, valid hex,
// 24 chars (96 bits), and unique across a batch.
func TestNewCompositionTrajectoryID_HexUnique(t *testing.T) {
	t.Parallel()
	seen := make(map[CompositionTrajectoryID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewCompositionTrajectoryID()
		if id.IsZero() {
			t.Fatal("NewCompositionTrajectoryID returned the empty id")
		}
		if len(string(id)) != 24 {
			t.Fatalf("id %q has len %d, want 24 hex chars", id, len(string(id)))
		}
		if _, err := hex.DecodeString(string(id)); err != nil {
			t.Fatalf("id %q is not valid hex: %v", id, err)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id generated: %q", id)
		}
		seen[id] = struct{}{}
	}
}

// TestNewRateMassContradictionID_HexUnique mirrors the trajectory ID test for
// the rate-mass ID constructor.
func TestNewRateMassContradictionID_HexUnique(t *testing.T) {
	t.Parallel()
	seen := make(map[RateMassContradictionID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewRateMassContradictionID()
		if id.IsZero() {
			t.Fatal("NewRateMassContradictionID returned the empty id")
		}
		if len(string(id)) != 24 {
			t.Fatalf("id %q has len %d, want 24 hex chars", id, len(string(id)))
		}
		if _, err := hex.DecodeString(string(id)); err != nil {
			t.Fatalf("id %q is not valid hex: %v", id, err)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id generated: %q", id)
		}
		seen[id] = struct{}{}
	}
}
