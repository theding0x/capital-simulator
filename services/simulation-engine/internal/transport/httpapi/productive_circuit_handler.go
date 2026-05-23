package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Vol. II Ch. 2 — The Circuit of Productive Capital.

type createProductiveCircuitRequest struct {
	AgentID                    string `json:"agent_id,omitempty"`
	ConstantPence              int64  `json:"constant_pence"`
	VariablePence              int64  `json:"variable_pence"`
	Mode                       string `json:"mode"`
	MinCapitalisationIncrement int64  `json:"min_capitalisation_increment_pence"`
}

type latentMoneyCapitalResponse struct {
	ProductiveCircuitID string `json:"productive_circuit_id"`
	Accumulated         int64  `json:"accumulated"`
	Threshold           int64  `json:"threshold"`
	IsArmed             bool   `json:"is_armed"`
}

type reserveFundResponse struct {
	ProductiveCircuitID string `json:"productive_circuit_id"`
	Balance             int64  `json:"balance"`
}

type revenueCircuitResponse struct {
	ProductiveCircuitID string `json:"productive_circuit_id"`
	Amount              int64  `json:"amount"`
	SpentAt             string `json:"spent_at"`
}

type capitalisationStepResponse struct {
	ProductiveCircuitID string `json:"productive_circuit_id"`
	AmountInjected      int64  `json:"amount_injected"`
	DeltaConstantPence  int64  `json:"delta_constant_pence"`
	DeltaVariablePence  int64  `json:"delta_variable_pence"`
	OccurredAt          string `json:"occurred_at"`
}

type reserveDrawResponse struct {
	ProductiveCircuitID string `json:"productive_circuit_id"`
	DrawnPence          int64  `json:"drawn_pence"`
	Reason              string `json:"reason"`
	OccurredAt          string `json:"occurred_at"`
}

type productiveCircuitResponse struct {
	ID                         string                       `json:"id"`
	AgentID                    string                       `json:"agent_id,omitempty"`
	ConstantPence              int64                        `json:"constant_pence"`
	VariablePence              int64                        `json:"variable_pence"`
	TotalAdvance               int64                        `json:"total_advance"`
	Mode                       string                       `json:"mode"`
	MinCapitalisationIncrement int64                        `json:"min_capitalisation_increment_pence"`
	LatentMoneyCapital         latentMoneyCapitalResponse   `json:"latent_money_capital"`
	ReserveFund                reserveFundResponse          `json:"reserve_fund"`
	RevenueCircuits            []revenueCircuitResponse     `json:"revenue_circuits"`
	CapitalisationSteps        []capitalisationStepResponse `json:"capitalisation_steps"`
	ReserveDraws               []reserveDrawResponse        `json:"reserve_draws"`
	CreatedAt                  string                       `json:"created_at,omitempty"`
}

func toProductiveCircuitResponse(pc circulation.ProductiveCircuit) productiveCircuitResponse {
	lmc := pc.LatentMoneyCapital
	rf := pc.ReserveFund

	revs := make([]revenueCircuitResponse, len(pc.RevenueCircuits))
	for i, r := range pc.RevenueCircuits {
		revs[i] = revenueCircuitResponse{
			ProductiveCircuitID: string(r.ProductiveCircuitID),
			Amount:              int64(r.Amount),
			SpentAt:             r.SpentAt.UTC().Format(time.RFC3339),
		}
	}

	steps := make([]capitalisationStepResponse, len(pc.CapitalisationSteps))
	for i, s := range pc.CapitalisationSteps {
		steps[i] = capitalisationStepResponse{
			ProductiveCircuitID: string(s.ProductiveCircuitID),
			AmountInjected:      int64(s.AmountInjected),
			DeltaConstantPence:  int64(s.DeltaConstantPence),
			DeltaVariablePence:  int64(s.DeltaVariablePence),
			OccurredAt:          s.OccurredAt.UTC().Format(time.RFC3339),
		}
	}

	draws := make([]reserveDrawResponse, len(pc.ReserveDraws))
	for i, d := range pc.ReserveDraws {
		draws[i] = reserveDrawResponse{
			ProductiveCircuitID: string(d.ProductiveCircuitID),
			DrawnPence:          int64(d.DrawnPence),
			Reason:              string(d.Reason),
			OccurredAt:          d.OccurredAt.UTC().Format(time.RFC3339),
		}
	}

	return productiveCircuitResponse{
		ID:                         string(pc.ID),
		AgentID:                    pc.AgentID,
		ConstantPence:              int64(pc.ConstantPence),
		VariablePence:              int64(pc.VariablePence),
		TotalAdvance:               int64(pc.TotalAdvance()),
		Mode:                       string(pc.Mode),
		MinCapitalisationIncrement: int64(pc.MinCapitalisationIncrement),
		LatentMoneyCapital: latentMoneyCapitalResponse{
			ProductiveCircuitID: string(lmc.ProductiveCircuitID),
			Accumulated:         int64(lmc.Accumulated),
			Threshold:           int64(lmc.Threshold),
			IsArmed:             lmc.IsArmed(),
		},
		ReserveFund: reserveFundResponse{
			ProductiveCircuitID: string(rf.ProductiveCircuitID),
			Balance:             int64(rf.Balance),
		},
		RevenueCircuits:     revs,
		CapitalisationSteps: steps,
		ReserveDraws:        draws,
		CreatedAt:           pc.CreatedAt.UTC().Format(time.RFC3339),
	}
}

