// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 35 ("Precious Metal and Rate of Exchange").
//
// Gold serves as world money in three functions: as the universal measure of
// value and means of payment (world money proper), as a reserve fund against
// domestic credit crises (reserve fund), and as the adjustment mechanism for
// balances of payments (international adjustment). The rate of exchange
// deviates from par — expressed in basis points — when balances of payments
// are persistently unfavourable or favourable. An adverse balance of payments
// drains gold; a favourable balance attracts it.
//
// Key concepts:
//   - WorldMoneyFunction distinguishes the three roles of gold at the world level.
//   - GoldReserve is a persisted entity recording the BoE gold stock at a moment.
//   - RateOfExchange is a persisted entity recording the deviation of the exchange
//     rate from mint par in basis points (negative = unfavourable for home country).
//   - BalanceOfPayments is a value type capturing net surplus or deficit.
//   - GoldFlow is a value type capturing a directional movement of gold.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 35 sentinel errors.
var (
	// ErrNegativeGoldReserveAmount is returned when a GoldReserve amount is negative.
	ErrNegativeGoldReserveAmount = errors.New("credit: gold reserve amount must be >= 0")

	// ErrInvalidGoldFlowDirection is returned when a GoldFlowDirection is not 1 or 2.
	ErrInvalidGoldFlowDirection = errors.New("credit: invalid gold flow direction")

	// ErrNegativeGoldFlowAmount is returned when a GoldFlow amount is negative.
	ErrNegativeGoldFlowAmount = errors.New("credit: gold flow amount must be >= 0")

	// ErrInvalidWorldMoneyFunction is returned when the world money function is not 1, 2, or 3.
	ErrInvalidWorldMoneyFunction = errors.New("credit: invalid world money function")
)

// ── WorldMoneyFunction ────────────────────────────────────────────────────────

// WorldMoneyFunction enumerates the three roles gold plays as world money
// (Vol. III Ch. 35).
type WorldMoneyFunction int

const (
	// FunctionWorldMoney is gold as universal means of payment and measure of value.
	FunctionWorldMoney WorldMoneyFunction = 1
	// FunctionReserveFund is gold held as a reserve fund against domestic credit crises.
	FunctionReserveFund WorldMoneyFunction = 2
	// FunctionAdjustment is gold as the mechanism for adjusting international balances.
	FunctionAdjustment WorldMoneyFunction = 3
)

// IsValid reports whether f is a defined WorldMoneyFunction constant.
func (f WorldMoneyFunction) IsValid() bool { return f >= 1 && f <= 3 }

// ── GoldFlowDirection ─────────────────────────────────────────────────────────

// GoldFlowDirection indicates whether gold is flowing into or out of a country.
type GoldFlowDirection int

const (
	// FlowInward is gold flowing into the domestic economy (favourable balance).
	FlowInward GoldFlowDirection = 1
	// FlowOutward is gold flowing out of the domestic economy (adverse balance).
	FlowOutward GoldFlowDirection = 2
)

// IsValid reports whether d is a defined GoldFlowDirection constant.
func (d GoldFlowDirection) IsValid() bool { return d >= 1 && d <= 2 }

// ── GoldReserveID ─────────────────────────────────────────────────────────────

// GoldReserveID identifies a persisted GoldReserve.
// 96-bit hex from crypto/rand.
type GoldReserveID string

// NewGoldReserveID returns a fresh random GoldReserveID.
func NewGoldReserveID() GoldReserveID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return GoldReserveID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty GoldReserveID.
func (id GoldReserveID) IsZero() bool { return id == "" }

// ── GoldReserve ───────────────────────────────────────────────────────────────

