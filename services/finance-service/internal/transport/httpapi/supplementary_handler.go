package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// compositionEffectRequest is the body for POST /v1/profit/composition-effect. Two
// enterprises A and B, each a decomposition s/v/c; with equal s and v the gap in
// their rates of profit measures the effect of the organic composition alone. The
// two rates and the gap are derived server-side.
type compositionEffectRequest struct {
	SA int64 `json:"s_a"`
	VA int64 `json:"v_a"`
	CA int64 `json:"c_a"`
	SB int64 `json:"s_b"`
	VB int64 `json:"v_b"`
	CB int64 `json:"c_b"`
}

// CreateCompositionEffect handles POST /v1/profit/composition-effect — compute the
// two rates of profit and the gap between them and persist.
func (h *Handler) CreateCompositionEffect(w http.ResponseWriter, r *http.Request) {
	var req compositionEffectRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.SA < 0 || req.VA < 0 || req.CA < 0 || req.SB < 0 || req.VB < 0 || req.CB < 0 {
		writeError(w, http.StatusUnprocessableEntity, "all magnitudes must be non-negative")
		return
	}
	analysis := profit.ComputeCompositionEffectAnalysis(profit.OrganicCompositionEffect{
		SA: profit.SurplusValue(req.SA),
		VA: profit.VariableCapital(req.VA),
		CA: profit.ConstantCapital(req.CA),
		SB: profit.SurplusValue(req.SB),
		VB: profit.VariableCapital(req.VB),
		CB: profit.ConstantCapital(req.CB),
	})
	saved, err := h.Store.CreateCompositionEffect(r.Context(), analysis)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/profit/composition-effect/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetCompositionEffect handles GET /v1/profit/composition-effect/{id}.
func (h *Handler) GetCompositionEffect(w http.ResponseWriter, r *http.Request) {
	id := profit.CompositionEffectID(r.PathValue("id"))
	a, err := h.Store.GetCompositionEffect(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// ListCompositionEffects handles GET /v1/profit/composition-effect.
func (h *Handler) ListCompositionEffects(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCompositionEffects(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []profit.CompositionEffectAnalysis{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// magnitudeChangeRequest is the body for POST /v1/profit/magnitude-change. A
// change in the absolute size of a capital (OriginalCapital yielding
// OriginalProfit) by Factor — a percentage of the original (100 = unchanged) — of
// a given Kind. Whether the rate of profit moves is derived server-side.
type magnitudeChangeRequest struct {
	Kind            string `json:"kind"`
	OriginalCapital int64  `json:"original_capital"`
	OriginalProfit  int64  `json:"original_profit"`
	Factor          int64  `json:"factor"`
}

// CreateMagnitudeChange handles POST /v1/profit/magnitude-change — compute the old
// and new rates of profit and whether the rate is unchanged, and persist.
func (h *Handler) CreateMagnitudeChange(w http.ResponseWriter, r *http.Request) {
	var req magnitudeChangeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	kind := profit.MagnitudeChangeKind(req.Kind)
	if !kind.Valid() {
		writeError(w, http.StatusUnprocessableEntity,
			"kind must be one of money_value_change, proportional_change, organic_change")
		return
	}
	if req.OriginalCapital < 0 || req.OriginalProfit < 0 {
		writeError(w, http.StatusUnprocessableEntity, "original_capital and original_profit must be non-negative")
		return
	}
	if req.Factor <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "factor must be positive (a percentage of the original magnitude: 100 = unchanged)")
		return
	}
	analysis := profit.ComputeMagnitudeChangeAnalysis(profit.CapitalMagnitudeChange{
		Kind:            kind,
		OriginalCapital: profit.TotalCapital(req.OriginalCapital),
		OriginalProfit:  profit.SurplusValue(req.OriginalProfit),
		Factor:          req.Factor,
	})
	saved, err := h.Store.CreateMagnitudeChange(r.Context(), analysis)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/profit/magnitude-change/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetMagnitudeChange handles GET /v1/profit/magnitude-change/{id}.
func (h *Handler) GetMagnitudeChange(w http.ResponseWriter, r *http.Request) {
	id := profit.MagnitudeChangeID(r.PathValue("id"))
	a, err := h.Store.GetMagnitudeChange(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// ListMagnitudeChanges handles GET /v1/profit/magnitude-change.
func (h *Handler) ListMagnitudeChanges(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListMagnitudeChanges(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []profit.MagnitudeChangeAnalysis{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetPartISummary handles GET /v1/profit/summary — a read-only snapshot gathering
// the rate-of-profit analyses recorded across Part I and naming the determinants
// the chapters have established. The analyses are drawn from the rate-of-profit
// store (Ch. 2).
func (h *Handler) GetPartISummary(w http.ResponseWriter, r *http.Request) {
	analyses, err := h.Store.ListProfitRates(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, profit.ComputePartISummary(analyses))
}
