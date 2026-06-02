// Package valorisation — Vol. II Ch. 16: The Turnover of Variable Capital.
//
// Marx derives the annual rate of surplus-value: where the per-circuit rate
// is s/v, the annual rate is n × s/v, where n is the number of turnovers per
// year of the variable capital. Two capitals with identical s/v but different
// n have dramatically different annual rates.
//
// Canonical contrast: Capitalist A (5-week turnover, n=10, £500 advance) vs
// Capitalist B (50-week turnover, n=1, £5,000 advance) — both produce £5,000
// of annual surplus, but A's annual rate is 1000% vs B's 100%.
//
// All monetary values are in pence (1 GBP = 240 pence). All rates and
// multipliers are expressed as basis points (100% = 10,000 bp) to avoid
// floating-point arithmetic.
package valorisation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// --- Sentinel errors ---------------------------------------------------------

var (
	// ErrNegativeSurplus is returned when surplus_pence is negative.
	ErrNegativeSurplus = errors.New("valorisation: surplus_pence cannot be negative")

	// ErrZeroVariable is returned when variable_pence is zero or negative.
	ErrZeroVariable = errors.New("valorisation: variable_pence must be positive")

	// ErrZeroTurnoverPeriod is returned when turnover_period_minutes is zero or negative.
	ErrZeroTurnoverPeriod = errors.New("valorisation: turnover_period_minutes must be positive")

	// ErrZeroAdvancePence is returned when the advance pence is zero or negative.
	ErrZeroAdvancePence = errors.New("valorisation: advance pence must be positive")
)

// --- ID generation -----------------------------------------------------------

func newHexID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}

// --- Pence -------------------------------------------------------------------

// Pence is the canonical monetary unit (1 GBP = 240 pence).
type Pence int64

// --- ID types ----------------------------------------------------------------

// VariableCapitalAdvanceID identifies a VariableCapitalAdvance.
type VariableCapitalAdvanceID string

// NewVariableCapitalAdvanceID returns a 96-bit random hex identifier.
func NewVariableCapitalAdvanceID() VariableCapitalAdvanceID {
	return VariableCapitalAdvanceID(newHexID())
}

// PerCircuitSurplusRateID identifies a PerCircuitSurplusRate record.
type PerCircuitSurplusRateID string

// NewPerCircuitSurplusRateID returns a 96-bit random hex identifier.
func NewPerCircuitSurplusRateID() PerCircuitSurplusRateID {
	return PerCircuitSurplusRateID(newHexID())
}

// AnnualSurplusRateID identifies an AnnualSurplusRate snapshot.
type AnnualSurplusRateID string

// NewAnnualSurplusRateID returns a 96-bit random hex identifier.
func NewAnnualSurplusRateID() AnnualSurplusRateID {
	return AnnualSurplusRateID(newHexID())
}

// AnnualSurplusRateContrastID identifies an AnnualSurplusRateContrast.
type AnnualSurplusRateContrastID string

// NewAnnualSurplusRateContrastID returns a 96-bit random hex identifier.
func NewAnnualSurplusRateContrastID() AnnualSurplusRateContrastID {
	return AnnualSurplusRateContrastID(newHexID())
}

// WorkingPeriodAdvanceID identifies a WorkingPeriodAdvance link.
type WorkingPeriodAdvanceID string

// NewWorkingPeriodAdvanceID returns a 96-bit random hex identifier.
func NewWorkingPeriodAdvanceID() WorkingPeriodAdvanceID {
	return WorkingPeriodAdvanceID(newHexID())
}

// --- Domain types ------------------------------------------------------------

// VariableCapitalAdvance is the variable capital (v) advanced per turnover
// period. Its turnover is the specific turnover of the wage-labour component
// of capital.
//
// TurnoversPerYearBasisPoints encodes n as 10,000 × n (so n=10 → 100,000 bp).
// This allows fractional turnovers to be represented exactly without float64.
type VariableCapitalAdvance struct {
	ID                          VariableCapitalAdvanceID `json:"id"`
	IndustrialCapitalID         string                   `json:"industrial_capital_id"`
	Pence                       Pence                    `json:"pence"`
	TurnoverPeriodMinutes       int64                    `json:"turnover_period_minutes"`
	TurnoversPerYearBasisPoints int64                    `json:"turnovers_per_year_basis_points"`
}

// Validate checks that the advance has a positive monetary amount and a
// positive turnover period.
func (a VariableCapitalAdvance) Validate() error {
	if a.Pence <= 0 {
		return ErrZeroAdvancePence
	}
	if a.TurnoverPeriodMinutes <= 0 {
		return ErrZeroTurnoverPeriod
	}
	return nil
}

// PerCircuitSurplusRate captures the per-circuit surplus rate s/v as basis
// points. The per-circuit rate is determined by the labour-process (Vol. I),
// not by the number of turnovers — a faster-turnover capital does not have a
// higher per-circuit rate.
//
// BasisPoints = 10,000 × SurplusPence / VariablePence (integer division).
type PerCircuitSurplusRate struct {
	ID                       PerCircuitSurplusRateID  `json:"id"`
	VariableCapitalAdvanceID VariableCapitalAdvanceID `json:"variable_capital_advance_id"`
	SurplusPence             Pence                    `json:"surplus_pence"`
	VariablePence            Pence                    `json:"variable_pence"`
	BasisPoints              int64                    `json:"basis_points"`
}

