package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createLandParcelRequest is the body for POST /v1/rent/parcels.
type createLandParcelRequest struct {
	OwnerID        string `json:"owner_id"`
	FertilityGrade int    `json:"fertility_grade"`
	AreaHectares   int64  `json:"area_hectares"`
	Location       string `json:"location"`
}

// CreateLandParcel handles POST /v1/rent/parcels.
// Returns 400 on bad JSON, 422 on invalid field values, 201 + Location on success.
func (h *Handler) CreateLandParcel(w http.ResponseWriter, r *http.Request) {
	var req createLandParcelRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.OwnerID == "" {
		writeError(w, http.StatusUnprocessableEntity, "owner_id must not be empty")
		return
	}
	if req.FertilityGrade < 1 || req.FertilityGrade > 4 {
		writeError(w, http.StatusUnprocessableEntity, "fertility_grade must be between 1 and 4")
		return
	}
	if req.AreaHectares < 0 {
		writeError(w, http.StatusUnprocessableEntity, "area_hectares must be >= 0")
		return
	}
	p := rent.LandParcel{
		OwnerID:        req.OwnerID,
		FertilityGrade: req.FertilityGrade,
		AreaHectares:   req.AreaHectares,
		Location:       req.Location,
	}
	saved, err := h.Store.CreateLandParcel(r.Context(), p)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/parcels/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListLandParcels handles GET /v1/rent/parcels.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListLandParcels(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListLandParcels(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.LandParcel{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetLandParcel handles GET /v1/rent/parcels/{id}.
func (h *Handler) GetLandParcel(w http.ResponseWriter, r *http.Request) {
	id := rent.LandParcelID(r.PathValue("id"))
	p, err := h.Store.GetLandParcel(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// createGroundRentRequest is the body for POST /v1/rent/rents.
type createGroundRentRequest struct {
	ParcelID            string       `json:"parcel_id"`
	CapitalistID        string       `json:"capitalist_id"`
	LandOwnerID         string       `json:"land_owner_id"`
	Form                rent.RentForm `json:"form"`
	AmountLabourMinutes int64        `json:"amount_labour_minutes"`
	PeriodYears         int          `json:"period_years"`
}

// CreateGroundRent handles POST /v1/rent/rents.
// Returns 400 on bad JSON, 422 on domain-rule violations, 201 + Location on success.
func (h *Handler) CreateGroundRent(w http.ResponseWriter, r *http.Request) {
	var req createGroundRentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	g := rent.GroundRent{
		ParcelID:            req.ParcelID,
		CapitalistID:        req.CapitalistID,
		LandOwnerID:         req.LandOwnerID,
		Form:                req.Form,
		AmountLabourMinutes: req.AmountLabourMinutes,
		PeriodYears:         req.PeriodYears,
	}
	if err := g.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	saved, err := h.Store.CreateGroundRent(r.Context(), g)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/rents/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListGroundRents handles GET /v1/rent/rents.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListGroundRents(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListGroundRents(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.GroundRent{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
