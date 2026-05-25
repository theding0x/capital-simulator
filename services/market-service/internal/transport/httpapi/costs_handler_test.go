package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/market-service/internal/costs"
	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
	"github.com/theding0x/capital-simulator/services/market-service/internal/transport/httpapi"
)

func newCostsHandler() *httpapi.Handler {
	return httpapi.New(store.NewMemory(), nil)
}

func postCosts(t *testing.T, handlerFn http.HandlerFunc, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handlerFn(rr, req)
	return rr
}

// POST /v1/circulation-costs — §I(b) book-keeping cost.
func TestCreateCirculationCost_201(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	rr := postCosts(t, h.CreateCirculationCost, map[string]any{
		"industrial_capital_id": "ic-001",
		"kind":                  "book_keeping",
		"pence":                 1200,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", rr.Code, rr.Body)
	}
	var c costs.CirculationCost
	_ = json.Unmarshal(rr.Body.Bytes(), &c)
	if c.ID == "" {
		t.Error("ID must be set")
	}
	if c.Nature != costs.NatureValuePreserving {
		t.Errorf("Nature = %s, want value_preserving", c.Nature)
	}
}

func TestCreateCirculationCost_Transportation_ValueAdding(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	rr := postCosts(t, h.CreateCirculationCost, map[string]any{
		"industrial_capital_id": "ic-002",
		"kind":                  "transportation",
		"pence":                 500,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d; body: %s", rr.Code, rr.Body)
	}
	var c costs.CirculationCost
	_ = json.Unmarshal(rr.Body.Bytes(), &c)
	if c.Nature != costs.NatureValueAdding {
		t.Errorf("Nature = %s, want value_adding", c.Nature)
	}
}

func TestCreateCirculationCost_400_BadJSON(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{bad}"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateCirculationCost(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}

func TestGetCirculationCost_404(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	req := httptest.NewRequest(http.MethodGet, "/v1/circulation-costs/nonexistent", nil)
	rr := httptest.NewRecorder()
	h.GetCirculationCost(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
}

func TestListCirculationCosts_NonNullItems(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	req := httptest.NewRequest(http.MethodGet, "/v1/circulation-costs", nil)
	rr := httptest.NewRecorder()
	h.ListCirculationCosts(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["items"] == nil {
		t.Error("items must not be null")
	}
}

// §I(a): Spinning Mill 1871 seller — 480 necessary + 120 surplus minutes.
func TestCreateCirculationAgent_201(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	rr := postCosts(t, h.CreateCirculationAgent, map[string]any{
		"industrial_capital_id":    "ic-spinner-1871",
		"role":                     "seller",
		"wage_rate_pence":          480,
		"labour_minutes_necessary": 480,
		"labour_minutes_surplus":   120,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d; body: %s", rr.Code, rr.Body)
	}
	var a costs.CirculationAgent
	_ = json.Unmarshal(rr.Body.Bytes(), &a)
	if a.ValueContribution() != 0 {
		t.Error("ValueContribution must be 0 — circulation labour adds no value")
	}
}

func TestCreateCommoditySupply_201(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	rr := postCosts(t, h.CreateCommoditySupply, map[string]any{
		"industrial_capital_id": "ic-001",
		"pence":                 100000,
		"is_voluntary":          true,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d; body: %s", rr.Code, rr.Body)
	}
	if rr.Header().Get("Location") == "" {
		t.Error("Location header must be set")
	}
}

func TestCreateCommoditySupply_400_MissingCapitalID(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	rr := postCosts(t, h.CreateCommoditySupply, map[string]any{"pence": 1000, "is_voluntary": true})
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}

// §III Birmingham glass: 3× breakage risk multiplier.
func TestSetTransportTariff_201_EffectiveRate(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	rr := postCosts(t, h.SetTransportTariff, map[string]any{
		"commodity_id":                          "glass-birmingham",
		"base_pence_per_ton_mile":               24,
		"fragility_multiplier_basis_points":     0,
		"breakage_risk_multiplier_basis_points": 30000,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d; body: %s", rr.Code, rr.Body)
	}
	var tt costs.TransportTariff
	_ = json.Unmarshal(rr.Body.Bytes(), &tt)
	if tt.EffectiveRatePencePerTonMile() != 72 {
		t.Errorf("EffectiveRate = %d, want 72 (3×24)", tt.EffectiveRatePencePerTonMile())
	}
}

func TestCreateTransportLeg_201_AddedValue(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	rr := postCosts(t, h.CreateTransportLeg, map[string]any{
		"commodity_id":             "glass-birmingham",
		"quantity":                 1,
		"origin_id":                "birmingham",
		"destination_id":           "london",
		"distance_meters":          80450,
		"weight_grams":             1000000,
		"labour_cost_pence":        100,
		"means_of_transport_pence": 20,
		"surplus_pence":            0,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d; body: %s", rr.Code, rr.Body)
	}
	var leg costs.TransportLeg
	_ = json.Unmarshal(rr.Body.Bytes(), &leg)
	if leg.AddedValue() != 120 {
		t.Errorf("AddedValue = %d, want 120", leg.AddedValue())
	}
}

// §I(c): 1,000,000 pence reserve, 1% annual replacement.
func TestCreateMoneyAsFauxFrais_201(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	rr := postCosts(t, h.CreateMoneyAsFauxFrais, map[string]any{
		"pence":                    1000000,
		"annual_replacement_pence": 10000,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d; body: %s", rr.Code, rr.Body)
	}
}

func TestSystemFauxFrais_NonNullItems(t *testing.T) {
	t.Parallel()
	h := newCostsHandler()
	req := httptest.NewRequest(http.MethodGet, "/v1/circulation-costs/system-faux-frais", nil)
	rr := httptest.NewRecorder()
	h.SystemFauxFrais(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var res costs.SystemFauxFraisResult
	_ = json.Unmarshal(rr.Body.Bytes(), &res)
	if res.Items == nil {
		t.Error("items must not be null")
	}
}
