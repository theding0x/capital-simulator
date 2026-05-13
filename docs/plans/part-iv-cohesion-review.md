# Part IV cohesion + best-practices review (Ch. 12–15)

Status: **plan, not executed.** Saved 2026-05-11 after an audit of the four
chapters of Part IV (Relative Surplus-Value) for cross-chapter cohesion and
Go / React best practices. Execute when credits reset.

## Scope of the review

The four chapters and where they live:

| Ch. | Title                                  | Service           | Domain package(s)                          |
|-----|----------------------------------------|-------------------|--------------------------------------------|
| 12  | The Concept of Relative Surplus-Value  | simulation-engine | `internal/production`                      |
| 13  | Co-operation                            | agent-service     | `internal/agent` (cooperation.go)          |
| 14  | Division of Labour and Manufacture      | agent-service     | `internal/agent` (manufacture.go)          |
| 15  | Machinery and Modern Industry           | simulation-engine | `internal/machinery`, `internal/engine`    |

Marx's argument across the section: relative SV (Ch. 12) requires
productivity gains; Ch. 13–15 enumerate the three historical forms those
gains take — simple co-operation, manufacture, machinery. The code base
**does not currently express that arc.** Each chapter is implemented as an
isolated kingdom with its own type for "productive power", its own
`LabourMinutes`, its own little ID factory.

## Findings, grouped

### A. Cross-chapter cohesion gaps

A1. **`LabourMinutes` is declared five times** (`grep -rn "^type LabourMinutes" services/`):
- `agent/labour_power.go:13`
- `commodity/labour.go:17`
- `machinery/machine.go:21`
- `production/production.go:8`
- `surplus/surplus.go:7`

