package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
)

// errSchedulerUnavailable is returned by the engine handlers when the scheduler
// (or its audit store) was not wired — a configuration error, surfaced as 500.
var errSchedulerUnavailable = errors.New("tick scheduler not configured")

// engineControlResponse is the body for POST /v1/engine/start and /stop. The
// Started/Stopped flag reports whether this call changed the run state; Status
// is always the resulting scheduler snapshot.
type engineControlResponse struct {
	Started bool                   `json:"started,omitempty"`
	Stopped bool                   `json:"stopped,omitempty"`
	Status  engine.SchedulerStatus `json:"status"`
}

// StartEngine handles POST /v1/engine/start. Idempotent: starting an
// already-running scheduler is a no-op that still returns 200 with status.
func (h *Handler) StartEngine(w http.ResponseWriter, _ *http.Request) {
	if h.Scheduler == nil {
		h.writeServerError(w, errSchedulerUnavailable)
		return
	}
	started := h.Scheduler.Start()
	writeJSON(w, http.StatusOK, engineControlResponse{Started: started, Status: h.Scheduler.Status()})
}

// StopEngine handles POST /v1/engine/stop. Idempotent: stopping an
// already-stopped scheduler is a no-op that still returns 200 with status.
func (h *Handler) StopEngine(w http.ResponseWriter, _ *http.Request) {
	if h.Scheduler == nil {
		h.writeServerError(w, errSchedulerUnavailable)
		return
	}
	stopped := h.Scheduler.Stop()
	writeJSON(w, http.StatusOK, engineControlResponse{Stopped: stopped, Status: h.Scheduler.Status()})
}

// EngineStatus handles GET /v1/engine/status: the current run state, pass
// number, last pass time, interval, and registered tickers.
func (h *Handler) EngineStatus(w http.ResponseWriter, _ *http.Request) {
	if h.Scheduler == nil {
		h.writeServerError(w, errSchedulerUnavailable)
		return
	}
	writeJSON(w, http.StatusOK, h.Scheduler.Status())
}

// ListEngineTicks handles GET /v1/engine/ticks?limit=N: the most-recent
// scheduler-pass audit rows, newest first. limit defaults to 50; limit=0
// returns every row. Always returns a JSON array, never null.
func (h *Handler) ListEngineTicks(w http.ResponseWriter, r *http.Request) {
	if h.EngineTicks == nil {
		h.writeServerError(w, errSchedulerUnavailable)
		return
	}
	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			writeError(w, http.StatusBadRequest, "limit must be a non-negative integer")
			return
		}
		limit = n
	}
	ticks, err := h.EngineTicks.ListEngineTicks(r.Context(), limit)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ticks)
}
