package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

type detailRoleInput struct {
	Name              string `json:"name"`
	SkillLevel        string `json:"skill_level"`
	OutputRatePerHour int64  `json:"output_rate_per_hour"`
	HeadCount         int    `json:"head_count"`
	ToolName          string `json:"tool_name"`
}

type createManufactureRequest struct {
	Name                        string            `json:"name"`
	CapitalistID                string            `json:"capitalist_id"`
	Form                        string            `json:"form"`
	Origin                      string            `json:"origin"`
	IndividualWorkingDayMinutes int64             `json:"individual_working_day_minutes"`
	Roles                       []detailRoleInput `json:"roles"`
}

type manufactureResponse struct {
	agent.Manufacture
	TotalWorkers           int                 `json:"total_workers"`
	CollectiveLabourer     collectiveLabourer  `json:"collective_labourer"`
	Hierarchy              []agent.DetailRole  `json:"hierarchy"`
	ProductivePowerFactor  float64             `json:"productive_power_factor"`
	ProductivityFactor     float64             `json:"productivity_factor"`
	ProductivePowerMinutes agent.LabourMinutes `json:"productive_power_minutes"`
	CooperationBaseline    agent.LabourMinutes `json:"cooperation_baseline_minutes"`
}

type collectiveLabourer struct {
	Roles                  []agent.DetailRole  `json:"roles"`
	TotalWorkers           int                 `json:"total_workers"`
	OutputPerPeriodMinutes agent.LabourMinutes `json:"output_per_period_minutes"`
	PeriodMinutes          int64               `json:"period_minutes"`
	IsParalysedIfAbsentOne bool                `json:"is_paralysed_if_absent_one"`
}

func buildManufactureResponse(m agent.Manufacture) manufactureResponse {
	if m.Roles == nil {
		m.Roles = []agent.DetailRole{}
	}
	cl := m.CollectiveLabourer()
	period := int64(m.IndividualWorkingDayMinutes)
	if period <= 0 {
		period = 480
	}
	hierarchy := []agent.DetailRole(m.Hierarchy())
	factor := agent.ManufactureProductivePowerFactor(m)
	return manufactureResponse{
		Manufacture:  m,
		TotalWorkers: m.TotalWorkers(),
		CollectiveLabourer: collectiveLabourer{
			Roles:                  cl.Roles,
			TotalWorkers:           cl.TotalWorkers(),
			OutputPerPeriodMinutes: cl.OutputPerPeriod(period),
			PeriodMinutes:          period,
			IsParalysedIfAbsentOne: cl.IsParalysed(1),
		},
		Hierarchy:              hierarchy,
		ProductivePowerFactor:  factor,
		ProductivityFactor:     factor, // Part IV bridge: feed into Ch. 12.
		ProductivePowerMinutes: agent.ManufactureProductivePower(m, period),
		CooperationBaseline: agent.CollectiveWorkingDay(
			agent.CooperationSize(m.TotalWorkers()), m.IndividualWorkingDayMinutes),
	}
}

func parseManufactureInput(req createManufactureRequest) (agent.Manufacture, error) {
	roles := make([]agent.DetailRole, 0, len(req.Roles))
	for _, r := range req.Roles {
		roles = append(roles, agent.DetailRole{
			Name:              strings.TrimSpace(r.Name),
			SkillLevel:        agent.SkillLevel(r.SkillLevel),
			OutputRatePerHour: r.OutputRatePerHour,
			HeadCount:         r.HeadCount,
			ToolName:          strings.TrimSpace(r.ToolName),
		})
	}
	m := agent.Manufacture{
		Name:                        strings.TrimSpace(req.Name),
		CapitalistID:                agent.AgentID(req.CapitalistID),
		Form:                        agent.ManufactureForm(req.Form),
		Origin:                      agent.ManufactureOrigin(req.Origin),
		IndividualWorkingDayMinutes: agent.LabourMinutes(req.IndividualWorkingDayMinutes),
		Roles:                       roles,
	}
	if err := m.Validate(); err != nil {
		return agent.Manufacture{}, err
	}
	return m, nil
}

func (h *Handler) CreateManufacture(w http.ResponseWriter, r *http.Request) {
	var req createManufactureRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	m, err := parseManufactureInput(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.ManufactureStore.CreateManufacture(r.Context(), m)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, buildManufactureResponse(saved))
}

func (h *Handler) GetManufacture(w http.ResponseWriter, r *http.Request) {
	id := agent.ManufactureID(r.PathValue("id"))
	m, err := h.ManufactureStore.GetManufacture(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildManufactureResponse(m))
}

