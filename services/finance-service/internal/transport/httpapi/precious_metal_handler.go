package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// createGoldReserveRequest is the body for POST /v1/credit/gold-reserves.
type createGoldReserveRequest struct {
	Name        string `json:"name"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

// createRateOfExchangeRequest is the body for POST /v1/credit/rates-of-exchange.
type createRateOfExchangeRequest struct {
	Name        string `json:"name"`
	DeviationBP int64  `json:"deviation_bp"`
	Description string `json:"description"`
}

// CreateGoldReserve handles POST /v1/credit/gold-reserves.
// Returns 400 on bad JSON, 422 on negative amount, 201 + Location on success.
func (h *Handler) CreateGoldReserve(w http.ResponseWriter, r *http.Request) {
	var req createGoldReserveRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	g, err := credit.NewGoldReserve(req.Name, req.Amount, req.Description)
	if err != nil {
		if errors.Is(err, credit.ErrNegativeGoldReserveAmount) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateGoldReserve(r.Context(), g)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/gold-reserves/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListGoldReserves handles GET /v1/credit/gold-reserves.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListGoldReserves(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListGoldReserves(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.GoldReserve{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetGoldReserve handles GET /v1/credit/gold-reserves/{id}.
func (h *Handler) GetGoldReserve(w http.ResponseWriter, r *http.Request) {
	id := credit.GoldReserveID(r.PathValue("id"))
	g, err := h.Store.GetGoldReserve(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, g)
}

// CreateRateOfExchange handles POST /v1/credit/rates-of-exchange.
// Returns 400 on bad JSON, 201 + Location on success.
func (h *Handler) CreateRateOfExchange(w http.ResponseWriter, r *http.Request) {
	var req createRateOfExchangeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	re := credit.NewRateOfExchange(req.Name, req.DeviationBP, req.Description)

	saved, err := h.Store.CreateRateOfExchange(r.Context(), re)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/rates-of-exchange/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListRatesOfExchange handles GET /v1/credit/rates-of-exchange.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListRatesOfExchange(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListRatesOfExchange(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.RateOfExchange{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetRateOfExchange handles GET /v1/credit/rates-of-exchange/{id}.
func (h *Handler) GetRateOfExchange(w http.ResponseWriter, r *http.Request) {
	id := credit.RateOfExchangeID(r.PathValue("id"))
	re, err := h.Store.GetRateOfExchange(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, re)
}
