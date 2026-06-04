package profit

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// This file carries Ch. 3's argument: the rate of profit is the rate of
// surplus-value seen through the prism of the total capital. Marx's master
// equation,
//
//	p' = s' · (v / C) = s' · v / (c + v),
//
// resolves the rate of profit into two independent factors — the rate of
// surplus-value s' and the value-composition of capital v/C. The whole chapter
// is the systematic variation of those factors: hold s' fixed and move v/C
// (§I); hold v/C fixed and move s' (§II); move both (§II.3, which reduces to
// the simpler cases). The types here let the simulator run any such variation
// and read off how p' shifts.

// ProfitRateFormula is one capital decomposition expressed as the variables of
// the rate of profit: the constant capital C (= c), the variable capital V
// (= v), and the rate of surplus-value SRate (= s').
//
// Marx writes the rate of profit as p' = s'(v/C), where C in that expression is
// the *total* capital c + v. To keep a full decomposition recoverable, this
// struct stores the constant component C and the variable component V
// separately; the denominator of the formula is therefore C + V, exposed as
// TotalCapital(). SRate is in basis points (100% = 10000 bp).
type ProfitRateFormula struct {
	C     ConstantCapital    `json:"c"`
	V     VariableCapital    `json:"v"`
	SRate RateOfSurplusValue `json:"s_rate"`
}

// Validate enforces non-negative magnitudes.
func (f ProfitRateFormula) Validate() error {
	if f.C < 0 || f.V < 0 || f.SRate < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// TotalCapital returns C = c + v, the denominator of the rate of profit.
func (f ProfitRateFormula) TotalCapital() TotalCapital {
	return TotalCapital(int64(f.C) + int64(f.V))
}

// ProfitRate returns p' = s' · v / C = s' · v / (c + v) in basis points. With
// no capital advanced (c = v = 0) the rate is zero; so too when v = 0, since
// variable capital is the sole source of the surplus-value the rate measures.
// The quotient is rounded half-up (the Vol. III house convention shared with
// avgprofit and tendency), so the same s'·v/C lands on identical basis points
// wherever it is computed.
func (f ProfitRateFormula) ProfitRate() int64 {
	total := int64(f.C) + int64(f.V)
	if total == 0 {
		return 0
	}
	return roundHalfUp(int64(f.SRate)*int64(f.V), total)
}

// CompositionRatio returns the value-composition v/C of this decomposition.
func (f ProfitRateFormula) CompositionRatio() CompositionRatio {
	return ComputeCompositionRatio(f.C, f.V)
}

// CompositionRatio is the value-composition of capital read from the
// profit-rate side: the variable capital's share of the total, v/C, in basis
// points. It is the second factor of p' = s' · (v/C); at a given s', a higher
// v/C (a lower organic composition) yields a higher rate of profit.
type CompositionRatio struct {
	Variable    VariableCapital `json:"variable"`
	Total       TotalCapital    `json:"total"`
	BasisPoints int64           `json:"basis_points"`
}

// ComputeCompositionRatio forms v/C in basis points from c and v. A zero total
// capital yields a zero ratio.
func ComputeCompositionRatio(c ConstantCapital, v VariableCapital) CompositionRatio {
	total := TotalCapital(int64(c) + int64(v))
	var bp int64
	if total > 0 {
		bp = int64(v) * basisPoints / int64(total)
	}
	return CompositionRatio{Variable: v, Total: total, BasisPoints: bp}
}

// VariationCase names which factors of p' = s'(v/C) move between two
// decompositions. Marx structures the chapter around these cases: first s'
// constant while v/C varies (§I), then v/C constant while s' varies (§II), then
// both (§II.3, reducible to the others).
type VariationCase string

const (
	// CaseSConstantVCVariable: the rate of surplus-value is unchanged and the
	// value-composition v/C moves (§I). p' then moves with v/C.
	CaseSConstantVCVariable VariationCase = "s_constant_vc_variable"
	// CaseSVariableVCConstant: the value-composition is unchanged and the rate
	// of surplus-value moves (§II). p' then moves in proportion to s'.
	CaseSVariableVCConstant VariationCase = "s_variable_vc_constant"
	// CaseBothVariable: both factors move (§II.3); the resulting movement of p'
	// reduces to a combination of the two simpler cases.
	CaseBothVariable VariationCase = "both_variable"
)

// DeriveVariationCase classifies the move from initial to changed by which
// factor — the rate of surplus-value or the value-composition v/C — varies.
func DeriveVariationCase(initial, changed ProfitRateFormula) VariationCase {
	sameSRate := initial.SRate == changed.SRate
	sameComposition := initial.CompositionRatio().BasisPoints == changed.CompositionRatio().BasisPoints
	switch {
	case sameSRate && !sameComposition:
		return CaseSConstantVCVariable
	case !sameSRate && sameComposition:
		return CaseSVariableVCConstant
	default:
		return CaseBothVariable
	}
}

// ProfitRateComparison sets two capitals side by side and reads off how their
// rates of profit stand to each other. Marx's §I theorem: two capitals with the
// same rate of surplus-value have rates of profit in the same proportion as
// their value-compositions, p'₁ : p'₂ = v₁/C₁ : v₂/C₂. EqualSurplusRate flags
// when that theorem applies.
type ProfitRateComparison struct {
	First             ProfitRateFormula `json:"first"`
	Second            ProfitRateFormula `json:"second"`
	FirstProfitRate   int64             `json:"first_profit_rate"`
	SecondProfitRate  int64             `json:"second_profit_rate"`
	FirstComposition  CompositionRatio  `json:"first_composition"`
	SecondComposition CompositionRatio  `json:"second_composition"`
	EqualSurplusRate  bool              `json:"equal_surplus_rate"`
}

// CompareProfitRates computes both rates of profit and value-compositions and
// flags whether the two capitals share a rate of surplus-value (in which case
// their rates of profit stand as their value-compositions).
func CompareProfitRates(first, second ProfitRateFormula) ProfitRateComparison {
	return ProfitRateComparison{
		First:             first,
		Second:            second,
		FirstProfitRate:   first.ProfitRate(),
		SecondProfitRate:  second.ProfitRate(),
		FirstComposition:  first.CompositionRatio(),
		SecondComposition: second.CompositionRatio(),
		EqualSurplusRate:  first.SRate == second.SRate,
	}
}

// VariationAnalysisID identifies a stored variation analysis. 96-bit hex from
// crypto/rand, like every other domain ID in the simulator.
type VariationAnalysisID string

// NewVariationAnalysisID returns a fresh random VariationAnalysisID.
func NewVariationAnalysisID() VariationAnalysisID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("profit: rand: %w", err))
	}
	return VariationAnalysisID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty VariationAnalysisID.
