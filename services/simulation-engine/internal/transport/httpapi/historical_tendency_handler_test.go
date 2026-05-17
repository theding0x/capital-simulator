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

// Ch. 32 — Historical Tendency of Capitalist Accumulation HTTP layer.

func newCh32Handler() *httpapi.Handler {
	mem := store.NewMemory()
	return httpapi.New(nil, httpapi.Deps{Trajectories: mem})
}

func TestGetNegationOfNegation_ReturnsThreeStages(t *testing.T) {
	t.Parallel()
	h := newCh32Handler()
	req := httptest.NewRequest(http.MethodGet, "/v1/accumulation/negation-of-negation", nil)
	w := httptest.NewRecorder()
	h.GetNegationOfNegation(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Stages []struct {
			Stage       string `json:"stage"`
			Description string `json:"description"`
		} `json:"stages"`
		Negation struct {
			Stage string `json:"stage"`
		} `json:"negation"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Stages) != 3 {
		t.Fatalf("want 3 stages, got %d", len(resp.Stages))
	}
	want := []string{
		simulation.NegationStagePettyProperty,
		simulation.NegationStageCapitalistExpropriation,
		simulation.NegationStageSocialisedProperty,
	}
	for i, w := range want {
		if resp.Stages[i].Stage != w {
			t.Errorf("stage[%d] = %q, want %q", i, resp.Stages[i].Stage, w)
		}
		if resp.Stages[i].Description == "" {
			t.Errorf("stage[%d] description must not be empty", i)
		}
	}
	if resp.Negation.Stage != simulation.NegationStageSocialisedProperty {
		t.Errorf("negation.stage = %q, want %q", resp.Negation.Stage, simulation.NegationStageSocialisedProperty)
	}
}

func TestRunCentralisation_HappyPath(t *testing.T) {
	t.Parallel()
	body := `{"firms":1000,"total_capital_pence":480000000,"wage_labourers":50000,"steps":6,"absorption_rate":0.2,"capital_growth_rate":0.05}`
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/centralisation", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newCh32Handler().RunCentralisation(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		InitialFirms      int64 `json:"initial_firms"`
		FinalFirms        int64 `json:"final_firms"`
		FinalCapitalPence int64 `json:"final_capital_pence"`
		ReserveArmySize   int64 `json:"reserve_army_size"`
		Steps             []struct {
			StepIndex                int64 `json:"step_index"`
			FirmsAbsorbed            int64 `json:"firms_absorbed"`
			CapitalConcentratedPence int64 `json:"capital_concentrated_pence"`
		} `json:"steps"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.InitialFirms != 1000 {
		t.Errorf("initial_firms = %d, want 1000", resp.InitialFirms)
	}
	if resp.FinalFirms >= resp.InitialFirms {
		t.Errorf("centralisation must reduce firms: initial %d, final %d", resp.InitialFirms, resp.FinalFirms)
	}
	if len(resp.Steps) == 0 {
		t.Fatal("trajectory must contain at least one step")
	}
	if resp.ReserveArmySize <= 0 {
		t.Errorf("reserve_army_size should grow, got %d", resp.ReserveArmySize)
	}
}

func TestRunCentralisation_BadRequest_FirmsNonPositive(t *testing.T) {
	t.Parallel()
	body := `{"firms":0,"total_capital_pence":1000,"wage_labourers":10,"steps":3,"absorption_rate":0.2,"capital_growth_rate":0.0}`
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/centralisation", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newCh32Handler().RunCentralisation(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRunCentralisation_BadRequest_AbsorptionRateOutOfRange(t *testing.T) {
	t.Parallel()
	for _, body := range []string{
		`{"firms":10,"total_capital_pence":1000,"wage_labourers":10,"steps":3,"absorption_rate":0,"capital_growth_rate":0.0}`,
		`{"firms":10,"total_capital_pence":1000,"wage_labourers":10,"steps":3,"absorption_rate":1.5,"capital_growth_rate":0.0}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/centralisation", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		newCh32Handler().RunCentralisation(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("body=%s: want 400, got %d", body, w.Code)
		}
	}
}

func TestRunCentralisation_MalformedJSON(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/centralisation", bytes.NewBufferString(`{nope`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newCh32Handler().RunCentralisation(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestRunCentralisation_PersistRequiresName(t *testing.T) {
	t.Parallel()
	body := `{"firms":100,"total_capital_pence":100000,"wage_labourers":500,"steps":3,"absorption_rate":0.2,"capital_growth_rate":0.0,"persist":true}`
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/centralisation", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newCh32Handler().RunCentralisation(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRunCentralisation_PersistAndRetrieve(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	h := httpapi.New(nil, httpapi.Deps{Trajectories: mem})

	body := `{"name":"Lancashire cotton, 1820-1880","firms":1000,"total_capital_pence":480000000,"wage_labourers":50000,"steps":6,"absorption_rate":0.2,"capital_growth_rate":0.05,"persist":true}`
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/centralisation", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.RunCentralisation(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	loc := w.Header().Get("Location")
	if loc == "" {
		t.Fatal("expected Location header on 201")
	}
	var created struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.ID == "" {
		t.Fatal("created.id must not be empty")
	}
	if created.Name != "Lancashire cotton, 1820-1880" {
		t.Errorf("created.name = %q, want Lancashire cotton, 1820-1880", created.Name)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/v1/accumulation/trajectories/"+created.ID, nil)
	getReq.SetPathValue("id", created.ID)
	getW := httptest.NewRecorder()
	h.GetAccumulationTrajectory(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("get want 200, got %d: %s", getW.Code, getW.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/v1/accumulation/trajectories", nil)
	listW := httptest.NewRecorder()
	h.ListAccumulationTrajectories(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("list want 200, got %d", listW.Code)
	}
	var list []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(listW.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list) != 1 || list[0].ID != created.ID {
		t.Fatalf("want list = [created.id], got %+v", list)
	}
}

func TestRunCentralisation_PersistDuplicateNameConflicts(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	h := httpapi.New(nil, httpapi.Deps{Trajectories: mem})
	body := `{"name":"Duplicate","firms":50,"total_capital_pence":1000,"wage_labourers":10,"steps":2,"absorption_rate":0.2,"capital_growth_rate":0.0,"persist":true}`
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/centralisation", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.RunCentralisation(w, req)
		if i == 0 && w.Code != http.StatusCreated {
			t.Fatalf("first want 201, got %d", w.Code)
		}
		if i == 1 && w.Code != http.StatusConflict {
			t.Fatalf("second want 409, got %d: %s", w.Code, w.Body.String())
		}
	}
}

func TestGetAccumulationTrajectory_NotFound(t *testing.T) {
	t.Parallel()
	h := newCh32Handler()
	req := httptest.NewRequest(http.MethodGet, "/v1/accumulation/trajectories/missingid", nil)
	req.SetPathValue("id", "missingid")
	w := httptest.NewRecorder()
	h.GetAccumulationTrajectory(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestListAccumulationTrajectories_EmptyReturnsArrayNotNull(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	newCh32Handler().ListAccumulationTrajectories(w, httptest.NewRequest(http.MethodGet, "/v1/accumulation/trajectories", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	body := bytes.TrimSpace(w.Body.Bytes())
	if string(body) == "null" {
		t.Fatalf("list endpoint returned null instead of []")
	}
	if string(body) != "[]" {
		t.Fatalf("want [], got %s", string(body))
	}
}

func TestGetAccumulationTrajectory_MissingID(t *testing.T) {
	t.Parallel()
	h := newCh32Handler()
	req := httptest.NewRequest(http.MethodGet, "/v1/accumulation/trajectories/", nil)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()
	h.GetAccumulationTrajectory(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}
