package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// PhysicalMaxMinutes is the hard ceiling for any working day: 24 h × 60 min.
const PhysicalMaxMinutes int64 = 1440

var (
	ErrWorkingDayExceedsPhysicalMax    = errors.New("agent: working day exceeds physical maximum of 24 hours")
	ErrWorkingDayExceedsStatutoryLimit = errors.New("agent: working day exceeds statutory limit")
)

// Named value types for the two segments of the working day.
type NecessaryLabourMinutes int64
type SurplusLabourMinutes int64
type StatutoryLimitMinutes int64

// WorkingDayID is a 96-bit hex identifier for stored WorkingDay records.
type WorkingDayID string

func (id WorkingDayID) IsZero() bool { return id == "" }

func NewWorkingDayID() WorkingDayID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return WorkingDayID(hex.EncodeToString(b))
}

// WorkingDay encodes the two segments of the working day: AB (necessary) and BC (surplus).
type WorkingDay struct {
	ID                     WorkingDayID           `json:"id"`
	NecessaryLabourMinutes NecessaryLabourMinutes `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   SurplusLabourMinutes   `json:"surplus_labour_minutes"`
	CreatedAt              time.Time              `json:"created_at"`
}

// TotalMinutes returns necessary + surplus labour time in minutes.
func (wd WorkingDay) TotalMinutes() int64 {
	return int64(wd.NecessaryLabourMinutes) + int64(wd.SurplusLabourMinutes)
}

// RateOfSurplusValue returns s/v — the degree of exploitation.
func (wd WorkingDay) RateOfSurplusValue() float64 {
	return float64(wd.SurplusLabourMinutes) / float64(wd.NecessaryLabourMinutes)
}

// Validate checks that necessary labour is positive and the total does not exceed
// the physical maximum of 24 hours.
func (wd WorkingDay) Validate() error {
	if wd.NecessaryLabourMinutes <= 0 {
		return errors.New("agent: necessary_labour_minutes must be positive")
	}
	if wd.TotalMinutes() > PhysicalMaxMinutes {
		return ErrWorkingDayExceedsPhysicalMax
	}
	return nil
}

// WorkingDayConstraint pairs a statutory limit with the jurisdiction or epoch it applies to.
type WorkingDayConstraint struct {
	Label                 string                `json:"label"`
	StatutoryLimitMinutes StatutoryLimitMinutes `json:"statutory_limit_minutes"`
}

// ValidateConstraint returns ErrWorkingDayExceedsStatutoryLimit when the working day
// total exceeds the constraint's statutory limit.
func ValidateConstraint(wd WorkingDay, c WorkingDayConstraint) error {
	if wd.TotalMinutes() > int64(c.StatutoryLimitMinutes) {
		return ErrWorkingDayExceedsStatutoryLimit
	}
	return nil
}

// NormalWorkingDay is a WorkingDay validated against a statutory limit.
type NormalWorkingDay struct {
	WorkingDay
	Constraint WorkingDayConstraint `json:"constraint"`
}

// Validate checks the working day invariants and the statutory limit.
func (n NormalWorkingDay) Validate() error {
	if err := n.WorkingDay.Validate(); err != nil {
		return err
	}
	return ValidateConstraint(n.WorkingDay, n.Constraint)
}

// Overwork represents minutes stolen above the statutory limit per working day.
type Overwork struct {
	MinutesPerDay int64 `json:"minutes_per_day"`
}

// AnnualMinutes returns total overwork minutes across workingDays working days per year.
func (o Overwork) AnnualMinutes(workingDays int) int64 {
	return o.MinutesPerDay * int64(workingDays)
}

// FactoryAct records a historical limit on working time.
type FactoryAct struct {
	Year                    int   `json:"year"`
	ChildLimitMinutes       int64 `json:"child_limit_minutes"`
	YoungPersonLimitMinutes int64 `json:"young_person_limit_minutes"`
	AdultLimitMinutes       int64 `json:"adult_limit_minutes"`
}

// ShiftKind distinguishes day and night shifts in a relay system.
type ShiftKind string

const (
	ShiftDay   ShiftKind = "day"
	ShiftNight ShiftKind = "night"
)

// TimerClass maps to adult vs. child workers in the relay system.
type TimerClass string

const (
	TimerFull TimerClass = "full"
	TimerHalf TimerClass = "half"
)

// RelaySet is one half of a relay system: a shift kind, a working day, and the set of workers.
type RelaySet struct {
	ShiftKind  ShiftKind  `json:"shift_kind"`
	WorkingDay WorkingDay `json:"working_day"`
	WorkerIDs  []AgentID  `json:"worker_ids"`
}

// RelayScheduleID is a 96-bit hex identifier for stored RelaySchedule records.
type RelayScheduleID string

func (id RelayScheduleID) IsZero() bool { return id == "" }

func NewRelayScheduleID() RelayScheduleID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return RelayScheduleID(hex.EncodeToString(b))
}

// RelaySchedule is two alternating worker sets covering complementary shifts.
type RelaySchedule struct {
	ID        RelayScheduleID `json:"id"`
	Sets      [2]RelaySet     `json:"sets"`
	CreatedAt time.Time       `json:"created_at"`
}

// Validate ensures each set's working day is individually valid and their combined
// total does not exceed the physical maximum of 24 hours.
func (rs RelaySchedule) Validate() error {
	if err := rs.Sets[0].WorkingDay.Validate(); err != nil {
		return err
	}
	if err := rs.Sets[1].WorkingDay.Validate(); err != nil {
		return err
	}
	combined := rs.Sets[0].WorkingDay.TotalMinutes() + rs.Sets[1].WorkingDay.TotalMinutes()
	if combined > PhysicalMaxMinutes {
		return ErrWorkingDayExceedsPhysicalMax
	}
	return nil
}
