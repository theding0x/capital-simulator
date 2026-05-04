---
chapter: 07
title: "The Labour-Process and the Process of Producing Surplus-Value"
status: proposed
primary_service: agent-service
secondary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| Worker / labourer | `Worker` | type | `agent` | Human agent who sells labour-power; carries `AgentID`, `LabourPowerValue`, `SkillLevel` |
| Capitalist | `Capitalist` | type | `agent` | Agent who purchases labour-power and owns means of production |
| Agent (union type) | `Agent` | type | `agent` | Sum type over `Worker` / `Capitalist`; used in store and HTTP |
| Agent ID | `AgentID` | type | `agent` | 96-bit hex string ID, same pattern as `commodity.ID` |
| New agent ID | `NewAgentID` | func | `agent` | `crypto/rand`-based constructor matching `commodity.NewID()` |
| Labour-power | `LabourPower` | type | `agent` | The capacity to labour; its value = reproduction cost in `LabourMinutes` |
| Value of labour-power | `LabourPowerValue` | type (alias) | `agent` | `LabourMinutes` required to reproduce the worker daily |
| Means of production | `MeansOfProduction` | type | `agent` | Raw materials + instruments consumed in one production run |
| Raw material | `RawMaterial` | type | `agent` | Input commodity consumed by the labour process (reference to `commodity.ID`) |
| Instruments of labour | `Instrument` | type | `agent` | Tools / machinery that transfer value to the product gradually |
| Labour process | `LabourProcess` | type | `agent` | One purposeful act: `WorkerID`, `MeansOfProduction`, `Duration LabourMinutes` → `Product` |
| Product | `Product` | type | `agent` | Output of a `LabourProcess`; new use-value with a derived SNLT |
| Necessary labour-time | `NecessaryLabour` | type (alias) | `agent` | `LabourMinutes` required to reproduce labour-power; equals `LabourPowerValue` per period |
| Surplus labour-time | `SurplusLabour` | func | `agent` | `WorkingDay - NecessaryLabour`; the unpaid portion |
| Working day | `WorkingDay` | type (alias) | `agent` | Total `LabourMinutes` performed in one period |
| Surplus-value | `SurplusValue` | type (alias) | `agent` | `LabourMinutes`; the value produced in surplus labour-time |
| Valorization process | `ValorizationProcess` | type | `agent` | Wraps `LabourProcess`; computes `SurplusValue` by comparing paid vs. total labour |
| Value transferred (constant capital) | `TransferredValue` | func | `agent` | Value of raw materials + wear on instruments transferred to product |
| Value added by living labour | `ValueAdded` | func | `agent` | New value created by `LabourProcess.Duration` of abstract labour |
| Total value of product | `ProductValue` | func | `agent` | `TransferredValue + ValueAdded`; denominated in `LabourMinutes` |
| Production run | `ProductionRun` | type | `simulation-engine/engine` | Engine-level record tying a `ValorizationProcess` to a simulation tick |

## Fixtures

Exact numbers and quotations Marx supplies in Chapter 7. These become test case names and values.

### §1 — The Labour-Process

- **§1.a** "Suppose that in the production of 10 lbs. of yarn … the spinner spins 10 lbs. of cotton, value 10 shillings."
  → `RawMaterial{CommodityID: <cotton>, Quantity: 10}` with SNLT equivalent 10 shillings; product is 10 lbs yarn.

- **§1.b** "The labour-process … ceases when the [product] is finished."
  → `LabourProcess` with a defined `Duration` (total minutes worked); `Product` is returned only on completion.

- **§1.c** Marx's three elementary factors: purposeful activity (labour), subject of labour (raw material), instruments of labour.
  → `LabourProcess` must carry all three: `WorkerID`, `RawMaterial`, `Instrument`.

- **§1.d** "The capitalist … checks that the work is done in a proper manner, and that the means of production are consumed with intelligence, so that no raw material is wasted."
  → `LabourProcess.Validate()` rejects zero-duration runs and nil raw materials.

### §2 — The Production of Surplus-Value

