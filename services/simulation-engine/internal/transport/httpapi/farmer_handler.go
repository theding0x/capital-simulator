package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Ch. 29 — Genesis of the Capitalist Farmer.

type createFarmTenureRequest struct {
	Form                 string  `json:"form"`
	LeasePeriodYears     int64   `json:"lease_period_years"`
	RentPence            int64   `json:"rent_pence"`
	CapitalAdvancedPence int64   `json:"capital_advanced_pence"`
	RevenuePence         int64   `json:"revenue_pence"`
	WageCostsPence       int64   `json:"wage_costs_pence"`
}

type farmingSurplusDTO struct {
	Revenue     int64 `json:"revenue"`
	NominalRent int64 `json:"nominal_rent"`
	WageCosts   int64 `json:"wage_costs"`
	Profit      int64 `json:"profit"`
}

type farmTenureResponse struct {
	ID                   string            `json:"id"`
	HistoricalStageID    string            `json:"historical_stage_id"`
	Form                 string            `json:"form"`
	LeasePeriodYears     int64             `json:"lease_period_years"`
	RentPence            int64             `json:"rent_pence"`
	CapitalAdvancedPence int64             `json:"capital_advanced_pence"`
	RevenuePence         int64             `json:"revenue_pence"`
	WageCostsPence       int64             `json:"wage_costs_pence"`
	Surplus              farmingSurplusDTO `json:"surplus"`
	CreatedAt            string            `json:"created_at"`
}

type realRentRequest struct {
	NominalRentPence   int64   `json:"nominal_rent_pence"`
	DepreciationFactor float64 `json:"depreciation_factor"`
}

type realRentResponse struct {
	NominalRentPence   int64   `json:"nominal_rent_pence"`
	RealRentPence      int64   `json:"real_rent_pence"`
	DepreciationFactor float64 `json:"depreciation_factor"`
}

func toFarmTenureResponse(f simulation.FarmTenure) farmTenureResponse {
	surplus := simulation.ComputeFarmingSurplus(f.RevenuePence, f.RentPence, f.WageCostsPence)
	return farmTenureResponse{
		ID:                   string(f.ID),
		HistoricalStageID:    string(f.HistoricalStageID),
		Form:                 string(f.Form),
		LeasePeriodYears:     f.LeasePeriodYears,
		RentPence:            int64(f.RentPence),
		CapitalAdvancedPence: int64(f.CapitalAdvancedPence),
		RevenuePence:         int64(f.RevenuePence),
		WageCostsPence:       int64(f.WageCostsPence),
		Surplus: farmingSurplusDTO{
			Revenue:     int64(surplus.Revenue),
			NominalRent: int64(surplus.NominalRent),
			WageCosts:   int64(surplus.WageCosts),
			Profit:      int64(surplus.Profit),
		},
		CreatedAt: f.CreatedAt.Format(time.RFC3339Nano),
	}
}

// CreateFarmTenure handles POST /v1/historical-stages/{id}/farm-tenures.
func (h *Handler) CreateFarmTenure(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))

	_, err := h.HistoricalStages.GetHistoricalStage(r.Context(), stageID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "historical stage not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var req createFarmTenureRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	ft := simulation.FarmTenure{
		HistoricalStageID:    stageID,
		Form:                 simulation.TenantForm(req.Form),
		LeasePeriodYears:     req.LeasePeriodYears,
		RentPence:            simulation.Pence(req.RentPence),
		CapitalAdvancedPence: simulation.Pence(req.CapitalAdvancedPence),
		RevenuePence:         simulation.Pence(req.RevenuePence),
		WageCostsPence:       simulation.Pence(req.WageCostsPence),
	}
	if err := ft.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.FarmTenures.CreateFarmTenure(r.Context(), ft)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/v1/historical-stages/%s/farm-tenures/%s", string(stageID), string(created.ID)))
	writeJSON(w, http.StatusCreated, toFarmTenureResponse(created))
}

// ListFarmTenures handles GET /v1/historical-stages/{id}/farm-tenures.
// Returns all tenures for the stage with computed surplus inline.
func (h *Handler) ListFarmTenures(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))

	_, err := h.HistoricalStages.GetHistoricalStage(r.Context(), stageID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "historical stage not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tenures, err := h.FarmTenures.ListFarmTenuresByStage(r.Context(), stageID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	out := make([]farmTenureResponse, 0, len(tenures))
	for _, ft := range tenures {
		out = append(out, toFarmTenureResponse(ft))
	}
	writeJSON(w, http.StatusOK, out)
}

// ComputeRealRent handles POST /v1/farm-tenures/real-rent.
// Stateless: deflates a nominal rent by the given depreciation factor.
func (h *Handler) ComputeRealRent(w http.ResponseWriter, r *http.Request) {
	var req realRentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.NominalRentPence < 0 {
		writeError(w, http.StatusBadRequest, "nominal_rent_pence cannot be negative")
		return
	}
	if req.DepreciationFactor <= 0 || req.DepreciationFactor > 1 {
		writeError(w, http.StatusBadRequest, "depreciation_factor must be in (0, 1]")
		return
	}

	md := simulation.MoneyDepreciation{DepreciationFactor: req.DepreciationFactor}
	realRent := simulation.ComputeRealRent(simulation.Pence(req.NominalRentPence), md)

	writeJSON(w, http.StatusOK, realRentResponse{
		NominalRentPence:   req.NominalRentPence,
		RealRentPence:      int64(realRent),
		DepreciationFactor: req.DepreciationFactor,
	})
}
