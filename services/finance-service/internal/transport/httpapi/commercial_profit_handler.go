package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
)

// commercialProfitRequest is the body for POST /v1/merchant/commercial-profit.
// The caller supplies productive_capital and commercial_capital (pence, both > 0)
// and surplus_value (pence, >= 0).
type commercialProfitRequest struct {
	ProductiveCapital int64 `json:"productive_capital"`
	CommercialCapital int64 `json:"commercial_capital"`
	SurplusValue      int64 `json:"surplus_value"`
}

// CreateCommercialProfit handles POST /v1/merchant/commercial-profit. Returns
// 400 on bad JSON, 422 on ErrNonPositiveCapital or ErrNegativeSurplusValue,
// 201 + Location on success.
func (h *Handler) CreateCommercialProfit(w http.ResponseWriter, r *http.Request) {
	var req commercialProfitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	cp, err := merchant.NewCommercialProfit(
		req.ProductiveCapital,
		req.CommercialCapital,
		req.SurplusValue,
	)
	if err != nil {
		if errors.Is(err, merchant.ErrNonPositiveCapital) || errors.Is(err, merchant.ErrNegativeSurplusValue) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateCommercialProfit(r.Context(), cp)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/merchant/commercial-profit/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListCommercialProfits handles GET /v1/merchant/commercial-profit. Returns a
// JSON object with an "items" array — never null.
func (h *Handler) ListCommercialProfits(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCommercialProfits(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []merchant.CommercialProfit{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetCommercialProfit handles GET /v1/merchant/commercial-profit/{id}.
func (h *Handler) GetCommercialProfit(w http.ResponseWriter, r *http.Request) {
	id := merchant.CommercialProfitID(r.PathValue("id"))
	cp, err := h.Store.GetCommercialProfit(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cp)
}
