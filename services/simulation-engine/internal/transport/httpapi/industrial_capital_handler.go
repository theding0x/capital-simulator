package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Vol. II Ch. 4 — The Three Formulas of the Circuit.
// IndustrialCapital is the synthesis: every industrial capital occupies all
// three stages (M, P, C) simultaneously — Marx's Nebeneinander.

// ----- response DTOs ---------------------------------------------------------

type industrialCapitalResponse struct {
	ID                       string                `json:"id"`
	AgentID                  string                `json:"agent_id,omitempty"`
	MoneyCircuitID           string                `json:"money_circuit_id,omitempty"`
	ProductiveCircuitID      string                `json:"productive_circuit_id,omitempty"`
	CommodityCircuitID       string                `json:"commodity_circuit_id,omitempty"`
	TotalPence               int64                 `json:"total_pence"`
	EconomyMode              string                `json:"economy_mode"`
	StagnationToleranceTicks int64                 `json:"stagnation_tolerance_ticks"`
	Status                   string                `json:"status"`
	LatestDistribution       *stageDistributionDTO `json:"latest_distribution,omitempty"`
	OpenBlocks               []stageBlockDTO       `json:"open_blocks"`
	CreatedAt                string                `json:"created_at"`
	UpdatedAt                string                `json:"updated_at"`
}

type stageDistributionDTO struct {
	ID                  string `json:"id"`
	IndustrialCapitalID string `json:"industrial_capital_id"`
	At                  string `json:"at"`
	MoneyPence          int64  `json:"money_pence"`
	ProductionPence     int64  `json:"production_pence"`
	CommodityPence      int64  `json:"commodity_pence"`
}

type stageBlockDTO struct {
	ID                  string  `json:"id"`
	IndustrialCapitalID string  `json:"industrial_capital_id"`
	Stage               string  `json:"stage"`
	Reason              string  `json:"reason"`
	OpenedAt            string  `json:"opened_at"`
	ClosedAt            *string `json:"closed_at,omitempty"`
}

type capitalPartDTO struct {
	ID                  string `json:"id"`
	IndustrialCapitalID string `json:"industrial_capital_id"`
	Pence               int64  `json:"pence"`
	Stage               string `json:"stage"`
	EnteredStageAt      string `json:"entered_stage_at"`
}

type valueRevolutionDTO struct {
	Event   valueRevolutionEventDTO `json:"event"`
	SetFree *moneyCapitalSetFreeDTO `json:"set_free,omitempty"`
	TiedUp  *moneyCapitalTiedUpDTO  `json:"tied_up,omitempty"`
}

type valueRevolutionEventDTO struct {
	ID                  string `json:"id"`
	IndustrialCapitalID string `json:"industrial_capital_id"`
	Affecting           string `json:"affecting"`
	DirectionPence      int64  `json:"direction_pence"`
	OccurredAt          string `json:"occurred_at"`
}

type moneyCapitalSetFreeDTO struct {
	ID                     string `json:"id"`
	ValueRevolutionEventID string `json:"value_revolution_event_id"`
	IndustrialCapitalID    string `json:"industrial_capital_id"`
	AmountPence            int64  `json:"amount_pence"`
}

type moneyCapitalTiedUpDTO struct {
	ID                     string `json:"id"`
	ValueRevolutionEventID string `json:"value_revolution_event_id"`
	IndustrialCapitalID    string `json:"industrial_capital_id"`
	AmountPence            int64  `json:"amount_pence"`
}

type metamorphosisInterlockDTO struct {
	ID                        string `json:"id"`
	BuyerIndustrialCapitalID  string `json:"buyer_industrial_capital_id"`
	SellerIndustrialCapitalID string `json:"seller_industrial_capital_id"`
	SellerMPSourceOrigin      string `json:"seller_mp_source_origin"`
	Pence                     int64  `json:"pence"`
	OccurredAt                string `json:"occurred_at"`
}

type supplyDemandImbalanceDTO struct {
	ID                  string `json:"id"`
	IndustrialCapitalID string `json:"industrial_capital_id"`
	Period              string `json:"period"`
	DemandPence         int64  `json:"demand_pence"`
	SupplyPence         int64  `json:"supply_pence"`
	ExcessPence         int64  `json:"excess_pence"`
}

