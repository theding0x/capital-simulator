package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

// wageEffectSphereInput carries user-supplied fields for one production sphere.
type wageEffectSphereInput struct {
	Name string `json:"name"`
	C    int64  `json:"c"`
	V    int64  `json:"v"`
}

// wageEffectRequest is the body for POST /v1/avgprofit/wage-effect.
type wageEffectRequest struct {
	BaseConstant int64                    `json:"base_constant"`
	BaseVariable int64                    `json:"base_variable"`
	SRate        int64                    `json:"s_rate"`
	WageFactor   int64                    `json:"wage_factor"`
	Spheres      []wageEffectSphereInput  `json:"spheres"`
}

// CreateWageEffectAnalysis handles POST /v1/avgprofit/wage-effect.
// Returns 400 on bad JSON, 422 on negative magnitude, 201+Location on success.
func (h *Handler) CreateWageEffectAnalysis(w http.ResponseWriter, r *http.Request) {
	var req wageEffectRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.BaseConstant < 0 || req.BaseVariable < 0 || req.SRate < 0 || req.WageFactor < 0 {
		writeError(w, http.StatusUnprocessableEntity, avgprofit.ErrNegativeMagnitude.Error())
		return
	}

	spheres := make([]avgprofit.SphereInput, 0, len(req.Spheres))
	for _, s := range req.Spheres {
		if s.C < 0 || s.V < 0 {
			writeError(w, http.StatusUnprocessableEntity, avgprofit.ErrNegativeMagnitude.Error())
			return
		}
		spheres = append(spheres, avgprofit.SphereInput{
			Name: s.Name,
			C:    avgprofit.ConstantCapital(s.C),
			V:    avgprofit.VariableCapital(s.V),
		})
	}

	analysis := avgprofit.ComputeWageEffectAnalysis(
		avgprofit.ConstantCapital(req.BaseConstant),
		avgprofit.VariableCapital(req.BaseVariable),
		avgprofit.RateOfSurplusValue(req.SRate),
		req.WageFactor,
		spheres,
	)

	saved, err := h.Store.CreateWageEffectAnalysis(r.Context(), analysis)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/avgprofit/wage-effect/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListWageEffectAnalyses handles GET /v1/avgprofit/wage-effect.
// Returns {"items":[...]} with a non-nil slice.
func (h *Handler) ListWageEffectAnalyses(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListWageEffectAnalyses(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []avgprofit.WageEffectAnalysis{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetWageEffectAnalysis handles GET /v1/avgprofit/wage-effect/{id}.
// Returns 404 when the ID is not found.
func (h *Handler) GetWageEffectAnalysis(w http.ResponseWriter, r *http.Request) {
	id := avgprofit.WageEffectAnalysisID(r.PathValue("id"))
	a, err := h.Store.GetWageEffectAnalysis(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}
