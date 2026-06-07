package observatory

import "github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"

// RateBP returns round-half-up(10000 * num / den) in basis points; 0 when den <= 0.
func RateBP(num, den int64) int64 {
	if den <= 0 {
		return 0
	}
	return (num*10000 + den/2) / den
}

// roundHalfUp64 returns round-half-up(num / den); 0 when den <= 0.
func roundHalfUp64(num, den int64) int64 {
	if den <= 0 {
		return 0
	}
	return (num + den/2) / den
}

// DeptCVS is the c/v/s decomposition of one department's total capital (pence).
type DeptCVS struct {
	C int64 `json:"c"`
	V int64 `json:"v"`
	S int64 `json:"s"`
}

// CirculationSnapshot is the Vol. II reproduction-scheme view derived from the
// current abode and field aggregate. Pure — no I/O, no mutations.
type CirculationSnapshot struct {
	Departments struct {
		I  DeptCVS `json:"I"`
		II DeptCVS `json:"II"`
	} `json:"departments"`
	KeystoneExchangePence int64 `json:"keystone_exchange_pence"`
	DeptIIConstantPence   int64 `json:"dept_ii_constant_pence"`
	AccumulationRateBP    int64 `json:"accumulation_rate_bp"`
	Extended              bool  `json:"extended"`
}

// DistSeriesPoint is one period entry in the rate-of-profit/crisis time series.
type DistSeriesPoint struct {
	Period          int64 `json:"period"`
	PprimeBP        int64 `json:"pprime_bp"`
	ProfitMassPence int64 `json:"profit_mass_pence"`
	CompositionBP   int64 `json:"composition_bp"`
	Crisis          bool  `json:"crisis"`
}

// SoilGrade is one band of the rent terrain — a soil of a given fertility worked
// with equal capital. The worst soil in cultivation regulates the market price
// and yields no differential rent (Vol. III Ch.39); a more fertile soil clears
// more output at that same price, and the surplus profit is differential rent.
type SoilGrade struct {
	ID                    string `json:"id"`                     // "A".."D", worst → best
	Name                  string `json:"name"`                   // "worst soil" …
	FertilityBP           int64  `json:"fertility_bp"`           // relative yield; worst soil = 10000
	YieldUnits            int64  `json:"yield_units"`            // output at equal capital
	IndividualPricePence  int64  `json:"individual_price_pence"` // the soil's own price of production per unit
	DifferentialRentPence int64  `json:"differential_rent_pence"`
	AbsoluteRentPence     int64  `json:"absolute_rent_pence"`
	TotalRentPence        int64  `json:"total_rent_pence"`
	LandPricePence        int64  `json:"land_price_pence"` // capitalised rent = rent ÷ interest
	Regulating            bool   `json:"regulating"`       // the worst soil that sets the price
}

// RentSnapshot is the Vol. III ground-rent terrain (Ch.37–47): differential rent
// across soil grades above a waterline set by the worst soil, plus the absolute
// rent landed property exacts on every soil out of agriculture's lower organic
// composition. A faithful derived illustration — the live run does not simulate
// agriculture, so the terrain is projected from the social capital and the
// general rate of profit, in the manner of the reproduction scheme and the TRPF
// series.
type RentSnapshot struct {
	Grades               []SoilGrade `json:"grades"`
	WaterlinePence       int64       `json:"waterline_pence"`         // regulating worst-soil price per unit
	CapitalPerGradePence int64       `json:"capital_per_grade_pence"` // equal capital advanced per soil
	DifferentialPence    int64       `json:"differential_pence"`      // Σ differential rent
	AbsolutePence        int64       `json:"absolute_pence"`          // Σ absolute rent (all soils)
	TotalRentPence       int64       `json:"total_rent_pence"`
	LandPricePence       int64       `json:"land_price_pence"` // Σ capitalised rent
	InterestRateBP       int64       `json:"interest_rate_bp"` // the capitalisation rate
}

// DistributionSnapshot is the Vol. III distribution view — the general rate of
// profit, interest, fictitious capital, ground-rent, the trinity formula, and
// the p′ series.
type DistributionSnapshot struct {
	GeneralRateBP   int64             `json:"general_rate_bp"`
	InterestRateBP  int64             `json:"interest_rate_bp"`
	LoanablePence   int64             `json:"loanable_pence"`
	FictitiousPence int64             `json:"fictitious_pence"`
	DistSeries      []DistSeriesPoint `json:"dist_series"`
	Rent            RentSnapshot      `json:"rent"`
	Trinity         struct {
		WagesPence  int64 `json:"wages_pence"`
		ProfitPence int64 `json:"profit_pence"`
		RentPence   int64 `json:"rent_pence"`
	} `json:"trinity"`
}

