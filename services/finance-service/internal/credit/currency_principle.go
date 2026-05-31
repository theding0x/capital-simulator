// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 34 ("The Currency Principle and the Bank Acts of 1844 and 1845").
//
// The Currency School held that the total note issue should always move in step
// with the gold reserve — as gold flows out, notes contract; as gold flows in,
// notes expand. The Bank Act of 1844 institutionalised this doctrine by splitting
// the Bank of England into an Issue Department (notes 100 % backed beyond the
// fixed fiduciary limit) and a Banking Department (commercial lending). During
// crises of 1847 and 1857 the Act was suspended because mechanical contraction
// aggravated the crisis rather than curing it.
//
// Key concepts:
//   - BankActRegime distinguishes Pre-Act-1844, Act-1844-in-force, and Suspended.
//   - NoteIssueConstraint records the statutory limits: gold backing + fiduciary
//     limit = MaxNotes; NotesBeyondMax > 0 only when the Act is suspended.
//   - CurrencyPrincipleAnalysis captures the empirical refutation: gold drains
//     correlate with high interest, not with high commodity prices.
//   - GoldDrain records a gold reserve movement (export or import, amount).
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 34 sentinel errors.
var (
	// ErrInvalidBankActRegime is returned when the regime is not one of the
	// three defined constants.
	ErrInvalidBankActRegime = errors.New("credit: invalid bank act regime")

	// ErrInvalidBankDepartment is returned when the department is not one of
	// the two defined constants.
	ErrInvalidBankDepartment = errors.New("credit: invalid bank department")

	// ErrNegativeGoldDrain is returned when the gold drain amount is negative.
	ErrNegativeGoldDrain = errors.New("credit: gold drain amount must be >= 0")

	// ErrGoldDrainMustNotCorrelate is returned when gold_drain_correlates_with_prices
	// is set to true — the empirical record refutes the Currency Principle claim.
	ErrGoldDrainMustNotCorrelate = errors.New("credit: gold drain must not correlate with commodity prices (empirically refuted)")
)

// ── BankActRegime ────────────────────────────────────────────────────────────

// BankActRegime classifies the legal framework governing the Bank of England's
// note issue at a given moment.
type BankActRegime int

const (
	// RegimePreAct1844 is the period before the Bank Charter Act 1844.
	RegimePreAct1844 BankActRegime = 1
	// RegimeAct1844 is the normal operation of the Bank Charter Act 1844.
	RegimeAct1844 BankActRegime = 2
	// RegimeSuspended is the period when the Act has been suspended by Treasury
	// letter (crises of 1847 and 1857).
	RegimeSuspended BankActRegime = 3
)

// IsValid reports whether r is a defined BankActRegime constant.
func (r BankActRegime) IsValid() bool { return r >= 1 && r <= 3 }

// ── BankDepartment ───────────────────────────────────────────────────────────

// BankDepartment identifies which half of the Bank of England is being referred
// to after the 1844 Act split the institution.
type BankDepartment int

const (
	// DeptIssue is the Issue Department — responsible for note issue backed by gold.
	DeptIssue BankDepartment = 1
	// DeptBanking is the Banking Department — responsible for commercial lending.
	DeptBanking BankDepartment = 2
)

// IsValid reports whether d is a defined BankDepartment constant.
func (d BankDepartment) IsValid() bool { return d >= 1 && d <= 2 }

// ── NoteIssueConstraintID ────────────────────────────────────────────────────

// NoteIssueConstraintID identifies a persisted NoteIssueConstraint.
// 96-bit hex from crypto/rand.
type NoteIssueConstraintID string

// NewNoteIssueConstraintID returns a fresh random NoteIssueConstraintID.
func NewNoteIssueConstraintID() NoteIssueConstraintID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return NoteIssueConstraintID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty NoteIssueConstraintID.
func (id NoteIssueConstraintID) IsZero() bool { return id == "" }

// ── NoteIssueConstraint ──────────────────────────────────────────────────────

// NoteIssueConstraint is the persisted entity for a statutory note-issue
// regime (Vol. III Ch. 34). MaxNotes is always derived as GoldBacking +
// UnbackedLimit; NotesBeyondMax > 0 only when the Act is suspended.
//
// Invariants:
//
//	Regime.IsValid()
//	MaxNotes == GoldBacking + UnbackedLimit  (set by NewNoteIssueConstraint)
type NoteIssueConstraint struct {
	ID             NoteIssueConstraintID `json:"id"`
	Name           string                `json:"name"`
	Regime         BankActRegime         `json:"regime"`
	MaxNotes       int64                 `json:"max_notes"`
	GoldBacking    int64                 `json:"gold_backing"`
	UnbackedLimit  int64                 `json:"unbacked_limit"`
	NotesBeyondMax int64                 `json:"notes_beyond_max"`
	Description    string                `json:"description"`
	CreatedAt      time.Time             `json:"created_at"`
}

