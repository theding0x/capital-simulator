// Package rent — unit tests for Vol. III Ch. 45 (Absolute Ground-Rent) pure
// functions and ID constructors.
//
// Fixture: agricultural organic composition (C:V = 60:40 = 15000 bp) is below
// the social average (40000 bp), so products sell near value, creating a
// value–price gap that landed property captures as absolute rent.
// Source: Marx, Capital Vol. III Ch. 45.
package rent

import (
	"testing"
)

// TestComputeAbsoluteRentBP verifies max(0, S−P) using Marx Ch. 45 numerics.
func TestComputeAbsoluteRentBP(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		surplusValueBP  int64
		averageProfitBP int64
		want            int64
	}{
		{"surplus exceeds average: rent = difference", 40, 20, 20},
		{"surplus equals average: rent = zero", 20, 20, 0},
		{"surplus below average: clamped to zero", 10, 20, 0},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ComputeAbsoluteRentBP(tc.surplusValueBP, tc.averageProfitBP)
			if got != tc.want {
				t.Errorf("ComputeAbsoluteRentBP(%d, %d) = %d, want %d",
					tc.surplusValueBP, tc.averageProfitBP, got, tc.want)
			}
		})
	}
}

// TestComputeCompositionRatioBP verifies C/V×10000 (round-half-up); 0 when V=0.
// Fixture: agricultural C=60,V=40 → 60/40×10000=15000;
//
//	manufacturing C=80,V=20 → 80/20×10000=40000;
//	V=0 → 0 (guard).
func TestComputeCompositionRatioBP(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name               string
		constantCapitalBP  int64
		variableCapitalBP  int64
		want               int64
	}{
		{"agriculture 60:40 → 15000", 60, 40, 15000},
		{"manufacturing 80:20 → 40000", 80, 20, 40000},
		{"variable capital zero → guard 0", 50, 0, 0},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ComputeCompositionRatioBP(tc.constantCapitalBP, tc.variableCapitalBP)
			if got != tc.want {
				t.Errorf("ComputeCompositionRatioBP(%d, %d) = %d, want %d",
					tc.constantCapitalBP, tc.variableCapitalBP, got, tc.want)
			}
		})
	}
}

// TestComputeValuePriceGap verifies value − price-of-production (unclamped).
// Fixture: agricultural product value=140, price-of-production=120 → gap=20;
//
//	at parity → 0; inverted (below average composition) → negative gap.
func TestComputeValuePriceGap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                           string
		valueLabourMinutes             int64
		priceOfProductionLabourMinutes int64
		want                           int64
	}{
		{"value above price: positive gap 20", 140, 120, 20},
		{"value equals price: gap zero", 120, 120, 0},
		{"value below price: negative gap", 100, 140, -40},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ComputeValuePriceGap(tc.valueLabourMinutes, tc.priceOfProductionLabourMinutes)
			if got != tc.want {
				t.Errorf("ComputeValuePriceGap(%d, %d) = %d, want %d",
					tc.valueLabourMinutes, tc.priceOfProductionLabourMinutes, got, tc.want)
			}
		})
	}
}

// TestCh45NewIDsDistinct verifies that each ID constructor produces 24-char hex
// strings and that two consecutive calls differ (no repeated values).
func TestCh45NewIDsDistinct(t *testing.T) {
	t.Parallel()

	t.Run("AgriculturalCompositionID", func(t *testing.T) {
		t.Parallel()
		a := NewAgriculturalCompositionID()
		b := NewAgriculturalCompositionID()
		if len(string(a)) != 24 {
			t.Errorf("len(AgriculturalCompositionID) = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two consecutive AgriculturalCompositionIDs must differ")
		}
	})

	t.Run("ValuePriceGapID", func(t *testing.T) {
		t.Parallel()
		a := NewValuePriceGapID()
		b := NewValuePriceGapID()
		if len(string(a)) != 24 {
			t.Errorf("len(ValuePriceGapID) = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two consecutive ValuePriceGapIDs must differ")
		}
	})

	t.Run("AbsoluteRentID", func(t *testing.T) {
		t.Parallel()
		a := NewAbsoluteRentID()
		b := NewAbsoluteRentID()
		if len(string(a)) != 24 {
			t.Errorf("len(AbsoluteRentID) = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two consecutive AbsoluteRentIDs must differ")
		}
	})

	t.Run("AbsoluteRentLimitID", func(t *testing.T) {
		t.Parallel()
		a := NewAbsoluteRentLimitID()
		b := NewAbsoluteRentLimitID()
		if len(string(a)) != 24 {
			t.Errorf("len(AbsoluteRentLimitID) = %d, want 24", len(string(a)))
		}
		if a == b {
			t.Error("two consecutive AbsoluteRentLimitIDs must differ")
		}
	})
}
