package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// economyRequest is the body for POST /v1/profit/economy. The economy applies to
// a capital decomposition — the constant capital c, the variable capital v and
// the surplus-value s, all LabourMinutes — and Saving is the economy realised in
// c. Kind names the form the economy takes (shared_conditions, waste_reduction,
// workday_extension, improved_machinery, raw_material_saving). The old rate of
// profit s/(c+v), the new rate s/(c−saving+v) and the gain are derived server-side.
type economyRequest struct {
	Kind         string `json:"kind"`
	Constant     int64  `json:"constant"`
	Variable     int64  `json:"variable"`
	SurplusValue int64  `json:"surplus_value"`
	Saving       int64  `json:"saving"`
}

// CreateEconomyAnalysis handles POST /v1/profit/economy — from an economy in
// constant capital, compute the rise it produces in the rate of profit s/(c+v)
// and persist.
func (h *Handler) CreateEconomyAnalysis(w http.ResponseWriter, r *http.Request) {
	var req economyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	kind := profit.EconomyKind(req.Kind)
	if !kind.Valid() {
		writeError(w, http.StatusUnprocessableEntity,
			"kind must be one of shared_conditions, waste_reduction, workday_extension, improved_machinery, raw_material_saving")
		return
	}
	if req.Constant < 0 || req.Variable < 0 || req.SurplusValue < 0 || req.Saving < 0 {
		writeError(w, http.StatusUnprocessableEntity, "constant, variable, surplus_value and saving must be non-negative")
		return
	}
	if req.Saving > req.Constant {
		writeError(w, http.StatusUnprocessableEntity, "saving cannot exceed constant capital")
		return
	}
	analysis := profit.ComputeEconomyAnalysis(profit.ConstantCapitalEconomy{
		Kind:            kind,
		ConstantCapital: profit.ConstantCapital(req.Constant),
		VariableCapital: profit.VariableCapital(req.Variable),
		SurplusValue:    profit.SurplusValue(req.SurplusValue),
		Saving:          profit.ConstantCapital(req.Saving),
	})
	saved, err := h.Store.CreateEconomyAnalysis(r.Context(), analysis)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/profit/economy/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetEconomyAnalysis handles GET /v1/profit/economy/{id}.
func (h *Handler) GetEconomyAnalysis(w http.ResponseWriter, r *http.Request) {
	id := profit.EconomyAnalysisID(r.PathValue("id"))
	a, err := h.Store.GetEconomyAnalysis(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// ListEconomyAnalyses handles GET /v1/profit/economy.
func (h *Handler) ListEconomyAnalyses(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListEconomyAnalyses(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []profit.EconomyAnalysis{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
