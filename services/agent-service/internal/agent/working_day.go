package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/theding0x/capital-simulator/pkg/surplusrate"
)

// PhysicalMaxMinutes is the hard ceiling for any working day: 24 h × 60 min.
const PhysicalMaxMinutes int64 = 1440

var (
	ErrNecessaryLabourRequired         = errors.New("agent: necessary_labour_minutes must be positive")
	ErrSurplusLabourNegative           = errors.New("agent: surplus_labour_minutes must not be negative")
	ErrWorkingDayExceedsPhysicalMax    = errors.New("agent: working day exceeds physical maximum of 24 hours")
	ErrWorkingDayExceedsStatutoryLimit = errors.New("agent: working day exceeds statutory limit")
	ErrNoStatutoryLimit                = errors.New("agent: factory act sets no statutory limit for this age band")
	ErrUnknownAgeBand                  = errors.New("agent: unknown age band")
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

// RateOfSurplusValue returns s/v — the degree of exploitation. Validate
// guarantees necessary labour is positive; on an unvalidated day the rate is
// undefined and reads as 0.
func (wd WorkingDay) RateOfSurplusValue() float64 {
	r, _ := surplusrate.Ratio(int64(wd.SurplusLabourMinutes), int64(wd.NecessaryLabourMinutes))
	return r
}

// Validate checks that necessary labour is positive, surplus labour is not negative
// (BC is an extension beyond AB, never a deduction from it; zero is the limiting case),
// and the total does not exceed the physical maximum of 24 hours.
func (wd WorkingDay) Validate() error {
	if wd.NecessaryLabourMinutes <= 0 {
		return ErrNecessaryLabourRequired
	}
	if wd.SurplusLabourMinutes < 0 {
		return ErrSurplusLabourNegative
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

// FactoryAct records the daily limits on working time in force once an Act is
// fully operative; later Acts inherit the limits they leave untouched. 1833:
// children 8 h (480), young persons 12 h (720). 1844: children 6½ h (390).
// 1847: young persons 10 h (600), definitive from 1 May 1848; children still 6½ h.
//
// A limit of 0 (or less) means "no statutory limit", not a zero-minute cap. None of the
// Acts of 1833, 1844 and 1847 limited the working day of the male worker over 18
// (Vol. I Ch. 10 §6), so AdultLimitMinutes is 0 in every historical Act. Women
// over 18 are covered through YoungPersonLimitMinutes (the 1844 and 1847 Acts put
// them on the young-person footing); there is deliberately no sex dimension.
//
// Only the general daily limit is modelled. A single daily number cannot
// represent the 7 h child exception, the silk-mill children, the 1845 Print Works
// Act, the 11 h preliminary limit of July 1847, the age phasing of 1834-36, or
// the 1850 Act (10½ h on weekdays, 7½ h on Saturday).
//
// It is a domain-only fixture: not persisted (migration 00008 has no limit column
// and the 1847 seed day carries no constraint), not exposed over HTTP (the
// handlers take a raw statutory_limit_minutes), and Constraint currently has no
// production caller.
type FactoryAct struct {
	Year                    int   `json:"year"`
	ChildLimitMinutes       int64 `json:"child_limit_minutes"`
	YoungPersonLimitMinutes int64 `json:"young_person_limit_minutes"`
	AdultLimitMinutes       int64 `json:"adult_limit_minutes"`
}

// AgeBand is the class of worker a Factory Act limits differently.
type AgeBand string

const (
	// AgeBandChild is a child under 13, limited to 8 h from 1833 and 6½ h from 1844.
	AgeBandChild AgeBand = "child"
	// AgeBandYoungPerson is a young person (13 to 18). Women over 18 belong here
	// too from 1844: the Acts put them on the young-person footing.
	AgeBandYoungPerson AgeBand = "young_person"
	// AgeBandAdult is the male worker over 18, whom none of the Acts limited. An
	// adult woman is AgeBandYoungPerson, not this band.
	AgeBandAdult AgeBand = "adult"
)

// Constraint derives the WorkingDayConstraint the Act sets for band, so the
// statutory number comes from the Act rather than a retyped literal. It returns
// ErrUnknownAgeBand for an unrecognised band and ErrNoStatutoryLimit where the
// Act sets no limit: ValidateConstraint with a zero limit would reject every
// positive day, so "no limit" must be an error, never a zero constraint.
func (fa FactoryAct) Constraint(band AgeBand) (WorkingDayConstraint, error) {
	var limit int64
	var who string
	switch band {
	case AgeBandChild:
		limit, who = fa.ChildLimitMinutes, "child"
	case AgeBandYoungPerson:
		limit, who = fa.YoungPersonLimitMinutes, "young person"
	case AgeBandAdult:
		limit, who = fa.AdultLimitMinutes, "adult"
	default:
		return WorkingDayConstraint{}, fmt.Errorf("%w: %q", ErrUnknownAgeBand, band)
	}
	if limit <= 0 {
		return WorkingDayConstraint{}, fmt.Errorf("%w: %d %s", ErrNoStatutoryLimit, fa.Year, band)
	}
	return WorkingDayConstraint{
		Label:                 fmt.Sprintf("Factory Act %d (%s)", fa.Year, who),
		StatutoryLimitMinutes: StatutoryLimitMinutes(limit),
	}, nil
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
