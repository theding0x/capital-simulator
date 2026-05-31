// Package rent — Ch. 38 types: Differential Rent, General Remarks.
// Extends the rent package opened in Ch. 37.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// roundHalfUp computes (numer + denom/2) / denom — canonical round-half-up
// division for basis-point arithmetic. Private to this package (same body as
// credit.roundHalfUp and merchant.roundHalfUp).
func roundHalfUp(numer, denom int64) int64 {
	return (numer + denom/2) / denom
}

// PriceOfProductionSurplusProfitID identifies a persisted PoP surplus-profit record.
type PriceOfProductionSurplusProfitID string

// NewPoPSurplusProfitID returns a fresh random PriceOfProductionSurplusProfitID.
func NewPoPSurplusProfitID() PriceOfProductionSurplusProfitID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return PriceOfProductionSurplusProfitID(hex.EncodeToString(b[:]))
}

// PriceOfProductionSurplusProfit records the gap between the market-regulating
// (general) price of production and the individual (waterfall-aided) price of
// production for one capital; the positive gap is the surplus-profit landed
// property converts into differential rent. Prices are in pence (1/100 £).
// SurplusProfitBP is derived server-side via ComputeSurplusProfit.
type PriceOfProductionSurplusProfit struct {
	ID                            PriceOfProductionSurplusProfitID `json:"id"`
	CapitalID                     string                           `json:"capital_id"`
	GeneralPriceOfProductionBP    int64                            `json:"general_price_of_production_bp"`
	IndividualPriceOfProductionBP int64                            `json:"individual_price_of_production_bp"`
	SurplusProfitBP               int64                            `json:"surplus_profit_bp"`
	CreatedAt                     time.Time                        `json:"created_at"`
}

// ComputeSurplusProfit returns the surplus-profit (general − individual) in the
// same unit as its inputs. May be zero or negative; not clamped.
func ComputeSurplusProfit(generalBP, individualBP int64) int64 {
	return generalBP - individualBP
}

// MonopolisedNaturalForceID identifies a persisted monopolised-natural-force record.
type MonopolisedNaturalForceID string

// NewMonopolisedNaturalForceID returns a fresh random MonopolisedNaturalForceID.
func NewMonopolisedNaturalForceID() MonopolisedNaturalForceID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return MonopolisedNaturalForceID(hex.EncodeToString(b[:]))
}

// MonopolisedNaturalForce is a natural productive force (waterfall, mine,
// fertile soil) whose monopolisation by landed property enables differential
// rent. IsMonopolisable is true when private ownership is legally possible.
type MonopolisedNaturalForce struct {
	ID              MonopolisedNaturalForceID `json:"id"`
	Kind            string                    `json:"kind"`
	IsMonopolisable bool                      `json:"is_monopolisable"`
	ParcelID        string                    `json:"parcel_id"`
	CreatedAt       time.Time                 `json:"created_at"`
}

// DifferentialRentID identifies a persisted differential-rent record.
type DifferentialRentID string

// NewDifferentialRentID returns a fresh random DifferentialRentID.
func NewDifferentialRentID() DifferentialRentID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return DifferentialRentID(hex.EncodeToString(b[:]))
}

// DifferentialRent is ground-rent arising from converting surplus-profit into
// rent. AnnualRentLabourMinutes is the captured magnitude in LabourMinutes
// (£1 = 480 LM).
type DifferentialRent struct {
	ID                      DifferentialRentID `json:"id"`
	ParcelID                string             `json:"parcel_id"`
	SurplusProfitBP         int64              `json:"surplus_profit_bp"`
	AnnualRentLabourMinutes int64              `json:"annual_rent_labour_minutes"`
	CreatedAt               time.Time          `json:"created_at"`
}

// CapitalisedRentPriceID identifies a persisted capitalised-rent-price record.
type CapitalisedRentPriceID string

// NewCapitalisedRentPriceID returns a fresh random CapitalisedRentPriceID.
func NewCapitalisedRentPriceID() CapitalisedRentPriceID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return CapitalisedRentPriceID(hex.EncodeToString(b[:]))
}

// CapitalisedRentPrice is the purchase price of land: annual rent capitalised at
// the interest rate. CapitalisedPriceLabourMinutes is derived via
// ComputeCapitalisedPrice.
type CapitalisedRentPrice struct {
	ID                            CapitalisedRentPriceID `json:"id"`
	ParcelID                      string                 `json:"parcel_id"`
	AnnualRentLabourMinutes       int64                  `json:"annual_rent_labour_minutes"`
	InterestRateBP                int64                  `json:"interest_rate_bp"`
	CapitalisedPriceLabourMinutes int64                  `json:"capitalised_price_labour_minutes"`
	CreatedAt                     time.Time              `json:"created_at"`
}

// ComputeCapitalisedPrice returns the capitalised land price in LabourMinutes:
// roundHalfUp(annualRent × 10000, interestRateBP). Returns 0 if interestRateBP <= 0.
func ComputeCapitalisedPrice(annualRentLabourMinutes, interestRateBP int64) int64 {
	if interestRateBP <= 0 {
		return 0
	}
	return roundHalfUp(annualRentLabourMinutes*10000, interestRateBP)
}
