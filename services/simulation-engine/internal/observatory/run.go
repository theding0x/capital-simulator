package observatory

import (
	"sync"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

const (
	// maxSeries caps the immiseration time-series retained per run (the abode
	// sparkline shows the most recent 60 periods).
	maxSeries = 60
	// maxAdvancePerPoll bounds how many General-Law periods a single poll may
	// advance, so a crafted advance count cannot pin the CPU.
	maxAdvancePerPoll = 10
)

// Run is one in-memory Atlas simulation run, owned by a single browser session.
// It is seeded from the Manager's template and advanced on poll; it is never
// persisted. All access is guarded by mu.
type Run struct {
	mu      sync.Mutex
	abode   simulation.AbodeState         // carries the live levers too
	field   []circulation.FieldCapital    // the orrery
	periods []simulation.GeneralLawPeriod // capped to the last maxSeries
	tick    int64
}

// RunSnapshot is the read-model projection of a Run at one instant; the
// transport layer maps it to the observatory snapshot response.
type RunSnapshot struct {
	Tick         int64
	Abode        simulation.AbodeState
	Readout      simulation.AbodeReadout
	Field        []circulation.FieldCapital
	Periods      []simulation.GeneralLawPeriod
	Circulation  CirculationSnapshot
	Distribution DistributionSnapshot
	// Field aggregate, summed once under the lock so the transport layer does not
	// re-walk the field. GeneralRateBP = round-half-up(ΣS / ΣC).
	SumTotal      int64
	SumCost       int64
	SumSurplus    int64
	GeneralRateBP int64
}

// Advance runs n periods of the General Law (clamped to [0, maxAdvancePerPoll]),
// growing the field by accumulation and appending each period to the capped
// immiseration series.
func (r *Run) Advance(n int) {
	if n < 0 {
		n = 0
	}
	if n > maxAdvancePerPoll {
		n = maxAdvancePerPoll
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := 0; i < n; i++ {
		next, period := simulation.AdvanceGeneralLaw(r.abode)
		r.abode = next
		r.periods = append(r.periods, period)
		if len(r.periods) > maxSeries {
			r.periods = r.periods[len(r.periods)-maxSeries:]
		}
		r.field = advanceField(r.field, r.abode.AccumulationRateBP)
		r.tick++
	}
}

// ApplyLevers perturbs the run's live abode (clamped) and returns the new abode.
func (r *Run) ApplyLevers(u simulation.LeverUpdate) simulation.AbodeState {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.abode = r.abode.ApplyLevers(u)
	return r.abode
}

// Snapshot returns a race-free copy of the run's current state.
func (r *Run) Snapshot() RunSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	field := make([]circulation.FieldCapital, len(r.field))
	copy(field, r.field)
	periods := make([]simulation.GeneralLawPeriod, len(r.periods))
	copy(periods, r.periods)
	readout := r.abode.Readout()

	var sumTotal, sumCost, sumSurplus int64
	for _, fc := range r.field {
		sumTotal += int64(fc.TotalPence)
		sumCost += int64(fc.CostPricePence)
		sumSurplus += int64(fc.SurplusPence)
	}
	generalRateBP := RateBP(sumSurplus, sumCost)

	return RunSnapshot{
		Tick:          r.tick,
		Abode:         r.abode,
		Readout:       readout,
		Field:         field,
		Periods:       periods,
		Circulation:   DeriveCirculation(r.abode, readout, sumSurplus, sumTotal),
		Distribution:  DeriveDistribution(periods, readout, sumSurplus, sumCost, sumTotal, generalRateBP),
		SumTotal:      sumTotal,
		SumCost:       sumCost,
		SumSurplus:    sumSurplus,
		GeneralRateBP: generalRateBP,
	}
}
