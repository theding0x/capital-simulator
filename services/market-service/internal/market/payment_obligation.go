package market

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// PaymentObligationID uniquely identifies a PaymentObligation.
type PaymentObligationID string

// NewPaymentObligationID returns a new random 96-bit hex obligation ID.
func NewPaymentObligationID() PaymentObligationID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("market: rand: %w", err))
	}
	return PaymentObligationID(hex.EncodeToString(b[:]))
}

// IsZero reports whether the ID is the empty string.
func (id PaymentObligationID) IsZero() bool { return id == "" }

// PaymentObligation is a deferred payment — money has functioned as means
// of payment, separating the moment of sale from the moment of settlement.
//
// "The seller turned his commodity into money, in order to satisfy by its
// means some want of his own; the hoarder, in order to keep it in the
// money-form, and the buyer-in-debt, in order to be able to pay." (Ch. 3 §3.b)
//
// The obligation exists from CreatedAt until PaidAt; while unpaid it is a
// claim on the debtor's future product.
type PaymentObligation struct {
	ID         PaymentObligationID `json:"id"`
	CreditorID OwnerID             `json:"creditor_id"`
	DebtorID   OwnerID             `json:"debtor_id"`
	Amount     Pence               `json:"amount"`
	CreatedAt  time.Time           `json:"created_at"`
	DueAt      time.Time           `json:"due_at"`
	PaidAt     *time.Time          `json:"paid_at,omitempty"`
}

// Validate checks that p is well-formed.
func (p PaymentObligation) Validate() error {
	if p.CreditorID.IsZero() {
		return errors.New("market: payment_obligation creditor_id is required")
	}
	if p.DebtorID.IsZero() {
		return errors.New("market: payment_obligation debtor_id is required")
	}
	if p.CreditorID == p.DebtorID {
		return errors.New("market: payment_obligation creditor and debtor cannot be the same owner")
	}
	if p.Amount <= 0 {
		return errors.New("market: payment_obligation amount must be positive")
	}
	if p.DueAt.IsZero() {
		return errors.New("market: payment_obligation due_at is required")
	}
	return nil
}

// IsSettled reports whether the obligation has been paid.
func (p PaymentObligation) IsSettled() bool { return p.PaidAt != nil }
