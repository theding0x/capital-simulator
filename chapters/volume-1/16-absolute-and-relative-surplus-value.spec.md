---
chapter: 16
title: "Absolute and Relative Surplus-Value"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| absolute surplus-value | `AbsoluteSurplusValue` | type | `simulation-engine/internal/surplus` | Surplus extracted by prolonging the working day beyond the equivalence point |
| relative surplus-value | `RelativeSurplusValue` | type | `simulation-engine/internal/surplus` | Surplus extracted by shortening necessary labour-time via productivity increase |
| surplus-value (general) | `SurplusValue` | type | `simulation-engine/internal/surplus` | Labour-time appropriated beyond the equivalent of labour-power; wraps `LabourMinutes` with an `Origin` tag |
| working day | `WorkingDay` | type | `simulation-engine/internal/surplus` | Partition: `NecessaryLabour + SurplusLabour == Total`; all fields `LabourMinutes` |
| necessary labour | `NecessaryLabour` | field (`LabourMinutes`) | `simulation-engine/internal/surplus` | Labour reproducing the value of labour-power; lower bound of the working day |
| surplus labour | `SurplusLabour` | field (`LabourMinutes`) | `simulation-engine/internal/surplus` | Labour beyond necessary; substance of surplus-value |
| rate of surplus-value | `RateSurplusValue` | func | `simulation-engine/internal/surplus` | `SurplusLabour / NecessaryLabour`; Marx's "rate of exploitation" |
| rate of profit | `RateOfProfit` | func | `simulation-engine/internal/surplus` | `SurplusValue / TotalCapitalAdvanced`; always < rate of surplus-value when constant capital > 0 |
| productive labour (capitalist) | `ProductiveLabour` | type | `simulation-engine/internal/surplus` | Labour that produces surplus-value for capital; not merely useful labour |
| collective labourer | `CollectiveLabourer` | type | `simulation-engine/internal/surplus` | Aggregate of `agent.ID` members whose combined output constitutes capitalist production |
| formal subjection of labour | `SubjectionKind` const | const (`string`) | `simulation-engine/internal/surplus` | `"formal"` — control without technical transformation of the labour process |
| real subjection of labour | `SubjectionKind` const | const (`string`) | `simulation-engine/internal/surplus` | `"real"` — technical process revolutionised by capital; precondition for relative surplus-value as a general method |
| productivity of labour | `ProductivityFactor` | type (`float64`) | `simulation-engine/internal/surplus` | Multiplier on output per unit time; >1 reduces SNLT and hence `NecessaryLabour` |
| prolongation of working day | `ProlongWorkingDay` | func | `simulation-engine/internal/surplus` | Raises `WorkingDay.Total`, holding `NecessaryLabour` fixed → produces `AbsoluteSurplusValue` |
| reduction of necessary labour | `ReduceNecessaryLabour` | func | `simulation-engine/internal/surplus` | Applies `ProductivityFactor` to lower `NecessaryLabour`, holding `WorkingDay.Total` fixed → produces `RelativeSurplusValue` |

## Fixtures

- **§1** `"The prolongation of the working-day beyond the point at which the labourer would have produced just an equivalent for the value of his labour-power, and the appropriation of that surplus-labour by capital, this is production of absolute surplus-value."` → `ProlongWorkingDay(wd, extraMinutes)` increases `SurplusLabour` by `extraMinutes` with `NecessaryLabour` unchanged; `AbsoluteSurplusValue.LabourMinutes == extraMinutes`.

- **§1** `"In order to prolong the surplus-labour, the necessary labour is shortened by methods whereby the equivalent for the wages is produced in less time."` → `ReduceNecessaryLabour(wd, factor=2.0)` halves `NecessaryLabour`, leaving `WorkingDay.Total` fixed; `RelativeSurplusValue.LabourMinutes == oldNecessaryLabour - newNecessaryLabour`.

- **§1** `"Suppose now such an eastern bread-cutter requires 12 working hours a week for the satisfaction of all his wants ... the honest fellow would perhaps have to work six days a week, in order to appropriate to himself the product of one working day."` → `WorkingDay{Total: 4320, NecessaryLabour: 720, SurplusLabour: 3600}` (minutes); `RateSurplusValue(3600, 720) == 5.0` (500%).

