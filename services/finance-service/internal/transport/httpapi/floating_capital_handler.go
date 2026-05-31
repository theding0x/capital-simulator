package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// floatingCapitalRequest is the body for POST /v1/credit/floating-capital.
// The embedded LoanCapitalVelocity and CirculationReserveSaving are flattened
// into top-level snake_case fields; loan_capital_created is derived server-side.
type floatingCapitalRequest struct {
	Name                     string                       `json:"name"`
	Amount                   int64                        `json:"amount"`
	Source                   credit.FloatingCapitalSource `json:"source"`
	VelocityMoneyAmount      int64                        `json:"velocity_money_amount"`
	VelocityTimesLent        int64                        `json:"velocity_times_lent"`
	ReserveSavingAmount      int64                        `json:"reserve_saving_amount"`
	ReserveSavingDescription string                       `json:"reserve_saving_description"`
	Description              string                       `json:"description"`
}

// CreateFloatingCapital handles POST /v1/credit/floating-capital.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateFloatingCapital(w http.ResponseWriter, r *http.Request) {
	var req floatingCapitalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	velocity, err := credit.NewLoanCapitalVelocity(req.VelocityMoneyAmount, req.VelocityTimesLent)
	if err != nil {
		if errors.Is(err, credit.ErrTimesLentTooLow) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	reserveSaving := credit.CirculationReserveSaving{
		Amount:      req.ReserveSavingAmount,
		Description: req.ReserveSavingDescription,
	}

	f, err := credit.NewFloatingCapital(
		req.Name,
		req.Amount,
		req.Source,
		velocity,
		reserveSaving,
		req.Description,
	)
	if err != nil {
		if errors.Is(err, credit.ErrNegativeFloatingAmount) ||
			errors.Is(err, credit.ErrInvalidFloatingSource) ||
			errors.Is(err, credit.ErrNegativeReserveSaving) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateFloatingCapital(r.Context(), f)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/floating-capital/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListFloatingCapitals handles GET /v1/credit/floating-capital.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListFloatingCapitals(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListFloatingCapitals(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.FloatingCapital{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetFloatingCapital handles GET /v1/credit/floating-capital/{id}.
func (h *Handler) GetFloatingCapital(w http.ResponseWriter, r *http.Request) {
	id := credit.FloatingCapitalID(r.PathValue("id"))
	f, err := h.Store.GetFloatingCapital(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, f)
}
