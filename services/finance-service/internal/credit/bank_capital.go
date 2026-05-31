// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 29 ("Component Parts of Bank Capital").
//
// Bank capital splits into two material parts:
//
//	(1) cash — gold coin and bank-notes actually held in the till;
//	(2) securities — commercial paper (bills of exchange) and public
//	    securities (government bonds, stocks, mortgages).
//
// The greater part of bank capital is FICTITIOUS. Deposits are mere book
// entries; bonds and stocks are titles to a future stream of surplus-value,
// not the underlying real capital. Their market value is set by capitalising
// the income they yield at the prevailing rate of interest:
//
//	MarketValue = AnnualIncome × 10000 / InterestRateBP.
//
// Because the interest rate sits in the denominator, the market value of every
// security moves INVERSELY with the rate of interest: a fall in the rate
// inflates paper values, a rise deflates them, while the real capital (where it
// exists at all) is untouched.
//
// Key invariants:
//   - BankCapital.TotalCapital == CashAmount + SecuritiesAmount.
//   - FictitiousCapitalValuation.MarketValue == roundHalfUp(AnnualIncome*10000, InterestRateBP).
//   - InterestRateBP > 0.
//   - DepositEntry.Amount >= 0 (a book entry can be zero but never negative).
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 29 sentinel errors.
var (
	// ErrBankCapitalMismatch is returned when TotalCapital != CashAmount + SecuritiesAmount.
	ErrBankCapitalMismatch = errors.New("credit: total_capital must equal cash_amount + securities_amount")

	// ErrInvalidComponentKind is returned when a BankCapitalComponent carries a kind outside 1..5.
	ErrInvalidComponentKind = errors.New("credit: bank_capital_component kind must be 1..5")

	// ErrNonPositiveInterestRate is returned when interest_rate_bp <= 0.
	ErrNonPositiveInterestRate = errors.New("credit: interest_rate_bp must be positive")

	// ErrNegativeDeposit is returned when a DepositEntry amount is < 0.
	ErrNegativeDeposit = errors.New("credit: deposit amount must be >= 0")
)

// BankCapitalComponentKind classifies a single component of bank capital
// (Vol. III Ch. 29).
type BankCapitalComponentKind int

const (
	// KindCash is gold coin or bank-notes held in the till.
	KindCash BankCapitalComponentKind = 1
	// KindBillOfExchange is commercial paper — a discounted bill.
	KindBillOfExchange BankCapitalComponentKind = 2
	// KindGovernmentBond is a public security — a title to future tax revenue.
	KindGovernmentBond BankCapitalComponentKind = 3
	// KindStock is a share — a title to future surplus-value.
	KindStock BankCapitalComponentKind = 4
	// KindMortgage is a mortgage — a title secured on land or buildings.
	KindMortgage BankCapitalComponentKind = 5
)

// IsValid reports whether k is a defined BankCapitalComponentKind (1..5).
func (k BankCapitalComponentKind) IsValid() bool {
	return k >= KindCash && k <= KindMortgage
}

// BankCapitalComponent is one constituent of a bank's capital — an amount of a
// particular kind (cash, bill, bond, stock, mortgage). Value type only.
type BankCapitalComponent struct {
	Amount      int64                    `json:"amount"` // pence
	Kind        BankCapitalComponentKind `json:"kind"`   // 1=cash,2=bill,3=bond,4=stock,5=mortgage
	Description string                   `json:"description"`
}

// ── BankCapitalID ─────────────────────────────────────────────────────────────

// BankCapitalID identifies a persisted BankCapital.
type BankCapitalID string

// NewBankCapitalID returns a fresh random BankCapitalID.
func NewBankCapitalID() BankCapitalID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return BankCapitalID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty BankCapitalID.
func (id BankCapitalID) IsZero() bool { return id == "" }

// ── BankCapital ───────────────────────────────────────────────────────────────

// BankCapital is the decomposition of a bank's capital into cash and securities
// (Vol. III Ch. 29). The greater part is fictitious — deposits are book entries
// and securities are capitalised income titles, not real capital.
//
// Invariant: TotalCapital == CashAmount + SecuritiesAmount (computed by the
// constructor, so it always holds).
type BankCapital struct {
	ID               BankCapitalID          `json:"id"`
	Name             string                 `json:"name"`
	CashAmount       int64                  `json:"cash_amount"`       // pence; gold + notes
	SecuritiesAmount int64                  `json:"securities_amount"` // pence; bills + public securities
	TotalCapital     int64                  `json:"total_capital"`     // DERIVED: cash + securities
	Components       []BankCapitalComponent `json:"components"`
	Description      string                 `json:"description"`
	CreatedAt        time.Time              `json:"created_at"`
}

