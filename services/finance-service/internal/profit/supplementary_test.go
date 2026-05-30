package profit

import (
	"errors"
	"testing"
)

// Marx's two-enterprise illustration: equal s = £1,000 on equal v = £1,000, the
// rate of profit set only by the organic composition. A (c = £9,000) shows 10%;
// B (c = £11,000) only 8 1/3%. The £2,000 of constant capital is the whole
// difference: 1000 bp − 833 bp = 167 bp.
func TestOrganicCompositionEffect(t *testing.T) {
	t.Parallel()

	e := OrganicCompositionEffect{SA: 1000, VA: 1000, CA: 9000, SB: 1000, VB: 1000, CB: 11000}
	if got := e.ProfitRateA(); got != 1000 {
		t.Errorf("p'_A: got %d bp, want 1000", got)
	}
	if got := e.ProfitRateB(); got != 833 {
		t.Errorf("p'_B: got %d bp, want 833", got)
	}
	if got := e.RateDifference(); got != 167 {
		t.Errorf("rate difference: got %d bp, want 167", got)
	}
	// The difference is symmetric: swapping the two enterprises yields the same gap.
	swapped := OrganicCompositionEffect{SA: 1000, VA: 1000, CA: 11000, SB: 1000, VB: 1000, CB: 9000}
	if swapped.RateDifference() != e.RateDifference() {
		t.Errorf("rate difference should be symmetric: %d vs %d", swapped.RateDifference(), e.RateDifference())
	}
}

func TestComputeCompositionEffectAnalysis(t *testing.T) {
	t.Parallel()

	a := ComputeCompositionEffectAnalysis(OrganicCompositionEffect{SA: 1000, VA: 1000, CA: 9000, SB: 1000, VB: 1000, CB: 11000})
	if a.ProfitRateA != 1000 || a.ProfitRateB != 833 || a.RateDifference != 167 {
		t.Errorf("got A=%d B=%d diff=%d, want 1000/833/167", a.ProfitRateA, a.ProfitRateB, a.RateDifference)
	}
	if a.ID != "" || !a.CreatedAt.IsZero() {
		t.Errorf("Compute must not assign ID/timestamp: id=%q createdAt=%v", a.ID, a.CreatedAt)
	}
}

func TestOrganicCompositionEffectValidate(t *testing.T) {
	t.Parallel()

	if err := (OrganicCompositionEffect{SA: 1000, VA: 1000, CA: 9000, SB: 1000, VB: 1000, CB: 11000}).Validate(); err != nil {
		t.Errorf("valid effect: %v", err)
	}
	if err := (OrganicCompositionEffect{CA: -1}).Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("negative magnitude: got %v, want ErrNegativeMagnitude", err)
	}
}

// Rodbertus refutation, money-value change: capital £100 yielding profit £20 (rate
// 20%). If gold falls by half, the same capital is reckoned at £200 and the profit
// at £40 — but the rate of profit is unmoved at 20%.
func TestCapitalMagnitudeChangeMoneyValue(t *testing.T) {
	t.Parallel()

	m := CapitalMagnitudeChange{Kind: KindMoneyValueChange, OriginalCapital: 100, OriginalProfit: 20, Factor: 200}
	if got := m.OldRate(); got != 2000 {
		t.Errorf("old rate: got %d bp, want 2000", got)
	}
	if got := m.NewRate(); got != 2000 {
		t.Errorf("new rate: got %d bp, want 2000 (nominal change leaves the rate untouched)", got)
	}
	if !m.RateUnchanged() {
		t.Error("a money-value change must leave the rate of profit unchanged")
	}
}

// A proportional change in the magnitude of capital at unchanged composition also
// leaves the rate untouched.
func TestCapitalMagnitudeChangeProportional(t *testing.T) {
	t.Parallel()

	m := CapitalMagnitudeChange{Kind: KindProportionalChange, OriginalCapital: 3000, OriginalProfit: 600, Factor: 250}
	if m.NewRate() != m.OldRate() {
		t.Errorf("proportional change: new %d != old %d", m.NewRate(), m.OldRate())
	}
	if !m.RateUnchanged() {
		t.Error("a proportional change must leave the rate of profit unchanged")
	}
}

