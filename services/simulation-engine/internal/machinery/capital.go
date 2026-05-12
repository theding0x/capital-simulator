package machinery

// VariableCapital is the portion of capital advanced in wages (labour-power),
// expressed in pence to match the agent-service money unit.
type VariableCapital int64

// ConstantCapital is the portion advanced in means of production (instruments,
// raw materials, auxiliaries), also in pence.
type ConstantCapital int64

// CapitalComposition pairs the variable and constant portions of total
// capital. §6: mechanisation converts variable into constant; total capital
// is conserved across the substitution period, only the composition shifts.
type CapitalComposition struct {
	Variable VariableCapital `json:"variable"`
	Constant ConstantCapital `json:"constant"`
}

// Total returns variable + constant (in pence).
func (c CapitalComposition) Total() int64 {
	return int64(c.Variable) + int64(c.Constant)
}

// Ratio returns the organic composition c/v. Zero variable returns 0.
func (c CapitalComposition) Ratio() float64 {
	if c.Variable == 0 {
		return 0
	}
	return float64(c.Constant) / float64(c.Variable)
}

// Substitute is the §6 mechanisation move: variable capital is reduced by
// `wageBillRemoved`, constant capital rises by `machineCost`. In the pure
// substitution period the two are equal and total capital is unchanged.
func (c CapitalComposition) Substitute(wageBillRemoved VariableCapital, machineCost ConstantCapital) CapitalComposition {
	return CapitalComposition{
		Variable: c.Variable - wageBillRemoved,
		Constant: c.Constant + machineCost,
	}
}

// IntensityFactor is the §3C labour-density multiplier. When the working day
// is shortened, each remaining hour is filled more densely with labour; an
// IntensityFactor > 1 captures that compression. The factor is bounded
// below by 1.0 — labour cannot fall below its nominal density without
// changing the meaning of the time-measure itself.
type IntensityFactor float64

// NewIntensityFactor clamps the multiplier to ≥ 1.0.
func NewIntensityFactor(v float64) IntensityFactor {
	if v < 1.0 {
		return 1.0
	}
	return IntensityFactor(v)
}

// EffectiveLabour scales nominal labour-minutes by the intensity factor. A 10-
// hour day at 1.5× intensity counts as 15 hours of average-intensity labour.
func EffectiveLabour(nominal LabourMinutes, factor IntensityFactor) LabourMinutes {
	if factor < 1.0 {
		factor = 1.0
	}
	return LabourMinutes(float64(nominal) * float64(factor))
}
