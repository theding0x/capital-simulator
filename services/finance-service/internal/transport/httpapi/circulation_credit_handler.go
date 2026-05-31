package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// clearingHouseSettlementRequest is the body for POST /v1/credit/clearing-house-settlements.
// net_balance is computed server-side and must not be supplied by the caller.
type clearingHouseSettlementRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TotalClaims int64  `json:"total_claims"`
	MoneyUsed   int64  `json:"money_used"`
}

// CreateClearingHouseSettlement handles POST /v1/credit/clearing-house-settlements.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateClearingHouseSettlement(w http.ResponseWriter, r *http.Request) {
	var req clearingHouseSettlementRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s, err := credit.NewClearingHouseSettlement(req.Name, req.Description, req.TotalClaims, req.MoneyUsed)
	if err != nil {
		if errors.Is(err, credit.ErrMoneyUsedExceedsClaims) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateClearingHouseSettlement(r.Context(), s)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/clearing-house-settlements/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListClearingHouseSettlements handles GET /v1/credit/clearing-house-settlements.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListClearingHouseSettlements(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListClearingHouseSettlements(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.ClearingHouseSettlement{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetClearingHouseSettlement handles GET /v1/credit/clearing-house-settlements/{id}.
func (h *Handler) GetClearingHouseSettlement(w http.ResponseWriter, r *http.Request) {
	id := credit.ClearingHouseSettlementID(r.PathValue("id"))
	s, err := h.Store.GetClearingHouseSettlement(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s)
}