// GoldReserve is a persisted entity recording the gold stock at the Bank of
// England at a specific historical moment (Vol. III Ch. 35). Amount is in
// pence (integer). ID and CreatedAt are assigned by the store.
//
// Invariant: Amount >= 0.
type GoldReserve struct {
	ID          GoldReserveID `json:"id"`
	Name        string        `json:"name"`
	Amount      int64         `json:"amount"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
}

// NewGoldReserve constructs a GoldReserve. ID and CreatedAt are left zero —
// the store assigns them.
//
//	amount < 0 → ErrNegativeGoldReserveAmount
func NewGoldReserve(name string, amount int64, description string) (GoldReserve, error) {
	if amount < 0 {
		return GoldReserve{}, ErrNegativeGoldReserveAmount
	}
	return GoldReserve{
		Name:        name,
		Amount:      amount,
		Description: description,
	}, nil
}

// ── RateOfExchangeID ──────────────────────────────────────────────────────────

// RateOfExchangeID identifies a persisted RateOfExchange.
// 96-bit hex from crypto/rand.
type RateOfExchangeID string

// NewRateOfExchangeID returns a fresh random RateOfExchangeID.
func NewRateOfExchangeID() RateOfExchangeID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return RateOfExchangeID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty RateOfExchangeID.
func (id RateOfExchangeID) IsZero() bool { return id == "" }

// ── RateOfExchange ────────────────────────────────────────────────────────────

// RateOfExchange is a persisted entity recording the deviation of an exchange
// rate from mint par, expressed in basis points (Vol. III Ch. 35). A negative
// DeviationBP indicates an unfavourable rate for the home country; a positive
// value indicates a favourable rate. ID and CreatedAt are assigned by the store.
type RateOfExchange struct {
	ID          RateOfExchangeID `json:"id"`
	Name        string           `json:"name"`
	DeviationBP int64            `json:"deviation_bp"`
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
}

// NewRateOfExchange constructs a RateOfExchange. The sign of DeviationBP is
// unrestricted. ID and CreatedAt are left zero — the store assigns them.
func NewRateOfExchange(name string, deviationBP int64, description string) RateOfExchange {
	return RateOfExchange{
		Name:        name,
		DeviationBP: deviationBP,
		Description: description,
	}
}

// ── BalanceOfPaymentsID ───────────────────────────────────────────────────────

// BalanceOfPaymentsID identifies a BalanceOfPayments value.
type BalanceOfPaymentsID string

// NewBalanceOfPaymentsID returns a fresh random BalanceOfPaymentsID.
func NewBalanceOfPaymentsID() BalanceOfPaymentsID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return BalanceOfPaymentsID(hex.EncodeToString(b[:]))
}

// ── BalanceOfPayments ─────────────────────────────────────────────────────────

// BalanceOfPayments is a value type (not persisted independently) capturing the
// net international balance of payments. A positive NetBalance is favourable
// (surplus); negative is adverse (deficit).
type BalanceOfPayments struct {
	ID          BalanceOfPaymentsID `json:"id"`
	NetBalance  int64               `json:"net_balance"`
	Favourable  bool                `json:"favourable"`
	Description string              `json:"description"`
}

// NewBalanceOfPayments constructs a BalanceOfPayments. ID is set; Favourable
// is derived as netBalance > 0.
func NewBalanceOfPayments(netBalance int64, description string) BalanceOfPayments {
	return BalanceOfPayments{
		ID:          NewBalanceOfPaymentsID(),
		NetBalance:  netBalance,
		Favourable:  netBalance > 0,
		Description: description,
	}
}

// ── GoldFlowID ────────────────────────────────────────────────────────────────

// GoldFlowID identifies a GoldFlow value.
type GoldFlowID string

// NewGoldFlowID returns a fresh random GoldFlowID.
func NewGoldFlowID() GoldFlowID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return GoldFlowID(hex.EncodeToString(b[:]))
}

// ── GoldFlow ──────────────────────────────────────────────────────────────────

// GoldFlow is a value type recording a directional movement of gold — inward
// (FlowInward) or outward (FlowOutward) — with the amount in pence.
//
// Invariants:
//
//	Direction.IsValid()
//	Amount >= 0
type GoldFlow struct {
	ID          GoldFlowID        `json:"id"`
	Direction   GoldFlowDirection `json:"direction"`
	Amount      int64             `json:"amount"`
	Description string            `json:"description"`
}

// NewGoldFlow constructs a GoldFlow. ID is set by the constructor.
//
//	!direction.IsValid() → ErrInvalidGoldFlowDirection
//	amount < 0          → ErrNegativeGoldFlowAmount
func NewGoldFlow(direction GoldFlowDirection, amount int64, description string) (GoldFlow, error) {
	if !direction.IsValid() {
		return GoldFlow{}, ErrInvalidGoldFlowDirection
	}
	if amount < 0 {
		return GoldFlow{}, ErrNegativeGoldFlowAmount
	}
	return GoldFlow{
		ID:          NewGoldFlowID(),
		Direction:   direction,
		Amount:      amount,
		Description: description,
	}, nil
}
