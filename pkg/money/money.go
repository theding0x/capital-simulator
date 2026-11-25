// Package money is the single money scale and the value<->money bridge for every
// service (ADR 0001, docs/adr/0001-value-money-composition.md; umbrella #369,
// #449). It owns two things and nothing else:
//
//   - one canonical £→pence scale: pre-decimal £sd, as Marx writes, so that
//     "3s 6d" and "6 shillings for a 12-hour day" are exact integers;
//   - the MELT, the monetary expression of labour-time: how many labour-minutes
//     one unit of money stands for.
//
// Keep it free of domain logic. Packages that still carry a legacy per-package
// scale (100, 1000) convert at their boundary with Rescale; existing fixtures and
// migrations are not rewritten.
package money

// Pre-decimal £sd: £1 = 20s = 240d.
const (
	ShillingsPerPound = 20
	PencePerShilling  = 12
	PencePerPound     = ShillingsPerPound * PencePerShilling // 240
)

// The MELT. Marx's Vol. I Ch. 17 example takes a 12-hour working day to create
// 6 shillings of value, so 1s = 2 hours = 120 labour-minutes, 1d = 10
// labour-minutes and £1 = 40 hours = 2,400 labour-minutes.
const (
	LabourMinutesPerPence = 10
	LabourMinutesPerPound = LabourMinutesPerPence * PencePerPound // 2400
)

// Pence is the canonical money-form magnitude: 1/240 of a pound.
type Pence int64

// FromLSD builds Pence from pounds, shillings and pence, e.g. FromLSD(0, 3, 6)
// for Marx's "3s. 6d.".
func FromLSD(pounds, shillings, pence int64) Pence {
	return Pence(pounds*PencePerPound + shillings*PencePerShilling + pence)
}

// ToLabourMinutes expresses p as the labour-time it stands for at the MELT.
// It is exact: one penny is a whole number of labour-minutes.
func (p Pence) ToLabourMinutes() int64 {
	return int64(p) * LabourMinutesPerPence
}

// FromLabourMinutes expresses a labour-time magnitude as money at the MELT,
// rounded half away from zero to the nearest penny.
func FromLabourMinutes(minutes int64) Pence {
	return Pence(roundDiv(minutes, LabourMinutesPerPence))
}

// Rescale converts an amount recorded at fromPerPound units to the pound
// (e.g. 100 for decimal pence, 1000 for the reproduction schemes) to canonical
// Pence, rounded half away from zero. It is the package-boundary adapter for
// legacy scales; fromPerPound must be positive.
func Rescale(amount, fromPerPound int64) Pence {
	return Pence(roundDiv(amount*PencePerPound, fromPerPound))
}

// roundDiv returns num/den rounded half away from zero. den must be positive.
func roundDiv(num, den int64) int64 {
	if num < 0 {
		return -roundDiv(-num, den)
	}
	q, r := num/den, num%den
	if 2*r >= den {
		q++
	}
	return q
}
