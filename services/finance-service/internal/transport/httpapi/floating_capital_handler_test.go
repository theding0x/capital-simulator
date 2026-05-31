package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestCreateFloatingCapital_Ipswich tests POST with the Ipswich-farmers fixture:
// source=2 (capital/revenue conversion, accompanies real accumulation), and the
// £20-lent-5-times velocity (2000 × 5 = 10000).
func TestCreateFloatingCapital_Ipswich(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Ipswich farmers","amount":2000000,"source":2,"velocity_money_amount":2000,"velocity_times_lent":5,"reserve_saving_amount":50000,"reserve_saving_description":"country-bank till economies","description":"rural surplus rediscounted to manufacturing"}`
	res, err := http.Post(ts.URL+"/v1/credit/floating-capital", "application/json", strings.NewReader(body))
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

	var f credit.FloatingCapital
	if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if f.ID.IsZero() {
		t.Error("ID must be set")
	}
	if f.Source != credit.SourceCapitalRevenueConversion {
		t.Errorf("Source: got %d, want %d", f.Source, credit.SourceCapitalRevenueConversion)
	}
	// loan_capital_created is derived server-side: 2000 × 5 = 10000.
	if f.Velocity.LoanCapitalCreated != 10000 {
		t.Errorf("Velocity.LoanCapitalCreated: got %d, want 10000", f.Velocity.LoanCapitalCreated)
	}
}

// TestGetFloatingCapital_OK tests GET /{id} returns 200 and round-trips the
// embedded velocity + reserve-saving sub-fields (proving Memory↔MySQL parity).
func TestGetFloatingCapital_OK(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Post-crisis idle","amount":3000000,"source":1,"velocity_money_amount":100000,"velocity_times_lent":1,"reserve_saving_amount":2000000,"reserve_saving_description":"banking economies","description":"money lying fallow after a crisis"}`
	res, err := http.Post(ts.URL+"/v1/credit/floating-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", res.StatusCode)
	}

	var created credit.FloatingCapital
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	getRes, err := http.Get(ts.URL + "/v1/credit/floating-capital/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status: got %d, want 200", getRes.StatusCode)
	}

	var fetched credit.FloatingCapital
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", fetched.ID, created.ID)
	}
	// Round-trip every embedded sub-field — Memory↔MySQL parity.
	if fetched.Velocity.MoneyAmount != 100000 {
		t.Errorf("Velocity.MoneyAmount round-trip: got %d, want 100000", fetched.Velocity.MoneyAmount)
	}
	if fetched.Velocity.TimesLent != 1 {
		t.Errorf("Velocity.TimesLent round-trip: got %d, want 1", fetched.Velocity.TimesLent)
	}
	if fetched.Velocity.LoanCapitalCreated != 100000 {
		t.Errorf("Velocity.LoanCapitalCreated round-trip: got %d, want 100000", fetched.Velocity.LoanCapitalCreated)
	}
	if fetched.ReserveSaving.Amount != 2000000 {
		t.Errorf("ReserveSaving.Amount round-trip: got %d, want 2000000", fetched.ReserveSaving.Amount)
	}
	if fetched.ReserveSaving.Description != "banking economies" {
		t.Errorf("ReserveSaving.Description round-trip: got %q, want %q", fetched.ReserveSaving.Description, "banking economies")
	}
}

// TestGetFloatingCapital_NotFound tests 404 for an unknown ID.
func TestGetFloatingCapital_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/floating-capital/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}

// TestCreateFloatingCapital_BadJSON tests 400 on malformed JSON.
func TestCreateFloatingCapital_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/floating-capital", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateFloatingCapital_InvalidSource tests 422 on source=3.
func TestCreateFloatingCapital_InvalidSource(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"x","amount":100,"source":3,"velocity_money_amount":2000,"velocity_times_lent":5,"reserve_saving_amount":0,"reserve_saving_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/floating-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateFloatingCapital_NegativeAmount tests 422 on amount=-1.
func TestCreateFloatingCapital_NegativeAmount(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"x","amount":-1,"source":1,"velocity_money_amount":2000,"velocity_times_lent":5,"reserve_saving_amount":0,"reserve_saving_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/floating-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateFloatingCapital_TimesLentTooLow tests 422 on velocity_times_lent=0.
func TestCreateFloatingCapital_TimesLentTooLow(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"x","amount":100,"source":1,"velocity_money_amount":2000,"velocity_times_lent":0,"reserve_saving_amount":0,"reserve_saving_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/floating-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestListFloatingCapitals_NeverNull verifies the list endpoint returns
// {"items":[]} rather than null when the store is empty.
func TestListFloatingCapitals_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/floating-capital")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.FloatingCapital `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}
