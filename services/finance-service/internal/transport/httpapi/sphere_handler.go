package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// productionSphereRequest is the body for POST /v1/avgprofit/spheres. The user
// supplies the branch name, total capital C, variable capital V, and the rate of
// surplus-value s′ in basis points. All derived fields — constant capital,
// surplus-value, individual rate of profit, variable percent, labour-power index
// — are computed server-side by ComputeProductionSphere.
type productionSphereRequest struct {
	Name  string `json:"name"`
	C     int64  `json:"c"`
	V     int64  `json:"v"`
	SRate int64  `json:"s_rate"`
}

// CreateProductionSphere handles POST /v1/avgprofit/spheres — compute all
// derived fields and persist the sphere. Returns 400 on bad JSON, 422 on invalid
// magnitudes, 201 with a Location header on success.
func (h *Handler) CreateProductionSphere(w http.ResponseWriter, r *http.Request) {
	var req productionSphereRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.C < 0 || req.V < 0 || req.SRate < 0 {
		writeError(w, http.StatusUnprocessableEntity, "c, v, and s_rate must be non-negative")
		return
	}
	if req.V > req.C {
		writeError(w, http.StatusUnprocessableEntity, "variable capital v cannot exceed total capital c")
		return
	}
	sphere := avgprofit.ComputeProductionSphere(
		req.Name,
		avgprofit.TotalCapital(req.C),
		avgprofit.VariableCapital(req.V),
		avgprofit.RateOfSurplusValue(req.SRate),
	)
	saved, err := h.Store.CreateProductionSphere(r.Context(), sphere)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/avgprofit/spheres/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetProductionSphere handles GET /v1/avgprofit/spheres/{id}.
func (h *Handler) GetProductionSphere(w http.ResponseWriter, r *http.Request) {
	id := avgprofit.ProductionSphereID(r.PathValue("id"))
	s, err := h.Store.GetProductionSphere(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s)
}

// ListProductionSpheres handles GET /v1/avgprofit/spheres — returns all stored
// spheres newest first. The items slice is normalised to non-nil so the JSON
// wire format always carries an array, never null.
func (h *Handler) ListProductionSpheres(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListProductionSpheres(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []avgprofit.ProductionSphere{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
