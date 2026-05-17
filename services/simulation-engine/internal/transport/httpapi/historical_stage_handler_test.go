package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

const englandBody = `{
  "name": "England 15th-18th c.",
  "description": "agrarian expropriation creating free proletariat",
  "primitive_accumulations": [
    {"period": "1450-1640", "method": "enclosure of common lands", "labourers_expropriated": 40000, "capital_formed": 80000},
    {"period": "1500-1640", "method": "expropriation of church property", "labourers_expropriated": 15000, "capital_formed": 30000},
    {"period": "1500-1800", "method": "colonial plunder and slave trade", "labourers_expropriated": 5000, "capital_formed": 90000}
  ]
}`

func TestCreateHistoricalStage_Success(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, nil, nil, nil, nil, st, st, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages", strings.NewReader(englandBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateHistoricalStage(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/v1/historical-stages/") {
		t.Errorf("missing or wrong Location header: %q", loc)
	}
	var resp struct {
		ID                         string `json:"id"`
		TotalCapitalFormed         int64  `json:"total_capital_formed"`
		TotalLabourersExpropriated int64  `json:"total_labourers_expropriated"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Error("response missing id")
	}
	if resp.TotalCapitalFormed != 200000 {
		t.Errorf("TotalCapitalFormed = %d, want 200000", resp.TotalCapitalFormed)
	}
	if resp.TotalLabourersExpropriated != 60000 {
		t.Errorf("TotalLabourersExpropriated = %d, want 60000", resp.TotalLabourersExpropriated)
	}
}

func TestCreateHistoricalStage_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, store.NewMemory(), store.NewMemory(), nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages", strings.NewReader(`{nope`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateHistoricalStage(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateHistoricalStage_BadRequest_MissingName(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, store.NewMemory(), store.NewMemory(), nil, nil, nil, nil, nil, nil, nil, nil, nil)
	body := `{"name":"","primitive_accumulations":[{"period":"x","method":"y","labourers_expropriated":1,"capital_formed":1}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateHistoricalStage(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateHistoricalStage_BadRequest_NoEpisodes(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, store.NewMemory(), store.NewMemory(), nil, nil, nil, nil, nil, nil, nil, nil, nil)
	body := `{"name":"Empty","description":"","primitive_accumulations":[]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateHistoricalStage(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateHistoricalStage_BadRequest_NonPositiveCapital(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, store.NewMemory(), store.NewMemory(), nil, nil, nil, nil, nil, nil, nil, nil, nil)
	body := `{"name":"X","description":"","primitive_accumulations":[{"period":"p","method":"m","labourers_expropriated":1,"capital_formed":0}]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateHistoricalStage(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateHistoricalStage_Conflict_DuplicateName(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, nil, nil, nil, nil, st, st, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages", strings.NewReader(englandBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreateHistoricalStage(w, req)
		if i == 0 && w.Code != http.StatusCreated {
			t.Fatalf("first create want 201, got %d: %s", w.Code, w.Body.String())
		}
		if i == 1 && w.Code != http.StatusConflict {
			t.Fatalf("second create want 409, got %d", w.Code)
		}
	}
}

func TestListHistoricalStages_ReturnsNonNilSlice(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, store.NewMemory(), store.NewMemory(), nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages", nil)
	w := httptest.NewRecorder()
	h.ListHistoricalStages(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if got := strings.TrimSpace(w.Body.String()); got == "null" {
		t.Errorf("response is null, want []: %q", got)
	}
}

func TestListHistoricalStages_ReturnsCreated(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	if _, err := st.CreateHistoricalStage(context.Background(), simulation.HistoricalStage{
		Name:        "Seed",
		Description: "fixture",
		PrimitiveAccumulations: []simulation.PrimitiveAccumulation{
			{Period: "p", Method: "m", LabourersExpropriated: 10, CapitalFormed: 100},
		},
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	h := httpapi.New(nil, nil, nil, nil, nil, st, st, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages", nil)
	w := httptest.NewRecorder()
	h.ListHistoricalStages(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var got []map[string]any
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 stage, got %d", len(got))
	}
	if got[0]["name"] != "Seed" {
		t.Errorf("name = %v, want Seed", got[0]["name"])
	}
}

func TestSeedScenarioFromStage_NotFound(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, store.NewMemory(), store.NewMemory(), nil, nil, nil, nil, nil, nil, nil, nil, nil)
	body := `{"pre_capitalist_workers":100000,"composition_ratio":0.8,"surplus_rate":1.0,"periods":3}`
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/missing/seed-scenario", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "missing")
	w := httptest.NewRecorder()
	h.SeedScenarioFromStage(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSeedScenarioFromStage_Success(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	created, err := st.CreateHistoricalStage(context.Background(), simulation.HistoricalStage{
		Name:        "England 15th-18th c.",
		Description: "agrarian expropriation",
		PrimitiveAccumulations: []simulation.PrimitiveAccumulation{
			{Period: "1450-1640", Method: "enclosure", LabourersExpropriated: 40000, CapitalFormed: 80000},
			{Period: "1500-1640", Method: "church", LabourersExpropriated: 15000, CapitalFormed: 30000},
			{Period: "1500-1800", Method: "plunder", LabourersExpropriated: 5000, CapitalFormed: 90000},
		},
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	h := httpapi.New(nil, nil, nil, nil, nil, st, st, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	body := `{"pre_capitalist_workers":100000,"composition_ratio":0.8,"surplus_rate":1.0,"periods":3}`
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(created.ID)+"/seed-scenario", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", string(created.ID))
	w := httptest.NewRecorder()
	h.SeedScenarioFromStage(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		StageID       string `json:"stage_id"`
		SeededCapital struct {
			ConstantCapital int64 `json:"constant_capital"`
			VariableCapital int64 `json:"variable_capital"`
		} `json:"seeded_capital"`
		Separation struct {
			DisplacedWorkers int64 `json:"displaced_workers"`
			FreeProletarians int64 `json:"free_proletarians"`
		} `json:"separation"`
		Cycles []struct {
			Period int64 `json:"period"`
		} `json:"cycles"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.StageID != string(created.ID) {
		t.Errorf("stage_id = %s, want %s", resp.StageID, created.ID)
	}
	if total := resp.SeededCapital.ConstantCapital + resp.SeededCapital.VariableCapital; total != 200000 {
		t.Errorf("seeded c+v = %d, want 200000", total)
	}
	if resp.Separation.FreeProletarians != resp.Separation.DisplacedWorkers {
		t.Errorf("invariant violated: FP=%d, DW=%d", resp.Separation.FreeProletarians, resp.Separation.DisplacedWorkers)
	}
	if len(resp.Cycles) != 3 {
		t.Errorf("cycles length = %d, want 3", len(resp.Cycles))
	}
}

func TestSeedScenarioFromStage_BadRequest_BadInputs(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, store.NewMemory(), store.NewMemory(), nil, nil, nil, nil, nil, nil, nil, nil, nil)
	cases := []struct {
		name string
		body string
	}{
		{"malformed", `{nope`},
		{"negative_pre_capitalist", `{"pre_capitalist_workers":-1,"composition_ratio":0.5,"surplus_rate":1.0,"periods":1}`},
		{"composition_ratio_ge_one", `{"pre_capitalist_workers":1,"composition_ratio":1.0,"surplus_rate":1.0,"periods":1}`},
		{"negative_surplus_rate", `{"pre_capitalist_workers":1,"composition_ratio":0.5,"surplus_rate":-0.1,"periods":1}`},
		{"zero_periods", `{"pre_capitalist_workers":1,"composition_ratio":0.5,"surplus_rate":1.0,"periods":0}`},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/x/seed-scenario", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("id", "x")
			w := httptest.NewRecorder()
			h.SeedScenarioFromStage(w, req)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("want 400, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}
