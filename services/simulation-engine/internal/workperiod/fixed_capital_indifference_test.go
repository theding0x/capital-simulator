package workperiod_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

// Vol. II Ch. 12: the working period is a matter of circulating capital. Marx
// is explicit that "the distinction between fixed and circulating capital has
// nothing to do with this", and that a steam-engine renewed only after twenty
// years "continues to function in the process of production for several
// working periods". These tests hold the engine fixed and vary the working
// period: the spinner's 7-day discrete period against the locomotive maker's
// 84-day connected one (Marx's twelve-to-one). The engine's wear per working
// day is the same in both, so only the circulating advance scales with the
// period.

const (
	ch12EngineLifeMinutes int64 = 20 * 525_600 // renewed "only after twenty years"
	ch12WorkingDayMinutes int64 = 600
	// Circulating outlay per working day, and the engine's wear per working day
	// (1 penny per minute of service x 600 minutes), giving a 1:5 wear-to-outlay
	// proportion that the working period must not change.
	ch12DailyOutlay = composition.Pence(3_000)
	ch12DailyWear   = composition.Pence(600)
)

// ch12SteamEngine is worth 1 penny per minute of its service life, so the
// linear-wear arithmetic is exact.
func ch12SteamEngine() composition.FixedCapitalItem {
	return composition.FixedCapitalItem{
		PencePurchased:     composition.Pence(ch12EngineLifeMinutes),
		ServiceLifeMinutes: ch12EngineLifeMinutes,
	}
}

type ch12Process struct {
	name string
	wp   workperiod.WorkingPeriod
}

func ch12Processes() []ch12Process {
	return []ch12Process{
		{"spinner, 7-day discrete", workperiod.WorkingPeriod{
			WorkingDayCount: 7, Mode: workperiod.ModeDiscrete, WorkingDayLengthMinutes: ch12WorkingDayMinutes,
		}},
		{"locomotive maker, 84-day connected", workperiod.WorkingPeriod{
			WorkingDayCount: 84, Mode: workperiod.ModeConnected, WorkingDayLengthMinutes: ch12WorkingDayMinutes,
		}},
	}
}

// Fixed capital is indifferent to working-period length: the engine wears the
// same per working day whether it serves a discrete or a connected process.
func TestFixedCapital_WearPerWorkingDayIndependentOfWorkingPeriod(t *testing.T) {
	t.Parallel()
	engine := ch12SteamEngine()
	oneDay := composition.ComputeLinearWear(engine, ch12WorkingDayMinutes, composition.WearUse).Pence
	if oneDay != ch12DailyWear {
		t.Fatalf("one working day's wear = %d, want %d", oneDay, ch12DailyWear)
	}
	for _, p := range ch12Processes() {
		t.Run(p.name, func(t *testing.T) {
			t.Parallel()
			periodMinutes := p.wp.WorkingDayCount * p.wp.WorkingDayLengthMinutes
			wear := composition.ComputeLinearWear(engine, periodMinutes, composition.WearUse).Pence
			if got := wear / composition.Pence(p.wp.WorkingDayCount); got != oneDay {
				t.Errorf("wear per working day over a %d-day period = %d, want %d (the single-day wear)",
					p.wp.WorkingDayCount, got, oneDay)
			}
		})
	}
}

// Only the circulating advance scales with the working period: it is twelve
// times as large for the locomotive maker, while the engine's wear stays the
// same 1:5 share of a period's outlay in both.
func TestCirculatingAdvanceScalesWithWorkingPeriod_FixedWearShareDoesNot(t *testing.T) {
	t.Parallel()
	engine := ch12SteamEngine()
	procs := ch12Processes()
	advances := make([]composition.Pence, len(procs))
	for i, p := range procs {
		periodMinutes := p.wp.WorkingDayCount * p.wp.WorkingDayLengthMinutes
		wear := composition.ComputeLinearWear(engine, periodMinutes, composition.WearUse).Pence
		advance := ch12DailyOutlay * composition.Pence(p.wp.Multiplier())
		if wear*5 != advance {
			t.Errorf("%s: wear %d is not 1/5 of the circulating advance %d", p.name, wear, advance)
		}
		advances[i] = advance
	}
	if advances[1] != 12*advances[0] {
		t.Errorf("circulating advance %d, want 12 x %d for the 84-day period against the 7-day one",
			advances[1], advances[0])
	}
}
