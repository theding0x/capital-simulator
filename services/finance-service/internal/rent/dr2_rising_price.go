// Package rent — Ch. 43 types: Differential Rent II, Third Case (Rising Price
// of Production). Engels-completed (Tables XI–XXIV). Price rises when added
// capital yields below average or a new inferior soil enters. Differential rent
// is an arithmetic progression of rent differences from the rentless regulating
// soil. THIS CHAPTER USES SHILLINGS + BUSHELS (Engels's units), kept as distinct
// fields from the pence/quarter types of earlier chapters.
//
// Engels Table XI: base=0, increment=12, grades=5 → [0,12,24,36,48] sum=120s.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// RisingPriceVariant enumerates the five DR II rising-price sub-cases.
type RisingPriceVariant int

const (
	RisingPriceEventualityAVariant1 RisingPriceVariant = 1 // A regulates, constant productivity 2nd investment
	RisingPriceEventualityAVariant2 RisingPriceVariant = 2 // A regulates, falling productivity 2nd investment
	RisingPriceEventualityBVariant1 RisingPriceVariant = 3 // inferior soil 'a' regulates, constant productivity
	RisingPriceEventualityBVariant2 RisingPriceVariant = 4 // inferior soil 'a' regulates, falling productivity
	RisingPriceEventualityBVariant3 RisingPriceVariant = 5 // inferior soil 'a' regulates, rising productivity
)

// IsValid reports whether v is a defined RisingPriceVariant.
func (v RisingPriceVariant) IsValid() bool {
	return v >= RisingPriceEventualityAVariant1 && v <= RisingPriceEventualityBVariant3
}

// RisingPriceDR2OutcomeID identifies a rising-price DR II outcome.
type RisingPriceDR2OutcomeID string

// NewRisingPriceDR2OutcomeID returns a fresh 96-bit hex ID.
func NewRisingPriceDR2OutcomeID() RisingPriceDR2OutcomeID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return RisingPriceDR2OutcomeID(hex.EncodeToString(b[:]))
}

// RisingPriceDR2Outcome tabulates a rising-price scenario in shillings + bushels.
type RisingPriceDR2Outcome struct {
	ID                      RisingPriceDR2OutcomeID `json:"id"`
	Variant                 RisingPriceVariant      `json:"variant"`
	TotalCapitalShillings   int64                   `json:"total_capital_shillings"`
	TotalRentShillings      int64                   `json:"total_rent_shillings"`
	TotalOutputBushels      int64                   `json:"total_output_bushels"`
	PricePerBushelShillings int64                   `json:"price_per_bushel_shillings"`
	CreatedAt               time.Time               `json:"created_at"`
}

// RentSequenceID identifies a rent-sequence record.
type RentSequenceID string

// NewRentSequenceID returns a fresh 96-bit hex ID.
func NewRentSequenceID() RentSequenceID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return RentSequenceID(hex.EncodeToString(b[:]))
}

// RentSequence is the arithmetic progression of rents across soil grades:
// rent[g] = BaseRent + (g-1)*Increment for g in 1..NumGrades.
type RentSequence struct {
	ID        RentSequenceID `json:"id"`
	BaseRent  int64          `json:"base_rent"`
	Increment int64          `json:"increment"`
	NumGrades int            `json:"num_grades"`
	CreatedAt time.Time      `json:"created_at"`
}

// EngelsTableEntryID identifies one Engels-table row.
type EngelsTableEntryID string

// NewEngelsTableEntryID returns a fresh 96-bit hex ID.
func NewEngelsTableEntryID() EngelsTableEntryID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return EngelsTableEntryID(hex.EncodeToString(b[:]))
}

// EngelsTableEntry is one soil-grade row from an Engels Ch. 43 table
// (Tables XI–XXIV), in shillings + bushels with a string grade label (A–E, 'a').
type EngelsTableEntry struct {
	ID               EngelsTableEntryID `json:"id"`
	TableNumber      int                `json:"table_number"`
	SoilGradeLabel   string             `json:"soil_grade_label"`
	OutputBushels    int64              `json:"output_bushels"`
	RentShillings    int64              `json:"rent_shillings"`
	CapitalShillings int64              `json:"capital_shillings"`
	CreatedAt        time.Time          `json:"created_at"`
}

// ComputeRentSequence returns the rents for grades 1..numGrades where
// rent[g] = baseRent + (g-1)*increment. numGrades <= 0 returns a non-nil
// empty slice. Deterministic.
//
// Engels Table XI: ComputeRentSequence(0, 12, 5) == [0,12,24,36,48] (sum 120).
func ComputeRentSequence(baseRent, increment int64, numGrades int) []int64 {
	capacity := numGrades
	if capacity < 0 {
		capacity = 0
	}
	out := make([]int64, 0, capacity)
	for g := 1; g <= numGrades; g++ {
		out = append(out, baseRent+int64(g-1)*increment)
	}
	return out
}
