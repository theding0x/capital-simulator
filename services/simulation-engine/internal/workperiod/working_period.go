// Package workperiod models Vol. II Ch. 12 — The Working Period.
// A working period is the span of connected working-days required to produce
// one finished commodity. Discrete processes (e.g. yarn) yield a finished unit
// each day; connected processes (e.g. a locomotive) yield nothing until the
// full period is complete. The key result: circulating-capital advance scales
// linearly with the working-period length for connected processes; natural
// constraints (annual crops, multi-year timber) set inviolable floors.
package workperiod

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// ErrPrematureDelivery is returned when a working period's working_day_count
// falls below the natural minimum for the specified commodity.
var ErrPrematureDelivery = errors.New("working_period: working_day_count is below the natural minimum for this commodity")

// newHexID returns a 96-bit (12-byte) random identifier as a hex string.
func newHexID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}

// --- ID types ---------------------------------------------------------------

// WorkingPeriodID is a 96-bit hex identifier for a WorkingPeriod.
type WorkingPeriodID string

// NewWorkingPeriodID returns a new random WorkingPeriodID.
func NewWorkingPeriodID() WorkingPeriodID { return WorkingPeriodID(newHexID()) }

// WorkingPeriodInterruptionID is a 96-bit hex identifier for an interruption.
type WorkingPeriodInterruptionID string

// NewWorkingPeriodInterruptionID returns a new random WorkingPeriodInterruptionID.
func NewWorkingPeriodInterruptionID() WorkingPeriodInterruptionID {
	return WorkingPeriodInterruptionID(newHexID())
}

// WorkingPeriodShorteningID is a 96-bit hex identifier for a shortening.
type WorkingPeriodShorteningID string

// NewWorkingPeriodShorteningID returns a new random WorkingPeriodShorteningID.
func NewWorkingPeriodShorteningID() WorkingPeriodShorteningID {
	return WorkingPeriodShorteningID(newHexID())
}

// NaturalWorkingPeriodConstraintID is a 96-bit hex identifier for a constraint.
type NaturalWorkingPeriodConstraintID string

// NewNaturalWorkingPeriodConstraintID returns a new random NaturalWorkingPeriodConstraintID.
func NewNaturalWorkingPeriodConstraintID() NaturalWorkingPeriodConstraintID {
	return NaturalWorkingPeriodConstraintID(newHexID())
}

// WorkingPeriodCreditFinancingID is a 96-bit hex identifier for a credit-financing record.
type WorkingPeriodCreditFinancingID string

// NewWorkingPeriodCreditFinancingID returns a new random WorkingPeriodCreditFinancingID.
func NewWorkingPeriodCreditFinancingID() WorkingPeriodCreditFinancingID {
	return WorkingPeriodCreditFinancingID(newHexID())
}

// --- WorkingPeriodMode ------------------------------------------------------

// WorkingPeriodMode distinguishes how a working period yields its product.
type WorkingPeriodMode string

const (
	// ModeDiscrete: each working-day yields one finished unit (e.g. yarn,
	// coal). The effective capital-advance multiplier is always 1.
	ModeDiscrete WorkingPeriodMode = "discrete"
	// ModeConnected: many working-days must be completed before any product
	// exists (e.g. locomotive, ship, road). The capital-advance multiplier
	// equals WorkingDayCount.
	ModeConnected WorkingPeriodMode = "connected"
)

// --- WorkingPeriod ----------------------------------------------------------

// WorkingPeriod is the span of connected working-days required to produce one
// finished commodity. Multiplier() expresses how many days' worth of
// circulating capital must be advanced simultaneously.
type WorkingPeriod struct {
	ID                      WorkingPeriodID               `json:"id"`
	IndustrialCapitalID     string                        `json:"industrial_capital_id"`
	CommodityID             string                        `json:"commodity_id"`
	WorkingDayCount         int64                         `json:"working_day_count"`
	Mode                    WorkingPeriodMode             `json:"mode"`
	WorkingDayLengthMinutes int64                         `json:"working_day_length_minutes"`
	CreatedAt               time.Time                     `json:"created_at"`
	Interruptions           []WorkingPeriodInterruption   `json:"interruptions,omitempty"`
	Shortenings             []WorkingPeriodShortening     `json:"shortenings,omitempty"`
	CreditFinancing         *WorkingPeriodCreditFinancing `json:"credit_financing,omitempty"`
}

