package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
)

// historicalMerchantCapitalRequest is the body for POST /v1/merchant/historical-forms.
type historicalMerchantCapitalRequest struct {
	Name               string `json:"name"`
	Stage              int    `json:"stage"`               // 1–3
	ProfitSource       string `json:"profit_source"`
	WageLabour         bool   `json:"wage_labour"`
	SubordinationIndex int64  `json:"subordination_index"` // bp [0, 10000]
}

// CreateHistoricalMerchantCapital handles POST /v1/merchant/historical-forms.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateHistoricalMerchantCapital(w http.ResponseWriter, r *http.Request) {
	var req historicalMerchantCapitalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	hm, err := merchant.NewHistoricalMerchantCapital(
		req.Name,
		merchant.DevelopmentStage(req.Stage),
		req.ProfitSource,
		req.WageLabour,
		req.SubordinationIndex,
	)
	if err != nil {
		if errors.Is(err, merchant.ErrInvalidDevelopmentStage) ||
			errors.Is(err, merchant.ErrSubordinationOutOfRange) ||
			errors.Is(err, merchant.ErrWageLabourInPreCapitalist) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateHistoricalMerchantCapital(r.Context(), hm)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/merchant/historical-forms/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListHistoricalMerchantCapitals handles GET /v1/merchant/historical-forms.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListHistoricalMerchantCapitals(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListHistoricalMerchantCapitals(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []merchant.HistoricalMerchantCapital{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetHistoricalMerchantCapital handles GET /v1/merchant/historical-forms/{id}.
func (h *Handler) GetHistoricalMerchantCapital(w http.ResponseWriter, r *http.Request) {
	id := merchant.HistoricalMerchantCapitalID(r.PathValue("id"))
	hm, err := h.Store.GetHistoricalMerchantCapital(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, hm)
}
