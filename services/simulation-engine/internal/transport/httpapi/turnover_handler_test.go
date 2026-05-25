// Vol. II Ch. 7 — handler tests for turnover time and number of turnovers.
package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
	tv "github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

func newTurnoverHandler() *httpapi.Handler {
	m := store.NewMemory()
	return httpapi.New(nil, httpapi.Deps{Turnovers: m})
}

func postTurnover(t *testing.T, h *httpapi.Handler, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/turnovers", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateTurnover(rr, req)
	return rr
}

// TestCreateTurnover_201_MoneyLens verifies the spinning-mill productive-lens case creates successfully.
func TestCreateTurnover_201_MoneyLens(t *testing.T) {
	t.Parallel()
	h := newTurnoverHandler()
	rr := postTurnover(t, h, map[string]any{
		"industrial_capital_id": "5eed000000000000000401",
		"lens":                  "productive",
		"turnover_time_minutes": 10080,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rr.Code, rr.Body)
	}
	if loc := rr.Header().Get("Location"); loc == "" {
		t.Fatal("want Location header, got empty")
	}
	var got tv.Turnover
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID == "" {
		t.Fatal("want non-empty ID in response")
	}
	if got.Lens != tv.LensProductive {
		t.Fatalf("want productive lens, got %q", got.Lens)
	}
}

// TestCreateTurnover_400_CommodityLens verifies Form-III (commodity) lens is rejected.
func TestCreateTurnover_400_CommodityLens(t *testing.T) {
	t.Parallel()
	h := newTurnoverHandler()
	rr := postTurnover(t, h, map[string]any{
		"industrial_capital_id": "",
		"lens":                  "commodity",
		"turnover_time_minutes": 10080,
	})
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rr.Code)
	}
}

// TestGetTurnover_404 verifies a missing turnover returns 404.
func TestGetTurnover_404(t *testing.T) {
	t.Parallel()
	h := newTurnoverHandler()
	req := httptest.NewRequest(http.MethodGet, "/v1/turnovers/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	rr := httptest.NewRecorder()
	h.GetTurnover(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rr.Code)
	}
}

