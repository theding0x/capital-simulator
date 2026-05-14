package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

type computeHourlyPriceRequest struct {
	DailyPence        int64 `json:"daily_pence"`
	WorkingDayMinutes int64 `json:"working_day_minutes"`
}

type hourlyPriceResponse struct {
	Numerator   int64   `json:"numerator"`
	Denominator int64   `json:"denominator"`
	AsFloat     float64 `json:"as_float"`
}

type createWorkingSessionRequest struct {
	AgentID               string `json:"agent_id"`
	DailyLabourPowerValue int64  `json:"daily_labour_power_value"`
	WorkingDayMinutes     int64  `json:"working_day_minutes"`
	OvertimeHours         int64  `json:"overtime_hours"`
	OvertimeRatePence     int64  `json:"overtime_rate_pence"`
	WagePeriod            string `json:"wage_period"`
}

type workingSessionResponse struct {
	agent.WorkingSession
	HourlyPrice hourlyPriceResponse `json:"hourly_price"`
	NominalWage agent.NominalWage   `json:"nominal_wage"`
}

func buildWorkingSessionResponse(s agent.WorkingSession) workingSessionResponse {
	hp := agent.ComputeHourlyPrice(s.DailyLabourPowerValue, s.WorkingDayMinutes)
	return workingSessionResponse{
		WorkingSession: s,
		HourlyPrice: hourlyPriceResponse{
			Numerator:   hp.Numerator,
			Denominator: hp.Denominator,
			AsFloat:     hp.AsFloat(),
		},
		NominalWage: agent.ComputeSessionWage(s),
	}
}

// ComputeHourlyPrice is a stateless endpoint: POST /v1/time-wages/hourly-price.
func (h *Handler) ComputeHourlyPrice(w http.ResponseWriter, r *http.Request) {
	var req computeHourlyPriceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.DailyPence <= 0 {
		writeError(w, http.StatusBadRequest, "daily_pence must be positive")
		return
	}
	if req.WorkingDayMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "working_day_minutes must be positive")
		return
	}
	hp := agent.ComputeHourlyPrice(
		agent.DailyLabourPowerValue{Pence: req.DailyPence},
		agent.WorkingDayMinutes{Minutes: agent.LabourMinutes(req.WorkingDayMinutes)},
	)
	writeJSON(w, http.StatusOK, hourlyPriceResponse{
		Numerator:   hp.Numerator,
		Denominator: hp.Denominator,
		AsFloat:     hp.AsFloat(),
	})
}

// CreateWorkingSession handles POST /v1/time-wages/sessions.
func (h *Handler) CreateWorkingSession(w http.ResponseWriter, r *http.Request) {
	var req createWorkingSessionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s := agent.WorkingSession{
		AgentID:               agent.AgentID(req.AgentID),
		DailyLabourPowerValue: agent.DailyLabourPowerValue{Pence: req.DailyLabourPowerValue},
		WorkingDayMinutes:     agent.WorkingDayMinutes{Minutes: agent.LabourMinutes(req.WorkingDayMinutes)},
		OvertimeHours:         agent.OvertimeHours{Hours: req.OvertimeHours},
		OvertimeRatePence:     agent.OvertimeRatePence{Pence: req.OvertimeRatePence},
		WagePeriod:            agent.WagePeriod(req.WagePeriod),
	}
	if err := s.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.TimeWageStore.CreateWorkingSession(r.Context(), s)
	if err != nil {
		writeAppError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/time-wages/sessions/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, buildWorkingSessionResponse(saved))
}

// GetWorkingSession handles GET /v1/time-wages/sessions/{id}.
func (h *Handler) GetWorkingSession(w http.ResponseWriter, r *http.Request) {
	id := agent.WorkingSessionID(r.PathValue("id"))
	s, err := h.TimeWageStore.GetWorkingSession(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildWorkingSessionResponse(s))
}
