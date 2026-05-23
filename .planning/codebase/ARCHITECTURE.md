<!-- refreshed: 2026-05-23 -->
# Architecture

**Analysis Date:** 2026-05-23

## System Overview

```text
┌──────────────────────────────────────────────────────────────────────┐
│                    web (React 18 + Vite + TS)                        │
│              nginx prod image · /api proxy → api-gateway             │
└─────────────────────────────┬────────────────────────────────────────┘
                              │ HTTP /v1/* (browser fetch → /api/*)
                              ▼
┌──────────────────────────────────────────────────────────────────────┐
│                   api-gateway  :8080                                 │
│  httputil.NewSingleHostReverseProxy per downstream service           │
│  `services/api-gateway/cmd/api-gateway/main.go`                     │
└───┬──────────┬──────────┬────────────┬─────────────┬────────────────┘
    │          │          │            │             │
    ▼          ▼          ▼            ▼             ▼
:8081      :8082      :8083        :8084          :8085
commodity  agent      market    simulation      finance
-service   -service   -service   -engine        -service
    │          │          │            │             │
    └──────────┴──────────┴────────────┴─────────────┘
                              │
              ┌───────────────┴──────────────┐
              ▼                              ▼
         ┌────────┐                    ┌─────────┐
         │ MySQL  │  one schema per    │  Redis  │
         │  :3306 │  service           │  :6379  │
         └────────┘                   └─────────┘
```

## Component Responsibilities

| Component | Port | Responsibility | Entry Point |
|-----------|------|----------------|-------------|
| `api-gateway` | 8080 | External entrypoint; fans all `/v1/*` requests out to domain services via `httputil.ReverseProxy` | `services/api-gateway/cmd/api-gateway/main.go` |
| `commodity-service` | 8081 | Vol. I — commodity, use-value, value, value-forms, fetishism, c+v decomposition, rate of surplus-value | `services/commodity-service/cmd/commodity-service/main.go` |
| `agent-service` | 8082 | Vol. I — workers, capitalists, labour-power, labour-process, working-day, cooperation, manufacture, wages (Ch. 4–22) | `services/agent-service/cmd/agent-service/main.go` |
| `market-service` | 8083 | Vol. I — exchange, money, prices; Vol. II — circulation phases | `services/market-service/cmd/market-service/main.go` |
| `simulation-engine` | 8084 | Vol. I — surplus-value, machinery, reproduction, accumulation, primitive accumulation; Vol. II — circuits of capital | `services/simulation-engine/cmd/simulation-engine/main.go` |
| `finance-service` | 8085 | Vol. III — profit, average rate of profit, prices of production, rent, interest, credit (scaffolded; no live routes yet) | `services/finance-service/cmd/finance-service/main.go` |
| `web` | 5173 (dev) / 80 (prod) | React 18 + Vite + TS SPA; chapter-by-chapter UI panels | `web/src/main.tsx` |

## Pattern Overview

**Overall:** Microservices monorepo behind a reverse-proxy gateway, chapter-driven domain evolution.

**Key Characteristics:**
- Single Go module (`github.com/theding0x/capital-simulator`) at repo root — no per-service `go.mod`
- Each domain service owns one MySQL schema and applies its own goose migrations on startup
- The gateway is a pure reverse proxy — no aggregation or transformation logic
- Domain logic is stateless pure functions in `internal/<pkg>/`; persistence is behind a Store interface
- Each new chapter of *Capital* is a branch + PR that adds domain types, migrations, HTTP handlers, and a React panel — all localized to specific services

## Layers

**Domain Layer:**
- Purpose: Pure functions modeling Marxist economic categories (value, labour, surplus-value, etc.)
- Location: `services/<svc>/internal/<pkg>/` (e.g., `services/commodity-service/internal/commodity/`)
- Contains: Exported structs, ID types, `New<Concept>ID()`, `Validate()`, pure computation functions
- Depends on: Nothing (no DB, no HTTP)
- Used by: HTTP transport layer, store layer

**Store Layer:**
- Purpose: Persistence boundary between domain and database
- Location: `services/<svc>/internal/store/`
- Contains: `store.go` (interface), `memory.go` (in-memory for tests), `mysql.go` (production)
- Depends on: Domain types, `pkg/mysql`
- Used by: HTTP transport layer (`Handler` struct holds a store reference)