- **§2.a** "We assumed … that the value of a day's labour-power is 3 shillings, and that 6 hours of labour are embodied in that sum; hence that 6 hours are necessary every day to produce the daily means of subsistence of the labourer."
  → `NecessaryLabour = 360` (minutes); `LabourPowerValue = 360 LabourMinutes`.

- **§2.b** "But our man is a capitalist. He … makes [the labourer] work … for 12 hours." The working day is 12 hours (720 minutes).
  → `WorkingDay = 720`; `SurplusLabour(720, 360) == 360`.

- **§2.c** "The value of the yarn … is 15 shillings. Let us now examine this product."
  - Cotton consumed: 10 shillings (= transferred value of raw material).
  - Spindle wear: 2 shillings (= instrument transfer).
  - New labour added: 3 shillings (= 6 hours of spinning, the necessary-labour portion).
  - **Total**: 15 shillings → `ProductValue = TransferredValue(10+2) + ValueAdded(3) = 15`.

- **§2.d** Full valorization: working day extended to 12 hours → new value created = 6 shillings; necessary labour = 3 shillings → `SurplusValue = 6 - 3 = 3 shillings`.
  → `SurplusValue(WorkingDay=720, NecessaryLabour=360, WagePerMinute) == 180 LabourMinutes` (3 hours).

- **§2.e** "The past labour that is embodied in the labour-power … is less than the living labour that it buys." The value of labour-power ≠ value it creates.
  → `LabourPowerValue < ValueAdded` must hold for any profitable `ValorizationProcess`.

- **§2.f** "The secret of the self-expansion of capital resolves itself into having the disposal of a definite quantity of other people's unpaid labour."
  → `ValorizationProcess.SurplusValue()` returns a non-zero `LabourMinutes` whenever `WorkingDay > NecessaryLabour`.

## Invariants

Mathematical laws established in Chapter 7 that tests must enforce:

- `SurplusLabour(wd, nl) == wd - nl` where `wd > nl > 0` [§2]
- `SurplusValue >= 0` for any valid `ValorizationProcess`; `SurplusValue == 0` iff `WorkingDay == NecessaryLabour` [§2]
- `ProductValue(run) == TransferredValue(run) + ValueAdded(run)` [§2, the 15-shillings arithmetic]
- `ValueAdded(run) == LabourMinutes(run.Duration)` — living labour creates value equal to its duration (abstract labour reduces via `AsAbstractLabour`) [§2]
- `TransferredValue(run) >= 0` — constant capital transfers value but does not create new value [§2]
- `LabourPowerValue < ValueAdded(run)` whenever `run.Duration > NecessaryLabour` (the core of surplus-value production) [§2.e]
- `NecessaryLabour + SurplusLabour == WorkingDay` — the working day partitions exhaustively [§2]

## Scope

### This chapter builds

**Services:** `agent-service`, `simulation-engine`

**New domain types** (`services/agent-service/internal/agent/`):

- `AgentID` — 96-bit hex typed string ID; `NewAgentID()` constructor.
- `Worker` — economic subject who sells labour-power: `AgentID`, `LabourPowerValue LabourMinutes`, `SkillLevel int` (uniform = 1 for now), `CreatedAt`, `UpdatedAt`.
- `Capitalist` — economic subject who buys labour-power and means of production: `AgentID`, `CreatedAt`, `UpdatedAt`.
- `LabourPower` — the capacity to labour sold to the capitalist; value (`LabourPowerValue LabourMinutes`) encodes reproduction cost.
- `MeansOfProduction` — composite: `RawMaterials []RawMaterial` + `Instruments []Instrument`.
- `RawMaterial` — reference to a `commodity.ID` plus `Quantity int64` and `SNLTPerUnit LabourMinutes` (snapshot at time of use).
- `Instrument` — reference to a `commodity.ID` plus `WearPerRun LabourMinutes` (value fragment transferred per production run).
- `LabourProcess` — purposeful act tying `WorkerID AgentID`, `CapitalistID AgentID`, `Means MeansOfProduction`, `Duration LabourMinutes` → `Product`.
- `Product` — output use-value: `CommodityKind string`, `Quantity int64`, `TotalValue LabourMinutes` (derived).
- `ValorizationProcess` — wraps `LabourProcess`; exposes `NecessaryLabour()`, `SurplusLabour()`, `SurplusValue()`, `ProductValue()`.

