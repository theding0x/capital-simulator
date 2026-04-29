// api-gateway is the external entrypoint to the capital-simulator economy.
// It will eventually fan requests out to the domain services (commodity, agent,
// market, simulation-engine) and aggregate responses for the React UI.
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
)

const serviceName = "api-gateway"

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	addr := getenv("SERVICE_ADDR", ":8080")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	srv.HandleFunc("/v1/info", handleInfo)
	srv.MarkReady(true)

	if err := srv.Run(context.Background()); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func handleInfo(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]any{
		"service":     serviceName,
		"status":      "scaffolded",
		"description": "External entrypoint to capital-simulator. Fans out to domain services.",
		"downstream": []string{
			"commodity-service",
			"agent-service",
			"market-service",
			"simulation-engine",
		},
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
