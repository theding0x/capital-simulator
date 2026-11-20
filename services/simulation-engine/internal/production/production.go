// Package production implements relative surplus-value concepts from
// Capital Vol. I, Ch. 12. All functions are pure; no persistence is needed.
package production

import (
	"errors"
	"math"

	"github.com/theding0x/capital-simulator/pkg/surplusrate"
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

// Partitioned reports whether the day satisfies the L-P1 invariant
// NecessaryLabour + SurplusLabour == Total. It is the one definition of that
// check; the surplus package and the HTTP layers use this type too.
func (wd WorkingDay) Partitioned() bool {
	return LabourMinutes(wd.NecessaryLabour)+LabourMinutes(wd.SurplusLabour) == wd.Total
}

// Sentinel errors for ShortenNecessaryLabour (the L-P1 partition guards).
var (
	ErrWorkingDayNotPartitioned = errors.New("production: necessary + surplus labour must equal the working day")
	ErrLPVNotPositive           = errors.New("production: new labour-power value must be positive")
	ErrLPVExceedsWorkingDay     = errors.New("production: new labour-power value exceeds the working day")
	ErrNotShortening            = errors.New("production: new labour-power value must be below the current necessary labour")
)

// ShortenNecessaryLabour recomputes the working-day split after a productivity
// rise reduces the value of labour-power to newLPV. Total day is unchanged [§1].
//
// Relative surplus-value arises only from a genuine shortening, so it rejects
// an unpartitioned day, a non-positive newLPV, a newLPV beyond the day (which
// would make surplus labour negative), and a newLPV that does not fall below
// the current necessary labour. SurplusLabour >= 0 and NL + SL == Total hold on
// every successful return.
func ShortenNecessaryLabour(wd WorkingDay, newLPV LabourPowerValue) (WorkingDay, error) {
	switch {
	case !wd.Partitioned():
		return WorkingDay{}, ErrWorkingDayNotPartitioned
	case newLPV <= 0:
		return WorkingDay{}, ErrLPVNotPositive
	case LabourMinutes(newLPV) > wd.Total:
		return WorkingDay{}, ErrLPVExceedsWorkingDay
	case NecessaryLabour(newLPV) >= wd.NecessaryLabour:
		return WorkingDay{}, ErrNotShortening
	}
	nl := NecessaryLabour(newLPV)
	sl := SurplusLabour(wd.Total - LabourMinutes(nl))
	return WorkingDay{Total: wd.Total, NecessaryLabour: nl, SurplusLabour: sl}, nil
}

// RateOfSurplusValue returns s/v as a float64. Returns 0 when nl == 0 [§1].
func RateOfSurplusValue(sl SurplusLabour, nl NecessaryLabour) float64 {
	r, _ := surplusrate.Ratio(int64(sl), int64(nl))
	return r
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
