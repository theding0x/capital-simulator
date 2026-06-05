package simulation

// SocialWorkingDayMinutes is the length of the aggregate social working day
// (Vol. I Ch. 9) against which the rate of surplus-value splits paid (necessary)
// from unpaid (surplus) labour. Ten hours — the post-1850 normal working day.
const SocialWorkingDayMinutes = 600

// maxMarginalCompositionBP caps the composition of newly-accumulated capital
// below full automation, so new capital always still hires some labour and the
// variable part of capital keeps growing in absolute magnitude (it just grows in
// an ever-diminishing proportion). Prevents the degenerate v→0 collapse.
const maxMarginalCompositionBP = 9500

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
		MarginalCompositionBP: 8000,   // new capital more machine-heavy than the 2:1 stock → composition rises from period 0
		DisplacementRateBP:    1800,   // legacy; no longer drives the law (kept for the persisted column)
		ProductivityGrowthBP:  500,    // marginal composition closes 5% of its gap to full automation per period
		PopulationGrowthBP:    150,    // legacy; supply now grows with the scale of capital (kept for the persisted column)
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
// [0, 9000]. The active labour army never fully vanishes, so pressure never pegs
// at 100% — the ceiling keeps the wage above absolute zero and s′ bounded. Zero
// reserve or workforce is no pressure.
func reservePressureBP(reserve, workforce int64) int64 {
	if reserve <= 0 || workforce <= 0 {
		return 0
	}
	bp := reserve * 10000 / workforce
	if bp > 9000 {
		return 9000
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
// (Vol. I Ch. 25): the heightened surplus is re-accumulated as new capital at
// the marginal (machine-heavier) composition, so the organic composition rises
// while c AND v both grow — variable capital is never destroyed. The working
// population grows with the scale of capital, but the rising composition makes
// the demand for labour grow more slowly, so a relatively redundant population —
// the reserve army — swells, pressing the wage below the value of labour-power
// and raising s′; the heightened surplus re-accumulates. The circle closes
// without driving living labour to zero.
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

	// Re-accumulate α·s as NEW capital at the marginal (machine-heavier)
	// composition. Both c and v grow in absolute magnitude; because new capital
	// is more constant-heavy than the existing stock, the organic composition
	// rises while variable capital (living labour) keeps growing — never
	// destroyed (Vol. I Ch. 25 §2: "the variable part... increases, but in a
	// constantly diminishing proportion").
	capitalised := r.TotalSurplusPence * s.AccumulationRateBP / 10000
	dc := capitalised * s.MarginalCompositionBP / 10000
	dv := capitalised - dc
	next.ConstantPence = s.ConstantPence + Pence(dc)
	next.VariablePence = s.VariablePence + Pence(dv)

	// The labouring population grows with the scale of accumulating capital
	// (capital draws population into its orbit). Because the rising composition
	// makes the demand for labour (v) grow more slowly than total capital, supply
	// outruns demand: a relatively redundant population — the industrial reserve
	// army — grows and its pressure rises, while the active army still grows in
	// absolute magnitude (§25.3).
	oldTotal := int64(s.ConstantPence + s.VariablePence)
	newTotal := int64(next.ConstantPence + next.VariablePence)
	if oldTotal > 0 {
		growthBP := (newTotal - oldTotal) * 10000 / oldTotal
		next.WorkerSupply = s.WorkerSupply + s.WorkerSupply*growthBP/10000
	}

	// The marginal composition climbs toward — but never reaches — full
	// automation, capped so new capital always still hires some labour.
	next.MarginalCompositionBP = s.MarginalCompositionBP +
		s.ProductivityGrowthBP*(maxMarginalCompositionBP-s.MarginalCompositionBP)/10000
	if next.MarginalCompositionBP > maxMarginalCompositionBP {
		next.MarginalCompositionBP = maxMarginalCompositionBP
	}

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
