package httpapi

import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// billOfExchangeRequest is the body for POST /v1/credit/bills-of-exchange.
type billOfExchangeRequest struct {
	Drawer       string `json:"drawer"`
	Drawee       string `json:"drawee"`
	FaceValue    int64  `json:"face_value"`
	MaturityDays int64  `json:"maturity_days"`
	Description  string `json:"description"`
}

// CreateBillOfExchange handles POST /v1/credit/bills-of-exchange.
// Returns 400 on bad JSON, 422 on validation errors, 201 + Location on success.
func (h *Handler) CreateBillOfExchange(w http.ResponseWriter, r *http.Request) {
	var req billOfExchangeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	bill, err := credit.NewBillOfExchange(req.Drawer, req.Drawee, req.FaceValue, req.MaturityDays, req.Description)
	if err != nil {
		if errors.Is(err, credit.ErrNonPositiveFaceValue) ||
			errors.Is(err, credit.ErrNonPositiveMaturityDays) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateBillOfExchange(r.Context(), bill)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/bills-of-exchange/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListBillsOfExchange handles GET /v1/credit/bills-of-exchange.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListBillsOfExchange(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListBillsOfExchange(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.BillOfExchange{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetBillOfExchange handles GET /v1/credit/bills-of-exchange/{id}.
func (h *Handler) GetBillOfExchange(w http.ResponseWriter, r *http.Request) {
	id := credit.BillOfExchangeID(r.PathValue("id"))
	bill, err := h.Store.GetBillOfExchange(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, bill)
}

// fictitiousCapitalRequest is the body for POST /v1/credit/fictitious-capital.
type fictitiousCapitalRequest struct {
	Kind         credit.FictitiousKind `json:"kind"`
	Name         string                `json:"name"`
	AnnualIncome int64                 `json:"annual_income"`
	RateBP       int64                 `json:"rate_bp"`
	Description  string                `json:"description"`
}

// CreateFictitiousCapital handles POST /v1/credit/fictitious-capital.
// CapitalisedValue and RealCapitalExists are derived server-side from the
// inputs (Vol. III Ch. 25). Returns 400 on bad JSON, 422 on validation
// errors, 201 + Location on success.
func (h *Handler) CreateFictitiousCapital(w http.ResponseWriter, r *http.Request) {
	var req fictitiousCapitalRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	fc, err := credit.NewFictitiousCapital(req.Kind, req.Name, req.AnnualIncome, req.RateBP, req.Description)
	if err != nil {
		if errors.Is(err, credit.ErrNonPositiveAnnualIncome) ||
			errors.Is(err, credit.ErrNonPositiveRateBPFict) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.Store.CreateFictitiousCapital(r.Context(), fc)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/credit/fictitious-capital/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, saved)
}

// ListFictitiousCapitals handles GET /v1/credit/fictitious-capital.
// Returns a JSON object with an "items" array — never null.
func (h *Handler) ListFictitiousCapitals(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListFictitiousCapitals(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if items == nil {
		items = []credit.FictitiousCapital{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetFictitiousCapital handles GET /v1/credit/fictitious-capital/{id}.
func (h *Handler) GetFictitiousCapital(w http.ResponseWriter, r *http.Request) {
	id := credit.FictitiousCapitalID(r.PathValue("id"))
	fc, err := h.Store.GetFictitiousCapital(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, fc)
}
