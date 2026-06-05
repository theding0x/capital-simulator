package store_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestMemoryAbodeRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	// Unseeded → the default initial abode.
	got, err := m.GetAbodeState(ctx)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Period != 0 || got.VariablePence != 300000 {
		t.Fatalf("default abode = %+v", got)
	}

	// Advance one period and persist it + its recorded series point.
	next, period := simulation.AdvanceGeneralLaw(got)
	if err := m.AdvanceAbode(ctx, next, period); err != nil {
		t.Fatalf("advance: %v", err)
	}
	got2, _ := m.GetAbodeState(ctx)
	if got2.Period != 1 {
		t.Errorf("period after advance = %d, want 1", got2.Period)
	}

	series, err := m.ListGeneralLawPeriods(ctx, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(series) != 1 {
		t.Fatalf("series len = %d, want 1", len(series))
	}
	if series[0].Period != 1 {
		t.Errorf("series[0].Period = %d, want 1", series[0].Period)
	}

	// Ascending order is preserved across several advances.
	cur := got2
	for i := 0; i < 5; i++ {
		n, p := simulation.AdvanceGeneralLaw(cur)
		_ = m.AdvanceAbode(ctx, n, p)
		cur = n
	}
	series, _ = m.ListGeneralLawPeriods(ctx, 3)
	if len(series) != 3 {
		t.Fatalf("limited series len = %d, want 3", len(series))
	}
	if series[0].Period >= series[2].Period {
		t.Errorf("series not ascending: %d..%d", series[0].Period, series[2].Period)
	}
}
