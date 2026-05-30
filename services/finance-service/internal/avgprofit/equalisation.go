package avgprofit

// This file carries the capital-flow and equalisation logic for Ch. 10.
// Competition pulls rates toward the general rate: capital flows into
// spheres with above-average rates (inflow), out of spheres with below-
// average rates (outflow), and stabilises when rates are equal (stable).
// An Equalisation record captures one sphere's journey toward the general
// rate, together with the CapitalFlow that drives it.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// CapitalFlowDirection encodes the direction of competitive capital movement.
type CapitalFlowDirection string

const (
	KindInflow  CapitalFlowDirection = "inflow"  // toward an above-average-rate sphere
	KindOutflow CapitalFlowDirection = "outflow" // away from a below-average-rate sphere
	KindStable  CapitalFlowDirection = "stable"  // sphere rate == general rate
)

// capitalFlowDirection maps signal vs benchmark to a direction constant.
// signal > benchmark -> KindInflow; signal < benchmark -> KindOutflow; equal -> KindStable.
func capitalFlowDirection(signal, benchmark int64) CapitalFlowDirection {
	switch {
	case signal > benchmark:
		return KindInflow
	case signal < benchmark:
		return KindOutflow
	default:
		return KindStable
	}
}

// CapitalFlow is a value object recording what drives competitive movement
// in a given sphere. Direction is derived from the rate comparison when
// generalRate != 0; otherwise from market price vs market value.
type CapitalFlow struct {
	SphereName        string               `json:"sphere_name"`
	CurrentSphereRate SphereProfitRate     `json:"current_sphere_rate"`
	GeneralRate       SphereProfitRate     `json:"general_rate"`
	MarketPrice       MarketPriceAmount    `json:"market_price"`
	MarketValue       CommodityValue       `json:"market_value"`
	Direction         CapitalFlowDirection `json:"direction"`
}

// ComputeCapitalFlow derives the CapitalFlow direction for a sphere.
// Priority: if generalRate != 0 use rate comparison; else if marketValue != 0
// use price vs value; else KindStable.
func ComputeCapitalFlow(name string, currentRate, generalRate SphereProfitRate, marketPrice MarketPriceAmount, marketValue CommodityValue) CapitalFlow {
	var dir CapitalFlowDirection
	switch {
	case generalRate != 0:
		dir = capitalFlowDirection(int64(currentRate), int64(generalRate))
	case marketValue != 0:
		dir = capitalFlowDirection(int64(marketPrice), int64(marketValue))
	default:
		dir = KindStable
	}
	return CapitalFlow{
		SphereName:        name,
		CurrentSphereRate: currentRate,
		GeneralRate:       generalRate,
		MarketPrice:       marketPrice,
		MarketValue:       marketValue,
		Direction:         dir,
	}
}

// EqualisationID identifies a persisted equalisation record.
// 96-bit hex from crypto/rand.
type EqualisationID string

// NewEqualisationID returns a fresh random EqualisationID.
func NewEqualisationID() EqualisationID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("avgprofit: rand: %w", err))
	}
	return EqualisationID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty EqualisationID.
func (id EqualisationID) IsZero() bool { return id == "" }

// Equalisation is a persisted record of competition pulling a sphere's rate
// toward the general rate. Flow carries the CapitalFlow value object that
// drove the equalisation; it is rebuilt from flat columns on read.
type Equalisation struct {
	ID           EqualisationID       `json:"id"`
	SphereName   string               `json:"sphere_name"`
	InitialRate  SphereProfitRate     `json:"initial_rate"`  // sphere's current rate
	TargetRate   SphereProfitRate     `json:"target_rate"`   // general rate
	Direction    CapitalFlowDirection `json:"direction"`
	IsConverging bool                 `json:"is_converging"` // InitialRate != TargetRate
	Flow         CapitalFlow          `json:"flow"`
	CreatedAt    time.Time            `json:"created_at"`
}

// Validate returns ErrNegativeMagnitude if either rate is negative.
func (e Equalisation) Validate() error {
	if e.InitialRate < 0 || e.TargetRate < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// ComputeEqualisation builds an Equalisation for a sphere converging toward
// the general rate. It calls ComputeCapitalFlow with zero market price/value
// so the rate-comparison path is used.
func ComputeEqualisation(name string, initialRate, targetRate SphereProfitRate) Equalisation {
	flow := ComputeCapitalFlow(name, initialRate, targetRate, 0, 0)
	return Equalisation{
		SphereName:   name,
		InitialRate:  initialRate,
		TargetRate:   targetRate,
		Direction:    flow.Direction,
		IsConverging: initialRate != targetRate,
		Flow:         flow,
	}
}
