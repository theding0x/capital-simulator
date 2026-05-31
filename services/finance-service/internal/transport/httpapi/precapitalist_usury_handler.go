package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// createUsurersCapitalRequest is the body for POST /v1/credit/usurers-capital.
type createUsurersCapitalRequest struct {
	Name           string   `json:"name"`
	Stage          int64    `json:"stage"`
	WageLabour     bool     `json:"wage_labour"`
	Borrowers      []string `json:"borrowers"`
	InterestRateBP int64    `json:"interest_rate_bp"`
	SubordinatedTo string   `json:"subordinated_to"`
	Description    string   `json:"description"`
}

// CreateUsurersCapital handles POST /v1/credit/usurers-capital.
// Returns 400 on bad JSON, 422 on domain-rule violations, 201 + Location on success.
func (h *Handler) CreateUsurersCapital(w http.ResponseWriter, r *http.Request) {
	var req createUsurersCapitalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	u, err := credit.NewUsurersCapital(
		req.Name,
		credit.UsuryDevelopmentStage(req.Stage),
		req.WageLabour,
		req.Borrowers,
		req.InterestRateBP,
		req.SubordinatedTo,
		req.Description,
	)
	if err != nil {
		if errors.Is(err, credit.ErrInvalidUsuryStage) ||
			errors.Is(err, credit.ErrWageLabourInPreCapitalist) ||
			errors.Is(err, credit.ErrNegativeHistoricalRate) ||
			errors.Is(err, credit.ErrMissingSubordination) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateUsurersCapital(r.Context(), u)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/usurers-capital/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListUsurersCapitals handles GET /v1/credit/usurers-capital.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListUsurersCapitals(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListUsurersCapitals(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.UsurersCapital{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetUsurersCapital handles GET /v1/credit/usurers-capital/{id}.
func (h *Handler) GetUsurersCapital(w http.ResponseWriter, r *http.Request) {
	id := credit.UsurersCapitalID(r.PathValue("id"))
	u, err := h.Store.GetUsurersCapital(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, u)
}
