package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// ── POST /v1/credit/currency-observations ────────────────────────────────────

// TestCreateCurrencyObservation_BoE1844 tests POST with the Bank of England 1844
// fixture: £20m outstanding, £8m reserve, £14m coin, £4m capital-transfer.
func TestCreateCurrencyObservation_BoE1844(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name":"Bank of England 1844",
		"total_currency_outstanding":2000000000,
		"coin_function_amount":1400000000,
		"capital_transfer_amount":400000000,
		"reserve_amount":800000000,
		"reserve_description":"Banking Department reserve 1844",
		"description":"Bank Charter Act 1844"
	}`
	res, err := http.Post(ts.URL+"/v1/credit/currency-observations", "application/json", strings.NewReader(body))
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

	var o credit.CurrencyObservation
	if err := json.NewDecoder(res.Body).Decode(&o); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if o.ID.IsZero() {
		t.Error("ID must be set")
	}
	if o.TotalCurrencyOutstanding != 2_000_000_000 {
		t.Errorf("Total: got %d, want 2000000000", o.TotalCurrencyOutstanding)
	}
	if o.CoinFunctionAmount != 1_400_000_000 {
		t.Errorf("Coin: got %d, want 1400000000", o.CoinFunctionAmount)
	}
	if o.CapitalTransferAmount != 400_000_000 {
		t.Errorf("CapitalTransfer: got %d, want 400000000", o.CapitalTransferAmount)
	}
	// reserve round-trip — proves Memory↔MySQL parity for nested struct
	if o.ReserveFund.Amount != 800_000_000 {
		t.Errorf("ReserveFund.Amount: got %d, want 800000000", o.ReserveFund.Amount)
	}
	if o.ReserveFund.Description != "Banking Department reserve 1844" {
		t.Errorf("ReserveFund.Description: got %q", o.ReserveFund.Description)
	}
}

// TestGetCurrencyObservation_OK tests GET /{id} returns 200 for an existing record.
func TestGetCurrencyObservation_OK(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Scottish Banking","total_currency_outstanding":1000000000,"coin_function_amount":100000000,"capital_transfer_amount":700000000,"reserve_amount":5000000,"reserve_description":"minimal metallic reserve","description":"credit velocity"}`
	res, err := http.Post(ts.URL+"/v1/credit/currency-observations", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", res.StatusCode)
	}

	var created credit.CurrencyObservation
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	getRes, err := http.Get(ts.URL + "/v1/credit/currency-observations/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status: got %d, want 200", getRes.StatusCode)
	}

	var fetched credit.CurrencyObservation
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", fetched.ID, created.ID)
	}
	// Full reserve round-trip.
	if fetched.ReserveFund.Amount != 5_000_000 {
		t.Errorf("ReserveFund.Amount: got %d, want 5000000", fetched.ReserveFund.Amount)
	}
	if fetched.ReserveFund.Description != "minimal metallic reserve" {
		t.Errorf("ReserveFund.Description: got %q", fetched.ReserveFund.Description)
	}
}

// TestGetCurrencyObservation_NotFound tests 404 for an unknown ID.
func TestGetCurrencyObservation_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/currency-observations/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}

// TestCreateCurrencyObservation_BadJSON tests 400 on malformed JSON.
func TestCreateCurrencyObservation_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/currency-observations", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateCurrencyObservation_NonPositiveTotal tests 422 on total <= 0.
func TestCreateCurrencyObservation_NonPositiveTotal(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","total_currency_outstanding":0,"coin_function_amount":0,"capital_transfer_amount":0,"reserve_amount":0,"reserve_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/currency-observations", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateCurrencyObservation_SplitExceedsTotal tests 422 when coin+capital > total.
func TestCreateCurrencyObservation_SplitExceedsTotal(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","total_currency_outstanding":1000,"coin_function_amount":700,"capital_transfer_amount":400,"reserve_amount":0,"reserve_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/currency-observations", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateCurrencyObservation_NegativeReserve tests 422 on reserve_amount < 0.
func TestCreateCurrencyObservation_NegativeReserve(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","total_currency_outstanding":1000,"coin_function_amount":100,"capital_transfer_amount":200,"reserve_amount":-1,"reserve_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/currency-observations", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestListCurrencyObservations_NeverNull verifies the list endpoint returns
// {"items":[]} rather than null when the store is empty.
func TestListCurrencyObservations_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/currency-observations")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.CurrencyObservation `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}
