package market

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Pence is the canonical money-magnitude unit for Vol. I Ch. 3 records.
// Money-prices in Marx's Ch. 3 examples are always integer pence; storing
// them as int64 avoids floating-point drift in arithmetic on hoards,
// obligations, and the circulation-volume formula M = ΣP / V.
type Pence int64

// HoardID uniquely identifies a Hoard.
type HoardID string

// NewHoardID returns a new random 96-bit hex hoard ID.
func NewHoardID() HoardID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("market: rand: %w", err))
	}
	return HoardID(hex.EncodeToString(b[:]))
}

// IsZero reports whether the ID is the empty string.
func (id HoardID) IsZero() bool { return id == "" }

// Hoard is money withdrawn from circulation as a store of value.
//
// "Money becomes petrified into a hoard, and the seller of commodities
// becomes a hoarder of money." (Capital Vol. I, Ch. 3 §3.a)
//
// The hoarder breaks the rhythm of C—M—C by stopping at M. The amount in
// pence is the gold-equivalent quantity withheld.
type Hoard struct {
	ID        HoardID   `json:"id"`
	OwnerID   OwnerID   `json:"owner_id"`
	Amount    Pence     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate checks that h is well-formed.
func (h Hoard) Validate() error {
	if h.OwnerID.IsZero() {
		return errors.New("market: hoard owner_id is required")
	}
	if h.Amount <= 0 {
		return errors.New("market: hoard amount must be positive")
	}
	return nil
}
