# Technology Stack

**Analysis Date:** 2026-05-23

## Languages

**Primary:**
- Go 1.25.7 — all backend microservices and shared libraries (`pkg/`, `services/`)
- TypeScript 5.6.3 — React frontend (`web/src/`)

**Secondary:**
- SQL — migration files at `services/<svc>/internal/store/migrations/*.sql`
- HCL/YAML — Kubernetes manifests (`deploy/k8s/`) and Docker Compose (`docker-compose.yml`)

## Runtime

**Environment:**
- Go 1.25.7 (specified in `go.mod`; Docker images use `golang:1.25-alpine` build stage)
- Node 20 (specified in `web/Dockerfile`; CI uses `actions/setup-node@v4` with `node-version: '20'`)

**Package Manager:**
- Go: stdlib module system — single `go.mod` at repo root; `github.com/theding0x/capital-simulator`
- Node: npm — `web/package.json`; `package-lock.json` present (CI uses `npm ci`)

**Lockfile:**
- Go: `go.sum` present — do not hand-edit; regenerate with `go mod tidy`
- Node: `web/package-lock.json`

## Frameworks

**Backend HTTP:**
- `pkg/httpx.Server` — internal wrapper around `net/http.ServeMux` with /healthz, /readyz, graceful shutdown, configurable timeouts (`pkg/httpx/server.go`)
- Go 1.22+ mux syntax used throughout: `mux.HandleFunc("POST /v1/...", handler)`

**Frontend:**
- React 18.3.1 — UI component library (`web/src/`)
- Vite 5.4.10 — dev server (port 5173) and production bundler (`web/vite.config.ts`)
- `@vitejs/plugin-react` 4.3.4 — Babel-based fast refresh

**Testing:**
- Go standard `testing` package — no third-party test framework
- `httptest.NewRecorder` for HTTP handler tests

**Build/Dev:**
- `make` — top-level `Makefile` for Go build, test, vet, per-service run targets, Docker, k8s
- `tsc -b` — TypeScript compiler (type-check only; `--noEmit`); aliased as `npm run lint`
- `vite build` — production bundle to `web/dist/`

## Key Dependencies

**Critical Go:**
- `github.com/go-sql-driver/mysql` v1.9.3 — MySQL driver; imported via `database/sql`, no ORM
- `github.com/pressly/goose/v3` v3.27.1 — SQL migration runner; used in `pkg/mysql/migrate.go`; migrations applied on service startup via `//go:embed`

**Indirect Go:**
- `filippo.io/edwards25519` v1.2.0 — transitive via MySQL driver
- `github.com/sethvargo/go-retry` v0.3.0 — transitive via goose
- `golang.org/x/sync` v0.20.0 — transitive via goose
- `go.uber.org/multierr` v1.11.0 — transitive via goose

**Critical Node:**
- `react` ^18.3.1 + `react-dom` ^18.3.1 — UI runtime
- `typescript` ^5.6.3 — type checking
- `vite` ^5.4.10 — dev + build tooling

## Configuration

**Environment (per service):**
- `SERVICE_ADDR` — bind address, e.g. `:8081`
- `LOG_LEVEL` — passed to `pkg/log`; parsed as `slog.Level` (default `info`)
- `MYSQL_DSN` — full DSN e.g. `root:capital@tcp(mysql:3306)/commodity?parseTime=true&loc=UTC`; fallback built by `pkg/mysql.ConfigFromEnv(serviceName)` (`pkg/mysql/client.go`)
- `REDIS_ADDR` — declared in Docker Compose and k8s manifests; not yet read by any Go service code (infra placeholder)
- `MYSQL_DISABLED=true` — bypasses MySQL and uses in-memory store (dev/test)
- `FALLBACK_MEMORY=true` — falls back to in-memory on MySQL connect failure

**Frontend:**
- `API_GATEWAY_URL` + `API_GATEWAY_HOST` — injected into nginx at container start via envsubst (`web/nginx.conf`)
- Dev: Vite proxies `/api` → `http://localhost:8080` (`web/vite.config.ts`)
- Prod: nginx proxies `/api/` → `${API_GATEWAY_URL}/` and strips the `/api` prefix

**Build config files:**
- `web/vite.config.ts` — Vite config (plugins, proxy, outDir, sourcemap)
- `web/tsconfig.json` — TypeScript config (ES2022 target, strict mode, no path aliases)

## Container / Orchestration

**Docker:**
- Each Go service: multi-stage Dockerfile — `golang:1.25-alpine` build → `gcr.io/distroless/static-debian12:nonroot` runtime; build context is repo root
- Web: `node:20-alpine` build → `nginx:1.27-alpine` runtime; build context is `./web`
- Compose file: `docker-compose.yml` — full local stack including MySQL 8, Redis 7-alpine, all six Go services, nginx web

**Kubernetes:**
- Kustomize — `deploy/k8s/kustomization.yaml` with namespace `capital-simulator`
- Per-service YAML in `deploy/k8s/services/`; infra YAML in `deploy/k8s/infra/` (MySQL, Redis)
- Resource limits defined per deployment: 50m–500m CPU, 64–256Mi memory

## Platform Requirements

**Development:**
- Go 1.25+
- Node 20+, npm
- Docker + Docker Compose (for local full-stack)
- GNU Make

**Production:**
- Google Cloud Run (managed platform, per `deploy.yml`)
- Google Artifact Registry (image registry)
- External MySQL 8 and Redis 7 (connection strings provided via Cloud Run env vars)

## Logging

**Library:** `log/slog` (Go stdlib) — wrapped in `pkg/log` (`pkg/log/log.go`)

**Format:** JSON to stdout; every record tagged with `"service": "<name>"`

**Level:** controlled by `LOG_LEVEL` env var; defaults to `info`

---

*Stack analysis: 2026-05-23*
