package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// §2: ValueComposition{8000,2000} → ratio=0.8
func TestComputeOrganicComposition_Success(t *testing.T) {
	t.Parallel()
	body := `{"constant_capital":8000,"variable_capital":2000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/organic-composition", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.ComputeOrganicComposition(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Ratio float64 `json:"ratio"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Ratio != 0.8 {
		t.Errorf("ratio = %v, want 0.8", resp.Ratio)
	}
}

func TestComputeOrganicComposition_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/organic-composition", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.ComputeOrganicComposition(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestComputeOrganicComposition_BadRequest_ZeroVariable(t *testing.T) {
	t.Parallel()
	body := `{"constant_capital":8000,"variable_capital":0}`
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/organic-composition", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.ComputeOrganicComposition(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

// §1: total=10000, OC=0.8, wage=200 → workers=10
func TestComputeLabourDemand_Success(t *testing.T) {
	t.Parallel()
	body := `{"total_capital":10000,"organic_composition_ratio":0.8,"wage_pence":200}`
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/labour-demand", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.ComputeLabourDemandEndpoint(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Workers int64 `json:"workers"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// 10000 * (1-0.8) / 200 = 10
	if resp.Workers != 10 {
		t.Errorf("workers = %d, want 10", resp.Workers)
	}
}

func TestComputeLabourDemand_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/labour-demand", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.ComputeLabourDemandEndpoint(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

// supply=50, demanded=40 → size=10, proportion=0.25
func TestComputeReserveArmy_Success(t *testing.T) {
	t.Parallel()
	body := `{"worker_supply":50,"workers_demanded":40}`
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/reserve-army", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.ComputeReserveArmyEndpoint(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		ReserveArmySize    int64   `json:"reserve_army_size"`
		RelativeProportion float64 `json:"relative_proportion"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ReserveArmySize != 10 {
		t.Errorf("reserve_army_size = %d, want 10", resp.ReserveArmySize)
	}
	want := float64(10) / float64(40)
	if resp.RelativeProportion != want {
		t.Errorf("relative_proportion = %v, want %v", resp.RelativeProportion, want)
	}
}

func TestComputeReserveArmy_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/reserve-army", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.ComputeReserveArmyEndpoint(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

// §1 scenario: create and verify 201 + Location header + series.
func TestCreateGeneralLawScenario_Success(t *testing.T) {
	t.Parallel()
	body := `{
		"name":"§1 test",
		"constant_capital":8000,
		"variable_capital":2000,
		"surplus_rate":1.0,
		"accumulation_rate":1.0,
		"productivity_growth":0.0,
		"wage_pence":200,
		"worker_supply":50,
		"periods":3
	}`
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/scenarios", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.CreateGeneralLawScenario(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	loc := w.Header().Get("Location")
	if !strings.HasPrefix(loc, "/v1/accumulation/scenarios/") {
		t.Errorf("Location = %q, want /v1/accumulation/scenarios/<id>", loc)
	}
	var resp struct {
		ID     string `json:"id"`
		Series []struct {
			Period int64 `json:"period"`
		} `json:"series"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Error("id should be non-empty")
	}
	if len(resp.Series) != 3 {
		t.Errorf("series len = %d, want 3", len(resp.Series))
	}
}

func TestCreateGeneralLawScenario_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/scenarios", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.CreateGeneralLawScenario(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateGeneralLawScenario_BadRequest_ZeroVariable(t *testing.T) {
	t.Parallel()
	body := `{"name":"x","constant_capital":8000,"variable_capital":0,"surplus_rate":1.0,"accumulation_rate":1.0,"productivity_growth":0.0,"wage_pence":200,"worker_supply":50,"periods":3}`
	req := httptest.NewRequest(http.MethodPost, "/v1/accumulation/scenarios", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.CreateGeneralLawScenario(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestGetGeneralLawScenario_NotFound(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/v1/accumulation/scenarios/nope", nil)
	req.SetPathValue("id", "nope")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: store.NewMemory(), HistoricalStages: store.NewMemory(), EnclosureEvents: store.NewMemory()})
	h.GetGeneralLawScenario(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

// Round-trip: create then get returns identical scenario with non-nil series.
func TestGetGeneralLawScenario_RoundTrip(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, httpapi.Deps{GeneralLaw: st, HistoricalStages: st, EnclosureEvents: st})

	createBody := `{
		"name":"§2 rising OC",
		"constant_capital":8000,
		"variable_capital":2000,
		"surplus_rate":1.0,
		"accumulation_rate":1.0,
		"productivity_growth":0.05,
		"wage_pence":200,
		"worker_supply":50,
		"periods":5
	}`
	createReq := httptest.NewRequest(http.MethodPost, "/v1/accumulation/scenarios", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	h.CreateGeneralLawScenario(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("create: want 201, got %d: %s", createW.Code, createW.Body.String())
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(createW.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/v1/accumulation/scenarios/"+created.ID, nil)
	getReq.SetPathValue("id", created.ID)
	getW := httptest.NewRecorder()
	h.GetGeneralLawScenario(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Fatalf("get: want 200, got %d: %s", getW.Code, getW.Body.String())
	}
	var got struct {
		ID     string `json:"id"`
		Series []struct {
			OrganicComposition float64 `json:"organic_composition"`
		} `json:"series"`
	}
	if err := json.NewDecoder(getW.Body).Decode(&got); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("id = %q, want %q", got.ID, created.ID)
	}
	if len(got.Series) != 5 {
		t.Errorf("series len = %d, want 5", len(got.Series))
	}
	// With rising OC, organic_composition should increase each period.
	for i := 1; i < len(got.Series); i++ {
		if got.Series[i].OrganicComposition <= got.Series[i-1].OrganicComposition {
			t.Errorf("period %d OC %v <= period %d OC %v (should rise)",
				i+1, got.Series[i].OrganicComposition, i, got.Series[i-1].OrganicComposition)
		}
	}
}