// Multiplier returns the capital-advance multiplier for circulating capital.
// In ModeConnected, daily outlays accumulate for WorkingDayCount days before
// any product exists; in ModeDiscrete, only a single day's outlay is tied up.
func (wp WorkingPeriod) Multiplier() int64 {
	if wp.Mode == ModeConnected {
		return wp.WorkingDayCount
	}
	return 1
}

// Validate reports whether the WorkingPeriod is internally consistent.
func (wp WorkingPeriod) Validate() error {
	if wp.IndustrialCapitalID == "" {
		return errors.New("working_period: industrial_capital_id is required")
	}
	if wp.CommodityID == "" {
		return errors.New("working_period: commodity_id is required")
	}
	if wp.WorkingDayCount <= 0 {
		return errors.New("working_period: working_day_count must be positive")
	}
	if wp.Mode != ModeDiscrete && wp.Mode != ModeConnected {
		return errors.New("working_period: mode must be 'discrete' or 'connected'")
	}
	if wp.WorkingDayLengthMinutes <= 0 {
		return errors.New("working_period: working_day_length_minutes must be positive")
	}
	return nil
}

// --- WorkingPeriodInterruption ----------------------------------------------

// WorkingPeriodInterruption records a pause in a connected or discrete working
// period. For discrete processes, DeteriorationPence must be zero — a one-day
// product does not deteriorate during an idle day.
type WorkingPeriodInterruption struct {
	ID                            WorkingPeriodInterruptionID `json:"id"`
	WorkingPeriodID               WorkingPeriodID             `json:"working_period_id"`
	InterruptedAt                 time.Time                   `json:"interrupted_at"`
	ResumedAt                     *time.Time                  `json:"resumed_at,omitempty"`
	DeteriorationPence            int64                       `json:"deterioration_pence"`
	WasteOfMeansOfProductionPence int64                       `json:"waste_of_means_of_production_pence"`
}

// Validate enforces invariants for an interruption. The caller must pass the
// owning WorkingPeriod's mode so the discrete-process rule can be checked.
func (i WorkingPeriodInterruption) Validate(mode WorkingPeriodMode) error {
	if mode == ModeDiscrete && i.DeteriorationPence != 0 {
		return errors.New("working_period: discrete interruptions cannot carry deterioration_pence")
	}
	if i.DeteriorationPence < 0 {
		return errors.New("working_period: deterioration_pence cannot be negative")
	}
	if i.WasteOfMeansOfProductionPence < 0 {
		return errors.New("working_period: waste_of_means_of_production_pence cannot be negative")
	}
	return nil
}

// --- WorkingPeriodShorteningKind --------------------------------------------

// WorkingPeriodShorteningKind names the mechanism by which a working period
// is shortened. Each kind carries distinct additional-cost invariants.
type WorkingPeriodShorteningKind string

const (
	// ShorteningCooperation: parallel workers reduce elapsed calendar time
	// without additional fixed capital.
	ShorteningCooperation WorkingPeriodShorteningKind = "cooperation"
	// ShorteningMachinery: new fixed capital replaces calendar time.
	// AdditionalFixedCapitalPence must be positive.
	ShorteningMachinery WorkingPeriodShorteningKind = "machinery"
	// ShorteningParallelLabour: additional workers run multiple operations
	// simultaneously. AdditionalLabourCount must be positive.
	ShorteningParallelLabour WorkingPeriodShorteningKind = "parallel_labour"
	// ShorteningBakewellBreeding: selective breeding shortens the natural
	// maturation period (Marx's Bakewell cattle example).
	ShorteningBakewellBreeding WorkingPeriodShorteningKind = "bakewell_breeding"
)

// --- WorkingPeriodShortening ------------------------------------------------

