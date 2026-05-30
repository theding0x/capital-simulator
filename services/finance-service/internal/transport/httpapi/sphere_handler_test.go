package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// Sphere A (600c+100v, s'=100%): create returns 201 + Location header and
// the correct derived fields — IndividualProfitRate must be 1429 (round-half-up).
func TestCreateProductionSphere_SphereA(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"name":"Sphere A (600c+100v)","c":700,"v":100,"s_rate":10000}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/spheres", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/avgprofit/spheres/") {
		t.Errorf("Location = %q, want /v1/avgprofit/spheres/ prefix", loc)
	}
	var sp avgprofit.ProductionSphere
	if err := json.NewDecoder(res.Body).Decode(&sp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sp.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if sp.SurplusValue != 100 {
		t.Errorf("surplus_value = %d, want 100", sp.SurplusValue)
	}
	// Critical: round-half-up gives 1429, not 1428.
	if sp.IndividualProfitRate != 1429 {
		t.Errorf("individual_profit_rate = %d, want 1429", sp.IndividualProfitRate)
	}
	if sp.ConstantCapital != 600 {
		t.Errorf("constant_capital = %d, want 600", sp.ConstantCapital)
	}
}

// Sphere B (100c+600v, s'=100%): same total capital, higher v/C, higher p'.
func TestCreateProductionSphere_SphereB(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"name":"Sphere B (100c+600v)","c":700,"v":600,"s_rate":10000}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/spheres", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	var sp avgprofit.ProductionSphere
	if err := json.NewDecoder(res.Body).Decode(&sp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sp.IndividualProfitRate != 8571 {
		t.Errorf("individual_profit_rate = %d, want 8571", sp.IndividualProfitRate)
	}
}

// Bad JSON returns 400.
func TestCreateProductionSphere_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/avgprofit/spheres", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// Negative c returns 422.
func TestCreateProductionSphere_NegativeC(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"name":"bad","c":-100,"v":50,"s_rate":10000}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/spheres", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// V > C returns 422.
func TestCreateProductionSphere_VExceedsC(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"name":"bad","c":100,"v":200,"s_rate":10000}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/spheres", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// GET /v1/avgprofit/spheres/{id} with an unknown id returns 404.
func TestGetProductionSphere_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/avgprofit/spheres/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// GET /v1/avgprofit/spheres returns 200 with a non-nil items array even when empty.
func TestListProductionSpheres_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/avgprofit/spheres")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []avgprofit.ProductionSphere `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}

// Create then GET round-trip via HTTP.
func TestCreateThenGetProductionSphere(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	// Create European capital sphere.
	res, err := http.Post(ts.URL+"/v1/avgprofit/spheres", "application/json",
		strings.NewReader(`{"name":"European capital (84c+16v)","c":100,"v":16,"s_rate":10000}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	var created avgprofit.ProductionSphere
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	if created.IndividualProfitRate != 1600 {
		t.Errorf("created individual_profit_rate = %d, want 1600", created.IndividualProfitRate)
	}

	// Fetch by id.
	got, err := http.Get(ts.URL + "/v1/avgprofit/spheres/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched avgprofit.ProductionSphere
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.IndividualProfitRate != 1600 {
		t.Errorf("fetched = %+v, want id %s and rate 1600", fetched, created.ID)
	}
}
