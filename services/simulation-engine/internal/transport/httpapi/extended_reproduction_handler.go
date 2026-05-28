package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Vol. II Ch. 21 — Accumulation and Reproduction on an Extended Scale.

// --- request types ---

type createExtendedReproductionSchemeRequest struct {
	Period                string `json:"period"`
	AccumulationRateIBps  int64  `json:"accumulation_rate_i_bps"`
	AccumulationRateIIBps int64  `json:"accumulation_rate_ii_bps"`
}

type createReinvestmentRequest struct {
	Department           string `json:"department"`
	DeltaConstantPence   int64  `json:"delta_constant_pence"`
	DeltaVariablePence   int64  `json:"delta_variable_pence"`
	ConsumedSurplusPence int64  `json:"consumed_surplus_pence"`
	CycleNumber          int64  `json:"cycle_number"`
}

type createMultiPeriodSchemeRequest struct {
	Label     string   `json:"label"`
	SchemeIDs []string `json:"scheme_ids"`
}

// --- response types ---

type extendedReproductionSchemeResponse struct {
	ID                    string                       `json:"id"`
	Period                string                       `json:"period"`
	DepartmentI           *departmentalCapitalResponse `json:"department_i,omitempty"`
	DepartmentII          *departmentalCapitalResponse `json:"department_ii,omitempty"`
	AccumulationRateIBps  int64                        `json:"accumulation_rate_i_bps"`
	AccumulationRateIIBps int64                        `json:"accumulation_rate_ii_bps"`
	IsBalanced            bool                         `json:"is_balanced"`
	TickCount             int64                        `json:"tick_count"`
}

type reinvestmentResponse struct {
	ID                   string `json:"id"`
	SchemeID             string `json:"scheme_id"`
	Department           string `json:"department"`
	DeltaConstantPence   int64  `json:"delta_constant_pence"`
	DeltaVariablePence   int64  `json:"delta_variable_pence"`
	ConsumedSurplusPence int64  `json:"consumed_surplus_pence"`
	CycleNumber          int64  `json:"cycle_number"`
}

type multiPeriodSchemeResponse struct {
	ID        string   `json:"id"`
	Label     string   `json:"label"`
	SchemeIDs []string `json:"scheme_ids"`
}

type compositionShiftResponse struct {
	ID                   string `json:"id"`
	MultiPeriodID        string `json:"multi_period_id"`
	CycleNumber          int64  `json:"cycle_number"`
	DeptICompositionBps  int64  `json:"dept_i_composition_bps"`
	DeptIICompositionBps int64  `json:"dept_ii_composition_bps"`
	SocialCompositionBps int64  `json:"social_composition_bps"`
}

type extendedMoneyRequirementResponse struct {
	ID                   string `json:"id"`
	SchemeID             string `json:"scheme_id"`
	BaseMoneyPence       int64  `json:"base_money_pence"`
	AdditionalMoneyPence int64  `json:"additional_money_pence"`
	TotalMoneyPence      int64  `json:"total_money_pence"`
}

type departmentIGrowthLeadResponse struct {
	ID              string `json:"id"`
	MultiPeriodID   string `json:"multi_period_id"`
	CycleNumber     int64  `json:"cycle_number"`
	DeptIGrowthBps  int64  `json:"dept_i_growth_bps"`
	DeptIIGrowthBps int64  `json:"dept_ii_growth_bps"`
	LeadConfirmed   bool   `json:"lead_confirmed"`
}

// --- helpers ---

func extendedSchemeToResponse(s repro.ExtendedReproductionScheme) extendedReproductionSchemeResponse {
	resp := extendedReproductionSchemeResponse{
		ID:                    string(s.ID),
		Period:                s.Period,
		AccumulationRateIBps:  s.AccumulationRateIBps,
		AccumulationRateIIBps: s.AccumulationRateIIBps,
		IsBalanced:            s.IsBalanced,
		TickCount:             s.TickCount,
	}
	if s.DepartmentI != nil {
		r := deptCapToResponse(*s.DepartmentI)
		resp.DepartmentI = &r
	}
	if s.DepartmentII != nil {
		r := deptCapToResponse(*s.DepartmentII)
		resp.DepartmentII = &r
	}
	return resp
}

func reinvestmentToResponse(r repro.Reinvestment) reinvestmentResponse {
	return reinvestmentResponse{
		ID:                   string(r.ID),
		SchemeID:             string(r.SchemeID),
		Department:           string(r.Department),
		DeltaConstantPence:   r.DeltaConstantPence,
		DeltaVariablePence:   r.DeltaVariablePence,
		ConsumedSurplusPence: r.ConsumedSurplusPence,
		CycleNumber:          r.CycleNumber,
	}
}

