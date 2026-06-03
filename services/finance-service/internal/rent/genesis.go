// Package rent — Ch. 47 types: Genesis of Capitalist Ground-Rent.
// Traces rent's historical evolution through three pre-capitalist forms —
// labour-rent (direct surplus-labour on the lord's land), product-rent
// (a portion of the product delivered), and money-rent (a money sum; the
// transitional form) — to fully developed capitalist ground-rent, which
// requires the three-class separation of landowner, capitalist, and labourer.
// The small peasant proprietor conflates rent, profit, and wages into one
// undifferentiated income.
package rent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 47 sentinel errors.
var (
	// ErrInvalidRentFormStage is returned when a stage is outside 1..4.
	ErrInvalidRentFormStage = errors.New("rent: stage must be 1..4")

	// ErrRentTransitionNotProgressive is returned when a transition is not strictly
	// progressive (FromStage < ToStage). Rent evolves one way through history.
	ErrRentTransitionNotProgressive = errors.New("rent: rent-form transition must be progressive (from_stage < to_stage)")
)

// RentFormStage enumerates the historically ordered stages of rent.
type RentFormStage int

const (
	RentFormLabour     RentFormStage = 1 // direct labour service on the lord's land
	RentFormProduct    RentFormStage = 2 // delivery of a portion of the product
	RentFormMoney      RentFormStage = 3 // money payment; transitional to capitalism
	RentFormCapitalist RentFormStage = 4 // fully developed capitalist ground-rent
)

// IsValid reports whether s is a defined RentFormStage.
func (s RentFormStage) IsValid() bool {
	return s >= RentFormLabour && s <= RentFormCapitalist
}

// RentHistoricalFormID identifies a rent-historical-form record.
type RentHistoricalFormID string

// NewRentHistoricalFormID returns a fresh 96-bit hex ID.
func NewRentHistoricalFormID() RentHistoricalFormID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return RentHistoricalFormID(hex.EncodeToString(b[:]))
}

// RentHistoricalForm records one stage of rent's historical development.
type RentHistoricalForm struct {
	ID            RentHistoricalFormID `json:"id"`
	Stage         RentFormStage        `json:"stage"`
	Name          string               `json:"name"`
	CoercionKind  string               `json:"coercion_kind"`
	SurplusForm   string               `json:"surplus_form"`
	ExampleRegion string               `json:"example_region"`
	ExamplePeriod string               `json:"example_period"`
	CreatedAt     time.Time            `json:"created_at"`
}

// Validate reports whether the historical form names a defined stage (1..4).
func (f RentHistoricalForm) Validate() error {
	if !f.Stage.IsValid() {
		return ErrInvalidRentFormStage
	}
	return nil
}

// HistoricalRentTransitionID identifies a historical-rent-transition record.
type HistoricalRentTransitionID string

// NewHistoricalRentTransitionID returns a fresh 96-bit hex ID.
func NewHistoricalRentTransitionID() HistoricalRentTransitionID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return HistoricalRentTransitionID(hex.EncodeToString(b[:]))
}

// HistoricalRentTransition records a progressive movement from one rent form to
// the next. Transitions are always progressive: FromStage < ToStage.
type HistoricalRentTransition struct {
	ID            HistoricalRentTransitionID `json:"id"`
	FromStage     RentFormStage              `json:"from_stage"`
	ToStage       RentFormStage              `json:"to_stage"`
	Preconditions string                     `json:"preconditions"`
	DrivingForce  string                     `json:"driving_force"`
	CreatedAt     time.Time                  `json:"created_at"`
}

// Validate enforces that a transition is well-formed and progressive: both stages
// must be defined (1..4) and the movement must advance (FromStage < ToStage). Rent
// evolves one way — labour → product → money → capitalist — so a transition that
// stays put or regresses is not a historical-rent transition.
func (tr HistoricalRentTransition) Validate() error {
	if !tr.FromStage.IsValid() || !tr.ToStage.IsValid() {
		return ErrInvalidRentFormStage
	}
	if tr.FromStage >= tr.ToStage {
		return ErrRentTransitionNotProgressive
	}
	return nil
}

// SmallPeasantProductionID identifies a small-peasant-production record.
type SmallPeasantProductionID string

// NewSmallPeasantProductionID returns a fresh 96-bit hex ID.
func NewSmallPeasantProductionID() SmallPeasantProductionID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("rent: rand: %w", err))
	}
	return SmallPeasantProductionID(hex.EncodeToString(b[:]))
}

// SmallPeasantProduction models the peasant proprietor who is simultaneously
// labourer, capitalist, and landowner. Rent, profit, and wages are conflated
// into one undifferentiated income: TotalIncomeBP = Rent + Profit + Wages.
type SmallPeasantProduction struct {
	ID                SmallPeasantProductionID `json:"id"`
	ParcelID          string                   `json:"parcel_id"`
	TotalIncomeBP     int64                    `json:"total_income_bp"`
	RentComponentBP   int64                    `json:"rent_component_bp"`
	ProfitComponentBP int64                    `json:"profit_component_bp"`
	WagesComponentBP  int64                    `json:"wages_component_bp"`
	IsConflated       bool                     `json:"is_conflated"`
	CreatedAt         time.Time                `json:"created_at"`
}

// ComputeTotalIncome returns the sum of the rent, profit, and wage components —
// the undifferentiated income of the small peasant proprietor.
func ComputeTotalIncome(rentComponentBP, profitComponentBP, wagesComponentBP int64) int64 {
	return rentComponentBP + profitComponentBP + wagesComponentBP
}
