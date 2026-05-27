package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
)

// --- request types ----------------------------------------------------------

type createSellingPhaseDetailedRequest struct {
	SellingPhaseID          string `json:"selling_phase_id"`
	IndustrialCapitalID     string `json:"industrial_capital_id"`
	MarketID                string `json:"market_id"`
	MarketDistanceMeters    int64  `json:"market_distance_meters"`
	CommunicationLagMinutes int64  `json:"communication_lag_minutes"`
	BuyerSearchMinutes      int64  `json:"buyer_search_minutes"`
	BargainingMinutes       int64  `json:"bargaining_minutes"`
	LegalSettlementMinutes  int64  `json:"legal_settlement_minutes"`
	TotalMinutes            int64  `json:"total_minutes"`
	// DistanceLagRelationID is optional. When provided, the lag-bound
	// invariant is validated against the referenced relation.
	DistanceLagRelationID string `json:"distance_lag_relation_id"`
}

type createBuyingPhaseDetailedRequest struct {
	BuyingPhaseID           string `json:"buying_phase_id"`
	IndustrialCapitalID     string `json:"industrial_capital_id"`
	MarketID                string `json:"market_id"`
	MarketDistanceMeters    int64  `json:"market_distance_meters"`
	CommunicationLagMinutes int64  `json:"communication_lag_minutes"`
	SupplierSearchMinutes   int64  `json:"supplier_search_minutes"`
	OrderProcessingMinutes  int64  `json:"order_processing_minutes"`
	LegalSettlementMinutes  int64  `json:"legal_settlement_minutes"`
	TotalMinutes            int64  `json:"total_minutes"`
	DistanceLagRelationID   string `json:"distance_lag_relation_id"`
}

type createDistanceLagRelationRequest struct {
	MarketID        string `json:"market_id"`
	DistanceMeters  int64  `json:"distance_meters"`
	LagMinutesPerKm int64  `json:"lag_minutes_per_km"`
}

type createCirculationSpeedImprovementRequest struct {
	MarketID           string `json:"market_id"`
	OldLagMinutesPerKm int64  `json:"old_lag_minutes_per_km"`
	NewLagMinutesPerKm int64  `json:"new_lag_minutes_per_km"`
	Reason             string `json:"reason"`
	EffectiveAt        string `json:"effective_at"` // RFC3339
}

// --- handlers ---------------------------------------------------------------

