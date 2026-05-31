// Package rent — Vol. III Ch. 42 unit tests: Differential Rent II, Falling
// Price of Production. Fixtures drawn from Marx Capital Vol. III Ch. 42
// Tables IV and IVd (Moore–Aveling edition).
package rent

import (
	"errors"
	"strings"
	"testing"
)

// ── FallingPriceVariant.IsValid ───────────────────────────────────────────────

func TestFallingPriceVariantIsValid(t *testing.T) {
	t.Parallel()

	valid := []FallingPriceVariant{
		FallingPriceConstantProductivity,
		FallingPriceDecreasingProductivity,
		FallingPriceRisingProductivity,
	}
	for _, v := range valid {
		if !v.IsValid() {
			t.Errorf("variant %d: IsValid() = false, want true", v)
		}
	}

	invalid := []FallingPriceVariant{0, 4}
	for _, v := range invalid {
		if v.IsValid() {
			t.Errorf("variant %d: IsValid() = true, want false", v)
		}
	}
}

// ── SoilExclusionEvent.Validate ───────────────────────────────────────────────

func TestSoilExclusionEventValidate(t *testing.T) {
	t.Parallel()

	// A(1) → B(2): new regulator > excluded → no error.
	evAB := SoilExclusionEvent{
		ExcludedGrade:      SoilGradeA,
		NewRegulatingGrade: SoilGradeB,
	}
	if err := evAB.Validate(); err != nil {
		t.Errorf("A→B: Validate() = %v, want nil", err)
	}

	// B(2) → A(1): new regulator < excluded → ErrInvalidExclusion.
	evBA := SoilExclusionEvent{
		ExcludedGrade:      SoilGradeB,
		NewRegulatingGrade: SoilGradeA,
	}
	if err := evBA.Validate(); !errors.Is(err, ErrInvalidExclusion) {
		t.Errorf("B→A: Validate() = %v, want ErrInvalidExclusion", err)
	}

	// A(1) → A(1): equal grades → ErrInvalidExclusion.
	evAA := SoilExclusionEvent{
		ExcludedGrade:      SoilGradeA,
		NewRegulatingGrade: SoilGradeA,
	}
	if err := evAA.Validate(); !errors.Is(err, ErrInvalidExclusion) {
		t.Errorf("A→A: Validate() = %v, want ErrInvalidExclusion", err)
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

// tableIVExclusion returns the Ch. 42 Table IV soil-exclusion event:
// Grade A excluded, Grade B becomes regulator at 45s (540 pence) per quarter.
func tableIVExclusion() SoilExclusionEvent {
	return SoilExclusionEvent{
		ID:                     NewSoilExclusionEventID(),
		ExcludedGrade:          SoilGradeA,
		NewRegulatingGrade:     SoilGradeB,
		NewPriceOfProductionBP: 540,
	}
}

// tableIVInvestments returns six successive investments from Marx Ch. 42
// Table IV: capital 250p each, output 400 1/100-qtrs each, tranches 1–6.
// SurplusProfitBP values (used as grain-rent in 1/100-qtr units): 0,0,50,50,75,75.
func tableIVInvestments() []SuccessiveInvestment {
	surplus := []int64{0, 0, 50, 50, 75, 75}
	invs := make([]SuccessiveInvestment, 6)
	for i, sp := range surplus {
		invs[i] = SuccessiveInvestment{
			ID:              NewSuccessiveInvestmentID(),
			TrancheNumber:   i + 1,
			CapitalBP:       250,
			OutputQuarters:  400,
			SurplusProfitBP: sp,
		}
	}
	return invs
}

// ── ComputeFallingPriceDR2 — Table IV ─────────────────────────────────────────

func TestComputeFallingPriceDR2_TableIV(t *testing.T) {
	t.Parallel()

	excl := tableIVExclusion()
	invs := tableIVInvestments()

	out := ComputeFallingPriceDR2(FallingPriceConstantProductivity, excl, invs)

	// TotalGrainRentQtrs = 0+0+50+50+75+75 = 250
	if out.TotalGrainRentQtrs != 250 {
		t.Errorf("TotalGrainRentQtrs = %d, want 250", out.TotalGrainRentQtrs)
	}
	// TotalMoneyRentalBP = roundHalfUp(250×540, 100) = 1350
	if out.TotalMoneyRentalBP != 1350 {
		t.Errorf("TotalMoneyRentalBP = %d, want 1350", out.TotalMoneyRentalBP)
	}
	// TotalCapitalBP = 6 × 250 = 1500
	if out.TotalCapitalBP != 1500 {
		t.Errorf("TotalCapitalBP = %d, want 1500", out.TotalCapitalBP)
	}
	if out.ExclusionEventID != string(excl.ID) {
		t.Errorf("ExclusionEventID = %q, want %q", out.ExclusionEventID, string(excl.ID))
	}
	if out.Variant != FallingPriceConstantProductivity {
		t.Errorf("Variant = %d, want FallingPriceConstantProductivity (%d)",
			out.Variant, FallingPriceConstantProductivity)
	}
}

// ── ComputeFallingPriceDR2 — Table IVd compensation ──────────────────────────

func TestComputeFallingPriceDR2_TableIVd_Compensation(t *testing.T) {
	t.Parallel()

	// Table IVd: price halved to 30s (360p); doubled grain-rent compensates.
	excl := SoilExclusionEvent{
		ID:                     NewSoilExclusionEventID(),
		ExcludedGrade:          SoilGradeA,
		NewRegulatingGrade:     SoilGradeB,
		NewPriceOfProductionBP: 360,
	}

	// Five investments whose SurplusProfitBP (grain-rent units) sum to 500.
	surplus := []int64{100, 100, 100, 100, 100}
	invs := make([]SuccessiveInvestment, len(surplus))
	for i, sp := range surplus {
		invs[i] = SuccessiveInvestment{
			ID:              NewSuccessiveInvestmentID(),
			TrancheNumber:   i + 1,
			CapitalBP:       250,
			OutputQuarters:  400,
			SurplusProfitBP: sp,
		}
	}

	out := ComputeFallingPriceDR2(FallingPriceDecreasingProductivity, excl, invs)

	// TotalMoneyRentalBP = roundHalfUp(500×360, 100) = 1800
	// Compensation: equals original Table I rental despite lower price.
	if out.TotalMoneyRentalBP != 1800 {
		t.Errorf("TotalMoneyRentalBP = %d, want 1800 (compensation law)", out.TotalMoneyRentalBP)
	}
}

// ── ComputeFallingPriceDR2 — empty investments ────────────────────────────────

func TestComputeFallingPriceDR2_Empty(t *testing.T) {
	t.Parallel()

	excl := tableIVExclusion()
	out := ComputeFallingPriceDR2(FallingPriceConstantProductivity, excl, nil)

	if out.TotalGrainRentQtrs != 0 {
		t.Errorf("TotalGrainRentQtrs = %d, want 0 on empty investments", out.TotalGrainRentQtrs)
	}
	if out.TotalMoneyRentalBP != 0 {
		t.Errorf("TotalMoneyRentalBP = %d, want 0 on empty investments", out.TotalMoneyRentalBP)
	}
	if out.TotalCapitalBP != 0 {
		t.Errorf("TotalCapitalBP = %d, want 0 on empty investments", out.TotalCapitalBP)
	}
	if out.TotalOutputQuarters != 0 {
		t.Errorf("TotalOutputQuarters = %d, want 0 on empty investments", out.TotalOutputQuarters)
	}
	if out.Variant != FallingPriceConstantProductivity {
		t.Errorf("Variant = %d, want FallingPriceConstantProductivity", out.Variant)
	}
	if out.ExclusionEventID != string(excl.ID) {
		t.Errorf("ExclusionEventID = %q, want %q", out.ExclusionEventID, string(excl.ID))
	}
}

// ── ComputeFallingPriceDR2 — zero price guard ─────────────────────────────────

func TestComputeFallingPriceDR2_ZeroPriceGuard(t *testing.T) {
	t.Parallel()

	excl := SoilExclusionEvent{
		ID:                     NewSoilExclusionEventID(),
		ExcludedGrade:          SoilGradeA,
		NewRegulatingGrade:     SoilGradeB,
		NewPriceOfProductionBP: 0,
	}
	invs := tableIVInvestments()
	out := ComputeFallingPriceDR2(FallingPriceConstantProductivity, excl, invs)

	if out.TotalMoneyRentalBP != 0 {
		t.Errorf("TotalMoneyRentalBP = %d, want 0 when NewPriceOfProductionBP=0", out.TotalMoneyRentalBP)
	}
}

// ── ID constructors ───────────────────────────────────────────────────────────

func TestCh42NewIDsDistinct(t *testing.T) {
	t.Parallel()

	// NewSoilExclusionEventID: 24-char hex, two calls differ.
	id1 := NewSoilExclusionEventID()
	id2 := NewSoilExclusionEventID()
	if len(string(id1)) != 24 {
		t.Errorf("SoilExclusionEventID len = %d, want 24", len(string(id1)))
	}
	if !isLowerHex(string(id1)) {
		t.Errorf("SoilExclusionEventID %q is not lowercase hex", id1)
	}
	if id1 == id2 {
		t.Error("two NewSoilExclusionEventID() calls returned the same value")
	}

	// NewFallingPriceDR2OutcomeID: 24-char hex, two calls differ.
	oid1 := NewFallingPriceDR2OutcomeID()
	oid2 := NewFallingPriceDR2OutcomeID()
	if len(string(oid1)) != 24 {
		t.Errorf("FallingPriceDR2OutcomeID len = %d, want 24", len(string(oid1)))
	}
	if !isLowerHex(string(oid1)) {
		t.Errorf("FallingPriceDR2OutcomeID %q is not lowercase hex", oid1)
	}
	if oid1 == oid2 {
		t.Error("two NewFallingPriceDR2OutcomeID() calls returned the same value")
	}

	// NewNormalCapitalPerAcreID: 24-char hex, two calls differ.
	nid1 := NewNormalCapitalPerAcreID()
	nid2 := NewNormalCapitalPerAcreID()
	if len(string(nid1)) != 24 {
		t.Errorf("NormalCapitalPerAcreID len = %d, want 24", len(string(nid1)))
	}
	if !isLowerHex(string(nid1)) {
		t.Errorf("NormalCapitalPerAcreID %q is not lowercase hex", nid1)
	}
	if nid1 == nid2 {
		t.Error("two NewNormalCapitalPerAcreID() calls returned the same value")
	}
}

// isLowerHex reports whether s consists entirely of [0-9a-f].
func isLowerHex(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !strings.ContainsRune("0123456789abcdef", r) {
			return false
		}
	}
	return true
}
