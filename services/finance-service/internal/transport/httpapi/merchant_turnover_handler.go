package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
)

// merchantTurnoverRequest is the body for POST /v1/merchant/turnover.
type merchantTurnoverRequest struct {
	MoneyAdvanced    int64 `json:"money_advanced"`
	GeneralRateBP    int64 `json:"general_rate_bp"`
	TurnoverCount    int64 `json:"turnover_count"`
	UnitsPerTurnover int64 `json:"units_per_turnover"`
}

// CreateMerchantTurnover handles POST /v1/merchant/turnover. Returns 400 on
// bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateMerchantTurnover(w http.ResponseWriter, r *http.Request) {
	var req merchantTurnoverRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	mt, err := merchant.NewMerchantTurnover(
		req.MoneyAdvanced,
		req.GeneralRateBP,
		req.TurnoverCount,
		req.UnitsPerTurnover,
	)
	if err != nil {
		if errors.Is(err, merchant.ErrNonPositiveCapital) ||
			errors.Is(err, merchant.ErrInvalidTurnoverCount) ||
			errors.Is(err, merchant.ErrNoUnits) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateMerchantTurnover(r.Context(), mt)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/merchant/turnover/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListMerchantTurnovers handles GET /v1/merchant/turnover. Returns a JSON
// object with an "items" array — never null.
func (h *Handler) ListMerchantTurnovers(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListMerchantTurnovers(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []merchant.MerchantTurnover{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetMerchantTurnover handles GET /v1/merchant/turnover/{id}.
func (h *Handler) GetMerchantTurnover(w http.ResponseWriter, r *http.Request) {
	id := merchant.MerchantTurnoverID(r.PathValue("id"))
	mt, err := h.Store.GetMerchantTurnover(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mt)
}

// turnoverEffectRequest is the body for POST /v1/merchant/turnover-effect.
type turnoverEffectRequest struct {
	MoneyAdvanced    int64   `json:"money_advanced"`
	GeneralRateBP    int64   `json:"general_rate_bp"`
	UnitsPerTurnover int64   `json:"units_per_turnover"`
	TurnoverCounts   []int64 `json:"turnover_counts"`
}

// ComputeTurnoverEffect handles POST /v1/merchant/turnover-effect. This is a
// pure compute endpoint — the result is NOT persisted. Returns 400 on bad
// JSON, 422 on validation errors, 200 with TurnoverEffectAnalysis on success.
func (h *Handler) ComputeTurnoverEffect(w http.ResponseWriter, r *http.Request) {
	var req turnoverEffectRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	analysis, err := merchant.NewTurnoverEffectAnalysis(
		req.MoneyAdvanced,
		req.GeneralRateBP,
		req.UnitsPerTurnover,
		req.TurnoverCounts,
	)
	if err != nil {
		if errors.Is(err, merchant.ErrNonPositiveCapital) ||
			errors.Is(err, merchant.ErrInvalidTurnoverCount) ||
			errors.Is(err, merchant.ErrNoUnits) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, analysis)
}
