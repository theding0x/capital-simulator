package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Vol. II Ch. 15 — The Effects of a Change of Prices.

// --- request types ---

type createPriceRevolutionRequest struct {
	IndustrialCapitalID      string `json:"industrial_capital_id"`
	Affecting                string `json:"affecting"`
	DirectionPence           int64  `json:"direction_pence"`
	OccurredAt               string `json:"occurred_at,omitempty"`
	BasisPointsChange        int64  `json:"basis_points_change"`
	OldAdvancePence          *int64 `json:"old_advance_pence,omitempty"`
	ExpectedRealisationPence *int64 `json:"expected_realisation_pence,omitempty"`
	ActualRealisationPence   *int64 `json:"actual_realisation_pence,omitempty"`
}

type recordPriceCaseRequest struct {
	FallCase *string `json:"fall_case,omitempty"`
	RiseCase *string `json:"rise_case,omitempty"`
}

type recordInventoryRevaluationRequest struct {
	IndustrialCapitalID string `json:"industrial_capital_id"`
	CommodityID         string `json:"commodity_id"`
	QuantityHeld        int64  `json:"quantity_held"`
	OldUnitPence        int64  `json:"old_unit_pence"`
	NewUnitPence        int64  `json:"new_unit_pence"`
}

type recordSpeculativeHoldRequest struct {
	IndustrialCapitalID  string `json:"industrial_capital_id"`
	CommodityID          string `json:"commodity_id"`
	HeldSince            string `json:"held_since,omitempty"`
	AnticipatedDirection int8   `json:"anticipated_direction"`
}

// --- response types ---

type valueRevolutionEventV2Response struct {
	ID                   string `json:"id"`
	IndustrialCapitalID  string `json:"industrial_capital_id"`
	Affecting            string `json:"affecting"`
	DirectionPence       int64  `json:"direction_pence"`
	OccurredAt           string `json:"occurred_at"`
	BasisPointsChange    int64  `json:"basis_points_change"`
	RevaluedCapitalPence int64  `json:"revalued_capital_pence"`
}

type revaluedCapitalResponse struct {
	ID                     string `json:"id"`
	ValueRevolutionEventID string `json:"value_revolution_event_id"`
	IndustrialCapitalID    string `json:"industrial_capital_id"`
	OldAdvancePence        int64  `json:"old_advance_pence"`
	NewAdvancePence        int64  `json:"new_advance_pence"`
	DeltaPence             int64  `json:"delta_pence"`
}

type realisedValueDeltaResponse struct {
	ID                       string `json:"id"`
	ValueRevolutionEventID   string `json:"value_revolution_event_id"`
	IndustrialCapitalID      string `json:"industrial_capital_id"`
	ExpectedRealisationPence int64  `json:"expected_realisation_pence"`
	ActualRealisationPence   int64  `json:"actual_realisation_pence"`
	DeltaPence               int64  `json:"delta_pence"`
}

type inventoryRevaluationResponse struct {
	ID                  string `json:"id"`
	IndustrialCapitalID string `json:"industrial_capital_id"`
	CommodityID         string `json:"commodity_id"`
	QuantityHeld        int64  `json:"quantity_held"`
	OldUnitPence        int64  `json:"old_unit_pence"`
	NewUnitPence        int64  `json:"new_unit_pence"`
	DeltaPence          int64  `json:"delta_pence"`
}

type speculativeHoldResponse struct {
	ID                   string `json:"id"`
	IndustrialCapitalID  string `json:"industrial_capital_id"`
	CommodityID          string `json:"commodity_id"`
	HeldSince            string `json:"held_since"`
	AnticipatedDirection int8   `json:"anticipated_direction"`
}

type priceCaseRecordResponse struct {
	ValueRevolutionEventID string  `json:"value_revolution_event_id"`
	FallCase               *string `json:"fall_case,omitempty"`
	RiseCase               *string `json:"rise_case,omitempty"`
}

