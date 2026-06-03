package rent

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// industrialSpheres is the Part II equalisation pool: two competing industrial
// branches of differing organic composition, same s′ = 100%. ΣS = 60, ΣC = 200,
// so the general rate p̄′ = 30%.
func industrialSpheres() []avgprofit.ProductionSphere {
	return []avgprofit.ProductionSphere{
		avgprofit.ComputeProductionSphere("textiles", 100, 20, 10000),  // s = 20
		avgprofit.ComputeProductionSphere("machinery", 100, 40, 10000), // s = 40
	}
}

// TestRentIsSurplusProfitAboveAverageProfit pins the first half of the Part VI
// law (#368): ground-rent is the agricultural surplus-value in excess of the
// average profit the capital draws at the general rate.
func TestRentIsSurplusProfitAboveAverageProfit(t *testing.T) {
	t.Parallel()
	// p̄′ = ΣS/ΣC = 60/200 = 30% over industry; agriculture advances C = 100 and
	// produces s = 50 (a lower-composition, naturally-favoured branch).
	excl := ComputeGeneralRateExclusion(industrialSpheres(), 100, 50)
	if excl.GeneralRateBP != 3000 {
		t.Fatalf("general rate = %d bp, want 3000 (30%%)", excl.GeneralRateBP)
	}
	if excl.AverageProfitAtGeneralRate != 30 {
		t.Errorf("average profit = %d, want 30 (C 100 × 30%%)", excl.AverageProfitAtGeneralRate)
	}
	// Rent = surplus-profit above average profit = 50 − 30 = 20.
	if excl.RentLabourMinutes != 20 {
		t.Errorf("rent = %d, want 20 (agri surplus 50 − average profit 30)", excl.RentLabourMinutes)
	}
}

// TestRentDoesNotEnterGeneralRateEqualisation pins the defining Part VI law: the
// agricultural surplus-profit is barred from the general-rate pool. The rate is
// computed over industry alone; had agriculture entered, the rate would rise and
// the surplus-profit would be competed away.
func TestRentDoesNotEnterGeneralRateEqualisation(t *testing.T) {
	t.Parallel()
	industry := industrialSpheres()
	excl := ComputeGeneralRateExclusion(industry, 100, 50)

	// The law, made explicit.
	if excl.RentEntersGeneralRate {
		t.Error("RentEntersGeneralRate = true, want false (rent is excluded from equalisation)")
	}

	// The exclusion rate is exactly the Part II general rate over industry alone —
	// the two are wired, not independent.
	industryRate := int64(avgprofit.ComputeGeneralProfitRate(industry).Rate)
	if excl.GeneralRateBP != industryRate {
		t.Errorf("exclusion rate %d != industry-only general rate %d", excl.GeneralRateBP, industryRate)
	}

	// Proof the exclusion is load-bearing: had the agricultural surplus-value
	// entered the pool, the general rate would rise — competing the surplus-profit
	// away. Excluding it is exactly what lets landed property capture it as rent.
	agriSphere := avgprofit.ComputeProductionSphere("agriculture", 100, 50, 10000) // s = 50
	pooled := append(append([]avgprofit.ProductionSphere{}, industry...), agriSphere)
	pooledRate := int64(avgprofit.ComputeGeneralProfitRate(pooled).Rate)
	if pooledRate <= industryRate {
		t.Fatalf("pooled rate %d should exceed industry rate %d (agri surplus-profit would lift p̄′)", pooledRate, industryRate)
	}

	// The rent the law preserves is the surplus-value above average profit at the
	// UNCHANGED industry rate — not at the pooled rate that would erase it.
	if excl.RentLabourMinutes != excl.AgriculturalSurplusValue-excl.AverageProfitAtGeneralRate {
		t.Errorf("rent %d != agri surplus %d − average profit %d", excl.RentLabourMinutes, excl.AgriculturalSurplusValue, excl.AverageProfitAtGeneralRate)
	}
}

// TestWorstSoilNoRent: on the regulating worst soil the agricultural capital
// produces no surplus-value above average profit, so there is no surplus-profit
// and no differential rent (clamped at 0).
func TestWorstSoilNoRent(t *testing.T) {
	t.Parallel()
	excl := ComputeGeneralRateExclusion(industrialSpheres(), 100, 30) // s == average profit
	if excl.RentLabourMinutes != 0 {
		t.Errorf("rent = %d, want 0 (worst soil yields no surplus-profit)", excl.RentLabourMinutes)
	}
}
