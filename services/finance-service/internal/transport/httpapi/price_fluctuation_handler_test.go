package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// Marx's §I rising branch: £320 fixed + £80 raw material + £100 v + £100 s, the
// raw-material price doubled (factor 200), drops p' from 20% (2000 bp) to 17.24%
// (1724 bp).
func TestCreatePriceFluctuationAnalysis_Marx(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"fixed":320,"original_material":80,"price_factor":200,"variable":100,"surplus_value":100}`
	res, err := http.Post(ts.URL+"/v1/profit/price-fluctuation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/profit/price-fluctuation/") {
		t.Errorf("Location = %q, want /v1/profit/price-fluctuation/ prefix", loc)
	}
	var a profit.PriceFluctuationAnalysis
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.OldProfitRate != 2000 || a.NewProfitRate != 1724 {
		t.Errorf("got old=%d new=%d, want 2000/1724", a.OldProfitRate, a.NewProfitRate)
	}
	if a.Kind != profit.KindRise {
		t.Errorf("kind = %q, want rise", a.Kind)
	}
	if a.ID.IsZero() {
		t.Error("expected an assigned id")
	}
}

func TestCreatePriceFluctuationAnalysis_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/price-fluctuation", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

func TestCreatePriceFluctuationAnalysis_NonPositiveFactor(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"fixed":320,"original_material":80,"price_factor":0,"variable":100,"surplus_value":100}`
	res, err := http.Post(ts.URL+"/v1/profit/price-fluctuation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestCreatePriceFluctuationAnalysis_Negative(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"fixed":-1,"original_material":80,"price_factor":200,"variable":100,"surplus_value":100}`
	res, err := http.Post(ts.URL+"/v1/profit/price-fluctuation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestGetPriceFluctuationAnalysis_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/price-fluctuation/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// §I falling branch: £20 fixed + £380 raw material, the price halved (factor 50),
// raises p' from 20% (2000 bp) to 32.26% (3225 bp, truncated).
func TestCreateThenGetPriceFluctuationAnalysis(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/price-fluctuation", "application/json",
		strings.NewReader(`{"fixed":20,"original_material":380,"price_factor":50,"variable":100,"surplus_value":100}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created profit.PriceFluctuationAnalysis
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.OldProfitRate != 2000 || created.NewProfitRate != 3225 {
		t.Errorf("got old=%d new=%d, want 2000/3225", created.OldProfitRate, created.NewProfitRate)
	}
	if created.Kind != profit.KindFall {
		t.Errorf("kind = %q, want fall", created.Kind)
	}

	got, err := http.Get(ts.URL + "/v1/profit/price-fluctuation/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched profit.PriceFluctuationAnalysis
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.NewProfitRate != created.NewProfitRate {
		t.Errorf("fetched = %+v, want id %s", fetched, created.ID)
	}
}

func TestListPriceFluctuationAnalyses_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/price-fluctuation")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []profit.PriceFluctuationAnalysis `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}
