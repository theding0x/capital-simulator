package httpapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	// s/v = 100%, so the year's surplus is also £5000: Marx's annual S for Capital A.
	if repro["total_surplus_annual_pence"].(float64) != 1200000 {
		t.Errorf("expected total_surplus_annual_pence=1200000, got %v",
			repro["total_surplus_annual_pence"])
	}
}

// --- Marx's Ch. 16 fixtures through the handlers (issue #450) ------------------

// valCall drives one handler and decodes its JSON object response.
func valCall(t *testing.T, fn http.HandlerFunc, method, target, id, body string, wantCode int) map[string]any {
	t.Helper()
	req := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if id != "" {
		req.SetPathValue("id", id)
	}
	w := httptest.NewRecorder()
	fn(w, req)
	if w.Code != wantCode {
		t.Fatalf("%s %s: expected %d, got %d: %s", method, target, wantCode, w.Code, w.Body.String())
	}
	var out map[string]any
	if err := json.NewDecoder(w.Body).Decode(&out); err != nil {
		t.Fatalf("%s %s: decode: %v", method, target, err)
	}
	return out
}

func valNum(t *testing.T, m map[string]any, key string) int64 {
	t.Helper()
	f, ok := m[key].(float64)
	if !ok {
		t.Fatalf("%s missing or not a number in %v", key, m)
	}
	return int64(f)
}

// valCircuit POSTs an advance and its per-circuit rate and returns the advance id.
func valCircuit(t *testing.T, h *httpapi.Handler, ic string, pence, periodMin, nBP, surplus int64) string {
	t.Helper()
	body := fmt.Sprintf(`{"industrial_capital_id": %q, "pence": %d, "turnover_period_minutes": %d, "turnovers_per_year_basis_points": %d}`,
		ic, pence, periodMin, nBP)
	created := valCall(t, h.CreateVariableCapitalAdvance, http.MethodPost, "/v1/variable-capital-advances", "", body, http.StatusCreated)
	id, _ := created["id"].(string)
	rate := fmt.Sprintf(`{"surplus_pence": %d, "variable_pence": %d}`, surplus, pence)
	valCall(t, h.RecordPerCircuitSurplusRate, http.MethodPost, "/v1/variable-capital-advances/"+id+"/per-circuit-rate", id, rate, http.StatusCreated)
	return id
}

func valAnnual(t *testing.T, h *httpapi.Handler, id string) map[string]any {
	t.Helper()
	return valCall(t, h.ComputeAnnualSurplusRateHandler, http.MethodPost,
		"/v1/variable-capital-advances/"+id+"/compute-annual-rate?period=1871", id, `{}`, http.StatusCreated)
}

func valRepro(t *testing.T, h *httpapi.Handler, id, query string) map[string]any {
	t.Helper()
	return valCall(t, h.GetVariableCapitalReproduction, http.MethodGet,
		"/v1/variable-capital-reproductions/"+id+query, id, "", http.StatusOK)
}

// Capital B (£5000, a 50-week period, n=1) reproduces the same £5000 of wages and of
// surplus a year as Capital A (£500, n=10): Ch. 16 §I's point that the annual amount
// of surplus-value is the same and only the annual rate differs.
func TestGetVariableCapitalReproduction_CapitalistB_SameAnnualSurplus(t *testing.T) {
	t.Parallel()
	h := newValHandler()
	a := valCircuit(t, h, "ic-a", 120000, 50400, 100000, 120000)
	b := valCircuit(t, h, "ic-b", 1200000, 504000, 10000, 1200000)
	reproA, reproB := valRepro(t, h, a, ""), valRepro(t, h, b, "")
	for _, key := range []string{"total_variable_paid_annual_pence", "total_surplus_annual_pence"} {
		if valNum(t, reproB, key) != 1200000 || valNum(t, reproA, key) != valNum(t, reproB, key) {
			t.Errorf("%s: A=%v B=%v, want 1200000 (£5000) each", key, reproA[key], reproB[key])
		}
	}
}

