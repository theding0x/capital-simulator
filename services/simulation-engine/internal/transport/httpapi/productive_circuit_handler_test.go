package httpapi_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

func newPCHandler(t *testing.T) *httpapi.Handler {
	t.Helper()
	mem := store.NewMemory()
	return httpapi.New(slog.Default(), httpapi.Deps{
		ProductiveCircuits: mem,
	})
}

func TestCreateProductiveCircuit_Created(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	body := `{"constant_pence":37200,"variable_pence":5000,"mode":"simple","min_capitalisation_increment_pence":100}`
	req := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateProductiveCircuit(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] == "" || resp["id"] == nil {
		t.Error("response missing id")
	}
	if resp["total_advance"] != float64(42200) {
		t.Errorf("total_advance = %v, want 42200", resp["total_advance"])
	}
	loc := w.Header().Get("Location")
	if !strings.HasPrefix(loc, "/v1/productive-circuits/") {
		t.Errorf("Location = %q, want /v1/productive-circuits/<id>", loc)
	}
}

func TestCreateProductiveCircuit_BadRequest_MissingMode(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	body := `{"constant_pence":37200,"variable_pence":5000,"min_capitalisation_increment_pence":100}`
	req := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateProductiveCircuit(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestCreateProductiveCircuit_BadRequest_InvalidJSON(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", strings.NewReader("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateProductiveCircuit(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestGetProductiveCircuit_NotFound(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/productive-circuits/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()
	h.GetProductiveCircuit(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestListProductiveCircuits_Empty(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/productive-circuits", nil)
	w := httptest.NewRecorder()
	h.ListProductiveCircuits(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp []any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp == nil {
		t.Error("response must be a non-nil slice, got null")
	}
}

func TestRecordRevenue_Created(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)

	createBody := `{"constant_pence":37200,"variable_pence":5000,"mode":"simple","min_capitalisation_increment_pence":100}`
	createReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	h.CreateProductiveCircuit(createW, createReq)
	var created map[string]any
	_ = json.NewDecoder(createW.Body).Decode(&created)
	id := created["id"].(string)

	revBody := `{"amount":7800}`
	revReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits/"+id+"/revenue", strings.NewReader(revBody))
	revReq.Header.Set("Content-Type", "application/json")
	revReq.SetPathValue("id", id)
	revW := httptest.NewRecorder()
	h.RecordRevenue(revW, revReq)
	if revW.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", revW.Code, revW.Body.String())
	}
	var rev map[string]any
	_ = json.NewDecoder(revW.Body).Decode(&rev)
	if rev["amount"] != float64(7800) {
		t.Errorf("amount = %v, want 7800", rev["amount"])
	}
}

func TestAccumulateLatentCapital_ArmsAtThreshold(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)

	createBody := `{"constant_pence":37200,"variable_pence":5000,"mode":"extended","min_capitalisation_increment_pence":100}`
	createReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	h.CreateProductiveCircuit(createW, createReq)
	var created map[string]any
	_ = json.NewDecoder(createW.Body).Decode(&created)
	id := created["id"].(string)

	accBody := `{"amount":99}`
	accReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits/"+id+"/accumulate", strings.NewReader(accBody))
	accReq.Header.Set("Content-Type", "application/json")
	accReq.SetPathValue("id", id)
	accW := httptest.NewRecorder()
	h.AccumulateLatentCapital(accW, accReq)
	var lmc map[string]any
	_ = json.NewDecoder(accW.Body).Decode(&lmc)
	if lmc["is_armed"] == true {
		t.Error("is_armed = true after 99 pence, want false (threshold=100)")
	}

	accBody2 := `{"amount":1}`
	accReq2 := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits/"+id+"/accumulate", strings.NewReader(accBody2))
	accReq2.Header.Set("Content-Type", "application/json")
	accReq2.SetPathValue("id", id)
	accW2 := httptest.NewRecorder()
	h.AccumulateLatentCapital(accW2, accReq2)
	var lmc2 map[string]any
	_ = json.NewDecoder(accW2.Body).Decode(&lmc2)
	if lmc2["is_armed"] != true {
		t.Error("is_armed = false after 100 pence, want true (threshold=100)")
	}
}

func TestCapitaliseCircuit_ErrIndivisible(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)

	createBody := `{"constant_pence":37200,"variable_pence":5000,"mode":"extended","min_capitalisation_increment_pence":100}`
	createReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	h.CreateProductiveCircuit(createW, createReq)
	var created map[string]any
	_ = json.NewDecoder(createW.Body).Decode(&created)
	id := created["id"].(string)

	capBody := `{"amount":50,"delta_constant_pence":50,"delta_variable_pence":0}`
	capReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits/"+id+"/capitalise", strings.NewReader(capBody))
	capReq.Header.Set("Content-Type", "application/json")
	capReq.SetPathValue("id", id)
	capW := httptest.NewRecorder()
	h.CapitaliseCircuit(capW, capReq)
	if capW.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422 (ErrIndivisibleProductionElement)", capW.Code)
	}
}

func TestCapitaliseCircuit_ConservationViolation(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)

	createBody := `{"constant_pence":37200,"variable_pence":5000,"mode":"extended","min_capitalisation_increment_pence":100}`
	createReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	h.CreateProductiveCircuit(createW, createReq)
	var created map[string]any
	_ = json.NewDecoder(createW.Body).Decode(&created)
	id := created["id"].(string)

	capBody := `{"amount":200,"delta_constant_pence":100,"delta_variable_pence":50}`
	capReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits/"+id+"/capitalise", strings.NewReader(capBody))
	capReq.Header.Set("Content-Type", "application/json")
	capReq.SetPathValue("id", id)
	capW := httptest.NewRecorder()
	h.CapitaliseCircuit(capW, capReq)
	if capW.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (conservation violation)", capW.Code)
	}
}

