package avgprofit

import (
	"errors"
	"testing"
)

// TestComputeMarketValue_CottonBulk verifies the spec fixture:
// ComputeMarketValue("cotton", 100, 80, 130).Value == 100 (bulk = regulating value).
func TestComputeMarketValue_CottonBulk(t *testing.T) {
	t.Parallel()
	mv := ComputeMarketValue("cotton", 100, 80, 130)

	if mv.Value != 100 {
		t.Errorf("Value = %d, want 100 (bulk condition)", mv.Value)
	}
	if mv.BulkConditionValue != 100 {
		t.Errorf("BulkConditionValue = %d, want 100", mv.BulkConditionValue)
	}
	if mv.BestConditionValue != 80 {
		t.Errorf("BestConditionValue = %d, want 80", mv.BestConditionValue)
	}
	if mv.WorstConditionValue != 130 {
		t.Errorf("WorstConditionValue = %d, want 130", mv.WorstConditionValue)
	}
	if mv.SphereName != "cotton" {
		t.Errorf("SphereName = %q, want %q", mv.SphereName, "cotton")
	}
}

// TestComputeMarketValue_ValueEqualsBulk verifies the invariant Value == BulkConditionValue
// for arbitrary inputs.
func TestComputeMarketValue_ValueEqualsBulk(t *testing.T) {
	t.Parallel()
	cases := []struct {
		bulk  CommodityValue
		best  CommodityValue
		worst CommodityValue
	}{
		{50, 40, 60},
		{0, 0, 0},
		{200, 150, 250},
	}
	for _, tc := range cases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			mv := ComputeMarketValue("sphere", tc.bulk, tc.best, tc.worst)
			if mv.Value != tc.bulk {
				t.Errorf("Value = %d, want %d (== BulkConditionValue)", mv.Value, tc.bulk)
			}
		})
	}
}

// TestComputeMarketPrice_Deviation verifies Deviation = amount − market value.
// With amount=95 and marketValue=100: Deviation = -5.
func TestComputeMarketPrice_Deviation(t *testing.T) {
	t.Parallel()
	mp := ComputeMarketPrice(95, 100)

	if mp.Deviation != -5 {
		t.Errorf("Deviation = %d, want -5 (95 − 100)", mp.Deviation)
	}
	if mp.Amount != 95 {
		t.Errorf("Amount = %d, want 95", mp.Amount)
	}
	if mp.MarketValue != 100 {
		t.Errorf("MarketValue = %d, want 100", mp.MarketValue)
	}
}

// TestComputeMarketPrice_AboveValue verifies a positive deviation when market price > value.
func TestComputeMarketPrice_AboveValue(t *testing.T) {
	t.Parallel()
	mp := ComputeMarketPrice(110, 100)

	if mp.Deviation != 10 {
		t.Errorf("Deviation = %d, want +10 (110 − 100)", mp.Deviation)
	}
}

// TestComputeMarketPrice_AtValue verifies zero deviation when price == market value.
func TestComputeMarketPrice_AtValue(t *testing.T) {
	t.Parallel()
	mp := ComputeMarketPrice(100, 100)

	if mp.Deviation != 0 {
		t.Errorf("Deviation = %d, want 0 when price == value", mp.Deviation)
	}
}

// TestComputeSurplusProfit_HighProductivityFirm verifies the spec fixture:
// ComputeSurplusProfit("firm", 60, 100, 1000, 2200).Amount == 40000
// because (100 − 60) × 1000 = 40000 > 0.
func TestComputeSurplusProfit_HighProductivityFirm(t *testing.T) {
	t.Parallel()
	sp := ComputeSurplusProfit("firm", 60, 100, 1000, 2200)

	if sp.Amount != 40000 {
		t.Errorf("Amount = %d, want 40000 ((100-60)*1000)", sp.Amount)
	}
	if sp.Amount <= 0 {
		t.Errorf("Amount = %d, want > 0 (high-productivity firm earns surplus)", sp.Amount)
	}
}

