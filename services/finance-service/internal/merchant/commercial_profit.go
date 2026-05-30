// Package merchant implements Capital Vol. III, Part IV.
// This file covers Ch. 17 ("Commercial Profit").
//
// Commercial profit is not new value created by the merchant but a portion of
// the surplus-value transferred from the industrial circuit. The merchant
// advances both productive capital (borrowed from the industrial sphere) and
// commercial capital (their own or borrowed trading fund). The general rate of
// profit is computed over the combined total capital — diluting the unadjusted
// rate that would obtain over the productive capital alone.
//
// Key invariants:
//   - AdjustedRateBP < UnadjustedRateBP whenever CommercialCapital > 0.
//   - NewValueCreated == 0 on every type in this file: commercial activity
//     does not produce value — it distributes already-produced surplus-value.
//   - Conservation: ProductiveCapitalistShare + MerchantShare == TotalSurplusValue.
package merchant

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// basisPoints is the fixed-point scale used throughout the merchant package:
// 1% = 100 bp, so a 20% rate = 2000 bp.
const basisPoints = int64(10000)

// Sentinel errors for Ch. 17 invariant violations.
var (
	// ErrNonPositiveCapital is returned when productive_capital or
	// commercial_capital is ≤ 0 — both must be strictly positive.
	ErrNonPositiveCapital = errors.New("merchant: capital must be positive")

	// ErrNegativeSurplusValue is returned when surplus_value is < 0.
	ErrNegativeSurplusValue = errors.New("merchant: surplus_value must be non-negative")

	// ErrNegativeLabourMinutes is returned when labour-time fields are
	// invalid (necessary ≤ 0, or workday < necessary, or minutes < 0).
	ErrNegativeLabourMinutes = errors.New("merchant: labour minutes out of range")

	// ErrCommercialCapitalExceedsTotal is reserved for future invariant
	// checks; currently unused but exported for completeness.
	ErrCommercialCapitalExceedsTotal = errors.New("merchant: commercial capital exceeds total")
)

// roundHalfUp computes (numer + denom/2) / denom — the canonical round-half-up
// idiom used throughout the simulator for basis-point quantities.
func roundHalfUp(numer, denom int64) int64 {
	return (numer + denom/2) / denom
}

// ── CommercialProfitID ───────────────────────────────────────────────────────

// CommercialProfitID identifies a persisted commercial-profit record.
// 96-bit hex from crypto/rand, matching every other domain ID in the simulator.
type CommercialProfitID string

// NewCommercialProfitID returns a fresh random CommercialProfitID.
func NewCommercialProfitID() CommercialProfitID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("merchant: rand: %w", err))
	}
	return CommercialProfitID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CommercialProfitID.
func (id CommercialProfitID) IsZero() bool { return id == "" }

// ── CommercialProfit ─────────────────────────────────────────────────────────

// CommercialProfit is the persisted entity for a commercial-profit computation
// (Vol. III Ch. 17). It records the productive and commercial capital advanced,
// the total surplus-value produced in the industrial circuit, and the two rates
// of profit — unadjusted (over productive capital alone) and adjusted (over the
// combined total). NewValueCreated is always 0: the merchant appropriates but
// does not produce surplus-value.
type CommercialProfit struct {
	ID                CommercialProfitID `json:"id"`
	ProductiveCapital int64              `json:"productive_capital"` // pence; > 0
	CommercialCapital int64              `json:"commercial_capital"` // pence; > 0
	SurplusValue      int64              `json:"surplus_value"`      // pence; >= 0
	UnadjustedRateBP  int64              `json:"unadjusted_rate_bp"` // s/C_productive in bp, round-half-up
	AdjustedRateBP    int64              `json:"adjusted_rate_bp"`   // s/(C_productive+C_commercial) in bp, round-half-up
	NewValueCreated   int64              `json:"new_value_created"`  // always 0
	CreatedAt         time.Time          `json:"created_at"`
}

