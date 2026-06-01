package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// ListCompetitionIllusions handles GET /v1/revenue/competition-illusions.
// Returns {"items":[...]} — never null.
func (h *Handler) ListCompetitionIllusions(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCompetitionIllusions(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []revenue.CompetitionIllusion{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createCostPriceFetishRequest is the body for POST /v1/revenue/cost-price-fetishes.
type createCostPriceFetishRequest struct {
	ConstantCapitalBP int64 `json:"constant_capital_bp"`
	VariableCapitalBP int64 `json:"variable_capital_bp"`
	ProfitBP          int64 `json:"profit_bp"`
}

// CreateCostPriceFetish handles POST /v1/revenue/cost-price-fetishes. It computes the
// cost-price (c+v) and the actual value (cost-price+profit). 400 bad JSON, 422 invalid, 201.
func (h *Handler) CreateCostPriceFetish(w http.ResponseWriter, r *http.Request) {
	var req createCostPriceFetishRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ConstantCapitalBP < 0 || req.VariableCapitalBP < 0 || req.ProfitBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "c, v, and profit must be >= 0")
		return
	}
	cost := revenue.ComputeCostPrice(req.ConstantCapitalBP, req.VariableCapitalBP)
	saved, err := h.Store.CreateCostPriceFetish(r.Context(), revenue.CostPriceFetish{
		ConstantCapitalBP: req.ConstantCapitalBP,
		VariableCapitalBP: req.VariableCapitalBP,
		CostPriceBP:       cost,
		ProfitBP:          req.ProfitBP,
		ActualValueBP:     revenue.ComputeActualValue(cost, req.ProfitBP),
	})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/revenue/cost-price-fetishes")
	writeJSON(w, http.StatusCreated, saved)
}

// createAnnualRevenueAccountRequest is the body for POST /v1/revenue/annual-accounts.
type createAnnualRevenueAccountRequest struct {
	ConstantCapitalLM    int64 `json:"constant_capital_lm"`
	WagesShareLM         int64 `json:"wages_share_lm"`
	ProfitAndRentShareLM int64 `json:"profit_and_rent_share_lm"`
}

// CreateAnnualRevenueAccount handles POST /v1/revenue/annual-accounts. It computes the
// new value (wages+profit&rent = v+s) and total social product (c + new value). 400/422/201.
func (h *Handler) CreateAnnualRevenueAccount(w http.ResponseWriter, r *http.Request) {
	var req createAnnualRevenueAccountRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ConstantCapitalLM < 0 || req.WagesShareLM < 0 || req.ProfitAndRentShareLM < 0 {
		writeError(w, http.StatusUnprocessableEntity, "c, wages, and profit&rent must be >= 0")
		return
	}
	newValue := revenue.ComputeNewValueCreated(req.WagesShareLM, req.ProfitAndRentShareLM)
	saved, err := h.Store.CreateAnnualRevenueAccount(r.Context(), revenue.AnnualRevenueAccount{
		NewValueCreatedLM:    newValue,
		WagesShareLM:         req.WagesShareLM,
		ProfitAndRentShareLM: req.ProfitAndRentShareLM,
		ConstantCapitalLM:    req.ConstantCapitalLM,
		TotalSocialProductLM: revenue.ComputeTotalSocialProduct(req.ConstantCapitalLM, newValue),
	})
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/revenue/annual-accounts")
	writeJSON(w, http.StatusCreated, saved)
}

// ListValueAppearanceContrasts handles GET /v1/revenue/value-contrasts.
// Returns {"items":[...]} — never null.
func (h *Handler) ListValueAppearanceContrasts(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListValueAppearanceContrasts(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []revenue.ValueAppearanceContrast{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
