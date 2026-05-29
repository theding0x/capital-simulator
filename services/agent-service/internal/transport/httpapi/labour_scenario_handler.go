package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// Ch. 17 — Changes of Magnitude in the Price of Labour-Power and in
// Surplus-Value. The endpoint is stateless: it accepts a LabourScenario
// (working day partition, intensity factor, productivity factor) and
// returns the ScenarioOutcome (daily value, partition, rate of surplus).

type labourScenarioRequest struct {
	WorkingDayMinutes      int64   `json:"working_day_minutes"`
	NecessaryLabourMinutes int64   `json:"necessary_labour_minutes"`
	IntensityFactor        float64 `json:"intensity_factor"`
	ProductivityFactor     float64 `json:"productivity_factor"`
}

type labourScenarioResponse struct {
	DailyValueMinutes       int64   `json:"daily_value_minutes"`
	NecessaryLabourMinutes  int64   `json:"necessary_labour_minutes"`
	SurplusLabourMinutes    int64   `json:"surplus_labour_minutes"`
	LabourPowerValueMinutes int64   `json:"labour_power_value_minutes"`
	RateOfSurplusValue      float64 `json:"rate_of_surplus_value"`
	LawConstantDailyValue   bool    `json:"law_constant_daily_value"`
}

func (h *Handler) ComputeLabourScenario(w http.ResponseWriter, r *http.Request) {
	var req labourScenarioRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.WorkingDayMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "working_day_minutes must be positive")
		return
	}
	if req.NecessaryLabourMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "necessary_labour_minutes must be positive")
		return
	}
	if req.NecessaryLabourMinutes > req.WorkingDayMinutes {
		writeError(w, http.StatusBadRequest, "necessary_labour_minutes must be <= working_day_minutes")
		return
	}
	intensity := agent.LabourIntensity{Factor: req.IntensityFactor}
	if req.IntensityFactor == 0 {
		intensity = agent.NormalLabourIntensity
	}
	productivity := agent.LabourProductivity{Factor: req.ProductivityFactor}
	if req.ProductivityFactor == 0 {
		productivity = agent.NormalLabourProductivity
	}
	scenario := agent.LabourScenario{
		WorkingDay: agent.WorkingDay{
			NecessaryLabourMinutes: agent.NecessaryLabourMinutes(req.NecessaryLabourMinutes),
			SurplusLabourMinutes:   agent.SurplusLabourMinutes(req.WorkingDayMinutes - req.NecessaryLabourMinutes),
		},
		Intensity:    intensity,
		Productivity: productivity,
	}
	if err := scenario.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	out := agent.ComputeOutcome(scenario)
	writeJSON(w, http.StatusOK, labourScenarioResponse{
		DailyValueMinutes:       int64(out.DailyValueMinutes),
		NecessaryLabourMinutes:  int64(out.NecessaryLabourMinutes),
		SurplusLabourMinutes:    int64(out.SurplusLabourMinutes),
		LabourPowerValueMinutes: int64(out.LabourPowerValueMinutes),
		RateOfSurplusValue:      out.RateOfSurplusValue,
		LawConstantDailyValue:   agent.LawConstantDailyValue(scenario.WorkingDay, scenario.Intensity, scenario.Productivity.Factor),
	})
}
