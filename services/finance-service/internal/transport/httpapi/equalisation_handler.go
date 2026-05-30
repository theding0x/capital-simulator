package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// capitalFlowRequest is the body for POST /v1/avgprofit/capital-flow.
// The handler builds a fuller CapitalFlow (with market price/value) then
// wraps it in an Equalisation computed from initialRate and targetRate.
type capitalFlowRequest struct {
	SphereName        string `json:"sphere_name"`
	CurrentSphereRate int64  `json:"current_sphere_rate"`
	GeneralRate       int64  `json:"general_rate"`
	MarketPrice       int64  `json:"market_price"`
	MarketValue       int64  `json:"market_value"`
}

// CreateCapitalFlow handles POST /v1/avgprofit/capital-flow.
// It computes a CapitalFlow (using all supplied inputs), builds an
// Equalisation via ComputeEqualisation, attaches the fuller Flow, persists
// the Equalisation, and returns 201 + Location.
func (h *Handler) CreateCapitalFlow(w http.ResponseWriter, r *http.Request) {
	var req capitalFlowRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Build the rich CapitalFlow first (carries market price/value).
	flow := avgprofit.ComputeCapitalFlow(
		req.SphereName,
		avgprofit.SphereProfitRate(req.CurrentSphereRate),
		avgprofit.SphereProfitRate(req.GeneralRate),
		avgprofit.MarketPriceAmount(req.MarketPrice),
		avgprofit.CommodityValue(req.MarketValue),
	)

	// Build Equalisation from rate inputs, then override Flow with the richer one.
	eq := avgprofit.ComputeEqualisation(
		req.SphereName,
		avgprofit.SphereProfitRate(req.CurrentSphereRate),
		avgprofit.SphereProfitRate(req.GeneralRate),
	)
	eq.Flow = flow
	// Keep the persisted Direction consistent with the richer Flow: when
	// general_rate is 0 the flow uses the market-price-vs-value path, which the
	// rate-only ComputeEqualisation would otherwise miss.
	eq.Direction = flow.Direction

	if err := eq.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	saved, err := h.Store.CreateEqualisation(r.Context(), eq)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/avgprofit/equalisation/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetEqualisation handles GET /v1/avgprofit/equalisation/{id}.
func (h *Handler) GetEqualisation(w http.ResponseWriter, r *http.Request) {
	id := avgprofit.EqualisationID(r.PathValue("id"))
	eq, err := h.Store.GetEqualisation(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, eq)
}