// CreateProductiveCircuit handles POST /v1/productive-circuits.
func (h *Handler) CreateProductiveCircuit(w http.ResponseWriter, r *http.Request) {
	var req createProductiveCircuitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ConstantPence <= 0 {
		writeError(w, http.StatusBadRequest, "constant_pence must be positive")
		return
	}
	if req.VariablePence <= 0 {
		writeError(w, http.StatusBadRequest, "variable_pence must be positive")
		return
	}
	if req.MinCapitalisationIncrement <= 0 {
		writeError(w, http.StatusBadRequest, "min_capitalisation_increment_pence must be positive")
		return
	}

	pc := circulation.ProductiveCircuit{
		AgentID:                    req.AgentID,
		ConstantPence:              circulation.Pence(req.ConstantPence),
		VariablePence:              circulation.Pence(req.VariablePence),
		Mode:                       circulation.ReproductionMode(req.Mode),
		MinCapitalisationIncrement: circulation.Pence(req.MinCapitalisationIncrement),
	}
	created, err := h.ProductiveCircuits.CreateProductiveCircuit(r.Context(), pc)
	if err != nil {
		if isValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/productive-circuits/"+string(created.ID))
	writeJSON(w, http.StatusCreated, toProductiveCircuitResponse(created))
}

// GetProductiveCircuit handles GET /v1/productive-circuits/{id}.
func (h *Handler) GetProductiveCircuit(w http.ResponseWriter, r *http.Request) {
	id := circulation.ProductiveCircuitID(r.PathValue("id"))
	pc, err := h.ProductiveCircuits.GetProductiveCircuit(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "productive circuit not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProductiveCircuitResponse(pc))
}

// ListProductiveCircuits handles GET /v1/productive-circuits.
func (h *Handler) ListProductiveCircuits(w http.ResponseWriter, r *http.Request) {
	list, err := h.ProductiveCircuits.ListProductiveCircuits(r.Context())
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	out := make([]productiveCircuitResponse, len(list))
	for i, pc := range list {
		out[i] = toProductiveCircuitResponse(pc)
	}
	writeJSON(w, http.StatusOK, out)
}

// RecordRevenue handles POST /v1/productive-circuits/{id}/revenue.
func (h *Handler) RecordRevenue(w http.ResponseWriter, r *http.Request) {
	id := circulation.ProductiveCircuitID(r.PathValue("id"))
	var req struct {
		Amount  int64  `json:"amount"`
		SpentAt string `json:"spent_at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	rc := circulation.RevenueCircuit{Amount: circulation.Pence(req.Amount)}
	if req.SpentAt != "" {
		t, err := time.Parse(time.RFC3339, req.SpentAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid spent_at: "+err.Error())
			return
		}
		rc.SpentAt = t
	}
	recorded, err := h.ProductiveCircuits.RecordRevenue(r.Context(), id, rc)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "productive circuit not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, revenueCircuitResponse{
		ProductiveCircuitID: string(recorded.ProductiveCircuitID),
		Amount:              int64(recorded.Amount),
		SpentAt:             recorded.SpentAt.UTC().Format(time.RFC3339),
	})
}

// AccumulateLatentCapital handles POST /v1/productive-circuits/{id}/accumulate.
func (h *Handler) AccumulateLatentCapital(w http.ResponseWriter, r *http.Request) {
	id := circulation.ProductiveCircuitID(r.PathValue("id"))
	var req struct {
		Amount int64 `json:"amount"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	lmc, err := h.ProductiveCircuits.Accumulate(r.Context(), id, circulation.Pence(req.Amount))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "productive circuit not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, latentMoneyCapitalResponse{
		ProductiveCircuitID: string(lmc.ProductiveCircuitID),
		Accumulated:         int64(lmc.Accumulated),
		Threshold:           int64(lmc.Threshold),
		IsArmed:             lmc.IsArmed(),
	})
}

// CapitaliseCircuit handles POST /v1/productive-circuits/{id}/capitalise.
func (h *Handler) CapitaliseCircuit(w http.ResponseWriter, r *http.Request) {
	id := circulation.ProductiveCircuitID(r.PathValue("id"))
	var req struct {
		Amount             int64 `json:"amount"`
		DeltaConstantPence int64 `json:"delta_constant_pence"`
		DeltaVariablePence int64 `json:"delta_variable_pence"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	if req.DeltaConstantPence+req.DeltaVariablePence != req.Amount {
		writeError(w, http.StatusBadRequest, "delta_constant_pence + delta_variable_pence must equal amount")
		return
	}
	step, err := h.ProductiveCircuits.Capitalise(r.Context(), id,
		circulation.Pence(req.Amount),
		circulation.Pence(req.DeltaConstantPence),
		circulation.Pence(req.DeltaVariablePence),
	)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "productive circuit not found")
			return
		}
		if errors.Is(err, circulation.ErrIndivisibleProductionElement) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, capitalisationStepResponse{
		ProductiveCircuitID: string(step.ProductiveCircuitID),
		AmountInjected:      int64(step.AmountInjected),
		DeltaConstantPence:  int64(step.DeltaConstantPence),
		DeltaVariablePence:  int64(step.DeltaVariablePence),
		OccurredAt:          step.OccurredAt.UTC().Format(time.RFC3339),
	})
}

// DepositReserve handles POST /v1/productive-circuits/{id}/reserve/deposit.
func (h *Handler) DepositReserve(w http.ResponseWriter, r *http.Request) {
	id := circulation.ProductiveCircuitID(r.PathValue("id"))
	var req struct {
		Amount int64 `json:"amount"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	fund, err := h.ProductiveCircuits.DepositReserve(r.Context(), id, circulation.Pence(req.Amount))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "productive circuit not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, reserveFundResponse{
		ProductiveCircuitID: string(fund.ProductiveCircuitID),
		Balance:             int64(fund.Balance),
	})
}

