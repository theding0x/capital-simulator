package circulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Pence is the Vol. II monetary unit: 100 pence = £1.
type Pence int64

// Errors returned by domain operations.
var (
	ErrIndivisibleProductionElement = errors.New("circulation: amount below minimum capitalisation increment")
	ErrInsufficientReserve          = errors.New("circulation: reserve draw exceeds balance")
	ErrPartialRealisation           = errors.New("circulation: commodity capital not fully realised")
	ErrValueRelationShift           = errors.New("circulation: value-relation shift detected mid-circuit")
)

// ReproductionMode describes how surplus-value is handled each circuit.
type ReproductionMode string

const (
	ReproductionSimple   ReproductionMode = "simple"
	ReproductionExtended ReproductionMode = "extended"
	ReproductionMixed    ReproductionMode = "mixed"
)

// ReserveDrawReason explains why a reserve was drawn.
type ReserveDrawReason string

const (
	ReserveDrawDelayedRealisation ReserveDrawReason = "delayed_realisation"
	ReserveDrawPriceRiseInMP      ReserveDrawReason = "price_rise_in_mp"
)

// ProductiveCircuitID is a 96-bit hex identifier.
type ProductiveCircuitID string

func (id ProductiveCircuitID) IsZero() bool { return id == "" }

// NewProductiveCircuitID generates a 96-bit hex identifier via crypto/rand.
func NewProductiveCircuitID() ProductiveCircuitID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return ProductiveCircuitID(hex.EncodeToString(b))
}

// ProductiveCircuit is the P…C'—M'—C…P form of the capital circuit.
// It starts and ends in production; circulation is the connecting link between
// two functionings of productive capital.
type ProductiveCircuit struct {
	ID                         ProductiveCircuitID  `json:"id"`
	AgentID                    string               `json:"agent_id,omitempty"`
	ConstantPence              Pence                `json:"constant_pence"`
	VariablePence              Pence                `json:"variable_pence"`
	Mode                       ReproductionMode     `json:"mode"`
	MinCapitalisationIncrement Pence                `json:"min_capitalisation_increment_pence"`
	LatentMoneyCapital         LatentMoneyCapital   `json:"latent_money_capital"`
	ReserveFund                ReserveFund          `json:"reserve_fund"`
	RevenueCircuits            []RevenueCircuit     `json:"revenue_circuits"`
	CapitalisationSteps        []CapitalisationStep `json:"capitalisation_steps"`
	ReserveDraws               []ReserveDraw        `json:"reserve_draws"`
	CreatedAt                  time.Time            `json:"created_at"`
}

// TotalAdvance returns the active money-capital (C+V) currently deployed in production.
func (p ProductiveCircuit) TotalAdvance() Pence {
	return p.ConstantPence + p.VariablePence
}

func (p ProductiveCircuit) Validate() error {
	if p.ConstantPence <= 0 {
		return errors.New("circulation: constant_pence must be positive")
	}
	if p.VariablePence <= 0 {
		return errors.New("circulation: variable_pence must be positive")
	}
	if p.Mode != ReproductionSimple && p.Mode != ReproductionExtended && p.Mode != ReproductionMixed {
		return errors.New("circulation: mode must be simple, extended, or mixed")
	}
	if p.MinCapitalisationIncrement <= 0 {
		return errors.New("circulation: min_capitalisation_increment_pence must be positive")
	}
	return nil
}

// RevenueCircuit records the c—m—c personal-consumption exit from C'—M'.
// This circuit falls outside the capital circuit entirely.
type RevenueCircuit struct {
	ProductiveCircuitID ProductiveCircuitID `json:"productive_circuit_id"`
	Amount              Pence               `json:"amount"`
	SpentAt             time.Time           `json:"spent_at"`
}

// LatentMoneyCapital is the hoarded m awaiting the minimum-capitalisation threshold.
// While hoarded it does not function as capital and earns no surplus-value.
type LatentMoneyCapital struct {
	ProductiveCircuitID ProductiveCircuitID `json:"productive_circuit_id"`
	Accumulated         Pence               `json:"accumulated"`
	Threshold           Pence               `json:"threshold"`
}

// IsArmed reports whether accumulated m has reached the minimum-capitalisation threshold.
func (l LatentMoneyCapital) IsArmed() bool { return l.Accumulated >= l.Threshold }

// CapitalisationStep records the moment a hoard crosses the threshold and is
// injected into the circuit as additional constant + variable capital.
// DeltaConstantPence + DeltaVariablePence must equal AmountInjected.
type CapitalisationStep struct {
	ProductiveCircuitID ProductiveCircuitID `json:"productive_circuit_id"`
	AmountInjected      Pence               `json:"amount_injected"`
	DeltaConstantPence  Pence               `json:"delta_constant_pence"`
	DeltaVariablePence  Pence               `json:"delta_variable_pence"`
	OccurredAt          time.Time           `json:"occurred_at"`
}

// ReserveFund is a part of the accumulation-fund held to absorb circulation
// disturbances. It is funded from accumulated surplus-value, never from the
// active money-capital M.
type ReserveFund struct {
	ProductiveCircuitID ProductiveCircuitID `json:"productive_circuit_id"`
	Balance             Pence               `json:"balance"`
}

// ReserveDraw records a single withdrawal from the reserve fund.
type ReserveDraw struct {
	ProductiveCircuitID ProductiveCircuitID `json:"productive_circuit_id"`
	DrawnPence          Pence               `json:"drawn_pence"`
	Reason              ReserveDrawReason   `json:"reason"`
	OccurredAt          time.Time           `json:"occurred_at"`
}
