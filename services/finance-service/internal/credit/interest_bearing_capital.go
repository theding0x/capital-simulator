// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 21 ("Interest-Bearing Capital").
//
// Interest-bearing capital is money-capital advanced to a functioning
// capitalist and returned with interest: the M—M′ circuit. Interest is a
// portion of the average profit on the borrowed capital paid by the borrower
// (functioning capitalist) to the lender (money capitalist). The crucial
// Marxian invariant is that interest-bearing capital creates no new value —
// it is a redistribution of already-produced surplus-value.
//
// Key invariants:
//   - InterestEarned = roundHalfUp(MoneyAdvanced × InterestRateBP, basisPoints).
//   - MoneyReturned == MoneyAdvanced + InterestEarned.
//   - NewValueCreated == 0 on every interest-bearing-capital record.
//   - Interest.Value < averageProfitOnCapital (interest is a portion of profit,
//     not profit itself; the functioning capitalist retains profit of enterprise).
//   - LoanTermDays > 0; MoneyAdvanced > 0; InterestRateBP > 0.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// basisPoints is the fixed-point scale used throughout the credit package:
// 1 bp = 0.01%; 10000 bp = 100%.
const basisPoints = int64(10000)

// Sentinel errors for the credit package.
var (
	// ErrNonPositiveMoney is returned when money_advanced is ≤ 0.
	ErrNonPositiveMoney = errors.New("credit: money_advanced must be positive")

	// ErrNonPositiveRate is returned when interest_rate_bp is ≤ 0.
	ErrNonPositiveRate = errors.New("credit: interest_rate_bp must be positive")

	// ErrNonPositiveTerm is returned when loan_term_days is ≤ 0.
	ErrNonPositiveTerm = errors.New("credit: loan_term_days must be positive")

	// ErrInterestExceedsProfit is returned when interest ≥ average profit on
	// capital. Marx's invariant: interest is a part of profit, so it must be
	// strictly less than the total profit available.
	ErrInterestExceedsProfit = errors.New("credit: interest must be less than average profit on capital")
)

// roundHalfUp computes (numer + denom/2) / denom — the canonical round-half-up
// division used throughout the simulator for basis-point arithmetic.
func roundHalfUp(numer, denom int64) int64 {
	return (numer + denom/2) / denom
}

// ── InterestBearingCapitalID ──────────────────────────────────────────────────

// InterestBearingCapitalID identifies a persisted interest-bearing-capital
// record. 96-bit hex from crypto/rand, matching every other domain ID in the
// simulator.
type InterestBearingCapitalID string

// NewInterestBearingCapitalID returns a fresh random InterestBearingCapitalID.
func NewInterestBearingCapitalID() InterestBearingCapitalID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return InterestBearingCapitalID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty InterestBearingCapitalID.
func (id InterestBearingCapitalID) IsZero() bool { return id == "" }

// ── InterestBearingCapital ────────────────────────────────────────────────────

// InterestBearingCapital is the persisted entity for an interest-bearing-capital
// computation (Vol. III Ch. 21). The money-capitalist (lender) advances
// MoneyAdvanced at InterestRateBP for LoanTermDays. On return, the functioning
// capitalist repays MoneyAdvanced + InterestEarned. The lender creates no new
// value — the interest is a portion of the average profit on the advanced capital.
//
//	InterestEarned  = roundHalfUp(MoneyAdvanced × InterestRateBP, basisPoints)
//	MoneyReturned   = MoneyAdvanced + InterestEarned
//	NewValueCreated = 0 always
type InterestBearingCapital struct {
	ID              InterestBearingCapitalID `json:"id"`
	MoneyAdvanced   int64                    `json:"money_advanced"`    // pence; > 0
	InterestRateBP  int64                    `json:"interest_rate_bp"`  // bp; > 0
	LoanTermDays    int64                    `json:"loan_term_days"`    // > 0
	InterestEarned  int64                    `json:"interest_earned"`   // pence; derived
	MoneyReturned   int64                    `json:"money_returned"`    // pence; derived
	NewValueCreated int64                    `json:"new_value_created"` // always 0
	CreatedAt       time.Time                `json:"created_at"`
}

// NewInterestBearingCapital constructs an InterestBearingCapital from the
// money advanced, the annual interest rate in basis points, and the loan term
// in days. No ID or timestamp is assigned — the store does that on persistence.
//
//	moneyAdvanced  ≤ 0 → ErrNonPositiveMoney
//	interestRateBP ≤ 0 → ErrNonPositiveRate
//	loanTermDays   ≤ 0 → ErrNonPositiveTerm
//
// Fixture (Marx, Vol. III Ch. 21):
//
//	£1,000 at 5% for one year:
//	moneyAdvanced=100000 pence, interestRateBP=500, loanTermDays=365
//	InterestEarned = roundHalfUp(100000*500, 10000) = roundHalfUp(50000000, 10000) = 5000
//	MoneyReturned  = 100000 + 5000 = 105000
func NewInterestBearingCapital(moneyAdvanced, interestRateBP, loanTermDays int64) (InterestBearingCapital, error) {
	if moneyAdvanced <= 0 {
		return InterestBearingCapital{}, ErrNonPositiveMoney
	}
	if interestRateBP <= 0 {
		return InterestBearingCapital{}, ErrNonPositiveRate
	}
	if loanTermDays <= 0 {
		return InterestBearingCapital{}, ErrNonPositiveTerm
	}
	interest := roundHalfUp(moneyAdvanced*interestRateBP, basisPoints)
	return InterestBearingCapital{
		MoneyAdvanced:   moneyAdvanced,
		InterestRateBP:  interestRateBP,
		LoanTermDays:    loanTermDays,
		InterestEarned:  interest,
		MoneyReturned:   moneyAdvanced + interest,
		NewValueCreated: 0,
	}, nil
}

