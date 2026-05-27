// Package circulation — Vol. II Ch. 15: Effects of a Change of Prices.
//
// When prices of raw materials, labour-power, or output shift mid-turnover,
// the capital advanced and surplus-value realised are both affected. This file
// adds the validation logic for ValueRevolutionEvent extensions, derives
// RevaluedCapital (new advance required to restart production),
// RealisedValueDelta (gap between expected and actual realisation), and
// InventoryRevaluation (paper gains/losses on stock held when prices move).
//
// The canonical invariant: surplus is measured against the *original* advance,
// not the revalued one.
// §: "If the value (price) of raw and auxiliary materials falls, then in the
// renewed process of production less money-capital is required to buy the same
// amount of them; a portion of the money-capital previously employed is set free."
package circulation

import (
	"fmt"
	"time"
)

// --- Enums ------------------------------------------------------------------

// PriceFallCase names the three possible responses to a fall in input prices.
// §: "He can (1) continue the same scale of production with less money-capital;
// (2) extend the scale of production; (3) accumulate the released capital as
// a hoard."
type PriceFallCase string

const (
	// CaseSameScale: production continues at the same scale; the released
	// capital sits idle or is hoarded.
	CaseSameScale PriceFallCase = "same_scale"
	// CaseExpandedScale: the freed capital is reinvested to expand output.
	CaseExpandedScale PriceFallCase = "expanded_scale"
	// CaseLargerStock: the freed capital is used to hold a larger raw-material
	// stock, anticipating a future price rebound.
	CaseLargerStock PriceFallCase = "larger_stock"
)

// PriceRiseCase names the three possible responses to a rise in input prices.
// §: "(1) Contraction of the scale of production; (2) additional money-capital;
// (3) accumulation fund is drawn upon."
type PriceRiseCase string

const (
	// CaseContractedScale: production shrinks because the same advance now
	// buys fewer inputs.
	CaseContractedScale PriceRiseCase = "contracted_scale"
	// CaseAdditionalCapital: new money-capital is injected to maintain scale.
	CaseAdditionalCapital PriceRiseCase = "additional_capital"
	// CaseAccumulationDrained: the capitalist draws on his accumulation fund
	// (sinking fund or reserve) to cover the higher advance.
	CaseAccumulationDrained PriceRiseCase = "accumulation_drained"
)

// --- ID types ---------------------------------------------------------------

// RevaluedCapitalID is a 96-bit hex identifier for a RevaluedCapital record.
type RevaluedCapitalID string

func (id RevaluedCapitalID) IsZero() bool { return id == "" }

// NewRevaluedCapitalID generates a 96-bit hex identifier via crypto/rand.
func NewRevaluedCapitalID() RevaluedCapitalID { return RevaluedCapitalID(newHexID()) }

// RealisedValueDeltaID is a 96-bit hex identifier for a RealisedValueDelta record.
type RealisedValueDeltaID string

func (id RealisedValueDeltaID) IsZero() bool { return id == "" }

// NewRealisedValueDeltaID generates a 96-bit hex identifier via crypto/rand.
func NewRealisedValueDeltaID() RealisedValueDeltaID {
	return RealisedValueDeltaID(newHexID())
}

// InventoryRevaluationID is a 96-bit hex identifier for an InventoryRevaluation record.
type InventoryRevaluationID string

func (id InventoryRevaluationID) IsZero() bool { return id == "" }

// NewInventoryRevaluationID generates a 96-bit hex identifier via crypto/rand.
func NewInventoryRevaluationID() InventoryRevaluationID {
	return InventoryRevaluationID(newHexID())
}

// SpeculativeHoldID is a 96-bit hex identifier for a SpeculativeHold record.
type SpeculativeHoldID string

func (id SpeculativeHoldID) IsZero() bool { return id == "" }

// NewSpeculativeHoldID generates a 96-bit hex identifier via crypto/rand.
func NewSpeculativeHoldID() SpeculativeHoldID { return SpeculativeHoldID(newHexID()) }

// CompoundCapitalAdjustmentID is a 96-bit hex identifier for a
// CompoundCapitalAdjustment record.
type CompoundCapitalAdjustmentID string

func (id CompoundCapitalAdjustmentID) IsZero() bool { return id == "" }

// NewCompoundCapitalAdjustmentID generates a 96-bit hex identifier via crypto/rand.
func NewCompoundCapitalAdjustmentID() CompoundCapitalAdjustmentID {
	return CompoundCapitalAdjustmentID(newHexID())
}

// --- RevaluedCapital --------------------------------------------------------

