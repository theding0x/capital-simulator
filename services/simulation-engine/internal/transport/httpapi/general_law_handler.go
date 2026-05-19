package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Ch. 25 — The General Law of Capitalist Accumulation.

// --- organic composition ---

type organicCompositionRequest struct {
	ConstantCapital int64 `json:"constant_capital"`
	VariableCapital int64 `json:"variable_capital"`
}

type organicCompositionResponse struct {
	ConstantCapital int64   `json:"constant_capital"`
	VariableCapital int64   `json:"variable_capital"`
	Ratio           float64 `json:"ratio"`
}

// ComputeOrganicComposition handles POST /v1/accumulation/organic-composition.
func (h *Handler) ComputeOrganicComposition(w http.ResponseWriter, r *http.Request) {
	var req organicCompositionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ConstantCapital < 0 {
		writeError(w, http.StatusBadRequest, "constant_capital cannot be negative")
		return
	}
	if req.VariableCapital <= 0 {
		writeError(w, http.StatusBadRequest, "variable_capital must be positive")
		return
	}

	vc := simulation.ValueComposition{
		ConstantCapital: simulation.Pence(req.ConstantCapital),
		VariableCapital: simulation.Pence(req.VariableCapital),
	}
	oc := simulation.ComputeOrganicComposition(vc)

	writeJSON(w, http.StatusOK, organicCompositionResponse{
		ConstantCapital: req.ConstantCapital,
		VariableCapital: req.VariableCapital,
		Ratio:           oc.Ratio,
	})
}

// --- labour demand ---

type labourDemandRequest struct {
	TotalCapital            int64   `json:"total_capital"`
	OrganicCompositionRatio float64 `json:"organic_composition_ratio"`
	WagePence               int64   `json:"wage_pence"`
}

type labourDemandResponse struct {
	TotalCapital            int64   `json:"total_capital"`
	OrganicCompositionRatio float64 `json:"organic_composition_ratio"`
	WagePence               int64   `json:"wage_pence"`
	Workers                 int64   `json:"workers"`
}

// ComputeLabourDemandEndpoint handles POST /v1/accumulation/labour-demand.
func (h *Handler) ComputeLabourDemandEndpoint(w http.ResponseWriter, r *http.Request) {
	var req labourDemandRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.TotalCapital <= 0 {
		writeError(w, http.StatusBadRequest, "total_capital must be positive")
		return
	}
	if req.OrganicCompositionRatio < 0 || req.OrganicCompositionRatio >= 1 {
		writeError(w, http.StatusBadRequest, "organic_composition_ratio must be in [0, 1)")
		return
	}
	if req.WagePence <= 0 {
		writeError(w, http.StatusBadRequest, "wage_pence must be positive")
		return
	}

	oc := simulation.OrganicComposition{Ratio: req.OrganicCompositionRatio}
	demand := simulation.ComputeLabourDemand(simulation.Pence(req.TotalCapital), oc, req.WagePence)

	writeJSON(w, http.StatusOK, labourDemandResponse{
		TotalCapital:            req.TotalCapital,
		OrganicCompositionRatio: req.OrganicCompositionRatio,
		WagePence:               req.WagePence,
		Workers:                 demand.Workers,
	})
}

// --- reserve army ---

type reserveArmyRequest struct {
	WorkerSupply    int64 `json:"worker_supply"`
	WorkersDemanded int64 `json:"workers_demanded"`
}

type reserveArmyResponse struct {
	WorkerSupply       int64   `json:"worker_supply"`
	WorkersDemanded    int64   `json:"workers_demanded"`
	ReserveArmySize    int64   `json:"reserve_army_size"`
	RelativeProportion float64 `json:"relative_proportion"`
}

// ComputeReserveArmyEndpoint handles POST /v1/accumulation/reserve-army.
func (h *Handler) ComputeReserveArmyEndpoint(w http.ResponseWriter, r *http.Request) {
	var req reserveArmyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.WorkerSupply < 0 {
		writeError(w, http.StatusBadRequest, "worker_supply cannot be negative")
		return
	}
	if req.WorkersDemanded < 0 {
		writeError(w, http.StatusBadRequest, "workers_demanded cannot be negative")
		return
	}

	demand := simulation.LabourDemand{Workers: req.WorkersDemanded}
	army := simulation.ComputeReserveArmy(req.WorkerSupply, demand)

	writeJSON(w, http.StatusOK, reserveArmyResponse{
		WorkerSupply:       req.WorkerSupply,
		WorkersDemanded:    req.WorkersDemanded,
		ReserveArmySize:    army.Size,
		RelativeProportion: army.RelativeProportion,
	})
}

// --- scenarios ---

type createGeneralLawScenarioRequest struct {
	Name               string  `json:"name"`
	ConstantCapital    int64   `json:"constant_capital"`
	VariableCapital    int64   `json:"variable_capital"`
	SurplusRate        float64 `json:"surplus_rate"`
	AccumulationRate   float64 `json:"accumulation_rate"`
	ProductivityGrowth float64 `json:"productivity_growth"`
	WagePence          int64   `json:"wage_pence"`
	WorkerSupply       int64   `json:"worker_supply"`
	Periods            int64   `json:"periods"`
}

