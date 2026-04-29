// commodity-service models the commodity: use-value, exchange-value, and value
// as crystallized socially-necessary labour-time. The cell-form of bourgeois
// wealth from which Capital Vol. I begins (Ch. 1).
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
)

const serviceName = "commodity-service"

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	addr := getenv("SERVICE_ADDR", ":8081")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	srv.HandleFunc("/v1/commodities", handleCommodities)
	srv.MarkReady(true)

	if err := srv.Run(context.Background()); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func handleCommodities(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]any{
		"service":     serviceName,
		"status":      "scaffolded",
		"description": "Will host commodities, their use-value, exchange-value, and value.",
		"chapter_ref": "Capital Vol. I, Ch. 1 - The Commodity",
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