// Spinning mill: Vol. II Ch. 1's weekly £50 wages (12,000p) and £78 surplus (18,720p),
// "156 per cent", turned over weekly (the week as its period is assumed) on Ch. 16's
// 50-week year (n=50): annual rate 7,800%, £2,500 of wages and £3,900 of surplus a
// year. The year_minutes parameter is only echoed: the totals come from n.
func TestVariableCapital_SpinningMill_EndToEnd_MarxWeekly(t *testing.T) {
	t.Parallel()
	h := newValHandler()
	id := valCircuit(t, h, "spinning-mill", 12000, 10080, 500000, 18720)

	annual := valAnnual(t, h, id)
	if valNum(t, annual, "n_basis_points") != 500000 || valNum(t, annual, "per_circuit_basis_points") != 15600 ||
		valNum(t, annual, "annual_basis_points") != 780000 {
		t.Errorf("annual rate = %v, want n=500000, s/v=15600, annual=780000", annual)
	}
	repro := valRepro(t, h, id, "?year_minutes=504000")
	if valNum(t, repro, "total_variable_paid_annual_pence") != 600000 || valNum(t, repro, "total_surplus_annual_pence") != 936000 {
		t.Errorf("reproduction = %v, want wages 600000p (£2500) and surplus 936000p (£3900)", repro)
	}
	if valNum(t, repro, "year_reference_minutes") != 504000 {
		t.Errorf("year_reference_minutes = %v, want the 504000 requested", repro["year_reference_minutes"])
	}
}

// Locomotive works, derived from Ch. 12 (the spinner's week for twelve weeks): £600
// advance, 156%, n=50/12 (41,667 bp, rounded), annual 65,001 bp (650% plus 1 bp of
// rounding n). Reproduction truncates (see #449), so it lands 4p and 6p off £2,500 and
// £3,900: 600,004p and 936,006p, pinned exactly.
func TestVariableCapital_LocomotiveWorks_EndToEnd_ChapterTwelve(t *testing.T) {
	t.Parallel()
	h := newValHandler()
	id := valCircuit(t, h, "locomotive-works", 144000, 120960, 41667, 224640)

	annual := valAnnual(t, h, id)
	if valNum(t, annual, "n_basis_points") != 41667 || valNum(t, annual, "per_circuit_basis_points") != 15600 ||
		valNum(t, annual, "annual_basis_points") != 65001 {
		t.Errorf("annual rate = %v, want n=41667, s/v=15600, annual=65001", annual)
	}
	repro := valRepro(t, h, id, "")
	gotV, gotS := valNum(t, repro, "total_variable_paid_annual_pence"), valNum(t, repro, "total_surplus_annual_pence")
	if gotV != 600004 || gotS != 936006 {
		t.Errorf("reproduction v=%d s=%d, want 600004 and 936006", gotV, gotS)
	}
}

// Marx's contrast: Capital A's annual rate is 1,000% and Capital B's 100%, ten to one.
func TestCreateAnnualSurplusRateContrast_MarxAB(t *testing.T) {
	t.Parallel()
	h := newValHandler()
	a := valCircuit(t, h, "ic-ab", 120000, 50400, 100000, 120000)
	b := valCircuit(t, h, "ic-ab", 1200000, 504000, 10000, 1200000)
	valAnnual(t, h, a)
	valAnnual(t, h, b)

	body := fmt.Sprintf(`{"capital_a_id": %q, "capital_b_id": %q, "period": "1871"}`, a, b)
	got := valCall(t, h.CreateAnnualSurplusRateContrast, http.MethodPost, "/v1/annual-surplus-rate-contrasts", "", body, http.StatusCreated)
	if valNum(t, got, "capital_a_annual_basis_points") != 100000 || valNum(t, got, "capital_b_annual_basis_points") != 10000 ||
		valNum(t, got, "ratio_basis_points") != 100000 {
		t.Errorf("contrast = %v, want A=100000 (1,000%%), B=10000 (100%%), ratio=100000 (10x)", got)
	}
}