type aggregateSupplyDemandDTO struct {
	Period      string `json:"period"`
	DemandPence int64  `json:"demand_pence"`
	SupplyPence int64  `json:"supply_pence"`
	ExcessPence int64  `json:"excess_pence"`
}

type sinkingFundDTO struct {
	ID                  string `json:"id"`
	IndustrialCapitalID string `json:"industrial_capital_id"`
	FixedCapitalPence   int64  `json:"fixed_capital_pence"`
	LifetimeYears       int64  `json:"lifetime_years"`
	AnnualPaymentPence  int64  `json:"annual_payment_pence"`
	AccumulatedPence    int64  `json:"accumulated_pence"`
}

// ----- to-DTO converters -----------------------------------------------------

func toIndustrialCapitalResponse(ic circulation.IndustrialCapital) industrialCapitalResponse {
	resp := industrialCapitalResponse{
		ID:                       string(ic.ID),
		AgentID:                  ic.AgentID,
		MoneyCircuitID:           string(ic.MoneyCircuitID),
		ProductiveCircuitID:      string(ic.ProductiveCircuitID),
		CommodityCircuitID:       string(ic.CommodityCircuitID),
		TotalPence:               int64(ic.TotalPence),
		EconomyMode:              string(ic.EconomyMode),
		StagnationToleranceTicks: ic.StagnationToleranceTicks,
		Status:                   string(ic.Status),
		OpenBlocks:               []stageBlockDTO{},
		CreatedAt:                ic.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:                ic.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if ic.Latest != nil {
		dto := toStageDistributionDTO(*ic.Latest)
		resp.LatestDistribution = &dto
	}
	for _, b := range ic.OpenBlocks {
		resp.OpenBlocks = append(resp.OpenBlocks, toStageBlockDTO(b))
	}
	return resp
}

func toStageDistributionDTO(sd circulation.StageDistribution) stageDistributionDTO {
	return stageDistributionDTO{
		ID:                  string(sd.ID),
		IndustrialCapitalID: string(sd.IndustrialCapitalID),
		At:                  sd.At.UTC().Format(time.RFC3339),
		MoneyPence:          int64(sd.MoneyPence),
		ProductionPence:     int64(sd.ProductionPence),
		CommodityPence:      int64(sd.CommodityPence),
	}
}

func toStageBlockDTO(b circulation.StageBlock) stageBlockDTO {
	dto := stageBlockDTO{
		ID:                  string(b.ID),
		IndustrialCapitalID: string(b.IndustrialCapitalID),
		Stage:               string(b.Stage),
		Reason:              string(b.Reason),
		OpenedAt:            b.OpenedAt.UTC().Format(time.RFC3339),
	}
	if b.ClosedAt != nil {
		s := b.ClosedAt.UTC().Format(time.RFC3339)
		dto.ClosedAt = &s
	}
	return dto
}

func toValueRevolutionDTO(res circulation.ValueRevolutionResult) valueRevolutionDTO {
	dto := valueRevolutionDTO{
		Event: valueRevolutionEventDTO{
			ID:                  string(res.Event.ID),
			IndustrialCapitalID: string(res.Event.IndustrialCapitalID),
			Affecting:           string(res.Event.Affecting),
			DirectionPence:      int64(res.Event.DirectionPence),
			OccurredAt:          res.Event.OccurredAt.UTC().Format(time.RFC3339),
		},
	}
	if res.SetFree != nil {
		dto.SetFree = &moneyCapitalSetFreeDTO{
			ID:                     string(res.SetFree.ID),
			ValueRevolutionEventID: string(res.SetFree.ValueRevolutionEventID),
			IndustrialCapitalID:    string(res.SetFree.IndustrialCapitalID),
			AmountPence:            int64(res.SetFree.AmountPence),
		}
	}
	if res.TiedUp != nil {
		dto.TiedUp = &moneyCapitalTiedUpDTO{
			ID:                     string(res.TiedUp.ID),
			ValueRevolutionEventID: string(res.TiedUp.ValueRevolutionEventID),
			IndustrialCapitalID:    string(res.TiedUp.IndustrialCapitalID),
			AmountPence:            int64(res.TiedUp.AmountPence),
		}
	}
	return dto
}

func toSinkingFundDTO(sf circulation.SinkingFund) sinkingFundDTO {
	return sinkingFundDTO{
		ID:                  string(sf.ID),
		IndustrialCapitalID: string(sf.IndustrialCapitalID),
		FixedCapitalPence:   int64(sf.FixedCapitalPence),
		LifetimeYears:       sf.LifetimeYears,
		AnnualPaymentPence:  int64(sf.AnnualPaymentPence),
		AccumulatedPence:    int64(sf.AccumulatedPence),
	}
}

// ----- handlers --------------------------------------------------------------

// CreateIndustrialCapital handles POST /v1/industrial-capitals.
func (h *Handler) CreateIndustrialCapital(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AgentID                  string `json:"agent_id,omitempty"`
		MoneyCircuitID           string `json:"money_circuit_id,omitempty"`
		ProductiveCircuitID      string `json:"productive_circuit_id,omitempty"`
		CommodityCircuitID       string `json:"commodity_circuit_id,omitempty"`
		TotalPence               int64  `json:"total_pence"`
		EconomyMode              string `json:"economy_mode"`
		StagnationToleranceTicks int64  `json:"stagnation_tolerance_ticks"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	ic := circulation.IndustrialCapital{
		AgentID:                  req.AgentID,
		MoneyCircuitID:           circulation.MoneyCircuitID(req.MoneyCircuitID),
		ProductiveCircuitID:      circulation.ProductiveCircuitID(req.ProductiveCircuitID),
		CommodityCircuitID:       circulation.CommodityCircuitID(req.CommodityCircuitID),
		TotalPence:               circulation.Pence(req.TotalPence),
		EconomyMode:              circulation.EconomyMode(req.EconomyMode),
		StagnationToleranceTicks: req.StagnationToleranceTicks,
	}
	created, err := h.IndustrialCapitals.CreateIndustrialCapital(r.Context(), ic)
	if err != nil {
		if isValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/industrial-capitals/"+string(created.ID))
	writeJSON(w, http.StatusCreated, toIndustrialCapitalResponse(created))
}

// GetIndustrialCapital handles GET /v1/industrial-capitals/{id}.
func (h *Handler) GetIndustrialCapital(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	ic, err := h.IndustrialCapitals.GetIndustrialCapital(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "industrial capital not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toIndustrialCapitalResponse(ic))
}

// ListIndustrialCapitals handles GET /v1/industrial-capitals.
func (h *Handler) ListIndustrialCapitals(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent_id")
	status := r.URL.Query().Get("status")
	economyMode := r.URL.Query().Get("economy_mode")
	list, err := h.IndustrialCapitals.ListIndustrialCapitals(r.Context(), agentID, status, economyMode)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	out := make([]industrialCapitalResponse, len(list))
	for i, ic := range list {
		out[i] = toIndustrialCapitalResponse(ic)
	}
	writeJSON(w, http.StatusOK, out)
}

// RecordCapitalPart handles POST /v1/industrial-capitals/{id}/parts.
func (h *Handler) RecordCapitalPart(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	var req struct {
		Pence          int64  `json:"pence"`
		Stage          string `json:"stage"`
		EnteredStageAt string `json:"entered_stage_at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Pence <= 0 {
		writeError(w, http.StatusBadRequest, "pence must be positive")
		return
	}
	part := circulation.CapitalPart{
		Pence: circulation.Pence(req.Pence),
		Stage: circulation.CapitalStage(req.Stage),
	}
	if req.EnteredStageAt != "" {
		t, err := time.Parse(time.RFC3339, req.EnteredStageAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid entered_stage_at: "+err.Error())
			return
		}
		part.EnteredStageAt = t
	}
	created, err := h.IndustrialCapitals.RecordCapitalPart(r.Context(), id, part)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "industrial capital not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, capitalPartDTO{
		ID:                  string(created.ID),
		IndustrialCapitalID: string(created.IndustrialCapitalID),
		Pence:               int64(created.Pence),
		Stage:               string(created.Stage),
		EnteredStageAt:      created.EnteredStageAt.UTC().Format(time.RFC3339),
	})
}

