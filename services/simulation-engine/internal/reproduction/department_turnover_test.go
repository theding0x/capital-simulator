package reproduction_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// Marx's canonical Ch. 20 per-period magnitudes (period 1865), parameterised by
// turnover number so the annualised-settlement behaviour can be exercised.
func deptI1865(turnoverBP int64) reproduction.DepartmentalCapital {
	return reproduction.DepartmentalCapital{
		Department:       reproduction.DepartmentI,
		ConstantPence:    4000,
		VariablePence:    1000,
		SurplusPence:     1000,
		TotalPence:       6000,
		TurnoverNumberBP: turnoverBP,
	}
}

func deptII1865(turnoverBP int64) reproduction.DepartmentalCapital {
	return reproduction.DepartmentalCapital{
		Department:       reproduction.DepartmentII,
		ConstantPence:    2000,
		VariablePence:    500,
		SurplusPence:     500,
		TotalPence:       3000,
		TurnoverNumberBP: turnoverBP,
	}
}

// TestAnnualisedSettlement_TurnoverTwoToOne is the issue #221 demonstration:
// Dept I turns over twice a year (20000 bp) and Dept II once (10000 bp). Per
// period the keystone balances (I(v+s)=2000 == II(c)=2000), but annualised it
// does not — I's flow doubles to 4000 while II's stays at 2000.
func TestAnnualisedSettlement_TurnoverTwoToOne(t *testing.T) {
	t.Parallel()
	deptI := deptI1865(20000)
	deptII := deptII1865(10000)

	if !reproduction.ComputeIsBalanced(deptI, deptII) {
		t.Fatal("raw keystone should balance: I(v+s)=2000 == II(c)=2000")
	}

	var res reproduction.BalanceCheckResult
	reproduction.ApplyAnnualisedSettlement(&res, &deptI, &deptII)

	if res.DeptIVplusSAnnualisedPence != 4000 {
		t.Fatalf("annualised I(v+s): got %d, want 4000", res.DeptIVplusSAnnualisedPence)
	}
	if res.DeptIIConstantAnnualisedPence != 2000 {
		t.Fatalf("annualised II(c): got %d, want 2000", res.DeptIIConstantAnnualisedPence)
	}
	if res.AnnualisedDeficitPence != 2000 {
		t.Fatalf("annualised deficit: got %d, want 2000", res.AnnualisedDeficitPence)
	}
	if res.AnnualisedBalanced {
		t.Fatal("annualised settlement must be unbalanced when turnovers differ")
	}
}

// TestAnnualisedSettlement_EqualTurnoverMatchesRaw: equal turnover ⇒ annualised
// equals raw, so a per-period-balanced scheme is annually balanced too.
func TestAnnualisedSettlement_EqualTurnoverMatchesRaw(t *testing.T) {
	t.Parallel()
	deptI := deptI1865(10000)
	deptII := deptII1865(10000)

	var res reproduction.BalanceCheckResult
	res.DeptIVplusSPence = deptI.VplusSPence()
	res.DeptIIConstantPence = deptII.ConstantPence
	reproduction.ApplyAnnualisedSettlement(&res, &deptI, &deptII)

	if res.DeptIVplusSAnnualisedPence != res.DeptIVplusSPence {
		t.Fatalf("equal turnover: annualised %d should equal raw %d",
			res.DeptIVplusSAnnualisedPence, res.DeptIVplusSPence)
	}
	if !res.AnnualisedBalanced {
		t.Fatal("equal turnover with balanced raw should be annually balanced")
	}
}

// TestAnnualisedHelpers_UnsetTreatedAsAnnual: a zero turnover is treated as
// annual (10000 bp), so annualised equals raw.
func TestAnnualisedHelpers_UnsetTreatedAsAnnual(t *testing.T) {
	t.Parallel()
	d := deptI1865(0)
	if d.AnnualisedVplusSPence() != d.VplusSPence() {
		t.Fatalf("unset turnover should be annual: annualised %d != raw %d",
			d.AnnualisedVplusSPence(), d.VplusSPence())
	}
	if d.AnnualisedConstantPence() != d.ConstantPence {
		t.Fatalf("unset turnover constant: got %d, want %d",
			d.AnnualisedConstantPence(), d.ConstantPence)
	}
}

// TestValidate_NegativeTurnover rejects a negative turnover number.
func TestValidate_NegativeTurnover(t *testing.T) {
	t.Parallel()
	d := deptI1865(-1)
	if err := d.Validate(); err == nil {
		t.Fatal("expected ErrNegativeTurnover for negative turnover_number_bp")
	}
}
