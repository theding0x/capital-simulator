---
chapter: 17
title: "Changes of Magnitude in the Price of Labour-Power and in Surplus-Value"
status: proposed
primary_service: agent-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| value of labour-power | `LabourPowerValue` | type | `agent` | `LabourMinutes int64` — the necessary labour-time the agent requires to reproduce itself |
| surplus-value magnitude | `SurplusLabour` | type | `agent` | `LabourMinutes int64` — working day minus necessary labour-time |
| working day | `WorkingDay` | type | `agent` | `TotalMinutes LabourMinutes` — extensive duration of one labour period |
| intensity of labour | `LabourIntensity` | type | `agent` | `Factor float64` — ratio to socially normal intensity (1.0 = normal) |
| productiveness of labour | `LabourProductivity` | type | `agent` | `Factor float64` — multiplicative change; mirrors `commodity.ProductivityChange.Factor` |
| rate of surplus-value | `RateOfSurplusValue` | func | `agent` | `func(surplusLabour, necessaryLabour LabourMinutes) float64` — returns s/v ratio |
| necessary labour-time | `NecessaryLabour` | type | `agent` | `LabourMinutes int64` — portion of working day replacing value of labour-power |
| value created per day | `DailyValueCreated` | func | `agent` | `func(day WorkingDay, intensity LabourIntensity) LabourMinutes` — §2: intensity scales value created |
| surplus-labour time | `SurplusLabourMinutes` | func | `agent` | `func(day WorkingDay, necessary NecessaryLabour) LabourMinutes` — working day minus necessary |
| magnitude scenario | `LabourScenario` | type | `agent` | Holds `WorkingDay`, `NecessaryLabour`, `LabourIntensity`, `LabourProductivity` — input to `ComputeOutcome` |
| scenario outcome | `ScenarioOutcome` | type | `agent` | Holds `DailyValue`, `SurplusLabour`, `LabourPowerValue`, `RateOfSurplusValue float64` |
| compute outcome | `ComputeOutcome` | func | `agent` | `func(s LabourScenario) ScenarioOutcome` — pure function; no side effects |
| law: constant working day | `LawConstantDailyValue` | func | `agent` | `func(day WorkingDay, intensity LabourIntensity, productivityFactor float64) bool` — §1 law 1 |
| inverse relation | `LawInverseRelation` | func | `agent` | `func(before, after ScenarioOutcome) bool` — §1 law 2: surplus-value and value of labour-power move in opposite directions when productivity changes |

## Fixtures

- **§1** `"A working day of 12 hours ... creates the same amount of value ... six shillings"` → `DailyValueCreated(WorkingDay{TotalMinutes:720}, LabourIntensity{Factor:1.0}) == 720` (constant regardless of productivity change)
- **§1** `"If the value of the labour-power fall to 3 shillings, or the necessary labour to 6 hours, the surplus-value will rise to 3 shillings"` → `ComputeOutcome(LabourScenario{WorkingDay:{720}, NecessaryLabour:{360}, LabourIntensity:{1.0}, LabourProductivity:{1.0}}).SurplusLabour == 360`
- **§1** `"the value of the labour-power falls from 4 shillings to 3, i.e. by 1/4 or 25%, the surplus-value rises from 2 shillings to 3, i.e. by 1/2 or 50%"` → proportional changes are asymmetric even though absolute shift is equal: `(4−3)/4.0 == 0.25` and `(3−2)/2.0 == 0.50`
- **§2** `"a working-day of more intense labour is embodied in more products ... value created may be 7, 8, or more shillings"` → `DailyValueCreated(WorkingDay{720}, LabourIntensity{Factor:1.333}) > DailyValueCreated(WorkingDay{720}, LabourIntensity{Factor:1.0})`
- **§3** `"necessary labour time be 6 hours ... working-day be lengthened by 2 hours ... surplus-value increases both absolutely and relatively"` → `ComputeOutcome(LabourScenario{WorkingDay:{840}, NecessaryLabour:{360}, ...}).SurplusLabour == 480`
- **§4A** `"working-day at 12 hours ... value of labour-power rises from 3 shillings to 4 ... if the day be lengthened by 4 hours ... absolute magnitude of surplus-value rises from 3 shillings to 4"` → `ComputeOutcome(LabourScenario{WorkingDay:{960}, NecessaryLabour:{480}, ...}).SurplusLabour == 480` confirming absolute surplus rises 33⅓%

## Invariants

- `scenario.WorkingDay.TotalMinutes == necessaryLabour.Minutes + surplusLabour.Minutes` [§1, §3] — the working day is fully partitioned between necessary and surplus labour
- `DailyValueCreated(day, intensity{1.0}) == day.TotalMinutes` [§1 law 1] — at normal intensity a constant working day creates a constant value regardless of productivity
- `surplusValue + labourPowerValue == dailyValue` [§1 law 2] — since daily value is constant under §1 conditions, any rise in one component requires an equal fall in the other
- `RateOfSurplusValue(s, v) == float64(s) / float64(v)` [§1] — rate of surplus-value is s/v, not s/(c+v); explicitly distinguished from the rate of profit

## Scope

### This chapter builds
- Services: `agent-service`
- New domain types:
  - `WorkingDay` — holds `TotalMinutes LabourMinutes`; the extensive duration of a labour period
  - `NecessaryLabour` — `LabourMinutes`; portion of the working day reproducing the value of labour-power
  - `LabourIntensity` — `Factor float64`; intensive magnitude scaling value created per unit time
  - `LabourProductivity` — `Factor float64`; multiplicative change in productiveness (mirrors `commodity.ProductivityChange`)
  - `LabourScenario` — composite input struct grouping the three variable factors
  - `ScenarioOutcome` — result struct: `DailyValue`, `SurplusLabour`, `LabourPowerValue`, `RateOfSurplusValue float64`
- New HTTP endpoints:
  - `POST /v1/labour-scenarios` — stateless; accepts `LabourScenario`, returns `ScenarioOutcome`; no persistence required
- React: add a "Ch. 17 — Magnitude Changes" panel with a form for the three factors (working-day length, intensity, productivity) and a result table showing daily value, necessary labour, surplus labour, and rate of surplus-value; highlight which law applies to the current combination

### Explicitly deferred to later chapters
- Rate of profit (s/C) — Marx explicitly defers this to Book III; Ch. 17 ends with only the rate of surplus-value (s/v). The rate of profit depends on the organic composition of capital not yet modelled.
- Moral/wear-and-tear deterioration of labour-power at prolonged working days — mentioned in §3 ("wear and tear increases in geometrical progression") but requires a degradation model on `Agent.LabourMinutes` that belongs to Ch. 10 (the working day) extended logic.
- International application of the law of value (intensity differences across nations) — §2 footnote; deferred until a multi-economy simulation tick is introduced.
- Differential rates for women, children, different grades of labour-power — explicitly excluded by Marx in the chapter preamble; deferred indefinitely.
