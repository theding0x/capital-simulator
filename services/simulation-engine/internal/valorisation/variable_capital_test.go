package valorisation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/valorisation"
)

// --- VariableCapitalAdvance.Validate -----------------------------------------

func TestVariableCapitalAdvance_Validate_Happy(t *testing.T) {
	t.Parallel()
	// Capitalist A: £500 advance, 5-week turnover (Marx Ch. 16 §2).
	adv := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		IndustrialCapitalID:         "test-ic-a",
		Pence:                       120000, // £500
		TurnoverPeriodMinutes:       50400,  // 5 weeks × 7 days × 1440 min
		TurnoversPerYearBasisPoints: 100000, // n=10
	}
	if err := adv.Validate(); err != nil {
		t.Fatalf("expected valid advance, got: %v", err)
	}
}

func TestVariableCapitalAdvance_Validate_ZeroPence(t *testing.T) {
	t.Parallel()
	adv := valorisation.VariableCapitalAdvance{
		Pence:                 0,
		TurnoverPeriodMinutes: 50400,
	}
	if err := adv.Validate(); err == nil {
		t.Fatal("expected error for zero pence")
	}
}

func TestVariableCapitalAdvance_Validate_ZeroTurnoverPeriod(t *testing.T) {
	t.Parallel()
	adv := valorisation.VariableCapitalAdvance{
		Pence:                 120000,
		TurnoverPeriodMinutes: 0,
	}
	if err := adv.Validate(); err == nil {
		t.Fatal("expected error for zero turnover period")
	}
}

// --- ComputePerCircuitRate ---------------------------------------------------

func TestComputePerCircuitRate_CapitalistA(t *testing.T) {
	t.Parallel()
	// Marx §2: s/v = 100% → BasisPoints = 10000.
	id := valorisation.NewVariableCapitalAdvanceID()
	rate, err := valorisation.ComputePerCircuitRate(id, 120000, 120000)
	if err != nil {
		t.Fatal(err)
	}
	if rate.BasisPoints != 10000 {
		t.Errorf("expected BasisPoints=10000, got %d", rate.BasisPoints)
	}
	if rate.VariableCapitalAdvanceID != id {
		t.Error("advance ID mismatch")
	}
}

func TestComputePerCircuitRate_SpinningMill1871(t *testing.T) {
	t.Parallel()
	// Spinning Mill 1871 calibrated example: s/v = 156% → BasisPoints = 15600.
	// Marx uses this to derive the annual rate of 8112%.
	id := valorisation.NewVariableCapitalAdvanceID()
	// £1 variable, £1.56 surplus.
	rate, err := valorisation.ComputePerCircuitRate(id, 15600, 10000)
	if err != nil {
		t.Fatal(err)
	}
	if rate.BasisPoints != 15600 {
		t.Errorf("expected BasisPoints=15600, got %d", rate.BasisPoints)
	}
}

func TestComputePerCircuitRate_NegativeSurplus(t *testing.T) {
	t.Parallel()
	id := valorisation.NewVariableCapitalAdvanceID()
	_, err := valorisation.ComputePerCircuitRate(id, -1, 10000)
	if err == nil {
		t.Fatal("expected error for negative surplus")
	}
}

func TestComputePerCircuitRate_ZeroVariable(t *testing.T) {
	t.Parallel()
	id := valorisation.NewVariableCapitalAdvanceID()
	_, err := valorisation.ComputePerCircuitRate(id, 5000, 0)
	if err == nil {
		t.Fatal("expected error for zero variable")
	}
}