type compoundCapitalAdjustmentResponse struct {
	ID                    string `json:"id"`
	IndustrialCapitalID   string `json:"industrial_capital_id"`
	Period                string `json:"period"`
	RevaluationDeltaPence int64  `json:"revaluation_delta_pence"`
	RealisationDeltaPence int64  `json:"realisation_delta_pence"`
	InventoryDeltaPence   int64  `json:"inventory_delta_pence"`
	NetEffectPence        int64  `json:"net_effect_pence"`
}

type priceRevolutionRecordResponse struct {
	Event           valueRevolutionEventV2Response `json:"event"`
	RevaluedCapital *revaluedCapitalResponse       `json:"revalued_capital,omitempty"`
	RealisedDelta   *realisedValueDeltaResponse    `json:"realised_delta,omitempty"`
	Revaluations    []inventoryRevaluationResponse `json:"inventory_revaluations"`
	FallCase        *string                        `json:"fall_case,omitempty"`
	RiseCase        *string                        `json:"rise_case,omitempty"`
}

// --- mapping helpers ---

func toPriceRevolutionRecordResponse(rec circulation.PriceRevolutionRecord) priceRevolutionRecordResponse {
	evt := rec.Event
	r := priceRevolutionRecordResponse{
		Event: valueRevolutionEventV2Response{
			ID:                   string(evt.ID),
			IndustrialCapitalID:  string(evt.IndustrialCapitalID),
			Affecting:            string(evt.Affecting),
			DirectionPence:       int64(evt.DirectionPence),
			OccurredAt:           evt.OccurredAt.UTC().Format(time.RFC3339),
			BasisPointsChange:    evt.BasisPointsChange,
			RevaluedCapitalPence: int64(evt.RevaluedCapitalPence),
		},
		Revaluations: make([]inventoryRevaluationResponse, 0, len(rec.Revaluations)),
	}
	if rec.RevaluedCapital != nil {
		rc := rec.RevaluedCapital
		resp := revaluedCapitalResponse{
			ID:                     string(rc.ID),
			ValueRevolutionEventID: string(rc.ValueRevolutionEventID),
			IndustrialCapitalID:    string(rc.IndustrialCapitalID),
			OldAdvancePence:        int64(rc.OldAdvancePence),
			NewAdvancePence:        int64(rc.NewAdvancePence),
			DeltaPence:             int64(rc.DeltaPence),
		}
		r.RevaluedCapital = &resp
	}
	if rec.RealisedDelta != nil {
		rvd := rec.RealisedDelta
		resp := realisedValueDeltaResponse{
			ID:                       string(rvd.ID),
			ValueRevolutionEventID:   string(rvd.ValueRevolutionEventID),
			IndustrialCapitalID:      string(rvd.IndustrialCapitalID),
			ExpectedRealisationPence: int64(rvd.ExpectedRealisationPence),
			ActualRealisationPence:   int64(rvd.ActualRealisationPence),
			DeltaPence:               int64(rvd.DeltaPence),
		}
		r.RealisedDelta = &resp
	}
	for _, inv := range rec.Revaluations {
		r.Revaluations = append(r.Revaluations, inventoryRevaluationResponse{
			ID:                  string(inv.ID),
			IndustrialCapitalID: string(inv.IndustrialCapitalID),
			CommodityID:         inv.CommodityID,
			QuantityHeld:        inv.QuantityHeld,
			OldUnitPence:        int64(inv.OldUnitPence),
			NewUnitPence:        int64(inv.NewUnitPence),
			DeltaPence:          int64(inv.DeltaPence),
		})
	}
	if rec.FallCase != nil {
		s := string(*rec.FallCase)
		r.FallCase = &s
	}
	if rec.RiseCase != nil {
		s := string(*rec.RiseCase)
		r.RiseCase = &s
	}
	return r
}

// --- handlers ---

