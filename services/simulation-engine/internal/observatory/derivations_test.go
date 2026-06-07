package observatory

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// seed returns a baseline AbodeState with Marx-faithful 2:1 c:v ratio and
// 100% surplus-value rate (Vol. I Ch. 9 numerical example).
func seedAbode() simulation.AbodeState {
	return simulation.AbodeState{
		Period:                1,
		ConstantPence:         600000,
		VariablePence:         300000,
		BaseWagePence:         2500,
		WorkerSupply:          150,
		SurplusRateBaseBP:     10000,
		AccumulationRateBP:    5000,
		MarginalCompositionBP: 8000,
		DisplacementRateBP:    1800,
		ProductivityGrowthBP:  500,
		PopulationGrowthBP:    150,
	}
}

func seedReadout(a simulation.AbodeState) simulation.AbodeReadout {
	return a.Readout()
}

// TestDeriveCirculation_Conservation verifies that department totals sum to
// within ±2 pence of sumTotal (integer rounding residuals are acceptable).
func TestDeriveCirculation_Conservation(t *testing.T) {
	t.Parallel()
	abode := seedAbode()
	r := seedReadout(abode)
	sumTotal := int64(900000)
	sumSurplus := int64(90000)

	snap := DeriveCirculation(abode, r, sumSurplus, sumTotal)

	dI := snap.Departments.I
	dII := snap.Departments.II
	got := dI.C + dI.V + dI.S + dII.C + dII.V + dII.S
	diff := got - sumTotal
	if diff < -2 || diff > 2 {
		t.Errorf("dept totals conserve to ±2: got sum=%d, want ~%d (diff=%d)", got, sumTotal, diff)
	}
}

// TestDeriveCirculation_Extended verifies that Extended == (AccumulationRateBP > 0).
func TestDeriveCirculation_Extended(t *testing.T) {
	t.Parallel()
	abode := seedAbode()
	r := seedReadout(abode)

	snapExtended := DeriveCirculation(abode, r, 90000, 900000)
	if !snapExtended.Extended {
		t.Error("Extended should be true when AccumulationRateBP=5000 > 0")
	}

	abodeSimple := abode
	abodeSimple.AccumulationRateBP = 0
	snapSimple := DeriveCirculation(abodeSimple, r, 90000, 900000)
	if snapSimple.Extended {
		t.Error("Extended should be false when AccumulationRateBP=0")
	}
}

// TestDeriveDistribution_InterestClamped verifies that InterestRateBP stays
// within [120, 1600] for a range of generalRateBP values including extremes.
func TestDeriveDistribution_InterestClamped(t *testing.T) {
	t.Parallel()
	abode := seedAbode()
	r := seedReadout(abode)
	periods := []simulation.GeneralLawPeriod{}

	cases := []int64{0, 1, 100, 1000, 5000, 10000, 50000, 100000}
	for _, gRate := range cases {
		dist := DeriveDistribution(periods, r, 90000, 810000, 900000, gRate)
		if dist.InterestRateBP < 120 || dist.InterestRateBP > 1600 {
			t.Errorf("InterestRateBP=%d out of [120,1600] for generalRateBP=%d",
				dist.InterestRateBP, gRate)
		}
	}
}

// TestDeriveDistribution_TrinityWagesAndProfit verifies that the trinity
// Wages+Profit equals TotalVariablePence+TotalSurplusPence, and that Rent is now
// the derived rent-terrain total (the honest zero closed by Ch.37–47).
func TestDeriveDistribution_TrinityWagesAndProfit(t *testing.T) {
	t.Parallel()
	abode := seedAbode()
	r := seedReadout(abode)

	dist := DeriveDistribution(nil, r, 90000, 810000, 900000, 1000)

	wantTotal := r.TotalVariablePence + r.TotalSurplusPence
	gotTotal := dist.Trinity.WagesPence + dist.Trinity.ProfitPence
	if gotTotal != wantTotal {
		t.Errorf("Trinity Wages+Profit=%d, want %d (TotalVariable+TotalSurplus)",
			gotTotal, wantTotal)
	}
	if dist.Trinity.RentPence != dist.Rent.TotalRentPence {
		t.Errorf("Trinity RentPence=%d, want %d (the rent-terrain total)",
			dist.Trinity.RentPence, dist.Rent.TotalRentPence)
	}
	if dist.Trinity.RentPence <= 0 {
		t.Errorf("Trinity RentPence=%d, want > 0 (rent now modelled)", dist.Trinity.RentPence)
	}
}

