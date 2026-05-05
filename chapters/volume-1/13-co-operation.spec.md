---
chapter: 13
title: "Co-operation"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| co-operation | `Cooperation` | type | `services/agent-service/internal/agent` | A named group of workers assembled by one capitalist to produce together |
| collective working-day | `CollectiveWorkingDay` | type | `services/agent-service/internal/agent` | Sum of individual working-days of all workers in a cooperation |
| worker (wage-labourer) | `Worker` (existing `Class`) | const | `services/agent-service/internal/agent` | Already defined; co-operation is valid only for `Worker` class |
| scale of co-operation | `CooperationSize` | type | `services/agent-service/internal/agent` | `int` count of workers assembled simultaneously under one capitalist |
| social productive power | `CollectiveProductivePower` | func | `services/agent-service/internal/agent` | Returns the output attributed to cooperation, over and above the sum of individual outputs |
| individual working-day | `WorkingDay` | type | `services/agent-service/internal/agent` | Duration in `LabourMinutes` a single worker contributes per period |
| average social labour | `AverageSocialLabour` | func | `services/agent-service/internal/agent` | `CollectiveWorkingDay / CooperationSize`; individual deviations cancel |
| minimum capital for co-operation | `MinimumCapital` | func | `services/agent-service/internal/agent` | Returns the least `Pence` a capitalist must command to assemble `n` workers simultaneously |
| supervision / directing authority | `Supervisor` | type | `services/agent-service/internal/agent` | A `Worker`-class agent with a supervisory role; wage-labourer, not capitalist |
| surplus-value from co-operation | (no new type — existing `SurplusValue` on `CapitalCircuit`) | — | `services/agent-service/internal/agent` | Co-operation raises use-value output without altering value-creation law; SurplusValue accounting is unchanged |
| capitalist command | `Capitalist` (existing `Class`) | const | `services/agent-service/internal/agent` | Co-operation requires a single capitalist directing; already modelled |
| productive power of capital (appearance) | `CooperativeOrigin` | func | `services/agent-service/internal/agent` | Returns `"capital"` — collective power appears as a property of capital, not of labour |

## Fixtures

- **§ (para 2)** `"If a working-day of 12 hours be embodied in six shillings, 1,200 such days will be embodied in 1,200 times 6 shillings."` → `CollectiveWorkingDay(workers=1200, individualMinutes=720) == 864_000` minutes; value scales linearly with worker count when no qualitative co-operation effect is present

- **§ (para 3)** `"the collective working-day of 12 men simultaneously employed, consists of 144 hours; and although the labour of each … may deviate more or less from average social labour … it possesses the qualities of an average social working-day"` → `AverageSocialLabour(collectiveMinutes=8640, n=12) == 720` (i.e. one 12-hour day per worker); individual deviations cancel

- **§ (para 3, Burke footnote 1)** `"any given five men will, in their total, afford a proportion of labour equal to any other five"` → `AverageSocialLabour(collectiveMinutes=3600, n=5) == 720`; minimum platoon size for deviation cancellation is 5

- **§ (para 5)** `"a dozen persons working together will, in their collective working-day of 144 hours, produce far more than twelve isolated men each working 12 hours"` → `CollectiveProductivePower(n=12, individualMinutes=720)` returns a surplus output factor `> 1.0` (qualitative cooperation bonus); isolated output factor is `1.0`

- **§ (para 7)** `"100 men co-operating extend the working-day to 1,200 hours"` → `CollectiveWorkingDay(workers=100, individualMinutes=720) == 72_000` minutes; time-critical tasks are addressed by massing labour at the decisive moment

- **§ (para 8)** `"The payment of 300 workmen at once … requires a greater outlay of capital than does the payment of a smaller number of men, week by week"` → `MinimumCapital(n=300, dailyWagePence=p)` is `300 * p`; for `n=10` it is `10 * p`; ratio 30:1 matches Marx's ratio for constant-capital outlay comparison

## Invariants

- `CollectiveWorkingDay(n, d) == n * d` [para 2] — value addition from co-operating workers is strictly additive from the standpoint of value; no value is created by the cooperative relation itself

- `AverageSocialLabour(CollectiveWorkingDay(n, d), n) == d` [para 3] — the collective working-day divided by n always reduces to the individual average; the social law holds for any n ≥ 1

- `CollectiveProductivePower(n, d) >= n * d` [para 5] — cooperative output of use-values is always ≥ the arithmetic sum of isolated outputs; equality holds only for n=1 (no cooperation)

- `MinimumCapital(n, dailyWage) == n * dailyWage` [para 8] — the capitalist must have at least `n × dailyWage` in pocket before workers assemble; the cooperation constraint is a capital-minimum constraint

## Scope

### This chapter builds
- Services: `agent-service`
- New domain types:
  - `Cooperation` — a named pool of `Worker` agents assembled by one `Capitalist` agent; holds `CooperationSize` and a `WorkingDay` per worker
  - `CollectiveWorkingDay` — value type (`LabourMinutes`) computed as `n × WorkingDay`; represents the aggregate labour contributed per period
  - `WorkingDay` — `LabourMinutes` alias scoped to a single worker's shift duration
  - `CooperationSize` — `int` count of workers in a cooperation; must be ≥ 1
  - `Supervisor` — a `Worker`-class member of a `Cooperation` with a `supervisory bool` flag; wage-labourer, not capitalist
- New HTTP endpoints:
  - `POST /v1/cooperations` — create a cooperation (capitalist ID, list of worker IDs, working-day minutes)
  - `GET /v1/cooperations/{id}` — retrieve a cooperation with its members
  - `GET /v1/cooperations` — list all cooperations; optional `?capitalist_id=` filter
  - `POST /v1/cooperations/{id}/collective-working-day` — compute and return `CollectiveWorkingDay` for the cooperation at current size
  - `POST /v1/cooperations/{id}/average-social-labour` — compute `AverageSocialLabour` and return the per-worker average
- React: add "Ch. 13 — Co-operation" panel; form to create a cooperation (select capitalist, add workers, set working-day); display collective working-day, average social labour, and a note that collective productive power appears as a property of capital

### Explicitly deferred to later chapters
- Division of labour within cooperation — Ch. 14 (Manufacture) takes this up explicitly; Ch. 13 covers only simple co-operation (same or same-kind work)
- Machinery as basis of co-operation — Ch. 15 (Machinery and Modern Industry)
- Rate of surplus-value and constant/variable capital ratios under co-operation — Ch. 8–9 accounting not reopened here
- World-money and monetary minimum capital calculation in non-pence denominations — deferred; Ch. 13 uses existing `Pence` type
- Resistance and counterpressure (class struggle within the cooperative) — Ch. 10 (Working Day) and later; Ch. 13 only names the antagonism
