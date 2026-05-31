// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 24 ("Externalisation of the Relations of Capital in
// the Form of Interest-Bearing Capital").
//
// Interest-bearing capital in the form M—M' is the most fetishised form of
// capital: money appears to breed money automatically, erasing all trace of
// production and surplus-labour. The compound-interest formula s=c(1+i)^n
// captures this fetish mathematically — Dr. Price's fantasy of a penny
// at 5% compound accumulating over 1776 years to exceed the gold value of
// 150 million earths illustrates the absurdity of mistaking arithmetic for
// economic law. The qualitative limit is the working day: surplus-value is
// bounded by living labour, not by a mathematical formula.
//
// Key invariants:
//   - CompoundInterestSchedule.FinalValue == integer iteration of
//     v_{k+1} = roundHalfUp(v_k*(10000+RateBP), 10000) for PeriodYears steps.
//   - CompoundInterestSchedule.NewValueCreated == 0 (fetish form).
//   - RateBP > 0 and PeriodYears > 0 and Principal > 0.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 24 sentinel errors.
var (
	// ErrNonPositivePrincipal is returned when principal <= 0.
	ErrNonPositivePrincipal = errors.New("credit: principal must be positive")

	// ErrNonPositiveRateBP is returned when rate_bp <= 0.
	ErrNonPositiveRateBP = errors.New("credit: rate_bp must be positive")

	// ErrNonPositivePeriodYears is returned when period_years <= 0.
	ErrNonPositivePeriodYears = errors.New("credit: period_years must be positive")
)

// ── CompoundInterestID ────────────────────────────────────────────────────────

// CompoundInterestID identifies a persisted compound-interest schedule.
// 96-bit hex from crypto/rand, matching every other domain ID in the simulator.
type CompoundInterestID string

// NewCompoundInterestID returns a fresh random CompoundInterestID.
func NewCompoundInterestID() CompoundInterestID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return CompoundInterestID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CompoundInterestID.
func (id CompoundInterestID) IsZero() bool { return id == "" }

// ── FetishForm ────────────────────────────────────────────────────────────────

// FetishForm classifies the degree of capital-fetishism in an interest form.
type FetishForm int

const (
	// FormSimpleInterest is the form where interest is calculated only on
	// the principal — one step removed from surplus-value.
	FormSimpleInterest FetishForm = 1
	// FormCompoundInterest is the compound form where interest compounds on
	// itself — the fetish intensifies as the formula obscures labour further.
	FormCompoundInterest FetishForm = 2
	// FormMM is the pure M—M' form: money advancing to return more money
	// with no visible intermediate circuit — the highest fetish form.
	FormMM FetishForm = 3
)

// ── AccumulationLimit ─────────────────────────────────────────────────────────

// AccumulationLimit expresses the qualitative boundary that the working day
// places on real accumulation regardless of the compound-interest formula.
// No mathematical series can exceed the surplus-labour time actually available.
type AccumulationLimit struct {
	WorkdayMinutes int64 `json:"workday_minutes"` // total working day; >= 0
}

// ── CompoundInterestSchedule ──────────────────────────────────────────────────

// CompoundInterestSchedule is the persisted entity for a compound-interest
// schedule (Vol. III Ch. 24). FinalValue is derived by integer iteration;
// NewValueCreated is always 0 — the fetish form conceals but does not create
// value.
//
// Invariants:
//
//	Principal > 0
//	RateBP > 0
//	PeriodYears > 0
//	FinalValue == integer compound (see computeFinalValue)
//	NewValueCreated == 0
type CompoundInterestSchedule struct {
	ID              CompoundInterestID `json:"id"`
	Principal       int64              `json:"principal"`         // > 0; base amount in smallest currency unit
	RateBP          int64              `json:"rate_bp"`           // > 0; annual interest rate in basis points
	PeriodYears     int64              `json:"period_years"`      // > 0; number of compounding periods
	FinalValue      int64              `json:"final_value"`       // DERIVED: integer compound
	NewValueCreated int64              `json:"new_value_created"` // always 0 — fetish form
	FetishForm      FetishForm         `json:"fetish_form"`       // classification
	CreatedAt       time.Time          `json:"created_at"`
}

// computeFinalValue applies integer fixed-point compound interest:
// start with v = principal; repeat periodYears times:
//
//	v = roundHalfUp(v*(10000+rateBP), 10000)
func computeFinalValue(principal, rateBP, periodYears int64) int64 {
	v := principal
	scale := basisPoints + rateBP
	for i := int64(0); i < periodYears; i++ {
		v = roundHalfUp(v*scale, basisPoints)
	}
	return v
}

// NewCompoundInterestSchedule constructs a CompoundInterestSchedule,
// deriving FinalValue by integer iteration. NewValueCreated is always 0.
//
//	principal <= 0   → ErrNonPositivePrincipal
//	rateBP <= 0      → ErrNonPositiveRateBP
//	periodYears <= 0 → ErrNonPositivePeriodYears
//
// Fixture (Marx, Vol. III Ch. 24 — Dr. Price):
//
//	Principal=1 pence, RateBP=500 (5%), PeriodYears=100
//	→ FinalValue=13150 (1 × 1.05^100 ≈ 131.5 pence, integer iteration)
func NewCompoundInterestSchedule(principal, rateBP, periodYears int64, form FetishForm) (CompoundInterestSchedule, error) {
	if principal <= 0 {
		return CompoundInterestSchedule{}, ErrNonPositivePrincipal
	}
	if rateBP <= 0 {
		return CompoundInterestSchedule{}, ErrNonPositiveRateBP
	}
	if periodYears <= 0 {
		return CompoundInterestSchedule{}, ErrNonPositivePeriodYears
	}
	return CompoundInterestSchedule{
		Principal:       principal,
		RateBP:          rateBP,
		PeriodYears:     periodYears,
		FinalValue:      computeFinalValue(principal, rateBP, periodYears),
		NewValueCreated: 0,
		FetishForm:      form,
	}, nil
}
