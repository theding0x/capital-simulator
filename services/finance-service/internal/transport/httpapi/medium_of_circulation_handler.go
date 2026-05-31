package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// currencyObservationRequest is the body for POST /v1/credit/currency-observations.
type currencyObservationRequest struct {
	Name                     string `json:"name"`
	TotalCurrencyOutstanding int64  `json:"total_currency_outstanding"`
	CoinFunctionAmount       int64  `json:"coin_function_amount"`
	CapitalTransferAmount    int64  `json:"capital_transfer_amount"`
	ReserveAmount            int64  `json:"reserve_amount"`
	ReserveDescription       string `json:"reserve_description"`
	Description              string `json:"description"`
}

// CreateCurrencyObservation handles POST /v1/credit/currency-observations.
// Returns 400 on bad JSON, 422 on invariant violations, 201 + Location on success.
func (h *Handler) CreateCurrencyObservation(w http.ResponseWriter, r *http.Request) {
	var req currencyObservationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	reserve := credit.ReserveFund{
		Amount:      req.ReserveAmount,
		Description: req.ReserveDescription,
	}

	o, err := credit.NewCurrencyObservation(
		req.Name,
		req.TotalCurrencyOutstanding,
		req.CoinFunctionAmount,
		req.CapitalTransferAmount,
		reserve,
		req.Description,
	)
	if err != nil {
		if errors.Is(err, credit.ErrNonPositiveCurrency) ||
			errors.Is(err, credit.ErrFunctionSplitExceedsTotal) ||
			errors.Is(err, credit.ErrNegativeReserve) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateCurrencyObservation(r.Context(), o)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/currency-observations/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListCurrencyObservations handles GET /v1/credit/currency-observations.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListCurrencyObservations(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCurrencyObservations(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.CurrencyObservation{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetCurrencyObservation handles GET /v1/credit/currency-observations/{id}.
func (h *Handler) GetCurrencyObservation(w http.ResponseWriter, r *http.Request) {
	id := credit.CurrencyObservationID(r.PathValue("id"))
	o, err := h.Store.GetCurrencyObservation(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, o)
}
