package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

// ── BankCapital ───────────────────────────────────────────────────────────────

// bankCapitalComponentRequest mirrors credit.BankCapitalComponent in a request body.
type bankCapitalComponentRequest struct {
	Amount      int64  `json:"amount"`
	Kind        int    `json:"kind"`
	Description string `json:"description"`
}

// bankCapitalRequest is the body for POST /v1/credit/bank-capital.
type bankCapitalRequest struct {
	Name             string                        `json:"name"`
	CashAmount       int64                         `json:"cash_amount"`
	SecuritiesAmount int64                         `json:"securities_amount"`
	Components       []bankCapitalComponentRequest `json:"components"`
	Description      string                        `json:"description"`
}

// CreateBankCapital handles POST /v1/credit/bank-capital.
// Returns 400 on bad JSON, 422 on invariant violations, 201+Location on success.
func (h *Handler) CreateBankCapital(w http.ResponseWriter, r *http.Request) {
	var req bankCapitalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	components := make([]credit.BankCapitalComponent, 0, len(req.Components))
	for _, c := range req.Components {
		components = append(components, credit.BankCapitalComponent{
			Amount:      c.Amount,
			Kind:        credit.BankCapitalComponentKind(c.Kind),
			Description: c.Description,
		})
	}

	bc, err := credit.NewBankCapital(req.Name, req.CashAmount, req.SecuritiesAmount, components, req.Description)
	if err != nil {
		if errors.Is(err, credit.ErrInvalidComponentKind) ||
			errors.Is(err, credit.ErrBankCapitalMismatch) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateBankCapital(r.Context(), bc)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/bank-capital/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListBankCapitals handles GET /v1/credit/bank-capital.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListBankCapitals(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListBankCapitals(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.BankCapital{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetBankCapital handles GET /v1/credit/bank-capital/{id}.
func (h *Handler) GetBankCapital(w http.ResponseWriter, r *http.Request) {
	id := credit.BankCapitalID(r.PathValue("id"))
	bc, err := h.Store.GetBankCapital(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, bc)
}

// ── FictitiousCapitalValuation ────────────────────────────────────────────────

// fictitiousCapitalValuationRequest is the body for
// POST /v1/credit/fictitious-capital-valuation.
type fictitiousCapitalValuationRequest struct {
	Name           string `json:"name"`
	NominalValue   int64  `json:"nominal_value"`
	AnnualIncome   int64  `json:"annual_income"`
	InterestRateBP int64  `json:"interest_rate_bp"`
	DividendRateBP int64  `json:"dividend_rate_bp"`
	Description    string `json:"description"`
}

// CreateFictitiousCapitalValuation handles POST /v1/credit/fictitious-capital-valuation.
// Returns 400 on bad JSON, 422 on invariant violations, 201+Location on success.
func (h *Handler) CreateFictitiousCapitalValuation(w http.ResponseWriter, r *http.Request) {
	var req fictitiousCapitalValuationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	v, err := credit.NewFictitiousCapitalValuation(
		req.Name,
		req.NominalValue, req.AnnualIncome, req.InterestRateBP, req.DividendRateBP,
		req.Description,
	)
	if err != nil {
		if errors.Is(err, credit.ErrNonPositiveInterestRate) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateFictitiousCapitalValuation(r.Context(), v)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/fictitious-capital-valuation/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListFictitiousCapitalValuations handles GET /v1/credit/fictitious-capital-valuation.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListFictitiousCapitalValuations(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListFictitiousCapitalValuations(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.FictitiousCapitalValuation{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetFictitiousCapitalValuation handles GET /v1/credit/fictitious-capital-valuation/{id}.
func (h *Handler) GetFictitiousCapitalValuation(w http.ResponseWriter, r *http.Request) {
	id := credit.FictitiousCapitalValuationID(r.PathValue("id"))
	v, err := h.Store.GetFictitiousCapitalValuation(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}
