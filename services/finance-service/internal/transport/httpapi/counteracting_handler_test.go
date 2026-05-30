package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/tendency"
)

// TestCreateCounteractingForce_Created posts a valid kind and expects 201 + a
// Location header, with the response carrying the derived ExcludedFromGeneralRate
// flag (false for an industrial force).
func TestCreateCounteractingForce_Created(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/tendency/counteracting-force", "application/json",
		strings.NewReader(`{"kind":"foreign_trade","note":"colonial inputs"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/tendency/counteracting-force/") {
		t.Errorf("Location = %q, want /v1/tendency/counteracting-force/ prefix", loc)
	}
	var f tendency.CounteractingForce
	if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if f.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if f.Kind != tendency.KindForeignTrade {
		t.Errorf("kind = %q, want foreign_trade", f.Kind)
	}
	if f.ExcludedFromGeneralRate {
		t.Error("ExcludedFromGeneralRate = true, want false for foreign_trade")
	}
}

// TestCreateCounteractingForce_StockExcluded posts the §VI stock-capital kind and
// asserts the response flags it excluded from the general rate.
func TestCreateCounteractingForce_StockExcluded(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/tendency/counteracting-force", "application/json",
		strings.NewReader(`{"kind":"stock_capital"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	var f tendency.CounteractingForce
	if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !f.ExcludedFromGeneralRate {
		t.Error("ExcludedFromGeneralRate = false, want true for stock_capital")
	}
}

// TestCreateCounteractingForce_InvalidKind returns 422 on an unknown kind.
func TestCreateCounteractingForce_InvalidKind(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/tendency/counteracting-force", "application/json",
		strings.NewReader(`{"kind":"usury"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestCreateCounteractingForce_BadJSON returns 400 on a malformed body.
func TestCreateCounteractingForce_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/tendency/counteracting-force", "application/json",
		strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestGetCounteractingForce_NotFound returns 404 for an unknown id.
func TestGetCounteractingForce_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/tendency/counteracting-force/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateThenGetCounteractingForce exercises the create→GET round-trip: a known
// id returns 200 and the same kind.
func TestCreateThenGetCounteractingForce(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/tendency/counteracting-force", "application/json",
		strings.NewReader(`{"kind":"intensified_exploitation"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created tendency.CounteractingForce
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/tendency/counteracting-force/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched tendency.CounteractingForce
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.Kind != tendency.KindIntensifiedExploitation {
		t.Errorf("fetched = %+v, want id %s and kind intensified_exploitation", fetched, created.ID)
	}
}

// scenarioBody is the §law scenario request with the five industrial force kinds.
const scenarioBody = `{"label":"All five industrial forces — law still holds",` +
	`"force_kinds":["intensified_exploitation","depressed_wages","cheaper_constant_capital","relative_overpopulation","foreign_trade"],` +
	`"steps":[[50,100],[100,100],[200,100],[300,100],[400,100]],"surplus_value_rate":10000}`

// TestCreateCounteractingScenario_Created posts the §law scenario and expects 201
// + a Location header, with the response carrying modified_rates
// [8865,6650,4431,3325,2660] and the law preserved (final < initial).
func TestCreateCounteractingScenario_Created(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/tendency/scenario", "application/json", strings.NewReader(scenarioBody))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/tendency/scenario/") {
		t.Errorf("Location = %q, want /v1/tendency/scenario/ prefix", loc)
	}
	var s tendency.CounteractingScenario
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if s.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	want := []int64{8865, 6650, 4431, 3325, 2660}
	if len(s.ModifiedRates) != len(want) {
		t.Fatalf("modified_rates len = %d, want %d", len(s.ModifiedRates), len(want))
	}
	for i, w := range want {
		if s.ModifiedRates[i] != w {
			t.Errorf("modified_rates[%d] = %d, want %d", i, s.ModifiedRates[i], w)
		}
	}
	if s.InitialProfitRate != 8865 || s.FinalProfitRate != 2660 {
		t.Errorf("rates = init %d final %d, want 8865/2660", s.InitialProfitRate, s.FinalProfitRate)
	}
	if s.FinalProfitRate >= s.InitialProfitRate {
		t.Errorf("law not preserved: final %d >= initial %d", s.FinalProfitRate, s.InitialProfitRate)
	}
}

// TestCreateCounteractingScenario_InvalidKind returns 422 when a force kind is
// unknown.
func TestCreateCounteractingScenario_InvalidKind(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"label":"bad","force_kinds":["usury"],"steps":[[50,100]],"surplus_value_rate":10000}`
	res, err := http.Post(ts.URL+"/v1/tendency/scenario", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestCreateCounteractingScenario_EmptyForceKinds returns 400 when force_kinds is
// empty.
func TestCreateCounteractingScenario_EmptyForceKinds(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"label":"x","force_kinds":[],"steps":[[50,100]],"surplus_value_rate":10000}`
	res, err := http.Post(ts.URL+"/v1/tendency/scenario", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateCounteractingScenario_EmptySteps returns 400 when steps is empty.
func TestCreateCounteractingScenario_EmptySteps(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"label":"x","force_kinds":["foreign_trade"],"steps":[],"surplus_value_rate":10000}`
	res, err := http.Post(ts.URL+"/v1/tendency/scenario", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateCounteractingScenario_BadJSON returns 400 on a malformed body.
func TestCreateCounteractingScenario_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/tendency/scenario", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestListCounteractingScenarios_NeverNull asserts the list endpoint returns 200
// and an items array that is non-null even when the store is empty.
func TestListCounteractingScenarios_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/tendency/scenario")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		t.Fatalf("decode raw: %v", err)
	}
	items, ok := raw["items"]
	if !ok {
		t.Fatal("response has no items key")
	}
	if string(items) == "null" {
		t.Errorf("items = null, want [] for an empty store")
	}
	var decoded []tendency.CounteractingScenario
	if err := json.Unmarshal(items, &decoded); err != nil {
		t.Fatalf("unmarshal items: %v", err)
	}
	if decoded == nil {
		t.Error("items decoded to nil, want non-nil slice")
	}
}

// TestGetCounteractingScenario_NotFound returns 404 for an unknown id.
func TestGetCounteractingScenario_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/tendency/scenario/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateThenGetCounteractingScenario exercises the create→GET round-trip over
// HTTP: a known id returns 200 and the same modified_rates.
func TestCreateThenGetCounteractingScenario(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/tendency/scenario", "application/json", strings.NewReader(scenarioBody))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created tendency.CounteractingScenario
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/tendency/scenario/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched tendency.CounteractingScenario
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("fetched id = %s, want %s", fetched.ID, created.ID)
	}
	if len(fetched.ModifiedRates) != 5 || fetched.ModifiedRates[0] != 8865 || fetched.ModifiedRates[4] != 2660 {
		t.Errorf("fetched modified_rates = %v, want [8865 ... 2660]", fetched.ModifiedRates)
	}
}
