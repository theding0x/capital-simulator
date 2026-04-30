# MongoDB → MySQL Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace every MongoDB dependency in capital-simulator with MySQL 8, preserving the existing Store interface contract and the MONGO_DISABLED / FALLBACK_MEMORY fallback pattern.

**Architecture:** A new `pkg/mysql` package mirrors `pkg/mongo` (env-driven DSN, ping-on-connect). Each active service (commodity-service, market-service) gets a `store/mysql.go` that implements the existing `Store` interface using `database/sql` + `github.com/go-sql-driver/mysql`. Schemas are auto-created via an init SQL file mounted into the MySQL Docker image on first boot.

**Tech Stack:** Go 1.22 · `database/sql` · `github.com/go-sql-driver/mysql v1.8.x` · MySQL 8 · Docker Compose · Kubernetes (kustomize)

---

## File Map

| Action   | Path |
|----------|------|
| Create   | `pkg/mysql/client.go` |
| Create   | `deploy/mysql/init.sql` |
| Delete   | `pkg/mongo/client.go` (entire package) |
| Create   | `services/commodity-service/internal/store/mysql.go` |
| Delete   | `services/commodity-service/internal/store/mongo.go` |
| Modify   | `services/commodity-service/cmd/commodity-service/main.go` |
| Create   | `services/market-service/internal/store/mysql.go` |
| Delete   | `services/market-service/internal/store/mongo.go` |
| Modify   | `services/market-service/cmd/market-service/main.go` |
| Modify   | `go.mod` |
| Modify   | `docker-compose.yml` |
| Replace  | `deploy/k8s/infra/mongo.yaml` → `deploy/k8s/infra/mysql.yaml` |
| Modify   | `deploy/k8s/services/commodity-service.yaml` |
| Modify   | `deploy/k8s/services/market-service.yaml` |
| Modify   | `deploy/k8s/services/agent-service.yaml` |
| Modify   | `deploy/k8s/services/simulation-engine.yaml` |
| Modify   | `deploy/k8s/services/api-gateway.yaml` |
| Modify   | `services/commodity-service/internal/commodity/commodity.go` (strip bson tags) |
| Modify   | `services/commodity-service/internal/commodity/labour.go` (strip bson tags) |
| Modify   | `services/market-service/internal/market/market.go` (strip bson tags) |

---

### Task 1: Swap Go dependencies

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Edit go.mod — remove mongo driver, add MySQL driver**

Replace the entire `require` block in `go.mod` with:

```
module github.com/theding0x/capital-simulator

go 1.22

require github.com/go-sql-driver/mysql v1.8.1
```

All the indirect dependencies listed currently (`golang.org/x/crypto`, `github.com/golang/snappy`, etc.) belong to the MongoDB driver and can be removed entirely.

- [ ] **Step 2: User runs go mod tidy**

> **Sandbox note:** `go mod tidy` cannot run in the sandbox. Ask the user to run:
> ```bash
> go mod tidy
> ```
> This will resolve the MySQL driver and write `go.sum`. Do not proceed to Go edits until this succeeds.

---

### Task 2: Create `pkg/mysql/client.go`

**Files:**
- Create: `pkg/mysql/client.go`

This package is a direct replacement for `pkg/mongo/client.go`. It provides env-driven config, a ping-on-connect `Connect` function, and a `Close`.

- [ ] **Step 1: Write the file**

```go
// Package mysql wraps database/sql with the conventions used across
// capital-simulator: env-driven DSN, ping-on-connect, and a single *sql.DB
// handle that callers receive.
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config holds connection settings for a MySQL client.
type Config struct {
	DSN            string        // e.g. root:capital@tcp(mysql:3306)/commodity?parseTime=true&loc=UTC
	ConnectTimeout time.Duration // default 10s
	AppName        string        // informational only
}

// ConfigFromEnv reads MYSQL_DSN from the environment. If unset, it builds a
// default DSN pointing at mysql:3306 using the service name as the schema.
func ConfigFromEnv(service string) Config {
	dsn := getenvDefault("MYSQL_DSN",
		fmt.Sprintf("root:capital@tcp(mysql:3306)/%s?parseTime=true&loc=UTC", service))
	return Config{
		DSN:            dsn,
		ConnectTimeout: 10 * time.Second,
		AppName:        service,
	}
}

// Validate ensures the config has the minimum required values.
func (c Config) Validate() error {
	if c.DSN == "" {
		return errors.New("mysql: DSN is required")
	}
	return nil
}

// DB wraps a connected *sql.DB.
type DB struct {
	SQL *sql.DB
}

// Connect opens a *sql.DB using cfg.DSN and pings the server to verify the
// connection is live. Callers should defer Close.
func Connect(ctx context.Context, cfg Config) (*DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 10 * time.Second
	}

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("mysql: open: %w", err)
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("mysql: ping: %w", err)
	}
	return &DB{SQL: db}, nil
}

// Close closes the underlying *sql.DB. Safe to call on nil.
func (d *DB) Close() error {
	if d == nil || d.SQL == nil {
		return nil
	}
	return d.SQL.Close()
}

func getenvDefault(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
```

- [ ] **Step 2: Commit**

```bash
git add pkg/mysql/client.go go.mod go.sum
git commit -m "feat(mysql): add pkg/mysql client wrapper, swap go.mod to MySQL driver

Replaces pkg/mongo with pkg/mysql. ConfigFromEnv reads MYSQL_DSN from the
environment; defaults to root:capital@tcp(mysql:3306)/<service>?parseTime=true&loc=UTC.
Removes go.mongodb.org/mongo-driver and all its transitive dependencies."
```

---

### Task 3: Create MySQL schema init file

**Files:**
- Create: `deploy/mysql/init.sql`

The MySQL Docker image runs every `.sql` file in `/docker-entrypoint-initdb.d/` on first boot. This file creates all service schemas up front so services can connect immediately.

- [ ] **Step 1: Write the init SQL**

