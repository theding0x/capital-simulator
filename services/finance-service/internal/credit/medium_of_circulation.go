// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 28 ("Medium of Circulation and Capital — Views of
// Tooke and Fullarton").
//
// Tooke and Fullarton conflate three distinct categories:
//  1. Money functioning as medium of circulation (coin function) — revenue
//     expenditure by consumers and retailers; determined by the volume of
//     retail transactions and price levels.
//  2. Money functioning as capital (capital-transfer function) — transfers of
//     productive capital between capitalists; the movement of money-capital in
//     the circuit M—C…P…C′—M′.
//  3. Interest-bearing capital proper — the distinct M—M′ circuit.
//
// Marx corrects Tooke by showing that the quantity of notes outstanding as
// currency (coin function) is governed by the velocity and volume of retail
// trade, NOT by the volume of capital transfers.  The Bank of England 1844
// data illustrates the split: £20m notes outstanding, £8m held as reserve
// against convertibility demands.
//
// Key invariants:
//   - CurrencyObservation.TotalCurrencyOutstanding > 0.
//   - CoinFunctionAmount + CapitalTransferAmount <= TotalCurrencyOutstanding.
//   - ReserveFund.Amount >= 0.
//   - CurrencyFunction ∈ {1, 2} (FunctionCoin, FunctionCapitalTransfer).
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 28 sentinel errors.
var (
	// ErrNonPositiveCurrency is returned when total_currency_outstanding <= 0.
	ErrNonPositiveCurrency = errors.New("credit: total_currency_outstanding must be > 0")

	// ErrFunctionSplitExceedsTotal is returned when coin + capital-transfer
	// amounts exceed total_currency_outstanding.
	ErrFunctionSplitExceedsTotal = errors.New("credit: coin_function_amount + capital_transfer_amount must be <= total_currency_outstanding")

	// ErrNegativeReserve is returned when reserve_fund.amount < 0.
	ErrNegativeReserve = errors.New("credit: reserve_fund.amount must be >= 0")

	// ErrInvalidCurrencyFunction is returned when a CurrencyFunction value is
	// outside the valid set {1, 2}.
	ErrInvalidCurrencyFunction = errors.New("credit: currency_function must be 1 (coin) or 2 (capital-transfer)")
)

// ── CurrencyFunction ──────────────────────────────────────────────────────────

// CurrencyFunction distinguishes the two roles that money plays in circulation.
type CurrencyFunction int

const (
	// FunctionCoin (1) — money as medium of circulation for revenue expenditure;
	// consumers and retailers; determined by retail volume and price levels.
	FunctionCoin CurrencyFunction = 1

	// FunctionCapitalTransfer (2) — money as the vehicle for transferring
	// productive capital between capitalists; part of the circuit M—C…P…C′—M′.
	FunctionCapitalTransfer CurrencyFunction = 2
)

// IsValid reports whether f is a recognised CurrencyFunction.
func (f CurrencyFunction) IsValid() bool {
	return f == FunctionCoin || f == FunctionCapitalTransfer
}

// ── Value objects ─────────────────────────────────────────────────────────────

// ReserveFund is the cash reserve held against convertibility demands
// (e.g. the Banking Department's reserve at the Bank of England 1844).
type ReserveFund struct {
	Amount      int64  `json:"amount"`      // pence; >= 0
	Description string `json:"description"` // e.g. "Banking Department reserve"
}

// TookeFullartonAnalysis documents the three categories Tooke/Fullarton
// conflated, together with the correct Marxian distinction.
type TookeFullartonAnalysis struct {
	ConflatedCategories []string `json:"conflated_categories"` // what Tooke/Fullarton merged
	CorrectDistinction  string   `json:"correct_distinction"`  // Marx's correction
}

// RevenueExpenditureCircuit describes the coin-function circuit:
// consumers receive wages/revenue → spend in retail → coin returns to banks.
type RevenueExpenditureCircuit struct {
	Description  string `json:"description"`
	CoinAmountBP int64  `json:"coin_amount_bp"` // coin function as share of total (bp)
}

// CapitalTransferCircuit describes the capital-function circuit:
// capitalists transfer productive capital via bills of exchange or bank balances.
type CapitalTransferCircuit struct {
	Description     string `json:"description"`
	CapitalAmountBP int64  `json:"capital_amount_bp"` // capital-transfer as share (bp)
}

// ── CurrencyObservationID ─────────────────────────────────────────────────────

