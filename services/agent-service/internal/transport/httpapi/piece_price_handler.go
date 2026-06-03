package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

type repricePieceWageRequest struct {
	ProductivityFactor float64 `json:"productivity_factor"`
}

type repricePieceWageResponse struct {
	PieceWage pieceWageResponse       `json:"piece_wage"`
	Current   agent.PiecePricePoint   `json:"current"`
	Appended  bool                    `json:"appended"`
	Points    []agent.PiecePricePoint `json:"points"`
}

// RepricePieceWage handles POST /v1/piece-wages/{agentID}/reprice. It reprices
// the agent's piece-wage from its immutable baseline by the given productivity
// factor — Marx's Ch. 21 law: the daily wage is held constant, so the price per
// piece falls in inverse proportion as the socially-normal output rises — and
// appends one point to the price series. When the factor is unchanged from the
// last recorded point the series is left untouched (appended = false), which is
// how the simulation-engine tick scheduler avoids spamming flat points.
func (h *Handler) RepricePieceWage(w http.ResponseWriter, r *http.Request) {
	agentID := agent.AgentID(r.PathValue("agentID"))
	var req repricePieceWageRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.ProductivityFactor <= 0 {
		writeError(w, http.StatusBadRequest, "productivity_factor must be positive")
		return
	}
	baseline, err := h.PieceWageStore.GetPieceWage(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	repriced := agent.RepricePieceWage(baseline, req.ProductivityFactor)
	current, appended, err := h.PiecePricePoints.AppendPiecePricePoint(r.Context(), agent.PiecePricePoint{
		AgentID:            agentID,
		ProductivityFactor: req.ProductivityFactor,
		PricePence:         repriced.PricePence,
		NormalOutput:       repriced.NormalOutput,
	})
	if err != nil {
		writeAppError(w, err)
		return
	}
	points, err := h.PiecePricePoints.ListPiecePricePoints(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, repricePieceWageResponse{
		PieceWage: buildPieceWageResponse(baseline),
		Current:   current,
		Appended:  appended,
		Points:    points,
	})
}

// ListPiecePricePoints handles GET /v1/piece-wages/{agentID}/price-points. It
// returns the piece-price time series for one agent's piece-wage, oldest first.
// The slice is always non-nil so an empty series marshals as [] not null.
func (h *Handler) ListPiecePricePoints(w http.ResponseWriter, r *http.Request) {
	agentID := agent.AgentID(r.PathValue("agentID"))
	points, err := h.PiecePricePoints.ListPiecePricePoints(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, points)
}

// ListPieceWages handles GET /v1/piece-wages — every piece-wage contract. The
// simulation-engine tick scheduler enumerates these to reprice each one as
// factory productivity rises.
func (h *Handler) ListPieceWages(w http.ResponseWriter, r *http.Request) {
	pws, err := h.PieceWageStore.ListPieceWages(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	out := make([]pieceWageResponse, 0, len(pws))
	for _, pw := range pws {
		out = append(out, buildPieceWageResponse(pw))
	}
	writeJSON(w, http.StatusOK, out)
}