```sql
-- Capital-simulator MySQL schema initialisation.
-- Mounted into /docker-entrypoint-initdb.d/ by docker-compose.
-- Runs once on first container boot (when the data volume is empty).

CREATE DATABASE IF NOT EXISTS commodity
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

CREATE DATABASE IF NOT EXISTS market
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

CREATE DATABASE IF NOT EXISTS agent
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

CREATE DATABASE IF NOT EXISTS simulation
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;
```

- [ ] **Step 2: Commit**

```bash
git add deploy/mysql/init.sql
git commit -m "chore(mysql): add database schema init SQL for Docker bootstrap"
```

---

### Task 4: Commodity-service MySQL store

**Files:**
- Create: `services/commodity-service/internal/store/mysql.go`
- Delete: `services/commodity-service/internal/store/mongo.go`

The new `MySQL` type implements the `Store` interface using flat SQL columns. `UseValue` and `ConcreteLabour` (nested in the domain type) are stored as flat columns; they are reassembled when scanning.

- [ ] **Step 1: Delete the old file**

```bash
rm services/commodity-service/internal/store/mongo.go
```

- [ ] **Step 2: Write `services/commodity-service/internal/store/mysql.go`**

```go
package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

const createCommoditiesTable = `
CREATE TABLE IF NOT EXISTS commodities (
    id               VARCHAR(24)  NOT NULL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    use_value_desc   TEXT         NOT NULL,
    use_value_unit   VARCHAR(100) NOT NULL,
    cl_kind          VARCHAR(100) NOT NULL DEFAULT '',
    cl_description   TEXT         NOT NULL DEFAULT '',
    snlt_per_unit    BIGINT       NOT NULL DEFAULT 0,
    created_at       DATETIME(6)  NOT NULL,
    updated_at       DATETIME(6)  NOT NULL,
    UNIQUE KEY uq_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`

// MySQL is a MySQL-backed Store. Construct via NewMySQL.
type MySQL struct {
	db  *sql.DB
	now func() time.Time
}

// NewMySQL returns a Store backed by db and ensures the commodities table exists.
func NewMySQL(ctx context.Context, db *sql.DB) (*MySQL, error) {
	if _, err := db.ExecContext(ctx, createCommoditiesTable); err != nil {
		return nil, err
	}
	return &MySQL{db: db, now: time.Now}, nil
}

func (m *MySQL) Create(ctx context.Context, c commodity.Commodity) (commodity.Commodity, error) {
	if err := c.Validate(); err != nil {
		return commodity.Commodity{}, err
	}
	if c.ID.IsZero() {
		c.ID = commodity.NewID()
	}
	now := m.now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now

	const q = `INSERT INTO commodities
		(id, name, use_value_desc, use_value_unit, cl_kind, cl_description, snlt_per_unit, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), c.Name,
		c.UseValue.Description, c.UseValue.Unit,
		c.ConcreteLabour.Kind, c.ConcreteLabour.Description,
		int64(c.SNLTPerUnit), c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return commodity.Commodity{}, ErrAlreadyExists
		}
		return commodity.Commodity{}, err
	}
	return c, nil
}

