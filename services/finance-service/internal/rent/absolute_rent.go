// Package rent — Ch. 45 types: Absolute Ground-Rent. Even the worst soil yields
// rent because agriculture has a LOWER organic composition of capital than the
// social average, so its products sell above the price of production (near value),
// generating surplus-value above average profit; landed property captures it.
// Units: composition/rate in bp (10000=100%); value/price in LabourMinutes.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type AgriculturalCompositionID string

func NewAgriculturalCompositionID() AgriculturalCompositionID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return AgriculturalCompositionID(hex.EncodeToString(b[:]))
}

// AgriculturalComposition records agricultural vs social-average organic composition.
type AgriculturalComposition struct {
	ID                   AgriculturalCompositionID `json:"id"`
	ConstantCapitalBP    int64                     `json:"constant_capital_bp"`
	VariableCapitalBP    int64                     `json:"variable_capital_bp"`
	CompositionRatioBP   int64                     `json:"composition_ratio_bp"`
	SocialAverageRatioBP int64                     `json:"social_average_ratio_bp"`
	IsBelowAverage       bool                      `json:"is_below_average"`
	CreatedAt            time.Time                 `json:"created_at"`
}

type ValuePriceGapID string

func NewValuePriceGapID() ValuePriceGapID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return ValuePriceGapID(hex.EncodeToString(b[:]))
}

// ValuePriceGap records value minus price of production. Gap not clamped.
type ValuePriceGap struct {
	ID                             ValuePriceGapID `json:"id"`
	ProductID                      string          `json:"product_id"`
	ValueLabourMinutes             int64           `json:"value_labour_minutes"`
	PriceOfProductionLabourMinutes int64           `json:"price_of_production_labour_minutes"`
	GapLabourMinutes               int64           `json:"gap_labour_minutes"`
	CreatedAt                      time.Time       `json:"created_at"`
}

type AbsoluteRentID string

func NewAbsoluteRentID() AbsoluteRentID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return AbsoluteRentID(hex.EncodeToString(b[:]))
}

// AbsoluteRent records the absolute rent on a parcel. AbsoluteRentBP=max(0,S−P).
type AbsoluteRent struct {
	ID              AbsoluteRentID `json:"id"`
	ParcelID        string         `json:"parcel_id"`
	ValuePriceGapID string         `json:"value_price_gap_id"`
	SurplusValueBP  int64          `json:"surplus_value_bp"`
	AverageProfitBP int64          `json:"average_profit_bp"`
	AbsoluteRentBP  int64          `json:"absolute_rent_bp"`
	CreatedAt       time.Time      `json:"created_at"`
}

type AbsoluteRentLimitID string

func NewAbsoluteRentLimitID() AbsoluteRentLimitID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return AbsoluteRentLimitID(hex.EncodeToString(b[:]))
}

// AbsoluteRentLimit bounds absolute rent by the value–price gap. IsAtLimit=actual>=max.
type AbsoluteRentLimit struct {
	ID           AbsoluteRentLimitID `json:"id"`
	MaxRentBP    int64               `json:"max_rent_bp"`
	ActualRentBP int64               `json:"actual_rent_bp"`
	IsAtLimit    bool                `json:"is_at_limit"`
	CreatedAt    time.Time           `json:"created_at"`
}

// ComputeAbsoluteRentBP returns max(0, surplusValueBP − averageProfitBP).
func ComputeAbsoluteRentBP(surplusValueBP, averageProfitBP int64) int64 {
	if surplusValueBP > averageProfitBP {
		return surplusValueBP - averageProfitBP
	}
	return 0
}

// ComputeCompositionRatioBP returns C/V × 10000 (round-half-up); 0 when V==0.
func ComputeCompositionRatioBP(constantCapitalBP, variableCapitalBP int64) int64 {
	if variableCapitalBP == 0 {
		return 0
	}
	return roundHalfUp(constantCapitalBP*10000, variableCapitalBP)
}

// ComputeValuePriceGap returns valueLabourMinutes − priceOfProductionLabourMinutes.
func ComputeValuePriceGap(valueLabourMinutes, priceOfProductionLabourMinutes int64) int64 {
	return valueLabourMinutes - priceOfProductionLabourMinutes
}