// SnapshotStageDistribution handles POST /v1/industrial-capitals/{id}/snapshot.
func (h *Handler) SnapshotStageDistribution(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	var req struct {
		MoneyPence      int64  `json:"money_pence"`
		ProductionPence int64  `json:"production_pence"`
		CommodityPence  int64  `json:"commodity_pence"`
		At              string `json:"at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sd := circulation.StageDistribution{
		MoneyPence:      circulation.Pence(req.MoneyPence),
		ProductionPence: circulation.Pence(req.ProductionPence),
		CommodityPence:  circulation.Pence(req.CommodityPence),
	}
	if req.At != "" {
		t, err := time.Parse(time.RFC3339, req.At)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid at: "+err.Error())
			return
		}
		sd.At = t
	}
	created, err := h.IndustrialCapitals.Snapshot(r.Context(), id, sd)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "industrial capital not found")
			return
		}
		if isValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toStageDistributionDTO(created))
}

// OpenStageBlock handles POST /v1/industrial-capitals/{id}/blocks.
func (h *Handler) OpenStageBlock(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	var req struct {
		Stage    string `json:"stage"`
		Reason   string `json:"reason"`
		OpenedAt string `json:"opened_at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	b := circulation.StageBlock{
		Stage:  circulation.CapitalStage(req.Stage),
		Reason: circulation.StageBlockReason(req.Reason),
	}
	if req.OpenedAt != "" {
		t, err := time.Parse(time.RFC3339, req.OpenedAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid opened_at: "+err.Error())
			return
		}
		b.OpenedAt = t
	}
	created, err := h.IndustrialCapitals.OpenBlock(r.Context(), id, b)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "industrial capital not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toStageBlockDTO(created))
}

