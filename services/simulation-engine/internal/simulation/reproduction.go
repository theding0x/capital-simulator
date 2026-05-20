package simulation

import (
	"crypto/rand"
	"encoding/hex"
	"math"
)

// Pence is the monetary unit for Part VII — The Accumulation of Capital.
type Pence int64

// CapitalStock is the value composition of capital at a point in time.
type CapitalStock struct {
	ConstantCapital Pence `json:"constant_capital"`
	VariableCapital Pence `json:"variable_capital"`
}

// SurplusValueFund records how the surplus-value produced in one cycle is split.
// Under simple reproduction Revenue == Total and Accumulated == 0.
type SurplusValueFund struct {
	Total       Pence `json:"total"`
	Revenue     Pence `json:"revenue"`
	Accumulated Pence `json:"accumulated"`
}

// ReproductionCycle is a snapshot of one circuit of production.
type ReproductionCycle struct {
	Period      int64            `json:"period"`
	Capital     CapitalStock     `json:"capital"`
	SurplusRate float64          `json:"surplus_rate"`
	Fund        SurplusValueFund `json:"fund"`
}

// ScenarioID is a 96-bit hex identifier for a persisted reproduction scenario.
type ScenarioID string

func (s ScenarioID) IsZero() bool { return s == "" }

// NewScenarioID generates a 96-bit hex identifier via crypto/rand.
func NewScenarioID() ScenarioID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return ScenarioID(hex.EncodeToString(b))
}

// RunSimpleReproduction simulates simple reproduction for the given number of
// periods. In each cycle the entire surplus-value is consumed as revenue;
// the capital stock does not grow.
func RunSimpleReproduction(initial CapitalStock, surplusRate float64, periods int64) []ReproductionCycle {
	if periods <= 0 {
		return []ReproductionCycle{}
	}
	cycles := make([]ReproductionCycle, periods)
	cap := initial
	for i := int64(0); i < periods; i++ {
		sv := Pence(math.Round(float64(cap.VariableCapital) * surplusRate))
		fund := SurplusValueFund{
			Total:       sv,
			Revenue:     sv,
			Accumulated: 0,
		}
		cycles[i] = ReproductionCycle{
			Period:      i + 1,
			Capital:     cap,
			SurplusRate: surplusRate,
			Fund:        fund,
		}
	}
	return cycles
}

// RepaymentPeriod returns the number of reproduction periods before the total
// consumed surplus-value equals or exceeds the original capital advanced.
// Integer division; fractional periods round up.
func RepaymentPeriod(capital, annualRevenue Pence) int64 {
	if annualRevenue <= 0 {
		return 0
	}
	q := int64(capital) / int64(annualRevenue)
	if int64(capital)%int64(annualRevenue) != 0 {
		q++
	}
	return q
}
