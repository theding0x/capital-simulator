// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 30 ("Money-Capital and Real Capital. I").
//
// The central question of Part V's closing chapters: does the accumulation of
// money-capital (loanable funds) correspond to the accumulation of real capital
// (productive capital plus commodities)? Marx's answer is: not necessarily. The
// two can and do diverge — money-capital can swell while real capital stagnates
// or contracts, and vice versa.
//
// Commercial credit — the bills of exchange that producers draw on one another,
// tied directly to the rhythm of reproduction — is the FOUNDATION of the whole
// credit system. In a crisis, blocked returns (commodities that will not sell,
// so C→M never completes) and the collapse of mutual confidence contract
// commercial credit regardless of how much money-capital is nominally available.
// The ultimate cause of real crises remains the restricted consumption of the
// masses set against the drive of production to expand without limit.
//
// Key invariants:
//   - RealCapitalAccumulation.Phase ∈ {1,2,3,4,5} (IndustrialCyclePhase.IsValid()).
//   - CreditLimit.ReserveCapital >= 0.
//   - CommercialCreditCircuit.TotalBillsAmount > 0.
//   - ReturnFlow.ExpectedDays >= 0 and ReturnFlow.ActualDays >= 0.
//
// IndustrialCyclePhase and ErrInvalidCyclePhase are reused from
// rate_of_interest.go (Ch. 22) — they are not redeclared here.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 30 sentinel errors.
var (
	// ErrNegativeReserveCapital is returned when credit_limit.reserve_capital < 0.
	ErrNegativeReserveCapital = errors.New("credit: reserve_capital must be >= 0")

	// ErrNonPositiveBills is returned when commercial_credit_circuit.total_bills_amount <= 0.
	ErrNonPositiveBills = errors.New("credit: total_bills_amount must be > 0")

	// ErrNegativeReturnDays is returned when a return-flow day count is < 0.
	ErrNegativeReturnDays = errors.New("credit: return-flow days must be >= 0")
)

// ── Value objects ─────────────────────────────────────────────────────────────

// AccumulationCorrespondence records whether the accumulation of money-capital
// (loanable funds) tracks the accumulation of real (productive) capital. The two
// growth rates MAY diverge — there is no equality invariant. A positive
// money-capital growth alongside a negative real-capital growth is exactly the
// crisis configuration Marx describes.
type AccumulationCorrespondence struct {
	MoneyCapitalGrowthBP int64  `json:"money_capital_growth_bp"` // bp; may be positive or negative
	RealCapitalGrowthBP  int64  `json:"real_capital_growth_bp"`  // bp; may be positive or negative
	Corresponds          bool   `json:"corresponds"`             // true when the two accumulations align
	Explanation          string `json:"explanation"`
}

// CreditLimit is the reserve of money-capital that sets the elastic limit of the
// credit system. In normal times credit expands far beyond this reserve; in a
// crisis the reserve becomes the hard floor and credit contracts back toward it.
type CreditLimit struct {
	ReserveCapital int64  `json:"reserve_capital"` // pence; >= 0
	Description    string `json:"description"`
}

// ReturnFlow describes the C→M realisation that must complete before a bill of
// exchange can be settled. When ActualDays exceeds ExpectedDays the return is
// "blocked" — the commodity has not sold on time, and the bill it backs becomes
// a demand for means of payment that cannot be met from sales.
type ReturnFlow struct {
	ExpectedDays int64 `json:"expected_days"` // days; >= 0
	ActualDays   int64 `json:"actual_days"`   // days; >= 0
}

// NewReturnFlow builds and validates a ReturnFlow.
//
//	expectedDays < 0 || actualDays < 0 → ErrNegativeReturnDays
func NewReturnFlow(expectedDays, actualDays int64) (ReturnFlow, error) {
	if expectedDays < 0 || actualDays < 0 {
		return ReturnFlow{}, ErrNegativeReturnDays
	}
	return ReturnFlow{ExpectedDays: expectedDays, ActualDays: actualDays}, nil
}

// IsBlocked reports whether the actual return ran longer than expected — the
// signature of a crisis in which commodities pile up unsold and bills fall due
// against returns that never arrive.
func (r ReturnFlow) IsBlocked() bool {
	return r.ActualDays > r.ExpectedDays
}

// ── CommercialCreditCircuitID ─────────────────────────────────────────────────

// CommercialCreditCircuitID identifies a CommercialCreditCircuit. The circuit is
// a value type — it is not persisted — but it carries an ID so it can be
// referenced and round-tripped through the API like the other domain types.
type CommercialCreditCircuitID string

// NewCommercialCreditCircuitID returns a fresh random CommercialCreditCircuitID.
func NewCommercialCreditCircuitID() CommercialCreditCircuitID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return CommercialCreditCircuitID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CommercialCreditCircuitID.
func (id CommercialCreditCircuitID) IsZero() bool { return id == "" }

// ── CommercialCreditCircuit ───────────────────────────────────────────────────

