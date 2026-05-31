// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 31 ("Money-Capital and Real Capital. II (Continuation)").
//
// Marx distinguishes two sources of loanable money-capital:
//
//  1. The simple transformation of money into loan capital. This draws on
//     industrial capital that has fallen idle — money lying fallow after a crisis,
//     plus the economies in the reserve of means of circulation that banking
//     concentrates. This source is INVERSELY related to real accumulation: it
//     swells precisely when production stagnates and capital cannot find
//     productive employment.
//
//  2. The transformation of capital or revenue into money that then becomes loan
//     capital. This ACCOMPANIES real accumulation: as the reproduction process
//     expands, a growing stream of surplus-value and revenue is momentarily
//     converted into money on its way to fresh investment, and a part is lent.
//
// Bill-brokers (Lombard Street) redistribute the surplus deposits of agricultural
// districts — money that sits idle in country banks — to the manufacturing
// centres by rediscounting bills. The mass of loanable capital is NOT the same as
// the quantity of currency: the same £20 lent five times over performs £100 of
// loan-capital service. Velocity multiplies a fixed mass of money into a far
// larger mass of loan capital.
//
// Key invariants:
//   - FloatingCapital.Amount >= 0 (pence).
//   - FloatingCapital.Source ∈ {1,2} (FloatingCapitalSource.IsValid()).
//   - LoanCapitalVelocity.TimesLent >= 1.
//   - LoanCapitalVelocity.LoanCapitalCreated == MoneyAmount × TimesLent (the
//     constructor computes it so the invariant always holds).
//   - CirculationReserveSaving.Amount >= 0 (pence).
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 31 sentinel errors.
var (
	// ErrNegativeFloatingAmount is returned when floating_capital.amount < 0.
	ErrNegativeFloatingAmount = errors.New("credit: floating-capital amount must be >= 0")

	// ErrInvalidFloatingSource is returned when floating_capital.source is not a
	// valid FloatingCapitalSource.
	ErrInvalidFloatingSource = errors.New("credit: invalid floating-capital source")

	// ErrTimesLentTooLow is returned when loan_capital_velocity.times_lent < 1.
	ErrTimesLentTooLow = errors.New("credit: times_lent must be >= 1")

	// ErrNegativeReserveSaving is returned when circulation_reserve_saving.amount < 0.
	ErrNegativeReserveSaving = errors.New("credit: reserve-saving amount must be >= 0")
)

// ── FloatingCapitalSource ──────────────────────────────────────────────────────

// FloatingCapitalSource names which of Marx's two channels a parcel of loanable
// money-capital flows through.
type FloatingCapitalSource int

const (
	// SourceSimpleTransformation is the simple transformation of money into loan
	// capital — idle post-crisis industrial capital plus banking economies in the
	// reserve of circulation. INVERSELY related to real accumulation.
	SourceSimpleTransformation FloatingCapitalSource = 1
	// SourceCapitalRevenueConversion is the transformation of capital or revenue
	// into money that becomes loan capital. ACCOMPANIES real accumulation.
	SourceCapitalRevenueConversion FloatingCapitalSource = 2
)

// IsValid reports whether s is one of the two recognised sources.
func (s FloatingCapitalSource) IsValid() bool {
	return s == SourceSimpleTransformation || s == SourceCapitalRevenueConversion
}

// ── LoanCapitalVelocityID ──────────────────────────────────────────────────────

// LoanCapitalVelocityID identifies a LoanCapitalVelocity. The velocity is a value
// type embedded in FloatingCapital and is never persisted independently, but it
// carries an ID so it can be referenced like the other domain types.
type LoanCapitalVelocityID string

// NewLoanCapitalVelocityID returns a fresh random LoanCapitalVelocityID.
func NewLoanCapitalVelocityID() LoanCapitalVelocityID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return LoanCapitalVelocityID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty LoanCapitalVelocityID.
func (id LoanCapitalVelocityID) IsZero() bool { return id == "" }

// ── LoanCapitalVelocity ────────────────────────────────────────────────────────

// LoanCapitalVelocity captures Marx's point that the mass of loanable capital is
// not the quantity of currency: the same money, lent repeatedly, performs a far
// larger mass of loan-capital service. The same £20 lent five times over is £100
// of loan capital. It is a value type, never persisted independently.
//
// Invariant:
//
//	TimesLent >= 1
//	LoanCapitalCreated == MoneyAmount × TimesLent
type LoanCapitalVelocity struct {
	MoneyAmount        int64 `json:"money_amount"`         // pence; the stock of money lent
	TimesLent          int64 `json:"times_lent"`           // >= 1
	LoanCapitalCreated int64 `json:"loan_capital_created"` // = MoneyAmount × TimesLent
}

// NewLoanCapitalVelocity builds and validates a LoanCapitalVelocity, computing
// LoanCapitalCreated = moneyAmount × timesLent so the invariant always holds.
//
//	timesLent < 1 → ErrTimesLentTooLow
//
// Fixture (Marx, Vol. III Ch. 31): £20 lent 5 times → £100 of loan capital
// (moneyAmount=2000, timesLent=5, LoanCapitalCreated=10000).
func NewLoanCapitalVelocity(moneyAmount, timesLent int64) (LoanCapitalVelocity, error) {
	if timesLent < 1 {
		return LoanCapitalVelocity{}, ErrTimesLentTooLow
	}
	return LoanCapitalVelocity{
		MoneyAmount:        moneyAmount,
		TimesLent:          timesLent,
		LoanCapitalCreated: moneyAmount * timesLent,
	}, nil
}

