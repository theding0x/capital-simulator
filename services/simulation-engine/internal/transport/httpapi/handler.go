package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
)

// Handler holds the logger plus the machinery store (added in Ch. 15)
// and the productivity fetcher (added in the Part IV bridge refactor).
// Pre-Ch.15 surplus and production endpoints remain stateless.
type Handler struct {
	Logger                *slog.Logger
	Machines              store.MachineStore
	Factories             store.FactoryStore
	Productivity          ProductivityFetcher
	GeneralLaw            store.GeneralLawStore
	HistoricalStages      store.HistoricalStageStore
	EnclosureEvents       store.EnclosureEventStore
	WageStatutes          store.WageStatuteStore
	VagrancyLaws          store.VagrancyLawStore
	FarmTenures           store.FarmTenureStore
	DomesticIndustries    store.DomesticIndustryStore
	CapitalOrigins        store.CapitalOriginStore
	ColonialTransfers     store.ColonialTransferStore
	NationalDebts         store.NationalDebtStore
	ProtectionSystems     store.ProtectionSystemStore
	Trajectories          store.AccumulationTrajectoryStore
	ColonialMarkets       store.ColonialLabourMarketStore
	ProductiveCircuits    store.ProductiveCircuitStore
	CommodityCircuits     store.CommodityCircuitStore
	MoneyCircuits         store.MoneyCircuitStore
	IndustrialCapitals    store.IndustrialCapitalStore
	AbodeStates           store.AbodeStateStore
	Turnovers             store.TurnoverStore
	Composition           store.CompositionStore
	AggregateTurnovers    store.AggregateTurnoverStore
	EconomistAttributions store.EconomistAttributionStore
	WorkingPeriods        store.WorkingPeriodStore
	ProductionTime        store.ProductionTimeStore
	PriceRevolutions      store.PriceRevolutionStore
	AnnualSurplusRates    store.AnnualSurplusRateStore
	SurplusCirculations   store.SurplusCirculationStore
	MoneyCapital          store.MoneyCapitalStore
	SimpleReproduction    store.SimpleReproductionSchemeStore
	ExtendedReproduction  store.ExtendedReproductionStore
	Scheduler             *engine.Scheduler
	EngineTicks           store.EngineTickStore
}

// Deps is the dependency bag passed to New. Each field maps to the
// matching Handler store; unset fields stay nil and the handlers that
// need them will fail at request time rather than at construction. Add
// new stores as fields here — callers that do not use them keep
// working unchanged.
type Deps struct {
	Machines              store.MachineStore
	Factories             store.FactoryStore
	Productivity          ProductivityFetcher
	GeneralLaw            store.GeneralLawStore
	HistoricalStages      store.HistoricalStageStore
	EnclosureEvents       store.EnclosureEventStore
	WageStatutes          store.WageStatuteStore
	VagrancyLaws          store.VagrancyLawStore
	FarmTenures           store.FarmTenureStore
	DomesticIndustries    store.DomesticIndustryStore
	CapitalOrigins        store.CapitalOriginStore
	ColonialTransfers     store.ColonialTransferStore
	NationalDebts         store.NationalDebtStore
	ProtectionSystems     store.ProtectionSystemStore
	Trajectories          store.AccumulationTrajectoryStore
	ColonialMarkets       store.ColonialLabourMarketStore
	ProductiveCircuits    store.ProductiveCircuitStore
	CommodityCircuits     store.CommodityCircuitStore
	MoneyCircuits         store.MoneyCircuitStore
	IndustrialCapitals    store.IndustrialCapitalStore
	AbodeStates           store.AbodeStateStore
	Turnovers             store.TurnoverStore
	Composition           store.CompositionStore
	AggregateTurnovers    store.AggregateTurnoverStore
	EconomistAttributions store.EconomistAttributionStore
	WorkingPeriods        store.WorkingPeriodStore
	ProductionTime        store.ProductionTimeStore
	PriceRevolutions      store.PriceRevolutionStore
	AnnualSurplusRates    store.AnnualSurplusRateStore
	SurplusCirculations   store.SurplusCirculationStore
	MoneyCapital          store.MoneyCapitalStore
	SimpleReproduction    store.SimpleReproductionSchemeStore
	ExtendedReproduction  store.ExtendedReproductionStore
	Scheduler             *engine.Scheduler
	EngineTicks           store.EngineTickStore
}

