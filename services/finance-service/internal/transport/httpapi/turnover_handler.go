package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// turnoverRequest is the body for POST /v1/profit/turnover-analysis. Capital
// magnitudes are LabourMinutes: the total capital C and the variable capital
// advanced v (the v in C). SRate is the rate of surplus-value s' in basis points
// (100% = 10000 bp), and N is the number of turnovers of the variable capital
// per year. The annual rate of profit p' = s'·n·(v/C) is derived server-side.
type turnoverRequest struct {
	C     int64 `json:"c"`
	V     int64 `json:"v"`
	SRate int64 `json:"s_rate"`
	N     int64 `json:"n"`
}

// CreateTurnoverAnalysis handles POST /v1/profit/turnover-analysis — from a
// capital (total C, variable advanced v), a rate of surplus-value s' and a
// turnover count n, compute the annual rate of profit p' = s'·n·(v/C) together
// with the single-turnover rate, the annual rate of surplus-value S' = s'·n and
// the year's wage bill vn, then persist.
func (h *Handler) CreateTurnoverAnalysis(w http.ResponseWriter, r *http.Request) {
	var req turnoverRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.C < 0 || req.V < 0 || req.SRate < 0 || req.N < 0 {
		writeError(w, http.StatusUnprocessableEntity, "c, v, s_rate and n must be non-negative")
		return
	}
	if req.V > req.C {
		writeError(w, http.StatusUnprocessableEntity, "variable capital v cannot exceed total capital c")
		return
	}
	analysis := profit.ComputeTurnoverAnalysis(
		profit.TotalCapital(req.C),
		profit.VariableCapitalPerTurnover(req.V),
		profit.RateOfSurplusValue(req.SRate),
		profit.TurnoverCount(req.N),
	)
	saved, err := h.Store.CreateTurnoverAnalysis(r.Context(), analysis)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/profit/turnover-analysis/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetTurnoverAnalysis handles GET /v1/profit/turnover-analysis/{id}.
func (h *Handler) GetTurnoverAnalysis(w http.ResponseWriter, r *http.Request) {
	id := profit.TurnoverAnalysisID(r.PathValue("id"))
	a, err := h.Store.GetTurnoverAnalysis(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// ListTurnoverAnalyses handles GET /v1/profit/turnover-analysis.
func (h *Handler) ListTurnoverAnalyses(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListTurnoverAnalyses(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []profit.TurnoverAnalysis{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
