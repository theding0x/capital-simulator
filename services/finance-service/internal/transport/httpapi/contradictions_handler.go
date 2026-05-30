package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/tendency"
)

// crisisRequest is the body for POST /v1/tendency/crisis. The caller supplies a
// constant-capital writedown (percent, 1–99) and a pre-crisis rate of profit
// (basis points); the restored post-crisis rate is derived server-side.
type crisisRequest struct {
	ConstantCapitalWritedown int64 `json:"constant_capital_writedown"`
	PreCrisisProfitRate      int64 `json:"pre_crisis_profit_rate"`
}

// contradictionRequest is the body for POST /v1/tendency/contradiction. The
// caller supplies one of the four valid kinds and an optional note; the
// IsCoexistent flag is derived server-side.
type contradictionRequest struct {
	Kind string `json:"kind"`
	Note string `json:"note"`
}

// CreateCrisis handles POST /v1/tendency/crisis. Returns 400 on bad JSON, 422 on
// an invalid writedown or negative pre-crisis rate, 201 + Location on success.
func (h *Handler) CreateCrisis(w http.ResponseWriter, r *http.Request) {
	var req crisisRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	crisis, err := tendency.NewCrisis(req.ConstantCapitalWritedown, req.PreCrisisProfitRate)
	if err != nil {
		if errors.Is(err, tendency.ErrInvalidWritedown) || errors.Is(err, tendency.ErrNegativeMagnitude) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateCrisis(r.Context(), crisis)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/tendency/crisis/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListCrises handles GET /v1/tendency/crisis. It returns a JSON object with an
// "items" array — never null.
func (h *Handler) ListCrises(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCrises(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []tendency.Crisis{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetCrisis handles GET /v1/tendency/crisis/{id}.
func (h *Handler) GetCrisis(w http.ResponseWriter, r *http.Request) {
	id := tendency.CrisisID(r.PathValue("id"))
	c, err := h.Store.GetCrisis(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// CreateInternalContradiction handles POST /v1/tendency/contradiction. Returns
// 400 on bad JSON, 422 on an invalid kind, 201 + Location on success.
func (h *Handler) CreateInternalContradiction(w http.ResponseWriter, r *http.Request) {
	var req contradictionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	contradiction, err := tendency.NewInternalContradiction(tendency.ContradictionKind(req.Kind), req.Note)
	if err != nil {
		if errors.Is(err, tendency.ErrInvalidContradictionKind) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateInternalContradiction(r.Context(), contradiction)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/tendency/contradiction/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetInternalContradiction handles GET /v1/tendency/contradiction/{id}.
func (h *Handler) GetInternalContradiction(w http.ResponseWriter, r *http.Request) {
	id := tendency.InternalContradictionID(r.PathValue("id"))
	c, err := h.Store.GetInternalContradiction(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// GetPartIIISummary handles GET /v1/tendency/summary. It assembles a
// PartIIISummary from the four Part III entity lists — composition trajectories
// (Ch. 13), counteracting scenarios (Ch. 14), crises and internal contradictions
// (Ch. 15) — reporting the count and latest (newest-first index 0) record of
// each. Pointer fields are nil when the corresponding list is empty. Returns 200.
func (h *Handler) GetPartIIISummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	trajectories, err := h.Store.ListCompositionTrajectories(ctx)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	scenarios, err := h.Store.ListCounteractingScenarios(ctx)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	crises, err := h.Store.ListCrises(ctx)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	contradictions, err := h.Store.ListInternalContradictions(ctx)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	summary := tendency.PartIIISummary{
		TrajectoryCount:    int64(len(trajectories)),
		ScenarioCount:      int64(len(scenarios)),
		CrisisCount:        int64(len(crises)),
		ContradictionCount: int64(len(contradictions)),
	}
	if len(trajectories) > 0 {
		summary.LatestTrajectory = &trajectories[0]
	}
	if len(scenarios) > 0 {
		summary.LatestScenario = &scenarios[0]
	}
	if len(crises) > 0 {
		summary.LatestCrisis = &crises[0]
	}
	if len(contradictions) > 0 {
		summary.LatestContradiction = &contradictions[0]
	}

	writeJSON(w, http.StatusOK, summary)
}
