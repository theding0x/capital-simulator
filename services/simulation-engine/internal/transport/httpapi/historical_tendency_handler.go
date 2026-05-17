package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Ch. 32 — The Historical Tendency of Capitalist Accumulation.

type negationResponse struct {
	Stage       string `json:"stage"`
	Description string `json:"description"`
}

type negationOfNegationResponse struct {
	Stages   []negationResponse `json:"stages"`
	Negation negationResponse   `json:"negation"`
}

type runCentralisationRequest struct {
	Name              string  `json:"name,omitempty"`
	Firms             int64   `json:"firms"`
	TotalCapitalPence int64   `json:"total_capital_pence"`
	WageLabourers     int64   `json:"wage_labourers"`
	Steps             int64   `json:"steps"`
	AbsorptionRate    float64 `json:"absorption_rate"`
	CapitalGrowthRate float64 `json:"capital_growth_rate"`
	Persist           bool    `json:"persist,omitempty"`
}

type centralisationStepResponse struct {
	StepIndex                int64 `json:"step_index"`
	FirmsAbsorbed            int64 `json:"firms_absorbed"`
	CapitalConcentratedPence int64 `json:"capital_concentrated_pence"`
}

type accumulationTrajectoryResponse struct {
	ID                  string                       `json:"id"`
	Name                string                       `json:"name"`
	InitialFirms        int64                        `json:"initial_firms"`
	InitialCapitalPence int64                        `json:"initial_capital_pence"`
	Steps               []centralisationStepResponse `json:"steps"`
	FinalFirms          int64                        `json:"final_firms"`
	FinalCapitalPence   int64                        `json:"final_capital_pence"`
	ReserveArmySize     int64                        `json:"reserve_army_size"`
	CreatedAt           string                       `json:"created_at,omitempty"`
}

func toCentralisationStepResponse(s simulation.CentralisationStep) centralisationStepResponse {
	return centralisationStepResponse{
		StepIndex:                s.StepIndex,
		FirmsAbsorbed:            s.FirmsAbsorbed,
		CapitalConcentratedPence: int64(s.CapitalConcentratedPence),
	}
}

func toAccumulationTrajectoryResponse(t simulation.AccumulationTrajectory) accumulationTrajectoryResponse {
	steps := make([]centralisationStepResponse, len(t.Steps))
	for i, s := range t.Steps {
		steps[i] = toCentralisationStepResponse(s)
	}
	resp := accumulationTrajectoryResponse{
		ID:                  string(t.ID),
		Name:                t.Name,
		InitialFirms:        t.InitialFirms,
		InitialCapitalPence: int64(t.InitialCapitalPence),
		Steps:               steps,
		FinalFirms:          t.FinalFirms,
		FinalCapitalPence:   int64(t.FinalCapitalPence),
		ReserveArmySize:     t.ReserveArmySize,
	}
	if !t.CreatedAt.IsZero() {
		resp.CreatedAt = t.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00")
	}
	return resp
}

// GetNegationOfNegation handles GET /v1/accumulation/negation-of-negation.
// Stateless: returns the three dialectical stages and the resolved
// third-stage negation.
func (h *Handler) GetNegationOfNegation(w http.ResponseWriter, _ *http.Request) {
	stages := simulation.DialecticalSequence()
	neg, err := simulation.NegationOfNegation(stages)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respStages := make([]negationResponse, len(stages))
	for i, s := range stages {
		respStages[i] = negationResponse{Stage: s.Stage, Description: s.Description}
	}
	writeJSON(w, http.StatusOK, negationOfNegationResponse{
		Stages:   respStages,
		Negation: negationResponse{Stage: neg.Stage, Description: neg.Description},
	})
}

// RunCentralisation handles POST /v1/accumulation/centralisation.
// Stateless compute: takes an initial CapitalistPrivateProperty plus
// simulation parameters and returns the resulting AccumulationTrajectory.
// If persist=true, the trajectory is also stored and the response carries
// its assigned id.
func (h *Handler) RunCentralisation(w http.ResponseWriter, r *http.Request) {
	var req runCentralisationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Firms <= 0 {
		writeError(w, http.StatusBadRequest, "firms must be positive")
		return
	}
	if req.TotalCapitalPence <= 0 {
		writeError(w, http.StatusBadRequest, "total_capital_pence must be positive")
		return
	}
	if req.WageLabourers < 0 {
		writeError(w, http.StatusBadRequest, "wage_labourers cannot be negative")
		return
	}
	if req.Steps <= 0 {
		writeError(w, http.StatusBadRequest, "steps must be positive")
		return
	}
	if req.AbsorptionRate <= 0 || req.AbsorptionRate > 1 {
		writeError(w, http.StatusBadRequest, "absorption_rate must be in (0, 1]")
		return
	}
	if req.CapitalGrowthRate < 0 {
		writeError(w, http.StatusBadRequest, "capital_growth_rate must be >= 0")
		return
	}

	initial := simulation.CapitalistPrivateProperty{
		Firms:             req.Firms,
		TotalCapitalPence: simulation.Pence(req.TotalCapitalPence),
		WageLabourers:     req.WageLabourers,
	}
	traj := simulation.RunCentralisation(initial, req.Steps, req.AbsorptionRate, req.CapitalGrowthRate)

	if req.Persist {
		if req.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required when persist=true")
			return
		}
		traj.Name = req.Name
		stored, err := h.Trajectories.CreateAccumulationTrajectory(r.Context(), traj)
		if err != nil {
			if errors.Is(err, store.ErrAlreadyExists) {
				writeError(w, http.StatusConflict, "trajectory with that name already exists")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Location", fmt.Sprintf("/v1/accumulation/trajectories/%s", string(stored.ID)))
		writeJSON(w, http.StatusCreated, toAccumulationTrajectoryResponse(stored))
		return
	}

	writeJSON(w, http.StatusOK, toAccumulationTrajectoryResponse(traj))
}

// ListAccumulationTrajectories handles GET /v1/accumulation/trajectories.
func (h *Handler) ListAccumulationTrajectories(w http.ResponseWriter, r *http.Request) {
	trajectories, err := h.Trajectories.ListAccumulationTrajectories(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]accumulationTrajectoryResponse, len(trajectories))
	for i, t := range trajectories {
		out[i] = toAccumulationTrajectoryResponse(t)
	}
	writeJSON(w, http.StatusOK, out)
}

// GetAccumulationTrajectory handles GET /v1/accumulation/trajectories/{id}.
func (h *Handler) GetAccumulationTrajectory(w http.ResponseWriter, r *http.Request) {
	id := simulation.AccumulationTrajectoryID(r.PathValue("id"))
	if id.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	traj, err := h.Trajectories.GetAccumulationTrajectory(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "accumulation trajectory not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toAccumulationTrajectoryResponse(traj))
}
