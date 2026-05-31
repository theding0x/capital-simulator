package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createBuildingSiteRentRequest is the body for POST /v1/rent/building-rents.
type createBuildingSiteRentRequest struct {
	ParcelID                string `json:"parcel_id"`
	LocationGrade           int    `json:"location_grade"`
	AnnualRentLabourMinutes int64  `json:"annual_rent_labour_minutes"`
	LeaseTermYears          int    `json:"lease_term_years"`
}

// CreateBuildingSiteRent handles POST /v1/rent/building-rents.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateBuildingSiteRent(w http.ResponseWriter, r *http.Request) {
	var req createBuildingSiteRentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.LocationGrade < 1 {
		writeError(w, http.StatusUnprocessableEntity, "location_grade must be >= 1")
		return
	}
	if req.AnnualRentLabourMinutes < 0 {
		writeError(w, http.StatusUnprocessableEntity, "annual_rent_labour_minutes must be >= 0")
		return
	}
	if req.LeaseTermYears < 1 {
		writeError(w, http.StatusUnprocessableEntity, "lease_term_years must be >= 1")
		return
	}
	rec := rent.BuildingSiteRent{
		ParcelID:                req.ParcelID,
		LocationGrade:           req.LocationGrade,
		AnnualRentLabourMinutes: req.AnnualRentLabourMinutes,
		LeaseTermYears:          req.LeaseTermYears,
	}
	saved, err := h.Store.CreateBuildingSiteRent(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/building-rents")
	writeJSON(w, http.StatusCreated, saved)
}

