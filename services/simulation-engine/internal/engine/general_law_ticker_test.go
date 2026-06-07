package engine_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestGeneralLawTickerAdvancesTheAbode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	tk := engine.NewGeneralLawTicker(m)
	if tk.Name() != "general-law" {
		t.Fatalf("name = %q", tk.Name())
	}

	before, _ := m.GetAbodeState(ctx)
	advanced, err := tk.Tick(ctx)
	if err != nil {
		t.Fatalf("tick: %v", err)
	}
	if advanced != 1 {
		t.Fatalf("advanced = %d, want 1", advanced)
	}

	after, _ := m.GetAbodeState(ctx)
	if after.Period != before.Period+1 {
		t.Errorf("period %d -> %d, want +1", before.Period, after.Period)
	}

	// Three more passes; the immiseration series grows and the reserve army is
	// larger at the end than at the start (the law's direction).
	for i := 0; i < 3; i++ {
		if _, err := tk.Tick(ctx); err != nil {
			t.Fatalf("tick %d: %v", i, err)
		}
	}
	series, _ := m.ListGeneralLawPeriods(ctx, 0)
	if len(series) != 4 {
		t.Fatalf("series len = %d, want 4", len(series))
	}
	if series[len(series)-1].ReserveArmyCount <= series[0].ReserveArmyCount {
		t.Errorf("reserve army did not grow: %d -> %d",
			series[0].ReserveArmyCount, series[len(series)-1].ReserveArmyCount)
	}
}
