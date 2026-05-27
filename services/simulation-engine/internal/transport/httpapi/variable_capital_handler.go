package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	val "github.com/theding0x/capital-simulator/services/simulation-engine/internal/valorisation"
)

// Vol. II Ch. 16 — The Turnover of Variable Capital.

// --- request types ---

type createVariableCapitalAdvanceRequest struct {
	IndustrialCapitalID         string `json:"industrial_capital_id"`
	Pence                       int64  `json:"pence"`
	TurnoverPeriodMinutes       int64  `json:"turnover_period_minutes"`
	TurnoversPerYearBasisPoints int64  `json:"turnovers_per_year_basis_points"`
}

type recordPerCircuitRateRequest struct {
	SurplusPence  int64 `json:"surplus_pence"`
	VariablePence int64 `json:"variable_pence"`
}

type createAnnualSurplusRateContrastRequest struct {
	CapitalAID string `json:"capital_a_id"`
	CapitalBID string `json:"capital_b_id"`
	Period     string `json:"period"`
}

// --- response types ---

type variableCapitalAdvanceResponse struct {
	ID                          string `json:"id"`
	IndustrialCapitalID         string `json:"industrial_capital_id"`
	Pence                       int64  `json:"pence"`
	TurnoverPeriodMinutes       int64  `json:"turnover_period_minutes"`
	TurnoversPerYearBasisPoints int64  `json:"turnovers_per_year_basis_points"`
}

type perCircuitSurplusRateResponse struct {
	ID                       string `json:"id"`
	VariableCapitalAdvanceID string `json:"variable_capital_advance_id"`
	SurplusPence             int64  `json:"surplus_pence"`
	VariablePence            int64  `json:"variable_pence"`
	BasisPoints              int64  `json:"basis_points"`
}

type annualSurplusRateResponse struct {
	ID                       string `json:"id"`
	VariableCapitalAdvanceID string `json:"variable_capital_advance_id"`
	Period                   string `json:"period"`
	NBasisPoints             int64  `json:"n_basis_points"`
	PerCircuitBasisPoints    int64  `json:"per_circuit_basis_points"`
	AnnualBasisPoints        int64  `json:"annual_basis_points"`
}

type annualSurplusRateContrastResponse struct {
	ID                        string `json:"id"`
	CapitalAID                string `json:"capital_a_id"`
	CapitalBID                string `json:"capital_b_id"`
	CapitalAAnnualBasisPoints int64  `json:"capital_a_annual_basis_points"`
	CapitalBAnnualBasisPoints int64  `json:"capital_b_annual_basis_points"`
	RatioBasisPoints          int64  `json:"ratio_basis_points"`
}

type variableCapitalReproductionResponse struct {
	VariableCapitalAdvanceID     string `json:"variable_capital_advance_id"`
	YearReferenceMinutes         int64  `json:"year_reference_minutes"`
	TotalVariablePaidAnnualPence int64  `json:"total_variable_paid_annual_pence"`
	TotalSurplusAnnualPence      int64  `json:"total_surplus_annual_pence"`
}

// --- helpers -----------------------------------------------------------------

func advanceToResponse(a val.VariableCapitalAdvance) variableCapitalAdvanceResponse {
	return variableCapitalAdvanceResponse{
		ID:                          string(a.ID),
		IndustrialCapitalID:         a.IndustrialCapitalID,
		Pence:                       int64(a.Pence),
		TurnoverPeriodMinutes:       a.TurnoverPeriodMinutes,
		TurnoversPerYearBasisPoints: a.TurnoversPerYearBasisPoints,
	}
}

func perCircuitRateToResponse(r val.PerCircuitSurplusRate) perCircuitSurplusRateResponse {
	return perCircuitSurplusRateResponse{
		ID:                       string(r.ID),
		VariableCapitalAdvanceID: string(r.VariableCapitalAdvanceID),
		SurplusPence:             int64(r.SurplusPence),
		VariablePence:            int64(r.VariablePence),
		BasisPoints:              r.BasisPoints,
	}
}

