package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// createRentHistoricalFormRequest is the body for POST /v1/rent/historical-forms.
type createRentHistoricalFormRequest struct {
	Stage         int    `json:"stage"`
	Name          string `json:"name"`
	CoercionKind  string `json:"coercion_kind"`
	SurplusForm   string `json:"surplus_form"`
	ExampleRegion string `json:"example_region"`
	ExamplePeriod string `json:"example_period"`
}

// CreateRentHistoricalForm handles POST /v1/rent/historical-forms.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateRentHistoricalForm(w http.ResponseWriter, r *http.Request) {
	var req createRentHistoricalFormRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !rent.RentFormStage(req.Stage).IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "stage is not a valid RentFormStage")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusUnprocessableEntity, "name must not be empty")
		return
	}
	rec := rent.RentHistoricalForm{
		Stage:         rent.RentFormStage(req.Stage),
		Name:          req.Name,
		CoercionKind:  req.CoercionKind,
		SurplusForm:   req.SurplusForm,
		ExampleRegion: req.ExampleRegion,
		ExamplePeriod: req.ExamplePeriod,
	}
	saved, err := h.Store.CreateRentHistoricalForm(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/historical-forms")
	writeJSON(w, http.StatusCreated, saved)
}

// ListRentHistoricalForms handles GET /v1/rent/historical-forms.
// Returns {"items":[...]} — never null.
func (h *Handler) ListRentHistoricalForms(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListRentHistoricalForms(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.RentHistoricalForm{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createHistoricalRentTransitionRequest is the body for POST /v1/rent/historical-transitions.
type createHistoricalRentTransitionRequest struct {
	FromStage     int    `json:"from_stage"`
	ToStage       int    `json:"to_stage"`
	Preconditions string `json:"preconditions"`
	DrivingForce  string `json:"driving_force"`
}

// CreateHistoricalRentTransition handles POST /v1/rent/historical-transitions.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateHistoricalRentTransition(w http.ResponseWriter, r *http.Request) {
	var req createHistoricalRentTransitionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !rent.RentFormStage(req.FromStage).IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "from_stage is not a valid RentFormStage")
		return
	}
	if !rent.RentFormStage(req.ToStage).IsValid() {
		writeError(w, http.StatusUnprocessableEntity, "to_stage is not a valid RentFormStage")
		return
	}
	if req.FromStage >= req.ToStage {
		writeError(w, http.StatusUnprocessableEntity, "from_stage must be less than to_stage")
		return
	}
	rec := rent.HistoricalRentTransition{
		FromStage:     rent.RentFormStage(req.FromStage),
		ToStage:       rent.RentFormStage(req.ToStage),
		Preconditions: req.Preconditions,
		DrivingForce:  req.DrivingForce,
	}
	saved, err := h.Store.CreateHistoricalRentTransition(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/historical-transitions")
	writeJSON(w, http.StatusCreated, saved)
}

// ListHistoricalRentTransitions handles GET /v1/rent/historical-transitions.
// Returns {"items":[...]} — never null.
func (h *Handler) ListHistoricalRentTransitions(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListHistoricalRentTransitions(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.HistoricalRentTransition{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// createSmallPeasantProductionRequest is the body for POST /v1/rent/small-peasant-productions.
type createSmallPeasantProductionRequest struct {
	ParcelID          string `json:"parcel_id"`
	RentComponentBP   int64  `json:"rent_component_bp"`
	ProfitComponentBP int64  `json:"profit_component_bp"`
	WagesComponentBP  int64  `json:"wages_component_bp"`
	IsConflated       bool   `json:"is_conflated"`
}

// CreateSmallPeasantProduction handles POST /v1/rent/small-peasant-productions.
// Returns 400 on bad JSON, 422 on invalid fields, 201 + Location on success.
func (h *Handler) CreateSmallPeasantProduction(w http.ResponseWriter, r *http.Request) {
	var req createSmallPeasantProductionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ParcelID == "" {
		writeError(w, http.StatusUnprocessableEntity, "parcel_id must not be empty")
		return
	}
	if req.RentComponentBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "rent_component_bp must be >= 0")
		return
	}
	if req.ProfitComponentBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "profit_component_bp must be >= 0")
		return
	}
	if req.WagesComponentBP < 0 {
		writeError(w, http.StatusUnprocessableEntity, "wages_component_bp must be >= 0")
		return
	}
	total := rent.ComputeTotalIncome(req.RentComponentBP, req.ProfitComponentBP, req.WagesComponentBP)
	rec := rent.SmallPeasantProduction{
		ParcelID:          req.ParcelID,
		TotalIncomeBP:     total,
		RentComponentBP:   req.RentComponentBP,
		ProfitComponentBP: req.ProfitComponentBP,
		WagesComponentBP:  req.WagesComponentBP,
		IsConflated:       req.IsConflated,
	}
	saved, err := h.Store.CreateSmallPeasantProduction(r.Context(), rec)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/rent/small-peasant-productions")
	writeJSON(w, http.StatusCreated, saved)
}

// ListSmallPeasantProductions handles GET /v1/rent/small-peasant-productions.
// Returns {"items":[...]} — never null.
func (h *Handler) ListSmallPeasantProductions(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListSmallPeasantProductions(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []rent.SmallPeasantProduction{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
