package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// noteIssueConstraintRequest is the body for POST /v1/credit/note-issue-constraints.
// max_notes is computed server-side (GoldBacking + UnbackedLimit) and must not
// be supplied by the caller.
type noteIssueConstraintRequest struct {
	Name           string               `json:"name"`
	Regime         credit.BankActRegime `json:"regime"`
	GoldBacking    int64                `json:"gold_backing"`
	UnbackedLimit  int64                `json:"unbacked_limit"`
	NotesBeyondMax int64                `json:"notes_beyond_max"`
	Description    string               `json:"description"`
}

// CreateNoteIssueConstraint handles POST /v1/credit/note-issue-constraints.
// Returns 400 on bad JSON, 422 on invalid regime, 201 + Location on success.
func (h *Handler) CreateNoteIssueConstraint(w http.ResponseWriter, r *http.Request) {
	var req noteIssueConstraintRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	n, err := credit.NewNoteIssueConstraint(req.Name, req.Regime, req.GoldBacking, req.UnbackedLimit, req.NotesBeyondMax, req.Description)
	if err != nil {
		if errors.Is(err, credit.ErrInvalidBankActRegime) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateNoteIssueConstraint(r.Context(), n)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/note-issue-constraints/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListNoteIssueConstraints handles GET /v1/credit/note-issue-constraints.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListNoteIssueConstraints(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListNoteIssueConstraints(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.NoteIssueConstraint{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetNoteIssueConstraint handles GET /v1/credit/note-issue-constraints/{id}.
func (h *Handler) GetNoteIssueConstraint(w http.ResponseWriter, r *http.Request) {
	id := credit.NoteIssueConstraintID(r.PathValue("id"))
	n, err := h.Store.GetNoteIssueConstraint(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, n)
}
