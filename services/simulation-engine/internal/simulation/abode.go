package simulation

// SocialWorkingDayMinutes is the length of the aggregate social working day
// (Vol. I Ch. 9) against which the rate of surplus-value splits paid (necessary)
// from unpaid (surplus) labour. Ten hours — the post-1850 normal working day.
const SocialWorkingDayMinutes = 600

// AbodeState is the evolving aggregate class relation beneath the field of
// capitals — "the hidden abode of production, on whose threshold there hangs the
// notice 'No admittance except on business'" (Vol. I Ch. 6 fin.). It is the
// total social capital partitioned into constant (c — dead labour) and variable
// (v — the wages bill) capital, the value of labour-power (BaseWagePence), the
// labouring population (WorkerSupply), and the parameters of the law. One
// scheduler pass advances it one period via AdvanceGeneralLaw — the General Law
// of Capitalist Accumulation in motion (Vol. I Ch. 25).
type AbodeState struct {
	Period                int64 `json:"period"`
	ConstantPence         Pence `json:"constant_pence"`          // c — dead labour
	VariablePence         Pence `json:"variable_pence"`          // v — the wages bill
	BaseWagePence         int64 `json:"base_wage_pence"`         // value of labour-power, per worker
	WorkerSupply          int64 `json:"worker_supply"`           // the labouring population
	SurplusRateBaseBP     int64 `json:"surplus_rate_base_bp"`    // s′ at full employment (10000 = 100%)
	AccumulationRateBP    int64 `json:"accumulation_rate_bp"`    // α — share of surplus capitalised
	MarginalCompositionBP int64 `json:"marginal_composition_bp"` // composition c/(c+v) of NEW capital
	DisplacementRateBP    int64 `json:"displacement_rate_bp"`    // share of v repelled to c each period
	ProductivityGrowthBP  int64 `json:"productivity_growth_bp"`  // gap to full automation closed per period
	PopulationGrowthBP    int64 `json:"population_growth_bp"`    // worker-supply growth per period
}

// NewAbodeState returns the seeded initial abode — a Marx-faithful aggregate
// (c:v = 2:1, s′ = 100%) with a small initial reserve army the law expands.
// These values are mirrored in migration 00066 so the MySQL-backed and
// memory-backed economies start identically.
func NewAbodeState() AbodeState {
	return AbodeState{
		Period:                0,
		ConstantPence:         600000, // £6,000 dead labour
		VariablePence:         300000, // £3,000 wages bill
		BaseWagePence:         2500,   // £25 value of labour-power → 120 employed
		WorkerSupply:          150,    // 30 already in the reserve army
		SurplusRateBaseBP:     10000,  // 100% rate of surplus-value
		AccumulationRateBP:    5000,   // reinvest 50% of surplus
		MarginalCompositionBP: 6667,   // new capital starts at the stock composition c/(c+v)
		DisplacementRateBP:    1800,   // machinery repels 18% of the wages bill per period
		ProductivityGrowthBP:  500,    // marginal composition closes 5% of its gap to full automation
		PopulationGrowthBP:    150,    // the labouring population grows 1.5% per period
	}
}

// GeneralLawPeriod is one recorded period of the general law — a point on the
// immiseration time-series surfaced in the abode.
type GeneralLawPeriod struct {
	Period               int64 `json:"period"`
	WagePence            int64 `json:"wage_pence"`              // paid wage (price of labour, compressed)
	RateOfExploitationBP int64 `json:"rate_of_exploitation_bp"` // s/v effective
	ReserveArmyCount     int64 `json:"reserve_army_count"`
	OrganicCompositionBP int64 `json:"organic_composition_bp"` // c/v
	EmployedCount        int64 `json:"employed_count"`
}

// AbodeReadout is the instantaneous class state derived from an AbodeState — the
// surface masks exactly this. Pence/minutes/basis-points integers throughout.
type AbodeReadout struct {
	TotalVariablePence     int64 // Σv — wages = paid labour
	TotalSurplusPence      int64 // Σs — unpaid labour
	RateOfExploitationBP   int64 // s/v
	NecessaryLabourMinutes int64
	SurplusLabourMinutes   int64
	OrganicCompositionBP   int64 // c/v
	ReserveArmyCount       int64
	ReserveArmyPressureBP  int64
	EmployedCount          int64
	WagePence              int64 // paid wage after reserve-army compression
}

// reservePressureBP returns reserve / workforce in basis points, clamped to
// [0, 10000]. Zero reserve or workforce is no pressure.
func reservePressureBP(reserve, workforce int64) int64 {
	if reserve <= 0 || workforce <= 0 {
		return 0
	}
	bp := reserve * 10000 / workforce
	if bp > 10000 {
		return 10000
	}
	return bp
}

// compressWage drives the price of labour below the value of labour-power as the
// reserve army grows (Vol. I Ch. 25 §3): nominal × (1 − pressure), rounded half
// up, floored at half the value of labour-power (subsistence).
func compressWage(baseWage, pressureBP int64) int64 {
	if baseWage <= 0 || pressureBP <= 0 {
		return baseWage
	}
	w := (baseWage*(10000-pressureBP) + 5000) / 10000
	if floor := baseWage / 2; w < floor {
		return floor
	}
	return w
}

