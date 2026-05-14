# Goose Migrations Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace inline `CREATE TABLE IF NOT EXISTS` strings in each service's `NewMySQL` with tracked SQL migrations using `github.com/pressly/goose/v3`.

**Architecture:** Each service embeds its own SQL migration files via `//go:embed`. A shared helper `pkg/mysql/migrate.go` calls the goose Provider API, which creates a `goose_db_version` tracking table and applies any unapplied migrations on service startup. Migration files live at `services/<svc>/internal/store/migrations/` so the embed path is relative to the store package.

**Tech Stack:** Go 1.22 · `github.com/pressly/goose/v3` (Provider API) · `embed.FS` · MySQL 8

---

## File Map

| Action | Path | Responsibility |
|--------|------|----------------|
| Create | `pkg/mysql/migrate.go` | `Migrate(ctx, db, fs.FS)` — calls goose Provider, shared by all services |
| Create | `services/commodity-service/internal/store/migrations/00001_ch01_commodities.sql` | DDL for the `commodities` table |
| Modify | `services/commodity-service/internal/store/mysql.go` | Remove `createCommoditiesTable` const; add embed + call `Migrate` |
| Create | `services/market-service/internal/store/migrations/00001_ch02_barter.sql` | DDL for `owners`, `offers`, `exchanges` |
| Create | `services/market-service/internal/store/migrations/00002_ch03_money.sql` | DDL for `market_config`, `prices` |
| Modify | `services/market-service/internal/store/mysql.go` | Remove 5 `initXxx` consts; add embed + call `Migrate` |

---

## Task 1: Add the goose dependency

**Files:**
- Modify: `go.mod` (via `go get`)

- [ ] **Step 1: Fetch the module**

```bash
go get github.com/pressly/goose/v3
go mod tidy
```

Expected: `go.mod` gains a `require github.com/pressly/goose/v3 vX.Y.Z` line; `go.sum` is updated.

---

## Task 2: Create `pkg/mysql/migrate.go`

**Files:**
- Create: `pkg/mysql/migrate.go`

- [ ] **Step 1: Create the file**

```go
package mysql

import (
	"context"
	"database/sql"
	"io/fs"

	"github.com/pressly/goose/v3"
)

// Migrate applies all pending SQL migrations from fsys using goose.
// fsys must contain .sql files at its root (pass fs.Sub if needed).
func Migrate(ctx context.Context, db *sql.DB, fsys fs.FS) error {
	p, err := goose.NewProvider(goose.DialectMySQL, db, fsys)
	if err != nil {
		return err
	}
	_, err = p.Up(ctx)
	return err
}
```

- [ ] **Step 2: Verify compilation**

```bash
make vet
```

Expected: no errors (the function is exported but not yet called by any service).

---

## Task 3: Commodity-service migration + updated store

**Files:**
- Create: `services/commodity-service/internal/store/migrations/00001_ch01_commodities.sql`
- Modify: `services/commodity-service/internal/store/mysql.go`

- [ ] **Step 1: Create the migration file**

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS commodities (
    id               VARCHAR(24)  NOT NULL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    use_value_desc   TEXT         NOT NULL,
    use_value_unit   VARCHAR(100) NOT NULL,
    cl_kind          VARCHAR(100) NOT NULL DEFAULT '',
    cl_description   TEXT         NOT NULL,
    snlt_per_unit    BIGINT       NOT NULL DEFAULT 0,
    created_at       DATETIME(6)  NOT NULL,
    updated_at       DATETIME(6)  NOT NULL,
    UNIQUE KEY uq_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS commodities;
```

- [ ] **Step 2: Update `services/commodity-service/internal/store/mysql.go`**

Replace the `createCommoditiesTable` constant and `NewMySQL` with:

```go
package store

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"time"

	pkgmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

//go:embed migrations
var migrationsFS embed.FS

// MySQL is a MySQL-backed Store. Construct via NewMySQL.
type MySQL struct {
	db  *sql.DB
	now func() time.Time
}

// NewMySQL returns a Store backed by db and runs any pending migrations.
func NewMySQL(ctx context.Context, db *sql.DB) (*MySQL, error) {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	if err := pkgmysql.Migrate(ctx, db, sub); err != nil {
		return nil, err
	}
	return &MySQL{db: db, now: time.Now}, nil
}
```

The rest of `mysql.go` (all CRUD methods, `isDuplicate`) stays unchanged.

- [ ] **Step 3: Verify**

```bash
make vet
```

Expected: no errors.

---

## Task 4: Market-service migrations + updated store

**Files:**
- Create: `services/market-service/internal/store/migrations/00001_ch02_barter.sql`
- Create: `services/market-service/internal/store/migrations/00002_ch03_money.sql`
- Modify: `services/market-service/internal/store/mysql.go`

- [ ] **Step 1: Create Ch.2 migration (owners, offers, exchanges)**

```sql
-- +goose Up
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

-- +goose Down
DROP TABLE IF EXISTS exchanges;
DROP TABLE IF EXISTS offers;
DROP TABLE IF EXISTS owners;
```

- [ ] **Step 2: Create Ch.3 migration (market_config, prices)**

```sql
-- +goose Up
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS prices;
DROP TABLE IF EXISTS market_config;
```

- [ ] **Step 3: Update `services/market-service/internal/store/mysql.go`**

Remove the five `initXxx` string constants (`initSchema`, `initOffers`, `initExchanges`, `initConfig`, `initPrices`) and replace `NewMySQL` with:

```go
package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"strings"
	"time"

	pkgmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

//go:embed migrations
var migrationsFS embed.FS

const (
	keyUniversalEquivalent = "universal_equivalent"
	keyMoneyCommodity      = "money_commodity"
)

// MySQL is a MySQL-backed Store for market-service. Construct via NewMySQL.
type MySQL struct {
	db  *sql.DB
	now func() time.Time
}

// NewMySQL returns a Store backed by db and runs any pending migrations.
func NewMySQL(ctx context.Context, db *sql.DB) (*MySQL, error) {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	if err := pkgmysql.Migrate(ctx, db, sub); err != nil {
		return nil, err
	}
	return &MySQL{db: db, now: time.Now}, nil
}
```

The rest of `mysql.go` (all CRUD/query methods, `scanOffer`, `scanExchange`, `scanExchangeRow`, `isDuplicate`) stays unchanged.

- [ ] **Step 4: Verify**

```bash
make vet test build
```

Expected: all pass. The `isDuplicate` helper in market's store is identical to commodity's — both can stay in their respective packages.

---

## Task 5: Rebuild Docker stack

- [ ] **Step 1: Tear down the existing volume and restart**

```bash
docker compose down -v
docker compose up --build
```

Expected: commodity-service and market-service logs show goose applying migrations (look for `goose: successfully migrated database` or similar), then start serving requests normally.

- [ ] **Step 2: Re-seed if needed**

```powershell
Get-Content deploy\mysql\seed.sql | docker compose exec -T mysql mysql -u root -pcapital
```

Expected: no errors (tables exist from goose; `INSERT IGNORE` skips if seed was already applied).
