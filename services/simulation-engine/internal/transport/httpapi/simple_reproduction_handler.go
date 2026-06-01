package httpapi

import (
	"context"
	"errors"
	"net/http"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Vol. II Ch. 20 — Simple Reproduction.

// --- request types ---

type createSimpleReproductionSchemeRequest struct {
	Period string `json:"period"`
}

type addDepartmentRequest struct {
	Department       string `json:"department"`
	ConstantPence    int64  `json:"constant_pence"`
	VariablePence    int64  `json:"variable_pence"`
	SurplusPence     int64  `json:"surplus_pence"`
	TotalPence       int64  `json:"total_pence"`
	TurnoverNumberBP int64  `json:"turnover_number_bp"`
}

type recordExchangeRequest struct {
	FromDepartment string `json:"from_department"`
	ToDepartment   string `json:"to_department"`
	Pence          int64  `json:"pence"`
	Kind           string `json:"kind"`
	Description    string `json:"description"`
}

type recordMoneyLoopRequest struct {
	Period       string `json:"period"`
	IsClosed     bool   `json:"is_closed"`
	NetFlowPence int64  `json:"net_flow_pence"`
}

// advanceTickRequest is the optional body for a tick. When worker_pool_size > 0
// the tick asserts the wage bill covers the pool's subsistence (issue #220).
type advanceTickRequest struct {
	WorkerPoolSize         int64 `json:"worker_pool_size"`
	SubsistenceBasketPence int64 `json:"subsistence_basket_pence"`
}

// --- response types ---

type departmentalCapitalResponse struct {
	ID               string `json:"id"`
	SchemeID         string `json:"scheme_id"`
	Department       string `json:"department"`
	ConstantPence    int64  `json:"constant_pence"`
	VariablePence    int64  `json:"variable_pence"`
	SurplusPence     int64  `json:"surplus_pence"`
	TotalPence       int64  `json:"total_pence"`
	TurnoverNumberBP int64  `json:"turnover_number_bp"`
}

type simpleReproductionSchemeResponse struct {
	ID           string                       `json:"id"`
	Period       string                       `json:"period"`
	DepartmentI  *departmentalCapitalResponse `json:"department_i,omitempty"`
	DepartmentII *departmentalCapitalResponse `json:"department_ii,omitempty"`
	IsBalanced   bool                         `json:"is_balanced"`
	TickCount    int64                        `json:"tick_count"`
}

type interdepartmentExchangeResponse struct {
	ID             string `json:"id"`
	SchemeID       string `json:"scheme_id"`
	FromDepartment string `json:"from_department"`
	ToDepartment   string `json:"to_department"`
	Pence          int64  `json:"pence"`
	Kind           string `json:"kind"`
	Description    string `json:"description"`
}

type reproductionTickResponse struct {
	ID                     string `json:"id"`
	SchemeID               string `json:"scheme_id"`
	TickNumber             int64  `json:"tick_number"`
	Period                 string `json:"period"`
	IsBalanced             bool   `json:"is_balanced"`
	WorkerPoolSize         int64  `json:"worker_pool_size"`
	WageBillPence          int64  `json:"wage_bill_pence"`
	SubsistenceBasketPence int64  `json:"subsistence_basket_pence"`
	SubsistenceCovered     bool   `json:"subsistence_covered"`
	Status                 string `json:"status"`
}

type balanceCheckResponse struct {
	SchemeID            string `json:"scheme_id"`
	IsBalanced          bool   `json:"is_balanced"`
	DeptIVplusSPence    int64  `json:"dept_i_v_plus_s_pence"`
	DeptIIConstantPence int64  `json:"dept_ii_constant_pence"`
	DeficitPence        int64  `json:"deficit_pence"`

	DeptITurnoverBP               int64 `json:"dept_i_turnover_bp"`
	DeptIITurnoverBP              int64 `json:"dept_ii_turnover_bp"`
	DeptIVplusSAnnualisedPence    int64 `json:"dept_i_v_plus_s_annualised_pence"`
	DeptIIConstantAnnualisedPence int64 `json:"dept_ii_constant_annualised_pence"`
	AnnualisedDeficitPence        int64 `json:"annualised_deficit_pence"`
	AnnualisedBalanced            bool  `json:"annualised_balanced"`
}

type moneyClosedLoopResponse struct {
	ID           string `json:"id"`
	SchemeID     string `json:"scheme_id"`
	Period       string `json:"period"`
	IsClosed     bool   `json:"is_closed"`
	NetFlowPence int64  `json:"net_flow_pence"`
}

// --- helpers ---

func deptCapToResponse(d repro.DepartmentalCapital) departmentalCapitalResponse {
	return departmentalCapitalResponse{
		ID:               string(d.ID),
		SchemeID:         string(d.SchemeID),
		Department:       string(d.Department),
		ConstantPence:    d.ConstantPence,
		VariablePence:    d.VariablePence,
		SurplusPence:     d.SurplusPence,
		TotalPence:       d.TotalPence,
		TurnoverNumberBP: d.TurnoverNumberBP,
	}
}

func schemeToResponse(s repro.SimpleReproductionScheme) simpleReproductionSchemeResponse {
	resp := simpleReproductionSchemeResponse{
		ID:         string(s.ID),
		Period:     s.Period,
		IsBalanced: s.IsBalanced,
		TickCount:  s.TickCount,
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

func exchangeToResponse(e repro.InterDepartmentExchange) interdepartmentExchangeResponse {
	return interdepartmentExchangeResponse{
		ID:             string(e.ID),
		SchemeID:       string(e.SchemeID),
		FromDepartment: string(e.FromDepartment),
		ToDepartment:   string(e.ToDepartment),
		Pence:          e.Pence,
		Kind:           string(e.Kind),
		Description:    e.Description,
	}
}

// projectDepartmentShares recomputes the Department I/II share split (basis
// points) from a reproduction scheme's two departments and persists it onto the
// period's SocialCapitalAggregate, wiring Ch. 20/21 reproduction back into the
// Ch. 17 aggregate placeholders (issue #215). A scheme missing either
// department — or with no advanced capital — is a no-op.
func (h *Handler) projectDepartmentShares(ctx context.Context, period string, deptI, deptII *repro.DepartmentalCapital) error {
	if deptI == nil || deptII == nil {
		return nil
	}
	deptIShareBP, deptIIShareBP := repro.ComputeDepartmentSharesBP(deptI.AdvancedPence(), deptII.AdvancedPence())
	if deptIShareBP == 0 && deptIIShareBP == 0 {
		return nil
	}
	_, err := h.SurplusCirculations.UpsertAggregateShares(ctx, period, deptIShareBP, deptIIShareBP)
	return err
}

// --- handlers ---

// CreateSimpleReproductionScheme handles POST /v1/reproduction/simple/schemes.
func (h *Handler) CreateSimpleReproductionScheme(w http.ResponseWriter, r *http.Request) {
	var req createSimpleReproductionSchemeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	s := repro.SimpleReproductionScheme{
		ID:     repro.NewSimpleReproductionSchemeID(),
		Period: req.Period,
	}
	created, err := h.SimpleReproduction.CreateSimpleReproductionScheme(r.Context(), s)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "scheme already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/reproduction/simple/schemes/"+string(created.ID))
	writeJSON(w, http.StatusCreated, schemeToResponse(created))
}

// GetSimpleReproductionScheme handles GET /v1/reproduction/simple/schemes/{id}.
func (h *Handler) GetSimpleReproductionScheme(w http.ResponseWriter, r *http.Request) {
	id := repro.SimpleReproductionSchemeID(r.PathValue("id"))
	s, err := h.SimpleReproduction.GetSimpleReproductionScheme(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "scheme not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, schemeToResponse(s))
}

// ListSimpleReproductionSchemes handles GET /v1/reproduction/simple/schemes.
func (h *Handler) ListSimpleReproductionSchemes(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	list, err := h.SimpleReproduction.ListSimpleReproductionSchemes(r.Context(), period)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	resp := make([]simpleReproductionSchemeResponse, len(list))
	for i, s := range list {
		resp[i] = schemeToResponse(s)
	}
	writeJSON(w, http.StatusOK, resp)
}

// AddDepartmentToScheme handles POST /v1/reproduction/simple/schemes/{id}/departments.
func (h *Handler) AddDepartmentToScheme(w http.ResponseWriter, r *http.Request) {
	id := repro.SimpleReproductionSchemeID(r.PathValue("id"))
	var req addDepartmentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	turnoverBP := req.TurnoverNumberBP
	if turnoverBP <= 0 {
		turnoverBP = 10000 // default: one turnover per year (annual)
	}
	d := repro.DepartmentalCapital{
		ID:               repro.NewDepartmentalCapitalID(),
		Department:       repro.CapitalDepartment(req.Department),
		ConstantPence:    req.ConstantPence,
		VariablePence:    req.VariablePence,
		SurplusPence:     req.SurplusPence,
		TotalPence:       req.TotalPence,
		TurnoverNumberBP: turnoverBP,
	}
	if err := d.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.SimpleReproduction.AddDepartmentToScheme(r.Context(), id, d)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			writeError(w, http.StatusNotFound, "scheme not found")
		case errors.Is(err, repro.ErrDepartmentAlreadySet):
			writeError(w, http.StatusConflict, "department already registered on this scheme")
		default:
			h.writeServerError(w, err)
		}
		return
	}
	writeJSON(w, http.StatusCreated, deptCapToResponse(created))
}

// RecordInterDepartmentExchange handles POST /v1/reproduction/simple/schemes/{id}/exchanges.
func (h *Handler) RecordInterDepartmentExchange(w http.ResponseWriter, r *http.Request) {
	id := repro.SimpleReproductionSchemeID(r.PathValue("id"))
	var req recordExchangeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	e := repro.InterDepartmentExchange{
		FromDepartment: repro.CapitalDepartment(req.FromDepartment),
		ToDepartment:   repro.CapitalDepartment(req.ToDepartment),
		Pence:          req.Pence,
		Kind:           repro.ExchangeKind(req.Kind),
		Description:    req.Description,
	}
	if err := e.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.SimpleReproduction.RecordInterDepartmentExchange(r.Context(), id, e)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "scheme not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, exchangeToResponse(created))
}

// AdvanceReproductionTick handles POST /v1/reproduction/simple/schemes/{id}/tick.
func (h *Handler) AdvanceReproductionTick(w http.ResponseWriter, r *http.Request) {
	id := repro.SimpleReproductionSchemeID(r.PathValue("id"))
	var req advanceTickRequest
	if r.ContentLength != 0 {
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	tick, err := h.SimpleReproduction.AdvanceReproductionTick(r.Context(), id, req.WorkerPoolSize, req.SubsistenceBasketPence)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "scheme not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	// Project the two departments' advanced capital back onto the period's
	// social-capital aggregate (issue #215).
	scheme, err := h.SimpleReproduction.GetSimpleReproductionScheme(r.Context(), id)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	if err := h.projectDepartmentShares(r.Context(), scheme.Period, scheme.DepartmentI, scheme.DepartmentII); err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, reproductionTickResponse{
		ID:                     string(tick.ID),
		SchemeID:               string(tick.SchemeID),
		TickNumber:             tick.TickNumber,
		Period:                 tick.Period,
		IsBalanced:             tick.IsBalanced,
		WorkerPoolSize:         tick.WorkerPoolSize,
		WageBillPence:          tick.WageBillPence,
		SubsistenceBasketPence: tick.SubsistenceBasketPence,
		SubsistenceCovered:     tick.SubsistenceCovered,
		Status:                 tick.Status,
	})
}

// GetSchemeBalanceCheck handles GET /v1/reproduction/simple/schemes/{id}/balance-check.
func (h *Handler) GetSchemeBalanceCheck(w http.ResponseWriter, r *http.Request) {
	id := repro.SimpleReproductionSchemeID(r.PathValue("id"))
	res, err := h.SimpleReproduction.CheckSchemeBalance(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "scheme not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, balanceCheckResponse{
		SchemeID:                      res.SchemeID,
		IsBalanced:                    res.IsBalanced,
		DeptIVplusSPence:              res.DeptIVplusSPence,
		DeptIIConstantPence:           res.DeptIIConstantPence,
		DeficitPence:                  res.DeficitPence,
		DeptITurnoverBP:               res.DeptITurnoverBP,
		DeptIITurnoverBP:              res.DeptIITurnoverBP,
		DeptIVplusSAnnualisedPence:    res.DeptIVplusSAnnualisedPence,
		DeptIIConstantAnnualisedPence: res.DeptIIConstantAnnualisedPence,
		AnnualisedDeficitPence:        res.AnnualisedDeficitPence,
		AnnualisedBalanced:            res.AnnualisedBalanced,
	})
}

// RecordMoneyClosedLoop handles POST /v1/reproduction/simple/schemes/{id}/money-loop.
func (h *Handler) RecordMoneyClosedLoop(w http.ResponseWriter, r *http.Request) {
	id := repro.SimpleReproductionSchemeID(r.PathValue("id"))
	var req recordMoneyLoopRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	loop := repro.MoneyClosedLoop{
		ID:           repro.NewMoneyClosedLoopID(),
		Period:       req.Period,
		IsClosed:     req.IsClosed,
		NetFlowPence: req.NetFlowPence,
	}
	created, err := h.SimpleReproduction.RecordMoneyClosedLoop(r.Context(), id, loop)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "scheme not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, moneyClosedLoopResponse{
		ID:           string(created.ID),
		SchemeID:     string(created.SchemeID),
		Period:       created.Period,
		IsClosed:     created.IsClosed,
		NetFlowPence: created.NetFlowPence,
	})
}
