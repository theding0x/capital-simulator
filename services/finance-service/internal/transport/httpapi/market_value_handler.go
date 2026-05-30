package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// marketValueRequest is the body for POST /v1/avgprofit/market-value.
type marketValueRequest struct {
	SphereName          string `json:"sphere_name"`
	BulkConditionValue  int64  `json:"bulk_condition_value"`
	BestConditionValue  int64  `json:"best_condition_value"`
	WorstConditionValue int64  `json:"worst_condition_value"`
}

// CreateMarketValue handles POST /v1/avgprofit/market-value.
// Returns 400 on bad JSON, 422 on negative magnitude, 201 + Location on success.
func (h *Handler) CreateMarketValue(w http.ResponseWriter, r *http.Request) {
	var req marketValueRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	mv := avgprofit.ComputeMarketValue(
		req.SphereName,
		avgprofit.CommodityValue(req.BulkConditionValue),
		avgprofit.CommodityValue(req.BestConditionValue),
		avgprofit.CommodityValue(req.WorstConditionValue),
	)
	if err := mv.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	saved, err := h.Store.CreateMarketValue(r.Context(), mv)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/avgprofit/market-value/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetMarketValue handles GET /v1/avgprofit/market-value/{id}.
func (h *Handler) GetMarketValue(w http.ResponseWriter, r *http.Request) {
	id := avgprofit.MarketValueID(r.PathValue("id"))
	mv, err := h.Store.GetMarketValue(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mv)
}

// surplusProfitRequest is the body for POST /v1/avgprofit/surplus-profit.
type surplusProfitRequest struct {
	FirmName        string `json:"firm_name"`
	IndividualValue int64  `json:"individual_value"`
	MarketValue     int64  `json:"market_value"`
	OutputQty       int64  `json:"output_qty"`
	GeneralRate     int64  `json:"general_rate"`
}

// CreateSurplusProfit handles POST /v1/avgprofit/surplus-profit.
// Returns 400 on bad JSON, 422 on negative magnitude, 201 + Location on success.
func (h *Handler) CreateSurplusProfit(w http.ResponseWriter, r *http.Request) {
	var req surplusProfitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	sp := avgprofit.ComputeSurplusProfit(
		req.FirmName,
		avgprofit.IndividualValue(req.IndividualValue),
		avgprofit.CommodityValue(req.MarketValue),
		req.OutputQty,
		avgprofit.SphereProfitRate(req.GeneralRate),
	)
	if err := sp.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	saved, err := h.Store.CreateSurplusProfit(r.Context(), sp)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/avgprofit/surplus-profit/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetSurplusProfit handles GET /v1/avgprofit/surplus-profit/{id}.
func (h *Handler) GetSurplusProfit(w http.ResponseWriter, r *http.Request) {
	id := avgprofit.SurplusProfitID(r.PathValue("id"))
	sp, err := h.Store.GetSurplusProfit(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sp)
}