func (h *Handler) ListManufactures(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capParam := r.URL.Query().Get("capitalist_id")
	formParam := r.URL.Query().Get("form")
	var (
		items []agent.Manufacture
		err   error
	)
	switch {
	case capParam != "":
		items, err = h.ManufactureStore.ListManufacturesByCapitalist(ctx, agent.AgentID(capParam))
	case formParam != "":
		items, err = h.ManufactureStore.ListManufacturesByForm(ctx, agent.ManufactureForm(formParam))
	default:
		items, err = h.ManufactureStore.ListManufactures(ctx)
	}
	if err != nil {
		writeAppError(w, err)
		return
	}
	out := make([]manufactureResponse, 0, len(items))
	for _, m := range items {
		out = append(out, buildManufactureResponse(m))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

type proportionalGroupSizeRequest struct {
	TargetOutputRate int64 `json:"target_output_rate"`
}

type proportionalGroupSizeResponse struct {
	ManufactureID    agent.ManufactureID `json:"manufacture_id"`
	TargetOutputRate int64               `json:"target_output_rate"`
	Headcounts       map[string]int      `json:"headcounts"`
}

func (h *Handler) ProportionalGroupSize(w http.ResponseWriter, r *http.Request) {
	id := agent.ManufactureID(r.PathValue("id"))
	var req proportionalGroupSizeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	m, err := h.ManufactureStore.GetManufacture(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	hc, err := agent.ProportionalGroupSize(m.Roles, req.TargetOutputRate)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, proportionalGroupSizeResponse{
		ManufactureID:    m.ID,
		TargetOutputRate: req.TargetOutputRate,
		Headcounts:       hc,
	})
}

type scaleManufactureRequest struct {
	Multiplier int `json:"multiplier"`
}

func (h *Handler) ScaleManufacture(w http.ResponseWriter, r *http.Request) {
	id := agent.ManufactureID(r.PathValue("id"))
	var req scaleManufactureRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	m, err := h.ManufactureStore.GetManufacture(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	scaled, err := agent.ScaleManufacture(m, req.Multiplier)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildManufactureResponse(scaled))
}

type collectiveLabourerResponse struct {
	ManufactureID          agent.ManufactureID `json:"manufacture_id"`
	TotalWorkers           int                 `json:"total_workers"`
	OutputPerPeriodMinutes agent.LabourMinutes `json:"output_per_period_minutes"`
	PeriodMinutes          int64               `json:"period_minutes"`
	IsParalysedIfAbsentOne bool                `json:"is_paralysed_if_absent_one"`
	Roles                  []agent.DetailRole  `json:"roles"`
}

func (h *Handler) GetCollectiveLabourer(w http.ResponseWriter, r *http.Request) {
	id := agent.ManufactureID(r.PathValue("id"))
	m, err := h.ManufactureStore.GetManufacture(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	cl := m.CollectiveLabourer()
	period := int64(480)
	writeJSON(w, http.StatusOK, collectiveLabourerResponse{
		ManufactureID:          m.ID,
		TotalWorkers:           cl.TotalWorkers(),
		OutputPerPeriodMinutes: cl.OutputPerPeriod(period),
		PeriodMinutes:          period,
		IsParalysedIfAbsentOne: cl.IsParalysed(1),
		Roles:                  cl.Roles,
	})
}

type minimumCapitalManufactureResponse struct {
	ManufactureID         agent.ManufactureID `json:"manufacture_id"`
	RawMaterialCostFactor float64             `json:"raw_material_cost_factor"`
	MinimumCapitalPence   agent.Pence         `json:"minimum_capital_pence"`
	CooperationBaseline   agent.Pence         `json:"cooperation_baseline_pence"`
}

func (h *Handler) GetManufactureMinimumCapital(w http.ResponseWriter, r *http.Request) {
	id := agent.ManufactureID(r.PathValue("id"))
	factorStr := r.URL.Query().Get("raw_material_cost_factor")
	factor := 0.0
	if factorStr != "" {
		v, err := strconv.ParseFloat(factorStr, 64)
		if err != nil || v < 0 {
			writeError(w, http.StatusBadRequest, "raw_material_cost_factor must be a non-negative number")
			return
		}
		factor = v
	}
	m, err := h.ManufactureStore.GetManufacture(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	coopBaseline := agent.MinimumCapital(
		agent.CooperationSize(m.TotalWorkers()),
		agent.SkilledDailyWagePence,
	)
	writeJSON(w, http.StatusOK, minimumCapitalManufactureResponse{
		ManufactureID:         m.ID,
		RawMaterialCostFactor: factor,
		MinimumCapitalPence:   agent.ManufactureMinimumCapital(m, factor),
		CooperationBaseline:   coopBaseline,
	})
}
