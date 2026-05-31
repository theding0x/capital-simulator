// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 32 ("Money-Capital and Real Capital. III (Conclusion)").
//
// The conclusion of the trilogy: accumulated loanable money-capital ALWAYS
// OVERSTATES real accumulation. Several streams swell the loan market without
// any new productive capital behind them:
//
//   - Revenue intended for consumption passes temporarily through a money form
//     and appears as loan capital, though it is no new productive capital.
//   - Capital released by falling input prices is freed and lent out; it
//     represents real accumulation ONLY when it is actually reinvested.
//   - Merchants' idle money and the savings of the growing rentier class swell
//     loan capital independently of any real expansion.
//
// In crisis the demand for loan capital is a demand for MEANS OF PAYMENT —
// convertibility into gold — NOT a demand for productive capital. Banking
// theory that confuses the two collapses on exactly this distinction.
//
// Key invariants:
//   - CapitalRelease.ReleasedAmount >= 0 (pence).
//   - CapitalRelease.Cause ∈ {1,2,3,4} (CapitalReleaseCause.IsValid()).
//   - LoanCapitalOverstatement.OverstatementPct >= 0 (basis points): loan
//     capital never understates real accumulation in the long run.
//   - RevenueLoanCapitalConversion.RevenueAmount > 0.
//   - MeansOfPaymentDemand in PhaseCrisis is for convertibility, not productive
//     capital (IsProductiveCapitalDemand == false).
//
// IndustrialCyclePhase, PhaseCrisis and ErrInvalidCyclePhase are reused from
// rate_of_interest.go (Ch. 22) — they are not redeclared here.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 32 sentinel errors.
var (
	// ErrNegativeReleasedAmount is returned when capital_release.released_amount < 0.
	ErrNegativeReleasedAmount = errors.New("credit: released_amount must be >= 0")

	// ErrInvalidReleaseCause is returned when capital_release.cause is not a valid
	// CapitalReleaseCause.
	ErrInvalidReleaseCause = errors.New("credit: invalid capital-release cause")

	// ErrNonPositiveRevenue is returned when revenue_loan_capital_conversion.revenue_amount <= 0.
	ErrNonPositiveRevenue = errors.New("credit: revenue_amount must be > 0")

	// ErrNegativeOverstatement is returned when loan_capital_overstatement.overstatement_pct < 0.
	ErrNegativeOverstatement = errors.New("credit: overstatement_pct must be >= 0")
)

// ── CapitalReleaseCause ─────────────────────────────────────────────────────────

// CapitalReleaseCause names why a parcel of money-capital is freed onto the loan
// market without any matching expansion of real (productive) capital.
type CapitalReleaseCause int

const (
	// CauseInputPriceFall is capital released because the price of raw materials
	// or other inputs has fallen — the same physical reproduction now ties up less
	// money, and the difference is freed.
	CauseInputPriceFall CapitalReleaseCause = 1
	// CauseMerchantIdleMoney is merchants' idle money that sits between purchases
	// and sales and is lent out in the interval.
	CauseMerchantIdleMoney CapitalReleaseCause = 2
	// CauseRentierSavings is the savings of the growing rentier class, withdrawn
	// from productive investment and placed at interest.
	CauseRentierSavings CapitalReleaseCause = 3
	// CauseRevenuePassthrough is revenue intended for consumption passing
	// temporarily through a money form, appearing as loan capital on the way.
	CauseRevenuePassthrough CapitalReleaseCause = 4
)

// IsValid reports whether c is one of the four recognised causes.
func (c CapitalReleaseCause) IsValid() bool {
	return c >= CauseInputPriceFall && c <= CauseRevenuePassthrough
}

// ── LoanCapitalOverstatement ────────────────────────────────────────────────────

// LoanCapitalOverstatement is the gap between apparent (loanable money-capital)
// accumulation and real (productive) accumulation, expressed in basis points. It
// is always >= 0: in the long run loan capital never understates real
// accumulation. A value type embedded in CapitalRelease, never persisted
// independently.
type LoanCapitalOverstatement struct {
	OverstatementPct int64  `json:"overstatement_pct"` // basis points; >= 0
	Description      string `json:"description"`
}

// ── RevenueLoanCapitalConversionID ──────────────────────────────────────────────

