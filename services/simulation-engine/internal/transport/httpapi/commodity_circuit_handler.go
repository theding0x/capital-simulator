package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Vol. II Ch. 3 — The Circuit of Commodity-Capital.

type openingCommodityCapitalResponse struct {
	ConstantPence      int64 `json:"constant_pence"`
	VariablePence      int64 `json:"variable_pence"`
	SurplusPence       int64 `json:"surplus_pence"`
	ValueOriginalPence int64 `json:"value_original_pence"`
	Total              int64 `json:"total"`
	PoundsTotal        int64 `json:"pounds_total"`
}

type commodityAugmentedResponse struct {
	CommodityCircuitID string `json:"commodity_circuit_id"`
	ConstantPence      int64  `json:"constant_pence"`
	VariablePence      int64  `json:"variable_pence"`
	SurplusPence       int64  `json:"surplus_pence"`
	CapitalisedPence   int64  `json:"capitalised_pence"`
	Total              int64  `json:"total"`
	PoundsTotal        int64  `json:"pounds_total"`
	IsExtended         bool   `json:"is_extended"`
	ClosedAt           string `json:"closed_at"`
}

type valueCompositionResponse struct {
	ConstantShare int64 `json:"constant_share"`
	VariableShare int64 `json:"variable_share"`
	SurplusShare  int64 `json:"surplus_share"`
	Total         int64 `json:"total"`
}

type successivePartialSaleResponse struct {
	CommodityCircuitID string                   `json:"commodity_circuit_id"`
	Quantity           int64                    `json:"quantity"`
	RealisedPence      int64                    `json:"realised_pence"`
	Decomposition      valueCompositionResponse `json:"decomposition"`
	SoldAt             string                   `json:"sold_at"`
}

type meansOfProductionSourceResponse struct {
	CommodityCircuitID       string `json:"commodity_circuit_id"`
	SourceCommodityCircuitID string `json:"source_commodity_circuit_id,omitempty"`
	SourceKind               string `json:"source_kind"`
}

type commodityCircuitResponse struct {
	ID                string                            `json:"id"`
	AgentID           string                            `json:"agent_id,omitempty"`
	Initial           openingCommodityCapitalResponse   `json:"initial"`
	Mode              string                            `json:"mode"`
	IsFirstInvestment bool                              `json:"is_first_investment"`
	SocialCapitalLens bool                              `json:"social_capital_lens"`
	PartialSales      []successivePartialSaleResponse   `json:"partial_sales"`
	MPSources         []meansOfProductionSourceResponse `json:"mp_sources"`
	Terminal          *commodityAugmentedResponse       `json:"terminal,omitempty"`
	CreatedAt         string                            `json:"created_at"`
}

func toCommodityCircuitResponse(cc circulation.CommodityCircuit) commodityCircuitResponse {
	sales := make([]successivePartialSaleResponse, len(cc.PartialSales))
	for i, s := range cc.PartialSales {
		sales[i] = successivePartialSaleResponse{
			CommodityCircuitID: string(s.CommodityCircuitID),
			Quantity:           s.Quantity,
			RealisedPence:      int64(s.RealisedPence),
			Decomposition: valueCompositionResponse{
				ConstantShare: int64(s.Decomposition.ConstantShare),
				VariableShare: int64(s.Decomposition.VariableShare),
				SurplusShare:  int64(s.Decomposition.SurplusShare),
				Total:         int64(s.Decomposition.Total()),
			},
			SoldAt: s.SoldAt.UTC().Format(time.RFC3339),
		}
	}

	sources := make([]meansOfProductionSourceResponse, len(cc.MPSources))
	for i, src := range cc.MPSources {
		sources[i] = meansOfProductionSourceResponse{
			CommodityCircuitID:       string(src.CommodityCircuitID),
			SourceCommodityCircuitID: src.SourceCommodityCircuitID,
			SourceKind:               string(src.SourceKind),
		}
	}

	resp := commodityCircuitResponse{
		ID:      string(cc.ID),
		AgentID: cc.AgentID,
		Initial: openingCommodityCapitalResponse{
			ConstantPence:      int64(cc.Initial.ConstantPence),
			VariablePence:      int64(cc.Initial.VariablePence),
			SurplusPence:       int64(cc.Initial.SurplusPence),
			ValueOriginalPence: int64(cc.Initial.ValueOriginalPence()),
			Total:              int64(cc.Initial.Total()),
			PoundsTotal:        cc.Initial.PoundsTotal,
		},
		Mode:              string(cc.Mode),
		IsFirstInvestment: cc.IsFirstInvestment,
		SocialCapitalLens: bool(cc.SocialCapitalLens),
		PartialSales:      sales,
		MPSources:         sources,
		CreatedAt:         cc.CreatedAt.UTC().Format(time.RFC3339),
	}

	if cc.Terminal != nil {
		aug := cc.Terminal
		resp.Terminal = &commodityAugmentedResponse{
			CommodityCircuitID: string(aug.CommodityCircuitID),
			ConstantPence:      int64(aug.ConstantPence),
			VariablePence:      int64(aug.VariablePence),
			SurplusPence:       int64(aug.SurplusPence),
			CapitalisedPence:   int64(aug.CapitalisedPence),
			Total:              int64(aug.Total()),
			PoundsTotal:        aug.PoundsTotal,
			IsExtended:         aug.IsExtended(cc.Initial.Total()),
			ClosedAt:           aug.ClosedAt.UTC().Format(time.RFC3339),
		}
	}

	return resp
}

