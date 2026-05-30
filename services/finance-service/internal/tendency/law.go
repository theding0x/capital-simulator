// Package tendency carries Capital Vol. III, Part III — the law of the
// tendential fall in the rate of profit. This file holds Ch. 13 ("The Law as
// Such").
//
// The law: with the growing productivity of social labour the organic
// composition of capital rises — constant capital (c) grows relative to
// variable capital (v). Variable capital is the sole source of surplus-value,
// so with the rate of surplus-value s′ held constant, the surplus-value s = s′·v
// no longer keeps pace with the total capital C = c + v. The rate of profit
//
//	p′ = s / (c + v)
//
// therefore falls, even though s′ is unchanged or rising. This is not a
// conjunctural fluctuation but a structural tendency of capitalist accumulation.
//
// A second, dialectically opposed movement accompanies the falling rate: the
// MASS of profit need not fall. As the rate falls the total capital grows, and
// the absolute mass of profit (C·p′) can stay constant or even rise. The
// RateMassContradiction value object exhibits this — the rate falls while the
// mass is preserved or grows. At the limit, ABSOLUTE over-accumulation: a fresh
// increment of capital sets no additional living labour in motion and yields
// zero (or negative) additional surplus-value.
package tendency

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// basisPoints is the fixed-point scale: 1% = 100 bp, so s′ of 100% = 10000 bp.
// Integer basis points keep rate computations exact across Marx's worked examples.
const basisPoints = 10000

// MassOfProfit is the absolute mass of profit in LabourMinutes — the total
// surplus-value appropriated, equal to C·p′. Distinct from the rate p′, it is
// the magnitude that may grow even as the rate falls.
type MassOfProfit int64

// Sentinel errors for invariant violations in this package.
var (
	ErrNegativeMagnitude = errors.New("tendency: magnitude must be non-negative")
	ErrEmptyTrajectory   = errors.New("tendency: trajectory must have at least one period")
)

// profitRateBP computes p′ = s/(c+v) in basis points, rounded half-up. Rounding
// matches the house convention used by avgprofit.individualProfitRate and the
// wage-effect rates: Period 1 (s=100, c+v=150) → 1,000,000/150 = 6666.67 → 6667.
// Returns 0 when the denominator is non-positive (guards divide-by-zero).
func profitRateBP(s, c, v int64) int64 {
	denom := c + v
	if denom <= 0 {
		return 0
	}
	return (s*basisPoints + denom/2) / denom
}

// CompositionTrajectoryID identifies a persisted composition trajectory.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type CompositionTrajectoryID string

// NewCompositionTrajectoryID returns a fresh random CompositionTrajectoryID.
func NewCompositionTrajectoryID() CompositionTrajectoryID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("tendency: rand: %w", err))
	}
	return CompositionTrajectoryID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CompositionTrajectoryID.
func (id CompositionTrajectoryID) IsZero() bool { return id == "" }

// RateMassContradictionID identifies a persisted rate-mass contradiction.
// 96-bit hex from crypto/rand.
type RateMassContradictionID string

// NewRateMassContradictionID returns a fresh random RateMassContradictionID.
func NewRateMassContradictionID() RateMassContradictionID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("tendency: rand: %w", err))
	}
	return RateMassContradictionID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty RateMassContradictionID.
func (id RateMassContradictionID) IsZero() bool { return id == "" }

// TrajectoryPeriod is one time-step in a rising-composition trajectory: a value
// object, not persisted independently (it lives inside CompositionTrajectory).
type TrajectoryPeriod struct {
	ConstantCapital int64 `json:"constant_capital"` // c in LabourMinutes
	VariableCapital int64 `json:"variable_capital"` // v in LabourMinutes
	SurplusValue    int64 `json:"surplus_value"`    // s = s′·v/10000; LabourMinutes
	ProfitRate      int64 `json:"profit_rate"`      // p′ = s/(c+v) in basis points
}

// CompositionTrajectory is the persisted entity capturing the law's trajectory:
// s′ held constant, c rising, v constant → p′ falling. Periods are serialised as
// JSON in the MySQL row; ProfitRates is always DERIVED from Periods on read and
// never stored separately.
type CompositionTrajectory struct {
	ID               CompositionTrajectoryID `json:"id"`
	Label            string                  `json:"label"`              // e.g. "Marx §law s'=100%"
	SurplusValueRate int64                   `json:"surplus_value_rate"` // s′ in basis points (10000 = 100%)
	Periods          []TrajectoryPeriod      `json:"periods"`
	ProfitRates      []int64                 `json:"profit_rates"` // derived: ProfitRate of each period
	CreatedAt        time.Time               `json:"created_at"`
}

// Validate enforces the invariants on a trajectory: it must have at least one
// period, and no magnitude may be negative.
func (t CompositionTrajectory) Validate() error {
	if len(t.Periods) == 0 {
		return ErrEmptyTrajectory
	}
	if t.SurplusValueRate < 0 {
		return ErrNegativeMagnitude
	}
	for _, p := range t.Periods {
		if p.ConstantCapital < 0 || p.VariableCapital < 0 || p.SurplusValue < 0 {
			return ErrNegativeMagnitude
		}
	}
	return nil
}

// DeriveProfitRates returns the ProfitRate of every period, in order. It is the
// single source of truth for the ProfitRates slice, recomputed on every read.
func (t CompositionTrajectory) DeriveProfitRates() []int64 {
	rates := make([]int64, 0, len(t.Periods))
	for _, p := range t.Periods {
		rates = append(rates, p.ProfitRate)
	}
	return rates
}

