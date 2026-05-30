package avgprofit

// This file carries Ch. 10's argument: within a branch of production the
// market value is set by the production conditions of the BULK of output —
// not the best or worst individual firm. A firm that produces below the
// market value (higher individual productivity) sells at the market value
// and pockets the difference as surplus-profit. The amount is signed:
// positive when individual value < market value (extra profit); negative
// when individual value > market value (below-average firm loses out).

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// MarketPriceAmount is a price in circulation (labour-minutes denomination).
type MarketPriceAmount LabourMinutes

// IndividualValue is one firm's own production cost (labour-minutes).
type IndividualValue LabourMinutes

// MarketValueID identifies a persisted market-value record.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type MarketValueID string

// NewMarketValueID returns a fresh random MarketValueID.
func NewMarketValueID() MarketValueID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("avgprofit: rand: %w", err))
	}
	return MarketValueID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty MarketValueID.
func (id MarketValueID) IsZero() bool { return id == "" }

// MarketValue: value set by the production conditions of the BULK of output.
// BestConditionValue is the lowest-cost (most productive) firm;
// WorstConditionValue is the highest-cost (least productive) firm.
// Value == BulkConditionValue (the regulating value for the branch).
type MarketValue struct {
	ID                  MarketValueID  `json:"id"`
	SphereName          string         `json:"sphere_name"`
	BulkConditionValue  CommodityValue `json:"bulk_condition_value"`
	BestConditionValue  CommodityValue `json:"best_condition_value"`
	WorstConditionValue CommodityValue `json:"worst_condition_value"`
	Value               CommodityValue `json:"value"` // == BulkConditionValue
	CreatedAt           time.Time      `json:"created_at"`
}

// Validate returns an error if any condition value is negative.
func (m MarketValue) Validate() error {
	if m.BulkConditionValue < 0 || m.BestConditionValue < 0 || m.WorstConditionValue < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// ComputeMarketValue builds a MarketValue with Value = bulk condition value.
func ComputeMarketValue(name string, bulk, best, worst CommodityValue) MarketValue {
	return MarketValue{
		SphereName:          name,
		BulkConditionValue:  bulk,
		BestConditionValue:  best,
		WorstConditionValue: worst,
		Value:               bulk,
	}
}

// MarketPrice is the actual circulation price versus the market value.
// It is a value object and is not persisted separately.
type MarketPrice struct {
	Amount      MarketPriceAmount `json:"amount"`
	MarketValue CommodityValue    `json:"market_value"`
	Deviation   ValueDeviation    `json:"deviation"` // amount − market value
}

// ComputeMarketPrice derives the signed deviation of a circulation price
// from the branch market value: Deviation = amount − marketValue.
func ComputeMarketPrice(amount MarketPriceAmount, marketValue CommodityValue) MarketPrice {
	return MarketPrice{
		Amount:      amount,
		MarketValue: marketValue,
		Deviation:   ValueDeviation(int64(amount) - int64(marketValue)),
	}
}

// SurplusProfitID identifies a persisted surplus-profit record.
// 96-bit hex from crypto/rand.
type SurplusProfitID string

// NewSurplusProfitID returns a fresh random SurplusProfitID.
func NewSurplusProfitID() SurplusProfitID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("avgprofit: rand: %w", err))
	}
	return SurplusProfitID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty SurplusProfitID.
func (id SurplusProfitID) IsZero() bool { return id == "" }

// SurplusProfit is the excess a below-market-value firm earns when it sells
// at the market value. Amount is signed: positive when individual value <
// market value (firm pockets extra); negative when individual value > market
// value (firm earns less than the average).
//
//	Amount = (MarketValue − IndividualValue) × OutputQty
type SurplusProfit struct {
	ID              SurplusProfitID  `json:"id"`
	FirmName        string           `json:"firm_name"`
	IndividualValue IndividualValue  `json:"individual_value"`
	MarketValue     CommodityValue   `json:"market_value"`
	OutputQty       int64            `json:"output_qty"`
	GeneralRate     SphereProfitRate `json:"general_rate"`
	Amount          Profit           `json:"amount"` // signed
	CreatedAt       time.Time        `json:"created_at"`
}

// Validate returns ErrNegativeMagnitude if any input magnitude is negative.
func (s SurplusProfit) Validate() error {
	if s.IndividualValue < 0 || s.MarketValue < 0 || s.OutputQty < 0 || s.GeneralRate < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// ComputeSurplusProfit computes the signed surplus-profit for a firm:
//
//	Amount = (marketVal − indVal) × outputQty
func ComputeSurplusProfit(firm string, indVal IndividualValue, marketVal CommodityValue, outputQty int64, generalRate SphereProfitRate) SurplusProfit {
	amount := Profit((int64(marketVal) - int64(indVal)) * outputQty)
	return SurplusProfit{
		FirmName:        firm,
		IndividualValue: indVal,
		MarketValue:     marketVal,
		OutputQty:       outputQty,
		GeneralRate:     generalRate,
		Amount:          amount,
	}
}
