# External Integrations

**Analysis Date:** 2026-05-23

## APIs & External Services

**None (application is self-contained):**
- No calls to third-party APIs (Stripe, SendGrid, etc.)
- All data originates from internal domain logic and seeded migrations

## Data Storage

**Databases:**
- MySQL 8
  - One schema per service; all schemas on a single MySQL instance
  - Schemas: `commodity`, `market`, `agent`, `simulation`, `finance`
  - Schema creation: `deploy/mysql/init.sql` (runs once on fresh volume; re-apply with `make mysql-bootstrap`)
  - Connection env var: `MYSQL_DSN` — full DSN string e.g. `root:capital@tcp(mysql:3306)/commodity?parseTime=true&loc=UTC`
  - Fallback: `MYSQL_DSN` absent → `pkg/mysql.ConfigFromEnv(service)` builds default DSN pointing at `mysql:3306/<service>`
  - Client: `database/sql` + `github.com/go-sql-driver/mysql` v1.9.3; no ORM
  - Connection pool: max 10 open, 5 idle, 3-minute max lifetime (`pkg/mysql/client.go`)
  - Migrations: `github.com/pressly/goose/v3` — SQL files embedded via `//go:embed` in each service's store package; applied on startup by `pkg/mysql.Migrate`

**File Storage:**
- None — no S3, GCS, or local file storage for domain data

**Caching:**
- Redis 7-alpine — declared in `docker-compose.yml` and k8s infra (`deploy/k8s/infra/redis.yaml`); `REDIS_ADDR` env var provisioned to all services
- Redis client code is **not yet implemented** in any Go service — the infra is scaffolded but no reads/writes occur

## Authentication & Identity

**Auth Provider:** None — no authentication layer implemented

- All API endpoints are open (unauthenticated)
- `deploy.yml` deploys Cloud Run with `--allow-unauthenticated`; a comment notes this should be tightened once the gateway mints OIDC tokens

## Internal Service Communication

**Pattern:** API Gateway reverse-proxy fan-out

The `api-gateway` (port 8080) is the single external entry point. It holds an `httputil.ReverseProxy` per downstream service. All `/v1/*` routes are proxied; there are no message queues or event buses.

**Reverse-proxy targets (api-gateway → domain services):**

| Gateway path prefix(es) | Target service | Port | Env var |
|--------------------------|---------------|------|---------|
| `/v1/commodities`, `/v1/exchange-ratio`, `/v1/capital`, `/v1/production-accounts`, `/v1/rate-of-surplus-value` | commodity-service | 8081 | `COMMODITY_SERVICE_URL` |
| `/v1/owners`, `/v1/offers`, `/v1/exchanges`, `/v1/universal-equivalent`, `/v1/money-commodity`, `/v1/prices` | market-service | 8083 | `MARKET_SERVICE_URL` |
| `/v1/agents`, `/v1/circuit-probes`, `/v1/workers`, `/v1/capitalists`, `/v1/labour-power`, `/v1/labour-processes`, `/v1/working-days`, `/v1/cooperations`, `/v1/manufactures`, `/v1/labour-scenarios`, `/v1/wage-forms`, `/v1/time-wages`, `/v1/piece-price`, `/v1/sub-contracts`, `/v1/intensities`, `/v1/wages`, `/v1/comparisons` | agent-service | 8082 | `AGENT_SERVICE_URL` |
| `/v1/sim/status`, `/v1/surplus/*`, `/v1/production/*`, `/v1/machines`, `/v1/factories`, `/v1/surplus-value/*`, `/v1/reproductions/*`, `/v1/accumulation`, `/v1/historical-stages`, `/v1/enclosure-events`, `/v1/statutory-wages`, `/v1/farm-tenures`, `/v1/market-formation`, `/v1/colonial-markets`, `/v1/commodity-circuits`, `/v1/productive-circuits` | simulation-engine | 8084 | `SIM_ENGINE_URL` |
| (Vol. III routes — commented out pending chapter PRs) | finance-service | 8085 | `FINANCE_SERVICE_URL` |

Proxy implementation: `services/api-gateway/internal/proxy/` using `net/http/httputil.ReverseProxy`.
Route registration: `services/api-gateway/cmd/api-gateway/main.go`.

**Direct service-to-service HTTP call (simulation-engine → agent-service):**
- `simulation-engine` calls `agent-service` directly (not through gateway) to fetch cooperation and manufacture productivity factors
- Implementation: `services/simulation-engine/internal/productivity/fetcher.go`
- Client: `net/http.Client` with 5-second timeout
- Paths called: `GET /v1/cooperations/{id}`, `GET /v1/manufactures/{id}`
- Config: `AGENT_SERVICE_URL` env var (default `http://agent-service:8082`)

