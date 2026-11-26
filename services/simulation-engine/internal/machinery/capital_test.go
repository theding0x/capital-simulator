package machinery_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
)

// Capital Vol. I, Ch. 15, §6 — Repulsion and Attraction of Workpeople.
// Marx's carpet factory: 100 workmen at £30 a year (variable £3,000) and £3,000
// of raw material a year; he discharges 50 and employs the other 50 with
// machinery costing £1,500. "Before the change, the total sum of £6,000
// consisted half of constant, and half of variable capital. After the change it
// consists of £4,500 constant (£3,000 raw material and £1,500 machinery), and
// £1,500 variable capital. The variable capital, instead of being one half, is
// only one quarter, of the total capital."
//
// Conservation: total capital is unchanged in the substitution period; only
// composition shifts from variable to constant.
func TestCapitalComposition_Para6_MechanisationConservesTotal(t *testing.T) {
	t.Parallel()
	const pencePerPound = 240
	const annualWagePence = 30 * pencePerPound // £30 = 7,200 pence
	before := machinery.CapitalComposition{
		Variable: machinery.VariableCapital(100 * annualWagePence),
		Constant: machinery.ConstantCapital(3_000 * pencePerPound), // £3,000 raw material
	}
	totalBefore := before.Total()
	if totalBefore != 6_000*pencePerPound {
		t.Fatalf("Total before = %d, want £6,000 = %d", totalBefore, 6_000*pencePerPound)
	}
	if 2*int64(before.Variable) != totalBefore {
		t.Fatalf("Variable before = %d, want one half of %d", before.Variable, totalBefore)
	}
	if got := before.Ratio(); got != 1.0 {
		t.Fatalf("c/v before = %v, want 1 (£3,000 : £3,000)", got)
	}

	const machineCostPence = 1_500 * pencePerPound // £1,500
	const workersDisplaced = 50
	after := before.Substitute(machinery.VariableCapital(workersDisplaced*annualWagePence),
		machinery.ConstantCapital(machineCostPence))

	if after.Variable != machinery.VariableCapital(1_500*pencePerPound) {
		t.Fatalf("Variable after = %d, want £1,500 = %d", after.Variable, 1_500*pencePerPound)
	}
	if after.Constant != machinery.ConstantCapital(4_500*pencePerPound) {
		t.Fatalf("Constant after = %d, want £4,500 (£3,000 raw material + £1,500 machinery) = %d",
			after.Constant, 4_500*pencePerPound)
	}
	if after.Total() != totalBefore {
		t.Fatalf("Total changed: before=%d, after=%d", totalBefore, after.Total())
	}
	if 4*int64(after.Variable) != after.Total() {
		t.Fatalf("Variable after = %d, want one quarter of %d", after.Variable, after.Total())
	}
	if got := after.Ratio(); got != 3.0 {
		t.Fatalf("c/v after = %v, want 3 (£4,500 : £1,500)", got)
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
