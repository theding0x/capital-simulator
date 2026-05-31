// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 23 ("Interest and Profit of Enterprise").
//
// The total profit of capital (= average profit) is divided into two
// qualitatively distinct categories: interest (paid to the money-capitalist
// as owner) and profit of enterprise (retained by the functioning capitalist
// as active user). This division appears as if interest is the natural fruit
// of capital-ownership and profit of enterprise is the reward of management —
// an ideological inversion that conceals the surplus-value relation.
//
// Key invariants:
//   - ProfitDivision.TotalProfitBP == ProfitDivision.InterestBP + ProfitDivision.ProfitOfEnterpriseBP.
//   - ProfitDivision.InterestBP >= 0 and ProfitDivision.ProfitOfEnterpriseBP >= 0.
//   - ProfitDivision.InterestBP < ProfitDivision.TotalProfitBP (interest is a portion, not the whole).
//   - CapitalOwnershipForm ∈ {FormOwnerOperator=1, FormSeparated=2}.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 23 sentinel errors.
var (
	// ErrNonPositiveTotalProfit is returned when total_profit_bp <= 0.
	ErrNonPositiveTotalProfit = errors.New("credit: total_profit_bp must be positive")

	// ErrNegativeInterest is returned when interest_bp < 0.
	ErrNegativeInterest = errors.New("credit: interest_bp must be >= 0")

	// ErrInterestNotLessThanTotal is returned when interest_bp >= total_profit_bp.
	// Marx's invariant: interest is a part of profit, strictly less than the whole.
	ErrInterestNotLessThanTotal = errors.New("credit: interest_bp must be less than total_profit_bp")

	// ErrInvalidOwnershipForm is returned when CapitalOwnershipForm is not 1 or 2.
	ErrInvalidOwnershipForm = errors.New("credit: capital_ownership_form must be 1 (OwnerOperator) or 2 (Separated)")
)

// ── CapitalOwnershipForm ──────────────────────────────────────────────────────

// CapitalOwnershipForm distinguishes whether the money-capitalist (lender) and
// the functioning capitalist (borrower) are the same person or separate persons.
type CapitalOwnershipForm int

const (
	// FormOwnerOperator is the form where a capitalist employs their own capital
	// in production; no interest is paid to an external lender.
	FormOwnerOperator CapitalOwnershipForm = 1
	// FormSeparated is the form where the functioning capitalist borrows from a
	// distinct money-capitalist; interest flows from borrower to lender.
	FormSeparated CapitalOwnershipForm = 2
)

// IsValid reports whether f is one of the two defined constants.
func (f CapitalOwnershipForm) IsValid() bool {
	return f == FormOwnerOperator || f == FormSeparated
}

// ── MystificationIndex ────────────────────────────────────────────────────────

// MystificationIndex is the degree to which interest conceals the surplus-value
// origin of profit. Expressed in basis points [0, 10000]: 10000 bp means the
// interest relation is entirely opaque; 0 bp means the surplus-value relation is
// fully transparent.
type MystificationIndex int64

// ── WagesOfSuperintendence ────────────────────────────────────────────────────

// WagesOfSuperintendence is the labour component of the functioning capitalist's
// income in large joint-stock companies. The salaried manager earns wages of
// skilled labour, not profit of enterprise. Marx's invariant: where ownership and
// management are fully separated, the "profit of enterprise" vanishes as a
// category distinct from wages of supervision.
type WagesOfSuperintendence struct {
	LabourMinutes int64 `json:"labour_minutes"` // canonical value-magnitude unit; >= 0
}

// ── ProfitOfEnterprise ────────────────────────────────────────────────────────

// ProfitOfEnterprise is the value object for the residual portion of the average
// profit retained by the functioning capitalist after paying interest to the
// money-capitalist.
//
//	ProfitOfEnterpriseBP = TotalProfitBP - InterestBP
type ProfitOfEnterprise struct {
	ProfitOfEnterpriseBP int64 `json:"profit_of_enterprise_bp"` // bp; = total - interest; >= 0
}

// NewProfitOfEnterprise constructs a ProfitOfEnterprise by subtraction.
//
//	totalProfitBP <= 0         → ErrNonPositiveTotalProfit
//	interestBP < 0             → ErrNegativeInterest
//	interestBP >= totalProfitBP → ErrInterestNotLessThanTotal
func NewProfitOfEnterprise(totalProfitBP, interestBP int64) (ProfitOfEnterprise, error) {
	if totalProfitBP <= 0 {
		return ProfitOfEnterprise{}, ErrNonPositiveTotalProfit
	}
	if interestBP < 0 {
		return ProfitOfEnterprise{}, ErrNegativeInterest
	}
	if interestBP >= totalProfitBP {
		return ProfitOfEnterprise{}, ErrInterestNotLessThanTotal
	}
	return ProfitOfEnterprise{ProfitOfEnterpriseBP: totalProfitBP - interestBP}, nil
}

