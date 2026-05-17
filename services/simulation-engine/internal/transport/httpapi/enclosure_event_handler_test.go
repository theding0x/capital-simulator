package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

const sutherlandBody = `{
  "period": "1814–1820",
  "acres_enclosed": 794000,
  "population_displaced": 15000,
  "beneficiary": "Countess of Sutherland"
}`

func TestCreateEnclosureEvent_Success(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, st, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/enclosure-events", strings.NewReader(sutherlandBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateEnclosureEvent(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/v1/enclosure-events/") {
		t.Errorf("missing or wrong Location header: %q", loc)
	}
	var resp struct {
		ID                  string `json:"id"`
		AcresEnclosed       int64  `json:"acres_enclosed"`
		PopulationDisplaced int64  `json:"population_displaced"`
		Beneficiary         string `json:"beneficiary"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Error("response missing id")
	}
	if resp.AcresEnclosed != 794000 {
		t.Errorf("acres_enclosed = %d, want 794000", resp.AcresEnclosed)
	}
	if resp.PopulationDisplaced != 15000 {
		t.Errorf("population_displaced = %d, want 15000", resp.PopulationDisplaced)
	}
	if resp.Beneficiary != "Countess of Sutherland" {
		t.Errorf("beneficiary = %q, want Countess of Sutherland", resp.Beneficiary)
	}
}

func TestCreateEnclosureEvent_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, store.NewMemory(), nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/v1/enclosure-events", strings.NewReader(`{nope`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateEnclosureEvent(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateEnclosureEvent_BadRequest_ZeroAcres(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, store.NewMemory(), nil, nil, nil, nil)
	body := `{"period":"1450–1540","acres_enclosed":0,"population_displaced":40000,"beneficiary":"landed gentry"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/enclosure-events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateEnclosureEvent(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateEnclosureEvent_BadRequest_MissingPeriod(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, store.NewMemory(), nil, nil, nil, nil)
	body := `{"period":"","acres_enclosed":500000,"population_displaced":40000,"beneficiary":"landed gentry"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/enclosure-events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateEnclosureEvent(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestListEnclosureEvents_ReturnsNonNilSlice(t *testing.T) {
	t.Parallel()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, store.NewMemory(), nil, nil, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/enclosure-events", nil)
	w := httptest.NewRecorder()
	h.ListEnclosureEvents(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	if got := strings.TrimSpace(w.Body.String()); got == "null" {
		t.Errorf("response is null, want []: %q", got)
	}
}

func TestListEnclosureEvents_ReturnsCreated(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, nil, nil, nil, nil, nil, st, nil, nil, nil, nil)

	createReq := httptest.NewRequest(http.MethodPost, "/v1/enclosure-events",
		strings.NewReader(`{"period":"1760–1820","acres_enclosed":3000000,"population_displaced":150000,"beneficiary":"improving landlords"}`))
	createReq.Header.Set("Content-Type", "application/json")
	cw := httptest.NewRecorder()
	h.CreateEnclosureEvent(cw, createReq)
	if cw.Code != http.StatusCreated {
		t.Fatalf("create want 201, got %d: %s", cw.Code, cw.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/enclosure-events", nil)
	w := httptest.NewRecorder()
	h.ListEnclosureEvents(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var got []map[string]any
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 event, got %d", len(got))
	}
	if got[0]["beneficiary"] != "improving landlords" {
		t.Errorf("beneficiary = %v, want improving landlords", got[0]["beneficiary"])
	}
}