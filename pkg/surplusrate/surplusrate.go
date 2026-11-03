// Package surplusrate is the single definition of the rate of surplus-value
// s' = s/v (Capital Vol. I Ch. 9) shared by every service. Magnitudes are
// labour-minute integers; the rate is either a float ratio (Ratio) or a scaled
// integer in basis points rounded half-up (BasisPoints), the convention Vol. III
// already uses for p' and s'.
package surplusrate

import "errors"

// BasisPointsPerUnit is the scale of BasisPoints: a rate of 1.0 (100%) is 10000.
const BasisPointsPerUnit = 10_000

// ErrZeroVariable is returned when v <= 0: surplus-value originates in variable
// capital, so without it there is no rate of exploitation to measure. Callers
// that cannot surface an error treat the rate as 0; callers at a validating
// boundary map it to a client error.
var ErrZeroVariable = errors.New("surplusrate: variable capital (necessary labour) must be positive")

// Ratio returns s/v as a float64. It returns ErrZeroVariable when v <= 0.
func Ratio(s, v int64) (float64, error) {
	if v <= 0 {
		return 0, ErrZeroVariable
	}
	return float64(s) / float64(v), nil
}

// BasisPoints returns s/v scaled by BasisPointsPerUnit and rounded half-up, so
// the 300% English agricultural labourer is 30000. It returns ErrZeroVariable
// when v <= 0.
func BasisPoints(s, v int64) (int64, error) {
	if v <= 0 {
		return 0, ErrZeroVariable
	}
	return (s*BasisPointsPerUnit + v/2) / v, nil
}
