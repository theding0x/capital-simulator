// This file holds Capital Vol. III, Part III, Ch. 15 ("Exposition of the
// Internal Contradictions of the Law"). It COMPLETES Part III.
//
// The chapter has four sections:
//
//	§I   The falling rate of profit and the rising mass of profit coexist —
//	     they do not cancel; they express the same underlying movement (rising
//	     composition) from two sides.
//	§II  Over-production as the form in which the contradiction appears in the
//	     sphere of circulation: excess commodities cannot be realised at value.
//	§III Over-accumulation / relative over-population: more capital than can
//	     be deployed at the average rate of profit; the reserve army grows.
//	§IV  Concentration and centralisation of capital: the fall of the rate
//	     accelerates both (small capitals squeezed out; credit centralises).
package tendency

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Sentinel errors for invariant violations in Ch. 15.
var (
	ErrInvalidContradictionKind = errors.New("tendency: invalid contradiction kind")
	ErrInvalidWritedown         = errors.New("tendency: constant_capital_writedown must be between 1 and 99 inclusive")
	ErrInvalidPeriods           = errors.New("tendency: periods must be at least 1")
)

// ContradictionKind enumerates the four faces of the law's internal
// contradiction Marx exposes in Vol. III Ch. 15, §I–§IV.
type ContradictionKind string

const (
	// KindFallingRateRisingMass — §I: the falling rate and rising mass of
	// profit are two sides of the same movement.
	KindFallingRateRisingMass ContradictionKind = "falling_rate_rising_mass"
	// KindOverproductionUnderConsumption — §II: rising productivity expands the
	// mass of commodities beyond the limits of consumption.
	KindOverproductionUnderConsumption ContradictionKind = "overproduction_underconsumption"
	// KindOveraccumulationRelativeOverpopulation — §III: excess capital cannot be
	// deployed at the average rate of profit; the surplus population grows.
	KindOveraccumulationRelativeOverpopulation ContradictionKind = "overaccumulation_relative_overpopulation"
	// KindConcentrationCentralisation — §IV: the fall in the rate accelerates the
	// concentration and centralisation of capital.
	KindConcentrationCentralisation ContradictionKind = "concentration_centralisation"
)

// IsValid reports whether k is one of the four named contradiction kinds.
func (k ContradictionKind) IsValid() bool {
	switch k {
	case KindFallingRateRisingMass,
		KindOverproductionUnderConsumption,
		KindOveraccumulationRelativeOverpopulation,
		KindConcentrationCentralisation:
		return true
	default:
		return false
	}
}

// Description returns a one-line Marx-faithful gloss for each kind, or the empty
// string for an unknown kind.
func (k ContradictionKind) Description() string {
	switch k {
	case KindFallingRateRisingMass:
		return "The falling rate and rising mass of profit are two sides of the same movement: both express the tendential fall in the rate of profit."
	case KindOverproductionUnderConsumption:
		return "Rising productivity expands the mass of commodities beyond the limits of consumption, producing periodic crises of over-production."
	case KindOveraccumulationRelativeOverpopulation:
		return "Excess capital cannot be deployed at the average rate of profit; the surplus population grows as the organic composition rises."
	case KindConcentrationCentralisation:
		return "The fall in the rate of profit accelerates the concentration and centralisation of capital: small capitals are squeezed out and swallowed by large ones."
	default:
		return ""
	}
}

// InternalContradictionID identifies a persisted internal contradiction.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type InternalContradictionID string

// NewInternalContradictionID returns a fresh random InternalContradictionID.
func NewInternalContradictionID() InternalContradictionID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("tendency: rand: %w", err))
	}
	return InternalContradictionID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty InternalContradictionID.
func (id InternalContradictionID) IsZero() bool { return id == "" }

// InternalContradiction is a persisted record of one face of the law's internal
// contradiction (Vol. III Ch. 15). Every contradiction in this chapter is
// dialectically paired — IsCoexistent is true for all four valid kinds.
type InternalContradiction struct {
	ID           InternalContradictionID `json:"id"`
	Kind         ContradictionKind       `json:"kind"`
	IsCoexistent bool                    `json:"is_coexistent"`
	Note         string                  `json:"note"`
	CreatedAt    time.Time               `json:"created_at"`
}

// NewInternalContradiction constructs an internal contradiction from a kind and
// an optional note. It returns ErrInvalidContradictionKind if kind is not one of
// the four named kinds, and sets IsCoexistent = true for every valid kind: the
// falling rate and rising mass are coexistent by §I, and the others are
// coexistent manifestations of the same tendency. It assigns no ID or
// timestamp — the store does that on persistence.
func NewInternalContradiction(kind ContradictionKind, note string) (InternalContradiction, error) {
	if !kind.IsValid() {
		return InternalContradiction{}, ErrInvalidContradictionKind
	}
	return InternalContradiction{
		Kind:         kind,
		IsCoexistent: true,
		Note:         note,
	}, nil
}

// Overproduction is the §II value object: the mass of commodities outruns the
// value that can be realised in the market. Computed on demand, not persisted.
type Overproduction struct {
	ExcessCommodityValue int64 `json:"excess_commodity_value"` // LabourMinutes
	RealisableValue      int64 `json:"realisable_value"`       // LabourMinutes
}

