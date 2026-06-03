package httpapi

import (
	"fmt"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// trinityStreamRequest is one branch of a trinity-formula creation request.
type trinityStreamRequest struct {
	ApparentRevenueBP int64 `json:"apparent_revenue_bp"`
	ActualSourceBP    int64 `json:"actual_source_bp"`
	IsFetishised      bool  `json:"is_fetishised"`
}

// createTrinityFormulaRequest is the body for POST /v1/revenue/trinity-formulas.
// The three slots fix the source of each stream (capital/land/labour).
type createTrinityFormulaRequest struct {
	CapitalStream trinityStreamRequest `json:"capital_stream"`
	LandStream    trinityStreamRequest `json:"land_stream"`
	LabourStream  trinityStreamRequest `json:"labour_stream"`
}

// trinityFormulaResponse bundles the saved formula with its three streams so the
// caller sees the partition of surplus-value the formula records. FetishOvershootBP
// surfaces the gap by which apparent revenues exceed the value actually created
// (v + s) — the ground-rent counted twice, i.e. the fetish made visible.
type trinityFormulaResponse struct {
	Formula           revenue.TrinityFormula  `json:"formula"`
	Streams           []revenue.RevenueStream `json:"streams"`
	FetishOvershootBP int64                   `json:"fetish_overshoot_bp"`
}

// CreateTrinityFormula handles POST /v1/revenue/trinity-formulas. It creates the
// three revenue streams, computes the formula totals from them, and persists the
// binding formula. Returns 400 on bad JSON, 422 on invalid fields, 201 + Location.
func (h *Handler) CreateTrinityFormula(w http.ResponseWriter, r *http.Request) {
	var req createTrinityFormulaRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.CapitalStream.ApparentRevenueBP < 0 || req.CapitalStream.ActualSourceBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "capital_stream amounts must be >= 0")
		return
	}
	if req.LandStream.ApparentRevenueBP < 0 || req.LandStream.ActualSourceBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "land_stream amounts must be >= 0")
		return
	}
	if req.LabourStream.ApparentRevenueBP < 0 || req.LabourStream.ActualSourceBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "labour_stream amounts must be >= 0")
		return
	}
	// Invariant: wages recover only variable capital — they are not surplus-value.
	if req.LabourStream.ActualSourceBP != 0 {
		writeError(w, http.StatusUnprocessableEntity, "labour_stream.actual_source_bp must be 0 (wages recover variable capital, not surplus-value)")
		return
	}

	// Conservation: profit + rent must partition a single surplus-value, so the
	// apparent-revenue overshoot above newly-created value (v + s) equals exactly
	// the ground-rent — the trinity's one fetish, surfaced as fetish_overshoot_bp
	// in the response rather than asserted away. Validate before persisting so a
	// malformed trinity leaves no orphan streams.
	reqStreams := []revenue.RevenueStream{
		{Source: revenue.RevenueSourceCapital, ApparentRevenueBP: req.CapitalStream.ApparentRevenueBP, ActualSourceBP: req.CapitalStream.ActualSourceBP},
		{Source: revenue.RevenueSourceLand, ApparentRevenueBP: req.LandStream.ApparentRevenueBP, ActualSourceBP: req.LandStream.ActualSourceBP},
		{Source: revenue.RevenueSourceLabour, ApparentRevenueBP: req.LabourStream.ApparentRevenueBP, ActualSourceBP: req.LabourStream.ActualSourceBP},
	}
	if err := revenue.CheckConservation(reqStreams); err != nil {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf(
			"%s: apparent overshoot %d bp != ground-rent %d bp",
			err.Error(), revenue.FetishOvershootBP(reqStreams), revenue.RentBP(reqStreams)))
		return
	}

	inputs := []struct {
		source revenue.RevenueSourceKind
		stream trinityStreamRequest
	}{
		{revenue.RevenueSourceCapital, req.CapitalStream},
		{revenue.RevenueSourceLand, req.LandStream},
		{revenue.RevenueSourceLabour, req.LabourStream},
	}
	streams := make([]revenue.RevenueStream, 0, len(inputs))
	for _, in := range inputs {
		saved, err := h.Store.CreateRevenueStream(r.Context(), revenue.RevenueStream{
			Source:            in.source,
			ApparentRevenueBP: in.stream.ApparentRevenueBP,
			ActualSourceBP:    in.stream.ActualSourceBP,
			IsFetishised:      in.stream.IsFetishised,
		})
		if err != nil {
			writeStoreError(w, err)
			return
		}
		streams = append(streams, saved)
	}

	tf := revenue.TrinityFormula{
		CapitalStreamID:        string(streams[0].ID),
		LandStreamID:           string(streams[1].ID),
		LabourStreamID:         string(streams[2].ID),
		TotalSurplusValueBP:    revenue.SumActualSourceBP(streams),
		TotalApparentRevenueBP: revenue.SumApparentRevenueBP(streams),
	}
	saved, err := h.Store.CreateTrinityFormula(r.Context(), tf)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/revenue/trinity-formulas/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, trinityFormulaResponse{
		Formula:           saved,
		Streams:           streams,
		FetishOvershootBP: revenue.FetishOvershootBP(streams),
	})
}

// GetTrinityFormula handles GET /v1/revenue/trinity-formulas/{id}.
// Returns 404 when no formula has the given ID, 200 + the formula otherwise.
func (h *Handler) GetTrinityFormula(w http.ResponseWriter, r *http.Request) {
	id := revenue.TrinityFormulaID(r.PathValue("id"))
	tf, err := h.Store.GetTrinityFormula(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tf)
}

// ListRevenueFetishForms handles GET /v1/revenue/fetish-forms.
// Returns {"items":[...]} — never null.
func (h *Handler) ListRevenueFetishForms(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListRevenueFetishForms(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []revenue.RevenueFetishForm{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
