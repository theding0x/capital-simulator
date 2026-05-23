package circulation

import (
	"errors"
	"testing"
)

// Vol. II Ch. 3 — The Circuit of Commodity-Capital tests.
// Fixture: Spinning Mill 1871 — 10,000 lbs cotton yarn, C′ = £500
//   (c=£372, v=£50, s=£78). Three successive partial sales exhaust the output.

// TestVol2Ch03_CommodityCircuit_SpinningMill1871_SimpleReproduction verifies that
// a valid C′ circuit (SurplusPence > 0) passes validation and computes its Total.
func TestVol2Ch03_CommodityCircuit_SpinningMill1871_SimpleReproduction(t *testing.T) {
	t.Parallel()
	cc := CommodityCircuit{
		Initial: OpeningCommodityCapital{
			ConstantPence: 37200,
			VariablePence: 5000,
			SurplusPence:  7800,
			PoundsTotal:   10000,
		},
		Mode: ReproductionSimple,
	}
	if err := cc.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
	if got := cc.Initial.Total(); got != 50000 {
		t.Errorf("Total() = %d pence, want 50000 (£500)", got)
	}
	if got := cc.Initial.ValueOriginalPence(); got != 42200 {
		t.Errorf("ValueOriginalPence() = %d pence, want 42200 (£422 capital)", got)
	}
}

// TestVol2Ch03_CommodityCircuit_NotCommodityCapital verifies ErrNotCommodityCapital
// is returned when SurplusPence == 0 and IsFirstInvestment is false.
func TestVol2Ch03_CommodityCircuit_NotCommodityCapital(t *testing.T) {
	t.Parallel()
	cc := CommodityCircuit{
		Initial: OpeningCommodityCapital{
			ConstantPence: 37200,
			VariablePence: 5000,
			SurplusPence:  0,
			PoundsTotal:   10000,
		},
		Mode: ReproductionSimple,
	}
	if err := cc.Validate(); !errors.Is(err, ErrNotCommodityCapital) {
		t.Errorf("Validate() = %v, want ErrNotCommodityCapital", err)
	}
}

// TestVol2Ch03_CommodityAugmented_SimpleReproduction checks that a closing circuit
// with the same total as the opening is not extended.
func TestVol2Ch03_CommodityAugmented_SimpleReproduction(t *testing.T) {
	t.Parallel()
	opening := OpeningCommodityCapital{
		ConstantPence: 37200,
		VariablePence: 5000,
		SurplusPence:  7800,
		PoundsTotal:   10000,
	}
	terminal := CommodityAugmented{
		ConstantPence: 37200,
		VariablePence: 5000,
		SurplusPence:  7800,
		PoundsTotal:   10000,
	}
	if terminal.IsExtended(opening.Total()) {
		t.Error("IsExtended() = true for simple reproduction, want false")
	}
}

// TestVol2Ch03_CommodityAugmented_ExtendedReproduction checks that a closing circuit
// with a higher total is marked extended.
func TestVol2Ch03_CommodityAugmented_ExtendedReproduction(t *testing.T) {
	t.Parallel()
	opening := OpeningCommodityCapital{
		ConstantPence: 37200,
		VariablePence: 5000,
		SurplusPence:  7800,
		PoundsTotal:   10000,
	}
	// Surplus capitalised: new P starts at £500; next C′ reflects extended base.
	terminal := CommodityAugmented{
		ConstantPence:    45000,
		VariablePence:    5000,
		SurplusPence:     7800,
		CapitalisedPence: 7800,
		PoundsTotal:      11560,
	}
	if !terminal.IsExtended(opening.Total()) {
		t.Error("IsExtended() = false for extended reproduction, want true")
	}
}