func TestComputePerCircuitRate_ZeroSurplus(t *testing.T) {
	t.Parallel()
	// Zero surplus is valid (s/v = 0%).
	id := valorisation.NewVariableCapitalAdvanceID()
	rate, err := valorisation.ComputePerCircuitRate(id, 0, 10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rate.BasisPoints != 0 {
		t.Errorf("expected BasisPoints=0, got %d", rate.BasisPoints)
	}
}

// --- ComputeAnnualSurplusRate ------------------------------------------------

func TestComputeAnnualSurplusRate_CapitalistA(t *testing.T) {
	t.Parallel()
	// A: n=10, s/v=100% → annual rate = 1000% (100000 bp).
	id := valorisation.NewVariableCapitalAdvanceID()
	rate := valorisation.ComputeAnnualSurplusRate(id, 100000, 10000, "1871")
	if rate.AnnualBasisPoints != 100000 {
		t.Errorf("Capitalist A: expected AnnualBasisPoints=100000, got %d", rate.AnnualBasisPoints)
	}
	if rate.Period != "1871" {
		t.Errorf("expected period=1871, got %s", rate.Period)
	}
}

func TestComputeAnnualSurplusRate_CapitalistB(t *testing.T) {
	t.Parallel()
	// B: n=1, s/v=100% → annual rate = 100% (10000 bp).
	id := valorisation.NewVariableCapitalAdvanceID()
	rate := valorisation.ComputeAnnualSurplusRate(id, 10000, 10000, "1871")
	if rate.AnnualBasisPoints != 10000 {
		t.Errorf("Capitalist B: expected AnnualBasisPoints=10000, got %d", rate.AnnualBasisPoints)
	}
}

func TestComputeAnnualSurplusRate_SpinningMill1871(t *testing.T) {
	t.Parallel()
	// Spinning Mill 1871: n=52, s/v=156% → annual = 52 × 156% = 8112%.
	// n_bp = 520000, per_circuit_bp = 15600: annual = 520000 × 15600 / 10000 = 811200.
	id := valorisation.NewVariableCapitalAdvanceID()
	rate := valorisation.ComputeAnnualSurplusRate(id, 520000, 15600, "1871")
	if rate.AnnualBasisPoints != 811200 {
		t.Errorf("Spinning Mill 1871: expected AnnualBasisPoints=811200 (8112%%), got %d", rate.AnnualBasisPoints)
	}
}

func TestComputeAnnualSurplusRate_Quadrupling(t *testing.T) {
	t.Parallel()
	// Invariant: doubling both n and s/v quadruples the annual rate.
	id := valorisation.NewVariableCapitalAdvanceID()
	base := valorisation.ComputeAnnualSurplusRate(id, 10000, 10000, "test")  // n=1, s/v=100% → 10000
	double := valorisation.ComputeAnnualSurplusRate(id, 20000, 20000, "test") // n=2, s/v=200% → 40000
	if double.AnnualBasisPoints != 4*base.AnnualBasisPoints {
		t.Errorf("expected quadrupling: base=%d, doubled=%d", base.AnnualBasisPoints, double.AnnualBasisPoints)
	}
}

func TestComputeAnnualSurplusRate_ZeroN(t *testing.T) {
	t.Parallel()
	id := valorisation.NewVariableCapitalAdvanceID()
	rate := valorisation.ComputeAnnualSurplusRate(id, 0, 10000, "test")
	if rate.AnnualBasisPoints != 0 {
		t.Errorf("expected AnnualBasisPoints=0 when n=0, got %d", rate.AnnualBasisPoints)
	}
}

// --- ComputeAnnualSurplusRateContrast ----------------------------------------

func TestComputeAnnualSurplusRateContrast_CanonicalAB(t *testing.T) {
	t.Parallel()
	// A: 1000% (100000 bp), B: 100% (10000 bp) → ratio = 10× (100000 bp).
	aID := valorisation.NewVariableCapitalAdvanceID()
	bID := valorisation.NewVariableCapitalAdvanceID()
	contrast := valorisation.ComputeAnnualSurplusRateContrast(aID, bID, 100000, 10000)
	if contrast.RatioBasisPoints != 100000 {
		t.Errorf("expected RatioBasisPoints=100000, got %d", contrast.RatioBasisPoints)
	}
	if contrast.CapitalAID != aID {
		t.Error("CapitalAID mismatch")
	}
	if contrast.CapitalBID != bID {
		t.Error("CapitalBID mismatch")
	}
}

func TestComputeAnnualSurplusRateContrast_ZeroDenominator(t *testing.T) {
	t.Parallel()
	// Should not panic; returns zero ratio.
	aID := valorisation.NewVariableCapitalAdvanceID()
	bID := valorisation.NewVariableCapitalAdvanceID()
	contrast := valorisation.ComputeAnnualSurplusRateContrast(aID, bID, 100000, 0)
	if contrast.RatioBasisPoints != 0 {
		t.Errorf("expected RatioBasisPoints=0 for zero denominator, got %d", contrast.RatioBasisPoints)
	}
}

// --- ComputeVariableCapitalReproduction --------------------------------------

func TestComputeVariableCapitalReproduction_CapitalistA(t *testing.T) {
	t.Parallel()
	// A: £500 × n=10 = £5000 annual variable; s/v=100% → £5000 annual surplus.
	adv := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		IndustrialCapitalID:         "test-ic-a",
		Pence:                       120000, // £500
		TurnoverPeriodMinutes:       50400,
		TurnoversPerYearBasisPoints: 100000, // n=10
	}
	repro := valorisation.ComputeVariableCapitalReproduction(adv, 10000, 525600)
	// TotalVariable = 120000 × 100000 / 10000 = 1200000p (£5000).
	if repro.TotalVariablePaidAnnualPence != 1200000 {
		t.Errorf("expected TotalVariable=1200000, got %d", repro.TotalVariablePaidAnnualPence)
	}
	// TotalSurplus = 1200000 × 10000 / 10000 = 1200000p (£5000).
	if repro.TotalSurplusAnnualPence != 1200000 {
		t.Errorf("expected TotalSurplus=1200000, got %d", repro.TotalSurplusAnnualPence)
	}
}