// NewCommercialProfit constructs a CommercialProfit from the productive capital,
// commercial capital, and surplus value. Both capital arguments must be > 0;
// surplus_value must be >= 0. No ID or timestamp is assigned — the store does
// that on persistence.
//
//	UnadjustedRateBP = round_half_up(surplusValue * 10000, productiveCapital)
//	AdjustedRateBP   = round_half_up(surplusValue * 10000, productiveCapital + commercialCapital)
//
// Fixture verification:
//
//	productive=720000, commercial=180000, surplus=144000
//	unadjusted = (144000*10000 + 360000) / 720000 = 1440360000/720000 = 2000 bp (20%)
//	adjusted   = (144000*10000 + 450000) / 900000 = 1440450000/900000 = 1600 bp (16%)
func NewCommercialProfit(productiveCapital, commercialCapital, surplusValue int64) (CommercialProfit, error) {
	if productiveCapital <= 0 || commercialCapital <= 0 {
		return CommercialProfit{}, ErrNonPositiveCapital
	}
	if surplusValue < 0 {
		return CommercialProfit{}, ErrNegativeSurplusValue
	}
	unadj := roundHalfUp(surplusValue*basisPoints, productiveCapital)
	adj := roundHalfUp(surplusValue*basisPoints, productiveCapital+commercialCapital)
	return CommercialProfit{
		ProductiveCapital: productiveCapital,
		CommercialCapital: commercialCapital,
		SurplusValue:      surplusValue,
		UnadjustedRateBP:  unadj,
		AdjustedRateBP:    adj,
		NewValueCreated:   0,
	}, nil
}

// ── AdjustedGeneralRate ──────────────────────────────────────────────────────

// AdjustedGeneralRate is a value object representing the general rate of profit
// once commercial capital is included in the denominator. It is computed on
// demand, not persisted.
type AdjustedGeneralRate struct {
	ProductiveCapital int64 `json:"productive_capital"`
	CommercialCapital int64 `json:"commercial_capital"`
	SurplusValue      int64 `json:"surplus_value"`
	UnadjustedRateBP  int64 `json:"unadjusted_rate_bp"`
	AdjustedRateBP    int64 `json:"adjusted_rate_bp"`
}

// Validate returns ErrNonPositiveCapital if either capital is ≤ 0, or
// ErrNegativeSurplusValue if surplus_value < 0.
func (a AdjustedGeneralRate) Validate() error {
	if a.ProductiveCapital <= 0 || a.CommercialCapital <= 0 {
		return ErrNonPositiveCapital
	}
	if a.SurplusValue < 0 {
		return ErrNegativeSurplusValue
	}
	return nil
}

// NewAdjustedGeneralRate builds and validates an AdjustedGeneralRate.
func NewAdjustedGeneralRate(productiveCapital, commercialCapital, surplusValue int64) (AdjustedGeneralRate, error) {
	if productiveCapital <= 0 || commercialCapital <= 0 {
		return AdjustedGeneralRate{}, ErrNonPositiveCapital
	}
	if surplusValue < 0 {
		return AdjustedGeneralRate{}, ErrNegativeSurplusValue
	}
	unadj := roundHalfUp(surplusValue*basisPoints, productiveCapital)
	adj := roundHalfUp(surplusValue*basisPoints, productiveCapital+commercialCapital)
	return AdjustedGeneralRate{
		ProductiveCapital: productiveCapital,
		CommercialCapital: commercialCapital,
		SurplusValue:      surplusValue,
		UnadjustedRateBP:  unadj,
		AdjustedRateBP:    adj,
	}, nil
}

// ── SurplusValueDeduction ────────────────────────────────────────────────────

// SurplusValueDeduction records the division of the total surplus-value between
// the productive capitalist and the merchant, in proportion to their respective
// capital advances (LabourMinutes as the unit for the value magnitudes).
//
// Conservation invariant: ProductiveCapitalistShare + MerchantShare == TotalSurplusValue.
type SurplusValueDeduction struct {
	TotalSurplusValue          int64 `json:"total_surplus_value"`          // LabourMinutes
	ProductiveCapitalistShare  int64 `json:"productive_capitalist_share"`  // LabourMinutes
	MerchantShare              int64 `json:"merchant_share"`               // LabourMinutes
}

// NewSurplusValueDeduction builds a SurplusValueDeduction, splitting
// totalSurplusValue in proportion to productiveCapital : commercialCapital.
//
//	MerchantShare              = round_half_up(total * commercial, productive + commercial)
//	ProductiveCapitalistShare  = total - MerchantShare
func NewSurplusValueDeduction(totalSurplusValue, productiveCapital, commercialCapital int64) (SurplusValueDeduction, error) {
	if productiveCapital <= 0 || commercialCapital <= 0 {
		return SurplusValueDeduction{}, ErrNonPositiveCapital
	}
	if totalSurplusValue < 0 {
		return SurplusValueDeduction{}, ErrNegativeSurplusValue
	}
	total := productiveCapital + commercialCapital
	merchantShare := roundHalfUp(totalSurplusValue*commercialCapital, total)
	return SurplusValueDeduction{
		TotalSurplusValue:         totalSurplusValue,
		ProductiveCapitalistShare: totalSurplusValue - merchantShare,
		MerchantShare:             merchantShare,
	}, nil
}