// TestComputeSurplusProfit_BelowAverageFirm verifies that when individual value 120
// exceeds market value 100, Amount is negative (firm earns less than average).
func TestComputeSurplusProfit_BelowAverageFirm(t *testing.T) {
	t.Parallel()
	// individual value 120 > market value 100 → Amount = (100 − 120) × 1000 = -20000.
	sp := ComputeSurplusProfit("below-average firm", 120, 100, 1000, 2200)

	if sp.Amount >= 0 {
		t.Errorf("Amount = %d, want < 0 (individual value 120 > market value 100)", sp.Amount)
	}
	if sp.Amount != -20000 {
		t.Errorf("Amount = %d, want -20000 ((100-120)*1000)", sp.Amount)
	}
}

// TestComputeSurplusProfit_AmountInvariant verifies Amount = (MarketValue − IndividualValue) × OutputQty
// for several cases.
func TestComputeSurplusProfit_AmountInvariant(t *testing.T) {
	t.Parallel()
	cases := []struct {
		indVal     IndividualValue
		marketVal  CommodityValue
		outputQty  int64
		wantAmount Profit
	}{
		{60, 100, 1000, 40000},
		{100, 100, 500, 0},
		{80, 70, 200, -2000},
		{0, 50, 100, 5000},
	}
	for _, tc := range cases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			sp := ComputeSurplusProfit("firm", tc.indVal, tc.marketVal, tc.outputQty, 2200)
			if sp.Amount != tc.wantAmount {
				t.Errorf("Amount = %d, want %d", sp.Amount, tc.wantAmount)
			}
		})
	}
}

// TestMarketValue_Validate_NegativeMagnitude verifies that any negative condition value
// returns ErrNegativeMagnitude.
func TestMarketValue_Validate_NegativeMagnitude(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		bulk  CommodityValue
		best  CommodityValue
		worst CommodityValue
	}{
		{"negative bulk", -1, 80, 130},
		{"negative best", 100, -1, 130},
		{"negative worst", 100, 80, -1},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mv := ComputeMarketValue("cotton", tc.bulk, tc.best, tc.worst)
			if err := mv.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
				t.Errorf("Validate() = %v, want ErrNegativeMagnitude", err)
			}
		})
	}
}

// TestMarketValue_Validate_Valid verifies a correctly constructed record passes validation.
func TestMarketValue_Validate_Valid(t *testing.T) {
	t.Parallel()
	mv := ComputeMarketValue("cotton", 100, 80, 130)
	if err := mv.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

// TestSurplusProfit_Validate_NegativeMagnitude verifies that negative input magnitudes
// return ErrNegativeMagnitude.
func TestSurplusProfit_Validate_NegativeMagnitude(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		indVal      IndividualValue
		marketVal   CommodityValue
		outputQty   int64
		generalRate SphereProfitRate
	}{
		{"negative individual value", -1, 100, 1000, 2200},
		{"negative market value", 60, -1, 1000, 2200},
		{"negative output qty", 60, 100, -1, 2200},
		{"negative general rate", 60, 100, 1000, -1},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sp := SurplusProfit{
				IndividualValue: tc.indVal,
				MarketValue:     tc.marketVal,
				OutputQty:       tc.outputQty,
				GeneralRate:     tc.generalRate,
			}
			if err := sp.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
				t.Errorf("Validate() = %v, want ErrNegativeMagnitude", err)
			}
		})
	}
}

// TestNewMarketValueID_Format verifies that NewMarketValueID returns a non-empty
// 24-character hex string (96-bit = 12 bytes = 24 hex chars).
func TestNewMarketValueID_Format(t *testing.T) {
	t.Parallel()
	id := NewMarketValueID()
	if id.IsZero() {
		t.Error("NewMarketValueID() returned a zero ID")
	}
	if len(string(id)) != 24 {
		t.Errorf("id length = %d, want 24 hex chars", len(string(id)))
	}
}

// TestNewSurplusProfitID_Format verifies that NewSurplusProfitID returns a non-empty
// 24-character hex string.
func TestNewSurplusProfitID_Format(t *testing.T) {
	t.Parallel()
	id := NewSurplusProfitID()
	if id.IsZero() {
		t.Error("NewSurplusProfitID() returned a zero ID")
	}
	if len(string(id)) != 24 {
		t.Errorf("id length = %d, want 24 hex chars", len(string(id)))
	}
}