func annualRateToResponse(r val.AnnualSurplusRate) annualSurplusRateResponse {
	return annualSurplusRateResponse{
		ID:                       string(r.ID),
		VariableCapitalAdvanceID: string(r.VariableCapitalAdvanceID),
		Period:                   r.Period,
		NBasisPoints:             r.NBasisPoints,
		PerCircuitBasisPoints:    r.PerCircuitBasisPoints,
		AnnualBasisPoints:        r.AnnualBasisPoints,
	}
}

// --- handlers ----------------------------------------------------------------

// CreateVariableCapitalAdvance handles POST /v1/variable-capital-advances.
func (h *Handler) CreateVariableCapitalAdvance(w http.ResponseWriter, r *http.Request) {
	var req createVariableCapitalAdvanceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.IndustrialCapitalID == "" {
		writeError(w, http.StatusBadRequest, "industrial_capital_id is required")
		return
	}
	if req.Pence <= 0 {
		writeError(w, http.StatusBadRequest, "pence must be positive")
		return
	}
	if req.TurnoverPeriodMinutes <= 0 {
		writeError(w, http.StatusBadRequest, "turnover_period_minutes must be positive")
		return
	}
	adv := val.VariableCapitalAdvance{
		ID:                          val.NewVariableCapitalAdvanceID(),
		IndustrialCapitalID:         req.IndustrialCapitalID,
		Pence:                       val.Pence(req.Pence),
		TurnoverPeriodMinutes:       req.TurnoverPeriodMinutes,
		TurnoversPerYearBasisPoints: req.TurnoversPerYearBasisPoints,
	}
	created, err := h.AnnualSurplusRates.CreateAdvance(r.Context(), adv)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "advance already exists")
			return
		}
		h.writeServerError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/variable-capital-advances/"+string(created.ID))
	writeJSON(w, http.StatusCreated, advanceToResponse(created))
}

// GetVariableCapitalAdvance handles GET /v1/variable-capital-advances/{id}.
func (h *Handler) GetVariableCapitalAdvance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	adv, err := h.AnnualSurplusRates.GetAdvance(r.Context(), val.VariableCapitalAdvanceID(id))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "advance not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, advanceToResponse(adv))
}

// ListVariableCapitalAdvances handles GET /v1/variable-capital-advances.
func (h *Handler) ListVariableCapitalAdvances(w http.ResponseWriter, r *http.Request) {
	icID := r.URL.Query().Get("industrial_capital_id")
	advances, err := h.AnnualSurplusRates.ListAdvances(r.Context(), icID)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	resp := make([]variableCapitalAdvanceResponse, len(advances))
	for i, a := range advances {
		resp[i] = advanceToResponse(a)
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": resp})
}

// RecordPerCircuitSurplusRate handles POST /v1/variable-capital-advances/{id}/per-circuit-rate.
func (h *Handler) RecordPerCircuitSurplusRate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req recordPerCircuitRateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	advanceID := val.VariableCapitalAdvanceID(id)
	rate, err := val.ComputePerCircuitRate(advanceID, val.Pence(req.SurplusPence), val.Pence(req.VariablePence))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.AnnualSurplusRates.RecordPerCircuitRate(r.Context(), rate)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "advance not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, perCircuitRateToResponse(saved))
}

// ComputeAnnualSurplusRateHandler handles POST /v1/variable-capital-advances/{id}/compute-annual-rate.
func (h *Handler) ComputeAnnualSurplusRateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	period := r.URL.Query().Get("period")
	if period == "" {
		writeError(w, http.StatusBadRequest, "period query parameter is required")
		return
	}
	advanceID := val.VariableCapitalAdvanceID(id)

	adv, err := h.AnnualSurplusRates.GetAdvance(r.Context(), advanceID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "advance not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	pcRate, err := h.AnnualSurplusRates.GetPerCircuitRate(r.Context(), advanceID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusConflict, "per-circuit rate not recorded yet; POST per-circuit-rate first")
			return
		}
		h.writeServerError(w, err)
		return
	}
	annualRate := val.ComputeAnnualSurplusRate(
		advanceID,
		adv.TurnoversPerYearBasisPoints,
		pcRate.BasisPoints,
		period,
	)
	saved, err := h.AnnualSurplusRates.RecordAnnualRate(r.Context(), annualRate)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, annualRateToResponse(saved))
}

