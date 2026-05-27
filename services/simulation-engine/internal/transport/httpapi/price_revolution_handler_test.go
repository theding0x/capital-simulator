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

func newPriceRevolutionHandler() *httpapi.Handler {
	mem := store.NewMemory()
	return httpapi.New(nil, httpapi.Deps{PriceRevolutions: mem})
}

func TestCreatePriceRevolution_MPEvent(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	// Spinning Mill 1871: cotton price falls 15%.
	body := `{
		"industrial_capital_id": "test-ic-001",
		"affecting": "means_of_production",
		"direction_pence": -7200,
		"basis_points_change": -1500,
		"old_advance_pence": 48000
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/price-revolutions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreatePriceRevolution(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["revalued_capital"] == nil {
		t.Error("expected revalued_capital in response")
	}
	if w.Header().Get("Location") == "" {
		t.Error("expected Location header")
	}
}

func TestCreatePriceRevolution_OutputEvent(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	// 1857 Crisis: output price falls 10%.
	body := `{
		"industrial_capital_id": "test-ic-002",
		"affecting": "output",
		"direction_pence": -6000,
		"basis_points_change": 0,
		"expected_realisation_pence": 60000,
		"actual_realisation_pence": 54000
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/price-revolutions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreatePriceRevolution(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["realised_delta"] == nil {
		t.Error("expected realised_delta in response")
	}
}

func TestCreatePriceRevolution_MissingIndustrialCapitalID(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	req := httptest.NewRequest(http.MethodPost, "/v1/price-revolutions",
		bytes.NewBufferString(`{"affecting":"output"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreatePriceRevolution(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreatePriceRevolution_BadJSON(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	req := httptest.NewRequest(http.MethodPost, "/v1/price-revolutions",
		bytes.NewBufferString(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreatePriceRevolution(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetPriceRevolution_RoundTrip(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	body := `{
		"industrial_capital_id": "test-ic-003",
		"affecting": "means_of_production",
		"direction_pence": -7200,
		"basis_points_change": -1500,
		"old_advance_pence": 48000
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/price-revolutions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreatePriceRevolution(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create failed: %d %s", w.Code, w.Body.String())
	}

	var created map[string]any
	_ = json.NewDecoder(w.Body).Decode(&created)
	evt, _ := created["event"].(map[string]any)
	id, _ := evt["id"].(string)

	getReq := httptest.NewRequest(http.MethodGet, "/v1/price-revolutions/"+id, nil)
	getReq.SetPathValue("id", id)
	gw := httptest.NewRecorder()
	h.GetPriceRevolution(gw, getReq)

	if gw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", gw.Code)
	}
	var got map[string]any
	_ = json.NewDecoder(gw.Body).Decode(&got)
	gotEvt, _ := got["event"].(map[string]any)
	if gotEvt["id"] != id {
		t.Errorf("id mismatch: want %s, got %v", id, gotEvt["id"])
	}
	// inventory_revaluations must be a non-nil slice, never null.
	if got["inventory_revaluations"] == nil {
		t.Error("inventory_revaluations must not be null")
	}
}

func TestGetPriceRevolution_NotFound(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	req := httptest.NewRequest(http.MethodGet, "/v1/price-revolutions/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()
	h.GetPriceRevolution(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestRecordPriceRevolutionCase(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	body := `{
		"industrial_capital_id": "test-ic-004",
		"affecting": "means_of_production",
		"direction_pence": -7200,
		"basis_points_change": -1500,
		"old_advance_pence": 48000
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/price-revolutions", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreatePriceRevolution(w, req)
	var created map[string]any
	_ = json.NewDecoder(w.Body).Decode(&created)
	evt, _ := created["event"].(map[string]any)
	id, _ := evt["id"].(string)

	caseReq := httptest.NewRequest(http.MethodPost,
		"/v1/price-revolutions/"+id+"/price-case",
		bytes.NewBufferString(`{"fall_case":"same_scale"}`))
	caseReq.Header.Set("Content-Type", "application/json")
	caseReq.SetPathValue("id", id)
	cw := httptest.NewRecorder()
	h.RecordPriceRevolutionCase(cw, caseReq)

	if cw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", cw.Code, cw.Body.String())
	}
}

func TestRecordInventoryRevaluation(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	// 5000 lbs × (10p − 12p) = −10000p paper loss.
	body := `{
		"industrial_capital_id": "test-ic-005",
		"commodity_id": "raw-cotton-1871",
		"quantity_held": 5000,
		"old_unit_pence": 12,
		"new_unit_pence": 10
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/inventory-revaluations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.RecordInventoryRevaluation(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["delta_pence"].(float64) != -10000 {
		t.Errorf("expected delta_pence=-10000, got %v", resp["delta_pence"])
	}
}

func TestRecordSpeculativeHold(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	body := `{
		"industrial_capital_id": "test-ic-006",
		"commodity_id": "manchester-cotton-1857",
		"anticipated_direction": 1
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/speculative-holds", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.RecordSpeculativeHold(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetCompoundCapitalAdjustment(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	req := httptest.NewRequest(http.MethodGet,
		"/v1/price-revolutions/compound?industrial_capital_id=test-ic-007&period=1871-Q1", nil)
	w := httptest.NewRecorder()
	h.GetCompoundCapitalAdjustment(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["net_effect_pence"] == nil {
		t.Error("expected net_effect_pence in response")
	}
}

func TestGetCompoundCapitalAdjustment_MissingParams(t *testing.T) {
	t.Parallel()
	h := newPriceRevolutionHandler()

	req := httptest.NewRequest(http.MethodGet, "/v1/price-revolutions/compound", nil)
	w := httptest.NewRecorder()
	h.GetCompoundCapitalAdjustment(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