// NewNoteIssueConstraint constructs a NoteIssueConstraint, computing MaxNotes
// as GoldBacking + UnbackedLimit. ID and CreatedAt are left zero — the store
// assigns them.
//
//	!regime.IsValid() → ErrInvalidBankActRegime
//
// Fixture (Marx, Vol. III Ch. 34): Bank Act 1844 — fiduciary limit £14m,
// gold backing £14m → MaxNotes £28m. Crisis of 1857: Act suspended, notes
// beyond £28m issued (£48 883 in excess).
func NewNoteIssueConstraint(name string, regime BankActRegime, goldBacking, unbackedLimit, notesBeyondMax int64, description string) (NoteIssueConstraint, error) {
	if !regime.IsValid() {
		return NoteIssueConstraint{}, ErrInvalidBankActRegime
	}
	return NoteIssueConstraint{
		Name:           name,
		Regime:         regime,
		MaxNotes:       goldBacking + unbackedLimit,
		GoldBacking:    goldBacking,
		UnbackedLimit:  unbackedLimit,
		NotesBeyondMax: notesBeyondMax,
		Description:    description,
	}, nil
}

// ── CurrencyPrincipleAnalysisID ──────────────────────────────────────────────

// CurrencyPrincipleAnalysisID identifies a CurrencyPrincipleAnalysis.
type CurrencyPrincipleAnalysisID string

// NewCurrencyPrincipleAnalysisID returns a fresh random CurrencyPrincipleAnalysisID.
func NewCurrencyPrincipleAnalysisID() CurrencyPrincipleAnalysisID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return CurrencyPrincipleAnalysisID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CurrencyPrincipleAnalysisID.
func (id CurrencyPrincipleAnalysisID) IsZero() bool { return id == "" }

// ── CurrencyPrincipleAnalysis ────────────────────────────────────────────────

// CurrencyPrincipleAnalysis is a value type (not persisted independently) that
// captures the empirical test of the Currency Principle. Marx's evidence shows
// gold drains correlate with high interest, not with high commodity prices —
// refuting the Currency School's mechanical quantity-theory claim.
//
// Invariant: GoldDrainCorrelatesWithPrices is always false (forced by constructor).
type CurrencyPrincipleAnalysis struct {
	ID                            CurrencyPrincipleAnalysisID `json:"id"`
	GoldReserveChange             int64                       `json:"gold_reserve_change"`
	CommodityPriceChange          int64                       `json:"commodity_price_change"`
	InterestRateChange            int64                       `json:"interest_rate_change"`
	GoldDrainCorrelatesWithPrices bool                        `json:"gold_drain_correlates_with_prices"`
	Description                   string                      `json:"description"`
}

// NewCurrencyPrincipleAnalysis builds a CurrencyPrincipleAnalysis. The
// goldDrainCorrelatesWithPrices flag must be false — the empirical record
// refutes the Currency Principle's claim.
//
//	goldDrainCorrelatesWithPrices == true → ErrGoldDrainMustNotCorrelate
func NewCurrencyPrincipleAnalysis(goldReserveChange, commodityPriceChange, interestRateChange int64, goldDrainCorrelatesWithPrices bool, description string) (CurrencyPrincipleAnalysis, error) {
	if goldDrainCorrelatesWithPrices {
		return CurrencyPrincipleAnalysis{}, ErrGoldDrainMustNotCorrelate
	}
	return CurrencyPrincipleAnalysis{
		GoldReserveChange:             goldReserveChange,
		CommodityPriceChange:          commodityPriceChange,
		InterestRateChange:            interestRateChange,
		GoldDrainCorrelatesWithPrices: false,
		Description:                   description,
	}, nil
}

// ── GoldDrain ────────────────────────────────────────────────────────────────

// GoldDrain is a value type (not persisted independently) recording a movement
// of gold reserves — either an export (IsExport == true) reducing the reserve,
// or an import (IsExport == false) replenishing it.
//
// Invariant: Amount >= 0.
type GoldDrain struct {
	Amount      int64  `json:"amount"`
	IsExport    bool   `json:"is_export"`
	Description string `json:"description"`
}

// NewGoldDrain builds and validates a GoldDrain.
//
//	amount < 0 → ErrNegativeGoldDrain
func NewGoldDrain(amount int64, isExport bool, description string) (GoldDrain, error) {
	if amount < 0 {
		return GoldDrain{}, ErrNegativeGoldDrain
	}
	return GoldDrain{
		Amount:      amount,
		IsExport:    isExport,
		Description: description,
	}, nil
}