// CreatePriceRevolution handles POST /v1/price-revolutions.
func (h *Handler) CreatePriceRevolution(w http.ResponseWriter, r *http.Request) {
	var req createPriceRevolutionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.IndustrialCapitalID == "" {
		writeError(w, http.StatusBadRequest, "industrial_capital_id is required")
		return
	}
	if req.Affecting == "" {
		writeError(w, http.StatusBadRequest, "affecting is required")
		return
	}

	evt := circulation.ValueRevolutionEvent{
		IndustrialCapitalID: circulation.IndustrialCapitalID(req.IndustrialCapitalID),
		Affecting:           circulation.ValueRevolutionAffecting(req.Affecting),
		DirectionPence:      circulation.Pence(req.DirectionPence),
		BasisPointsChange:   req.BasisPointsChange,
	}
	if req.OccurredAt != "" {
		t, err := time.Parse(time.RFC3339, req.OccurredAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid occurred_at: "+err.Error())
			return
		}
		evt.OccurredAt = t
	}

	rec := circulation.PriceRevolutionRecord{Event: evt}

	// Derive RevaluedCapital for MP/LP events when old_advance_pence is provided.
	affecting := circulation.ValueRevolutionAffecting(req.Affecting)
	if req.OldAdvancePence != nil && *req.OldAdvancePence > 0 &&
		(affecting == circulation.AffectingMeansOfProduction || affecting == circulation.AffectingLabourPower) {
		rc := circulation.DeriveRevaluedCapital(evt, circulation.Pence(*req.OldAdvancePence))
		rec.RevaluedCapital = &rc
		rec.Event.RevaluedCapitalPence = rc.NewAdvancePence
	}

	// Derive RealisedValueDelta for output events.
	if affecting == circulation.AffectingOutput &&
		req.ExpectedRealisationPence != nil && req.ActualRealisationPence != nil {
		rvd := circulation.RealisedValueDelta{
			IndustrialCapitalID:      evt.IndustrialCapitalID,
			ExpectedRealisationPence: circulation.Pence(*req.ExpectedRealisationPence),
			ActualRealisationPence:   circulation.Pence(*req.ActualRealisationPence),
			DeltaPence:               circulation.Pence(*req.ActualRealisationPence - *req.ExpectedRealisationPence),
		}
		if err := rvd.Validate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		rec.RealisedDelta = &rvd
	}

	created, err := h.PriceRevolutions.CreatePriceRevolution(r.Context(), rec)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "price revolution already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/price-revolutions/"+string(created.Event.ID))
	writeJSON(w, http.StatusCreated, toPriceRevolutionRecordResponse(created))
}

// GetPriceRevolution handles GET /v1/price-revolutions/{id}.
func (h *Handler) GetPriceRevolution(w http.ResponseWriter, r *http.Request) {
	id := circulation.ValueRevolutionEventID(r.PathValue("id"))
	rec, err := h.PriceRevolutions.GetPriceRevolution(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "price revolution not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toPriceRevolutionRecordResponse(rec))
}

// RecordPriceRevolutionCase handles POST /v1/price-revolutions/{id}/price-case.
func (h *Handler) RecordPriceRevolutionCase(w http.ResponseWriter, r *http.Request) {
	id := circulation.ValueRevolutionEventID(r.PathValue("id"))
	var req recordPriceCaseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	pc := circulation.PriceCaseRecord{ValueRevolutionEventID: id}
	if req.FallCase != nil {
		fc := circulation.PriceFallCase(*req.FallCase)
		pc.FallCase = &fc
	}
	if req.RiseCase != nil {
		rc := circulation.PriceRiseCase(*req.RiseCase)
		pc.RiseCase = &rc
	}
	recorded, err := h.PriceRevolutions.RecordPriceCase(r.Context(), pc)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "price revolution event not found")
			return
		}
		if isPriceRevolutionValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	resp := priceCaseRecordResponse{
		ValueRevolutionEventID: string(recorded.ValueRevolutionEventID),
	}
	if recorded.FallCase != nil {
		s := string(*recorded.FallCase)
		resp.FallCase = &s
	}
	if recorded.RiseCase != nil {
		s := string(*recorded.RiseCase)
		resp.RiseCase = &s
	}
	writeJSON(w, http.StatusOK, resp)
}

