package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestCreateProfitDivision_BorrowedCapital asserts 201 + Location on valid input.
// Fixture: borrowed capital — total 2000 bp, interest 500 bp → enterprise 1500 bp.
func TestCreateProfitDivision_BorrowedCapital(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	body, _ := json.Marshal(map[string]any{
		"total_profit_bp":     2000,
		"interest_bp":         500,
		"ownership_form":      int(credit.FormSeparated),
		"wages_labour_minutes": 0,
		"mystification_index": 5000,
	})
	resp, err := http.Post(ts.URL+"/v1/credit/profit-division", "application/json", bytes.NewReader(body))
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

	var pd credit.ProfitDivision
	if err := json.NewDecoder(resp.Body).Decode(&pd); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if pd.ProfitOfEnterpriseBP != 1500 {
		t.Errorf("ProfitOfEnterpriseBP = %d, want 1500", pd.ProfitOfEnterpriseBP)
	}
	// Partition invariant.
	if pd.TotalProfitBP != pd.InterestBP+pd.ProfitOfEnterpriseBP {
		t.Errorf("partition violated: %d != %d + %d",
			pd.TotalProfitBP, pd.InterestBP, pd.ProfitOfEnterpriseBP)
	}
}

// TestCreateProfitDivision_OwnCapital asserts 201 for own-capital case (interest 0).
func TestCreateProfitDivision_OwnCapital(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	body, _ := json.Marshal(map[string]any{
		"total_profit_bp":     2000,
		"interest_bp":         0,
		"ownership_form":      int(credit.FormOwnerOperator),
		"wages_labour_minutes": 480,
		"mystification_index": 2000,
	})
	resp, err := http.Post(ts.URL+"/v1/credit/profit-division", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d, want 201", resp.StatusCode)
	}

	var pd credit.ProfitDivision
	if err := json.NewDecoder(resp.Body).Decode(&pd); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if pd.ProfitOfEnterpriseBP != 2000 {
		t.Errorf("ProfitOfEnterpriseBP = %d, want 2000", pd.ProfitOfEnterpriseBP)
	}
}

// TestCreateProfitDivision_BadJSON asserts 400 on malformed JSON body.
func TestCreateProfitDivision_BadJSON(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	resp, err := http.Post(ts.URL+"/v1/credit/profit-division", "application/json", bytes.NewReader([]byte(`{bad`)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

// TestCreateProfitDivision_ValidationError asserts 422 when interest_bp >= total_profit_bp.
func TestCreateProfitDivision_ValidationError(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	body, _ := json.Marshal(map[string]any{
		"total_profit_bp":     2000,
		"interest_bp":         2000,
		"ownership_form":      int(credit.FormSeparated),
		"wages_labour_minutes": 0,
		"mystification_index": 0,
	})
	resp, err := http.Post(ts.URL+"/v1/credit/profit-division", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", resp.StatusCode)
	}
}

// TestCreateProfitDivision_GeneralRateProvenance asserts total_profit_bp is sourced
// from the stored general rate when general_rate_id is supplied, ignoring any
// caller-supplied total_profit_bp.
func TestCreateProfitDivision_GeneralRateProvenance(t *testing.T) {
	t.Parallel()

	ts, st := newFinanceTestServer(t)
	gr, err := st.CreateGeneralProfitRate(t.Context(), avgprofit.GeneralProfitRate{Rate: 2000})
	if err != nil {
		t.Fatalf("seed general rate: %v", err)
	}

	body, _ := json.Marshal(map[string]any{
		"total_profit_bp":      9999, // wrong on purpose — must be ignored
		"general_rate_id":      string(gr.ID),
		"interest_bp":          500,
		"ownership_form":       int(credit.FormSeparated),
		"wages_labour_minutes": 0,
		"mystification_index":  5000,
	})
	resp, err := http.Post(ts.URL+"/v1/credit/profit-division", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", resp.StatusCode)
	}
	var pd credit.ProfitDivision
	if err := json.NewDecoder(resp.Body).Decode(&pd); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if pd.TotalProfitBP != 2000 {
		t.Errorf("TotalProfitBP = %d, want 2000 (from stored general rate, not the supplied 9999)", pd.TotalProfitBP)
	}
	if pd.ProfitOfEnterpriseBP != 1500 {
		t.Errorf("ProfitOfEnterpriseBP = %d, want 1500 (2000 - 500)", pd.ProfitOfEnterpriseBP)
	}
}

// TestCreateProfitDivision_UnknownGeneralRate asserts an unknown general_rate_id yields 404.
func TestCreateProfitDivision_UnknownGeneralRate(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	body, _ := json.Marshal(map[string]any{
		"general_rate_id":      "doesnotexist",
		"interest_bp":          500,
		"ownership_form":       int(credit.FormSeparated),
		"wages_labour_minutes": 0,
		"mystification_index":  5000,
	})
	resp, err := http.Post(ts.URL+"/v1/credit/profit-division", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

// TestListProfitDivisions_NeverNull asserts GET /v1/credit/profit-division
// returns {"items":[...]} with a non-null items array even on an empty store.
func TestListProfitDivisions_NeverNull(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	resp, err := http.Get(ts.URL + "/v1/credit/profit-division")
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

// TestGetProfitDivision_NotFound asserts 404 for a missing ID.
func TestGetProfitDivision_NotFound(t *testing.T) {
	t.Parallel()

	ts, _ := newFinanceTestServer(t)

	resp, err := http.Get(ts.URL + "/v1/credit/profit-division/doesnotexist")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}
