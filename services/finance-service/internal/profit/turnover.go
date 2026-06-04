package profit

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// This file carries Ch. 4's argument: turnover corrects the rate of profit.
//
// Ch. 2–3 measured the rate of profit on a single turnover, p' = s'(v/C). But
// the surplus-value a variable capital yields in a year is not s once — it is s
// once for every time that capital turns over. With S = sn and the annual rate
// of surplus-value S' = s'n (Vol. II, Ch. XVI), the rate of profit corrected for
// turnover becomes
//
//	p' = s' · n · (v / C),
//
// the annual rate of profit. Two capitals of equal composition and equal s'
// then have rates of profit "related inversely as their periods of turnover":
// the one that turns over twice a year earns twice the annual rate of the one
// that turns over once. v is the variable capital *advanced* (the v that sits in
// C = c + v); the capitalist's books show only the year's whole wage bill vn,
// never v itself — a distinction Marx stresses and that this package keeps.

// TurnoverCount is n: the number of times the variable capital turns over in a
// year. It is a plain count of whole turnovers; fractional turnovers (Marx's
// Manchester 8½) are an approximation outside the integer model.
type TurnoverCount int64

// AnnualSurplusValueRate is S' = s' · n in basis points (100% = 10000 bp): the
// surplus-value appropriated across a whole year relative to the variable
// capital advanced, as against s', the rate over a single turnover. It is the
// rate the annual rate of profit implies once turnover is taken into account.
type AnnualSurplusValueRate int64

// VariableCapitalPerTurnover is the variable capital advanced in a single
// turnover period — the v that figures in C = c + v and in the rate of profit.
// It is what the formula p' = s'·n·(v/C) requires, and what the capitalist
// (who sees only the year's total wage bill vn = v·n) does not directly know.
type VariableCapitalPerTurnover LabourMinutes

// AnnualProfitRate is the central magnitude of this chapter: the rate of profit
// corrected for turnover, p' = s' · n · (v / C). It carries the four factors so
// the correction stays legible beside the single-turnover rate s'·(v/C): the
// rate of surplus-value, the number of turnovers, the variable capital advanced,
// and the total capital. BasisPoints is the rate itself (100% = 10000 bp).
type AnnualProfitRate struct {
	SurplusValueRate RateOfSurplusValue         `json:"surplus_value_rate"`
	Turnovers        TurnoverCount              `json:"turnovers"`
	VariableCapital  VariableCapitalPerTurnover `json:"variable_capital"`
	TotalCapital     TotalCapital               `json:"total_capital"`
	BasisPoints      int64                      `json:"basis_points"`
}

// ComputeAnnualProfitRate forms p' = s' · n · (v / C) in basis points. With no
// capital advanced (C = 0) the rate is zero. The multiplication is carried out
// before the division, and the quotient is rounded half-up (the Vol. III house
// convention), so the integer result matches Marx's worked rates exactly
// (Capital A: 100% · 2 · 20/100 = 40%).
func ComputeAnnualProfitRate(sRate RateOfSurplusValue, n TurnoverCount, v VariableCapitalPerTurnover, c TotalCapital) AnnualProfitRate {
	var bp int64
	if c > 0 {
		bp = roundHalfUp(int64(sRate)*int64(n)*int64(v), int64(c))
	}
	return AnnualProfitRate{
		SurplusValueRate: sRate,
		Turnovers:        n,
		VariableCapital:  v,
		TotalCapital:     c,
		BasisPoints:      bp,
	}
}

// TurnoverAnalysisID identifies a stored turnover analysis. 96-bit hex from
// crypto/rand, like every other domain ID in the simulator.
type TurnoverAnalysisID string

// NewTurnoverAnalysisID returns a fresh random TurnoverAnalysisID.
func NewTurnoverAnalysisID() TurnoverAnalysisID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("profit: rand: %w", err))
	}
	return TurnoverAnalysisID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty TurnoverAnalysisID.
func (id TurnoverAnalysisID) IsZero() bool { return id == "" }

// TurnoverAnalysis records one capital and the effect of its turnover on the
// rate of profit. It stores the inputs — total capital C, variable capital
// advanced V (the v in C), rate of surplus-value SRate (s', basis points) and
// the number of turnovers N — and the magnitudes they give rise to: the annual
// rate of profit p' = s'·n·(v/C), the single-turnover rate s'·(v/C) it would be
// without turnover, the annual rate of surplus-value S' = s'·n, and the year's
// whole wage bill vn (AnnualWages) that the capitalist sees in place of v.
type TurnoverAnalysis struct {
	ID                       TurnoverAnalysisID         `json:"id"`
	C                        TotalCapital               `json:"c"`
	V                        VariableCapitalPerTurnover `json:"v"`
	SRate                    RateOfSurplusValue         `json:"s_rate"`
	N                        TurnoverCount              `json:"n"`
	AnnualProfitRate         AnnualProfitRate           `json:"annual_profit_rate"`
	SingleTurnoverProfitRate int64                      `json:"single_turnover_profit_rate"`
	AnnualSurplusValueRate   AnnualSurplusValueRate     `json:"annual_surplus_value_rate"`
	AnnualWages              LabourMinutes              `json:"annual_wages"`
	CreatedAt                time.Time                  `json:"created_at"`
}

// Validate enforces the magnitude invariants: non-negative inputs, and the
// variable capital advanced cannot exceed the total capital it is a part of.
func (a TurnoverAnalysis) Validate() error {
	if a.C < 0 || a.V < 0 || a.SRate < 0 || a.N < 0 {
		return ErrNegativeMagnitude
	}
	if int64(a.V) > int64(a.C) {
		return ErrVariableExceedsTotal
	}
	return nil
}

// ComputeTurnoverAnalysis derives the annual rate of profit and its companions
// from a capital (total C, variable advanced v), a rate of surplus-value s' and
// a turnover count n. It is a pure function: it assigns no ID or timestamp — the
// store does that on persistence.
//
// The single-turnover rate is p' as it would stand without the turnover factor,
// s'·(v/C); the annual rate multiplies it by n. The annual rate of surplus-value
// is S' = s'·n, and the annual wage bill is vn — the figure the capitalist's
// books actually record, as distinct from the variable capital v it conceals.
func ComputeTurnoverAnalysis(c TotalCapital, v VariableCapitalPerTurnover, sRate RateOfSurplusValue, n TurnoverCount) TurnoverAnalysis {
	annual := ComputeAnnualProfitRate(sRate, n, v, c)
	var single int64
	if c > 0 {
		single = roundHalfUp(int64(sRate)*int64(v), int64(c))
	}
	return TurnoverAnalysis{
		C:                        c,
		V:                        v,
		SRate:                    sRate,
		N:                        n,
		AnnualProfitRate:         annual,
		SingleTurnoverProfitRate: single,
		AnnualSurplusValueRate:   AnnualSurplusValueRate(int64(sRate) * int64(n)),
		AnnualWages:              LabourMinutes(int64(v) * int64(n)),
	}
}
