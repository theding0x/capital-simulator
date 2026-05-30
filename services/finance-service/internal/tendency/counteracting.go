// This file holds Capital Vol. III, Part III, Ch. 14 ("Counteracting
// Influences"). Ch. 13 established the law as such: a rising organic
// composition drives the general rate of profit down. This chapter exhibits the
// six counter-tendencies Marx enumerates — they do not abolish the law, they
// hamper, retard, and partly paralyse the fall.
//
// The six counteracting forces (§I–§VI):
//
//	I.   intensified exploitation  — more surplus-value from a given labour-force
//	II.  depressed wages           — wages forced below the value of labour-power
//	III. cheaper constant capital  — rising productivity cheapens the elements of c
//	IV.  relative overpopulation   — the reserve army keeps wages low, new branches low-comp
//	V.   foreign trade             — cheaper imported inputs and a higher rate of surplus-value
//	VI.  stock capital             — interest-bearing capital excluded from the general rate
//
// The five industrial forces (§I–§V) slow the fall of the rate of profit along a
// rising-composition trajectory; the sixth (§VI) is excluded from the general
// rate entirely, so it contributes no uplift to the industrial path. Even with
// all five industrial forces active, the law persists: the final rate of profit
// remains below the initial rate.
package tendency

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Sentinel errors for invariant violations in Ch. 14.
var (
	ErrInvalidForceKind        = errors.New("tendency: invalid counteracting force kind")
	ErrModifiedRateNotHigher   = errors.New("tendency: modified profit rate must exceed baseline")
	ErrNewRateNotHigher        = errors.New("tendency: new profit rate must exceed old when constant capital is cheapened")
	ErrInvalidCheapeningFactor = errors.New("tendency: cheapening factor must be between 1 and 100 inclusive")
	ErrEmptyForces             = errors.New("tendency: scenario must have at least one counteracting force")
	ErrLawNotPreserved         = errors.New("tendency: final profit rate must remain below initial rate even with counteracting forces")
)

// CounteractingForceKind enumerates Marx's six counter-tendencies to the falling
// rate of profit (Vol. III Ch. 14, §I–§VI).
type CounteractingForceKind string

const (
	// KindIntensifiedExploitation — §I: more surplus-value from a given labour-force.
	KindIntensifiedExploitation CounteractingForceKind = "intensified_exploitation"
	// KindDepressedWages — §II: wages forced below the value of labour-power.
	KindDepressedWages CounteractingForceKind = "depressed_wages"
	// KindCheaperConstantCapital — §III: rising productivity cheapens the elements of c.
	KindCheaperConstantCapital CounteractingForceKind = "cheaper_constant_capital"
	// KindRelativeOverpopulation — §IV: the reserve army keeps wages low; new branches low-comp.
	KindRelativeOverpopulation CounteractingForceKind = "relative_overpopulation"
	// KindForeignTrade — §V: cheaper imported inputs and a higher rate of surplus-value.
	KindForeignTrade CounteractingForceKind = "foreign_trade"
	// KindStockCapital — §VI: interest-bearing capital excluded from the general rate.
	KindStockCapital CounteractingForceKind = "stock_capital"
)

// IsValid reports whether k is one of the six named counteracting force kinds.
func (k CounteractingForceKind) IsValid() bool {
	switch k {
	case KindIntensifiedExploitation,
		KindDepressedWages,
		KindCheaperConstantCapital,
		KindRelativeOverpopulation,
		KindForeignTrade,
		KindStockCapital:
		return true
	default:
		return false
	}
}

// Description returns a one-line Marx-faithful gloss for each kind, or the empty
// string for an unknown kind.
func (k CounteractingForceKind) Description() string {
	switch k {
	case KindIntensifiedExploitation:
		return "More surplus-value is extracted from a given labour-force by intensifying exploitation above the average."
	case KindDepressedWages:
		return "Wages are forced below the value of labour-power, raising surplus-value without raising productivity."
	case KindCheaperConstantCapital:
		return "Rising productivity cheapens the elements of constant capital, so c grows in mass more slowly than in value."
	case KindRelativeOverpopulation:
		return "The reserve army of labour keeps wages low and supplies cheap labour to new, low-composition branches of production."
	case KindForeignTrade:
		return "Foreign trade cheapens both the elements of constant capital and the means of subsistence, and raises the rate of surplus-value."
	case KindStockCapital:
		return "Interest-bearing (stock) capital yields interest, not the average profit, and is excluded from the equalisation of the general rate."
	default:
		return ""
	}
}

