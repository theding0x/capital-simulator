package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestCreateCapitalRelease_PriceFall tests POST with the raw-material-price-fall
// fixture: cause=1 (input price fall), not reinvested → overstates accumulation,
// overstatement_pct=2500.
func TestCreateCapitalRelease_PriceFall(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Capital released by price fall","released_amount":5000000,"cause":1,"reinvested":false,"overstatement_pct":2500,"overstatement_description":"freed money lent out, not reinvested","description":"raw-material prices fall; the difference is lent out"}`
	res, err := http.Post(ts.URL+"/v1/credit/capital-release", "application/json", strings.NewReader(body))
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

	var c credit.CapitalRelease
	if err := json.NewDecoder(res.Body).Decode(&c); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if c.ID.IsZero() {
		t.Error("ID must be set")
	}
	if c.Cause != credit.CauseInputPriceFall {
		t.Errorf("Cause: got %d, want %d", c.Cause, credit.CauseInputPriceFall)
	}
	if c.Overstatement.OverstatementPct != 2500 {
		t.Errorf("Overstatement.OverstatementPct: got %d, want 2500", c.Overstatement.OverstatementPct)
	}
}

// TestGetCapitalRelease_OK tests GET /{id} returns 200 and round-trips the
// embedded overstatement sub-fields (proving Memory↔MySQL parity).
func TestGetCapitalRelease_OK(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Rentier savings","released_amount":2000000,"cause":3,"reinvested":false,"overstatement_pct":10000,"overstatement_description":"the whole sum is apparent accumulation","description":"rentier class places savings at interest"}`
	res, err := http.Post(ts.URL+"/v1/credit/capital-release", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", res.StatusCode)
	}

	var created credit.CapitalRelease
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	getRes, err := http.Get(ts.URL + "/v1/credit/capital-release/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status: got %d, want 200", getRes.StatusCode)
	}

	var fetched credit.CapitalRelease
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", fetched.ID, created.ID)
	}
	// Round-trip every embedded sub-field — Memory↔MySQL parity.
	if fetched.Overstatement.OverstatementPct != 10000 {
		t.Errorf("Overstatement.OverstatementPct round-trip: got %d, want 10000", fetched.Overstatement.OverstatementPct)
	}
	if fetched.Overstatement.Description != "the whole sum is apparent accumulation" {
		t.Errorf("Overstatement.Description round-trip: got %q, want %q", fetched.Overstatement.Description, "the whole sum is apparent accumulation")
	}
	if fetched.Cause != credit.CauseRentierSavings {
		t.Errorf("Cause round-trip: got %d, want %d", fetched.Cause, credit.CauseRentierSavings)
	}
}

// TestGetCapitalRelease_NotFound tests 404 for an unknown ID.
func TestGetCapitalRelease_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/capital-release/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}

// TestCreateCapitalRelease_BadJSON tests 400 on malformed JSON.
func TestCreateCapitalRelease_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/capital-release", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateCapitalRelease_InvalidCause tests 422 on cause=5.
func TestCreateCapitalRelease_InvalidCause(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"x","released_amount":100,"cause":5,"reinvested":false,"overstatement_pct":0,"overstatement_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/capital-release", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateCapitalRelease_NegativeAmount tests 422 on released_amount=-1.
func TestCreateCapitalRelease_NegativeAmount(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"x","released_amount":-1,"cause":1,"reinvested":false,"overstatement_pct":0,"overstatement_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/capital-release", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateCapitalRelease_NegativeOverstatement tests 422 on overstatement_pct=-1.
func TestCreateCapitalRelease_NegativeOverstatement(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"x","released_amount":100,"cause":1,"reinvested":false,"overstatement_pct":-1,"overstatement_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/capital-release", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestListCapitalReleases_NeverNull verifies the list endpoint returns
// {"items":[]} rather than null when the store is empty.
func TestListCapitalReleases_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/capital-release")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.CapitalRelease `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}