// UpgradeSellingPhase handles POST /v1/turnover-time/{id}/selling-phase/{phaseId}/detail.
// It is non-destructive: it does not alter the parent SellingPhase's OpenedAt,
// ClosedAt, or Outcome.
func (h *Handler) UpgradeSellingPhase(w http.ResponseWriter, r *http.Request) {
	turnoverID := circulation.TurnoverTimeID(r.PathValue("id"))
	phaseID := circulation.SellingPhaseID(r.PathValue("phaseId"))

	var req createSellingPhaseDetailedRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	d := circulation.SellingPhaseDetailed{
		SellingPhaseID:          phaseID,
		TurnoverTimeID:          turnoverID,
		IndustrialCapitalID:     circulation.IndustrialCapitalID(req.IndustrialCapitalID),
		MarketID:                circulation.MarketID(req.MarketID),
		MarketDistanceMeters:    req.MarketDistanceMeters,
		CommunicationLagMinutes: req.CommunicationLagMinutes,
		BuyerSearchMinutes:      req.BuyerSearchMinutes,
		BargainingMinutes:       req.BargainingMinutes,
		LegalSettlementMinutes:  req.LegalSettlementMinutes,
		TotalMinutes:            req.TotalMinutes,
	}

	if err := d.ValidateSubSpans(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Optional lag-bound check against a DistanceLagRelation.
	if req.DistanceLagRelationID != "" {
		rel, err := h.Store.GetDistanceLagRelation(r.Context(), circulation.DistanceLagRelationID(req.DistanceLagRelationID))
		if err != nil {
			writeStoreError(w, err)
			return
		}
		if err := d.ValidateLagBound(rel.LagMinutesPerKm); err != nil {
			if errors.Is(err, circulation.ErrLagExceedsMaximum) {
				writeError(w, http.StatusUnprocessableEntity, err.Error())
				return
			}
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	created, err := h.Store.CreateSellingPhaseDetailed(r.Context(), d)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// UpgradeBuyingPhase handles POST /v1/turnover-time/{id}/buying-phase/{phaseId}/detail.
func (h *Handler) UpgradeBuyingPhase(w http.ResponseWriter, r *http.Request) {
	turnoverID := circulation.TurnoverTimeID(r.PathValue("id"))
	phaseID := circulation.BuyingPhaseID(r.PathValue("phaseId"))

	var req createBuyingPhaseDetailedRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	d := circulation.BuyingPhaseDetailed{
		BuyingPhaseID:           phaseID,
		TurnoverTimeID:          turnoverID,
		IndustrialCapitalID:     circulation.IndustrialCapitalID(req.IndustrialCapitalID),
		MarketID:                circulation.MarketID(req.MarketID),
		MarketDistanceMeters:    req.MarketDistanceMeters,
		CommunicationLagMinutes: req.CommunicationLagMinutes,
		SupplierSearchMinutes:   req.SupplierSearchMinutes,
		OrderProcessingMinutes:  req.OrderProcessingMinutes,
		LegalSettlementMinutes:  req.LegalSettlementMinutes,
		TotalMinutes:            req.TotalMinutes,
	}

	if err := d.ValidateSubSpans(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.DistanceLagRelationID != "" {
		rel, err := h.Store.GetDistanceLagRelation(r.Context(), circulation.DistanceLagRelationID(req.DistanceLagRelationID))
		if err != nil {
			writeStoreError(w, err)
			return
		}
		if err := d.ValidateLagBound(rel.LagMinutesPerKm); err != nil {
			if errors.Is(err, circulation.ErrLagExceedsMaximum) {
				writeError(w, http.StatusUnprocessableEntity, err.Error())
				return
			}
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	created, err := h.Store.CreateBuyingPhaseDetailed(r.Context(), d)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// CreateDistanceLagRelation handles POST /v1/distance-lag.
func (h *Handler) CreateDistanceLagRelation(w http.ResponseWriter, r *http.Request) {
	var req createDistanceLagRelationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.MarketID == "" {
		writeError(w, http.StatusBadRequest, "market_id is required")
		return
	}
	rel := circulation.DistanceLagRelation{
		MarketID:        circulation.MarketID(req.MarketID),
		DistanceMeters:  req.DistanceMeters,
		LagMinutesPerKm: req.LagMinutesPerKm,
	}
	created, err := h.Store.CreateDistanceLagRelation(r.Context(), rel)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// ListDistanceLagRelations handles GET /v1/distance-lag.
func (h *Handler) ListDistanceLagRelations(w http.ResponseWriter, r *http.Request) {
	out, err := h.Store.ListDistanceLagRelations(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out == nil {
		out = []circulation.DistanceLagRelation{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// CreateCirculationSpeedImprovement handles POST /v1/circulation-speed-improvements.
func (h *Handler) CreateCirculationSpeedImprovement(w http.ResponseWriter, r *http.Request) {
	var req createCirculationSpeedImprovementRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.MarketID == "" {
		writeError(w, http.StatusBadRequest, "market_id is required")
		return
	}
	effectiveAt, err := parseRFC3339(req.EffectiveAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, "effective_at: "+err.Error())
		return
	}
	csi := circulation.CirculationSpeedImprovement{
		MarketID:           circulation.MarketID(req.MarketID),
		OldLagMinutesPerKm: req.OldLagMinutesPerKm,
		NewLagMinutesPerKm: req.NewLagMinutesPerKm,
		Reason:             circulation.CirculationSpeedReason(req.Reason),
		EffectiveAt:        effectiveAt,
	}
	created, err := h.Store.CreateCirculationSpeedImprovement(r.Context(), csi)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// ListCirculationSpeedImprovements handles GET /v1/circulation-speed-improvements.
// Optional query param: ?market_id=<id> to filter by market.
func (h *Handler) ListCirculationSpeedImprovements(w http.ResponseWriter, r *http.Request) {
	marketID := circulation.MarketID(r.URL.Query().Get("market_id"))
	out, err := h.Store.ListCirculationSpeedImprovements(r.Context(), marketID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out == nil {
		out = []circulation.CirculationSpeedImprovement{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// GetAnnualSurplusPenalty handles GET /v1/turnover-time/{id}/annual-surplus-penalty.
// Query param: ?period=<label> (e.g. "1871").
// Returns a stub; full calculation is deferred to Vol. II Ch. 16.
func (h *Handler) GetAnnualSurplusPenalty(w http.ResponseWriter, r *http.Request) {
	turnoverID := r.PathValue("id")
	period := r.URL.Query().Get("period")
	if period == "" {
		writeError(w, http.StatusBadRequest, "period query parameter is required")
		return
	}
	// Look up the TurnoverTime to get the IndustrialCapitalID.
	tt, err := h.Store.GetTurnoverTime(r.Context(), circulation.TurnoverTimeID(turnoverID))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	asp, err := h.Store.GetAnnualSurplusPenalty(r.Context(), tt.IndustrialCapitalID, period)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// Return a zero-penalty stub if no record exists yet.
			asp = circulation.AnnualSurplusPenalty{
				IndustrialCapitalID:         tt.IndustrialCapitalID,
				Period:                      period,
				IdealAnnualRateBasisPoints:  10000,
				ActualAnnualRateBasisPoints: 10000,
				PenaltyBasisPoints:          0,
			}
		} else {
			writeStoreError(w, err)
			return
		}
	}
	writeJSON(w, http.StatusOK, asp)
}
