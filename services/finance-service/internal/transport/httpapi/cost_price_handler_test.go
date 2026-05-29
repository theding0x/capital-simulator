package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

func newFinanceTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/profit/cost-price", h.CreateCostPrice)
	mux.HandleFunc("GET /v1/profit/cost-price", h.ListCostPrices)
	mux.HandleFunc("GET /v1/profit/cost-price/{id}", h.GetCostPrice)
	mux.HandleFunc("POST /v1/profit/profit-form", h.ComputeProfitForm)
	mux.HandleFunc("POST /v1/profit/rate", h.CreateProfitRate)
	mux.HandleFunc("GET /v1/profit/rate", h.ListProfitRates)
	mux.HandleFunc("GET /v1/profit/rate/{id}", h.GetProfitRate)
	mux.HandleFunc("POST /v1/profit/variation", h.CreateVariation)
	mux.HandleFunc("GET /v1/profit/variation", h.ListVariations)
	mux.HandleFunc("GET /v1/profit/variation/{id}", h.GetVariation)
	mux.HandleFunc("POST /v1/profit/compare", h.CompareProfitRates)
	mux.HandleFunc("POST /v1/profit/turnover-analysis", h.CreateTurnoverAnalysis)
	mux.HandleFunc("GET /v1/profit/turnover-analysis", h.ListTurnoverAnalyses)
	mux.HandleFunc("GET /v1/profit/turnover-analysis/{id}", h.GetTurnoverAnalysis)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

func TestCreateCostPrice_Marx(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"constant":400,"variable":100,"fixed_wear_and_tear":20,"fixed_advanced":1200}`
	res, err := http.Post(ts.URL+"/v1/profit/cost-price", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/profit/cost-price/") {
		t.Errorf("Location = %q, want /v1/profit/cost-price/ prefix", loc)
	}
	var cp profit.CostPrice
	if err := json.NewDecoder(res.Body).Decode(&cp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if cp.K != 500 || cp.FixedComponent != 20 || cp.CirculatingComponent != 480 {
		t.Errorf("got k=%d fixed=%d circ=%d, want 500/20/480", cp.K, cp.FixedComponent, cp.CirculatingComponent)
	}
	if cp.ID.IsZero() {
		t.Error("expected an assigned id")
	}
}

func TestCreateCostPrice_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/cost-price", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

func TestCreateCostPrice_InvalidOutlay(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	// fixed wear-and-tear exceeds consumed constant capital → 422
	body := `{"constant":10,"variable":100,"fixed_wear_and_tear":20}`
	res, err := http.Post(ts.URL+"/v1/profit/cost-price", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestGetCostPrice_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/cost-price/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

func TestCreateThenGetCostPrice(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/cost-price", "application/json",
		strings.NewReader(`{"constant":400,"variable":100}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created profit.CostPrice
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/profit/cost-price/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched profit.CostPrice
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.K != 500 {
		t.Errorf("fetched = %+v, want id %s and k 500", fetched, created.ID)
	}
}

func TestListCostPrices_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/cost-price")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []profit.CostPrice `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}

func TestComputeProfitForm_Marx(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/profit-form", "application/json",
		strings.NewReader(`{"cost_price":500,"surplus_value":100}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got profitFormResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Profit != 100 || got.CommodityValue != 600 || !got.MystifiesOrigin {
		t.Errorf("profit-form = %+v, want profit 100, C 600, mystifies", got.ProfitForm)
	}
	if len(got.SellingPriceScenarios) != 5 {
		t.Fatalf("scenarios = %d, want 5", len(got.SellingPriceScenarios))
	}
	last := got.SellingPriceScenarios[len(got.SellingPriceScenarios)-1]
	if last.SellingPrice != 600 || last.Amount != 100 || last.BelowValue {
		t.Errorf("final scenario = %+v, want sell 600, amount 100, not below value", last)
	}
}