// CounteractingForceID identifies a persisted counteracting force.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type CounteractingForceID string

// NewCounteractingForceID returns a fresh random CounteractingForceID.
func NewCounteractingForceID() CounteractingForceID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("tendency: rand: %w", err))
	}
	return CounteractingForceID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CounteractingForceID.
func (id CounteractingForceID) IsZero() bool { return id == "" }

// CounteractingForce is a single counter-tendency record (Vol. III Ch. 14).
// ExcludedFromGeneralRate is true only for KindStockCapital — interest-bearing
// capital does not enter the equalisation of the general rate of profit.
type CounteractingForce struct {
	ID                      CounteractingForceID   `json:"id"`
	Kind                    CounteractingForceKind `json:"kind"`
	ExcludedFromGeneralRate bool                   `json:"excluded_from_general_rate"`
	Note                    string                 `json:"note"`
	CreatedAt               time.Time              `json:"created_at"`
}

// NewCounteractingForce constructs a counteracting force from a kind and an
// optional note. It returns ErrInvalidForceKind if kind is not one of the six
// named kinds, and sets ExcludedFromGeneralRate = (kind == KindStockCapital).
// It assigns no ID or timestamp — the store does that on persistence.
func NewCounteractingForce(kind CounteractingForceKind, note string) (CounteractingForce, error) {
	if !kind.IsValid() {
		return CounteractingForce{}, ErrInvalidForceKind
	}
	return CounteractingForce{
		Kind:                    kind,
		ExcludedFromGeneralRate: kind == KindStockCapital,
		Note:                    note,
	}, nil
}

// ExploitationIntensityEffect is the §I value object: raising the rate of
// surplus-value on a fixed c and v lifts the rate of profit. Computed on demand,
// not persisted.
type ExploitationIntensityEffect struct {
	OldSRate int64 `json:"old_s_rate"` // old s' in basis points
	NewSRate int64 `json:"new_s_rate"` // new s' in basis points (> OldSRate)
	C        int64 `json:"c"`          // constant capital; LabourMinutes
	V        int64 `json:"v"`          // variable capital; LabourMinutes
}

// BaselineProfitRate is p′ at the old rate of surplus-value, in basis points.
func (e ExploitationIntensityEffect) BaselineProfitRate() int64 {
	return profitRateBP(e.OldSRate*e.V/basisPoints, e.C, e.V)
}

// ModifiedProfitRate is p′ at the new (higher) rate of surplus-value, in basis points.
func (e ExploitationIntensityEffect) ModifiedProfitRate() int64 {
	return profitRateBP(e.NewSRate*e.V/basisPoints, e.C, e.V)
}

// Validate enforces non-negativity and that the modified rate exceeds the baseline.
func (e ExploitationIntensityEffect) Validate() error {
	if e.OldSRate < 0 || e.NewSRate < 0 || e.C < 0 || e.V < 0 {
		return ErrNegativeMagnitude
	}
	if e.ModifiedProfitRate() <= e.BaselineProfitRate() {
		return ErrModifiedRateNotHigher
	}
	return nil
}

// ComputeExploitationIntensityEffect builds and validates an
// ExploitationIntensityEffect from the old and new s′ (basis points) and c, v.
func ComputeExploitationIntensityEffect(oldSRate, newSRate, c, v int64) (ExploitationIntensityEffect, error) {
	e := ExploitationIntensityEffect{
		OldSRate: oldSRate,
		NewSRate: newSRate,
		C:        c,
		V:        v,
	}
	if err := e.Validate(); err != nil {
		return ExploitationIntensityEffect{}, err
	}
	return e, nil
}

// ConstantCapitalCheapening is the §III value object: cheapening the elements of
// constant capital lowers the denominator C and so lifts the rate of profit.
// Computed on demand, not persisted.
type ConstantCapitalCheapening struct {
	OriginalC int64 `json:"original_c"` // original constant capital; LabourMinutes
	NewC      int64 `json:"new_c"`      // cheapened constant capital (< OriginalC); LabourMinutes
	V         int64 `json:"v"`          // variable capital; LabourMinutes
	S         int64 `json:"s"`          // surplus value; LabourMinutes
}