- **§1** Mill's error: `"if the rate of surplus-value be 20%, the rate of profit will be 20:500, i.e., 4% and not 20%"` (£400 constant + £100 variable + £20 surplus) → `RateSurplusValue(surplusLabour=20, necessaryLabour=100) == 0.20`; `RateOfProfit(surplusValue=20, totalCapitalAdvanced=500) == 0.04`. The two functions must return different results for the same surplus magnitude.

- **§1** `"From one standpoint, any distinction between absolute and relative surplus-value appears illusory. Relative surplus-value is absolute, since it compels the absolute prolongation of the working-day beyond the labour-time necessary ... Absolute surplus-value is relative, since it makes necessary such a development of the productiveness of labour."` → both `AbsoluteSurplusValue` and `RelativeSurplusValue` embed a `SurplusValue` with the same `LabourMinutes` arithmetic; `Origin` tag (`"absolute"` / `"relative"`) distinguishes them, not the magnitude type.

## Invariants

- `NecessaryLabour + SurplusLabour == WorkingDay.Total` [§1] — partition invariant; must hold for both `ProlongWorkingDay` and `ReduceNecessaryLabour` paths.
- `AbsoluteSurplusValue.LabourMinutes == WorkingDay.SurplusLabour` [§1] — surplus magnitude equals surplus-labour time in the absolute case.
- `RelativeSurplusValue.LabourMinutes == (oldNecessaryLabour - newNecessaryLabour)` [§1] — surplus gain equals the reduction in necessary labour when total day is fixed.
- `RateOfProfit(sv, totalCapital) < RateSurplusValue(sv, necessaryLabour)` when `totalCapital > necessaryLabour` [§1, Mill critique] — profit rate is always lower than rate of surplus-value because denominator includes constant capital.
- `NecessaryLabour > 0` — precondition for any wage relation; `SurplusLabour` is undefined if the worker produces nothing for themselves.

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `WorkingDay` — total, necessary-labour, and surplus-labour partition; all fields `LabourMinutes`
  - `SurplusValue` — named `LabourMinutes` wrapper; carries `Origin string` tag (`"absolute"` / `"relative"`)
  - `AbsoluteSurplusValue` — result of `ProlongWorkingDay`; embeds `SurplusValue`
  - `RelativeSurplusValue` — result of `ReduceNecessaryLabour`; embeds `SurplusValue`
  - `ProductivityFactor` — `float64`; applied to `NecessaryLabour` via `ReduceNecessaryLabour`
  - `ProductiveLabour` — struct `{ AgentID agent.ID; ProducesSurplusValue bool }` recording whether labour counts as productive in the capitalist sense
  - `CollectiveLabourer` — `[]agent.ID` grouping; placeholder for cooperation chapters
  - `SubjectionKind` — `string` enum `"formal"` / `"real"` recording degree of labour's subsumption to capital
- New HTTP endpoints:
  - `POST /v1/surplus-value/absolute` — compute `AbsoluteSurplusValue` given `working_day_minutes` and `extension_minutes`
  - `POST /v1/surplus-value/relative` — compute `RelativeSurplusValue` given `working_day_minutes`, `necessary_labour_minutes`, and `productivity_factor`
  - `GET  /v1/surplus-value/rate` — return `rate_of_surplus_value` and `rate_of_profit` for supplied inputs; surfaces the Mill distinction
- React: add "Ch. 16 — Absolute and Relative Surplus-Value" panel with a working-day configurator (numeric inputs for `NecessaryLabour` and `WorkingDay.Total`), a side-by-side calculator for absolute vs relative surplus-value, and a rate comparison table showing `rate_of_surplus_value` vs `rate_of_profit`.

### Explicitly deferred to later chapters
- Machinery and large-scale industry as the technical mechanism producing `ProductivityFactor` — Ch. 15; this chapter treats it as an external input only.
- Cooperation as a source of collective surplus-labour — Ch. 11–13; `CollectiveLabourer` is a grouping type only here.
- Accumulation: reinvestment of surplus-value into new capital (M′ → expanded M) — Ch. 23–24.
- Wages as a transformed form of the value of labour-power — Ch. 17–20; `NecessaryLabour` is taken as given rather than derived from wage contracts.
- Decomposition of `TotalCapitalAdvanced` into constant + variable components — Ch. 8–9; the rate-of-profit endpoint accepts the total as a caller-supplied parameter.
