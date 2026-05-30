package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/tendency"
)

// trajectoryRequest is the body for POST /v1/tendency/trajectory. The caller
// supplies a label, a constant rate of surplus-value (s′ in basis points), and a
// series of (c, v) steps. The per-period surplus-value and rate of profit — and
// the falling-rate trajectory — are derived server-side.
type trajectoryRequest struct {
	Label            string     `json:"label"`
	SurplusValueRate int64      `json:"surplus_value_rate"` // s′ in basis points
	Steps            [][2]int64 `json:"steps"`              // each [c, v]
}

// rateMassRequest is the body for POST /v1/tendency/rate-mass. The caller
// supplies two (C, rate) states; the masses and their signed change are derived.
type rateMassRequest struct {
	OldC    int64 `json:"old_c"`
	OldRate int64 `json:"old_rate"` // basis points
	NewC    int64 `json:"new_c"`
	NewRate int64 `json:"new_rate"` // basis points
}

// CreateCompositionTrajectory handles POST /v1/tendency/trajectory.
// Returns 400 on bad JSON or empty steps, 422 on a negative magnitude,
// 201 + Location on success.
func (h *Handler) CreateCompositionTrajectory(w http.ResponseWriter, r *http.Request) {
	var req trajectoryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Steps) == 0 {
		writeError(w, http.StatusBadRequest, "steps must contain at least one [c, v] pair")
		return
	}
	if req.SurplusValueRate < 0 {
		writeError(w, http.StatusUnprocessableEntity, "surplus_value_rate must be non-negative")
		return
	}
	for _, step := range req.Steps {
		if step[0] < 0 || step[1] < 0 {
			writeError(w, http.StatusUnprocessableEntity, "constant and variable capital must be non-negative")
			return
		}
	}

	trajectory := tendency.BuildCompositionTrajectory(req.Label, req.SurplusValueRate, req.Steps)
	saved, err := h.Store.CreateCompositionTrajectory(r.Context(), trajectory)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/tendency/trajectory/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListCompositionTrajectories handles GET /v1/tendency/trajectory. It returns a
// JSON object with an "items" array — never null.
func (h *Handler) ListCompositionTrajectories(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCompositionTrajectories(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []tendency.CompositionTrajectory{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetCompositionTrajectory handles GET /v1/tendency/trajectory/{id}.
func (h *Handler) GetCompositionTrajectory(w http.ResponseWriter, r *http.Request) {
	id := tendency.CompositionTrajectoryID(r.PathValue("id"))
	t, err := h.Store.GetCompositionTrajectory(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// CreateRateMassContradiction handles POST /v1/tendency/rate-mass.
// Returns 400 on bad JSON, 422 on a negative magnitude, 201 + Location on success.
func (h *Handler) CreateRateMassContradiction(w http.ResponseWriter, r *http.Request) {
	var req rateMassRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.OldC < 0 || req.OldRate < 0 || req.NewC < 0 || req.NewRate < 0 {
		writeError(w, http.StatusUnprocessableEntity, "old_c, old_rate, new_c and new_rate must be non-negative")
		return
	}

	contradiction := tendency.ComputeRateMassContradiction(req.OldC, req.OldRate, req.NewC, req.NewRate)
	saved, err := h.Store.CreateRateMassContradiction(r.Context(), contradiction)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/tendency/rate-mass/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetRateMassContradiction handles GET /v1/tendency/rate-mass/{id}.
func (h *Handler) GetRateMassContradiction(w http.ResponseWriter, r *http.Request) {
	id := tendency.RateMassContradictionID(r.PathValue("id"))
	c, err := h.Store.GetRateMassContradiction(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}
