package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// interestBearingCapitalRequest is the body for POST /v1/credit/interest-bearing-capital.
type interestBearingCapitalRequest struct {
	MoneyAdvanced  int64 `json:"money_advanced"`
	InterestRateBP int64 `json:"interest_rate_bp"`
	LoanTermDays   int64 `json:"loan_term_days"`
}

// CreateInterestBearingCapital handles POST /v1/credit/interest-bearing-capital.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateInterestBearingCapital(w http.ResponseWriter, r *http.Request) {
	var req interestBearingCapitalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	ibc, err := credit.NewInterestBearingCapital(req.MoneyAdvanced, req.InterestRateBP, req.LoanTermDays)
	if err != nil {
		if errors.Is(err, credit.ErrNonPositiveMoney) ||
			errors.Is(err, credit.ErrNonPositiveRate) ||
			errors.Is(err, credit.ErrNonPositiveTerm) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateInterestBearingCapital(r.Context(), ibc)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/interest-bearing-capital/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListInterestBearingCapitals handles GET /v1/credit/interest-bearing-capital.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListInterestBearingCapitals(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListInterestBearingCapitals(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.InterestBearingCapital{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetInterestBearingCapital handles GET /v1/credit/interest-bearing-capital/{id}.
func (h *Handler) GetInterestBearingCapital(w http.ResponseWriter, r *http.Request) {
	id := credit.InterestBearingCapitalID(r.PathValue("id"))
	ibc, err := h.Store.GetInterestBearingCapital(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ibc)
}