type generalLawSnapshotResponse struct {
	Period             int64   `json:"period"`
	ConstantCapital    int64   `json:"constant_capital"`
	VariableCapital    int64   `json:"variable_capital"`
	OrganicComposition float64 `json:"organic_composition"`
	Workers            int64   `json:"workers"`
	ReserveArmySize    int64   `json:"reserve_army_size"`
	RelativeProportion float64 `json:"relative_proportion"`
}

type generalLawScenarioResponse struct {
	ID                 string                       `json:"id"`
	Name               string                       `json:"name"`
	ConstantCapital    int64                        `json:"constant_capital"`
	VariableCapital    int64                        `json:"variable_capital"`
	SurplusRate        float64                      `json:"surplus_rate"`
	AccumulationRate   float64                      `json:"accumulation_rate"`
	ProductivityGrowth float64                      `json:"productivity_growth"`
	WagePence          int64                        `json:"wage_pence"`
	WorkerSupply       int64                        `json:"worker_supply"`
	Periods            int64                        `json:"periods"`
	CreatedAt          string                       `json:"created_at"`
	Series             []generalLawSnapshotResponse `json:"series"`
}

func toGeneralLawScenarioResponse(s simulation.GeneralLawScenario) generalLawScenarioResponse {
	series := simulation.RunGeneralLaw(s)
	snaps := make([]generalLawSnapshotResponse, len(series))
	for i, sn := range series {
		snaps[i] = generalLawSnapshotResponse{
			Period:             sn.Period,
			ConstantCapital:    int64(sn.ConstantCapital),
			VariableCapital:    int64(sn.VariableCapital),
			OrganicComposition: sn.OrganicComposition,
			Workers:            sn.Workers,
			ReserveArmySize:    sn.ReserveArmySize,
			RelativeProportion: sn.RelativeProportion,
		}
	}
	return generalLawScenarioResponse{
		ID:                 string(s.ID),
		Name:               s.Name,
		ConstantCapital:    int64(s.ConstantCapital),
		VariableCapital:    int64(s.VariableCapital),
		SurplusRate:        s.SurplusRate,
		AccumulationRate:   s.AccumulationRate,
		ProductivityGrowth: s.ProductivityGrowth,
		WagePence:          s.WagePence,
		WorkerSupply:       s.WorkerSupply,
		Periods:            s.Periods,
		CreatedAt:          s.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
		Series:             snaps,
	}
}

// CreateGeneralLawScenario handles POST /v1/accumulation/scenarios.
func (h *Handler) CreateGeneralLawScenario(w http.ResponseWriter, r *http.Request) {
	var req createGeneralLawScenarioRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.ConstantCapital < 0 {
		writeError(w, http.StatusBadRequest, "constant_capital cannot be negative")
		return
	}
	if req.VariableCapital <= 0 {
		writeError(w, http.StatusBadRequest, "variable_capital must be positive")
		return
	}
	if req.SurplusRate < 0 {
		writeError(w, http.StatusBadRequest, "surplus_rate cannot be negative")
		return
	}
	if req.AccumulationRate < 0 || req.AccumulationRate > 1 {
		writeError(w, http.StatusBadRequest, "accumulation_rate must be between 0 and 1")
		return
	}
	if req.ProductivityGrowth < 0 {
		writeError(w, http.StatusBadRequest, "productivity_growth cannot be negative")
		return
	}
	if req.WagePence <= 0 {
		writeError(w, http.StatusBadRequest, "wage_pence must be positive")
		return
	}
	if req.WorkerSupply <= 0 {
		writeError(w, http.StatusBadRequest, "worker_supply must be positive")
		return
	}
	if req.Periods <= 0 {
		writeError(w, http.StatusBadRequest, "periods must be positive")
		return
	}

	s := simulation.GeneralLawScenario{
		Name:               req.Name,
		ConstantCapital:    simulation.Pence(req.ConstantCapital),
		VariableCapital:    simulation.Pence(req.VariableCapital),
		SurplusRate:        req.SurplusRate,
		AccumulationRate:   req.AccumulationRate,
		ProductivityGrowth: req.ProductivityGrowth,
		WagePence:          req.WagePence,
		WorkerSupply:       req.WorkerSupply,
		Periods:            req.Periods,
	}
	created, err := h.GeneralLaw.CreateGeneralLawScenario(r.Context(), s)
	if err != nil {
		h.writeServerError(w, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/v1/accumulation/scenarios/%s", string(created.ID)))
	writeJSON(w, http.StatusCreated, toGeneralLawScenarioResponse(created))
}

// GetGeneralLawScenario handles GET /v1/accumulation/scenarios/{id}.
func (h *Handler) GetGeneralLawScenario(w http.ResponseWriter, r *http.Request) {
	id := simulation.GeneralLawScenarioID(r.PathValue("id"))
	if id.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	s, err := h.GeneralLaw.GetGeneralLawScenario(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "scenario not found")
			return
		}
		h.writeServerError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toGeneralLawScenarioResponse(s))
}
