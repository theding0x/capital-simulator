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

func newProductionTimeHandler(t *testing.T) (*httpapi.Handler, *store.Memory) {
	t.Helper()
	mem := store.NewMemory()
	return httpapi.New(nil, httpapi.Deps{ProductionTime: mem, WorkingPeriods: mem}), mem
}

// --- ProductionLabourGap ---

func TestRecordProductionLabourGap_Created(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)

	body := map[string]any{
		"production_time_nanos": 24192000000000000,
		"labour_time_nanos":     864000000000000,
		"gap_nanos":             23328000000000000,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods/wp-corn/production-labour-gap", bytes.NewReader(b))
	req.SetPathValue("id", "wp-corn")
	w := httptest.NewRecorder()
	h.RecordProductionLabourGap(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRecordProductionLabourGap_BadJSON(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods/wp1/production-labour-gap", bytes.NewBufferString("{bad}"))
	req.SetPathValue("id", "wp1")
	w := httptest.NewRecorder()
	h.RecordProductionLabourGap(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRecordProductionLabourGap_InvalidGap(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	body := map[string]any{
		"production_time_nanos": 1000,
		"labour_time_nanos":     2000, // labour > production — invalid
		"gap_nanos":             -1000,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods/wp1/production-labour-gap", bytes.NewReader(b))
	req.SetPathValue("id", "wp1")
	w := httptest.NewRecorder()
	h.RecordProductionLabourGap(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetProductionLabourGap_NotFound(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/working-periods/missing/production-labour-gap", nil)
	req.SetPathValue("id", "missing")
	w := httptest.NewRecorder()
	h.GetProductionLabourGap(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// --- NaturalProcessActivation ---

func TestRecordNaturalProcessActivation_Created(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	body := map[string]any{
		"kind":                        "agricultural_growth",
		"started_at":                  "1871-04-01T00:00:00Z",
		"requires_labour_maintenance": false,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods/wp-corn/natural-activations", bytes.NewReader(b))
	req.SetPathValue("id", "wp-corn")
	w := httptest.NewRecorder()
	h.RecordNaturalProcessActivation(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRecordNaturalProcessActivation_UnknownKind(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	body := map[string]any{
		"kind":       "not_a_real_kind",
		"started_at": "1871-04-01T00:00:00Z",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods/wp1/natural-activations", bytes.NewReader(b))
	req.SetPathValue("id", "wp1")
	w := httptest.NewRecorder()
	h.RecordNaturalProcessActivation(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListNaturalProcessActivations_EmptySlice(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/working-periods/wp-empty/natural-activations", nil)
	req.SetPathValue("id", "wp-empty")
	w := httptest.NewRecorder()
	h.ListNaturalProcessActivations(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var result []any
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("body not JSON array: %v", err)
	}
	if result == nil {
		t.Error("expected non-nil slice, got null")
	}
}

// --- NaturalSubject ---

func TestRecordNaturalSubject_Created(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	body := map[string]any{
		"commodity_id": "corn",
		"description":  "Corn crop; seed planted April",
		"is_living":    false,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods/wp-corn/natural-subject", bytes.NewReader(b))
	req.SetPathValue("id", "wp-corn")
	w := httptest.NewRecorder()
	h.RecordNaturalSubject(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetNaturalSubject_NotFound(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/working-periods/missing/natural-subject", nil)
	req.SetPathValue("id", "missing")
	w := httptest.NewRecorder()
	h.GetNaturalSubject(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// --- IndustryBenchmark ---

func TestCreateIndustryBenchmark_Created(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	body := map[string]any{
		"industry_name":                 "Agriculture (corn)",
		"typical_production_days_count": 280,
		"typical_labour_days_count":     10,
		"ratio_basis_points":            357,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/industry-benchmarks", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.CreateIndustryBenchmark(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	if w.Header().Get("Location") == "" {
		t.Error("expected Location header")
	}
}

func TestCreateIndustryBenchmark_BadRatio(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	body := map[string]any{
		"industry_name":                 "bad",
		"typical_production_days_count": 10,
		"typical_labour_days_count":     5,
		"ratio_basis_points":            99999,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/industry-benchmarks", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.CreateIndustryBenchmark(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListIndustryBenchmarks_NonNilSlice(t *testing.T) {
	t.Parallel()
	h, _ := newProductionTimeHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/industry-benchmarks", nil)
	w := httptest.NewRecorder()
	h.ListIndustryBenchmarks(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var result []any
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("body not JSON array: %v", err)
	}
	if result == nil {
		t.Error("expected non-nil slice, got null")
	}
}
