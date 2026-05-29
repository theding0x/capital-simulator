package profit

import "testing"

// Marx's §I.1 worked figures: 80c + 20v with s' = 100% gives p' = 20%, and the
// same composition with s' = 50% gives p' = 10%. The struct's C field carries
// the constant capital (80); the denominator of p' = s'(v/C) is C + V = 100.
func TestProfitRateFormula_ProfitRate_Marx(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		f    ProfitRateFormula
		want int64
	}{
		{"s'=100%", ProfitRateFormula{C: 80, V: 20, SRate: 10000}, 2000},
		{"s'=50%", ProfitRateFormula{C: 80, V: 20, SRate: 5000}, 1000},
		{"§I.3 C=120", ProfitRateFormula{C: 100, V: 20, SRate: 10000}, 1666},
		{"§I.3 C=80", ProfitRateFormula{C: 60, V: 20, SRate: 10000}, 2500},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.f.ProfitRate(); got != tc.want {
				t.Errorf("ProfitRate() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestProfitRateFormula_TotalCapital(t *testing.T) {
	t.Parallel()
	f := ProfitRateFormula{C: 80, V: 20, SRate: 10000}
	if got := f.TotalCapital(); got != 100 {
		t.Errorf("TotalCapital() = %d, want 100", got)
	}
}

// Invariant [§intro]: p' = s' · v / C. The rate is zero with no variable
// capital, since variable capital alone produces the surplus the rate measures.
func TestProfitRateFormula_Invariant(t *testing.T) {
	t.Parallel()
	f := ProfitRateFormula{C: 84, V: 16, SRate: 12500}
	total := int64(f.C) + int64(f.V)
	want := int64(f.SRate) * int64(f.V) / total
	if got := f.ProfitRate(); got != want {
		t.Errorf("ProfitRate() = %d, want %d (= s'·v/C)", got, want)
	}
}

func TestProfitRateFormula_ZeroVariable(t *testing.T) {
	t.Parallel()
	if got := (ProfitRateFormula{C: 100, V: 0, SRate: 10000}).ProfitRate(); got != 0 {
		t.Errorf("ProfitRate() with v=0 = %d, want 0", got)
	}
}

func TestComputeCompositionRatio(t *testing.T) {
	t.Parallel()
	r := ComputeCompositionRatio(80, 20)
	if r.Total != 100 || r.BasisPoints != 2000 {
		t.Errorf("v/C of 80c+20v = {total %d, %d bp}, want {100, 2000}", r.Total, r.BasisPoints)
	}
	if zero := ComputeCompositionRatio(0, 0); zero.BasisPoints != 0 {
		t.Errorf("v/C of zero capital = %d bp, want 0", zero.BasisPoints)
	}
}

// §I.1: two capitals with equal s' have rates of profit in the same proportion
// as their value-compositions. Marx's figures: 12,000c + 3,000v → p' = 20%;
// 13,000c + 2,000v → p' = 13⅓%, both at s' = 100%.
func TestCompareProfitRates_EqualSurplus(t *testing.T) {
	t.Parallel()
	c := CompareProfitRates(
		ProfitRateFormula{C: 12000, V: 3000, SRate: 10000},
		ProfitRateFormula{C: 13000, V: 2000, SRate: 10000},
	)
	if !c.EqualSurplusRate {
		t.Fatal("EqualSurplusRate = false, want true")
	}
	if c.FirstProfitRate != 2000 || c.SecondProfitRate != 1333 {
		t.Errorf("rates = %d, %d, want 2000, 1333", c.FirstProfitRate, c.SecondProfitRate)
	}
	// p'₁ : p'₂ = v₁/C₁ : v₂/C₂ — cross-multiply to avoid division.
	lhs := c.FirstProfitRate * c.SecondComposition.BasisPoints
	rhs := c.SecondProfitRate * c.FirstComposition.BasisPoints
	if lhs != rhs {
		t.Errorf("p'₁:p'₂ (%d:%d) not proportional to v₁/C₁:v₂/C₂ (%d:%d)",
			c.FirstProfitRate, c.SecondProfitRate,
			c.FirstComposition.BasisPoints, c.SecondComposition.BasisPoints)
	}
}

func TestCompareProfitRates_UnequalSurplus(t *testing.T) {
	t.Parallel()
	c := CompareProfitRates(
		ProfitRateFormula{C: 80, V: 20, SRate: 10000},
		ProfitRateFormula{C: 80, V: 20, SRate: 5000},
	)
	if c.EqualSurplusRate {
		t.Error("EqualSurplusRate = true, want false for differing s'")
	}
	if c.FirstProfitRate != 2000 || c.SecondProfitRate != 1000 {
		t.Errorf("rates = %d, %d, want 2000, 1000", c.FirstProfitRate, c.SecondProfitRate)
	}
}

// §I.3: with s' constant, raising the value-composition (lowering v/C) lowers
// p'; lowering it raises p'. Invariants [§I.4a, §I.4c].
func TestComputeVariationAnalysis_SConstant(t *testing.T) {
	t.Parallel()
	base := ProfitRateFormula{C: 80, V: 20, SRate: 10000} // C=100, p'=20%

	// v/C falls (c rises 80→100): p' falls 20% → 16⅔%.
	falling := ComputeVariationAnalysis(base, ProfitRateFormula{C: 100, V: 20, SRate: 10000})
	if falling.Case != CaseSConstantVCVariable {
		t.Errorf("case = %q, want %q", falling.Case, CaseSConstantVCVariable)
	}
	if !(falling.NewProfitRate < falling.OldProfitRate) {
		t.Errorf("v/C falling: NewProfitRate %d should be < OldProfitRate %d",
			falling.NewProfitRate, falling.OldProfitRate)
	}
	if falling.OldProfitRate != 2000 || falling.NewProfitRate != 1666 {
		t.Errorf("rates = %d → %d, want 2000 → 1666", falling.OldProfitRate, falling.NewProfitRate)
	}

	// v/C rises (c falls 80→60): p' rises 20% → 25%.
	rising := ComputeVariationAnalysis(base, ProfitRateFormula{C: 60, V: 20, SRate: 10000})
	if rising.Case != CaseSConstantVCVariable {
		t.Errorf("case = %q, want %q", rising.Case, CaseSConstantVCVariable)
	}
	if !(rising.NewProfitRate > rising.OldProfitRate) {
		t.Errorf("v/C rising: NewProfitRate %d should be > OldProfitRate %d",
			rising.NewProfitRate, rising.OldProfitRate)
	}
	if rising.NewProfitRate != 2500 {
		t.Errorf("NewProfitRate = %d, want 2500", rising.NewProfitRate)
	}
	if rising.ProportionalChange {
		t.Error("ProportionalChange = true, want false when s' is held constant")
	}
}

// §II.1: with v/C constant, p' moves in the same proportion as s'. Marx's
// figures: 80c+20v at s'=100% → p'=20%; 160c+40v at s'=50% → p'=10%; the
// composition v/C = 20% is unchanged and 100:50 = 20:10.
func TestComputeVariationAnalysis_SVariableProportional(t *testing.T) {
	t.Parallel()
	a := ComputeVariationAnalysis(
		ProfitRateFormula{C: 80, V: 20, SRate: 10000},
		ProfitRateFormula{C: 160, V: 40, SRate: 5000},
	)
	if a.Case != CaseSVariableVCConstant {
		t.Errorf("case = %q, want %q", a.Case, CaseSVariableVCConstant)
	}
	if a.OldCompositionBP != a.NewCompositionBP {
		t.Errorf("composition changed %d → %d, want constant", a.OldCompositionBP, a.NewCompositionBP)
	}
	if a.OldProfitRate != 2000 || a.NewProfitRate != 1000 {
		t.Errorf("rates = %d → %d, want 2000 → 1000", a.OldProfitRate, a.NewProfitRate)
	}
	if !a.ProportionalChange {
		t.Error("ProportionalChange = false, want true when v/C is held constant")
	}
}

func TestDeriveVariationCase_BothVariable(t *testing.T) {
	t.Parallel()
	got := DeriveVariationCase(
		ProfitRateFormula{C: 80, V: 20, SRate: 10000}, // v/C = 2000 bp
		ProfitRateFormula{C: 120, V: 40, SRate: 5000}, // v/C = 2500 bp, s' changed
	)
	if got != CaseBothVariable {
		t.Errorf("case = %q, want %q", got, CaseBothVariable)
	}
}

func TestNewVariationAnalysisID_Unique(t *testing.T) {
	t.Parallel()
	a, b := NewVariationAnalysisID(), NewVariationAnalysisID()
	if a == b {
		t.Error("two IDs collided")
	}
	if a.IsZero() || len(string(a)) != 24 {
		t.Errorf("id = %q, want 24 hex chars", a)
	}
}
