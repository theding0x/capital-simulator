package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
)

// --- POST /v1/merchant/commercial-profit ---

// TestCreateCommercialProfit_MarxFixture posts the 720+180 Marx fixture
// (productive=720000, commercial=180000, surplus=144000) and expects 201 +
// a Location header. The response must carry:
//   - adjusted_rate_bp (1600) < unadjusted_rate_bp (2000) — invariant 2
//   - new_value_created == 0
func TestCreateCommercialProfit_MarxFixture(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"productive_capital":720000,"commercial_capital":180000,"surplus_value":144000}`
	res, err := http.Post(ts.URL+"/v1/merchant/commercial-profit", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/merchant/commercial-profit/") {
		t.Errorf("Location = %q, want /v1/merchant/commercial-profit/ prefix", loc)
	}

	var cp merchant.CommercialProfit
	if err := json.NewDecoder(res.Body).Decode(&cp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if cp.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if cp.CreatedAt.IsZero() {
		t.Error("CreatedAt not set by store")
	}
	// Central invariant: commercial activity creates no new value.
	if cp.NewValueCreated != 0 {
		t.Errorf("new_value_created = %d, want 0", cp.NewValueCreated)
	}
	// Invariant 2: adjusted rate is diluted by commercial capital.
	if cp.AdjustedRateBP >= cp.UnadjustedRateBP {
		t.Errorf("invariant 2 violated: adjusted_rate_bp %d >= unadjusted_rate_bp %d",
			cp.AdjustedRateBP, cp.UnadjustedRateBP)
	}
	if cp.UnadjustedRateBP != 2000 {
		t.Errorf("unadjusted_rate_bp = %d, want 2000", cp.UnadjustedRateBP)
	}
	if cp.AdjustedRateBP != 1600 {
		t.Errorf("adjusted_rate_bp = %d, want 1600", cp.AdjustedRateBP)
	}
}

// TestCreateCommercialProfit_BadJSON returns 400 on a malformed body.
func TestCreateCommercialProfit_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/merchant/commercial-profit", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateCommercialProfit_NonPositiveCapital returns 422 when productive_capital
// or commercial_capital is 0 (handler maps ErrNonPositiveCapital → 422).
func TestCreateCommercialProfit_NonPositiveCapital(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	for _, body := range []string{
		`{"productive_capital":0,"commercial_capital":180000,"surplus_value":144000}`,
		`{"productive_capital":720000,"commercial_capital":0,"surplus_value":144000}`,
	} {
		res, err := http.Post(ts.URL+"/v1/merchant/commercial-profit", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("post: %v", err)
		}
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("body=%s: status = %d, want 422", body, res.StatusCode)
		}
		res.Body.Close()
	}
}

// --- GET /v1/merchant/commercial-profit ---

// TestListCommercialProfits_NeverNull asserts the list endpoint returns 200 and
// an items array that is non-null even when the store is empty.
func TestListCommercialProfits_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/merchant/commercial-profit")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		t.Fatalf("decode raw: %v", err)
	}
	items, ok := raw["items"]
	if !ok {
		t.Fatal("response has no items key")
	}
	if string(items) == "null" {
		t.Errorf("items = null, want [] for an empty store")
	}
	var decoded []merchant.CommercialProfit
	if err := json.Unmarshal(items, &decoded); err != nil {
		t.Fatalf("unmarshal items: %v", err)
	}
	if decoded == nil {
		t.Error("items decoded to nil, want non-nil slice")
	}
}

// --- GET /v1/merchant/commercial-profit/{id} ---

// TestGetCommercialProfit_NotFound returns 404 for an unknown id.
func TestGetCommercialProfit_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/merchant/commercial-profit/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateThenGetCommercialProfit exercises the create→GET round-trip using the
// 800+100 fixture (productive=800000, commercial=100000, surplus=160000 →
// adjusted=1778 bp, round-half-up).
func TestCreateThenGetCommercialProfit(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"productive_capital":800000,"commercial_capital":100000,"surplus_value":160000}`
	createRes, err := http.Post(ts.URL+"/v1/merchant/commercial-profit", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer createRes.Body.Close()

	var created merchant.CommercialProfit
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	// GET the same record.
	getRes, err := http.Get(ts.URL + "/v1/merchant/commercial-profit/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", getRes.StatusCode)
	}

	var fetched merchant.CommercialProfit
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("fetched.ID = %s, want %s", fetched.ID, created.ID)
	}
	if fetched.ProductiveCapital != 800000 {
		t.Errorf("fetched.ProductiveCapital = %d, want 800000", fetched.ProductiveCapital)
	}
	if fetched.AdjustedRateBP != 1778 {
		t.Errorf("fetched.AdjustedRateBP = %d, want 1778 (round-half-up)", fetched.AdjustedRateBP)
	}
	if fetched.NewValueCreated != 0 {
		t.Errorf("fetched.NewValueCreated = %d, want 0", fetched.NewValueCreated)
	}
}
