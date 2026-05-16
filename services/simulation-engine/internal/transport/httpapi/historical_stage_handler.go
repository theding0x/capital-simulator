package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Ch. 26 — The Secret of Primitive Accumulation.

type primitiveAccumulationDTO struct {
	Period                string `json:"period"`
	Method                string `json:"method"`
	LabourersExpropriated int64  `json:"labourers_expropriated"`
	CapitalFormed         int64  `json:"capital_formed"`
}

type createHistoricalStageRequest struct {
	Name                   string                     `json:"name"`
	Description            string                     `json:"description"`
	PrimitiveAccumulations []primitiveAccumulationDTO `json:"primitive_accumulations"`
}

type historicalStageResponse struct {
	ID                         string                     `json:"id"`
	Name                       string                     `json:"name"`
	Description                string                     `json:"description"`
	PrimitiveAccumulations     []primitiveAccumulationDTO `json:"primitive_accumulations"`
	TotalCapitalFormed         int64                      `json:"total_capital_formed"`
	TotalLabourersExpropriated int64                      `json:"total_labourers_expropriated"`
	CreatedAt                  string                     `json:"created_at"`
}

type producerSeparationDTO struct {
	PreCapitalistWorkers int64 `json:"pre_capitalist_workers"`
	DisplacedWorkers     int64 `json:"displaced_workers"`
	FreeProletarians     int64 `json:"free_proletarians"`
}

type seedScenarioRequest struct {
	PreCapitalistWorkers int64   `json:"pre_capitalist_workers"`
	CompositionRatio     float64 `json:"composition_ratio"`
	SurplusRate          float64 `json:"surplus_rate"`
	Periods              int64   `json:"periods"`
}

type seededReproductionCycleDTO struct {
	Period      int64                     `json:"period"`
	Capital     seededCapitalStockDTO     `json:"capital"`
	SurplusRate float64                   `json:"surplus_rate"`
	Fund        seededSurplusValueFundDTO `json:"fund"`
}

type seededCapitalStockDTO struct {
	ConstantCapital int64 `json:"constant_capital"`
	VariableCapital int64 `json:"variable_capital"`
}

type seededSurplusValueFundDTO struct {
	Total       int64 `json:"total"`
	Revenue     int64 `json:"revenue"`
	Accumulated int64 `json:"accumulated"`
}

type seedScenarioResponse struct {
	StageID       string                       `json:"stage_id"`
	StageName     string                       `json:"stage_name"`
	SeededCapital seededCapitalStockDTO        `json:"seeded_capital"`
	Separation    producerSeparationDTO        `json:"separation"`
	SurplusRate   float64                      `json:"surplus_rate"`
	Periods       int64                        `json:"periods"`
	Cycles        []seededReproductionCycleDTO `json:"cycles"`
}

func episodesFromDTO(in []primitiveAccumulationDTO) []simulation.PrimitiveAccumulation {
	out := make([]simulation.PrimitiveAccumulation, len(in))
	for i, p := range in {
		out[i] = simulation.PrimitiveAccumulation{
			Period:                p.Period,
			Method:                p.Method,
			LabourersExpropriated: p.LabourersExpropriated,
			CapitalFormed:         simulation.Pence(p.CapitalFormed),
		}
	}
	return out
}

func episodesToDTO(in []simulation.PrimitiveAccumulation) []primitiveAccumulationDTO {
	out := make([]primitiveAccumulationDTO, len(in))
	for i, p := range in {
		out[i] = primitiveAccumulationDTO{
			Period:                p.Period,
			Method:                p.Method,
			LabourersExpropriated: p.LabourersExpropriated,
			CapitalFormed:         int64(p.CapitalFormed),
		}
	}
	return out
}

func toHistoricalStageResponse(h simulation.HistoricalStage) historicalStageResponse {
	return historicalStageResponse{
		ID:                         string(h.ID),
		Name:                       h.Name,
		Description:                h.Description,
		PrimitiveAccumulations:     episodesToDTO(h.PrimitiveAccumulations),
		TotalCapitalFormed:         int64(h.TotalCapitalFormed()),
		TotalLabourersExpropriated: h.TotalLabourersExpropriated(),
		CreatedAt:                  h.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
	}
}

