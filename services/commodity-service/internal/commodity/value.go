package commodity

import (
	"errors"
	"fmt"
	"math"
)

// Quantity is a non-negative amount of a commodity in its socially-recognized
// unit (e.g. yards of linen, qtrs of corn). It is allowed to be fractional.
type Quantity float64

// Value computes the value of a given quantity of c, expressed in
// LabourMinutes. It is SNLTPerUnit * quantity, with the concrete labour that
// produced c reduced to abstract human labour-time first.
//
// Capital I, Ch. 1, §1: "A use value, or useful article, therefore, has
// value only because human labour in the abstract has been embodied or
// materialised in it. How, then, is the magnitude of this value to be
// measured? Plainly, by the quantity of the value-creating substance, the
// labour, contained in the article."
func (c Commodity) Value(qty Quantity) LabourMinutes {
	if qty < 0 {
		qty = 0
	}
	abstract := AsAbstractLabour(c.ConcreteLabour, c.SNLTPerUnit)
	v := float64(abstract) * float64(qty)
	return LabourMinutes(math.Round(v))
}

// ExchangeRatio represents a pairwise exchange relation: BaseQty units of Base
// exchange for QuoteQty units of Quote at parity of value.
//
// Capital I, Ch. 1, §1: "1 quarter corn = x cwt iron." The equation expresses
// "that in two different things ... there exists in equal quantities something
// common to both" - namely, the abstract human labour-time congealed in each.
type ExchangeRatio struct {
	Base     Commodity `json:"base"`
	Quote    Commodity `json:"quote"`
	BaseQty  Quantity  `json:"base_qty"`
	QuoteQty Quantity  `json:"quote_qty"`
	// CommonValue is the labour-time both sides equate to.
	CommonValue LabourMinutes `json:"common_value"`
}

// ExchangeRatioBetween derives the simple equivalence x A = y B for the given
// base quantity. "First: the valid exchange values of a given commodity express
// something equal" (Ch. 1 §1).
//
// If quote.SNLTPerUnit is zero (a free good) the function returns an error
// because no finite quantity of it equates to a positive value of base.
func ExchangeRatioBetween(base Commodity, quote Commodity, baseQty Quantity) (ExchangeRatio, error) {
	if baseQty <= 0 {
		return ExchangeRatio{}, errors.New("commodity: base_qty must be positive")
	}
	if quote.SNLTPerUnit <= 0 {
		return ExchangeRatio{}, fmt.Errorf("commodity: quote %q has no value (snlt_per_unit=%d); no finite ratio exists", quote.Name, quote.SNLTPerUnit)
	}
	common := base.Value(baseQty)
	quoteQty := Quantity(float64(common) / float64(quote.SNLTPerUnit))
	return ExchangeRatio{
		Base:        base,
		Quote:       quote,
		BaseQty:     baseQty,
		QuoteQty:    quoteQty,
		CommonValue: common,
	}, nil
}

// ProductivityChange records that the productivity of the concrete labour
// producing c has changed by Factor. Factor > 1 means productivity rose
// (less SNLT required per unit); 0 < Factor < 1 means productivity fell.
//
// Capital I, Ch. 1, §1: "in general, the greater the productiveness of
// labour, the less is the labour-time required for the production of an
// article, the less is the amount of labour crystallised in that article,
// and the less is its value; and vice versa, the less the productiveness
// of labour, the greater is the labour-time required for the production
// of an article, and the greater is its value. The value of a commodity,
// therefore, varies directly as the quantity, and inversely as the
// productiveness, of the labour incorporated in it."
type ProductivityChange struct {
	// Factor is the multiplicative change in productivity. SNLTPerUnit is
	// divided by Factor: doubled productivity halves SNLT.
	Factor float64
}

// Apply returns a copy of c with SNLTPerUnit divided by p.Factor.
// Returns an error for non-positive Factor.
func (p ProductivityChange) Apply(c Commodity) (Commodity, error) {
	if p.Factor <= 0 {
		return Commodity{}, errors.New("commodity: productivity factor must be positive")
	}
	out := c
	out.SNLTPerUnit = LabourMinutes(math.Round(float64(c.SNLTPerUnit) / p.Factor))
	return out, nil
}
