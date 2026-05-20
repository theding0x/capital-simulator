package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Ch. 30 — Reaction of the Agricultural Revolution on Industry.

type createDomesticIndustryRequest struct {
	Name              string `json:"name"`
	HouseholdsEngaged int64  `json:"households_engaged"`
	AnnualOutputPence int64  `json:"annual_output_pence"`
}

type domesticIndustryResponse struct {
	ID                string `json:"id"`
	HistoricalStageID string `json:"historical_stage_id"`
	Name              string `json:"name"`
	HouseholdsEngaged int64  `json:"households_engaged"`
	AnnualOutputPence int64  `json:"annual_output_pence"`
	CreatedAt         string `json:"created_at"`
}

type marketFormationRequest struct {
	Name              string `json:"name"`
	HouseholdsEngaged int64  `json:"households_engaged"`
	AnnualOutputPence int64  `json:"annual_output_pence"`
	Expropriated      int64  `json:"expropriated"`
}

type marketFormationResponse struct {
	Period                    string `json:"period"`
	LabourersProletarianised  int64  `json:"labourers_proletarianised"`
	DomesticIndustryDestroyed bool   `json:"domestic_industry_destroyed"`
	NewWageLabourers          int64  `json:"new_wage_labourers"`
	MarketSizePence           int64  `json:"market_size_pence"`
}

type homeMarketSizeResponse struct {
	CommoditisedOutputPence int64 `json:"commoditised_output_pence"`
	WageLabourers           int64 `json:"wage_labourers"`
}

func toDomesticIndustryResponse(d simulation.DomesticIndustry) domesticIndustryResponse {
	return domesticIndustryResponse{
		ID:                string(d.ID),
		HistoricalStageID: string(d.HistoricalStageID),
		Name:              d.Name,
		HouseholdsEngaged: d.HouseholdsEngaged,
		AnnualOutputPence: int64(d.AnnualOutputPence),
		CreatedAt:         d.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
	}
}

// CreateDomesticIndustry handles POST /v1/historical-stages/{id}/domestic-industries.
func (h *Handler) CreateDomesticIndustry(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))
	if stageID.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req createDomesticIndustryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	d := simulation.DomesticIndustry{
		HistoricalStageID: stageID,
		Name:              req.Name,
		HouseholdsEngaged: req.HouseholdsEngaged,
		AnnualOutputPence: simulation.Pence(req.AnnualOutputPence),
	}
	if err := d.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.DomesticIndustries.CreateDomesticIndustry(r.Context(), d)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "domestic industry already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/v1/historical-stages/%s/domestic-industries", string(stageID)))
	writeJSON(w, http.StatusCreated, toDomesticIndustryResponse(created))
}

// ComputeMarketFormation handles POST /v1/market-formation.
// Stateless: derives a HomeMarketFormation from an inline domestic industry definition
// and an expropriated count without persisting anything.
func (h *Handler) ComputeMarketFormation(w http.ResponseWriter, r *http.Request) {
	var req marketFormationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	di := simulation.DomesticIndustry{
		Name:              req.Name,
		HouseholdsEngaged: req.HouseholdsEngaged,
		AnnualOutputPence: simulation.Pence(req.AnnualOutputPence),
	}
	if err := di.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Expropriated < 0 {
		writeError(w, http.StatusBadRequest, "expropriated cannot be negative")
		return
	}

	formation := simulation.ComputeMarketFormation(di, req.Expropriated)
	writeJSON(w, http.StatusOK, marketFormationResponse{
		Period:                    formation.Period,
		LabourersProletarianised:  formation.LabourersProletarianised,
		DomesticIndustryDestroyed: formation.DomesticIndustryDestroyed,
		NewWageLabourers:          formation.NewWageLabourers,
		MarketSizePence:           int64(formation.MarketSizePence),
	})
}

// GetHomeMarket handles GET /v1/historical-stages/{id}/home-market.
// Returns the aggregate home market size derived from all domestic industries
// registered for the given historical stage.
func (h *Handler) GetHomeMarket(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))
	if stageID.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	_, err := h.HistoricalStages.GetHistoricalStage(r.Context(), stageID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "historical stage not found")
			return
		}
		h.writeServerError(w, err)
		return
	}

	industries, err := h.DomesticIndustries.ListDomesticIndustriesByStage(r.Context(), stageID)
	if err != nil {
		h.writeServerError(w, err)
		return
	}

	size := simulation.AggregateHomeMarket(industries)
	writeJSON(w, http.StatusOK, homeMarketSizeResponse{
		CommoditisedOutputPence: int64(size.CommoditisedOutputPence),
		WageLabourers:           size.WageLabourers,
	})
}
