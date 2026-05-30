// Package merchant implements Capital Vol. III, Part IV.
// This file covers Ch. 18 ("The Turnover of Merchant's Capital; Prices").
//
// The key insight is that the annual profit on merchant's capital depends on
// the total money advanced and the general rate of profit — NOT on the number
// of turnovers or units sold. Turnover affects only the per-unit markup: the
// faster the merchant turns their capital, the smaller the markup on each
// article must be to yield the same annual profit.
//
// Key invariants:
//   - AnnualMerchantProfit = roundHalfUp(moneyAdvanced * generalRateBP, basisPoints).
//   - MarkupPerUnit = AnnualMerchantProfit / (TurnoverCount * UnitsPerTurnover)
//     (integer division — truncation, not round-half-up).
//   - AnnualMerchantProfit is invariant to TurnoverCount and UnitsPerTurnover.
//   - MerchantPriceOfProduction = CostPrice + AverageProfit (both non-negative).
package merchant

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Sentinel errors for Ch. 18 invariant violations.
var (
	// ErrInvalidTurnoverCount is returned when turnover_count < 1.
	ErrInvalidTurnoverCount = errors.New("merchant: turnover_count must be >= 1")

	// ErrNoUnits is returned when units_per_turnover <= 0.
	ErrNoUnits = errors.New("merchant: units_per_turnover must be positive")
)

// ── MerchantTurnoverID ───────────────────────────────────────────────────────

// MerchantTurnoverID identifies a persisted merchant-turnover record.
// 96-bit hex from crypto/rand, matching every other domain ID in the simulator.
type MerchantTurnoverID string

// NewMerchantTurnoverID returns a fresh random MerchantTurnoverID.
func NewMerchantTurnoverID() MerchantTurnoverID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("merchant: rand: %w", err))
	}
	return MerchantTurnoverID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty MerchantTurnoverID.
func (id MerchantTurnoverID) IsZero() bool { return id == "" }

// ── MerchantTurnover ─────────────────────────────────────────────────────────

// MerchantTurnover is the persisted entity for a merchant-turnover computation
// (Vol. III Ch. 18). It records the money advanced, the general rate of profit,
// and the turnover count and units, deriving the annual merchant profit and the
// per-unit markup.
//
// AnnualMerchantProfit is invariant to turnover: increasing TurnoverCount or
// UnitsPerTurnover lowers MarkupPerUnit but leaves AnnualMerchantProfit unchanged.
type MerchantTurnover struct {
	ID                  MerchantTurnoverID `json:"id"`
	MoneyAdvanced       int64              `json:"money_advanced"`        // pence; > 0
	GeneralRateBP       int64              `json:"general_rate_bp"`       // bp; > 0
	TurnoverCount       int64              `json:"turnover_count"`        // count; >= 1
	UnitsPerTurnover    int64              `json:"units_per_turnover"`    // count; > 0
	AnnualMerchantProfit int64             `json:"annual_merchant_profit"` // pence; derived
	MarkupPerUnit       int64              `json:"markup_per_unit"`       // pence; derived (integer div)
	CreatedAt           time.Time          `json:"created_at"`
}

// NewMerchantTurnover constructs a MerchantTurnover from the money advanced,
// general rate of profit, turnover count, and units per turnover.
//
//	AnnualMerchantProfit = roundHalfUp(moneyAdvanced * generalRateBP, basisPoints)
//	MarkupPerUnit        = AnnualMerchantProfit / (turnoverCount * unitsPerTurnover)
//	                       (integer division — truncation)
//
// Fixture verification (seed row 1801):
//
//	moneyAdvanced=100000, generalRateBP=1500, turnoverCount=1, unitsPerTurnover=300
//	annual = roundHalfUp(100000*1500, 10000) = roundHalfUp(150000000, 10000) = 15000
//	markup = 15000 / (1 * 300) = 50
//
// Fixture verification (seed row 1803):
//
//	turnoverCount=3, unitsPerTurnover=300
//	markup = 15000 / (3 * 300) = 15000 / 900 = 16 (integer div; remainder 600 discarded)
func NewMerchantTurnover(moneyAdvanced, generalRateBP, turnoverCount, unitsPerTurnover int64) (MerchantTurnover, error) {
	if moneyAdvanced <= 0 || generalRateBP <= 0 {
		return MerchantTurnover{}, ErrNonPositiveCapital
	}
	if turnoverCount < 1 {
		return MerchantTurnover{}, ErrInvalidTurnoverCount
	}
	if unitsPerTurnover <= 0 {
		return MerchantTurnover{}, ErrNoUnits
	}
	annual := roundHalfUp(moneyAdvanced*generalRateBP, basisPoints)
	markup := annual / (turnoverCount * unitsPerTurnover)
	return MerchantTurnover{
		MoneyAdvanced:       moneyAdvanced,
		GeneralRateBP:       generalRateBP,
		TurnoverCount:       turnoverCount,
		UnitsPerTurnover:    unitsPerTurnover,
		AnnualMerchantProfit: annual,
		MarkupPerUnit:       markup,
	}, nil
}