// DrawReserve handles POST /v1/productive-circuits/{id}/reserve/draw.
func (h *Handler) DrawReserve(w http.ResponseWriter, r *http.Request) {
	id := circulation.ProductiveCircuitID(r.PathValue("id"))
	var req struct {
		Amount int64  `json:"amount"`
		Reason string `json:"reason"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	reason := circulation.ReserveDrawReason(req.Reason)
	if reason != circulation.ReserveDrawDelayedRealisation && reason != circulation.ReserveDrawPriceRiseInMP {
		writeError(w, http.StatusBadRequest, "reason must be delayed_realisation or price_rise_in_mp")
		return
	}
	draw, err := h.ProductiveCircuits.WithdrawReserve(r.Context(), id, circulation.Pence(req.Amount), reason)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "productive circuit not found")
			return
		}
		if errors.Is(err, circulation.ErrInsufficientReserve) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, reserveDrawResponse{
		ProductiveCircuitID: string(draw.ProductiveCircuitID),
		DrawnPence:          int64(draw.DrawnPence),
		Reason:              string(draw.Reason),
		OccurredAt:          draw.OccurredAt.UTC().Format(time.RFC3339),
	})
}

// isValidationError reports whether err originates from domain validation.
func isValidationError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return len(msg) > 0 && (msg[:len("circulation:")] == "circulation:" || msg == store.ErrAlreadyExists.Error())
}
