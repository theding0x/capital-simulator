package avgprofit

import (
	"testing"
)

// fiveSpheres returns the canonical five production spheres from Ch. 9.
// Each sphere: C=100, s'=100%, so s = v; k = C = 100.
// ΣS = 110, ΣC = 500 → general rate = 110·10000/500 = 2200 bp.
func fiveSpheres() []ProductionSphere {
	return []ProductionSphere{
		ComputeProductionSphere("Sphere I (80c+20v)", 100, 20, 10000),
		ComputeProductionSphere("Sphere II (70c+30v)", 100, 30, 10000),
		ComputeProductionSphere("Sphere III (60c+40v)", 100, 40, 10000),
		ComputeProductionSphere("Sphere IV (85c+15v)", 100, 15, 10000),
		ComputeProductionSphere("Sphere V (95c+5v)", 100, 5, 10000),
	}
}

// TestComputeGeneralProfitRate_FiveSpheres verifies the five-sphere fixture from
// Ch. 9 fixture table: ΣS=110, ΣC=500 → Rate=2200, AverageVariablePercent=22.
func TestComputeGeneralProfitRate_FiveSpheres(t *testing.T) {
	t.Parallel()
	g := ComputeGeneralProfitRate(fiveSpheres())

	if g.Rate != 2200 {
		t.Errorf("Rate = %d, want 2200", g.Rate)
	}
	if g.SumSurplusValues != 110 {
		t.Errorf("SumSurplusValues = %d, want 110", g.SumSurplusValues)
	}
	if g.SumTotalCapitals != 500 {
		t.Errorf("SumTotalCapitals = %d, want 500", g.SumTotalCapitals)
	}
	if g.AverageVariablePercent != 22 {
		t.Errorf("AverageVariablePercent = %d, want 22", g.AverageVariablePercent)
	}
}

// TestComputeAverageProfit_FiveSphereRate verifies the canonical Ch. 9 calculation:
// 100 · 2200 / 10000 = 22.
func TestComputeAverageProfit_FiveSphereRate(t *testing.T) {
	t.Parallel()
	ap := ComputeAverageProfit(100, 2200)

	if ap.Amount != 22 {
		t.Errorf("Amount = %d, want 22", ap.Amount)
	}
	if ap.CapitalAdvanced != 100 {
		t.Errorf("CapitalAdvanced = %d, want 100", ap.CapitalAdvanced)
	}
	if ap.GeneralRate != 2200 {
		t.Errorf("GeneralRate = %d, want 2200", ap.GeneralRate)
	}
}

// TestComputeSocialCapitalAggregate_FiveSpheres verifies the conservation laws:
// Σ values = 610 = Σ prices; ValueConserved true; Σ average profits = Σ surplus-values = 110;
// Σ deviations = +2 − 8 − 18 + 7 + 17 = 0; WeightedGeneralRate = 2200.
func TestComputeSocialCapitalAggregate_FiveSpheres(t *testing.T) {
	t.Parallel()
	agg := ComputeSocialCapitalAggregate(fiveSpheres())

	if agg.TotalValues != 610 {
		t.Errorf("TotalValues = %d, want 610", agg.TotalValues)
	}
	if agg.TotalPricesOfProduction != 610 {
		t.Errorf("TotalPricesOfProduction = %d, want 610", agg.TotalPricesOfProduction)
	}
	if !agg.ValueConserved {
		t.Error("ValueConserved = false, want true")
	}
	if agg.SumAverageProfits != 110 {
		t.Errorf("SumAverageProfits = %d, want 110 (= SumSurplusValues)", agg.SumAverageProfits)
	}
	if agg.SumSurplusValues != 110 {
		t.Errorf("SumSurplusValues = %d, want 110", agg.SumSurplusValues)
	}
	if agg.SumAverageProfits != Profit(agg.SumSurplusValues) {
		t.Errorf("SumAverageProfits %d != SumSurplusValues %d (conservation broken)",
			agg.SumAverageProfits, agg.SumSurplusValues)
	}
	if !agg.ProfitConserved {
		t.Errorf("ProfitConserved = false, want true (Σprofit %d == Σsurplus-value %d)",
			agg.SumAverageProfits, agg.SumSurplusValues)
	}
	if !agg.ValueConserved {
		t.Error("ValueConserved = false, want true (Σprice == Σvalue)")
	}
	if agg.WeightedGeneralRate != 2200 {
		t.Errorf("WeightedGeneralRate = %d, want 2200", agg.WeightedGeneralRate)
	}

	// Σ deviations across five prices of production must equal zero.
	if len(agg.PricesOfProduction) != 5 {
		t.Fatalf("len(PricesOfProduction) = %d, want 5", len(agg.PricesOfProduction))
	}
	var sumDev ValueDeviation
	for _, p := range agg.PricesOfProduction {
		sumDev += p.Deviation
	}
	if sumDev != 0 {
		t.Errorf("sum of deviations = %d, want 0", sumDev)
	}
}

