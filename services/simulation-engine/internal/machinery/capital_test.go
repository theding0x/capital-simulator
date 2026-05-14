package machinery_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
)

// Capital Vol. I, Ch. 15, §6 — Repulsion and Attraction of Workpeople.
// "variable capital of £3,000 ... 100 workmen at £30 a year ... machinery
//  costing £1,500 displaces 50 men; variable capital becomes £1,500."
//
// Conservation: total capital is unchanged in the substitution period; only
// composition shifts from variable to constant.
func TestCapitalComposition_Para6_MechanisationConservesTotal(t *testing.T) {
	t.Parallel()
	const annualWagePence = 30 * 240 // £30 = 7,200 pence
	before := machinery.CapitalComposition{
		Variable: machinery.VariableCapital(100 * annualWagePence),
		Constant: machinery.ConstantCapital(100 * annualWagePence), // £3,000
	}
	totalBefore := before.Total()
	const machineCostPence = 1_500 * 240 // £1,500
	const workersDisplaced = 50
	after := before.Substitute(machinery.VariableCapital(workersDisplaced*annualWagePence),
		machinery.ConstantCapital(machineCostPence))

	if after.Variable != machinery.VariableCapital(50*annualWagePence) {
		t.Fatalf("Variable after = %d, want %d", after.Variable, 50*annualWagePence)
	}
	if after.Total() != totalBefore {
		t.Fatalf("Total changed: before=%d, after=%d", totalBefore, after.Total())
	}
	if !(after.Constant > before.Constant) {
		t.Fatalf("Constant should rise: before=%d, after=%d", before.Constant, after.Constant)
	}
}

// Invariant: VariableCapital + ConstantCapital == TotalCapital (§6).
func TestCapitalComposition_TotalIsSum(t *testing.T) {
	t.Parallel()
	c := machinery.CapitalComposition{
		Variable: machinery.VariableCapital(1500),
		Constant: machinery.ConstantCapital(4500),
	}
	if c.Total() != 6000 {
		t.Fatalf("Total = %d, want 6000", c.Total())
	}
}

// §3C — Intensification of labour: when the working day is shortened, the
// remaining hours are filled more densely. The IntensityFactor scales
// effective labour per hour above 1.0.
func TestIntensityFactor_BoundedAboveOne(t *testing.T) {
	t.Parallel()
	if got := machinery.NewIntensityFactor(1.4); float64(got) != 1.4 {
		t.Fatalf("NewIntensityFactor(1.4) = %v, want 1.4", got)
	}
	if got := machinery.NewIntensityFactor(0.5); float64(got) != 1.0 {
		t.Fatalf("NewIntensityFactor(0.5) should clamp to 1.0; got %v", got)
	}
}

// EffectiveLabour applies the intensity multiplier to nominal labour time.
func TestEffectiveLabour_IntensifiedDay(t *testing.T) {
	t.Parallel()
	got := machinery.EffectiveLabour(machinery.LabourMinutes(600), machinery.NewIntensityFactor(1.5))
	if got != machinery.LabourMinutes(900) {
		t.Fatalf("EffectiveLabour(600, 1.5) = %d, want 900", got)
	}
}