// DeriveCirculation projects the Vol. II simple/extended reproduction scheme
// from the abode parameters and the field aggregate totals (sumSurplus and
// sumTotal, both in pence). Pure — never mutates its inputs.
func DeriveCirculation(abode simulation.AbodeState, r simulation.AbodeReadout, sumSurplus, sumTotal int64) CirculationSnapshot {
	const (
		seedOCbp = 20000
		maxOCbp  = 300000
	)

	ocBP := r.OrganicCompositionBP

	var compShiftBP int64
	if ocBP > seedOCbp {
		compShiftBP = (ocBP - seedOCbp) * 10000 / (maxOCbp - seedOCbp)
		if compShiftBP > 10000 {
			compShiftBP = 10000
		}
	}

	fracIbp := int64(5500) + compShiftBP*12/100
	totalI := sumTotal * fracIbp / 10000
	totalII := sumTotal - totalI

	cShareBP := int64(6200) + compShiftBP*18/100

	splitDept := func(total int64) DeptCVS {
		c := roundHalfUp64(total*cShareBP, 10000)
		rest := total - c
		effRate := abode.SurplusRateBaseBP
		v := roundHalfUp64(rest*10000, 10000+effRate)
		s := rest - v
		return DeptCVS{C: c, V: v, S: s}
	}

	dI := splitDept(totalI)
	dII := splitDept(totalII)

	var snap CirculationSnapshot
	snap.Departments.I = dI
	snap.Departments.II = dII
	snap.KeystoneExchangePence = dI.V + dI.S
	snap.DeptIIConstantPence = dII.C
	snap.AccumulationRateBP = abode.AccumulationRateBP
	snap.Extended = abode.AccumulationRateBP > 0
	return snap
}

// DeriveDistribution projects the Vol. III distribution view from the period
// history, current readout, and field aggregate sums. Pure — never mutates
// its inputs.
func DeriveDistribution(periods []simulation.GeneralLawPeriod, r simulation.AbodeReadout, sumSurplus, sumCost, sumTotal, generalRateBP int64) DistributionSnapshot {
	interestBP := generalRateBP * 42 / 100
	if interestBP < 120 {
		interestBP = 120
	}
	if interestBP > 1600 {
		interestBP = 1600
	}

	loanable := sumSurplus * 18 / 10

	var cs int64
	if r.OrganicCompositionBP > 20000 {
		cs = (r.OrganicCompositionBP - 20000) * 100 / (300000 - 20000)
		if cs > 100 {
			cs = 100
		}
	}
	// Reordered (divide before multiply) so the intermediate cannot overflow int64.
	fictitious := sumTotal / 1000 * (400 + cs*8)

	// dist_series — the tendential fall of the rate of profit in motion (Vol III
	// Ch.13–15) as a faithful illustration, keyed on the absolute period so the
	// trajectory is stable across the rolling window the run keeps. With a constant
	// rate of exploitation a rising organic composition drives p′ down across each
	// ~12-period cycle while the mass of profit climbs (the rate–mass contradiction);
	// at the cycle's close a crisis devalues constant capital and restores the rate —
	// the sawtooth.
	const (
		trpfStartCompBP = 8000 // opening c/v ≈ 0.8 : 1 — a high opening rate of profit
		trpfStepCompBP  = 5500 // composition climbs each period within a cycle
		trpfCycleLen    = 12   // periods per accumulation cycle before a crisis
	)
	sRate := r.RateOfExploitationBP // constant rate of exploitation s/v
	if sRate <= 0 {
		sRate = 10000
	}
	massBase := sumSurplus
	if massBase <= 0 {
		massBase = r.TotalSurplusPence
	}
	if massBase <= 0 {
		massBase = 1
	}

	series := make([]DistSeriesPoint, len(periods))
	for i, p := range periods {
		phase := (p.Period - 1) % trpfCycleLen
		if phase < 0 {
			phase = 0
		}
		compBP := int64(trpfStartCompBP) + phase*trpfStepCompBP
		// At each cycle boundary constant capital is devalued and the rate restored.
		crisis := phase == 0 && p.Period > 1
		series[i] = DistSeriesPoint{
			Period:          p.Period,
			PprimeBP:        RateBP(sRate, 10000+compBP),
			ProfitMassPence: massBase * (100 + (p.Period-1)*4) / 100,
			CompositionBP:   compBP,
			Crisis:          crisis,
		}
	}

	rent := DeriveRent(r, sumCost, generalRateBP, interestBP)

	var dist DistributionSnapshot
	dist.GeneralRateBP = generalRateBP
	dist.InterestRateBP = interestBP
	dist.LoanablePence = loanable
	dist.FictitiousPence = fictitious
	dist.DistSeries = series
	dist.Rent = rent
	dist.Trinity.WagesPence = r.TotalVariablePence
	dist.Trinity.ProfitPence = r.TotalSurplusPence
	// Ground-rent is now modelled as a derived terrain (Ch.37–47), so the
	// trinity's third revenue is no longer a hard zero.
	dist.Trinity.RentPence = rent.TotalRentPence
	return dist
}

