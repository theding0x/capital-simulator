package observatory

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

func seedRun() *Run {
	return &Run{
		abode: simulation.NewAbodeState(),
		field: []circulation.FieldCapital{{
			ID: "ic1", TotalPence: 300000, MoneyPence: 300000, SurplusPence: 60000, CostPricePence: 240000,
		}},
	}
}

func TestRunAdvanceIncrementsTickGrowsFieldAndSeries(t *testing.T) {
	t.Parallel()
	r := seedRun()
	r.Advance(3)
	snap := r.Snapshot()
	if snap.Tick != 3 {
		t.Fatalf("Tick = %d, want 3", snap.Tick)
	}
	if len(snap.Periods) != 3 {
		t.Fatalf("len(Periods) = %d, want 3", len(snap.Periods))
	}
	if snap.Field[0].TotalPence <= 300000 {
		t.Fatalf("field did not accumulate: %d", snap.Field[0].TotalPence)
	}
}

func TestRunAdvanceZeroIsNoOp(t *testing.T) {
	t.Parallel()
	r := seedRun()
	r.Advance(0)
	snap := r.Snapshot()
	if snap.Tick != 0 || len(snap.Periods) != 0 {
		t.Fatalf("Advance(0) changed state: tick=%d periods=%d", snap.Tick, len(snap.Periods))
	}
}

func TestRunAdvanceClampsAndCapsSeries(t *testing.T) {
	t.Parallel()
	r := seedRun()
	for i := 0; i < 11; i++ {
		r.Advance(1000) // each call clamps to maxAdvancePerPoll
	}
	snap := r.Snapshot()
	if len(snap.Periods) > maxSeries {
		t.Fatalf("len(Periods) = %d, want <= %d", len(snap.Periods), maxSeries)
	}
	if snap.Tick != int64(11*maxAdvancePerPoll) {
		t.Fatalf("Tick = %d, want %d", snap.Tick, 11*maxAdvancePerPoll)
	}
}

func TestRunApplyLeversClampsAndPersistsInSession(t *testing.T) {
	t.Parallel()
	r := seedRun()
	over := int64(999999)
	abode := r.ApplyLevers(simulation.LeverUpdate{SurplusRateBaseBP: &over})
	if abode.SurplusRateBaseBP != 100000 {
		t.Fatalf("returned SurplusRateBaseBP = %d, want clamped 100000", abode.SurplusRateBaseBP)
	}
	if got := r.Snapshot().Abode.SurplusRateBaseBP; got != 100000 {
		t.Fatalf("snapshot SurplusRateBaseBP = %d, want 100000", got)
	}
}
