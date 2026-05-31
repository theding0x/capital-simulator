package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// ── POST /v1/credit/gold-reserves ────────────────────────────────────────────

// TestCreateGoldReserve_BoE1857 verifies the BoE Nov 1857 crisis fixture:
// amount=800 000 000 pence → 201, Location set, body amount==800000000.
func TestCreateGoldReserve_BoE1857(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"BoE Nov 1857","amount":800000000,"description":"Bank of England gold reserve during the 1857 crisis."}`
	res, err := http.Post(ts.URL+"/v1/credit/gold-reserves", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status: got %d, want 201", res.StatusCode)
	}
	if res.Header.Get("Location") == "" {
		t.Error("expected Location header")
	}

	var g credit.GoldReserve
	if err := json.NewDecoder(res.Body).Decode(&g); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if g.ID.IsZero() {
		t.Error("ID must be set")
	}
	if g.Amount != 800_000_000 {
		t.Errorf("amount: got %d, want 800000000", g.Amount)
	}
}

// TestCreateGoldReserve_BadJSON tests 400 on malformed JSON.
func TestCreateGoldReserve_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/gold-reserves", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateGoldReserve_NegativeAmount tests 422 on negative amount.
func TestCreateGoldReserve_NegativeAmount(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","amount":-1,"description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/gold-reserves", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// ── GET /v1/credit/gold-reserves ─────────────────────────────────────────────

// TestListGoldReserves_NeverNull verifies the list endpoint returns
// {"items":[]} not null when the store is empty.
func TestListGoldReserves_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/gold-reserves")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.GoldReserve `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}

// ── GET /v1/credit/gold-reserves/{id} ────────────────────────────────────────

// TestGetGoldReserve_NotFound tests 404 for an unknown ID.
func TestGetGoldReserve_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/gold-reserves/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}

// ── POST /v1/credit/rates-of-exchange ────────────────────────────────────────

// TestCreateRateOfExchange_LondonHamburg1847 verifies the London-Hamburg 1847
// adverse-balance fixture: deviation_bp=−150 → 201, Location set, body deviation_bp==-150.
func TestCreateRateOfExchange_LondonHamburg1847(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"London-Hamburg 1847","deviation_bp":-150,"description":"Adverse balance from corn imports."}`
	res, err := http.Post(ts.URL+"/v1/credit/rates-of-exchange", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status: got %d, want 201", res.StatusCode)
	}
	if res.Header.Get("Location") == "" {
		t.Error("expected Location header")
	}

	var re credit.RateOfExchange
	if err := json.NewDecoder(res.Body).Decode(&re); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if re.ID.IsZero() {
		t.Error("ID must be set")
	}
	if re.DeviationBP != -150 {
		t.Errorf("deviation_bp: got %d, want -150", re.DeviationBP)
	}
}

// TestCreateRateOfExchange_BadJSON tests 400 on malformed JSON.
func TestCreateRateOfExchange_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/rates-of-exchange", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// ── GET /v1/credit/rates-of-exchange ─────────────────────────────────────────

// TestListRatesOfExchange_NeverNull verifies the list endpoint returns
// {"items":[]} not null when the store is empty.
func TestListRatesOfExchange_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/rates-of-exchange")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.RateOfExchange `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}

// ── GET /v1/credit/rates-of-exchange/{id} ────────────────────────────────────

// TestGetRateOfExchange_NotFound tests 404 for an unknown ID.
func TestGetRateOfExchange_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/rates-of-exchange/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}

// TestRateOfExchange_PostThenGetByID verifies that a created record can be
// retrieved via its ID with deviation_bp intact.
func TestRateOfExchange_PostThenGetByID(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"London-India 1857","deviation_bp":-150,"description":"Persistent adverse drain."}`
	createRes, err := http.Post(ts.URL+"/v1/credit/rates-of-exchange", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer createRes.Body.Close()
	if createRes.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", createRes.StatusCode)
	}

	var created credit.RateOfExchange
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	getRes, err := http.Get(ts.URL + "/v1/credit/rates-of-exchange/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status: got %d, want 200", getRes.StatusCode)
	}

	var fetched credit.RateOfExchange
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", fetched.ID, created.ID)
	}
	if fetched.DeviationBP != -150 {
		t.Errorf("deviation_bp: got %d, want -150", fetched.DeviationBP)
	}
}
