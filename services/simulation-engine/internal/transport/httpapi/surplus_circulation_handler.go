package httpapi

import (
	"errors"
	"net/http"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Vol. II Ch. 17 — The Circulation of Surplus-Value.

// --- request types ---

type createSurplusCirculationRequest struct {
	Period            string `json:"period"`
	TotalSurplusPence int64  `json:"total_surplus_pence"`
}

type recordRealisationSourceRequest struct {
	Source string `json:"source"`
	Pence  int64  `json:"pence"`
}

// --- response types ---

type realisationSourceEntryResponse struct {
	SurplusCirculationID     string `json:"surplus_circulation_id"`
	Source                   string `json:"source"`
	Pence                    int64  `json:"pence"`
	DeferredRealisationPence int64  `json:"deferred_realisation_pence"`
}

type surplusCirculationResponse struct {
	ID                     string                           `json:"id"`
	Period                 string                           `json:"period"`
	TotalSurplusPence      int64                            `json:"total_surplus_pence"`
	RealisedSurplusPence   int64                            `json:"realised_surplus_pence"`
	UnrealisedSurplusPence int64                            `json:"unrealised_surplus_pence"`
	RealisationSources     []realisationSourceEntryResponse `json:"realisation_sources"`
}

type socialCapitalAggregateResponse struct {
	Period                  string `json:"period"`
	TotalAdvancedPence      int64  `json:"total_advanced_pence"`
	TotalAnnualOutputPence  int64  `json:"total_annual_output_pence"`
	TotalAnnualSurplusPence int64  `json:"total_annual_surplus_pence"`
	DepartmentIShareBP      int64  `json:"department_i_share_bp"`
	DepartmentIIShareBP     int64  `json:"department_ii_share_bp"`
}

type realisationPuzzleResponse struct {
	SurplusCirculationID string `json:"surplus_circulation_id"`
	PuzzleStatement      string `json:"puzzle_statement"`
	ResolvedInChapter    int64  `json:"resolved_in_chapter"`
}

// --- helpers ---

func scToResponse(sc repro.SurplusCirculation) surplusCirculationResponse {
	sources := make([]realisationSourceEntryResponse, len(sc.RealisationSources))
	for i, e := range sc.RealisationSources {
		sources[i] = realisationSourceEntryResponse{
			SurplusCirculationID:     string(e.SurplusCirculationID),
			Source:                   string(e.Source),
			Pence:                    int64(e.Pence),
			DeferredRealisationPence: int64(e.DeferredRealisationPence),
		}
	}
	return surplusCirculationResponse{
		ID:                     string(sc.ID),
		Period:                 sc.Period,
		TotalSurplusPence:      int64(sc.TotalSurplusPence),
		RealisedSurplusPence:   int64(sc.RealisedSurplusPence),
		UnrealisedSurplusPence: int64(sc.UnrealisedSurplusPence),
		RealisationSources:     sources,
	}
}

// --- handlers ---

// CreateSurplusCirculation handles POST /v1/reproduction/surplus-circulations.
func (h *Handler) CreateSurplusCirculation(w http.ResponseWriter, r *http.Request) {
	var req createSurplusCirculationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	if req.TotalSurplusPence < 0 {
		writeError(w, http.StatusBadRequest, "total_surplus_pence cannot be negative")
		return
	}
	sc := repro.SurplusCirculation{
		Period:            req.Period,
		TotalSurplusPence: repro.Pence(req.TotalSurplusPence),
	}
	created, err := h.SurplusCirculations.CreateSurplusCirculation(r.Context(), sc)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "surplus circulation already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/reproduction/surplus-circulations/"+string(created.ID))
	writeJSON(w, http.StatusCreated, scToResponse(created))
}

// GetSurplusCirculation handles GET /v1/reproduction/surplus-circulations/{id}.
func (h *Handler) GetSurplusCirculation(w http.ResponseWriter, r *http.Request) {
	id := repro.SurplusCirculationID(r.PathValue("id"))
	sc, err := h.SurplusCirculations.GetSurplusCirculation(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "surplus circulation not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, scToResponse(sc))
}

// ListSurplusCirculations handles GET /v1/reproduction/surplus-circulations.
func (h *Handler) ListSurplusCirculations(w http.ResponseWriter, r *http.Request) {
	list, err := h.SurplusCirculations.ListSurplusCirculations(r.Context())
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	items := make([]surplusCirculationResponse, len(list))
	for i, sc := range list {
		items[i] = scToResponse(sc)
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// RecordRealisationSource handles
// POST /v1/reproduction/surplus-circulations/{id}/realisation-source.
func (h *Handler) RecordRealisationSource(w http.ResponseWriter, r *http.Request) {
	id := repro.SurplusCirculationID(r.PathValue("id"))
	var req recordRealisationSourceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Source == "" {
		writeError(w, http.StatusBadRequest, "source is required")
		return
	}
	if req.Pence < 0 {
		writeError(w, http.StatusBadRequest, "pence cannot be negative")
		return
	}
	entry := repro.RealisationSourceEntry{
		Source: repro.RealisationSource(req.Source),
		Pence:  repro.Pence(req.Pence),
	}
	recorded, err := h.SurplusCirculations.RecordRealisationSource(r.Context(), id, entry)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "surplus circulation not found")
			return
		}
		if errors.Is(err, repro.ErrUnknownSource) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, repro.ErrRealisationExceedsTotal) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, realisationSourceEntryResponse{
		SurplusCirculationID:     string(recorded.SurplusCirculationID),
		Source:                   string(recorded.Source),
		Pence:                    int64(recorded.Pence),
		DeferredRealisationPence: int64(recorded.DeferredRealisationPence),
	})
}

// GetSocialAggregate handles GET /v1/reproduction/social-aggregate?period=.
func (h *Handler) GetSocialAggregate(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		writeError(w, http.StatusBadRequest, "period query parameter is required")
		return
	}
	agg, err := h.SurplusCirculations.SnapshotAggregate(r.Context(), period)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, socialCapitalAggregateResponse{
		Period:                  agg.Period,
		TotalAdvancedPence:      int64(agg.TotalAdvancedPence),
		TotalAnnualOutputPence:  int64(agg.TotalAnnualOutputPence),
		TotalAnnualSurplusPence: int64(agg.TotalAnnualSurplusPence),
		DepartmentIShareBP:      agg.DepartmentIShareBP,
		DepartmentIIShareBP:     agg.DepartmentIIShareBP,
	})
}

// ListRealisationPuzzles handles GET /v1/reproduction/realisation-puzzles.
func (h *Handler) ListRealisationPuzzles(w http.ResponseWriter, r *http.Request) {
	puzzles, err := h.SurplusCirculations.ListRealisationPuzzles(r.Context())
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	items := make([]realisationPuzzleResponse, len(puzzles))
	for i, p := range puzzles {
		items[i] = realisationPuzzleResponse{
			SurplusCirculationID: string(p.SurplusCirculationID),
			PuzzleStatement:      p.PuzzleStatement,
			ResolvedInChapter:    p.ResolvedInChapter,
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
