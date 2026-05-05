package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

// ── Ch. 8: Constant & Variable Capital ──────────────────────────────────────

type constantCapitalInput struct {
	OriginalValue   commodity.LabourMinutes `json:"original_value"`
	Kind            commodity.ConstantKind  `json:"kind"`
	ServiceLifeDays int64                   `json:"service_life_days"`
}

type variableCapitalInput struct {
	WageValue  commodity.LabourMinutes `json:"wage_value"`
	WorkingDay commodity.LabourMinutes `json:"working_day"`
}

type decomposeRequest struct {
	Constants []constantCapitalInput `json:"constant_capitals"`
	Variable  variableCapitalInput   `json:"variable_capital"`
}

type decomposeResponse struct {
	ProductValue       commodity.ProductValue       `json:"product_value"`
	CapitalComposition commodity.CapitalComposition `json:"capital_composition"`
	CompositionRatio   float64                      `json:"composition_ratio"`
}

// DecomposeCapital handles POST /v1/capital/decompose [Ch.8].
func (h *Handler) DecomposeCapital(w http.ResponseWriter, r *http.Request) {
	var req decomposeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Constants) == 0 {
		writeError(w, http.StatusBadRequest, "constant_capitals must not be empty")
		return
	}
	constants := make([]commodity.ConstantCapital, len(req.Constants))
	for i, c := range req.Constants {
		constants[i] = commodity.ConstantCapital{
			OriginalValue:   c.OriginalValue,
			Kind:            c.Kind,
			ServiceLifeDays: c.ServiceLifeDays,
		}
		if err := constants[i].Validate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	v := commodity.VariableCapital{
		WageValue:  req.Variable.WageValue,
		WorkingDay: req.Variable.WorkingDay,
	}
	if err := v.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	pv := commodity.DecomposeProductValue(constants, v)
	cc := commodity.CapitalComposition{Constant: pv.Constant, Variable: v.WageValue}
	writeJSON(w, http.StatusOK, decomposeResponse{
		ProductValue:       pv,
		CapitalComposition: cc,
		CompositionRatio:   cc.Ratio(),
	})
}

// CapitalCompositionHandler handles GET /v1/capital/composition?constant_value=&variable_value= [Ch.8].
func (h *Handler) CapitalCompositionHandler(w http.ResponseWriter, r *http.Request) {
	cStr := r.URL.Query().Get("constant_value")
	vStr := r.URL.Query().Get("variable_value")
	if cStr == "" || vStr == "" {
		writeError(w, http.StatusBadRequest, "constant_value and variable_value are required")
		return
	}
	c, err := strconv.ParseInt(cStr, 10, 64)
	if err != nil || c < 0 {
		writeError(w, http.StatusBadRequest, "constant_value must be a non-negative integer")
		return
	}
	v, err := strconv.ParseInt(vStr, 10, 64)
	if err != nil || v < 0 {
		writeError(w, http.StatusBadRequest, "variable_value must be a non-negative integer")
		return
	}
	cc := commodity.CapitalComposition{
		Constant: commodity.LabourMinutes(c),
		Variable: commodity.LabourMinutes(v),
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"constant":          cc.Constant,
		"variable":          cc.Variable,
		"composition_ratio": cc.Ratio(),
	})
}

// ── Ch. 9: Rate of Surplus-Value ─────────────────────────────────────────────

type createProductionAccountRequest struct {
	Constant commodity.LabourMinutes `json:"constant"`
	Variable commodity.LabourMinutes `json:"variable"`
	Surplus  commodity.SurplusValue  `json:"surplus"`
}

type productionAccountResponse struct {
	commodity.ProductionAccount
	RateOfSurplusValue float64                 `json:"rate_of_surplus_value"`
	ValueProduct       commodity.ValueProduct  `json:"value_product"`
	ExpandedCapital    commodity.LabourMinutes `json:"expanded_capital"`
	SurplusProduceFrac float64                 `json:"surplus_produce_fraction"`
}

func buildAccountResponse(a commodity.ProductionAccount) productionAccountResponse {
	rate, _ := a.Rate()
	return productionAccountResponse{
		ProductionAccount:  a,
		RateOfSurplusValue: float64(rate),
		ValueProduct:       a.ValueProduct(),
		ExpandedCapital:    a.ExpandedCapital(),
		SurplusProduceFrac: commodity.SurplusProduceRatio(
			commodity.LabourMinutes(a.Surplus),
			a.Variable,
		),
	}
}

// CreateProductionAccount handles POST /v1/production-accounts [Ch.9].
func (h *Handler) CreateProductionAccount(w http.ResponseWriter, r *http.Request) {
	var req createProductionAccountRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a := commodity.ProductionAccount{
		Constant:  req.Constant,
		Variable:  req.Variable,
		Surplus:   req.Surplus,
		CreatedAt: time.Now().UTC(),
	}
	if err := a.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if a.Variable == 0 {
		writeError(w, http.StatusBadRequest, commodity.ErrDivisionByZero.Error())
		return
	}
	saved, err := h.ProductionAccountStore.CreateProductionAccount(r.Context(), a)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, buildAccountResponse(saved))
}

// ListProductionAccounts handles GET /v1/production-accounts [Ch.9].
func (h *Handler) ListProductionAccounts(w http.ResponseWriter, r *http.Request) {
	all, err := h.ProductionAccountStore.ListProductionAccounts(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}
	if all == nil {
		all = []commodity.ProductionAccount{}
	}
	items := make([]productionAccountResponse, len(all))
	for i, a := range all {
		items[i] = buildAccountResponse(a)
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

// GetProductionAccount handles GET /v1/production-accounts/{id} [Ch.9].
func (h *Handler) GetProductionAccount(w http.ResponseWriter, r *http.Request) {
	id := commodity.ProductionAccountID(r.PathValue("id"))
	a, err := h.ProductionAccountStore.GetProductionAccount(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildAccountResponse(a))
}

type rateOfSurplusValueRequest struct {
	Surplus  commodity.SurplusValue  `json:"surplus"`
	Variable commodity.LabourMinutes `json:"variable"`
}

type rateOfSurplusValueResponse struct {
	Rate               float64                 `json:"rate"`
	Surplus            commodity.SurplusValue  `json:"surplus"`
	Variable           commodity.LabourMinutes `json:"variable"`
	ValueProduct       commodity.ValueProduct  `json:"value_product"`
	SurplusProduceFrac float64                 `json:"surplus_produce_fraction"`
}

// RateOfSurplusValue handles POST /v1/rate-of-surplus-value [Ch.9] — stateless probe.
func (h *Handler) RateOfSurplusValue(w http.ResponseWriter, r *http.Request) {
	var req rateOfSurplusValueRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rate, err := commodity.ComputeRate(req.Surplus, req.Variable)
	if errors.Is(err, commodity.ErrDivisionByZero) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rateOfSurplusValueResponse{
		Rate:               float64(rate),
		Surplus:            req.Surplus,
		Variable:           req.Variable,
		ValueProduct:       commodity.ComputeValueProduct(req.Variable, req.Surplus),
		SurplusProduceFrac: commodity.SurplusProduceRatio(
			commodity.LabourMinutes(req.Surplus), req.Variable,
		),
	})
}
