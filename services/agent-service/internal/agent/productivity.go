package agent

import "errors"

// LabourProductivity is a multiplicative change in the productiveness of
// labour: a Factor > 1.0 means a unit of labour-time yields proportionally
// more use-values than before. Capital Vol. I, Ch. 17 §1.
//
// Productivity does not change the value created by a working day at
// normal intensity — only how that value partitions between necessary
// and surplus labour. A rise in productivity shrinks the necessary
// labour-time required to reproduce labour-power's value (§1 law 2:
// surplus-value and value of labour-power move in opposite directions
// when productivity changes).
//
// Mirrors commodity.ProductivityChange.Factor on the commodity side; this
// type lives in agent because Ch. 17 frames productivity from the worker's
// labour-power-reproduction perspective rather than the commodity's value.
type LabourProductivity struct {
	Factor float64 `json:"factor"`
}

// NormalLabourProductivity is the no-change reference (Factor == 1.0).
var NormalLabourProductivity = LabourProductivity{Factor: 1.0}

var ErrProductivityFactorNonPositive = errors.New("agent: labour productivity factor must be positive")

func (lp LabourProductivity) Validate() error {
	if lp.Factor <= 0 {
		return ErrProductivityFactorNonPositive
	}
	return nil
}

// ApplyToNecessaryLabour returns the new necessary labour-time after a
// productivity change: necessary' = necessary / factor. Inversely
// proportional, per Marx §1.
func (lp LabourProductivity) ApplyToNecessaryLabour(necessary LabourMinutes) LabourMinutes {
	if lp.Factor <= 0 {
		return necessary
	}
	return LabourMinutes(float64(necessary) / lp.Factor)
}