// Readout projects the current AbodeState to the instantaneous class state.
func (a AbodeState) Readout() AbodeReadout {
	oc := ComputeOrganicComposition(CapitalStock{ConstantCapital: a.ConstantPence, VariableCapital: a.VariablePence})
	employed := ComputeLabourDemand(a.ConstantPence+a.VariablePence, oc, a.BaseWagePence)
	reserve := ComputeReserveArmy(a.WorkerSupply, employed)
	pressure := reservePressureBP(reserve.Size, a.WorkerSupply)
	effRate := a.SurplusRateBaseBP * (10000 + pressure) / 10000
	surplus := int64(a.VariablePence) * effRate / 10000
	var ocBP int64
	if a.VariablePence > 0 {
		ocBP = 10000 * int64(a.ConstantPence) / int64(a.VariablePence)
	}
	necessary := int64(SocialWorkingDayMinutes) * 10000 / (10000 + effRate)
	return AbodeReadout{
		TotalVariablePence:     int64(a.VariablePence),
		TotalSurplusPence:      surplus,
		RateOfExploitationBP:   effRate,
		NecessaryLabourMinutes: necessary,
		SurplusLabourMinutes:   int64(SocialWorkingDayMinutes) - necessary,
		OrganicCompositionBP:   ocBP,
		ReserveArmyCount:       reserve.Size,
		ReserveArmyPressureBP:  pressure,
		EmployedCount:          employed.Workers,
		WagePence:              compressWage(a.BaseWagePence, pressure),
	}
}

// AdvanceGeneralLaw runs one period of the general law, returning the next state
// and the GeneralLawPeriod recorded for the immiseration series. The loop
// (Vol. I Ch. 25): the heightened surplus is re-accumulated at the marginal
// (machine-heavier) composition; machinery then repels a share of the wages bill
// into constant capital, holding v down and swelling the reserve army; the
// labouring population grows and the marginal composition climbs toward — but
// never reaches — full automation. The circle closes.
func AdvanceGeneralLaw(s AbodeState) (AbodeState, GeneralLawPeriod) {
	r := s.Readout()
	period := GeneralLawPeriod{
		Period:               s.Period + 1,
		WagePence:            r.WagePence,
		RateOfExploitationBP: r.RateOfExploitationBP,
		ReserveArmyCount:     r.ReserveArmyCount,
		OrganicCompositionBP: r.OrganicCompositionBP,
		EmployedCount:        r.EmployedCount,
	}

	next := s
	next.Period = s.Period + 1

	// Re-accumulate α·s at the marginal composition (new capital is more
	// machine-heavy than the existing stock).
	capitalised := r.TotalSurplusPence * s.AccumulationRateBP / 10000
	dc := capitalised * s.MarginalCompositionBP / 10000
	dv := capitalised - dc
	next.ConstantPence = s.ConstantPence + Pence(dc)
	next.VariablePence = s.VariablePence + Pence(dv)

	// Machinery repels living labour: a share of the wages bill is converted to
	// constant capital — the working population "produces the means by which it
	// is made relatively superfluous" (§25.3). This holds v down so the reserve
	// army grows even as total capital accumulates.
	displaced := int64(next.VariablePence) * s.DisplacementRateBP / 10000
	next.VariablePence -= Pence(displaced)
	next.ConstantPence += Pence(displaced)

	// The labouring population grows; the marginal composition asymptotes below
	// full automation (no fictional ceiling — it closes a fixed fraction of the
	// remaining gap each period).
	next.WorkerSupply = s.WorkerSupply + s.WorkerSupply*s.PopulationGrowthBP/10000
	next.MarginalCompositionBP = s.MarginalCompositionBP +
		s.ProductivityGrowthBP*(10000-s.MarginalCompositionBP)/10000

	return next, period
}

// LeverUpdate is a partial perturbation of the live abode's law parameters
// (Slice 3 — the levers). Non-nil fields are applied; nil fields are left
// unchanged. The three levers are the working day (the rate of surplus-value),
// the wage (the value of labour-power), and the accumulation rate α.
type LeverUpdate struct {
	SurplusRateBaseBP  *int64 `json:"surplus_rate_base_bp,omitempty"` // working day: necessary↔surplus
	BaseWagePence      *int64 `json:"base_wage_pence,omitempty"`      // value of labour-power
	AccumulationRateBP *int64 `json:"accumulation_rate_bp,omitempty"` // α
}

// IsEmpty reports whether the update would change nothing.
func (u LeverUpdate) IsEmpty() bool {
	return u.SurplusRateBaseBP == nil && u.BaseWagePence == nil && u.AccumulationRateBP == nil
}

// ApplyLevers returns a copy of the state with the supplied levers applied, each
// clamped so the law cannot be driven to a degenerate state: the base rate of
// surplus-value to [0, 100000] (0–1000%), the wage to at least 1 pence (a
// positive value of labour-power, so employment v/wage stays finite), and α to
// [0, 10000] basis points.
func (a AbodeState) ApplyLevers(u LeverUpdate) AbodeState {
	next := a
	if u.SurplusRateBaseBP != nil {
		v := *u.SurplusRateBaseBP
		if v < 0 {
			v = 0
		}
		if v > 100000 {
			v = 100000
		}
		next.SurplusRateBaseBP = v
	}
	if u.BaseWagePence != nil {
		v := *u.BaseWagePence
		if v < 1 {
			v = 1
		}
		next.BaseWagePence = v
	}
	if u.AccumulationRateBP != nil {
		v := *u.AccumulationRateBP
		if v < 0 {
			v = 0
		}
		if v > 10000 {
			v = 10000
		}
		next.AccumulationRateBP = v
	}
	return next
}
