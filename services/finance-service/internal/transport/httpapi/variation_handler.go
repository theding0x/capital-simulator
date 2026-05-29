package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// formulaRequest is one capital decomposition in a Ch. 3 request body: the
// constant capital c, the variable capital v, and the rate of surplus-value s'
// in basis points (100% = 10000 bp). The rate of profit p' = s'(v/C) is derived
// server-side.
type formulaRequest struct {
	C     int64 `json:"c"`
	V     int64 `json:"v"`
	SRate int64 `json:"s_rate"`
}

func (r formulaRequest) valid() bool { return r.C >= 0 && r.V >= 0 && r.SRate >= 0 }

func (r formulaRequest) toFormula() profit.ProfitRateFormula {
	return profit.ProfitRateFormula{
		C:     profit.ConstantCapital(r.C),
		V:     profit.VariableCapital(r.V),
		SRate: profit.RateOfSurplusValue(r.SRate),
	}
}

// variationRequest is the body for POST /v1/profit/variation: an initial and a
// changed decomposition whose rate-of-profit movement is to be analysed.
type variationRequest struct {
	Initial formulaRequest `json:"initial"`
	Changed formulaRequest `json:"changed"`
}

// CreateVariation handles POST /v1/profit/variation — from an initial and a
// changed decomposition, classify the variation case, compute the old and new
// rate of profit and value-composition, flag a proportional change, and persist.
func (h *Handler) CreateVariation(w http.ResponseWriter, r *http.Request) {
	var req variationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !req.Initial.valid() || !req.Changed.valid() {
		writeError(w, http.StatusUnprocessableEntity, "c, v and s_rate must be non-negative")
		return
	}
	analysis := profit.ComputeVariationAnalysis(req.Initial.toFormula(), req.Changed.toFormula())
	saved, err := h.Store.CreateVariation(r.Context(), analysis)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/profit/variation/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// GetVariation handles GET /v1/profit/variation/{id}.
func (h *Handler) GetVariation(w http.ResponseWriter, r *http.Request) {
	id := profit.VariationAnalysisID(r.PathValue("id"))
	a, err := h.Store.GetVariation(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

// ListVariations handles GET /v1/profit/variation.
func (h *Handler) ListVariations(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListVariations(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []profit.VariationAnalysis{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// compareRequest is the body for POST /v1/profit/compare: two capitals set side
// by side to read off how their rates of profit relate.
type compareRequest struct {
	First  formulaRequest `json:"first"`
	Second formulaRequest `json:"second"`
}

// CompareProfitRates handles POST /v1/profit/compare — compute both rates of
// profit and value-compositions and flag whether the capitals share s' (in
// which case their rates of profit stand as their value-compositions). This is
// a compute-only endpoint; nothing is persisted.
func (h *Handler) CompareProfitRates(w http.ResponseWriter, r *http.Request) {
	var req compareRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !req.First.valid() || !req.Second.valid() {
		writeError(w, http.StatusUnprocessableEntity, "c, v and s_rate must be non-negative")
		return
	}
	comparison := profit.CompareProfitRates(req.First.toFormula(), req.Second.toFormula())
	writeJSON(w, http.StatusOK, comparison)
}
