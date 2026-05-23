# Codebase Concerns

**Analysis Date:** 2026-05-23

---

## HIGH SEVERITY

### Ch. 3 Backend Is Missing — UI Calls 404

**What's wrong:** `docs/architecture.md` marks Ch. 3 (Money: hoarding, means of payment, world money) as Done. The React UI (`web/src/chapters/vol1/Ch03Money.tsx`) renders four panels — C-M-C circuit recorder, hoard panel, payment-obligation tracker, world-money transfer log — and `web/src/api.ts` implements the corresponding calls. None of these endpoints exist on the backend.

**Missing on the server side:**
- `GET/POST /v1/circuits` (C-M-C circuits)
- `GET/POST /v1/hoards` (miser hoarding)
- `GET/POST /v1/payment-obligations`, `POST /v1/payment-obligations/{id}/settle`
- `GET/POST /v1/world-money-transfers`
- `GET /v1/circulation/money-required`

**Evidence:**
- `services/market-service/internal/market/market.go` (268 lines) has `CircuitLeg`, `UniversalEquivalent`, `MoneyCommodity`, `Price` but no `Hoard`, `PaymentObligation`, `WorldMoneyTransfer`, or `Circuit` types.
- `services/market-service/internal/store/store.go` — Store interface has no Ch.3 methods.
- `services/market-service/internal/transport/httpapi/routes.go` — only registers Ch.2 endpoints.
- `services/api-gateway/cmd/api-gateway/main.go` — no proxy rules for `/v1/hoards`, `/v1/circuits`, `/v1/payment-obligations`, `/v1/world-money-transfers`, `/v1/circulation`.

**Impact:** Every API call from the Ch.3 panel returns 404. The dashboard appears populated (the Ch.3 panel renders) but all interactive forms fail silently or with network errors.

**Fix approach:** Add domain types to `services/market-service/internal/market/market.go`, Store interface methods to `services/market-service/internal/store/store.go`, Memory/MySQL implementations, migration, HTTP handler, routes registration, and api-gateway proxy rules. Seed migration needed for hoards and payment obligations.

---

### Vol. II Ch. 1 Has No Implementation — Skipped Test as Anchor Only

**What's wrong:** `docs/architecture.md` marks Vol. II Ch. 1 (The Circuit of Money-Capital) as "In progress" but there is no Go domain code, no store methods, no HTTP handler, no migration, and no React panel for it. The only artifact is a single skipped test:

- `services/simulation-engine/internal/circulation/circulation_test.go:214`:
  ```go
  t.Skip("TODO: implement Vol. II Ch. 1 — The Circuit of Money-Capital")
  ```

**Impact:** The registry entry (`v2-ch01`) shows as `"pending"` and renders the "Not yet implemented" placeholder, which is correct. But the status in `docs/architecture.md` says "In progress" — this is a stale label that will mislead chapter planners.

**Fix approach:** Either implement the chapter (following Vol. II Ch. 2 in `services/simulation-engine/internal/circulation/productive_circuit.go` as a model) or change the `docs/architecture.md` status from "In progress" to "Pending" until a branch is created.

---

### `finance-service` Is Entirely Empty — Vol. III Blocked

**What's wrong:** `finance-service` (port 8085) is the designated home for all of Vol. III (52 chapters: profit, average rate of profit, prices of production, rent, interest, credit, fictitious capital, the trinity formula). The service is scaffolded but has zero domain logic, zero migrations, and an empty Store interface:

- `services/finance-service/internal/store/store.go:27`: `type Store interface{}`
- `services/finance-service/internal/store/memory.go` — 21 lines, scaffolding comments only
- `services/finance-service/internal/store/mysql.go:25` — `NewMySQL` is a no-op; no `//go:embed`, no `Migrate` call
- `services/finance-service/internal/store/migrations/` — empty directory
- `services/finance-service/internal/transport/httpapi/routes.go:12`: `func Register(_ *httpx.Server, _ *Handler) {}`

