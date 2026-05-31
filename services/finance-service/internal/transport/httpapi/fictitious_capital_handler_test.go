package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestCreateBillOfExchange_1839Bills tests POST /v1/credit/bills-of-exchange
// with Marx's 1839 UK bills fixture: £132m face value, 90-day maturity.
func TestCreateBillOfExchange_1839Bills(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"drawer":"Manchester Spinner","drawee":"London Merchant","face_value":13212346000,"maturity_days":90,"description":"1839 UK bill"}`
	res, err := http.Post(ts.URL+"/v1/credit/bills-of-exchange", "application/json", strings.NewReader(body))
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

	var bill credit.BillOfExchange
	if err := json.NewDecoder(res.Body).Decode(&bill); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if bill.ID.IsZero() {
		t.Error("ID must be set")
	}
	if bill.FaceValue != 13212346000 {
		t.Errorf("FaceValue: got %d, want 13212346000", bill.FaceValue)
	}
	if bill.MaturityDays != 90 {
		t.Errorf("MaturityDays: got %d, want 90", bill.MaturityDays)
	}
}

// TestCreateBillOfExchange_BadJSON tests 400 on malformed JSON.
func TestCreateBillOfExchange_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/bills-of-exchange", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateBillOfExchange_ZeroMaturity tests 422 on maturity_days=0.
func TestCreateBillOfExchange_ZeroMaturity(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"drawer":"A","drawee":"B","face_value":1000,"maturity_days":0,"description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/bills-of-exchange", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestListBillsOfExchange_NeverNull verifies the list endpoint returns
// {"items":[]} rather than null when the store is empty.
func TestListBillsOfExchange_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/bills-of-exchange")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.BillOfExchange `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}

// TestCreateFictitiousCapital_Consol tests POST /v1/credit/fictitious-capital
// with Marx's government consol fixture: £5 income at 5% → £100 capitalised,
// real_capital_exists=false (the original loan is consumed).
func TestCreateFictitiousCapital_Consol(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"kind":3,"name":"UK Consol","annual_income":500,"rate_bp":500,"description":"government bond"}`
	res, err := http.Post(ts.URL+"/v1/credit/fictitious-capital", "application/json", strings.NewReader(body))
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

	var fc credit.FictitiousCapital
	if err := json.NewDecoder(res.Body).Decode(&fc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if fc.CapitalisedValue != 10000 {
		t.Errorf("CapitalisedValue: got %d, want 10000", fc.CapitalisedValue)
	}
	if fc.RealCapitalExists {
		t.Error("RealCapitalExists must be false for public debt")
	}
}

// TestCreateFictitiousCapital_ZeroRate tests 422 on rate_bp=0.
func TestCreateFictitiousCapital_ZeroRate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"kind":4,"name":"x","annual_income":1000,"rate_bp":0,"description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/fictitious-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestListFictitiousCapitals_NeverNull verifies the list endpoint returns
// {"items":[]} rather than null when the store is empty.
func TestListFictitiousCapitals_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/fictitious-capital")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.FictitiousCapital `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}

// TestGetFictitiousCapital_NotFound tests 404 for an unknown ID.
func TestGetFictitiousCapital_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/fictitious-capital/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}
