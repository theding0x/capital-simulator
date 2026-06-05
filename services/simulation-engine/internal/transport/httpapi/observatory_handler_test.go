package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

	h := New(nil, Deps{IndustrialCapitals: m, AbodeStates: m}) // Scheduler nil → tick 0, not running

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

	// Snapshot v2 carries the hidden-abode block.
	ab := resp.Abode
	if ab.TotalVariablePence <= 0 {
		t.Errorf("abode Σv = %d, want > 0", ab.TotalVariablePence)
	}
	if ab.LawSeries == nil {
		t.Error("abode.law_series must be non-nil (never null)")
	}
	// s/v == round(10000 * Σs / Σv) — the rate of exploitation.
	wantSV := (ab.TotalSurplusPence*10000 + ab.TotalVariablePence/2) / ab.TotalVariablePence
	if ab.RateOfExploitationBP != wantSV {
		t.Errorf("s/v = %d, want %d (round 10000*Σs/Σv)", ab.RateOfExploitationBP, wantSV)
	}
	if ab.NecessaryLabourMinutes+ab.SurplusLabourMinutes == 0 {
		t.Error("working day minutes not populated")
	}
	if resp.Abode.AccumulationRateBP == 0 && resp.Abode.BaseWagePence == 0 {
		t.Error("abode block missing base lever values (accumulation_rate_bp / base_wage_pence)")
	}
}

func TestSetObservatoryLevers(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	h := New(nil, Deps{AbodeStates: m}) // mirror TestGetObservatorySnapshot's handler construction

	body := `{"accumulation_rate_bp": 0, "surplus_rate_base_bp": 30000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/observatory/levers", strings.NewReader(body))
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
		t.Errorf("untouched wage = %d, want 2500 (default)", resp.BaseWagePence)
	}

	// Persisted on the live abode.
	st, _ := m.GetAbodeState(req.Context())
	if st.AccumulationRateBP != 0 || st.SurplusRateBaseBP != 30000 {
		t.Errorf("not persisted: %+v", st)
	}

	// Malformed body → 400.
	bad := httptest.NewRequest(http.MethodPost, "/v1/observatory/levers", strings.NewReader("{"))
	badRec := httptest.NewRecorder()
	h.SetObservatoryLevers(badRec, bad)
	if badRec.Code != http.StatusBadRequest {
		t.Errorf("malformed code = %d, want 400", badRec.Code)
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
