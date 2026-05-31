package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// compoundInterestRequest is the body for POST /v1/credit/compound-interest.
type compoundInterestRequest struct {
	Principal   int64             `json:"principal"`
	RateBP      int64             `json:"rate_bp"`
	PeriodYears int64             `json:"period_years"`
	FetishForm  credit.FetishForm `json:"fetish_form"`
}

// CreateCompoundInterestSchedule handles POST /v1/credit/compound-interest.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateCompoundInterestSchedule(w http.ResponseWriter, r *http.Request) {
	var req compoundInterestRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	sched, err := credit.NewCompoundInterestSchedule(req.Principal, req.RateBP, req.PeriodYears, req.FetishForm)
	if err != nil {
		if errors.Is(err, credit.ErrNonPositivePrincipal) ||
			errors.Is(err, credit.ErrNonPositiveRateBP) ||
			errors.Is(err, credit.ErrNonPositivePeriodYears) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateCompoundInterestSchedule(r.Context(), sched)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/compound-interest/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListCompoundInterestSchedules handles GET /v1/credit/compound-interest.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListCompoundInterestSchedules(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCompoundInterestSchedules(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.CompoundInterestSchedule{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetCompoundInterestSchedule handles GET /v1/credit/compound-interest/{id}.
func (h *Handler) GetCompoundInterestSchedule(w http.ResponseWriter, r *http.Request) {
	id := credit.CompoundInterestID(r.PathValue("id"))
	sched, err := h.Store.GetCompoundInterestSchedule(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sched)
}