func multiPeriodToResponse(mp repro.MultiPeriodScheme) multiPeriodSchemeResponse {
	ids := make([]string, len(mp.SchemeIDs))
	for i, sid := range mp.SchemeIDs {
		ids[i] = string(sid)
	}
	return multiPeriodSchemeResponse{
		ID:        string(mp.ID),
		Label:     mp.Label,
		SchemeIDs: ids,
	}
}

// --- handlers ---

// handleCreateExtendedReproductionScheme handles POST /v1/reproduction/extended/schemes.
func (h *Handler) handleCreateExtendedReproductionScheme(w http.ResponseWriter, r *http.Request) {
	var req createExtendedReproductionSchemeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	s := repro.ExtendedReproductionScheme{
		ID:                    repro.NewExtendedReproductionSchemeID(),
		Period:                req.Period,
		AccumulationRateIBps:  req.AccumulationRateIBps,
		AccumulationRateIIBps: req.AccumulationRateIIBps,
	}
	created, err := h.ExtendedReproduction.CreateExtendedReproductionScheme(r.Context(), s)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "scheme already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/reproduction/extended/schemes/"+string(created.ID))
	writeJSON(w, http.StatusCreated, extendedSchemeToResponse(created))
}

// handleListExtendedReproductionSchemes handles GET /v1/reproduction/extended/schemes.
func (h *Handler) handleListExtendedReproductionSchemes(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	list, err := h.ExtendedReproduction.ListExtendedReproductionSchemes(r.Context(), period)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	resp := make([]extendedReproductionSchemeResponse, len(list))
	for i, s := range list {
		resp[i] = extendedSchemeToResponse(s)
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleGetExtendedReproductionScheme handles GET /v1/reproduction/extended/schemes/{id}.
func (h *Handler) handleGetExtendedReproductionScheme(w http.ResponseWriter, r *http.Request) {
	id := repro.ExtendedReproductionSchemeID(r.PathValue("id"))
	s, err := h.ExtendedReproduction.GetExtendedReproductionScheme(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "scheme not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, extendedSchemeToResponse(s))
}

// handleAddExtendedDepartmentToScheme handles POST /v1/reproduction/extended/schemes/{id}/departments.
func (h *Handler) handleAddExtendedDepartmentToScheme(w http.ResponseWriter, r *http.Request) {
	id := repro.ExtendedReproductionSchemeID(r.PathValue("id"))
	var req addDepartmentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	d := repro.DepartmentalCapital{
		ID:            repro.NewDepartmentalCapitalID(),
		Department:    repro.CapitalDepartment(req.Department),
		ConstantPence: req.ConstantPence,
		VariablePence: req.VariablePence,
		SurplusPence:  req.SurplusPence,
		TotalPence:    req.TotalPence,
	}
	if err := d.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.ExtendedReproduction.AddExtendedDepartmentToScheme(r.Context(), id, d)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeError(w, http.StatusNotFound, "scheme not found")
		case errors.Is(err, store.ErrAlreadyExists):
			writeError(w, http.StatusConflict, "department already registered on this scheme")
		default:
			h.writeServerError(w, err)
		}
		return
	}
	writeJSON(w, http.StatusCreated, deptCapToResponse(created))
}

// handleCreateReinvestment handles POST /v1/reproduction/extended/schemes/{id}/reinvestments.
func (h *Handler) handleCreateReinvestment(w http.ResponseWriter, r *http.Request) {
	id := repro.ExtendedReproductionSchemeID(r.PathValue("id"))
	var req createReinvestmentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	reinvest := repro.Reinvestment{
		ID:                   repro.NewReinvestmentID(),
		Department:           repro.CapitalDepartment(req.Department),
		DeltaConstantPence:   req.DeltaConstantPence,
		DeltaVariablePence:   req.DeltaVariablePence,
		ConsumedSurplusPence: req.ConsumedSurplusPence,
		CycleNumber:          req.CycleNumber,
	}
	created, err := h.ExtendedReproduction.CreateReinvestment(r.Context(), id, reinvest)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "scheme not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, reinvestmentToResponse(created))
}

// handleTickExtendedScheme handles POST /v1/reproduction/extended/schemes/{id}/tick.
func (h *Handler) handleTickExtendedScheme(w http.ResponseWriter, r *http.Request) {
	id := repro.ExtendedReproductionSchemeID(r.PathValue("id"))
	s, err := h.ExtendedReproduction.TickExtendedScheme(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeError(w, http.StatusNotFound, "scheme not found or departments missing")
		case errors.Is(err, repro.ErrNegativeAccumulationRate),
			errors.Is(err, repro.ErrAccumulationRateExceeds100Pct),
			errors.Is(err, repro.ErrExtendedImbalance):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			h.writeServerError(w, err)
		}
		return
	}
	writeJSON(w, http.StatusCreated, extendedSchemeToResponse(s))
}

