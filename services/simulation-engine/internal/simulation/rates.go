// Package simulation implements Marx's three formulae for the rate of
// surplus-value from Capital Vol. I, Ch. 18. All functions are pure;
// no persistence is required.
package simulation

import "github.com/theding0x/capital-simulator/services/simulation-engine/internal/labour"

// NecessaryLabour is the paid portion of the working day — the time the
// labourer needs to reproduce the value of their labour-power.
type NecessaryLabour struct {
	Minutes labour.LabourMinutes `json:"minutes"`
}

// SurplusLabour is the unpaid portion — time worked beyond necessary
// labour, appropriated by capital as surplus-value.
type SurplusLabour struct {
	Minutes labour.LabourMinutes `json:"minutes"`
}

// WorkingDay is the total time the labourer is at the capitalist's
// disposal. TotalMinutes equals NecessaryLabour.Minutes + SurplusLabour.Minutes.
type WorkingDay struct {
	TotalMinutes labour.LabourMinutes `json:"total_minutes"`
}

// VariableCapital is the money equivalent of necessary labour — the
// wage advanced by capital. Expressed in pence as a proxy for the
// labour-time ratio; absolute pricing belongs to Ch. 2-3.
type VariableCapital struct {
	Pence int64 `json:"pence"`
}

// SurplusValue is the money equivalent of surplus labour appropriated
// by capital. Expressed in pence as a proxy for the labour-time ratio.
type SurplusValue struct {
	Pence int64 `json:"pence"`
}

// ValueOfProduct is the newly created value in one working day:
// variable capital plus surplus-value (v + s). Constant capital is
// excluded — it merely reappears in the product unchanged.
type ValueOfProduct struct {
	Pence int64 `json:"pence"`
}

// RateScenario is the input to all three rate formulae: the two
// labour-time partitions of a working day.
type RateScenario struct {
	NecessaryLabour NecessaryLabour `json:"necessary_labour"`
	SurplusLabour   SurplusLabour   `json:"surplus_labour"`
}

// RateResult holds the three equivalent expressions of the rate of
// surplus-value computed by ComputeRates.
type RateResult struct {
	FormulaI   float64 `json:"formula_i"`
	FormulaII  float64 `json:"formula_ii"`
	FormulaIII float64 `json:"formula_iii"`
}

// FormulaI returns surplus-labour divided by necessary-labour: s/v.
// Unbounded above 1.0 — a 300% rate of exploitation is possible
// (English agricultural labourer, §).
func FormulaI(s SurplusLabour, v NecessaryLabour) float64 {
	if v.Minutes == 0 {
		return 0
	}
	return float64(s.Minutes) / float64(v.Minutes)
}

// FormulaII returns surplus-labour divided by the total working day:
// s/(s+v). Always strictly less than 1.0, because surplus-labour is
// always less than the full working day.
func FormulaII(s SurplusLabour, wd WorkingDay) float64 {
	if wd.TotalMinutes == 0 {
		return 0
	}
	return float64(s.Minutes) / float64(wd.TotalMinutes)
}

// FormulaIII returns unpaid labour divided by paid labour. Renames
// Formula I's ratio (unpaid = surplus, paid = necessary) but the
// arithmetic is identical.
func FormulaIII(unpaid, paid labour.LabourMinutes) float64 {
	if paid == 0 {
		return 0
	}
	return float64(unpaid) / float64(paid)
}

// ComputeRates applies all three formulae to scenario and returns the
// results. The working day is derived from the two partitions.
func ComputeRates(s RateScenario) RateResult {
	wd := WorkingDay{TotalMinutes: s.NecessaryLabour.Minutes + s.SurplusLabour.Minutes}
	return RateResult{
		FormulaI:   FormulaI(s.SurplusLabour, s.NecessaryLabour),
		FormulaII:  FormulaII(s.SurplusLabour, wd),
		FormulaIII: FormulaIII(s.SurplusLabour.Minutes, s.NecessaryLabour.Minutes),
	}
}
