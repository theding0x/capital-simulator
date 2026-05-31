package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// rateOfInterestRequest is the body for POST /v1/credit/rate-of-interest.
type rateOfInterestRequest struct {
	RateBP              int64                       `json:"rate_bp"`
	AverageProfitRateBP int64                       `json:"average_profit_rate_bp"`
	CyclePhase          credit.IndustrialCyclePhase `json:"cycle_phase"`
	Period              string                      `json:"period"`
}

// interestRateAnalysisRequest is the body for POST /v1/credit/interest-rate-analysis.
type interestRateAnalysisRequest struct {
	AverageProfitRateBP int64 `json:"average_profit_rate_bp"`
	InterestRateBP      int64 `json:"interest_rate_bp"`
}

// CreateRateOfInterest handles POST /v1/credit/rate-of-interest.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateRateOfInterest(w http.ResponseWriter, r *http.Request) {
	var req rateOfInterestRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	roi, err := credit.NewRateOfInterest(req.RateBP, req.AverageProfitRateBP, req.CyclePhase, req.Period)
	if err != nil {
		if errors.Is(err, credit.ErrNegativeRate) ||
			errors.Is(err, credit.ErrNonPositiveAverageProfit) ||
			errors.Is(err, credit.ErrRateExceedsAverageProfit) ||
			errors.Is(err, credit.ErrInvalidCyclePhase) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateRateOfInterest(r.Context(), roi)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/rate-of-interest/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListRatesOfInterest handles GET /v1/credit/rate-of-interest.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListRatesOfInterest(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListRatesOfInterest(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.RateOfInterest{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetRateOfInterest handles GET /v1/credit/rate-of-interest/{id}.
func (h *Handler) GetRateOfInterest(w http.ResponseWriter, r *http.Request) {
	id := credit.RateOfInterestID(r.PathValue("id"))
	roi, err := h.Store.GetRateOfInterest(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, roi)
}

// ComputeInterestRateAnalysis handles POST /v1/credit/interest-rate-analysis.
// Stateless: computes profit of enterprise without persisting anything.
// Returns 400 on bad JSON, 422 on validation errors, 200 on success.
func (h *Handler) ComputeInterestRateAnalysis(w http.ResponseWriter, r *http.Request) {
	var req interestRateAnalysisRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	analysis, err := credit.NewInterestRateAnalysis(req.AverageProfitRateBP, req.InterestRateBP)
	if err != nil {
		if errors.Is(err, credit.ErrNegativeRate) ||
			errors.Is(err, credit.ErrNonPositiveAverageProfit) ||
			errors.Is(err, credit.ErrRateExceedsAverageProfit) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, analysis)
}
