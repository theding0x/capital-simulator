package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/tendency"
)

// counteractingForceRequest is the body for POST /v1/tendency/counteracting-force.
// The caller supplies one of the six valid kinds and an optional note; the
// ExcludedFromGeneralRate flag is derived server-side.
type counteractingForceRequest struct {
	Kind string `json:"kind"`
	Note string `json:"note"`
}

// counteractingScenarioRequest is the body for POST /v1/tendency/scenario. The
// caller supplies a label, a set of force kinds, the (c, v) steps of a base
// trajectory, and a constant s′ (basis points). The modified rate-of-profit path
// — and the proof that the law persists — are derived server-side.
type counteractingScenarioRequest struct {
	Label            string     `json:"label"`
	ForceKinds       []string   `json:"force_kinds"`
	Steps            [][2]int64 `json:"steps"`
	SurplusValueRate int64      `json:"surplus_value_rate"` // s′ in basis points
}

// CreateCounteractingForce handles POST /v1/tendency/counteracting-force.
// Returns 400 on bad JSON, 422 on an invalid kind, 201 + Location on success.
func (h *Handler) CreateCounteractingForce(w http.ResponseWriter, r *http.Request) {
	var req counteractingForceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	force, err := tendency.NewCounteractingForce(tendency.CounteractingForceKind(req.Kind), req.Note)
	if err != nil {
		if errors.Is(err, tendency.ErrInvalidForceKind) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateCounteractingForce(r.Context(), force)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/tendency/counteracting-force/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetCounteractingForce handles GET /v1/tendency/counteracting-force/{id}.
func (h *Handler) GetCounteractingForce(w http.ResponseWriter, r *http.Request) {
	id := tendency.CounteractingForceID(r.PathValue("id"))
	f, err := h.Store.GetCounteractingForce(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, f)
}

// CreateCounteractingScenario handles POST /v1/tendency/scenario. It validates
// the force kinds, builds the base trajectory, applies the forces, and persists
// the scenario. Returns 400 on bad JSON or empty force_kinds/steps, 422 on an
// invalid kind or when the law is not preserved, 201 + Location on success.
func (h *Handler) CreateCounteractingScenario(w http.ResponseWriter, r *http.Request) {
	var req counteractingScenarioRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.ForceKinds) == 0 {
		writeError(w, http.StatusBadRequest, "force_kinds must contain at least one kind")
		return
	}
	if len(req.Steps) == 0 {
		writeError(w, http.StatusBadRequest, "steps must contain at least one [c, v] pair")
		return
	}

	forces := make([]tendency.CounteractingForce, 0, len(req.ForceKinds))
	for _, k := range req.ForceKinds {
		force, err := tendency.NewCounteractingForce(tendency.CounteractingForceKind(k), "")
		if err != nil {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		forces = append(forces, force)
	}

	trajectory := tendency.BuildCompositionTrajectory(req.Label, req.SurplusValueRate, req.Steps)
	scenario, err := tendency.BuildCounteractingScenario(req.Label, trajectory, forces)
	if err != nil {
		switch {
		case errors.Is(err, tendency.ErrEmptyForces), errors.Is(err, tendency.ErrEmptyTrajectory):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		}
		return
	}

	saved, err := h.Store.CreateCounteractingScenario(r.Context(), scenario)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/tendency/scenario/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListCounteractingScenarios handles GET /v1/tendency/scenario. It returns a
// JSON object with an "items" array — never null.
func (h *Handler) ListCounteractingScenarios(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCounteractingScenarios(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []tendency.CounteractingScenario{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetCounteractingScenario handles GET /v1/tendency/scenario/{id}.
func (h *Handler) GetCounteractingScenario(w http.ResponseWriter, r *http.Request) {
	id := tendency.CounteractingScenarioID(r.PathValue("id"))
	s, err := h.Store.GetCounteractingScenario(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s)
}
