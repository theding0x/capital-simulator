---
chapter: 11
title: "Rate and Mass of Surplus Value"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| mass of surplus-value (S) | `SurplusValueMass` | type (`int64`) | `simulation-engine/internal/surplus` | Total surplus-value extracted from all workers simultaneously; in LabourMinutes |
| rate of surplus-value (s/v or a′/a) | `SurplusValueRate` | type (`struct`) | `simulation-engine/internal/surplus` | Ratio of surplus-labour to necessary-labour; numerator and denominator in LabourMinutes |
| variable capital (V) | `VariableCapital` | type (`int64`) | `simulation-engine/internal/surplus` | Total money-value advanced for labour-power across all workers; in Pence to mirror agent.Pence |
| value of one labour-power (v / P) | `LabourPowerValue` | type (`int64`) | `simulation-engine/internal/surplus` | Daily cost of reproducing a single worker; in Pence |
| number of labourers (n) | `WorkerCount` | type (`int`) | `simulation-engine/internal/surplus` | Count of simultaneously employed workers |
| individual surplus-value (s) | `IndividualSurplus` | func | `simulation-engine/internal/surplus` | `func IndividualSurplus(rate SurplusValueRate, v LabourPowerValue) SurplusValueMass` |
| mass formula S = (s/v) × V | `MassByRate` | func | `simulation-engine/internal/surplus` | `func MassByRate(rate SurplusValueRate, totalVariableCapital VariableCapital) SurplusValueMass` |
| mass formula S = P × (a′/a) × n | `MassByWorkers` | func | `simulation-engine/internal/surplus` | `func MassByWorkers(v LabourPowerValue, rate SurplusValueRate, n WorkerCount) SurplusValueMass` |
| minimum capital (threshold for capitalist) | `MinimumCapital` | func | `simulation-engine/internal/surplus` | `func MinimumCapital(v LabourPowerValue, n WorkerCount) VariableCapital` |
| absolute limit of working day (< 24 hrs) | `AbsoluteWorkdayLimit` | const | `simulation-engine/internal/surplus` | `const AbsoluteWorkdayLimit LabourMinutes = 24 * 60` |
| SurplusValueSnapshot (aggregate calc result) | `SurplusValueSnapshot` | type (`struct`) | `simulation-engine/internal/surplus` | Carries Rate, VariableCapital, WorkerCount, Mass; returned by POST endpoint |
| rate expressed as (a′/a) | `Rate()` | method on `SurplusValueRate` | `simulation-engine/internal/surplus` | Returns `float64`; SurplusLabour / NecessaryLabour |

## Fixtures

- **§1** `"if the rate of surplus-value be = 100%, this variable capital of 3s. produces a mass of surplus-value of 3s."` → `MassByRate(SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}, VariableCapital(3)) == SurplusValueMass(3)`

- **§1** `"a variable capital of 300s. will produce a daily surplus-value of 300s."` (100 labourers, rate 100%) → `MassByWorkers(LabourPowerValue(3), SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}, WorkerCount(100)) == SurplusValueMass(300)`

- **§1** `"if the rate of surplus-value doubles … and at the same time variable capital is lessened by half … it yields also a surplus-value of 150s."` (rate 100%→200%, capital 300s→150s, workers 100→50) → `MassByRate(SurplusValueRate{SurplusLabour: 12, NecessaryLabour: 6}, VariableCapital(150)) == SurplusValueMass(300)` *(note: the Marx quote says "150s" but S = (12/6)×150 = 300; the compensation law keeps S constant when rate doubles and workers halve — the expected value is 300, not 150)*

- **§1** `"a variable capital of 1,500s., that employs 500 labourers at a rate of surplus-value of 100% … produces daily a surplus-value of 1,500s."` → `MassByRate(SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}, VariableCapital(1500)) == SurplusValueMass(1500)`

- **§1** `"A capital of 300s. that employs 100 labourers a day with a rate of surplus-value of 200% … produces only a mass of surplus-value of 600s."` → `MassByRate(SurplusValueRate{SurplusLabour: 12, NecessaryLabour: 6}, VariableCapital(300)) == SurplusValueMass(600)`

- **§1 (minimum capital)** `"he would have to employ two labourers in order to live … as well as … a labourer"` → `MinimumCapital(LabourPowerValue(3), WorkerCount(2)) == VariableCapital(6)`

## Invariants

- `MassByRate(rate, V) == MassByWorkers(v, rate, n)` when `V == LabourPowerValue * n` — both formulations of S must agree [§1]
- `IndividualSurplus(rate, v) * int64(n) == int64(MassByWorkers(v, rate, n))` — aggregate mass equals individual surplus scaled by worker count [§1]
- `SurplusValueMass(int64(rate.SurplusLabour) * int64(n)) <= SurplusValueMass(int64(AbsoluteWorkdayLimit) * int64(n))` — total surplus-labour can never exceed the absolute limit of 24 hours times the number of workers [§1]
- `rate.NecessaryLabour + rate.SurplusLabour < int64(AbsoluteWorkdayLimit)` — the full working day (necessary + surplus) must remain strictly less than 24 hours [§1]

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `SurplusValueRate` — ratio of surplus-labour to necessary-labour in LabourMinutes, with a `Rate() float64` method
  - `SurplusValueMass` — aggregate surplus-value extracted from all simultaneously employed workers
  - `VariableCapital` — total money-value advanced for labour-power (Pence)
  - `LabourPowerValue` — daily reproduction cost of one worker (Pence)
  - `WorkerCount` — number of simultaneously employed workers
  - `SurplusValueSnapshot` — result struct returned by the aggregate calculation endpoint
- New HTTP endpoints:
  - `POST /v1/surplus/mass` — accepts `{rate, variable_capital}` or `{labour_power_value, rate, worker_count}`; returns `SurplusValueSnapshot` with computed mass and both formula results for cross-validation
  - `GET /v1/surplus/limits` — returns `AbsoluteWorkdayLimit` and the minimum variable capital given query params `?labour_power_value=&worker_count=`
- React: Add "Ch. 11 — Rate and Mass of Surplus Value" panel in `web/src/App.tsx`; inputs for rate (surplus/necessary labour minutes), variable capital or worker count; displays computed mass and illustrates the compensation law (decrease in workers ↔ rise in rate)

### Explicitly deferred to later chapters
- Relative surplus-value (changing the working day by altering necessary labour) — Ch. 12 introduces it explicitly as the next concept
- Cooperation (large-scale simultaneous employment as a qualitative change in production) — Ch. 13
- Rate of profit vs. rate of surplus-value — Ch. 9 / Book III (rate of profit proper)
- Population growth as the limit to social surplus-value — Ch. 23 (general law of capitalist accumulation)
- Minimum capital threshold as a barrier to entry across branches of production — Ch. 25