// ListAnnualSurplusRates handles GET /v1/annual-surplus-rates.
func (h *Handler) ListAnnualSurplusRates(w http.ResponseWriter, r *http.Request) {
	icID := r.URL.Query().Get("industrial_capital_id")
	period := r.URL.Query().Get("period")
	rates, err := h.AnnualSurplusRates.ListAnnualRates(r.Context(), icID, period)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	resp := make([]annualSurplusRateResponse, len(rates))
	for i, r2 := range rates {
		resp[i] = annualRateToResponse(r2)
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": resp})
}

// CreateAnnualSurplusRateContrast handles POST /v1/annual-surplus-rate-contrasts.
func (h *Handler) CreateAnnualSurplusRateContrast(w http.ResponseWriter, r *http.Request) {
	var req createAnnualSurplusRateContrastRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.CapitalAID == "" || req.CapitalBID == "" {
		writeError(w, http.StatusBadRequest, "capital_a_id and capital_b_id are required")
		return
	}
	if req.Period == "" {
		writeError(w, http.StatusBadRequest, "period is required")
		return
	}
	aID := val.VariableCapitalAdvanceID(req.CapitalAID)
	bID := val.VariableCapitalAdvanceID(req.CapitalBID)

	// Fetch annual rates for the requested period to get the basis points.
	allRates, err := h.AnnualSurplusRates.ListAnnualRates(r.Context(), "", req.Period)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	var aAnnualBP, bAnnualBP int64
	for _, rt := range allRates {
		if rt.VariableCapitalAdvanceID == aID {
			aAnnualBP = rt.AnnualBasisPoints
		}
		if rt.VariableCapitalAdvanceID == bID {
			bAnnualBP = rt.AnnualBasisPoints
		}
	}
	contrast := val.ComputeAnnualSurplusRateContrast(aID, bID, aAnnualBP, bAnnualBP)
	saved, err := h.AnnualSurplusRates.RecordContrast(r.Context(), contrast)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, annualSurplusRateContrastResponse{
		ID:                        string(saved.ID),
		CapitalAID:                string(saved.CapitalAID),
		CapitalBID:                string(saved.CapitalBID),
		CapitalAAnnualBasisPoints: saved.CapitalAAnnualBasisPoints,
		CapitalBAnnualBasisPoints: saved.CapitalBAnnualBasisPoints,
		RatioBasisPoints:          saved.RatioBasisPoints,
	})
}

// GetVariableCapitalReproduction handles GET /v1/variable-capital-reproductions/{id}.
func (h *Handler) GetVariableCapitalReproduction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	yearMinutes := int64(525600) // standard year in minutes
	if ym := r.URL.Query().Get("year_minutes"); ym != "" {
		parsed, err := strconv.ParseInt(ym, 10, 64)
		if err != nil || parsed <= 0 {
			writeError(w, http.StatusBadRequest, "year_minutes must be a positive integer")
			return
		}
		yearMinutes = parsed
	}

	advanceID := val.VariableCapitalAdvanceID(id)
	adv, err := h.AnnualSurplusRates.GetAdvance(r.Context(), advanceID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusNotFound, "advance not found")
			return
		}
		h.writeServerError(w, err)
		return
	}
	pcRate, err := h.AnnualSurplusRates.GetPerCircuitRate(r.Context(), advanceID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeError(w, http.StatusConflict, "per-circuit rate not recorded yet")
			return
		}
		h.writeServerError(w, err)
		return
	}
	repro := val.ComputeVariableCapitalReproduction(adv, pcRate.BasisPoints, yearMinutes)
	writeJSON(w, http.StatusOK, variableCapitalReproductionResponse{
		VariableCapitalAdvanceID:     string(repro.VariableCapitalAdvanceID),
		YearReferenceMinutes:         repro.YearReferenceMinutes,
		TotalVariablePaidAnnualPence: int64(repro.TotalVariablePaidAnnualPence),
		TotalSurplusAnnualPence:      int64(repro.TotalSurplusAnnualPence),
	})
}
