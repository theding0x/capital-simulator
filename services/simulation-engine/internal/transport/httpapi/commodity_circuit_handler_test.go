package httpapi_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

func newCCHandler(t *testing.T) *httpapi.Handler {
	t.Helper()
	mem := store.NewMemory()
	return httpapi.New(slog.Default(), httpapi.Deps{
		CommodityCircuits: mem,
	})
}

// createSpinningMill1871 is a helper that POSTs the Spinning Mill 1871 fixture
// and returns the created circuit's id.
func createSpinningMill1871(t *testing.T, h *httpapi.Handler) string {
	t.Helper()
	body := `{"constant_pence":37200,"variable_pence":5000,"surplus_pence":7800,"pounds_total":10000,"mode":"simple"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/commodity-circuits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateCommodityCircuit(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	id, _ := resp["id"].(string)
	if id == "" {
		t.Fatal("create response missing id")
	}
	return id
}

// TestCreateCommodityCircuit_Created verifies the Spinning Mill 1871 fixture is
// accepted and the response contains computed fields.
func TestCreateCommodityCircuit_Created(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	body := `{"constant_pence":37200,"variable_pence":5000,"surplus_pence":7800,"pounds_total":10000,"mode":"simple"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/commodity-circuits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateCommodityCircuit(w, req)
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
	initial, _ := resp["initial"].(map[string]any)
	if initial == nil {
		t.Fatal("response missing initial")
	}
	if initial["total"] != float64(50000) {
		t.Errorf("initial.total = %v, want 50000 (£500)", initial["total"])
	}
	if initial["value_original_pence"] != float64(42200) {
		t.Errorf("initial.value_original_pence = %v, want 42200 (£422)", initial["value_original_pence"])
	}
	loc := w.Header().Get("Location")
	if !strings.HasPrefix(loc, "/v1/commodity-circuits/") {
		t.Errorf("Location = %q, want /v1/commodity-circuits/<id>", loc)
	}
}

// TestCreateCommodityCircuit_NotCommodityCapital verifies that a zero
// surplus_pence without is_first_investment returns 400.
func TestCreateCommodityCircuit_NotCommodityCapital(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	body := `{"constant_pence":37200,"variable_pence":5000,"surplus_pence":0,"pounds_total":10000,"mode":"simple"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/commodity-circuits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateCommodityCircuit(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (ErrNotCommodityCapital)", w.Code)
	}
}

// TestCreateCommodityCircuit_FirstInvestment verifies that is_first_investment=true
// bypasses the surplus_pence check.
func TestCreateCommodityCircuit_FirstInvestment(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	body := `{"constant_pence":37200,"variable_pence":5000,"surplus_pence":0,"pounds_total":10000,"mode":"simple","is_first_investment":true}`
	req := httptest.NewRequest(http.MethodPost, "/v1/commodity-circuits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateCommodityCircuit(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201 for first investment; body: %s", w.Code, w.Body.String())
	}
}

// TestCreateCommodityCircuit_BadJSON verifies 400 on malformed body.
func TestCreateCommodityCircuit_BadJSON(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/v1/commodity-circuits", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateCommodityCircuit(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

// TestGetCommodityCircuit_RoundTrip verifies that a created circuit is
// retrievable with its original magnitudes intact.
func TestGetCommodityCircuit_RoundTrip(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	id := createSpinningMill1871(t, h)

	req := httptest.NewRequest(http.MethodGet, "/v1/commodity-circuits/"+id, nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.GetCommodityCircuit(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["id"] != id {
		t.Errorf("id = %v, want %s", resp["id"], id)
	}
}

// TestGetCommodityCircuit_NotFound returns 404 for a missing id.
func TestGetCommodityCircuit_NotFound(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/commodity-circuits/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()
	h.GetCommodityCircuit(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

// TestListCommodityCircuits_Empty verifies the list endpoint returns a non-nil
// slice (never JSON null) when no circuits exist.
func TestListCommodityCircuits_Empty(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/commodity-circuits", nil)
	w := httptest.NewRecorder()
	h.ListCommodityCircuits(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp []any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp == nil {
		t.Error("list response must be non-nil slice, got null")
	}
}

// TestRecordPartialSale_DecompositionCorrect verifies that the server computes
// the aliquot decomposition for the first tranche of the Spinning Mill fixture:
// 7,440 lbs realised at £372 → decomposition must sum to 37200 pence.
func TestRecordPartialSale_DecompositionCorrect(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	id := createSpinningMill1871(t, h)

	saleBody := `{"quantity":7440,"realised_pence":37200}`
	req := httptest.NewRequest(http.MethodPost, "/v1/commodity-circuits/"+id+"/partial-sales", strings.NewReader(saleBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.RecordCommodityPartialSale(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", w.Code, w.Body.String())
	}
	var sale map[string]any
	_ = json.NewDecoder(w.Body).Decode(&sale)
	decomp, _ := sale["decomposition"].(map[string]any)
	if decomp == nil {
		t.Fatal("response missing decomposition")
	}
	if decomp["total"] != float64(37200) {
		t.Errorf("decomposition.total = %v, want 37200", decomp["total"])
	}
}

// TestRecordPartialSale_NotFound returns 404 for a missing circuit.
func TestRecordPartialSale_NotFound(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	body := `{"quantity":100,"realised_pence":500}`
	req := httptest.NewRequest(http.MethodPost, "/v1/commodity-circuits/ghost/partial-sales", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "ghost")
	w := httptest.NewRecorder()
	h.RecordCommodityPartialSale(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

// TestLinkMPSource_Created verifies an import source is accepted.
func TestLinkMPSource_Created(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	id := createSpinningMill1871(t, h)

	srcBody := `{"source_kind":"import"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/commodity-circuits/"+id+"/mp-source", strings.NewReader(srcBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.LinkCommodityMPSource(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", w.Code, w.Body.String())
	}
	var src map[string]any
	_ = json.NewDecoder(w.Body).Decode(&src)
	if src["source_kind"] != "import" {
		t.Errorf("source_kind = %v, want import", src["source_kind"])
	}
}

