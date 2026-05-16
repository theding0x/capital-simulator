package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// Ch. 24 — The Transformation of Surplus-Value into Capital.
// Both endpoints are stateless compute endpoints.

type extendedReproductionRequest struct {
	ConstantCapital  int64   `json:"constant_capital"`
	VariableCapital  int64   `json:"variable_capital"`
	SurplusRate      float64 `json:"surplus_rate"`
	AccumRate        float64 `json:"accum_rate"`
	CompositionRatio float64 `json:"composition_ratio"`
	Periods          int64   `json:"periods"`
}

type accumulationPeriodResponse struct {
	Period          int64   `json:"period"`
	ConstantCapital int64   `json:"constant_capital"`
	VariableCapital int64   `json:"variable_capital"`
	SurplusProduced int64   `json:"surplus_produced"`
	NewConstant     int64   `json:"new_constant"`
	NewVariable     int64   `json:"new_variable"`
	Revenue         int64   `json:"revenue"`
}

type extendedReproductionResponse struct {
	ConstantCapital  int64                        `json:"constant_capital"`
	VariableCapital  int64                        `json:"variable_capital"`
	SurplusRate      float64                      `json:"surplus_rate"`
	AccumRate        float64                      `json:"accum_rate"`
	CompositionRatio float64                      `json:"composition_ratio"`
	Periods          int64                        `json:"periods"`
	Cycles           []accumulationPeriodResponse `json:"cycles"`
}

type splitSurplusRequest struct {
	Surplus          int64   `json:"surplus"`
	AccumRate        float64 `json:"accum_rate"`
	CompositionRatio float64 `json:"composition_ratio"`
}

type splitSurplusResponse struct {
	Surplus          int64   `json:"surplus"`
	AccumRate        float64 `json:"accum_rate"`
	CompositionRatio float64 `json:"composition_ratio"`
	NewConstant      int64   `json:"new_constant"`
	NewVariable      int64   `json:"new_variable"`
	Revenue          int64   `json:"revenue"`
}

// RunExtendedReproduction handles POST /v1/reproductions/extended.
func (h *Handler) RunExtendedReproduction(w http.ResponseWriter, r *http.Request) {
	var req extendedReproductionRequest
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
	if req.SurplusRate < 0 {
		writeError(w, http.StatusBadRequest, "surplus_rate cannot be negative")
		return
	}
	if req.AccumRate < 0 || req.AccumRate > 1 {
		writeError(w, http.StatusBadRequest, "accum_rate must be between 0 and 1")
		return
	}
	if req.CompositionRatio < 0 || req.CompositionRatio > 1 {
		writeError(w, http.StatusBadRequest, "composition_ratio must be between 0 and 1")
		return
	}
	if req.Periods <= 0 {
		writeError(w, http.StatusBadRequest, "periods must be positive")
		return
	}

	initial := simulation.CapitalStock{
		ConstantCapital: simulation.Pence(req.ConstantCapital),
		VariableCapital: simulation.Pence(req.VariableCapital),
	}
	cycles := simulation.RunExtendedReproduction(initial, req.SurplusRate, req.AccumRate, req.CompositionRatio, req.Periods)

	responses := make([]accumulationPeriodResponse, len(cycles))
	for i, c := range cycles {
		responses[i] = accumulationPeriodResponse{
			Period:          c.Period,
			ConstantCapital: int64(c.CapitalStock.ConstantCapital),
			VariableCapital: int64(c.CapitalStock.VariableCapital),
			SurplusProduced: int64(c.SurplusProduced),
			NewConstant:     int64(c.NewConstant),
			NewVariable:     int64(c.NewVariable),
			Revenue:         int64(c.Revenue),
		}
	}

	writeJSON(w, http.StatusOK, extendedReproductionResponse{
		ConstantCapital:  req.ConstantCapital,
		VariableCapital:  req.VariableCapital,
		SurplusRate:      req.SurplusRate,
		AccumRate:        req.AccumRate,
		CompositionRatio: req.CompositionRatio,
		Periods:          req.Periods,
		Cycles:           responses,
	})
}

// SplitSurplusValue handles POST /v1/reproductions/split-surplus.
func (h *Handler) SplitSurplusValue(w http.ResponseWriter, r *http.Request) {
	var req splitSurplusRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Surplus < 0 {
		writeError(w, http.StatusBadRequest, "surplus cannot be negative")
		return
	}
	if req.AccumRate < 0 || req.AccumRate > 1 {
		writeError(w, http.StatusBadRequest, "accum_rate must be between 0 and 1")
		return
	}
	if req.CompositionRatio < 0 || req.CompositionRatio > 1 {
		writeError(w, http.StatusBadRequest, "composition_ratio must be between 0 and 1")
		return
	}

	additional := simulation.SplitSurplus(simulation.Pence(req.Surplus), req.AccumRate, req.CompositionRatio)
	revenue := simulation.Pence(req.Surplus) - additional.Constant - additional.Variable

	writeJSON(w, http.StatusOK, splitSurplusResponse{
		Surplus:          req.Surplus,
		AccumRate:        req.AccumRate,
		CompositionRatio: req.CompositionRatio,
		NewConstant:      int64(additional.Constant),
		NewVariable:      int64(additional.Variable),
		Revenue:          int64(revenue),
	})
}
