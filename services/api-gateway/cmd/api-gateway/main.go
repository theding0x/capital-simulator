// api-gateway is the external entrypoint to the capital-simulator economy.
// It fans requests out to the domain services (commodity, agent, market,
// simulation-engine) and aggregates responses for the React UI.
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
	"github.com/theding0x/capital-simulator/services/api-gateway/internal/proxy"
)

const serviceName = "api-gateway"

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	addr := getenv("SERVICE_ADDR", ":8080")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	srv.HandleFunc("GET /v1/info", handleInfo)

	// Reverse-proxy routes to commodity-service.
	commodityURL := getenv("COMMODITY_SERVICE_URL", "http://commodity-service:8081")
	commodityProxy, err := proxy.New(commodityURL, logger)
	if err != nil {
		logger.Error("failed to build commodity proxy", "err", err)
		os.Exit(1)
	}
	srv.Handle("/v1/commodities", commodityProxy)
	srv.Handle("/v1/commodities/{rest...}", commodityProxy)
	srv.Handle("/v1/exchange-ratio", commodityProxy)

	// Ch. 8-9 — capital decomposition and production accounts → commodity-service
	srv.Handle("/v1/capital", commodityProxy)
	srv.Handle("/v1/capital/{rest...}", commodityProxy)
	srv.Handle("/v1/production-accounts", commodityProxy)
	srv.Handle("/v1/production-accounts/{rest...}", commodityProxy)
	srv.Handle("/v1/rate-of-surplus-value", commodityProxy)

	// Reverse-proxy routes to market-service.
	marketURL := getenv("MARKET_SERVICE_URL", "http://market-service:8083")
	marketProxy, err := proxy.New(marketURL, logger)
	if err != nil {
		logger.Error("failed to build market proxy", "err", err)
		os.Exit(1)
	}
	srv.Handle("/v1/owners", marketProxy)
	srv.Handle("/v1/owners/{rest...}", marketProxy)
	srv.Handle("/v1/offers", marketProxy)
	srv.Handle("/v1/offers/{rest...}", marketProxy)
	srv.Handle("/v1/exchanges", marketProxy)
	srv.Handle("/v1/exchanges/{rest...}", marketProxy)
	srv.Handle("/v1/universal-equivalent", marketProxy)
	srv.Handle("/v1/money-commodity", marketProxy)
	srv.Handle("/v1/prices", marketProxy)
	srv.Handle("/v1/prices/{rest...}", marketProxy)
	srv.Handle("/v1/turnover-time", marketProxy)
	srv.Handle("/v1/turnover-time/{rest...}", marketProxy)
	srv.Handle("/v1/perishability", marketProxy)
	srv.Handle("/v1/market-separation", marketProxy)

	// Vol. I Ch. 3 — Money domain: hoards, payment obligations, world-money,
	// M = ΣP / V circulation calculator → market-service
	srv.Handle("/v1/hoards", marketProxy)
	srv.Handle("/v1/hoards/{rest...}", marketProxy)
	srv.Handle("/v1/payment-obligations", marketProxy)
	srv.Handle("/v1/payment-obligations/{rest...}", marketProxy)
	srv.Handle("/v1/world-money-transfers", marketProxy)
	srv.Handle("/v1/world-money-transfers/{rest...}", marketProxy)
	srv.Handle("/v1/circulation/money-required", marketProxy)

	// Vol. II Ch. 6 — Costs of Circulation → market-service
	srv.Handle("/v1/circulation-costs", marketProxy)
	srv.Handle("/v1/circulation-costs/{rest...}", marketProxy)
	srv.Handle("/v1/circulation-agents", marketProxy)
	srv.Handle("/v1/money-as-faux-frais", marketProxy)
	srv.Handle("/v1/commodity-supplies", marketProxy)
	srv.Handle("/v1/commodity-supplies/{rest...}", marketProxy)
	srv.Handle("/v1/transport-tariffs", marketProxy)
	srv.Handle("/v1/transport-legs", marketProxy)
	srv.Handle("/v1/transport-legs/{rest...}", marketProxy)

	// Reverse-proxy routes to agent-service.
	agentURL := getenv("AGENT_SERVICE_URL", "http://agent-service:8082")
	agentProxy, err := proxy.New(agentURL, logger)
	if err != nil {
		logger.Error("failed to build agent proxy", "err", err)
		os.Exit(1)
	}
	srv.Handle("/v1/agents", agentProxy)
	srv.Handle("/v1/agents/{rest...}", agentProxy)
	srv.Handle("/v1/circuit-probes", agentProxy)
	srv.Handle("/v1/exchange-simulations", agentProxy)

	// Ch. 6 — labour-power routes proxy to agent-service
	srv.Handle("/v1/workers", agentProxy)
	srv.Handle("/v1/workers/{rest...}", agentProxy)
	srv.Handle("/v1/capitalists", agentProxy)
	srv.Handle("/v1/capitalists/{rest...}", agentProxy)
	srv.Handle("/v1/labour-power", agentProxy)
	srv.Handle("/v1/labour-power/{rest...}", agentProxy)

	// Ch. 7 — labour-process routes proxy to agent-service
	srv.Handle("/v1/labour-processes", agentProxy)
	srv.Handle("/v1/labour-processes/{rest...}", agentProxy)

	// Ch. 10 — working-day routes proxy to agent-service
	srv.Handle("/v1/working-days", agentProxy)
	srv.Handle("/v1/working-days/{rest...}", agentProxy)
	srv.Handle("/v1/relay-schedules", agentProxy)
	srv.Handle("/v1/relay-schedules/{rest...}", agentProxy)

	// Ch. 13 — co-operation routes proxy to agent-service
	srv.Handle("/v1/cooperations", agentProxy)
	srv.Handle("/v1/cooperations/{rest...}", agentProxy)

	// Ch. 14 — manufacture routes proxy to agent-service
	srv.Handle("/v1/manufactures", agentProxy)
	srv.Handle("/v1/manufactures/{rest...}", agentProxy)

	// Ch. 17 — labour-scenario routes proxy to agent-service
	srv.Handle("/v1/labour-scenarios", agentProxy)
	srv.Handle("/v1/labour-scenarios/{rest...}", agentProxy)

	// Ch. 19 — wage-form routes proxy to agent-service
	srv.Handle("/v1/wage-forms", agentProxy)
	srv.Handle("/v1/wage-forms/{rest...}", agentProxy)

	// Ch. 20 — time-wage routes proxy to agent-service
	srv.Handle("/v1/time-wages", agentProxy)
	srv.Handle("/v1/time-wages/{rest...}", agentProxy)

	// Ch. 21 — piece-wage routes proxy to agent-service
	srv.Handle("/v1/piece-price", agentProxy)
	srv.Handle("/v1/sub-contracts", agentProxy)
	srv.Handle("/v1/sub-contracts/{rest...}", agentProxy)
	// Ch. 22 — National Differences in Wages
	srv.Handle("/v1/intensities", agentProxy)
	srv.Handle("/v1/intensities/{rest...}", agentProxy)
	srv.Handle("/v1/wages", agentProxy)
	srv.Handle("/v1/wages/{rest...}", agentProxy)
	srv.Handle("/v1/comparisons", agentProxy)

	// Reverse-proxy routes to simulation-engine.
	simURL := getenv("SIM_ENGINE_URL", "http://simulation-engine:8084")
	simProxy, err := proxy.New(simURL, logger)
	if err != nil {
		logger.Error("failed to build simulation-engine proxy", "err", err)
		os.Exit(1)
	}
	srv.Handle("/v1/sim/status", simProxy)

	// Issue #214 — automatic tick scheduler (operator control) → simulation-engine
	srv.Handle("/v1/engine/start", simProxy)
	srv.Handle("/v1/engine/stop", simProxy)
	srv.Handle("/v1/engine/status", simProxy)
	srv.Handle("/v1/engine/ticks", simProxy)

	// Ch. 11 — Rate and Mass of Surplus-Value → simulation-engine
	srv.Handle("/v1/surplus/mass", simProxy)
	srv.Handle("/v1/surplus/limits", simProxy)

	// Ch. 12 — Relative Surplus-Value → simulation-engine
	srv.Handle("/v1/production/working-day", simProxy)
	srv.Handle("/v1/production/working-day/shorten", simProxy)
	srv.Handle("/v1/production/rate-of-surplus-value", simProxy)
	srv.Handle("/v1/production/extra-surplus-value", simProxy)
	// Part IV bridge — Ch.13/14/15 productivity → Ch.12 relative surplus-value
	srv.Handle("/v1/production/relative-surplus", simProxy)
	srv.Handle("/v1/production/relative-surplus-from-productivity", simProxy)

	// Ch. 15 — Machinery and Modern Industry → simulation-engine
	srv.Handle("/v1/machines", simProxy)
	srv.Handle("/v1/machines/{rest...}", simProxy)
	srv.Handle("/v1/factories", simProxy)
	srv.Handle("/v1/factories/{rest...}", simProxy)

	// Ch. 16 — Absolute and Relative Surplus-Value → simulation-engine
	srv.Handle("/v1/surplus-value/absolute", simProxy)
	srv.Handle("/v1/surplus-value/relative", simProxy)
	srv.Handle("/v1/surplus-value/rate", simProxy)

	// Ch. 18 — Various Formula for the Rate of Surplus-Value → simulation-engine
	srv.Handle("/v1/surplus-value/rates", simProxy)

	// Ch. 23 — Simple Reproduction → simulation-engine
	srv.Handle("/v1/reproductions/simple", simProxy)
	srv.Handle("/v1/reproductions/repayment-period", simProxy)

	// Ch. 24 — The Transformation of Surplus-Value into Capital → simulation-engine
	srv.Handle("/v1/reproductions/extended", simProxy)
	srv.Handle("/v1/reproductions/split-surplus", simProxy)

	// Ch. 25 — The General Law of Capitalist Accumulation → simulation-engine
	srv.Handle("/v1/accumulation", simProxy)
	srv.Handle("/v1/accumulation/{rest...}", simProxy)

	// Ch. 26 — The Secret of Primitive Accumulation → simulation-engine
	srv.Handle("/v1/historical-stages", simProxy)
	srv.Handle("/v1/historical-stages/{rest...}", simProxy)

	// Ch. 27 — Expropriation of the Agricultural Population → simulation-engine
	srv.Handle("/v1/enclosure-events", simProxy)
	srv.Handle("/v1/enclosure-events/{rest...}", simProxy)

	// Ch. 28 — Bloody Legislation Against the Expropriated → simulation-engine
	srv.Handle("/v1/statutory-wages", simProxy)
	srv.Handle("/v1/statutory-wages/{rest...}", simProxy)

	// Ch. 29 — Genesis of the Capitalist Farmer → simulation-engine
	srv.Handle("/v1/farm-tenures", simProxy)
	srv.Handle("/v1/farm-tenures/{rest...}", simProxy)

	// Ch. 30 — Reaction of the Agricultural Revolution on Industry → simulation-engine
	srv.Handle("/v1/market-formation", simProxy)

	// Ch. 33 — The Modern Theory of Colonisation → simulation-engine
	srv.Handle("/v1/colonial-markets", simProxy)
	srv.Handle("/v1/colonial-markets/{rest...}", simProxy)

	// Vol. II Ch. 3 — The Circuit of Commodity-Capital → simulation-engine
	srv.Handle("/v1/commodity-circuits", simProxy)
	srv.Handle("/v1/commodity-circuits/{rest...}", simProxy)

	// Vol. II Ch. 1 — The Circuit of Money-Capital → simulation-engine
	srv.Handle("/v1/money-circuits", simProxy)
	srv.Handle("/v1/money-circuits/{rest...}", simProxy)

	// Vol. II Ch. 2 — The Circuit of Productive Capital → simulation-engine
	srv.Handle("/v1/productive-circuits", simProxy)
	srv.Handle("/v1/productive-circuits/{rest...}", simProxy)

	// Vol. II Ch. 4 — The Three Formulas of the Circuit → simulation-engine
	srv.Handle("/v1/industrial-capitals", simProxy)
	srv.Handle("/v1/industrial-capitals/{rest...}", simProxy)

	// Vol. II Ch. 7 — The Turnover Time and the Number of Turnovers → simulation-engine
	srv.Handle("/v1/turnovers", simProxy)
	srv.Handle("/v1/turnovers/{rest...}", simProxy)

	// Vol. II Ch. 8 — Fixed Capital and Circulating Capital → simulation-engine
	srv.Handle("/v1/capital-components", simProxy)
	srv.Handle("/v1/capital-components/{rest...}", simProxy)
	srv.Handle("/v1/fixed-items", simProxy)
	srv.Handle("/v1/fixed-items/{rest...}", simProxy)
	srv.Handle("/v1/circulating-cycles", simProxy)

	// Vol. II Ch. 9 — The Aggregate Turnover of Advanced Capital → simulation-engine
	srv.Handle("/v1/aggregate-turnovers", simProxy)
	srv.Handle("/v1/aggregate-turnovers/{rest...}", simProxy)

	// Vol. II Ch. 10 — Theories of Fixed and Circulating Capital: Physiocrats and Adam Smith → simulation-engine
	srv.Handle("/v1/circulation/economist-attributions", simProxy)

	// Vol. II Ch. 12 — The Working Period → simulation-engine
	srv.Handle("/v1/working-periods", simProxy)
	srv.Handle("/v1/working-periods/{rest...}", simProxy)
	srv.Handle("/v1/natural-constraints", simProxy)
	srv.Handle("/v1/natural-constraints/{rest...}", simProxy)

	// Vol. II Ch. 13 — The Time of Production → simulation-engine
	srv.Handle("/v1/industry-benchmarks", simProxy)
	srv.Handle("/v1/industry-benchmarks/{rest...}", simProxy)

	// Vol. II Ch. 15 — The Effects of a Change of Prices → simulation-engine
	srv.Handle("/v1/price-revolutions", simProxy)
	srv.Handle("/v1/price-revolutions/{rest...}", simProxy)
	srv.Handle("/v1/inventory-revaluations", simProxy)
	srv.Handle("/v1/speculative-holds", simProxy)

	// Vol. II Ch. 16 — The Turnover of Variable Capital → simulation-engine
	srv.Handle("/v1/variable-capital-advances", simProxy)
	srv.Handle("/v1/variable-capital-advances/{rest...}", simProxy)
	srv.Handle("/v1/annual-surplus-rates", simProxy)
	srv.Handle("/v1/annual-surplus-rate-contrasts", simProxy)
	srv.Handle("/v1/variable-capital-reproductions/{rest...}", simProxy)

	// Vol. II Ch. 17 — The Circulation of Surplus-Value → simulation-engine
	srv.Handle("/v1/reproduction/surplus-circulations", simProxy)
	srv.Handle("/v1/reproduction/surplus-circulations/{rest...}", simProxy)
	srv.Handle("/v1/reproduction/social-aggregate", simProxy)
	srv.Handle("/v1/reproduction/realisation-puzzles", simProxy)

	// Vol. II Ch. 18 — The Role of Money-Capital in Reproduction → simulation-engine
	srv.Handle("/v1/reproduction/apportionments", simProxy)
	srv.Handle("/v1/reproduction/apportionments/{rest...}", simProxy)
	srv.Handle("/v1/reproduction/department-reserves", simProxy)
	srv.Handle("/v1/reproduction/department-reserves/{rest...}", simProxy)
	srv.Handle("/v1/reproduction/circulating-money-mass", simProxy)
	srv.Handle("/v1/reproduction/circulating-money-mass/{rest...}", simProxy)
	srv.Handle("/v1/reproduction/wage-rotation-funds", simProxy)
	srv.Handle("/v1/reproduction/wage-rotation-funds/{rest...}", simProxy)
	srv.Handle("/v1/reproduction/inter-department-settlements", simProxy)
	srv.Handle("/v1/reproduction/inter-department-settlements/{rest...}", simProxy)

	// Vol. II Ch. 20 — Simple Reproduction (two-department scheme) → simulation-engine
	srv.Handle("/v1/reproduction/simple/schemes", simProxy)
	srv.Handle("/v1/reproduction/simple/schemes/{rest...}", simProxy)

	// Vol. II Ch. 21 — Accumulation and Reproduction on an Extended Scale → simulation-engine
	srv.Handle("/v1/reproduction/extended/schemes", simProxy)
	srv.Handle("/v1/reproduction/extended/schemes/{rest...}", simProxy)
	srv.Handle("/v1/reproduction/extended/multi-period", simProxy)
	srv.Handle("/v1/reproduction/extended/multi-period/{rest...}", simProxy)

	// Reverse-proxy routes to finance-service (Vol. III — profit, rent,
	// interest, credit, fictitious capital, the trinity formula). The
	// service is scaffolded empty in foundation Phase 3; each Vol. III
	// chapter PR uncomments / adds its route block below.
	financeURL := getenv("FINANCE_SERVICE_URL", "http://finance-service:8085")
	financeProxy, err := proxy.New(financeURL, logger)
	if err != nil {
		logger.Error("failed to build finance proxy", "err", err)
		os.Exit(1)
	}
	// Vol. III Ch. 1 — Cost-Price and Profit → finance-service
	srv.Handle("/v1/profit/cost-price", financeProxy)
	srv.Handle("/v1/profit/cost-price/{rest...}", financeProxy)
	srv.Handle("/v1/profit/profit-form", financeProxy)
	// Vol. III Ch. 2 — The Rate of Profit → finance-service
	srv.Handle("/v1/profit/rate", financeProxy)
	srv.Handle("/v1/profit/rate/{rest...}", financeProxy)
	// Vol. III, Ch. 9 — General Rate of Profit + Prices of Production → finance-service
	// srv.Handle("/v1/general-profit-rate", financeProxy)
	// srv.Handle("/v1/prices-of-production", financeProxy)
	// Vol. III, Ch. 21-25 — Interest-Bearing Capital + Credit → finance-service
	// srv.Handle("/v1/interest-rates", financeProxy)
	// srv.Handle("/v1/credit", financeProxy)
	// Vol. III, Ch. 37-47 — Ground-Rent (differential I/II, absolute) → finance-service
	// srv.Handle("/v1/ground-rent", financeProxy)
	// srv.Handle("/v1/ground-rent/{rest...}", financeProxy)

	srv.MarkReady(true)

	if err := srv.Run(context.Background()); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func handleInfo(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]any{
		"service":     serviceName,
		"status":      "ch-25-general-law",
		"description": "External entrypoint. Forwards commodity routes to commodity-service; owner/offer/exchange/price routes to market-service; agent/circuit/exchange-simulation/wage-form routes to agent-service; machine/factory/surplus-value/reproduction/accumulation routes to simulation-engine.",
		"downstream": []string{
			"commodity-service",
			"agent-service",
			"market-service",
			"simulation-engine",
			"finance-service",
		},
		"chapter": "Capital Vol. I, Ch. 23 — Simple Reproduction",
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