// RevaluedCapital records the change in the advance of capital required to
// restart production after a price revolution on means of production or
// labour-power. DeltaPence is signed: positive means more capital is needed
// (price rose); negative means capital is set free (price fell).
// §: "In the case of a fall in value of raw materials, the advance of capital
// in the next process of reproduction is diminished."
type RevaluedCapital struct {
	ID                     RevaluedCapitalID      `json:"id"`
	ValueRevolutionEventID ValueRevolutionEventID `json:"value_revolution_event_id"`
	IndustrialCapitalID    IndustrialCapitalID    `json:"industrial_capital_id"`
	OldAdvancePence        Pence                  `json:"old_advance_pence"`
	NewAdvancePence        Pence                  `json:"new_advance_pence"`
	DeltaPence             Pence                  `json:"delta_pence"`
}

// Validate checks that DeltaPence == NewAdvancePence - OldAdvancePence.
func (rc RevaluedCapital) Validate() error {
	if rc.OldAdvancePence <= 0 {
		return fmt.Errorf("circulation: old_advance_pence must be positive")
	}
	if rc.NewAdvancePence <= 0 {
		return fmt.Errorf("circulation: new_advance_pence must be positive")
	}
	expected := rc.NewAdvancePence - rc.OldAdvancePence
	if rc.DeltaPence != expected {
		return fmt.Errorf("circulation: delta_pence %d != new_advance_pence - old_advance_pence %d",
			rc.DeltaPence, expected)
	}
	return nil
}

// DeriveRevaluedCapital constructs a RevaluedCapital from an event and the
// old advance, applying the basis-points change:
// NewAdvancePence = OldAdvancePence * (10000 + BasisPointsChange) / 10000.
func DeriveRevaluedCapital(evt ValueRevolutionEvent, oldAdvancePence Pence) RevaluedCapital {
	newAdv := Pence(int64(oldAdvancePence) * (10000 + evt.BasisPointsChange) / 10000)
	return RevaluedCapital{
		ID:                     NewRevaluedCapitalID(),
		ValueRevolutionEventID: evt.ID,
		IndustrialCapitalID:    evt.IndustrialCapitalID,
		OldAdvancePence:        oldAdvancePence,
		NewAdvancePence:        newAdv,
		DeltaPence:             newAdv - oldAdvancePence,
	}
}

// --- RealisedValueDelta -----------------------------------------------------

// RealisedValueDelta records the gap between the expected and actual
// realisation of commodity-capital when output prices shift.
// AffectingOutput events change RealisedValueDelta, not RevaluedCapital.
// §: "If the value of the commodities falls after they have been produced but
// before they have been sold, the capitalist loses part of his surplus-value."
type RealisedValueDelta struct {
	ID                       RealisedValueDeltaID   `json:"id"`
	ValueRevolutionEventID   ValueRevolutionEventID `json:"value_revolution_event_id"`
	IndustrialCapitalID      IndustrialCapitalID    `json:"industrial_capital_id"`
	ExpectedRealisationPence Pence                  `json:"expected_realisation_pence"`
	ActualRealisationPence   Pence                  `json:"actual_realisation_pence"`
	DeltaPence               Pence                  `json:"delta_pence"`
}

// SurplusAgainstOriginalAdvance returns the surplus measured against the
// original advance, not the revalued capital. This is the canonical invariant
// of Ch. 15: surplus = actual_realisation - old_advance_pence.
func (r RealisedValueDelta) SurplusAgainstOriginalAdvance(oldAdvancePence Pence) Pence {
	return r.ActualRealisationPence - oldAdvancePence
}

// Validate checks that DeltaPence == ActualRealisationPence - ExpectedRealisationPence.
func (r RealisedValueDelta) Validate() error {
	if r.ExpectedRealisationPence <= 0 {
		return fmt.Errorf("circulation: expected_realisation_pence must be positive")
	}
	expected := r.ActualRealisationPence - r.ExpectedRealisationPence
	if r.DeltaPence != expected {
		return fmt.Errorf("circulation: delta_pence %d != actual - expected %d",
			r.DeltaPence, expected)
	}
	return nil
}

// --- InventoryRevaluation ---------------------------------------------------

// InventoryRevaluation records a paper gain or loss on commodity stock held
// when prices move. The delta is unrealised until the inventory leaves stock.
// §: "The owner of a stock of raw material … suffers a paper loss on his
// existing stock corresponding to the fall in price."
type InventoryRevaluation struct {
	ID                  InventoryRevaluationID `json:"id"`
	IndustrialCapitalID IndustrialCapitalID    `json:"industrial_capital_id"`
	CommodityID         string                 `json:"commodity_id"`
	QuantityHeld        int64                  `json:"quantity_held"`
	OldUnitPence        Pence                  `json:"old_unit_pence"`
	NewUnitPence        Pence                  `json:"new_unit_pence"`
	DeltaPence          Pence                  `json:"delta_pence"`
}

