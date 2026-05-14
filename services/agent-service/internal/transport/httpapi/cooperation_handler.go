package httpapi

import (
	"net/http"
	"strings"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

type cooperationMemberInput struct {
	WorkerID          string `json:"worker_id"`
	Supervisory       bool   `json:"supervisory"`
	WorkingDayMinutes int64  `json:"working_day_minutes"`
}

type createCooperationRequest struct {
	Name         string                   `json:"name"`
	CapitalistID string                   `json:"capitalist_id"`
	Members      []cooperationMemberInput `json:"members"`
}

type cooperationResponse struct {
	agent.Cooperation
	Size                  agent.CooperationSize `json:"size"`
	CollectiveWorkingDay  agent.LabourMinutes   `json:"collective_working_day_minutes"`
	AverageSocialLabour   agent.LabourMinutes   `json:"average_social_labour_minutes"`
	CollectivePowerFactor float64               `json:"collective_power_factor"`
	ProductivityFactor    float64               `json:"productivity_factor"`
	CooperativeOrigin     string                `json:"cooperative_origin"`
	Supervisors           []agent.Supervisor    `json:"supervisors"`
}

func buildCooperationResponse(c agent.Cooperation) cooperationResponse {
	n := c.Size()
	d := c.AverageIndividualWorkingDay()
	collective := agent.CollectiveWorkingDay(n, d)
	if c.Members == nil {
		c.Members = []agent.CooperationMember{}
	}
	factor := agent.CollectiveProductivePower(n, d)
	return cooperationResponse{
		Cooperation:           c,
		Size:                  n,
		CollectiveWorkingDay:  collective,
		AverageSocialLabour:   agent.AverageSocialLabour(collective, n),
		CollectivePowerFactor: factor,
		ProductivityFactor:    factor, // Part IV bridge: feed into Ch. 12.
		CooperativeOrigin:     agent.CooperativeOrigin(),
		Supervisors:           c.Supervisors(),
	}
}

func (h *Handler) CreateCooperation(w http.ResponseWriter, r *http.Request) {
	var req createCooperationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	members := make([]agent.CooperationMember, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, agent.CooperationMember{
			WorkerID:          agent.AgentID(m.WorkerID),
			Supervisory:       m.Supervisory,
			WorkingDayMinutes: agent.LabourMinutes(m.WorkingDayMinutes),
		})
	}
	c := agent.Cooperation{
		Name:         strings.TrimSpace(req.Name),
		CapitalistID: agent.AgentID(req.CapitalistID),
		Members:      members,
	}
	if err := c.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.CooperationStore.CreateCooperation(r.Context(), c)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, buildCooperationResponse(saved))
}

func (h *Handler) GetCooperation(w http.ResponseWriter, r *http.Request) {
	id := agent.CooperationID(r.PathValue("id"))
	c, err := h.CooperationStore.GetCooperation(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildCooperationResponse(c))
}

func (h *Handler) ListCooperations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capParam := r.URL.Query().Get("capitalist_id")
	var (
		coops []agent.Cooperation
		err   error
	)
	if capParam != "" {
		coops, err = h.CooperationStore.ListCooperationsByCapitalist(ctx, agent.AgentID(capParam))
	} else {
		coops, err = h.CooperationStore.ListCooperations(ctx)
	}
	if err != nil {
		writeAppError(w, err)
		return
	}
	items := make([]cooperationResponse, 0, len(coops))
	for _, c := range coops {
		items = append(items, buildCooperationResponse(c))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

type collectiveWorkingDayResponse struct {
	CooperationID                 agent.CooperationID   `json:"cooperation_id"`
	Size                          agent.CooperationSize `json:"size"`
	IndividualWorkingDayMinutes   agent.LabourMinutes   `json:"individual_working_day_minutes"`
	CollectiveWorkingDayMinutes   agent.LabourMinutes   `json:"collective_working_day_minutes"`
	CollectivePowerFactor         float64               `json:"collective_power_factor"`
	CollectiveOutputUseValueUnits float64               `json:"collective_output_use_value_units"`
	CooperativeOrigin             string                `json:"cooperative_origin"`
}

func (h *Handler) ComputeCollectiveWorkingDay(w http.ResponseWriter, r *http.Request) {
	id := agent.CooperationID(r.PathValue("id"))
	c, err := h.CooperationStore.GetCooperation(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	n := c.Size()
	d := c.AverageIndividualWorkingDay()
	collective := agent.CollectiveWorkingDay(n, d)
	factor := agent.CollectiveProductivePower(n, d)
	writeJSON(w, http.StatusOK, collectiveWorkingDayResponse{
		CooperationID:                 c.ID,
		Size:                          n,
		IndividualWorkingDayMinutes:   d,
		CollectiveWorkingDayMinutes:   collective,
		CollectivePowerFactor:         factor,
		CollectiveOutputUseValueUnits: float64(collective) * factor,
		CooperativeOrigin:             agent.CooperativeOrigin(),
	})
}

type averageSocialLabourResponse struct {
	CooperationID                 agent.CooperationID   `json:"cooperation_id"`
	Size                          agent.CooperationSize `json:"size"`
	CollectiveWorkingDayMinutes   agent.LabourMinutes   `json:"collective_working_day_minutes"`
	AverageSocialLabourMinutes    agent.LabourMinutes   `json:"average_social_labour_minutes"`
	MeetsBurkeMinimumPlatoon      bool                  `json:"meets_burke_minimum_platoon"`
}

func (h *Handler) ComputeAverageSocialLabour(w http.ResponseWriter, r *http.Request) {
	id := agent.CooperationID(r.PathValue("id"))
	c, err := h.CooperationStore.GetCooperation(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	n := c.Size()
	d := c.AverageIndividualWorkingDay()
	collective := agent.CollectiveWorkingDay(n, d)
	writeJSON(w, http.StatusOK, averageSocialLabourResponse{
		CooperationID:                c.ID,
		Size:                         n,
		CollectiveWorkingDayMinutes:  collective,
		AverageSocialLabourMinutes:   agent.AverageSocialLabour(collective, n),
		MeetsBurkeMinimumPlatoon:     n >= agent.CooperationMinSize,
	})
}

type minimumCapitalRequest struct {
	WorkerCount    int64 `json:"worker_count"`
	DailyWagePence int64 `json:"daily_wage_pence"`
}

type minimumCapitalResponse struct {
	WorkerCount    agent.CooperationSize `json:"worker_count"`
	DailyWagePence agent.Pence           `json:"daily_wage_pence"`
	MinimumCapital agent.Pence           `json:"minimum_capital_pence"`
}

func (h *Handler) ComputeMinimumCapital(w http.ResponseWriter, r *http.Request) {
	var req minimumCapitalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.WorkerCount <= 0 {
		writeError(w, http.StatusBadRequest, "worker_count must be positive")
		return
	}
	if req.DailyWagePence <= 0 {
		writeError(w, http.StatusBadRequest, "daily_wage_pence must be positive")
		return
	}
	n := agent.CooperationSize(req.WorkerCount)
	wage := agent.Pence(req.DailyWagePence)
	writeJSON(w, http.StatusOK, minimumCapitalResponse{
		WorkerCount:    n,
		DailyWagePence: wage,
		MinimumCapital: agent.MinimumCapital(n, wage),
	})
}