// CreateHistoricalStage handles POST /v1/historical-stages.
func (h *Handler) CreateHistoricalStage(w http.ResponseWriter, r *http.Request) {
	var req createHistoricalStageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if len(req.PrimitiveAccumulations) == 0 {
		writeError(w, http.StatusBadRequest, "at least one primitive accumulation is required")
		return
	}

	stage := simulation.HistoricalStage{
		Name:                   req.Name,
		Description:            req.Description,
		PrimitiveAccumulations: episodesFromDTO(req.PrimitiveAccumulations),
	}
	if err := stage.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.HistoricalStages.CreateHistoricalStage(r.Context(), stage)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "historical stage with that name already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/v1/historical-stages/%s", string(created.ID)))
	writeJSON(w, http.StatusCreated, toHistoricalStageResponse(created))
}

// ListHistoricalStages handles GET /v1/historical-stages.
func (h *Handler) ListHistoricalStages(w http.ResponseWriter, r *http.Request) {
	stages, err := h.HistoricalStages.ListHistoricalStages(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]historicalStageResponse, 0, len(stages))
	for _, s := range stages {
		out = append(out, toHistoricalStageResponse(s))
	}
	writeJSON(w, http.StatusOK, out)
}

// SeedScenarioFromStage handles POST /v1/historical-stages/{id}/seed-scenario.
// It derives a starting CapitalStock and ProducerSeparation from the stage,
// then runs a simple-reproduction series. The series is returned inline; no
// new persisted scenario is created — Ch. 23/24 reproduction remain stateless.
func (h *Handler) SeedScenarioFromStage(w http.ResponseWriter, r *http.Request) {
	id := simulation.HistoricalStageID(r.PathValue("id"))
	if id.IsZero() {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req seedScenarioRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.PreCapitalistWorkers < 0 {
		writeError(w, http.StatusBadRequest, "pre_capitalist_workers cannot be negative")
		return
	}
	if req.CompositionRatio < 0 || req.CompositionRatio >= 1 {
		writeError(w, http.StatusBadRequest, "composition_ratio must be in [0, 1)")
		return
	}
	if req.SurplusRate < 0 {
		writeError(w, http.StatusBadRequest, "surplus_rate cannot be negative")
		return
	}
	if req.Periods <= 0 {
		writeError(w, http.StatusBadRequest, "periods must be positive")
		return
	}

	stage, err := h.HistoricalStages.GetHistoricalStage(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "historical stage not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	separation := simulation.SeparationFromStage(stage, req.PreCapitalistWorkers)
	seeded := simulation.SeedCapitalStock(stage, req.CompositionRatio)
	cycles := simulation.RunSimpleReproduction(seeded, req.SurplusRate, req.Periods)

	dto := make([]seededReproductionCycleDTO, len(cycles))
	for i, c := range cycles {
		dto[i] = seededReproductionCycleDTO{
			Period: c.Period,
			Capital: seededCapitalStockDTO{
				ConstantCapital: int64(c.Capital.ConstantCapital),
				VariableCapital: int64(c.Capital.VariableCapital),
			},
			SurplusRate: c.SurplusRate,
			Fund: seededSurplusValueFundDTO{
				Total:       int64(c.Fund.Total),
				Revenue:     int64(c.Fund.Revenue),
				Accumulated: int64(c.Fund.Accumulated),
			},
		}
	}

	writeJSON(w, http.StatusOK, seedScenarioResponse{
		StageID:   string(stage.ID),
		StageName: stage.Name,
		SeededCapital: seededCapitalStockDTO{
			ConstantCapital: int64(seeded.ConstantCapital),
			VariableCapital: int64(seeded.VariableCapital),
		},
		Separation: producerSeparationDTO{
			PreCapitalistWorkers: separation.PreCapitalistWorkers,
			DisplacedWorkers:     separation.DisplacedWorkers,
			FreeProletarians:     separation.FreeProletarians,
		},
		SurplusRate: req.SurplusRate,
		Periods:     req.Periods,
		Cycles:      dto,
	})
}
