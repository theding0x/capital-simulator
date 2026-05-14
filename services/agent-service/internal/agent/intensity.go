package agent

import "errors"

// LabourIntensity is the intensive magnitude of labour: a multiplier on
// value created per unit of labour-time. Capital Vol. I, Ch. 17 §2.
//
// At normal social intensity Factor = 1.0 and a working day of N minutes
// creates N minutes of value. Intensified labour (Factor > 1.0) packs
// more expended labour into the same hour and so creates more value per
// minute. Marx treats this as raising the "average socially necessary"
// floor; the simulation keeps it as a free factor on the scenario input.
type LabourIntensity struct {
	Factor float64 `json:"factor"`
}

// NormalLabourIntensity is the reference value used by the §1 law that
// holds daily value constant under productivity variation.
var NormalLabourIntensity = LabourIntensity{Factor: 1.0}

var ErrIntensityFactorNonPositive = errors.New("agent: labour intensity factor must be positive")

func (li LabourIntensity) Validate() error {
	if li.Factor <= 0 {
		return ErrIntensityFactorNonPositive
	}
	return nil
}

// IsNormal reports whether intensity equals the social normal (Factor == 1.0).
func (li LabourIntensity) IsNormal() bool { return li.Factor == 1.0 }
