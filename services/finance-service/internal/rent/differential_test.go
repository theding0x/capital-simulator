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

	t.Run("waterfall capital earns surplus-profit", func(t *testing.T) {
		t.Parallel()
		// General PoP 11500 bp; individual (waterfall-aided) 9000 bp.
		// Surplus-profit = 11500 − 9000 = 2500.
		got := ComputeSurplusProfit(11500, 9000)
		if got != 2500 {
			t.Errorf("ComputeSurplusProfit(11500,9000) = %d, want 2500", got)
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
