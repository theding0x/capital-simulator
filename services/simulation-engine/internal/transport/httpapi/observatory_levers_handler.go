package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// leversRequest is the POST /v1/observatory/levers body. Each field is optional;
// omitted fields leave that lever unchanged.
type leversRequest struct {
	SurplusRateBaseBP  *int64 `json:"surplus_rate_base_bp"`
	BaseWagePence      *int64 `json:"base_wage_pence"`
	AccumulationRateBP *int64 `json:"accumulation_rate_bp"`
}

// leversResponse echoes the applied (clamped) lever values.
type leversResponse struct {
	SurplusRateBaseBP  int64 `json:"surplus_rate_base_bp"`
	BaseWagePence      int64 `json:"base_wage_pence"`
	AccumulationRateBP int64 `json:"accumulation_rate_bp"`
}

// SetObservatoryLevers handles POST /v1/observatory/levers (Slice 3 — the
// levers). It applies a partial perturbation of the abode's law parameters —
// the working day (the rate of surplus-value), the wage (the value of
// labour-power), and the accumulation rate α — to the live state; the effects
// appear in subsequent snapshots as the General Law responds. Returns the
// applied (clamped) lever values.
func (h *Handler) SetObservatoryLevers(w http.ResponseWriter, r *http.Request) {
	if h.AbodeStates == nil {
		h.writeServerError(w, errors.New("abode state store not configured"))
		return
	}
	var body leversRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "malformed JSON body")
		return
	}
	state, err := h.AbodeStates.SetAbodeLevers(r.Context(), simulation.LeverUpdate{
		SurplusRateBaseBP:  body.SurplusRateBaseBP,
		BaseWagePence:      body.BaseWagePence,
		AccumulationRateBP: body.AccumulationRateBP,
	})
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, leversResponse{
		SurplusRateBaseBP:  state.SurplusRateBaseBP,
		BaseWagePence:      state.BaseWagePence,
		AccumulationRateBP: state.AccumulationRateBP,
	})
}