// ── CirculationReserveSaving ───────────────────────────────────────────────────

// CirculationReserveSaving is the economy in the reserve of means of circulation
// that banking concentration produces — idle till-money that the credit system
// pools and lends out instead of leaving fallow. It feeds the simple
// transformation of money into loan capital. A value type, never persisted
// independently.
type CirculationReserveSaving struct {
	Amount      int64  `json:"amount"` // pence; >= 0
	Description string `json:"description"`
}

// ── BillBrokerID ───────────────────────────────────────────────────────────────

// BillBrokerID identifies a BillBroker. The broker is a value type — it is not
// persisted — but it carries an ID so it can be referenced like the other domain
// types.
type BillBrokerID string

// NewBillBrokerID returns a fresh random BillBrokerID.
func NewBillBrokerID() BillBrokerID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return BillBrokerID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty BillBrokerID.
func (id BillBrokerID) IsZero() bool { return id == "" }

// ── BillBroker ─────────────────────────────────────────────────────────────────

// BillBroker is a Lombard Street bill-broker who redistributes the surplus
// deposits of agricultural districts to the manufacturing centres by
// rediscounting bills. A value type, never persisted.
type BillBroker struct {
	ID               BillBrokerID `json:"id"`
	Name             string       `json:"name"`
	RediscountVolume int64        `json:"rediscount_volume"` // pence
	Description      string       `json:"description"`
}

// ── Rediscount ─────────────────────────────────────────────────────────────────

// Rediscount records a single bill passed from its original discounter through a
// broker to the centre that ultimately holds it. A value type, never persisted.
type Rediscount struct {
	OriginalDiscounter string `json:"original_discounter"`
	Broker             string `json:"broker"`
	BillAmount         int64  `json:"bill_amount"` // pence
}

// ── FloatingCapitalID ──────────────────────────────────────────────────────────

// FloatingCapitalID identifies a persisted FloatingCapital.
type FloatingCapitalID string

// NewFloatingCapitalID returns a fresh random FloatingCapitalID.
func NewFloatingCapitalID() FloatingCapitalID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return FloatingCapitalID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty FloatingCapitalID.
func (id FloatingCapitalID) IsZero() bool { return id == "" }

// ── FloatingCapital ────────────────────────────────────────────────────────────

// FloatingCapital is the persisted entity for a parcel of loanable money-capital
// (Vol. III Ch. 31). It records the amount, which of the two sources it flows
// through, the velocity that multiplies a fixed money mass into a larger mass of
// loan capital, and the economy in the circulation reserve that feeds the simple
// transformation. The only persisted entity in this chapter.
//
// Invariants:
//
//	Amount >= 0
//	Source.IsValid()
//	Velocity.TimesLent >= 1 (enforced by NewLoanCapitalVelocity)
//	ReserveSaving.Amount >= 0
type FloatingCapital struct {
	ID            FloatingCapitalID        `json:"id"`
	Name          string                   `json:"name"`
	Amount        int64                    `json:"amount"` // pence; >= 0
	Source        FloatingCapitalSource    `json:"source"` // 1 or 2
	Velocity      LoanCapitalVelocity      `json:"loan_capital_velocity"`
	ReserveSaving CirculationReserveSaving `json:"circulation_reserve_saving"`
	Description   string                   `json:"description"`
	CreatedAt     time.Time                `json:"created_at"`
}

// NewFloatingCapital constructs a FloatingCapital, enforcing invariants. ID and
// CreatedAt are left zero — the store assigns them.
//
//	amount < 0                  → ErrNegativeFloatingAmount
//	!source.IsValid()           → ErrInvalidFloatingSource
//	reserveSaving.Amount < 0    → ErrNegativeReserveSaving
//
// The velocity must already be constructed via NewLoanCapitalVelocity, which
// guarantees TimesLent >= 1 and the LoanCapitalCreated identity.
//
// Fixture (Marx, Vol. III Ch. 31): Ipswich farmers — surplus rural deposits
// redistributed to manufacturing (Source=SourceCapitalRevenueConversion,
// accompanies real accumulation).
func NewFloatingCapital(
	name string,
	amount int64,
	source FloatingCapitalSource,
	velocity LoanCapitalVelocity,
	reserveSaving CirculationReserveSaving,
	description string,
) (FloatingCapital, error) {
	if amount < 0 {
		return FloatingCapital{}, ErrNegativeFloatingAmount
	}
	if !source.IsValid() {
		return FloatingCapital{}, ErrInvalidFloatingSource
	}
	if reserveSaving.Amount < 0 {
		return FloatingCapital{}, ErrNegativeReserveSaving
	}
	return FloatingCapital{
		Name:          name,
		Amount:        amount,
		Source:        source,
		Velocity:      velocity,
		ReserveSaving: reserveSaving,
		Description:   description,
	}, nil
}

// RelationToAccumulation returns whether this parcel of loanable capital moves
// inversely to real accumulation (the simple transformation of idle money) or
// accompanies it (the conversion of capital and revenue into money). This is the
// central distinction of the chapter.
func (f FloatingCapital) RelationToAccumulation() string {
	switch f.Source {
	case SourceSimpleTransformation:
		return "inverse to real accumulation — idle money transformed into loan capital swells when production stagnates"
	case SourceCapitalRevenueConversion:
		return "accompanies real accumulation — capital and revenue converted into money on the way to fresh investment"
	default:
		return "unknown source"
	}
}
