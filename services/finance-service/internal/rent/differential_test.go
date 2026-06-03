// Package rent — Ch. 38 tests: Differential Rent, General Remarks.
// Fixtures drawn from Marx, Vol. III Ch. 38: the waterfall example and the
// capitalised-rent formula (annual rent / interest rate).
package rent

import (
	"testing"
)

// TestComputeSurplusProfit verifies the surplus-profit formula
// (general − individual price of production).
func TestComputeSurplusProfit(t *testing.T) {
	t.Parallel()

	t.Run("waterfall surplus-profit = general PoP − individual PRICE OF PRODUCTION", func(t *testing.T) {
		t.Parallel()
		// Waterfall factory: individual cost-price £90 (9000) + absolute average
		// profit £15 (1500) = individual price of production £105 (10500). General
		// (market-regulating) PoP £115 (11500). Surplus-profit = 11500 − 10500 =
		// 1000 (£10) — NOT 11500 − 9000 = 2500 (£25), which wrongly fed the bare
		// cost-price.
		individualPoP := IndividualPriceOfProduction(9000, 1500)
		if individualPoP != 10500 {
			t.Fatalf("IndividualPriceOfProduction(9000,1500) = %d, want 10500", individualPoP)
		}
		got := ComputeSurplusProfit(11500, individualPoP)
		if got != 1000 {
			t.Errorf("ComputeSurplusProfit(11500, individual PoP 10500) = %d, want 1000 (£10)", got)
		}
	})

	t.Run("equal prices yield zero surplus-profit", func(t *testing.T) {
		t.Parallel()
		got := ComputeSurplusProfit(9000, 9000)
		if got != 0 {
			t.Errorf("ComputeSurplusProfit(9000,9000) = %d, want 0", got)
		}
	})

	t.Run("individual above general gives negative result (not clamped)", func(t *testing.T) {
		t.Parallel()
		// Individual 11500 > general 9000 → −2500.
		got := ComputeSurplusProfit(9000, 11500)
		if got != -2500 {
			t.Errorf("ComputeSurplusProfit(9000,11500) = %d, want -2500", got)
		}
	})
}

// TestSurplusProfitEqualsDifferentialRent pins the Ch. 38 identity (#378): the
// waterfall surplus-profit (£10 = 1000 pence) equals the differential rent it
// becomes (£10 = 4800 LabourMinutes), and that rent capitalised at 5% is the
// £200 (96000 LM) land-price seeded alongside.
func TestSurplusProfitEqualsDifferentialRent(t *testing.T) {
	t.Parallel()
	surplusProfitBP := ComputeSurplusProfit(11500, IndividualPriceOfProduction(9000, 1500)) // 1000
	const rentLM = 4800                                                                      // £10
	if !SurplusProfitMatchesRent(surplusProfitBP, rentLM) {
		t.Errorf("surplus-profit %d pence does not match rent %d LM (should both be £10)", surplusProfitBP, rentLM)
	}
	// A DifferentialRent built from the matched magnitudes is internally consistent.
	dr := DifferentialRent{SurplusProfitBP: surplusProfitBP, AnnualRentLabourMinutes: rentLM}
	if err := dr.Validate(); err != nil {
		t.Errorf("DifferentialRent.Validate() = %v, want nil", err)
	}
	// The old contradictory pairing (£25 surplus next to £10 rent) is rejected.
	bad := DifferentialRent{SurplusProfitBP: 2500, AnnualRentLabourMinutes: 4800}
	if err := bad.Validate(); err == nil {
		t.Error("DifferentialRent{2500 pence, 4800 LM}.Validate() = nil, want mismatch error")
	}
	// And the rent capitalises to the seeded £200 land-price at 5%.
	if got := ComputeCapitalisedPrice(rentLM, 500); got != 96000 {
		t.Errorf("ComputeCapitalisedPrice(4800,500) = %d, want 96000 (£200)", got)
	}
}

// TestComputeCapitalisedPrice verifies the land-purchase-price formula:
// roundHalfUp(annualRent × 10000, interestRateBP).
func TestComputeCapitalisedPrice(t *testing.T) {
	t.Parallel()

	t.Run("Marx fixture: £10 rent at 5% interest", func(t *testing.T) {
		t.Parallel()
		// £10 annual rent = 4800 LM; 5% = 500 bp.
		// Capitalised price = roundHalfUp(4800×10000, 500) = 96 000 000/500 = 96000 LM.
		got := ComputeCapitalisedPrice(4800, 500)
		if got != 96000 {
			t.Errorf("ComputeCapitalisedPrice(4800,500) = %d, want 96000", got)
		}
	})

	t.Run("non-exact: rounds half-up", func(t *testing.T) {
		t.Parallel()
		// roundHalfUp(1000×10000, 300) = roundHalfUp(10_000_000, 300)
		// = (10_000_000 + 150) / 300 = 10_000_150/300 = 33333.
		got := ComputeCapitalisedPrice(1000, 300)
		if got != 33333 {
			t.Errorf("ComputeCapitalisedPrice(1000,300) = %d, want 33333", got)
		}
	})

	t.Run("zero interest rate returns 0", func(t *testing.T) {
		t.Parallel()
		got := ComputeCapitalisedPrice(4800, 0)
		if got != 0 {
			t.Errorf("ComputeCapitalisedPrice(4800,0) = %d, want 0", got)
		}
	})

	t.Run("negative interest rate returns 0", func(t *testing.T) {
		t.Parallel()
		got := ComputeCapitalisedPrice(4800, -5)
		if got != 0 {
			t.Errorf("ComputeCapitalisedPrice(4800,-5) = %d, want 0", got)
		}
	})
}

// TestCh38NewIDsDistinct verifies that each constructor returns a 24-char hex
// string and that two successive calls produce different values.
func TestCh38NewIDsDistinct(t *testing.T) {
	t.Parallel()

	t.Run("NewPoPSurplusProfitID", func(t *testing.T) {
		t.Parallel()
		a, b := NewPoPSurplusProfitID(), NewPoPSurplusProfitID()
		if len(string(a)) != 24 {
			t.Errorf("id length = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two calls returned identical IDs")
		}
	})

	t.Run("NewMonopolisedNaturalForceID", func(t *testing.T) {
		t.Parallel()
		a, b := NewMonopolisedNaturalForceID(), NewMonopolisedNaturalForceID()
		if len(string(a)) != 24 {
			t.Errorf("id length = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two calls returned identical IDs")
		}
	})

	t.Run("NewDifferentialRentID", func(t *testing.T) {
		t.Parallel()
		a, b := NewDifferentialRentID(), NewDifferentialRentID()
		if len(string(a)) != 24 {
			t.Errorf("id length = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two calls returned identical IDs")
		}
	})

	t.Run("NewCapitalisedRentPriceID", func(t *testing.T) {
		t.Parallel()
		a, b := NewCapitalisedRentPriceID(), NewCapitalisedRentPriceID()
		if len(string(a)) != 24 {
			t.Errorf("id length = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two calls returned identical IDs")
		}
	})
}
