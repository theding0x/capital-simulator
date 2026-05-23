package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
)

// --- request/response types -------------------------------------------------

type createTurnoverTimeRequest struct {
	IndustrialCapitalID string `json:"industrial_capital_id"`
}

type addLabourTimeRequest struct {
	Nanos int64 `json:"nanos"`
}

type addLabourInterruptionRequest struct {
	Nanos int64 `json:"nanos"`
}

type recordLatentMPRequest struct {
	IndustrialCapitalID string `json:"industrial_capital_id"`
	Pence               int64  `json:"pence"`
	HeldAt              string `json:"held_at"`
}

type recordNaturalProcessRequest struct {
	IndustrialCapitalID string `json:"industrial_capital_id"`
	Process             string `json:"process"`
	DurationNanos       int64  `json:"duration_nanos"`
	StartedAt           string `json:"started_at"`
}

// sellingPhaseRequest covers both open and close.
// If SellingPhaseID is non-empty, it is a close request.
type sellingPhaseRequest struct {
	// Open fields.
	IndustrialCapitalID string `json:"industrial_capital_id"`
	CommodityCircuitID  string `json:"commodity_circuit_id"`
	OpenedAt            string `json:"opened_at"`
	// Close fields.
	SellingPhaseID string `json:"selling_phase_id"`
	Outcome        string `json:"outcome"`
}

// buyingPhaseRequest covers both open and close.
// If BuyingPhaseID is non-empty, it is a close request.
type buyingPhaseRequest struct {
	// Open fields.
	IndustrialCapitalID string `json:"industrial_capital_id"`
	MoneyCircuitID      string `json:"money_circuit_id"`
	OpenedAt            string `json:"opened_at"`
	MarketLocation      string `json:"market_location"`
	// Close field.
	BuyingPhaseID string `json:"buying_phase_id"`
}

type createPerishabilityRequest struct {
	CommodityID string `json:"commodity_id"`
	WindowNanos int64  `json:"window_nanos"`
}

type createMarketSeparationRequest struct {
	IndustrialCapitalID string `json:"industrial_capital_id"`
	SellingMarketID     string `json:"selling_market_id"`
	BuyingMarketID      string `json:"buying_market_id"`
}

type activeFractionResponse struct {
	TurnoverTimeID         string `json:"turnover_time_id"`
	ProductionTimeNanos    int64  `json:"production_time_nanos"`
	TotalNanos             int64  `json:"total_nanos"`
	ActiveFractionBasisPts int64  `json:"active_fraction_basis_points"`
}

// --- handlers ---------------------------------------------------------------

