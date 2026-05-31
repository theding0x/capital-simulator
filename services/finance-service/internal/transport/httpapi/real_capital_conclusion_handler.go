package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// capitalReleaseRequest is the body for POST /v1/credit/capital-release.
// The embedded LoanCapitalOverstatement is flattened into top-level snake_case
// fields (overstatement_pct, overstatement_description).
type capitalReleaseRequest struct {
	Name                     string                     `json:"name"`
	ReleasedAmount           int64                      `json:"released_amount"`
	Cause                    credit.CapitalReleaseCause `json:"cause"`
	Reinvested               bool                       `json:"reinvested"`
	OverstatementPct         int64                      `json:"overstatement_pct"`
	OverstatementDescription string                     `json:"overstatement_description"`
	Description              string                     `json:"description"`
}

// CreateCapitalRelease handles POST /v1/credit/capital-release.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateCapitalRelease(w http.ResponseWriter, r *http.Request) {
	var req capitalReleaseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	overstatement := credit.LoanCapitalOverstatement{
		OverstatementPct: req.OverstatementPct,
		Description:      req.OverstatementDescription,
	}

	c, err := credit.NewCapitalRelease(
		req.Name,
		req.ReleasedAmount,
		req.Cause,
		req.Reinvested,
		overstatement,
		req.Description,
	)
	if err != nil {
		if errors.Is(err, credit.ErrNegativeReleasedAmount) ||
			errors.Is(err, credit.ErrInvalidReleaseCause) ||
			errors.Is(err, credit.ErrNegativeOverstatement) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateCapitalRelease(r.Context(), c)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/capital-release/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListCapitalReleases handles GET /v1/credit/capital-release.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListCapitalReleases(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCapitalReleases(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.CapitalRelease{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetCapitalRelease handles GET /v1/credit/capital-release/{id}.
func (h *Handler) GetCapitalRelease(w http.ResponseWriter, r *http.Request) {
	id := credit.CapitalReleaseID(r.PathValue("id"))
	c, err := h.Store.GetCapitalRelease(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}