// FallingProfitRateLaw is the single-period expression of the law: with constant
// s′, a rising c/v composition causes a falling p′. It is a value object,
// computed on demand and not persisted.
type FallingProfitRateLaw struct {
	ConstantCapital  int64 `json:"constant_capital"`   // c; LabourMinutes
	VariableCapital  int64 `json:"variable_capital"`   // v; LabourMinutes
	SurplusValue     int64 `json:"surplus_value"`      // s = s′·v; LabourMinutes
	SurplusValueRate int64 `json:"surplus_value_rate"` // s′ in basis points
	ProfitRate       int64 `json:"profit_rate"`        // p′ = s/(c+v) in basis points
}

// ComputeFallingProfitRateLaw derives the single-period law from c, v and s′.
// s = surplusValueRateBP·v/10000; p′ = round_half_up(s·10000/(c+v)).
// For all c > 0 the result satisfies ProfitRate < SurplusValueRate.
func ComputeFallingProfitRateLaw(c, v, surplusValueRateBP int64) FallingProfitRateLaw {
	s := surplusValueRateBP * v / basisPoints
	return FallingProfitRateLaw{
		ConstantCapital:  c,
		VariableCapital:  v,
		SurplusValue:     s,
		SurplusValueRate: surplusValueRateBP,
		ProfitRate:       profitRateBP(s, c, v),
	}
}

// BuildCompositionTrajectory constructs a CompositionTrajectory from a label, a
// constant s′ (in basis points), and a slice of (c, v) pairs. For each pair it
// computes s = surplusValueRateBP·v/10000 and p′ = round_half_up(s·10000/(c+v)),
// then fills both Periods and the derived ProfitRates slice. It assigns no ID or
// timestamp — the store does that on persistence.
func BuildCompositionTrajectory(label string, surplusValueRateBP int64, steps [][2]int64) CompositionTrajectory {
	periods := make([]TrajectoryPeriod, 0, len(steps))
	rates := make([]int64, 0, len(steps))
	for _, step := range steps {
		c, v := step[0], step[1]
		s := surplusValueRateBP * v / basisPoints
		rate := profitRateBP(s, c, v)
		periods = append(periods, TrajectoryPeriod{
			ConstantCapital: c,
			VariableCapital: v,
			SurplusValue:    s,
			ProfitRate:      rate,
		})
		rates = append(rates, rate)
	}
	return CompositionTrajectory{
		Label:            label,
		SurplusValueRate: surplusValueRateBP,
		Periods:          periods,
		ProfitRates:      rates,
	}
}

// RateMassContradiction is the persisted entity exhibiting the law's internal
// contradiction: as the rate of profit falls, the absolute mass of profit need
// not fall. Mass = C·rate/10000; MassChange = NewMass − OldMass (signed).
type RateMassContradiction struct {
	ID         RateMassContradictionID `json:"id"`
	OldC       int64                   `json:"old_c"`       // old total capital C; LabourMinutes
	OldRate    int64                   `json:"old_rate"`    // old p′ in basis points
	NewC       int64                   `json:"new_c"`       // new total capital C; LabourMinutes
	NewRate    int64                   `json:"new_rate"`    // new p′ in basis points
	OldMass    int64                   `json:"old_mass"`    // OldC·OldRate/10000; LabourMinutes
	NewMass    int64                   `json:"new_mass"`    // NewC·NewRate/10000; LabourMinutes
	MassChange int64                   `json:"mass_change"` // NewMass − OldMass (signed)
	CreatedAt  time.Time               `json:"created_at"`
}

// Validate returns ErrNegativeMagnitude if any input magnitude is negative.
func (r RateMassContradiction) Validate() error {
	if r.OldC < 0 || r.OldRate < 0 || r.NewC < 0 || r.NewRate < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// ComputeRateMassContradiction derives the masses and their signed change from
// the two (C, rate) states. Mass = C·rate/10000 (integer division). It assigns
// no ID or timestamp — the store does that on persistence.
func ComputeRateMassContradiction(oldC, oldRate, newC, newRate int64) RateMassContradiction {
	oldMass := oldC * oldRate / basisPoints
	newMass := newC * newRate / basisPoints
	return RateMassContradiction{
		OldC:       oldC,
		OldRate:    oldRate,
		NewC:       newC,
		NewRate:    newRate,
		OldMass:    oldMass,
		NewMass:    newMass,
		MassChange: newMass - oldMass,
	}
}

// AbsoluteOveraccumulation is the limit-case value object: a fresh increment of
// capital (DeltaC) that yields no additional surplus-value (DeltaSurplusValue
// ≤ 0). It is computed on demand and not persisted.
type AbsoluteOveraccumulation struct {
	DeltaC            int64 `json:"delta_c"`             // additional capital advanced; LabourMinutes
	DeltaSurplusValue int64 `json:"delta_surplus_value"` // additional surplus-value (may be 0 or negative)
	IsAbsolute        bool  `json:"is_absolute"`         // DeltaSurplusValue <= 0
}

// ComputeAbsoluteOveraccumulation reports whether an advance of additional
// capital is absolute over-accumulation: it is, exactly when the additional
// surplus-value it sets in motion is zero or negative.
func ComputeAbsoluteOveraccumulation(deltaC, deltaSurplusValue int64) AbsoluteOveraccumulation {
	return AbsoluteOveraccumulation{
		DeltaC:            deltaC,
		DeltaSurplusValue: deltaSurplusValue,
		IsAbsolute:        deltaSurplusValue <= 0,
	}
}
