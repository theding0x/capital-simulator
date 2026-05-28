package httpapi

import (
	"errors"
	"net/http"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Vol. II Ch. 18 — The Role of Money-Capital in Reproduction.

// --- request types ---

type createMoneySupplyApportionmentRequest struct {
	TotalSocialMoneyPence int64  `json:"total_social_money_pence"`
	DeptIReservePence     int64  `json:"dept_i_reserve_pence"`
	DeptIIReservePence    int64  `json:"dept_ii_reserve_pence"`
	WageRotationPence     int64  `json:"wage_rotation_pence"`
	IdleHoardPence        int64  `json:"idle_hoard_pence"`
	Period                string `json:"period"`
}

type createDepartmentMoneyReserveRequest struct {
	Department     string `json:"department"`
	ReservePurpose string `json:"reserve_purpose"`
	ReservePence   int64  `json:"reserve_pence"`
	Period         string `json:"period"`
}

type createCirculatingMoneyMassRequest struct {
	MoneyStockPence            int64  `json:"money_stock_pence"`
	VelocityPerYearBasisPoints int64  `json:"velocity_per_year_basis_points"`
	Period                     string `json:"period"`
}

type createWageRotationFundRequest struct {
	FundPence          int64  `json:"fund_pence"`
	WageCycleFrequency int64  `json:"wage_cycle_frequency"`
	Department         string `json:"department"`
	Period             string `json:"period"`
}

type createInterDepartmentSettlementRequest struct {
	FromDepartment    string `json:"from_department"`
	ToDepartment      string `json:"to_department"`
	SettledPence      int64  `json:"settled_pence"`
	SettlementPurpose string `json:"settlement_purpose"`
	Period            string `json:"period"`
}

// --- response types ---

type moneySupplyApportionmentResponse struct {
	ID                    string `json:"id"`
	TotalSocialMoneyPence int64  `json:"total_social_money_pence"`
	DeptIReservePence     int64  `json:"dept_i_reserve_pence"`
	DeptIIReservePence    int64  `json:"dept_ii_reserve_pence"`
	WageRotationPence     int64  `json:"wage_rotation_pence"`
	IdleHoardPence        int64  `json:"idle_hoard_pence"`
	Period                string `json:"period"`
}

type departmentMoneyReserveResponse struct {
	ID             string `json:"id"`
	Department     string `json:"department"`
	ReservePurpose string `json:"reserve_purpose"`
	ReservePence   int64  `json:"reserve_pence"`
	Period         string `json:"period"`
}

type circulatingMoneyMassResponse struct {
	ID                             string `json:"id"`
	MoneyStockPence                int64  `json:"money_stock_pence"`
	VelocityPerYearBasisPoints     int64  `json:"velocity_per_year_basis_points"`
	EffectiveCirculatingValuePence int64  `json:"effective_circulating_value_pence"`
	Period                         string `json:"period"`
}

type wageRotationFundResponse struct {
	ID                 string `json:"id"`
	FundPence          int64  `json:"fund_pence"`
	WageCycleFrequency int64  `json:"wage_cycle_frequency"`
	Department         string `json:"department"`
	Period             string `json:"period"`
}

type interDepartmentSettlementResponse struct {
	ID                string `json:"id"`
	FromDepartment    string `json:"from_department"`
	ToDepartment      string `json:"to_department"`
	SettledPence      int64  `json:"settled_pence"`
	SettlementPurpose string `json:"settlement_purpose"`
	Period            string `json:"period"`
}

type apportionmentBalanceCheckResponse struct {
	Balanced        bool   `json:"balanced"`
	ApportionmentID string `json:"apportionment_id"`
}

// --- helpers ---

func apportionmentToResponse(a repro.MoneySupplyApportionment) moneySupplyApportionmentResponse {
	return moneySupplyApportionmentResponse{
		ID:                    string(a.ID),
		TotalSocialMoneyPence: a.TotalSocialMoneyPence,
		DeptIReservePence:     a.DeptIReservePence,
		DeptIIReservePence:    a.DeptIIReservePence,
		WageRotationPence:     a.WageRotationPence,
		IdleHoardPence:        a.IdleHoardPence,
		Period:                a.Period,
	}
}

func reserveToResponse(r repro.DepartmentMoneyReserve) departmentMoneyReserveResponse {
	return departmentMoneyReserveResponse{
		ID:             string(r.ID),
		Department:     string(r.Department),
		ReservePurpose: string(r.ReservePurpose),
		ReservePence:   r.ReservePence,
		Period:         r.Period,
	}
}

func massToResponse(mass repro.CirculatingMoneyMass) circulatingMoneyMassResponse {
	return circulatingMoneyMassResponse{
		ID:                             string(mass.ID),
		MoneyStockPence:                mass.MoneyStockPence,
		VelocityPerYearBasisPoints:     mass.VelocityPerYearBasisPoints,
		EffectiveCirculatingValuePence: mass.EffectiveCirculatingValuePence,
		Period:                         mass.Period,
	}
}

func wageToResponse(w repro.WageRotationFund) wageRotationFundResponse {
	return wageRotationFundResponse{
		ID:                 string(w.ID),
		FundPence:          w.FundPence,
		WageCycleFrequency: w.WageCycleFrequency,
		Department:         string(w.Department),
		Period:             w.Period,
	}
}

func settlementToResponse(s repro.InterDepartmentSettlement) interDepartmentSettlementResponse {
	return interDepartmentSettlementResponse{
		ID:                string(s.ID),
		FromDepartment:    string(s.FromDepartment),
		ToDepartment:      string(s.ToDepartment),
		SettledPence:      s.SettledPence,
		SettlementPurpose: string(s.SettlementPurpose),
		Period:            s.Period,
	}
}

// --- handlers ---

// CreateApportionment handles POST /v1/reproduction/apportionments.
func (h *Handler) CreateApportionment(w http.ResponseWriter, r *http.Request) {
	var req createMoneySupplyApportionmentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	a := repro.MoneySupplyApportionment{
		ID:                    repro.NewMoneySupplyApportionmentID(),
		TotalSocialMoneyPence: req.TotalSocialMoneyPence,
		DeptIReservePence:     req.DeptIReservePence,
		DeptIIReservePence:    req.DeptIIReservePence,
		WageRotationPence:     req.WageRotationPence,
		IdleHoardPence:        req.IdleHoardPence,
		Period:                req.Period,
	}
	if err := a.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.MoneyCapital.CreateMoneySupplyApportionment(r.Context(), a)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "apportionment already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/reproduction/apportionments/"+string(created.ID))
	writeJSON(w, http.StatusCreated, apportionmentToResponse(created))
}

