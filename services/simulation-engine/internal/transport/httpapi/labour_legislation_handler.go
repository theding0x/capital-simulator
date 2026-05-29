package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Ch. 28 — Bloody Legislation Against the Expropriated.

// --- WageStatute DTOs ---

type wageStatuteResponse struct {
	ID                 string `json:"id"`
	HistoricalStageID  string `json:"historical_stage_id"`
	Period             string `json:"period"`
	Jurisdiction       string `json:"jurisdiction"`
	MaxWagePence       int64  `json:"max_wage_pence"`
	MinWagePence       int64  `json:"min_wage_pence"`
	EnforcementPenalty string `json:"enforcement_penalty"`
	CreatedAt          string `json:"created_at"`
}

type createWageStatuteRequest struct {
	Period             string `json:"period"`
	Jurisdiction       string `json:"jurisdiction"`
	MaxWagePence       int64  `json:"max_wage_pence"`
	MinWagePence       int64  `json:"min_wage_pence"`
	EnforcementPenalty string `json:"enforcement_penalty"`
}

func toWageStatuteResponse(w simulation.WageStatute) wageStatuteResponse {
	return wageStatuteResponse{
		ID:                 string(w.ID),
		HistoricalStageID:  string(w.HistoricalStageID),
		Period:             w.Period,
		Jurisdiction:       w.Jurisdiction,
		MaxWagePence:       w.MaxWagePence,
		MinWagePence:       w.MinWagePence,
		EnforcementPenalty: w.EnforcementPenalty,
		CreatedAt:          w.CreatedAt.Format(time.RFC3339Nano),
	}
}

// --- VagrancyLaw DTOs ---

type vagrancyLawResponse struct {
	ID                string `json:"id"`
	HistoricalStageID string `json:"historical_stage_id"`
	Period            string `json:"period"`
	Jurisdiction      string `json:"jurisdiction"`
	Punishment        string `json:"punishment"`
	TargetPopulation  string `json:"target_population"`
	CreatedAt         string `json:"created_at"`
}

type createVagrancyLawRequest struct {
	Period           string `json:"period"`
	Jurisdiction     string `json:"jurisdiction"`
	Punishment       string `json:"punishment"`
	TargetPopulation string `json:"target_population"`
}

func toVagrancyLawResponse(v simulation.VagrancyLaw) vagrancyLawResponse {
	return vagrancyLawResponse{
		ID:                string(v.ID),
		HistoricalStageID: string(v.HistoricalStageID),
		Period:            v.Period,
		Jurisdiction:      v.Jurisdiction,
		Punishment:        v.Punishment,
		TargetPopulation:  v.TargetPopulation,
		CreatedAt:         v.CreatedAt.Format(time.RFC3339Nano),
	}
}

// --- LabourDisciplineRegime DTO ---

type labourDisciplineResponse struct {
	Period       string                `json:"period"`
	Mechanisms   []string              `json:"mechanisms"`
	WageStatutes []wageStatuteResponse `json:"wage_statutes"`
	VagrancyLaws []vagrancyLawResponse `json:"vagrancy_laws"`
}

// --- StatutoryWage DTO ---

type compareStatutoryWageRequest struct {
	ActedWagePence  int64 `json:"acted_wage_pence"`
	MarketWagePence int64 `json:"market_wage_pence"`
}

// CreateWageStatute handles POST /v1/historical-stages/{id}/wage-statutes.
func (h *Handler) CreateWageStatute(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))

	_, err := h.HistoricalStages.GetHistoricalStage(r.Context(), stageID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "historical stage not found")
		return
	}
	if err != nil {
		h.writeServerError(w, err)
		return
	}

	var req createWageStatuteRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	ws := simulation.WageStatute{
		HistoricalStageID:  stageID,
		Period:             req.Period,
		Jurisdiction:       req.Jurisdiction,
		MaxWagePence:       req.MaxWagePence,
		MinWagePence:       req.MinWagePence,
		EnforcementPenalty: req.EnforcementPenalty,
	}
	if err := ws.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.WageStatutes.CreateWageStatute(r.Context(), ws)
	if err != nil {
		h.writeServerError(w, err)
		return
	}

	w.Header().Set("Location", "/v1/historical-stages/"+string(stageID)+"/wage-statutes/"+string(created.ID))
	writeJSON(w, http.StatusCreated, toWageStatuteResponse(created))
}

// CreateVagrancyLaw handles POST /v1/historical-stages/{id}/vagrancy-laws.
func (h *Handler) CreateVagrancyLaw(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))

	_, err := h.HistoricalStages.GetHistoricalStage(r.Context(), stageID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "historical stage not found")
		return
	}
	if err != nil {
		h.writeServerError(w, err)
		return
	}

	var req createVagrancyLawRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	vl := simulation.VagrancyLaw{
		HistoricalStageID: stageID,
		Period:            req.Period,
		Jurisdiction:      req.Jurisdiction,
		Punishment:        req.Punishment,
		TargetPopulation:  req.TargetPopulation,
	}
	if err := vl.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.VagrancyLaws.CreateVagrancyLaw(r.Context(), vl)
	if err != nil {
		h.writeServerError(w, err)
		return
	}

	w.Header().Set("Location", "/v1/historical-stages/"+string(stageID)+"/vagrancy-laws/"+string(created.ID))
	writeJSON(w, http.StatusCreated, toVagrancyLawResponse(created))
}

// GetLabourDiscipline handles GET /v1/historical-stages/{id}/labour-discipline.
// Returns the assembled LabourDisciplineRegime for the stage.
func (h *Handler) GetLabourDiscipline(w http.ResponseWriter, r *http.Request) {
	stageID := simulation.HistoricalStageID(r.PathValue("id"))

	stage, err := h.HistoricalStages.GetHistoricalStage(r.Context(), stageID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "historical stage not found")
		return
	}
	if err != nil {
		h.writeServerError(w, err)
		return
	}

	ws, err := h.WageStatutes.ListWageStatutesByStage(r.Context(), stageID)
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	vl, err := h.VagrancyLaws.ListVagrancyLawsByStage(r.Context(), stageID)
	if err != nil {
		h.writeServerError(w, err)
		return
	}

	regime := simulation.AssembleRegime(stage.Name, ws, vl)

	wsResp := make([]wageStatuteResponse, 0, len(regime.WageStatutes))
	for _, item := range regime.WageStatutes {
		wsResp = append(wsResp, toWageStatuteResponse(item))
	}
	vlResp := make([]vagrancyLawResponse, 0, len(regime.VagrancyLaws))
	for _, item := range regime.VagrancyLaws {
		vlResp = append(vlResp, toVagrancyLawResponse(item))
	}

	writeJSON(w, http.StatusOK, labourDisciplineResponse{
		Period:       regime.Period,
		Mechanisms:   regime.Mechanisms,
		WageStatutes: wsResp,
		VagrancyLaws: vlResp,
	})
}

// CompareStatutoryWage handles POST /v1/statutory-wages/compare.
// Stateless: computes deviation between an acted (legally imposed) wage and
// the prevailing market wage.
func (h *Handler) CompareStatutoryWage(w http.ResponseWriter, r *http.Request) {
	var req compareStatutoryWageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	sw, err := simulation.ComputeStatutoryWage(req.ActedWagePence, req.MarketWagePence)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, sw)
}