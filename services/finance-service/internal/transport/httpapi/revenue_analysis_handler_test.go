package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

func createCh49Composition(t *testing.T, h *Handler) annualCompositionResponse {
	t.Helper()
	body := `{"constant_capital_labour_minutes":6000,"variable_capital_labour_minutes":1500,"surplus_value_labour_minutes":1500}`
	rec := httptest.NewRecorder()
	h.CreateAnnualValueComposition(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/annual-compositions", strings.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("create composition status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var resp annualCompositionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return resp
}

func TestCreateAnnualValueCompositionCreated(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	resp := createCh49Composition(t, h)
	if resp.Composition.TotalValueLabourMinutes != 9000 {
		t.Errorf("total value = %d, want 9000", resp.Composition.TotalValueLabourMinutes)
	}
	if resp.Composition.NewValueLabourMinutes != 3000 {
		t.Errorf("new value = %d, want 3000", resp.Composition.NewValueLabourMinutes)
	}
	if resp.RevenueLimit.MaxDistributableRevenueLM != 3000 {
		t.Errorf("max distributable = %d, want 3000", resp.RevenueLimit.MaxDistributableRevenueLM)
	}
	if resp.RevenueLimit.ConstantCapitalReplacement != 6000 {
		t.Errorf("constant replacement = %d, want 6000", resp.RevenueLimit.ConstantCapitalReplacement)
	}
}

func TestCreateAnnualValueCompositionBadJSON(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.CreateAnnualValueComposition(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/annual-compositions", strings.NewReader("{bad")))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestCreateAnnualValueCompositionNegative(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	body := `{"constant_capital_labour_minutes":-1,"variable_capital_labour_minutes":1500,"surplus_value_labour_minutes":1500}`
	h.CreateAnnualValueComposition(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/annual-compositions", strings.NewReader(body)))
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rec.Code)
	}
}

func TestGetAnnualValueCompositionRoundTrip(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	created := createCh49Composition(t, h)
	id := string(created.Composition.ID)

	req := httptest.NewRequest(http.MethodGet, "/v1/revenue/annual-compositions/"+id, nil)
	req.SetPathValue("id", id)
	rec := httptest.NewRecorder()
	h.GetAnnualValueComposition(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var got annualCompositionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if string(got.Composition.ID) != id {
		t.Errorf("id = %q, want %q", got.Composition.ID, id)
	}
	if got.RevenueLimit.MaxDistributableRevenueLM != 3000 {
		t.Errorf("limit max = %d, want 3000", got.RevenueLimit.MaxDistributableRevenueLM)
	}
}

func TestGetAnnualValueCompositionNotFound(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/revenue/annual-compositions/nope", nil)
	req.SetPathValue("id", "nope")
	rec := httptest.NewRecorder()
	h.GetAnnualValueComposition(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestCreateSurplusValuePartitionCreated(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	comp := createCh49Composition(t, h)
	body := `{"composition_id":"` + string(comp.Composition.ID) + `","profit_enterprise_lm":900,"interest_lm":300,"ground_rent_lm":300,"residual_surplus_lm":0}`
	rec := httptest.NewRecorder()
	h.CreateSurplusValuePartition(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/surplus-partitions", strings.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var got revenue.SurplusValuePartition
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.TotalSurplusValueLM != 1500 {
		t.Errorf("total = %d, want 1500", got.TotalSurplusValueLM)
	}
}

func TestCreateSurplusValuePartitionExceedsSurplus(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	comp := createCh49Composition(t, h)
	// 1000 + 300 + 300 = 1600 > produced surplus-value 1500.
	body := `{"composition_id":"` + string(comp.Composition.ID) + `","profit_enterprise_lm":1000,"interest_lm":300,"ground_rent_lm":300,"residual_surplus_lm":0}`
	rec := httptest.NewRecorder()
	h.CreateSurplusValuePartition(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/surplus-partitions", strings.NewReader(body)))
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreateSurplusValuePartitionCompositionNotFound(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	body := `{"composition_id":"nope","profit_enterprise_lm":900,"interest_lm":300,"ground_rent_lm":300,"residual_surplus_lm":0}`
	rec := httptest.NewRecorder()
	h.CreateSurplusValuePartition(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/surplus-partitions", strings.NewReader(body)))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}
