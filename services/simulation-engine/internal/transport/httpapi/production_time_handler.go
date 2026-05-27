package httpapi

// Vol. II Ch. 13 — The Time of Production.
// Handlers for production-labour gaps, natural-process activations,
// natural subjects, and industry benchmarks.

import (
	"errors"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	wp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

// --- request / response types -----------------------------------------------

type createProductionLabourGapRequest struct {
	ProductionTimeNanos int64 `json:"production_time_nanos"`
	LabourTimeNanos     int64 `json:"labour_time_nanos"`
	GapNanos            int64 `json:"gap_nanos"`
}

type createNaturalProcessActivationRequest struct {
	Kind                      string  `json:"kind"`
	StartedAt                 string  `json:"started_at"`
	EndedAt                   *string `json:"ended_at,omitempty"`
	RequiresLabourMaintenance bool    `json:"requires_labour_maintenance"`
}

type createNaturalSubjectRequest struct {
	CommodityID string `json:"commodity_id"`
	Description string `json:"description"`
	IsLiving    bool   `json:"is_living"`
}

type createIndustryBenchmarkRequest struct {
	IndustryName               string `json:"industry_name"`
	TypicalProductionDaysCount int64  `json:"typical_production_days_count"`
	TypicalLabourDaysCount     int64  `json:"typical_labour_days_count"`
	RatioBasisPoints           int64  `json:"ratio_basis_points"`
}

// --- handlers ---------------------------------------------------------------

// RecordProductionLabourGap handles POST /v1/working-periods/{id}/production-labour-gap.
// Returns 201 on success; 400 on validation error; 409 on duplicate.
func (h *Handler) RecordProductionLabourGap(w http.ResponseWriter, r *http.Request) {
	id := wp.WorkingPeriodID(r.PathValue("id"))

	var req createProductionLabourGapRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	gap := wp.ProductionLabourGap{
		ID:                  wp.NewProductionLabourGapID(),
		WorkingPeriodID:     id,
		ProductionTimeNanos: req.ProductionTimeNanos,
		LabourTimeNanos:     req.LabourTimeNanos,
		GapNanos:            req.GapNanos,
	}
	if err := gap.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.ProductionTime.RecordProductionLabourGap(r.Context(), gap)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// GetProductionLabourGap handles GET /v1/working-periods/{id}/production-labour-gap.
func (h *Handler) GetProductionLabourGap(w http.ResponseWriter, r *http.Request) {
	id := wp.WorkingPeriodID(r.PathValue("id"))
	gap, err := h.ProductionTime.GetProductionLabourGap(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "production labour gap not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, gap)
}

// RecordNaturalProcessActivation handles POST /v1/working-periods/{id}/natural-activations.
// Returns 201 on success; 400 on validation error.
func (h *Handler) RecordNaturalProcessActivation(w http.ResponseWriter, r *http.Request) {
	id := wp.WorkingPeriodID(r.PathValue("id"))

	var req createNaturalProcessActivationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	startedAt, err := time.Parse(time.RFC3339, req.StartedAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, "started_at must be RFC 3339: "+err.Error())
		return
	}

	act := wp.NaturalProcessActivation{
		ID:                        wp.NewNaturalProcessActivationID(),
		WorkingPeriodID:           id,
		Kind:                      wp.NaturalProcessKind(req.Kind),
		StartedAt:                 startedAt,
		RequiresLabourMaintenance: req.RequiresLabourMaintenance,
	}
	if req.EndedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.EndedAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "ended_at must be RFC 3339: "+err.Error())
			return
		}
		act.EndedAt = &t
	}
	if err := act.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.ProductionTime.RecordNaturalProcessActivation(r.Context(), act)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// ListNaturalProcessActivations handles GET /v1/working-periods/{id}/natural-activations.
func (h *Handler) ListNaturalProcessActivations(w http.ResponseWriter, r *http.Request) {
	id := wp.WorkingPeriodID(r.PathValue("id"))
	list, err := h.ProductionTime.ListNaturalProcessActivations(r.Context(), id)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// RecordNaturalSubject handles POST /v1/working-periods/{id}/natural-subject.
// Returns 201 on success; 400 on validation; 409 on duplicate.
func (h *Handler) RecordNaturalSubject(w http.ResponseWriter, r *http.Request) {
	id := wp.WorkingPeriodID(r.PathValue("id"))

	var req createNaturalSubjectRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	sub := wp.NaturalSubject{
		ID:              wp.NewNaturalSubjectID(),
		WorkingPeriodID: id,
		CommodityID:     req.CommodityID,
		Description:     req.Description,
		IsLiving:        req.IsLiving,
	}
	if err := sub.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.ProductionTime.RecordNaturalSubject(r.Context(), sub)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// GetNaturalSubject handles GET /v1/working-periods/{id}/natural-subject.
func (h *Handler) GetNaturalSubject(w http.ResponseWriter, r *http.Request) {
	id := wp.WorkingPeriodID(r.PathValue("id"))
	sub, err := h.ProductionTime.GetNaturalSubject(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "natural subject not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sub)
}

// CreateIndustryBenchmark handles POST /v1/industry-benchmarks.
// Returns 201 + Location header on success.
func (h *Handler) CreateIndustryBenchmark(w http.ResponseWriter, r *http.Request) {
	var req createIndustryBenchmarkRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	ind := wp.NaturalProcessIndustry{
		ID:                         wp.NewNaturalProcessIndustryID(),
		IndustryName:               req.IndustryName,
		TypicalProductionDaysCount: req.TypicalProductionDaysCount,
		TypicalLabourDaysCount:     req.TypicalLabourDaysCount,
		RatioBasisPoints:           req.RatioBasisPoints,
	}
	if err := ind.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.ProductionTime.CreateIndustryBenchmark(r.Context(), ind)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/industry-benchmarks/"+string(created.ID))
	writeJSON(w, http.StatusCreated, created)
}

// ListIndustryBenchmarks handles GET /v1/industry-benchmarks.
// Returns all benchmarks ordered by ratio_basis_points ascending.
func (h *Handler) ListIndustryBenchmarks(w http.ResponseWriter, r *http.Request) {
	list, err := h.ProductionTime.ListIndustryBenchmarks(r.Context())
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
