package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

// AgentStore is the composite store interface every agent-service handler
// reaches into. Combining the per-chapter stores into one composite lets
// New take a single argument and the test harness pass `store.NewMemory()`
// once instead of seven times.
type AgentStore interface {
	store.Store
	store.CircuitStore
	store.LabourPowerStore
	store.LabourProcessStore
	store.WorkingDayStore
	store.CooperationStore
	store.ManufactureStore
}

type Handler struct {
	Store              store.Store
	CircuitStore       store.CircuitStore
	LabourPowerStore   store.LabourPowerStore
	LabourProcessStore store.LabourProcessStore
	WorkingDayStore    store.WorkingDayStore
	CooperationStore   store.CooperationStore
	ManufactureStore   store.ManufactureStore
	Logger             *slog.Logger
}

// New constructs a Handler. AgentStore is a composite interface that the
// in-memory and MySQL stores both satisfy.
func New(s AgentStore, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		Store:              s,
		CircuitStore:       s,
		LabourPowerStore:   s,
		LabourProcessStore: s,
		WorkingDayStore:    s,
		CooperationStore:   s,
		ManufactureStore:   s,
		Logger:             logger,
	}
}

type createAgentRequest struct {
	Name          string      `json:"name"`
	Class         agent.Class `json:"class"`
	MoneyBalance  agent.Pence `json:"money_balance"`
	LabourMinutes int64       `json:"labour_minutes"`
}

type updateAgentRequest struct {
	Name          *string      `json:"name,omitempty"`
	MoneyBalance  *agent.Pence `json:"money_balance,omitempty"`
	LabourMinutes *int64       `json:"labour_minutes,omitempty"`
}

type createCircuitRequest struct {
	MAdvanced   agent.Pence       `json:"m_advanced"`
	CommodityID string            `json:"commodity_id"`
	MReturned   agent.Pence       `json:"m_returned"`
	CircuitType agent.CircuitType `json:"circuit_type"`
}