Three of those live **inside one service** (simulation-engine: surplus,
production, machinery). Within sim-engine that means a Ch. 12 value
cannot be passed to a Ch. 15 function without an explicit cast — even
though they encode the identical Marxist primitive. The comments in
`surplus.go` and `machinery/machine.go` literally apologise ("redeclared
here to avoid an import cycle"). The cycle is artificial.

A2. **No bridge from productivity gain (Ch. 12) to its empirical mechanisms (Ch. 13/14/15).**
`production.ApplyProductivityToSNLT(snlt, pf ProductivityFactor)` accepts
a `ProductivityFactor` but **nothing in the codebase produces one.** Each
of Ch. 13/14/15 carries its own float64 multiplier with a different name
and a different home:

| Ch. | Name                            | Package    | Returns        |
|-----|---------------------------------|------------|----------------|
| 12  | `ProductivityFactor`            | production | float64        |
| 13  | `CollectiveProductivePower`     | agent      | float64        |
| 14  | `ManufactureProductivePowerFactor` | agent   | float64        |
| 15  | `ProductivePower`               | machinery  | int64 (units!) |

So three of the four chapters compute a productivity multiplier, but none
of them feeds Ch. 12's relative-SV calculator. The narrative Marx is
building in Part IV is invisible in the wiring.

A3. **Two `CapitalComposition` structs.**
- `commodity/capital.go:114` — Constant/Variable in LabourMinutes-style.
- `machinery/capital.go:14` — Constant/Variable in pence (int64).
Same name, different units, no relation. Marx's argument is that
mechanisation transforms a *single* category (the organic composition);
the code presents two unrelated types in two services.

A4. **Machine ↔ Manufacture link missing.** Marx (Ch. 15 §1): the
`SpecialisedTool` of manufacture is the seed that becomes the Machine's
`WorkingTool`. `agent.SpecialisedTool` and `machinery.Machine.WorkingTool`
both exist; neither references the other. A Machine cannot be created
*from* a Manufacture's tools, even though that is exactly the §1 story.

A5. **Cooperation ↔ Factory link missing.** Same critique: Marx says
Factory is the most-developed form of Cooperation; the two records have
no FK, no shared interface, no API path that converts one to the other.

A6. **`Pence` is declared only in agent-service** (`agent/agent.go:13`).
`machinery/capital.go` uses bare `int64` "in pence" for its `VariableCapital`
/ `ConstantCapital` types. Inconsistent with the rest of the codebase
and forces every machinery handler to multiply-by-pence in comments instead
of in the type system.

### B. Dead code from the Ch. 15 spec

B1. **`engine.Tick` type — zero callers.** Declared in the Ch. 15 spec,
implemented in `engine/engine.go:21`, but `AdvanceTick` returns an inline
`(factoryID, sequence, occurredAt)` tuple persisted directly into
`factory_ticks` rows. The named type is orphan.

B2. **`engine.CyclePhase` enum — zero callers outside its own test.**
Defined in `engine/engine.go:30`, exercised by one test, exposed by no
endpoint, transitioned by no domain logic. Part of Ch. 15 §7 but never
wired into the simulation.

B3. **`machinery.IntensityFactor` / `EffectiveLabour` — zero callers
outside their own test.** §3C concept. Either wire to a tick path (a
working day intensifies under machinery) or delete.

B4. **`factory_ticks` table is write-only.** `AdvanceTick` inserts a row
each call; no endpoint reads them back; the React panel keeps the
last-tick result in component state only. The persisted history is
inaccessible.

### C. Test-coverage asymmetry

C1. **Agent-service has zero HTTP handler tests.**
- `services/agent-service/internal/transport/httpapi/` — no `_test.go` files.
- `cooperation_handler.go` and `manufacture_handler.go` both do non-trivial
  JSON shaping in `buildCooperationResponse` / `buildManufactureResponse`
  and run pure-function computations in handlers, none of which is
  covered.

Compare to sim-engine, which has `production_handler_test.go` (Ch. 12)
and `machinery_handler_test.go` (Ch. 15). The section is asymmetric: half
its handlers are tested, half are not.

C2. **`engine` package has one trivial test.** Only the `CyclePhase`
validity check. Tick (orphan) has none.

### D. Best-practice issues

D1. **`agent.httpapi.New` takes 7 positional store interfaces.** Easy to
swap by accident. Either accept a config struct, or accept a single
embedded interface (main.go's `agentStore` interface already exists for
exactly this).

D2. **`NewXxxID()` helpers duplicated 10+ times.** Each ID type
re-implements 96-bit hex from `crypto/rand`. Could live as
`pkg/ids.New() string` and be imported by every domain package.

D3. **React format helpers duplicated.**
- `minutesToHours` appears in Ch13 + Ch14, byte-identical.
- `poundsFromPence` appears in Ch13 + Ch14, byte-identical.
- Ch15 has its own `compactNumber` and `poundsFromLabourMinutes`.
- `web/src/format.ts` already exports `fmtMinutes` and `fmtQty` but the
  chapter files don't use them.

D4. **Ch. 15 `createFactoryRequest` accepts both inline machine fields
and bare ID references** in the same `machines: [...]` array. Clever but
undocumented. Two choices:
- Document the dual-mode in the handler comment + OpenAPI.
- Drop the inline path and require `POST /v1/machines` first.

D5. **No "Part IV overview" UI affordance.** The chapter sidebar jumps
12 → 13 → 14 → 15 without surfacing the Marx narrative ("three forms by
which capital realises relative SV"). Optional but on-brand.

## Concrete change list (executable in order)

### P0 — section cohesion (must, before next chapter ships)

**P0-1. Unify `LabourMinutes` within simulation-engine.**
- New file: `services/simulation-engine/internal/labour/labour.go` defining
  `type LabourMinutes int64` (canonical, with one godoc paragraph quoting
  Capital Vol. I).
- Rewrite `surplus`, `production`, `machinery` to import and use
  `labour.LabourMinutes`. Delete the three local declarations.
- Update tests; update store JSON marshalling (the wire format should not
  change — `int64` is the same on the wire).
- **Out of scope** for this step: the agent-service and commodity-service
  declarations. Those are separate services and the IPC boundary
  legitimately re-types. Note this in a follow-up (see P2-1).

Files touched (read-only audit count): ~12.

**P0-2. Introduce `production.ProductivityFactor` as the lingua franca and let Ch. 13/14/15 emit it.**
- Promote `production.ProductivityFactor` (float64) to the new
  `labour` package so every chapter can import it without a cycle.
- Add a tiny adapter per chapter:
  - `agent.CooperationProductivityFactor(coop) ProductivityFactor`
    (wraps `CollectiveProductivePower`).
  - `agent.ManufactureProductivityFactor(m) ProductivityFactor`
    (wraps `ManufactureProductivePowerFactor`).
  - `machinery.FactoryProductivityFactor(f) ProductivityFactor`
    (new — derived from `TotalProductivePower` × hand-labour saved).
- Add `POST /v1/production/relative-surplus-from-productivity` that
  accepts `{working_day, current_lpv, source: "cooperation"|"manufacture"|"factory", source_id}`
  and returns the shortened WorkingDay plus the relative SV delta. The
  endpoint reads the source from its respective store (cross-service
  HTTP call — gateway-mediated). This is the bridge.
- React: add a "feed productivity into Ch. 12" affordance on each of the
  Ch. 13/14/15 panels (a button that opens the Ch. 12 panel pre-filled).

This is the single most important change for narrative cohesion across
Part IV.

**P0-3. Resolve the Ch. 15 orphans — three sub-tasks, one per type.**

- **Tick:** Wire `engine.Tick` as the return shape of `AdvanceTick`. Store
  it in the `factory_ticks` table (which already has the right shape).
  Add `GET /v1/factories/{id}/ticks` (returns `{items: []Tick}`,
  ordered by sequence). The React Ch. 15 panel shows the last 10 ticks
  in a strip-chart-or-table below the factory floor.
- **CyclePhase:** Either (a) implement a minimal phase-transition rule
  on the simulation-engine tick loop (e.g. "after 100 ticks, prosperity →
  overproduction"), expose `GET /v1/sim/phase`, and surface it in the
  Ch. 15 header; or (b) delete the type and its test, mark deferred to
  a future "industrial cycle" chapter. **Recommendation:** (b). The
  industrial cycle is genuinely a later-volume topic; killing the
  orphan is the honest move.
- **IntensityFactor / EffectiveLabour:** Wire into the factory tick: if
  the factory has an associated working day shorter than the §10
  statutory limit, apply `EffectiveLabour` to the labour cost reported
  by the tick. Otherwise delete. **Recommendation:** wire — Marx's §3C
  is short but real, and we already have the type.

**P0-4. Drop the duplicate React format helpers.**
- Extend `web/src/format.ts` with `fmtPounds(pence: number)` and
  `fmtCompact(n: number)`.
- Replace the local copies in Ch13, Ch14, Ch15. Verify Ch12 (no
  helpers currently).
- Mechanical change; zero behaviour delta.

### P1 — section parity (should, while in Part IV)

**P1-1. Add HTTP handler tests for Cooperation and Manufacture.**
- New files:
  - `services/agent-service/internal/transport/httpapi/cooperation_handler_test.go`
  - `services/agent-service/internal/transport/httpapi/manufacture_handler_test.go`
- Mirror the structure of `production_handler_test.go` /
  `machinery_handler_test.go`: spin up `httptest.NewServer`, exercise
  each route with a Marx fixture (Burke 5-man platoon, type-foundry
  4/2/1 proportionality, Birmingham combination, paralysis-by-absence).
- Asserts that JSON shape, error codes, and computed fields match.

**P1-2. Add a `Machine.SourceManufactureID *ManufactureID` field.**
- Captures the Ch. 15 §1 lineage: a Machine often emerges from a
  Manufacture's specialised tool.
- Nullable — many machines arrive in the simulation without a prior
  manufacture (steam-plough has none).
- DB migration: `00003_ch15_machine_source.sql` adds the column;
  Down clause drops it.
- React Ch. 15 panel: optional "from manufacture …" selector on
  the registration form, populated from `api.listManufactures()`.

**P1-3. Add a `Factory.SourceCooperationID *CooperationID` field.**
- Symmetric to P1-2. Captures the Ch. 13 → Ch. 15 succession.

**P1-4. Decide on agent-service handler constructor.**
- Option A: introduce a config struct
  `type Deps struct { Store, CircuitStore, … }` and `New(d Deps)`.
- Option B: define a single composite interface (mirroring the
  `agentStore` interface in `main.go`) and have `New(s agentStore)`
  pull individual sub-stores via embedding.
- **Recommendation:** Option B — fewer parameters, lets the test
  harness pass `store.NewMemory()` once.

### P2 — repo-wide refactors (later, with their own spike)

**P2-1. Lift `LabourMinutes` to `pkg/labour`.**
- After P0-1 has proven the in-service unification, lift the type one
  more level to `pkg/labour` and re-import across all services. The
  IPC boundary still uses `int64` on the wire so no API break.
- Open question: do we want a single `pkg/money` (Pence + Pounds
  conversion helpers) too? Probably yes; combine into one PR.

**P2-2. Consolidate `NewXxxID()` into `pkg/ids`.**
- `pkg/ids.New() string` returning 24-char hex from `crypto/rand`.
- Each `New<X>ID()` becomes `MachineID(ids.New())`. Type-safety
  preserved; the implementation moves.

**P2-3. Reconcile the two `CapitalComposition` structs.**
- They serve different roles: commodity-side decomposes a *product* (c/v/s
  per production run); machinery-side decomposes the *capital advanced*
  (V vs. C across the firm). Marx uses the same term for both.
- Either (a) rename one to disambiguate (`ProductValueComposition` for
  commodity-service, keep `CapitalComposition` in machinery), or
  (b) give them a common interface in `pkg/capital`. **Recommendation:**
  (a); the analytical levels really are different.

### P3 — UI polish

**P3-1. Part IV overview panel.**
- New `web/src/chapters/PartIVOverview.tsx` (or a dashboard-style card
  at the top of Ch.12).
- Shows: current relative-SV rate (Ch. 12), current productivity factor
  (highest among Ch. 13/14/15 sources), capital composition trend
  (Ch. 15 §6 substitution). Useful as a closing-shot for the section.
- Wire under sidebar group "Part IV — Relative Surplus-Value".

**P3-2. Disambiguate "productive power" in the React copy.**
- The four chapters use the phrase with three different meanings.
  Standardise: "output factor" for the float multiplier, "productive
  power" for absolute units/day, "productivity factor" only when feeding
  Ch. 12.

## What is explicitly NOT in this plan

- Ch. 11 (Rate and Mass of Surplus-Value) is Part III; out of scope here
  even though its `surplus` package shares simulation-engine with Ch. 12
  and 15. The `LabourMinutes` unification in P0-1 will touch surplus too,
  but no new endpoints / domain types for Ch. 11.
- Migrating the `agent.WorkingDay` (Ch. 10) and `production.WorkingDay`
  (Ch. 12) into one type. They're co-named but live in different
  services with different invariants; merging is a Ch. 17–20 problem.
- Replacing JSON `int64` LabourMinutes with a typed wire format. The
  current convention is fine.
- Adding a router to the React app. CLAUDE.md anti-pattern.

## Acceptance criteria

When this plan is executed, the following must hold:

- [ ] One `LabourMinutes` type per service (P0-1).
- [ ] At least one end-to-end path where a Ch. 13/14/15 record drives a
      Ch. 12 calculation, with a test that proves it (P0-2).
- [ ] No symbols in `services/simulation-engine/internal/engine` or
      `internal/machinery` that have zero callers outside their own
      test file (P0-3).
- [ ] `factory_ticks` is readable via HTTP, and the React panel shows
      the last N ticks for a factory (P0-3 Tick).
- [ ] No `minutesToHours` / `poundsFromPence` definitions outside
      `web/src/format.ts` (P0-4).
- [ ] `cooperation_handler_test.go` and `manufacture_handler_test.go`
      exist and exercise every route in their handler files (P1-1).
- [ ] `make vet test build` and `npm run lint && npm run build` pass.
- [ ] `docs/architecture.md` Ch. 15 section is updated to reflect the
      new wiring (Tick endpoint, productivity bridge, etc.).
