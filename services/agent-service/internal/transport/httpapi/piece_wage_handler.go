package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

type createPieceWageRequest struct {
	WageFormID   string `json:"wage_form_id,omitempty"`
	PricePence   int64  `json:"price_pence,omitempty"`
	NormalOutput int64  `json:"normal_output"`
}

type pieceWageResponse struct {
	agent.PieceWage
}

type computePiecePriceRequest struct {
	DailyWageFarthings       int64  `json:"daily_wage_farthings"`
	DayValueProductFarthings int64  `json:"day_value_product_farthings"`
	NormalOutput             int64  `json:"normal_output"`
	PiecesProduced           int64  `json:"pieces_produced"`
	QualityOutcome           string `json:"quality_outcome"`
}

type computePiecePriceResponse struct {
	PiecePrice     int64 `json:"piece_price"`
	PieceValue     int64 `json:"piece_value"`
	ActualEarnings int64 `json:"actual_earnings"`
}

type createSubContractRequest struct {
	HeadLabourerID     string   `json:"head_labourer_id"`
	AssistantIDs       []string `json:"assistant_ids"`
	PieceRatePence     int64    `json:"piece_rate_pence"`
	AssistantRatePence int64    `json:"assistant_rate_pence"`
}

type subContractResponse struct {
	agent.SubContract
	Spread int64 `json:"spread"`
}

func buildSubContractResponse(sc agent.SubContract) subContractResponse {
	return subContractResponse{
		SubContract: sc,
		Spread:      agent.SubContractSpread(sc),
	}
}

// CreatePieceWage handles POST /v1/agents/{id}/piece-wages.
func (h *Handler) CreatePieceWage(w http.ResponseWriter, r *http.Request) {
	agentID := agent.AgentID(r.PathValue("id"))
	var req createPieceWageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	pw := agent.PieceWage{
		AgentID:      agentID,
		WageFormID:   agent.WageFormID(req.WageFormID),
		PricePence:   req.PricePence,
		NormalOutput: req.NormalOutput,
	}
	if !pw.WageFormID.IsZero() {
		wf, err := h.WageFormStore.GetWageForm(r.Context(), agentID)
		if err != nil {
			writeAppError(w, err)
			return
		}
		pw.PricePence = int64(wf.Wage.DailyPence) / req.NormalOutput
	}
	if err := pw.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.PieceWageStore.CreatePieceWage(r.Context(), pw)
	if err != nil {
		writeAppError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/agents/"+string(agentID)+"/piece-wages")
	writeJSON(w, http.StatusCreated, pieceWageResponse{saved})
}

// GetPieceWage handles GET /v1/agents/{id}/piece-wages.
func (h *Handler) GetPieceWage(w http.ResponseWriter, r *http.Request) {
	agentID := agent.AgentID(r.PathValue("id"))
	pw, err := h.PieceWageStore.GetPieceWage(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, pieceWageResponse{pw})
}

// ComputePiecePrice handles POST /v1/piece-price — stateless computation.
func (h *Handler) ComputePiecePrice(w http.ResponseWriter, r *http.Request) {
	var req computePiecePriceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.DailyWageFarthings <= 0 {
		writeError(w, http.StatusBadRequest, "daily_wage_farthings must be positive")
		return
	}
	if req.DayValueProductFarthings <= 0 {
		writeError(w, http.StatusBadRequest, "day_value_product_farthings must be positive")
		return
	}
	if req.NormalOutput <= 0 {
		writeError(w, http.StatusBadRequest, "normal_output must be positive")
		return
	}
	piecePrice := agent.ComputePiecePrice(req.DailyWageFarthings, req.NormalOutput)
	pieceValue := agent.ComputePieceValue(req.DayValueProductFarthings, req.NormalOutput)

	var actualEarnings int64
	if req.PiecesProduced > 0 && req.QualityOutcome != "" {
		pw := agent.PieceWage{PricePence: piecePrice, NormalOutput: req.NormalOutput}
		s := agent.PieceSession{
			PiecesProduced: req.PiecesProduced,
			QualityOutcome: agent.QualityOutcome(req.QualityOutcome),
		}
		actualEarnings = agent.ComputePieceEarnings(s, pw)
	}

	writeJSON(w, http.StatusOK, computePiecePriceResponse{
		PiecePrice:     piecePrice,
		PieceValue:     pieceValue,
		ActualEarnings: actualEarnings,
	})
}

// CreateSubContract handles POST /v1/sub-contracts.
func (h *Handler) CreateSubContract(w http.ResponseWriter, r *http.Request) {
	var req createSubContractRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	assistantIDs := make([]agent.AgentID, len(req.AssistantIDs))
	for i, id := range req.AssistantIDs {
		assistantIDs[i] = agent.AgentID(id)
	}
	sc := agent.SubContract{
		HeadLabourerID:     agent.AgentID(req.HeadLabourerID),
		AssistantIDs:       assistantIDs,
		PieceRatePence:     req.PieceRatePence,
		AssistantRatePence: req.AssistantRatePence,
	}
	if err := sc.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.PieceWageStore.CreateSubContract(r.Context(), sc)
	if err != nil {
		writeAppError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/sub-contracts/"+string(saved.ID))
	writeJSON(w, http.StatusCreated, buildSubContractResponse(saved))
}

// GetSubContract handles GET /v1/sub-contracts/{id}.
func (h *Handler) GetSubContract(w http.ResponseWriter, r *http.Request) {
	id := agent.SubContractID(r.PathValue("id"))
	sc, err := h.PieceWageStore.GetSubContract(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildSubContractResponse(sc))
}
