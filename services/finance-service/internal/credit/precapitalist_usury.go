// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 36 ("Pre-Capitalist Relationships").
//
// Usury is the historical form of interest-bearing capital that predates
// industrial capitalism. In antiquity and the Middle Ages it exploited
// petty producers and the wealthy alike at rates that could ruin a peasant
// or a nobleman. As industrial capital rose it subordinated usury: the usurer
// became a mere department of the credit system, lending to industrial
// capitalists at rates disciplined by the average rate of profit.
//
// Key concepts:
//   - UsuryDevelopmentStage distinguishes Antiquity, Mediaeval, Transition, and
//     Subordinated phases.
//   - UsurersCapital is a persisted entity recording the usurer's capital at a
//     historical stage, its borrowers, and the interest rate in basis points.
//   - PreCapitalistCreditForm is a value type naming an early credit instrument
//     (temple loans, montes pietatis, etc.).
//   - UsurySubordination records the disciplining of usury under industrial capital.
//   - HistoricalInterestRate records an observed rate in a historical period.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 36 sentinel errors.
var (
	// ErrInvalidUsuryStage is returned when a UsuryDevelopmentStage is outside 1..4.
	ErrInvalidUsuryStage = errors.New("credit: invalid usury development stage")

	// ErrWageLabourInPreCapitalist is returned when wage labour is asserted for
	// a stage that predates industrial capitalism (Antiquity or Mediaeval).
	ErrWageLabourInPreCapitalist = errors.New("credit: wage labour cannot exist in pre-capitalist usury stages")

	// ErrNegativeHistoricalRate is returned when interest_rate_bp < 0.
	ErrNegativeHistoricalRate = errors.New("credit: historical interest rate must be >= 0")

	// ErrMissingSubordination is returned when stage is StageSubordinated but
	// subordinated_to is not "industrial_capital".
	ErrMissingSubordination = errors.New("credit: subordinated stage requires subordinated_to == \"industrial_capital\"")
)

// ── UsuryDevelopmentStage ─────────────────────────────────────────────────────

// UsuryDevelopmentStage enumerates the four historical phases of usury
// (Vol. III Ch. 36).
type UsuryDevelopmentStage int

const (
	// StageAntiquity is usury in the ancient world (Rome, Greece) — exploiting
	// small producers and the nobility before wage labour exists.
	StageAntiquity UsuryDevelopmentStage = 1
	// StageMediaeval is usury under feudalism — church-prohibited but ubiquitous,
	// financing kings and peasants alike.
	StageMediaeval UsuryDevelopmentStage = 2
	// StageTransition is usury during the transition to capitalism — the merchant
	// and proto-industrial phases where wage labour begins to emerge.
	StageTransition UsuryDevelopmentStage = 3
	// StageSubordinated is usury subordinated under industrial capital — the usurer
	// becomes a department of the banking system, rates are disciplined by the
	// average rate of profit.
	StageSubordinated UsuryDevelopmentStage = 4
)

// IsValid reports whether s is a defined UsuryDevelopmentStage constant (1..4).
func (s UsuryDevelopmentStage) IsValid() bool { return s >= 1 && s <= 4 }

// ── UsurersCapitalID ──────────────────────────────────────────────────────────

// UsurersCapitalID identifies a persisted UsurersCapital.
// 96-bit hex from crypto/rand.
type UsurersCapitalID string

// NewUsurersCapitalID returns a fresh random UsurersCapitalID.
func NewUsurersCapitalID() UsurersCapitalID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return UsurersCapitalID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty UsurersCapitalID.
func (id UsurersCapitalID) IsZero() bool { return id == "" }

// ── UsurersCapital ────────────────────────────────────────────────────────────

// UsurersCapital is a persisted entity recording a form of usurers' capital at
// a specific historical stage (Vol. III Ch. 36). ID and CreatedAt are assigned
// by the store.
//
// Invariants:
//
//	Stage.IsValid()
//	WageLabour == false when Stage ∈ {StageAntiquity, StageMediaeval}
//	InterestRateBP >= 0
//	Stage == StageSubordinated → SubordinatedTo == "industrial_capital"
type UsurersCapital struct {
	ID             UsurersCapitalID      `json:"id"`
	Name           string                `json:"name"`
	Stage          UsuryDevelopmentStage `json:"stage"`
	WageLabour     bool                  `json:"wage_labour"`
	Borrowers      []string              `json:"borrowers"`
	InterestRateBP int64                 `json:"interest_rate_bp"`
	SubordinatedTo string                `json:"subordinated_to"`
	Description    string                `json:"description"`
	CreatedAt      time.Time             `json:"created_at"`
}

