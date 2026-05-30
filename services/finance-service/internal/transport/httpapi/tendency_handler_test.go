package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/tendency"
)

// TestCreateCompositionTrajectory_Created verifies POST with the §law fixture
// returns 201 + a Location header, and that the response carries the derived
// falling profit_rates [6667, 5000, 3333, 2500, 2000].
func TestCreateCompositionTrajectory_Created(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"label":"Marx §law s'=100%","surplus_value_rate":10000,` +
		`"steps":[[50,100],[100,100],[200,100],[300,100],[400,100]]}`
	res, err := http.Post(ts.URL+"/v1/tendency/trajectory", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/tendency/trajectory/") {
		t.Errorf("Location = %q, want /v1/tendency/trajectory/ prefix", loc)
	}

	var tr tendency.CompositionTrajectory
	if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if tr.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	want := []int64{6667, 5000, 3333, 2500, 2000}
	if len(tr.ProfitRates) != len(want) {
		t.Fatalf("profit_rates len = %d, want %d", len(tr.ProfitRates), len(want))
	}
	for i, w := range want {
		if tr.ProfitRates[i] != w {
			t.Errorf("profit_rates[%d] = %d, want %d", i, tr.ProfitRates[i], w)
		}
	}
}

// TestCreateCompositionTrajectory_BadJSON returns 400 on a malformed body.
func TestCreateCompositionTrajectory_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/tendency/trajectory", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateCompositionTrajectory_EmptySteps returns 400 when steps is empty.
func TestCreateCompositionTrajectory_EmptySteps(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"label":"empty","surplus_value_rate":10000,"steps":[]}`
	res, err := http.Post(ts.URL+"/v1/tendency/trajectory", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestListCompositionTrajectories_NeverNull asserts the list endpoint returns
// 200 and an items array that is non-null even when the store is empty.
func TestListCompositionTrajectories_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/tendency/trajectory")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []tendency.CompositionTrajectory `json:"items"`
	}
	// Decode into a struct whose Items defaults to nil; assert the JSON carried
	// a real array, not null, by re-reading the raw body.
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
	if err := json.Unmarshal(items, &payload.Items); err != nil {
		t.Fatalf("unmarshal items: %v", err)
	}
	if payload.Items == nil {
		t.Error("items decoded to nil, want non-nil slice")
	}
}

// TestGetCompositionTrajectory_NotFound returns 404 for an unknown id.
func TestGetCompositionTrajectory_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/tendency/trajectory/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateThenGetCompositionTrajectory exercises the create→GET round-trip over
// HTTP: known id returns 200 and the same derived profit_rates.
func TestCreateThenGetCompositionTrajectory(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"label":"law","surplus_value_rate":10000,"steps":[[50,100],[400,100]]}`
	res, err := http.Post(ts.URL+"/v1/tendency/trajectory", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created tendency.CompositionTrajectory
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/tendency/trajectory/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched tendency.CompositionTrajectory
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("fetched id = %s, want %s", fetched.ID, created.ID)
	}
	if len(fetched.ProfitRates) != 2 || fetched.ProfitRates[0] != 6667 || fetched.ProfitRates[1] != 2000 {
		t.Errorf("fetched profit_rates = %v, want [6667 2000]", fetched.ProfitRates)
	}
}

// TestCreateRateMassContradiction_Created verifies POST with the mass-preserved
// fixture returns 201 + Location and the computed masses.
func TestCreateRateMassContradiction_Created(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"old_c":1000000,"old_rate":2000,"new_c":2000000,"new_rate":1000}`
	res, err := http.Post(ts.URL+"/v1/tendency/rate-mass", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/tendency/rate-mass/") {
		t.Errorf("Location = %q, want /v1/tendency/rate-mass/ prefix", loc)
	}
	var rm tendency.RateMassContradiction
	if err := json.NewDecoder(res.Body).Decode(&rm); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if rm.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if rm.OldMass != 200000 || rm.NewMass != 200000 || rm.MassChange != 0 {
		t.Errorf("masses = old %d new %d change %d, want 200000/200000/0",
			rm.OldMass, rm.NewMass, rm.MassChange)
	}
}

// TestCreateRateMassContradiction_BadJSON returns 400 on a malformed body.
func TestCreateRateMassContradiction_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/tendency/rate-mass", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestGetRateMassContradiction_NotFound returns 404 for an unknown id.
func TestGetRateMassContradiction_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/tendency/rate-mass/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateThenGetRateMassContradiction exercises the create→GET round-trip for
// the rate-mass endpoint (mass-rose-while-rate-fell fixture).
func TestCreateThenGetRateMassContradiction(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"old_c":1000000,"old_rate":2000,"new_c":3000000,"new_rate":1000}`
	res, err := http.Post(ts.URL+"/v1/tendency/rate-mass", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created tendency.RateMassContradiction
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/tendency/rate-mass/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched tendency.RateMassContradiction
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("fetched id = %s, want %s", fetched.ID, created.ID)
	}
	if fetched.NewMass != 300000 || fetched.MassChange <= 0 {
		t.Errorf("fetched NewMass=%d MassChange=%d, want 300000 and >0", fetched.NewMass, fetched.MassChange)
	}
}
