package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// Spinning Mill 1871 (Vol. II Ch. 2 §I): C' = 10,000 lbs of yarn = £500,
// £422 capital-value + £78 surplus-value (legacy pence 42,200 + 7,800).
const (
	exitCommodityBody = `{"value_original":42200,"value_surplus":7800}`
	notRealisedYet    = `{"realised_pence":42200}`
	fullRealisation   = `{"realised_pence":50000,"sold_at":"1871-01-06T00:00:00Z"}`
)

func createMill(t *testing.T, h *httpapi.Handler) string {
	t.Helper()
	body := `{"constant_pence":37200,"variable_pence":5000,"mode":"simple","min_capitalisation_increment_pence":100}`
	req := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateProductiveCircuit(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status = %d; body: %s", w.Code, w.Body.String())
	}
	var created map[string]any
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	return created["id"].(string)
}

func postExitCommodity(h *httpapi.Handler, id, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits/"+id+"/commodity", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.RecordExitCommodity(w, req)
	return w
}

func postExitRealisation(h *httpapi.Handler, id, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/v1/productive-circuits/"+id+"/realise", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.RecordExitRealisation(w, req)
	return w
}

func getCircuit(t *testing.T, h *httpapi.Handler, id string) map[string]any {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/productive-circuits/"+id, nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.GetProductiveCircuit(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get status = %d; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	return resp
}

func TestRecordExitCommodity_Created(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	id := createMill(t, h)

	w := postExitCommodity(h, id, exitCommodityBody)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["productive_circuit_id"] != id {
		t.Errorf("productive_circuit_id = %v, want %s", resp["productive_circuit_id"], id)
	}
	if resp["value_original"] != float64(42200) || resp["value_surplus"] != float64(7800) {
		t.Errorf("body = %v", resp)
	}
}

func TestRecordExitCommodity_NotFound(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	if w := postExitCommodity(h, "nonexistent", exitCommodityBody); w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestRecordExitCommodity_SecondIsConflict(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	id := createMill(t, h)
	if w := postExitCommodity(h, id, exitCommodityBody); w.Code != http.StatusCreated {
		t.Fatalf("first status = %d, want 201", w.Code)
	}
	if w := postExitCommodity(h, id, exitCommodityBody); w.Code != http.StatusConflict {
		t.Fatalf("second status = %d, want 409", w.Code)
	}
}

func TestRecordExitCommodity_BadRequest(t *testing.T) {
	t.Parallel()
	tests := []struct{ name, body string }{
		{"malformed json", "{bad"},
		{"zero value_original", `{"value_original":0,"value_surplus":7800}`},
		{"negative value_original", `{"value_original":-1,"value_surplus":7800}`},
		{"negative value_surplus", `{"value_original":42200,"value_surplus":-1}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := newPCHandler(t)
			id := createMill(t, h)
			if w := postExitCommodity(h, id, tt.body); w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestRecordExitRealisation_FullSaleCreated(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	id := createMill(t, h)
	postExitCommodity(h, id, exitCommodityBody)

	w := postExitRealisation(h, id, fullRealisation)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["productive_circuit_id"] != id {
		t.Errorf("productive_circuit_id = %v, want %s", resp["productive_circuit_id"], id)
	}
	if resp["realised_pence"] != float64(50000) {
		t.Errorf("realised_pence = %v, want 50000", resp["realised_pence"])
	}
	if resp["sold_at"] != "1871-01-06T00:00:00Z" {
		t.Errorf("sold_at = %v, want 1871-01-06T00:00:00Z", resp["sold_at"])
	}
}

func TestRecordExitRealisation_PartialIs422AndNotRecorded(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	id := createMill(t, h)
	postExitCommodity(h, id, exitCommodityBody)

	if w := postExitRealisation(h, id, notRealisedYet); w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body: %s", w.Code, w.Body.String())
	}
	circuit := getCircuit(t, h, id)
	if _, present := circuit["exit_realisation"]; present {
		t.Errorf("exit_realisation present after rejected partial sale: %v", circuit["exit_realisation"])
	}
	if circuit["exit_commodity"] == nil {
		t.Error("exit_commodity missing: the circuit must stay at C'")
	}
}

func TestRecordExitRealisation_WithoutCommodityIs422(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	id := createMill(t, h)
	if w := postExitRealisation(h, id, fullRealisation); w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body: %s", w.Code, w.Body.String())
	}
}

func TestRecordExitRealisation_NotFound(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	if w := postExitRealisation(h, "nonexistent", fullRealisation); w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestRecordExitRealisation_SecondIsConflict(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	id := createMill(t, h)
	postExitCommodity(h, id, exitCommodityBody)
	if w := postExitRealisation(h, id, fullRealisation); w.Code != http.StatusCreated {
		t.Fatalf("first status = %d, want 201", w.Code)
	}
	if w := postExitRealisation(h, id, fullRealisation); w.Code != http.StatusConflict {
		t.Fatalf("second status = %d, want 409", w.Code)
	}
}

func TestRecordExitRealisation_BadRequest(t *testing.T) {
	t.Parallel()
	tests := []struct{ name, body string }{
		{"malformed json", "{bad"},
		{"zero realised_pence", `{"realised_pence":0}`},
		{"negative realised_pence", `{"realised_pence":-50000}`},
		{"bad sold_at", `{"realised_pence":50000,"sold_at":"yesterday"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := newPCHandler(t)
			id := createMill(t, h)
			postExitCommodity(h, id, exitCommodityBody)
			if w := postExitRealisation(h, id, tt.body); w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestGetProductiveCircuit_ExitFieldsAppearInOrder(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	id := createMill(t, h)

	c := getCircuit(t, h, id)
	if _, ok := c["exit_commodity"]; ok {
		t.Error("exit_commodity present before C' recorded")
	}
	if _, ok := c["exit_realisation"]; ok {
		t.Error("exit_realisation present before C' recorded")
	}

	postExitCommodity(h, id, exitCommodityBody)
	c = getCircuit(t, h, id)
	ec, ok := c["exit_commodity"].(map[string]any)
	if !ok {
		t.Fatalf("exit_commodity = %v, want object", c["exit_commodity"])
	}
	if ec["value_original"] != float64(42200) || ec["value_surplus"] != float64(7800) {
		t.Errorf("exit_commodity = %v", ec)
	}
	if _, ok := c["exit_realisation"]; ok {
		t.Error("exit_realisation present before any sale")
	}

	postExitRealisation(h, id, fullRealisation)
	c = getCircuit(t, h, id)
	er, ok := c["exit_realisation"].(map[string]any)
	if !ok {
		t.Fatalf("exit_realisation = %v, want object after full sale", c["exit_realisation"])
	}
	if er["realised_pence"] != float64(50000) {
		t.Errorf("exit_realisation = %v", er)
	}
}

func TestListProductiveCircuits_IncludesExitFieldsAndIsNeverNull(t *testing.T) {
	t.Parallel()
	h := newPCHandler(t)
	id := createMill(t, h)
	postExitCommodity(h, id, exitCommodityBody)
	postExitRealisation(h, id, fullRealisation)

	req := httptest.NewRequest(http.MethodGet, "/v1/productive-circuits", nil)
	w := httptest.NewRecorder()
	h.ListProductiveCircuits(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var list []map[string]any
	if err := json.NewDecoder(w.Body).Decode(&list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if list == nil || len(list) != 1 {
		t.Fatalf("list = %v, want one non-nil element", list)
	}
	if list[0]["exit_commodity"] == nil || list[0]["exit_realisation"] == nil {
		t.Errorf("list entry lacks exit fields: %v", list[0])
	}
}