**Impact:** All commented-out Vol. III proxy routes in `services/api-gateway/cmd/api-gateway/main.go` (lines 214-229) will 404 when uncommented until the service implements them. The `financeProxy` variable is kept live only via `_ = financeProxy`.

**Fix approach:** This is by design and documented as such — each Vol. III chapter PR adds to the empty scaffold. It is not a bug but a new developer must understand the empty state is intentional, not an oversight.

---

### `api-gateway` Info Endpoint Shows Stale Status

**What's wrong:** `services/api-gateway/cmd/api-gateway/main.go:239-248` reports:

```go
"status":  "ch-25-general-law",
"chapter": "Capital Vol. I, Ch. 23 — Simple Reproduction",
```

The codebase is at Vol. II Ch. 3 (33 chapters of Vol. I + 2 chapters of Vol. II done). The status string is 12+ chapters behind and the chapter description is inconsistent with the status key.

**Files:** `services/api-gateway/cmd/api-gateway/main.go:239-248`

**Impact:** Low runtime impact but confuses observability tooling and any consumer of `GET /v1/info` (e.g., health dashboards or CI status scraping).

**Fix approach:** Update `"status"` to `"v2-ch03-commodity-circuit"` and `"chapter"` to `"Capital Vol. II, Ch. 3 — The Circuit of Commodity-Capital"` after each chapter PR merges.

---

## MEDIUM SEVERITY

### `pkg/redis` Does Not Exist — Architecture Diagram Is Wrong

**What's wrong:** `docs/architecture.md` (lines 36, 47-48, 56, 63) shows Redis in the topology diagram, lists `market-service` and `simulation-engine` as "MySQL + Redis", and claims `pkg/redis` exists. None of this is true:

- `pkg/` contains only `httpx/`, `log/`, and `mysql/` — no `redis/` subdirectory.
- `go.mod` has no Redis dependency (neither `go-redis/redis` nor `redis/rueidis`).
- No service file imports or references Redis at runtime.

**Impact:** The architecture document creates a false expectation of caching. Tick state and hot data are not cached. Any future developer who reads the docs will spend time looking for Redis integration that does not exist. The k8s manifest `deploy/k8s/infra/redis.yaml` deploys a Redis pod that nothing connects to.

**Fix approach:** Either implement `pkg/redis` with a real client and wire it into simulation-engine's tick path, or update `docs/architecture.md` to remove Redis from the topology until a chapter requires it. The k8s Redis deployment can stay as infrastructure-forward scaffolding but should be annotated accordingly.

---

### Ch. 23 and Ch. 24 Endpoints Are Stateless — Seed Migrations Were Dropped

**What's wrong:** `docs/architecture.md` for Ch. 23 and Ch. 24 describes persisted tables (`reproduction_cycles`, `accumulation_scenarios`) and seeded fixtures (Marx's £1,000/£200 example). Migration `00026_v1_drop_dead_ch23_ch24_tables.sql` drops both tables with the comment "No Go code reads or writes" these tables.

The seed migrations (`00005_v1_ch23_seed.sql`, `00007_v1_ch24_seed.sql`) still run but insert into tables that no longer exist — this will cause goose to fail at startup unless those seeds were also updated.

**Files:**
- `services/simulation-engine/internal/store/migrations/00026_v1_drop_dead_ch23_ch24_tables.sql`
- `services/simulation-engine/internal/store/migrations/00005_v1_ch23_seed.sql`
- `services/simulation-engine/internal/store/migrations/00007_v1_ch24_seed.sql`

**Impact:** If the seed files reference the now-dropped tables and goose runs them in sequence, the migration will fail. If goose tracks by integer prefix and the seeds already ran before the drop, this is moot in production. Needs verification on a clean schema.

**Fix approach:** Verify that the seed migrations for Ch.23 and Ch.24 have `-- +goose Down` that handles the absence of the tables gracefully, or convert the seed files to no-ops with a comment explaining the table was dropped.

---

### JSON-Stored Slices in MySQL Are Not Queryable

