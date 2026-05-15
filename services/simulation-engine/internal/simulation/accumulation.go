package simulation

import "math"

// AccumulationRate is the fraction of surplus-value reinvested as new capital.
// 0.0 == simple reproduction; 1.0 == full accumulation.
type AccumulationRate float64

// CompositionRatio is the fraction of new capital that is constant capital,
// fixed by the technical conditions of production (e.g. 0.8 for 4c + 1v).
type CompositionRatio float64

// AdditionalCapital is the portion of surplus-value converted into new means
// of production and new labour-power.
type AdditionalCapital struct {
	Constant Pence `json:"constant"`
	Variable Pence `json:"variable"`
}

// Accumulation is one period of extended reproduction: the capital stock
// entering the period, the surplus-value produced, and how it is divided
// between new capital and revenue.
type Accumulation struct {
	Period          int64        `json:"period"`
	CapitalStock    CapitalStock `json:"capital_stock"`
	SurplusProduced Pence        `json:"surplus_produced"`
	NewConstant     Pence        `json:"new_constant"`
	NewVariable     Pence        `json:"new_variable"`
	Revenue         Pence        `json:"revenue"`
}

// SplitSurplus divides surplus-value into new constant capital, new variable
// capital, and revenue consumed by the capitalist, according to the
// accumulation rate and organic composition ratio.
// accumulated is rounded first; variable is the exact complement of constant
// so that constant + variable == accumulated in all cases.
func SplitSurplus(surplus Pence, accumRate float64, compositionRatio float64) AdditionalCapital {
	accumulated := int64(math.Round(float64(surplus) * accumRate))
	constant := Pence(int64(math.Round(float64(accumulated) * compositionRatio)))
	variable := Pence(accumulated - int64(constant))
	return AdditionalCapital{Constant: constant, Variable: variable}
}

// RunExtendedReproduction traces the genealogy of accumulated capital for
// periods steps. Each step's CapitalStock is the AdditionalCapital from the
// prior step, tracing how each layer of accumulated surplus-value in turn
// begets new surplus-value — the "Abraham begat Isaac" sequence from §1.
func RunExtendedReproduction(initial CapitalStock, surplusRate, accumRate, compositionRatio float64, periods int64) []Accumulation {
	if periods <= 0 {
		return []Accumulation{}
	}
	result := make([]Accumulation, periods)
	cap := initial
	for i := int64(0); i < periods; i++ {
		sv := Pence(math.Round(float64(cap.VariableCapital) * surplusRate))
		additional := SplitSurplus(sv, accumRate, compositionRatio)
		revenue := sv - additional.Constant - additional.Variable
		result[i] = Accumulation{
			Period:          i + 1,
			CapitalStock:    cap,
			SurplusProduced: sv,
			NewConstant:     additional.Constant,
			NewVariable:     additional.Variable,
			Revenue:         revenue,
		}
		cap = CapitalStock{
			ConstantCapital: additional.Constant,
			VariableCapital: additional.Variable,
		}
	}
	return result
}
