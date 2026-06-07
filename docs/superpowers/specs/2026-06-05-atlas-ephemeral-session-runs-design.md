# Atlas Observatory — ephemeral, per-session, in-memory runs

**Date:** 2026-06-05
**Branch:** `feature/atlas-observatory`
**Status:** design — awaiting review

## Goal

The Atlas Observatory simulation should **reset on every page reload** while still
**persisting the data that defines a run (the seed) and the UI's preferences**.

Concretely: when a page loads, the seed is loaded into memory and the run advances
**entirely in memory** — it does **not** write per-tick progress to MySQL. A reload
starts a brand-new run from the seed. Each page load gets its **own** run, isolated
from any other tab.

## Decisions locked during brainstorming

1. **Runtime:** server-side, **per-session**, **in-memory** run. The General Law
   tick math stays in the tested Go domain (`internal/simulation`); we do not port
   it to TypeScript.
2. **Scope:** **Atlas Observatory only** — the orrery field of capitals, the
   aggregate vitals, the hidden abode (General Law), and the immiseration series.
   The other scheduler-driven domains (factories Ch.15, reproduction Ch.20/21,
   piece-prices Ch.19) keep today's MySQL-backed behaviour; they are static CRUD
   exhibits, not a "running simulation".
3. **Persist (durable):** MySQL holds the **seed** + all non-Atlas chapter data;
   `localStorage` holds **UI preferences** (currency toggle, speed, reduced-motion).
4. **Reset:** a page reload mints a new session id, so the run restarts from seed.
   **Levers reset to seed defaults** on reload (a fully clean run).
5. **Advance cadence:** **advance-on-poll** — the run advances only while a tab is
   actively polling. Idle ⇒ zero compute, zero writes.

## Why (UX + cloud cost)

