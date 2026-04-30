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

	srv.MarkReady(true)

	if err := srv.Run(context.Background()); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func handleInfo(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]any{
		"service":     serviceName,
		"status":      "ch-1-the-commodity",
		"description": "External entrypoint. Forwards /v1/commodities/* and /v1/exchange-ratio to commodity-service.",
		"downstream": []string{
			"commodity-service",
			"agent-service",
			"market-service",
			"simulation-engine",
		},
		"chapter": "Capital Vol. I, Ch. 1 - The Commodity",
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