type reinvestRequest struct {
	CommodityID string      `json:"commodity_id"`
	MReturned   agent.Pence `json:"m_returned"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAgentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a := agent.Agent{
		Name:          strings.TrimSpace(req.Name),
		Class:         req.Class,
		MoneyBalance:  req.MoneyBalance,
		LabourMinutes: req.LabourMinutes,
	}
	if err := a.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.Create(r.Context(), a)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	classParam := r.URL.Query().Get("class")
	var (
		agents []agent.Agent
		err    error
	)
	if classParam != "" {
		agents, err = h.Store.ListByClass(ctx, agent.Class(classParam))
	} else {
		agents, err = h.Store.List(ctx)
	}
	if err != nil {
		writeAppError(w, err)
		return
	}
	if agents == nil {
		agents = []agent.Agent{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": agents})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := agent.ID(r.PathValue("id"))
	a, err := h.Store.Get(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := agent.ID(r.PathValue("id"))
	var req updateAgentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}
	updated, err := h.Store.Update(r.Context(), id, store.Update{
		Name:          req.Name,
		MoneyBalance:  req.MoneyBalance,
		LabourMinutes: req.LabourMinutes,
	})
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := agent.ID(r.PathValue("id"))
	if err := h.Store.Delete(r.Context(), id); err != nil {
		writeAppError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateCircuit(w http.ResponseWriter, r *http.Request) {
	agentID := agent.ID(r.PathValue("id"))
	var req createCircuitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a, err := h.Store.Get(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	if (a.Class == agent.WorkerClass || a.Class == agent.Owner) && req.CircuitType == agent.CircuitMCM {
		writeError(w, http.StatusBadRequest, agent.ErrWrongClass.Error())
		return
	}
	c := agent.CapitalCircuit{
		AgentID:      agentID,
		MAdvanced:    req.MAdvanced,
		CommodityID:  req.CommodityID,
		MReturned:    req.MReturned,
		SurplusValue: req.MReturned - req.MAdvanced,
		CircuitType:  req.CircuitType,
	}
	if err := c.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.CircuitStore.CreateCircuit(r.Context(), c)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}

func (h *Handler) ListCircuits(w http.ResponseWriter, r *http.Request) {
	agentID := agent.ID(r.PathValue("id"))
	circuits, err := h.CircuitStore.ListCircuits(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	if circuits == nil {
		circuits = []agent.CapitalCircuit{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": circuits})
}

func (h *Handler) Reinvest(w http.ResponseWriter, r *http.Request) {
	agentID := agent.ID(r.PathValue("id"))
	var req reinvestRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a, err := h.Store.Get(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	circuit, _, err := a.Reinvest(req.CommodityID, req.MReturned)
	if err != nil {
		writeAppError(w, err)
		return
	}
	saved, err := h.CircuitStore.CreateCircuit(r.Context(), circuit)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}

func (h *Handler) Hoard(w http.ResponseWriter, r *http.Request) {
	agentID := agent.ID(r.PathValue("id"))
	a, err := h.Store.Get(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	if _, err := a.Hoard(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	hoarding := true
	saved, err := h.Store.Update(r.Context(), agentID, store.Update{Hoarding: &hoarding})
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, saved)
}

type computeCircuitRequest struct {
	M           agent.Pence `json:"m"`
	CommodityID string      `json:"commodity_id"`
	MPrime      agent.Pence `json:"m_prime"`
}

type computeCircuitResponse struct {
	M            agent.Pence `json:"m"`
	CommodityID  string      `json:"commodity_id,omitempty"`
	MPrime       agent.Pence `json:"m_prime"`
	SurplusValue agent.Pence `json:"surplus_value"`
	Origin       string      `json:"origin"`
}

type computeExchangeRequest struct {
	AValue agent.Pence `json:"a_value"`
	BValue agent.Pence `json:"b_value"`
}

func (h *Handler) ComputeCircuit(w http.ResponseWriter, r *http.Request) {
	var req computeCircuitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.M <= 0 {
		writeError(w, http.StatusBadRequest, "m must be positive")
		return
	}
	if req.MPrime < 0 {
		writeError(w, http.StatusBadRequest, "m_prime cannot be negative")
		return
	}
	var surplusValue agent.Pence
	var origin string
	if req.CommodityID != "" {
		mc := agent.MerchantsCapital{M: req.M, CommodityID: req.CommodityID, MPrime: req.MPrime}
		surplusValue = mc.SurplusValue()
		origin = mc.Origin()
	} else {
		uc := agent.UsurersCapital{M: req.M, MPrime: req.MPrime}
		surplusValue = uc.SurplusValue()
		origin = uc.Origin()
	}
	writeJSON(w, http.StatusOK, computeCircuitResponse{
		M:            req.M,
		CommodityID:  req.CommodityID,
		MPrime:       req.MPrime,
		SurplusValue: surplusValue,
		Origin:       origin,
	})
}

func (h *Handler) ComputeExchange(w http.ResponseWriter, r *http.Request) {
	var req computeExchangeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.AValue < 0 || req.BValue < 0 {
		writeError(w, http.StatusBadRequest, "values cannot be negative")
		return
	}
	var result agent.ExchangeResult
	if req.AValue == req.BValue {
		result = agent.ExchangeEquivalents(req.AValue, req.BValue)
	} else {
		result = agent.ExchangeNonEquivalents(req.AValue, req.BValue)
	}
	writeJSON(w, http.StatusOK, result)
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return errors.New("invalid json: " + err.Error())
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeAppError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, store.ErrAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, agent.ErrInsufficientFunds),
		errors.Is(err, agent.ErrNotCapitalist),
		errors.Is(err, agent.ErrWrongClass),
		errors.Is(err, agent.ErrInvalidProcess),
		errors.Is(err, agent.ErrInvalidContract),
		errors.Is(err, agent.ErrWorkingDayExceedsPhysicalMax),
		errors.Is(err, agent.ErrWorkingDayExceedsStatutoryLimit),
		errors.Is(err, agent.ErrNecessaryLabourRequired),
		errors.Is(err, agent.ErrCooperationCapitalistID),
		errors.Is(err, agent.ErrCooperationNoMembers),
		errors.Is(err, agent.ErrCooperationMemberInvalid),
		errors.Is(err, agent.ErrCooperationWorkingDay),
		errors.Is(err, agent.ErrCooperationSizeTooSmall),
		errors.Is(err, agent.ErrManufactureCapitalistID),
		errors.Is(err, agent.ErrManufactureForm),
		errors.Is(err, agent.ErrManufactureOrigin),
		errors.Is(err, agent.ErrManufactureWorkingDay),
		errors.Is(err, agent.ErrManufactureNoRoles),
		errors.Is(err, agent.ErrManufactureRoleName),
		errors.Is(err, agent.ErrManufactureSkillLevel),
		errors.Is(err, agent.ErrManufactureOutputRate),
		errors.Is(err, agent.ErrManufactureHeadCount),
		errors.Is(err, agent.ErrBottleneck),
		errors.Is(err, agent.ErrInvalidMultiplier):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
