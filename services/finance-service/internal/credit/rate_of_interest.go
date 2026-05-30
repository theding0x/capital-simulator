// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 22 ("Division of Profit. Rate of Interest. Natural
// Rate of Interest").
//
// The rate of interest has no "natural" or law-governed level — it is
// determined empirically by supply and demand in the money market. Its
// maximum limit is the total average profit (above which no borrower could
// profit); its minimum is indeterminate and can fall nearly to zero.
//
// Key invariants:
//   - RateOfInterest.RateBP >= 0.
//   - RateOfInterest.RateBP < AverageProfitRateBP (interest is a portion of
//     profit, strictly less than the whole).
//   - InterestRateAnalysis.ProfitOfEnterpriseBP == AverageProfitRateBP - InterestRateBP.
//   - IndustrialCyclePhase ∈ {1,2,3,4,5}.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 22 sentinel errors.
var (
	// ErrNegativeRate is returned when rate_bp < 0.
	ErrNegativeRate = errors.New("credit: rate_bp must be >= 0")

	// ErrNonPositiveAverageProfit is returned when average_profit_rate_bp <= 0.
	ErrNonPositiveAverageProfit = errors.New("credit: average_profit_rate_bp must be positive")

	// ErrRateExceedsAverageProfit is returned when rate_bp >= average_profit_rate_bp.
	ErrRateExceedsAverageProfit = errors.New("credit: rate_bp must be less than average_profit_rate_bp")

	// ErrInvalidCyclePhase is returned when a phase value is outside 1..5.
	ErrInvalidCyclePhase = errors.New("credit: industrial_cycle_phase must be 1..5")
)

// ── IndustrialCyclePhase ──────────────────────────────────────────────────────

// IndustrialCyclePhase is the phase of the industrial cycle at the time an
// interest-rate observation is recorded. Marx identifies five phases; the
// interest rate behaves differently in each.
type IndustrialCyclePhase int

const (
	PhaseInactivity     IndustrialCyclePhase = 1 // stagnation; low demand for money-capital
	PhaseRevival        IndustrialCyclePhase = 2 // recovery; rising demand
	PhaseProsperity     IndustrialCyclePhase = 3 // boom; moderate interest, high profit
	PhaseOverproduction IndustrialCyclePhase = 4 // overproduction; credit stretched
	PhaseCrisis         IndustrialCyclePhase = 5 // crisis; sharp spike in interest rate
)

// IsValid reports whether p is one of the five defined constants.
func (p IndustrialCyclePhase) IsValid() bool {
	return p >= PhaseInactivity && p <= PhaseCrisis
}

// ── RateOfInterestID ─────────────────────────────────────────────────────────

// RateOfInterestID identifies a persisted rate-of-interest observation.
// 96-bit hex from crypto/rand, matching every other domain ID in the simulator.
type RateOfInterestID string

// NewRateOfInterestID returns a fresh random RateOfInterestID.
func NewRateOfInterestID() RateOfInterestID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return RateOfInterestID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty RateOfInterestID.
func (id RateOfInterestID) IsZero() bool { return id == "" }

// ── RateOfInterest ────────────────────────────────────────────────────────────

// RateOfInterest is the persisted entity for an interest-rate observation
// (Vol. III Ch. 22). It records the observed market rate of interest in basis
// points, the corresponding average profit rate, and the phase of the
// industrial cycle.
//
// Invariants:
//
//	RateBP >= 0
//	AverageProfitRateBP > 0
//	RateBP < AverageProfitRateBP
//	CyclePhase.IsValid()
type RateOfInterest struct {
	ID                  RateOfInterestID     `json:"id"`
	RateBP              int64                `json:"rate_bp"`                // bp; >= 0
	AverageProfitRateBP int64                `json:"average_profit_rate_bp"` // bp; > 0
	CyclePhase          IndustrialCyclePhase `json:"cycle_phase"`            // 1..5
	Period              string               `json:"period"`                 // e.g. "1843", "1847 crisis"
	CreatedAt           time.Time            `json:"created_at"`
}

// NewRateOfInterest constructs a RateOfInterest from observed values.
// No ID or timestamp is assigned — the store does that on persistence.
//
//	rateBP < 0                         → ErrNegativeRate
//	averageProfitRateBP <= 0           → ErrNonPositiveAverageProfit
//	rateBP >= averageProfitRateBP      → ErrRateExceedsAverageProfit
//	!phase.IsValid()                   → ErrInvalidCyclePhase
//
// Fixture (Marx, Vol. III Ch. 22):
//
//	Prosperity 1843: avg 2000 bp, interest 150 bp → phase PhaseProsperity
//	Crisis 1847:     avg 2000 bp, interest 800 bp → phase PhaseCrisis
func NewRateOfInterest(rateBP, averageProfitRateBP int64, phase IndustrialCyclePhase, period string) (RateOfInterest, error) {
	if rateBP < 0 {
		return RateOfInterest{}, ErrNegativeRate
	}
	if averageProfitRateBP <= 0 {
		return RateOfInterest{}, ErrNonPositiveAverageProfit
	}
	if rateBP >= averageProfitRateBP {
		return RateOfInterest{}, ErrRateExceedsAverageProfit
	}
	if !phase.IsValid() {
		return RateOfInterest{}, ErrInvalidCyclePhase
	}
	return RateOfInterest{
		RateBP:              rateBP,
		AverageProfitRateBP: averageProfitRateBP,
		CyclePhase:          phase,
		Period:              period,
	}, nil
}

