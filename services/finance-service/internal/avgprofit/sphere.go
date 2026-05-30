package avgprofit

// This file carries Ch. 8's argument: with the same rate of surplus-value in
// every branch but different organic compositions (different v/C splits), the
// individual rate of profit p′ = s/C differs across branches. Living labour —
// and therefore surplus-value, and therefore profit — is set in motion by the
// variable capital alone. A branch with a larger v/C sets more living labour in
// motion per £100 of total capital and shows a higher individual rate of profit.
//
// The equalisation of these individual rates into a general rate of profit is
// deferred to Ch. 9. This chapter only exhibits the divergence: Sphere B
// (100c + 600v, s′ = 100%) yields six times the surplus-value of Sphere A
// (600c + 100v, s′ = 100%) on the same total capital.

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

// labourMinutesPerPound is the value of one unit of variable capital (£1) in
// labour-minutes. Marx's §"£100 ... 6,000 working-hours" — £1 of variable
// capital = one labour-power-week = 60 working-hours = 3600 labour-minutes.
const labourMinutesPerPound = 3600

// LabourMinutes is the canonical value-magnitude unit for this package.
type LabourMinutes int64

// TotalCapital is C = c + v: the whole advance; denominator of the rate of profit.
type TotalCapital LabourMinutes

// ConstantCapital is c: the value of means of production consumed.
type ConstantCapital LabourMinutes

// VariableCapital is v: the value of labour-power purchased; the sole index of
// living labour and therefore the sole source of surplus-value.
type VariableCapital LabourMinutes

// SurplusValue is s: the unpaid labour objectified in the product.
type SurplusValue LabourMinutes

// RateOfSurplusValue is s′ = s/v in basis points (10000 = 100%).
type RateOfSurplusValue int64

// SphereProfitRate is p′ = s/C for an individual branch, in basis points.
// Because C ≥ v, p′ ≤ s′; the gap between them measures the concealment wrought
// by the profit-form.
type SphereProfitRate int64

// LabourPowerIndex is the living labour-minutes set in motion per £100 of total
// capital. It makes visible the sole origin of surplus-value: the branch with
// the higher v/C index sets more living labour in motion and thus generates more
// surplus-value per pound of capital advanced.
type LabourPowerIndex int64

// Sentinel errors for invariant violations.
var (
	ErrNegativeMagnitude    = errors.New("avgprofit: magnitude must be non-negative")
	ErrVariableExceedsTotal = errors.New("avgprofit: variable capital exceeds total capital")
)

// ProductionSphereID identifies a persisted production sphere (branch).
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type ProductionSphereID string

// NewProductionSphereID returns a fresh random ProductionSphereID.
func NewProductionSphereID() ProductionSphereID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("avgprofit: rand: %w", err))
	}
	return ProductionSphereID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty ProductionSphereID.
func (id ProductionSphereID) IsZero() bool { return id == "" }

// ProductionSphere is the persisted entity: a branch of production. Inputs are
// Name, C (total capital), V (variable capital), and SRate (s′ in basis points).
// The derived fields — ConstantCapital, SurplusValue, IndividualProfitRate,
// VariablePercent, LabourPowerIndex — are filled by ComputeProductionSphere.
type ProductionSphere struct {
	ID                   ProductionSphereID `json:"id"`
	Name                 string             `json:"name"`
	C                    TotalCapital       `json:"c"`
	V                    VariableCapital    `json:"v"`
	SRate                RateOfSurplusValue `json:"s_rate"`
	ConstantCapital      ConstantCapital    `json:"constant_capital"`
	SurplusValue         SurplusValue       `json:"surplus_value"`
	IndividualProfitRate SphereProfitRate   `json:"individual_profit_rate"`
	VariablePercent      int64              `json:"variable_percent"`
	LabourPowerIndex     LabourPowerIndex   `json:"labour_power_index"`
	CreatedAt            time.Time          `json:"created_at"`
}

// Validate enforces the invariants on the user-supplied input fields.
func (s ProductionSphere) Validate() error {
	if s.C < 0 || s.V < 0 {
		return ErrNegativeMagnitude
	}
	if s.SRate < 0 {
		return ErrNegativeMagnitude
	}
	if int64(s.V) > int64(s.C) {
		return ErrVariableExceedsTotal
	}
	return nil
}

// surplusValue computes s = s′·v / 10000.
// Integer division is exact for all of Marx's worked fixtures.
func surplusValue(sRate RateOfSurplusValue, v VariableCapital) SurplusValue {
	return SurplusValue(int64(sRate) * int64(v) / basisPoints)
}

// individualProfitRate computes p′ = s/C in basis points, rounded half-up.
// Rounding is required: Sphere A (s=100, C=700) → 1,000,000/700 = 1428.57 → 1429.
// Truncation gives the wrong 1428.
// Implementation: if C <= 0 { 0 } else (s*10000 + C/2) / C.
func individualProfitRate(s SurplusValue, c TotalCapital) SphereProfitRate {
	if c <= 0 {
		return 0
	}
	return SphereProfitRate((int64(s)*basisPoints + int64(c)/2) / int64(c))
}

// variablePercent computes v·100 / C as an integer percentage.
func variablePercent(v VariableCapital, c TotalCapital) int64 {
	if c <= 0 {
		return 0
	}
	return int64(v) * 100 / int64(c)
}

// labourPowerIndex computes (v·labourMinutesPerPound)·100 / C.
// It measures the living labour-minutes set in motion per £100 of total capital.
func labourPowerIndex(v VariableCapital, c TotalCapital) LabourPowerIndex {
	if c <= 0 {
		return 0
	}
	return LabourPowerIndex(int64(v) * labourMinutesPerPound * 100 / int64(c))
}

// ComputeProductionSphere derives all fields from the user-supplied inputs.
// It is a pure constructor: it assigns no ID or timestamp — the store does that
// on persistence.
func ComputeProductionSphere(name string, c TotalCapital, v VariableCapital, sRate RateOfSurplusValue) ProductionSphere {
	s := surplusValue(sRate, v)
	return ProductionSphere{
		Name:                 name,
		C:                    c,
		V:                    v,
		SRate:                sRate,
		ConstantCapital:      ConstantCapital(int64(c) - int64(v)),
		SurplusValue:         s,
		IndividualProfitRate: individualProfitRate(s, c),
		VariablePercent:      variablePercent(v, c),
		LabourPowerIndex:     labourPowerIndex(v, c),
	}
}
