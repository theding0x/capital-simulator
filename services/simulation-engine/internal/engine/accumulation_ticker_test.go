package engine_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestAccumulationTickerGrowsCapital(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	ic, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 500000, EconomyMode: circulation.EconomyMoney,
	})
	_, _ = m.Snapshot(ctx, ic.ID, circulation.StageDistribution{
		MoneyPence: 100000, ProductionPence: 300000, CommodityPence: 100000,
	})
	_, _ = m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{
		IndustrialCapitalID: ic.ID, Period: "1871",
		DemandPence: 400000, SupplyPence: 480000, ExcessPence: 80000, // s = 80000
	})

	// α = 5000 bp (50% of surplus reinvested per pass).
	tk := engine.NewAccumulationTicker(m, 5000)
	if tk.Name() != "accumulation" {
		t.Fatalf("name = %q", tk.Name())
	}

	advanced, err := tk.Tick(ctx)
	if err != nil {
		t.Fatalf("tick: %v", err)
	}
	if advanced != 1 {
		t.Fatalf("advanced = %d, want 1", advanced)
	}

	// total grew by α·s = 0.5 * 80000 = 40000 → 540000.
	got, _ := m.GetIndustrialCapital(ctx, ic.ID)
	if got.TotalPence != 540000 {
		t.Errorf("total after one pass = %d, want 540000", got.TotalPence)
	}

	// A second pass grows it again (linear in Slice 1).
	_, _ = tk.Tick(ctx)
	got, _ = m.GetIndustrialCapital(ctx, ic.ID)
	if got.TotalPence != 580000 {
		t.Errorf("total after two passes = %d, want 580000", got.TotalPence)
	}
}