// ── CommercialWorker ─────────────────────────────────────────────────────────

// CommercialWorker is the value object for a commercial wage-worker employed by
// the merchant (Vol. III Ch. 17). The worker performs unpaid labour for the
// merchant (enabling a larger profit share) but creates no new value — all value
// is produced in the industrial circuit.
//
//	NewValueCreated  = 0 always.
//	UnpaidLabourMinutes = WorkdayMinutes - NecessaryLabourMinutes.
type CommercialWorker struct {
	NecessaryLabourMinutes int64 `json:"necessary_labour_minutes"` // LabourMinutes; > 0
	WorkdayMinutes         int64 `json:"workday_minutes"`          // LabourMinutes; >= Necessary
	UnpaidLabourMinutes    int64 `json:"unpaid_labour_minutes"`    // = Workday - Necessary
	NewValueCreated        int64 `json:"new_value_created"`        // always 0
}

// NewCommercialWorker builds and validates a CommercialWorker.
// necessaryLabourMinutes must be > 0; workdayMinutes must be >= necessaryLabourMinutes.
func NewCommercialWorker(necessaryLabourMinutes, workdayMinutes int64) (CommercialWorker, error) {
	if necessaryLabourMinutes <= 0 {
		return CommercialWorker{}, ErrNegativeLabourMinutes
	}
	if workdayMinutes < necessaryLabourMinutes {
		return CommercialWorker{}, ErrNegativeLabourMinutes
	}
	return CommercialWorker{
		NecessaryLabourMinutes: necessaryLabourMinutes,
		WorkdayMinutes:         workdayMinutes,
		UnpaidLabourMinutes:    workdayMinutes - necessaryLabourMinutes,
		NewValueCreated:        0,
	}, nil
}

// ── UnpaidCommercialLabour ───────────────────────────────────────────────────

// UnpaidCommercialLabour is the value object for the unpaid portion of
// commercial labour time. Like all commercial activity, it creates no new value
// — it is a form of unpaid labour that enables the merchant to appropriate a
// greater share of the industrial surplus-value.
type UnpaidCommercialLabour struct {
	Minutes         int64 `json:"minutes"`          // LabourMinutes; >= 0
	NewValueCreated int64 `json:"new_value_created"` // always 0
}

// NewUnpaidCommercialLabour builds and validates an UnpaidCommercialLabour.
// minutes must be >= 0.
func NewUnpaidCommercialLabour(minutes int64) (UnpaidCommercialLabour, error) {
	if minutes < 0 {
		return UnpaidCommercialLabour{}, ErrNegativeLabourMinutes
	}
	return UnpaidCommercialLabour{
		Minutes:         minutes,
		NewValueCreated: 0,
	}, nil
}

// ── WholesaleRetailSpread ────────────────────────────────────────────────────

// WholesaleRetailSpread is the value object for the price spread between the
// wholesale price paid by the merchant and the retail price charged to the
// consumer. The spread is a form of commercial profit — a deduction from the
// surplus-value — not new value created in trade.
type WholesaleRetailSpread struct {
	WholesalePrice int64 `json:"wholesale_price"` // pence; > 0
	RetailPrice    int64 `json:"retail_price"`    // pence; > WholesalePrice
	SpreadPence    int64 `json:"spread_pence"`    // = RetailPrice - WholesalePrice
}

// NewWholesaleRetailSpread builds and validates a WholesaleRetailSpread.
// wholesalePrice must be > 0; retailPrice must be > wholesalePrice.
func NewWholesaleRetailSpread(wholesalePrice, retailPrice int64) (WholesaleRetailSpread, error) {
	if wholesalePrice <= 0 {
		return WholesaleRetailSpread{}, ErrNonPositiveCapital
	}
	if retailPrice <= wholesalePrice {
		return WholesaleRetailSpread{}, ErrNonPositiveCapital
	}
	return WholesaleRetailSpread{
		WholesalePrice: wholesalePrice,
		RetailPrice:    retailPrice,
		SpreadPence:    retailPrice - wholesalePrice,
	}, nil
}
