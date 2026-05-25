// Package turnover — Vol. II Ch. 7: The Turnover Time and the Number of Turnovers.
//
// Marx establishes that the turnover of capital is the periodically repeated circuit,
// formalising n = T/t (number of turnovers = year / turnover-time). Only the
// money-form (LensMoney) and productive-form (LensProductive) circuits are admissible
// for turnover analysis; the commodity-capital form (LensCommodity) is explicitly
// rejected because it does not begin with the advance of capital-value.
package turnover

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// YearReferenceMinutes is the natural unit of measure for turnover time.
// "Just as the working day is the natural unit for measuring the function of
// labour-power, so the year is the natural unit for measuring the turnovers
// of functioning capital." 365 × 24 × 60 = 525,600 minutes.
const YearReferenceMinutes int64 = 525_600

// Sentinel errors.
var (
	// ErrLensNotApplicableForTurnover is returned when Form III (commodity-capital)
	// is used to create a Turnover. Marx: "But it is not to be used in connection
	// with the turnover of capital, which always begins with the advance of
	// capital-value."
	ErrLensNotApplicableForTurnover = errors.New("commodity-capital lens (Form III) cannot be used for turnover analysis")

	// ErrCycleNotComplete is returned by RecordCycle when the returned capital has
	// not yet recovered its advance in the starting form.
	ErrCycleNotComplete = errors.New("turnover cycle is not complete: returned pence < advance pence")

	// ErrNoCompletedCycles is returned when RecomputeNumber is called on a Turnover
	// with no recorded cycles.
	ErrNoCompletedCycles = errors.New("cannot recompute turnover number: no completed cycles recorded")
)

// --- ID types ----------------------------------------------------------------

// TurnoverID identifies a Turnover.
type TurnoverID string

// NewTurnoverID returns a 96-bit random hex identifier.
func NewTurnoverID() TurnoverID {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return TurnoverID(hex.EncodeToString(b))
}

// TurnoverCycleID identifies a single completed TurnoverCycle.
type TurnoverCycleID string

// NewTurnoverCycleID returns a 96-bit random hex identifier.
func NewTurnoverCycleID() TurnoverCycleID {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return TurnoverCycleID(hex.EncodeToString(b))
}

// TurnoverNumberID identifies a persisted TurnoverNumber snapshot.
type TurnoverNumberID string

// NewTurnoverNumberID returns a 96-bit random hex identifier.
func NewTurnoverNumberID() TurnoverNumberID {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return TurnoverNumberID(hex.EncodeToString(b))
}

// --- Pence -------------------------------------------------------------------

// Pence is the canonical monetary unit (1 GBP = 240 pence).
type Pence int64

// --- CircuitLens -------------------------------------------------------------

// CircuitLens is the circuit-form used to read a Turnover.
// Only LensMoney and LensProductive are valid for turnover analysis.
type CircuitLens string

const (
	// LensMoney (Form I: M…M') — used when studying surplus-value formation.
	// Marx: "Circuit I is of service in a study primarily of the influence
	// of the turnover on the formation of surplus-value."
	LensMoney CircuitLens = "money"

	// LensProductive (Form II: P…P) — used when studying product formation.
	// Marx: "Circuit II … in a study of its influence on the creation of the product."
	LensProductive CircuitLens = "productive"

	// LensCommodity (Form III: C'…C') — explicitly not usable for turnover.
	LensCommodity CircuitLens = "commodity"
)

// ValidateLens returns ErrLensNotApplicableForTurnover for Form III.
func ValidateLens(l CircuitLens) error {
	if l == LensCommodity {
		return ErrLensNotApplicableForTurnover
	}
	if l != LensMoney && l != LensProductive {
		return errors.New("unknown circuit lens: " + string(l))
	}
	return nil
}

// --- Turnover ----------------------------------------------------------------

// Turnover represents the periodically repeated circuit of capital.
// "An individual circuit is but a constantly repeated section in the life of a capital."
type Turnover struct {
	ID                  TurnoverID        `json:"id"`
	IndustrialCapitalID string            `json:"industrial_capital_id"`
	Lens                CircuitLens       `json:"lens"`
	TurnoverTimeMinutes int64             `json:"turnover_time_minutes"`
	Cycles              []TurnoverCycleID `json:"cycles"`
}

// --- TurnoverCycle -----------------------------------------------------------

// TurnoverCycle is one completed period of the periodically repeated circuit.
// Closing condition: ReturnedPence ≥ AdvancePence AND the capital is back in
// the starting functional form of its Lens.
type TurnoverCycle struct {
	ID                 TurnoverCycleID `json:"id"`
	TurnoverID         TurnoverID      `json:"turnover_id"`
	StartedAt          time.Time       `json:"started_at"`
	EndedAt            time.Time       `json:"ended_at"`
	AdvancePence       Pence           `json:"advance_pence"`
	ReturnedPence      Pence           `json:"returned_pence"`
	ProductionMinutes  int64           `json:"production_minutes"`
	CirculationMinutes int64           `json:"circulation_minutes"`
}

// TotalMinutes returns ProductionMinutes + CirculationMinutes.
func (c TurnoverCycle) TotalMinutes() int64 {
	return c.ProductionMinutes + c.CirculationMinutes
}

// Validate checks the close-condition for a completed cycle.
func (c TurnoverCycle) Validate() error {
	if c.ReturnedPence < c.AdvancePence {
		return ErrCycleNotComplete
	}
	return nil
}

// --- TurnoverNumber ----------------------------------------------------------

// TurnoverNumber captures the formula n = T/t as integer basis-points.
// "If we designate the year as the unit of measure of the turnover time by T,
// the time of turnover of a given capital by t, and the number of its turnovers
// by n, then n = T/t." — Marx.
//
// BasisPoints = 10000 × T / t (integer division). This avoids float64 while
// preserving the exact ratio; numerator and denominator are both kept for exact
// reconstruction and human display.
type TurnoverNumber struct {
	ID                   TurnoverNumberID `json:"id"`
	TurnoverID           TurnoverID       `json:"turnover_id"`
	YearReferenceMinutes int64            `json:"year_reference_minutes"`
	TurnoverTimeMinutes  int64            `json:"turnover_time_minutes"`
	BasisPoints          int64            `json:"basis_points"`
}

// ComputeTurnoverNumber builds a TurnoverNumber for the given turnover-time.
// BasisPoints = 10000 × YearReferenceMinutes / turnoverTimeMinutes (integer division).
func ComputeTurnoverNumber(turnoverID TurnoverID, turnoverTimeMinutes int64) TurnoverNumber {
	var bp int64
	if turnoverTimeMinutes > 0 {
		bp = 10_000 * YearReferenceMinutes / turnoverTimeMinutes
	}
	return TurnoverNumber{
		ID:                   NewTurnoverNumberID(),
		TurnoverID:           turnoverID,
		YearReferenceMinutes: YearReferenceMinutes,
		TurnoverTimeMinutes:  turnoverTimeMinutes,
		BasisPoints:          bp,
	}
}

// AverageTurnoverTime computes the average total minutes across a set of cycles.
// Returns 0 if cycles is empty.
func AverageTurnoverTime(cycles []TurnoverCycle) int64 {
	if len(cycles) == 0 {
		return 0
	}
	var total int64
	for _, c := range cycles {
		total += c.TotalMinutes()
	}
	return total / int64(len(cycles))
}
