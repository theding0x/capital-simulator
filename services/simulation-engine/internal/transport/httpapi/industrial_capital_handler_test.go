package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// Spinning Mill 1871 — Marx's industrial capital fixture.
// Total advance: £4000 (= 400000 pence).

func spinningMillHandler(t *testing.T) *httpapi.Handler {
	t.Helper()
	return httpapi.New(nil, httpapi.Deps{IndustrialCapitals: store.NewMemory()})
}

func TestCreateIndustrialCapital_Created(t *testing.T) {
	t.Parallel()
	h := spinningMillHandler(t)
	body := `{"total_pence":400000,"economy_mode":"money","stagnation_tolerance_ticks":3,"agent_id":"spinning-mill-1871"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/industrial-capitals", bytes.NewBufferString(body))
	h.CreateIndustrialCapital(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Location") == "" {
		t.Fatal("missing Location header")
	}
	var resp map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["id"] == "" {
		t.Error("expected non-empty id")
	}
	if resp["status"] != "active" {
		t.Errorf("expected status active, got %v", resp["status"])
	}
}

func TestCreateIndustrialCapital_BadRequest(t *testing.T) {
	t.Parallel()
	h := spinningMillHandler(t)
	// total_pence = 0 fails validation
	body := `{"total_pence":0,"economy_mode":"money","stagnation_tolerance_ticks":3}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/industrial-capitals", bytes.NewBufferString(body))
	h.CreateIndustrialCapital(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestGetIndustrialCapital_NotFound(t *testing.T) {
	t.Parallel()
	h := spinningMillHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/industrial-capitals/doesnotexist", nil)
	req.SetPathValue("id", "doesnotexist")
	h.GetIndustrialCapital(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestListIndustrialCapitals_NonNull(t *testing.T) {
	t.Parallel()
	h := spinningMillHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/industrial-capitals", nil)
	h.ListIndustrialCapitals(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var list []any
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if list == nil {
		t.Error("expected non-null slice")
	}
}

func createTestIC(t *testing.T, h *httpapi.Handler) string {
	t.Helper()
	body := `{"total_pence":400000,"economy_mode":"money","stagnation_tolerance_ticks":3,"agent_id":"spinning-mill-1871"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/industrial-capitals", bytes.NewBufferString(body))
	h.CreateIndustrialCapital(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("setup: create IC failed %d: %s", rec.Code, rec.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	return resp["id"].(string)
}

func TestSnapshotStageDistribution_Created(t *testing.T) {
	t.Parallel()
	// Nebeneinander: 133333+133334+133333 = 400000
	h := spinningMillHandler(t)
	id := createTestIC(t, h)
	body := `{"money_pence":133333,"production_pence":133334,"commodity_pence":133333}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/industrial-capitals/"+id+"/snapshot", bytes.NewBufferString(body))
	req.SetPathValue("id", id)
	h.SnapshotStageDistribution(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestSnapshotStageDistribution_TotalMismatch(t *testing.T) {
	t.Parallel()
	h := spinningMillHandler(t)
	id := createTestIC(t, h)
	// 100000+100000+100000 = 300000 ≠ 400000
	body := `{"money_pence":100000,"production_pence":100000,"commodity_pence":100000}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/industrial-capitals/"+id+"/snapshot", bytes.NewBufferString(body))
	req.SetPathValue("id", id)
	h.SnapshotStageDistribution(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestSetSinkingFund_Created(t *testing.T) {
	t.Parallel()
	// Marx: £4000 fixed capital, 10-year lifetime → £400/yr = 40000 pence/yr
	h := spinningMillHandler(t)
	id := createTestIC(t, h)
	body := `{"fixed_capital_pence":400000,"lifetime_years":10,"annual_payment_pence":40000}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.SetPathValue("id", id)
	h.SetSinkingFund(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestTickSinkingFund_Accumulates(t *testing.T) {
	t.Parallel()
	h := spinningMillHandler(t)
	id := createTestIC(t, h)
	body := `{"fixed_capital_pence":400000,"lifetime_years":10,"annual_payment_pence":40000}`
	recSet := httptest.NewRecorder()
	reqSet := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	reqSet.SetPathValue("id", id)
	h.SetSinkingFund(recSet, reqSet)
	recTick := httptest.NewRecorder()
	reqTick := httptest.NewRequest(http.MethodPost, "/", nil)
	reqTick.SetPathValue("id", id)
	h.TickSinkingFund(recTick, reqTick)
	if recTick.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recTick.Code, recTick.Body.String())
	}
	var sf map[string]any
	_ = json.NewDecoder(recTick.Body).Decode(&sf)
	if sf["accumulated_pence"].(float64) != 40000 {
		t.Errorf("expected accumulated_pence=40000 after one tick, got %v", sf["accumulated_pence"])
	}
}

func TestTickSinkingFund_NotFound(t *testing.T) {
	t.Parallel()
	h := spinningMillHandler(t)
	id := createTestIC(t, h)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.SetPathValue("id", id)
	h.TickSinkingFund(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 (no sinking fund yet), got %d", rec.Code)
	}
}

func TestRecordValueRevolution_SetFree(t *testing.T) {
	t.Parallel()
	// Prices fall → money set free
	h := spinningMillHandler(t)
	id := createTestIC(t, h)
	body := `{"affecting":"means_of_production","direction_pence":-5000}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.SetPathValue("id", id)
	h.RecordValueRevolution(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp["set_free"] == nil {
		t.Error("expected set_free to be populated when prices fall")
	}
	if resp["tied_up"] != nil {
		t.Error("expected tied_up to be nil when prices fall")
	}
}

func TestComputeAndRecordSupplyDemand(t *testing.T) {
	t.Parallel()
	// Marx §: "if composition is 80c+20v+20s, demand=100, supply=120"
	h := spinningMillHandler(t)
	id := createTestIC(t, h)
	body := `{"constant_pence":80,"variable_pence":20,"surplus_pence":20,"period":"1871-Q1"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.SetPathValue("id", id)
	h.ComputeAndRecordSupplyDemand(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var sdi map[string]any
	_ = json.NewDecoder(rec.Body).Decode(&sdi)
	if sdi["demand_pence"].(float64) != 100 {
		t.Errorf("demand_pence: want 100, got %v", sdi["demand_pence"])
	}
	if sdi["supply_pence"].(float64) != 120 {
		t.Errorf("supply_pence: want 120, got %v", sdi["supply_pence"])
	}
	if sdi["excess_pence"].(float64) != 20 {
		t.Errorf("excess_pence: want 20, got %v", sdi["excess_pence"])
	}
}

func TestAggregateSupplyDemand_MissingPeriod(t *testing.T) {
	t.Parallel()
	h := spinningMillHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/industrial-capitals/supply-demand/aggregate", nil)
	h.AggregateSupplyDemand(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing period, got %d", rec.Code)
	}
}