## Service Discovery

**Mechanism:** Environment variables + DNS (Docker Compose service names / k8s ClusterIP names)

Each service's URL is injected as an env var. No service registry or sidecar.

**Default URLs (used when env var is absent):**

| Env Var | Default |
|---------|---------|
| `COMMODITY_SERVICE_URL` | `http://commodity-service:8081` |
| `MARKET_SERVICE_URL` | `http://market-service:8083` |
| `AGENT_SERVICE_URL` | `http://agent-service:8082` |
| `SIM_ENGINE_URL` | `http://simulation-engine:8084` |
| `FINANCE_SERVICE_URL` | `http://finance-service:8085` |
| `MYSQL_DSN` | `root:capital@tcp(mysql:3306)/<service>?parseTime=true&loc=UTC` |
| `REDIS_ADDR` | `redis:6379` (env set; not yet read by Go code) |

## Monitoring & Observability

**Health endpoints:** Every service exposes `/healthz` (liveness) and `/readyz` (readiness) via `pkg/httpx.Server` (`pkg/httpx/server.go`)

**Kubernetes probes:**
- Readiness: `GET /readyz` — `initialDelaySeconds: 3`, `periodSeconds: 5`
- Liveness: `GET /healthz` — `initialDelaySeconds: 10`, `periodSeconds: 20`

**Error Tracking:** None — no Sentry, Datadog, or similar

**Metrics:** None — no Prometheus, OpenTelemetry, or similar

**Logs:** JSON structured logs to stdout via `log/slog` (`pkg/log/log.go`); consumed by Cloud Run's log aggregation

## CI/CD & Deployment

**CI Pipeline:** GitHub Actions

- Trigger: push/PR to `main`
- Config: `.github/workflows/ci.yml`
- Go job: `go mod tidy` (clean check), `go vet ./...`, `go build ./...`, `go test ./...`
- Web job: `npm ci`, `npm run lint` (tsc typecheck), `npm run build`
- Go version: read from `go.mod` via `actions/setup-go@v5`
- Node version: 20 via `actions/setup-node@v4`

**Deployment:** Google Cloud Run

- Config: `.github/workflows/deploy.yml`
- Trigger: push to `main` or manual `workflow_dispatch`
- Change detection: `dorny/paths-filter@v3` — only services with changed files are rebuilt/deployed
- Registry: Google Artifact Registry (`GCP_REGION-docker.pkg.dev/GCP_PROJECT_ID/AR_REPOSITORY/<service>`)
- Auth: Workload Identity Federation (`GCP_WIF_PROVIDER`, `GCP_DEPLOY_SA` secrets)
- Image builder: `docker/build-push-action@v6` with GitHub Actions cache (`type=gha`)
- Deployer: `google-github-actions/deploy-cloudrun@v2`
- Runtime config (DSNs, secrets): managed out-of-band on Cloud Run services; not set by the deploy workflow

## Frontend ↔ Backend Communication

**Protocol:** HTTP/JSON — REST; no GraphQL, WebSockets, or gRPC

**Path:** Browser → nginx `/api/*` → api-gateway `/v1/*` → domain service

- Dev: Vite dev server proxies `/api` → `http://localhost:8080` (strips `/api` prefix via `changeOrigin`)
- Prod: nginx envsubst template proxies `/api/` → `${API_GATEWAY_URL}/` (`web/nginx.conf`)

**React API client:** `web/src/api.ts`
- Base path: `const BASE = "/api"` (line 181)
- All calls: `fetch(\`${BASE}${path}\`, ...)` — no third-party HTTP client

## Webhooks & Callbacks

**Incoming:** None

**Outgoing:** None

## External Tooling Integrations

**Obsidian MCP server (development tooling only):**
- Used by Claude agents to read chapter source text and specs from the `red-vault` Obsidian vault
- Server runs at `127.0.0.1:27300`
- Not part of the running application; development workflow tool only
- Paths: `marx-engels/<year>/capital-volume-<roman>/{texts,specs}/`

**GitHub Container Registry (ghcr.io):**
- Docker image tag target: `ghcr.io/theding0x/capital-simulator/<service>:dev` (local builds via `make docker`)
- Production images pushed to Google Artifact Registry (see CI/CD above)

---

*Integration audit: 2026-05-23*
