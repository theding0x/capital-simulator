package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createDR2OutcomeRequest is the body for POST /v1/rent/dr2-outcomes.
type createDR2OutcomeRequest struct {
	Case        int                        `json:"case"`
	Investments []dr2InvestmentRequestItem `json:"investments"`
}

type dr2InvestmentRequestItem struct {
	ParcelID        string `json:"parcel_id"`
	TrancheNumber   int    `json:"tranche_number"`
	CapitalBP       int64  `json:"capital_bp"`
	OutputQuarters  int64  `json:"output_quarters"`
	SurplusProfitBP int64  `json:"surplus_profit_bp"`
}

// CreateDR2Outcome handles POST /v1/rent/dr2-outcomes.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateDR2Outcome(w http.ResponseWriter, r *http.Request) {
	var req createDR2OutcomeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	c := rent.CapitalProductivityCase(req.Case)
	if !c.IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "case must be 1–4")
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
	outcome := rent.ComputeDR2Outcome(c, invs)
	saved, err := h.Store.CreateDifferentialRentIIOutcome(r.Context(), outcome)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/dr2-outcomes")
	writeJSON(w, http.StatusCreated, saved)
}

// ListDR2Outcomes handles GET /v1/rent/dr2-outcomes.
// Returns {"items":[...]} — never null.
func (h *Handler) ListDR2Outcomes(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListDifferentialRentIIOutcomes(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.DifferentialRentIIOutcome{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createIntensiveExtensiveRequest is the body for POST /v1/rent/intensive-extensive.
type createIntensiveExtensiveRequest struct {
	IntensiveRentPerAcreBP int64 `json:"intensive_rent_per_acre_bp"`
	ExtensiveRentPerAcreBP int64 `json:"extensive_rent_per_acre_bp"`
	TotalCapitalBP         int64 `json:"total_capital_bp"`
	TotalRentalSameBP      int64 `json:"total_rental_same_bp"`
}

// CreateIntensiveExtensiveComparison handles POST /v1/rent/intensive-extensive.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateIntensiveExtensiveComparison(w http.ResponseWriter, r *http.Request) {
	var req createIntensiveExtensiveRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.IntensiveRentPerAcreBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "intensive_rent_per_acre_bp must be > 0")
		return
	}
	if req.ExtensiveRentPerAcreBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "extensive_rent_per_acre_bp must be > 0")
		return
	}
	if req.TotalCapitalBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "total_capital_bp must be > 0")
		return
	}
	if req.TotalRentalSameBP <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "total_rental_same_bp must be > 0")
		return
	}
	if req.IntensiveRentPerAcreBP <= req.ExtensiveRentPerAcreBP {
		writeError(w, http.StatusUnprocessableEntity, "intensive_rent_per_acre_bp must be > extensive_rent_per_acre_bp")
		return
	}
	cmp := rent.IntensiveExtensiveComparison{
		IntensiveRentPerAcreBP: req.IntensiveRentPerAcreBP,
		ExtensiveRentPerAcreBP: req.ExtensiveRentPerAcreBP,
		TotalCapitalBP:         req.TotalCapitalBP,
		TotalRentalSameBP:      req.TotalRentalSameBP,
	}
	saved, err := h.Store.CreateIntensiveExtensiveComparison(r.Context(), cmp)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/intensive-extensive")
	writeJSON(w, http.StatusCreated, saved)
}

// ListIntensiveExtensiveComparisons handles GET /v1/rent/intensive-extensive.
// Returns {"items":[...]} — never null.
func (h *Handler) ListIntensiveExtensiveComparisons(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListIntensiveExtensiveComparisons(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.IntensiveExtensiveComparison{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
