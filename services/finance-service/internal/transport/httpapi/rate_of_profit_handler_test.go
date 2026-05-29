package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

func TestCreateProfitRate_Marx(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"constant":400,"variable":100,"surplus_value":100}`
	res, err := http.Post(ts.URL+"/v1/profit/rate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/profit/rate/") {
		t.Errorf("Location = %q, want /v1/profit/rate/ prefix", loc)
	}
	var a profit.ProfitRateAnalysis
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.ProfitRate.BasisPoints != 2000 || a.SurplusValueRate != 10000 || a.Mystification != 8000 {
		t.Errorf("got p'=%d s'=%d myst=%d, want 2000/10000/8000",
			a.ProfitRate.BasisPoints, a.SurplusValueRate, a.Mystification)
	}
	if a.ProfitRate.TotalCapital != 500 {
		t.Errorf("total capital = %d, want 500", a.ProfitRate.TotalCapital)
	}
	if a.ID.IsZero() {
		t.Error("expected an assigned id")
	}
}

func TestCreateProfitRate_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/rate", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

func TestCreateProfitRate_Negative(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"constant":-1,"variable":100,"surplus_value":100}`
	res, err := http.Post(ts.URL+"/v1/profit/rate", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestGetProfitRate_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/rate/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

func TestCreateThenGetProfitRate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/rate", "application/json",
		strings.NewReader(`{"constant":0,"variable":100,"surplus_value":100}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created profit.ProfitRateAnalysis
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// v = C, the limiting case: the two rates coincide and nothing is mystified.
	if created.ProfitRate.BasisPoints != int64(created.SurplusValueRate) || created.Mystification != 0 {
		t.Errorf("v=C case: p'=%d s'=%d myst=%d, want equal rates and 0 mystification",
			created.ProfitRate.BasisPoints, created.SurplusValueRate, created.Mystification)
	}

	got, err := http.Get(ts.URL + "/v1/profit/rate/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched profit.ProfitRateAnalysis
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.ProfitRate.BasisPoints != created.ProfitRate.BasisPoints {
		t.Errorf("fetched = %+v, want id %s", fetched, created.ID)
	}
}

func TestListProfitRates_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/rate")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []profit.ProfitRateAnalysis `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}
