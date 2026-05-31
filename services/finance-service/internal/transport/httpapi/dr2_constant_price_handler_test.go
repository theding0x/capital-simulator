// Package httpapi — handler tests for Vol. III Ch. 41 DR II Constant Price
// endpoints. Fixtures drawn from Marx Vol. III Ch. 41 Table II: Soil D,
// two successive tranches each 250bp capital, 900bp surplus-profit →
// TotalRentalBP=1800, RateOfSurplusProfitBP=36000.
package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"log/slog"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

// newDR2ConstantTestServer builds a minimal httptest.Server that registers
// only the Vol. III Ch. 41 DR II constant-price routes.
func newDR2ConstantTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/dr2-outcomes", h.CreateDR2Outcome)
	mux.HandleFunc("GET /v1/rent/dr2-outcomes", h.ListDR2Outcomes)
	mux.HandleFunc("POST /v1/rent/intensive-extensive", h.CreateIntensiveExtensiveComparison)
	mux.HandleFunc("GET /v1/rent/intensive-extensive", h.ListIntensiveExtensiveComparisons)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// marxDR2TwoTranches is the canonical body for a Case II (proportional)
// request with two Soil-D tranches from Marx Ch. 41 Table II.
const marxDR2TwoTranches = `{
  "case": 2,
  "investments": [
    {"parcel_id":"PD","tranche_number":1,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":900},
    {"parcel_id":"PD","tranche_number":2,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":900}
  ]
}`

// ── POST /v1/rent/dr2-outcomes — happy path ───────────────────────────────────

func TestCreateDR2Outcome_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2ConstantTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/dr2-outcomes", "application/json",
		strings.NewReader(marxDR2TwoTranches))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.DifferentialRentIIOutcome
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.TotalRentalBP != 1800 {
		t.Errorf("total_rental_bp = %d, want 1800", saved.TotalRentalBP)
	}
	if saved.RateOfSurplusProfitBP != 36000 {
		t.Errorf("rate_of_surplus_profit_bp = %d, want 36000", saved.RateOfSurplusProfitBP)
	}
}

// ── POST /v1/rent/dr2-outcomes — validation failures ─────────────────────────

func TestCreateDR2Outcome_CaseZero_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2ConstantTestServer(t)

	body := `{"case":0,"investments":[{"parcel_id":"PD","tranche_number":1,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":900}]}`
	res, err := http.Post(ts.URL+"/v1/rent/dr2-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for case=0", res.StatusCode)
	}
}

func TestCreateDR2Outcome_CaseFive_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2ConstantTestServer(t)

	body := `{"case":5,"investments":[{"parcel_id":"PD","tranche_number":1,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":900}]}`
	res, err := http.Post(ts.URL+"/v1/rent/dr2-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for case=5", res.StatusCode)
	}
}

func TestCreateDR2Outcome_EmptyInvestments_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2ConstantTestServer(t)

	body := `{"case":2,"investments":[]}`
	res, err := http.Post(ts.URL+"/v1/rent/dr2-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty investments", res.StatusCode)
	}
}

func TestCreateDR2Outcome_MalformedJSON_Returns400(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2ConstantTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/dr2-outcomes", "application/json",
		strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", res.StatusCode)
	}
}

// ── GET /v1/rent/dr2-outcomes ─────────────────────────────────────────────────

func TestListDR2Outcomes_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2ConstantTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/dr2-outcomes")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.DifferentialRentIIOutcome `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── POST /v1/rent/intensive-extensive — happy path ───────────────────────────

// TestCreateIntensiveExtensive_HappyPath uses Marx Ch. 41 values:
// intensive 1800bp > extensive 900bp for same total rental 1800bp,
// total capital 500bp.
func TestCreateIntensiveExtensive_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2ConstantTestServer(t)

	body := `{"intensive_rent_per_acre_bp":1800,"extensive_rent_per_acre_bp":900,"total_capital_bp":500,"total_rental_same_bp":1800}`
	res, err := http.Post(ts.URL+"/v1/rent/intensive-extensive", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.IntensiveExtensiveComparison
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
}

// TestCreateIntensiveExtensive_IntensiveNotGreaterThanExtensive verifies that
// equal intensive and extensive rents (both 900bp) return 422.
func TestCreateIntensiveExtensive_IntensiveNotGreaterThanExtensive_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2ConstantTestServer(t)

	body := `{"intensive_rent_per_acre_bp":900,"extensive_rent_per_acre_bp":900,"total_capital_bp":500,"total_rental_same_bp":1800}`
	res, err := http.Post(ts.URL+"/v1/rent/intensive-extensive", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 when intensive <= extensive", res.StatusCode)
	}
}

// ── GET /v1/rent/intensive-extensive ─────────────────────────────────────────

func TestListIntensiveExtensive_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2ConstantTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/intensive-extensive")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.IntensiveExtensiveComparison `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
