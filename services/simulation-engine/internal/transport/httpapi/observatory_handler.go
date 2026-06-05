package httpapi

import (
	"errors"
	"net/http"
	"sort"
)

// observatorySnapshotResponse is the GET /v1/observatory/snapshot body: the
// whole field of industrial capitals, the aggregate vital-signs, and the
// engine's current tick/run state. Consumed by the Atlas page.
type observatorySnapshotResponse struct {
	Tick       int64              `json:"tick"`
	Running    bool               `json:"running"`
	IntervalMS int64              `json:"interval_ms"`
	Capitals   []fieldCapitalDTO  `json:"capitals"`
	Aggregate  aggregateVitalsDTO `json:"aggregate"`
}

type fieldCapitalDTO struct {
	ID              string `json:"id"`
	TotalPence      int64  `json:"total_pence"`
	MoneyPence      int64  `json:"money_pence"`
	ProductionPence int64  `json:"production_pence"`
	CommodityPence  int64  `json:"commodity_pence"`
	CostPricePence  int64  `json:"cost_price_pence"`
	SurplusPence    int64  `json:"surplus_pence"`
	Status          string `json:"status"`
	TurnoverNumber  int64  `json:"turnover_number"`
}

type aggregateVitalsDTO struct {
	TotalSocialCapitalPence int64 `json:"total_social_capital_pence"`
	CostPricePence          int64 `json:"cost_price_pence"`
	SurplusPence            int64 `json:"surplus_pence"`
	AvgRateOfProfitBP       int64 `json:"avg_rate_of_profit_bp"`
}

// GetObservatorySnapshot handles GET /v1/observatory/snapshot. The capitals
// array is always non-null; p̄′ = ΣS/ΣC (round-half-up basis points), 0 when
// ΣC == 0. Engine tick/running come from the scheduler when configured.
func (h *Handler) GetObservatorySnapshot(w http.ResponseWriter, r *http.Request) {
	if h.IndustrialCapitals == nil {
		h.writeServerError(w, errors.New("industrial capital store not configured"))
		return
	}
	field, err := h.IndustrialCapitals.FieldSnapshot(r.Context())
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	sort.Slice(field, func(i, j int) bool { return field[i].ID < field[j].ID })

	resp := observatorySnapshotResponse{Capitals: make([]fieldCapitalDTO, len(field))}
	var sumTotal, sumCost, sumSurplus int64
	for i, fc := range field {
		resp.Capitals[i] = fieldCapitalDTO{
			ID:              string(fc.ID),
			TotalPence:      int64(fc.TotalPence),
			MoneyPence:      int64(fc.MoneyPence),
			ProductionPence: int64(fc.ProductionPence),
			CommodityPence:  int64(fc.CommodityPence),
			CostPricePence:  int64(fc.CostPricePence),
			SurplusPence:    int64(fc.SurplusPence),
			Status:          string(fc.Status),
			TurnoverNumber:  fc.TurnoverNumber,
		}
		sumTotal += int64(fc.TotalPence)
		sumCost += int64(fc.CostPricePence)
		sumSurplus += int64(fc.SurplusPence)
	}
	resp.Aggregate = aggregateVitalsDTO{
		TotalSocialCapitalPence: sumTotal,
		CostPricePence:          sumCost,
		SurplusPence:            sumSurplus,
		AvgRateOfProfitBP:       rateBP(sumSurplus, sumCost),
	}
	if h.Scheduler != nil {
		st := h.Scheduler.Status()
		resp.Tick = st.Tick
		resp.Running = st.Running
		resp.IntervalMS = st.IntervalMS
	}
	writeJSON(w, http.StatusOK, resp)
}

// rateBP returns round-half-up(10000 * num / den) basis points; 0 when den <= 0.
func rateBP(num, den int64) int64 {
	if den <= 0 {
		return 0
	}
	return (num*10000 + den/2) / den
}