// TestComputeSocialCapitalAggregate_WeightedFixture checks the design contract's
// weighted fixture: A=200(v50), B=300(v120), C=1000(v150), D=4000(v400), all s'=100%.
// ΣS = 720, ΣC = 5500 → WeightedGeneralRate = round_half_up(720·10000/5500) = 1309.
// The naive simple average of individual rates (not capital-weighted) would be ≈2250 bp —
// verifying the capital-weighted rate is the chapter's core argument.
func TestComputeSocialCapitalAggregate_WeightedFixture(t *testing.T) {
	t.Parallel()
	spheres := []ProductionSphere{
		ComputeProductionSphere("A (v50)", 200, 50, 10000),
		ComputeProductionSphere("B (v120)", 300, 120, 10000),
		ComputeProductionSphere("C (v150)", 1000, 150, 10000),
		ComputeProductionSphere("D (v400)", 4000, 400, 10000),
	}
	agg := ComputeSocialCapitalAggregate(spheres)

	if agg.WeightedGeneralRate != 1309 {
		t.Errorf("WeightedGeneralRate = %d, want 1309 (round-half-up of 720·10000/5500 = 1309.09)", agg.WeightedGeneralRate)
	}
}

// TestComputeGeneralProfitRate_EmptySlice ensures an empty input is valid with rate 0.
func TestComputeGeneralProfitRate_EmptySlice(t *testing.T) {
	t.Parallel()
	g := ComputeGeneralProfitRate(nil)

	if g.Rate != 0 {
		t.Errorf("Rate = %d, want 0 for empty spheres", g.Rate)
	}
	if g.SumSurplusValues != 0 {
		t.Errorf("SumSurplusValues = %d, want 0", g.SumSurplusValues)
	}
	if g.SumTotalCapitals != 0 {
		t.Errorf("SumTotalCapitals = %d, want 0", g.SumTotalCapitals)
	}
}

// TestComputeGeneralProfitRate_ZeroVariable confirms a sphere with V=0 contributes
// zero surplus-value to the total.
func TestComputeGeneralProfitRate_ZeroVariable(t *testing.T) {
	t.Parallel()
	spheres := []ProductionSphere{
		ComputeProductionSphere("All-constant (100c+0v)", 100, 0, 10000),
	}
	g := ComputeGeneralProfitRate(spheres)

	if g.SumSurplusValues != 0 {
		t.Errorf("SumSurplusValues = %d, want 0 when V=0", g.SumSurplusValues)
	}
	if g.Rate != 0 {
		t.Errorf("Rate = %d, want 0 when only zero-variable sphere", g.Rate)
	}
}

// TestGeneralProfitRate_Validate_Valid verifies valid five-sphere input passes.
func TestGeneralProfitRate_Validate_Valid(t *testing.T) {
	t.Parallel()
	g := ComputeGeneralProfitRate(fiveSpheres())
	if err := g.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

// TestGeneralProfitRate_Validate_EmptySpheresOK verifies empty sphere slice is valid.
func TestGeneralProfitRate_Validate_EmptySpheresOK(t *testing.T) {
	t.Parallel()
	g := GeneralProfitRate{}
	if err := g.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for empty spheres", err)
	}
}

// TestAverageProfit_Amount_Invariant verifies Amount = C·rate/10000 for multiple cases.
func TestAverageProfit_Amount_Invariant(t *testing.T) {
	t.Parallel()
	cases := []struct {
		c    TotalCapital
		rate SphereProfitRate
		want Profit
	}{
		{100, 2200, 22},
		{300, 1500, 45},
		{500, 2000, 100},
		{0, 2200, 0},
	}
	for _, tc := range cases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			ap := ComputeAverageProfit(tc.c, tc.rate)
			if ap.Amount != tc.want {
				t.Errorf("ComputeAverageProfit(%d, %d).Amount = %d, want %d",
					tc.c, tc.rate, ap.Amount, tc.want)
			}
		})
	}
}
