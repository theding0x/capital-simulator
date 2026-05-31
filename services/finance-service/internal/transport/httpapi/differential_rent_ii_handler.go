package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createSuccessiveInvestmentRequest is the body for POST /v1/rent/successive-investments.
type createSuccessiveInvestmentRequest struct {
	ParcelID        string `json:"parcel_id"`
	TrancheNumber   int    `json:"tranche_number"`
	CapitalBP       int64  `json:"capital_bp"`
	OutputQuarters  int64  `json:"output_quarters"`
	SurplusProfitBP int64  `json:"surplus_profit_bp"`
}

// CreateSuccessiveInvestment handles POST /v1/rent/successive-investments.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateSuccessiveInvestment(w http.ResponseWriter, r *http.Request) {
	var req createSuccessiveInvestmentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.TrancheNumber < 1 {
		writeError(w, http.StatusUnprocessableEntity, "tranche_number must be >= 1")
		return
	}
	if req.CapitalBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "capital_bp must be > 0")
		return
	}
	if req.OutputQuarters <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "output_quarters must be > 0")
		return
	}
	if req.SurplusProfitBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "surplus_profit_bp must be >= 0")
		return
	}
	inv := rent.SuccessiveInvestment{
		ParcelID:        req.ParcelID,
		TrancheNumber:   req.TrancheNumber,
		CapitalBP:       req.CapitalBP,
		OutputQuarters:  req.OutputQuarters,
		SurplusProfitBP: req.SurplusProfitBP,
	}
	saved, err := h.Store.CreateSuccessiveInvestment(r.Context(), inv)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/successive-investments/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListSuccessiveInvestments handles GET /v1/rent/successive-investments.
// Returns {"items":[...]} — never null.
func (h *Handler) ListSuccessiveInvestments(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListSuccessiveInvestments(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.SuccessiveInvestment{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetDR2Table handles GET /v1/rent/dr2-tables/{parcel_id}.
// Fetches all successive investments for the parcel, computes the DR II aggregate,
// and returns {"table":..., "investments":[...]}. Unknown parcel yields a zero table
// with an empty investments array — no 404.
func (h *Handler) GetDR2Table(w http.ResponseWriter, r *http.Request) {
	parcelID := r.PathValue("parcel_id")
	invs, err := h.Store.ListSuccessiveInvestmentsByParcel(r.Context(), parcelID)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if invs == nil {
		invs = []rent.SuccessiveInvestment{}
	}
	tbl := rent.ComputeDifferentialRentII(parcelID, invs, 1)
	tbl.ID = rent.NewDifferentialRentIITableID()
	writeJSON(w, http.StatusOK, map[string]any{
		"table":       tbl,
		"investments": invs,
	})
}

// createLeaseTermRequest is the body for POST /v1/rent/lease-terms.
type createLeaseTermRequest struct {
	ParcelID      string `json:"parcel_id"`
	StartYear     int    `json:"start_year"`
	DurationYears int    `json:"duration_years"`
	FixedRentBP   int64  `json:"fixed_rent_bp"`
}

// CreateLeaseTerm handles POST /v1/rent/lease-terms.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateLeaseTerm(w http.ResponseWriter, r *http.Request) {
	var req createLeaseTermRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.StartYear <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "start_year must be > 0")
		return
	}
	if req.DurationYears <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "duration_years must be > 0")
		return
	}
	if req.FixedRentBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "fixed_rent_bp must be >= 0")
		return
	}
	lt := rent.LeaseTerm{
		ParcelID:      req.ParcelID,
		StartYear:     req.StartYear,
		DurationYears: req.DurationYears,
		FixedRentBP:   req.FixedRentBP,
	}
	saved, err := h.Store.CreateLeaseTerm(r.Context(), lt)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/lease-terms/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListLeaseTerms handles GET /v1/rent/lease-terms.
// Returns {"items":[...]} — never null.
func (h *Handler) ListLeaseTerms(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListLeaseTerms(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.LeaseTerm{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
