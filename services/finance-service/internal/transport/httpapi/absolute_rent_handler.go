package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createAgriculturalCompositionRequest is the body for POST /v1/rent/compositions.
type createAgriculturalCompositionRequest struct {
	ConstantCapitalBP    int64 `json:"constant_capital_bp"`
	VariableCapitalBP    int64 `json:"variable_capital_bp"`
	SocialAverageRatioBP int64 `json:"social_average_ratio_bp"`
}

// CreateAgriculturalComposition handles POST /v1/rent/compositions.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateAgriculturalComposition(w http.ResponseWriter, r *http.Request) {
	var req createAgriculturalCompositionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.VariableCapitalBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "variable_capital_bp must be > 0")
		return
	}
	if req.SocialAverageRatioBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "social_average_ratio_bp must be > 0")
		return
	}
	if req.ConstantCapitalBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "constant_capital_bp must be >= 0")
		return
	}
	ratio := rent.ComputeCompositionRatioBP(req.ConstantCapitalBP, req.VariableCapitalBP)
	rec := rent.AgriculturalComposition{
		ConstantCapitalBP:    req.ConstantCapitalBP,
		VariableCapitalBP:    req.VariableCapitalBP,
		CompositionRatioBP:   ratio,
		SocialAverageRatioBP: req.SocialAverageRatioBP,
		IsBelowAverage:       ratio < req.SocialAverageRatioBP,
	}
	saved, err := h.Store.CreateAgriculturalComposition(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/compositions")
	writeJSON(w, http.StatusCreated, saved)
}

// ListAgriculturalCompositions handles GET /v1/rent/compositions.
// Returns {"items":[...]} — never null.
func (h *Handler) ListAgriculturalCompositions(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListAgriculturalCompositions(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.AgriculturalComposition{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createValuePriceGapRequest is the body for POST /v1/rent/value-gaps.
type createValuePriceGapRequest struct {
	ProductID                      string `json:"product_id"`
	ValueLabourMinutes             int64  `json:"value_labour_minutes"`
	PriceOfProductionLabourMinutes int64  `json:"price_of_production_labour_minutes"`
}

// CreateValuePriceGap handles POST /v1/rent/value-gaps.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateValuePriceGap(w http.ResponseWriter, r *http.Request) {
	var req createValuePriceGapRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ProductID == "" {
		writeError(w, http.StatusUnprocessableEntity, "product_id must not be empty")
		return
	}
	if req.ValueLabourMinutes < 0 {
		writeError(w, http.StatusUnprocessableEntity, "value_labour_minutes must be >= 0")
		return
	}
	if req.PriceOfProductionLabourMinutes < 0 {
		writeError(w, http.StatusUnprocessableEntity, "price_of_production_labour_minutes must be >= 0")
		return
	}
	gap := rent.ComputeValuePriceGap(req.ValueLabourMinutes, req.PriceOfProductionLabourMinutes)
	rec := rent.ValuePriceGap{
		ProductID:                      req.ProductID,
		ValueLabourMinutes:             req.ValueLabourMinutes,
		PriceOfProductionLabourMinutes: req.PriceOfProductionLabourMinutes,
		GapLabourMinutes:               gap,
	}
	saved, err := h.Store.CreateValuePriceGap(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/value-gaps")
	writeJSON(w, http.StatusCreated, saved)
}

// ListValuePriceGaps handles GET /v1/rent/value-gaps.
// Returns {"items":[...]} — never null.
func (h *Handler) ListValuePriceGaps(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListValuePriceGaps(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.ValuePriceGap{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createAbsoluteRentRequest is the body for POST /v1/rent/absolute-rents.
type createAbsoluteRentRequest struct {
	ParcelID        string `json:"parcel_id"`
	ValuePriceGapID string `json:"value_price_gap_id"`
	SurplusValueBP  int64  `json:"surplus_value_bp"`
	AverageProfitBP int64  `json:"average_profit_bp"`
}

// CreateAbsoluteRent handles POST /v1/rent/absolute-rents.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateAbsoluteRent(w http.ResponseWriter, r *http.Request) {
	var req createAbsoluteRentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.SurplusValueBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "surplus_value_bp must be >= 0")
		return
	}
	if req.AverageProfitBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "average_profit_bp must be >= 0")
		return
	}
	absoluteRent := rent.ComputeAbsoluteRentBP(req.SurplusValueBP, req.AverageProfitBP)
	rec := rent.AbsoluteRent{
		ParcelID:        req.ParcelID,
		ValuePriceGapID: req.ValuePriceGapID,
		SurplusValueBP:  req.SurplusValueBP,
		AverageProfitBP: req.AverageProfitBP,
		AbsoluteRentBP:  absoluteRent,
	}
	saved, err := h.Store.CreateAbsoluteRent(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/absolute-rents")
	writeJSON(w, http.StatusCreated, saved)
}

// GetAbsoluteRent handles GET /v1/rent/absolute-rents/{id}.
// Returns 200 or 404 via writeStoreError.
func (h *Handler) GetAbsoluteRent(w http.ResponseWriter, r *http.Request) {
	id := rent.AbsoluteRentID(r.PathValue("id"))
	ar, err := h.Store.GetAbsoluteRent(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ar)
}

// ListAbsoluteRents handles GET /v1/rent/absolute-rents.
// Returns {"items":[...]} — never null.
func (h *Handler) ListAbsoluteRents(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListAbsoluteRents(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.AbsoluteRent{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createAbsoluteRentLimitRequest is the body for POST /v1/rent/rent-limits.
type createAbsoluteRentLimitRequest struct {
	MaxRentBP    int64 `json:"max_rent_bp"`
	ActualRentBP int64 `json:"actual_rent_bp"`
}

// CreateAbsoluteRentLimit handles POST /v1/rent/rent-limits.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateAbsoluteRentLimit(w http.ResponseWriter, r *http.Request) {
	var req createAbsoluteRentLimitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.MaxRentBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "max_rent_bp must be >= 0")
		return
	}
	if req.ActualRentBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "actual_rent_bp must be >= 0")
		return
	}
	rec := rent.AbsoluteRentLimit{
		MaxRentBP:    req.MaxRentBP,
		ActualRentBP: req.ActualRentBP,
		IsAtLimit:    req.ActualRentBP >= req.MaxRentBP,
	}
	saved, err := h.Store.CreateAbsoluteRentLimit(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/rent-limits")
	writeJSON(w, http.StatusCreated, saved)
}

// ListAbsoluteRentLimits handles GET /v1/rent/rent-limits.
// Returns {"items":[...]} — never null.
func (h *Handler) ListAbsoluteRentLimits(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListAbsoluteRentLimits(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.AbsoluteRentLimit{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
