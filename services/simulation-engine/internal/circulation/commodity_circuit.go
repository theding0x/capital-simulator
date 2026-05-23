package circulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Vol. II Ch. 3 — The Circuit of Commodity-Capital.

var (
	ErrNotCommodityCapital        = errors.New("circulation: opening extreme must carry surplus-value (C′ not C)")
	ErrInsufficientMaterialBase   = errors.New("circulation: capitalised surplus has no material elements in MP sources")
	ErrUnsourcedMeansOfProduction = errors.New("circulation: means of production must trace to a producer, merchant, or import")
)

// CommodityCircuitID is a 96-bit hex identifier.
type CommodityCircuitID string

func (id CommodityCircuitID) IsZero() bool { return id == "" }

// NewCommodityCircuitID generates a 96-bit hex identifier via crypto/rand.
func NewCommodityCircuitID() CommodityCircuitID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return CommodityCircuitID(hex.EncodeToString(b))
}

// SocialCapitalLens is a flag that enables aggregate-across-circuits endpoints.
type SocialCapitalLens bool

// MeansOfProductionSourceKind names the origin of the second-extreme C in M—C.
type MeansOfProductionSourceKind string

const (
	SourceProducerCircuit MeansOfProductionSourceKind = "producer_circuit"
	SourceMerchantHolding MeansOfProductionSourceKind = "merchant_holding"
	SourceImport          MeansOfProductionSourceKind = "import"
)

// OpeningCommodityCapital is the C′ that opens the circuit.
// Unlike Form I and Form II, this extreme already contains both capital-value
// (c + v) and surplus-value (s) in commodity-form.
type OpeningCommodityCapital struct {
	ConstantPence Pence `json:"constant_pence"`
	VariablePence Pence `json:"variable_pence"`
	SurplusPence  Pence `json:"surplus_pence"`
	PoundsTotal   int64 `json:"pounds_total"`
}

// Total returns the full value of C′ = (c + v) + s.
func (o OpeningCommodityCapital) Total() Pence {
	return o.ConstantPence + o.VariablePence + o.SurplusPence
}

// ValueOriginalPence returns the capital-value portion c + v.
func (o OpeningCommodityCapital) ValueOriginalPence() Pence {
	return o.ConstantPence + o.VariablePence
}

// CommodityAugmented is the closing C′ produced by the P phase.
// Simple reproduction: Total() == opening Total().
// Extended reproduction: Total() > opening Total() and CapitalisedPence > 0.
type CommodityAugmented struct {
	CommodityCircuitID CommodityCircuitID `json:"commodity_circuit_id"`
	ConstantPence      Pence              `json:"constant_pence"`
	VariablePence      Pence              `json:"variable_pence"`
	SurplusPence       Pence              `json:"surplus_pence"`
	CapitalisedPence   Pence              `json:"capitalised_pence"`
	PoundsTotal        int64              `json:"pounds_total"`
	ClosedAt           time.Time          `json:"closed_at"`
}

// Total returns the full value of the closing C′.
func (a CommodityAugmented) Total() Pence {
	return a.ConstantPence + a.VariablePence + a.SurplusPence
}

// IsExtended reports whether the closing circuit is larger than its opening.
func (a CommodityAugmented) IsExtended(openingTotal Pence) bool {
	return a.Total() > openingTotal
}

// ValueComposition holds the proportional c / v / s decomposition of a partial sale.
// ConstantShare + VariableShare + SurplusShare == the sale's RealisedPence.
type ValueComposition struct {
	ConstantShare Pence `json:"constant_share"`
	VariableShare Pence `json:"variable_share"`
	SurplusShare  Pence `json:"surplus_share"`
}

// Total returns ConstantShare + VariableShare + SurplusShare.
func (v ValueComposition) Total() Pence {
	return v.ConstantShare + v.VariableShare + v.SurplusShare
}