// Validate checks that DeltaPence == QuantityHeld * (NewUnitPence - OldUnitPence).
func (inv InventoryRevaluation) Validate() error {
	if inv.QuantityHeld <= 0 {
		return fmt.Errorf("circulation: quantity_held must be positive")
	}
	if inv.OldUnitPence <= 0 {
		return fmt.Errorf("circulation: old_unit_pence must be positive")
	}
	if inv.NewUnitPence <= 0 {
		return fmt.Errorf("circulation: new_unit_pence must be positive")
	}
	expected := Pence(inv.QuantityHeld) * (inv.NewUnitPence - inv.OldUnitPence)
	if inv.DeltaPence != expected {
		return fmt.Errorf("circulation: delta_pence %d != quantity * (new - old) %d",
			inv.DeltaPence, expected)
	}
	return nil
}

// --- SpeculativeHold --------------------------------------------------------

// SpeculativeHold registers that a capitalist is holding a commodity in
// anticipation of a favourable price movement. The key Marxian point:
// any monetary gain from speculation is appropriation of value already
// produced by others, not creation of new surplus-value.
// §: "This is a form of the concentration of wealth … not creation of value."
type SpeculativeHold struct {
	ID                   SpeculativeHoldID   `json:"id"`
	IndustrialCapitalID  IndustrialCapitalID `json:"industrial_capital_id"`
	CommodityID          string              `json:"commodity_id"`
	HeldSince            time.Time           `json:"held_since"`
	AnticipatedDirection int8                `json:"anticipated_direction"`
}

// Validate checks that AnticipatedDirection is +1 (expecting price rise) or
// -1 (expecting price fall).
func (s SpeculativeHold) Validate() error {
	if s.AnticipatedDirection != 1 && s.AnticipatedDirection != -1 {
		return fmt.Errorf("circulation: anticipated_direction must be +1 or -1")
	}
	if s.CommodityID == "" {
		return fmt.Errorf("circulation: commodity_id is required")
	}
	return nil
}

// --- CompoundCapitalAdjustment ----------------------------------------------

// CompoundCapitalAdjustment aggregates the three delta streams — revaluation,
// realisation, and inventory — for one industrial capital over one period.
// Invariant: NetEffectPence == RevaluationDeltaPence + RealisationDeltaPence + InventoryDeltaPence.
// §: "The effect of a revolution in value is different for different parts of the
// industrial capital."
type CompoundCapitalAdjustment struct {
	ID                    CompoundCapitalAdjustmentID `json:"id"`
	IndustrialCapitalID   IndustrialCapitalID         `json:"industrial_capital_id"`
	Period                string                      `json:"period"`
	RevaluationDeltaPence Pence                       `json:"revaluation_delta_pence"`
	RealisationDeltaPence Pence                       `json:"realisation_delta_pence"`
	InventoryDeltaPence   Pence                       `json:"inventory_delta_pence"`
	NetEffectPence        Pence                       `json:"net_effect_pence"`
}

// NetEffect returns the algebraic sum of the three delta streams.
func (c CompoundCapitalAdjustment) NetEffect() Pence {
	return c.RevaluationDeltaPence + c.RealisationDeltaPence + c.InventoryDeltaPence
}

// Validate checks the net-effect invariant.
func (c CompoundCapitalAdjustment) Validate() error {
	expected := c.NetEffect()
	if c.NetEffectPence != expected {
		return fmt.Errorf("circulation: net_effect_pence %d != sum of deltas %d",
			c.NetEffectPence, expected)
	}
	return nil
}

// --- PriceCaseRecord --------------------------------------------------------

// PriceCaseRecord stores the operator-selected outcome case for a
// ValueRevolutionEvent. FallCase and RiseCase are mutually exclusive.
type PriceCaseRecord struct {
	ValueRevolutionEventID ValueRevolutionEventID `json:"value_revolution_event_id"`
	FallCase               *PriceFallCase         `json:"fall_case,omitempty"`
	RiseCase               *PriceRiseCase         `json:"rise_case,omitempty"`
}

// Validate checks mutual exclusivity of FallCase and RiseCase.
func (p PriceCaseRecord) Validate() error {
	if p.FallCase != nil && p.RiseCase != nil {
		return fmt.Errorf("circulation: fall_case and rise_case are mutually exclusive")
	}
	if p.FallCase == nil && p.RiseCase == nil {
		return fmt.Errorf("circulation: one of fall_case or rise_case is required")
	}
	return nil
}

// --- PriceRevolutionRecord --------------------------------------------------

// PriceRevolutionRecord bundles a ValueRevolutionEvent with all derived rows,
// returned by GET /v1/price-revolutions/{id}.
type PriceRevolutionRecord struct {
	Event           ValueRevolutionEvent   `json:"event"`
	RevaluedCapital *RevaluedCapital       `json:"revalued_capital,omitempty"`
	RealisedDelta   *RealisedValueDelta    `json:"realised_delta,omitempty"`
	Revaluations    []InventoryRevaluation `json:"inventory_revaluations"`
	FallCase        *PriceFallCase         `json:"fall_case,omitempty"`
	RiseCase        *PriceRiseCase         `json:"rise_case,omitempty"`
}
