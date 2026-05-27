package circulation_test

// Vol. II Ch. 15 — Effects of a Change of Prices.
// All fixtures are drawn from Marx's Spinning Mill 1871 example and the
// 1857-depression-year compound-adjustment example from Ch. 15.

import (
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// --- RevaluedCapital --------------------------------------------------------

func TestDeriveRevaluedCapital_MPFall20Pct(t *testing.T) {
	t.Parallel()
	// Spinning Mill 1871: £370 in cotton (MP), 20% price fall → £74 freed.
	// BasisPointsChange = -2000 (−20%)
	// NewAdvancePence = 37000 * (10000 - 2000) / 10000 = 29600
	// DeltaPence = 29600 - 37000 = -7400 (£74 freed)
	icID := circulation.NewIndustrialCapitalID()
	evtID := circulation.NewValueRevolutionEventID()
	evt := circulation.ValueRevolutionEvent{
		ID:                  evtID,
		IndustrialCapitalID: icID,
		Affecting:           circulation.AffectingMeansOfProduction,
		DirectionPence:      -7400,
		BasisPointsChange:   -2000,
		OccurredAt:          time.Date(1871, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	rc := circulation.DeriveRevaluedCapital(evt, 37000)

	if rc.OldAdvancePence != 37000 {
		t.Errorf("OldAdvancePence = %d, want 37000", rc.OldAdvancePence)
	}
	if rc.NewAdvancePence != 29600 {
		t.Errorf("NewAdvancePence = %d, want 29600", rc.NewAdvancePence)
	}
	if rc.DeltaPence != -7400 {
		t.Errorf("DeltaPence = %d, want -7400 (£74 freed)", rc.DeltaPence)
	}
	if rc.ValueRevolutionEventID != evtID {
		t.Error("ValueRevolutionEventID not propagated")
	}
	if rc.IndustrialCapitalID != icID {
		t.Error("IndustrialCapitalID not propagated")
	}
	if err := rc.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestDeriveRevaluedCapital_MPRise20Pct(t *testing.T) {
	t.Parallel()
	// Spinning Mill 1871: cotton rises 20% → £74 tied up.
	// NewAdvancePence = 37000 * (10000 + 2000) / 10000 = 44400
	// DeltaPence = +7400
	evt := circulation.ValueRevolutionEvent{
		ID:                  circulation.NewValueRevolutionEventID(),
		IndustrialCapitalID: circulation.NewIndustrialCapitalID(),
		Affecting:           circulation.AffectingMeansOfProduction,
		DirectionPence:      7400,
		BasisPointsChange:   2000,
		OccurredAt:          time.Date(1871, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	rc := circulation.DeriveRevaluedCapital(evt, 37000)

	if rc.NewAdvancePence != 44400 {
		t.Errorf("NewAdvancePence = %d, want 44400", rc.NewAdvancePence)
	}
	if rc.DeltaPence != 7400 {
		t.Errorf("DeltaPence = %d, want +7400 (£74 tied up)", rc.DeltaPence)
	}
	if err := rc.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestRevaluedCapital_Validate_DeltaMismatch(t *testing.T) {
	t.Parallel()
	rc := circulation.RevaluedCapital{
		ID:                     circulation.NewRevaluedCapitalID(),
		ValueRevolutionEventID: circulation.NewValueRevolutionEventID(),
		IndustrialCapitalID:    circulation.NewIndustrialCapitalID(),
		OldAdvancePence:        37000,
		NewAdvancePence:        29600,
		DeltaPence:             9999, // wrong: should be -7400
	}
	if err := rc.Validate(); err == nil {
		t.Error("Validate() should fail when DeltaPence != NewAdvancePence - OldAdvancePence")
	}
}

func TestRevaluedCapital_Validate_ZeroOldAdvance(t *testing.T) {
	t.Parallel()
	rc := circulation.RevaluedCapital{
		OldAdvancePence: 0,
		NewAdvancePence: 10000,
		DeltaPence:      10000,
	}
	if err := rc.Validate(); err == nil {
		t.Error("Validate() should fail for OldAdvancePence == 0")
	}
}

// --- RealisedValueDelta -----------------------------------------------------

func TestRealisedValueDelta_SurplusAgainstOriginalAdvance_Loss(t *testing.T) {
	t.Parallel()
	// Catastrophic fall in output prices: expected £500, received £300.
	// Old advance was £400. Surplus against original advance = 300 - 400 = -100 (a loss).
	rvd := circulation.RealisedValueDelta{
		ID:                       circulation.NewRealisedValueDeltaID(),
		ValueRevolutionEventID:   circulation.NewValueRevolutionEventID(),
		IndustrialCapitalID:      circulation.NewIndustrialCapitalID(),
		ExpectedRealisationPence: 50000,
		ActualRealisationPence:   30000,
		DeltaPence:               -20000,
	}
	if err := rvd.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
	surplus := rvd.SurplusAgainstOriginalAdvance(40000)
	if surplus != -10000 {
		t.Errorf("SurplusAgainstOriginalAdvance(40000) = %d, want -10000", surplus)
	}
}

func TestRealisedValueDelta_SurplusAgainstOriginalAdvance_Positive(t *testing.T) {
	t.Parallel()
	// Output prices rose after production. Expected £500, got £600.
	// Surplus against £400 original advance = 600 - 400 = 200.
	rvd := circulation.RealisedValueDelta{
		ExpectedRealisationPence: 50000,
		ActualRealisationPence:   60000,
		DeltaPence:               10000,
	}
	if err := rvd.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
	surplus := rvd.SurplusAgainstOriginalAdvance(40000)
	if surplus != 20000 {
		t.Errorf("SurplusAgainstOriginalAdvance(40000) = %d, want 20000", surplus)
	}
}

func TestRealisedValueDelta_Validate_DeltaMismatch(t *testing.T) {
	t.Parallel()
	rvd := circulation.RealisedValueDelta{
		ExpectedRealisationPence: 50000,
		ActualRealisationPence:   30000,
		DeltaPence:               9999, // wrong: should be -20000
	}
	if err := rvd.Validate(); err == nil {
		t.Error("Validate() should fail when DeltaPence != actual - expected")
	}
}

func TestRealisedValueDelta_Validate_ZeroExpected(t *testing.T) {
	t.Parallel()
	rvd := circulation.RealisedValueDelta{
		ExpectedRealisationPence: 0,
		ActualRealisationPence:   30000,
		DeltaPence:               30000,
	}
	if err := rvd.Validate(); err == nil {
		t.Error("Validate() should fail for ExpectedRealisationPence == 0")
	}
}

// --- InventoryRevaluation ---------------------------------------------------

func TestInventoryRevaluation_Validate_CottonStockFall(t *testing.T) {
	t.Parallel()
	// Spinning Mill 1871: 1000 lbs cotton stock at 10p/lb → falls to 8p/lb.
	// DeltaPence = 1000 * (8 - 10) = -2000 (paper loss of £20).
	inv := circulation.InventoryRevaluation{
		ID:                  circulation.NewInventoryRevaluationID(),
		IndustrialCapitalID: circulation.NewIndustrialCapitalID(),
		CommodityID:         "cotton-raw",
		QuantityHeld:        1000,
		OldUnitPence:        10,
		NewUnitPence:        8,
		DeltaPence:          -2000,
	}
	if err := inv.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestInventoryRevaluation_Validate_YarnStockRise(t *testing.T) {
	t.Parallel()
	// Finished yarn: 500 lbs at 20p/lb → rises to 25p/lb.
	// DeltaPence = 500 * (25 - 20) = +2500.
	inv := circulation.InventoryRevaluation{
		ID:                  circulation.NewInventoryRevaluationID(),
		IndustrialCapitalID: circulation.NewIndustrialCapitalID(),
		CommodityID:         "yarn-finished",
		QuantityHeld:        500,
		OldUnitPence:        20,
		NewUnitPence:        25,
		DeltaPence:          2500,
	}
	if err := inv.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestInventoryRevaluation_Validate_DeltaMismatch(t *testing.T) {
	t.Parallel()
	inv := circulation.InventoryRevaluation{
		QuantityHeld: 1000,
		OldUnitPence: 10,
		NewUnitPence: 8,
		DeltaPence:   9999, // wrong
	}
	if err := inv.Validate(); err == nil {
		t.Error("Validate() should fail when DeltaPence != quantity * (new - old)")
	}
}

func TestInventoryRevaluation_Validate_ZeroQuantity(t *testing.T) {
	t.Parallel()
	inv := circulation.InventoryRevaluation{
		QuantityHeld: 0,
		OldUnitPence: 10,
		NewUnitPence: 8,
		DeltaPence:   0,
	}
	if err := inv.Validate(); err == nil {
		t.Error("Validate() should fail for QuantityHeld == 0")
	}
}

// --- SpeculativeHold --------------------------------------------------------

func TestSpeculativeHold_Validate_AnticipateRise(t *testing.T) {
	t.Parallel()
	// A spinner hoards cotton expecting prices to rise.
	sh := circulation.SpeculativeHold{
		ID:                   circulation.NewSpeculativeHoldID(),
		IndustrialCapitalID:  circulation.NewIndustrialCapitalID(),
		CommodityID:          "cotton-raw",
		HeldSince:            time.Date(1857, 1, 1, 0, 0, 0, 0, time.UTC),
		AnticipatedDirection: 1,
	}
	if err := sh.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestSpeculativeHold_Validate_Anticipatefall(t *testing.T) {
	t.Parallel()
	sh := circulation.SpeculativeHold{
		ID:                   circulation.NewSpeculativeHoldID(),
		IndustrialCapitalID:  circulation.NewIndustrialCapitalID(),
		CommodityID:          "cotton-raw",
		HeldSince:            time.Date(1857, 1, 1, 0, 0, 0, 0, time.UTC),
		AnticipatedDirection: -1,
	}
	if err := sh.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for direction -1", err)
	}
}

func TestSpeculativeHold_Validate_InvalidDirection(t *testing.T) {
	t.Parallel()
	sh := circulation.SpeculativeHold{
		ID:                   circulation.NewSpeculativeHoldID(),
		IndustrialCapitalID:  circulation.NewIndustrialCapitalID(),
		CommodityID:          "cotton-raw",
		HeldSince:            time.Now(),
		AnticipatedDirection: 0, // invalid
	}
	if err := sh.Validate(); err == nil {
		t.Error("Validate() should fail for AnticipatedDirection == 0")
	}
}

func TestSpeculativeHold_Validate_MissingCommodityID(t *testing.T) {
	t.Parallel()
	sh := circulation.SpeculativeHold{
		AnticipatedDirection: 1,
	}
	if err := sh.Validate(); err == nil {
		t.Error("Validate() should fail for empty CommodityID")
	}
}

// --- CompoundCapitalAdjustment ----------------------------------------------

func TestCompoundCapitalAdjustment_Validate_1857Depression(t *testing.T) {
	t.Parallel()
	// 1857 depression year: three simultaneous shocks to the Spinning Mill.
	// revaluation delta: -7400 (cotton fell, capital freed)
	// realisation delta: -5000 (output prices also fell)
	// inventory delta: -2000 (cotton stock paper loss)
	// net: -14400
	cca := circulation.CompoundCapitalAdjustment{
		ID:                    circulation.NewCompoundCapitalAdjustmentID(),
		IndustrialCapitalID:   circulation.NewIndustrialCapitalID(),
		Period:                "1857",
		RevaluationDeltaPence: -7400,
		RealisationDeltaPence: -5000,
		InventoryDeltaPence:   -2000,
		NetEffectPence:        -14400,
	}
	if err := cca.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
	if cca.NetEffect() != -14400 {
		t.Errorf("NetEffect() = %d, want -14400", cca.NetEffect())
	}
}

func TestCompoundCapitalAdjustment_Validate_NetMismatch(t *testing.T) {
	t.Parallel()
	cca := circulation.CompoundCapitalAdjustment{
		RevaluationDeltaPence: -7400,
		RealisationDeltaPence: -5000,
		InventoryDeltaPence:   -2000,
		NetEffectPence:        9999, // wrong: should be -14400
	}
	if err := cca.Validate(); err == nil {
		t.Error("Validate() should fail when NetEffectPence != sum of deltas")
	}
}

func TestCompoundCapitalAdjustment_NetEffect_Zero(t *testing.T) {
	t.Parallel()
	cca := circulation.CompoundCapitalAdjustment{
		RevaluationDeltaPence: 0,
		RealisationDeltaPence: 0,
		InventoryDeltaPence:   0,
		NetEffectPence:        0,
	}
	if err := cca.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for all-zero", err)
	}
}

// --- PriceCaseRecord --------------------------------------------------------

func TestPriceCaseRecord_Validate_FallCase(t *testing.T) {
	t.Parallel()
	fc := circulation.CaseSameScale
	p := circulation.PriceCaseRecord{
		ValueRevolutionEventID: circulation.NewValueRevolutionEventID(),
		FallCase:               &fc,
	}
	if err := p.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestPriceCaseRecord_Validate_RiseCase(t *testing.T) {
	t.Parallel()
	rc := circulation.CaseContractedScale
	p := circulation.PriceCaseRecord{
		ValueRevolutionEventID: circulation.NewValueRevolutionEventID(),
		RiseCase:               &rc,
	}
	if err := p.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestPriceCaseRecord_Validate_BothSet(t *testing.T) {
	t.Parallel()
	fc := circulation.CaseSameScale
	rc := circulation.CaseContractedScale
	p := circulation.PriceCaseRecord{FallCase: &fc, RiseCase: &rc}
	if err := p.Validate(); err == nil {
		t.Error("Validate() should fail when both FallCase and RiseCase are set")
	}
}

func TestPriceCaseRecord_Validate_NeitherSet(t *testing.T) {
	t.Parallel()
	p := circulation.PriceCaseRecord{}
	if err := p.Validate(); err == nil {
		t.Error("Validate() should fail when neither FallCase nor RiseCase is set")
	}
}