// SuccessivePartialSale is one tranche of the piecemeal realisation of C′.
// Marx's worked example: 10,000 lbs yarn sold in three tranches whose
// Decomposition sums to the full C′ (c=£372, v=£50, s=£78).
type SuccessivePartialSale struct {
	CommodityCircuitID CommodityCircuitID `json:"commodity_circuit_id"`
	Quantity           int64              `json:"quantity"`
	RealisedPence      Pence              `json:"realised_pence"`
	Decomposition      ValueComposition   `json:"decomposition"`
	SoldAt             time.Time          `json:"sold_at"`
}

// MeansOfProductionSource records the external dependency for M—C(L+MP).
// C (= L + MP) must come from outside this circuit — another producer, a merchant,
// or an import — because this circuit's opening C′ already contains produced value.
type MeansOfProductionSource struct {
	CommodityCircuitID       CommodityCircuitID          `json:"commodity_circuit_id"`
	SourceCommodityCircuitID string                      `json:"source_commodity_circuit_id,omitempty"`
	SourceKind               MeansOfProductionSourceKind `json:"source_kind"`
}

func (s MeansOfProductionSource) Validate() error {
	switch s.SourceKind {
	case SourceProducerCircuit, SourceMerchantHolding, SourceImport:
		return nil
	default:
		return ErrUnsourcedMeansOfProduction
	}
}

// CommodityCircuit is the C′—M′—C…P…C′ form of the capital circuit.
// Capital opens and closes in commodity-form; surplus-value is present at the
// starting extreme, distinguishing this form from Form I (M…M′) and Form II (P…P).
type CommodityCircuit struct {
	ID                CommodityCircuitID        `json:"id"`
	AgentID           string                    `json:"agent_id,omitempty"`
	Initial           OpeningCommodityCapital   `json:"initial"`
	Mode              ReproductionMode          `json:"mode"`
	IsFirstInvestment bool                      `json:"is_first_investment"`
	SocialCapitalLens SocialCapitalLens         `json:"social_capital_lens"`
	PartialSales      []SuccessivePartialSale   `json:"partial_sales"`
	MPSources         []MeansOfProductionSource `json:"mp_sources"`
	Terminal          *CommodityAugmented       `json:"terminal,omitempty"`
	CreatedAt         time.Time                 `json:"created_at"`
}

func (c CommodityCircuit) Validate() error {
	if !c.IsFirstInvestment && c.Initial.SurplusPence <= 0 {
		return ErrNotCommodityCapital
	}
	if c.Mode != ReproductionSimple && c.Mode != ReproductionExtended && c.Mode != ReproductionMixed {
		return errors.New("circulation: mode must be simple, extended, or mixed")
	}
	return nil
}

// EnsureMaterialAdequacy returns ErrInsufficientMaterialBase when capitalisation
// is requested but no MP sources have been linked to provide the physical elements.
func (c CommodityCircuit) EnsureMaterialAdequacy(capitalisedAmount Pence) error {
	if capitalisedAmount > 0 && len(c.MPSources) == 0 {
		return ErrInsufficientMaterialBase
	}
	return nil
}

// ComputeDecomposition derives the proportional c/v/s breakdown of a partial sale
// from the circuit's opening ratios. SurplusShare absorbs integer rounding so the
// three components sum exactly to realisedPence.
func (c CommodityCircuit) ComputeDecomposition(quantity int64, realisedPence Pence) ValueComposition {
	total := c.Initial.PoundsTotal
	if total == 0 {
		return ValueComposition{SurplusShare: realisedPence}
	}
	cShare := Pence(int64(c.Initial.ConstantPence) * quantity / total)
	vShare := Pence(int64(c.Initial.VariablePence) * quantity / total)
	sShare := realisedPence - cShare - vShare
	return ValueComposition{
		ConstantShare: cShare,
		VariableShare: vShare,
		SurplusShare:  sShare,
	}
}