// OldProfitRate is p′ at the original constant capital, in basis points.
func (c ConstantCapitalCheapening) OldProfitRate() int64 {
	return profitRateBP(c.S, c.OriginalC, c.V)
}

// NewProfitRate is p′ at the cheapened constant capital, in basis points.
func (c ConstantCapitalCheapening) NewProfitRate() int64 {
	return profitRateBP(c.S, c.NewC, c.V)
}

// Validate enforces non-negativity and, when c is genuinely cheapened
// (NewC < OriginalC), that the new rate exceeds the old. No invariant beyond
// non-negativity is enforced when NewC >= OriginalC.
func (c ConstantCapitalCheapening) Validate() error {
	if c.OriginalC < 0 || c.NewC < 0 || c.V < 0 || c.S < 0 {
		return ErrNegativeMagnitude
	}
	if c.NewC < c.OriginalC && c.NewProfitRate() <= c.OldProfitRate() {
		return ErrNewRateNotHigher
	}
	return nil
}

// ComputeConstantCapitalCheapening builds and validates a
// ConstantCapitalCheapening from the original and cheapened c, plus v and s.
func ComputeConstantCapitalCheapening(originalC, newC, v, s int64) (ConstantCapitalCheapening, error) {
	c := ConstantCapitalCheapening{
		OriginalC: originalC,
		NewC:      newC,
		V:         v,
		S:         s,
	}
	if err := c.Validate(); err != nil {
		return ConstantCapitalCheapening{}, err
	}
	return c, nil
}

// ForeignTradeEffect is the §V value object: foreign trade cheapens the elements
// of constant capital (a factor < 100% on c) and raises the rate of
// surplus-value (a boost added to s′). Both movements lift the rate of profit;
// NetProfitRateEffect is the signed change in p′ (basis points). Computed on
// demand, not persisted.
type ForeignTradeEffect struct {
	ReferenceC          int64 `json:"reference_c"`            // reference constant capital; LabourMinutes
	ReferenceV          int64 `json:"reference_v"`            // reference variable capital; LabourMinutes
	ReferenceSRate      int64 `json:"reference_s_rate"`       // reference s' in basis points
	CheapeningFactor    int64 `json:"cheapening_factor"`      // integer percent: 80 means inputs cost 80% of before
	SRateBoost          int64 `json:"s_rate_boost"`           // basis points added to ReferenceSRate
	NetProfitRateEffect int64 `json:"net_profit_rate_effect"` // derived; signed bp; >0 counteracts the fall
}

// Validate enforces non-negativity and that the cheapening factor lies in (0, 100].
func (f ForeignTradeEffect) Validate() error {
	if f.ReferenceC < 0 || f.ReferenceV < 0 || f.ReferenceSRate < 0 || f.SRateBoost < 0 {
		return ErrNegativeMagnitude
	}
	if f.CheapeningFactor <= 0 || f.CheapeningFactor > 100 {
		return ErrInvalidCheapeningFactor
	}
	return nil
}

// ComputeForeignTradeEffect builds and validates a ForeignTradeEffect, deriving
// NetProfitRateEffect from the reference scenario, the cheapening factor on c,
// and the boost to s′.
func ComputeForeignTradeEffect(referenceC, referenceV, referenceSRate, cheapeningFactor, sRateBoost int64) (ForeignTradeEffect, error) {
	f := ForeignTradeEffect{
		ReferenceC:       referenceC,
		ReferenceV:       referenceV,
		ReferenceSRate:   referenceSRate,
		CheapeningFactor: cheapeningFactor,
		SRateBoost:       sRateBoost,
	}
	if err := f.Validate(); err != nil {
		return ForeignTradeEffect{}, err
	}
	baselineRate := profitRateBP(f.ReferenceSRate*f.ReferenceV/basisPoints, f.ReferenceC, f.ReferenceV)
	cheapenedC := f.ReferenceC * f.CheapeningFactor / 100
	boostedSRate := f.ReferenceSRate + f.SRateBoost
	boostedS := boostedSRate * f.ReferenceV / basisPoints
	modifiedRate := profitRateBP(boostedS, cheapenedC, f.ReferenceV)
	f.NetProfitRateEffect = modifiedRate - baselineRate
	return f, nil
}

// CounteractingScenarioID identifies a persisted counteracting scenario.
// 96-bit hex from crypto/rand.
type CounteractingScenarioID string

