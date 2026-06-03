package rent

import (
	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// basisPoints is the fixed-point scale (1% = 100 bp), mirroring avgprofit, used
// for the average-profit-at-the-general-rate computation below.
const basisPoints = 10000

// GeneralRateExclusion expresses the defining law of Part VI: ground-rent is the
// agricultural surplus-profit that the landed-property barrier bars from entering
// the general-rate-of-profit equalisation.
//
// The general rate p̄′ is fixed by industry (Part II, Ch. 9): competing
// industrial capitals pool their surplus-value and each draws average profit at
// p̄′ = ΣS / ΣC. Agricultural capital sits behind the landed-property barrier —
// it draws the same average profit at p̄′ on its advance, but the EXCESS of the
// surplus-value it actually produces over that average profit is NOT thrown back
// into the pool. Were it pooled, competition would equalise it away and lift p̄′;
// landed property prevents that, capturing the excess — the surplus-profit — as
// ground-rent. That is precisely why agricultural surplus-profit is *not*
// competed away into p̄′: rent does not enter the equalisation.
type GeneralRateExclusion struct {
	// GeneralRateBP is p̄′ in basis points, computed over the INDUSTRIAL spheres
	// alone (agriculture never enters the equalisation).
	GeneralRateBP int64 `json:"general_rate_bp"`
	// AgriculturalCapital is C_agri (LabourMinutes) advanced behind the barrier.
	AgriculturalCapital int64 `json:"agricultural_capital"`
	// AgriculturalSurplusValue is s_agri actually produced (LabourMinutes).
	AgriculturalSurplusValue int64 `json:"agricultural_surplus_value"`
	// AverageProfitAtGeneralRate is C_agri × p̄′ — all the agricultural capital
	// would receive if its surplus-profit entered the pool.
	AverageProfitAtGeneralRate int64 `json:"average_profit_at_general_rate"`
	// RentLabourMinutes is the surplus-profit above average profit
	// (s_agri − average profit) captured as ground-rent.
	RentLabourMinutes int64 `json:"rent_labour_minutes"`
	// RentEntersGeneralRate is always false: the law made explicit.
	RentEntersGeneralRate bool `json:"rent_enters_general_rate"`
}

// AverageProfitAtRate returns the average profit a capital draws at the general
// rate: capital × rateBP / 10000 (integer truncation, matching avgprofit's
// ComputeAverageProfit).
func AverageProfitAtRate(capital, rateBP int64) int64 {
	return capital * rateBP / basisPoints
}

// ComputeGeneralRateExclusion wires the Part VI rent law to the Part II general
// rate. It computes p̄′ over the INDUSTRIAL spheres alone — agriculture's
// surplus-profit never enters the equalisation — then derives the agricultural
// capital's average profit at p̄′ and the rent it pays: the surplus-profit above
// average profit that landed property excludes from the pool. Rent is clamped at
// 0 (the regulating worst soil yields no differential surplus-profit, hence no
// differential rent).
func ComputeGeneralRateExclusion(industrialSpheres []avgprofit.ProductionSphere, agriculturalCapital, agriculturalSurplusValue int64) GeneralRateExclusion {
	g := avgprofit.ComputeGeneralProfitRate(industrialSpheres)
	rate := int64(g.Rate)
	avgProfit := AverageProfitAtRate(agriculturalCapital, rate)
	rent := agriculturalSurplusValue - avgProfit
	if rent < 0 {
		rent = 0
	}
	return GeneralRateExclusion{
		GeneralRateBP:              rate,
		AgriculturalCapital:        agriculturalCapital,
		AgriculturalSurplusValue:   agriculturalSurplusValue,
		AverageProfitAtGeneralRate: avgProfit,
		RentLabourMinutes:          rent,
		RentEntersGeneralRate:      false,
	}
}
