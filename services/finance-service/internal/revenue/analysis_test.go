package revenue

import (
	"errors"
	"strings"
	"testing"
)

// Ch. 49 fixture: aggregate annual social capital c=6000, v=1500, s=1500.
// Total value 9000; new value 3000. Surplus-value partition: profit of
// enterprise 900 + interest 300 + ground-rent 300 = 1500.
func ch49Composition() AnnualValueComposition {
	return AnnualValueComposition{
		ConstantCapitalLabourMinutes: 6000,
		VariableCapitalLabourMinutes: 1500,
		SurplusValueLabourMinutes:    1500,
		TotalValueLabourMinutes:      int64(ComputeTotalValue(6000, 1500, 1500)),
		NewValueLabourMinutes:        int64(ComputeNewValue(1500, 1500)),
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

func TestAnnualValueCompositionValidate(t *testing.T) {
	t.Parallel()
	good := ch49Composition()

	badTotal := good
	badTotal.TotalValueLabourMinutes = 3000 // Smith's error: v+s taken as the whole product

	badNew := good
	badNew.NewValueLabourMinutes = 9000 // c smuggled into new value

	negative := AnnualValueComposition{ConstantCapitalLabourMinutes: -1}

	tests := []struct {
		name    string
		comp    AnnualValueComposition
		wantErr bool
	}{
		{"Ch. 49 fixture 6000+1500+1500", good, false},
		{"zero composition", AnnualValueComposition{}, false},
		{"total != c+v+s", badTotal, true},
		{"new != v+s", badNew, true},
		{"negative component", negative, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.comp.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate() err = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr && !errors.Is(err, ErrCompositionNotConserved) {
				t.Errorf("err = %v, want ErrCompositionNotConserved", err)
			}
		})
	}
}

// Smith's adding-up error (Ch. 49): resolving the whole annual product, 9000,
// into wages + profit + rent omits the constant capital c = 6000 that must be
// replaced. Only v + s = 3000 is newly created value, so only 3000 can be revenue
// and the error is 6000.
func TestCheckRevenueResolutionSmithsError(t *testing.T) {
	t.Parallel()
	comp := ch49Composition()
	comp.ID = NewAnnualValueCompositionID()

	err := CheckRevenueResolution(comp, 9000)
	if err == nil {
		t.Fatal("CheckRevenueResolution(9000) = nil, want Smith's adding-up error")
	}
	var smith *SmithsAddingUpError
	if !errors.As(err, &smith) {
		t.Fatalf("err = %v (%T), want *SmithsAddingUpError", err, err)
	}
	if smith.CompositionID != string(comp.ID) {
		t.Errorf("CompositionID = %q, want %q", smith.CompositionID, comp.ID)
	}
	if smith.ClaimedRevenueLM != 9000 {
		t.Errorf("ClaimedRevenueLM = %d, want 9000 (the whole product)", smith.ClaimedRevenueLM)
	}
	if smith.NewValueLM != 3000 {
		t.Errorf("NewValueLM = %d, want 3000 (v+s)", smith.NewValueLM)
	}
	if smith.OmittedConstantCapitalLM != 6000 {
		t.Errorf("OmittedConstantCapitalLM = %d, want 6000 (the error)", smith.OmittedConstantCapitalLM)
	}
	// Resolving the whole product omits exactly the constant capital.
	if smith.OmittedConstantCapitalLM != comp.ConstantCapitalLabourMinutes {
		t.Errorf("OmittedConstantCapitalLM = %d, want c = %d", smith.OmittedConstantCapitalLM, comp.ConstantCapitalLabourMinutes)
	}
	if msg := smith.Error(); !strings.Contains(msg, "6000") {
		t.Errorf("Error() = %q, want it to state the 6000 error", msg)
	}
}

func TestCheckRevenueResolutionBoundaries(t *testing.T) {
	t.Parallel()
	comp := ch49Composition() // c=6000, v=1500, s=1500; new value 3000

	tests := []struct {
		name        string
		claimed     int64
		wantOmitted int64 // 0 means no Smith's error expected
	}{
		{"nothing claimed as revenue", 0, 0},
		{"wages alone (v)", 1500, 0},
		{"surplus-value alone (s)", 1500, 0},
		{"exactly the new value v+s", 3000, 0},
		{"one labour-minute past v+s", 3001, 1},
		{"v+s plus half of c", 6000, 3000},
		{"the whole product c+v+s", 9000, 6000},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := CheckRevenueResolution(comp, tc.claimed)
			if tc.wantOmitted == 0 {
				if err != nil {
					t.Fatalf("CheckRevenueResolution(%d) = %v, want nil", tc.claimed, err)
				}
				return
			}
			var smith *SmithsAddingUpError
			if !errors.As(err, &smith) {
				t.Fatalf("CheckRevenueResolution(%d) = %v, want *SmithsAddingUpError", tc.claimed, err)
			}
			if smith.OmittedConstantCapitalLM != tc.wantOmitted {
				t.Errorf("OmittedConstantCapitalLM = %d, want %d", smith.OmittedConstantCapitalLM, tc.wantOmitted)
			}
			// The error is always the claim less the new value: ClaimedRevenueLM
			// = NewValueLM + OmittedConstantCapitalLM.
			if smith.ClaimedRevenueLM != smith.NewValueLM+smith.OmittedConstantCapitalLM {
				t.Errorf("claimed %d != new %d + omitted %d", smith.ClaimedRevenueLM, smith.NewValueLM, smith.OmittedConstantCapitalLM)
			}
		})
	}
}

func TestCheckRevenueResolutionRejectsBadInput(t *testing.T) {
	t.Parallel()

	t.Run("negative claim is not Smith's error", func(t *testing.T) {
		t.Parallel()
		err := CheckRevenueResolution(ch49Composition(), -1)
		if err == nil {
			t.Fatal("CheckRevenueResolution(-1) = nil, want error")
		}
		var smith *SmithsAddingUpError
		if errors.As(err, &smith) {
			t.Errorf("negative claim reported as Smith's error: %v", err)
		}
		if !errors.Is(err, ErrNegativeRevenueClaim) {
			t.Errorf("err = %v, want ErrNegativeRevenueClaim", err)
		}
	})

	t.Run("non-conserving composition is rejected before comparing", func(t *testing.T) {
		t.Parallel()
		drifted := ch49Composition()
		drifted.NewValueLabourMinutes = 9000 // c smuggled into new value
		err := CheckRevenueResolution(drifted, 9000)
		if !errors.Is(err, ErrCompositionNotConserved) {
			t.Errorf("err = %v, want ErrCompositionNotConserved", err)
		}
	})
}
