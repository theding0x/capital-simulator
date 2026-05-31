// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 25 ("Credit and Fictitious Capital").
//
// The credit system develops from money's function as a means of payment.
// Bills of exchange circulate as commercial money before their maturity date.
// Fictitious capital arises when an income stream is capitalised at the
// prevailing rate of interest: CapitalisedValue = AnnualIncome × 10000 / RateBP.
// Government bonds represent consumed (non-existent) capital — the original
// loan is spent; only the claim to future tax revenue persists. Stocks
// represent titles to future surplus-value, not the underlying real capital.
//
// Key invariants:
//   - FictitiousCapital.CapitalisedValue == roundHalfUp(AnnualIncome*10000, RateBP).
//   - KindPublicDebt has RealCapitalExists == false (consumed capital).
//   - BillOfExchange.FaceValue > 0 and MaturityDays > 0.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 25 sentinel errors.
var (
	ErrNonPositiveFaceValue   = errors.New("credit: face_value must be positive")
	ErrNonPositiveMaturityDays = errors.New("credit: maturity_days must be positive")
	ErrNonPositiveAnnualIncome = errors.New("credit: annual_income must be positive")
	ErrNonPositiveRateBPFict   = errors.New("credit: rate_bp must be positive")
)

// FictitiousKind classifies the species of fictitious capital.
type FictitiousKind int

const (
	// KindBill is a bill of exchange treated as fictitious capital.
	KindBill FictitiousKind = 1
	// KindBankNote is a bank note circulating beyond real backing.
	KindBankNote FictitiousKind = 2
	// KindPublicDebt is a government bond — the original capital has been
	// consumed; only a claim on future tax revenue remains.
	KindPublicDebt FictitiousKind = 3
	// KindShare is a stock representing a title to future surplus-value.
	KindShare FictitiousKind = 4
)

// ── BillOfExchangeID ──────────────────────────────────────────────────────────

// BillOfExchangeID identifies a persisted BillOfExchange.
type BillOfExchangeID string

// NewBillOfExchangeID returns a fresh random BillOfExchangeID.
func NewBillOfExchangeID() BillOfExchangeID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return BillOfExchangeID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty BillOfExchangeID.
func (id BillOfExchangeID) IsZero() bool { return id == "" }

// ── BillOfExchange ────────────────────────────────────────────────────────────

// BillOfExchange is a promissory note with a fixed maturity date, circulating
// as commercial money before the due date (Vol. III Ch. 25).
//
// Invariants:
//
//	FaceValue > 0
//	MaturityDays > 0
type BillOfExchange struct {
	ID          BillOfExchangeID `json:"id"`
	Drawer      string           `json:"drawer"`        // name of issuer
	Drawee      string           `json:"drawee"`        // name of debtor
	FaceValue   int64            `json:"face_value"`    // > 0; pence
	MaturityDays int64           `json:"maturity_days"` // > 0
	Description string           `json:"description"`
	CreatedAt   time.Time        `json:"created_at"`
}

// NewBillOfExchange constructs a BillOfExchange, enforcing invariants.
//
//	faceValue <= 0   → ErrNonPositiveFaceValue
//	maturityDays <= 0 → ErrNonPositiveMaturityDays
func NewBillOfExchange(drawer, drawee string, faceValue, maturityDays int64, description string) (BillOfExchange, error) {
	if faceValue <= 0 {
		return BillOfExchange{}, ErrNonPositiveFaceValue
	}
	if maturityDays <= 0 {
		return BillOfExchange{}, ErrNonPositiveMaturityDays
	}
	return BillOfExchange{
		Drawer:       drawer,
		Drawee:       drawee,
		FaceValue:    faceValue,
		MaturityDays: maturityDays,
		Description:  description,
	}, nil
}

// ── FictitiousCapitalID ───────────────────────────────────────────────────────

// FictitiousCapitalID identifies a persisted FictitiousCapital.
type FictitiousCapitalID string

// NewFictitiousCapitalID returns a fresh random FictitiousCapitalID.
func NewFictitiousCapitalID() FictitiousCapitalID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return FictitiousCapitalID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty FictitiousCapitalID.
func (id FictitiousCapitalID) IsZero() bool { return id == "" }

// ── FictitiousCapital ─────────────────────────────────────────────────────────

// FictitiousCapital is a capitalised income stream (Vol. III Ch. 25).
// CapitalisedValue is derived: roundHalfUp(AnnualIncome*10000, RateBP).
// For KindPublicDebt, RealCapitalExists == false — the original loan is spent.
//
// Invariants:
//
//	AnnualIncome > 0
//	RateBP > 0
//	CapitalisedValue == roundHalfUp(AnnualIncome*10000, RateBP)
//	KindPublicDebt → RealCapitalExists == false
type FictitiousCapital struct {
	ID                 FictitiousCapitalID `json:"id"`
	Kind               FictitiousKind      `json:"kind"`                 // 1=bill,2=bank-note,3=public-debt,4=share
	Name               string              `json:"name"`
	AnnualIncome       int64               `json:"annual_income"`        // > 0; pence
	RateBP             int64               `json:"rate_bp"`              // > 0; prevailing interest rate in bp
	CapitalisedValue   int64               `json:"capitalised_value"`    // DERIVED: roundHalfUp(income*10000, rate_bp)
	RealCapitalExists  bool                `json:"real_capital_exists"`  // false for public-debt (consumed)
	Description        string              `json:"description"`
	CreatedAt          time.Time           `json:"created_at"`
}

// NewFictitiousCapital constructs a FictitiousCapital, deriving CapitalisedValue.
//
//	annualIncome <= 0 → ErrNonPositiveAnnualIncome
//	rateBP <= 0       → ErrNonPositiveRateBPFict
//
// Fixture (Marx, Vol. III Ch. 25 — government consol):
//
//	AnnualIncome=500 p (£5), RateBP=500 (5%) → CapitalisedValue=10000 p (£100).
//	The original £100 is consumed; real_capital_exists=false.
func NewFictitiousCapital(kind FictitiousKind, name string, annualIncome, rateBP int64, description string) (FictitiousCapital, error) {
	if annualIncome <= 0 {
		return FictitiousCapital{}, ErrNonPositiveAnnualIncome
	}
	if rateBP <= 0 {
		return FictitiousCapital{}, ErrNonPositiveRateBPFict
	}
	capitalised := roundHalfUp(annualIncome*basisPoints, rateBP)
	realCapital := kind != KindPublicDebt
	return FictitiousCapital{
		Kind:              kind,
		Name:              name,
		AnnualIncome:      annualIncome,
		RateBP:            rateBP,
		CapitalisedValue:  capitalised,
		RealCapitalExists: realCapital,
		Description:       description,
	}, nil
}
