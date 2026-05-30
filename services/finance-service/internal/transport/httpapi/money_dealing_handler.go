package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
)

// moneyDealingCapitalRequest is the body for POST /v1/merchant/money-dealing.
// operation_kinds is an array of int values in range 1–5.
type moneyDealingCapitalRequest struct {
	MoneyAdvanced  int64   `json:"money_advanced"`
	GeneralRateBP  int64   `json:"general_rate_bp"`
	OperationKinds []int   `json:"operation_kinds"`
}

// CreateMoneyDealingCapital handles POST /v1/merchant/money-dealing.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateMoneyDealingCapital(w http.ResponseWriter, r *http.Request) {
	var req moneyDealingCapitalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	kinds := make([]merchant.MoneyOperationKind, len(req.OperationKinds))
	for i, k := range req.OperationKinds {
		kinds[i] = merchant.MoneyOperationKind(k)
	}

	md, err := merchant.NewMoneyDealingCapital(req.MoneyAdvanced, req.GeneralRateBP, kinds)
	if err != nil {
		if errors.Is(err, merchant.ErrNonPositiveMoney) ||
			errors.Is(err, merchant.ErrNonPositiveCapital) ||
			errors.Is(err, merchant.ErrInvalidOperationKind) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateMoneyDealingCapital(r.Context(), md)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/merchant/money-dealing/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListMoneyDealingCapitals handles GET /v1/merchant/money-dealing.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListMoneyDealingCapitals(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListMoneyDealingCapitals(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []merchant.MoneyDealingCapital{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetMoneyDealingCapital handles GET /v1/merchant/money-dealing/{id}.
func (h *Handler) GetMoneyDealingCapital(w http.ResponseWriter, r *http.Request) {
	id := merchant.MoneyDealingCapitalID(r.PathValue("id"))
	md, err := h.Store.GetMoneyDealingCapital(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, md)
}
