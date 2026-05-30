package avgprofit

import "testing"

// fiveSpheres (the canonical Ch. 9 five-sphere fixture) is declared in
// general_rate_test.go and reused here.

// fivePartIIPrices builds the matching prices of production: k=100, general rate
// 2200, each price 122, values 120/130/140/115/105 → Σ values = Σ prices = 610.
func fivePartIIPrices() []PriceOfProduction {
	return []PriceOfProduction{
		ComputePriceOfProduction("Sphere I (80c+20v)", 100, 2200, 120),
		ComputePriceOfProduction("Sphere II (70c+30v)", 100, 2200, 130),
		ComputePriceOfProduction("Sphere III (60c+40v)", 100, 2200, 140),
		ComputePriceOfProduction("Sphere IV (85c+15v)", 100, 2200, 115),
		ComputePriceOfProduction("Sphere V (95c+5v)", 100, 2200, 105),
	}
}

// TestComputePartIISummary_FiveSphere verifies the Part II conservation snapshot
// against the Ch. 9 five-sphere fixture: Σ values == Σ prices == 610, the
// general rate is 2200, and surplus-value produced (110) equals average profit
// drawn (110) for the social aggregate.
func TestComputePartIISummary_FiveSphere(t *testing.T) {
	t.Parallel()
	summary := ComputePartIISummary(fiveSpheres(), fivePartIIPrices(), nil)

	if summary.TotalValues != 610 {
		t.Errorf("TotalValues = %d, want 610", summary.TotalValues)
	}
	if summary.TotalPricesOfProduction != 610 {
		t.Errorf("TotalPricesOfProduction = %d, want 610", summary.TotalPricesOfProduction)
	}
	if !summary.ValueConserved {
		t.Error("ValueConserved = false, want true (610 == 610)")
	}
	if summary.GeneralRate != 2200 {
		t.Errorf("GeneralRate = %d, want 2200", summary.GeneralRate)
	}
	if !summary.AverageCompositionCheck.SurplusValueEqualsAverageProfit {
		t.Error("SurplusValueEqualsAverageProfit = false, want true")
	}
	if summary.AverageCompositionCheck.TotalSurplusValue != 110 {
		t.Errorf("TotalSurplusValue = %d, want 110", summary.AverageCompositionCheck.TotalSurplusValue)
	}
	if summary.AverageCompositionCheck.TotalAverageProfit != 110 {
		t.Errorf("TotalAverageProfit = %d, want 110", summary.AverageCompositionCheck.TotalAverageProfit)
	}
	if summary.SphereCount != 5 {
		t.Errorf("SphereCount = %d, want 5", summary.SphereCount)
	}
	if summary.PricesOfProductionCount != 5 {
		t.Errorf("PricesOfProductionCount = %d, want 5", summary.PricesOfProductionCount)
	}
	if summary.WageEffectCount != 0 {
		t.Errorf("WageEffectCount = %d, want 0", summary.WageEffectCount)
	}
}

// TestComputePartIISummary_Conserved asserts the invariant: when Σ prices ==
// Σ values, both ValueConserved and SurplusValueEqualsAverageProfit are true and
// TotalAverageProfit == TotalSurplusValue.
func TestComputePartIISummary_Conserved(t *testing.T) {
	t.Parallel()
	summary := ComputePartIISummary(fiveSpheres(), fivePartIIPrices(), nil)
	if int64(summary.TotalPricesOfProduction) != int64(summary.TotalValues) {
		t.Fatalf("fixture not conserved: prices %d != values %d",
			summary.TotalPricesOfProduction, summary.TotalValues)
	}
	if !summary.ValueConserved {
		t.Error("ValueConserved = false, want true")
	}
	if int64(summary.AverageCompositionCheck.TotalAverageProfit) != int64(summary.AverageCompositionCheck.TotalSurplusValue) {
		t.Errorf("avg profit %d != surplus %d when conserved",
			summary.AverageCompositionCheck.TotalAverageProfit,
			summary.AverageCompositionCheck.TotalSurplusValue)
	}
}

// TestComputePartIISummary_Empty: nil inputs must not panic; with no prices,
// 0 == 0 so value is conserved and every count is zero.
func TestComputePartIISummary_Empty(t *testing.T) {
	t.Parallel()
	summary := ComputePartIISummary(nil, nil, nil)
	if summary.TotalValues != 0 {
		t.Errorf("TotalValues = %d, want 0", summary.TotalValues)
	}
	if summary.TotalPricesOfProduction != 0 {
		t.Errorf("TotalPricesOfProduction = %d, want 0", summary.TotalPricesOfProduction)
	}
	if !summary.ValueConserved {
		t.Error("ValueConserved = false, want true (0 == 0)")
	}
	if summary.SphereCount != 0 || summary.PricesOfProductionCount != 0 || summary.WageEffectCount != 0 {
		t.Errorf("counts = %d/%d/%d, want 0/0/0",
			summary.SphereCount, summary.PricesOfProductionCount, summary.WageEffectCount)
	}
}
