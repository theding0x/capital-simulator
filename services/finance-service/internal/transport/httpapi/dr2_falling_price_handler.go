package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createSoilExclusionEventRequest is the body for POST /v1/rent/exclusion-events.
type createSoilExclusionEventRequest struct {
	ExcludedGrade          int   `json:"excluded_grade"`
	NewRegulatingGrade     int   `json:"new_regulating_grade"`
	NewPriceOfProductionBP int64 `json:"new_price_of_production_bp"`
}

// CreateSoilExclusionEvent handles POST /v1/rent/exclusion-events.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateSoilExclusionEvent(w http.ResponseWriter, r *http.Request) {
	var req createSoilExclusionEventRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	excluded := rent.SoilGrade(req.ExcludedGrade)
	if !excluded.IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "excluded_grade must be 1–4")
		return
	}
	newReg := rent.SoilGrade(req.NewRegulatingGrade)
	if !newReg.IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "new_regulating_grade must be 1–4")
		return
	}
	if req.NewPriceOfProductionBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "new_price_of_production_bp must be > 0")
		return
	}
	event := rent.SoilExclusionEvent{
		ExcludedGrade:          excluded,
		NewRegulatingGrade:     newReg,
		NewPriceOfProductionBP: req.NewPriceOfProductionBP,
	}
	if err := event.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	saved, err := h.Store.CreateSoilExclusionEvent(r.Context(), event)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/exclusion-events")
	writeJSON(w, http.StatusCreated, saved)
}

// ListSoilExclusionEvents handles GET /v1/rent/exclusion-events.
// Returns {"items":[...]} — never null.
func (h *Handler) ListSoilExclusionEvents(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListSoilExclusionEvents(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.SoilExclusionEvent{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// fallingDR2InvestmentItem is one investment in the POST /v1/rent/dr2-falling-outcomes body.
type fallingDR2InvestmentItem struct {
	ParcelID        string `json:"parcel_id"`
	TrancheNumber   int    `json:"tranche_number"`
	CapitalBP       int64  `json:"capital_bp"`
	OutputQuarters  int64  `json:"output_quarters"`
	SurplusProfitBP int64  `json:"surplus_profit_bp"`
}

// createFallingPriceDR2OutcomeRequest is the body for POST /v1/rent/dr2-falling-outcomes.
type createFallingPriceDR2OutcomeRequest struct {
	Variant          int                        `json:"variant"`
	ExclusionEventID string                     `json:"exclusion_event_id"`
	Investments      []fallingDR2InvestmentItem `json:"investments"`
}

// CreateFallingPriceDR2Outcome handles POST /v1/rent/dr2-falling-outcomes.
// Returns 400 on bad JSON, 422 on invalid fields, 404 if exclusion event missing, 201 + Location on success.
func (h *Handler) CreateFallingPriceDR2Outcome(w http.ResponseWriter, r *http.Request) {
	var req createFallingPriceDR2OutcomeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	variant := rent.FallingPriceVariant(req.Variant)
	if !variant.IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "variant must be 1–3")
		return
	}
	if req.ExclusionEventID == "" {
		writeError(w, http.StatusUnprocessableEntity, "exclusion_event_id must not be empty")
		return
	}
	if len(req.Investments) < 1 {
		writeError(w, http.StatusUnprocessableEntity, "investments must have at least 1 item")
		return
	}
	invs := make([]rent.SuccessiveInvestment, 0, len(req.Investments))
	for i, item := range req.Investments {
		if item.ParcelID == "" {
			writeError(w, http.StatusUnprocessableEntity, "investments["+itoa(i)+"]: parcel_id must not be empty")
			return
		}
		if item.TrancheNumber < 1 {
			writeError(w, http.StatusUnprocessableEntity, "investments["+itoa(i)+"]: tranche_number must be >= 1")
			return
		}
		if item.CapitalBP <= 0 {
			writeError(w, http.StatusUnprocessableEntity, "investments["+itoa(i)+"]: capital_bp must be > 0")
			return
		}
		if item.OutputQuarters <= 0 {
			writeError(w, http.StatusUnprocessableEntity, "investments["+itoa(i)+"]: output_quarters must be > 0")
			return
		}
		if item.SurplusProfitBP < 0 {
			writeError(w, http.StatusUnprocessableEntity, "investments["+itoa(i)+"]: surplus_profit_bp must be >= 0")
			return
		}
		invs = append(invs, rent.SuccessiveInvestment{
			ParcelID:        item.ParcelID,
			TrancheNumber:   item.TrancheNumber,
			CapitalBP:       item.CapitalBP,
			OutputQuarters:  item.OutputQuarters,
			SurplusProfitBP: item.SurplusProfitBP,
		})
	}
	excl, err := h.Store.GetSoilExclusionEvent(r.Context(), rent.SoilExclusionEventID(req.ExclusionEventID))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	out := rent.ComputeFallingPriceDR2(variant, excl, invs)
	saved, err := h.Store.CreateFallingPriceDR2Outcome(r.Context(), out)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/dr2-falling-outcomes")
	writeJSON(w, http.StatusCreated, saved)
}

// ListFallingPriceDR2Outcomes handles GET /v1/rent/dr2-falling-outcomes.
// Returns {"items":[...]} — never null.
func (h *Handler) ListFallingPriceDR2Outcomes(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListFallingPriceDR2Outcomes(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.FallingPriceDR2Outcome{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createNormalCapitalPerAcreRequest is the body for POST /v1/rent/normal-capital.
type createNormalCapitalPerAcreRequest struct {
	PeriodLabel      string `json:"period_label"`
	CapitalPerAcreBP int64  `json:"capital_per_acre_bp"`
}

// CreateNormalCapitalPerAcre handles POST /v1/rent/normal-capital.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateNormalCapitalPerAcre(w http.ResponseWriter, r *http.Request) {
	var req createNormalCapitalPerAcreRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.PeriodLabel == "" {
		writeError(w, http.StatusUnprocessableEntity, "period_label must not be empty")
		return
	}
	if req.CapitalPerAcreBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "capital_per_acre_bp must be > 0")
		return
	}
	n := rent.NormalCapitalPerAcre{
		PeriodLabel:      req.PeriodLabel,
		CapitalPerAcreBP: req.CapitalPerAcreBP,
	}
	saved, err := h.Store.CreateNormalCapitalPerAcre(r.Context(), n)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/normal-capital")
	writeJSON(w, http.StatusCreated, saved)
}

// ListNormalCapitalPerAcre handles GET /v1/rent/normal-capital.
// Returns {"items":[...]} — never null.
func (h *Handler) ListNormalCapitalPerAcre(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListNormalCapitalPerAcre(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.NormalCapitalPerAcre{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
