package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// §I.4a: s' held at 100%, value-composition raised (c 80→100); p' falls
// 20% → 16⅔% (2000 → 1666 bp).
func TestCreateVariation_Marx(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"initial":{"c":80,"v":20,"s_rate":10000},"changed":{"c":100,"v":20,"s_rate":10000}}`
	res, err := http.Post(ts.URL+"/v1/profit/variation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/profit/variation/") {
		t.Errorf("Location = %q, want /v1/profit/variation/ prefix", loc)
	}
	var a profit.VariationAnalysis
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.Case != profit.CaseSConstantVCVariable {
		t.Errorf("case = %q, want %q", a.Case, profit.CaseSConstantVCVariable)
	}
	if a.OldProfitRate != 2000 || a.NewProfitRate != 1666 {
		t.Errorf("rates = %d → %d, want 2000 → 1666", a.OldProfitRate, a.NewProfitRate)
	}
	if a.ProportionalChange {
		t.Error("ProportionalChange = true, want false for the s'-constant case")
	}
	if a.ID.IsZero() {
		t.Error("expected an assigned id")
	}
}

func TestCreateVariation_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/variation", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

func TestCreateVariation_Negative(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"initial":{"c":-1,"v":20,"s_rate":10000},"changed":{"c":100,"v":20,"s_rate":10000}}`
	res, err := http.Post(ts.URL+"/v1/profit/variation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestGetVariation_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/variation/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

func TestCreateThenGetVariation(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	// §II.1: v/C held at 20%, s' halved 100% → 50%; p' falls in proportion.
	res, err := http.Post(ts.URL+"/v1/profit/variation", "application/json",
		strings.NewReader(`{"initial":{"c":80,"v":20,"s_rate":10000},"changed":{"c":160,"v":40,"s_rate":5000}}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created profit.VariationAnalysis
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.Case != profit.CaseSVariableVCConstant || !created.ProportionalChange {
		t.Errorf("got case %q proportional=%v, want s_variable_vc_constant + true",
			created.Case, created.ProportionalChange)
	}

	got, err := http.Get(ts.URL + "/v1/profit/variation/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched profit.VariationAnalysis
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.NewProfitRate != created.NewProfitRate {
		t.Errorf("fetched = %+v, want id %s", fetched, created.ID)
	}
}

func TestListVariations_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/variation")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []profit.VariationAnalysis `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}

// §I.1: two capitals with equal s' have rates of profit proportional to their
// value-compositions. Marx's 15,000 capital: 12,000c+3,000v → 20%; 13,000c+
// 2,000v → 13⅓%.
func TestCompareProfitRates_EqualSurplus(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"first":{"c":12000,"v":3000,"s_rate":10000},"second":{"c":13000,"v":2000,"s_rate":10000}}`
	res, err := http.Post(ts.URL+"/v1/profit/compare", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var c profit.ProfitRateComparison
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !c.EqualSurplusRate {
		t.Error("equal_surplus_rate = false, want true")
	}
	if c.FirstProfitRate != 2000 || c.SecondProfitRate != 1333 {
		t.Errorf("rates = %d, %d, want 2000, 1333", c.FirstProfitRate, c.SecondProfitRate)
	}
}

func TestCompareProfitRates_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/compare", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}
