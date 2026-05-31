// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 26 ("Accumulation of Money-Capital and Its Influence
// on the Rate of Interest").
//
// As money-capital accumulates during post-crisis inactivity, loanable funds
// become abundant while demand for them is low, driving interest rates down.
// When revival begins, idle industrial capital seeks loanable outlets and
// demand rises. By prosperity, credit is stretched and the rate climbs above
// the average. In overproduction the rate peaks; in crisis it spikes sharply
// as every capitalist scrambles for means of payment.
//
// Key invariants:
//   - CycleObservation.Phase ∈ {1,2,3,4,5} (IndustrialCyclePhase.IsValid()).
//   - CycleObservation.InterestRateBP >= 0.
//   - IdleCapital.Amount >= 0.
//   - LoanableSupply.Amount >= 0.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 26 sentinel errors.
var (
	// ErrNegativeInterestRate is returned when interest_rate_bp < 0.
	ErrNegativeInterestRate = errors.New("credit: interest_rate_bp must be >= 0")

	// ErrNegativeIdleAmount is returned when idle_capital.amount < 0 or
	// loanable_supply.amount < 0.
	ErrNegativeIdleAmount = errors.New("credit: amount must be >= 0")
)

// ── Value objects ─────────────────────────────────────────────────────────────

// LoanableCapitalSupply records the aggregate supply of loanable money-capital
// at the observed phase of the industrial cycle.
type LoanableCapitalSupply struct {
	Amount   int64 `json:"amount"`   // pence; >= 0
	Abundant bool  `json:"abundant"` // true when supply exceeds demand
}

// InterestCycleObservation records an interest-rate reading tied to a phase of
// the industrial cycle and a named historical period.
type InterestCycleObservation struct {
	Phase          IndustrialCyclePhase `json:"phase"`            // 1..5
	InterestRateBP int64                `json:"interest_rate_bp"` // bp; >= 0
	Period         string               `json:"period"`
}

// IdleIndustrialCapital is the portion of industrial capital that is
// temporarily idle — sitting as money-capital waiting for profitable
// employment. Its accumulation swells the supply of loanable funds.
type IdleIndustrialCapital struct {
	Amount      int64  `json:"amount"` // pence; >= 0
	Description string `json:"description"`
}

// AccumulationGap is the signed difference between the rate of accumulation of
// money-capital and the rate of accumulation of productive capital. A positive
// gap means money-capital grows faster; a negative gap means productive
// capital outpaces loanable supply (as in crisis).
type AccumulationGap struct {
	GapBP       int64  `json:"gap_bp"` // bp; may be positive or negative
	Description string `json:"description"`
}

// ── MoneyCapitalAccumulationID ────────────────────────────────────────────────

// MoneyCapitalAccumulationID identifies a persisted MoneyCapitalAccumulation.
type MoneyCapitalAccumulationID string

// NewMoneyCapitalAccumulationID returns a fresh random MoneyCapitalAccumulationID.
func NewMoneyCapitalAccumulationID() MoneyCapitalAccumulationID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return MoneyCapitalAccumulationID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty MoneyCapitalAccumulationID.
func (id MoneyCapitalAccumulationID) IsZero() bool { return id == "" }

// ── MoneyCapitalAccumulation ──────────────────────────────────────────────────

// MoneyCapitalAccumulation is the persisted entity for an observation of
// money-capital accumulation across the industrial cycle (Vol. III Ch. 26).
// It records the cycle phase, prevailing interest rate, the supply of loanable
// funds, idle industrial capital, and the accumulation gap.
//
// Invariants:
//
//	CycleObservation.Phase.IsValid()
//	CycleObservation.InterestRateBP >= 0
//	IdleCapital.Amount >= 0
//	LoanableSupply.Amount >= 0
type MoneyCapitalAccumulation struct {
	ID               MoneyCapitalAccumulationID `json:"id"`
	Name             string                     `json:"name"`
	CycleObservation InterestCycleObservation   `json:"cycle_observation"`
	LoanableSupply   LoanableCapitalSupply      `json:"loanable_supply"`
	IdleCapital      IdleIndustrialCapital      `json:"idle_capital"`
	Gap              AccumulationGap            `json:"gap"`
	Description      string                     `json:"description"`
	CreatedAt        time.Time                  `json:"created_at"`
}

// NewMoneyCapitalAccumulation constructs a MoneyCapitalAccumulation, enforcing
// invariants. ID and CreatedAt are left zero — the store assigns them.
//
//	!cycleObs.Phase.IsValid()  → ErrInvalidCyclePhase
//	cycleObs.InterestRateBP<0 → ErrNegativeInterestRate
//	idleCapital.Amount<0      → ErrNegativeIdleAmount
//	loanableSupply.Amount<0   → ErrNegativeIdleAmount
func NewMoneyCapitalAccumulation(
	name string,
	cycleObs InterestCycleObservation,
	loanableSupply LoanableCapitalSupply,
	idleCapital IdleIndustrialCapital,
	gap AccumulationGap,
	description string,
) (MoneyCapitalAccumulation, error) {
	if !cycleObs.Phase.IsValid() {
		return MoneyCapitalAccumulation{}, ErrInvalidCyclePhase
	}
	if cycleObs.InterestRateBP < 0 {
		return MoneyCapitalAccumulation{}, ErrNegativeInterestRate
	}
	if idleCapital.Amount < 0 {
		return MoneyCapitalAccumulation{}, ErrNegativeIdleAmount
	}
	if loanableSupply.Amount < 0 {
		return MoneyCapitalAccumulation{}, ErrNegativeIdleAmount
	}
	return MoneyCapitalAccumulation{
		Name:             name,
		CycleObservation: cycleObs,
		LoanableSupply:   loanableSupply,
		IdleCapital:      idleCapital,
		Gap:              gap,
		Description:      description,
	}, nil
}

// ClassifyAccumulation returns a descriptive string for the accumulation
// dynamic at the given cycle phase and loanable-capital supply.
func ClassifyAccumulation(obs InterestCycleObservation, _ LoanableCapitalSupply) string {
	switch obs.Phase {
	case PhaseInactivity:
		return "post-crisis idle: loanable money abundant, low demand, depressed rate"
	case PhaseRevival:
		return "revival: loanable money still abundant, rising demand, moderate rate"
	case PhaseProsperity:
		return "prosperity: credit stretched, rate rising above average"
	case PhaseOverproduction:
		return "overproduction: credit fully extended, rate high"
	case PhaseCrisis:
		return "crisis: acute demand for means of payment, rate peaks"
	default:
		return "unknown phase"
	}
}
