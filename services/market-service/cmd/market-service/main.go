// market-service models exchange and circulation - where commodities meet
// money and the M-C-M' circuit unfolds. Capital Vol. I, Ch. 2-3 territory.
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
)

const serviceName = "market-service"

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	addr := getenv("SERVICE_ADDR", ":8083")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	srv.HandleFunc("/v1/markets", handleMarkets)
	srv.MarkReady(true)

	if err := srv.Run(context.Background()); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func handleMarkets(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]any{
		"service":     serviceName,
		"status":      "scaffolded",
		"description": "Will host exchange and circulation: prices, trades, the C-M-C and M-C-M' circuits.",
		"chapter_ref": "Capital Vol. I, Ch. 2-3 - Exchange and Money",
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
