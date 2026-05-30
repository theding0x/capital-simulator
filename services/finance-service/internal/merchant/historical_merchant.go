// Package merchant implements Capital Vol. III, Part IV.
// This file covers Ch. 20 ("Historical Facts about Merchant's Capital").
//
// Merchant's capital is older than the capitalist mode of production — it is in
// fact historically the oldest free state of existence of capital. Marx traces
// three stages of development: (1) pre-capitalist merchant capital that profits
// by unequal exchange and employs no wage-labour; (2) a transition stage as
// industrial production grows; (3) the subordination of merchant capital under
// industrial capital as the general rate of profit equalises.
//
// Key invariants:
//   - Invariant 1: pre-capitalist merchant capital (StagePreCapitalist) cannot
//     employ wage-labour. Setting WageLabour=true with stage=1 returns
//     ErrWageLabourInPreCapitalist.
//   - Invariant 2 (inverse-development law): as industrial capitalism advances
//     (rising stage), the SubordinationIndex of merchant capital rises — i.e.
//     merchant capital is drawn more tightly into the orbit of industrial capital.
//     InverseDevelopmentRelation.Monotonic captures this.
//   - SubordinationIndex is in basis points [0, 10000]: merchant capital as a
//     share of total social capital × 10000.
//
// basisPoints and roundHalfUp are declared in commercial_profit.go; they are
// not redeclared here.
package merchant

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"time"
)

// Sentinel errors for Ch. 20 invariant violations.
var (
	// ErrInvalidDevelopmentStage is returned when a DevelopmentStage value is
	// not in the valid range 1–3.
	ErrInvalidDevelopmentStage = errors.New("merchant: invalid development_stage — must be 1–3")

	// ErrSubordinationOutOfRange is returned when SubordinationIndex is outside
	// [0, 10000] basis points.
	ErrSubordinationOutOfRange = errors.New("merchant: subordination_index must be within [0,10000]")

	// ErrWageLabourInPreCapitalist is returned when a pre-capitalist merchant
	// capital record attempts to employ wage-labour (invariant 1).
	ErrWageLabourInPreCapitalist = errors.New("merchant: pre-capitalist merchant capital cannot employ wage-labour")
)

// DevelopmentStage enumerates the three historical stages through which
// merchant capital passes in its relation to industrial capital (Vol. III Ch. 20).
// Modelled exactly on MoneyOperationKind in money_dealing.go.
type DevelopmentStage int

const (
	// StagePreCapitalist — merchant capital operates before or outside the
	// capitalist mode of production. Profit arises from unequal exchange, not
	// wage-labour exploitation.
	StagePreCapitalist DevelopmentStage = 1
	// StageTransition — industrial production is growing; merchant capital begins
	// to be drawn into the circuit of capital but has not yet been fully
	// subordinated.
	StageTransition DevelopmentStage = 2
	// StageSubordinated — merchant capital is fully subordinated under industrial
	// capital; it participates in the general rate of profit and may employ
	// wage-labour (commercial workers).
	StageSubordinated DevelopmentStage = 3
)

// IsValid reports whether s is a recognised DevelopmentStage (1–3).
func (s DevelopmentStage) IsValid() bool {
	return s >= StagePreCapitalist && s <= StageSubordinated
}

// String returns the snake_case name of the stage.
func (s DevelopmentStage) String() string {
	switch s {
	case StagePreCapitalist:
		return "pre_capitalist"
	case StageTransition:
		return "transition"
	case StageSubordinated:
		return "subordinated"
	default:
		return fmt.Sprintf("DevelopmentStage(%d)", int(s))
	}
}

// ── HistoricalMerchantCapitalID ──────────────────────────────────────────────

// HistoricalMerchantCapitalID identifies a persisted historical-merchant-capital
// record. 96-bit hex from crypto/rand, matching every other domain ID in the
// simulator.
type HistoricalMerchantCapitalID string

// NewHistoricalMerchantCapitalID returns a fresh random HistoricalMerchantCapitalID.
func NewHistoricalMerchantCapitalID() HistoricalMerchantCapitalID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("merchant: rand: %w", err))
	}
	return HistoricalMerchantCapitalID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty HistoricalMerchantCapitalID.
func (id HistoricalMerchantCapitalID) IsZero() bool { return id == "" }

// ── HistoricalMerchantCapital ────────────────────────────────────────────────

// HistoricalMerchantCapital is the persisted entity for a historical form of
// merchant capital (Vol. III Ch. 20). It records the development stage, profit
// source, whether wage-labour is employed, and the subordination index in basis
// points [0, 10000] (merchant capital ÷ total social capital × 10000).
//
// Invariant 1: StagePreCapitalist && WageLabour == true is forbidden.
// Invariant 2: SubordinationIndex rises with advancing Stage (checked across
// records by InverseDevelopmentRelation).
type HistoricalMerchantCapital struct {
	ID                HistoricalMerchantCapitalID `json:"id"`
	Name              string                      `json:"name"`               // e.g. "Venice carrying trade"
	Stage             DevelopmentStage            `json:"stage"`              // 1–3
	ProfitSource      string                      `json:"profit_source"`      // unequal_exchange|monopoly_trade|carrying_monopoly|industrial_subordination
	WageLabour        bool                        `json:"wage_labour"`        // false for pre-capitalist
	SubordinationIndex int64                      `json:"subordination_index"` // bp [0, 10000]
	CreatedAt         time.Time                   `json:"created_at"`
}

