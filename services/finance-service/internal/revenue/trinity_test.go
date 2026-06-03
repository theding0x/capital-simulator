package revenue

import (
	"errors"
	"testing"
)

// trinityFixture builds the Ch. 48 seed trinity. The capitalist farms rented
// land: capital's *apparent* revenue is the average profit 20 (interest 5 +
// profit of enterprise 15), but of that 20 the ground-rent 3 is paid out to the
// landlord, so capital's *actual* surplus share is 17 (profit of enterprise 12 +
// interest 5). Land's actual share is the rent 3. Profit + rent = 17 + 3 = 20 =
// the single surplus-value s — partitioned, not double-counted. Wages 80
// (= variable capital v) recover v and represent no surplus-value. Apparent
// revenue 20 + 3 + 80 = 103 overshoots newly-created value v + s = 100 by exactly
// the rent 3 — that overshoot is the fetish, not extra wealth.
func trinityFixture() []RevenueStream {
	return []RevenueStream{
		{Source: RevenueSourceCapital, ApparentRevenueBP: 20, ActualSourceBP: 17, IsFetishised: true},
		{Source: RevenueSourceLand, ApparentRevenueBP: 3, ActualSourceBP: 3, IsFetishised: true},
		{Source: RevenueSourceLabour, ApparentRevenueBP: 80, ActualSourceBP: 0, IsFetishised: true},
	}
}

func TestRevenueSourceKindIsValid(t *testing.T) {
	t.Parallel()
	valid := []RevenueSourceKind{RevenueSourceCapital, RevenueSourceLand, RevenueSourceLabour}
	for _, k := range valid {
		if !k.IsValid() {
			t.Errorf("RevenueSourceKind(%d).IsValid() = false, want true", k)
		}
	}
	for _, k := range []RevenueSourceKind{0, 4, -1} {
		if k.IsValid() {
			t.Errorf("RevenueSourceKind(%d).IsValid() = true, want false", k)
		}
	}
}

func TestSumApparentRevenueBP(t *testing.T) {
	t.Parallel()
	got := SumApparentRevenueBP(trinityFixture())
	if want := int64(103); got != want {
		t.Errorf("SumApparentRevenueBP = %d, want %d", got, want)
	}
	if got := SumApparentRevenueBP(nil); got != 0 {
		t.Errorf("SumApparentRevenueBP(nil) = %d, want 0", got)
	}
}

func TestSumActualSourceBP(t *testing.T) {
	t.Parallel()
	got := SumActualSourceBP(trinityFixture())
	if want := int64(20); got != want {
		t.Errorf("SumActualSourceBP = %d, want %d (profit 17 + rent 3 = one s, no double-count)", got, want)
	}
}

func TestVariableCapitalBP(t *testing.T) {
	t.Parallel()
	// Wages (the labour stream's apparent revenue) are the variable capital.
	if got, want := VariableCapitalBP(trinityFixture()), int64(80); got != want {
		t.Errorf("VariableCapitalBP = %d, want %d", got, want)
	}
}

// TestTrinityInvariants asserts the chapter's stated invariants on the seed fixture.
func TestTrinityInvariants(t *testing.T) {
	t.Parallel()
	streams := trinityFixture()
	tf := TrinityFormula{
		TotalApparentRevenueBP: SumApparentRevenueBP(streams),
		TotalSurplusValueBP:    SumActualSourceBP(streams),
	}

	// Invariant 1: total apparent revenue is the sum of the streams' apparent yields.
	if tf.TotalApparentRevenueBP != SumApparentRevenueBP(streams) {
		t.Errorf("TotalApparentRevenueBP = %d, want %d", tf.TotalApparentRevenueBP, SumApparentRevenueBP(streams))
	}
	// Invariant 2: total surplus-value is partitioned, not created — sum of actual sources.
	if tf.TotalSurplusValueBP != SumActualSourceBP(streams) {
		t.Errorf("TotalSurplusValueBP = %d, want %d", tf.TotalSurplusValueBP, SumActualSourceBP(streams))
	}
	// Invariant 3: the apparent revenues *do* exceed newly-created value (v + s),
	// and the overshoot equals exactly the ground-rent — the trinity's fetish
	// (rent appears as income with no value behind it). The old assertion
	// (apparent <= s + v) only held because s was double-counted as 23; with the
	// correct s = 20 it would FAIL, masking the very mystification the chapter
	// exists to demolish. We surface the gap, not assert it away.
	if got, want := FetishOvershootBP(streams), RentBP(streams); got != want {
		t.Errorf("fetish overshoot = %d, want it to equal ground-rent %d", got, want)
	}
	if got := FetishOvershootBP(streams); got != 3 {
		t.Errorf("fetish overshoot = %d, want 3 (the doubly-counted ground-rent)", got)
	}
	// Invariant 4: wages recover only variable capital — no surplus-value.
	for _, s := range streams {
		if s.Source == RevenueSourceLabour && s.ActualSourceBP != 0 {
			t.Errorf("labour stream ActualSourceBP = %d, want 0 (wages are v, not s)", s.ActualSourceBP)
		}
	}
	// Invariant 5: every surface stream is fetishised.
	for _, s := range streams {
		if !s.IsFetishised {
			t.Errorf("stream source %d IsFetishised = false, want true", s.Source)
		}
	}
}

