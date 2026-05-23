# Codebase Structure

**Analysis Date:** 2026-05-23

## Directory Layout

```
capital-simulator/
├── pkg/                        # Shared Go libraries (used by all services)
│   ├── httpx/server.go         # HTTP server with /healthz, /readyz, graceful shutdown
│   ├── log/log.go              # slog wrapper
│   └── mysql/                  # MySQL client + goose migration runner
│       ├── client.go
│       └── migrate.go
├── services/                   # One directory per microservice
│   ├── api-gateway/
│   ├── commodity-service/
│   ├── agent-service/
│   ├── market-service/
│   ├── simulation-engine/
│   └── finance-service/
├── web/                        # Vite + React 18 + TypeScript SPA
│   └── src/
│       ├── App.tsx             # Root component — no router
│       ├── api.ts              # All fetch() calls to /api/v1/*
│       ├── types.ts            # TypeScript mirror of Go response structs
│       ├── chapters/
│       │   ├── registry.ts     # Authoritative chapter list (CHAPTERS[])
│       │   ├── vol1/           # Ch01–Ch33 components
│       │   ├── vol2/           # Vol. II chapter components
│       │   └── vol3/           # Vol. III (placeholder)
│       └── components/         # Shared UI components (Sidebar, ChapterShell, etc.)
├── deploy/
│   ├── k8s/                    # Kustomize manifests (namespace + infra + per-service)
│   └── mysql/
│       ├── init.sql            # CREATE DATABASE per service — auto-run on first boot
│       └── seed.sql            # Optional bootstrap seeds
├── docs/
│   └── architecture.md         # Authoritative roadmap table + topology diagram
├── go.mod                      # Single module: github.com/theding0x/capital-simulator
├── go.sum
├── Makefile                    # vet, test, build, run-<service>, mysql-bootstrap
└── docker-compose.yml          # Local full-stack: mysql, redis, all services, web
```

## Service Internal Layout

Every service follows the same internal structure:

```
services/<svc>/
├── Dockerfile                              # Multi-stage, distroless; build context = repo root
├── cmd/<svc>/main.go                       # Service entrypoint
└── internal/
    ├── <pkg>/                              # Domain package(s) — pure functions, no I/O
    │   ├── <concept>.go                    # Exported types, New<Concept>ID(), Validate()
    │   └── <concept>_test.go              # Unit tests (t.Parallel(), Marx fixtures)
    ├── store/
    │   ├── store.go                        # Store interface(s) + sentinel errors + Update types
    │   ├── memory.go                       # In-memory implementation (used in tests)
    │   ├── memory_test.go                  # Store interface compliance tests
    │   ├── mysql.go                        # MySQL implementation; runs migrations on New
    │   └── migrations/
    │       ├── NNNNN_v<V>_ch<NN>_<slug>.sql  # Schema DDL (goose Up/Down)
    │       └── NNNNN_v<V>_ch<NN>_seed.sql   # Marx-faithful seed data (IDs: 5eed…)
    └── transport/httpapi/
        ├── handler.go                      # Handler struct; New(store, logger); shared helpers
        ├── routes.go                       # Register(s *httpx.Server, h *Handler)
        └── <concept>_handler.go           # Per-resource handler methods + test file
```

## Directory Purposes

**`pkg/`:**
- Purpose: Shared infrastructure imported by every service; no domain logic here
- Key files: `pkg/httpx/server.go`, `pkg/mysql/client.go`, `pkg/mysql/migrate.go`, `pkg/log/log.go`

**`services/api-gateway/`:**
- Purpose: Pure reverse proxy — routes `/v1/*` to domain services by path prefix
- Key files: `services/api-gateway/cmd/api-gateway/main.go` (all route registrations live here), `services/api-gateway/internal/proxy/proxy.go`
- No domain logic, no store, no migrations

**`services/commodity-service/`:**
- Purpose: Commodity, value, value-forms, fetishism, c+v decomposition (Vol. I Ch. 1, 8-9)
- Domain package: `services/commodity-service/internal/commodity/`
- Key domain files: `commodity.go`, `value.go`, `valueform.go`, `fetishism.go`, `capital.go`, `production_account.go`, `labour.go`

**`services/agent-service/`:**
- Purpose: Workers, capitalists, labour-power, labour-process, working-day, cooperation, manufacture, wages (Vol. I Ch. 4–22)
- Domain package: `services/agent-service/internal/agent/`
- Key domain files: `agent.go`, `labour_power.go`, `labour_process.go`, `working_day.go`, `cooperation.go`, `manufacture.go`, `wage.go`, `time_wage.go`, `piece_wage.go`, `national_wage.go`, `labour_scenario.go`