**Transport Layer:**
- Purpose: HTTP handlers, request/response mapping, error translation
- Location: `services/<svc>/internal/transport/httpapi/`
- Contains: `handler.go` (Handler struct + helpers), `routes.go` (route registration), `<concept>_handler.go` per resource
- Depends on: Domain layer, store layer
- Used by: Service `main.go`

**Shared Packages:**
- Purpose: Cross-cutting infrastructure used by all services
- Location: `pkg/`
- Contains: `pkg/httpx/server.go` (HTTP server with `/healthz`, `/readyz`, graceful shutdown), `pkg/log/log.go` (slog wrapper), `pkg/mysql/client.go` + `migrate.go` (MySQL connect + goose runner)

## Data Flow

### REST Request Path (Browser → Domain Service)

1. Browser sends `fetch("/api/v1/commodities")` to nginx (`web/`)
2. nginx proxies `/api/*` → `api-gateway:8080`
3. `api-gateway` routes by path prefix using `httputil.NewSingleHostReverseProxy` (`services/api-gateway/cmd/api-gateway/main.go:35`)
4. Domain service handler receives request, calls Store method
5. Store method executes SQL (MySQL) or returns in-memory data
6. Handler writes JSON response; gateway transparently forwards it to nginx → browser

### Cross-Service HTTP Call (simulation-engine → agent-service)

The only current cross-service call is in `simulation-engine/internal/productivity/fetcher.go`:
- `Fetcher.Fetch()` issues `GET http://agent-service:8082/v1/cooperations/{id}` or `/v1/manufactures/{id}` to read productivity factors for the Ch. 12 bridge endpoint
- The `AGENT_SERVICE_URL` env var configures this target

### Migration Path (service startup)

1. `main.go` calls `pmysql.Connect()` → `pkg/mysql/client.go`
2. On success, calls `store.NewMySQL(ctx, cli.SQL)` → `services/<svc>/internal/store/mysql.go`
3. `NewMySQL` embeds SQL files via `//go:embed migrations` and calls `pkgmysql.Migrate(ctx, db, sub)`
4. `pkg/mysql/migrate.go` uses goose v3 to apply pending `.sql` files in numeric order

### Store Fallback Path

Every domain service `main.go` checks `MYSQL_DISABLED=true` (skip MySQL entirely) or `FALLBACK_MEMORY=true` (fall back on connection failure). In either case the in-memory store is used — enabling all services to run without a database for development.

## Key Abstractions

**Store Interface:**
- Purpose: Seam between HTTP/domain layer and database; enables Memory/MySQL swapping
- Examples: `services/commodity-service/internal/store/store.go`, `services/simulation-engine/internal/store/store.go`
- Pattern: Interface defined in `store.go`; implementations in `memory.go` and `mysql.go`; sentinel errors `ErrNotFound` / `ErrAlreadyExists` in `store.go`

**httpx.Server:**
- Purpose: Shared HTTP server scaffolding — all services use the same `pkg/httpx.New(cfg, logger)` wrapper
- Location: `pkg/httpx/server.go`
- Pattern: Wraps `http.ServeMux`; adds `/healthz` and `/readyz`; `MarkReady(bool)` is called after store init; `Run(ctx)` blocks until SIGINT/SIGTERM

**Commodity Circuit Model:**
- Purpose: The organizing concept M—C(Lp+Mp)…P…C'—M' maps directly to service boundaries
- Vol. I (production / P): commodity-service, agent-service, simulation-engine
- Vol. II (circulation / M—C and C'—M'): market-service, simulation-engine
- Vol. III (distribution of surplus-value / delta-M): finance-service (pending)
- Frontend: `web/src/chapters/registry.ts` tags each chapter with `circuitNode[]` for sidebar filtering

**ChapterDef / Registry:**
- Purpose: Single source of truth for which chapters exist and their status
- Location: `web/src/chapters/registry.ts`
- Pattern: `CHAPTERS: ChapterDef[]` with `id` (e.g., `"v1-ch01"`), `volume`, `number`, `title`, `circuitNode[]`, `status: "done" | "pending"`

## Entry Points

**Go service bootstrap:**
- Location: `services/<svc>/cmd/<svc>/main.go`
- Pattern: `applog.New(serviceName)` → `openStore()` (MySQL or Memory fallback) → `httpx.New(cfg, logger)` → `httpapi.Register(srv, handler)` → `srv.MarkReady(true)` → `srv.Run(ctx)`