// TestLinkMPSource_InvalidKind returns 400 for an unrecognised source_kind.
func TestLinkMPSource_InvalidKind(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	id := createSpinningMill1871(t, h)

	srcBody := `{"source_kind":"fantasy_source"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/commodity-circuits/"+id+"/mp-source", strings.NewReader(srcBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.LinkCommodityMPSource(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (ErrUnsourcedMeansOfProduction)", w.Code)
	}
}

// TestCloseCommodityCircuit_SimpleReproduction closes the Spinning Mill 1871
// circuit with the same opening value — IsExtended must be false.
func TestCloseCommodityCircuit_SimpleReproduction(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	id := createSpinningMill1871(t, h)

	closeBody := `{"constant_pence":37200,"variable_pence":5000,"surplus_pence":7800,"capitalised_pence":0,"pounds_total":10000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/commodity-circuits/"+id+"/close", strings.NewReader(closeBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.CloseCommodityCircuit(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	terminal, _ := resp["terminal"].(map[string]any)
	if terminal == nil {
		t.Fatal("response missing terminal")
	}
	if terminal["is_extended"] != false {
		t.Errorf("terminal.is_extended = %v, want false (simple reproduction)", terminal["is_extended"])
	}
}

// TestGetCommodityCircuitAggregate_Empty returns zero counts on a fresh store.
func TestGetCommodityCircuitAggregate_Empty(t *testing.T) {
	t.Parallel()
	h := newCCHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/commodity-circuits/aggregate", nil)
	w := httptest.NewRecorder()
	h.GetCommodityCircuitAggregate(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["circuit_count"] != float64(0) {
		t.Errorf("circuit_count = %v, want 0", resp["circuit_count"])
	}
}
