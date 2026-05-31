// Package rent — Ch. 46 types: Building Site Rent, Rent in Mining, Price of Land.
// Differential rent applies wherever a monopolisable natural force exists.
// Building-site rent is location-driven; mining rent follows agricultural DR;
// monopoly-price rent is the excess of sale price over value. The price of land
// is capitalised rent: annual rent ÷ interest rate.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// LandPriceScenarioKind enumerates the conditions under which the price of land rises.
type LandPriceScenarioKind int

const (
	LandPriceFallingInterestRate  LandPriceScenarioKind = 1 // price rises, rent unchanged
	LandPriceCapitalImprovements  LandPriceScenarioKind = 2 // price rises, rent unchanged
	LandPriceRentRisesWithProduct LandPriceScenarioKind = 3 // product price rises → rent rises
	LandPriceRentRisesMoreCapital LandPriceScenarioKind = 4 // more capital → more rent mass
	LandPriceRentRisesPriceFalls  LandPriceScenarioKind = 5 // product price falls but DR rises
)

// IsValid reports whether k is a defined LandPriceScenarioKind.
func (k LandPriceScenarioKind) IsValid() bool {
	return k >= LandPriceFallingInterestRate && k <= LandPriceRentRisesPriceFalls
}

// BuildingSiteRentID identifies a building-site-rent record.
type BuildingSiteRentID string

// NewBuildingSiteRentID returns a fresh 96-bit hex ID.
func NewBuildingSiteRentID() BuildingSiteRentID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return BuildingSiteRentID(hex.EncodeToString(b[:]))
}

// BuildingSiteRent is location-driven differential rent (no fertility component).
type BuildingSiteRent struct {
	ID                      BuildingSiteRentID `json:"id"`
	ParcelID                string             `json:"parcel_id"`
	LocationGrade           int                `json:"location_grade"`
	AnnualRentLabourMinutes int64              `json:"annual_rent_labour_minutes"`
	LeaseTermYears          int                `json:"lease_term_years"`
	CreatedAt               time.Time          `json:"created_at"`
}

// MiningRentID identifies a mining-rent record.
type MiningRentID string

// NewMiningRentID returns a fresh 96-bit hex ID.
func NewMiningRentID() MiningRentID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return MiningRentID(hex.EncodeToString(b[:]))
}

// MiningRent is differential rent from a mine's quality/richness gradient.
// The worst (marginal) mine yields zero rent — the same law as worst soil.
type MiningRent struct {
	ID                      MiningRentID `json:"id"`
	ParcelID                string       `json:"parcel_id"`
	OreGrade                int          `json:"ore_grade"`
	AnnualRentLabourMinutes int64        `json:"annual_rent_labour_minutes"`
	SurplusProfitBP         int64        `json:"surplus_profit_bp"`
	CreatedAt               time.Time    `json:"created_at"`
}

// MonopolyPriceRentID identifies a monopoly-price-rent record.
type MonopolyPriceRentID string

// NewMonopolyPriceRentID returns a fresh 96-bit hex ID.
func NewMonopolyPriceRentID() MonopolyPriceRentID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return MonopolyPriceRentID(hex.EncodeToString(b[:]))
}

// MonopolyPriceRent is the excess of monopoly sale price over value
// (buyer willingness to pay determines price). MonopolyRent = sale − value (>=0).
type MonopolyPriceRent struct {
	ID                             MonopolyPriceRentID `json:"id"`
	ParcelID                       string              `json:"parcel_id"`
	ProductKind                    string              `json:"product_kind"`
	ValueLabourMinutes             int64               `json:"value_labour_minutes"`
	MonopolySalePriceLabourMinutes int64               `json:"monopoly_sale_price_labour_minutes"`
	MonopolyRentLabourMinutes      int64               `json:"monopoly_rent_labour_minutes"`
	CreatedAt                      time.Time           `json:"created_at"`
}

// LandPriceScenarioID identifies a land-price-scenario record.
type LandPriceScenarioID string

// NewLandPriceScenarioID returns a fresh 96-bit hex ID.
func NewLandPriceScenarioID() LandPriceScenarioID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return LandPriceScenarioID(hex.EncodeToString(b[:]))
}

// LandPriceScenario records conditions under which the price of land rises.
// NewLandPriceBP = CapitalisedLandPrice(NewAnnualRentBP, NewInterestRateBP).
type LandPriceScenario struct {
	ID                LandPriceScenarioID   `json:"id"`
	Kind              LandPriceScenarioKind `json:"kind"`
	ParcelID          string                `json:"parcel_id"`
	OldAnnualRentBP   int64                 `json:"old_annual_rent_bp"`
	NewAnnualRentBP   int64                 `json:"new_annual_rent_bp"`
	OldInterestRateBP int64                 `json:"old_interest_rate_bp"`
	NewInterestRateBP int64                 `json:"new_interest_rate_bp"`
	OldLandPriceBP    int64                 `json:"old_land_price_bp"`
	NewLandPriceBP    int64                 `json:"new_land_price_bp"`
	CreatedAt         time.Time             `json:"created_at"`
}

// ComputeMonopolyRent returns max(0, salePriceLM − valueLM).
func ComputeMonopolyRent(salePriceLM, valueLM int64) int64 {
	if salePriceLM > valueLM {
		return salePriceLM - valueLM
	}
	return 0
}

// CapitalisedLandPrice returns annualRent × 10000 ÷ interestRateBP (round-half-up);
// 0 when interestRateBP <= 0. E.g. rent 10, 250 bp (2.5%) → 400.
func CapitalisedLandPrice(annualRent, interestRateBP int64) int64 {
	if interestRateBP <= 0 {
		return 0
	}
	return roundHalfUp(annualRent*10000, interestRateBP)
}
