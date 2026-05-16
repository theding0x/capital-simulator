package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// seedStage creates a historical stage in the given memory store and returns
// its ID — used to satisfy the 404-guard in the wage-statute/vagrancy handlers.
func seedStage(t *testing.T, st *store.Memory) simulation.HistoricalStageID {
	t.Helper()
	stage := simulation.HistoricalStage{
		Name:        "England 15th–18th c.",
		Description: "Paradigm case of primitive accumulation.",
		PrimitiveAccumulations: []simulation.PrimitiveAccumulation{
			{Period: "1450–1640", Method: "enclosure", LabourersExpropriated: 40000, CapitalFormed: 80000},
		},
	}
	created, err := st.CreateHistoricalStage(context.Background(), stage)
	if err != nil {
		t.Fatalf("seedStage: %v", err)
	}
	return created.ID
}

func TestCreateWageStatute_Created(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStage(t, st)
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, st, st)

	body := map[string]any{
		"period":              "1349",
		"jurisdiction":        "England",
		"max_wage_pence":      24,
		"min_wage_pence":      0,
		"enforcement_penalty": "imprisonment",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(stageID)+"/wage-statutes", bytes.NewReader(b))
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.CreateWageStatute(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	if w.Header().Get("Location") == "" {
		t.Error("Location header missing")
	}
}

func TestCreateWageStatute_StageNotFound(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, st, st)

	body := map[string]any{
		"period": "1349", "jurisdiction": "England",
		"max_wage_pence": 24, "min_wage_pence": 0,
		"enforcement_penalty": "imprisonment",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/nonexistent/wage-statutes", bytes.NewReader(b))
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()
	h.CreateWageStatute(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestCreateWageStatute_BadJSON(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStage(t, st)
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, st, st)

	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(stageID)+"/wage-statutes", bytes.NewReader([]byte("notjson")))
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.CreateWageStatute(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateWageStatute_InvalidBounds(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStage(t, st)
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, st, st)

	body := map[string]any{
		"period": "1600", "jurisdiction": "England",
		"max_wage_pence": 5, "min_wage_pence": 10,
		"enforcement_penalty": "fine",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(stageID)+"/wage-statutes", bytes.NewReader(b))
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.CreateWageStatute(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateVagrancyLaw_Created(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStage(t, st)
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, st, st)

	body := map[string]any{
		"period":            "1530",
		"jurisdiction":      "England",
		"punishment":        "whipping+imprisonment",
		"target_population": "sturdy vagabonds",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(stageID)+"/vagrancy-laws", bytes.NewReader(b))
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.CreateVagrancyLaw(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateVagrancyLaw_StageNotFound(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, st, st)

	body := map[string]any{
		"period": "1530", "jurisdiction": "England",
		"punishment": "whipping", "target_population": "vagabonds",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/missing/vagrancy-laws", bytes.NewReader(b))
	req.SetPathValue("id", "missing")
	w := httptest.NewRecorder()
	h.CreateVagrancyLaw(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestGetLabourDiscipline_Empty(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStage(t, st)
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, st, st)

	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/"+string(stageID)+"/labour-discipline", nil)
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.GetLabourDiscipline(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["wage_statutes"] == nil {
		t.Error("wage_statutes must be non-null")
	}
	if resp["vagrancy_laws"] == nil {
		t.Error("vagrancy_laws must be non-null")
	}
}

func TestGetLabourDiscipline_NotFound(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, st, st)

	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/ghost/labour-discipline", nil)
	req.SetPathValue("id", "ghost")
	w := httptest.NewRecorder()
	h.GetLabourDiscipline(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestGetLabourDiscipline_RoundTrip(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStage(t, st)
	h := httpapi.New(nil, nil, nil, nil, nil, st, nil, st, st)

	wsBody := map[string]any{
		"period": "1349", "jurisdiction": "England",
		"max_wage_pence": 24, "min_wage_pence": 0,
		"enforcement_penalty": "imprisonment",
	}
	b, _ := json.Marshal(wsBody)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.SetPathValue("id", string(stageID))
	rec := httptest.NewRecorder()
	h.CreateWageStatute(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create wage statute: %d", rec.Code)
	}

	vlBody := map[string]any{
		"period": "1530", "jurisdiction": "England",
		"punishment": "whipping+imprisonment", "target_population": "sturdy vagabonds",
	}
	b, _ = json.Marshal(vlBody)
	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.SetPathValue("id", string(stageID))
	rec = httptest.NewRecorder()
	h.CreateVagrancyLaw(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create vagrancy law: %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.GetLabourDiscipline(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	ws := resp["wage_statutes"].([]any)
	vl := resp["vagrancy_laws"].([]any)
	if len(ws) != 1 {
		t.Errorf("want 1 wage statute, got %d", len(ws))
	}
	if len(vl) != 1 {
		t.Errorf("want 1 vagrancy law, got %d", len(vl))
	}
	mechs := resp["mechanisms"].([]any)
	if len(mechs) != 2 {
		t.Errorf("want 2 mechanisms, got %d", len(mechs))
	}
}

func TestCompareStatutoryWage_Suppression(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	body := map[string]any{"acted_wage_pence": 31, "market_wage_pence": 35}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/statutory-wages/compare", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.CompareStatutoryWage(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp simulation.StatutoryWage
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Deviation >= 0 {
		t.Errorf("want negative deviation (wage suppression), got %f", resp.Deviation)
	}
}

func TestCompareStatutoryWage_ZeroMarket(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	body := map[string]any{"acted_wage_pence": 10, "market_wage_pence": 0}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/statutory-wages/compare", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.CompareStatutoryWage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCompareStatutoryWage_BadJSON(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/statutory-wages/compare", bytes.NewReader([]byte("bad")))
	w := httptest.NewRecorder()
	h.CompareStatutoryWage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}