// RecordInventoryRevaluation handles POST /v1/inventory-revaluations.
func (h *Handler) RecordInventoryRevaluation(w http.ResponseWriter, r *http.Request) {
	var req recordInventoryRevaluationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.IndustrialCapitalID == "" {
		writeError(w, http.StatusBadRequest, "industrial_capital_id is required")
		return
	}
	if req.CommodityID == "" {
		writeError(w, http.StatusBadRequest, "commodity_id is required")
		return
	}
	inv := circulation.InventoryRevaluation{
		IndustrialCapitalID: circulation.IndustrialCapitalID(req.IndustrialCapitalID),
		CommodityID:         req.CommodityID,
		QuantityHeld:        req.QuantityHeld,
		OldUnitPence:        circulation.Pence(req.OldUnitPence),
		NewUnitPence:        circulation.Pence(req.NewUnitPence),
		DeltaPence:          circulation.Pence(req.QuantityHeld) * (circulation.Pence(req.NewUnitPence) - circulation.Pence(req.OldUnitPence)),
	}
	if err := inv.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.PriceRevolutions.RecordInventoryRevaluation(r.Context(), inv)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, inventoryRevaluationResponse{
		ID:                  string(created.ID),
		IndustrialCapitalID: string(created.IndustrialCapitalID),
		CommodityID:         created.CommodityID,
		QuantityHeld:        created.QuantityHeld,
		OldUnitPence:        int64(created.OldUnitPence),
		NewUnitPence:        int64(created.NewUnitPence),
		DeltaPence:          int64(created.DeltaPence),
	})
}

// RecordSpeculativeHold handles POST /v1/speculative-holds.
func (h *Handler) RecordSpeculativeHold(w http.ResponseWriter, r *http.Request) {
	var req recordSpeculativeHoldRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.IndustrialCapitalID == "" {
		writeError(w, http.StatusBadRequest, "industrial_capital_id is required")
		return
	}
	if req.CommodityID == "" {
		writeError(w, http.StatusBadRequest, "commodity_id is required")
		return
	}
	sh := circulation.SpeculativeHold{
		IndustrialCapitalID:  circulation.IndustrialCapitalID(req.IndustrialCapitalID),
		CommodityID:          req.CommodityID,
		AnticipatedDirection: req.AnticipatedDirection,
	}
	if req.HeldSince != "" {
		t, err := time.Parse(time.RFC3339, req.HeldSince)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid held_since: "+err.Error())
			return
		}
		sh.HeldSince = t
	}
	if err := sh.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.PriceRevolutions.RecordSpeculativeHold(r.Context(), sh)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, speculativeHoldResponse{
		ID:                   string(created.ID),
		IndustrialCapitalID:  string(created.IndustrialCapitalID),
		CommodityID:          created.CommodityID,
		HeldSince:            created.HeldSince.UTC().Format(time.RFC3339),
		AnticipatedDirection: created.AnticipatedDirection,
	})
}

// GetCompoundCapitalAdjustment handles GET /v1/price-revolutions/compound.
// Query params: industrial_capital_id, period.
func (h *Handler) GetCompoundCapitalAdjustment(w http.ResponseWriter, r *http.Request) {
	icID := r.URL.Query().Get("industrial_capital_id")
	period := r.URL.Query().Get("period")
	if icID == "" {
		writeError(w, http.StatusBadRequest, "industrial_capital_id query param is required")
		return
	}
	if period == "" {
		writeError(w, http.StatusBadRequest, "period query param is required")
		return
	}
	adj, err := h.PriceRevolutions.AggregateCompound(r.Context(),
		circulation.IndustrialCapitalID(icID), period)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, compoundCapitalAdjustmentResponse{
		ID:                    string(adj.ID),
		IndustrialCapitalID:   string(adj.IndustrialCapitalID),
		Period:                adj.Period,
		RevaluationDeltaPence: int64(adj.RevaluationDeltaPence),
		RealisationDeltaPence: int64(adj.RealisationDeltaPence),
		InventoryDeltaPence:   int64(adj.InventoryDeltaPence),
		NetEffectPence:        int64(adj.NetEffectPence),
	})
}

// isPriceRevolutionValidationError reports whether err originates from price
// revolution domain validation (circulation: prefix).
func isPriceRevolutionValidationError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return len(msg) >= 12 && msg[:12] == "circulation:"
}
