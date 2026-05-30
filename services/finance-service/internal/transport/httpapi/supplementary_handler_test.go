package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// Marx's two-enterprise illustration: equal s = £1,000 on equal v = £1,000;
// A (c = £9,000) shows 10% (1000 bp), B (c = £11,000) 8 1/3% (833 bp); the gap is
// 167 bp.
func TestCreateCompositionEffect_Marx(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"s_a":1000,"v_a":1000,"c_a":9000,"s_b":1000,"v_b":1000,"c_b":11000}`
	res, err := http.Post(ts.URL+"/v1/profit/composition-effect", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/profit/composition-effect/") {
		t.Errorf("Location = %q, want /v1/profit/composition-effect/ prefix", loc)
	}
	var a profit.CompositionEffectAnalysis
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.ProfitRateA != 1000 || a.ProfitRateB != 833 || a.RateDifference != 167 {
		t.Errorf("got A=%d B=%d diff=%d, want 1000/833/167", a.ProfitRateA, a.ProfitRateB, a.RateDifference)
	}
	if a.ID.IsZero() {
		t.Error("expected an assigned id")
	}
}

func TestCreateCompositionEffect_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/composition-effect", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

func TestCreateCompositionEffect_Negative(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"s_a":1000,"v_a":1000,"c_a":-9000,"s_b":1000,"v_b":1000,"c_b":11000}`
	res, err := http.Post(ts.URL+"/v1/profit/composition-effect", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestGetCompositionEffect_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/composition-effect/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

func TestListCompositionEffects_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/composition-effect")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	var body struct {
		Items []profit.CompositionEffectAnalysis `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}

// Rodbertus refutation: capital £100, profit £20 (20%); a money-value change by
// factor 200 (gold halves) leaves the rate at 20% (2000 bp).
func TestCreateMagnitudeChange_Rodbertus(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"kind":"money_value_change","original_capital":100,"original_profit":20,"factor":200}`
	res, err := http.Post(ts.URL+"/v1/profit/magnitude-change", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	var a profit.MagnitudeChangeAnalysis
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.OldRate != 2000 || a.NewRate != 2000 || !a.RateUnchanged {
		t.Errorf("got old=%d new=%d unchanged=%v, want 2000/2000/true", a.OldRate, a.NewRate, a.RateUnchanged)
	}
}

// An organic change moves the rate: doubling the capital with the same profit
// halves the rate, 20% → 10%.
func TestCreateMagnitudeChange_Organic(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"kind":"organic_change","original_capital":100,"original_profit":20,"factor":200}`
	res, err := http.Post(ts.URL+"/v1/profit/magnitude-change", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var a profit.MagnitudeChangeAnalysis
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.NewRate != 1000 || a.RateUnchanged {
		t.Errorf("got new=%d unchanged=%v, want 1000/false", a.NewRate, a.RateUnchanged)
	}
}

func TestCreateMagnitudeChange_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/magnitude-change", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

func TestCreateMagnitudeChange_UnknownKind(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"kind":"bogus","original_capital":100,"original_profit":20,"factor":200}`
	res, err := http.Post(ts.URL+"/v1/profit/magnitude-change", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestListMagnitudeChanges_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/magnitude-change")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	var body struct {
		Items []profit.MagnitudeChangeAnalysis `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}

// The Part I summary names the determinants of the rate of profit and gathers the
// rate-of-profit analyses recorded so far. After recording one rate it carries at
// least that analysis.
func TestGetPartISummary(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	// Record one rate-of-profit analysis (Ch. 2) so the summary has something to gather.
	_, err := http.Post(ts.URL+"/v1/profit/rate", "application/json",
		strings.NewReader(`{"constant":9000,"variable":1000,"surplus_value":1000}`))
	if err != nil {
		t.Fatalf("seed rate: %v", err)
	}

	res, err := http.Get(ts.URL + "/v1/profit/summary")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var s profit.PartISummary
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(s.Determinants) == 0 {
		t.Error("summary must name the determinants of the rate of profit")
	}
	if len(s.Analyses) < 1 {
		t.Errorf("summary analyses: got %d, want >= 1", len(s.Analyses))
	}
	if s.ID.IsZero() {
		t.Error("summary must carry an ID")
	}
}
