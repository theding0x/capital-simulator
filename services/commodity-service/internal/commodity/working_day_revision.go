package commodity

import "errors"

// MaxWorkingDayMinutes is the 24 hours of a natural day, the physical bound on
// any working-day (Ch. 10).
const MaxWorkingDayMinutes LabourMinutes = 24 * 60

var (
	// ErrWorkingDayBelowNecessary is returned when a working-day is shorter
	// than the necessary labour: the surplus labour would be negative [§3].
	ErrWorkingDayBelowNecessary = errors.New("commodity: working day is shorter than the necessary labour")
	// ErrWorkingDayExceedsNaturalDay is returned when a working-day is longer
	// than the 24 hours of a natural day [§3].
	ErrWorkingDayExceedsNaturalDay = errors.New("commodity: working day exceeds the 24 hours of a natural day")
)

// RevisedWorkingDay is a working-day re-cut with the necessary labour held
// fixed, and the rate of surplus-value that results [§3].
type RevisedWorkingDay struct {
	WorkingDay      LabourMinutes      `json:"working_day"`
	NecessaryLabour NecessaryLabour    `json:"necessary_labour"`
	SurplusLabour   SurplusLabour      `json:"surplus_labour"`
	Rate            RateOfSurplusValue `json:"rate"`
	RateBasisPoints int64              `json:"rate_basis_points"`
}

// ReviseWorkingDay re-cuts the working-day to newDay with necessary held
// fixed: only the surplus labour moves, and the rate of surplus-value is
// surplus labour over necessary labour [§3].
//
// This is Marx's answer to Senior's "last hour". Senior claimed the whole net
// profit comes from the last hour of an 11½-hour day, so that a 13-hour day
// would more than double it and a one-hour cut would destroy it. On
// Senior's own figures the workman produces his wages in 5¾ hours and the net
// profit in the other 5¾ (a rate of 100%). Lengthening the day to 13 hours
// lifts the surplus labour to 7¼ hours, a rate of 126 2/23%; cutting it by an
// hour leaves 4¾ hours, a rate of 82 14/23%.
//
// It returns ErrDivisionByZero if necessary is not positive,
// ErrWorkingDayExceedsNaturalDay if newDay is longer than 24 hours, and
// ErrWorkingDayBelowNecessary if newDay is shorter than necessary. A day
// exactly as long as the necessary labour is accepted and has a zero rate.
func ReviseWorkingDay(necessary NecessaryLabour, newDay LabourMinutes) (RevisedWorkingDay, error) {
	if necessary <= 0 {
		return RevisedWorkingDay{}, ErrDivisionByZero
	}
	if newDay > MaxWorkingDayMinutes {
		return RevisedWorkingDay{}, ErrWorkingDayExceedsNaturalDay
	}
	if newDay < LabourMinutes(necessary) {
		return RevisedWorkingDay{}, ErrWorkingDayBelowNecessary
	}
	surplus := SurplusLabour(newDay - LabourMinutes(necessary))
	rate, err := ComputeRate(SurplusValue(surplus), LabourMinutes(necessary))
	if err != nil {
		return RevisedWorkingDay{}, err
	}
	bp, err := ComputeRateBasisPoints(SurplusValue(surplus), LabourMinutes(necessary))
	if err != nil {
		return RevisedWorkingDay{}, err
	}
	return RevisedWorkingDay{
		WorkingDay:      newDay,
		NecessaryLabour: necessary,
		SurplusLabour:   surplus,
		Rate:            rate,
		RateBasisPoints: bp,
	}, nil
}
