package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// moneyCapitalAccumulationRequest is the body for POST /v1/credit/money-capital-accumulation.
type moneyCapitalAccumulationRequest struct {
	Name             string                      `json:"name"`
	Phase            credit.IndustrialCyclePhase `json:"phase"`
	InterestRateBP   int64                       `json:"interest_rate_bp"`
	Period           string                      `json:"period"`
	LoanableAmount   int64                       `json:"loanable_amount"`
	LoanableAbundant bool                        `json:"loanable_abundant"`
	IdleAmount       int64                       `json:"idle_amount"`
	IdleDescription  string                      `json:"idle_description"`
	GapBP            int64                       `json:"gap_bp"`
	GapDescription   string                      `json:"gap_description"`
	Description      string                      `json:"description"`
}

// CreateMoneyCapitalAccumulation handles POST /v1/credit/money-capital-accumulation.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateMoneyCapitalAccumulation(w http.ResponseWriter, r *http.Request) {
	var req moneyCapitalAccumulationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	cycleObs := credit.InterestCycleObservation{
		Phase:          req.Phase,
		InterestRateBP: req.InterestRateBP,
		Period:         req.Period,
	}
	loanableSupply := credit.LoanableCapitalSupply{
		Amount:   req.LoanableAmount,
		Abundant: req.LoanableAbundant,
	}
	idleCapital := credit.IdleIndustrialCapital{
		Amount:      req.IdleAmount,
		Description: req.IdleDescription,
	}
	gap := credit.AccumulationGap{
		GapBP:       req.GapBP,
		Description: req.GapDescription,
	}

	a, err := credit.NewMoneyCapitalAccumulation(req.Name, cycleObs, loanableSupply, idleCapital, gap, req.Description)
	if err != nil {
		if errors.Is(err, credit.ErrInvalidCyclePhase) ||
			errors.Is(err, credit.ErrNegativeInterestRate) ||
			errors.Is(err, credit.ErrNegativeIdleAmount) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateMoneyCapitalAccumulation(r.Context(), a)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/money-capital-accumulation/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListMoneyCapitalAccumulations handles GET /v1/credit/money-capital-accumulation.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListMoneyCapitalAccumulations(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListMoneyCapitalAccumulations(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.MoneyCapitalAccumulation{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetMoneyCapitalAccumulation handles GET /v1/credit/money-capital-accumulation/{id}.
func (h *Handler) GetMoneyCapitalAccumulation(w http.ResponseWriter, r *http.Request) {
	id := credit.MoneyCapitalAccumulationID(r.PathValue("id"))
	a, err := h.Store.GetMoneyCapitalAccumulation(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}
