package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

// createAnnualCompositionRequest is the body for POST /v1/revenue/annual-compositions.
type createAnnualCompositionRequest struct {
	ConstantCapitalLabourMinutes int64 `json:"constant_capital_labour_minutes"`
	VariableCapitalLabourMinutes int64 `json:"variable_capital_labour_minutes"`
	SurplusValueLabourMinutes    int64 `json:"surplus_value_labour_minutes"`
}

// annualCompositionResponse bundles a composition with its derived revenue limit.
type annualCompositionResponse struct {
	Composition  revenue.AnnualValueComposition `json:"composition"`
	RevenueLimit revenue.RevenueLimit           `json:"revenue_limit"`
}

// CreateAnnualValueComposition handles POST /v1/revenue/annual-compositions. It
// computes the total value (c+v+s) and new value (v+s), persists the composition,
// and derives + persists the revenue limit. 400 bad JSON, 422 invalid, 201 + Location.
func (h *Handler) CreateAnnualValueComposition(w http.ResponseWriter, r *http.Request) {
	var req createAnnualCompositionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ConstantCapitalLabourMinutes < 0 || req.VariableCapitalLabourMinutes < 0 || req.SurplusValueLabourMinutes < 0 {
		writeError(w, http.StatusUnprocessableEntity, "c, v, and s must be >= 0")
		return
	}
	comp := revenue.AnnualValueComposition{
		ConstantCapitalLabourMinutes: req.ConstantCapitalLabourMinutes,
		VariableCapitalLabourMinutes: req.VariableCapitalLabourMinutes,
		SurplusValueLabourMinutes:    req.SurplusValueLabourMinutes,
		TotalValueLabourMinutes:      revenue.ComputeTotalValue(req.ConstantCapitalLabourMinutes, req.VariableCapitalLabourMinutes, req.SurplusValueLabourMinutes),
		NewValueLabourMinutes:        revenue.ComputeNewValue(req.VariableCapitalLabourMinutes, req.SurplusValueLabourMinutes),
	}
	savedComp, err := h.Store.CreateAnnualValueComposition(r.Context(), comp)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	savedLimit, err := h.Store.CreateRevenueLimit(r.Context(), revenue.DeriveRevenueLimit(savedComp))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/revenue/annual-compositions/"+string(savedComp.ID))
	writeJSON(w, http.StatusCreated, annualCompositionResponse{Composition: savedComp, RevenueLimit: savedLimit})
}

// GetAnnualValueComposition handles GET /v1/revenue/annual-compositions/{id}.
// Returns 404 when no composition has the given ID, 200 + {composition, revenue_limit} otherwise.
func (h *Handler) GetAnnualValueComposition(w http.ResponseWriter, r *http.Request) {
	id := revenue.AnnualValueCompositionID(r.PathValue("id"))
	comp, err := h.Store.GetAnnualValueComposition(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	limit, lerr := h.Store.GetRevenueLimitByComposition(r.Context(), string(comp.ID))
	if lerr != nil {
		if !errors.Is(lerr, store.ErrNotFound) {
			writeStoreError(w, lerr)
			return
		}
		limit = revenue.DeriveRevenueLimit(comp)
	}
	writeJSON(w, http.StatusOK, annualCompositionResponse{Composition: comp, RevenueLimit: limit})
}

// createSurplusPartitionRequest is the body for POST /v1/revenue/surplus-partitions.
type createSurplusPartitionRequest struct {
	CompositionID      string `json:"composition_id"`
	ProfitEnterpriseLM int64  `json:"profit_enterprise_lm"`
	InterestLM         int64  `json:"interest_lm"`
	GroundRentLM       int64  `json:"ground_rent_lm"`
	ResidualSurplusLM  int64  `json:"residual_surplus_lm"`
}

// CreateSurplusValuePartition handles POST /v1/revenue/surplus-partitions. It looks up
// the named composition, builds the partition (total = sum of parts), and validates that
// the partition does not distribute more surplus-value than was produced. 400 bad JSON,
// 404 unknown composition, 422 invalid partition, 201 + Location.
func (h *Handler) CreateSurplusValuePartition(w http.ResponseWriter, r *http.Request) {
	var req createSurplusPartitionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.CompositionID == "" {
		writeError(w, http.StatusUnprocessableEntity, "composition_id must not be empty")
		return
	}
	if req.ProfitEnterpriseLM < 0 || req.InterestLM < 0 || req.GroundRentLM < 0 || req.ResidualSurplusLM < 0 {
		writeError(w, http.StatusUnprocessableEntity, "partition components must be >= 0")
		return
	}
	comp, err := h.Store.GetAnnualValueComposition(r.Context(), revenue.AnnualValueCompositionID(req.CompositionID))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	partition := revenue.SurplusValuePartition{
		CompositionID:      req.CompositionID,
		ProfitEnterpriseLM: req.ProfitEnterpriseLM,
		InterestLM:         req.InterestLM,
		GroundRentLM:       req.GroundRentLM,
		ResidualSurplusLM:  req.ResidualSurplusLM,
	}
	partition.TotalSurplusValueLM = revenue.SumSurplusPartition(partition)
	if err := revenue.ValidateRevenueLimit(comp, partition); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	saved, err := h.Store.CreateSurplusValuePartition(r.Context(), partition)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/revenue/surplus-partitions")
	writeJSON(w, http.StatusCreated, saved)
}
