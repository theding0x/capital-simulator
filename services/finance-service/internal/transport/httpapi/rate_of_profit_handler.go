package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// profitRateRequest is the body for POST /v1/profit/rate. All magnitudes are
// LabourMinutes: the constant capital c, the variable capital v, and the
// surplus-value s of one capital decomposition.
type profitRateRequest struct {
	Constant     int64 `json:"constant"`
	Variable     int64 `json:"variable"`
	SurplusValue int64 `json:"surplus_value"`
}

// CreateProfitRate handles POST /v1/profit/rate — from a decomposition
// c + v + s, compute the rate of profit p' = s/C and the rate of surplus-value
// s' = s/v, record the mystification between them, and persist.
func (h *Handler) CreateProfitRate(w http.ResponseWriter, r *http.Request) {
	var req profitRateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Constant < 0 || req.Variable < 0 || req.SurplusValue < 0 {
		writeError(w, http.StatusUnprocessableEntity, "constant, variable and surplus_value must be non-negative")
		return
	}
	analysis := profit.ComputeProfitRateAnalysis(
		profit.ConstantCapital(req.Constant),
		profit.VariableCapital(req.Variable),
		profit.SurplusValue(req.SurplusValue),
	)
	saved, err := h.Store.CreateProfitRate(r.Context(), analysis)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/profit/rate/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetProfitRate handles GET /v1/profit/rate/{id}.
func (h *Handler) GetProfitRate(w http.ResponseWriter, r *http.Request) {
	id := profit.ProfitRateID(r.PathValue("id"))
	a, err := h.Store.GetProfitRate(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// ListProfitRates handles GET /v1/profit/rate.
func (h *Handler) ListProfitRates(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListProfitRates(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []profit.ProfitRateAnalysis{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
