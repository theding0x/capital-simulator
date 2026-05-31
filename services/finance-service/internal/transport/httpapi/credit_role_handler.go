package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

// ── StockCompany ──────────────────────────────────────────────────────────────

// stockCompanyRequest is the body for POST /v1/credit/stock-companies.
type stockCompanyRequest struct {
	Name            string `json:"name"`
	FixedCapital    int64  `json:"fixed_capital"`
	FloatingCapital int64  `json:"floating_capital"`
	TotalCapital    int64  `json:"total_capital"`
	ManagerName     string `json:"manager_name"`
	ManagerTitle    string `json:"manager_title"`
	ManagerWage     int64  `json:"manager_wage_minutes"`
	Description     string `json:"description"`
}

// CreateStockCompany handles POST /v1/credit/stock-companies.
// Returns 400 on bad JSON, 422 on invariant violations, 201+Location on success.
func (h *Handler) CreateStockCompany(w http.ResponseWriter, r *http.Request) {
	var req stockCompanyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	mgr := credit.Manager{
		Name:        req.ManagerName,
		Title:       req.ManagerTitle,
		WageMinutes: req.ManagerWage,
	}
	sc, err := credit.NewStockCompany(
		req.Name,
		req.FixedCapital, req.FloatingCapital, req.TotalCapital,
		mgr,
		req.Description,
	)
	if err != nil {
		if errors.Is(err, credit.ErrCapitalMismatch) ||
			errors.Is(err, credit.ErrNonPositiveWage) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateStockCompany(r.Context(), sc)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/stock-companies/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListStockCompanies handles GET /v1/credit/stock-companies.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListStockCompanies(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListStockCompanies(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.StockCompany{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetStockCompany handles GET /v1/credit/stock-companies/{id}.
func (h *Handler) GetStockCompany(w http.ResponseWriter, r *http.Request) {
	id := credit.StockCompanyID(r.PathValue("id"))
	sc, err := h.Store.GetStockCompany(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sc)
}

// ── CooperativeFactory ────────────────────────────────────────────────────────

// cooperativeFactoryRequest is the body for POST /v1/credit/cooperative-factories.
type cooperativeFactoryRequest struct {
	Name            string `json:"name"`
	WorkerOwned     bool   `json:"worker_owned"`
	WorkerCount     int64  `json:"worker_count"`
	CapitalAdvanced int64  `json:"capital_advanced"`
	Description     string `json:"description"`
}

// CreateCooperativeFactory handles POST /v1/credit/cooperative-factories.
// Returns 400 on bad JSON, 422 on invariant violations, 201+Location on success.
func (h *Handler) CreateCooperativeFactory(w http.ResponseWriter, r *http.Request) {
	var req cooperativeFactoryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	cf, err := credit.NewCooperativeFactory(
		req.Name,
		req.WorkerOwned,
		req.WorkerCount,
		req.CapitalAdvanced,
		req.Description,
	)
	if err != nil {
		if errors.Is(err, credit.ErrNotWorkerOwned) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateCooperativeFactory(r.Context(), cf)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/cooperative-factories/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListCooperativeFactories handles GET /v1/credit/cooperative-factories.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListCooperativeFactories(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCooperativeFactories(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.CooperativeFactory{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
