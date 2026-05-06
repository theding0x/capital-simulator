// simulation-engine is the time-step orchestrator. It advances the simulated
// economy one period at a time, telling the domain services to produce,
// exchange, and accumulate. Without it, the rest of the world is static.
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

const serviceName = "simulation-engine"

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	addr := getenv("SERVICE_ADDR", ":8084")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	srv.HandleFunc("/v1/sim/status", handleStatus)

	// Ch. 11 — Rate and Mass of Surplus-Value
	h := httpapi.New(logger)
	httpapi.Register(srv, h)

	srv.MarkReady(true)

	if err := srv.Run(context.Background()); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func handleStatus(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]any{
		"service":     serviceName,
		"status":      "ch-11-rate-mass-surplus-value",
		"description": "Drives the simulated economy forward one period at a time.",
		"tick":        0,
		"running":     false,
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
