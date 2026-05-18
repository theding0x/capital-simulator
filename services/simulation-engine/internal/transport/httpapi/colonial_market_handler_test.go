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

// Ch. 33 — The Modern Theory of Colonisation HTTP layer.

func newCh33Handler() (*httpapi.Handler, *store.Memory) {
	mem := store.NewMemory()
	return httpapi.New(nil, httpapi.Deps{ColonialMarkets: mem}), mem
}

func TestCreateColonialMarket_PeelAtSwanRiver(t *testing.T) {
	t.Parallel()
	h, _ := newCh33Handler()
	body := `{"colony":"Swan River, 1829","free_labourers":300,"annual_wage_pence":4800,"land_available":true}`
	req := httptest.NewRequest(http.MethodPost, "/v1/colonial-markets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateColonialMarket(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); loc == "" {
		t.Errorf("Location header is empty")
	}
	var resp struct {
		ID                       string `json:"id"`
		Colony                   string `json:"colony"`
		FreeLabourers            int64  `json:"free_labourers"`
		WakefieldSchemeApplied   bool   `json:"wakefield_scheme_applied"`
		SurplusLabourExtractable bool   `json:"surplus_labour_extractable"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Errorf("ID was not assigned")
	}
	if resp.Colony != "Swan River, 1829" {
		t.Errorf("Colony = %q", resp.Colony)
	}
	if resp.FreeLabourers != 300 {
		t.Errorf("FreeLabourers = %d, want 300", resp.FreeLabourers)
	}
	if resp.WakefieldSchemeApplied {
		t.Errorf("WakefieldSchemeApplied = true, want false (no scheme yet)")
	}
	if resp.SurplusLabourExtractable {
		t.Errorf("SurplusLabourExtractable = true, want false (Peel's fiasco)")
	}
}

func TestCreateColonialMarket_DuplicateColonyConflicts(t *testing.T) {
	t.Parallel()
	h, mem := newCh33Handler()
	_, err := mem.CreateColonialLabourMarket(context.Background(), simulation.ColonialLabourMarket{
		Colony:          "Saint Domingo",
		FreeLabourers:   50,
		AnnualWagePence: 3600,
		LandAvailable:   true,
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	body := `{"colony":"saint domingo","free_labourers":100,"annual_wage_pence":3600,"land_available":true}`
	req := httptest.NewRequest(http.MethodPost, "/v1/colonial-markets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateColonialMarket(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("want 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateColonialMarket_MalformedJSON400(t *testing.T) {
	t.Parallel()
	h, _ := newCh33Handler()
	req := httptest.NewRequest(http.MethodPost, "/v1/colonial-markets", bytes.NewBufferString(`{"colony":`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateColonialMarket(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateColonialMarket_MissingColony400(t *testing.T) {
	t.Parallel()
	h, _ := newCh33Handler()
	body := `{"colony":"","free_labourers":300,"annual_wage_pence":4800,"land_available":true}`
	req := httptest.NewRequest(http.MethodPost, "/v1/colonial-markets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateColonialMarket(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestListColonialMarkets_ReturnsNonNilSlice(t *testing.T) {
	t.Parallel()
	h, _ := newCh33Handler()
	req := httptest.NewRequest(http.MethodGet, "/v1/colonial-markets", nil)
	w := httptest.NewRecorder()

	h.ListColonialMarkets(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if body != "[]\n" {
		t.Errorf("empty list body = %q, want []", body)
	}
}

func TestGetColonialMarket_NotFound(t *testing.T) {
	t.Parallel()
	h, _ := newCh33Handler()
	req := httptest.NewRequest(http.MethodGet, "/v1/colonial-markets/missing", nil)
	req.SetPathValue("id", "missing")
	w := httptest.NewRecorder()

	h.GetColonialMarket(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegulateColonialMarket_AppliesSufficientPrice(t *testing.T) {
	t.Parallel()
	h, mem := newCh33Handler()
	stored, err := mem.CreateColonialLabourMarket(context.Background(), simulation.ColonialLabourMarket{
		Colony:          "South Australia, 1836",
		FreeLabourers:   5000,
		AnnualWagePence: 4800,
		LandAvailable:   true,
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	body := `{"sufficient_price":{"price_per_acre_pence":240,"desired_acres":50}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/colonial-markets/"+string(stored.ID)+"/regulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", string(stored.ID))
	w := httptest.NewRecorder()

	h.RegulateColonialMarket(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		WakefieldSchemeApplied   bool  `json:"wakefield_scheme_applied"`
		IndependenceYears        int64 `json:"independence_years"`
		SurplusLabourExtractable bool  `json:"surplus_labour_extractable"`
		FreeLabourers            int64 `json:"free_labourers"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !resp.WakefieldSchemeApplied {
		t.Errorf("WakefieldSchemeApplied = false, want true")
	}
	if resp.IndependenceYears != 5 {
		t.Errorf("IndependenceYears = %d, want 5", resp.IndependenceYears)
	}
	if !resp.SurplusLabourExtractable {
		t.Errorf("SurplusLabourExtractable = false, want true")
	}
	if resp.FreeLabourers != 5000 {
		t.Errorf("FreeLabourers = %d, want 5000 (invariant)", resp.FreeLabourers)
	}
}

func TestRegulateColonialMarket_RejectsZeroPrice(t *testing.T) {
	t.Parallel()
	h, mem := newCh33Handler()
	stored, err := mem.CreateColonialLabourMarket(context.Background(), simulation.ColonialLabourMarket{
		Colony:          "Swan River, 1829",
		FreeLabourers:   300,
		AnnualWagePence: 4800,
		LandAvailable:   true,
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	body := `{"sufficient_price":{"price_per_acre_pence":0,"desired_acres":50}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/colonial-markets/"+string(stored.ID)+"/regulate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", string(stored.ID))
	w := httptest.NewRecorder()

	h.RegulateColonialMarket(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestComputeIndependence_HappyPath(t *testing.T) {
	t.Parallel()
	h, mem := newCh33Handler()
	stored, err := mem.CreateColonialLabourMarket(context.Background(), simulation.ColonialLabourMarket{
		Colony:          "South Australia, 1836",
		FreeLabourers:   5000,
		AnnualWagePence: 4800,
		LandAvailable:   true,
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	body := `{"worker_id":"settler-001","sufficient_price":{"price_per_acre_pence":240,"desired_acres":50}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/colonial-markets/"+string(stored.ID)+"/independence", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", string(stored.ID))
	w := httptest.NewRecorder()

	h.ComputeIndependence(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		YearsWorked       int64 `json:"years_worked"`
		SavingsPence      int64 `json:"savings_pence"`
		TargetRansomPence int64 `json:"target_ransom_pence"`
		BecameLandowner   bool  `json:"became_landowner"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.YearsWorked != 5 {
		t.Errorf("YearsWorked = %d, want 5", resp.YearsWorked)
	}
	if !resp.BecameLandowner {
		t.Errorf("BecameLandowner = false, want true")
	}
	if resp.TargetRansomPence != 12000 {
		t.Errorf("TargetRansomPence = %d, want 12000", resp.TargetRansomPence)
	}
}

func TestComputeIndependence_NotFound(t *testing.T) {
	t.Parallel()
	h, _ := newCh33Handler()
	body := `{"worker_id":"x","sufficient_price":{"price_per_acre_pence":240,"desired_acres":50}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/colonial-markets/missing/independence", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "missing")
	w := httptest.NewRecorder()

	h.ComputeIndependence(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestComputeIndependence_RequiresWorkerID(t *testing.T) {
	t.Parallel()
	h, mem := newCh33Handler()
	stored, err := mem.CreateColonialLabourMarket(context.Background(), simulation.ColonialLabourMarket{
		Colony:          "South Australia, 1836",
		FreeLabourers:   5000,
		AnnualWagePence: 4800,
		LandAvailable:   true,
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	body := `{"worker_id":"","sufficient_price":{"price_per_acre_pence":240,"desired_acres":50}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/colonial-markets/"+string(stored.ID)+"/independence", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", string(stored.ID))
	w := httptest.NewRecorder()

	h.ComputeIndependence(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}
