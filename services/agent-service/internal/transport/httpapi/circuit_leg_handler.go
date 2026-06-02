package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// Ch. 4 cross-service settlement (#216). market-service POSTs a completed
// exchange here; the handler records the two complementary circuit legs — a
// sale on the giver and a purchase on the receiver.

type recordSettlementRequest struct {
	ExchangeID          string `json:"exchange_id"`
	GiverID             string `json:"giver_id"`
	ReceiverID          string `json:"receiver_id"`
	GiverCommodityID    string `json:"giver_commodity_id"`
	ReceiverCommodityID string `json:"receiver_commodity_id"`
	RealisedValue       int64  `json:"realised_value"`
}

type settlementResponse struct {
	Legs []agent.CircuitLeg `json:"legs"`
}

// RecordExchangeSettlement handles POST /v1/circuit-legs.
func (h *Handler) RecordExchangeSettlement(w http.ResponseWriter, r *http.Request) {
	var req recordSettlementRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	settlement := agent.ExchangeSettlement{
		ExchangeID:          req.ExchangeID,
		GiverID:             agent.ID(req.GiverID),
		ReceiverID:          agent.ID(req.ReceiverID),
		GiverCommodityID:    req.GiverCommodityID,
		ReceiverCommodityID: req.ReceiverCommodityID,
		RealisedValue:       agent.LabourMinutes(req.RealisedValue),
	}
	if err := settlement.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	giver, receiver := settlement.Legs()
	savedGiver, err := h.CircuitLegStore.CreateCircuitLeg(r.Context(), giver)
	if err != nil {
		writeAppError(w, err)
		return
	}
	savedReceiver, err := h.CircuitLegStore.CreateCircuitLeg(r.Context(), receiver)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, settlementResponse{Legs: []agent.CircuitLeg{savedGiver, savedReceiver}})
}

// ListCircuitLegs handles GET /v1/agents/{id}/circuit-legs.
func (h *Handler) ListCircuitLegs(w http.ResponseWriter, r *http.Request) {
	agentID := agent.ID(r.PathValue("id"))
	legs, err := h.CircuitLegStore.ListCircuitLegs(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	if legs == nil {
		legs = []agent.CircuitLeg{}
	}
	writeJSON(w, http.StatusOK, legs)
}
