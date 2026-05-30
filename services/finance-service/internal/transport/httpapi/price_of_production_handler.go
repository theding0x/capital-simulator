package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// priceOfProductionRequest is the body for POST /v1/avgprofit/price-of-production.
// The caller supplies sphere name, cost-price k, general rate p̄′, and
// commodity value (0 = unknown; composition will then be "").
type priceOfProductionRequest struct {
	SphereName     string `json:"sphere_name"`
	CostPrice      int64  `json:"cost_price"`
	GeneralRate    int64  `json:"general_rate"`
	CommodityValue int64  `json:"commodity_value"`
}

// CreatePriceOfProduction handles POST /v1/avgprofit/price-of-production.
// Returns 400 on bad JSON, 422 on negative magnitude, 201 + Location on success.
func (h *Handler) CreatePriceOfProduction(w http.ResponseWriter, r *http.Request) {
	var req priceOfProductionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.CostPrice < 0 || req.GeneralRate < 0 || req.CommodityValue < 0 {
		writeError(w, http.StatusUnprocessableEntity, "cost_price, general_rate, and commodity_value must be non-negative")
		return
	}

	pop := avgprofit.ComputePriceOfProduction(
		req.SphereName,
		avgprofit.CostPrice(req.CostPrice),
		avgprofit.SphereProfitRate(req.GeneralRate),
		avgprofit.CommodityValue(req.CommodityValue),
	)
	saved, err := h.Store.CreatePriceOfProduction(r.Context(), pop)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/avgprofit/price-of-production/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetPriceOfProduction handles GET /v1/avgprofit/price-of-production/{id}.
func (h *Handler) GetPriceOfProduction(w http.ResponseWriter, r *http.Request) {
	id := avgprofit.PriceOfProductionID(r.PathValue("id"))
	p, err := h.Store.GetPriceOfProduction(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// ListPricesOfProduction handles GET /v1/avgprofit/price-of-production —
// returns all stored records newest first. The items slice is normalised to
// non-nil so the JSON wire format always carries an array, never null.
func (h *Handler) ListPricesOfProduction(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListPricesOfProduction(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []avgprofit.PriceOfProduction{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
