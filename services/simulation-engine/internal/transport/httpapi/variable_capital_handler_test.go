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

func newValHandler() *httpapi.Handler {
	mem := store.NewMemory()
	return httpapi.New(nil, httpapi.Deps{AnnualSurplusRates: mem})
}

func TestCreateVariableCapitalAdvance_Created(t *testing.T) {
	t.Parallel()
	h := newValHandler()

	// Capitalist A: £500 (120000p), 5-week period, n=10.
	body := `{
		"industrial_capital_id": "test-ic-001",
		"pence": 120000,
		"turnover_period_minutes": 50400,
		"turnovers_per_year_basis_points": 100000
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/variable-capital-advances", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateVariableCapitalAdvance(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["id"] == nil {
		t.Error("expected id in response")
	}
	if w.Header().Get("Location") == "" {
		t.Error("expected Location header")
	}
}

func TestCreateVariableCapitalAdvance_InvalidBody(t *testing.T) {
	t.Parallel()
	h := newValHandler()

	req := httptest.NewRequest(http.MethodPost, "/v1/variable-capital-advances",
		bytes.NewBufferString(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateVariableCapitalAdvance(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetVariableCapitalAdvance_NotFound(t *testing.T) {
	t.Parallel()
	h := newValHandler()

	req := httptest.NewRequest(http.MethodGet, "/v1/variable-capital-advances/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()
	h.GetVariableCapitalAdvance(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestRecordPerCircuitSurplusRate_Created(t *testing.T) {
	t.Parallel()
	h := newValHandler()

	// First create an advance.
	createBody := `{
		"industrial_capital_id": "test-ic-002",
		"pence": 120000,
		"turnover_period_minutes": 50400,
		"turnovers_per_year_basis_points": 100000
	}`
	createReq := httptest.NewRequest(http.MethodPost, "/v1/variable-capital-advances",
		bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	cw := httptest.NewRecorder()
	h.CreateVariableCapitalAdvance(cw, createReq)
	if cw.Code != http.StatusCreated {
		t.Fatalf("create advance failed: %d", cw.Code)
	}
	var created map[string]any
	_ = json.NewDecoder(cw.Body).Decode(&created)
	id, _ := created["id"].(string)

	// Record per-circuit rate: s/v = 100% → 10000 bp.
	rateBody := `{"surplus_pence": 120000, "variable_pence": 120000}`
	rateReq := httptest.NewRequest(http.MethodPost,
		"/v1/variable-capital-advances/"+id+"/per-circuit-rate",
		bytes.NewBufferString(rateBody))
	rateReq.Header.Set("Content-Type", "application/json")
	rateReq.SetPathValue("id", id)
	rw := httptest.NewRecorder()
	h.RecordPerCircuitSurplusRate(rw, rateReq)

	if rw.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rw.Code, rw.Body.String())
	}
	var rateResp map[string]any
	_ = json.NewDecoder(rw.Body).Decode(&rateResp)
	if rateResp["basis_points"].(float64) != 10000 {
		t.Errorf("expected basis_points=10000, got %v", rateResp["basis_points"])
	}
}

func TestComputeAnnualSurplusRate_ConflictBeforePerCircuit(t *testing.T) {
	t.Parallel()
	h := newValHandler()

	// Create advance without recording per-circuit rate.
	createBody := `{
		"industrial_capital_id": "test-ic-003",
		"pence": 120000,
		"turnover_period_minutes": 50400,
		"turnovers_per_year_basis_points": 100000
	}`
	createReq := httptest.NewRequest(http.MethodPost, "/v1/variable-capital-advances",
		bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	cw := httptest.NewRecorder()
	h.CreateVariableCapitalAdvance(cw, createReq)
	var created map[string]any
	_ = json.NewDecoder(cw.Body).Decode(&created)
	id, _ := created["id"].(string)

	// Try to compute annual rate without a per-circuit rate → 409.
	annualReq := httptest.NewRequest(http.MethodPost,
		"/v1/variable-capital-advances/"+id+"/compute-annual-rate?period=1871",
		bytes.NewBufferString(`{}`))
	annualReq.Header.Set("Content-Type", "application/json")
	annualReq.SetPathValue("id", id)
	aw := httptest.NewRecorder()
	h.ComputeAnnualSurplusRateHandler(aw, annualReq)

	if aw.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", aw.Code, aw.Body.String())
	}
}

func TestListAnnualSurplusRates_NonNullEmpty(t *testing.T) {
	t.Parallel()
	h := newValHandler()

	req := httptest.NewRequest(http.MethodGet, "/v1/annual-surplus-rates", nil)
	w := httptest.NewRecorder()
	h.ListAnnualSurplusRates(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["items"] == nil {
		t.Error("items must not be null")
	}
}

func TestGetVariableCapitalReproduction_CapitalistA(t *testing.T) {
	t.Parallel()
	h := newValHandler()

	// Capitalist A: £500 (120000p), n=10, s/v=100%.
	// total_variable_paid_annual = 120000 × (100000/10000) = 1200000p (£5000).
	createBody := `{
		"industrial_capital_id": "test-ic-004",
		"pence": 120000,
		"turnover_period_minutes": 50400,
		"turnovers_per_year_basis_points": 100000
	}`
	createReq := httptest.NewRequest(http.MethodPost, "/v1/variable-capital-advances",
		bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	cw := httptest.NewRecorder()
	h.CreateVariableCapitalAdvance(cw, createReq)
	var created map[string]any
	_ = json.NewDecoder(cw.Body).Decode(&created)
	id, _ := created["id"].(string)

	// Record per-circuit rate: s/v = 100%.
	rateReq := httptest.NewRequest(http.MethodPost,
		"/v1/variable-capital-advances/"+id+"/per-circuit-rate",
		bytes.NewBufferString(`{"surplus_pence": 120000, "variable_pence": 120000}`))
	rateReq.Header.Set("Content-Type", "application/json")
	rateReq.SetPathValue("id", id)
	rw := httptest.NewRecorder()
	h.RecordPerCircuitSurplusRate(rw, rateReq)
	if rw.Code != http.StatusCreated {
		t.Fatalf("record per-circuit rate failed: %d", rw.Code)
	}

	// GET reproduction.
	reproReq := httptest.NewRequest(http.MethodGet,
		"/v1/variable-capital-reproductions/"+id, nil)
	reproReq.SetPathValue("id", id)
	rr := httptest.NewRecorder()
	h.GetVariableCapitalReproduction(rr, reproReq)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var repro map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&repro)
	if repro["total_variable_paid_annual_pence"].(float64) != 1200000 {
		t.Errorf("expected total_variable_paid_annual_pence=1200000, got %v",
			repro["total_variable_paid_annual_pence"])
	}
}