- **UX:** the whole live page resets coherently as one unit; the static exhibits
  stay populated from MySQL (the repo invariant "the dashboard must come up
  populated — empty panels are a regression" still holds). Per-session means one
  tab's reload never yanks another tab's view.
- **Cost:** today the global scheduler runs and writes through to MySQL on every
  pass — billing compute + write-IOPS + unbounded table growth
  (`general_law_periods`, `engine_ticks`) even with zero viewers. The new model
  reads MySQL **once at boot** to build a seed template, advances each session in
  RAM only while it is polled, and never writes per-tick progress. Compute and DB
  write-load for the Atlas feature scale to ~zero at idle.

## Current architecture (what changes)

```
Browser (useSnapshot, polls /v1/observatory/snapshot every 2s)
   └─ Atlas "Run" button → POST /v1/engine/start|stop  (GLOBAL scheduler)
api-gateway  → simulation-engine
simulation-engine:
   engine.Scheduler (one global loop, fixed interval, writes engine_ticks)
     ├─ AccumulationTicker  → store.AccumulateCapital  (MySQL write per pass)
     ├─ GeneralLawTicker    → store.AdvanceAbode        (MySQL write per pass)
     ├─ FactoryTicker / ReproductionTicker / PiecePriceTicker  (out of scope)
   GET /v1/observatory/snapshot → store.FieldSnapshot + store.GetAbodeState +
                                  store.ListGeneralLawPeriods + Scheduler.Status()
   POST /v1/observatory/levers  → store.SetAbodeLevers  (MySQL write)
MySQL: abode_state (singleton row), general_law_periods (grows), industrial_capitals
```

Relevant pure domain (kept as-is, reused):

- `simulation.NewAbodeState()` — the seeded initial abode (mirrors migration 00066/00067).
- `simulation.AdvanceGeneralLaw(s) → (next, period)` — one period of the law.
- `simulation.AbodeState.Readout()` — instantaneous class state.
- `simulation.AbodeState.ApplyLevers(u)` — clamped lever perturbation.
- `circulation.FieldCapital` — per-capital read-model (orrery orbit).
- Orrery growth rule today lives in `store.Memory.AccumulateCapital` (rescale M/P/C
  to the new total, Commodity absorbs rounding); we lift this rule into a pure
  function so a run can advance its field without the store.

## Target architecture

```
Browser
  atlas/session.ts: const atlasSessionId = randomUUID()  // module-scoped, NOT persisted
  useSnapshot → GET /v1/observatory/snapshot?advance=N   header X-Atlas-Session: <id>
                (N = running ? speed : 0)
  levers      → POST /v1/observatory/levers              header X-Atlas-Session: <id>
  localStorage: currency (modern), speed, reduced-motion
api-gateway  → simulation-engine   (existing /v1/observatory/* routes; headers forwarded)
simulation-engine:
  internal/observatory.Manager
    seed template (built once at boot from store.GetAbodeState + store.FieldSnapshot)
    map[sessionID]*Run  (mutex; lastAccess; TTL sweeper; max-session cap)
  GET  snapshot → Manager.GetOrCreate(id).Advance(N).Snapshot()   (no DB I/O)
  POST levers   → Manager.GetOrCreate(id).ApplyLevers(u)          (no DB I/O)
  engine.Scheduler: keeps Factory/Reproduction/PiecePrice tickers ONLY
                    (AccumulationTicker + GeneralLawTicker removed from registration)
MySQL: read ONCE at boot for the seed; no per-tick Atlas writes thereafter
```

### Session identity & lifecycle

- The browser mints `atlasSessionId = crypto.randomUUID()` at module load
  (`web/src/atlas/session.ts`), held in memory only. It is **not** written to
  `localStorage`/`sessionStorage`, so a reload re-runs the module ⇒ new id ⇒ fresh
  run. Fallback to `crypto.getRandomValues` if `randomUUID` is unavailable
  (it requires a secure context; `localhost` qualifies).
- The id is sent on every observatory request as `X-Atlas-Session`.
- Server-side `Manager` keeps `map[string]*Run`. An unknown or evicted id lazily
  creates a fresh run from the seed template — so an evicted session transparently
  restarts from seed (correct: same as a reload).
- A TTL sweeper goroutine evicts runs whose `lastAccess` is older than `sessionTTL`
  (default 15m), plus a `maxSessions` cap (default 500) evicting the oldest. Both
  guard against RAM leaks.

### In-memory Run (`internal/observatory/run.go`)

```go
type Run struct {
    mu         sync.Mutex
    abode      simulation.AbodeState         // holds the live levers too
    field      []circulation.FieldCapital    // the orrery
    periods    []simulation.GeneralLawPeriod // capped to last maxSeries (60)
    tick       int64
    lastAccess time.Time
}
```

- `Advance(n int)` — clamp `n` to `[0, maxAdvancePerPoll]` (10); for each step:
  `abode, period = simulation.AdvanceGeneralLaw(abode)`; append `period` (cap 60);
  `field = advanceField(field, abode.AccumulationRateBP)`; `tick++`.
- `ApplyLevers(u)` — `abode = abode.ApplyLevers(u)`; returns applied (clamped) values.
- `Snapshot()` — projects to the existing `observatorySnapshotResponse` shape:
  `abode.Readout()` + lever scalars + the (capped) `periods` as `law_series`; the
  field as `capitals`; the aggregate `p̄′ = ΣS/ΣC` computed over the field (same
  `rateBP` round-half-up as today); `tick` = `run.tick`. The **handler** (which
  knows the advance count) sets `running` = `advance > 0` and `interval_ms`
  (informational constant) on the response.

### Pure field advance (`internal/observatory/field.go`)

Lifts the existing `AccumulateCapital` rescale rule into a pure function so it has
no store dependency and is unit-testable in isolation:

```go
func advanceField(in []circulation.FieldCapital, alphaBP int64) []circulation.FieldCapital
```

For each capital: `delta = SurplusPence * alphaBP / 10000`; skip if `delta <= 0`;
`newTotal = TotalPence + delta`; rescale `Money/Production/Commodity` proportionally
to `newTotal` with `Commodity` absorbing integer rounding (so the parts sum to
`newTotal`); `CostPrice/Surplus/Status/TurnoverNumber` unchanged. This reproduces
today's orrery motion exactly, just per-session and in RAM.

### Seed template (built once at boot)

In `main.go`, after the store opens:

```go
seedAbode, _ := st.GetAbodeState(ctx)     // the seeded singleton (or NewAbodeState())
seedField, _ := st.FieldSnapshot(ctx)     // the seeded industrial capitals
mgr := observatory.NewManager(seedAbode, seedField, logger)
```

`Manager.GetOrCreate(id)` **deep-copies** the template into a new `Run`
(`AbodeState` is all-scalar so it copies by value; `field` is copied with
`append([]FieldCapital(nil), seedField...)`; `periods` starts empty). MySQL is
touched only here (and by the unrelated CRUD chapters).

### Advance-on-poll (no new routes, no gateway change)

- `GET /v1/observatory/snapshot?advance=N` — read `X-Atlas-Session`; `GetOrCreate`;
  `Advance(N)`; return `Snapshot()`. `advance` defaults to 1, clamped `[0,10]`.
- Run / pause / speed are **client-side**: `N = running ? speed : 0`. The server
  holds the *run* in memory but is stateless about the *running toggle* — so we do
  **not** need a per-session start/stop endpoint (and therefore no new gateway
  line, avoiding the per-path 502 footgun).
- `POST /v1/observatory/levers` — read `X-Atlas-Session`; `GetOrCreate`;
  `ApplyLevers`; echo applied values. No DB write.
- **Reset without reload** (optional UI nicety) is purely client-side: regenerate
  `atlasSessionId`; the next poll hits a fresh run. No backend endpoint.

### MySQL write-path removal

- Remove `AccumulationTicker` and `GeneralLawTicker` from the scheduler
  registration in `main.go`. The remaining tickers are untouched and, as today,
  only run if an operator starts the engine (`SIM_TICK_AUTOSTART` defaults off).
- The store methods `AdvanceAbode`, `SetAbodeLevers`, `AccumulateCapital` and the
  two tickers remain in the codebase (still compiled and unit-tested) but are no
  longer driven by the runtime loop. `abode_state` is read once at boot;
  `general_law_periods` is no longer appended at runtime. **No migrations are added
  or edited** (we introduce no new domain type or table).

### Frontend changes

- `web/src/atlas/session.ts` — module-scoped `atlasSessionId` (+ `getRandomValues`
  fallback). New file.
- `web/src/api.ts` — `getObservatorySnapshot(advance)` and `setObservatoryLevers(u)`
  send the `X-Atlas-Session` header (thread an optional `headers` arg through the
  `http()` helper); snapshot appends `?advance=`.
- `web/src/atlas/useSnapshot.ts` — accept `advance` (derived from running×speed);
  send the header; poll interval comes from persisted prefs (default 2000ms).
  Default `running = true` on a fresh load (the run animates immediately).
- `web/src/atlas/Atlas.tsx` — `toggleRun` becomes pure client state (no
  `startEngine`/`stopEngine` calls); `running` defaults to `true` and is no longer
  synced from `snapshot.running`; `speed` maps to the `advance` count.
- `web/src/CurrencyContext.tsx` — initialise `modern` from `localStorage`, persist
  on toggle.
- `web/src/atlas/TickHeartbeat.tsx` (Transport) — persist `speed`/`reduced` to
  `localStorage`.

### Wire types

The `ObservatorySnapshot` / `AbodeReadout` / `LeverState` shapes are **unchanged**,
so `web/src/types.ts` needs no structural edits. The only behavioural change is the
header + `advance` query param, which are transport concerns, not body shape.

## New / changed files

| File | Change |
|------|--------|
| `services/simulation-engine/internal/observatory/run.go` | **new** — `Run`, `Advance`, `ApplyLevers`, `Snapshot` |
| `services/simulation-engine/internal/observatory/manager.go` | **new** — `Manager`, seed template, get-or-create, TTL + cap eviction |
| `services/simulation-engine/internal/observatory/field.go` | **new** — pure `advanceField` |
| `services/simulation-engine/internal/observatory/*_test.go` | **new** — field/run/manager unit tests |
| `services/simulation-engine/internal/transport/httpapi/observatory_handler.go` | session-aware; advance-on-poll; reads `Manager` not store |
| `services/simulation-engine/internal/transport/httpapi/observatory_levers_handler.go` | session-aware; mutates `Manager` not store |
| `services/simulation-engine/internal/transport/httpapi/handler.go` | `Deps.Observatory *observatory.Manager` |
| `services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go` | rewrite for session/advance |
| `services/simulation-engine/cmd/simulation-engine/main.go` | build seed template + `Manager`; drop 2 tickers from scheduler |
| `web/src/atlas/session.ts` | **new** — ephemeral session id |
| `web/src/api.ts` | header + `advance` on observatory calls |
| `web/src/atlas/useSnapshot.ts` | advance/header/interval |
| `web/src/atlas/Atlas.tsx` | client-side run/pause/speed; drop engine start/stop |
| `web/src/CurrencyContext.tsx` | persist currency |
| `web/src/atlas/TickHeartbeat.tsx` | persist speed/reduced |
| `docs/architecture.md` | note Atlas runs are ephemeral, per-session, in-memory |

## API surface

| Method/Path | Before | After |
|-------------|--------|-------|
| `GET /v1/observatory/snapshot` | reads store + global scheduler status | reads per-session in-memory run; `?advance=N` advances it; `X-Atlas-Session` header |
| `POST /v1/observatory/levers` | `store.SetAbodeLevers` (MySQL write) | mutates per-session in-memory run; `X-Atlas-Session` header |
| `POST /v1/engine/start`,`/stop` | toggled by Atlas Run button | unchanged endpoints, but **Atlas no longer calls them** (operator-only now) |

No new paths ⇒ **no api-gateway change** and **no new migration**.

## Persistence model (where each datum lives)

| Datum | Location | Lifetime |
|-------|----------|----------|
| Seed abode + seeded field | MySQL (read once at boot) | durable |
| Non-Atlas chapter data | MySQL | durable |
| Live run (abode, field, immiseration series, tick, levers) | server RAM, per session | dies with the session / on reload |
| UI prefs (currency, speed, reduced-motion) | browser `localStorage` | durable across reloads |
| Session id | browser memory only | one page load |

## Testing strategy

- `field_test.go` — `advanceField`: totals grow by `α·surplus`; M/P/C proportions
  preserved; parts sum to total; zero/negative surplus is a no-op.
- `run_test.go` — `Advance(n)` grows the field, appends ≤60 capped periods, and
  increments `tick`; `ApplyLevers` clamps and is reflected in the next snapshot;
  `Advance(0)` is a no-op.
- `manager_test.go` — `GetOrCreate` returns a fresh run deep-copied from the
  template (mutating one run leaves the template and other sessions untouched);
  unknown id ⇒ fresh-from-seed; TTL eviction; max-session cap eviction.
- `observatory_handler_test.go` — `X-Atlas-Session` selects/creates a run;
  `advance=N` advances, `advance=0` does not; two different ids get independent
  runs; levers are per-session; list/`capitals` is never `null`; no store writes
  occur (in-memory store, assert `AdvanceAbode` count stays 0).
- Frontend — `npm run lint && npm run build`; Playwright on the booted stack:
  reload yields `period`/`tick` reset to seed; the run animates while "running";
  pause halts advancement; currency toggle survives reload; levers return to seed
  defaults after reload.

## Non-goals

- No change to factory / reproduction / piece-price tickers or any non-Atlas
  domain.
- No save / restore / run-history.
- No shared/multi-tab synchronised run.
- No new tables or migrations; the existing `abode_state` / `general_law_periods`
  tables become seed-read-only at runtime.

## Risks & mitigations

- **Gateway 502 on a new path** — avoided entirely: we add no new path (advance via
  query param on the already-proxied snapshot route). `ReverseProxy` forwards the
  custom header by default.
- **`crypto.randomUUID` secure-context requirement** — `localhost` is a secure
  context; include a `getRandomValues` fallback for safety.
- **Session RAM leak** — TTL sweeper + max-session cap, both unit-tested.
- **Seed staleness** — the template is read once at boot; changing the MySQL seed
  needs a service restart. Acceptable (documented); a refresh hook is YAGNI.
- **Rapid double-polls double-advance** — acceptable for a visualisation; the run
  mutex serialises concurrent polls for one session.

## Done when

`make vet test build` passes · `cd web && npm run lint && npm run build` passes ·
Playwright confirms reload resets the run to seed while currency persists ·
`docs/architecture.md` updated · no MySQL writes occur during a run (verified
against the running stack).