// ListBuildingSiteRents handles GET /v1/rent/building-rents.
// Returns {"items":[...]} — never null.
func (h *Handler) ListBuildingSiteRents(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListBuildingSiteRents(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.BuildingSiteRent{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createMiningRentRequest is the body for POST /v1/rent/mining-rents.
type createMiningRentRequest struct {
	ParcelID                string `json:"parcel_id"`
	OreGrade                int    `json:"ore_grade"`
	AnnualRentLabourMinutes int64  `json:"annual_rent_labour_minutes"`
	SurplusProfitBP         int64  `json:"surplus_profit_bp"`
}

// CreateMiningRent handles POST /v1/rent/mining-rents.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateMiningRent(w http.ResponseWriter, r *http.Request) {
	var req createMiningRentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.OreGrade < 1 {
		writeError(w, http.StatusUnprocessableEntity, "ore_grade must be >= 1")
		return
	}
	if req.AnnualRentLabourMinutes < 0 {
		writeError(w, http.StatusUnprocessableEntity, "annual_rent_labour_minutes must be >= 0")
		return
	}
	if req.SurplusProfitBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "surplus_profit_bp must be >= 0")
		return
	}
	rec := rent.MiningRent{
		ParcelID:                req.ParcelID,
		OreGrade:                req.OreGrade,
		AnnualRentLabourMinutes: req.AnnualRentLabourMinutes,
		SurplusProfitBP:         req.SurplusProfitBP,
	}
	saved, err := h.Store.CreateMiningRent(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/mining-rents")
	writeJSON(w, http.StatusCreated, saved)
}

// ListMiningRents handles GET /v1/rent/mining-rents.
// Returns {"items":[...]} — never null.
func (h *Handler) ListMiningRents(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListMiningRents(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.MiningRent{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createMonopolyPriceRentRequest is the body for POST /v1/rent/monopoly-rents.
type createMonopolyPriceRentRequest struct {
	ParcelID                       string `json:"parcel_id"`
	ProductKind                    string `json:"product_kind"`
	ValueLabourMinutes             int64  `json:"value_labour_minutes"`
	MonopolySalePriceLabourMinutes int64  `json:"monopoly_sale_price_labour_minutes"`
}

// CreateMonopolyPriceRent handles POST /v1/rent/monopoly-rents.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateMonopolyPriceRent(w http.ResponseWriter, r *http.Request) {
	var req createMonopolyPriceRentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.ProductKind == "" {
		writeError(w, http.StatusUnprocessableEntity, "product_kind must not be empty")
		return
	}
	if req.ValueLabourMinutes < 0 {
		writeError(w, http.StatusUnprocessableEntity, "value_labour_minutes must be >= 0")
		return
	}
	if req.MonopolySalePriceLabourMinutes < 0 {
		writeError(w, http.StatusUnprocessableEntity, "monopoly_sale_price_labour_minutes must be >= 0")
		return
	}
	monopolyRent := rent.ComputeMonopolyRent(req.MonopolySalePriceLabourMinutes, req.ValueLabourMinutes)
	rec := rent.MonopolyPriceRent{
		ParcelID:                       req.ParcelID,
		ProductKind:                    req.ProductKind,
		ValueLabourMinutes:             req.ValueLabourMinutes,
		MonopolySalePriceLabourMinutes: req.MonopolySalePriceLabourMinutes,
		MonopolyRentLabourMinutes:      monopolyRent,
	}
	saved, err := h.Store.CreateMonopolyPriceRent(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/monopoly-rents")
	writeJSON(w, http.StatusCreated, saved)
}

// ListMonopolyPriceRents handles GET /v1/rent/monopoly-rents.
// Returns {"items":[...]} — never null.
func (h *Handler) ListMonopolyPriceRents(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListMonopolyPriceRents(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.MonopolyPriceRent{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createLandPriceScenarioRequest is the body for POST /v1/rent/land-price-scenarios.
type createLandPriceScenarioRequest struct {
	Kind              int    `json:"kind"`
	ParcelID          string `json:"parcel_id"`
	OldAnnualRentBP   int64  `json:"old_annual_rent_bp"`
	NewAnnualRentBP   int64  `json:"new_annual_rent_bp"`
	OldInterestRateBP int64  `json:"old_interest_rate_bp"`
	NewInterestRateBP int64  `json:"new_interest_rate_bp"`
	OldLandPriceBP    int64  `json:"old_land_price_bp"`
}

// CreateLandPriceScenario handles POST /v1/rent/land-price-scenarios.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateLandPriceScenario(w http.ResponseWriter, r *http.Request) {
	var req createLandPriceScenarioRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !rent.LandPriceScenarioKind(req.Kind).IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "kind is not a valid LandPriceScenarioKind")
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.NewInterestRateBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "new_interest_rate_bp must be > 0")
		return
	}
	if req.OldInterestRateBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "old_interest_rate_bp must be >= 0")
		return
	}
	if req.OldAnnualRentBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "old_annual_rent_bp must be >= 0")
		return
	}
	if req.NewAnnualRentBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "new_annual_rent_bp must be >= 0")
		return
	}
	if req.OldLandPriceBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "old_land_price_bp must be >= 0")
		return
	}
	newLandPrice := rent.CapitalisedLandPrice(req.NewAnnualRentBP, req.NewInterestRateBP)
	rec := rent.LandPriceScenario{
		Kind:              rent.LandPriceScenarioKind(req.Kind),
		ParcelID:          req.ParcelID,
		OldAnnualRentBP:   req.OldAnnualRentBP,
		NewAnnualRentBP:   req.NewAnnualRentBP,
		OldInterestRateBP: req.OldInterestRateBP,
		NewInterestRateBP: req.NewInterestRateBP,
		OldLandPriceBP:    req.OldLandPriceBP,
		NewLandPriceBP:    newLandPrice,
	}
	saved, err := h.Store.CreateLandPriceScenario(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/land-price-scenarios")
	writeJSON(w, http.StatusCreated, saved)
}

// ListLandPriceScenarios handles GET /v1/rent/land-price-scenarios.
// Returns {"items":[...]} — never null.
func (h *Handler) ListLandPriceScenarios(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListLandPriceScenarios(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.LandPriceScenario{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