// UnrealisedValue is the value that cannot be realised: the excess of commodity
// value over what the market can absorb (ExcessCommodityValue − RealisableValue).
func (o Overproduction) UnrealisedValue() int64 {
	return o.ExcessCommodityValue - o.RealisableValue
}

// Validate returns ErrNegativeMagnitude if either field is negative.
func (o Overproduction) Validate() error {
	if o.ExcessCommodityValue < 0 || o.RealisableValue < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// ComputeOverproduction builds and validates an Overproduction value object.
func ComputeOverproduction(excessCommodityValue, realisableValue int64) (Overproduction, error) {
	o := Overproduction{
		ExcessCommodityValue: excessCommodityValue,
		RealisableValue:      realisableValue,
	}
	if err := o.Validate(); err != nil {
		return Overproduction{}, err
	}
	return o, nil
}

// Overaccumulation is the §III value object: capital advanced in excess of what
// can be deployed at the average rate of profit. When the actual rate on the
// excess falls to zero or below, the over-accumulation is absolute. Computed on
// demand, not persisted.
type Overaccumulation struct {
	ExcessCapital     int64 `json:"excess_capital"`      // LabourMinutes; capital in excess
	AverageProfitRate int64 `json:"average_profit_rate"` // p̄' in basis points
	ActualProfitRate  int64 `json:"actual_profit_rate"`  // actual p' on excess capital; bp
}

// IsAbsolute reports absolute over-accumulation: the fresh capital sets no
// additional surplus-value in motion, so its actual rate of profit is ≤ 0.
func (o Overaccumulation) IsAbsolute() bool {
	return o.ActualProfitRate <= 0
}

// ShortfallBP is the signed gap between the average rate and the actual rate on
// the excess capital (AverageProfitRate − ActualProfitRate); positive = shortfall.
func (o Overaccumulation) ShortfallBP() int64 {
	return o.AverageProfitRate - o.ActualProfitRate
}

// Validate returns ErrNegativeMagnitude if the excess capital or either rate is negative.
func (o Overaccumulation) Validate() error {
	if o.ExcessCapital < 0 || o.AverageProfitRate < 0 || o.ActualProfitRate < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// ComputeOveraccumulation builds and validates an Overaccumulation value object.
func ComputeOveraccumulation(excessCapital, averageProfitRate, actualProfitRate int64) (Overaccumulation, error) {
	o := Overaccumulation{
		ExcessCapital:     excessCapital,
		AverageProfitRate: averageProfitRate,
		ActualProfitRate:  actualProfitRate,
	}
	if err := o.Validate(); err != nil {
		return Overaccumulation{}, err
	}
	return o, nil
}

// CapitalConcentration is the §IV value object for the concentration of capital:
// a single capital grows by accumulating a fraction of itself each period. Each
// period C ← C + C·SurplusValueRate/10000 (integer arithmetic, truncation per
// step). Computed on demand, not persisted.
type CapitalConcentration struct {
	InitialC         int64 `json:"initial_c"`          // LabourMinutes
	SurplusValueRate int64 `json:"surplus_value_rate"` // s' in basis points; fraction accumulated each period
	Periods          int64 `json:"periods"`            // number of accumulation periods
}

// FinalC is the capital after every accumulation period, growing geometrically
// at SurplusValueRate per period.
func (c CapitalConcentration) FinalC() int64 {
	capital := c.InitialC
	for i := int64(0); i < c.Periods; i++ {
		capital += capital * c.SurplusValueRate / basisPoints
	}
	return capital
}

// Path returns the per-period capital values, with InitialC as element [0] and
// one entry appended for each of the Periods accumulation steps.
func (c CapitalConcentration) Path() []int64 {
	path := make([]int64, 0, c.Periods+1)
	capital := c.InitialC
	path = append(path, capital)
	for i := int64(0); i < c.Periods; i++ {
		capital += capital * c.SurplusValueRate / basisPoints
		path = append(path, capital)
	}
	return path
}

// Validate returns ErrNegativeMagnitude if any field is negative, and
// ErrInvalidPeriods if Periods < 1.
func (c CapitalConcentration) Validate() error {
	if c.InitialC < 0 || c.SurplusValueRate < 0 || c.Periods < 0 {
		return ErrNegativeMagnitude
	}
	if c.Periods < 1 {
		return ErrInvalidPeriods
	}
	return nil
}

// ComputeCapitalConcentration builds and validates a CapitalConcentration.
func ComputeCapitalConcentration(initialC, surplusValueRate, periods int64) (CapitalConcentration, error) {
	c := CapitalConcentration{
		InitialC:         initialC,
		SurplusValueRate: surplusValueRate,
		Periods:          periods,
	}
	if err := c.Validate(); err != nil {
		return CapitalConcentration{}, err
	}
	return c, nil
}

// CapitalCentralisation is the §IV value object for the centralisation of
// capital: many existing capitals are absorbed into one. The result is the sum
// of the absorbed capitals (AbsorbedCapitals · AbsorbedSize). Unlike
// concentration, centralisation creates no new capital — it merely redistributes
// the existing mass. Computed on demand, not persisted.
type CapitalCentralisation struct {
	AbsorbedCapitals int64 `json:"absorbed_capitals"` // count of capitals absorbed
	AbsorbedSize     int64 `json:"absorbed_size"`     // average size of each; LabourMinutes
}

// ResultingCapital is the centralised capital: the count of absorbed capitals
// times their average size.
func (c CapitalCentralisation) ResultingCapital() int64 {
	return c.AbsorbedCapitals * c.AbsorbedSize
}

// Validate returns ErrNegativeMagnitude if either field is negative.
func (c CapitalCentralisation) Validate() error {
	if c.AbsorbedCapitals < 0 || c.AbsorbedSize < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// ComputeCapitalCentralisation builds and validates a CapitalCentralisation.
func ComputeCapitalCentralisation(absorbedCapitals, absorbedSize int64) (CapitalCentralisation, error) {
	c := CapitalCentralisation{
		AbsorbedCapitals: absorbedCapitals,
		AbsorbedSize:     absorbedSize,
	}
	if err := c.Validate(); err != nil {
		return CapitalCentralisation{}, err
	}
	return c, nil
}

// CrisisID identifies a persisted crisis record.
// 96-bit hex from crypto/rand.
type CrisisID string

// NewCrisisID returns a fresh random CrisisID.
func NewCrisisID() CrisisID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("tendency: rand: %w", err))
	}
	return CrisisID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CrisisID.
func (id CrisisID) IsZero() bool { return id == "" }

// Crisis is the persisted entity exhibiting the violent resolution of the law's
// contradiction (Vol. III Ch. 15): a crisis devalues (writes down) a percentage
// of constant capital, which shrinks the denominator of p′ = s/(c+v) and so
// restores the rate of profit. PostCrisisProfitRate is derived and always
// exceeds PreCrisisProfitRate for a writedown in (0,100) and pre > 0.
type Crisis struct {
	ID                       CrisisID  `json:"id"`
	ConstantCapitalWritedown int64     `json:"constant_capital_writedown"` // percent integer; 1–99
	PreCrisisProfitRate      int64     `json:"pre_crisis_profit_rate"`     // bp
	PostCrisisProfitRate     int64     `json:"post_crisis_profit_rate"`    // bp; derived; always > Pre
	CreatedAt                time.Time `json:"created_at"`
}

// restoredProfitRate computes the post-crisis rate of profit, round-half-up.
// Writing down a fraction w of constant capital shrinks the denominator from
// 10000 bp to (10000 − w·100) bp, lifting the rate by 10000/(10000 − w·100):
//
//	post = roundHalfUp(pre·10000 / (10000 − writedown·100))
//
// For writedown ∈ (0,100) the denominator lies in (0, 10000), so post > pre.
func restoredProfitRate(preCrisisProfitRate, constantCapitalWritedown int64) int64 {
	denom := basisPoints - constantCapitalWritedown*100
	if denom <= 0 {
		return preCrisisProfitRate
	}
	numer := preCrisisProfitRate * basisPoints
	return (numer + denom/2) / denom
}

// NewCrisis constructs a crisis from a constant-capital writedown (percent) and
// a pre-crisis rate of profit (basis points), deriving the restored post-crisis
// rate. It returns ErrInvalidWritedown if writedown is not in (0,100) exclusive,
// and ErrNegativeMagnitude if the pre-crisis rate is negative. It assigns no ID
// or timestamp — the store does that on persistence.
func NewCrisis(constantCapitalWritedown, preCrisisProfitRate int64) (Crisis, error) {
	if constantCapitalWritedown <= 0 || constantCapitalWritedown >= 100 {
		return Crisis{}, ErrInvalidWritedown
	}
	if preCrisisProfitRate < 0 {
		return Crisis{}, ErrNegativeMagnitude
	}
	return Crisis{
		ConstantCapitalWritedown: constantCapitalWritedown,
		PreCrisisProfitRate:      preCrisisProfitRate,
		PostCrisisProfitRate:     restoredProfitRate(preCrisisProfitRate, constantCapitalWritedown),
	}, nil
}

// PartIIISummary is the value object that closes Part III: the running counts of
// every Part III entity type and the latest record of each. Assembled by the
// summary handler from the store's list calls; never persisted. Pointer fields
// are nil when the corresponding list is empty, so the JSON omits them.
type PartIIISummary struct {
	TrajectoryCount     int64                  `json:"trajectory_count"`
	ScenarioCount       int64                  `json:"scenario_count"`
	CrisisCount         int64                  `json:"crisis_count"`
	ContradictionCount  int64                  `json:"contradiction_count"`
	LatestTrajectory    *CompositionTrajectory `json:"latest_trajectory,omitempty"`
	LatestScenario      *CounteractingScenario `json:"latest_scenario,omitempty"`
	LatestCrisis        *Crisis                `json:"latest_crisis,omitempty"`
	LatestContradiction *InternalContradiction `json:"latest_contradiction,omitempty"`
}
