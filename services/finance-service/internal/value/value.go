// Package value declares the Vol. I value-magnitude vocabulary once for
// finance-service: c, v and s as distinct named types over the canonical
// LabourMinutes unit. Vol. III's profit, avgprofit and revenue packages all
// speak this vocabulary, so a Vol. I surplus-value s is the same type as the
// s that Vol. III divides into profit, interest and rent.
package value

// LabourMinutes is the canonical value-magnitude unit. Marx's worked examples
// are denominated in money (£, shillings), but he fixes the correspondence —
// "the value produced by one labourer during an average social working-day is
// represented by a money sum of 6s." — so the £-figures and the labour-time
// figures are the same magnitude under two names. We keep the labour-time unit.
type LabourMinutes int64

// ConstantCapital (c) is the value transferred from the consumed means of
// production to the product. It reappears in the product's value but creates
// no new value of its own.
type ConstantCapital LabourMinutes

// VariableCapital (v) is the value advanced for labour-power; the sole index of
// living labour and therefore the sole source of surplus-value.
type VariableCapital LabourMinutes

// SurplusValue (s) is the unpaid labour appropriated by the capitalist — the
// excess of the new value created over the variable capital advanced.
type SurplusValue LabourMinutes
