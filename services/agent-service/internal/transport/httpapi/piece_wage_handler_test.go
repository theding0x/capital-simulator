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

// § 12-hour day, 24 pieces, 144 farthings daily wage → piece price = 6, piece value = 12
func TestComputePiecePrice_OK(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)
	body := map[string]any{
		"daily_wage_farthings":        144,
		"day_value_product_farthings": 288,
		"normal_output":               24,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/piece-price", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ComputePiecePrice(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		PiecePrice int64 `json:"piece_price"`
		PieceValue int64 `json:"piece_value"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.PiecePrice != 6 {
		t.Errorf("piece_price = %d, want 6", resp.PiecePrice)
	}
	if resp.PieceValue != 12 {
		t.Errorf("piece_value = %d, want 12", resp.PieceValue)
	}
	if resp.PieceValue <= resp.PiecePrice {
		t.Errorf("piece_value %d must exceed piece_price %d", resp.PieceValue, resp.PiecePrice)
	}
}

func TestComputePiecePrice_BadRequest(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)

	t.Run("bad json", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodPost, "/v1/piece-price", bytes.NewReader([]byte("{bad}")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.ComputePiecePrice(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("zero normal_output", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"daily_wage_farthings":        144,
			"day_value_product_farthings": 288,
			"normal_output":               0,
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/v1/piece-price", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.ComputePiecePrice(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})
}

// § create and retrieve a piece-wage contract for an agent
func TestCreateAndGetPieceWage(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)
	agentID := "test-agent-001"

	// Create
	createBody := map[string]any{"price_pence": 6, "normal_output": 24}
	b, _ := json.Marshal(createBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/agents/"+agentID+"/piece-wages", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", agentID)
	rr := httptest.NewRecorder()
	h.CreatePieceWage(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body: %s", rr.Code, rr.Body.String())
	}

	// Get
	req2 := httptest.NewRequest(http.MethodGet, "/v1/agents/"+agentID+"/piece-wages", nil)
	req2.SetPathValue("id", agentID)
	rr2 := httptest.NewRecorder()
	h.GetPieceWage(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200; body: %s", rr2.Code, rr2.Body.String())
	}
	var pw struct {
		PricePence       int64 `json:"price_pence"`
		NormalOutput     int64 `json:"normal_output"`
		ImpliedDailyWage int64 `json:"implied_daily_wage"`
	}
	if err := json.NewDecoder(rr2.Body).Decode(&pw); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if pw.PricePence != 6 || pw.NormalOutput != 24 {
		t.Errorf("round-trip: got price=%d output=%d, want 6/24", pw.PricePence, pw.NormalOutput)
	}
	if pw.ImpliedDailyWage != 144 {
		t.Errorf("implied_daily_wage = %d, want 144 (6×24)", pw.ImpliedDailyWage)
	}
}

// § 404 when no piece-wage contract exists for agent
func TestGetPieceWage_NotFound(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/agents/no-such-agent/piece-wages", nil)
	req.SetPathValue("id", "no-such-agent")
	rr := httptest.NewRecorder()
	h.GetPieceWage(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
}

// § sub-contract sweating: head keeps spread = piece_rate - assistant_rate
func TestCreateAndGetSubContract(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)
	createBody := map[string]any{
		"head_labourer_id":     "head-x",
		"assistant_ids":        []string{"assist-a", "assist-b"},
		"piece_rate_pence":     6,
		"assistant_rate_pence": 4,
	}
	b, _ := json.Marshal(createBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/sub-contracts", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateSubContract(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201; body: %s", rr.Code, rr.Body.String())
	}
	var created struct {
		ID     string `json:"id"`
		Spread int64  `json:"spread"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	if created.Spread != 2 {
		t.Errorf("spread = %d, want 2", created.Spread)
	}

	// Get by ID
	req2 := httptest.NewRequest(http.MethodGet, "/v1/sub-contracts/"+created.ID, nil)
	req2.SetPathValue("id", created.ID)
	rr2 := httptest.NewRecorder()
	h.GetSubContract(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Fatalf("get status = %d, want 200; body: %s", rr2.Code, rr2.Body.String())
	}
}

// § wage_form_id derivation: price_pence derived as wf.Wage.DailyPence / normal_output.
// Fixture: daily_pence=144, normal_output=24 → price_pence=6.
func TestCreatePieceWage_DerivedFromWageForm(t *testing.T) {
	t.Parallel()

	mem := store.NewMemory()
	h := httpapi.New(mem, nil)
	agentID := "worker-pw-wf"

	wf := agent.WageForm{
		AgentID: agent.AgentID(agentID),
		Wage:    agent.Wage{DailyPence: 144, WorkingDayMinutes: 720},
		LabourPowerValue: agent.WageLabourValue{
			DailyPence:       144,
			NecessaryMinutes: 360,
		},
	}
	saved, err := mem.CreateWageForm(t.Context(), wf)
	if err != nil {
		t.Fatalf("seed wage form: %v", err)
	}

	body := map[string]any{
		"wage_form_id":  string(saved.ID),
		"normal_output": 24,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/agents/"+agentID+"/piece-wages", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", agentID)
	rr := httptest.NewRecorder()
	h.CreatePieceWage(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		PricePence int64 `json:"price_pence"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.PricePence != 6 {
		t.Errorf("price_pence = %d, want 6 (144/24, derived from wage form)", resp.PricePence)
	}
}

// § 400 when SubContract spread is inverted (assistant earns more than head)
func TestCreateSubContract_InvalidSpread(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)
	body := map[string]any{
		"head_labourer_id":     "head-x",
		"assistant_ids":        []string{},
		"piece_rate_pence":     3,
		"assistant_rate_pence": 6,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/sub-contracts", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateSubContract(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
}
