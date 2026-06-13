package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// generalRateRequest is the body for POST /v1/avgprofit/general-rate and
// POST /v1/avgprofit/social-aggregate. The caller supplies a slice of spheres;
// each sphere is computed server-side by ComputeProductionSphere.
type generalRateRequest struct {
	Spheres []sphereInput `json:"spheres"`
}

// sphereInput carries the user-supplied fields for one production sphere.
type sphereInput struct {
	Name  string `json:"name"`
	C     int64  `json:"c"`
	V     int64  `json:"v"`
	SRate int64  `json:"s_rate"`
}

// CreateGeneralProfitRate handles POST /v1/avgprofit/general-rate — compute
// the capital-weighted general rate over the supplied spheres and persist it.
// Returns 400 on bad JSON, 422 if any sphere is invalid, 201 + Location on success.
func (h *Handler) CreateGeneralProfitRate(w http.ResponseWriter, r *http.Request) {
	var req generalRateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	spheres := make([]avgprofit.ProductionSphere, 0, len(req.Spheres))
	for _, inp := range req.Spheres {
		if inp.C < 0 || inp.V < 0 || inp.SRate < 0 {
			writeError(w, http.StatusUnprocessableEntity, "c, v, and s_rate must be non-negative")
			return
		}
		if inp.V > inp.C {
			writeError(w, http.StatusUnprocessableEntity, "variable capital v cannot exceed total capital c")
			return
		}
		spheres = append(spheres, avgprofit.ComputeProductionSphere(
			inp.Name,
			avgprofit.TotalCapital(inp.C),
			avgprofit.VariableCapital(inp.V),
			avgprofit.RateOfSurplusValue(inp.SRate),
		))
	}

	g := avgprofit.ComputeGeneralProfitRate(spheres)
	saved, err := h.Store.CreateGeneralProfitRate(r.Context(), g)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/avgprofit/general-rate/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetGeneralProfitRate handles GET /v1/avgprofit/general-rate/{id}.
func (h *Handler) GetGeneralProfitRate(w http.ResponseWriter, r *http.Request) {
	id := avgprofit.GeneralProfitRateID(r.PathValue("id"))
	g, err := h.Store.GetGeneralProfitRate(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, g)
}

// ComputeSocialAggregate handles POST /v1/avgprofit/social-aggregate — pure
// compute only (no store). Returns 200 with GeneralRateAggregate.
func (h *Handler) ComputeSocialAggregate(w http.ResponseWriter, r *http.Request) {
	var req generalRateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	spheres := make([]avgprofit.ProductionSphere, 0, len(req.Spheres))
	for _, inp := range req.Spheres {
		if inp.C < 0 || inp.V < 0 || inp.SRate < 0 {
			writeError(w, http.StatusUnprocessableEntity, "c, v, and s_rate must be non-negative")
			return
		}
		if inp.V > inp.C {
			writeError(w, http.StatusUnprocessableEntity, "variable capital v cannot exceed total capital c")
			return
		}
		spheres = append(spheres, avgprofit.ComputeProductionSphere(
			inp.Name,
			avgprofit.TotalCapital(inp.C),
			avgprofit.VariableCapital(inp.V),
			avgprofit.RateOfSurplusValue(inp.SRate),
		))
	}

	agg := avgprofit.ComputeSocialCapitalAggregate(spheres)
	writeJSON(w, http.StatusOK, agg)
}
