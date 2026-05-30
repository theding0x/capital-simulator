package avgprofit

// This file carries the two composition concepts that underpin Ch. 8's argument.
//
// TechnicalComposition is the physical ratio of means of production to workers —
// the real, material basis of the organic composition. A more mechanised branch
// has a higher technical composition (more Mp per worker).
//
// OrganicComposition is the value expression of that physical ratio: c/(c+v), the
// share of total capital that is constant. A branch with a higher organic
// composition has more constant capital per unit of variable capital and therefore
// sets less living labour in motion per pound of total capital — and therefore
// generates less surplus-value and shows a lower individual rate of profit.

// CompositionKind names the relation of a branch's organic composition to the
// social average. HIGHER organic composition means MORE constant capital, i.e. a
// LOWER variable percent — the branch sets less living labour in motion per £100
// of total capital and shows a lower individual rate of profit.
type CompositionKind string

const (
	// KindHigher: more constant capital than the social average → lower p′.
	KindHigher CompositionKind = "higher"
	// KindAverage: organic composition equal to the social average.
	KindAverage CompositionKind = "average"
	// KindLower: less constant capital than the social average → higher p′.
	KindLower CompositionKind = "lower"
)

// TechnicalComposition is the physical ratio of means of production to workers —
// the real foundation of the organic composition. It changes whenever a branch
// adopts new machinery or methods; the organic composition follows in value terms.
type TechnicalComposition struct {
	MeansOfProduction int64 `json:"means_of_production"` // physical units of Mp
	Workers           int64 `json:"workers"`             // number of labourers
}

// MpPerWorker returns the number of means-of-production units per worker.
// Returns 0 when there are no workers (avoiding division by zero).
func (t TechnicalComposition) MpPerWorker() int64 {
	if t.Workers <= 0 {
		return 0
	}
	return t.MeansOfProduction / t.Workers
}

// OrganicComposition is the value-ratio c/(c+v), expressed as the absolute
// magnitudes of constant and variable capital. The percentage methods derive
// the share of each component in total capital C = c + v.
type OrganicComposition struct {
	Constant ConstantCapital `json:"constant"`
	Variable VariableCapital `json:"variable"`
}

// Total returns C = c + v.
func (o OrganicComposition) Total() TotalCapital {
	return TotalCapital(int64(o.Constant) + int64(o.Variable))
}

// ConstantPercent returns c·100/(c+v) as an integer percentage.
// Returns 0 when total capital is zero.
func (o OrganicComposition) ConstantPercent() int64 {
	total := o.Total()
	if total <= 0 {
		return 0
	}
	return int64(o.Constant) * 100 / int64(total)
}

// VariablePercent returns v·100/(c+v) as an integer percentage.
// Returns 0 when total capital is zero. Note: VariablePercent + ConstantPercent
// may not sum to exactly 100 due to integer truncation.
func (o OrganicComposition) VariablePercent() int64 {
	total := o.Total()
	if total <= 0 {
		return 0
	}
	return int64(o.Variable) * 100 / int64(total)
}

// ClassifyComposition compares a sphere's variable-percent to the social average:
//   - below average → KindHigher (more constant capital; less living labour per £100)
//   - above average → KindLower  (less constant capital; more living labour per £100)
//   - equal         → KindAverage
func ClassifyComposition(sphereVariablePercent, averageVariablePercent int64) CompositionKind {
	switch {
	case sphereVariablePercent < averageVariablePercent:
		return KindHigher
	case sphereVariablePercent > averageVariablePercent:
		return KindLower
	default:
		return KindAverage
	}
}