// An organic change — the magnitude of capital altered with the composition, the
// profit not scaling with it — does move the rate. Doubling the capital with the
// same profit halves the rate: 20% → 10%.
func TestCapitalMagnitudeChangeOrganic(t *testing.T) {
	t.Parallel()

	m := CapitalMagnitudeChange{Kind: KindOrganicChange, OriginalCapital: 100, OriginalProfit: 20, Factor: 200}
	if got := m.OldRate(); got != 2000 {
		t.Errorf("old rate: got %d bp, want 2000", got)
	}
	if got := m.NewRate(); got != 1000 {
		t.Errorf("new rate: got %d bp, want 1000", got)
	}
	if m.RateUnchanged() {
		t.Error("an organic change must move the rate of profit")
	}
}

func TestComputeMagnitudeChangeAnalysis(t *testing.T) {
	t.Parallel()

	a := ComputeMagnitudeChangeAnalysis(CapitalMagnitudeChange{Kind: KindMoneyValueChange, OriginalCapital: 100, OriginalProfit: 20, Factor: 200})
	if a.OldRate != 2000 || a.NewRate != 2000 || !a.RateUnchanged {
		t.Errorf("got old=%d new=%d unchanged=%v, want 2000/2000/true", a.OldRate, a.NewRate, a.RateUnchanged)
	}
	if a.ID != "" || !a.CreatedAt.IsZero() {
		t.Errorf("Compute must not assign ID/timestamp: id=%q createdAt=%v", a.ID, a.CreatedAt)
	}
}

func TestMagnitudeChangeKindValid(t *testing.T) {
	t.Parallel()

	for _, k := range []MagnitudeChangeKind{KindMoneyValueChange, KindProportionalChange, KindOrganicChange} {
		if !k.Valid() {
			t.Errorf("kind %q should be valid", k)
		}
	}
	if MagnitudeChangeKind("nominal").Valid() {
		t.Error(`"nominal" should be invalid`)
	}
}

func TestCapitalMagnitudeChangeValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		m    CapitalMagnitudeChange
		want error
	}{
		{
			name: "valid",
			m:    CapitalMagnitudeChange{Kind: KindMoneyValueChange, OriginalCapital: 100, OriginalProfit: 20, Factor: 200},
			want: nil,
		},
		{
			name: "negative magnitude",
			m:    CapitalMagnitudeChange{Kind: KindMoneyValueChange, OriginalCapital: -1, Factor: 100},
			want: ErrNegativeMagnitude,
		},
		{
			name: "non-positive factor",
			m:    CapitalMagnitudeChange{Kind: KindMoneyValueChange, OriginalCapital: 100, Factor: 0},
			want: ErrNonPositivePriceFactor,
		},
		{
			name: "unknown kind",
			m:    CapitalMagnitudeChange{Kind: "bogus", OriginalCapital: 100, Factor: 100},
			want: ErrUnknownMagnitudeKind,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.m.Validate(); !errors.Is(err, tc.want) {
				t.Errorf("Validate() = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestProfitSurfaceAppearance(t *testing.T) {
	t.Parallel()

	a := NewProfitSurfaceAppearance(2000, 10000)
	if a.ProfitRate != 2000 || a.SurplusValueRate != 10000 {
		t.Errorf("got p'=%d s'=%d, want 2000/10000", a.ProfitRate, a.SurplusValueRate)
	}
	if a.ApparentSource == "" || a.RealSource == "" {
		t.Error("surface appearance must name both the apparent and real source")
	}
}

func TestComputePartISummary(t *testing.T) {
	t.Parallel()

	// Structural invariant: the summary always names the determinants, and carries
	// the analyses it is given (here one Ch. 2 analysis).
	analyses := []ProfitRateAnalysis{ComputeProfitRateAnalysis(9000, 1000, 1000)}
	s := ComputePartISummary(analyses)
	if len(s.Determinants) == 0 {
		t.Error("summary must name the determinants of the rate of profit")
	}
	if len(s.Analyses) < 1 {
		t.Errorf("summary analyses: got %d, want >= 1", len(s.Analyses))
	}
	if s.ID.IsZero() {
		t.Error("summary must carry an ID")
	}

	// A nil analyses slice is normalised to a non-nil empty slice.
	empty := ComputePartISummary(nil)
	if empty.Analyses == nil {
		t.Error("analyses should be a non-nil slice, not nil")
	}
}
