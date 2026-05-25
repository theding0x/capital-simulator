// Vol. II Ch. 6 — HTTP handlers for costs of circulation.
package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/costs"
	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
)

// --- request/response types -------------------------------------------------

type createCirculationCostRequest struct {
	IndustrialCapitalID costs.IndustrialCapitalID `json:"industrial_capital_id"`
	Kind                costs.CirculationCostKind `json:"kind"`
	Pence               costs.Pence               `json:"pence"`
	IncurredAt          *time.Time                `json:"incurred_at,omitempty"`
}

type createCirculationAgentRequest struct {
	IndustrialCapitalID    costs.IndustrialCapitalID  `json:"industrial_capital_id"`
	AgentID                costs.AgentID              `json:"agent_id,omitempty"`
	Role                   costs.CirculationAgentRole `json:"role"`
	WageRatePence          costs.Pence                `json:"wage_rate_pence"`
	LabourMinutesNecessary int64                      `json:"labour_minutes_necessary"`
	LabourMinutesSurplus   int64                      `json:"labour_minutes_surplus"`
}

type createMoneyFauxFraisRequest struct {
	IndustrialCapitalID    costs.IndustrialCapitalID `json:"industrial_capital_id,omitempty"`
	Pence                  costs.Pence               `json:"pence"`
	AnnualReplacementPence costs.Pence               `json:"annual_replacement_pence"`
}

type createCommoditySupplyRequest struct {
	IndustrialCapitalID costs.IndustrialCapitalID `json:"industrial_capital_id"`
	CommodityID         costs.CommodityID         `json:"commodity_id,omitempty"`
	Pence               costs.Pence               `json:"pence"`
	IsVoluntary         bool                      `json:"is_voluntary"`
}

type createStorageCostRequest struct {
	BuildingPence           costs.Pence `json:"building_pence"`
	LabourPence             costs.Pence `json:"labour_pence"`
	PreservationLabourPence costs.Pence `json:"preservation_labour_pence"`
}

type setTransportTariffRequest struct {
	CommodityID                       costs.CommodityID `json:"commodity_id"`
	BasePencePerTonMile               costs.Pence       `json:"base_pence_per_ton_mile"`
	FragilityMultiplierBasisPoints    int64             `json:"fragility_multiplier_basis_points"`
	BreakageRiskMultiplierBasisPoints int64             `json:"breakage_risk_multiplier_basis_points"`
}

type createTransportLegRequest struct {
	CommodityID            costs.CommodityID `json:"commodity_id"`
	Quantity               int64             `json:"quantity"`
	OriginID               string            `json:"origin_id"`
	DestinationID          string            `json:"destination_id"`
	DistanceMeters         int64             `json:"distance_meters"`
	WeightGrams            int64             `json:"weight_grams"`
	VolumeCubicMillimeters int64             `json:"volume_cubic_millimeters"`
	LabourCostPence        costs.Pence       `json:"labour_cost_pence"`
	MeansOfTransportPence  costs.Pence       `json:"means_of_transport_pence"`
	SurplusPence           costs.Pence       `json:"surplus_pence"`
}

// --- CirculationCost handlers -----------------------------------------------

// CreateCirculationCost handles POST /v1/circulation-costs.
func (h *Handler) CreateCirculationCost(w http.ResponseWriter, r *http.Request) {
	var req createCirculationCostRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	at := time.Now().UTC()
	if req.IncurredAt != nil {
		at = *req.IncurredAt
	}
	c, err := costs.NewCirculationCost(req.IndustrialCapitalID, req.Kind, req.Pence, at)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.CreateCirculationCost(r.Context(), c)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Location", "/v1/circulation-costs/"+string(created.ID))
	writeJSON(w, http.StatusCreated, created)
}

