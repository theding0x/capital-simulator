package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/transport/httpapi"
)

// seedPieceWage installs the Spinning Mill baseline (6 farthings at 24 pieces a
// day → 144 farthings implied daily wage) for the given agent.
func seedPieceWage(t *testing.T, mem *store.Memory, agentID string) {
	t.Helper()
	if _, err := mem.CreatePieceWage(t.Context(), agent.PieceWage{
		AgentID:      agent.AgentID(agentID),
		PricePence:   6,
		NormalOutput: 24,
	}); err != nil {
		t.Fatalf("seed piece wage: %v", err)
	}
}

func reprice(t *testing.T, h *httpapi.Handler, agentID string, factor float64) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(map[string]any{"productivity_factor": factor})
	req := httptest.NewRequest(http.MethodPost, "/v1/piece-wages/"+agentID+"/reprice", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("agentID", agentID)
	rr := httptest.NewRecorder()
	h.RepricePieceWage(rr, req)
	return rr
}

// § Marx's law of piece-wages: at 2× productivity the daily wage is unchanged
// (144f) and the price per piece is exactly halved, 6f → 3f, output 24 → 48.
func TestRepricePieceWage_DoubleProductivityHalvesPrice(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	h := httpapi.New(mem, nil)
	const agentID = "spinning-mill-worker"
	seedPieceWage(t, mem, agentID)

	rr := reprice(t, h, agentID, 2.0)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		Current  agent.PiecePricePoint   `json:"current"`
		Appended bool                    `json:"appended"`
		Points   []agent.PiecePricePoint `json:"points"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Current.PricePence != 3 {
		t.Errorf("price = %d, want 3 (halved from 6)", resp.Current.PricePence)
	}
	if resp.Current.NormalOutput != 48 {
		t.Errorf("output = %d, want 48 (doubled from 24)", resp.Current.NormalOutput)
	}
	if dw := resp.Current.PricePence * resp.Current.NormalOutput; dw != 144 {
		t.Errorf("implied daily wage = %d, want 144 (unchanged)", dw)
	}
	if !resp.Appended {
		t.Error("expected appended = true for a first reprice")
	}
	if len(resp.Points) != 1 {
		t.Errorf("points = %d, want 1", len(resp.Points))
	}
}

// § the scheduler reprices every pass; an unchanged productivity factor must not
// spam the series — the second identical reprice is a no-op.
func TestRepricePieceWage_DedupOnUnchangedFactor(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	h := httpapi.New(mem, nil)
	const agentID = "dedup-worker"
	seedPieceWage(t, mem, agentID)

	if rr := reprice(t, h, agentID, 1.5); rr.Code != http.StatusOK {
		t.Fatalf("first reprice status = %d; body: %s", rr.Code, rr.Body.String())
	}
	rr := reprice(t, h, agentID, 1.5)
	if rr.Code != http.StatusOK {
		t.Fatalf("second reprice status = %d; body: %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		Appended bool                    `json:"appended"`
		Points   []agent.PiecePricePoint `json:"points"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Appended {
		t.Error("expected appended = false on an unchanged factor")
	}
	if len(resp.Points) != 1 {
		t.Errorf("points = %d, want 1 (no duplicate appended)", len(resp.Points))
	}
}

func TestRepricePieceWage_NotFound(t *testing.T) {
	t.Parallel()
	h := httpapi.New(store.NewMemory(), nil)
	rr := reprice(t, h, "no-such-agent", 2.0)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
}

func TestRepricePieceWage_BadRequest(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	h := httpapi.New(mem, nil)
	const agentID = "bad-req-worker"
	seedPieceWage(t, mem, agentID)

	t.Run("malformed json", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodPost, "/v1/piece-wages/"+agentID+"/reprice", bytes.NewReader([]byte("{bad}")))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("agentID", agentID)
		rr := httptest.NewRecorder()
		h.RepricePieceWage(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("non-positive factor", func(t *testing.T) {
		t.Parallel()
		rr := reprice(t, h, agentID, 0)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})
}

// § an empty series must marshal as [] not null so the panel can render it.
func TestListPiecePricePoints_EmptyIsNonNull(t *testing.T) {
	t.Parallel()
	h := httpapi.New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/piece-wages/empty/price-points", nil)
	req.SetPathValue("agentID", "empty")
	rr := httptest.NewRecorder()
	h.ListPiecePricePoints(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var points []agent.PiecePricePoint
	if err := json.NewDecoder(rr.Body).Decode(&points); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if points == nil {
		t.Error("expected non-nil empty slice, got null")
	}
}

func TestListPieceWages_ReturnsSeeded(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	h := httpapi.New(mem, nil)
	seedPieceWage(t, mem, "listed-worker")

	req := httptest.NewRequest(http.MethodGet, "/v1/piece-wages", nil)
	rr := httptest.NewRecorder()
	h.ListPieceWages(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var pws []struct {
		AgentID    string `json:"agent_id"`
		PricePence int64  `json:"price_pence"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&pws); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(pws) != 1 || pws[0].AgentID != "listed-worker" || pws[0].PricePence != 6 {
		t.Errorf("got %+v, want one piece-wage for listed-worker at 6f", pws)
	}
}