// NewHistoricalMerchantCapital constructs a HistoricalMerchantCapital. No ID
// or timestamp is assigned — the store does that on persistence.
//
//	!stage.IsValid()                           → ErrInvalidDevelopmentStage
//	subordinationIndex < 0 or > 10000          → ErrSubordinationOutOfRange
//	stage == StagePreCapitalist && wageLabour  → ErrWageLabourInPreCapitalist
func NewHistoricalMerchantCapital(
	name string,
	stage DevelopmentStage,
	profitSource string,
	wageLabour bool,
	subordinationIndex int64,
) (HistoricalMerchantCapital, error) {
	if !stage.IsValid() {
		return HistoricalMerchantCapital{}, ErrInvalidDevelopmentStage
	}
	if subordinationIndex < 0 || subordinationIndex > basisPoints {
		return HistoricalMerchantCapital{}, ErrSubordinationOutOfRange
	}
	if stage == StagePreCapitalist && wageLabour {
		return HistoricalMerchantCapital{}, ErrWageLabourInPreCapitalist
	}
	return HistoricalMerchantCapital{
		Name:               name,
		Stage:              stage,
		ProfitSource:       profitSource,
		WageLabour:         wageLabour,
		SubordinationIndex: subordinationIndex,
	}, nil
}

// ── StagePoint ───────────────────────────────────────────────────────────────

// StagePoint is a helper value object pairing a DevelopmentStage with a
// SubordinationIndex for use in InverseDevelopmentRelation.
type StagePoint struct {
	Stage DevelopmentStage `json:"stage"`
	Index int64            `json:"index"` // bp [0, 10000]
}

// ── InverseDevelopmentRelation ───────────────────────────────────────────────

// InverseDevelopmentRelation is the value object verifying the inverse-development
// law (Vol. III Ch. 20): as industrial capitalism advances (rising stage), the
// SubordinationIndex of merchant capital rises — i.e. merchant capital is drawn
// ever more tightly into the orbit of industrial capital.
//
// Monotonic is true iff the SubordinationIndices strictly increase as stages
// advance (invariant 2).
type InverseDevelopmentRelation struct {
	Stages              []DevelopmentStage `json:"stages"`
	SubordinationIndices []int64           `json:"subordination_indices"` // bp; parallel to Stages
	Monotonic           bool               `json:"monotonic"`             // true iff indices strictly increasing
}

// NewInverseDevelopmentRelation builds an InverseDevelopmentRelation from a
// slice of StagePoints. It sorts by stage, validates each stage and index, then
// computes Monotonic.
//
//	any !stage.IsValid()                   → ErrInvalidDevelopmentStage
//	any index < 0 or index > 10000         → ErrSubordinationOutOfRange
func NewInverseDevelopmentRelation(points []StagePoint) (InverseDevelopmentRelation, error) {
	for _, p := range points {
		if !p.Stage.IsValid() {
			return InverseDevelopmentRelation{}, ErrInvalidDevelopmentStage
		}
		if p.Index < 0 || p.Index > basisPoints {
			return InverseDevelopmentRelation{}, ErrSubordinationOutOfRange
		}
	}

	// Sort by stage ascending so Monotonic reflects stage order.
	sorted := make([]StagePoint, len(points))
	copy(sorted, points)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Stage < sorted[j].Stage
	})

	stages := make([]DevelopmentStage, len(sorted))
	indices := make([]int64, len(sorted))
	for i, p := range sorted {
		stages[i] = p.Stage
		indices[i] = p.Index
	}

	monotonic := true
	for i := 1; i < len(indices); i++ {
		if indices[i] <= indices[i-1] {
			monotonic = false
			break
		}
	}

	return InverseDevelopmentRelation{
		Stages:               stages,
		SubordinationIndices: indices,
		Monotonic:            monotonic,
	}, nil
}

// ── HistoricalTradingSystem ──────────────────────────────────────────────────

// HistoricalTradingSystem is a named historical exemplar of merchant capital
// (e.g. Venice, Genoa, the Dutch carrying trade). It is a value object — not
// persisted — used in API responses and tests to represent the concrete
// historical forms Marx enumerates in Ch. 20.
type HistoricalTradingSystem struct {
	Name         string           `json:"name"`
	Stage        DevelopmentStage `json:"stage"`
	ProfitSource string           `json:"profit_source"`
	WageLabour   bool             `json:"wage_labour"`
}

// NewHistoricalTradingSystem builds and validates a HistoricalTradingSystem.
//
//	!stage.IsValid()                           → ErrInvalidDevelopmentStage
//	stage == StagePreCapitalist && wageLabour  → ErrWageLabourInPreCapitalist
func NewHistoricalTradingSystem(
	name string,
	stage DevelopmentStage,
	profitSource string,
	wageLabour bool,
) (HistoricalTradingSystem, error) {
	if !stage.IsValid() {
		return HistoricalTradingSystem{}, ErrInvalidDevelopmentStage
	}
	if stage == StagePreCapitalist && wageLabour {
		return HistoricalTradingSystem{}, ErrWageLabourInPreCapitalist
	}
	return HistoricalTradingSystem{
		Name:         name,
		Stage:        stage,
		ProfitSource: profitSource,
		WageLabour:   wageLabour,
	}, nil
}

// ── PreCapitalistForm ────────────────────────────────────────────────────────

// PreCapitalistForm is the value object expressing invariant 1 structurally:
// pre-capitalist merchant capital profits by unequal exchange, not wage-labour
// exploitation. ExploitsWageLabour is always false; the constructor enforces it.
type PreCapitalistForm struct {
	ProfitSource       string `json:"profit_source"`
	ExploitsWageLabour bool   `json:"exploits_wage_labour"` // always false
}

// NewPreCapitalistForm builds a PreCapitalistForm, setting ExploitsWageLabour=false
// regardless of any caller intent.
func NewPreCapitalistForm(profitSource string) (PreCapitalistForm, error) {
	return PreCapitalistForm{
		ProfitSource:       profitSource,
		ExploitsWageLabour: false,
	}, nil
}
