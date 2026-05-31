// Package rent — unit tests for Vol. III Ch. 46 types: Building Site Rent,
// Rent in Mining, Price of Land.
//
// Fixtures drawn from Marx, Capital Vol. III Ch. 46:
//   - Monopoly-price rent: Burgundy grand cru vineyard, sale price 400
//     labour-minutes, value 200 → monopoly_rent = 200.
//   - Capitalised land price: annual rent 10 bp, interest rate 250 bp (2.5%)
//     → land price = 400 bp; interest rate 500 bp (5%) → land price = 200 bp.
//   - Zero interest rate → land price = 0 (indeterminate).
package rent

import (
	"testing"
)

// TestLandPriceScenarioKindIsValid covers the full valid range (1–5) and the
// two invalid sentinels (0 and 6).
func TestLandPriceScenarioKindIsValid(t *testing.T) {
	t.Parallel()

	valid := []LandPriceScenarioKind{
		LandPriceFallingInterestRate,
		LandPriceCapitalImprovements,
		LandPriceRentRisesWithProduct,
		LandPriceRentRisesMoreCapital,
		LandPriceRentRisesPriceFalls,
	}
	for _, k := range valid {
		if !k.IsValid() {
			t.Errorf("kind %d: IsValid() = false, want true", k)
		}
	}

	invalid := []LandPriceScenarioKind{0, 6}
	for _, k := range invalid {
		if k.IsValid() {
			t.Errorf("kind %d: IsValid() = true, want false", k)
		}
	}
}

// TestComputeMonopolyRent verifies max(0, sale−value).
// Marx Ch. 46: burgundy grand cru, sale price 400, value 200 → rent 200.
func TestComputeMonopolyRent(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		salePrice int64
		value     int64
		want      int64
	}{
		{"burgundy_grand_cru sale>value", 400, 200, 200},
		{"sale==value no rent", 200, 200, 0},
		{"sale<value floors at zero", 100, 200, 0},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ComputeMonopolyRent(tc.salePrice, tc.value)
			if got != tc.want {
				t.Errorf("ComputeMonopolyRent(%d, %d) = %d, want %d",
					tc.salePrice, tc.value, got, tc.want)
			}
		})
	}
}

// TestCapitalisedLandPrice verifies annualRent × 10000 ÷ interestRateBP.
// Marx Ch. 46: rent 10 bp, rate 250 bp (2.5%) → price 400 bp;
//
//	rent 10 bp, rate 500 bp (5%)   → price 200 bp;
//	rate == 0                      → price   0 bp (undefined).
func TestCapitalisedLandPrice(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		annualRent int64
		rateBP     int64
		want       int64
	}{
		{"rent10_rate250bp_price400", 10, 250, 400},
		{"rent10_rate500bp_price200", 10, 500, 200},
		{"zero_rate_returns_zero", 10, 0, 0},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := CapitalisedLandPrice(tc.annualRent, tc.rateBP)
			if got != tc.want {
				t.Errorf("CapitalisedLandPrice(%d, %d) = %d, want %d",
					tc.annualRent, tc.rateBP, got, tc.want)
			}
		})
	}
}

// TestCh46NewIDsDistinct verifies that each New*ID constructor returns a
// 24-character hex string and that two successive calls differ.
func TestCh46NewIDsDistinct(t *testing.T) {
	t.Parallel()

	t.Run("BuildingSiteRentID", func(t *testing.T) {
		t.Parallel()
		a, b := NewBuildingSiteRentID(), NewBuildingSiteRentID()
		if len(string(a)) != 24 {
			t.Errorf("len(BuildingSiteRentID) = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two NewBuildingSiteRentID() calls must differ")
		}
	})

	t.Run("MiningRentID", func(t *testing.T) {
		t.Parallel()
		a, b := NewMiningRentID(), NewMiningRentID()
		if len(string(a)) != 24 {
			t.Errorf("len(MiningRentID) = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two NewMiningRentID() calls must differ")
		}
	})

	t.Run("MonopolyPriceRentID", func(t *testing.T) {
		t.Parallel()
		a, b := NewMonopolyPriceRentID(), NewMonopolyPriceRentID()
		if len(string(a)) != 24 {
			t.Errorf("len(MonopolyPriceRentID) = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two NewMonopolyPriceRentID() calls must differ")
		}
	})

	t.Run("LandPriceScenarioID", func(t *testing.T) {
		t.Parallel()
		a, b := NewLandPriceScenarioID(), NewLandPriceScenarioID()
		if len(string(a)) != 24 {
			t.Errorf("len(LandPriceScenarioID) = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two NewLandPriceScenarioID() calls must differ")
		}
	})
}
