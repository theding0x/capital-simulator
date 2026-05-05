---
chapter: 24
title: "Conversion of Surplus-Value into Capital"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| accumulation of capital | `Accumulation` | type | `simulation` | `Period int64; CapitalStock CapitalStock; SurplusProduced Pence; AccumulationRate float64; NewConstant Pence; NewVariable Pence` — one extended reproduction cycle |
| accumulation rate | `AccumulationRate` | type | `simulation` | `float64` — fraction of surplus-value reinvested as new capital (0.0 = simple reproduction; 1.0 = full accumulation) |
| additional capital | `AdditionalCapital` | type | `simulation` | `Constant Pence; Variable Pence` — the portion of surplus-value converted into new means of production and new labour-power |
| revenue | `Revenue` | type | `simulation` | `Pence int64` — portion of surplus-value consumed by the capitalist; complement of AdditionalCapital |
| run extended reproduction | `RunExtendedReproduction` | func | `simulation` | `func(initial CapitalStock, surplusRate float64, accumRate float64, periods int64) []Accumulation` — compound growth simulation |
| split surplus | `SplitSurplus` | func | `simulation` | `func(surplus Pence, accumRate float64, compositionRatio float64) AdditionalCapital` — divides surplus into new constant and variable capital per organic composition |
| organic composition ratio | `CompositionRatio` | type | `simulation` | `float64` — fraction of new capital that is constant capital (e.g. 0.8 for 4/5 constant + 1/5 variable) |

## Fixtures

- **§ spinner example §1** `"capital of £10,000; four-fifths (£8,000) in cotton/machinery; one-fifth (£2,000) in wages; surplus-value 100%; surplus product = £2,000"` → `CapitalStock{ConstantCapital:8000, VariableCapital:2000}; SurplusRate:1.0 → SurplusProduced:2000`
- **§ reinvestment §1** `"to convert £2,000 into capital: four-fifths (£1,600) in cotton/machinery, one-fifth (£400) in wages"` → `SplitSurplus(2000, 1.0, 0.8) == AdditionalCapital{Constant:1600, Variable:400}`
- **§ compound §1** second cycle on new capital: `AdditionalCapital{1600,400}` generates surplus of `400 × 1.0 == £400`; `SplitSurplus(400, 1.0, 0.8) == AdditionalCapital{320,80}`
- **§ Abraham begat Isaac** `"£10,000 → £2,000 surplus → £400 new surplus → £80 further surplus"`: first three iterations of `RunExtendedReproduction(CapitalStock{8000,2000}, 1.0, 1.0, 3)` should match these values
- **§ partial accumulation §3** if capitalist consumes half the surplus: `AccumulationRate:0.5`; `SplitSurplus(2000, 0.5, 0.8) == AdditionalCapital{800,200}`; Revenue == 1000

## Invariants

- `SplitSurplus(surplus, rate, ratio).Constant == int64(float64(surplus) * rate * ratio)` — constant portion of additional capital
- `SplitSurplus(surplus, rate, ratio).Variable == int64(float64(surplus) * rate * (1 - ratio))` — variable portion
- `SplitSurplus(surplus, rate, ratio).Constant + SplitSurplus(surplus, rate, ratio).Variable == int64(float64(surplus) * rate)` — additional capital totals the accumulated portion
- Under full accumulation (`rate == 1.0`): `Revenue == 0`; entire surplus-value is converted to capital
- `cycle.CapitalStock.ConstantCapital + cycle.CapitalStock.VariableCapital > prev.CapitalStock.ConstantCapital + prev.CapitalStock.VariableCapital` for any `rate > 0` — capital grows each period

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `AccumulationRate` — `float64`; 0.0 = simple reproduction; 1.0 = all surplus converted to capital
  - `AdditionalCapital` — `Constant, Variable Pence`; the newly formed capital from surplus-value
  - `Revenue` — `Pence int64`; the capitalist's consumed portion; complement of accumulated surplus
  - `CompositionRatio` — `float64`; fraction of new capital that is constant (fixed by technical conditions)
  - `Accumulation` — one period's record: prior stock, surplus produced, split into new capital and revenue
- New functions:
  - `SplitSurplus(surplus Pence, accumRate float64, compositionRatio float64) AdditionalCapital` — pure; splits surplus into constant and variable additional capital
  - `RunExtendedReproduction(initial CapitalStock, surplusRate, accumRate, compositionRatio float64, periods int64) []Accumulation` — compound growth simulation
- New HTTP endpoints:
  - `POST /v1/reproductions/extended` — run extended reproduction scenario; body includes composition ratio and accumulation rate
  - `POST /v1/reproductions/split-surplus` — stateless; compute AdditionalCapital from surplus, rate, and composition
- Migration: `00002_ch24_accumulation_scenarios.sql` — `accumulation_scenarios` table for named scenarios
- React: extend the Ch. 23 reproduction panel or add a "Ch. 24 — Accumulation" panel; inputs for accumulation rate and composition ratio; side-by-side comparison of simple vs. extended reproduction; compound growth chart over N periods

### Explicitly deferred to later chapters
- Concentration vs. centralisation of capital — mentioned in §4 but fully developed in Ch. 25
- Composition of capital (organic composition) changing with accumulation — touched in §4 but the central topic of Ch. 25
- Credit system and joint-stock companies — mentioned but deferred to Book III
- The "so-called labour fund" (§5) — Marx critiques the concept; deferred as the reserve army model is in Ch. 25
