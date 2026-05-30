package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
)

// --- POST /v1/merchant/turnover ---

// TestCreateMerchantTurnover_1TurnoverFixture posts the 1-turnover sugar-merchant
// fixture (Vol. III Ch. 18): moneyAdvanced=100000p, generalRateBP=1500,
// turnoverCount=1, unitsPerTurnover=300.
//
//	annual = 15000; markup = 50
//
// Expects 201 + Location header with the assigned ID.
func TestCreateMerchantTurnover_1TurnoverFixture(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":100000,"general_rate_bp":1500,"turnover_count":1,"units_per_turnover":300}`
	res, err := http.Post(ts.URL+"/v1/merchant/turnover", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/merchant/turnover/") {
		t.Errorf("Location = %q, want /v1/merchant/turnover/ prefix", loc)
	}

	var mt merchant.MerchantTurnover
	if err := json.NewDecoder(res.Body).Decode(&mt); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if mt.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if mt.CreatedAt.IsZero() {
		t.Error("CreatedAt not set by store")
	}
	if mt.AnnualMerchantProfit != 15000 {
		t.Errorf("annual_merchant_profit = %d, want 15000", mt.AnnualMerchantProfit)
	}
	if mt.MarkupPerUnit != 50 {
		t.Errorf("markup_per_unit = %d, want 50", mt.MarkupPerUnit)
	}
}

// TestCreateMerchantTurnover_BadJSON returns 400 on a malformed body.
func TestCreateMerchantTurnover_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/merchant/turnover", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateMerchantTurnover_InvalidTurnoverCount returns 422 when turnover_count
// is 0 (handler maps ErrInvalidTurnoverCount → 422).
func TestCreateMerchantTurnover_InvalidTurnoverCount(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":100000,"general_rate_bp":1500,"turnover_count":0,"units_per_turnover":300}`
	res, err := http.Post(ts.URL+"/v1/merchant/turnover", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// --- GET /v1/merchant/turnover ---

// TestListMerchantTurnovers_NeverNull asserts the list endpoint returns 200 and
// an items array that is non-null even when the store is empty.
func TestListMerchantTurnovers_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/merchant/turnover")
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
	var decoded []merchant.MerchantTurnover
	if err := json.Unmarshal(items, &decoded); err != nil {
		t.Fatalf("unmarshal items: %v", err)
	}
	if decoded == nil {
		t.Error("items decoded to nil, want non-nil slice")
	}
}

// --- GET /v1/merchant/turnover/{id} ---

// TestGetMerchantTurnover_NotFound returns 404 for an unknown id.
func TestGetMerchantTurnover_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/merchant/turnover/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// --- POST /v1/merchant/turnover-effect ---

// TestComputeTurnoverEffect_MarxFixture posts the [1,5]-turnover analysis
// (moneyAdvanced=100000, generalRateBP=1500, unitsPerTurnover=300,
// turnoverCounts=[1,5]) and expects 200 with:
//
//	points[0].markup_per_unit == 50
//	points[1].markup_per_unit == 10
//
// Also confirms the endpoint does NOT persist: a subsequent GET /v1/merchant/turnover
// list is empty (items == []).
func TestComputeTurnoverEffect_MarxFixture(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":100000,"general_rate_bp":1500,"units_per_turnover":300,"turnover_counts":[1,5]}`
	res, err := http.Post(ts.URL+"/v1/merchant/turnover-effect", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}

	var analysis merchant.TurnoverEffectAnalysis
	if err := json.NewDecoder(res.Body).Decode(&analysis); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if analysis.AnnualMerchantProfit != 15000 {
		t.Errorf("annual_merchant_profit = %d, want 15000", analysis.AnnualMerchantProfit)
	}
	if len(analysis.Points) != 2 {
		t.Fatalf("len(points) = %d, want 2", len(analysis.Points))
	}
	if analysis.Points[0].MarkupPerUnit != 50 {
		t.Errorf("points[0].markup_per_unit = %d, want 50", analysis.Points[0].MarkupPerUnit)
	}
	if analysis.Points[1].MarkupPerUnit != 10 {
		t.Errorf("points[1].markup_per_unit = %d, want 10", analysis.Points[1].MarkupPerUnit)
	}

	// Confirm the pure-compute endpoint did NOT persist anything.
	listRes, err := http.Get(ts.URL + "/v1/merchant/turnover")
	if err != nil {
		t.Fatalf("list get: %v", err)
	}
	defer listRes.Body.Close()

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(listRes.Body).Decode(&raw); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	var items []merchant.MerchantTurnover
	if err := json.Unmarshal(raw["items"], &items); err != nil {
		t.Fatalf("unmarshal items: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("after turnover-effect, list len = %d, want 0 (not persisted)", len(items))
	}
}

// TestComputeTurnoverEffect_BadJSON returns 400 on a malformed body.
func TestComputeTurnoverEffect_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/merchant/turnover-effect", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestComputeTurnoverEffect_EmptyTurnoverCounts returns 422 when turnover_counts
// is empty (handler maps ErrInvalidTurnoverCount → 422).
func TestComputeTurnoverEffect_EmptyTurnoverCounts(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":100000,"general_rate_bp":1500,"units_per_turnover":300,"turnover_counts":[]}`
	res, err := http.Post(ts.URL+"/v1/merchant/turnover-effect", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}
