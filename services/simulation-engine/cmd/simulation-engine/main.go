// simulation-engine is the time-step orchestrator. It advances the simulated
// economy one period at a time, telling the domain services to produce,
// exchange, and accumulate. Ch. 15 adds the first persistent domain: the
// machinery / factory loop.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
	pmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/productivity"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

const serviceName = "simulation-engine"

type machineryStore interface {
	store.MachineStore
	store.FactoryStore
	store.GeneralLawStore
	store.HistoricalStageStore
	store.EnclosureEventStore
	store.WageStatuteStore
	store.VagrancyLawStore
	store.FarmTenureStore
	store.DomesticIndustryStore
	store.CapitalOriginStore
	store.ColonialTransferStore
	store.NationalDebtStore
	store.ProtectionSystemStore
	store.AccumulationTrajectoryStore
	store.ColonialLabourMarketStore
	store.ProductiveCircuitStore
	store.CommodityCircuitStore
	store.MoneyCircuitStore
	store.IndustrialCapitalStore
	store.TurnoverStore
	store.CompositionStore
	store.AggregateTurnoverStore
	store.EconomistAttributionStore
	store.WorkingPeriodStore
	store.ProductionTimeStore
	store.PriceRevolutionStore
	store.AnnualSurplusRateStore
	store.SurplusCirculationStore
	store.MoneyCapitalStore
	store.SimpleReproductionSchemeStore
	store.ExtendedReproductionStore
	store.EngineTickStore
}

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, mysqlDB, err := openStore(ctx, logger)
	if err != nil {
		logger.Error("could not open any store", "err", err)
		os.Exit(1)
	}
	if mysqlDB != nil {
		defer func() { _ = mysqlDB.Close() }()
	}

	addr := getenv("SERVICE_ADDR", ":8084")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	agentURL := getenv("AGENT_SERVICE_URL", "http://agent-service:8082")
	pf := productivity.New(agentURL, st)

	// Issue #214 — the automatic tick scheduler. It advances every persisted
	// factory (Ch. 15) and reproduction scheme (Ch. 20/21) one period per pass
	// on a fixed interval, writing one engine_ticks audit row per pass.
	// Operator-controlled via POST /v1/engine/start and /stop; set
	// SIM_TICK_AUTOSTART=true to start it on boot. It is stopped gracefully
	// after the HTTP server drains on SIGTERM.
	scheduler := engine.NewScheduler(tickInterval(), []engine.Ticker{
		engine.NewFactoryTicker(st),
		engine.NewReproductionTicker(st),
	}, st, logger)

	srv.HandleFunc("/v1/sim/status", func(w http.ResponseWriter, _ *http.Request) {
		handleStatus(w, scheduler.Status())
	})

	h := httpapi.New(logger, httpapi.Deps{
		Machines:              st,
		Factories:             st,
		Productivity:          pf,
		GeneralLaw:            st,
		HistoricalStages:      st,
		EnclosureEvents:       st,
		WageStatutes:          st,
		VagrancyLaws:          st,
		FarmTenures:           st,
		DomesticIndustries:    st,
		CapitalOrigins:        st,
		ColonialTransfers:     st,
		NationalDebts:         st,
		ProtectionSystems:     st,
		Trajectories:          st,
		ColonialMarkets:       st,
		ProductiveCircuits:    st,
		CommodityCircuits:     st,
		MoneyCircuits:         st,
		IndustrialCapitals:    st,
		Turnovers:             st,
		Composition:           st,
		AggregateTurnovers:    st,
		EconomistAttributions: st,
		WorkingPeriods:        st,
		ProductionTime:        st,
		PriceRevolutions:      st,
		AnnualSurplusRates:    st,
		SurplusCirculations:   st,
		MoneyCapital:          st,
		SimpleReproduction:    st,
		ExtendedReproduction:  st,
		Scheduler:             scheduler,
		EngineTicks:           st,
	})
	httpapi.Register(srv, h)

	srv.MarkReady(true)

	if strings.EqualFold(os.Getenv("SIM_TICK_AUTOSTART"), "true") {
		scheduler.Start()
	}

	err = srv.Run(ctx)
	scheduler.Stop()
	if err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func openStore(ctx context.Context, logger *slog.Logger) (machineryStore, *pmysql.DB, error) {
	if strings.EqualFold(os.Getenv("MYSQL_DISABLED"), "true") {
		logger.Warn("MYSQL_DISABLED=true; using in-memory store")
		return store.NewMemory(), nil, nil
	}
	cfg := pmysql.ConfigFromEnv(serviceName)
	dialCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()
	cli, err := pmysql.Connect(dialCtx, cfg)
	if err != nil {
		if strings.EqualFold(os.Getenv("FALLBACK_MEMORY"), "true") {
			logger.Warn("mysql connect failed; falling back to in-memory store", "err", err)
			return store.NewMemory(), nil, nil
		}
		return nil, nil, err
	}
	initCtx, initCancel := context.WithTimeout(ctx, 30*time.Second)
	defer initCancel()
	mstore, err := store.NewMySQL(initCtx, cli.SQL)
	if err != nil {
		_ = cli.Close()
		return nil, nil, err
	}
	logger.Info("mysql store ready", "dsn_prefix", cfg.DSN[:min(len(cfg.DSN), 30)])
	return mstore, cli, nil
}

func handleStatus(w http.ResponseWriter, st engine.SchedulerStatus) {
	resp := map[string]any{
		"service":     serviceName,
		"status":      "ch-21-extended-reproduction",
		"description": "Drives the simulated economy forward one period at a time; persists machinery, factory, annual surplus-rate, surplus-circulation, money-capital, simple-reproduction and extended-reproduction scheme state.",
		"tick":        st.Tick,
		"running":     st.Running,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// tickInterval reads SIM_TICK_INTERVAL (a Go duration like "5s" or "250ms")
// and falls back to engine.DefaultTickInterval when unset or invalid.
func tickInterval() time.Duration {
	if v := os.Getenv("SIM_TICK_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return engine.DefaultTickInterval
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
