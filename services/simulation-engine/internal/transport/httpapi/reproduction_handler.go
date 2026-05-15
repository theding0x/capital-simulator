package httpapi

import (
	"math"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// Ch. 23 — Simple Reproduction. Both endpoints are stateless compute endpoints.

type simpleReproductionRequest struct {
	ConstantCapital int64   `json:"constant_capital"`
	VariableCapital int64   `json:"variable_capital"`
	SurplusRate     float64 `json:"surplus_rate"`
	Periods         int64   `json:"periods"`
}

type cycleResponse struct {
	Period          int64   `json:"period"`
	ConstantCapital int64   `json:"constant_capital"`
	VariableCapital int64   `json:"variable_capital"`
	SurplusRate     float64 `json:"surplus_rate"`
	SurplusTotal    int64   `json:"surplus_total"`
	Revenue         int64   `json:"revenue"`
	Accumulated     int64   `json:"accumulated"`
}

type simpleReproductionResponse struct {
	ConstantCapital int64           `json:"constant_capital"`
	VariableCapital int64           `json:"variable_capital"`
	SurplusRate     float64         `json:"surplus_rate"`
	Periods         int64           `json:"periods"`
	Cycles          []cycleResponse `json:"cycles"`
	RepaymentPeriod int64           `json:"repayment_period"`
}

type repaymentPeriodRequest struct {
	Capital       int64 `json:"capital"`
	AnnualRevenue int64 `json:"annual_revenue"`
}

type repaymentPeriodResponse struct {
	Capital         int64 `json:"capital"`
	AnnualRevenue   int64 `json:"annual_revenue"`
	RepaymentPeriod int64 `json:"repayment_period"`
}

// RunSimpleReproduction handles POST /v1/reproductions/simple.
func (h *Handler) RunSimpleReproduction(w http.ResponseWriter, r *http.Request) {
	var req simpleReproductionRequest
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
	if req.Periods <= 0 {
		writeError(w, http.StatusBadRequest, "periods must be positive")
		return
	}

	initial := simulation.CapitalStock{
		ConstantCapital: simulation.Pence(req.ConstantCapital),
		VariableCapital: simulation.Pence(req.VariableCapital),
	}
	cycles := simulation.RunSimpleReproduction(initial, req.SurplusRate, req.Periods)

	responses := make([]cycleResponse, len(cycles))
	for i, c := range cycles {
		responses[i] = cycleResponse{
			Period:          c.Period,
			ConstantCapital: int64(c.Capital.ConstantCapital),
			VariableCapital: int64(c.Capital.VariableCapital),
			SurplusRate:     c.SurplusRate,
			SurplusTotal:    int64(c.Fund.Total),
			Revenue:         int64(c.Fund.Revenue),
			Accumulated:     int64(c.Fund.Accumulated),
		}
	}

	annualRevenue := simulation.Pence(math.Round(float64(req.VariableCapital) * req.SurplusRate))
	rp := simulation.RepaymentPeriod(
		simulation.Pence(req.ConstantCapital+req.VariableCapital),
		annualRevenue,
	)

	writeJSON(w, http.StatusOK, simpleReproductionResponse{
		ConstantCapital: req.ConstantCapital,
		VariableCapital: req.VariableCapital,
		SurplusRate:     req.SurplusRate,
		Periods:         req.Periods,
		Cycles:          responses,
		RepaymentPeriod: rp,
	})
}

// ComputeRepaymentPeriod handles POST /v1/reproductions/repayment-period.
func (h *Handler) ComputeRepaymentPeriod(w http.ResponseWriter, r *http.Request) {
	var req repaymentPeriodRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Capital <= 0 {
		writeError(w, http.StatusBadRequest, "capital must be positive")
		return
	}
	if req.AnnualRevenue <= 0 {
		writeError(w, http.StatusBadRequest, "annual_revenue must be positive")
		return
	}

	rp := simulation.RepaymentPeriod(
		simulation.Pence(req.Capital),
		simulation.Pence(req.AnnualRevenue),
	)
	writeJSON(w, http.StatusOK, repaymentPeriodResponse{
		Capital:         req.Capital,
		AnnualRevenue:   req.AnnualRevenue,
		RepaymentPeriod: rp,
	})
}