// CreateCommodityCircuit handles POST /v1/commodity-circuits.
func (h *Handler) CreateCommodityCircuit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AgentID           string `json:"agent_id,omitempty"`
		ConstantPence     int64  `json:"constant_pence"`
		VariablePence     int64  `json:"variable_pence"`
		SurplusPence      int64  `json:"surplus_pence"`
		PoundsTotal       int64  `json:"pounds_total"`
		Mode              string `json:"mode"`
		IsFirstInvestment bool   `json:"is_first_investment"`
		SocialCapitalLens bool   `json:"social_capital_lens"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ConstantPence < 0 {
		writeError(w, http.StatusBadRequest, "constant_pence cannot be negative")
		return
	}
	if req.VariablePence < 0 {
		writeError(w, http.StatusBadRequest, "variable_pence cannot be negative")
		return
	}

	cc := circulation.CommodityCircuit{
		AgentID: req.AgentID,
		Initial: circulation.OpeningCommodityCapital{
			ConstantPence: circulation.Pence(req.ConstantPence),
			VariablePence: circulation.Pence(req.VariablePence),
			SurplusPence:  circulation.Pence(req.SurplusPence),
			PoundsTotal:   req.PoundsTotal,
		},
		Mode:              circulation.ReproductionMode(req.Mode),
		IsFirstInvestment: req.IsFirstInvestment,
		SocialCapitalLens: circulation.SocialCapitalLens(req.SocialCapitalLens),
	}
	created, err := h.CommodityCircuits.CreateCommodityCircuit(r.Context(), cc)
	if err != nil {
		if isCommodityCircuitValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/commodity-circuits/"+string(created.ID))
	writeJSON(w, http.StatusCreated, toCommodityCircuitResponse(created))
}

// GetCommodityCircuit handles GET /v1/commodity-circuits/{id}.
func (h *Handler) GetCommodityCircuit(w http.ResponseWriter, r *http.Request) {
	id := circulation.CommodityCircuitID(r.PathValue("id"))
	cc, err := h.CommodityCircuits.GetCommodityCircuit(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "commodity circuit not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toCommodityCircuitResponse(cc))
}

// ListCommodityCircuits handles GET /v1/commodity-circuits.
func (h *Handler) ListCommodityCircuits(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent_id")
	list, err := h.CommodityCircuits.ListCommodityCircuits(r.Context(), agentID)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	out := make([]commodityCircuitResponse, len(list))
	for i, cc := range list {
		out[i] = toCommodityCircuitResponse(cc)
	}
	writeJSON(w, http.StatusOK, out)
}

// RecordCommodityPartialSale handles POST /v1/commodity-circuits/{id}/partial-sales.
// Decomposition is computed server-side from the circuit's per-pound ratios.
func (h *Handler) RecordCommodityPartialSale(w http.ResponseWriter, r *http.Request) {
	id := circulation.CommodityCircuitID(r.PathValue("id"))
	var req struct {
		Quantity      int64  `json:"quantity"`
		RealisedPence int64  `json:"realised_pence"`
		SoldAt        string `json:"sold_at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "quantity must be positive")
		return
	}
	if req.RealisedPence <= 0 {
		writeError(w, http.StatusBadRequest, "realised_pence must be positive")
		return
	}

	cc, err := h.CommodityCircuits.GetCommodityCircuit(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "commodity circuit not found")
			return
		}
		h.writeServerError(w, err)
		return
	}

	decomp := cc.ComputeDecomposition(req.Quantity, circulation.Pence(req.RealisedPence))
	sale := circulation.SuccessivePartialSale{
		Quantity:      req.Quantity,
		RealisedPence: circulation.Pence(req.RealisedPence),
		Decomposition: decomp,
	}
	if req.SoldAt != "" {
		t, err := time.Parse(time.RFC3339, req.SoldAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid sold_at: "+err.Error())
			return
		}
		sale.SoldAt = t
	}

	recorded, err := h.CommodityCircuits.RecordPartialSale(r.Context(), id, sale)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "commodity circuit not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, successivePartialSaleResponse{
		CommodityCircuitID: string(recorded.CommodityCircuitID),
		Quantity:           recorded.Quantity,
		RealisedPence:      int64(recorded.RealisedPence),
		Decomposition: valueCompositionResponse{
			ConstantShare: int64(recorded.Decomposition.ConstantShare),
			VariableShare: int64(recorded.Decomposition.VariableShare),
			SurplusShare:  int64(recorded.Decomposition.SurplusShare),
			Total:         int64(recorded.Decomposition.Total()),
		},
		SoldAt: recorded.SoldAt.UTC().Format(time.RFC3339),
	})
}

