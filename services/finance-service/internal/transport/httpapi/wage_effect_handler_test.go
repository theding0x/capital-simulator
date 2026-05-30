package httpapi

// Tests for Ch. 11 — Effects of General Wage Fluctuations on Prices of Production.
// All fixtures from design-ch11.md; prices in centi-units (×100).

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// TestCreateWageEffectAnalysis_201Location verifies 201 + Location header.
func TestCreateWageEffectAnalysis_201Location(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"base_constant":80,"base_variable":20,"s_rate":10000,"wage_factor":125,"spheres":[]}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/wage-effect", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/avgprofit/wage-effect/") {
		t.Errorf("Location = %q, want /v1/avgprofit/wage-effect/ prefix", loc)
	}

	var a avgprofit.WageEffectAnalysis
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.ID.IsZero() {
		t.Error("expected an assigned id")
	}
}

// TestCreateWageEffectAnalysis_BadJSON returns 400 on malformed body.
func TestCreateWageEffectAnalysis_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/avgprofit/wage-effect", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateWageEffectAnalysis_NegativeMagnitude returns 422 on negative base_constant.
func TestCreateWageEffectAnalysis_NegativeMagnitude(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"base_constant":-80,"base_variable":20,"s_rate":10000,"wage_factor":125,"spheres":[]}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/wage-effect", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestCreateWageEffectAnalysis_NegativeSphereC returns 422 when a sphere c is negative.
func TestCreateWageEffectAnalysis_NegativeSphereC(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"base_constant":80,"base_variable":20,"s_rate":10000,"wage_factor":125,"spheres":[{"name":"bad","c":-1,"v":50}]}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/wage-effect", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestGetWageEffectAnalysis_NotFound returns 404 for an unknown id.
func TestGetWageEffectAnalysis_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/avgprofit/wage-effect/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestListWageEffectAnalyses_NonNilItems verifies 200 with non-nil items array.
func TestListWageEffectAnalyses_NonNilItems(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/avgprofit/wage-effect")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []avgprofit.WageEffectAnalysis `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be non-nil (never JSON null)")
	}
}

// TestCreateWageEffectAnalysis_EndToEnd_MainFixture is a full round-trip test using
// the §main fixture from the design contract (base 80c+20v, factor 125, lower 50c+50v
// and higher 92c+8v). Asserts new_general_rate 1429 and correct rise/fall directions.
func TestCreateWageEffectAnalysis_EndToEnd_MainFixture(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"base_constant": 80,
		"base_variable": 20,
		"s_rate":        10000,
		"wage_factor":   125,
		"spheres": [
			{"name": "lower (50c+50v)",  "c": 50, "v": 50},
			{"name": "higher (92c+8v)", "c": 92, "v":  8}
		]
	}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/wage-effect", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	var a avgprofit.WageEffectAnalysis
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// General rates
	if a.OldGeneralRate != 2000 {
		t.Errorf("old_general_rate = %d, want 2000", a.OldGeneralRate)
	}
	if a.NewGeneralRate != 1429 {
		t.Errorf("new_general_rate = %d, want 1429", a.NewGeneralRate)
	}
	if a.Kind != avgprofit.KindWageRise {
		t.Errorf("kind = %q, want KindWageRise", a.Kind)
	}

	// Average outcome: price unchanged
	if a.AverageOutcome.Direction != avgprofit.KindPriceUnchanged {
		t.Errorf("average_outcome.direction = %q, want unchanged", a.AverageOutcome.Direction)
	}

	// Outcomes slice
	if len(a.Outcomes) != 2 {
		t.Fatalf("outcomes len = %d, want 2", len(a.Outcomes))
	}

	// Lower sphere (50c+50v) rises
	lower := a.Outcomes[0]
	if lower.Direction != avgprofit.KindPriceRises {
		t.Errorf("lower direction = %q, want rises", lower.Direction)
	}
	if lower.NewPrice != 12857 {
		t.Errorf("lower new_price = %d, want 12857", lower.NewPrice)
	}

	// Higher sphere (92c+8v) falls
	higher := a.Outcomes[1]
	if higher.Direction != avgprofit.KindPriceFalls {
		t.Errorf("higher direction = %q, want falls", higher.Direction)
	}
	if higher.NewPrice != 11657 {
		t.Errorf("higher new_price = %d, want 11657", higher.NewPrice)
	}

	// GET the saved record
	id := string(a.ID)
	getRes, err := http.Get(ts.URL + "/v1/avgprofit/wage-effect/" + id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", getRes.StatusCode)
	}
	var fetched avgprofit.WageEffectAnalysis
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != a.ID || fetched.NewGeneralRate != a.NewGeneralRate {
		t.Errorf("get round-trip mismatch: id=%s rate=%d", fetched.ID, fetched.NewGeneralRate)
	}
}
