// Package rent — Ch. 42 types: Differential Rent II, Second Case (Falling
// Price of Production). When additional capital in superior soils excludes the
// worst soil (A) from cultivation, a better soil becomes the new regulator and
// the price of production FALLS. The compensation law: doubled grain-rent can
// offset the lower price so total money-rent is preserved.
//
// Unit scheme: CapitalBP in pence; OutputQuarters in 1/100 qtr; for THIS chapter
// SuccessiveInvestment.SurplusProfitBP carries grain-rent in 1/100-qtr units (not
// pence); NewPriceOfProductionBP is pence per WHOLE quarter (30s→360, 45s→540).
// TotalMoneyRentalBP = roundHalfUp(TotalGrainRentQtrs × NewPriceOfProductionBP, 100).
//
// Table IV (variant 1, 540p): grain 250 → roundHalfUp(250×540,100)=1350 (£13.50).
// Table IVd (variant 3, 360p): grain 500 → roundHalfUp(500×360,100)=1800 (£18).
// Compensation: IVd money-rent (1800) equals the original Table I rental despite
// the halved price — doubled grain-rent compensates exactly.
//
// Reuses SoilGrade, SuccessiveInvestment, roundHalfUp from this package.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// ErrInvalidExclusion is returned when NewRegulatingGrade <= ExcludedGrade.
var ErrInvalidExclusion = errors.New("rent: new_regulating_grade must be > excluded_grade")

// FallingPriceVariant enumerates the three DR II falling-price sub-cases.
type FallingPriceVariant int

const (
	FallingPriceConstantProductivity   FallingPriceVariant = 1
	FallingPriceDecreasingProductivity FallingPriceVariant = 2
	FallingPriceRisingProductivity     FallingPriceVariant = 3
)

// IsValid reports whether v is a defined FallingPriceVariant.
func (v FallingPriceVariant) IsValid() bool {
	return v >= FallingPriceConstantProductivity && v <= FallingPriceRisingProductivity
}

// SoilExclusionEventID identifies a soil-exclusion event.
type SoilExclusionEventID string

// NewSoilExclusionEventID returns a fresh 96-bit hex ID.
func NewSoilExclusionEventID() SoilExclusionEventID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return SoilExclusionEventID(hex.EncodeToString(b[:]))
}

// SoilExclusionEvent records a soil grade driven from cultivation; the
// NewRegulatingGrade becomes the new worst-cultivated soil at NewPriceOfProductionBP
// (pence per whole quarter).
type SoilExclusionEvent struct {
	ID                     SoilExclusionEventID `json:"id"`
	ExcludedGrade          SoilGrade            `json:"excluded_grade"`
	NewRegulatingGrade     SoilGrade            `json:"new_regulating_grade"`
	NewPriceOfProductionBP int64                `json:"new_price_of_production_bp"`
	CreatedAt              time.Time            `json:"created_at"`
}

// Validate enforces the Ch. 42 invariant: the new regulating grade must be
// strictly superior (higher number) to the excluded grade.
func (e SoilExclusionEvent) Validate() error {
	if e.NewRegulatingGrade <= e.ExcludedGrade {
		return ErrInvalidExclusion
	}
	return nil
}

// FallingPriceDR2OutcomeID identifies a falling-price DR II outcome.
type FallingPriceDR2OutcomeID string

// NewFallingPriceDR2OutcomeID returns a fresh 96-bit hex ID.
func NewFallingPriceDR2OutcomeID() FallingPriceDR2OutcomeID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return FallingPriceDR2OutcomeID(hex.EncodeToString(b[:]))
}

// FallingPriceDR2Outcome is the computed falling-price result: total money
// rental valued at the new (lower) regulating price.
type FallingPriceDR2Outcome struct {
	ID                  FallingPriceDR2OutcomeID `json:"id"`
	Variant             FallingPriceVariant      `json:"variant"`
	ExclusionEventID    string                   `json:"exclusion_event_id"`
	TotalMoneyRentalBP  int64                    `json:"total_money_rental_bp"`
	TotalGrainRentQtrs  int64                    `json:"total_grain_rent_qtrs"`
	TotalOutputQuarters int64                    `json:"total_output_quarters"`
	TotalCapitalBP      int64                    `json:"total_capital_bp"`
	CreatedAt           time.Time                `json:"created_at"`
}

// NormalCapitalPerAcreID identifies a normal-capital-per-acre record.
type NormalCapitalPerAcreID string

// NewNormalCapitalPerAcreID returns a fresh 96-bit hex ID.
func NewNormalCapitalPerAcreID() NormalCapitalPerAcreID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return NormalCapitalPerAcreID(hex.EncodeToString(b[:]))
}

// NormalCapitalPerAcre records the minimum capital required per acre in a
// historical period (Marx: pre-1848 England £8/acre = 800p; post-1848 £12 = 1200p).
type NormalCapitalPerAcre struct {
	ID               NormalCapitalPerAcreID `json:"id"`
	PeriodLabel      string                 `json:"period_label"`
	CapitalPerAcreBP int64                  `json:"capital_per_acre_bp"`
	CreatedAt        time.Time              `json:"created_at"`
}

// ComputeFallingPriceDR2 aggregates successive investments after a soil-exclusion
// event, valuing grain-rent at the new (fallen) regulating price.
// TotalGrainRentQtrs = Σ SurplusProfitBP (1/100-qtr units);
// TotalMoneyRentalBP = roundHalfUp(TotalGrainRentQtrs × exclusion.NewPriceOfProductionBP, 100).
// Empty investments → zero outcome with Variant + ExclusionEventID set.
func ComputeFallingPriceDR2(variant FallingPriceVariant, exclusion SoilExclusionEvent, investments []SuccessiveInvestment) FallingPriceDR2Outcome {
	out := FallingPriceDR2Outcome{
		Variant:          variant,
		ExclusionEventID: string(exclusion.ID),
	}
	for _, inv := range investments {
		out.TotalCapitalBP += inv.CapitalBP
		out.TotalOutputQuarters += inv.OutputQuarters
		out.TotalGrainRentQtrs += inv.SurplusProfitBP
	}
	if exclusion.NewPriceOfProductionBP > 0 {
		out.TotalMoneyRentalBP = roundHalfUp(out.TotalGrainRentQtrs*exclusion.NewPriceOfProductionBP, 100)
	}
	return out
}