// LinkCommodityMPSource handles POST /v1/commodity-circuits/{id}/mp-source.
func (h *Handler) LinkCommodityMPSource(w http.ResponseWriter, r *http.Request) {
	id := circulation.CommodityCircuitID(r.PathValue("id"))
	var req struct {
		SourceCommodityCircuitID string `json:"source_commodity_circuit_id,omitempty"`
		SourceKind               string `json:"source_kind"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	source := circulation.MeansOfProductionSource{
		SourceCommodityCircuitID: req.SourceCommodityCircuitID,
		SourceKind:               circulation.MeansOfProductionSourceKind(req.SourceKind),
	}
	linked, err := h.CommodityCircuits.LinkMPSource(r.Context(), id, source)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "commodity circuit not found")
			return
		}
		if errors.Is(err, circulation.ErrUnsourcedMeansOfProduction) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, meansOfProductionSourceResponse{
		CommodityCircuitID:       string(linked.CommodityCircuitID),
		SourceCommodityCircuitID: linked.SourceCommodityCircuitID,
		SourceKind:               string(linked.SourceKind),
	})
}

// CloseCommodityCircuit handles POST /v1/commodity-circuits/{id}/close.
func (h *Handler) CloseCommodityCircuit(w http.ResponseWriter, r *http.Request) {
	id := circulation.CommodityCircuitID(r.PathValue("id"))
	var req struct {
		ConstantPence    int64  `json:"constant_pence"`
		VariablePence    int64  `json:"variable_pence"`
		SurplusPence     int64  `json:"surplus_pence"`
		CapitalisedPence int64  `json:"capitalised_pence"`
		PoundsTotal      int64  `json:"pounds_total"`
		ClosedAt         string `json:"closed_at,omitempty"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ConstantPence < 0 || req.VariablePence < 0 || req.SurplusPence < 0 {
		writeError(w, http.StatusBadRequest, "terminal pence values cannot be negative")
		return
	}

	aug := circulation.CommodityAugmented{
		ConstantPence:    circulation.Pence(req.ConstantPence),
		VariablePence:    circulation.Pence(req.VariablePence),
		SurplusPence:     circulation.Pence(req.SurplusPence),
		CapitalisedPence: circulation.Pence(req.CapitalisedPence),
		PoundsTotal:      req.PoundsTotal,
	}
	if req.ClosedAt != "" {
		t, err := time.Parse(time.RFC3339, req.ClosedAt)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid closed_at: "+err.Error())
			return
		}
		aug.ClosedAt = t
	}

	closed, err := h.CommodityCircuits.CloseCommodityCircuit(r.Context(), id, aug)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "commodity circuit not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toCommodityCircuitResponse(closed))
}

// GetCommodityCircuitAggregate handles GET /v1/commodity-circuits/aggregate.
// Returns totals across all circuits for the social-capital lens.
func (h *Handler) GetCommodityCircuitAggregate(w http.ResponseWriter, r *http.Request) {
	list, err := h.CommodityCircuits.ListCommodityCircuits(r.Context(), "")
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	var (
		circuitCount       int
		closedCount        int
		totalInitialPence  int64
		totalTerminalPence int64
	)
	for _, cc := range list {
		circuitCount++
		totalInitialPence += int64(cc.Initial.Total())
		if cc.Terminal != nil {
			totalTerminalPence += int64(cc.Terminal.Total())
			closedCount++
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"circuit_count":        circuitCount,
		"closed_count":         closedCount,
		"total_initial_pence":  totalInitialPence,
		"total_terminal_pence": totalTerminalPence,
	})
}

func isCommodityCircuitValidationError(err error) bool {
	return errors.Is(err, circulation.ErrNotCommodityCapital) ||
		errors.Is(err, circulation.ErrInsufficientMaterialBase) ||
		errors.Is(err, circulation.ErrUnsourcedMeansOfProduction) ||
		isValidationError(err)
}
