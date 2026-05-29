package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/labour"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// Ch. 18 — Various Formula for the Rate of Surplus-Value. Stateless endpoint.

type ratesRequest struct {
	NecessaryLabourMinutes int64 `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   int64 `json:"surplus_labour_minutes"`
}

type ratesResponse struct {
	NecessaryLabourMinutes int64   `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   int64   `json:"surplus_labour_minutes"`
	WorkingDayMinutes      int64   `json:"working_day_minutes"`
	FormulaI               float64 `json:"formula_i"`
	FormulaII              float64 `json:"formula_ii"`
	FormulaIII             float64 `json:"formula_iii"`
}

// ComputeRatesOfSurplusValue handles POST /v1/surplus-value/rates.
// Applies Marx's three equivalent formulae from Ch. 18 to the given
// labour-time partition and returns the results side by side.
func (h *Handler) ComputeRatesOfSurplusValue(w http.ResponseWriter, r *http.Request) {
	var req ratesRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.NecessaryLabourMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "necessary_labour_minutes must be positive")
		return
	}
	if req.SurplusLabourMinutes < 0 {
		writeError(w, http.StatusBadRequest, "surplus_labour_minutes cannot be negative")
		return
	}
	sc := simulation.RateScenario{
		NecessaryLabour: simulation.NecessaryLabour{Minutes: labour.LabourMinutes(req.NecessaryLabourMinutes)},
		SurplusLabour:   simulation.SurplusLabour{Minutes: labour.LabourMinutes(req.SurplusLabourMinutes)},
	}
	result := simulation.ComputeRates(sc)
	writeJSON(w, http.StatusOK, ratesResponse{
		NecessaryLabourMinutes: req.NecessaryLabourMinutes,
		SurplusLabourMinutes:   req.SurplusLabourMinutes,
		WorkingDayMinutes:      req.NecessaryLabourMinutes + req.SurplusLabourMinutes,
		FormulaI:               result.FormulaI,
		FormulaII:              result.FormulaII,
		FormulaIII:             result.FormulaIII,
	})
}
