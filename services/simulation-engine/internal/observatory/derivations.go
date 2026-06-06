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

// DistributionSnapshot is the Vol. III distribution view — the general rate of
// profit, interest, fictitious capital, the trinity formula, and the p′ series.
type DistributionSnapshot struct {
	GeneralRateBP   int64             `json:"general_rate_bp"`
	InterestRateBP  int64             `json:"interest_rate_bp"`
	LoanablePence   int64             `json:"loanable_pence"`
	FictitiousPence int64             `json:"fictitious_pence"`
	DistSeries      []DistSeriesPoint `json:"dist_series"`
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

	var dist DistributionSnapshot
	dist.GeneralRateBP = generalRateBP
	dist.InterestRateBP = interestBP
	dist.LoanablePence = loanable
	dist.FictitiousPence = fictitious
	dist.DistSeries = series
	dist.Trinity.WagesPence = r.TotalVariablePence
	dist.Trinity.ProfitPence = r.TotalSurplusPence
	dist.Trinity.RentPence = 0
	return dist
}
