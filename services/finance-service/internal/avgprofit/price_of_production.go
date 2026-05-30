package avgprofit

// This file carries the price-of-production computation for Ch. 9. Once the
// general rate p̄′ is known, each capital's price of production is:
//
//   price = k + k·p̄′/10000   (k = cost-price = c + v consumed)
//
// Deviation = price − value: positive for higher-composition branches
// (receivers of transferred surplus-value), negative for lower-composition
// branches (donors). The sign follows the FIXTURES, not the spec prose.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// PriceOfProductionID identifies a persisted price-of-production record.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type PriceOfProductionID string

// NewPriceOfProductionID returns a fresh random PriceOfProductionID.
func NewPriceOfProductionID() PriceOfProductionID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("avgprofit: rand: %w", err))
	}
	return PriceOfProductionID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty PriceOfProductionID.
func (id PriceOfProductionID) IsZero() bool { return id == "" }

// PriceOfProduction is the persisted price-of-production entity. Inputs are
// SphereName, CostPrice k, GeneralRate p̄′, and CommodityValue (0 = unknown).
// Derived fields — AverageProfit, Price, Deviation, Composition — are filled
// by ComputePriceOfProduction.
type PriceOfProduction struct {
	ID             PriceOfProductionID `json:"id"`
	SphereName     string              `json:"sphere_name"`
	CostPrice      CostPrice           `json:"cost_price"`
	GeneralRate    SphereProfitRate    `json:"general_rate"`
	CommodityValue CommodityValue      `json:"commodity_value"` // input; 0 = unknown
	AverageProfit  Profit              `json:"average_profit"`  // k*rate/10000
	Price          ProductionPrice     `json:"price"`           // k + average profit
	Deviation      ValueDeviation      `json:"deviation"`       // price − value
	Composition    CompositionKind     `json:"composition"`     // from deviation sign; "" if value<=0
	CreatedAt      time.Time           `json:"created_at"`
}

// Validate returns an error if any input magnitude is negative.
func (p PriceOfProduction) Validate() error {
	if p.CostPrice < 0 || p.GeneralRate < 0 || p.CommodityValue < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// ComputePriceOfProduction derives AverageProfit, Price, Deviation, and
// Composition from the supplied inputs:
//
//	avgProfit   = Profit(int64(k) * int64(rate) / 10000)
//	price       = ProductionPrice(int64(k) + int64(avgProfit))
//	deviation   = ValueDeviation(int64(price) - int64(value))
//	composition: if value<=0 -> ""; dev>0->KindHigher, dev<0->KindLower, ==0->KindAverage
func ComputePriceOfProduction(name string, k CostPrice, rate SphereProfitRate, value CommodityValue) PriceOfProduction {
	avgProfit := Profit(int64(k) * int64(rate) / basisPoints)
	price := ProductionPrice(int64(k) + int64(avgProfit))
	dev := ValueDeviation(int64(price) - int64(value))

	var kind CompositionKind
	if value > 0 {
		switch {
		case dev > 0:
			kind = KindHigher
		case dev < 0:
			kind = KindLower
		default:
			kind = KindAverage
		}
	}

	return PriceOfProduction{
		SphereName:     name,
		CostPrice:      k,
		GeneralRate:    rate,
		CommodityValue: value,
		AverageProfit:  avgProfit,
		Price:          price,
		Deviation:      dev,
		Composition:    kind,
	}
}

// ValuePriceDeviation is a simple helper that records the signed difference
// between a commodity's value and its price of production.
type ValuePriceDeviation struct {
	Value             CommodityValue  `json:"value"`
	PriceOfProduction ProductionPrice `json:"price_of_production"`
	Deviation         ValueDeviation  `json:"deviation"` // = price − value
}

// ComputeValuePriceDeviation returns the signed deviation price − value.
func ComputeValuePriceDeviation(value CommodityValue, price ProductionPrice) ValuePriceDeviation {
	return ValuePriceDeviation{
		Value:             value,
		PriceOfProduction: price,
		Deviation:         ValueDeviation(int64(price) - int64(value)),
	}
}
