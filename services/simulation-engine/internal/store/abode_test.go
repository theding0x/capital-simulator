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

func TestMemorySetAbodeLevers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	zero := int64(0)
	got, err := m.SetAbodeLevers(ctx, simulation.LeverUpdate{AccumulationRateBP: &zero})
	if err != nil {
		t.Fatalf("set: %v", err)
	}
	if got.AccumulationRateBP != 0 {
		t.Errorf("alpha = %d, want 0", got.AccumulationRateBP)
	}
	// The other parameters are untouched (default base wage 2500).
	if got.BaseWagePence != 2500 {
		t.Errorf("base wage bled = %d, want 2500", got.BaseWagePence)
	}

	// Persisted: a later Get reflects the lever.
	st, _ := m.GetAbodeState(ctx)
	if st.AccumulationRateBP != 0 {
		t.Errorf("not persisted: %d", st.AccumulationRateBP)
	}

	// With α = 0 the law performs simple reproduction: no surplus is
	// capitalised, so total social capital (c+v) does not grow when advanced
	// (displacement only shifts value between c and v).
	next, _ := simulation.AdvanceGeneralLaw(st)
	if next.ConstantPence+next.VariablePence != st.ConstantPence+st.VariablePence {
		t.Errorf("alpha=0 should not grow total capital: %d -> %d",
			st.ConstantPence+st.VariablePence, next.ConstantPence+next.VariablePence)
	}
}
