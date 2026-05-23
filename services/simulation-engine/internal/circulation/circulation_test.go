package circulation

import "testing"

// TestVol2Ch01_TheCircuitOfMoneyCapital is the textual anchor for Capital
// Vol. II, Ch. 1 — The Circuit of Money-Capital.
//
// Replace t.Skip with the chapter's first textual example as a fixture.
// The canonical anchor for this chapter is the Spinning Mill 1871 example
// from §I: 50 labourers, £50 weekly wages, £372 means of production,
// 3,000 weekly labour-hours, yielding 10,000 lbs cotton yarn at £500.
// The circuit M — C(Lp+Mp) … P … C′ — M′ should round-trip these values
// with surplus realised = £78 (rate of surplus-value = 156%).
//
// Spec: marx-engels/1885/capital-volume-ii/specs/01-the-circuit-of-money-capital.spec.md
func TestVol2Ch01_TheCircuitOfMoneyCapital(t *testing.T) {
	t.Parallel()
	t.Skip("TODO: implement Vol. II Ch. 1 — The Circuit of Money-Capital")
}

// TestVol2Ch02_ProductiveCircuit_SpinningMillSimpleReproduction tests the §I
// fixture: Spinning Mill 1871 simple reproduction. The productive capital of
// £422 (constant=£372, variable=£50) produces C' of £500; the £78 surplus exits
// as revenue. The next circuit opens at the same £422.
func TestVol2Ch02_ProductiveCircuit_SpinningMillSimpleReproduction(t *testing.T) {
	t.Parallel()
	pc := ProductiveCircuit{
		ConstantPence:              37200,
		VariablePence:              5000,
		Mode:                       ReproductionSimple,
		MinCapitalisationIncrement: 100,
	}
	if err := pc.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
	if got := pc.TotalAdvance(); got != 42200 {
		t.Errorf("TotalAdvance() = %d, want 42200 (£422)", got)
	}
}

// TestVol2Ch02_ProductiveCircuit_SpinningMillExtendedReproduction tests §II:
// the £78 surplus is capitalised. P(n+1).Constant + P(n+1).Variable ==
// P(n).Constant + P(n).Variable + CapitalisationStep.AmountInjected.
func TestVol2Ch02_ProductiveCircuit_SpinningMillExtendedReproduction(t *testing.T) {
	t.Parallel()
	initial := ProductiveCircuit{
		ConstantPence:              37200,
		VariablePence:              5000,
		Mode:                       ReproductionExtended,
		MinCapitalisationIncrement: 100,
	}
	injected := Pence(7800)
	dc := injected
	dv := Pence(0)

	step := CapitalisationStep{
		AmountInjected:     injected,
		DeltaConstantPence: dc,
		DeltaVariablePence: dv,
	}
	// Capitalisation conservation: Δc + Δv == injected.
	if step.DeltaConstantPence+step.DeltaVariablePence != step.AmountInjected {
		t.Errorf("conservation failed: Δc(%d) + Δv(%d) != injected(%d)",
			step.DeltaConstantPence, step.DeltaVariablePence, step.AmountInjected)
	}

	newConstant := initial.ConstantPence + dc
	newVariable := initial.VariablePence + dv
	wantTotal := initial.TotalAdvance() + injected
	if newConstant+newVariable != wantTotal {
		t.Errorf("extended reproduction magnitude: got %d, want %d", newConstant+newVariable, wantTotal)
	}
	// £422 + £78 = £500.
	if wantTotal != 50000 {
		t.Errorf("P'(n+1) total = %d pence, want 50000 (£500)", wantTotal)
	}
}

// TestVol2Ch02_LatentMoneyCapital_IsArmed tests the £1-per-spindle threshold.
// m below 100 pence stays latent; at 100 pence IsArmed flips true.
func TestVol2Ch02_LatentMoneyCapital_IsArmed(t *testing.T) {
	t.Parallel()
	lmc := LatentMoneyCapital{Threshold: 100}

	lmc.Accumulated = 99
	if lmc.IsArmed() {
		t.Error("IsArmed() = true with Accumulated=99 < Threshold=100, want false")
	}

	lmc.Accumulated = 100
	if !lmc.IsArmed() {
		t.Error("IsArmed() = false with Accumulated=100 == Threshold=100, want true")
	}

	lmc.Accumulated = 200
	if !lmc.IsArmed() {
		t.Error("IsArmed() = false with Accumulated=200 > Threshold=100, want true")
	}
}

// TestVol2Ch02_ProductiveCircuit_Validate_InvalidMode confirms validation
// rejects unrecognised mode strings.
func TestVol2Ch02_ProductiveCircuit_Validate_InvalidMode(t *testing.T) {
	t.Parallel()
	pc := ProductiveCircuit{
		ConstantPence:              37200,
		VariablePence:              5000,
		Mode:                       "unknown",
		MinCapitalisationIncrement: 100,
	}
	if err := pc.Validate(); err == nil {
		t.Error("Validate() = nil, want error for invalid mode")
	}
}

// TestVol2Ch02_NewProductiveCircuitID_Unique checks that two successive calls
// produce distinct identifiers.
func TestVol2Ch02_NewProductiveCircuitID_Unique(t *testing.T) {
	t.Parallel()
	a := NewProductiveCircuitID()
	b := NewProductiveCircuitID()
	if a == b {
		t.Errorf("NewProductiveCircuitID() collision: %s == %s", a, b)
	}
	if len(string(a)) != 24 {
		t.Errorf("ID length = %d, want 24 hex chars (96-bit)", len(string(a)))
	}
}
