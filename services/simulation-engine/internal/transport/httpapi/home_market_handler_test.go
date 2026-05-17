package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// § Westphalian flax: peasants who all spun flax are forcibly expropriated;
// large establishments arise; their output enters the commodity market.
const westphalianFlaxBody = `{
  "name": "Westphalian flax spinning",
  "households_engaged": 2000,
  "annual_output_pence": 48000
}`

// seedStageForMarket creates a historical stage in the memory store and returns its ID.
func seedStageForMarket(t *testing.T, st *store.Memory) simulation.HistoricalStageID {
	t.Helper()
	stage, err := st.CreateHistoricalStage(t.Context(), simulation.HistoricalStage{
		Name:        "England 15th–18th c.",
		Description: "Primitive accumulation",
		PrimitiveAccumulations: []simulation.PrimitiveAccumulation{
			{Period: "1450–1640", Method: "enclosure", LabourersExpropriated: 40000, CapitalFormed: 80000},
		},
	})
	if err != nil {
		t.Fatalf("seedStageForMarket: %v", err)
	}
	return stage.ID
}

func TestCreateDomesticIndustry_Success(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForMarket(t, st)
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, nil, nil, nil, st, nil, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/"+string(stageID)+"/domestic-industries",
		strings.NewReader(westphalianFlaxBody))
	req.SetPathValue("id", string(stageID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateDomesticIndustry(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); !strings.Contains(loc, "domestic-industries") {
		t.Errorf("missing or wrong Location header: %q", loc)
	}
	var resp struct {
		ID                string `json:"id"`
		HouseholdsEngaged int64  `json:"households_engaged"`
		AnnualOutputPence int64  `json:"annual_output_pence"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Error("response missing id")
	}
	if resp.HouseholdsEngaged != 2000 {
		t.Errorf("households_engaged = %d, want 2000", resp.HouseholdsEngaged)
	}
	if resp.AnnualOutputPence != 48000 {
		t.Errorf("annual_output_pence = %d, want 48000", resp.AnnualOutputPence)
	}
}

func TestCreateDomesticIndustry_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, nil, nil, nil, st, nil, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/abc/domestic-industries",
		strings.NewReader(`{nope`))
	req.SetPathValue("id", "abc")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateDomesticIndustry(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateDomesticIndustry_BadRequest_ZeroHouseholds(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, nil, nil, nil, st, nil, nil, nil, nil, nil)

	body := `{"name":"test","households_engaged":0,"annual_output_pence":1000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/abc/domestic-industries",
		strings.NewReader(body))
	req.SetPathValue("id", "abc")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateDomesticIndustry(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestComputeMarketFormation_Success(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	// § Westphalian flax: 1500 of 2000 households expropriated.
	body := `{"name":"Westphalian flax spinning","households_engaged":2000,"annual_output_pence":48000,"expropriated":1500}`
	req := httptest.NewRequest(http.MethodPost, "/v1/market-formation", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ComputeMarketFormation(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		LabourersProletarianised  int64  `json:"labourers_proletarianised"`
		NewWageLabourers          int64  `json:"new_wage_labourers"`
		DomesticIndustryDestroyed bool   `json:"domestic_industry_destroyed"`
		MarketSizePence           int64  `json:"market_size_pence"`
		Period                    string `json:"period"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.NewWageLabourers != resp.LabourersProletarianised {
		t.Errorf("invariant: new_wage_labourers(%d) != labourers_proletarianised(%d)",
			resp.NewWageLabourers, resp.LabourersProletarianised)
	}
	if resp.NewWageLabourers != 1500 {
		t.Errorf("want 1500 new wage labourers, got %d", resp.NewWageLabourers)
	}
	if resp.MarketSizePence != 48000 {
		t.Errorf("want market_size_pence 48000, got %d", resp.MarketSizePence)
	}
	if resp.DomesticIndustryDestroyed {
		t.Error("partial expropriation should not destroy the industry")
	}
}

func TestComputeMarketFormation_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/market-formation", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ComputeMarketFormation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestGetHomeMarket_Success(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForMarket(t, st)

	_, _ = st.CreateDomesticIndustry(t.Context(), simulation.DomesticIndustry{
		HistoricalStageID: stageID,
		Name:              "Westphalian flax spinning",
		HouseholdsEngaged: 2000,
		AnnualOutputPence: 48000,
	})
	_, _ = st.CreateDomesticIndustry(t.Context(), simulation.DomesticIndustry{
		HistoricalStageID: stageID,
		Name:              "Westphalian linen weaving",
		HouseholdsEngaged: 1000,
		AnnualOutputPence: 24000,
	})

	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, nil, nil, nil, st, nil, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/"+string(stageID)+"/home-market", nil)
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.GetHomeMarket(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		CommoditisedOutputPence int64 `json:"commoditised_output_pence"`
		WageLabourers           int64 `json:"wage_labourers"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.CommoditisedOutputPence != 72000 {
		t.Errorf("want commoditised_output_pence 72000, got %d", resp.CommoditisedOutputPence)
	}
	if resp.WageLabourers != 3000 {
		t.Errorf("want 3000 wage_labourers, got %d", resp.WageLabourers)
	}
}

func TestGetHomeMarket_NotFound(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, nil, nil, nil, st, nil, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/nonexistent/home-market", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()
	h.GetHomeMarket(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetHomeMarket_EmptyReturnsZero(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForMarket(t, st)

	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, nil, nil, nil, st, nil, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/"+string(stageID)+"/home-market", nil)
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.GetHomeMarket(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		CommoditisedOutputPence int64 `json:"commoditised_output_pence"`
		WageLabourers           int64 `json:"wage_labourers"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.CommoditisedOutputPence != 0 || resp.WageLabourers != 0 {
		t.Errorf("empty stage should yield zero home market size, got %+v", resp)
	}
}
