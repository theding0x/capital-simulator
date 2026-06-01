package revenue

import "testing"

// trinityFixture builds the Ch. 48 seed trinity: average profit 20 (interest 5 +
// profit of enterprise 15), ground-rent 3 deducted as a further claim on surplus-
// value, wages 80 (= variable capital). Total surplus-value behind the formula is
// 23 (20 profit + 3 rent); wages recover v and represent no surplus-value.
func trinityFixture() []RevenueStream {
	return []RevenueStream{
		{Source: RevenueSourceCapital, ApparentRevenueBP: 20, ActualSourceBP: 20, IsFetishised: true},
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
	if want := int64(23); got != want {
		t.Errorf("SumActualSourceBP = %d, want %d", got, want)
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
	// Invariant 3: apparent revenues cannot exceed newly created value (s + v).
	if maxNew := tf.TotalSurplusValueBP + VariableCapitalBP(streams); tf.TotalApparentRevenueBP > maxNew {
		t.Errorf("TotalApparentRevenueBP %d exceeds new value s+v = %d", tf.TotalApparentRevenueBP, maxNew)
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
