// market-service models exchange and the emergence of money — where commodities
// meet one another as values and the universal equivalent crystallises.
// Capital Vol. I, Ch. 2: Exchange.
package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
	pmongo "github.com/theding0x/capital-simulator/pkg/mongo"
	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
	"github.com/theding0x/capital-simulator/services/market-service/internal/transport/httpapi"
)

const serviceName = "market-service"

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, mongoCli, err := openStore(ctx, logger)
	if err != nil {
		logger.Error("could not open any store", "err", err)
		os.Exit(1)
	}
	if mongoCli != nil {
		defer func() {
			shutdownCtx, c := context.WithTimeout(context.Background(), 5*time.Second)
			defer c()
			_ = mongoCli.Close(shutdownCtx)
		}()
	}

	addr := getenv("SERVICE_ADDR", ":8083")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	httpapi.Register(srv, httpapi.New(st, logger))
	srv.MarkReady(true)

	if err := srv.Run(ctx); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func openStore(ctx context.Context, logger *slog.Logger) (store.Store, *pmongo.Client, error) {
	if strings.EqualFold(os.Getenv("MONGO_DISABLED"), "true") {
		logger.Warn("MONGO_DISABLED=true; using in-memory store")
		return store.NewMemory(), nil, nil
	}

	cfg := pmongo.ConfigFromEnv(serviceName)
	dialCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	cli, err := pmongo.Connect(dialCtx, cfg)
	if err != nil {
		if strings.EqualFold(os.Getenv("FALLBACK_MEMORY"), "true") {
			logger.Warn("mongo connect failed; falling back to in-memory store", "err", err)
			return store.NewMemory(), nil, nil
		}
		return nil, nil, err
	}

	mstore, err := store.NewMongo(ctx, cli.DB)
	if err != nil {
		_ = cli.Close(context.Background())
		return nil, nil, err
	}
	logger.Info("mongo store ready", "uri", cfg.URI, "database", cfg.Database)
	return mstore, cli, nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
