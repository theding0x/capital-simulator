package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

// trinityValidBody is the faithful Ch. 48 trinity: capital's apparent revenue is
// the average profit 20, but its actual surplus share is 17 (profit of enterprise
// 12 + interest 5) after the ground-rent 3 is deducted. Profit 17 + rent 3 = one
// surplus-value s = 20; apparent 103 overshoots v+s = 100 by exactly the rent 3.
const trinityValidBody = `{"capital_stream":{"apparent_revenue_bp":20,"actual_source_bp":17,"is_fetishised":true},` +
	`"land_stream":{"apparent_revenue_bp":3,"actual_source_bp":3,"is_fetishised":true},` +
	`"labour_stream":{"apparent_revenue_bp":80,"actual_source_bp":0,"is_fetishised":true}}`

func TestCreateTrinityFormulaCreated(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodPost, "/v1/revenue/trinity-formulas", strings.NewReader(trinityValidBody))
	rec := httptest.NewRecorder()
	h.CreateTrinityFormula(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	if loc := rec.Header().Get("Location"); !strings.HasPrefix(loc, "/v1/revenue/trinity-formulas/") {
		t.Errorf("Location = %q, want /v1/revenue/trinity-formulas/ prefix", loc)
	}
	var resp trinityFormulaResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Formula.TotalApparentRevenueBP != 103 {
		t.Errorf("TotalApparentRevenueBP = %d, want 103", resp.Formula.TotalApparentRevenueBP)
	}
	if resp.Formula.TotalSurplusValueBP != 20 {
		t.Errorf("TotalSurplusValueBP = %d, want 20 (profit 17 + rent 3, one s)", resp.Formula.TotalSurplusValueBP)
	}
	// The fetish overshoot is surfaced and equals the doubly-counted ground-rent.
	if resp.FetishOvershootBP != 3 {
		t.Errorf("FetishOvershootBP = %d, want 3 (apparent 103 - new value 100 = rent)", resp.FetishOvershootBP)
	}
	if len(resp.Streams) != 3 {
		t.Fatalf("streams = %d, want 3", len(resp.Streams))
	}
}

// TestCreateTrinityFormulaDoubleCountRejected pins the #377 regression: a trinity
// whose capital actual double-counts the rent (s would be 23) no longer conserves
// surplus-value and must be rejected with 422.
func TestCreateTrinityFormulaDoubleCountRejected(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	body := `{"capital_stream":{"apparent_revenue_bp":20,"actual_source_bp":20,"is_fetishised":true},` +
		`"land_stream":{"apparent_revenue_bp":3,"actual_source_bp":3,"is_fetishised":true},` +
		`"labour_stream":{"apparent_revenue_bp":80,"actual_source_bp":0,"is_fetishised":true}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/revenue/trinity-formulas", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.CreateTrinityFormula(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreateTrinityFormulaBadJSON(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodPost, "/v1/revenue/trinity-formulas", strings.NewReader("{bad"))
	rec := httptest.NewRecorder()
	h.CreateTrinityFormula(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestCreateTrinityFormulaLabourSurplusRejected(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	body := `{"capital_stream":{"apparent_revenue_bp":20,"actual_source_bp":20,"is_fetishised":true},` +
		`"land_stream":{"apparent_revenue_bp":3,"actual_source_bp":3,"is_fetishised":true},` +
		`"labour_stream":{"apparent_revenue_bp":80,"actual_source_bp":5,"is_fetishised":true}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/revenue/trinity-formulas", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.CreateTrinityFormula(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body=%s", rec.Code, rec.Body.String())
	}
}

func TestGetTrinityFormulaNotFound(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/revenue/trinity-formulas/nope", nil)
	req.SetPathValue("id", "nope")
	rec := httptest.NewRecorder()
	h.GetTrinityFormula(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestGetTrinityFormulaRoundTrip(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	cRec := httptest.NewRecorder()
	h.CreateTrinityFormula(cRec, httptest.NewRequest(http.MethodPost, "/v1/revenue/trinity-formulas", strings.NewReader(trinityValidBody)))
	var created trinityFormulaResponse
	if err := json.Unmarshal(cRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create: %v", err)
	}
	id := string(created.Formula.ID)

	req := httptest.NewRequest(http.MethodGet, "/v1/revenue/trinity-formulas/"+id, nil)
	req.SetPathValue("id", id)
	rec := httptest.NewRecorder()
	h.GetTrinityFormula(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var got revenue.TrinityFormula
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal get: %v", err)
	}
	if string(got.ID) != id {
		t.Errorf("id = %q, want %q", got.ID, id)
	}
}

func TestListRevenueFetishFormsEmptyNonNull(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/revenue/fetish-forms", nil)
	rec := httptest.NewRecorder()
	h.ListRevenueFetishForms(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Body.String(); !strings.Contains(got, `"items":[]`) {
		t.Errorf("body = %s, want items:[]", got)
	}
}

func TestListRevenueFetishFormsNonEmpty(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	if _, err := mem.CreateRevenueFetishForm(context.Background(), revenue.RevenueFetishForm{
		Source:         revenue.RevenueSourceCapital,
		SurfaceFormula: "Capital — Interest",
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	h := New(mem, nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/revenue/fetish-forms", nil)
	rec := httptest.NewRecorder()
	h.ListRevenueFetishForms(rec, req)
	var resp struct {
		Items []revenue.RevenueFetishForm `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(resp.Items))
	}
}