**`services/market-service/`:**
- Purpose: Exchange, money, prices (Vol. I Ch. 2-3); Vol. II circulation phases
- Domain package: `services/market-service/internal/market/`

**`services/simulation-engine/`:**
- Purpose: Surplus-value computation, machinery, reproduction, accumulation, primitive accumulation (Vol. I Ch. 11–33); Vol. II circuits of capital
- Multiple domain packages under `services/simulation-engine/internal/`:
  - `engine/` — tick engine
  - `machinery/` — machines, factories (Ch. 15)
  - `production/` — working-day and relative surplus-value (Ch. 12)
  - `surplus/` — absolute/relative surplus-value (Ch. 11, 16, 18)
  - `simulation/` — accumulation, reproduction, primitive accumulation, historical stages (Ch. 23–33)
  - `circulation/` — circuits of capital (Vol. II Ch. 2-3)
  - `productivity/` — cross-service HTTP fetcher for cooperation/manufacture/factory productivity factors
  - `labour/` — shared labour types

**`services/finance-service/`:**
- Purpose: Vol. III — profit, prices of production, rent, interest (scaffolded empty; routes commented out in gateway)
- Status: Store interface and MySQL scaffolding present; no domain logic or active routes yet

**`web/src/`:**
- Purpose: React SPA; one panel per chapter of Capital
- Key files: `App.tsx` (root shell), `api.ts` (all HTTP calls), `types.ts` (TypeScript mirror of Go structs)
- Chapter panels: `web/src/chapters/vol1/Ch<NN><Title>.tsx`, `web/src/chapters/vol2/Ch<NN><Title>.tsx`
- Chapter registry: `web/src/chapters/registry.ts` (CHAPTERS array — authoritative list)
- Shared components: `web/src/components/` (Sidebar, ChapterShell, CircuitSpine, TurnoverPlayer)

**`deploy/k8s/`:**
- Purpose: Kustomize manifests for Kubernetes deployment
- Layout: `namespace.yaml`, `infra/` (mysql.yaml, redis.yaml), `services/<svc>.yaml` per service

**`docs/`:**
- Purpose: Human-readable architecture docs and the authoritative chapter roadmap
- Key file: `docs/architecture.md` — roadmap table with `Done` / `Pending` status; next chapter to implement is the first `Pending` row

## Key File Locations

**Entry Points:**
- `services/<svc>/cmd/<svc>/main.go` — service main (one per service)
- `web/src/main.tsx` — React app mount
- `web/src/App.tsx` — root component

**Route Definitions:**
- Gateway routes: `services/api-gateway/cmd/api-gateway/main.go` (all `srv.Handle(...)` calls)
- Per-service routes: `services/<svc>/internal/transport/httpapi/routes.go`

**Store Interfaces:**
- `services/<svc>/internal/store/store.go` (interface + sentinel errors)

**Migrations:**
- `services/<svc>/internal/store/migrations/NNNNN_v<V>_ch<NN>_<slug>.sql`
- Database init: `deploy/mysql/init.sql` (`CREATE DATABASE IF NOT EXISTS <name>`)

**React API Client:**
- `web/src/api.ts` — all `fetch()` calls; add a new function per new endpoint

**React Type Definitions:**
- `web/src/types.ts` — TypeScript mirror of Go response structs (keep in sync with Go)

**Chapter Components:**
- Vol. 1: `web/src/chapters/vol1/Ch<NN><Title>.tsx`
- Vol. 2: `web/src/chapters/vol2/Ch<NN><Title>.tsx`
- Vol. 3: `web/src/chapters/vol3/Ch<NN><Title>.tsx`

**Chapter Registry:**
- `web/src/chapters/registry.ts` — CHAPTERS[] array; IDs are `v<V>-ch<NN>`; `status: "done" | "pending"`

**Environment Config Pattern:**
- `pkg/mysql/client.go` — `ConfigFromEnv(service)` reads `MYSQL_DSN`, defaults to `root:capital@tcp(mysql:3306)/<service>?parseTime=true&loc=UTC`
- Gateway service URLs: `COMMODITY_SERVICE_URL`, `AGENT_SERVICE_URL`, `MARKET_SERVICE_URL`, `SIM_ENGINE_URL`, `FINANCE_SERVICE_URL`