// TestDeriveRent_DifferentialOverWaterline verifies Marx's Differential Rent I
// (Ch.39): the worst soil regulates the price (the waterline) and yields no
// differential rent; each more fertile soil yields strictly more differential
// rent than the last. With yields 1:2:3:4 the differentials run 0 : w : 2w : 3w.
func TestDeriveRent_DifferentialOverWaterline(t *testing.T) {
	t.Parallel()
	abode := seedAbode()
	r := seedReadout(abode)

	rent := DeriveRent(r, 810000, 1000, 500)

	if len(rent.Grades) != 4 {
		t.Fatalf("expected 4 soil grades, got %d", len(rent.Grades))
	}
	worst := rent.Grades[0]
	if !worst.Regulating {
		t.Error("the worst soil (grade A) must be the regulating soil")
	}
	if worst.DifferentialRentPence != 0 {
		t.Errorf("worst soil differential rent=%d, want 0", worst.DifferentialRentPence)
	}
	// Differential rent strictly increases with fertility.
	for i := 1; i < len(rent.Grades); i++ {
		if rent.Grades[i].DifferentialRentPence <= rent.Grades[i-1].DifferentialRentPence {
			t.Errorf("differential rent not strictly rising at grade %d: %d <= %d",
				i, rent.Grades[i].DifferentialRentPence, rent.Grades[i-1].DifferentialRentPence)
		}
	}
	// The differentials are multiples of the waterline (0, w, 2w, 3w).
	for i, g := range rent.Grades {
		want := int64(i) * rent.WaterlinePence
		if g.DifferentialRentPence != want {
			t.Errorf("grade %s differential=%d, want %d (=%d×waterline)",
				g.ID, g.DifferentialRentPence, want, i)
		}
	}
}

// TestDeriveRent_AbsoluteOnEverySoil verifies that absolute rent (from
// agriculture's lower organic composition, Ch.45) is positive and exacted
// equally on every soil — including the worst, which pays no differential rent.
func TestDeriveRent_AbsoluteOnEverySoil(t *testing.T) {
	t.Parallel()
	abode := seedAbode()
	r := seedReadout(abode)

	rent := DeriveRent(r, 810000, 1000, 500)

	ar := rent.Grades[0].AbsoluteRentPence
	if ar <= 0 {
		t.Fatalf("absolute rent=%d, want > 0 (agriculture's lower composition)", ar)
	}
	for _, g := range rent.Grades {
		if g.AbsoluteRentPence != ar {
			t.Errorf("grade %s absolute rent=%d, want %d (uniform across soils)",
				g.ID, g.AbsoluteRentPence, ar)
		}
		if g.TotalRentPence != g.DifferentialRentPence+g.AbsoluteRentPence {
			t.Errorf("grade %s total=%d != differential+absolute=%d",
				g.ID, g.TotalRentPence, g.DifferentialRentPence+g.AbsoluteRentPence)
		}
	}
	// The worst soil's whole rent is absolute rent.
	if rent.Grades[0].TotalRentPence != ar {
		t.Errorf("worst soil total rent=%d, want %d (absolute only)",
			rent.Grades[0].TotalRentPence, ar)
	}
}

// TestDeriveRent_Conservation verifies the terrain totals sum their grades and
// that land price is each grade's rent capitalised at the interest rate.
func TestDeriveRent_Conservation(t *testing.T) {
	t.Parallel()
	abode := seedAbode()
	r := seedReadout(abode)

	const interestBP = 500
	rent := DeriveRent(r, 810000, 1000, interestBP)

	var sumDR, sumTotal, sumLand int64
	for _, g := range rent.Grades {
		sumDR += g.DifferentialRentPence
		sumTotal += g.TotalRentPence
		sumLand += g.LandPricePence
		wantLand := (g.TotalRentPence*10000 + interestBP/2) / interestBP
		if g.LandPricePence != wantLand {
			t.Errorf("grade %s land price=%d, want %d (rent ÷ interest)",
				g.ID, g.LandPricePence, wantLand)
		}
	}
	if rent.DifferentialPence != sumDR {
		t.Errorf("DifferentialPence=%d, want %d", rent.DifferentialPence, sumDR)
	}
	if rent.TotalRentPence != sumTotal {
		t.Errorf("TotalRentPence=%d, want %d", rent.TotalRentPence, sumTotal)
	}
	if rent.LandPricePence != sumLand {
		t.Errorf("LandPricePence=%d, want %d", rent.LandPricePence, sumLand)
	}
	if rent.AbsolutePence != rent.Grades[0].AbsoluteRentPence*int64(len(rent.Grades)) {
		t.Errorf("AbsolutePence=%d, want %d (per-soil × grades)",
			rent.AbsolutePence, rent.Grades[0].AbsoluteRentPence*int64(len(rent.Grades)))
	}
}