func TestDrawReserve_ErrInsufficientReserve(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)

	createBody := `{"constant_pence":37200,"variable_pence":5000,"mode":"simple","min_capitalisation_increment_pence":100}`
	createReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", strings.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	h.CreateProductiveCircuit(createW, createReq)
	var created map[string]any
	_ = json.NewDecoder(createW.Body).Decode(&created)
	id := created["id"].(string)

	drawBody := `{"amount":5000,"reason":"delayed_realisation"}`
	drawReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits/"+id+"/reserve/draw", strings.NewReader(drawBody))
	drawReq.Header.Set("Content-Type", "application/json")
	drawReq.SetPathValue("id", id)
	drawW := httptest.NewRecorder()
	h.DrawReserve(drawW, drawReq)
	if drawW.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422 (ErrInsufficientReserve)", drawW.Code)
	}
}

func TestDepositAndDrawReserve_RoundTrip(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)

	body, _ := json.Marshal(map[string]any{
		"constant_pence":                     37200,
		"variable_pence":                     5000,
		"mode":                               "simple",
		"min_capitalisation_increment_pence": 100,
	})
	createReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	h.CreateProductiveCircuit(createW, createReq)
	var created map[string]any
	_ = json.NewDecoder(createW.Body).Decode(&created)
	id := created["id"].(string)

	depBody := `{"amount":7800}`
	depReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits/"+id+"/reserve/deposit", strings.NewReader(depBody))
	depReq.Header.Set("Content-Type", "application/json")
	depReq.SetPathValue("id", id)
	depW := httptest.NewRecorder()
	h.DepositReserve(depW, depReq)
	if depW.Code != http.StatusOK {
		t.Fatalf("deposit status = %d, want 200; body: %s", depW.Code, depW.Body.String())
	}

	drawBody := `{"amount":5000,"reason":"delayed_realisation"}`
	drawReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits/"+id+"/reserve/draw", strings.NewReader(drawBody))
	drawReq.Header.Set("Content-Type", "application/json")
	drawReq.SetPathValue("id", id)
	drawW := httptest.NewRecorder()
	h.DrawReserve(drawW, drawReq)
	if drawW.Code != http.StatusCreated {
		t.Fatalf("draw status = %d, want 201; body: %s", drawW.Code, drawW.Body.String())
	}
	var draw map[string]any
	_ = json.NewDecoder(drawW.Body).Decode(&draw)
	if draw["drawn_pence"] != float64(5000) {
		t.Errorf("drawn_pence = %v, want 5000", draw["drawn_pence"])
	}
}

func TestGetProductiveCircuit_RoundTrip(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)

	body, _ := json.Marshal(map[string]any{
		"constant_pence":                     37200,
		"variable_pence":                     5000,
		"mode":                               "simple",
		"min_capitalisation_increment_pence": 100,
	})
	createReq := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", bytes.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	h.CreateProductiveCircuit(createW, createReq)
	var created map[string]any
	_ = json.NewDecoder(createW.Body).Decode(&created)
	id := created["id"].(string)

	getReq := httptest.NewRequest(http.MethodGet, "/v1/productive-circuits/"+id, nil)
	getReq.SetPathValue("id", id)
	getW := httptest.NewRecorder()
	h.GetProductiveCircuit(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", getW.Code, getW.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(getW.Body).Decode(&resp)
	if resp["id"] != id {
		t.Errorf("id = %v, want %s", resp["id"], id)
	}
	if resp["constant_pence"] != float64(37200) {
		t.Errorf("constant_pence = %v, want 37200", resp["constant_pence"])
	}
}