**What's wrong:** Several domain types serialize nested slices directly to JSON columns in MySQL rather than using normalized join tables. This is a pattern used for convenience but violates first normal form and makes SQL-level querying impossible.

**Affected files:**
- `services/agent-service/internal/store/mysql.go:587` — `json.Marshal(lp.Means)` for `MeansOfProduction` in `labour_processes`
- `services/agent-service/internal/store/mysql.go:726` — `json.Marshal(c.Members)` for `CooperationMember` slice in `cooperations`
- `services/agent-service/internal/store/mysql.go:688,692` — `json.Marshal(rs.Sets[0/1].WorkerIDs)` for relay schedule worker IDs

**Impact:** Cannot filter, index, or join on individual members or means of production via SQL. For the current read patterns (always fetch the parent record and unmarshal) this is acceptable. If any future chapter needs "which cooperations include worker X?" or "which labour processes use commodity Y?", this requires a schema migration.

**Fix approach:** When a chapter requires querying into these nested structures, migrate to a proper join table. No immediate action needed.

---

### `services/api-gateway/internal/routes/routes.go` Is a Placeholder

**What's wrong:** The file is empty except for a package declaration and a comment:

```go
// Package routes will hold the HTTP route registration for the api-gateway as
// it grows. Today it is a placeholder so the import path exists from the
// initial scaffold and future chapters can extend it without churn.
package routes
```

All actual routing lives in `services/api-gateway/cmd/api-gateway/main.go`, which is now 260+ lines of `srv.Handle(...)` calls. There is no mechanism forcing new routes to be added in one place vs. the other.

**Files:** `services/api-gateway/internal/routes/routes.go`, `services/api-gateway/cmd/api-gateway/main.go`

**Impact:** `main.go` will continue to grow unboundedly. At the current rate (~6 route pairs per chapter), it will exceed 500 lines within a few Vol. II chapters.

**Fix approach:** Move route registration blocks into `services/api-gateway/internal/routes/routes.go` as a `Register(srv *httpx.Server, proxies map[string]*proxy.Proxy)` function. Or at minimum document the intentional split to avoid confusion for new contributors.

---

### Orphaned TODO Struct Types in `simulation/general_law.go`

**What's wrong:** `services/simulation-engine/internal/simulation/general_law.go` exports four struct types that exist but are never used by any handler, store method, or test:

- `TechnicalComposition` — "TODO(Vol.I Ch.15): activated when machinery chapters extend..." (line 11)
- `Concentration` — "TODO(Vol.III Ch.15/27): activated when credit chapters model..." (line 25)
- `Centralisation` — "TODO(Vol.III Ch.27): activated when credit system chapter models..." (line 42)
- `RelativeSurplusPopulation` — "TODO(Vol.I Ch.25 extension): partition IRA into three strata" (line 52)

These are documented TODOs for future chapters but currently compile only because they're exported and unused types in Go do not cause compiler errors. They contribute to the file's cognitive load.

**File:** `services/simulation-engine/internal/simulation/general_law.go:11-60`

**Impact:** Low. Go compiles them fine. A new developer reading this file will spend time understanding which types are active and which are dormant.

**Fix approach:** Add a `// Dormant: activated in Vol. III Ch. NN` comment block header grouping all deferred types, or move them to a `deferred.go` file in the same package.

---

### `ChapterShell.tsx` Uses `ComponentType<any>` to Bypass TypeScript

**What's wrong:** `web/src/components/ChapterShell.tsx:55`:

```typescript
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type AnyPanel = ComponentType<any>;
```

Every chapter panel is cast to `AnyPanel` when inserted into the `CHAPTER_PANELS` map. This means passing wrong props to any chapter panel will not be caught at compile time.

**File:** `web/src/components/ChapterShell.tsx:54-57`

**Impact:** The `npm run lint` (tsc typecheck) will not catch prop mismatches. Bugs will surface at runtime. The suppression also disables the ESLint `no-explicit-any` rule for this scope.

