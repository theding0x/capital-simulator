package surplus

import (
	"math"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/labour"
)

// Capital Vol. I, Ch. 16 — Absolute and Relative Surplus-Value.
//
// The synthesis of Chs. 11–15: surplus-value comes in two analytically
// distinct forms which nonetheless share a single magnitude type. Marx's
// closing argument is that the distinction "appears illusory" from one
// angle — they share the same arithmetic — yet remains real, because
// each presupposes a different mechanism (prolongation vs. productivity)
// and a different historical phase of capitalist development (formal vs.
// real subjection of labour).

// Origin tags a SurplusValue with the mechanism that produced it.
type Origin string

const (
	OriginAbsolute Origin = "absolute"
	OriginRelative Origin = "relative"
)

// SurplusValue is labour-time appropriated beyond the equivalent of
// labour-power. It wraps a LabourMinutes magnitude with an Origin tag so
// callers can distinguish absolute from relative gains without inspecting
// the call-site that produced them.
type SurplusValue struct {
	LabourMinutes labour.LabourMinutes `json:"labour_minutes"`
	Origin        Origin               `json:"origin"`
}

// AbsoluteSurplusValue is surplus produced by prolonging the working day
// past necessary labour-time (§1, "production of absolute surplus-value").
type AbsoluteSurplusValue struct {
	SurplusValue
}

// RelativeSurplusValue is surplus produced by shortening necessary labour
// via a productivity rise (§1, "in order to prolong the surplus-labour,
// the necessary labour is shortened").
type RelativeSurplusValue struct {
	SurplusValue
}

// WorkingDay is the §1 partition: total time the labourer is at the
// capitalist's disposal, split into the necessary-labour share (which
// reproduces labour-power) and the surplus-labour share (appropriated by
// capital). All fields are LabourMinutes; the invariant
//
//	NecessaryLabour + SurplusLabour == Total
//
// must hold after every transition.
type WorkingDay struct {
	Total           labour.LabourMinutes `json:"total"`
	NecessaryLabour labour.LabourMinutes `json:"necessary_labour"`
	SurplusLabour   labour.LabourMinutes `json:"surplus_labour"`
}

// ProductivityFactor is a multiplicative scalar on output per unit time.
// >1 reduces the SNLT of subsistence goods and therefore necessary labour.
// Re-export of the canonical labour package type so Ch. 16 call-sites
// can use a single import.
type ProductivityFactor = labour.ProductivityFactor

// ProlongWorkingDay extends Total by extraMinutes, holding NecessaryLabour
// fixed. The increment lands entirely in SurplusLabour, producing an
// AbsoluteSurplusValue equal to extraMinutes (§1).
func ProlongWorkingDay(wd WorkingDay, extraMinutes labour.LabourMinutes) (WorkingDay, AbsoluteSurplusValue) {
	if extraMinutes <= 0 {
		return wd, AbsoluteSurplusValue{SurplusValue: SurplusValue{Origin: OriginAbsolute}}
	}
	newWD := WorkingDay{
		Total:           wd.Total + extraMinutes,
		NecessaryLabour: wd.NecessaryLabour,
		SurplusLabour:   wd.SurplusLabour + extraMinutes,
	}
	return newWD, AbsoluteSurplusValue{
		SurplusValue: SurplusValue{LabourMinutes: extraMinutes, Origin: OriginAbsolute},
	}
}

// ReduceNecessaryLabour applies a ProductivityFactor to NecessaryLabour,
// holding Total fixed. The freed minutes shift into SurplusLabour, producing
// a RelativeSurplusValue equal to the necessary-labour reduction (§1).
// Factor ≤ 0 is rejected as a no-op; factor ≥ enough-to-zero is floored at
// NecessaryLabour = 1 to preserve the wage relation.
func ReduceNecessaryLabour(wd WorkingDay, factor ProductivityFactor) (WorkingDay, RelativeSurplusValue) {
	if factor <= 0 {
		return wd, RelativeSurplusValue{SurplusValue: SurplusValue{Origin: OriginRelative}}
	}
	scaled := labour.LabourMinutes(math.Round(float64(wd.NecessaryLabour) / float64(factor)))
	if scaled < 1 {
		scaled = 1
	}
	if scaled >= wd.NecessaryLabour {
		return wd, RelativeSurplusValue{SurplusValue: SurplusValue{Origin: OriginRelative}}
	}
	diff := wd.NecessaryLabour - scaled
	newWD := WorkingDay{
		Total:           wd.Total,
		NecessaryLabour: scaled,
		SurplusLabour:   wd.SurplusLabour + diff,
	}
	return newWD, RelativeSurplusValue{
		SurplusValue: SurplusValue{LabourMinutes: diff, Origin: OriginRelative},
	}
}

// RateSurplusValue is Marx's "rate of exploitation": surplus-labour divided
// by necessary-labour. Returns 0 when necessary is 0 (undefined wage relation).
func RateSurplusValue(surplusLabour, necessaryLabour labour.LabourMinutes) float64 {
	if necessaryLabour == 0 {
		return 0
	}
	return float64(surplusLabour) / float64(necessaryLabour)
}

// RateOfProfit is surplus-value divided by the *total* capital advanced
// (constant + variable). Always strictly less than RateSurplusValue when
// constant capital is positive — Marx's correction of Mill (§1):
// "if the rate of surplus-value be 20%, the rate of profit will be
// 20:500, i.e., 4% and not 20%".
func RateOfProfit(surplusValue, totalCapitalAdvanced labour.LabourMinutes) float64 {
	if totalCapitalAdvanced == 0 {
		return 0
	}
	return float64(surplusValue) / float64(totalCapitalAdvanced)
}

