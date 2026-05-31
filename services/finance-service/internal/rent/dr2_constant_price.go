// Package rent — Ch. 41 types: Differential Rent II, First Case (Constant
// Price of Production). When soil A keeps regulating price, additional capital
// on better soils B-D always RAISES absolute rent per acre, whether productivity
// of added capital is proportional (Case II), decreasing (Case III), or
// increasing (Case IV). Case I (only average profit, no surplus) forms no new rent.
//
// Marx Table II (Case II, Soil D, two £2.50 tranches each surplus £9):
//
//	TotalCapitalBP=500, TotalRentalBP=1800, AdditionalCapitalBP=250,
//	NewRentPerAcreBP=1800, RateOfSurplusProfitBP=roundHalfUp(1800*10000,500)=36000.
//
// Units: *BP in pence. Reuses SuccessiveInvestment + roundHalfUp from this package.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// CapitalProductivityCase enumerates the four DR II constant-price sub-cases.
type CapitalProductivityCase int

const (
	CaseAverageOnly      CapitalProductivityCase = 1 // only average profit → no new rent
	CaseProportional     CapitalProductivityCase = 2 // output ∝ capital
	CaseDecreasingRate   CapitalProductivityCase = 3 // surplus-profit rate falls, absolute rises
	CaseIncreasingReturn CapitalProductivityCase = 4 // improvements → output rises faster than capital
)

// IsValid reports whether c is a defined CapitalProductivityCase.
func (c CapitalProductivityCase) IsValid() bool {
	return c >= CaseAverageOnly && c <= CaseIncreasingReturn
}

// DifferentialRentIIOutcomeID identifies a DR II constant-price outcome.
type DifferentialRentIIOutcomeID string

// NewDifferentialRentIIOutcomeID returns a fresh 96-bit hex ID.
func NewDifferentialRentIIOutcomeID() DifferentialRentIIOutcomeID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return DifferentialRentIIOutcomeID(hex.EncodeToString(b[:]))
}

// DifferentialRentIIOutcome is the computed result of applying a
// CapitalProductivityCase to successive investments on one parcel (pence units).
type DifferentialRentIIOutcome struct {
	ID                    DifferentialRentIIOutcomeID `json:"id"`
	ParcelID              string                      `json:"parcel_id"`
	Case                  CapitalProductivityCase     `json:"case"`
	AdditionalCapitalBP   int64                       `json:"additional_capital_bp"`
	NewRentPerAcreBP      int64                       `json:"new_rent_per_acre_bp"`
	TotalRentalBP         int64                       `json:"total_rental_bp"`
	RateOfSurplusProfitBP int64                       `json:"rate_of_surplus_profit_bp"`
	CreatedAt             time.Time                   `json:"created_at"`
}

// IntensiveExtensiveComparisonID identifies an intensive-vs-extensive comparison.
type IntensiveExtensiveComparisonID string

// NewIntensiveExtensiveComparisonID returns a fresh 96-bit hex ID.
func NewIntensiveExtensiveComparisonID() IntensiveExtensiveComparisonID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return IntensiveExtensiveComparisonID(hex.EncodeToString(b[:]))
}

// IntensiveExtensiveComparison shows that concentrating capital on a smaller
// acreage yields higher rent per acre than dispersing equal capital over more
// acres at the same total rental (pence units).
type IntensiveExtensiveComparison struct {
	ID                     IntensiveExtensiveComparisonID `json:"id"`
	IntensiveRentPerAcreBP int64                          `json:"intensive_rent_per_acre_bp"`
	ExtensiveRentPerAcreBP int64                          `json:"extensive_rent_per_acre_bp"`
	TotalCapitalBP         int64                          `json:"total_capital_bp"`
	TotalRentalSameBP      int64                          `json:"total_rental_same_bp"`
	CreatedAt              time.Time                      `json:"created_at"`
}

// ComputeDR2Outcome aggregates successive investments under a CapitalProductivityCase.
// TotalRentalBP = Σ SurplusProfitBP; AdditionalCapitalBP = Σ CapitalBP of tranches
// with TrancheNumber>=2; NewRentPerAcreBP = TotalRentalBP (1-acre normalisation);
// RateOfSurplusProfitBP = roundHalfUp(TotalRentalBP*10000, ΣCapitalBP) when ΣCapitalBP>0
// else 0. ParcelID from the first investment. Empty input → zero outcome with Case set.
func ComputeDR2Outcome(c CapitalProductivityCase, investments []SuccessiveInvestment) DifferentialRentIIOutcome {
	out := DifferentialRentIIOutcome{Case: c}
	if len(investments) == 0 {
		return out
	}
	out.ParcelID = investments[0].ParcelID
	var totalCapital, totalRental, additionalCapital int64
	for _, inv := range investments {
		totalCapital += inv.CapitalBP
		totalRental += inv.SurplusProfitBP
		if inv.TrancheNumber >= 2 {
			additionalCapital += inv.CapitalBP
		}
	}
	out.TotalRentalBP = totalRental
	out.AdditionalCapitalBP = additionalCapital
	out.NewRentPerAcreBP = totalRental
	if totalCapital > 0 {
		out.RateOfSurplusProfitBP = roundHalfUp(totalRental*10000, totalCapital)
	}
	return out
}
