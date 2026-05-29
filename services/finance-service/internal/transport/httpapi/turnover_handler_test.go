package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// Marx's Capital A: 100C, 20v, s' = 100%, n = 2 → annual p' = 40% (4000 bp),
// single-turnover p' = 20% (2000 bp), S' = 200% (20000 bp).
func TestCreateTurnoverAnalysis_Marx(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"c":100,"v":20,"s_rate":10000,"n":2}`
	res, err := http.Post(ts.URL+"/v1/profit/turnover-analysis", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/profit/turnover-analysis/") {
		t.Errorf("Location = %q, want /v1/profit/turnover-analysis/ prefix", loc)
	}
	var a profit.TurnoverAnalysis
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.AnnualProfitRate.BasisPoints != 4000 || a.SingleTurnoverProfitRate != 2000 || a.AnnualSurplusValueRate != 20000 {
		t.Errorf("got annual=%d single=%d S'=%d, want 4000/2000/20000",
			a.AnnualProfitRate.BasisPoints, a.SingleTurnoverProfitRate, a.AnnualSurplusValueRate)
	}
	if a.AnnualWages != 40 {
		t.Errorf("annual wages = %d, want 40", a.AnnualWages)
	}
	if a.ID.IsZero() {
		t.Error("expected an assigned id")
	}
}

func TestCreateTurnoverAnalysis_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/turnover-analysis", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

func TestCreateTurnoverAnalysis_Negative(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"c":100,"v":20,"s_rate":10000,"n":-1}`
	res, err := http.Post(ts.URL+"/v1/profit/turnover-analysis", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestCreateTurnoverAnalysis_VariableExceedsTotal(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"c":100,"v":200,"s_rate":10000,"n":2}`
	res, err := http.Post(ts.URL+"/v1/profit/turnover-analysis", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestGetTurnoverAnalysis_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/turnover-analysis/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

func TestCreateThenGetTurnoverAnalysis(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	// Capital B: 200C, 40v, s' = 100%, n = 1 → annual p' = 20%.
	res, err := http.Post(ts.URL+"/v1/profit/turnover-analysis", "application/json",
		strings.NewReader(`{"c":200,"v":40,"s_rate":10000,"n":1}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created profit.TurnoverAnalysis
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.AnnualProfitRate.BasisPoints != 2000 {
		t.Errorf("annual p' = %d, want 2000", created.AnnualProfitRate.BasisPoints)
	}

	got, err := http.Get(ts.URL + "/v1/profit/turnover-analysis/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched profit.TurnoverAnalysis
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.AnnualProfitRate.BasisPoints != created.AnnualProfitRate.BasisPoints {
		t.Errorf("fetched = %+v, want id %s", fetched, created.ID)
	}
}

func TestListTurnoverAnalyses_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/turnover-analysis")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []profit.TurnoverAnalysis `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}