// CreateDepartmentReserve handles POST /v1/reproduction/department-reserves.
func (h *Handler) CreateDepartmentReserve(w http.ResponseWriter, r *http.Request) {
	var req createDepartmentMoneyReserveRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	reserve := repro.DepartmentMoneyReserve{
		ID:             repro.NewDepartmentMoneyReserveID(),
		Department:     repro.CapitalDepartment(req.Department),
		ReservePurpose: repro.ReservePurpose(req.ReservePurpose),
		ReservePence:   req.ReservePence,
		Period:         req.Period,
	}
	if err := reserve.Department.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := reserve.ReservePurpose.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.MoneyCapital.CreateDepartmentMoneyReserve(r.Context(), reserve)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "department reserve already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/reproduction/department-reserves/"+string(created.ID))
	writeJSON(w, http.StatusCreated, reserveToResponse(created))
}

// CreateCirculatingMoneyMass handles POST /v1/reproduction/circulating-money-masses.
// The effective_circulating_value_pence is computed server-side as stock * velocity / 10000.
func (h *Handler) CreateCirculatingMoneyMass(w http.ResponseWriter, r *http.Request) {
	var req createCirculatingMoneyMassRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	if req.MoneyStockPence < 0 {
		writeError(w, http.StatusBadRequest, repro.ErrNegativeMoneyStock.Error())
		return
	}
	if req.VelocityPerYearBasisPoints <= 0 {
		writeError(w, http.StatusBadRequest, repro.ErrZeroVelocity.Error())
		return
	}
	effective := repro.ComputeEffectiveCirculatingValue(req.MoneyStockPence, req.VelocityPerYearBasisPoints)
	mass := repro.CirculatingMoneyMass{
		ID:                             repro.NewCirculatingMoneyMassID(),
		MoneyStockPence:                req.MoneyStockPence,
		VelocityPerYearBasisPoints:     req.VelocityPerYearBasisPoints,
		EffectiveCirculatingValuePence: effective,
		Period:                         req.Period,
	}
	created, err := h.MoneyCapital.CreateCirculatingMoneyMass(r.Context(), mass)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "circulating money mass already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/reproduction/circulating-money-masses/"+string(created.ID))
	writeJSON(w, http.StatusCreated, massToResponse(created))
}

