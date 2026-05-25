package circulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Vol. II Ch. 1 — The Circuit of Money-Capital.

var (
	// ErrLabourMeansMismatch is returned when MeansCapacityHours ≠ LabourPowerHours.
	// Marx §I: "If the means of production at hand were insufficient, the excess labour
	// at the disposal of the purchaser could not be utilised."
	ErrLabourMeansMismatch = errors.New("circulation: means_capacity_hours must equal labour_power_hours")

	// ErrNoProductionPhase is returned when a commodity-capital record is attempted
	// without an intervening productive phase. Surplus-value arises only in P.
	ErrNoProductionPhase = errors.New("circulation: surplus-value requires a productive phase between M—C and C′")

	// ErrInvalidPhaseTransition is returned on an attempt to skip a phase or
	// re-enter a completed phase.
	ErrInvalidPhaseTransition = errors.New("circulation: invalid phase transition or re-entry of a completed phase")

	// ErrMagnitudeNotPreserved is returned when the advance magnitude is not
	// preserved across a form-change. Marx §I: "M is the same capital-value as P,
	// only it has a different mode of existence."
	ErrMagnitudeNotPreserved = errors.New("circulation: advance magnitude not preserved across form-change")
)

// MoneyCircuitID is a 96-bit hex identifier.
type MoneyCircuitID string

func (id MoneyCircuitID) IsZero() bool { return id == "" }

// NewMoneyCircuitID generates a 96-bit hex identifier via crypto/rand.
func NewMoneyCircuitID() MoneyCircuitID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return MoneyCircuitID(hex.EncodeToString(b))
}

// CircuitMoment names each of the six moments in M—C(Lp+Mp)…P…C′—M′.
type CircuitMoment string

const (
	MomentM       CircuitMoment = "M"         // opening money-capital
	MomentMC      CircuitMoment = "M-C"       // purchase phase
	MomentP       CircuitMoment = "P"          // productive state
	MomentCPrime  CircuitMoment = "C-prime"   // commodity-capital
	MomentCMPrime CircuitMoment = "C-M-prime" // realisation phase
	MomentMPrime  CircuitMoment = "M-prime"   // closed, M′ returned
)

// phaseOrder maps each moment to its linear sequence index.
var phaseOrder = map[CircuitMoment]int{
	MomentM:       0,
	MomentMC:      1,
	MomentP:       2,
	MomentCPrime:  3,
	MomentCMPrime: 4,
	MomentMPrime:  5,
}

// MoneyAdvance is M: the opening money-capital advanced to begin the circuit.
type MoneyAdvance struct {
	Amount     Pence     `json:"amount"`
	AgentID    string    `json:"agent_id,omitempty"`
	AdvancedAt time.Time `json:"advanced_at"`
}

// LabourLeg is the M—L component of the purchase phase.
type LabourLeg struct {
	Amount           Pence `json:"amount"`
	LabourPowerHours int64 `json:"labour_power_hours"`
}

// MeansLeg is the M—MP component of the purchase phase.
type MeansLeg struct {
	Amount             Pence `json:"amount"`
	MeansCapacityHours int64 `json:"means_capacity_hours"`
}

// PurchasePhase records M—C: the money advance split into L and MP.
type PurchasePhase struct {
	MoneyCircuitID MoneyCircuitID `json:"money_circuit_id"`
	LabourLeg      LabourLeg      `json:"labour_leg"`
	MeansLeg       MeansLeg       `json:"means_leg"`
}

// Validate enforces the proportional law: means must absorb exactly the
// labour purchased. Insufficient means waste labour; excess means starve.
func (p PurchasePhase) Validate() error {
	if p.LabourLeg.LabourPowerHours != p.MeansLeg.MeansCapacityHours {
		return ErrLabourMeansMismatch
	}
	return nil
}

// Total returns LabourLeg.Amount + MeansLeg.Amount.
func (p PurchasePhase) Total() Pence {
	return p.LabourLeg.Amount + p.MeansLeg.Amount
}

// ProductiveState is P: labour-power × means of production combined and
// functioning. ConstantPence + VariablePence must equal the original advance.
type ProductiveState struct {
	MoneyCircuitID MoneyCircuitID `json:"money_circuit_id"`
	EnteredAt      time.Time      `json:"entered_at"`
	ConstantPence  Pence          `json:"constant_pence"`
	VariablePence  Pence          `json:"variable_pence"`
}

// TotalAdvance returns c + v, which must equal the original MoneyAdvance.Amount.
func (ps ProductiveState) TotalAdvance() Pence {
	return ps.ConstantPence + ps.VariablePence
}

// CommodityCapital is C′ = C + c: the commodity-form carrying surplus-value.
// ValueOriginal is the carried-over capital-value (= the original advance).
// ValueSurplus is the surplus-value crystallised in the commodity body.
type CommodityCapital struct {
	MoneyCircuitID MoneyCircuitID `json:"money_circuit_id"`
	CommodityID    string         `json:"commodity_id,omitempty"`
	ValueOriginal  Pence          `json:"value_original"`
	ValueSurplus   Pence          `json:"value_surplus"`
}

// Total returns C′ = ValueOriginal + ValueSurplus.
func (cc CommodityCapital) Total() Pence {
	return cc.ValueOriginal + cc.ValueSurplus
}

// Realisation records C′—M′: the sale that closes the circuit.
type Realisation struct {
	MoneyCircuitID MoneyCircuitID `json:"money_circuit_id"`
	SoldAt         time.Time      `json:"sold_at"`
	RealisedPence  Pence          `json:"realised_pence"`
}

// MoneyCircuit is the M—C(Lp+Mp)…P…C′—M′ form of capital in motion.
// Capital opens and closes in money-form; valorisation is M′ − M.
type MoneyCircuit struct {
	ID          MoneyCircuitID    `json:"id"`
	Advance     MoneyAdvance      `json:"advance"`
	Moment      CircuitMoment     `json:"moment"`
	Purchase    *PurchasePhase    `json:"purchase,omitempty"`
	Productive  *ProductiveState  `json:"productive,omitempty"`
	Commodity   *CommodityCapital `json:"commodity,omitempty"`
	Realisation *Realisation      `json:"realisation,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

// SurplusRealised returns M′ − M. Positive ⇒ valorising; zero ⇒ simple
// reproduction; negative ⇒ failed realisation. Returns 0 before M′.
func (mc MoneyCircuit) SurplusRealised() Pence {
	if mc.Realisation == nil {
		return 0
	}
	return mc.Realisation.RealisedPence - mc.Advance.Amount
}

// CanTransitionTo returns nil if the circuit can legally advance to next.
// Phases must appear in strict linear order; skipping or reversing is forbidden.
func (mc MoneyCircuit) CanTransitionTo(next CircuitMoment) error {
	currentRank, ok := phaseOrder[mc.Moment]
	if !ok {
		return ErrInvalidPhaseTransition
	}
	nextRank, ok := phaseOrder[next]
	if !ok {
		return ErrInvalidPhaseTransition
	}
	if nextRank != currentRank+1 {
		return ErrInvalidPhaseTransition
	}
	return nil
}