// CloseStageBlock handles POST /v1/industrial-capitals/{id}/blocks/{block_id}/close.
func (h *Handler) CloseStageBlock(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	blockID := circulation.StageBlockID(r.PathValue("block_id"))
	closed, err := h.IndustrialCapitals.CloseBlock(r.Context(), id, blockID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "block not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toStageBlockDTO(closed))
}

// RecordValueRevolution handles POST /v1/industrial-capitals/{id}/value-revolution.
func (h *Handler) RecordValueRevolution(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	var req struct {
		Affecting      string `json:"affecting"`
		DirectionPence int64  `json:"direction_pence"`
		OccurredAt     string `json:"occurred_at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	evt := circulation.ValueRevolutionEvent{
		IndustrialCapitalID: id,
		Affecting:           circulation.ValueRevolutionAffecting(req.Affecting),
		DirectionPence:      circulation.Pence(req.DirectionPence),
	}
	if req.OccurredAt != "" {
		t, err := time.Parse(time.RFC3339, req.OccurredAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid occurred_at: "+err.Error())
			return
		}
		evt.OccurredAt = t
	}
	res, err := circulation.ApplyValueRevolution(evt)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.IndustrialCapitals.RecordValueRevolution(r.Context(), res)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "industrial capital not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toValueRevolutionDTO(saved))
}

// RecordInterlock handles POST /v1/industrial-capitals/{id}/interlocks.
func (h *Handler) RecordInterlock(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	var req struct {
		SellerIndustrialCapitalID string `json:"seller_industrial_capital_id"`
		SellerMPSourceOrigin      string `json:"seller_mp_source_origin"`
		Pence                     int64  `json:"pence"`
		OccurredAt                string `json:"occurred_at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	mi := circulation.MetamorphosisInterlock{
		SellerIndustrialCapitalID: circulation.IndustrialCapitalID(req.SellerIndustrialCapitalID),
		SellerMPSourceOrigin:      circulation.MPSourceOrigin(req.SellerMPSourceOrigin),
		Pence:                     circulation.Pence(req.Pence),
	}
	if req.OccurredAt != "" {
		t, err := time.Parse(time.RFC3339, req.OccurredAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid occurred_at: "+err.Error())
			return
		}
		mi.OccurredAt = t
	}
	saved, err := h.IndustrialCapitals.RecordInterlock(r.Context(), id, mi)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "industrial capital not found")
			return
		}
		if errors.Is(err, circulation.ErrNonCapitalCounterparty) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, metamorphosisInterlockDTO{
		ID:                        string(saved.ID),
		BuyerIndustrialCapitalID:  string(saved.BuyerIndustrialCapitalID),
		SellerIndustrialCapitalID: string(saved.SellerIndustrialCapitalID),
		SellerMPSourceOrigin:      string(saved.SellerMPSourceOrigin),
		Pence:                     int64(saved.Pence),
		OccurredAt:                saved.OccurredAt.UTC().Format(time.RFC3339),
	})
}

