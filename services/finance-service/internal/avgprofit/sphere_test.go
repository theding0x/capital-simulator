package avgprofit

import (
	"errors"
	"testing"
)

// --- ComputeProductionSphere fixture tests ---

// Marx's Sphere A: 600c + 100v, s' = 100%, total capital = 700.
// surplusValue = 10000*100/10000 = 100.
// individualProfitRate = round_half_up(100*10000/700) = round(1428.57) = 1429.
func TestComputeProductionSphere_SphereA(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("Sphere A (600c+100v)", 700, 100, 10000)
	if sp.SurplusValue != 100 {
		t.Errorf("SurplusValue = %d, want 100", sp.SurplusValue)
	}
	// Critical round-half-up assertion: truncation would give 1428.
	if sp.IndividualProfitRate != 1429 {
		t.Errorf("IndividualProfitRate = %d, want 1429 (round-half-up, not 1428)", sp.IndividualProfitRate)
	}
	if sp.ConstantCapital != 600 {
		t.Errorf("ConstantCapital = %d, want 600", sp.ConstantCapital)
	}
}

// Marx's Sphere B: 100c + 600v, s' = 100%, total capital = 700.
// surplusValue = 10000*600/10000 = 600.
// individualProfitRate = round_half_up(600*10000/700) = round(8571.43) = 8571.
func TestComputeProductionSphere_SphereB(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("Sphere B (100c+600v)", 700, 600, 10000)
	if sp.SurplusValue != 600 {
		t.Errorf("SurplusValue = %d, want 600", sp.SurplusValue)
	}
	if sp.IndividualProfitRate != 8571 {
		t.Errorf("IndividualProfitRate = %d, want 8571", sp.IndividualProfitRate)
	}
	if sp.ConstantCapital != 100 {
		t.Errorf("ConstantCapital = %d, want 100", sp.ConstantCapital)
	}
}

// European capital: 84c + 16v, s' = 100%, total capital = 100.
// surplusValue = 10000*16/10000 = 16.
// individualProfitRate = round_half_up(16*10000/100) = 1600.
func TestComputeProductionSphere_European(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("European capital (84c+16v)", 100, 16, 10000)
	if sp.SurplusValue != 16 {
		t.Errorf("SurplusValue = %d, want 16", sp.SurplusValue)
	}
	if sp.IndividualProfitRate != 1600 {
		t.Errorf("IndividualProfitRate = %d, want 1600", sp.IndividualProfitRate)
	}
}

// Asian capital: 16c + 84v, s' = 25% (2500 bp), total capital = 100.
// surplusValue = 2500*84/10000 = 21.
// individualProfitRate = round_half_up(21*10000/100) = 2100.
// Key point: lower s' than European (2500 vs 10000) but HIGHER p' (2100 vs 1600).
func TestComputeProductionSphere_Asian(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("Asian capital (16c+84v)", 100, 84, 2500)
	if sp.SurplusValue != 21 {
		t.Errorf("SurplusValue = %d, want 21", sp.SurplusValue)
	}
	if sp.IndividualProfitRate != 2100 {
		t.Errorf("IndividualProfitRate = %d, want 2100", sp.IndividualProfitRate)
	}
}

// Low organic composition: 90c + 10v, s' = 100%, total capital = 100.
// surplusValue = 10000*10/10000 = 10.
// individualProfitRate = round_half_up(10*10000/100) = 1000 (exactly 10%).
func TestComputeProductionSphere_LowComposition(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("Low composition (90c+10v)", 100, 10, 10000)
	if sp.SurplusValue != 10 {
		t.Errorf("SurplusValue = %d, want 10", sp.SurplusValue)
	}
	if sp.IndividualProfitRate != 1000 {
		t.Errorf("IndividualProfitRate = %d, want 1000", sp.IndividualProfitRate)
	}
}

// High organic composition: 10c + 90v, s' = 100%, total capital = 100.
// surplusValue = 10000*90/10000 = 90.
// individualProfitRate = round_half_up(90*10000/100) = 9000 (exactly 90%).
func TestComputeProductionSphere_HighComposition(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("High composition (10c+90v)", 100, 90, 10000)
	if sp.SurplusValue != 90 {
		t.Errorf("SurplusValue = %d, want 90", sp.SurplusValue)
	}
	if sp.IndividualProfitRate != 9000 {
		t.Errorf("IndividualProfitRate = %d, want 9000", sp.IndividualProfitRate)
	}
}

