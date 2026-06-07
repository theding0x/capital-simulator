package store_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestMemoryAccumulateCapital(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	ic, err := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 500000, EconomyMode: circulation.EconomyMoney,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := m.Snapshot(ctx, ic.ID, circulation.StageDistribution{
		MoneyPence: 100000, ProductionPence: 300000, CommodityPence: 100000,
	}); err != nil {
		t.Fatalf("snapshot: %v", err)
	}

	// Accumulate £1,000 (100000 pence) into the capital.
	grown, err := m.AccumulateCapital(ctx, ic.ID, 100000)
	if err != nil {
		t.Fatalf("accumulate: %v", err)
	}
	if grown.TotalPence != 600000 {
		t.Errorf("total = %d, want 600000", grown.TotalPence)
	}

	// The field snapshot's latest distribution rescales, preserving proportions
	// (money 20%, production 60%, commodity 20% of 600000) and summing to total.
	field, _ := m.FieldSnapshot(ctx)
	var fc circulation.FieldCapital
	for _, f := range field {
		if f.ID == ic.ID {
			fc = f
		}
	}
	if fc.TotalPence != 600000 {
		t.Fatalf("field total = %d, want 600000", fc.TotalPence)
	}
	if got := fc.MoneyPence + fc.ProductionPence + fc.CommodityPence; got != 600000 {
		t.Errorf("distribution sum = %d, want 600000", got)
	}
	if fc.ProductionPence != 360000 {
		t.Errorf("production = %d, want 360000 (60%% of 600000)", fc.ProductionPence)
	}

	// Missing capital → ErrNotFound.
	if _, err := m.AccumulateCapital(ctx, circulation.IndustrialCapitalID("nope"), 1000); err != store.ErrNotFound {
		t.Errorf("missing capital err = %v, want ErrNotFound", err)
	}
}
