package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
)

func newV2Handler() *Handler {
	return New(store.NewMemory(), nil)
}

// --- UpgradeSellingPhase ----------------------------------------------------

func TestUpgradeSellingPhase_Created(t *testing.T) {
	t.Parallel()
	h := newV2Handler()

	body := map[string]any{
		"selling_phase_id":          "phase-001",
		"industrial_capital_id":     "ic-001",
		"market_id":                 "london-bombay",
		"market_distance_meters":    9000000,
		"communication_lag_minutes": 360000,
		"buyer_search_minutes":      72000,
		"bargaining_minutes":        72000,
		"legal_settlement_minutes":  21600,
		"total_minutes":             525600,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/turnover-time/tt-001/selling-phase/phase-001/detail", bytes.NewReader(b))
	req.SetPathValue("id", "tt-001")
	req.SetPathValue("phaseId", "phase-001")
	w := httptest.NewRecorder()

	h.UpgradeSellingPhase(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["id"] == "" || resp["id"] == nil {
		t.Fatal("expected non-empty id in response")
	}
}

func TestUpgradeSellingPhase_BadJSON(t *testing.T) {
	t.Parallel()
	h := newV2Handler()
	req := httptest.NewRequest("POST", "/v1/turnover-time/tt/selling-phase/sp/detail", bytes.NewReader([]byte("not-json")))
	req.SetPathValue("id", "tt")
	req.SetPathValue("phaseId", "sp")
	w := httptest.NewRecorder()
	h.UpgradeSellingPhase(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpgradeSellingPhase_SubSpanMismatch(t *testing.T) {
	t.Parallel()
	h := newV2Handler()
	body := map[string]any{
		"selling_phase_id":          "sp",
		"total_minutes":             1000,
		"buyer_search_minutes":      100,
		"bargaining_minutes":        100,
		"communication_lag_minutes": 100,
		"legal_settlement_minutes":  100,
		// sum = 400 != 1000
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/turnover-time/tt/selling-phase/sp/detail", bytes.NewReader(b))
	req.SetPathValue("id", "tt")
	req.SetPathValue("phaseId", "sp")
	w := httptest.NewRecorder()
	h.UpgradeSellingPhase(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpgradeSellingPhase_LagExceedsMax_Returns422(t *testing.T) {
	t.Parallel()
	h := newV2Handler()

	// Create a distance-lag relation: london-manchester 320 km, 1 min/km.
	relBody := map[string]any{
		"market_id":          "london-manchester",
		"distance_meters":    320000,
		"lag_minutes_per_km": 1,
	}
	rb, _ := json.Marshal(relBody)
	rreq := httptest.NewRequest("POST", "/v1/distance-lag", bytes.NewReader(rb))
	rw := httptest.NewRecorder()
	h.CreateDistanceLagRelation(rw, rreq)
	if rw.Code != http.StatusCreated {
		t.Fatalf("setup: expected 201, got %d", rw.Code)
	}
	var relResp map[string]any
	json.NewDecoder(rw.Body).Decode(&relResp)
	relID := relResp["id"].(string)

	// Phase with comm lag 1000 > 1*320=320 min.
	body := map[string]any{
		"selling_phase_id":          "sp",
		"market_id":                 "london-manchester",
		"market_distance_meters":    320000,
		"communication_lag_minutes": 1000,
		"buyer_search_minutes":      100,
		"bargaining_minutes":        100,
		"legal_settlement_minutes":  100,
		"total_minutes":             1300,
		"distance_lag_relation_id":  relID,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/turnover-time/tt/selling-phase/sp/detail", bytes.NewReader(b))
	req.SetPathValue("id", "tt")
	req.SetPathValue("phaseId", "sp")
	w := httptest.NewRecorder()
	h.UpgradeSellingPhase(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d: %s", w.Code, w.Body.String())
	}
}

// --- CreateDistanceLagRelation -----------------------------------------------

func TestCreateDistanceLagRelation_Created(t *testing.T) {
	t.Parallel()
	h := newV2Handler()
	body := map[string]any{
		"market_id":          "london-bombay",
		"distance_meters":    9000000,
		"lag_minutes_per_km": 58,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/distance-lag", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.CreateDistanceLagRelation(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateDistanceLagRelation_MissingMarketID(t *testing.T) {
	t.Parallel()
	h := newV2Handler()
	body := map[string]any{"distance_meters": 9000000, "lag_minutes_per_km": 58}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/distance-lag", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.CreateDistanceLagRelation(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListDistanceLagRelations_ReturnsNonNullSlice(t *testing.T) {
	t.Parallel()
	h := newV2Handler()
	req := httptest.NewRequest("GET", "/v1/distance-lag", nil)
	w := httptest.NewRecorder()
	h.ListDistanceLagRelations(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["items"] == nil {
		t.Fatal("expected non-null items slice")
	}
}

// --- CreateCirculationSpeedImprovement --------------------------------------

func TestCreateCirculationSpeedImprovement_Created(t *testing.T) {
	t.Parallel()
	h := newV2Handler()
	body := map[string]any{
		"market_id":              "london-bombay",
		"old_lag_minutes_per_km": 58,
		"new_lag_minutes_per_km": 1,
		"reason":                 "telegraph",
		"effective_at":           "1870-01-01T00:00:00Z",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/circulation-speed-improvements", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.CreateCirculationSpeedImprovement(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateCirculationSpeedImprovement_BadEffectiveAt(t *testing.T) {
	t.Parallel()
	h := newV2Handler()
	body := map[string]any{
		"market_id":    "london-bombay",
		"reason":       "telegraph",
		"effective_at": "not-a-date",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/circulation-speed-improvements", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.CreateCirculationSpeedImprovement(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListCirculationSpeedImprovements_ReturnsNonNullSlice(t *testing.T) {
	t.Parallel()
	h := newV2Handler()
	req := httptest.NewRequest("GET", "/v1/circulation-speed-improvements", nil)
	w := httptest.NewRecorder()
	h.ListCirculationSpeedImprovements(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["items"] == nil {
		t.Fatal("expected non-null items slice")
	}
}

// --- GetAnnualSurplusPenalty ------------------------------------------------

func TestGetAnnualSurplusPenalty_MissingPeriod(t *testing.T) {
	t.Parallel()
	h := newV2Handler()
	req := httptest.NewRequest("GET", "/v1/turnover-time/tt/annual-surplus-penalty", nil)
	req.SetPathValue("id", "tt")
	w := httptest.NewRecorder()
	h.GetAnnualSurplusPenalty(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetAnnualSurplusPenalty_NotFoundTurnover(t *testing.T) {
	t.Parallel()
	h := newV2Handler()
	req := httptest.NewRequest("GET", "/v1/turnover-time/missing/annual-surplus-penalty?period=1871", nil)
	req.SetPathValue("id", "missing-id")
	w := httptest.NewRecorder()
	h.GetAnnualSurplusPenalty(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
