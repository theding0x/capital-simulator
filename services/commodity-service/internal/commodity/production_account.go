package commodity

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/theding0x/capital-simulator/pkg/surplusrate"
)

// ProductionAccountID is a 96-bit hex ID for ProductionAccount records.
type ProductionAccountID string

func (id ProductionAccountID) IsZero() bool { return id == "" }

func NewProductionAccountID() ProductionAccountID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return ProductionAccountID(hex.EncodeToString(b[:]))
}

// SurplusValue is the labour-time produced above the value of labour-power [§1].
type SurplusValue LabourMinutes

// RateOfSurplusValue is s/v — the degree of exploitation of labour-power [§1].
type RateOfSurplusValue float64

// ValueProduct is the new value created by living labour: v + s [§1].
type ValueProduct LabourMinutes

// NecessaryLabour is the portion of the working day that reproduces the value
// of labour-power (variable capital) [§1].
type NecessaryLabour LabourMinutes

// SurplusLabour is the portion of the working day that produces surplus-value [§1].
type SurplusLabour LabourMinutes

// WorkingDay is the total labour-time: NecessaryLabour + SurplusLabour [§4].
type WorkingDay LabourMinutes

// ErrDivisionByZero is returned when computing a rate against zero variable capital.
var ErrDivisionByZero = errors.New("commodity: variable capital cannot be zero")

// ProductionAccount records the capital magnitudes for a single production run [§1].
type ProductionAccount struct {
	ID        ProductionAccountID `json:"id"`
	Constant  LabourMinutes       `json:"constant"`  // c: transferred value of means of production
	Variable  LabourMinutes       `json:"variable"`  // v: value of labour-power (wage)
	Surplus   SurplusValue        `json:"surplus"`   // s: surplus-value produced
	CreatedAt time.Time           `json:"created_at"`
}

// ValueProduct returns v + s — the new value created by living labour [§1].
func (a ProductionAccount) ValueProduct() ValueProduct {
	return ValueProduct(a.Variable + LabourMinutes(a.Surplus))
}

// ExpandedCapital returns c + v + s — the total capital after production [§1].
func (a ProductionAccount) ExpandedCapital() LabourMinutes {
	return a.Constant + a.Variable + LabourMinutes(a.Surplus)
}

// Rate returns s/v — the rate of surplus-value [§1].
func (a ProductionAccount) Rate() (RateOfSurplusValue, error) {
	return ComputeRate(a.Surplus, a.Variable)
}

// Validate reports if a is well-formed.
func (a ProductionAccount) Validate() error {
	if a.Constant < 0 {
		return errors.New("commodity: production_account constant must be non-negative")
	}
	if a.Variable < 0 {
		return errors.New("commodity: production_account variable must be non-negative")
	}
	if a.Surplus < 0 {
		return errors.New("commodity: production_account surplus must be non-negative")
	}
	return nil
}

// ComputeRate returns s/v = surplus-value / variable-capital [§1].
// Returns ErrDivisionByZero if v <= 0.
func ComputeRate(s SurplusValue, v LabourMinutes) (RateOfSurplusValue, error) {
	r, err := surplusrate.Ratio(int64(s), int64(v))
	if err != nil {
		return 0, ErrDivisionByZero
	}
	return RateOfSurplusValue(r), nil
}

// ComputeRateBasisPoints returns s/v in basis points rounded half-up (10000 =
// 100%), the same scaled-integer form Vol. III's profit.ComputeRateOfSurplusValue
// uses, so Vol. I's s feeds Vol. III's p' without a unit seam [§1].
// Returns ErrDivisionByZero if v <= 0.
func ComputeRateBasisPoints(s SurplusValue, v LabourMinutes) (int64, error) {
	bp, err := surplusrate.BasisPoints(int64(s), int64(v))
	if err != nil {
		return 0, ErrDivisionByZero
	}
	return bp, nil
}

// ComputeValueProduct returns v + s — the new value created [§1].
func ComputeValueProduct(v LabourMinutes, s SurplusValue) ValueProduct {
	return ValueProduct(v + LabourMinutes(s))
}

// SurplusProduceRatio returns surplus-labour / necessary-labour [§4] — the
// ratio of surplus-produce to the part of the product in which necessary
// labour is incorporated, not to the total product. It equals s' by
// construction. Returns 0 if necessaryLabour == 0.
func SurplusProduceRatio(surplusLabour, necessaryLabour LabourMinutes) float64 {
	if necessaryLabour == 0 {
		return 0
	}
	return float64(surplusLabour) / float64(necessaryLabour)
}
