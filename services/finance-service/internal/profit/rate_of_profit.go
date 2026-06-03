package profit

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// This file carries Ch. 2's argument: the same surplus-value, measured against
// two different denominators, yields two different rates. Measured against the
// variable capital that alone produced it, s gives the rate of surplus-value
// s' = s/v — Vol. I's "degree of exploitation". Measured against the *total*
// capital advanced, the identical s gives the rate of profit p' = s/C = s/(c+v).
// Because C >= v, the rate of profit is always the smaller of the two, and the
// gap between them measures how thoroughly the profit-form conceals the origin
// of surplus-value in variable capital alone.

// basisPoints is the fixed-point scale for both rates: 1% = 100 bp, so a rate
// of surplus-value of 100% is 10000 bp. Integer basis points keep the
// magnitudes exact (Marx's worked rates — 20%, 100% — are whole percentages)
// and keep floating point out of the value layer.
const basisPoints = 10000

// roundHalfUp computes (numer + denom/2) / denom — the canonical round-half-up
// division used throughout the simulator for basis-point arithmetic (also in
// credit.roundHalfUp, merchant.roundHalfUp, rent.roundHalfUp). Callers guard
// denom > 0 before invoking. Capital III's worked rates are whole percentages,
// so the rounding only bites at the margins, but it keeps profit's basis points
// exactly consistent with the avgprofit general rate that bounds the rate of
// interest in Part V.
func roundHalfUp(numer, denom int64) int64 {
	return (numer + denom/2) / denom
}

// ProfitRateID identifies a stored rate-of-profit analysis. 96-bit hex from
// crypto/rand, like every other domain ID in the simulator.
type ProfitRateID string

// NewProfitRateID returns a fresh random ProfitRateID.
func NewProfitRateID() ProfitRateID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("profit: rand: %w", err))
	}
	return ProfitRateID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty ProfitRateID.
func (id ProfitRateID) IsZero() bool { return id == "" }

// TotalCapital is C = c + v: the whole capital advanced for production, and the
// denominator against which the rate of profit measures surplus-value. Where the
// rate of surplus-value divides by variable capital alone, the rate of profit
// divides by the total — which is precisely why it is the smaller rate.
type TotalCapital LabourMinutes

// RateOfSurplusValue is s' = s/v, expressed in basis points (100% = 10000 bp).
// It is the Vol. I measure — the degree of exploitation of labour-power —
// carried into Vol. III so it can be set beside the rate of profit.
type RateOfSurplusValue int64

// MystificationDegree quantifies how much the profit-form conceals the origin of
// surplus-value: (s' - p') / s' in basis points. It is 0 when c = 0 (the two
// rates coincide and nothing is hidden) and rises towards 10000 as constant
// capital comes to dominate the advance, dwarfing the variable capital that is
// the sole source of the surplus.
type MystificationDegree int64

// RateOfProfit is the central magnitude of this chapter: p' = s / C = s / (c+v),
// the surplus-value measured against the total capital advanced. It carries the
// surplus-value and the total capital that produced it so the inversion at the
// heart of the profit-form stays legible — the same surplus, divided by a larger
// denominator, yields a smaller and mystified rate.
type RateOfProfit struct {
	SurplusValue SurplusValue `json:"surplus_value"`
	TotalCapital TotalCapital `json:"total_capital"`
	BasisPoints  int64        `json:"basis_points"`
}

// ComputeRateOfProfit forms p' = s / C in basis points from the surplus-value
// and the capital decomposition. A zero total capital (c = v = 0) yields a zero
// rate: with nothing advanced there is no rate of self-expansion to express.
func ComputeRateOfProfit(s SurplusValue, c ConstantCapital, v VariableCapital) RateOfProfit {
	total := TotalCapital(int64(c) + int64(v))
	var bp int64
	if total > 0 {
		bp = roundHalfUp(int64(s)*basisPoints, int64(total))
	}
	return RateOfProfit{SurplusValue: s, TotalCapital: total, BasisPoints: bp}
}

// ComputeRateOfSurplusValue forms s' = s / v in basis points. A zero variable
// capital yields a zero rate: surplus-value originates in variable capital, so
// without it there is no rate of exploitation to measure.
func ComputeRateOfSurplusValue(s SurplusValue, v VariableCapital) RateOfSurplusValue {
	if v <= 0 {
		return 0
	}
	return RateOfSurplusValue(roundHalfUp(int64(s)*basisPoints, int64(v)))
}

// ProfitRateAnalysis records a single capital decomposition c + v + s alongside
// the two rates it gives rise to and the gap between them. It is the persisted
// entity behind /v1/profit/rate: the inputs (constant and variable capital), the
// rate of profit p' (which itself carries s and C = c + v), the rate of
// surplus-value s', and the degree to which the former mystifies the latter.
type ProfitRateAnalysis struct {
	ID               ProfitRateID        `json:"id"`
	ConstantCapital  ConstantCapital     `json:"constant_capital"`
	VariableCapital  VariableCapital     `json:"variable_capital"`
	ProfitRate       RateOfProfit        `json:"profit_rate"`
	SurplusValueRate RateOfSurplusValue  `json:"surplus_value_rate"`
	Mystification    MystificationDegree `json:"mystification"`
	CreatedAt        time.Time           `json:"created_at"`
}

// ComputeProfitRateAnalysis derives both rates and the mystification degree from
// a capital decomposition c + v + s. The rate of profit p' = s/C and the rate of
// surplus-value s' = s/v are computed independently; the mystification degree is
// (s' - p')/s', which is 0 when there is no constant capital (the rates coincide)
// and rises towards 10000 bp as constant capital dominates. A zero surplus-value
// rate (v = 0) leaves the mystification at 0: there is no rate to conceal.
func ComputeProfitRateAnalysis(c ConstantCapital, v VariableCapital, s SurplusValue) ProfitRateAnalysis {
	pPrime := ComputeRateOfProfit(s, c, v)
	sPrime := ComputeRateOfSurplusValue(s, v)
	var myst MystificationDegree
	if sPrime > 0 {
		myst = MystificationDegree(roundHalfUp((int64(sPrime)-pPrime.BasisPoints)*basisPoints, int64(sPrime)))
	}
	return ProfitRateAnalysis{
		ConstantCapital:  c,
		VariableCapital:  v,
		ProfitRate:       pPrime,
		SurplusValueRate: sPrime,
		Mystification:    myst,
	}
}
