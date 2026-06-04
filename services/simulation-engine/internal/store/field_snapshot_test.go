package store_test

import (
	"context"
	"sort"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestMemoryFieldSnapshot(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	// Two capitals of differing organic composition (Vol. III Ch. 9 spheres).
	a, err := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 500000, EconomyMode: circulation.EconomyMoney,
	})
	if err != nil {
		t.Fatalf("create a: %v", err)
	}
	b, err := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 300000, EconomyMode: circulation.EconomyMoney,
	})
	if err != nil {
		t.Fatalf("create b: %v", err)
	}

	if _, err := m.Snapshot(ctx, a.ID, circulation.StageDistribution{
		MoneyPence: 100000, ProductionPence: 300000, CommodityPence: 100000,
	}); err != nil {
		t.Fatalf("snapshot a: %v", err)
	}
	if _, err := m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{
		IndustrialCapitalID: a.ID, Period: "1871",
		DemandPence: 400000, SupplyPence: 480000, ExcessPence: 80000, // p' = 20%
	}); err != nil {
		t.Fatalf("supply-demand a: %v", err)
	}
	if _, err := m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{
		IndustrialCapitalID: b.ID, Period: "1871",
		DemandPence: 250000, SupplyPence: 275000, ExcessPence: 25000, // p' = 10%
	}); err != nil {
		t.Fatalf("supply-demand b: %v", err)
	}

	field, err := m.FieldSnapshot(ctx)
	if err != nil {
		t.Fatalf("FieldSnapshot: %v", err)
	}
	if len(field) != 2 {
		t.Fatalf("want 2 capitals, got %d", len(field))
	}
	sort.Slice(field, func(i, j int) bool { return field[i].ID < field[j].ID })

	byID := map[circulation.IndustrialCapitalID]circulation.FieldCapital{}
	for _, fc := range field {
		byID[fc.ID] = fc
	}
	fa := byID[a.ID]
	if fa.MoneyPence != 100000 || fa.ProductionPence != 300000 || fa.CommodityPence != 100000 {
		t.Errorf("a distribution: %+v", fa)
	}
	if fa.CostPricePence != 400000 || fa.SurplusPence != 80000 {
		t.Errorf("a cost/surplus: cost=%d surplus=%d", fa.CostPricePence, fa.SurplusPence)
	}

	// Capital with no distribution recorded defaults to all-money.
	c, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 120000, EconomyMode: circulation.EconomyMoney,
	})
	field, _ = m.FieldSnapshot(ctx)
	for _, fc := range field {
		if fc.ID == c.ID && fc.MoneyPence != 120000 {
			t.Errorf("default-money capital: money=%d want 120000", fc.MoneyPence)
		}
	}
}
