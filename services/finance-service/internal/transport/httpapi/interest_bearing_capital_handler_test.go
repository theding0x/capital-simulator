package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// --- POST /v1/credit/interest-bearing-capital ---

// TestCreateInterestBearingCapital_Marx1000Pounds posts the canonical Marx
// fixture (Vol. III Ch. 21): £1,000 at 5% for one year.
//
//	MoneyAdvanced=100000 pence, InterestRateBP=500, LoanTermDays=365
//	InterestEarned = roundHalfUp(100000*500, 10000) = 5000
//	MoneyReturned  = 105000; NewValueCreated = 0
//
// Expects 201 + Location header with the assigned ID.
func TestCreateInterestBearingCapital_Marx1000Pounds(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":100000,"interest_rate_bp":500,"loan_term_days":365}`
	res, err := http.Post(ts.URL+"/v1/credit/interest-bearing-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/credit/interest-bearing-capital/") {
		t.Errorf("Location = %q, want /v1/credit/interest-bearing-capital/ prefix", loc)
	}

	var ibc credit.InterestBearingCapital
	if err := json.NewDecoder(res.Body).Decode(&ibc); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if ibc.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if ibc.CreatedAt.IsZero() {
		t.Error("CreatedAt not set by store")
	}
	if ibc.InterestEarned != 5000 {
		t.Errorf("interest_earned = %d, want 5000", ibc.InterestEarned)
	}
	if ibc.MoneyReturned != 105000 {
		t.Errorf("money_returned = %d, want 105000", ibc.MoneyReturned)
	}
	if ibc.NewValueCreated != 0 {
		t.Errorf("new_value_created = %d, want 0", ibc.NewValueCreated)
	}
}

// TestCreateInterestBearingCapital_BadJSON returns 400 on a malformed body.
func TestCreateInterestBearingCapital_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/interest-bearing-capital", "application/json", strings.NewReader("{nope"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateInterestBearingCapital_ZeroMoneyAdvanced returns 422 when
// money_advanced=0 (handler maps ErrNonPositiveMoney → 422).
func TestCreateInterestBearingCapital_ZeroMoneyAdvanced(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":0,"interest_rate_bp":500,"loan_term_days":365}`
	res, err := http.Post(ts.URL+"/v1/credit/interest-bearing-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for money_advanced=0", res.StatusCode)
	}
}

// TestCreateInterestBearingCapital_ZeroRate returns 422 when interest_rate_bp=0.
func TestCreateInterestBearingCapital_ZeroRate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":100000,"interest_rate_bp":0,"loan_term_days":365}`
	res, err := http.Post(ts.URL+"/v1/credit/interest-bearing-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for interest_rate_bp=0", res.StatusCode)
	}
}

// TestCreateInterestBearingCapital_ZeroTerm returns 422 when loan_term_days=0.
func TestCreateInterestBearingCapital_ZeroTerm(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"money_advanced":100000,"interest_rate_bp":500,"loan_term_days":0}`
	res, err := http.Post(ts.URL+"/v1/credit/interest-bearing-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for loan_term_days=0", res.StatusCode)
	}
}

// --- GET /v1/credit/interest-bearing-capital ---

// TestListInterestBearingCapitals_NeverNull asserts the list endpoint returns
// 200 and an items array that is non-null even when the store is empty.
func TestListInterestBearingCapitals_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/interest-bearing-capital")
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
	var decoded []credit.InterestBearingCapital
	if err := json.Unmarshal(items, &decoded); err != nil {
		t.Fatalf("unmarshal items: %v", err)
	}
	if decoded == nil {
		t.Error("items decoded to nil, want non-nil slice")
	}
}

// --- GET /v1/credit/interest-bearing-capital/{id} ---

// TestGetInterestBearingCapital_NotFound returns 404 for an unknown id.
func TestGetInterestBearingCapital_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/interest-bearing-capital/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}