**Fix approach:** Define a union type `PanelComponent = typeof Ch01Commodity | typeof Ch02Exchange | ...` or use a discriminated prop interface. Either approach is more verbose but type-safe.

---

### 10 `react-hooks/exhaustive-deps` Suppressions

**What's wrong:** Ten `useEffect` calls across Ch.26-Ch.33 and Ch.02-Ch.03 of Vol. II suppress the exhaustive-deps lint rule:

```tsx
// eslint-disable-next-line react-hooks/exhaustive-deps
```

**Files:**
- `web/src/chapters/vol1/Ch26PrimitiveAccumulation.tsx:71`
- `web/src/chapters/vol1/Ch27EnclosureEvents.tsx:54`
- `web/src/chapters/vol1/Ch28BloodyLegislation.tsx:75`
- `web/src/chapters/vol1/Ch29GenesisFarmer.tsx:64`
- `web/src/chapters/vol1/Ch30HomeMarket.tsx:55`
- `web/src/chapters/vol1/Ch31IndustrialCapitalist.tsx:61`
- `web/src/chapters/vol1/Ch32HistoricalTendency.tsx:46`
- `web/src/chapters/vol1/Ch33ModernColonisation.tsx:45`
- `web/src/chapters/vol2/Ch02CircuitProductiveCapital.tsx:42`
- `web/src/chapters/vol2/Ch03CircuitCommodityCapital.tsx:41`

**Impact:** Each suppression hides a potentially stale closure or an infinite-render loop. The `// lint` tool (`npm run lint` = `tsc`) does not check ESLint rules, so these are never reported in CI.

**Fix approach:** Audit each suppression. The typical pattern is `useEffect(() => { load() }, [])` where `load` is defined inside the component — wrap `load` in `useCallback` with the correct deps or move it outside the component.

---

## LOW SEVERITY

### `simulation-engine` MySQL Store Is 1,966 Lines

**What's wrong:** `services/simulation-engine/internal/store/mysql.go` is 1,966 lines, growing by ~80-130 lines per chapter. The `agent-service` MySQL store is 1,400 lines. These files implement dozens of domain concepts sequentially.

**Files:**
- `services/simulation-engine/internal/store/mysql.go` (1,966 lines)
- `services/agent-service/internal/store/mysql.go` (1,400 lines)

**Impact:** Navigation is difficult. Grep is necessary to find any method. Code review for a new chapter's store additions requires scanning the entire file to verify no duplicate method or conflicting transaction.

**Fix approach:** Per-chapter or per-concept MySQL store files (`mysql_machinery.go`, `mysql_circulation.go`) composed via embedding or by splitting the `MySQL` struct across files (Go allows methods on a type in multiple files within a package). This can be done incrementally when adding new chapters.

---

### Agent-Service Has No Missing-Agent Test for Some Handlers

**What's wrong:** Not all handlers in `services/agent-service/internal/transport/httpapi/` test the 404 path for missing agent IDs. For example, `piece_wage_handler_test.go` tests the missing-agent 404 but some handlers in `labour_power_handler.go` don't have corresponding 404 tests.

**Impact:** Low. The store's `ErrNotFound` sentinel is correctly mapped in every handler via `errors.Is`, but the test gap means regressions in error mapping could go undetected.

---

### `TurnoverPlayer` Animates UI Without Backend Time Concept

**What's wrong:** `web/src/components/TurnoverPlayer.tsx` provides play/pause/rewind controls that animate the circuit diagram through nodes on a client-side timer. There is no backend simulation tick that this player reflects. The UI animates `M → M-C → P → C′ → C-M′ → M′` as a purely cosmetic sequence.

**File:** `web/src/components/TurnoverPlayer.tsx`

**Impact:** Cosmetic only. The player does not create a false impression to the user that the economy is "running" — it purely navigates the circuit diagram. However, as Vol. II turnover chapters are implemented, this component will need to be connected to real tick data rather than a client-side interval.

---