// NewUsurersCapital constructs a UsurersCapital. ID and CreatedAt are left zero
// — the store assigns them. Validation is fail-fast in declaration order.
//
//	!stage.IsValid()                                        → ErrInvalidUsuryStage
//	wageLabour && stage ∈ {StageAntiquity,StageMediaeval}  → ErrWageLabourInPreCapitalist
//	interestRateBP < 0                                      → ErrNegativeHistoricalRate
//	stage == StageSubordinated && subordinatedTo != "industrial_capital" → ErrMissingSubordination
func NewUsurersCapital(
	name string,
	stage UsuryDevelopmentStage,
	wageLabour bool,
	borrowers []string,
	interestRateBP int64,
	subordinatedTo, description string,
) (UsurersCapital, error) {
	if !stage.IsValid() {
		return UsurersCapital{}, ErrInvalidUsuryStage
	}
	if wageLabour && (stage == StageAntiquity || stage == StageMediaeval) {
		return UsurersCapital{}, ErrWageLabourInPreCapitalist
	}
	if interestRateBP < 0 {
		return UsurersCapital{}, ErrNegativeHistoricalRate
	}
	if stage == StageSubordinated && subordinatedTo != "industrial_capital" {
		return UsurersCapital{}, ErrMissingSubordination
	}
	return UsurersCapital{
		Name:           name,
		Stage:          stage,
		WageLabour:     wageLabour,
		Borrowers:      borrowers,
		InterestRateBP: interestRateBP,
		SubordinatedTo: subordinatedTo,
		Description:    description,
	}, nil
}

// ── PreCapitalistCreditFormID ─────────────────────────────────────────────────

// PreCapitalistCreditFormID identifies a PreCapitalistCreditForm value.
// 96-bit hex from crypto/rand.
type PreCapitalistCreditFormID string

// NewPreCapitalistCreditFormID returns a fresh random PreCapitalistCreditFormID.
func NewPreCapitalistCreditFormID() PreCapitalistCreditFormID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return PreCapitalistCreditFormID(hex.EncodeToString(b[:]))
}

// ── PreCapitalistCreditForm ───────────────────────────────────────────────────

// PreCapitalistCreditForm names a historical credit instrument that predates
// the capitalist credit system (temple loans, montes pietatis, hawala, etc.).
// It is a value type: ID is set by the constructor; it is not persisted
// independently.
type PreCapitalistCreditForm struct {
	ID          PreCapitalistCreditFormID `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
}

// NewPreCapitalistCreditForm constructs a PreCapitalistCreditForm with a fresh
// ID. Infallible.
func NewPreCapitalistCreditForm(name, description string) PreCapitalistCreditForm {
	return PreCapitalistCreditForm{
		ID:          NewPreCapitalistCreditFormID(),
		Name:        name,
		Description: description,
	}
}

// ── UsurySubordination ────────────────────────────────────────────────────────

// UsurySubordination records the disciplining of usury under industrial capital
// (Vol. III Ch. 36). It is a value type with no constructor — compose directly.
type UsurySubordination struct {
	SubordinatedTo string `json:"subordinated_to"`
	Disciplined    bool   `json:"disciplined"`
	Description    string `json:"description"`
}

// ── HistoricalInterestRate ────────────────────────────────────────────────────

// HistoricalInterestRate records an observed interest rate in a named historical
// period, expressed in basis points (Vol. III Ch. 36).
//
// Invariant: RateBP >= 0.
type HistoricalInterestRate struct {
	Period      string `json:"period"`
	RateBP      int64  `json:"rate_bp"`
	Description string `json:"description"`
}

// NewHistoricalInterestRate constructs a HistoricalInterestRate.
//
//	rateBP < 0 → ErrNegativeHistoricalRate
func NewHistoricalInterestRate(period string, rateBP int64, description string) (HistoricalInterestRate, error) {
	if rateBP < 0 {
		return HistoricalInterestRate{}, ErrNegativeHistoricalRate
	}
	return HistoricalInterestRate{
		Period:      period,
		RateBP:      rateBP,
		Description: description,
	}, nil
}
