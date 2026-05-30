package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// TestCreateGeneralProfitRate_FiveSpheres verifies 201 + Location header and that
// the canonical five-sphere fixture produces Rate=2200.
func TestCreateGeneralProfitRate_FiveSpheres(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"spheres":[
		{"name":"Sphere I (80c+20v)","c":100,"v":20,"s_rate":10000},
		{"name":"Sphere II (70c+30v)","c":100,"v":30,"s_rate":10000},
		{"name":"Sphere III (60c+40v)","c":100,"v":40,"s_rate":10000},
		{"name":"Sphere IV (85c+15v)","c":100,"v":15,"s_rate":10000},
		{"name":"Sphere V (95c+5v)","c":100,"v":5,"s_rate":10000}
	]}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/general-rate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/avgprofit/general-rate/") {
		t.Errorf("Location = %q, want /v1/avgprofit/general-rate/ prefix", loc)
	}
	var g avgprofit.GeneralProfitRate
	if err := json.NewDecoder(res.Body).Decode(&g); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if g.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if g.Rate != 2200 {
		t.Errorf("rate = %d, want 2200", g.Rate)
	}
	if g.SumSurplusValues != 110 {
		t.Errorf("sum_surplus_values = %d, want 110", g.SumSurplusValues)
	}
	if g.SumTotalCapitals != 500 {
		t.Errorf("sum_total_capitals = %d, want 500", g.SumTotalCapitals)
	}
}

// TestCreateGeneralProfitRate_BadJSON returns 400 on malformed request body.
func TestCreateGeneralProfitRate_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/avgprofit/general-rate", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateGeneralProfitRate_NegativeC returns 422 when c is negative.
func TestCreateGeneralProfitRate_NegativeC(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"spheres":[{"name":"bad","c":-100,"v":50,"s_rate":10000}]}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/general-rate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestCreateGeneralProfitRate_VExceedsC returns 422 when v > c.
func TestCreateGeneralProfitRate_VExceedsC(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"spheres":[{"name":"bad","c":100,"v":200,"s_rate":10000}]}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/general-rate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestGetGeneralProfitRate_NotFound returns 404 for an unknown id.
func TestGetGeneralProfitRate_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/avgprofit/general-rate/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateThenGetGeneralProfitRate verifies the create→GET round-trip over HTTP.
func TestCreateThenGetGeneralProfitRate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"spheres":[{"name":"Sphere I (80c+20v)","c":100,"v":20,"s_rate":10000}]}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/general-rate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	var created avgprofit.GeneralProfitRate
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/avgprofit/general-rate/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched avgprofit.GeneralProfitRate
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.Rate != created.Rate {
		t.Errorf("fetched = %+v, want id %s rate %d", fetched, created.ID, created.Rate)
	}
}

// TestComputeSocialAggregate_FiveSpheres verifies that the social-aggregate endpoint
// returns 200 with ValueConserved=true for the canonical five spheres.
func TestComputeSocialAggregate_FiveSpheres(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"spheres":[
		{"name":"Sphere I (80c+20v)","c":100,"v":20,"s_rate":10000},
		{"name":"Sphere II (70c+30v)","c":100,"v":30,"s_rate":10000},
		{"name":"Sphere III (60c+40v)","c":100,"v":40,"s_rate":10000},
		{"name":"Sphere IV (85c+15v)","c":100,"v":15,"s_rate":10000},
		{"name":"Sphere V (95c+5v)","c":100,"v":5,"s_rate":10000}
	]}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/social-aggregate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var agg avgprofit.SocialCapitalAggregate
	if err := json.NewDecoder(res.Body).Decode(&agg); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !agg.ValueConserved {
		t.Error("value_conserved = false, want true for five-sphere fixture")
	}
	if agg.WeightedGeneralRate != 2200 {
		t.Errorf("weighted_general_rate = %d, want 2200", agg.WeightedGeneralRate)
	}
	if agg.TotalValues != 610 {
		t.Errorf("total_values = %d, want 610", agg.TotalValues)
	}
	if agg.TotalPricesOfProduction != 610 {
		t.Errorf("total_prices_of_production = %d, want 610", agg.TotalPricesOfProduction)
	}
	if agg.SumAverageProfits != 110 {
		t.Errorf("sum_average_profits = %d, want 110", agg.SumAverageProfits)
	}
}

// TestComputeSocialAggregate_BadJSON returns 400 on malformed request body.
func TestComputeSocialAggregate_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/avgprofit/social-aggregate", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestComputeSocialAggregate_VExceedsC returns 422 when v > c in any sphere.
func TestComputeSocialAggregate_VExceedsC(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"spheres":[{"name":"bad","c":50,"v":100,"s_rate":10000}]}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/social-aggregate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}