// CreateTurnoverTime handles POST /v1/turnover-time.
func (h *Handler) CreateTurnoverTime(w http.ResponseWriter, r *http.Request) {
	var req createTurnoverTimeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	tt := circulation.TurnoverTime{
		IndustrialCapitalID: circulation.IndustrialCapitalID(req.IndustrialCapitalID),
	}
	if err := tt.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.CreateTurnoverTime(r.Context(), tt)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// GetTurnoverTime handles GET /v1/turnover-time/{id}.
func (h *Handler) GetTurnoverTime(w http.ResponseWriter, r *http.Request) {
	id := circulation.TurnoverTimeID(r.PathValue("id"))
	tt, err := h.Store.GetTurnoverTime(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tt)
}

// ListTurnoverTimes handles GET /v1/turnover-time.
func (h *Handler) ListTurnoverTimes(w http.ResponseWriter, r *http.Request) {
	out, err := h.Store.ListTurnoverTimes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if out == nil {
		out = []circulation.TurnoverTime{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": out})
}

// AddLabourTime handles POST /v1/turnover-time/{id}/labour-time.
func (h *Handler) AddLabourTime(w http.ResponseWriter, r *http.Request) {
	id := circulation.TurnoverTimeID(r.PathValue("id"))
	var req addLabourTimeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	tt, err := h.Store.AddLabourTime(r.Context(), id, req.Nanos)
	if err != nil {
		writeCirculationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tt)
}

// AddLabourInterruption handles POST /v1/turnover-time/{id}/labour-interruption.
func (h *Handler) AddLabourInterruption(w http.ResponseWriter, r *http.Request) {
	id := circulation.TurnoverTimeID(r.PathValue("id"))
	var req addLabourInterruptionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	tt, err := h.Store.AddLabourInterruption(r.Context(), id, req.Nanos)
	if err != nil {
		writeCirculationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tt)
}

// RecordLatentMP handles POST /v1/turnover-time/{id}/latent-mp.
func (h *Handler) RecordLatentMP(w http.ResponseWriter, r *http.Request) {
	id := circulation.TurnoverTimeID(r.PathValue("id"))
	var req recordLatentMPRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	heldAt, err := parseRFC3339(req.HeldAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid held_at: "+err.Error())
		return
	}
	lpc := circulation.LatentProductiveCapital{
		TurnoverTimeID:      id,
		IndustrialCapitalID: circulation.IndustrialCapitalID(req.IndustrialCapitalID),
		Pence:               circulation.Pence(req.Pence),
		HeldAt:              heldAt,
	}
	created, err := h.Store.RecordLatentMP(r.Context(), lpc)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// RecordNaturalProcess handles POST /v1/turnover-time/{id}/natural-process.
func (h *Handler) RecordNaturalProcess(w http.ResponseWriter, r *http.Request) {
	id := circulation.TurnoverTimeID(r.PathValue("id"))
	var req recordNaturalProcessRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	startedAt, err := parseRFC3339(req.StartedAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid started_at: "+err.Error())
		return
	}
	nps := circulation.NaturalProcessSpan{
		TurnoverTimeID:      id,
		IndustrialCapitalID: circulation.IndustrialCapitalID(req.IndustrialCapitalID),
		Process:             circulation.NaturalProcessKind(req.Process),
		DurationNanos:       req.DurationNanos,
		StartedAt:           startedAt,
	}
	created, err := h.Store.RecordNaturalProcess(r.Context(), nps)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// RecordSellingPhase handles POST /v1/turnover-time/{id}/selling-phase.
// Discriminated by presence of selling_phase_id: close request if set, open otherwise.
func (h *Handler) RecordSellingPhase(w http.ResponseWriter, r *http.Request) {
	id := circulation.TurnoverTimeID(r.PathValue("id"))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read body")
		return
	}
	var req sellingPhaseRequest
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	if req.SellingPhaseID != "" {
		// Close.
		sp, err := h.Store.CloseSellingPhase(r.Context(), id,
			circulation.SellingPhaseID(req.SellingPhaseID),
			circulation.SellingOutcome(req.Outcome))
		if err != nil {
			writeCirculationError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, sp)
		return
	}

	// Open.
	openedAt, err := parseRFC3339(req.OpenedAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid opened_at: "+err.Error())
		return
	}
	sp := circulation.SellingPhase{
		TurnoverTimeID:      id,
		IndustrialCapitalID: circulation.IndustrialCapitalID(req.IndustrialCapitalID),
		CommodityCircuitID:  circulation.CommodityCircuitID(req.CommodityCircuitID),
		OpenedAt:            openedAt,
	}
	created, err := h.Store.OpenSellingPhase(r.Context(), sp)
	if err != nil {
		writeCirculationError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// RecordBuyingPhase handles POST /v1/turnover-time/{id}/buying-phase.
func (h *Handler) RecordBuyingPhase(w http.ResponseWriter, r *http.Request) {
	id := circulation.TurnoverTimeID(r.PathValue("id"))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read body")
		return
	}
	var req buyingPhaseRequest
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}

	if req.BuyingPhaseID != "" {
		// Close.
		bp, err := h.Store.CloseBuyingPhase(r.Context(), id,
			circulation.BuyingPhaseID(req.BuyingPhaseID))
		if err != nil {
			writeCirculationError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, bp)
		return
	}

	// Open.
	openedAt, err := parseRFC3339(req.OpenedAt)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid opened_at: "+err.Error())
		return
	}
	bp := circulation.BuyingPhase{
		TurnoverTimeID:      id,
		IndustrialCapitalID: circulation.IndustrialCapitalID(req.IndustrialCapitalID),
		MoneyCircuitID:      circulation.MoneyCircuitID(req.MoneyCircuitID),
		OpenedAt:            openedAt,
		MarketLocation:      req.MarketLocation,
	}
	created, err := h.Store.OpenBuyingPhase(r.Context(), bp)
	if err != nil {
		writeCirculationError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// SetPerishability handles POST /v1/perishability.
func (h *Handler) SetPerishability(w http.ResponseWriter, r *http.Request) {
	var req createPerishabilityRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	p := circulation.Perishability{
		CommodityID: circulation.CommodityID(req.CommodityID),
		WindowNanos: req.WindowNanos,
	}
	created, err := h.Store.SetPerishability(r.Context(), p)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// SetMarketSeparation handles POST /v1/market-separation.
func (h *Handler) SetMarketSeparation(w http.ResponseWriter, r *http.Request) {
	var req createMarketSeparationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	ms := circulation.MarketSeparation{
		IndustrialCapitalID: circulation.IndustrialCapitalID(req.IndustrialCapitalID),
		SellingMarketID:     req.SellingMarketID,
		BuyingMarketID:      req.BuyingMarketID,
	}
	created, err := h.Store.SetMarketSeparation(r.Context(), ms)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// GetActiveFraction handles GET /v1/turnover-time/{id}/active-fraction.
func (h *Handler) GetActiveFraction(w http.ResponseWriter, r *http.Request) {
	id := circulation.TurnoverTimeID(r.PathValue("id"))
	tt, err := h.Store.GetTurnoverTime(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	resp := activeFractionResponse{
		TurnoverTimeID:         string(tt.ID),
		ProductionTimeNanos:    tt.Production.TotalNanos(),
		TotalNanos:             tt.TotalNanos(),
		ActiveFractionBasisPts: tt.ActiveFractionBasisPoints(),
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- helpers ----------------------------------------------------------------

func writeCirculationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, store.ErrAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, circulation.ErrConcurrentProductionAndCirculation),
		errors.Is(err, circulation.ErrSellingPhaseAlreadyOpen),
		errors.Is(err, circulation.ErrBuyingPhaseAlreadyOpen),
		errors.Is(err, circulation.ErrNoOpenSellingPhase),
		errors.Is(err, circulation.ErrNoOpenBuyingPhase):
		writeError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func parseRFC3339(s string) (time.Time, error) {
	if s == "" {
		return time.Now().UTC(), nil
	}
	return time.Parse(time.RFC3339, s)
}
