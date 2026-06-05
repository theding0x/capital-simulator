package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestGetObservatorySnapshot(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	a, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 500000, EconomyMode: circulation.EconomyMoney,
	})
	b, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 300000, EconomyMode: circulation.EconomyMoney,
	})
	_, _ = m.Snapshot(ctx, a.ID, circulation.StageDistribution{
		MoneyPence: 100000, ProductionPence: 300000, CommodityPence: 100000,
	})
	_, _ = m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{
		IndustrialCapitalID: a.ID, Period: "1871",
		DemandPence: 400000, SupplyPence: 480000, ExcessPence: 80000,
	})
	_, _ = m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{
		IndustrialCapitalID: b.ID, Period: "1871",
		DemandPence: 250000, SupplyPence: 275000, ExcessPence: 25000,
	})

	h := New(nil, Deps{IndustrialCapitals: m}) // Scheduler nil → tick 0, not running

	req := httptest.NewRequest(http.MethodGet, "/v1/observatory/snapshot", nil)
	rec := httptest.NewRecorder()
	h.GetObservatorySnapshot(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp observatorySnapshotResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Capitals == nil {
		t.Fatal("capitals must never be null")
	}
	if len(resp.Capitals) != 2 {
		t.Fatalf("capitals = %d, want 2", len(resp.Capitals))
	}
	// ΣS = 105000, ΣC = 650000 → p̄′ = round(10000*105000/650000) = 1615 bp.
	if got := resp.Aggregate.AvgRateOfProfitBP; got != 1615 {
		t.Errorf("avg_rate_of_profit_bp = %d, want 1615", got)
	}
	if resp.Aggregate.TotalSocialCapitalPence != 800000 {
		t.Errorf("total social capital = %d, want 800000", resp.Aggregate.TotalSocialCapitalPence)
	}
	if resp.Aggregate.SurplusPence != 105000 || resp.Aggregate.CostPricePence != 650000 {
		t.Errorf("aggregate cost/surplus = %d/%d", resp.Aggregate.CostPricePence, resp.Aggregate.SurplusPence)
	}
	for _, c := range resp.Capitals {
		if c.TurnoverNumber < 1 {
			t.Errorf("capital %s turnover_number = %d, want >= 1 (memory default)", c.ID, c.TurnoverNumber)
		}
	}
}

func TestGetObservatorySnapshotEmptyFieldNeverNull(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{IndustrialCapitals: store.NewMemory()})
	rec := httptest.NewRecorder()
	h.GetObservatorySnapshot(rec, httptest.NewRequest(http.MethodGet, "/v1/observatory/snapshot", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var resp observatorySnapshotResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Capitals == nil {
		t.Fatal("empty field: capitals must be [] not null")
	}
	if resp.Aggregate.AvgRateOfProfitBP != 0 {
		t.Errorf("empty field p̄′ = %d, want 0 (no divide-by-zero)", resp.Aggregate.AvgRateOfProfitBP)
	}
}
