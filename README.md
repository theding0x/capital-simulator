# Capital Simulator

A microservices simulation of an economy, built chapter-by-chapter from
Karl Marx's *Capital, Volume I*. Each chapter of the source text becomes
one feature branch, one PR, and one set of domain types, endpoints, and
React panels. The goal is a formal, runnable model of the categories
Marx develops — commodity, value, money, capital, surplus-value,
accumulation — with the analysis exposed both as HTTP APIs and as a
dashboard you can drive interactively.

It is part code project, part close reading. Where the text gives a
worked numerical example (Marx's 20 yards of linen = 1 coat, the
spinner with £8,000c + £2,000v, the agricultural labourer at 300%
exploitation, Redgrave's 1866 spindle ratios) the example is preserved
verbatim as a test fixture and as a seed row in MySQL, so the dashboard
comes up populated with the same numbers Marx wrote down.

The source text throughout is the **Moore–Aveling English translation
of 1887**, digitised and maintained by the [Marxists Internet
Archive](https://www.marxists.org/archive/marx/works/1867-c1/). All
chapter implementations, test fixtures, and seed rows trace back to
that text. See [Sources and acknowledgements](#sources-and-acknowledgements)
below.

**Status.** Volume I is implemented through the end of Part VII (Ch. 25,
"The General Law of Capitalist Accumulation"). See
[`docs/architecture.md`](docs/architecture.md) for the full chapter
roadmap and a per-chapter summary of what was built.

---

## What it does

A running stack lets you:

- **Compute** — hit stateless endpoints with Marx's formulae: rate of
  surplus-value (three forms), exact piece-wage in farthings, the
  Petty law on prices, the quantity-of-money formula `M = ΣP / V`,
  organic composition `c/(c+v)`, the size of the industrial reserve
  army, repayment-period ceiling division, and many more.
- **Persist** — register commodities, owners, workers, capitalists,
  cooperations, manufactures, factories, wage contracts, national wage
  scales, and accumulation scenarios. Each domain type lives behind a
  Store interface with an in-memory implementation (for tests) and a
  MySQL implementation (for production).
- **Step** — advance a factory one tick (Ch. 15), watch wear and tear
  accumulate by bits across the lifespan, see the moral-depreciation
  ceiling, watch labour displaced per period, and watch the reserve
  army grow as the organic composition rises (Ch. 25).
- **Inspect** — the React dashboard renders one panel per chapter
  ([`web/src/chapters/`](web/src/chapters)). Each panel exposes the
  chapter's inputs, calls the matching service, and shows the result
  in a form that mirrors Marx's argument.

The dashboard is also where the system's seams become visible: the
Part IV "relative surplus-value bridge" component appears inside the
Ch. 13 cooperation inspector, the Ch. 14 manufacture inspector, and
the Ch. 15 factory floor, because cooperation, manufacture, and
machinery are the three forms by which capital realises relative
surplus-value (Marx, Part IV intro).

---

## Architecture

```
                ┌──────────────┐
        browser │     web      │  React 18 + Vite + TypeScript, served by nginx
                └──────┬───────┘
                       │  /api/*
                       ▼
                ┌──────────────┐
                │ api-gateway  │  external HTTP entrypoint, fans out
                └──┬─┬─┬─┬─────┘
                   │ │ │ │
       ┌───────────┘ │ │ └────────────┐
       ▼             ▼ ▼              ▼
┌────────────┐ ┌──────────┐ ┌──────────────────┐
│ commodity- │ │  agent-  │ │ market-service   │
│  service   │ │ service  │ │                  │
└─────┬──────┘ └────┬─────┘ └────────┬─────────┘
      │             │                │
      └─────────────┼────────────────┘
                    │
            ┌───────▼────────────┐
            │ simulation-engine  │  drives ticks
            └───────┬────────────┘
                    │
   ┌────────────────┴───────────────┐
   ▼                                ▼
┌────────┐                       ┌───────┐
│ MySQL  │  durable state        │ Redis │  hot tick state
└────────┘                       └───────┘
```

| Service             | Port | Role                                                                    |
|---------------------|------|-------------------------------------------------------------------------|
| `api-gateway`       | 8080 | External entrypoint; reverse-proxies to domain services.                 |
| `commodity-service` | 8081 | Use-value, exchange-value, value, value-forms, fetishism (Ch. 1, 8, 9).  |
| `market-service`    | 8083 | Exchange, owners, money, circuits, prices (Ch. 2, 3).                    |
| `agent-service`     | 8082 | Workers, capitalists, wage-forms, cooperation, manufacture (Ch. 4–22).   |
| `simulation-engine` | 8084 | Surplus-value mechanics, machinery, reproduction, accumulation (Ch. 11–25). |

All five services share a single root `go.mod` at
`github.com/theding0x/capital-simulator`. Cross-cutting concerns live
under [`pkg/`](pkg):

- `pkg/log` — structured logging via `log/slog`.
- `pkg/httpx` — HTTP server scaffolding with `/healthz`, `/readyz`, graceful shutdown.
- `pkg/mysql` — MySQL driver, connection config, and `Migrate` helper (goose v3).

---

## Stack

Fixed and intentional. Do not propose alternatives in PRs.

- **Go 1.25** monorepo. `database/sql` directly, no ORM. `pressly/goose/v3`
  for migrations. IDs are 96-bit hex via `crypto/rand` — no `google/uuid`.
- **React 18 + Vite + TypeScript** for the dashboard. No router yet:
  one `App.tsx` composes one panel per chapter.
- **MySQL 8** with `utf8mb4_unicode_ci` collation. Schema lives in
  per-service `internal/store/migrations/*.sql` files, embedded via
  `//go:embed` and applied by `pkg/mysql.Migrate` at service startup.
  Migrations are append-only.
- **Redis 7** for hot tick state and caches.
- **Docker** for local dev. **Kubernetes** (kustomize) for the deploy
  manifests under [`deploy/k8s/`](deploy/k8s).

---

## Quick start

### Bring up the whole stack with Docker Compose

```sh
docker compose up --build
```

This builds the five service images and the web image, starts MySQL +
Redis, runs all migrations + seeds at first boot, and exposes:

- Dashboard — `http://localhost:5173`
- API gateway — `http://localhost:8080`
- Individual services — `localhost:8081`–`8084` (for direct probing)
- MySQL — `localhost:3306` (root password `capital`)
- Redis — `localhost:6379`

Each service waits for MySQL's `healthcheck` before starting, so a cold
`docker compose up` takes about 20 seconds to be fully ready.

### Run one service from source

```sh
make run-commodity-service     # or agent-service, market-service, simulation-engine, api-gateway
```

You still need MySQL and Redis running — the easiest path is
`docker compose up mysql redis`, then run the Go service against them
with the env vars from `docker-compose.yml`.

### Build + verify Go

```sh
make vet test build
```

Runs `go vet ./...`, `go test ./...`, and builds every service binary
into `bin/`. The full test suite is required to be green before any PR
merges.

### Web dev loop

```sh
cd web
npm install
npm run dev        # Vite dev server, proxies /api to api-gateway
npm run lint       # tsc typecheck (no eslint — strict TS only)
npm run build      # production build into web/dist
```

### Deploy

```sh
make k8s-apply     # kubectl apply -k deploy/k8s
make k8s-delete    # kubectl delete -k deploy/k8s
```

`deploy/k8s/` is a Kustomize tree with a namespace overlay,
infrastructure overlay (MySQL, Redis), and one overlay per service.
Service images are pulled from `ghcr.io/theding0x/capital-simulator/*`.

---

## Repository layout

```
services/
  api-gateway/         reverse-proxies /v1/* to the appropriate domain service
  commodity-service/   Ch. 1 commodity / value-forms; Ch. 8–9 c, v, s
  agent-service/       Ch. 4–22 — class positions, labour-power, wages, cooperation, manufacture
  market-service/      Ch. 2–3 — exchange, money, circuits, prices
  simulation-engine/   Ch. 11–25 — surplus-value, machinery, reproduction, accumulation

  <svc>/cmd/<svc>/     main.go (one per service)
  <svc>/internal/
    <pkg>/             domain types and pure functions
    store/             store.go interface, memory.go, mysql.go, migrations/
    transport/httpapi/ handlers and routes

pkg/                   shared Go libraries (log, httpx, mysql)

web/                   React 18 + Vite + TypeScript dashboard
  src/chapters/        one panel per chapter (Ch01Commodity.tsx ... Ch25GeneralLaw.tsx)
  src/components/      shared widgets (e.g. RelativeSurplusBridge)
  src/api.ts           the API client (mirror of every backend endpoint)
  src/types.ts         TypeScript mirror of every Go response struct

deploy/k8s/            Kustomize tree for k8s deploys
deploy/mysql/          init.sql + seed.sql consumed by the docker-compose mysql image
docs/architecture.md   authoritative topology + per-chapter summary
docs/plans/            design-direction docs for cross-cutting refactors
CLAUDE.md              contributor conventions, chapter workflow, anti-patterns
```

The chapter **source text** and **specs** are deliberately not in this
repo. They live in a separate Obsidian vault
(`marx-engels/1867/capital-volume-i/{texts,specs}/`) and are read at
implementation time via the `obsidian` MCP server. The text files in
that vault are mirrored from the Marxists Internet Archive's *Capital,
Vol. I* collection at
`https://www.marxists.org/archive/marx/works/1867-c1/`. This keeps the
repo focused on running code, not commentary.

---

## Chapter workflow

Each chapter becomes one branch and one PR, named `volume-X/chapter-Y`.
The full procedure is documented in [`CLAUDE.md`](CLAUDE.md), but the
short form is:

1. Branch off `main` as `volume-1/chapter-NN`.
2. Open a **draft** PR populated from the chapter spec — the *planned*
   services, endpoints, tests, and notes.
3. Implement: new domain types in the relevant service(s), tests using
   Marx's textual examples as fixtures, a `NNNNN_chNN_<slug>.sql`
   schema migration, a paired `NNNNN_chNN_seed.sql` with Marx-faithful
   exemplars (seed IDs prefixed `5eed00000000000000<CC>` so they never
   collide with `crypto/rand` output), and a React panel registered in
   `web/src/chapters/registry.ts`.
4. Flip the chapter's row in `docs/architecture.md` from `Pending` to
   `Done` and add a `### Ch. NN — what was built` section.
5. Commit signed (`commit.gpgsign=true` is set per-repo) with a
   conventional-commit message.
6. Update the draft PR description so it reflects what actually shipped.
7. Wait for GitHub Actions. Fix any failures and push again.
8. When CI is green, mark ready for review.

A few rules from `CLAUDE.md` worth surfacing here:

- **Migrations are append-only.** Never edit existing SQL files — add a
  new numbered file instead.
- **All DDL lives in `.sql` files.** No inline `CREATE TABLE` in Go.
- **No `interface{}` / `any`** in domain structs. Concrete types only.
- **JSON tags are `snake_case`.** TypeScript mirrors match the wire format.
- **`LabourMinutes int64`** is the canonical value-magnitude unit.

---

## Currency conventions

Value-magnitudes inside the simulation are `LabourMinutes` (`int64`)
because that is what Marx's analysis quantifies. Money expressions —
wages, prices, capital stocks — are in **pence** (`Pence int64`), since
Marx's worked examples are in £-s-d and pence is the smallest unit that
keeps everything in integer arithmetic. The dashboard renders in pounds
via `web/src/format.ts` helpers (`fmtPounds`, `fmtLabourMinutes`).

Time-wage and piece-wage chapters (Ch. 20, 21) extend this to
**farthings** (¼d.) where the historical wage tables required quarter-
pence precision. The exact-rational `HourlyPriceOfLabour{Numerator,
Denominator}` type in Ch. 20 preserves Marx's point that the *exact*
hourly price is rarely a round number — only the integer paid wage is.

---

## Where to read next

- [`docs/architecture.md`](docs/architecture.md) — the authoritative
  source of truth for what's built. Includes the topology diagram,
  the service table, the chapter roadmap, and a `### Ch. NN — what was
  built` section for every merged chapter.
- [`CLAUDE.md`](CLAUDE.md) — full contributor playbook: conventions,
  chapter workflow (manual and swarm variants), anti-patterns, where
  to look for X.
- [`docs/plans/`](docs/plans) — design-direction documents for
  cross-cutting refactors (e.g. the Part IV cohesion review that
  unified `LabourMinutes` and added the relative-surplus bridge).
- [`web/src/chapters/registry.ts`](web/src/chapters/registry.ts) — the
  index of chapter panels; the fastest way to see what the dashboard
  currently exposes.

---

## Contributing

Issues and pull requests are welcome. **Interpretation issues are
especially welcome** — if you have read *Capital* and you spot a place
where this code misrepresents Marx's argument, please open an issue.
You do not need to write code or run the project.

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the full contributor
guide, including issue templates, citation conventions, and how
disagreements over interpretation are resolved.

---

## Sources and acknowledgements

Every line of Marx's prose this project quotes, paraphrases, or models
comes from the **Marxists Internet Archive** (MIA) at
[marxists.org](https://www.marxists.org). Specifically:

- **Primary source.** *Capital: A Critique of Political Economy,
  Volume I*, Karl Marx, 1867. English translation by Samuel Moore and
  Edward Aveling, 1887, edited by Frederick Engels.
  [marxists.org/archive/marx/works/1867-c1/](https://www.marxists.org/archive/marx/works/1867-c1/)
- **Chapter numbering** follows the Moore–Aveling edition as served by
  MIA, not the Penguin / Fowkes translation. See
  [`CONTRIBUTING.md`](CONTRIBUTING.md#a-note-before-you-read-further-chapter-numbering)
  for the implications and for the policy on cross-edition citation.

The texts themselves are public domain — Marx died in 1883, Moore in
1922, Aveling in 1898, and the translation has been out of copyright
for over a century. The work for which MIA deserves explicit credit is
the **digitisation, transcription, proofing, indexing, and free
hosting** of these texts, sustained by volunteers since 1990. Without
that work, this project would have to begin by typing *Capital* back in
from a paper edition. It did not.

If you find this project useful, please consider
[donating to MIA](https://www.marxists.org/admin/intro/general/donate.htm)
to support their ongoing transcription and hosting work. They keep the
text available to everyone; this project just builds on top of it.

---

## License

MIT — see [`LICENSE`](LICENSE). Fork, modify, host, and sell the code
freely, provided you preserve the copyright notice.
