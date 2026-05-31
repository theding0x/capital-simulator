package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestCreateRealCapitalAccumulation_Crisis1847 tests POST with the 1847 crisis
// fixture: phase=5 (Crisis), commercial credit contracts, cause is NOT money
// scarcity, money-/real-capital growth diverge.
func TestCreateRealCapitalAccumulation_Crisis1847(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Crisis 1847","phase":5,"money_capital_growth_bp":800,"real_capital_growth_bp":-600,"corresponds":false,"correspondence_explanation":"loanable money piles up while real capital contracts","reserve_capital":500000,"credit_limit_description":"Bank of England reserve","commercial_credit_contracts":true,"cause_is_money_scarcity":false,"description":"blocked returns and lost confidence"}`
	res, err := http.Post(ts.URL+"/v1/credit/real-capital-accumulation", "application/json", strings.NewReader(body))
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

	var a credit.RealCapitalAccumulation
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.ID.IsZero() {
		t.Error("ID must be set")
	}
	if a.Phase != credit.PhaseCrisis {
		t.Errorf("Phase: got %d, want %d", a.Phase, credit.PhaseCrisis)
	}
	if !a.CommercialCreditContracts {
		t.Error("CommercialCreditContracts must be true")
	}
	if a.CauseIsMoneyScarcity {
		t.Error("CauseIsMoneyScarcity must be false")
	}
	if a.Correspondence.MoneyCapitalGrowthBP != 800 {
		t.Errorf("MoneyCapitalGrowthBP: got %d, want 800", a.Correspondence.MoneyCapitalGrowthBP)
	}
	if a.Correspondence.RealCapitalGrowthBP != -600 {
		t.Errorf("RealCapitalGrowthBP: got %d, want -600", a.Correspondence.RealCapitalGrowthBP)
	}
}

// TestGetRealCapitalAccumulation_OK tests GET /{id} returns 200 for an existing
// record and round-trips a correspondence sub-field (proving Memory↔MySQL parity
// of the embedded value objects).
func TestGetRealCapitalAccumulation_OK(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Revival","phase":2,"money_capital_growth_bp":400,"real_capital_growth_bp":450,"corresponds":true,"correspondence_explanation":"funds and investment expand together","reserve_capital":3000000,"credit_limit_description":"ample reserve","commercial_credit_contracts":false,"cause_is_money_scarcity":false,"description":"revival"}`
	res, err := http.Post(ts.URL+"/v1/credit/real-capital-accumulation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", res.StatusCode)
	}

	var created credit.RealCapitalAccumulation
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	getRes, err := http.Get(ts.URL + "/v1/credit/real-capital-accumulation/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status: got %d, want 200", getRes.StatusCode)
	}

	var fetched credit.RealCapitalAccumulation
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", fetched.ID, created.ID)
	}
	// Round-trip the embedded correspondence and credit-limit sub-fields.
	if fetched.Correspondence.RealCapitalGrowthBP != 450 {
		t.Errorf("RealCapitalGrowthBP round-trip: got %d, want 450", fetched.Correspondence.RealCapitalGrowthBP)
	}
	if !fetched.Correspondence.Corresponds {
		t.Error("Corresponds round-trip: got false, want true")
	}
	if fetched.CreditLimit.ReserveCapital != 3000000 {
		t.Errorf("ReserveCapital round-trip: got %d, want 3000000", fetched.CreditLimit.ReserveCapital)
	}
}

// TestGetRealCapitalAccumulation_NotFound tests 404 for an unknown ID.
func TestGetRealCapitalAccumulation_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/real-capital-accumulation/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}

// TestCreateRealCapitalAccumulation_BadJSON tests 400 on malformed JSON.
func TestCreateRealCapitalAccumulation_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/real-capital-accumulation", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateRealCapitalAccumulation_InvalidPhase tests 422 on phase=0.
func TestCreateRealCapitalAccumulation_InvalidPhase(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"x","phase":0,"money_capital_growth_bp":0,"real_capital_growth_bp":0,"corresponds":false,"correspondence_explanation":"","reserve_capital":0,"credit_limit_description":"","commercial_credit_contracts":false,"cause_is_money_scarcity":false,"description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/real-capital-accumulation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateRealCapitalAccumulation_NegativeReserve tests 422 on reserve_capital=-1.
func TestCreateRealCapitalAccumulation_NegativeReserve(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"x","phase":1,"money_capital_growth_bp":0,"real_capital_growth_bp":0,"corresponds":false,"correspondence_explanation":"","reserve_capital":-1,"credit_limit_description":"","commercial_credit_contracts":false,"cause_is_money_scarcity":false,"description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/real-capital-accumulation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestListRealCapitalAccumulations_NeverNull verifies the list endpoint returns
// {"items":[]} rather than null when the store is empty.
func TestListRealCapitalAccumulations_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/real-capital-accumulation")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.RealCapitalAccumulation `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}