// ── InterestRateAnalysis ─────────────────────────────────────────────────────

// InterestRateAnalysis is the compute-only (not persisted) result of splitting
// the average profit rate into interest and profit of enterprise.
//
//	ProfitOfEnterpriseBP = AverageProfitRateBP - InterestRateBP
//
// Invariants:
//
//	InterestRateBP >= 0
//	InterestRateBP < AverageProfitRateBP
//	ProfitOfEnterpriseBP == AverageProfitRateBP - InterestRateBP
type InterestRateAnalysis struct {
	AverageProfitRateBP  int64 `json:"average_profit_rate_bp"`  // bp; > 0
	InterestRateBP       int64 `json:"interest_rate_bp"`        // bp; >= 0
	ProfitOfEnterpriseBP int64 `json:"profit_of_enterprise_bp"` // bp; = avg - interest
}

// NewInterestRateAnalysis computes the split of average profit into interest
// and profit of enterprise. Not persisted — stateless POST endpoint only.
//
//	interestRateBP < 0                       → ErrNegativeRate
//	averageProfitRateBP <= 0                 → ErrNonPositiveAverageProfit
//	interestRateBP >= averageProfitRateBP    → ErrRateExceedsAverageProfit
//
// Fixture (Marx, Vol. III Ch. 22):
//
//	avg 2000 bp, interest 500 bp (¼ of 2000):
//	ProfitOfEnterpriseBP = 2000 - 500 = 1500
func NewInterestRateAnalysis(averageProfitRateBP, interestRateBP int64) (InterestRateAnalysis, error) {
	if interestRateBP < 0 {
		return InterestRateAnalysis{}, ErrNegativeRate
	}
	if averageProfitRateBP <= 0 {
		return InterestRateAnalysis{}, ErrNonPositiveAverageProfit
	}
	if interestRateBP >= averageProfitRateBP {
		return InterestRateAnalysis{}, ErrRateExceedsAverageProfit
	}
	return InterestRateAnalysis{
		AverageProfitRateBP:  averageProfitRateBP,
		InterestRateBP:       interestRateBP,
		ProfitOfEnterpriseBP: averageProfitRateBP - interestRateBP,
	}, nil
}

// ── InterestRateBound ─────────────────────────────────────────────────────────

// InterestRateBound is the value object for the theoretical bounds on the rate
// of interest. The maximum is the average profit rate (above which no borrower
// could profit). The minimum approaches but never reaches zero.
type InterestRateBound struct {
	MaxBP int64 `json:"max_bp"` // bp; == average profit rate
	MinBP int64 `json:"min_bp"` // bp; >= 0, <= MaxBP
}

// NewInterestRateBound builds and validates an InterestRateBound.
//
//	maxBP <= 0     → ErrNonPositiveAverageProfit
//	minBP < 0      → ErrNegativeRate
//	minBP > maxBP  → ErrRateExceedsAverageProfit
func NewInterestRateBound(maxBP, minBP int64) (InterestRateBound, error) {
	if maxBP <= 0 {
		return InterestRateBound{}, ErrNonPositiveAverageProfit
	}
	if minBP < 0 {
		return InterestRateBound{}, ErrNegativeRate
	}
	if minBP > maxBP {
		return InterestRateBound{}, ErrRateExceedsAverageProfit
	}
	return InterestRateBound{MaxBP: maxBP, MinBP: minBP}, nil
}

// ── MarketInterestRate ────────────────────────────────────────────────────────

// MarketInterestRate is the value object for the daily-quoted rate in the
// money market. Unlike the rate of profit (a tendency), the market interest
// rate is directly observable and uniform at any given moment.
type MarketInterestRate struct {
	RateBP int64 `json:"rate_bp"` // bp; >= 0
}

// NewMarketInterestRate builds and validates a MarketInterestRate.
// rateBP < 0 → ErrNegativeRate.
func NewMarketInterestRate(rateBP int64) (MarketInterestRate, error) {
	if rateBP < 0 {
		return MarketInterestRate{}, ErrNegativeRate
	}
	return MarketInterestRate{RateBP: rateBP}, nil
}
