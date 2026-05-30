package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// TestCreateMarketValue_CottonFixture verifies 201 + Location and the spec fixture:
// bulk=100, best=80, worst=130 → Value=100.
func TestCreateMarketValue_CottonFixture(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"cotton","bulk_condition_value":100,"best_condition_value":80,"worst_condition_value":130}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/market-value", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/avgprofit/market-value/") {
		t.Errorf("Location = %q, want /v1/avgprofit/market-value/ prefix", loc)
	}

	var mv avgprofit.MarketValue
	if err := json.NewDecoder(res.Body).Decode(&mv); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if mv.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if mv.Value != 100 {
		t.Errorf("Value = %d, want 100 (bulk condition)", mv.Value)
	}
	if mv.BestConditionValue != 80 {
		t.Errorf("BestConditionValue = %d, want 80", mv.BestConditionValue)
	}
	if mv.WorstConditionValue != 130 {
		t.Errorf("WorstConditionValue = %d, want 130", mv.WorstConditionValue)
	}
}

// TestCreateMarketValue_BadJSON returns 400 on malformed body.
func TestCreateMarketValue_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/avgprofit/market-value", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateMarketValue_NegativeBulk returns 422 on negative bulk_condition_value.
func TestCreateMarketValue_NegativeBulk(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"cotton","bulk_condition_value":-1,"best_condition_value":80,"worst_condition_value":130}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/market-value", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestCreateMarketValue_NegativeBest returns 422 on negative best_condition_value.
func TestCreateMarketValue_NegativeBest(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"cotton","bulk_condition_value":100,"best_condition_value":-1,"worst_condition_value":130}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/market-value", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestGetMarketValue_NotFound returns 404 for an unknown id.
func TestGetMarketValue_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/avgprofit/market-value/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateThenGetMarketValue verifies the create → GET round-trip over HTTP.
func TestCreateThenGetMarketValue(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"cotton","bulk_condition_value":100,"best_condition_value":80,"worst_condition_value":130}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/market-value", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created avgprofit.MarketValue
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/avgprofit/market-value/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched avgprofit.MarketValue
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.Value != 100 {
		t.Errorf("fetched = %+v, want id %s value 100", fetched, created.ID)
	}
}

// TestCreateSurplusProfit_HighProductivityFirm verifies 201 + Location and the spec
// fixture: indVal=60, marketVal=100, qty=1000, generalRate=2200 → amount=40000.
func TestCreateSurplusProfit_HighProductivityFirm(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"firm_name":"High-productivity firm","individual_value":60,"market_value":100,"output_qty":1000,"general_rate":2200}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/surplus-profit", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/avgprofit/surplus-profit/") {
		t.Errorf("Location = %q, want /v1/avgprofit/surplus-profit/ prefix", loc)
	}

	var sp avgprofit.SurplusProfit
	if err := json.NewDecoder(res.Body).Decode(&sp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sp.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if sp.Amount != 40000 {
		t.Errorf("Amount = %d, want 40000 ((100-60)*1000)", sp.Amount)
	}
}

// TestCreateSurplusProfit_BadJSON returns 400 on malformed body.
func TestCreateSurplusProfit_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/avgprofit/surplus-profit", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateSurplusProfit_NegativeIndividualValue returns 422 on negative individual_value.
func TestCreateSurplusProfit_NegativeIndividualValue(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"firm_name":"bad","individual_value":-1,"market_value":100,"output_qty":1000,"general_rate":2200}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/surplus-profit", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestGetSurplusProfit_NotFound returns 404 for an unknown id.
func TestGetSurplusProfit_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/avgprofit/surplus-profit/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}
