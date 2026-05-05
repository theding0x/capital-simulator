---
chapter: 18
title: "Various Formula for the Rate of Surplus-Value"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| necessary labour | `NecessaryLabour` | type | `simulation` | `Minutes LabourMinutes` — the paid portion of the working day |
| surplus labour | `SurplusLabour` | type | `simulation` | `Minutes LabourMinutes` — the unpaid portion of the working day |
| working day | `WorkingDay` | type | `simulation` | `TotalMinutes LabourMinutes` — necessary + surplus labour combined |
| variable capital | `VariableCapital` | type | `simulation` | `Pence int64` — money equivalent of necessary labour |
| surplus-value | `SurplusValue` | type | `simulation` | `Pence int64` — money equivalent of surplus labour |
| value of product | `ValueOfProduct` | type | `simulation` | `Pence int64` — newly created value in one working day (variable capital + surplus-value) |
| rate of surplus-value (Formula I) | `FormulaI` | func | `simulation` | `func(s SurplusLabour, v NecessaryLabour) float64` — returns s/v; unbounded above 1.0 |
| rate of surplus-value (Formula II) | `FormulaII` | func | `simulation` | `func(s SurplusLabour, wd WorkingDay) float64` — returns surplus/working-day; always < 1.0 |
| rate of surplus-value (Formula III) | `FormulaIII` | func | `simulation` | `func(unpaid, paid LabourMinutes) float64` — unpaid/paid labour; identical to Formula I |
| rate scenario | `RateScenario` | type | `simulation` | Holds `NecessaryLabour`, `SurplusLabour`; input to all three formula functions |
| rate result | `RateResult` | type | `simulation` | Holds `FormulaI`, `FormulaII`, `FormulaIII float64`; output of `ComputeRates` |
| compute rates | `ComputeRates` | func | `simulation` | `func(s RateScenario) RateResult` — pure function; applies all three formulae |

## Fixtures

- **§ formula I** `"6 hours surplus-labour / 6 hours necessary labour = 100%"` → `FormulaI(SurplusLabour{360}, NecessaryLabour{360}) == 1.0`
- **§ formula II** `"6 hours surplus-labour / 12-hour working-day = 50%"` → `FormulaII(SurplusLabour{360}, WorkingDay{720}) == 0.5`
- **§ formula II bound** Formula II can never reach or exceed 1.0, because surplus-labour is always less than the full working day → `FormulaII(SurplusLabour{719}, WorkingDay{720}) < 1.0`
- **§ formula I unbounded** English agricultural labourer: `"the labourer gets only 1/4, the capitalist ... 3/4 ... this surplus-labour of the English agricultural labourer is to his necessary labour as 3:1, which gives a rate of exploitation of 300%"` → `FormulaI(SurplusLabour{540}, NecessaryLabour{180}) == 3.0`
- **§ formula III** `"paid labour / unpaid labour"` is the same ratio as Formula I with renamed terms → `FormulaIII(SurplusLabour{360}.Minutes, NecessaryLabour{360}.Minutes) == FormulaI(SurplusLabour{360}, NecessaryLabour{360})`

## Invariants

- `scenario.NecessaryLabour.Minutes + scenario.SurplusLabour.Minutes == workingDay.TotalMinutes` — working day partitions entirely into necessary and surplus labour
- `FormulaI(s, v) == FormulaIII(unpaid, paid)` — Formulae I and III are identical in magnitude; only the labels differ (value terms vs. time terms)
- `FormulaII(s, wd) < 1.0` always, because `s.Minutes < wd.TotalMinutes` by definition
- `FormulaI(s, v) > FormulaII(s, wd)` whenever `v.Minutes < wd.TotalMinutes` (i.e., whenever surplus-labour is positive)
- `result.FormulaI == float64(s.Minutes) / float64(v.Minutes)` — Formula I is s/v, not s/(c+v); no constant capital enters

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `NecessaryLabour` — `Minutes LabourMinutes`; the paid portion of the working day
  - `SurplusLabour` — `Minutes LabourMinutes`; the unpaid, appropriated portion
  - `RateScenario` — composite input: `NecessaryLabour` + `SurplusLabour`
  - `RateResult` — output struct holding `FormulaI`, `FormulaII`, `FormulaIII float64`
- New functions:
  - `FormulaI(s SurplusLabour, v NecessaryLabour) float64` — s/v ratio
  - `FormulaII(s SurplusLabour, wd WorkingDay) float64` — s/(s+v) ratio; bounded below 1.0
  - `FormulaIII(unpaid, paid LabourMinutes) float64` — unpaid/paid; same result as FormulaI
  - `ComputeRates(s RateScenario) RateResult` — stateless; applies all three formulae
- New HTTP endpoints:
  - `POST /v1/surplus-value/rates` — stateless; accepts `RateScenario`, returns `RateResult` with all three formula values
- React: add a "Ch. 18 — Rate of Surplus-Value Formulae" panel; three numeric inputs (necessary labour minutes, surplus labour minutes); display FormulaI, FormulaII, FormulaIII side-by-side with a note that FormulaII is always below 100% while FormulaI is unbounded

### Explicitly deferred to later chapters
- Rate of profit (s/C) — Marx explicitly sets it aside; belongs to Book III which is outside this simulation's scope
- Absolute money values of surplus-value and variable capital — the chapter treats them only as proxies for labour-time ratios; pricing belongs to earlier chapters
- Constant capital — the chapter explicitly excludes it from the value-of-product denominator; the simulation follows suit