// RevenueLoanCapitalConversionID identifies a RevenueLoanCapitalConversion. The
// conversion is a value type — it is not persisted — but it carries an ID so it
// can be referenced like the other domain types.
type RevenueLoanCapitalConversionID string

// NewRevenueLoanCapitalConversionID returns a fresh random RevenueLoanCapitalConversionID.
func NewRevenueLoanCapitalConversionID() RevenueLoanCapitalConversionID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return RevenueLoanCapitalConversionID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty RevenueLoanCapitalConversionID.
func (id RevenueLoanCapitalConversionID) IsZero() bool { return id == "" }

// ── RevenueLoanCapitalConversion ────────────────────────────────────────────────

// RevenueLoanCapitalConversion captures revenue intended for consumption that
// passes temporarily through a money form and appears as loan capital, though it
// represents no new productive capital. A value type, never persisted.
//
// Invariant:
//
//	RevenueAmount > 0
type RevenueLoanCapitalConversion struct {
	ID            RevenueLoanCapitalConversionID `json:"id"`
	RevenueAmount int64                          `json:"revenue_amount"` // pence; > 0
	Description   string                         `json:"description"`
}

// NewRevenueLoanCapitalConversion builds and validates a
// RevenueLoanCapitalConversion. The ID is left zero; callers may assign one with
// NewRevenueLoanCapitalConversionID.
//
//	revenueAmount <= 0 → ErrNonPositiveRevenue
//
// Fixture (Marx, Vol. III Ch. 32): a spinner realises the proceeds of his yarn
// as money and sets a portion aside as revenue; deposited in a bank it becomes
// loan capital for the period, though no new productive capital exists.
func NewRevenueLoanCapitalConversion(revenueAmount int64, description string) (RevenueLoanCapitalConversion, error) {
	if revenueAmount <= 0 {
		return RevenueLoanCapitalConversion{}, ErrNonPositiveRevenue
	}
	return RevenueLoanCapitalConversion{
		RevenueAmount: revenueAmount,
		Description:   description,
	}, nil
}

// ── MeansOfPaymentDemandID ──────────────────────────────────────────────────────

// MeansOfPaymentDemandID identifies a MeansOfPaymentDemand. The demand is a value
// type — it is not persisted — but it carries an ID so it can be referenced like
// the other domain types.
type MeansOfPaymentDemandID string

// NewMeansOfPaymentDemandID returns a fresh random MeansOfPaymentDemandID.
func NewMeansOfPaymentDemandID() MeansOfPaymentDemandID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return MeansOfPaymentDemandID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty MeansOfPaymentDemandID.
func (id MeansOfPaymentDemandID) IsZero() bool { return id == "" }

// ── MeansOfPaymentDemand ────────────────────────────────────────────────────────

// MeansOfPaymentDemand records, for a phase of the industrial cycle, whether the
// demand for loan capital is a demand for productive capital or merely a demand
// for means of payment — convertibility into gold. In a crisis it is the latter:
// IsProductiveCapitalDemand is false. A value type, never persisted.
type MeansOfPaymentDemand struct {
	ID                        MeansOfPaymentDemandID `json:"id"`
	Phase                     IndustrialCyclePhase   `json:"phase"`  // 1..5
	Amount                    int64                  `json:"amount"` // pence
	IsProductiveCapitalDemand bool                   `json:"is_productive_capital_demand"`
	Description               string                 `json:"description"`
}

// NewMeansOfPaymentDemand builds and validates a MeansOfPaymentDemand. The ID is
// left zero; callers may assign one with NewMeansOfPaymentDemandID. When the
// phase is PhaseCrisis the demand is forced to means of payment only
// (IsProductiveCapitalDemand is set to false), echoing Marx's point that in a
// crisis the cry is for convertibility, not for fresh productive capital.
//
//	!phase.IsValid() → ErrInvalidCyclePhase
//
// Fixture (Marx, Vol. III Ch. 32): the 1857 crisis — demand for loan capital is
// demand for convertibility into gold, not for expansion of productive capacity.
func NewMeansOfPaymentDemand(phase IndustrialCyclePhase, amount int64, isProductiveCapitalDemand bool, description string) (MeansOfPaymentDemand, error) {
	if !phase.IsValid() {
		return MeansOfPaymentDemand{}, ErrInvalidCyclePhase
	}
	if phase == PhaseCrisis {
		isProductiveCapitalDemand = false
	}
	return MeansOfPaymentDemand{
		Phase:                     phase,
		Amount:                    amount,
		IsProductiveCapitalDemand: isProductiveCapitalDemand,
		Description:               description,
	}, nil
}

