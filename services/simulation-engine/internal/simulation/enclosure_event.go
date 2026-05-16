package simulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Ch. 27 — Expropriation of the Agricultural Population.
//
// The English enclosures (15th–18th c.) are Marx's paradigm case of how
// capital dispossesses the peasantry from the land. Each EnclosureEvent
// records one wave: the period, the acreage seized, the population driven
// off, and the beneficiary landlord class. Together they give substance to
// the abstract "expropriation of the agricultural producer" that underlies
// all subsequent capitalist development.

// EnclosureEventID is a 96-bit hex identifier for a persisted enclosure event.
type EnclosureEventID string

// IsZero reports whether the ID is unset.
func (id EnclosureEventID) IsZero() bool { return id == "" }

// NewEnclosureEventID generates a 96-bit hex identifier via crypto/rand.
func NewEnclosureEventID() EnclosureEventID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return EnclosureEventID(hex.EncodeToString(b))
}

var (
	// ErrAcresNonPositive — an enclosure event must seize a positive acreage.
	ErrAcresNonPositive = errors.New("simulation: acres_enclosed must be positive")

	// ErrPopulationNegative — population_displaced cannot be negative.
	ErrPopulationNegative = errors.New("simulation: population_displaced cannot be negative")

	// ErrBeneficiaryRequired — the beneficiary landlord class must be named.
	ErrBeneficiaryRequired = errors.New("simulation: beneficiary is required")

	// ErrPeriodRequired — the historical period must be named.
	ErrPeriodRequired = errors.New("simulation: period is required")
)

// EnclosureEvent records a single wave of agricultural expropriation: the
// period in which it occurred, the acreage of common or peasant land seized,
// the population of producers driven from the land, and the class of
// beneficiaries who appropriated it.
type EnclosureEvent struct {
	ID                  EnclosureEventID `json:"id"`
	Period              string           `json:"period"`
	AcresEnclosed       int64            `json:"acres_enclosed"`
	PopulationDisplaced int64            `json:"population_displaced"`
	Beneficiary         string           `json:"beneficiary"`
	CreatedAt           time.Time        `json:"created_at"`
}

// Validate enforces invariants on an enclosure event.
func (e EnclosureEvent) Validate() error {
	if e.Period == "" {
		return ErrPeriodRequired
	}
	if e.AcresEnclosed <= 0 {
		return ErrAcresNonPositive
	}
	if e.PopulationDisplaced < 0 {
		return ErrPopulationNegative
	}
	if e.Beneficiary == "" {
		return ErrBeneficiaryRequired
	}
	return nil
}