// WorkingPeriodShortening records a reduction of the working period by one
// of the four recognised mechanisms.
type WorkingPeriodShortening struct {
	ID                          WorkingPeriodShorteningID   `json:"id"`
	WorkingPeriodID             WorkingPeriodID             `json:"working_period_id"`
	Kind                        WorkingPeriodShorteningKind `json:"kind"`
	OldWorkingDayCount          int64                       `json:"old_working_day_count"`
	NewWorkingDayCount          int64                       `json:"new_working_day_count"`
	AdditionalFixedCapitalPence int64                       `json:"additional_fixed_capital_pence"`
	AdditionalLabourCount       int64                       `json:"additional_labour_count"`
	OccurredAt                  time.Time                   `json:"occurred_at"`
}

// Validate enforces the invariants specific to each shortening kind.
func (s WorkingPeriodShortening) Validate() error {
	if s.OldWorkingDayCount <= 0 || s.NewWorkingDayCount <= 0 {
		return errors.New("working_period: old_working_day_count and new_working_day_count must be positive")
	}
	if s.NewWorkingDayCount >= s.OldWorkingDayCount {
		return errors.New("working_period: new_working_day_count must be less than old_working_day_count")
	}
	if s.Kind == ShorteningMachinery && s.AdditionalFixedCapitalPence <= 0 {
		return errors.New("working_period: machinery shortening requires additional_fixed_capital_pence > 0")
	}
	if s.Kind == ShorteningParallelLabour && s.AdditionalLabourCount <= 0 {
		return errors.New("working_period: parallel-labour shortening requires additional_labour_count > 0")
	}
	return nil
}

// --- NaturalWorkingPeriodConstraint -----------------------------------------

// NaturalWorkingPeriodConstraint records the minimum working-period length for
// commodities subject to natural constraints (annual crops, multi-year timber,
// livestock gestation). Creating a WorkingPeriod with WorkingDayCount below
// MinWorkingDays for the constrained commodity returns ErrPrematureDelivery.
type NaturalWorkingPeriodConstraint struct {
	ID             NaturalWorkingPeriodConstraintID `json:"id"`
	CommodityID    string                           `json:"commodity_id"`
	MinWorkingDays int64                            `json:"min_working_days"`
	Reason         string                           `json:"reason"`
	CreatedAt      time.Time                        `json:"created_at"`
}

// Validate reports whether the constraint is well-formed.
func (nc NaturalWorkingPeriodConstraint) Validate() error {
	if nc.CommodityID == "" {
		return errors.New("working_period: commodity_id is required for natural constraint")
	}
	if nc.MinWorkingDays <= 0 {
		return errors.New("working_period: min_working_days must be positive")
	}
	if nc.Reason == "" {
		return errors.New("working_period: reason is required for natural constraint")
	}
	return nil
}

// --- WorkingPeriodCreditFinancing -------------------------------------------

// WorkingPeriodCreditFinancing records how a long connected working period is
// financed. OwnCapitalPence covers the fraction the capitalist advances from
// reserves; BorrowedCapitalPence covers the rest (e.g. a mortgage, as in
// Marx's London builder example). Vol. III Part Five will model interest
// charges; MortgageProviderID is a forward-reference placeholder.
type WorkingPeriodCreditFinancing struct {
	ID                   WorkingPeriodCreditFinancingID `json:"id"`
	WorkingPeriodID      WorkingPeriodID                `json:"working_period_id"`
	OwnCapitalPence      int64                          `json:"own_capital_pence"`
	BorrowedCapitalPence int64                          `json:"borrowed_capital_pence"`
	MortgageProviderID   string                         `json:"mortgage_provider_id,omitempty"`
	CreatedAt            time.Time                      `json:"created_at"`
}

// Validate reports whether the credit-financing record is well-formed.
func (cf WorkingPeriodCreditFinancing) Validate() error {
	if cf.OwnCapitalPence < 0 {
		return errors.New("working_period: own_capital_pence cannot be negative")
	}
	if cf.BorrowedCapitalPence < 0 {
		return errors.New("working_period: borrowed_capital_pence cannot be negative")
	}
	if cf.OwnCapitalPence+cf.BorrowedCapitalPence <= 0 {
		return errors.New("working_period: total financing (own + borrowed) must be positive")
	}
	return nil
}
