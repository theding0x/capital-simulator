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

// Ch. 2 fixture: the weaver (…0003) sells 20 yards of linen to the tailor
// (…0004), realised at 1200 labour-minutes. The settlement lands a sale leg on
// the weaver and a purchase leg on the tailor.
func TestRecordExchangeSettlement_OK(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	h := httpapi.New(mem, nil)
	body := map[string]any{
		"exchange_id":        "exch-001",
		"giver_id":           "5eed00000000000000000003",
		"receiver_id":        "5eed00000000000000000004",
		"giver_commodity_id": "000000000000000000000001",
		"realised_value":     1200,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/circuit-legs", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.RecordExchangeSettlement(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		Legs []agent.CircuitLeg `json:"legs"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Legs) != 2 {
		t.Fatalf("want 2 legs, got %d", len(resp.Legs))
	}
	if resp.Legs[0].Kind != agent.CircuitLegSale || resp.Legs[1].Kind != agent.CircuitLegPurchase {
		t.Errorf("legs = %s,%s; want sale,purchase", resp.Legs[0].Kind, resp.Legs[1].Kind)
	}
	// Each owner can list exactly their own leg — the counterparty's circuit is closed.
	for _, owner := range []string{"5eed00000000000000000003", "5eed00000000000000000004"} {
		legs, err := mem.ListCircuitLegs(t.Context(), agent.ID(owner))
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(legs) != 1 {
			t.Errorf("owner %s should have 1 leg, got %d", owner, len(legs))
		}
	}
}

func TestRecordExchangeSettlement_BadRequest(t *testing.T) {
	t.Parallel()
	h := httpapi.New(store.NewMemory(), nil)

	t.Run("bad json", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodPost, "/v1/circuit-legs", bytes.NewReader([]byte("{bad}")))
		rr := httptest.NewRecorder()
		h.RecordExchangeSettlement(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("self exchange", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"exchange_id": "x", "giver_id": "a1", "receiver_id": "a1",
			"giver_commodity_id": "c1", "realised_value": 100,
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/v1/circuit-legs", bytes.NewReader(b))
		rr := httptest.NewRecorder()
		h.RecordExchangeSettlement(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("non-positive value", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"exchange_id": "x", "giver_id": "a1", "receiver_id": "a2",
			"giver_commodity_id": "c1", "realised_value": 0,
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/v1/circuit-legs", bytes.NewReader(b))
		rr := httptest.NewRecorder()
		h.RecordExchangeSettlement(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})
}

func TestListCircuitLegs_Empty(t *testing.T) {
	t.Parallel()
	h := httpapi.New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/agents/nobody/circuit-legs", nil)
	req.SetPathValue("id", "nobody")
	rr := httptest.NewRecorder()
	h.ListCircuitLegs(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var legs []agent.CircuitLeg
	if err := json.NewDecoder(rr.Body).Decode(&legs); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if legs == nil {
		t.Error("body must be [] not null")
	}
}
