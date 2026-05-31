package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestCreateMoneyCapitalAccumulation_Inactivity tests POST with the post-1847
// inactivity fixture: phase=1 (Inactivity), rate=100bp, abundant loanable supply.
func TestCreateMoneyCapitalAccumulation_Inactivity(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Post-1847 Inactivity","phase":1,"interest_rate_bp":100,"period":"1847 crisis","loanable_amount":5000000,"loanable_abundant":true,"idle_amount":2000000,"idle_description":"idle cotton-mill capital","gap_bp":500,"gap_description":"money-capital accumulates faster","description":"post-crisis stagnation"}`
	res, err := http.Post(ts.URL+"/v1/credit/money-capital-accumulation", "application/json", strings.NewReader(body))
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

	var a credit.MoneyCapitalAccumulation
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.ID.IsZero() {
		t.Error("ID must be set")
	}
	if a.CycleObservation.Phase != credit.PhaseInactivity {
		t.Errorf("Phase: got %d, want %d", a.CycleObservation.Phase, credit.PhaseInactivity)
	}
	if a.CycleObservation.InterestRateBP != 100 {
		t.Errorf("InterestRateBP: got %d, want 100", a.CycleObservation.InterestRateBP)
	}
	if a.CycleObservation.Period != "1847 crisis" {
		t.Errorf("Period: got %q, want %q", a.CycleObservation.Period, "1847 crisis")
	}
	if !a.LoanableSupply.Abundant {
		t.Error("LoanableSupply.Abundant must be true")
	}
}

// TestCreateMoneyCapitalAccumulation_Crisis tests POST with the 1857 crisis fixture:
// phase=5, rate=1000bp, scarce loanable capital, no idle capital.
func TestCreateMoneyCapitalAccumulation_Crisis(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Crisis Peak 1857","phase":5,"interest_rate_bp":1000,"loanable_amount":500000,"loanable_abundant":false,"idle_amount":0,"idle_description":"no idle capital","gap_bp":-800,"gap_description":"productive capital outpaces loanable supply","description":"1857 crisis"}`
	res, err := http.Post(ts.URL+"/v1/credit/money-capital-accumulation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status: got %d, want 201", res.StatusCode)
	}

	var a credit.MoneyCapitalAccumulation
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.CycleObservation.Phase != credit.PhaseCrisis {
		t.Errorf("Phase: got %d, want %d", a.CycleObservation.Phase, credit.PhaseCrisis)
	}
	if a.Gap.GapBP != -800 {
		t.Errorf("GapBP: got %d, want -800", a.Gap.GapBP)
	}
}

// TestGetMoneyCapitalAccumulation_OK tests GET /{id} returns 200 for an existing record.
func TestGetMoneyCapitalAccumulation_OK(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	// Create first.
	body := `{"name":"Revival","phase":2,"interest_rate_bp":300,"loanable_amount":4000000,"loanable_abundant":true,"idle_amount":500000,"idle_description":"re-engaged","gap_bp":100,"gap_description":"gap narrows","description":"revival"}`
	res, err := http.Post(ts.URL+"/v1/credit/money-capital-accumulation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", res.StatusCode)
	}

	var created credit.MoneyCapitalAccumulation
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	// Get by ID.
	getRes, err := http.Get(ts.URL + "/v1/credit/money-capital-accumulation/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status: got %d, want 200", getRes.StatusCode)
	}

	var fetched credit.MoneyCapitalAccumulation
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", fetched.ID, created.ID)
	}
}

// TestGetMoneyCapitalAccumulation_NotFound tests 404 for an unknown ID.
func TestGetMoneyCapitalAccumulation_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/money-capital-accumulation/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}

// TestCreateMoneyCapitalAccumulation_BadJSON tests 400 on malformed JSON.
func TestCreateMoneyCapitalAccumulation_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/money-capital-accumulation", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateMoneyCapitalAccumulation_InvalidPhase tests 422 on phase=0.
func TestCreateMoneyCapitalAccumulation_InvalidPhase(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"x","phase":0,"interest_rate_bp":100,"loanable_amount":100,"loanable_abundant":false,"idle_amount":0,"idle_description":"","gap_bp":0,"gap_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/money-capital-accumulation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateMoneyCapitalAccumulation_NegativeRate tests 422 on interest_rate_bp=-1.
func TestCreateMoneyCapitalAccumulation_NegativeRate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"x","phase":1,"interest_rate_bp":-1,"loanable_amount":100,"loanable_abundant":false,"idle_amount":0,"idle_description":"","gap_bp":0,"gap_description":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/money-capital-accumulation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestListMoneyCapitalAccumulations_NeverNull verifies the list endpoint returns
// {"items":[]} rather than null when the store is empty.
func TestListMoneyCapitalAccumulations_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/money-capital-accumulation")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.MoneyCapitalAccumulation `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}
