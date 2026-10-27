// Package surplus implements the rate and mass of surplus-value from
// Capital Vol. I, Ch. 11. All functions are pure; no persistence is needed.
package surplus

import (
	"errors"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/labour"
)

// AbsoluteWorkdayLimit is the physical ceiling of any working day: 24 h × 60 min.
const AbsoluteWorkdayLimit labour.LabourMinutes = 24 * 60

// ErrExceedsWorkdayLimit is returned when necessary + surplus labour reaches or
// exceeds AbsoluteWorkdayLimit. Marx, Ch. 11 second law: the absolute limit of
// the working day "sets an absolute limit to the compensation of a reduction of
// variable capital by a higher rate of surplus-value".
var ErrExceedsWorkdayLimit = errors.New("surplus: necessary + surplus labour must be less than the 24-hour absolute workday limit")

// SurplusValueMass is the total surplus-value extracted from all simultaneously
// employed workers. Expressed in the same abstract unit as VariableCapital (§1).
type SurplusValueMass int64

// SurplusValueRate is the ratio of surplus-labour to necessary-labour, both in
// LabourMinutes. Corresponds to Marx's s/v or a′/a notation [§1].
type SurplusValueRate struct {
	SurplusLabour   int64 `json:"surplus_labour"`
	NecessaryLabour int64 `json:"necessary_labour"`
}

// Validate enforces necessary + surplus < AbsoluteWorkdayLimit [§1, second law].
func (r SurplusValueRate) Validate() error {
	if labour.LabourMinutes(r.NecessaryLabour+r.SurplusLabour) >= AbsoluteWorkdayLimit {
		return ErrExceedsWorkdayLimit
	}
	return nil
}

// Rate returns the dimensionless ratio s/v as a float64.
func (r SurplusValueRate) Rate() float64 {
	if r.NecessaryLabour == 0 {
		return 0
	}
	return float64(r.SurplusLabour) / float64(r.NecessaryLabour)
}

// VariableCapital is the total money-value advanced for labour-power across all
// simultaneously employed workers [§1]. Expressed in the same integer unit as
// SurplusValueMass so that S = rate × V holds exactly.
type VariableCapital int64

// LabourPowerValue is the daily cost of reproducing a single worker [§1].
type LabourPowerValue int64

// WorkerCount is the number of simultaneously employed workers [§1].
type WorkerCount int

// IndividualSurplus returns the surplus-value produced by a single worker
// given the rate and the value of their labour-power [§1].
// s_individual = v × (a′/a) = v × SurplusLabour / NecessaryLabour
func IndividualSurplus(rate SurplusValueRate, v LabourPowerValue) SurplusValueMass {
	if rate.NecessaryLabour == 0 {
		return 0
	}
	return SurplusValueMass(int64(v) * rate.SurplusLabour / rate.NecessaryLabour)
}

// MassByRate computes S = (s/v) × V — total surplus-value from total variable
// capital and the rate of surplus-value [§1, formula I].
func MassByRate(rate SurplusValueRate, totalVariableCapital VariableCapital) SurplusValueMass {
	if rate.NecessaryLabour == 0 {
		return 0
	}
	return SurplusValueMass(rate.SurplusLabour * int64(totalVariableCapital) / rate.NecessaryLabour)
}

// MassByWorkers computes S = P × (a′/a) × n — total surplus-value from per-worker
// labour-power value, the rate, and the number of workers [§1, formula II].
func MassByWorkers(v LabourPowerValue, rate SurplusValueRate, n WorkerCount) SurplusValueMass {
	if rate.NecessaryLabour == 0 {
		return 0
	}
	return SurplusValueMass(int64(v) * rate.SurplusLabour * int64(n) / rate.NecessaryLabour)
}

// MinimumCapital returns the minimum variable capital required to employ n workers
// each at daily reproduction cost v [§1].
// V_min = v × n
func MinimumCapital(v LabourPowerValue, n WorkerCount) VariableCapital {
	return VariableCapital(int64(v) * int64(n))
}
