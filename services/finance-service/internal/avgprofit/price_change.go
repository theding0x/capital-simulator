package avgprofit

// This file carries the first of Ch. 12's two closing remarks on prices of
// production. A commodity's price of production changes for one of only two
// reasons:
//
//   (a) the GENERAL rate of profit changes — even though the commodity's own
//       value is unaltered (a shift somewhere else in the social capital), or
//   (b) the commodity's OWN value changes — a productivity shift in its sphere
//       or in the spheres furnishing its constant capital.
//
// Either suffices; both together is also possible. The price-change record only
// reports WHICH cause is at work, deriving the boolean flags from the cause.
//
// The second remark — compensation grounds — is the CompensationGround value
// object below: a capitalist's claimed "grounds" for a higher price (slow
// turnover, risk, skill) create no new value; they are only a claim on a larger
// share of the aggregate surplus-value already produced.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// PriceChangeCause names the reason a price of production has moved.
type PriceChangeCause string

const (
	// CauseGeneralRateChange: the general rate of profit changed; the
	// commodity's own value is unaltered.
	CauseGeneralRateChange PriceChangeCause = "general_rate_change"
	// CauseIndividualValueChange: the commodity's own value changed (a
	// productivity shift in its sphere or its inputs).
	CauseIndividualValueChange PriceChangeCause = "individual_value_change"
	// CauseBoth: both the general rate and the commodity's own value changed.
	CauseBoth PriceChangeCause = "both"
)

// PriceOfProductionChangeID identifies a persisted price-of-production change.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type PriceOfProductionChangeID string

// NewPriceOfProductionChangeID returns a fresh random PriceOfProductionChangeID.
func NewPriceOfProductionChangeID() PriceOfProductionChangeID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("avgprofit: rand: %w", err))
	}
	return PriceOfProductionChangeID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty PriceOfProductionChangeID.
func (id PriceOfProductionChangeID) IsZero() bool { return id == "" }

// PriceOfProductionChange is the persisted entity recording a movement in a
// sphere's price of production and which of the two admissible causes produced
// it. The RateChanged/ValueChanged/PriceChanged flags are derived from Cause by
// ComputePriceOfProductionChange. OldPrice and NewPrice are whole-unit context.
type PriceOfProductionChange struct {
	ID           PriceOfProductionChangeID `json:"id"`
	SphereName   string                    `json:"sphere_name"`
	Cause        PriceChangeCause          `json:"cause"`
	RateChanged  bool                      `json:"rate_changed"`  // derived from Cause
	ValueChanged bool                      `json:"value_changed"` // derived from Cause
	PriceChanged bool                      `json:"price_changed"` // Cause != ""
	OldPrice     ProductionPrice           `json:"old_price"`     // whole units (optional context)
	NewPrice     ProductionPrice           `json:"new_price"`
	CreatedAt    time.Time                 `json:"created_at"`
}

// Validate returns ErrNegativeMagnitude if either price magnitude is negative.
func (c PriceOfProductionChange) Validate() error {
	if c.OldPrice < 0 || c.NewPrice < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// ComputePriceOfProductionChange derives the boolean flags from the cause:
//
//	RateChanged  = cause == CauseGeneralRateChange || cause == CauseBoth
//	ValueChanged = cause == CauseIndividualValueChange || cause == CauseBoth
//	PriceChanged = cause != ""
//
// It assigns no ID or timestamp — the store does that on persistence.
func ComputePriceOfProductionChange(name string, cause PriceChangeCause, oldPrice, newPrice ProductionPrice) PriceOfProductionChange {
	rateChanged := cause == CauseGeneralRateChange || cause == CauseBoth
	valueChanged := cause == CauseIndividualValueChange || cause == CauseBoth
	priceChanged := cause != ""
	return PriceOfProductionChange{
		SphereName:   name,
		Cause:        cause,
		RateChanged:  rateChanged,
		ValueChanged: valueChanged,
		PriceChanged: priceChanged,
		OldPrice:     oldPrice,
		NewPrice:     newPrice,
	}
}

// compensationActuallyIs is the invariant truth about every compensation
// ground: it is never value creation, only a claim on already-produced surplus.
const compensationActuallyIs = "share in aggregate surplus-value"

// CompensationGround is a value object exposing Ch. 12's second remark. A
// capitalist's claimed "grounds" for a higher price APPEAR to be a creation of
// value but ACTUALLY are only a redistribution — a share in the aggregate
// surplus-value. It is computed on demand and is NOT persisted.
type CompensationGround struct {
	Reason       string `json:"reason"`
	AppearsToBe  string `json:"appears_to_be"`
	ActuallyIs   string `json:"actually_is"`   // always compensationActuallyIs
	CreatesValue bool   `json:"creates_value"` // always false
}

// ComputeCompensationGround returns the compensation ground for a claimed
// reason. Whatever the reason or appearance, ActuallyIs is the constant
// "share in aggregate surplus-value" and CreatesValue is false.
func ComputeCompensationGround(reason, appearsToBe string) CompensationGround {
	return CompensationGround{
		Reason:       reason,
		AppearsToBe:  appearsToBe,
		ActuallyIs:   compensationActuallyIs,
		CreatesValue: false,
	}
}
