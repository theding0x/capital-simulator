package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

type runLabourProcessRequest struct {
	WorkerID          agent.AgentID           `json:"worker_id"`
	CapitalistID      agent.AgentID           `json:"capitalist_id"`
	MeansOfProduction agent.MeansOfProduction `json:"means_of_production"`
	DurationMinutes   agent.LabourMinutes     `json:"duration_minutes"`
	ProductKind       string                  `json:"product_kind"`
	ProductQuantity   int64                   `json:"product_quantity"`
}

type valorizationSummary struct {
	NecessaryLabour agent.LabourMinutes `json:"necessary_labour"`
	SurplusLabour   agent.LabourMinutes `json:"surplus_labour"`
	SurplusValue    agent.LabourMinutes `json:"surplus_value"`
	ProductValue    agent.LabourMinutes `json:"product_value"`
}

type labourProcessResponse struct {
	LabourProcess agent.LabourProcess `json:"labour_process"`
	Product       agent.Product       `json:"product"`
	Valorization  valorizationSummary `json:"valorization"`
}

func buildLabourProcessResponse(lp agent.LabourProcess) labourProcessResponse {
	vp := agent.ValorizationProcess{
		Process:                lp,
		NecessaryLabourMinutes: lp.NecessaryLabourMinutes,
	}
	product := agent.Product{
		CommodityKind: lp.ProductKind,
		Quantity:      lp.ProductQuantity,
		TotalValue:    vp.ProductValue(),
	}
	return labourProcessResponse{
		LabourProcess: lp,
		Product:       product,
		Valorization: valorizationSummary{
			NecessaryLabour: vp.NecessaryLabour(),
			SurplusLabour:   vp.SurplusLabour(),
			SurplusValue:    vp.SurplusValue(),
			ProductValue:    vp.ProductValue(),
		},
	}
}

func (h *Handler) validateProcessParties(r *http.Request, workerID, capitalistID agent.AgentID) (agent.Worker, error) {
	worker, err := h.LabourPowerStore.GetWorker(r.Context(), workerID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return agent.Worker{}, agent.ErrInvalidProcess
		}
		return agent.Worker{}, err
	}
	if _, err := h.LabourPowerStore.GetCapitalist(r.Context(), capitalistID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return agent.Worker{}, agent.ErrInvalidProcess
		}
		return agent.Worker{}, err
	}
	return worker, nil
}

func (h *Handler) RunLabourProcess(w http.ResponseWriter, r *http.Request) {
	var req runLabourProcessRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	worker, err := h.validateProcessParties(r, req.WorkerID, req.CapitalistID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	lp := agent.LabourProcess{
		WorkerID:               req.WorkerID,
		CapitalistID:           req.CapitalistID,
		Means:                  req.MeansOfProduction,
		Duration:               req.DurationMinutes,
		NecessaryLabourMinutes: worker.LabourPowerValueMinutes,
		ProductKind:            req.ProductKind,
		ProductQuantity:        req.ProductQuantity,
	}
	if err := lp.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.LabourProcessStore.CreateLabourProcess(r.Context(), lp)
	if err != nil {
		writeAppError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, buildLabourProcessResponse(saved))
}

func (h *Handler) GetLabourProcessRecord(w http.ResponseWriter, r *http.Request) {
	id := agent.LabourProcessID(r.PathValue("id"))
	lp, err := h.LabourProcessStore.GetLabourProcess(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildLabourProcessResponse(lp))
}