// ── MoneyCapitalCircuit ────────────────────────────────────────────────────────

// MoneyCapitalCircuit is the value object for the M—M′ circuit of
// interest-bearing capital. The money-capitalist advances MoneyAdvanced and
// recovers MoneyReturned (= MoneyAdvanced + InterestEarned). Computed on
// demand, not persisted.
//
//	Invariant: MoneyReturned == MoneyAdvanced + InterestEarned
type MoneyCapitalCircuit struct {
	MoneyAdvanced  int64 `json:"money_advanced"`  // pence; > 0
	InterestEarned int64 `json:"interest_earned"` // pence; >= 0
	MoneyReturned  int64 `json:"money_returned"`  // pence; = MoneyAdvanced + InterestEarned
}

// NewMoneyCapitalCircuit builds and validates a MoneyCapitalCircuit.
//
//	moneyAdvanced  ≤ 0                             → ErrNonPositiveMoney
//	moneyReturned != moneyAdvanced+interestEarned  → error describing the invariant
func NewMoneyCapitalCircuit(moneyAdvanced, interestEarned, moneyReturned int64) (MoneyCapitalCircuit, error) {
	if moneyAdvanced <= 0 {
		return MoneyCapitalCircuit{}, ErrNonPositiveMoney
	}
	if moneyReturned != moneyAdvanced+interestEarned {
		return MoneyCapitalCircuit{}, errors.New("credit: MoneyReturned must equal MoneyAdvanced + InterestEarned")
	}
	return MoneyCapitalCircuit{
		MoneyAdvanced:  moneyAdvanced,
		InterestEarned: interestEarned,
		MoneyReturned:  moneyReturned,
	}, nil
}

// ── LoanableCapital ────────────────────────────────────────────────────────────

// LoanableCapital is the value object for the pool of money-capital available
// for lending. Computed on demand, not persisted.
type LoanableCapital struct {
	Amount int64 `json:"amount"` // pence; > 0
}

// NewLoanableCapital builds and validates a LoanableCapital.
// amount ≤ 0 → ErrNonPositiveMoney.
func NewLoanableCapital(amount int64) (LoanableCapital, error) {
	if amount <= 0 {
		return LoanableCapital{}, ErrNonPositiveMoney
	}
	return LoanableCapital{Amount: amount}, nil
}

// ── Interest ──────────────────────────────────────────────────────────────────

// Interest is the value object for the interest payment — the portion of
// average profit transferred from the functioning capitalist to the money
// capitalist. Interest must be strictly less than average profit; the remainder
// is the profit of enterprise retained by the borrower.
type Interest struct {
	Value int64 `json:"value"` // pence; >= 0
}

// NewInterest builds and validates an Interest.
//
//	value >= averageProfitOnCapital → ErrInterestExceedsProfit
//
// The invariant reflects Marx: interest is a part (Bruchteil) of profit, not
// profit itself. The functioning capitalist must retain some profit of enterprise
// or the separation of ownership from function collapses.
func NewInterest(value, averageProfitOnCapital int64) (Interest, error) {
	if value >= averageProfitOnCapital {
		return Interest{}, ErrInterestExceedsProfit
	}
	return Interest{Value: value}, nil
}

// ── LoanTerm ──────────────────────────────────────────────────────────────────

// LoanTerm is the value object for the duration of a loan in days.
type LoanTerm struct {
	Days int64 `json:"days"` // > 0
}

// NewLoanTerm builds and validates a LoanTerm.
// days ≤ 0 → ErrNonPositiveTerm.
func NewLoanTerm(days int64) (LoanTerm, error) {
	if days <= 0 {
		return LoanTerm{}, ErrNonPositiveTerm
	}
	return LoanTerm{Days: days}, nil
}

// ── MoneyCapitalist ───────────────────────────────────────────────────────────

// MoneyCapitalist represents the passive lender who owns the money-capital
// and earns interest by placing it at the disposal of functioning capitalists.
// The money capitalist does not participate in the production process; their
// return is interest, a portion of the surplus-value produced elsewhere.
type MoneyCapitalist struct {
	// Role describes the money-capitalist's position in the circuit.
	Role string `json:"role"` // e.g. "lender", "rentier"
	// EarnsInterest flags that the money-capitalist's income is interest,
	// not profit of enterprise.
	EarnsInterest bool `json:"earns_interest"`
}

// ── FunctioningCapitalist ─────────────────────────────────────────────────────

// FunctioningCapitalist represents the active borrower who employs the
// borrowed capital in production or trade and earns the profit of enterprise
// (the residual after paying interest to the money-capitalist).
type FunctioningCapitalist struct {
	// Role describes the functioning capitalist's position in the circuit.
	Role string `json:"role"` // e.g. "industrial capitalist", "commercial capitalist"
	// EarnsInterest is always false — the functioning capitalist earns profit
	// of enterprise, not interest.
	EarnsInterest bool `json:"earns_interest"`
}
