package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Vol. II Ch. 1 — The Circuit of Money-Capital.

type moneyAdvanceResponse struct {
	Amount     int64  `json:"amount"`
	AgentID    string `json:"agent_id,omitempty"`
	AdvancedAt string `json:"advanced_at"`
}

type labourLegResponse struct {
	Amount           int64 `json:"amount"`
	LabourPowerHours int64 `json:"labour_power_hours"`
}

type meansLegResponse struct {
	Amount             int64 `json:"amount"`
	MeansCapacityHours int64 `json:"means_capacity_hours"`
}

type purchasePhaseResponse struct {
	MoneyCircuitID string            `json:"money_circuit_id"`
	LabourLeg      labourLegResponse `json:"labour_leg"`
	MeansLeg       meansLegResponse  `json:"means_leg"`
	Total          int64             `json:"total"`
}

type productiveStateResponse struct {
	MoneyCircuitID string `json:"money_circuit_id"`
	EnteredAt      string `json:"entered_at"`
	ConstantPence  int64  `json:"constant_pence"`
	VariablePence  int64  `json:"variable_pence"`
	TotalAdvance   int64  `json:"total_advance"`
}

type commodityCapitalResponse struct {
	MoneyCircuitID string `json:"money_circuit_id"`
	CommodityID    string `json:"commodity_id,omitempty"`
	ValueOriginal  int64  `json:"value_original"`
	ValueSurplus   int64  `json:"value_surplus"`
	Total          int64  `json:"total"`
}

type realisationResponse struct {
	MoneyCircuitID  string `json:"money_circuit_id"`
	SoldAt          string `json:"sold_at"`
	RealisedPence   int64  `json:"realised_pence"`
	SurplusRealised int64  `json:"surplus_realised"`
}

type moneyCircuitResponse struct {
	ID          string                    `json:"id"`
	Advance     moneyAdvanceResponse      `json:"advance"`
	Moment      string                    `json:"moment"`
	Purchase    *purchasePhaseResponse    `json:"purchase,omitempty"`
	Productive  *productiveStateResponse  `json:"productive,omitempty"`
	Commodity   *commodityCapitalResponse `json:"commodity,omitempty"`
	Realisation *realisationResponse      `json:"realisation,omitempty"`
	CreatedAt   string                    `json:"created_at"`
}

func toMoneyCircuitResponse(mc circulation.MoneyCircuit) moneyCircuitResponse {
	resp := moneyCircuitResponse{
		ID: string(mc.ID),
		Advance: moneyAdvanceResponse{
			Amount:     int64(mc.Advance.Amount),
			AgentID:    mc.Advance.AgentID,
			AdvancedAt: mc.Advance.AdvancedAt.UTC().Format(time.RFC3339),
		},
		Moment:    string(mc.Moment),
		CreatedAt: mc.CreatedAt.UTC().Format(time.RFC3339),
	}
	if mc.Purchase != nil {
		p := mc.Purchase
		resp.Purchase = &purchasePhaseResponse{
			MoneyCircuitID: string(p.MoneyCircuitID),
			LabourLeg:      labourLegResponse{Amount: int64(p.LabourLeg.Amount), LabourPowerHours: p.LabourLeg.LabourPowerHours},
			MeansLeg:       meansLegResponse{Amount: int64(p.MeansLeg.Amount), MeansCapacityHours: p.MeansLeg.MeansCapacityHours},
			Total:          int64(p.Total()),
		}
	}
	if mc.Productive != nil {
		ps := mc.Productive
		resp.Productive = &productiveStateResponse{
			MoneyCircuitID: string(ps.MoneyCircuitID),
			EnteredAt:      ps.EnteredAt.UTC().Format(time.RFC3339),
			ConstantPence:  int64(ps.ConstantPence),
			VariablePence:  int64(ps.VariablePence),
			TotalAdvance:   int64(ps.TotalAdvance()),
		}
	}
	if mc.Commodity != nil {
		cc := mc.Commodity
		resp.Commodity = &commodityCapitalResponse{
			MoneyCircuitID: string(cc.MoneyCircuitID),
			CommodityID:    cc.CommodityID,
			ValueOriginal:  int64(cc.ValueOriginal),
			ValueSurplus:   int64(cc.ValueSurplus),
			Total:          int64(cc.Total()),
		}
	}
	if mc.Realisation != nil {
		rl := mc.Realisation
		resp.Realisation = &realisationResponse{
			MoneyCircuitID:  string(rl.MoneyCircuitID),
			SoldAt:          rl.SoldAt.UTC().Format(time.RFC3339),
			RealisedPence:   int64(rl.RealisedPence),
			SurplusRealised: int64(mc.SurplusRealised()),
		}
	}
	return resp
}

