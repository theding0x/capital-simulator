package httpapi

import (
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

// SetObservatoryLevers handles POST /v1/observatory/levers. It applies a partial
// perturbation of the abode's law parameters — the working day (s′), the wage
// (the value of labour-power), and the accumulation rate α — to the caller's
// in-memory session run (keyed by X-Atlas-Session). Returns the applied
// (clamped) values. No store I/O.
func (h *Handler) SetObservatoryLevers(w http.ResponseWriter, r *http.Request) {
	if h.Observatory == nil {
		h.writeServerError(w, errors.New("observatory not configured"))
		return
	}
	var body leversRequest
	if err := decodeJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	run := h.Observatory.GetOrCreate(r.Header.Get("X-Atlas-Session"))
	abode := run.ApplyLevers(simulation.LeverUpdate{
		SurplusRateBaseBP:  body.SurplusRateBaseBP,
		BaseWagePence:      body.BaseWagePence,
		AccumulationRateBP: body.AccumulationRateBP,
	})
	writeJSON(w, http.StatusOK, leversResponse{
		SurplusRateBaseBP:  abode.SurplusRateBaseBP,
		BaseWagePence:      abode.BaseWagePence,
		AccumulationRateBP: abode.AccumulationRateBP,
	})
}
