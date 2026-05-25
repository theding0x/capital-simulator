// Vol. II Ch. 7 — HTTP handlers for turnover time and number of turnovers.
package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	tv "github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

type createTurnoverRequest struct {
	IndustrialCapitalID string         `json:"industrial_capital_id"`
	Lens                tv.CircuitLens `json:"lens"`
	TurnoverTimeMinutes int64          `json:"turnover_time_minutes"`
}

type recordCycleRequest struct {
	StartedAt          time.Time `json:"started_at"`
	EndedAt            time.Time `json:"ended_at"`
	AdvancePence       tv.Pence  `json:"advance_pence"`
	ReturnedPence      tv.Pence  `json:"returned_pence"`
	ProductionMinutes  int64     `json:"production_minutes"`
	CirculationMinutes int64     `json:"circulation_minutes"`
}

// CreateTurnover handles POST /v1/turnovers.
func (h *Handler) CreateTurnover(w http.ResponseWriter, r *http.Request) {
	var req createTurnoverRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	t := tv.Turnover{
		IndustrialCapitalID: req.IndustrialCapitalID,
		Lens:                req.Lens,
		TurnoverTimeMinutes: req.TurnoverTimeMinutes,
	}
	created, err := h.Turnovers.CreateTurnover(r.Context(), t)
	if err != nil {
		if errors.Is(err, tv.ErrLensNotApplicableForTurnover) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Location", "/v1/turnovers/"+string(created.ID))
	writeJSON(w, http.StatusCreated, created)
}

// GetTurnover handles GET /v1/turnovers/{id}.
func (h *Handler) GetTurnover(w http.ResponseWriter, r *http.Request) {
	id := tv.TurnoverID(r.PathValue("id"))
	t, err := h.Turnovers.GetTurnover(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// ListTurnovers handles GET /v1/turnovers.
func (h *Handler) ListTurnovers(w http.ResponseWriter, r *http.Request) {
	icID := r.URL.Query().Get("industrial_capital_id")
	lens := tv.CircuitLens(r.URL.Query().Get("lens"))
	items, err := h.Turnovers.ListTurnovers(r.Context(), icID, lens)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// RecordCycle handles POST /v1/turnovers/{id}/cycles.
func (h *Handler) RecordCycle(w http.ResponseWriter, r *http.Request) {
	id := tv.TurnoverID(r.PathValue("id"))
	var req recordCycleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	c := tv.TurnoverCycle{
		StartedAt:          req.StartedAt,
		EndedAt:            req.EndedAt,
		AdvancePence:       req.AdvancePence,
		ReturnedPence:      req.ReturnedPence,
		ProductionMinutes:  req.ProductionMinutes,
		CirculationMinutes: req.CirculationMinutes,
	}
	created, err := h.Turnovers.RecordCycle(r.Context(), id, c)
	if err != nil {
		if errors.Is(err, tv.ErrCycleNotComplete) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// RecomputeNumber handles POST /v1/turnovers/{id}/recompute-number.
func (h *Handler) RecomputeNumber(w http.ResponseWriter, r *http.Request) {
	id := tv.TurnoverID(r.PathValue("id"))
	tn, err := h.Turnovers.RecomputeNumber(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, tv.ErrNoCompletedCycles) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tn)
}

// GetNumber handles GET /v1/turnovers/{id}/number.
func (h *Handler) GetNumber(w http.ResponseWriter, r *http.Request) {
	id := tv.TurnoverID(r.PathValue("id"))
	tn, err := h.Turnovers.GetNumber(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tn)
}