**New domain functions** (`services/agent-service/internal/agent/`):

- `TransferredValue(mp MeansOfProduction) LabourMinutes` — sums raw-material SNLT + instrument wear.
- `ValueAdded(duration LabourMinutes) LabourMinutes` — `AsAbstractLabour`-based; equals `duration` for uniform skill.
- `SurplusLabour(workingDay, necessaryLabour LabourMinutes) LabourMinutes` — `workingDay - necessaryLabour`.
- `SurplusValue(vp ValorizationProcess) LabourMinutes` — living labour created minus labour-power value paid.

**Store** (`services/agent-service/internal/store/`):

- `store.go` — `Store` interface: `CreateWorker`, `GetWorker`, `ListWorkers`, `CreateCapitalist`, `GetCapitalist`.
- `memory.go` — in-memory implementation for tests.
- `mysql.go` — MySQL implementation; `NewMySQL` runs `CREATE TABLE IF NOT EXISTS` for `workers` and `capitalists`.

**New HTTP endpoints** (`services/agent-service/internal/transport/httpapi/`):

- `POST /v1/workers` — register a worker (body: `name`, `labour_power_value`).
- `GET /v1/workers` — list all workers.
- `GET /v1/workers/{id}` — fetch one worker.
- `POST /v1/capitalists` — register a capitalist (body: `name`).
- `GET /v1/capitalists` — list all capitalists.
- `GET /v1/capitalists/{id}` — fetch one capitalist.
- `POST /v1/labour-processes` — run a labour process (body: `worker_id`, `capitalist_id`, `means_of_production`, `duration_minutes`); returns `Product` + `ValorizationProcess` summary.
- `GET /v1/labour-processes/{id}` — fetch a recorded process and its valorization result.

**simulation-engine additions** (`services/simulation-engine/internal/engine/`):

- `ProductionRun` type — records a `ValorizationProcess` result for a given simulation tick (`TickID`, `LabourProcessID`, `SurplusValue LabourMinutes`).
- Engine stub wired to accept a tick advance call; no full scheduler yet (that is Ch. 10+).

**api-gateway:**

- Proxy `/v1/workers/*`, `/v1/capitalists/*`, `/v1/labour-processes/*` → `agent-service:8082`.

**React UI** (`web/src/`):

- Add `Worker`, `Capitalist`, `LabourProcess`, `Product`, `ValorizationResult` to `types.ts`.
- Add Ch. 07 section to `App.tsx`: worker registration form, capitalist registration form, "Run Labour Process" form (pick worker + capitalist, enter duration + means), and a valorization result card showing necessary / surplus labour breakdown.
- Extend `api.ts` with `createWorker`, `listWorkers`, `createCapitalist`, `listCapitalists`, `runLabourProcess`.

### Explicitly deferred to later chapters

- **Rate of surplus-value (s/v)** — Ch. 8-9 introduce constant vs. variable capital formally and compute the rate. The `SurplusValue` magnitude is produced here; the rate ratio waits.
- **Mass of surplus-value across many workers** — relates to cooperation and the scale of exploitation; picked up in Ch. 11.
- **Machinery and relative surplus-value** — changing `NecessaryLabour` via productivity increases; deferred to Ch. 12-15.
- **The working day as a political object** — Ch. 10 models the contested length of `WorkingDay` between capital and labour; here `WorkingDay` is an input constant.
- **Simulation clock / tick scheduler** — the engine `ProductionRun` is recorded but the automatic tick loop is Ch. 10+.
- **Wages as a category** — the form of value of labour-power (time-wages, piece-wages) is Ch. 19-20; here `LabourPowerValue` is simply set at construction time.
- **Redis hot state** — caching valorization results in Redis is deferred until the tick loop is live (Ch. 10+).