// DeriveRent projects the Vol. III ground-rent terrain from the social capital
// and the general rate of profit. Four soils of rising fertility (Marx's
// canonical 1:2:3:4 yields, Ch.39) take equal capital; the worst regulates the
// market price — the waterline. A more fertile soil clears more output at that
// price, and the surplus profit is differential rent. Agriculture's lower
// organic composition makes its value exceed its price of production; that
// excess, withheld by the monopoly of landed property, is absolute rent, exacted
// on every soil including the worst. Land price is the rent capitalised at the
// rate of interest. Pure — never mutates its inputs.
func DeriveRent(r simulation.AbodeReadout, sumCost, generalRateBP, interestBP int64) RentSnapshot {
	const (
		nGrades   = 4
		floorK    = 250000 // £2,500 — keep the terrain legible on an empty run
		agCompNum = 60     // agriculture's organic composition ≈ 60% of the social average
	)

	capRate := interestBP
	if capRate <= 0 {
		capRate = 500 // 5% — a sane capitalisation floor
	}

	k := sumCost / 12
	if k < floorK {
		k = floorK
	}
	avgProfit := roundHalfUp64(k*generalRateBP, 10000)
	priceProd := k + avgProfit // price of production per farm — equal across soils

	// Absolute rent — agriculture's value above its price of production. Its lower
	// composition leaves more variable capital, hence more surplus, than the
	// average; landed property withholds that excess on every soil.
	sRate := r.RateOfExploitationBP
	if sRate <= 0 {
		sRate = 10000
	}
	agOCbp := r.OrganicCompositionBP * agCompNum / 100
	vAg := roundHalfUp64(k*10000, 10000+agOCbp)
	sAg := roundHalfUp64(vAg*sRate, 10000)
	absolute := sAg - avgProfit // value − price of production
	if absolute < 0 {
		absolute = 0
	}

	// The regulating market price — the waterline — is the worst soil's price of
	// production plus the absolute rent it must yield to landed property: every
	// soil sells here; the worst pays it all out as absolute rent, the more
	// fertile keep the surplus above it as differential rent. (So the price-of-
	// production band + the absolute-rent band reach exactly the waterline.)
	waterline := priceProd + absolute

	names := [nGrades]string{"worst soil", "poorer soil", "richer soil", "best soil"}
	ids := [nGrades]string{"A", "B", "C", "D"}

	grades := make([]SoilGrade, nGrades)
	var sumDR, sumLand, sumTotal int64
	for i := 0; i < nGrades; i++ {
		yield := int64(i + 1) // 1, 2, 3, 4
		differential := int64(i) * waterline // (yield − worst yield) × the market price
		total := differential + absolute
		land := roundHalfUp64(total*10000, capRate)
		grades[i] = SoilGrade{
			ID:                    ids[i],
			Name:                  names[i],
			FertilityBP:           yield * 10000,
			YieldUnits:            yield,
			IndividualPricePence:  roundHalfUp64(priceProd, yield),
			DifferentialRentPence: differential,
			AbsoluteRentPence:     absolute,
			TotalRentPence:        total,
			LandPricePence:        land,
			Regulating:            i == 0,
		}
		sumDR += differential
		sumLand += land
		sumTotal += total
	}

	return RentSnapshot{
		Grades:               grades,
		WaterlinePence:       waterline,
		CapitalPerGradePence: k,
		DifferentialPence:    sumDR,
		AbsolutePence:        absolute * nGrades,
		TotalRentPence:       sumTotal,
		LandPricePence:       sumLand,
		InterestRateBP:       capRate,
	}
}