func New(logger *slog.Logger, d Deps) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		Logger:                logger,
		Machines:              d.Machines,
		Factories:             d.Factories,
		Productivity:          d.Productivity,
		GeneralLaw:            d.GeneralLaw,
		HistoricalStages:      d.HistoricalStages,
		EnclosureEvents:       d.EnclosureEvents,
		WageStatutes:          d.WageStatutes,
		VagrancyLaws:          d.VagrancyLaws,
		FarmTenures:           d.FarmTenures,
		DomesticIndustries:    d.DomesticIndustries,
		CapitalOrigins:        d.CapitalOrigins,
		ColonialTransfers:     d.ColonialTransfers,
		NationalDebts:         d.NationalDebts,
		ProtectionSystems:     d.ProtectionSystems,
		Trajectories:          d.Trajectories,
		ColonialMarkets:       d.ColonialMarkets,
		ProductiveCircuits:    d.ProductiveCircuits,
		CommodityCircuits:     d.CommodityCircuits,
		MoneyCircuits:         d.MoneyCircuits,
		IndustrialCapitals:    d.IndustrialCapitals,
		AbodeStates:           d.AbodeStates,
		Turnovers:             d.Turnovers,
		Composition:           d.Composition,
		AggregateTurnovers:    d.AggregateTurnovers,
		EconomistAttributions: d.EconomistAttributions,
		WorkingPeriods:        d.WorkingPeriods,
		ProductionTime:        d.ProductionTime,
		PriceRevolutions:      d.PriceRevolutions,
		AnnualSurplusRates:    d.AnnualSurplusRates,
		SurplusCirculations:   d.SurplusCirculations,
		MoneyCapital:          d.MoneyCapital,
		SimpleReproduction:    d.SimpleReproduction,
		ExtendedReproduction:  d.ExtendedReproduction,
		Scheduler:             d.Scheduler,
		EngineTicks:           d.EngineTicks,
	}
}

// massRequest accepts either {rate, variable_capital} or
// {rate, labour_power_value, worker_count}, or all fields for cross-validation.
type massRequest struct {
	SurplusLabour    int64  `json:"surplus_labour"`
	NecessaryLabour  int64  `json:"necessary_labour"`
	VariableCapital  *int64 `json:"variable_capital,omitempty"`
	LabourPowerValue *int64 `json:"labour_power_value,omitempty"`
	WorkerCount      *int   `json:"worker_count,omitempty"`
}

// massResponse is the response DTO for POST /v1/surplus/mass.
// Optional fields use pointers so they are omitted when not computed.
type massResponse struct {
	Rate            surplus.SurplusValueRate `json:"rate"`
	VariableCapital int64                    `json:"variable_capital"`
	WorkerCount     *int                     `json:"worker_count,omitempty"`
	Mass            int64                    `json:"mass"`
	MassByRate      int64                    `json:"mass_by_rate"`
	MassByWorkers   *int64                   `json:"mass_by_workers,omitempty"`
}