// CreateWageRotationFund handles POST /v1/reproduction/wage-rotation-funds.
func (h *Handler) CreateWageRotationFund(w http.ResponseWriter, r *http.Request) {
	var req createWageRotationFundRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	fund := repro.WageRotationFund{
		ID:                 repro.NewWageRotationFundID(),
		FundPence:          req.FundPence,
		WageCycleFrequency: req.WageCycleFrequency,
		Department:         repro.CapitalDepartment(req.Department),
		Period:             req.Period,
	}
	created, err := h.MoneyCapital.CreateWageRotationFund(r.Context(), fund)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "wage rotation fund already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/reproduction/wage-rotation-funds/"+string(created.ID))
	writeJSON(w, http.StatusCreated, wageToResponse(created))
}

// CreateInterDepartmentSettlement handles POST /v1/reproduction/inter-department-settlements.
func (h *Handler) CreateInterDepartmentSettlement(w http.ResponseWriter, r *http.Request) {
	var req createInterDepartmentSettlementRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	s := repro.InterDepartmentSettlement{
		ID:                repro.NewInterDepartmentSettlementID(),
		FromDepartment:    repro.CapitalDepartment(req.FromDepartment),
		ToDepartment:      repro.CapitalDepartment(req.ToDepartment),
		SettledPence:      req.SettledPence,
		SettlementPurpose: repro.ReservePurpose(req.SettlementPurpose),
		Period:            req.Period,
	}
	if err := s.FromDepartment.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.ToDepartment.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.SettlementPurpose.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.MoneyCapital.CreateInterDepartmentSettlement(r.Context(), s)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "inter-department settlement already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/reproduction/inter-department-settlements/"+string(created.ID))
	writeJSON(w, http.StatusCreated, settlementToResponse(created))
}

// ListApportionments handles GET /v1/reproduction/apportionments.
func (h *Handler) ListApportionments(w http.ResponseWriter, r *http.Request) {
	list, err := h.MoneyCapital.ListMoneySupplyApportionments(r.Context())
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	resp := make([]moneySupplyApportionmentResponse, len(list))
	for i, a := range list {
		resp[i] = apportionmentToResponse(a)
	}
	writeJSON(w, http.StatusOK, resp)
}

// ListInterDepartmentSettlements handles GET /v1/reproduction/inter-department-settlements.
func (h *Handler) ListInterDepartmentSettlements(w http.ResponseWriter, r *http.Request) {
	list, err := h.MoneyCapital.ListInterDepartmentSettlements(r.Context())
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	resp := make([]interDepartmentSettlementResponse, len(list))
	for i, s := range list {
		resp[i] = settlementToResponse(s)
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetApportionmentBalanceCheck handles GET /v1/reproduction/apportionments/{id}/balance-check.
func (h *Handler) GetApportionmentBalanceCheck(w http.ResponseWriter, r *http.Request) {
	id := repro.MoneySupplyApportionmentID(r.PathValue("id"))
	balanced, err := h.MoneyCapital.CheckApportionmentBalance(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "apportionment not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apportionmentBalanceCheckResponse{
		Balanced:        balanced,
		ApportionmentID: string(id),
	})
}
