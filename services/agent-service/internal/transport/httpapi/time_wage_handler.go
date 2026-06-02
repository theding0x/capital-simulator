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
	WageFormID            string `json:"wage_form_id,omitempty"`
	DailyLabourPowerValue int64  `json:"daily_labour_power_value,omitempty"`
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
		WageFormID:            agent.WageFormID(req.WageFormID),
		DailyLabourPowerValue: agent.DailyLabourPowerValue{Pence: req.DailyLabourPowerValue},
		WorkingDayMinutes:     agent.WorkingDayMinutes{Minutes: agent.LabourMinutes(req.WorkingDayMinutes)},
		OvertimeHours:         agent.OvertimeHours{Hours: req.OvertimeHours},
		OvertimeRatePence:     agent.OvertimeRatePence{Pence: req.OvertimeRatePence},
		WagePeriod:            agent.WagePeriod(req.WagePeriod),
	}
	if !s.WageFormID.IsZero() {
		wf, err := h.WageFormStore.GetWageForm(r.Context(), s.AgentID)
		if err != nil {
			writeAppError(w, err)
			return
		}
		s.DailyLabourPowerValue = agent.DailyLabourPowerValue{Pence: int64(wf.LabourPowerValue.DailyPence)}
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

// ListWorkingSessions handles GET /v1/agents/{id}/time-wages/sessions.
func (h *Handler) ListWorkingSessions(w http.ResponseWriter, r *http.Request) {
	agentID := agent.AgentID(r.PathValue("id"))
	sessions, err := h.TimeWageStore.ListWorkingSessions(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	out := make([]workingSessionResponse, len(sessions))
	for i, s := range sessions {
		out[i] = buildWorkingSessionResponse(s)
	}
	writeJSON(w, http.StatusOK, out)
}

type nominalWageRequest struct {
	DailyLabourPowerValue int64 `json:"daily_labour_power_value"`
	WorkingDayMinutes     int64 `json:"working_day_minutes"`
	OvertimeHours         int64 `json:"overtime_hours,omitempty"`
	OvertimeRatePence     int64 `json:"overtime_rate_pence,omitempty"`

	// Optional Ch. 25 reserve-army overlay (see agent.ReserveArmyPressure).
	ReserveArmyMagnitude      int64 `json:"reserve_army_magnitude,omitempty"`
	ReserveArmyWorkforceTotal int64 `json:"reserve_army_workforce_total,omitempty"`
}

type nominalWageResponse struct {
	NominalWage           agent.NominalWage `json:"nominal_wage"`
	ReserveArmyPressureBP int64             `json:"reserve_army_pressure_bp,omitempty"`
	CompressedWagePence   int64             `json:"compressed_wage_pence,omitempty"`
}

// ComputeNominalWage is a stateless endpoint: POST /v1/time-wages/nominal-wage.
// It returns the nominal money-wage for a working session and, when a reserve
// army overlay is supplied, the wage compressed below it (Ch. 25 → Ch. 20).
func (h *Handler) ComputeNominalWage(w http.ResponseWriter, r *http.Request) {
	var req nominalWageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.DailyLabourPowerValue <= 0 {
		writeError(w, http.StatusBadRequest, "daily_labour_power_value must be positive")
		return
	}
	if req.WorkingDayMinutes < 60 {
		writeError(w, http.StatusBadRequest, "working_day_minutes must be at least 60")
		return
	}
	if req.OvertimeHours < 0 || req.OvertimeRatePence < 0 {
		writeError(w, http.StatusBadRequest, "overtime values cannot be negative")
		return
	}
	session := agent.WorkingSession{
		DailyLabourPowerValue: agent.DailyLabourPowerValue{Pence: req.DailyLabourPowerValue},
		WorkingDayMinutes:     agent.WorkingDayMinutes{Minutes: agent.LabourMinutes(req.WorkingDayMinutes)},
		OvertimeHours:         agent.OvertimeHours{Hours: req.OvertimeHours},
		OvertimeRatePence:     agent.OvertimeRatePence{Pence: req.OvertimeRatePence},
		WagePeriod:            agent.WagePeriodDaily,
	}
	nominal := agent.ComputeSessionWage(session)
	resp := nominalWageResponse{NominalWage: nominal}
	if req.ReserveArmyMagnitude != 0 || req.ReserveArmyWorkforceTotal != 0 {
		pressure := agent.ReserveArmyPressure{
			Magnitude:      req.ReserveArmyMagnitude,
			WorkforceTotal: req.ReserveArmyWorkforceTotal,
		}
		if err := pressure.Validate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		resp.ReserveArmyPressureBP = pressure.PressureBP()
		resp.CompressedWagePence = pressure.CompressPence(nominal.Pence)
	}
	writeJSON(w, http.StatusOK, resp)
}