// ComputeMass handles POST /v1/surplus/mass.
func (h *Handler) ComputeMass(w http.ResponseWriter, r *http.Request) {
	var req massRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.NecessaryLabour <= 0 {
		writeError(w, http.StatusBadRequest, "necessary_labour must be positive")
		return
	}
	if req.SurplusLabour < 0 {
		writeError(w, http.StatusBadRequest, "surplus_labour cannot be negative")
		return
	}
	if req.VariableCapital != nil && *req.VariableCapital <= 0 {
		writeError(w, http.StatusBadRequest, "variable_capital must be positive")
		return
	}
	if req.LabourPowerValue != nil && *req.LabourPowerValue <= 0 {
		writeError(w, http.StatusBadRequest, "labour_power_value must be positive")
		return
	}
	if req.WorkerCount != nil && *req.WorkerCount <= 0 {
		writeError(w, http.StatusBadRequest, "worker_count must be positive")
		return
	}
	if req.VariableCapital == nil && (req.LabourPowerValue == nil || req.WorkerCount == nil) {
		writeError(w, http.StatusBadRequest, "provide variable_capital or both labour_power_value and worker_count")
		return
	}

	rate := surplus.SurplusValueRate{
		SurplusLabour:   req.SurplusLabour,
		NecessaryLabour: req.NecessaryLabour,
	}

	resp := massResponse{Rate: rate}

	var vc surplus.VariableCapital
	if req.VariableCapital != nil {
		vc = surplus.VariableCapital(*req.VariableCapital)
	}

	if req.LabourPowerValue != nil && req.WorkerCount != nil {
		v := surplus.LabourPowerValue(*req.LabourPowerValue)
		n := surplus.WorkerCount(*req.WorkerCount)
		wc := int(*req.WorkerCount)
		resp.WorkerCount = &wc
		mw := int64(surplus.MassByWorkers(v, rate, n))
		resp.MassByWorkers = &mw
		if req.VariableCapital == nil {
			vc = surplus.MinimumCapital(v, n)
		}
	}

	resp.VariableCapital = int64(vc)
	resp.MassByRate = int64(surplus.MassByRate(rate, vc))

	// Primary mass: prefer worker formula when available, else rate formula.
	if resp.MassByWorkers != nil {
		resp.Mass = *resp.MassByWorkers
	} else {
		resp.Mass = resp.MassByRate
	}

	writeJSON(w, http.StatusOK, resp)
}

// limitsResponse is the GET /v1/surplus/limits response shape.
type limitsResponse struct {
	AbsoluteWorkdayLimit int64  `json:"absolute_workday_limit"`
	MinimumCapital       *int64 `json:"minimum_capital,omitempty"`
	LabourPowerValue     *int64 `json:"labour_power_value,omitempty"`
	WorkerCount          *int   `json:"worker_count,omitempty"`
}

// GetLimits handles GET /v1/surplus/limits.
func (h *Handler) GetLimits(w http.ResponseWriter, r *http.Request) {
	limit := int64(surplus.AbsoluteWorkdayLimit)
	resp := limitsResponse{AbsoluteWorkdayLimit: limit}

	lpvStr := r.URL.Query().Get("labour_power_value")
	wcStr := r.URL.Query().Get("worker_count")
	if lpvStr != "" && wcStr != "" {
		lpv, err1 := strconv.ParseInt(lpvStr, 10, 64)
		wc, err2 := strconv.Atoi(wcStr)
		if err1 != nil || err2 != nil || lpv <= 0 || wc <= 0 {
			writeError(w, http.StatusBadRequest, "labour_power_value and worker_count must be positive integers")
			return
		}
		mc := int64(surplus.MinimumCapital(surplus.LabourPowerValue(lpv), surplus.WorkerCount(wc)))
		resp.MinimumCapital = &mc
		resp.LabourPowerValue = &lpv
		resp.WorkerCount = &wc
	}

	writeJSON(w, http.StatusOK, resp)
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return errors.New("invalid json: " + err.Error())
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// writeServerError logs err via the handler's structured logger and writes an
// opaque 500 response with a random correlation ID. Callers must not write
// anything to w after calling this.
func (h *Handler) writeServerError(w http.ResponseWriter, err error) {
	b := make([]byte, 8)
	if _, randErr := rand.Read(b); randErr != nil {
		panic("crypto/rand unavailable: " + randErr.Error())
	}
	id := hex.EncodeToString(b)
	h.Logger.Error("internal error", "error", err, "correlation_id", id)
	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"error":          "internal error",
		"correlation_id": id,
	})
}