**React app bootstrap:**
- Location: `web/src/main.tsx` → `web/src/App.tsx`
- Pattern: Single `App` component with `Sidebar` (chapter list filtered by `CircuitNode`), `CircuitSpine` (M—C…P…C'—M' diagram), and `ChapterShell` (renders the active chapter component). No router — `activeChapterId` state drives which panel is shown.

**API gateway bootstrap:**
- Location: `services/api-gateway/cmd/api-gateway/main.go`
- Pattern: One `proxy.New(url, logger)` per downstream service → `srv.Handle(pathPrefix, proxy)` for each route group. All routes are registered at startup; no dynamic routing.

## Architectural Constraints

- **Single module:** All Go code shares `github.com/theding0x/capital-simulator` module root. Per-service `go.mod` files are forbidden.
- **No ORM:** `database/sql` with `github.com/go-sql-driver/mysql` directly. All SQL is inline string constants in `mysql.go`.
- **No TLS termination:** Services listen on plain HTTP. HTTPS is handled at the cluster/ingress level.
- **No web router:** React SPA uses no react-router. One `App.tsx` with `useState` for active chapter.
- **Migrations are append-only:** Existing `.sql` migration files must never be edited. Add new numbered files only. Goose tracks by the integer prefix.
- **DDL in SQL files only:** `CREATE TABLE`, `ALTER TABLE`, etc. must not appear in Go source.
- **ID generation:** 96-bit hex from `crypto/rand`. Pattern: `New<Concept>ID()` returning a named string type. No `google/uuid`, no `math/rand`.
- **LabourMinutes:** `type LabourMinutes int64` is the canonical value-magnitude unit — never `float64`.
- **Global state:** None. Each service is stateless beyond its store connection. No module-level singletons except the logger (`applog.SetDefault`).
- **Threading:** Standard Go goroutine-per-request HTTP server. No worker threads.
- **Cross-service calls:** Only `simulation-engine → agent-service` (productivity fetcher). All other inter-domain communication goes through the UI / API gateway, not service-to-service.

## Anti-Patterns

### Inline DDL in Go source

**What happens:** Writing `CREATE TABLE ...` as a string in a `.go` file instead of a migration file.
**Why it's wrong:** Bypasses goose versioning; cannot be rolled back; not tracked by the migration sequence.
**Do this instead:** Add `services/<svc>/internal/store/migrations/NNNNN_v<V>_ch<NN>_<slug>.sql` with the DDL and reference it via `//go:embed migrations` in `mysql.go`.

### Editing existing migration files

**What happens:** Modifying a numbered `.sql` file that has already been applied.
**Why it's wrong:** Goose tracks by integer prefix; changing an applied migration creates schema drift between environments.
**Do this instead:** Add a new migration file with the next integer prefix.

### Putting business logic in handlers

**What happens:** Computing surplus-value, value-form, or other domain logic inside a handler function.
**Why it's wrong:** Handlers are untestable without HTTP scaffolding; domain logic must be testable as pure functions.
**Do this instead:** Implement computation in `services/<svc>/internal/<pkg>/` and call it from the handler.

## Error Handling

**Strategy:** Sentinel errors in store; `errors.Is` mapping in transport layer.

**Patterns:**
- Store layer returns `ErrNotFound` or `ErrAlreadyExists` (defined in `store.go`)
- Handler calls `writeStoreError(w, err)` which maps via `errors.Is`: `ErrNotFound` → 404, `ErrAlreadyExists` → 409, anything else → 500
- Domain layer validates inputs and returns descriptive `errors.New(...)` strings; handlers map these to 400
- `json.Decoder.DisallowUnknownFields()` is used everywhere; parse errors map to 400

## Cross-Cutting Concerns

**Logging:** `log/slog` via `pkg/log/log.go`. `applog.New(serviceName)` returns a structured logger. `applog.SetDefault(logger)` sets the global default. All handlers receive the logger via `Handler.Logger`.

**Validation:** Domain structs implement `Validate() error`. Called by handlers before store operations. Validation is a pure domain concern — not in handlers, not in stores.

**Health checks:** `pkg/httpx.Server` registers `/healthz` (always 200) and `/readyz` (200 only after `MarkReady(true)` is called post-store-init). Kubernetes liveness/readiness probes target these endpoints.

**Authentication:** None. All endpoints are open. Authentication is deferred to future volumes/chapters.

---

*Architecture analysis: 2026-05-23*
