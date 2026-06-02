package reproduction_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// TestBuildReproductionFromAnnualRates covers the Ch.16→20 bridge (issue #222):
// surplus is derived from each department's annual rate, and the Marx fixture
// (100% rates) yields a balanced keystone I(v+s)=2000 == II(c)=2000.
func TestBuildReproductionFromAnnualRates(t *testing.T) {
	t.Parallel()
	scheme := reproduction.BuildReproductionFromAnnualRates("1865", []reproduction.AnnualRateDepartment{
		{Department: reproduction.DepartmentI, ConstantPence: 4000, VariablePence: 1000, AnnualSurplusBP: 10000},
		{Department: reproduction.DepartmentII, ConstantPence: 2000, VariablePence: 500, AnnualSurplusBP: 10000},
	})
	if scheme.DepartmentI == nil || scheme.DepartmentII == nil {
		t.Fatal("both departments should be set")
	}
	if scheme.DepartmentI.SurplusPence != 1000 {
		t.Fatalf("dept I surplus: got %d, want 1000 (1000 × 10000/10000)", scheme.DepartmentI.SurplusPence)
	}
	if scheme.DepartmentII.SurplusPence != 500 {
		t.Fatalf("dept II surplus: got %d, want 500", scheme.DepartmentII.SurplusPence)
	}
	if scheme.DepartmentI.TotalPence != 6000 {
		t.Fatalf("dept I total: got %d, want 6000 (c+v+s)", scheme.DepartmentI.TotalPence)
	}
	if !scheme.IsBalanced {
		t.Fatal("expected balanced: I(v+s)=2000 == II(c)=2000")
	}
}

// TestBuildReproductionFromAnnualRates_Imbalanced: a 200% Dept I rate produces
// more surplus than Dept II's constant capital can absorb, breaking the keystone.
func TestBuildReproductionFromAnnualRates_Imbalanced(t *testing.T) {
	t.Parallel()
	scheme := reproduction.BuildReproductionFromAnnualRates("1865", []reproduction.AnnualRateDepartment{
		{Department: reproduction.DepartmentI, ConstantPence: 4000, VariablePence: 1000, AnnualSurplusBP: 20000}, // s=2000
		{Department: reproduction.DepartmentII, ConstantPence: 2000, VariablePence: 500, AnnualSurplusBP: 10000},
	})
	// I(v+s) = 1000 + 2000 = 3000 != II(c) = 2000 → imbalanced.
	if scheme.IsBalanced {
		t.Fatal("expected imbalanced when Dept I annual rate is 200%")
	}
}
