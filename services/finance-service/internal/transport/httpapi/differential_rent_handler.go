package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createPoPSurplusProfitRequest is the body for POST /v1/rent/surplus-profits.
type createPoPSurplusProfitRequest struct {
	CapitalID                     string `json:"capital_id"`
	GeneralPriceOfProductionBP    int64  `json:"general_price_of_production_bp"`
	IndividualPriceOfProductionBP int64  `json:"individual_price_of_production_bp"`
}

// CreatePoPSurplusProfit handles POST /v1/rent/surplus-profits.
// SurplusProfitBP is computed server-side; any client-supplied value is ignored.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreatePoPSurplusProfit(w http.ResponseWriter, r *http.Request) {
	var req createPoPSurplusProfitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.CapitalID == "" {
		writeError(w, http.StatusUnprocessableEntity, "capital_id must not be empty")
		return
	}
	if req.GeneralPriceOfProductionBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "general_price_of_production_bp must be >= 0")
		return
	}
	if req.IndividualPriceOfProductionBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "individual_price_of_production_bp must be >= 0")
		return
	}
	s := rent.PriceOfProductionSurplusProfit{
		CapitalID:                     req.CapitalID,
		GeneralPriceOfProductionBP:    req.GeneralPriceOfProductionBP,
		IndividualPriceOfProductionBP: req.IndividualPriceOfProductionBP,
		SurplusProfitBP:               rent.ComputeSurplusProfit(req.GeneralPriceOfProductionBP, req.IndividualPriceOfProductionBP),
	}
	saved, err := h.Store.CreatePoPSurplusProfit(r.Context(), s)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/surplus-profits/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListPoPSurplusProfits handles GET /v1/rent/surplus-profits.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListPoPSurplusProfits(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListPoPSurplusProfits(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.PriceOfProductionSurplusProfit{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createMonopolisedNaturalForceRequest is the body for POST /v1/rent/monopolised-forces.
type createMonopolisedNaturalForceRequest struct {
	Kind            string `json:"kind"`
	IsMonopolisable bool   `json:"is_monopolisable"`
	ParcelID        string `json:"parcel_id"`
}

// CreateMonopolisedNaturalForce handles POST /v1/rent/monopolised-forces.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateMonopolisedNaturalForce(w http.ResponseWriter, r *http.Request) {
	var req createMonopolisedNaturalForceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Kind == "" {
		writeError(w, http.StatusUnprocessableEntity, "kind must not be empty")
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	f := rent.MonopolisedNaturalForce{
		Kind:            req.Kind,
		IsMonopolisable: req.IsMonopolisable,
		ParcelID:        req.ParcelID,
	}
	saved, err := h.Store.CreateMonopolisedNaturalForce(r.Context(), f)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/monopolised-forces/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListMonopolisedNaturalForces handles GET /v1/rent/monopolised-forces.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListMonopolisedNaturalForces(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListMonopolisedNaturalForces(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.MonopolisedNaturalForce{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createDifferentialRentRequest is the body for POST /v1/rent/differential-rents.
type createDifferentialRentRequest struct {
	ParcelID                string `json:"parcel_id"`
	SurplusProfitBP         int64  `json:"surplus_profit_bp"`
	AnnualRentLabourMinutes int64  `json:"annual_rent_labour_minutes"`
}

// CreateDifferentialRent handles POST /v1/rent/differential-rents.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateDifferentialRent(w http.ResponseWriter, r *http.Request) {
	var req createDifferentialRentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.SurplusProfitBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "surplus_profit_bp must be >= 0")
		return
	}
	if req.AnnualRentLabourMinutes < 0 {
		writeError(w, http.StatusUnprocessableEntity, "annual_rent_labour_minutes must be >= 0")
		return
	}
	dr := rent.DifferentialRent{
		ParcelID:                req.ParcelID,
		SurplusProfitBP:         req.SurplusProfitBP,
		AnnualRentLabourMinutes: req.AnnualRentLabourMinutes,
	}
	// Ch. 38 identity: the surplus-profit IS the differential rent it becomes —
	// reject a rent whose surplus-profit (pence) and annual rent (LabourMinutes)
	// denote different money-magnitudes.
	if err := dr.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	saved, err := h.Store.CreateDifferentialRent(r.Context(), dr)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/differential-rents/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListDifferentialRents handles GET /v1/rent/differential-rents.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListDifferentialRents(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListDifferentialRents(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.DifferentialRent{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createCapitalisedRentPriceRequest is the body for POST /v1/rent/capitalised-prices.
type createCapitalisedRentPriceRequest struct {
	ParcelID                string `json:"parcel_id"`
	AnnualRentLabourMinutes int64  `json:"annual_rent_labour_minutes"`
	InterestRateBP          int64  `json:"interest_rate_bp"`
}

// CreateCapitalisedRentPrice handles POST /v1/rent/capitalised-prices.
// CapitalisedPriceLabourMinutes is computed server-side; any client-supplied value is ignored.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateCapitalisedRentPrice(w http.ResponseWriter, r *http.Request) {
	var req createCapitalisedRentPriceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.AnnualRentLabourMinutes < 0 {
		writeError(w, http.StatusUnprocessableEntity, "annual_rent_labour_minutes must be >= 0")
		return
	}
	if req.InterestRateBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "interest_rate_bp must be > 0")
		return
	}
	crp := rent.CapitalisedRentPrice{
		ParcelID:                      req.ParcelID,
		AnnualRentLabourMinutes:       req.AnnualRentLabourMinutes,
		InterestRateBP:                req.InterestRateBP,
		CapitalisedPriceLabourMinutes: rent.ComputeCapitalisedPrice(req.AnnualRentLabourMinutes, req.InterestRateBP),
	}
	saved, err := h.Store.CreateCapitalisedRentPrice(r.Context(), crp)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/capitalised-prices/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListCapitalisedRentPrices handles GET /v1/rent/capitalised-prices.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListCapitalisedRentPrices(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCapitalisedRentPrices(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.CapitalisedRentPrice{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
