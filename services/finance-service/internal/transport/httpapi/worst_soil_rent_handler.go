package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createRentOnWorstSoilRequest is the body for POST /v1/rent/worst-soil-rents.
type createRentOnWorstSoilRequest struct {
	ParcelID               string `json:"parcel_id"`
	Mechanism              int    `json:"mechanism"`
	OldPriceOfProductionBP int64  `json:"old_price_of_production_bp"`
	NewPriceOfProductionBP int64  `json:"new_price_of_production_bp"`
}

// CreateRentOnWorstSoil handles POST /v1/rent/worst-soil-rents.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateRentOnWorstSoil(w http.ResponseWriter, r *http.Request) {
	var req createRentOnWorstSoilRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if !rent.RentEmergenceMechanism(req.Mechanism).IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "mechanism must be 1–3")
		return
	}
	if req.OldPriceOfProductionBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "old_price_of_production_bp must be > 0")
		return
	}
	if req.NewPriceOfProductionBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "new_price_of_production_bp must be > 0")
		return
	}
	rentPerAcre := rent.ComputeWorstSoilRentPerAcre(req.OldPriceOfProductionBP, req.NewPriceOfProductionBP)
	rec := rent.RentOnWorstSoil{
		ParcelID:               req.ParcelID,
		Mechanism:              rent.RentEmergenceMechanism(req.Mechanism),
		OldPriceOfProductionBP: req.OldPriceOfProductionBP,
		NewPriceOfProductionBP: req.NewPriceOfProductionBP,
		RentPerAcreBP:          rentPerAcre,
	}
	saved, err := h.Store.CreateRentOnWorstSoil(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/worst-soil-rents")
	writeJSON(w, http.StatusCreated, saved)
}

// ListRentOnWorstSoils handles GET /v1/rent/worst-soil-rents.
// Returns {"items":[...]} — never null.
func (h *Handler) ListRentOnWorstSoils(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListRentOnWorstSoils(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.RentOnWorstSoil{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createSoilImprovementCapitalRequest is the body for POST /v1/rent/soil-improvements.
type createSoilImprovementCapitalRequest struct {
	ParcelID           string `json:"parcel_id"`
	CapitalInvestedBP  int64  `json:"capital_invested_bp"`
	ProductivityGainBP int64  `json:"productivity_gain_bp"`
	AmortisationYears  int    `json:"amortisation_years"`
	AmortisedYears     int    `json:"amortised_years"`
}

// CreateSoilImprovementCapital handles POST /v1/rent/soil-improvements.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateSoilImprovementCapital(w http.ResponseWriter, r *http.Request) {
	var req createSoilImprovementCapitalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.CapitalInvestedBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "capital_invested_bp must be > 0")
		return
	}
	if req.ProductivityGainBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "productivity_gain_bp must be >= 0")
		return
	}
	if req.AmortisationYears < 1 {
		writeError(w, http.StatusUnprocessableEntity, "amortisation_years must be >= 1")
		return
	}
	if req.AmortisedYears < 0 {
		writeError(w, http.StatusUnprocessableEntity, "amortised_years must be >= 0")
		return
	}
	rec := rent.SoilImprovementCapital{
		ParcelID:           req.ParcelID,
		CapitalInvestedBP:  req.CapitalInvestedBP,
		ProductivityGainBP: req.ProductivityGainBP,
		AmortisationYears:  req.AmortisationYears,
		AmortisedYears:     req.AmortisedYears,
	}
	saved, err := h.Store.CreateSoilImprovementCapital(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/soil-improvements")
	writeJSON(w, http.StatusCreated, saved)
}

// ListSoilImprovementCapitals handles GET /v1/rent/soil-improvements.
// Returns {"items":[...]} — never null.
func (h *Handler) ListSoilImprovementCapitals(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListSoilImprovementCapitals(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.SoilImprovementCapital{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createRegulatingPriceShiftRequest is the body for POST /v1/rent/regulating-price-shifts.
type createRegulatingPriceShiftRequest struct {
	FromGrade  int   `json:"from_grade"`
	ToGrade    int   `json:"to_grade"`
	OldPriceBP int64 `json:"old_price_bp"`
	NewPriceBP int64 `json:"new_price_bp"`
}

// CreateRegulatingPriceShift handles POST /v1/rent/regulating-price-shifts.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateRegulatingPriceShift(w http.ResponseWriter, r *http.Request) {
	var req createRegulatingPriceShiftRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !rent.SoilGrade(req.FromGrade).IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "from_grade must be 1–4")
		return
	}
	if !rent.SoilGrade(req.ToGrade).IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "to_grade must be 1–4")
		return
	}
	if req.OldPriceBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "old_price_bp must be > 0")
		return
	}
	if req.NewPriceBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "new_price_bp must be > 0")
		return
	}
	rec := rent.RegulatingPriceShift{
		FromGrade:  rent.SoilGrade(req.FromGrade),
		ToGrade:    rent.SoilGrade(req.ToGrade),
		OldPriceBP: req.OldPriceBP,
		NewPriceBP: req.NewPriceBP,
	}
	saved, err := h.Store.CreateRegulatingPriceShift(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/regulating-price-shifts")
	writeJSON(w, http.StatusCreated, saved)
}

// ListRegulatingPriceShifts handles GET /v1/rent/regulating-price-shifts.
// Returns {"items":[...]} — never null.
func (h *Handler) ListRegulatingPriceShifts(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListRegulatingPriceShifts(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.RegulatingPriceShift{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
