package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

func newFarmerHandler(t *testing.T) (*httpapi.Handler, *store.Memory) {
	t.Helper()
	mem := store.NewMemory()
	h := &httpapi.Handler{
		HistoricalStages: mem,
		FarmTenures:      mem,
	}
	return h, mem
}

func seedStageForFarmer(t *testing.T, mem *store.Memory) simulation.HistoricalStage {
	t.Helper()
	stage, err := mem.CreateHistoricalStage(t.Context(), simulation.HistoricalStage{
		Name:        "England 15th-18th c.",
		Description: "Ch.29 test stage",
		PrimitiveAccumulations: []simulation.PrimitiveAccumulation{
			{Period: "1450-1640", Method: "enclosure", LabourersExpropriated: 100, CapitalFormed: 5000},
		},
	})
	if err != nil {
		t.Fatalf("seed stage: %v", err)
	}
	return stage
}

func TestCreateFarmTenure_Created(t *testing.T) {
	t.Parallel()
	h, mem := newFarmerHandler(t)
	stage := seedStageForFarmer(t, mem)

	body, _ := json.Marshal(map[string]any{
		"form":                   "metayer",
		"lease_period_years":     1,
		"rent_pence":             0,
		"capital_advanced_pence": 500,
		"revenue_pence":          1000,
		"wage_costs_pence":       200,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(stage.ID)+"/farm-tenures", bytes.NewReader(body))
	req.SetPathValue("id", string(stage.ID))
	w := httptest.NewRecorder()
	h.CreateFarmTenure(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	if w.Header().Get("Location") == "" {
		t.Fatal("expected Location header")
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] == "" {
		t.Fatal("expected id in response")
	}
}

// L-H2 (#412): once a stage has recorded the capitalist-farmer form, a later
// regression to bailiff is rejected at the API with 409 Conflict; recording the
// same (highest) form again is still allowed.
func TestCreateFarmTenure_FormRegression(t *testing.T) {
	t.Parallel()
	h, mem := newFarmerHandler(t)
	stage := seedStageForFarmer(t, mem)

	post := func(form string) int {
		body, _ := json.Marshal(map[string]any{
			"form": form, "lease_period_years": 1, "capital_advanced_pence": 100,
			"revenue_pence": 200, "wage_costs_pence": 50,
		})
		req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(stage.ID)+"/farm-tenures", bytes.NewReader(body))
		req.SetPathValue("id", string(stage.ID))
		w := httptest.NewRecorder()
		h.CreateFarmTenure(w, req)
		return w.Code
	}

	if code := post("capitalist-farmer"); code != http.StatusCreated {
		t.Fatalf("seed capitalist-farmer: want 201, got %d", code)
	}
	if code := post("bailiff"); code != http.StatusConflict {
		t.Fatalf("regress to bailiff: want 409, got %d", code)
	}
	if code := post("capitalist-farmer"); code != http.StatusCreated {
		t.Fatalf("repeat capitalist-farmer: want 201, got %d", code)
	}
}

func TestCreateFarmTenure_StageNotFound(t *testing.T) {
	t.Parallel()
	h, _ := newFarmerHandler(t)

	body, _ := json.Marshal(map[string]any{
		"form": "metayer", "lease_period_years": 1,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/nonexistent/farm-tenures", bytes.NewReader(body))
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()
	h.CreateFarmTenure(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestCreateFarmTenure_BadJSON(t *testing.T) {
	t.Parallel()
	h, mem := newFarmerHandler(t)
	stage := seedStageForFarmer(t, mem)

	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(stage.ID)+"/farm-tenures", bytes.NewReader([]byte("{")))
	req.SetPathValue("id", string(stage.ID))
	w := httptest.NewRecorder()
	h.CreateFarmTenure(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateFarmTenure_InvalidForm(t *testing.T) {
	t.Parallel()
	h, mem := newFarmerHandler(t)
	stage := seedStageForFarmer(t, mem)

	body, _ := json.Marshal(map[string]any{
		"form":               "lord",
		"lease_period_years": 10,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(stage.ID)+"/farm-tenures", bytes.NewReader(body))
	req.SetPathValue("id", string(stage.ID))
	w := httptest.NewRecorder()
	h.CreateFarmTenure(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestListFarmTenures_Empty(t *testing.T) {
	t.Parallel()
	h, mem := newFarmerHandler(t)
	stage := seedStageForFarmer(t, mem)

	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/"+string(stage.ID)+"/farm-tenures", nil)
	req.SetPathValue("id", string(stage.ID))
	w := httptest.NewRecorder()
	h.ListFarmTenures(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var out []any
	if err := json.NewDecoder(w.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out == nil {
		t.Fatal("list must not be null")
	}
}

func TestListFarmTenures_StageNotFound(t *testing.T) {
	t.Parallel()
	h, _ := newFarmerHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/missing/farm-tenures", nil)
	req.SetPathValue("id", "missing")
	w := httptest.NewRecorder()
	h.ListFarmTenures(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestComputeRealRent_HalfFactor(t *testing.T) {
	t.Parallel()
	h, _ := newFarmerHandler(t)

	body, _ := json.Marshal(map[string]any{
		"nominal_rent_pence":  1200,
		"depreciation_factor": 0.5,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/farm-tenures/real-rent", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ComputeRealRent(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["real_rent_pence"].(float64) != 600 {
		t.Fatalf("want real_rent_pence 600, got %v", resp["real_rent_pence"])
	}
}

func TestComputeRealRent_InvalidFactor(t *testing.T) {
	t.Parallel()
	h, _ := newFarmerHandler(t)

	body, _ := json.Marshal(map[string]any{
		"nominal_rent_pence":  1000,
		"depreciation_factor": 0.0,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/farm-tenures/real-rent", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ComputeRealRent(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestComputeRealRent_BadJSON(t *testing.T) {
	t.Parallel()
	h, _ := newFarmerHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/v1/farm-tenures/real-rent", bytes.NewReader([]byte("{")))
	w := httptest.NewRecorder()
	h.ComputeRealRent(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}
