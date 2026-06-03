package market

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// ErrCircuitNotEquivalent is returned when a Circuit's two legs do not carry
// the same value-magnitude. C—M—C is a circuit of equivalents: the same value
// leaves in commodity-form (the sale) and returns in commodity-form (the
// purchase). A mismatch means the recorded movement is not a simple-circulation
// circuit at all.
var ErrCircuitNotEquivalent = errors.New("market: circuit is not a circuit of equivalents (sale value != purchase value)")

// CircuitID uniquely identifies a Circuit.
type CircuitID string

// NewCircuitID returns a new random 96-bit hex circuit ID.
func NewCircuitID() CircuitID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("market: rand: %w", err))
	}
	return CircuitID(hex.EncodeToString(b[:]))
}

// IsZero reports whether the ID is the empty string.
func (id CircuitID) IsZero() bool { return id == "" }

// Validate checks that a single leg of a circuit is well-formed.
func (l CircuitLeg) Validate() error {
	switch l.Kind {
	case KindSale, KindPurchase:
	default:
		return fmt.Errorf("market: circuit leg kind must be %q or %q", KindSale, KindPurchase)
	}
	if l.OwnerID.IsZero() {
		return errors.New("market: circuit leg owner_id is required")
	}
	if l.CommodityID == "" {
		return errors.New("market: circuit leg commodity_id is required")
	}
	if l.MoneyID == "" {
		return errors.New("market: circuit leg money_id is required")
	}
	if l.Value <= 0 {
		return errors.New("market: circuit leg value must be positive")
	}
	return nil
}

// Circuit is the two-leg C—M—C movement of simple circulation: a sale (C → M)
// followed by a purchase (M → C).
//
// "The circuit C—M—C starts with one commodity, and finishes with another,
// which falls out of circulation and into consumption. … its final result is
// the exchange of one commodity for another." (Capital Vol. I, Ch. 3 §2.a)
//
// The defining invariant is that the two legs carry the same value-magnitude:
// the weaver sells 20 yards of linen for £2 and lays out that same £2 on a
// family Bible. Value is conserved across the circuit; only its use-value form
// changes (linen → Bible).
type Circuit struct {
	ID          CircuitID  `json:"id"`
	SaleLeg     CircuitLeg `json:"sale_leg"`
	PurchaseLeg CircuitLeg `json:"purchase_leg"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Validate checks that c is a well-formed circuit of equivalents.
func (c Circuit) Validate() error {
	if c.SaleLeg.Kind != KindSale {
		return fmt.Errorf("market: sale_leg kind must be %q", KindSale)
	}
	if c.PurchaseLeg.Kind != KindPurchase {
		return fmt.Errorf("market: purchase_leg kind must be %q", KindPurchase)
	}
	if err := c.SaleLeg.Validate(); err != nil {
		return err
	}
	if err := c.PurchaseLeg.Validate(); err != nil {
		return err
	}
	// The load-bearing Ch. 3 invariant: the same value-magnitude leaves in
	// commodity-form and returns in commodity-form.
	if c.SaleLeg.Value != c.PurchaseLeg.Value {
		return ErrCircuitNotEquivalent
	}
	return nil
}
