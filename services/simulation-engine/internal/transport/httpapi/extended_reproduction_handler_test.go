package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// newExtendedHandler returns an httpapi.Handler wired to a fresh in-memory store
// with the ExtendedReproduction store set. SurplusCirculations is wired to the
// same store because the tick endpoint projects department shares onto the
// social-capital aggregate (issue #215).
func newExtendedHandler(t *testing.T) *Handler {
	t.Helper()
	mem := store.NewMemory()
	return New(nil, Deps{ExtendedReproduction: mem, SurplusCirculations: mem})
}

// createExtendedSchemeInner calls the handler and returns the created scheme ID
// and HTTP status. Useful for building test state.
func createExtendedSchemeInner(t *testing.T, h *Handler, body string) (string, int) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/schemes",
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.handleCreateExtendedReproductionScheme(w, req)
	if w.Code != http.StatusCreated {
		return "", w.Code
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode created scheme: %v", err)
	}
	id, _ := resp["id"].(string)
	return id, w.Code
}

// addExtendedDeptInner adds a department to an extended scheme. Returns the status.
func addExtendedDeptInner(t *testing.T, h *Handler, schemeID, deptJSON string) int {
	t.Helper()
	path := fmt.Sprintf("/v1/reproduction/extended/schemes/%s/departments", schemeID)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(deptJSON))
	req.SetPathValue("id", schemeID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.handleAddExtendedDepartmentToScheme(w, req)
	return w.Code
}

// --- CreateExtendedReproductionScheme ---------------------------------------

func TestHandleCreateExtendedReproductionScheme_Created(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)
	body := `{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`
	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/schemes",
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateExtendedReproductionScheme(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if id, _ := resp["id"].(string); id == "" {
		t.Fatal("expected non-empty id in response")
	}
	if resp["period"] != "Year I" {
		t.Fatalf("period: want %q, got %v", "Year I", resp["period"])
	}
	if loc := w.Header().Get("Location"); loc == "" {
		t.Fatal("expected Location header")
	}
}

func TestHandleCreateExtendedReproductionScheme_BadJSON(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)
	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/schemes",
		bytes.NewBufferString(`{bad json}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateExtendedReproductionScheme(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad JSON, got %d", w.Code)
	}
}

func TestHandleCreateExtendedReproductionScheme_MissingPeriod(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)
	body := `{"accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`
	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/schemes",
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateExtendedReproductionScheme(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing period, got %d", w.Code)
	}
}

// --- GetExtendedReproductionScheme ------------------------------------------

func TestHandleGetExtendedReproductionScheme_OK(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	// Create first.
	id, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id == "" {
		t.Fatal("could not create scheme for Get test")
	}

	req := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/schemes/"+id, nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()

	h.handleGetExtendedReproductionScheme(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] != id {
		t.Fatalf("id: want %q, got %v", id, resp["id"])
	}
}

func TestHandleGetExtendedReproductionScheme_NotFound(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	req := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/schemes/ghost-id", nil)
	req.SetPathValue("id", "ghost-id")
	w := httptest.NewRecorder()

	h.handleGetExtendedReproductionScheme(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- AddExtendedDepartmentToScheme ------------------------------------------

func TestHandleAddExtendedDepartmentToScheme_Created(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	id, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id == "" {
		t.Fatal("could not create scheme")
	}

	deptBody := `{"department":"department_i","constant_pence":4000000,"variable_pence":1000000,"surplus_pence":1000000,"total_pence":6000000}`
	code := addExtendedDeptInner(t, h, id, deptBody)
	if code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", code)
	}
}

func TestHandleAddExtendedDepartmentToScheme_BadJSON(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	id, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id == "" {
		t.Fatal("could not create scheme")
	}

	code := addExtendedDeptInner(t, h, id, `{not valid json`)
	if code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad JSON, got %d", code)
	}
}

func TestHandleAddExtendedDepartmentToScheme_AlreadyExists(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	id, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id == "" {
		t.Fatal("could not create scheme")
	}

	deptBody := `{"department":"department_i","constant_pence":4000000,"variable_pence":1000000,"surplus_pence":1000000,"total_pence":6000000}`
	if code := addExtendedDeptInner(t, h, id, deptBody); code != http.StatusCreated {
		t.Fatalf("first add: expected 201, got %d", code)
	}
	// Second add of same department must return 409.
	code := addExtendedDeptInner(t, h, id, deptBody)
	if code != http.StatusConflict {
		t.Fatalf("duplicate add: expected 409, got %d", code)
	}
}

func TestHandleAddExtendedDepartmentToScheme_SchemeNotFound(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	deptBody := `{"department":"department_i","constant_pence":4000000,"variable_pence":1000000,"surplus_pence":1000000,"total_pence":6000000}`
	code := addExtendedDeptInner(t, h, "ghost-scheme-id", deptBody)
	if code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", code)
	}
}