func TestComputeVariableCapitalReproduction_CapitalistB(t *testing.T) {
	t.Parallel()
	// B: £5000 × n=1 = £5000 annual variable; s/v=100% → £5000 annual surplus.
	// Same total annual surplus as A — the key invariant.
	adv := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		IndustrialCapitalID:         "test-ic-b",
		Pence:                       1200000, // £5000
		TurnoverPeriodMinutes:       504000,
		TurnoversPerYearBasisPoints: 10000, // n=1
	}
	repro := valorisation.ComputeVariableCapitalReproduction(adv, 10000, 525600)
	if repro.TotalVariablePaidAnnualPence != 1200000 {
		t.Errorf("expected TotalVariable=1200000, got %d", repro.TotalVariablePaidAnnualPence)
	}
	if repro.TotalSurplusAnnualPence != 1200000 {
		t.Errorf("expected TotalSurplus=1200000, got %d", repro.TotalSurplusAnnualPence)
	}
}

func TestComputeVariableCapitalReproduction_SameAnnualSurplusRegardlessOfN(t *testing.T) {
	t.Parallel()
	// Invariant: same weekly labour → same total annual surplus regardless of n.
	// A and B both arrive at £5000 annual surplus.
	advA := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		Pence:                       120000,  // £500
		TurnoverPeriodMinutes:       50400,
		TurnoversPerYearBasisPoints: 100000, // n=10
	}
	advB := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		Pence:                       1200000, // £5000
		TurnoverPeriodMinutes:       504000,
		TurnoversPerYearBasisPoints: 10000, // n=1
	}
	reproA := valorisation.ComputeVariableCapitalReproduction(advA, 10000, 525600)
	reproB := valorisation.ComputeVariableCapitalReproduction(advB, 10000, 525600)
	if reproA.TotalSurplusAnnualPence != reproB.TotalSurplusAnnualPence {
		t.Errorf(
			"same weekly labour must produce same annual surplus: A=%d B=%d",
			reproA.TotalSurplusAnnualPence, reproB.TotalSurplusAnnualPence,
		)
	}
}

func TestComputeVariableCapitalReproduction_AtFullRate(t *testing.T) {
	t.Parallel()
	// At PerCircuitBasisPoints == 10000 (100%):
	// TotalSurplusAnnualPence == TotalVariablePaidAnnualPence.
	adv := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		Pence:                       240000,
		TurnoverPeriodMinutes:       10080,
		TurnoversPerYearBasisPoints: 520000, // n=52 (weekly)
	}
	repro := valorisation.ComputeVariableCapitalReproduction(adv, 10000, 525600)
	if repro.TotalSurplusAnnualPence != repro.TotalVariablePaidAnnualPence {
		t.Errorf(
			"at 100%% per-circuit rate, surplus should equal variable: surplus=%d variable=%d",
			repro.TotalSurplusAnnualPence, repro.TotalVariablePaidAnnualPence,
		)
	}
}
