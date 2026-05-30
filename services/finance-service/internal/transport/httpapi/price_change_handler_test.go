package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// TestCreatePriceOfProductionChange_Created verifies 201 + Location header and
// that a general-rate change derives RateChanged true / ValueChanged false /
// PriceChanged true.
func TestCreatePriceOfProductionChange_Created(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"cotton","cause":"general_rate_change","old_price":122,"new_price":130}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-change", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/avgprofit/price-change/") {
		t.Errorf("Location = %q, want /v1/avgprofit/price-change/ prefix", loc)
	}
	var c avgprofit.PriceOfProductionChange
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if c.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if !c.RateChanged || c.ValueChanged || !c.PriceChanged {
		t.Errorf("flags = rate %v value %v price %v, want true/false/true",
			c.RateChanged, c.ValueChanged, c.PriceChanged)
	}
}

// TestGetPriceOfProductionChange_NotFound returns 404 for an unknown id.
func TestGetPriceOfProductionChange_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/avgprofit/price-change/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreatePriceOfProductionChange_BadJSON returns 400 on malformed body.
func TestCreatePriceOfProductionChange_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-change", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreatePriceOfProductionChange_NegativePrice returns 422 on a negative price.
func TestCreatePriceOfProductionChange_NegativePrice(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"sphere_name":"cotton","cause":"general_rate_change","old_price":-1,"new_price":130}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-change", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestCreateThenGetPriceOfProductionChange verifies the create→GET round-trip
// over HTTP.
func TestCreateThenGetPriceOfProductionChange(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"linen","cause":"individual_value_change","old_price":122,"new_price":110}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/price-change", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created avgprofit.PriceOfProductionChange
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/avgprofit/price-change/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched avgprofit.PriceOfProductionChange
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || !fetched.ValueChanged || fetched.RateChanged {
		t.Errorf("fetched = %+v, want id %s value-changed", fetched, created.ID)
	}
}

// TestComputeCompensationGround_OK verifies the compensation-ground myth-buster
// returns 200 with the constant ActuallyIs and CreatesValue false.
func TestComputeCompensationGround_OK(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"reason":"slow turnover","appears_to_be":"extra value creation"}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/compensation-ground", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var g avgprofit.CompensationGround
	if err := json.NewDecoder(res.Body).Decode(&g); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if g.ActuallyIs != "share in aggregate surplus-value" {
		t.Errorf("ActuallyIs = %q, want %q", g.ActuallyIs, "share in aggregate surplus-value")
	}
	if g.CreatesValue {
		t.Error("CreatesValue = true, want false")
	}
}

// TestComputeCompensationGround_BadJSON returns 400 on malformed body.
func TestComputeCompensationGround_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/avgprofit/compensation-ground", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}