// NewBankCapital constructs a BankCapital, computing TotalCapital so the
// invariant TotalCapital == CashAmount + SecuritiesAmount always holds. Every
// component's kind is validated.
//
//	any component kind outside 1..5 → ErrInvalidComponentKind
//
// Fixture (Marx, Vol. III Ch. 29 — a Scottish bank):
//
//	CashAmount=300000000 (£3m), SecuritiesAmount=2400000000 (£24m)
//	→ TotalCapital=2700000000 (£27m); the cash is a small fraction of the whole.
func NewBankCapital(
	name string,
	cashAmount, securitiesAmount int64,
	components []BankCapitalComponent,
	description string,
) (BankCapital, error) {
	for _, c := range components {
		if !c.Kind.IsValid() {
			return BankCapital{}, ErrInvalidComponentKind
		}
	}
	return BankCapital{
		Name:             name,
		CashAmount:       cashAmount,
		SecuritiesAmount: securitiesAmount,
		TotalCapital:     cashAmount + securitiesAmount,
		Components:       components,
		Description:      description,
	}, nil
}

// ── FictitiousCapitalValuationID ──────────────────────────────────────────────

// FictitiousCapitalValuationID identifies a persisted FictitiousCapitalValuation.
type FictitiousCapitalValuationID string

// NewFictitiousCapitalValuationID returns a fresh random FictitiousCapitalValuationID.
func NewFictitiousCapitalValuationID() FictitiousCapitalValuationID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return FictitiousCapitalValuationID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty FictitiousCapitalValuationID.
func (id FictitiousCapitalValuationID) IsZero() bool { return id == "" }

// ── FictitiousCapitalValuation ────────────────────────────────────────────────

// FictitiousCapitalValuation capitalises a security's income stream at the
// prevailing rate of interest (Vol. III Ch. 29). MarketValue is derived:
// roundHalfUp(AnnualIncome*10000, InterestRateBP). Because InterestRateBP is the
// denominator, MarketValue moves inversely with the rate of interest.
//
// For a stock, AnnualIncome is typically NominalValue × DividendRateBP / 10000;
// DividendRateBP is recorded for reference but does not alter MarketValue, which
// is set purely by the income and the prevailing rate.
//
// Invariants:
//
//	InterestRateBP > 0
//	MarketValue == roundHalfUp(AnnualIncome*10000, InterestRateBP)
type FictitiousCapitalValuation struct {
	ID             FictitiousCapitalValuationID `json:"id"`
	Name           string                       `json:"name"`
	NominalValue   int64                        `json:"nominal_value"`    // pence; face / par value
	AnnualIncome   int64                        `json:"annual_income"`    // pence; yearly yield
	InterestRateBP int64                        `json:"interest_rate_bp"` // > 0; prevailing rate in bp
	DividendRateBP int64                        `json:"dividend_rate_bp"` // optional; for stocks, in bp
	MarketValue    int64                        `json:"market_value"`     // DERIVED: roundHalfUp(income*10000, rate_bp)
	Description    string                       `json:"description"`
	CreatedAt      time.Time                    `json:"created_at"`
}

// NewFictitiousCapitalValuation constructs a FictitiousCapitalValuation,
// deriving MarketValue so the invariant always holds.
//
//	interestRateBP <= 0 → ErrNonPositiveInterestRate
//
// Fixture (Marx, Vol. III Ch. 29 — government bond):
//
//	AnnualIncome=500 p (£5), InterestRateBP=500 (5%) → MarketValue=10000 p (£100).
//	The same income at 250 bp → 20000; at 1000 bp → 5000 (inverse relation).
func NewFictitiousCapitalValuation(
	name string,
	nominalValue, annualIncome, interestRateBP, dividendRateBP int64,
	description string,
) (FictitiousCapitalValuation, error) {
	if interestRateBP <= 0 {
		return FictitiousCapitalValuation{}, ErrNonPositiveInterestRate
	}
	market := roundHalfUp(annualIncome*basisPoints, interestRateBP)
	return FictitiousCapitalValuation{
		Name:           name,
		NominalValue:   nominalValue,
		AnnualIncome:   annualIncome,
		InterestRateBP: interestRateBP,
		DividendRateBP: dividendRateBP,
		MarketValue:    market,
		Description:    description,
	}, nil
}

// ── DepositEntry ──────────────────────────────────────────────────────────────

// DepositEntryID identifies a DepositEntry value.
type DepositEntryID string

// NewDepositEntryID returns a fresh random DepositEntryID.
func NewDepositEntryID() DepositEntryID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return DepositEntryID(hex.EncodeToString(b[:]))
}

// DepositEntry is a book-entry deposit — the purest form of fictitious bank
// capital (Vol. III Ch. 29). It is a value type only; it is never persisted, as
// a deposit is an accounting entry rather than a thing the bank owns.
//
// Invariant: Amount >= 0.
type DepositEntry struct {
	ID          DepositEntryID `json:"id"`
	Amount      int64          `json:"amount"` // pence; >= 0
	Description string         `json:"description"`
}

// NewDepositEntry constructs a DepositEntry, enforcing Amount >= 0.
//
//	amount < 0 → ErrNegativeDeposit
func NewDepositEntry(amount int64, description string) (DepositEntry, error) {
	if amount < 0 {
		return DepositEntry{}, ErrNegativeDeposit
	}
	return DepositEntry{
		ID:          NewDepositEntryID(),
		Amount:      amount,
		Description: description,
	}, nil
}
