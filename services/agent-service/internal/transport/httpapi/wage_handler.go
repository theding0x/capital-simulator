package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

type createWageFormRequest struct {
	AgentID           string `json:"agent_id"`
	DailyPence        int64  `json:"daily_pence"`
	WorkingDayMinutes int64  `json:"working_day_minutes"`
	LPVDailyPence     int64  `json:"lpv_daily_pence"`
	NecessaryMinutes  int64  `json:"necessary_minutes"`
}

type wageFormResponse struct {
	agent.WageForm
	HourlyWage    agent.Pence               `json:"hourly_wage"`
	Appearance    agent.WageAppearance      `json:"appearance"`
	Decomposition agent.LabourDecomposition `json:"decomposition"`
}

func buildWageFormResponse(wf agent.WageForm) wageFormResponse {
	return wageFormResponse{
		WageForm:      wf,
		HourlyWage:    agent.HourlyWage(wf.Wage),
		Appearance:    agent.Appearance(wf.Wage),
		Decomposition: agent.Decompose(wf.Wage, wf.LabourPowerValue),
	}
}

func (h *Handler) CreateWageForm(w http.ResponseWriter, r *http.Request) {
	var req createWageFormRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	wf := agent.WageForm{
		AgentID: agent.AgentID(req.AgentID),
		Wage: agent.Wage{
			DailyPence:        agent.Pence(req.DailyPence),
			WorkingDayMinutes: agent.LabourMinutes(req.WorkingDayMinutes),
		},
		LabourPowerValue: agent.WageLabourValue{
			DailyPence:       agent.Pence(req.LPVDailyPence),
			NecessaryMinutes: agent.LabourMinutes(req.NecessaryMinutes),
		},
	}
	if err := wf.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	// If this agent has a Ch.6 LabourWorker record, NecessaryMinutes must
	// equal labour_power_value_minutes — both express the same reproduction cost.
	worker, err := h.LabourPowerStore.GetWorker(r.Context(), wf.AgentID)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		writeAppError(w, err)
		return
	}
	if err == nil && worker.LabourPowerValueMinutes != wf.LabourPowerValue.NecessaryMinutes {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf(
			"necessary_minutes %d does not match worker labour_power_value_minutes %d",
			wf.LabourPowerValue.NecessaryMinutes, worker.LabourPowerValueMinutes,
		))
		return
	}
	saved, err := h.WageFormStore.CreateWageForm(r.Context(), wf)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, buildWageFormResponse(saved))
}

func (h *Handler) GetWageForm(w http.ResponseWriter, r *http.Request) {
	agentID := agent.AgentID(r.PathValue("agentID"))
	wf, err := h.WageFormStore.GetWageForm(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildWageFormResponse(wf))
}

// ListWageForms handles GET /v1/wage-forms.
func (h *Handler) ListWageForms(w http.ResponseWriter, r *http.Request) {
	wfs, err := h.WageFormStore.ListWageForms(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	out := make([]wageFormResponse, len(wfs))
	for i, wf := range wfs {
		out[i] = buildWageFormResponse(wf)
	}
	writeJSON(w, http.StatusOK, out)
}