// --- LabourPowerIndex exact value tests ---

// Sphere A: v=100, C=700 → labourPowerIndex = 100*480*100/700 = 6857.
func TestLabourPowerIndex_SphereA(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("Sphere A", 700, 100, 10000)
	if sp.LabourPowerIndex != 6857 {
		t.Errorf("LabourPowerIndex = %d, want 6857", sp.LabourPowerIndex)
	}
}

// Sphere B: v=600, C=700 → labourPowerIndex = 600*480*100/700 = 41142.
// B carries 6× Sphere A's variable capital at the same total capital, and at the
// canonical 480 scale the truncated index lands exactly on 6× A's (6857·6 =
// 41142) — the integer-division divergence the old 3600 scale showed is gone.
func TestLabourPowerIndex_SphereB(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("Sphere B", 700, 600, 10000)
	if sp.LabourPowerIndex != 41142 {
		t.Errorf("LabourPowerIndex = %d, want 41142", sp.LabourPowerIndex)
	}
	// B's index is exactly 6× Sphere A's (same C, 6× v).
	spA := ComputeProductionSphere("Sphere A", 700, 100, 10000)
	if sp.LabourPowerIndex != LabourPowerIndex(6*int64(spA.LabourPowerIndex)) {
		t.Errorf("B.LabourPowerIndex = %d, want 6×A = %d", sp.LabourPowerIndex, 6*int64(spA.LabourPowerIndex))
	}
}

// --- Invariant tests ---

// Invariant: SurplusValue == SRate*V / 10000.
func TestInvariant_SurplusValue(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		c     TotalCapital
		v     VariableCapital
		sRate RateOfSurplusValue
	}{
		{"Sphere A", 700, 100, 10000},
		{"Sphere B", 700, 600, 10000},
		{"European", 100, 16, 10000},
		{"Asian", 100, 84, 2500},
		{"Low composition", 100, 10, 10000},
		{"High composition", 100, 90, 10000},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sp := ComputeProductionSphere(tc.name, tc.c, tc.v, tc.sRate)
			want := SurplusValue(int64(tc.sRate) * int64(tc.v) / basisPoints)
			if sp.SurplusValue != want {
				t.Errorf("SurplusValue = %d, want %d", sp.SurplusValue, want)
			}
		})
	}
}

// Invariant: IndividualProfitRate == round_half_up(SurplusValue*10000 / C).
func TestInvariant_IndividualProfitRate_RoundHalfUp(t *testing.T) {
	t.Parallel()
	// The round-half-up formula: (s*10000 + C/2) / C
	roundHalfUp := func(s SurplusValue, c TotalCapital) SphereProfitRate {
		if c <= 0 {
			return 0
		}
		return SphereProfitRate((int64(s)*basisPoints + int64(c)/2) / int64(c))
	}
	cases := []struct {
		name  string
		c     TotalCapital
		v     VariableCapital
		sRate RateOfSurplusValue
	}{
		{"Sphere A", 700, 100, 10000},
		{"Sphere B", 700, 600, 10000},
		{"European", 100, 16, 10000},
		{"Asian", 100, 84, 2500},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sp := ComputeProductionSphere(tc.name, tc.c, tc.v, tc.sRate)
			want := roundHalfUp(sp.SurplusValue, tc.c)
			if sp.IndividualProfitRate != want {
				t.Errorf("IndividualProfitRate = %d, want %d", sp.IndividualProfitRate, want)
			}
		})
	}
}

// Invariant: higher v/C implies higher IndividualProfitRate.
// Sphere B (v=600, C=700) > Sphere A (v=100, C=700) — same s' and same C.
func TestInvariant_Monotonicity_SphereB_GT_SphereA(t *testing.T) {
	t.Parallel()
	spA := ComputeProductionSphere("Sphere A", 700, 100, 10000)
	spB := ComputeProductionSphere("Sphere B", 700, 600, 10000)
	if spB.IndividualProfitRate <= spA.IndividualProfitRate {
		t.Errorf("Sphere B p'=%d must be > Sphere A p'=%d (higher v/C gives higher p')",
			spB.IndividualProfitRate, spA.IndividualProfitRate)
	}
	if spB.SurplusValue <= spA.SurplusValue {
		t.Errorf("Sphere B s=%d must be > Sphere A s=%d", spB.SurplusValue, spA.SurplusValue)
	}
}