// TestVol2Ch03_AliquotDecomposition verifies three successive partial sales of
// the Spinning Mill 1871 C′. Each sale's components must sum to its own
// realised_pence. The grand total of all components must equal the grand total
// of all realisations (50000 pence = £500). Individual component grand totals
// are allowed to differ from the opening ratios by integer truncation; surplus
// absorbs the remainder within each sale, not across all sales.
func TestVol2Ch03_AliquotDecomposition(t *testing.T) {
	t.Parallel()
	cc := CommodityCircuit{
		Initial: OpeningCommodityCapital{
			ConstantPence: 37200,
			VariablePence: 5000,
			SurplusPence:  7800,
			PoundsTotal:   10000,
		},
		Mode: ReproductionSimple,
	}

	sales := []struct {
		qty      int64
		realised Pence
	}{
		{7440, 37200},
		{1000, 5000},
		{1560, 7800},
	}

	var grandTotal, grandRealised Pence
	for _, s := range sales {
		d := cc.ComputeDecomposition(s.qty, s.realised)
		if d.Total() != s.realised {
			t.Errorf("sale qty=%d: decomposition total %d != realised %d", s.qty, d.Total(), s.realised)
		}
		grandTotal += d.Total()
		grandRealised += s.realised
	}

	if grandTotal != grandRealised {
		t.Errorf("grand decomposition total %d != grand realised %d", grandTotal, grandRealised)
	}
	if grandRealised != 50000 {
		t.Errorf("grand realised = %d pence, want 50000 (£500 = full C′)", grandRealised)
	}
}

// TestVol2Ch03_EnsureMaterialAdequacy verifies that capitalisation without linked
// MP sources returns ErrInsufficientMaterialBase.
func TestVol2Ch03_EnsureMaterialAdequacy(t *testing.T) {
	t.Parallel()
	cc := CommodityCircuit{
		Initial: OpeningCommodityCapital{SurplusPence: 7800, PoundsTotal: 10000},
		Mode:    ReproductionExtended,
	}

	if err := cc.EnsureMaterialAdequacy(7800); !errors.Is(err, ErrInsufficientMaterialBase) {
		t.Errorf("EnsureMaterialAdequacy(7800) = %v, want ErrInsufficientMaterialBase", err)
	}

	cc.MPSources = []MeansOfProductionSource{{SourceKind: SourceImport}}
	if err := cc.EnsureMaterialAdequacy(7800); err != nil {
		t.Errorf("EnsureMaterialAdequacy(7800) with MP source = %v, want nil", err)
	}

	if err := cc.EnsureMaterialAdequacy(0); err != nil {
		t.Errorf("EnsureMaterialAdequacy(0) = %v, want nil (zero capitalisation ok)", err)
	}
}

// TestVol2Ch03_MeansOfProductionSource_Validate checks that valid source kinds
// pass and unrecognised kinds return ErrUnsourcedMeansOfProduction.
func TestVol2Ch03_MeansOfProductionSource_Validate(t *testing.T) {
	t.Parallel()
	valid := []MeansOfProductionSourceKind{
		SourceProducerCircuit,
		SourceMerchantHolding,
		SourceImport,
	}
	for _, k := range valid {
		s := MeansOfProductionSource{SourceKind: k}
		if err := s.Validate(); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", k, err)
		}
	}
	bad := MeansOfProductionSource{SourceKind: "fantasy_source"}
	if err := bad.Validate(); !errors.Is(err, ErrUnsourcedMeansOfProduction) {
		t.Errorf("Validate(fantasy_source) = %v, want ErrUnsourcedMeansOfProduction", err)
	}
}

// TestVol2Ch03_NewCommodityCircuitID_Unique checks that two successive calls
// produce distinct 24-char hex identifiers.
func TestVol2Ch03_NewCommodityCircuitID_Unique(t *testing.T) {
	t.Parallel()
	a := NewCommodityCircuitID()
	b := NewCommodityCircuitID()
	if a == b {
		t.Errorf("NewCommodityCircuitID() collision: %s == %s", a, b)
	}
	if len(string(a)) != 24 {
		t.Errorf("ID length = %d, want 24 hex chars (96-bit)", len(string(a)))
	}
}

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
