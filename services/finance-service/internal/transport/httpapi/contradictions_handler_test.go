package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/tendency"
)

// --- POST/GET /v1/tendency/crisis ---

// TestCreateCrisis_Created posts the §crisis fixture {30,1000} and expects 201 +
// a Location header, with the response carrying the derived post-crisis rate 1429.
func TestCreateCrisis_Created(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/tendency/crisis", "application/json",
		strings.NewReader(`{"constant_capital_writedown":30,"pre_crisis_profit_rate":1000}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/tendency/crisis/") {
		t.Errorf("Location = %q, want /v1/tendency/crisis/ prefix", loc)
	}
	var c tendency.Crisis
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if c.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if c.PostCrisisProfitRate != 1429 {
		t.Errorf("post_crisis_profit_rate = %d, want 1429", c.PostCrisisProfitRate)
	}
	if c.PostCrisisProfitRate <= c.PreCrisisProfitRate {
		t.Errorf("post %d must exceed pre %d", c.PostCrisisProfitRate, c.PreCrisisProfitRate)
	}
}

// TestCreateCrisis_InvalidWritedown returns 422 when the writedown is at the
// (0,100) boundary (0 or 100).
func TestCreateCrisis_InvalidWritedown(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	for _, w := range []string{"0", "100"} {
		res, err := http.Post(ts.URL+"/v1/tendency/crisis", "application/json",
			strings.NewReader(`{"constant_capital_writedown":`+w+`,"pre_crisis_profit_rate":1000}`))
		if err != nil {
			t.Fatalf("post writedown=%s: %v", w, err)
		}
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("writedown=%s: status = %d, want 422", w, res.StatusCode)
		}
		res.Body.Close()
	}
}

// TestCreateCrisis_BadJSON returns 400 on a malformed body.
func TestCreateCrisis_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/tendency/crisis", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestListCrises_NeverNull asserts the list endpoint returns 200 and an items
// array that is non-null even when the store is empty.
func TestListCrises_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/tendency/crisis")
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
	var decoded []tendency.Crisis
	if err := json.Unmarshal(items, &decoded); err != nil {
		t.Fatalf("unmarshal items: %v", err)
	}
	if decoded == nil {
		t.Error("items decoded to nil, want non-nil slice")
	}
}

// TestGetCrisis_NotFound returns 404 for an unknown id.
func TestGetCrisis_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/tendency/crisis/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateThenGetCrisis exercises the create→GET round-trip: a known id returns
// 200 and the same restored rate.
func TestCreateThenGetCrisis(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/tendency/crisis", "application/json",
		strings.NewReader(`{"constant_capital_writedown":30,"pre_crisis_profit_rate":1000}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created tendency.Crisis
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/tendency/crisis/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched tendency.Crisis
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.PostCrisisProfitRate != 1429 {
		t.Errorf("fetched = %+v, want id %s and post-rate 1429", fetched, created.ID)
	}
}

// --- POST/GET /v1/tendency/contradiction ---

// TestCreateInternalContradiction_Created posts a valid kind and note and expects
// 201 + a Location header, with the derived is_coexistent flag true.
func TestCreateInternalContradiction_Created(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/tendency/contradiction", "application/json",
		strings.NewReader(`{"kind":"falling_rate_rising_mass","note":"Marx §I"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/tendency/contradiction/") {
		t.Errorf("Location = %q, want /v1/tendency/contradiction/ prefix", loc)
	}
	var c tendency.InternalContradiction
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if c.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if c.Kind != tendency.KindFallingRateRisingMass {
		t.Errorf("kind = %q, want falling_rate_rising_mass", c.Kind)
	}
	if !c.IsCoexistent {
		t.Error("is_coexistent = false, want true")
	}
	if c.Note != "Marx §I" {
		t.Errorf("note = %q, want %q", c.Note, "Marx §I")
	}
}

// TestCreateInternalContradiction_InvalidKind returns 422 on an unknown kind.
func TestCreateInternalContradiction_InvalidKind(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/tendency/contradiction", "application/json",
		strings.NewReader(`{"kind":"usury"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestCreateInternalContradiction_BadJSON returns 400 on a malformed body.
func TestCreateInternalContradiction_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/tendency/contradiction", "application/json",
		strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestGetInternalContradiction_NotFound returns 404 for an unknown id.
func TestGetInternalContradiction_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/tendency/contradiction/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateThenGetInternalContradiction exercises the create→GET round-trip: a
// known id returns 200 and the same kind.
func TestCreateThenGetInternalContradiction(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/tendency/contradiction", "application/json",
		strings.NewReader(`{"kind":"concentration_centralisation"}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created tendency.InternalContradiction
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/tendency/contradiction/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched tendency.InternalContradiction
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.Kind != tendency.KindConcentrationCentralisation {
		t.Errorf("fetched = %+v, want id %s and kind concentration_centralisation", fetched, created.ID)
	}
}

// --- GET /v1/tendency/summary ---

// TestGetPartIIISummary_Empty asserts the summary endpoint returns 200 with zero
// counts and no populated latest fields on an empty store — no panic, and the
// count fields are present numbers rather than null.
func TestGetPartIIISummary_Empty(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/tendency/summary")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}

	// Count fields must be present numbers (not null) on an empty store.
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		t.Fatalf("decode raw: %v", err)
	}
	for _, k := range []string{"trajectory_count", "scenario_count", "crisis_count", "contradiction_count"} {
		v, ok := raw[k]
		if !ok {
			t.Errorf("summary missing %q", k)
			continue
		}
		if string(v) == "null" {
			t.Errorf("%q = null, want a number", k)
		}
	}

	var summary tendency.PartIIISummary
	if err := json.Unmarshal(rawJSON(raw), &summary); err != nil {
		t.Fatalf("unmarshal summary: %v", err)
	}
	if summary.CrisisCount != 0 || summary.ContradictionCount != 0 {
		t.Errorf("empty-store counts = crisis %d contradiction %d, want 0/0",
			summary.CrisisCount, summary.ContradictionCount)
	}
	if summary.LatestCrisis != nil || summary.LatestContradiction != nil {
		t.Error("empty-store latest fields should be nil")
	}
}

// TestGetPartIIISummary_AfterCreate creates a crisis and a contradiction, then
// asserts the summary reflects them: counts of at least one each and the
// latest-of-each fields populated. (The in-memory store loads no SQL seeds, so
// the counts here are exactly the records this test created.)
func TestGetPartIIISummary_AfterCreate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	crisisRes, err := http.Post(ts.URL+"/v1/tendency/crisis", "application/json",
		strings.NewReader(`{"constant_capital_writedown":30,"pre_crisis_profit_rate":1000}`))
	if err != nil {
		t.Fatalf("post crisis: %v", err)
	}
	crisisRes.Body.Close()

	contraRes, err := http.Post(ts.URL+"/v1/tendency/contradiction", "application/json",
		strings.NewReader(`{"kind":"falling_rate_rising_mass","note":"§I"}`))
	if err != nil {
		t.Fatalf("post contradiction: %v", err)
	}
	contraRes.Body.Close()

	res, err := http.Get(ts.URL + "/v1/tendency/summary")
	if err != nil {
		t.Fatalf("get summary: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var summary tendency.PartIIISummary
	if err := json.NewDecoder(res.Body).Decode(&summary); err != nil {
		t.Fatalf("decode summary: %v", err)
	}

	if summary.CrisisCount < 1 {
		t.Errorf("crisis_count = %d, want >= 1", summary.CrisisCount)
	}
	if summary.ContradictionCount < 1 {
		t.Errorf("contradiction_count = %d, want >= 1", summary.ContradictionCount)
	}
	if summary.LatestCrisis == nil {
		t.Error("latest_crisis is nil, want the crisis just created")
	} else if summary.LatestCrisis.PostCrisisProfitRate != 1429 {
		t.Errorf("latest_crisis.post_crisis_profit_rate = %d, want 1429", summary.LatestCrisis.PostCrisisProfitRate)
	}
	if summary.LatestContradiction == nil {
		t.Error("latest_contradiction is nil, want the contradiction just created")
	} else if summary.LatestContradiction.Kind != tendency.KindFallingRateRisingMass {
		t.Errorf("latest_contradiction.kind = %q, want falling_rate_rising_mass", summary.LatestContradiction.Kind)
	}
}

// rawJSON re-marshals a decoded raw map so it can be unmarshalled into a typed
// struct, letting one HTTP read drive both the presence check and the typed
// assertion.
func rawJSON(m map[string]json.RawMessage) []byte {
	b, _ := json.Marshal(m)
	return b
}
