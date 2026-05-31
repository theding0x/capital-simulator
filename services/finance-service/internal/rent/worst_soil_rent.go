// Package rent — Ch. 44 types: Differential Rent on the Worst Cultivated Soil
// (hitherto rentless soil A acquires rent by three mechanisms). Also models
// soil-improvement capital and regulating-price shifts.
//
// Price fields use basis-point pence (£1 = 10000 bp): £3.00=30000, £3.50=35000.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// RentEmergenceMechanism enumerates the three routes by which differential rent
// first appears on the hitherto rentless worst soil.
type RentEmergenceMechanism int

const (
	MechanismBetterSoilRegulates    RentEmergenceMechanism = 1 // better soil's under-productive investment sets new higher price
	MechanismIncreasingProductivity RentEmergenceMechanism = 2 // A's average price falls → new lower regulating price
	MechanismDecreasingProductivity RentEmergenceMechanism = 3 // A's average price rises → rent emerges on A
)

// IsValid reports whether m is a defined RentEmergenceMechanism.
func (m RentEmergenceMechanism) IsValid() bool {
	return m >= MechanismBetterSoilRegulates && m <= MechanismDecreasingProductivity
}

// RentOnWorstSoilID identifies a rent-on-worst-soil record.
type RentOnWorstSoilID string

// NewRentOnWorstSoilID returns a fresh 96-bit hex ID.
func NewRentOnWorstSoilID() RentOnWorstSoilID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return RentOnWorstSoilID(hex.EncodeToString(b[:]))
}

// RentOnWorstSoil records differential rent emerging on the rentless worst soil A.
// RentPerAcreBP is computed server-side via ComputeWorstSoilRentPerAcre (never negative).
type RentOnWorstSoil struct {
	ID                     RentOnWorstSoilID      `json:"id"`
	ParcelID               string                 `json:"parcel_id"`
	Mechanism              RentEmergenceMechanism `json:"mechanism"`
	OldPriceOfProductionBP int64                  `json:"old_price_of_production_bp"`
	NewPriceOfProductionBP int64                  `json:"new_price_of_production_bp"`
	RentPerAcreBP          int64                  `json:"rent_per_acre_bp"`
	CreatedAt              time.Time              `json:"created_at"`
}

// SoilImprovementCapitalID identifies a soil-improvement-capital record.
type SoilImprovementCapitalID string

// NewSoilImprovementCapitalID returns a fresh 96-bit hex ID.
func NewSoilImprovementCapitalID() SoilImprovementCapitalID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return SoilImprovementCapitalID(hex.EncodeToString(b[:]))
}

// SoilImprovementCapital records a permanent soil improvement (drainage, levelling)
// that yields a rent component until amortised, then dissolves into pure DR.
type SoilImprovementCapital struct {
	ID                 SoilImprovementCapitalID `json:"id"`
	ParcelID           string                   `json:"parcel_id"`
	CapitalInvestedBP  int64                    `json:"capital_invested_bp"`
	ProductivityGainBP int64                    `json:"productivity_gain_bp"`
	AmortisationYears  int                      `json:"amortisation_years"`
	AmortisedYears     int                      `json:"amortised_years"`
	CreatedAt          time.Time                `json:"created_at"`
}

// IsAmortised reports whether the improvement has been fully amortised.
func (s SoilImprovementCapital) IsAmortised() bool {
	return s.AmortisedYears >= s.AmortisationYears
}

// RegulatingPriceShiftID identifies a regulating-price-shift record.
type RegulatingPriceShiftID string

// NewRegulatingPriceShiftID returns a fresh 96-bit hex ID.
func NewRegulatingPriceShiftID() RegulatingPriceShiftID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return RegulatingPriceShiftID(hex.EncodeToString(b[:]))
}

// RegulatingPriceShift records a movement in which soil grade regulates the
// market price of production (old/new prices in bp).
type RegulatingPriceShift struct {
	ID         RegulatingPriceShiftID `json:"id"`
	FromGrade  SoilGrade              `json:"from_grade"`
	ToGrade    SoilGrade              `json:"to_grade"`
	OldPriceBP int64                  `json:"old_price_bp"`
	NewPriceBP int64                  `json:"new_price_bp"`
	CreatedAt  time.Time              `json:"created_at"`
}

// ComputeWorstSoilRentPerAcre returns the differential rent per acre opened by a
// shift from oldPriceBP to newPriceBP: the positive gap, or 0 if not positive.
//
//	ComputeWorstSoilRentPerAcre(30000, 35000) == 5000 ; (35000, 30000) == 0.
func ComputeWorstSoilRentPerAcre(oldPriceBP, newPriceBP int64) int64 {
	if newPriceBP > oldPriceBP {
		return newPriceBP - oldPriceBP
	}
	return 0
}
