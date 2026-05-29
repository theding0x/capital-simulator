package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// costPriceRequest is the body for POST /v1/profit/cost-price. All magnitudes
// are LabourMinutes. FixedWearAndTear and FixedAdvanced are optional; when both
// are zero the constant capital is treated as wholly circulating.
type costPriceRequest struct {
	Constant         int64 `json:"constant"`
	Variable         int64 `json:"variable"`
	FixedWearAndTear int64 `json:"fixed_wear_and_tear"`
	FixedAdvanced    int64 `json:"fixed_advanced"`
}

func (req costPriceRequest) outlay() profit.CapitalOutlay {
	return profit.CapitalOutlay{
		Constant:         profit.ConstantCapital(req.Constant),
		Variable:         profit.VariableCapital(req.Variable),
		FixedWearAndTear: profit.ConstantCapital(req.FixedWearAndTear),
		FixedAdvanced:    profit.ConstantCapital(req.FixedAdvanced),
	}
}

// CreateCostPrice handles POST /v1/profit/cost-price — compute k = c + v from a
// capital outlay and persist it.
func (h *Handler) CreateCostPrice(w http.ResponseWriter, r *http.Request) {
	var req costPriceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	outlay := req.outlay()
	if err := outlay.Validate(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	cp, err := h.Store.CreateCostPrice(r.Context(), profit.ComputeCostPrice(outlay))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	w.Header().Set("Location", "/v1/profit/cost-price/"+string(cp.ID))
	writeJSON(w, http.StatusCreated, cp)
}

// GetCostPrice handles GET /v1/profit/cost-price/{id}.
func (h *Handler) GetCostPrice(w http.ResponseWriter, r *http.Request) {
	id := profit.CostPriceID(r.PathValue("id"))
	cp, err := h.Store.GetCostPrice(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cp)
}

// ListCostPrices handles GET /v1/profit/cost-price.
func (h *Handler) ListCostPrices(w http.ResponseWriter, r *http.Request) {
	items, err := h.Store.ListCostPrices(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []profit.CostPrice{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// profitFormRequest is the body for POST /v1/profit/profit-form.
type profitFormRequest struct {
	CostPrice    int64 `json:"cost_price"`
	SurplusValue int64 `json:"surplus_value"`
}

// profitFormResponse carries the C = k + p transformation plus the selling-price
// scenarios that demonstrate a profit can be realised even below value.
type profitFormResponse struct {
	profit.ProfitForm
	SellingPriceScenarios []profit.Profit `json:"selling_price_scenarios"`
}

// ComputeProfitForm handles POST /v1/profit/profit-form — regroup C = k + s as
// C = k + p, and compute the spread of profitable selling prices between the
// cost-price and the commodity's value.
func (h *Handler) ComputeProfitForm(w http.ResponseWriter, r *http.Request) {
	var req profitFormRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.CostPrice < 0 || req.SurplusValue < 0 {
		writeError(w, http.StatusUnprocessableEntity, "cost_price and surplus_value must be non-negative")
		return
	}
	k := profit.LabourMinutes(req.CostPrice)
	s := profit.SurplusValue(req.SurplusValue)
	form := profit.ComputeProfitForm(k, s)
	writeJSON(w, http.StatusOK, profitFormResponse{
		ProfitForm:            form,
		SellingPriceScenarios: sellingPriceScenarios(k, form.CommodityValue),
	})
}

// sellingPriceScenarios samples selling prices across the interval (k, C],
// illustrating Marx's point that every price above the cost-price yields a
// profit — including those below the commodity's value. The final sample, at C
// itself, realises the full surplus-value.
func sellingPriceScenarios(k, commodityValue profit.LabourMinutes) []profit.Profit {
	out := []profit.Profit{}
	if commodityValue <= k {
		return out
	}
	span := commodityValue - k
	const samples = 5
	for i := 1; i <= samples; i++ {
		sell := k + span*profit.LabourMinutes(i)/samples
		out = append(out, profit.ComputeProfit(sell, k, commodityValue))
	}
	return out
}
