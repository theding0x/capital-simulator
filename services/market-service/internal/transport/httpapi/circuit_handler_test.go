package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

// weaverCircuitBody is the Ch. 3 §2 fixture as a wire payload: the weaver sells
// 20 yards of linen for £2 (200 pence) and lays the £2 out on a family Bible.
func weaverCircuitBody() map[string]any {
	return map[string]any{
		"sale_leg": map[string]any{
			"kind": "sale", "commodity_id": "linen", "money_id": "gold", "owner_id": "weaver", "value": 200,
		},
		"purchase_leg": map[string]any{
			"kind": "purchase", "commodity_id": "bible", "money_id": "gold", "owner_id": "weaver", "value": 200,
		},
	}
}

// TestCreateCircuit_Created exercises the happy path and the equivalence
// invariant: an equivalent two-leg circuit is recorded with a server-assigned id.
func TestCreateCircuit_Created(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	w := postJSON(t, "/v1/circuits", weaverCircuitBody(), h.CreateCircuit)
	if w.Code != http.StatusCreated {
		t.Fatalf("status: want 201, got %d (%s)", w.Code, w.Body.String())
	}
	var got market.Circuit
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID.IsZero() {
		t.Fatal("expected server-assigned id")
	}
	if got.SaleLeg.Value != got.PurchaseLeg.Value {
		t.Errorf("legs not equivalent: sale %d != purchase %d", got.SaleLeg.Value, got.PurchaseLeg.Value)
	}
}

// TestCreateCircuit_NonEquivalent422 verifies that a non-equivalent circuit is
// rejected with 422 — the difference would be surplus-value (Ch. 4), not simple
// circulation.
func TestCreateCircuit_NonEquivalent422(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	body := weaverCircuitBody()
	body["purchase_leg"].(map[string]any)["value"] = 250
	w := postJSON(t, "/v1/circuits", body, h.CreateCircuit)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status: want 422, got %d (%s)", w.Code, w.Body.String())
	}
}

// TestCreateCircuit_BadJSON400 verifies malformed bodies surface a 400.
func TestCreateCircuit_BadJSON400(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/v1/circuits", strings.NewReader("{"))
	r.Header.Set("Content-Type", "application/json")
	h.CreateCircuit(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", w.Code)
	}
}

// TestListCircuits_EmptyShape verifies the response is shaped `{items: []}`
// (not null) on fresh state so the React client always receives an iterable.
func TestListCircuits_EmptyShape(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/circuits", nil)
	h.ListCircuits(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", w.Code)
	}
	var got struct {
		Items []market.Circuit `json:"items"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Items == nil {
		t.Error("items: want empty array, got nil")
	}
}

// TestCreateThenListCircuit_RoundTrip records a circuit then lists it back.
func TestCreateThenListCircuit_RoundTrip(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	if w := postJSON(t, "/v1/circuits", weaverCircuitBody(), h.CreateCircuit); w.Code != http.StatusCreated {
		t.Fatalf("create: want 201, got %d (%s)", w.Code, w.Body.String())
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/circuits", nil)
	h.ListCircuits(w, r)
	var got struct {
		Items []market.Circuit `json:"items"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got.Items) != 1 {
		t.Fatalf("want 1 circuit, got %d", len(got.Items))
	}
}
