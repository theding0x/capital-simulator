// Package rent — Ch. 40 types: Differential Rent II (Second Form).
// Successive capital investments on the SAME plot (intensive cultivation).
// DR II presupposes DR I: the worst soil still regulates the market price.
// With more capital concentrated on better soil, rent per acre RISES even
// when the rate of surplus-profit per tranche is unchanged.
//
// Marx doubled-investment arithmetic (Soil D, Ch. 40):
//
//	Tranche 1: capital £2.50 (250), output 4 qtr (400), surplus-profit £9 (900)
//	Tranche 2: capital £2.50 (250), output 4 qtr (400), surplus-profit £9 (900)
//	TotalCapital=500, TotalOutput=800, TotalRent=1800 (£18), RentPerAcre=1800 (1 acre)
//
// Units mirror differential_rent_i.go: *BP in pence; OutputQuarters in 1/100 qtr.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// SuccessiveInvestmentID identifies one capital tranche on a parcel.
type SuccessiveInvestmentID string

// NewSuccessiveInvestmentID returns a fresh 96-bit hex ID.
func NewSuccessiveInvestmentID() SuccessiveInvestmentID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return SuccessiveInvestmentID(hex.EncodeToString(b[:]))
}

// SuccessiveInvestment is one capital tranche applied to a parcel for intensive
// cultivation. SurplusProfitBP is the surplus-profit this tranche generates
// above the regulating (worst-soil) price of production.
type SuccessiveInvestment struct {
	ID              SuccessiveInvestmentID `json:"id"`
	ParcelID        string                 `json:"parcel_id"`
	TrancheNumber   int                    `json:"tranche_number"`
	CapitalBP       int64                  `json:"capital_bp"`
	OutputQuarters  int64                  `json:"output_quarters"`
	SurplusProfitBP int64                  `json:"surplus_profit_bp"`
	CreatedAt       time.Time              `json:"created_at"`
}

// DifferentialRentIITableID identifies a computed DR II aggregate result.
type DifferentialRentIITableID string

// NewDifferentialRentIITableID returns a fresh 96-bit hex ID. The DR II table
// is NOT persisted; this stamps the computed result for response identity.
func NewDifferentialRentIITableID() DifferentialRentIITableID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return DifferentialRentIITableID(hex.EncodeToString(b[:]))
}

// DifferentialRentIITable is the computed aggregate for all successive
// investments on one parcel. NOT persisted; built on demand by
// ComputeDifferentialRentII and returned in the GET dr2-tables response.
type DifferentialRentIITable struct {
	ID                  DifferentialRentIITableID `json:"id"`
	ParcelID            string                    `json:"parcel_id"`
	TotalCapitalBP      int64                     `json:"total_capital_bp"`
	TotalOutputQuarters int64                     `json:"total_output_quarters"`
	TotalRentBP         int64                     `json:"total_rent_bp"`
	RentPerAcreBP       int64                     `json:"rent_per_acre_bp"`
}

// LeaseTermID identifies a lease-term record.
type LeaseTermID string

// NewLeaseTermID returns a fresh 96-bit hex ID.
func NewLeaseTermID() LeaseTermID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return LeaseTermID(hex.EncodeToString(b[:]))
}

// LeaseTerm records the lease under which a capitalist farms a parcel.
// FixedRentBP is locked for the duration; improvements made during the term
// accrue to the capitalist, not the landowner — the central DR II observation.
type LeaseTerm struct {
	ID            LeaseTermID `json:"id"`
	ParcelID      string      `json:"parcel_id"`
	StartYear     int         `json:"start_year"`
	DurationYears int         `json:"duration_years"`
	FixedRentBP   int64       `json:"fixed_rent_bp"`
	CreatedAt     time.Time   `json:"created_at"`
}

// ComputeDifferentialRentII aggregates successive investments on one parcel
// into a DifferentialRentIITable. acres is the parcel size for rent-per-acre;
// acres <= 0 is treated as 1 (divide-by-zero guard). The caller sets the
// result ID (NewDifferentialRentIITableID()). The input slice is not mutated.
func ComputeDifferentialRentII(parcelID string, investments []SuccessiveInvestment, acres int64) DifferentialRentIITable {
	if acres <= 0 {
		acres = 1
	}
	var totalCapital, totalOutput, totalRent int64
	for _, inv := range investments {
		totalCapital += inv.CapitalBP
		totalOutput += inv.OutputQuarters
		totalRent += inv.SurplusProfitBP
	}
	return DifferentialRentIITable{
		ParcelID:            parcelID,
		TotalCapitalBP:      totalCapital,
		TotalOutputQuarters: totalOutput,
		TotalRentBP:         totalRent,
		RentPerAcreBP:       totalRent / acres,
	}
}
