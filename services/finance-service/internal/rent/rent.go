// Package rent implements Capital Vol. III, Part VI, Ch. 37
// ("Introduction" to Ground-Rent). Part VI opens here. Ground-rent is the
// portion of surplus-value that accrues to the landowner as a condition of
// access to land. Three forms are distinguished: differential (from unequal
// fertility/location), absolute (private property charges even the worst soil),
// and monopoly (unique natural/social monopolies). The canonical three-class
// fixture is a Landowner (receives rent), an AgriculturalCapitalist (advances
// capital, pays rent), and the LandParcel they contest.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNegativeRent      = errors.New("rent: amount_labour_minutes must be >= 0")
	ErrSameOwnerOperator = errors.New("rent: land_owner_id and capitalist_id must differ")
	ErrInvalidRentForm   = errors.New("rent: rent form is not valid")
)

// RentForm enumerates the three forms of ground-rent (Vol. III Part VI).
type RentForm int

const (
	RentFormDifferential RentForm = 1
	RentFormAbsolute     RentForm = 2
	RentFormMonopoly     RentForm = 3
)

func (f RentForm) IsValid() bool {
	return f == RentFormDifferential || f == RentFormAbsolute || f == RentFormMonopoly
}

type LandParcelID string

func NewLandParcelID() LandParcelID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return LandParcelID(hex.EncodeToString(b[:]))
}

// LandParcel is a plot of agricultural land. FertilityGrade runs 1 (A, worst —
// the regulating parcel) to 4 (D, best). OwnerID references the Landowner.
type LandParcel struct {
	ID             LandParcelID `json:"id"`
	OwnerID        string       `json:"owner_id"`
	FertilityGrade int          `json:"fertility_grade"`
	AreaHectares   int64        `json:"area_hectares"`
	Location       string       `json:"location"`
	CreatedAt      time.Time    `json:"created_at"`
}

type LandownerID string

func NewLandownerID() LandownerID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return LandownerID(hex.EncodeToString(b[:]))
}

// Landowner holds legal title and collects ground-rent without productive labour.
type Landowner struct {
	ID        LandownerID `json:"id"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"created_at"`
}

type AgriculturalCapitalistID string

func NewAgriculturalCapitalistID() AgriculturalCapitalistID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return AgriculturalCapitalistID(hex.EncodeToString(b[:]))
}

// AgriculturalCapitalist advances capital on leased land and pays ground-rent.
// CapitalAdvanced is in LabourMinutes.
type AgriculturalCapitalist struct {
	ID              AgriculturalCapitalistID `json:"id"`
	CapitalAdvanced int64                    `json:"capital_advanced"`
	LeaseParcelID   string                   `json:"lease_parcel_id"`
	CreatedAt       time.Time                `json:"created_at"`
}

type GroundRentID string

func NewGroundRentID() GroundRentID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return GroundRentID(hex.EncodeToString(b[:]))
}

// GroundRent records a rent payment from an AgriculturalCapitalist to a
// Landowner for a LandParcel. AmountLabourMinutes is the canonical value unit.
type GroundRent struct {
	ID                  GroundRentID `json:"id"`
	ParcelID            string       `json:"parcel_id"`
	CapitalistID        string       `json:"capitalist_id"`
	LandOwnerID         string       `json:"land_owner_id"`
	Form                RentForm     `json:"form"`
	AmountLabourMinutes int64        `json:"amount_labour_minutes"`
	PeriodYears         int          `json:"period_years"`
	CreatedAt           time.Time    `json:"created_at"`
}

// Validate enforces Ch. 37 invariants: non-negative rent, distinct
// owner/operator, and a valid rent form.
func (g GroundRent) Validate() error {
	if g.AmountLabourMinutes < 0 {
		return ErrNegativeRent
	}
	if g.LandOwnerID == g.CapitalistID {
		return ErrSameOwnerOperator
	}
	if !g.Form.IsValid() {
		return ErrInvalidRentForm
	}
	return nil
}
