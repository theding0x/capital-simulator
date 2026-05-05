package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

type createWorkerRequest struct {
	OwnsLabourPower         bool                `json:"owns_labour_power"`
	OwnsCommoditiesToSell   bool                `json:"owns_commodities_to_sell"`
	CapacityMinutesPerDay   agent.LabourMinutes `json:"capacity_minutes_per_day"`
	LabourPowerValueMinutes agent.LabourMinutes `json:"labour_power_value_minutes"`
}

type createCapitalistRequest struct {
	MoneyCapital agent.LabourMinutes `json:"money_capital"`
}

type createOfferingRequest struct {
	OwnerID               agent.AgentID       `json:"owner_id"`
	CapacityMinutesPerDay agent.LabourMinutes `json:"capacity_minutes_per_day"`
	ContractDays          int64               `json:"contract_days"`
	AskingWage            agent.LabourMinutes `json:"asking_wage"`
}

type createPurchaseRequest struct {
	SellerID     agent.AgentID       `json:"seller_id"`
	BuyerID      agent.AgentID       `json:"buyer_id"`
	WageMinutes  agent.LabourMinutes `json:"wage_minutes"`
	ContractDays int64               `json:"contract_days"`
}

func (h *Handler) CreateWorker(w http.ResponseWriter, r *http.Request) {
	var req createWorkerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	worker := agent.Worker{
		OwnsLabourPower:         req.OwnsLabourPower,
		OwnsCommoditiesToSell:   req.OwnsCommoditiesToSell,
		LabourPower:             agent.LabourPower{CapacityMinutesPerDay: req.CapacityMinutesPerDay},
		LabourPowerValueMinutes: req.LabourPowerValueMinutes,
	}
	if err := worker.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.LabourPowerStore.CreateWorker(r.Context(), worker)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) GetWorker(w http.ResponseWriter, r *http.Request) {
	id := agent.AgentID(r.PathValue("id"))
	worker, err := h.LabourPowerStore.GetWorker(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, worker)
}

func (h *Handler) ListWorkers(w http.ResponseWriter, r *http.Request) {
	workers, err := h.LabourPowerStore.ListWorkers(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	if workers == nil {
		workers = []agent.Worker{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": workers})
}

func (h *Handler) CreateCapitalist(w http.ResponseWriter, r *http.Request) {
	var req createCapitalistRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	cap := agent.Capitalist{
		MoneyCapital: req.MoneyCapital,
	}
	if err := cap.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.LabourPowerStore.CreateCapitalist(r.Context(), cap)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) GetCapitalist(w http.ResponseWriter, r *http.Request) {
	id := agent.AgentID(r.PathValue("id"))
	cap, err := h.LabourPowerStore.GetCapitalist(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cap)
}

func (h *Handler) ListCapitalists(w http.ResponseWriter, r *http.Request) {
	caps, err := h.LabourPowerStore.ListCapitalists(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	if caps == nil {
		caps = []agent.Capitalist{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": caps})
}

func (h *Handler) CreateOffering(w http.ResponseWriter, r *http.Request) {
	var req createOfferingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	offering := agent.LabourPowerOffering{
		OwnerID:               req.OwnerID,
		CapacityMinutesPerDay: req.CapacityMinutesPerDay,
		ContractDays:          req.ContractDays,
		AskingWage:            req.AskingWage,
	}
	if err := offering.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.LabourPowerStore.CreateOffering(r.Context(), offering)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) ListOfferings(w http.ResponseWriter, r *http.Request) {
	offerings, err := h.LabourPowerStore.ListOfferings(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	if offerings == nil {
		offerings = []agent.LabourPowerOffering{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": offerings})
}

func (h *Handler) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	var req createPurchaseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	purchase := agent.LabourPowerPurchase{
		SellerID:     req.SellerID,
		BuyerID:      req.BuyerID,
		WageMinutes:  req.WageMinutes,
		ContractDays: req.ContractDays,
	}
	if err := purchase.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validatePurchaseParties(r, req.SellerID, req.BuyerID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.LabourPowerStore.CreatePurchase(r.Context(), purchase)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) validatePurchaseParties(r *http.Request, sellerID, buyerID agent.AgentID) error {
	if _, err := h.LabourPowerStore.GetWorker(r.Context(), sellerID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return agent.ErrInvalidContract
		}
		return err
	}
	if _, err := h.LabourPowerStore.GetCapitalist(r.Context(), buyerID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return agent.ErrInvalidContract
		}
		return err
	}
	return nil
}

func (h *Handler) GetPurchase(w http.ResponseWriter, r *http.Request) {
	id := agent.PurchaseID(r.PathValue("id"))
	purchase, err := h.LabourPowerStore.GetPurchase(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, purchase)
}

func (h *Handler) ListPurchases(w http.ResponseWriter, r *http.Request) {
	purchases, err := h.LabourPowerStore.ListPurchases(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	if purchases == nil {
		purchases = []agent.LabourPowerPurchase{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": purchases})
}