func (m *MySQL) Get(ctx context.Context, id commodity.ID) (commodity.Commodity, error) {
	const q = `SELECT id, name, use_value_desc, use_value_unit, cl_kind, cl_description,
		snlt_per_unit, created_at, updated_at
		FROM commodities WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanCommodity(row)
}

func (m *MySQL) List(ctx context.Context) ([]commodity.Commodity, error) {
	const q = `SELECT id, name, use_value_desc, use_value_unit, cl_kind, cl_description,
		snlt_per_unit, created_at, updated_at
		FROM commodities ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []commodity.Commodity
	for rows.Next() {
		c, err := scanCommodityRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (m *MySQL) Update(ctx context.Context, id commodity.ID, u Update) (commodity.Commodity, error) {
	if u.IsEmpty() {
		return m.Get(ctx, id)
	}
	cur, err := m.Get(ctx, id)
	if err != nil {
		return commodity.Commodity{}, err
	}
	next := u.Apply(cur)
	if err := next.Validate(); err != nil {
		return commodity.Commodity{}, err
	}
	next.UpdatedAt = m.now().UTC()

	const q = `UPDATE commodities SET
		name = ?, use_value_desc = ?, use_value_unit = ?,
		cl_kind = ?, cl_description = ?, snlt_per_unit = ?, updated_at = ?
		WHERE id = ?`
	res, err := m.db.ExecContext(ctx, q,
		next.Name,
		next.UseValue.Description, next.UseValue.Unit,
		next.ConcreteLabour.Kind, next.ConcreteLabour.Description,
		int64(next.SNLTPerUnit), next.UpdatedAt,
		string(id),
	)
	if err != nil {
		if isDuplicate(err) {
			return commodity.Commodity{}, ErrAlreadyExists
		}
		return commodity.Commodity{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return commodity.Commodity{}, ErrNotFound
	}
	return next, nil
}

func (m *MySQL) Delete(ctx context.Context, id commodity.ID) error {
	res, err := m.db.ExecContext(ctx, `DELETE FROM commodities WHERE id = ?`, string(id))
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// scanCommodity scans a *sql.Row (single row) into a Commodity.
func scanCommodity(row *sql.Row) (commodity.Commodity, error) {
	var c commodity.Commodity
	var id string
	err := row.Scan(
		&id,
		&c.Name,
		&c.UseValue.Description,
		&c.UseValue.Unit,
		&c.ConcreteLabour.Kind,
		&c.ConcreteLabour.Description,
		&c.SNLTPerUnit,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return commodity.Commodity{}, ErrNotFound
	}
	if err != nil {
		return commodity.Commodity{}, err
	}
	c.ID = commodity.ID(id)
	return c, nil
}

// scanCommodityRow scans a *sql.Rows (cursor) into a Commodity.
func scanCommodityRow(rows *sql.Rows) (commodity.Commodity, error) {
	var c commodity.Commodity
	var id string
	err := rows.Scan(
		&id,
		&c.Name,
		&c.UseValue.Description,
		&c.UseValue.Unit,
		&c.ConcreteLabour.Kind,
		&c.ConcreteLabour.Description,
		&c.SNLTPerUnit,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return commodity.Commodity{}, err
	}
	c.ID = commodity.ID(id)
	return c, nil
}

// isDuplicate reports whether err is a MySQL duplicate-key error (1062).
func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	// go-sql-driver/mysql wraps driver errors; the message contains "1062" or "Duplicate entry".
	return containsAny(err.Error(), "1062", "Duplicate entry")
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}
```

> **Note on `isDuplicate`:** The proper way to detect duplicate-key errors with the MySQL driver is to type-assert to `*mysql.MySQLError` and check `Number == 1062`. However, that would require importing `github.com/go-sql-driver/mysql` in the store package just for the error type. Using the error string is a pragmatic alternative for this project. If the project later needs tighter coupling, replace with:
> ```go
> import driver "github.com/go-sql-driver/mysql"
> var me *driver.MySQLError
> if errors.As(err, &me) && me.Number == 1062 { return true }
> ```

- [ ] **Step 3: Commit**

```bash
git add services/commodity-service/internal/store/mysql.go
git rm services/commodity-service/internal/store/mongo.go
git commit -m "feat(commodity): replace MongoDB store with MySQL store

NewMySQL creates the commodities table on first connect. UseValue and
ConcreteLabour nested structs are stored as flat columns. The Store interface
contract (Create/Get/List/Update/Delete) is unchanged."
```

---

### Task 5: Update commodity-service main.go

**Files:**
- Modify: `services/commodity-service/cmd/commodity-service/main.go`

Replace the `pmongo` import and `openStore` function with a MySQL equivalent.

- [ ] **Step 1: Write the updated main.go**

Replace the entire file contents with:

```go
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

	httpapi.Register(srv, httpapi.New(st, logger))
	srv.MarkReady(true)

	if err := srv.Run(ctx); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

// openStore tries MySQL first; if MYSQL_DISABLED=true or the dial fails and
// FALLBACK_MEMORY=true, returns an in-memory store. Returns (store, *mysql.DB or nil, error).
func openStore(ctx context.Context, logger *slog.Logger) (store.Store, *pmysql.DB, error) {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

- [ ] **Step 2: Commit**

```bash
git add services/commodity-service/cmd/commodity-service/main.go
git commit -m "refactor(commodity): wire MySQL store in main.go

Replaces pmongo.Connect + store.NewMongo with pmysql.Connect + store.NewMySQL.
Fallback env vars renamed MONGO_DISABLED→MYSQL_DISABLED, FALLBACK_MEMORY unchanged."
```

---

### Task 6: Market-service MySQL store

**Files:**
- Create: `services/market-service/internal/store/mysql.go`
- Delete: `services/market-service/internal/store/mongo.go`

The market-service store has five tables: owners, offers, exchanges, market_config (two singletons), and prices.

- [ ] **Step 1: Delete the old file**

```bash
rm services/market-service/internal/store/mongo.go
```

- [ ] **Step 2: Write `services/market-service/internal/store/mysql.go`**

```go
package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

const initSchema = `
CREATE TABLE IF NOT EXISTS owners (
    id         VARCHAR(24)  NOT NULL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at DATETIME(6)  NOT NULL,
    updated_at DATETIME(6)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS offers (
    id                 VARCHAR(24)  NOT NULL PRIMARY KEY,
    owner_id           VARCHAR(24)  NOT NULL,
    commodity_id       VARCHAR(24)  NOT NULL,
    quantity           DOUBLE       NOT NULL,
    seeks_kind         VARCHAR(20)  NOT NULL DEFAULT '',
    seeks_commodity_id VARCHAR(24)  NOT NULL DEFAULT '',
    created_at         DATETIME(6)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS exchanges (
    id                     VARCHAR(24) NOT NULL PRIMARY KEY,
    giver_id               VARCHAR(24) NOT NULL,
    receiver_id            VARCHAR(24) NOT NULL,
    giver_commodity_id     VARCHAR(24) NOT NULL,
    giver_qty              DOUBLE      NOT NULL,
    receiver_commodity_id  VARCHAR(24) NOT NULL,
    receiver_qty           DOUBLE      NOT NULL,
    realised_value         BIGINT      NOT NULL,
    created_at             DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS market_config (
    key_name     VARCHAR(50) NOT NULL PRIMARY KEY,
    commodity_id VARCHAR(24) NOT NULL DEFAULT '',
    ts           DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS prices (
    commodity_id       VARCHAR(24) NOT NULL PRIMARY KEY,
    money_commodity_id VARCHAR(24) NOT NULL,
    amount             BIGINT      NOT NULL,
    updated_at         DATETIME(6) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`

const (
	keyUniversalEquivalent = "universal_equivalent"
	keyMoneyCommodity      = "money_commodity"
)

// MySQL is a MySQL-backed Store for market-service. Construct via NewMySQL.
type MySQL struct {
	db  *sql.DB
	now func() time.Time
}

// NewMySQL returns a Store backed by db and ensures all market tables exist.
// initSchema contains multiple statements, so the DSN must include multiStatements=true.
func NewMySQL(ctx context.Context, db *sql.DB) (*MySQL, error) {
	stmts := splitStatements(initSchema)
	for _, stmt := range stmts {
		if stmt == "" {
			continue
		}
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return nil, err
		}
	}
	return &MySQL{db: db, now: time.Now}, nil
}

// splitStatements splits a multi-statement SQL string on semicolons.
func splitStatements(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == ';' {
			out = append(out, trim(cur))
			cur = ""
		} else {
			cur += string(r)
		}
	}
	if t := trim(cur); t != "" {
		out = append(out, t)
	}
	return out
}

func trim(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\n' || s[start] == '\r' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\n' || s[end-1] == '\r' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// --- Owner ------------------------------------------------------------------

func (m *MySQL) CreateOwner(ctx context.Context, o market.Owner) (market.Owner, error) {
	if err := o.Validate(); err != nil {
		return market.Owner{}, err
	}
	if o.ID.IsZero() {
		o.ID = market.NewOwnerID()
	}
	now := m.now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now
	const q = `INSERT INTO owners (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q, string(o.ID), o.Name, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		if isDuplicate(err) {
			return market.Owner{}, ErrAlreadyExists
		}
		return market.Owner{}, err
	}
	return o, nil
}

func (m *MySQL) GetOwner(ctx context.Context, id market.OwnerID) (market.Owner, error) {
	const q = `SELECT id, name, created_at, updated_at FROM owners WHERE id = ?`
	var o market.Owner
	var rawID string
	err := m.db.QueryRowContext(ctx, q, string(id)).Scan(&rawID, &o.Name, &o.CreatedAt, &o.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.Owner{}, ErrNotFound
	}
	if err != nil {
		return market.Owner{}, err
	}
	o.ID = market.OwnerID(rawID)
	return o, nil
}

func (m *MySQL) ListOwners(ctx context.Context) ([]market.Owner, error) {
	const q = `SELECT id, name, created_at, updated_at FROM owners ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.Owner
	for rows.Next() {
		var o market.Owner
		var rawID string
		if err := rows.Scan(&rawID, &o.Name, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		o.ID = market.OwnerID(rawID)
		out = append(out, o)
	}
	return out, rows.Err()
}

// --- Offer ------------------------------------------------------------------

func (m *MySQL) CreateOffer(ctx context.Context, o market.Offer) (market.Offer, error) {
	if err := o.Validate(); err != nil {
		return market.Offer{}, err
	}
	if o.ID.IsZero() {
		o.ID = market.NewOfferID()
	}
	o.CreatedAt = m.now().UTC()
	const q = `INSERT INTO offers
		(id, owner_id, commodity_id, quantity, seeks_kind, seeks_commodity_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(o.ID), string(o.OwnerID), string(o.CommodityID),
		o.Quantity, o.SeeksKind, string(o.SeeksCommodityID), o.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return market.Offer{}, ErrAlreadyExists
		}
		return market.Offer{}, err
	}
	return o, nil
}

func (m *MySQL) GetOffer(ctx context.Context, id market.OfferID) (market.Offer, error) {
	const q = `SELECT id, owner_id, commodity_id, quantity, seeks_kind, seeks_commodity_id, created_at
		FROM offers WHERE id = ?`
	return scanOffer(m.db.QueryRowContext(ctx, q, string(id)))
}

func (m *MySQL) ListOffers(ctx context.Context) ([]market.Offer, error) {
	const q = `SELECT id, owner_id, commodity_id, quantity, seeks_kind, seeks_commodity_id, created_at
		FROM offers ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.Offer
	for rows.Next() {
		var o market.Offer
		var id, ownerID, commodityID, seeksCommodityID string
		if err := rows.Scan(&id, &ownerID, &commodityID, &o.Quantity,
			&o.SeeksKind, &seeksCommodityID, &o.CreatedAt); err != nil {
			return nil, err
		}
		o.ID = market.OfferID(id)
		o.OwnerID = market.OwnerID(ownerID)
		o.CommodityID = market.CommodityID(commodityID)
		o.SeeksCommodityID = market.CommodityID(seeksCommodityID)
		out = append(out, o)
	}
	return out, rows.Err()
}

func (m *MySQL) DeleteOffer(ctx context.Context, id market.OfferID) error {
	res, err := m.db.ExecContext(ctx, `DELETE FROM offers WHERE id = ?`, string(id))
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func scanOffer(row *sql.Row) (market.Offer, error) {
	var o market.Offer
	var id, ownerID, commodityID, seeksCommodityID string
	err := row.Scan(&id, &ownerID, &commodityID, &o.Quantity,
		&o.SeeksKind, &seeksCommodityID, &o.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.Offer{}, ErrNotFound
	}
	if err != nil {
		return market.Offer{}, err
	}
	o.ID = market.OfferID(id)
	o.OwnerID = market.OwnerID(ownerID)
	o.CommodityID = market.CommodityID(commodityID)
	o.SeeksCommodityID = market.CommodityID(seeksCommodityID)
	return o, nil
}

// --- Exchange ---------------------------------------------------------------

func (m *MySQL) CreateExchange(ctx context.Context, e market.Exchange) (market.Exchange, error) {
	if err := e.Validate(); err != nil {
		return market.Exchange{}, err
	}
	if e.ID.IsZero() {
		e.ID = market.NewExchangeID()
	}
	e.CreatedAt = m.now().UTC()
	const q = `INSERT INTO exchanges
		(id, giver_id, receiver_id, giver_commodity_id, giver_qty, receiver_commodity_id, receiver_qty, realised_value, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(e.ID),
		string(e.GiverID), string(e.ReceiverID),
		string(e.GiverCommodityID), e.GiverQty,
		string(e.ReceiverCommodityID), e.ReceiverQty,
		int64(e.RealisedValue), e.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return market.Exchange{}, ErrAlreadyExists
		}
		return market.Exchange{}, err
	}
	return e, nil
}

func (m *MySQL) GetExchange(ctx context.Context, id market.ExchangeID) (market.Exchange, error) {
	const q = `SELECT id, giver_id, receiver_id, giver_commodity_id, giver_qty,
		receiver_commodity_id, receiver_qty, realised_value, created_at
		FROM exchanges WHERE id = ?`
	return scanExchange(m.db.QueryRowContext(ctx, q, string(id)))
}

func (m *MySQL) ListExchanges(ctx context.Context) ([]market.Exchange, error) {
	const q = `SELECT id, giver_id, receiver_id, giver_commodity_id, giver_qty,
		receiver_commodity_id, receiver_qty, realised_value, created_at
		FROM exchanges ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.Exchange
	for rows.Next() {
		e, err := scanExchangeRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func scanExchange(row *sql.Row) (market.Exchange, error) {
	var e market.Exchange
	var id, giverID, receiverID, giverCID, receiverCID string
	var rv int64
	err := row.Scan(&id, &giverID, &receiverID, &giverCID, &e.GiverQty,
		&receiverCID, &e.ReceiverQty, &rv, &e.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.Exchange{}, ErrNotFound
	}
	if err != nil {
		return market.Exchange{}, err
	}
	e.ID = market.ExchangeID(id)
	e.GiverID = market.OwnerID(giverID)
	e.ReceiverID = market.OwnerID(receiverID)
	e.GiverCommodityID = market.CommodityID(giverCID)
	e.ReceiverCommodityID = market.CommodityID(receiverCID)
	e.RealisedValue = market.RealisedValue(rv)
	return e, nil
}

func scanExchangeRow(rows *sql.Rows) (market.Exchange, error) {
	var e market.Exchange
	var id, giverID, receiverID, giverCID, receiverCID string
	var rv int64
	if err := rows.Scan(&id, &giverID, &receiverID, &giverCID, &e.GiverQty,
		&receiverCID, &e.ReceiverQty, &rv, &e.CreatedAt); err != nil {
		return market.Exchange{}, err
	}
	e.ID = market.ExchangeID(id)
	e.GiverID = market.OwnerID(giverID)
	e.ReceiverID = market.OwnerID(receiverID)
	e.GiverCommodityID = market.CommodityID(giverCID)
	e.ReceiverCommodityID = market.CommodityID(receiverCID)
	e.RealisedValue = market.RealisedValue(rv)
	return e, nil
}

// --- UniversalEquivalent (singleton in market_config) -----------------------

func (m *MySQL) SetUniversalEquivalent(ctx context.Context, ue market.UniversalEquivalent) (market.UniversalEquivalent, error) {
	existing, err := m.GetUniversalEquivalent(ctx)
	if err == nil && existing.CommodityID == ue.CommodityID {
		return existing, nil
	}
	if ue.SetAt.IsZero() {
		ue.SetAt = m.now().UTC()
	}
	const q = `INSERT INTO market_config (key_name, commodity_id, ts) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE commodity_id = VALUES(commodity_id), ts = VALUES(ts)`
	_, err = m.db.ExecContext(ctx, q, keyUniversalEquivalent, string(ue.CommodityID), ue.SetAt)
	if err != nil {
		return market.UniversalEquivalent{}, err
	}
	return ue, nil
}

func (m *MySQL) GetUniversalEquivalent(ctx context.Context) (market.UniversalEquivalent, error) {
	const q = `SELECT commodity_id, ts FROM market_config WHERE key_name = ?`
	var commodityID string
	var ue market.UniversalEquivalent
	err := m.db.QueryRowContext(ctx, q, keyUniversalEquivalent).Scan(&commodityID, &ue.SetAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.UniversalEquivalent{}, ErrNotFound
	}
	if err != nil {
		return market.UniversalEquivalent{}, err
	}
	ue.CommodityID = market.CommodityID(commodityID)
	return ue, nil
}

// --- MoneyCommodity (singleton in market_config) ----------------------------

func (m *MySQL) SetMoneyCommodity(ctx context.Context, mc market.MoneyCommodity) (market.MoneyCommodity, error) {
	if mc.CreatedAt.IsZero() {
		mc.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO market_config (key_name, commodity_id, ts) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE commodity_id = VALUES(commodity_id), ts = VALUES(ts)`
	_, err := m.db.ExecContext(ctx, q, keyMoneyCommodity, string(mc.CommodityID), mc.CreatedAt)
	if err != nil {
		return market.MoneyCommodity{}, err
	}
	return mc, nil
}

func (m *MySQL) GetMoneyCommodity(ctx context.Context) (market.MoneyCommodity, error) {
	const q = `SELECT commodity_id, ts FROM market_config WHERE key_name = ?`
	var commodityID string
	var mc market.MoneyCommodity
	err := m.db.QueryRowContext(ctx, q, keyMoneyCommodity).Scan(&commodityID, &mc.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.MoneyCommodity{}, ErrNotFound
	}
	if err != nil {
		return market.MoneyCommodity{}, err
	}
	mc.CommodityID = market.CommodityID(commodityID)
	return mc, nil
}

// --- Price ------------------------------------------------------------------

func (m *MySQL) SetPrice(ctx context.Context, p market.Price) (market.Price, error) {
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = m.now().UTC()
	}
	const q = `INSERT INTO prices (commodity_id, money_commodity_id, amount, updated_at) VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE money_commodity_id = VALUES(money_commodity_id),
		amount = VALUES(amount), updated_at = VALUES(updated_at)`
	_, err := m.db.ExecContext(ctx, q,
		string(p.CommodityID), string(p.MoneyCommodityID), int64(p.Amount), p.UpdatedAt)
	if err != nil {
		return market.Price{}, err
	}
	return p, nil
}

func (m *MySQL) GetPrice(ctx context.Context, commodityID market.CommodityID) (market.Price, error) {
	const q = `SELECT commodity_id, money_commodity_id, amount, updated_at FROM prices WHERE commodity_id = ?`
	var p market.Price
	var cid, mcid string
	var amount int64
	err := m.db.QueryRowContext(ctx, q, string(commodityID)).Scan(&cid, &mcid, &amount, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.Price{}, ErrNotFound
	}
	if err != nil {
		return market.Price{}, err
	}
	p.CommodityID = market.CommodityID(cid)
	p.MoneyCommodityID = market.CommodityID(mcid)
	p.Amount = market.PriceAmount(amount)
	return p, nil
}

func (m *MySQL) ListPrices(ctx context.Context) ([]market.Price, error) {
	const q = `SELECT commodity_id, money_commodity_id, amount, updated_at FROM prices ORDER BY commodity_id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.Price
	for rows.Next() {
		var p market.Price
		var cid, mcid string
		var amount int64
		if err := rows.Scan(&cid, &mcid, &amount, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.CommodityID = market.CommodityID(cid)
		p.MoneyCommodityID = market.CommodityID(mcid)
		p.Amount = market.PriceAmount(amount)
		out = append(out, p)
	}
	return out, rows.Err()
}

// isDuplicate reports whether err is a MySQL duplicate-key error (1062).
func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return contains(s, "1062") || contains(s, "Duplicate entry")
}

func contains(s, sub string) bool {
	if len(sub) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
```

- [ ] **Step 3: Commit**

```bash
git add services/market-service/internal/store/mysql.go
git rm services/market-service/internal/store/mongo.go
git commit -m "feat(market): replace MongoDB store with MySQL store

Five tables: owners, offers, exchanges, market_config (two singletons keyed
by key_name), prices. Singleton upserts use INSERT ... ON DUPLICATE KEY UPDATE.
Store interface contract unchanged."
```

---

### Task 7: Update market-service main.go

**Files:**
- Modify: `services/market-service/cmd/market-service/main.go`

Same pattern as Task 5: replace `pmongo` with `pmysql`.

- [ ] **Step 1: Write the updated main.go**

Replace the entire file contents with:

```go
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
	pmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/market-service/internal/store"
	"github.com/theding0x/capital-simulator/services/market-service/internal/transport/httpapi"
)

const serviceName = "market-service"

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

	addr := getenv("SERVICE_ADDR", ":8083")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	httpapi.Register(srv, httpapi.New(st, logger))
	srv.MarkReady(true)

	if err := srv.Run(ctx); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func openStore(ctx context.Context, logger *slog.Logger) (store.Store, *pmysql.DB, error) {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

- [ ] **Step 2: Commit**

```bash
git add services/market-service/cmd/market-service/main.go
git commit -m "refactor(market): wire MySQL store in main.go"
```

---

### Task 8: Update docker-compose.yml

**Files:**
- Modify: `docker-compose.yml`

Replace the `mongo` service with `mysql`, mount the init SQL, and update all service env vars from `MONGO_URI`/`MONGO_DATABASE` to `MYSQL_DSN`.

- [ ] **Step 1: Write the updated docker-compose.yml**

Replace the entire file with:

```yaml
# Local development stack for capital-simulator.
# Bring up the whole economy with: docker compose up --build
name: capital-simulator

services:
  mysql:
    image: mysql:8
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: capital
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./deploy/mysql/init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-pcapital"]
      interval: 10s
      timeout: 5s
      retries: 10

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  api-gateway:
    build:
      context: .
      dockerfile: services/api-gateway/Dockerfile
    environment:
      LOG_LEVEL: info
      SERVICE_ADDR: ":8080"
      COMMODITY_SERVICE_URL: http://commodity-service:8081
      REDIS_ADDR: redis:6379
    ports:
      - "8080:8080"
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
      commodity-service:
        condition: service_started

  commodity-service:
    build:
      context: .
      dockerfile: services/commodity-service/Dockerfile
    environment:
      LOG_LEVEL: info
      SERVICE_ADDR: ":8081"
      MYSQL_DSN: "root:capital@tcp(mysql:3306)/commodity?parseTime=true&loc=UTC"
      REDIS_ADDR: redis:6379
    ports:
      - "8081:8081"
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy

  agent-service:
    build:
      context: .
      dockerfile: services/agent-service/Dockerfile
    environment:
      LOG_LEVEL: info
      SERVICE_ADDR: ":8082"
      MYSQL_DSN: "root:capital@tcp(mysql:3306)/agent?parseTime=true&loc=UTC"
      REDIS_ADDR: redis:6379
    ports:
      - "8082:8082"
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy

  market-service:
    build:
      context: .
      dockerfile: services/market-service/Dockerfile
    environment:
      LOG_LEVEL: info
      SERVICE_ADDR: ":8083"
      MYSQL_DSN: "root:capital@tcp(mysql:3306)/market?parseTime=true&loc=UTC"
      REDIS_ADDR: redis:6379
    ports:
      - "8083:8083"
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy

  simulation-engine:
    build:
      context: .
      dockerfile: services/simulation-engine/Dockerfile
    environment:
      LOG_LEVEL: info
      SERVICE_ADDR: ":8084"
      MYSQL_DSN: "root:capital@tcp(mysql:3306)/simulation?parseTime=true&loc=UTC"
      REDIS_ADDR: redis:6379
    ports:
      - "8084:8084"
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy

  web:
    build:
      context: ./web
    ports:
      - "5173:80"
    depends_on:
      - api-gateway

volumes:
  mysql_data:
  redis_data:
```

- [ ] **Step 2: Commit**

```bash
git add docker-compose.yml
git commit -m "chore(compose): replace MongoDB with MySQL 8

MySQL init SQL mounted at /docker-entrypoint-initdb.d/init.sql creates all
service schemas on first boot. MONGO_URI/MONGO_DATABASE replaced with MYSQL_DSN
per service. Health check uses mysqladmin ping."
```

---

### Task 9: Update Kubernetes infra

**Files:**
- Delete: `deploy/k8s/infra/mongo.yaml`
- Create: `deploy/k8s/infra/mysql.yaml`
- Modify: `deploy/k8s/services/commodity-service.yaml`
- Modify: `deploy/k8s/services/market-service.yaml`
- Modify: `deploy/k8s/services/agent-service.yaml`
- Modify: `deploy/k8s/services/simulation-engine.yaml`
- Modify: `deploy/k8s/services/api-gateway.yaml`

- [ ] **Step 1: Delete mongo.yaml, create mysql.yaml**

```bash
git rm deploy/k8s/infra/mongo.yaml
```

Write `deploy/k8s/infra/mysql.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql
  labels:
    app: mysql
spec:
  ports:
    - port: 3306
      targetPort: 3306
      name: mysql
  selector:
    app: mysql
  clusterIP: None
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: mysql
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - name: mysql
          image: mysql:8
          env:
            - name: MYSQL_ROOT_PASSWORD
              value: capital
          ports:
            - containerPort: 3306
              name: mysql
          volumeMounts:
            - name: data
              mountPath: /var/lib/mysql
            - name: initdb
              mountPath: /docker-entrypoint-initdb.d
          readinessProbe:
            exec:
              command: ["mysqladmin", "ping", "-h", "localhost", "-u", "root", "-pcapital"]
            initialDelaySeconds: 20
            periodSeconds: 10
          livenessProbe:
            exec:
              command: ["mysqladmin", "ping", "-h", "localhost", "-u", "root", "-pcapital"]
            initialDelaySeconds: 30
            periodSeconds: 30
      volumes:
        - name: initdb
          configMap:
            name: mysql-initdb
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 5Gi
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: mysql-initdb
data:
  init.sql: |
    CREATE DATABASE IF NOT EXISTS commodity
      DEFAULT CHARACTER SET utf8mb4
      DEFAULT COLLATE utf8mb4_unicode_ci;
    CREATE DATABASE IF NOT EXISTS market
      DEFAULT CHARACTER SET utf8mb4
      DEFAULT COLLATE utf8mb4_unicode_ci;
    CREATE DATABASE IF NOT EXISTS agent
      DEFAULT CHARACTER SET utf8mb4
      DEFAULT COLLATE utf8mb4_unicode_ci;
    CREATE DATABASE IF NOT EXISTS simulation
      DEFAULT CHARACTER SET utf8mb4
      DEFAULT COLLATE utf8mb4_unicode_ci;
```

- [ ] **Step 2: Update commodity-service.yaml**

In `deploy/k8s/services/commodity-service.yaml`, replace the env block:

```yaml
          env:
            - name: SERVICE_ADDR
              value: ":8081"
            - name: LOG_LEVEL
              value: info
            - name: MYSQL_DSN
              value: "root:capital@tcp(mysql:3306)/commodity?parseTime=true&loc=UTC"
            - name: REDIS_ADDR
              value: redis:6379
```

(Remove the MONGO_URI and MONGO_DATABASE entries.)

- [ ] **Step 3: Update market-service.yaml**

In `deploy/k8s/services/market-service.yaml`, replace the env block:

```yaml
          env:
            - name: SERVICE_ADDR
              value: ":8083"
            - name: LOG_LEVEL
              value: info
            - name: MYSQL_DSN
              value: "root:capital@tcp(mysql:3306)/market?parseTime=true&loc=UTC"
            - name: REDIS_ADDR
              value: redis:6379
```

- [ ] **Step 4: Update agent-service.yaml**

Replace env block:

```yaml
          env:
            - name: SERVICE_ADDR
              value: ":8082"
            - name: LOG_LEVEL
              value: info
            - name: MYSQL_DSN
              value: "root:capital@tcp(mysql:3306)/agent?parseTime=true&loc=UTC"
            - name: REDIS_ADDR
              value: redis:6379
```

- [ ] **Step 5: Update simulation-engine.yaml**

Replace env block:

```yaml
          env:
            - name: SERVICE_ADDR
              value: ":8084"
            - name: LOG_LEVEL
              value: info
            - name: MYSQL_DSN
              value: "root:capital@tcp(mysql:3306)/simulation?parseTime=true&loc=UTC"
            - name: REDIS_ADDR
              value: redis:6379
```

- [ ] **Step 6: Update api-gateway.yaml**

Remove MONGO_URI from the env block (api-gateway doesn't use it directly).

- [ ] **Step 7: Commit**

```bash
git add deploy/k8s/infra/mysql.yaml
git rm deploy/k8s/infra/mongo.yaml
git add deploy/k8s/services/
git commit -m "chore(k8s): replace MongoDB StatefulSet with MySQL 8

Init SQL is delivered via ConfigMap mounted at /docker-entrypoint-initdb.d.
All service deployments updated from MONGO_URI/MONGO_DATABASE to MYSQL_DSN."
```

---

### Task 10: Strip bson tags from domain types

**Files:**
- Modify: `services/commodity-service/internal/commodity/commodity.go`
- Modify: `services/commodity-service/internal/commodity/labour.go`
- Modify: `services/market-service/internal/market/market.go`

The `bson:"..."` struct tags are dead metadata now. Removing them keeps the types clean and avoids confusing future readers.

- [ ] **Step 1: Edit commodity.go — remove all bson tags**

In `services/commodity-service/internal/commodity/commodity.go`, change:

```go
type Commodity struct {
	ID             ID             `json:"id" bson:"_id"`
	Name           string         `json:"name" bson:"name"`
	UseValue       UseValue       `json:"use_value" bson:"use_value"`
	ConcreteLabour ConcreteLabour `json:"concrete_labour" bson:"concrete_labour"`
	SNLTPerUnit LabourMinutes `json:"snlt_per_unit" bson:"snlt_per_unit"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type UseValue struct {
	Description string `json:"description" bson:"description"`
	Unit string `json:"unit" bson:"unit"`
}
```

To:

```go
type Commodity struct {
	ID             ID             `json:"id"`
	Name           string         `json:"name"`
	UseValue       UseValue       `json:"use_value"`
	ConcreteLabour ConcreteLabour `json:"concrete_labour"`
	SNLTPerUnit LabourMinutes `json:"snlt_per_unit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UseValue struct {
	Description string `json:"description"`
	Unit string `json:"unit"`
}
```

- [ ] **Step 2: Edit labour.go — remove bson tags from ConcreteLabour**

In `services/commodity-service/internal/commodity/labour.go`, change:

```go
type ConcreteLabour struct {
	Kind string `json:"kind" bson:"kind"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}
```

To:

```go
type ConcreteLabour struct {
	Kind        string `json:"kind"`
	Description string `json:"description,omitempty"`
}
```

- [ ] **Step 3: Edit market.go — remove all bson tags**

In `services/market-service/internal/market/market.go`, remove all `bson:"..."` from every struct field. The json tags remain unchanged. For reference, the updated structs are:

```go
type Owner struct {
	ID        OwnerID   `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Offer struct {
	ID               OfferID     `json:"id"`
	OwnerID          OwnerID     `json:"owner_id"`
	CommodityID      CommodityID `json:"commodity_id"`
	Quantity         float64     `json:"quantity"`
	SeeksKind        string      `json:"seeks_kind"`
	SeeksCommodityID CommodityID `json:"seeks_commodity_id,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
}

type BarterRatio struct {
	CommodityA CommodityID `json:"commodity_a"`
	QtyA       float64     `json:"qty_a"`
	CommodityB CommodityID `json:"commodity_b"`
	QtyB       float64     `json:"qty_b"`
}

type Exchange struct {
	ID                  ExchangeID    `json:"id"`
	GiverID             OwnerID       `json:"giver_id"`
	ReceiverID          OwnerID       `json:"receiver_id"`
	GiverCommodityID    CommodityID   `json:"giver_commodity_id"`
	GiverQty            float64       `json:"giver_qty"`
	ReceiverCommodityID CommodityID   `json:"receiver_commodity_id"`
	ReceiverQty         float64       `json:"receiver_qty"`
	RealisedValue       RealisedValue `json:"realised_value"`
	CreatedAt           time.Time     `json:"created_at"`
}

type CircuitLeg struct {
	Kind        LegKind       `json:"kind"`
	CommodityID CommodityID   `json:"commodity_id"`
	MoneyID     CommodityID   `json:"money_id"`
	OwnerID     OwnerID       `json:"owner_id"`
	Value       RealisedValue `json:"value"`
}

type UniversalEquivalent struct {
	CommodityID CommodityID `json:"commodity_id"`
	SetAt       time.Time   `json:"set_at"`
}

type MoneyCommodity struct {
	CommodityID CommodityID `json:"commodity_id"`
	CreatedAt   time.Time   `json:"created_at"`
}

type Price struct {
	CommodityID      CommodityID `json:"commodity_id"`
	MoneyCommodityID CommodityID `json:"money_commodity_id"`
	Amount           PriceAmount `json:"amount"`
	UpdatedAt        time.Time   `json:"updated_at"`
}
```

- [ ] **Step 4: Commit**

```bash
git add services/commodity-service/internal/commodity/commodity.go \
        services/commodity-service/internal/commodity/labour.go \
        services/market-service/internal/market/market.go
git commit -m "chore: remove dead bson struct tags from domain types"
```

---

### Task 11: Delete pkg/mongo

**Files:**
- Delete: `pkg/mongo/client.go` (and the directory)

- [ ] **Step 1: Remove the package**

```bash
git rm pkg/mongo/client.go
rmdir pkg/mongo 2>/dev/null || true
```

- [ ] **Step 2: Commit**

```bash
git commit -m "chore: delete pkg/mongo — fully replaced by pkg/mysql"
```

---

### Task 12: Verify and hand off

- [ ] **Step 1: Ask user to run Go build and test**

The sandbox cannot run `go build` or `go test`. Ask the user to run:

```bash
go mod tidy
make vet test build
```

Expected: all packages compile, all tests pass. The store tests use `store.NewMemory()` and don't require a live database.

- [ ] **Step 2: Ask user to verify Docker stack**

```bash
docker compose down -v   # remove old mongo_data volume
docker compose up --build
```

Expected: mysql container becomes healthy, commodity-service and market-service log "mysql store ready", all services pass readiness probes.

- [ ] **Step 3: Update CLAUDE.md**

In `CLAUDE.md`, find the stack line:

```
Go 1.22 monorepo · React 18 + Vite + TS · MongoDB · Redis · Docker · k8s
```

Change to:

```
Go 1.22 monorepo · React 18 + Vite + TS · MySQL 8 · Redis · Docker · k8s
```

Also update the Anti-patterns section — remove the "An ORM. Use the official Mongo driver directly via the Store interface." line and replace with:

```
- An ORM. Use database/sql with github.com/go-sql-driver/mysql directly via the Store interface.
```

- [ ] **Step 4: Final commit**

```bash
git add CLAUDE.md
git commit -m "docs: update CLAUDE.md — MongoDB → MySQL in stack and anti-patterns"
```

---

## Self-Review

**Spec coverage check:**

| Requirement | Task |
|-------------|------|
| Remove go.mongodb.org/mongo-driver | Task 1 |
| New pkg/mysql client (env-driven, ping-on-connect) | Task 2 |
| Schema creation on first boot (Docker + k8s) | Tasks 3, 9 |
| commodity-service MySQL store (full CRUD) | Task 4 |
| commodity-service main.go wired to MySQL | Task 5 |
| market-service MySQL store (5 tables, singletons, upserts) | Task 6 |
| market-service main.go wired to MySQL | Task 7 |
| docker-compose.yml updated | Task 8 |
| k8s infra and service manifests updated | Task 9 |
| bson tags stripped | Task 10 |
| pkg/mongo deleted | Task 11 |
| CLAUDE.md updated | Task 12 |

**Placeholder scan:** No TBDs, no "add validation", no forward-references to undefined types. All code blocks are complete.

**Type consistency check:**
- `commodity.ID` → stored as `string` in SQL, cast back with `commodity.ID(id)` ✓
- `commodity.LabourMinutes` → stored as `int64`, scanned directly into the `int64`-based type ✓
- `market.OwnerID`, `OfferID`, `ExchangeID`, `CommodityID` → all stored as `string`, cast on read ✓
- `market.RealisedValue` (int64) → stored as `BIGINT`, scanned into `int64`, cast ✓
- `market.PriceAmount` (int64) → same pattern ✓
- `time.Time` → `DATETIME(6)` with `parseTime=true&loc=UTC` in DSN ✓
- `NewMySQL` in both stores called with `ctx, db.SQL` (`*sql.DB`) ✓
- `isDuplicate` helper exists in both store packages (not shared — avoids cross-package dep) ✓
