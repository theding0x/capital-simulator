package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createRisingPriceDR2OutcomeRequest is the body for POST /v1/rent/rising-price-outcomes.
type createRisingPriceDR2OutcomeRequest struct {
	Variant                 int   `json:"variant"`
	TotalCapitalShillings   int64 `json:"total_capital_shillings"`
	TotalRentShillings      int64 `json:"total_rent_shillings"`
	TotalOutputBushels      int64 `json:"total_output_bushels"`
	PricePerBushelShillings int64 `json:"price_per_bushel_shillings"`
}

// CreateRisingPriceDR2Outcome handles POST /v1/rent/rising-price-outcomes.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateRisingPriceDR2Outcome(w http.ResponseWriter, r *http.Request) {
	var req createRisingPriceDR2OutcomeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	variant := rent.RisingPriceVariant(req.Variant)
	if !variant.IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "variant must be 1–5")
		return
	}
	if req.TotalCapitalShillings <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "total_capital_shillings must be > 0")
		return
	}
	if req.TotalOutputBushels <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "total_output_bushels must be > 0")
		return
	}
	if req.PricePerBushelShillings <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "price_per_bushel_shillings must be > 0")
		return
	}
	if req.TotalRentShillings < 0 {
		writeError(w, http.StatusUnprocessableEntity, "total_rent_shillings must be >= 0")
		return
	}
	o := rent.RisingPriceDR2Outcome{
		Variant:                 variant,
		TotalCapitalShillings:   req.TotalCapitalShillings,
		TotalRentShillings:      req.TotalRentShillings,
		TotalOutputBushels:      req.TotalOutputBushels,
		PricePerBushelShillings: req.PricePerBushelShillings,
	}
	saved, err := h.Store.CreateRisingPriceDR2Outcome(r.Context(), o)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/rising-price-outcomes")
	writeJSON(w, http.StatusCreated, saved)
}

// ListRisingPriceDR2Outcomes handles GET /v1/rent/rising-price-outcomes.
// Returns {"items":[...]} — never null.
func (h *Handler) ListRisingPriceDR2Outcomes(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListRisingPriceDR2Outcomes(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.RisingPriceDR2Outcome{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createRentSequenceRequest is the body for POST /v1/rent/rent-sequences.
type createRentSequenceRequest struct {
	BaseRent  int64 `json:"base_rent"`
	Increment int64 `json:"increment"`
	NumGrades int   `json:"num_grades"`
}

// rentSequenceResponse embeds a stored RentSequence and adds the computed sequence slice.
type rentSequenceResponse struct {
	rent.RentSequence
	Sequence []int64 `json:"sequence"`
}

// CreateRentSequence handles POST /v1/rent/rent-sequences.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
// Response includes the stored record fields plus a "sequence" array of computed rents.
func (h *Handler) CreateRentSequence(w http.ResponseWriter, r *http.Request) {
	var req createRentSequenceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.NumGrades < 1 {
		writeError(w, http.StatusUnprocessableEntity, "num_grades must be >= 1")
		return
	}
	if req.Increment < 0 {
		writeError(w, http.StatusUnprocessableEntity, "increment must be >= 0")
		return
	}
	seq := rent.ComputeRentSequence(req.BaseRent, req.Increment, req.NumGrades)
	rs := rent.RentSequence{
		BaseRent:  req.BaseRent,
		Increment: req.Increment,
		NumGrades: req.NumGrades,
	}
	saved, err := h.Store.CreateRentSequence(r.Context(), rs)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/rent-sequences")
	writeJSON(w, http.StatusCreated, rentSequenceResponse{RentSequence: saved, Sequence: seq})
}

// ListRentSequences handles GET /v1/rent/rent-sequences.
// Returns {"items":[...]} — never null. Items are bare RentSequence records.
func (h *Handler) ListRentSequences(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListRentSequences(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.RentSequence{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createEngelsTableEntryRequest is the body for POST /v1/rent/engels-tables.
type createEngelsTableEntryRequest struct {
	TableNumber      int    `json:"table_number"`
	SoilGradeLabel   string `json:"soil_grade_label"`
	OutputBushels    int64  `json:"output_bushels"`
	RentShillings    int64  `json:"rent_shillings"`
	CapitalShillings int64  `json:"capital_shillings"`
}

// CreateEngelsTableEntry handles POST /v1/rent/engels-tables.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateEngelsTableEntry(w http.ResponseWriter, r *http.Request) {
	var req createEngelsTableEntryRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.TableNumber < 1 {
		writeError(w, http.StatusUnprocessableEntity, "table_number must be >= 1")
		return
	}
	if req.SoilGradeLabel == "" {
		writeError(w, http.StatusUnprocessableEntity, "soil_grade_label must not be empty")
		return
	}
	if req.OutputBushels <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "output_bushels must be > 0")
		return
	}
	if req.CapitalShillings <= 0 {
		writeError(w, http.StatusUnprocessableEntity, "capital_shillings must be > 0")
		return
	}
	if req.RentShillings < 0 {
		writeError(w, http.StatusUnprocessableEntity, "rent_shillings must be >= 0")
		return
	}
	e := rent.EngelsTableEntry{
		TableNumber:      req.TableNumber,
		SoilGradeLabel:   req.SoilGradeLabel,
		OutputBushels:    req.OutputBushels,
		RentShillings:    req.RentShillings,
		CapitalShillings: req.CapitalShillings,
	}
	saved, err := h.Store.CreateEngelsTableEntry(r.Context(), e)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/engels-tables")
	writeJSON(w, http.StatusCreated, saved)
}

// ListEngelsTableEntries handles GET /v1/rent/engels-tables.
// Returns {"items":[...]} — never null.
func (h *Handler) ListEngelsTableEntries(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListEngelsTableEntries(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.EngelsTableEntry{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