// AnnualSurplusRate is the product n × (s/v) expressed in basis points.
// "The annual rate of surplus-value is … equal to the rate of surplus-value
// produced in one period of turnover, multiplied by the number of turnovers
// of the variable capital." — Marx.
//
// AnnualBasisPoints = NBasisPoints × PerCircuitBasisPoints / 10,000.
type AnnualSurplusRate struct {
	ID                       AnnualSurplusRateID      `json:"id"`
	VariableCapitalAdvanceID VariableCapitalAdvanceID `json:"variable_capital_advance_id"`
	Period                   string                   `json:"period"`
	NBasisPoints             int64                    `json:"n_basis_points"`
	PerCircuitBasisPoints    int64                    `json:"per_circuit_basis_points"`
	AnnualBasisPoints        int64                    `json:"annual_basis_points"`
}

// AnnualSurplusRateContrast compares two capitals' annual rates.
// "Capital A with £500 variable capital and 10 turnovers achieves an annual
// rate of 1000%, whereas Capital B with £5,000 and 1 turnover achieves only
// 100%, despite each paying out £5,000 in wages over the year." — Marx.
//
// RatioBasisPoints = 10,000 × CapitalAAnnualBasisPoints / CapitalBAnnualBasisPoints.
type AnnualSurplusRateContrast struct {
	ID                        AnnualSurplusRateContrastID `json:"id"`
	CapitalAID                VariableCapitalAdvanceID    `json:"capital_a_id"`
	CapitalBID                VariableCapitalAdvanceID    `json:"capital_b_id"`
	CapitalAAnnualBasisPoints int64                       `json:"capital_a_annual_basis_points"`
	CapitalBAnnualBasisPoints int64                       `json:"capital_b_annual_basis_points"`
	RatioBasisPoints          int64                       `json:"ratio_basis_points"`
}

// WorkingPeriodAdvance links a VariableCapitalAdvance to a Ch. 12 working
// period, recording the advance multiplier for connected-mode production.
// For discrete-mode production the multiplier is always 1.
type WorkingPeriodAdvance struct {
	ID                       WorkingPeriodAdvanceID   `json:"id"`
	VariableCapitalAdvanceID VariableCapitalAdvanceID `json:"variable_capital_advance_id"`
	WorkingPeriodID          string                   `json:"working_period_id"`
	AdvanceMultiplier        int64                    `json:"advance_multiplier"`
}

// VariableCapitalReproduction captures the reproduction of variable capital
// over a full year.
//
// Invariant: TotalVariablePaidAnnualPence = Pence × TurnoversPerYearBasisPoints / 10,000.
// Invariant (at 100% per-circuit rate): TotalSurplusAnnualPence == TotalVariablePaidAnnualPence.
type VariableCapitalReproduction struct {
	VariableCapitalAdvanceID     VariableCapitalAdvanceID `json:"variable_capital_advance_id"`
	YearReferenceMinutes         int64                    `json:"year_reference_minutes"`
	TotalVariablePaidAnnualPence Pence                    `json:"total_variable_paid_annual_pence"`
	TotalSurplusAnnualPence      Pence                    `json:"total_surplus_annual_pence"`
}

// --- Pure functions ----------------------------------------------------------

// ComputePerCircuitRate derives the per-circuit surplus rate for the given
// surplus and variable-capital magnitudes.
//
// Returns ErrNegativeSurplus if surplusPence < 0, ErrZeroVariable if
// variablePence ≤ 0.
func ComputePerCircuitRate(
	advanceID VariableCapitalAdvanceID,
	surplusPence, variablePence Pence,
) (PerCircuitSurplusRate, error) {
	if surplusPence < 0 {
		return PerCircuitSurplusRate{}, ErrNegativeSurplus
	}
	if variablePence <= 0 {
		return PerCircuitSurplusRate{}, ErrZeroVariable
	}
	bp := 10_000 * int64(surplusPence) / int64(variablePence)
	return PerCircuitSurplusRate{
		ID:                       NewPerCircuitSurplusRateID(),
		VariableCapitalAdvanceID: advanceID,
		SurplusPence:             surplusPence,
		VariablePence:            variablePence,
		BasisPoints:              bp,
	}, nil
}

// ComputeAnnualSurplusRate computes the annual rate of surplus-value.
//
// AnnualBasisPoints = nBasisPoints × perCircuitBasisPoints / 10,000.
// "Doubling n and s/v simultaneously quadruples the annual rate" —
// this multiplicative structure is why faster-turnover capital is more
// powerful despite identical per-circuit exploitation.
func ComputeAnnualSurplusRate(
	advanceID VariableCapitalAdvanceID,
	nBasisPoints, perCircuitBasisPoints int64,
	period string,
) AnnualSurplusRate {
	var annual int64
	if nBasisPoints > 0 && perCircuitBasisPoints > 0 {
		annual = nBasisPoints * perCircuitBasisPoints / 10_000
	}
	return AnnualSurplusRate{
		ID:                       NewAnnualSurplusRateID(),
		VariableCapitalAdvanceID: advanceID,
		Period:                   period,
		NBasisPoints:             nBasisPoints,
		PerCircuitBasisPoints:    perCircuitBasisPoints,
		AnnualBasisPoints:        annual,
	}
}

