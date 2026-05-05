---
chapter: 09
title: "The Rate of Surplus-Value"
status: proposed
primary_service: commodity-service
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| constant capital (c) | `ConstantCapital` | type | `commodity` | `LabourMinutes` alias; value of means of production consumed; transferred, not created |
| variable capital (v) | `VariableCapital` | type | `commodity` | `LabourMinutes` alias; value of labour-power purchased; the only value-creating component |
| surplus-value (s) | `SurplusValue` | type | `commodity` | `LabourMinutes` alias; s = value-product − variable capital |
| rate of surplus-value | `RateOfSurplusValue` | type | `commodity` | `float64`; s/v ratio; the degree of exploitation of labour-power |
| value-product | `ValueProduct` | type | `commodity` | `LabourMinutes`; the new value actually created: v + s (c excluded) |
| capital advanced (C) | `CapitalAdvanced` | type | `commodity` | `LabourMinutes`; C = c + v; total outlay before production |
| expanded capital (C′) | `ExpandedCapital` | type | `commodity` | `LabourMinutes`; C′ = c + v + s; capital after production |
| necessary labour-time | `NecessaryLabour` | type | `commodity` | `LabourMinutes`; portion of working-day that reproduces v |
| surplus-labour-time | `SurplusLabour` | type | `commodity` | `LabourMinutes`; portion of working-day that produces s |
| working-day | `WorkingDay` | type | `commodity` | `LabourMinutes`; NecessaryLabour + SurplusLabour |
| surplus-produce | `SurplusProduce` | type | `commodity` | `Quantity`; fraction of total product embodying surplus-value |
| production account | `ProductionAccount` | type | `commodity` | Aggregates c, v, s for a single production run; all fields `LabourMinutes` |
| rate of surplus-value (func) | `ComputeRate` | func | `commodity` | `func(s SurplusValue, v VariableCapital) (RateOfSurplusValue, error)` |
| value-product (func) | `ComputeValueProduct` | func | `commodity` | `func(v VariableCapital, s SurplusValue) ValueProduct` |
| surplus-produce ratio | `SurplusProduceRatio` | func | `commodity` | `func(surplusLabour, necessaryLabour LabourMinutes) float64`; s/(v+s) as fraction of total product |

## Fixtures

- **§1** `C = £410 const. + £90 var.; surplus-value = £90` → `ProductionAccount{ConstantCapital: 410, VariableCapital: 90, SurplusValue: 90}` is valid; `ComputeRate(90, 90)` returns `1.0` (100%)

- **§1** `rate of surplus-value is not s/C or s/(c+v) but s/v; not 90/500 but 90/90 = 100%` → `ComputeRate(SurplusValue(90), VariableCapital(500))` returns `0.18`; `ComputeRate(SurplusValue(90), VariableCapital(90))` returns `1.0`; these are distinct and only the latter is the rate of surplus-value

- **§1** `s/v = (surplus labour)/(necessary labour)` → given `NecessaryLabour(360)` and `SurplusLabour(360)` (6 h each in a 12 h day), `ComputeRate` applied to both ratios yields identical `RateOfSurplusValue(1.0)` (100%)

- **§1 (spinning mill, 1871)** `£52 var. + £80 surpl.; rate = 80/52 = 153 11/13%` → `ComputeRate(SurplusValue(80), VariableCapital(52))` ≈ `1.5385`; `NecessaryLabour` ≈ 231m, `SurplusLabour` ≈ 369m in a 600-minute (10 h) day

- **§1 (Jacob's wheat, 1815)** `wages £3 10s., surplus £3 11s.; s/v > 100%` → `ComputeRate(SurplusValue(211), VariableCapital(210))` > `1.0` (pence units; 42s = 504d, 43s = 516d; ratio > 1)

- **§4** `surplus-produce = 1/10 of 20 lbs. of yarn = 2 lbs.` → `SurplusProduceRatio(surplusLabour=360, necessaryLabour=360)` returns `0.5`; 2 lbs out of 4 lbs of new-value product (one half of the value-creating portion)

## Invariants

- `account.ValueProduct() == account.VariableCapital + account.SurplusValue` — constant capital merely re-appears; new value is v + s only [§1]

- `account.ExpandedCapital() == account.ConstantCapital + account.VariableCapital + account.SurplusValue` — C′ = c + v + s [§1]

- `ComputeRate(s, v) == float64(s)/float64(v)` and equals `float64(surplusLabour)/float64(necessaryLabour)` — both forms are equivalent [§1]

- `workingDay == necessaryLabour + surplusLabour` — no labour is outside these two categories [§4]

- `SurplusProduceRatio(sl, nl) == float64(sl)/float64(nl+sl)` — surplus-produce fraction equals surplus-labour fraction of the working-day [§4]

## Scope

### This chapter builds
- Services: commodity-service
- New domain types:
  - `ConstantCapital` — `LabourMinutes` alias; value of consumed means of production
  - `VariableCapital` — `LabourMinutes` alias; value of purchased labour-power
  - `SurplusValue` — `LabourMinutes` alias; value newly created beyond variable capital
  - `RateOfSurplusValue` — `float64`; the s/v ratio expressing degree of exploitation
  - `ValueProduct` — `LabourMinutes`; new value created in production: v + s
  - `NecessaryLabour` — `LabourMinutes`; working-time reproducing variable capital
  - `SurplusLabour` — `LabourMinutes`; working-time producing surplus-value
  - `WorkingDay` — `LabourMinutes`; NecessaryLabour + SurplusLabour
  - `ProductionAccount` — struct: c, v, s plus methods `ValueProduct()`, `ExpandedCapital()`, `Rate()`
- New HTTP endpoints:
  - `POST /v1/production-accounts` — record a production run (c, v, s in LabourMinutes); returns full account with rate, value-product, expanded capital
  - `GET /v1/production-accounts` — list recorded accounts
  - `GET /v1/production-accounts/{id}` — fetch one account
  - `POST /v1/rate-of-surplus-value` — stateless: given s and v (LabourMinutes), return rate and equivalent labour-time split
- React: "Ch. 09 — Rate of Surplus-Value" panel; input form for c/v/s (labour-minutes); displays value-product, rate as %, necessary/surplus labour split, surplus-produce fraction; Marx's 1871 spinning-mill and Jacob's 1815 wheat examples as canned fixtures

### Explicitly deferred to later chapters
- Absolute vs. relative surplus-value distinction — Ch. 10 (working-day) and Ch. 12 introduce this split explicitly
- Rate of profit (s/(c+v)) — Book III, Ch. 3; Marx flags this out of scope here
- Organic composition of capital (c/v) — Ch. 25 (general law of capitalist accumulation)
- Valorisation process (how surplus-value arises in production) — Ch. 7–8 must precede Ch. 9 implementation; labour-process domain not yet built
- Division of surplus-value into profit, interest, rent — Books II–III
