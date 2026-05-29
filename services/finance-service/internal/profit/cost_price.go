// Package profit models the conversion of surplus-value into profit and of the
// rate of surplus-value into the rate of profit — Capital Vol. III, Part I.
//
// Ch. 1 ("Cost-Price and Profit") has two movements:
//
//  1. The cost-price k. Of a commodity's value C = c + v + s, the portion that
//     merely replaces the consumed capital (c + v) is what the commodity costs
//     the capitalist. Marx names it k, so C = c + v + s becomes C = k + s.
//  2. The profit-form. Surplus-value, re-attributed from variable capital to
//     the total capital advanced, takes the mystified form of profit p, and
//     C = k + s becomes C = k + p. That transformation lives in profit.go.
//
// This file defines the value magnitudes, the capital outlay, the cost-price,
// and the Vol. III value formula CommodityValueV3.
package profit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// LabourMinutes is the canonical value-magnitude unit across the simulator:
// the socially necessary labour-time congealed in a value. Marx's worked
// examples in this chapter are denominated in money (£, shillings), but he
// fixes the correspondence — "the value produced by one labourer during an
// average social working-day is represented by a money sum of 6s." — so the
// £-figures and the labour-time figures are the same magnitude under two
// names. We keep the labour-time unit throughout.
type LabourMinutes int64

// ConstantCapital (c) is the value transferred from the consumed means of
// production to the product. It reappears in the product's value but creates
// no new value of its own.
type ConstantCapital LabourMinutes

// VariableCapital (v) is the value advanced for labour-power. In the labour
// process, living labour-power replaces it with a greater new value, so v is
// the sole source of surplus-value.
type VariableCapital LabourMinutes

// SurplusValue (s) is the unpaid labour appropriated by the capitalist — the
// excess of the new value created over the variable capital advanced.
type SurplusValue LabourMinutes

// Sentinel validation errors for a CapitalOutlay.
var (
	ErrNegativeMagnitude        = errors.New("profit: value magnitude cannot be negative")
	ErrFixedWearExceedsConstant = errors.New("profit: fixed wear-and-tear exceeds consumed constant capital")
	ErrWearExceedsAdvanced      = errors.New("profit: fixed wear-and-tear exceeds advanced fixed capital")
)

// CostPriceID identifies a stored cost-price computation. 96-bit hex from
// crypto/rand, like every other domain ID in the simulator.
type CostPriceID string

// NewCostPriceID returns a fresh random CostPriceID.
func NewCostPriceID() CostPriceID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("profit: rand: %w", err))
	}
	return CostPriceID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CostPriceID.
func (id CostPriceID) IsZero() bool { return id == "" }

// CapitalOutlay is the capital a commodity's production consumes, decomposed so
// that the cost-price can separate the fixed and circulating components.
//
// Constant is the total constant capital consumed this period (c). Of it,
// FixedWearAndTear is the depreciation of the instruments of labour — the only
// part of fixed capital that enters the cost-price — and the remainder
// (Constant − FixedWearAndTear) is circulating constant capital (raw and
// auxiliary materials). FixedAdvanced is the full value of the instruments of
// labour; it is informational, recording that only its wear-and-tear, not its
// whole value, enters the product (£20 of £1,200 in Marx's example).
type CapitalOutlay struct {
	Constant         ConstantCapital `json:"constant"`
	Variable         VariableCapital `json:"variable"`
	FixedWearAndTear ConstantCapital `json:"fixed_wear_and_tear"`
	FixedAdvanced    ConstantCapital `json:"fixed_advanced"`
}

// Validate enforces the magnitude invariants of a CapitalOutlay.
func (o CapitalOutlay) Validate() error {
	if o.Constant < 0 || o.Variable < 0 || o.FixedWearAndTear < 0 || o.FixedAdvanced < 0 {
		return ErrNegativeMagnitude
	}
	if o.FixedWearAndTear > o.Constant {
		// Fixed depreciation is one component of the consumed constant capital,
		// so it cannot exceed it.
		return ErrFixedWearExceedsConstant
	}
	if o.FixedAdvanced > 0 && o.FixedWearAndTear > o.FixedAdvanced {
		// You cannot wear away more value in a period than the fixed capital
		// advanced holds.
		return ErrWearExceedsAdvanced
	}
	return nil
}

// CirculatingConstant is the constant capital that enters the cost-price in
// full this period — raw and auxiliary materials (c minus fixed wear-and-tear).
func (o CapitalOutlay) CirculatingConstant() ConstantCapital {
	return o.Constant - o.FixedWearAndTear
}

// CostPrice (k) is the portion of a commodity's value that merely replaces the
// capital consumed in its production: k = c + v. It is what the commodity costs
// the capitalist, as distinct from what it costs in labour (its value).
//
// FixedComponent is the part of k contributed by fixed-capital depreciation;
// CirculatingComponent is the rest (circulating constant capital plus variable
// capital). Their sum is K.
type CostPrice struct {
	ID                   CostPriceID   `json:"id"`
	Outlay               CapitalOutlay `json:"outlay"`
	K                    LabourMinutes `json:"k"`
	FixedComponent       LabourMinutes `json:"fixed_component"`
	CirculatingComponent LabourMinutes `json:"circulating_component"`
	CreatedAt            time.Time     `json:"created_at"`
}

// ComputeCostPrice derives the cost-price k = c + v from a capital outlay,
// splitting it into the fixed (depreciation) and circulating components. It is
// a pure function: it does not assign an ID or timestamp — the store does that
// on persistence.
func ComputeCostPrice(outlay CapitalOutlay) CostPrice {
	k := LabourMinutes(outlay.Constant) + LabourMinutes(outlay.Variable)
	fixed := LabourMinutes(outlay.FixedWearAndTear)
	return CostPrice{
		Outlay:               outlay,
		K:                    k,
		FixedComponent:       fixed,
		CirculatingComponent: k - fixed,
	}
}

// CommodityValueV3 is the Vol. III framing of a commodity's value:
// C = c + v + s, regrouped as C = k + s. Total is C, CostPrice is k, and
// SurplusValue is s.
type CommodityValueV3 struct {
	CostPrice    LabourMinutes `json:"cost_price"`
	SurplusValue SurplusValue  `json:"surplus_value"`
	Total        LabourMinutes `json:"total"`
}

// ComputeCommodityValue assembles C = k + s. Because the cost-price excludes
// the surplus-value, C exceeds k whenever s > 0 — which is always the case in
// capitalist production.
func ComputeCommodityValue(k LabourMinutes, s SurplusValue) CommodityValueV3 {
	return CommodityValueV3{
		CostPrice:    k,
		SurplusValue: s,
		Total:        k + LabourMinutes(s),
	}
}
