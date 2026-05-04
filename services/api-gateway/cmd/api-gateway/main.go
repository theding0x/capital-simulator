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

	// Reverse-proxy routes to agent-service.
	agentURL := getenv("AGENT_SERVICE_URL", "http://agent-service:8082")
	agentProxy, err := proxy.New(agentURL, logger)
	if err != nil {
		logger.Error("failed to build agent proxy", "err", err)
		os.Exit(1)
	}
	srv.Handle("/v1/agents", agentProxy)
	srv.Handle("/v1/agents/{rest...}", agentProxy)

	srv.MarkReady(true)

	if err := srv.Run(context.Background()); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func handleInfo(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]any{
		"service":     serviceName,
		"status":      "ch-4-capital",
		"description": "External entrypoint. Forwards commodity routes to commodity-service; owner/offer/exchange/price routes to market-service.",
		"downstream": []string{
			"commodity-service",
			"agent-service",
			"market-service",
			"simulation-engine",
		},
		"chapter": "Capital Vol. I, Ch. 4 - The General Formula for Capital",
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
