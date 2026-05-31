// Package rent — Ch. 39 types: Differential Rent I (First Form).
// Soils of unequal fertility cultivated side by side with equal capital. The
// worst soil (Grade A) that must be cultivated regulates the market price and
// yields NO rent; better soils yield a surplus-product appropriated as rent.
//
// Unit scheme (monetary amounts in PENCE; quantities in 1/100 quarter):
//
//	CapitalAdvancedBP   — pence (£2.50 → 250)
//	OutputQuarters      — 1/100 qtr (1 qtr → 100)
//	PriceOfProductionBP — pence per quarter (£3/qtr → 300)
//	SurplusProductQtrs  — 1/100 qtr
//	MoneyRentBP         — pence = SurplusProductQtrs × PriceOfProductionBP / 100
//	TotalRentalBP       — pence = Σ MoneyRentBP
//
// Marx Table I (Ch. 39): A out=100 rent=0; B out=200 rent=300; C out=300 rent=600;
// D out=400 rent=900; total rental = 1800 pence = £18.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// SoilGrade enumerates fertility grades A–D (1 = worst, regulates price).
type SoilGrade int

const (
	SoilGradeA SoilGrade = 1
	SoilGradeB SoilGrade = 2
	SoilGradeC SoilGrade = 3
	SoilGradeD SoilGrade = 4
)

// IsValid reports whether s is a defined SoilGrade.
func (s SoilGrade) IsValid() bool { return s >= SoilGradeA && s <= SoilGradeD }

// String returns the letter label for a SoilGrade.
func (s SoilGrade) String() string {
	switch s {
	case SoilGradeA:
		return "A"
	case SoilGradeB:
		return "B"
	case SoilGradeC:
		return "C"
	case SoilGradeD:
		return "D"
	default:
		return "unknown"
	}
}

// DifferentialRentITableID identifies a DR I table.
type DifferentialRentITableID string

// NewDifferentialRentITableID returns a fresh 96-bit hex ID.
func NewDifferentialRentITableID() DifferentialRentITableID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return DifferentialRentITableID(hex.EncodeToString(b[:]))
}

// DifferentialRentIEntryID identifies one soil row within a DR I table.
type DifferentialRentIEntryID string

// NewDifferentialRentIEntryID returns a fresh 96-bit hex ID.
func NewDifferentialRentIEntryID() DifferentialRentIEntryID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return DifferentialRentIEntryID(hex.EncodeToString(b[:]))
}

// DifferentialRentIEntry is one soil row in a DR I table (pence; 1/100-qtr).
// MoneyRentBP = SurplusProductQtrs × PriceOfProductionBP / 100.
type DifferentialRentIEntry struct {
	ID                  DifferentialRentIEntryID `json:"id"`
	TableID             string                   `json:"table_id"`
	Grade               SoilGrade                `json:"grade"`
	CapitalAdvancedBP   int64                    `json:"capital_advanced_bp"`
	OutputQuarters      int64                    `json:"output_quarters"`
	PriceOfProductionBP int64                    `json:"price_of_production_bp"`
	SurplusProductQtrs  int64                    `json:"surplus_product_qtrs"`
	MoneyRentBP         int64                    `json:"money_rent_bp"`
	CreatedAt           time.Time                `json:"created_at"`
}

// DifferentialRentITable is the computed aggregate for a set of DR I soil rows.
type DifferentialRentITable struct {
	ID              DifferentialRentITableID `json:"id"`
	Description     string                   `json:"description"`
	TotalRentalBP   int64                    `json:"total_rental_bp"`
	RegulatingGrade SoilGrade                `json:"regulating_grade"`
	CreatedAt       time.Time                `json:"created_at"`
}

// LocationRentFactorID identifies a location-rent-factor record.
type LocationRentFactorID string

// NewLocationRentFactorID returns a fresh 96-bit hex ID.
func NewLocationRentFactorID() LocationRentFactorID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return LocationRentFactorID(hex.EncodeToString(b[:]))
}

// LocationRentFactor records a transport-cost differential giving rise to a
// location-based component of differential rent (pence).
type LocationRentFactor struct {
	ID               LocationRentFactorID `json:"id"`
	ParcelID         string               `json:"parcel_id"`
	TransportCostBP  int64                `json:"transport_cost_bp"`
	RentEquivalentBP int64                `json:"rent_equivalent_bp"`
	CreatedAt        time.Time            `json:"created_at"`
}

// ComputeDifferentialRentI derives SurplusProductQtrs and MoneyRentBP for each
// entry and returns the filled table header plus the computed entries. The
// regulating entry is the one with the smallest OutputQuarters (ties broken by
// lowest Grade), making the result deterministic regardless of input order.
// The caller's slice is not mutated. Empty input returns a zero table and nil.
func ComputeDifferentialRentI(entries []DifferentialRentIEntry) (DifferentialRentITable, []DifferentialRentIEntry) {
	if len(entries) == 0 {
		return DifferentialRentITable{}, nil
	}
	regIdx := 0
	for i := 1; i < len(entries); i++ {
		e, reg := entries[i], entries[regIdx]
		if e.OutputQuarters < reg.OutputQuarters ||
			(e.OutputQuarters == reg.OutputQuarters && e.Grade < reg.Grade) {
			regIdx = i
		}
	}
	regulatingOutput := entries[regIdx].OutputQuarters
	regulatingGrade := entries[regIdx].Grade

	computed := make([]DifferentialRentIEntry, len(entries))
	var totalRentalBP int64
	for i, e := range entries {
		e.SurplusProductQtrs = e.OutputQuarters - regulatingOutput
		e.MoneyRentBP = e.SurplusProductQtrs * e.PriceOfProductionBP / 100
		totalRentalBP += e.MoneyRentBP
		computed[i] = e
	}
	return DifferentialRentITable{TotalRentalBP: totalRentalBP, RegulatingGrade: regulatingGrade}, computed
}
