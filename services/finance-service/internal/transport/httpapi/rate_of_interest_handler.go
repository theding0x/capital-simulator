package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// rateOfInterestRequest is the body for POST /v1/credit/rate-of-interest.
//
// The ceiling of the rate of interest is the *general* average rate of profit
// (Vol. III Part V), so average_profit_rate_bp must be that general rate
// (ΣS/ΣC from /v1/avgprofit/general-rate), not an individual firm's p′. Set
// general_rate_id to source it from a stored general-rate record instead, which
// enforces that provenance; when general_rate_id is set, average_profit_rate_bp
// is ignored.
type rateOfInterestRequest struct {
	RateBP              int64                       `json:"rate_bp"`
	AverageProfitRateBP int64                       `json:"average_profit_rate_bp"`
	GeneralRateID       string                      `json:"general_rate_id"`
	CyclePhase          credit.IndustrialCyclePhase `json:"cycle_phase"`
	Period              string                      `json:"period"`
}

// interestRateAnalysisRequest is the body for POST /v1/credit/interest-rate-analysis.
//
// As with rate-of-interest, average_profit_rate_bp must be the general average
// rate of profit; set general_rate_id to read it from a stored general-rate
// record instead (which then takes precedence over average_profit_rate_bp).
type interestRateAnalysisRequest struct {
	AverageProfitRateBP int64  `json:"average_profit_rate_bp"`
	GeneralRateID       string `json:"general_rate_id"`
	InterestRateBP      int64  `json:"interest_rate_bp"`
}

// resolveAverageProfitRate determines the average/total profit rate that bounds
// the Part V split between interest and profit of enterprise. Marx makes the
// *general* average rate of profit (ΣS/ΣC, round-half-up) the ceiling of the
// rate of interest — not an individual firm's p′. When generalRateID is set, the
// rate is read from the stored general-rate record, enforcing that provenance,
// and fallbackBP is ignored; an id that names no record yields store.ErrNotFound.
// When generalRateID is empty, fallbackBP is used (the request field is
// documented as requiring the general rate).
func (h *Handler) resolveAverageProfitRate(ctx context.Context, generalRateID string, fallbackBP int64) (int64, error) {
	if generalRateID == "" {
		return fallbackBP, nil
	}
	g, err := h.Store.GetGeneralProfitRate(ctx, avgprofit.GeneralProfitRateID(generalRateID))
	if err != nil {
		return 0, err
	}
	return int64(g.Rate), nil
}

// CreateRateOfInterest handles POST /v1/credit/rate-of-interest.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateRateOfInterest(w http.ResponseWriter, r *http.Request) {
	var req rateOfInterestRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	avgBP, err := h.resolveAverageProfitRate(r.Context(), req.GeneralRateID, req.AverageProfitRateBP)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	roi, err := credit.NewRateOfInterest(req.RateBP, avgBP, req.CyclePhase, req.Period)
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

	avgBP, err := h.resolveAverageProfitRate(r.Context(), req.GeneralRateID, req.AverageProfitRateBP)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	analysis, err := credit.NewInterestRateAnalysis(avgBP, req.InterestRateBP)
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