// TestDeriveDistribution_CompositionNonDecreasing verifies that over a short run
// (fewer periods than the crisis interval), the illustrative composition only
// climbs — CompositionBP in the series is non-decreasing before any crisis.
func TestDeriveDistribution_CompositionNonDecreasing(t *testing.T) {
	t.Parallel()
	abode := seedAbode()
	r := seedReadout(abode)

	// Four periods — short enough that no crisis fires; composition only rises.
	periods := []simulation.GeneralLawPeriod{
		{Period: 1, WagePence: 2500, RateOfExploitationBP: 10000, ReserveArmyCount: 30, OrganicCompositionBP: 20000, EmployedCount: 120},
		{Period: 2, WagePence: 2450, RateOfExploitationBP: 10100, ReserveArmyCount: 32, OrganicCompositionBP: 22000, EmployedCount: 121},
		{Period: 3, WagePence: 2400, RateOfExploitationBP: 10200, ReserveArmyCount: 35, OrganicCompositionBP: 25000, EmployedCount: 122},
		{Period: 4, WagePence: 2350, RateOfExploitationBP: 10300, ReserveArmyCount: 38, OrganicCompositionBP: 28000, EmployedCount: 123},
	}
	dist := DeriveDistribution(periods, r, 90000, 810000, 900000, 1000)

	for i := 1; i < len(dist.DistSeries); i++ {
		prev := dist.DistSeries[i-1].CompositionBP
		curr := dist.DistSeries[i].CompositionBP
		if curr < prev {
			t.Errorf("CompositionBP decreased from period %d to %d: %d -> %d",
				dist.DistSeries[i-1].Period, dist.DistSeries[i].Period, prev, curr)
		}
	}
}

// TestDeriveDistribution_CrisisSawtooth verifies the illustrative TRPF trajectory:
// over a long-enough run the rising composition drives p′ to the floor, a crisis
// fires (devaluing constant capital), and the rate is restored — the sawtooth.
// It also checks the mass of profit rises monotonically (the rate–mass contradiction).
func TestDeriveDistribution_CrisisSawtooth(t *testing.T) {
	t.Parallel()
	abode := seedAbode()
	r := seedReadout(abode)

	// 24 periods is well past the ~12-period crisis interval at the seed rate.
	periods := make([]simulation.GeneralLawPeriod, 24)
	for i := range periods {
		periods[i] = simulation.GeneralLawPeriod{
			Period:               int64(i + 1),
			WagePence:            2000,
			RateOfExploitationBP: 10000,
			OrganicCompositionBP: 20000,
			EmployedCount:        100,
		}
	}
	dist := DeriveDistribution(periods, r, 90000, 810000, 900000, 1000)

	firstCrisis := -1
	for i, p := range dist.DistSeries {
		if p.Crisis {
			firstCrisis = i
			break
		}
	}
	if firstCrisis <= 0 {
		t.Fatalf("expected a crisis after the rate falls to the floor; got firstCrisis=%d", firstCrisis)
	}
	// The crisis restores the rate: the crisis point's p′ exceeds the declining
	// point just before it.
	if dist.DistSeries[firstCrisis].PprimeBP <= dist.DistSeries[firstCrisis-1].PprimeBP {
		t.Errorf("crisis should restore p′: point %d=%d not > prior %d=%d",
			firstCrisis, dist.DistSeries[firstCrisis].PprimeBP,
			firstCrisis-1, dist.DistSeries[firstCrisis-1].PprimeBP)
	}
	// Mass of profit rises monotonically across the cycle.
	for i := 1; i < len(dist.DistSeries); i++ {
		if dist.DistSeries[i].ProfitMassPence < dist.DistSeries[i-1].ProfitMassPence {
			t.Errorf("mass of profit decreased at %d: %d -> %d",
				i, dist.DistSeries[i-1].ProfitMassPence, dist.DistSeries[i].ProfitMassPence)
		}
	}
}

// TestRateBP verifies round-half-up semantics and zero-denominator guard.
func TestRateBP(t *testing.T) {
	t.Parallel()
	cases := []struct {
		num, den int64
		want     int64
	}{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 0},
		{10000, 10000, 10000},
		{5000, 10000, 5000},
		{1, 2, 5000},     // 0.5 rounds to 5000
		{3, 4, 7500},     // 0.75
		{300000, 900000, 3333}, // Marx 2:1 c:v → p′≈33.33%
	}
	for _, tc := range cases {
		got := RateBP(tc.num, tc.den)
		if got != tc.want {
			t.Errorf("RateBP(%d,%d)=%d, want %d", tc.num, tc.den, got, tc.want)
		}
	}
}
