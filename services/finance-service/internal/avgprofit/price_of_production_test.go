package avgprofit

import (
	"errors"
	"testing"
)

// TestComputePriceOfProduction_Formula_300k1500 verifies the general price formula:
// k=300, rate=1500 → avgProfit = 300·1500/10000 = 45; price = 300+45 = 345.
func TestComputePriceOfProduction_Formula_300k1500(t *testing.T) {
	t.Parallel()
	pop := ComputePriceOfProduction("", 300, 1500, 0)

	if pop.Price != 345 {
		t.Errorf("Price = %d, want 345", pop.Price)
	}
	if pop.AverageProfit != 45 {
		t.Errorf("AverageProfit = %d, want 45", pop.AverageProfit)
	}
}

// TestComputePriceOfProduction_Formula_100k2200 verifies the Ch. 9 canonical calculation:
// k=100, rate=2200 → avgProfit = 100·2200/10000 = 22; price = 100+22 = 122.
func TestComputePriceOfProduction_Formula_100k2200(t *testing.T) {
	t.Parallel()
	pop := ComputePriceOfProduction("", 100, 2200, 0)

	if pop.Price != 122 {
		t.Errorf("Price = %d, want 122", pop.Price)
	}
	if pop.AverageProfit != 22 {
		t.Errorf("AverageProfit = %d, want 22", pop.AverageProfit)
	}
	if pop.CostPrice != 100 {
		t.Errorf("CostPrice = %d, want 100", pop.CostPrice)
	}
	if pop.GeneralRate != 2200 {
		t.Errorf("GeneralRate = %d, want 2200", pop.GeneralRate)
	}
}

// TestComputePriceOfProduction_SphereI verifies Sphere I (80c+20v): value=120, price=122,
// deviation=+2, composition=KindHigher (price > value → receiver).
func TestComputePriceOfProduction_SphereI(t *testing.T) {
	t.Parallel()
	// k=100, rate=2200, value=120 (C+s = 100+20).
	pop := ComputePriceOfProduction("Sphere I (80c+20v)", 100, 2200, 120)

	if pop.Price != 122 {
		t.Errorf("Price = %d, want 122", pop.Price)
	}
	if pop.Deviation != 2 {
		t.Errorf("Deviation = %d, want +2", pop.Deviation)
	}
	if pop.Composition != KindHigher {
		t.Errorf("Composition = %q, want %q (deviation > 0)", pop.Composition, KindHigher)
	}
}

// TestComputePriceOfProduction_SphereIII verifies Sphere III (60c+40v): value=140, price=122,
// deviation=-18, composition=KindLower (price < value → donor).
func TestComputePriceOfProduction_SphereIII(t *testing.T) {
	t.Parallel()
	// k=100, rate=2200, value=140 (C+s = 100+40).
	pop := ComputePriceOfProduction("Sphere III (60c+40v)", 100, 2200, 140)

	if pop.Price != 122 {
		t.Errorf("Price = %d, want 122", pop.Price)
	}
	if pop.Deviation != -18 {
		t.Errorf("Deviation = %d, want -18", pop.Deviation)
	}
	if pop.Composition != KindLower {
		t.Errorf("Composition = %q, want %q (deviation < 0)", pop.Composition, KindLower)
	}
}

// TestComputePriceOfProduction_AverageCase verifies a sphere where value equals price:
// deviation=0 → KindAverage.
func TestComputePriceOfProduction_AverageCase(t *testing.T) {
	t.Parallel()
	// k=100, rate=2200, price=122, value=122: deviation=0.
	pop := ComputePriceOfProduction("Average sphere", 100, 2200, 122)

	if pop.Deviation != 0 {
		t.Errorf("Deviation = %d, want 0", pop.Deviation)
	}
	if pop.Composition != KindAverage {
		t.Errorf("Composition = %q, want %q (deviation == 0)", pop.Composition, KindAverage)
	}
}

// TestComputePriceOfProduction_ValueZero verifies that when value<=0, Composition is "".
func TestComputePriceOfProduction_ValueZero(t *testing.T) {
	t.Parallel()
	pop := ComputePriceOfProduction("unknown value sphere", 100, 2200, 0)

	if pop.Composition != CompositionKind("") {
		t.Errorf("Composition = %q, want empty string when value=0", pop.Composition)
	}
}

