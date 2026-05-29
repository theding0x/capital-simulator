package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
)

func newCh3Handler() *Handler {
	return New(store.NewMemory(), nil)
}

func postJSON(t *testing.T, route string, body any, fn http.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		t.Fatalf("encode body: %v", err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, route, buf)
	r.Header.Set("Content-Type", "application/json")
	fn(w, r)
	return w
}

// TestCreateHoard_Created exercises the happy-path POST /v1/hoards using the
// Marx fixture from Ch.3 §3.a — a weaver-turned-miser withholding 200 pence
// of realised sale-money from circulation.
func TestCreateHoard_Created(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	w := postJSON(t, "/v1/hoards", map[string]any{
		"owner_id": "weaver-001",
		"amount":   200,
	}, h.CreateHoard)
	if w.Code != http.StatusCreated {
		t.Fatalf("status: want 201, got %d (%s)", w.Code, w.Body.String())
	}
	var got market.Hoard
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID.IsZero() {
		t.Fatal("expected server-assigned id")
	}
	if got.Amount != 200 {
		t.Errorf("amount: want 200, got %d", got.Amount)
	}
}

// TestCreateHoard_BadOwner verifies the validator surfaces a 400 when the
// owner id is missing.
func TestCreateHoard_BadOwner(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	w := postJSON(t, "/v1/hoards", map[string]any{"owner_id": "", "amount": 10}, h.CreateHoard)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", w.Code)
	}
}

// TestListHoards_EmptyShape verifies the response is shaped as `{items: []}`
// (not `null`) so the React client receives an iterable on fresh state.
func TestListHoards_EmptyShape(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/hoards", nil)
	h.ListHoards(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d", w.Code)
	}
	var got struct {
		Items []market.Hoard `json:"items"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Items == nil {
		t.Error("items: want empty array, got nil")
	}
}

// TestSettlePaymentObligation_NotFound exercises the ErrNotFound → 404
// mapping for an unknown obligation id.
func TestSettlePaymentObligation_NotFound(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/v1/payment-obligations/unknown/settle", nil)
	r.SetPathValue("id", "unknown")
	h.SettlePaymentObligation(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want 404, got %d", w.Code)
	}
}

// TestCreatePaymentObligation_RoundTrip creates an obligation then settles it.
// Marx Ch.3 §3.b: a linen-weaver advances cloth to a tailor on 30-day credit;
// settlement transforms the obligation back into money.
func TestCreatePaymentObligation_RoundTrip(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	due := time.Now().Add(30 * 24 * time.Hour).UTC()
	w := postJSON(t, "/v1/payment-obligations", map[string]any{
		"creditor_id": "weaver-001",
		"debtor_id":   "tailor-001",
		"amount":      4800,
		"due_at":      due,
	}, h.CreatePaymentObligation)
	if w.Code != http.StatusCreated {
		t.Fatalf("create: want 201, got %d (%s)", w.Code, w.Body.String())
	}
	var created market.PaymentObligation
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.PaidAt != nil {
		t.Error("PaidAt should be nil on creation")
	}

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "/v1/payment-obligations/"+string(created.ID)+"/settle", nil)
	r2.SetPathValue("id", string(created.ID))
	h.SettlePaymentObligation(w2, r2)
	if w2.Code != http.StatusOK {
		t.Fatalf("settle: want 200, got %d (%s)", w2.Code, w2.Body.String())
	}
	var settled market.PaymentObligation
	if err := json.Unmarshal(w2.Body.Bytes(), &settled); err != nil {
		t.Fatalf("decode settled: %v", err)
	}
	if settled.PaidAt == nil {
		t.Error("PaidAt should be set after settle")
	}
}

// TestMoneyRequired exercises Marx's quantity-of-money formula M = ΣP / V.
// Fixture from Ch.3 §2.b: if the sum of prices of commodities to be circulated
// is 2400 pence and each money-unit averages 4 turnovers per period, then
// 600 pence of currency are required.
func TestMoneyRequired(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/circulation/money-required?sum_of_prices=2400&velocity=4", nil)
	h.MoneyRequired(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status: want 200, got %d (%s)", w.Code, w.Body.String())
	}
	var got struct {
		MoneyRequired int64 `json:"money_required"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.MoneyRequired != 600 {
		t.Errorf("money_required: want 600, got %d", got.MoneyRequired)
	}
}

// TestMoneyRequired_BadVelocity verifies the velocity guard.
func TestMoneyRequired_BadVelocity(t *testing.T) {
	t.Parallel()
	h := newCh3Handler()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/circulation/money-required?sum_of_prices=100&velocity=0", nil)
	h.MoneyRequired(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want 400, got %d", w.Code)
	}
}
