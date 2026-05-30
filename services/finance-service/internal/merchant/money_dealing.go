// Package merchant implements Capital Vol. III, Part IV.
// This file covers Ch. 19 ("Money-Dealing Capital").
//
// Money-dealing capital is a form of capital specialised in the purely
// technical operations of money circulation — receiving and making payments,
// keeping accounts, exchanging currencies, safekeeping cash — that arise
// necessarily from the circulation of commodity-capital and the payment of
// wages. Like commercial capital, it appropriates a share of the general rate
// of profit; unlike industrial capital, it creates no new value.
//
// Key invariants:
//   - MoneyDealingProfit = roundHalfUp(moneyAdvanced × generalRateBP, basisPoints).
//   - NewValueCreated == 0 on every money-dealing capital record.
//   - The M—M′ circuit (MoneyDealingCircuit) is valid only when
//     MoneyReturned > MoneyAdvanced.
//   - OperationKinds enumerates the five technical functions; each must be 1–5.
package merchant

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// ErrInvalidOperationKind is returned when a MoneyOperationKind value is not
// in the valid range 1–5.
var ErrInvalidOperationKind = errors.New("merchant: invalid money_operation_kind — must be 1–5")

// MoneyOperationKind enumerates the five technical operations performed by
// money-dealing capital (Vol. III Ch. 19). None creates new value.
type MoneyOperationKind int

const (
	// KindReceipt — §1: collecting and registering incoming payments.
	KindReceipt MoneyOperationKind = 1
	// KindPayment — §2: executing outgoing payments on behalf of industrial
	// and commercial capitalists.
	KindPayment MoneyOperationKind = 2
	// KindExchange — §3: converting between currencies or coin denominations;
	// a technical requirement of world-market transactions.
	KindExchange MoneyOperationKind = 3
	// KindSafekeeping — §4: holding cash reserves and other monetary
	// instruments for clients.
	KindSafekeeping MoneyOperationKind = 4
	// KindBookkeeping — §5: account-keeping, ledger entries, and the clerical
	// work of recording monetary flows.
	KindBookkeeping MoneyOperationKind = 5
)

// IsValid reports whether k is a recognised MoneyOperationKind (1–5).
func (k MoneyOperationKind) IsValid() bool {
	return k >= KindReceipt && k <= KindBookkeeping
}

// ── MoneyDealingCapitalID ────────────────────────────────────────────────────

// MoneyDealingCapitalID identifies a persisted money-dealing-capital record.
// 96-bit hex from crypto/rand, matching every other domain ID in the simulator.
type MoneyDealingCapitalID string

// NewMoneyDealingCapitalID returns a fresh random MoneyDealingCapitalID.
func NewMoneyDealingCapitalID() MoneyDealingCapitalID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("merchant: rand: %w", err))
	}
	return MoneyDealingCapitalID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty MoneyDealingCapitalID.
func (id MoneyDealingCapitalID) IsZero() bool { return id == "" }

// ── MoneyDealingCapital ──────────────────────────────────────────────────────

// MoneyDealingCapital is the persisted entity for a money-dealing-capital
// computation (Vol. III Ch. 19). The money-dealer advances capital to cover
// the costs of purely technical monetary operations (receipt, payment,
// exchange, safekeeping, bookkeeping). Like all merchant capital, it
// participates in the formation of the general rate of profit and earns
// MoneyDealingProfit from that rate — but creates no new value.
//
//	MoneyDealingProfit = roundHalfUp(moneyAdvanced × generalRateBP, basisPoints)
//	NewValueCreated    = 0 always
type MoneyDealingCapital struct {
	ID                 MoneyDealingCapitalID `json:"id"`
	MoneyAdvanced      int64                 `json:"money_advanced"`       // pence; > 0
	GeneralRateBP      int64                 `json:"general_rate_bp"`      // bp; > 0
	OperationKinds     []MoneyOperationKind  `json:"operation_kinds"`      // each validated 1–5
	MoneyDealingProfit int64                 `json:"money_dealing_profit"` // pence; derived
	NewValueCreated    int64                 `json:"new_value_created"`    // always 0
	CreatedAt          time.Time             `json:"created_at"`
}

