package revenue

import "testing"

// Ch. 49 fixture: aggregate annual social capital c=6000, v=1500, s=1500.
// Total value 9000; new value 3000. Surplus-value partition: profit of
// enterprise 900 + interest 300 + ground-rent 300 = 1500.
func ch49Composition() AnnualValueComposition {
	return AnnualValueComposition{
		ConstantCapitalLabourMinutes: 6000,
		VariableCapitalLabourMinutes: 1500,
		SurplusValueLabourMinutes:    1500,
		TotalValueLabourMinutes:      ComputeTotalValue(6000, 1500, 1500),
		NewValueLabourMinutes:        ComputeNewValue(1500, 1500),
	}
}

func TestComputeTotalAndNewValue(t *testing.T) {
	t.Parallel()
	if got := ComputeTotalValue(6000, 1500, 1500); got != 9000 {
		t.Errorf("ComputeTotalValue = %d, want 9000", got)
	}
	if got := ComputeNewValue(1500, 1500); got != 3000 {
		t.Errorf("ComputeNewValue = %d, want 3000", got)
	}
}

func TestDeriveRevenueLimit(t *testing.T) {
	t.Parallel()
	comp := ch49Composition()
	comp.ID = NewAnnualValueCompositionID()
	rl := DeriveRevenueLimit(comp)
	if rl.MaxDistributableRevenueLM != 3000 {
		t.Errorf("MaxDistributableRevenueLM = %d, want 3000 (v+s)", rl.MaxDistributableRevenueLM)
	}
	if rl.ConstantCapitalReplacement != 6000 {
		t.Errorf("ConstantCapitalReplacement = %d, want 6000 (c)", rl.ConstantCapitalReplacement)
	}
	if rl.CompositionID != string(comp.ID) {
		t.Errorf("CompositionID = %q, want %q", rl.CompositionID, comp.ID)
	}
}

func TestValidateRevenueLimit(t *testing.T) {
	t.Parallel()
	comp := ch49Composition()

	good := SurplusValuePartition{
		TotalSurplusValueLM: 1500,
		ProfitEnterpriseLM:  900,
		InterestLM:          300,
		GroundRentLM:        300,
		ResidualSurplusLM:   0,
	}
	if err := ValidateRevenueLimit(comp, good); err != nil {
		t.Errorf("ValidateRevenueLimit(valid) = %v, want nil", err)
	}

	// Parts do not sum to the declared total.
	badSum := good
	badSum.ProfitEnterpriseLM = 1000
	if err := ValidateRevenueLimit(comp, badSum); err == nil {
		t.Error("ValidateRevenueLimit(parts != total) = nil, want error")
	}

	// Total exceeds the surplus-value produced (Smith's error).
	tooMuch := SurplusValuePartition{
		TotalSurplusValueLM: 1600,
		ProfitEnterpriseLM:  1000,
		InterestLM:          300,
		GroundRentLM:        300,
		ResidualSurplusLM:   0,
	}
	if err := ValidateRevenueLimit(comp, tooMuch); err == nil {
		t.Error("ValidateRevenueLimit(total > s) = nil, want error")
	}
}

func TestSumSurplusPartition(t *testing.T) {
	t.Parallel()
	got := SumSurplusPartition(SurplusValuePartition{
		ProfitEnterpriseLM: 900, InterestLM: 300, GroundRentLM: 300, ResidualSurplusLM: 0,
	})
	if got != 1500 {
		t.Errorf("SumSurplusPartition = %d, want 1500", got)
	}
}

func TestCh49IDsUniqueAndHex(t *testing.T) {
	t.Parallel()
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		for _, id := range []string{
			string(NewAnnualValueCompositionID()),
			string(NewRevenueLimitID()),
			string(NewSurplusValuePartitionID()),
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