// ── ProfitDivisionID ──────────────────────────────────────────────────────────

// ProfitDivisionID identifies a persisted profit-division record.
// 96-bit hex from crypto/rand, matching every other domain ID in the simulator.
type ProfitDivisionID string

// NewProfitDivisionID returns a fresh random ProfitDivisionID.
func NewProfitDivisionID() ProfitDivisionID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return ProfitDivisionID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty ProfitDivisionID.
func (id ProfitDivisionID) IsZero() bool { return id == "" }

// ── ProfitOfEnterpriseID ──────────────────────────────────────────────────────

// ProfitOfEnterpriseID identifies a profit-of-enterprise value object.
// 96-bit hex from crypto/rand.
type ProfitOfEnterpriseID string

// NewProfitOfEnterpriseID returns a fresh random ProfitOfEnterpriseID.
func NewProfitOfEnterpriseID() ProfitOfEnterpriseID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return ProfitOfEnterpriseID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty ProfitOfEnterpriseID.
func (id ProfitOfEnterpriseID) IsZero() bool { return id == "" }

// ── ProfitDivision ────────────────────────────────────────────────────────────

// ProfitDivision is the persisted entity for a profit-division record
// (Vol. III Ch. 23). It splits the total average profit into interest paid to
// the money-capitalist and profit of enterprise retained by the functioning
// capitalist. ProfitOfEnterpriseBP is always derived by subtraction.
//
// Invariants:
//
//	TotalProfitBP > 0
//	InterestBP >= 0
//	ProfitOfEnterpriseBP >= 0
//	TotalProfitBP == InterestBP + ProfitOfEnterpriseBP
//	InterestBP < TotalProfitBP
//	OwnershipForm.IsValid()
type ProfitDivision struct {
	ID                       ProfitDivisionID     `json:"id"`
	TotalProfitBP            int64                `json:"total_profit_bp"`              // bp; > 0
	InterestBP               int64                `json:"interest_bp"`                  // bp; >= 0
	ProfitOfEnterpriseBP     int64                `json:"profit_of_enterprise_bp"`       // bp; = total - interest; >= 0
	OwnershipForm            CapitalOwnershipForm `json:"ownership_form"`               // 1 or 2
	WagesOfSuperintendence   WagesOfSuperintendence `json:"wages_of_superintendence"`   // labour component
	MystificationIndex       MystificationIndex   `json:"mystification_index"`          // bp [0,10000]
	CreatedAt                time.Time            `json:"created_at"`
}

// NewProfitDivision constructs a ProfitDivision, deriving ProfitOfEnterpriseBP
// by subtraction. No ID or timestamp is assigned — the store does that on persistence.
//
//	totalProfitBP <= 0          → ErrNonPositiveTotalProfit
//	interestBP < 0              → ErrNegativeInterest
//	interestBP >= totalProfitBP → ErrInterestNotLessThanTotal
//	!ownershipForm.IsValid()    → ErrInvalidOwnershipForm
//
// Fixtures (Marx, Vol. III Ch. 23):
//
//	Borrowed capital: totalProfitBP=2000, interestBP=500, ownershipForm=FormSeparated
//	  → ProfitOfEnterpriseBP=1500
//	Own capital:      totalProfitBP=2000, interestBP=0,   ownershipForm=FormOwnerOperator
//	  → ERROR: interestBP=0 is valid only when totalProfitBP>0, but interestBP < totalProfitBP is required;
//	  own-capital fixture: use interestBP=0 with FormOwnerOperator is valid since 0 < 2000.
func NewProfitDivision(
	totalProfitBP, interestBP int64,
	ownershipForm CapitalOwnershipForm,
	wages WagesOfSuperintendence,
	mystification MystificationIndex,
) (ProfitDivision, error) {
	if totalProfitBP <= 0 {
		return ProfitDivision{}, ErrNonPositiveTotalProfit
	}
	if interestBP < 0 {
		return ProfitDivision{}, ErrNegativeInterest
	}
	if interestBP >= totalProfitBP {
		return ProfitDivision{}, ErrInterestNotLessThanTotal
	}
	if !ownershipForm.IsValid() {
		return ProfitDivision{}, ErrInvalidOwnershipForm
	}
	return ProfitDivision{
		TotalProfitBP:          totalProfitBP,
		InterestBP:             interestBP,
		ProfitOfEnterpriseBP:   totalProfitBP - interestBP,
		OwnershipForm:          ownershipForm,
		WagesOfSuperintendence: wages,
		MystificationIndex:     mystification,
	}, nil
}
