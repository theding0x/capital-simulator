package httpapi_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
	"github.com/theding0x/capital-simulator/services/market-service/internal/transport/httpapi"
)

func newTTHandler(t *testing.T) *httpapi.Handler {
	t.Helper()
	return httpapi.New(store.NewMemory(), slog.Default())
}

// createSpinningMill1871TT creates a TurnoverTime for the SpinningMill1871 fixture
// and returns the created id.
func createSpinningMill1871TT(t *testing.T, h *httpapi.Handler) string {
	t.Helper()
	body := `{"industrial_capital_id":"5eed000000000000000502"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/turnover-time", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateTurnoverTime(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create TurnoverTime status = %d, want 201; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	id, _ := resp["id"].(string)
	if id == "" {
		t.Fatal("create TurnoverTime response missing id")
	}
	return id
}

func TestCreateTurnoverTime_Created(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	body := `{"industrial_capital_id":"5eed000000000000000502"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/turnover-time", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateTurnoverTime(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] == nil || resp["id"] == "" {
		t.Error("response missing id")
	}
	if resp["industrial_capital_id"] != "5eed000000000000000502" {
		t.Errorf("industrial_capital_id = %v, want 5eed000000000000000502", resp["industrial_capital_id"])
	}
}

func TestCreateTurnoverTime_MissingIndustrialCapitalID(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	body := `{"industrial_capital_id":""}`
	req := httptest.NewRequest(http.MethodPost, "/v1/turnover-time", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateTurnoverTime(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for missing industrial_capital_id", w.Code)
	}
}

func TestCreateTurnoverTime_BadJSON(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/v1/turnover-time", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateTurnoverTime(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", w.Code)
	}
}

func TestGetTurnoverTime_NotFound(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/turnover-time/doesnotexist", nil)
	req.SetPathValue("id", "doesnotexist")
	w := httptest.NewRecorder()
	h.GetTurnoverTime(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestListTurnoverTimes_NonNullEmpty(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/turnover-time", nil)
	w := httptest.NewRecorder()
	h.ListTurnoverTimes(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["items"] == nil {
		t.Error("items must not be null — empty slice must be []")
	}
}

func TestAddLabourTime_UpdatesProduction(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	id := createSpinningMill1871TT(t, h)

	// 1-week labour time in nanoseconds.
	weekNanos := int64(7 * 24 * 60 * 60 * 1_000_000_000)
	body := `{"nanos":604800000000000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/turnover-time/"+id+"/labour-time", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.AddLabourTime(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("AddLabourTime status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	prod, _ := resp["production"].(map[string]any)
	if prod == nil {
		t.Fatal("response missing production")
	}
	got, _ := prod["labour_time_nanos"].(float64)
	if int64(got) != weekNanos {
		t.Errorf("labour_time_nanos = %d, want %d", int64(got), weekNanos)
	}
}

func TestAddLabourTime_NotFound(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	body := `{"nanos":1000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/turnover-time/ghost/labour-time", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "ghost")
	w := httptest.NewRecorder()
	h.AddLabourTime(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

// TestAddLabourTime_ConcurrentCirculationRejected verifies the invariant:
// production sub-spans cannot be added while a selling phase is open.
// §: "production time and circulation time are mutually exclusive."
func TestAddLabourTime_ConcurrentCirculationRejected(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	id := createSpinningMill1871TT(t, h)

	// Open a selling phase.
	sellBody := `{"industrial_capital_id":"5eed000000000000000502","commodity_circuit_id":"ccid","opened_at":"1871-01-01T00:00:00Z"}`
	sr := httptest.NewRequest(http.MethodPost, "/v1/turnover-time/"+id+"/selling-phase", strings.NewReader(sellBody))
	sr.Header.Set("Content-Type", "application/json")
	sr.SetPathValue("id", id)
	sw := httptest.NewRecorder()
	h.RecordSellingPhase(sw, sr)
	if sw.Code != http.StatusCreated {
		t.Fatalf("open selling phase status = %d, want 201", sw.Code)
	}

	// Try to add labour time while selling phase is open — must be rejected.
	labourBody := `{"nanos":1000}`
	lr := httptest.NewRequest(http.MethodPost, "/v1/turnover-time/"+id+"/labour-time", strings.NewReader(labourBody))
	lr.Header.Set("Content-Type", "application/json")
	lr.SetPathValue("id", id)
	lw := httptest.NewRecorder()
	h.AddLabourTime(lw, lr)
	if lw.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 (concurrent production and circulation)", lw.Code)
	}
}

func TestGetActiveFraction_ReturnsValidRange(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	id := createSpinningMill1871TT(t, h)

	// Add 1-week labour time.
	labourBody := `{"nanos":604800000000000}`
	lr := httptest.NewRequest(http.MethodPost, "/v1/turnover-time/"+id+"/labour-time", strings.NewReader(labourBody))
	lr.Header.Set("Content-Type", "application/json")
	lr.SetPathValue("id", id)
	lw := httptest.NewRecorder()
	h.AddLabourTime(lw, lr)
	if lw.Code != http.StatusOK {
		t.Fatalf("AddLabourTime: %d", lw.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/turnover-time/"+id+"/active-fraction", nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.GetActiveFraction(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GetActiveFraction status = %d; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	bp, _ := resp["active_fraction_basis_points"].(float64)
	// With no circulation, active fraction must be 10000 (100%).
	if int64(bp) != 10000 {
		t.Errorf("active_fraction_basis_points = %v, want 10000 (no circulation)", bp)
	}
}

func TestGetActiveFraction_NotFound(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/turnover-time/ghost/active-fraction", nil)
	req.SetPathValue("id", "ghost")
	w := httptest.NewRecorder()
	h.GetActiveFraction(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestSetPerishability_Created(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	// Milk: 2-day window (172800000000000 ns).
	body := `{"commodity_id":"milk-cid","window_nanos":172800000000000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/perishability", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.SetPerishability(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("SetPerishability status = %d, want 201; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] == nil || resp["id"] == "" {
		t.Error("response missing id")
	}
	if resp["commodity_id"] != "milk-cid" {
		t.Errorf("commodity_id = %v, want milk-cid", resp["commodity_id"])
	}
	if resp["window_nanos"] != float64(172800000000000) {
		t.Errorf("window_nanos = %v, want 172800000000000", resp["window_nanos"])
	}
}

func TestSetMarketSeparation_Created(t *testing.T) {
	t.Parallel()
	h := newTTHandler(t)
	body := `{"industrial_capital_id":"5eed000000000000000502","selling_market_id":"manchester","buying_market_id":"liverpool"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/market-separation", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.SetMarketSeparation(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("SetMarketSeparation status = %d, want 201; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] == nil || resp["id"] == "" {
		t.Error("response missing id")
	}
}