// handleGetExtendedMoneyRequirement handles GET /v1/reproduction/extended/schemes/{id}/money-requirement.
func (h *Handler) handleGetExtendedMoneyRequirement(w http.ResponseWriter, r *http.Request) {
	id := repro.ExtendedReproductionSchemeID(r.PathValue("id"))
	basePenceStr := r.URL.Query().Get("base_money_pence")
	var baseMoneyPence int64
	if basePenceStr != "" {
		v, err := strconv.ParseInt(basePenceStr, 10, 64)
		if err != nil || v < 0 {
			writeError(w, http.StatusBadRequest, "base_money_pence must be a non-negative integer")
			return
		}
		baseMoneyPence = v
	}
	req, err := h.ExtendedReproduction.GetExtendedMoneyRequirement(r.Context(), id, baseMoneyPence)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeError(w, http.StatusNotFound, "scheme not found")
		case errors.Is(err, repro.ErrExtendedImbalance):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			h.writeServerError(w, err)
		}
		return
	}
	writeJSON(w, http.StatusOK, extendedMoneyRequirementResponse{
		ID:                   string(req.ID),
		SchemeID:             string(req.SchemeID),
		BaseMoneyPence:       req.BaseMoneyPence,
		AdditionalMoneyPence: req.AdditionalMoneyPence,
		TotalMoneyPence:      req.TotalMoneyPence,
	})
}

// handleCreateMultiPeriodScheme handles POST /v1/reproduction/extended/multi-period.
func (h *Handler) handleCreateMultiPeriodScheme(w http.ResponseWriter, r *http.Request) {
	var req createMultiPeriodSchemeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Label == "" {
		writeError(w, http.StatusBadRequest, "label is required")
		return
	}
	schemeIDs := make([]repro.ExtendedReproductionSchemeID, len(req.SchemeIDs))
	for i, sid := range req.SchemeIDs {
		schemeIDs[i] = repro.ExtendedReproductionSchemeID(sid)
	}
	mp := repro.MultiPeriodScheme{
		ID:        repro.NewMultiPeriodSchemeID(),
		Label:     req.Label,
		SchemeIDs: schemeIDs,
	}
	created, err := h.ExtendedReproduction.CreateMultiPeriodScheme(r.Context(), mp)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeError(w, http.StatusNotFound, "one or more linked schemes not found")
		case errors.Is(err, store.ErrAlreadyExists):
			writeError(w, http.StatusConflict, "multi-period scheme already exists")
		default:
			h.writeServerError(w, err)
		}
		return
	}
	w.Header().Set("Location", "/v1/reproduction/extended/multi-period/"+string(created.ID))
	writeJSON(w, http.StatusCreated, multiPeriodToResponse(created))
}

// handleGetCompositionShift handles GET /v1/reproduction/extended/multi-period/{id}/composition-shift.
func (h *Handler) handleGetCompositionShift(w http.ResponseWriter, r *http.Request) {
	id := repro.MultiPeriodSchemeID(r.PathValue("id"))
	shifts, err := h.ExtendedReproduction.ListCompositionShifts(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "multi-period scheme not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	resp := make([]compositionShiftResponse, len(shifts))
	for i, cs := range shifts {
		resp[i] = compositionShiftResponse{
			ID:                   string(cs.ID),
			MultiPeriodID:        string(cs.MultiPeriodID),
			CycleNumber:          cs.CycleNumber,
			DeptICompositionBps:  cs.DeptICompositionBps,
			DeptIICompositionBps: cs.DeptIICompositionBps,
			SocialCompositionBps: cs.SocialCompositionBps,
		}
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleGetGrowthLead handles GET /v1/reproduction/extended/multi-period/{id}/growth-lead.
func (h *Handler) handleGetGrowthLead(w http.ResponseWriter, r *http.Request) {
	id := repro.MultiPeriodSchemeID(r.PathValue("id"))
	leads, err := h.ExtendedReproduction.ListGrowthLeads(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "multi-period scheme not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	resp := make([]departmentIGrowthLeadResponse, len(leads))
	for i, gl := range leads {
		resp[i] = departmentIGrowthLeadResponse{
			ID:              string(gl.ID),
			MultiPeriodID:   string(gl.MultiPeriodID),
			CycleNumber:     gl.CycleNumber,
			DeptIGrowthBps:  gl.DeptIGrowthBps,
			DeptIIGrowthBps: gl.DeptIIGrowthBps,
			LeadConfirmed:   gl.LeadConfirmed,
		}
	}
	writeJSON(w, http.StatusOK, resp)
}
