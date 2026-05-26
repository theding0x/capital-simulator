package httpapi

// Vol. II Ch. 12 — The Working Period.
// Handlers for working periods, their interruptions / shortenings /
// credit-financing records, and natural-working-period constraints.

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	wp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

// --- request / response types ------------------------------------------------

type createWorkingPeriodRequest struct {
	IndustrialCapitalID     string `json:"industrial_capital_id"`
	CommodityID             string `json:"commodity_id"`
	WorkingDayCount         int64  `json:"working_day_count"`
	Mode                    string `json:"mode"`
	WorkingDayLengthMinutes int64  `json:"working_day_length_minutes"`
}

type createInterruptionRequest struct {
	DeteriorationPence            int64 `json:"deterioration_pence"`
	WasteOfMeansOfProductionPence int64 `json:"waste_of_means_of_production_pence"`
}

type createShorteningRequest struct {
	Kind                        string `json:"kind"`
	OldWorkingDayCount          int64  `json:"old_working_day_count"`
	NewWorkingDayCount          int64  `json:"new_working_day_count"`
	AdditionalFixedCapitalPence int64  `json:"additional_fixed_capital_pence"`
	AdditionalLabourCount       int64  `json:"additional_labour_count"`
}

type createCreditFinancingRequest struct {
	OwnCapitalPence      int64  `json:"own_capital_pence"`
	BorrowedCapitalPence int64  `json:"borrowed_capital_pence"`
	MortgageProviderID   string `json:"mortgage_provider_id,omitempty"`
}

type createNaturalConstraintRequest struct {
	CommodityID    string `json:"commodity_id"`
	MinWorkingDays int64  `json:"min_working_days"`
	Reason         string `json:"reason"`
}

// --- handlers ----------------------------------------------------------------

// CreateWorkingPeriod handles POST /v1/working-periods.
// Returns 201 + Location header on success; 400 on validation or natural-
// constraint violation; 409 on duplicate; 500 otherwise.
func (h *Handler) CreateWorkingPeriod(w http.ResponseWriter, r *http.Request) {
	var req createWorkingPeriodRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	period := wp.WorkingPeriod{
		ID:                      wp.NewWorkingPeriodID(),
		IndustrialCapitalID:     req.IndustrialCapitalID,
		CommodityID:             req.CommodityID,
		WorkingDayCount:         req.WorkingDayCount,
		Mode:                    wp.WorkingPeriodMode(req.Mode),
		WorkingDayLengthMinutes: req.WorkingDayLengthMinutes,
	}
	if err := period.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.WorkingPeriods.CreateWorkingPeriod(r.Context(), period)
	if err != nil {
		if errors.Is(err, wp.ErrPrematureDelivery) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}

	w.Header().Set("Location", "/v1/working-periods/"+string(created.ID))
	writeJSON(w, http.StatusCreated, created)
}

// GetWorkingPeriod handles GET /v1/working-periods/{id}.
func (h *Handler) GetWorkingPeriod(w http.ResponseWriter, r *http.Request) {
	id := wp.WorkingPeriodID(r.PathValue("id"))
	period, err := h.WorkingPeriods.GetWorkingPeriod(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "working period not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, period)
}

// ListWorkingPeriods handles GET /v1/working-periods.
// Optional query params: industrial_capital_id, mode, commodity_id.
func (h *Handler) ListWorkingPeriods(w http.ResponseWriter, r *http.Request) {
	icID := r.URL.Query().Get("industrial_capital_id")
	mode := r.URL.Query().Get("mode")
	commodityID := r.URL.Query().Get("commodity_id")

	list, err := h.WorkingPeriods.ListWorkingPeriods(r.Context(), icID, mode, commodityID)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// RecordWorkingPeriodInterruption handles POST /v1/working-periods/{id}/interruptions.
// Returns 201 with the created interruption on success.
func (h *Handler) RecordWorkingPeriodInterruption(w http.ResponseWriter, r *http.Request) {
	id := wp.WorkingPeriodID(r.PathValue("id"))

	// Fetch the parent so we can validate the mode invariant.
	period, err := h.WorkingPeriods.GetWorkingPeriod(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "working period not found")
			return
		}
		h.writeServerError(w, err)
		return
	}

	var req createInterruptionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	intr := wp.WorkingPeriodInterruption{
		ID:                            wp.NewWorkingPeriodInterruptionID(),
		WorkingPeriodID:               id,
		DeteriorationPence:            req.DeteriorationPence,
		WasteOfMeansOfProductionPence: req.WasteOfMeansOfProductionPence,
	}
	if err := intr.Validate(period.Mode); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.WorkingPeriods.RecordInterruption(r.Context(), id, intr)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "working period not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// RecordWorkingPeriodShortening handles POST /v1/working-periods/{id}/shortenings.
// Returns 201 with the created shortening on success.
func (h *Handler) RecordWorkingPeriodShortening(w http.ResponseWriter, r *http.Request) {
	id := wp.WorkingPeriodID(r.PathValue("id"))

	var req createShorteningRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s := wp.WorkingPeriodShortening{
		ID:                          wp.NewWorkingPeriodShorteningID(),
		WorkingPeriodID:             id,
		Kind:                        wp.WorkingPeriodShorteningKind(req.Kind),
		OldWorkingDayCount:          req.OldWorkingDayCount,
		NewWorkingDayCount:          req.NewWorkingDayCount,
		AdditionalFixedCapitalPence: req.AdditionalFixedCapitalPence,
		AdditionalLabourCount:       req.AdditionalLabourCount,
	}
	if err := s.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.WorkingPeriods.RecordShortening(r.Context(), id, s)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "working period not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// RecordWorkingPeriodCreditFinancing handles POST /v1/working-periods/{id}/credit-financing.
// Returns 201 with the created credit-financing record on success.
func (h *Handler) RecordWorkingPeriodCreditFinancing(w http.ResponseWriter, r *http.Request) {
	id := wp.WorkingPeriodID(r.PathValue("id"))

	var req createCreditFinancingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	cf := wp.WorkingPeriodCreditFinancing{
		ID:                   wp.NewWorkingPeriodCreditFinancingID(),
		WorkingPeriodID:      id,
		OwnCapitalPence:      req.OwnCapitalPence,
		BorrowedCapitalPence: req.BorrowedCapitalPence,
		MortgageProviderID:   req.MortgageProviderID,
	}
	if err := cf.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.WorkingPeriods.RecordCreditFinancing(r.Context(), id, cf)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "working period not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// CreateNaturalConstraint handles POST /v1/natural-constraints.
// Returns 201 + Location header on success.
func (h *Handler) CreateNaturalConstraint(w http.ResponseWriter, r *http.Request) {
	var req createNaturalConstraintRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	nc := wp.NaturalWorkingPeriodConstraint{
		ID:             wp.NewNaturalWorkingPeriodConstraintID(),
		CommodityID:    req.CommodityID,
		MinWorkingDays: req.MinWorkingDays,
		Reason:         req.Reason,
	}
	if err := nc.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.WorkingPeriods.CreateNaturalConstraint(r.Context(), nc)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}

	w.Header().Set("Location", "/v1/natural-constraints/"+string(created.ID))
	writeJSON(w, http.StatusCreated, created)
}

// ListNaturalConstraints handles GET /v1/natural-constraints.
// Returns all registered natural working-period constraints.
func (h *Handler) ListNaturalConstraints(w http.ResponseWriter, r *http.Request) {
	list, err := h.WorkingPeriods.ListNaturalConstraints(r.Context())
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