// CurrencyObservationID identifies a persisted CurrencyObservation.
type CurrencyObservationID string

// NewCurrencyObservationID returns a fresh random CurrencyObservationID.
func NewCurrencyObservationID() CurrencyObservationID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return CurrencyObservationID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CurrencyObservationID.
func (id CurrencyObservationID) IsZero() bool { return id == "" }

// ── CurrencyObservation ───────────────────────────────────────────────────────

// CurrencyObservation is the persisted entity for an observation of currency
// in circulation, split between its coin function (revenue expenditure) and its
// capital-transfer function (productive-capital movement), with a reserve fund
// held against convertibility demands.
//
// Invariants:
//
//	TotalCurrencyOutstanding > 0
//	CoinFunctionAmount + CapitalTransferAmount <= TotalCurrencyOutstanding
//	ReserveFund.Amount >= 0
type CurrencyObservation struct {
	ID                       CurrencyObservationID `json:"id"`
	Name                     string                `json:"name"`
	TotalCurrencyOutstanding int64                 `json:"total_currency_outstanding"` // pence; > 0
	CoinFunctionAmount       int64                 `json:"coin_function_amount"`       // pence; revenue expenditure
	CapitalTransferAmount    int64                 `json:"capital_transfer_amount"`    // pence; productive-capital transfers
	ReserveFund              ReserveFund           `json:"reserve_fund"`
	Description              string                `json:"description"`
	CreatedAt                time.Time             `json:"created_at"`
}

// NewCurrencyObservation constructs a CurrencyObservation, enforcing invariants.
// ID and CreatedAt are left zero — the store assigns them.
//
//	total <= 0                             → ErrNonPositiveCurrency
//	coin + capitalTransfer > total         → ErrFunctionSplitExceedsTotal
//	reserveFund.Amount < 0                 → ErrNegativeReserve
func NewCurrencyObservation(
	name string,
	total int64,
	coinAmount int64,
	capitalTransferAmount int64,
	reserveFund ReserveFund,
	description string,
) (CurrencyObservation, error) {
	if total <= 0 {
		return CurrencyObservation{}, ErrNonPositiveCurrency
	}
	if coinAmount+capitalTransferAmount > total {
		return CurrencyObservation{}, ErrFunctionSplitExceedsTotal
	}
	if reserveFund.Amount < 0 {
		return CurrencyObservation{}, ErrNegativeReserve
	}
	return CurrencyObservation{
		Name:                     name,
		TotalCurrencyOutstanding: total,
		CoinFunctionAmount:       coinAmount,
		CapitalTransferAmount:    capitalTransferAmount,
		ReserveFund:              reserveFund,
		Description:              description,
	}, nil
}

// CoinShareBP returns the coin-function share of total currency outstanding
// in basis points (rounded down).
//
//	CoinShareBP = (CoinFunctionAmount * 10000) / TotalCurrencyOutstanding
//
// Returns 0 when TotalCurrencyOutstanding == 0 to prevent division by zero.
func (obs CurrencyObservation) CoinShareBP() int64 {
	if obs.TotalCurrencyOutstanding == 0 {
		return 0
	}
	return (obs.CoinFunctionAmount * 10000) / obs.TotalCurrencyOutstanding
}

// ClassifyCurrency returns a brief description of the coin-vs-capital split
// for the given observation, mirroring Marx's distinction.
func ClassifyCurrency(obs CurrencyObservation) string {
	coinBP := obs.CoinShareBP()
	switch {
	case coinBP >= 7000:
		return "predominantly coin function: currency serves mainly revenue expenditure (retail/consumer)"
	case coinBP >= 2500:
		return "mixed function: currency serves both revenue expenditure and capital transfer"
	default:
		return "predominantly capital-transfer function: currency moves productive capital between capitalists"
	}
}

// NewTookeFullartonAnalysis returns the standard analysis of what Tooke and
// Fullarton conflated and Marx's correction.
func NewTookeFullartonAnalysis() TookeFullartonAnalysis {
	return TookeFullartonAnalysis{
		ConflatedCategories: []string{
			"money as medium of circulation (coin — retail/consumer revenue expenditure)",
			"money as capital (transfer of productive capital between capitalists)",
			"interest-bearing capital (the M—M′ circuit proper)",
		},
		CorrectDistinction: "Only the coin function governs the quantity of notes needed for circulation; " +
			"capital transfers between capitalists require no net addition to the note circulation — " +
			"they cancel via book entries or bills of exchange.",
	}
}