// NewCounteractingScenarioID returns a fresh random CounteractingScenarioID.
func NewCounteractingScenarioID() CounteractingScenarioID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("tendency: rand: %w", err))
	}
	return CounteractingScenarioID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CounteractingScenarioID.
func (id CounteractingScenarioID) IsZero() bool { return id == "" }

// CounteractingScenario is the persisted entity (Vol. III Ch. 14): a set of
// counteracting forces applied to a rising-composition base trajectory. The
// modified path slows the fall of the rate of profit but does not abolish it —
// FinalProfitRate remains below InitialProfitRate.
type CounteractingScenario struct {
	ID                CounteractingScenarioID `json:"id"`
	Label             string                  `json:"label"`
	Forces            []CounteractingForce    `json:"forces"`              // len >= 1; inline value objects
	BaseTrajectory    CompositionTrajectory   `json:"base_trajectory"`     // the underlying law trajectory (embedded value)
	InitialProfitRate int64                   `json:"initial_profit_rate"` // first period p' after uplift; bp
	FinalProfitRate   int64                   `json:"final_profit_rate"`   // last period p' after uplift; bp
	ModifiedRates     []int64                 `json:"modified_rates"`      // one bp value per period after uplift
	CreatedAt         time.Time               `json:"created_at"`
}

// Validate enforces that the scenario has at least one force, at least one
// trajectory period, and that the law persists: the final rate of profit must
// remain below the initial rate.
func (s CounteractingScenario) Validate() error {
	if len(s.Forces) == 0 {
		return ErrEmptyForces
	}
	if len(s.BaseTrajectory.Periods) == 0 {
		return ErrEmptyTrajectory
	}
	if s.FinalProfitRate >= s.InitialProfitRate {
		return ErrLawNotPreserved
	}
	return nil
}

// forceUplift returns the additive per-period uplift δ contributed by a single
// force, given the period's surplus-value s and total capital C = c + v. The
// uplift is denominator-normalised so it scales with the rate, not the absolute
// magnitudes. KindStockCapital contributes zero — it is excluded from the
// general rate and does not modify the industrial path.
func forceUplift(kind CounteractingForceKind, s, denom int64) int64 {
	if denom <= 0 {
		return 0
	}
	switch kind {
	case KindIntensifiedExploitation, KindDepressedWages:
		return s * 500 / denom
	case KindCheaperConstantCapital, KindForeignTrade:
		return s * 1000 / denom
	case KindRelativeOverpopulation:
		return s * 300 / denom
	default: // KindStockCapital and any unknown kind
		return 0
	}
}

// applyCounteractingForces computes the modified per-period rate of profit by
// adding every force's uplift to the base rate of each period.
func applyCounteractingForces(trajectory CompositionTrajectory, forces []CounteractingForce) []int64 {
	rates := make([]int64, 0, len(trajectory.Periods))
	for _, p := range trajectory.Periods {
		denom := p.ConstantCapital + p.VariableCapital
		base := profitRateBP(p.SurplusValue, p.ConstantCapital, p.VariableCapital)
		var totalUplift int64
		for _, f := range forces {
			totalUplift += forceUplift(f.Kind, p.SurplusValue, denom)
		}
		rates = append(rates, base+totalUplift)
	}
	return rates
}

// BuildCounteractingScenario applies the given forces to the base trajectory and
// returns the scenario. It returns ErrEmptyForces if forces is empty, and
// ErrLawNotPreserved if the modified path does not keep the final rate of profit
// below the initial rate. It assigns no ID or timestamp — the store does that on
// persistence.
func BuildCounteractingScenario(label string, trajectory CompositionTrajectory, forces []CounteractingForce) (CounteractingScenario, error) {
	if len(forces) == 0 {
		return CounteractingScenario{}, ErrEmptyForces
	}
	if len(trajectory.Periods) == 0 {
		return CounteractingScenario{}, ErrEmptyTrajectory
	}

	modified := applyCounteractingForces(trajectory, forces)
	scenario := CounteractingScenario{
		Label:             label,
		Forces:            forces,
		BaseTrajectory:    trajectory,
		InitialProfitRate: modified[0],
		FinalProfitRate:   modified[len(modified)-1],
		ModifiedRates:     modified,
	}
	if scenario.FinalProfitRate >= scenario.InitialProfitRate {
		return CounteractingScenario{}, ErrLawNotPreserved
	}
	return scenario, nil
}