func (id VariationAnalysisID) IsZero() bool { return id == "" }

// VariationAnalysis records the movement of the rate of profit between an
// initial and a changed decomposition. It stores both decompositions, the case
// that classifies the move, the old and new rates of profit and
// value-compositions (basis points), and whether p' moved in the same
// proportion as s' — the signature of the §II case where v/C is held constant.
type VariationAnalysis struct {
	ID                 VariationAnalysisID `json:"id"`
	Case               VariationCase       `json:"case"`
	Initial            ProfitRateFormula   `json:"initial"`
	Changed            ProfitRateFormula   `json:"changed"`
	OldProfitRate      int64               `json:"old_profit_rate"`
	NewProfitRate      int64               `json:"new_profit_rate"`
	OldCompositionBP   int64               `json:"old_composition_bp"`
	NewCompositionBP   int64               `json:"new_composition_bp"`
	ProportionalChange bool                `json:"proportional_change"`
	CreatedAt          time.Time           `json:"created_at"`
}

// ComputeVariationAnalysis derives the case, both rates of profit, both
// value-compositions, and the proportional-change flag from an initial and a
// changed decomposition. It is a pure function: it assigns no ID or timestamp —
// the store does that on persistence.
//
// ProportionalChange is true when the rate of profit moved in exact proportion
// to the rate of surplus-value, i.e. p'old : p'new = s'old : s'new. By Marx's
// §II this holds precisely when the value-composition v/C is unchanged; it is
// the case in which "p' increases or decreases in the same proportion as s'".
func ComputeVariationAnalysis(initial, changed ProfitRateFormula) VariationAnalysis {
	oldP := initial.ProfitRate()
	newP := changed.ProfitRate()
	proportional := int64(initial.SRate) != 0 &&
		oldP*int64(changed.SRate) == newP*int64(initial.SRate)
	return VariationAnalysis{
		Case:               DeriveVariationCase(initial, changed),
		Initial:            initial,
		Changed:            changed,
		OldProfitRate:      oldP,
		NewProfitRate:      newP,
		OldCompositionBP:   initial.CompositionRatio().BasisPoints,
		NewCompositionBP:   changed.CompositionRatio().BasisPoints,
		ProportionalChange: proportional,
	}
}
