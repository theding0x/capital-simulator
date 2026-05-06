// commodity-service models the commodity: use-value, exchange-value, and value
// as crystallized socially-necessary labour-time. The cell-form of bourgeois
// wealth from which Capital Vol. I begins (Ch. 1).
package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
	pmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/commodity-service/internal/store"
	"github.com/theding0x/capital-simulator/services/commodity-service/internal/transport/httpapi"
)

const serviceName = "commodity-service"

// commodityStore embeds both store interfaces, enabling the openStore return
// type to be used by handlers requiring either Store or ProductionAccountStore.
type commodityStore interface {
	store.Store
	store.ProductionAccountStore
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
		defer func() {
			_ = mysqlDB.Close()
		}()
	}

	addr := getenv("SERVICE_ADDR", ":8081")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	httpapi.Register(srv, httpapi.New(st, st, logger))
	srv.MarkReady(true)

	if err := srv.Run(ctx); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

// openStore tries MySQL first; if MYSQL_DISABLED=true or the dial fails and
// FALLBACK_MEMORY=true, returns an in-memory store. Returns (commodityStore, *mysql.DB or nil, error).
func openStore(ctx context.Context, logger *slog.Logger) (commodityStore, *pmysql.DB, error) {
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

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