// ComputeAnnualSurplusRateContrast builds the canonical A-vs-B comparison.
//
// RatioBasisPoints = 10,000 × CapitalAAnnualBasisPoints / CapitalBAnnualBasisPoints.
// Returns zero ratio when CapitalBAnnualBasisPoints == 0.
func ComputeAnnualSurplusRateContrast(
	aID, bID VariableCapitalAdvanceID,
	aAnnualBP, bAnnualBP int64,
) AnnualSurplusRateContrast {
	var ratio int64
	if bAnnualBP > 0 {
		ratio = 10_000 * aAnnualBP / bAnnualBP
	}
	return AnnualSurplusRateContrast{
		ID:                        NewAnnualSurplusRateContrastID(),
		CapitalAID:                aID,
		CapitalBID:                bID,
		CapitalAAnnualBasisPoints: aAnnualBP,
		CapitalBAnnualBasisPoints: bAnnualBP,
		RatioBasisPoints:          ratio,
	}
}

// PriceAdjustedAnnualRate records how a Ch. 15 value revolution on output
// shifts the Ch. 16 annual rate of surplus-value (issue #222). A price change
// of PriceBasisPointsChange on the realised product re-prices the value added
// (v + s); because the advance is fixed, the whole swing accrues to surplus
// measured against the original advance, which is then re-annualised by n.
type PriceAdjustedAnnualRate struct {
	VariablePence                 Pence `json:"variable_pence"`
	BaseSurplusPence              Pence `json:"base_surplus_pence"`
	PriceBasisPointsChange        int64 `json:"price_basis_points_change"`
	AdjustedSurplusPence          Pence `json:"adjusted_surplus_pence"`
	NBasisPoints                  int64 `json:"n_basis_points"`
	BasePerCircuitBasisPoints     int64 `json:"base_per_circuit_basis_points"`
	AdjustedPerCircuitBasisPoints int64 `json:"adjusted_per_circuit_basis_points"`
	BaseAnnualBasisPoints         int64 `json:"base_annual_basis_points"`
	AdjustedAnnualBasisPoints     int64 `json:"adjusted_annual_basis_points"`
}

// ComputeAdjustedAnnualRate applies a value-revolution price change (priceBP, in
// basis points, signed) to the per-circuit surplus and re-annualises by n.
// AdjustedSurplus = baseSurplus + (v + baseSurplus) × priceBP / 10000; the
// per-circuit and annual rates are recomputed from it. A non-positive variable
// capital yields zero rates (issue #222).
func ComputeAdjustedAnnualRate(variablePence, baseSurplusPence Pence, nBasisPoints, priceBasisPointsChange int64) PriceAdjustedAnnualRate {
	v := int64(variablePence)
	s := int64(baseSurplusPence)
	adjustedSurplus := s + (v+s)*priceBasisPointsChange/10_000

	var basePerCircuit, adjPerCircuit int64
	if v > 0 {
		basePerCircuit = 10_000 * s / v
		adjPerCircuit = 10_000 * adjustedSurplus / v
	}
	return PriceAdjustedAnnualRate{
		VariablePence:                 variablePence,
		BaseSurplusPence:              baseSurplusPence,
		PriceBasisPointsChange:        priceBasisPointsChange,
		AdjustedSurplusPence:          Pence(adjustedSurplus),
		NBasisPoints:                  nBasisPoints,
		BasePerCircuitBasisPoints:     basePerCircuit,
		AdjustedPerCircuitBasisPoints: adjPerCircuit,
		BaseAnnualBasisPoints:         nBasisPoints * basePerCircuit / 10_000,
		AdjustedAnnualBasisPoints:     nBasisPoints * adjPerCircuit / 10_000,
	}
}

// ComputeVariableCapitalReproduction derives the annual reproduction magnitudes.
//
// TotalVariablePaidAnnualPence = adv.Pence × adv.TurnoversPerYearBasisPoints / 10,000.
// TotalSurplusAnnualPence = TotalVariablePaid × perCircuitBP / 10,000.
func ComputeVariableCapitalReproduction(
	adv VariableCapitalAdvance,
	perCircuitBP int64,
	yearMinutes int64,
) VariableCapitalReproduction {
	totalVariable := Pence(int64(adv.Pence) * adv.TurnoversPerYearBasisPoints / 10_000)
	totalSurplus := Pence(int64(totalVariable) * perCircuitBP / 10_000)
	return VariableCapitalReproduction{
		VariableCapitalAdvanceID:     adv.ID,
		YearReferenceMinutes:         yearMinutes,
		TotalVariablePaidAnnualPence: totalVariable,
		TotalSurplusAnnualPence:      totalSurplus,
	}
}
