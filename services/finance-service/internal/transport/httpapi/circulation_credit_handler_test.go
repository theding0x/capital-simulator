package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestCreateClearingHouseSettlement_London tests POST with the London Clearing
// House 1840s fixture: total_claims=15 000 000 000, money_used=1 000 000 000
// → net_balance=14 000 000 000, 201 with Location header set.
func TestCreateClearingHouseSettlement_London(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"London Clearing House 1840s","description":"Bills cancel; only net balance requires coin.","total_claims":15000000000,"money_used":1000000000}`
	res, err := http.Post(ts.URL+"/v1/credit/clearing-house-settlements", "application/json", strings.NewReader(body))
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

	var s credit.ClearingHouseSettlement
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if s.ID.IsZero() {
		t.Error("ID must be set")
	}
	if s.NetBalance != 14_000_000_000 {
		t.Errorf("NetBalance: got %d, want 14000000000", s.NetBalance)
	}
}

// TestCreateClearingHouseSettlement_BadJSON tests 400 on malformed JSON.
func TestCreateClearingHouseSettlement_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/clearing-house-settlements", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateClearingHouseSettlement_MoneyUsedExceedsClaims tests 422 when
// money_used > total_claims.
func TestCreateClearingHouseSettlement_MoneyUsedExceedsClaims(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","description":"","total_claims":1000,"money_used":1001}`
	res, err := http.Post(ts.URL+"/v1/credit/clearing-house-settlements", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestListClearingHouseSettlements_NeverNull verifies the list endpoint returns
// {"items":[]} not null when the store is empty.
func TestListClearingHouseSettlements_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/clearing-house-settlements")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.ClearingHouseSettlement `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}

// TestGetClearingHouseSettlement_NotFound tests 404 for an unknown ID.
func TestGetClearingHouseSettlement_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/clearing-house-settlements/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}

// TestGetClearingHouseSettlement_RoundTrip creates a settlement via POST then
// retrieves it via the Location header and verifies all fields match.
func TestGetClearingHouseSettlement_RoundTrip(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Scotland velocity-9","description":"£1 note passes through 9 hands in a week.","total_claims":900000,"money_used":100000}`
	createRes, err := http.Post(ts.URL+"/v1/credit/clearing-house-settlements", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer createRes.Body.Close()
	if createRes.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", createRes.StatusCode)
	}

	var created credit.ClearingHouseSettlement
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	getRes, err := http.Get(ts.URL + "/v1/credit/clearing-house-settlements/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status: got %d, want 200", getRes.StatusCode)
	}

	var fetched credit.ClearingHouseSettlement
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", fetched.ID, created.ID)
	}
	if fetched.Name != created.Name {
		t.Errorf("Name mismatch: got %q, want %q", fetched.Name, created.Name)
	}
	if fetched.TotalClaims != 900_000 {
		t.Errorf("TotalClaims: got %d, want 900000", fetched.TotalClaims)
	}
	if fetched.MoneyUsed != 100_000 {
		t.Errorf("MoneyUsed: got %d, want 100000", fetched.MoneyUsed)
	}
	if fetched.NetBalance != 800_000 {
		t.Errorf("NetBalance: got %d, want 800000", fetched.NetBalance)
	}
}