// TestComputePriceOfProduction_AllFiveSpheres verifies all five spheres from the
// Ch. 9 fixture table against verbatim numbers.
func TestComputePriceOfProduction_AllFiveSpheres(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		v           int64
		value       CommodityValue
		wantDev     ValueDeviation
		wantKind    CompositionKind
	}{
		{"Sphere I (80c+20v)", 20, 120, 2, KindHigher},
		{"Sphere II (70c+30v)", 30, 130, -8, KindLower},
		{"Sphere III (60c+40v)", 40, 140, -18, KindLower},
		{"Sphere IV (85c+15v)", 15, 115, 7, KindHigher},
		{"Sphere V (95c+5v)", 5, 105, 17, KindHigher},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			pop := ComputePriceOfProduction(tc.name, 100, 2200, tc.value)

			if pop.Price != 122 {
				t.Errorf("[%s] Price = %d, want 122", tc.name, pop.Price)
			}
			if pop.AverageProfit != 22 {
				t.Errorf("[%s] AverageProfit = %d, want 22", tc.name, pop.AverageProfit)
			}
			if pop.Deviation != tc.wantDev {
				t.Errorf("[%s] Deviation = %d, want %d", tc.name, pop.Deviation, tc.wantDev)
			}
			if pop.Composition != tc.wantKind {
				t.Errorf("[%s] Composition = %q, want %q", tc.name, pop.Composition, tc.wantKind)
			}
		})
	}
}

// TestComputeValuePriceDeviation_SphereIII verifies the canonical fixture from the
// design contract: ComputeValuePriceDeviation(140, 122).Deviation == -18.
func TestComputeValuePriceDeviation_SphereIII(t *testing.T) {
	t.Parallel()
	dev := ComputeValuePriceDeviation(140, 122)

	if dev.Deviation != -18 {
		t.Errorf("Deviation = %d, want -18", dev.Deviation)
	}
	if dev.Value != 140 {
		t.Errorf("Value = %d, want 140", dev.Value)
	}
	if dev.PriceOfProduction != 122 {
		t.Errorf("PriceOfProduction = %d, want 122", dev.PriceOfProduction)
	}
}

// TestComputeValuePriceDeviation_SphereI verifies: ComputeValuePriceDeviation(120, 122).Deviation == 2.
func TestComputeValuePriceDeviation_SphereI(t *testing.T) {
	t.Parallel()
	dev := ComputeValuePriceDeviation(120, 122)

	if dev.Deviation != 2 {
		t.Errorf("Deviation = %d, want +2", dev.Deviation)
	}
}

// TestComputeValuePriceDeviation_Zero verifies zero deviation when value == price.
func TestComputeValuePriceDeviation_Zero(t *testing.T) {
	t.Parallel()
	dev := ComputeValuePriceDeviation(122, 122)

	if dev.Deviation != 0 {
		t.Errorf("Deviation = %d, want 0 when value == price", dev.Deviation)
	}
}

// TestPriceOfProduction_Validate_Valid ensures a correctly constructed record passes.
func TestPriceOfProduction_Validate_Valid(t *testing.T) {
	t.Parallel()
	pop := ComputePriceOfProduction("Sphere I", 100, 2200, 120)
	if err := pop.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

// TestPriceOfProduction_Validate_NegativeCostPrice ensures negative CostPrice returns
// ErrNegativeMagnitude.
func TestPriceOfProduction_Validate_NegativeCostPrice(t *testing.T) {
	t.Parallel()
	pop := PriceOfProduction{CostPrice: -1, GeneralRate: 2200, CommodityValue: 120}
	if err := pop.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("Validate() = %v, want ErrNegativeMagnitude", err)
	}
}

// TestPriceOfProduction_Validate_NegativeGeneralRate ensures negative GeneralRate returns
// ErrNegativeMagnitude.
func TestPriceOfProduction_Validate_NegativeGeneralRate(t *testing.T) {
	t.Parallel()
	pop := PriceOfProduction{CostPrice: 100, GeneralRate: -1, CommodityValue: 120}
	if err := pop.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("Validate() = %v, want ErrNegativeMagnitude", err)
	}
}

// TestPriceOfProduction_Validate_NegativeCommodityValue ensures negative CommodityValue
// returns ErrNegativeMagnitude.
func TestPriceOfProduction_Validate_NegativeCommodityValue(t *testing.T) {
	t.Parallel()
	pop := PriceOfProduction{CostPrice: 100, GeneralRate: 2200, CommodityValue: -1}
	if err := pop.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("Validate() = %v, want ErrNegativeMagnitude", err)
	}
}

// TestPriceOfProduction_Price_Invariant verifies price = k + k·rate/10000 for
// several combinations.
func TestPriceOfProduction_Price_Invariant(t *testing.T) {
	t.Parallel()
	cases := []struct {
		k     CostPrice
		rate  SphereProfitRate
		price ProductionPrice
	}{
		{100, 2200, 122},
		{300, 1500, 345},
		{500, 2000, 600},
		{0, 2200, 0},
	}
	for _, tc := range cases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			pop := ComputePriceOfProduction("", tc.k, tc.rate, 0)
			if pop.Price != tc.price {
				t.Errorf("ComputePriceOfProduction(%d,%d).Price = %d, want %d",
					tc.k, tc.rate, pop.Price, tc.price)
			}
		})
	}
}
