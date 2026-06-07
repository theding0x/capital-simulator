package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/observatory"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// newTestObservatory seeds a memory store with two Marx-faithful capitals, then
// builds a Manager from the seed snapshot (as main.go does at boot).
func newTestObservatory(t *testing.T) *observatory.Manager {
	t.Helper()
	ctx := context.Background()
	m := store.NewMemory()
	a, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{TotalPence: 500000, EconomyMode: circulation.EconomyMoney})
	b, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{TotalPence: 300000, EconomyMode: circulation.EconomyMoney})
	_, _ = m.Snapshot(ctx, a.ID, circulation.StageDistribution{MoneyPence: 100000, ProductionPence: 300000, CommodityPence: 100000})
	_, _ = m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{IndustrialCapitalID: a.ID, Period: "1871", DemandPence: 400000, SupplyPence: 480000, ExcessPence: 80000})
	_, _ = m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{IndustrialCapitalID: b.ID, Period: "1871", DemandPence: 250000, SupplyPence: 275000, ExcessPence: 25000})
	seedAbode, _ := m.GetAbodeState(ctx)
	seedField, _ := m.FieldSnapshot(ctx)
	return observatory.NewManager(seedAbode, seedField, slog.Default())
}

func snapshot(t *testing.T, h *Handler, sess string, advance int) observatorySnapshotResponse {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/observatory/snapshot?advance=%d", advance), nil)
	if sess != "" {
		req.Header.Set("X-Atlas-Session", sess)
	}
	rec := httptest.NewRecorder()
	h.GetObservatorySnapshot(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("snapshot status = %d; body=%s", rec.Code, rec.Body.String())
	}
	var resp observatorySnapshotResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return resp
}

func TestGetObservatorySnapshotSeed(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{Observatory: newTestObservatory(t)})
	resp := snapshot(t, h, "seed", 0) // advance=0 → pristine seed

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
	if resp.Abode.TotalVariablePence <= 0 {
		t.Errorf("abode Σv = %d, want > 0", resp.Abode.TotalVariablePence)
	}
	if resp.Abode.LawSeries == nil {
		t.Error("abode.law_series must be non-nil (never null)")
	}
	if resp.Abode.AccumulationRateBP == 0 && resp.Abode.BaseWagePence == 0 {
		t.Error("abode block missing base lever values")
	}
}

func TestGetObservatorySnapshotAdvanceAccumulates(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{Observatory: newTestObservatory(t)})
	base := snapshot(t, h, "grow", 0)
	after := snapshot(t, h, "grow", 3) // same session advances 3 periods
	if after.Tick != 3 {
		t.Fatalf("tick = %d, want 3", after.Tick)
	}
	if after.Aggregate.TotalSocialCapitalPence <= base.Aggregate.TotalSocialCapitalPence {
		t.Fatalf("capital did not accumulate: %d <= %d", after.Aggregate.TotalSocialCapitalPence, base.Aggregate.TotalSocialCapitalPence)
	}
}

func TestGetObservatorySnapshotSessionsIndependent(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{Observatory: newTestObservatory(t)})
	_ = snapshot(t, h, "A", 5)
	b := snapshot(t, h, "B", 0)
	if b.Tick != 0 {
		t.Fatalf("session B tick = %d, want 0 (independent of A)", b.Tick)
	}
}

func TestGetObservatorySnapshotEmptyFieldNeverNull(t *testing.T) {
	t.Parallel()
	empty := observatory.NewManager(simulation.NewAbodeState(), nil, slog.Default())
	h := New(nil, Deps{Observatory: empty})
	resp := snapshot(t, h, "e", 0)
	if resp.Capitals == nil {
		t.Fatal("empty field: capitals must be [] not null")
	}
	if resp.Aggregate.AvgRateOfProfitBP != 0 {
		t.Errorf("empty field p̄′ = %d, want 0 (no divide-by-zero)", resp.Aggregate.AvgRateOfProfitBP)
	}
}

func TestSetObservatoryLevers(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{Observatory: newTestObservatory(t)})

	body := `{"accumulation_rate_bp": 0, "surplus_rate_base_bp": 30000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/observatory/levers", strings.NewReader(body))
	req.Header.Set("X-Atlas-Session", "lev")
	rec := httptest.NewRecorder()
	h.SetObservatoryLevers(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		SurplusRateBaseBP  int64 `json:"surplus_rate_base_bp"`
		BaseWagePence      int64 `json:"base_wage_pence"`
		AccumulationRateBP int64 `json:"accumulation_rate_bp"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.AccumulationRateBP != 0 || resp.SurplusRateBaseBP != 30000 {
		t.Errorf("levers not applied: %+v", resp)
	}
	if resp.BaseWagePence != 2500 {
		t.Errorf("untouched wage = %d, want 2500 (seed default)", resp.BaseWagePence)
	}

	// reflected in the SAME session's snapshot
	same := snapshot(t, h, "lev", 0)
	if same.Abode.SurplusRateBaseBP != 30000 || same.Abode.AccumulationRateBP != 0 {
		t.Errorf("levers not reflected in session snapshot: %+v", same.Abode)
	}

	// a DIFFERENT session still sees the seed default (10000)
	other := snapshot(t, h, "other", 0)
	if other.Abode.SurplusRateBaseBP == 30000 {
		t.Error("levers leaked across sessions")
	}

	// malformed body → 400
	bad := httptest.NewRequest(http.MethodPost, "/v1/observatory/levers", strings.NewReader("{"))
	bad.Header.Set("X-Atlas-Session", "lev")
	badRec := httptest.NewRecorder()
	h.SetObservatoryLevers(badRec, bad)
	if badRec.Code != http.StatusBadRequest {
		t.Errorf("malformed code = %d, want 400", badRec.Code)
	}
}
