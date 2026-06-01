package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// getBalanceCheck calls the simple-scheme balance-check endpoint and decodes it.
func getBalanceCheck(t *testing.T, h *Handler, id string) balanceCheckResponse {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet,
		"/v1/reproduction/simple/schemes/"+id+"/balance-check", nil)
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.GetSchemeBalanceCheck(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("balance-check: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp balanceCheckResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode balance-check: %v", err)
	}
	return resp
}

// TestBalanceCheck_AnnualisedSettlement is the issue #221 acceptance test: a
// scheme where Dept I turns over twice as fast as Dept II is balanced per
// period but unbalanced on annualised flows, and the balance-check surfaces
// both the raw and the annualised figures.
func TestBalanceCheck_AnnualisedSettlement(t *testing.T) {
	t.Parallel()
	h := newShareProjectionHandler(t)
	id := createSimpleScheme(t, h, "1865-turnover")
	// Dept I: 4000c+1000v+1000s, turns over twice a year (20000 bp).
	addSimpleDept(t, h, id, `{"department":"department_i","constant_pence":4000,"variable_pence":1000,"surplus_pence":1000,"total_pence":6000,"turnover_number_bp":20000}`)
	// Dept II: 2000c+500v+500s, annual (10000 bp).
	addSimpleDept(t, h, id, `{"department":"department_ii","constant_pence":2000,"variable_pence":500,"surplus_pence":500,"total_pence":3000,"turnover_number_bp":10000}`)

	res := getBalanceCheck(t, h, id)

	// Raw keystone balances (I(v+s)=2000 == II(c)=2000).
	if !res.IsBalanced || res.DeficitPence != 0 {
		t.Fatalf("raw keystone should balance: isBalanced=%v deficit=%d", res.IsBalanced, res.DeficitPence)
	}
	if res.DeptITurnoverBP != 20000 || res.DeptIITurnoverBP != 10000 {
		t.Fatalf("turnover surfaced wrong: I=%d II=%d, want 20000/10000", res.DeptITurnoverBP, res.DeptIITurnoverBP)
	}
	if res.DeptIVplusSAnnualisedPence != 4000 || res.DeptIIConstantAnnualisedPence != 2000 {
		t.Fatalf("annualised flows: I=%d II=%d, want 4000/2000",
			res.DeptIVplusSAnnualisedPence, res.DeptIIConstantAnnualisedPence)
	}
	if res.AnnualisedDeficitPence != 2000 || res.AnnualisedBalanced {
		t.Fatalf("annualised settlement: deficit=%d balanced=%v, want 2000/false",
			res.AnnualisedDeficitPence, res.AnnualisedBalanced)
	}
}

// TestAddDepartment_TurnoverDefaultsToAnnual confirms an omitted turnover
// number round-trips as annual (10000 bp), keeping raw == annualised.
func TestAddDepartment_TurnoverDefaultsToAnnual(t *testing.T) {
	t.Parallel()
	h := newShareProjectionHandler(t)
	id := createSimpleScheme(t, h, "1870")

	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/simple/schemes/"+id+"/departments",
		bytes.NewBufferString(`{"department":"department_i","constant_pence":4000,"variable_pence":1000,"surplus_pence":1000,"total_pence":6000}`))
	req.SetPathValue("id", id)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.AddDepartmentToScheme(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("add dept: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp departmentalCapitalResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.TurnoverNumberBP != 10000 {
		t.Fatalf("omitted turnover should default to 10000, got %d", resp.TurnoverNumberBP)
	}
}
