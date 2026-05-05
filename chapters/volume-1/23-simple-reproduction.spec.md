---
chapter: 23
title: "Simple Reproduction"
status: proposed
primary_service: simulation-engine
---

## Concepts → types

| Marx term | Go identifier | Kind | Package | Notes |
|---|---|---|---|---|
| capital stock | `CapitalStock` | type | `simulation` | `ConstantCapital Pence; VariableCapital Pence` — value composition of capital at a point in time |
| surplus-value fund | `SurplusValueFund` | type | `simulation` | `Total Pence; Revenue Pence; Accumulated Pence` — split of surplus-value between consumption and reinvestment |
| reproduction cycle | `ReproductionCycle` | type | `simulation` | `Period int64; Capital CapitalStock; SurplusRate float64; Fund SurplusValueFund` — one circuit of production |
| original capital | `OriginalCapital` | type | `simulation` | `Pence int64` — the starting value advanced before any surplus-value is produced |
| run simple reproduction | `RunSimpleReproduction` | func | `simulation` | `func(initial CapitalStock, surplusRate float64, periods int64) []ReproductionCycle` — produces cycle-by-cycle snapshot |
| repayment period | `RepaymentPeriod` | func | `simulation` | `func(capital Pence, annualRevenue Pence) int64` — number of periods before original capital is fully replaced by consumed surplus-value |
| individual consumption | `IndividualConsumption` | type | `simulation` | `Pence int64` — capitalist's personal expenditure out of revenue; does not enlarge capital |
| productive consumption | `ProductiveConsumption` | type | `simulation` | `Pence int64` — consumption of means of production and labour-power in production |

## Fixtures

- **§ core example** `"capital of £1,000 begets yearly a surplus-value of £200"` → `RunSimpleReproduction(CapitalStock{800,200}, 1.0, 5)` should show SurplusValueFund.Total == 200 each period; Revenue == 200 (all consumed); Accumulated == 0
- **§ repayment period** `"at the end of 5 years the surplus-value consumed will amount to 5×£200 or the £1,000 originally advanced"` → `RepaymentPeriod(1000, 200) == 5`
- **§ half-consumed** `"if only a part, say one half, were consumed ... the same result would follow at the end of 10 years"` → `RepaymentPeriod(1000, 100) == 10`
- **§ general rule** `"value of capital advanced / surplus-value annually consumed = number of reproduction periods"` → `RepaymentPeriod(C, S) == C / S` for integer-divisible cases
- **§ variable capital reconstituted** `"it is the labourer's own labour, realised in a product, which is advanced to him"` — VariableCapital in each cycle is drawn from surplus-value of prior cycles, not from original owner's fund after the first cycle; simulation records this in `cycle.Capital.VariableCapital` as the reconstituted amount

## Invariants

- `RepaymentPeriod(capital, annualRevenue) == capital / annualRevenue` — integer division; fractional periods round up
- Under simple reproduction: `fund.Accumulated == 0` and `fund.Revenue == fund.Total` each cycle — all surplus-value is consumed as revenue
- `sum(cycle.Fund.Total for all cycles through RepaymentPeriod) >= initial.ConstantCapital + initial.VariableCapital` — after repayment period, total consumed surplus-value equals original capital
- `cycle.Capital.ConstantCapital + cycle.Capital.VariableCapital == prev.ConstantCapital + prev.VariableCapital` — under simple reproduction the capital stock does not grow

## Scope

### This chapter builds
- Services: `simulation-engine`
- New domain types:
  - `CapitalStock` — `ConstantCapital, VariableCapital Pence`; value composition at a snapshot
  - `SurplusValueFund` — `Total, Revenue, Accumulated Pence`; the split of produced surplus-value
  - `ReproductionCycle` — one period's record: capital, surplus fund, period number
  - `OriginalCapital` — `Pence int64`; the historically given starting advance
  - `IndividualConsumption`, `ProductiveConsumption` — tagged Pence types distinguishing the two modes
- New functions:
  - `RunSimpleReproduction(initial CapitalStock, surplusRate float64, periods int64) []ReproductionCycle` — stateless simulation
  - `RepaymentPeriod(capital, annualRevenue Pence) int64` — computes how many periods before original capital is dissolved
- New HTTP endpoints:
  - `POST /v1/reproductions/simple` — run a simple reproduction simulation; body: `{constant_capital, variable_capital, surplus_rate, periods}`; returns array of cycle snapshots
  - `POST /v1/reproductions/repayment-period` — stateless; compute repayment period
- Migration: `00001_ch23_reproduction_cycles.sql` — `reproduction_cycles` table for persisting named scenarios
- React: add a "Ch. 23 — Simple Reproduction" panel; numeric inputs for capital stock and surplus rate; slider for number of periods; table showing each cycle's constant capital, variable capital, surplus-value, revenue consumed, and running total of surplus consumed vs. original capital; highlight when repayment period is reached

### Explicitly deferred to later chapters
- Extended reproduction (accumulation) — deferred to Ch. 24; here all surplus-value is consumed as revenue
- Splitting of surplus-value among profit, interest, rent — explicitly deferred to Book III
- Circulation of capital (M→C→M') — deferred to Book II
- Class-level reproduction of the wage-worker relation — described in the chapter but not modelled as new domain logic; it is a qualitative invariant, not a computable quantity
