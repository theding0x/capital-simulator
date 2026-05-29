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

func newMoneyCircuitHandler(t *testing.T) *httpapi.Handler {
	t.Helper()
	mem := store.NewMemory()
	return httpapi.New(slog.Default(), httpapi.Deps{
		MoneyCircuits: mem,
	})
}

// createSpinningMill1871MC creates a SpinningMill1871 money circuit (advance=£422, 156% rate).
func createSpinningMill1871MC(t *testing.T, h *httpapi.Handler) string {
	t.Helper()
	body := `{"amount":42200}`
	req := httptest.NewRequest(http.MethodPost, "/v1/money-circuits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateMoneyCircuit(w, req)
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

func TestCreateMoneyCircuit_Created(t *testing.T) {
	t.Parallel()
	h := newMoneyCircuitHandler(t)
	body := `{"amount":42200}`
	req := httptest.NewRequest(http.MethodPost, "/v1/money-circuits", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateMoneyCircuit(w, req)
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
	if resp["moment"] != "M" {
		t.Errorf("moment = %v, want M", resp["moment"])
	}
	advance, _ := resp["advance"].(map[string]any)
	if advance == nil {
		t.Fatal("response missing advance")
	}
	if advance["amount"] != float64(42200) {
		t.Errorf("advance.amount = %v, want 42200", advance["amount"])
	}
	loc := w.Header().Get("Location")
	if !strings.HasPrefix(loc, "/v1/money-circuits/") {
		t.Errorf("Location = %q, want /v1/money-circuits/<id>", loc)
	}
}

func TestCreateMoneyCircuit_BadRequest(t *testing.T) {
	t.Parallel()
	h := newMoneyCircuitHandler(t)

	cases := []struct {
		name string
		body string
	}{
		{"malformed json", `{not json`},
		{"zero amount", `{"amount":0}`},
		{"negative amount", `{"amount":-1}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodPost, "/v1/money-circuits", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h.CreateMoneyCircuit(w, req)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", w.Code)
			}
		})
	}
}

func TestGetMoneyCircuit_NotFound(t *testing.T) {
	t.Parallel()
	h := newMoneyCircuitHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/money-circuits/doesnotexist", nil)
	req.SetPathValue("id", "doesnotexist")
	w := httptest.NewRecorder()
	h.GetMoneyCircuit(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestGetMoneyCircuit_RoundTrip(t *testing.T) {
	t.Parallel()
	h := newMoneyCircuitHandler(t)
	id := createSpinningMill1871MC(t, h)

	req := httptest.NewRequest(http.MethodGet, "/v1/money-circuits/"+id, nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.GetMoneyCircuit(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] != id {
		t.Errorf("id = %v, want %s", resp["id"], id)
	}
}

func TestListMoneyCircuits_NonNull(t *testing.T) {
	t.Parallel()
	h := newMoneyCircuitHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/money-circuits", nil)
	w := httptest.NewRecorder()
	h.ListMoneyCircuits(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp []any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp == nil {
		t.Error("list response must not be null")
	}
}

func TestRecordMoneyPurchase_AdvancesToMC(t *testing.T) {
	t.Parallel()
	h := newMoneyCircuitHandler(t)
	id := createSpinningMill1871MC(t, h)

	// SpinningMill1871: Lp=£265.40, Mp=£156.60 → total=£422
	body := `{"labour_amount":26540,"labour_power_hours":1440,"means_amount":15660,"means_capacity_hours":1440}`
	req := httptest.NewRequest(http.MethodPost, "/v1/money-circuits/"+id+"/purchase", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.RecordMoneyPurchase(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["moment"] != "M-C" {
		t.Errorf("moment = %v, want M-C", resp["moment"])
	}
	purchase, _ := resp["purchase"].(map[string]any)
	if purchase == nil {
		t.Fatal("response missing purchase")
	}
	if purchase["total"] != float64(42200) {
		t.Errorf("purchase.total = %v, want 42200", purchase["total"])
	}
}

func TestRecordMoneyPurchase_NotFound(t *testing.T) {
	t.Parallel()
	h := newMoneyCircuitHandler(t)
	body := `{"labour_amount":26540,"labour_power_hours":1440,"means_amount":15660,"means_capacity_hours":1440}`
	req := httptest.NewRequest(http.MethodPost, "/v1/money-circuits/missing/purchase", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "missing")
	w := httptest.NewRecorder()
	h.RecordMoneyPurchase(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestRecordMoneyProduce_AdvancesToP(t *testing.T) {
	t.Parallel()
	h := newMoneyCircuitHandler(t)
	id := createSpinningMill1871MC(t, h)

	// advance to M-C first
	purchaseBody := `{"labour_amount":26540,"labour_power_hours":1440,"means_amount":15660,"means_capacity_hours":1440}`
	pr := httptest.NewRequest(http.MethodPost, "/v1/money-circuits/"+id+"/purchase", strings.NewReader(purchaseBody))
	pr.Header.Set("Content-Type", "application/json")
	pr.SetPathValue("id", id)
	h.RecordMoneyPurchase(httptest.NewRecorder(), pr)

	body := `{"constant_pence":15660,"variable_pence":26540}`
	req := httptest.NewRequest(http.MethodPost, "/v1/money-circuits/"+id+"/produce", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.RecordMoneyProduce(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["moment"] != "P" {
		t.Errorf("moment = %v, want P", resp["moment"])
	}
}
