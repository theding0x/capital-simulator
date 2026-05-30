package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
)

// --- POST /v1/merchant/historical-forms ---

// TestCreateHistoricalMerchantCapital_Venice posts the Venice seed fixture
// (Vol. III Ch. 20, seed 2001):
//
//	name="Venice carrying trade", stage=1 (StagePreCapitalist),
//	profit_source="unequal_exchange", wage_labour=false, subordination_index=500
//
// Expects 201 + Location header with the assigned ID.
func TestCreateHistoricalMerchantCapital_Venice(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Venice carrying trade","stage":1,"profit_source":"unequal_exchange","wage_labour":false,"subordination_index":500}`
	res, err := http.Post(ts.URL+"/v1/merchant/historical-forms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/merchant/historical-forms/") {
		t.Errorf("Location = %q, want /v1/merchant/historical-forms/ prefix", loc)
	}

	var hm merchant.HistoricalMerchantCapital
	if err := json.NewDecoder(res.Body).Decode(&hm); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if hm.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if hm.CreatedAt.IsZero() {
		t.Error("CreatedAt not set by store")
	}
	if hm.Stage != merchant.StagePreCapitalist {
		t.Errorf("stage = %d, want %d (StagePreCapitalist)", hm.Stage, merchant.StagePreCapitalist)
	}
	if hm.SubordinationIndex != 500 {
		t.Errorf("subordination_index = %d, want 500", hm.SubordinationIndex)
	}
	if hm.WageLabour {
		t.Error("wage_labour = true, want false for pre-capitalist")
	}
	if hm.ProfitSource != "unequal_exchange" {
		t.Errorf("profit_source = %q, want unequal_exchange", hm.ProfitSource)
	}
}

// TestCreateHistoricalMerchantCapital_BadJSON returns 400 on a malformed body.
func TestCreateHistoricalMerchantCapital_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/merchant/historical-forms", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateHistoricalMerchantCapital_InvalidStageZero returns 422 when
// stage=0 (handler maps ErrInvalidDevelopmentStage → 422).
func TestCreateHistoricalMerchantCapital_InvalidStageZero(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","stage":0,"profit_source":"unequal_exchange","wage_labour":false,"subordination_index":500}`
	res, err := http.Post(ts.URL+"/v1/merchant/historical-forms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for stage=0", res.StatusCode)
	}
}

// TestCreateHistoricalMerchantCapital_SubordinationIndexTooHigh returns 422
// when subordination_index=10001 (ErrSubordinationOutOfRange → 422).
func TestCreateHistoricalMerchantCapital_SubordinationIndexTooHigh(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","stage":2,"profit_source":"carrying_monopoly","wage_labour":false,"subordination_index":10001}`
	res, err := http.Post(ts.URL+"/v1/merchant/historical-forms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for subordination_index=10001", res.StatusCode)
	}
}

// TestCreateHistoricalMerchantCapital_PreCapitalistWageLabour returns 422 when
// stage=1 and wage_labour=true (ErrWageLabourInPreCapitalist → 422).
func TestCreateHistoricalMerchantCapital_PreCapitalistWageLabour(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","stage":1,"profit_source":"unequal_exchange","wage_labour":true,"subordination_index":500}`
	res, err := http.Post(ts.URL+"/v1/merchant/historical-forms", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for pre-capitalist+wage_labour=true", res.StatusCode)
	}
}

// --- GET /v1/merchant/historical-forms ---

// TestListHistoricalMerchantCapitals_NeverNull asserts the list endpoint
// returns 200 and an items array that is non-null even when the store is empty.
func TestListHistoricalMerchantCapitals_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/merchant/historical-forms")
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
	var decoded []merchant.HistoricalMerchantCapital
	if err := json.Unmarshal(items, &decoded); err != nil {
		t.Fatalf("unmarshal items: %v", err)
	}
	if decoded == nil {
		t.Error("items decoded to nil, want non-nil slice")
	}
}

// --- GET /v1/merchant/historical-forms/{id} ---

// TestGetHistoricalMerchantCapital_NotFound returns 404 for an unknown id.
func TestGetHistoricalMerchantCapital_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/merchant/historical-forms/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}
