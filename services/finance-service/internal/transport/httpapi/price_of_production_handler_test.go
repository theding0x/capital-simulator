package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// TestCreatePriceOfProduction_SphereI verifies 201 + Location header and that
// Sphere I (80c+20v) produces Price=122, Deviation=+2, Composition=higher.
func TestCreatePriceOfProduction_SphereI(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"Sphere I (80c+20v)","cost_price":100,"general_rate":2200,"commodity_value":120}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-of-production", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/avgprofit/price-of-production/") {
		t.Errorf("Location = %q, want /v1/avgprofit/price-of-production/ prefix", loc)
	}
	var pop avgprofit.PriceOfProduction
	if err := json.NewDecoder(res.Body).Decode(&pop); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if pop.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if pop.Price != 122 {
		t.Errorf("price = %d, want 122", pop.Price)
	}
	if pop.AverageProfit != 22 {
		t.Errorf("average_profit = %d, want 22", pop.AverageProfit)
	}
	if pop.Deviation != 2 {
		t.Errorf("deviation = %d, want +2", pop.Deviation)
	}
	if pop.Composition != avgprofit.KindHigher {
		t.Errorf("composition = %q, want %q", pop.Composition, avgprofit.KindHigher)
	}
}

// TestCreatePriceOfProduction_SphereIII verifies Sphere III fixture:
// Price=122, Deviation=-18, Composition=lower.
func TestCreatePriceOfProduction_SphereIII(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"Sphere III (60c+40v)","cost_price":100,"general_rate":2200,"commodity_value":140}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-of-production", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	var pop avgprofit.PriceOfProduction
	if err := json.NewDecoder(res.Body).Decode(&pop); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if pop.Deviation != -18 {
		t.Errorf("deviation = %d, want -18", pop.Deviation)
	}
	if pop.Composition != avgprofit.KindLower {
		t.Errorf("composition = %q, want %q", pop.Composition, avgprofit.KindLower)
	}
}

// TestCreatePriceOfProduction_BadJSON returns 400 on malformed request body.
func TestCreatePriceOfProduction_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-of-production", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreatePriceOfProduction_NegativeCostPrice returns 422 on negative cost_price.
func TestCreatePriceOfProduction_NegativeCostPrice(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"sphere_name":"bad","cost_price":-1,"general_rate":2200,"commodity_value":120}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-of-production", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestCreatePriceOfProduction_NegativeGeneralRate returns 422 on negative general_rate.
func TestCreatePriceOfProduction_NegativeGeneralRate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"sphere_name":"bad","cost_price":100,"general_rate":-1,"commodity_value":120}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-of-production", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestCreatePriceOfProduction_NegativeCommodityValue returns 422 on negative commodity_value.
func TestCreatePriceOfProduction_NegativeCommodityValue(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"sphere_name":"bad","cost_price":100,"general_rate":2200,"commodity_value":-5}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-of-production", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestGetPriceOfProduction_NotFound returns 404 for an unknown id.
func TestGetPriceOfProduction_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/avgprofit/price-of-production/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestListPricesOfProduction_NeverNull returns 200 with a non-nil items array even when empty.
func TestListPricesOfProduction_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/avgprofit/price-of-production")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []avgprofit.PriceOfProduction `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}

// TestCreateThenGetPriceOfProduction verifies the create→GET round-trip over HTTP.
func TestCreateThenGetPriceOfProduction(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	createBody := `{"sphere_name":"Sphere I (80c+20v)","cost_price":100,"general_rate":2200,"commodity_value":120}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-of-production", "application/json", strings.NewReader(createBody))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created avgprofit.PriceOfProduction
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/avgprofit/price-of-production/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched avgprofit.PriceOfProduction
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.Price != 122 || fetched.Deviation != 2 {
		t.Errorf("fetched = %+v, want id %s price 122 deviation +2", fetched, created.ID)
	}
}

// TestListPricesOfProduction_ContainsCreated verifies that after creating a record
// it appears in the list response.
func TestListPricesOfProduction_ContainsCreated(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	createBody := `{"sphere_name":"Sphere IV (85c+15v)","cost_price":100,"general_rate":2200,"commodity_value":115}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-of-production", "application/json", strings.NewReader(createBody))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	res.Body.Close()

	list, err := http.Get(ts.URL + "/v1/avgprofit/price-of-production")
	if err != nil {
		t.Fatalf("get list: %v", err)
	}
	defer list.Body.Close()
	var body struct {
		Items []avgprofit.PriceOfProduction `json:"items"`
	}
	if err := json.NewDecoder(list.Body).Decode(&body); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(body.Items) != 1 {
		t.Fatalf("items len = %d, want 1", len(body.Items))
	}
	if body.Items[0].Deviation != 7 {
		t.Errorf("items[0].deviation = %d, want +7 (Sphere IV)", body.Items[0].Deviation)
	}
}
