package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

type createWorkingDayRequest struct {
	NecessaryLabourMinutes int64  `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   int64  `json:"surplus_labour_minutes"`
	StatutoryLimitMinutes  *int64 `json:"statutory_limit_minutes,omitempty"`
}

type workingDayResponse struct {
	WorkingDay         agent.WorkingDay `json:"working_day"`
	TotalMinutes       int64            `json:"total_minutes"`
	RateOfSurplusValue float64          `json:"rate_of_surplus_value"`
	ExceedsStatutory   bool             `json:"exceeds_statutory,omitempty"`
}

func buildWorkingDayResponse(wd agent.WorkingDay, limit *int64) workingDayResponse {
	resp := workingDayResponse{
		WorkingDay:         wd,
		TotalMinutes:       wd.TotalMinutes(),
		RateOfSurplusValue: wd.RateOfSurplusValue(),
	}
	if limit != nil {
		resp.ExceedsStatutory = wd.TotalMinutes() > *limit
	}
	return resp
}

func (h *Handler) CreateWorkingDay(w http.ResponseWriter, r *http.Request) {
	var req createWorkingDayRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	wd := agent.WorkingDay{
		NecessaryLabourMinutes: agent.NecessaryLabourMinutes(req.NecessaryLabourMinutes),
		SurplusLabourMinutes:   agent.SurplusLabourMinutes(req.SurplusLabourMinutes),
	}
	if err := wd.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.StatutoryLimitMinutes != nil {
		c := agent.WorkingDayConstraint{StatutoryLimitMinutes: agent.StatutoryLimitMinutes(*req.StatutoryLimitMinutes)}
		if err := agent.ValidateConstraint(wd, c); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	saved, err := h.WorkingDayStore.CreateWorkingDay(r.Context(), wd)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, buildWorkingDayResponse(saved, req.StatutoryLimitMinutes))
}

func (h *Handler) GetWorkingDay(w http.ResponseWriter, r *http.Request) {
	id := agent.WorkingDayID(r.PathValue("id"))
	wd, err := h.WorkingDayStore.GetWorkingDay(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildWorkingDayResponse(wd, nil))
}

type validateWorkingDayRequest struct {
	NecessaryLabourMinutes int64  `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   int64  `json:"surplus_labour_minutes"`
	StatutoryLimitMinutes  *int64 `json:"statutory_limit_minutes,omitempty"`
}

type validateWorkingDayResponse struct {
	TotalMinutes       int64   `json:"total_minutes"`
	RateOfSurplusValue float64 `json:"rate_of_surplus_value"`
	Valid              bool    `json:"valid"`
	Error              string  `json:"error,omitempty"`
}

func (h *Handler) ValidateWorkingDay(w http.ResponseWriter, r *http.Request) {
	var req validateWorkingDayRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	wd := agent.WorkingDay{
		NecessaryLabourMinutes: agent.NecessaryLabourMinutes(req.NecessaryLabourMinutes),
		SurplusLabourMinutes:   agent.SurplusLabourMinutes(req.SurplusLabourMinutes),
	}
	resp := validateWorkingDayResponse{
		TotalMinutes: wd.TotalMinutes(),
	}
	if wd.NecessaryLabourMinutes > 0 {
		resp.RateOfSurplusValue = wd.RateOfSurplusValue()
	}
	if err := wd.Validate(); err != nil {
		resp.Error = err.Error()
		writeJSON(w, http.StatusOK, resp)
		return
	}
	if req.StatutoryLimitMinutes != nil {
		c := agent.WorkingDayConstraint{StatutoryLimitMinutes: agent.StatutoryLimitMinutes(*req.StatutoryLimitMinutes)}
		if err := agent.ValidateConstraint(wd, c); err != nil {
			resp.Error = err.Error()
			writeJSON(w, http.StatusOK, resp)
			return
		}
	}
	resp.Valid = true
	writeJSON(w, http.StatusOK, resp)
}

type relaySetInput struct {
	ShiftKind              string   `json:"shift_kind"`
	NecessaryLabourMinutes int64    `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   int64    `json:"surplus_labour_minutes"`
	WorkerIDs              []string `json:"worker_ids"`
}

type createRelayScheduleRequest struct {
	Sets [2]relaySetInput `json:"sets"`
}

func (h *Handler) CreateRelaySchedule(w http.ResponseWriter, r *http.Request) {
	var req createRelayScheduleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rs := agent.RelaySchedule{}
	for i, s := range req.Sets {
		wids := make([]agent.AgentID, len(s.WorkerIDs))
		for j, id := range s.WorkerIDs {
			wids[j] = agent.AgentID(id)
		}
		rs.Sets[i] = agent.RelaySet{
			ShiftKind: agent.ShiftKind(s.ShiftKind),
			WorkingDay: agent.WorkingDay{
				NecessaryLabourMinutes: agent.NecessaryLabourMinutes(s.NecessaryLabourMinutes),
				SurplusLabourMinutes:   agent.SurplusLabourMinutes(s.SurplusLabourMinutes),
			},
			WorkerIDs: wids,
		}
	}
	if err := rs.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.WorkingDayStore.CreateRelaySchedule(r.Context(), rs)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}

func (h *Handler) GetRelaySchedule(w http.ResponseWriter, r *http.Request) {
	id := agent.RelayScheduleID(r.PathValue("id"))
	rs, err := h.WorkingDayStore.GetRelaySchedule(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rs)
}