// ── Rentier ─────────────────────────────────────────────────────────────────────

// Rentier is a capitalist living off interest — a member of the growing class
// that withdraws its savings from productive investment and places them at
// interest, swelling loan capital without expanding real capital. A value type,
// never persisted.
type Rentier struct {
	Name           string `json:"name"`
	InterestIncome int64  `json:"interest_income"` // pence
	Description    string `json:"description"`
}

// ── CapitalReleaseID ────────────────────────────────────────────────────────────

// CapitalReleaseID identifies a persisted CapitalRelease.
type CapitalReleaseID string

// NewCapitalReleaseID returns a fresh random CapitalReleaseID.
func NewCapitalReleaseID() CapitalReleaseID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return CapitalReleaseID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CapitalReleaseID.
func (id CapitalReleaseID) IsZero() bool { return id == "" }

// ── CapitalRelease ──────────────────────────────────────────────────────────────

// CapitalRelease is the persisted entity for a parcel of money-capital freed onto
// the loan market (Vol. III Ch. 32). It records the amount, the cause, whether it
// was reinvested (a release represents real accumulation only when reinvested),
// and the embedded overstatement — the gap between apparent and real
// accumulation. The only persisted entity in this chapter.
//
// Invariants:
//
//	ReleasedAmount >= 0
//	Cause.IsValid()
//	Overstatement.OverstatementPct >= 0
type CapitalRelease struct {
	ID             CapitalReleaseID         `json:"id"`
	Name           string                   `json:"name"`
	ReleasedAmount int64                    `json:"released_amount"` // pence; >= 0
	Cause          CapitalReleaseCause      `json:"cause"`           // 1..4
	Reinvested     bool                     `json:"reinvested"`      // real accumulation only when true
	Overstatement  LoanCapitalOverstatement `json:"loan_capital_overstatement"`
	Description    string                   `json:"description"`
	CreatedAt      time.Time                `json:"created_at"`
}

// NewCapitalRelease constructs a CapitalRelease, enforcing invariants. ID and
// CreatedAt are left zero — the store assigns them.
//
//	releasedAmount < 0                   → ErrNegativeReleasedAmount
//	!cause.IsValid()                     → ErrInvalidReleaseCause
//	overstatement.OverstatementPct < 0   → ErrNegativeOverstatement
//
// Fixture (Marx, Vol. III Ch. 32): capital released by a fall in raw-material
// prices (Cause=CauseInputPriceFall). Lent out but not reinvested it merely
// overstates real accumulation; once reinvested it reflects real accumulation.
func NewCapitalRelease(
	name string,
	releasedAmount int64,
	cause CapitalReleaseCause,
	reinvested bool,
	overstatement LoanCapitalOverstatement,
	description string,
) (CapitalRelease, error) {
	if releasedAmount < 0 {
		return CapitalRelease{}, ErrNegativeReleasedAmount
	}
	if !cause.IsValid() {
		return CapitalRelease{}, ErrInvalidReleaseCause
	}
	if overstatement.OverstatementPct < 0 {
		return CapitalRelease{}, ErrNegativeOverstatement
	}
	return CapitalRelease{
		Name:           name,
		ReleasedAmount: releasedAmount,
		Cause:          cause,
		Reinvested:     reinvested,
		Overstatement:  overstatement,
		Description:    description,
	}, nil
}

// AccumulationVerdict returns whether this release reflects real accumulation or
// merely overstates it. A release is real accumulation only when the freed
// capital is actually reinvested in production; otherwise it swells loan capital
// while no new productive capital exists behind it.
func (c CapitalRelease) AccumulationVerdict() string {
	if c.Reinvested {
		return "reflects real accumulation — the released capital is reinvested in production"
	}
	return "overstates real accumulation — the released capital swells loan capital but is not reinvested"
}
