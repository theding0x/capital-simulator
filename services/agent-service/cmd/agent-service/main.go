// agent-service models the economic agents of the simulation: workers,
// capitalists, and (later) landowners. Agents are the bearers of the class
// relations through which capital reproduces itself.
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
)

const serviceName = "agent-service"

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	addr := getenv("SERVICE_ADDR", ":8082")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	srv.HandleFunc("/v1/agents", handleAgents)
	srv.MarkReady(true)

	if err := srv.Run(context.Background()); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func handleAgents(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]any{
		"service":     serviceName,
		"status":      "scaffolded",
		"description": "Will host workers, capitalists, and other agents who bear class relations.",
		"chapter_ref": "Introduced incrementally; relevant from Capital Vol. I, Ch. 4 onward.",
		"items":       []any{},
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
