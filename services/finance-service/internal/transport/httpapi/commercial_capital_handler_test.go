package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
)

// --- POST /v1/merchant/commercial-capital ---

// TestCreateCommercialCapital_LinenMerchant posts the Marx §1 linen-merchant
// fixture (£3,000 advanced, realisation function) and expects 201 + a Location
// header. The response must carry SurplusValueProduced == 0 (the central
// invariant of commercial capital).
func TestCreateCommercialCapital_LinenMerchant(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":300000,"commodity_description":"30,000 yards linen at 10p","function":0}`
	res, err := http.Post(ts.URL+"/v1/merchant/commercial-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/merchant/commercial-capital/") {
		t.Errorf("Location = %q, want /v1/merchant/commercial-capital/ prefix", loc)
	}

	var cc merchant.CommercialCapital
	if err := json.NewDecoder(res.Body).Decode(&cc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if cc.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if cc.CreatedAt.IsZero() {
		t.Error("CreatedAt not set by store")
	}
	// Central invariant: commercial capital creates no new surplus-value.
	if cc.SurplusValueProduced != 0 {
		t.Errorf("surplus_value_produced = %d, want 0", cc.SurplusValueProduced)
	}
	if cc.MoneyAdvanced != 300000 {
		t.Errorf("money_advanced = %d, want 300000", cc.MoneyAdvanced)
	}
	if cc.Function != merchant.FunctionRealisation {
		t.Errorf("function = %d, want FunctionRealisation (0)", cc.Function)
	}
}

// TestCreateCommercialCapital_BadJSON returns 400 on a malformed body.
func TestCreateCommercialCapital_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/merchant/commercial-capital", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateCommercialCapital_NonPositiveMoney returns 422 when money_advanced is 0
// (handler maps ErrNonPositiveMoney → 422 Unprocessable Entity).
func TestCreateCommercialCapital_NonPositiveMoney(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	for _, body := range []string{
		`{"money_advanced":0,"commodity_description":"linen","function":0}`,
		`{"money_advanced":-1,"commodity_description":"linen","function":0}`,
	} {
		res, err := http.Post(ts.URL+"/v1/merchant/commercial-capital", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("post: %v", err)
		}
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("body=%s: status = %d, want 422", body, res.StatusCode)
		}
		res.Body.Close()
	}
}

// --- GET /v1/merchant/commercial-capital ---

// TestListCommercialCapitals_NeverNull asserts the list endpoint returns 200 and
// an items array that is non-null even when the store is empty.
func TestListCommercialCapitals_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/merchant/commercial-capital")
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
	var decoded []merchant.CommercialCapital
	if err := json.Unmarshal(items, &decoded); err != nil {
		t.Fatalf("unmarshal items: %v", err)
	}
	if decoded == nil {
		t.Error("items decoded to nil, want non-nil slice")
	}
}

// --- GET /v1/merchant/commercial-capital/{id} ---

// TestGetCommercialCapital_NotFound returns 404 for an unknown id.
func TestGetCommercialCapital_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/merchant/commercial-capital/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateThenGetCommercialCapital exercises the create→GET round-trip: a known
// id returns 200 with the same money_advanced and SurplusValueProduced == 0.
func TestCreateThenGetCommercialCapital(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	// Create: Marx §2 storage example — £3,000 advanced to store linen stock.
	body := `{"money_advanced":300000,"commodity_description":"linen warehouse stock","function":1}`
	createRes, err := http.Post(ts.URL+"/v1/merchant/commercial-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer createRes.Body.Close()

	var created merchant.CommercialCapital
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	// GET the same record.
	getRes, err := http.Get(ts.URL + "/v1/merchant/commercial-capital/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", getRes.StatusCode)
	}

	var fetched merchant.CommercialCapital
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("fetched.ID = %s, want %s", fetched.ID, created.ID)
	}
	if fetched.MoneyAdvanced != 300000 {
		t.Errorf("fetched.MoneyAdvanced = %d, want 300000", fetched.MoneyAdvanced)
	}
	if fetched.SurplusValueProduced != 0 {
		t.Errorf("fetched.SurplusValueProduced = %d, want 0", fetched.SurplusValueProduced)
	}
	if fetched.Function != merchant.FunctionStorage {
		t.Errorf("fetched.Function = %d, want FunctionStorage (1)", fetched.Function)
	}
}