// CommercialCreditCircuit is the chain of bills of exchange that producers draw
// on one another as a commodity passes down the line of production — spinner →
// weaver → manufacturer → exporter. Each stage advances credit to the next; the
// whole chain is the foundation of the credit system and is tied directly to the
// reproduction process. It is a value type, never persisted.
//
// Invariant:
//
//	TotalBillsAmount > 0
type CommercialCreditCircuit struct {
	ID               CommercialCreditCircuitID `json:"id"`
	Stages           []string                  `json:"stages"`             // e.g. spinner→weaver→manufacturer→exporter
	TotalBillsAmount int64                     `json:"total_bills_amount"` // pence; > 0
	Description      string                    `json:"description"`
}

// NewCommercialCreditCircuit builds and validates a CommercialCreditCircuit.
// The ID is left zero; callers may assign one with NewCommercialCreditCircuitID.
//
//	totalBillsAmount <= 0 → ErrNonPositiveBills
func NewCommercialCreditCircuit(stages []string, totalBillsAmount int64, description string) (CommercialCreditCircuit, error) {
	if totalBillsAmount <= 0 {
		return CommercialCreditCircuit{}, ErrNonPositiveBills
	}
	return CommercialCreditCircuit{
		Stages:           stages,
		TotalBillsAmount: totalBillsAmount,
		Description:      description,
	}, nil
}

// StageCount returns the number of stages in the credit chain.
func (c CommercialCreditCircuit) StageCount() int { return len(c.Stages) }

// ── RealCapitalAccumulationID ─────────────────────────────────────────────────

// RealCapitalAccumulationID identifies a persisted RealCapitalAccumulation.
type RealCapitalAccumulationID string

// NewRealCapitalAccumulationID returns a fresh random RealCapitalAccumulationID.
func NewRealCapitalAccumulationID() RealCapitalAccumulationID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return RealCapitalAccumulationID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty RealCapitalAccumulationID.
func (id RealCapitalAccumulationID) IsZero() bool { return id == "" }

// ── RealCapitalAccumulation ───────────────────────────────────────────────────

// RealCapitalAccumulation is the persisted entity for an observation of whether
// money-capital accumulation corresponds to real-capital accumulation at a given
// phase of the industrial cycle (Vol. III Ch. 30). It records the correspondence
// between the two accumulations, the credit limit (reserve), whether commercial
// credit is contracting, and whether the cause of any contraction is a scarcity
// of money (Marx says: usually not — it is blocked returns and lost confidence).
//
// Invariants:
//
//	Phase.IsValid()
//	CreditLimit.ReserveCapital >= 0
type RealCapitalAccumulation struct {
	ID                        RealCapitalAccumulationID  `json:"id"`
	Name                      string                     `json:"name"`
	Phase                     IndustrialCyclePhase       `json:"phase"` // 1..5
	Correspondence            AccumulationCorrespondence `json:"accumulation_correspondence"`
	CreditLimit               CreditLimit                `json:"credit_limit"`
	CommercialCreditContracts bool                       `json:"commercial_credit_contracts"`
	CauseIsMoneyScarcity      bool                       `json:"cause_is_money_scarcity"`
	Description               string                     `json:"description"`
	CreatedAt                 time.Time                  `json:"created_at"`
}

// NewRealCapitalAccumulation constructs a RealCapitalAccumulation, enforcing
// invariants. ID and CreatedAt are left zero — the store assigns them.
//
//	!phase.IsValid()                  → ErrInvalidCyclePhase
//	creditLimit.ReserveCapital < 0    → ErrNegativeReserveCapital
//
// Fixture (Marx, Vol. III Ch. 30):
//
//	Crisis 1847: phase PhaseCrisis, commercial credit contracts, cause is NOT
//	money scarcity (it is blocked returns + lost confidence), money-capital and
//	real-capital growth diverge (Corresponds=false).
func NewRealCapitalAccumulation(
	name string,
	phase IndustrialCyclePhase,
	correspondence AccumulationCorrespondence,
	creditLimit CreditLimit,
	commercialCreditContracts bool,
	causeIsMoneyScarcity bool,
	description string,
) (RealCapitalAccumulation, error) {
	if !phase.IsValid() {
		return RealCapitalAccumulation{}, ErrInvalidCyclePhase
	}
	if creditLimit.ReserveCapital < 0 {
		return RealCapitalAccumulation{}, ErrNegativeReserveCapital
	}
	return RealCapitalAccumulation{
		Name:                      name,
		Phase:                     phase,
		Correspondence:            correspondence,
		CreditLimit:               creditLimit,
		CommercialCreditContracts: commercialCreditContracts,
		CauseIsMoneyScarcity:      causeIsMoneyScarcity,
		Description:               description,
	}, nil
}

// CorrespondenceVerdict returns a descriptive string for whether the
// accumulation of money-capital aligns with the accumulation of real capital.
// When the observation records a contraction of commercial credit, the verdict
// also names the cause (money scarcity vs. blocked returns / lost confidence),
// echoing Marx's insistence that the credit contraction in a crisis is rarely a
// mere shortage of currency.
func (a RealCapitalAccumulation) CorrespondenceVerdict() string {
	if a.CommercialCreditContracts {
		if a.CauseIsMoneyScarcity {
			return "commercial credit contracting — attributed to a scarcity of money-capital"
		}
		return "commercial credit contracting — cause is blocked returns and lost confidence, not a money shortage"
	}
	if a.Correspondence.Corresponds {
		return "money-capital and real-capital accumulation correspond"
	}
	return "money-capital and real-capital accumulation diverge"
}
