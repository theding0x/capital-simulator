package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// profitDivisionRequest is the body for POST /v1/credit/profit-division.
//
// total_profit_bp is the average profit being split into interest and profit of
// enterprise; like the rate of interest, its proper source is the general
// average rate of profit. Set general_rate_id to read total_profit_bp from a
// stored general-rate record instead (which then takes precedence).
type profitDivisionRequest struct {
	TotalProfitBP      int64                       `json:"total_profit_bp"`
	GeneralRateID      string                      `json:"general_rate_id"`
	InterestBP         int64                       `json:"interest_bp"`
	OwnershipForm      credit.CapitalOwnershipForm `json:"ownership_form"`
	WagesLabourMinutes int64                       `json:"wages_labour_minutes"`
	MystificationIndex int64                       `json:"mystification_index"`
}

// CreateProfitDivision handles POST /v1/credit/profit-division.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateProfitDivision(w http.ResponseWriter, r *http.Request) {
	var req profitDivisionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	totalBP, err := h.resolveAverageProfitRate(r.Context(), req.GeneralRateID, req.TotalProfitBP)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	wages := credit.WagesOfSuperintendence{LabourMinutes: req.WagesLabourMinutes}
	pd, err := credit.NewProfitDivision(
		totalBP,
		req.InterestBP,
		req.OwnershipForm,
		wages,
		credit.MystificationIndex(req.MystificationIndex),
	)
	if err != nil {
		if errors.Is(err, credit.ErrNonPositiveTotalProfit) ||
			errors.Is(err, credit.ErrNegativeInterest) ||
			errors.Is(err, credit.ErrInterestNotLessThanTotal) ||
			errors.Is(err, credit.ErrInvalidOwnershipForm) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateProfitDivision(r.Context(), pd)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/profit-division/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListProfitDivisions handles GET /v1/credit/profit-division.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListProfitDivisions(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListProfitDivisions(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.ProfitDivision{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetProfitDivision handles GET /v1/credit/profit-division/{id}.
func (h *Handler) GetProfitDivision(w http.ResponseWriter, r *http.Request) {
	id := credit.ProfitDivisionID(r.PathValue("id"))
	pd, err := h.Store.GetProfitDivision(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, pd)
}