// TestProfitRentPartitionSurplus asserts profit and rent partition one s,
// rather than rent being added on top of profit (the #377 defect).
func TestProfitRentPartitionSurplus(t *testing.T) {
	t.Parallel()
	streams := trinityFixture()
	if got, want := ProfitBP(streams), int64(17); got != want {
		t.Errorf("ProfitBP = %d, want 17 (profit of enterprise 12 + interest 5)", got)
	}
	if got, want := RentBP(streams), int64(3); got != want {
		t.Errorf("RentBP = %d, want 3", got)
	}
	// profit + rent == the single surplus-value s (nothing added on top).
	if got, want := ProfitBP(streams)+RentBP(streams), SumActualSourceBP(streams); got != want {
		t.Errorf("profit+rent = %d, want = s %d", got, want)
	}
}

// TestNewValueAndConservation asserts the load-bearing conservation invariant
// wages + profit + rent ≙ v + s on the faithful fixture.
func TestNewValueAndConservation(t *testing.T) {
	t.Parallel()
	streams := trinityFixture()
	if got, want := VariableCapitalBP(streams)+ProfitBP(streams)+RentBP(streams), NewValueBP(streams); got != want {
		t.Errorf("wages+profit+rent = %d, want v+s = %d", got, want)
	}
	if got := NewValueBP(streams); got != 100 {
		t.Errorf("NewValueBP = %d, want 100 (v 80 + s 20)", got)
	}
	if err := CheckConservation(streams); err != nil {
		t.Errorf("CheckConservation(faithful trinity) = %v, want nil", err)
	}
}

// TestCheckConservation_DoubleCountRejected pins the #377 regression: the old
// fixture set capital actual = 20, double-counting the rent so s = 23 and the
// apparent overshoot collapsed to 0. That trinity must now be rejected.
func TestCheckConservation_DoubleCountRejected(t *testing.T) {
	t.Parallel()
	bad := []RevenueStream{
		{Source: RevenueSourceCapital, ApparentRevenueBP: 20, ActualSourceBP: 20, IsFetishised: true},
		{Source: RevenueSourceLand, ApparentRevenueBP: 3, ActualSourceBP: 3, IsFetishised: true},
		{Source: RevenueSourceLabour, ApparentRevenueBP: 80, ActualSourceBP: 0, IsFetishised: true},
	}
	if err := CheckConservation(bad); !errors.Is(err, ErrSurplusNotConserved) {
		t.Errorf("CheckConservation(double-count) = %v, want ErrSurplusNotConserved", err)
	}
}

func TestNewIDsUniqueAndHex(t *testing.T) {
	t.Parallel()
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		for _, id := range []string{
			string(NewRevenueStreamID()),
			string(NewTrinityFormulaID()),
			string(NewRevenueFetishFormID()),
		} {
			if len(id) != 24 {
				t.Fatalf("id %q length = %d, want 24", id, len(id))
			}
			if seen[id] {
				t.Fatalf("duplicate id %q", id)
			}
			seen[id] = true
		}
	}
}
