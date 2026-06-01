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

func TestListCompetitionIllusionsEmptyNonNull(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.ListCompetitionIllusions(rec, httptest.NewRequest(http.MethodGet, "/v1/revenue/competition-illusions", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"items":[]`) {
		t.Errorf("body = %s, want items:[]", rec.Body.String())
	}
}

func TestListCompetitionIllusionsNonEmpty(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	if _, err := mem.CreateCompetitionIllusion(context.Background(), revenue.CompetitionIllusion{Name: "Profit arises in exchange", ProducedBy: "price_competition"}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	h := New(mem, nil)
	rec := httptest.NewRecorder()
	h.ListCompetitionIllusions(rec, httptest.NewRequest(http.MethodGet, "/v1/revenue/competition-illusions", nil))
	var resp struct {
		Items []revenue.CompetitionIllusion `json:"items"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(resp.Items))
	}
}

func TestCreateCostPriceFetishCreated(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	body := `{"constant_capital_bp":70,"variable_capital_bp":30,"profit_bp":20}`
	rec := httptest.NewRecorder()
	h.CreateCostPriceFetish(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/cost-price-fetishes", strings.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var got revenue.CostPriceFetish
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.CostPriceBP != 100 {
		t.Errorf("cost-price = %d, want 100", got.CostPriceBP)
	}
	if got.ActualValueBP != 120 {
		t.Errorf("actual value = %d, want 120", got.ActualValueBP)
	}
}

func TestCreateCostPriceFetishBadJSON(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.CreateCostPriceFetish(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/cost-price-fetishes", strings.NewReader("{bad")))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestCreateCostPriceFetishNegative(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	body := `{"constant_capital_bp":-1,"variable_capital_bp":30,"profit_bp":20}`
	h.CreateCostPriceFetish(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/cost-price-fetishes", strings.NewReader(body)))
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", rec.Code)
	}
}

func TestCreateAnnualRevenueAccountCreated(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	body := `{"constant_capital_lm":6000,"wages_share_lm":1500,"profit_and_rent_share_lm":1500}`
	rec := httptest.NewRecorder()
	h.CreateAnnualRevenueAccount(rec, httptest.NewRequest(http.MethodPost, "/v1/revenue/annual-accounts", strings.NewReader(body)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	var got revenue.AnnualRevenueAccount
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.NewValueCreatedLM != 3000 {
		t.Errorf("new value = %d, want 3000", got.NewValueCreatedLM)
	}
	if got.TotalSocialProductLM != 9000 {
		t.Errorf("total = %d, want 9000", got.TotalSocialProductLM)
	}
}

func TestListValueAppearanceContrastsEmptyNonNull(t *testing.T) {
	t.Parallel()
	h := New(store.NewMemory(), nil)
	rec := httptest.NewRecorder()
	h.ListValueAppearanceContrasts(rec, httptest.NewRequest(http.MethodGet, "/v1/revenue/value-contrasts", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"items":[]`) {
		t.Errorf("body = %s, want items:[]", rec.Body.String())
	}
}
