package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// priceChangeRequest is the body for POST /v1/avgprofit/price-change. The caller
// supplies the sphere name, the cause of the change, and the old/new prices for
// context. The rate_changed/value_changed/price_changed flags are derived from
// the cause server-side.
type priceChangeRequest struct {
	SphereName string `json:"sphere_name"`
	Cause      string `json:"cause"`
	OldPrice   int64  `json:"old_price"`
	NewPrice   int64  `json:"new_price"`
}

// compensationGroundRequest is the body for POST /v1/avgprofit/compensation-ground.
// The caller supplies a claimed reason and what it appears to be; the response
// asserts what it actually is — a share in aggregate surplus-value, not value
// creation.
type compensationGroundRequest struct {
	Reason      string `json:"reason"`
	AppearsToBe string `json:"appears_to_be"`
}

// CreatePriceOfProductionChange handles POST /v1/avgprofit/price-change.
// Returns 400 on bad JSON, 422 on negative price, 201 + Location on success.
func (h *Handler) CreatePriceOfProductionChange(w http.ResponseWriter, r *http.Request) {
	var req priceChangeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.OldPrice < 0 || req.NewPrice < 0 {
		writeError(w, http.StatusUnprocessableEntity, "old_price and new_price must be non-negative")
		return
	}

	change := avgprofit.ComputePriceOfProductionChange(
		req.SphereName,
		avgprofit.PriceChangeCause(req.Cause),
		avgprofit.ProductionPrice(req.OldPrice),
		avgprofit.ProductionPrice(req.NewPrice),
	)
	saved, err := h.Store.CreatePriceOfProductionChange(r.Context(), change)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/avgprofit/price-change/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetPriceOfProductionChange handles GET /v1/avgprofit/price-change/{id}.
func (h *Handler) GetPriceOfProductionChange(w http.ResponseWriter, r *http.Request) {
	id := avgprofit.PriceOfProductionChangeID(r.PathValue("id"))
	c, err := h.Store.GetPriceOfProductionChange(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// ComputeCompensationGround handles POST /v1/avgprofit/compensation-ground. It
// returns the compensation-ground value object (not persisted): whatever the
// claimed reason, it is a share in aggregate surplus-value and creates no value.
// Returns 400 on bad JSON, 200 on success.
func (h *Handler) ComputeCompensationGround(w http.ResponseWriter, r *http.Request) {
	var req compensationGroundRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	ground := avgprofit.ComputeCompensationGround(req.Reason, req.AppearsToBe)
	writeJSON(w, http.StatusOK, ground)
}