// NewMoneyDealingCapital constructs a MoneyDealingCapital from the money
// advanced, the general rate of profit in basis points, and a slice of
// operation kinds. No ID or timestamp is assigned — the store does that on
// persistence.
//
//	moneyAdvanced ≤ 0   → ErrNonPositiveMoney   (from commercial_capital.go)
//	generalRateBP ≤ 0   → ErrNonPositiveCapital (from commercial_profit.go)
//	any kind !IsValid() → ErrInvalidOperationKind
//
// Fixture verification (seed 1901):
//
//	moneyAdvanced=500000, generalRateBP=1500
//	profit = roundHalfUp(500000*1500, 10000) = roundHalfUp(750000000, 10000) = 75000
func NewMoneyDealingCapital(moneyAdvanced, generalRateBP int64, operationKinds []MoneyOperationKind) (MoneyDealingCapital, error) {
	if moneyAdvanced <= 0 {
		return MoneyDealingCapital{}, ErrNonPositiveMoney
	}
	if generalRateBP <= 0 {
		return MoneyDealingCapital{}, ErrNonPositiveCapital
	}
	for _, k := range operationKinds {
		if !k.IsValid() {
			return MoneyDealingCapital{}, ErrInvalidOperationKind
		}
	}
	profit := roundHalfUp(moneyAdvanced*generalRateBP, basisPoints)
	kinds := make([]MoneyOperationKind, len(operationKinds))
	copy(kinds, operationKinds)
	return MoneyDealingCapital{
		MoneyAdvanced:      moneyAdvanced,
		GeneralRateBP:      generalRateBP,
		OperationKinds:     kinds,
		MoneyDealingProfit: profit,
		NewValueCreated:    0,
	}, nil
}

// ── CashReserve ──────────────────────────────────────────────────────────────

// CashReserve is the value object for a cash reserve held by a money-dealer on
// behalf of clients (KindSafekeeping). The amount must be non-negative.
// Computed on demand, not persisted.
type CashReserve struct {
	Amount int64 `json:"amount"` // pence; >= 0
}

// NewCashReserve builds and validates a CashReserve.
// amount < 0 → ErrNegativeMagnitude.
func NewCashReserve(amount int64) (CashReserve, error) {
	if amount < 0 {
		return CashReserve{}, ErrNegativeMagnitude
	}
	return CashReserve{Amount: amount}, nil
}

// ── MoneyDealingCircuit ───────────────────────────────────────────────────────

// MoneyDealingCircuit is the value object for the money-dealer's M—M′ circuit.
// The money-dealer advances MoneyAdvanced and recovers MoneyReturned (which
// includes the money-dealing profit drawn from the general rate). The circuit
// is distinct from MerchantCircuit (M—C—M′): no commodity changes hands; the
// movement is purely monetary — M advances, M′ returns.
//
//	Profit() = MoneyReturned − MoneyAdvanced
//
// The circuit is valid only when MoneyReturned > MoneyAdvanced — i.e. the
// money-dealer extracts profit from the general pool of surplus-value.
// Computed on demand, not persisted.
type MoneyDealingCircuit struct {
	MoneyAdvanced int64 `json:"money_advanced"` // pence; > 0
	MoneyReturned int64 `json:"money_returned"` // pence; > MoneyAdvanced
}

// Profit is the money-dealer's gain: MoneyReturned − MoneyAdvanced.
func (c MoneyDealingCircuit) Profit() int64 { return c.MoneyReturned - c.MoneyAdvanced }

// NewMoneyDealingCircuit builds and validates a MoneyDealingCircuit.
//
//	moneyAdvanced ≤ 0                      → ErrNonPositiveMoney
//	moneyReturned ≤ moneyAdvanced           → ErrNoProfit
func NewMoneyDealingCircuit(moneyAdvanced, moneyReturned int64) (MoneyDealingCircuit, error) {
	if moneyAdvanced <= 0 {
		return MoneyDealingCircuit{}, ErrNonPositiveMoney
	}
	if moneyReturned <= moneyAdvanced {
		return MoneyDealingCircuit{}, ErrNoProfit
	}
	return MoneyDealingCircuit{
		MoneyAdvanced: moneyAdvanced,
		MoneyReturned: moneyReturned,
	}, nil
}

// ── MoneyOperation ───────────────────────────────────────────────────────────

// MoneyOperation is a value object describing a single technical money
// operation performed by a money-dealer. It is used in API responses and
// tests to annotate the five kinds with human-readable descriptions; it is
// NOT persisted separately — operation kinds are stored as a JSON array
// in MoneyDealingCapital.OperationKinds.
type MoneyOperation struct {
	Kind        MoneyOperationKind `json:"kind"`        // 1–5
	Description string             `json:"description"` // human-readable
}

// NewMoneyOperation builds and validates a MoneyOperation.
// kind must satisfy IsValid() — otherwise ErrInvalidOperationKind is returned.
func NewMoneyOperation(kind MoneyOperationKind, description string) (MoneyOperation, error) {
	if !kind.IsValid() {
		return MoneyOperation{}, ErrInvalidOperationKind
	}
	return MoneyOperation{Kind: kind, Description: description}, nil
}
