package circulation

import (
	"errors"
	"math"
	"time"
)

// ExitCommodityCapital is C' as the exit form of one productive circuit's P:
// the commodity-product that carries the capital-value consumed in P plus the
// surplus-value created there (Vol. II Ch. 2 §I). Spinning Mill 1871: 10,000
// lbs of yarn valued at £500 is ValueOriginal 42,200 (8,440 lbs) plus
// ValueSurplus 7,800 (1,560 lbs). It is keyed by ProductiveCircuitID, unlike
// Ch. 1's CommodityCapital, which is keyed by MoneyCircuitID.
type ExitCommodityCapital struct {
	ProductiveCircuitID ProductiveCircuitID `json:"productive_circuit_id"`
	ValueOriginal       Pence               `json:"value_original"`
	ValueSurplus        Pence               `json:"value_surplus"`
}

// Total returns the value of C' (C + c): the sum a sale must realise.
func (c ExitCommodityCapital) Total() Pence { return c.ValueOriginal + c.ValueSurplus }

// Validate checks that C' carries a positive capital-value and a non-negative
// surplus-value, and that their sum is representable: a wrapped (negative)
// Total would let a one-penny sale pass the "sold in its entirety" guard.
func (c ExitCommodityCapital) Validate() error {
	if c.ValueOriginal <= 0 {
		return errors.New("circulation: value_original must be positive")
	}
	if c.ValueSurplus < 0 {
		return errors.New("circulation: value_surplus must not be negative")
	}
	if c.ValueOriginal > math.MaxInt64-c.ValueSurplus {
		return errors.New("circulation: value_original + value_surplus overflows")
	}
	return nil
}

// ExitRealisation is the C'—M' sale that closes a productive circuit's
// circulation phase. It exists only for a full sale: a rejected partial sale
// is never recorded.
type ExitRealisation struct {
	ProductiveCircuitID ProductiveCircuitID `json:"productive_circuit_id"`
	SoldAt              time.Time           `json:"sold_at"`
	RealisedPence       Pence               `json:"realised_pence"`
}

// ValidateRealisation checks a proposed C'—M' sale against the circuit's
// recorded C'. "C' must be sold at its value and in its entirety" (Vol. II
// Ch. 2 §I), so the sale must realise at least the whole value of C'; selling
// only the capital-value yarn leaves the surplus-value unrealised.
func (p ProductiveCircuit) ValidateRealisation(realised Pence) error {
	return ValidateRealisation(p.ExitCommodity, realised)
}

// ValidateRealisation is the state-free form of
// ProductiveCircuit.ValidateRealisation, for stores that hold C' as columns
// rather than as a ProductiveCircuit value. It returns ErrNoCommodityCapital
// when commodity is nil and ErrPartialRealisation when realised falls short of
// commodity.Total(). Realising more than the value is accepted here: a sale
// above value belongs to Vol. II Ch. 15.
func ValidateRealisation(commodity *ExitCommodityCapital, realised Pence) error {
	if commodity == nil {
		return ErrNoCommodityCapital
	}
	if realised < commodity.Total() {
		return ErrPartialRealisation
	}
	return nil
}