// ── TurnoverPoint ────────────────────────────────────────────────────────────

// TurnoverPoint is a single data point in a TurnoverEffectAnalysis — the
// per-unit markup at a given turnover count and unit volume. As TurnoverCount
// rises, MarkupPerUnit falls (AnnualMerchantProfit remains constant).
type TurnoverPoint struct {
	TurnoverCount int64 `json:"turnover_count"`
	TotalUnits    int64 `json:"total_units"`    // = TurnoverCount * UnitsPerTurnover
	MarkupPerUnit int64 `json:"markup_per_unit"` // pence; AnnualMerchantProfit / TotalUnits (integer div)
}

// ── TurnoverEffectAnalysis ───────────────────────────────────────────────────

// TurnoverEffectAnalysis is a value object (NOT persisted) that shows how
// MarkupPerUnit shrinks as TurnoverCount rises, while AnnualMerchantProfit
// remains constant. It is computed on demand and returned directly from the
// POST /v1/merchant/turnover-effect endpoint.
type TurnoverEffectAnalysis struct {
	MoneyAdvanced       int64           `json:"money_advanced"`
	GeneralRateBP       int64           `json:"general_rate_bp"`
	AnnualMerchantProfit int64          `json:"annual_merchant_profit"`
	Points              []TurnoverPoint `json:"points"`
}

// NewTurnoverEffectAnalysis builds a TurnoverEffectAnalysis from the given
// money advanced, general rate, units per turnover, and a slice of turnover
// counts to evaluate. At least one count must be supplied; each must be >= 1.
//
//	AnnualMerchantProfit = roundHalfUp(moneyAdvanced * generalRateBP, basisPoints)
//	for each count: TotalUnits = count * unitsPerTurnover
//	                MarkupPerUnit = AnnualMerchantProfit / TotalUnits (integer div)
func NewTurnoverEffectAnalysis(moneyAdvanced, generalRateBP, unitsPerTurnover int64, turnoverCounts []int64) (TurnoverEffectAnalysis, error) {
	if moneyAdvanced <= 0 || generalRateBP <= 0 {
		return TurnoverEffectAnalysis{}, ErrNonPositiveCapital
	}
	if unitsPerTurnover <= 0 {
		return TurnoverEffectAnalysis{}, ErrNoUnits
	}
	if len(turnoverCounts) < 1 {
		return TurnoverEffectAnalysis{}, ErrInvalidTurnoverCount
	}
	for _, c := range turnoverCounts {
		if c < 1 {
			return TurnoverEffectAnalysis{}, ErrInvalidTurnoverCount
		}
	}
	annual := roundHalfUp(moneyAdvanced*generalRateBP, basisPoints)
	points := make([]TurnoverPoint, len(turnoverCounts))
	for i, count := range turnoverCounts {
		total := count * unitsPerTurnover
		points[i] = TurnoverPoint{
			TurnoverCount: count,
			TotalUnits:    total,
			MarkupPerUnit: annual / total,
		}
	}
	return TurnoverEffectAnalysis{
		MoneyAdvanced:       moneyAdvanced,
		GeneralRateBP:       generalRateBP,
		AnnualMerchantProfit: annual,
		Points:              points,
	}, nil
}

// ── MerchantPriceOfProduction ────────────────────────────────────────────────

// MerchantPriceOfProduction is the value object for the price of production
// as seen from the merchant's standpoint: the cost-price paid to the producer
// plus the average profit that the merchant's capital must earn at the general
// rate (Vol. III Ch. 18). This is not new value — it is a claim on the
// surplus-value already produced in the industrial circuit.
//
//	PriceOfProduction = CostPrice + AverageProfit
type MerchantPriceOfProduction struct {
	CostPrice           int64 `json:"cost_price"`            // pence; >= 0
	AverageProfit       int64 `json:"average_profit"`        // pence; >= 0
	PriceOfProduction   int64 `json:"price_of_production"`   // pence; = CostPrice + AverageProfit
}

// NewMerchantPriceOfProduction builds and validates a MerchantPriceOfProduction.
// costPrice and averageProfit must both be >= 0.
func NewMerchantPriceOfProduction(costPrice, averageProfit int64) (MerchantPriceOfProduction, error) {
	if costPrice < 0 || averageProfit < 0 {
		return MerchantPriceOfProduction{}, ErrNonPositiveCapital
	}
	return MerchantPriceOfProduction{
		CostPrice:         costPrice,
		AverageProfit:     averageProfit,
		PriceOfProduction: costPrice + averageProfit,
	}, nil
}