**Next Chapter to Implement:**
- `docs/architecture.md` roadmap table — first row with status `Pending` is the next chapter

## Naming Conventions

**Go files:**
- Domain types: `<concept>.go` (e.g., `commodity.go`, `labour_power.go`)
- Tests: `<concept>_test.go` co-located in same package
- Handlers: `<concept>_handler.go` and `<concept>_handler_test.go`
- Store files always named `store.go`, `memory.go`, `mysql.go`

**Migration files:**
- Format: `NNNNN_v<V>_ch<NN>_<slug>.sql` where `V` ∈ `{1,2,3}`, `NN` is zero-padded chapter, `slug` is descriptive
- Seed files: `NNNNN_v<V>_ch<NN>_seed.sql`
- Integer prefix controls goose ordering; token `v<V>_ch<NN>` is informational only

**Go types:**
- ID types: `type <Concept>ID string` with `New<Concept>ID() <Concept>ID` using `crypto/rand` 96-bit hex
- Value magnitudes: `type LabourMinutes int64`
- Seed IDs: `5eed00000000000000<CC><XX>` (recognisable; never collide with `NewID()` output)

**React files:**
- Chapter components: `Ch<NN><PascalTitle>.tsx` in the appropriate `vol<V>/` subdirectory
- CSS (if needed): `Ch<NN><PascalTitle>.css` co-located with the `.tsx`

**Chapter IDs in registry:**
- Format: `v<V>-ch<NN>` (e.g., `v1-ch01`, `v2-ch03`)

## Where to Add New Code

**New chapter in Vol. I (agent-service example):**
- Domain type: `services/agent-service/internal/agent/<concept>.go`
- Tests: `services/agent-service/internal/agent/<concept>_test.go`
- Store interface additions: `services/agent-service/internal/store/store.go`
- Memory store additions: `services/agent-service/internal/store/memory.go`
- MySQL store additions: `services/agent-service/internal/store/mysql.go`
- Schema migration: `services/agent-service/internal/store/migrations/NNNNN_v1_ch<NN>_<slug>.sql`
- Seed migration: `services/agent-service/internal/store/migrations/NNNNN_v1_ch<NN>_seed.sql`
- HTTP handler: `services/agent-service/internal/transport/httpapi/<concept>_handler.go`
- Route registration: `services/agent-service/internal/transport/httpapi/routes.go`
- Gateway proxy rules: `services/api-gateway/cmd/api-gateway/main.go` (add `srv.Handle(...)` blocks)
- React types: `web/src/types.ts`
- React API functions: `web/src/api.ts`
- React component: `web/src/chapters/vol1/Ch<NN><Title>.tsx`
- Register component: `web/src/chapters/registry.ts` (add entry with `status: "done"`)
- Import in App: `web/src/App.tsx`
- Update roadmap: `docs/architecture.md` (flip `Pending` → `Done`, add `### Ch. NN` section)

**New shared library:**
- Location: `pkg/<name>/` — only for truly cross-cutting concerns used by multiple services

**New service (Vol. III chapter unlocking finance-service):**
- Uncomment the relevant `srv.Handle(...)` block in `services/api-gateway/cmd/api-gateway/main.go`
- Add domain package under `services/finance-service/internal/<pkg>/`
- Add migrations under `services/finance-service/internal/store/migrations/`
- Add entry to `deploy/mysql/init.sql` if the schema doesn't exist yet; run `make mysql-bootstrap`

## Special Directories

**`.planning/codebase/`:**
- Purpose: Codebase map documents consumed by GSD planning tools
- Generated: Yes (by mapping agents)
- Committed: Yes

**`deploy/k8s/`:**
- Purpose: Kustomize manifests for Kubernetes
- Generated: No
- Committed: Yes

**`web/dist/`:**
- Purpose: Vite production build output
- Generated: Yes (`npm run build`)
- Committed: No (gitignored)

**`web/node_modules/`:**
- Purpose: npm dependencies
- Generated: Yes
- Committed: No

**`bin/`:**
- Purpose: Compiled Go binaries (`make build` output)
- Generated: Yes
- Committed: No

**`services/<svc>/internal/store/migrations/`:**
- Purpose: SQL migration files tracked by goose
- Generated: No (hand-authored per chapter)
- Committed: Yes — never delete or edit; only add

---

*Structure analysis: 2026-05-23*