// TestListTurnovers_NonNullItems verifies the list endpoint returns a non-null items array.
func TestListTurnovers_NonNullItems(t *testing.T) {
	t.Parallel()
	h := newTurnoverHandler()
	req := httptest.NewRequest(http.MethodGet, "/v1/turnovers", nil)
	rr := httptest.NewRecorder()
	h.ListTurnovers(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	var resp map[string]json.RawMessage
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var items []tv.Turnover
	if err := json.Unmarshal(resp["items"], &items); err != nil {
		t.Fatalf("items unmarshal: %v", err)
	}
	if items == nil {
		t.Fatal("items must not be null")
	}
}

// TestRecordCycle_400_Incomplete verifies that a cycle where returned < advance is rejected.
func TestRecordCycle_400_Incomplete(t *testing.T) {
	t.Parallel()
	h := newTurnoverHandler()
	rr := postTurnover(t, h, map[string]any{
		"lens":                  "money",
		"turnover_time_minutes": 777600,
	})
	var created tv.Turnover
	json.NewDecoder(rr.Body).Decode(&created)

	body, _ := json.Marshal(map[string]any{
		"started_at":          "1870-01-01T00:00:00Z",
		"ended_at":            "1871-07-01T00:00:00Z",
		"advance_pence":       500000,
		"returned_pence":      400000, // less than advance — incomplete
		"production_minutes":  720000,
		"circulation_minutes": 57600,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/turnovers/"+string(created.ID)+"/cycles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", string(created.ID))
	rr2 := httptest.NewRecorder()
	h.RecordCycle(rr2, req)
	if rr2.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", rr2.Code, rr2.Body)
	}
}

// TestRecomputeNumber_400_NoCycles verifies that recomputing with no cycles returns 400.
func TestRecomputeNumber_400_NoCycles(t *testing.T) {
	t.Parallel()
	h := newTurnoverHandler()
	rr := postTurnover(t, h, map[string]any{
		"lens":                  "productive",
		"turnover_time_minutes": 2102400,
	})
	var created tv.Turnover
	json.NewDecoder(rr.Body).Decode(&created)

	req := httptest.NewRequest(http.MethodPost, "/v1/turnovers/"+string(created.ID)+"/recompute-number", nil)
	req.SetPathValue("id", string(created.ID))
	rr2 := httptest.NewRecorder()
	h.RecomputeNumber(rr2, req)
	if rr2.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d: %s", rr2.Code, rr2.Body)
	}
}

// TestTurnover_RoundTrip_SpinningMill recreates Marx's spinning-mill fixture end-to-end:
// weekly turnover, t = 10080 min, n ≈ 52/year, BasisPoints = 521428.
func TestTurnover_RoundTrip_SpinningMill(t *testing.T) {
	t.Parallel()
	h := newTurnoverHandler()

	// Create turnover.
	rr := postTurnover(t, h, map[string]any{
		"industrial_capital_id": "5eed000000000000000401",
		"lens":                  "productive",
		"turnover_time_minutes": 10080,
	})
	if rr.Code != http.StatusCreated {
		t.Fatalf("create: want 201, got %d", rr.Code)
	}
	var turnover tv.Turnover
	json.NewDecoder(rr.Body).Decode(&turnover)

	// Record two complete cycles (production 8000 + circulation 2080 = 10080 min each).
	for i := 0; i < 2; i++ {
		start := time.Date(1871, 1, 1+i*7, 0, 0, 0, 0, time.UTC)
		end := start.Add(7 * 24 * time.Hour)
		body, _ := json.Marshal(map[string]any{
			"started_at":          start.Format(time.RFC3339),
			"ended_at":            end.Format(time.RFC3339),
			"advance_pence":       48000,
			"returned_pence":      57600,
			"production_minutes":  8000,
			"circulation_minutes": 2080,
		})
		req := httptest.NewRequest(http.MethodPost, "/v1/turnovers/"+string(turnover.ID)+"/cycles", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", string(turnover.ID))
		rr2 := httptest.NewRecorder()
		h.RecordCycle(rr2, req)
		if rr2.Code != http.StatusCreated {
			t.Fatalf("cycle %d: want 201, got %d: %s", i, rr2.Code, rr2.Body)
		}
	}

	// Recompute number.
	req := httptest.NewRequest(http.MethodPost, "/v1/turnovers/"+string(turnover.ID)+"/recompute-number", nil)
	req.SetPathValue("id", string(turnover.ID))
	rr3 := httptest.NewRecorder()
	h.RecomputeNumber(rr3, req)
	if rr3.Code != http.StatusOK {
		t.Fatalf("recompute: want 200, got %d: %s", rr3.Code, rr3.Body)
	}
	var tn tv.TurnoverNumber
	json.NewDecoder(rr3.Body).Decode(&tn)
	if tn.BasisPoints != 521428 {
		t.Fatalf("spinning mill: want BasisPoints 521428, got %d", tn.BasisPoints)
	}

	// GetNumber confirms the stored value.
	req2 := httptest.NewRequest(http.MethodGet, "/v1/turnovers/"+string(turnover.ID)+"/number", nil)
	req2.SetPathValue("id", string(turnover.ID))
	rr4 := httptest.NewRecorder()
	h.GetNumber(rr4, req2)
	if rr4.Code != http.StatusOK {
		t.Fatalf("get number: want 200, got %d", rr4.Code)
	}
	var tn2 tv.TurnoverNumber
	json.NewDecoder(rr4.Body).Decode(&tn2)
	if tn2.BasisPoints != 521428 {
		t.Fatalf("get number: want BasisPoints 521428, got %d", tn2.BasisPoints)
	}
}

// TestCreateTurnover_400_BadJSON verifies malformed JSON returns 400.
func TestCreateTurnover_400_BadJSON(t *testing.T) {
	t.Parallel()
	h := newTurnoverHandler()
	req := httptest.NewRequest(http.MethodPost, "/v1/turnovers", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateTurnover(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rr.Code)
	}
}
