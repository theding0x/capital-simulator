package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// CooperationID is a 96-bit hex identifier for stored Cooperation records.
type CooperationID string

func (id CooperationID) IsZero() bool { return id == "" }

func NewCooperationID() CooperationID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return CooperationID(hex.EncodeToString(b))
}

// CooperationSize is the number of workers assembled simultaneously under one capitalist.
type CooperationSize int

// CooperationMinSize is Marx's minimum platoon for the law of averages to hold
// (Burke fn. 1, §3): five workers cancel out individual deviations.
const CooperationMinSize CooperationSize = 5

var (
	ErrCooperationSizeTooSmall    = errors.New("cooperation: size must be at least 1")
	ErrCooperationCapitalistID    = errors.New("cooperation: capitalist_id is required")
	ErrCooperationWorkingDay      = errors.New("cooperation: working_day_minutes must be positive")
	ErrCooperationNoMembers       = errors.New("cooperation: at least one member is required")
	ErrCooperationMemberInvalid   = errors.New("cooperation: member worker_id is required")
	ErrCooperationProductivePower = errors.New("cooperation: productive_power_factor must be 0 (unrecorded) or at least 1")
)

// Supervisor is a Worker-class member of a Cooperation with directing authority.
// In Marx's terms, the supervisor is still a wage-labourer, not a capitalist —
// his function is the social organisation of cooperative labour.
type Supervisor struct {
	WorkerID AgentID `json:"worker_id"`
}

// CooperationMember is one worker's participation in a Cooperation.
// Each member contributes a working-day measured in LabourMinutes;
// if Supervisory is true the member directs the others on the capitalist's behalf.
type CooperationMember struct {
	WorkerID          AgentID       `json:"worker_id"`
	Supervisory       bool          `json:"supervisory"`
	WorkingDayMinutes LabourMinutes `json:"working_day_minutes"`
}

// Cooperation is a named pool of Worker agents assembled by one Capitalist
// agent. The members all do "the same or the same kind of work" (Ch. 13, §1).
type Cooperation struct {
	ID           CooperationID       `json:"id"`
	Name         string              `json:"name"`
	CapitalistID AgentID             `json:"capitalist_id"`
	Members      []CooperationMember `json:"members"`
	// ProductivePower is the observed use-value output factor of this
	// co-operation relative to the same workers isolated (Ch. 13, §1: "far more
	// than twelve isolated men"). Marx gives no formula for it, so it is
	// recorded, never derived from the headcount. Zero means no gain recorded;
	// PowerFactor then reports exactly 1.0.
	ProductivePower float64   `json:"productive_power_factor"`
	CreatedAt       time.Time `json:"created_at"`
}

// PowerFactor returns the recorded productive-power factor, or 1.0 (no gain
// asserted) when none was recorded.
func (c Cooperation) PowerFactor() float64 {
	if c.ProductivePower == 0 {
		return 1.0
	}
	return c.ProductivePower
}

// Size returns the count of workers in the cooperation.
func (c Cooperation) Size() CooperationSize {
	return CooperationSize(len(c.Members))
}

// Supervisors returns the directing-authority subset of the members.
// Always returns a non-nil slice so JSON marshalling produces [] not null.
func (c Cooperation) Supervisors() []Supervisor {
	out := make([]Supervisor, 0)
	for _, m := range c.Members {
		if m.Supervisory {
			out = append(out, Supervisor{WorkerID: m.WorkerID})
		}
	}
	return out
}

// AverageIndividualWorkingDay returns the mean per-worker working day in minutes,
// rounding down via integer division. This is the empirical input to
// AverageSocialLabour and CollectiveWorkingDay computations for heterogeneous shifts.
func (c Cooperation) AverageIndividualWorkingDay() LabourMinutes {
	n := len(c.Members)
	if n == 0 {
		return 0
	}
	var total int64
	for _, m := range c.Members {
		total += int64(m.WorkingDayMinutes)
	}
	return LabourMinutes(total / int64(n))
}

// Validate enforces the cooperation invariants.
func (c Cooperation) Validate() error {
	if c.CapitalistID.IsZero() {
		return ErrCooperationCapitalistID
	}
	if len(c.Members) == 0 {
		return ErrCooperationNoMembers
	}
	if c.ProductivePower != 0 && c.ProductivePower < 1 {
		return ErrCooperationProductivePower
	}
	for _, m := range c.Members {
		if m.WorkerID.IsZero() {
			return ErrCooperationMemberInvalid
		}
		if m.WorkingDayMinutes <= 0 {
			return ErrCooperationWorkingDay
		}
	}
	return nil
}

// CollectiveWorkingDay returns the sum of individual working-days for n workers
// each contributing d minutes: it is the aggregate labour-time the capitalist
// commands per period. From the standpoint of value (Ch. 13, §2) this scales
// linearly with worker count; no value is created by the cooperative relation itself.
func CollectiveWorkingDay(n CooperationSize, d LabourMinutes) LabourMinutes {
	if n <= 0 || d <= 0 {
		return 0
	}
	return LabourMinutes(int64(n) * int64(d))
}

// AverageSocialLabour reduces the collective working-day back to a per-worker average.
// Individual deviations cancel for platoons of CooperationMinSize or more (Burke fn. 1).
func AverageSocialLabour(collective LabourMinutes, n CooperationSize) LabourMinutes {
	if n <= 0 {
		return 0
	}
	return LabourMinutes(int64(collective) / int64(n))
}

// MinimumCapital returns the least Pence a capitalist must command to pay n
// workers simultaneously at the given daily wage (Ch. 13, §8). The cooperation
// constraint is, before anything else, a capital-minimum constraint: "the law
// that the minimum amount of capital advanced increases with the scale of
// co-operation."
func MinimumCapital(n CooperationSize, dailyWage Pence) Pence {
	if n <= 0 || dailyWage <= 0 {
		return 0
	}
	return Pence(int64(n) * int64(dailyWage))
}

// CooperativeOrigin names where the social productive power of cooperation
// historically appears. It appears as a property of capital (Ch. 13, §13):
// "Because this power costs capital nothing, and because, on the other hand,
// the workman himself does not develop it before his labour belongs to capital,
// it appears as a power with which capital is endowed by nature."
func CooperativeOrigin() string {
	return "capital"
}

// CooperationProductivityFactor returns the recorded productivity gain of a
// cooperation as a multiplicative factor on SNLT. It is the bridge that lets a
// Ch. 13 record drive a Ch. 12 relative-surplus-value calculation: the factor
// is whatever was observed and recorded on the cooperation (1.0 when none), not
// a function of its size. The float64 return type matches the wire format
// expected by simulation-engine's production.ProductivityFactor.
func CooperationProductivityFactor(c Cooperation) float64 {
	return c.PowerFactor()
}
