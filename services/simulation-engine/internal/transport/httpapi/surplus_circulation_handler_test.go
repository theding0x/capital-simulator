package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// newSCHandler returns an httpapi.Handler wired to an in-memory store.
func newSCHandler(t *testing.T) *httpapi.Handler {
	t.Helper()
	mem := store.NewMemory()
	return httpapi.New(nil, httpapi.Deps{SurplusCirculations: mem})
}

func TestCreateSurplusCirculation_Created(t *testing.T) {
	t.Parallel()
	h := newSCHandler(t)
	body := `{"period":"1880","total_surplus_pence":1000000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/surplus-circulations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateSurplusCirculation(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["id"] == "" {
		t.Fatal("expected id in response")
	}
	if resp["period"] != "1880" {
		t.Fatalf("period: want 1880, got %v", resp["period"])
	}
	if w.Header().Get("Location") == "" {
		t.Fatal("expected Location header")
	}
}

func TestCreateSurplusCirculation_BadRequest_MissingPeriod(t *testing.T) {
	t.Parallel()
	h := newSCHandler(t)
	body := `{"total_surplus_pence":1000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/surplus-circulations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateSurplusCirculation(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateSurplusCirculation_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	h := newSCHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/surplus-circulations", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateSurplusCirculation(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetSurplusCirculation_NotFound(t *testing.T) {
	t.Parallel()
	h := newSCHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/reproduction/surplus-circulations/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()
	h.GetSurplusCirculation(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestListSurplusCirculations_NonNullSlice(t *testing.T) {
	t.Parallel()
	h := newSCHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/reproduction/surplus-circulations", nil)
	w := httptest.NewRecorder()
	h.ListSurplusCirculations(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	items, ok := resp["items"]
	if !ok {
		t.Fatal("expected items key")
	}
	if items == nil {
		t.Fatal("items must not be null")
	}
}

func TestRecordRealisationSource_Created(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	h := httpapi.New(nil, httpapi.Deps{SurplusCirculations: mem})

	// Create a circulation first.
	createBody := `{"period":"1880","total_surplus_pence":1000000}`
	cr := httptest.NewRequest(http.MethodPost, "/v1/reproduction/surplus-circulations", bytes.NewBufferString(createBody))
	cr.Header.Set("Content-Type", "application/json")
	cw := httptest.NewRecorder()
	h.CreateSurplusCirculation(cw, cr)
	if cw.Code != http.StatusCreated {
		t.Fatalf("create failed: %d %s", cw.Code, cw.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(cw.Body).Decode(&created)
	id := created["id"].(string)

	// Record a realisation source.
	srcBody := `{"source":"capitalist_consumption","pence":500000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/surplus-circulations/"+id+"/realisation-source", bytes.NewBufferString(srcBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.RecordRealisationSource(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["source"] != "capitalist_consumption" {
		t.Fatalf("source: got %v", resp["source"])
	}
}

func TestRecordRealisationSource_UnknownSource(t *testing.T) {
	t.Parallel()
	h := newSCHandler(t)
	srcBody := `{"source":"bad_source","pence":100}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproduction/surplus-circulations/5eed0000000000001701/realisation-source", bytes.NewBufferString(srcBody))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "5eed0000000000001701")
	w := httptest.NewRecorder()
	h.RecordRealisationSource(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetSocialAggregate_SeedPeriod(t *testing.T) {
	t.Parallel()
	h := newSCHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/reproduction/social-aggregate?period=1871", nil)
	w := httptest.NewRecorder()
	h.GetSocialAggregate(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp["period"] != "1871" {
		t.Fatalf("period: got %v", resp["period"])
	}
}

func TestGetSocialAggregate_MissingPeriod(t *testing.T) {
	t.Parallel()
	h := newSCHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/reproduction/social-aggregate", nil)
	w := httptest.NewRecorder()
	h.GetSocialAggregate(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListRealisationPuzzles_SeedReturned(t *testing.T) {
	t.Parallel()
	h := newSCHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/reproduction/realisation-puzzles", nil)
	w := httptest.NewRecorder()
	h.ListRealisationPuzzles(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.NewDecoder(w.Body).Decode(&resp)
	items, ok := resp["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatal("expected at least one realisation puzzle in seed")
	}
	puzzle := items[0].(map[string]any)
	if puzzle["resolved_in_chapter"] != float64(20) {
		t.Fatalf("resolved_in_chapter: got %v", puzzle["resolved_in_chapter"])
	}
}
