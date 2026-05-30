package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
)

// commercialCapitalRequest is the body for POST /v1/merchant/commercial-capital.
// The caller supplies the money advanced (pence), a commodity description, and
// the merchant's function (0=realisation, 1=storage, 2=transport, 3=distribution).
type commercialCapitalRequest struct {
	MoneyAdvanced        int64  `json:"money_advanced"`
	CommodityDescription string `json:"commodity_description"`
	Function             int    `json:"function"`
}

// CreateCommercialCapital handles POST /v1/merchant/commercial-capital. Returns
// 400 on bad JSON, 422 on ErrNonPositiveMoney, 201 + Location on success.
func (h *Handler) CreateCommercialCapital(w http.ResponseWriter, r *http.Request) {
	var req commercialCapitalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	cc, err := merchant.NewCommercialCapital(
		req.MoneyAdvanced,
		req.CommodityDescription,
		merchant.CommercialCapitalFunction(req.Function),
	)
	if err != nil {
		if errors.Is(err, merchant.ErrNonPositiveMoney) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateCommercialCapital(r.Context(), cc)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/merchant/commercial-capital/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListCommercialCapitals handles GET /v1/merchant/commercial-capital. It returns
// a JSON object with an "items" array — never null.
func (h *Handler) ListCommercialCapitals(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCommercialCapitals(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []merchant.CommercialCapital{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetCommercialCapital handles GET /v1/merchant/commercial-capital/{id}.
func (h *Handler) GetCommercialCapital(w http.ResponseWriter, r *http.Request) {
	id := merchant.CommercialCapitalID(r.PathValue("id"))
	cc, err := h.Store.GetCommercialCapital(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cc)
}