// ComputeAndRecordSupplyDemand handles POST /v1/industrial-capitals/{id}/supply-demand.
func (h *Handler) ComputeAndRecordSupplyDemand(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	var req struct {
		ConstantPence int64  `json:"constant_pence"`
		VariablePence int64  `json:"variable_pence"`
		SurplusPence  int64  `json:"surplus_pence"`
		Period        string `json:"period"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	oc := circulation.OrganicComposition{
		ConstantPence: circulation.Pence(req.ConstantPence),
		VariablePence: circulation.Pence(req.VariablePence),
	}
	sdi := circulation.ComputeSupplyDemandImbalance(id, req.Period, oc, circulation.Pence(req.SurplusPence))
	saved, err := h.IndustrialCapitals.RecordSupplyDemand(r.Context(), sdi)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "industrial capital not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, supplyDemandImbalanceDTO{
		ID:                  string(saved.ID),
		IndustrialCapitalID: string(saved.IndustrialCapitalID),
		Period:              saved.Period,
		DemandPence:         int64(saved.DemandPence),
		SupplyPence:         int64(saved.SupplyPence),
		ExcessPence:         int64(saved.ExcessPence),
	})
}

// GetSupplyDemand handles GET /v1/industrial-capitals/{id}/supply-demand.
func (h *Handler) GetSupplyDemand(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	period := r.URL.Query().Get("period")
	if period == "" {
		writeError(w, http.StatusBadRequest, "period query parameter is required")
		return
	}
	sdi, err := h.IndustrialCapitals.GetSupplyDemand(r.Context(), id, period)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "supply demand imbalance not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, supplyDemandImbalanceDTO{
		ID:                  string(sdi.ID),
		IndustrialCapitalID: string(sdi.IndustrialCapitalID),
		Period:              sdi.Period,
		DemandPence:         int64(sdi.DemandPence),
		SupplyPence:         int64(sdi.SupplyPence),
		ExcessPence:         int64(sdi.ExcessPence),
	})
}

// AggregateSupplyDemand handles GET /v1/industrial-capitals/supply-demand/aggregate.
func (h *Handler) AggregateSupplyDemand(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		writeError(w, http.StatusBadRequest, "period query parameter is required")
		return
	}
	agg, err := h.IndustrialCapitals.AggregateSupplyDemand(r.Context(), period)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, aggregateSupplyDemandDTO{
		Period:      agg.Period,
		DemandPence: int64(agg.DemandPence),
		SupplyPence: int64(agg.SupplyPence),
		ExcessPence: int64(agg.ExcessPence),
	})
}

// SetSinkingFund handles POST /v1/industrial-capitals/{id}/sinking-fund.
func (h *Handler) SetSinkingFund(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	var req struct {
		FixedCapitalPence  int64 `json:"fixed_capital_pence"`
		LifetimeYears      int64 `json:"lifetime_years"`
		AnnualPaymentPence int64 `json:"annual_payment_pence"`
		AccumulatedPence   int64 `json:"accumulated_pence"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	sf := circulation.SinkingFund{
		FixedCapitalPence:  circulation.Pence(req.FixedCapitalPence),
		LifetimeYears:      req.LifetimeYears,
		AnnualPaymentPence: circulation.Pence(req.AnnualPaymentPence),
		AccumulatedPence:   circulation.Pence(req.AccumulatedPence),
	}
	saved, err := h.IndustrialCapitals.SetSinkingFund(r.Context(), id, sf)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "industrial capital not found")
			return
		}
		if isValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toSinkingFundDTO(saved))
}

// TickSinkingFund handles POST /v1/industrial-capitals/{id}/sinking-fund/tick.
func (h *Handler) TickSinkingFund(w http.ResponseWriter, r *http.Request) {
	id := circulation.IndustrialCapitalID(r.PathValue("id"))
	updated, err := h.IndustrialCapitals.TickSinkingFund(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "sinking fund not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toSinkingFundDTO(updated))
}