// Invariant: Asian p' > European p' despite Asian s' < European s'.
// (Higher v/C wins over lower s'.)
func TestInvariant_Monotonicity_Asian_GT_European(t *testing.T) {
	t.Parallel()
	european := ComputeProductionSphere("European capital (84c+16v)", 100, 16, 10000)
	asian := ComputeProductionSphere("Asian capital (16c+84v)", 100, 84, 2500)
	// Asian s' (2500 bp = 25%) < European s' (10000 bp = 100%).
	if asian.SRate >= european.SRate {
		t.Errorf("Asian s'=%d must be < European s'=%d", asian.SRate, european.SRate)
	}
	// Yet Asian p' (2100 bp) > European p' (1600 bp).
	if asian.IndividualProfitRate <= european.IndividualProfitRate {
		t.Errorf("Asian p'=%d must be > European p'=%d despite lower s'",
			asian.IndividualProfitRate, european.IndividualProfitRate)
	}
}

// Invariant: V == 0 implies IndividualProfitRate == 0 (no living labour, no surplus).
func TestInvariant_ZeroVariable_ZeroProfitRate(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("No variable capital", 700, 0, 10000)
	if sp.IndividualProfitRate != 0 {
		t.Errorf("IndividualProfitRate = %d, want 0 when V=0", sp.IndividualProfitRate)
	}
	if sp.SurplusValue != 0 {
		t.Errorf("SurplusValue = %d, want 0 when V=0", sp.SurplusValue)
	}
}

// Boundary: V == C (all-variable capital; constant capital is zero).
func TestBoundary_VariableEqualsTotal(t *testing.T) {
	t.Parallel()
	// 0c + 100v, s' = 100% → s = 100, p' = 10000 bp (= 100%).
	sp := ComputeProductionSphere("All-variable capital", 100, 100, 10000)
	if sp.ConstantCapital != 0 {
		t.Errorf("ConstantCapital = %d, want 0", sp.ConstantCapital)
	}
	if sp.SurplusValue != 100 {
		t.Errorf("SurplusValue = %d, want 100", sp.SurplusValue)
	}
	if sp.IndividualProfitRate != 10000 {
		t.Errorf("IndividualProfitRate = %d, want 10000", sp.IndividualProfitRate)
	}
}

// Boundary: C == 0 (degenerate case; guard against divide-by-zero).
func TestBoundary_ZeroTotalCapital(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("Zero total capital", 0, 0, 10000)
	if sp.IndividualProfitRate != 0 {
		t.Errorf("IndividualProfitRate = %d, want 0 when C=0", sp.IndividualProfitRate)
	}
	if sp.LabourPowerIndex != 0 {
		t.Errorf("LabourPowerIndex = %d, want 0 when C=0", sp.LabourPowerIndex)
	}
}

// --- Validate tests ---

func TestValidate_Valid(t *testing.T) {
	t.Parallel()
	sp := ComputeProductionSphere("Sphere A", 700, 100, 10000)
	sp.C = 700
	sp.V = 100
	sp.SRate = 10000
	if err := sp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestValidate_NegativeC(t *testing.T) {
	t.Parallel()
	sp := ProductionSphere{C: -1, V: 0, SRate: 10000}
	if err := sp.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("Validate() = %v, want ErrNegativeMagnitude", err)
	}
}

func TestValidate_NegativeV(t *testing.T) {
	t.Parallel()
	sp := ProductionSphere{C: 100, V: -1, SRate: 10000}
	if err := sp.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("Validate() = %v, want ErrNegativeMagnitude", err)
	}
}

func TestValidate_NegativeSRate(t *testing.T) {
	t.Parallel()
	sp := ProductionSphere{C: 100, V: 50, SRate: -1}
	if err := sp.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("Validate() = %v, want ErrNegativeMagnitude", err)
	}
}

func TestValidate_VariableExceedsTotal(t *testing.T) {
	t.Parallel()
	sp := ProductionSphere{C: 100, V: 200, SRate: 10000}
	if err := sp.Validate(); !errors.Is(err, ErrVariableExceedsTotal) {
		t.Errorf("Validate() = %v, want ErrVariableExceedsTotal", err)
	}
}
