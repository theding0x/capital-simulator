package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// realCapitalAccumulationRequest is the body for POST /v1/credit/real-capital-accumulation.
type realCapitalAccumulationRequest struct {
	Name                      string                      `json:"name"`
	Phase                     credit.IndustrialCyclePhase `json:"phase"`
	MoneyCapitalGrowthBP      int64                       `json:"money_capital_growth_bp"`
	RealCapitalGrowthBP       int64                       `json:"real_capital_growth_bp"`
	Corresponds               bool                        `json:"corresponds"`
	CorrespondenceExplanation string                      `json:"correspondence_explanation"`
	ReserveCapital            int64                       `json:"reserve_capital"`
	CreditLimitDescription    string                      `json:"credit_limit_description"`
	CommercialCreditContracts bool                        `json:"commercial_credit_contracts"`
	CauseIsMoneyScarcity      bool                        `json:"cause_is_money_scarcity"`
	Description               string                      `json:"description"`
}

// CreateRealCapitalAccumulation handles POST /v1/credit/real-capital-accumulation.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateRealCapitalAccumulation(w http.ResponseWriter, r *http.Request) {
	var req realCapitalAccumulationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	correspondence := credit.AccumulationCorrespondence{
		MoneyCapitalGrowthBP: req.MoneyCapitalGrowthBP,
		RealCapitalGrowthBP:  req.RealCapitalGrowthBP,
		Corresponds:          req.Corresponds,
		Explanation:          req.CorrespondenceExplanation,
	}
	creditLimit := credit.CreditLimit{
		ReserveCapital: req.ReserveCapital,
		Description:    req.CreditLimitDescription,
	}

	a, err := credit.NewRealCapitalAccumulation(
		req.Name,
		req.Phase,
		correspondence,
		creditLimit,
		req.CommercialCreditContracts,
		req.CauseIsMoneyScarcity,
		req.Description,
	)
	if err != nil {
		if errors.Is(err, credit.ErrInvalidCyclePhase) ||
			errors.Is(err, credit.ErrNegativeReserveCapital) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateRealCapitalAccumulation(r.Context(), a)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/real-capital-accumulation/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListRealCapitalAccumulations handles GET /v1/credit/real-capital-accumulation.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListRealCapitalAccumulations(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListRealCapitalAccumulations(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.RealCapitalAccumulation{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetRealCapitalAccumulation handles GET /v1/credit/real-capital-accumulation/{id}.
func (h *Handler) GetRealCapitalAccumulation(w http.ResponseWriter, r *http.Request) {
	id := credit.RealCapitalAccumulationID(r.PathValue("id"))
	a, err := h.Store.GetRealCapitalAccumulation(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}
