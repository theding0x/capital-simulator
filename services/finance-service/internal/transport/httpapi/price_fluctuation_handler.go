package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// priceFluctuationRequest is the body for POST /v1/profit/price-fluctuation. The
// constant capital is split into a Fixed component the price movement does not
// touch and an OriginalMaterial component — the circulating outlay on raw
// material — that scales with the price. PriceFactor is that price as a
// percentage of its original level (100 = unchanged, 200 = doubled, 50 = halved).
// Variable and SurplusValue are held fixed. The direction of the movement and the
// old and new rates of profit are derived server-side.
type priceFluctuationRequest struct {
	Fixed            int64 `json:"fixed"`
	OriginalMaterial int64 `json:"original_material"`
	PriceFactor      int64 `json:"price_factor"`
	Variable         int64 `json:"variable"`
	SurplusValue     int64 `json:"surplus_value"`
}

// CreatePriceFluctuationAnalysis handles POST /v1/profit/price-fluctuation — from
// a re-pricing of the raw-material element of constant capital, compute the
// movement it produces in the rate of profit s/(c+v) and persist.
func (h *Handler) CreatePriceFluctuationAnalysis(w http.ResponseWriter, r *http.Request) {
	var req priceFluctuationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Fixed < 0 || req.OriginalMaterial < 0 || req.Variable < 0 || req.SurplusValue < 0 {
		writeError(w, http.StatusUnprocessableEntity, "fixed, original_material, variable and surplus_value must be non-negative")
		return
	}
	if req.PriceFactor <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "price_factor must be positive (a percentage of the original price: 100 = unchanged, 200 = doubled, 50 = halved)")
		return
	}
	analysis := profit.ComputePriceFluctuationAnalysis(profit.ConstantCapitalPriceEffect{
		FixedCapital:            profit.ConstantCapital(req.Fixed),
		OriginalMaterialCapital: profit.ConstantCapital(req.OriginalMaterial),
		PriceFactor:             req.PriceFactor,
		VariableCapital:         profit.VariableCapital(req.Variable),
		SurplusValue:            profit.SurplusValue(req.SurplusValue),
	})
	saved, err := h.Store.CreatePriceFluctuationAnalysis(r.Context(), analysis)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/profit/price-fluctuation/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetPriceFluctuationAnalysis handles GET /v1/profit/price-fluctuation/{id}.
func (h *Handler) GetPriceFluctuationAnalysis(w http.ResponseWriter, r *http.Request) {
	id := profit.PriceFluctuationAnalysisID(r.PathValue("id"))
	a, err := h.Store.GetPriceFluctuationAnalysis(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// ListPriceFluctuationAnalyses handles GET /v1/profit/price-fluctuation.
func (h *Handler) ListPriceFluctuationAnalyses(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListPriceFluctuationAnalyses(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []profit.PriceFluctuationAnalysis{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
