package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestCreateRateOfInterest_Created asserts 201 + Location on valid input.
// Fixture: Prosperity 1843 — avg 2000 bp, interest 150 bp.
func TestCreateRateOfInterest_Created(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	body, _ := json.Marshal(map[string]any{
		"rate_bp":                150,
		"average_profit_rate_bp": 2000,
		"cycle_phase":            int(credit.PhaseProsperity),
		"period":                 "1843",
	})
	resp, err := http.Post(ts.URL+"/v1/credit/rate-of-interest", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d, want 201", resp.StatusCode)
	}
	if loc := resp.Header.Get("Location"); loc == "" {
		t.Error("Location header missing")
	}
}

// TestCreateRateOfInterest_BadJSON asserts 400 on malformed JSON body.
func TestCreateRateOfInterest_BadJSON(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	resp, err := http.Post(ts.URL+"/v1/credit/rate-of-interest", "application/json", bytes.NewReader([]byte(`{bad`)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

// TestCreateRateOfInterest_ValidationError asserts 422 when rate_bp >= average_profit_rate_bp.
func TestCreateRateOfInterest_ValidationError(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	body, _ := json.Marshal(map[string]any{
		"rate_bp":                2000,
		"average_profit_rate_bp": 2000,
		"cycle_phase":            int(credit.PhaseProsperity),
		"period":                 "test",
	})
	resp, err := http.Post(ts.URL+"/v1/credit/rate-of-interest", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", resp.StatusCode)
	}
}

// TestListRatesOfInterest_NeverNull asserts GET /v1/credit/rate-of-interest
// returns {"items":[...]} with a non-null items array even on an empty store.
func TestListRatesOfInterest_NeverNull(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	resp, err := http.Get(ts.URL + "/v1/credit/rate-of-interest")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	var out map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if string(out["items"]) == "null" {
		t.Error("items is null, want [] for empty store")
	}
}

// TestGetRateOfInterest_NotFound asserts 404 for a missing ID.
func TestGetRateOfInterest_NotFound(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	resp, err := http.Get(ts.URL + "/v1/credit/rate-of-interest/doesnotexist")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

// TestComputeInterestRateAnalysis_MarxFixture covers the canonical Ch. 22 fixture:
// avg 2000 bp, interest 500 bp (¼ of 2000) → profit of enterprise 1500 bp.
func TestComputeInterestRateAnalysis_MarxFixture(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	body, _ := json.Marshal(map[string]any{
		"average_profit_rate_bp": 2000,
		"interest_rate_bp":       500,
	})
	resp, err := http.Post(ts.URL+"/v1/credit/interest-rate-analysis", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	var out credit.InterestRateAnalysis
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.ProfitOfEnterpriseBP != 1500 {
		t.Errorf("ProfitOfEnterpriseBP = %d, want 1500", out.ProfitOfEnterpriseBP)
	}
	// Invariant: enterprise = avg - interest
	if out.ProfitOfEnterpriseBP != out.AverageProfitRateBP-out.InterestRateBP {
		t.Errorf("invariant violated: %d != %d - %d",
			out.ProfitOfEnterpriseBP, out.AverageProfitRateBP, out.InterestRateBP)
	}
}

// TestComputeInterestRateAnalysis_BadJSON asserts 400 on malformed JSON.
func TestComputeInterestRateAnalysis_BadJSON(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	resp, err := http.Post(ts.URL+"/v1/credit/interest-rate-analysis", "application/json", bytes.NewReader([]byte(`{bad`)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

// TestComputeInterestRateAnalysis_ValidationError asserts 422 when interest >= avg profit.
func TestComputeInterestRateAnalysis_ValidationError(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	body, _ := json.Marshal(map[string]any{
		"average_profit_rate_bp": 2000,
		"interest_rate_bp":       2000,
	})
	resp, err := http.Post(ts.URL+"/v1/credit/interest-rate-analysis", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", resp.StatusCode)
	}
}
