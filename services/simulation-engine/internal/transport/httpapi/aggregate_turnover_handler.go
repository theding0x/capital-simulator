package httpapi

import (
	"errors"
	"net/http"

	comp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	tv "github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

// --- request / response types ------------------------------------------------

type createAggregateTurnoverRequest struct {
	IndustrialCapitalID string             `json:"industrial_capital_id"`
	Contributions       []contributionInput `json:"contributions"`
}

type contributionInput struct {
	CapitalComponentID        string `json:"capital_component_id"`
	Kind                      string `json:"kind"`
	AdvancedPence             int64  `json:"advanced_pence"`
	TurnoverNumberBasisPoints int64  `json:"turnover_number_basis_points"`
	DifferenceKind            string `json:"difference_kind"`
}

type lifetimeCycleRequest struct {
	DurationMinutes int64 `json:"duration_minutes"`
}

type crisisPhaseRequest struct {
	Phase string `json:"phase"`
	Notes string `json:"notes,omitempty"`
}

// --- handlers ----------------------------------------------------------------

// CreateAggregateTurnover handles POST /v1/aggregate-turnovers.
// Request must include industrial_capital_id and at least one contribution.
// AnnualReturnPence is derived server-side: AdvancedPence × TurnoverNumberBasisPoints / 10000.
// Lens is forced to LensMoney (money-form homogeneity invariant, §2).
func (h *Handler) CreateAggregateTurnover(w http.ResponseWriter, r *http.Request) {
	var req createAggregateTurnoverRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.IndustrialCapitalID == "" {
		writeError(w, http.StatusBadRequest, "industrial_capital_id is required")
		return
	}
	if len(req.Contributions) == 0 {
		writeError(w, http.StatusBadRequest, "at least one contribution is required")
		return
	}

	at := tv.AggregateTurnover{
		ID:                  tv.NewAggregateTurnoverID(),
		IndustrialCapitalID: req.IndustrialCapitalID,
		Lens:                tv.LensMoney,
		AdvancedTotalPence:  1, // placeholder; recomputed after contributions are recorded
	}
	created, err := h.AggregateTurnovers.CreateAggregateTurnover(r.Context(), at)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}

	for _, ci := range req.Contributions {
		diffKind := tv.TurnoverDifferenceKind(ci.DifferenceKind)
		if diffKind == "" {
			diffKind = tv.DifferenceReal
		}
		c := tv.ComponentTurnoverContribution{
			CapitalComponentID:        comp.CapitalComponentID(ci.CapitalComponentID),
			Kind:                      comp.CapitalKind(ci.Kind),
			AdvancedPence:             tv.Pence(ci.AdvancedPence),
			TurnoverNumberBasisPoints: ci.TurnoverNumberBasisPoints,
			DifferenceKind:            diffKind,
		}
		if _, err := h.AggregateTurnovers.RecordContribution(r.Context(), created.ID, c); err != nil {
			h.writeServerError(w, err)
			return
		}
	}

	result, err := h.AggregateTurnovers.RecomputeAggregate(r.Context(), created.ID)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/aggregate-turnovers/"+string(result.ID))
	writeJSON(w, http.StatusCreated, result)
}

// GetAggregateTurnover handles GET /v1/aggregate-turnovers/{id}.
func (h *Handler) GetAggregateTurnover(w http.ResponseWriter, r *http.Request) {
	id := tv.AggregateTurnoverID(r.PathValue("id"))
	at, err := h.AggregateTurnovers.GetAggregateTurnover(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "aggregate turnover not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, at)
}

// ListAggregateTurnovers handles GET /v1/aggregate-turnovers.
// Optional query param: ?industrial_capital_id=
func (h *Handler) ListAggregateTurnovers(w http.ResponseWriter, r *http.Request) {
	icID := r.URL.Query().Get("industrial_capital_id")
	list, err := h.AggregateTurnovers.ListAggregateTurnovers(r.Context(), icID)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// RecomputeAggregateTurnover handles POST /v1/aggregate-turnovers/{id}/recompute.
func (h *Handler) RecomputeAggregateTurnover(w http.ResponseWriter, r *http.Request) {
	id := tv.AggregateTurnoverID(r.PathValue("id"))
	result, err := h.AggregateTurnovers.RecomputeAggregate(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "aggregate turnover not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// RecordAggregateTurnoverLifetimeCycle handles POST /v1/aggregate-turnovers/{id}/lifetime-cycle.
func (h *Handler) RecordAggregateTurnoverLifetimeCycle(w http.ResponseWriter, r *http.Request) {
	id := tv.AggregateTurnoverID(r.PathValue("id"))
	var req lifetimeCycleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.DurationMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "duration_minutes must be positive")
		return
	}
	lc := tv.LifetimeCycle{DurationMinutes: req.DurationMinutes}
	result, err := h.AggregateTurnovers.RecordLifetimeCycle(r.Context(), id, lc)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "aggregate turnover not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

// RecordCrisisPhase handles POST /v1/aggregate-turnovers/{id}/crisis-phase.
// The {id} is the aggregate_turnover_id; the server looks up the lifetime cycle.
func (h *Handler) RecordCrisisPhase(w http.ResponseWriter, r *http.Request) {
	atID := tv.AggregateTurnoverID(r.PathValue("id"))
	var req crisisPhaseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Phase == "" {
		writeError(w, http.StatusBadRequest, "phase is required")
		return
	}
	phase := tv.CrisisCyclePhase(req.Phase)
	switch phase {
	case tv.PhaseDepression, tv.PhaseMediumActivity, tv.PhasePrecipitancy, tv.PhaseCrisis:
	default:
		writeError(w, http.StatusBadRequest, "phase must be one of: depression, medium_activity, precipitancy, crisis")
		return
	}

	at, err := h.AggregateTurnovers.GetAggregateTurnover(r.Context(), atID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "aggregate turnover not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	if at.LifetimeCycle == nil {
		writeError(w, http.StatusBadRequest, "lifetime cycle must be registered before recording crisis phases")
		return
	}

	record := tv.CrisisCyclePhaseRecord{Phase: phase, Notes: req.Notes}
	result, err := h.AggregateTurnovers.RecordCrisisPhase(r.Context(), at.LifetimeCycle.ID, record)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "lifetime cycle not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}
