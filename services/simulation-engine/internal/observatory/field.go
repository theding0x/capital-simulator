// Package observatory runs the Atlas Observatory as ephemeral, per-session,
// in-memory simulation runs. A Manager seeds each browser session's Run from an
// immutable template (loaded once from the store at boot) and advances it on
// poll; nothing is written back to the store, so a page reload (new session)
// starts a clean run from seed.
package observatory

import "github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"

// advanceField capitalises alphaBP basis points of each capital's surplus back
// into it — the spiral of accumulation (Vol. I Ch. 24) applied to the orrery
// read-model. TotalPence grows and the M/P/C arcs rescale to the new total
// (CommodityPence absorbs integer rounding so the parts still sum to the total).
// CostPrice, Surplus, Status and TurnoverNumber are unchanged. Pure: returns a
// new slice and never mutates its input. Mirrors store.Memory.AccumulateCapital.
func advanceField(in []circulation.FieldCapital, alphaBP int64) []circulation.FieldCapital {
	out := make([]circulation.FieldCapital, len(in))
	copy(out, in)
	if alphaBP <= 0 {
		return out
	}
	for i := range out {
		fc := out[i]
		if fc.SurplusPence <= 0 {
			continue
		}
		delta := circulation.Pence(int64(fc.SurplusPence) * alphaBP / 10000)
		if delta <= 0 {
			continue
		}
		oldTotal := fc.TotalPence
		newTotal := oldTotal + delta
		if oldTotal > 0 {
			money := fc.MoneyPence * newTotal / oldTotal
			prod := fc.ProductionPence * newTotal / oldTotal
			fc.MoneyPence = money
			fc.ProductionPence = prod
			fc.CommodityPence = newTotal - money - prod
		} else {
			fc.MoneyPence = newTotal
			fc.ProductionPence = 0
			fc.CommodityPence = 0
		}
		fc.TotalPence = newTotal
		out[i] = fc
	}
	return out
}
