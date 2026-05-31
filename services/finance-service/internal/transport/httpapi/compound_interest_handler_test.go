package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestCreateCompoundInterestSchedule_DrPrice tests POST /v1/credit/compound-interest
// using Marx's Dr. Price fixture: 10000 units at 5% (500 bp) for 100 periods
// (principal scaled for integer precision; conceptual penny example).
func TestCreateCompoundInterestSchedule_DrPrice(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"principal":10000,"rate_bp":500,"period_years":100,"fetish_form":2}`
	res, err := http.Post(ts.URL+"/v1/credit/compound-interest", "application/json", strings.NewReader(body))
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

	var sched credit.CompoundInterestSchedule
	if err := json.NewDecoder(res.Body).Decode(&sched); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if sched.ID.IsZero() {
		t.Error("ID must be set")
	}
	if sched.NewValueCreated != 0 {
		t.Errorf("NewValueCreated: got %d, want 0", sched.NewValueCreated)
	}
	if sched.FinalValue <= sched.Principal {
		t.Errorf("FinalValue %d must exceed Principal %d", sched.FinalValue, sched.Principal)
	}
}

// TestCreateCompoundInterestSchedule_BadJSON tests 400 on malformed JSON.
func TestCreateCompoundInterestSchedule_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/compound-interest", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateCompoundInterestSchedule_ZeroPrincipal tests 422 on principal=0.
func TestCreateCompoundInterestSchedule_ZeroPrincipal(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"principal":0,"rate_bp":500,"period_years":10,"fetish_form":2}`
	res, err := http.Post(ts.URL+"/v1/credit/compound-interest", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestListCompoundInterestSchedules_NeverNull verifies the list endpoint
// returns {"items":[]} rather than null when the store is empty.
func TestListCompoundInterestSchedules_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/compound-interest")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.CompoundInterestSchedule `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}

// TestGetCompoundInterestSchedule_NotFound tests 404 for an unknown ID.
func TestGetCompoundInterestSchedule_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/compound-interest/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}
