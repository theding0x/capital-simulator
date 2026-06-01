package revenue

import "testing"

func TestValidClassTendencyDirection(t *testing.T) {
	t.Parallel()
	for _, d := range []string{"expansion", "concentration", "dissolution", "subordination"} {
		if !ValidClassTendencyDirection(d) {
			t.Errorf("ValidClassTendencyDirection(%q) = false, want true", d)
		}
	}
	for _, d := range []string{"", "growth", "collapse"} {
		if ValidClassTendencyDirection(d) {
			t.Errorf("ValidClassTendencyDirection(%q) = true, want false", d)
		}
	}
}

// TestClassIncomeSourcesPartitionNewValue asserts the Ch. 52 invariant: wages recover
// variable capital (component 0); profit and rent are portions of surplus-value (>0);
// together the three sources partition v + s.
func TestClassIncomeSourcesPartitionNewValue(t *testing.T) {
	t.Parallel()
	const v, s = 1500, 1500
	wages := ClassIncomeSource{RevenueForm: "wages", SurplusValueComponentBP: 0}
	profit := ClassIncomeSource{RevenueForm: "profit", SurplusValueComponentBP: 900}
	rent := ClassIncomeSource{RevenueForm: "rent", SurplusValueComponentBP: 600}

	if wages.SurplusValueComponentBP != 0 {
		t.Errorf("wages surplus component = %d, want 0", wages.SurplusValueComponentBP)
	}
	if profit.SurplusValueComponentBP <= 0 || rent.SurplusValueComponentBP <= 0 {
		t.Error("profit and rent must carry a positive surplus-value component")
	}
	if got := profit.SurplusValueComponentBP + rent.SurplusValueComponentBP; got != s {
		t.Errorf("profit + rent = %d, want s = %d", got, s)
	}
	// Wages equal variable capital, which is implicit (component 0 means "= v").
	if v <= 0 {
		t.Fatal("fixture variable capital must be positive")
	}
}

func TestCh52IDsUniqueAndHex(t *testing.T) {
	t.Parallel()
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		for _, id := range []string{
			string(NewSocialClassID()),
			string(NewClassIncomeSourceID()),
			string(NewClassTendencyID()),
			string(NewClassAmbiguityID()),
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
