// Package production implements relative surplus-value concepts from
// Capital Vol. I, Ch. 12. All functions are pure; no persistence is needed.
package production

import (
	"math"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/labour"
)

// LabourMinutes and ProductivityFactor live in the labour package so every
// chapter implemented in simulation-engine speaks the same value-magnitude
// vocabulary. The aliases here keep call-sites compiling.
type LabourMinutes = labour.LabourMinutes
type ProductivityFactor = labour.ProductivityFactor

// Named LabourMinutes types for distinct domain roles (Ch. 12, §1).
type AbsoluteSurplusValue LabourMinutes
type RelativeSurplusValue LabourMinutes
type NecessaryLabour LabourMinutes
type SurplusLabour LabourMinutes
type LabourPowerValue LabourMinutes
type IndividualValue LabourMinutes
type SocialValue LabourMinutes

// Quantity is the number of articles produced or sold.
type Quantity int64

// Compile-time guards: keep the alias-vs-defined-type distinction honest.
// A bare LabourMinutes value must round-trip via the labour package.
var _ labour.LabourMinutes = LabourMinutes(0)
var _ labour.ProductivityFactor = ProductivityFactor(0)

// WorkingDay holds the tripartite split of the total working day [§1].
// Invariant: NecessaryLabour + SurplusLabour == Total.
type WorkingDay struct {
	Total           LabourMinutes   `json:"total"`
	NecessaryLabour NecessaryLabour `json:"necessary_labour"`
	SurplusLabour   SurplusLabour   `json:"surplus_labour"`
}

// ShortenNecessaryLabour recomputes the working-day split after a productivity
// rise reduces the value of labour-power to newLPV. Total day is unchanged [§1].
func ShortenNecessaryLabour(wd WorkingDay, newLPV LabourPowerValue) WorkingDay {
	nl := NecessaryLabour(newLPV)
	sl := SurplusLabour(wd.Total - LabourMinutes(nl))
	return WorkingDay{Total: wd.Total, NecessaryLabour: nl, SurplusLabour: sl}
}

// RateOfSurplusValue returns s/v as a float64. Returns 0 when nl == 0 [§1].
func RateOfSurplusValue(sl SurplusLabour, nl NecessaryLabour) float64 {
	if nl == 0 {
		return 0
	}
	return float64(sl) / float64(nl)
}

// ExtraSurplusValue is the total extra gain captured by a capitalist whose
// individual value is below the social value, selling at social value.
// Returns 0 when iv >= sv [§1].
func ExtraSurplusValue(iv IndividualValue, sv SocialValue, qty Quantity) LabourMinutes {
	if iv >= IndividualValue(sv) {
		return 0
	}
	return LabourMinutes(sv-SocialValue(iv)) * LabourMinutes(qty)
}

// ApplyProductivityToSNLT implements the inverse law: value ∝ 1/productivity.
// Returns snlt unchanged for non-positive factor [§1].
func ApplyProductivityToSNLT(snlt LabourMinutes, pf ProductivityFactor) LabourMinutes {
	if pf <= 0 {
		return snlt
	}
	return LabourMinutes(math.Round(float64(snlt) / float64(pf)))
}
