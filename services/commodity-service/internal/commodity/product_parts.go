package commodity

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrEmptyAccount is returned when an account's c+v+s is zero: there is no
	// value to apportion the product by [§2].
	ErrEmptyAccount = errors.New("commodity: production account has no value to apportion the product by")
	// ErrNegativeProduct is returned for a total product below zero [§2].
	ErrNegativeProduct = errors.New("commodity: total product must be non-negative")
)

// ProductParts is the total product divided into the parts in which the
// constant capital, the variable capital and the surplus-value are
// respectively embodied [§2].
//
// Marx: "the various component parts of this value, can be represented as
// being contained respectively in corresponding parts of the product." Of
// 20 lbs. of yarn worth 24s. c + 3s. v + 3s. s, 16 lbs. are "nothing more than
// so much cotton and spindle-waste in disguise", and the remaining 4 lbs.
// hold all the value newly created: 2 lbs. the equivalent of wages, 2 lbs. the
// surplus-value.
type ProductParts struct {
	Constant Quantity `json:"constant"`
	Variable Quantity `json:"variable"`
	Surplus  Quantity `json:"surplus"`
}

// ApportionProduct divides total, a quantity of the product of a, into the
// parts embodying a's c, v and s. Each part is the total times the share of
// the value it embodies, so the parts are proportional to the value
// components and together make up the whole product [§2].
//
// It returns ErrEmptyAccount if c+v+s is zero, ErrNegativeProduct if total is
// below zero, and a's own validation error if a is malformed.
func ApportionProduct(a ProductionAccount, total Quantity) (ProductParts, error) {
	if err := a.Validate(); err != nil {
		return ProductParts{}, err
	}
	if math.IsNaN(float64(total)) || math.IsInf(float64(total), 0) {
		return ProductParts{}, fmt.Errorf("commodity: total product %v is not a finite number", total)
	}
	if total < 0 {
		return ProductParts{}, ErrNegativeProduct
	}
	// Sum in float64: c+v+s of three int64 magnitudes can exceed MaxInt64.
	c, v, s := float64(a.Constant), float64(a.Variable), float64(a.Surplus)
	value := c + v + s
	if value == 0 {
		return ProductParts{}, ErrEmptyAccount
	}
	// Scale the total by each part's share of the value. A share is at most 1,
	// so a finite total can never overflow; multiplying total by the part
	// before dividing by the whole could.
	product := float64(total)
	return ProductParts{
		Constant: Quantity(product * (c / value)),
		Variable: Quantity(product * (v / value)),
		Surplus:  Quantity(product * (s / value)),
	}, nil
}