// GetCirculationCost handles GET /v1/circulation-costs/{id}.
func (h *Handler) GetCirculationCost(w http.ResponseWriter, r *http.Request) {
	id := costs.CirculationCostID(r.PathValue("id"))
	c, err := h.Store.GetCirculationCost(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// ListCirculationCosts handles GET /v1/circulation-costs.
// Optional query params: industrial_capital_id, kind, nature.
func (h *Handler) ListCirculationCosts(w http.ResponseWriter, r *http.Request) {
	icID := costs.IndustrialCapitalID(r.URL.Query().Get("industrial_capital_id"))
	kind := costs.CirculationCostKind(r.URL.Query().Get("kind"))
	nature := costs.CostNature(r.URL.Query().Get("nature"))
	out, err := h.Store.ListCirculationCosts(r.Context(), icID, kind, nature)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out == nil {
		out = []costs.CirculationCost{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// AggregateCirculationCosts handles GET /v1/circulation-costs/aggregate.
func (h *Handler) AggregateCirculationCosts(w http.ResponseWriter, r *http.Request) {
	icID := costs.IndustrialCapitalID(r.URL.Query().Get("industrial_capital_id"))
	period := r.URL.Query().Get("period")
	res, err := h.Store.AggregateForCapital(r.Context(), icID, period)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// SystemFauxFrais handles GET /v1/circulation-costs/system-faux-frais.
func (h *Handler) SystemFauxFrais(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	res, err := h.Store.SystemFauxFrais(r.Context(), period)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// --- CirculationAgent handler -----------------------------------------------

// CreateCirculationAgent handles POST /v1/circulation-agents.
func (h *Handler) CreateCirculationAgent(w http.ResponseWriter, r *http.Request) {
	var req createCirculationAgentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a := costs.CirculationAgent{
		ID:                     costs.NewCirculationAgentID(),
		IndustrialCapitalID:    req.IndustrialCapitalID,
		AgentID:                req.AgentID,
		Role:                   req.Role,
		WageRatePence:          req.WageRatePence,
		LabourMinutesNecessary: req.LabourMinutesNecessary,
		LabourMinutesSurplus:   req.LabourMinutesSurplus,
	}
	created, err := h.Store.CreateCirculationAgent(r.Context(), a)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// --- MoneyAsFauxFrais handler -----------------------------------------------

// CreateMoneyAsFauxFrais handles POST /v1/money-as-faux-frais.
func (h *Handler) CreateMoneyAsFauxFrais(w http.ResponseWriter, r *http.Request) {
	var req createMoneyFauxFraisRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	mf := costs.MoneyAsFauxFrais{
		IndustrialCapitalID:    req.IndustrialCapitalID,
		Pence:                  req.Pence,
		AnnualReplacementPence: req.AnnualReplacementPence,
	}
	created, err := h.Store.CreateMoneyAsFauxFrais(r.Context(), mf)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// --- CommoditySupply handlers -----------------------------------------------

// CreateCommoditySupply handles POST /v1/commodity-supplies.
func (h *Handler) CreateCommoditySupply(w http.ResponseWriter, r *http.Request) {
	var req createCommoditySupplyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.IndustrialCapitalID == "" {
		writeError(w, http.StatusBadRequest, "industrial_capital_id is required")
		return
	}
	s := costs.CommoditySupply{
		IndustrialCapitalID: req.IndustrialCapitalID,
		CommodityID:         req.CommodityID,
		Pence:               req.Pence,
		HeldSince:           time.Now().UTC(),
		IsVoluntary:         req.IsVoluntary,
	}
	created, err := h.Store.CreateCommoditySupply(r.Context(), s)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Location", "/v1/commodity-supplies/"+string(created.ID))
	writeJSON(w, http.StatusCreated, created)
}

// AddStorageCost handles POST /v1/commodity-supplies/{id}/storage-cost.
func (h *Handler) AddStorageCost(w http.ResponseWriter, r *http.Request) {
	supplyID := costs.CommoditySupplyID(r.PathValue("id"))
	if _, err := h.Store.GetCommoditySupply(r.Context(), supplyID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var req createStorageCostRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sc := costs.StorageCost{
		CommoditySupplyID:       supplyID,
		BuildingPence:           req.BuildingPence,
		LabourPence:             req.LabourPence,
		PreservationLabourPence: req.PreservationLabourPence,
	}
	created, err := h.Store.CreateStorageCost(r.Context(), sc)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// --- TransportTariff handler ------------------------------------------------

// SetTransportTariff handles POST /v1/transport-tariffs.
func (h *Handler) SetTransportTariff(w http.ResponseWriter, r *http.Request) {
	var req setTransportTariffRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	t := costs.TransportTariff{
		CommodityID:                       req.CommodityID,
		BasePencePerTonMile:               req.BasePencePerTonMile,
		FragilityMultiplierBasisPoints:    req.FragilityMultiplierBasisPoints,
		BreakageRiskMultiplierBasisPoints: req.BreakageRiskMultiplierBasisPoints,
	}
	created, err := h.Store.SetTransportTariff(r.Context(), t)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// --- TransportLeg handler ---------------------------------------------------

// CreateTransportLeg handles POST /v1/transport-legs.
func (h *Handler) CreateTransportLeg(w http.ResponseWriter, r *http.Request) {
	var req createTransportLegRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	leg := costs.TransportLeg{
		CommodityID:            req.CommodityID,
		Quantity:               req.Quantity,
		OriginID:               req.OriginID,
		DestinationID:          req.DestinationID,
		DistanceMeters:         req.DistanceMeters,
		WeightGrams:            req.WeightGrams,
		VolumeCubicMillimeters: req.VolumeCubicMillimeters,
		LabourCostPence:        req.LabourCostPence,
		MeansOfTransportPence:  req.MeansOfTransportPence,
		SurplusPence:           req.SurplusPence,
	}
	created, err := h.Store.CreateTransportLeg(r.Context(), leg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Location", "/v1/transport-legs/"+string(created.ID))
	writeJSON(w, http.StatusCreated, created)
}