### `vol3/` Chapter Directory Is Empty

**What's wrong:** `web/src/chapters/vol3/` exists but contains no component files. All 52 Vol. III chapters show the "Not yet implemented" placeholder, which is correct. The empty directory was pre-created for future use.

**Impact:** None. Expected state.

---

### `CapitalCompositionForm.tsx` Lives at Chapters Root, Not in a Shared Component Directory

**What's wrong:** `web/src/chapters/CapitalCompositionForm.tsx` is a shared form component used by Ch.23 and Ch.24. By convention, shared React components live in `web/src/components/`. Its location at the chapters root breaks this expectation.

**File:** `web/src/chapters/CapitalCompositionForm.tsx`

**Impact:** Low. Discoverable via import trace. No functional issue.

**Fix approach:** Move to `web/src/components/CapitalCompositionForm.tsx` and update two import statements in `Ch23SimpleReproduction.tsx` and `Ch24AccumulationOfCapital.tsx`.

---

### `GET /v1/info` Chapter Description Never Updated After Ch. 25

**What's wrong:** The `handleInfo` function in `services/api-gateway/cmd/api-gateway/main.go` has two stale fields (detailed above under HIGH severity). The description field also only mentions commodity/agent/market/simulation-engine services and omits finance-service from the meaningful routing description.

**File:** `services/api-gateway/cmd/api-gateway/main.go:236-252`

---

## Missing Critical Features

### Vol. II Ch. 1 Backend — Blocks Proper Circuit Triptych

The three circuits (M…M′, P…P, C′…C′) are meant to be the conceptual spine of Vol. II. Ch. 2 (P…P) and Ch. 3 (C′…C′) are implemented. Ch. 1 (M…M′) is not. Until Ch. 1 is implemented, the circuit triptych cannot be displayed in full, and the aggregate/social-capital views across all three forms are incomplete.

**What's missing:** `MoneyCircuit` domain type in `services/simulation-engine/internal/circulation/`, store interface methods, handler, migration, seed, api-gateway proxy rules, React panel `web/src/chapters/vol2/Ch01CircuitMoneyCapital.tsx`.

---

### Ch. 3 (Money) Backend — 5 Endpoint Groups Missing

Detailed above under HIGH severity. The money-functions of Capital Vol. I Ch. 3 (hoarding, means of payment, world money, C-M-C circulation quantity) are depicted in the architecture doc and wired in the React UI but have no server-side implementation.

---

## Test Coverage Gaps

### Vol. II Ch. 1 Test Is a Placeholder

**What's not tested:** The money-circuit (M—C…P…C′—M′) domain — there is a test function anchor but it calls `t.Skip`.

**File:** `services/simulation-engine/internal/circulation/circulation_test.go:201-215`

**Risk:** High — when Ch. 1 is implemented, the test skeleton exists but may be forgotten or the `t.Skip` may accidentally remain.

**Priority:** High — remove the `t.Skip` when Ch. 1 domain code is written.

---

### Ch. 3 Backend Has Zero Tests (Because It Has Zero Code)

**What's not tested:** All Ch. 3 money-function domain code and HTTP handlers.

**Risk:** High — the panel makes API calls that always fail; no tests exist to catch regressions when the backend is implemented.

---

### Finance-Service Has No Domain Tests

**What's not tested:** `services/finance-service/` has no domain files to test (all pending), but also has no test files at all.

**Risk:** Low now, high as Vol. III chapters land. Each chapter PR must add tests or the convention breaks.

---

## Dependencies at Risk

### No Redis Client in `go.mod`

`go.mod` lists only two runtime dependencies: `github.com/go-sql-driver/mysql` and `github.com/pressly/goose/v3`. Redis is in the k8s manifests and architecture diagram but not in go.mod. When any chapter requires Redis caching (turnover state, tick data), this dependency must be added and the `pkg/redis` package must be created from scratch.

**Impact:** Medium. The Redis pod is running in k8s but nothing connects to it.

---

*Concerns audit: 2026-05-23*