// --- CreateReinvestment -----------------------------------------------------

func TestHandleCreateReinvestment_Created(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	id, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id == "" {
		t.Fatal("could not create scheme")
	}

	rBody := `{"department":"dept_i","delta_constant_pence":333333,"delta_variable_pence":166667,"consumed_surplus_pence":500000,"cycle_number":1}`
	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/schemes/"+id+"/reinvestments",
		bytes.NewBufferString(rBody))
	req.SetPathValue("id", id)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateReinvestment(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := resp["id"]; !ok {
		t.Fatal("expected id in response")
	}
}

func TestHandleCreateReinvestment_BadJSON(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	id, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id == "" {
		t.Fatal("could not create scheme")
	}

	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/schemes/"+id+"/reinvestments",
		bytes.NewBufferString(`{bad json}`))
	req.SetPathValue("id", id)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateReinvestment(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleCreateReinvestment_SchemeNotFound(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	rBody := `{"department":"dept_i","delta_constant_pence":333333,"delta_variable_pence":166667,"consumed_surplus_pence":500000,"cycle_number":1}`
	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/schemes/ghost/reinvestments",
		bytes.NewBufferString(rBody))
	req.SetPathValue("id", "ghost")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateReinvestment(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- TickExtendedScheme -----------------------------------------------------

func TestHandleTickExtendedScheme_Created(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	id, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id == "" {
		t.Fatal("could not create scheme")
	}
	// Add both departments.
	dI := `{"department":"department_i","constant_pence":4000000,"variable_pence":1000000,"surplus_pence":1000000,"total_pence":6000000}`
	dII := `{"department":"department_ii","constant_pence":1500000,"variable_pence":750000,"surplus_pence":750000,"total_pence":3000000}`
	if code := addExtendedDeptInner(t, h, id, dI); code != http.StatusCreated {
		t.Fatalf("add deptI: %d", code)
	}
	if code := addExtendedDeptInner(t, h, id, dII); code != http.StatusCreated {
		t.Fatalf("add deptII: %d", code)
	}

	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/schemes/"+id+"/tick", nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()

	h.handleTickExtendedScheme(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 after tick, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	tickCount, _ := resp["tick_count"].(float64)
	if tickCount != 1 {
		t.Fatalf("tick_count: want 1, got %v", resp["tick_count"])
	}
}

func TestHandleTickExtendedScheme_NotFound(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/schemes/ghost/tick", nil)
	req.SetPathValue("id", "ghost")
	w := httptest.NewRecorder()

	h.handleTickExtendedScheme(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- GetExtendedMoneyRequirement --------------------------------------------

func TestHandleGetExtendedMoneyRequirement_OK(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	id, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id == "" {
		t.Fatal("could not create scheme")
	}
	dI := `{"department":"department_i","constant_pence":4000000,"variable_pence":1000000,"surplus_pence":1000000,"total_pence":6000000}`
	dII := `{"department":"department_ii","constant_pence":1500000,"variable_pence":750000,"surplus_pence":750000,"total_pence":3000000}`
	if code := addExtendedDeptInner(t, h, id, dI); code != http.StatusCreated {
		t.Fatalf("add deptI: %d", code)
	}
	if code := addExtendedDeptInner(t, h, id, dII); code != http.StatusCreated {
		t.Fatalf("add deptII: %d", code)
	}

	req := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/schemes/"+id+"/money-requirement?base_money_pence=100000",
		nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()

	h.handleGetExtendedMoneyRequirement(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	basePence, _ := resp["base_money_pence"].(float64)
	if basePence != 100000 {
		t.Fatalf("base_money_pence: want 100000, got %v", resp["base_money_pence"])
	}
	totalPence, _ := resp["total_money_pence"].(float64)
	if totalPence < basePence {
		t.Fatalf("total_money_pence (%v) must be >= base_money_pence (%v)",
			resp["total_money_pence"], resp["base_money_pence"])
	}
}

func TestHandleGetExtendedMoneyRequirement_NotFound(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	req := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/schemes/ghost/money-requirement?base_money_pence=100000",
		nil)
	req.SetPathValue("id", "ghost")
	w := httptest.NewRecorder()

	h.handleGetExtendedMoneyRequirement(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandleGetExtendedMoneyRequirement_BadBasePence(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	id, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id == "" {
		t.Fatal("could not create scheme")
	}

	req := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/schemes/"+id+"/money-requirement?base_money_pence=-1",
		nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()

	h.handleGetExtendedMoneyRequirement(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for negative base_money_pence, got %d", w.Code)
	}
}

// --- CreateMultiPeriodScheme ------------------------------------------------

func TestHandleCreateMultiPeriodScheme_Created(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	id1, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	id2, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year II","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id1 == "" || id2 == "" {
		t.Fatal("could not create schemes")
	}

	mpBody := fmt.Sprintf(`{"label":"Marx Year I → Year II","scheme_ids":["%s","%s"]}`, id1, id2)
	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/multi-period",
		bytes.NewBufferString(mpBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateMultiPeriodScheme(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := resp["id"]; !ok {
		t.Fatal("expected id in multi-period response")
	}
	if loc := w.Header().Get("Location"); loc == "" {
		t.Fatal("expected Location header")
	}
}

func TestHandleCreateMultiPeriodScheme_BadJSON(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/multi-period",
		bytes.NewBufferString(`{bad json}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateMultiPeriodScheme(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleCreateMultiPeriodScheme_MissingLabel(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/multi-period",
		bytes.NewBufferString(`{"label":"","scheme_ids":[]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateMultiPeriodScheme(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing label, got %d", w.Code)
	}
}

func TestHandleCreateMultiPeriodScheme_SchemeNotFound(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	mpBody := `{"label":"bad link","scheme_ids":["ghost-scheme-id"]}`
	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/multi-period",
		bytes.NewBufferString(mpBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleCreateMultiPeriodScheme(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for ghost scheme, got %d", w.Code)
	}
}

// --- GetCompositionShift ----------------------------------------------------

func TestHandleGetCompositionShift_OK(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	// Build two extended schemes and a multi-period linking them.
	id1, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	id2, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year II","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id1 == "" || id2 == "" {
		t.Fatal("could not create schemes")
	}

	mpBody := fmt.Sprintf(`{"label":"two periods","scheme_ids":["%s","%s"]}`, id1, id2)
	mpReq := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/multi-period",
		bytes.NewBufferString(mpBody))
	mpReq.Header.Set("Content-Type", "application/json")
	mpW := httptest.NewRecorder()
	h.handleCreateMultiPeriodScheme(mpW, mpReq)
	if mpW.Code != http.StatusCreated {
		t.Fatalf("create multi-period: %d: %s", mpW.Code, mpW.Body.String())
	}
	var mpResp map[string]any
	if err := json.NewDecoder(mpW.Body).Decode(&mpResp); err != nil {
		t.Fatalf("decode mp: %v", err)
	}
	mpID, _ := mpResp["id"].(string)

	req := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/multi-period/"+mpID+"/composition-shift",
		nil)
	req.SetPathValue("id", mpID)
	w := httptest.NewRecorder()

	h.handleGetCompositionShift(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp []any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode shifts: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-null array for composition shifts")
	}
}

func TestHandleGetCompositionShift_NotFound(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	req := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/multi-period/ghost/composition-shift",
		nil)
	req.SetPathValue("id", "ghost")
	w := httptest.NewRecorder()

	h.handleGetCompositionShift(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- GetGrowthLead ----------------------------------------------------------

func TestHandleGetGrowthLead_OK(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	id1, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	id2, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year II","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id1 == "" || id2 == "" {
		t.Fatal("could not create schemes")
	}

	mpBody := fmt.Sprintf(`{"label":"two periods","scheme_ids":["%s","%s"]}`, id1, id2)
	mpReq := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/multi-period",
		bytes.NewBufferString(mpBody))
	mpReq.Header.Set("Content-Type", "application/json")
	mpW := httptest.NewRecorder()
	h.handleCreateMultiPeriodScheme(mpW, mpReq)
	var mpResp map[string]any
	if err := json.NewDecoder(mpW.Body).Decode(&mpResp); err != nil {
		t.Fatalf("decode mp: %v", err)
	}
	mpID, _ := mpResp["id"].(string)

	req := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/multi-period/"+mpID+"/growth-lead",
		nil)
	req.SetPathValue("id", mpID)
	w := httptest.NewRecorder()

	h.handleGetGrowthLead(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp []any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode growth leads: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-null array for growth leads")
	}
}

func TestHandleGetGrowthLead_NotFound(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	req := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/multi-period/ghost/growth-lead",
		nil)
	req.SetPathValue("id", "ghost")
	w := httptest.NewRecorder()

	h.handleGetGrowthLead(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// --- List endpoints never return null ---------------------------------------

func TestHandleListEndpoints_NeverReturnNull(t *testing.T) {
	t.Parallel()
	h := newExtendedHandler(t)

	// Create a multi-period scheme with a single extended scheme.
	id1, _ := createExtendedSchemeInner(t, h,
		`{"period":"Year I","accumulation_rate_i_bps":5000,"accumulation_rate_ii_bps":3000}`)
	if id1 == "" {
		t.Fatal("could not create scheme")
	}

	mpBody := fmt.Sprintf(`{"label":"single","scheme_ids":["%s"]}`, id1)
	mpReq := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/extended/multi-period",
		bytes.NewBufferString(mpBody))
	mpReq.Header.Set("Content-Type", "application/json")
	mpW := httptest.NewRecorder()
	h.handleCreateMultiPeriodScheme(mpW, mpReq)
	if mpW.Code != http.StatusCreated {
		t.Fatalf("create multi-period: %d: %s", mpW.Code, mpW.Body.String())
	}
	var mpResp map[string]any
	if err := json.NewDecoder(mpW.Body).Decode(&mpResp); err != nil {
		t.Fatalf("decode mp: %v", err)
	}
	mpID, _ := mpResp["id"].(string)

	// /composition-shift must return [] not null.
	csReq := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/multi-period/"+mpID+"/composition-shift", nil)
	csReq.SetPathValue("id", mpID)
	csW := httptest.NewRecorder()
	h.handleGetCompositionShift(csW, csReq)
	if csW.Code != http.StatusOK {
		t.Fatalf("composition-shift: expected 200, got %d", csW.Code)
	}
	if body := csW.Body.String(); body == "null\n" || body == "null" {
		t.Fatal("composition-shift must not return null — got null")
	}

	// /growth-lead must return [] not null.
	glReq := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/extended/multi-period/"+mpID+"/growth-lead", nil)
	glReq.SetPathValue("id", mpID)
	glW := httptest.NewRecorder()
	h.handleGetGrowthLead(glW, glReq)
	if glW.Code != http.StatusOK {
		t.Fatalf("growth-lead: expected 200, got %d", glW.Code)
	}
	if body := glW.Body.String(); body == "null\n" || body == "null" {
		t.Fatal("growth-lead must not return null — got null")
	}
}