// CreateMoneyCircuit handles POST /v1/money-circuits.
func (h *Handler) CreateMoneyCircuit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AgentID    string `json:"agent_id,omitempty"`
		Amount     int64  `json:"amount"`
		AdvancedAt string `json:"advanced_at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	advance := circulation.MoneyAdvance{
		Amount:  circulation.Pence(req.Amount),
		AgentID: req.AgentID,
	}
	if req.AdvancedAt != "" {
		t, err := time.Parse(time.RFC3339, req.AdvancedAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid advanced_at: "+err.Error())
			return
		}
		advance.AdvancedAt = t
	}
	mc := circulation.MoneyCircuit{Advance: advance}
	created, err := h.MoneyCircuits.CreateMoneyCircuit(r.Context(), mc)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/money-circuits/"+string(created.ID))
	writeJSON(w, http.StatusCreated, toMoneyCircuitResponse(created))
}

// GetMoneyCircuit handles GET /v1/money-circuits/{id}.
func (h *Handler) GetMoneyCircuit(w http.ResponseWriter, r *http.Request) {
	id := circulation.MoneyCircuitID(r.PathValue("id"))
	mc, err := h.MoneyCircuits.GetMoneyCircuit(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "money circuit not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toMoneyCircuitResponse(mc))
}

// ListMoneyCircuits handles GET /v1/money-circuits.
func (h *Handler) ListMoneyCircuits(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent_id")
	moment := circulation.CircuitMoment(r.URL.Query().Get("moment"))
	list, err := h.MoneyCircuits.ListMoneyCircuits(r.Context(), agentID, moment)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	out := make([]moneyCircuitResponse, len(list))
	for i, mc := range list {
		out[i] = toMoneyCircuitResponse(mc)
	}
	writeJSON(w, http.StatusOK, out)
}

// RecordMoneyPurchase handles POST /v1/money-circuits/{id}/purchase.
func (h *Handler) RecordMoneyPurchase(w http.ResponseWriter, r *http.Request) {
	id := circulation.MoneyCircuitID(r.PathValue("id"))
	var req struct {
		LabourAmount       int64 `json:"labour_amount"`
		LabourPowerHours   int64 `json:"labour_power_hours"`
		MeansAmount        int64 `json:"means_amount"`
		MeansCapacityHours int64 `json:"means_capacity_hours"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.LabourAmount <= 0 || req.MeansAmount <= 0 {
		writeError(w, http.StatusBadRequest, "labour_amount and means_amount must be positive")
		return
	}
	p := circulation.PurchasePhase{
		LabourLeg: circulation.LabourLeg{Amount: circulation.Pence(req.LabourAmount), LabourPowerHours: req.LabourPowerHours},
		MeansLeg:  circulation.MeansLeg{Amount: circulation.Pence(req.MeansAmount), MeansCapacityHours: req.MeansCapacityHours},
	}
	mc, err := h.MoneyCircuits.RecordPurchase(r.Context(), id, p)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "money circuit not found")
			return
		}
		if isMoneyCircuitValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toMoneyCircuitResponse(mc))
}

// RecordMoneyProduce handles POST /v1/money-circuits/{id}/produce.
func (h *Handler) RecordMoneyProduce(w http.ResponseWriter, r *http.Request) {
	id := circulation.MoneyCircuitID(r.PathValue("id"))
	var req struct {
		ConstantPence int64  `json:"constant_pence"`
		VariablePence int64  `json:"variable_pence"`
		EnteredAt     string `json:"entered_at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ConstantPence < 0 || req.VariablePence < 0 {
		writeError(w, http.StatusBadRequest, "constant_pence and variable_pence cannot be negative")
		return
	}
	ps := circulation.ProductiveState{
		ConstantPence: circulation.Pence(req.ConstantPence),
		VariablePence: circulation.Pence(req.VariablePence),
	}
	if req.EnteredAt != "" {
		t, err := time.Parse(time.RFC3339, req.EnteredAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid entered_at: "+err.Error())
			return
		}
		ps.EnteredAt = t
	}
	mc, err := h.MoneyCircuits.RecordProductive(r.Context(), id, ps)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "money circuit not found")
			return
		}
		if isMoneyCircuitValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toMoneyCircuitResponse(mc))
}

// RecordMoneyCommodity handles POST /v1/money-circuits/{id}/commodity.
func (h *Handler) RecordMoneyCommodity(w http.ResponseWriter, r *http.Request) {
	id := circulation.MoneyCircuitID(r.PathValue("id"))
	var req struct {
		CommodityID   string `json:"commodity_id,omitempty"`
		ValueOriginal int64  `json:"value_original"`
		ValueSurplus  int64  `json:"value_surplus"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ValueOriginal <= 0 {
		writeError(w, http.StatusBadRequest, "value_original must be positive")
		return
	}
	cc := circulation.CommodityCapital{
		CommodityID:   req.CommodityID,
		ValueOriginal: circulation.Pence(req.ValueOriginal),
		ValueSurplus:  circulation.Pence(req.ValueSurplus),
	}
	mc, err := h.MoneyCircuits.RecordCommodity(r.Context(), id, cc)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "money circuit not found")
			return
		}
		if isMoneyCircuitValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toMoneyCircuitResponse(mc))
}

// RecordMoneyRealise handles POST /v1/money-circuits/{id}/realise.
func (h *Handler) RecordMoneyRealise(w http.ResponseWriter, r *http.Request) {
	id := circulation.MoneyCircuitID(r.PathValue("id"))
	var req struct {
		RealisedPence int64  `json:"realised_pence"`
		SoldAt        string `json:"sold_at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.RealisedPence <= 0 {
		writeError(w, http.StatusBadRequest, "realised_pence must be positive")
		return
	}
	rl := circulation.Realisation{RealisedPence: circulation.Pence(req.RealisedPence)}
	if req.SoldAt != "" {
		t, err := time.Parse(time.RFC3339, req.SoldAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid sold_at: "+err.Error())
			return
		}
		rl.SoldAt = t
	}
	mc, err := h.MoneyCircuits.RecordRealisation(r.Context(), id, rl)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "money circuit not found")
			return
		}
		if isMoneyCircuitValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toMoneyCircuitResponse(mc))
}

func isMoneyCircuitValidationError(err error) bool {
	return errors.Is(err, circulation.ErrLabourMeansMismatch) ||
		errors.Is(err, circulation.ErrNoProductionPhase) ||
		errors.Is(err, circulation.ErrInvalidPhaseTransition) ||
		errors.Is(err, circulation.ErrMagnitudeNotPreserved) ||
		isValidationError(err)
}
